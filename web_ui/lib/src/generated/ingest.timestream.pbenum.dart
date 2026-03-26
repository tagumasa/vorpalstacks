// This is a generated file - do not edit.
//
// Generated from ingest.timestream.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class BatchLoadDataFormat extends $pb.ProtobufEnum {
  static const BatchLoadDataFormat BATCH_LOAD_DATA_FORMAT_CSV =
      BatchLoadDataFormat._(
          0, _omitEnumNames ? '' : 'BATCH_LOAD_DATA_FORMAT_CSV');

  static const $core.List<BatchLoadDataFormat> values = <BatchLoadDataFormat>[
    BATCH_LOAD_DATA_FORMAT_CSV,
  ];

  static final $core.List<BatchLoadDataFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static BatchLoadDataFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BatchLoadDataFormat._(super.value, super.name);
}

class BatchLoadStatus extends $pb.ProtobufEnum {
  static const BatchLoadStatus BATCH_LOAD_STATUS_PENDING_RESUME =
      BatchLoadStatus._(
          0, _omitEnumNames ? '' : 'BATCH_LOAD_STATUS_PENDING_RESUME');
  static const BatchLoadStatus BATCH_LOAD_STATUS_IN_PROGRESS =
      BatchLoadStatus._(
          1, _omitEnumNames ? '' : 'BATCH_LOAD_STATUS_IN_PROGRESS');
  static const BatchLoadStatus BATCH_LOAD_STATUS_SUCCEEDED =
      BatchLoadStatus._(2, _omitEnumNames ? '' : 'BATCH_LOAD_STATUS_SUCCEEDED');
  static const BatchLoadStatus BATCH_LOAD_STATUS_PROGRESS_STOPPED =
      BatchLoadStatus._(
          3, _omitEnumNames ? '' : 'BATCH_LOAD_STATUS_PROGRESS_STOPPED');
  static const BatchLoadStatus BATCH_LOAD_STATUS_CREATED =
      BatchLoadStatus._(4, _omitEnumNames ? '' : 'BATCH_LOAD_STATUS_CREATED');
  static const BatchLoadStatus BATCH_LOAD_STATUS_FAILED =
      BatchLoadStatus._(5, _omitEnumNames ? '' : 'BATCH_LOAD_STATUS_FAILED');

  static const $core.List<BatchLoadStatus> values = <BatchLoadStatus>[
    BATCH_LOAD_STATUS_PENDING_RESUME,
    BATCH_LOAD_STATUS_IN_PROGRESS,
    BATCH_LOAD_STATUS_SUCCEEDED,
    BATCH_LOAD_STATUS_PROGRESS_STOPPED,
    BATCH_LOAD_STATUS_CREATED,
    BATCH_LOAD_STATUS_FAILED,
  ];

  static final $core.List<BatchLoadStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static BatchLoadStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BatchLoadStatus._(super.value, super.name);
}

class DimensionValueType extends $pb.ProtobufEnum {
  static const DimensionValueType DIMENSION_VALUE_TYPE_VARCHAR =
      DimensionValueType._(
          0, _omitEnumNames ? '' : 'DIMENSION_VALUE_TYPE_VARCHAR');

  static const $core.List<DimensionValueType> values = <DimensionValueType>[
    DIMENSION_VALUE_TYPE_VARCHAR,
  ];

  static final $core.List<DimensionValueType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static DimensionValueType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DimensionValueType._(super.value, super.name);
}

