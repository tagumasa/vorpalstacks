// This is a generated file - do not edit.
//
// Generated from monitoring.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class ActionsSuppressedBy extends $pb.ProtobufEnum {
  static const ActionsSuppressedBy ACTIONS_SUPPRESSED_BY_WAITPERIOD =
      ActionsSuppressedBy._(
          0, _omitEnumNames ? '' : 'ACTIONS_SUPPRESSED_BY_WAITPERIOD');
  static const ActionsSuppressedBy ACTIONS_SUPPRESSED_BY_ALARM =
      ActionsSuppressedBy._(
          1, _omitEnumNames ? '' : 'ACTIONS_SUPPRESSED_BY_ALARM');
  static const ActionsSuppressedBy ACTIONS_SUPPRESSED_BY_EXTENSIONPERIOD =
      ActionsSuppressedBy._(
          2, _omitEnumNames ? '' : 'ACTIONS_SUPPRESSED_BY_EXTENSIONPERIOD');

  static const $core.List<ActionsSuppressedBy> values = <ActionsSuppressedBy>[
    ACTIONS_SUPPRESSED_BY_WAITPERIOD,
    ACTIONS_SUPPRESSED_BY_ALARM,
    ACTIONS_SUPPRESSED_BY_EXTENSIONPERIOD,
  ];

  static final $core.List<ActionsSuppressedBy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ActionsSuppressedBy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ActionsSuppressedBy._(super.value, super.name);
}

class AlarmMuteRuleStatus extends $pb.ProtobufEnum {
  static const AlarmMuteRuleStatus ALARM_MUTE_RULE_STATUS_ACTIVE =
      AlarmMuteRuleStatus._(
          0, _omitEnumNames ? '' : 'ALARM_MUTE_RULE_STATUS_ACTIVE');
  static const AlarmMuteRuleStatus ALARM_MUTE_RULE_STATUS_SCHEDULED =
      AlarmMuteRuleStatus._(
          1, _omitEnumNames ? '' : 'ALARM_MUTE_RULE_STATUS_SCHEDULED');
  static const AlarmMuteRuleStatus ALARM_MUTE_RULE_STATUS_EXPIRED =
      AlarmMuteRuleStatus._(
          2, _omitEnumNames ? '' : 'ALARM_MUTE_RULE_STATUS_EXPIRED');

  static const $core.List<AlarmMuteRuleStatus> values = <AlarmMuteRuleStatus>[
    ALARM_MUTE_RULE_STATUS_ACTIVE,
    ALARM_MUTE_RULE_STATUS_SCHEDULED,
    ALARM_MUTE_RULE_STATUS_EXPIRED,
  ];

  static final $core.List<AlarmMuteRuleStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AlarmMuteRuleStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AlarmMuteRuleStatus._(super.value, super.name);
}

class AlarmType extends $pb.ProtobufEnum {
  static const AlarmType ALARM_TYPE_COMPOSITEALARM =
      AlarmType._(0, _omitEnumNames ? '' : 'ALARM_TYPE_COMPOSITEALARM');
  static const AlarmType ALARM_TYPE_METRICALARM =
      AlarmType._(1, _omitEnumNames ? '' : 'ALARM_TYPE_METRICALARM');

  static const $core.List<AlarmType> values = <AlarmType>[
    ALARM_TYPE_COMPOSITEALARM,
    ALARM_TYPE_METRICALARM,
  ];

  static final $core.List<AlarmType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AlarmType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AlarmType._(super.value, super.name);
}

