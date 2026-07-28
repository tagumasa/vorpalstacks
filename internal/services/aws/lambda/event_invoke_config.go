package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// PutFunctionEventInvokeConfig creates or updates the configuration for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) PutFunctionEventInvokeConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if qualifier == "" {
		qualifier = "$LATEST"
	}

	config := &lambdastore.EventInvokeConfig{}

	if _, ok := req.Parameters["MaximumEventAgeInSeconds"]; ok {
		maxEventAge := request.GetIntParam(req.Parameters, "MaximumEventAgeInSeconds")
		if maxEventAge < 60 || maxEventAge > 21600 {
			return nil, NewInvalidParameter("MaximumEventAgeInSeconds",
				"MaximumEventAgeInSeconds must be between 60 and 21600 seconds")
		}
		config.MaximumEventAgeInSeconds = int32(maxEventAge)
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		maxRetry := request.GetIntParam(req.Parameters, "MaximumRetryAttempts")
		if maxRetry < 0 || maxRetry > 2 {
			return nil, NewInvalidParameter("MaximumRetryAttempts",
				"MaximumRetryAttempts must be between 0 and 2")
		}
		config.MaximumRetryAttempts = int32(maxRetry)
	}
	if destMap := request.GetMapParam(req.Parameters, "DestinationConfig"); destMap != nil {
		config.DestinationConfig = parseDestinationConfig(destMap)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.Functions.SetEventInvokeConfig(functionName, qualifier, config); err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return s.toEventInvokeConfig(config), nil
}

// GetFunctionEventInvokeConfig retrieves the configuration for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) GetFunctionEventInvokeConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if qualifier == "" {
		qualifier = "$LATEST"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := store.Functions.GetEventInvokeConfig(functionName, qualifier)
	if err != nil {
		if err == lambdastore.ErrEventInvokeConfigNotFound || err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return s.toEventInvokeConfig(config), nil
}

// DeleteFunctionEventInvokeConfig deletes the configuration for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) DeleteFunctionEventInvokeConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if qualifier == "" {
		qualifier = "$LATEST"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.Functions.DeleteEventInvokeConfig(functionName, qualifier); err != nil {
		if err == lambdastore.ErrEventInvokeConfigNotFound || err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListFunctionEventInvokeConfigs lists all configurations for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) ListFunctionEventInvokeConfigs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configs, err := store.Functions.ListEventInvokeConfigs(functionName)
	if err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(configs))
	for _, c := range configs {
		items = append(items, s.toEventInvokeConfig(&c))
	}

	return map[string]interface{}{
		"FunctionEventInvokeConfigs": items,
	}, nil
}

// UpdateFunctionEventInvokeConfig updates the configuration for asynchronous invocation
// of the specified Lambda function. Only fields provided in the request are modified;
// existing values for unprovided fields are preserved. DestinationConfig is replaced
// atomically when provided. If no config exists for the qualifier, a new one is created.
func (s *LambdaService) UpdateFunctionEventInvokeConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	if functionName == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}

	qualifier := request.GetStringParam(req.Parameters, "Qualifier")
	if qualifier == "" {
		qualifier = "$LATEST"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Attempt to load existing config; if not found, create a blank one.
	config, err := store.Functions.GetEventInvokeConfig(functionName, qualifier)
	if err != nil && err != lambdastore.ErrEventInvokeConfigNotFound && err != lambdastore.ErrFunctionNotFound {
		return nil, err
	}
	if config == nil {
		config = &lambdastore.EventInvokeConfig{}
	}

	// Overwrite only the fields that were explicitly provided in the request.
	if _, ok := req.Parameters["MaximumEventAgeInSeconds"]; ok {
		config.MaximumEventAgeInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumEventAgeInSeconds"))
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		config.MaximumRetryAttempts = int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
	}
	if destMap := request.GetMapParam(req.Parameters, "DestinationConfig"); destMap != nil {
		config.DestinationConfig = parseDestinationConfig(destMap)
	}

	if err := store.Functions.SetEventInvokeConfig(functionName, qualifier, config); err != nil {
		if err == lambdastore.ErrFunctionNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return s.toEventInvokeConfig(config), nil
}

func (s *LambdaService) toEventInvokeConfig(c *lambdastore.EventInvokeConfig) map[string]interface{} {
	result := map[string]interface{}{
		"FunctionName": c.FunctionName,
		"Qualifier":    c.Qualifier,
		"LastModified": float64(c.LastModified.Unix()),
	}

	if c.MaximumEventAgeInSeconds > 0 {
		result["MaximumEventAgeInSeconds"] = c.MaximumEventAgeInSeconds
	}
	if c.MaximumRetryAttempts >= 0 {
		result["MaximumRetryAttempts"] = c.MaximumRetryAttempts
	}
	if c.DestinationConfig != nil {
		result["DestinationConfig"] = toDestinationConfig(c.DestinationConfig)
	}

	return result
}
