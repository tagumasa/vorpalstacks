package timestream

import (
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/chunk"
	pb "vorpalstacks/internal/pb/storage/storage_timestream"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func databaseBucketName(region string) string {
	return "timestream_databases-" + region
}

func tableBucketName(region string) string {
	return "timestream_tables-" + region
}



func scheduledQueryBucketName(region string) string {
	return "timestream_scheduled_queries-" + region
}

func scheduledQueryRunBucketName(region string) string {
	return "timestream_scheduled_query_runs-" + region
}

func batchLoadTaskBucketName(region string) string {
	return "timestream_batch_load_tasks-" + region
}

func accountSettingsBucketName(region string) string {
	return "timestream_account_settings-" + region
}

// Store manages Timestream databases.
type Store struct {
	*common.BaseStore
	*common.TagStore
	storage    storage.BasicStorage
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	createMu   sync.Mutex
}

// NewStore creates a new Timestream database store.
func NewStore(store storage.BasicStorage, accountID, region string) *Store {
	return &Store{
		BaseStore:  common.NewBaseStore(store.Bucket(databaseBucketName(region)), "timestream"),
		TagStore:   common.NewTagStoreWithRegion(store, "timestream", region),
		storage:    store,
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		accountID:  accountID,
		region:     region,
	}
}

// GetAccountID returns the AWS account ID.
func (s *Store) GetAccountID() string { return s.accountID }

// GetRegion returns the AWS region.
func (s *Store) GetRegion() string { return s.region }

// CreateDatabase creates a new Timestream database.
func (s *Store) CreateDatabase(name, kmsKeyID string) (*Database, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := name
	if s.Exists(key) {
		return nil, ErrDatabaseAlreadyExists
	}

	now := time.Now().UTC()
	db := &Database{
		DatabaseName:    name,
		ARN:             s.arnBuilder.Timestream().Database(name),
		TableCount:      0,
		KmsKeyId:        kmsKeyID,
		CreationTime:    now,
		LastUpdatedTime: now,
	}

	if err := s.PutProto(key, DatabaseToProto(db)); err != nil {
		return nil, err
	}

	return db, nil
}

// GetDatabase retrieves a database by name.
func (s *Store) GetDatabase(name string) (*Database, error) {
	var pbDb pb.Database
	if err := s.GetProto(name, &pbDb); err != nil {
		return nil, ErrDatabaseNotFound
	}
	return ProtoToDatabase(&pbDb), nil
}

// UpdateDatabase updates an existing database.
func (s *Store) UpdateDatabase(name, kmsKeyID string) (*Database, error) {
	db, err := s.GetDatabase(name)
	if err != nil {
		return nil, err
	}

	db.KmsKeyId = kmsKeyID
	db.LastUpdatedTime = time.Now().UTC()

	if err := s.PutProto(name, DatabaseToProto(db)); err != nil {
		return nil, err
	}

	return db, nil
}

// DeleteDatabase deletes a database by name.
func (s *Store) DeleteDatabase(name string) error {
	if !s.Exists(name) {
		return ErrDatabaseNotFound
	}

	tableStore := common.NewBaseStore(s.storage.Bucket(tableBucketName(s.region)), "timestream")
	prefix := name + "#"
	hasTables := false
	if err := tableStore.ScanPrefix(prefix, func(key string, value []byte) error {
		hasTables = true
		return nil
	}); err != nil {
		return fmt.Errorf("failed to scan tables: %w", err)
	}
	if hasTables {
		return ErrDatabaseNotEmpty
	}

	return s.BaseStore.Delete(name)
}

// ListDatabases returns a list of databases.
func (s *Store) ListDatabases(opts common.ListOptions) (*common.ListResult[Database], error) {
	result, err := common.ListProto[*pb.Database](s.BaseStore, opts, func() *pb.Database { return &pb.Database{} }, nil)
	if err != nil {
		return nil, err
	}

	databases := make([]*Database, 0, len(result.Items))
	for _, pbDb := range result.Items {
		databases = append(databases, ProtoToDatabase(pbDb))
	}

	return &common.ListResult[Database]{
		Items:       databases,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}

type chunkBuffer struct {
	entries  []chunk.TimestreamEntry
	seenKeys map[string]int64
}

type indexedEntry struct {
	entry       *chunk.TimestreamEntry
	originalIdx int
}
