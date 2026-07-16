package iot

import (
	"strings"
	pb "vorpalstacks/internal/pb/storage/storage_iot"
	"vorpalstacks/internal/store/aws/common"
)

func shadowKey(thingName, shadowName string) string {
	if shadowName == "" || shadowName == "classic" {
		return thingName + "/$current"
	}
	return thingName + "/" + shadowName
}

// thingShadowLockKey returns the lock key used to serialise thing-level
// shadow mutations. All shadow write operations and DeleteThing's shadow
// cleanup acquire this lock to prevent concurrent modification.
// Lock ordering is always: thing-level lock first, then per-shadow lock.
func thingShadowLockKey(thingName string) string {
	return thingName + "/"
}

func (s *IotStore) GetShadow(thingName, shadowName string) (*ShadowDocument, error) {
	pb := &pb.ShadowDocument{}
	if err := s.shadowsBase.GetProto(shadowKey(thingName, shadowName), pb); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrShadowNotFound
		}
		return nil, err
	}
	return ProtoToShadow(pb), nil
}

func (s *IotStore) PutShadow(thingName, shadowName string, doc *ShadowDocument) error {
	s.shadowLocker.Lock(thingShadowLockKey(thingName))
	defer s.shadowLocker.Unlock(thingShadowLockKey(thingName))

	lockKey := shadowKey(thingName, shadowName)
	s.shadowLocker.Lock(lockKey)
	defer s.shadowLocker.Unlock(lockKey)
	return s.shadowsBase.PutProto(shadowKey(thingName, shadowName), ShadowToProto(doc))
}

// PutShadowWithVersion atomically reads the current shadow, checks the version,
// and writes the new document if the version matches.  Returns ErrVersionConflict
// if clientVersion > 0 and does not match the stored version.
func (s *IotStore) PutShadowWithVersion(thingName, shadowName string, doc *ShadowDocument, clientVersion int64) error {
	s.shadowLocker.Lock(thingShadowLockKey(thingName))
	defer s.shadowLocker.Unlock(thingShadowLockKey(thingName))

	lockKey := shadowKey(thingName, shadowName)
	s.shadowLocker.Lock(lockKey)
	defer s.shadowLocker.Unlock(lockKey)

	pb := &pb.ShadowDocument{}
	err := s.shadowsBase.GetProto(shadowKey(thingName, shadowName), pb)
	if err == nil {
		stored := ProtoToShadow(pb)
		if clientVersion > 0 && stored.VersionNumber != clientVersion {
			return ErrVersionConflict
		}
	} else if !common.IsNotFound(err) {
		return err
	}
	return s.shadowsBase.PutProto(shadowKey(thingName, shadowName), ShadowToProto(doc))
}

func (s *IotStore) DeleteShadow(thingName, shadowName string) error {
	s.shadowLocker.Lock(thingShadowLockKey(thingName))
	defer s.shadowLocker.Unlock(thingShadowLockKey(thingName))

	lockKey := shadowKey(thingName, shadowName)
	s.shadowLocker.Lock(lockKey)
	defer s.shadowLocker.Unlock(lockKey)
	return s.shadowsBase.Delete(shadowKey(thingName, shadowName))
}

// ListShadowNames returns the distinct named shadows for a thing, paginated.
// The marker is the last shadow name of the previous page; pagination is
// applied in-memory after collecting the distinct names (per-thing shadow
// counts are small). An empty nextMarker means no further page.
func (s *IotStore) ListShadowNames(thingName string, opts common.ListOptions) ([]string, string, error) {
	prefix := thingName + "/"
	var names []string
	seen := make(map[string]bool)
	if err := s.shadowsBase.ScanPrefix(prefix, func(key string, _ []byte) error {
		parts := strings.SplitN(key[len(prefix):], "/", 2)
		name := parts[0]
		if name == "$current" {
			return nil
		}
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
		return nil
	}); err != nil {
		return nil, "", err
	}
	// Apply marker-based pagination. Names arrive sorted from the scan.
	start := 0
	if opts.Marker != "" {
		found := false
		for i, n := range names {
			if n == opts.Marker {
				start = i + 1
				found = true
				break
			}
		}
		// Marker not found — return empty page instead of repeating first page.
		if !found {
			return []string{}, "", nil
		}
	}
	page := names[start:]
	var nextMarker string
	maxItems := opts.MaxItems
	if maxItems <= 0 {
		maxItems = 100
	}
	if maxItems > 0 && len(page) > maxItems {
		page = page[:maxItems]
		nextMarker = page[len(page)-1]
	}
	return page, nextMarker, nil
}
