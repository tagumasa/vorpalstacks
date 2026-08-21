package sfn

import (
	"context"
	"fmt"

	types "vorpalstacks/internal/common/tags"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
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
