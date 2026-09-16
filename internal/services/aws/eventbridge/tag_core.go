package eventbridge

import (
	"context"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// resolveTaggableResourceCore validates that the ARN refers to an existing
// taggable EventBridge resource. The taggable set is exactly event buses
// and rules — the tagging page enumerates nothing else ("In EventBridge,
// you can assign tags to rule and event buses", eb-tagging, fetched
// 2026-09-15) — so archive, connection and API-destination ARNs are not
// taggable and answer with the same resource-not-found shape as any other
// unsupported resource form. The ResourceARN requirement of the tag family
// is enforced here so every calling plane sees it.
func (s *EventsService) resolveTaggableResourceCore(ctx context.Context, store *eventsstore.EventsStore, resourceArn string) error {
	if resourceArn == "" {
		return awserrors.NewValidationException("ResourceARN is required")
	}
	_, _, _, _, resource := svcarn.SplitARN(resourceArn)

	switch {
	case strings.HasPrefix(resource, "event-bus/"):
		name := svcarn.ExtractEventBusNameFromARN(resourceArn)
		if _, err := store.GetEventBus(ctx, name); err != nil {
			return NewResourceNotFoundException("Event bus '" + name + "' does not exist")
		}
	case strings.HasPrefix(resource, "rule/"):
		eventBusName, ruleName := extractRuleInfoFromArn(resourceArn)
		if _, err := store.GetRule(ctx, eventBusName, ruleName); err != nil {
			return NewResourceNotFoundException("Rule '" + ruleName + "' does not exist")
		}
	default:
		return NewResourceNotFoundException("Resource not found: " + resourceArn)
	}

	return nil
}

// validateResourceTags enforces the documented tag limits: at most 50 tags
// per resource, keys of up to 128 Unicode characters, values of up to 256,
// and the aws: prefix reserved (eb-tagging, fetched 2026-09-15).
func validateResourceTags(tags []tagutil.Tag) error {
	switch v, key := tagutil.CheckTags(tags, tagutil.StandardLimits()); v {
	case tagutil.TooManyTags:
		return awserrors.NewValidationException(fmt.Sprintf("too many tags (maximum %d)", tagutil.MaxTagsPerResource))
	case tagutil.TagKeyTooShort:
		return awserrors.NewValidationException("tag key cannot be empty")
	case tagutil.TagKeyTooLong:
		return awserrors.NewValidationException(fmt.Sprintf("tag key %q cannot exceed %d characters", key, tagutil.MaxTagKeyLength))
	case tagutil.TagValueTooLong:
		return awserrors.NewValidationException(fmt.Sprintf("tag value for key %q cannot exceed %d characters", key, tagutil.MaxTagValueLength))
	case tagutil.ReservedTagKey:
		return awserrors.NewValidationException(fmt.Sprintf("tag key %q cannot start with 'aws:' (reserved prefix)", key))
	}
	return nil
}

// tagResourceCore validates the required members, the resource's
// taggability and the tag limits, then applies the tag map. Tags is a
// required member of TagResourceRequest; the presence rejection lives
// here so every calling plane sees it.
func (s *EventsService) tagResourceCore(ctx context.Context, store *eventsstore.EventsStore, resourceArn string, tags []tagutil.Tag) error {
	if len(tags) == 0 {
		return awserrors.NewValidationException("Tags are required")
	}
	if err := s.resolveTaggableResourceCore(ctx, store, resourceArn); err != nil {
		return err
	}

	if err := validateResourceTags(tags); err != nil {
		return err
	}

	tagMap := make(map[string]string, len(tags))
	for _, t := range tags {
		tagMap[t.Key] = t.Value
	}
	return store.TagStore.Tag(resourceArn, tagMap)
}

// untagResourceCore validates the required members and the resource's
// taggability, then removes the tag keys. TagKeys is a required member of
// UntagResourceRequest; the presence rejection lives here so every
// calling plane sees it.
func (s *EventsService) untagResourceCore(ctx context.Context, store *eventsstore.EventsStore, resourceArn string, tagKeysMap map[string]bool) error {
	if len(tagKeysMap) == 0 {
		return awserrors.NewValidationException("TagKeys are required")
	}
	if err := s.resolveTaggableResourceCore(ctx, store, resourceArn); err != nil {
		return err
	}

	keys := make([]string, 0, len(tagKeysMap))
	for k := range tagKeysMap {
		keys = append(keys, k)
	}
	return store.TagStore.Untag(resourceArn, keys)
}

// listTagsForResourceCore validates the resource exists and returns its tags.
func (s *EventsService) listTagsForResourceCore(ctx context.Context, store *eventsstore.EventsStore, resourceArn string) ([]tagutil.Tag, error) {
	if err := s.resolveTaggableResourceCore(ctx, store, resourceArn); err != nil {
		return nil, err
	}
	return store.TagStore.ListAsSlice(resourceArn)
}

// tagListToMaps converts a tag slice into the wire tag list shape shared by
// the tag-family responses.
func tagListToMaps(tagSlice []tagutil.Tag) []map[string]string {
	tagMaps := make([]map[string]string, 0, len(tagSlice))
	for _, t := range tagSlice {
		tagMaps = append(tagMaps, map[string]string{
			"Key":   t.Key,
			"Value": t.Value,
		})
	}
	return tagMaps
}
