// Package timestream provides Timestream storage functionality for vorpalstacks.
package timestream

import "time"

// Database represents an Amazon Timestream database.
type Database struct {
	DatabaseName    string    `json:"databaseName"`
	ARN             string    `json:"arn"`
	TableCount      int64     `json:"tableCount"`
	KmsKeyId        string    `json:"kmsKeyId,omitempty"`
	CreationTime    time.Time `json:"creationTime"`
	LastUpdatedTime time.Time `json:"lastUpdatedTime"`
}

// Table represents an Amazon Timestream table.
type Table struct {
	TableName           string               `json:"tableName"`
	DatabaseName        string               `json:"databaseName"`
	ARN                 string               `json:"arn"`
	TableStatus         TableStatus          `json:"tableStatus"`
	RetentionProperties *RetentionProperties `json:"retentionProperties,omitempty"`
	Schema              *Schema              `json:"schema,omitempty"`
	CreationTime        time.Time            `json:"creationTime"`
	LastUpdatedTime     time.Time            `json:"lastUpdatedTime"`
}

// TableStatus represents the status of a Timestream table.
type TableStatus string

// TableStatus constants define the possible statuses of a Timestream table.
const (
	TableStatusActive    TableStatus = "ACTIVE"
	TableStatusDeleting  TableStatus = "DELETING"
	TableStatusRestoring TableStatus = "RESTORING"
)

// RetentionProperties defines the retention period for data in memory and magnetic storage.
type RetentionProperties struct {
	MemoryStoreRetentionPeriodInHours  int64 `json:"memoryStoreRetentionPeriodInHours"`
	MagneticStoreRetentionPeriodInDays int64 `json:"magneticStoreRetentionPeriodInDays"`
}

// Schema defines the schema for a Timestream table.
type Schema struct {
	CompositePartitionKey []PartitionKey `json:"compositePartitionKey,omitempty"`
}

// PartitionKey defines a partition key for a Timestream table.
type PartitionKey struct {
	Type                PartitionKeyType    `json:"type"`
	Name                string              `json:"name,omitempty"`
	EnforcementInRecord EnforcementInRecord `json:"enforcementInRecord,omitempty"`
}

// PartitionKeyType defines the type of partition key.
type PartitionKeyType string

const (
	// PartitionKeyTypeDimension represents a dimension-based partition key.
	PartitionKeyTypeDimension PartitionKeyType = "DIMENSION"
	// PartitionKeyTypeMeasure represents a measure-based partition key.
	PartitionKeyTypeMeasure PartitionKeyType = "MEASURE"
)

// EnforcementInRecord defines enforcement requirements for a partition key.
type EnforcementInRecord string

const (
	// EnforcementInRecordRequired indicates the partition key is required in every record.
	EnforcementInRecordRequired EnforcementInRecord = "REQUIRED"
	// EnforcementInRecordOptional indicates the partition key is optional in records.
	EnforcementInRecordOptional EnforcementInRecord = "OPTIONAL"
)

// Record represents a Timestream record containing time series data.
type Record struct {
	Dimensions       []Dimension      `json:"dimensions,omitempty"`
	MeasureName      string           `json:"measureName"`
	MeasureValue     string           `json:"measureValue,omitempty"`
	MeasureValueType MeasureValueType `json:"measureValueType,omitempty"`
	MeasureValues    []MeasureValue   `json:"measureValues,omitempty"`
	Time             string           `json:"time"`
	TimeUnit         TimeUnit         `json:"timeUnit,omitempty"`
	Version          int64            `json:"version,omitempty"`
}

// Dimension represents a dimension in a Timestream record.
type Dimension struct {
	Name               string             `json:"name"`
	Value              string             `json:"value"`
	DimensionValueType DimensionValueType `json:"dimensionValueType,omitempty"`
}

// DimensionValueType defines the type of dimension value.
type DimensionValueType string