class AnomalyDetectorStateValue extends $pb.ProtobufEnum {
  static const AnomalyDetectorStateValue ANOMALY_DETECTOR_STATE_VALUE_TRAINED =
      AnomalyDetectorStateValue._(
          0, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_STATE_VALUE_TRAINED');
  static const AnomalyDetectorStateValue
      ANOMALY_DETECTOR_STATE_VALUE_PENDING_TRAINING =
      AnomalyDetectorStateValue._(
          1,
          _omitEnumNames
              ? ''
              : 'ANOMALY_DETECTOR_STATE_VALUE_PENDING_TRAINING');
  static const AnomalyDetectorStateValue
      ANOMALY_DETECTOR_STATE_VALUE_TRAINED_INSUFFICIENT_DATA =
      AnomalyDetectorStateValue._(
          2,
          _omitEnumNames
              ? ''
              : 'ANOMALY_DETECTOR_STATE_VALUE_TRAINED_INSUFFICIENT_DATA');

  static const $core.List<AnomalyDetectorStateValue> values =
      <AnomalyDetectorStateValue>[
    ANOMALY_DETECTOR_STATE_VALUE_TRAINED,
    ANOMALY_DETECTOR_STATE_VALUE_PENDING_TRAINING,
    ANOMALY_DETECTOR_STATE_VALUE_TRAINED_INSUFFICIENT_DATA,
  ];

  static final $core.List<AnomalyDetectorStateValue?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AnomalyDetectorStateValue? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AnomalyDetectorStateValue._(super.value, super.name);
}

class AnomalyDetectorType extends $pb.ProtobufEnum {
  static const AnomalyDetectorType ANOMALY_DETECTOR_TYPE_SINGLE_METRIC =
      AnomalyDetectorType._(
          0, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_TYPE_SINGLE_METRIC');
  static const AnomalyDetectorType ANOMALY_DETECTOR_TYPE_METRIC_MATH =
      AnomalyDetectorType._(
          1, _omitEnumNames ? '' : 'ANOMALY_DETECTOR_TYPE_METRIC_MATH');

  static const $core.List<AnomalyDetectorType> values = <AnomalyDetectorType>[
    ANOMALY_DETECTOR_TYPE_SINGLE_METRIC,
    ANOMALY_DETECTOR_TYPE_METRIC_MATH,
  ];

  static final $core.List<AnomalyDetectorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AnomalyDetectorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AnomalyDetectorType._(super.value, super.name);
}

class ComparisonOperator extends $pb.ProtobufEnum {
  static const ComparisonOperator
      COMPARISON_OPERATOR_LESSTHANLOWERORGREATERTHANUPPERTHRESHOLD =
      ComparisonOperator._(
          0,
          _omitEnumNames
              ? ''
              : 'COMPARISON_OPERATOR_LESSTHANLOWERORGREATERTHANUPPERTHRESHOLD');
  static const ComparisonOperator COMPARISON_OPERATOR_LESSTHANLOWERTHRESHOLD =
      ComparisonOperator._(1,
          _omitEnumNames ? '' : 'COMPARISON_OPERATOR_LESSTHANLOWERTHRESHOLD');
  static const ComparisonOperator
      COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD = ComparisonOperator._(
          2,
          _omitEnumNames
              ? ''
              : 'COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD');
  static const ComparisonOperator
      COMPARISON_OPERATOR_GREATERTHANUPPERTHRESHOLD = ComparisonOperator._(
          3,
          _omitEnumNames
              ? ''
              : 'COMPARISON_OPERATOR_GREATERTHANUPPERTHRESHOLD');
  static const ComparisonOperator COMPARISON_OPERATOR_GREATERTHANTHRESHOLD =
      ComparisonOperator._(
          4, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_GREATERTHANTHRESHOLD');
  static const ComparisonOperator COMPARISON_OPERATOR_LESSTHANTHRESHOLD =
      ComparisonOperator._(
          5, _omitEnumNames ? '' : 'COMPARISON_OPERATOR_LESSTHANTHRESHOLD');
  static const ComparisonOperator
      COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD = ComparisonOperator._(
          6,
          _omitEnumNames
              ? ''
              : 'COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD');

  static const $core.List<ComparisonOperator> values = <ComparisonOperator>[
    COMPARISON_OPERATOR_LESSTHANLOWERORGREATERTHANUPPERTHRESHOLD,
    COMPARISON_OPERATOR_LESSTHANLOWERTHRESHOLD,
    COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD,
    COMPARISON_OPERATOR_GREATERTHANUPPERTHRESHOLD,
    COMPARISON_OPERATOR_GREATERTHANTHRESHOLD,
    COMPARISON_OPERATOR_LESSTHANTHRESHOLD,
    COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD,
  ];

  static final $core.List<ComparisonOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static ComparisonOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComparisonOperator._(super.value, super.name);
}

class EvaluationState extends $pb.ProtobufEnum {
  static const EvaluationState EVALUATION_STATE_PARTIAL_DATA =
      EvaluationState._(
          0, _omitEnumNames ? '' : 'EVALUATION_STATE_PARTIAL_DATA');
  static const EvaluationState EVALUATION_STATE_EVALUATION_ERROR =
      EvaluationState._(
          1, _omitEnumNames ? '' : 'EVALUATION_STATE_EVALUATION_ERROR');
  static const EvaluationState EVALUATION_STATE_EVALUATION_FAILURE =
      EvaluationState._(
          2, _omitEnumNames ? '' : 'EVALUATION_STATE_EVALUATION_FAILURE');

