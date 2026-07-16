package cloudtrail

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// GetEventConfiguration retrieves the event configuration for a trail or EDS.
func (s *CloudTrailService) GetEventConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trailName := request.GetStringParam(req.Parameters, "TrailName")
	edsID := request.GetStringParam(req.Parameters, "EventDataStore")

	if trailName == "" && edsID == "" {
		return nil, awserrors.NewAWSError("InvalidParameter",
			"Either TrailName or EventDataStore is required", 400)
	}

	config, err := store.GetEventConfiguration(trailName, edsID)
	if err != nil {
		return map[string]interface{}{
			"TrailARN":          "",
			"EventDataStoreArn": "",
		}, nil
	}

	return config, nil
}

// PutEventConfiguration stores event configuration for a trail or EDS.
func (s *CloudTrailService) PutEventConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	trailName := request.GetStringParam(req.Parameters, "TrailName")
	edsID := request.GetStringParam(req.Parameters, "EventDataStore")

	if trailName == "" && edsID == "" {
		return nil, awserrors.NewAWSError("InvalidParameter",
			"Either TrailName or EventDataStore is required", 400)
	}

	config := map[string]interface{}{}
	if trailName != "" {
		config["TrailARN"] = trailName
	}
	if edsID != "" {
		config["EventDataStoreArn"] = edsID
	}

	if v, ok := req.Parameters["MaxEventSize"]; ok {
		config["MaxEventSize"] = v
	}
	if v, ok := req.Parameters["AggregationConfigurations"]; ok {
		config["AggregationConfigurations"] = v
	}
	if v, ok := req.Parameters["ContextKeySelectors"]; ok {
		config["ContextKeySelectors"] = v
	}

	if err := store.PutEventConfiguration(trailName, edsID, config); err != nil {
		return nil, s.mapStoreError(err)
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

	memberAccountID := request.GetStringParam(req.Parameters, "MemberAccountId")
	if memberAccountID == "" {
		return nil, awserrors.NewAWSError("InvalidParameter",
			"MemberAccountId is required", 400)
	}

	if err := store.RegisterDelegatedAdmin(memberAccountID); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"DelegatedAdminAccountId": memberAccountID,
	}, nil
}

// DeregisterOrganizationDelegatedAdmin removes a delegated administrator
// account for CloudTrail.
func (s *CloudTrailService) DeregisterOrganizationDelegatedAdmin(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	memberAccountID := request.GetStringParam(req.Parameters, "DelegatedAdminAccountId")
	if memberAccountID == "" {
		return nil, awserrors.NewAWSError("InvalidParameter",
			"DelegatedAdminAccountId is required", 400)
	}

	if err := store.DeregisterDelegatedAdmin(memberAccountID); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}