const (
	// DimensionValueTypeVarchar represents a VARCHAR dimension value type.
	DimensionValueTypeVarchar DimensionValueType = "VARCHAR"
)

// MeasureValueType defines the type of measure value.
type MeasureValueType string

const (
	// MeasureValueTypeDouble represents a DOUBLE measure value type.
	MeasureValueTypeDouble MeasureValueType = "DOUBLE"
	// MeasureValueTypeBigint represents a BIGINT measure value type.
	MeasureValueTypeBigint MeasureValueType = "BIGINT"
	// MeasureValueTypeVarchar represents a VARCHAR measure value type.
	MeasureValueTypeVarchar MeasureValueType = "VARCHAR"
	// MeasureValueTypeBoolean represents a BOOLEAN measure value type.
	MeasureValueTypeBoolean MeasureValueType = "BOOLEAN"
	// MeasureValueTypeTimestamp represents a TIMESTAMP measure value type.
	MeasureValueTypeTimestamp MeasureValueType = "TIMESTAMP"
	// MeasureValueTypeMulti represents a MULTI measure value type.
	MeasureValueTypeMulti MeasureValueType = "MULTI"
)

// MeasureValue represents a measure value in a multi-measure record.
type MeasureValue struct {
	Name  string           `json:"name"`
	Value string           `json:"value"`
	Type  MeasureValueType `json:"type,omitempty"`
}

// TimeUnit defines the unit of time for Timestream records.
type TimeUnit string

const (
	// TimeUnitMilliseconds represents milliseconds time unit.
	TimeUnitMilliseconds TimeUnit = "MILLISECONDS"
	// TimeUnitSeconds represents seconds time unit.
	TimeUnitSeconds TimeUnit = "SECONDS"
	// TimeUnitMicroseconds represents microseconds time unit.
	TimeUnitMicroseconds TimeUnit = "MICROSECONDS"
	// TimeUnitNanoseconds represents nanoseconds time unit.
	TimeUnitNanoseconds TimeUnit = "NANOSECONDS"
)

// Endpoint represents a Timestream endpoint for database connectivity.
type Endpoint struct {
	Address              string `json:"address"`
	CachePeriodInMinutes int64  `json:"cachePeriodInMinutes"`
}

// RejectedRecord represents a record that was rejected during ingestion.
type RejectedRecord struct {
	RecordIndex     int64  `json:"recordIndex"`
	Reason          string `json:"reason"`
	ExistingVersion int64  `json:"existingVersion,omitempty"`
}

// ScheduledQueryState represents the state of a scheduled query.
type ScheduledQueryState string

const (
	// ScheduledQueryStateEnabled indicates the scheduled query is enabled.
	ScheduledQueryStateEnabled ScheduledQueryState = "ENABLED"
	// ScheduledQueryStateDisabled indicates the scheduled query is disabled.
	ScheduledQueryStateDisabled ScheduledQueryState = "DISABLED"
)

// ScheduledQueryStatus represents the status of a scheduled query.
type ScheduledQueryStatus string

const (
	// ScheduledQueryStatusEnabled indicates the scheduled query is enabled.
	ScheduledQueryStatusEnabled ScheduledQueryStatus = "ENABLED"
	// ScheduledQueryStatusDisabled indicates the scheduled query is disabled.
	ScheduledQueryStatusDisabled ScheduledQueryStatus = "DISABLED"
	// ScheduledQueryStatusCreating indicates the scheduled query is being created.
	ScheduledQueryStatusCreating ScheduledQueryStatus = "CREATING"
	// ScheduledQueryStatusUpdating indicates the scheduled query is being updated.
	ScheduledQueryStatusUpdating ScheduledQueryStatus = "UPDATING"
	// ScheduledQueryStatusDeleting indicates the scheduled query is being deleted.
	ScheduledQueryStatusDeleting ScheduledQueryStatus = "DELETING"
	// ScheduledQueryStatusCreateFailed indicates the scheduled query creation failed.
	ScheduledQueryStatusCreateFailed ScheduledQueryStatus = "CREATE_FAILED"
	// ScheduledQueryStatusUpdateFailed indicates the scheduled query update failed.
	ScheduledQueryStatusUpdateFailed ScheduledQueryStatus = "UPDATE_FAILED"
	// ScheduledQueryStatusDeleteFailed indicates the scheduled query deletion failed.
	ScheduledQueryStatusDeleteFailed ScheduledQueryStatus = "DELETE_FAILED"
)