  static const $core.List<EvaluationState> values = <EvaluationState>[
    EVALUATION_STATE_PARTIAL_DATA,
    EVALUATION_STATE_EVALUATION_ERROR,
    EVALUATION_STATE_EVALUATION_FAILURE,
  ];

  static final $core.List<EvaluationState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EvaluationState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EvaluationState._(super.value, super.name);
}

class HistoryItemType extends $pb.ProtobufEnum {
  static const HistoryItemType HISTORY_ITEM_TYPE_STATEUPDATE =
      HistoryItemType._(
          0, _omitEnumNames ? '' : 'HISTORY_ITEM_TYPE_STATEUPDATE');
  static const HistoryItemType HISTORY_ITEM_TYPE_ALARMCONTRIBUTORACTION =
      HistoryItemType._(
          1, _omitEnumNames ? '' : 'HISTORY_ITEM_TYPE_ALARMCONTRIBUTORACTION');
  static const HistoryItemType HISTORY_ITEM_TYPE_ACTION =
      HistoryItemType._(2, _omitEnumNames ? '' : 'HISTORY_ITEM_TYPE_ACTION');
  static const HistoryItemType HISTORY_ITEM_TYPE_ALARMCONTRIBUTORSTATEUPDATE =
      HistoryItemType._(
          3,
          _omitEnumNames
              ? ''
              : 'HISTORY_ITEM_TYPE_ALARMCONTRIBUTORSTATEUPDATE');
  static const HistoryItemType HISTORY_ITEM_TYPE_CONFIGURATIONUPDATE =
      HistoryItemType._(
          4, _omitEnumNames ? '' : 'HISTORY_ITEM_TYPE_CONFIGURATIONUPDATE');

  static const $core.List<HistoryItemType> values = <HistoryItemType>[
    HISTORY_ITEM_TYPE_STATEUPDATE,
    HISTORY_ITEM_TYPE_ALARMCONTRIBUTORACTION,
    HISTORY_ITEM_TYPE_ACTION,
    HISTORY_ITEM_TYPE_ALARMCONTRIBUTORSTATEUPDATE,
    HISTORY_ITEM_TYPE_CONFIGURATIONUPDATE,
  ];

  static final $core.List<HistoryItemType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static HistoryItemType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HistoryItemType._(super.value, super.name);
}

class MetricStreamOutputFormat extends $pb.ProtobufEnum {
  static const MetricStreamOutputFormat
      METRIC_STREAM_OUTPUT_FORMAT_OPEN_TELEMETRY_0_7 =
      MetricStreamOutputFormat._(
          0,
          _omitEnumNames
              ? ''
              : 'METRIC_STREAM_OUTPUT_FORMAT_OPEN_TELEMETRY_0_7');
  static const MetricStreamOutputFormat METRIC_STREAM_OUTPUT_FORMAT_JSON =
      MetricStreamOutputFormat._(
          1, _omitEnumNames ? '' : 'METRIC_STREAM_OUTPUT_FORMAT_JSON');
  static const MetricStreamOutputFormat
      METRIC_STREAM_OUTPUT_FORMAT_OPEN_TELEMETRY_1_0 =
      MetricStreamOutputFormat._(
          2,
          _omitEnumNames
              ? ''
              : 'METRIC_STREAM_OUTPUT_FORMAT_OPEN_TELEMETRY_1_0');

