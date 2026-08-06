package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

// CreateStateMachine creates a new state machine.
func (s *StepFunctionService) CreateStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "name")
	definition := request.GetParamLowerFirst(req.Parameters, "definition")
	roleArn := request.GetParamLowerFirst(req.Parameters, "roleArn")
	smType := request.GetParamLowerFirst(req.Parameters, "type")
	publish := false
	if v, ok := req.Parameters["publish"]; ok {
		if vBool, ok := v.(bool); ok {
			publish = vBool
		}
	}
	versionDescription := request.GetParamLowerFirst(req.Parameters, "versionDescription")

	// Role validation requires the IAM validator from the request context.
	if err := validateStateMachineRole(ctx, reqCtx, roleArn); err != nil {
		return nil, err
	}

	loggingConfig, err := parseLoggingConfigurationFromJSON(req.Parameters["loggingConfiguration"])
	if err != nil {
		return nil, err
	}

	encryptionConfig, err := parseEncryptionConfigurationFromJSON(req.Parameters["encryptionConfiguration"])
	if err != nil {
		return nil, err
	}

	tracingConfig, err := parseTracingConfigurationFromJSON(req.Parameters["tracingConfiguration"])
	if err != nil {
		return nil, err
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags"))

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createStateMachineCore(ctx, store, CreateStateMachineInput{
		Name:                 name,
		Definition:           definition,
		RoleArn:              roleArn,
		Type:                 smType,
		LoggingConfiguration: loggingConfig,
		EncryptionConfig:     encryptionConfig,
		TracingConfig:        tracingConfig,
		Tags:                 tags,
		Publish:              publish,
		VersionDescription:   versionDescription,
	})
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"stateMachineArn": result.StateMachineArn,
		"creationDate":    result.CreationDate.Unix(),
	}
	if result.StateMachineVersionArn != "" {
		resp["stateMachineVersionArn"] = result.StateMachineVersionArn
	}
	return resp, nil
}

func generateRevisionId() string {
	return uuid.New().String()
}

