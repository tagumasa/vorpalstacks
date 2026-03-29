package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/iam"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arncommon "vorpalstacks/internal/utils/aws/arn"
)

// CreateStateMachine creates a new state machine.
func (s *StepFunctionService) CreateStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "name")
	definition := request.GetParamLowerFirst(req.Parameters, "definition")
	roleArn := request.GetParamLowerFirst(req.Parameters, "roleArn")
	smType := request.GetParamLowerFirst(req.Parameters, "type")
	description := request.GetParamLowerFirst(req.Parameters, "description")

	if name == "" {
		return nil, NewInvalidName("State Machine name is required")
	}
	if definition == "" {
		return nil, NewInvalidDefinitionException("State Machine definition is required")
	}

	if smType == "" {
		smType = "STANDARD"
	}

	if roleArn != "" {
		if err := validateStateMachineRole(ctx, reqCtx, roleArn); err != nil {
			return nil, err
		}
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags"))

	if err := validateDefinitionJSONataFields(definition); err != nil {
		return nil, err
	}

	sm := &sfnstore.StateMachine{
		Name:               name,
		Definition:         definition,
		RoleArn:            roleArn,
		Type:               smType,
		Description:        description,
		Tags:               tags,
		VariableReferences: extractVariableReferences(definition),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stateMachineArn": sm.StateMachineArn,
		"creationDate":    sm.CreationDate.Unix(),
	}, nil
}

// DeleteStateMachine deletes a state machine.
func (s *StepFunctionService) DeleteStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteStateMachine(ctx, arn); err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ValidateStateMachineDefinition validates a state machine definition without creating it.
func (s *StepFunctionService) ValidateStateMachineDefinition(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	definition := request.GetParamLowerFirst(req.Parameters, "definition")
	if definition == "" {
		return nil, NewInvalidDefinitionException("definition is required")
	}

	var def map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return map[string]interface{}{
			"result":      "FAIL",
			"diagnostics": []map[string]string{{"severity": "ERROR", "code": "InvalidDefinition", "message": "definition is not valid JSON: " + err.Error()}},
		}, nil
	}

	if _, ok := def["StartAt"]; !ok {
		return map[string]interface{}{
			"result":      "FAIL",
			"diagnostics": []map[string]string{{"severity": "ERROR", "code": "MissingStartAt", "message": "State machine definition must include 'StartAt'"}},
		}, nil
	}

	states, ok := def["States"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"result":      "FAIL",
			"diagnostics": []map[string]string{{"severity": "ERROR", "code": "MissingStates", "message": "State machine definition must include 'States'"}},
		}, nil
	}

	startAt, _ := def["StartAt"].(string)
	if _, exists := states[startAt]; !exists {
		return map[string]interface{}{
			"result":      "FAIL",
			"diagnostics": []map[string]string{{"severity": "ERROR", "code": "InvalidStartAt", "message": fmt.Sprintf("StartAt '%s' does not reference a valid state", startAt)}},
		}, nil
	}

	return map[string]interface{}{
		"result":      "OK",
		"diagnostics": []map[string]string{},
	}, nil
}

// DescribeStateMachine returns the details of a state machine.
func (s *StepFunctionService) DescribeStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sm, err := store.GetStateMachine(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
		}
		return nil, err
	}

	response := stateMachineToResponse(sm)
	if sm.Tags != nil {
		response["tags"] = sm.Tags
	}
	refs := extractVariableReferences(sm.Definition)
	if len(refs) > 0 {
		response["variableReferences"] = refs
	}

	return response, nil
}

// DescribeStateMachineForExecution retrieves the state machine associated with an execution.
func (s *StepFunctionService) DescribeStateMachineForExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	executionArn := request.GetParamLowerFirst(req.Parameters, "executionArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	exec, err := store.GetExecution(ctx, executionArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + executionArn)
		}
		return nil, err
	}

	sm, err := store.GetStateMachine(ctx, exec.StateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + exec.StateMachineArn)
		}
		return nil, err
	}

	response := stateMachineToResponse(sm)
	if sm.Tags != nil {
		response["tags"] = sm.Tags
	}

	return response, nil
}

