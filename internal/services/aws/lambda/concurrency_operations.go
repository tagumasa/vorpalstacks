package lambda

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// PutFunctionConcurrency sets the reserved concurrent execution limit for the specified Lambda function.
func (s *LambdaService) PutFunctionConcurrency(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	concurrency := int64(request.GetIntParam(req.Parameters, "ReservedConcurrentExecutions"))
	if concurrency < 0 {
		return nil, NewInvalidParameter("ReservedConcurrentExecutions", "Must be non-negative. Use DeleteFunctionConcurrency to remove concurrency limits.")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	concurrencyPtr := &concurrency
	if err := store.Functions.SetReservedConcurrency(functionName, concurrencyPtr); err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"ReservedConcurrentExecutions": concurrency,
	}, nil
}

// GetFunctionConcurrency retrieves the reserved concurrent execution limit for the specified Lambda function.
func (s *LambdaService) GetFunctionConcurrency(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	concurrency, err := store.Functions.GetReservedConcurrency(functionName)
	if err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	if concurrency == nil {
		return response.EmptyResponse(), nil
	}

	return map[string]interface{}{
		"ReservedConcurrentExecutions": *concurrency,
	}, nil
}

// DeleteFunctionConcurrency removes the reserved concurrent execution limit from the specified Lambda function.
func (s *LambdaService) DeleteFunctionConcurrency(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Functions.SetReservedConcurrency(functionName, nil); err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}
