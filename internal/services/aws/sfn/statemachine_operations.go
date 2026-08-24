package sfn

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
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
	return s.validateStateMachineDefinitionCore(ValidateStateMachineDefinitionInput{
		Definition: request.GetParamLowerFirst(req.Parameters, "definition"),
		SMType:     request.GetParamLowerFirst(req.Parameters, "type"),
		Severity:   request.GetParamLowerFirst(req.Parameters, "severity"),
		MaxResults: int32(request.GetIntParam(req.Parameters, "maxResults")),
	})
}

// DescribeStateMachine returns the details of a state machine.
func (s *StepFunctionService) DescribeStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeStateMachineCore(ctx, store, DescribeStateMachineInput{
		StateMachineArn: request.GetParamLowerFirst(req.Parameters, "stateMachineArn"),
		IncludedData:    request.GetParamLowerFirst(req.Parameters, "includedData"),
	})
}

// DescribeStateMachineForExecution retrieves the state machine associated with an execution.
func (s *StepFunctionService) DescribeStateMachineForExecution(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeStateMachineForExecutionCore(ctx, store, DescribeStateMachineForExecutionInput{
		ExecutionArn: request.GetParamLowerFirst(req.Parameters, "executionArn"),
		IncludedData: request.GetParamLowerFirst(req.Parameters, "includedData"),
	})
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

// UpdateStateMachine updates an existing state machine. The type of a
// state machine is fixed at creation and revisionId is not an input member
// of this operation, so neither request parameter is honoured.
func (s *StepFunctionService) UpdateStateMachine(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamLowerFirst(req.Parameters, "stateMachineArn")
	definition := request.GetParamLowerFirst(req.Parameters, "definition")
	roleArn := request.GetParamLowerFirst(req.Parameters, "roleArn")
	publish := false
	if v, ok := req.Parameters["publish"]; ok {
		if vBool, ok := v.(bool); ok {
			publish = vBool
		}
	}
	versionDescription := request.GetParamLowerFirst(req.Parameters, "versionDescription")

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
		LoggingConfiguration: loggingConfig,
		EncryptionConfig:     encryptionConfig,
		TracingConfig:        tracingConfig,
		Publish:              publish,
		VersionDescription:   versionDescription,
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
