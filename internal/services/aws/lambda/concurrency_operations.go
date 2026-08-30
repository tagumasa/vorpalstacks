package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
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

	concurrency, err := s.putFunctionConcurrencyCore(reqCtx, &ConcurrencyInput{
		FunctionName: functionName,
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
	functionName, err := resolveConcurrencyFunction(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	concurrency, err := s.getFunctionConcurrencyCore(store, functionName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ReservedConcurrentExecutions": concurrency,
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

	if err := s.deleteFunctionConcurrencyCore(store, functionName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
