package sfn

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// TestState tests a single state within a state machine definition; the
// in-memory run lives in the TestState Core.
func (s *StepFunctionService) TestState(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	mock, err := parseTestStateMock(req.Parameters["mock"])
	if err != nil {
		return nil, err
	}
	stateConfig, err := parseTestStateConfiguration(req.Parameters["stateConfiguration"])
	if err != nil {
		return nil, err
	}

	revealSecrets := false
	if v, ok := req.Parameters["revealSecrets"].(bool); ok {
		revealSecrets = v
	}

	return s.testStateCore(ctx, store, TestStateInput{
		Definition:      request.GetParamLowerFirst(req.Parameters, "definition"),
		StateName:       request.GetParamLowerFirst(req.Parameters, "stateName"),
		Input:           request.GetParamLowerFirst(req.Parameters, "input"),
		InspectionLevel: request.GetParamLowerFirst(req.Parameters, "inspectionLevel"),
		Variables:       request.GetParamLowerFirst(req.Parameters, "variables"),
		RoleArn:         request.GetParamLowerFirst(req.Parameters, "roleArn"),
		Context:         request.GetParamLowerFirst(req.Parameters, "context"),
		RevealSecrets:   revealSecrets,
		Mock:            mock,
		StateConfig:     stateConfig,
	})
}

// parseTestStateMock reads the optional MockInput object from the wire.
func parseTestStateMock(raw interface{}) (*TestStateMock, error) {
	if raw == nil {
		return nil, nil
	}
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return nil, NewValidationException("mock must be an object")
	}
	mock := &TestStateMock{}
	if v, ok := obj["result"].(string); ok && v != "" {
		mock.Result = v
		mock.ResultProvided = true
	}
	if eo, ok := obj["errorOutput"].(map[string]interface{}); ok {
		mock.Error, _ = eo["error"].(string)
		mock.Cause, _ = eo["cause"].(string)
		if mock.Error != "" || mock.Cause != "" {
			mock.ErrorProvided = true
		}
	}
	if v, ok := obj["fieldValidationMode"].(string); ok {
		mock.FieldValidationMode = v
	}
	return mock, nil
}

// parseTestStateConfiguration reads the optional TestStateConfiguration
// object from the wire.
func parseTestStateConfiguration(raw interface{}) (*TestStateConfiguration, error) {
	if raw == nil {
		return nil, nil
	}
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return nil, NewValidationException("stateConfiguration must be an object")
	}
	cfg := &TestStateConfiguration{}
	if v, ok := obj["errorCausedByState"].(string); ok {
		cfg.ErrorCausedByState = v
	}
	if v, ok := obj["mapItemReaderData"].(string); ok {
		cfg.MapItemReaderData = v
	}
	if v, ok := obj["mapIterationFailureCount"].(float64); ok {
		cfg.MapIterationFailureCount = int32(v)
	}
	if v, ok := obj["retrierRetryCount"].(float64); ok {
		cfg.RetrierRetryCount = int32(v)
	}
	return cfg, nil
}
