package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"

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
	if _, err := store.GetEventBus(ctx, name); err != nil {
		return nil, mapStoreError(err, name)
	}
	// Cascade-delete: rules → targets (paginated), then archives
	ruleToken := ""
	for {
		rulesResult, err := store.ListRules(ctx, name, "", 1000, ruleToken)
		if err != nil {
			break
		}
		for _, rule := range rulesResult.Rules {
			targetToken := ""
			for {
				targetsResult, tErr := store.ListTargetsByRule(ctx, name, rule.Name, 1000, targetToken)
				if tErr != nil {
					break
				}
				for _, t := range targetsResult.Targets {
					_ = store.DeleteTarget(ctx, name, rule.Name, t.ID)
				}
				if targetsResult.NextToken == "" {
					break
				}
				targetToken = targetsResult.NextToken
			}
			_ = store.DeleteRule(ctx, name, rule.Name)
			lastFireTimes.Delete(rule.ARN)
			_ = store.TagStore.Delete(rule.ARN)
		}
		if rulesResult.NextToken == "" {
			break
		}
		ruleToken = rulesResult.NextToken
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

// PutPermission adds a resource policy statement to the specified event bus,
// granting the given principal permission to put events. Supports two modes:
// (1) individual parameters (Principal, StatementId, Action, Condition) and
// (2) a complete policy document via the Policy parameter.
func (s *EventsService) PutPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	busName := request.GetParamLowerFirst(req.Parameters, "EventBusName")
	if busName == "" {
		busName = "default"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	eventBus, err := store.GetEventBus(ctx, busName)
	if err != nil {
		return nil, mapStoreError(err, busName)
	}

	// Mode 1: Full policy document provided via the Policy parameter.
	if policyStr, ok := req.Parameters["Policy"].(string); ok && policyStr != "" {
		var policyDoc map[string]interface{}
		if err := json.Unmarshal([]byte(policyStr), &policyDoc); err != nil {
			return nil, awserrors.NewValidationException("Invalid policy document")
		}
		if _, ok := policyDoc["Version"]; !ok {
			policyDoc["Version"] = "2012-10-17"
		}
		policyBytes, err := json.Marshal(policyDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal policy: %w", err)
		}
		eventBus.Policy = string(policyBytes)
		if err := store.UpdateEventBus(ctx, eventBus); err != nil {
			return nil, err
		}
		return response.EmptyResponse(), nil
	}

	// Mode 2: Individual parameters (Principal, StatementId, Action, Condition).
	principal := request.GetStringParam(req.Parameters, "Principal")
	statementID := request.GetStringParam(req.Parameters, "StatementId")
	action := request.GetStringParam(req.Parameters, "Action")
	if action == "" {
		action = "events:PutEvents"
	}

	if principal == "" || statementID == "" {
		return nil, awserrors.NewValidationException("Principal and StatementId are required")
	}

	var policyDoc map[string]interface{}
	if eventBus.Policy != "" {
		if err := json.Unmarshal([]byte(eventBus.Policy), &policyDoc); err != nil {
			policyDoc = make(map[string]interface{})
		}
	}
	if _, ok := policyDoc["Version"]; !ok {
		policyDoc["Version"] = "2012-10-17"
	}

	statement := map[string]interface{}{
		"Sid":       statementID,
		"Effect":    "Allow",
		"Principal": map[string]interface{}{"AWS": principal},
		"Action":    action,
		"Resource":  eventBus.ARN,
	}
	if condition, ok := req.Parameters["Condition"].(string); ok && condition != "" {
		var cond map[string]interface{}
		if err := json.Unmarshal([]byte(condition), &cond); err == nil {
			statement["Condition"] = cond
		}
	}

	statements, _ := policyDoc["Statement"].([]interface{})
	replaced := false
	for i, s := range statements {
		if stmt, ok := s.(map[string]interface{}); ok {
			if sid, _ := stmt["Sid"].(string); sid == statementID {
				statements[i] = statement
				replaced = true
				break
			}
		}
	}
	if !replaced {
		statements = append(statements, statement)
	}
	policyDoc["Statement"] = statements

	policyBytes, err := json.Marshal(policyDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal policy: %w", err)
	}
	eventBus.Policy = string(policyBytes)

	if err := store.UpdateEventBus(ctx, eventBus); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemovePermission removes a resource policy statement from the specified
// event bus identified by its StatementId.
func (s *EventsService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	busName := request.GetParamLowerFirst(req.Parameters, "EventBusName")
	if busName == "" {
		busName = "default"
	}

	statementID := request.GetStringParam(req.Parameters, "StatementId")
	if statementID == "" {
		return nil, awserrors.NewValidationException("StatementId is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	eventBus, err := store.GetEventBus(ctx, busName)
	if err != nil {
		return nil, mapStoreError(err, busName)
	}

	if eventBus.Policy == "" {
		return response.EmptyResponse(), nil
	}

	var policyDoc map[string]interface{}
	if err := json.Unmarshal([]byte(eventBus.Policy), &policyDoc); err != nil {
		return response.EmptyResponse(), nil
	}

	statements, ok := policyDoc["Statement"].([]interface{})
	if !ok {
		return response.EmptyResponse(), nil
	}

	filtered := make([]interface{}, 0, len(statements))
	for _, s := range statements {
		if stmt, ok := s.(map[string]interface{}); ok {
			if sid, _ := stmt["Sid"].(string); sid != statementID {
				filtered = append(filtered, s)
			}
		}
	}

	if len(filtered) == 0 {
		delete(policyDoc, "Statement")
	} else {
		policyDoc["Statement"] = filtered
	}

	policyBytes, err := json.Marshal(policyDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal policy: %w", err)
	}
	eventBus.Policy = string(policyBytes)

	if err := store.UpdateEventBus(ctx, eventBus); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
