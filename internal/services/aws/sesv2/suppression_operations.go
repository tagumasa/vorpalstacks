package sesv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// GetSuppressedDestination retrieves details about a suppressed email destination.
func (s *SESv2Service) GetSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getSuppressedDestinationCore(store, request.GetStringParam(req.Parameters, "EmailAddress"))
}

// PutSuppressedDestination adds or updates a suppressed destination.
func (s *SESv2Service) PutSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putSuppressedDestinationCore(store, PutSuppressedDestinationInput{
		EmailAddress: request.GetStringParam(req.Parameters, "EmailAddress"),
		Reason:       request.GetStringParam(req.Parameters, "Reason"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// DeleteSuppressedDestination removes an email address from the suppression list.
func (s *SESv2Service) DeleteSuppressedDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteSuppressedDestinationCore(store, request.GetStringParam(req.Parameters, "EmailAddress")); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListSuppressedDestinations returns a list of suppressed destinations.
func (s *SESv2Service) ListSuppressedDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectTenantName(req.Parameters); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listSuppressedDestinationsCore(store, ListSuppressedDestinationsInput{
		MaxItems:  pagination.GetMaxItems(req.Parameters, 100, "PageSize"),
		NextToken: pagination.GetMarker(req.Parameters, "NextToken"),
		Reasons:   request.GetStringList(req.Parameters, "Reasons"),
		StartDate: parseTimestampParam(req.Parameters, "StartDate"),
		EndDate:   parseTimestampParam(req.Parameters, "EndDate"),
	})
}
