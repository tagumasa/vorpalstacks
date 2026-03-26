package lambda

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/timeutils"
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

	if maxEventAge := request.GetIntParam(req.Parameters, "MaximumEventAgeInSeconds"); maxEventAge > 0 {
		config.MaximumEventAgeInSeconds = int32(maxEventAge)
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		config.MaximumRetryAttempts = int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
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

func (s *LambdaService) toEventInvokeConfig(c *lambdastore.EventInvokeConfig) map[string]interface{} {
	result := map[string]interface{}{
		"FunctionName": c.FunctionName,
		"Qualifier":    c.Qualifier,
		"LastModified": c.LastModified.Format(timeutils.ISO8601UTCFormat),
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
