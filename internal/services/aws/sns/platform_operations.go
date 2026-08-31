package sns

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
)

// CreatePlatformApplication creates a platform application for push notifications.
func (s *SNSService) CreatePlatformApplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.createPlatformApplicationCore(store, CreatePlatformApplicationInput{
		Name:       request.GetStringParam(req.Parameters, "Name"),
		Platform:   request.GetStringParam(req.Parameters, "Platform"),
		Attributes: parseAttributes(req.Parameters),
	})
}

// DeletePlatformApplication deletes a platform application.
func (s *SNSService) DeletePlatformApplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.deletePlatformApplicationCore(store, DeletePlatformApplicationInput{
		PlatformApplicationArn: request.GetStringParam(req.Parameters, "PlatformApplicationArn"),
	})
}

// GetPlatformApplicationAttributes retrieves the attributes of a platform application.
func (s *SNSService) GetPlatformApplicationAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.getPlatformApplicationAttributesCore(store, GetPlatformApplicationAttributesInput{
		PlatformApplicationArn: request.GetStringParam(req.Parameters, "PlatformApplicationArn"),
	})
}

// SetPlatformApplicationAttributes sets the attributes of a platform application.
func (s *SNSService) SetPlatformApplicationAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.setPlatformApplicationAttributesCore(store, SetPlatformApplicationAttributesInput{
		PlatformApplicationArn: request.GetStringParam(req.Parameters, "PlatformApplicationArn"),
		Attributes:             parseAttributes(req.Parameters),
	})
}

// ListPlatformApplications lists all platform applications.
func (s *SNSService) ListPlatformApplications(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.listPlatformApplicationsCore(store, ListPlatformApplicationsInput{
		NextToken: pagination.GetMarker(req.Parameters, "NextToken"),
	})
}

// CreatePlatformEndpoint creates a platform endpoint for push notifications.
func (s *SNSService) CreatePlatformEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.createPlatformEndpointCore(store, CreatePlatformEndpointInput{
		PlatformApplicationArn: request.GetStringParam(req.Parameters, "PlatformApplicationArn"),
		Token:                  request.GetStringParam(req.Parameters, "Token"),
		CustomUserData:         request.GetStringParam(req.Parameters, "CustomUserData"),
		Attributes:             parseAttributes(req.Parameters),
	})
}

// DeleteEndpoint deletes a platform endpoint.
func (s *SNSService) DeleteEndpoint(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.deleteEndpointCore(store, DeleteEndpointInput{
		EndpointArn: request.GetStringParam(req.Parameters, "EndpointArn"),
	})
}

// GetEndpointAttributes retrieves the attributes of a platform endpoint.
func (s *SNSService) GetEndpointAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.getEndpointAttributesCore(store, GetEndpointAttributesInput{
		EndpointArn: request.GetStringParam(req.Parameters, "EndpointArn"),
	})
}

// SetEndpointAttributes sets the attributes of a platform endpoint.
func (s *SNSService) SetEndpointAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.setEndpointAttributesCore(store, SetEndpointAttributesInput{
		EndpointArn: request.GetStringParam(req.Parameters, "EndpointArn"),
		Attributes:  parseAttributes(req.Parameters),
	})
}

// ListEndpointsByPlatformApplication lists endpoints for a platform application.
func (s *SNSService) ListEndpointsByPlatformApplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.listEndpointsByPlatformApplicationCore(store, ListEndpointsByPlatformApplicationInput{
		PlatformApplicationArn: request.GetStringParam(req.Parameters, "PlatformApplicationArn"),
		NextToken:              request.GetStringParam(req.Parameters, "NextToken"),
	})
}
