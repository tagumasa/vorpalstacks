package timestream

import (
	"log"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_timestream"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TableStore manages Timestream tables.
type TableStore struct {
	*common.BaseStore
	arnBuilder    *svcarn.ARNBuilder
	accountID     string
	region        string
	databaseStore *Store
	createMu      sync.Mutex
}

// NewTableStore creates a new table store.
func NewTableStore(store storage.BasicStorage, databaseStore *Store, accountID, region string) *TableStore {
	return &TableStore{
		BaseStore:     common.NewBaseStore(store.Bucket(tableBucketName(region)), "timestream"),
		arnBuilder:    svcarn.NewARNBuilder(accountID, region),
		accountID:     accountID,
		region:        region,
		databaseStore: databaseStore,
	}
}

// CreateTable creates a new table in a database.
func (s *TableStore) CreateTable(databaseName, tableName string, retentionProperties *RetentionProperties, schema *Schema) (*Table, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := databaseName + "#" + tableName
	if s.Exists(key) {
		return nil, ErrTableAlreadyExists
	}

	if !s.databaseStore.Exists(databaseName) {
		return nil, ErrDatabaseNotFound
	}

	now := time.Now().UTC()
	table := &Table{
		TableName:       tableName,
		DatabaseName:    databaseName,
		ARN:             s.arnBuilder.Timestream().Table(databaseName, tableName),
		TableStatus:     TableStatusActive,
		CreationTime:    now,
		LastUpdatedTime: now,
	}

	if retentionProperties != nil {
		table.RetentionProperties = retentionProperties
	} else {
		table.RetentionProperties = &RetentionProperties{
			MemoryStoreRetentionPeriodInHours:  6,
			MagneticStoreRetentionPeriodInDays: 7300,
		}
	}

	if schema != nil {
		table.Schema = schema
	}

	if err := s.PutProto(key, TableToProto(table)); err != nil {
		return nil, err
	}

	if db, err := s.databaseStore.GetDatabase(databaseName); err == nil {
		db.TableCount++
		if err := s.databaseStore.PutProto(databaseName, DatabaseToProto(db)); err != nil {
			log.Printf("Failed to update database table count: %v", err)
		}
	}

	return table, nil
}

// GetTable retrieves a table by database and table name.
func (s *TableStore) GetTable(databaseName, tableName string) (*Table, error) {
	key := databaseName + "#" + tableName
	var pbTable pb.Table
	if err := s.GetProto(key, &pbTable); err != nil {
		return nil, ErrTableNotFound
	}
	return ProtoToTable(&pbTable), nil
}

// UpdateTable updates an existing table.
func (s *TableStore) UpdateTable(databaseName, tableName string, retentionProperties *RetentionProperties, schema *Schema) (*Table, error) {
	table, err := s.GetTable(databaseName, tableName)
	if err != nil {
		return nil, err
	}

	if retentionProperties != nil {
		table.RetentionProperties = retentionProperties
	}
	if schema != nil {
		table.Schema = schema
	}
	table.LastUpdatedTime = time.Now().UTC()

	key := databaseName + "#" + tableName
	if err := s.PutProto(key, TableToProto(table)); err != nil {
		return nil, err
	}

	return table, nil
}

// DeleteTable deletes a table by database and table name.
func (s *TableStore) DeleteTable(databaseName, tableName string) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := databaseName + "#" + tableName
	if !s.Exists(key) {
		return ErrTableNotFound
	}
	if err := s.BaseStore.Delete(key); err != nil {
		return err
	}

	if db, err := s.databaseStore.GetDatabase(databaseName); err == nil {
		if db.TableCount > 0 {
			db.TableCount--
			if err := s.databaseStore.PutProto(databaseName, DatabaseToProto(db)); err != nil {
				log.Printf("Failed to update database table count: %v", err)
			}
		}
	}

	return nil
}

// ListTables returns tables in a database.
func (s *TableStore) ListTables(databaseName string, opts common.ListOptions) (*common.ListResult[Table], error) {
	if databaseName != "" {
		opts.Prefix = databaseName + "#"
	}
	result, err := common.ListProto[*pb.Table](s.BaseStore, opts, func() *pb.Table { return &pb.Table{} }, nil)
	if err != nil {
		return nil, err
	}

	tables := make([]*Table, 0, len(result.Items))
	for _, pbTable := range result.Items {
		tables = append(tables, ProtoToTable(pbTable))
	}

	return &common.ListResult[Table]{
		Items:       tables,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}
