// This is a generated file - do not edit.
//
// Generated from kinesis.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ConsumerStatus extends $pb.ProtobufEnum {
  static const ConsumerStatus CONSUMER_STATUS_ACTIVE =
      ConsumerStatus._(0, _omitEnumNames ? '' : 'CONSUMER_STATUS_ACTIVE');
  static const ConsumerStatus CONSUMER_STATUS_DELETING =
      ConsumerStatus._(1, _omitEnumNames ? '' : 'CONSUMER_STATUS_DELETING');
  static const ConsumerStatus CONSUMER_STATUS_CREATING =
      ConsumerStatus._(2, _omitEnumNames ? '' : 'CONSUMER_STATUS_CREATING');

  static const $core.List<ConsumerStatus> values = <ConsumerStatus>[
    CONSUMER_STATUS_ACTIVE,
    CONSUMER_STATUS_DELETING,
    CONSUMER_STATUS_CREATING,
  ];

  static final $core.List<ConsumerStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ConsumerStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConsumerStatus._(super.value, super.name);
}

class EncryptionType extends $pb.ProtobufEnum {
  static const EncryptionType ENCRYPTION_TYPE_NONE =
      EncryptionType._(0, _omitEnumNames ? '' : 'ENCRYPTION_TYPE_NONE');
  static const EncryptionType ENCRYPTION_TYPE_KMS =
      EncryptionType._(1, _omitEnumNames ? '' : 'ENCRYPTION_TYPE_KMS');

  static const $core.List<EncryptionType> values = <EncryptionType>[
    ENCRYPTION_TYPE_NONE,
    ENCRYPTION_TYPE_KMS,
  ];

  static final $core.List<EncryptionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static EncryptionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EncryptionType._(super.value, super.name);
}

class MetricsName extends $pb.ProtobufEnum {
  static const MetricsName METRICS_NAME_INCOMING_RECORDS =
      MetricsName._(0, _omitEnumNames ? '' : 'METRICS_NAME_INCOMING_RECORDS');
  static const MetricsName METRICS_NAME_INCOMING_BYTES =
      MetricsName._(1, _omitEnumNames ? '' : 'METRICS_NAME_INCOMING_BYTES');
  static const MetricsName METRICS_NAME_WRITE_PROVISIONED_THROUGHPUT_EXCEEDED =
      MetricsName._(
          2,
          _omitEnumNames
              ? ''
              : 'METRICS_NAME_WRITE_PROVISIONED_THROUGHPUT_EXCEEDED');
  static const MetricsName METRICS_NAME_READ_PROVISIONED_THROUGHPUT_EXCEEDED =
      MetricsName._(
          3,
          _omitEnumNames
              ? ''
              : 'METRICS_NAME_READ_PROVISIONED_THROUGHPUT_EXCEEDED');
  static const MetricsName METRICS_NAME_OUTGOING_BYTES =
      MetricsName._(4, _omitEnumNames ? '' : 'METRICS_NAME_OUTGOING_BYTES');
  static const MetricsName METRICS_NAME_OUTGOING_RECORDS =
      MetricsName._(5, _omitEnumNames ? '' : 'METRICS_NAME_OUTGOING_RECORDS');
  static const MetricsName METRICS_NAME_ITERATOR_AGE_MILLISECONDS =
      MetricsName._(
          6, _omitEnumNames ? '' : 'METRICS_NAME_ITERATOR_AGE_MILLISECONDS');
  static const MetricsName METRICS_NAME_ALL =
      MetricsName._(7, _omitEnumNames ? '' : 'METRICS_NAME_ALL');

  static const $core.List<MetricsName> values = <MetricsName>[
    METRICS_NAME_INCOMING_RECORDS,
    METRICS_NAME_INCOMING_BYTES,
    METRICS_NAME_WRITE_PROVISIONED_THROUGHPUT_EXCEEDED,
    METRICS_NAME_READ_PROVISIONED_THROUGHPUT_EXCEEDED,
    METRICS_NAME_OUTGOING_BYTES,
    METRICS_NAME_OUTGOING_RECORDS,
    METRICS_NAME_ITERATOR_AGE_MILLISECONDS,
    METRICS_NAME_ALL,
  ];

  static final $core.List<MetricsName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static MetricsName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetricsName._(super.value, super.name);
}

