package lambda

import (
	"fmt"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// FunctionUrlConfigInput carries the wire members of a
// CreateFunctionUrlConfig request.
type FunctionUrlConfigInput struct {
	AuthType   string
	InvokeMode string
	Qualifier  string
	Cors       *lambdastore.CorsConfig
}

// FunctionUrlConfigUpdateInput carries the wire members of an
// UpdateFunctionUrlConfig request. The Cors member keeps its raw wire form
// because the update merges the provided fields into the stored CORS
// configuration.
type FunctionUrlConfigUpdateInput struct {
	AuthType   string
	InvokeMode string
	Cors       map[string]interface{}
	Qualifier  string
}

// createFunctionUrlConfigCore creates a function URL for a function. The
// URL qualifier must name an alias of the function; numeric versions are
// not valid URL qualifiers.
func (s *LambdaService) createFunctionUrlConfigCore(stores *lambdaStore, function *lambdastore.Function, in *FunctionUrlConfigInput) (*lambdastore.FunctionUrlConfig, error) {
	if err := validateAuthType(in.AuthType); err != nil {
		return nil, err
	}

	// AWS returns ResourceConflictException when a URL config already exists.
	if function.UrlConfig != nil {
		return nil, NewResourceConflict(fmt.Sprintf("Function URL already exists for function: %s", function.FunctionName))
	}

	if err := validateInvokeMode(in.InvokeMode); err != nil {
		return nil, err
	}

	qualifier := in.Qualifier
	if qualifier != "" && qualifier != "$LATEST" {
		if _, _, _, err := stores.Functions.ResolveQualifier(function.FunctionName, qualifier); err != nil {
			return nil, NewInvalidParameter("Qualifier", "The function URL qualifier must name an alias of the function")
		}
	}

	config := &lambdastore.FunctionUrlConfig{
		AuthType:   in.AuthType,
		InvokeMode: in.InvokeMode,
		Qualifier:  qualifier,
		Cors:       in.Cors,
	}

	if err := stores.Functions.SetFunctionUrlConfig(function.FunctionName, config); err != nil {
		return nil, mapStoreError(err)
	}

	return config, nil
}

// deleteFunctionUrlConfigCore deletes a function's URL configuration.
func (s *LambdaService) deleteFunctionUrlConfigCore(stores *lambdaStore, functionNameRaw string) error {
	if functionNameRaw == "" {
		return NewInvalidParameter("FunctionName", "Function name is required")
	}
	functionName := extractFunctionName(functionNameRaw)

	if err := stores.Functions.DeleteFunctionUrlConfig(functionName); err != nil {
		return NewResourceNotFound("Function", functionName)
	}
	return nil
}

// updateFunctionUrlConfigCore applies the provided auth type, invoke mode,
// CORS and qualifier updates to a function's URL configuration.
func (s *LambdaService) updateFunctionUrlConfigCore(stores *lambdaStore, function *lambdastore.Function, in *FunctionUrlConfigUpdateInput) (*lambdastore.FunctionUrlConfig, error) {
	if in.AuthType != "" {
		if err := validateAuthType(in.AuthType); err != nil {
			return nil, err
		}
		function.UrlConfig.AuthType = in.AuthType
	}
	if in.InvokeMode != "" {
		if err := validateInvokeMode(in.InvokeMode); err != nil {
			return nil, err
		}
		function.UrlConfig.InvokeMode = in.InvokeMode
	}
	if in.Cors != nil {
		if function.UrlConfig.Cors == nil {
			function.UrlConfig.Cors = &lambdastore.CorsConfig{}
		}
		function.UrlConfig.Cors = updateCorsConfig(function.UrlConfig.Cors, in.Cors)
	}

	// The URL qualifier names an alias of the function; numeric versions
	// are rejected like on create.
	if in.Qualifier != "" {
		if in.Qualifier != "$LATEST" {
			if _, _, _, err := stores.Functions.ResolveQualifier(function.FunctionName, in.Qualifier); err != nil {
				return nil, NewInvalidParameter("Qualifier", "The function URL qualifier must name an alias of the function")
			}
		}
		function.UrlConfig.Qualifier = in.Qualifier
	}

	if err := stores.Functions.SetFunctionUrlConfig(function.FunctionName, function.UrlConfig); err != nil {
		return nil, mapStoreError(err)
	}

	updatedFunction, err := stores.Functions.Get(function.FunctionName)
	if err != nil {
		return nil, err
	}

	return updatedFunction.UrlConfig, nil
}
