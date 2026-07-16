// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

// TableStatus represents the status of a DynamoDB table.
type TableStatus string

// TableStatus constants define the possible statuses of a DynamoDB table.
const (
	TableStatusCreating TableStatus = "CREATING"
	TableStatusActive   TableStatus = "ACTIVE"
	TableStatusUpdating TableStatus = "UPDATING"
	TableStatusDeleting TableStatus = "DELETING"
	TableStatusArchived TableStatus = "ARCHIVED"
)

// BillingMode represents the billing mode for a DynamoDB table.
type BillingMode string

// BillingMode constants define the supported billing modes.
const (
	BillingModeProvisioned   BillingMode = "PROVISIONED"
	BillingModePayPerRequest BillingMode = "PAY_PER_REQUEST"
)

// KeyType represents the type of a key in a DynamoDB key schema.
type KeyType string

// KeyType constants define the supported key types.
const (
	KeyTypeHash  KeyType = "HASH"
	KeyTypeRange KeyType = "RANGE"
)

// ScalarAttributeType represents the data type of a scalar attribute.
type ScalarAttributeType string

// ScalarAttributeType constants define the supported scalar attribute types.
const (
	ScalarAttributeTypeS ScalarAttributeType = "S"
	ScalarAttributeTypeN ScalarAttributeType = "N"
	ScalarAttributeTypeB ScalarAttributeType = "B"
)

// StreamViewType represents the type of view for DynamoDB streams.
type StreamViewType string

// StreamViewType constants define the supported stream view types.
const (
	StreamViewTypeNewImage        StreamViewType = "NEW_IMAGE"
	StreamViewTypeOldImage        StreamViewType = "OLD_IMAGE"
	StreamViewTypeNewAndOldImages StreamViewType = "NEW_AND_OLD_IMAGES"
	StreamViewTypeKeysOnly        StreamViewType = "KEYS_ONLY"
)

// SSEType represents the server-side encryption type.
type SSEType string

// SSEType constants define the supported server-side encryption types.
const (
	SSETypeAES256 SSEType = "AES256"
	SSETypeKMS    SSEType = "KMS"
)

// IndexStatus represents the status of a DynamoDB index.
type IndexStatus string

// IndexStatus constants define the possible statuses of an index.
const (
	IndexStatusCreating IndexStatus = "CREATING"
	IndexStatusActive   IndexStatus = "ACTIVE"
	IndexStatusUpdating IndexStatus = "UPDATING"
	IndexStatusDeleting IndexStatus = "DELETING"
)

// BackupStatus represents the status of a DynamoDB backup.
type BackupStatus string

// BackupStatus constants define the possible statuses of a backup.
const (
	BackupStatusCreating  BackupStatus = "CREATING"
	BackupStatusAvailable BackupStatus = "AVAILABLE"
	BackupStatusDeleted   BackupStatus = "DELETED"
)

// BackupType represents the type of a DynamoDB backup.
type BackupType string

// BackupType constants define the supported backup types.
const (
	BackupTypeUser      BackupType = "USER"
	BackupTypeSystem    BackupType = "SYSTEM"
	BackupTypeAWSBackup BackupType = "AWS_BACKUP"
)

// TTLStatus represents the status of TTL for a DynamoDB table.
type TTLStatus string

// TTLStatus constants define the possible TTL statuses.
const (
	TTLStatusEnabling  TTLStatus = "ENABLING"
	TTLStatusEnabled   TTLStatus = "ENABLED"
	TTLStatusDisabling TTLStatus = "DISABLING"
	TTLStatusDisabled  TTLStatus = "DISABLED"
)

// PointInTimeRecoveryStatus represents the status of point-in-time recovery.
type PointInTimeRecoveryStatus string

// PointInTimeRecoveryStatus constants define the possible point-in-time recovery statuses.
const (
	PITRStatusEnabled  PointInTimeRecoveryStatus = "ENABLED"
	PITRStatusDisabled PointInTimeRecoveryStatus = "DISABLED"
)

