package eventbridge

import (
	"context"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// taggableResource represents a validated EventBridge resource that can be tagged.
type taggableResource struct {
	arn string
}

// resolveTaggableResource validates that the ARN refers to an existing EventBridge
// resource (event bus, rule, archive, connection, or API destination).
func (s *EventsService) resolveTaggableResource(ctx context.Context, reqCtx *request.RequestContext, resourceArn string) (*taggableResource, error) {
	_, _, _, _, resource := svcarn.SplitARN(resourceArn)
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasPrefix(resource, "event-bus/"):
		name := svcarn.ExtractEventBusNameFromARN(resourceArn)
		if _, err := store.GetEventBus(ctx, name); err != nil {
			return nil, NewResourceNotFoundException("Event bus '" + name + "' does not exist")
		}
	case strings.HasPrefix(resource, "rule/"):
		eventBusName, ruleName := extractRuleInfoFromArn(resourceArn)
		if _, err := store.GetRule(ctx, eventBusName, ruleName); err != nil {
			return nil, NewResourceNotFoundException("Rule '" + ruleName + "' does not exist")
		}
	case strings.HasPrefix(resource, "archive/"):
		name := strings.TrimPrefix(resource, "archive/")
		if _, err := store.GetArchive(ctx, name); err != nil {
			return nil, NewResourceNotFoundException("Archive '" + name + "' does not exist")
		}
	case strings.HasPrefix(resource, "connection/"):
		name := strings.TrimPrefix(resource, "connection/")
		if _, err := store.GetConnection(ctx, name); err != nil {
			return nil, NewResourceNotFoundException("Connection '" + name + "' does not exist")
		}
	case strings.HasPrefix(resource, "api-destination/"):
		name := strings.TrimPrefix(resource, "api-destination/")
		if _, err := store.GetApiDestination(ctx, name); err != nil {
			return nil, NewResourceNotFoundException("API destination '" + name + "' does not exist")
		}
	default:
		return nil, NewResourceNotFoundException("Resource not found: " + resourceArn)
	}

	return &taggableResource{arn: resourceArn}, nil
}

// TagResource adds tags to an EventBridge resource.
func (s *EventsService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, awserrors.NewValidationException("ResourceARN is required")
	}

	newTags := tagutil.ParseTags(req.Parameters, "Tags")
	if len(newTags) == 0 {
		return nil, awserrors.NewValidationException("Tags are required")
	}

	resolved, err := s.resolveTaggableResource(ctx, reqCtx, resourceArn)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tagMap := make(map[string]string, len(newTags))
	for _, t := range newTags {
		tagMap[t.Key] = t.Value
	}
	if err := store.TagStore.Tag(resolved.arn, tagMap); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from an EventBridge resource.
func (s *EventsService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, awserrors.NewValidationException("ResourceARN is required")
	}

	tagKeysMap := tagutil.ParseTagKeys(req.Parameters, "TagKeys")
	if len(tagKeysMap) == 0 {
		tagKeysMap = tagutil.ParseTagKeys(req.Parameters, "tagKeys")
	}
	if len(tagKeysMap) == 0 {
		return nil, awserrors.NewValidationException("TagKeys are required")
	}

	resolved, err := s.resolveTaggableResource(ctx, reqCtx, resourceArn)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(tagKeysMap))
	for k := range tagKeysMap {
		keys = append(keys, k)
	}
	if err := store.TagStore.Untag(resolved.arn, keys); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for an EventBridge resource.
func (s *EventsService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, awserrors.NewValidationException("ResourceARN is required")
	}

	resolved, err := s.resolveTaggableResource(ctx, reqCtx, resourceArn)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tagSlice, err := store.TagStore.ListAsSlice(resolved.arn)
	if err != nil {
		return nil, err
	}
	tagMaps := make([]map[string]string, 0, len(tagSlice))
	for _, t := range tagSlice {
		tagMaps = append(tagMaps, map[string]string{
			"Key":   t.Key,
			"Value": t.Value,
		})
	}
	return map[string]interface{}{"Tags": tagMaps}, nil
}

func extractRuleInfoFromArn(arn string) (eventBusName, ruleName string) {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	if resource == "" {
		return "", ""
	}
	parts := strings.Split(resource, "/")
	if len(parts) >= 3 && parts[0] == "rule" {
		eventBusName = parts[1]
		ruleName = parts[2]
	} else if len(parts) >= 2 {
		eventBusName = "default"
		ruleName = parts[len(parts)-1]
	}
	return eventBusName, ruleName
}
