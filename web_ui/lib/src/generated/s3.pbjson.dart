// This is a generated file - do not edit.
//
// Generated from s3.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use analyticsS3ExportFileFormatDescriptor instead')
const AnalyticsS3ExportFileFormat$json = {
  '1': 'AnalyticsS3ExportFileFormat',
  '2': [
    {'1': 'ANALYTICS_S3_EXPORT_FILE_FORMAT_CSV', '2': 0},
  ],
};

/// Descriptor for `AnalyticsS3ExportFileFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List analyticsS3ExportFileFormatDescriptor =
    $convert.base64Decode(
        'ChtBbmFseXRpY3NTM0V4cG9ydEZpbGVGb3JtYXQSJwojQU5BTFlUSUNTX1MzX0VYUE9SVF9GSU'
        'xFX0ZPUk1BVF9DU1YQAA==');

@$core.Deprecated('Use archiveStatusDescriptor instead')
const ArchiveStatus$json = {
  '1': 'ArchiveStatus',
  '2': [
    {'1': 'ARCHIVE_STATUS_DEEP_ARCHIVE_ACCESS', '2': 0},
    {'1': 'ARCHIVE_STATUS_ARCHIVE_ACCESS', '2': 1},
  ],
};

/// Descriptor for `ArchiveStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List archiveStatusDescriptor = $convert.base64Decode(
    'Cg1BcmNoaXZlU3RhdHVzEiYKIkFSQ0hJVkVfU1RBVFVTX0RFRVBfQVJDSElWRV9BQ0NFU1MQAB'
    'IhCh1BUkNISVZFX1NUQVRVU19BUkNISVZFX0FDQ0VTUxAB');

@$core.Deprecated('Use bucketAbacStatusDescriptor instead')
const BucketAbacStatus$json = {
  '1': 'BucketAbacStatus',
  '2': [
    {'1': 'BUCKET_ABAC_STATUS_DISABLED', '2': 0},
    {'1': 'BUCKET_ABAC_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `BucketAbacStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bucketAbacStatusDescriptor = $convert.base64Decode(
    'ChBCdWNrZXRBYmFjU3RhdHVzEh8KG0JVQ0tFVF9BQkFDX1NUQVRVU19ESVNBQkxFRBAAEh4KGk'
    'JVQ0tFVF9BQkFDX1NUQVRVU19FTkFCTEVEEAE=');

@$core.Deprecated('Use bucketAccelerateStatusDescriptor instead')
const BucketAccelerateStatus$json = {
  '1': 'BucketAccelerateStatus',
  '2': [
    {'1': 'BUCKET_ACCELERATE_STATUS_SUSPENDED', '2': 0},
    {'1': 'BUCKET_ACCELERATE_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `BucketAccelerateStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bucketAccelerateStatusDescriptor =
    $convert.base64Decode(
        'ChZCdWNrZXRBY2NlbGVyYXRlU3RhdHVzEiYKIkJVQ0tFVF9BQ0NFTEVSQVRFX1NUQVRVU19TVV'
        'NQRU5ERUQQABIkCiBCVUNLRVRfQUNDRUxFUkFURV9TVEFUVVNfRU5BQkxFRBAB');

@$core.Deprecated('Use bucketCannedACLDescriptor instead')
const BucketCannedACL$json = {
  '1': 'BucketCannedACL',
  '2': [
    {'1': 'BUCKET_CANNED_A_C_L_AUTHENTICATED_READ', '2': 0},
    {'1': 'BUCKET_CANNED_A_C_L_PRIVATE', '2': 1},
    {'1': 'BUCKET_CANNED_A_C_L_PUBLIC_READ', '2': 2},
    {'1': 'BUCKET_CANNED_A_C_L_PUBLIC_READ_WRITE', '2': 3},
  ],
};

/// Descriptor for `BucketCannedACL`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bucketCannedACLDescriptor = $convert.base64Decode(
    'Cg9CdWNrZXRDYW5uZWRBQ0wSKgomQlVDS0VUX0NBTk5FRF9BX0NfTF9BVVRIRU5USUNBVEVEX1'
    'JFQUQQABIfChtCVUNLRVRfQ0FOTkVEX0FfQ19MX1BSSVZBVEUQARIjCh9CVUNLRVRfQ0FOTkVE'
    'X0FfQ19MX1BVQkxJQ19SRUFEEAISKQolQlVDS0VUX0NBTk5FRF9BX0NfTF9QVUJMSUNfUkVBRF'
    '9XUklURRAD');

@$core.Deprecated('Use bucketLocationConstraintDescriptor instead')
const BucketLocationConstraint$json = {
  '1': 'BucketLocationConstraint',
  '2': [
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_3', '2': 0},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU_WEST_3', '2': 1},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_US_GOV_EAST_1', '2': 2},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_1', '2': 3},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_IL_CENTRAL_1', '2': 4},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU_NORTH_1', '2': 5},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_CN_NORTH_1', '2': 6},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_5', '2': 7},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_2', '2': 8},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_US_WEST_1', '2': 9},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTH_1', '2': 10},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_2', '2': 11},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU_SOUTH_1', '2': 12},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU_WEST_2', '2': 13},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_US_WEST_2', '2': 14},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_ME_CENTRAL_1', '2': 15},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_CN_NORTHWEST_1', '2': 16},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU_SOUTH_2', '2': 17},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_EAST_1', '2': 18},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU_CENTRAL_1', '2': 19},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_4', '2': 20},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_CA_CENTRAL_1', '2': 21},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_SA_EAST_1', '2': 22},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTH_2', '2': 23},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_US_GOV_WEST_1', '2': 24},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU_CENTRAL_2', '2': 25},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU', '2': 26},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_1', '2': 27},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AF_SOUTH_1', '2': 28},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_ME_SOUTH_1', '2': 29},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_EU_WEST_1', '2': 30},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_3', '2': 31},
    {'1': 'BUCKET_LOCATION_CONSTRAINT_US_EAST_2', '2': 32},
  ],
};

/// Descriptor for `BucketLocationConstraint`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bucketLocationConstraintDescriptor = $convert.base64Decode(
    'ChhCdWNrZXRMb2NhdGlvbkNvbnN0cmFpbnQSLQopQlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTl'
    'RfQVBfU09VVEhFQVNUXzMQABIoCiRCVUNLRVRfTE9DQVRJT05fQ09OU1RSQUlOVF9FVV9XRVNU'
    'XzMQARIsCihCVUNLRVRfTE9DQVRJT05fQ09OU1RSQUlOVF9VU19HT1ZfRUFTVF8xEAISLQopQl'
    'VDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfQVBfTk9SVEhFQVNUXzEQAxIrCidCVUNLRVRfTE9D'
    'QVRJT05fQ09OU1RSQUlOVF9JTF9DRU5UUkFMXzEQBBIpCiVCVUNLRVRfTE9DQVRJT05fQ09OU1'
    'RSQUlOVF9FVV9OT1JUSF8xEAUSKQolQlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfQ05fTk9S'
    'VEhfMRAGEi0KKUJVQ0tFVF9MT0NBVElPTl9DT05TVFJBSU5UX0FQX1NPVVRIRUFTVF81EAcSLQ'
    'opQlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfQVBfTk9SVEhFQVNUXzIQCBIoCiRCVUNLRVRf'
    'TE9DQVRJT05fQ09OU1RSQUlOVF9VU19XRVNUXzEQCRIpCiVCVUNLRVRfTE9DQVRJT05fQ09OU1'
    'RSQUlOVF9BUF9TT1VUSF8xEAoSLQopQlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfQVBfU09V'
    'VEhFQVNUXzIQCxIpCiVCVUNLRVRfTE9DQVRJT05fQ09OU1RSQUlOVF9FVV9TT1VUSF8xEAwSKA'
    'okQlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfRVVfV0VTVF8yEA0SKAokQlVDS0VUX0xPQ0FU'
    'SU9OX0NPTlNUUkFJTlRfVVNfV0VTVF8yEA4SKwonQlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTl'
    'RfTUVfQ0VOVFJBTF8xEA8SLQopQlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfQ05fTk9SVEhX'
    'RVNUXzEQEBIpCiVCVUNLRVRfTE9DQVRJT05fQ09OU1RSQUlOVF9FVV9TT1VUSF8yEBESKAokQl'
    'VDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfQVBfRUFTVF8xEBISKwonQlVDS0VUX0xPQ0FUSU9O'
    'X0NPTlNUUkFJTlRfRVVfQ0VOVFJBTF8xEBMSLQopQlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTl'
    'RfQVBfU09VVEhFQVNUXzQQFBIrCidCVUNLRVRfTE9DQVRJT05fQ09OU1RSQUlOVF9DQV9DRU5U'
    'UkFMXzEQFRIoCiRCVUNLRVRfTE9DQVRJT05fQ09OU1RSQUlOVF9TQV9FQVNUXzEQFhIpCiVCVU'
    'NLRVRfTE9DQVRJT05fQ09OU1RSQUlOVF9BUF9TT1VUSF8yEBcSLAooQlVDS0VUX0xPQ0FUSU9O'
    'X0NPTlNUUkFJTlRfVVNfR09WX1dFU1RfMRAYEisKJ0JVQ0tFVF9MT0NBVElPTl9DT05TVFJBSU'
    '5UX0VVX0NFTlRSQUxfMhAZEiEKHUJVQ0tFVF9MT0NBVElPTl9DT05TVFJBSU5UX0VVEBoSLQop'
    'QlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfQVBfU09VVEhFQVNUXzEQGxIpCiVCVUNLRVRfTE'
    '9DQVRJT05fQ09OU1RSQUlOVF9BRl9TT1VUSF8xEBwSKQolQlVDS0VUX0xPQ0FUSU9OX0NPTlNU'
    'UkFJTlRfTUVfU09VVEhfMRAdEigKJEJVQ0tFVF9MT0NBVElPTl9DT05TVFJBSU5UX0VVX1dFU1'
    'RfMRAeEi0KKUJVQ0tFVF9MT0NBVElPTl9DT05TVFJBSU5UX0FQX05PUlRIRUFTVF8zEB8SKAok'
    'QlVDS0VUX0xPQ0FUSU9OX0NPTlNUUkFJTlRfVVNfRUFTVF8yECA=');

@$core.Deprecated('Use bucketLogsPermissionDescriptor instead')
const BucketLogsPermission$json = {
  '1': 'BucketLogsPermission',
  '2': [
    {'1': 'BUCKET_LOGS_PERMISSION_READ', '2': 0},
    {'1': 'BUCKET_LOGS_PERMISSION_WRITE', '2': 1},
    {'1': 'BUCKET_LOGS_PERMISSION_FULL_CONTROL', '2': 2},
  ],
};

/// Descriptor for `BucketLogsPermission`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bucketLogsPermissionDescriptor = $convert.base64Decode(
    'ChRCdWNrZXRMb2dzUGVybWlzc2lvbhIfChtCVUNLRVRfTE9HU19QRVJNSVNTSU9OX1JFQUQQAB'
    'IgChxCVUNLRVRfTE9HU19QRVJNSVNTSU9OX1dSSVRFEAESJwojQlVDS0VUX0xPR1NfUEVSTUlT'
    'U0lPTl9GVUxMX0NPTlRST0wQAg==');

@$core.Deprecated('Use bucketTypeDescriptor instead')
const BucketType$json = {
  '1': 'BucketType',
  '2': [
    {'1': 'BUCKET_TYPE_DIRECTORY', '2': 0},
  ],
};

/// Descriptor for `BucketType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bucketTypeDescriptor = $convert
    .base64Decode('CgpCdWNrZXRUeXBlEhkKFUJVQ0tFVF9UWVBFX0RJUkVDVE9SWRAA');

@$core.Deprecated('Use bucketVersioningStatusDescriptor instead')
const BucketVersioningStatus$json = {
  '1': 'BucketVersioningStatus',
  '2': [
    {'1': 'BUCKET_VERSIONING_STATUS_SUSPENDED', '2': 0},
    {'1': 'BUCKET_VERSIONING_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `BucketVersioningStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bucketVersioningStatusDescriptor =
    $convert.base64Decode(
        'ChZCdWNrZXRWZXJzaW9uaW5nU3RhdHVzEiYKIkJVQ0tFVF9WRVJTSU9OSU5HX1NUQVRVU19TVV'
        'NQRU5ERUQQABIkCiBCVUNLRVRfVkVSU0lPTklOR19TVEFUVVNfRU5BQkxFRBAB');

@$core.Deprecated('Use checksumAlgorithmDescriptor instead')
const ChecksumAlgorithm$json = {
  '1': 'ChecksumAlgorithm',
  '2': [
    {'1': 'CHECKSUM_ALGORITHM_SHA256', '2': 0},
    {'1': 'CHECKSUM_ALGORITHM_SHA1', '2': 1},
    {'1': 'CHECKSUM_ALGORITHM_CRC32', '2': 2},
    {'1': 'CHECKSUM_ALGORITHM_CRC64NVME', '2': 3},
    {'1': 'CHECKSUM_ALGORITHM_CRC32C', '2': 4},
  ],
};

/// Descriptor for `ChecksumAlgorithm`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List checksumAlgorithmDescriptor = $convert.base64Decode(
    'ChFDaGVja3N1bUFsZ29yaXRobRIdChlDSEVDS1NVTV9BTEdPUklUSE1fU0hBMjU2EAASGwoXQ0'
    'hFQ0tTVU1fQUxHT1JJVEhNX1NIQTEQARIcChhDSEVDS1NVTV9BTEdPUklUSE1fQ1JDMzIQAhIg'
    'ChxDSEVDS1NVTV9BTEdPUklUSE1fQ1JDNjROVk1FEAMSHQoZQ0hFQ0tTVU1fQUxHT1JJVEhNX0'
    'NSQzMyQxAE');

@$core.Deprecated('Use checksumModeDescriptor instead')
const ChecksumMode$json = {
  '1': 'ChecksumMode',
  '2': [
    {'1': 'CHECKSUM_MODE_ENABLED', '2': 0},
  ],
};

/// Descriptor for `ChecksumMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List checksumModeDescriptor = $convert
    .base64Decode('CgxDaGVja3N1bU1vZGUSGQoVQ0hFQ0tTVU1fTU9ERV9FTkFCTEVEEAA=');

@$core.Deprecated('Use checksumTypeDescriptor instead')
const ChecksumType$json = {
  '1': 'ChecksumType',
  '2': [
    {'1': 'CHECKSUM_TYPE_FULL_OBJECT', '2': 0},
    {'1': 'CHECKSUM_TYPE_COMPOSITE', '2': 1},
  ],
};

/// Descriptor for `ChecksumType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List checksumTypeDescriptor = $convert.base64Decode(
    'CgxDaGVja3N1bVR5cGUSHQoZQ0hFQ0tTVU1fVFlQRV9GVUxMX09CSkVDVBAAEhsKF0NIRUNLU1'
    'VNX1RZUEVfQ09NUE9TSVRFEAE=');

@$core.Deprecated('Use compressionTypeDescriptor instead')
const CompressionType$json = {
  '1': 'CompressionType',
  '2': [
    {'1': 'COMPRESSION_TYPE_NONE', '2': 0},
    {'1': 'COMPRESSION_TYPE_GZIP', '2': 1},
    {'1': 'COMPRESSION_TYPE_BZIP2', '2': 2},
  ],
};

/// Descriptor for `CompressionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List compressionTypeDescriptor = $convert.base64Decode(
    'Cg9Db21wcmVzc2lvblR5cGUSGQoVQ09NUFJFU1NJT05fVFlQRV9OT05FEAASGQoVQ09NUFJFU1'
    'NJT05fVFlQRV9HWklQEAESGgoWQ09NUFJFU1NJT05fVFlQRV9CWklQMhAC');

@$core.Deprecated('Use dataRedundancyDescriptor instead')
const DataRedundancy$json = {
  '1': 'DataRedundancy',
  '2': [
    {'1': 'DATA_REDUNDANCY_SINGLEAVAILABILITYZONE', '2': 0},
    {'1': 'DATA_REDUNDANCY_SINGLELOCALZONE', '2': 1},
  ],
};

/// Descriptor for `DataRedundancy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dataRedundancyDescriptor = $convert.base64Decode(
    'Cg5EYXRhUmVkdW5kYW5jeRIqCiZEQVRBX1JFRFVOREFOQ1lfU0lOR0xFQVZBSUxBQklMSVRZWk'
    '9ORRAAEiMKH0RBVEFfUkVEVU5EQU5DWV9TSU5HTEVMT0NBTFpPTkUQAQ==');

@$core.Deprecated('Use deleteMarkerReplicationStatusDescriptor instead')
const DeleteMarkerReplicationStatus$json = {
  '1': 'DeleteMarkerReplicationStatus',
  '2': [
    {'1': 'DELETE_MARKER_REPLICATION_STATUS_DISABLED', '2': 0},
    {'1': 'DELETE_MARKER_REPLICATION_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `DeleteMarkerReplicationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deleteMarkerReplicationStatusDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVNYXJrZXJSZXBsaWNhdGlvblN0YXR1cxItCilERUxFVEVfTUFSS0VSX1JFUExJQ0'
        'FUSU9OX1NUQVRVU19ESVNBQkxFRBAAEiwKKERFTEVURV9NQVJLRVJfUkVQTElDQVRJT05fU1RB'
        'VFVTX0VOQUJMRUQQAQ==');

@$core.Deprecated('Use encodingTypeDescriptor instead')
const EncodingType$json = {
  '1': 'EncodingType',
  '2': [
    {'1': 'ENCODING_TYPE_URL', '2': 0},
  ],
};

/// Descriptor for `EncodingType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List encodingTypeDescriptor = $convert
    .base64Decode('CgxFbmNvZGluZ1R5cGUSFQoRRU5DT0RJTkdfVFlQRV9VUkwQAA==');

@$core.Deprecated('Use encryptionTypeDescriptor instead')
const EncryptionType$json = {
  '1': 'EncryptionType',
  '2': [
    {'1': 'ENCRYPTION_TYPE_NONE', '2': 0},
    {'1': 'ENCRYPTION_TYPE_SSE_C', '2': 1},
  ],
};

/// Descriptor for `EncryptionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List encryptionTypeDescriptor = $convert.base64Decode(
    'Cg5FbmNyeXB0aW9uVHlwZRIYChRFTkNSWVBUSU9OX1RZUEVfTk9ORRAAEhkKFUVOQ1JZUFRJT0'
    '5fVFlQRV9TU0VfQxAB');

@$core.Deprecated('Use eventDescriptor instead')
const Event$json = {
  '1': 'Event',
  '2': [
    {'1': 'EVENT_S3_OBJECTRESTORE_', '2': 0},
    {'1': 'EVENT_S3_OBJECTTAGGING_', '2': 1},
    {'1': 'EVENT_S3_REPLICATION_OPERATIONFAILEDREPLICATION', '2': 2},
    {'1': 'EVENT_S3_OBJECTCREATED_COPY', '2': 3},
    {'1': 'EVENT_S3_OBJECTCREATED_POST', '2': 4},
    {'1': 'EVENT_S3_REDUCEDREDUNDANCYLOSTOBJECT', '2': 5},
    {'1': 'EVENT_S3_OBJECTCREATED_PUT', '2': 6},
    {'1': 'EVENT_S3_OBJECTRESTORE_COMPLETED', '2': 7},
    {'1': 'EVENT_S3_LIFECYCLEEXPIRATION_DELETE', '2': 8},
    {'1': 'EVENT_S3_OBJECTTAGGING_PUT', '2': 9},
    {'1': 'EVENT_S3_OBJECTREMOVED_DELETEMARKERCREATED', '2': 10},
    {'1': 'EVENT_S3_OBJECTTAGGING_DELETE', '2': 11},
    {'1': 'EVENT_S3_REPLICATION_OPERATIONMISSEDTHRESHOLD', '2': 12},
    {'1': 'EVENT_S3_LIFECYCLETRANSITION', '2': 13},
    {'1': 'EVENT_S3_OBJECTCREATED_COMPLETEMULTIPARTUPLOAD', '2': 14},
    {'1': 'EVENT_S3_INTELLIGENTTIERING', '2': 15},
    {'1': 'EVENT_S3_OBJECTCREATED_', '2': 16},
    {'1': 'EVENT_S3_REPLICATION_OPERATIONREPLICATEDAFTERTHRESHOLD', '2': 17},
    {'1': 'EVENT_S3_OBJECTRESTORE_DELETE', '2': 18},
    {'1': 'EVENT_S3_LIFECYCLEEXPIRATION_', '2': 19},
    {'1': 'EVENT_S3_REPLICATION_OPERATIONNOTTRACKED', '2': 20},
    {'1': 'EVENT_S3_LIFECYCLEEXPIRATION_DELETEMARKERCREATED', '2': 21},
    {'1': 'EVENT_S3_REPLICATION_', '2': 22},
    {'1': 'EVENT_S3_OBJECTREMOVED_', '2': 23},
    {'1': 'EVENT_S3_OBJECTRESTORE_POST', '2': 24},
    {'1': 'EVENT_S3_OBJECTACL_PUT', '2': 25},
    {'1': 'EVENT_S3_OBJECTREMOVED_DELETE', '2': 26},
  ],
};

/// Descriptor for `Event`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventDescriptor = $convert.base64Decode(
    'CgVFdmVudBIbChdFVkVOVF9TM19PQkpFQ1RSRVNUT1JFXxAAEhsKF0VWRU5UX1MzX09CSkVDVF'
    'RBR0dJTkdfEAESMwovRVZFTlRfUzNfUkVQTElDQVRJT05fT1BFUkFUSU9ORkFJTEVEUkVQTElD'
    'QVRJT04QAhIfChtFVkVOVF9TM19PQkpFQ1RDUkVBVEVEX0NPUFkQAxIfChtFVkVOVF9TM19PQk'
    'pFQ1RDUkVBVEVEX1BPU1QQBBIoCiRFVkVOVF9TM19SRURVQ0VEUkVEVU5EQU5DWUxPU1RPQkpF'
    'Q1QQBRIeChpFVkVOVF9TM19PQkpFQ1RDUkVBVEVEX1BVVBAGEiQKIEVWRU5UX1MzX09CSkVDVF'
    'JFU1RPUkVfQ09NUExFVEVEEAcSJwojRVZFTlRfUzNfTElGRUNZQ0xFRVhQSVJBVElPTl9ERUxF'
    'VEUQCBIeChpFVkVOVF9TM19PQkpFQ1RUQUdHSU5HX1BVVBAJEi4KKkVWRU5UX1MzX09CSkVDVF'
    'JFTU9WRURfREVMRVRFTUFSS0VSQ1JFQVRFRBAKEiEKHUVWRU5UX1MzX09CSkVDVFRBR0dJTkdf'
    'REVMRVRFEAsSMQotRVZFTlRfUzNfUkVQTElDQVRJT05fT1BFUkFUSU9OTUlTU0VEVEhSRVNIT0'
    'xEEAwSIAocRVZFTlRfUzNfTElGRUNZQ0xFVFJBTlNJVElPThANEjIKLkVWRU5UX1MzX09CSkVD'
    'VENSRUFURURfQ09NUExFVEVNVUxUSVBBUlRVUExPQUQQDhIfChtFVkVOVF9TM19JTlRFTExJR0'
    'VOVFRJRVJJTkcQDxIbChdFVkVOVF9TM19PQkpFQ1RDUkVBVEVEXxAQEjoKNkVWRU5UX1MzX1JF'
    'UExJQ0FUSU9OX09QRVJBVElPTlJFUExJQ0FURURBRlRFUlRIUkVTSE9MRBAREiEKHUVWRU5UX1'
    'MzX09CSkVDVFJFU1RPUkVfREVMRVRFEBISIQodRVZFTlRfUzNfTElGRUNZQ0xFRVhQSVJBVElP'
    'Tl8QExIsCihFVkVOVF9TM19SRVBMSUNBVElPTl9PUEVSQVRJT05OT1RUUkFDS0VEEBQSNAowRV'
    'ZFTlRfUzNfTElGRUNZQ0xFRVhQSVJBVElPTl9ERUxFVEVNQVJLRVJDUkVBVEVEEBUSGQoVRVZF'
    'TlRfUzNfUkVQTElDQVRJT05fEBYSGwoXRVZFTlRfUzNfT0JKRUNUUkVNT1ZFRF8QFxIfChtFVk'
    'VOVF9TM19PQkpFQ1RSRVNUT1JFX1BPU1QQGBIaChZFVkVOVF9TM19PQkpFQ1RBQ0xfUFVUEBkS'
    'IQodRVZFTlRfUzNfT0JKRUNUUkVNT1ZFRF9ERUxFVEUQGg==');

@$core.Deprecated('Use existingObjectReplicationStatusDescriptor instead')
const ExistingObjectReplicationStatus$json = {
  '1': 'ExistingObjectReplicationStatus',
  '2': [
    {'1': 'EXISTING_OBJECT_REPLICATION_STATUS_DISABLED', '2': 0},
    {'1': 'EXISTING_OBJECT_REPLICATION_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ExistingObjectReplicationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List existingObjectReplicationStatusDescriptor =
    $convert.base64Decode(
        'Ch9FeGlzdGluZ09iamVjdFJlcGxpY2F0aW9uU3RhdHVzEi8KK0VYSVNUSU5HX09CSkVDVF9SRV'
        'BMSUNBVElPTl9TVEFUVVNfRElTQUJMRUQQABIuCipFWElTVElOR19PQkpFQ1RfUkVQTElDQVRJ'
        'T05fU1RBVFVTX0VOQUJMRUQQAQ==');

@$core.Deprecated('Use expirationStateDescriptor instead')
const ExpirationState$json = {
  '1': 'ExpirationState',
  '2': [
    {'1': 'EXPIRATION_STATE_DISABLED', '2': 0},
    {'1': 'EXPIRATION_STATE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ExpirationState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List expirationStateDescriptor = $convert.base64Decode(
    'Cg9FeHBpcmF0aW9uU3RhdGUSHQoZRVhQSVJBVElPTl9TVEFURV9ESVNBQkxFRBAAEhwKGEVYUE'
    'lSQVRJT05fU1RBVEVfRU5BQkxFRBAB');

@$core.Deprecated('Use expirationStatusDescriptor instead')
const ExpirationStatus$json = {
  '1': 'ExpirationStatus',
  '2': [
    {'1': 'EXPIRATION_STATUS_DISABLED', '2': 0},
    {'1': 'EXPIRATION_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ExpirationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List expirationStatusDescriptor = $convert.base64Decode(
    'ChBFeHBpcmF0aW9uU3RhdHVzEh4KGkVYUElSQVRJT05fU1RBVFVTX0RJU0FCTEVEEAASHQoZRV'
    'hQSVJBVElPTl9TVEFUVVNfRU5BQkxFRBAB');

@$core.Deprecated('Use expressionTypeDescriptor instead')
const ExpressionType$json = {
  '1': 'ExpressionType',
  '2': [
    {'1': 'EXPRESSION_TYPE_SQL', '2': 0},
  ],
};

/// Descriptor for `ExpressionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List expressionTypeDescriptor = $convert
    .base64Decode('Cg5FeHByZXNzaW9uVHlwZRIXChNFWFBSRVNTSU9OX1RZUEVfU1FMEAA=');

@$core.Deprecated('Use fileHeaderInfoDescriptor instead')
const FileHeaderInfo$json = {
  '1': 'FileHeaderInfo',
  '2': [
    {'1': 'FILE_HEADER_INFO_IGNORE', '2': 0},
    {'1': 'FILE_HEADER_INFO_NONE', '2': 1},
    {'1': 'FILE_HEADER_INFO_USE', '2': 2},
  ],
};

/// Descriptor for `FileHeaderInfo`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List fileHeaderInfoDescriptor = $convert.base64Decode(
    'Cg5GaWxlSGVhZGVySW5mbxIbChdGSUxFX0hFQURFUl9JTkZPX0lHTk9SRRAAEhkKFUZJTEVfSE'
    'VBREVSX0lORk9fTk9ORRABEhgKFEZJTEVfSEVBREVSX0lORk9fVVNFEAI=');

@$core.Deprecated('Use filterRuleNameDescriptor instead')
const FilterRuleName$json = {
  '1': 'FilterRuleName',
  '2': [
    {'1': 'FILTER_RULE_NAME_SUFFIX', '2': 0},
    {'1': 'FILTER_RULE_NAME_PREFIX', '2': 1},
  ],
};

/// Descriptor for `FilterRuleName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List filterRuleNameDescriptor = $convert.base64Decode(
    'Cg5GaWx0ZXJSdWxlTmFtZRIbChdGSUxURVJfUlVMRV9OQU1FX1NVRkZJWBAAEhsKF0ZJTFRFUl'
    '9SVUxFX05BTUVfUFJFRklYEAE=');

@$core.Deprecated('Use intelligentTieringAccessTierDescriptor instead')
const IntelligentTieringAccessTier$json = {
  '1': 'IntelligentTieringAccessTier',
  '2': [
    {'1': 'INTELLIGENT_TIERING_ACCESS_TIER_DEEP_ARCHIVE_ACCESS', '2': 0},
    {'1': 'INTELLIGENT_TIERING_ACCESS_TIER_ARCHIVE_ACCESS', '2': 1},
  ],
};

/// Descriptor for `IntelligentTieringAccessTier`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List intelligentTieringAccessTierDescriptor =
    $convert.base64Decode(
        'ChxJbnRlbGxpZ2VudFRpZXJpbmdBY2Nlc3NUaWVyEjcKM0lOVEVMTElHRU5UX1RJRVJJTkdfQU'
        'NDRVNTX1RJRVJfREVFUF9BUkNISVZFX0FDQ0VTUxAAEjIKLklOVEVMTElHRU5UX1RJRVJJTkdf'
        'QUNDRVNTX1RJRVJfQVJDSElWRV9BQ0NFU1MQAQ==');

@$core.Deprecated('Use intelligentTieringStatusDescriptor instead')
const IntelligentTieringStatus$json = {
  '1': 'IntelligentTieringStatus',
  '2': [
    {'1': 'INTELLIGENT_TIERING_STATUS_DISABLED', '2': 0},
    {'1': 'INTELLIGENT_TIERING_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `IntelligentTieringStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List intelligentTieringStatusDescriptor =
    $convert.base64Decode(
        'ChhJbnRlbGxpZ2VudFRpZXJpbmdTdGF0dXMSJwojSU5URUxMSUdFTlRfVElFUklOR19TVEFUVV'
        'NfRElTQUJMRUQQABImCiJJTlRFTExJR0VOVF9USUVSSU5HX1NUQVRVU19FTkFCTEVEEAE=');

@$core.Deprecated('Use inventoryConfigurationStateDescriptor instead')
const InventoryConfigurationState$json = {
  '1': 'InventoryConfigurationState',
  '2': [
    {'1': 'INVENTORY_CONFIGURATION_STATE_DISABLED', '2': 0},
    {'1': 'INVENTORY_CONFIGURATION_STATE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `InventoryConfigurationState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventoryConfigurationStateDescriptor =
    $convert.base64Decode(
        'ChtJbnZlbnRvcnlDb25maWd1cmF0aW9uU3RhdGUSKgomSU5WRU5UT1JZX0NPTkZJR1VSQVRJT0'
        '5fU1RBVEVfRElTQUJMRUQQABIpCiVJTlZFTlRPUllfQ09ORklHVVJBVElPTl9TVEFURV9FTkFC'
        'TEVEEAE=');

@$core.Deprecated('Use inventoryFormatDescriptor instead')
const InventoryFormat$json = {
  '1': 'InventoryFormat',
  '2': [
    {'1': 'INVENTORY_FORMAT_ORC', '2': 0},
    {'1': 'INVENTORY_FORMAT_CSV', '2': 1},
    {'1': 'INVENTORY_FORMAT_PARQUET', '2': 2},
  ],
};

/// Descriptor for `InventoryFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventoryFormatDescriptor = $convert.base64Decode(
    'Cg9JbnZlbnRvcnlGb3JtYXQSGAoUSU5WRU5UT1JZX0ZPUk1BVF9PUkMQABIYChRJTlZFTlRPUl'
    'lfRk9STUFUX0NTVhABEhwKGElOVkVOVE9SWV9GT1JNQVRfUEFSUVVFVBAC');

@$core.Deprecated('Use inventoryFrequencyDescriptor instead')
const InventoryFrequency$json = {
  '1': 'InventoryFrequency',
  '2': [
    {'1': 'INVENTORY_FREQUENCY_DAILY', '2': 0},
    {'1': 'INVENTORY_FREQUENCY_WEEKLY', '2': 1},
  ],
};

/// Descriptor for `InventoryFrequency`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventoryFrequencyDescriptor = $convert.base64Decode(
    'ChJJbnZlbnRvcnlGcmVxdWVuY3kSHQoZSU5WRU5UT1JZX0ZSRVFVRU5DWV9EQUlMWRAAEh4KGk'
    'lOVkVOVE9SWV9GUkVRVUVOQ1lfV0VFS0xZEAE=');

@$core.Deprecated('Use inventoryIncludedObjectVersionsDescriptor instead')
const InventoryIncludedObjectVersions$json = {
  '1': 'InventoryIncludedObjectVersions',
  '2': [
    {'1': 'INVENTORY_INCLUDED_OBJECT_VERSIONS_CURRENT', '2': 0},
    {'1': 'INVENTORY_INCLUDED_OBJECT_VERSIONS_ALL', '2': 1},
  ],
};

/// Descriptor for `InventoryIncludedObjectVersions`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventoryIncludedObjectVersionsDescriptor =
    $convert.base64Decode(
        'Ch9JbnZlbnRvcnlJbmNsdWRlZE9iamVjdFZlcnNpb25zEi4KKklOVkVOVE9SWV9JTkNMVURFRF'
        '9PQkpFQ1RfVkVSU0lPTlNfQ1VSUkVOVBAAEioKJklOVkVOVE9SWV9JTkNMVURFRF9PQkpFQ1Rf'
        'VkVSU0lPTlNfQUxMEAE=');

@$core.Deprecated('Use inventoryOptionalFieldDescriptor instead')
const InventoryOptionalField$json = {
  '1': 'InventoryOptionalField',
  '2': [
    {'1': 'INVENTORY_OPTIONAL_FIELD_LIFECYCLEEXPIRATIONDATE', '2': 0},
    {'1': 'INVENTORY_OPTIONAL_FIELD_BUCKETKEYSTATUS', '2': 1},
    {'1': 'INVENTORY_OPTIONAL_FIELD_STORAGECLASS', '2': 2},
    {'1': 'INVENTORY_OPTIONAL_FIELD_CHECKSUMALGORITHM', '2': 3},
    {'1': 'INVENTORY_OPTIONAL_FIELD_REPLICATIONSTATUS', '2': 4},
    {'1': 'INVENTORY_OPTIONAL_FIELD_SIZE', '2': 5},
    {'1': 'INVENTORY_OPTIONAL_FIELD_OBJECTLOCKLEGALHOLDSTATUS', '2': 6},
    {'1': 'INVENTORY_OPTIONAL_FIELD_OBJECTOWNER', '2': 7},
    {'1': 'INVENTORY_OPTIONAL_FIELD_OBJECTACCESSCONTROLLIST', '2': 8},
    {'1': 'INVENTORY_OPTIONAL_FIELD_ENCRYPTIONSTATUS', '2': 9},
    {'1': 'INVENTORY_OPTIONAL_FIELD_OBJECTLOCKMODE', '2': 10},
    {'1': 'INVENTORY_OPTIONAL_FIELD_ETAG', '2': 11},
    {'1': 'INVENTORY_OPTIONAL_FIELD_ISMULTIPARTUPLOADED', '2': 12},
    {'1': 'INVENTORY_OPTIONAL_FIELD_INTELLIGENTTIERINGACCESSTIER', '2': 13},
    {'1': 'INVENTORY_OPTIONAL_FIELD_OBJECTLOCKRETAINUNTILDATE', '2': 14},
    {'1': 'INVENTORY_OPTIONAL_FIELD_LASTMODIFIEDDATE', '2': 15},
  ],
};

/// Descriptor for `InventoryOptionalField`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventoryOptionalFieldDescriptor = $convert.base64Decode(
    'ChZJbnZlbnRvcnlPcHRpb25hbEZpZWxkEjQKMElOVkVOVE9SWV9PUFRJT05BTF9GSUVMRF9MSU'
    'ZFQ1lDTEVFWFBJUkFUSU9OREFURRAAEiwKKElOVkVOVE9SWV9PUFRJT05BTF9GSUVMRF9CVUNL'
    'RVRLRVlTVEFUVVMQARIpCiVJTlZFTlRPUllfT1BUSU9OQUxfRklFTERfU1RPUkFHRUNMQVNTEA'
    'ISLgoqSU5WRU5UT1JZX09QVElPTkFMX0ZJRUxEX0NIRUNLU1VNQUxHT1JJVEhNEAMSLgoqSU5W'
    'RU5UT1JZX09QVElPTkFMX0ZJRUxEX1JFUExJQ0FUSU9OU1RBVFVTEAQSIQodSU5WRU5UT1JZX0'
    '9QVElPTkFMX0ZJRUxEX1NJWkUQBRI2CjJJTlZFTlRPUllfT1BUSU9OQUxfRklFTERfT0JKRUNU'
    'TE9DS0xFR0FMSE9MRFNUQVRVUxAGEigKJElOVkVOVE9SWV9PUFRJT05BTF9GSUVMRF9PQkpFQ1'
    'RPV05FUhAHEjQKMElOVkVOVE9SWV9PUFRJT05BTF9GSUVMRF9PQkpFQ1RBQ0NFU1NDT05UUk9M'
    'TElTVBAIEi0KKUlOVkVOVE9SWV9PUFRJT05BTF9GSUVMRF9FTkNSWVBUSU9OU1RBVFVTEAkSKw'
    'onSU5WRU5UT1JZX09QVElPTkFMX0ZJRUxEX09CSkVDVExPQ0tNT0RFEAoSIQodSU5WRU5UT1JZ'
    'X09QVElPTkFMX0ZJRUxEX0VUQUcQCxIwCixJTlZFTlRPUllfT1BUSU9OQUxfRklFTERfSVNNVU'
    'xUSVBBUlRVUExPQURFRBAMEjkKNUlOVkVOVE9SWV9PUFRJT05BTF9GSUVMRF9JTlRFTExJR0VO'
    'VFRJRVJJTkdBQ0NFU1NUSUVSEA0SNgoySU5WRU5UT1JZX09QVElPTkFMX0ZJRUxEX09CSkVDVE'
    'xPQ0tSRVRBSU5VTlRJTERBVEUQDhItCilJTlZFTlRPUllfT1BUSU9OQUxfRklFTERfTEFTVE1P'
    'RElGSUVEREFURRAP');

@$core.Deprecated('Use jSONTypeDescriptor instead')
const JSONType$json = {
  '1': 'JSONType',
  '2': [
    {'1': 'J_S_O_N_TYPE_LINES', '2': 0},
    {'1': 'J_S_O_N_TYPE_DOCUMENT', '2': 1},
  ],
};

/// Descriptor for `JSONType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List jSONTypeDescriptor = $convert.base64Decode(
    'CghKU09OVHlwZRIWChJKX1NfT19OX1RZUEVfTElORVMQABIZChVKX1NfT19OX1RZUEVfRE9DVU'
    '1FTlQQAQ==');

@$core.Deprecated('Use locationTypeDescriptor instead')
const LocationType$json = {
  '1': 'LocationType',
  '2': [
    {'1': 'LOCATION_TYPE_AVAILABILITYZONE', '2': 0},
    {'1': 'LOCATION_TYPE_LOCALZONE', '2': 1},
  ],
};

/// Descriptor for `LocationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List locationTypeDescriptor = $convert.base64Decode(
    'CgxMb2NhdGlvblR5cGUSIgoeTE9DQVRJT05fVFlQRV9BVkFJTEFCSUxJVFlaT05FEAASGwoXTE'
    '9DQVRJT05fVFlQRV9MT0NBTFpPTkUQAQ==');

@$core.Deprecated('Use mFADeleteDescriptor instead')
const MFADelete$json = {
  '1': 'MFADelete',
  '2': [
    {'1': 'M_F_A_DELETE_DISABLED', '2': 0},
    {'1': 'M_F_A_DELETE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `MFADelete`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mFADeleteDescriptor = $convert.base64Decode(
    'CglNRkFEZWxldGUSGQoVTV9GX0FfREVMRVRFX0RJU0FCTEVEEAASGAoUTV9GX0FfREVMRVRFX0'
    'VOQUJMRUQQAQ==');

@$core.Deprecated('Use mFADeleteStatusDescriptor instead')
const MFADeleteStatus$json = {
  '1': 'MFADeleteStatus',
  '2': [
    {'1': 'M_F_A_DELETE_STATUS_DISABLED', '2': 0},
    {'1': 'M_F_A_DELETE_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `MFADeleteStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mFADeleteStatusDescriptor = $convert.base64Decode(
    'Cg9NRkFEZWxldGVTdGF0dXMSIAocTV9GX0FfREVMRVRFX1NUQVRVU19ESVNBQkxFRBAAEh8KG0'
    '1fRl9BX0RFTEVURV9TVEFUVVNfRU5BQkxFRBAB');

@$core.Deprecated('Use metadataDirectiveDescriptor instead')
const MetadataDirective$json = {
  '1': 'MetadataDirective',
  '2': [
    {'1': 'METADATA_DIRECTIVE_REPLACE', '2': 0},
    {'1': 'METADATA_DIRECTIVE_COPY', '2': 1},
  ],
};

/// Descriptor for `MetadataDirective`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metadataDirectiveDescriptor = $convert.base64Decode(
    'ChFNZXRhZGF0YURpcmVjdGl2ZRIeChpNRVRBREFUQV9ESVJFQ1RJVkVfUkVQTEFDRRAAEhsKF0'
    '1FVEFEQVRBX0RJUkVDVElWRV9DT1BZEAE=');

@$core.Deprecated('Use metricsStatusDescriptor instead')
const MetricsStatus$json = {
  '1': 'MetricsStatus',
  '2': [
    {'1': 'METRICS_STATUS_DISABLED', '2': 0},
    {'1': 'METRICS_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `MetricsStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricsStatusDescriptor = $convert.base64Decode(
    'Cg1NZXRyaWNzU3RhdHVzEhsKF01FVFJJQ1NfU1RBVFVTX0RJU0FCTEVEEAASGgoWTUVUUklDU1'
    '9TVEFUVVNfRU5BQkxFRBAB');

@$core.Deprecated('Use objectAttributesDescriptor instead')
const ObjectAttributes$json = {
  '1': 'ObjectAttributes',
  '2': [
    {'1': 'OBJECT_ATTRIBUTES_OBJECT_PARTS', '2': 0},
    {'1': 'OBJECT_ATTRIBUTES_CHECKSUM', '2': 1},
    {'1': 'OBJECT_ATTRIBUTES_STORAGE_CLASS', '2': 2},
    {'1': 'OBJECT_ATTRIBUTES_OBJECT_SIZE', '2': 3},
    {'1': 'OBJECT_ATTRIBUTES_ETAG', '2': 4},
  ],
};

/// Descriptor for `ObjectAttributes`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectAttributesDescriptor = $convert.base64Decode(
    'ChBPYmplY3RBdHRyaWJ1dGVzEiIKHk9CSkVDVF9BVFRSSUJVVEVTX09CSkVDVF9QQVJUUxAAEh'
    '4KGk9CSkVDVF9BVFRSSUJVVEVTX0NIRUNLU1VNEAESIwofT0JKRUNUX0FUVFJJQlVURVNfU1RP'
    'UkFHRV9DTEFTUxACEiEKHU9CSkVDVF9BVFRSSUJVVEVTX09CSkVDVF9TSVpFEAMSGgoWT0JKRU'
    'NUX0FUVFJJQlVURVNfRVRBRxAE');

@$core.Deprecated('Use objectCannedACLDescriptor instead')
const ObjectCannedACL$json = {
  '1': 'ObjectCannedACL',
  '2': [
    {'1': 'OBJECT_CANNED_A_C_L_AUTHENTICATED_READ', '2': 0},
    {'1': 'OBJECT_CANNED_A_C_L_BUCKET_OWNER_READ', '2': 1},
    {'1': 'OBJECT_CANNED_A_C_L_PRIVATE', '2': 2},
    {'1': 'OBJECT_CANNED_A_C_L_AWS_EXEC_READ', '2': 3},
    {'1': 'OBJECT_CANNED_A_C_L_BUCKET_OWNER_FULL_CONTROL', '2': 4},
    {'1': 'OBJECT_CANNED_A_C_L_PUBLIC_READ', '2': 5},
    {'1': 'OBJECT_CANNED_A_C_L_PUBLIC_READ_WRITE', '2': 6},
  ],
};

/// Descriptor for `ObjectCannedACL`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectCannedACLDescriptor = $convert.base64Decode(
    'Cg9PYmplY3RDYW5uZWRBQ0wSKgomT0JKRUNUX0NBTk5FRF9BX0NfTF9BVVRIRU5USUNBVEVEX1'
    'JFQUQQABIpCiVPQkpFQ1RfQ0FOTkVEX0FfQ19MX0JVQ0tFVF9PV05FUl9SRUFEEAESHwobT0JK'
    'RUNUX0NBTk5FRF9BX0NfTF9QUklWQVRFEAISJQohT0JKRUNUX0NBTk5FRF9BX0NfTF9BV1NfRV'
    'hFQ19SRUFEEAMSMQotT0JKRUNUX0NBTk5FRF9BX0NfTF9CVUNLRVRfT1dORVJfRlVMTF9DT05U'
    'Uk9MEAQSIwofT0JKRUNUX0NBTk5FRF9BX0NfTF9QVUJMSUNfUkVBRBAFEikKJU9CSkVDVF9DQU'
    '5ORURfQV9DX0xfUFVCTElDX1JFQURfV1JJVEUQBg==');

@$core.Deprecated('Use objectLockEnabledDescriptor instead')
const ObjectLockEnabled$json = {
  '1': 'ObjectLockEnabled',
  '2': [
    {'1': 'OBJECT_LOCK_ENABLED_ENABLED', '2': 0},
  ],
};

/// Descriptor for `ObjectLockEnabled`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectLockEnabledDescriptor = $convert.base64Decode(
    'ChFPYmplY3RMb2NrRW5hYmxlZBIfChtPQkpFQ1RfTE9DS19FTkFCTEVEX0VOQUJMRUQQAA==');

@$core.Deprecated('Use objectLockLegalHoldStatusDescriptor instead')
const ObjectLockLegalHoldStatus$json = {
  '1': 'ObjectLockLegalHoldStatus',
  '2': [
    {'1': 'OBJECT_LOCK_LEGAL_HOLD_STATUS_OFF', '2': 0},
    {'1': 'OBJECT_LOCK_LEGAL_HOLD_STATUS_ON', '2': 1},
  ],
};

/// Descriptor for `ObjectLockLegalHoldStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectLockLegalHoldStatusDescriptor =
    $convert.base64Decode(
        'ChlPYmplY3RMb2NrTGVnYWxIb2xkU3RhdHVzEiUKIU9CSkVDVF9MT0NLX0xFR0FMX0hPTERfU1'
        'RBVFVTX09GRhAAEiQKIE9CSkVDVF9MT0NLX0xFR0FMX0hPTERfU1RBVFVTX09OEAE=');

@$core.Deprecated('Use objectLockModeDescriptor instead')
const ObjectLockMode$json = {
  '1': 'ObjectLockMode',
  '2': [
    {'1': 'OBJECT_LOCK_MODE_GOVERNANCE', '2': 0},
    {'1': 'OBJECT_LOCK_MODE_COMPLIANCE', '2': 1},
  ],
};

/// Descriptor for `ObjectLockMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectLockModeDescriptor = $convert.base64Decode(
    'Cg5PYmplY3RMb2NrTW9kZRIfChtPQkpFQ1RfTE9DS19NT0RFX0dPVkVSTkFOQ0UQABIfChtPQk'
    'pFQ1RfTE9DS19NT0RFX0NPTVBMSUFOQ0UQAQ==');

@$core.Deprecated('Use objectLockRetentionModeDescriptor instead')
const ObjectLockRetentionMode$json = {
  '1': 'ObjectLockRetentionMode',
  '2': [
    {'1': 'OBJECT_LOCK_RETENTION_MODE_GOVERNANCE', '2': 0},
    {'1': 'OBJECT_LOCK_RETENTION_MODE_COMPLIANCE', '2': 1},
  ],
};

/// Descriptor for `ObjectLockRetentionMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectLockRetentionModeDescriptor = $convert.base64Decode(
    'ChdPYmplY3RMb2NrUmV0ZW50aW9uTW9kZRIpCiVPQkpFQ1RfTE9DS19SRVRFTlRJT05fTU9ERV'
    '9HT1ZFUk5BTkNFEAASKQolT0JKRUNUX0xPQ0tfUkVURU5USU9OX01PREVfQ09NUExJQU5DRRAB');

@$core.Deprecated('Use objectOwnershipDescriptor instead')
const ObjectOwnership$json = {
  '1': 'ObjectOwnership',
  '2': [
    {'1': 'OBJECT_OWNERSHIP_BUCKETOWNERPREFERRED', '2': 0},
    {'1': 'OBJECT_OWNERSHIP_OBJECTWRITER', '2': 1},
    {'1': 'OBJECT_OWNERSHIP_BUCKETOWNERENFORCED', '2': 2},
  ],
};

/// Descriptor for `ObjectOwnership`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectOwnershipDescriptor = $convert.base64Decode(
    'Cg9PYmplY3RPd25lcnNoaXASKQolT0JKRUNUX09XTkVSU0hJUF9CVUNLRVRPV05FUlBSRUZFUl'
    'JFRBAAEiEKHU9CSkVDVF9PV05FUlNISVBfT0JKRUNUV1JJVEVSEAESKAokT0JKRUNUX09XTkVS'
    'U0hJUF9CVUNLRVRPV05FUkVORk9SQ0VEEAI=');

@$core.Deprecated('Use objectStorageClassDescriptor instead')
const ObjectStorageClass$json = {
  '1': 'ObjectStorageClass',
  '2': [
    {'1': 'OBJECT_STORAGE_CLASS_OUTPOSTS', '2': 0},
    {'1': 'OBJECT_STORAGE_CLASS_EXPRESS_ONEZONE', '2': 1},
    {'1': 'OBJECT_STORAGE_CLASS_REDUCED_REDUNDANCY', '2': 2},
    {'1': 'OBJECT_STORAGE_CLASS_DEEP_ARCHIVE', '2': 3},
    {'1': 'OBJECT_STORAGE_CLASS_SNOW', '2': 4},
    {'1': 'OBJECT_STORAGE_CLASS_ONEZONE_IA', '2': 5},
    {'1': 'OBJECT_STORAGE_CLASS_STANDARD_IA', '2': 6},
    {'1': 'OBJECT_STORAGE_CLASS_GLACIER', '2': 7},
    {'1': 'OBJECT_STORAGE_CLASS_STANDARD', '2': 8},
    {'1': 'OBJECT_STORAGE_CLASS_FSX_ONTAP', '2': 9},
    {'1': 'OBJECT_STORAGE_CLASS_GLACIER_IR', '2': 10},
    {'1': 'OBJECT_STORAGE_CLASS_FSX_OPENZFS', '2': 11},
    {'1': 'OBJECT_STORAGE_CLASS_INTELLIGENT_TIERING', '2': 12},
  ],
};

/// Descriptor for `ObjectStorageClass`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectStorageClassDescriptor = $convert.base64Decode(
    'ChJPYmplY3RTdG9yYWdlQ2xhc3MSIQodT0JKRUNUX1NUT1JBR0VfQ0xBU1NfT1VUUE9TVFMQAB'
    'IoCiRPQkpFQ1RfU1RPUkFHRV9DTEFTU19FWFBSRVNTX09ORVpPTkUQARIrCidPQkpFQ1RfU1RP'
    'UkFHRV9DTEFTU19SRURVQ0VEX1JFRFVOREFOQ1kQAhIlCiFPQkpFQ1RfU1RPUkFHRV9DTEFTU1'
    '9ERUVQX0FSQ0hJVkUQAxIdChlPQkpFQ1RfU1RPUkFHRV9DTEFTU19TTk9XEAQSIwofT0JKRUNU'
    'X1NUT1JBR0VfQ0xBU1NfT05FWk9ORV9JQRAFEiQKIE9CSkVDVF9TVE9SQUdFX0NMQVNTX1NUQU'
    '5EQVJEX0lBEAYSIAocT0JKRUNUX1NUT1JBR0VfQ0xBU1NfR0xBQ0lFUhAHEiEKHU9CSkVDVF9T'
    'VE9SQUdFX0NMQVNTX1NUQU5EQVJEEAgSIgoeT0JKRUNUX1NUT1JBR0VfQ0xBU1NfRlNYX09OVE'
    'FQEAkSIwofT0JKRUNUX1NUT1JBR0VfQ0xBU1NfR0xBQ0lFUl9JUhAKEiQKIE9CSkVDVF9TVE9S'
    'QUdFX0NMQVNTX0ZTWF9PUEVOWkZTEAsSLAooT0JKRUNUX1NUT1JBR0VfQ0xBU1NfSU5URUxMSU'
    'dFTlRfVElFUklORxAM');

@$core.Deprecated('Use objectVersionStorageClassDescriptor instead')
const ObjectVersionStorageClass$json = {
  '1': 'ObjectVersionStorageClass',
  '2': [
    {'1': 'OBJECT_VERSION_STORAGE_CLASS_STANDARD', '2': 0},
  ],
};

/// Descriptor for `ObjectVersionStorageClass`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List objectVersionStorageClassDescriptor =
    $convert.base64Decode(
        'ChlPYmplY3RWZXJzaW9uU3RvcmFnZUNsYXNzEikKJU9CSkVDVF9WRVJTSU9OX1NUT1JBR0VfQ0'
        'xBU1NfU1RBTkRBUkQQAA==');

@$core.Deprecated('Use optionalObjectAttributesDescriptor instead')
const OptionalObjectAttributes$json = {
  '1': 'OptionalObjectAttributes',
  '2': [
    {'1': 'OPTIONAL_OBJECT_ATTRIBUTES_RESTORE_STATUS', '2': 0},
  ],
};

/// Descriptor for `OptionalObjectAttributes`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List optionalObjectAttributesDescriptor =
    $convert.base64Decode(
        'ChhPcHRpb25hbE9iamVjdEF0dHJpYnV0ZXMSLQopT1BUSU9OQUxfT0JKRUNUX0FUVFJJQlVURV'
        'NfUkVTVE9SRV9TVEFUVVMQAA==');

@$core.Deprecated('Use ownerOverrideDescriptor instead')
const OwnerOverride$json = {
  '1': 'OwnerOverride',
  '2': [
    {'1': 'OWNER_OVERRIDE_DESTINATION', '2': 0},
  ],
};

/// Descriptor for `OwnerOverride`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List ownerOverrideDescriptor = $convert.base64Decode(
    'Cg1Pd25lck92ZXJyaWRlEh4KGk9XTkVSX09WRVJSSURFX0RFU1RJTkFUSU9OEAA=');

@$core.Deprecated('Use partitionDateSourceDescriptor instead')
const PartitionDateSource$json = {
  '1': 'PartitionDateSource',
  '2': [
    {'1': 'PARTITION_DATE_SOURCE_DELIVERYTIME', '2': 0},
    {'1': 'PARTITION_DATE_SOURCE_EVENTTIME', '2': 1},
  ],
};

/// Descriptor for `PartitionDateSource`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List partitionDateSourceDescriptor = $convert.base64Decode(
    'ChNQYXJ0aXRpb25EYXRlU291cmNlEiYKIlBBUlRJVElPTl9EQVRFX1NPVVJDRV9ERUxJVkVSWV'
    'RJTUUQABIjCh9QQVJUSVRJT05fREFURV9TT1VSQ0VfRVZFTlRUSU1FEAE=');

@$core.Deprecated('Use payerDescriptor instead')
const Payer$json = {
  '1': 'Payer',
  '2': [
    {'1': 'PAYER_REQUESTER', '2': 0},
    {'1': 'PAYER_BUCKETOWNER', '2': 1},
  ],
};

/// Descriptor for `Payer`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List payerDescriptor = $convert.base64Decode(
    'CgVQYXllchITCg9QQVlFUl9SRVFVRVNURVIQABIVChFQQVlFUl9CVUNLRVRPV05FUhAB');

@$core.Deprecated('Use permissionDescriptor instead')
const Permission$json = {
  '1': 'Permission',
  '2': [
    {'1': 'PERMISSION_READ', '2': 0},
    {'1': 'PERMISSION_WRITE_ACP', '2': 1},
    {'1': 'PERMISSION_WRITE', '2': 2},
    {'1': 'PERMISSION_FULL_CONTROL', '2': 3},
    {'1': 'PERMISSION_READ_ACP', '2': 4},
  ],
};

/// Descriptor for `Permission`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List permissionDescriptor = $convert.base64Decode(
    'CgpQZXJtaXNzaW9uEhMKD1BFUk1JU1NJT05fUkVBRBAAEhgKFFBFUk1JU1NJT05fV1JJVEVfQU'
    'NQEAESFAoQUEVSTUlTU0lPTl9XUklURRACEhsKF1BFUk1JU1NJT05fRlVMTF9DT05UUk9MEAMS'
    'FwoTUEVSTUlTU0lPTl9SRUFEX0FDUBAE');

@$core.Deprecated('Use protocolDescriptor instead')
const Protocol$json = {
  '1': 'Protocol',
  '2': [
    {'1': 'PROTOCOL_HTTPS', '2': 0},
    {'1': 'PROTOCOL_HTTP', '2': 1},
  ],
};

/// Descriptor for `Protocol`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List protocolDescriptor = $convert.base64Decode(
    'CghQcm90b2NvbBISCg5QUk9UT0NPTF9IVFRQUxAAEhEKDVBST1RPQ09MX0hUVFAQAQ==');

@$core.Deprecated('Use quoteFieldsDescriptor instead')
const QuoteFields$json = {
  '1': 'QuoteFields',
  '2': [
    {'1': 'QUOTE_FIELDS_ALWAYS', '2': 0},
    {'1': 'QUOTE_FIELDS_ASNEEDED', '2': 1},
  ],
};

/// Descriptor for `QuoteFields`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List quoteFieldsDescriptor = $convert.base64Decode(
    'CgtRdW90ZUZpZWxkcxIXChNRVU9URV9GSUVMRFNfQUxXQVlTEAASGQoVUVVPVEVfRklFTERTX0'
    'FTTkVFREVEEAE=');

@$core.Deprecated('Use replicaModificationsStatusDescriptor instead')
const ReplicaModificationsStatus$json = {
  '1': 'ReplicaModificationsStatus',
  '2': [
    {'1': 'REPLICA_MODIFICATIONS_STATUS_DISABLED', '2': 0},
    {'1': 'REPLICA_MODIFICATIONS_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ReplicaModificationsStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List replicaModificationsStatusDescriptor =
    $convert.base64Decode(
        'ChpSZXBsaWNhTW9kaWZpY2F0aW9uc1N0YXR1cxIpCiVSRVBMSUNBX01PRElGSUNBVElPTlNfU1'
        'RBVFVTX0RJU0FCTEVEEAASKAokUkVQTElDQV9NT0RJRklDQVRJT05TX1NUQVRVU19FTkFCTEVE'
        'EAE=');

@$core.Deprecated('Use replicationRuleStatusDescriptor instead')
const ReplicationRuleStatus$json = {
  '1': 'ReplicationRuleStatus',
  '2': [
    {'1': 'REPLICATION_RULE_STATUS_DISABLED', '2': 0},
    {'1': 'REPLICATION_RULE_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ReplicationRuleStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List replicationRuleStatusDescriptor = $convert.base64Decode(
    'ChVSZXBsaWNhdGlvblJ1bGVTdGF0dXMSJAogUkVQTElDQVRJT05fUlVMRV9TVEFUVVNfRElTQU'
    'JMRUQQABIjCh9SRVBMSUNBVElPTl9SVUxFX1NUQVRVU19FTkFCTEVEEAE=');

@$core.Deprecated('Use replicationStatusDescriptor instead')
const ReplicationStatus$json = {
  '1': 'ReplicationStatus',
  '2': [
    {'1': 'REPLICATION_STATUS_PENDING', '2': 0},
    {'1': 'REPLICATION_STATUS_REPLICA', '2': 1},
    {'1': 'REPLICATION_STATUS_COMPLETE', '2': 2},
    {'1': 'REPLICATION_STATUS_COMPLETED', '2': 3},
    {'1': 'REPLICATION_STATUS_FAILED', '2': 4},
  ],
};

/// Descriptor for `ReplicationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List replicationStatusDescriptor = $convert.base64Decode(
    'ChFSZXBsaWNhdGlvblN0YXR1cxIeChpSRVBMSUNBVElPTl9TVEFUVVNfUEVORElORxAAEh4KGl'
    'JFUExJQ0FUSU9OX1NUQVRVU19SRVBMSUNBEAESHwobUkVQTElDQVRJT05fU1RBVFVTX0NPTVBM'
    'RVRFEAISIAocUkVQTElDQVRJT05fU1RBVFVTX0NPTVBMRVRFRBADEh0KGVJFUExJQ0FUSU9OX1'
    'NUQVRVU19GQUlMRUQQBA==');

@$core.Deprecated('Use replicationTimeStatusDescriptor instead')
const ReplicationTimeStatus$json = {
  '1': 'ReplicationTimeStatus',
  '2': [
    {'1': 'REPLICATION_TIME_STATUS_DISABLED', '2': 0},
    {'1': 'REPLICATION_TIME_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ReplicationTimeStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List replicationTimeStatusDescriptor = $convert.base64Decode(
    'ChVSZXBsaWNhdGlvblRpbWVTdGF0dXMSJAogUkVQTElDQVRJT05fVElNRV9TVEFUVVNfRElTQU'
    'JMRUQQABIjCh9SRVBMSUNBVElPTl9USU1FX1NUQVRVU19FTkFCTEVEEAE=');

@$core.Deprecated('Use requestChargedDescriptor instead')
const RequestCharged$json = {
  '1': 'RequestCharged',
  '2': [
    {'1': 'REQUEST_CHARGED_REQUESTER', '2': 0},
  ],
};

/// Descriptor for `RequestCharged`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List requestChargedDescriptor = $convert.base64Decode(
    'Cg5SZXF1ZXN0Q2hhcmdlZBIdChlSRVFVRVNUX0NIQVJHRURfUkVRVUVTVEVSEAA=');

@$core.Deprecated('Use requestPayerDescriptor instead')
const RequestPayer$json = {
  '1': 'RequestPayer',
  '2': [
    {'1': 'REQUEST_PAYER_REQUESTER', '2': 0},
  ],
};

/// Descriptor for `RequestPayer`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List requestPayerDescriptor = $convert.base64Decode(
    'CgxSZXF1ZXN0UGF5ZXISGwoXUkVRVUVTVF9QQVlFUl9SRVFVRVNURVIQAA==');

@$core.Deprecated('Use restoreRequestTypeDescriptor instead')
const RestoreRequestType$json = {
  '1': 'RestoreRequestType',
  '2': [
    {'1': 'RESTORE_REQUEST_TYPE_SELECT', '2': 0},
  ],
};

/// Descriptor for `RestoreRequestType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List restoreRequestTypeDescriptor = $convert.base64Decode(
    'ChJSZXN0b3JlUmVxdWVzdFR5cGUSHwobUkVTVE9SRV9SRVFVRVNUX1RZUEVfU0VMRUNUEAA=');

@$core.Deprecated('Use s3TablesBucketTypeDescriptor instead')
const S3TablesBucketType$json = {
  '1': 'S3TablesBucketType',
  '2': [
    {'1': 'S3_TABLES_BUCKET_TYPE_CUSTOMER', '2': 0},
    {'1': 'S3_TABLES_BUCKET_TYPE_AWS', '2': 1},
  ],
};

/// Descriptor for `S3TablesBucketType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List s3TablesBucketTypeDescriptor = $convert.base64Decode(
    'ChJTM1RhYmxlc0J1Y2tldFR5cGUSIgoeUzNfVEFCTEVTX0JVQ0tFVF9UWVBFX0NVU1RPTUVSEA'
    'ASHQoZUzNfVEFCTEVTX0JVQ0tFVF9UWVBFX0FXUxAB');

@$core.Deprecated('Use serverSideEncryptionDescriptor instead')
const ServerSideEncryption$json = {
  '1': 'ServerSideEncryption',
  '2': [
    {'1': 'SERVER_SIDE_ENCRYPTION_AWS_KMS', '2': 0},
    {'1': 'SERVER_SIDE_ENCRYPTION_AWS_KMS_DSSE', '2': 1},
    {'1': 'SERVER_SIDE_ENCRYPTION_AES256', '2': 2},
    {'1': 'SERVER_SIDE_ENCRYPTION_AWS_FSX', '2': 3},
  ],
};

/// Descriptor for `ServerSideEncryption`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List serverSideEncryptionDescriptor = $convert.base64Decode(
    'ChRTZXJ2ZXJTaWRlRW5jcnlwdGlvbhIiCh5TRVJWRVJfU0lERV9FTkNSWVBUSU9OX0FXU19LTV'
    'MQABInCiNTRVJWRVJfU0lERV9FTkNSWVBUSU9OX0FXU19LTVNfRFNTRRABEiEKHVNFUlZFUl9T'
    'SURFX0VOQ1JZUFRJT05fQUVTMjU2EAISIgoeU0VSVkVSX1NJREVfRU5DUllQVElPTl9BV1NfRl'
    'NYEAM=');

@$core.Deprecated('Use sessionModeDescriptor instead')
const SessionMode$json = {
  '1': 'SessionMode',
  '2': [
    {'1': 'SESSION_MODE_READONLY', '2': 0},
    {'1': 'SESSION_MODE_READWRITE', '2': 1},
  ],
};

/// Descriptor for `SessionMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sessionModeDescriptor = $convert.base64Decode(
    'CgtTZXNzaW9uTW9kZRIZChVTRVNTSU9OX01PREVfUkVBRE9OTFkQABIaChZTRVNTSU9OX01PRE'
    'VfUkVBRFdSSVRFEAE=');

@$core.Deprecated('Use sseKmsEncryptedObjectsStatusDescriptor instead')
const SseKmsEncryptedObjectsStatus$json = {
  '1': 'SseKmsEncryptedObjectsStatus',
  '2': [
    {'1': 'SSE_KMS_ENCRYPTED_OBJECTS_STATUS_DISABLED', '2': 0},
    {'1': 'SSE_KMS_ENCRYPTED_OBJECTS_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `SseKmsEncryptedObjectsStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sseKmsEncryptedObjectsStatusDescriptor =
    $convert.base64Decode(
        'ChxTc2VLbXNFbmNyeXB0ZWRPYmplY3RzU3RhdHVzEi0KKVNTRV9LTVNfRU5DUllQVEVEX09CSk'
        'VDVFNfU1RBVFVTX0RJU0FCTEVEEAASLAooU1NFX0tNU19FTkNSWVBURURfT0JKRUNUU19TVEFU'
        'VVNfRU5BQkxFRBAB');

@$core.Deprecated('Use storageClassDescriptor instead')
const StorageClass$json = {
  '1': 'StorageClass',
  '2': [
    {'1': 'STORAGE_CLASS_OUTPOSTS', '2': 0},
    {'1': 'STORAGE_CLASS_EXPRESS_ONEZONE', '2': 1},
    {'1': 'STORAGE_CLASS_REDUCED_REDUNDANCY', '2': 2},
    {'1': 'STORAGE_CLASS_DEEP_ARCHIVE', '2': 3},
    {'1': 'STORAGE_CLASS_SNOW', '2': 4},
    {'1': 'STORAGE_CLASS_ONEZONE_IA', '2': 5},
    {'1': 'STORAGE_CLASS_STANDARD_IA', '2': 6},
    {'1': 'STORAGE_CLASS_GLACIER', '2': 7},
    {'1': 'STORAGE_CLASS_STANDARD', '2': 8},
    {'1': 'STORAGE_CLASS_FSX_ONTAP', '2': 9},
    {'1': 'STORAGE_CLASS_GLACIER_IR', '2': 10},
    {'1': 'STORAGE_CLASS_FSX_OPENZFS', '2': 11},
    {'1': 'STORAGE_CLASS_INTELLIGENT_TIERING', '2': 12},
  ],
};

/// Descriptor for `StorageClass`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List storageClassDescriptor = $convert.base64Decode(
    'CgxTdG9yYWdlQ2xhc3MSGgoWU1RPUkFHRV9DTEFTU19PVVRQT1NUUxAAEiEKHVNUT1JBR0VfQ0'
    'xBU1NfRVhQUkVTU19PTkVaT05FEAESJAogU1RPUkFHRV9DTEFTU19SRURVQ0VEX1JFRFVOREFO'
    'Q1kQAhIeChpTVE9SQUdFX0NMQVNTX0RFRVBfQVJDSElWRRADEhYKElNUT1JBR0VfQ0xBU1NfU0'
    '5PVxAEEhwKGFNUT1JBR0VfQ0xBU1NfT05FWk9ORV9JQRAFEh0KGVNUT1JBR0VfQ0xBU1NfU1RB'
    'TkRBUkRfSUEQBhIZChVTVE9SQUdFX0NMQVNTX0dMQUNJRVIQBxIaChZTVE9SQUdFX0NMQVNTX1'
    'NUQU5EQVJEEAgSGwoXU1RPUkFHRV9DTEFTU19GU1hfT05UQVAQCRIcChhTVE9SQUdFX0NMQVNT'
    'X0dMQUNJRVJfSVIQChIdChlTVE9SQUdFX0NMQVNTX0ZTWF9PUEVOWkZTEAsSJQohU1RPUkFHRV'
    '9DTEFTU19JTlRFTExJR0VOVF9USUVSSU5HEAw=');

@$core.Deprecated('Use storageClassAnalysisSchemaVersionDescriptor instead')
const StorageClassAnalysisSchemaVersion$json = {
  '1': 'StorageClassAnalysisSchemaVersion',
  '2': [
    {'1': 'STORAGE_CLASS_ANALYSIS_SCHEMA_VERSION_V_1', '2': 0},
  ],
};

/// Descriptor for `StorageClassAnalysisSchemaVersion`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List storageClassAnalysisSchemaVersionDescriptor =
    $convert.base64Decode(
        'CiFTdG9yYWdlQ2xhc3NBbmFseXNpc1NjaGVtYVZlcnNpb24SLQopU1RPUkFHRV9DTEFTU19BTk'
        'FMWVNJU19TQ0hFTUFfVkVSU0lPTl9WXzEQAA==');

@$core.Deprecated('Use tableSseAlgorithmDescriptor instead')
const TableSseAlgorithm$json = {
  '1': 'TableSseAlgorithm',
  '2': [
    {'1': 'TABLE_SSE_ALGORITHM_AWS_KMS', '2': 0},
    {'1': 'TABLE_SSE_ALGORITHM_AES256', '2': 1},
  ],
};

/// Descriptor for `TableSseAlgorithm`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tableSseAlgorithmDescriptor = $convert.base64Decode(
    'ChFUYWJsZVNzZUFsZ29yaXRobRIfChtUQUJMRV9TU0VfQUxHT1JJVEhNX0FXU19LTVMQABIeCh'
    'pUQUJMRV9TU0VfQUxHT1JJVEhNX0FFUzI1NhAB');

@$core.Deprecated('Use taggingDirectiveDescriptor instead')
const TaggingDirective$json = {
  '1': 'TaggingDirective',
  '2': [
    {'1': 'TAGGING_DIRECTIVE_REPLACE', '2': 0},
    {'1': 'TAGGING_DIRECTIVE_COPY', '2': 1},
  ],
};

/// Descriptor for `TaggingDirective`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List taggingDirectiveDescriptor = $convert.base64Decode(
    'ChBUYWdnaW5nRGlyZWN0aXZlEh0KGVRBR0dJTkdfRElSRUNUSVZFX1JFUExBQ0UQABIaChZUQU'
    'dHSU5HX0RJUkVDVElWRV9DT1BZEAE=');

@$core.Deprecated('Use tierDescriptor instead')
const Tier$json = {
  '1': 'Tier',
  '2': [
    {'1': 'TIER_BULK', '2': 0},
    {'1': 'TIER_EXPEDITED', '2': 1},
    {'1': 'TIER_STANDARD', '2': 2},
  ],
};

/// Descriptor for `Tier`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tierDescriptor = $convert.base64Decode(
    'CgRUaWVyEg0KCVRJRVJfQlVMSxAAEhIKDlRJRVJfRVhQRURJVEVEEAESEQoNVElFUl9TVEFORE'
    'FSRBAC');

@$core.Deprecated('Use transitionDefaultMinimumObjectSizeDescriptor instead')
const TransitionDefaultMinimumObjectSize$json = {
  '1': 'TransitionDefaultMinimumObjectSize',
  '2': [
    {
      '1': 'TRANSITION_DEFAULT_MINIMUM_OBJECT_SIZE_ALL_STORAGE_CLASSES_128K',
      '2': 0
    },
    {
      '1': 'TRANSITION_DEFAULT_MINIMUM_OBJECT_SIZE_VARIES_BY_STORAGE_CLASS',
      '2': 1
    },
  ],
};

/// Descriptor for `TransitionDefaultMinimumObjectSize`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List transitionDefaultMinimumObjectSizeDescriptor =
    $convert.base64Decode(
        'CiJUcmFuc2l0aW9uRGVmYXVsdE1pbmltdW1PYmplY3RTaXplEkMKP1RSQU5TSVRJT05fREVGQV'
        'VMVF9NSU5JTVVNX09CSkVDVF9TSVpFX0FMTF9TVE9SQUdFX0NMQVNTRVNfMTI4SxAAEkIKPlRS'
        'QU5TSVRJT05fREVGQVVMVF9NSU5JTVVNX09CSkVDVF9TSVpFX1ZBUklFU19CWV9TVE9SQUdFX0'
        'NMQVNTEAE=');

@$core.Deprecated('Use transitionStorageClassDescriptor instead')
const TransitionStorageClass$json = {
  '1': 'TransitionStorageClass',
  '2': [
    {'1': 'TRANSITION_STORAGE_CLASS_DEEP_ARCHIVE', '2': 0},
    {'1': 'TRANSITION_STORAGE_CLASS_ONEZONE_IA', '2': 1},
    {'1': 'TRANSITION_STORAGE_CLASS_STANDARD_IA', '2': 2},
    {'1': 'TRANSITION_STORAGE_CLASS_GLACIER', '2': 3},
    {'1': 'TRANSITION_STORAGE_CLASS_GLACIER_IR', '2': 4},
    {'1': 'TRANSITION_STORAGE_CLASS_INTELLIGENT_TIERING', '2': 5},
  ],
};

/// Descriptor for `TransitionStorageClass`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List transitionStorageClassDescriptor = $convert.base64Decode(
    'ChZUcmFuc2l0aW9uU3RvcmFnZUNsYXNzEikKJVRSQU5TSVRJT05fU1RPUkFHRV9DTEFTU19ERU'
    'VQX0FSQ0hJVkUQABInCiNUUkFOU0lUSU9OX1NUT1JBR0VfQ0xBU1NfT05FWk9ORV9JQRABEigK'
    'JFRSQU5TSVRJT05fU1RPUkFHRV9DTEFTU19TVEFOREFSRF9JQRACEiQKIFRSQU5TSVRJT05fU1'
    'RPUkFHRV9DTEFTU19HTEFDSUVSEAMSJwojVFJBTlNJVElPTl9TVE9SQUdFX0NMQVNTX0dMQUNJ'
    'RVJfSVIQBBIwCixUUkFOU0lUSU9OX1NUT1JBR0VfQ0xBU1NfSU5URUxMSUdFTlRfVElFUklORx'
    'AF');

@$core.Deprecated('Use typeDescriptor instead')
const Type$json = {
  '1': 'Type',
  '2': [
    {'1': 'TYPE_GROUP', '2': 0},
    {'1': 'TYPE_AMAZONCUSTOMERBYEMAIL', '2': 1},
    {'1': 'TYPE_CANONICALUSER', '2': 2},
  ],
};

/// Descriptor for `Type`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List typeDescriptor = $convert.base64Decode(
    'CgRUeXBlEg4KClRZUEVfR1JPVVAQABIeChpUWVBFX0FNQVpPTkNVU1RPTUVSQllFTUFJTBABEh'
    'YKElRZUEVfQ0FOT05JQ0FMVVNFUhAC');

@$core.Deprecated('Use abacStatusDescriptor instead')
const AbacStatus$json = {
  '1': 'AbacStatus',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketAbacStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `AbacStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List abacStatusDescriptor = $convert.base64Decode(
    'CgpBYmFjU3RhdHVzEi8KBnN0YXR1cxiQ5PsCIAEoDjIULnMzLkJ1Y2tldEFiYWNTdGF0dXNSBn'
    'N0YXR1cw==');

@$core.Deprecated('Use abortIncompleteMultipartUploadDescriptor instead')
const AbortIncompleteMultipartUpload$json = {
  '1': 'AbortIncompleteMultipartUpload',
  '2': [
    {
      '1': 'daysafterinitiation',
      '3': 522354315,
      '4': 1,
      '5': 5,
      '10': 'daysafterinitiation'
    },
  ],
};

/// Descriptor for `AbortIncompleteMultipartUpload`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List abortIncompleteMultipartUploadDescriptor =
    $convert.base64Decode(
        'Ch5BYm9ydEluY29tcGxldGVNdWx0aXBhcnRVcGxvYWQSNAoTZGF5c2FmdGVyaW5pdGlhdGlvbh'
        'iL/Yn5ASABKAVSE2RheXNhZnRlcmluaXRpYXRpb24=');

@$core.Deprecated('Use abortMultipartUploadOutputDescriptor instead')
const AbortMultipartUploadOutput$json = {
  '1': 'AbortMultipartUploadOutput',
  '2': [
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `AbortMultipartUploadOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List abortMultipartUploadOutputDescriptor =
    $convert.base64Decode(
        'ChpBYm9ydE11bHRpcGFydFVwbG9hZE91dHB1dBI+Cg5yZXF1ZXN0Y2hhcmdlZBiT0Ku5ASABKA'
        '4yEi5zMy5SZXF1ZXN0Q2hhcmdlZFIOcmVxdWVzdGNoYXJnZWQ=');

@$core.Deprecated('Use abortMultipartUploadRequestDescriptor instead')
const AbortMultipartUploadRequest$json = {
  '1': 'AbortMultipartUploadRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'ifmatchinitiatedtime',
      '3': 444607790,
      '4': 1,
      '5': 9,
      '10': 'ifmatchinitiatedtime'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'uploadid', '3': 449040722, '4': 1, '5': 9, '10': 'uploadid'},
  ],
};

/// Descriptor for `AbortMultipartUploadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List abortMultipartUploadRequestDescriptor = $convert.base64Decode(
    'ChtBYm9ydE11bHRpcGFydFVwbG9hZFJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZX'
    'QSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchI2'
    'ChRpZm1hdGNoaW5pdGlhdGVkdGltZRiu2oDUASABKAlSFGlmbWF0Y2hpbml0aWF0ZWR0aW1lEh'
    'MKA2tleRiNkutoIAEoCVIDa2V5EjgKDHJlcXVlc3RwYXllchik5uH1ASABKA4yEC5zMy5SZXF1'
    'ZXN0UGF5ZXJSDHJlcXVlc3RwYXllchIeCgh1cGxvYWRpZBjSoo/WASABKAlSCHVwbG9hZGlk');

@$core.Deprecated('Use accelerateConfigurationDescriptor instead')
const AccelerateConfiguration$json = {
  '1': 'AccelerateConfiguration',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketAccelerateStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `AccelerateConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accelerateConfigurationDescriptor =
    $convert.base64Decode(
        'ChdBY2NlbGVyYXRlQ29uZmlndXJhdGlvbhI1CgZzdGF0dXMYkOT7AiABKA4yGi5zMy5CdWNrZX'
        'RBY2NlbGVyYXRlU3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use accessControlPolicyDescriptor instead')
const AccessControlPolicy$json = {
  '1': 'AccessControlPolicy',
  '2': [
    {
      '1': 'grants',
      '3': 226910555,
      '4': 3,
      '5': 11,
      '6': '.s3.Grant',
      '10': 'grants'
    },
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
  ],
};

/// Descriptor for `AccessControlPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessControlPolicyDescriptor = $convert.base64Decode(
    'ChNBY2Nlc3NDb250cm9sUG9saWN5EiQKBmdyYW50cxjbwplsIAMoCzIJLnMzLkdyYW50UgZncm'
    'FudHMSIwoFb3duZXIY9fyK2QEgASgLMgkuczMuT3duZXJSBW93bmVy');

@$core.Deprecated('Use accessControlTranslationDescriptor instead')
const AccessControlTranslation$json = {
  '1': 'AccessControlTranslation',
  '2': [
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 14,
      '6': '.s3.OwnerOverride',
      '10': 'owner'
    },
  ],
};

/// Descriptor for `AccessControlTranslation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessControlTranslationDescriptor =
    $convert.base64Decode(
        'ChhBY2Nlc3NDb250cm9sVHJhbnNsYXRpb24SKwoFb3duZXIY9fyK2QEgASgOMhEuczMuT3duZX'
        'JPdmVycmlkZVIFb3duZXI=');

@$core.Deprecated('Use accessDeniedDescriptor instead')
const AccessDenied$json = {
  '1': 'AccessDenied',
};

/// Descriptor for `AccessDenied`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDeniedDescriptor =
    $convert.base64Decode('CgxBY2Nlc3NEZW5pZWQ=');

@$core.Deprecated('Use analyticsAndOperatorDescriptor instead')
const AnalyticsAndOperator$json = {
  '1': 'AnalyticsAndOperator',
  '2': [
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `AnalyticsAndOperator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyticsAndOperatorDescriptor = $convert.base64Decode(
    'ChRBbmFseXRpY3NBbmRPcGVyYXRvchIaCgZwcmVmaXgY6rPTggEgASgJUgZwcmVmaXgSHwoEdG'
    'FncxjBwfa1ASADKAsyBy5zMy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use analyticsConfigurationDescriptor instead')
const AnalyticsConfiguration$json = {
  '1': 'AnalyticsConfiguration',
  '2': [
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.AnalyticsFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'storageclassanalysis',
      '3': 242394117,
      '4': 1,
      '5': 11,
      '6': '.s3.StorageClassAnalysis',
      '10': 'storageclassanalysis'
    },
  ],
};

/// Descriptor for `AnalyticsConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyticsConfigurationDescriptor = $convert.base64Decode(
    'ChZBbmFseXRpY3NDb25maWd1cmF0aW9uEi8KBmZpbHRlchiYgaelASABKAsyEy5zMy5BbmFseX'
    'RpY3NGaWx0ZXJSBmZpbHRlchISCgJpZBiB8qK3ASABKAlSAmlkEk8KFHN0b3JhZ2VjbGFzc2Fu'
    'YWx5c2lzGIXIynMgASgLMhguczMuU3RvcmFnZUNsYXNzQW5hbHlzaXNSFHN0b3JhZ2VjbGFzc2'
    'FuYWx5c2lz');

@$core.Deprecated('Use analyticsExportDestinationDescriptor instead')
const AnalyticsExportDestination$json = {
  '1': 'AnalyticsExportDestination',
  '2': [
    {
      '1': 's3bucketdestination',
      '3': 94381700,
      '4': 1,
      '5': 11,
      '6': '.s3.AnalyticsS3BucketDestination',
      '10': 's3bucketdestination'
    },
  ],
};

/// Descriptor for `AnalyticsExportDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyticsExportDestinationDescriptor =
    $convert.base64Decode(
        'ChpBbmFseXRpY3NFeHBvcnREZXN0aW5hdGlvbhJVChNzM2J1Y2tldGRlc3RpbmF0aW9uGITNgC'
        '0gASgLMiAuczMuQW5hbHl0aWNzUzNCdWNrZXREZXN0aW5hdGlvblITczNidWNrZXRkZXN0aW5h'
        'dGlvbg==');

@$core.Deprecated('Use analyticsFilterDescriptor instead')
const AnalyticsFilter$json = {
  '1': 'AnalyticsFilter',
  '2': [
    {
      '1': 'and',
      '3': 297135431,
      '4': 1,
      '5': 11,
      '6': '.s3.AnalyticsAndOperator',
      '10': 'and'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {'1': 'tag', '3': 411259956, '4': 1, '5': 11, '6': '.s3.Tag', '10': 'tag'},
  ],
};

/// Descriptor for `AnalyticsFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyticsFilterDescriptor = $convert.base64Decode(
    'Cg9BbmFseXRpY3NGaWx0ZXISLgoDYW5kGMfa140BIAEoCzIYLnMzLkFuYWx5dGljc0FuZE9wZX'
    'JhdG9yUgNhbmQSGgoGcHJlZml4GOqz04IBIAEoCVIGcHJlZml4Eh0KA3RhZxi0qI3EASABKAsy'
    'By5zMy5UYWdSA3RhZw==');

@$core.Deprecated('Use analyticsS3BucketDestinationDescriptor instead')
const AnalyticsS3BucketDestination$json = {
  '1': 'AnalyticsS3BucketDestination',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bucketaccountid',
      '3': 235438976,
      '4': 1,
      '5': 9,
      '10': 'bucketaccountid'
    },
    {
      '1': 'format',
      '3': 531693427,
      '4': 1,
      '5': 14,
      '6': '.s3.AnalyticsS3ExportFileFormat',
      '10': 'format'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
  ],
};

/// Descriptor for `AnalyticsS3BucketDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyticsS3BucketDestinationDescriptor = $convert.base64Decode(
    'ChxBbmFseXRpY3NTM0J1Y2tldERlc3RpbmF0aW9uEhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2'
    'V0EisKD2J1Y2tldGFjY291bnRpZBiAh6JwIAEoCVIPYnVja2V0YWNjb3VudGlkEjsKBmZvcm1h'
    'dBjz/sP9ASABKA4yHy5zMy5BbmFseXRpY3NTM0V4cG9ydEZpbGVGb3JtYXRSBmZvcm1hdBIaCg'
    'ZwcmVmaXgY6rPTggEgASgJUgZwcmVmaXg=');

@$core.Deprecated('Use blockedEncryptionTypesDescriptor instead')
const BlockedEncryptionTypes$json = {
  '1': 'BlockedEncryptionTypes',
  '2': [
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 3,
      '5': 14,
      '6': '.s3.EncryptionType',
      '10': 'encryptiontype'
    },
  ],
};

/// Descriptor for `BlockedEncryptionTypes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockedEncryptionTypesDescriptor =
    $convert.base64Decode(
        'ChZCbG9ja2VkRW5jcnlwdGlvblR5cGVzEj0KDmVuY3J5cHRpb250eXBlGLXf8X0gAygOMhIucz'
        'MuRW5jcnlwdGlvblR5cGVSDmVuY3J5cHRpb250eXBl');

@$core.Deprecated('Use bucketDescriptor instead')
const Bucket$json = {
  '1': 'Bucket',
  '2': [
    {'1': 'bucketarn', '3': 255683899, '4': 1, '5': 9, '10': 'bucketarn'},
    {'1': 'bucketregion', '3': 309298816, '4': 1, '5': 9, '10': 'bucketregion'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `Bucket`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bucketDescriptor = $convert.base64Decode(
    'CgZCdWNrZXQSHwoJYnVja2V0YXJuGLva9XkgASgJUglidWNrZXRhcm4SJgoMYnVja2V0cmVnaW'
    '9uGICNvpMBIAEoCVIMYnVja2V0cmVnaW9uEiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNy'
    'ZWF0aW9uZGF0ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use bucketAlreadyExistsDescriptor instead')
const BucketAlreadyExists$json = {
  '1': 'BucketAlreadyExists',
};

/// Descriptor for `BucketAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bucketAlreadyExistsDescriptor =
    $convert.base64Decode('ChNCdWNrZXRBbHJlYWR5RXhpc3Rz');

@$core.Deprecated('Use bucketAlreadyOwnedByYouDescriptor instead')
const BucketAlreadyOwnedByYou$json = {
  '1': 'BucketAlreadyOwnedByYou',
};

/// Descriptor for `BucketAlreadyOwnedByYou`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bucketAlreadyOwnedByYouDescriptor =
    $convert.base64Decode('ChdCdWNrZXRBbHJlYWR5T3duZWRCeVlvdQ==');

@$core.Deprecated('Use bucketInfoDescriptor instead')
const BucketInfo$json = {
  '1': 'BucketInfo',
  '2': [
    {
      '1': 'dataredundancy',
      '3': 526607767,
      '4': 1,
      '5': 14,
      '6': '.s3.DataRedundancy',
      '10': 'dataredundancy'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `BucketInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bucketInfoDescriptor = $convert.base64Decode(
    'CgpCdWNrZXRJbmZvEj4KDmRhdGFyZWR1bmRhbmN5GJfLjfsBIAEoDjISLnMzLkRhdGFSZWR1bm'
    'RhbmN5Ug5kYXRhcmVkdW5kYW5jeRImCgR0eXBlGO6g14oBIAEoDjIOLnMzLkJ1Y2tldFR5cGVS'
    'BHR5cGU=');

@$core.Deprecated('Use bucketLifecycleConfigurationDescriptor instead')
const BucketLifecycleConfiguration$json = {
  '1': 'BucketLifecycleConfiguration',
  '2': [
    {
      '1': 'rules',
      '3': 42675585,
      '4': 3,
      '5': 11,
      '6': '.s3.LifecycleRule',
      '10': 'rules'
    },
  ],
};

/// Descriptor for `BucketLifecycleConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bucketLifecycleConfigurationDescriptor =
    $convert.base64Decode(
        'ChxCdWNrZXRMaWZlY3ljbGVDb25maWd1cmF0aW9uEioKBXJ1bGVzGIHbrBQgAygLMhEuczMuTG'
        'lmZWN5Y2xlUnVsZVIFcnVsZXM=');

@$core.Deprecated('Use bucketLoggingStatusDescriptor instead')
const BucketLoggingStatus$json = {
  '1': 'BucketLoggingStatus',
  '2': [
    {
      '1': 'loggingenabled',
      '3': 119537244,
      '4': 1,
      '5': 11,
      '6': '.s3.LoggingEnabled',
      '10': 'loggingenabled'
    },
  ],
};

/// Descriptor for `BucketLoggingStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bucketLoggingStatusDescriptor = $convert.base64Decode(
    'ChNCdWNrZXRMb2dnaW5nU3RhdHVzEj0KDmxvZ2dpbmdlbmFibGVkGNz8/zggASgLMhIuczMuTG'
    '9nZ2luZ0VuYWJsZWRSDmxvZ2dpbmdlbmFibGVk');

@$core.Deprecated('Use cORSConfigurationDescriptor instead')
const CORSConfiguration$json = {
  '1': 'CORSConfiguration',
  '2': [
    {
      '1': 'corsrules',
      '3': 280917788,
      '4': 3,
      '5': 11,
      '6': '.s3.CORSRule',
      '10': 'corsrules'
    },
  ],
};

/// Descriptor for `CORSConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cORSConfigurationDescriptor = $convert.base64Decode(
    'ChFDT1JTQ29uZmlndXJhdGlvbhIuCgljb3JzcnVsZXMYnO75hQEgAygLMgwuczMuQ09SU1J1bG'
    'VSCWNvcnNydWxlcw==');

@$core.Deprecated('Use cORSRuleDescriptor instead')
const CORSRule$json = {
  '1': 'CORSRule',
  '2': [
    {
      '1': 'allowedheaders',
      '3': 63126900,
      '4': 3,
      '5': 9,
      '10': 'allowedheaders'
    },
    {
      '1': 'allowedmethods',
      '3': 56383476,
      '4': 3,
      '5': 9,
      '10': 'allowedmethods'
    },
    {
      '1': 'allowedorigins',
      '3': 19474107,
      '4': 3,
      '5': 9,
      '10': 'allowedorigins'
    },
    {
      '1': 'exposeheaders',
      '3': 290364554,
      '4': 3,
      '5': 9,
      '10': 'exposeheaders'
    },
    {'1': 'id', '3': 384363361, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'maxageseconds',
      '3': 315057526,
      '4': 1,
      '5': 5,
      '10': 'maxageseconds'
    },
  ],
};

/// Descriptor for `CORSRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cORSRuleDescriptor = $convert.base64Decode(
    'CghDT1JTUnVsZRIpCg5hbGxvd2VkaGVhZGVycxj0+oweIAMoCVIOYWxsb3dlZGhlYWRlcnMSKQ'
    'oOYWxsb3dlZG1ldGhvZHMY9K/xGiADKAlSDmFsbG93ZWRtZXRob2RzEikKDmFsbG93ZWRvcmln'
    'aW5zGLvNpAkgAygJUg5hbGxvd2Vkb3JpZ2lucxIoCg1leHBvc2VoZWFkZXJzGIq5uooBIAMoCV'
    'INZXhwb3NlaGVhZGVycxISCgJpZBjh1qO3ASABKAlSAmlkEigKDW1heGFnZXNlY29uZHMY9sqd'
    'lgEgASgFUg1tYXhhZ2VzZWNvbmRz');

@$core.Deprecated('Use cSVInputDescriptor instead')
const CSVInput$json = {
  '1': 'CSVInput',
  '2': [
    {
      '1': 'allowquotedrecorddelimiter',
      '3': 268967449,
      '4': 1,
      '5': 8,
      '10': 'allowquotedrecorddelimiter'
    },
    {'1': 'comments', '3': 307768056, '4': 1, '5': 9, '10': 'comments'},
    {
      '1': 'fielddelimiter',
      '3': 88917725,
      '4': 1,
      '5': 9,
      '10': 'fielddelimiter'
    },
    {
      '1': 'fileheaderinfo',
      '3': 43578219,
      '4': 1,
      '5': 14,
      '6': '.s3.FileHeaderInfo',
      '10': 'fileheaderinfo'
    },
    {
      '1': 'quotecharacter',
      '3': 24331231,
      '4': 1,
      '5': 9,
      '10': 'quotecharacter'
    },
    {
      '1': 'quoteescapecharacter',
      '3': 378770328,
      '4': 1,
      '5': 9,
      '10': 'quoteescapecharacter'
    },
    {
      '1': 'recorddelimiter',
      '3': 299337270,
      '4': 1,
      '5': 9,
      '10': 'recorddelimiter'
    },
  ],
};

/// Descriptor for `CSVInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cSVInputDescriptor = $convert.base64Decode(
    'CghDU1ZJbnB1dBJCChphbGxvd3F1b3RlZHJlY29yZGRlbGltaXRlchiZvKCAASABKAhSGmFsbG'
    '93cXVvdGVkcmVjb3JkZGVsaW1pdGVyEh4KCGNvbW1lbnRzGPjV4JIBIAEoCVIIY29tbWVudHMS'
    'KQoOZmllbGRkZWxpbWl0ZXIY3Y2zKiABKAlSDmZpZWxkZGVsaW1pdGVyEj0KDmZpbGVoZWFkZX'
    'JpbmZvGOvm4xQgASgOMhIuczMuRmlsZUhlYWRlckluZm9SDmZpbGVoZWFkZXJpbmZvEikKDnF1'
    'b3RlY2hhcmFjdGVyGN+HzQsgASgJUg5xdW90ZWNoYXJhY3RlchI2ChRxdW90ZWVzY2FwZWNoYX'
    'JhY3RlchiYp860ASABKAlSFHF1b3RlZXNjYXBlY2hhcmFjdGVyEiwKD3JlY29yZGRlbGltaXRl'
    'chi2jN6OASABKAlSD3JlY29yZGRlbGltaXRlcg==');

@$core.Deprecated('Use cSVOutputDescriptor instead')
const CSVOutput$json = {
  '1': 'CSVOutput',
  '2': [
    {
      '1': 'fielddelimiter',
      '3': 88917725,
      '4': 1,
      '5': 9,
      '10': 'fielddelimiter'
    },
    {
      '1': 'quotecharacter',
      '3': 24331231,
      '4': 1,
      '5': 9,
      '10': 'quotecharacter'
    },
    {
      '1': 'quoteescapecharacter',
      '3': 378770328,
      '4': 1,
      '5': 9,
      '10': 'quoteescapecharacter'
    },
    {
      '1': 'quotefields',
      '3': 136837903,
      '4': 1,
      '5': 14,
      '6': '.s3.QuoteFields',
      '10': 'quotefields'
    },
    {
      '1': 'recorddelimiter',
      '3': 299337270,
      '4': 1,
      '5': 9,
      '10': 'recorddelimiter'
    },
  ],
};

/// Descriptor for `CSVOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cSVOutputDescriptor = $convert.base64Decode(
    'CglDU1ZPdXRwdXQSKQoOZmllbGRkZWxpbWl0ZXIY3Y2zKiABKAlSDmZpZWxkZGVsaW1pdGVyEi'
    'kKDnF1b3RlY2hhcmFjdGVyGN+HzQsgASgJUg5xdW90ZWNoYXJhY3RlchI2ChRxdW90ZWVzY2Fw'
    'ZWNoYXJhY3RlchiYp860ASABKAlSFHF1b3RlZXNjYXBlY2hhcmFjdGVyEjQKC3F1b3RlZmllbG'
    'RzGI/2n0EgASgOMg8uczMuUXVvdGVGaWVsZHNSC3F1b3RlZmllbGRzEiwKD3JlY29yZGRlbGlt'
    'aXRlchi2jN6OASABKAlSD3JlY29yZGRlbGltaXRlcg==');

@$core.Deprecated('Use checksumDescriptor instead')
const Checksum$json = {
  '1': 'Checksum',
  '2': [
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
  ],
};

/// Descriptor for `Checksum`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checksumDescriptor = $convert.base64Decode(
    'CghDaGVja3N1bRInCg1jaGVja3N1bWNyYzMyGMKjzTMgASgJUg1jaGVja3N1bWNyYzMyEikKDm'
    'NoZWNrc3VtY3JjMzJjGKefpUwgASgJUg5jaGVja3N1bWNyYzMyYxIvChFjaGVja3N1bWNyYzY0'
    'bnZtZRjNvIs4IAEoCVIRY2hlY2tzdW1jcmM2NG52bWUSJgoMY2hlY2tzdW1zaGExGMTs4IoBIA'
    'EoCVIMY2hlY2tzdW1zaGExEikKDmNoZWNrc3Vtc2hhMjU2GL753EQgASgJUg5jaGVja3N1bXNo'
    'YTI1NhI3CgxjaGVja3N1bXR5cGUYw77ZLiABKA4yEC5zMy5DaGVja3N1bVR5cGVSDGNoZWNrc3'
    'VtdHlwZQ==');

@$core.Deprecated('Use commonPrefixDescriptor instead')
const CommonPrefix$json = {
  '1': 'CommonPrefix',
  '2': [
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
  ],
};

/// Descriptor for `CommonPrefix`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List commonPrefixDescriptor = $convert
    .base64Decode('CgxDb21tb25QcmVmaXgSGgoGcHJlZml4GOqz04IBIAEoCVIGcHJlZml4');

@$core.Deprecated('Use completeMultipartUploadOutputDescriptor instead')
const CompleteMultipartUploadOutput$json = {
  '1': 'CompleteMultipartUploadOutput',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `CompleteMultipartUploadOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completeMultipartUploadOutputDescriptor = $convert.base64Decode(
    'Ch1Db21wbGV0ZU11bHRpcGFydFVwbG9hZE91dHB1dBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2'
    'tldBIuChBidWNrZXRrZXllbmFibGVkGOyijM8BIAEoCFIQYnVja2V0a2V5ZW5hYmxlZBInCg1j'
    'aGVja3N1bWNyYzMyGMKjzTMgASgJUg1jaGVja3N1bWNyYzMyEikKDmNoZWNrc3VtY3JjMzJjGK'
    'efpUwgASgJUg5jaGVja3N1bWNyYzMyYxIvChFjaGVja3N1bWNyYzY0bnZtZRjNvIs4IAEoCVIR'
    'Y2hlY2tzdW1jcmM2NG52bWUSJgoMY2hlY2tzdW1zaGExGMTs4IoBIAEoCVIMY2hlY2tzdW1zaG'
    'ExEikKDmNoZWNrc3Vtc2hhMjU2GL753EQgASgJUg5jaGVja3N1bXNoYTI1NhI3CgxjaGVja3N1'
    'bXR5cGUYw77ZLiABKA4yEC5zMy5DaGVja3N1bVR5cGVSDGNoZWNrc3VtdHlwZRIWCgRldGFnGI'
    'Hfs5UBIAEoCVIEZXRhZxIhCgpleHBpcmF0aW9uGImpn3UgASgJUgpleHBpcmF0aW9uEhMKA2tl'
    'eRiNkutoIAEoCVIDa2V5Eh4KCGxvY2F0aW9uGMebgt4BIAEoCVIIbG9jYXRpb24SPgoOcmVxdW'
    'VzdGNoYXJnZWQYk9CruQEgASgOMhIuczMuUmVxdWVzdENoYXJnZWRSDnJlcXVlc3RjaGFyZ2Vk'
    'EiQKC3NzZWttc2tleWlkGNiI1NQBIAEoCVILc3Nla21za2V5aWQSTwoUc2VydmVyc2lkZWVuY3'
    'J5cHRpb24YsY6fBCABKA4yGC5zMy5TZXJ2ZXJTaWRlRW5jcnlwdGlvblIUc2VydmVyc2lkZWVu'
    'Y3J5cHRpb24SIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlk');

@$core.Deprecated('Use completeMultipartUploadRequestDescriptor instead')
const CompleteMultipartUploadRequest$json = {
  '1': 'CompleteMultipartUploadRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'ifnonematch', '3': 231211830, '4': 1, '5': 9, '10': 'ifnonematch'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'mpuobjectsize',
      '3': 522430190,
      '4': 1,
      '5': 3,
      '10': 'mpuobjectsize'
    },
    {
      '1': 'multipartupload',
      '3': 362466695,
      '4': 1,
      '5': 11,
      '6': '.s3.CompletedMultipartUpload',
      '8': {},
      '10': 'multipartupload'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'uploadid', '3': 449040722, '4': 1, '5': 9, '10': 'uploadid'},
  ],
};

/// Descriptor for `CompleteMultipartUploadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completeMultipartUploadRequestDescriptor = $convert.base64Decode(
    'Ch5Db21wbGV0ZU11bHRpcGFydFVwbG9hZFJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidW'
    'NrZXQSJwoNY2hlY2tzdW1jcmMzMhjCo80zIAEoCVINY2hlY2tzdW1jcmMzMhIpCg5jaGVja3N1'
    'bWNyYzMyYxinn6VMIAEoCVIOY2hlY2tzdW1jcmMzMmMSLwoRY2hlY2tzdW1jcmM2NG52bWUYzb'
    'yLOCABKAlSEWNoZWNrc3VtY3JjNjRudm1lEiYKDGNoZWNrc3Vtc2hhMRjE7OCKASABKAlSDGNo'
    'ZWNrc3Vtc2hhMRIpCg5jaGVja3N1bXNoYTI1Nhi++dxEIAEoCVIOY2hlY2tzdW1zaGEyNTYSNw'
    'oMY2hlY2tzdW10eXBlGMO+2S4gASgOMhAuczMuQ2hlY2tzdW1UeXBlUgxjaGVja3N1bXR5cGUS'
    'MwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchIbCg'
    'dpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNoEiMKC2lmbm9uZW1hdGNoGLaGoG4gASgJUgtpZm5v'
    'bmVtYXRjaBITCgNrZXkYjZLraCABKAlSA2tleRIoCg1tcHVvYmplY3RzaXplGO7NjvkBIAEoA1'
    'INbXB1b2JqZWN0c2l6ZRJQCg9tdWx0aXBhcnR1cGxvYWQYh5vrrAEgASgLMhwuczMuQ29tcGxl'
    'dGVkTXVsdGlwYXJ0VXBsb2FkQgSItRgBUg9tdWx0aXBhcnR1cGxvYWQSOAoMcmVxdWVzdHBheW'
    'VyGKTm4fUBIAEoDjIQLnMzLlJlcXVlc3RQYXllclIMcmVxdWVzdHBheWVyEjUKFHNzZWN1c3Rv'
    'bWVyYWxnb3JpdGhtGNDJgSsgASgJUhRzc2VjdXN0b21lcmFsZ29yaXRobRIpCg5zc2VjdXN0b2'
    '1lcmtleRia/vQ7IAEoCVIOc3NlY3VzdG9tZXJrZXkSLgoRc3NlY3VzdG9tZXJrZXltZDUY6NEX'
    'IAEoCVIRc3NlY3VzdG9tZXJrZXltZDUSHgoIdXBsb2FkaWQY0qKP1gEgASgJUgh1cGxvYWRpZA'
    '==');

@$core.Deprecated('Use completedMultipartUploadDescriptor instead')
const CompletedMultipartUpload$json = {
  '1': 'CompletedMultipartUpload',
  '2': [
    {
      '1': 'parts',
      '3': 213028806,
      '4': 3,
      '5': 11,
      '6': '.s3.CompletedPart',
      '10': 'parts'
    },
  ],
};

/// Descriptor for `CompletedMultipartUpload`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completedMultipartUploadDescriptor =
    $convert.base64Decode(
        'ChhDb21wbGV0ZWRNdWx0aXBhcnRVcGxvYWQSKgoFcGFydHMYxp/KZSADKAsyES5zMy5Db21wbG'
        'V0ZWRQYXJ0UgVwYXJ0cw==');

@$core.Deprecated('Use completedPartDescriptor instead')
const CompletedPart$json = {
  '1': 'CompletedPart',
  '2': [
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'partnumber', '3': 372082310, '4': 1, '5': 5, '10': 'partnumber'},
  ],
};

/// Descriptor for `CompletedPart`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List completedPartDescriptor = $convert.base64Decode(
    'Cg1Db21wbGV0ZWRQYXJ0EicKDWNoZWNrc3VtY3JjMzIYwqPNMyABKAlSDWNoZWNrc3VtY3JjMz'
    'ISKQoOY2hlY2tzdW1jcmMzMmMYp5+lTCABKAlSDmNoZWNrc3VtY3JjMzJjEi8KEWNoZWNrc3Vt'
    'Y3JjNjRudm1lGM28izggASgJUhFjaGVja3N1bWNyYzY0bnZtZRImCgxjaGVja3N1bXNoYTEYxO'
    'zgigEgASgJUgxjaGVja3N1bXNoYTESKQoOY2hlY2tzdW1zaGEyNTYYvvncRCABKAlSDmNoZWNr'
    'c3Vtc2hhMjU2EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEiIKCnBhcnRudW1iZXIYho22sQEgAS'
    'gFUgpwYXJ0bnVtYmVy');

@$core.Deprecated('Use conditionDescriptor instead')
const Condition$json = {
  '1': 'Condition',
  '2': [
    {
      '1': 'httperrorcodereturnedequals',
      '3': 172790163,
      '4': 1,
      '5': 9,
      '10': 'httperrorcodereturnedequals'
    },
    {
      '1': 'keyprefixequals',
      '3': 334847946,
      '4': 1,
      '5': 9,
      '10': 'keyprefixequals'
    },
  ],
};

/// Descriptor for `Condition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conditionDescriptor = $convert.base64Decode(
    'CglDb25kaXRpb24SQwobaHR0cGVycm9yY29kZXJldHVybmVkZXF1YWxzGJOjslIgASgJUhtodH'
    'RwZXJyb3Jjb2RlcmV0dXJuZWRlcXVhbHMSLAoPa2V5cHJlZml4ZXF1YWxzGMq/1Z8BIAEoCVIP'
    'a2V5cHJlZml4ZXF1YWxz');

@$core.Deprecated('Use continuationEventDescriptor instead')
const ContinuationEvent$json = {
  '1': 'ContinuationEvent',
};

/// Descriptor for `ContinuationEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List continuationEventDescriptor =
    $convert.base64Decode('ChFDb250aW51YXRpb25FdmVudA==');

@$core.Deprecated('Use copyObjectOutputDescriptor instead')
const CopyObjectOutput$json = {
  '1': 'CopyObjectOutput',
  '2': [
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {
      '1': 'copyobjectresult',
      '3': 84805307,
      '4': 1,
      '5': 11,
      '6': '.s3.CopyObjectResult',
      '8': {},
      '10': 'copyobjectresult'
    },
    {
      '1': 'copysourceversionid',
      '3': 257134375,
      '4': 1,
      '5': 9,
      '10': 'copysourceversionid'
    },
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {
      '1': 'ssekmsencryptioncontext',
      '3': 149030970,
      '4': 1,
      '5': 9,
      '10': 'ssekmsencryptioncontext'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `CopyObjectOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List copyObjectOutputDescriptor = $convert.base64Decode(
    'ChBDb3B5T2JqZWN0T3V0cHV0Ei4KEGJ1Y2tldGtleWVuYWJsZWQY7KKMzwEgASgIUhBidWNrZX'
    'RrZXllbmFibGVkEkkKEGNvcHlvYmplY3RyZXN1bHQYu424KCABKAsyFC5zMy5Db3B5T2JqZWN0'
    'UmVzdWx0QgSItRgBUhBjb3B5b2JqZWN0cmVzdWx0EjMKE2NvcHlzb3VyY2V2ZXJzaW9uaWQYp5'
    '7OeiABKAlSE2NvcHlzb3VyY2V2ZXJzaW9uaWQSIQoKZXhwaXJhdGlvbhiJqZ91IAEoCVIKZXhw'
    'aXJhdGlvbhI+Cg5yZXF1ZXN0Y2hhcmdlZBiT0Ku5ASABKA4yEi5zMy5SZXF1ZXN0Q2hhcmdlZF'
    'IOcmVxdWVzdGNoYXJnZWQSNQoUc3NlY3VzdG9tZXJhbGdvcml0aG0Y0MmBKyABKAlSFHNzZWN1'
    'c3RvbWVyYWxnb3JpdGhtEi4KEXNzZWN1c3RvbWVya2V5bWQ1GOjRFyABKAlSEXNzZWN1c3RvbW'
    'Vya2V5bWQ1EjsKF3NzZWttc2VuY3J5cHRpb25jb250ZXh0GLqQiEcgASgJUhdzc2VrbXNlbmNy'
    'eXB0aW9uY29udGV4dBIkCgtzc2VrbXNrZXlpZBjYiNTUASABKAlSC3NzZWttc2tleWlkEk8KFH'
    'NlcnZlcnNpZGVlbmNyeXB0aW9uGLGOnwQgASgOMhguczMuU2VydmVyU2lkZUVuY3J5cHRpb25S'
    'FHNlcnZlcnNpZGVlbmNyeXB0aW9uEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZA'
    '==');

@$core.Deprecated('Use copyObjectRequestDescriptor instead')
const CopyObjectRequest$json = {
  '1': 'CopyObjectRequest',
  '2': [
    {
      '1': 'acl',
      '3': 394696836,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectCannedACL',
      '10': 'acl'
    },
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {'1': 'cachecontrol', '3': 288966655, '4': 1, '5': 9, '10': 'cachecontrol'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'contentdisposition',
      '3': 120040130,
      '4': 1,
      '5': 9,
      '10': 'contentdisposition'
    },
    {
      '1': 'contentencoding',
      '3': 317106228,
      '4': 1,
      '5': 9,
      '10': 'contentencoding'
    },
    {
      '1': 'contentlanguage',
      '3': 108485649,
      '4': 1,
      '5': 9,
      '10': 'contentlanguage'
    },
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'copysource', '3': 152315650, '4': 1, '5': 9, '10': 'copysource'},
    {
      '1': 'copysourceifmatch',
      '3': 512868332,
      '4': 1,
      '5': 9,
      '10': 'copysourceifmatch'
    },
    {
      '1': 'copysourceifmodifiedsince',
      '3': 301576082,
      '4': 1,
      '5': 9,
      '10': 'copysourceifmodifiedsince'
    },
    {
      '1': 'copysourceifnonematch',
      '3': 237403010,
      '4': 1,
      '5': 9,
      '10': 'copysourceifnonematch'
    },
    {
      '1': 'copysourceifunmodifiedsince',
      '3': 130556417,
      '4': 1,
      '5': 9,
      '10': 'copysourceifunmodifiedsince'
    },
    {
      '1': 'copysourcessecustomeralgorithm',
      '3': 465017156,
      '4': 1,
      '5': 9,
      '10': 'copysourcessecustomeralgorithm'
    },
    {
      '1': 'copysourcessecustomerkey',
      '3': 289313102,
      '4': 1,
      '5': 9,
      '10': 'copysourcessecustomerkey'
    },
    {
      '1': 'copysourcessecustomerkeymd5',
      '3': 408962492,
      '4': 1,
      '5': 9,
      '10': 'copysourcessecustomerkeymd5'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'expectedsourcebucketowner',
      '3': 283238582,
      '4': 1,
      '5': 9,
      '10': 'expectedsourcebucketowner'
    },
    {'1': 'expires', '3': 128582948, '4': 1, '5': 9, '10': 'expires'},
    {
      '1': 'grantfullcontrol',
      '3': 102486874,
      '4': 1,
      '5': 9,
      '10': 'grantfullcontrol'
    },
    {'1': 'grantread', '3': 4518126, '4': 1, '5': 9, '10': 'grantread'},
    {'1': 'grantreadacp', '3': 208410960, '4': 1, '5': 9, '10': 'grantreadacp'},
    {
      '1': 'grantwriteacp',
      '3': 157594957,
      '4': 1,
      '5': 9,
      '10': 'grantwriteacp'
    },
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'ifnonematch', '3': 231211830, '4': 1, '5': 9, '10': 'ifnonematch'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 3,
      '5': 11,
      '6': '.s3.CopyObjectRequest.MetadataEntry',
      '10': 'metadata'
    },
    {
      '1': 'metadatadirective',
      '3': 282181534,
      '4': 1,
      '5': 14,
      '6': '.s3.MetadataDirective',
      '10': 'metadatadirective'
    },
    {
      '1': 'objectlocklegalholdstatus',
      '3': 536561974,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockLegalHoldStatus',
      '10': 'objectlocklegalholdstatus'
    },
    {
      '1': 'objectlockmode',
      '3': 189255203,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockMode',
      '10': 'objectlockmode'
    },
    {
      '1': 'objectlockretainuntildate',
      '3': 264584249,
      '4': 1,
      '5': 9,
      '10': 'objectlockretainuntildate'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {
      '1': 'ssekmsencryptioncontext',
      '3': 149030970,
      '4': 1,
      '5': 9,
      '10': 'ssekmsencryptioncontext'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'tagging', '3': 33436541, '4': 1, '5': 9, '10': 'tagging'},
    {
      '1': 'taggingdirective',
      '3': 498516426,
      '4': 1,
      '5': 14,
      '6': '.s3.TaggingDirective',
      '10': 'taggingdirective'
    },
    {
      '1': 'websiteredirectlocation',
      '3': 71844662,
      '4': 1,
      '5': 9,
      '10': 'websiteredirectlocation'
    },
  ],
  '3': [CopyObjectRequest_MetadataEntry$json],
};

@$core.Deprecated('Use copyObjectRequestDescriptor instead')
const CopyObjectRequest_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CopyObjectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List copyObjectRequestDescriptor = $convert.base64Decode(
    'ChFDb3B5T2JqZWN0UmVxdWVzdBIpCgNhY2wYhLGavAEgASgOMhMuczMuT2JqZWN0Q2FubmVkQU'
    'NMUgNhY2wSGQoGYnVja2V0GNjquBogASgJUgZidWNrZXQSLgoQYnVja2V0a2V5ZW5hYmxlZBjs'
    'oozPASABKAhSEGJ1Y2tldGtleWVuYWJsZWQSJgoMY2FjaGVjb250cm9sGP+P5YkBIAEoCVIMY2'
    'FjaGVjb250cm9sEkYKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMuQ2hlY2tzdW1B'
    'bGdvcml0aG1SEWNoZWNrc3VtYWxnb3JpdGhtEjEKEmNvbnRlbnRkaXNwb3NpdGlvbhjC1Z45IA'
    'EoCVISY29udGVudGRpc3Bvc2l0aW9uEiwKD2NvbnRlbnRlbmNvZGluZxi00JqXASABKAlSD2Nv'
    'bnRlbnRlbmNvZGluZxIrCg9jb250ZW50bGFuZ3VhZ2UYkbjdMyABKAlSD2NvbnRlbnRsYW5ndW'
    'FnZRIkCgtjb250ZW50dHlwZRiT1eieASABKAlSC2NvbnRlbnR0eXBlEiEKCmNvcHlzb3VyY2UY'
    'gs7QSCABKAlSCmNvcHlzb3VyY2USMAoRY29weXNvdXJjZWlmbWF0Y2gY7P/G9AEgASgJUhFjb3'
    'B5c291cmNlaWZtYXRjaBJAChljb3B5c291cmNlaWZtb2RpZmllZHNpbmNlGJLf5o8BIAEoCVIZ'
    'Y29weXNvdXJjZWlmbW9kaWZpZWRzaW5jZRI3ChVjb3B5c291cmNlaWZub25lbWF0Y2gYgveZcS'
    'ABKAlSFWNvcHlzb3VyY2VpZm5vbmVtYXRjaBJDChtjb3B5c291cmNlaWZ1bm1vZGlmaWVkc2lu'
    'Y2UYgcSgPiABKAlSG2NvcHlzb3VyY2VpZnVubW9kaWZpZWRzaW5jZRJKCh5jb3B5c291cmNlc3'
    'NlY3VzdG9tZXJhbGdvcml0aG0YxLLe3QEgASgJUh5jb3B5c291cmNlc3NlY3VzdG9tZXJhbGdv'
    'cml0aG0SPgoYY29weXNvdXJjZXNzZWN1c3RvbWVya2V5GM6i+okBIAEoCVIYY29weXNvdXJjZX'
    'NzZWN1c3RvbWVya2V5EkQKG2NvcHlzb3VyY2Vzc2VjdXN0b21lcmtleW1kNRi8i4HDASABKAlS'
    'G2NvcHlzb3VyY2Vzc2VjdXN0b21lcmtleW1kNRIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D'
    '4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEkAKGWV4cGVjdGVkc291cmNlYnVja2V0b3duZXIY'
    'tsGHhwEgASgJUhlleHBlY3RlZHNvdXJjZWJ1Y2tldG93bmVyEhsKB2V4cGlyZXMYpIqoPSABKA'
    'lSB2V4cGlyZXMSLQoQZ3JhbnRmdWxsY29udHJvbBjapu8wIAEoCVIQZ3JhbnRmdWxsY29udHJv'
    'bBIfCglncmFudHJlYWQY7uGTAiABKAlSCWdyYW50cmVhZBIlCgxncmFudHJlYWRhY3AY0LKwYy'
    'ABKAlSDGdyYW50cmVhZGFjcBInCg1ncmFudHdyaXRlYWNwGM3qkksgASgJUg1ncmFudHdyaXRl'
    'YWNwEhsKB2lmbWF0Y2gY0Ja3LCABKAlSB2lmbWF0Y2gSIwoLaWZub25lbWF0Y2gYtoagbiABKA'
    'lSC2lmbm9uZW1hdGNoEhMKA2tleRiNkutoIAEoCVIDa2V5EkMKCG1ldGFkYXRhGOHij+ABIAMo'
    'CzIjLnMzLkNvcHlPYmplY3RSZXF1ZXN0Lk1ldGFkYXRhRW50cnlSCG1ldGFkYXRhEkcKEW1ldG'
    'FkYXRhZGlyZWN0aXZlGJ7/xoYBIAEoDjIVLnMzLk1ldGFkYXRhRGlyZWN0aXZlUhFtZXRhZGF0'
    'YWRpcmVjdGl2ZRJfChlvYmplY3Rsb2NrbGVnYWxob2xkc3RhdHVzGLaS7f8BIAEoDjIdLnMzLk'
    '9iamVjdExvY2tMZWdhbEhvbGRTdGF0dXNSGW9iamVjdGxvY2tsZWdhbGhvbGRzdGF0dXMSPQoO'
    'b2JqZWN0bG9ja21vZGUYo5yfWiABKA4yEi5zMy5PYmplY3RMb2NrTW9kZVIOb2JqZWN0bG9ja2'
    '1vZGUSPwoZb2JqZWN0bG9ja3JldGFpbnVudGlsZGF0ZRi5+JR+IAEoCVIZb2JqZWN0bG9ja3Jl'
    'dGFpbnVudGlsZGF0ZRI4CgxyZXF1ZXN0cGF5ZXIYpObh9QEgASgOMhAuczMuUmVxdWVzdFBheW'
    'VyUgxyZXF1ZXN0cGF5ZXISNQoUc3NlY3VzdG9tZXJhbGdvcml0aG0Y0MmBKyABKAlSFHNzZWN1'
    'c3RvbWVyYWxnb3JpdGhtEikKDnNzZWN1c3RvbWVya2V5GJr+9DsgASgJUg5zc2VjdXN0b21lcm'
    'tleRIuChFzc2VjdXN0b21lcmtleW1kNRjo0RcgASgJUhFzc2VjdXN0b21lcmtleW1kNRI7Chdz'
    'c2VrbXNlbmNyeXB0aW9uY29udGV4dBi6kIhHIAEoCVIXc3Nla21zZW5jcnlwdGlvbmNvbnRleH'
    'QSJAoLc3Nla21za2V5aWQY2IjU1AEgASgJUgtzc2VrbXNrZXlpZBJPChRzZXJ2ZXJzaWRlZW5j'
    'cnlwdGlvbhixjp8EIAEoDjIYLnMzLlNlcnZlclNpZGVFbmNyeXB0aW9uUhRzZXJ2ZXJzaWRlZW'
    '5jcnlwdGlvbhI4CgxzdG9yYWdlY2xhc3MYx4jEuwEgASgOMhAuczMuU3RvcmFnZUNsYXNzUgxz'
    'dG9yYWdlY2xhc3MSGwoHdGFnZ2luZxj95vgPIAEoCVIHdGFnZ2luZxJEChB0YWdnaW5nZGlyZW'
    'N0aXZlGMqD2+0BIAEoDjIULnMzLlRhZ2dpbmdEaXJlY3RpdmVSEHRhZ2dpbmdkaXJlY3RpdmUS'
    'OwoXd2Vic2l0ZXJlZGlyZWN0bG9jYXRpb24YtoahIiABKAlSF3dlYnNpdGVyZWRpcmVjdGxvY2'
    'F0aW9uGjsKDU1ldGFkYXRhRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlS'
    'BXZhbHVlOgI4AQ==');

@$core.Deprecated('Use copyObjectResultDescriptor instead')
const CopyObjectResult$json = {
  '1': 'CopyObjectResult',
  '2': [
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
  ],
};

/// Descriptor for `CopyObjectResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List copyObjectResultDescriptor = $convert.base64Decode(
    'ChBDb3B5T2JqZWN0UmVzdWx0EicKDWNoZWNrc3VtY3JjMzIYwqPNMyABKAlSDWNoZWNrc3VtY3'
    'JjMzISKQoOY2hlY2tzdW1jcmMzMmMYp5+lTCABKAlSDmNoZWNrc3VtY3JjMzJjEi8KEWNoZWNr'
    'c3VtY3JjNjRudm1lGM28izggASgJUhFjaGVja3N1bWNyYzY0bnZtZRImCgxjaGVja3N1bXNoYT'
    'EYxOzgigEgASgJUgxjaGVja3N1bXNoYTESKQoOY2hlY2tzdW1zaGEyNTYYvvncRCABKAlSDmNo'
    'ZWNrc3Vtc2hhMjU2EjcKDGNoZWNrc3VtdHlwZRjDvtkuIAEoDjIQLnMzLkNoZWNrc3VtVHlwZV'
    'IMY2hlY2tzdW10eXBlEhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEiYKDGxhc3Rtb2RpZmllZBin'
    'nPzOASABKAlSDGxhc3Rtb2RpZmllZA==');

@$core.Deprecated('Use copyPartResultDescriptor instead')
const CopyPartResult$json = {
  '1': 'CopyPartResult',
  '2': [
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
  ],
};

/// Descriptor for `CopyPartResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List copyPartResultDescriptor = $convert.base64Decode(
    'Cg5Db3B5UGFydFJlc3VsdBInCg1jaGVja3N1bWNyYzMyGMKjzTMgASgJUg1jaGVja3N1bWNyYz'
    'MyEikKDmNoZWNrc3VtY3JjMzJjGKefpUwgASgJUg5jaGVja3N1bWNyYzMyYxIvChFjaGVja3N1'
    'bWNyYzY0bnZtZRjNvIs4IAEoCVIRY2hlY2tzdW1jcmM2NG52bWUSJgoMY2hlY2tzdW1zaGExGM'
    'Ts4IoBIAEoCVIMY2hlY2tzdW1zaGExEikKDmNoZWNrc3Vtc2hhMjU2GL753EQgASgJUg5jaGVj'
    'a3N1bXNoYTI1NhIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxImCgxsYXN0bW9kaWZpZWQYp5z8zg'
    'EgASgJUgxsYXN0bW9kaWZpZWQ=');

@$core.Deprecated('Use createBucketConfigurationDescriptor instead')
const CreateBucketConfiguration$json = {
  '1': 'CreateBucketConfiguration',
  '2': [
    {
      '1': 'bucket',
      '3': 55457112,
      '4': 1,
      '5': 11,
      '6': '.s3.BucketInfo',
      '10': 'bucket'
    },
    {
      '1': 'location',
      '3': 465604039,
      '4': 1,
      '5': 11,
      '6': '.s3.LocationInfo',
      '10': 'location'
    },
    {
      '1': 'locationconstraint',
      '3': 245375838,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketLocationConstraint',
      '10': 'locationconstraint'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateBucketConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBucketConfigurationDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVCdWNrZXRDb25maWd1cmF0aW9uEikKBmJ1Y2tldBjY6rgaIAEoCzIOLnMzLkJ1Y2'
    'tldEluZm9SBmJ1Y2tldBIwCghsb2NhdGlvbhjHm4LeASABKAsyEC5zMy5Mb2NhdGlvbkluZm9S'
    'CGxvY2F0aW9uEk8KEmxvY2F0aW9uY29uc3RyYWludBjexoB1IAEoDjIcLnMzLkJ1Y2tldExvY2'
    'F0aW9uQ29uc3RyYWludFISbG9jYXRpb25jb25zdHJhaW50Eh8KBHRhZ3MYwcH2tQEgAygLMgcu'
    'czMuVGFnUgR0YWdz');

@$core.Deprecated(
    'Use createBucketMetadataConfigurationRequestDescriptor instead')
const CreateBucketMetadataConfigurationRequest$json = {
  '1': 'CreateBucketMetadataConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'metadataconfiguration',
      '3': 361784403,
      '4': 1,
      '5': 11,
      '6': '.s3.MetadataConfiguration',
      '8': {},
      '10': 'metadataconfiguration'
    },
  ],
};

/// Descriptor for `CreateBucketMetadataConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBucketMetadataConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CihDcmVhdGVCdWNrZXRNZXRhZGF0YUNvbmZpZ3VyYXRpb25SZXF1ZXN0EhkKBmJ1Y2tldBjY6r'
        'gaIAEoCVIGYnVja2V0EkYKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMuQ2hlY2tz'
        'dW1BbGdvcml0aG1SEWNoZWNrc3VtYWxnb3JpdGhtEiIKCmNvbnRlbnRtZDUYi+rq8QEgASgJUg'
        'pjb250ZW50bWQ1EjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVj'
        'a2V0b3duZXISWQoVbWV0YWRhdGFjb25maWd1cmF0aW9uGNPIwawBIAEoCzIZLnMzLk1ldGFkYX'
        'RhQ29uZmlndXJhdGlvbkIEiLUYAVIVbWV0YWRhdGFjb25maWd1cmF0aW9u');

@$core.Deprecated(
    'Use createBucketMetadataTableConfigurationRequestDescriptor instead')
const CreateBucketMetadataTableConfigurationRequest$json = {
  '1': 'CreateBucketMetadataTableConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'metadatatableconfiguration',
      '3': 289222191,
      '4': 1,
      '5': 11,
      '6': '.s3.MetadataTableConfiguration',
      '8': {},
      '10': 'metadatatableconfiguration'
    },
  ],
};

/// Descriptor for `CreateBucketMetadataTableConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createBucketMetadataTableConfigurationRequestDescriptor =
    $convert.base64Decode(
        'Ci1DcmVhdGVCdWNrZXRNZXRhZGF0YVRhYmxlQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2'
        'V0GNjquBogASgJUgZidWNrZXQSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5D'
        'aGVja3N1bUFsZ29yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxAS'
        'ABKAlSCmNvbnRlbnRtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0'
        'ZWRidWNrZXRvd25lchJoChptZXRhZGF0YXRhYmxlY29uZmlndXJhdGlvbhiv3PSJASABKAsyHi'
        '5zMy5NZXRhZGF0YVRhYmxlQ29uZmlndXJhdGlvbkIEiLUYAVIabWV0YWRhdGF0YWJsZWNvbmZp'
        'Z3VyYXRpb24=');

@$core.Deprecated('Use createBucketOutputDescriptor instead')
const CreateBucketOutput$json = {
  '1': 'CreateBucketOutput',
  '2': [
    {'1': 'bucketarn', '3': 255683899, '4': 1, '5': 9, '10': 'bucketarn'},
    {'1': 'location', '3': 465604039, '4': 1, '5': 9, '10': 'location'},
  ],
};

/// Descriptor for `CreateBucketOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBucketOutputDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVCdWNrZXRPdXRwdXQSHwoJYnVja2V0YXJuGLva9XkgASgJUglidWNrZXRhcm4SHg'
    'oIbG9jYXRpb24Yx5uC3gEgASgJUghsb2NhdGlvbg==');

@$core.Deprecated('Use createBucketRequestDescriptor instead')
const CreateBucketRequest$json = {
  '1': 'CreateBucketRequest',
  '2': [
    {
      '1': 'acl',
      '3': 394696836,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketCannedACL',
      '10': 'acl'
    },
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'createbucketconfiguration',
      '3': 27550528,
      '4': 1,
      '5': 11,
      '6': '.s3.CreateBucketConfiguration',
      '8': {},
      '10': 'createbucketconfiguration'
    },
    {
      '1': 'grantfullcontrol',
      '3': 102486874,
      '4': 1,
      '5': 9,
      '10': 'grantfullcontrol'
    },
    {'1': 'grantread', '3': 4518126, '4': 1, '5': 9, '10': 'grantread'},
    {'1': 'grantreadacp', '3': 208410960, '4': 1, '5': 9, '10': 'grantreadacp'},
    {'1': 'grantwrite', '3': 286850821, '4': 1, '5': 9, '10': 'grantwrite'},
    {
      '1': 'grantwriteacp',
      '3': 157594957,
      '4': 1,
      '5': 9,
      '10': 'grantwriteacp'
    },
    {
      '1': 'objectlockenabledforbucket',
      '3': 248510328,
      '4': 1,
      '5': 8,
      '10': 'objectlockenabledforbucket'
    },
    {
      '1': 'objectownership',
      '3': 448301184,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectOwnership',
      '10': 'objectownership'
    },
  ],
};

/// Descriptor for `CreateBucketRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBucketRequestDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVCdWNrZXRSZXF1ZXN0EikKA2FjbBiEsZq8ASABKA4yEy5zMy5CdWNrZXRDYW5uZW'
    'RBQ0xSA2FjbBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBJkChljcmVhdGVidWNrZXRjb25m'
    'aWd1cmF0aW9uGMDGkQ0gASgLMh0uczMuQ3JlYXRlQnVja2V0Q29uZmlndXJhdGlvbkIEiLUYAV'
    'IZY3JlYXRlYnVja2V0Y29uZmlndXJhdGlvbhItChBncmFudGZ1bGxjb250cm9sGNqm7zAgASgJ'
    'UhBncmFudGZ1bGxjb250cm9sEh8KCWdyYW50cmVhZBju4ZMCIAEoCVIJZ3JhbnRyZWFkEiUKDG'
    'dyYW50cmVhZGFjcBjQsrBjIAEoCVIMZ3JhbnRyZWFkYWNwEiIKCmdyYW50d3JpdGUYhf7jiAEg'
    'ASgJUgpncmFudHdyaXRlEicKDWdyYW50d3JpdGVhY3AYzeqSSyABKAlSDWdyYW50d3JpdGVhY3'
    'ASQQoab2JqZWN0bG9ja2VuYWJsZWRmb3JidWNrZXQY+O6/diABKAhSGm9iamVjdGxvY2tlbmFi'
    'bGVkZm9yYnVja2V0EkEKD29iamVjdG93bmVyc2hpcBiAkeLVASABKA4yEy5zMy5PYmplY3RPd2'
    '5lcnNoaXBSD29iamVjdG93bmVyc2hpcA==');

@$core.Deprecated('Use createMultipartUploadOutputDescriptor instead')
const CreateMultipartUploadOutput$json = {
  '1': 'CreateMultipartUploadOutput',
  '2': [
    {'1': 'abortdate', '3': 232475318, '4': 1, '5': 9, '10': 'abortdate'},
    {'1': 'abortruleid', '3': 232462739, '4': 1, '5': 9, '10': 'abortruleid'},
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {
      '1': 'ssekmsencryptioncontext',
      '3': 149030970,
      '4': 1,
      '5': 9,
      '10': 'ssekmsencryptioncontext'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {'1': 'uploadid', '3': 449040722, '4': 1, '5': 9, '10': 'uploadid'},
  ],
};

/// Descriptor for `CreateMultipartUploadOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createMultipartUploadOutputDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVNdWx0aXBhcnRVcGxvYWRPdXRwdXQSHwoJYWJvcnRkYXRlGLaV7W4gASgJUglhYm'
    '9ydGRhdGUSIwoLYWJvcnRydWxlaWQYk7PsbiABKAlSC2Fib3J0cnVsZWlkEhkKBmJ1Y2tldBjY'
    '6rgaIAEoCVIGYnVja2V0Ei4KEGJ1Y2tldGtleWVuYWJsZWQY7KKMzwEgASgIUhBidWNrZXRrZX'
    'llbmFibGVkEkYKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMuQ2hlY2tzdW1BbGdv'
    'cml0aG1SEWNoZWNrc3VtYWxnb3JpdGhtEjcKDGNoZWNrc3VtdHlwZRjDvtkuIAEoDjIQLnMzLk'
    'NoZWNrc3VtVHlwZVIMY2hlY2tzdW10eXBlEhMKA2tleRiNkutoIAEoCVIDa2V5Ej4KDnJlcXVl'
    'c3RjaGFyZ2VkGJPQq7kBIAEoDjISLnMzLlJlcXVlc3RDaGFyZ2VkUg5yZXF1ZXN0Y2hhcmdlZB'
    'I1ChRzc2VjdXN0b21lcmFsZ29yaXRobRjQyYErIAEoCVIUc3NlY3VzdG9tZXJhbGdvcml0aG0S'
    'LgoRc3NlY3VzdG9tZXJrZXltZDUY6NEXIAEoCVIRc3NlY3VzdG9tZXJrZXltZDUSOwoXc3Nla2'
    '1zZW5jcnlwdGlvbmNvbnRleHQYupCIRyABKAlSF3NzZWttc2VuY3J5cHRpb25jb250ZXh0EiQK'
    'C3NzZWttc2tleWlkGNiI1NQBIAEoCVILc3Nla21za2V5aWQSTwoUc2VydmVyc2lkZWVuY3J5cH'
    'Rpb24YsY6fBCABKA4yGC5zMy5TZXJ2ZXJTaWRlRW5jcnlwdGlvblIUc2VydmVyc2lkZWVuY3J5'
    'cHRpb24SHgoIdXBsb2FkaWQY0qKP1gEgASgJUgh1cGxvYWRpZA==');

@$core.Deprecated('Use createMultipartUploadRequestDescriptor instead')
const CreateMultipartUploadRequest$json = {
  '1': 'CreateMultipartUploadRequest',
  '2': [
    {
      '1': 'acl',
      '3': 394696836,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectCannedACL',
      '10': 'acl'
    },
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {'1': 'cachecontrol', '3': 288966655, '4': 1, '5': 9, '10': 'cachecontrol'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {
      '1': 'contentdisposition',
      '3': 120040130,
      '4': 1,
      '5': 9,
      '10': 'contentdisposition'
    },
    {
      '1': 'contentencoding',
      '3': 317106228,
      '4': 1,
      '5': 9,
      '10': 'contentencoding'
    },
    {
      '1': 'contentlanguage',
      '3': 108485649,
      '4': 1,
      '5': 9,
      '10': 'contentlanguage'
    },
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'expires', '3': 128582948, '4': 1, '5': 9, '10': 'expires'},
    {
      '1': 'grantfullcontrol',
      '3': 102486874,
      '4': 1,
      '5': 9,
      '10': 'grantfullcontrol'
    },
    {'1': 'grantread', '3': 4518126, '4': 1, '5': 9, '10': 'grantread'},
    {'1': 'grantreadacp', '3': 208410960, '4': 1, '5': 9, '10': 'grantreadacp'},
    {
      '1': 'grantwriteacp',
      '3': 157594957,
      '4': 1,
      '5': 9,
      '10': 'grantwriteacp'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 3,
      '5': 11,
      '6': '.s3.CreateMultipartUploadRequest.MetadataEntry',
      '10': 'metadata'
    },
    {
      '1': 'objectlocklegalholdstatus',
      '3': 536561974,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockLegalHoldStatus',
      '10': 'objectlocklegalholdstatus'
    },
    {
      '1': 'objectlockmode',
      '3': 189255203,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockMode',
      '10': 'objectlockmode'
    },
    {
      '1': 'objectlockretainuntildate',
      '3': 264584249,
      '4': 1,
      '5': 9,
      '10': 'objectlockretainuntildate'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {
      '1': 'ssekmsencryptioncontext',
      '3': 149030970,
      '4': 1,
      '5': 9,
      '10': 'ssekmsencryptioncontext'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'tagging', '3': 33436541, '4': 1, '5': 9, '10': 'tagging'},
    {
      '1': 'websiteredirectlocation',
      '3': 71844662,
      '4': 1,
      '5': 9,
      '10': 'websiteredirectlocation'
    },
  ],
  '3': [CreateMultipartUploadRequest_MetadataEntry$json],
};

@$core.Deprecated('Use createMultipartUploadRequestDescriptor instead')
const CreateMultipartUploadRequest_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateMultipartUploadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createMultipartUploadRequestDescriptor = $convert.base64Decode(
    'ChxDcmVhdGVNdWx0aXBhcnRVcGxvYWRSZXF1ZXN0EikKA2FjbBiEsZq8ASABKA4yEy5zMy5PYm'
    'plY3RDYW5uZWRBQ0xSA2FjbBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIuChBidWNrZXRr'
    'ZXllbmFibGVkGOyijM8BIAEoCFIQYnVja2V0a2V5ZW5hYmxlZBImCgxjYWNoZWNvbnRyb2wY/4'
    '/liQEgASgJUgxjYWNoZWNvbnRyb2wSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5z'
    'My5DaGVja3N1bUFsZ29yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SNwoMY2hlY2tzdW10eXBlGM'
    'O+2S4gASgOMhAuczMuQ2hlY2tzdW1UeXBlUgxjaGVja3N1bXR5cGUSMQoSY29udGVudGRpc3Bv'
    'c2l0aW9uGMLVnjkgASgJUhJjb250ZW50ZGlzcG9zaXRpb24SLAoPY29udGVudGVuY29kaW5nGL'
    'TQmpcBIAEoCVIPY29udGVudGVuY29kaW5nEisKD2NvbnRlbnRsYW5ndWFnZRiRuN0zIAEoCVIP'
    'Y29udGVudGxhbmd1YWdlEiQKC2NvbnRlbnR0eXBlGJPV6J4BIAEoCVILY29udGVudHR5cGUSMw'
    'oTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchIbCgdl'
    'eHBpcmVzGKSKqD0gASgJUgdleHBpcmVzEi0KEGdyYW50ZnVsbGNvbnRyb2wY2qbvMCABKAlSEG'
    'dyYW50ZnVsbGNvbnRyb2wSHwoJZ3JhbnRyZWFkGO7hkwIgASgJUglncmFudHJlYWQSJQoMZ3Jh'
    'bnRyZWFkYWNwGNCysGMgASgJUgxncmFudHJlYWRhY3ASJwoNZ3JhbnR3cml0ZWFjcBjN6pJLIA'
    'EoCVINZ3JhbnR3cml0ZWFjcBITCgNrZXkYjZLraCABKAlSA2tleRJOCghtZXRhZGF0YRjh4o/g'
    'ASADKAsyLi5zMy5DcmVhdGVNdWx0aXBhcnRVcGxvYWRSZXF1ZXN0Lk1ldGFkYXRhRW50cnlSCG'
    '1ldGFkYXRhEl8KGW9iamVjdGxvY2tsZWdhbGhvbGRzdGF0dXMYtpLt/wEgASgOMh0uczMuT2Jq'
    'ZWN0TG9ja0xlZ2FsSG9sZFN0YXR1c1IZb2JqZWN0bG9ja2xlZ2FsaG9sZHN0YXR1cxI9Cg5vYm'
    'plY3Rsb2NrbW9kZRijnJ9aIAEoDjISLnMzLk9iamVjdExvY2tNb2RlUg5vYmplY3Rsb2NrbW9k'
    'ZRI/ChlvYmplY3Rsb2NrcmV0YWludW50aWxkYXRlGLn4lH4gASgJUhlvYmplY3Rsb2NrcmV0YW'
    'ludW50aWxkYXRlEjgKDHJlcXVlc3RwYXllchik5uH1ASABKA4yEC5zMy5SZXF1ZXN0UGF5ZXJS'
    'DHJlcXVlc3RwYXllchI1ChRzc2VjdXN0b21lcmFsZ29yaXRobRjQyYErIAEoCVIUc3NlY3VzdG'
    '9tZXJhbGdvcml0aG0SKQoOc3NlY3VzdG9tZXJrZXkYmv70OyABKAlSDnNzZWN1c3RvbWVya2V5'
    'Ei4KEXNzZWN1c3RvbWVya2V5bWQ1GOjRFyABKAlSEXNzZWN1c3RvbWVya2V5bWQ1EjsKF3NzZW'
    'ttc2VuY3J5cHRpb25jb250ZXh0GLqQiEcgASgJUhdzc2VrbXNlbmNyeXB0aW9uY29udGV4dBIk'
    'Cgtzc2VrbXNrZXlpZBjYiNTUASABKAlSC3NzZWttc2tleWlkEk8KFHNlcnZlcnNpZGVlbmNyeX'
    'B0aW9uGLGOnwQgASgOMhguczMuU2VydmVyU2lkZUVuY3J5cHRpb25SFHNlcnZlcnNpZGVlbmNy'
    'eXB0aW9uEjgKDHN0b3JhZ2VjbGFzcxjHiMS7ASABKA4yEC5zMy5TdG9yYWdlQ2xhc3NSDHN0b3'
    'JhZ2VjbGFzcxIbCgd0YWdnaW5nGP3m+A8gASgJUgd0YWdnaW5nEjsKF3dlYnNpdGVyZWRpcmVj'
    'dGxvY2F0aW9uGLaGoSIgASgJUhd3ZWJzaXRlcmVkaXJlY3Rsb2NhdGlvbho7Cg1NZXRhZGF0YU'
    'VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createSessionOutputDescriptor instead')
const CreateSessionOutput$json = {
  '1': 'CreateSessionOutput',
  '2': [
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.s3.SessionCredentials',
      '10': 'credentials'
    },
    {
      '1': 'ssekmsencryptioncontext',
      '3': 149030970,
      '4': 1,
      '5': 9,
      '10': 'ssekmsencryptioncontext'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
  ],
};

/// Descriptor for `CreateSessionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSessionOutputDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVTZXNzaW9uT3V0cHV0Ei4KEGJ1Y2tldGtleWVuYWJsZWQY7KKMzwEgASgIUhBidW'
    'NrZXRrZXllbmFibGVkEjwKC2NyZWRlbnRpYWxzGPKajrYBIAEoCzIWLnMzLlNlc3Npb25DcmVk'
    'ZW50aWFsc1ILY3JlZGVudGlhbHMSOwoXc3Nla21zZW5jcnlwdGlvbmNvbnRleHQYupCIRyABKA'
    'lSF3NzZWttc2VuY3J5cHRpb25jb250ZXh0EiQKC3NzZWttc2tleWlkGNiI1NQBIAEoCVILc3Nl'
    'a21za2V5aWQSTwoUc2VydmVyc2lkZWVuY3J5cHRpb24YsY6fBCABKA4yGC5zMy5TZXJ2ZXJTaW'
    'RlRW5jcnlwdGlvblIUc2VydmVyc2lkZWVuY3J5cHRpb24=');

@$core.Deprecated('Use createSessionRequestDescriptor instead')
const CreateSessionRequest$json = {
  '1': 'CreateSessionRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {
      '1': 'ssekmsencryptioncontext',
      '3': 149030970,
      '4': 1,
      '5': 9,
      '10': 'ssekmsencryptioncontext'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {
      '1': 'sessionmode',
      '3': 43104977,
      '4': 1,
      '5': 14,
      '6': '.s3.SessionMode',
      '10': 'sessionmode'
    },
  ],
};

/// Descriptor for `CreateSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSessionRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVTZXNzaW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIuChBidW'
    'NrZXRrZXllbmFibGVkGOyijM8BIAEoCFIQYnVja2V0a2V5ZW5hYmxlZBI7Chdzc2VrbXNlbmNy'
    'eXB0aW9uY29udGV4dBi6kIhHIAEoCVIXc3Nla21zZW5jcnlwdGlvbmNvbnRleHQSJAoLc3Nla2'
    '1za2V5aWQY2IjU1AEgASgJUgtzc2VrbXNrZXlpZBJPChRzZXJ2ZXJzaWRlZW5jcnlwdGlvbhix'
    'jp8EIAEoDjIYLnMzLlNlcnZlclNpZGVFbmNyeXB0aW9uUhRzZXJ2ZXJzaWRlZW5jcnlwdGlvbh'
    'I0CgtzZXNzaW9ubW9kZRjR9cYUIAEoDjIPLnMzLlNlc3Npb25Nb2RlUgtzZXNzaW9ubW9kZQ==');

@$core.Deprecated('Use defaultRetentionDescriptor instead')
const DefaultRetention$json = {
  '1': 'DefaultRetention',
  '2': [
    {'1': 'days', '3': 494075051, '4': 1, '5': 5, '10': 'days'},
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockRetentionMode',
      '10': 'mode'
    },
    {'1': 'years', '3': 325150814, '4': 1, '5': 5, '10': 'years'},
  ],
};

/// Descriptor for `DefaultRetention`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List defaultRetentionDescriptor = $convert.base64Decode(
    'ChBEZWZhdWx0UmV0ZW50aW9uEhYKBGRheXMYq/nL6wEgASgFUgRkYXlzEjMKBG1vZGUYs+65mg'
    'EgASgOMhsuczMuT2JqZWN0TG9ja1JldGVudGlvbk1vZGVSBG1vZGUSGAoFeWVhcnMY3tCFmwEg'
    'ASgFUgV5ZWFycw==');

@$core.Deprecated('Use deleteDescriptor instead')
const Delete$json = {
  '1': 'Delete',
  '2': [
    {
      '1': 'objects',
      '3': 136869388,
      '4': 3,
      '5': 11,
      '6': '.s3.ObjectIdentifier',
      '10': 'objects'
    },
    {'1': 'quiet', '3': 267428792, '4': 1, '5': 8, '10': 'quiet'},
  ],
};

/// Descriptor for `Delete`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDescriptor = $convert.base64Decode(
    'CgZEZWxldGUSMQoHb2JqZWN0cxiM7KFBIAMoCzIULnMzLk9iamVjdElkZW50aWZpZXJSB29iam'
    'VjdHMSFwoFcXVpZXQYuMfCfyABKAhSBXF1aWV0');

@$core.Deprecated(
    'Use deleteBucketAnalyticsConfigurationRequestDescriptor instead')
const DeleteBucketAnalyticsConfigurationRequest$json = {
  '1': 'DeleteBucketAnalyticsConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteBucketAnalyticsConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteBucketAnalyticsConfigurationRequestDescriptor = $convert.base64Decode(
        'CilEZWxldGVCdWNrZXRBbmFseXRpY3NDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2O'
        'q4GiABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3Rl'
        'ZGJ1Y2tldG93bmVyEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use deleteBucketCorsRequestDescriptor instead')
const DeleteBucketCorsRequest$json = {
  '1': 'DeleteBucketCorsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketCorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketCorsRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVCdWNrZXRDb3JzUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzCh'
        'NleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use deleteBucketEncryptionRequestDescriptor instead')
const DeleteBucketEncryptionRequest$json = {
  '1': 'DeleteBucketEncryptionRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketEncryptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketEncryptionRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVCdWNrZXRFbmNyeXB0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2'
        'tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated(
    'Use deleteBucketIntelligentTieringConfigurationRequestDescriptor instead')
const DeleteBucketIntelligentTieringConfigurationRequest$json = {
  '1': 'DeleteBucketIntelligentTieringConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteBucketIntelligentTieringConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteBucketIntelligentTieringConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CjJEZWxldGVCdWNrZXRJbnRlbGxpZ2VudFRpZXJpbmdDb25maWd1cmF0aW9uUmVxdWVzdBIZCg'
        'ZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJ'
        'UhNleHBlY3RlZGJ1Y2tldG93bmVyEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated(
    'Use deleteBucketInventoryConfigurationRequestDescriptor instead')
const DeleteBucketInventoryConfigurationRequest$json = {
  '1': 'DeleteBucketInventoryConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteBucketInventoryConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteBucketInventoryConfigurationRequestDescriptor = $convert.base64Decode(
        'CilEZWxldGVCdWNrZXRJbnZlbnRvcnlDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2O'
        'q4GiABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3Rl'
        'ZGJ1Y2tldG93bmVyEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use deleteBucketLifecycleRequestDescriptor instead')
const DeleteBucketLifecycleRequest$json = {
  '1': 'DeleteBucketLifecycleRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketLifecycleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketLifecycleRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVCdWNrZXRMaWZlY3ljbGVSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2'
        'V0EjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXI=');

@$core.Deprecated(
    'Use deleteBucketMetadataConfigurationRequestDescriptor instead')
const DeleteBucketMetadataConfigurationRequest$json = {
  '1': 'DeleteBucketMetadataConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketMetadataConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketMetadataConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CihEZWxldGVCdWNrZXRNZXRhZGF0YUNvbmZpZ3VyYXRpb25SZXF1ZXN0EhkKBmJ1Y2tldBjY6r'
        'gaIAEoCVIGYnVja2V0EjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVk'
        'YnVja2V0b3duZXI=');

@$core.Deprecated(
    'Use deleteBucketMetadataTableConfigurationRequestDescriptor instead')
const DeleteBucketMetadataTableConfigurationRequest$json = {
  '1': 'DeleteBucketMetadataTableConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketMetadataTableConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteBucketMetadataTableConfigurationRequestDescriptor =
    $convert.base64Decode(
        'Ci1EZWxldGVCdWNrZXRNZXRhZGF0YVRhYmxlQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2'
        'V0GNjquBogASgJUgZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhw'
        'ZWN0ZWRidWNrZXRvd25lcg==');

@$core
    .Deprecated('Use deleteBucketMetricsConfigurationRequestDescriptor instead')
const DeleteBucketMetricsConfigurationRequest$json = {
  '1': 'DeleteBucketMetricsConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteBucketMetricsConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketMetricsConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CidEZWxldGVCdWNrZXRNZXRyaWNzQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquB'
        'ogASgJUgZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRi'
        'dWNrZXRvd25lchISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use deleteBucketOwnershipControlsRequestDescriptor instead')
const DeleteBucketOwnershipControlsRequest$json = {
  '1': 'DeleteBucketOwnershipControlsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketOwnershipControlsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketOwnershipControlsRequestDescriptor =
    $convert.base64Decode(
        'CiREZWxldGVCdWNrZXRPd25lcnNoaXBDb250cm9sc1JlcXVlc3QSGQoGYnVja2V0GNjquBogAS'
        'gJUgZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNr'
        'ZXRvd25lcg==');

@$core.Deprecated('Use deleteBucketPolicyRequestDescriptor instead')
const DeleteBucketPolicyRequest$json = {
  '1': 'DeleteBucketPolicyRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVCdWNrZXRQb2xpY3lSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0Ej'
        'MKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXI=');

@$core.Deprecated('Use deleteBucketReplicationRequestDescriptor instead')
const DeleteBucketReplicationRequest$json = {
  '1': 'DeleteBucketReplicationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketReplicationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketReplicationRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVCdWNrZXRSZXBsaWNhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidW'
        'NrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25l'
        'cg==');

@$core.Deprecated('Use deleteBucketRequestDescriptor instead')
const DeleteBucketRequest$json = {
  '1': 'DeleteBucketRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketRequestDescriptor = $convert.base64Decode(
    'ChNEZWxldGVCdWNrZXRSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EjMKE2V4cG'
    'VjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXI=');

@$core.Deprecated('Use deleteBucketTaggingRequestDescriptor instead')
const DeleteBucketTaggingRequest$json = {
  '1': 'DeleteBucketTaggingRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketTaggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketTaggingRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVCdWNrZXRUYWdnaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
        'IzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use deleteBucketWebsiteRequestDescriptor instead')
const DeleteBucketWebsiteRequest$json = {
  '1': 'DeleteBucketWebsiteRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeleteBucketWebsiteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteBucketWebsiteRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVCdWNrZXRXZWJzaXRlUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
        'IzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use deleteMarkerEntryDescriptor instead')
const DeleteMarkerEntry$json = {
  '1': 'DeleteMarkerEntry',
  '2': [
    {'1': 'islatest', '3': 80355831, '4': 1, '5': 8, '10': 'islatest'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `DeleteMarkerEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMarkerEntryDescriptor = $convert.base64Decode(
    'ChFEZWxldGVNYXJrZXJFbnRyeRIdCghpc2xhdGVzdBj3w6gmIAEoCFIIaXNsYXRlc3QSEwoDa2'
    'V5GI2S62ggASgJUgNrZXkSJgoMbGFzdG1vZGlmaWVkGKec/M4BIAEoCVIMbGFzdG1vZGlmaWVk'
    'EiMKBW93bmVyGPX8itkBIAEoCzIJLnMzLk93bmVyUgVvd25lchIgCgl2ZXJzaW9uaWQYm+GZoQ'
    'EgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use deleteMarkerReplicationDescriptor instead')
const DeleteMarkerReplication$json = {
  '1': 'DeleteMarkerReplication',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.DeleteMarkerReplicationStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `DeleteMarkerReplication`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMarkerReplicationDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVNYXJrZXJSZXBsaWNhdGlvbhI8CgZzdGF0dXMYkOT7AiABKA4yIS5zMy5EZWxldG'
        'VNYXJrZXJSZXBsaWNhdGlvblN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use deleteObjectOutputDescriptor instead')
const DeleteObjectOutput$json = {
  '1': 'DeleteObjectOutput',
  '2': [
    {'1': 'deletemarker', '3': 5472257, '4': 1, '5': 8, '10': 'deletemarker'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `DeleteObjectOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteObjectOutputDescriptor = $convert.base64Decode(
    'ChJEZWxldGVPYmplY3RPdXRwdXQSJQoMZGVsZXRlbWFya2VyGIGAzgIgASgIUgxkZWxldGVtYX'
    'JrZXISPgoOcmVxdWVzdGNoYXJnZWQYk9CruQEgASgOMhIuczMuUmVxdWVzdENoYXJnZWRSDnJl'
    'cXVlc3RjaGFyZ2VkEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated('Use deleteObjectRequestDescriptor instead')
const DeleteObjectRequest$json = {
  '1': 'DeleteObjectRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bypassgovernanceretention',
      '3': 129938942,
      '4': 1,
      '5': 8,
      '10': 'bypassgovernanceretention'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'ifmatchlastmodifiedtime',
      '3': 94688686,
      '4': 1,
      '5': 9,
      '10': 'ifmatchlastmodifiedtime'
    },
    {'1': 'ifmatchsize', '3': 398307783, '4': 1, '5': 3, '10': 'ifmatchsize'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'mfa', '3': 325407842, '4': 1, '5': 9, '10': 'mfa'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `DeleteObjectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteObjectRequestDescriptor = $convert.base64Decode(
    'ChNEZWxldGVPYmplY3RSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0Ej8KGWJ5cG'
    'Fzc2dvdmVybmFuY2VyZXRlbnRpb24Y/uv6PSABKAhSGWJ5cGFzc2dvdmVybmFuY2VyZXRlbnRp'
    'b24SMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lch'
    'IbCgdpZm1hdGNoGNCWtywgASgJUgdpZm1hdGNoEjsKF2lmbWF0Y2hsYXN0bW9kaWZpZWR0aW1l'
    'GK6rky0gASgJUhdpZm1hdGNobGFzdG1vZGlmaWVkdGltZRIkCgtpZm1hdGNoc2l6ZRjH4/a9AS'
    'ABKANSC2lmbWF0Y2hzaXplEhMKA2tleRiNkutoIAEoCVIDa2V5EhQKA21mYRjiqJWbASABKAlS'
    'A21mYRI4CgxyZXF1ZXN0cGF5ZXIYpObh9QEgASgOMhAuczMuUmVxdWVzdFBheWVyUgxyZXF1ZX'
    'N0cGF5ZXISIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlk');

@$core.Deprecated('Use deleteObjectTaggingOutputDescriptor instead')
const DeleteObjectTaggingOutput$json = {
  '1': 'DeleteObjectTaggingOutput',
  '2': [
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `DeleteObjectTaggingOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteObjectTaggingOutputDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVPYmplY3RUYWdnaW5nT3V0cHV0EiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcn'
        'Npb25pZA==');

@$core.Deprecated('Use deleteObjectTaggingRequestDescriptor instead')
const DeleteObjectTaggingRequest$json = {
  '1': 'DeleteObjectTaggingRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `DeleteObjectTaggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteObjectTaggingRequestDescriptor = $convert.base64Decode(
    'ChpEZWxldGVPYmplY3RUYWdnaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
    'IzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhMK'
    'A2tleRiNkutoIAEoCVIDa2V5EiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated('Use deleteObjectsOutputDescriptor instead')
const DeleteObjectsOutput$json = {
  '1': 'DeleteObjectsOutput',
  '2': [
    {
      '1': 'deleted',
      '3': 304444427,
      '4': 3,
      '5': 11,
      '6': '.s3.DeletedObject',
      '10': 'deleted'
    },
    {
      '1': 'errors',
      '3': 166551719,
      '4': 3,
      '5': 11,
      '6': '.s3.Error',
      '10': 'errors'
    },
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `DeleteObjectsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteObjectsOutputDescriptor = $convert.base64Decode(
    'ChNEZWxldGVPYmplY3RzT3V0cHV0Ei8KB2RlbGV0ZWQYi+iVkQEgAygLMhEuczMuRGVsZXRlZE'
    '9iamVjdFIHZGVsZXRlZBIkCgZlcnJvcnMYp8G1TyADKAsyCS5zMy5FcnJvclIGZXJyb3JzEj4K'
    'DnJlcXVlc3RjaGFyZ2VkGJPQq7kBIAEoDjISLnMzLlJlcXVlc3RDaGFyZ2VkUg5yZXF1ZXN0Y2'
    'hhcmdlZA==');

@$core.Deprecated('Use deleteObjectsRequestDescriptor instead')
const DeleteObjectsRequest$json = {
  '1': 'DeleteObjectsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bypassgovernanceretention',
      '3': 129938942,
      '4': 1,
      '5': 8,
      '10': 'bypassgovernanceretention'
    },
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'delete',
      '3': 395831915,
      '4': 1,
      '5': 11,
      '6': '.s3.Delete',
      '8': {},
      '10': 'delete'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'mfa', '3': 325407842, '4': 1, '5': 9, '10': 'mfa'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
  ],
};

/// Descriptor for `DeleteObjectsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteObjectsRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVPYmplY3RzUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBI/ChlieX'
    'Bhc3Nnb3Zlcm5hbmNlcmV0ZW50aW9uGP7r+j0gASgIUhlieXBhc3Nnb3Zlcm5hbmNlcmV0ZW50'
    'aW9uEkYKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMuQ2hlY2tzdW1BbGdvcml0aG'
    '1SEWNoZWNrc3VtYWxnb3JpdGhtEiwKBmRlbGV0ZRjr1N+8ASABKAsyCi5zMy5EZWxldGVCBIi1'
    'GAFSBmRlbGV0ZRIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2'
    'tldG93bmVyEhQKA21mYRjiqJWbASABKAlSA21mYRI4CgxyZXF1ZXN0cGF5ZXIYpObh9QEgASgO'
    'MhAuczMuUmVxdWVzdFBheWVyUgxyZXF1ZXN0cGF5ZXI=');

@$core.Deprecated('Use deletePublicAccessBlockRequestDescriptor instead')
const DeletePublicAccessBlockRequest$json = {
  '1': 'DeletePublicAccessBlockRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `DeletePublicAccessBlockRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePublicAccessBlockRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVQdWJsaWNBY2Nlc3NCbG9ja1JlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidW'
        'NrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25l'
        'cg==');

@$core.Deprecated('Use deletedObjectDescriptor instead')
const DeletedObject$json = {
  '1': 'DeletedObject',
  '2': [
    {'1': 'deletemarker', '3': 5472257, '4': 1, '5': 8, '10': 'deletemarker'},
    {
      '1': 'deletemarkerversionid',
      '3': 492617282,
      '4': 1,
      '5': 9,
      '10': 'deletemarkerversionid'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `DeletedObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletedObjectDescriptor = $convert.base64Decode(
    'Cg1EZWxldGVkT2JqZWN0EiUKDGRlbGV0ZW1hcmtlchiBgM4CIAEoCFIMZGVsZXRlbWFya2VyEj'
    'gKFWRlbGV0ZW1hcmtlcnZlcnNpb25pZBjC/PLqASABKAlSFWRlbGV0ZW1hcmtlcnZlcnNpb25p'
    'ZBITCgNrZXkYjZLraCABKAlSA2tleRIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaW'
    'Q=');

@$core.Deprecated('Use destinationDescriptor instead')
const Destination$json = {
  '1': 'Destination',
  '2': [
    {
      '1': 'accesscontroltranslation',
      '3': 450579688,
      '4': 1,
      '5': 11,
      '6': '.s3.AccessControlTranslation',
      '10': 'accesscontroltranslation'
    },
    {'1': 'account', '3': 435725053, '4': 1, '5': 9, '10': 'account'},
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.s3.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {
      '1': 'metrics',
      '3': 436365847,
      '4': 1,
      '5': 11,
      '6': '.s3.Metrics',
      '10': 'metrics'
    },
    {
      '1': 'replicationtime',
      '3': 98550465,
      '4': 1,
      '5': 11,
      '6': '.s3.ReplicationTime',
      '10': 'replicationtime'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
  ],
};

/// Descriptor for `Destination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List destinationDescriptor = $convert.base64Decode(
    'CgtEZXN0aW5hdGlvbhJcChhhY2Nlc3Njb250cm9sdHJhbnNsYXRpb24Y6Jnt1gEgASgLMhwucz'
    'MuQWNjZXNzQ29udHJvbFRyYW5zbGF0aW9uUhhhY2Nlc3Njb250cm9sdHJhbnNsYXRpb24SHAoH'
    'YWNjb3VudBj9xeLPASABKAlSB2FjY291bnQSGQoGYnVja2V0GNjquBogASgJUgZidWNrZXQSWA'
    'oXZW5jcnlwdGlvbmNvbmZpZ3VyYXRpb24Y98bTayABKAsyGy5zMy5FbmNyeXB0aW9uQ29uZmln'
    'dXJhdGlvblIXZW5jcnlwdGlvbmNvbmZpZ3VyYXRpb24SKQoHbWV0cmljcxiX1InQASABKAsyCy'
    '5zMy5NZXRyaWNzUgdtZXRyaWNzEkAKD3JlcGxpY2F0aW9udGltZRjBhf8uIAEoCzITLnMzLlJl'
    'cGxpY2F0aW9uVGltZVIPcmVwbGljYXRpb250aW1lEjgKDHN0b3JhZ2VjbGFzcxjHiMS7ASABKA'
    '4yEC5zMy5TdG9yYWdlQ2xhc3NSDHN0b3JhZ2VjbGFzcw==');

@$core.Deprecated('Use destinationResultDescriptor instead')
const DestinationResult$json = {
  '1': 'DestinationResult',
  '2': [
    {
      '1': 'tablebucketarn',
      '3': 457996261,
      '4': 1,
      '5': 9,
      '10': 'tablebucketarn'
    },
    {
      '1': 'tablebuckettype',
      '3': 37403814,
      '4': 1,
      '5': 14,
      '6': '.s3.S3TablesBucketType',
      '10': 'tablebuckettype'
    },
    {
      '1': 'tablenamespace',
      '3': 272949507,
      '4': 1,
      '5': 9,
      '10': 'tablenamespace'
    },
  ],
};

/// Descriptor for `DestinationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List destinationResultDescriptor = $convert.base64Decode(
    'ChFEZXN0aW5hdGlvblJlc3VsdBIqCg50YWJsZWJ1Y2tldGFybhjl77HaASABKAlSDnRhYmxlYn'
    'Vja2V0YXJuEkMKD3RhYmxlYnVja2V0dHlwZRim+eoRIAEoDjIWLnMzLlMzVGFibGVzQnVja2V0'
    'VHlwZVIPdGFibGVidWNrZXR0eXBlEioKDnRhYmxlbmFtZXNwYWNlGIPCk4IBIAEoCVIOdGFibG'
    'VuYW1lc3BhY2U=');

@$core.Deprecated('Use encryptionDescriptor instead')
const Encryption$json = {
  '1': 'Encryption',
  '2': [
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'encryptiontype'
    },
    {'1': 'kmscontext', '3': 161224782, '4': 1, '5': 9, '10': 'kmscontext'},
    {'1': 'kmskeyid', '3': 13237581, '4': 1, '5': 9, '10': 'kmskeyid'},
  ],
};

/// Descriptor for `Encryption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptionDescriptor = $convert.base64Decode(
    'CgpFbmNyeXB0aW9uEkMKDmVuY3J5cHRpb250eXBlGLXf8X0gASgOMhguczMuU2VydmVyU2lkZU'
    'VuY3J5cHRpb25SDmVuY3J5cHRpb250eXBlEiEKCmttc2NvbnRleHQYzrDwTCABKAlSCmttc2Nv'
    'bnRleHQSHQoIa21za2V5aWQYzfqnBiABKAlSCGttc2tleWlk');

@$core.Deprecated('Use encryptionConfigurationDescriptor instead')
const EncryptionConfiguration$json = {
  '1': 'EncryptionConfiguration',
  '2': [
    {
      '1': 'replicakmskeyid',
      '3': 77830171,
      '4': 1,
      '5': 9,
      '10': 'replicakmskeyid'
    },
  ],
};

/// Descriptor for `EncryptionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptionConfigurationDescriptor =
    $convert.base64Decode(
        'ChdFbmNyeXB0aW9uQ29uZmlndXJhdGlvbhIrCg9yZXBsaWNha21za2V5aWQYm7COJSABKAlSD3'
        'JlcGxpY2FrbXNrZXlpZA==');

@$core.Deprecated('Use encryptionTypeMismatchDescriptor instead')
const EncryptionTypeMismatch$json = {
  '1': 'EncryptionTypeMismatch',
};

/// Descriptor for `EncryptionTypeMismatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptionTypeMismatchDescriptor =
    $convert.base64Decode('ChZFbmNyeXB0aW9uVHlwZU1pc21hdGNo');

@$core.Deprecated('Use endEventDescriptor instead')
const EndEvent$json = {
  '1': 'EndEvent',
};

/// Descriptor for `EndEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endEventDescriptor =
    $convert.base64Decode('CghFbmRFdmVudA==');

@$core.Deprecated('Use errorDescriptor instead')
const Error$json = {
  '1': 'Error',
  '2': [
    {'1': 'code', '3': 425572629, '4': 1, '5': 9, '10': 'code'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `Error`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorDescriptor = $convert.base64Decode(
    'CgVFcnJvchIWCgRjb2RlGJXy9soBIAEoCVIEY29kZRITCgNrZXkYjZLraCABKAlSA2tleRIbCg'
    'dtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdlEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNp'
    'b25pZA==');

@$core.Deprecated('Use errorDetailsDescriptor instead')
const ErrorDetails$json = {
  '1': 'ErrorDetails',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
  ],
};

/// Descriptor for `ErrorDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorDetailsDescriptor = $convert.base64Decode(
    'CgxFcnJvckRldGFpbHMSHwoJZXJyb3Jjb2RlGJnWwxAgASgJUgllcnJvcmNvZGUSJgoMZXJyb3'
    'JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3JtZXNzYWdl');

@$core.Deprecated('Use errorDocumentDescriptor instead')
const ErrorDocument$json = {
  '1': 'ErrorDocument',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
  ],
};

/// Descriptor for `ErrorDocument`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorDocumentDescriptor =
    $convert.base64Decode('Cg1FcnJvckRvY3VtZW50EhMKA2tleRiNkutoIAEoCVIDa2V5');

@$core.Deprecated('Use eventBridgeConfigurationDescriptor instead')
const EventBridgeConfiguration$json = {
  '1': 'EventBridgeConfiguration',
};

/// Descriptor for `EventBridgeConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventBridgeConfigurationDescriptor =
    $convert.base64Decode('ChhFdmVudEJyaWRnZUNvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use existingObjectReplicationDescriptor instead')
const ExistingObjectReplication$json = {
  '1': 'ExistingObjectReplication',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.ExistingObjectReplicationStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ExistingObjectReplication`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List existingObjectReplicationDescriptor =
    $convert.base64Decode(
        'ChlFeGlzdGluZ09iamVjdFJlcGxpY2F0aW9uEj4KBnN0YXR1cxiQ5PsCIAEoDjIjLnMzLkV4aX'
        'N0aW5nT2JqZWN0UmVwbGljYXRpb25TdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use filterRuleDescriptor instead')
const FilterRule$json = {
  '1': 'FilterRule',
  '2': [
    {
      '1': 'name',
      '3': 266367751,
      '4': 1,
      '5': 14,
      '6': '.s3.FilterRuleName',
      '10': 'name'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `FilterRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterRuleDescriptor = $convert.base64Decode(
    'CgpGaWx0ZXJSdWxlEikKBG5hbWUYh+aBfyABKA4yEi5zMy5GaWx0ZXJSdWxlTmFtZVIEbmFtZR'
    'IYCgV2YWx1ZRjr8p+KASABKAlSBXZhbHVl');

@$core.Deprecated('Use getBucketAbacOutputDescriptor instead')
const GetBucketAbacOutput$json = {
  '1': 'GetBucketAbacOutput',
  '2': [
    {
      '1': 'abacstatus',
      '3': 436325625,
      '4': 1,
      '5': 11,
      '6': '.s3.AbacStatus',
      '8': {},
      '10': 'abacstatus'
    },
  ],
};

/// Descriptor for `GetBucketAbacOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketAbacOutputDescriptor = $convert.base64Decode(
    'ChNHZXRCdWNrZXRBYmFjT3V0cHV0EjgKCmFiYWNzdGF0dXMY+ZmH0AEgASgLMg4uczMuQWJhY1'
    'N0YXR1c0IEiLUYAVIKYWJhY3N0YXR1cw==');

@$core.Deprecated('Use getBucketAbacRequestDescriptor instead')
const GetBucketAbacRequest$json = {
  '1': 'GetBucketAbacRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketAbacRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketAbacRequestDescriptor = $convert.base64Decode(
    'ChRHZXRCdWNrZXRBYmFjUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzChNleH'
    'BlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core
    .Deprecated('Use getBucketAccelerateConfigurationOutputDescriptor instead')
const GetBucketAccelerateConfigurationOutput$json = {
  '1': 'GetBucketAccelerateConfigurationOutput',
  '2': [
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketAccelerateStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetBucketAccelerateConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketAccelerateConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiZHZXRCdWNrZXRBY2NlbGVyYXRlQ29uZmlndXJhdGlvbk91dHB1dBI+Cg5yZXF1ZXN0Y2hhcm'
        'dlZBiT0Ku5ASABKA4yEi5zMy5SZXF1ZXN0Q2hhcmdlZFIOcmVxdWVzdGNoYXJnZWQSNQoGc3Rh'
        'dHVzGJDk+wIgASgOMhouczMuQnVja2V0QWNjZWxlcmF0ZVN0YXR1c1IGc3RhdHVz');

@$core
    .Deprecated('Use getBucketAccelerateConfigurationRequestDescriptor instead')
const GetBucketAccelerateConfigurationRequest$json = {
  '1': 'GetBucketAccelerateConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
  ],
};

/// Descriptor for `GetBucketAccelerateConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketAccelerateConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CidHZXRCdWNrZXRBY2NlbGVyYXRlQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquB'
        'ogASgJUgZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRi'
        'dWNrZXRvd25lchI4CgxyZXF1ZXN0cGF5ZXIYpObh9QEgASgOMhAuczMuUmVxdWVzdFBheWVyUg'
        'xyZXF1ZXN0cGF5ZXI=');

@$core.Deprecated('Use getBucketAclOutputDescriptor instead')
const GetBucketAclOutput$json = {
  '1': 'GetBucketAclOutput',
  '2': [
    {
      '1': 'grants',
      '3': 226910555,
      '4': 3,
      '5': 11,
      '6': '.s3.Grant',
      '10': 'grants'
    },
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
  ],
};

/// Descriptor for `GetBucketAclOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketAclOutputDescriptor = $convert.base64Decode(
    'ChJHZXRCdWNrZXRBY2xPdXRwdXQSJAoGZ3JhbnRzGNvCmWwgAygLMgkuczMuR3JhbnRSBmdyYW'
    '50cxIjCgVvd25lchj1/IrZASABKAsyCS5zMy5Pd25lclIFb3duZXI=');

@$core.Deprecated('Use getBucketAclRequestDescriptor instead')
const GetBucketAclRequest$json = {
  '1': 'GetBucketAclRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketAclRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketAclRequestDescriptor = $convert.base64Decode(
    'ChNHZXRCdWNrZXRBY2xSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EjMKE2V4cG'
    'VjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXI=');

@$core.Deprecated('Use getBucketAnalyticsConfigurationOutputDescriptor instead')
const GetBucketAnalyticsConfigurationOutput$json = {
  '1': 'GetBucketAnalyticsConfigurationOutput',
  '2': [
    {
      '1': 'analyticsconfiguration',
      '3': 229750388,
      '4': 1,
      '5': 11,
      '6': '.s3.AnalyticsConfiguration',
      '8': {},
      '10': 'analyticsconfiguration'
    },
  ],
};

/// Descriptor for `GetBucketAnalyticsConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketAnalyticsConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiVHZXRCdWNrZXRBbmFseXRpY3NDb25maWd1cmF0aW9uT3V0cHV0ElsKFmFuYWx5dGljc2Nvbm'
        'ZpZ3VyYXRpb24Y9OzGbSABKAsyGi5zMy5BbmFseXRpY3NDb25maWd1cmF0aW9uQgSItRgBUhZh'
        'bmFseXRpY3Njb25maWd1cmF0aW9u');

@$core
    .Deprecated('Use getBucketAnalyticsConfigurationRequestDescriptor instead')
const GetBucketAnalyticsConfigurationRequest$json = {
  '1': 'GetBucketAnalyticsConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetBucketAnalyticsConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketAnalyticsConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiZHZXRCdWNrZXRBbmFseXRpY3NDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4Gi'
        'ABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1'
        'Y2tldG93bmVyEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getBucketCorsOutputDescriptor instead')
const GetBucketCorsOutput$json = {
  '1': 'GetBucketCorsOutput',
  '2': [
    {
      '1': 'corsrules',
      '3': 280917788,
      '4': 3,
      '5': 11,
      '6': '.s3.CORSRule',
      '10': 'corsrules'
    },
  ],
};

/// Descriptor for `GetBucketCorsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketCorsOutputDescriptor = $convert.base64Decode(
    'ChNHZXRCdWNrZXRDb3JzT3V0cHV0Ei4KCWNvcnNydWxlcxic7vmFASADKAsyDC5zMy5DT1JTUn'
    'VsZVIJY29yc3J1bGVz');

@$core.Deprecated('Use getBucketCorsRequestDescriptor instead')
const GetBucketCorsRequest$json = {
  '1': 'GetBucketCorsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketCorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketCorsRequestDescriptor = $convert.base64Decode(
    'ChRHZXRCdWNrZXRDb3JzUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzChNleH'
    'BlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use getBucketEncryptionOutputDescriptor instead')
const GetBucketEncryptionOutput$json = {
  '1': 'GetBucketEncryptionOutput',
  '2': [
    {
      '1': 'serversideencryptionconfiguration',
      '3': 314603363,
      '4': 1,
      '5': 11,
      '6': '.s3.ServerSideEncryptionConfiguration',
      '8': {},
      '10': 'serversideencryptionconfiguration'
    },
  ],
};

/// Descriptor for `GetBucketEncryptionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketEncryptionOutputDescriptor = $convert.base64Decode(
    'ChlHZXRCdWNrZXRFbmNyeXB0aW9uT3V0cHV0En0KIXNlcnZlcnNpZGVlbmNyeXB0aW9uY29uZm'
    'lndXJhdGlvbhjj7oGWASABKAsyJS5zMy5TZXJ2ZXJTaWRlRW5jcnlwdGlvbkNvbmZpZ3VyYXRp'
    'b25CBIi1GAFSIXNlcnZlcnNpZGVlbmNyeXB0aW9uY29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use getBucketEncryptionRequestDescriptor instead')
const GetBucketEncryptionRequest$json = {
  '1': 'GetBucketEncryptionRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketEncryptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketEncryptionRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRCdWNrZXRFbmNyeXB0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
        'IzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated(
    'Use getBucketIntelligentTieringConfigurationOutputDescriptor instead')
const GetBucketIntelligentTieringConfigurationOutput$json = {
  '1': 'GetBucketIntelligentTieringConfigurationOutput',
  '2': [
    {
      '1': 'intelligenttieringconfiguration',
      '3': 381057245,
      '4': 1,
      '5': 11,
      '6': '.s3.IntelligentTieringConfiguration',
      '8': {},
      '10': 'intelligenttieringconfiguration'
    },
  ],
};

/// Descriptor for `GetBucketIntelligentTieringConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getBucketIntelligentTieringConfigurationOutputDescriptor =
    $convert.base64Decode(
        'Ci5HZXRCdWNrZXRJbnRlbGxpZ2VudFRpZXJpbmdDb25maWd1cmF0aW9uT3V0cHV0EncKH2ludG'
        'VsbGlnZW50dGllcmluZ2NvbmZpZ3VyYXRpb24Y3fHZtQEgASgLMiMuczMuSW50ZWxsaWdlbnRU'
        'aWVyaW5nQ29uZmlndXJhdGlvbkIEiLUYAVIfaW50ZWxsaWdlbnR0aWVyaW5nY29uZmlndXJhdG'
        'lvbg==');

@$core.Deprecated(
    'Use getBucketIntelligentTieringConfigurationRequestDescriptor instead')
const GetBucketIntelligentTieringConfigurationRequest$json = {
  '1': 'GetBucketIntelligentTieringConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetBucketIntelligentTieringConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getBucketIntelligentTieringConfigurationRequestDescriptor =
    $convert.base64Decode(
        'Ci9HZXRCdWNrZXRJbnRlbGxpZ2VudFRpZXJpbmdDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidW'
        'NrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNl'
        'eHBlY3RlZGJ1Y2tldG93bmVyEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getBucketInventoryConfigurationOutputDescriptor instead')
const GetBucketInventoryConfigurationOutput$json = {
  '1': 'GetBucketInventoryConfigurationOutput',
  '2': [
    {
      '1': 'inventoryconfiguration',
      '3': 401695392,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryConfiguration',
      '8': {},
      '10': 'inventoryconfiguration'
    },
  ],
};

/// Descriptor for `GetBucketInventoryConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketInventoryConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiVHZXRCdWNrZXRJbnZlbnRvcnlDb25maWd1cmF0aW9uT3V0cHV0ElwKFmludmVudG9yeWNvbm'
        'ZpZ3VyYXRpb24YoMXFvwEgASgLMhouczMuSW52ZW50b3J5Q29uZmlndXJhdGlvbkIEiLUYAVIW'
        'aW52ZW50b3J5Y29uZmlndXJhdGlvbg==');

@$core
    .Deprecated('Use getBucketInventoryConfigurationRequestDescriptor instead')
const GetBucketInventoryConfigurationRequest$json = {
  '1': 'GetBucketInventoryConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetBucketInventoryConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketInventoryConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiZHZXRCdWNrZXRJbnZlbnRvcnlDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4Gi'
        'ABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1'
        'Y2tldG93bmVyEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use getBucketLifecycleConfigurationOutputDescriptor instead')
const GetBucketLifecycleConfigurationOutput$json = {
  '1': 'GetBucketLifecycleConfigurationOutput',
  '2': [
    {
      '1': 'rules',
      '3': 42675585,
      '4': 3,
      '5': 11,
      '6': '.s3.LifecycleRule',
      '10': 'rules'
    },
    {
      '1': 'transitiondefaultminimumobjectsize',
      '3': 81614768,
      '4': 1,
      '5': 14,
      '6': '.s3.TransitionDefaultMinimumObjectSize',
      '10': 'transitiondefaultminimumobjectsize'
    },
  ],
};

/// Descriptor for `GetBucketLifecycleConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketLifecycleConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiVHZXRCdWNrZXRMaWZlY3ljbGVDb25maWd1cmF0aW9uT3V0cHV0EioKBXJ1bGVzGIHbrBQgAy'
        'gLMhEuczMuTGlmZWN5Y2xlUnVsZVIFcnVsZXMSeQoidHJhbnNpdGlvbmRlZmF1bHRtaW5pbXVt'
        'b2JqZWN0c2l6ZRiwr/UmIAEoDjImLnMzLlRyYW5zaXRpb25EZWZhdWx0TWluaW11bU9iamVjdF'
        'NpemVSInRyYW5zaXRpb25kZWZhdWx0bWluaW11bW9iamVjdHNpemU=');

@$core
    .Deprecated('Use getBucketLifecycleConfigurationRequestDescriptor instead')
const GetBucketLifecycleConfigurationRequest$json = {
  '1': 'GetBucketLifecycleConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketLifecycleConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketLifecycleConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiZHZXRCdWNrZXRMaWZlY3ljbGVDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4Gi'
        'ABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1'
        'Y2tldG93bmVy');

@$core.Deprecated('Use getBucketLocationOutputDescriptor instead')
const GetBucketLocationOutput$json = {
  '1': 'GetBucketLocationOutput',
  '2': [
    {
      '1': 'locationconstraint',
      '3': 245375838,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketLocationConstraint',
      '10': 'locationconstraint'
    },
  ],
};

/// Descriptor for `GetBucketLocationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketLocationOutputDescriptor = $convert.base64Decode(
    'ChdHZXRCdWNrZXRMb2NhdGlvbk91dHB1dBJPChJsb2NhdGlvbmNvbnN0cmFpbnQY3saAdSABKA'
    '4yHC5zMy5CdWNrZXRMb2NhdGlvbkNvbnN0cmFpbnRSEmxvY2F0aW9uY29uc3RyYWludA==');

@$core.Deprecated('Use getBucketLocationRequestDescriptor instead')
const GetBucketLocationRequest$json = {
  '1': 'GetBucketLocationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketLocationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketLocationRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRCdWNrZXRMb2NhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZXQSMw'
        'oTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lcg==');

@$core.Deprecated('Use getBucketLoggingOutputDescriptor instead')
const GetBucketLoggingOutput$json = {
  '1': 'GetBucketLoggingOutput',
  '2': [
    {
      '1': 'loggingenabled',
      '3': 119537244,
      '4': 1,
      '5': 11,
      '6': '.s3.LoggingEnabled',
      '10': 'loggingenabled'
    },
  ],
};

/// Descriptor for `GetBucketLoggingOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketLoggingOutputDescriptor =
    $convert.base64Decode(
        'ChZHZXRCdWNrZXRMb2dnaW5nT3V0cHV0Ej0KDmxvZ2dpbmdlbmFibGVkGNz8/zggASgLMhIucz'
        'MuTG9nZ2luZ0VuYWJsZWRSDmxvZ2dpbmdlbmFibGVk');

@$core.Deprecated('Use getBucketLoggingRequestDescriptor instead')
const GetBucketLoggingRequest$json = {
  '1': 'GetBucketLoggingRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketLoggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketLoggingRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRCdWNrZXRMb2dnaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzCh'
        'NleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use getBucketMetadataConfigurationOutputDescriptor instead')
const GetBucketMetadataConfigurationOutput$json = {
  '1': 'GetBucketMetadataConfigurationOutput',
  '2': [
    {
      '1': 'getbucketmetadataconfigurationresult',
      '3': 525444508,
      '4': 1,
      '5': 11,
      '6': '.s3.GetBucketMetadataConfigurationResult',
      '8': {},
      '10': 'getbucketmetadataconfigurationresult'
    },
  ],
};

/// Descriptor for `GetBucketMetadataConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketMetadataConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiRHZXRCdWNrZXRNZXRhZGF0YUNvbmZpZ3VyYXRpb25PdXRwdXQShgEKJGdldGJ1Y2tldG1ldG'
        'FkYXRhY29uZmlndXJhdGlvbnJlc3VsdBicy8b6ASABKAsyKC5zMy5HZXRCdWNrZXRNZXRhZGF0'
        'YUNvbmZpZ3VyYXRpb25SZXN1bHRCBIi1GAFSJGdldGJ1Y2tldG1ldGFkYXRhY29uZmlndXJhdG'
        'lvbnJlc3VsdA==');

@$core.Deprecated('Use getBucketMetadataConfigurationRequestDescriptor instead')
const GetBucketMetadataConfigurationRequest$json = {
  '1': 'GetBucketMetadataConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketMetadataConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketMetadataConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiVHZXRCdWNrZXRNZXRhZGF0YUNvbmZpZ3VyYXRpb25SZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIA'
        'EoCVIGYnVja2V0EjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVj'
        'a2V0b3duZXI=');

@$core.Deprecated('Use getBucketMetadataConfigurationResultDescriptor instead')
const GetBucketMetadataConfigurationResult$json = {
  '1': 'GetBucketMetadataConfigurationResult',
  '2': [
    {
      '1': 'metadataconfigurationresult',
      '3': 151901996,
      '4': 1,
      '5': 11,
      '6': '.s3.MetadataConfigurationResult',
      '10': 'metadataconfigurationresult'
    },
  ],
};

/// Descriptor for `GetBucketMetadataConfigurationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketMetadataConfigurationResultDescriptor =
    $convert.base64Decode(
        'CiRHZXRCdWNrZXRNZXRhZGF0YUNvbmZpZ3VyYXRpb25SZXN1bHQSZAobbWV0YWRhdGFjb25maW'
        'd1cmF0aW9ucmVzdWx0GKyut0ggASgLMh8uczMuTWV0YWRhdGFDb25maWd1cmF0aW9uUmVzdWx0'
        'UhttZXRhZGF0YWNvbmZpZ3VyYXRpb25yZXN1bHQ=');

@$core.Deprecated(
    'Use getBucketMetadataTableConfigurationOutputDescriptor instead')
const GetBucketMetadataTableConfigurationOutput$json = {
  '1': 'GetBucketMetadataTableConfigurationOutput',
  '2': [
    {
      '1': 'getbucketmetadatatableconfigurationresult',
      '3': 330020248,
      '4': 1,
      '5': 11,
      '6': '.s3.GetBucketMetadataTableConfigurationResult',
      '8': {},
      '10': 'getbucketmetadatatableconfigurationresult'
    },
  ],
};

/// Descriptor for `GetBucketMetadataTableConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getBucketMetadataTableConfigurationOutputDescriptor = $convert.base64Decode(
        'CilHZXRCdWNrZXRNZXRhZGF0YVRhYmxlQ29uZmlndXJhdGlvbk91dHB1dBKVAQopZ2V0YnVja2'
        'V0bWV0YWRhdGF0YWJsZWNvbmZpZ3VyYXRpb25yZXN1bHQYmOuunQEgASgLMi0uczMuR2V0QnVj'
        'a2V0TWV0YWRhdGFUYWJsZUNvbmZpZ3VyYXRpb25SZXN1bHRCBIi1GAFSKWdldGJ1Y2tldG1ldG'
        'FkYXRhdGFibGVjb25maWd1cmF0aW9ucmVzdWx0');

@$core.Deprecated(
    'Use getBucketMetadataTableConfigurationRequestDescriptor instead')
const GetBucketMetadataTableConfigurationRequest$json = {
  '1': 'GetBucketMetadataTableConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketMetadataTableConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getBucketMetadataTableConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CipHZXRCdWNrZXRNZXRhZGF0YVRhYmxlQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2V0GN'
        'jquBogASgJUgZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0'
        'ZWRidWNrZXRvd25lcg==');

@$core.Deprecated(
    'Use getBucketMetadataTableConfigurationResultDescriptor instead')
const GetBucketMetadataTableConfigurationResult$json = {
  '1': 'GetBucketMetadataTableConfigurationResult',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.s3.ErrorDetails',
      '10': 'error'
    },
    {
      '1': 'metadatatableconfigurationresult',
      '3': 180591208,
      '4': 1,
      '5': 11,
      '6': '.s3.MetadataTableConfigurationResult',
      '10': 'metadatatableconfigurationresult'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `GetBucketMetadataTableConfigurationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getBucketMetadataTableConfigurationResultDescriptor = $convert.base64Decode(
        'CilHZXRCdWNrZXRNZXRhZGF0YVRhYmxlQ29uZmlndXJhdGlvblJlc3VsdBIqCgVlcnJvchjyub'
        'acASABKAsyEC5zMy5FcnJvckRldGFpbHNSBWVycm9yEnMKIG1ldGFkYXRhdGFibGVjb25maWd1'
        'cmF0aW9ucmVzdWx0GOi0jlYgASgLMiQuczMuTWV0YWRhdGFUYWJsZUNvbmZpZ3VyYXRpb25SZX'
        'N1bHRSIG1ldGFkYXRhdGFibGVjb25maWd1cmF0aW9ucmVzdWx0EhkKBnN0YXR1cxiQ5PsCIAEo'
        'CVIGc3RhdHVz');

@$core.Deprecated('Use getBucketMetricsConfigurationOutputDescriptor instead')
const GetBucketMetricsConfigurationOutput$json = {
  '1': 'GetBucketMetricsConfigurationOutput',
  '2': [
    {
      '1': 'metricsconfiguration',
      '3': 14504541,
      '4': 1,
      '5': 11,
      '6': '.s3.MetricsConfiguration',
      '8': {},
      '10': 'metricsconfiguration'
    },
  ],
};

/// Descriptor for `GetBucketMetricsConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketMetricsConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiNHZXRCdWNrZXRNZXRyaWNzQ29uZmlndXJhdGlvbk91dHB1dBJVChRtZXRyaWNzY29uZmlndX'
        'JhdGlvbhjdpPUGIAEoCzIYLnMzLk1ldHJpY3NDb25maWd1cmF0aW9uQgSItRgBUhRtZXRyaWNz'
        'Y29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use getBucketMetricsConfigurationRequestDescriptor instead')
const GetBucketMetricsConfigurationRequest$json = {
  '1': 'GetBucketMetricsConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `GetBucketMetricsConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketMetricsConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRCdWNrZXRNZXRyaWNzQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquBogAS'
        'gJUgZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNr'
        'ZXRvd25lchISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated(
    'Use getBucketNotificationConfigurationRequestDescriptor instead')
const GetBucketNotificationConfigurationRequest$json = {
  '1': 'GetBucketNotificationConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketNotificationConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getBucketNotificationConfigurationRequestDescriptor = $convert.base64Decode(
        'CilHZXRCdWNrZXROb3RpZmljYXRpb25Db25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2O'
        'q4GiABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3Rl'
        'ZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use getBucketOwnershipControlsOutputDescriptor instead')
const GetBucketOwnershipControlsOutput$json = {
  '1': 'GetBucketOwnershipControlsOutput',
  '2': [
    {
      '1': 'ownershipcontrols',
      '3': 6169057,
      '4': 1,
      '5': 11,
      '6': '.s3.OwnershipControls',
      '8': {},
      '10': 'ownershipcontrols'
    },
  ],
};

/// Descriptor for `GetBucketOwnershipControlsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketOwnershipControlsOutputDescriptor =
    $convert.base64Decode(
        'CiBHZXRCdWNrZXRPd25lcnNoaXBDb250cm9sc091dHB1dBJMChFvd25lcnNoaXBjb250cm9scx'
        'jhw/gCIAEoCzIVLnMzLk93bmVyc2hpcENvbnRyb2xzQgSItRgBUhFvd25lcnNoaXBjb250cm9s'
        'cw==');

@$core.Deprecated('Use getBucketOwnershipControlsRequestDescriptor instead')
const GetBucketOwnershipControlsRequest$json = {
  '1': 'GetBucketOwnershipControlsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketOwnershipControlsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketOwnershipControlsRequestDescriptor =
    $convert.base64Decode(
        'CiFHZXRCdWNrZXRPd25lcnNoaXBDb250cm9sc1JlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUg'
        'ZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRv'
        'd25lcg==');

@$core.Deprecated('Use getBucketPolicyOutputDescriptor instead')
const GetBucketPolicyOutput$json = {
  '1': 'GetBucketPolicyOutput',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '8': {}, '10': 'policy'},
  ],
};

/// Descriptor for `GetBucketPolicyOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketPolicyOutputDescriptor = $convert.base64Decode(
    'ChVHZXRCdWNrZXRQb2xpY3lPdXRwdXQSIAoGcG9saWN5GKDv8OABIAEoCUIEiLUYAVIGcG9saW'
    'N5');

@$core.Deprecated('Use getBucketPolicyRequestDescriptor instead')
const GetBucketPolicyRequest$json = {
  '1': 'GetBucketPolicyRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketPolicyRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRCdWNrZXRQb2xpY3lSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EjMKE2'
        'V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXI=');

@$core.Deprecated('Use getBucketPolicyStatusOutputDescriptor instead')
const GetBucketPolicyStatusOutput$json = {
  '1': 'GetBucketPolicyStatusOutput',
  '2': [
    {
      '1': 'policystatus',
      '3': 256036138,
      '4': 1,
      '5': 11,
      '6': '.s3.PolicyStatus',
      '8': {},
      '10': 'policystatus'
    },
  ],
};

/// Descriptor for `GetBucketPolicyStatusOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketPolicyStatusOutputDescriptor =
    $convert.base64Decode(
        'ChtHZXRCdWNrZXRQb2xpY3lTdGF0dXNPdXRwdXQSPQoMcG9saWN5c3RhdHVzGKqai3ogASgLMh'
        'AuczMuUG9saWN5U3RhdHVzQgSItRgBUgxwb2xpY3lzdGF0dXM=');

@$core.Deprecated('Use getBucketPolicyStatusRequestDescriptor instead')
const GetBucketPolicyStatusRequest$json = {
  '1': 'GetBucketPolicyStatusRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketPolicyStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketPolicyStatusRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRCdWNrZXRQb2xpY3lTdGF0dXNSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2'
        'V0EjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXI=');

@$core.Deprecated('Use getBucketReplicationOutputDescriptor instead')
const GetBucketReplicationOutput$json = {
  '1': 'GetBucketReplicationOutput',
  '2': [
    {
      '1': 'replicationconfiguration',
      '3': 115861462,
      '4': 1,
      '5': 11,
      '6': '.s3.ReplicationConfiguration',
      '8': {},
      '10': 'replicationconfiguration'
    },
  ],
};

/// Descriptor for `GetBucketReplicationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketReplicationOutputDescriptor =
    $convert.base64Decode(
        'ChpHZXRCdWNrZXRSZXBsaWNhdGlvbk91dHB1dBJhChhyZXBsaWNhdGlvbmNvbmZpZ3VyYXRpb2'
        '4Y1s+fNyABKAsyHC5zMy5SZXBsaWNhdGlvbkNvbmZpZ3VyYXRpb25CBIi1GAFSGHJlcGxpY2F0'
        'aW9uY29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use getBucketReplicationRequestDescriptor instead')
const GetBucketReplicationRequest$json = {
  '1': 'GetBucketReplicationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketReplicationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketReplicationRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRCdWNrZXRSZXBsaWNhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZX'
        'QSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lcg==');

@$core.Deprecated('Use getBucketRequestPaymentOutputDescriptor instead')
const GetBucketRequestPaymentOutput$json = {
  '1': 'GetBucketRequestPaymentOutput',
  '2': [
    {
      '1': 'payer',
      '3': 174621923,
      '4': 1,
      '5': 14,
      '6': '.s3.Payer',
      '10': 'payer'
    },
  ],
};

/// Descriptor for `GetBucketRequestPaymentOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketRequestPaymentOutputDescriptor =
    $convert.base64Decode(
        'Ch1HZXRCdWNrZXRSZXF1ZXN0UGF5bWVudE91dHB1dBIiCgVwYXllchjjiaJTIAEoDjIJLnMzLl'
        'BheWVyUgVwYXllcg==');

@$core.Deprecated('Use getBucketRequestPaymentRequestDescriptor instead')
const GetBucketRequestPaymentRequest$json = {
  '1': 'GetBucketRequestPaymentRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketRequestPaymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketRequestPaymentRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXRCdWNrZXRSZXF1ZXN0UGF5bWVudFJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidW'
        'NrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25l'
        'cg==');

@$core.Deprecated('Use getBucketTaggingOutputDescriptor instead')
const GetBucketTaggingOutput$json = {
  '1': 'GetBucketTaggingOutput',
  '2': [
    {
      '1': 'tagset',
      '3': 454361330,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tagset'
    },
  ],
};

/// Descriptor for `GetBucketTaggingOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketTaggingOutputDescriptor =
    $convert.base64Decode(
        'ChZHZXRCdWNrZXRUYWdnaW5nT3V0cHV0EiMKBnRhZ3NldBjygdTYASADKAsyBy5zMy5UYWdSBn'
        'RhZ3NldA==');

@$core.Deprecated('Use getBucketTaggingRequestDescriptor instead')
const GetBucketTaggingRequest$json = {
  '1': 'GetBucketTaggingRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketTaggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketTaggingRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRCdWNrZXRUYWdnaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzCh'
        'NleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use getBucketVersioningOutputDescriptor instead')
const GetBucketVersioningOutput$json = {
  '1': 'GetBucketVersioningOutput',
  '2': [
    {
      '1': 'mfadelete',
      '3': 354119319,
      '4': 1,
      '5': 14,
      '6': '.s3.MFADeleteStatus',
      '10': 'mfadelete'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketVersioningStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetBucketVersioningOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketVersioningOutputDescriptor = $convert.base64Decode(
    'ChlHZXRCdWNrZXRWZXJzaW9uaW5nT3V0cHV0EjUKCW1mYWRlbGV0ZRiX3e2oASABKA4yEy5zMy'
    '5NRkFEZWxldGVTdGF0dXNSCW1mYWRlbGV0ZRI1CgZzdGF0dXMYkOT7AiABKA4yGi5zMy5CdWNr'
    'ZXRWZXJzaW9uaW5nU3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use getBucketVersioningRequestDescriptor instead')
const GetBucketVersioningRequest$json = {
  '1': 'GetBucketVersioningRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketVersioningRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketVersioningRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRCdWNrZXRWZXJzaW9uaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
        'IzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use getBucketWebsiteOutputDescriptor instead')
const GetBucketWebsiteOutput$json = {
  '1': 'GetBucketWebsiteOutput',
  '2': [
    {
      '1': 'errordocument',
      '3': 329194857,
      '4': 1,
      '5': 11,
      '6': '.s3.ErrorDocument',
      '10': 'errordocument'
    },
    {
      '1': 'indexdocument',
      '3': 23264623,
      '4': 1,
      '5': 11,
      '6': '.s3.IndexDocument',
      '10': 'indexdocument'
    },
    {
      '1': 'redirectallrequeststo',
      '3': 317142372,
      '4': 1,
      '5': 11,
      '6': '.s3.RedirectAllRequestsTo',
      '10': 'redirectallrequeststo'
    },
    {
      '1': 'routingrules',
      '3': 83486933,
      '4': 3,
      '5': 11,
      '6': '.s3.RoutingRule',
      '10': 'routingrules'
    },
  ],
};

/// Descriptor for `GetBucketWebsiteOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketWebsiteOutputDescriptor = $convert.base64Decode(
    'ChZHZXRCdWNrZXRXZWJzaXRlT3V0cHV0EjsKDWVycm9yZG9jdW1lbnQY6br8nAEgASgLMhEucz'
    'MuRXJyb3JEb2N1bWVudFINZXJyb3Jkb2N1bWVudBI6Cg1pbmRleGRvY3VtZW50GO/6iwsgASgL'
    'MhEuczMuSW5kZXhEb2N1bWVudFINaW5kZXhkb2N1bWVudBJTChVyZWRpcmVjdGFsbHJlcXVlc3'
    'RzdG8Y5OqclwEgASgLMhkuczMuUmVkaXJlY3RBbGxSZXF1ZXN0c1RvUhVyZWRpcmVjdGFsbHJl'
    'cXVlc3RzdG8SNgoMcm91dGluZ3J1bGVzGNXR5ycgAygLMg8uczMuUm91dGluZ1J1bGVSDHJvdX'
    'RpbmdydWxlcw==');

@$core.Deprecated('Use getBucketWebsiteRequestDescriptor instead')
const GetBucketWebsiteRequest$json = {
  '1': 'GetBucketWebsiteRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetBucketWebsiteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBucketWebsiteRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRCdWNrZXRXZWJzaXRlUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzCh'
        'NleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use getObjectAclOutputDescriptor instead')
const GetObjectAclOutput$json = {
  '1': 'GetObjectAclOutput',
  '2': [
    {
      '1': 'grants',
      '3': 226910555,
      '4': 3,
      '5': 11,
      '6': '.s3.Grant',
      '10': 'grants'
    },
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `GetObjectAclOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectAclOutputDescriptor = $convert.base64Decode(
    'ChJHZXRPYmplY3RBY2xPdXRwdXQSJAoGZ3JhbnRzGNvCmWwgAygLMgkuczMuR3JhbnRSBmdyYW'
    '50cxIjCgVvd25lchj1/IrZASABKAsyCS5zMy5Pd25lclIFb3duZXISPgoOcmVxdWVzdGNoYXJn'
    'ZWQYk9CruQEgASgOMhIuczMuUmVxdWVzdENoYXJnZWRSDnJlcXVlc3RjaGFyZ2Vk');

@$core.Deprecated('Use getObjectAclRequestDescriptor instead')
const GetObjectAclRequest$json = {
  '1': 'GetObjectAclRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetObjectAclRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectAclRequestDescriptor = $convert.base64Decode(
    'ChNHZXRPYmplY3RBY2xSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EjMKE2V4cG'
    'VjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXISEwoDa2V5GI2S'
    '62ggASgJUgNrZXkSOAoMcmVxdWVzdHBheWVyGKTm4fUBIAEoDjIQLnMzLlJlcXVlc3RQYXllcl'
    'IMcmVxdWVzdHBheWVyEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated('Use getObjectAttributesOutputDescriptor instead')
const GetObjectAttributesOutput$json = {
  '1': 'GetObjectAttributesOutput',
  '2': [
    {
      '1': 'checksum',
      '3': 434578355,
      '4': 1,
      '5': 11,
      '6': '.s3.Checksum',
      '10': 'checksum'
    },
    {'1': 'deletemarker', '3': 5472257, '4': 1, '5': 8, '10': 'deletemarker'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'objectparts',
      '3': 402544889,
      '4': 1,
      '5': 11,
      '6': '.s3.GetObjectAttributesParts',
      '10': 'objectparts'
    },
    {'1': 'objectsize', '3': 130133988, '4': 1, '5': 3, '10': 'objectsize'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetObjectAttributesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectAttributesOutputDescriptor = $convert.base64Decode(
    'ChlHZXRPYmplY3RBdHRyaWJ1dGVzT3V0cHV0EiwKCGNoZWNrc3VtGLPHnM8BIAEoCzIMLnMzLk'
    'NoZWNrc3VtUghjaGVja3N1bRIlCgxkZWxldGVtYXJrZXIYgYDOAiABKAhSDGRlbGV0ZW1hcmtl'
    'chIWCgRldGFnGIHfs5UBIAEoCVIEZXRhZxImCgxsYXN0bW9kaWZpZWQYp5z8zgEgASgJUgxsYX'
    'N0bW9kaWZpZWQSQgoLb2JqZWN0cGFydHMY+bH5vwEgASgLMhwuczMuR2V0T2JqZWN0QXR0cmli'
    'dXRlc1BhcnRzUgtvYmplY3RwYXJ0cxIhCgpvYmplY3RzaXplGOTfhj4gASgDUgpvYmplY3RzaX'
    'plEj4KDnJlcXVlc3RjaGFyZ2VkGJPQq7kBIAEoDjISLnMzLlJlcXVlc3RDaGFyZ2VkUg5yZXF1'
    'ZXN0Y2hhcmdlZBI4CgxzdG9yYWdlY2xhc3MYx4jEuwEgASgOMhAuczMuU3RvcmFnZUNsYXNzUg'
    'xzdG9yYWdlY2xhc3MSIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlk');

@$core.Deprecated('Use getObjectAttributesPartsDescriptor instead')
const GetObjectAttributesParts$json = {
  '1': 'GetObjectAttributesParts',
  '2': [
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'maxparts', '3': 398005914, '4': 1, '5': 5, '10': 'maxparts'},
    {
      '1': 'nextpartnumbermarker',
      '3': 28931219,
      '4': 1,
      '5': 9,
      '10': 'nextpartnumbermarker'
    },
    {
      '1': 'partnumbermarker',
      '3': 376535672,
      '4': 1,
      '5': 9,
      '10': 'partnumbermarker'
    },
    {
      '1': 'parts',
      '3': 213028806,
      '4': 3,
      '5': 11,
      '6': '.s3.ObjectPart',
      '10': 'parts'
    },
    {
      '1': 'totalpartscount',
      '3': 56278317,
      '4': 1,
      '5': 5,
      '10': 'totalpartscount'
    },
  ],
};

/// Descriptor for `GetObjectAttributesParts`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectAttributesPartsDescriptor = $convert.base64Decode(
    'ChhHZXRPYmplY3RBdHRyaWJ1dGVzUGFydHMSIwoLaXN0cnVuY2F0ZWQY2p+4cyABKAhSC2lzdH'
    'J1bmNhdGVkEh4KCG1heHBhcnRzGJqt5L0BIAEoBVIIbWF4cGFydHMSNQoUbmV4dHBhcnRudW1i'
    'ZXJtYXJrZXIYk+nlDSABKAlSFG5leHRwYXJ0bnVtYmVybWFya2VyEi4KEHBhcnRudW1iZXJtYX'
    'JrZXIY+PTFswEgASgJUhBwYXJ0bnVtYmVybWFya2VyEicKBXBhcnRzGMafymUgAygLMg4uczMu'
    'T2JqZWN0UGFydFIFcGFydHMSKwoPdG90YWxwYXJ0c2NvdW50GK366hogASgFUg90b3RhbHBhcn'
    'RzY291bnQ=');

@$core.Deprecated('Use getObjectAttributesRequestDescriptor instead')
const GetObjectAttributesRequest$json = {
  '1': 'GetObjectAttributesRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'maxparts', '3': 398005914, '4': 1, '5': 5, '10': 'maxparts'},
    {
      '1': 'objectattributes',
      '3': 370999888,
      '4': 3,
      '5': 14,
      '6': '.s3.ObjectAttributes',
      '10': 'objectattributes'
    },
    {
      '1': 'partnumbermarker',
      '3': 376535672,
      '4': 1,
      '5': 9,
      '10': 'partnumbermarker'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetObjectAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectAttributesRequestDescriptor = $convert.base64Decode(
    'ChpHZXRPYmplY3RBdHRyaWJ1dGVzUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
    'IzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhMK'
    'A2tleRiNkutoIAEoCVIDa2V5Eh4KCG1heHBhcnRzGJqt5L0BIAEoBVIIbWF4cGFydHMSRAoQb2'
    'JqZWN0YXR0cmlidXRlcxjQhPSwASADKA4yFC5zMy5PYmplY3RBdHRyaWJ1dGVzUhBvYmplY3Rh'
    'dHRyaWJ1dGVzEi4KEHBhcnRudW1iZXJtYXJrZXIY+PTFswEgASgJUhBwYXJ0bnVtYmVybWFya2'
    'VyEjgKDHJlcXVlc3RwYXllchik5uH1ASABKA4yEC5zMy5SZXF1ZXN0UGF5ZXJSDHJlcXVlc3Rw'
    'YXllchI1ChRzc2VjdXN0b21lcmFsZ29yaXRobRjQyYErIAEoCVIUc3NlY3VzdG9tZXJhbGdvcm'
    'l0aG0SKQoOc3NlY3VzdG9tZXJrZXkYmv70OyABKAlSDnNzZWN1c3RvbWVya2V5Ei4KEXNzZWN1'
    'c3RvbWVya2V5bWQ1GOjRFyABKAlSEXNzZWN1c3RvbWVya2V5bWQ1EiAKCXZlcnNpb25pZBib4Z'
    'mhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated('Use getObjectLegalHoldOutputDescriptor instead')
const GetObjectLegalHoldOutput$json = {
  '1': 'GetObjectLegalHoldOutput',
  '2': [
    {
      '1': 'legalhold',
      '3': 141145940,
      '4': 1,
      '5': 11,
      '6': '.s3.ObjectLockLegalHold',
      '8': {},
      '10': 'legalhold'
    },
  ],
};

/// Descriptor for `GetObjectLegalHoldOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectLegalHoldOutputDescriptor =
    $convert.base64Decode(
        'ChhHZXRPYmplY3RMZWdhbEhvbGRPdXRwdXQSPgoJbGVnYWxob2xkGNTupkMgASgLMhcuczMuT2'
        'JqZWN0TG9ja0xlZ2FsSG9sZEIEiLUYAVIJbGVnYWxob2xk');

@$core.Deprecated('Use getObjectLegalHoldRequestDescriptor instead')
const GetObjectLegalHoldRequest$json = {
  '1': 'GetObjectLegalHoldRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetObjectLegalHoldRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectLegalHoldRequestDescriptor = $convert.base64Decode(
    'ChlHZXRPYmplY3RMZWdhbEhvbGRSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0Ej'
    'MKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXISEwoD'
    'a2V5GI2S62ggASgJUgNrZXkSOAoMcmVxdWVzdHBheWVyGKTm4fUBIAEoDjIQLnMzLlJlcXVlc3'
    'RQYXllclIMcmVxdWVzdHBheWVyEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated('Use getObjectLockConfigurationOutputDescriptor instead')
const GetObjectLockConfigurationOutput$json = {
  '1': 'GetObjectLockConfigurationOutput',
  '2': [
    {
      '1': 'objectlockconfiguration',
      '3': 108348298,
      '4': 1,
      '5': 11,
      '6': '.s3.ObjectLockConfiguration',
      '8': {},
      '10': 'objectlockconfiguration'
    },
  ],
};

/// Descriptor for `GetObjectLockConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectLockConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiBHZXRPYmplY3RMb2NrQ29uZmlndXJhdGlvbk91dHB1dBJeChdvYmplY3Rsb2NrY29uZmlndX'
        'JhdGlvbhiKh9UzIAEoCzIbLnMzLk9iamVjdExvY2tDb25maWd1cmF0aW9uQgSItRgBUhdvYmpl'
        'Y3Rsb2NrY29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use getObjectLockConfigurationRequestDescriptor instead')
const GetObjectLockConfigurationRequest$json = {
  '1': 'GetObjectLockConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetObjectLockConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectLockConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiFHZXRPYmplY3RMb2NrQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUg'
        'ZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRv'
        'd25lcg==');

@$core.Deprecated('Use getObjectOutputDescriptor instead')
const GetObjectOutput$json = {
  '1': 'GetObjectOutput',
  '2': [
    {'1': 'acceptranges', '3': 464620960, '4': 1, '5': 9, '10': 'acceptranges'},
    {'1': 'body', '3': 42602646, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {'1': 'cachecontrol', '3': 288966655, '4': 1, '5': 9, '10': 'cachecontrol'},
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {
      '1': 'contentdisposition',
      '3': 120040130,
      '4': 1,
      '5': 9,
      '10': 'contentdisposition'
    },
    {
      '1': 'contentencoding',
      '3': 317106228,
      '4': 1,
      '5': 9,
      '10': 'contentencoding'
    },
    {
      '1': 'contentlanguage',
      '3': 108485649,
      '4': 1,
      '5': 9,
      '10': 'contentlanguage'
    },
    {
      '1': 'contentlength',
      '3': 227596631,
      '4': 1,
      '5': 3,
      '10': 'contentlength'
    },
    {'1': 'contentrange', '3': 11089360, '4': 1, '5': 9, '10': 'contentrange'},
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'deletemarker', '3': 5472257, '4': 1, '5': 8, '10': 'deletemarker'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {'1': 'expires', '3': 128582948, '4': 1, '5': 9, '10': 'expires'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 3,
      '5': 11,
      '6': '.s3.GetObjectOutput.MetadataEntry',
      '10': 'metadata'
    },
    {'1': 'missingmeta', '3': 79140523, '4': 1, '5': 5, '10': 'missingmeta'},
    {
      '1': 'objectlocklegalholdstatus',
      '3': 536561974,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockLegalHoldStatus',
      '10': 'objectlocklegalholdstatus'
    },
    {
      '1': 'objectlockmode',
      '3': 189255203,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockMode',
      '10': 'objectlockmode'
    },
    {
      '1': 'objectlockretainuntildate',
      '3': 264584249,
      '4': 1,
      '5': 9,
      '10': 'objectlockretainuntildate'
    },
    {'1': 'partscount', '3': 154996373, '4': 1, '5': 5, '10': 'partscount'},
    {
      '1': 'replicationstatus',
      '3': 529093900,
      '4': 1,
      '5': 14,
      '6': '.s3.ReplicationStatus',
      '10': 'replicationstatus'
    },
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {'1': 'restore', '3': 267943794, '4': 1, '5': 9, '10': 'restore'},
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'tagcount', '3': 339592595, '4': 1, '5': 5, '10': 'tagcount'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
    {
      '1': 'websiteredirectlocation',
      '3': 71844662,
      '4': 1,
      '5': 9,
      '10': 'websiteredirectlocation'
    },
  ],
  '3': [GetObjectOutput_MetadataEntry$json],
};

@$core.Deprecated('Use getObjectOutputDescriptor instead')
const GetObjectOutput_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetObjectOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectOutputDescriptor = $convert.base64Decode(
    'Cg9HZXRPYmplY3RPdXRwdXQSJgoMYWNjZXB0cmFuZ2VzGKCbxt0BIAEoCVIMYWNjZXB0cmFuZ2'
    'VzEhsKBGJvZHkYlqGoFCABKAxCBIi1GAFSBGJvZHkSLgoQYnVja2V0a2V5ZW5hYmxlZBjsoozP'
    'ASABKAhSEGJ1Y2tldGtleWVuYWJsZWQSJgoMY2FjaGVjb250cm9sGP+P5YkBIAEoCVIMY2FjaG'
    'Vjb250cm9sEicKDWNoZWNrc3VtY3JjMzIYwqPNMyABKAlSDWNoZWNrc3VtY3JjMzISKQoOY2hl'
    'Y2tzdW1jcmMzMmMYp5+lTCABKAlSDmNoZWNrc3VtY3JjMzJjEi8KEWNoZWNrc3VtY3JjNjRudm'
    '1lGM28izggASgJUhFjaGVja3N1bWNyYzY0bnZtZRImCgxjaGVja3N1bXNoYTEYxOzgigEgASgJ'
    'UgxjaGVja3N1bXNoYTESKQoOY2hlY2tzdW1zaGEyNTYYvvncRCABKAlSDmNoZWNrc3Vtc2hhMj'
    'U2EjcKDGNoZWNrc3VtdHlwZRjDvtkuIAEoDjIQLnMzLkNoZWNrc3VtVHlwZVIMY2hlY2tzdW10'
    'eXBlEjEKEmNvbnRlbnRkaXNwb3NpdGlvbhjC1Z45IAEoCVISY29udGVudGRpc3Bvc2l0aW9uEi'
    'wKD2NvbnRlbnRlbmNvZGluZxi00JqXASABKAlSD2NvbnRlbnRlbmNvZGluZxIrCg9jb250ZW50'
    'bGFuZ3VhZ2UYkbjdMyABKAlSD2NvbnRlbnRsYW5ndWFnZRInCg1jb250ZW50bGVuZ3RoGNeyw2'
    'wgASgDUg1jb250ZW50bGVuZ3RoEiUKDGNvbnRlbnRyYW5nZRjQ66QFIAEoCVIMY29udGVudHJh'
    'bmdlEiQKC2NvbnRlbnR0eXBlGJPV6J4BIAEoCVILY29udGVudHR5cGUSJQoMZGVsZXRlbWFya2'
    'VyGIGAzgIgASgIUgxkZWxldGVtYXJrZXISFgoEZXRhZxiB37OVASABKAlSBGV0YWcSIQoKZXhw'
    'aXJhdGlvbhiJqZ91IAEoCVIKZXhwaXJhdGlvbhIbCgdleHBpcmVzGKSKqD0gASgJUgdleHBpcm'
    'VzEiYKDGxhc3Rtb2RpZmllZBinnPzOASABKAlSDGxhc3Rtb2RpZmllZBJBCghtZXRhZGF0YRjh'
    '4o/gASADKAsyIS5zMy5HZXRPYmplY3RPdXRwdXQuTWV0YWRhdGFFbnRyeVIIbWV0YWRhdGESIw'
    'oLbWlzc2luZ21ldGEYq63eJSABKAVSC21pc3NpbmdtZXRhEl8KGW9iamVjdGxvY2tsZWdhbGhv'
    'bGRzdGF0dXMYtpLt/wEgASgOMh0uczMuT2JqZWN0TG9ja0xlZ2FsSG9sZFN0YXR1c1IZb2JqZW'
    'N0bG9ja2xlZ2FsaG9sZHN0YXR1cxI9Cg5vYmplY3Rsb2NrbW9kZRijnJ9aIAEoDjISLnMzLk9i'
    'amVjdExvY2tNb2RlUg5vYmplY3Rsb2NrbW9kZRI/ChlvYmplY3Rsb2NrcmV0YWludW50aWxkYX'
    'RlGLn4lH4gASgJUhlvYmplY3Rsb2NrcmV0YWludW50aWxkYXRlEiEKCnBhcnRzY291bnQYlZ30'
    'SSABKAVSCnBhcnRzY291bnQSRwoRcmVwbGljYXRpb25zdGF0dXMYjKql/AEgASgOMhUuczMuUm'
    'VwbGljYXRpb25TdGF0dXNSEXJlcGxpY2F0aW9uc3RhdHVzEj4KDnJlcXVlc3RjaGFyZ2VkGJPQ'
    'q7kBIAEoDjISLnMzLlJlcXVlc3RDaGFyZ2VkUg5yZXF1ZXN0Y2hhcmdlZBIbCgdyZXN0b3JlGP'
    'L+4X8gASgJUgdyZXN0b3JlEjUKFHNzZWN1c3RvbWVyYWxnb3JpdGhtGNDJgSsgASgJUhRzc2Vj'
    'dXN0b21lcmFsZ29yaXRobRIuChFzc2VjdXN0b21lcmtleW1kNRjo0RcgASgJUhFzc2VjdXN0b2'
    '1lcmtleW1kNRIkCgtzc2VrbXNrZXlpZBjYiNTUASABKAlSC3NzZWttc2tleWlkEk8KFHNlcnZl'
    'cnNpZGVlbmNyeXB0aW9uGLGOnwQgASgOMhguczMuU2VydmVyU2lkZUVuY3J5cHRpb25SFHNlcn'
    'ZlcnNpZGVlbmNyeXB0aW9uEjgKDHN0b3JhZ2VjbGFzcxjHiMS7ASABKA4yEC5zMy5TdG9yYWdl'
    'Q2xhc3NSDHN0b3JhZ2VjbGFzcxIeCgh0YWdjb3VudBiTi/ehASABKAVSCHRhZ2NvdW50EiAKCX'
    'ZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZBI7Chd3ZWJzaXRlcmVkaXJlY3Rsb2NhdGlv'
    'bhi2hqEiIAEoCVIXd2Vic2l0ZXJlZGlyZWN0bG9jYXRpb24aOwoNTWV0YWRhdGFFbnRyeRIQCg'
    'NrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use getObjectRequestDescriptor instead')
const GetObjectRequest$json = {
  '1': 'GetObjectRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksummode',
      '3': 350288954,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumMode',
      '10': 'checksummode'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'ifmodifiedsince',
      '3': 376471190,
      '4': 1,
      '5': 9,
      '10': 'ifmodifiedsince'
    },
    {'1': 'ifnonematch', '3': 231211830, '4': 1, '5': 9, '10': 'ifnonematch'},
    {
      '1': 'ifunmodifiedsince',
      '3': 311345013,
      '4': 1,
      '5': 9,
      '10': 'ifunmodifiedsince'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'partnumber', '3': 372082310, '4': 1, '5': 5, '10': 'partnumber'},
    {'1': 'range', '3': 51505011, '4': 1, '5': 9, '10': 'range'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'responsecachecontrol',
      '3': 150103426,
      '4': 1,
      '5': 9,
      '10': 'responsecachecontrol'
    },
    {
      '1': 'responsecontentdisposition',
      '3': 394384183,
      '4': 1,
      '5': 9,
      '10': 'responsecontentdisposition'
    },
    {
      '1': 'responsecontentencoding',
      '3': 67832847,
      '4': 1,
      '5': 9,
      '10': 'responsecontentencoding'
    },
    {
      '1': 'responsecontentlanguage',
      '3': 398143842,
      '4': 1,
      '5': 9,
      '10': 'responsecontentlanguage'
    },
    {
      '1': 'responsecontenttype',
      '3': 388099976,
      '4': 1,
      '5': 9,
      '10': 'responsecontenttype'
    },
    {
      '1': 'responseexpires',
      '3': 236546855,
      '4': 1,
      '5': 9,
      '10': 'responseexpires'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetObjectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectRequestDescriptor = $convert.base64Decode(
    'ChBHZXRPYmplY3RSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EjgKDGNoZWNrc3'
    'VtbW9kZRi6+IOnASABKA4yEC5zMy5DaGVja3N1bU1vZGVSDGNoZWNrc3VtbW9kZRIzChNleHBl'
    'Y3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhsKB2lmbWF0Y2'
    'gY0Ja3LCABKAlSB2lmbWF0Y2gSLAoPaWZtb2RpZmllZHNpbmNlGJb9wbMBIAEoCVIPaWZtb2Rp'
    'ZmllZHNpbmNlEiMKC2lmbm9uZW1hdGNoGLaGoG4gASgJUgtpZm5vbmVtYXRjaBIwChFpZnVubW'
    '9kaWZpZWRzaW5jZRj1/rqUASABKAlSEWlmdW5tb2RpZmllZHNpbmNlEhMKA2tleRiNkutoIAEo'
    'CVIDa2V5EiIKCnBhcnRudW1iZXIYho22sQEgASgFUgpwYXJ0bnVtYmVyEhcKBXJhbmdlGPPOxx'
    'ggASgJUgVyYW5nZRI4CgxyZXF1ZXN0cGF5ZXIYpObh9QEgASgOMhAuczMuUmVxdWVzdFBheWVy'
    'UgxyZXF1ZXN0cGF5ZXISNQoUcmVzcG9uc2VjYWNoZWNvbnRyb2wYgsvJRyABKAlSFHJlc3Bvbn'
    'NlY2FjaGVjb250cm9sEkIKGnJlc3BvbnNlY29udGVudGRpc3Bvc2l0aW9uGLemh7wBIAEoCVIa'
    'cmVzcG9uc2Vjb250ZW50ZGlzcG9zaXRpb24SOwoXcmVzcG9uc2Vjb250ZW50ZW5jb2RpbmcYj5'
    'isICABKAlSF3Jlc3BvbnNlY29udGVudGVuY29kaW5nEjwKF3Jlc3BvbnNlY29udGVudGxhbmd1'
    'YWdlGOLi7L0BIAEoCVIXcmVzcG9uc2Vjb250ZW50bGFuZ3VhZ2USNAoTcmVzcG9uc2Vjb250ZW'
    '50dHlwZRiI34e5ASABKAlSE3Jlc3BvbnNlY29udGVudHR5cGUSKwoPcmVzcG9uc2VleHBpcmVz'
    'GKfW5XAgASgJUg9yZXNwb25zZWV4cGlyZXMSNQoUc3NlY3VzdG9tZXJhbGdvcml0aG0Y0MmBKy'
    'ABKAlSFHNzZWN1c3RvbWVyYWxnb3JpdGhtEikKDnNzZWN1c3RvbWVya2V5GJr+9DsgASgJUg5z'
    'c2VjdXN0b21lcmtleRIuChFzc2VjdXN0b21lcmtleW1kNRjo0RcgASgJUhFzc2VjdXN0b21lcm'
    'tleW1kNRIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use getObjectRetentionOutputDescriptor instead')
const GetObjectRetentionOutput$json = {
  '1': 'GetObjectRetentionOutput',
  '2': [
    {
      '1': 'retention',
      '3': 299600946,
      '4': 1,
      '5': 11,
      '6': '.s3.ObjectLockRetention',
      '8': {},
      '10': 'retention'
    },
  ],
};

/// Descriptor for `GetObjectRetentionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectRetentionOutputDescriptor =
    $convert.base64Decode(
        'ChhHZXRPYmplY3RSZXRlbnRpb25PdXRwdXQSPwoJcmV0ZW50aW9uGLKY7o4BIAEoCzIXLnMzLk'
        '9iamVjdExvY2tSZXRlbnRpb25CBIi1GAFSCXJldGVudGlvbg==');

@$core.Deprecated('Use getObjectRetentionRequestDescriptor instead')
const GetObjectRetentionRequest$json = {
  '1': 'GetObjectRetentionRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetObjectRetentionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectRetentionRequestDescriptor = $convert.base64Decode(
    'ChlHZXRPYmplY3RSZXRlbnRpb25SZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0Ej'
    'MKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXISEwoD'
    'a2V5GI2S62ggASgJUgNrZXkSOAoMcmVxdWVzdHBheWVyGKTm4fUBIAEoDjIQLnMzLlJlcXVlc3'
    'RQYXllclIMcmVxdWVzdHBheWVyEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated('Use getObjectTaggingOutputDescriptor instead')
const GetObjectTaggingOutput$json = {
  '1': 'GetObjectTaggingOutput',
  '2': [
    {
      '1': 'tagset',
      '3': 454361330,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tagset'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetObjectTaggingOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectTaggingOutputDescriptor =
    $convert.base64Decode(
        'ChZHZXRPYmplY3RUYWdnaW5nT3V0cHV0EiMKBnRhZ3NldBjygdTYASADKAsyBy5zMy5UYWdSBn'
        'RhZ3NldBIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use getObjectTaggingRequestDescriptor instead')
const GetObjectTaggingRequest$json = {
  '1': 'GetObjectTaggingRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `GetObjectTaggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectTaggingRequestDescriptor = $convert.base64Decode(
    'ChdHZXRPYmplY3RUYWdnaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzCh'
    'NleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhMKA2tl'
    'eRiNkutoIAEoCVIDa2V5EjgKDHJlcXVlc3RwYXllchik5uH1ASABKA4yEC5zMy5SZXF1ZXN0UG'
    'F5ZXJSDHJlcXVlc3RwYXllchIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use getObjectTorrentOutputDescriptor instead')
const GetObjectTorrentOutput$json = {
  '1': 'GetObjectTorrentOutput',
  '2': [
    {'1': 'body', '3': 42602646, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `GetObjectTorrentOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectTorrentOutputDescriptor = $convert.base64Decode(
    'ChZHZXRPYmplY3RUb3JyZW50T3V0cHV0EhsKBGJvZHkYlqGoFCABKAxCBIi1GAFSBGJvZHkSPg'
    'oOcmVxdWVzdGNoYXJnZWQYk9CruQEgASgOMhIuczMuUmVxdWVzdENoYXJnZWRSDnJlcXVlc3Rj'
    'aGFyZ2Vk');

@$core.Deprecated('Use getObjectTorrentRequestDescriptor instead')
const GetObjectTorrentRequest$json = {
  '1': 'GetObjectTorrentRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
  ],
};

/// Descriptor for `GetObjectTorrentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getObjectTorrentRequestDescriptor = $convert.base64Decode(
    'ChdHZXRPYmplY3RUb3JyZW50UmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzCh'
    'NleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhMKA2tl'
    'eRiNkutoIAEoCVIDa2V5EjgKDHJlcXVlc3RwYXllchik5uH1ASABKA4yEC5zMy5SZXF1ZXN0UG'
    'F5ZXJSDHJlcXVlc3RwYXllcg==');

@$core.Deprecated('Use getPublicAccessBlockOutputDescriptor instead')
const GetPublicAccessBlockOutput$json = {
  '1': 'GetPublicAccessBlockOutput',
  '2': [
    {
      '1': 'publicaccessblockconfiguration',
      '3': 136498568,
      '4': 1,
      '5': 11,
      '6': '.s3.PublicAccessBlockConfiguration',
      '8': {},
      '10': 'publicaccessblockconfiguration'
    },
  ],
};

/// Descriptor for `GetPublicAccessBlockOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicAccessBlockOutputDescriptor =
    $convert.base64Decode(
        'ChpHZXRQdWJsaWNBY2Nlc3NCbG9ja091dHB1dBJzCh5wdWJsaWNhY2Nlc3NibG9ja2NvbmZpZ3'
        'VyYXRpb24YiJuLQSABKAsyIi5zMy5QdWJsaWNBY2Nlc3NCbG9ja0NvbmZpZ3VyYXRpb25CBIi1'
        'GAFSHnB1YmxpY2FjY2Vzc2Jsb2NrY29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use getPublicAccessBlockRequestDescriptor instead')
const GetPublicAccessBlockRequest$json = {
  '1': 'GetPublicAccessBlockRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `GetPublicAccessBlockRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicAccessBlockRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRQdWJsaWNBY2Nlc3NCbG9ja1JlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZX'
        'QSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lcg==');

@$core.Deprecated('Use glacierJobParametersDescriptor instead')
const GlacierJobParameters$json = {
  '1': 'GlacierJobParameters',
  '2': [
    {
      '1': 'tier',
      '3': 519596586,
      '4': 1,
      '5': 14,
      '6': '.s3.Tier',
      '10': 'tier'
    },
  ],
};

/// Descriptor for `GlacierJobParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List glacierJobParametersDescriptor = $convert.base64Decode(
    'ChRHbGFjaWVySm9iUGFyYW1ldGVycxIgCgR0aWVyGKrU4fcBIAEoDjIILnMzLlRpZXJSBHRpZX'
    'I=');

@$core.Deprecated('Use grantDescriptor instead')
const Grant$json = {
  '1': 'Grant',
  '2': [
    {
      '1': 'grantee',
      '3': 89454056,
      '4': 1,
      '5': 11,
      '6': '.s3.Grantee',
      '10': 'grantee'
    },
    {
      '1': 'permission',
      '3': 465848867,
      '4': 1,
      '5': 14,
      '6': '.s3.Permission',
      '10': 'permission'
    },
  ],
};

/// Descriptor for `Grant`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grantDescriptor = $convert.base64Decode(
    'CgVHcmFudBIoCgdncmFudGVlGOjr0yogASgLMgsuczMuR3JhbnRlZVIHZ3JhbnRlZRIyCgpwZX'
    'JtaXNzaW9uGKOUkd4BIAEoDjIOLnMzLlBlcm1pc3Npb25SCnBlcm1pc3Npb24=');

@$core.Deprecated('Use granteeDescriptor instead')
const Grantee$json = {
  '1': 'Grantee',
  '2': [
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {'1': 'id', '3': 384363361, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.s3.Type',
      '10': 'type'
    },
    {'1': 'uri', '3': 443116318, '4': 1, '5': 9, '10': 'uri'},
  ],
};

/// Descriptor for `Grantee`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List granteeDescriptor = $convert.base64Decode(
    'CgdHcmFudGVlEiQKC2Rpc3BsYXluYW1lGLfJsscBIAEoCVILZGlzcGxheW5hbWUSJgoMZW1haW'
    'xhZGRyZXNzGNbxiukBIAEoCVIMZW1haWxhZGRyZXNzEhIKAmlkGOHWo7cBIAEoCVICaWQSIAoE'
    'dHlwZRjuoNeKASABKA4yCC5zMy5UeXBlUgR0eXBlEhQKA3VyaRie1qXTASABKAlSA3VyaQ==');

@$core.Deprecated('Use headBucketOutputDescriptor instead')
const HeadBucketOutput$json = {
  '1': 'HeadBucketOutput',
  '2': [
    {
      '1': 'accesspointalias',
      '3': 360968558,
      '4': 1,
      '5': 8,
      '10': 'accesspointalias'
    },
    {'1': 'bucketarn', '3': 255683899, '4': 1, '5': 9, '10': 'bucketarn'},
    {
      '1': 'bucketlocationname',
      '3': 126439780,
      '4': 1,
      '5': 9,
      '10': 'bucketlocationname'
    },
    {
      '1': 'bucketlocationtype',
      '3': 306578505,
      '4': 1,
      '5': 14,
      '6': '.s3.LocationType',
      '10': 'bucketlocationtype'
    },
    {'1': 'bucketregion', '3': 309298816, '4': 1, '5': 9, '10': 'bucketregion'},
  ],
};

/// Descriptor for `HeadBucketOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List headBucketOutputDescriptor = $convert.base64Decode(
    'ChBIZWFkQnVja2V0T3V0cHV0Ei4KEGFjY2Vzc3BvaW50YWxpYXMY7uKPrAEgASgIUhBhY2Nlc3'
    'Nwb2ludGFsaWFzEh8KCWJ1Y2tldGFybhi72vV5IAEoCVIJYnVja2V0YXJuEjEKEmJ1Y2tldGxv'
    'Y2F0aW9ubmFtZRjkoqU8IAEoCVISYnVja2V0bG9jYXRpb25uYW1lEkQKEmJ1Y2tldGxvY2F0aW'
    '9udHlwZRjJiJiSASABKA4yEC5zMy5Mb2NhdGlvblR5cGVSEmJ1Y2tldGxvY2F0aW9udHlwZRIm'
    'CgxidWNrZXRyZWdpb24YgI2+kwEgASgJUgxidWNrZXRyZWdpb24=');

@$core.Deprecated('Use headBucketRequestDescriptor instead')
const HeadBucketRequest$json = {
  '1': 'HeadBucketRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `HeadBucketRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List headBucketRequestDescriptor = $convert.base64Decode(
    'ChFIZWFkQnVja2V0UmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzChNleHBlY3'
    'RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVy');

@$core.Deprecated('Use headObjectOutputDescriptor instead')
const HeadObjectOutput$json = {
  '1': 'HeadObjectOutput',
  '2': [
    {'1': 'acceptranges', '3': 464620960, '4': 1, '5': 9, '10': 'acceptranges'},
    {
      '1': 'archivestatus',
      '3': 533231984,
      '4': 1,
      '5': 14,
      '6': '.s3.ArchiveStatus',
      '10': 'archivestatus'
    },
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {'1': 'cachecontrol', '3': 288966655, '4': 1, '5': 9, '10': 'cachecontrol'},
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {
      '1': 'contentdisposition',
      '3': 120040130,
      '4': 1,
      '5': 9,
      '10': 'contentdisposition'
    },
    {
      '1': 'contentencoding',
      '3': 317106228,
      '4': 1,
      '5': 9,
      '10': 'contentencoding'
    },
    {
      '1': 'contentlanguage',
      '3': 108485649,
      '4': 1,
      '5': 9,
      '10': 'contentlanguage'
    },
    {
      '1': 'contentlength',
      '3': 227596631,
      '4': 1,
      '5': 3,
      '10': 'contentlength'
    },
    {'1': 'contentrange', '3': 11089360, '4': 1, '5': 9, '10': 'contentrange'},
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'deletemarker', '3': 5472257, '4': 1, '5': 8, '10': 'deletemarker'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {'1': 'expires', '3': 128582948, '4': 1, '5': 9, '10': 'expires'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 3,
      '5': 11,
      '6': '.s3.HeadObjectOutput.MetadataEntry',
      '10': 'metadata'
    },
    {'1': 'missingmeta', '3': 79140523, '4': 1, '5': 5, '10': 'missingmeta'},
    {
      '1': 'objectlocklegalholdstatus',
      '3': 536561974,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockLegalHoldStatus',
      '10': 'objectlocklegalholdstatus'
    },
    {
      '1': 'objectlockmode',
      '3': 189255203,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockMode',
      '10': 'objectlockmode'
    },
    {
      '1': 'objectlockretainuntildate',
      '3': 264584249,
      '4': 1,
      '5': 9,
      '10': 'objectlockretainuntildate'
    },
    {'1': 'partscount', '3': 154996373, '4': 1, '5': 5, '10': 'partscount'},
    {
      '1': 'replicationstatus',
      '3': 529093900,
      '4': 1,
      '5': 14,
      '6': '.s3.ReplicationStatus',
      '10': 'replicationstatus'
    },
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {'1': 'restore', '3': 267943794, '4': 1, '5': 9, '10': 'restore'},
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'tagcount', '3': 339592595, '4': 1, '5': 5, '10': 'tagcount'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
    {
      '1': 'websiteredirectlocation',
      '3': 71844662,
      '4': 1,
      '5': 9,
      '10': 'websiteredirectlocation'
    },
  ],
  '3': [HeadObjectOutput_MetadataEntry$json],
};

@$core.Deprecated('Use headObjectOutputDescriptor instead')
const HeadObjectOutput_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `HeadObjectOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List headObjectOutputDescriptor = $convert.base64Decode(
    'ChBIZWFkT2JqZWN0T3V0cHV0EiYKDGFjY2VwdHJhbmdlcxigm8bdASABKAlSDGFjY2VwdHJhbm'
    'dlcxI7Cg1hcmNoaXZlc3RhdHVzGPDyof4BIAEoDjIRLnMzLkFyY2hpdmVTdGF0dXNSDWFyY2hp'
    'dmVzdGF0dXMSLgoQYnVja2V0a2V5ZW5hYmxlZBjsoozPASABKAhSEGJ1Y2tldGtleWVuYWJsZW'
    'QSJgoMY2FjaGVjb250cm9sGP+P5YkBIAEoCVIMY2FjaGVjb250cm9sEicKDWNoZWNrc3VtY3Jj'
    'MzIYwqPNMyABKAlSDWNoZWNrc3VtY3JjMzISKQoOY2hlY2tzdW1jcmMzMmMYp5+lTCABKAlSDm'
    'NoZWNrc3VtY3JjMzJjEi8KEWNoZWNrc3VtY3JjNjRudm1lGM28izggASgJUhFjaGVja3N1bWNy'
    'YzY0bnZtZRImCgxjaGVja3N1bXNoYTEYxOzgigEgASgJUgxjaGVja3N1bXNoYTESKQoOY2hlY2'
    'tzdW1zaGEyNTYYvvncRCABKAlSDmNoZWNrc3Vtc2hhMjU2EjcKDGNoZWNrc3VtdHlwZRjDvtku'
    'IAEoDjIQLnMzLkNoZWNrc3VtVHlwZVIMY2hlY2tzdW10eXBlEjEKEmNvbnRlbnRkaXNwb3NpdG'
    'lvbhjC1Z45IAEoCVISY29udGVudGRpc3Bvc2l0aW9uEiwKD2NvbnRlbnRlbmNvZGluZxi00JqX'
    'ASABKAlSD2NvbnRlbnRlbmNvZGluZxIrCg9jb250ZW50bGFuZ3VhZ2UYkbjdMyABKAlSD2Nvbn'
    'RlbnRsYW5ndWFnZRInCg1jb250ZW50bGVuZ3RoGNeyw2wgASgDUg1jb250ZW50bGVuZ3RoEiUK'
    'DGNvbnRlbnRyYW5nZRjQ66QFIAEoCVIMY29udGVudHJhbmdlEiQKC2NvbnRlbnR0eXBlGJPV6J'
    '4BIAEoCVILY29udGVudHR5cGUSJQoMZGVsZXRlbWFya2VyGIGAzgIgASgIUgxkZWxldGVtYXJr'
    'ZXISFgoEZXRhZxiB37OVASABKAlSBGV0YWcSIQoKZXhwaXJhdGlvbhiJqZ91IAEoCVIKZXhwaX'
    'JhdGlvbhIbCgdleHBpcmVzGKSKqD0gASgJUgdleHBpcmVzEiYKDGxhc3Rtb2RpZmllZBinnPzO'
    'ASABKAlSDGxhc3Rtb2RpZmllZBJCCghtZXRhZGF0YRjh4o/gASADKAsyIi5zMy5IZWFkT2JqZW'
    'N0T3V0cHV0Lk1ldGFkYXRhRW50cnlSCG1ldGFkYXRhEiMKC21pc3NpbmdtZXRhGKut3iUgASgF'
    'UgttaXNzaW5nbWV0YRJfChlvYmplY3Rsb2NrbGVnYWxob2xkc3RhdHVzGLaS7f8BIAEoDjIdLn'
    'MzLk9iamVjdExvY2tMZWdhbEhvbGRTdGF0dXNSGW9iamVjdGxvY2tsZWdhbGhvbGRzdGF0dXMS'
    'PQoOb2JqZWN0bG9ja21vZGUYo5yfWiABKA4yEi5zMy5PYmplY3RMb2NrTW9kZVIOb2JqZWN0bG'
    '9ja21vZGUSPwoZb2JqZWN0bG9ja3JldGFpbnVudGlsZGF0ZRi5+JR+IAEoCVIZb2JqZWN0bG9j'
    'a3JldGFpbnVudGlsZGF0ZRIhCgpwYXJ0c2NvdW50GJWd9EkgASgFUgpwYXJ0c2NvdW50EkcKEX'
    'JlcGxpY2F0aW9uc3RhdHVzGIyqpfwBIAEoDjIVLnMzLlJlcGxpY2F0aW9uU3RhdHVzUhFyZXBs'
    'aWNhdGlvbnN0YXR1cxI+Cg5yZXF1ZXN0Y2hhcmdlZBiT0Ku5ASABKA4yEi5zMy5SZXF1ZXN0Q2'
    'hhcmdlZFIOcmVxdWVzdGNoYXJnZWQSGwoHcmVzdG9yZRjy/uF/IAEoCVIHcmVzdG9yZRI1ChRz'
    'c2VjdXN0b21lcmFsZ29yaXRobRjQyYErIAEoCVIUc3NlY3VzdG9tZXJhbGdvcml0aG0SLgoRc3'
    'NlY3VzdG9tZXJrZXltZDUY6NEXIAEoCVIRc3NlY3VzdG9tZXJrZXltZDUSJAoLc3Nla21za2V5'
    'aWQY2IjU1AEgASgJUgtzc2VrbXNrZXlpZBJPChRzZXJ2ZXJzaWRlZW5jcnlwdGlvbhixjp8EIA'
    'EoDjIYLnMzLlNlcnZlclNpZGVFbmNyeXB0aW9uUhRzZXJ2ZXJzaWRlZW5jcnlwdGlvbhI4Cgxz'
    'dG9yYWdlY2xhc3MYx4jEuwEgASgOMhAuczMuU3RvcmFnZUNsYXNzUgxzdG9yYWdlY2xhc3MSHg'
    'oIdGFnY291bnQYk4v3oQEgASgFUgh0YWdjb3VudBIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2'
    'ZXJzaW9uaWQSOwoXd2Vic2l0ZXJlZGlyZWN0bG9jYXRpb24YtoahIiABKAlSF3dlYnNpdGVyZW'
    'RpcmVjdGxvY2F0aW9uGjsKDU1ldGFkYXRhRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFs'
    'dWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use headObjectRequestDescriptor instead')
const HeadObjectRequest$json = {
  '1': 'HeadObjectRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksummode',
      '3': 350288954,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumMode',
      '10': 'checksummode'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {
      '1': 'ifmodifiedsince',
      '3': 376471190,
      '4': 1,
      '5': 9,
      '10': 'ifmodifiedsince'
    },
    {'1': 'ifnonematch', '3': 231211830, '4': 1, '5': 9, '10': 'ifnonematch'},
    {
      '1': 'ifunmodifiedsince',
      '3': 311345013,
      '4': 1,
      '5': 9,
      '10': 'ifunmodifiedsince'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'partnumber', '3': 372082310, '4': 1, '5': 5, '10': 'partnumber'},
    {'1': 'range', '3': 51505011, '4': 1, '5': 9, '10': 'range'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'responsecachecontrol',
      '3': 150103426,
      '4': 1,
      '5': 9,
      '10': 'responsecachecontrol'
    },
    {
      '1': 'responsecontentdisposition',
      '3': 394384183,
      '4': 1,
      '5': 9,
      '10': 'responsecontentdisposition'
    },
    {
      '1': 'responsecontentencoding',
      '3': 67832847,
      '4': 1,
      '5': 9,
      '10': 'responsecontentencoding'
    },
    {
      '1': 'responsecontentlanguage',
      '3': 398143842,
      '4': 1,
      '5': 9,
      '10': 'responsecontentlanguage'
    },
    {
      '1': 'responsecontenttype',
      '3': 388099976,
      '4': 1,
      '5': 9,
      '10': 'responsecontenttype'
    },
    {
      '1': 'responseexpires',
      '3': 236546855,
      '4': 1,
      '5': 9,
      '10': 'responseexpires'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `HeadObjectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List headObjectRequestDescriptor = $convert.base64Decode(
    'ChFIZWFkT2JqZWN0UmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBI4CgxjaGVja3'
    'N1bW1vZGUYuviDpwEgASgOMhAuczMuQ2hlY2tzdW1Nb2RlUgxjaGVja3N1bW1vZGUSMwoTZXhw'
    'ZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchIbCgdpZm1hdG'
    'NoGNCWtywgASgJUgdpZm1hdGNoEiwKD2lmbW9kaWZpZWRzaW5jZRiW/cGzASABKAlSD2lmbW9k'
    'aWZpZWRzaW5jZRIjCgtpZm5vbmVtYXRjaBi2hqBuIAEoCVILaWZub25lbWF0Y2gSMAoRaWZ1bm'
    '1vZGlmaWVkc2luY2UY9f66lAEgASgJUhFpZnVubW9kaWZpZWRzaW5jZRITCgNrZXkYjZLraCAB'
    'KAlSA2tleRIiCgpwYXJ0bnVtYmVyGIaNtrEBIAEoBVIKcGFydG51bWJlchIXCgVyYW5nZRjzzs'
    'cYIAEoCVIFcmFuZ2USOAoMcmVxdWVzdHBheWVyGKTm4fUBIAEoDjIQLnMzLlJlcXVlc3RQYXll'
    'clIMcmVxdWVzdHBheWVyEjUKFHJlc3BvbnNlY2FjaGVjb250cm9sGILLyUcgASgJUhRyZXNwb2'
    '5zZWNhY2hlY29udHJvbBJCChpyZXNwb25zZWNvbnRlbnRkaXNwb3NpdGlvbhi3poe8ASABKAlS'
    'GnJlc3BvbnNlY29udGVudGRpc3Bvc2l0aW9uEjsKF3Jlc3BvbnNlY29udGVudGVuY29kaW5nGI'
    '+YrCAgASgJUhdyZXNwb25zZWNvbnRlbnRlbmNvZGluZxI8ChdyZXNwb25zZWNvbnRlbnRsYW5n'
    'dWFnZRji4uy9ASABKAlSF3Jlc3BvbnNlY29udGVudGxhbmd1YWdlEjQKE3Jlc3BvbnNlY29udG'
    'VudHR5cGUYiN+HuQEgASgJUhNyZXNwb25zZWNvbnRlbnR0eXBlEisKD3Jlc3BvbnNlZXhwaXJl'
    'cxin1uVwIAEoCVIPcmVzcG9uc2VleHBpcmVzEjUKFHNzZWN1c3RvbWVyYWxnb3JpdGhtGNDJgS'
    'sgASgJUhRzc2VjdXN0b21lcmFsZ29yaXRobRIpCg5zc2VjdXN0b21lcmtleRia/vQ7IAEoCVIO'
    'c3NlY3VzdG9tZXJrZXkSLgoRc3NlY3VzdG9tZXJrZXltZDUY6NEXIAEoCVIRc3NlY3VzdG9tZX'
    'JrZXltZDUSIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlk');

@$core.Deprecated('Use idempotencyParameterMismatchDescriptor instead')
const IdempotencyParameterMismatch$json = {
  '1': 'IdempotencyParameterMismatch',
};

/// Descriptor for `IdempotencyParameterMismatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List idempotencyParameterMismatchDescriptor =
    $convert.base64Decode('ChxJZGVtcG90ZW5jeVBhcmFtZXRlck1pc21hdGNo');

@$core.Deprecated('Use indexDocumentDescriptor instead')
const IndexDocument$json = {
  '1': 'IndexDocument',
  '2': [
    {'1': 'suffix', '3': 287160941, '4': 1, '5': 9, '10': 'suffix'},
  ],
};

/// Descriptor for `IndexDocument`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List indexDocumentDescriptor = $convert.base64Decode(
    'Cg1JbmRleERvY3VtZW50EhoKBnN1ZmZpeBjt9PaIASABKAlSBnN1ZmZpeA==');

@$core.Deprecated('Use initiatorDescriptor instead')
const Initiator$json = {
  '1': 'Initiator',
  '2': [
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {'1': 'id', '3': 384363361, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `Initiator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List initiatorDescriptor = $convert.base64Decode(
    'CglJbml0aWF0b3ISJAoLZGlzcGxheW5hbWUYt8myxwEgASgJUgtkaXNwbGF5bmFtZRISCgJpZB'
    'jh1qO3ASABKAlSAmlk');

@$core.Deprecated('Use inputSerializationDescriptor instead')
const InputSerialization$json = {
  '1': 'InputSerialization',
  '2': [
    {
      '1': 'csv',
      '3': 455189056,
      '4': 1,
      '5': 11,
      '6': '.s3.CSVInput',
      '10': 'csv'
    },
    {
      '1': 'compressiontype',
      '3': 335198942,
      '4': 1,
      '5': 14,
      '6': '.s3.CompressionType',
      '10': 'compressiontype'
    },
    {
      '1': 'json',
      '3': 532512708,
      '4': 1,
      '5': 11,
      '6': '.s3.JSONInput',
      '10': 'json'
    },
    {
      '1': 'parquet',
      '3': 308729196,
      '4': 1,
      '5': 11,
      '6': '.s3.ParquetInput',
      '10': 'parquet'
    },
  ],
};

/// Descriptor for `InputSerialization`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inputSerializationDescriptor = $convert.base64Decode(
    'ChJJbnB1dFNlcmlhbGl6YXRpb24SIgoDY3N2GMDEhtkBIAEoCzIMLnMzLkNTVklucHV0UgNjc3'
    'YSQQoPY29tcHJlc3Npb250eXBlGN716p8BIAEoDjITLnMzLkNvbXByZXNzaW9uVHlwZVIPY29t'
    'cHJlc3Npb250eXBlEiUKBGpzb24YxP/1/QEgASgLMg0uczMuSlNPTklucHV0UgRqc29uEi4KB3'
    'BhcnF1ZXQY7KqbkwEgASgLMhAuczMuUGFycXVldElucHV0UgdwYXJxdWV0');

@$core.Deprecated('Use intelligentTieringAndOperatorDescriptor instead')
const IntelligentTieringAndOperator$json = {
  '1': 'IntelligentTieringAndOperator',
  '2': [
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `IntelligentTieringAndOperator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List intelligentTieringAndOperatorDescriptor =
    $convert.base64Decode(
        'Ch1JbnRlbGxpZ2VudFRpZXJpbmdBbmRPcGVyYXRvchIaCgZwcmVmaXgY6rPTggEgASgJUgZwcm'
        'VmaXgSHwoEdGFncxjBwfa1ASADKAsyBy5zMy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use intelligentTieringConfigurationDescriptor instead')
const IntelligentTieringConfiguration$json = {
  '1': 'IntelligentTieringConfiguration',
  '2': [
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.IntelligentTieringFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.IntelligentTieringStatus',
      '10': 'status'
    },
    {
      '1': 'tierings',
      '3': 92603853,
      '4': 3,
      '5': 11,
      '6': '.s3.Tiering',
      '10': 'tierings'
    },
  ],
};

/// Descriptor for `IntelligentTieringConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List intelligentTieringConfigurationDescriptor =
    $convert.base64Decode(
        'Ch9JbnRlbGxpZ2VudFRpZXJpbmdDb25maWd1cmF0aW9uEjgKBmZpbHRlchiYgaelASABKAsyHC'
        '5zMy5JbnRlbGxpZ2VudFRpZXJpbmdGaWx0ZXJSBmZpbHRlchISCgJpZBiB8qK3ASABKAlSAmlk'
        'EjcKBnN0YXR1cxiQ5PsCIAEoDjIcLnMzLkludGVsbGlnZW50VGllcmluZ1N0YXR1c1IGc3RhdH'
        'VzEioKCHRpZXJpbmdzGM2LlCwgAygLMgsuczMuVGllcmluZ1IIdGllcmluZ3M=');

@$core.Deprecated('Use intelligentTieringFilterDescriptor instead')
const IntelligentTieringFilter$json = {
  '1': 'IntelligentTieringFilter',
  '2': [
    {
      '1': 'and',
      '3': 297135431,
      '4': 1,
      '5': 11,
      '6': '.s3.IntelligentTieringAndOperator',
      '10': 'and'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {'1': 'tag', '3': 411259956, '4': 1, '5': 11, '6': '.s3.Tag', '10': 'tag'},
  ],
};

/// Descriptor for `IntelligentTieringFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List intelligentTieringFilterDescriptor = $convert.base64Decode(
    'ChhJbnRlbGxpZ2VudFRpZXJpbmdGaWx0ZXISNwoDYW5kGMfa140BIAEoCzIhLnMzLkludGVsbG'
    'lnZW50VGllcmluZ0FuZE9wZXJhdG9yUgNhbmQSGgoGcHJlZml4GOqz04IBIAEoCVIGcHJlZml4'
    'Eh0KA3RhZxi0qI3EASABKAsyBy5zMy5UYWdSA3RhZw==');

@$core.Deprecated('Use invalidObjectStateDescriptor instead')
const InvalidObjectState$json = {
  '1': 'InvalidObjectState',
  '2': [
    {
      '1': 'accesstier',
      '3': 263116188,
      '4': 1,
      '5': 14,
      '6': '.s3.IntelligentTieringAccessTier',
      '10': 'accesstier'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
  ],
};

/// Descriptor for `InvalidObjectState`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidObjectStateDescriptor = $convert.base64Decode(
    'ChJJbnZhbGlkT2JqZWN0U3RhdGUSQwoKYWNjZXNzdGllchicq7t9IAEoDjIgLnMzLkludGVsbG'
    'lnZW50VGllcmluZ0FjY2Vzc1RpZXJSCmFjY2Vzc3RpZXISOAoMc3RvcmFnZWNsYXNzGMeIxLsB'
    'IAEoDjIQLnMzLlN0b3JhZ2VDbGFzc1IMc3RvcmFnZWNsYXNz');

@$core.Deprecated('Use invalidRequestDescriptor instead')
const InvalidRequest$json = {
  '1': 'InvalidRequest',
};

/// Descriptor for `InvalidRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRequestDescriptor =
    $convert.base64Decode('Cg5JbnZhbGlkUmVxdWVzdA==');

@$core.Deprecated('Use invalidWriteOffsetDescriptor instead')
const InvalidWriteOffset$json = {
  '1': 'InvalidWriteOffset',
};

/// Descriptor for `InvalidWriteOffset`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidWriteOffsetDescriptor =
    $convert.base64Decode('ChJJbnZhbGlkV3JpdGVPZmZzZXQ=');

@$core.Deprecated('Use inventoryConfigurationDescriptor instead')
const InventoryConfiguration$json = {
  '1': 'InventoryConfiguration',
  '2': [
    {
      '1': 'destination',
      '3': 457443680,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryDestination',
      '10': 'destination'
    },
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'includedobjectversions',
      '3': 487647084,
      '4': 1,
      '5': 14,
      '6': '.s3.InventoryIncludedObjectVersions',
      '10': 'includedobjectversions'
    },
    {'1': 'isenabled', '3': 319799247, '4': 1, '5': 8, '10': 'isenabled'},
    {
      '1': 'optionalfields',
      '3': 5014913,
      '4': 3,
      '5': 14,
      '6': '.s3.InventoryOptionalField',
      '10': 'optionalfields'
    },
    {
      '1': 'schedule',
      '3': 66697965,
      '4': 1,
      '5': 11,
      '6': '.s3.InventorySchedule',
      '10': 'schedule'
    },
  ],
};

/// Descriptor for `InventoryConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryConfigurationDescriptor = $convert.base64Decode(
    'ChZJbnZlbnRvcnlDb25maWd1cmF0aW9uEj4KC2Rlc3RpbmF0aW9uGOCSkNoBIAEoCzIYLnMzLk'
    'ludmVudG9yeURlc3RpbmF0aW9uUgtkZXN0aW5hdGlvbhIvCgZmaWx0ZXIYmIGnpQEgASgLMhMu'
    'czMuSW52ZW50b3J5RmlsdGVyUgZmaWx0ZXISEgoCaWQYgfKitwEgASgJUgJpZBJfChZpbmNsdW'
    'RlZG9iamVjdHZlcnNpb25zGOzOw+gBIAEoDjIjLnMzLkludmVudG9yeUluY2x1ZGVkT2JqZWN0'
    'VmVyc2lvbnNSFmluY2x1ZGVkb2JqZWN0dmVyc2lvbnMSIAoJaXNlbmFibGVkGM//vpgBIAEoCF'
    'IJaXNlbmFibGVkEkUKDm9wdGlvbmFsZmllbGRzGIGLsgIgAygOMhouczMuSW52ZW50b3J5T3B0'
    'aW9uYWxGaWVsZFIOb3B0aW9uYWxmaWVsZHMSNAoIc2NoZWR1bGUY7fXmHyABKAsyFS5zMy5Jbn'
    'ZlbnRvcnlTY2hlZHVsZVIIc2NoZWR1bGU=');

@$core.Deprecated('Use inventoryDestinationDescriptor instead')
const InventoryDestination$json = {
  '1': 'InventoryDestination',
  '2': [
    {
      '1': 's3bucketdestination',
      '3': 94381700,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryS3BucketDestination',
      '10': 's3bucketdestination'
    },
  ],
};

/// Descriptor for `InventoryDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryDestinationDescriptor = $convert.base64Decode(
    'ChRJbnZlbnRvcnlEZXN0aW5hdGlvbhJVChNzM2J1Y2tldGRlc3RpbmF0aW9uGITNgC0gASgLMi'
    'AuczMuSW52ZW50b3J5UzNCdWNrZXREZXN0aW5hdGlvblITczNidWNrZXRkZXN0aW5hdGlvbg==');

@$core.Deprecated('Use inventoryEncryptionDescriptor instead')
const InventoryEncryption$json = {
  '1': 'InventoryEncryption',
  '2': [
    {
      '1': 'ssekms',
      '3': 230338916,
      '4': 1,
      '5': 11,
      '6': '.s3.SSEKMS',
      '10': 'ssekms'
    },
    {
      '1': 'sses3',
      '3': 4241889,
      '4': 1,
      '5': 11,
      '6': '.s3.SSES3',
      '10': 'sses3'
    },
  ],
};

/// Descriptor for `InventoryEncryption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryEncryptionDescriptor = $convert.base64Decode(
    'ChNJbnZlbnRvcnlFbmNyeXB0aW9uEiUKBnNzZWttcxjk4uptIAEoCzIKLnMzLlNTRUtNU1IGc3'
    'Nla21zEiIKBXNzZXMzGOHzggIgASgLMgkuczMuU1NFUzNSBXNzZXMz');

@$core.Deprecated('Use inventoryFilterDescriptor instead')
const InventoryFilter$json = {
  '1': 'InventoryFilter',
  '2': [
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
  ],
};

/// Descriptor for `InventoryFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryFilterDescriptor = $convert.base64Decode(
    'Cg9JbnZlbnRvcnlGaWx0ZXISGgoGcHJlZml4GOqz04IBIAEoCVIGcHJlZml4');

@$core.Deprecated('Use inventoryS3BucketDestinationDescriptor instead')
const InventoryS3BucketDestination$json = {
  '1': 'InventoryS3BucketDestination',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'encryption',
      '3': 137702205,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryEncryption',
      '10': 'encryption'
    },
    {
      '1': 'format',
      '3': 531693427,
      '4': 1,
      '5': 14,
      '6': '.s3.InventoryFormat',
      '10': 'format'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
  ],
};

/// Descriptor for `InventoryS3BucketDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryS3BucketDestinationDescriptor = $convert.base64Decode(
    'ChxJbnZlbnRvcnlTM0J1Y2tldERlc3RpbmF0aW9uEh8KCWFjY291bnRpZBjSwbkfIAEoCVIJYW'
    'Njb3VudGlkEhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EjoKCmVuY3J5cHRpb24YvdbUQSAB'
    'KAsyFy5zMy5JbnZlbnRvcnlFbmNyeXB0aW9uUgplbmNyeXB0aW9uEi8KBmZvcm1hdBjz/sP9AS'
    'ABKA4yEy5zMy5JbnZlbnRvcnlGb3JtYXRSBmZvcm1hdBIaCgZwcmVmaXgY6rPTggEgASgJUgZw'
    'cmVmaXg=');

@$core.Deprecated('Use inventoryScheduleDescriptor instead')
const InventorySchedule$json = {
  '1': 'InventorySchedule',
  '2': [
    {
      '1': 'frequency',
      '3': 227673762,
      '4': 1,
      '5': 14,
      '6': '.s3.InventoryFrequency',
      '10': 'frequency'
    },
  ],
};

/// Descriptor for `InventorySchedule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryScheduleDescriptor = $convert.base64Decode(
    'ChFJbnZlbnRvcnlTY2hlZHVsZRI3CglmcmVxdWVuY3kYoo3IbCABKA4yFi5zMy5JbnZlbnRvcn'
    'lGcmVxdWVuY3lSCWZyZXF1ZW5jeQ==');

@$core.Deprecated('Use inventoryTableConfigurationDescriptor instead')
const InventoryTableConfiguration$json = {
  '1': 'InventoryTableConfiguration',
  '2': [
    {
      '1': 'configurationstate',
      '3': 458329211,
      '4': 1,
      '5': 14,
      '6': '.s3.InventoryConfigurationState',
      '10': 'configurationstate'
    },
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.s3.MetadataTableEncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
  ],
};

/// Descriptor for `InventoryTableConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryTableConfigurationDescriptor = $convert.base64Decode(
    'ChtJbnZlbnRvcnlUYWJsZUNvbmZpZ3VyYXRpb24SUwoSY29uZmlndXJhdGlvbnN0YXRlGPuYxt'
    'oBIAEoDjIfLnMzLkludmVudG9yeUNvbmZpZ3VyYXRpb25TdGF0ZVISY29uZmlndXJhdGlvbnN0'
    'YXRlEmUKF2VuY3J5cHRpb25jb25maWd1cmF0aW9uGPfG02sgASgLMiguczMuTWV0YWRhdGFUYW'
    'JsZUVuY3J5cHRpb25Db25maWd1cmF0aW9uUhdlbmNyeXB0aW9uY29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use inventoryTableConfigurationResultDescriptor instead')
const InventoryTableConfigurationResult$json = {
  '1': 'InventoryTableConfigurationResult',
  '2': [
    {
      '1': 'configurationstate',
      '3': 458329211,
      '4': 1,
      '5': 14,
      '6': '.s3.InventoryConfigurationState',
      '10': 'configurationstate'
    },
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.s3.ErrorDetails',
      '10': 'error'
    },
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {'1': 'tablestatus', '3': 207908810, '4': 1, '5': 9, '10': 'tablestatus'},
  ],
};

/// Descriptor for `InventoryTableConfigurationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryTableConfigurationResultDescriptor = $convert.base64Decode(
    'CiFJbnZlbnRvcnlUYWJsZUNvbmZpZ3VyYXRpb25SZXN1bHQSUwoSY29uZmlndXJhdGlvbnN0YX'
    'RlGPuYxtoBIAEoDjIfLnMzLkludmVudG9yeUNvbmZpZ3VyYXRpb25TdGF0ZVISY29uZmlndXJh'
    'dGlvbnN0YXRlEioKBWVycm9yGPK5tpwBIAEoCzIQLnMzLkVycm9yRGV0YWlsc1IFZXJyb3ISHg'
    'oIdGFibGVhcm4Y44DrzQEgASgJUgh0YWJsZWFybhIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0'
    'YWJsZW5hbWUSIwoLdGFibGVzdGF0dXMYyt+RYyABKAlSC3RhYmxlc3RhdHVz');

@$core.Deprecated('Use inventoryTableConfigurationUpdatesDescriptor instead')
const InventoryTableConfigurationUpdates$json = {
  '1': 'InventoryTableConfigurationUpdates',
  '2': [
    {
      '1': 'configurationstate',
      '3': 458329211,
      '4': 1,
      '5': 14,
      '6': '.s3.InventoryConfigurationState',
      '10': 'configurationstate'
    },
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.s3.MetadataTableEncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
  ],
};

/// Descriptor for `InventoryTableConfigurationUpdates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryTableConfigurationUpdatesDescriptor =
    $convert.base64Decode(
        'CiJJbnZlbnRvcnlUYWJsZUNvbmZpZ3VyYXRpb25VcGRhdGVzElMKEmNvbmZpZ3VyYXRpb25zdG'
        'F0ZRj7mMbaASABKA4yHy5zMy5JbnZlbnRvcnlDb25maWd1cmF0aW9uU3RhdGVSEmNvbmZpZ3Vy'
        'YXRpb25zdGF0ZRJlChdlbmNyeXB0aW9uY29uZmlndXJhdGlvbhj3xtNrIAEoCzIoLnMzLk1ldG'
        'FkYXRhVGFibGVFbmNyeXB0aW9uQ29uZmlndXJhdGlvblIXZW5jcnlwdGlvbmNvbmZpZ3VyYXRp'
        'b24=');

@$core.Deprecated('Use jSONInputDescriptor instead')
const JSONInput$json = {
  '1': 'JSONInput',
  '2': [
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.s3.JSONType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `JSONInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List jSONInputDescriptor = $convert.base64Decode(
    'CglKU09OSW5wdXQSJAoEdHlwZRjuoNeKASABKA4yDC5zMy5KU09OVHlwZVIEdHlwZQ==');

@$core.Deprecated('Use jSONOutputDescriptor instead')
const JSONOutput$json = {
  '1': 'JSONOutput',
  '2': [
    {
      '1': 'recorddelimiter',
      '3': 299337270,
      '4': 1,
      '5': 9,
      '10': 'recorddelimiter'
    },
  ],
};

/// Descriptor for `JSONOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List jSONOutputDescriptor = $convert.base64Decode(
    'CgpKU09OT3V0cHV0EiwKD3JlY29yZGRlbGltaXRlchi2jN6OASABKAlSD3JlY29yZGRlbGltaX'
    'Rlcg==');

@$core.Deprecated('Use journalTableConfigurationDescriptor instead')
const JournalTableConfiguration$json = {
  '1': 'JournalTableConfiguration',
  '2': [
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.s3.MetadataTableEncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {
      '1': 'recordexpiration',
      '3': 327741658,
      '4': 1,
      '5': 11,
      '6': '.s3.RecordExpiration',
      '10': 'recordexpiration'
    },
  ],
};

/// Descriptor for `JournalTableConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List journalTableConfigurationDescriptor = $convert.base64Decode(
    'ChlKb3VybmFsVGFibGVDb25maWd1cmF0aW9uEmUKF2VuY3J5cHRpb25jb25maWd1cmF0aW9uGP'
    'fG02sgASgLMiguczMuTWV0YWRhdGFUYWJsZUVuY3J5cHRpb25Db25maWd1cmF0aW9uUhdlbmNy'
    'eXB0aW9uY29uZmlndXJhdGlvbhJEChByZWNvcmRleHBpcmF0aW9uGNrho5wBIAEoCzIULnMzLl'
    'JlY29yZEV4cGlyYXRpb25SEHJlY29yZGV4cGlyYXRpb24=');

@$core.Deprecated('Use journalTableConfigurationResultDescriptor instead')
const JournalTableConfigurationResult$json = {
  '1': 'JournalTableConfigurationResult',
  '2': [
    {
      '1': 'error',
      '3': 328047858,
      '4': 1,
      '5': 11,
      '6': '.s3.ErrorDetails',
      '10': 'error'
    },
    {
      '1': 'recordexpiration',
      '3': 327741658,
      '4': 1,
      '5': 11,
      '6': '.s3.RecordExpiration',
      '10': 'recordexpiration'
    },
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {'1': 'tablestatus', '3': 207908810, '4': 1, '5': 9, '10': 'tablestatus'},
  ],
};

/// Descriptor for `JournalTableConfigurationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List journalTableConfigurationResultDescriptor = $convert.base64Decode(
    'Ch9Kb3VybmFsVGFibGVDb25maWd1cmF0aW9uUmVzdWx0EioKBWVycm9yGPK5tpwBIAEoCzIQLn'
    'MzLkVycm9yRGV0YWlsc1IFZXJyb3ISRAoQcmVjb3JkZXhwaXJhdGlvbhja4aOcASABKAsyFC5z'
    'My5SZWNvcmRFeHBpcmF0aW9uUhByZWNvcmRleHBpcmF0aW9uEh4KCHRhYmxlYXJuGOOA680BIA'
    'EoCVIIdGFibGVhcm4SIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1lEiMKC3RhYmxl'
    'c3RhdHVzGMrfkWMgASgJUgt0YWJsZXN0YXR1cw==');

@$core.Deprecated('Use journalTableConfigurationUpdatesDescriptor instead')
const JournalTableConfigurationUpdates$json = {
  '1': 'JournalTableConfigurationUpdates',
  '2': [
    {
      '1': 'recordexpiration',
      '3': 327741658,
      '4': 1,
      '5': 11,
      '6': '.s3.RecordExpiration',
      '10': 'recordexpiration'
    },
  ],
};

/// Descriptor for `JournalTableConfigurationUpdates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List journalTableConfigurationUpdatesDescriptor =
    $convert.base64Decode(
        'CiBKb3VybmFsVGFibGVDb25maWd1cmF0aW9uVXBkYXRlcxJEChByZWNvcmRleHBpcmF0aW9uGN'
        'rho5wBIAEoCzIULnMzLlJlY29yZEV4cGlyYXRpb25SEHJlY29yZGV4cGlyYXRpb24=');

@$core.Deprecated('Use lambdaFunctionConfigurationDescriptor instead')
const LambdaFunctionConfiguration$json = {
  '1': 'LambdaFunctionConfiguration',
  '2': [
    {
      '1': 'events',
      '3': 3416229,
      '4': 3,
      '5': 14,
      '6': '.s3.Event',
      '10': 'events'
    },
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.NotificationConfigurationFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'lambdafunctionarn',
      '3': 533137274,
      '4': 1,
      '5': 9,
      '10': 'lambdafunctionarn'
    },
  ],
};

/// Descriptor for `LambdaFunctionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionConfigurationDescriptor = $convert.base64Decode(
    'ChtMYW1iZGFGdW5jdGlvbkNvbmZpZ3VyYXRpb24SJAoGZXZlbnRzGKXB0AEgAygOMgkuczMuRX'
    'ZlbnRSBmV2ZW50cxI/CgZmaWx0ZXIYmIGnpQEgASgLMiMuczMuTm90aWZpY2F0aW9uQ29uZmln'
    'dXJhdGlvbkZpbHRlclIGZmlsdGVyEhIKAmlkGIHyorcBIAEoCVICaWQSMAoRbGFtYmRhZnVuY3'
    'Rpb25hcm4Y+o6c/gEgASgJUhFsYW1iZGFmdW5jdGlvbmFybg==');

@$core.Deprecated('Use lifecycleExpirationDescriptor instead')
const LifecycleExpiration$json = {
  '1': 'LifecycleExpiration',
  '2': [
    {'1': 'date', '3': 458388346, '4': 1, '5': 9, '10': 'date'},
    {'1': 'days', '3': 494075051, '4': 1, '5': 5, '10': 'days'},
    {
      '1': 'expiredobjectdeletemarker',
      '3': 232733109,
      '4': 1,
      '5': 8,
      '10': 'expiredobjectdeletemarker'
    },
  ],
};

/// Descriptor for `LifecycleExpiration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lifecycleExpirationDescriptor = $convert.base64Decode(
    'ChNMaWZlY3ljbGVFeHBpcmF0aW9uEhYKBGRhdGUY+ubJ2gEgASgJUgRkYXRlEhYKBGRheXMYq/'
    'nL6wEgASgFUgRkYXlzEj8KGWV4cGlyZWRvYmplY3RkZWxldGVtYXJrZXIYtfP8biABKAhSGWV4'
    'cGlyZWRvYmplY3RkZWxldGVtYXJrZXI=');

@$core.Deprecated('Use lifecycleRuleDescriptor instead')
const LifecycleRule$json = {
  '1': 'LifecycleRule',
  '2': [
    {
      '1': 'abortincompletemultipartupload',
      '3': 145054547,
      '4': 1,
      '5': 11,
      '6': '.s3.AbortIncompleteMultipartUpload',
      '10': 'abortincompletemultipartupload'
    },
    {
      '1': 'expiration',
      '3': 245879945,
      '4': 1,
      '5': 11,
      '6': '.s3.LifecycleExpiration',
      '10': 'expiration'
    },
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.LifecycleRuleFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384363361, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'noncurrentversionexpiration',
      '3': 73713275,
      '4': 1,
      '5': 11,
      '6': '.s3.NoncurrentVersionExpiration',
      '10': 'noncurrentversionexpiration'
    },
    {
      '1': 'noncurrentversiontransitions',
      '3': 254236814,
      '4': 3,
      '5': 11,
      '6': '.s3.NoncurrentVersionTransition',
      '10': 'noncurrentversiontransitions'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.ExpirationStatus',
      '10': 'status'
    },
    {
      '1': 'transitions',
      '3': 435714020,
      '4': 3,
      '5': 11,
      '6': '.s3.Transition',
      '10': 'transitions'
    },
  ],
};

/// Descriptor for `LifecycleRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lifecycleRuleDescriptor = $convert.base64Decode(
    'Cg1MaWZlY3ljbGVSdWxlEm0KHmFib3J0aW5jb21wbGV0ZW11bHRpcGFydHVwbG9hZBjTtpVFIA'
    'EoCzIiLnMzLkFib3J0SW5jb21wbGV0ZU11bHRpcGFydFVwbG9hZFIeYWJvcnRpbmNvbXBsZXRl'
    'bXVsdGlwYXJ0dXBsb2FkEjoKCmV4cGlyYXRpb24YiamfdSABKAsyFy5zMy5MaWZlY3ljbGVFeH'
    'BpcmF0aW9uUgpleHBpcmF0aW9uEjMKBmZpbHRlchiYgaelASABKAsyFy5zMy5MaWZlY3ljbGVS'
    'dWxlRmlsdGVyUgZmaWx0ZXISEgoCaWQY4dajtwEgASgJUgJpZBJkChtub25jdXJyZW50dmVyc2'
    'lvbmV4cGlyYXRpb24Y+4yTIyABKAsyHy5zMy5Ob25jdXJyZW50VmVyc2lvbkV4cGlyYXRpb25S'
    'G25vbmN1cnJlbnR2ZXJzaW9uZXhwaXJhdGlvbhJmChxub25jdXJyZW50dmVyc2lvbnRyYW5zaX'
    'Rpb25zGI6xnXkgAygLMh8uczMuTm9uY3VycmVudFZlcnNpb25UcmFuc2l0aW9uUhxub25jdXJy'
    'ZW50dmVyc2lvbnRyYW5zaXRpb25zEhoKBnByZWZpeBjqs9OCASABKAlSBnByZWZpeBIvCgZzdG'
    'F0dXMYkOT7AiABKA4yFC5zMy5FeHBpcmF0aW9uU3RhdHVzUgZzdGF0dXMSNAoLdHJhbnNpdGlv'
    'bnMY5O/hzwEgAygLMg4uczMuVHJhbnNpdGlvblILdHJhbnNpdGlvbnM=');

@$core.Deprecated('Use lifecycleRuleAndOperatorDescriptor instead')
const LifecycleRuleAndOperator$json = {
  '1': 'LifecycleRuleAndOperator',
  '2': [
    {
      '1': 'objectsizegreaterthan',
      '3': 283649577,
      '4': 1,
      '5': 3,
      '10': 'objectsizegreaterthan'
    },
    {
      '1': 'objectsizelessthan',
      '3': 436793998,
      '4': 1,
      '5': 3,
      '10': 'objectsizelessthan'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `LifecycleRuleAndOperator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lifecycleRuleAndOperatorDescriptor = $convert.base64Decode(
    'ChhMaWZlY3ljbGVSdWxlQW5kT3BlcmF0b3ISOAoVb2JqZWN0c2l6ZWdyZWF0ZXJ0aGFuGKnMoI'
    'cBIAEoA1IVb2JqZWN0c2l6ZWdyZWF0ZXJ0aGFuEjIKEm9iamVjdHNpemVsZXNzdGhhbhiO5aPQ'
    'ASABKANSEm9iamVjdHNpemVsZXNzdGhhbhIaCgZwcmVmaXgY6rPTggEgASgJUgZwcmVmaXgSHw'
    'oEdGFncxjBwfa1ASADKAsyBy5zMy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use lifecycleRuleFilterDescriptor instead')
const LifecycleRuleFilter$json = {
  '1': 'LifecycleRuleFilter',
  '2': [
    {
      '1': 'and',
      '3': 297135431,
      '4': 1,
      '5': 11,
      '6': '.s3.LifecycleRuleAndOperator',
      '10': 'and'
    },
    {
      '1': 'objectsizegreaterthan',
      '3': 283649577,
      '4': 1,
      '5': 3,
      '10': 'objectsizegreaterthan'
    },
    {
      '1': 'objectsizelessthan',
      '3': 436793998,
      '4': 1,
      '5': 3,
      '10': 'objectsizelessthan'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {'1': 'tag', '3': 411259956, '4': 1, '5': 11, '6': '.s3.Tag', '10': 'tag'},
  ],
};

/// Descriptor for `LifecycleRuleFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lifecycleRuleFilterDescriptor = $convert.base64Decode(
    'ChNMaWZlY3ljbGVSdWxlRmlsdGVyEjIKA2FuZBjH2teNASABKAsyHC5zMy5MaWZlY3ljbGVSdW'
    'xlQW5kT3BlcmF0b3JSA2FuZBI4ChVvYmplY3RzaXplZ3JlYXRlcnRoYW4YqcyghwEgASgDUhVv'
    'YmplY3RzaXplZ3JlYXRlcnRoYW4SMgoSb2JqZWN0c2l6ZWxlc3N0aGFuGI7lo9ABIAEoA1ISb2'
    'JqZWN0c2l6ZWxlc3N0aGFuEhoKBnByZWZpeBjqs9OCASABKAlSBnByZWZpeBIdCgN0YWcYtKiN'
    'xAEgASgLMgcuczMuVGFnUgN0YWc=');

@$core
    .Deprecated('Use listBucketAnalyticsConfigurationsOutputDescriptor instead')
const ListBucketAnalyticsConfigurationsOutput$json = {
  '1': 'ListBucketAnalyticsConfigurationsOutput',
  '2': [
    {
      '1': 'analyticsconfigurationlist',
      '3': 223424476,
      '4': 3,
      '5': 11,
      '6': '.s3.AnalyticsConfiguration',
      '10': 'analyticsconfigurationlist'
    },
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'nextcontinuationtoken',
      '3': 260840781,
      '4': 1,
      '5': 9,
      '10': 'nextcontinuationtoken'
    },
  ],
};

/// Descriptor for `ListBucketAnalyticsConfigurationsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBucketAnalyticsConfigurationsOutputDescriptor =
    $convert.base64Decode(
        'CidMaXN0QnVja2V0QW5hbHl0aWNzQ29uZmlndXJhdGlvbnNPdXRwdXQSXQoaYW5hbHl0aWNzY2'
        '9uZmlndXJhdGlvbmxpc3QY3N/EaiADKAsyGi5zMy5BbmFseXRpY3NDb25maWd1cmF0aW9uUhph'
        'bmFseXRpY3Njb25maWd1cmF0aW9ubGlzdBIwChFjb250aW51YXRpb250b2tlbhjoysCIASABKA'
        'lSEWNvbnRpbnVhdGlvbnRva2VuEiMKC2lzdHJ1bmNhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRl'
        'ZBI3ChVuZXh0Y29udGludWF0aW9udG9rZW4YzbqwfCABKAlSFW5leHRjb250aW51YXRpb250b2'
        'tlbg==');

@$core.Deprecated(
    'Use listBucketAnalyticsConfigurationsRequestDescriptor instead')
const ListBucketAnalyticsConfigurationsRequest$json = {
  '1': 'ListBucketAnalyticsConfigurationsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `ListBucketAnalyticsConfigurationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBucketAnalyticsConfigurationsRequestDescriptor =
    $convert.base64Decode(
        'CihMaXN0QnVja2V0QW5hbHl0aWNzQ29uZmlndXJhdGlvbnNSZXF1ZXN0EhkKBmJ1Y2tldBjY6r'
        'gaIAEoCVIGYnVja2V0EjAKEWNvbnRpbnVhdGlvbnRva2VuGOjKwIgBIAEoCVIRY29udGludWF0'
        'aW9udG9rZW4SMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZX'
        'Rvd25lcg==');

@$core.Deprecated(
    'Use listBucketIntelligentTieringConfigurationsOutputDescriptor instead')
const ListBucketIntelligentTieringConfigurationsOutput$json = {
  '1': 'ListBucketIntelligentTieringConfigurationsOutput',
  '2': [
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {
      '1': 'intelligenttieringconfigurationlist',
      '3': 406526373,
      '4': 3,
      '5': 11,
      '6': '.s3.IntelligentTieringConfiguration',
      '10': 'intelligenttieringconfigurationlist'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'nextcontinuationtoken',
      '3': 260840781,
      '4': 1,
      '5': 9,
      '10': 'nextcontinuationtoken'
    },
  ],
};

/// Descriptor for `ListBucketIntelligentTieringConfigurationsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listBucketIntelligentTieringConfigurationsOutputDescriptor =
    $convert.base64Decode(
        'CjBMaXN0QnVja2V0SW50ZWxsaWdlbnRUaWVyaW5nQ29uZmlndXJhdGlvbnNPdXRwdXQSMAoRY2'
        '9udGludWF0aW9udG9rZW4Y6MrAiAEgASgJUhFjb250aW51YXRpb250b2tlbhJ5CiNpbnRlbGxp'
        'Z2VudHRpZXJpbmdjb25maWd1cmF0aW9ubGlzdBils+zBASADKAsyIy5zMy5JbnRlbGxpZ2VudF'
        'RpZXJpbmdDb25maWd1cmF0aW9uUiNpbnRlbGxpZ2VudHRpZXJpbmdjb25maWd1cmF0aW9ubGlz'
        'dBIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSNwoVbmV4dGNvbnRpbnVhdG'
        'lvbnRva2VuGM26sHwgASgJUhVuZXh0Y29udGludWF0aW9udG9rZW4=');

@$core.Deprecated(
    'Use listBucketIntelligentTieringConfigurationsRequestDescriptor instead')
const ListBucketIntelligentTieringConfigurationsRequest$json = {
  '1': 'ListBucketIntelligentTieringConfigurationsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `ListBucketIntelligentTieringConfigurationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listBucketIntelligentTieringConfigurationsRequestDescriptor =
    $convert.base64Decode(
        'CjFMaXN0QnVja2V0SW50ZWxsaWdlbnRUaWVyaW5nQ29uZmlndXJhdGlvbnNSZXF1ZXN0EhkKBm'
        'J1Y2tldBjY6rgaIAEoCVIGYnVja2V0EjAKEWNvbnRpbnVhdGlvbnRva2VuGOjKwIgBIAEoCVIR'
        'Y29udGludWF0aW9udG9rZW4SMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZW'
        'N0ZWRidWNrZXRvd25lcg==');

@$core
    .Deprecated('Use listBucketInventoryConfigurationsOutputDescriptor instead')
const ListBucketInventoryConfigurationsOutput$json = {
  '1': 'ListBucketInventoryConfigurationsOutput',
  '2': [
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {
      '1': 'inventoryconfigurationlist',
      '3': 389500552,
      '4': 3,
      '5': 11,
      '6': '.s3.InventoryConfiguration',
      '10': 'inventoryconfigurationlist'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'nextcontinuationtoken',
      '3': 260840781,
      '4': 1,
      '5': 9,
      '10': 'nextcontinuationtoken'
    },
  ],
};

/// Descriptor for `ListBucketInventoryConfigurationsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBucketInventoryConfigurationsOutputDescriptor =
    $convert.base64Decode(
        'CidMaXN0QnVja2V0SW52ZW50b3J5Q29uZmlndXJhdGlvbnNPdXRwdXQSMAoRY29udGludWF0aW'
        '9udG9rZW4Y6MrAiAEgASgJUhFjb250aW51YXRpb250b2tlbhJeChppbnZlbnRvcnljb25maWd1'
        'cmF0aW9ubGlzdBiInd25ASADKAsyGi5zMy5JbnZlbnRvcnlDb25maWd1cmF0aW9uUhppbnZlbn'
        'Rvcnljb25maWd1cmF0aW9ubGlzdBIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0'
        'ZWQSNwoVbmV4dGNvbnRpbnVhdGlvbnRva2VuGM26sHwgASgJUhVuZXh0Y29udGludWF0aW9udG'
        '9rZW4=');

@$core.Deprecated(
    'Use listBucketInventoryConfigurationsRequestDescriptor instead')
const ListBucketInventoryConfigurationsRequest$json = {
  '1': 'ListBucketInventoryConfigurationsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `ListBucketInventoryConfigurationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBucketInventoryConfigurationsRequestDescriptor =
    $convert.base64Decode(
        'CihMaXN0QnVja2V0SW52ZW50b3J5Q29uZmlndXJhdGlvbnNSZXF1ZXN0EhkKBmJ1Y2tldBjY6r'
        'gaIAEoCVIGYnVja2V0EjAKEWNvbnRpbnVhdGlvbnRva2VuGOjKwIgBIAEoCVIRY29udGludWF0'
        'aW9udG9rZW4SMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZX'
        'Rvd25lcg==');

@$core.Deprecated('Use listBucketMetricsConfigurationsOutputDescriptor instead')
const ListBucketMetricsConfigurationsOutput$json = {
  '1': 'ListBucketMetricsConfigurationsOutput',
  '2': [
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {
      '1': 'metricsconfigurationlist',
      '3': 374737701,
      '4': 3,
      '5': 11,
      '6': '.s3.MetricsConfiguration',
      '10': 'metricsconfigurationlist'
    },
    {
      '1': 'nextcontinuationtoken',
      '3': 260840781,
      '4': 1,
      '5': 9,
      '10': 'nextcontinuationtoken'
    },
  ],
};

/// Descriptor for `ListBucketMetricsConfigurationsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBucketMetricsConfigurationsOutputDescriptor =
    $convert.base64Decode(
        'CiVMaXN0QnVja2V0TWV0cmljc0NvbmZpZ3VyYXRpb25zT3V0cHV0EjAKEWNvbnRpbnVhdGlvbn'
        'Rva2VuGOjKwIgBIAEoCVIRY29udGludWF0aW9udG9rZW4SIwoLaXN0cnVuY2F0ZWQY2p+4cyAB'
        'KAhSC2lzdHJ1bmNhdGVkElgKGG1ldHJpY3Njb25maWd1cmF0aW9ubGlzdBilltiyASADKAsyGC'
        '5zMy5NZXRyaWNzQ29uZmlndXJhdGlvblIYbWV0cmljc2NvbmZpZ3VyYXRpb25saXN0EjcKFW5l'
        'eHRjb250aW51YXRpb250b2tlbhjNurB8IAEoCVIVbmV4dGNvbnRpbnVhdGlvbnRva2Vu');

@$core
    .Deprecated('Use listBucketMetricsConfigurationsRequestDescriptor instead')
const ListBucketMetricsConfigurationsRequest$json = {
  '1': 'ListBucketMetricsConfigurationsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `ListBucketMetricsConfigurationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBucketMetricsConfigurationsRequestDescriptor =
    $convert.base64Decode(
        'CiZMaXN0QnVja2V0TWV0cmljc0NvbmZpZ3VyYXRpb25zUmVxdWVzdBIZCgZidWNrZXQY2Oq4Gi'
        'ABKAlSBmJ1Y2tldBIwChFjb250aW51YXRpb250b2tlbhjoysCIASABKAlSEWNvbnRpbnVhdGlv'
        'bnRva2VuEjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3'
        'duZXI=');

@$core.Deprecated('Use listBucketsOutputDescriptor instead')
const ListBucketsOutput$json = {
  '1': 'ListBucketsOutput',
  '2': [
    {
      '1': 'buckets',
      '3': 404596653,
      '4': 3,
      '5': 11,
      '6': '.s3.Bucket',
      '10': 'buckets'
    },
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
  ],
};

/// Descriptor for `ListBucketsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBucketsOutputDescriptor = $convert.base64Decode(
    'ChFMaXN0QnVja2V0c091dHB1dBIoCgdidWNrZXRzGK3P9sABIAMoCzIKLnMzLkJ1Y2tldFIHYn'
    'Vja2V0cxIwChFjb250aW51YXRpb250b2tlbhjoysCIASABKAlSEWNvbnRpbnVhdGlvbnRva2Vu'
    'EiMKBW93bmVyGPX8itkBIAEoCzIJLnMzLk93bmVyUgVvd25lchIaCgZwcmVmaXgY6rPTggEgAS'
    'gJUgZwcmVmaXg=');

@$core.Deprecated('Use listBucketsRequestDescriptor instead')
const ListBucketsRequest$json = {
  '1': 'ListBucketsRequest',
  '2': [
    {'1': 'bucketregion', '3': 309298816, '4': 1, '5': 9, '10': 'bucketregion'},
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {'1': 'maxbuckets', '3': 264636545, '4': 1, '5': 5, '10': 'maxbuckets'},
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
  ],
};

/// Descriptor for `ListBucketsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBucketsRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0QnVja2V0c1JlcXVlc3QSJgoMYnVja2V0cmVnaW9uGICNvpMBIAEoCVIMYnVja2V0cm'
    'VnaW9uEjAKEWNvbnRpbnVhdGlvbnRva2VuGOjKwIgBIAEoCVIRY29udGludWF0aW9udG9rZW4S'
    'IQoKbWF4YnVja2V0cxiBkZh+IAEoBVIKbWF4YnVja2V0cxIaCgZwcmVmaXgY6rPTggEgASgJUg'
    'ZwcmVmaXg=');

@$core.Deprecated('Use listDirectoryBucketsOutputDescriptor instead')
const ListDirectoryBucketsOutput$json = {
  '1': 'ListDirectoryBucketsOutput',
  '2': [
    {
      '1': 'buckets',
      '3': 404596653,
      '4': 3,
      '5': 11,
      '6': '.s3.Bucket',
      '10': 'buckets'
    },
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
  ],
};

/// Descriptor for `ListDirectoryBucketsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDirectoryBucketsOutputDescriptor =
    $convert.base64Decode(
        'ChpMaXN0RGlyZWN0b3J5QnVja2V0c091dHB1dBIoCgdidWNrZXRzGK3P9sABIAMoCzIKLnMzLk'
        'J1Y2tldFIHYnVja2V0cxIwChFjb250aW51YXRpb250b2tlbhjoysCIASABKAlSEWNvbnRpbnVh'
        'dGlvbnRva2Vu');

@$core.Deprecated('Use listDirectoryBucketsRequestDescriptor instead')
const ListDirectoryBucketsRequest$json = {
  '1': 'ListDirectoryBucketsRequest',
  '2': [
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {
      '1': 'maxdirectorybuckets',
      '3': 424013420,
      '4': 1,
      '5': 5,
      '10': 'maxdirectorybuckets'
    },
  ],
};

/// Descriptor for `ListDirectoryBucketsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDirectoryBucketsRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0RGlyZWN0b3J5QnVja2V0c1JlcXVlc3QSMAoRY29udGludWF0aW9udG9rZW4Y6MrAiA'
        'EgASgJUhFjb250aW51YXRpb250b2tlbhI0ChNtYXhkaXJlY3RvcnlidWNrZXRzGOzcl8oBIAEo'
        'BVITbWF4ZGlyZWN0b3J5YnVja2V0cw==');

@$core.Deprecated('Use listMultipartUploadsOutputDescriptor instead')
const ListMultipartUploadsOutput$json = {
  '1': 'ListMultipartUploadsOutput',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'commonprefixes',
      '3': 83168757,
      '4': 3,
      '5': 11,
      '6': '.s3.CommonPrefix',
      '10': 'commonprefixes'
    },
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'encodingtype',
      '3': 532628025,
      '4': 1,
      '5': 14,
      '6': '.s3.EncodingType',
      '10': 'encodingtype'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'keymarker', '3': 485505207, '4': 1, '5': 9, '10': 'keymarker'},
    {'1': 'maxuploads', '3': 111552074, '4': 1, '5': 5, '10': 'maxuploads'},
    {
      '1': 'nextkeymarker',
      '3': 433290006,
      '4': 1,
      '5': 9,
      '10': 'nextkeymarker'
    },
    {
      '1': 'nextuploadidmarker',
      '3': 528782759,
      '4': 1,
      '5': 9,
      '10': 'nextuploadidmarker'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'uploadidmarker',
      '3': 31676084,
      '4': 1,
      '5': 9,
      '10': 'uploadidmarker'
    },
    {
      '1': 'uploads',
      '3': 75826334,
      '4': 3,
      '5': 11,
      '6': '.s3.MultipartUpload',
      '10': 'uploads'
    },
  ],
};

/// Descriptor for `ListMultipartUploadsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMultipartUploadsOutputDescriptor = $convert.base64Decode(
    'ChpMaXN0TXVsdGlwYXJ0VXBsb2Fkc091dHB1dBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
    'I7Cg5jb21tb25wcmVmaXhlcxj1m9QnIAMoCzIQLnMzLkNvbW1vblByZWZpeFIOY29tbW9ucHJl'
    'Zml4ZXMSIAoJZGVsaW1pdGVyGJvZiJABIAEoCVIJZGVsaW1pdGVyEjgKDGVuY29kaW5ndHlwZR'
    'i5hP39ASABKA4yEC5zMy5FbmNvZGluZ1R5cGVSDGVuY29kaW5ndHlwZRIjCgtpc3RydW5jYXRl'
    'ZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSIAoJa2V5bWFya2VyGLfxwOcBIAEoCVIJa2V5bWFya2'
    'VyEiEKCm1heHVwbG9hZHMYysyYNSABKAVSCm1heHVwbG9hZHMSKAoNbmV4dGtleW1hcmtlchiW'
    '9s3OASABKAlSDW5leHRrZXltYXJrZXISMgoSbmV4dHVwbG9hZGlkbWFya2VyGKerkvwBIAEoCV'
    'ISbmV4dHVwbG9hZGlkbWFya2VyEhoKBnByZWZpeBjqs9OCASABKAlSBnByZWZpeBI+Cg5yZXF1'
    'ZXN0Y2hhcmdlZBiT0Ku5ASABKA4yEi5zMy5SZXF1ZXN0Q2hhcmdlZFIOcmVxdWVzdGNoYXJnZW'
    'QSKQoOdXBsb2FkaWRtYXJrZXIYtK2NDyABKAlSDnVwbG9hZGlkbWFya2VyEjAKB3VwbG9hZHMY'
    'nomUJCADKAsyEy5zMy5NdWx0aXBhcnRVcGxvYWRSB3VwbG9hZHM=');

@$core.Deprecated('Use listMultipartUploadsRequestDescriptor instead')
const ListMultipartUploadsRequest$json = {
  '1': 'ListMultipartUploadsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'encodingtype',
      '3': 532628025,
      '4': 1,
      '5': 14,
      '6': '.s3.EncodingType',
      '10': 'encodingtype'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'keymarker', '3': 485505207, '4': 1, '5': 9, '10': 'keymarker'},
    {'1': 'maxuploads', '3': 111552074, '4': 1, '5': 5, '10': 'maxuploads'},
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'uploadidmarker',
      '3': 31676084,
      '4': 1,
      '5': 9,
      '10': 'uploadidmarker'
    },
  ],
};

/// Descriptor for `ListMultipartUploadsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMultipartUploadsRequestDescriptor = $convert.base64Decode(
    'ChtMaXN0TXVsdGlwYXJ0VXBsb2Fkc1JlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZX'
    'QSIAoJZGVsaW1pdGVyGJvZiJABIAEoCVIJZGVsaW1pdGVyEjgKDGVuY29kaW5ndHlwZRi5hP39'
    'ASABKA4yEC5zMy5FbmNvZGluZ1R5cGVSDGVuY29kaW5ndHlwZRIzChNleHBlY3RlZGJ1Y2tldG'
    '93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEiAKCWtleW1hcmtlchi38cDnASAB'
    'KAlSCWtleW1hcmtlchIhCgptYXh1cGxvYWRzGMrMmDUgASgFUgptYXh1cGxvYWRzEhoKBnByZW'
    'ZpeBjqs9OCASABKAlSBnByZWZpeBI4CgxyZXF1ZXN0cGF5ZXIYpObh9QEgASgOMhAuczMuUmVx'
    'dWVzdFBheWVyUgxyZXF1ZXN0cGF5ZXISKQoOdXBsb2FkaWRtYXJrZXIYtK2NDyABKAlSDnVwbG'
    '9hZGlkbWFya2Vy');

@$core.Deprecated('Use listObjectVersionsOutputDescriptor instead')
const ListObjectVersionsOutput$json = {
  '1': 'ListObjectVersionsOutput',
  '2': [
    {
      '1': 'commonprefixes',
      '3': 83168757,
      '4': 3,
      '5': 11,
      '6': '.s3.CommonPrefix',
      '10': 'commonprefixes'
    },
    {
      '1': 'deletemarkers',
      '3': 376648970,
      '4': 3,
      '5': 11,
      '6': '.s3.DeleteMarkerEntry',
      '10': 'deletemarkers'
    },
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'encodingtype',
      '3': 532628025,
      '4': 1,
      '5': 14,
      '6': '.s3.EncodingType',
      '10': 'encodingtype'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'keymarker', '3': 485505207, '4': 1, '5': 9, '10': 'keymarker'},
    {'1': 'maxkeys', '3': 247655034, '4': 1, '5': 5, '10': 'maxkeys'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nextkeymarker',
      '3': 433290006,
      '4': 1,
      '5': 9,
      '10': 'nextkeymarker'
    },
    {
      '1': 'nextversionidmarker',
      '3': 449586812,
      '4': 1,
      '5': 9,
      '10': 'nextversionidmarker'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'versionidmarker',
      '3': 335417777,
      '4': 1,
      '5': 9,
      '10': 'versionidmarker'
    },
    {
      '1': 'versions',
      '3': 252099085,
      '4': 3,
      '5': 11,
      '6': '.s3.ObjectVersion',
      '10': 'versions'
    },
  ],
};

/// Descriptor for `ListObjectVersionsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listObjectVersionsOutputDescriptor = $convert.base64Decode(
    'ChhMaXN0T2JqZWN0VmVyc2lvbnNPdXRwdXQSOwoOY29tbW9ucHJlZml4ZXMY9ZvUJyADKAsyEC'
    '5zMy5Db21tb25QcmVmaXhSDmNvbW1vbnByZWZpeGVzEj8KDWRlbGV0ZW1hcmtlcnMYiurMswEg'
    'AygLMhUuczMuRGVsZXRlTWFya2VyRW50cnlSDWRlbGV0ZW1hcmtlcnMSIAoJZGVsaW1pdGVyGJ'
    'vZiJABIAEoCVIJZGVsaW1pdGVyEjgKDGVuY29kaW5ndHlwZRi5hP39ASABKA4yEC5zMy5FbmNv'
    'ZGluZ1R5cGVSDGVuY29kaW5ndHlwZRIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2'
    'F0ZWQSIAoJa2V5bWFya2VyGLfxwOcBIAEoCVIJa2V5bWFya2VyEhsKB21heGtleXMY+tSLdiAB'
    'KAVSB21heGtleXMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIoCg1uZXh0a2V5bWFya2VyGJb2zc'
    '4BIAEoCVINbmV4dGtleW1hcmtlchI0ChNuZXh0dmVyc2lvbmlkbWFya2VyGPzMsNYBIAEoCVIT'
    'bmV4dHZlcnNpb25pZG1hcmtlchIaCgZwcmVmaXgY6rPTggEgASgJUgZwcmVmaXgSPgoOcmVxdW'
    'VzdGNoYXJnZWQYk9CruQEgASgOMhIuczMuUmVxdWVzdENoYXJnZWRSDnJlcXVlc3RjaGFyZ2Vk'
    'EiwKD3ZlcnNpb25pZG1hcmtlchixo/ifASABKAlSD3ZlcnNpb25pZG1hcmtlchIwCgh2ZXJzaW'
    '9ucxiN9Jp4IAMoCzIRLnMzLk9iamVjdFZlcnNpb25SCHZlcnNpb25z');

@$core.Deprecated('Use listObjectVersionsRequestDescriptor instead')
const ListObjectVersionsRequest$json = {
  '1': 'ListObjectVersionsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'encodingtype',
      '3': 532628025,
      '4': 1,
      '5': 14,
      '6': '.s3.EncodingType',
      '10': 'encodingtype'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'keymarker', '3': 485505207, '4': 1, '5': 9, '10': 'keymarker'},
    {'1': 'maxkeys', '3': 247655034, '4': 1, '5': 5, '10': 'maxkeys'},
    {
      '1': 'optionalobjectattributes',
      '3': 147719620,
      '4': 3,
      '5': 14,
      '6': '.s3.OptionalObjectAttributes',
      '10': 'optionalobjectattributes'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'versionidmarker',
      '3': 335417777,
      '4': 1,
      '5': 9,
      '10': 'versionidmarker'
    },
  ],
};

/// Descriptor for `ListObjectVersionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listObjectVersionsRequestDescriptor = $convert.base64Decode(
    'ChlMaXN0T2JqZWN0VmVyc2lvbnNSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0Ei'
    'AKCWRlbGltaXRlchib2YiQASABKAlSCWRlbGltaXRlchI4CgxlbmNvZGluZ3R5cGUYuYT9/QEg'
    'ASgOMhAuczMuRW5jb2RpbmdUeXBlUgxlbmNvZGluZ3R5cGUSMwoTZXhwZWN0ZWRidWNrZXRvd2'
    '5lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchIgCglrZXltYXJrZXIYt/HA5wEgASgJ'
    'UglrZXltYXJrZXISGwoHbWF4a2V5cxj61It2IAEoBVIHbWF4a2V5cxJbChhvcHRpb25hbG9iam'
    'VjdGF0dHJpYnV0ZXMYxIu4RiADKA4yHC5zMy5PcHRpb25hbE9iamVjdEF0dHJpYnV0ZXNSGG9w'
    'dGlvbmFsb2JqZWN0YXR0cmlidXRlcxIaCgZwcmVmaXgY6rPTggEgASgJUgZwcmVmaXgSOAoMcm'
    'VxdWVzdHBheWVyGKTm4fUBIAEoDjIQLnMzLlJlcXVlc3RQYXllclIMcmVxdWVzdHBheWVyEiwK'
    'D3ZlcnNpb25pZG1hcmtlchixo/ifASABKAlSD3ZlcnNpb25pZG1hcmtlcg==');

@$core.Deprecated('Use listObjectsOutputDescriptor instead')
const ListObjectsOutput$json = {
  '1': 'ListObjectsOutput',
  '2': [
    {
      '1': 'commonprefixes',
      '3': 83168757,
      '4': 3,
      '5': 11,
      '6': '.s3.CommonPrefix',
      '10': 'commonprefixes'
    },
    {
      '1': 'contents',
      '3': 119498692,
      '4': 3,
      '5': 11,
      '6': '.s3.Object',
      '10': 'contents'
    },
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'encodingtype',
      '3': 532628025,
      '4': 1,
      '5': 14,
      '6': '.s3.EncodingType',
      '10': 'encodingtype'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxkeys', '3': 247655034, '4': 1, '5': 5, '10': 'maxkeys'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `ListObjectsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listObjectsOutputDescriptor = $convert.base64Decode(
    'ChFMaXN0T2JqZWN0c091dHB1dBI7Cg5jb21tb25wcmVmaXhlcxj1m9QnIAMoCzIQLnMzLkNvbW'
    '1vblByZWZpeFIOY29tbW9ucHJlZml4ZXMSKQoIY29udGVudHMYxM/9OCADKAsyCi5zMy5PYmpl'
    'Y3RSCGNvbnRlbnRzEiAKCWRlbGltaXRlchib2YiQASABKAlSCWRlbGltaXRlchI4CgxlbmNvZG'
    'luZ3R5cGUYuYT9/QEgASgOMhAuczMuRW5jb2RpbmdUeXBlUgxlbmNvZGluZ3R5cGUSIwoLaXN0'
    'cnVuY2F0ZWQY2p+4cyABKAhSC2lzdHJ1bmNhdGVkEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2'
    'VyEhsKB21heGtleXMY+tSLdiABKAVSB21heGtleXMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIi'
    'CgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIaCgZwcmVmaXgY6rPTggEgASgJUg'
    'ZwcmVmaXgSPgoOcmVxdWVzdGNoYXJnZWQYk9CruQEgASgOMhIuczMuUmVxdWVzdENoYXJnZWRS'
    'DnJlcXVlc3RjaGFyZ2Vk');

@$core.Deprecated('Use listObjectsRequestDescriptor instead')
const ListObjectsRequest$json = {
  '1': 'ListObjectsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'encodingtype',
      '3': 532628025,
      '4': 1,
      '5': 14,
      '6': '.s3.EncodingType',
      '10': 'encodingtype'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {'1': 'maxkeys', '3': 247655034, '4': 1, '5': 5, '10': 'maxkeys'},
    {
      '1': 'optionalobjectattributes',
      '3': 147719620,
      '4': 3,
      '5': 14,
      '6': '.s3.OptionalObjectAttributes',
      '10': 'optionalobjectattributes'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
  ],
};

/// Descriptor for `ListObjectsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listObjectsRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0T2JqZWN0c1JlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZXQSIAoJZGVsaW'
    '1pdGVyGJvZiJABIAEoCVIJZGVsaW1pdGVyEjgKDGVuY29kaW5ndHlwZRi5hP39ASABKA4yEC5z'
    'My5FbmNvZGluZ1R5cGVSDGVuY29kaW5ndHlwZRIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D'
    '4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEhsK'
    'B21heGtleXMY+tSLdiABKAVSB21heGtleXMSWwoYb3B0aW9uYWxvYmplY3RhdHRyaWJ1dGVzGM'
    'SLuEYgAygOMhwuczMuT3B0aW9uYWxPYmplY3RBdHRyaWJ1dGVzUhhvcHRpb25hbG9iamVjdGF0'
    'dHJpYnV0ZXMSGgoGcHJlZml4GOqz04IBIAEoCVIGcHJlZml4EjgKDHJlcXVlc3RwYXllchik5u'
    'H1ASABKA4yEC5zMy5SZXF1ZXN0UGF5ZXJSDHJlcXVlc3RwYXllcg==');

@$core.Deprecated('Use listObjectsV2OutputDescriptor instead')
const ListObjectsV2Output$json = {
  '1': 'ListObjectsV2Output',
  '2': [
    {
      '1': 'commonprefixes',
      '3': 83168757,
      '4': 3,
      '5': 11,
      '6': '.s3.CommonPrefix',
      '10': 'commonprefixes'
    },
    {
      '1': 'contents',
      '3': 119498692,
      '4': 3,
      '5': 11,
      '6': '.s3.Object',
      '10': 'contents'
    },
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'encodingtype',
      '3': 532628025,
      '4': 1,
      '5': 14,
      '6': '.s3.EncodingType',
      '10': 'encodingtype'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'keycount', '3': 253426132, '4': 1, '5': 5, '10': 'keycount'},
    {'1': 'maxkeys', '3': 247655034, '4': 1, '5': 5, '10': 'maxkeys'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nextcontinuationtoken',
      '3': 260840781,
      '4': 1,
      '5': 9,
      '10': 'nextcontinuationtoken'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {'1': 'startafter', '3': 339670328, '4': 1, '5': 9, '10': 'startafter'},
  ],
};

/// Descriptor for `ListObjectsV2Output`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listObjectsV2OutputDescriptor = $convert.base64Decode(
    'ChNMaXN0T2JqZWN0c1YyT3V0cHV0EjsKDmNvbW1vbnByZWZpeGVzGPWb1CcgAygLMhAuczMuQ2'
    '9tbW9uUHJlZml4Ug5jb21tb25wcmVmaXhlcxIpCghjb250ZW50cxjEz/04IAMoCzIKLnMzLk9i'
    'amVjdFIIY29udGVudHMSMAoRY29udGludWF0aW9udG9rZW4Y6MrAiAEgASgJUhFjb250aW51YX'
    'Rpb250b2tlbhIgCglkZWxpbWl0ZXIYm9mIkAEgASgJUglkZWxpbWl0ZXISOAoMZW5jb2Rpbmd0'
    'eXBlGLmE/f0BIAEoDjIQLnMzLkVuY29kaW5nVHlwZVIMZW5jb2Rpbmd0eXBlEiMKC2lzdHJ1bm'
    'NhdGVkGNqfuHMgASgIUgtpc3RydW5jYXRlZBIdCghrZXljb3VudBjU8+t4IAEoBVIIa2V5Y291'
    'bnQSGwoHbWF4a2V5cxj61It2IAEoBVIHbWF4a2V5cxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEj'
    'cKFW5leHRjb250aW51YXRpb250b2tlbhjNurB8IAEoCVIVbmV4dGNvbnRpbnVhdGlvbnRva2Vu'
    'EhoKBnByZWZpeBjqs9OCASABKAlSBnByZWZpeBI+Cg5yZXF1ZXN0Y2hhcmdlZBiT0Ku5ASABKA'
    '4yEi5zMy5SZXF1ZXN0Q2hhcmdlZFIOcmVxdWVzdGNoYXJnZWQSIgoKc3RhcnRhZnRlchi46vuh'
    'ASABKAlSCnN0YXJ0YWZ0ZXI=');

@$core.Deprecated('Use listObjectsV2RequestDescriptor instead')
const ListObjectsV2Request$json = {
  '1': 'ListObjectsV2Request',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'continuationtoken',
      '3': 286270824,
      '4': 1,
      '5': 9,
      '10': 'continuationtoken'
    },
    {'1': 'delimiter', '3': 302132379, '4': 1, '5': 9, '10': 'delimiter'},
    {
      '1': 'encodingtype',
      '3': 532628025,
      '4': 1,
      '5': 14,
      '6': '.s3.EncodingType',
      '10': 'encodingtype'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'fetchowner', '3': 402910551, '4': 1, '5': 8, '10': 'fetchowner'},
    {'1': 'maxkeys', '3': 247655034, '4': 1, '5': 5, '10': 'maxkeys'},
    {
      '1': 'optionalobjectattributes',
      '3': 147719620,
      '4': 3,
      '5': 14,
      '6': '.s3.OptionalObjectAttributes',
      '10': 'optionalobjectattributes'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'startafter', '3': 339670328, '4': 1, '5': 9, '10': 'startafter'},
  ],
};

/// Descriptor for `ListObjectsV2Request`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listObjectsV2RequestDescriptor = $convert.base64Decode(
    'ChRMaXN0T2JqZWN0c1YyUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIwChFjb2'
    '50aW51YXRpb250b2tlbhjoysCIASABKAlSEWNvbnRpbnVhdGlvbnRva2VuEiAKCWRlbGltaXRl'
    'chib2YiQASABKAlSCWRlbGltaXRlchI4CgxlbmNvZGluZ3R5cGUYuYT9/QEgASgOMhAuczMuRW'
    '5jb2RpbmdUeXBlUgxlbmNvZGluZ3R5cGUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEo'
    'CVITZXhwZWN0ZWRidWNrZXRvd25lchIiCgpmZXRjaG93bmVyGNfaj8ABIAEoCFIKZmV0Y2hvd2'
    '5lchIbCgdtYXhrZXlzGPrUi3YgASgFUgdtYXhrZXlzElsKGG9wdGlvbmFsb2JqZWN0YXR0cmli'
    'dXRlcxjEi7hGIAMoDjIcLnMzLk9wdGlvbmFsT2JqZWN0QXR0cmlidXRlc1IYb3B0aW9uYWxvYm'
    'plY3RhdHRyaWJ1dGVzEhoKBnByZWZpeBjqs9OCASABKAlSBnByZWZpeBI4CgxyZXF1ZXN0cGF5'
    'ZXIYpObh9QEgASgOMhAuczMuUmVxdWVzdFBheWVyUgxyZXF1ZXN0cGF5ZXISIgoKc3RhcnRhZn'
    'Rlchi46vuhASABKAlSCnN0YXJ0YWZ0ZXI=');

@$core.Deprecated('Use listPartsOutputDescriptor instead')
const ListPartsOutput$json = {
  '1': 'ListPartsOutput',
  '2': [
    {'1': 'abortdate', '3': 232475318, '4': 1, '5': 9, '10': 'abortdate'},
    {'1': 'abortruleid', '3': 232462739, '4': 1, '5': 9, '10': 'abortruleid'},
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {
      '1': 'initiator',
      '3': 414951451,
      '4': 1,
      '5': 11,
      '6': '.s3.Initiator',
      '10': 'initiator'
    },
    {'1': 'istruncated', '3': 242094042, '4': 1, '5': 8, '10': 'istruncated'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'maxparts', '3': 398005914, '4': 1, '5': 5, '10': 'maxparts'},
    {
      '1': 'nextpartnumbermarker',
      '3': 28931219,
      '4': 1,
      '5': 9,
      '10': 'nextpartnumbermarker'
    },
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
    {
      '1': 'partnumbermarker',
      '3': 376535672,
      '4': 1,
      '5': 9,
      '10': 'partnumbermarker'
    },
    {
      '1': 'parts',
      '3': 213028806,
      '4': 3,
      '5': 11,
      '6': '.s3.Part',
      '10': 'parts'
    },
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'uploadid', '3': 449040722, '4': 1, '5': 9, '10': 'uploadid'},
  ],
};

/// Descriptor for `ListPartsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPartsOutputDescriptor = $convert.base64Decode(
    'Cg9MaXN0UGFydHNPdXRwdXQSHwoJYWJvcnRkYXRlGLaV7W4gASgJUglhYm9ydGRhdGUSIwoLYW'
    'JvcnRydWxlaWQYk7PsbiABKAlSC2Fib3J0cnVsZWlkEhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVj'
    'a2V0EkYKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMuQ2hlY2tzdW1BbGdvcml0aG'
    '1SEWNoZWNrc3VtYWxnb3JpdGhtEjcKDGNoZWNrc3VtdHlwZRjDvtkuIAEoDjIQLnMzLkNoZWNr'
    'c3VtVHlwZVIMY2hlY2tzdW10eXBlEi8KCWluaXRpYXRvchib0O7FASABKAsyDS5zMy5Jbml0aW'
    'F0b3JSCWluaXRpYXRvchIjCgtpc3RydW5jYXRlZBjan7hzIAEoCFILaXN0cnVuY2F0ZWQSEwoD'
    'a2V5GI2S62ggASgJUgNrZXkSHgoIbWF4cGFydHMYmq3kvQEgASgFUghtYXhwYXJ0cxI1ChRuZX'
    'h0cGFydG51bWJlcm1hcmtlchiT6eUNIAEoCVIUbmV4dHBhcnRudW1iZXJtYXJrZXISIwoFb3du'
    'ZXIY9fyK2QEgASgLMgkuczMuT3duZXJSBW93bmVyEi4KEHBhcnRudW1iZXJtYXJrZXIY+PTFsw'
    'EgASgJUhBwYXJ0bnVtYmVybWFya2VyEiEKBXBhcnRzGMafymUgAygLMgguczMuUGFydFIFcGFy'
    'dHMSPgoOcmVxdWVzdGNoYXJnZWQYk9CruQEgASgOMhIuczMuUmVxdWVzdENoYXJnZWRSDnJlcX'
    'Vlc3RjaGFyZ2VkEjgKDHN0b3JhZ2VjbGFzcxjHiMS7ASABKA4yEC5zMy5TdG9yYWdlQ2xhc3NS'
    'DHN0b3JhZ2VjbGFzcxIeCgh1cGxvYWRpZBjSoo/WASABKAlSCHVwbG9hZGlk');

@$core.Deprecated('Use listPartsRequestDescriptor instead')
const ListPartsRequest$json = {
  '1': 'ListPartsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'maxparts', '3': 398005914, '4': 1, '5': 5, '10': 'maxparts'},
    {
      '1': 'partnumbermarker',
      '3': 376535672,
      '4': 1,
      '5': 9,
      '10': 'partnumbermarker'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'uploadid', '3': 449040722, '4': 1, '5': 9, '10': 'uploadid'},
  ],
};

/// Descriptor for `ListPartsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPartsRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0UGFydHNSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EjMKE2V4cGVjdG'
    'VkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXISEwoDa2V5GI2S62gg'
    'ASgJUgNrZXkSHgoIbWF4cGFydHMYmq3kvQEgASgFUghtYXhwYXJ0cxIuChBwYXJ0bnVtYmVybW'
    'Fya2VyGPj0xbMBIAEoCVIQcGFydG51bWJlcm1hcmtlchI4CgxyZXF1ZXN0cGF5ZXIYpObh9QEg'
    'ASgOMhAuczMuUmVxdWVzdFBheWVyUgxyZXF1ZXN0cGF5ZXISNQoUc3NlY3VzdG9tZXJhbGdvcm'
    'l0aG0Y0MmBKyABKAlSFHNzZWN1c3RvbWVyYWxnb3JpdGhtEikKDnNzZWN1c3RvbWVya2V5GJr+'
    '9DsgASgJUg5zc2VjdXN0b21lcmtleRIuChFzc2VjdXN0b21lcmtleW1kNRjo0RcgASgJUhFzc2'
    'VjdXN0b21lcmtleW1kNRIeCgh1cGxvYWRpZBjSoo/WASABKAlSCHVwbG9hZGlk');

@$core.Deprecated('Use locationInfoDescriptor instead')
const LocationInfo$json = {
  '1': 'LocationInfo',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.s3.LocationType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `LocationInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List locationInfoDescriptor = $convert.base64Decode(
    'CgxMb2NhdGlvbkluZm8SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIoCgR0eXBlGO6g14oBIAEoDj'
    'IQLnMzLkxvY2F0aW9uVHlwZVIEdHlwZQ==');

@$core.Deprecated('Use loggingEnabledDescriptor instead')
const LoggingEnabled$json = {
  '1': 'LoggingEnabled',
  '2': [
    {'1': 'targetbucket', '3': 385492919, '4': 1, '5': 9, '10': 'targetbucket'},
    {
      '1': 'targetgrants',
      '3': 211776104,
      '4': 3,
      '5': 11,
      '6': '.s3.TargetGrant',
      '10': 'targetgrants'
    },
    {
      '1': 'targetobjectkeyformat',
      '3': 368711798,
      '4': 1,
      '5': 11,
      '6': '.s3.TargetObjectKeyFormat',
      '10': 'targetobjectkeyformat'
    },
    {'1': 'targetprefix', '3': 214607973, '4': 1, '5': 9, '10': 'targetprefix'},
  ],
};

/// Descriptor for `LoggingEnabled`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loggingEnabledDescriptor = $convert.base64Decode(
    'Cg5Mb2dnaW5nRW5hYmxlZBImCgx0YXJnZXRidWNrZXQYt8/otwEgASgJUgx0YXJnZXRidWNrZX'
    'QSNgoMdGFyZ2V0Z3JhbnRzGOjk/WQgAygLMg8uczMuVGFyZ2V0R3JhbnRSDHRhcmdldGdyYW50'
    'cxJTChV0YXJnZXRvYmplY3RrZXlmb3JtYXQY9rDorwEgASgLMhkuczMuVGFyZ2V0T2JqZWN0S2'
    'V5Rm9ybWF0UhV0YXJnZXRvYmplY3RrZXlmb3JtYXQSJQoMdGFyZ2V0cHJlZml4GOXQqmYgASgJ'
    'Ugx0YXJnZXRwcmVmaXg=');

@$core.Deprecated('Use metadataConfigurationDescriptor instead')
const MetadataConfiguration$json = {
  '1': 'MetadataConfiguration',
  '2': [
    {
      '1': 'inventorytableconfiguration',
      '3': 82018446,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryTableConfiguration',
      '10': 'inventorytableconfiguration'
    },
    {
      '1': 'journaltableconfiguration',
      '3': 70721911,
      '4': 1,
      '5': 11,
      '6': '.s3.JournalTableConfiguration',
      '10': 'journaltableconfiguration'
    },
  ],
};

/// Descriptor for `MetadataConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataConfigurationDescriptor = $convert.base64Decode(
    'ChVNZXRhZGF0YUNvbmZpZ3VyYXRpb24SZAobaW52ZW50b3J5dGFibGVjb25maWd1cmF0aW9uGI'
    '6BjicgASgLMh8uczMuSW52ZW50b3J5VGFibGVDb25maWd1cmF0aW9uUhtpbnZlbnRvcnl0YWJs'
    'ZWNvbmZpZ3VyYXRpb24SXgoZam91cm5hbHRhYmxlY29uZmlndXJhdGlvbhj3wtwhIAEoCzIdLn'
    'MzLkpvdXJuYWxUYWJsZUNvbmZpZ3VyYXRpb25SGWpvdXJuYWx0YWJsZWNvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use metadataConfigurationResultDescriptor instead')
const MetadataConfigurationResult$json = {
  '1': 'MetadataConfigurationResult',
  '2': [
    {
      '1': 'destinationresult',
      '3': 362787975,
      '4': 1,
      '5': 11,
      '6': '.s3.DestinationResult',
      '10': 'destinationresult'
    },
    {
      '1': 'inventorytableconfigurationresult',
      '3': 288024093,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryTableConfigurationResult',
      '10': 'inventorytableconfigurationresult'
    },
    {
      '1': 'journaltableconfigurationresult',
      '3': 389355872,
      '4': 1,
      '5': 11,
      '6': '.s3.JournalTableConfigurationResult',
      '10': 'journaltableconfigurationresult'
    },
  ],
};

/// Descriptor for `MetadataConfigurationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataConfigurationResultDescriptor = $convert.base64Decode(
    'ChtNZXRhZGF0YUNvbmZpZ3VyYXRpb25SZXN1bHQSRwoRZGVzdGluYXRpb25yZXN1bHQYh+n+rA'
    'EgASgLMhUuczMuRGVzdGluYXRpb25SZXN1bHRSEWRlc3RpbmF0aW9ucmVzdWx0EncKIWludmVu'
    'dG9yeXRhYmxlY29uZmlndXJhdGlvbnJlc3VsdBidzKuJASABKAsyJS5zMy5JbnZlbnRvcnlUYW'
    'JsZUNvbmZpZ3VyYXRpb25SZXN1bHRSIWludmVudG9yeXRhYmxlY29uZmlndXJhdGlvbnJlc3Vs'
    'dBJxCh9qb3VybmFsdGFibGVjb25maWd1cmF0aW9ucmVzdWx0GOCy1LkBIAEoCzIjLnMzLkpvdX'
    'JuYWxUYWJsZUNvbmZpZ3VyYXRpb25SZXN1bHRSH2pvdXJuYWx0YWJsZWNvbmZpZ3VyYXRpb25y'
    'ZXN1bHQ=');

@$core.Deprecated('Use metadataEntryDescriptor instead')
const MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `MetadataEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataEntryDescriptor = $convert.base64Decode(
    'Cg1NZXRhZGF0YUVudHJ5EhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGAoFdmFsdWUY6/KfigEgAS'
    'gJUgV2YWx1ZQ==');

@$core.Deprecated('Use metadataTableConfigurationDescriptor instead')
const MetadataTableConfiguration$json = {
  '1': 'MetadataTableConfiguration',
  '2': [
    {
      '1': 's3tablesdestination',
      '3': 398468423,
      '4': 1,
      '5': 11,
      '6': '.s3.S3TablesDestination',
      '10': 's3tablesdestination'
    },
  ],
};

/// Descriptor for `MetadataTableConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataTableConfigurationDescriptor =
    $convert.base64Decode(
        'ChpNZXRhZGF0YVRhYmxlQ29uZmlndXJhdGlvbhJNChNzM3RhYmxlc2Rlc3RpbmF0aW9uGMfKgL'
        '4BIAEoCzIXLnMzLlMzVGFibGVzRGVzdGluYXRpb25SE3MzdGFibGVzZGVzdGluYXRpb24=');

@$core.Deprecated('Use metadataTableConfigurationResultDescriptor instead')
const MetadataTableConfigurationResult$json = {
  '1': 'MetadataTableConfigurationResult',
  '2': [
    {
      '1': 's3tablesdestinationresult',
      '3': 310170672,
      '4': 1,
      '5': 11,
      '6': '.s3.S3TablesDestinationResult',
      '10': 's3tablesdestinationresult'
    },
  ],
};

/// Descriptor for `MetadataTableConfigurationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataTableConfigurationResultDescriptor =
    $convert.base64Decode(
        'CiBNZXRhZGF0YVRhYmxlQ29uZmlndXJhdGlvblJlc3VsdBJfChlzM3RhYmxlc2Rlc3RpbmF0aW'
        '9ucmVzdWx0GLCo85MBIAEoCzIdLnMzLlMzVGFibGVzRGVzdGluYXRpb25SZXN1bHRSGXMzdGFi'
        'bGVzZGVzdGluYXRpb25yZXN1bHQ=');

@$core.Deprecated('Use metadataTableEncryptionConfigurationDescriptor instead')
const MetadataTableEncryptionConfiguration$json = {
  '1': 'MetadataTableEncryptionConfiguration',
  '2': [
    {'1': 'kmskeyarn', '3': 110041649, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {
      '1': 'ssealgorithm',
      '3': 290263248,
      '4': 1,
      '5': 14,
      '6': '.s3.TableSseAlgorithm',
      '10': 'ssealgorithm'
    },
  ],
};

/// Descriptor for `MetadataTableEncryptionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataTableEncryptionConfigurationDescriptor =
    $convert.base64Decode(
        'CiRNZXRhZGF0YVRhYmxlRW5jcnlwdGlvbkNvbmZpZ3VyYXRpb24SHwoJa21za2V5YXJuGLG0vD'
        'QgASgJUglrbXNrZXlhcm4SPQoMc3NlYWxnb3JpdGhtGNChtIoBIAEoDjIVLnMzLlRhYmxlU3Nl'
        'QWxnb3JpdGhtUgxzc2VhbGdvcml0aG0=');

@$core.Deprecated('Use metricsDescriptor instead')
const Metrics$json = {
  '1': 'Metrics',
  '2': [
    {
      '1': 'eventthreshold',
      '3': 178842779,
      '4': 1,
      '5': 11,
      '6': '.s3.ReplicationTimeValue',
      '10': 'eventthreshold'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.MetricsStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `Metrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricsDescriptor = $convert.base64Decode(
    'CgdNZXRyaWNzEkMKDmV2ZW50dGhyZXNob2xkGJvZo1UgASgLMhguczMuUmVwbGljYXRpb25UaW'
    '1lVmFsdWVSDmV2ZW50dGhyZXNob2xkEiwKBnN0YXR1cxiQ5PsCIAEoDjIRLnMzLk1ldHJpY3NT'
    'dGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use metricsAndOperatorDescriptor instead')
const MetricsAndOperator$json = {
  '1': 'MetricsAndOperator',
  '2': [
    {
      '1': 'accesspointarn',
      '3': 211889319,
      '4': 1,
      '5': 9,
      '10': 'accesspointarn'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `MetricsAndOperator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricsAndOperatorDescriptor = $convert.base64Decode(
    'ChJNZXRyaWNzQW5kT3BlcmF0b3ISKQoOYWNjZXNzcG9pbnRhcm4Yp9mEZSABKAlSDmFjY2Vzc3'
    'BvaW50YXJuEhoKBnByZWZpeBjqs9OCASABKAlSBnByZWZpeBIfCgR0YWdzGMHB9rUBIAMoCzIH'
    'LnMzLlRhZ1IEdGFncw==');

@$core.Deprecated('Use metricsConfigurationDescriptor instead')
const MetricsConfiguration$json = {
  '1': 'MetricsConfiguration',
  '2': [
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.MetricsFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `MetricsConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricsConfigurationDescriptor = $convert.base64Decode(
    'ChRNZXRyaWNzQ29uZmlndXJhdGlvbhItCgZmaWx0ZXIYmIGnpQEgASgLMhEuczMuTWV0cmljc0'
    'ZpbHRlclIGZmlsdGVyEhIKAmlkGIHyorcBIAEoCVICaWQ=');

@$core.Deprecated('Use metricsFilterDescriptor instead')
const MetricsFilter$json = {
  '1': 'MetricsFilter',
  '2': [
    {
      '1': 'accesspointarn',
      '3': 211889319,
      '4': 1,
      '5': 9,
      '10': 'accesspointarn'
    },
    {
      '1': 'and',
      '3': 297135431,
      '4': 1,
      '5': 11,
      '6': '.s3.MetricsAndOperator',
      '10': 'and'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {'1': 'tag', '3': 411259956, '4': 1, '5': 11, '6': '.s3.Tag', '10': 'tag'},
  ],
};

/// Descriptor for `MetricsFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricsFilterDescriptor = $convert.base64Decode(
    'Cg1NZXRyaWNzRmlsdGVyEikKDmFjY2Vzc3BvaW50YXJuGKfZhGUgASgJUg5hY2Nlc3Nwb2ludG'
    'FybhIsCgNhbmQYx9rXjQEgASgLMhYuczMuTWV0cmljc0FuZE9wZXJhdG9yUgNhbmQSGgoGcHJl'
    'Zml4GOqz04IBIAEoCVIGcHJlZml4Eh0KA3RhZxi0qI3EASABKAsyBy5zMy5UYWdSA3RhZw==');

@$core.Deprecated('Use multipartUploadDescriptor instead')
const MultipartUpload$json = {
  '1': 'MultipartUpload',
  '2': [
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {'1': 'initiated', '3': 114631595, '4': 1, '5': 9, '10': 'initiated'},
    {
      '1': 'initiator',
      '3': 414951451,
      '4': 1,
      '5': 11,
      '6': '.s3.Initiator',
      '10': 'initiator'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'uploadid', '3': 449040722, '4': 1, '5': 9, '10': 'uploadid'},
  ],
};

/// Descriptor for `MultipartUpload`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List multipartUploadDescriptor = $convert.base64Decode(
    'Cg9NdWx0aXBhcnRVcGxvYWQSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaG'
    'Vja3N1bUFsZ29yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SNwoMY2hlY2tzdW10eXBlGMO+2S4g'
    'ASgOMhAuczMuQ2hlY2tzdW1UeXBlUgxjaGVja3N1bXR5cGUSHwoJaW5pdGlhdGVkGKvH1DYgAS'
    'gJUglpbml0aWF0ZWQSLwoJaW5pdGlhdG9yGJvQ7sUBIAEoCzINLnMzLkluaXRpYXRvclIJaW5p'
    'dGlhdG9yEhMKA2tleRiNkutoIAEoCVIDa2V5EiMKBW93bmVyGPX8itkBIAEoCzIJLnMzLk93bm'
    'VyUgVvd25lchI4CgxzdG9yYWdlY2xhc3MYx4jEuwEgASgOMhAuczMuU3RvcmFnZUNsYXNzUgxz'
    'dG9yYWdlY2xhc3MSHgoIdXBsb2FkaWQY0qKP1gEgASgJUgh1cGxvYWRpZA==');

@$core.Deprecated('Use noSuchBucketDescriptor instead')
const NoSuchBucket$json = {
  '1': 'NoSuchBucket',
};

/// Descriptor for `NoSuchBucket`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchBucketDescriptor =
    $convert.base64Decode('CgxOb1N1Y2hCdWNrZXQ=');

@$core.Deprecated('Use noSuchKeyDescriptor instead')
const NoSuchKey$json = {
  '1': 'NoSuchKey',
};

/// Descriptor for `NoSuchKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchKeyDescriptor =
    $convert.base64Decode('CglOb1N1Y2hLZXk=');

@$core.Deprecated('Use noSuchUploadDescriptor instead')
const NoSuchUpload$json = {
  '1': 'NoSuchUpload',
};

/// Descriptor for `NoSuchUpload`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noSuchUploadDescriptor =
    $convert.base64Decode('CgxOb1N1Y2hVcGxvYWQ=');

@$core.Deprecated('Use noncurrentVersionExpirationDescriptor instead')
const NoncurrentVersionExpiration$json = {
  '1': 'NoncurrentVersionExpiration',
  '2': [
    {
      '1': 'newernoncurrentversions',
      '3': 488220096,
      '4': 1,
      '5': 5,
      '10': 'newernoncurrentversions'
    },
    {
      '1': 'noncurrentdays',
      '3': 409704199,
      '4': 1,
      '5': 5,
      '10': 'noncurrentdays'
    },
  ],
};

/// Descriptor for `NoncurrentVersionExpiration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noncurrentVersionExpirationDescriptor =
    $convert.base64Decode(
        'ChtOb25jdXJyZW50VmVyc2lvbkV4cGlyYXRpb24SPAoXbmV3ZXJub25jdXJyZW50dmVyc2lvbn'
        'MYwMvm6AEgASgFUhduZXdlcm5vbmN1cnJlbnR2ZXJzaW9ucxIqCg5ub25jdXJyZW50ZGF5cxiH'
        'rq7DASABKAVSDm5vbmN1cnJlbnRkYXlz');

@$core.Deprecated('Use noncurrentVersionTransitionDescriptor instead')
const NoncurrentVersionTransition$json = {
  '1': 'NoncurrentVersionTransition',
  '2': [
    {
      '1': 'newernoncurrentversions',
      '3': 488220096,
      '4': 1,
      '5': 5,
      '10': 'newernoncurrentversions'
    },
    {
      '1': 'noncurrentdays',
      '3': 409704199,
      '4': 1,
      '5': 5,
      '10': 'noncurrentdays'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.TransitionStorageClass',
      '10': 'storageclass'
    },
  ],
};

/// Descriptor for `NoncurrentVersionTransition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noncurrentVersionTransitionDescriptor = $convert.base64Decode(
    'ChtOb25jdXJyZW50VmVyc2lvblRyYW5zaXRpb24SPAoXbmV3ZXJub25jdXJyZW50dmVyc2lvbn'
    'MYwMvm6AEgASgFUhduZXdlcm5vbmN1cnJlbnR2ZXJzaW9ucxIqCg5ub25jdXJyZW50ZGF5cxiH'
    'rq7DASABKAVSDm5vbmN1cnJlbnRkYXlzEkIKDHN0b3JhZ2VjbGFzcxjHiMS7ASABKA4yGi5zMy'
    '5UcmFuc2l0aW9uU3RvcmFnZUNsYXNzUgxzdG9yYWdlY2xhc3M=');

@$core.Deprecated('Use notFoundDescriptor instead')
const NotFound$json = {
  '1': 'NotFound',
};

/// Descriptor for `NotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notFoundDescriptor =
    $convert.base64Decode('CghOb3RGb3VuZA==');

@$core.Deprecated('Use notificationConfigurationDescriptor instead')
const NotificationConfiguration$json = {
  '1': 'NotificationConfiguration',
  '2': [
    {
      '1': 'eventbridgeconfiguration',
      '3': 82453893,
      '4': 1,
      '5': 11,
      '6': '.s3.EventBridgeConfiguration',
      '10': 'eventbridgeconfiguration'
    },
    {
      '1': 'lambdafunctionconfigurations',
      '3': 506928806,
      '4': 3,
      '5': 11,
      '6': '.s3.LambdaFunctionConfiguration',
      '10': 'lambdafunctionconfigurations'
    },
    {
      '1': 'queueconfigurations',
      '3': 487547852,
      '4': 3,
      '5': 11,
      '6': '.s3.QueueConfiguration',
      '10': 'queueconfigurations'
    },
    {
      '1': 'topicconfigurations',
      '3': 275598696,
      '4': 3,
      '5': 11,
      '6': '.s3.TopicConfiguration',
      '10': 'topicconfigurations'
    },
  ],
};

/// Descriptor for `NotificationConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notificationConfigurationDescriptor = $convert.base64Decode(
    'ChlOb3RpZmljYXRpb25Db25maWd1cmF0aW9uElsKGGV2ZW50YnJpZGdlY29uZmlndXJhdGlvbh'
    'iFy6gnIAEoCzIcLnMzLkV2ZW50QnJpZGdlQ29uZmlndXJhdGlvblIYZXZlbnRicmlkZ2Vjb25m'
    'aWd1cmF0aW9uEmcKHGxhbWJkYWZ1bmN0aW9uY29uZmlndXJhdGlvbnMYpr3c8QEgAygLMh8ucz'
    'MuTGFtYmRhRnVuY3Rpb25Db25maWd1cmF0aW9uUhxsYW1iZGFmdW5jdGlvbmNvbmZpZ3VyYXRp'
    'b25zEkwKE3F1ZXVlY29uZmlndXJhdGlvbnMYzMe96AEgAygLMhYuczMuUXVldWVDb25maWd1cm'
    'F0aW9uUhNxdWV1ZWNvbmZpZ3VyYXRpb25zEkwKE3RvcGljY29uZmlndXJhdGlvbnMY6Jq1gwEg'
    'AygLMhYuczMuVG9waWNDb25maWd1cmF0aW9uUhN0b3BpY2NvbmZpZ3VyYXRpb25z');

@$core.Deprecated('Use notificationConfigurationFilterDescriptor instead')
const NotificationConfigurationFilter$json = {
  '1': 'NotificationConfigurationFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 11,
      '6': '.s3.S3KeyFilter',
      '10': 'key'
    },
  ],
};

/// Descriptor for `NotificationConfigurationFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notificationConfigurationFilterDescriptor =
    $convert.base64Decode(
        'Ch9Ob3RpZmljYXRpb25Db25maWd1cmF0aW9uRmlsdGVyEiQKA2tleRiNkutoIAEoCzIPLnMzLl'
        'MzS2V5RmlsdGVyUgNrZXk=');

@$core.Deprecated('Use objectDescriptor instead')
const Object$json = {
  '1': 'Object',
  '2': [
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 3,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
    {
      '1': 'restorestatus',
      '3': 456059636,
      '4': 1,
      '5': 11,
      '6': '.s3.RestoreStatus',
      '10': 'restorestatus'
    },
    {'1': 'size', '3': 105352829, '4': 1, '5': 3, '10': 'size'},
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectStorageClass',
      '10': 'storageclass'
    },
  ],
};

/// Descriptor for `Object`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectDescriptor = $convert.base64Decode(
    'CgZPYmplY3QSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiADKA4yFS5zMy5DaGVja3N1bUFsZ2'
    '9yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SNwoMY2hlY2tzdW10eXBlGMO+2S4gASgOMhAuczMu'
    'Q2hlY2tzdW1UeXBlUgxjaGVja3N1bXR5cGUSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSEwoDa2'
    'V5GI2S62ggASgJUgNrZXkSJgoMbGFzdG1vZGlmaWVkGKec/M4BIAEoCVIMbGFzdG1vZGlmaWVk'
    'EiMKBW93bmVyGPX8itkBIAEoCzIJLnMzLk93bmVyUgVvd25lchI7Cg1yZXN0b3Jlc3RhdHVzGP'
    'TVu9kBIAEoCzIRLnMzLlJlc3RvcmVTdGF0dXNSDXJlc3RvcmVzdGF0dXMSFQoEc2l6ZRj9nJ4y'
    'IAEoA1IEc2l6ZRI+CgxzdG9yYWdlY2xhc3MYx4jEuwEgASgOMhYuczMuT2JqZWN0U3RvcmFnZU'
    'NsYXNzUgxzdG9yYWdlY2xhc3M=');

@$core.Deprecated('Use objectAlreadyInActiveTierErrorDescriptor instead')
const ObjectAlreadyInActiveTierError$json = {
  '1': 'ObjectAlreadyInActiveTierError',
};

/// Descriptor for `ObjectAlreadyInActiveTierError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectAlreadyInActiveTierErrorDescriptor =
    $convert.base64Decode('Ch5PYmplY3RBbHJlYWR5SW5BY3RpdmVUaWVyRXJyb3I=');

@$core.Deprecated('Use objectEncryptionDescriptor instead')
const ObjectEncryption$json = {
  '1': 'ObjectEncryption',
  '2': [
    {
      '1': 'ssekms',
      '3': 230338916,
      '4': 1,
      '5': 11,
      '6': '.s3.SSEKMSEncryption',
      '10': 'ssekms'
    },
  ],
};

/// Descriptor for `ObjectEncryption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectEncryptionDescriptor = $convert.base64Decode(
    'ChBPYmplY3RFbmNyeXB0aW9uEi8KBnNzZWttcxjk4uptIAEoCzIULnMzLlNTRUtNU0VuY3J5cH'
    'Rpb25SBnNzZWttcw==');

@$core.Deprecated('Use objectIdentifierDescriptor instead')
const ObjectIdentifier$json = {
  '1': 'ObjectIdentifier',
  '2': [
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'size', '3': 105352829, '4': 1, '5': 3, '10': 'size'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `ObjectIdentifier`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectIdentifierDescriptor = $convert.base64Decode(
    'ChBPYmplY3RJZGVudGlmaWVyEhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEhMKA2tleRiNkutoIA'
    'EoCVIDa2V5Ei0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWUS'
    'FQoEc2l6ZRj9nJ4yIAEoA1IEc2l6ZRIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaW'
    'Q=');

@$core.Deprecated('Use objectLockConfigurationDescriptor instead')
const ObjectLockConfiguration$json = {
  '1': 'ObjectLockConfiguration',
  '2': [
    {
      '1': 'objectlockenabled',
      '3': 495741839,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockEnabled',
      '10': 'objectlockenabled'
    },
    {
      '1': 'rule',
      '3': 475696372,
      '4': 1,
      '5': 11,
      '6': '.s3.ObjectLockRule',
      '10': 'rule'
    },
  ],
};

/// Descriptor for `ObjectLockConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectLockConfigurationDescriptor = $convert.base64Decode(
    'ChdPYmplY3RMb2NrQ29uZmlndXJhdGlvbhJHChFvYmplY3Rsb2NrZW5hYmxlZBiP17HsASABKA'
    '4yFS5zMy5PYmplY3RMb2NrRW5hYmxlZFIRb2JqZWN0bG9ja2VuYWJsZWQSKgoEcnVsZRj0meri'
    'ASABKAsyEi5zMy5PYmplY3RMb2NrUnVsZVIEcnVsZQ==');

@$core.Deprecated('Use objectLockLegalHoldDescriptor instead')
const ObjectLockLegalHold$json = {
  '1': 'ObjectLockLegalHold',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockLegalHoldStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ObjectLockLegalHold`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectLockLegalHoldDescriptor = $convert.base64Decode(
    'ChNPYmplY3RMb2NrTGVnYWxIb2xkEjgKBnN0YXR1cxiQ5PsCIAEoDjIdLnMzLk9iamVjdExvY2'
    'tMZWdhbEhvbGRTdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use objectLockRetentionDescriptor instead')
const ObjectLockRetention$json = {
  '1': 'ObjectLockRetention',
  '2': [
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockRetentionMode',
      '10': 'mode'
    },
    {
      '1': 'retainuntildate',
      '3': 252881225,
      '4': 1,
      '5': 9,
      '10': 'retainuntildate'
    },
  ],
};

/// Descriptor for `ObjectLockRetention`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectLockRetentionDescriptor = $convert.base64Decode(
    'ChNPYmplY3RMb2NrUmV0ZW50aW9uEjMKBG1vZGUYs+65mgEgASgOMhsuczMuT2JqZWN0TG9ja1'
    'JldGVudGlvbk1vZGVSBG1vZGUSKwoPcmV0YWludW50aWxkYXRlGMnSynggASgJUg9yZXRhaW51'
    'bnRpbGRhdGU=');

@$core.Deprecated('Use objectLockRuleDescriptor instead')
const ObjectLockRule$json = {
  '1': 'ObjectLockRule',
  '2': [
    {
      '1': 'defaultretention',
      '3': 165021749,
      '4': 1,
      '5': 11,
      '6': '.s3.DefaultRetention',
      '10': 'defaultretention'
    },
  ],
};

/// Descriptor for `ObjectLockRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectLockRuleDescriptor = $convert.base64Decode(
    'Cg5PYmplY3RMb2NrUnVsZRJDChBkZWZhdWx0cmV0ZW50aW9uGLWQ2E4gASgLMhQuczMuRGVmYX'
    'VsdFJldGVudGlvblIQZGVmYXVsdHJldGVudGlvbg==');

@$core.Deprecated('Use objectNotInActiveTierErrorDescriptor instead')
const ObjectNotInActiveTierError$json = {
  '1': 'ObjectNotInActiveTierError',
};

/// Descriptor for `ObjectNotInActiveTierError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectNotInActiveTierErrorDescriptor =
    $convert.base64Decode('ChpPYmplY3ROb3RJbkFjdGl2ZVRpZXJFcnJvcg==');

@$core.Deprecated('Use objectPartDescriptor instead')
const ObjectPart$json = {
  '1': 'ObjectPart',
  '2': [
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {'1': 'partnumber', '3': 372082310, '4': 1, '5': 5, '10': 'partnumber'},
    {'1': 'size', '3': 105352829, '4': 1, '5': 3, '10': 'size'},
  ],
};

/// Descriptor for `ObjectPart`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectPartDescriptor = $convert.base64Decode(
    'CgpPYmplY3RQYXJ0EicKDWNoZWNrc3VtY3JjMzIYwqPNMyABKAlSDWNoZWNrc3VtY3JjMzISKQ'
    'oOY2hlY2tzdW1jcmMzMmMYp5+lTCABKAlSDmNoZWNrc3VtY3JjMzJjEi8KEWNoZWNrc3VtY3Jj'
    'NjRudm1lGM28izggASgJUhFjaGVja3N1bWNyYzY0bnZtZRImCgxjaGVja3N1bXNoYTEYxOzgig'
    'EgASgJUgxjaGVja3N1bXNoYTESKQoOY2hlY2tzdW1zaGEyNTYYvvncRCABKAlSDmNoZWNrc3Vt'
    'c2hhMjU2EiIKCnBhcnRudW1iZXIYho22sQEgASgFUgpwYXJ0bnVtYmVyEhUKBHNpemUY/ZyeMi'
    'ABKANSBHNpemU=');

@$core.Deprecated('Use objectVersionDescriptor instead')
const ObjectVersion$json = {
  '1': 'ObjectVersion',
  '2': [
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 3,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'islatest', '3': 80355831, '4': 1, '5': 8, '10': 'islatest'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.s3.Owner',
      '10': 'owner'
    },
    {
      '1': 'restorestatus',
      '3': 456059636,
      '4': 1,
      '5': 11,
      '6': '.s3.RestoreStatus',
      '10': 'restorestatus'
    },
    {'1': 'size', '3': 105352829, '4': 1, '5': 3, '10': 'size'},
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectVersionStorageClass',
      '10': 'storageclass'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `ObjectVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List objectVersionDescriptor = $convert.base64Decode(
    'Cg1PYmplY3RWZXJzaW9uEkYKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogAygOMhUuczMuQ2hlY2'
    'tzdW1BbGdvcml0aG1SEWNoZWNrc3VtYWxnb3JpdGhtEjcKDGNoZWNrc3VtdHlwZRjDvtkuIAEo'
    'DjIQLnMzLkNoZWNrc3VtVHlwZVIMY2hlY2tzdW10eXBlEhYKBGV0YWcYgd+zlQEgASgJUgRldG'
    'FnEh0KCGlzbGF0ZXN0GPfDqCYgASgIUghpc2xhdGVzdBITCgNrZXkYjZLraCABKAlSA2tleRIm'
    'CgxsYXN0bW9kaWZpZWQYp5z8zgEgASgJUgxsYXN0bW9kaWZpZWQSIwoFb3duZXIY9fyK2QEgAS'
    'gLMgkuczMuT3duZXJSBW93bmVyEjsKDXJlc3RvcmVzdGF0dXMY9NW72QEgASgLMhEuczMuUmVz'
    'dG9yZVN0YXR1c1INcmVzdG9yZXN0YXR1cxIVCgRzaXplGP2cnjIgASgDUgRzaXplEkUKDHN0b3'
    'JhZ2VjbGFzcxjHiMS7ASABKA4yHS5zMy5PYmplY3RWZXJzaW9uU3RvcmFnZUNsYXNzUgxzdG9y'
    'YWdlY2xhc3MSIAoJdmVyc2lvbmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlk');

@$core.Deprecated('Use outputLocationDescriptor instead')
const OutputLocation$json = {
  '1': 'OutputLocation',
  '2': [
    {
      '1': 's3',
      '3': 100795332,
      '4': 1,
      '5': 11,
      '6': '.s3.S3Location',
      '10': 's3'
    },
  ],
};

/// Descriptor for `OutputLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List outputLocationDescriptor = $convert.base64Decode(
    'Cg5PdXRwdXRMb2NhdGlvbhIhCgJzMxjEh4gwIAEoCzIOLnMzLlMzTG9jYXRpb25SAnMz');

@$core.Deprecated('Use outputSerializationDescriptor instead')
const OutputSerialization$json = {
  '1': 'OutputSerialization',
  '2': [
    {
      '1': 'csv',
      '3': 455189056,
      '4': 1,
      '5': 11,
      '6': '.s3.CSVOutput',
      '10': 'csv'
    },
    {
      '1': 'json',
      '3': 532512708,
      '4': 1,
      '5': 11,
      '6': '.s3.JSONOutput',
      '10': 'json'
    },
  ],
};

/// Descriptor for `OutputSerialization`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List outputSerializationDescriptor = $convert.base64Decode(
    'ChNPdXRwdXRTZXJpYWxpemF0aW9uEiMKA2NzdhjAxIbZASABKAsyDS5zMy5DU1ZPdXRwdXRSA2'
    'NzdhImCgRqc29uGMT/9f0BIAEoCzIOLnMzLkpTT05PdXRwdXRSBGpzb24=');

@$core.Deprecated('Use ownerDescriptor instead')
const Owner$json = {
  '1': 'Owner',
  '2': [
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {'1': 'id', '3': 384363361, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `Owner`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ownerDescriptor = $convert.base64Decode(
    'CgVPd25lchIkCgtkaXNwbGF5bmFtZRi3ybLHASABKAlSC2Rpc3BsYXluYW1lEhIKAmlkGOHWo7'
    'cBIAEoCVICaWQ=');

@$core.Deprecated('Use ownershipControlsDescriptor instead')
const OwnershipControls$json = {
  '1': 'OwnershipControls',
  '2': [
    {
      '1': 'rules',
      '3': 42675585,
      '4': 3,
      '5': 11,
      '6': '.s3.OwnershipControlsRule',
      '10': 'rules'
    },
  ],
};

/// Descriptor for `OwnershipControls`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ownershipControlsDescriptor = $convert.base64Decode(
    'ChFPd25lcnNoaXBDb250cm9scxIyCgVydWxlcxiB26wUIAMoCzIZLnMzLk93bmVyc2hpcENvbn'
    'Ryb2xzUnVsZVIFcnVsZXM=');

@$core.Deprecated('Use ownershipControlsRuleDescriptor instead')
const OwnershipControlsRule$json = {
  '1': 'OwnershipControlsRule',
  '2': [
    {
      '1': 'objectownership',
      '3': 448301184,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectOwnership',
      '10': 'objectownership'
    },
  ],
};

/// Descriptor for `OwnershipControlsRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ownershipControlsRuleDescriptor = $convert.base64Decode(
    'ChVPd25lcnNoaXBDb250cm9sc1J1bGUSQQoPb2JqZWN0b3duZXJzaGlwGICR4tUBIAEoDjITLn'
    'MzLk9iamVjdE93bmVyc2hpcFIPb2JqZWN0b3duZXJzaGlw');

@$core.Deprecated('Use parquetInputDescriptor instead')
const ParquetInput$json = {
  '1': 'ParquetInput',
};

/// Descriptor for `ParquetInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parquetInputDescriptor =
    $convert.base64Decode('CgxQYXJxdWV0SW5wdXQ=');

@$core.Deprecated('Use partDescriptor instead')
const Part$json = {
  '1': 'Part',
  '2': [
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {'1': 'partnumber', '3': 372082310, '4': 1, '5': 5, '10': 'partnumber'},
    {'1': 'size', '3': 105352829, '4': 1, '5': 3, '10': 'size'},
  ],
};

/// Descriptor for `Part`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List partDescriptor = $convert.base64Decode(
    'CgRQYXJ0EicKDWNoZWNrc3VtY3JjMzIYwqPNMyABKAlSDWNoZWNrc3VtY3JjMzISKQoOY2hlY2'
    'tzdW1jcmMzMmMYp5+lTCABKAlSDmNoZWNrc3VtY3JjMzJjEi8KEWNoZWNrc3VtY3JjNjRudm1l'
    'GM28izggASgJUhFjaGVja3N1bWNyYzY0bnZtZRImCgxjaGVja3N1bXNoYTEYxOzgigEgASgJUg'
    'xjaGVja3N1bXNoYTESKQoOY2hlY2tzdW1zaGEyNTYYvvncRCABKAlSDmNoZWNrc3Vtc2hhMjU2'
    'EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEiYKDGxhc3Rtb2RpZmllZBinnPzOASABKAlSDGxhc3'
    'Rtb2RpZmllZBIiCgpwYXJ0bnVtYmVyGIaNtrEBIAEoBVIKcGFydG51bWJlchIVCgRzaXplGP2c'
    'njIgASgDUgRzaXpl');

@$core.Deprecated('Use partitionedPrefixDescriptor instead')
const PartitionedPrefix$json = {
  '1': 'PartitionedPrefix',
  '2': [
    {
      '1': 'partitiondatesource',
      '3': 127339161,
      '4': 1,
      '5': 14,
      '6': '.s3.PartitionDateSource',
      '10': 'partitiondatesource'
    },
  ],
};

/// Descriptor for `PartitionedPrefix`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List partitionedPrefixDescriptor = $convert.base64Decode(
    'ChFQYXJ0aXRpb25lZFByZWZpeBJMChNwYXJ0aXRpb25kYXRlc291cmNlGJmV3DwgASgOMhcucz'
    'MuUGFydGl0aW9uRGF0ZVNvdXJjZVITcGFydGl0aW9uZGF0ZXNvdXJjZQ==');

@$core.Deprecated('Use policyStatusDescriptor instead')
const PolicyStatus$json = {
  '1': 'PolicyStatus',
  '2': [
    {'1': 'ispublic', '3': 274891537, '4': 1, '5': 8, '10': 'ispublic'},
  ],
};

/// Descriptor for `PolicyStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyStatusDescriptor = $convert.base64Decode(
    'CgxQb2xpY3lTdGF0dXMSHgoIaXNwdWJsaWMYkYaKgwEgASgIUghpc3B1YmxpYw==');

@$core.Deprecated('Use progressDescriptor instead')
const Progress$json = {
  '1': 'Progress',
  '2': [
    {
      '1': 'bytesprocessed',
      '3': 487068657,
      '4': 1,
      '5': 3,
      '10': 'bytesprocessed'
    },
    {
      '1': 'bytesreturned',
      '3': 121984684,
      '4': 1,
      '5': 3,
      '10': 'bytesreturned'
    },
    {'1': 'bytesscanned', '3': 186950329, '4': 1, '5': 3, '10': 'bytesscanned'},
  ],
};

/// Descriptor for `Progress`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List progressDescriptor = $convert.base64Decode(
    'CghQcm9ncmVzcxIqCg5ieXRlc3Byb2Nlc3NlZBjxp6DoASABKANSDmJ5dGVzcHJvY2Vzc2VkEi'
    'cKDWJ5dGVzcmV0dXJuZWQYrK2VOiABKANSDWJ5dGVzcmV0dXJuZWQSJQoMYnl0ZXNzY2FubmVk'
    'GLnFklkgASgDUgxieXRlc3NjYW5uZWQ=');

@$core.Deprecated('Use progressEventDescriptor instead')
const ProgressEvent$json = {
  '1': 'ProgressEvent',
  '2': [
    {
      '1': 'details',
      '3': 247611974,
      '4': 1,
      '5': 11,
      '6': '.s3.Progress',
      '10': 'details'
    },
  ],
};

/// Descriptor for `ProgressEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List progressEventDescriptor = $convert.base64Decode(
    'Cg1Qcm9ncmVzc0V2ZW50EikKB2RldGFpbHMYxoSJdiABKAsyDC5zMy5Qcm9ncmVzc1IHZGV0YW'
    'lscw==');

@$core.Deprecated('Use publicAccessBlockConfigurationDescriptor instead')
const PublicAccessBlockConfiguration$json = {
  '1': 'PublicAccessBlockConfiguration',
  '2': [
    {
      '1': 'blockpublicacls',
      '3': 31292983,
      '4': 1,
      '5': 8,
      '10': 'blockpublicacls'
    },
    {
      '1': 'blockpublicpolicy',
      '3': 505379326,
      '4': 1,
      '5': 8,
      '10': 'blockpublicpolicy'
    },
    {
      '1': 'ignorepublicacls',
      '3': 406560036,
      '4': 1,
      '5': 8,
      '10': 'ignorepublicacls'
    },
    {
      '1': 'restrictpublicbuckets',
      '3': 258930040,
      '4': 1,
      '5': 8,
      '10': 'restrictpublicbuckets'
    },
  ],
};

/// Descriptor for `PublicAccessBlockConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicAccessBlockConfigurationDescriptor = $convert.base64Decode(
    'Ch5QdWJsaWNBY2Nlc3NCbG9ja0NvbmZpZ3VyYXRpb24SKwoPYmxvY2twdWJsaWNhY2xzGLf89Q'
    '4gASgIUg9ibG9ja3B1YmxpY2FjbHMSMAoRYmxvY2twdWJsaWNwb2xpY3kY/vP98AEgASgIUhFi'
    'bG9ja3B1YmxpY3BvbGljeRIuChBpZ25vcmVwdWJsaWNhY2xzGKS67sEBIAEoCFIQaWdub3JlcH'
    'VibGljYWNscxI3ChVyZXN0cmljdHB1YmxpY2J1Y2tldHMY+Oq7eyABKAhSFXJlc3RyaWN0cHVi'
    'bGljYnVja2V0cw==');

@$core.Deprecated('Use putBucketAbacRequestDescriptor instead')
const PutBucketAbacRequest$json = {
  '1': 'PutBucketAbacRequest',
  '2': [
    {
      '1': 'abacstatus',
      '3': 436325625,
      '4': 1,
      '5': 11,
      '6': '.s3.AbacStatus',
      '8': {},
      '10': 'abacstatus'
    },
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `PutBucketAbacRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketAbacRequestDescriptor = $convert.base64Decode(
    'ChRQdXRCdWNrZXRBYmFjUmVxdWVzdBI4CgphYmFjc3RhdHVzGPmZh9ABIAEoCzIOLnMzLkFiYW'
    'NTdGF0dXNCBIi1GAFSCmFiYWNzdGF0dXMSGQoGYnVja2V0GNjquBogASgJUgZidWNrZXQSRgoR'
    'Y2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFsZ29yaXRobVIRY2hlY2'
    'tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbnRlbnRtZDUSMwoTZXhw'
    'ZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lcg==');

@$core
    .Deprecated('Use putBucketAccelerateConfigurationRequestDescriptor instead')
const PutBucketAccelerateConfigurationRequest$json = {
  '1': 'PutBucketAccelerateConfigurationRequest',
  '2': [
    {
      '1': 'accelerateconfiguration',
      '3': 376003075,
      '4': 1,
      '5': 11,
      '6': '.s3.AccelerateConfiguration',
      '8': {},
      '10': 'accelerateconfiguration'
    },
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `PutBucketAccelerateConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketAccelerateConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CidQdXRCdWNrZXRBY2NlbGVyYXRlQ29uZmlndXJhdGlvblJlcXVlc3QSXwoXYWNjZWxlcmF0ZW'
        'NvbmZpZ3VyYXRpb24Yg7SlswEgASgLMhsuczMuQWNjZWxlcmF0ZUNvbmZpZ3VyYXRpb25CBIi1'
        'GAFSF2FjY2VsZXJhdGVjb25maWd1cmF0aW9uEhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0Ek'
        'YKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMuQ2hlY2tzdW1BbGdvcml0aG1SEWNo'
        'ZWNrc3VtYWxnb3JpdGhtEjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdG'
        'VkYnVja2V0b3duZXI=');

@$core.Deprecated('Use putBucketAclRequestDescriptor instead')
const PutBucketAclRequest$json = {
  '1': 'PutBucketAclRequest',
  '2': [
    {
      '1': 'acl',
      '3': 394696836,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketCannedACL',
      '10': 'acl'
    },
    {
      '1': 'accesscontrolpolicy',
      '3': 302514423,
      '4': 1,
      '5': 11,
      '6': '.s3.AccessControlPolicy',
      '8': {},
      '10': 'accesscontrolpolicy'
    },
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'grantfullcontrol',
      '3': 102486874,
      '4': 1,
      '5': 9,
      '10': 'grantfullcontrol'
    },
    {'1': 'grantread', '3': 4518126, '4': 1, '5': 9, '10': 'grantread'},
    {'1': 'grantreadacp', '3': 208410960, '4': 1, '5': 9, '10': 'grantreadacp'},
    {'1': 'grantwrite', '3': 286850821, '4': 1, '5': 9, '10': 'grantwrite'},
    {
      '1': 'grantwriteacp',
      '3': 157594957,
      '4': 1,
      '5': 9,
      '10': 'grantwriteacp'
    },
  ],
};

/// Descriptor for `PutBucketAclRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketAclRequestDescriptor = $convert.base64Decode(
    'ChNQdXRCdWNrZXRBY2xSZXF1ZXN0EikKA2FjbBiEsZq8ASABKA4yEy5zMy5CdWNrZXRDYW5uZW'
    'RBQ0xSA2FjbBJTChNhY2Nlc3Njb250cm9scG9saWN5GPeBoJABIAEoCzIXLnMzLkFjY2Vzc0Nv'
    'bnRyb2xQb2xpY3lCBIi1GAFSE2FjY2Vzc2NvbnRyb2xwb2xpY3kSGQoGYnVja2V0GNjquBogAS'
    'gJUgZidWNrZXQSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFs'
    'Z29yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbn'
    'RlbnRtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRv'
    'd25lchItChBncmFudGZ1bGxjb250cm9sGNqm7zAgASgJUhBncmFudGZ1bGxjb250cm9sEh8KCW'
    'dyYW50cmVhZBju4ZMCIAEoCVIJZ3JhbnRyZWFkEiUKDGdyYW50cmVhZGFjcBjQsrBjIAEoCVIM'
    'Z3JhbnRyZWFkYWNwEiIKCmdyYW50d3JpdGUYhf7jiAEgASgJUgpncmFudHdyaXRlEicKDWdyYW'
    '50d3JpdGVhY3AYzeqSSyABKAlSDWdyYW50d3JpdGVhY3A=');

@$core
    .Deprecated('Use putBucketAnalyticsConfigurationRequestDescriptor instead')
const PutBucketAnalyticsConfigurationRequest$json = {
  '1': 'PutBucketAnalyticsConfigurationRequest',
  '2': [
    {
      '1': 'analyticsconfiguration',
      '3': 229750388,
      '4': 1,
      '5': 11,
      '6': '.s3.AnalyticsConfiguration',
      '8': {},
      '10': 'analyticsconfiguration'
    },
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `PutBucketAnalyticsConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketAnalyticsConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiZQdXRCdWNrZXRBbmFseXRpY3NDb25maWd1cmF0aW9uUmVxdWVzdBJbChZhbmFseXRpY3Njb2'
        '5maWd1cmF0aW9uGPTsxm0gASgLMhouczMuQW5hbHl0aWNzQ29uZmlndXJhdGlvbkIEiLUYAVIW'
        'YW5hbHl0aWNzY29uZmlndXJhdGlvbhIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzChNleH'
        'BlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhIKAmlkGIHy'
        'orcBIAEoCVICaWQ=');

@$core.Deprecated('Use putBucketCorsRequestDescriptor instead')
const PutBucketCorsRequest$json = {
  '1': 'PutBucketCorsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'corsconfiguration',
      '3': 359712351,
      '4': 1,
      '5': 11,
      '6': '.s3.CORSConfiguration',
      '8': {},
      '10': 'corsconfiguration'
    },
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `PutBucketCorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketCorsRequestDescriptor = $convert.base64Decode(
    'ChRQdXRCdWNrZXRDb3JzUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBJNChFjb3'
    'JzY29uZmlndXJhdGlvbhjfjMOrASABKAsyFS5zMy5DT1JTQ29uZmlndXJhdGlvbkIEiLUYAVIR'
    'Y29yc2NvbmZpZ3VyYXRpb24SRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaG'
    'Vja3N1bUFsZ29yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASAB'
    'KAlSCmNvbnRlbnRtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZW'
    'RidWNrZXRvd25lcg==');

@$core.Deprecated('Use putBucketEncryptionRequestDescriptor instead')
const PutBucketEncryptionRequest$json = {
  '1': 'PutBucketEncryptionRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'serversideencryptionconfiguration',
      '3': 314603363,
      '4': 1,
      '5': 11,
      '6': '.s3.ServerSideEncryptionConfiguration',
      '8': {},
      '10': 'serversideencryptionconfiguration'
    },
  ],
};

/// Descriptor for `PutBucketEncryptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketEncryptionRequestDescriptor = $convert.base64Decode(
    'ChpQdXRCdWNrZXRFbmNyeXB0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
    'JGChFjaGVja3N1bWFsZ29yaXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3VtQWxnb3JpdGhtUhFj'
    'aGVja3N1bWFsZ29yaXRobRIiCgpjb250ZW50bWQ1GIvq6vEBIAEoCVIKY29udGVudG1kNRIzCh'
    'NleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEn0KIXNl'
    'cnZlcnNpZGVlbmNyeXB0aW9uY29uZmlndXJhdGlvbhjj7oGWASABKAsyJS5zMy5TZXJ2ZXJTaW'
    'RlRW5jcnlwdGlvbkNvbmZpZ3VyYXRpb25CBIi1GAFSIXNlcnZlcnNpZGVlbmNyeXB0aW9uY29u'
    'ZmlndXJhdGlvbg==');

@$core.Deprecated(
    'Use putBucketIntelligentTieringConfigurationRequestDescriptor instead')
const PutBucketIntelligentTieringConfigurationRequest$json = {
  '1': 'PutBucketIntelligentTieringConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'intelligenttieringconfiguration',
      '3': 381057245,
      '4': 1,
      '5': 11,
      '6': '.s3.IntelligentTieringConfiguration',
      '8': {},
      '10': 'intelligenttieringconfiguration'
    },
  ],
};

/// Descriptor for `PutBucketIntelligentTieringConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putBucketIntelligentTieringConfigurationRequestDescriptor =
    $convert.base64Decode(
        'Ci9QdXRCdWNrZXRJbnRlbGxpZ2VudFRpZXJpbmdDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidW'
        'NrZXQY2Oq4GiABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNl'
        'eHBlY3RlZGJ1Y2tldG93bmVyEhIKAmlkGIHyorcBIAEoCVICaWQSdwofaW50ZWxsaWdlbnR0aW'
        'VyaW5nY29uZmlndXJhdGlvbhjd8dm1ASABKAsyIy5zMy5JbnRlbGxpZ2VudFRpZXJpbmdDb25m'
        'aWd1cmF0aW9uQgSItRgBUh9pbnRlbGxpZ2VudHRpZXJpbmdjb25maWd1cmF0aW9u');

@$core
    .Deprecated('Use putBucketInventoryConfigurationRequestDescriptor instead')
const PutBucketInventoryConfigurationRequest$json = {
  '1': 'PutBucketInventoryConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'inventoryconfiguration',
      '3': 401695392,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryConfiguration',
      '8': {},
      '10': 'inventoryconfiguration'
    },
  ],
};

/// Descriptor for `PutBucketInventoryConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketInventoryConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiZQdXRCdWNrZXRJbnZlbnRvcnlDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4Gi'
        'ABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1'
        'Y2tldG93bmVyEhIKAmlkGIHyorcBIAEoCVICaWQSXAoWaW52ZW50b3J5Y29uZmlndXJhdGlvbh'
        'igxcW/ASABKAsyGi5zMy5JbnZlbnRvcnlDb25maWd1cmF0aW9uQgSItRgBUhZpbnZlbnRvcnlj'
        'b25maWd1cmF0aW9u');

@$core.Deprecated('Use putBucketLifecycleConfigurationOutputDescriptor instead')
const PutBucketLifecycleConfigurationOutput$json = {
  '1': 'PutBucketLifecycleConfigurationOutput',
  '2': [
    {
      '1': 'transitiondefaultminimumobjectsize',
      '3': 81614768,
      '4': 1,
      '5': 14,
      '6': '.s3.TransitionDefaultMinimumObjectSize',
      '10': 'transitiondefaultminimumobjectsize'
    },
  ],
};

/// Descriptor for `PutBucketLifecycleConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketLifecycleConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiVQdXRCdWNrZXRMaWZlY3ljbGVDb25maWd1cmF0aW9uT3V0cHV0EnkKInRyYW5zaXRpb25kZW'
        'ZhdWx0bWluaW11bW9iamVjdHNpemUYsK/1JiABKA4yJi5zMy5UcmFuc2l0aW9uRGVmYXVsdE1p'
        'bmltdW1PYmplY3RTaXplUiJ0cmFuc2l0aW9uZGVmYXVsdG1pbmltdW1vYmplY3RzaXpl');

@$core
    .Deprecated('Use putBucketLifecycleConfigurationRequestDescriptor instead')
const PutBucketLifecycleConfigurationRequest$json = {
  '1': 'PutBucketLifecycleConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'lifecycleconfiguration',
      '3': 8468208,
      '4': 1,
      '5': 11,
      '6': '.s3.BucketLifecycleConfiguration',
      '8': {},
      '10': 'lifecycleconfiguration'
    },
    {
      '1': 'transitiondefaultminimumobjectsize',
      '3': 81614768,
      '4': 1,
      '5': 14,
      '6': '.s3.TransitionDefaultMinimumObjectSize',
      '10': 'transitiondefaultminimumobjectsize'
    },
  ],
};

/// Descriptor for `PutBucketLifecycleConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketLifecycleConfigurationRequestDescriptor = $convert.base64Decode(
    'CiZQdXRCdWNrZXRMaWZlY3ljbGVDb25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4Gi'
    'ABKAlSBmJ1Y2tldBJGChFjaGVja3N1bWFsZ29yaXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3Vt'
    'QWxnb3JpdGhtUhFjaGVja3N1bWFsZ29yaXRobRIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D'
    '4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEmEKFmxpZmVjeWNsZWNvbmZpZ3VyYXRpb24Y8O2E'
    'BCABKAsyIC5zMy5CdWNrZXRMaWZlY3ljbGVDb25maWd1cmF0aW9uQgSItRgBUhZsaWZlY3ljbG'
    'Vjb25maWd1cmF0aW9uEnkKInRyYW5zaXRpb25kZWZhdWx0bWluaW11bW9iamVjdHNpemUYsK/1'
    'JiABKA4yJi5zMy5UcmFuc2l0aW9uRGVmYXVsdE1pbmltdW1PYmplY3RTaXplUiJ0cmFuc2l0aW'
    '9uZGVmYXVsdG1pbmltdW1vYmplY3RzaXpl');

@$core.Deprecated('Use putBucketLoggingRequestDescriptor instead')
const PutBucketLoggingRequest$json = {
  '1': 'PutBucketLoggingRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bucketloggingstatus',
      '3': 289553275,
      '4': 1,
      '5': 11,
      '6': '.s3.BucketLoggingStatus',
      '8': {},
      '10': 'bucketloggingstatus'
    },
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
  ],
};

/// Descriptor for `PutBucketLoggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketLoggingRequestDescriptor = $convert.base64Decode(
    'ChdQdXRCdWNrZXRMb2dnaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBJTCh'
    'NidWNrZXRsb2dnaW5nc3RhdHVzGPv2iIoBIAEoCzIXLnMzLkJ1Y2tldExvZ2dpbmdTdGF0dXNC'
    'BIi1GAFSE2J1Y2tldGxvZ2dpbmdzdGF0dXMSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA'
    '4yFS5zMy5DaGVja3N1bUFsZ29yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1k'
    'NRiL6urxASABKAlSCmNvbnRlbnRtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCV'
    'ITZXhwZWN0ZWRidWNrZXRvd25lcg==');

@$core.Deprecated('Use putBucketMetricsConfigurationRequestDescriptor instead')
const PutBucketMetricsConfigurationRequest$json = {
  '1': 'PutBucketMetricsConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'metricsconfiguration',
      '3': 14504541,
      '4': 1,
      '5': 11,
      '6': '.s3.MetricsConfiguration',
      '8': {},
      '10': 'metricsconfiguration'
    },
  ],
};

/// Descriptor for `PutBucketMetricsConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketMetricsConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CiRQdXRCdWNrZXRNZXRyaWNzQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquBogAS'
        'gJUgZidWNrZXQSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNr'
        'ZXRvd25lchISCgJpZBiB8qK3ASABKAlSAmlkElUKFG1ldHJpY3Njb25maWd1cmF0aW9uGN2k9Q'
        'YgASgLMhguczMuTWV0cmljc0NvbmZpZ3VyYXRpb25CBIi1GAFSFG1ldHJpY3Njb25maWd1cmF0'
        'aW9u');

@$core.Deprecated(
    'Use putBucketNotificationConfigurationRequestDescriptor instead')
const PutBucketNotificationConfigurationRequest$json = {
  '1': 'PutBucketNotificationConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'notificationconfiguration',
      '3': 290208045,
      '4': 1,
      '5': 11,
      '6': '.s3.NotificationConfiguration',
      '8': {},
      '10': 'notificationconfiguration'
    },
    {
      '1': 'skipdestinationvalidation',
      '3': 114697590,
      '4': 1,
      '5': 8,
      '10': 'skipdestinationvalidation'
    },
  ],
};

/// Descriptor for `PutBucketNotificationConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putBucketNotificationConfigurationRequestDescriptor = $convert.base64Decode(
        'CilQdXRCdWNrZXROb3RpZmljYXRpb25Db25maWd1cmF0aW9uUmVxdWVzdBIZCgZidWNrZXQY2O'
        'q4GiABKAlSBmJ1Y2tldBIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3Rl'
        'ZGJ1Y2tldG93bmVyEmUKGW5vdGlmaWNhdGlvbmNvbmZpZ3VyYXRpb24YrfKwigEgASgLMh0ucz'
        'MuTm90aWZpY2F0aW9uQ29uZmlndXJhdGlvbkIEiLUYAVIZbm90aWZpY2F0aW9uY29uZmlndXJh'
        'dGlvbhI/Chlza2lwZGVzdGluYXRpb252YWxpZGF0aW9uGPbK2DYgASgIUhlza2lwZGVzdGluYX'
        'Rpb252YWxpZGF0aW9u');

@$core.Deprecated('Use putBucketOwnershipControlsRequestDescriptor instead')
const PutBucketOwnershipControlsRequest$json = {
  '1': 'PutBucketOwnershipControlsRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'ownershipcontrols',
      '3': 6169057,
      '4': 1,
      '5': 11,
      '6': '.s3.OwnershipControls',
      '8': {},
      '10': 'ownershipcontrols'
    },
  ],
};

/// Descriptor for `PutBucketOwnershipControlsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketOwnershipControlsRequestDescriptor = $convert.base64Decode(
    'CiFQdXRCdWNrZXRPd25lcnNoaXBDb250cm9sc1JlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUg'
    'ZidWNrZXQSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFsZ29y'
    'aXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbnRlbn'
    'RtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25l'
    'chJMChFvd25lcnNoaXBjb250cm9scxjhw/gCIAEoCzIVLnMzLk93bmVyc2hpcENvbnRyb2xzQg'
    'SItRgBUhFvd25lcnNoaXBjb250cm9scw==');

@$core.Deprecated('Use putBucketPolicyRequestDescriptor instead')
const PutBucketPolicyRequest$json = {
  '1': 'PutBucketPolicyRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'confirmremoveselfbucketaccess',
      '3': 227046446,
      '4': 1,
      '5': 8,
      '10': 'confirmremoveselfbucketaccess'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '8': {}, '10': 'policy'},
  ],
};

/// Descriptor for `PutBucketPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketPolicyRequestDescriptor = $convert.base64Decode(
    'ChZQdXRCdWNrZXRQb2xpY3lSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EkYKEW'
    'NoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMuQ2hlY2tzdW1BbGdvcml0aG1SEWNoZWNr'
    'c3VtYWxnb3JpdGhtEkcKHWNvbmZpcm1yZW1vdmVzZWxmYnVja2V0YWNjZXNzGK7ooWwgASgIUh'
    '1jb25maXJtcmVtb3Zlc2VsZmJ1Y2tldGFjY2VzcxIiCgpjb250ZW50bWQ1GIvq6vEBIAEoCVIK'
    'Y29udGVudG1kNRIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2'
    'tldG93bmVyEiAKBnBvbGljeRig7/DgASABKAlCBIi1GAFSBnBvbGljeQ==');

@$core.Deprecated('Use putBucketReplicationRequestDescriptor instead')
const PutBucketReplicationRequest$json = {
  '1': 'PutBucketReplicationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'replicationconfiguration',
      '3': 115861462,
      '4': 1,
      '5': 11,
      '6': '.s3.ReplicationConfiguration',
      '8': {},
      '10': 'replicationconfiguration'
    },
    {'1': 'token', '3': 439704531, '4': 1, '5': 9, '10': 'token'},
  ],
};

/// Descriptor for `PutBucketReplicationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketReplicationRequestDescriptor = $convert.base64Decode(
    'ChtQdXRCdWNrZXRSZXBsaWNhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZX'
    'QSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFsZ29yaXRobVIR'
    'Y2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbnRlbnRtZDUSMw'
    'oTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchJhChhy'
    'ZXBsaWNhdGlvbmNvbmZpZ3VyYXRpb24Y1s+fNyABKAsyHC5zMy5SZXBsaWNhdGlvbkNvbmZpZ3'
    'VyYXRpb25CBIi1GAFSGHJlcGxpY2F0aW9uY29uZmlndXJhdGlvbhIYCgV0b2tlbhjTt9XRASAB'
    'KAlSBXRva2Vu');

@$core.Deprecated('Use putBucketRequestPaymentRequestDescriptor instead')
const PutBucketRequestPaymentRequest$json = {
  '1': 'PutBucketRequestPaymentRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'requestpaymentconfiguration',
      '3': 456139707,
      '4': 1,
      '5': 11,
      '6': '.s3.RequestPaymentConfiguration',
      '8': {},
      '10': 'requestpaymentconfiguration'
    },
  ],
};

/// Descriptor for `PutBucketRequestPaymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketRequestPaymentRequestDescriptor = $convert.base64Decode(
    'Ch5QdXRCdWNrZXRSZXF1ZXN0UGF5bWVudFJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidW'
    'NrZXQSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFsZ29yaXRo'
    'bVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbnRlbnRtZD'
    'USMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchJr'
    'ChtyZXF1ZXN0cGF5bWVudGNvbmZpZ3VyYXRpb24Yu8fA2QEgASgLMh8uczMuUmVxdWVzdFBheW'
    '1lbnRDb25maWd1cmF0aW9uQgSItRgBUhtyZXF1ZXN0cGF5bWVudGNvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use putBucketTaggingRequestDescriptor instead')
const PutBucketTaggingRequest$json = {
  '1': 'PutBucketTaggingRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'tagging',
      '3': 33436541,
      '4': 1,
      '5': 11,
      '6': '.s3.Tagging',
      '8': {},
      '10': 'tagging'
    },
  ],
};

/// Descriptor for `PutBucketTaggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketTaggingRequestDescriptor = $convert.base64Decode(
    'ChdQdXRCdWNrZXRUYWdnaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBJGCh'
    'FjaGVja3N1bWFsZ29yaXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3VtQWxnb3JpdGhtUhFjaGVj'
    'a3N1bWFsZ29yaXRobRIiCgpjb250ZW50bWQ1GIvq6vEBIAEoCVIKY29udGVudG1kNRIzChNleH'
    'BlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEi4KB3RhZ2dp'
    'bmcY/eb4DyABKAsyCy5zMy5UYWdnaW5nQgSItRgBUgd0YWdnaW5n');

@$core.Deprecated('Use putBucketVersioningRequestDescriptor instead')
const PutBucketVersioningRequest$json = {
  '1': 'PutBucketVersioningRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'mfa', '3': 325407842, '4': 1, '5': 9, '10': 'mfa'},
    {
      '1': 'versioningconfiguration',
      '3': 124078014,
      '4': 1,
      '5': 11,
      '6': '.s3.VersioningConfiguration',
      '8': {},
      '10': 'versioningconfiguration'
    },
  ],
};

/// Descriptor for `PutBucketVersioningRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketVersioningRequestDescriptor = $convert.base64Decode(
    'ChpQdXRCdWNrZXRWZXJzaW9uaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
    'JGChFjaGVja3N1bWFsZ29yaXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3VtQWxnb3JpdGhtUhFj'
    'aGVja3N1bWFsZ29yaXRobRIiCgpjb250ZW50bWQ1GIvq6vEBIAEoCVIKY29udGVudG1kNRIzCh'
    'NleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhQKA21m'
    'YRjiqJWbASABKAlSA21mYRJeChd2ZXJzaW9uaW5nY29uZmlndXJhdGlvbhi+j5U7IAEoCzIbLn'
    'MzLlZlcnNpb25pbmdDb25maWd1cmF0aW9uQgSItRgBUhd2ZXJzaW9uaW5nY29uZmlndXJhdGlv'
    'bg==');

@$core.Deprecated('Use putBucketWebsiteRequestDescriptor instead')
const PutBucketWebsiteRequest$json = {
  '1': 'PutBucketWebsiteRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'websiteconfiguration',
      '3': 505552313,
      '4': 1,
      '5': 11,
      '6': '.s3.WebsiteConfiguration',
      '8': {},
      '10': 'websiteconfiguration'
    },
  ],
};

/// Descriptor for `PutBucketWebsiteRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putBucketWebsiteRequestDescriptor = $convert.base64Decode(
    'ChdQdXRCdWNrZXRXZWJzaXRlUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBJGCh'
    'FjaGVja3N1bWFsZ29yaXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3VtQWxnb3JpdGhtUhFjaGVj'
    'a3N1bWFsZ29yaXRobRIiCgpjb250ZW50bWQ1GIvq6vEBIAEoCVIKY29udGVudG1kNRIzChNleH'
    'BlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyElYKFHdlYnNp'
    'dGVjb25maWd1cmF0aW9uGLm7iPEBIAEoCzIYLnMzLldlYnNpdGVDb25maWd1cmF0aW9uQgSItR'
    'gBUhR3ZWJzaXRlY29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use putObjectAclOutputDescriptor instead')
const PutObjectAclOutput$json = {
  '1': 'PutObjectAclOutput',
  '2': [
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `PutObjectAclOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectAclOutputDescriptor = $convert.base64Decode(
    'ChJQdXRPYmplY3RBY2xPdXRwdXQSPgoOcmVxdWVzdGNoYXJnZWQYk9CruQEgASgOMhIuczMuUm'
    'VxdWVzdENoYXJnZWRSDnJlcXVlc3RjaGFyZ2Vk');

@$core.Deprecated('Use putObjectAclRequestDescriptor instead')
const PutObjectAclRequest$json = {
  '1': 'PutObjectAclRequest',
  '2': [
    {
      '1': 'acl',
      '3': 394696836,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectCannedACL',
      '10': 'acl'
    },
    {
      '1': 'accesscontrolpolicy',
      '3': 302514423,
      '4': 1,
      '5': 11,
      '6': '.s3.AccessControlPolicy',
      '8': {},
      '10': 'accesscontrolpolicy'
    },
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'grantfullcontrol',
      '3': 102486874,
      '4': 1,
      '5': 9,
      '10': 'grantfullcontrol'
    },
    {'1': 'grantread', '3': 4518126, '4': 1, '5': 9, '10': 'grantread'},
    {'1': 'grantreadacp', '3': 208410960, '4': 1, '5': 9, '10': 'grantreadacp'},
    {'1': 'grantwrite', '3': 286850821, '4': 1, '5': 9, '10': 'grantwrite'},
    {
      '1': 'grantwriteacp',
      '3': 157594957,
      '4': 1,
      '5': 9,
      '10': 'grantwriteacp'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `PutObjectAclRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectAclRequestDescriptor = $convert.base64Decode(
    'ChNQdXRPYmplY3RBY2xSZXF1ZXN0EikKA2FjbBiEsZq8ASABKA4yEy5zMy5PYmplY3RDYW5uZW'
    'RBQ0xSA2FjbBJTChNhY2Nlc3Njb250cm9scG9saWN5GPeBoJABIAEoCzIXLnMzLkFjY2Vzc0Nv'
    'bnRyb2xQb2xpY3lCBIi1GAFSE2FjY2Vzc2NvbnRyb2xwb2xpY3kSGQoGYnVja2V0GNjquBogAS'
    'gJUgZidWNrZXQSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFs'
    'Z29yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbn'
    'RlbnRtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRv'
    'd25lchItChBncmFudGZ1bGxjb250cm9sGNqm7zAgASgJUhBncmFudGZ1bGxjb250cm9sEh8KCW'
    'dyYW50cmVhZBju4ZMCIAEoCVIJZ3JhbnRyZWFkEiUKDGdyYW50cmVhZGFjcBjQsrBjIAEoCVIM'
    'Z3JhbnRyZWFkYWNwEiIKCmdyYW50d3JpdGUYhf7jiAEgASgJUgpncmFudHdyaXRlEicKDWdyYW'
    '50d3JpdGVhY3AYzeqSSyABKAlSDWdyYW50d3JpdGVhY3ASEwoDa2V5GI2S62ggASgJUgNrZXkS'
    'OAoMcmVxdWVzdHBheWVyGKTm4fUBIAEoDjIQLnMzLlJlcXVlc3RQYXllclIMcmVxdWVzdHBheW'
    'VyEiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb25pZA==');

@$core.Deprecated('Use putObjectLegalHoldOutputDescriptor instead')
const PutObjectLegalHoldOutput$json = {
  '1': 'PutObjectLegalHoldOutput',
  '2': [
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `PutObjectLegalHoldOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectLegalHoldOutputDescriptor =
    $convert.base64Decode(
        'ChhQdXRPYmplY3RMZWdhbEhvbGRPdXRwdXQSPgoOcmVxdWVzdGNoYXJnZWQYk9CruQEgASgOMh'
        'IuczMuUmVxdWVzdENoYXJnZWRSDnJlcXVlc3RjaGFyZ2Vk');

@$core.Deprecated('Use putObjectLegalHoldRequestDescriptor instead')
const PutObjectLegalHoldRequest$json = {
  '1': 'PutObjectLegalHoldRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'legalhold',
      '3': 141145940,
      '4': 1,
      '5': 11,
      '6': '.s3.ObjectLockLegalHold',
      '8': {},
      '10': 'legalhold'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `PutObjectLegalHoldRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectLegalHoldRequestDescriptor = $convert.base64Decode(
    'ChlQdXRPYmplY3RMZWdhbEhvbGRSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0Ek'
    'YKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMuQ2hlY2tzdW1BbGdvcml0aG1SEWNo'
    'ZWNrc3VtYWxnb3JpdGhtEiIKCmNvbnRlbnRtZDUYi+rq8QEgASgJUgpjb250ZW50bWQ1EjMKE2'
    'V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXISEwoDa2V5'
    'GI2S62ggASgJUgNrZXkSPgoJbGVnYWxob2xkGNTupkMgASgLMhcuczMuT2JqZWN0TG9ja0xlZ2'
    'FsSG9sZEIEiLUYAVIJbGVnYWxob2xkEjgKDHJlcXVlc3RwYXllchik5uH1ASABKA4yEC5zMy5S'
    'ZXF1ZXN0UGF5ZXJSDHJlcXVlc3RwYXllchIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW'
    '9uaWQ=');

@$core.Deprecated('Use putObjectLockConfigurationOutputDescriptor instead')
const PutObjectLockConfigurationOutput$json = {
  '1': 'PutObjectLockConfigurationOutput',
  '2': [
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `PutObjectLockConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectLockConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CiBQdXRPYmplY3RMb2NrQ29uZmlndXJhdGlvbk91dHB1dBI+Cg5yZXF1ZXN0Y2hhcmdlZBiT0K'
        'u5ASABKA4yEi5zMy5SZXF1ZXN0Q2hhcmdlZFIOcmVxdWVzdGNoYXJnZWQ=');

@$core.Deprecated('Use putObjectLockConfigurationRequestDescriptor instead')
const PutObjectLockConfigurationRequest$json = {
  '1': 'PutObjectLockConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'objectlockconfiguration',
      '3': 108348298,
      '4': 1,
      '5': 11,
      '6': '.s3.ObjectLockConfiguration',
      '8': {},
      '10': 'objectlockconfiguration'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'token', '3': 439704531, '4': 1, '5': 9, '10': 'token'},
  ],
};

/// Descriptor for `PutObjectLockConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectLockConfigurationRequestDescriptor = $convert.base64Decode(
    'CiFQdXRPYmplY3RMb2NrQ29uZmlndXJhdGlvblJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUg'
    'ZidWNrZXQSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFsZ29y'
    'aXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbnRlbn'
    'RtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25l'
    'chJeChdvYmplY3Rsb2NrY29uZmlndXJhdGlvbhiKh9UzIAEoCzIbLnMzLk9iamVjdExvY2tDb2'
    '5maWd1cmF0aW9uQgSItRgBUhdvYmplY3Rsb2NrY29uZmlndXJhdGlvbhI4CgxyZXF1ZXN0cGF5'
    'ZXIYpObh9QEgASgOMhAuczMuUmVxdWVzdFBheWVyUgxyZXF1ZXN0cGF5ZXISGAoFdG9rZW4Y07'
    'fV0QEgASgJUgV0b2tlbg==');

@$core.Deprecated('Use putObjectOutputDescriptor instead')
const PutObjectOutput$json = {
  '1': 'PutObjectOutput',
  '2': [
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'checksumtype',
      '3': 97935171,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumType',
      '10': 'checksumtype'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {
      '1': 'ssekmsencryptioncontext',
      '3': 149030970,
      '4': 1,
      '5': 9,
      '10': 'ssekmsencryptioncontext'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {'1': 'size', '3': 105352829, '4': 1, '5': 3, '10': 'size'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `PutObjectOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectOutputDescriptor = $convert.base64Decode(
    'Cg9QdXRPYmplY3RPdXRwdXQSLgoQYnVja2V0a2V5ZW5hYmxlZBjsoozPASABKAhSEGJ1Y2tldG'
    'tleWVuYWJsZWQSJwoNY2hlY2tzdW1jcmMzMhjCo80zIAEoCVINY2hlY2tzdW1jcmMzMhIpCg5j'
    'aGVja3N1bWNyYzMyYxinn6VMIAEoCVIOY2hlY2tzdW1jcmMzMmMSLwoRY2hlY2tzdW1jcmM2NG'
    '52bWUYzbyLOCABKAlSEWNoZWNrc3VtY3JjNjRudm1lEiYKDGNoZWNrc3Vtc2hhMRjE7OCKASAB'
    'KAlSDGNoZWNrc3Vtc2hhMRIpCg5jaGVja3N1bXNoYTI1Nhi++dxEIAEoCVIOY2hlY2tzdW1zaG'
    'EyNTYSNwoMY2hlY2tzdW10eXBlGMO+2S4gASgOMhAuczMuQ2hlY2tzdW1UeXBlUgxjaGVja3N1'
    'bXR5cGUSFgoEZXRhZxiB37OVASABKAlSBGV0YWcSIQoKZXhwaXJhdGlvbhiJqZ91IAEoCVIKZX'
    'hwaXJhdGlvbhI+Cg5yZXF1ZXN0Y2hhcmdlZBiT0Ku5ASABKA4yEi5zMy5SZXF1ZXN0Q2hhcmdl'
    'ZFIOcmVxdWVzdGNoYXJnZWQSNQoUc3NlY3VzdG9tZXJhbGdvcml0aG0Y0MmBKyABKAlSFHNzZW'
    'N1c3RvbWVyYWxnb3JpdGhtEi4KEXNzZWN1c3RvbWVya2V5bWQ1GOjRFyABKAlSEXNzZWN1c3Rv'
    'bWVya2V5bWQ1EjsKF3NzZWttc2VuY3J5cHRpb25jb250ZXh0GLqQiEcgASgJUhdzc2VrbXNlbm'
    'NyeXB0aW9uY29udGV4dBIkCgtzc2VrbXNrZXlpZBjYiNTUASABKAlSC3NzZWttc2tleWlkEk8K'
    'FHNlcnZlcnNpZGVlbmNyeXB0aW9uGLGOnwQgASgOMhguczMuU2VydmVyU2lkZUVuY3J5cHRpb2'
    '5SFHNlcnZlcnNpZGVlbmNyeXB0aW9uEhUKBHNpemUY/ZyeMiABKANSBHNpemUSIAoJdmVyc2lv'
    'bmlkGJvhmaEBIAEoCVIJdmVyc2lvbmlk');

@$core.Deprecated('Use putObjectRequestDescriptor instead')
const PutObjectRequest$json = {
  '1': 'PutObjectRequest',
  '2': [
    {
      '1': 'acl',
      '3': 394696836,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectCannedACL',
      '10': 'acl'
    },
    {'1': 'body', '3': 42602646, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {'1': 'cachecontrol', '3': 288966655, '4': 1, '5': 9, '10': 'cachecontrol'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'contentdisposition',
      '3': 120040130,
      '4': 1,
      '5': 9,
      '10': 'contentdisposition'
    },
    {
      '1': 'contentencoding',
      '3': 317106228,
      '4': 1,
      '5': 9,
      '10': 'contentencoding'
    },
    {
      '1': 'contentlanguage',
      '3': 108485649,
      '4': 1,
      '5': 9,
      '10': 'contentlanguage'
    },
    {
      '1': 'contentlength',
      '3': 227596631,
      '4': 1,
      '5': 3,
      '10': 'contentlength'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'expires', '3': 128582948, '4': 1, '5': 9, '10': 'expires'},
    {
      '1': 'grantfullcontrol',
      '3': 102486874,
      '4': 1,
      '5': 9,
      '10': 'grantfullcontrol'
    },
    {'1': 'grantread', '3': 4518126, '4': 1, '5': 9, '10': 'grantread'},
    {'1': 'grantreadacp', '3': 208410960, '4': 1, '5': 9, '10': 'grantreadacp'},
    {
      '1': 'grantwriteacp',
      '3': 157594957,
      '4': 1,
      '5': 9,
      '10': 'grantwriteacp'
    },
    {'1': 'ifmatch', '3': 93178704, '4': 1, '5': 9, '10': 'ifmatch'},
    {'1': 'ifnonematch', '3': 231211830, '4': 1, '5': 9, '10': 'ifnonematch'},
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 3,
      '5': 11,
      '6': '.s3.PutObjectRequest.MetadataEntry',
      '10': 'metadata'
    },
    {
      '1': 'objectlocklegalholdstatus',
      '3': 536561974,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockLegalHoldStatus',
      '10': 'objectlocklegalholdstatus'
    },
    {
      '1': 'objectlockmode',
      '3': 189255203,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockMode',
      '10': 'objectlockmode'
    },
    {
      '1': 'objectlockretainuntildate',
      '3': 264584249,
      '4': 1,
      '5': 9,
      '10': 'objectlockretainuntildate'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {
      '1': 'ssekmsencryptioncontext',
      '3': 149030970,
      '4': 1,
      '5': 9,
      '10': 'ssekmsencryptioncontext'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'tagging', '3': 33436541, '4': 1, '5': 9, '10': 'tagging'},
    {
      '1': 'websiteredirectlocation',
      '3': 71844662,
      '4': 1,
      '5': 9,
      '10': 'websiteredirectlocation'
    },
    {
      '1': 'writeoffsetbytes',
      '3': 496666419,
      '4': 1,
      '5': 3,
      '10': 'writeoffsetbytes'
    },
  ],
  '3': [PutObjectRequest_MetadataEntry$json],
};

@$core.Deprecated('Use putObjectRequestDescriptor instead')
const PutObjectRequest_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PutObjectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectRequestDescriptor = $convert.base64Decode(
    'ChBQdXRPYmplY3RSZXF1ZXN0EikKA2FjbBiEsZq8ASABKA4yEy5zMy5PYmplY3RDYW5uZWRBQ0'
    'xSA2FjbBIbCgRib2R5GJahqBQgASgMQgSItRgBUgRib2R5EhkKBmJ1Y2tldBjY6rgaIAEoCVIG'
    'YnVja2V0Ei4KEGJ1Y2tldGtleWVuYWJsZWQY7KKMzwEgASgIUhBidWNrZXRrZXllbmFibGVkEi'
    'YKDGNhY2hlY29udHJvbBj/j+WJASABKAlSDGNhY2hlY29udHJvbBJGChFjaGVja3N1bWFsZ29y'
    'aXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3VtQWxnb3JpdGhtUhFjaGVja3N1bWFsZ29yaXRobR'
    'InCg1jaGVja3N1bWNyYzMyGMKjzTMgASgJUg1jaGVja3N1bWNyYzMyEikKDmNoZWNrc3VtY3Jj'
    'MzJjGKefpUwgASgJUg5jaGVja3N1bWNyYzMyYxIvChFjaGVja3N1bWNyYzY0bnZtZRjNvIs4IA'
    'EoCVIRY2hlY2tzdW1jcmM2NG52bWUSJgoMY2hlY2tzdW1zaGExGMTs4IoBIAEoCVIMY2hlY2tz'
    'dW1zaGExEikKDmNoZWNrc3Vtc2hhMjU2GL753EQgASgJUg5jaGVja3N1bXNoYTI1NhIxChJjb2'
    '50ZW50ZGlzcG9zaXRpb24YwtWeOSABKAlSEmNvbnRlbnRkaXNwb3NpdGlvbhIsCg9jb250ZW50'
    'ZW5jb2RpbmcYtNCalwEgASgJUg9jb250ZW50ZW5jb2RpbmcSKwoPY29udGVudGxhbmd1YWdlGJ'
    'G43TMgASgJUg9jb250ZW50bGFuZ3VhZ2USJwoNY29udGVudGxlbmd0aBjXssNsIAEoA1INY29u'
    'dGVudGxlbmd0aBIiCgpjb250ZW50bWQ1GIvq6vEBIAEoCVIKY29udGVudG1kNRIkCgtjb250ZW'
    '50dHlwZRiT1eieASABKAlSC2NvbnRlbnR0eXBlEjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938'
    'PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXISGwoHZXhwaXJlcxikiqg9IAEoCVIHZXhwaXJlcx'
    'ItChBncmFudGZ1bGxjb250cm9sGNqm7zAgASgJUhBncmFudGZ1bGxjb250cm9sEh8KCWdyYW50'
    'cmVhZBju4ZMCIAEoCVIJZ3JhbnRyZWFkEiUKDGdyYW50cmVhZGFjcBjQsrBjIAEoCVIMZ3Jhbn'
    'RyZWFkYWNwEicKDWdyYW50d3JpdGVhY3AYzeqSSyABKAlSDWdyYW50d3JpdGVhY3ASGwoHaWZt'
    'YXRjaBjQlrcsIAEoCVIHaWZtYXRjaBIjCgtpZm5vbmVtYXRjaBi2hqBuIAEoCVILaWZub25lbW'
    'F0Y2gSEwoDa2V5GI2S62ggASgJUgNrZXkSQgoIbWV0YWRhdGEY4eKP4AEgAygLMiIuczMuUHV0'
    'T2JqZWN0UmVxdWVzdC5NZXRhZGF0YUVudHJ5UghtZXRhZGF0YRJfChlvYmplY3Rsb2NrbGVnYW'
    'xob2xkc3RhdHVzGLaS7f8BIAEoDjIdLnMzLk9iamVjdExvY2tMZWdhbEhvbGRTdGF0dXNSGW9i'
    'amVjdGxvY2tsZWdhbGhvbGRzdGF0dXMSPQoOb2JqZWN0bG9ja21vZGUYo5yfWiABKA4yEi5zMy'
    '5PYmplY3RMb2NrTW9kZVIOb2JqZWN0bG9ja21vZGUSPwoZb2JqZWN0bG9ja3JldGFpbnVudGls'
    'ZGF0ZRi5+JR+IAEoCVIZb2JqZWN0bG9ja3JldGFpbnVudGlsZGF0ZRI4CgxyZXF1ZXN0cGF5ZX'
    'IYpObh9QEgASgOMhAuczMuUmVxdWVzdFBheWVyUgxyZXF1ZXN0cGF5ZXISNQoUc3NlY3VzdG9t'
    'ZXJhbGdvcml0aG0Y0MmBKyABKAlSFHNzZWN1c3RvbWVyYWxnb3JpdGhtEikKDnNzZWN1c3RvbW'
    'Vya2V5GJr+9DsgASgJUg5zc2VjdXN0b21lcmtleRIuChFzc2VjdXN0b21lcmtleW1kNRjo0Rcg'
    'ASgJUhFzc2VjdXN0b21lcmtleW1kNRI7Chdzc2VrbXNlbmNyeXB0aW9uY29udGV4dBi6kIhHIA'
    'EoCVIXc3Nla21zZW5jcnlwdGlvbmNvbnRleHQSJAoLc3Nla21za2V5aWQY2IjU1AEgASgJUgtz'
    'c2VrbXNrZXlpZBJPChRzZXJ2ZXJzaWRlZW5jcnlwdGlvbhixjp8EIAEoDjIYLnMzLlNlcnZlcl'
    'NpZGVFbmNyeXB0aW9uUhRzZXJ2ZXJzaWRlZW5jcnlwdGlvbhI4CgxzdG9yYWdlY2xhc3MYx4jE'
    'uwEgASgOMhAuczMuU3RvcmFnZUNsYXNzUgxzdG9yYWdlY2xhc3MSGwoHdGFnZ2luZxj95vgPIA'
    'EoCVIHdGFnZ2luZxI7Chd3ZWJzaXRlcmVkaXJlY3Rsb2NhdGlvbhi2hqEiIAEoCVIXd2Vic2l0'
    'ZXJlZGlyZWN0bG9jYXRpb24SLgoQd3JpdGVvZmZzZXRieXRlcxizjursASABKANSEHdyaXRlb2'
    'Zmc2V0Ynl0ZXMaOwoNTWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgC'
    'IAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use putObjectRetentionOutputDescriptor instead')
const PutObjectRetentionOutput$json = {
  '1': 'PutObjectRetentionOutput',
  '2': [
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `PutObjectRetentionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectRetentionOutputDescriptor =
    $convert.base64Decode(
        'ChhQdXRPYmplY3RSZXRlbnRpb25PdXRwdXQSPgoOcmVxdWVzdGNoYXJnZWQYk9CruQEgASgOMh'
        'IuczMuUmVxdWVzdENoYXJnZWRSDnJlcXVlc3RjaGFyZ2Vk');

@$core.Deprecated('Use putObjectRetentionRequestDescriptor instead')
const PutObjectRetentionRequest$json = {
  '1': 'PutObjectRetentionRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'bypassgovernanceretention',
      '3': 129938942,
      '4': 1,
      '5': 8,
      '10': 'bypassgovernanceretention'
    },
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'retention',
      '3': 299600946,
      '4': 1,
      '5': 11,
      '6': '.s3.ObjectLockRetention',
      '8': {},
      '10': 'retention'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `PutObjectRetentionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectRetentionRequestDescriptor = $convert.base64Decode(
    'ChlQdXRPYmplY3RSZXRlbnRpb25SZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0Ej'
    '8KGWJ5cGFzc2dvdmVybmFuY2VyZXRlbnRpb24Y/uv6PSABKAhSGWJ5cGFzc2dvdmVybmFuY2Vy'
    'ZXRlbnRpb24SRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFsZ2'
    '9yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbnRl'
    'bnRtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd2'
    '5lchITCgNrZXkYjZLraCABKAlSA2tleRI4CgxyZXF1ZXN0cGF5ZXIYpObh9QEgASgOMhAuczMu'
    'UmVxdWVzdFBheWVyUgxyZXF1ZXN0cGF5ZXISPwoJcmV0ZW50aW9uGLKY7o4BIAEoCzIXLnMzLk'
    '9iamVjdExvY2tSZXRlbnRpb25CBIi1GAFSCXJldGVudGlvbhIgCgl2ZXJzaW9uaWQYm+GZoQEg'
    'ASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use putObjectTaggingOutputDescriptor instead')
const PutObjectTaggingOutput$json = {
  '1': 'PutObjectTaggingOutput',
  '2': [
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `PutObjectTaggingOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectTaggingOutputDescriptor =
    $convert.base64Decode(
        'ChZQdXRPYmplY3RUYWdnaW5nT3V0cHV0EiAKCXZlcnNpb25pZBib4ZmhASABKAlSCXZlcnNpb2'
        '5pZA==');

@$core.Deprecated('Use putObjectTaggingRequestDescriptor instead')
const PutObjectTaggingRequest$json = {
  '1': 'PutObjectTaggingRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'tagging',
      '3': 33436541,
      '4': 1,
      '5': 11,
      '6': '.s3.Tagging',
      '8': {},
      '10': 'tagging'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `PutObjectTaggingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putObjectTaggingRequestDescriptor = $convert.base64Decode(
    'ChdQdXRPYmplY3RUYWdnaW5nUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBJGCh'
    'FjaGVja3N1bWFsZ29yaXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3VtQWxnb3JpdGhtUhFjaGVj'
    'a3N1bWFsZ29yaXRobRIiCgpjb250ZW50bWQ1GIvq6vEBIAEoCVIKY29udGVudG1kNRIzChNleH'
    'BlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhMKA2tleRiN'
    'kutoIAEoCVIDa2V5EjgKDHJlcXVlc3RwYXllchik5uH1ASABKA4yEC5zMy5SZXF1ZXN0UGF5ZX'
    'JSDHJlcXVlc3RwYXllchIuCgd0YWdnaW5nGP3m+A8gASgLMgsuczMuVGFnZ2luZ0IEiLUYAVIH'
    'dGFnZ2luZxIgCgl2ZXJzaW9uaWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use putPublicAccessBlockRequestDescriptor instead')
const PutPublicAccessBlockRequest$json = {
  '1': 'PutPublicAccessBlockRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'publicaccessblockconfiguration',
      '3': 136498568,
      '4': 1,
      '5': 11,
      '6': '.s3.PublicAccessBlockConfiguration',
      '8': {},
      '10': 'publicaccessblockconfiguration'
    },
  ],
};

/// Descriptor for `PutPublicAccessBlockRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putPublicAccessBlockRequestDescriptor = $convert.base64Decode(
    'ChtQdXRQdWJsaWNBY2Nlc3NCbG9ja1JlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZX'
    'QSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiABKA4yFS5zMy5DaGVja3N1bUFsZ29yaXRobVIR'
    'Y2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbnRlbnRtZDUSMw'
    'oTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchJzCh5w'
    'dWJsaWNhY2Nlc3NibG9ja2NvbmZpZ3VyYXRpb24YiJuLQSABKAsyIi5zMy5QdWJsaWNBY2Nlc3'
    'NCbG9ja0NvbmZpZ3VyYXRpb25CBIi1GAFSHnB1YmxpY2FjY2Vzc2Jsb2NrY29uZmlndXJhdGlv'
    'bg==');

@$core.Deprecated('Use queueConfigurationDescriptor instead')
const QueueConfiguration$json = {
  '1': 'QueueConfiguration',
  '2': [
    {
      '1': 'events',
      '3': 3416229,
      '4': 3,
      '5': 14,
      '6': '.s3.Event',
      '10': 'events'
    },
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.NotificationConfigurationFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'queuearn', '3': 400520576, '4': 1, '5': 9, '10': 'queuearn'},
  ],
};

/// Descriptor for `QueueConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queueConfigurationDescriptor = $convert.base64Decode(
    'ChJRdWV1ZUNvbmZpZ3VyYXRpb24SJAoGZXZlbnRzGKXB0AEgAygOMgkuczMuRXZlbnRSBmV2ZW'
    '50cxI/CgZmaWx0ZXIYmIGnpQEgASgLMiMuczMuTm90aWZpY2F0aW9uQ29uZmlndXJhdGlvbkZp'
    'bHRlclIGZmlsdGVyEhIKAmlkGIHyorcBIAEoCVICaWQSHgoIcXVldWVhcm4YgOv9vgEgASgJUg'
    'hxdWV1ZWFybg==');

@$core.Deprecated('Use recordExpirationDescriptor instead')
const RecordExpiration$json = {
  '1': 'RecordExpiration',
  '2': [
    {'1': 'days', '3': 494075051, '4': 1, '5': 5, '10': 'days'},
    {
      '1': 'expiration',
      '3': 245879945,
      '4': 1,
      '5': 14,
      '6': '.s3.ExpirationState',
      '10': 'expiration'
    },
  ],
};

/// Descriptor for `RecordExpiration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recordExpirationDescriptor = $convert.base64Decode(
    'ChBSZWNvcmRFeHBpcmF0aW9uEhYKBGRheXMYq/nL6wEgASgFUgRkYXlzEjYKCmV4cGlyYXRpb2'
    '4YiamfdSABKA4yEy5zMy5FeHBpcmF0aW9uU3RhdGVSCmV4cGlyYXRpb24=');

@$core.Deprecated('Use recordsEventDescriptor instead')
const RecordsEvent$json = {
  '1': 'RecordsEvent',
  '2': [
    {'1': 'payload', '3': 6526790, '4': 1, '5': 12, '10': 'payload'},
  ],
};

/// Descriptor for `RecordsEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recordsEventDescriptor = $convert.base64Decode(
    'CgxSZWNvcmRzRXZlbnQSGwoHcGF5bG9hZBjGro4DIAEoDFIHcGF5bG9hZA==');

@$core.Deprecated('Use redirectDescriptor instead')
const Redirect$json = {
  '1': 'Redirect',
  '2': [
    {'1': 'hostname', '3': 119731117, '4': 1, '5': 9, '10': 'hostname'},
    {
      '1': 'httpredirectcode',
      '3': 348964153,
      '4': 1,
      '5': 9,
      '10': 'httpredirectcode'
    },
    {
      '1': 'protocol',
      '3': 173534166,
      '4': 1,
      '5': 14,
      '6': '.s3.Protocol',
      '10': 'protocol'
    },
    {
      '1': 'replacekeyprefixwith',
      '3': 532110043,
      '4': 1,
      '5': 9,
      '10': 'replacekeyprefixwith'
    },
    {
      '1': 'replacekeywith',
      '3': 406460263,
      '4': 1,
      '5': 9,
      '10': 'replacekeywith'
    },
  ],
};

/// Descriptor for `Redirect`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redirectDescriptor = $convert.base64Decode(
    'CghSZWRpcmVjdBIdCghob3N0bmFtZRit54s5IAEoCVIIaG9zdG5hbWUSLgoQaHR0cHJlZGlyZW'
    'N0Y29kZRi5irOmASABKAlSEGh0dHByZWRpcmVjdGNvZGUSKwoIcHJvdG9jb2wY1tffUiABKA4y'
    'DC5zMy5Qcm90b2NvbFIIcHJvdG9jb2wSNgoUcmVwbGFjZWtleXByZWZpeHdpdGgY27Xd/QEgAS'
    'gJUhRyZXBsYWNla2V5cHJlZml4d2l0aBIqCg5yZXBsYWNla2V5d2l0aBjnrujBASABKAlSDnJl'
    'cGxhY2VrZXl3aXRo');

@$core.Deprecated('Use redirectAllRequestsToDescriptor instead')
const RedirectAllRequestsTo$json = {
  '1': 'RedirectAllRequestsTo',
  '2': [
    {'1': 'hostname', '3': 119731117, '4': 1, '5': 9, '10': 'hostname'},
    {
      '1': 'protocol',
      '3': 173534166,
      '4': 1,
      '5': 14,
      '6': '.s3.Protocol',
      '10': 'protocol'
    },
  ],
};

/// Descriptor for `RedirectAllRequestsTo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redirectAllRequestsToDescriptor = $convert.base64Decode(
    'ChVSZWRpcmVjdEFsbFJlcXVlc3RzVG8SHQoIaG9zdG5hbWUYreeLOSABKAlSCGhvc3RuYW1lEi'
    'sKCHByb3RvY29sGNbX31IgASgOMgwuczMuUHJvdG9jb2xSCHByb3RvY29s');

@$core.Deprecated('Use renameObjectOutputDescriptor instead')
const RenameObjectOutput$json = {
  '1': 'RenameObjectOutput',
};

/// Descriptor for `RenameObjectOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List renameObjectOutputDescriptor =
    $convert.base64Decode('ChJSZW5hbWVPYmplY3RPdXRwdXQ=');

@$core.Deprecated('Use renameObjectRequestDescriptor instead')
const RenameObjectRequest$json = {
  '1': 'RenameObjectRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'destinationifmatch',
      '3': 494878870,
      '4': 1,
      '5': 9,
      '10': 'destinationifmatch'
    },
    {
      '1': 'destinationifmodifiedsince',
      '3': 253871164,
      '4': 1,
      '5': 9,
      '10': 'destinationifmodifiedsince'
    },
    {
      '1': 'destinationifnonematch',
      '3': 176301672,
      '4': 1,
      '5': 9,
      '10': 'destinationifnonematch'
    },
    {
      '1': 'destinationifunmodifiedsince',
      '3': 58294063,
      '4': 1,
      '5': 9,
      '10': 'destinationifunmodifiedsince'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'renamesource', '3': 72928743, '4': 1, '5': 9, '10': 'renamesource'},
    {
      '1': 'sourceifmatch',
      '3': 528463881,
      '4': 1,
      '5': 9,
      '10': 'sourceifmatch'
    },
    {
      '1': 'sourceifmodifiedsince',
      '3': 533054399,
      '4': 1,
      '5': 9,
      '10': 'sourceifmodifiedsince'
    },
    {
      '1': 'sourceifnonematch',
      '3': 56549307,
      '4': 1,
      '5': 9,
      '10': 'sourceifnonematch'
    },
    {
      '1': 'sourceifunmodifiedsince',
      '3': 449515432,
      '4': 1,
      '5': 9,
      '10': 'sourceifunmodifiedsince'
    },
  ],
};

/// Descriptor for `RenameObjectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List renameObjectRequestDescriptor = $convert.base64Decode(
    'ChNSZW5hbWVPYmplY3RSZXF1ZXN0EhkKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EiMKC2NsaW'
    'VudHRva2VuGMz7u0EgASgJUgtjbGllbnR0b2tlbhIyChJkZXN0aW5hdGlvbmlmbWF0Y2gYloH9'
    '6wEgASgJUhJkZXN0aW5hdGlvbmlmbWF0Y2gSQQoaZGVzdGluYXRpb25pZm1vZGlmaWVkc2luY2'
    'UYvIiHeSABKAlSGmRlc3RpbmF0aW9uaWZtb2RpZmllZHNpbmNlEjkKFmRlc3RpbmF0aW9uaWZu'
    'b25lbWF0Y2gY6MyIVCABKAlSFmRlc3RpbmF0aW9uaWZub25lbWF0Y2gSRQocZGVzdGluYXRpb2'
    '5pZnVubW9kaWZpZWRzaW5jZRiv/uUbIAEoCVIcZGVzdGluYXRpb25pZnVubW9kaWZpZWRzaW5j'
    'ZRITCgNrZXkYjZLraCABKAlSA2tleRIlCgxyZW5hbWVzb3VyY2UY55vjIiABKAlSDHJlbmFtZX'
    'NvdXJjZRIoCg1zb3VyY2VpZm1hdGNoGInw/vsBIAEoCVINc291cmNlaWZtYXRjaBI4ChVzb3Vy'
    'Y2VpZm1vZGlmaWVkc2luY2UYv4eX/gEgASgJUhVzb3VyY2VpZm1vZGlmaWVkc2luY2USLwoRc2'
    '91cmNlaWZub25lbWF0Y2gYu7/7GiABKAlSEXNvdXJjZWlmbm9uZW1hdGNoEjwKF3NvdXJjZWlm'
    'dW5tb2RpZmllZHNpbmNlGKifrNYBIAEoCVIXc291cmNlaWZ1bm1vZGlmaWVkc2luY2U=');

@$core.Deprecated('Use replicaModificationsDescriptor instead')
const ReplicaModifications$json = {
  '1': 'ReplicaModifications',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.ReplicaModificationsStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ReplicaModifications`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicaModificationsDescriptor = $convert.base64Decode(
    'ChRSZXBsaWNhTW9kaWZpY2F0aW9ucxI5CgZzdGF0dXMYkOT7AiABKA4yHi5zMy5SZXBsaWNhTW'
    '9kaWZpY2F0aW9uc1N0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use replicationConfigurationDescriptor instead')
const ReplicationConfiguration$json = {
  '1': 'ReplicationConfiguration',
  '2': [
    {'1': 'role', '3': 271285818, '4': 1, '5': 9, '10': 'role'},
    {
      '1': 'rules',
      '3': 42675585,
      '4': 3,
      '5': 11,
      '6': '.s3.ReplicationRule',
      '10': 'rules'
    },
  ],
};

/// Descriptor for `ReplicationConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicationConfigurationDescriptor =
    $convert.base64Decode(
        'ChhSZXBsaWNhdGlvbkNvbmZpZ3VyYXRpb24SFgoEcm9sZRi6/K2BASABKAlSBHJvbGUSLAoFcn'
        'VsZXMYgdusFCADKAsyEy5zMy5SZXBsaWNhdGlvblJ1bGVSBXJ1bGVz');

@$core.Deprecated('Use replicationRuleDescriptor instead')
const ReplicationRule$json = {
  '1': 'ReplicationRule',
  '2': [
    {
      '1': 'deletemarkerreplication',
      '3': 18229087,
      '4': 1,
      '5': 11,
      '6': '.s3.DeleteMarkerReplication',
      '10': 'deletemarkerreplication'
    },
    {
      '1': 'destination',
      '3': 457443680,
      '4': 1,
      '5': 11,
      '6': '.s3.Destination',
      '10': 'destination'
    },
    {
      '1': 'existingobjectreplication',
      '3': 500093148,
      '4': 1,
      '5': 11,
      '6': '.s3.ExistingObjectReplication',
      '10': 'existingobjectreplication'
    },
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.ReplicationRuleFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384363361, '4': 1, '5': 9, '10': 'id'},
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {
      '1': 'sourceselectioncriteria',
      '3': 201951000,
      '4': 1,
      '5': 11,
      '6': '.s3.SourceSelectionCriteria',
      '10': 'sourceselectioncriteria'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.ReplicationRuleStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ReplicationRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicationRuleDescriptor = $convert.base64Decode(
    'Cg9SZXBsaWNhdGlvblJ1bGUSWAoXZGVsZXRlbWFya2VycmVwbGljYXRpb24Y387YCCABKAsyGy'
    '5zMy5EZWxldGVNYXJrZXJSZXBsaWNhdGlvblIXZGVsZXRlbWFya2VycmVwbGljYXRpb24SNQoL'
    'ZGVzdGluYXRpb24Y4JKQ2gEgASgLMg8uczMuRGVzdGluYXRpb25SC2Rlc3RpbmF0aW9uEl8KGW'
    'V4aXN0aW5nb2JqZWN0cmVwbGljYXRpb24Y3KG77gEgASgLMh0uczMuRXhpc3RpbmdPYmplY3RS'
    'ZXBsaWNhdGlvblIZZXhpc3RpbmdvYmplY3RyZXBsaWNhdGlvbhI1CgZmaWx0ZXIYmIGnpQEgAS'
    'gLMhkuczMuUmVwbGljYXRpb25SdWxlRmlsdGVyUgZmaWx0ZXISEgoCaWQY4dajtwEgASgJUgJp'
    'ZBIaCgZwcmVmaXgY6rPTggEgASgJUgZwcmVmaXgSHQoIcHJpb3JpdHkYqr62NCABKAVSCHByaW'
    '9yaXR5ElgKF3NvdXJjZXNlbGVjdGlvbmNyaXRlcmlhGJiOpmAgASgLMhsuczMuU291cmNlU2Vs'
    'ZWN0aW9uQ3JpdGVyaWFSF3NvdXJjZXNlbGVjdGlvbmNyaXRlcmlhEjQKBnN0YXR1cxiQ5PsCIA'
    'EoDjIZLnMzLlJlcGxpY2F0aW9uUnVsZVN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use replicationRuleAndOperatorDescriptor instead')
const ReplicationRuleAndOperator$json = {
  '1': 'ReplicationRuleAndOperator',
  '2': [
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ReplicationRuleAndOperator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicationRuleAndOperatorDescriptor =
    $convert.base64Decode(
        'ChpSZXBsaWNhdGlvblJ1bGVBbmRPcGVyYXRvchIaCgZwcmVmaXgY6rPTggEgASgJUgZwcmVmaX'
        'gSHwoEdGFncxjBwfa1ASADKAsyBy5zMy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use replicationRuleFilterDescriptor instead')
const ReplicationRuleFilter$json = {
  '1': 'ReplicationRuleFilter',
  '2': [
    {
      '1': 'and',
      '3': 297135431,
      '4': 1,
      '5': 11,
      '6': '.s3.ReplicationRuleAndOperator',
      '10': 'and'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {'1': 'tag', '3': 411259956, '4': 1, '5': 11, '6': '.s3.Tag', '10': 'tag'},
  ],
};

/// Descriptor for `ReplicationRuleFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicationRuleFilterDescriptor = $convert.base64Decode(
    'ChVSZXBsaWNhdGlvblJ1bGVGaWx0ZXISNAoDYW5kGMfa140BIAEoCzIeLnMzLlJlcGxpY2F0aW'
    '9uUnVsZUFuZE9wZXJhdG9yUgNhbmQSGgoGcHJlZml4GOqz04IBIAEoCVIGcHJlZml4Eh0KA3Rh'
    'Zxi0qI3EASABKAsyBy5zMy5UYWdSA3RhZw==');

@$core.Deprecated('Use replicationTimeDescriptor instead')
const ReplicationTime$json = {
  '1': 'ReplicationTime',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.ReplicationTimeStatus',
      '10': 'status'
    },
    {
      '1': 'time',
      '3': 535094277,
      '4': 1,
      '5': 11,
      '6': '.s3.ReplicationTimeValue',
      '10': 'time'
    },
  ],
};

/// Descriptor for `ReplicationTime`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicationTimeDescriptor = $convert.base64Decode(
    'Cg9SZXBsaWNhdGlvblRpbWUSNAoGc3RhdHVzGJDk+wIgASgOMhkuczMuUmVwbGljYXRpb25UaW'
    '1lU3RhdHVzUgZzdGF0dXMSMAoEdGltZRiFyJP/ASABKAsyGC5zMy5SZXBsaWNhdGlvblRpbWVW'
    'YWx1ZVIEdGltZQ==');

@$core.Deprecated('Use replicationTimeValueDescriptor instead')
const ReplicationTimeValue$json = {
  '1': 'ReplicationTimeValue',
  '2': [
    {'1': 'minutes', '3': 383358607, '4': 1, '5': 5, '10': 'minutes'},
  ],
};

/// Descriptor for `ReplicationTimeValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicationTimeValueDescriptor = $convert.base64Decode(
    'ChRSZXBsaWNhdGlvblRpbWVWYWx1ZRIcCgdtaW51dGVzGI+t5rYBIAEoBVIHbWludXRlcw==');

@$core.Deprecated('Use requestPaymentConfigurationDescriptor instead')
const RequestPaymentConfiguration$json = {
  '1': 'RequestPaymentConfiguration',
  '2': [
    {
      '1': 'payer',
      '3': 174621923,
      '4': 1,
      '5': 14,
      '6': '.s3.Payer',
      '10': 'payer'
    },
  ],
};

/// Descriptor for `RequestPaymentConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestPaymentConfigurationDescriptor =
    $convert.base64Decode(
        'ChtSZXF1ZXN0UGF5bWVudENvbmZpZ3VyYXRpb24SIgoFcGF5ZXIY44miUyABKA4yCS5zMy5QYX'
        'llclIFcGF5ZXI=');

@$core.Deprecated('Use requestProgressDescriptor instead')
const RequestProgress$json = {
  '1': 'RequestProgress',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
  ],
};

/// Descriptor for `RequestProgress`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestProgressDescriptor = $convert.base64Decode(
    'Cg9SZXF1ZXN0UHJvZ3Jlc3MSHAoHZW5hYmxlZBi/yJvkASABKAhSB2VuYWJsZWQ=');

@$core.Deprecated('Use restoreObjectOutputDescriptor instead')
const RestoreObjectOutput$json = {
  '1': 'RestoreObjectOutput',
  '2': [
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'restoreoutputpath',
      '3': 4884176,
      '4': 1,
      '5': 9,
      '10': 'restoreoutputpath'
    },
  ],
};

/// Descriptor for `RestoreObjectOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreObjectOutputDescriptor = $convert.base64Decode(
    'ChNSZXN0b3JlT2JqZWN0T3V0cHV0Ej4KDnJlcXVlc3RjaGFyZ2VkGJPQq7kBIAEoDjISLnMzLl'
    'JlcXVlc3RDaGFyZ2VkUg5yZXF1ZXN0Y2hhcmdlZBIvChFyZXN0b3Jlb3V0cHV0cGF0aBjQjaoC'
    'IAEoCVIRcmVzdG9yZW91dHB1dHBhdGg=');

@$core.Deprecated('Use restoreObjectRequestDescriptor instead')
const RestoreObjectRequest$json = {
  '1': 'RestoreObjectRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'restorerequest',
      '3': 66503431,
      '4': 1,
      '5': 11,
      '6': '.s3.RestoreRequest',
      '8': {},
      '10': 'restorerequest'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `RestoreObjectRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreObjectRequestDescriptor = $convert.base64Decode(
    'ChRSZXN0b3JlT2JqZWN0UmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldBJGChFjaG'
    'Vja3N1bWFsZ29yaXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3VtQWxnb3JpdGhtUhFjaGVja3N1'
    'bWFsZ29yaXRobRIzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2'
    'tldG93bmVyEhMKA2tleRiNkutoIAEoCVIDa2V5EjgKDHJlcXVlc3RwYXllchik5uH1ASABKA4y'
    'EC5zMy5SZXF1ZXN0UGF5ZXJSDHJlcXVlc3RwYXllchJDCg5yZXN0b3JlcmVxdWVzdBiHhtsfIA'
    'EoCzISLnMzLlJlc3RvcmVSZXF1ZXN0QgSItRgBUg5yZXN0b3JlcmVxdWVzdBIgCgl2ZXJzaW9u'
    'aWQYm+GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use restoreRequestDescriptor instead')
const RestoreRequest$json = {
  '1': 'RestoreRequest',
  '2': [
    {'1': 'days', '3': 494075051, '4': 1, '5': 5, '10': 'days'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'glacierjobparameters',
      '3': 456012858,
      '4': 1,
      '5': 11,
      '6': '.s3.GlacierJobParameters',
      '10': 'glacierjobparameters'
    },
    {
      '1': 'outputlocation',
      '3': 67991028,
      '4': 1,
      '5': 11,
      '6': '.s3.OutputLocation',
      '10': 'outputlocation'
    },
    {
      '1': 'selectparameters',
      '3': 337564850,
      '4': 1,
      '5': 11,
      '6': '.s3.SelectParameters',
      '10': 'selectparameters'
    },
    {
      '1': 'tier',
      '3': 519596586,
      '4': 1,
      '5': 14,
      '6': '.s3.Tier',
      '10': 'tier'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.s3.RestoreRequestType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `RestoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreRequestDescriptor = $convert.base64Decode(
    'Cg5SZXN0b3JlUmVxdWVzdBIWCgRkYXlzGKv5y+sBIAEoBVIEZGF5cxIjCgtkZXNjcmlwdGlvbh'
    'iK9Pk2IAEoCVILZGVzY3JpcHRpb24SUAoUZ2xhY2llcmpvYnBhcmFtZXRlcnMYuui42QEgASgL'
    'MhguczMuR2xhY2llckpvYlBhcmFtZXRlcnNSFGdsYWNpZXJqb2JwYXJhbWV0ZXJzEj0KDm91dH'
    'B1dGxvY2F0aW9uGPTrtSAgASgLMhIuczMuT3V0cHV0TG9jYXRpb25SDm91dHB1dGxvY2F0aW9u'
    'EkQKEHNlbGVjdHBhcmFtZXRlcnMYsqn7oAEgASgLMhQuczMuU2VsZWN0UGFyYW1ldGVyc1IQc2'
    'VsZWN0cGFyYW1ldGVycxIgCgR0aWVyGKrU4fcBIAEoDjIILnMzLlRpZXJSBHRpZXISLgoEdHlw'
    'ZRjuoNeKASABKA4yFi5zMy5SZXN0b3JlUmVxdWVzdFR5cGVSBHR5cGU=');

@$core.Deprecated('Use restoreStatusDescriptor instead')
const RestoreStatus$json = {
  '1': 'RestoreStatus',
  '2': [
    {
      '1': 'isrestoreinprogress',
      '3': 494957054,
      '4': 1,
      '5': 8,
      '10': 'isrestoreinprogress'
    },
    {
      '1': 'restoreexpirydate',
      '3': 276676935,
      '4': 1,
      '5': 9,
      '10': 'restoreexpirydate'
    },
  ],
};

/// Descriptor for `RestoreStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreStatusDescriptor = $convert.base64Decode(
    'Cg1SZXN0b3JlU3RhdHVzEjQKE2lzcmVzdG9yZWlucHJvZ3Jlc3MY/uOB7AEgASgIUhNpc3Jlc3'
    'RvcmVpbnByb2dyZXNzEjAKEXJlc3RvcmVleHBpcnlkYXRlGMeC94MBIAEoCVIRcmVzdG9yZWV4'
    'cGlyeWRhdGU=');

@$core.Deprecated('Use routingRuleDescriptor instead')
const RoutingRule$json = {
  '1': 'RoutingRule',
  '2': [
    {
      '1': 'condition',
      '3': 212017247,
      '4': 1,
      '5': 11,
      '6': '.s3.Condition',
      '10': 'condition'
    },
    {
      '1': 'redirect',
      '3': 533074130,
      '4': 1,
      '5': 11,
      '6': '.s3.Redirect',
      '10': 'redirect'
    },
  ],
};

/// Descriptor for `RoutingRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List routingRuleDescriptor = $convert.base64Decode(
    'CgtSb3V0aW5nUnVsZRIuCgljb25kaXRpb24Y38CMZSABKAsyDS5zMy5Db25kaXRpb25SCWNvbm'
    'RpdGlvbhIsCghyZWRpcmVjdBjSoZj+ASABKAsyDC5zMy5SZWRpcmVjdFIIcmVkaXJlY3Q=');

@$core.Deprecated('Use s3KeyFilterDescriptor instead')
const S3KeyFilter$json = {
  '1': 'S3KeyFilter',
  '2': [
    {
      '1': 'filterrules',
      '3': 484396335,
      '4': 3,
      '5': 11,
      '6': '.s3.FilterRule',
      '10': 'filterrules'
    },
  ],
};

/// Descriptor for `S3KeyFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3KeyFilterDescriptor = $convert.base64Decode(
    'CgtTM0tleUZpbHRlchI0CgtmaWx0ZXJydWxlcxivmv3mASADKAsyDi5zMy5GaWx0ZXJSdWxlUg'
    'tmaWx0ZXJydWxlcw==');

@$core.Deprecated('Use s3LocationDescriptor instead')
const S3Location$json = {
  '1': 'S3Location',
  '2': [
    {
      '1': 'accesscontrollist',
      '3': 63422661,
      '4': 3,
      '5': 11,
      '6': '.s3.Grant',
      '10': 'accesscontrollist'
    },
    {'1': 'bucketname', '3': 208117045, '4': 1, '5': 9, '10': 'bucketname'},
    {
      '1': 'cannedacl',
      '3': 103981981,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectCannedACL',
      '10': 'cannedacl'
    },
    {
      '1': 'encryption',
      '3': 137702205,
      '4': 1,
      '5': 11,
      '6': '.s3.Encryption',
      '10': 'encryption'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {
      '1': 'tagging',
      '3': 33436541,
      '4': 1,
      '5': 11,
      '6': '.s3.Tagging',
      '10': 'tagging'
    },
    {
      '1': 'usermetadata',
      '3': 518679176,
      '4': 3,
      '5': 11,
      '6': '.s3.MetadataEntry',
      '10': 'usermetadata'
    },
  ],
};

/// Descriptor for `S3Location`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3LocationDescriptor = $convert.base64Decode(
    'CgpTM0xvY2F0aW9uEjoKEWFjY2Vzc2NvbnRyb2xsaXN0GMWBnx4gAygLMgkuczMuR3JhbnRSEW'
    'FjY2Vzc2NvbnRyb2xsaXN0EiEKCmJ1Y2tldG5hbWUYtbqeYyABKAlSCmJ1Y2tldG5hbWUSNAoJ'
    'Y2FubmVkYWNsGJ3HyjEgASgOMhMuczMuT2JqZWN0Q2FubmVkQUNMUgljYW5uZWRhY2wSMQoKZW'
    '5jcnlwdGlvbhi91tRBIAEoCzIOLnMzLkVuY3J5cHRpb25SCmVuY3J5cHRpb24SGgoGcHJlZml4'
    'GOqz04IBIAEoCVIGcHJlZml4EjgKDHN0b3JhZ2VjbGFzcxjHiMS7ASABKA4yEC5zMy5TdG9yYW'
    'dlQ2xhc3NSDHN0b3JhZ2VjbGFzcxIoCgd0YWdnaW5nGP3m+A8gASgLMgsuczMuVGFnZ2luZ1IH'
    'dGFnZ2luZxI5Cgx1c2VybWV0YWRhdGEYiNWp9wEgAygLMhEuczMuTWV0YWRhdGFFbnRyeVIMdX'
    'Nlcm1ldGFkYXRh');

@$core.Deprecated('Use s3TablesDestinationDescriptor instead')
const S3TablesDestination$json = {
  '1': 'S3TablesDestination',
  '2': [
    {
      '1': 'tablebucketarn',
      '3': 457996261,
      '4': 1,
      '5': 9,
      '10': 'tablebucketarn'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `S3TablesDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3TablesDestinationDescriptor = $convert.base64Decode(
    'ChNTM1RhYmxlc0Rlc3RpbmF0aW9uEioKDnRhYmxlYnVja2V0YXJuGOXvsdoBIAEoCVIOdGFibG'
    'VidWNrZXRhcm4SIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use s3TablesDestinationResultDescriptor instead')
const S3TablesDestinationResult$json = {
  '1': 'S3TablesDestinationResult',
  '2': [
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {
      '1': 'tablebucketarn',
      '3': 457996261,
      '4': 1,
      '5': 9,
      '10': 'tablebucketarn'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'tablenamespace',
      '3': 272949507,
      '4': 1,
      '5': 9,
      '10': 'tablenamespace'
    },
  ],
};

/// Descriptor for `S3TablesDestinationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3TablesDestinationResultDescriptor = $convert.base64Decode(
    'ChlTM1RhYmxlc0Rlc3RpbmF0aW9uUmVzdWx0Eh4KCHRhYmxlYXJuGOOA680BIAEoCVIIdGFibG'
    'Vhcm4SKgoOdGFibGVidWNrZXRhcm4Y5e+x2gEgASgJUg50YWJsZWJ1Y2tldGFybhIgCgl0YWJs'
    'ZW5hbWUY3eTagQEgASgJUgl0YWJsZW5hbWUSKgoOdGFibGVuYW1lc3BhY2UYg8KTggEgASgJUg'
    '50YWJsZW5hbWVzcGFjZQ==');

@$core.Deprecated('Use sSEKMSDescriptor instead')
const SSEKMS$json = {
  '1': 'SSEKMS',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `SSEKMS`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sSEKMSDescriptor =
    $convert.base64Decode('CgZTU0VLTVMSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZA==');

@$core.Deprecated('Use sSEKMSEncryptionDescriptor instead')
const SSEKMSEncryption$json = {
  '1': 'SSEKMSEncryption',
  '2': [
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {'1': 'kmskeyarn', '3': 117627377, '4': 1, '5': 9, '10': 'kmskeyarn'},
  ],
};

/// Descriptor for `SSEKMSEncryption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sSEKMSEncryptionDescriptor = $convert.base64Decode(
    'ChBTU0VLTVNFbmNyeXB0aW9uEi4KEGJ1Y2tldGtleWVuYWJsZWQY7KKMzwEgASgIUhBidWNrZX'
    'RrZXllbmFibGVkEh8KCWttc2tleWFybhjxs4s4IAEoCVIJa21za2V5YXJu');

@$core.Deprecated('Use sSES3Descriptor instead')
const SSES3$json = {
  '1': 'SSES3',
};

/// Descriptor for `SSES3`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sSES3Descriptor =
    $convert.base64Decode('CgVTU0VTMw==');

@$core.Deprecated('Use scanRangeDescriptor instead')
const ScanRange$json = {
  '1': 'ScanRange',
  '2': [
    {'1': 'end', '3': 261322315, '4': 1, '5': 3, '10': 'end'},
    {'1': 'start', '3': 182978944, '4': 1, '5': 3, '10': 'start'},
  ],
};

/// Descriptor for `ScanRange`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scanRangeDescriptor = $convert.base64Decode(
    'CglTY2FuUmFuZ2USEwoDZW5kGMvszXwgASgDUgNlbmQSFwoFc3RhcnQYgJOgVyABKANSBXN0YX'
    'J0');

@$core.Deprecated('Use selectObjectContentEventStreamDescriptor instead')
const SelectObjectContentEventStream$json = {
  '1': 'SelectObjectContentEventStream',
  '2': [
    {
      '1': 'cont',
      '3': 242004558,
      '4': 1,
      '5': 11,
      '6': '.s3.ContinuationEvent',
      '10': 'cont'
    },
    {
      '1': 'end',
      '3': 261322315,
      '4': 1,
      '5': 11,
      '6': '.s3.EndEvent',
      '10': 'end'
    },
    {
      '1': 'progress',
      '3': 439787879,
      '4': 1,
      '5': 11,
      '6': '.s3.ProgressEvent',
      '10': 'progress'
    },
    {
      '1': 'records',
      '3': 423557454,
      '4': 1,
      '5': 11,
      '6': '.s3.RecordsEvent',
      '10': 'records'
    },
    {
      '1': 'stats',
      '3': 267161229,
      '4': 1,
      '5': 11,
      '6': '.s3.StatsEvent',
      '10': 'stats'
    },
  ],
};

/// Descriptor for `SelectObjectContentEventStream`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectObjectContentEventStreamDescriptor = $convert.base64Decode(
    'Ch5TZWxlY3RPYmplY3RDb250ZW50RXZlbnRTdHJlYW0SLAoEY29udBjO5LJzIAEoCzIVLnMzLk'
    'NvbnRpbnVhdGlvbkV2ZW50UgRjb250EiEKA2VuZBjL7M18IAEoCzIMLnMzLkVuZEV2ZW50UgNl'
    'bmQSMQoIcHJvZ3Jlc3MY58La0QEgASgLMhEuczMuUHJvZ3Jlc3NFdmVudFIIcHJvZ3Jlc3MSLg'
    'oHcmVjb3JkcxjO8vvJASABKAsyEC5zMy5SZWNvcmRzRXZlbnRSB3JlY29yZHMSJwoFc3RhdHMY'
    'jZ2yfyABKAsyDi5zMy5TdGF0c0V2ZW50UgVzdGF0cw==');

@$core.Deprecated('Use selectObjectContentOutputDescriptor instead')
const SelectObjectContentOutput$json = {
  '1': 'SelectObjectContentOutput',
  '2': [
    {
      '1': 'payload',
      '3': 6526790,
      '4': 1,
      '5': 11,
      '6': '.s3.SelectObjectContentEventStream',
      '8': {},
      '10': 'payload'
    },
  ],
};

/// Descriptor for `SelectObjectContentOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectObjectContentOutputDescriptor =
    $convert.base64Decode(
        'ChlTZWxlY3RPYmplY3RDb250ZW50T3V0cHV0EkUKB3BheWxvYWQYxq6OAyABKAsyIi5zMy5TZW'
        'xlY3RPYmplY3RDb250ZW50RXZlbnRTdHJlYW1CBIi1GAFSB3BheWxvYWQ=');

@$core.Deprecated('Use selectObjectContentRequestDescriptor instead')
const SelectObjectContentRequest$json = {
  '1': 'SelectObjectContentRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'expression', '3': 193051916, '4': 1, '5': 9, '10': 'expression'},
    {
      '1': 'expressiontype',
      '3': 438626996,
      '4': 1,
      '5': 14,
      '6': '.s3.ExpressionType',
      '10': 'expressiontype'
    },
    {
      '1': 'inputserialization',
      '3': 181173046,
      '4': 1,
      '5': 11,
      '6': '.s3.InputSerialization',
      '10': 'inputserialization'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'outputserialization',
      '3': 376385057,
      '4': 1,
      '5': 11,
      '6': '.s3.OutputSerialization',
      '10': 'outputserialization'
    },
    {
      '1': 'requestprogress',
      '3': 409526766,
      '4': 1,
      '5': 11,
      '6': '.s3.RequestProgress',
      '10': 'requestprogress'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {
      '1': 'scanrange',
      '3': 132827166,
      '4': 1,
      '5': 11,
      '6': '.s3.ScanRange',
      '10': 'scanrange'
    },
  ],
};

/// Descriptor for `SelectObjectContentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectObjectContentRequestDescriptor = $convert.base64Decode(
    'ChpTZWxlY3RPYmplY3RDb250ZW50UmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2tldB'
    'IzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEiEK'
    'CmV4cHJlc3Npb24YjPqGXCABKAlSCmV4cHJlc3Npb24SPgoOZXhwcmVzc2lvbnR5cGUYtNWT0Q'
    'EgASgOMhIuczMuRXhwcmVzc2lvblR5cGVSDmV4cHJlc3Npb250eXBlEkkKEmlucHV0c2VyaWFs'
    'aXphdGlvbhi29rFWIAEoCzIWLnMzLklucHV0U2VyaWFsaXphdGlvblISaW5wdXRzZXJpYWxpem'
    'F0aW9uEhMKA2tleRiNkutoIAEoCVIDa2V5Ek0KE291dHB1dHNlcmlhbGl6YXRpb24Yody8swEg'
    'ASgLMhcuczMuT3V0cHV0U2VyaWFsaXphdGlvblITb3V0cHV0c2VyaWFsaXphdGlvbhJBCg9yZX'
    'F1ZXN0cHJvZ3Jlc3MY7sOjwwEgASgLMhMuczMuUmVxdWVzdFByb2dyZXNzUg9yZXF1ZXN0cHJv'
    'Z3Jlc3MSNQoUc3NlY3VzdG9tZXJhbGdvcml0aG0Y0MmBKyABKAlSFHNzZWN1c3RvbWVyYWxnb3'
    'JpdGhtEikKDnNzZWN1c3RvbWVya2V5GJr+9DsgASgJUg5zc2VjdXN0b21lcmtleRIuChFzc2Vj'
    'dXN0b21lcmtleW1kNRjo0RcgASgJUhFzc2VjdXN0b21lcmtleW1kNRIuCglzY2FucmFuZ2UYnp'
    'CrPyABKAsyDS5zMy5TY2FuUmFuZ2VSCXNjYW5yYW5nZQ==');

@$core.Deprecated('Use selectParametersDescriptor instead')
const SelectParameters$json = {
  '1': 'SelectParameters',
  '2': [
    {'1': 'expression', '3': 193051916, '4': 1, '5': 9, '10': 'expression'},
    {
      '1': 'expressiontype',
      '3': 438626996,
      '4': 1,
      '5': 14,
      '6': '.s3.ExpressionType',
      '10': 'expressiontype'
    },
    {
      '1': 'inputserialization',
      '3': 181173046,
      '4': 1,
      '5': 11,
      '6': '.s3.InputSerialization',
      '10': 'inputserialization'
    },
    {
      '1': 'outputserialization',
      '3': 376385057,
      '4': 1,
      '5': 11,
      '6': '.s3.OutputSerialization',
      '10': 'outputserialization'
    },
  ],
};

/// Descriptor for `SelectParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectParametersDescriptor = $convert.base64Decode(
    'ChBTZWxlY3RQYXJhbWV0ZXJzEiEKCmV4cHJlc3Npb24YjPqGXCABKAlSCmV4cHJlc3Npb24SPg'
    'oOZXhwcmVzc2lvbnR5cGUYtNWT0QEgASgOMhIuczMuRXhwcmVzc2lvblR5cGVSDmV4cHJlc3Np'
    'b250eXBlEkkKEmlucHV0c2VyaWFsaXphdGlvbhi29rFWIAEoCzIWLnMzLklucHV0U2VyaWFsaX'
    'phdGlvblISaW5wdXRzZXJpYWxpemF0aW9uEk0KE291dHB1dHNlcmlhbGl6YXRpb24Yody8swEg'
    'ASgLMhcuczMuT3V0cHV0U2VyaWFsaXphdGlvblITb3V0cHV0c2VyaWFsaXphdGlvbg==');

@$core.Deprecated('Use serverSideEncryptionByDefaultDescriptor instead')
const ServerSideEncryptionByDefault$json = {
  '1': 'ServerSideEncryptionByDefault',
  '2': [
    {
      '1': 'kmsmasterkeyid',
      '3': 521472339,
      '4': 1,
      '5': 9,
      '10': 'kmsmasterkeyid'
    },
    {
      '1': 'ssealgorithm',
      '3': 308868112,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'ssealgorithm'
    },
  ],
};

/// Descriptor for `ServerSideEncryptionByDefault`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverSideEncryptionByDefaultDescriptor =
    $convert.base64Decode(
        'Ch1TZXJ2ZXJTaWRlRW5jcnlwdGlvbkJ5RGVmYXVsdBIqCg5rbXNtYXN0ZXJrZXlpZBjTktT4AS'
        'ABKAlSDmttc21hc3RlcmtleWlkEkAKDHNzZWFsZ29yaXRobRiQ6KOTASABKA4yGC5zMy5TZXJ2'
        'ZXJTaWRlRW5jcnlwdGlvblIMc3NlYWxnb3JpdGht');

@$core.Deprecated('Use serverSideEncryptionConfigurationDescriptor instead')
const ServerSideEncryptionConfiguration$json = {
  '1': 'ServerSideEncryptionConfiguration',
  '2': [
    {
      '1': 'rules',
      '3': 42675585,
      '4': 3,
      '5': 11,
      '6': '.s3.ServerSideEncryptionRule',
      '10': 'rules'
    },
  ],
};

/// Descriptor for `ServerSideEncryptionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverSideEncryptionConfigurationDescriptor =
    $convert.base64Decode(
        'CiFTZXJ2ZXJTaWRlRW5jcnlwdGlvbkNvbmZpZ3VyYXRpb24SNQoFcnVsZXMYgdusFCADKAsyHC'
        '5zMy5TZXJ2ZXJTaWRlRW5jcnlwdGlvblJ1bGVSBXJ1bGVz');

@$core.Deprecated('Use serverSideEncryptionRuleDescriptor instead')
const ServerSideEncryptionRule$json = {
  '1': 'ServerSideEncryptionRule',
  '2': [
    {
      '1': 'applyserversideencryptionbydefault',
      '3': 453555041,
      '4': 1,
      '5': 11,
      '6': '.s3.ServerSideEncryptionByDefault',
      '10': 'applyserversideencryptionbydefault'
    },
    {
      '1': 'blockedencryptiontypes',
      '3': 399314360,
      '4': 1,
      '5': 11,
      '6': '.s3.BlockedEncryptionTypes',
      '10': 'blockedencryptiontypes'
    },
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
  ],
};

/// Descriptor for `ServerSideEncryptionRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverSideEncryptionRuleDescriptor = $convert.base64Decode(
    'ChhTZXJ2ZXJTaWRlRW5jcnlwdGlvblJ1bGUSdQoiYXBwbHlzZXJ2ZXJzaWRlZW5jcnlwdGlvbm'
    'J5ZGVmYXVsdBjh5qLYASABKAsyIS5zMy5TZXJ2ZXJTaWRlRW5jcnlwdGlvbkJ5RGVmYXVsdFIi'
    'YXBwbHlzZXJ2ZXJzaWRlZW5jcnlwdGlvbmJ5ZGVmYXVsdBJWChZibG9ja2VkZW5jcnlwdGlvbn'
    'R5cGVzGLibtL4BIAEoCzIaLnMzLkJsb2NrZWRFbmNyeXB0aW9uVHlwZXNSFmJsb2NrZWRlbmNy'
    'eXB0aW9udHlwZXMSLgoQYnVja2V0a2V5ZW5hYmxlZBjsoozPASABKAhSEGJ1Y2tldGtleWVuYW'
    'JsZWQ=');

@$core.Deprecated('Use sessionCredentialsDescriptor instead')
const SessionCredentials$json = {
  '1': 'SessionCredentials',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {
      '1': 'secretaccesskey',
      '3': 172288487,
      '4': 1,
      '5': 9,
      '10': 'secretaccesskey'
    },
    {'1': 'sessiontoken', '3': 211161069, '4': 1, '5': 9, '10': 'sessiontoken'},
  ],
};

/// Descriptor for `SessionCredentials`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionCredentialsDescriptor = $convert.base64Decode(
    'ChJTZXNzaW9uQ3JlZGVudGlhbHMSJAoLYWNjZXNza2V5aWQYoLe32AEgASgJUgthY2Nlc3NrZX'
    'lpZBIhCgpleHBpcmF0aW9uGImpn3UgASgJUgpleHBpcmF0aW9uEisKD3NlY3JldGFjY2Vzc2tl'
    'eRjn05NSIAEoCVIPc2VjcmV0YWNjZXNza2V5EiUKDHNlc3Npb250b2tlbhjtn9hkIAEoCVIMc2'
    'Vzc2lvbnRva2Vu');

@$core.Deprecated('Use simplePrefixDescriptor instead')
const SimplePrefix$json = {
  '1': 'SimplePrefix',
};

/// Descriptor for `SimplePrefix`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List simplePrefixDescriptor =
    $convert.base64Decode('CgxTaW1wbGVQcmVmaXg=');

@$core.Deprecated('Use sourceSelectionCriteriaDescriptor instead')
const SourceSelectionCriteria$json = {
  '1': 'SourceSelectionCriteria',
  '2': [
    {
      '1': 'replicamodifications',
      '3': 340349161,
      '4': 1,
      '5': 11,
      '6': '.s3.ReplicaModifications',
      '10': 'replicamodifications'
    },
    {
      '1': 'ssekmsencryptedobjects',
      '3': 55585072,
      '4': 1,
      '5': 11,
      '6': '.s3.SseKmsEncryptedObjects',
      '10': 'ssekmsencryptedobjects'
    },
  ],
};

/// Descriptor for `SourceSelectionCriteria`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sourceSelectionCriteriaDescriptor = $convert.base64Decode(
    'ChdTb3VyY2VTZWxlY3Rpb25Dcml0ZXJpYRJQChRyZXBsaWNhbW9kaWZpY2F0aW9ucxjpoaWiAS'
    'ABKAsyGC5zMy5SZXBsaWNhTW9kaWZpY2F0aW9uc1IUcmVwbGljYW1vZGlmaWNhdGlvbnMSVQoW'
    'c3Nla21zZW5jcnlwdGVkb2JqZWN0cxiw0sAaIAEoCzIaLnMzLlNzZUttc0VuY3J5cHRlZE9iam'
    'VjdHNSFnNzZWttc2VuY3J5cHRlZG9iamVjdHM=');

@$core.Deprecated('Use sseKmsEncryptedObjectsDescriptor instead')
const SseKmsEncryptedObjects$json = {
  '1': 'SseKmsEncryptedObjects',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.SseKmsEncryptedObjectsStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `SseKmsEncryptedObjects`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sseKmsEncryptedObjectsDescriptor =
    $convert.base64Decode(
        'ChZTc2VLbXNFbmNyeXB0ZWRPYmplY3RzEjsKBnN0YXR1cxiQ5PsCIAEoDjIgLnMzLlNzZUttc0'
        'VuY3J5cHRlZE9iamVjdHNTdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use statsDescriptor instead')
const Stats$json = {
  '1': 'Stats',
  '2': [
    {
      '1': 'bytesprocessed',
      '3': 487068657,
      '4': 1,
      '5': 3,
      '10': 'bytesprocessed'
    },
    {
      '1': 'bytesreturned',
      '3': 121984684,
      '4': 1,
      '5': 3,
      '10': 'bytesreturned'
    },
    {'1': 'bytesscanned', '3': 186950329, '4': 1, '5': 3, '10': 'bytesscanned'},
  ],
};

/// Descriptor for `Stats`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statsDescriptor = $convert.base64Decode(
    'CgVTdGF0cxIqCg5ieXRlc3Byb2Nlc3NlZBjxp6DoASABKANSDmJ5dGVzcHJvY2Vzc2VkEicKDW'
    'J5dGVzcmV0dXJuZWQYrK2VOiABKANSDWJ5dGVzcmV0dXJuZWQSJQoMYnl0ZXNzY2FubmVkGLnF'
    'klkgASgDUgxieXRlc3NjYW5uZWQ=');

@$core.Deprecated('Use statsEventDescriptor instead')
const StatsEvent$json = {
  '1': 'StatsEvent',
  '2': [
    {
      '1': 'details',
      '3': 247611974,
      '4': 1,
      '5': 11,
      '6': '.s3.Stats',
      '10': 'details'
    },
  ],
};

/// Descriptor for `StatsEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statsEventDescriptor = $convert.base64Decode(
    'CgpTdGF0c0V2ZW50EiYKB2RldGFpbHMYxoSJdiABKAsyCS5zMy5TdGF0c1IHZGV0YWlscw==');

@$core.Deprecated('Use storageClassAnalysisDescriptor instead')
const StorageClassAnalysis$json = {
  '1': 'StorageClassAnalysis',
  '2': [
    {
      '1': 'dataexport',
      '3': 235554260,
      '4': 1,
      '5': 11,
      '6': '.s3.StorageClassAnalysisDataExport',
      '10': 'dataexport'
    },
  ],
};

/// Descriptor for `StorageClassAnalysis`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List storageClassAnalysisDescriptor = $convert.base64Decode(
    'ChRTdG9yYWdlQ2xhc3NBbmFseXNpcxJFCgpkYXRhZXhwb3J0GNSLqXAgASgLMiIuczMuU3Rvcm'
    'FnZUNsYXNzQW5hbHlzaXNEYXRhRXhwb3J0UgpkYXRhZXhwb3J0');

@$core.Deprecated('Use storageClassAnalysisDataExportDescriptor instead')
const StorageClassAnalysisDataExport$json = {
  '1': 'StorageClassAnalysisDataExport',
  '2': [
    {
      '1': 'destination',
      '3': 457443680,
      '4': 1,
      '5': 11,
      '6': '.s3.AnalyticsExportDestination',
      '10': 'destination'
    },
    {
      '1': 'outputschemaversion',
      '3': 159901762,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClassAnalysisSchemaVersion',
      '10': 'outputschemaversion'
    },
  ],
};

/// Descriptor for `StorageClassAnalysisDataExport`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List storageClassAnalysisDataExportDescriptor =
    $convert.base64Decode(
        'Ch5TdG9yYWdlQ2xhc3NBbmFseXNpc0RhdGFFeHBvcnQSRAoLZGVzdGluYXRpb24Y4JKQ2gEgAS'
        'gLMh4uczMuQW5hbHl0aWNzRXhwb3J0RGVzdGluYXRpb25SC2Rlc3RpbmF0aW9uEloKE291dHB1'
        'dHNjaGVtYXZlcnNpb24YwtCfTCABKA4yJS5zMy5TdG9yYWdlQ2xhc3NBbmFseXNpc1NjaGVtYV'
        'ZlcnNpb25SE291dHB1dHNjaGVtYXZlcnNpb24=');

@$core.Deprecated('Use tagDescriptor instead')
const Tag$json = {
  '1': 'Tag',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Tag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagDescriptor = $convert.base64Decode(
    'CgNUYWcSEwoDa2V5GI2S62ggASgJUgNrZXkSGAoFdmFsdWUY6/KfigEgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use taggingDescriptor instead')
const Tagging$json = {
  '1': 'Tagging',
  '2': [
    {
      '1': 'tagset',
      '3': 454361330,
      '4': 3,
      '5': 11,
      '6': '.s3.Tag',
      '10': 'tagset'
    },
  ],
};

/// Descriptor for `Tagging`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taggingDescriptor = $convert.base64Decode(
    'CgdUYWdnaW5nEiMKBnRhZ3NldBjygdTYASADKAsyBy5zMy5UYWdSBnRhZ3NldA==');

@$core.Deprecated('Use targetGrantDescriptor instead')
const TargetGrant$json = {
  '1': 'TargetGrant',
  '2': [
    {
      '1': 'grantee',
      '3': 89454056,
      '4': 1,
      '5': 11,
      '6': '.s3.Grantee',
      '10': 'grantee'
    },
    {
      '1': 'permission',
      '3': 465848867,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketLogsPermission',
      '10': 'permission'
    },
  ],
};

/// Descriptor for `TargetGrant`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetGrantDescriptor = $convert.base64Decode(
    'CgtUYXJnZXRHcmFudBIoCgdncmFudGVlGOjr0yogASgLMgsuczMuR3JhbnRlZVIHZ3JhbnRlZR'
    'I8CgpwZXJtaXNzaW9uGKOUkd4BIAEoDjIYLnMzLkJ1Y2tldExvZ3NQZXJtaXNzaW9uUgpwZXJt'
    'aXNzaW9u');

@$core.Deprecated('Use targetObjectKeyFormatDescriptor instead')
const TargetObjectKeyFormat$json = {
  '1': 'TargetObjectKeyFormat',
  '2': [
    {
      '1': 'partitionedprefix',
      '3': 189627261,
      '4': 1,
      '5': 11,
      '6': '.s3.PartitionedPrefix',
      '10': 'partitionedprefix'
    },
    {
      '1': 'simpleprefix',
      '3': 495593044,
      '4': 1,
      '5': 11,
      '6': '.s3.SimplePrefix',
      '10': 'simpleprefix'
    },
  ],
};

/// Descriptor for `TargetObjectKeyFormat`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetObjectKeyFormatDescriptor = $convert.base64Decode(
    'ChVUYXJnZXRPYmplY3RLZXlGb3JtYXQSRgoRcGFydGl0aW9uZWRwcmVmaXgY/fa1WiABKAsyFS'
    '5zMy5QYXJ0aXRpb25lZFByZWZpeFIRcGFydGl0aW9uZWRwcmVmaXgSOAoMc2ltcGxlcHJlZml4'
    'GNTMqOwBIAEoCzIQLnMzLlNpbXBsZVByZWZpeFIMc2ltcGxlcHJlZml4');

@$core.Deprecated('Use tieringDescriptor instead')
const Tiering$json = {
  '1': 'Tiering',
  '2': [
    {
      '1': 'accesstier',
      '3': 263116188,
      '4': 1,
      '5': 14,
      '6': '.s3.IntelligentTieringAccessTier',
      '10': 'accesstier'
    },
    {'1': 'days', '3': 494075051, '4': 1, '5': 5, '10': 'days'},
  ],
};

/// Descriptor for `Tiering`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tieringDescriptor = $convert.base64Decode(
    'CgdUaWVyaW5nEkMKCmFjY2Vzc3RpZXIYnKu7fSABKA4yIC5zMy5JbnRlbGxpZ2VudFRpZXJpbm'
    'dBY2Nlc3NUaWVyUgphY2Nlc3N0aWVyEhYKBGRheXMYq/nL6wEgASgFUgRkYXlz');

@$core.Deprecated('Use tooManyPartsDescriptor instead')
const TooManyParts$json = {
  '1': 'TooManyParts',
};

/// Descriptor for `TooManyParts`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyPartsDescriptor =
    $convert.base64Decode('CgxUb29NYW55UGFydHM=');

@$core.Deprecated('Use topicConfigurationDescriptor instead')
const TopicConfiguration$json = {
  '1': 'TopicConfiguration',
  '2': [
    {
      '1': 'events',
      '3': 3416229,
      '4': 3,
      '5': 14,
      '6': '.s3.Event',
      '10': 'events'
    },
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.s3.NotificationConfigurationFilter',
      '10': 'filter'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `TopicConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List topicConfigurationDescriptor = $convert.base64Decode(
    'ChJUb3BpY0NvbmZpZ3VyYXRpb24SJAoGZXZlbnRzGKXB0AEgAygOMgkuczMuRXZlbnRSBmV2ZW'
    '50cxI/CgZmaWx0ZXIYmIGnpQEgASgLMiMuczMuTm90aWZpY2F0aW9uQ29uZmlndXJhdGlvbkZp'
    'bHRlclIGZmlsdGVyEhIKAmlkGIHyorcBIAEoCVICaWQSHQoIdG9waWNhcm4YnPTODiABKAlSCH'
    'RvcGljYXJu');

@$core.Deprecated('Use transitionDescriptor instead')
const Transition$json = {
  '1': 'Transition',
  '2': [
    {'1': 'date', '3': 458388346, '4': 1, '5': 9, '10': 'date'},
    {'1': 'days', '3': 494075051, '4': 1, '5': 5, '10': 'days'},
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.TransitionStorageClass',
      '10': 'storageclass'
    },
  ],
};

/// Descriptor for `Transition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transitionDescriptor = $convert.base64Decode(
    'CgpUcmFuc2l0aW9uEhYKBGRhdGUY+ubJ2gEgASgJUgRkYXRlEhYKBGRheXMYq/nL6wEgASgFUg'
    'RkYXlzEkIKDHN0b3JhZ2VjbGFzcxjHiMS7ASABKA4yGi5zMy5UcmFuc2l0aW9uU3RvcmFnZUNs'
    'YXNzUgxzdG9yYWdlY2xhc3M=');

@$core.Deprecated(
    'Use updateBucketMetadataInventoryTableConfigurationRequestDescriptor instead')
const UpdateBucketMetadataInventoryTableConfigurationRequest$json = {
  '1': 'UpdateBucketMetadataInventoryTableConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'inventorytableconfiguration',
      '3': 82018446,
      '4': 1,
      '5': 11,
      '6': '.s3.InventoryTableConfigurationUpdates',
      '8': {},
      '10': 'inventorytableconfiguration'
    },
  ],
};

/// Descriptor for `UpdateBucketMetadataInventoryTableConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateBucketMetadataInventoryTableConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CjZVcGRhdGVCdWNrZXRNZXRhZGF0YUludmVudG9yeVRhYmxlQ29uZmlndXJhdGlvblJlcXVlc3'
        'QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZXQSRgoRY2hlY2tzdW1hbGdvcml0aG0YsIHYeiAB'
        'KA4yFS5zMy5DaGVja3N1bUFsZ29yaXRobVIRY2hlY2tzdW1hbGdvcml0aG0SIgoKY29udGVudG'
        '1kNRiL6urxASABKAlSCmNvbnRlbnRtZDUSMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEo'
        'CVITZXhwZWN0ZWRidWNrZXRvd25lchJxChtpbnZlbnRvcnl0YWJsZWNvbmZpZ3VyYXRpb24Yjo'
        'GOJyABKAsyJi5zMy5JbnZlbnRvcnlUYWJsZUNvbmZpZ3VyYXRpb25VcGRhdGVzQgSItRgBUhtp'
        'bnZlbnRvcnl0YWJsZWNvbmZpZ3VyYXRpb24=');

@$core.Deprecated(
    'Use updateBucketMetadataJournalTableConfigurationRequestDescriptor instead')
const UpdateBucketMetadataJournalTableConfigurationRequest$json = {
  '1': 'UpdateBucketMetadataJournalTableConfigurationRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'journaltableconfiguration',
      '3': 70721911,
      '4': 1,
      '5': 11,
      '6': '.s3.JournalTableConfigurationUpdates',
      '8': {},
      '10': 'journaltableconfiguration'
    },
  ],
};

/// Descriptor for `UpdateBucketMetadataJournalTableConfigurationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateBucketMetadataJournalTableConfigurationRequestDescriptor =
    $convert.base64Decode(
        'CjRVcGRhdGVCdWNrZXRNZXRhZGF0YUpvdXJuYWxUYWJsZUNvbmZpZ3VyYXRpb25SZXF1ZXN0Eh'
        'kKBmJ1Y2tldBjY6rgaIAEoCVIGYnVja2V0EkYKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgO'
        'MhUuczMuQ2hlY2tzdW1BbGdvcml0aG1SEWNoZWNrc3VtYWxnb3JpdGhtEiIKCmNvbnRlbnRtZD'
        'UYi+rq8QEgASgJUgpjb250ZW50bWQ1EjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938PiABKAlS'
        'E2V4cGVjdGVkYnVja2V0b3duZXISawoZam91cm5hbHRhYmxlY29uZmlndXJhdGlvbhj3wtwhIA'
        'EoCzIkLnMzLkpvdXJuYWxUYWJsZUNvbmZpZ3VyYXRpb25VcGRhdGVzQgSItRgBUhlqb3VybmFs'
        'dGFibGVjb25maWd1cmF0aW9u');

@$core.Deprecated('Use updateObjectEncryptionRequestDescriptor instead')
const UpdateObjectEncryptionRequest$json = {
  '1': 'UpdateObjectEncryptionRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'objectencryption',
      '3': 64620504,
      '4': 1,
      '5': 11,
      '6': '.s3.ObjectEncryption',
      '8': {},
      '10': 'objectencryption'
    },
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
};

/// Descriptor for `UpdateObjectEncryptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateObjectEncryptionRequestDescriptor = $convert.base64Decode(
    'Ch1VcGRhdGVPYmplY3RFbmNyeXB0aW9uUmVxdWVzdBIZCgZidWNrZXQY2Oq4GiABKAlSBmJ1Y2'
    'tldBJGChFjaGVja3N1bWFsZ29yaXRobRiwgdh6IAEoDjIVLnMzLkNoZWNrc3VtQWxnb3JpdGht'
    'UhFjaGVja3N1bWFsZ29yaXRobRIiCgpjb250ZW50bWQ1GIvq6vEBIAEoCVIKY29udGVudG1kNR'
    'IzChNleHBlY3RlZGJ1Y2tldG93bmVyGKfd/D4gASgJUhNleHBlY3RlZGJ1Y2tldG93bmVyEhMK'
    'A2tleRiNkutoIAEoCVIDa2V5EkkKEG9iamVjdGVuY3J5cHRpb24Y2I/oHiABKAsyFC5zMy5PYm'
    'plY3RFbmNyeXB0aW9uQgSItRgBUhBvYmplY3RlbmNyeXB0aW9uEjgKDHJlcXVlc3RwYXllchik'
    '5uH1ASABKA4yEC5zMy5SZXF1ZXN0UGF5ZXJSDHJlcXVlc3RwYXllchIgCgl2ZXJzaW9uaWQYm+'
    'GZoQEgASgJUgl2ZXJzaW9uaWQ=');

@$core.Deprecated('Use updateObjectEncryptionResponseDescriptor instead')
const UpdateObjectEncryptionResponse$json = {
  '1': 'UpdateObjectEncryptionResponse',
  '2': [
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
  ],
};

/// Descriptor for `UpdateObjectEncryptionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateObjectEncryptionResponseDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVPYmplY3RFbmNyeXB0aW9uUmVzcG9uc2USPgoOcmVxdWVzdGNoYXJnZWQYk9CruQ'
        'EgASgOMhIuczMuUmVxdWVzdENoYXJnZWRSDnJlcXVlc3RjaGFyZ2Vk');

@$core.Deprecated('Use uploadPartCopyOutputDescriptor instead')
const UploadPartCopyOutput$json = {
  '1': 'UploadPartCopyOutput',
  '2': [
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {
      '1': 'copypartresult',
      '3': 405767057,
      '4': 1,
      '5': 11,
      '6': '.s3.CopyPartResult',
      '8': {},
      '10': 'copypartresult'
    },
    {
      '1': 'copysourceversionid',
      '3': 257134375,
      '4': 1,
      '5': 9,
      '10': 'copysourceversionid'
    },
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
  ],
};

/// Descriptor for `UploadPartCopyOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadPartCopyOutputDescriptor = $convert.base64Decode(
    'ChRVcGxvYWRQYXJ0Q29weU91dHB1dBIuChBidWNrZXRrZXllbmFibGVkGOyijM8BIAEoCFIQYn'
    'Vja2V0a2V5ZW5hYmxlZBJECg5jb3B5cGFydHJlc3VsdBiRh77BASABKAsyEi5zMy5Db3B5UGFy'
    'dFJlc3VsdEIEiLUYAVIOY29weXBhcnRyZXN1bHQSMwoTY29weXNvdXJjZXZlcnNpb25pZBinns'
    '56IAEoCVITY29weXNvdXJjZXZlcnNpb25pZBI+Cg5yZXF1ZXN0Y2hhcmdlZBiT0Ku5ASABKA4y'
    'Ei5zMy5SZXF1ZXN0Q2hhcmdlZFIOcmVxdWVzdGNoYXJnZWQSNQoUc3NlY3VzdG9tZXJhbGdvcm'
    'l0aG0Y0MmBKyABKAlSFHNzZWN1c3RvbWVyYWxnb3JpdGhtEi4KEXNzZWN1c3RvbWVya2V5bWQ1'
    'GOjRFyABKAlSEXNzZWN1c3RvbWVya2V5bWQ1EiQKC3NzZWttc2tleWlkGNiI1NQBIAEoCVILc3'
    'Nla21za2V5aWQSTwoUc2VydmVyc2lkZWVuY3J5cHRpb24YsY6fBCABKA4yGC5zMy5TZXJ2ZXJT'
    'aWRlRW5jcnlwdGlvblIUc2VydmVyc2lkZWVuY3J5cHRpb24=');

@$core.Deprecated('Use uploadPartCopyRequestDescriptor instead')
const UploadPartCopyRequest$json = {
  '1': 'UploadPartCopyRequest',
  '2': [
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {'1': 'copysource', '3': 152315650, '4': 1, '5': 9, '10': 'copysource'},
    {
      '1': 'copysourceifmatch',
      '3': 512868332,
      '4': 1,
      '5': 9,
      '10': 'copysourceifmatch'
    },
    {
      '1': 'copysourceifmodifiedsince',
      '3': 301576082,
      '4': 1,
      '5': 9,
      '10': 'copysourceifmodifiedsince'
    },
    {
      '1': 'copysourceifnonematch',
      '3': 237403010,
      '4': 1,
      '5': 9,
      '10': 'copysourceifnonematch'
    },
    {
      '1': 'copysourceifunmodifiedsince',
      '3': 130556417,
      '4': 1,
      '5': 9,
      '10': 'copysourceifunmodifiedsince'
    },
    {
      '1': 'copysourcerange',
      '3': 51610607,
      '4': 1,
      '5': 9,
      '10': 'copysourcerange'
    },
    {
      '1': 'copysourcessecustomeralgorithm',
      '3': 465017156,
      '4': 1,
      '5': 9,
      '10': 'copysourcessecustomeralgorithm'
    },
    {
      '1': 'copysourcessecustomerkey',
      '3': 289313102,
      '4': 1,
      '5': 9,
      '10': 'copysourcessecustomerkey'
    },
    {
      '1': 'copysourcessecustomerkeymd5',
      '3': 408962492,
      '4': 1,
      '5': 9,
      '10': 'copysourcessecustomerkeymd5'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'expectedsourcebucketowner',
      '3': 283238582,
      '4': 1,
      '5': 9,
      '10': 'expectedsourcebucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'partnumber', '3': 372082310, '4': 1, '5': 5, '10': 'partnumber'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'uploadid', '3': 449040722, '4': 1, '5': 9, '10': 'uploadid'},
  ],
};

/// Descriptor for `UploadPartCopyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadPartCopyRequestDescriptor = $convert.base64Decode(
    'ChVVcGxvYWRQYXJ0Q29weVJlcXVlc3QSGQoGYnVja2V0GNjquBogASgJUgZidWNrZXQSIQoKY2'
    '9weXNvdXJjZRiCztBIIAEoCVIKY29weXNvdXJjZRIwChFjb3B5c291cmNlaWZtYXRjaBjs/8b0'
    'ASABKAlSEWNvcHlzb3VyY2VpZm1hdGNoEkAKGWNvcHlzb3VyY2VpZm1vZGlmaWVkc2luY2UYkt'
    '/mjwEgASgJUhljb3B5c291cmNlaWZtb2RpZmllZHNpbmNlEjcKFWNvcHlzb3VyY2VpZm5vbmVt'
    'YXRjaBiC95lxIAEoCVIVY29weXNvdXJjZWlmbm9uZW1hdGNoEkMKG2NvcHlzb3VyY2VpZnVubW'
    '9kaWZpZWRzaW5jZRiBxKA+IAEoCVIbY29weXNvdXJjZWlmdW5tb2RpZmllZHNpbmNlEisKD2Nv'
    'cHlzb3VyY2VyYW5nZRjvh84YIAEoCVIPY29weXNvdXJjZXJhbmdlEkoKHmNvcHlzb3VyY2Vzc2'
    'VjdXN0b21lcmFsZ29yaXRobRjEst7dASABKAlSHmNvcHlzb3VyY2Vzc2VjdXN0b21lcmFsZ29y'
    'aXRobRI+Chhjb3B5c291cmNlc3NlY3VzdG9tZXJrZXkYzqL6iQEgASgJUhhjb3B5c291cmNlc3'
    'NlY3VzdG9tZXJrZXkSRAobY29weXNvdXJjZXNzZWN1c3RvbWVya2V5bWQ1GLyLgcMBIAEoCVIb'
    'Y29weXNvdXJjZXNzZWN1c3RvbWVya2V5bWQ1EjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938Pi'
    'ABKAlSE2V4cGVjdGVkYnVja2V0b3duZXISQAoZZXhwZWN0ZWRzb3VyY2VidWNrZXRvd25lchi2'
    'wYeHASABKAlSGWV4cGVjdGVkc291cmNlYnVja2V0b3duZXISEwoDa2V5GI2S62ggASgJUgNrZX'
    'kSIgoKcGFydG51bWJlchiGjbaxASABKAVSCnBhcnRudW1iZXISOAoMcmVxdWVzdHBheWVyGKTm'
    '4fUBIAEoDjIQLnMzLlJlcXVlc3RQYXllclIMcmVxdWVzdHBheWVyEjUKFHNzZWN1c3RvbWVyYW'
    'xnb3JpdGhtGNDJgSsgASgJUhRzc2VjdXN0b21lcmFsZ29yaXRobRIpCg5zc2VjdXN0b21lcmtl'
    'eRia/vQ7IAEoCVIOc3NlY3VzdG9tZXJrZXkSLgoRc3NlY3VzdG9tZXJrZXltZDUY6NEXIAEoCV'
    'IRc3NlY3VzdG9tZXJrZXltZDUSHgoIdXBsb2FkaWQY0qKP1gEgASgJUgh1cGxvYWRpZA==');

@$core.Deprecated('Use uploadPartOutputDescriptor instead')
const UploadPartOutput$json = {
  '1': 'UploadPartOutput',
  '2': [
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
  ],
};

/// Descriptor for `UploadPartOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadPartOutputDescriptor = $convert.base64Decode(
    'ChBVcGxvYWRQYXJ0T3V0cHV0Ei4KEGJ1Y2tldGtleWVuYWJsZWQY7KKMzwEgASgIUhBidWNrZX'
    'RrZXllbmFibGVkEicKDWNoZWNrc3VtY3JjMzIYwqPNMyABKAlSDWNoZWNrc3VtY3JjMzISKQoO'
    'Y2hlY2tzdW1jcmMzMmMYp5+lTCABKAlSDmNoZWNrc3VtY3JjMzJjEi8KEWNoZWNrc3VtY3JjNj'
    'Rudm1lGM28izggASgJUhFjaGVja3N1bWNyYzY0bnZtZRImCgxjaGVja3N1bXNoYTEYxOzgigEg'
    'ASgJUgxjaGVja3N1bXNoYTESKQoOY2hlY2tzdW1zaGEyNTYYvvncRCABKAlSDmNoZWNrc3Vtc2'
    'hhMjU2EhYKBGV0YWcYgd+zlQEgASgJUgRldGFnEj4KDnJlcXVlc3RjaGFyZ2VkGJPQq7kBIAEo'
    'DjISLnMzLlJlcXVlc3RDaGFyZ2VkUg5yZXF1ZXN0Y2hhcmdlZBI1ChRzc2VjdXN0b21lcmFsZ2'
    '9yaXRobRjQyYErIAEoCVIUc3NlY3VzdG9tZXJhbGdvcml0aG0SLgoRc3NlY3VzdG9tZXJrZXlt'
    'ZDUY6NEXIAEoCVIRc3NlY3VzdG9tZXJrZXltZDUSJAoLc3Nla21za2V5aWQY2IjU1AEgASgJUg'
    'tzc2VrbXNrZXlpZBJPChRzZXJ2ZXJzaWRlZW5jcnlwdGlvbhixjp8EIAEoDjIYLnMzLlNlcnZl'
    'clNpZGVFbmNyeXB0aW9uUhRzZXJ2ZXJzaWRlZW5jcnlwdGlvbg==');

@$core.Deprecated('Use uploadPartRequestDescriptor instead')
const UploadPartRequest$json = {
  '1': 'UploadPartRequest',
  '2': [
    {'1': 'body', '3': 42602646, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {'1': 'bucket', '3': 55457112, '4': 1, '5': 9, '10': 'bucket'},
    {
      '1': 'checksumalgorithm',
      '3': 257294512,
      '4': 1,
      '5': 14,
      '6': '.s3.ChecksumAlgorithm',
      '10': 'checksumalgorithm'
    },
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'contentlength',
      '3': 227596631,
      '4': 1,
      '5': 3,
      '10': 'contentlength'
    },
    {'1': 'contentmd5', '3': 507163915, '4': 1, '5': 9, '10': 'contentmd5'},
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'partnumber', '3': 372082310, '4': 1, '5': 5, '10': 'partnumber'},
    {
      '1': 'requestpayer',
      '3': 515404580,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestPayer',
      '10': 'requestpayer'
    },
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkey',
      '3': 125648666,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkey'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'uploadid', '3': 449040722, '4': 1, '5': 9, '10': 'uploadid'},
  ],
};

/// Descriptor for `UploadPartRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List uploadPartRequestDescriptor = $convert.base64Decode(
    'ChFVcGxvYWRQYXJ0UmVxdWVzdBIbCgRib2R5GJahqBQgASgMQgSItRgBUgRib2R5EhkKBmJ1Y2'
    'tldBjY6rgaIAEoCVIGYnVja2V0EkYKEWNoZWNrc3VtYWxnb3JpdGhtGLCB2HogASgOMhUuczMu'
    'Q2hlY2tzdW1BbGdvcml0aG1SEWNoZWNrc3VtYWxnb3JpdGhtEicKDWNoZWNrc3VtY3JjMzIYwq'
    'PNMyABKAlSDWNoZWNrc3VtY3JjMzISKQoOY2hlY2tzdW1jcmMzMmMYp5+lTCABKAlSDmNoZWNr'
    'c3VtY3JjMzJjEi8KEWNoZWNrc3VtY3JjNjRudm1lGM28izggASgJUhFjaGVja3N1bWNyYzY0bn'
    'ZtZRImCgxjaGVja3N1bXNoYTEYxOzgigEgASgJUgxjaGVja3N1bXNoYTESKQoOY2hlY2tzdW1z'
    'aGEyNTYYvvncRCABKAlSDmNoZWNrc3Vtc2hhMjU2EicKDWNvbnRlbnRsZW5ndGgY17LDbCABKA'
    'NSDWNvbnRlbnRsZW5ndGgSIgoKY29udGVudG1kNRiL6urxASABKAlSCmNvbnRlbnRtZDUSMwoT'
    'ZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVITZXhwZWN0ZWRidWNrZXRvd25lchITCgNrZX'
    'kYjZLraCABKAlSA2tleRIiCgpwYXJ0bnVtYmVyGIaNtrEBIAEoBVIKcGFydG51bWJlchI4Cgxy'
    'ZXF1ZXN0cGF5ZXIYpObh9QEgASgOMhAuczMuUmVxdWVzdFBheWVyUgxyZXF1ZXN0cGF5ZXISNQ'
    'oUc3NlY3VzdG9tZXJhbGdvcml0aG0Y0MmBKyABKAlSFHNzZWN1c3RvbWVyYWxnb3JpdGhtEikK'
    'DnNzZWN1c3RvbWVya2V5GJr+9DsgASgJUg5zc2VjdXN0b21lcmtleRIuChFzc2VjdXN0b21lcm'
    'tleW1kNRjo0RcgASgJUhFzc2VjdXN0b21lcmtleW1kNRIeCgh1cGxvYWRpZBjSoo/WASABKAlS'
    'CHVwbG9hZGlk');

@$core.Deprecated('Use versioningConfigurationDescriptor instead')
const VersioningConfiguration$json = {
  '1': 'VersioningConfiguration',
  '2': [
    {
      '1': 'mfadelete',
      '3': 354119319,
      '4': 1,
      '5': 14,
      '6': '.s3.MFADelete',
      '10': 'mfadelete'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.s3.BucketVersioningStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `VersioningConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List versioningConfigurationDescriptor = $convert.base64Decode(
    'ChdWZXJzaW9uaW5nQ29uZmlndXJhdGlvbhIvCgltZmFkZWxldGUYl93tqAEgASgOMg0uczMuTU'
    'ZBRGVsZXRlUgltZmFkZWxldGUSNQoGc3RhdHVzGJDk+wIgASgOMhouczMuQnVja2V0VmVyc2lv'
    'bmluZ1N0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use websiteConfigurationDescriptor instead')
const WebsiteConfiguration$json = {
  '1': 'WebsiteConfiguration',
  '2': [
    {
      '1': 'errordocument',
      '3': 329194857,
      '4': 1,
      '5': 11,
      '6': '.s3.ErrorDocument',
      '10': 'errordocument'
    },
    {
      '1': 'indexdocument',
      '3': 23264623,
      '4': 1,
      '5': 11,
      '6': '.s3.IndexDocument',
      '10': 'indexdocument'
    },
    {
      '1': 'redirectallrequeststo',
      '3': 317142372,
      '4': 1,
      '5': 11,
      '6': '.s3.RedirectAllRequestsTo',
      '10': 'redirectallrequeststo'
    },
    {
      '1': 'routingrules',
      '3': 83486933,
      '4': 3,
      '5': 11,
      '6': '.s3.RoutingRule',
      '10': 'routingrules'
    },
  ],
};

/// Descriptor for `WebsiteConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List websiteConfigurationDescriptor = $convert.base64Decode(
    'ChRXZWJzaXRlQ29uZmlndXJhdGlvbhI7Cg1lcnJvcmRvY3VtZW50GOm6/JwBIAEoCzIRLnMzLk'
    'Vycm9yRG9jdW1lbnRSDWVycm9yZG9jdW1lbnQSOgoNaW5kZXhkb2N1bWVudBjv+osLIAEoCzIR'
    'LnMzLkluZGV4RG9jdW1lbnRSDWluZGV4ZG9jdW1lbnQSUwoVcmVkaXJlY3RhbGxyZXF1ZXN0c3'
    'RvGOTqnJcBIAEoCzIZLnMzLlJlZGlyZWN0QWxsUmVxdWVzdHNUb1IVcmVkaXJlY3RhbGxyZXF1'
    'ZXN0c3RvEjYKDHJvdXRpbmdydWxlcxjV0ecnIAMoCzIPLnMzLlJvdXRpbmdSdWxlUgxyb3V0aW'
    '5ncnVsZXM=');

@$core.Deprecated('Use writeGetObjectResponseRequestDescriptor instead')
const WriteGetObjectResponseRequest$json = {
  '1': 'WriteGetObjectResponseRequest',
  '2': [
    {'1': 'acceptranges', '3': 464620960, '4': 1, '5': 9, '10': 'acceptranges'},
    {'1': 'body', '3': 42602646, '4': 1, '5': 12, '8': {}, '10': 'body'},
    {
      '1': 'bucketkeyenabled',
      '3': 434311532,
      '4': 1,
      '5': 8,
      '10': 'bucketkeyenabled'
    },
    {'1': 'cachecontrol', '3': 288966655, '4': 1, '5': 9, '10': 'cachecontrol'},
    {
      '1': 'checksumcrc32',
      '3': 108220866,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32'
    },
    {
      '1': 'checksumcrc32c',
      '3': 159993767,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc32c'
    },
    {
      '1': 'checksumcrc64nvme',
      '3': 117628493,
      '4': 1,
      '5': 9,
      '10': 'checksumcrc64nvme'
    },
    {'1': 'checksumsha1', '3': 290993732, '4': 1, '5': 9, '10': 'checksumsha1'},
    {
      '1': 'checksumsha256',
      '3': 144129214,
      '4': 1,
      '5': 9,
      '10': 'checksumsha256'
    },
    {
      '1': 'contentdisposition',
      '3': 120040130,
      '4': 1,
      '5': 9,
      '10': 'contentdisposition'
    },
    {
      '1': 'contentencoding',
      '3': 317106228,
      '4': 1,
      '5': 9,
      '10': 'contentencoding'
    },
    {
      '1': 'contentlanguage',
      '3': 108485649,
      '4': 1,
      '5': 9,
      '10': 'contentlanguage'
    },
    {
      '1': 'contentlength',
      '3': 227596631,
      '4': 1,
      '5': 3,
      '10': 'contentlength'
    },
    {'1': 'contentrange', '3': 11089360, '4': 1, '5': 9, '10': 'contentrange'},
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'deletemarker', '3': 5472257, '4': 1, '5': 8, '10': 'deletemarker'},
    {'1': 'etag', '3': 313323393, '4': 1, '5': 9, '10': 'etag'},
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {'1': 'expires', '3': 128582948, '4': 1, '5': 9, '10': 'expires'},
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 3,
      '5': 11,
      '6': '.s3.WriteGetObjectResponseRequest.MetadataEntry',
      '10': 'metadata'
    },
    {'1': 'missingmeta', '3': 79140523, '4': 1, '5': 5, '10': 'missingmeta'},
    {
      '1': 'objectlocklegalholdstatus',
      '3': 536561974,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockLegalHoldStatus',
      '10': 'objectlocklegalholdstatus'
    },
    {
      '1': 'objectlockmode',
      '3': 189255203,
      '4': 1,
      '5': 14,
      '6': '.s3.ObjectLockMode',
      '10': 'objectlockmode'
    },
    {
      '1': 'objectlockretainuntildate',
      '3': 264584249,
      '4': 1,
      '5': 9,
      '10': 'objectlockretainuntildate'
    },
    {'1': 'partscount', '3': 154996373, '4': 1, '5': 5, '10': 'partscount'},
    {
      '1': 'replicationstatus',
      '3': 529093900,
      '4': 1,
      '5': 14,
      '6': '.s3.ReplicationStatus',
      '10': 'replicationstatus'
    },
    {
      '1': 'requestcharged',
      '3': 388687891,
      '4': 1,
      '5': 14,
      '6': '.s3.RequestCharged',
      '10': 'requestcharged'
    },
    {'1': 'requestroute', '3': 41111156, '4': 1, '5': 9, '10': 'requestroute'},
    {'1': 'requesttoken', '3': 262268424, '4': 1, '5': 9, '10': 'requesttoken'},
    {'1': 'restore', '3': 267943794, '4': 1, '5': 9, '10': 'restore'},
    {
      '1': 'ssecustomeralgorithm',
      '3': 90203344,
      '4': 1,
      '5': 9,
      '10': 'ssecustomeralgorithm'
    },
    {
      '1': 'ssecustomerkeymd5',
      '3': 387304,
      '4': 1,
      '5': 9,
      '10': 'ssecustomerkeymd5'
    },
    {'1': 'ssekmskeyid', '3': 445973592, '4': 1, '5': 9, '10': 'ssekmskeyid'},
    {
      '1': 'serversideencryption',
      '3': 8898353,
      '4': 1,
      '5': 14,
      '6': '.s3.ServerSideEncryption',
      '10': 'serversideencryption'
    },
    {'1': 'statuscode', '3': 303830783, '4': 1, '5': 5, '10': 'statuscode'},
    {
      '1': 'storageclass',
      '3': 393282631,
      '4': 1,
      '5': 14,
      '6': '.s3.StorageClass',
      '10': 'storageclass'
    },
    {'1': 'tagcount', '3': 339592595, '4': 1, '5': 5, '10': 'tagcount'},
    {'1': 'versionid', '3': 338063515, '4': 1, '5': 9, '10': 'versionid'},
  ],
  '3': [WriteGetObjectResponseRequest_MetadataEntry$json],
};

@$core.Deprecated('Use writeGetObjectResponseRequestDescriptor instead')
const WriteGetObjectResponseRequest_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `WriteGetObjectResponseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List writeGetObjectResponseRequestDescriptor = $convert.base64Decode(
    'Ch1Xcml0ZUdldE9iamVjdFJlc3BvbnNlUmVxdWVzdBImCgxhY2NlcHRyYW5nZXMYoJvG3QEgAS'
    'gJUgxhY2NlcHRyYW5nZXMSGwoEYm9keRiWoagUIAEoDEIEiLUYAVIEYm9keRIuChBidWNrZXRr'
    'ZXllbmFibGVkGOyijM8BIAEoCFIQYnVja2V0a2V5ZW5hYmxlZBImCgxjYWNoZWNvbnRyb2wY/4'
    '/liQEgASgJUgxjYWNoZWNvbnRyb2wSJwoNY2hlY2tzdW1jcmMzMhjCo80zIAEoCVINY2hlY2tz'
    'dW1jcmMzMhIpCg5jaGVja3N1bWNyYzMyYxinn6VMIAEoCVIOY2hlY2tzdW1jcmMzMmMSLwoRY2'
    'hlY2tzdW1jcmM2NG52bWUYzbyLOCABKAlSEWNoZWNrc3VtY3JjNjRudm1lEiYKDGNoZWNrc3Vt'
    'c2hhMRjE7OCKASABKAlSDGNoZWNrc3Vtc2hhMRIpCg5jaGVja3N1bXNoYTI1Nhi++dxEIAEoCV'
    'IOY2hlY2tzdW1zaGEyNTYSMQoSY29udGVudGRpc3Bvc2l0aW9uGMLVnjkgASgJUhJjb250ZW50'
    'ZGlzcG9zaXRpb24SLAoPY29udGVudGVuY29kaW5nGLTQmpcBIAEoCVIPY29udGVudGVuY29kaW'
    '5nEisKD2NvbnRlbnRsYW5ndWFnZRiRuN0zIAEoCVIPY29udGVudGxhbmd1YWdlEicKDWNvbnRl'
    'bnRsZW5ndGgY17LDbCABKANSDWNvbnRlbnRsZW5ndGgSJQoMY29udGVudHJhbmdlGNDrpAUgAS'
    'gJUgxjb250ZW50cmFuZ2USJAoLY29udGVudHR5cGUYk9XongEgASgJUgtjb250ZW50dHlwZRIl'
    'CgxkZWxldGVtYXJrZXIYgYDOAiABKAhSDGRlbGV0ZW1hcmtlchIWCgRldGFnGIHfs5UBIAEoCV'
    'IEZXRhZxIfCgllcnJvcmNvZGUYmdbDECABKAlSCWVycm9yY29kZRImCgxlcnJvcm1lc3NhZ2UY'
    'qYqr9wEgASgJUgxlcnJvcm1lc3NhZ2USIQoKZXhwaXJhdGlvbhiJqZ91IAEoCVIKZXhwaXJhdG'
    'lvbhIbCgdleHBpcmVzGKSKqD0gASgJUgdleHBpcmVzEiYKDGxhc3Rtb2RpZmllZBinnPzOASAB'
    'KAlSDGxhc3Rtb2RpZmllZBJPCghtZXRhZGF0YRjh4o/gASADKAsyLy5zMy5Xcml0ZUdldE9iam'
    'VjdFJlc3BvbnNlUmVxdWVzdC5NZXRhZGF0YUVudHJ5UghtZXRhZGF0YRIjCgttaXNzaW5nbWV0'
    'YRirrd4lIAEoBVILbWlzc2luZ21ldGESXwoZb2JqZWN0bG9ja2xlZ2FsaG9sZHN0YXR1cxi2ku'
    '3/ASABKA4yHS5zMy5PYmplY3RMb2NrTGVnYWxIb2xkU3RhdHVzUhlvYmplY3Rsb2NrbGVnYWxo'
    'b2xkc3RhdHVzEj0KDm9iamVjdGxvY2ttb2RlGKOcn1ogASgOMhIuczMuT2JqZWN0TG9ja01vZG'
    'VSDm9iamVjdGxvY2ttb2RlEj8KGW9iamVjdGxvY2tyZXRhaW51bnRpbGRhdGUYufiUfiABKAlS'
    'GW9iamVjdGxvY2tyZXRhaW51bnRpbGRhdGUSIQoKcGFydHNjb3VudBiVnfRJIAEoBVIKcGFydH'
    'Njb3VudBJHChFyZXBsaWNhdGlvbnN0YXR1cxiMqqX8ASABKA4yFS5zMy5SZXBsaWNhdGlvblN0'
    'YXR1c1IRcmVwbGljYXRpb25zdGF0dXMSPgoOcmVxdWVzdGNoYXJnZWQYk9CruQEgASgOMhIucz'
    'MuUmVxdWVzdENoYXJnZWRSDnJlcXVlc3RjaGFyZ2VkEiUKDHJlcXVlc3Ryb3V0ZRj0nM0TIAEo'
    'CVIMcmVxdWVzdHJvdXRlEiUKDHJlcXVlc3R0b2tlbhiIzId9IAEoCVIMcmVxdWVzdHRva2VuEh'
    'sKB3Jlc3RvcmUY8v7hfyABKAlSB3Jlc3RvcmUSNQoUc3NlY3VzdG9tZXJhbGdvcml0aG0Y0MmB'
    'KyABKAlSFHNzZWN1c3RvbWVyYWxnb3JpdGhtEi4KEXNzZWN1c3RvbWVya2V5bWQ1GOjRFyABKA'
    'lSEXNzZWN1c3RvbWVya2V5bWQ1EiQKC3NzZWttc2tleWlkGNiI1NQBIAEoCVILc3Nla21za2V5'
    'aWQSTwoUc2VydmVyc2lkZWVuY3J5cHRpb24YsY6fBCABKA4yGC5zMy5TZXJ2ZXJTaWRlRW5jcn'
    'lwdGlvblIUc2VydmVyc2lkZWVuY3J5cHRpb24SIgoKc3RhdHVzY29kZRj/rfCQASABKAVSCnN0'
    'YXR1c2NvZGUSOAoMc3RvcmFnZWNsYXNzGMeIxLsBIAEoDjIQLnMzLlN0b3JhZ2VDbGFzc1IMc3'
    'RvcmFnZWNsYXNzEh4KCHRhZ2NvdW50GJOL96EBIAEoBVIIdGFnY291bnQSIAoJdmVyc2lvbmlk'
    'GJvhmaEBIAEoCVIJdmVyc2lvbmlkGjsKDU1ldGFkYXRhRW50cnkSEAoDa2V5GAEgASgJUgNrZX'
    'kSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');
