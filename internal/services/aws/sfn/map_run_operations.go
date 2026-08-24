package sfn

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// DescribeMapRun retrieves the details of a Step Functions map run.
func (s *StepFunctionService) DescribeMapRun(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.describeMapRunCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "mapRunArn"))
}

// ListMapRuns lists map runs, optionally filtered by execution ARN.
func (s *StepFunctionService) ListMapRuns(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := parsePageLimit(req)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listMapRunsCore(ctx, store,
		request.GetParamLowerFirst(req.Parameters, "executionArn"),
		limit,
		request.GetParamLowerFirst(req.Parameters, "nextToken"))
}

// UpdateMapRun modifies the concurrency and failure tolerance settings of a
// running map execution.
func (s *StepFunctionService) UpdateMapRun(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := UpdateMapRunInput{
		MapRunArn: request.GetParamLowerFirst(req.Parameters, "mapRunArn"),
	}
	if _, ok := req.Parameters["maxConcurrency"]; ok {
		v := int64(request.GetIntParam(req.Parameters, "maxConcurrency"))
		in.MaxConcurrency = &v
	}
	if _, ok := req.Parameters["toleratedFailureCount"]; ok {
		v := request.GetInt64Param(req.Parameters, "toleratedFailureCount")
		in.ToleratedFailureCount = &v
	}
	if _, ok := req.Parameters["toleratedFailurePercentage"]; ok {
		// The Smithy model types toleratedFailurePercentage as a
		// single-precision float; the cast marks the mandated precision
		// boundary between generic JSON numbers and the wire type.
		v := float32(request.GetFloatParam(req.Parameters, "toleratedFailurePercentage"))
		in.ToleratedFailurePercentage = &v
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.updateMapRunCore(ctx, store, in); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
