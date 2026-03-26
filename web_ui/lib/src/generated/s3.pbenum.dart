// This is a generated file - do not edit.
//
// Generated from s3.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AnalyticsS3ExportFileFormat extends $pb.ProtobufEnum {
  static const AnalyticsS3ExportFileFormat ANALYTICS_S3_EXPORT_FILE_FORMAT_CSV =
      AnalyticsS3ExportFileFormat._(
          0, _omitEnumNames ? '' : 'ANALYTICS_S3_EXPORT_FILE_FORMAT_CSV');

  static const $core.List<AnalyticsS3ExportFileFormat> values =
      <AnalyticsS3ExportFileFormat>[
    ANALYTICS_S3_EXPORT_FILE_FORMAT_CSV,
  ];

  static final $core.List<AnalyticsS3ExportFileFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static AnalyticsS3ExportFileFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AnalyticsS3ExportFileFormat._(super.value, super.name);
}

class ArchiveStatus extends $pb.ProtobufEnum {
  static const ArchiveStatus ARCHIVE_STATUS_DEEP_ARCHIVE_ACCESS =
      ArchiveStatus._(
          0, _omitEnumNames ? '' : 'ARCHIVE_STATUS_DEEP_ARCHIVE_ACCESS');
  static const ArchiveStatus ARCHIVE_STATUS_ARCHIVE_ACCESS =
      ArchiveStatus._(1, _omitEnumNames ? '' : 'ARCHIVE_STATUS_ARCHIVE_ACCESS');

  static const $core.List<ArchiveStatus> values = <ArchiveStatus>[
    ARCHIVE_STATUS_DEEP_ARCHIVE_ACCESS,
    ARCHIVE_STATUS_ARCHIVE_ACCESS,
  ];

  static final $core.List<ArchiveStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ArchiveStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ArchiveStatus._(super.value, super.name);
}

class BucketAbacStatus extends $pb.ProtobufEnum {
  static const BucketAbacStatus BUCKET_ABAC_STATUS_DISABLED =
      BucketAbacStatus._(
          0, _omitEnumNames ? '' : 'BUCKET_ABAC_STATUS_DISABLED');
  static const BucketAbacStatus BUCKET_ABAC_STATUS_ENABLED =
      BucketAbacStatus._(1, _omitEnumNames ? '' : 'BUCKET_ABAC_STATUS_ENABLED');

  static const $core.List<BucketAbacStatus> values = <BucketAbacStatus>[
    BUCKET_ABAC_STATUS_DISABLED,
    BUCKET_ABAC_STATUS_ENABLED,
  ];

  static final $core.List<BucketAbacStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static BucketAbacStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BucketAbacStatus._(super.value, super.name);
}

class BucketAccelerateStatus extends $pb.ProtobufEnum {
  static const BucketAccelerateStatus BUCKET_ACCELERATE_STATUS_SUSPENDED =
      BucketAccelerateStatus._(
          0, _omitEnumNames ? '' : 'BUCKET_ACCELERATE_STATUS_SUSPENDED');
  static const BucketAccelerateStatus BUCKET_ACCELERATE_STATUS_ENABLED =
      BucketAccelerateStatus._(
          1, _omitEnumNames ? '' : 'BUCKET_ACCELERATE_STATUS_ENABLED');

  static const $core.List<BucketAccelerateStatus> values =
      <BucketAccelerateStatus>[
    BUCKET_ACCELERATE_STATUS_SUSPENDED,
    BUCKET_ACCELERATE_STATUS_ENABLED,
  ];

  static final $core.List<BucketAccelerateStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static BucketAccelerateStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BucketAccelerateStatus._(super.value, super.name);
}

class BucketCannedACL extends $pb.ProtobufEnum {
  static const BucketCannedACL BUCKET_CANNED_A_C_L_AUTHENTICATED_READ =
      BucketCannedACL._(
          0, _omitEnumNames ? '' : 'BUCKET_CANNED_A_C_L_AUTHENTICATED_READ');
  static const BucketCannedACL BUCKET_CANNED_A_C_L_PRIVATE =
      BucketCannedACL._(1, _omitEnumNames ? '' : 'BUCKET_CANNED_A_C_L_PRIVATE');
  static const BucketCannedACL BUCKET_CANNED_A_C_L_PUBLIC_READ =
      BucketCannedACL._(
          2, _omitEnumNames ? '' : 'BUCKET_CANNED_A_C_L_PUBLIC_READ');
  static const BucketCannedACL BUCKET_CANNED_A_C_L_PUBLIC_READ_WRITE =
      BucketCannedACL._(
          3, _omitEnumNames ? '' : 'BUCKET_CANNED_A_C_L_PUBLIC_READ_WRITE');

  static const $core.List<BucketCannedACL> values = <BucketCannedACL>[
    BUCKET_CANNED_A_C_L_AUTHENTICATED_READ,
    BUCKET_CANNED_A_C_L_PRIVATE,
    BUCKET_CANNED_A_C_L_PUBLIC_READ,
    BUCKET_CANNED_A_C_L_PUBLIC_READ_WRITE,
  ];

  static final $core.List<BucketCannedACL?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static BucketCannedACL? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BucketCannedACL._(super.value, super.name);
}

