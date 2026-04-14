package route53

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const recordSetBucketName = "route53_record_sets"

// RecordSetStore manages Route 53 record sets.
type RecordSetStore struct {
	*common.BaseStore
	mu sync.RWMutex
}

// NewRecordSetStore creates a new RecordSetStore instance with the specified storage.
func NewRecordSetStore(store storage.BasicStorage) *RecordSetStore {
	return &RecordSetStore{
		BaseStore: common.NewBaseStore(store.Bucket(recordSetBucketName), "route53"),
	}
}

// Get retrieves a resource record set from the store.
// Returns the record set or an error if not found.
func (s *RecordSetStore) Get(hostedZoneId, name, recordType, setIdentifier string) (*ResourceRecordSet, error) {
	key := recordSetKey(hostedZoneId, name, recordType, setIdentifier)
	var recordSet ResourceRecordSet
	if err := s.BaseStore.Get(key, &recordSet); err != nil {
		return nil, NewStoreError("get_record_set", err)
	}
	return &recordSet, nil
}

// List returns all resource record sets for a hosted zone.
func (s *RecordSetStore) List(hostedZoneId string) ([]*ResourceRecordSet, error) {
	prefix := hostedZoneId + "#"
	var recordSets []*ResourceRecordSet
	err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var recordSet ResourceRecordSet
		if err := json.Unmarshal(value, &recordSet); err != nil {
			return err
		}
		recordSets = append(recordSets, &recordSet)
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_record_sets", err)
	}
	return recordSets, nil
}

// Create creates a new resource record set in the store.
// Returns an error if a record set with the same key already exists.
func (s *RecordSetStore) Create(hostedZoneId string, recordSet *ResourceRecordSet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := recordSetKey(hostedZoneId, recordSet.Name, recordSet.Type, recordSet.SetIdentifier)
	if s.BaseStore.Exists(key) {
		return NewStoreError("create_record_set", common.ErrAlreadyExists)
	}
	if err := s.BaseStore.Put(key, recordSet); err != nil {
		return NewStoreError("create_record_set", err)
	}
	return nil
}

// Upsert creates or updates a resource record set in the store.
func (s *RecordSetStore) Upsert(hostedZoneId string, recordSet *ResourceRecordSet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := recordSetKey(hostedZoneId, recordSet.Name, recordSet.Type, recordSet.SetIdentifier)
	if err := s.BaseStore.Put(key, recordSet); err != nil {
		return NewStoreError("upsert_record_set", err)
	}
	return nil
}

// Delete deletes a resource record set from the store.
// Returns an error if the record set does not exist.
func (s *RecordSetStore) Delete(hostedZoneId, name, recordType, setIdentifier string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := recordSetKey(hostedZoneId, name, recordType, setIdentifier)
	if !s.BaseStore.Exists(key) {
		return NewStoreError("delete_record_set", common.ErrNotFound)
	}
	if err := s.BaseStore.Delete(key); err != nil {
		return NewStoreError("delete_record_set", err)
	}
	return nil
}

// Exists checks whether a resource record set exists in the store.
func (s *RecordSetStore) Exists(key string) bool {
	return s.BaseStore.Exists(key)
}

func recordSetKey(hostedZoneId, name, recordType, setIdentifier string) string {
	if setIdentifier != "" {
		return fmt.Sprintf("%s#%s#%s#%s", hostedZoneId, escapeKeyPart(name), recordType, escapeKeyPart(setIdentifier))
	}
	return fmt.Sprintf("%s#%s#%s", hostedZoneId, escapeKeyPart(name), recordType)
}

func escapeKeyPart(s string) string {
	return strings.ReplaceAll(s, "#", "\\#")
}


