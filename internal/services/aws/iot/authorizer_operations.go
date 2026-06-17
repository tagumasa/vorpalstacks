package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// CreateAuthorizer creates a custom authorizer for MQTT connections.
func (s *IoTService) CreateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	auth := &iotstore.Authorizer{
		AuthorizerName:        name,
		AuthorizerFunctionARN: request.GetParamCaseInsensitive(req.Parameters, "authorizerFunctionArn"),
		TokenName:             request.GetParamCaseInsensitive(req.Parameters, "tokenKeyName"),
		TokenSignature:        request.GetParamCaseInsensitive(req.Parameters, "tokenSignature"),
		Status:                true,
		EnableCachingForHTTP:  true,
		CreationDate:          time.Now().UTC(),
		LastModifiedDate:      time.Now().UTC(),
	}

	if statusStr := request.GetParamCaseInsensitive(req.Parameters, "status"); statusStr != "" {
		if err := ValidateAuthorizerStatus(statusStr); err != nil {
			return nil, err
		}
		auth.Status = statusStr == "ACTIVE"
	}

	if request.HasParam(req.Parameters, "enableCachingForHttp") {
		auth.EnableCachingForHTTP = request.GetBoolParam(req.Parameters, "enableCachingForHttp")
	}

	created, err := store.CreateAuthorizer(auth)
	if err != nil {
		return nil, err
	}

	return authorizerResponse(created), nil
}

// DescribeAuthorizer retrieves details of a custom authorizer.
func (s *IoTService) DescribeAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	auth, err := store.GetAuthorizer(name)
	if err != nil {
		return nil, iotstore.ErrAuthorizerNotFound
	}

	return map[string]interface{}{"authorizerDescription": authorizerResponse(auth)}, nil
}

// UpdateAuthorizer modifies a custom authorizer configuration.
func (s *IoTService) UpdateAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := iotstore.AuthorizerUpdateOpts{
		FunctionARN:    request.GetParamCaseInsensitive(req.Parameters, "authorizerFunctionArn"),
		TokenName:      request.GetParamCaseInsensitive(req.Parameters, "tokenKeyName"),
		TokenSignature: request.GetParamCaseInsensitive(req.Parameters, "tokenSignature"),
		Status:         request.GetParamCaseInsensitive(req.Parameters, "status"),
	}
	if opts.Status != "" {
		if err := ValidateAuthorizerStatus(opts.Status); err != nil {
			return nil, err
		}
	}
	if request.HasParam(req.Parameters, "enableCachingForHttp") {
		v := request.GetBoolParam(req.Parameters, "enableCachingForHttp")
		opts.EnableCaching = &v
	}

	existing, err := store.UpdateAuthorizer(name, opts)
	if err != nil {
		return nil, err
	}

	return authorizerResponse(existing), nil
}

// DeleteAuthorizer removes a custom authorizer.
func (s *IoTService) DeleteAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "authorizerName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteAuthorizer(name); err != nil {
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

	auths, err := store.ListAuthorizers(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(auths.Items))
	for _, a := range auths.Items {
		result = append(result, authorizerResponse(a))
	}

	return listResponse("authorizers", result, auths.NextMarker), nil
}
