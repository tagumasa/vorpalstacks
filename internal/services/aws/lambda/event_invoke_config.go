package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// resolveEventInvokeTarget resolves the FunctionName reference forms and
// merges the Qualifier parameter with an embedded qualifier, defaulting to
// $LATEST when neither is present. Shared by the event-invoke-config
// operations.
func resolveEventInvokeTarget(params map[string]interface{}) (string, string, error) {
	functionNameRaw := request.GetStringParam(params, "FunctionName")
	if functionNameRaw == "" {
		return "", "", NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName, embeddedQualifier := resolveFunctionRef(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return "", "", err
	}
	qualifier := mergeQualifier(request.GetStringParam(params, "Qualifier"), embeddedQualifier)
	if qualifier == "" {
		qualifier = "$LATEST"
	}
	return functionName, qualifier, nil
}

// eventInvokeConfigInput builds the Core input from the wire request.
func eventInvokeConfigInput(req *request.ParsedRequest, functionName, qualifier string) *EventInvokeConfigInput {
	in := &EventInvokeConfigInput{
		FunctionName: functionName,
		Qualifier:    qualifier,
	}
	if _, ok := req.Parameters["MaximumEventAgeInSeconds"]; ok {
		in.HasMaximumEventAgeInSeconds = true
		in.MaximumEventAgeInSeconds = int32(request.GetIntParam(req.Parameters, "MaximumEventAgeInSeconds"))
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		in.HasMaximumRetryAttempts = true
		in.MaximumRetryAttempts = int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
	}
	if destMap := request.GetMapParam(req.Parameters, "DestinationConfig"); destMap != nil {
		in.DestinationConfig = parseDestinationConfig(destMap)
	}
	return in
}

// PutFunctionEventInvokeConfig creates or updates the configuration for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) PutFunctionEventInvokeConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, qualifier, err := resolveEventInvokeTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	config, err := s.putFunctionEventInvokeConfigCore(reqCtx, eventInvokeConfigInput(req, functionName, qualifier))
	if err != nil {
		return nil, err
	}

	return s.toEventInvokeConfig(config), nil
}

// GetFunctionEventInvokeConfig retrieves the configuration for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) GetFunctionEventInvokeConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, qualifier, err := resolveEventInvokeTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := s.getFunctionEventInvokeConfigCore(store, functionName, qualifier)
	if err != nil {
		return nil, err
	}

	return s.toEventInvokeConfig(config), nil
}

// DeleteFunctionEventInvokeConfig deletes the configuration for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) DeleteFunctionEventInvokeConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, qualifier, err := resolveEventInvokeTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteFunctionEventInvokeConfigCore(store, functionName, qualifier); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListFunctionEventInvokeConfigs lists all configurations for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) ListFunctionEventInvokeConfigs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionNameRaw := request.GetStringParam(req.Parameters, "FunctionName")
	if functionNameRaw == "" {
		return nil, NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName := extractFunctionName(functionNameRaw)
	if err := validateFunctionName(functionName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	configs, err := s.listFunctionEventInvokeConfigsCore(store, functionName)
	if err != nil {
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
	functionName, qualifier, err := resolveEventInvokeTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := s.updateFunctionEventInvokeConfigCore(store, eventInvokeConfigInput(req, functionName, qualifier))
	if err != nil {
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