class MeasureValueType extends $pb.ProtobufEnum {
  static const MeasureValueType MEASURE_VALUE_TYPE_BIGINT =
      MeasureValueType._(0, _omitEnumNames ? '' : 'MEASURE_VALUE_TYPE_BIGINT');
  static const MeasureValueType MEASURE_VALUE_TYPE_MULTI =
      MeasureValueType._(1, _omitEnumNames ? '' : 'MEASURE_VALUE_TYPE_MULTI');
  static const MeasureValueType MEASURE_VALUE_TYPE_VARCHAR =
      MeasureValueType._(2, _omitEnumNames ? '' : 'MEASURE_VALUE_TYPE_VARCHAR');
  static const MeasureValueType MEASURE_VALUE_TYPE_BOOLEAN =
      MeasureValueType._(3, _omitEnumNames ? '' : 'MEASURE_VALUE_TYPE_BOOLEAN');
  static const MeasureValueType MEASURE_VALUE_TYPE_DOUBLE =
      MeasureValueType._(4, _omitEnumNames ? '' : 'MEASURE_VALUE_TYPE_DOUBLE');
  static const MeasureValueType MEASURE_VALUE_TYPE_TIMESTAMP =
      MeasureValueType._(
          5, _omitEnumNames ? '' : 'MEASURE_VALUE_TYPE_TIMESTAMP');

  static const $core.List<MeasureValueType> values = <MeasureValueType>[
    MEASURE_VALUE_TYPE_BIGINT,
    MEASURE_VALUE_TYPE_MULTI,
    MEASURE_VALUE_TYPE_VARCHAR,
    MEASURE_VALUE_TYPE_BOOLEAN,
    MEASURE_VALUE_TYPE_DOUBLE,
    MEASURE_VALUE_TYPE_TIMESTAMP,
  ];

  static final $core.List<MeasureValueType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static MeasureValueType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MeasureValueType._(super.value, super.name);
}

class PartitionKeyEnforcementLevel extends $pb.ProtobufEnum {
  static const PartitionKeyEnforcementLevel
      PARTITION_KEY_ENFORCEMENT_LEVEL_OPTIONAL = PartitionKeyEnforcementLevel._(
          0, _omitEnumNames ? '' : 'PARTITION_KEY_ENFORCEMENT_LEVEL_OPTIONAL');
  static const PartitionKeyEnforcementLevel
      PARTITION_KEY_ENFORCEMENT_LEVEL_REQUIRED = PartitionKeyEnforcementLevel._(
          1, _omitEnumNames ? '' : 'PARTITION_KEY_ENFORCEMENT_LEVEL_REQUIRED');

  static const $core.List<PartitionKeyEnforcementLevel> values =
      <PartitionKeyEnforcementLevel>[
    PARTITION_KEY_ENFORCEMENT_LEVEL_OPTIONAL,
    PARTITION_KEY_ENFORCEMENT_LEVEL_REQUIRED,
  ];

  static final $core.List<PartitionKeyEnforcementLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PartitionKeyEnforcementLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PartitionKeyEnforcementLevel._(super.value, super.name);
}

class PartitionKeyType extends $pb.ProtobufEnum {
  static const PartitionKeyType PARTITION_KEY_TYPE_DIMENSION =
      PartitionKeyType._(
          0, _omitEnumNames ? '' : 'PARTITION_KEY_TYPE_DIMENSION');
  static const PartitionKeyType PARTITION_KEY_TYPE_MEASURE =
      PartitionKeyType._(1, _omitEnumNames ? '' : 'PARTITION_KEY_TYPE_MEASURE');

  static const $core.List<PartitionKeyType> values = <PartitionKeyType>[
    PARTITION_KEY_TYPE_DIMENSION,
    PARTITION_KEY_TYPE_MEASURE,
  ];

  static final $core.List<PartitionKeyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PartitionKeyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PartitionKeyType._(super.value, super.name);
}

class S3EncryptionOption extends $pb.ProtobufEnum {
  static const S3EncryptionOption S3_ENCRYPTION_OPTION_SSE_S3 =
      S3EncryptionOption._(
          0, _omitEnumNames ? '' : 'S3_ENCRYPTION_OPTION_SSE_S3');
  static const S3EncryptionOption S3_ENCRYPTION_OPTION_SSE_KMS =
      S3EncryptionOption._(
          1, _omitEnumNames ? '' : 'S3_ENCRYPTION_OPTION_SSE_KMS');

  static const $core.List<S3EncryptionOption> values = <S3EncryptionOption>[
    S3_ENCRYPTION_OPTION_SSE_S3,
    S3_ENCRYPTION_OPTION_SSE_KMS,
  ];

  static final $core.List<S3EncryptionOption?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static S3EncryptionOption? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const S3EncryptionOption._(super.value, super.name);
}

