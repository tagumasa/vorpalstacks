// This is a generated file - do not edit.
//
// Generated from query.timestream.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ComputeMode extends $pb.ProtobufEnum {
  static const ComputeMode COMPUTE_MODE_PROVISIONED =
      ComputeMode._(0, _omitEnumNames ? '' : 'COMPUTE_MODE_PROVISIONED');
  static const ComputeMode COMPUTE_MODE_ON_DEMAND =
      ComputeMode._(1, _omitEnumNames ? '' : 'COMPUTE_MODE_ON_DEMAND');

  static const $core.List<ComputeMode> values = <ComputeMode>[
    COMPUTE_MODE_PROVISIONED,
    COMPUTE_MODE_ON_DEMAND,
  ];

  static final $core.List<ComputeMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ComputeMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComputeMode._(super.value, super.name);
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

class LastUpdateStatus extends $pb.ProtobufEnum {
  static const LastUpdateStatus LAST_UPDATE_STATUS_PENDING =
      LastUpdateStatus._(0, _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_PENDING');
  static const LastUpdateStatus LAST_UPDATE_STATUS_SUCCEEDED =
      LastUpdateStatus._(
          1, _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_SUCCEEDED');
  static const LastUpdateStatus LAST_UPDATE_STATUS_FAILED =
      LastUpdateStatus._(2, _omitEnumNames ? '' : 'LAST_UPDATE_STATUS_FAILED');

  static const $core.List<LastUpdateStatus> values = <LastUpdateStatus>[
    LAST_UPDATE_STATUS_PENDING,
    LAST_UPDATE_STATUS_SUCCEEDED,
    LAST_UPDATE_STATUS_FAILED,
  ];

  static final $core.List<LastUpdateStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static LastUpdateStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LastUpdateStatus._(super.value, super.name);
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

  static const $core.List<MeasureValueType> values = <MeasureValueType>[
    MEASURE_VALUE_TYPE_BIGINT,
    MEASURE_VALUE_TYPE_MULTI,
    MEASURE_VALUE_TYPE_VARCHAR,
    MEASURE_VALUE_TYPE_BOOLEAN,
    MEASURE_VALUE_TYPE_DOUBLE,
  ];

  static final $core.List<MeasureValueType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static MeasureValueType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MeasureValueType._(super.value, super.name);
}

class QueryInsightsMode extends $pb.ProtobufEnum {
  static const QueryInsightsMode QUERY_INSIGHTS_MODE_DISABLED =
      QueryInsightsMode._(
          0, _omitEnumNames ? '' : 'QUERY_INSIGHTS_MODE_DISABLED');
  static const QueryInsightsMode QUERY_INSIGHTS_MODE_ENABLED_WITH_RATE_CONTROL =
      QueryInsightsMode._(
          1,
          _omitEnumNames
              ? ''
              : 'QUERY_INSIGHTS_MODE_ENABLED_WITH_RATE_CONTROL');

  static const $core.List<QueryInsightsMode> values = <QueryInsightsMode>[
    QUERY_INSIGHTS_MODE_DISABLED,
    QUERY_INSIGHTS_MODE_ENABLED_WITH_RATE_CONTROL,
  ];

  static final $core.List<QueryInsightsMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static QueryInsightsMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueryInsightsMode._(super.value, super.name);
}

class QueryPricingModel extends $pb.ProtobufEnum {
  static const QueryPricingModel QUERY_PRICING_MODEL_COMPUTE_UNITS =
      QueryPricingModel._(
          0, _omitEnumNames ? '' : 'QUERY_PRICING_MODEL_COMPUTE_UNITS');
  static const QueryPricingModel QUERY_PRICING_MODEL_BYTES_SCANNED =
      QueryPricingModel._(
          1, _omitEnumNames ? '' : 'QUERY_PRICING_MODEL_BYTES_SCANNED');

  static const $core.List<QueryPricingModel> values = <QueryPricingModel>[
    QUERY_PRICING_MODEL_COMPUTE_UNITS,
    QUERY_PRICING_MODEL_BYTES_SCANNED,
  ];

  static final $core.List<QueryPricingModel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static QueryPricingModel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueryPricingModel._(super.value, super.name);
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

class ScalarType extends $pb.ProtobufEnum {
  static const ScalarType SCALAR_TYPE_BIGINT =
      ScalarType._(0, _omitEnumNames ? '' : 'SCALAR_TYPE_BIGINT');
  static const ScalarType SCALAR_TYPE_UNKNOWN =
      ScalarType._(1, _omitEnumNames ? '' : 'SCALAR_TYPE_UNKNOWN');
  static const ScalarType SCALAR_TYPE_INTERVAL_DAY_TO_SECOND = ScalarType._(
      2, _omitEnumNames ? '' : 'SCALAR_TYPE_INTERVAL_DAY_TO_SECOND');
  static const ScalarType SCALAR_TYPE_INTEGER =
      ScalarType._(3, _omitEnumNames ? '' : 'SCALAR_TYPE_INTEGER');
  static const ScalarType SCALAR_TYPE_TIME =
      ScalarType._(4, _omitEnumNames ? '' : 'SCALAR_TYPE_TIME');
  static const ScalarType SCALAR_TYPE_VARCHAR =
      ScalarType._(5, _omitEnumNames ? '' : 'SCALAR_TYPE_VARCHAR');
  static const ScalarType SCALAR_TYPE_BOOLEAN =
      ScalarType._(6, _omitEnumNames ? '' : 'SCALAR_TYPE_BOOLEAN');
  static const ScalarType SCALAR_TYPE_DOUBLE =
      ScalarType._(7, _omitEnumNames ? '' : 'SCALAR_TYPE_DOUBLE');
  static const ScalarType SCALAR_TYPE_TIMESTAMP =
      ScalarType._(8, _omitEnumNames ? '' : 'SCALAR_TYPE_TIMESTAMP');
  static const ScalarType SCALAR_TYPE_DATE =
      ScalarType._(9, _omitEnumNames ? '' : 'SCALAR_TYPE_DATE');
  static const ScalarType SCALAR_TYPE_INTERVAL_YEAR_TO_MONTH = ScalarType._(
      10, _omitEnumNames ? '' : 'SCALAR_TYPE_INTERVAL_YEAR_TO_MONTH');

  static const $core.List<ScalarType> values = <ScalarType>[
    SCALAR_TYPE_BIGINT,
    SCALAR_TYPE_UNKNOWN,
    SCALAR_TYPE_INTERVAL_DAY_TO_SECOND,
    SCALAR_TYPE_INTEGER,
    SCALAR_TYPE_TIME,
    SCALAR_TYPE_VARCHAR,
    SCALAR_TYPE_BOOLEAN,
    SCALAR_TYPE_DOUBLE,
    SCALAR_TYPE_TIMESTAMP,
    SCALAR_TYPE_DATE,
    SCALAR_TYPE_INTERVAL_YEAR_TO_MONTH,
  ];

  static final $core.List<ScalarType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 10);
  static ScalarType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScalarType._(super.value, super.name);
}

class ScheduledQueryInsightsMode extends $pb.ProtobufEnum {
  static const ScheduledQueryInsightsMode
      SCHEDULED_QUERY_INSIGHTS_MODE_DISABLED = ScheduledQueryInsightsMode._(
          0, _omitEnumNames ? '' : 'SCHEDULED_QUERY_INSIGHTS_MODE_DISABLED');
  static const ScheduledQueryInsightsMode
      SCHEDULED_QUERY_INSIGHTS_MODE_ENABLED_WITH_RATE_CONTROL =
      ScheduledQueryInsightsMode._(
          1,
          _omitEnumNames
              ? ''
              : 'SCHEDULED_QUERY_INSIGHTS_MODE_ENABLED_WITH_RATE_CONTROL');

  static const $core.List<ScheduledQueryInsightsMode> values =
      <ScheduledQueryInsightsMode>[
    SCHEDULED_QUERY_INSIGHTS_MODE_DISABLED,
    SCHEDULED_QUERY_INSIGHTS_MODE_ENABLED_WITH_RATE_CONTROL,
  ];

  static final $core.List<ScheduledQueryInsightsMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ScheduledQueryInsightsMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScheduledQueryInsightsMode._(super.value, super.name);
}

class ScheduledQueryRunStatus extends $pb.ProtobufEnum {
  static const ScheduledQueryRunStatus
      SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_SUCCESS =
      ScheduledQueryRunStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_SUCCESS');
  static const ScheduledQueryRunStatus
      SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_FAILURE =
      ScheduledQueryRunStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_FAILURE');
  static const ScheduledQueryRunStatus
      SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_FAILURE =
      ScheduledQueryRunStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_FAILURE');
  static const ScheduledQueryRunStatus
      SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_SUCCESS =
      ScheduledQueryRunStatus._(
          3,
          _omitEnumNames
              ? ''
              : 'SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_SUCCESS');

  static const $core.List<ScheduledQueryRunStatus> values =
      <ScheduledQueryRunStatus>[
    SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_SUCCESS,
    SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_FAILURE,
    SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_FAILURE,
    SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_SUCCESS,
  ];

  static final $core.List<ScheduledQueryRunStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ScheduledQueryRunStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScheduledQueryRunStatus._(super.value, super.name);
}

class ScheduledQueryState extends $pb.ProtobufEnum {
  static const ScheduledQueryState SCHEDULED_QUERY_STATE_DISABLED =
      ScheduledQueryState._(
          0, _omitEnumNames ? '' : 'SCHEDULED_QUERY_STATE_DISABLED');
  static const ScheduledQueryState SCHEDULED_QUERY_STATE_ENABLED =
      ScheduledQueryState._(
          1, _omitEnumNames ? '' : 'SCHEDULED_QUERY_STATE_ENABLED');

  static const $core.List<ScheduledQueryState> values = <ScheduledQueryState>[
    SCHEDULED_QUERY_STATE_DISABLED,
    SCHEDULED_QUERY_STATE_ENABLED,
  ];

  static final $core.List<ScheduledQueryState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ScheduledQueryState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScheduledQueryState._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
