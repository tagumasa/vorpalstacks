package sfn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	types "vorpalstacks/internal/common/tags"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// This file holds the tag-operation Core methods shared by the HTTP API and
// the admin console handler. Both entry points converge on the same
// validation, quota enforcement and persistence logic, so behaviour cannot
// drift between the two surfaces.

// enforceTagQuota verifies that applying the given tags keeps the resource
// within the fifty-tags-per-resource quota, counting the keys already
// present on the resource.
func enforceTagQuota(store *sfnstore.StepFunctionStore, arn string, tags []types.Tag) error {
	existing, err := store.ListAsSlice(arn)
	if err != nil {
		return err
	}
	keys := make(map[string]struct{}, len(existing)+len(tags))
	for _, t := range existing {
		keys[t.Key] = struct{}{}
	}
	for _, t := range tags {
		keys[t.Key] = struct{}{}
	}
	if len(keys) > sfnstore.MaxTagsPerResource {
		return NewTooManyTags(fmt.Sprintf("Too many tags: %d, maximum allowed %d", len(keys), sfnstore.MaxTagsPerResource))
	}
	return nil
}

// tagResourceCore validates the resource, enforces the tag quota and applies
// the given tags. It is the single mutation path for TagResource.
func (s *StepFunctionService) tagResourceCore(ctx context.Context, store *sfnstore.StepFunctionStore, arn string, tags []types.Tag) error {
	if err := validateTaggableResource(ctx, store, arn); err != nil {
		return err
	}
	if err := enforceTagQuota(store, arn, tags); err != nil {
		return err
	}
	return store.TagFromSlice(arn, tags)
}

// untagResourceCore validates the resource and removes the given tag keys.
func (s *StepFunctionService) untagResourceCore(ctx context.Context, store *sfnstore.StepFunctionStore, arn string, tagKeys []string) error {
	if err := validateTaggableResource(ctx, store, arn); err != nil {
		return err
	}
	return store.Untag(arn, tagKeys)
}

// listTagsForResourceCore validates the resource and returns its tags.
func (s *StepFunctionService) listTagsForResourceCore(ctx context.Context, store *sfnstore.StepFunctionStore, arn string) ([]types.Tag, error) {
	if err := validateTaggableResource(ctx, store, arn); err != nil {
		return nil, err
	}
	return store.ListAsSlice(arn)
}

// validateTaggableResource probes the resource identified by an ARN for a
// tag operation. Tagging applies to state machines and activities only
// (TagResource API reference: resourceArn is "the Step Functions state
// machine or activity"); version and alias ARNs extend the state machine
// ARN with a qualifier and map runs, executions and any non-States ARN
// have no tag store behind them — accepting them would persist phantom
// tag records, so they are rejected with ResourceNotFound.
func validateTaggableResource(ctx context.Context, store *sfnstore.StepFunctionStore, arn string) error {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	switch {
	case strings.HasPrefix(resource, "stateMachine:"):
		if strings.Contains(strings.TrimPrefix(resource, "stateMachine:"), ":") {
			// A qualified ARN (version or alias) does not name a taggable
			// resource.
			return NewResourceNotFound("Resource does not exist: " + arn)
		}
		if _, err := store.GetStateMachine(ctx, arn); err != nil {
			if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
				return NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
			}
			return err
		}
	case strings.HasPrefix(resource, "activity:"):
		if _, err := store.GetActivity(ctx, arn); err != nil {
			if errors.Is(err, sfnstore.ErrActivityNotFound) {
				return NewActivityDoesNotExist("Activity Does not exist: " + arn)
			}
			return err
		}
	default:
		return NewResourceNotFound("Resource does not exist: " + arn)
	}
	return nil
}

// tagHandlerConfig wires the generic tag handler framework to the Core
// path: the parse and format closures only translate wire parameters,
// while validation, quota enforcement and persistence run through the
// same Core functions the admin console handler calls.
func (s *StepFunctionService) tagHandlerConfig(store *sfnstore.StepFunctionStore) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:   "resourceArn",
			TagsParam:       "tags",
			TagKeysParam:    "tagKeys",
			RequireTags:     true,
			RequireTagKeys:  true,
			RequireResource: true,
		},
		ResourceKey: func(rawKey string) string { return rawKey },
		ValidateResource: func(ctx context.Context, arn string) error {
			return validateTaggableResource(ctx, store, arn)
		},
		ParseTags: func(params map[string]interface{}) []tagutil.Tag {
			return tagutil.MapToTags(tagutil.ToMap(tagutil.ParseTags(params, "tags")))
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return tagutil.ParseTagKeysAsSlice(params, "tagKeys")
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagSlice []tagutil.Tag) error {
			return s.tagResourceCore(ctx, store, resourceKey, tagSlice)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return s.untagResourceCore(ctx, store, resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return s.listTagsForResourceCore(ctx, store, resourceKey)
		},
		FormatResponse: func(tagSlice []tagutil.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"tags": tagutil.ToResponseWithKeyNames(tagSlice, "key", "value"),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
	}
}