  static const $core.List<MetricStreamOutputFormat> values =
      <MetricStreamOutputFormat>[
    METRIC_STREAM_OUTPUT_FORMAT_OPEN_TELEMETRY_0_7,
    METRIC_STREAM_OUTPUT_FORMAT_JSON,
    METRIC_STREAM_OUTPUT_FORMAT_OPEN_TELEMETRY_1_0,
  ];

  static final $core.List<MetricStreamOutputFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static MetricStreamOutputFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetricStreamOutputFormat._(super.value, super.name);
}

class RecentlyActive extends $pb.ProtobufEnum {
  static const RecentlyActive RECENTLY_ACTIVE_PT3H =
      RecentlyActive._(0, _omitEnumNames ? '' : 'RECENTLY_ACTIVE_PT3H');

  static const $core.List<RecentlyActive> values = <RecentlyActive>[
    RECENTLY_ACTIVE_PT3H,
  ];

  static final $core.List<RecentlyActive?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static RecentlyActive? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RecentlyActive._(super.value, super.name);
}

class ScanBy extends $pb.ProtobufEnum {
  static const ScanBy SCAN_BY_TIMESTAMP_DESCENDING =
      ScanBy._(0, _omitEnumNames ? '' : 'SCAN_BY_TIMESTAMP_DESCENDING');
  static const ScanBy SCAN_BY_TIMESTAMP_ASCENDING =
      ScanBy._(1, _omitEnumNames ? '' : 'SCAN_BY_TIMESTAMP_ASCENDING');

  static const $core.List<ScanBy> values = <ScanBy>[
    SCAN_BY_TIMESTAMP_DESCENDING,
    SCAN_BY_TIMESTAMP_ASCENDING,
  ];

