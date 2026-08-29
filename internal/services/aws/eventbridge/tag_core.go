package eventbridge

import (
	"context"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// resolveTaggableResourceCore validates that the ARN refers to an existing
// EventBridge resource (event bus, rule, archive, connection, or API
// destination).
func (s *EventsService) resolveTaggableResourceCore(ctx context.Context, store *eventsstore.EventsStore, resourceArn string) error {
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
	case strings.HasPrefix(resource, "archive/"):
		name := strings.TrimPrefix(resource, "archive/")
		if _, err := store.GetArchive(ctx, name); err != nil {
			return NewResourceNotFoundException("Archive '" + name + "' does not exist")
		}
	case strings.HasPrefix(resource, "connection/"):
		name := strings.TrimPrefix(resource, "connection/")
		if _, err := store.GetConnection(ctx, name); err != nil {
			return NewResourceNotFoundException("Connection '" + name + "' does not exist")
		}
	case strings.HasPrefix(resource, "api-destination/"):
		name := strings.TrimPrefix(resource, "api-destination/")
		if _, err := store.GetApiDestination(ctx, name); err != nil {
			return NewResourceNotFoundException("API destination '" + name + "' does not exist")
		}
	default:
		return NewResourceNotFoundException("Resource not found: " + resourceArn)
	}

	return nil
}

// tagResourceCore validates the resource exists and applies the tag map.
func (s *EventsService) tagResourceCore(ctx context.Context, store *eventsstore.EventsStore, resourceArn string, tags []tagutil.Tag) error {
	if err := s.resolveTaggableResourceCore(ctx, store, resourceArn); err != nil {
		return err
	}

	tagMap := make(map[string]string, len(tags))
	for _, t := range tags {
		tagMap[t.Key] = t.Value
	}
	return store.TagStore.Tag(resourceArn, tagMap)
}

// untagResourceCore validates the resource exists and removes the tag keys.
func (s *EventsService) untagResourceCore(ctx context.Context, store *eventsstore.EventsStore, resourceArn string, tagKeysMap map[string]bool) error {
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

// validateResourceArnParam enforces the common ResourceARN requirement of
// the tag family.
func validateResourceArnParam(resourceArn string) error {
	if resourceArn == "" {
		return awserrors.NewValidationException("ResourceARN is required")
	}
	return nil
}
