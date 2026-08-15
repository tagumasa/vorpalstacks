package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
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
	if err := validateMaxResults(maxResults, 0, sfnstore.MaxValidateDefinitionResults, "maxResults"); err != nil {
		return nil, err
	}
	if maxResults == 0 {
		maxResults = sfnstore.MaxValidateDefinitionResults
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
	if err := validateArnRequired(arn, "stateMachineArn"); err != nil {
		return nil, err
	}

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
	if err := validateArnRequired(executionArn, "executionArn"); err != nil {
		return nil, err
	}

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

func validateStateMachineRole(ctx context.Context, reqCtx *request.RequestContext, roleArn string) error {
	validator := reqCtx.GetIAMValidator()
	return validator.ValidateRoleForServiceWithErrors(ctx, roleArn, iam.ServicePrincipalStates, &iam.RoleErrorFactories{
		RoleNotFoundError:        sfnRoleNotFoundError,
		RoleCannotBeAssumedError: sfnRoleCannotBeAssumedError,
		InvalidArnError:          sfnInvalidRoleArnError,
	})
}

func sfnRoleNotFoundError(roleArn string) error {
	// The Smithy model has no InvalidParameterException for SFN; an
	// unresolvable role is an input-constraint failure of the create call.
	return NewValidationException(fmt.Sprintf("Role Arn is not valid for State Machine: %s", roleArn))
}

func sfnRoleCannotBeAssumedError(roleArn string) error {
	// AccessDeniedException is an AWS-common auth-class error rather than a
	// Smithy-modelled SFN operation error; a role that exists but cannot be
	// assumed is an authorisation failure, not an input-constraint failure.
	return awserrors.NewAWSError("AccessDeniedException", fmt.Sprintf("Role %s is invalid or cannot be assumed.", roleArn), 403)
}

func sfnInvalidRoleArnError(roleArn string) error {
	return NewInvalidArnException(fmt.Sprintf("Invalid Role Arn: %s", roleArn))
}
