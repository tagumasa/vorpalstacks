package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, roleArn, iam.ServicePrincipalStates); err != nil {
			return nil, err
		}
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags"))

	sm := &sfnstore.StateMachine{
		Name:        name,
		Definition:  definition,
		RoleArn:     roleArn,
		Type:        smType,
		Description: description,
		Tags:        tags,
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
	}
	if roleArn != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForService(ctx, roleArn, iam.ServicePrincipalStates); err != nil {
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
