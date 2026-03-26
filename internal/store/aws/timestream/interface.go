package timestream

import (
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// DatabaseStoreInterface defines operations for managing Timestream databases.
type DatabaseStoreInterface interface {
	GetAccountID() string
	GetRegion() string
	CreateDatabase(name, kmsKeyID string) (*Database, error)
	GetDatabase(name string) (*Database, error)
	UpdateDatabase(name, kmsKeyID string) (*Database, error)
	DeleteDatabase(name string) error
	ListDatabases(opts common.ListOptions) (*common.ListResult[Database], error)
	Raw() *Store
}

// TableStoreInterface defines operations for managing Timestream tables.
type TableStoreInterface interface {
	CreateTable(databaseName, tableName string, retentionProperties *RetentionProperties, schema *Schema) (*Table, error)
	GetTable(databaseName, tableName string) (*Table, error)
	UpdateTable(databaseName, tableName string, retentionProperties *RetentionProperties, schema *Schema) (*Table, error)
	DeleteTable(databaseName, tableName string) error
	ListTables(databaseName string, opts common.ListOptions) (*common.ListResult[Table], error)
	Raw() *TableStore
}

// RecordStoreInterface defines operations for managing Timestream records.
type RecordStoreInterface interface {
	WriteRecords(databaseName, tableName string, records []Record) ([]RejectedRecord, error)
	QueryRecords(databaseName, tableName string, startTime, endTime time.Time) ([]*StoredRecord, error)
	FlushAllBuffers() error
	Raw() *RecordStore
}

// BatchLoadTaskStoreInterface defines operations for managing Timestream batch load tasks.
type BatchLoadTaskStoreInterface interface {
	CreateBatchLoadTask(taskId, targetDatabaseName, targetTableName string, dataSourceConfig *DataSourceConfiguration, dataModelConfig *DataModelConfiguration, reportConfig *ReportConfiguration, recordVersion int64) (*BatchLoadTaskDescription, error)
	GetBatchLoadTask(taskId string) (*BatchLoadTaskDescription, error)
	UpdateBatchLoadTaskStatus(taskId string, status BatchLoadStatus, errorMessage string) error
	UpdateBatchLoadTaskProgress(taskId string, progress *BatchLoadProgressReport) error
	DeleteBatchLoadTask(taskId string) error
	ListBatchLoadTasks(taskStatus BatchLoadStatus) ([]*BatchLoadTask, error)
	Raw() *BatchLoadTaskStore
}

// ScheduledQueryStoreInterface defines operations for managing scheduled queries.
type ScheduledQueryStoreInterface interface {
	CreateScheduledQuery(name, queryString string, scheduleConfig *ScheduleConfiguration, notificationConfig *NotificationConfiguration, roleARN, kmsKeyID string, errorReportConfig *ErrorReportConfiguration, targetConfig *TargetConfiguration, clientToken string) (*ScheduledQuery, error)
	GetScheduledQuery(name string) (*ScheduledQuery, error)
	UpdateScheduledQuery(name string, state ScheduledQueryState, scheduleConfig *ScheduleConfiguration, notificationConfig *NotificationConfiguration, kmsKeyID string, errorReportConfig *ErrorReportConfiguration, targetConfig *TargetConfiguration) (*ScheduledQuery, error)
	DeleteScheduledQuery(name string) error
	ListScheduledQueries() ([]*ScheduledQuery, error)
	UpdateNextRunTime(name string, nextRunTime time.Time) error
	UpdateLastRun(name string, status string, runTime time.Time) error
	Raw() *ScheduledQueryStore
}

// ScheduledQueryRunStoreInterface defines operations for managing scheduled query runs.
type ScheduledQueryRunStoreInterface interface {
	CreateRun(scheduledQueryARN string, invocationTime, triggerTime time.Time) (*ScheduledQueryRun, error)
	GetRun(arn string) (*ScheduledQueryRun, error)
	UpdateRunStatus(arn string, status ScheduleRunStatus, errStr string, stats *ExecutionStats) error
	ListRuns(scheduledQueryARN string) ([]*ScheduledQueryRun, error)
	Raw() *ScheduledQueryRunStore
}

// AccountSettingsStoreInterface defines operations for managing account settings.
type AccountSettingsStoreInterface interface {
	GetAccountSettings() (*AccountSettings, error)
	UpdateAccountSettings(maxQueryTCU *int64, queryPricingMode, queryComputeType, encryptionConfiguration string) (*AccountSettings, error)
	Raw() *AccountSettingsStore
}

// TimestreamStoresInterface defines access to all Timestream stores.
type TimestreamStoresInterface interface {
	DatabaseStore() DatabaseStoreInterface
	TableStore() TableStoreInterface
	RecordStore() RecordStoreInterface
	BatchLoadTaskStore() BatchLoadTaskStoreInterface
	ScheduledQueryStore() ScheduledQueryStoreInterface
	ScheduledQueryRunStore() ScheduledQueryRunStoreInterface
	AccountSettingsStore() AccountSettingsStoreInterface
}

// TimestreamStores provides access to all Timestream stores.
type TimestreamStores struct {
	databaseStore          *Store
	tableStore             *TableStore
	recordStore            *RecordStore
	batchLoadTaskStore     *BatchLoadTaskStore
	scheduledQueryStore    *ScheduledQueryStore
	scheduledQueryRunStore *ScheduledQueryRunStore
	accountSettingsStore   *AccountSettingsStore
}

// NewTimestreamStores creates a new TimestreamStores with the given stores.
func NewTimestreamStores(database *Store, table *TableStore, record *RecordStore, batchLoad *BatchLoadTaskStore, scheduledQuery *ScheduledQueryStore, scheduledQueryRun *ScheduledQueryRunStore, accountSettings *AccountSettingsStore) *TimestreamStores {
	return &TimestreamStores{
		databaseStore:          database,
		tableStore:             table,
		recordStore:            record,
		batchLoadTaskStore:     batchLoad,
		scheduledQueryStore:    scheduledQuery,
		scheduledQueryRunStore: scheduledQueryRun,
		accountSettingsStore:   accountSettings,
	}
}

// DatabaseStore returns the database store.
func (s *TimestreamStores) DatabaseStore() DatabaseStoreInterface {
	if s.databaseStore == nil {
		return nil
	}
	return s.databaseStore
}

// TableStore returns the table store.
func (s *TimestreamStores) TableStore() TableStoreInterface {
	if s.tableStore == nil {
		return nil
	}
	return s.tableStore
}

// RecordStore returns the record store.
func (s *TimestreamStores) RecordStore() RecordStoreInterface {
	if s.recordStore == nil {
		return nil
	}
	return s.recordStore
}

// BatchLoadTaskStore returns the batch load task store.
func (s *TimestreamStores) BatchLoadTaskStore() BatchLoadTaskStoreInterface {
	if s.batchLoadTaskStore == nil {
		return nil
	}
	return s.batchLoadTaskStore
}

// ScheduledQueryStore returns the scheduled query store.
func (s *TimestreamStores) ScheduledQueryStore() ScheduledQueryStoreInterface {
	if s.scheduledQueryStore == nil {
		return nil
	}
	return s.scheduledQueryStore
}

// ScheduledQueryRunStore returns the scheduled query run store.
func (s *TimestreamStores) ScheduledQueryRunStore() ScheduledQueryRunStoreInterface {
	if s.scheduledQueryRunStore == nil {
		return nil
	}
	return s.scheduledQueryRunStore
}

// AccountSettingsStore returns the account settings store.
func (s *TimestreamStores) AccountSettingsStore() AccountSettingsStoreInterface {
	if s.accountSettingsStore == nil {
		return nil
	}
	return s.accountSettingsStore
}

// Raw returns the underlying database store.
func (s *Store) Raw() *Store {
	return s
}

// Raw returns the underlying table store.
func (s *TableStore) Raw() *TableStore {
	return s
}

// Raw returns the underlying record store.
func (s *RecordStore) Raw() *RecordStore {
	return s
}

// Raw returns the underlying batch load task store.
func (s *BatchLoadTaskStore) Raw() *BatchLoadTaskStore {
	return s
}

// Raw returns the underlying scheduled query store.
func (s *ScheduledQueryStore) Raw() *ScheduledQueryStore {
	return s
}

// Raw returns the underlying scheduled query run store.
func (s *ScheduledQueryRunStore) Raw() *ScheduledQueryRunStore {
	return s
}

// Raw returns the underlying account settings store.
func (s *AccountSettingsStore) Raw() *AccountSettingsStore {
	return s
}
