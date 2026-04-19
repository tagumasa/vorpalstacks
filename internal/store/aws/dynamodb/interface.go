package dynamodb

import (
	"context"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

// TableStoreInterface defines operations for managing DynamoDB tables.
type TableStoreInterface interface {
	Get(name string) (*Table, error)
	Create(
		name string,
		keySchema []*KeySchemaElement,
		attributeDefinitions []*AttributeDefinition,
		billingMode BillingMode,
		provisionedThroughput *ProvisionedThroughput,
		gsi []*GlobalSecondaryIndex,
		lsi []*LocalSecondaryIndex,
		streamSpec *StreamSpecification,
		tags []types.Tag,
		deletionProtectionEnabled bool,
	) (*Table, error)
	Put(table *Table) error
	Delete(name string) error
	Exists(name string) bool
	List(marker string, limit int) ([]*Table, string, error)
	UpdateItemCount(name string, delta int64) error
	UpdateTableSize(name string, delta int64) error
	Tags() *common.TagStore
	ARNBuilder() *svcarn.DynamoDBBuilder
	GetPartitionKey(table *Table) string
	GetSortKey(table *Table) string
	SetTimeToLive(name string, ttl *TimeToLiveSpecification) error
	GetTimeToLive(name string) (*TimeToLiveSpecification, error)
	SetPointInTimeRecovery(name string, pitr *PointInTimeRecoveryDescription) error
	GetPointInTimeRecovery(name string) (*PointInTimeRecoveryDescription, error)
	SetResourcePolicy(name string, policy string) error
	GetResourcePolicy(name string) (string, error)
	DeleteResourcePolicy(name string) error
	SetKinesisStreamingDestination(name string, destinations []*KinesisDataStreamDestination) error
	SetContributorInsights(name string, enabled bool) error
}

// ItemStoreInterface defines operations for managing DynamoDB items.
type ItemStoreInterface interface {
	Get(tableName string, key map[string]*AttributeValue) (*Item, error)
	Put(tableName string, key map[string]*AttributeValue, attributes map[string]*AttributeValue) (*Item, error)
	Delete(tableName string, key map[string]*AttributeValue) error
	Exists(tableName string, key map[string]*AttributeValue) bool
	List(tableName string, marker string, limit int) ([]*Item, string, error)
	Scan(tableName string, fn func(item *Item) error) error
	ScanByPartitionKey(tableName, partitionKeyValue string, fn func(item *Item) error) error
	ScanByPartitionKeyWithTable(tableName string, table *Table, partitionKeyValue string, fn func(item *Item) error) error
	Count(tableName string) (int64, error)
	DeleteAllForTable(tableName string) error
}

// BackupStoreInterface defines operations for managing DynamoDB backups.
type BackupStoreInterface interface {
	Get(backupArn string) (*Backup, error)
	GetByName(backupName string) (*Backup, error)
	Create(backupName, tableName, tableArn string, tableSize int64) (*Backup, error)
	Put(backup *Backup) error
	Delete(backupName string) error
	Exists(backupName string) bool
	List(marker string, limit int, tableName string) ([]*Backup, string, error)
	ARNBuilder() *svcarn.DynamoDBBuilder
}

// GlobalTableStoreInterface defines operations for managing DynamoDB global tables.
type GlobalTableStoreInterface interface {
	Get(name string) (*GlobalTable, error)
	Create(name string, replicationGroup []*Replica) (*GlobalTable, error)
	Put(globalTable *GlobalTable) error
	Delete(name string) error
	Exists(name string) bool
	List(marker string, limit int) ([]*GlobalTable, string, error)
	ARNBuilder() *svcarn.DynamoDBBuilder
}

// ExportStoreInterface defines operations for managing DynamoDB exports.
type ExportStoreInterface interface {
	Get(exportArn string) (*ExportDescription, error)
	Create(tableArn, tableId, exportFormat string) (*ExportDescription, error)
	Put(export *ExportDescription) error
	List(tableArn string) ([]*ExportDescription, error)
}

// ImportStoreInterface defines operations for managing DynamoDB imports.
type ImportStoreInterface interface {
	Get(importArn string) (*ImportTableDescription, error)
	Create(tableArn, tableId string) (*ImportTableDescription, error)
	Put(imp *ImportTableDescription) error
	List(tableArn string) ([]*ImportTableDescription, error)
}

// DynamoDBTxnInterface defines operations for DynamoDB transactions.
type DynamoDBTxnInterface interface {
	GetTable(name string) (*Table, error)
	PutTable(table *Table) error
	GetItem(tableName string, key map[string]*AttributeValue) (*Item, error)
	PutItem(tableName string, key map[string]*AttributeValue, attributes map[string]*AttributeValue) error
	DeleteItem(tableName string, key map[string]*AttributeValue) error
	ItemExists(tableName string, key map[string]*AttributeValue) (bool, error)
	UpdateItemCount(tableName string, delta int64) error
	UpdateTableSize(tableName string, delta int64) error
	PutIndexEntries(tableName string, item *Item) error
	DeleteIndexEntries(tableName string, item *Item) error
	QueryByGSI(tableName, indexName, hashKeyValue string, opts IndexQueryOptions) ([]*Item, error)
	QueryByLSI(tableName, indexName, hashKeyValue string, opts IndexQueryOptions) ([]*Item, error)
	Scan(tableName string, fn func(item *Item) error) error
	ScanByPartitionKey(tableName, partitionKeyValue string, fn func(item *Item) error) error
}

// DynamoDBStoreInterface defines access to all DynamoDB stores.
type DynamoDBStoreInterface interface {
	Tables() TableStoreInterface
	Items() ItemStoreInterface
	Backups() BackupStoreInterface
	GlobalTables() GlobalTableStoreInterface
	Exports() ExportStoreInterface
	Imports() ImportStoreInterface
	Storage() storage.TransactionalStorageWith2PC
	View(ctx context.Context, fn func(txn *DynamoDBTxn) error) error
	Update(ctx context.Context, fn func(txn *DynamoDBTxn) error) error
	TwoPhaseTransaction() storage.TwoPhaseTransaction
}

var (
	_ TableStoreInterface       = (*TableStore)(nil)
	_ ItemStoreInterface        = (*ItemStore)(nil)
	_ BackupStoreInterface      = (*BackupStore)(nil)
	_ GlobalTableStoreInterface = (*GlobalTableStore)(nil)
	_ ExportStoreInterface      = (*ExportStore)(nil)
	_ ImportStoreInterface      = (*ImportStore)(nil)
	_ DynamoDBTxnInterface      = (*DynamoDBTxn)(nil)
	_ DynamoDBStoreInterface    = (*DynamoDBStore)(nil)
)
