package cloudtrail

import (
	stderrors "errors"

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
		return nil, newInvalidParameterException(
			"Either TrailName or EventDataStore is required")
	}
	// The two target members are mutually exclusive, exactly as the Put
	// path enforces.
	if in.TrailName != "" && in.EventDataStore != "" {
		return nil, newInvalidParameterCombinationException(
			"TrailName and EventDataStore cannot be used together")
	}

	// The target resolves first and the configuration is looked up under
	// the resolved identity — the same key the Put path persists under —
	// so a configuration written through a name is found through an ARN
	// and vice versa, and a resource that does not exist answers the
	// resource type's declared not-found error from the resolution itself.
	var trailName, edsID string
	if in.TrailName != "" {
		trail, err := s.resolveTrailCore(store, in.TrailName)
		if err != nil {
			return nil, err
		}
		trailName = trail.Name
	} else {
		eds, err := s.resolveEventDataStore(store, in.EventDataStore)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		edsID = eds.EventDataStoreID
	}

	config, err := store.GetEventConfiguration(trailName, edsID)
	if err != nil {
		// A resolved resource without a stored configuration answers the
		// resource type's declared not-found error — the platform records
		// no default configuration to serve in its place; every other
		// store failure stays on the mapStoreError path.
		if stderrors.Is(err, cloudtrailstore.ErrEventConfigurationNotFound) {
			if trailName != "" {
				return nil, ErrTrailNotFound
			}
			return nil, ErrEventDataStoreNotFoundException
		}
		return nil, s.mapStoreError(err)
	}

	return config, nil
}