// GetStateMachine returns the details of a state machine (alias for DescribeStateMachine).
func (s *StepFunctionService) GetStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeStateMachine(ctx, reqCtx, req)
}

// ListStateMachines returns a list of state machines.
func (s *StepFunctionService) ListStateMachines(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit == 0 {
		limit = 100
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListStateMachines(ctx, limit, nextToken)
	if err != nil {
		return nil, err
	}

	stateMachines := make([]map[string]interface{}, len(result.StateMachines))
	for i, sm := range result.StateMachines {
		stateMachines[i] = map[string]interface{}{
			"stateMachineArn": sm.StateMachineArn,
			"name":            sm.Name,
			"type":            sm.Type,
			"creationDate":    sm.CreationDate.Unix(),
		}
	}

	response := map[string]interface{}{
		"stateMachines": stateMachines,
	}

	if result.NextToken != "" {
		response["nextToken"] = result.NextToken
	}

	return response, nil
}

// UpdateStateMachine updates an existing state machine.
func (s *StepFunctionService) UpdateStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	definition := request.GetParamLowerFirst(req.Parameters, "definition")
	roleArn := request.GetParamLowerFirst(req.Parameters, "roleArn")
	description := request.GetParamLowerFirst(req.Parameters, "description")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sm, err := store.GetStateMachine(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
		}
		return nil, err
	}

	if definition != "" {
		sm.Definition = definition
		if err := validateDefinitionJSONataFields(definition); err != nil {
			return nil, err
		}
		sm.VariableReferences = extractVariableReferences(definition)
	}
	if roleArn != "" {
		if err := validateStateMachineRole(ctx, reqCtx, roleArn); err != nil {
			return nil, err
		}
		sm.RoleArn = roleArn
	}
	if description != "" {
		sm.Description = description
	}
	sm.UpdateDate = time.Now().UTC()

	if err := store.UpdateStateMachine(ctx, sm); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stateMachineArn": sm.StateMachineArn,
		"updateDate":      sm.UpdateDate.Unix(),
	}, nil
}

// StartExecution starts an execution of a state machine.
func (s *StepFunctionService) StartExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stateMachineArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	name := request.GetParamLowerFirst(req.Parameters, "name")
	input := request.GetParamLowerFirst(req.Parameters, "input")
	traceHeader := request.GetParamLowerFirst(req.Parameters, "traceHeader")

	if stateMachineArn == "" {
		stateMachineArn = request.GetParamLowerFirst(req.Parameters, "StateMachineArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sm, err := store.GetStateMachine(ctx, stateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + stateMachineArn)
		}
		return nil, err
	}

	if name == "" {
		name = generateExecutionName()
	}

	executionArn := fmt.Sprintf("arn:aws:states:%s:%s:execution:%s:%s",
		reqCtx.GetRegion(), s.accountID,
		arncommon.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, input, traceHeader)
	exec.ExecutionArn = executionArn

	if err := store.CreateExecution(ctx, exec); err != nil {
		return nil, err
	}

	sqsStore := s.sqsStore
	if sqsStore == nil {
		sqsStore = reqCtx.GetSQSStore()
	}
	snsStore := s.snsStore
	if snsStore == nil {
		snsStore = reqCtx.GetSNSStore()
	}
	eventsStore := s.eventsStore
	if eventsStore == nil {
		if ebStore := reqCtx.GetEventBridgeStore(); ebStore != nil {
			if concrete, ok := ebStore.(*eventsstore.EventsStore); ok {
				eventsStore = concrete
			}
		}
	}

	executor := NewExecutorWithStores(store, s.lambdaInvoker, sqsStore, snsStore, eventsStore, s.accountID, reqCtx.GetRegion())
	execCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(executionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(executionArn)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in execution", logs.String("arn", executionArn), logs.Any("panic", r))
			}
		}()
		_ = executor.ExecuteStateMachine(execCtx, exec)
	}()

	return map[string]interface{}{
		"executionArn":    exec.ExecutionArn,
		"startDate":       exec.StartDate.Unix(),
		"stateMachineArn": exec.StateMachineArn,
	}, nil
}

