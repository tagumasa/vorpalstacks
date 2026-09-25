// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// UpdateGlobalTableSettings request list bounds, from the API model's
// length traits: the replica settings update list carries at most 50
// entries, and each replica's per-index settings update list at most 20.
const (
	MaxReplicaSettingsUpdates    = 50
	MaxReplicaGSISettingsUpdates = 20
)

// globalTableBucketName is the bucket the global-table record lives in.
// Global-table membership is account-global metadata — one record every
// region reads and writes identically — so the bucket has no region
// suffix: it sits in the account-global storage, and a region-suffixed
// name would recreate the per-region record the account-global binding
// exists to prevent.
func globalTableBucketName() string {
	return "dynamodb_global_tables"
}

// globalTableMu serialises read-modify-write cycles over the global-table
// record space. Every store instance binds to the same account-global
// bucket, so the lock spans all of them: a per-instance mutex could not
// exclude a concurrent Update running through another region's store.
var globalTableMu sync.Mutex

// GlobalTableStore manages DynamoDB global tables in persistent storage.
type GlobalTableStore struct {
	*common.BaseStore
	arnBuilder *svcarn.DynamoDBBuilder
}

// NewGlobalTableStore creates a new store for DynamoDB global tables. The
// store argument is the account-global storage handle, not the caller's
// regional storage: the record must resolve identically from every region.
func NewGlobalTableStore(store storage.BasicStorage, accountId, region string) *GlobalTableStore {
	return &GlobalTableStore{
		BaseStore:  common.NewBaseStore(store.Bucket(globalTableBucketName()), "dynamodb_global_tables"),
		arnBuilder: svcarn.NewARNBuilder(accountId, region).DynamoDB(),
	}
}

// Get retrieves a global table by name.
func (s *GlobalTableStore) Get(name string) (*GlobalTable, error) {
	var pbGlobalTable storage_dynamodb.GlobalTable
	if err := s.BaseStore.GetProto(name, &pbGlobalTable); err != nil {
		return nil, err
	}
	return ProtoToGlobalTable(&pbGlobalTable), nil
}

// Create creates a new global table.
func (s *GlobalTableStore) Create(name string, replicationGroup []*Replica) (*GlobalTable, error) {
	globalTableMu.Lock()
	defer globalTableMu.Unlock()
	if s.Exists(name) {
		return nil, ErrTableAlreadyExists
	}

	now := time.Now().UTC()
	globalTable := &GlobalTable{
		GlobalTableName:   name,
		GlobalTableArn:    s.arnBuilder.GlobalTable(name),
		GlobalTableStatus: GlobalTableStatusActive,
		CreationDateTime:  now,
		ReplicationGroup:  replicationGroup,
	}

	if err := s.put(globalTable); err != nil {
		return nil, err
	}

	return globalTable, nil
}

// put writes the global-table record without locking; the caller must
// already hold the store lock.
func (s *GlobalTableStore) put(globalTable *GlobalTable) error {
	return s.BaseStore.PutProto(globalTable.GlobalTableName, GlobalTableToProto(globalTable))
}

// Update loads the global table under the store lock, applies mutate to it,
// and persists the result. A mutate error aborts without writing. This is
// the locked read-modify-write path for global-table records: the storage
// layer is read-committed, so an unlocked Get→Put sequence can lose
// concurrent replica or settings changes.
func (s *GlobalTableStore) Update(name string, mutate func(*GlobalTable) error) (*GlobalTable, error) {
	globalTableMu.Lock()
	defer globalTableMu.Unlock()
	globalTable, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	if err := mutate(globalTable); err != nil {
		return nil, err
	}
	if err := s.put(globalTable); err != nil {
		return nil, err
	}
	return globalTable, nil
}

// Delete deletes a global table by name. The deletion takes the store
// lock: an unlocked delete could interleave with a locked
// read-modify-write cycle over the same record and resurrect state the
// cycle was removing.
func (s *GlobalTableStore) Delete(name string) error {
	globalTableMu.Lock()
	defer globalTableMu.Unlock()
	return s.BaseStore.Delete(name)
}

// DeleteIfEmpty removes the global table record only when its replication
// group has no members left, holding the store lock across the emptiness
// read and the delete — the atomic conditional delete. An unlocked
// check-then-delete could erase a record a concurrent member re-add had
// just repopulated between the two steps. The boolean reports whether the
// record was deleted; an absent record answers (false, nil).
func (s *GlobalTableStore) DeleteIfEmpty(name string) (bool, error) {
	globalTableMu.Lock()
	defer globalTableMu.Unlock()
	globalTable, err := s.Get(name)
	if err != nil {
		if common.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	if len(globalTable.ReplicationGroup) != 0 {
		return false, nil
	}
	if err := s.BaseStore.Delete(name); err != nil {
		if common.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Exists checks if a global table exists.
func (s *GlobalTableStore) Exists(name string) bool {
	return s.BaseStore.Exists(name)
}

// List lists global tables.
func (s *GlobalTableStore) List(marker string, limit int) ([]*GlobalTable, string, error) {
	return listProtoConverted(s.BaseStore, marker, limit, func() *storage_dynamodb.GlobalTable { return &storage_dynamodb.GlobalTable{} }, ProtoToGlobalTable, nil)
}