// ScheduledQuery represents a scheduled query in Timestream.
type ScheduledQuery struct {
	ARN                            string                     `json:"arn"`
	Name                           string                     `json:"name"`
	QueryString                    string                     `json:"queryString"`
	ScheduleConfiguration          *ScheduleConfiguration     `json:"scheduleConfiguration"`
	NotificationConfiguration      *NotificationConfiguration `json:"notificationConfiguration,omitempty"`
	ScheduledQueryExecutionRoleARN string                     `json:"scheduledQueryExecutionRoleArn,omitempty"`
	KmsKeyID                       string                     `json:"kmsKeyId,omitempty"`
	ErrorReportConfiguration       *ErrorReportConfiguration  `json:"errorReportConfiguration,omitempty"`
	TargetConfiguration            *TargetConfiguration       `json:"targetConfiguration,omitempty"`
	ClientToken                    string                     `json:"clientToken,omitempty"`
	ScheduledQueryStatus           ScheduledQueryStatus       `json:"scheduledQueryStatus"`
	CreationTime                   time.Time                  `json:"creationTime"`
	LastRunStatus                  string                     `json:"lastRunStatus,omitempty"`
	PreviousRunTime                time.Time                  `json:"previousRunTime,omitempty"`
	NextRunTime                    time.Time                  `json:"nextRunTime,omitempty"`
}

// ScheduleConfiguration defines the schedule for a scheduled query.
type ScheduleConfiguration struct {
	ScheduleExpression string `json:"scheduleExpression"`
}

// NotificationConfiguration defines the notification settings for a scheduled query.
type NotificationConfiguration struct {
	SNSConfiguration *SNSConfiguration `json:"snsConfiguration,omitempty"`
}

// SNSConfiguration defines the SNS topic configuration for notifications.
type SNSConfiguration struct {
	TopicARN string `json:"topicArn"`
}

// ErrorReportConfiguration defines where to store error reports for failed query executions.
type ErrorReportConfiguration struct {
	S3Configuration *S3ErrorReportConfiguration `json:"s3Configuration,omitempty"`
}

// S3ErrorReportConfiguration defines the S3 location for error reports.
type S3ErrorReportConfiguration struct {
	BucketName      string `json:"bucketName"`
	ObjectKeyPrefix string `json:"objectKeyPrefix,omitempty"`
}

// TargetConfiguration defines the target configuration for a scheduled query.
type TargetConfiguration struct {
	TimestreamConfiguration *TimestreamConfiguration `json:"timestreamConfiguration"`
}

// TimestreamConfiguration defines the target Timestream table configuration.
type TimestreamConfiguration struct {
	DatabaseName         string                `json:"databaseName"`
	TableName            string                `json:"tableName"`
	TimeColumn           string                `json:"timeColumn,omitempty"`
	DimensionMappings    []DimensionMapping    `json:"dimensionMappings,omitempty"`
	MultiMeasureMappings *MultiMeasureMappings `json:"multiMeasureMappings,omitempty"`
	MixedMeasureMappings []MixedMeasureMapping `json:"mixedMeasureMappings,omitempty"`
}

// DimensionMapping defines the mapping from source columns to dimensions.
type DimensionMapping struct {
	SourceColumn      *SourceColumn      `json:"sourceColumn,omitempty"`
	DestinationColumn *DestinationColumn `json:"destinationColumn,omitempty"`
}

// SourceColumn defines the source column name for mapping.
type SourceColumn struct {
	Name string `json:"name"`
}

