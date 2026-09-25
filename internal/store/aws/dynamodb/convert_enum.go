package dynamodb

import (
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

// enumWire binds one domain value to the persisted protobuf enum value it
// serialises as. The slice per domain is that domain's single value-set
// definition: both conversion directions read it, so the pair cannot drift
// apart, and a value unknown to the table round-trips as the zero value on
// either side (UNSPECIFIED in, empty string out).
type enumWire[P ~int32, T ~string] struct {
	proto P
	value T
}

// enumToProto resolves a domain value through its wire table, falling back
// to the domain's UNSPECIFIED member.
func enumToProto[P ~int32, T ~string](table []enumWire[P, T], v T, unspecified P) P {
	for _, entry := range table {
		if entry.value == v {
			return entry.proto
		}
	}
	return unspecified
}

// enumFromProto resolves a persisted enum value through its wire table,
// falling back to the domain's zero value.
func enumFromProto[P ~int32, T ~string](table []enumWire[P, T], p P) T {
	for _, entry := range table {
		if entry.proto == p {
			return entry.value
		}
	}
	return ""
}

var tableStatusWire = []enumWire[pb.TableStatus, TableStatus]{
	{pb.TableStatus_TABLE_STATUS_CREATING, TableStatusCreating},
	{pb.TableStatus_TABLE_STATUS_ACTIVE, TableStatusActive},
	{pb.TableStatus_TABLE_STATUS_UPDATING, TableStatusUpdating},
	{pb.TableStatus_TABLE_STATUS_DELETING, TableStatusDeleting},
}

var billingModeWire = []enumWire[pb.BillingMode, BillingMode]{
	{pb.BillingMode_BILLING_MODE_PROVISIONED, BillingModeProvisioned},
	{pb.BillingMode_BILLING_MODE_PAY_PER_REQUEST, BillingModePayPerRequest},
}

var tableClassWire = []enumWire[pb.TableClass, TableClass]{
	{pb.TableClass_TABLE_CLASS_STANDARD, TableClassStandard},
	{pb.TableClass_TABLE_CLASS_STANDARD_INFREQUENT_ACCESS, TableClassStandardInfrequentAccess},
}

var sseStatusWire = []enumWire[pb.SSEStatus, SSEStatus]{
	{pb.SSEStatus_SSE_STATUS_ENABLING, SSEStatusEnabling},
	{pb.SSEStatus_SSE_STATUS_ENABLED, SSEStatusEnabled},
	{pb.SSEStatus_SSE_STATUS_DISABLING, SSEStatusDisabling},
	{pb.SSEStatus_SSE_STATUS_DISABLED, SSEStatusDisabled},
	{pb.SSEStatus_SSE_STATUS_UPDATING, SSEStatusUpdating},
}

var sseTypeWire = []enumWire[pb.SSEType, SSEType]{
	{pb.SSEType_SSE_TYPE_AES256, SSETypeAES256},
	{pb.SSEType_SSE_TYPE_KMS, SSETypeKMS},
}

var destinationStatusWire = []enumWire[pb.DestinationStatus, DestinationStatus]{
	{pb.DestinationStatus_DESTINATION_STATUS_ENABLING, DestinationStatusEnabling},
	{pb.DestinationStatus_DESTINATION_STATUS_ACTIVE, DestinationStatusActive},
	{pb.DestinationStatus_DESTINATION_STATUS_DISABLING, DestinationStatusDisabling},
	{pb.DestinationStatus_DESTINATION_STATUS_DISABLED, DestinationStatusDisabled},
	{pb.DestinationStatus_DESTINATION_STATUS_ENABLE_FAILED, DestinationStatusEnableFailed},
	{pb.DestinationStatus_DESTINATION_STATUS_UPDATING, DestinationStatusUpdating},
}

var contributorInsightsModeWire = []enumWire[pb.ContributorInsightsMode, ContributorInsightsMode]{
	{pb.ContributorInsightsMode_CONTRIBUTOR_INSIGHTS_MODE_THROTTLED_KEYS, ContributorInsightsModeThrottledKeys},
	{pb.ContributorInsightsMode_CONTRIBUTOR_INSIGHTS_MODE_ACCESSED_AND_THROTTLED_KEYS, ContributorInsightsModeAccessedAndThrottledKeys},
}

var acdtPrecisionWire = []enumWire[pb.ApproximateCreationDateTimePrecision, ApproximateCreationDateTimePrecision]{
	{pb.ApproximateCreationDateTimePrecision_APPROXIMATE_CREATION_DATE_TIME_PRECISION_MILLISECOND, ACDTPrecisionMillisecond},
	{pb.ApproximateCreationDateTimePrecision_APPROXIMATE_CREATION_DATE_TIME_PRECISION_MICROSECOND, ACDTPrecisionMicrosecond},
}

var streamViewTypeWire = []enumWire[pb.StreamViewType, StreamViewType]{
	{pb.StreamViewType_STREAM_VIEW_TYPE_NEW_IMAGE, StreamViewTypeNewImage},
	{pb.StreamViewType_STREAM_VIEW_TYPE_OLD_IMAGE, StreamViewTypeOldImage},
	{pb.StreamViewType_STREAM_VIEW_TYPE_NEW_AND_OLD_IMAGES, StreamViewTypeNewAndOldImages},
	{pb.StreamViewType_STREAM_VIEW_TYPE_KEYS_ONLY, StreamViewTypeKeysOnly},
}

var ttlStatusWire = []enumWire[pb.TTLStatus, TTLStatus]{
	{pb.TTLStatus_TTL_STATUS_ENABLING, TTLStatusEnabling},
	{pb.TTLStatus_TTL_STATUS_ENABLED, TTLStatusEnabled},
	{pb.TTLStatus_TTL_STATUS_DISABLING, TTLStatusDisabling},
	{pb.TTLStatus_TTL_STATUS_DISABLED, TTLStatusDisabled},
}

var pitrStatusWire = []enumWire[pb.PointInTimeRecoveryStatus, PointInTimeRecoveryStatus]{
	{pb.PointInTimeRecoveryStatus_POINT_IN_TIME_RECOVERY_STATUS_ENABLED, PITRStatusEnabled},
	{pb.PointInTimeRecoveryStatus_POINT_IN_TIME_RECOVERY_STATUS_DISABLED, PITRStatusDisabled},
}

var backupStatusWire = []enumWire[pb.BackupStatus, BackupStatus]{
	{pb.BackupStatus_BACKUP_STATUS_CREATING, BackupStatusCreating},
	{pb.BackupStatus_BACKUP_STATUS_AVAILABLE, BackupStatusAvailable},
	{pb.BackupStatus_BACKUP_STATUS_DELETED, BackupStatusDeleted},
}

var backupTypeWire = []enumWire[pb.BackupType, BackupType]{
	{pb.BackupType_BACKUP_TYPE_USER, BackupTypeUser},
	{pb.BackupType_BACKUP_TYPE_SYSTEM, BackupTypeSystem},
	{pb.BackupType_BACKUP_TYPE_AWS_BACKUP, BackupTypeAWSBackup},
}

var globalTableStatusWire = []enumWire[pb.GlobalTableStatus, GlobalTableStatus]{
	{pb.GlobalTableStatus_GLOBAL_TABLE_STATUS_CREATING, GlobalTableStatusCreating},
	{pb.GlobalTableStatus_GLOBAL_TABLE_STATUS_ACTIVE, GlobalTableStatusActive},
	{pb.GlobalTableStatus_GLOBAL_TABLE_STATUS_UPDATING, GlobalTableStatusUpdating},
	{pb.GlobalTableStatus_GLOBAL_TABLE_STATUS_DELETING, GlobalTableStatusDeleting},
}

var replicaStatusWire = []enumWire[pb.ReplicaStatus, ReplicaStatus]{
	{pb.ReplicaStatus_REPLICA_STATUS_CREATING, ReplicaStatusCreating},
	{pb.ReplicaStatus_REPLICA_STATUS_CREATION_FAILED, ReplicaStatusCreationFailed},
	{pb.ReplicaStatus_REPLICA_STATUS_UPDATING, ReplicaStatusUpdating},
	{pb.ReplicaStatus_REPLICA_STATUS_DELETING, ReplicaStatusDeleting},
	{pb.ReplicaStatus_REPLICA_STATUS_ACTIVE, ReplicaStatusActive},
	{pb.ReplicaStatus_REPLICA_STATUS_ARCHIVING, ReplicaStatusArchiving},
	{pb.ReplicaStatus_REPLICA_STATUS_ARCHIVED, ReplicaStatusArchived},
	{pb.ReplicaStatus_REPLICA_STATUS_INACCESSIBLE_ENCRYPTION_CREDENTIALS, ReplicaStatusInaccessibleEncryptionCredentials},
	{pb.ReplicaStatus_REPLICA_STATUS_REGION_DISABLED, ReplicaStatusRegionDisabled},
	{pb.ReplicaStatus_REPLICA_STATUS_REPLICATION_NOT_AUTHORIZED, ReplicaStatusReplicationNotAuthorized},
}

var importStatusWire = []enumWire[pb.ImportStatus, ImportStatus]{
	{pb.ImportStatus_IMPORT_STATUS_IN_PROGRESS, ImportStatusInProgress},
	{pb.ImportStatus_IMPORT_STATUS_COMPLETED, ImportStatusCompleted},
	{pb.ImportStatus_IMPORT_STATUS_CANCELLING, ImportStatusCancelling},
	{pb.ImportStatus_IMPORT_STATUS_CANCELLED, ImportStatusCancelled},
	{pb.ImportStatus_IMPORT_STATUS_FAILED, ImportStatusFailed},
}

var exportStatusWire = []enumWire[pb.ExportStatus, ExportStatus]{
	{pb.ExportStatus_EXPORT_STATUS_IN_PROGRESS, ExportStatusInProgress},
	{pb.ExportStatus_EXPORT_STATUS_COMPLETED, ExportStatusCompleted},
	{pb.ExportStatus_EXPORT_STATUS_FAILED, ExportStatusFailed},
}

var inputFormatWire = []enumWire[pb.InputFormat, InputFormat]{
	{pb.InputFormat_INPUT_FORMAT_DYNAMODB_JSON, InputFormatDynamoDBJSON},
	{pb.InputFormat_INPUT_FORMAT_ION, InputFormatIon},
	{pb.InputFormat_INPUT_FORMAT_CSV, InputFormatCSV},
}

var inputCompressionTypeWire = []enumWire[pb.InputCompressionType, InputCompressionType]{
	{pb.InputCompressionType_INPUT_COMPRESSION_TYPE_GZIP, InputCompressionTypeGzip},
	{pb.InputCompressionType_INPUT_COMPRESSION_TYPE_ZSTD, InputCompressionTypeZstd},
	{pb.InputCompressionType_INPUT_COMPRESSION_TYPE_NONE, InputCompressionTypeNone},
}

var exportFormatWire = []enumWire[pb.ExportFormat, ExportFormat]{
	{pb.ExportFormat_EXPORT_FORMAT_DYNAMODB_JSON, ExportFormatDynamoDBJSON},
	{pb.ExportFormat_EXPORT_FORMAT_ION, ExportFormatIon},
}

var exportTypeWire = []enumWire[pb.ExportType, ExportType]{
	{pb.ExportType_EXPORT_TYPE_FULL_EXPORT, ExportTypeFullExport},
	{pb.ExportType_EXPORT_TYPE_INCREMENTAL_EXPORT, ExportTypeIncrementalExport},
}

var exportViewTypeWire = []enumWire[pb.ExportViewType, ExportViewType]{
	{pb.ExportViewType_EXPORT_VIEW_TYPE_NEW_IMAGE, ExportViewTypeNewImage},
	{pb.ExportViewType_EXPORT_VIEW_TYPE_NEW_AND_OLD_IMAGES, ExportViewTypeNewAndOldImages},
}

var s3SseAlgorithmWire = []enumWire[pb.S3SseAlgorithm, S3SseAlgorithm]{
	{pb.S3SseAlgorithm_S3_SSE_ALGORITHM_AES256, S3SseAlgorithmAES256},
	{pb.S3SseAlgorithm_S3_SSE_ALGORITHM_KMS, S3SseAlgorithmKMS},
}

var projectionTypeWire = []enumWire[pb.ProjectionType, ProjectionType]{
	{pb.ProjectionType_PROJECTION_TYPE_ALL, ProjectionTypeAll},
	{pb.ProjectionType_PROJECTION_TYPE_KEYS_ONLY, ProjectionTypeKeysOnly},
	{pb.ProjectionType_PROJECTION_TYPE_INCLUDE, ProjectionTypeInclude},
}

var keyTypeWire = []enumWire[pb.KeyType, KeyType]{
	{pb.KeyType_KEY_TYPE_HASH, KeyTypeHash},
	{pb.KeyType_KEY_TYPE_RANGE, KeyTypeRange},
}

var scalarAttributeTypeWire = []enumWire[pb.ScalarAttributeType, ScalarAttributeType]{
	{pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_S, ScalarAttributeTypeS},
	{pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_N, ScalarAttributeTypeN},
	{pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_B, ScalarAttributeTypeB},
}

var indexStatusWire = []enumWire[pb.IndexStatus, IndexStatus]{
	{pb.IndexStatus_INDEX_STATUS_CREATING, IndexStatusCreating},
	{pb.IndexStatus_INDEX_STATUS_ACTIVE, IndexStatusActive},
	{pb.IndexStatus_INDEX_STATUS_UPDATING, IndexStatusUpdating},
	{pb.IndexStatus_INDEX_STATUS_DELETING, IndexStatusDeleting},
}

var vectorDistanceFunctionWire = []enumWire[pb.VectorDistanceFunction, VectorDistanceFunction]{
	{pb.VectorDistanceFunction_VECTOR_DISTANCE_FUNCTION_COSINE, VectorDistanceFunctionCosine},
	{pb.VectorDistanceFunction_VECTOR_DISTANCE_FUNCTION_EUCLIDEAN, VectorDistanceFunctionEuclidean},
	{pb.VectorDistanceFunction_VECTOR_DISTANCE_FUNCTION_DOT_PRODUCT, VectorDistanceFunctionDotProduct},
}

var searchSchemaElementTypeWire = []enumWire[pb.SearchSchemaElementType, SearchSchemaElementType]{
	{pb.SearchSchemaElementType_SEARCH_SCHEMA_ELEMENT_TYPE_HASH, SearchSchemaElementTypeHash},
	{pb.SearchSchemaElementType_SEARCH_SCHEMA_ELEMENT_TYPE_INLINE_FILTER, SearchSchemaElementTypeInlineFilter},
}

func tableStatusToProto(s TableStatus) pb.TableStatus {
	return enumToProto(tableStatusWire, s, pb.TableStatus_TABLE_STATUS_UNSPECIFIED)
}

func protoToTableStatus(s pb.TableStatus) TableStatus {
	return enumFromProto(tableStatusWire, s)
}

func billingModeToProto(b BillingMode) pb.BillingMode {
	return enumToProto(billingModeWire, b, pb.BillingMode_BILLING_MODE_UNSPECIFIED)
}

func protoToBillingMode(b pb.BillingMode) BillingMode {
	return enumFromProto(billingModeWire, b)
}

func tableClassToProto(s TableClass) pb.TableClass {
	return enumToProto(tableClassWire, s, pb.TableClass_TABLE_CLASS_UNSPECIFIED)
}

func protoToTableClass(s pb.TableClass) TableClass {
	return enumFromProto(tableClassWire, s)
}

func sseStatusToProto(s SSEStatus) pb.SSEStatus {
	return enumToProto(sseStatusWire, s, pb.SSEStatus_SSE_STATUS_UNSPECIFIED)
}

func protoToSSEStatus(s pb.SSEStatus) SSEStatus {
	return enumFromProto(sseStatusWire, s)
}

func sseTypeToProto(s SSEType) pb.SSEType {
	return enumToProto(sseTypeWire, s, pb.SSEType_SSE_TYPE_UNSPECIFIED)
}

func protoToSSEType(s pb.SSEType) SSEType {
	return enumFromProto(sseTypeWire, s)
}

func destinationStatusToProto(s DestinationStatus) pb.DestinationStatus {
	return enumToProto(destinationStatusWire, s, pb.DestinationStatus_DESTINATION_STATUS_UNSPECIFIED)
}

func protoToDestinationStatus(s pb.DestinationStatus) DestinationStatus {
	return enumFromProto(destinationStatusWire, s)
}

func contributorInsightsModeToProto(s ContributorInsightsMode) pb.ContributorInsightsMode {
	return enumToProto(contributorInsightsModeWire, s, pb.ContributorInsightsMode_CONTRIBUTOR_INSIGHTS_MODE_UNSPECIFIED)
}

func protoToContributorInsightsMode(s pb.ContributorInsightsMode) ContributorInsightsMode {
	return enumFromProto(contributorInsightsModeWire, s)
}

func approximateCreationDateTimePrecisionToProto(s ApproximateCreationDateTimePrecision) pb.ApproximateCreationDateTimePrecision {
	return enumToProto(acdtPrecisionWire, s, pb.ApproximateCreationDateTimePrecision_APPROXIMATE_CREATION_DATE_TIME_PRECISION_UNSPECIFIED)
}

func protoToApproximateCreationDateTimePrecision(s pb.ApproximateCreationDateTimePrecision) ApproximateCreationDateTimePrecision {
	return enumFromProto(acdtPrecisionWire, s)
}

func streamViewTypeToProto(s StreamViewType) pb.StreamViewType {
	return enumToProto(streamViewTypeWire, s, pb.StreamViewType_STREAM_VIEW_TYPE_UNSPECIFIED)
}

func protoToStreamViewType(s pb.StreamViewType) StreamViewType {
	return enumFromProto(streamViewTypeWire, s)
}

func ttlStatusToProto(s TTLStatus) pb.TTLStatus {
	return enumToProto(ttlStatusWire, s, pb.TTLStatus_TTL_STATUS_UNSPECIFIED)
}

func protoToTTLStatus(s pb.TTLStatus) TTLStatus {
	return enumFromProto(ttlStatusWire, s)
}

func pointInTimeRecoveryStatusToProto(s PointInTimeRecoveryStatus) pb.PointInTimeRecoveryStatus {
	return enumToProto(pitrStatusWire, s, pb.PointInTimeRecoveryStatus_POINT_IN_TIME_RECOVERY_STATUS_UNSPECIFIED)
}

func protoToPointInTimeRecoveryStatus(s pb.PointInTimeRecoveryStatus) PointInTimeRecoveryStatus {
	return enumFromProto(pitrStatusWire, s)
}

func backupStatusToProto(s BackupStatus) pb.BackupStatus {
	return enumToProto(backupStatusWire, s, pb.BackupStatus_BACKUP_STATUS_UNSPECIFIED)
}

func protoToBackupStatus(s pb.BackupStatus) BackupStatus {
	return enumFromProto(backupStatusWire, s)
}

func backupTypeToProto(t BackupType) pb.BackupType {
	return enumToProto(backupTypeWire, t, pb.BackupType_BACKUP_TYPE_UNSPECIFIED)
}

func protoToBackupType(t pb.BackupType) BackupType {
	return enumFromProto(backupTypeWire, t)
}

func globalTableStatusToProto(s GlobalTableStatus) pb.GlobalTableStatus {
	return enumToProto(globalTableStatusWire, s, pb.GlobalTableStatus_GLOBAL_TABLE_STATUS_UNSPECIFIED)
}

func protoToGlobalTableStatus(s pb.GlobalTableStatus) GlobalTableStatus {
	return enumFromProto(globalTableStatusWire, s)
}

func replicaStatusToProto(s ReplicaStatus) pb.ReplicaStatus {
	return enumToProto(replicaStatusWire, s, pb.ReplicaStatus_REPLICA_STATUS_UNSPECIFIED)
}

func protoToReplicaStatus(s pb.ReplicaStatus) ReplicaStatus {
	return enumFromProto(replicaStatusWire, s)
}

func importStatusToProto(s ImportStatus) pb.ImportStatus {
	return enumToProto(importStatusWire, s, pb.ImportStatus_IMPORT_STATUS_UNSPECIFIED)
}

func protoToImportStatus(s pb.ImportStatus) ImportStatus {
	return enumFromProto(importStatusWire, s)
}

func exportStatusToProto(s ExportStatus) pb.ExportStatus {
	return enumToProto(exportStatusWire, s, pb.ExportStatus_EXPORT_STATUS_UNSPECIFIED)
}

func protoToExportStatus(s pb.ExportStatus) ExportStatus {
	return enumFromProto(exportStatusWire, s)
}

func inputFormatToProto(s InputFormat) pb.InputFormat {
	return enumToProto(inputFormatWire, s, pb.InputFormat_INPUT_FORMAT_UNSPECIFIED)
}

func protoToInputFormat(s pb.InputFormat) InputFormat {
	return enumFromProto(inputFormatWire, s)
}

func inputCompressionTypeToProto(s InputCompressionType) pb.InputCompressionType {
	return enumToProto(inputCompressionTypeWire, s, pb.InputCompressionType_INPUT_COMPRESSION_TYPE_UNSPECIFIED)
}

func protoToInputCompressionType(s pb.InputCompressionType) InputCompressionType {
	return enumFromProto(inputCompressionTypeWire, s)
}

func exportFormatToProto(s ExportFormat) pb.ExportFormat {
	return enumToProto(exportFormatWire, s, pb.ExportFormat_EXPORT_FORMAT_UNSPECIFIED)
}

func protoToExportFormat(s pb.ExportFormat) ExportFormat {
	return enumFromProto(exportFormatWire, s)
}

func exportTypeToProto(s ExportType) pb.ExportType {
	return enumToProto(exportTypeWire, s, pb.ExportType_EXPORT_TYPE_UNSPECIFIED)
}

func protoToExportType(s pb.ExportType) ExportType {
	return enumFromProto(exportTypeWire, s)
}

func exportViewTypeToProto(s ExportViewType) pb.ExportViewType {
	return enumToProto(exportViewTypeWire, s, pb.ExportViewType_EXPORT_VIEW_TYPE_UNSPECIFIED)
}

func protoToExportViewType(s pb.ExportViewType) ExportViewType {
	return enumFromProto(exportViewTypeWire, s)
}

func s3SseAlgorithmToProto(s S3SseAlgorithm) pb.S3SseAlgorithm {
	return enumToProto(s3SseAlgorithmWire, s, pb.S3SseAlgorithm_S3_SSE_ALGORITHM_UNSPECIFIED)
}

func protoToS3SseAlgorithm(s pb.S3SseAlgorithm) S3SseAlgorithm {
	return enumFromProto(s3SseAlgorithmWire, s)
}

func projectionTypeToProto(s ProjectionType) pb.ProjectionType {
	return enumToProto(projectionTypeWire, s, pb.ProjectionType_PROJECTION_TYPE_UNSPECIFIED)
}

func protoToProjectionType(s pb.ProjectionType) ProjectionType {
	return enumFromProto(projectionTypeWire, s)
}

func keyTypeToProto(k KeyType) pb.KeyType {
	return enumToProto(keyTypeWire, k, pb.KeyType_KEY_TYPE_UNSPECIFIED)
}

func protoToKeyType(k pb.KeyType) KeyType {
	return enumFromProto(keyTypeWire, k)
}

func scalarAttributeTypeToProto(s ScalarAttributeType) pb.ScalarAttributeType {
	return enumToProto(scalarAttributeTypeWire, s, pb.ScalarAttributeType_SCALAR_ATTRIBUTE_TYPE_UNSPECIFIED)
}

func protoToScalarAttributeType(s pb.ScalarAttributeType) ScalarAttributeType {
	return enumFromProto(scalarAttributeTypeWire, s)
}

func indexStatusToProto(s IndexStatus) pb.IndexStatus {
	return enumToProto(indexStatusWire, s, pb.IndexStatus_INDEX_STATUS_UNSPECIFIED)
}

func protoToIndexStatus(s pb.IndexStatus) IndexStatus {
	return enumFromProto(indexStatusWire, s)
}

func vectorDistanceFunctionToProto(s VectorDistanceFunction) pb.VectorDistanceFunction {
	return enumToProto(vectorDistanceFunctionWire, s, pb.VectorDistanceFunction_VECTOR_DISTANCE_FUNCTION_UNSPECIFIED)
}

func protoToVectorDistanceFunction(s pb.VectorDistanceFunction) VectorDistanceFunction {
	return enumFromProto(vectorDistanceFunctionWire, s)
}

func searchSchemaElementTypeToProto(s SearchSchemaElementType) pb.SearchSchemaElementType {
	return enumToProto(searchSchemaElementTypeWire, s, pb.SearchSchemaElementType_SEARCH_SCHEMA_ELEMENT_TYPE_UNSPECIFIED)
}

func protoToSearchSchemaElementType(s pb.SearchSchemaElementType) SearchSchemaElementType {
	return enumFromProto(searchSchemaElementTypeWire, s)
}
