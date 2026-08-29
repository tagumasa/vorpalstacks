package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// GetEventConfiguration retrieves the event configuration for a trail or EDS.
func (s *CloudTrailService) GetEventConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.getEventConfigurationCore(store, EventConfigurationResourceInput{
		TrailName:      request.GetStringParam(req.Parameters, "TrailName"),
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
	})
}

// PutEventConfiguration stores event configuration for a trail or EDS.
func (s *CloudTrailService) PutEventConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := s.putEventConfigurationCore(store, PutEventConfigurationInput{
		TrailName:      request.GetStringParam(req.Parameters, "TrailName"),
		EventDataStore: request.GetStringParam(req.Parameters, "EventDataStore"),
		Params:         req.Parameters,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RegisterOrganizationDelegatedAdmin registers a delegated administrator
// account for CloudTrail in an AWS Organization.
func (s *CloudTrailService) RegisterOrganizationDelegatedAdmin(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.registerOrganizationDelegatedAdminCore(store, RegisterOrganizationDelegatedAdminInput{
		MemberAccountID: request.GetStringParam(req.Parameters, "MemberAccountId"),
	})
}

// DeregisterOrganizationDelegatedAdmin removes a delegated administrator
// account for CloudTrail.
func (s *CloudTrailService) DeregisterOrganizationDelegatedAdmin(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if err := s.deregisterOrganizationDelegatedAdminCore(store, DeregisterOrganizationDelegatedAdminInput{
		DelegatedAdminAccountID: request.GetStringParam(req.Parameters, "DelegatedAdminAccountId"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
