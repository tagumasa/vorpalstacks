package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// PutFunctionConcurrency sets the reserved concurrent execution limit for the specified Lambda function.
func (s *LambdaService) PutFunctionConcurrency(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	concurrency, err := s.putFunctionConcurrencyCore(reqCtx, &ConcurrencyInput{
		FunctionName: request.GetStringParam(req.Parameters, "FunctionName"),
		Reserved:     int64(request.GetIntParam(req.Parameters, "ReservedConcurrentExecutions")),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ReservedConcurrentExecutions": concurrency,
	}, nil
}

// GetFunctionConcurrency retrieves the reserved concurrent execution limit for the specified Lambda function.
func (s *LambdaService) GetFunctionConcurrency(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	concurrency, err := s.getFunctionConcurrencyCore(reqCtx, request.GetStringParam(req.Parameters, "FunctionName"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ReservedConcurrentExecutions": concurrency,
	}, nil
}

// DeleteFunctionConcurrency removes the reserved concurrent execution limit from the specified Lambda function.
func (s *LambdaService) DeleteFunctionConcurrency(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteFunctionConcurrencyCore(reqCtx, request.GetStringParam(req.Parameters, "FunctionName")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