// putEventConfigurationCore is the single entry point for
// PutEventConfiguration: it resolves exactly one target resource (a trail
// or an event data store — never both), validates the MaxEventSize enum,
// the aggregation configuration shapes and event categories, and the
// context key selectors before persisting the configuration. The persisted
// configuration is the response — the operation's response members are
// exactly the stored ones.
func (s *CloudTrailService) putEventConfigurationCore(store cloudtrailstore.CloudTrailStoreInterface, in PutEventConfigurationInput) (map[string]interface{}, error) {
	if in.TrailName == "" && in.EventDataStore == "" {
		return nil, newInvalidParameterException(
			"Either TrailName or EventDataStore is required")
	}
	// The two target members are mutually exclusive; the operation declares
	// InvalidParameterCombinationException for a broken combination.
	if in.TrailName != "" && in.EventDataStore != "" {
		return nil, newInvalidParameterCombinationException(
			"TrailName and EventDataStore cannot be used together")
	}

	config := map[string]interface{}{}
	// The resolved identities are also the storage keys: the Put persists
	// under the canonical trail name and event data store ID so the Get
	// path — which resolves the same way — always finds the record.
	var keyTrailName, keyEDSID string
	if in.TrailName != "" {
		// The response member is the trail's ARN, so the stored record must
		// carry the ARN of the resolved trail, never the caller-supplied
		// name; a trail that does not exist is rejected with the operation's
		// declared not-found error.
		trail, err := s.resolveTrailCore(store, in.TrailName)
		if err != nil {
			return nil, err
		}
		config["TrailARN"] = trail.TrailARN
		keyTrailName = trail.Name
	}
	if in.EventDataStore != "" {
		// The event data store target resolves through the store like the
		// trail target: an absent store is the operation's declared
		// not-found error, and the stored configuration carries the store's
		// canonical ARN.
		eds, err := s.resolveEventDataStore(store, in.EventDataStore)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		config["EventDataStoreArn"] = eds.EventDataStoreARN
		keyEDSID = eds.EventDataStoreID
	}

	if v, ok := in.Params["MaxEventSize"]; ok {
		sizeStr, _ := v.(string)
		if sizeStr != "Standard" && sizeStr != "Large" {
			return nil, newInvalidParameterException(
				"MaxEventSize must be 'Standard' or 'Large'")
		}
		config["MaxEventSize"] = v
	}

	if v, ok := in.Params["AggregationConfigurations"]; ok {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, newInvalidParameterException(
				"AggregationConfigurations must be a list")
		}
		// The model bounds the list to a single configuration.
		if len(arr) > 1 {
			return nil, newInvalidParameterException(
				"AggregationConfigurations must contain at most one configuration")
		}
		for _, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				return nil, newInvalidParameterException(
					"Each AggregationConfiguration must be a map")
			}
			// Templates is model-required with length 1-50 and the Template
			// enum values.
			templatesRaw, hasTemplates := m["Templates"]
			if !hasTemplates {
				return nil, newInvalidParameterException(
					"AggregationConfiguration.Templates is required")
			}
			templates, ok := templatesRaw.([]interface{})
			if !ok || len(templates) < 1 || len(templates) > 50 {
				return nil, newInvalidParameterException(
					"AggregationConfiguration.Templates must contain between 1 and 50 templates")
			}
			for _, tRaw := range templates {
				tStr, ok := tRaw.(string)
				if !ok || !validAggregationTemplates[tStr] {
					return nil, newInvalidParameterException(
						"Templates must contain only API_ACTIVITY, RESOURCE_ACCESS, or USER_ACTIONS")
				}
			}
			// EventCategory is a model-required scalar carrying the single
			// enum value Data.
			ecStr, hasEC := m["EventCategory"].(string)
			if !hasEC || ecStr != "Data" {
				return nil, newInvalidParameterException(
					"AggregationConfiguration.EventCategory must be Data")
			}
		}
		config["AggregationConfigurations"] = v
	}

	if v, ok := in.Params["ContextKeySelectors"]; ok {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, newInvalidParameterException(
				"ContextKeySelectors must be a list")
		}
		// The model bounds the list to two selectors.
		if len(arr) > 2 {
			return nil, newInvalidParameterException(
				"ContextKeySelectors must contain at most two selectors")
		}
		for _, item := range arr {
			m, ok := item.(map[string]interface{})
			if !ok {
				return nil, newInvalidParameterException(
					"Each ContextKeySelector must be a map")
			}
			typeStr, hasType := m["Type"].(string)
			if !hasType || !validContextKeyTypes[typeStr] {
				return nil, newInvalidParameterException(
					"ContextKeySelector.Type must be TagContext or RequestContext")
			}
			// Equals is model-required with length 1-50.
			equalsRaw, hasEquals := m["Equals"]
			if !hasEquals {
				return nil, newInvalidParameterException(
					"ContextKeySelector.Equals is required")
			}
			equals, ok := equalsRaw.([]interface{})
			if !ok || len(equals) < 1 || len(equals) > 50 {
				return nil, newInvalidParameterException(
					"ContextKeySelector.Equals must contain between 1 and 50 values")
			}
			for _, eRaw := range equals {
				if _, ok := eRaw.(string); !ok {
					return nil, newInvalidParameterException(
						"ContextKeySelector.Equals values must be strings")
				}
			}
		}
		config["ContextKeySelectors"] = v
	}

	if err := store.PutEventConfiguration(keyTrailName, keyEDSID, config); err != nil {
		return nil, s.mapStoreError(err)
	}

	return config, nil
}

// registerOrganizationDelegatedAdminCore is the single entry point for
// RegisterOrganizationDelegatedAdmin. Accounts on this platform never
// belong to an organization — no Organizations substrate exists — so the
// operation's documented refusal for a non-member account is the only
// reachable outcome: OrganizationsNotInUseException answers after the
// member validation.
func (s *CloudTrailService) registerOrganizationDelegatedAdminCore(in RegisterOrganizationDelegatedAdminInput) (map[string]interface{}, error) {
	if in.MemberAccountID == "" {
		return nil, newInvalidParameterException(
			"MemberAccountId is required")
	}

	return nil, newOrganizationsNotInUseException(
		"The request is made from an account that is not a member of an organization")
}

// deregisterOrganizationDelegatedAdminCore is the single entry point for
// DeregisterOrganizationDelegatedAdmin. Like the register direction, the
// account is never an organization member, so the documented non-member
// refusal answers.
func (s *CloudTrailService) deregisterOrganizationDelegatedAdminCore(in DeregisterOrganizationDelegatedAdminInput) error {
	if in.DelegatedAdminAccountID == "" {
		return newInvalidParameterException(
			"DelegatedAdminAccountId is required")
	}

	return newOrganizationsNotInUseException(
		"The request is made from an account that is not a member of an organization")
}
