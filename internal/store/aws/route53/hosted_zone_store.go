package route53

// Package route53 provides Route 53 data store implementations for vorpalstacks.

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const hostedZoneBucketName = "route53_hosted_zones"

// HostedZoneStore manages Route 53 hosted zones.
type HostedZoneStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewHostedZoneStore creates a new HostedZoneStore instance with the specified storage and account ID.
func NewHostedZoneStore(store storage.BasicStorage, accountId string) *HostedZoneStore {
	return &HostedZoneStore{
		BaseStore:  common.NewBaseStore(store.Bucket(hostedZoneBucketName), "route53"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a hosted zone by its ID from the store.
// Returns the hosted zone or an error if not found.
func (s *HostedZoneStore) Get(id string) (*HostedZone, error) {
	var zone HostedZone
	if err := s.BaseStore.Get(id, &zone); err != nil {
		return nil, NewStoreError("get_hosted_zone", err)
	}
	return &zone, nil
}

// GetByName retrieves a hosted zone by its name from the store.
// Returns the hosted zone or an error if not found.
func (s *HostedZoneStore) GetByName(name string) (*HostedZone, error) {
	name = normalizeZoneName(name)
	var found *HostedZone
	err := s.ForEach(func(key string, value []byte) error {
		var zone HostedZone
		if err := json.Unmarshal(value, &zone); err != nil {
			return err
		}
		if normalizeZoneName(zone.Name) == name && found == nil {
			found = &zone
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_hosted_zone_by_name", err)
	}
	if found == nil {
		return nil, NewStoreError("get_hosted_zone_by_name", common.ErrNotFound)
	}
	return found, nil
}

// List returns a list of hosted zones from the store with pagination support.
func (s *HostedZoneStore) List(marker string, maxItems int) (*HostedZoneListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var zones []*HostedZone
	count := 0
	started := marker == ""
	hasMore := false
	var lastKey string

	err := s.ForEach(func(key string, value []byte) error {
		var zone HostedZone
		if err := json.Unmarshal(value, &zone); err != nil {
			return err
		}

		if !started {
			if zone.ID == marker || zone.ID > marker {
				started = true
			}
			if !started {
				return nil
			}
		}

		if count < maxItems {
			zones = append(zones, &zone)
			count++
			lastKey = zone.ID
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_hosted_zones", err)
	}

	return &HostedZoneListResult{
		HostedZones: zones,
		IsTruncated: hasMore,
		Marker:      lastKey,
	}, nil
}

// ListByName returns a list of hosted zones sorted by name.
func (s *HostedZoneStore) ListByName() ([]*HostedZone, error) {
	var zones []*HostedZone
	err := s.ForEach(func(key string, value []byte) error {
		var zone HostedZone
		if err := json.Unmarshal(value, &zone); err != nil {
			return err
		}
		zones = append(zones, &zone)
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_hosted_zones_by_name", err)
	}

	for i := 0; i < len(zones)-1; i++ {
		for j := i + 1; j < len(zones); j++ {
			if strings.Compare(zones[i].Name, zones[j].Name) > 0 {
				zones[i], zones[j] = zones[j], zones[i]
			}
		}
	}

	return zones, nil
}

// Create creates a new hosted zone in the store.
// Returns an error if the hosted zone already exists.
func (s *HostedZoneStore) Create(zone *HostedZone) error {
	if s.Exists(zone.ID) {
		return NewStoreError("create_hosted_zone", common.ErrAlreadyExists)
	}
	zone.CreatedAt = time.Now()
	if err := s.BaseStore.Put(zone.ID, zone); err != nil {
		return NewStoreError("create_hosted_zone", err)
	}
	return nil
}

// Update updates an existing hosted zone in the store.
// Returns an error if the hosted zone does not exist.
func (s *HostedZoneStore) Update(zone *HostedZone) error {
	if !s.Exists(zone.ID) {
		return NewStoreError("update_hosted_zone", common.ErrNotFound)
	}
	if err := s.BaseStore.Put(zone.ID, zone); err != nil {
		return NewStoreError("update_hosted_zone", err)
	}
	return nil
}

// Delete deletes a hosted zone by its ID from the store.
// Returns an error if the hosted zone does not exist.
func (s *HostedZoneStore) Delete(id string) error {
	if !s.Exists(id) {
		return NewStoreError("delete_hosted_zone", common.ErrNotFound)
	}
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_hosted_zone", err)
	}
	return nil
}

// Exists checks whether a hosted zone exists in the store by its ID.
func (s *HostedZoneStore) Exists(id string) bool {
	return s.BaseStore.Exists(id)
}

func normalizeZoneName(name string) string {
	name = strings.ToLower(name)
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	return name
}
