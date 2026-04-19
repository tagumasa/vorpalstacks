package eventbridge

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TagResource adds tags to an EventBridge resource.
func (s *EventsService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, NewValidationException("ResourceARN is required")
	}

	newTags := tagutil.ParseTags(req.Parameters, "Tags")
	if len(newTags) == 0 {
		return nil, NewValidationException("Tags are required")
	}

	_, _, _, _, resource := svcarn.SplitARN(resourceArn)
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(resource, "rule/") {
		if _, err := s.getRuleByArn(ctx, reqCtx, resourceArn); err == nil {
			tagMap := make(map[string]string, len(newTags))
			for _, t := range newTags {
				tagMap[t.Key] = t.Value
			}
			if err := store.TagStore.Tag(resourceArn, tagMap); err != nil {
				return nil, err
			}
			return response.EmptyResponse(), nil
		}
	}

	if _, err := s.getEventBusByArn(ctx, reqCtx, resourceArn); err == nil {
		tagMap := make(map[string]string, len(newTags))
		for _, t := range newTags {
			tagMap[t.Key] = t.Value
		}
		if err := store.TagStore.Tag(resourceArn, tagMap); err != nil {
			return nil, err
		}
		return response.EmptyResponse(), nil
	}

	return nil, NewResourceNotFoundException("Resource not found")
}

// UntagResource removes tags from an EventBridge resource.
func (s *EventsService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, NewValidationException("ResourceARN is required")
	}

	tagKeysMap := tagutil.ParseTagKeys(req.Parameters, "TagKeys")
	if len(tagKeysMap) == 0 {
		tagKeysMap = tagutil.ParseTagKeys(req.Parameters, "tagKeys")
	}
	if len(tagKeysMap) == 0 {
		return nil, NewValidationException("TagKeys are required")
	}

	_, _, _, _, resource := svcarn.SplitARN(resourceArn)
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(resource, "rule/") {
		if _, err := s.getRuleByArn(ctx, reqCtx, resourceArn); err == nil {
			keys := make([]string, 0, len(tagKeysMap))
			for k := range tagKeysMap {
				keys = append(keys, k)
			}
			if err := store.TagStore.Untag(resourceArn, keys); err != nil {
				return nil, err
			}
			return response.EmptyResponse(), nil
		}
	}

	if _, err := s.getEventBusByArn(ctx, reqCtx, resourceArn); err == nil {
		keys := make([]string, 0, len(tagKeysMap))
		for k := range tagKeysMap {
			keys = append(keys, k)
		}
		if err := store.TagStore.Untag(resourceArn, keys); err != nil {
			return nil, err
		}
		return response.EmptyResponse(), nil
	}

	return nil, NewResourceNotFoundException("Resource not found")
}

// ListTagsForResource lists tags for an EventBridge resource.
func (s *EventsService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, NewValidationException("ResourceARN is required")
	}

	_, _, _, _, resource := svcarn.SplitARN(resourceArn)
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	isRule := strings.HasPrefix(resource, "rule/")
	if isRule {
		if _, err := s.getRuleByArn(ctx, reqCtx, resourceArn); err != nil {
			return nil, NewResourceNotFoundException("Resource not found: " + resourceArn)
		}
	} else {
		if _, err := s.getEventBusByArn(ctx, reqCtx, resourceArn); err != nil {
			return nil, NewResourceNotFoundException("Resource not found: " + resourceArn)
		}
	}

	tagSlice, err := store.TagStore.ListAsSlice(resourceArn)
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

func (s *EventsService) getEventBusByArn(ctx context.Context, reqCtx *request.RequestContext, arnStr string) (*eventsstore.EventBus, error) {
	name := svcarn.ExtractEventBusNameFromARN(arnStr)
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return store.GetEventBus(ctx, name)
}

func (s *EventsService) getRuleByArn(ctx context.Context, reqCtx *request.RequestContext, arn string) (*eventsstore.Rule, error) {
	eventBusName, ruleName := extractRuleInfoFromArn(arn)
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return store.GetRule(ctx, eventBusName, ruleName)
}

func extractRuleInfoFromArn(arn string) (eventBusName, ruleName string) {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	if resource == "" {
		return "", ""
	}
	// resource format: rule/event-bus-name/rule-name
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