// DestinationColumn defines the destination column name for mapping.
type DestinationColumn struct {
	Name string `json:"name"`
}

// MultiMeasureMappings defines mappings for multi-measure records.
type MultiMeasureMappings struct {
	MultiMeasureAttributeMappings []MultiMeasureAttributeMapping `json:"multiMeasureAttributeMappings"`
}

// MultiMeasureAttributeMapping defines the mapping for a multi-measure attribute.
type MultiMeasureAttributeMapping struct {
	SourceColumn                    *SourceColumn    `json:"sourceColumn"`
	TargetMultiMeasureAttributeName string           `json:"targetMultiMeasureAttributeName"`
	MeasureValueMeasureValueType    MeasureValueType `json:"measureValueType"`
}

// MixedMeasureMapping defines the mapping for mixed measure values.
type MixedMeasureMapping struct {
	MeasureName                   string                         `json:"measureName,omitempty"`
	SourceColumn                  string                         `json:"sourceColumn,omitempty"`
	TargetMeasureName             string                         `json:"targetMeasureName,omitempty"`
	MeasureValueMeasureValueType  MeasureValueType               `json:"measureValueType"`
	MultiMeasureAttributeMappings []MultiMeasureAttributeMapping `json:"multiMeasureAttributeMappings,omitempty"`
}

// ScheduleRunStatus represents the execution status of a scheduled query run.
type ScheduleRunStatus string

const (
	// ScheduleRunStatusPending indicates the scheduled query run is pending.
	ScheduleRunStatusPending ScheduleRunStatus = "PENDING"
	// ScheduleRunStatusRunning indicates the scheduled query run is in progress.
	ScheduleRunStatusRunning ScheduleRunStatus = "RUNNING"
	// ScheduleRunStatusSucceeded indicates the scheduled query run succeeded.
	ScheduleRunStatusSucceeded ScheduleRunStatus = "SUCCEEDED"
	// ScheduleRunStatusFailed indicates the scheduled query run failed.
	ScheduleRunStatusFailed ScheduleRunStatus = "FAILED"
	// ScheduleRunStatusCancelled indicates the scheduled query run was cancelled.
	ScheduleRunStatusCancelled ScheduleRunStatus = "CANCELLED"
)

// ScheduledQueryRun represents the execution of a scheduled query.
type ScheduledQueryRun struct {
	ARN               string            `json:"arn"`
	ScheduledQueryARN string            `json:"scheduledQueryArn"`
	InvocationTime    time.Time         `json:"invocationTime"`
	TriggerTime       time.Time         `json:"triggerTime"`
	RunStatus         ScheduleRunStatus `json:"runStatus"`
	ExecutionStats    *ExecutionStats   `json:"executionStats,omitempty"`
	Error             string            `json:"error,omitempty"`
	CompletionTime    time.Time         `json:"completionTime,omitempty"`
}

// ExecutionStats contains statistics about a scheduled query execution.
type ExecutionStats struct {
	TotalBytesProcessed    int64 `json:"totalBytesProcessed,omitempty"`
	DataWrites             int64 `json:"dataWrites,omitempty"`
	TotalRecordsProcessed  int64 `json:"totalRecordsProcessed,omitempty"`
	BytesMetered           int64 `json:"bytesMetered,omitempty"`
	QueryResultRows        int64 `json:"queryResultRows,omitempty"`
	CumulativeBytesScanned int64 `json:"cumulativeBytesScanned,omitempty"`
	ExecutionTimeInMillis  int64 `json:"executionTimeInMillis,omitempty"`
}

// BatchLoadStatus represents the status of a batch load task.
type BatchLoadStatus string

