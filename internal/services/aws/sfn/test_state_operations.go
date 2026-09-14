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
