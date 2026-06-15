package iot

import (
		"github.com/google/uuid"
	"time"
	"vorpalstacks/internal/store/aws/common"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
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
		return nil, ErrBillingGroupNotFound
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

// CreateAuthorizer persists a new custom authorizer.
