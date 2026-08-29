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

	input := CreateEventBusInput{
		Name: name,
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		if !validateDescription(desc) {
			return nil, errDescriptionTooLong()
		}
		input.Description = desc
	}

	if policy, ok := req.Parameters["Policy"].(string); ok {
		input.Policy = policy
	}

	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok {
		input.KmsKeyIdentifier = kms
	}
	if dlc, ok := req.Parameters["DeadLetterConfig"].(map[string]interface{}); ok {
		input.DeadLetterConfig = &eventsstore.DeadLetterConfig{}
		if arn, ok := dlc["Arn"].(string); ok {
			input.DeadLetterConfig.Arn = arn
		}
	}
	if lc, ok := req.Parameters["LogConfig"].(map[string]interface{}); ok {
		input.LogConfig = &BusLogConfigInput{}
		if id, ok := lc["IncludeDetail"].(string); ok {
			input.LogConfig.IncludeDetailSet = true
			input.LogConfig.IncludeDetail = id
		}
		if lvl, ok := lc["Level"].(string); ok {
			input.LogConfig.LevelSet = true
			input.LogConfig.Level = lvl
		}
	}

	if tags := tagutil.ParseTags(req.Parameters, "Tags"); len(tags) > 0 {
		input.Tags = tags
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createEventBusCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	eventBus := result.EventBus

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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteEventBusCore(ctx, store, DeleteEventBusInput{Name: name}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// eventBusNameParam extracts the raw EventBusName value with presence
// semantics across the case variants, distinguishing an omitted member
// from one provided as an empty (or non-string) value. The defaulting and
// the empty-value rejection run at the Core layer through
// resolveEventBusNameCore.
func eventBusNameParam(req *request.ParsedRequest) (string, bool) {
	for _, k := range []string{"EventBusName", "eventBusName"} {
		if _, ok := req.Parameters[k]; ok {
			str, _ := req.Parameters[k].(string)
			return str, true
		}
	}
	return "", false
}

// DescribeEventBus returns information about an event bus.
func (s *EventsService) DescribeEventBus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		if _, ok := req.Parameters["Name"]; ok {
			return nil, awserrors.NewValidationException("Name must not be empty")
		}
		name = "default"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeEventBusCore(ctx, store, name)
	if err != nil {
		return nil, err
	}

	response := eventBusToDescribeMap(result.EventBus)
	if len(result.Tags) > 0 {
		response["Tags"] = tagListToMaps(result.Tags)
	}

	return response, nil
}

// ListEventBuses returns a list of event buses.
func (s *EventsService) ListEventBuses(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listEventBusesCore(ctx, store, ListEventBusesInput{
		NamePrefix: namePrefix,
		Limit:      limit,
		NextToken:  nextToken,
	})
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

// parseUpdateEventBusInput reads the UpdateEventBus wire request into the
// transport-agnostic Core input.
func parseUpdateEventBusInput(req *request.ParsedRequest) UpdateEventBusInput {
	input := UpdateEventBusInput{
		Name: request.GetParamLowerFirst(req.Parameters, "Name"),
	}
	if desc, ok := req.Parameters["Description"].(string); ok {
		input.DescriptionSet = true
		input.Description = desc
	}
	if policy, ok := req.Parameters["Policy"].(string); ok {
		input.PolicySet = true
		input.Policy = policy
	}
	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok {
		input.KmsKeyIdentifierSet = true
		input.KmsKeyIdentifier = kms
	}
	if dlc, ok := req.Parameters["DeadLetterConfig"].(map[string]interface{}); ok {
		input.DeadLetterConfigSet = true
		input.DeadLetterConfig = &eventsstore.DeadLetterConfig{}
		if arn, ok := dlc["Arn"].(string); ok {
			input.DeadLetterConfig.Arn = arn
		}
	}
	if lc, ok := req.Parameters["LogConfig"].(map[string]interface{}); ok {
		input.LogConfigSet = true
		input.LogConfig = &BusLogConfigInput{}
		if id, ok := lc["IncludeDetail"].(string); ok {
			input.LogConfig.IncludeDetailSet = true
			input.LogConfig.IncludeDetail = id
		}
		if lvl, ok := lc["Level"].(string); ok {
			input.LogConfig.LevelSet = true
			input.LogConfig.Level = lvl
		}
	}
	return input
}

// UpdateEventBus updates an event bus.
func (s *EventsService) UpdateEventBus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := parseUpdateEventBusInput(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	eventBus, err := s.updateEventBusCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	return eventBusToUpdateMap(eventBus), nil
}

// PutPermission adds a resource policy statement to the specified event bus,
// granting the given principal permission to put events. Supports two modes:
// (1) individual parameters (Principal, StatementId, Action, Condition) and
// (2) a complete policy document via the Policy parameter.
func (s *EventsService) PutPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	busName, busNameProvided := eventBusNameParam(req)

	input := PutPermissionInput{
		BusName:         busName,
		BusNameProvided: busNameProvided,
		Principal:       request.GetStringParam(req.Parameters, "Principal"),
		StatementId:     request.GetStringParam(req.Parameters, "StatementId"),
		Action:          request.GetStringParam(req.Parameters, "Action"),
		Condition:       request.GetStringParam(req.Parameters, "Condition"),
	}
	if policyStr, ok := req.Parameters["Policy"].(string); ok {
		input.PolicySet = true
		input.Policy = policyStr
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.putPermissionCore(ctx, store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemovePermission removes a resource policy statement from the specified
// event bus identified by its StatementId.
func (s *EventsService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	busName, busNameProvided := eventBusNameParam(req)

	input := RemovePermissionInput{
		BusName:         busName,
		BusNameProvided: busNameProvided,
		StatementId:     request.GetStringParam(req.Parameters, "StatementId"),
	}
	input.RemoveAll, _ = req.Parameters["RemoveAllPermissions"].(bool)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.removePermissionCore(ctx, store, input); err != nil {
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
