package sns

import (
	"context"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// validPlatforms lists the push notification platforms accepted by SNS.
var validPlatforms = map[string]bool{
	"APNS":         true,
	"APNS_SANDBOX": true,
	"GCM":          true,
	"ADM":          true,
	"BAIDU":        true,
	"WNS":          true,
	"MPNS":         true,
}

// CreatePlatformApplication creates a platform application for push notifications.
func (s *SNSService) CreatePlatformApplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	platform := request.GetStringParam(req.Parameters, "Platform")

	if name == "" {
		return nil, awserrors.NewInvalidParameterException("Name is required")
	}
	if len(name) > 100 {
		return nil, awserrors.NewInvalidParameterException(fmt.Sprintf("Name too long: %d characters (maximum 100)", len(name)))
	}
	if platform == "" {
		return nil, awserrors.NewInvalidParameterException("Platform is required")
	}
	normalisedPlatform := strings.ToUpper(platform)
	if !validPlatforms[normalisedPlatform] {
		return nil, awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid Platform: %s. Valid values: APNS, APNS_SANDBOX, GCM, ADM, BAIDU, WNS, MPNS", platform))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	attrs := parseAttributes(req.Parameters)
	for attrName, value := range attrs {
		if err := validatePlatformAttributeValue(attrName, value); err != nil {
			return nil, err
		}
	}

	app := &snsstore.PlatformApplication{
		Name:       name,
		Platform:   normalisedPlatform,
		Attributes: attrs,
	}

	created, err := store.CreatePlatformApplication(app)
	if err != nil {
		if err == snsstore.ErrPlatformApplicationAlreadyExists {
			return nil, awserrors.NewInvalidParameterException("Platform application already exists with the same name")
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
		return nil, awserrors.NewInvalidParameterException("PlatformApplicationArn is required")
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
		return nil, awserrors.NewInvalidParameterException("PlatformApplicationArn is required")
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
		return nil, awserrors.NewInvalidParameterException("PlatformApplicationArn is required")
	}

	attrs := parseAttributes(req.Parameters)
	if len(attrs) == 0 {
		return nil, awserrors.NewInvalidParameterException("Attributes is required")
	}

	// Enforce value length cap for DoS protection.
	for name, value := range attrs {
		if err := validatePlatformAttributeValue(name, value); err != nil {
			return nil, err
		}
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

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
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

	// Validate PlatformApplicationArn format before proceeding.
	if err := validatePlatformApplicationArn(platformApplicationArn); err != nil {
		return nil, err
	}
	if token == "" {
		return nil, awserrors.NewInvalidParameterException("Token is required")
	}
	if len(token) > 2048 {
		return nil, awserrors.NewInvalidParameterException(fmt.Sprintf("Token too long: %d bytes (maximum 2048)", len(token)))
	}
	if len(customUserData) > 2048 {
		return nil, awserrors.NewInvalidParameterException(fmt.Sprintf("CustomUserData too long: %d bytes (maximum 2048)", len(customUserData)))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	attrs := parseAttributes(req.Parameters)
	for attrName, value := range attrs {
		if err := validatePlatformAttributeValue(attrName, value); err != nil {
			return nil, err
		}
	}

	endpoint := &snsstore.PlatformEndpoint{
		PlatformApplicationArn: platformApplicationArn,
		Token:                  token,
		CustomUserData:         customUserData,
		Attributes:             attrs,
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
		return nil, awserrors.NewInvalidParameterException("EndpointArn is required")
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
		return nil, awserrors.NewInvalidParameterException("EndpointArn is required")
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
		return nil, awserrors.NewInvalidParameterException("EndpointArn is required")
	}

	attrs := parseAttributes(req.Parameters)
	if len(attrs) == 0 {
		return nil, awserrors.NewInvalidParameterException("Attributes is required")
	}

	// Enforce value length cap for DoS protection.
	for name, value := range attrs {
		if err := validatePlatformAttributeValue(name, value); err != nil {
			return nil, err
		}
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
		return nil, awserrors.NewInvalidParameterException("PlatformApplicationArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Verify platform application existence before listing endpoints.
	// Without this check, a non-existent ARN returns an empty result instead
	// of NotFoundException.
	if _, err := store.GetPlatformApplication(platformApplicationArn); err != nil {
		if err == snsstore.ErrPlatformApplicationNotFound {
			return nil, ErrPlatformAppNotFound
		}
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