class BucketLocationConstraint extends $pb.ProtobufEnum {
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_3 = BucketLocationConstraint._(
          0, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_3');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_EU_WEST_3 =
      BucketLocationConstraint._(
          1, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU_WEST_3');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_US_GOV_EAST_1 = BucketLocationConstraint._(
          2, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_US_GOV_EAST_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_1 = BucketLocationConstraint._(
          3, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_IL_CENTRAL_1 = BucketLocationConstraint._(
          4, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_IL_CENTRAL_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_EU_NORTH_1 =
      BucketLocationConstraint._(
          5, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU_NORTH_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_CN_NORTH_1 =
      BucketLocationConstraint._(
          6, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_CN_NORTH_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_5 = BucketLocationConstraint._(
          7, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_5');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_2 = BucketLocationConstraint._(
          8, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_2');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_US_WEST_1 =
      BucketLocationConstraint._(
          9, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_US_WEST_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_AP_SOUTH_1 =
      BucketLocationConstraint._(
          10, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTH_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_2 = BucketLocationConstraint._(11,
          _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_2');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_EU_SOUTH_1 =
      BucketLocationConstraint._(
          12, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU_SOUTH_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_EU_WEST_2 =
      BucketLocationConstraint._(
          13, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU_WEST_2');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_US_WEST_2 =
      BucketLocationConstraint._(
          14, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_US_WEST_2');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_ME_CENTRAL_1 = BucketLocationConstraint._(
          15, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_ME_CENTRAL_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_CN_NORTHWEST_1 = BucketLocationConstraint._(16,
          _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_CN_NORTHWEST_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_EU_SOUTH_2 =
      BucketLocationConstraint._(
          17, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU_SOUTH_2');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_AP_EAST_1 =
      BucketLocationConstraint._(
          18, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_EAST_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_EU_CENTRAL_1 = BucketLocationConstraint._(
          19, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU_CENTRAL_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_4 = BucketLocationConstraint._(20,
          _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_4');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_CA_CENTRAL_1 = BucketLocationConstraint._(
          21, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_CA_CENTRAL_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_SA_EAST_1 =
      BucketLocationConstraint._(
          22, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_SA_EAST_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_AP_SOUTH_2 =
      BucketLocationConstraint._(
          23, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTH_2');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_US_GOV_WEST_1 = BucketLocationConstraint._(
          24, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_US_GOV_WEST_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_EU_CENTRAL_2 = BucketLocationConstraint._(
          25, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU_CENTRAL_2');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_EU =
      BucketLocationConstraint._(
          26, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_1 = BucketLocationConstraint._(27,
          _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_AF_SOUTH_1 =
      BucketLocationConstraint._(
          28, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AF_SOUTH_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_ME_SOUTH_1 =
      BucketLocationConstraint._(
          29, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_ME_SOUTH_1');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_EU_WEST_1 =
      BucketLocationConstraint._(
          30, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_EU_WEST_1');
  static const BucketLocationConstraint
      BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_3 = BucketLocationConstraint._(31,
          _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_3');
  static const BucketLocationConstraint BUCKET_LOCATION_CONSTRAINT_US_EAST_2 =
      BucketLocationConstraint._(
          32, _omitEnumNames ? '' : 'BUCKET_LOCATION_CONSTRAINT_US_EAST_2');

  static const $core.List<BucketLocationConstraint> values =
      <BucketLocationConstraint>[
    BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_3,
    BUCKET_LOCATION_CONSTRAINT_EU_WEST_3,
    BUCKET_LOCATION_CONSTRAINT_US_GOV_EAST_1,
    BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_1,
    BUCKET_LOCATION_CONSTRAINT_IL_CENTRAL_1,
    BUCKET_LOCATION_CONSTRAINT_EU_NORTH_1,
    BUCKET_LOCATION_CONSTRAINT_CN_NORTH_1,
    BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_5,
    BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_2,
    BUCKET_LOCATION_CONSTRAINT_US_WEST_1,
    BUCKET_LOCATION_CONSTRAINT_AP_SOUTH_1,
    BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_2,
    BUCKET_LOCATION_CONSTRAINT_EU_SOUTH_1,
    BUCKET_LOCATION_CONSTRAINT_EU_WEST_2,
    BUCKET_LOCATION_CONSTRAINT_US_WEST_2,
    BUCKET_LOCATION_CONSTRAINT_ME_CENTRAL_1,
    BUCKET_LOCATION_CONSTRAINT_CN_NORTHWEST_1,
    BUCKET_LOCATION_CONSTRAINT_EU_SOUTH_2,
    BUCKET_LOCATION_CONSTRAINT_AP_EAST_1,
    BUCKET_LOCATION_CONSTRAINT_EU_CENTRAL_1,
    BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_4,
    BUCKET_LOCATION_CONSTRAINT_CA_CENTRAL_1,
    BUCKET_LOCATION_CONSTRAINT_SA_EAST_1,
    BUCKET_LOCATION_CONSTRAINT_AP_SOUTH_2,
    BUCKET_LOCATION_CONSTRAINT_US_GOV_WEST_1,
    BUCKET_LOCATION_CONSTRAINT_EU_CENTRAL_2,
    BUCKET_LOCATION_CONSTRAINT_EU,
    BUCKET_LOCATION_CONSTRAINT_AP_SOUTHEAST_1,
    BUCKET_LOCATION_CONSTRAINT_AF_SOUTH_1,
    BUCKET_LOCATION_CONSTRAINT_ME_SOUTH_1,
    BUCKET_LOCATION_CONSTRAINT_EU_WEST_1,
    BUCKET_LOCATION_CONSTRAINT_AP_NORTHEAST_3,
    BUCKET_LOCATION_CONSTRAINT_US_EAST_2,
  ];

  static final $core.List<BucketLocationConstraint?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 32);
  static BucketLocationConstraint? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BucketLocationConstraint._(super.value, super.name);
}

class BucketLogsPermission extends $pb.ProtobufEnum {
  static const BucketLogsPermission BUCKET_LOGS_PERMISSION_READ =
      BucketLogsPermission._(
          0, _omitEnumNames ? '' : 'BUCKET_LOGS_PERMISSION_READ');
  static const BucketLogsPermission BUCKET_LOGS_PERMISSION_WRITE =
      BucketLogsPermission._(
          1, _omitEnumNames ? '' : 'BUCKET_LOGS_PERMISSION_WRITE');
  static const BucketLogsPermission BUCKET_LOGS_PERMISSION_FULL_CONTROL =
      BucketLogsPermission._(
          2, _omitEnumNames ? '' : 'BUCKET_LOGS_PERMISSION_FULL_CONTROL');

  static const $core.List<BucketLogsPermission> values = <BucketLogsPermission>[
    BUCKET_LOGS_PERMISSION_READ,
    BUCKET_LOGS_PERMISSION_WRITE,
    BUCKET_LOGS_PERMISSION_FULL_CONTROL,
  ];

  static final $core.List<BucketLogsPermission?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static BucketLogsPermission? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BucketLogsPermission._(super.value, super.name);
}

class BucketType extends $pb.ProtobufEnum {
  static const BucketType BUCKET_TYPE_DIRECTORY =
      BucketType._(0, _omitEnumNames ? '' : 'BUCKET_TYPE_DIRECTORY');

  static const $core.List<BucketType> values = <BucketType>[
    BUCKET_TYPE_DIRECTORY,
  ];

  static final $core.List<BucketType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static BucketType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BucketType._(super.value, super.name);
}

class BucketVersioningStatus extends $pb.ProtobufEnum {
  static const BucketVersioningStatus BUCKET_VERSIONING_STATUS_SUSPENDED =
      BucketVersioningStatus._(
          0, _omitEnumNames ? '' : 'BUCKET_VERSIONING_STATUS_SUSPENDED');
  static const BucketVersioningStatus BUCKET_VERSIONING_STATUS_ENABLED =
      BucketVersioningStatus._(
          1, _omitEnumNames ? '' : 'BUCKET_VERSIONING_STATUS_ENABLED');

  static const $core.List<BucketVersioningStatus> values =
      <BucketVersioningStatus>[
    BUCKET_VERSIONING_STATUS_SUSPENDED,
    BUCKET_VERSIONING_STATUS_ENABLED,
  ];

  static final $core.List<BucketVersioningStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static BucketVersioningStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BucketVersioningStatus._(super.value, super.name);
}

class ChecksumAlgorithm extends $pb.ProtobufEnum {
  static const ChecksumAlgorithm CHECKSUM_ALGORITHM_SHA256 =
      ChecksumAlgorithm._(0, _omitEnumNames ? '' : 'CHECKSUM_ALGORITHM_SHA256');
  static const ChecksumAlgorithm CHECKSUM_ALGORITHM_SHA1 =
      ChecksumAlgorithm._(1, _omitEnumNames ? '' : 'CHECKSUM_ALGORITHM_SHA1');
  static const ChecksumAlgorithm CHECKSUM_ALGORITHM_CRC32 =
      ChecksumAlgorithm._(2, _omitEnumNames ? '' : 'CHECKSUM_ALGORITHM_CRC32');
  static const ChecksumAlgorithm CHECKSUM_ALGORITHM_CRC64NVME =
      ChecksumAlgorithm._(
          3, _omitEnumNames ? '' : 'CHECKSUM_ALGORITHM_CRC64NVME');
  static const ChecksumAlgorithm CHECKSUM_ALGORITHM_CRC32C =
      ChecksumAlgorithm._(4, _omitEnumNames ? '' : 'CHECKSUM_ALGORITHM_CRC32C');

  static const $core.List<ChecksumAlgorithm> values = <ChecksumAlgorithm>[
    CHECKSUM_ALGORITHM_SHA256,
    CHECKSUM_ALGORITHM_SHA1,
    CHECKSUM_ALGORITHM_CRC32,
    CHECKSUM_ALGORITHM_CRC64NVME,
    CHECKSUM_ALGORITHM_CRC32C,
  ];

  static final $core.List<ChecksumAlgorithm?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ChecksumAlgorithm? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChecksumAlgorithm._(super.value, super.name);
}

class ChecksumMode extends $pb.ProtobufEnum {
  static const ChecksumMode CHECKSUM_MODE_ENABLED =
      ChecksumMode._(0, _omitEnumNames ? '' : 'CHECKSUM_MODE_ENABLED');

  static const $core.List<ChecksumMode> values = <ChecksumMode>[
    CHECKSUM_MODE_ENABLED,
  ];

  static final $core.List<ChecksumMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ChecksumMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChecksumMode._(super.value, super.name);
}

class ChecksumType extends $pb.ProtobufEnum {
  static const ChecksumType CHECKSUM_TYPE_FULL_OBJECT =
      ChecksumType._(0, _omitEnumNames ? '' : 'CHECKSUM_TYPE_FULL_OBJECT');
  static const ChecksumType CHECKSUM_TYPE_COMPOSITE =
      ChecksumType._(1, _omitEnumNames ? '' : 'CHECKSUM_TYPE_COMPOSITE');

  static const $core.List<ChecksumType> values = <ChecksumType>[
    CHECKSUM_TYPE_FULL_OBJECT,
    CHECKSUM_TYPE_COMPOSITE,
  ];

  static final $core.List<ChecksumType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ChecksumType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChecksumType._(super.value, super.name);
}

class CompressionType extends $pb.ProtobufEnum {
  static const CompressionType COMPRESSION_TYPE_NONE =
      CompressionType._(0, _omitEnumNames ? '' : 'COMPRESSION_TYPE_NONE');
  static const CompressionType COMPRESSION_TYPE_GZIP =
      CompressionType._(1, _omitEnumNames ? '' : 'COMPRESSION_TYPE_GZIP');
  static const CompressionType COMPRESSION_TYPE_BZIP2 =
      CompressionType._(2, _omitEnumNames ? '' : 'COMPRESSION_TYPE_BZIP2');

  static const $core.List<CompressionType> values = <CompressionType>[
    COMPRESSION_TYPE_NONE,
    COMPRESSION_TYPE_GZIP,
    COMPRESSION_TYPE_BZIP2,
  ];

  static final $core.List<CompressionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static CompressionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CompressionType._(super.value, super.name);
}

class DataRedundancy extends $pb.ProtobufEnum {
  static const DataRedundancy DATA_REDUNDANCY_SINGLEAVAILABILITYZONE =
      DataRedundancy._(
          0, _omitEnumNames ? '' : 'DATA_REDUNDANCY_SINGLEAVAILABILITYZONE');
  static const DataRedundancy DATA_REDUNDANCY_SINGLELOCALZONE =
      DataRedundancy._(
          1, _omitEnumNames ? '' : 'DATA_REDUNDANCY_SINGLELOCALZONE');

  static const $core.List<DataRedundancy> values = <DataRedundancy>[
    DATA_REDUNDANCY_SINGLEAVAILABILITYZONE,
    DATA_REDUNDANCY_SINGLELOCALZONE,
  ];

  static final $core.List<DataRedundancy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DataRedundancy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DataRedundancy._(super.value, super.name);
}

class DeleteMarkerReplicationStatus extends $pb.ProtobufEnum {
  static const DeleteMarkerReplicationStatus
      DELETE_MARKER_REPLICATION_STATUS_DISABLED =
      DeleteMarkerReplicationStatus._(
          0, _omitEnumNames ? '' : 'DELETE_MARKER_REPLICATION_STATUS_DISABLED');
  static const DeleteMarkerReplicationStatus
      DELETE_MARKER_REPLICATION_STATUS_ENABLED =
      DeleteMarkerReplicationStatus._(
          1, _omitEnumNames ? '' : 'DELETE_MARKER_REPLICATION_STATUS_ENABLED');

  static const $core.List<DeleteMarkerReplicationStatus> values =
      <DeleteMarkerReplicationStatus>[
    DELETE_MARKER_REPLICATION_STATUS_DISABLED,
    DELETE_MARKER_REPLICATION_STATUS_ENABLED,
  ];

  static final $core.List<DeleteMarkerReplicationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DeleteMarkerReplicationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeleteMarkerReplicationStatus._(super.value, super.name);
}

class EncodingType extends $pb.ProtobufEnum {
  static const EncodingType ENCODING_TYPE_URL =
      EncodingType._(0, _omitEnumNames ? '' : 'ENCODING_TYPE_URL');

  static const $core.List<EncodingType> values = <EncodingType>[
    ENCODING_TYPE_URL,
  ];

  static final $core.List<EncodingType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static EncodingType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EncodingType._(super.value, super.name);
}

class EncryptionType extends $pb.ProtobufEnum {
  static const EncryptionType ENCRYPTION_TYPE_NONE =
      EncryptionType._(0, _omitEnumNames ? '' : 'ENCRYPTION_TYPE_NONE');
  static const EncryptionType ENCRYPTION_TYPE_SSE_C =
      EncryptionType._(1, _omitEnumNames ? '' : 'ENCRYPTION_TYPE_SSE_C');

  static const $core.List<EncryptionType> values = <EncryptionType>[
    ENCRYPTION_TYPE_NONE,
    ENCRYPTION_TYPE_SSE_C,
  ];

  static final $core.List<EncryptionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static EncryptionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EncryptionType._(super.value, super.name);
}

class Event extends $pb.ProtobufEnum {
  static const Event EVENT_S3_OBJECTRESTORE_ =
      Event._(0, _omitEnumNames ? '' : 'EVENT_S3_OBJECTRESTORE_');
  static const Event EVENT_S3_OBJECTTAGGING_ =
      Event._(1, _omitEnumNames ? '' : 'EVENT_S3_OBJECTTAGGING_');
  static const Event EVENT_S3_REPLICATION_OPERATIONFAILEDREPLICATION = Event._(
      2,
      _omitEnumNames ? '' : 'EVENT_S3_REPLICATION_OPERATIONFAILEDREPLICATION');
  static const Event EVENT_S3_OBJECTCREATED_COPY =
      Event._(3, _omitEnumNames ? '' : 'EVENT_S3_OBJECTCREATED_COPY');
  static const Event EVENT_S3_OBJECTCREATED_POST =
      Event._(4, _omitEnumNames ? '' : 'EVENT_S3_OBJECTCREATED_POST');
  static const Event EVENT_S3_REDUCEDREDUNDANCYLOSTOBJECT =
      Event._(5, _omitEnumNames ? '' : 'EVENT_S3_REDUCEDREDUNDANCYLOSTOBJECT');
  static const Event EVENT_S3_OBJECTCREATED_PUT =
      Event._(6, _omitEnumNames ? '' : 'EVENT_S3_OBJECTCREATED_PUT');
  static const Event EVENT_S3_OBJECTRESTORE_COMPLETED =
      Event._(7, _omitEnumNames ? '' : 'EVENT_S3_OBJECTRESTORE_COMPLETED');
  static const Event EVENT_S3_LIFECYCLEEXPIRATION_DELETE =
      Event._(8, _omitEnumNames ? '' : 'EVENT_S3_LIFECYCLEEXPIRATION_DELETE');
  static const Event EVENT_S3_OBJECTTAGGING_PUT =
      Event._(9, _omitEnumNames ? '' : 'EVENT_S3_OBJECTTAGGING_PUT');
  static const Event EVENT_S3_OBJECTREMOVED_DELETEMARKERCREATED = Event._(
      10, _omitEnumNames ? '' : 'EVENT_S3_OBJECTREMOVED_DELETEMARKERCREATED');
  static const Event EVENT_S3_OBJECTTAGGING_DELETE =
      Event._(11, _omitEnumNames ? '' : 'EVENT_S3_OBJECTTAGGING_DELETE');
  static const Event EVENT_S3_REPLICATION_OPERATIONMISSEDTHRESHOLD = Event._(12,
      _omitEnumNames ? '' : 'EVENT_S3_REPLICATION_OPERATIONMISSEDTHRESHOLD');
  static const Event EVENT_S3_LIFECYCLETRANSITION =
      Event._(13, _omitEnumNames ? '' : 'EVENT_S3_LIFECYCLETRANSITION');
  static const Event EVENT_S3_OBJECTCREATED_COMPLETEMULTIPARTUPLOAD = Event._(
      14,
      _omitEnumNames ? '' : 'EVENT_S3_OBJECTCREATED_COMPLETEMULTIPARTUPLOAD');
  static const Event EVENT_S3_INTELLIGENTTIERING =
      Event._(15, _omitEnumNames ? '' : 'EVENT_S3_INTELLIGENTTIERING');
  static const Event EVENT_S3_OBJECTCREATED_ =
      Event._(16, _omitEnumNames ? '' : 'EVENT_S3_OBJECTCREATED_');
  static const Event EVENT_S3_REPLICATION_OPERATIONREPLICATEDAFTERTHRESHOLD =
      Event._(
          17,
          _omitEnumNames
              ? ''
              : 'EVENT_S3_REPLICATION_OPERATIONREPLICATEDAFTERTHRESHOLD');
  static const Event EVENT_S3_OBJECTRESTORE_DELETE =
      Event._(18, _omitEnumNames ? '' : 'EVENT_S3_OBJECTRESTORE_DELETE');
  static const Event EVENT_S3_LIFECYCLEEXPIRATION_ =
      Event._(19, _omitEnumNames ? '' : 'EVENT_S3_LIFECYCLEEXPIRATION_');
  static const Event EVENT_S3_REPLICATION_OPERATIONNOTTRACKED = Event._(
      20, _omitEnumNames ? '' : 'EVENT_S3_REPLICATION_OPERATIONNOTTRACKED');
  static const Event EVENT_S3_LIFECYCLEEXPIRATION_DELETEMARKERCREATED = Event._(
      21,
      _omitEnumNames ? '' : 'EVENT_S3_LIFECYCLEEXPIRATION_DELETEMARKERCREATED');
  static const Event EVENT_S3_REPLICATION_ =
      Event._(22, _omitEnumNames ? '' : 'EVENT_S3_REPLICATION_');
  static const Event EVENT_S3_OBJECTREMOVED_ =
      Event._(23, _omitEnumNames ? '' : 'EVENT_S3_OBJECTREMOVED_');
  static const Event EVENT_S3_OBJECTRESTORE_POST =
      Event._(24, _omitEnumNames ? '' : 'EVENT_S3_OBJECTRESTORE_POST');
  static const Event EVENT_S3_OBJECTACL_PUT =
      Event._(25, _omitEnumNames ? '' : 'EVENT_S3_OBJECTACL_PUT');
  static const Event EVENT_S3_OBJECTREMOVED_DELETE =
      Event._(26, _omitEnumNames ? '' : 'EVENT_S3_OBJECTREMOVED_DELETE');

  static const $core.List<Event> values = <Event>[
    EVENT_S3_OBJECTRESTORE_,
    EVENT_S3_OBJECTTAGGING_,
    EVENT_S3_REPLICATION_OPERATIONFAILEDREPLICATION,
    EVENT_S3_OBJECTCREATED_COPY,
    EVENT_S3_OBJECTCREATED_POST,
    EVENT_S3_REDUCEDREDUNDANCYLOSTOBJECT,
    EVENT_S3_OBJECTCREATED_PUT,
    EVENT_S3_OBJECTRESTORE_COMPLETED,
    EVENT_S3_LIFECYCLEEXPIRATION_DELETE,
    EVENT_S3_OBJECTTAGGING_PUT,
    EVENT_S3_OBJECTREMOVED_DELETEMARKERCREATED,
    EVENT_S3_OBJECTTAGGING_DELETE,
    EVENT_S3_REPLICATION_OPERATIONMISSEDTHRESHOLD,
    EVENT_S3_LIFECYCLETRANSITION,
    EVENT_S3_OBJECTCREATED_COMPLETEMULTIPARTUPLOAD,
    EVENT_S3_INTELLIGENTTIERING,
    EVENT_S3_OBJECTCREATED_,
    EVENT_S3_REPLICATION_OPERATIONREPLICATEDAFTERTHRESHOLD,
    EVENT_S3_OBJECTRESTORE_DELETE,
    EVENT_S3_LIFECYCLEEXPIRATION_,
    EVENT_S3_REPLICATION_OPERATIONNOTTRACKED,
    EVENT_S3_LIFECYCLEEXPIRATION_DELETEMARKERCREATED,
    EVENT_S3_REPLICATION_,
    EVENT_S3_OBJECTREMOVED_,
    EVENT_S3_OBJECTRESTORE_POST,
    EVENT_S3_OBJECTACL_PUT,
    EVENT_S3_OBJECTREMOVED_DELETE,
  ];

  static final $core.List<Event?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 26);
  static Event? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Event._(super.value, super.name);
}

class ExistingObjectReplicationStatus extends $pb.ProtobufEnum {
  static const ExistingObjectReplicationStatus
      EXISTING_OBJECT_REPLICATION_STATUS_DISABLED =
      ExistingObjectReplicationStatus._(0,
          _omitEnumNames ? '' : 'EXISTING_OBJECT_REPLICATION_STATUS_DISABLED');
  static const ExistingObjectReplicationStatus
      EXISTING_OBJECT_REPLICATION_STATUS_ENABLED =
      ExistingObjectReplicationStatus._(1,
          _omitEnumNames ? '' : 'EXISTING_OBJECT_REPLICATION_STATUS_ENABLED');

  static const $core.List<ExistingObjectReplicationStatus> values =
      <ExistingObjectReplicationStatus>[
    EXISTING_OBJECT_REPLICATION_STATUS_DISABLED,
    EXISTING_OBJECT_REPLICATION_STATUS_ENABLED,
  ];

  static final $core.List<ExistingObjectReplicationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExistingObjectReplicationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExistingObjectReplicationStatus._(super.value, super.name);
}

class ExpirationState extends $pb.ProtobufEnum {
  static const ExpirationState EXPIRATION_STATE_DISABLED =
      ExpirationState._(0, _omitEnumNames ? '' : 'EXPIRATION_STATE_DISABLED');
  static const ExpirationState EXPIRATION_STATE_ENABLED =
      ExpirationState._(1, _omitEnumNames ? '' : 'EXPIRATION_STATE_ENABLED');

  static const $core.List<ExpirationState> values = <ExpirationState>[
    EXPIRATION_STATE_DISABLED,
    EXPIRATION_STATE_ENABLED,
  ];

  static final $core.List<ExpirationState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExpirationState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExpirationState._(super.value, super.name);
}

class ExpirationStatus extends $pb.ProtobufEnum {
  static const ExpirationStatus EXPIRATION_STATUS_DISABLED =
      ExpirationStatus._(0, _omitEnumNames ? '' : 'EXPIRATION_STATUS_DISABLED');
  static const ExpirationStatus EXPIRATION_STATUS_ENABLED =
      ExpirationStatus._(1, _omitEnumNames ? '' : 'EXPIRATION_STATUS_ENABLED');

  static const $core.List<ExpirationStatus> values = <ExpirationStatus>[
    EXPIRATION_STATUS_DISABLED,
    EXPIRATION_STATUS_ENABLED,
  ];

  static final $core.List<ExpirationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExpirationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExpirationStatus._(super.value, super.name);
}

class ExpressionType extends $pb.ProtobufEnum {
  static const ExpressionType EXPRESSION_TYPE_SQL =
      ExpressionType._(0, _omitEnumNames ? '' : 'EXPRESSION_TYPE_SQL');

  static const $core.List<ExpressionType> values = <ExpressionType>[
    EXPRESSION_TYPE_SQL,
  ];

  static final $core.List<ExpressionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ExpressionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExpressionType._(super.value, super.name);
}

class FileHeaderInfo extends $pb.ProtobufEnum {
  static const FileHeaderInfo FILE_HEADER_INFO_IGNORE =
      FileHeaderInfo._(0, _omitEnumNames ? '' : 'FILE_HEADER_INFO_IGNORE');
  static const FileHeaderInfo FILE_HEADER_INFO_NONE =
      FileHeaderInfo._(1, _omitEnumNames ? '' : 'FILE_HEADER_INFO_NONE');
  static const FileHeaderInfo FILE_HEADER_INFO_USE =
      FileHeaderInfo._(2, _omitEnumNames ? '' : 'FILE_HEADER_INFO_USE');

  static const $core.List<FileHeaderInfo> values = <FileHeaderInfo>[
    FILE_HEADER_INFO_IGNORE,
    FILE_HEADER_INFO_NONE,
    FILE_HEADER_INFO_USE,
  ];

  static final $core.List<FileHeaderInfo?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static FileHeaderInfo? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FileHeaderInfo._(super.value, super.name);
}

class FilterRuleName extends $pb.ProtobufEnum {
  static const FilterRuleName FILTER_RULE_NAME_SUFFIX =
      FilterRuleName._(0, _omitEnumNames ? '' : 'FILTER_RULE_NAME_SUFFIX');
  static const FilterRuleName FILTER_RULE_NAME_PREFIX =
      FilterRuleName._(1, _omitEnumNames ? '' : 'FILTER_RULE_NAME_PREFIX');

  static const $core.List<FilterRuleName> values = <FilterRuleName>[
    FILTER_RULE_NAME_SUFFIX,
    FILTER_RULE_NAME_PREFIX,
  ];

  static final $core.List<FilterRuleName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FilterRuleName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FilterRuleName._(super.value, super.name);
}

class IntelligentTieringAccessTier extends $pb.ProtobufEnum {
  static const IntelligentTieringAccessTier
      INTELLIGENT_TIERING_ACCESS_TIER_DEEP_ARCHIVE_ACCESS =
      IntelligentTieringAccessTier._(
          0,
          _omitEnumNames
              ? ''
              : 'INTELLIGENT_TIERING_ACCESS_TIER_DEEP_ARCHIVE_ACCESS');
  static const IntelligentTieringAccessTier
      INTELLIGENT_TIERING_ACCESS_TIER_ARCHIVE_ACCESS =
      IntelligentTieringAccessTier._(
          1,
          _omitEnumNames
              ? ''
              : 'INTELLIGENT_TIERING_ACCESS_TIER_ARCHIVE_ACCESS');

  static const $core.List<IntelligentTieringAccessTier> values =
      <IntelligentTieringAccessTier>[
    INTELLIGENT_TIERING_ACCESS_TIER_DEEP_ARCHIVE_ACCESS,
    INTELLIGENT_TIERING_ACCESS_TIER_ARCHIVE_ACCESS,
  ];

  static final $core.List<IntelligentTieringAccessTier?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IntelligentTieringAccessTier? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IntelligentTieringAccessTier._(super.value, super.name);
}

class IntelligentTieringStatus extends $pb.ProtobufEnum {
  static const IntelligentTieringStatus INTELLIGENT_TIERING_STATUS_DISABLED =
      IntelligentTieringStatus._(
          0, _omitEnumNames ? '' : 'INTELLIGENT_TIERING_STATUS_DISABLED');
  static const IntelligentTieringStatus INTELLIGENT_TIERING_STATUS_ENABLED =
      IntelligentTieringStatus._(
          1, _omitEnumNames ? '' : 'INTELLIGENT_TIERING_STATUS_ENABLED');

  static const $core.List<IntelligentTieringStatus> values =
      <IntelligentTieringStatus>[
    INTELLIGENT_TIERING_STATUS_DISABLED,
    INTELLIGENT_TIERING_STATUS_ENABLED,
  ];

  static final $core.List<IntelligentTieringStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IntelligentTieringStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IntelligentTieringStatus._(super.value, super.name);
}

class InventoryConfigurationState extends $pb.ProtobufEnum {
  static const InventoryConfigurationState
      INVENTORY_CONFIGURATION_STATE_DISABLED = InventoryConfigurationState._(
          0, _omitEnumNames ? '' : 'INVENTORY_CONFIGURATION_STATE_DISABLED');
  static const InventoryConfigurationState
      INVENTORY_CONFIGURATION_STATE_ENABLED = InventoryConfigurationState._(
          1, _omitEnumNames ? '' : 'INVENTORY_CONFIGURATION_STATE_ENABLED');

  static const $core.List<InventoryConfigurationState> values =
      <InventoryConfigurationState>[
    INVENTORY_CONFIGURATION_STATE_DISABLED,
    INVENTORY_CONFIGURATION_STATE_ENABLED,
  ];

  static final $core.List<InventoryConfigurationState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InventoryConfigurationState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventoryConfigurationState._(super.value, super.name);
}

class InventoryFormat extends $pb.ProtobufEnum {
  static const InventoryFormat INVENTORY_FORMAT_ORC =
      InventoryFormat._(0, _omitEnumNames ? '' : 'INVENTORY_FORMAT_ORC');
  static const InventoryFormat INVENTORY_FORMAT_CSV =
      InventoryFormat._(1, _omitEnumNames ? '' : 'INVENTORY_FORMAT_CSV');
  static const InventoryFormat INVENTORY_FORMAT_PARQUET =
      InventoryFormat._(2, _omitEnumNames ? '' : 'INVENTORY_FORMAT_PARQUET');

  static const $core.List<InventoryFormat> values = <InventoryFormat>[
    INVENTORY_FORMAT_ORC,
    INVENTORY_FORMAT_CSV,
    INVENTORY_FORMAT_PARQUET,
  ];

  static final $core.List<InventoryFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static InventoryFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventoryFormat._(super.value, super.name);
}

class InventoryFrequency extends $pb.ProtobufEnum {
  static const InventoryFrequency INVENTORY_FREQUENCY_DAILY =
      InventoryFrequency._(
          0, _omitEnumNames ? '' : 'INVENTORY_FREQUENCY_DAILY');
  static const InventoryFrequency INVENTORY_FREQUENCY_WEEKLY =
      InventoryFrequency._(
          1, _omitEnumNames ? '' : 'INVENTORY_FREQUENCY_WEEKLY');

  static const $core.List<InventoryFrequency> values = <InventoryFrequency>[
    INVENTORY_FREQUENCY_DAILY,
    INVENTORY_FREQUENCY_WEEKLY,
  ];

  static final $core.List<InventoryFrequency?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InventoryFrequency? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventoryFrequency._(super.value, super.name);
}

class InventoryIncludedObjectVersions extends $pb.ProtobufEnum {
  static const InventoryIncludedObjectVersions
      INVENTORY_INCLUDED_OBJECT_VERSIONS_CURRENT =
      InventoryIncludedObjectVersions._(0,
          _omitEnumNames ? '' : 'INVENTORY_INCLUDED_OBJECT_VERSIONS_CURRENT');
  static const InventoryIncludedObjectVersions
      INVENTORY_INCLUDED_OBJECT_VERSIONS_ALL =
      InventoryIncludedObjectVersions._(
          1, _omitEnumNames ? '' : 'INVENTORY_INCLUDED_OBJECT_VERSIONS_ALL');

  static const $core.List<InventoryIncludedObjectVersions> values =
      <InventoryIncludedObjectVersions>[
    INVENTORY_INCLUDED_OBJECT_VERSIONS_CURRENT,
    INVENTORY_INCLUDED_OBJECT_VERSIONS_ALL,
  ];

  static final $core.List<InventoryIncludedObjectVersions?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InventoryIncludedObjectVersions? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventoryIncludedObjectVersions._(super.value, super.name);
}

class InventoryOptionalField extends $pb.ProtobufEnum {
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_LIFECYCLEEXPIRATIONDATE =
      InventoryOptionalField._(
          0,
          _omitEnumNames
              ? ''
              : 'INVENTORY_OPTIONAL_FIELD_LIFECYCLEEXPIRATIONDATE');
  static const InventoryOptionalField INVENTORY_OPTIONAL_FIELD_BUCKETKEYSTATUS =
      InventoryOptionalField._(
          1, _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_BUCKETKEYSTATUS');
  static const InventoryOptionalField INVENTORY_OPTIONAL_FIELD_STORAGECLASS =
      InventoryOptionalField._(
          2, _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_STORAGECLASS');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_CHECKSUMALGORITHM = InventoryOptionalField._(3,
          _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_CHECKSUMALGORITHM');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_REPLICATIONSTATUS = InventoryOptionalField._(4,
          _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_REPLICATIONSTATUS');
  static const InventoryOptionalField INVENTORY_OPTIONAL_FIELD_SIZE =
      InventoryOptionalField._(
          5, _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_SIZE');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_OBJECTLOCKLEGALHOLDSTATUS =
      InventoryOptionalField._(
          6,
          _omitEnumNames
              ? ''
              : 'INVENTORY_OPTIONAL_FIELD_OBJECTLOCKLEGALHOLDSTATUS');
  static const InventoryOptionalField INVENTORY_OPTIONAL_FIELD_OBJECTOWNER =
      InventoryOptionalField._(
          7, _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_OBJECTOWNER');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_OBJECTACCESSCONTROLLIST =
      InventoryOptionalField._(
          8,
          _omitEnumNames
              ? ''
              : 'INVENTORY_OPTIONAL_FIELD_OBJECTACCESSCONTROLLIST');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_ENCRYPTIONSTATUS = InventoryOptionalField._(
          9, _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_ENCRYPTIONSTATUS');
  static const InventoryOptionalField INVENTORY_OPTIONAL_FIELD_OBJECTLOCKMODE =
      InventoryOptionalField._(
          10, _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_OBJECTLOCKMODE');
  static const InventoryOptionalField INVENTORY_OPTIONAL_FIELD_ETAG =
      InventoryOptionalField._(
          11, _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_ETAG');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_ISMULTIPARTUPLOADED = InventoryOptionalField._(
          12,
          _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_ISMULTIPARTUPLOADED');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_INTELLIGENTTIERINGACCESSTIER =
      InventoryOptionalField._(
          13,
          _omitEnumNames
              ? ''
              : 'INVENTORY_OPTIONAL_FIELD_INTELLIGENTTIERINGACCESSTIER');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_OBJECTLOCKRETAINUNTILDATE =
      InventoryOptionalField._(
          14,
          _omitEnumNames
              ? ''
              : 'INVENTORY_OPTIONAL_FIELD_OBJECTLOCKRETAINUNTILDATE');
  static const InventoryOptionalField
      INVENTORY_OPTIONAL_FIELD_LASTMODIFIEDDATE = InventoryOptionalField._(15,
          _omitEnumNames ? '' : 'INVENTORY_OPTIONAL_FIELD_LASTMODIFIEDDATE');

  static const $core.List<InventoryOptionalField> values =
      <InventoryOptionalField>[
    INVENTORY_OPTIONAL_FIELD_LIFECYCLEEXPIRATIONDATE,
    INVENTORY_OPTIONAL_FIELD_BUCKETKEYSTATUS,
    INVENTORY_OPTIONAL_FIELD_STORAGECLASS,
    INVENTORY_OPTIONAL_FIELD_CHECKSUMALGORITHM,
    INVENTORY_OPTIONAL_FIELD_REPLICATIONSTATUS,
    INVENTORY_OPTIONAL_FIELD_SIZE,
    INVENTORY_OPTIONAL_FIELD_OBJECTLOCKLEGALHOLDSTATUS,
    INVENTORY_OPTIONAL_FIELD_OBJECTOWNER,
    INVENTORY_OPTIONAL_FIELD_OBJECTACCESSCONTROLLIST,
    INVENTORY_OPTIONAL_FIELD_ENCRYPTIONSTATUS,
    INVENTORY_OPTIONAL_FIELD_OBJECTLOCKMODE,
    INVENTORY_OPTIONAL_FIELD_ETAG,
    INVENTORY_OPTIONAL_FIELD_ISMULTIPARTUPLOADED,
    INVENTORY_OPTIONAL_FIELD_INTELLIGENTTIERINGACCESSTIER,
    INVENTORY_OPTIONAL_FIELD_OBJECTLOCKRETAINUNTILDATE,
    INVENTORY_OPTIONAL_FIELD_LASTMODIFIEDDATE,
  ];

  static final $core.List<InventoryOptionalField?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 15);
  static InventoryOptionalField? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventoryOptionalField._(super.value, super.name);
}

class JSONType extends $pb.ProtobufEnum {
  static const JSONType J_S_O_N_TYPE_LINES =
      JSONType._(0, _omitEnumNames ? '' : 'J_S_O_N_TYPE_LINES');
  static const JSONType J_S_O_N_TYPE_DOCUMENT =
      JSONType._(1, _omitEnumNames ? '' : 'J_S_O_N_TYPE_DOCUMENT');

  static const $core.List<JSONType> values = <JSONType>[
    J_S_O_N_TYPE_LINES,
    J_S_O_N_TYPE_DOCUMENT,
  ];

  static final $core.List<JSONType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static JSONType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const JSONType._(super.value, super.name);
}

class LocationType extends $pb.ProtobufEnum {
  static const LocationType LOCATION_TYPE_AVAILABILITYZONE =
      LocationType._(0, _omitEnumNames ? '' : 'LOCATION_TYPE_AVAILABILITYZONE');
  static const LocationType LOCATION_TYPE_LOCALZONE =
      LocationType._(1, _omitEnumNames ? '' : 'LOCATION_TYPE_LOCALZONE');

  static const $core.List<LocationType> values = <LocationType>[
    LOCATION_TYPE_AVAILABILITYZONE,
    LOCATION_TYPE_LOCALZONE,
  ];

  static final $core.List<LocationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static LocationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LocationType._(super.value, super.name);
}

class MFADelete extends $pb.ProtobufEnum {
  static const MFADelete M_F_A_DELETE_DISABLED =
      MFADelete._(0, _omitEnumNames ? '' : 'M_F_A_DELETE_DISABLED');
  static const MFADelete M_F_A_DELETE_ENABLED =
      MFADelete._(1, _omitEnumNames ? '' : 'M_F_A_DELETE_ENABLED');

  static const $core.List<MFADelete> values = <MFADelete>[
    M_F_A_DELETE_DISABLED,
    M_F_A_DELETE_ENABLED,
  ];

  static final $core.List<MFADelete?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MFADelete? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MFADelete._(super.value, super.name);
}

class MFADeleteStatus extends $pb.ProtobufEnum {
  static const MFADeleteStatus M_F_A_DELETE_STATUS_DISABLED = MFADeleteStatus._(
      0, _omitEnumNames ? '' : 'M_F_A_DELETE_STATUS_DISABLED');
  static const MFADeleteStatus M_F_A_DELETE_STATUS_ENABLED =
      MFADeleteStatus._(1, _omitEnumNames ? '' : 'M_F_A_DELETE_STATUS_ENABLED');

  static const $core.List<MFADeleteStatus> values = <MFADeleteStatus>[
    M_F_A_DELETE_STATUS_DISABLED,
    M_F_A_DELETE_STATUS_ENABLED,
  ];

  static final $core.List<MFADeleteStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MFADeleteStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MFADeleteStatus._(super.value, super.name);
}

class MetadataDirective extends $pb.ProtobufEnum {
  static const MetadataDirective METADATA_DIRECTIVE_REPLACE =
      MetadataDirective._(
          0, _omitEnumNames ? '' : 'METADATA_DIRECTIVE_REPLACE');
  static const MetadataDirective METADATA_DIRECTIVE_COPY =
      MetadataDirective._(1, _omitEnumNames ? '' : 'METADATA_DIRECTIVE_COPY');

  static const $core.List<MetadataDirective> values = <MetadataDirective>[
    METADATA_DIRECTIVE_REPLACE,
    METADATA_DIRECTIVE_COPY,
  ];

  static final $core.List<MetadataDirective?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MetadataDirective? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetadataDirective._(super.value, super.name);
}

class MetricsStatus extends $pb.ProtobufEnum {
  static const MetricsStatus METRICS_STATUS_DISABLED =
      MetricsStatus._(0, _omitEnumNames ? '' : 'METRICS_STATUS_DISABLED');
  static const MetricsStatus METRICS_STATUS_ENABLED =
      MetricsStatus._(1, _omitEnumNames ? '' : 'METRICS_STATUS_ENABLED');

  static const $core.List<MetricsStatus> values = <MetricsStatus>[
    METRICS_STATUS_DISABLED,
    METRICS_STATUS_ENABLED,
  ];

  static final $core.List<MetricsStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MetricsStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetricsStatus._(super.value, super.name);
}

class ObjectAttributes extends $pb.ProtobufEnum {
  static const ObjectAttributes OBJECT_ATTRIBUTES_OBJECT_PARTS =
      ObjectAttributes._(
          0, _omitEnumNames ? '' : 'OBJECT_ATTRIBUTES_OBJECT_PARTS');
  static const ObjectAttributes OBJECT_ATTRIBUTES_CHECKSUM =
      ObjectAttributes._(1, _omitEnumNames ? '' : 'OBJECT_ATTRIBUTES_CHECKSUM');
  static const ObjectAttributes OBJECT_ATTRIBUTES_STORAGE_CLASS =
      ObjectAttributes._(
          2, _omitEnumNames ? '' : 'OBJECT_ATTRIBUTES_STORAGE_CLASS');
  static const ObjectAttributes OBJECT_ATTRIBUTES_OBJECT_SIZE =
      ObjectAttributes._(
          3, _omitEnumNames ? '' : 'OBJECT_ATTRIBUTES_OBJECT_SIZE');
  static const ObjectAttributes OBJECT_ATTRIBUTES_ETAG =
      ObjectAttributes._(4, _omitEnumNames ? '' : 'OBJECT_ATTRIBUTES_ETAG');

  static const $core.List<ObjectAttributes> values = <ObjectAttributes>[
    OBJECT_ATTRIBUTES_OBJECT_PARTS,
    OBJECT_ATTRIBUTES_CHECKSUM,
    OBJECT_ATTRIBUTES_STORAGE_CLASS,
    OBJECT_ATTRIBUTES_OBJECT_SIZE,
    OBJECT_ATTRIBUTES_ETAG,
  ];

  static final $core.List<ObjectAttributes?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ObjectAttributes? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectAttributes._(super.value, super.name);
}

class ObjectCannedACL extends $pb.ProtobufEnum {
  static const ObjectCannedACL OBJECT_CANNED_A_C_L_AUTHENTICATED_READ =
      ObjectCannedACL._(
          0, _omitEnumNames ? '' : 'OBJECT_CANNED_A_C_L_AUTHENTICATED_READ');
  static const ObjectCannedACL OBJECT_CANNED_A_C_L_BUCKET_OWNER_READ =
      ObjectCannedACL._(
          1, _omitEnumNames ? '' : 'OBJECT_CANNED_A_C_L_BUCKET_OWNER_READ');
  static const ObjectCannedACL OBJECT_CANNED_A_C_L_PRIVATE =
      ObjectCannedACL._(2, _omitEnumNames ? '' : 'OBJECT_CANNED_A_C_L_PRIVATE');
  static const ObjectCannedACL OBJECT_CANNED_A_C_L_AWS_EXEC_READ =
      ObjectCannedACL._(
          3, _omitEnumNames ? '' : 'OBJECT_CANNED_A_C_L_AWS_EXEC_READ');
  static const ObjectCannedACL OBJECT_CANNED_A_C_L_BUCKET_OWNER_FULL_CONTROL =
      ObjectCannedACL._(
          4,
          _omitEnumNames
              ? ''
              : 'OBJECT_CANNED_A_C_L_BUCKET_OWNER_FULL_CONTROL');
  static const ObjectCannedACL OBJECT_CANNED_A_C_L_PUBLIC_READ =
      ObjectCannedACL._(
          5, _omitEnumNames ? '' : 'OBJECT_CANNED_A_C_L_PUBLIC_READ');
  static const ObjectCannedACL OBJECT_CANNED_A_C_L_PUBLIC_READ_WRITE =
      ObjectCannedACL._(
          6, _omitEnumNames ? '' : 'OBJECT_CANNED_A_C_L_PUBLIC_READ_WRITE');

  static const $core.List<ObjectCannedACL> values = <ObjectCannedACL>[
    OBJECT_CANNED_A_C_L_AUTHENTICATED_READ,
    OBJECT_CANNED_A_C_L_BUCKET_OWNER_READ,
    OBJECT_CANNED_A_C_L_PRIVATE,
    OBJECT_CANNED_A_C_L_AWS_EXEC_READ,
    OBJECT_CANNED_A_C_L_BUCKET_OWNER_FULL_CONTROL,
    OBJECT_CANNED_A_C_L_PUBLIC_READ,
    OBJECT_CANNED_A_C_L_PUBLIC_READ_WRITE,
  ];

  static final $core.List<ObjectCannedACL?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static ObjectCannedACL? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectCannedACL._(super.value, super.name);
}

class ObjectLockEnabled extends $pb.ProtobufEnum {
  static const ObjectLockEnabled OBJECT_LOCK_ENABLED_ENABLED =
      ObjectLockEnabled._(
          0, _omitEnumNames ? '' : 'OBJECT_LOCK_ENABLED_ENABLED');

  static const $core.List<ObjectLockEnabled> values = <ObjectLockEnabled>[
    OBJECT_LOCK_ENABLED_ENABLED,
  ];

  static final $core.List<ObjectLockEnabled?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ObjectLockEnabled? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectLockEnabled._(super.value, super.name);
}

class ObjectLockLegalHoldStatus extends $pb.ProtobufEnum {
  static const ObjectLockLegalHoldStatus OBJECT_LOCK_LEGAL_HOLD_STATUS_OFF =
      ObjectLockLegalHoldStatus._(
          0, _omitEnumNames ? '' : 'OBJECT_LOCK_LEGAL_HOLD_STATUS_OFF');
  static const ObjectLockLegalHoldStatus OBJECT_LOCK_LEGAL_HOLD_STATUS_ON =
      ObjectLockLegalHoldStatus._(
          1, _omitEnumNames ? '' : 'OBJECT_LOCK_LEGAL_HOLD_STATUS_ON');

  static const $core.List<ObjectLockLegalHoldStatus> values =
      <ObjectLockLegalHoldStatus>[
    OBJECT_LOCK_LEGAL_HOLD_STATUS_OFF,
    OBJECT_LOCK_LEGAL_HOLD_STATUS_ON,
  ];

  static final $core.List<ObjectLockLegalHoldStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ObjectLockLegalHoldStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectLockLegalHoldStatus._(super.value, super.name);
}

class ObjectLockMode extends $pb.ProtobufEnum {
  static const ObjectLockMode OBJECT_LOCK_MODE_GOVERNANCE =
      ObjectLockMode._(0, _omitEnumNames ? '' : 'OBJECT_LOCK_MODE_GOVERNANCE');
  static const ObjectLockMode OBJECT_LOCK_MODE_COMPLIANCE =
      ObjectLockMode._(1, _omitEnumNames ? '' : 'OBJECT_LOCK_MODE_COMPLIANCE');

  static const $core.List<ObjectLockMode> values = <ObjectLockMode>[
    OBJECT_LOCK_MODE_GOVERNANCE,
    OBJECT_LOCK_MODE_COMPLIANCE,
  ];

  static final $core.List<ObjectLockMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ObjectLockMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectLockMode._(super.value, super.name);
}

class ObjectLockRetentionMode extends $pb.ProtobufEnum {
  static const ObjectLockRetentionMode OBJECT_LOCK_RETENTION_MODE_GOVERNANCE =
      ObjectLockRetentionMode._(
          0, _omitEnumNames ? '' : 'OBJECT_LOCK_RETENTION_MODE_GOVERNANCE');
  static const ObjectLockRetentionMode OBJECT_LOCK_RETENTION_MODE_COMPLIANCE =
      ObjectLockRetentionMode._(
          1, _omitEnumNames ? '' : 'OBJECT_LOCK_RETENTION_MODE_COMPLIANCE');

  static const $core.List<ObjectLockRetentionMode> values =
      <ObjectLockRetentionMode>[
    OBJECT_LOCK_RETENTION_MODE_GOVERNANCE,
    OBJECT_LOCK_RETENTION_MODE_COMPLIANCE,
  ];

  static final $core.List<ObjectLockRetentionMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ObjectLockRetentionMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectLockRetentionMode._(super.value, super.name);
}

class ObjectOwnership extends $pb.ProtobufEnum {
  static const ObjectOwnership OBJECT_OWNERSHIP_BUCKETOWNERPREFERRED =
      ObjectOwnership._(
          0, _omitEnumNames ? '' : 'OBJECT_OWNERSHIP_BUCKETOWNERPREFERRED');
  static const ObjectOwnership OBJECT_OWNERSHIP_OBJECTWRITER =
      ObjectOwnership._(
          1, _omitEnumNames ? '' : 'OBJECT_OWNERSHIP_OBJECTWRITER');
  static const ObjectOwnership OBJECT_OWNERSHIP_BUCKETOWNERENFORCED =
      ObjectOwnership._(
          2, _omitEnumNames ? '' : 'OBJECT_OWNERSHIP_BUCKETOWNERENFORCED');

  static const $core.List<ObjectOwnership> values = <ObjectOwnership>[
    OBJECT_OWNERSHIP_BUCKETOWNERPREFERRED,
    OBJECT_OWNERSHIP_OBJECTWRITER,
    OBJECT_OWNERSHIP_BUCKETOWNERENFORCED,
  ];

  static final $core.List<ObjectOwnership?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ObjectOwnership? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectOwnership._(super.value, super.name);
}

class ObjectStorageClass extends $pb.ProtobufEnum {
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_OUTPOSTS =
      ObjectStorageClass._(
          0, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_OUTPOSTS');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_EXPRESS_ONEZONE =
      ObjectStorageClass._(
          1, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_EXPRESS_ONEZONE');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_REDUCED_REDUNDANCY =
      ObjectStorageClass._(
          2, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_REDUCED_REDUNDANCY');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_DEEP_ARCHIVE =
      ObjectStorageClass._(
          3, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_DEEP_ARCHIVE');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_SNOW =
      ObjectStorageClass._(
          4, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_SNOW');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_ONEZONE_IA =
      ObjectStorageClass._(
          5, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_ONEZONE_IA');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_STANDARD_IA =
      ObjectStorageClass._(
          6, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_STANDARD_IA');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_GLACIER =
      ObjectStorageClass._(
          7, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_GLACIER');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_STANDARD =
      ObjectStorageClass._(
          8, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_STANDARD');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_FSX_ONTAP =
      ObjectStorageClass._(
          9, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_FSX_ONTAP');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_GLACIER_IR =
      ObjectStorageClass._(
          10, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_GLACIER_IR');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_FSX_OPENZFS =
      ObjectStorageClass._(
          11, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_FSX_OPENZFS');
  static const ObjectStorageClass OBJECT_STORAGE_CLASS_INTELLIGENT_TIERING =
      ObjectStorageClass._(
          12, _omitEnumNames ? '' : 'OBJECT_STORAGE_CLASS_INTELLIGENT_TIERING');

  static const $core.List<ObjectStorageClass> values = <ObjectStorageClass>[
    OBJECT_STORAGE_CLASS_OUTPOSTS,
    OBJECT_STORAGE_CLASS_EXPRESS_ONEZONE,
    OBJECT_STORAGE_CLASS_REDUCED_REDUNDANCY,
    OBJECT_STORAGE_CLASS_DEEP_ARCHIVE,
    OBJECT_STORAGE_CLASS_SNOW,
    OBJECT_STORAGE_CLASS_ONEZONE_IA,
    OBJECT_STORAGE_CLASS_STANDARD_IA,
    OBJECT_STORAGE_CLASS_GLACIER,
    OBJECT_STORAGE_CLASS_STANDARD,
    OBJECT_STORAGE_CLASS_FSX_ONTAP,
    OBJECT_STORAGE_CLASS_GLACIER_IR,
    OBJECT_STORAGE_CLASS_FSX_OPENZFS,
    OBJECT_STORAGE_CLASS_INTELLIGENT_TIERING,
  ];

  static final $core.List<ObjectStorageClass?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 12);
  static ObjectStorageClass? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectStorageClass._(super.value, super.name);
}

class ObjectVersionStorageClass extends $pb.ProtobufEnum {
  static const ObjectVersionStorageClass OBJECT_VERSION_STORAGE_CLASS_STANDARD =
      ObjectVersionStorageClass._(
          0, _omitEnumNames ? '' : 'OBJECT_VERSION_STORAGE_CLASS_STANDARD');

  static const $core.List<ObjectVersionStorageClass> values =
      <ObjectVersionStorageClass>[
    OBJECT_VERSION_STORAGE_CLASS_STANDARD,
  ];

  static final $core.List<ObjectVersionStorageClass?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ObjectVersionStorageClass? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ObjectVersionStorageClass._(super.value, super.name);
}

class OptionalObjectAttributes extends $pb.ProtobufEnum {
  static const OptionalObjectAttributes
      OPTIONAL_OBJECT_ATTRIBUTES_RESTORE_STATUS = OptionalObjectAttributes._(
          0, _omitEnumNames ? '' : 'OPTIONAL_OBJECT_ATTRIBUTES_RESTORE_STATUS');

  static const $core.List<OptionalObjectAttributes> values =
      <OptionalObjectAttributes>[
    OPTIONAL_OBJECT_ATTRIBUTES_RESTORE_STATUS,
  ];

  static final $core.List<OptionalObjectAttributes?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static OptionalObjectAttributes? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OptionalObjectAttributes._(super.value, super.name);
}

class OwnerOverride extends $pb.ProtobufEnum {
  static const OwnerOverride OWNER_OVERRIDE_DESTINATION =
      OwnerOverride._(0, _omitEnumNames ? '' : 'OWNER_OVERRIDE_DESTINATION');

  static const $core.List<OwnerOverride> values = <OwnerOverride>[
    OWNER_OVERRIDE_DESTINATION,
  ];

  static final $core.List<OwnerOverride?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static OwnerOverride? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OwnerOverride._(super.value, super.name);
}

class PartitionDateSource extends $pb.ProtobufEnum {
  static const PartitionDateSource PARTITION_DATE_SOURCE_DELIVERYTIME =
      PartitionDateSource._(
          0, _omitEnumNames ? '' : 'PARTITION_DATE_SOURCE_DELIVERYTIME');
  static const PartitionDateSource PARTITION_DATE_SOURCE_EVENTTIME =
      PartitionDateSource._(
          1, _omitEnumNames ? '' : 'PARTITION_DATE_SOURCE_EVENTTIME');

  static const $core.List<PartitionDateSource> values = <PartitionDateSource>[
    PARTITION_DATE_SOURCE_DELIVERYTIME,
    PARTITION_DATE_SOURCE_EVENTTIME,
  ];

  static final $core.List<PartitionDateSource?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PartitionDateSource? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PartitionDateSource._(super.value, super.name);
}

class Payer extends $pb.ProtobufEnum {
  static const Payer PAYER_REQUESTER =
      Payer._(0, _omitEnumNames ? '' : 'PAYER_REQUESTER');
  static const Payer PAYER_BUCKETOWNER =
      Payer._(1, _omitEnumNames ? '' : 'PAYER_BUCKETOWNER');

  static const $core.List<Payer> values = <Payer>[
    PAYER_REQUESTER,
    PAYER_BUCKETOWNER,
  ];

  static final $core.List<Payer?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static Payer? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Payer._(super.value, super.name);
}

class Permission extends $pb.ProtobufEnum {
  static const Permission PERMISSION_READ =
      Permission._(0, _omitEnumNames ? '' : 'PERMISSION_READ');
  static const Permission PERMISSION_WRITE_ACP =
      Permission._(1, _omitEnumNames ? '' : 'PERMISSION_WRITE_ACP');
  static const Permission PERMISSION_WRITE =
      Permission._(2, _omitEnumNames ? '' : 'PERMISSION_WRITE');
  static const Permission PERMISSION_FULL_CONTROL =
      Permission._(3, _omitEnumNames ? '' : 'PERMISSION_FULL_CONTROL');
  static const Permission PERMISSION_READ_ACP =
      Permission._(4, _omitEnumNames ? '' : 'PERMISSION_READ_ACP');

  static const $core.List<Permission> values = <Permission>[
    PERMISSION_READ,
    PERMISSION_WRITE_ACP,
    PERMISSION_WRITE,
    PERMISSION_FULL_CONTROL,
    PERMISSION_READ_ACP,
  ];

  static final $core.List<Permission?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static Permission? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Permission._(super.value, super.name);
}

class Protocol extends $pb.ProtobufEnum {
  static const Protocol PROTOCOL_HTTPS =
      Protocol._(0, _omitEnumNames ? '' : 'PROTOCOL_HTTPS');
  static const Protocol PROTOCOL_HTTP =
      Protocol._(1, _omitEnumNames ? '' : 'PROTOCOL_HTTP');

  static const $core.List<Protocol> values = <Protocol>[
    PROTOCOL_HTTPS,
    PROTOCOL_HTTP,
  ];

  static final $core.List<Protocol?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static Protocol? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Protocol._(super.value, super.name);
}

class QuoteFields extends $pb.ProtobufEnum {
  static const QuoteFields QUOTE_FIELDS_ALWAYS =
      QuoteFields._(0, _omitEnumNames ? '' : 'QUOTE_FIELDS_ALWAYS');
  static const QuoteFields QUOTE_FIELDS_ASNEEDED =
      QuoteFields._(1, _omitEnumNames ? '' : 'QUOTE_FIELDS_ASNEEDED');

  static const $core.List<QuoteFields> values = <QuoteFields>[
    QUOTE_FIELDS_ALWAYS,
    QUOTE_FIELDS_ASNEEDED,
  ];

  static final $core.List<QuoteFields?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static QuoteFields? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QuoteFields._(super.value, super.name);
}

class ReplicaModificationsStatus extends $pb.ProtobufEnum {
  static const ReplicaModificationsStatus
      REPLICA_MODIFICATIONS_STATUS_DISABLED = ReplicaModificationsStatus._(
          0, _omitEnumNames ? '' : 'REPLICA_MODIFICATIONS_STATUS_DISABLED');
  static const ReplicaModificationsStatus REPLICA_MODIFICATIONS_STATUS_ENABLED =
      ReplicaModificationsStatus._(
          1, _omitEnumNames ? '' : 'REPLICA_MODIFICATIONS_STATUS_ENABLED');

  static const $core.List<ReplicaModificationsStatus> values =
      <ReplicaModificationsStatus>[
    REPLICA_MODIFICATIONS_STATUS_DISABLED,
    REPLICA_MODIFICATIONS_STATUS_ENABLED,
  ];

  static final $core.List<ReplicaModificationsStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ReplicaModificationsStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReplicaModificationsStatus._(super.value, super.name);
}

class ReplicationRuleStatus extends $pb.ProtobufEnum {
  static const ReplicationRuleStatus REPLICATION_RULE_STATUS_DISABLED =
      ReplicationRuleStatus._(
          0, _omitEnumNames ? '' : 'REPLICATION_RULE_STATUS_DISABLED');
  static const ReplicationRuleStatus REPLICATION_RULE_STATUS_ENABLED =
      ReplicationRuleStatus._(
          1, _omitEnumNames ? '' : 'REPLICATION_RULE_STATUS_ENABLED');

  static const $core.List<ReplicationRuleStatus> values =
      <ReplicationRuleStatus>[
    REPLICATION_RULE_STATUS_DISABLED,
    REPLICATION_RULE_STATUS_ENABLED,
  ];

  static final $core.List<ReplicationRuleStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ReplicationRuleStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReplicationRuleStatus._(super.value, super.name);
}

class ReplicationStatus extends $pb.ProtobufEnum {
  static const ReplicationStatus REPLICATION_STATUS_PENDING =
      ReplicationStatus._(
          0, _omitEnumNames ? '' : 'REPLICATION_STATUS_PENDING');
  static const ReplicationStatus REPLICATION_STATUS_REPLICA =
      ReplicationStatus._(
          1, _omitEnumNames ? '' : 'REPLICATION_STATUS_REPLICA');
  static const ReplicationStatus REPLICATION_STATUS_COMPLETE =
      ReplicationStatus._(
          2, _omitEnumNames ? '' : 'REPLICATION_STATUS_COMPLETE');
  static const ReplicationStatus REPLICATION_STATUS_COMPLETED =
      ReplicationStatus._(
          3, _omitEnumNames ? '' : 'REPLICATION_STATUS_COMPLETED');
  static const ReplicationStatus REPLICATION_STATUS_FAILED =
      ReplicationStatus._(4, _omitEnumNames ? '' : 'REPLICATION_STATUS_FAILED');

  static const $core.List<ReplicationStatus> values = <ReplicationStatus>[
    REPLICATION_STATUS_PENDING,
    REPLICATION_STATUS_REPLICA,
    REPLICATION_STATUS_COMPLETE,
    REPLICATION_STATUS_COMPLETED,
    REPLICATION_STATUS_FAILED,
  ];

  static final $core.List<ReplicationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ReplicationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReplicationStatus._(super.value, super.name);
}

class ReplicationTimeStatus extends $pb.ProtobufEnum {
  static const ReplicationTimeStatus REPLICATION_TIME_STATUS_DISABLED =
      ReplicationTimeStatus._(
          0, _omitEnumNames ? '' : 'REPLICATION_TIME_STATUS_DISABLED');
  static const ReplicationTimeStatus REPLICATION_TIME_STATUS_ENABLED =
      ReplicationTimeStatus._(
          1, _omitEnumNames ? '' : 'REPLICATION_TIME_STATUS_ENABLED');

  static const $core.List<ReplicationTimeStatus> values =
      <ReplicationTimeStatus>[
    REPLICATION_TIME_STATUS_DISABLED,
    REPLICATION_TIME_STATUS_ENABLED,
  ];

  static final $core.List<ReplicationTimeStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ReplicationTimeStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReplicationTimeStatus._(super.value, super.name);
}

class RequestCharged extends $pb.ProtobufEnum {
  static const RequestCharged REQUEST_CHARGED_REQUESTER =
      RequestCharged._(0, _omitEnumNames ? '' : 'REQUEST_CHARGED_REQUESTER');

  static const $core.List<RequestCharged> values = <RequestCharged>[
    REQUEST_CHARGED_REQUESTER,
  ];

  static final $core.List<RequestCharged?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static RequestCharged? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RequestCharged._(super.value, super.name);
}

class RequestPayer extends $pb.ProtobufEnum {
  static const RequestPayer REQUEST_PAYER_REQUESTER =
      RequestPayer._(0, _omitEnumNames ? '' : 'REQUEST_PAYER_REQUESTER');

  static const $core.List<RequestPayer> values = <RequestPayer>[
    REQUEST_PAYER_REQUESTER,
  ];

  static final $core.List<RequestPayer?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static RequestPayer? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RequestPayer._(super.value, super.name);
}

class RestoreRequestType extends $pb.ProtobufEnum {
  static const RestoreRequestType RESTORE_REQUEST_TYPE_SELECT =
      RestoreRequestType._(
          0, _omitEnumNames ? '' : 'RESTORE_REQUEST_TYPE_SELECT');

  static const $core.List<RestoreRequestType> values = <RestoreRequestType>[
    RESTORE_REQUEST_TYPE_SELECT,
  ];

  static final $core.List<RestoreRequestType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static RestoreRequestType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RestoreRequestType._(super.value, super.name);
}

class S3TablesBucketType extends $pb.ProtobufEnum {
  static const S3TablesBucketType S3_TABLES_BUCKET_TYPE_CUSTOMER =
      S3TablesBucketType._(
          0, _omitEnumNames ? '' : 'S3_TABLES_BUCKET_TYPE_CUSTOMER');
  static const S3TablesBucketType S3_TABLES_BUCKET_TYPE_AWS =
      S3TablesBucketType._(
          1, _omitEnumNames ? '' : 'S3_TABLES_BUCKET_TYPE_AWS');

  static const $core.List<S3TablesBucketType> values = <S3TablesBucketType>[
    S3_TABLES_BUCKET_TYPE_CUSTOMER,
    S3_TABLES_BUCKET_TYPE_AWS,
  ];

  static final $core.List<S3TablesBucketType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static S3TablesBucketType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const S3TablesBucketType._(super.value, super.name);
}

class ServerSideEncryption extends $pb.ProtobufEnum {
  static const ServerSideEncryption SERVER_SIDE_ENCRYPTION_AWS_KMS =
      ServerSideEncryption._(
          0, _omitEnumNames ? '' : 'SERVER_SIDE_ENCRYPTION_AWS_KMS');
  static const ServerSideEncryption SERVER_SIDE_ENCRYPTION_AWS_KMS_DSSE =
      ServerSideEncryption._(
          1, _omitEnumNames ? '' : 'SERVER_SIDE_ENCRYPTION_AWS_KMS_DSSE');
  static const ServerSideEncryption SERVER_SIDE_ENCRYPTION_AES256 =
      ServerSideEncryption._(
          2, _omitEnumNames ? '' : 'SERVER_SIDE_ENCRYPTION_AES256');
  static const ServerSideEncryption SERVER_SIDE_ENCRYPTION_AWS_FSX =
      ServerSideEncryption._(
          3, _omitEnumNames ? '' : 'SERVER_SIDE_ENCRYPTION_AWS_FSX');

  static const $core.List<ServerSideEncryption> values = <ServerSideEncryption>[
    SERVER_SIDE_ENCRYPTION_AWS_KMS,
    SERVER_SIDE_ENCRYPTION_AWS_KMS_DSSE,
    SERVER_SIDE_ENCRYPTION_AES256,
    SERVER_SIDE_ENCRYPTION_AWS_FSX,
  ];

  static final $core.List<ServerSideEncryption?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ServerSideEncryption? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ServerSideEncryption._(super.value, super.name);
}

class SessionMode extends $pb.ProtobufEnum {
  static const SessionMode SESSION_MODE_READONLY =
      SessionMode._(0, _omitEnumNames ? '' : 'SESSION_MODE_READONLY');
  static const SessionMode SESSION_MODE_READWRITE =
      SessionMode._(1, _omitEnumNames ? '' : 'SESSION_MODE_READWRITE');

  static const $core.List<SessionMode> values = <SessionMode>[
    SESSION_MODE_READONLY,
    SESSION_MODE_READWRITE,
  ];

  static final $core.List<SessionMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SessionMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SessionMode._(super.value, super.name);
}

class SseKmsEncryptedObjectsStatus extends $pb.ProtobufEnum {
  static const SseKmsEncryptedObjectsStatus
      SSE_KMS_ENCRYPTED_OBJECTS_STATUS_DISABLED =
      SseKmsEncryptedObjectsStatus._(
          0, _omitEnumNames ? '' : 'SSE_KMS_ENCRYPTED_OBJECTS_STATUS_DISABLED');
  static const SseKmsEncryptedObjectsStatus
      SSE_KMS_ENCRYPTED_OBJECTS_STATUS_ENABLED = SseKmsEncryptedObjectsStatus._(
          1, _omitEnumNames ? '' : 'SSE_KMS_ENCRYPTED_OBJECTS_STATUS_ENABLED');

  static const $core.List<SseKmsEncryptedObjectsStatus> values =
      <SseKmsEncryptedObjectsStatus>[
    SSE_KMS_ENCRYPTED_OBJECTS_STATUS_DISABLED,
    SSE_KMS_ENCRYPTED_OBJECTS_STATUS_ENABLED,
  ];

  static final $core.List<SseKmsEncryptedObjectsStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SseKmsEncryptedObjectsStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SseKmsEncryptedObjectsStatus._(super.value, super.name);
}

class StorageClass extends $pb.ProtobufEnum {
  static const StorageClass STORAGE_CLASS_OUTPOSTS =
      StorageClass._(0, _omitEnumNames ? '' : 'STORAGE_CLASS_OUTPOSTS');
  static const StorageClass STORAGE_CLASS_EXPRESS_ONEZONE =
      StorageClass._(1, _omitEnumNames ? '' : 'STORAGE_CLASS_EXPRESS_ONEZONE');
  static const StorageClass STORAGE_CLASS_REDUCED_REDUNDANCY = StorageClass._(
      2, _omitEnumNames ? '' : 'STORAGE_CLASS_REDUCED_REDUNDANCY');
  static const StorageClass STORAGE_CLASS_DEEP_ARCHIVE =
      StorageClass._(3, _omitEnumNames ? '' : 'STORAGE_CLASS_DEEP_ARCHIVE');
  static const StorageClass STORAGE_CLASS_SNOW =
      StorageClass._(4, _omitEnumNames ? '' : 'STORAGE_CLASS_SNOW');
  static const StorageClass STORAGE_CLASS_ONEZONE_IA =
      StorageClass._(5, _omitEnumNames ? '' : 'STORAGE_CLASS_ONEZONE_IA');
  static const StorageClass STORAGE_CLASS_STANDARD_IA =
      StorageClass._(6, _omitEnumNames ? '' : 'STORAGE_CLASS_STANDARD_IA');
  static const StorageClass STORAGE_CLASS_GLACIER =
      StorageClass._(7, _omitEnumNames ? '' : 'STORAGE_CLASS_GLACIER');
  static const StorageClass STORAGE_CLASS_STANDARD =
      StorageClass._(8, _omitEnumNames ? '' : 'STORAGE_CLASS_STANDARD');
  static const StorageClass STORAGE_CLASS_FSX_ONTAP =
      StorageClass._(9, _omitEnumNames ? '' : 'STORAGE_CLASS_FSX_ONTAP');
  static const StorageClass STORAGE_CLASS_GLACIER_IR =
      StorageClass._(10, _omitEnumNames ? '' : 'STORAGE_CLASS_GLACIER_IR');
  static const StorageClass STORAGE_CLASS_FSX_OPENZFS =
      StorageClass._(11, _omitEnumNames ? '' : 'STORAGE_CLASS_FSX_OPENZFS');
  static const StorageClass STORAGE_CLASS_INTELLIGENT_TIERING = StorageClass._(
      12, _omitEnumNames ? '' : 'STORAGE_CLASS_INTELLIGENT_TIERING');

  static const $core.List<StorageClass> values = <StorageClass>[
    STORAGE_CLASS_OUTPOSTS,
    STORAGE_CLASS_EXPRESS_ONEZONE,
    STORAGE_CLASS_REDUCED_REDUNDANCY,
    STORAGE_CLASS_DEEP_ARCHIVE,
    STORAGE_CLASS_SNOW,
    STORAGE_CLASS_ONEZONE_IA,
    STORAGE_CLASS_STANDARD_IA,
    STORAGE_CLASS_GLACIER,
    STORAGE_CLASS_STANDARD,
    STORAGE_CLASS_FSX_ONTAP,
    STORAGE_CLASS_GLACIER_IR,
    STORAGE_CLASS_FSX_OPENZFS,
    STORAGE_CLASS_INTELLIGENT_TIERING,
  ];

  static final $core.List<StorageClass?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 12);
  static StorageClass? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StorageClass._(super.value, super.name);
}

class StorageClassAnalysisSchemaVersion extends $pb.ProtobufEnum {
  static const StorageClassAnalysisSchemaVersion
      STORAGE_CLASS_ANALYSIS_SCHEMA_VERSION_V_1 =
      StorageClassAnalysisSchemaVersion._(
          0, _omitEnumNames ? '' : 'STORAGE_CLASS_ANALYSIS_SCHEMA_VERSION_V_1');

  static const $core.List<StorageClassAnalysisSchemaVersion> values =
      <StorageClassAnalysisSchemaVersion>[
    STORAGE_CLASS_ANALYSIS_SCHEMA_VERSION_V_1,
  ];

  static final $core.List<StorageClassAnalysisSchemaVersion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static StorageClassAnalysisSchemaVersion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StorageClassAnalysisSchemaVersion._(super.value, super.name);
}

class TableSseAlgorithm extends $pb.ProtobufEnum {
  static const TableSseAlgorithm TABLE_SSE_ALGORITHM_AWS_KMS =
      TableSseAlgorithm._(
          0, _omitEnumNames ? '' : 'TABLE_SSE_ALGORITHM_AWS_KMS');
  static const TableSseAlgorithm TABLE_SSE_ALGORITHM_AES256 =
      TableSseAlgorithm._(
          1, _omitEnumNames ? '' : 'TABLE_SSE_ALGORITHM_AES256');

  static const $core.List<TableSseAlgorithm> values = <TableSseAlgorithm>[
    TABLE_SSE_ALGORITHM_AWS_KMS,
    TABLE_SSE_ALGORITHM_AES256,
  ];

  static final $core.List<TableSseAlgorithm?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static TableSseAlgorithm? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TableSseAlgorithm._(super.value, super.name);
}

class TaggingDirective extends $pb.ProtobufEnum {
  static const TaggingDirective TAGGING_DIRECTIVE_REPLACE =
      TaggingDirective._(0, _omitEnumNames ? '' : 'TAGGING_DIRECTIVE_REPLACE');
  static const TaggingDirective TAGGING_DIRECTIVE_COPY =
      TaggingDirective._(1, _omitEnumNames ? '' : 'TAGGING_DIRECTIVE_COPY');

  static const $core.List<TaggingDirective> values = <TaggingDirective>[
    TAGGING_DIRECTIVE_REPLACE,
    TAGGING_DIRECTIVE_COPY,
  ];

  static final $core.List<TaggingDirective?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static TaggingDirective? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TaggingDirective._(super.value, super.name);
}

class Tier extends $pb.ProtobufEnum {
  static const Tier TIER_BULK = Tier._(0, _omitEnumNames ? '' : 'TIER_BULK');
  static const Tier TIER_EXPEDITED =
      Tier._(1, _omitEnumNames ? '' : 'TIER_EXPEDITED');
  static const Tier TIER_STANDARD =
      Tier._(2, _omitEnumNames ? '' : 'TIER_STANDARD');

  static const $core.List<Tier> values = <Tier>[
    TIER_BULK,
    TIER_EXPEDITED,
    TIER_STANDARD,
  ];

  static final $core.List<Tier?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static Tier? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Tier._(super.value, super.name);
}

class TransitionDefaultMinimumObjectSize extends $pb.ProtobufEnum {
  static const TransitionDefaultMinimumObjectSize
      TRANSITION_DEFAULT_MINIMUM_OBJECT_SIZE_ALL_STORAGE_CLASSES_128K =
      TransitionDefaultMinimumObjectSize._(
          0,
          _omitEnumNames
              ? ''
              : 'TRANSITION_DEFAULT_MINIMUM_OBJECT_SIZE_ALL_STORAGE_CLASSES_128K');
  static const TransitionDefaultMinimumObjectSize
      TRANSITION_DEFAULT_MINIMUM_OBJECT_SIZE_VARIES_BY_STORAGE_CLASS =
      TransitionDefaultMinimumObjectSize._(
          1,
          _omitEnumNames
              ? ''
              : 'TRANSITION_DEFAULT_MINIMUM_OBJECT_SIZE_VARIES_BY_STORAGE_CLASS');

  static const $core.List<TransitionDefaultMinimumObjectSize> values =
      <TransitionDefaultMinimumObjectSize>[
    TRANSITION_DEFAULT_MINIMUM_OBJECT_SIZE_ALL_STORAGE_CLASSES_128K,
    TRANSITION_DEFAULT_MINIMUM_OBJECT_SIZE_VARIES_BY_STORAGE_CLASS,
  ];

  static final $core.List<TransitionDefaultMinimumObjectSize?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static TransitionDefaultMinimumObjectSize? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TransitionDefaultMinimumObjectSize._(super.value, super.name);
}

class TransitionStorageClass extends $pb.ProtobufEnum {
  static const TransitionStorageClass TRANSITION_STORAGE_CLASS_DEEP_ARCHIVE =
      TransitionStorageClass._(
          0, _omitEnumNames ? '' : 'TRANSITION_STORAGE_CLASS_DEEP_ARCHIVE');
  static const TransitionStorageClass TRANSITION_STORAGE_CLASS_ONEZONE_IA =
      TransitionStorageClass._(
          1, _omitEnumNames ? '' : 'TRANSITION_STORAGE_CLASS_ONEZONE_IA');
  static const TransitionStorageClass TRANSITION_STORAGE_CLASS_STANDARD_IA =
      TransitionStorageClass._(
          2, _omitEnumNames ? '' : 'TRANSITION_STORAGE_CLASS_STANDARD_IA');
  static const TransitionStorageClass TRANSITION_STORAGE_CLASS_GLACIER =
      TransitionStorageClass._(
          3, _omitEnumNames ? '' : 'TRANSITION_STORAGE_CLASS_GLACIER');
  static const TransitionStorageClass TRANSITION_STORAGE_CLASS_GLACIER_IR =
      TransitionStorageClass._(
          4, _omitEnumNames ? '' : 'TRANSITION_STORAGE_CLASS_GLACIER_IR');
  static const TransitionStorageClass
      TRANSITION_STORAGE_CLASS_INTELLIGENT_TIERING = TransitionStorageClass._(5,
          _omitEnumNames ? '' : 'TRANSITION_STORAGE_CLASS_INTELLIGENT_TIERING');

  static const $core.List<TransitionStorageClass> values =
      <TransitionStorageClass>[
    TRANSITION_STORAGE_CLASS_DEEP_ARCHIVE,
    TRANSITION_STORAGE_CLASS_ONEZONE_IA,
    TRANSITION_STORAGE_CLASS_STANDARD_IA,
    TRANSITION_STORAGE_CLASS_GLACIER,
    TRANSITION_STORAGE_CLASS_GLACIER_IR,
    TRANSITION_STORAGE_CLASS_INTELLIGENT_TIERING,
  ];

  static final $core.List<TransitionStorageClass?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static TransitionStorageClass? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TransitionStorageClass._(super.value, super.name);
}

class Type extends $pb.ProtobufEnum {
  static const Type TYPE_GROUP = Type._(0, _omitEnumNames ? '' : 'TYPE_GROUP');
  static const Type TYPE_AMAZONCUSTOMERBYEMAIL =
      Type._(1, _omitEnumNames ? '' : 'TYPE_AMAZONCUSTOMERBYEMAIL');
  static const Type TYPE_CANONICALUSER =
      Type._(2, _omitEnumNames ? '' : 'TYPE_CANONICALUSER');

  static const $core.List<Type> values = <Type>[
    TYPE_GROUP,
    TYPE_AMAZONCUSTOMERBYEMAIL,
    TYPE_CANONICALUSER,
  ];

  static final $core.List<Type?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static Type? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Type._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
