package iot

import (
		"github.com/google/uuid"
	"time"
	"strings"
	"context"
	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
)
func (s *IotStore) CreateThingGroup(group *ThingGroup) (*ThingGroup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if group.GroupName == "" {
		return nil, ErrInvalidRequest
	}
	if group.ParentGroupName != "" && !s.thingGroupsBase.Exists(group.ParentGroupName) {
		return nil, ErrThingGroupNotFound
	}
	if group.GroupID == "" {
		group.GroupID = uuid.New().String()
	}
	group.GroupARN = BuildThingGroupARN(s.accountID, s.region, group.GroupName)
	if group.Version == 0 {
		group.Version = 1
	}
	return group, s.thingGroupPS.Create(group)
}

func (s *IotStore) GetThingGroup(name string) (*ThingGroup, error) {
	return s.thingGroupPS.Get(name)
}

func (s *IotStore) DeleteThingGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hasMembers := false
	s.groupThingMemberBase.ScanPrefix(name+"\x00", func(key string, _ []byte) error {
		hasMembers = true
		return errFound
	})
	if hasMembers {
		return ErrDeleteConflict
	}
	hasChildren := false
	s.thingGroupsBase.ScanPrefix("", func(key string, val []byte) error {
		g := &pb.ThingGroup{}
		if proto.Unmarshal(val, g) == nil && g.ParentGroupName == name {
			hasChildren = true
			return errFound
		}
		return nil
	})
	if hasChildren {
		return ErrDeleteConflict
	}
	return s.thingGroupPS.DeleteIfExists(name)
}

func (s *IotStore) ListThingGroups(opts common.ListOptions, parentGroupName string) (*common.ListResult[ThingGroup], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var filter func(*pb.ThingGroup) bool
	if parentGroupName != "" {
		filter = func(g *pb.ThingGroup) bool { return g.ParentGroupName == parentGroupName }
	}
	result, err := common.ListProto(s.thingGroupsBase, opts, func() *pb.ThingGroup { return &pb.ThingGroup{} }, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*ThingGroup, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToThingGroup(p))
	}
	return &common.ListResult[ThingGroup]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

func (s *IotStore) UpdateThingGroup(groupName string, opts ThingGroupUpdateOpts) (*ThingGroup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.thingGroupPS.Get(groupName)
	if err != nil {
		return nil, ErrThingGroupNotFound
	}
	if opts.ExpectedVersion > 0 && existing.Version != opts.ExpectedVersion {
		return nil, ErrVersionConflict
	}
	if opts.Description != "" {
		existing.Description = opts.Description
	}
	if len(opts.Attributes) > 0 {
		if existing.Attributes == nil {
			existing.Attributes = make(map[string]string)
		}
		for k, v := range opts.Attributes {
			existing.Attributes[k] = v
		}
	}
	existing.Version++
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.thingGroupPS.Update(existing)
}

// CreateBillingGroup persists a new billing group.
func (s *IotStore) AddThingToThingGroup(thingName, groupName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.thingsBase.Exists(thingName) {
		return ErrThingNotFound
	}
	if !s.thingGroupsBase.Exists(groupName) {
		return ErrThingGroupNotFound
	}
	tgk := thingName + "\x00" + groupName
	gtk := groupName + "\x00" + thingName
	tgmBucket := "iot-thing-group-member" + s.rs
	gtmBucket := "iot-group-thing-member" + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(tgmBucket).Put([]byte(tgk), []byte("1")); err != nil {
			return err
		}
		return txn.Bucket(gtmBucket).Put([]byte(gtk), []byte("1"))
	})
}

func (s *IotStore) RemoveThingFromThingGroup(thingName, groupName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.thingsBase.Exists(thingName) {
		return ErrThingNotFound
	}
	if !s.thingGroupsBase.Exists(groupName) {
		return ErrThingGroupNotFound
	}
	tgk := thingName + "\x00" + groupName
	gtk := groupName + "\x00" + thingName
	tgmBucket := "iot-thing-group-member" + s.rs
	gtmBucket := "iot-group-thing-member" + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(tgmBucket).Delete([]byte(tgk)); err != nil {
			return err
		}
		return txn.Bucket(gtmBucket).Delete([]byte(gtk))
	})
}

func (s *IotStore) ListThingsInGroup(groupName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var things []string
	err := s.groupThingMemberBase.ScanPrefix(groupName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			things = append(things, parts[1])
		}
		return nil
	})
	return things, err
}

func (s *IotStore) ListGroupsForThing(thingName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var groups []string
	err := s.thingGroupMemberBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			groups = append(groups, parts[1])
		}
		return nil
	})
	return groups, err
}
