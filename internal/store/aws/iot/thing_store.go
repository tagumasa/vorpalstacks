package iot

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func (s *IotStore) CreateThing(thing *Thing) (*Thing, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if thing.ThingName == "" {
		return nil, ErrInvalidRequest
	}
	if thing.ThingTypeName != "" {
		tt := &pb.ThingType{}
		if err := s.thingTypesBase.GetProto(thing.ThingTypeName, tt); err != nil {
			if common.IsNotFound(err) {
				return nil, ErrThingTypeNotFound
			}
			return nil, err
		}
		if tt.Deprecated {
			return nil, ErrInvalidRequest
		}
	}
	if thing.ThingID == "" {
		thing.ThingID = uuid.New().String()
	}
	thing.ThingARN = BuildThingARN(s.accountID, s.region, thing.ThingName)
	return thing, s.thingPS.Create(thing)
}

func (s *IotStore) GetThing(thingName string) (*Thing, error) {
	return s.thingPS.Get(thingName)
}

func (s *IotStore) UpdateThing(thingName string, opts ThingUpdateOpts) (*Thing, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, err := s.thingPS.Get(thingName)
	if err != nil {
		return nil, err
	}
	if opts.RemoveThingType {
		existing.ThingTypeName = ""
	} else if opts.ThingTypeName != "" {
		tt := &pb.ThingType{}
		if err := s.thingTypesBase.GetProto(opts.ThingTypeName, tt); err != nil {
			if common.IsNotFound(err) {
				return nil, ErrThingTypeNotFound
			}
			return nil, err
		}
		if tt.Deprecated {
			return nil, ErrInvalidRequest
		}
		existing.ThingTypeName = opts.ThingTypeName
	}
	if opts.PayloadProvided {
		if opts.MergeAttributes {
			if existing.Attributes == nil {
				existing.Attributes = make(map[string]string)
			}
			for k, v := range opts.Attributes {
				if v == "" {
					delete(existing.Attributes, k)
				} else {
					existing.Attributes[k] = v
				}
			}
		} else {
			existing.Attributes = make(map[string]string)
			for k, v := range opts.Attributes {
				if v != "" {
					existing.Attributes[k] = v
				}
			}
		}
		existing.AttributeNames = mapKeys(existing.Attributes)
	}
	existing.Version++
	existing.LastModifiedDate = time.Now().UTC()
	return existing, s.thingPS.Update(existing)
}

// errFound is a sentinel error used to short-circuit ScanPrefix
// once the first matching entry has been found.
var errFound = errors.New("found")

func (s *IotStore) DeleteThing(thingName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.thingsBase.Exists(thingName) {
		return ErrThingNotFound
	}
	hasAttachments := false
	s.thingPrincipalBase.ScanPrefix(thingName+"\x00", func(key string, _ []byte) error {
		hasAttachments = true
		return errFound
	})
	if hasAttachments {
		return ErrDeleteConflict
	}
	return s.deleteThingTxn(thingName)
}

