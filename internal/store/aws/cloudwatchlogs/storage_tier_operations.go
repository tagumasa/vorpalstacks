package cloudwatchlogs

import (
	"time"
)

// The account-level storage tier policy family: one singleton record per
// region store. Put refreshes the update stamp; Get reports the stored
// record or the not-found sentinel — the absent policy has never been
// set, and the getter's modelled ResourceNotFoundException is the
// service Core's mapping of that sentinel.

// PutStorageTierPolicy stores the account's storage tier policy,
// stamping the update time.
func (s *Store) PutStorageTierPolicy(tier string) (*StorageTierPolicy, error) {
	policy := &StorageTierPolicy{
		StorageTier:     tier,
		LastUpdatedTime: time.Now().UnixMilli(),
	}
	if err := s.Put(s.storageTierPolicyKey(), policy); err != nil {
		return nil, err
	}
	return policy, nil
}

// GetStorageTierPolicy reads the account's storage tier policy. An
// absent record is the not-found sentinel.
func (s *Store) GetStorageTierPolicy() (*StorageTierPolicy, error) {
	return getJSONRecord[StorageTierPolicy](s, s.storageTierPolicyKey(), ErrResourceNotFound)
}
