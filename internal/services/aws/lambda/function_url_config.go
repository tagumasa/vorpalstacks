package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateFunctionUrlConfig creates a function URL for a Lambda function.
// A function URL is a dedicated HTTP(S) endpoint that can be used to invoke
// the function without authentication or with AWS IAM or AWS_IAM auth.
func (s *LambdaService) CreateFunctionUrlConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	// The URL qualifier names an alias of the function ("The alias name" —
	// API model); numeric versions are not valid URL qualifiers. An
	// embedded ":qualifier" suffix on the FunctionName reference is
	// equivalent to the Qualifier parameter.
	_, embeddedQualifier := resolveFunctionRef(request.GetStringParam(req.Parameters, "FunctionName"))
	qualifier := mergeQualifier(request.GetStringParam(req.Parameters, "Qualifier"), embeddedQualifier)

	in := &FunctionUrlConfigInput{
		AuthType:   request.GetStringParam(req.Parameters, "AuthType"),
		InvokeMode: request.GetStringParam(req.Parameters, "InvokeMode"),
		Qualifier:  qualifier,
	}
	if corsMap := request.GetMapParam(req.Parameters, "Cors"); corsMap != nil {
		in.Cors = parseCorsConfig(corsMap)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := s.createFunctionUrlConfigCore(store, function, in)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"FunctionUrl":      config.FunctionUrl,
		"FunctionArn":      config.FunctionArn,
		"AuthType":         config.AuthType,
		"CreationTime":     config.CreationTime.Format(timeutils.ISO8601UTCFormat),
		"LastModifiedTime": config.LastModifiedTime.Format(timeutils.ISO8601UTCFormat),
	}

	if config.InvokeMode != "" {
		result["InvokeMode"] = config.InvokeMode
	}
	if config.Cors != nil {
		result["Cors"] = toCorsConfig(config.Cors)
	}

	return result, nil
}

// DeleteFunctionUrlConfig deletes the function URL configuration for a Lambda function.
// This removes the dedicated HTTP(S) endpoint associated with the function.
func (s *LambdaService) DeleteFunctionUrlConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := request.GetStringParam(req.Parameters, "FunctionName")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteFunctionUrlConfigCore(store, functionName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetFunctionUrlConfig retrieves the function URL configuration for a Lambda function.
// Returns the URL, auth type, CORS settings, and other configuration details.
func (s *LambdaService) GetFunctionUrlConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if function.UrlConfig == nil {
		return nil, ErrResourceNotFound
	}

	return s.toFunctionUrlConfig(function.UrlConfig), nil
}

// UpdateFunctionUrlConfig updates the function URL configuration for a Lambda function.
// Allows modification of auth type, invoke mode, and CORS settings.
func (s *LambdaService) UpdateFunctionUrlConfig(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if function.UrlConfig == nil {
		return nil, ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := s.updateFunctionUrlConfigCore(store, function, &FunctionUrlConfigUpdateInput{
		AuthType:   request.GetStringParam(req.Parameters, "AuthType"),
		InvokeMode: request.GetStringParam(req.Parameters, "InvokeMode"),
		Cors:       request.GetMapParam(req.Parameters, "Cors"),
		Qualifier:  request.GetStringParam(req.Parameters, "Qualifier"),
	})
	if err != nil {
		return nil, err
	}

	return s.toFunctionUrlConfig(config), nil
}

// ListFunctionUrlConfigs lists the function URL configurations for a Lambda function.
// Since a function can only have one URL config, this returns an array with at most one element.
func (s *LambdaService) ListFunctionUrlConfigs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	function, err := s.validateAndGetFunction(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if function.UrlConfig == nil {
		return map[string]interface{}{
			"FunctionUrlConfigs": []interface{}{},
		}, nil
	}

	return map[string]interface{}{
		"FunctionUrlConfigs": []interface{}{s.toFunctionUrlConfig(function.UrlConfig)},
	}, nil
}

func (s *LambdaService) toFunctionUrlConfig(c *lambdastore.FunctionUrlConfig) map[string]interface{} {
	if c == nil {
		return response.EmptyResponse()
	}

	result := map[string]interface{}{
		"FunctionUrl": c.FunctionUrl,
		"FunctionArn": c.FunctionArn,
		"AuthType":    c.AuthType,
		"InvokeMode":  c.InvokeMode,
	}

	if !c.CreationTime.IsZero() {
		result["CreationTime"] = c.CreationTime.Format(timeutils.ISO8601UTCFormat)
	}
	if !c.LastModifiedTime.IsZero() {
		result["LastModifiedTime"] = c.LastModifiedTime.Format(timeutils.ISO8601UTCFormat)
	}

	if c.Cors != nil {
		result["Cors"] = toCorsConfig(c.Cors)
	}

	return result
}

func applyCorsFields(config *lambdastore.CorsConfig, corsMap map[string]interface{}) {
	if allowCreds, ok := corsMap["AllowCredentials"].(bool); ok {
		config.AllowCredentials = allowCreds
	}
	if allowHeaders, ok := corsMap["AllowHeaders"].([]interface{}); ok {
		config.AllowHeaders = make([]string, 0, len(allowHeaders))
		for _, h := range allowHeaders {
			if hs, ok := h.(string); ok {
				config.AllowHeaders = append(config.AllowHeaders, hs)
			}
		}
	}
	if allowMethods, ok := corsMap["AllowMethods"].([]interface{}); ok {
		config.AllowMethods = make([]string, 0, len(allowMethods))
		for _, m := range allowMethods {
			if ms, ok := m.(string); ok {
				config.AllowMethods = append(config.AllowMethods, ms)
			}
		}
	}
	if allowOrigins, ok := corsMap["AllowOrigins"].([]interface{}); ok {
		config.AllowOrigins = make([]string, 0, len(allowOrigins))
		for _, o := range allowOrigins {
			if os, ok := o.(string); ok {
				config.AllowOrigins = append(config.AllowOrigins, os)
			}
		}
	}
	if exposeHeaders, ok := corsMap["ExposeHeaders"].([]interface{}); ok {
		config.ExposeHeaders = make([]string, 0, len(exposeHeaders))
		for _, h := range exposeHeaders {
			if hs, ok := h.(string); ok {
				config.ExposeHeaders = append(config.ExposeHeaders, hs)
			}
		}
	}
	if maxAge, ok := corsMap["MaxAge"].(int); ok {
		config.MaxAge = int32(maxAge)
	}
	if maxAgeFloat, ok := corsMap["MaxAge"].(float64); ok {
		config.MaxAge = int32(maxAgeFloat)
	}
}

func parseCorsConfig(corsMap map[string]interface{}) *lambdastore.CorsConfig {
	config := &lambdastore.CorsConfig{}
	applyCorsFields(config, corsMap)
	return config
}

func toCorsConfig(c *lambdastore.CorsConfig) map[string]interface{} {
	return map[string]interface{}{
		"AllowCredentials": c.AllowCredentials,
		"AllowHeaders":     c.AllowHeaders,
		"AllowMethods":     c.AllowMethods,
		"AllowOrigins":     c.AllowOrigins,
		"ExposeHeaders":    c.ExposeHeaders,
		"MaxAge":           c.MaxAge,
	}
}

func updateCorsConfig(existing *lambdastore.CorsConfig, corsMap map[string]interface{}) *lambdastore.CorsConfig {
	applyCorsFields(existing, corsMap)
	return existing
}
