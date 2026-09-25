// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"time"
)

// TableStatus represents the status of a DynamoDB table.
type TableStatus string

// TableStatus constants define the possible statuses of a DynamoDB table.
const (
	TableStatusCreating TableStatus = "CREATING"
	TableStatusActive   TableStatus = "ACTIVE"
	TableStatusUpdating TableStatus = "UPDATING"
	TableStatusDeleting TableStatus = "DELETING"
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

// DeletedTableBackupSuffix is the standard naming convention of the system
// backup a recovery-enabled table's deletion writes: the developer guide
// ("Delete a table with PITR enabled") states "All system backups follow a
// standard naming convention of table-name$DeletedTableBackup".
const DeletedTableBackupSuffix = "$DeletedTableBackup"

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

// TableClass represents the storage class of a table or replica.
type TableClass string

// TableClass constants define the supported table classes; the class is
// applied to a replica through ReplicaTableClass.
const (
	TableClassStandard                 TableClass = "STANDARD"
	TableClassStandardInfrequentAccess TableClass = "STANDARD_INFREQUENT_ACCESS"
)

// SSEStatus represents the status of server-side encryption on a table.
type SSEStatus string

// SSEStatus constants define the possible server-side encryption statuses.
const (
	SSEStatusEnabling  SSEStatus = "ENABLING"
	SSEStatusEnabled   SSEStatus = "ENABLED"
	SSEStatusDisabling SSEStatus = "DISABLING"
	SSEStatusDisabled  SSEStatus = "DISABLED"
	SSEStatusUpdating  SSEStatus = "UPDATING"
)

// DestinationStatus represents the status of a Kinesis data stream
// destination attached to a table.
type DestinationStatus string

// DestinationStatus constants define the possible destination statuses.
const (
	DestinationStatusEnabling     DestinationStatus = "ENABLING"
	DestinationStatusActive       DestinationStatus = "ACTIVE"
	DestinationStatusDisabling    DestinationStatus = "DISABLING"
	DestinationStatusDisabled     DestinationStatus = "DISABLED"
	DestinationStatusEnableFailed DestinationStatus = "ENABLE_FAILED"
	DestinationStatusUpdating     DestinationStatus = "UPDATING"
)

// ContributorInsightsMode represents the mode of contributor insights on a
// table or index.
type ContributorInsightsMode string

// ContributorInsightsMode constants define the supported modes.
const (
	ContributorInsightsModeThrottledKeys            ContributorInsightsMode = "THROTTLED_KEYS"
	ContributorInsightsModeAccessedAndThrottledKeys ContributorInsightsMode = "ACCESSED_AND_THROTTLED_KEYS"
)

// ApproximateCreationDateTimePrecision represents the precision of the
// approximate creation timestamp a Kinesis destination records.
type ApproximateCreationDateTimePrecision string

// ApproximateCreationDateTimePrecision constants define the supported
// precisions.
const (
	ACDTPrecisionMillisecond ApproximateCreationDateTimePrecision = "MILLISECOND"
	ACDTPrecisionMicrosecond ApproximateCreationDateTimePrecision = "MICROSECOND"
)

// GlobalTableStatus represents the status of a global table.
type GlobalTableStatus string

// GlobalTableStatus constants define the possible global table statuses.
const (
	GlobalTableStatusCreating GlobalTableStatus = "CREATING"
	GlobalTableStatusActive   GlobalTableStatus = "ACTIVE"
	GlobalTableStatusUpdating GlobalTableStatus = "UPDATING"
	GlobalTableStatusDeleting GlobalTableStatus = "DELETING"
)

// ReplicaStatus represents the status of one replica in a global table's
// replication group.
type ReplicaStatus string

// ReplicaStatus constants define the possible replica statuses.
const (
	ReplicaStatusCreating                          ReplicaStatus = "CREATING"
	ReplicaStatusCreationFailed                    ReplicaStatus = "CREATION_FAILED"
	ReplicaStatusUpdating                          ReplicaStatus = "UPDATING"
	ReplicaStatusDeleting                          ReplicaStatus = "DELETING"
	ReplicaStatusActive                            ReplicaStatus = "ACTIVE"
	ReplicaStatusArchiving                         ReplicaStatus = "ARCHIVING"
	ReplicaStatusArchived                          ReplicaStatus = "ARCHIVED"
	ReplicaStatusInaccessibleEncryptionCredentials ReplicaStatus = "INACCESSIBLE_ENCRYPTION_CREDENTIALS"
	ReplicaStatusRegionDisabled                    ReplicaStatus = "REGION_DISABLED"
	ReplicaStatusReplicationNotAuthorized          ReplicaStatus = "REPLICATION_NOT_AUTHORIZED"
)

// ImportStatus represents the status of a table import.
type ImportStatus string

// ImportStatus constants define the possible import statuses.
const (
	ImportStatusInProgress ImportStatus = "IN_PROGRESS"
	ImportStatusCompleted  ImportStatus = "COMPLETED"
	ImportStatusCancelling ImportStatus = "CANCELLING"
	ImportStatusCancelled  ImportStatus = "CANCELLED"
	ImportStatusFailed     ImportStatus = "FAILED"
)

// ExportStatus represents the status of a table export.
type ExportStatus string

// ExportStatus constants define the possible export statuses.
const (
	ExportStatusInProgress ExportStatus = "IN_PROGRESS"
	ExportStatusCompleted  ExportStatus = "COMPLETED"
	ExportStatusFailed     ExportStatus = "FAILED"
)

// InputFormat represents the format of the source data of an import.
type InputFormat string

// InputFormat constants define the supported input formats.
const (
	InputFormatDynamoDBJSON InputFormat = "DYNAMODB_JSON"
	InputFormatIon          InputFormat = "ION"
	InputFormatCSV          InputFormat = "CSV"
)

// InputCompressionType represents the compression of the source data of an
// import.
type InputCompressionType string

// InputCompressionType constants define the supported compressions.
const (
	InputCompressionTypeGzip InputCompressionType = "GZIP"
	InputCompressionTypeZstd InputCompressionType = "ZSTD"
	InputCompressionTypeNone InputCompressionType = "NONE"
)

// ExportFormat represents the output format of a table export.
type ExportFormat string

// ExportFormat constants define the supported export formats.
const (
	ExportFormatDynamoDBJSON ExportFormat = "DYNAMODB_JSON"
	ExportFormatIon          ExportFormat = "ION"
)

// ExportType represents the kind of a table export.
type ExportType string

// ExportType constants define the supported export types.
const (
	ExportTypeFullExport        ExportType = "FULL_EXPORT"
	ExportTypeIncrementalExport ExportType = "INCREMENTAL_EXPORT"
)

// ExportViewType represents the image view an export writes.
type ExportViewType string

// ExportViewType constants define the supported views.
const (
	ExportViewTypeNewImage        ExportViewType = "NEW_IMAGE"
	ExportViewTypeNewAndOldImages ExportViewType = "NEW_AND_OLD_IMAGES"
)

// S3SseAlgorithm represents the server-side encryption algorithm an export
// applies on its S3 output.
type S3SseAlgorithm string

// S3SseAlgorithm constants define the supported algorithms.
const (
	S3SseAlgorithmAES256 S3SseAlgorithm = "AES256"
	S3SseAlgorithmKMS    S3SseAlgorithm = "KMS"
)

// VectorDistanceFunction represents the distance metric a vector index
// scores matches by.
type VectorDistanceFunction string

// VectorDistanceFunction constants define the supported metrics.
const (
	VectorDistanceFunctionCosine     VectorDistanceFunction = "COSINE"
	VectorDistanceFunctionEuclidean  VectorDistanceFunction = "EUCLIDEAN"
	VectorDistanceFunctionDotProduct VectorDistanceFunction = "DOT_PRODUCT"
)

// SearchSchemaElementType represents an attribute's role in a vector index
// search schema: HASH partitions the index, INLINE_FILTER is projected for
// filtering.
type SearchSchemaElementType string

// SearchSchemaElementType constants define the supported roles.
const (
	SearchSchemaElementTypeHash         SearchSchemaElementType = "HASH"
	SearchSchemaElementTypeInlineFilter SearchSchemaElementType = "INLINE_FILTER"
)

// ProjectionType represents which attributes a secondary index projects.
type ProjectionType string

// ProjectionType constants define the supported projection types.
const (
	ProjectionTypeAll      ProjectionType = "ALL"
	ProjectionTypeKeysOnly ProjectionType = "KEYS_ONLY"
	ProjectionTypeInclude  ProjectionType = "INCLUDE"
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
	ProjectionType   ProjectionType `json:"projection_type,omitempty"`
	NonKeyAttributes []string       `json:"non_key_attributes,omitempty"`
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
	OnDemandThroughput    *OnDemandThroughput    `json:"on_demand_throughput,omitempty"`
	WarmThroughput        *WarmThroughput        `json:"warm_throughput,omitempty"`
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
	Status                         SSEStatus `json:"status,omitempty"`
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

// PointInTimeRecoveryDescription represents the point-in-time recovery description for a table.
type PointInTimeRecoveryDescription struct {
	Status                     PointInTimeRecoveryStatus `json:"status"`
	EarliestRestorableDateTime time.Time                 `json:"earliest_restorable_date_time,omitempty"`
	LatestRestorableDateTime   time.Time                 `json:"latest_restorable_date_time,omitempty"`
	RecoveryPeriodInDays       int                       `json:"recovery_period_in_days,omitempty"`
}

// SearchSchemaElement defines an attribute's role in a vector index search
// schema: HASH partitions the index, INLINE_FILTER is projected for filtering.
type SearchSchemaElement struct {
	AttributeName           string                  `json:"attribute_name"`
	SearchSchemaElementType SearchSchemaElementType `json:"search_schema_element_type,omitempty"`
}

// VectorIndex represents a vector index on a table. A vector index enables
// similarity search over one vector attribute, whose item value is a list of
// numbers.
type VectorIndex struct {
	IndexName           string                 `json:"index_name"`
	IndexArn            string                 `json:"index_arn,omitempty"`
	VectorAttributeName string                 `json:"vector_attribute_name"`
	Dimensions          int64                  `json:"dimensions"`
	DistanceFunction    VectorDistanceFunction `json:"distance_function"`
	Projection          *Projection            `json:"projection"`
	SearchSchema        []*SearchSchemaElement `json:"search_schema,omitempty"`
	IndexStatus         IndexStatus            `json:"index_status,omitempty"`
	Backfilling         bool                   `json:"backfilling,omitempty"`
	IndexSizeBytes      int64                  `json:"index_size_bytes,omitempty"`
	ItemCount           int64                  `json:"item_count,omitempty"`
}

// Vector index bounds. The quotas are documented in the DynamoDB quotas
// page's Vector indexes table: dimensions 1..4096 (not adjustable), at most
// 5 vector indexes per table, at most 18 inline-filter elements per index,
// and one partition key (HASH) per index; the TopK range is documented
// 1..100 on the SearchVectors API reference.
const (
	VectorDimensionsMin    = 1
	VectorDimensionsMax    = 4096
	VectorAttributeNameMax = 255
	VectorSearchTopKMin    = 1
	VectorSearchTopKMax    = 100
	VectorIndexesPerTable  = 5
	VectorInlineFiltersMax = 18
)

// Secondary-index and billing-mode quotas from the DynamoDB quotas page:
// at most 20 global secondary indexes and 5 local secondary indexes per
// table, at most 100 user-specified projected attribute names combined
// across a table's secondary indexes, and at most four provisioned to
// on-demand billing mode switches inside a rolling 24-hour window.
const (
	GlobalSecondaryIndexesPerTable = 20
	LocalSecondaryIndexesPerTable  = 5
	ProjectedAttributesPerTable    = 100
	BillingModeSwitchesPerDay      = 4
)

// Table represents a DynamoDB table.
type Table struct {
	Name                          string                          `json:"name"`
	ARN                           string                          `json:"arn"`
	TableId                       string                          `json:"table_id,omitempty"`
	Status                        TableStatus                     `json:"status"`
	CreationDateTime              time.Time                       `json:"creation_date_time"`
	LastUpdatedDateTime           time.Time                       `json:"last_updated_date_time,omitempty"`
	KeySchema                     []*KeySchemaElement             `json:"key_schema"`
	AttributeDefinitions          []*AttributeDefinition          `json:"attribute_definitions"`
	ProvisionedThroughput         *ProvisionedThroughput          `json:"provisioned_throughput,omitempty"`
	BillingMode                   BillingMode                     `json:"billing_mode"`
	GlobalSecondaryIndexes        []*GlobalSecondaryIndex         `json:"global_secondary_indexes,omitempty"`
	LocalSecondaryIndexes         []*LocalSecondaryIndex          `json:"local_secondary_indexes,omitempty"`
	VectorIndexes                 []*VectorIndex                  `json:"vector_indexes,omitempty"`
	StreamSpecification           *StreamSpecification            `json:"stream_specification,omitempty"`
	SSEDescription                *SSEDescription                 `json:"sse_description,omitempty"`
	TableSizeBytes                int64                           `json:"table_size_bytes"`
	ItemCount                     int64                           `json:"item_count"`
	DeletionProtectionEnabled     bool                            `json:"deletion_protection_enabled,omitempty"`
	StreamArn                     string                          `json:"stream_arn,omitempty"`
	LatestStreamLabel             string                          `json:"latest_stream_label,omitempty"`
	TimeToLive                    *TimeToLiveSpecification        `json:"time_to_live,omitempty"`
	PointInTimeRecovery           *PointInTimeRecoveryDescription `json:"point_in_time_recovery,omitempty"`
	ResourcePolicy                string                          `json:"resource_policy,omitempty"`
	ResourcePolicyRevisionId      int                             `json:"resource_policy_revision_id,omitempty"`
	KinesisDataStreamDestinations []*KinesisDataStreamDestination `json:"kinesis_data_stream_destinations,omitempty"`
	ContributorInsightsEnabled    bool                            `json:"contributor_insights_enabled,omitempty"`
	ContributorInsightsMode       ContributorInsightsMode         `json:"contributor_insights_mode,omitempty"`
	ContributorInsightsUpdatedAt  time.Time                       `json:"contributor_insights_updated_at,omitempty"`
	WarmThroughput                *WarmThroughput                 `json:"warm_throughput,omitempty"`
	OnDemandThroughput            *OnDemandThroughput             `json:"on_demand_throughput,omitempty"`
	GlobalTableSourceArn          string                          `json:"global_table_source_arn,omitempty"`
	TableClass                    TableClass                      `json:"table_class,omitempty"`
	RestoreSummary                *RestoreSummary                 `json:"restore_summary,omitempty"`
	// BillingModeSwitches records the provisioned-to-on-demand switch
	// timestamps inside the rolling 24-hour quota window.
	BillingModeSwitches []time.Time `json:"billing_mode_switches,omitempty"`
}

// RestoreSummary records how a table was created by a restore operation.
type RestoreSummary struct {
	SourceBackupArn   string    `json:"source_backup_arn,omitempty"`
	SourceTableArn    string    `json:"source_table_arn,omitempty"`
	RestoreDateTime   time.Time `json:"restore_date_time,omitempty"`
	RestoreInProgress bool      `json:"restore_in_progress"`
}

// Backup represents a DynamoDB table backup.
type Backup struct {
	BackupName              string                  `json:"backup_name"`
	BackupArn               string                  `json:"backup_arn"`
	SourceTableName         string                  `json:"source_table_name"`
	SourceTableArn          string                  `json:"source_table_arn"`
	SourceTableId           string                  `json:"source_table_id,omitempty"`
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
	VectorIndexes           []*VectorIndex          `json:"vector_indexes,omitempty"`
}

// GlobalTable represents a global table in DynamoDB.
type GlobalTable struct {
	GlobalTableName   string            `json:"global_table_name"`
	GlobalTableArn    string            `json:"global_table_arn"`
	GlobalTableStatus GlobalTableStatus `json:"global_table_status"`
	CreationDateTime  time.Time         `json:"creation_date_time"`
	ReplicationGroup  []*Replica        `json:"replication_group"`
	// Write-capacity auto-scaling applied at the global level; echoed on
	// every replica description.
	WriteAutoScalingSettings *AutoScalingSettingsDescription `json:"write_auto_scaling_settings,omitempty"`
	// Per-index write-capacity settings applied at the global level; they
	// merge with each replica's per-index read settings at echo time.
	GlobalSecondaryIndexWriteSettings []IndexAutoScalingSettings `json:"global_secondary_index_write_settings,omitempty"`
}

// Replica represents a replica of a global table in a specific region.
type Replica struct {
	RegionName                    string        `json:"region_name"`
	ReplicaStatus                 ReplicaStatus `json:"replica_status"`
	BillingMode                   BillingMode   `json:"billing_mode,omitempty"`
	ProvisionedReadCapacityUnits  int64         `json:"provisioned_read_capacity_units,omitempty"`
	ProvisionedWriteCapacityUnits int64         `json:"provisioned_write_capacity_units,omitempty"`
	// Read-capacity auto-scaling applied per replica, echoed through
	// ReplicaProvisionedReadCapacityAutoScalingSettings.
	ReadAutoScalingSettings *AutoScalingSettingsDescription `json:"read_auto_scaling_settings,omitempty"`
	// Per-index read-capacity settings applied per replica; the write side
	// lives on the global table and the two merge at echo time.
	GlobalSecondaryIndexReadSettings []IndexAutoScalingSettings `json:"global_secondary_index_read_settings,omitempty"`
	// ReplicaTableClass and its update time, echoed through
	// ReplicaTableClassSummary.
	TableClass            TableClass `json:"table_class,omitempty"`
	TableClassLastUpdated *time.Time `json:"table_class_last_updated,omitempty"`
	// The replica Update action's recorded overrides: the KMS key
	// identifier the replica encrypts under, the on-demand read maximum
	// override, and the per-index capacity overrides.
	KMSMasterKeyId                string                         `json:"kms_master_key_id,omitempty"`
	OnDemandThroughputOverride    *OnDemandThroughput            `json:"on_demand_throughput_override,omitempty"`
	GlobalSecondaryIndexOverrides []*ReplicaGlobalSecondaryIndex `json:"global_secondary_index_overrides,omitempty"`
}

// ReplicaGlobalSecondaryIndex is a replica's per-index capacity override,
// applied through the replica Update action's GlobalSecondaryIndexes
// member: the read-side provisioned override and the on-demand maximum
// override an index carries in one replica region.
type ReplicaGlobalSecondaryIndex struct {
	IndexName                    string              `json:"index_name"`
	ProvisionedReadCapacityUnits int64               `json:"provisioned_read_capacity_units,omitempty"`
	OnDemandThroughputOverride   *OnDemandThroughput `json:"on_demand_throughput_override,omitempty"`
}

// TargetTrackingScalingPolicyConfiguration mirrors the model's
// AutoScalingTargetTrackingScalingPolicyConfigurationDescription.
type TargetTrackingScalingPolicyConfiguration struct {
	DisableScaleIn   *bool
	ScaleInCooldown  *int32
	ScaleOutCooldown *int32
	TargetValue      float64
}

// AutoScalingPolicyDescription mirrors the model's AutoScalingPolicyDescription.
type AutoScalingPolicyDescription struct {
	PolicyName                               *string
	TargetTrackingScalingPolicyConfiguration *TargetTrackingScalingPolicyConfiguration
}

// AutoScalingSettingsDescription mirrors the model's
// AutoScalingSettingsDescription: one capacity dimension's auto-scaling
// state. Members are pointers because absence is meaningful on the wire.
type AutoScalingSettingsDescription struct {
	MinimumUnits        *int64
	MaximumUnits        *int64
	AutoScalingDisabled *bool
	AutoScalingRoleArn  *string
	ScalingPolicies     []AutoScalingPolicyDescription
}

// IndexAutoScalingSettings is one index's stored settings: capacity units
// plus the read and write auto-scaling descriptions. A write-side entry
// leaves the read members nil and vice versa.
type IndexAutoScalingSettings struct {
	IndexName                     string
	ProvisionedReadCapacityUnits  *int64
	ProvisionedWriteCapacityUnits *int64
	Read                          *AutoScalingSettingsDescription
	Write                         *AutoScalingSettingsDescription
}

// ReplicaAutoScalingDescription is one replica's stored auto-scaling state
// for UpdateTableReplicaAutoScaling / DescribeTableReplicaAutoScaling.
type ReplicaAutoScalingDescription struct {
	RegionName             string
	Read                   *AutoScalingSettingsDescription
	Write                  *AutoScalingSettingsDescription
	GlobalSecondaryIndexes []IndexAutoScalingSettings
}

// TableReplicaAutoScalingSettings is the stored table-level auto-scaling
// record: the per-replica descriptions keyed by region.
type TableReplicaAutoScalingSettings struct {
	Replicas []ReplicaAutoScalingDescription
}

// KinesisDataStreamDestination represents a Kinesis data stream destination for a table.
type KinesisDataStreamDestination struct {
	StreamArn                            string                               `json:"stream_arn"`
	DestinationStatus                    DestinationStatus                    `json:"destination_status"`
	DestinationStatusDescription         string                               `json:"destination_status_description,omitempty"`
	ApproximateCreationDateTimePrecision ApproximateCreationDateTimePrecision `json:"approximate_creation_date_time_precision,omitempty"`
}

// WarmThroughput represents the warm throughput for a DynamoDB table.
type WarmThroughput struct {
	ReadUnitsPerSecond  int64 `json:"read_units_per_second,omitempty"`
	WriteUnitsPerSecond int64 `json:"write_units_per_second,omitempty"`
}

// OnDemandThroughput represents the on-demand throughput settings.
type OnDemandThroughput struct {
	MaxReadRequestUnits  int64 `json:"max_read_request_units,omitempty"`
	MaxWriteRequestUnits int64 `json:"max_write_request_units,omitempty"`
}

// The model's on-demand maximum bounds: a present member is either a
// limit of at least OnDemandThroughputMinUnits or the documented removal
// sentinel OnDemandThroughputRemoveValue ("set the value ... to -1").
const (
	OnDemandThroughputMinUnits    = 1
	OnDemandThroughputRemoveValue = -1
)

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

// IsNull returns true if the attribute value is null. Pointer presence is
// the discriminator, matching the wire parser and the sibling Is* checks.
func (av *AttributeValue) IsNull() bool {
	return av.NULL != nil
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

// ImportTableDescription represents the description of a table import.
type ImportTableDescription struct {
	ImportArn            string               `json:"import_arn"`
	ImportStatus         ImportStatus         `json:"import_status"`
	TableArn             string               `json:"table_arn,omitempty"`
	TableId              string               `json:"table_id,omitempty"`
	StartTime            time.Time            `json:"start_time,omitempty"`
	EndTime              time.Time            `json:"end_time,omitempty"`
	ProcessedItemCount   int64                `json:"processed_item_count,omitempty"`
	ProcessedSizeBytes   int64                `json:"processed_size_bytes,omitempty"`
	ImportedItemCount    int64                `json:"imported_item_count,omitempty"`
	ErrorCount           int64                `json:"error_count,omitempty"`
	InputFormat          InputFormat          `json:"input_format,omitempty"`
	S3BucketSource       *S3BucketSource      `json:"s3_bucket_source,omitempty"`
	FailureCode          string               `json:"failure_code,omitempty"`
	FailureMessage       string               `json:"failure_message,omitempty"`
	ClientToken          string               `json:"client_token,omitempty"`
	InputCompressionType InputCompressionType `json:"input_compression_type,omitempty"`
}

// S3BucketSource represents an S3 bucket source for table import.
type S3BucketSource struct {
	S3Bucket      string `json:"s3_bucket,omitempty"`
	S3Prefix      string `json:"s3_prefix,omitempty"`
	S3BucketOwner string `json:"s3_bucket_owner,omitempty"`
}

// ExportDescription represents the description of a table export.
type ExportDescription struct {
	ExportArn         string         `json:"export_arn"`
	ExportStatus      ExportStatus   `json:"export_status"`
	StartTime         time.Time      `json:"start_time,omitempty"`
	EndTime           time.Time      `json:"end_time,omitempty"`
	ManifestFilesSize int64          `json:"manifest_files_size,omitempty"`
	ItemCount         int64          `json:"item_count,omitempty"`
	BilledSizeBytes   int64          `json:"billed_size_bytes,omitempty"`
	ExportTime        time.Time      `json:"export_time,omitempty"`
	TableArn          string         `json:"table_arn,omitempty"`
	TableId           string         `json:"table_id,omitempty"`
	ExportFormat      ExportFormat   `json:"export_format,omitempty"`
	S3Bucket          string         `json:"s3_bucket,omitempty"`
	S3Prefix          string         `json:"s3_prefix,omitempty"`
	FailureCode       string         `json:"failure_code,omitempty"`
	FailureMessage    string         `json:"failure_message,omitempty"`
	ClientToken       string         `json:"client_token,omitempty"`
	S3BucketOwner     string         `json:"s3_bucket_owner,omitempty"`
	S3SseKmsKeyId     string         `json:"s3_sse_kms_key_id,omitempty"`
	ExportManifest    string         `json:"export_manifest,omitempty"`
	ExportType        ExportType     `json:"export_type,omitempty"`
	S3SseAlgorithm    S3SseAlgorithm `json:"s3_sse_algorithm,omitempty"`
	ExportViewType    ExportViewType `json:"export_view_type,omitempty"`
	ExportFromTime    time.Time      `json:"export_from_time,omitempty"`
	ExportToTime      time.Time      `json:"export_to_time,omitempty"`
}

// ContributorInsightsSummary represents the contributor insights summary for a table or index.
type ContributorInsightsSummary struct {
	TableName                 string `json:"table_name"`
	IndexName                 string `json:"index_name,omitempty"`
	ContributorInsightsStatus string `json:"contributor_insights_status"`
}