const (
	// BatchLoadStatusCreated indicates the batch load task has been created.
	BatchLoadStatusCreated BatchLoadStatus = "CREATED"
	// BatchLoadStatusInProgress indicates the batch load task is in progress.
	BatchLoadStatusInProgress BatchLoadStatus = "IN_PROGRESS"
	// BatchLoadStatusFailed indicates the batch load task failed.
	BatchLoadStatusFailed BatchLoadStatus = "FAILED"
	// BatchLoadStatusSucceeded indicates the batch load task succeeded.
	BatchLoadStatusSucceeded BatchLoadStatus = "SUCCEEDED"
	// BatchLoadStatusProgressStopped indicates the batch load task progress has stopped.
	BatchLoadStatusProgressStopped BatchLoadStatus = "PROGRESS_STOPPED"
	// BatchLoadStatusPendingResume indicates the batch load task is pending resume.
	BatchLoadStatusPendingResume BatchLoadStatus = "PENDING_RESUME"
)

// BatchLoadDataFormat defines the data format for batch loading.
type BatchLoadDataFormat string

const (
	// BatchLoadDataFormatCsv indicates CSV data format for batch loading.
	BatchLoadDataFormatCsv BatchLoadDataFormat = "CSV"
)

// BatchLoadTask represents a batch load task for importing data into Timestream.
type BatchLoadTask struct {
	TaskId          string          `json:"taskId"`
	DatabaseName    string          `json:"databaseName,omitempty"`
	TableName       string          `json:"tableName,omitempty"`
	TaskStatus      BatchLoadStatus `json:"taskStatus"`
	CreationTime    time.Time       `json:"creationTime"`
	LastUpdatedTime time.Time       `json:"lastUpdatedTime,omitempty"`
	ResumableUntil  time.Time       `json:"resumableUntil,omitempty"`
}

// BatchLoadTaskDescription represents the detailed description of a batch load task.
type BatchLoadTaskDescription struct {
	TaskId                  string                   `json:"taskId"`
	TargetDatabaseName      string                   `json:"targetDatabaseName,omitempty"`
	TargetTableName         string                   `json:"targetTableName,omitempty"`
	TaskStatus              BatchLoadStatus          `json:"taskStatus"`
	CreationTime            time.Time                `json:"creationTime"`
	LastUpdatedTime         time.Time                `json:"lastUpdatedTime,omitempty"`
	ResumableUntil          time.Time                `json:"resumableUntil,omitempty"`
	ErrorMessage            string                   `json:"errorMessage,omitempty"`
	RecordVersion           int64                    `json:"recordVersion,omitempty"`
	DataSourceConfiguration *DataSourceConfiguration `json:"dataSourceConfiguration,omitempty"`
	DataModelConfiguration  *DataModelConfiguration  `json:"dataModelConfiguration,omitempty"`
	ReportConfiguration     *ReportConfiguration     `json:"reportConfiguration,omitempty"`
	ProgressReport          *BatchLoadProgressReport `json:"progressReport,omitempty"`
}

// BatchLoadProgressReport contains progress information for a batch load task.
type BatchLoadProgressReport struct {
	BytesMetered            int64 `json:"bytesMetered,omitempty"`
	FileFailures            int64 `json:"fileFailures,omitempty"`
	ParseFailures           int64 `json:"parseFailures,omitempty"`
	RecordIngestionFailures int64 `json:"recordIngestionFailures,omitempty"`
	RecordsIngested         int64 `json:"recordsIngested,omitempty"`
	RecordsProcessed        int64 `json:"recordsProcessed,omitempty"`
}

// DataSourceConfiguration defines the data source configuration for batch loading.
type DataSourceConfiguration struct {
	DataFormat                BatchLoadDataFormat        `json:"dataFormat"`
	DataSourceS3Configuration *DataSourceS3Configuration `json:"dataSourceS3Configuration,omitempty"`
	CsvConfiguration          *CsvConfiguration          `json:"csvConfiguration,omitempty"`
}

// DataSourceS3Configuration defines the S3 location for the data source.
type DataSourceS3Configuration struct {
	BucketName      string `json:"bucketName"`
	ObjectKeyPrefix string `json:"objectKeyPrefix,omitempty"`
}

