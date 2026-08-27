package iot

import (
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the fleet-indexing configuration pair. Handlers on
// both protocol planes are thin adapters; validation and persistence live
// here only.

// UpdateIndexingConfigurationInput carries the indexing modes with
// explicit presence: a nil member leaves the stored mode untouched, and
// absent configuration altogether starts from a fresh (empty) record —
// matching the previous partial-update semantics.
type UpdateIndexingConfigurationInput struct {
	ThingIndexingMode             *string
	ThingGroupIndexingMode        *string
	ThingConnectivityIndexingMode *string
}

// getIndexingConfigurationCore loads the indexing configuration. A nil
// result means no configuration is readable or stored; AWS reports the
// default OFF mode in that case.
func (s *IoTService) getIndexingConfigurationCore(store iotstore.IotStoreInterface) *iotstore.IndexingConfiguration {
	ic, err := store.GetIndexingConfiguration()
	if err != nil {
		return nil
	}
	return ic
}

// updateIndexingConfigurationCore merges the supplied modes into the
// stored configuration, creating it when absent.
func (s *IoTService) updateIndexingConfigurationCore(store iotstore.IotStoreInterface, in UpdateIndexingConfigurationInput) error {
	ic, err := store.GetIndexingConfiguration()
	if err != nil {
		ic = &iotstore.IndexingConfiguration{}
	}
	if in.ThingIndexingMode != nil {
		ic.ThingIndexingMode = *in.ThingIndexingMode
	}
	if in.ThingGroupIndexingMode != nil {
		ic.ThingGroupIndexingMode = *in.ThingGroupIndexingMode
	}
	if in.ThingConnectivityIndexingMode != nil {
		ic.ThingConnectivityIndexingMode = *in.ThingConnectivityIndexingMode
	}
	return store.UpdateIndexingConfiguration(ic)
}