  static final $core.List<ScanBy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ScanBy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScanBy._(super.value, super.name);
}

class StandardUnit extends $pb.ProtobufEnum {
  static const StandardUnit STANDARD_UNIT_MEGABITS_SECOND =
      StandardUnit._(0, _omitEnumNames ? '' : 'STANDARD_UNIT_MEGABITS_SECOND');
  static const StandardUnit STANDARD_UNIT_KILOBYTES_SECOND =
      StandardUnit._(1, _omitEnumNames ? '' : 'STANDARD_UNIT_KILOBYTES_SECOND');
  static const StandardUnit STANDARD_UNIT_KILOBITS_SECOND =
      StandardUnit._(2, _omitEnumNames ? '' : 'STANDARD_UNIT_KILOBITS_SECOND');
  static const StandardUnit STANDARD_UNIT_GIGABITS =
      StandardUnit._(3, _omitEnumNames ? '' : 'STANDARD_UNIT_GIGABITS');
  static const StandardUnit STANDARD_UNIT_NONE =
      StandardUnit._(4, _omitEnumNames ? '' : 'STANDARD_UNIT_NONE');
  static const StandardUnit STANDARD_UNIT_MILLISECONDS =
      StandardUnit._(5, _omitEnumNames ? '' : 'STANDARD_UNIT_MILLISECONDS');
  static const StandardUnit STANDARD_UNIT_TERABITS_SECOND =
      StandardUnit._(6, _omitEnumNames ? '' : 'STANDARD_UNIT_TERABITS_SECOND');
  static const StandardUnit STANDARD_UNIT_BYTES =
      StandardUnit._(7, _omitEnumNames ? '' : 'STANDARD_UNIT_BYTES');
  static const StandardUnit STANDARD_UNIT_GIGABYTES =
      StandardUnit._(8, _omitEnumNames ? '' : 'STANDARD_UNIT_GIGABYTES');
  static const StandardUnit STANDARD_UNIT_MEGABYTES_SECOND =
      StandardUnit._(9, _omitEnumNames ? '' : 'STANDARD_UNIT_MEGABYTES_SECOND');
  static const StandardUnit STANDARD_UNIT_BITS =
      StandardUnit._(10, _omitEnumNames ? '' : 'STANDARD_UNIT_BITS');
  static const StandardUnit STANDARD_UNIT_GIGABITS_SECOND =
      StandardUnit._(11, _omitEnumNames ? '' : 'STANDARD_UNIT_GIGABITS_SECOND');
  static const StandardUnit STANDARD_UNIT_MEGABYTES =
      StandardUnit._(12, _omitEnumNames ? '' : 'STANDARD_UNIT_MEGABYTES');
  static const StandardUnit STANDARD_UNIT_BYTES_SECOND =
      StandardUnit._(13, _omitEnumNames ? '' : 'STANDARD_UNIT_BYTES_SECOND');
  static const StandardUnit STANDARD_UNIT_MEGABITS =
      StandardUnit._(14, _omitEnumNames ? '' : 'STANDARD_UNIT_MEGABITS');
  static const StandardUnit STANDARD_UNIT_SECONDS =
      StandardUnit._(15, _omitEnumNames ? '' : 'STANDARD_UNIT_SECONDS');
  static const StandardUnit STANDARD_UNIT_COUNT_SECOND =
      StandardUnit._(16, _omitEnumNames ? '' : 'STANDARD_UNIT_COUNT_SECOND');
  static const StandardUnit STANDARD_UNIT_PERCENT =
      StandardUnit._(17, _omitEnumNames ? '' : 'STANDARD_UNIT_PERCENT');
  static const StandardUnit STANDARD_UNIT_COUNT =
      StandardUnit._(18, _omitEnumNames ? '' : 'STANDARD_UNIT_COUNT');
  static const StandardUnit STANDARD_UNIT_KILOBITS =
      StandardUnit._(19, _omitEnumNames ? '' : 'STANDARD_UNIT_KILOBITS');
  static const StandardUnit STANDARD_UNIT_KILOBYTES =
      StandardUnit._(20, _omitEnumNames ? '' : 'STANDARD_UNIT_KILOBYTES');
  static const StandardUnit STANDARD_UNIT_TERABYTES =
      StandardUnit._(21, _omitEnumNames ? '' : 'STANDARD_UNIT_TERABYTES');
  static const StandardUnit STANDARD_UNIT_GIGABYTES_SECOND = StandardUnit._(
      22, _omitEnumNames ? '' : 'STANDARD_UNIT_GIGABYTES_SECOND');
  static const StandardUnit STANDARD_UNIT_MICROSECONDS =
      StandardUnit._(23, _omitEnumNames ? '' : 'STANDARD_UNIT_MICROSECONDS');
  static const StandardUnit STANDARD_UNIT_TERABITS =
      StandardUnit._(24, _omitEnumNames ? '' : 'STANDARD_UNIT_TERABITS');
  static const StandardUnit STANDARD_UNIT_BITS_SECOND =
      StandardUnit._(25, _omitEnumNames ? '' : 'STANDARD_UNIT_BITS_SECOND');
  static const StandardUnit STANDARD_UNIT_TERABYTES_SECOND = StandardUnit._(
      26, _omitEnumNames ? '' : 'STANDARD_UNIT_TERABYTES_SECOND');

  static const $core.List<StandardUnit> values = <StandardUnit>[
    STANDARD_UNIT_MEGABITS_SECOND,
    STANDARD_UNIT_KILOBYTES_SECOND,
    STANDARD_UNIT_KILOBITS_SECOND,
    STANDARD_UNIT_GIGABITS,
    STANDARD_UNIT_NONE,
    STANDARD_UNIT_MILLISECONDS,
    STANDARD_UNIT_TERABITS_SECOND,
    STANDARD_UNIT_BYTES,
    STANDARD_UNIT_GIGABYTES,
    STANDARD_UNIT_MEGABYTES_SECOND,
    STANDARD_UNIT_BITS,
    STANDARD_UNIT_GIGABITS_SECOND,
    STANDARD_UNIT_MEGABYTES,
    STANDARD_UNIT_BYTES_SECOND,
    STANDARD_UNIT_MEGABITS,
    STANDARD_UNIT_SECONDS,
    STANDARD_UNIT_COUNT_SECOND,
    STANDARD_UNIT_PERCENT,
    STANDARD_UNIT_COUNT,
    STANDARD_UNIT_KILOBITS,
    STANDARD_UNIT_KILOBYTES,
    STANDARD_UNIT_TERABYTES,
    STANDARD_UNIT_GIGABYTES_SECOND,
    STANDARD_UNIT_MICROSECONDS,
    STANDARD_UNIT_TERABITS,
    STANDARD_UNIT_BITS_SECOND,
    STANDARD_UNIT_TERABYTES_SECOND,
  ];

