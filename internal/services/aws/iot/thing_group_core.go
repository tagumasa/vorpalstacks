package iot

import (
	"time"

	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Core functions for the thing-group family (static groups). Handlers on
// both protocol planes are thin adapters that parse the wire request into
// the Input structs below, call these functions, and serialise the
// results; validation and persistence live here only.

// CreateThingGroupInput carries the create parameters for a static thing
// group.
type CreateThingGroupInput struct {
	GroupName       string
	ParentGroupName string
	Description     string
	Attributes      map[string]string
}

// UpdateThingGroupInput carries the mutable fields of a thing group plus
// the compare-and-swap version. PropertiesProvided distinguishes an absent
// thingGroupProperties member (rejected: the model marks it required) from
// a present-but-empty one.
type UpdateThingGroupInput struct {
	GroupName          string
	Description        string
	Attributes         map[string]string
	ExpectedVersion    int64
	PropertiesProvided bool
}

// createThingGroupCore validates and persists a new thing group, returning
// the stored record including its ARN and version.
func (s *IoTService) createThingGroupCore(store iotstore.IotStoreInterface, in CreateThingGroupInput) (*iotstore.ThingGroup, error) {
	if in.GroupName == "" {
		return nil, iotstore.ErrMissingParam
	}
	now := time.Now().UTC()
	group := &iotstore.ThingGroup{
		GroupName:        in.GroupName,
		ParentGroupName:  in.ParentGroupName,
		Description:      in.Description,
		Attributes:       in.Attributes,
		CreationDate:     now,
		LastModifiedDate: now,
	}
	return store.CreateThingGroup(group)
}

// describeThingGroupCore resolves a thing group by name.
func (s *IoTService) describeThingGroupCore(store iotstore.IotStoreInterface, groupName string) (*iotstore.ThingGroup, error) {
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.GetThingGroup(groupName)
}

// updateThingGroupCore applies mutable-field updates under the
// expected-version compare-and-swap check enforced by the store.
func (s *IoTService) updateThingGroupCore(store iotstore.IotStoreInterface, in UpdateThingGroupInput) (*iotstore.ThingGroup, error) {
	if in.GroupName == "" || !in.PropertiesProvided {
		return nil, iotstore.ErrMissingParam
	}
	opts := iotstore.ThingGroupUpdateOpts{
		Description:     in.Description,
		Attributes:      in.Attributes,
		ExpectedVersion: in.ExpectedVersion,
	}
	return store.UpdateThingGroup(in.GroupName, opts)
}

// deleteThingGroupCore removes a group's tags and the group itself.
func (s *IoTService) deleteThingGroupCore(store iotstore.IotStoreInterface, groupName string) error {
	if groupName == "" {
		return iotstore.ErrMissingParam
	}
	arn := iotstore.BuildThingGroupARN(store.GetAccountID(), store.GetRegion(), groupName)
	_ = store.DeleteAllTags(arn)
	return store.DeleteThingGroup(groupName)
}

// listThingGroupsCore lists groups page by page with the optional parent
// filter.
func (s *IoTService) listThingGroupsCore(store iotstore.IotStoreInterface, opts storecommon.ListOptions, parentFilter string) (*storecommon.ListResult[iotstore.ThingGroup], error) {
	return store.ListThingGroups(opts, parentFilter)
}

// addThingToThingGroupCore attaches a thing to a static group.
func (s *IoTService) addThingToThingGroupCore(store iotstore.IotStoreInterface, thingName, groupName string) error {
	if thingName == "" || groupName == "" {
		return iotstore.ErrMissingParam
	}
	return store.AddThingToThingGroup(thingName, groupName)
}

// removeThingFromThingGroupCore detaches a thing from a static group.
func (s *IoTService) removeThingFromThingGroupCore(store iotstore.IotStoreInterface, thingName, groupName string) error {
	if thingName == "" || groupName == "" {
		return iotstore.ErrMissingParam
	}
	return store.RemoveThingFromThingGroup(thingName, groupName)
}

// listThingsInThingGroupCore returns the thing names belonging to a group.
func (s *IoTService) listThingsInThingGroupCore(store iotstore.IotStoreInterface, groupName string) ([]string, error) {
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.ListThingsInGroup(groupName)
}

// listThingGroupsForThingCore returns the group names a thing belongs to.
func (s *IoTService) listThingGroupsForThingCore(store iotstore.IotStoreInterface, thingName string) ([]string, error) {
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}
	return store.ListGroupsForThing(thingName)
}

// updateThingGroupsForThingCore atomically applies a batch of group
// memberships to add and remove for one thing.
func (s *IoTService) updateThingGroupsForThingCore(store iotstore.IotStoreInterface, thingName string, add, remove []string) error {
	if thingName == "" {
		return iotstore.ErrMissingParam
	}
	for _, g := range add {
		if err := store.AddThingToThingGroup(thingName, g); err != nil {
			return err
		}
	}
	for _, g := range remove {
		if err := store.RemoveThingFromThingGroup(thingName, g); err != nil {
			return err
		}
	}
	return nil
}
