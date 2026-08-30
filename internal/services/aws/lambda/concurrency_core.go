package lambda

import (
	"errors"

	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// ConcurrencyInput carries the fields needed by the reserved-concurrency
// operations. The function name arrives already resolved from its wire
// reference forms by the handler.
type ConcurrencyInput struct {
	FunctionName string
	Reserved     int64
}

// putFunctionConcurrencyCore sets the reserved concurrent execution limit
// for a function.
func (s *LambdaService) putFunctionConcurrencyCore(reqCtx *request.RequestContext, in *ConcurrencyInput) (int64, error) {
	if in.Reserved < 0 {
		return 0, NewInvalidParameter("ReservedConcurrentExecutions", "Must be non-negative. Use DeleteFunctionConcurrency to remove concurrency limits.")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return 0, err
	}
	if err := stores.Functions.SetReservedConcurrency(in.FunctionName, &in.Reserved); err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return 0, ErrResourceNotFound
		}
		return 0, err
	}

	return in.Reserved, nil
}

// getFunctionConcurrencyCore retrieves the reserved concurrent execution
// limit for a function.
func (s *LambdaService) getFunctionConcurrencyCore(stores *lambdaStore, functionName string) (int64, error) {
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
func (s *LambdaService) deleteFunctionConcurrencyCore(stores *lambdaStore, functionName string) error {
	if err := stores.Functions.SetReservedConcurrency(functionName, nil); err != nil {
		if errors.Is(err, lambdastore.ErrFunctionNotFound) {
			return ErrResourceNotFound
		}
		return err
	}
	return nil
}