// KeySchemaElement represents an element of a key schema.
type KeySchemaElement struct {
	AttributeName string  `json:"attribute_name"`
	KeyType       KeyType `json:"key_type"`
}

// AttributeDefinition represents the definition of an attribute.
type AttributeDefinition struct {
	AttributeName string              `json:"attribute_name"`
	AttributeType ScalarAttributeType `json:"attribute_type"`
}

// ProvisionedThroughput represents the provisioned throughput for a table or index.
type ProvisionedThroughput struct {
	ReadCapacityUnits      int64     `json:"read_capacity_units"`
	WriteCapacityUnits     int64     `json:"write_capacity_units"`
	LastDecreaseDateTime   time.Time `json:"last_decrease_date_time,omitempty"`
	LastIncreaseDateTime   time.Time `json:"last_increase_date_time,omitempty"`
	NumberOfDecreasesToday int64     `json:"number_of_decreases_today,omitempty"`
}

// Projection represents the attributes that are projected from an index.
type Projection struct {
	ProjectionType   string   `json:"projection_type,omitempty"`
	NonKeyAttributes []string `json:"non_key_attributes,omitempty"`
}

// LocalSecondaryIndex represents a local secondary index.
type LocalSecondaryIndex struct {
	IndexName      string              `json:"index_name"`
	KeySchema      []*KeySchemaElement `json:"key_schema"`
	Projection     *Projection         `json:"projection"`
	IndexSizeBytes int64               `json:"index_size_bytes,omitempty"`
	ItemCount      int64               `json:"item_count,omitempty"`
}

// GlobalSecondaryIndex represents a global secondary index.
type GlobalSecondaryIndex struct {
	IndexName             string                 `json:"index_name"`
	IndexArn              string                 `json:"index_arn,omitempty"`
	KeySchema             []*KeySchemaElement    `json:"key_schema"`
	Projection            *Projection            `json:"projection"`
	ProvisionedThroughput *ProvisionedThroughput `json:"provisioned_throughput,omitempty"`
	IndexStatus           IndexStatus            `json:"index_status,omitempty"`
	IndexSizeBytes        int64                  `json:"index_size_bytes,omitempty"`
	ItemCount             int64                  `json:"item_count,omitempty"`
}

// StreamSpecification represents the DynamoDB streams specification for a table.
type StreamSpecification struct {
	StreamEnabled  bool           `json:"stream_enabled"`
	StreamViewType StreamViewType `json:"stream_view_type,omitempty"`
}

// SSEDescription represents the server-side encryption description for a table.
type SSEDescription struct {
	Status                         string    `json:"status,omitempty"`
	SSEType                        SSEType   `json:"sse_type,omitempty"`
	KMSMasterKeyArn                string    `json:"kms_master_key_arn,omitempty"`
	InaccessibleEncryptionDateTime time.Time `json:"inaccessible_encryption_date_time,omitempty"`
}

// TimeToLiveSpecification represents the time-to-live specification for a table.
type TimeToLiveSpecification struct {
	Enabled       bool      `json:"enabled"`
	AttributeName string    `json:"attribute_name,omitempty"`
	Status        TTLStatus `json:"status,omitempty"`
}

// PointInTimeRecoverySpecification represents the point-in-time recovery specification for a table.
type PointInTimeRecoverySpecification struct {
	Enabled bool `json:"enabled"`
}

// PointInTimeRecoveryDescription represents the point-in-time recovery description for a table.
type PointInTimeRecoveryDescription struct {
	Status                     PointInTimeRecoveryStatus `json:"status"`
	EarliestRestorableDateTime time.Time                 `json:"earliest_restorable_date_time,omitempty"`
	LatestRestorableDateTime   time.Time                 `json:"latest_restorable_date_time,omitempty"`
}

