package iot

import (
	"context"
	"github.com/google/uuid"
	"strings"
	"time"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateBillingGroup(bg *BillingGroup) (*BillingGroup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if bg.GroupName == "" {
		return nil, ErrInvalidRequest
	}
	existing := &pb.BillingGroup{}
	if err := s.billingGroupsBase.GetProto(bg.GroupName, existing); err == nil {
		return nil, ErrBillingGroupAlreadyExists
	}
	if bg.GroupID == "" {
		bg.GroupID = uuid.New().String()
	}
	bg.GroupARN = BuildBillingGroupARN(s.accountID, s.region, bg.GroupName)
	if bg.Version == 0 {
		bg.Version = 1
	}
	return bg, s.billingGroupPS.Create(bg)
}

func (s *IotStore) GetBillingGroup(name string) (*BillingGroup, error) {
	return s.billingGroupPS.Get(name)
}

func (s *IotStore) UpdateBillingGroup(name string, opts BillingGroupUpdateOpts) (*BillingGroup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.billingGroupPS.Get(name)
	if err != nil {
		return nil, err
	}
	if opts.ExpectedVersion > 0 && existing.Version != opts.ExpectedVersion {
		return nil, ErrVersionConflict
	}
	if opts.Description != "" {
		existing.Description = opts.Description
	}
	existing.Version++
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.billingGroupPS.Update(existing)
}

func (s *IotStore) DeleteBillingGroup(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hasMembers := false
	if err := s.billingThingMemberBase.ScanPrefix(name+"\x00", func(key string, _ []byte) error {
		hasMembers = true
		return errFound
	}); err != nil && err != errFound {
		return err
	}
	if hasMembers {
		return ErrDeleteConflict
	}
	return s.billingGroupPS.DeleteIfExists(name)
}

// ListBillingGroups returns all billing groups.
func (s *IotStore) ListBillingGroups(opts common.ListOptions) (*common.ListResult[BillingGroup], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.billingGroupsBase, opts, func() *pb.BillingGroup { return &pb.BillingGroup{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*BillingGroup, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToBillingGroup(p))
	}
	return &common.ListResult[BillingGroup]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

func (s *IotStore) AddThingToBillingGroup(thingName, billingGroup string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.thingsBase.Exists(thingName) {
		return ErrThingNotFound
	}
	if !s.billingGroupsBase.Exists(billingGroup) {
		return ErrBillingGroupNotFound
	}
	// AWS: a thing can belong to at most one billing group. If already in
	// a different group, reject with InvalidRequestException.
	existingGroup := ""
	if err := s.thingBillingMemberBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			existingGroup = parts[1]
		}
		return errFound
	}); err != nil && err != errFound {
		return err
	}
	if existingGroup != "" && existingGroup != billingGroup {
		return ErrThingAlreadyInBillingGroup
	}
	tbk := thingName + "\x00" + billingGroup
	btk := billingGroup + "\x00" + thingName
	tbmBucket := bucketThingBillingMember + s.rs
	btmBucket := bucketBillingThingMember + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(tbmBucket).Put([]byte(tbk), []byte("1")); err != nil {
			return err
		}
		return txn.Bucket(btmBucket).Put([]byte(btk), []byte("1"))
	})
}

func (s *IotStore) RemoveThingFromBillingGroup(thingName, billingGroup string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.thingsBase.Exists(thingName) {
		return ErrThingNotFound
	}
	if !s.billingGroupsBase.Exists(billingGroup) {
		return ErrBillingGroupNotFound
	}
	tbk := thingName + "\x00" + billingGroup
	btk := billingGroup + "\x00" + thingName
	tbmBucket := bucketThingBillingMember + s.rs
	btmBucket := bucketBillingThingMember + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		if err := txn.Bucket(tbmBucket).Delete([]byte(tbk)); err != nil {
			return err
		}
		return txn.Bucket(btmBucket).Delete([]byte(btk))
	})
}

func (s *IotStore) ListThingsInBillingGroup(billingGroup string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var things []string
	err := s.billingThingMemberBase.ScanPrefix(billingGroup+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			things = append(things, parts[1])
		}
		return nil
	})
	return things, err
}

func (s *IotStore) ListBillingGroupsForThing(thingName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var groups []string
	err := s.thingBillingMemberBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		parts := strings.SplitN(key, "\x00", 2)
		if len(parts) == 2 {
			groups = append(groups, parts[1])
		}
		return nil
	})
	return groups, err
}