// CsvConfiguration defines the CSV parsing configuration for batch loading.
type CsvConfiguration struct {
	ColumnSeparator string `json:"columnSeparator,omitempty"`
	EscapeChar      string `json:"escapeChar,omitempty"`
	NullValue       string `json:"nullValue,omitempty"`
	QuoteChar       string `json:"quoteChar,omitempty"`
	TrimWhiteSpace  bool   `json:"trimWhiteSpace,omitempty"`
}

// DataModelConfiguration defines the data model configuration for batch loading.
type DataModelConfiguration struct {
	DataModel                *DataModel                `json:"dataModel,omitempty"`
	DataModelS3Configuration *DataModelS3Configuration `json:"dataModelS3Configuration,omitempty"`
}

// DataModel defines the data model for batch loading.
type DataModel struct {
	DimensionMappings    []DimensionMapping    `json:"dimensionMappings,omitempty"`
	MeasureNameColumn    string                `json:"measureNameColumn,omitempty"`
	MixedMeasureMappings []MixedMeasureMapping `json:"mixedMeasureMappings,omitempty"`
	MultiMeasureMappings *MultiMeasureMappings `json:"multiMeasureMappings,omitempty"`
	TimeColumn           string                `json:"timeColumn,omitempty"`
	TimeUnit             TimeUnit              `json:"timeUnit,omitempty"`
}

// DataModelS3Configuration defines the S3 location for the data model.
type DataModelS3Configuration struct {
	BucketName string `json:"bucketName,omitempty"`
	ObjectKey  string `json:"objectKey,omitempty"`
}

// ReportConfiguration defines the configuration for batch load reports.
type ReportConfiguration struct {
	ReportS3Configuration *ReportS3Configuration `json:"reportS3Configuration,omitempty"`
}

// ReportS3Configuration defines the S3 location for batch load reports.
type ReportS3Configuration struct {
	BucketName       string             `json:"bucketName"`
	EncryptionOption S3EncryptionOption `json:"encryptionOption,omitempty"`
	KmsKeyId         string             `json:"kmsKeyId,omitempty"`
	ObjectKeyPrefix  string             `json:"objectKeyPrefix,omitempty"`
}

// S3EncryptionOption defines the S3 encryption option.
type S3EncryptionOption string

const (
	// S3EncryptionOptionSSES3 indicates S3-managed encryption.
	S3EncryptionOptionSSES3 S3EncryptionOption = "SSE_S3"
	// S3EncryptionOptionSSEKMS indicates KMS-managed encryption.
	S3EncryptionOptionSSEKMS S3EncryptionOption = "SSE_KMS"
)

// AccountSettings represents the account settings for Timestream.
type AccountSettings struct {
	MaxQueryTCU             int64  `json:"maxQueryTCU,omitempty"`
	QueryPricingMode        string `json:"queryPricingMode,omitempty"`
	EncryptionConfiguration string `json:"encryptionConfiguration,omitempty"`
	QueryComputeType        string `json:"queryComputeType,omitempty"`
}

const (
	// QueryPricingModeBytesScanned indicates pricing based on bytes scanned.
	QueryPricingModeBytesScanned string = "BYTES_SCANNED"
	// QueryPricingModeComputeUnits indicates pricing based on compute units.
	QueryPricingModeComputeUnits string = "COMPUTE_UNITS"
	// QueryComputeTypeProvisioned indicates provisioned compute capacity.
	QueryComputeTypeProvisioned string = "PROVISIONED"
	// QueryComputeTypeOnDemand indicates on-demand compute capacity.
	QueryComputeTypeOnDemand string = "ON_DEMAND"
	// EncryptionConfigurationSSEKMS indicates KMS encryption for account settings.
	EncryptionConfigurationSSEKMS string = "SSE_KMS"
	// EncryptionConfigurationSSES3 indicates S3 encryption for account settings.
	EncryptionConfigurationSSES3 string = "SSE_S3"
)