// Table represents a DynamoDB table.
type Table struct {
	Name                          string                          `json:"name"`
	ARN                           string                          `json:"arn"`
	Status                        TableStatus                     `json:"status"`
	CreationDateTime              time.Time                       `json:"creation_date_time"`
	LastUpdatedDateTime           time.Time                       `json:"last_updated_date_time,omitempty"`
	KeySchema                     []*KeySchemaElement             `json:"key_schema"`
	AttributeDefinitions          []*AttributeDefinition          `json:"attribute_definitions"`
	ProvisionedThroughput         *ProvisionedThroughput          `json:"provisioned_throughput,omitempty"`
	BillingMode                   BillingMode                     `json:"billing_mode"`
	GlobalSecondaryIndexes        []*GlobalSecondaryIndex         `json:"global_secondary_indexes,omitempty"`
	LocalSecondaryIndexes         []*LocalSecondaryIndex          `json:"local_secondary_indexes,omitempty"`
	StreamSpecification           *StreamSpecification            `json:"stream_specification,omitempty"`
	SSEDescription                *SSEDescription                 `json:"sse_description,omitempty"`
	TableSizeBytes                int64                           `json:"table_size_bytes"`
	ItemCount                     int64                           `json:"item_count"`
	Tags                          []types.Tag                     `json:"tags,omitempty"`
	DeletionProtectionEnabled     bool                            `json:"deletion_protection_enabled,omitempty"`
	StreamArn                     string                          `json:"stream_arn,omitempty"`
	LatestStreamLabel             string                          `json:"latest_stream_label,omitempty"`
	TimeToLive                    *TimeToLiveSpecification        `json:"time_to_live,omitempty"`
	PointInTimeRecovery           *PointInTimeRecoveryDescription `json:"point_in_time_recovery,omitempty"`
	ResourcePolicy                string                          `json:"resource_policy,omitempty"`
	KinesisDataStreamDestinations []*KinesisDataStreamDestination `json:"kinesis_data_stream_destinations,omitempty"`
	ContributorInsightsEnabled    bool                            `json:"contributor_insights_enabled,omitempty"`
	TableClass                    string                          `json:"table_class,omitempty"`
}

// Backup represents a DynamoDB table backup.
type Backup struct {
	BackupName              string                  `json:"backup_name"`
	BackupArn               string                  `json:"backup_arn"`
	SourceTableName         string                  `json:"source_table_name"`
	SourceTableArn          string                  `json:"source_table_arn"`
	SourceTableCreationTime time.Time               `json:"source_table_creation_time,omitempty"`
	SourceTableSizeBytes    int64                   `json:"source_table_size_bytes,omitempty"`
	SourceTableItemCount    int64                   `json:"source_table_item_count,omitempty"`
	BackupStatus            BackupStatus            `json:"backup_status"`
	BackupType              BackupType              `json:"backup_type"`
	BackupCreationDateTime  time.Time               `json:"backup_creation_date_time"`
	BackupSizeBytes         int64                   `json:"backup_size_bytes"`
	BackupExpiryDateTime    time.Time               `json:"backup_expiry_date_time,omitempty"`
	KeySchema               []*KeySchemaElement     `json:"key_schema,omitempty"`
	AttributeDefinitions    []*AttributeDefinition  `json:"attribute_definitions,omitempty"`
	BillingMode             BillingMode             `json:"billing_mode,omitempty"`
	ProvisionedThroughput   *ProvisionedThroughput  `json:"provisioned_throughput,omitempty"`
	GlobalSecondaryIndexes  []*GlobalSecondaryIndex `json:"global_secondary_indexes,omitempty"`
	LocalSecondaryIndexes   []*LocalSecondaryIndex  `json:"local_secondary_indexes,omitempty"`
}

// GlobalTable represents a global table in DynamoDB.
type GlobalTable struct {
	GlobalTableName   string     `json:"global_table_name"`
	GlobalTableArn    string     `json:"global_table_arn"`
	GlobalTableStatus string     `json:"global_table_status"`
	CreationDateTime  time.Time  `json:"creation_date_time"`
	ReplicationGroup  []*Replica `json:"replication_group"`
}