// DeleteStateMachine deletes a state machine.
func (s *StepFunctionService) DeleteStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteStateMachineCore(ctx, store, DeleteStateMachineInput{
		StateMachineArn: arn,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ValidateStateMachineDefinition validates a state machine definition without creating it.
func (s *StepFunctionService) ValidateStateMachineDefinition(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	definition := request.GetParamLowerFirst(req.Parameters, "definition")
	if err := validateDefinitionJSON(definition); err != nil {
		return nil, err
	}

	severity, err := validateSeverity(request.GetParamLowerFirst(req.Parameters, "severity"))
	if err != nil {
		return nil, err
	}

	smType, err := validateStateMachineType(request.GetParamLowerFirst(req.Parameters, "type"))
	if err != nil {
		return nil, err
	}
	_ = smType

	maxResults := int32(request.GetIntParam(req.Parameters, "maxResults"))
	if err := validateMaxResults(maxResults, 0, 100, "maxResults"); err != nil {
		return nil, err
	}
	if maxResults == 0 {
		maxResults = 100
	}

	var def map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return map[string]interface{}{
			"result":      "FAIL",
			"diagnostics": []map[string]string{{"severity": "ERROR", "code": "InvalidDefinition", "message": "definition is not valid JSON: " + err.Error()}},
		}, nil
	}

	diagnostics := []map[string]string{}

	if _, ok := def["StartAt"]; !ok {
		diagnostics = append(diagnostics, map[string]string{"severity": "ERROR", "code": "MissingStartAt", "message": "State machine definition must include 'StartAt'"})
	}

	states, ok := def["States"].(map[string]interface{})
	if !ok {
		diagnostics = append(diagnostics, map[string]string{"severity": "ERROR", "code": "MissingStates", "message": "State machine definition must include 'States'"})
	}

	startAt, _ := def["StartAt"].(string)
	if states != nil && startAt != "" {
		if _, exists := states[startAt]; !exists {
			diagnostics = append(diagnostics, map[string]string{"severity": "ERROR", "code": "InvalidStartAt", "message": fmt.Sprintf("StartAt '%s' does not reference a valid state", startAt)})
		}
	}

	// Per-state validation: check Type enum and required fields per type.
	validTypes := map[string]bool{
		"Pass": true, "Task": true, "Choice": true, "Wait": true,
		"Succeed": true, "Fail": true, "Parallel": true, "Map": true,
	}
	if states != nil {
		for name, rawState := range states {
			stateMap, ok := rawState.(map[string]interface{})
			if !ok {
				diagnostics = append(diagnostics, map[string]string{"severity": "ERROR", "code": "InvalidState", "message": fmt.Sprintf("State '%s' is not an object", name)})
				continue
			}
			stateType, _ := stateMap["Type"].(string)
			if stateType == "" {
				diagnostics = append(diagnostics, map[string]string{"severity": "ERROR", "code": "MissingStateType", "message": fmt.Sprintf("State '%s' is missing Type", name)})
				continue
			}
			if !validTypes[stateType] {
				diagnostics = append(diagnostics, map[string]string{"severity": "ERROR", "code": "InvalidStateType", "message": fmt.Sprintf("State '%s' has invalid Type '%s'", name, stateType)})
				continue
			}
			if next, ok := stateMap["Next"].(string); ok {
				if _, exists := states[next]; !exists {
					diagnostics = append(diagnostics, map[string]string{"severity": "ERROR", "code": "InvalidNext", "message": fmt.Sprintf("State '%s' references unknown Next state '%s'", name, next)})
				}
			}
			if stateType == "Task" {
				if _, ok := stateMap["Resource"]; !ok {
					diagnostics = append(diagnostics, map[string]string{"severity": "ERROR", "code": "MissingResource", "message": fmt.Sprintf("Task state '%s' is missing Resource", name)})
				}
			}
			if stateType == "Fail" {
				if _, ok := stateMap["Cause"]; !ok {
					if _, ok2 := stateMap["Error"]; !ok2 {
						diagnostics = append(diagnostics, map[string]string{"severity": "WARNING", "code": "MissingFailDetails", "message": fmt.Sprintf("Fail state '%s' should specify Cause or Error", name)})
					}
				}
			}
		}
	}

	if severity == "ERROR" {
		filtered := []map[string]string{}
		for _, d := range diagnostics {
			if d["severity"] == "ERROR" {
				filtered = append(filtered, d)
			}
		}
		diagnostics = filtered
	}

	truncated := false
	if maxResults > 0 && len(diagnostics) > int(maxResults) {
		diagnostics = diagnostics[:maxResults]
		truncated = true
	}

	result := "OK"
	if len(diagnostics) > 0 {
		result = "FAIL"
	}

	return map[string]interface{}{
		"result":      result,
		"diagnostics": diagnostics,
		"truncated":   truncated,
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
	tagsMap, _ := store.ListAsSlice(sm.StateMachineArn)
	if len(tagsMap) > 0 {
		response["tags"] = tagsMap
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
	tagsMap, _ := store.ListAsSlice(sm.StateMachineArn)
	if len(tagsMap) > 0 {
		response["tags"] = tagsMap
	}

	return response, nil
}

// GetStateMachine returns the details of a state machine (alias for DescribeStateMachine).
func (s *StepFunctionService) GetStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeStateMachine(ctx, reqCtx, req)
}

// ListStateMachines returns a list of state machines.
func (s *StepFunctionService) ListStateMachines(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listStateMachinesCore(ctx, store, ListStateMachinesInput{
		MaxResults: limit,
		NextToken:  nextToken,
	})
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

	resp := map[string]interface{}{
		"stateMachines": stateMachines,
	}

	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

// UpdateStateMachine updates an existing state machine.
func (s *StepFunctionService) UpdateStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	definition := request.GetParamLowerFirst(req.Parameters, "definition")
	roleArn := request.GetParamLowerFirst(req.Parameters, "roleArn")
	smType := request.GetParamLowerFirst(req.Parameters, "type")
	publish := false
	if v, ok := req.Parameters["publish"]; ok {
		if vBool, ok := v.(bool); ok {
			publish = vBool
		}
	}
	versionDescription := request.GetParamLowerFirst(req.Parameters, "versionDescription")
	revisionId := request.GetParamLowerFirst(req.Parameters, "revisionId")

	// Role validation requires the IAM validator from the request context.
	if roleArn != "" {
		if err := validateStateMachineRole(ctx, reqCtx, roleArn); err != nil {
			return nil, err
		}
	}

	loggingConfig, err := parseLoggingConfigurationFromJSON(req.Parameters["loggingConfiguration"])
	if err != nil {
		return nil, err
	}

	encryptionConfig, err := parseEncryptionConfigurationFromJSON(req.Parameters["encryptionConfiguration"])
	if err != nil {
		return nil, err
	}

	tracingConfig, err := parseTracingConfigurationFromJSON(req.Parameters["tracingConfiguration"])
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.updateStateMachineCore(ctx, store, UpdateStateMachineInput{
		StateMachineArn:      arn,
		Definition:           definition,
		DefinitionProvided:   definition != "",
		RoleArn:              roleArn,
		RoleArnProvided:      roleArn != "",
		Type:                 smType,
		TypeProvided:         smType != "",
		LoggingConfiguration: loggingConfig,
		EncryptionConfig:     encryptionConfig,
		TracingConfig:        tracingConfig,
		Publish:              publish,
		VersionDescription:   versionDescription,
		RevisionId:           revisionId,
	})
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"stateMachineArn": result.StateMachineArn,
		"updateDate":      result.UpdateDate.Unix(),
		"revisionId":      result.RevisionId,
	}
	if result.StateMachineVersionArn != "" {
		resp["stateMachineVersionArn"] = result.StateMachineVersionArn
	}
	return resp, nil
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

	executionArn := arnutil.NewARNBuilder(s.accountID, reqCtx.GetRegion()).StepFunctions().Execution(arnutil.ExtractStateMachineNameFromARN(sm.StateMachineArn), name)

	exec := sfnstore.NewExecution(sm.StateMachineArn, name, input, traceHeader)
	exec.ExecutionArn = executionArn

	if err := store.CreateExecution(ctx, exec); err != nil {
		if errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return nil, awserrors.NewAWSError("ExecutionAlreadyExists", "An execution with the same name already exists: "+executionArn, 400)
		}
		return nil, err
	}

	executor := NewExecutorWithStores(store, s.bus, s.accountID, reqCtx.GetRegion())
	execCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(executionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(executionArn)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in execution", logs.String("arn", executionArn), logs.Any("panic", r))
				exec.Status = "FAILED"
				exec.Error = "States.InternalError"
				exec.Cause = fmt.Sprintf("internal panic: %v", r)
				exec.StopDate = time.Now().UTC()
				_ = store.UpdateExecution(context.Background(), exec)
			}
		}()
		if err := executor.ExecuteStateMachine(execCtx, exec); err != nil {
			logs.Error("sfn: execution error", logs.String("arn", executionArn), logs.Err(err))
		}
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

	if isTerminalStatus(exec.Status) {
		return map[string]interface{}{
			"stopDate": exec.StopDate.Unix(),
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
	mapRunArn := request.GetParamLowerFirst(req.Parameters, "mapRunArn")
	redriveFilter := request.GetParamLowerFirst(req.Parameters, "redriveFilter")
	if redriveFilter != "" && redriveFilter != "REDRIVEN" && redriveFilter != "NOT_REDRIVEN" {
		return nil, NewInvalidParameterValue("redriveFilter must be REDRIVEN or NOT_REDRIVEN, got " + redriveFilter)
	}
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.GetStateMachine(ctx, stateMachineArn)
	if err != nil {
		return nil, NewStateMachineDoesNotExist("State Machine Does not exist: " + stateMachineArn)
	}
	result, err := store.ListExecutions(ctx, stateMachineArn, statusFilter, mapRunArn, redriveFilter, limit, nextToken)
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
		if exec.MapRunArn != "" {
			executions[i]["mapRunArn"] = exec.MapRunArn
		}
		if exec.ItemCount != 0 {
			executions[i]["itemCount"] = exec.ItemCount
		}
		if exec.RedriveCount != 0 {
			executions[i]["redriveCount"] = exec.RedriveCount
		}
		if !exec.RedriveDate.IsZero() {
			executions[i]["redriveDate"] = exec.RedriveDate.Unix()
		}
		if exec.StateMachineAliasArn != "" {
			executions[i]["stateMachineAliasArn"] = exec.StateMachineAliasArn
		}
		if exec.StateMachineVersionArn != "" {
			executions[i]["stateMachineVersionArn"] = exec.StateMachineVersionArn
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
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "nextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.GetExecution(ctx, arn)
	if err != nil {
		return nil, NewExecutionDoesNotExist("Execution Does not exist: " + arn)
	}
	events, nextTokenResult, err := store.GetExecutionHistory(ctx, arn, limit, nextToken)
	if err != nil {
		return nil, err
	}

	includeExecutionData := true
	if v, ok := req.Parameters["includeExecutionData"]; ok {
		if vBool, ok := v.(bool); ok {
			includeExecutionData = vBool
		}
	}

	reverseOrder := false
	if v, ok := req.Parameters["reverseOrder"]; ok {
		if vBool, ok := v.(bool); ok {
			reverseOrder = vBool
		}
	}

	history := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		history = append(history, historyEventToResponse(event, includeExecutionData))
	}

	if reverseOrder {
		for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
			history[i], history[j] = history[j], history[i]
		}
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

	if err := validateResourceName(name); err != nil {
		return nil, err
	}

	activity := &sfnstore.Activity{
		Name: name,
	}

	if ec, err := parseEncryptionConfigurationFromJSON(req.Parameters["encryptionConfiguration"]); err != nil {
		return nil, err
	} else {
		activity.EncryptionConfiguration = ec
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateActivity(ctx, activity); err != nil {
		return nil, err
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags"))
	if len(tags) > 0 {
		if err := store.Tag(activity.ActivityArn, tags); err != nil {
			return nil, err
		}
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
		if errors.Is(err, sfnstore.ErrActivityNotFound) {
			return nil, NewActivityDoesNotExist("Activity Does not exist: " + arn)
		}
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
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
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

func (s *StepFunctionService) tagHandlerConfig(store *sfnstore.StepFunctionStore) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:   "resourceArn",
			TagsParam:       "tags",
			TagKeysParam:    "tagKeys",
			RequireTags:     true,
			RequireTagKeys:  true,
			RequireResource: true,
		},
		ResourceKey: func(rawKey string) string { return rawKey },
		ValidateResource: func(ctx context.Context, arn string) error {
			if !strings.Contains(arn, ":stateMachine:") {
				return nil
			}
			if _, err := store.GetStateMachine(ctx, arn); err != nil {
				if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
					return NewStateMachineDoesNotExist("State Machine Does not exist: " + arn)
				}
				return err
			}
			return nil
		},
		ParseTags: func(params map[string]interface{}) []types.Tag {
			return tagutil.MapToTags(tagutil.ToMap(tagutil.ParseTags(params, "tags")))
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return tagutil.ParseTagKeysAsSlice(params, "tagKeys")
		},
		TagFunc: func(_ context.Context, resourceKey string, tagSlice []types.Tag) error {
			return store.TagFromSlice(resourceKey, tagSlice)
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			return store.ListAsSlice(resourceKey)
		},
		FormatResponse: func(tagSlice []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"tags": tagutil.ToResponseWithKeyNames(tagSlice, "key", "value"),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
	}
}

// TagResource adds tags to a state machine.
func (s *StepFunctionService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.tagHandlerConfig(store))
}

// UntagResource removes tags from a state machine.
func (s *StepFunctionService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.tagHandlerConfig(store))
}

// ListTagsForResource returns the tags for a state machine.
func (s *StepFunctionService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.tagHandlerConfig(store))
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
	return awserrors.NewAWSError("InvalidParameterException", fmt.Sprintf("Role Arn is not valid for State Machine: %s", roleArn), 400)
}

func sfnRoleCannotBeAssumedError(roleArn string) error {
	return awserrors.NewAWSError("AccessDeniedException", fmt.Sprintf("Role %s is invalid or cannot be assumed.", roleArn), 403)
}

func sfnInvalidRoleArnError(roleArn string) error {
	return awserrors.NewAWSError("InvalidArn", fmt.Sprintf("Invalid Role Arn: %s", roleArn), 400)
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