// StopExecution stops a running execution of a state machine.
func (s *StepFunctionService) StopExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "executionArn")
	errorMsg := request.GetParamLowerFirst(req.Parameters, "error")
	cause := request.GetParamLowerFirst(req.Parameters, "cause")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exec, err := store.GetExecution(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + arn)
		}
		return nil, err
	}

	terminalStates := map[string]bool{
		"SUCCEEDED": true,
		"FAILED":    true,
		"TIMED_OUT": true,
		"ABORTED":   true,
	}
	if terminalStates[exec.Status] {
		stopDate := time.Now().UTC()
		return map[string]interface{}{
			"stopDate": stopDate.Unix(),
		}, nil
	}

	store.CancelExecution(arn)

	exec.Status = "ABORTED"
	exec.StopDate = time.Now().UTC()
	exec.Error = errorMsg
	exec.Cause = cause

	if err := store.UpdateExecution(ctx, exec); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"stopDate": exec.StopDate.Unix(),
	}, nil
}

// DescribeExecution returns the details of an execution.
func (s *StepFunctionService) DescribeExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "executionArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exec, err := store.GetExecution(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrExecutionNotFound) {
			return nil, NewExecutionDoesNotExist("Execution Does not exist: " + arn)
		}
		return nil, err
	}

	return executionToResponse(exec), nil
}

// ListExecutions returns a list of executions for a state machine.
func (s *StepFunctionService) ListExecutions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stateMachineArn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	statusFilter := request.GetParamLowerFirst(req.Parameters, "statusFilter")
	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit == 0 {
		limit = 100
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListExecutions(ctx, stateMachineArn, statusFilter, limit, nextToken)
	if err != nil {
		return nil, err
	}

	executions := make([]map[string]interface{}, len(result.Executions))
	for i, exec := range result.Executions {
		executions[i] = map[string]interface{}{
			"executionArn":    exec.ExecutionArn,
			"stateMachineArn": exec.StateMachineArn,
			"name":            exec.Name,
			"status":          exec.Status,
			"startDate":       exec.StartDate.Unix(),
		}
		if !exec.StopDate.IsZero() {
			executions[i]["stopDate"] = exec.StopDate.Unix()
		}
	}

	response := map[string]interface{}{
		"executions": executions,
	}

	if result.NextToken != "" {
		response["nextToken"] = result.NextToken
	}

	return response, nil
}

// GetExecutionHistory returns the history of an execution.
func (s *StepFunctionService) GetExecutionHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "executionArn")
	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit == 0 {
		limit = 100
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	events, nextTokenResult, err := store.GetExecutionHistory(ctx, arn, limit, nextToken)
	if err != nil {
		return nil, err
	}

	history := make([]map[string]interface{}, len(events))
	for i, event := range events {
		history[i] = historyEventToResponse(event)
	}

	response := map[string]interface{}{
		"events": history,
	}

	if nextTokenResult != "" {
		response["nextToken"] = nextTokenResult
	}

	return response, nil
}

// CreateActivity creates a new activity.
func (s *StepFunctionService) CreateActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "name")

	activity := &sfnstore.Activity{
		Name: name,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateActivity(ctx, activity); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"activityArn":  activity.ActivityArn,
		"creationDate": activity.CreationDate.Unix(),
	}, nil
}

// DeleteActivity deletes an activity.
func (s *StepFunctionService) DeleteActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "activityArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteActivity(ctx, arn); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeActivity returns the details of an activity.
func (s *StepFunctionService) DescribeActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "activityArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	activity, err := store.GetActivity(ctx, arn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			return nil, NewActivityDoesNotExist("Activity Does not exist: " + arn)
		}
		return nil, err
	}

	return activityToResponse(activity), nil
}