// Replica represents a replica of a global table in a specific region.
type Replica struct {
	RegionName                    string `json:"region_name"`
	ReplicaStatus                 string `json:"replica_status"`
	BillingMode                   string `json:"billing_mode,omitempty"`
	ProvisionedReadCapacityUnits  int64  `json:"provisioned_read_capacity_units,omitempty"`
	ProvisionedWriteCapacityUnits int64  `json:"provisioned_write_capacity_units,omitempty"`
}

// KinesisDataStreamDestination represents a Kinesis data stream destination for a table.
type KinesisDataStreamDestination struct {
	StreamArn                    string `json:"stream_arn"`
	DestinationStatus            string `json:"destination_status"`
	DestinationStatusDescription string `json:"destination_status_description,omitempty"`
}

// AttributeValue represents a DynamoDB attribute value.
type AttributeValue struct {
	S    *string                    `json:"s,omitempty"`
	N    *string                    `json:"n,omitempty"`
	B    []byte                     `json:"b,omitempty"`
	SS   []string                   `json:"ss,omitempty"`
	NS   []string                   `json:"ns,omitempty"`
	BS   [][]byte                   `json:"bs,omitempty"`
	M    map[string]*AttributeValue `json:"m,omitempty"`
	L    []*AttributeValue          `json:"l,omitempty"`
	NULL *bool                      `json:"null,omitempty"`
	BOOL *bool                      `json:"bool,omitempty"`
}

// Item represents a DynamoDB item.
type Item struct {
	TableName  string                     `json:"table_name"`
	Key        map[string]*AttributeValue `json:"key"`
	Attributes map[string]*AttributeValue `json:"attributes"`
}

// IsString returns true if the attribute value is a string.
func (av *AttributeValue) IsString() bool {
	return av.S != nil
}

// IsNumber returns true if the attribute value is a number.
func (av *AttributeValue) IsNumber() bool {
	return av.N != nil
}

// IsBinary returns true if the attribute value is binary data.
func (av *AttributeValue) IsBinary() bool {
	return av.B != nil
}

// IsBool returns true if the attribute value is a boolean.
func (av *AttributeValue) IsBool() bool {
	return av.BOOL != nil
}

// IsNull returns true if the attribute value is null.
func (av *AttributeValue) IsNull() bool {
	return av.NULL != nil && *av.NULL
}

// IsMap returns true if the attribute value is a map.
func (av *AttributeValue) IsMap() bool {
	return av.M != nil
}

// IsList returns true if the attribute value is a list.
func (av *AttributeValue) IsList() bool {
	return av.L != nil
}

// IsStringSet returns true if the attribute value is a string set.
func (av *AttributeValue) IsStringSet() bool {
	return av.SS != nil
}

// IsNumberSet returns true if the attribute value is a number set.
func (av *AttributeValue) IsNumberSet() bool {
	return av.NS != nil
}

// IsBinarySet returns true if the attribute value is a binary set.
func (av *AttributeValue) IsBinarySet() bool {
	return av.BS != nil
}

// StringValue creates a new string attribute value.
func StringValue(s string) *AttributeValue {
	return &AttributeValue{S: &s}
}

// NumberValue creates a new number attribute value.
func NumberValue(n string) *AttributeValue {
	return &AttributeValue{N: &n}
}

// BinaryValue creates a new binary attribute value.
func BinaryValue(b []byte) *AttributeValue {
	return &AttributeValue{B: b}
}

// BoolValue creates a new boolean attribute value.
func BoolValue(b bool) *AttributeValue {
	return &AttributeValue{BOOL: &b}
}

// NullValue creates a new null attribute value.
func NullValue() *AttributeValue {
	t := true
	return &AttributeValue{NULL: &t}
}

// MapValue creates a new map attribute value.
func MapValue(m map[string]*AttributeValue) *AttributeValue {
	return &AttributeValue{M: m}
}

// ListValue creates a new list attribute value.
func ListValue(l []*AttributeValue) *AttributeValue {
	return &AttributeValue{L: l}
}

