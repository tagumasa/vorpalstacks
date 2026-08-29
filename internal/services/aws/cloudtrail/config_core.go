package cloudtrail

import (
	awserrors "vorpalstacks/internal/common/errors"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// EventConfigurationResourceInput carries the trail-name or event-data-store
// selector for the event-configuration operations.
type EventConfigurationResourceInput struct {
	TrailName      string
	EventDataStore string
}

// PutEventConfigurationInput carries the raw configuration members for
// PutEventConfiguration. The members are presence-checked by the Core so the
// validation and stored-configuration shape stay on the Core layer.
type PutEventConfigurationInput struct {
	TrailName      string
	EventDataStore string
	Params         map[string]interface{}
}

// RegisterOrganizationDelegatedAdminInput carries the member account ID.
type RegisterOrganizationDelegatedAdminInput struct {
	MemberAccountID string
}

// DeregisterOrganizationDelegatedAdminInput carries the delegated admin
// account ID.
type DeregisterOrganizationDelegatedAdminInput struct {
	DelegatedAdminAccountID string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// getEventConfigurationCore is the single entry point for
// GetEventConfiguration.
func (s *CloudTrailService) getEventConfigurationCore(store cloudtrailstore.CloudTrailStoreInterface, in EventConfigurationResourceInput) (map[string]interface{}, error) {
	if in.TrailName == "" && in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter",
			"Either TrailName or EventDataStore is required", 400)
	}

	config, err := store.GetEventConfiguration(in.TrailName, in.EventDataStore)
	if err != nil {
		return nil, awserrors.NewAWSError("ConfigurationException",
			"No event configuration found for the specified resource", 404)
	}

	return config, nil
}

// putEventConfigurationCore is the single entry point for
// PutEventConfiguration: it validates the MaxEventSize enum, the aggregation
// configuration shapes and event categories, and the context key selectors
// before persisting the configuration.
func (s *CloudTrailService) putEventConfigurationCore(store cloudtrailstore.CloudTrailStoreInterface, in PutEventConfigurationInput) error {
	if in.TrailName == "" && in.EventDataStore == "" {
		return awserrors.NewAWSError("InvalidParameter",
			"Either TrailName or EventDataStore is required", 400)
	}

	config := map[string]interface{}{}
	if in.TrailName != "" {
		config["TrailARN"] = in.TrailName
	}
	if in.EventDataStore != "" {
		config["EventDataStoreArn"] = in.EventDataStore
	}

	if v, ok := in.Params["MaxEventSize"]; ok {
		sizeStr, _ := v.(string)
		if sizeStr != "Standard" && sizeStr != "Large" {
			return awserrors.NewAWSError("InvalidParameterException",
				"MaxEventSize must be 'Standard' or 'Large'", 400)
		}
		config["MaxEventSize"] = v
	}

	if v, ok := in.Params["AggregationConfigurations"]; ok {
		arr, ok := v.([]interface{})
		if !ok {
			return awserrors.NewAWSError("InvalidParameterException",
				"AggregationConfigurations must be a list", 400)
		}
		for _, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				return awserrors.NewAWSError("InvalidParameterException",
					"Each AggregationConfiguration must be a map", 400)
			}
			if ec, hasEC := m["EventCategory"]; hasEC {
				ecArr, ok := ec.([]interface{})
				if !ok {
					return awserrors.NewAWSError("InvalidParameterException",
						"AggregationConfiguration.EventCategory must be a list", 400)
				}
				for _, ecItem := range ecArr {
					ecStr, ok := ecItem.(string)
					if !ok || (ecStr != "insight" && ecStr != "lap" && ecStr != "management" && ecStr != "data") {
						return awserrors.NewAWSError("InvalidEventCategoryException",
							"EventCategory must be one of: insight, lap, management, data", 400)
					}
				}
			}
		}
		config["AggregationConfigurations"] = v
	}

	if v, ok := in.Params["ContextKeySelectors"]; ok {
		arr, ok := v.([]interface{})
		if !ok {
			return awserrors.NewAWSError("InvalidParameterException",
				"ContextKeySelectors must be a list", 400)
		}
		for _, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				return awserrors.NewAWSError("InvalidParameterException",
					"Each ContextKeySelector must be a map", 400)
			}
			if _, hasType := m["Type"]; !hasType {
				return awserrors.NewAWSError("InvalidParameterException",
					"ContextKeySelector.Type is required", 400)
			}
		}
		config["ContextKeySelectors"] = v
	}

	if err := store.PutEventConfiguration(in.TrailName, in.EventDataStore, config); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// registerOrganizationDelegatedAdminCore is the single entry point for
// RegisterOrganizationDelegatedAdmin.
func (s *CloudTrailService) registerOrganizationDelegatedAdminCore(store cloudtrailstore.CloudTrailStoreInterface, in RegisterOrganizationDelegatedAdminInput) (map[string]interface{}, error) {
	if in.MemberAccountID == "" {
		return nil, awserrors.NewAWSError("InvalidParameter",
			"MemberAccountId is required", 400)
	}

	if err := store.RegisterDelegatedAdmin(in.MemberAccountID); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"DelegatedAdminAccountId": in.MemberAccountID,
	}, nil
}

// deregisterOrganizationDelegatedAdminCore is the single entry point for
// DeregisterOrganizationDelegatedAdmin.
func (s *CloudTrailService) deregisterOrganizationDelegatedAdminCore(store cloudtrailstore.CloudTrailStoreInterface, in DeregisterOrganizationDelegatedAdminInput) error {
	if in.DelegatedAdminAccountID == "" {
		return awserrors.NewAWSError("InvalidParameter",
			"DelegatedAdminAccountId is required", 400)
	}

	if err := store.DeregisterDelegatedAdmin(in.DelegatedAdminAccountID); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}
