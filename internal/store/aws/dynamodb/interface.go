package dynamodb

import (
	"context"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TableStoreInterface defines operations for managing DynamoDB tables.
type TableStoreInterface interface {
	Get(name string) (*Table, error)
	Create(params CreateTableParams) (*Table, error)
	Put(table *Table) error
	Update(name string, mutate func(*Table) error) (*Table, error)
	WithTableLock(name string, fn func() error) error
	Exists(name string) bool
	List(marker string, limit int) ([]*Table, string, error)
	Tags() *common.TagStore
	ARNBuilder() *svcarn.DynamoDBBuilder
	SetTimeToLive(name string, ttl *TimeToLiveSpecification) error
	GetTimeToLive(name string) (*TimeToLiveSpecification, error)
	SetPointInTimeRecovery(name string, pitr *PointInTimeRecoveryDescription) error
	GetPointInTimeRecovery(name string) (*PointInTimeRecoveryDescription, error)
	SetResourcePolicyExpected(name string, policy string, expectedRev int) (int, error)
	GetResourcePolicyRevisionId(name string) (int, error)
	GetResourcePolicy(name string) (string, error)
	DeleteResourcePolicyExpected(name string, expectedRev int) (int, error)
	SetKinesisStreamingDestination(name string, destinations []*KinesisDataStreamDestination) error
	SetContributorInsights(name string, enabled bool, mode ContributorInsightsMode) error
	UpdateAutoScalingSettings(name string, merge func(existing *TableReplicaAutoScalingSettings) *TableReplicaAutoScalingSettings) (*TableReplicaAutoScalingSettings, error)
	GetAutoScalingSettings(name string) (*TableReplicaAutoScalingSettings, error)
}

// ItemStoreInterface defines operations for managing DynamoDB items. Item
// writes go through DynamoDBStoreInterface.Update and the DynamoDBTxn write
// methods, which journal and account the writes with the item.
type ItemStoreInterface interface {
	Get(tableName string, key map[string]*AttributeValue) (*Item, error)
	Exists(tableName string, key map[string]*AttributeValue) bool
	List(tableName string, marker string, limit int) ([]*Item, string, error)
	Scan(tableName string, fn func(item *Item) error) error
	ScanWithOptions(tableName string, opts ScanOptions, fn func(item *Item) error) (string, error)
	ScanByPartitionKey(tableName, partitionKeyValue string, fn func(item *Item) error) error
	ScanByPartitionKeyWithTable(tableName string, table *Table, partitionKeyValue string, opts ScanOptions, fn func(item *Item) error) (string, error)
	Count(tableName string) (int64, error)
}

// BackupStoreInterface defines operations for managing DynamoDB backups.
// A backup's identity is its ARN's generated id segment — the name is a
// label, so every keyed operation addresses the ARN.
type BackupStoreInterface interface {
	Get(backupArn string) (*Backup, error)
	Create(backupName, tableName, tableArn string, tableSize int64) (*Backup, error)
	Put(backup *Backup) error
	Delete(backupArn string) error
	List(marker string, limit int, tableName string) ([]*Backup, string, error)
	SaveSnapshot(backupArn string, items []*Item) error
	GetSnapshot(backupArn string) ([]*Item, error)
	DeleteSnapshot(backupArn string) error
}

// GlobalTableStoreInterface defines operations for managing DynamoDB global tables.
type GlobalTableStoreInterface interface {
	Get(name string) (*GlobalTable, error)
	Create(name string, replicationGroup []*Replica) (*GlobalTable, error)
	Update(name string, mutate func(*GlobalTable) error) (*GlobalTable, error)
	Delete(name string) error
	DeleteIfEmpty(name string) (bool, error)
	Exists(name string) bool
	List(marker string, limit int) ([]*GlobalTable, string, error)
}

// ExportStoreInterface defines operations for managing DynamoDB exports.
type ExportStoreInterface interface {
	Get(exportArn string) (*ExportDescription, error)
	Create(tableArn, tableId string, exportFormat ExportFormat) (*ExportDescription, error)
	Put(export *ExportDescription) error
	List(tableArn, marker string, maxItems int) ([]*ExportDescription, string, error)
}

// ImportStoreInterface defines operations for managing DynamoDB imports.
type ImportStoreInterface interface {
	Get(importArn string) (*ImportTableDescription, error)
	Create(tableArn, tableId string) (*ImportTableDescription, error)
	Put(imp *ImportTableDescription) error
	List(tableArn, marker string, maxItems int) ([]*ImportTableDescription, string, error)
}

// DynamoDBStoreInterface defines access to all DynamoDB stores.
type DynamoDBStoreInterface interface {
	Tables() TableStoreInterface
	Items() ItemStoreInterface
	Indexes() *IndexStore
	Backups() BackupStoreInterface
	GlobalTables() GlobalTableStoreInterface
	Exports() ExportStoreInterface
	Imports() ImportStoreInterface
	Streams() *StreamStore
	Contributors() *ContributorStore
	Idempotency() *IdempotencyStore
	Journal() *JournalStore
	Storage() storage.TransactionalStorageWith2PC
	View(ctx context.Context, fn func(txn *DynamoDBTxn) error) error
	Update(ctx context.Context, fn func(txn *DynamoDBTxn) error) error
	TwoPhaseTransaction() storage.TwoPhaseTransaction
	NewTxn(txn storage.Transaction) *DynamoDBTxn
	RecordContributorReads(ctx context.Context, tableName string, keys []map[string]*AttributeValue) error
	RecordContributorQuery(ctx context.Context, tableName string, key map[string]*AttributeValue) error
	FlushContributorWrites(ctx context.Context, events []ContributorWriteEvent)
	FlushTableMetrics(deltas map[string]TableMetricDelta)
}

var (
	_ TableStoreInterface       = (*TableStore)(nil)
	_ ItemStoreInterface        = (*ItemStore)(nil)
	_ BackupStoreInterface      = (*BackupStore)(nil)
	_ GlobalTableStoreInterface = (*GlobalTableStore)(nil)
	_ ExportStoreInterface      = (*ExportStore)(nil)
	_ ImportStoreInterface      = (*ImportStore)(nil)
	_ DynamoDBStoreInterface    = (*DynamoDBStore)(nil)
)