class ScalarMeasureValueType extends $pb.ProtobufEnum {
  static const ScalarMeasureValueType SCALAR_MEASURE_VALUE_TYPE_BIGINT =
      ScalarMeasureValueType._(
          0, _omitEnumNames ? '' : 'SCALAR_MEASURE_VALUE_TYPE_BIGINT');
  static const ScalarMeasureValueType SCALAR_MEASURE_VALUE_TYPE_VARCHAR =
      ScalarMeasureValueType._(
          1, _omitEnumNames ? '' : 'SCALAR_MEASURE_VALUE_TYPE_VARCHAR');
  static const ScalarMeasureValueType SCALAR_MEASURE_VALUE_TYPE_BOOLEAN =
      ScalarMeasureValueType._(
          2, _omitEnumNames ? '' : 'SCALAR_MEASURE_VALUE_TYPE_BOOLEAN');
  static const ScalarMeasureValueType SCALAR_MEASURE_VALUE_TYPE_DOUBLE =
      ScalarMeasureValueType._(
          3, _omitEnumNames ? '' : 'SCALAR_MEASURE_VALUE_TYPE_DOUBLE');
  static const ScalarMeasureValueType SCALAR_MEASURE_VALUE_TYPE_TIMESTAMP =
      ScalarMeasureValueType._(
          4, _omitEnumNames ? '' : 'SCALAR_MEASURE_VALUE_TYPE_TIMESTAMP');

  static const $core.List<ScalarMeasureValueType> values =
      <ScalarMeasureValueType>[
    SCALAR_MEASURE_VALUE_TYPE_BIGINT,
    SCALAR_MEASURE_VALUE_TYPE_VARCHAR,
    SCALAR_MEASURE_VALUE_TYPE_BOOLEAN,
    SCALAR_MEASURE_VALUE_TYPE_DOUBLE,
    SCALAR_MEASURE_VALUE_TYPE_TIMESTAMP,
  ];

  static final $core.List<ScalarMeasureValueType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ScalarMeasureValueType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScalarMeasureValueType._(super.value, super.name);
}

class TableStatus extends $pb.ProtobufEnum {
  static const TableStatus TABLE_STATUS_RESTORING =
      TableStatus._(0, _omitEnumNames ? '' : 'TABLE_STATUS_RESTORING');
  static const TableStatus TABLE_STATUS_ACTIVE =
      TableStatus._(1, _omitEnumNames ? '' : 'TABLE_STATUS_ACTIVE');
  static const TableStatus TABLE_STATUS_DELETING =
      TableStatus._(2, _omitEnumNames ? '' : 'TABLE_STATUS_DELETING');

  static const $core.List<TableStatus> values = <TableStatus>[
    TABLE_STATUS_RESTORING,
    TABLE_STATUS_ACTIVE,
    TABLE_STATUS_DELETING,
  ];

  static final $core.List<TableStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static TableStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TableStatus._(super.value, super.name);
}

class TimeUnit extends $pb.ProtobufEnum {
  static const TimeUnit TIME_UNIT_MILLISECONDS =
      TimeUnit._(0, _omitEnumNames ? '' : 'TIME_UNIT_MILLISECONDS');
  static const TimeUnit TIME_UNIT_SECONDS =
      TimeUnit._(1, _omitEnumNames ? '' : 'TIME_UNIT_SECONDS');
  static const TimeUnit TIME_UNIT_MICROSECONDS =
      TimeUnit._(2, _omitEnumNames ? '' : 'TIME_UNIT_MICROSECONDS');
  static const TimeUnit TIME_UNIT_NANOSECONDS =
      TimeUnit._(3, _omitEnumNames ? '' : 'TIME_UNIT_NANOSECONDS');

  static const $core.List<TimeUnit> values = <TimeUnit>[
    TIME_UNIT_MILLISECONDS,
    TIME_UNIT_SECONDS,
    TIME_UNIT_MICROSECONDS,
    TIME_UNIT_NANOSECONDS,
  ];

  static final $core.List<TimeUnit?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static TimeUnit? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TimeUnit._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
