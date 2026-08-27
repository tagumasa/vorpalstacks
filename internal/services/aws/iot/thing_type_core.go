package iot

import (
	"time"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the thing-type family. Handlers on both protocol
// planes are thin adapters; validation and persistence live here only.

// CreateThingTypeInput carries the create parameters for a thing type.
type CreateThingTypeInput struct {
	ThingTypeName        string
	Description          string
	SearchableAttributes []string
}

// UpdateThingTypeInput carries the mutable fields of a thing type.
type UpdateThingTypeInput struct {
	ThingTypeName        string
	Description          string
	SearchableAttributes []string
}

// createThingTypeCore validates and persists a new thing type.
func (s *IoTService) createThingTypeCore(store iotstore.IotStoreInterface, in CreateThingTypeInput) (*iotstore.ThingType, error) {
	if in.ThingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}
	now := time.Now().UTC()
	tt := &iotstore.ThingType{
		ThingTypeName:        in.ThingTypeName,
		Description:          in.Description,
		SearchableAttributes: in.SearchableAttributes,
		Version:              1,
		CreationDate:         now,
		LastModifiedDate:     now,
	}
	return store.CreateThingType(tt)
}

// describeThingTypeCore resolves a thing type by name.
func (s *IoTService) describeThingTypeCore(store iotstore.IotStoreInterface, thingTypeName string) (*iotstore.ThingType, error) {
	if thingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.GetThingType(thingTypeName)
}

// updateThingTypeCore applies mutable-field updates to a thing type.
func (s *IoTService) updateThingTypeCore(store iotstore.IotStoreInterface, in UpdateThingTypeInput) (*iotstore.ThingType, error) {
	if in.ThingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}
	opts := iotstore.ThingTypeUpdateOpts{
		Description:          in.Description,
		SearchableAttributes: in.SearchableAttributes,
	}
	return store.UpdateThingType(in.ThingTypeName, opts)
}

// deleteThingTypeCore removes a type's tags and the type itself.
func (s *IoTService) deleteThingTypeCore(store iotstore.IotStoreInterface, thingTypeName string) error {
	if thingTypeName == "" {
		return iotstore.ErrMissingParam
	}
	arn := iotstore.BuildThingTypeARN(store.GetAccountID(), store.GetRegion(), thingTypeName)
	_ = store.DeleteAllTags(arn)
	return store.DeleteThingType(thingTypeName)
}

// deprecateThingTypeCore flips the deprecation flag, honouring
// undoDeprecate.
func (s *IoTService) deprecateThingTypeCore(store iotstore.IotStoreInterface, thingTypeName string, undoDeprecate bool) error {
	if thingTypeName == "" {
		return iotstore.ErrMissingParam
	}
	_, err := store.SetThingTypeDeprecation(thingTypeName, !undoDeprecate)
	return err
}

// listThingTypesCore lists thing types page by page.
func (s *IoTService) listThingTypesCore(store iotstore.IotStoreInterface, opts storecommon.ListOptions) (*storecommon.ListResult[iotstore.ThingType], error) {
	return store.ListThingTypes(opts)
}
