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
	UpdateScheduledQuery(name string, state ScheduledQueryStatus, scheduleConfig *ScheduleConfiguration, notificationConfig *NotificationConfiguration, kmsKeyID string, errorReportConfig *ErrorReportConfiguration, targetConfig *TargetConfiguration) (*ScheduledQuery, error)
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