  static final $core.List<StandardUnit?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 26);
  static StandardUnit? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StandardUnit._(super.value, super.name);
}

class StateValue extends $pb.ProtobufEnum {
  static const StateValue STATE_VALUE_ALARM =
      StateValue._(0, _omitEnumNames ? '' : 'STATE_VALUE_ALARM');
  static const StateValue STATE_VALUE_OK =
      StateValue._(1, _omitEnumNames ? '' : 'STATE_VALUE_OK');
  static const StateValue STATE_VALUE_INSUFFICIENT_DATA =
      StateValue._(2, _omitEnumNames ? '' : 'STATE_VALUE_INSUFFICIENT_DATA');

  static const $core.List<StateValue> values = <StateValue>[
    STATE_VALUE_ALARM,
    STATE_VALUE_OK,
    STATE_VALUE_INSUFFICIENT_DATA,
  ];

  static final $core.List<StateValue?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static StateValue? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StateValue._(super.value, super.name);
}

class Statistic extends $pb.ProtobufEnum {
  static const Statistic STATISTIC_SUM =
      Statistic._(0, _omitEnumNames ? '' : 'STATISTIC_SUM');
  static const Statistic STATISTIC_SAMPLECOUNT =
      Statistic._(1, _omitEnumNames ? '' : 'STATISTIC_SAMPLECOUNT');
  static const Statistic STATISTIC_AVERAGE =
      Statistic._(2, _omitEnumNames ? '' : 'STATISTIC_AVERAGE');
  static const Statistic STATISTIC_MAXIMUM =
      Statistic._(3, _omitEnumNames ? '' : 'STATISTIC_MAXIMUM');
  static const Statistic STATISTIC_MINIMUM =
      Statistic._(4, _omitEnumNames ? '' : 'STATISTIC_MINIMUM');

  static const $core.List<Statistic> values = <Statistic>[
    STATISTIC_SUM,
    STATISTIC_SAMPLECOUNT,
    STATISTIC_AVERAGE,
    STATISTIC_MAXIMUM,
    STATISTIC_MINIMUM,
  ];

  static final $core.List<Statistic?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static Statistic? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Statistic._(super.value, super.name);
}

class StatusCode extends $pb.ProtobufEnum {
  static const StatusCode STATUS_CODE_PARTIAL_DATA =
      StatusCode._(0, _omitEnumNames ? '' : 'STATUS_CODE_PARTIAL_DATA');
  static const StatusCode STATUS_CODE_FORBIDDEN =
      StatusCode._(1, _omitEnumNames ? '' : 'STATUS_CODE_FORBIDDEN');
  static const StatusCode STATUS_CODE_COMPLETE =
      StatusCode._(2, _omitEnumNames ? '' : 'STATUS_CODE_COMPLETE');
  static const StatusCode STATUS_CODE_INTERNAL_ERROR =
      StatusCode._(3, _omitEnumNames ? '' : 'STATUS_CODE_INTERNAL_ERROR');

  static const $core.List<StatusCode> values = <StatusCode>[
    STATUS_CODE_PARTIAL_DATA,
    STATUS_CODE_FORBIDDEN,
    STATUS_CODE_COMPLETE,
    STATUS_CODE_INTERNAL_ERROR,
  ];

  static final $core.List<StatusCode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static StatusCode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StatusCode._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