// StringSet creates a new string set attribute value.
func StringSet(ss []string) *AttributeValue {
	return &AttributeValue{SS: ss}
}

// NumberSet creates a new number set attribute value.
func NumberSet(ns []string) *AttributeValue {
	return &AttributeValue{NS: ns}
}

// BinarySet creates a new binary set attribute value.
func BinarySet(bs [][]byte) *AttributeValue {
	return &AttributeValue{BS: bs}
}

// GetKeyAttributeValue retrieves a key attribute value by name.
func (i *Item) GetKeyAttributeValue(attrName string) *AttributeValue {
	if i.Key != nil {
		return i.Key[attrName]
	}
	return nil
}

// GetAttribute retrieves an attribute value by name.
func (i *Item) GetAttribute(attrName string) *AttributeValue {
	if i.Attributes != nil {
		return i.Attributes[attrName]
	}
	return nil
}

// TableListResult represents the result of listing DynamoDB tables.
type TableListResult struct {
	Tables      []*Table
	NextToken   string
	IsTruncated bool
}

// TTLSpecification represents the TTL specification for a DynamoDB table.
type TTLSpecification struct {
	TableName     string `json:"table_name"`
	Enabled       bool   `json:"enabled"`
	AttributeName string `json:"attribute_name,omitempty"`
}

// Endpoint represents a DynamoDB endpoint.
type Endpoint struct {
	Address              string `json:"address"`
	CachePeriodInMinutes int64  `json:"cache_period_in_minutes"`
}

// ImportTableDescription represents the description of a table import.
type ImportTableDescription struct {
	ImportArn          string          `json:"import_arn"`
	ImportStatus       string          `json:"import_status"`
	TableArn           string          `json:"table_arn,omitempty"`
	TableId            string          `json:"table_id,omitempty"`
	StartTime          time.Time       `json:"start_time,omitempty"`
	EndTime            time.Time       `json:"end_time,omitempty"`
	ProcessedItemCount int64           `json:"processed_item_count,omitempty"`
	ProcessedSizeBytes int64           `json:"processed_size_bytes,omitempty"`
	InputFormat        string          `json:"input_format,omitempty"`
	S3BucketSource     *S3BucketSource `json:"s3_bucket_source,omitempty"`
	FailureCode        string          `json:"failure_code,omitempty"`
	FailureMessage     string          `json:"failure_message,omitempty"`
}

// S3BucketSource represents an S3 bucket source for table import.
type S3BucketSource struct {
	S3Bucket      string `json:"s3_bucket,omitempty"`
	S3Prefix      string `json:"s3_prefix,omitempty"`
	S3BucketOwner string `json:"s3_bucket_owner,omitempty"`
}

// ExportDescription represents the description of a table export.
type ExportDescription struct {
	ExportArn         string    `json:"export_arn"`
	ExportStatus      string    `json:"export_status"`
	StartTime         time.Time `json:"start_time,omitempty"`
	EndTime           time.Time `json:"end_time,omitempty"`
	ManifestFilesSize int64     `json:"manifest_files_size,omitempty"`
	ItemCount         int64     `json:"item_count,omitempty"`
	BilledSizeBytes   int64     `json:"billed_size_bytes,omitempty"`
	TableArn          string    `json:"table_arn,omitempty"`
	TableId           string    `json:"table_id,omitempty"`
	ExportFormat      string    `json:"export_format,omitempty"`
	S3Bucket          string    `json:"s3_bucket,omitempty"`
	S3Prefix          string    `json:"s3_prefix,omitempty"`
	FailureCode       string    `json:"failure_code,omitempty"`
	FailureMessage    string    `json:"failure_message,omitempty"`
}

// ContributorInsightsSummary represents the contributor insights summary for a table or index.
type ContributorInsightsSummary struct {
	TableName                 string `json:"table_name"`
	IndexName                 string `json:"index_name,omitempty"`
	ContributorInsightsStatus string `json:"contributor_insights_status"`
}
