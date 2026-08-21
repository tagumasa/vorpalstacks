package lambda

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// resolveConcurrencyFunction resolves the FunctionName reference forms for
// the function-level concurrency operations (name, name:qualifier, full or
// partial ARN — the qualifier is irrelevant because reserved concurrency
// applies to the whole function).
func resolveConcurrencyFunction(params map[string]interface{}) (string, error) {
	functionNameRaw := request.GetStringParam(params, "FunctionName")
	if functionNameRaw == "" {
		return "", NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName := extractFunctionName(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return "", err
	}
	return functionName, nil
}

// PutFunctionConcurrency sets the reserved concurrent execution limit for the specified Lambda function.
func (s *LambdaService) PutFunctionConcurrency(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, err := resolveConcurrencyFunction(req.Parameters)
	if err != nil {
		return nil, err
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
	functionName, err := resolveConcurrencyFunction(req.Parameters)
	if err != nil {
		return nil, err
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
		// AWS answers with ResourceNotFoundException when the function has
		// never had reserved concurrency configured: the concurrency
		// sub-resource does not exist until PutFunctionConcurrency sets it.
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"ReservedConcurrentExecutions": *concurrency,
	}, nil
}

// DeleteFunctionConcurrency removes the reserved concurrent execution limit from the specified Lambda function.
func (s *LambdaService) DeleteFunctionConcurrency(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, err := resolveConcurrencyFunction(req.Parameters)
	if err != nil {
		return nil, err
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
