package iot

import (
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/store/aws/common"
)

// GenericKV provides JSON-backed persistence for gap-implementation entities
// that do not warrant a dedicated proto message. All state is stored under
// the "iot-generic-kv" bucket with keys of the form "<category>/<name>".

func (s *IotStore) PutGeneric(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.genericKVBase.Put(key, value)
}

func (s *IotStore) GetGeneric(key string, dest interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.genericKVBase.Get(key, dest)
}

// GetGenericExists fetches a generic-KV record and reports whether it exists.
// A missing key returns (false, nil) so callers can distinguish "not found"
// from genuine store errors.
func (s *IotStore) GetGenericExists(key string, dest interface{}) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.genericKVBase.Get(key, dest); err != nil {
		if common.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *IotStore) DeleteGeneric(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.genericKVBase.Delete(key)
}

func (s *IotStore) ListGeneric(prefix string) ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var items []map[string]interface{}
	err := s.genericKVBase.ScanPrefix(prefix, func(_ string, value []byte) error {
		item := map[string]interface{}{}
		// Surface unmarshal failures instead of silently skipping corrupt
		// entries; hiding them would make data loss invisible to callers.
		if err := json.Unmarshal(value, &item); err != nil {
			return fmt.Errorf("generic-kv: decode entry for prefix %q: %w", prefix, err)
		}
		items = append(items, item)
		return nil
	})
	return items, err
}