// GetActivity returns the details of an activity (alias for DescribeActivity).
func (s *StepFunctionService) GetActivity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeActivity(ctx, reqCtx, req)
}

// ListActivities returns a list of activities.
func (s *StepFunctionService) ListActivities(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if limit == 0 {
		limit = 100
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListActivities(ctx, limit, nextToken)
	if err != nil {
		return nil, err
	}

	activities := make([]map[string]interface{}, len(result.Activities))
	for i, activity := range result.Activities {
		activities[i] = activityToResponse(activity)
	}

	response := map[string]interface{}{
		"activities": activities,
	}

	if result.NextToken != "" {
		response["nextToken"] = result.NextToken
	}

	return response, nil
}

// TagResource adds tags to a state machine.
func (s *StepFunctionService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "resourceArn")
	tags := tagutil.ParseTags(req.Parameters, "tags")

	if strings.Contains(arn, ":stateMachine:") {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		sm, err := store.GetStateMachine(ctx, arn)
		if err != nil {
			if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
				return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
			}
			return nil, err
		}
		if sm.Tags == nil {
			sm.Tags = make(map[string]string)
		}
		for k, v := range tagutil.ToMap(tags) {
			sm.Tags[k] = v
		}
		if err := store.UpdateStateMachine(ctx, sm); err != nil {
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a state machine.
func (s *StepFunctionService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "resourceArn")
	tagKeys := tagutil.ParseTagKeysAsSlice(req.Parameters, "tagKeys")

	if strings.Contains(arn, ":stateMachine:") {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		sm, err := store.GetStateMachine(ctx, arn)
		if err != nil {
			if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
				return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
			}
			return nil, err
		}
		if sm.Tags != nil {
			for _, k := range tagKeys {
				delete(sm.Tags, k)
			}
		}
		if err := store.UpdateStateMachine(ctx, sm); err != nil {
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource returns the tags for a state machine.
func (s *StepFunctionService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "resourceArn")

	tags := make([]map[string]string, 0)

	if strings.Contains(arn, ":stateMachine:") {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		sm, err := store.GetStateMachine(ctx, arn)
		if err != nil {
			if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
				return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
			}
			return nil, err
		}
		if sm.Tags != nil {
			tags = tagutil.ConvertToMapSlice(tagutil.MapToTags(sm.Tags))
		}
	}

	return map[string]interface{}{
		"tags": tags,
	}, nil
}

func generateExecutionName() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func validateStateMachineRole(ctx context.Context, reqCtx *request.RequestContext, roleArn string) error {
	validator := reqCtx.GetIAMValidator()
	return validator.ValidateRoleForServiceWithErrors(ctx, roleArn, iam.ServicePrincipalStates, &iam.RoleErrorFactories{
		RoleNotFoundError:        sfnRoleNotFoundError,
		RoleCannotBeAssumedError: sfnRoleCannotBeAssumedError,
		InvalidArnError:          sfnInvalidRoleArnError,
	})
}

func sfnRoleNotFoundError(roleArn string) error {
	return &SFNError{awserrors.NewAWSError("InvalidParameterException", fmt.Sprintf("Role Arn is not valid for State Machine: %s", roleArn), 400)}
}

func sfnRoleCannotBeAssumedError(roleArn string) error {
	return &SFNError{awserrors.NewAWSError("AccessDeniedException", fmt.Sprintf("Role %s is invalid or cannot be assumed.", roleArn), 400)}
}

func sfnInvalidRoleArnError(roleArn string) error {
	return &SFNError{awserrors.NewAWSError("InvalidArn", fmt.Sprintf("Invalid Role Arn: %s", roleArn), 400)}
}

func validateDefinitionJSONataFields(definition string) error {
	var def map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return nil
	}

	topQL, _ := def["QueryLanguage"].(string)
	if topQL == "" {
		topQL = "JSONPath"
	}

	states, ok := def["States"].(map[string]interface{})
	if !ok {
		return nil
	}

	for stateName, stateData := range states {
		stateMap, ok := stateData.(map[string]interface{})
		if !ok {
			continue
		}

		stateType, _ := stateMap["Type"].(string)
		stateQL, _ := stateMap["QueryLanguage"].(string)
		if stateQL == "" {
			stateQL = topQL
		}

		if stateQL == "JSONata" {
			jsonPathOnlyFields := getJSONPathOnlyFields(stateType, stateMap)
			if len(jsonPathOnlyFields) > 0 {
				return NewInvalidDefinitionException(fmt.Sprintf(
					"State '%s' uses JSONata QueryLanguage but contains JSONPath-only field(s): %s",
					stateName, strings.Join(jsonPathOnlyFields, ", ")))
			}
		} else {
			jsonataOnlyFields := getJSONataOnlyFields(stateType, stateMap)
			if len(jsonataOnlyFields) > 0 {
				return NewInvalidDefinitionException(fmt.Sprintf(
					"State '%s' uses JSONPath QueryLanguage but contains JSONata-only field(s): %s",
					stateName, strings.Join(jsonataOnlyFields, ", ")))
			}
		}
	}

	return nil
}

