package iot

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"time"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateThingType(tt *ThingType) (*ThingType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if tt.ThingTypeName == "" {
		return nil, ErrInvalidRequest
	}
	if tt.ThingTypeID == "" {
		tt.ThingTypeID = uuid.New().String()
	}
	tt.ThingTypeARN = BuildThingTypeARN(s.accountID, s.region, tt.ThingTypeName)
	return tt, s.thingTypePS.Create(tt)
}

func (s *IotStore) GetThingType(name string) (*ThingType, error) {
	return s.thingTypePS.Get(name)
}

func (s *IotStore) DeleteThingType(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hasThings := false
	s.thingsBase.ScanPrefix("", func(key string, val []byte) error {
		p := &pb.Thing{}
		if proto.Unmarshal(val, p) == nil && p.ThingTypeName == name {
			hasThings = true
			return errFound
		}
		return nil
	})
	if hasThings {
		return ErrDeleteConflict
	}
	return s.thingTypePS.DeleteIfExists(name)
}

func (s *IotStore) ListThingTypes(opts common.ListOptions) (*common.ListResult[ThingType], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, err := common.ListProto(s.thingTypesBase, opts, func() *pb.ThingType { return &pb.ThingType{} }, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*ThingType, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToThingType(p))
	}
	return &common.ListResult[ThingType]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

func (s *IotStore) UpdateThingType(name string, opts ThingTypeUpdateOpts) (*ThingType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.thingTypePS.Get(name)
	if err != nil {
		return nil, ErrThingTypeNotFound
	}
	if opts.Description != "" {
		existing.Description = opts.Description
	}
	existing.Version++
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.thingTypePS.Update(existing)
}

func (s *IotStore) SetThingTypeDeprecation(name string, deprecated bool) (*ThingType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.thingTypePS.Get(name)
	if err != nil {
		return nil, ErrThingTypeNotFound
	}
	existing.Deprecated = deprecated
	if deprecated {
		existing.DeprecationDate = time.Now().UTC()
	} else {
		existing.DeprecationDate = time.Time{}
	}
	return existing, s.thingTypePS.Update(existing)
}
