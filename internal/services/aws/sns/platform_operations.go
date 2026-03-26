package sns

import (
	"context"
	"strconv"
	"strings"

	"vorpalstacks/internal/services/aws/common/pagination"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

func parsePlatformAttributes(params map[string]interface{}) map[string]string {
	result := make(map[string]string)

	if attrs, ok := params["Attributes"].(map[string]interface{}); ok {
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
	}
	if attrs, ok := params["attributes"].(map[string]interface{}); ok {
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
	}

	for i := 1; ; i++ {
		keyKey := "Attributes.entry." + strconv.Itoa(i) + ".key"
		valueKey := "Attributes.entry." + strconv.Itoa(i) + ".value"

		key := request.GetStringParam(params, keyKey)
		if key == "" {
			break
		}

		value := request.GetStringParam(params, valueKey)
		result[key] = value
	}

	return result
}

func parseEndpointAttributes(params map[string]interface{}) map[string]string {
	result := make(map[string]string)

	if attrs, ok := params["Attributes"].(map[string]interface{}); ok {
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
	}
	if attrs, ok := params["attributes"].(map[string]interface{}); ok {
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
	}

	for i := 1; ; i++ {
		keyKey := "Attributes.entry." + strconv.Itoa(i) + ".key"
		valueKey := "Attributes.entry." + strconv.Itoa(i) + ".value"

		key := request.GetStringParam(params, keyKey)
		if key == "" {
			break
		}

		value := request.GetStringParam(params, valueKey)
		result[key] = value
	}

	return result
}

// CreatePlatformApplication creates a platform application for push notifications.
func (s *SNSService) CreatePlatformApplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	platform := request.GetStringParam(req.Parameters, "Platform")

	if name == "" {
		return nil, NewInvalidParameterException("Name is required")
	}
	if platform == "" {
		return nil, NewInvalidParameterException("Platform is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	app := &snsstore.PlatformApplication{
		Name:       name,
		Platform:   strings.ToUpper(platform),
		Attributes: parsePlatformAttributes(req.Parameters),
	}

	created, err := store.CreatePlatformApplication(app)
	if err != nil {
		if err == snsstore.ErrPlatformApplicationAlreadyExists {
			return nil, NewInvalidParameterException("Platform application already exists with the same name")
		}
		return nil, err
	}

	return map[string]interface{}{
		"PlatformApplicationArn": created.PlatformApplicationArn,
	}, nil
}

// DeletePlatformApplication deletes a platform application.
func (s *SNSService) DeletePlatformApplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	platformApplicationArn := request.GetStringParam(req.Parameters, "PlatformApplicationArn")
	if platformApplicationArn == "" {
		return nil, NewInvalidParameterException("PlatformApplicationArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeletePlatformApplication(platformApplicationArn); err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetPlatformApplicationAttributes retrieves the attributes of a platform application.
func (s *SNSService) GetPlatformApplicationAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	platformApplicationArn := request.GetStringParam(req.Parameters, "PlatformApplicationArn")
	if platformApplicationArn == "" {
		return nil, NewInvalidParameterException("PlatformApplicationArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	attrs, err := store.GetPlatformApplicationAttributes(platformApplicationArn)
	if err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// SetPlatformApplicationAttributes sets the attributes of a platform application.
func (s *SNSService) SetPlatformApplicationAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	platformApplicationArn := request.GetStringParam(req.Parameters, "PlatformApplicationArn")
	if platformApplicationArn == "" {
		return nil, NewInvalidParameterException("PlatformApplicationArn is required")
	}

	attrs := parsePlatformAttributes(req.Parameters)
	if len(attrs) == 0 {
		return nil, NewInvalidParameterException("Attributes is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.SetPlatformApplicationAttributes(platformApplicationArn, attrs); err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListPlatformApplications lists all platform applications.
func (s *SNSService) ListPlatformApplications(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	result, err := store.ListPlatformApplications(common.ListOptions{Marker: nextToken})
	if err != nil {
		return nil, err
	}

	apps := make([]map[string]interface{}, 0, len(result.Items))
	for _, app := range result.Items {
		apps = append(apps, map[string]interface{}{
			"PlatformApplicationArn": app.PlatformApplicationArn,
			"Attributes":             app.Attributes,
		})
	}

	token := ""
	if result.IsTruncated && result.NextMarker != "" {
		token = result.NextMarker
	}
	return pagination.BuildListResponse("PlatformApplications", apps, token), nil
}

// CreatePlatformEndpoint creates a platform endpoint for push notifications.
func (s *SNSService) CreatePlatformEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	platformApplicationArn := request.GetStringParam(req.Parameters, "PlatformApplicationArn")
	token := request.GetStringParam(req.Parameters, "Token")
	customUserData := request.GetStringParam(req.Parameters, "CustomUserData")

	if platformApplicationArn == "" {
		return nil, NewInvalidParameterException("PlatformApplicationArn is required")
	}
	if token == "" {
		return nil, NewInvalidParameterException("Token is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	endpoint := &snsstore.PlatformEndpoint{
		PlatformApplicationArn: platformApplicationArn,
		Token:                  token,
		CustomUserData:         customUserData,
		Attributes:             parseEndpointAttributes(req.Parameters),
	}

	created, err := store.CreatePlatformEndpoint(endpoint)
	if err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"EndpointArn": created.EndpointArn,
	}, nil
}

// DeleteEndpoint deletes a platform endpoint.
func (s *SNSService) DeleteEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	endpointArn := request.GetStringParam(req.Parameters, "EndpointArn")
	if endpointArn == "" {
		return nil, NewInvalidParameterException("EndpointArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteEndpoint(endpointArn); err != nil {
		if err == snsstore.ErrEndpointNotFound {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetEndpointAttributes retrieves the attributes of a platform endpoint.
func (s *SNSService) GetEndpointAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	endpointArn := request.GetStringParam(req.Parameters, "EndpointArn")
	if endpointArn == "" {
		return nil, NewInvalidParameterException("EndpointArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	attrs, err := store.GetEndpointAttributes(endpointArn)
	if err != nil {
		if err == snsstore.ErrEndpointNotFound {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// SetEndpointAttributes sets the attributes of a platform endpoint.
func (s *SNSService) SetEndpointAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	endpointArn := request.GetStringParam(req.Parameters, "EndpointArn")
	if endpointArn == "" {
		return nil, NewInvalidParameterException("EndpointArn is required")
	}

	attrs := parseEndpointAttributes(req.Parameters)
	if len(attrs) == 0 {
		return nil, NewInvalidParameterException("Attributes is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.SetEndpointAttributes(endpointArn, attrs); err != nil {
		if err == snsstore.ErrEndpointNotFound {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListEndpointsByPlatformApplication lists endpoints for a platform application.
func (s *SNSService) ListEndpointsByPlatformApplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	platformApplicationArn := request.GetStringParam(req.Parameters, "PlatformApplicationArn")
	if platformApplicationArn == "" {
		return nil, NewInvalidParameterException("PlatformApplicationArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	result, err := store.ListEndpointsByPlatformApplication(platformApplicationArn, common.ListOptions{Marker: nextToken})
	if err != nil {
		return nil, err
	}

	endpoints := make([]map[string]interface{}, 0, len(result.Items))
	for _, ep := range result.Items {
		endpoints = append(endpoints, map[string]interface{}{
			"EndpointArn": ep.EndpointArn,
			"Attributes":  ep.Attributes,
		})
	}

	token := ""
	if result.IsTruncated && result.NextMarker != "" {
		token = result.NextMarker
	}
	return pagination.BuildListResponse("Endpoints", endpoints, token), nil
}
