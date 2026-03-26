// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
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
// Get retrieves a global table by name.
func (s *GlobalTableStore) Get(name string) (*GlobalTable, error) {
	var globalTable GlobalTable
	if err := s.BaseStore.Get(name, &globalTable); err != nil {
		return nil, err
	}
	return &globalTable, nil
}

// Create creates a new global table.
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

	if err := s.BaseStore.Put(name, globalTable); err != nil {
		return nil, err
	}

	return globalTable, nil
}

// Put stores a global table.
func (s *GlobalTableStore) Put(globalTable *GlobalTable) error {
	return s.BaseStore.Put(globalTable.GlobalTableName, globalTable)
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

	result, err := common.List[GlobalTable](s.BaseStore, opts, nil)
	if err != nil {
		return nil, "", err
	}

	if !result.IsTruncated {
		return result.Items, "", nil
	}
	return result.Items, result.NextMarker, nil
}

// ARNBuilder returns the ARN builder for DynamoDB.
func (s *GlobalTableStore) ARNBuilder() *svcarn.DynamoDBBuilder {
	return s.arnBuilder
}
