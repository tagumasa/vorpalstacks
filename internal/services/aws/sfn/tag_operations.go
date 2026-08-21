package sfn

import (
	"context"
	"errors"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

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
			// Resource existence is verified by ValidateResource above;
			// this closure applies quota enforcement and persistence via
			// the shared tag Core path.
			if err := enforceTagQuota(store, resourceKey, tagSlice); err != nil {
				return err
			}
			return store.TagFromSlice(resourceKey, tagSlice)
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return store.ListAsSlice(resourceKey)
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

// validateTaggableResource probes the resource identified by an ARN for a
// tag operation. SFN supports tagging on state machines, activities, state
// machine aliases, state machine versions and map runs; each type is
// resolved through its store getter and mapped to the documented
// not-found error. The resource type is read from the ARN resource field
// (stateMachine:<name>, stateMachineAlias:<name>, mapRun:<id>,
// activity:<id>). Any other ARN shape — executions, or ARNs that are not
// States resources at all — is rejected outright.
func validateTaggableResource(ctx context.Context, store *sfnstore.StepFunctionStore, arn string) error {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	switch {
	case strings.HasPrefix(resource, "stateMachineAlias:"):
		if _, err := store.GetStateMachineAlias(ctx, arn); err != nil {
			if errors.Is(err, sfnstore.ErrStateMachineAliasNotFound) {
				return NewResourceNotFound("State Machine Alias Does not exist: " + arn)
			}
			return err
		}
	case strings.HasPrefix(resource, "mapRun:"):
		if _, err := store.GetMapRun(ctx, arn); err != nil {
			if errors.Is(err, sfnstore.ErrMapRunNotFound) {
				return NewResourceNotFound("Map Run does not exist: " + arn)
			}
			return err
		}
	case strings.HasPrefix(resource, "stateMachine:"):
		// Version ARNs append ":<number>" to the state machine ARN and
		// do not resolve through the plain state-machine getter; probe the
		// version store when the state machine lookup misses.
		if _, err := store.GetStateMachine(ctx, arn); err != nil {
			if !errors.Is(err, sfnstore.ErrStateMachineNotFound) {
				return err
			}
			if _, verr := store.GetStateMachineVersion(ctx, arn); verr != nil {
				if errors.Is(verr, sfnstore.ErrStateMachineVersionNotFound) {
					return NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
				}
				return verr
			}
		}
	case strings.HasPrefix(resource, "activity:"):
		if _, err := store.GetActivity(ctx, arn); err != nil {
			if errors.Is(err, sfnstore.ErrActivityNotFound) {
				return NewActivityDoesNotExist("Activity Does not exist: " + arn)
			}
			return err
		}
	default:
		// Executions and any non-States ARN have no tag store behind
		// them; accepting them would persist phantom tag records.
		return NewResourceNotFound("Resource does not exist: " + arn)
	}
	return nil
}

// TagResource adds tags to a state machine.
func (s *StepFunctionService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.tagHandlerConfig(store))
}

// UntagResource removes tags from a state machine.
func (s *StepFunctionService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.tagHandlerConfig(store))
}

// ListTagsForResource returns the tags for a state machine.
func (s *StepFunctionService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.tagHandlerConfig(store))
}
