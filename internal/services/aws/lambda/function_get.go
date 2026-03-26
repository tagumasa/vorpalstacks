// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	"vorpalstacks/internal/store/aws/common"
)

// GetFunction retrieves information about the specified Lambda function.
func (s *LambdaService) GetFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, version, _, err := s.validateAndGetFunctionWithQualifier(reqCtx, req.Parameters)
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
			"RepositoryType": "S3",
			"ImageUri":       function.ImageUri,
		},
		"Tags": tagutil.ToMap(function.Tags),
	}, nil
}

// GetFunctionConfiguration retrieves the configuration of the specified Lambda function.
func (s *LambdaService) GetFunctionConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, version, _, err := s.validateAndGetFunctionWithQualifier(reqCtx, req.Parameters)
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
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 {
		maxItems = 50
	}
	marker := request.GetStringParam(req.Parameters, "Marker")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Functions.List(common.ListOptions{
		Prefix:   "",
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	functions := make([]interface{}, 0, len(result.Items))
	for _, fn := range result.Items {
		functions = append(functions, s.toFunctionConfiguration(fn))
	}

	response := map[string]interface{}{
		"Functions": functions,
	}

	if result.IsTruncated {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}