func getJSONPathOnlyFields(stateType string, stateMap map[string]interface{}) []string {
	var forbidden []string

	switch stateType {
	case "Pass", "Task", "Parallel", "Map", "Succeed":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
		if _, ok := stateMap["OutputPath"]; ok {
			forbidden = append(forbidden, "OutputPath")
		}
		if _, ok := stateMap["Parameters"]; ok {
			forbidden = append(forbidden, "Parameters")
		}
	case "Choice":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
	case "Wait":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
		if _, ok := stateMap["OutputPath"]; ok {
			forbidden = append(forbidden, "OutputPath")
		}
		if _, ok := stateMap["SecondsPath"]; ok {
			forbidden = append(forbidden, "SecondsPath")
		}
		if _, ok := stateMap["TimestampPath"]; ok {
			forbidden = append(forbidden, "TimestampPath")
		}
	}

	switch stateType {
	case "Task":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
		if _, ok := stateMap["TimeoutSecondsPath"]; ok {
			forbidden = append(forbidden, "TimeoutSecondsPath")
		}
		if _, ok := stateMap["HeartbeatSecondsPath"]; ok {
			forbidden = append(forbidden, "HeartbeatSecondsPath")
		}
	case "Pass":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
	case "Map":
		if _, ok := stateMap["ItemsPath"]; ok {
			forbidden = append(forbidden, "ItemsPath")
		}
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
	case "Parallel":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
	case "Fail":
		if _, ok := stateMap["CausePath"]; ok {
			forbidden = append(forbidden, "CausePath")
		}
		if _, ok := stateMap["ErrorPath"]; ok {
			forbidden = append(forbidden, "ErrorPath")
		}
	}

	return forbidden
}

func getJSONataOnlyFields(stateType string, stateMap map[string]interface{}) []string {
	var forbidden []string

	if _, ok := stateMap["Arguments"]; ok {
		if stateType == "Task" || stateType == "Parallel" {
			forbidden = append(forbidden, "Arguments")
		}
	}

	if _, ok := stateMap["Items"]; ok {
		if stateType == "Map" {
			forbidden = append(forbidden, "Items")
		}
	}

	if _, ok := stateMap["Condition"]; ok {
		if stateType == "Choice" {
			forbidden = append(forbidden, "Condition")
		}
	}

	if stateType == "Choice" {
		if choices, ok := stateMap["Choices"].([]interface{}); ok {
			for _, choice := range choices {
				if choiceMap, ok := choice.(map[string]interface{}); ok {
					if _, ok := choiceMap["Condition"]; ok {
						forbidden = append(forbidden, "Condition")
						break
					}
				}
			}
		}
	}

	return forbidden
}

