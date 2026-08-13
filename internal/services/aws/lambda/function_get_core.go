package lambda

import (
	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// GetFunctionInput carries the fields needed for GetFunction.
type GetFunctionInput struct {
	FunctionName string
	Qualifier    string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// getFunctionCore retrieves a function (optionally by qualifier) along with
// its tags. It is the single entry point shared by the HTTP API handler and
// the admin gRPC handler. The raw function name or ARN is resolved and
// validated internally so that all callers share a single validation path.
func (s *LambdaService) getFunctionCore(stores *lambdaStore, in *GetFunctionInput) (*lambdastore.Function, *lambdastore.Version, *lambdastore.Alias, map[string]string, error) {
	functionName := extractFunctionName(in.FunctionName)
	if err := validateFunctionName(functionName); err != nil {
		return nil, nil, nil, nil, err
	}

	function, version, alias, err := s.resolveQualifier(stores.Functions, functionName, in.Qualifier)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tags, err := stores.Functions.TagStore.List(function.FunctionName)
	if err != nil {
		logs.Warn("Failed to fetch tags for function",
			logs.String("function", function.FunctionName),
			logs.Err(err))
		tags = map[string]string{}
	}

	return function, version, alias, tags, nil
}

// getFunctionConfigurationCore retrieves a function configuration (optionally
// by qualifier). It is the single entry point shared by the HTTP API handler
// and the admin gRPC handler.
func (s *LambdaService) getFunctionConfigurationCore(stores *lambdaStore, in *GetFunctionInput) (*lambdastore.Function, *lambdastore.Version, *lambdastore.Alias, error) {
	functionName := extractFunctionName(in.FunctionName)
	if err := validateFunctionName(functionName); err != nil {
		return nil, nil, nil, err
	}

	function, version, alias, err := s.resolveQualifier(stores.Functions, functionName, in.Qualifier)
	if err != nil {
		return nil, nil, nil, err
	}
	return function, version, alias, nil
}
