package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateAuthorizer creates a custom authorizer for MQTT connections.
func (s *IoTService) CreateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	status := request.GetParamCaseInsensitive(req.Parameters, "status")
	created, err := s.createAuthorizerCore(store, CreateAuthorizerInput{
		AuthorizerName:         request.GetParamCaseInsensitive(req.Parameters, "authorizerName"),
		AuthorizerFunctionARN:  request.GetParamCaseInsensitive(req.Parameters, "authorizerFunctionArn"),
		TokenKeyName:           request.GetParamCaseInsensitive(req.Parameters, "tokenKeyName"),
		TokenSigningPublicKeys: request.ParseAttributes(req.Parameters, "tokenSigningPublicKeys"),
		SigningDisabled:        request.GetBoolParam(req.Parameters, "signingDisabled"),
		Status:                 status,
		StatusProvided:         status != "",
		EnableCachingForHTTP:   request.GetBoolParam(req.Parameters, "enableCachingForHttp"),
		EnableCachingProvided:  request.HasParam(req.Parameters, "enableCachingForHttp"),
	})
	if err != nil {
		return nil, err
	}

	return authorizerResponse(created), nil
}

// DescribeAuthorizer retrieves details of a custom authorizer.
func (s *IoTService) DescribeAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	auth, err := s.describeAuthorizerCore(store, request.GetParamCaseInsensitive(req.Parameters, "authorizerName"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"authorizerDescription": authorizerResponse(auth)}, nil
}

// UpdateAuthorizer modifies a custom authorizer configuration.
func (s *IoTService) UpdateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// ParseAttributes always returns a non-nil (but possibly empty) map.
	// Guard with len() so absent keys are not overwritten with an empty map.
	var signingKeys map[string]string
	if parsed := request.ParseAttributes(req.Parameters, "tokenSigningPublicKeys"); len(parsed) > 0 {
		signingKeys = parsed
	}

	var enableCaching *bool
	if request.HasParam(req.Parameters, "enableCachingForHttp") {
		v := request.GetBoolParam(req.Parameters, "enableCachingForHttp")
		enableCaching = &v
	}

	existing, err := s.updateAuthorizerCore(store, UpdateAuthorizerInput{
		AuthorizerName:         request.GetParamCaseInsensitive(req.Parameters, "authorizerName"),
		AuthorizerFunctionARN:  request.GetParamCaseInsensitive(req.Parameters, "authorizerFunctionArn"),
		TokenKeyName:           request.GetParamCaseInsensitive(req.Parameters, "tokenKeyName"),
		TokenSigningPublicKeys: signingKeys,
		Status:                 request.GetParamCaseInsensitive(req.Parameters, "status"),
		EnableCaching:          enableCaching,
	})
	if err != nil {
		return nil, err
	}

	return authorizerResponse(existing), nil
}

// DeleteAuthorizer removes a custom authorizer.
func (s *IoTService) DeleteAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteAuthorizerCore(store, request.GetParamCaseInsensitive(req.Parameters, "authorizerName")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListAuthorizers returns all custom authorizers.
func (s *IoTService) ListAuthorizers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := s.listAuthorizersCore(store, opts.Marker, opts.MaxItems)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.Authorizers))
	for _, a := range result.Authorizers {
		items = append(items, authorizerResponse(a))
	}

	return listResponse("authorizers", items, result.NextMarker), nil
}

// SetDefaultAuthorizer designates an existing authorizer as the account's
// default.
func (s *IoTService) SetDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.setDefaultAuthorizerCore(store, request.GetParamCaseInsensitive(req.Parameters, "authorizerName"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"authorizerName": result.AuthorizerName,
		"authorizerArn":  result.AuthorizerARN,
	}, nil
}

// ClearDefaultAuthorizer removes the default-authorizer designation.
func (s *IoTService) ClearDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.clearDefaultAuthorizerCore(store); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// DescribeDefaultAuthorizer retrieves the default authorizer's full
// description.
func (s *IoTService) DescribeDefaultAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	auth, err := s.describeDefaultAuthorizerCore(store)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"authorizerDescription": authorizerResponse(auth),
	}, nil
}
