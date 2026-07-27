package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// eventBusToListItem serialises an EventBus for the ListEventBuses response.
// Per Smithy EventBus shape: Name, Arn, Description, Policy, CreationTime,
// LastModifiedTime.
func eventBusToListItem(eb *eventsstore.EventBus) map[string]interface{} {
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

// eventBusToDescribeMap serialises an EventBus for the
// DescribeEventBusResponse shape. Includes the Describe-only fields:
// KmsKeyIdentifier, DeadLetterConfig, LogConfig.
func eventBusToDescribeMap(eb *eventsstore.EventBus) map[string]interface{} {
	result := eventBusToListItem(eb)
	if eb.KmsKeyIdentifier != "" {
		result["KmsKeyIdentifier"] = eb.KmsKeyIdentifier
	}
	if eb.DeadLetterConfig != nil {
		dlc := map[string]interface{}{}
		if eb.DeadLetterConfig.Arn != "" {
			dlc["Arn"] = eb.DeadLetterConfig.Arn
		}
		result["DeadLetterConfig"] = dlc
	}
	if eb.LogConfig != nil {
		lc := map[string]interface{}{}
		if eb.LogConfig.IncludeDetail != "" {
			lc["IncludeDetail"] = eb.LogConfig.IncludeDetail
		}
		if eb.LogConfig.Level != "" {
			lc["Level"] = eb.LogConfig.Level
		}
		result["LogConfig"] = lc
	}
	return result
}

// eventBusToUpdateMap serialises an EventBus for the UpdateEventBusResponse
// shape per Smithy: Arn, Name, KmsKeyIdentifier, Description,
// DeadLetterConfig, LogConfig.
func eventBusToUpdateMap(eb *eventsstore.EventBus) map[string]interface{} {
	result := map[string]interface{}{
		"Arn":  eb.ARN,
		"Name": eb.Name,
	}
	if eb.Description != "" {
		result["Description"] = eb.Description
	}
	if eb.KmsKeyIdentifier != "" {
		result["KmsKeyIdentifier"] = eb.KmsKeyIdentifier
	}
	if eb.DeadLetterConfig != nil {
		dlc := map[string]interface{}{}
		if eb.DeadLetterConfig.Arn != "" {
			dlc["Arn"] = eb.DeadLetterConfig.Arn
		}
		result["DeadLetterConfig"] = dlc
	}
	if eb.LogConfig != nil {
		lc := map[string]interface{}{}
		if eb.LogConfig.IncludeDetail != "" {
			lc["IncludeDetail"] = eb.LogConfig.IncludeDetail
		}
		if eb.LogConfig.Level != "" {
			lc["Level"] = eb.LogConfig.Level
		}
		result["LogConfig"] = lc
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

	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok {
		eventBus.KmsKeyIdentifier = kms
	}
	if dlc, ok := req.Parameters["DeadLetterConfig"].(map[string]interface{}); ok {
		eventBus.DeadLetterConfig = &eventsstore.DeadLetterConfig{}
		if arn, ok := dlc["Arn"].(string); ok {
			eventBus.DeadLetterConfig.Arn = arn
		}
	}
	if lc, ok := req.Parameters["LogConfig"].(map[string]interface{}); ok {
		eventBus.LogConfig = &eventsstore.BusLogConfig{}
		if id, ok := lc["IncludeDetail"].(string); ok {
			if !isValidLogIncludeDetail(id) {
				return nil, awserrors.NewValidationException("LogConfig.IncludeDetail must be one of: NONE, FULL")
			}
			eventBus.LogConfig.IncludeDetail = id
		}
		if lvl, ok := lc["Level"].(string); ok {
			if !isValidLogLevel(lvl) {
				return nil, awserrors.NewValidationException("LogConfig.Level must be one of: OFF, ERROR, INFO, TRACE")
			}
			eventBus.LogConfig.Level = lvl
		}
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

	// CreateEventBusResponse shape: EventBusArn, Description,
	// KmsKeyIdentifier, DeadLetterConfig, LogConfig.
	resp := map[string]interface{}{
		"EventBusArn": eventBus.ARN,
	}
	if eventBus.Description != "" {
		resp["Description"] = eventBus.Description
	}
	if eventBus.KmsKeyIdentifier != "" {
		resp["KmsKeyIdentifier"] = eventBus.KmsKeyIdentifier
	}
	if eventBus.DeadLetterConfig != nil {
		dlc := map[string]interface{}{}
		if eventBus.DeadLetterConfig.Arn != "" {
			dlc["Arn"] = eventBus.DeadLetterConfig.Arn
		}
		resp["DeadLetterConfig"] = dlc
	}
	if eventBus.LogConfig != nil {
		lc := map[string]interface{}{}
		if eventBus.LogConfig.IncludeDetail != "" {
			lc["IncludeDetail"] = eventBus.LogConfig.IncludeDetail
		}
		if eventBus.LogConfig.Level != "" {
			lc["Level"] = eventBus.LogConfig.Level
		}
		resp["LogConfig"] = lc
	}
	return resp, nil
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
	// Cascade-delete: rules → targets (paginated), then archives.
	// AWS aborts the entire DeleteEventBus when any cascade step fails
	// (ConcurrentModificationException or InternalException), so that the
	// bus remains queryable for follow-up diagnosis. We follow the same
	// contract: collect the first cascade error and abort without deleting
	// the bus when one occurs.
	var cascadeErr error

	ruleToken := ""
	for cascadeErr == nil {
		rulesResult, err := store.ListRules(ctx, name, "", 1000, ruleToken)
		if err != nil {
			cascadeErr = fmt.Errorf("DeleteEventBus: list rules: %w", err)
			break
		}
		for _, rule := range rulesResult.Rules {
			targetToken := ""
			for cascadeErr == nil {
				targetsResult, tErr := store.ListTargetsByRule(ctx, name, rule.Name, 1000, targetToken)
				if tErr != nil {
					cascadeErr = fmt.Errorf("DeleteEventBus: list targets for rule %s: %w", rule.Name, tErr)
					break
				}
				for _, t := range targetsResult.Targets {
					if err := store.DeleteTarget(ctx, name, rule.Name, t.ID); err != nil {
						cascadeErr = fmt.Errorf("DeleteEventBus: delete target %s: %w", t.ID, err)
						break
					}
				}
				if cascadeErr != nil {
					break
				}
				if targetsResult.NextToken == "" {
					break
				}
				targetToken = targetsResult.NextToken
			}
			if cascadeErr != nil {
				break
			}
			if err := store.DeleteRule(ctx, name, rule.Name); err != nil {
				cascadeErr = fmt.Errorf("DeleteEventBus: delete rule %s: %w", rule.Name, err)
				break
			}
			lastFireTimes.Delete(rule.ARN)
			_ = store.TagStore.Delete(rule.ARN)
		}
		if cascadeErr != nil {
			break
		}
		if rulesResult.NextToken == "" {
			break
		}
		ruleToken = rulesResult.NextToken
	}

	if cascadeErr == nil {
		archives, err := store.ListArchivesForEventBus(ctx, name)
		if err != nil {
			cascadeErr = fmt.Errorf("DeleteEventBus: list archives: %w", err)
		} else {
			for _, a := range archives {
				if err := store.DeleteArchiveEvents(ctx, a.Name); err != nil {
					cascadeErr = fmt.Errorf("DeleteEventBus: delete archive events %s: %w", a.Name, err)
					break
				}
				if err := store.DeleteArchive(ctx, a.Name); err != nil {
					cascadeErr = fmt.Errorf("DeleteEventBus: delete archive %s: %w", a.Name, err)
					break
				}
			}
		}
	}

	if cascadeErr != nil {
		// Leave the bus in place so callers can inspect orphaned resources.
		return nil, awserrors.NewAWSError(
			"InternalException",
			cascadeErr.Error(),
			http.StatusInternalServerError,
		)
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

	result := eventBusToDescribeMap(eventBus)

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
		eventBuses[i] = eventBusToListItem(eb)
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

	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok {
		eventBus.KmsKeyIdentifier = kms
	}
	if dlc, ok := req.Parameters["DeadLetterConfig"].(map[string]interface{}); ok {
		eventBus.DeadLetterConfig = &eventsstore.DeadLetterConfig{}
		if arn, ok := dlc["Arn"].(string); ok {
			eventBus.DeadLetterConfig.Arn = arn
		}
	}
	if lc, ok := req.Parameters["LogConfig"].(map[string]interface{}); ok {
		eventBus.LogConfig = &eventsstore.BusLogConfig{}
		if id, ok := lc["IncludeDetail"].(string); ok {
			if !isValidLogIncludeDetail(id) {
				return nil, awserrors.NewValidationException("LogConfig.IncludeDetail must be one of: NONE, FULL")
			}
			eventBus.LogConfig.IncludeDetail = id
		}
		if lvl, ok := lc["Level"].(string); ok {
			if !isValidLogLevel(lvl) {
				return nil, awserrors.NewValidationException("LogConfig.Level must be one of: OFF, ERROR, INFO, TRACE")
			}
			eventBus.LogConfig.Level = lvl
		}
	}

	if err := store.UpdateEventBus(ctx, eventBus); err != nil {
		return nil, err
	}

	return eventBusToUpdateMap(eventBus), nil
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
		// AWS enforces an 8192-byte ceiling on the resource policy.
		if len(policyStr) > 8192 {
			return nil, awserrors.NewPolicyLengthExceededException(
				fmt.Sprintf("Event bus policy length %d exceeds the maximum allowed length of 8192 bytes", len(policyStr)),
			)
		}
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
	removeAll, _ := req.Parameters["RemoveAllPermissions"].(bool)

	// AWS requires either StatementId or RemoveAllPermissions=true.
	if !removeAll && statementID == "" {
		return nil, awserrors.NewValidationException("StatementId is required when RemoveAllPermissions is not true")
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

	// RemoveAllPermissions clears the policy entirely.
	if removeAll {
		eventBus.Policy = ""
		if err := store.UpdateEventBus(ctx, eventBus); err != nil {
			return nil, err
		}
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

// isValidLogIncludeDetail validates IncludeDetail against the Smithy enum
// values: NONE or FULL.
func isValidLogIncludeDetail(v string) bool {
	return v == "NONE" || v == "FULL"
}

// isValidLogLevel validates Level against the Smithy enum values:
// OFF, ERROR, INFO, or TRACE.
func isValidLogLevel(v string) bool {
	switch v {
	case "OFF", "ERROR", "INFO", "TRACE":
		return true
	}
	return false
}