class MinimumThroughputBillingCommitmentInputStatus extends $pb.ProtobufEnum {
  static const MinimumThroughputBillingCommitmentInputStatus
      MINIMUM_THROUGHPUT_BILLING_COMMITMENT_INPUT_STATUS_DISABLED =
      MinimumThroughputBillingCommitmentInputStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_INPUT_STATUS_DISABLED');
  static const MinimumThroughputBillingCommitmentInputStatus
      MINIMUM_THROUGHPUT_BILLING_COMMITMENT_INPUT_STATUS_ENABLED =
      MinimumThroughputBillingCommitmentInputStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_INPUT_STATUS_ENABLED');

  static const $core.List<MinimumThroughputBillingCommitmentInputStatus>
      values = <MinimumThroughputBillingCommitmentInputStatus>[
    MINIMUM_THROUGHPUT_BILLING_COMMITMENT_INPUT_STATUS_DISABLED,
    MINIMUM_THROUGHPUT_BILLING_COMMITMENT_INPUT_STATUS_ENABLED,
  ];

  static final $core.List<MinimumThroughputBillingCommitmentInputStatus?>
      _byValue = $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MinimumThroughputBillingCommitmentInputStatus? valueOf(
          $core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MinimumThroughputBillingCommitmentInputStatus._(
      super.value, super.name);
}

class MinimumThroughputBillingCommitmentOutputStatus extends $pb.ProtobufEnum {
  static const MinimumThroughputBillingCommitmentOutputStatus
      MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_DISABLED =
      MinimumThroughputBillingCommitmentOutputStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_DISABLED');
  static const MinimumThroughputBillingCommitmentOutputStatus
      MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_ENABLED =
      MinimumThroughputBillingCommitmentOutputStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_ENABLED');
  static const MinimumThroughputBillingCommitmentOutputStatus
      MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_ENABLED_UNTIL_EARLIEST_ALLOWED_END =
      MinimumThroughputBillingCommitmentOutputStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_ENABLED_UNTIL_EARLIEST_ALLOWED_END');

  static const $core.List<MinimumThroughputBillingCommitmentOutputStatus>
      values = <MinimumThroughputBillingCommitmentOutputStatus>[
    MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_DISABLED,
    MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_ENABLED,
    MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_ENABLED_UNTIL_EARLIEST_ALLOWED_END,
  ];

  static final $core.List<MinimumThroughputBillingCommitmentOutputStatus?>
      _byValue = $pb.ProtobufEnum.$_initByValueList(values, 2);
  static MinimumThroughputBillingCommitmentOutputStatus? valueOf(
          $core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MinimumThroughputBillingCommitmentOutputStatus._(
      super.value, super.name);
}

class ScalingType extends $pb.ProtobufEnum {
  static const ScalingType SCALING_TYPE_UNIFORM_SCALING =
      ScalingType._(0, _omitEnumNames ? '' : 'SCALING_TYPE_UNIFORM_SCALING');

  static const $core.List<ScalingType> values = <ScalingType>[
    SCALING_TYPE_UNIFORM_SCALING,
  ];

  static final $core.List<ScalingType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ScalingType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScalingType._(super.value, super.name);
}

class ShardFilterType extends $pb.ProtobufEnum {
  static const ShardFilterType SHARD_FILTER_TYPE_AT_TIMESTAMP =
      ShardFilterType._(
          0, _omitEnumNames ? '' : 'SHARD_FILTER_TYPE_AT_TIMESTAMP');
  static const ShardFilterType SHARD_FILTER_TYPE_FROM_TRIM_HORIZON =
      ShardFilterType._(
          1, _omitEnumNames ? '' : 'SHARD_FILTER_TYPE_FROM_TRIM_HORIZON');
  static const ShardFilterType SHARD_FILTER_TYPE_AT_LATEST =
      ShardFilterType._(2, _omitEnumNames ? '' : 'SHARD_FILTER_TYPE_AT_LATEST');
  static const ShardFilterType SHARD_FILTER_TYPE_AT_TRIM_HORIZON =
      ShardFilterType._(
          3, _omitEnumNames ? '' : 'SHARD_FILTER_TYPE_AT_TRIM_HORIZON');
  static const ShardFilterType SHARD_FILTER_TYPE_AFTER_SHARD_ID =
      ShardFilterType._(
          4, _omitEnumNames ? '' : 'SHARD_FILTER_TYPE_AFTER_SHARD_ID');
  static const ShardFilterType SHARD_FILTER_TYPE_FROM_TIMESTAMP =
      ShardFilterType._(
          5, _omitEnumNames ? '' : 'SHARD_FILTER_TYPE_FROM_TIMESTAMP');

