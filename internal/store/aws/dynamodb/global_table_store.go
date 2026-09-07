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

func globalTableBucketName(region string) string {
	return "dynamodb_global_tables-" + region
}

// GlobalTableStore manages DynamoDB global tables in persistent storage.
type GlobalTableStore struct {
	*common.BaseStore
	arnBuilder *svcarn.DynamoDBBuilder
	mu         sync.Mutex
}

// NewGlobalTableStore creates a new store for DynamoDB global tables.
func NewGlobalTableStore(store storage.BasicStorage, accountId, region string) *GlobalTableStore {
	return &GlobalTableStore{
		BaseStore:  common.NewBaseStore(store.Bucket(globalTableBucketName(region)), "dynamodb_global_tables"),
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Exists(name) {
		return nil, ErrTableAlreadyExists
	}

	now := time.Now().UTC()
	globalTable := &GlobalTable{
		GlobalTableName:   name,
		GlobalTableArn:    s.arnBuilder.GlobalTable(name),
		GlobalTableStatus: "ACTIVE",
		CreationDateTime:  now,
		ReplicationGroup:  replicationGroup,
	}

	if err := s.put(globalTable); err != nil {
		return nil, err
	}

	return globalTable, nil
}

// Put stores a global table under the store lock. Callers that read the
// record before writing it must use Update instead, which holds the lock
// across both halves of the read-modify-write.
func (s *GlobalTableStore) Put(globalTable *GlobalTable) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.put(globalTable)
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
	s.mu.Lock()
	defer s.mu.Unlock()
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

// Delete deletes a global table by name.
func (s *GlobalTableStore) Delete(name string) error {
	return s.BaseStore.Delete(name)
}

// Exists checks if a global table exists.
func (s *GlobalTableStore) Exists(name string) bool {
	return s.BaseStore.Exists(name)
}

// List lists global tables.
func (s *GlobalTableStore) List(marker string, limit int) ([]*GlobalTable, string, error) {
	opts := common.ListOptions{
		Marker:   marker,
		MaxItems: limit,
	}

	result, err := common.ListProto[*storage_dynamodb.GlobalTable](s.BaseStore, opts, func() *storage_dynamodb.GlobalTable { return &storage_dynamodb.GlobalTable{} }, nil)
	if err != nil {
		return nil, "", err
	}

	globalTables := make([]*GlobalTable, len(result.Items))
	for i, pbGlobalTable := range result.Items {
		globalTables[i] = ProtoToGlobalTable(pbGlobalTable)
	}

	if !result.IsTruncated {
		return globalTables, "", nil
	}
	return globalTables, result.NextMarker, nil
}

// ARNBuilder returns the ARN builder for DynamoDB.
func (s *GlobalTableStore) ARNBuilder() *svcarn.DynamoDBBuilder {
	return s.arnBuilder
}
