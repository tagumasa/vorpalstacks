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

// PutFunctionEventInvokeConfig creates or updates the configuration for asynchronous invocation of the specified Lambda function.
func (s *LambdaService) PutFunctionEventInvokeConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName, qualifier, err := resolveEventInvokeTarget(req.Parameters)
	if err != nil {
		return nil, err
	}

	config := &lambdastore.EventInvokeConfig{}

	if _, ok := req.Parameters["MaximumEventAgeInSeconds"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "MaximumEventAgeInSeconds"))
		if err := validateMaximumEventAgeInSeconds(val); err != nil {
			return nil, err
		}
		config.MaximumEventAgeInSeconds = val
	} else {
		// Put replaces the whole configuration; omitted members fall back
		// to the AWS defaults rather than zero.
		config.MaximumEventAgeInSeconds = lambdastore.DefaultMaximumEventAgeInSeconds
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
		if err := validateMaximumRetryAttempts(val); err != nil {
			return nil, err
		}
		config.MaximumRetryAttempts = val
	} else {
		config.MaximumRetryAttempts = lambdastore.DefaultMaximumRetryAttempts
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
	functionName, qualifier, err := resolveEventInvokeTarget(req.Parameters)
	if err != nil {
		return nil, err
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
	functionName, qualifier, err := resolveEventInvokeTarget(req.Parameters)
	if err != nil {
		return nil, err
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
	functionName, qualifier, err := resolveEventInvokeTarget(req.Parameters)
	if err != nil {
		return nil, err
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
		// Creating the configuration through an update: fields the update
		// does not specify take the AWS defaults.
		config = &lambdastore.EventInvokeConfig{
			MaximumEventAgeInSeconds: lambdastore.DefaultMaximumEventAgeInSeconds,
			MaximumRetryAttempts:     lambdastore.DefaultMaximumRetryAttempts,
		}
	}

	// Overwrite only the fields that were explicitly provided in the request.
	if _, ok := req.Parameters["MaximumEventAgeInSeconds"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "MaximumEventAgeInSeconds"))
		if err := validateMaximumEventAgeInSeconds(val); err != nil {
			return nil, err
		}
		config.MaximumEventAgeInSeconds = val
	}
	if _, ok := req.Parameters["MaximumRetryAttempts"]; ok {
		val := int32(request.GetIntParam(req.Parameters, "MaximumRetryAttempts"))
		if err := validateMaximumRetryAttempts(val); err != nil {
			return nil, err
		}
		config.MaximumRetryAttempts = val
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