var variableRefRegex = regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)`)

func extractVariableReferences(definition string) map[string][]string {
	var def map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return nil
	}

	states, ok := def["States"].(map[string]interface{})
	if !ok {
		return nil
	}

	assignedVars := make(map[string]bool)
	for _, stateData := range states {
		stateMap, ok := stateData.(map[string]interface{})
		if !ok {
			continue
		}
		if assign, ok := stateMap["Assign"].(map[string]interface{}); ok {
			for name := range assign {
				clean := strings.TrimPrefix(name, "$")
				assignedVars[clean] = true
			}
		}
		if choices, ok := stateMap["Choices"].([]interface{}); ok {
			for _, choice := range choices {
				if choiceMap, ok := choice.(map[string]interface{}); ok {
					if assign, ok := choiceMap["Assign"].(map[string]interface{}); ok {
						for name := range assign {
							clean := strings.TrimPrefix(name, "$")
							assignedVars[clean] = true
						}
					}
				}
			}
		}
	}

	result := make(map[string][]string)

	for stateName, stateData := range states {
		stateMap, ok := stateData.(map[string]interface{})
		if !ok {
			continue
		}

		refs := collectVariableRefsFromState(stateMap, assignedVars)
		if len(refs) > 0 {
			result[stateName] = refs
		}
	}

	return result
}

func collectVariableRefsFromState(stateMap map[string]interface{}, assignedVars map[string]bool) []string {
	seen := make(map[string]bool)
	var refs []string

	for _, field := range []string{"Assign", "Output", "Arguments", "Items", "Condition"} {
		if field == "Condition" {
			if choices, ok := stateMap["Choices"].([]interface{}); ok {
				for _, choice := range choices {
					if choiceMap, ok := choice.(map[string]interface{}); ok {
						scanValueForVariableRefs(choiceMap["Condition"], seen, &refs, assignedVars)
					}
				}
			}
			continue
		}
		if val, ok := stateMap[field]; ok {
			scanValueForVariableRefs(val, seen, &refs, assignedVars)
		}
	}

	return refs
}

var jsonataBuiltins = map[string]bool{
	"states": true, "context": true,
	"abs": true, "count": true, "sum": true, "max": true, "min": true, "average": true,
	"string": true, "substring": true, "length": true, "uppercase": true, "lowercase": true,
	"trim": true, "pad": true, "contains": true, "split": true, "join": true,
	"match": true, "replace": true, "base64encode": true, "base64decode": true,
	"number": true, "round": true, "floor": true, "ceil": true, "sqrt": true, "power": true,
	"random": true, "boolean": true, "not": true, "exists": true, "type": true,
	"each": true, "filter": true, "flatten": true, "keys": true, "lookup": true,
	"map": true, "merge": true, "reverse": true, "sort": true, "spread": true,
	"sift": true, "distinct": true, "single": true, "tail": true, "append": true,
	"errors": true, "fromMillis": true, "toMillis": true, "millis": true,
	"now": true, "uuid": true, "parse": true, "hash": true, "partition": true, "range": true,
	"description": true, "url": true, "encodeUrlComponent": true,
	"assert": true, "error": true, "order": true,
}

func scanValueForVariableRefs(v interface{}, seen map[string]bool, refs *[]string, assignedVars map[string]bool) {
	switch val := v.(type) {
	case string:
		if strings.HasPrefix(val, "{%") && strings.HasSuffix(val, "%}") {
			expr := strings.TrimPrefix(strings.TrimSuffix(val, "%}"), "{%")
			for _, match := range variableRefRegex.FindAllStringSubmatch(expr, -1) {
				name := match[1]
				if !seen[name] && (!jsonataBuiltins[name] || assignedVars[name]) {
					seen[name] = true
					*refs = append(*refs, name)
				}
			}
		}
	case map[string]interface{}:
		for _, child := range val {
			scanValueForVariableRefs(child, seen, refs, assignedVars)
		}
	case []interface{}:
		for _, child := range val {
			scanValueForVariableRefs(child, seen, refs, assignedVars)
		}
	}
}
