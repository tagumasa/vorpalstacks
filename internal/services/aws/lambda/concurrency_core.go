package lambda

import (
	"errors"

	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// ConcurrencyInput carries the fields needed by the reserved-concurrency
// operations. FunctionName arrives in its raw wire form; the Core
// resolves every reference shape the API accepts.
type ConcurrencyInput struct {
	FunctionName string
	Reserved     int64
}

// resolveConcurrencyFunctionName resolves the FunctionName reference forms
// for the function-level concurrency operations (name, name:qualifier,
// full or partial ARN — the qualifier is irrelevant because reserved
// concurrency applies to the whole function).
func resolveConcurrencyFunctionName(functionNameRaw string) (string, error) {
	functionName := extractFunctionName(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return "", err
	}
	return functionName, nil
}

// putFunctionConcurrencyCore sets the reserved concurrent execution limit
// for a function. Member validation precedes the store acquisition so an
// invalid request never surfaces a storage error.
func (s *LambdaService) putFunctionConcurrencyCore(reqCtx *request.RequestContext, in *ConcurrencyInput) (int64, error) {
	functionName, err := resolveConcurrencyFunctionName(in.FunctionName)
	if err != nil {
		return 0, err
	}
	if in.Reserved < 0 {
		return 0, NewInvalidParameter("ReservedConcurrentExecutions", "Must be non-negative. Use DeleteFunctionConcurrency to remove concurrency limits.")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return 0, err
	}
	if err := stores.Functions.SetReservedConcurrency(functionName, &in.Reserved); err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return 0, ErrResourceNotFound
		}
		return 0, err
	}

	return in.Reserved, nil
}

// getFunctionConcurrencyCore retrieves the reserved concurrent execution
// limit for a function.
func (s *LambdaService) getFunctionConcurrencyCore(reqCtx *request.RequestContext, functionNameRaw string) (int64, error) {
	functionName, err := resolveConcurrencyFunctionName(functionNameRaw)
	if err != nil {
		return 0, err
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return 0, err
	}
	concurrency, err := stores.Functions.GetReservedConcurrency(functionName)
	if err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return 0, ErrResourceNotFound
		}
		return 0, err
	}
	if concurrency == nil {
		// AWS answers with ResourceNotFoundException when the function has
		// never had reserved concurrency configured: the concurrency
		// sub-resource does not exist until PutFunctionConcurrency sets it.
		return 0, ErrResourceNotFound
	}

	return *concurrency, nil
}

// deleteFunctionConcurrencyCore removes the reserved concurrent execution
// limit from a function.
func (s *LambdaService) deleteFunctionConcurrencyCore(reqCtx *request.RequestContext, functionNameRaw string) error {
	functionName, err := resolveConcurrencyFunctionName(functionNameRaw)
	if err != nil {
		return err
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	if err := stores.Functions.SetReservedConcurrency(functionName, nil); err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return ErrResourceNotFound
		}
		return err
	}
	return nil
}
