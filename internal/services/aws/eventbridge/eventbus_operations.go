package eventbridge

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

func eventBusToMap(eb *eventsstore.EventBus) map[string]interface{} {
	result := map[string]interface{}{
		"Arn":              eb.ARN,
		"Name":             eb.Name,
		"CreationTime":     eb.CreatedAt.Unix(),
		"LastModifiedTime": eb.LastModifiedAt.Unix(),
	}
	if eb.Description != "" {
		result["Description"] = eb.Description
	}
	if eb.Policy != "" {
		result["Policy"] = eb.Policy
	}
	return result
}

// CreateEventBus creates a new event bus.
func (s *EventsService) CreateEventBus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Event bus name is required")
	}
	if name == "default" {
		return nil, awserrors.NewValidationException("Cannot create event bus named 'default'")
	}

	eventBus := &eventsstore.EventBus{
		Name: name,
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		eventBus.Description = desc
	}

	if policy, ok := req.Parameters["Policy"].(string); ok {
		eventBus.Policy = policy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateEventBus(ctx, eventBus); err != nil {
		return nil, mapStoreError(err, name)
	}

	if tags := tagutil.ParseTags(req.Parameters, "Tags"); len(tags) > 0 {
		if err := store.TagStore.TagFromSlice(eventBus.ARN, tags); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"EventBusArn": eventBus.ARN,
	}, nil
}

// DeleteEventBus deletes an event bus.
func (s *EventsService) DeleteEventBus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Event bus name is required")
	}
	if name == "default" {
		return nil, awserrors.NewValidationException("Cannot delete event bus 'default'")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	eventBus, err := store.GetEventBus(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}
	// Cascade-delete: rules → targets, then archives
	rulesResult, err := store.ListRules(ctx, name, "", 1000, "")
	if err == nil {
		for _, rule := range rulesResult.Rules {
			targetsResult, tErr := store.ListTargetsByRule(ctx, name, rule.Name, 1000, "")
			if tErr == nil {
				for _, t := range targetsResult.Targets {
					_ = store.DeleteTarget(ctx, name, rule.Name, t.ID)
				}
			}
			_ = store.DeleteRule(ctx, name, rule.Name)
			store.TagStore.Delete(rule.ARN)
		}
	}

	archives, err := store.ListArchivesForEventBus(ctx, name)
	if err == nil {
		for _, a := range archives {
			_ = store.DeleteArchiveEvents(ctx, a.Name)
			_ = store.DeleteArchive(ctx, a.Name)
		}
	}

	if err := store.DeleteEventBus(ctx, name); err != nil {
		return nil, err
	}

	store.TagStore.Delete(eventBus.ARN)

	return response.EmptyResponse(), nil
}

// DescribeEventBus returns information about an event bus.
func (s *EventsService) DescribeEventBus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		name = "default"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	eventBus, err := store.GetEventBus(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	result := eventBusToMap(eventBus)

	if tagSlice, err := store.TagStore.ListAsSlice(eventBus.ARN); err == nil && len(tagSlice) > 0 {
		tagMaps := make([]map[string]string, 0, len(tagSlice))
		for _, t := range tagSlice {
			tagMaps = append(tagMaps, map[string]string{"Key": t.Key, "Value": t.Value})
		}
		result["Tags"] = tagMaps
	}

	return result, nil
}

// ListEventBuses returns a list of event buses.
func (s *EventsService) ListEventBuses(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListEventBuses(ctx, namePrefix, limit, nextToken)
	if err != nil {
		return nil, err
	}

	eventBuses := make([]map[string]interface{}, len(result.EventBuses))
	for i, eb := range result.EventBuses {
		eventBuses[i] = eventBusToMap(eb)
	}

	response := map[string]interface{}{
		"EventBuses": eventBuses,
	}

	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}

// UpdateEventBus updates an event bus.
func (s *EventsService) UpdateEventBus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Event bus name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	eventBus, err := store.GetEventBus(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		eventBus.Description = desc
	}

	if policy, ok := req.Parameters["Policy"].(string); ok {
		eventBus.Policy = policy
	}

	if err := store.UpdateEventBus(ctx, eventBus); err != nil {
		return nil, err
	}

	return eventBusToMap(eventBus), nil
}
