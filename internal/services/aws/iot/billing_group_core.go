package iot

import (
	"time"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the billing-group family, including the thing
// membership operations served from the device-management surface.
// Handlers on both protocol planes are thin adapters; validation and
// persistence live here only.

// CreateBillingGroupInput carries the create parameters for a billing
// group.
type CreateBillingGroupInput struct {
	GroupName   string
	Description string
}

// UpdateBillingGroupInput carries the mutable fields of a billing group
// plus the compare-and-swap version. PropertiesProvided distinguishes an
// absent billingGroupProperties member (rejected: the model marks it
// required) from a present-but-empty one.
type UpdateBillingGroupInput struct {
	GroupName          string
	Description        string
	ExpectedVersion    int64
	PropertiesProvided bool
}

// createBillingGroupCore validates and persists a new billing group.
func (s *IoTService) createBillingGroupCore(store iotstore.IotStoreInterface, in CreateBillingGroupInput) (*iotstore.BillingGroup, error) {
	if in.GroupName == "" {
		return nil, iotstore.ErrMissingParam
	}
	now := time.Now().UTC()
	bg := &iotstore.BillingGroup{
		GroupName:        in.GroupName,
		Description:      in.Description,
		CreationDate:     now,
		LastModifiedDate: now,
	}
	return store.CreateBillingGroup(bg)
}

// describeBillingGroupCore resolves a billing group by name.
func (s *IoTService) describeBillingGroupCore(store iotstore.IotStoreInterface, groupName string) (*iotstore.BillingGroup, error) {
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.GetBillingGroup(groupName)
}

// updateBillingGroupCore applies mutable-field updates under the
// expected-version compare-and-swap check enforced by the store.
func (s *IoTService) updateBillingGroupCore(store iotstore.IotStoreInterface, in UpdateBillingGroupInput) (*iotstore.BillingGroup, error) {
	if in.GroupName == "" || !in.PropertiesProvided {
		return nil, iotstore.ErrMissingParam
	}
	opts := iotstore.BillingGroupUpdateOpts{
		Description:     in.Description,
		ExpectedVersion: in.ExpectedVersion,
	}
	return store.UpdateBillingGroup(in.GroupName, opts)
}

// deleteBillingGroupCore removes a group's tags and the group itself.
func (s *IoTService) deleteBillingGroupCore(store iotstore.IotStoreInterface, groupName string) error {
	if groupName == "" {
		return iotstore.ErrMissingParam
	}
	arn := iotstore.BuildBillingGroupARN(store.GetAccountID(), store.GetRegion(), groupName)
	_ = store.DeleteAllTags(arn)
	return store.DeleteBillingGroup(groupName)
}

// listBillingGroupsCore lists billing groups page by page.
func (s *IoTService) listBillingGroupsCore(store iotstore.IotStoreInterface, opts storecommon.ListOptions) (*storecommon.ListResult[iotstore.BillingGroup], error) {
	return store.ListBillingGroups(opts)
}

// addThingToBillingGroupCore attaches a thing to a billing group.
func (s *IoTService) addThingToBillingGroupCore(store iotstore.IotStoreInterface, thingName, billingGroup string) error {
	if thingName == "" || billingGroup == "" {
		return iotstore.ErrMissingParam
	}
	return store.AddThingToBillingGroup(thingName, billingGroup)
}

// removeThingFromBillingGroupCore detaches a thing from a billing group.
func (s *IoTService) removeThingFromBillingGroupCore(store iotstore.IotStoreInterface, thingName, billingGroup string) error {
	if thingName == "" || billingGroup == "" {
		return iotstore.ErrMissingParam
	}
	return store.RemoveThingFromBillingGroup(thingName, billingGroup)
}

// listThingsInBillingGroupCore resolves the group (AWS returns
// ResourceNotFoundException for a non-existent billing group) and returns
// its thing names.
func (s *IoTService) listThingsInBillingGroupCore(store iotstore.IotStoreInterface, billingGroup string) ([]string, error) {
	if billingGroup == "" {
		return nil, iotstore.ErrMissingParam
	}
	// AWS: returns ResourceNotFoundException for a non-existent billing group.
	if _, err := store.GetBillingGroup(billingGroup); err != nil {
		return nil, err
	}
	return store.ListThingsInBillingGroup(billingGroup)
}
