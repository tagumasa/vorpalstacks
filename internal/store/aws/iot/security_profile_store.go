package iot

import (
	"google.golang.org/protobuf/proto"
	"time"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateSecurityProfile(sp *SecurityProfile) (*SecurityProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sp.SecurityProfileName == "" {
		return nil, ErrInvalidRequest
	}
	sp.SecurityProfileARN = BuildSecurityProfileARN(s.accountID, s.region, sp.SecurityProfileName)
	if sp.CreationDate.IsZero() {
		sp.CreationDate = time.Now().UTC()
	}
	if sp.LastModifiedDate.IsZero() {
		sp.LastModifiedDate = sp.CreationDate
	}
	return sp, s.securityProfilePS.Create(sp)
}

func (s *IotStore) GetSecurityProfile(name string) (*SecurityProfile, error) {
	return s.securityProfilePS.Get(name)
}

// UpdateSecurityProfile replaces the stored record wholesale. The service
// Core owns versioning and timestamps and performs the read-modify-write,
// so a full replace is required for cleared members (nil behaviours and
// the delete-flag semantics) to persist.
func (s *IotStore) UpdateSecurityProfile(name string, sp *SecurityProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := s.securityProfilePS.Get(name); err != nil {
		return err
	}
	return s.securityProfilePS.Update(sp)
}

func (s *IotStore) DeleteSecurityProfile(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.securityProfilePS.DeleteIfExists(name)
}

func (s *IotStore) ListSecurityProfiles(opts common.ListOptions) (*common.ListResult[SecurityProfile], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.securityProfilesBase, opts, func() *pb.SecurityProfile { return &pb.SecurityProfile{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*SecurityProfile, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToSecurityProfile(p))
	}
	return &common.ListResult[SecurityProfile]{Items: items, NextMarker: result.NextMarker}, nil
}

func (s *IotStore) ListActiveViolations(thingName string) ([]*ViolationEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var events []*ViolationEvent
	err := s.violationEventsBase.ScanPrefix("", func(key string, value []byte) error {
		pbEvent := &pb.ViolationEvent{}
		if proto.Unmarshal(value, pbEvent) == nil {
			if thingName == "" || pbEvent.ThingName == thingName {
				events = append(events, ProtoToViolationEvent(pbEvent))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *IotStore) ListViolationEvents(opts common.ListOptions, securityProfileName, thingName string) ([]*ViolationEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.violationEventsBase, opts, func() *pb.ViolationEvent { return &pb.ViolationEvent{} }, nil)
	if err != nil {
		return nil, err
	}
	events := make([]*ViolationEvent, 0, len(result.Items))
	for _, p := range result.Items {
		e := ProtoToViolationEvent(p)
		if securityProfileName != "" && e.SecurityProfileName != securityProfileName {
			continue
		}
		if thingName != "" && e.ThingName != thingName {
			continue
		}
		events = append(events, e)
	}
	return events, nil
}
