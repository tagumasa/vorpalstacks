// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// GetFunction retrieves information about the specified Lambda function.
func (s *LambdaService) GetFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	function, version, _, tags, err := s.getFunctionCore(store, &GetFunctionInput{
		FunctionName: request.GetStringParam(req.Parameters, "FunctionName"),
		Qualifier:    request.GetStringParam(req.Parameters, "Qualifier"),
	})
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if version != nil {
		config = s.toVersionConfiguration(version)
	} else {
		config = s.toFunctionConfiguration(function)
	}

	return map[string]interface{}{
		"Configuration": config,
		"Code": map[string]interface{}{
			"Location":       function.CodeLocation,
			"RepositoryType": repositoryType(function),
			"ImageUri":       function.ImageUri,
		},
		"Tags": tags,
	}, nil
}

// GetFunctionConfiguration retrieves the configuration of the specified Lambda function.
func (s *LambdaService) GetFunctionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	function, version, _, err := s.getFunctionConfigurationCore(store, &GetFunctionInput{
		FunctionName: request.GetStringParam(req.Parameters, "FunctionName"),
		Qualifier:    request.GetStringParam(req.Parameters, "Qualifier"),
	})
	if err != nil {
		return nil, err
	}

	if version != nil {
		return s.toVersionConfiguration(version), nil
	}

	return s.toFunctionConfiguration(function), nil
}

// ListFunctions lists all Lambda functions in the current account.
func (s *LambdaService) ListFunctions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.listFunctionsCore(store, &ListFunctionsInput{
		Marker:   request.GetStringParam(req.Parameters, "Marker"),
		MaxItems: request.GetIntParam(req.Parameters, "MaxItems"),
	})
	if err != nil {
		return nil, err
	}

	functions := make([]interface{}, 0, len(items))
	for _, fn := range items {
		functions = append(functions, s.toFunctionConfiguration(fn))
	}

	response := map[string]interface{}{
		"Functions": functions,
	}

	if nextMarker != "" {
		response["NextMarker"] = nextMarker
	}

	return response, nil
}
