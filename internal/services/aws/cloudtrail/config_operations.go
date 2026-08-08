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
		return nil, awserrors.NewAWSError("ConfigurationException",
			"No event configuration found for the specified resource", 404)
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
		sizeStr, _ := v.(string)
		if sizeStr != "Standard" && sizeStr != "Large" {
			return nil, awserrors.NewAWSError("InvalidParameterException",
				"MaxEventSize must be 'Standard' or 'Large'", 400)
		}
		config["MaxEventSize"] = v
	}

	if v, ok := req.Parameters["AggregationConfigurations"]; ok {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, awserrors.NewAWSError("InvalidParameterException",
				"AggregationConfigurations must be a list", 400)
		}
		for _, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				return nil, awserrors.NewAWSError("InvalidParameterException",
					"Each AggregationConfiguration must be a map", 400)
			}
			if ec, hasEC := m["EventCategory"]; hasEC {
				ecArr, ok := ec.([]interface{})
				if !ok {
					return nil, awserrors.NewAWSError("InvalidParameterException",
						"AggregationConfiguration.EventCategory must be a list", 400)
				}
				for _, ecItem := range ecArr {
					ecStr, ok := ecItem.(string)
					if !ok || (ecStr != "insight" && ecStr != "lap" && ecStr != "management" && ecStr != "data") {
						return nil, awserrors.NewAWSError("InvalidEventCategoryException",
							"EventCategory must be one of: insight, lap, management, data", 400)
					}
				}
			}
		}
		config["AggregationConfigurations"] = v
	}

	if v, ok := req.Parameters["ContextKeySelectors"]; ok {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, awserrors.NewAWSError("InvalidParameterException",
				"ContextKeySelectors must be a list", 400)
		}
		for _, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				return nil, awserrors.NewAWSError("InvalidParameterException",
					"Each ContextKeySelector must be a map", 400)
			}
			if _, hasType := m["Type"]; !hasType {
				return nil, awserrors.NewAWSError("InvalidParameterException",
					"ContextKeySelector.Type is required", 400)
			}
		}
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