  static const $core.List<ShardFilterType> values = <ShardFilterType>[
    SHARD_FILTER_TYPE_AT_TIMESTAMP,
    SHARD_FILTER_TYPE_FROM_TRIM_HORIZON,
    SHARD_FILTER_TYPE_AT_LATEST,
    SHARD_FILTER_TYPE_AT_TRIM_HORIZON,
    SHARD_FILTER_TYPE_AFTER_SHARD_ID,
    SHARD_FILTER_TYPE_FROM_TIMESTAMP,
  ];

  static final $core.List<ShardFilterType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ShardFilterType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ShardFilterType._(super.value, super.name);
}

class ShardIteratorType extends $pb.ProtobufEnum {
  static const ShardIteratorType SHARD_ITERATOR_TYPE_AT_TIMESTAMP =
      ShardIteratorType._(
          0, _omitEnumNames ? '' : 'SHARD_ITERATOR_TYPE_AT_TIMESTAMP');
  static const ShardIteratorType SHARD_ITERATOR_TYPE_TRIM_HORIZON =
      ShardIteratorType._(
          1, _omitEnumNames ? '' : 'SHARD_ITERATOR_TYPE_TRIM_HORIZON');
  static const ShardIteratorType SHARD_ITERATOR_TYPE_AFTER_SEQUENCE_NUMBER =
      ShardIteratorType._(
          2, _omitEnumNames ? '' : 'SHARD_ITERATOR_TYPE_AFTER_SEQUENCE_NUMBER');
  static const ShardIteratorType SHARD_ITERATOR_TYPE_AT_SEQUENCE_NUMBER =
      ShardIteratorType._(
          3, _omitEnumNames ? '' : 'SHARD_ITERATOR_TYPE_AT_SEQUENCE_NUMBER');
  static const ShardIteratorType SHARD_ITERATOR_TYPE_LATEST =
      ShardIteratorType._(
          4, _omitEnumNames ? '' : 'SHARD_ITERATOR_TYPE_LATEST');

  static const $core.List<ShardIteratorType> values = <ShardIteratorType>[
    SHARD_ITERATOR_TYPE_AT_TIMESTAMP,
    SHARD_ITERATOR_TYPE_TRIM_HORIZON,
    SHARD_ITERATOR_TYPE_AFTER_SEQUENCE_NUMBER,
    SHARD_ITERATOR_TYPE_AT_SEQUENCE_NUMBER,
    SHARD_ITERATOR_TYPE_LATEST,
  ];

  static final $core.List<ShardIteratorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ShardIteratorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ShardIteratorType._(super.value, super.name);
}

class StreamMode extends $pb.ProtobufEnum {
  static const StreamMode STREAM_MODE_PROVISIONED =
      StreamMode._(0, _omitEnumNames ? '' : 'STREAM_MODE_PROVISIONED');
  static const StreamMode STREAM_MODE_ON_DEMAND =
      StreamMode._(1, _omitEnumNames ? '' : 'STREAM_MODE_ON_DEMAND');

  static const $core.List<StreamMode> values = <StreamMode>[
    STREAM_MODE_PROVISIONED,
    STREAM_MODE_ON_DEMAND,
  ];

  static final $core.List<StreamMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static StreamMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StreamMode._(super.value, super.name);
}

class StreamStatus extends $pb.ProtobufEnum {
  static const StreamStatus STREAM_STATUS_UPDATING =
      StreamStatus._(0, _omitEnumNames ? '' : 'STREAM_STATUS_UPDATING');
  static const StreamStatus STREAM_STATUS_ACTIVE =
      StreamStatus._(1, _omitEnumNames ? '' : 'STREAM_STATUS_ACTIVE');
  static const StreamStatus STREAM_STATUS_DELETING =
      StreamStatus._(2, _omitEnumNames ? '' : 'STREAM_STATUS_DELETING');
  static const StreamStatus STREAM_STATUS_CREATING =
      StreamStatus._(3, _omitEnumNames ? '' : 'STREAM_STATUS_CREATING');

  static const $core.List<StreamStatus> values = <StreamStatus>[
    STREAM_STATUS_UPDATING,
    STREAM_STATUS_ACTIVE,
    STREAM_STATUS_DELETING,
    STREAM_STATUS_CREATING,
  ];

  static final $core.List<StreamStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static StreamStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StreamStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