func (s *IotStore) deleteThingTxn(thingName string) error {
	// Acquire thing-level shadow lock to prevent concurrent shadow
	// operations from racing with shadow cleanup (AWS spec: DeleteThing
	// atomically removes all shadows for the thing).
	s.shadowLocker.Lock(thingShadowLockKey(thingName))
	defer s.shadowLocker.Unlock(thingShadowLockKey(thingName))

	// Remove ThingGroup and BillingGroup memberships to prevent dangling
	// references in their respective member buckets.
	thingsBucket := bucketThings + s.rs
	shadowsBucket := bucketShadows + s.rs
	tgmBucket := bucketThingGroupMember + s.rs
	gtmBucket := bucketGroupThingMember + s.rs
	tbmBucket := bucketThingBillingMember + s.rs
	btmBucket := bucketBillingThingMember + s.rs
	return s.ts.Update(context.Background(), func(txn storage.Transaction) error {
		tgmB := txn.Bucket(tgmBucket)
		gtmB := txn.Bucket(gtmBucket)
		tbmB := txn.Bucket(tbmBucket)
		btmB := txn.Bucket(btmBucket)

		// --- ThingGroup membership cleanup ---
		iter := tgmB.ScanPrefix([]byte(thingName + "\x00"))
		var groupThingKeys [][]byte
		var thingGroupKeys [][]byte
		for iter.Next() {
			k := iter.Key()
			copied := make([]byte, len(k))
			copy(copied, k)
			thingGroupKeys = append(thingGroupKeys, copied)
			parts := strings.SplitN(string(k), "\x00", 2)
			if len(parts) >= 2 {
				gtk := []byte(parts[1] + "\x00" + thingName)
				copied2 := make([]byte, len(gtk))
				copy(copied2, gtk)
				groupThingKeys = append(groupThingKeys, copied2)
			}
		}
		if err := iter.Error(); err != nil {
			return err
		}
		for _, k := range groupThingKeys {
			if err := gtmB.Delete(k); err != nil {
				return err
			}
		}
		for _, k := range thingGroupKeys {
			if err := tgmB.Delete(k); err != nil {
				return err
			}
		}

		// --- BillingGroup membership cleanup ---
		biter := tbmB.ScanPrefix([]byte(thingName + "\x00"))
		var billingThingKeys [][]byte
		var thingBillingKeys [][]byte
		for biter.Next() {
			k := biter.Key()
			copied := make([]byte, len(k))
			copy(copied, k)
			thingBillingKeys = append(thingBillingKeys, copied)
			parts := strings.SplitN(string(k), "\x00", 2)
			if len(parts) >= 2 {
				btk := []byte(parts[1] + "\x00" + thingName)
				copied2 := make([]byte, len(btk))
				copy(copied2, btk)
				billingThingKeys = append(billingThingKeys, copied2)
			}
		}
		if err := biter.Error(); err != nil {
			return err
		}
		for _, k := range billingThingKeys {
			if err := btmB.Delete(k); err != nil {
				return err
			}
		}
		for _, k := range thingBillingKeys {
			if err := tbmB.Delete(k); err != nil {
				return err
			}
		}

		// --- Shadow cleanup ---
		shadowB := txn.Bucket(shadowsBucket)
		sIter := shadowB.ScanPrefix([]byte(thingName + "/"))
		var shadowKeys [][]byte
		for sIter.Next() {
			k := sIter.Key()
			copied := make([]byte, len(k))
			copy(copied, k)
			shadowKeys = append(shadowKeys, copied)
		}
		if err := sIter.Error(); err != nil {
			return err
		}
		for _, k := range shadowKeys {
			if err := shadowB.Delete(k); err != nil {
				return err
			}
		}
		return txn.Bucket(thingsBucket).Delete([]byte(thingName))
	})
}

func (s *IotStore) ListThings(opts common.ListOptions, attributeName, attributeValue string) (*common.ListResult[Thing], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var filter func(*pb.Thing) bool
	if attributeName != "" {
		filter = func(p *pb.Thing) bool {
			v, ok := p.Attributes[attributeName]
			return ok && v == attributeValue
		}
	}
	result, err := common.ListProto(s.thingsBase, opts, func() *pb.Thing { return &pb.Thing{} }, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*Thing, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToThing(p))
	}
	return &common.ListResult[Thing]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}

func (s *IotStore) ListThingsForThingType(thingTypeName string, opts common.ListOptions) (*common.ListResult[Thing], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	filter := func(p *pb.Thing) bool { return p.ThingTypeName == thingTypeName }
	result, err := common.ListProto(s.thingsBase, opts, func() *pb.Thing { return &pb.Thing{} }, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*Thing, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProtoToThing(p))
	}
	return &common.ListResult[Thing]{
		Items:      items,
		NextMarker: result.NextMarker,
	}, nil
}
