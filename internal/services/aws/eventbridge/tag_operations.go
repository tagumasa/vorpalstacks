package eventbridge

import (
	"context"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
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
		if rule, err := s.getRuleByArn(ctx, reqCtx, resourceArn); err == nil {
			rule.Tags = tagutil.Apply(rule.Tags, newTags)
			if err := store.UpdateRule(ctx, rule); err != nil {
				return nil, err
			}
			return response.EmptyResponse(), nil
		}
	}

	if eventBus, err := s.getEventBusByArn(ctx, reqCtx, resourceArn); err == nil {
		eventBus.Tags = tagutil.Apply(eventBus.Tags, newTags)
		if err := store.UpdateEventBus(ctx, eventBus); err != nil {
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
		if rule, err := s.getRuleByArn(ctx, reqCtx, resourceArn); err == nil {
			rule.Tags = tagutil.Remove(rule.Tags, tagKeysMap)
			if err := store.UpdateRule(ctx, rule); err != nil {
				return nil, err
			}
			return response.EmptyResponse(), nil
		}
	}

	if eventBus, err := s.getEventBusByArn(ctx, reqCtx, resourceArn); err == nil {
		eventBus.Tags = tagutil.Remove(eventBus.Tags, tagKeysMap)
		if err := store.UpdateEventBus(ctx, eventBus); err != nil {
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

	var tags []map[string]string

	_, _, _, _, resource := svcarn.SplitARN(resourceArn)

	if strings.HasPrefix(resource, "rule/") {
		rule, err := s.getRuleByArn(ctx, reqCtx, resourceArn)
		if err == nil {
			for _, t := range rule.Tags {
				tags = append(tags, map[string]string{
					"Key":   t.Key,
					"Value": t.Value,
				})
			}
			return map[string]interface{}{"Tags": tags}, nil
		}
	}

	if eventBus, err := s.getEventBusByArn(ctx, reqCtx, resourceArn); err == nil {
		for _, t := range eventBus.Tags {
			tags = append(tags, map[string]string{
				"Key":   t.Key,
				"Value": t.Value,
			})
		}
		return map[string]interface{}{"Tags": tags}, nil
	}

	return nil, NewResourceNotFoundException("Resource not found: " + resourceArn)
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
