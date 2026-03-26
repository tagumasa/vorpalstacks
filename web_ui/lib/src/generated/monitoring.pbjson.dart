// This is a generated file - do not edit.
//
// Generated from monitoring.proto.

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

@$core.Deprecated('Use actionsSuppressedByDescriptor instead')
const ActionsSuppressedBy$json = {
  '1': 'ActionsSuppressedBy',
  '2': [
    {'1': 'ACTIONS_SUPPRESSED_BY_WAITPERIOD', '2': 0},
    {'1': 'ACTIONS_SUPPRESSED_BY_ALARM', '2': 1},
    {'1': 'ACTIONS_SUPPRESSED_BY_EXTENSIONPERIOD', '2': 2},
  ],
};

/// Descriptor for `ActionsSuppressedBy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List actionsSuppressedByDescriptor = $convert.base64Decode(
    'ChNBY3Rpb25zU3VwcHJlc3NlZEJ5EiQKIEFDVElPTlNfU1VQUFJFU1NFRF9CWV9XQUlUUEVSSU'
    '9EEAASHwobQUNUSU9OU19TVVBQUkVTU0VEX0JZX0FMQVJNEAESKQolQUNUSU9OU19TVVBQUkVT'
    'U0VEX0JZX0VYVEVOU0lPTlBFUklPRBAC');

@$core.Deprecated('Use alarmMuteRuleStatusDescriptor instead')
const AlarmMuteRuleStatus$json = {
  '1': 'AlarmMuteRuleStatus',
  '2': [
    {'1': 'ALARM_MUTE_RULE_STATUS_ACTIVE', '2': 0},
    {'1': 'ALARM_MUTE_RULE_STATUS_SCHEDULED', '2': 1},
    {'1': 'ALARM_MUTE_RULE_STATUS_EXPIRED', '2': 2},
  ],
};

/// Descriptor for `AlarmMuteRuleStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List alarmMuteRuleStatusDescriptor = $convert.base64Decode(
    'ChNBbGFybU11dGVSdWxlU3RhdHVzEiEKHUFMQVJNX01VVEVfUlVMRV9TVEFUVVNfQUNUSVZFEA'
    'ASJAogQUxBUk1fTVVURV9SVUxFX1NUQVRVU19TQ0hFRFVMRUQQARIiCh5BTEFSTV9NVVRFX1JV'
    'TEVfU1RBVFVTX0VYUElSRUQQAg==');

@$core.Deprecated('Use alarmTypeDescriptor instead')
const AlarmType$json = {
  '1': 'AlarmType',
  '2': [
    {'1': 'ALARM_TYPE_COMPOSITEALARM', '2': 0},
    {'1': 'ALARM_TYPE_METRICALARM', '2': 1},
  ],
};

/// Descriptor for `AlarmType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List alarmTypeDescriptor = $convert.base64Decode(
    'CglBbGFybVR5cGUSHQoZQUxBUk1fVFlQRV9DT01QT1NJVEVBTEFSTRAAEhoKFkFMQVJNX1RZUE'
    'VfTUVUUklDQUxBUk0QAQ==');

@$core.Deprecated('Use anomalyDetectorStateValueDescriptor instead')
const AnomalyDetectorStateValue$json = {
  '1': 'AnomalyDetectorStateValue',
  '2': [
    {'1': 'ANOMALY_DETECTOR_STATE_VALUE_TRAINED', '2': 0},
    {'1': 'ANOMALY_DETECTOR_STATE_VALUE_PENDING_TRAINING', '2': 1},
    {'1': 'ANOMALY_DETECTOR_STATE_VALUE_TRAINED_INSUFFICIENT_DATA', '2': 2},
  ],
};

/// Descriptor for `AnomalyDetectorStateValue`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List anomalyDetectorStateValueDescriptor = $convert.base64Decode(
    'ChlBbm9tYWx5RGV0ZWN0b3JTdGF0ZVZhbHVlEigKJEFOT01BTFlfREVURUNUT1JfU1RBVEVfVk'
    'FMVUVfVFJBSU5FRBAAEjEKLUFOT01BTFlfREVURUNUT1JfU1RBVEVfVkFMVUVfUEVORElOR19U'
    'UkFJTklORxABEjoKNkFOT01BTFlfREVURUNUT1JfU1RBVEVfVkFMVUVfVFJBSU5FRF9JTlNVRk'
    'ZJQ0lFTlRfREFUQRAC');

@$core.Deprecated('Use anomalyDetectorTypeDescriptor instead')
const AnomalyDetectorType$json = {
  '1': 'AnomalyDetectorType',
  '2': [
    {'1': 'ANOMALY_DETECTOR_TYPE_SINGLE_METRIC', '2': 0},
    {'1': 'ANOMALY_DETECTOR_TYPE_METRIC_MATH', '2': 1},
  ],
};

/// Descriptor for `AnomalyDetectorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List anomalyDetectorTypeDescriptor = $convert.base64Decode(
    'ChNBbm9tYWx5RGV0ZWN0b3JUeXBlEicKI0FOT01BTFlfREVURUNUT1JfVFlQRV9TSU5HTEVfTU'
    'VUUklDEAASJQohQU5PTUFMWV9ERVRFQ1RPUl9UWVBFX01FVFJJQ19NQVRIEAE=');

@$core.Deprecated('Use comparisonOperatorDescriptor instead')
const ComparisonOperator$json = {
  '1': 'ComparisonOperator',
  '2': [
    {
      '1': 'COMPARISON_OPERATOR_LESSTHANLOWERORGREATERTHANUPPERTHRESHOLD',
      '2': 0
    },
    {'1': 'COMPARISON_OPERATOR_LESSTHANLOWERTHRESHOLD', '2': 1},
    {'1': 'COMPARISON_OPERATOR_GREATERTHANOREQUALTOTHRESHOLD', '2': 2},
    {'1': 'COMPARISON_OPERATOR_GREATERTHANUPPERTHRESHOLD', '2': 3},
    {'1': 'COMPARISON_OPERATOR_GREATERTHANTHRESHOLD', '2': 4},
    {'1': 'COMPARISON_OPERATOR_LESSTHANTHRESHOLD', '2': 5},
    {'1': 'COMPARISON_OPERATOR_LESSTHANOREQUALTOTHRESHOLD', '2': 6},
  ],
};

/// Descriptor for `ComparisonOperator`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List comparisonOperatorDescriptor = $convert.base64Decode(
    'ChJDb21wYXJpc29uT3BlcmF0b3ISQAo8Q09NUEFSSVNPTl9PUEVSQVRPUl9MRVNTVEhBTkxPV0'
    'VST1JHUkVBVEVSVEhBTlVQUEVSVEhSRVNIT0xEEAASLgoqQ09NUEFSSVNPTl9PUEVSQVRPUl9M'
    'RVNTVEhBTkxPV0VSVEhSRVNIT0xEEAESNQoxQ09NUEFSSVNPTl9PUEVSQVRPUl9HUkVBVEVSVE'
    'hBTk9SRVFVQUxUT1RIUkVTSE9MRBACEjEKLUNPTVBBUklTT05fT1BFUkFUT1JfR1JFQVRFUlRI'
    'QU5VUFBFUlRIUkVTSE9MRBADEiwKKENPTVBBUklTT05fT1BFUkFUT1JfR1JFQVRFUlRIQU5USF'
    'JFU0hPTEQQBBIpCiVDT01QQVJJU09OX09QRVJBVE9SX0xFU1NUSEFOVEhSRVNIT0xEEAUSMgou'
    'Q09NUEFSSVNPTl9PUEVSQVRPUl9MRVNTVEhBTk9SRVFVQUxUT1RIUkVTSE9MRBAG');

@$core.Deprecated('Use evaluationStateDescriptor instead')
const EvaluationState$json = {
  '1': 'EvaluationState',
  '2': [
    {'1': 'EVALUATION_STATE_PARTIAL_DATA', '2': 0},
    {'1': 'EVALUATION_STATE_EVALUATION_ERROR', '2': 1},
    {'1': 'EVALUATION_STATE_EVALUATION_FAILURE', '2': 2},
  ],
};

/// Descriptor for `EvaluationState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List evaluationStateDescriptor = $convert.base64Decode(
    'Cg9FdmFsdWF0aW9uU3RhdGUSIQodRVZBTFVBVElPTl9TVEFURV9QQVJUSUFMX0RBVEEQABIlCi'
    'FFVkFMVUFUSU9OX1NUQVRFX0VWQUxVQVRJT05fRVJST1IQARInCiNFVkFMVUFUSU9OX1NUQVRF'
    'X0VWQUxVQVRJT05fRkFJTFVSRRAC');

@$core.Deprecated('Use historyItemTypeDescriptor instead')
const HistoryItemType$json = {
  '1': 'HistoryItemType',
  '2': [
    {'1': 'HISTORY_ITEM_TYPE_STATEUPDATE', '2': 0},
    {'1': 'HISTORY_ITEM_TYPE_ALARMCONTRIBUTORACTION', '2': 1},
    {'1': 'HISTORY_ITEM_TYPE_ACTION', '2': 2},
    {'1': 'HISTORY_ITEM_TYPE_ALARMCONTRIBUTORSTATEUPDATE', '2': 3},
    {'1': 'HISTORY_ITEM_TYPE_CONFIGURATIONUPDATE', '2': 4},
  ],
};

/// Descriptor for `HistoryItemType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List historyItemTypeDescriptor = $convert.base64Decode(
    'Cg9IaXN0b3J5SXRlbVR5cGUSIQodSElTVE9SWV9JVEVNX1RZUEVfU1RBVEVVUERBVEUQABIsCi'
    'hISVNUT1JZX0lURU1fVFlQRV9BTEFSTUNPTlRSSUJVVE9SQUNUSU9OEAESHAoYSElTVE9SWV9J'
    'VEVNX1RZUEVfQUNUSU9OEAISMQotSElTVE9SWV9JVEVNX1RZUEVfQUxBUk1DT05UUklCVVRPUl'
    'NUQVRFVVBEQVRFEAMSKQolSElTVE9SWV9JVEVNX1RZUEVfQ09ORklHVVJBVElPTlVQREFURRAE');

@$core.Deprecated('Use metricStreamOutputFormatDescriptor instead')
const MetricStreamOutputFormat$json = {
  '1': 'MetricStreamOutputFormat',
  '2': [
    {'1': 'METRIC_STREAM_OUTPUT_FORMAT_OPEN_TELEMETRY_0_7', '2': 0},
    {'1': 'METRIC_STREAM_OUTPUT_FORMAT_JSON', '2': 1},
    {'1': 'METRIC_STREAM_OUTPUT_FORMAT_OPEN_TELEMETRY_1_0', '2': 2},
  ],
};

/// Descriptor for `MetricStreamOutputFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricStreamOutputFormatDescriptor = $convert.base64Decode(
    'ChhNZXRyaWNTdHJlYW1PdXRwdXRGb3JtYXQSMgouTUVUUklDX1NUUkVBTV9PVVRQVVRfRk9STU'
    'FUX09QRU5fVEVMRU1FVFJZXzBfNxAAEiQKIE1FVFJJQ19TVFJFQU1fT1VUUFVUX0ZPUk1BVF9K'
    'U09OEAESMgouTUVUUklDX1NUUkVBTV9PVVRQVVRfRk9STUFUX09QRU5fVEVMRU1FVFJZXzFfMB'
    'AC');

@$core.Deprecated('Use recentlyActiveDescriptor instead')
const RecentlyActive$json = {
  '1': 'RecentlyActive',
  '2': [
    {'1': 'RECENTLY_ACTIVE_PT3H', '2': 0},
  ],
};

/// Descriptor for `RecentlyActive`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List recentlyActiveDescriptor = $convert
    .base64Decode('Cg5SZWNlbnRseUFjdGl2ZRIYChRSRUNFTlRMWV9BQ1RJVkVfUFQzSBAA');

@$core.Deprecated('Use scanByDescriptor instead')
const ScanBy$json = {
  '1': 'ScanBy',
  '2': [
    {'1': 'SCAN_BY_TIMESTAMP_DESCENDING', '2': 0},
    {'1': 'SCAN_BY_TIMESTAMP_ASCENDING', '2': 1},
  ],
};

/// Descriptor for `ScanBy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scanByDescriptor = $convert.base64Decode(
    'CgZTY2FuQnkSIAocU0NBTl9CWV9USU1FU1RBTVBfREVTQ0VORElORxAAEh8KG1NDQU5fQllfVE'
    'lNRVNUQU1QX0FTQ0VORElORxAB');

@$core.Deprecated('Use standardUnitDescriptor instead')
const StandardUnit$json = {
  '1': 'StandardUnit',
  '2': [
    {'1': 'STANDARD_UNIT_MEGABITS_SECOND', '2': 0},
    {'1': 'STANDARD_UNIT_KILOBYTES_SECOND', '2': 1},
    {'1': 'STANDARD_UNIT_KILOBITS_SECOND', '2': 2},
    {'1': 'STANDARD_UNIT_GIGABITS', '2': 3},
    {'1': 'STANDARD_UNIT_NONE', '2': 4},
    {'1': 'STANDARD_UNIT_MILLISECONDS', '2': 5},
    {'1': 'STANDARD_UNIT_TERABITS_SECOND', '2': 6},
    {'1': 'STANDARD_UNIT_BYTES', '2': 7},
    {'1': 'STANDARD_UNIT_GIGABYTES', '2': 8},
    {'1': 'STANDARD_UNIT_MEGABYTES_SECOND', '2': 9},
    {'1': 'STANDARD_UNIT_BITS', '2': 10},
    {'1': 'STANDARD_UNIT_GIGABITS_SECOND', '2': 11},
    {'1': 'STANDARD_UNIT_MEGABYTES', '2': 12},
    {'1': 'STANDARD_UNIT_BYTES_SECOND', '2': 13},
    {'1': 'STANDARD_UNIT_MEGABITS', '2': 14},
    {'1': 'STANDARD_UNIT_SECONDS', '2': 15},
    {'1': 'STANDARD_UNIT_COUNT_SECOND', '2': 16},
    {'1': 'STANDARD_UNIT_PERCENT', '2': 17},
    {'1': 'STANDARD_UNIT_COUNT', '2': 18},
    {'1': 'STANDARD_UNIT_KILOBITS', '2': 19},
    {'1': 'STANDARD_UNIT_KILOBYTES', '2': 20},
    {'1': 'STANDARD_UNIT_TERABYTES', '2': 21},
    {'1': 'STANDARD_UNIT_GIGABYTES_SECOND', '2': 22},
    {'1': 'STANDARD_UNIT_MICROSECONDS', '2': 23},
    {'1': 'STANDARD_UNIT_TERABITS', '2': 24},
    {'1': 'STANDARD_UNIT_BITS_SECOND', '2': 25},
    {'1': 'STANDARD_UNIT_TERABYTES_SECOND', '2': 26},
  ],
};

/// Descriptor for `StandardUnit`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List standardUnitDescriptor = $convert.base64Decode(
    'CgxTdGFuZGFyZFVuaXQSIQodU1RBTkRBUkRfVU5JVF9NRUdBQklUU19TRUNPTkQQABIiCh5TVE'
    'FOREFSRF9VTklUX0tJTE9CWVRFU19TRUNPTkQQARIhCh1TVEFOREFSRF9VTklUX0tJTE9CSVRT'
    'X1NFQ09ORBACEhoKFlNUQU5EQVJEX1VOSVRfR0lHQUJJVFMQAxIWChJTVEFOREFSRF9VTklUX0'
    '5PTkUQBBIeChpTVEFOREFSRF9VTklUX01JTExJU0VDT05EUxAFEiEKHVNUQU5EQVJEX1VOSVRf'
    'VEVSQUJJVFNfU0VDT05EEAYSFwoTU1RBTkRBUkRfVU5JVF9CWVRFUxAHEhsKF1NUQU5EQVJEX1'
    'VOSVRfR0lHQUJZVEVTEAgSIgoeU1RBTkRBUkRfVU5JVF9NRUdBQllURVNfU0VDT05EEAkSFgoS'
    'U1RBTkRBUkRfVU5JVF9CSVRTEAoSIQodU1RBTkRBUkRfVU5JVF9HSUdBQklUU19TRUNPTkQQCx'
    'IbChdTVEFOREFSRF9VTklUX01FR0FCWVRFUxAMEh4KGlNUQU5EQVJEX1VOSVRfQllURVNfU0VD'
    'T05EEA0SGgoWU1RBTkRBUkRfVU5JVF9NRUdBQklUUxAOEhkKFVNUQU5EQVJEX1VOSVRfU0VDT0'
    '5EUxAPEh4KGlNUQU5EQVJEX1VOSVRfQ09VTlRfU0VDT05EEBASGQoVU1RBTkRBUkRfVU5JVF9Q'
    'RVJDRU5UEBESFwoTU1RBTkRBUkRfVU5JVF9DT1VOVBASEhoKFlNUQU5EQVJEX1VOSVRfS0lMT0'
    'JJVFMQExIbChdTVEFOREFSRF9VTklUX0tJTE9CWVRFUxAUEhsKF1NUQU5EQVJEX1VOSVRfVEVS'
    'QUJZVEVTEBUSIgoeU1RBTkRBUkRfVU5JVF9HSUdBQllURVNfU0VDT05EEBYSHgoaU1RBTkRBUk'
    'RfVU5JVF9NSUNST1NFQ09ORFMQFxIaChZTVEFOREFSRF9VTklUX1RFUkFCSVRTEBgSHQoZU1RB'
    'TkRBUkRfVU5JVF9CSVRTX1NFQ09ORBAZEiIKHlNUQU5EQVJEX1VOSVRfVEVSQUJZVEVTX1NFQ0'
    '9ORBAa');

@$core.Deprecated('Use stateValueDescriptor instead')
const StateValue$json = {
  '1': 'StateValue',
  '2': [
    {'1': 'STATE_VALUE_ALARM', '2': 0},
    {'1': 'STATE_VALUE_OK', '2': 1},
    {'1': 'STATE_VALUE_INSUFFICIENT_DATA', '2': 2},
  ],
};

/// Descriptor for `StateValue`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stateValueDescriptor = $convert.base64Decode(
    'CgpTdGF0ZVZhbHVlEhUKEVNUQVRFX1ZBTFVFX0FMQVJNEAASEgoOU1RBVEVfVkFMVUVfT0sQAR'
    'IhCh1TVEFURV9WQUxVRV9JTlNVRkZJQ0lFTlRfREFUQRAC');

@$core.Deprecated('Use statisticDescriptor instead')
const Statistic$json = {
  '1': 'Statistic',
  '2': [
    {'1': 'STATISTIC_SUM', '2': 0},
    {'1': 'STATISTIC_SAMPLECOUNT', '2': 1},
    {'1': 'STATISTIC_AVERAGE', '2': 2},
    {'1': 'STATISTIC_MAXIMUM', '2': 3},
    {'1': 'STATISTIC_MINIMUM', '2': 4},
  ],
};

/// Descriptor for `Statistic`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List statisticDescriptor = $convert.base64Decode(
    'CglTdGF0aXN0aWMSEQoNU1RBVElTVElDX1NVTRAAEhkKFVNUQVRJU1RJQ19TQU1QTEVDT1VOVB'
    'ABEhUKEVNUQVRJU1RJQ19BVkVSQUdFEAISFQoRU1RBVElTVElDX01BWElNVU0QAxIVChFTVEFU'
    'SVNUSUNfTUlOSU1VTRAE');

@$core.Deprecated('Use statusCodeDescriptor instead')
const StatusCode$json = {
  '1': 'StatusCode',
  '2': [
    {'1': 'STATUS_CODE_PARTIAL_DATA', '2': 0},
    {'1': 'STATUS_CODE_FORBIDDEN', '2': 1},
    {'1': 'STATUS_CODE_COMPLETE', '2': 2},
    {'1': 'STATUS_CODE_INTERNAL_ERROR', '2': 3},
  ],
};

/// Descriptor for `StatusCode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List statusCodeDescriptor = $convert.base64Decode(
    'CgpTdGF0dXNDb2RlEhwKGFNUQVRVU19DT0RFX1BBUlRJQUxfREFUQRAAEhkKFVNUQVRVU19DT0'
    'RFX0ZPUkJJRERFThABEhgKFFNUQVRVU19DT0RFX0NPTVBMRVRFEAISHgoaU1RBVFVTX0NPREVf'
    'SU5URVJOQUxfRVJST1IQAw==');

@$core.Deprecated('Use alarmContributorDescriptor instead')
const AlarmContributor$json = {
  '1': 'AlarmContributor',
  '2': [
    {
      '1': 'contributorattributes',
      '3': 82097856,
      '4': 3,
      '5': 11,
      '6': '.monitoring.AlarmContributor.ContributorattributesEntry',
      '10': 'contributorattributes'
    },
    {
      '1': 'contributorid',
      '3': 33454376,
      '4': 1,
      '5': 9,
      '10': 'contributorid'
    },
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
    {
      '1': 'statetransitionedtimestamp',
      '3': 46811093,
      '4': 1,
      '5': 9,
      '10': 'statetransitionedtimestamp'
    },
  ],
  '3': [AlarmContributor_ContributorattributesEntry$json],
};

@$core.Deprecated('Use alarmContributorDescriptor instead')
const AlarmContributor_ContributorattributesEntry$json = {
  '1': 'ContributorattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AlarmContributor`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alarmContributorDescriptor = $convert.base64Decode(
    'ChBBbGFybUNvbnRyaWJ1dG9yEnAKFWNvbnRyaWJ1dG9yYXR0cmlidXRlcxjA7ZInIAMoCzI3Lm'
    '1vbml0b3JpbmcuQWxhcm1Db250cmlidXRvci5Db250cmlidXRvcmF0dHJpYnV0ZXNFbnRyeVIV'
    'Y29udHJpYnV0b3JhdHRyaWJ1dGVzEicKDWNvbnRyaWJ1dG9yaWQYqPL5DyABKAlSDWNvbnRyaW'
    'J1dG9yaWQSJAoLc3RhdGVyZWFzb24Y89WtswEgASgJUgtzdGF0ZXJlYXNvbhJBChpzdGF0ZXRy'
    'YW5zaXRpb25lZHRpbWVzdGFtcBjVj6kWIAEoCVIac3RhdGV0cmFuc2l0aW9uZWR0aW1lc3RhbX'
    'AaSAoaQ29udHJpYnV0b3JhdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFs'
    'dWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use alarmHistoryItemDescriptor instead')
const AlarmHistoryItem$json = {
  '1': 'AlarmHistoryItem',
  '2': [
    {
      '1': 'alarmcontributorattributes',
      '3': 34317575,
      '4': 3,
      '5': 11,
      '6': '.monitoring.AlarmHistoryItem.AlarmcontributorattributesEntry',
      '10': 'alarmcontributorattributes'
    },
    {
      '1': 'alarmcontributorid',
      '3': 305558919,
      '4': 1,
      '5': 9,
      '10': 'alarmcontributorid'
    },
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {
      '1': 'alarmtype',
      '3': 301805339,
      '4': 1,
      '5': 14,
      '6': '.monitoring.AlarmType',
      '10': 'alarmtype'
    },
    {'1': 'historydata', '3': 30118600, '4': 1, '5': 9, '10': 'historydata'},
    {
      '1': 'historyitemtype',
      '3': 442439537,
      '4': 1,
      '5': 14,
      '6': '.monitoring.HistoryItemType',
      '10': 'historyitemtype'
    },
    {
      '1': 'historysummary',
      '3': 53483226,
      '4': 1,
      '5': 9,
      '10': 'historysummary'
    },
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
  ],
  '3': [AlarmHistoryItem_AlarmcontributorattributesEntry$json],
};

@$core.Deprecated('Use alarmHistoryItemDescriptor instead')
const AlarmHistoryItem_AlarmcontributorattributesEntry$json = {
  '1': 'AlarmcontributorattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AlarmHistoryItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alarmHistoryItemDescriptor = $convert.base64Decode(
    'ChBBbGFybUhpc3RvcnlJdGVtEn8KGmFsYXJtY29udHJpYnV0b3JhdHRyaWJ1dGVzGIfKrhAgAy'
    'gLMjwubW9uaXRvcmluZy5BbGFybUhpc3RvcnlJdGVtLkFsYXJtY29udHJpYnV0b3JhdHRyaWJ1'
    'dGVzRW50cnlSGmFsYXJtY29udHJpYnV0b3JhdHRyaWJ1dGVzEjIKEmFsYXJtY29udHJpYnV0b3'
    'JpZBiH69mRASABKAlSEmFsYXJtY29udHJpYnV0b3JpZBIfCglhbGFybW5hbWUY4t/lGSABKAlS'
    'CWFsYXJtbmFtZRI3CglhbGFybXR5cGUYm970jwEgASgOMhUubW9uaXRvcmluZy5BbGFybVR5cG'
    'VSCWFsYXJtdHlwZRIjCgtoaXN0b3J5ZGF0YRjIpa4OIAEoCVILaGlzdG9yeWRhdGESSQoPaGlz'
    'dG9yeWl0ZW10eXBlGPGu/NIBIAEoDjIbLm1vbml0b3JpbmcuSGlzdG9yeUl0ZW1UeXBlUg9oaX'
    'N0b3J5aXRlbXR5cGUSKQoOaGlzdG9yeXN1bW1hcnkY2q3AGSABKAlSDmhpc3RvcnlzdW1tYXJ5'
    'Eh8KCXRpbWVzdGFtcBjEw7dNIAEoCVIJdGltZXN0YW1wGk0KH0FsYXJtY29udHJpYnV0b3JhdH'
    'RyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4'
    'AQ==');

@$core.Deprecated('Use alarmMuteRuleSummaryDescriptor instead')
const AlarmMuteRuleSummary$json = {
  '1': 'AlarmMuteRuleSummary',
  '2': [
    {
      '1': 'alarmmuterulearn',
      '3': 493590237,
      '4': 1,
      '5': 9,
      '10': 'alarmmuterulearn'
    },
    {'1': 'expiredate', '3': 454143207, '4': 1, '5': 9, '10': 'expiredate'},
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {'1': 'mutetype', '3': 111213245, '4': 1, '5': 9, '10': 'mutetype'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.monitoring.AlarmMuteRuleStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `AlarmMuteRuleSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alarmMuteRuleSummaryDescriptor = $convert.base64Decode(
    'ChRBbGFybU11dGVSdWxlU3VtbWFyeRIuChBhbGFybW11dGVydWxlYXJuGN2trusBIAEoCVIQYW'
    'xhcm1tdXRlcnVsZWFybhIiCgpleHBpcmVkYXRlGOfZxtgBIAEoCVIKZXhwaXJlZGF0ZRI1ChRs'
    'YXN0dXBkYXRlZHRpbWVzdGFtcBiVy8g/IAEoCVIUbGFzdHVwZGF0ZWR0aW1lc3RhbXASHQoIbX'
    'V0ZXR5cGUYvfWDNSABKAlSCG11dGV0eXBlEjoKBnN0YXR1cxiQ5PsCIAEoDjIfLm1vbml0b3Jp'
    'bmcuQWxhcm1NdXRlUnVsZVN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use anomalyDetectorDescriptor instead')
const AnomalyDetector$json = {
  '1': 'AnomalyDetector',
  '2': [
    {
      '1': 'configuration',
      '3': 442426458,
      '4': 1,
      '5': 11,
      '6': '.monitoring.AnomalyDetectorConfiguration',
      '10': 'configuration'
    },
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {
      '1': 'metriccharacteristics',
      '3': 129404842,
      '4': 1,
      '5': 11,
      '6': '.monitoring.MetricCharacteristics',
      '10': 'metriccharacteristics'
    },
    {
      '1': 'metricmathanomalydetector',
      '3': 439381475,
      '4': 1,
      '5': 11,
      '6': '.monitoring.MetricMathAnomalyDetector',
      '10': 'metricmathanomalydetector'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {
      '1': 'singlemetricanomalydetector',
      '3': 418531933,
      '4': 1,
      '5': 11,
      '6': '.monitoring.SingleMetricAnomalyDetector',
      '10': 'singlemetricanomalydetector'
    },
    {'1': 'stat', '3': 325549752, '4': 1, '5': 9, '10': 'stat'},
    {
      '1': 'statevalue',
      '3': 334526008,
      '4': 1,
      '5': 14,
      '6': '.monitoring.AnomalyDetectorStateValue',
      '10': 'statevalue'
    },
  ],
};

/// Descriptor for `AnomalyDetector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anomalyDetectorDescriptor = $convert.base64Decode(
    'Cg9Bbm9tYWx5RGV0ZWN0b3ISUgoNY29uZmlndXJhdGlvbhjayPvSASABKAsyKC5tb25pdG9yaW'
    '5nLkFub21hbHlEZXRlY3RvckNvbmZpZ3VyYXRpb25SDWNvbmZpZ3VyYXRpb24SOQoKZGltZW5z'
    'aW9ucxjRm9/cASADKAsyFS5tb25pdG9yaW5nLkRpbWVuc2lvblIKZGltZW5zaW9ucxJaChVtZX'
    'RyaWNjaGFyYWN0ZXJpc3RpY3MYqp/aPSABKAsyIS5tb25pdG9yaW5nLk1ldHJpY0NoYXJhY3Rl'
    'cmlzdGljc1IVbWV0cmljY2hhcmFjdGVyaXN0aWNzEmcKGW1ldHJpY21hdGhhbm9tYWx5ZGV0ZW'
    'N0b3IY49vB0QEgASgLMiUubW9uaXRvcmluZy5NZXRyaWNNYXRoQW5vbWFseURldGVjdG9yUhlt'
    'ZXRyaWNtYXRoYW5vbWFseWRldGVjdG9yEiEKCm1ldHJpY25hbWUY+77aMiABKAlSCm1ldHJpY2'
    '5hbWUSIAoJbmFtZXNwYWNlGMGEuakBIAEoCVIJbmFtZXNwYWNlEm0KG3NpbmdsZW1ldHJpY2Fu'
    'b21hbHlkZXRlY3RvchjdlMnHASABKAsyJy5tb25pdG9yaW5nLlNpbmdsZU1ldHJpY0Fub21hbH'
    'lEZXRlY3RvclIbc2luZ2xlbWV0cmljYW5vbWFseWRldGVjdG9yEhYKBHN0YXQYuP2dmwEgASgJ'
    'UgRzdGF0EkkKCnN0YXRldmFsdWUYuOzBnwEgASgOMiUubW9uaXRvcmluZy5Bbm9tYWx5RGV0ZW'
    'N0b3JTdGF0ZVZhbHVlUgpzdGF0ZXZhbHVl');

@$core.Deprecated('Use anomalyDetectorConfigurationDescriptor instead')
const AnomalyDetectorConfiguration$json = {
  '1': 'AnomalyDetectorConfiguration',
  '2': [
    {
      '1': 'excludedtimeranges',
      '3': 31303297,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Range',
      '10': 'excludedtimeranges'
    },
    {
      '1': 'metrictimezone',
      '3': 528921359,
      '4': 1,
      '5': 9,
      '10': 'metrictimezone'
    },
  ],
};

/// Descriptor for `AnomalyDetectorConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anomalyDetectorConfigurationDescriptor =
    $convert.base64Decode(
        'ChxBbm9tYWx5RGV0ZWN0b3JDb25maWd1cmF0aW9uEkQKEmV4Y2x1ZGVkdGltZXJhbmdlcxiBzf'
        'YOIAMoCzIRLm1vbml0b3JpbmcuUmFuZ2VSEmV4Y2x1ZGVkdGltZXJhbmdlcxIqCg5tZXRyaWN0'
        'aW1lem9uZRiP5pr8ASABKAlSDm1ldHJpY3RpbWV6b25l');

@$core.Deprecated('Use compositeAlarmDescriptor instead')
const CompositeAlarm$json = {
  '1': 'CompositeAlarm',
  '2': [
    {
      '1': 'actionsenabled',
      '3': 126027878,
      '4': 1,
      '5': 8,
      '10': 'actionsenabled'
    },
    {
      '1': 'actionssuppressedby',
      '3': 377750064,
      '4': 1,
      '5': 14,
      '6': '.monitoring.ActionsSuppressedBy',
      '10': 'actionssuppressedby'
    },
    {
      '1': 'actionssuppressedreason',
      '3': 20815779,
      '4': 1,
      '5': 9,
      '10': 'actionssuppressedreason'
    },
    {
      '1': 'actionssuppressor',
      '3': 486417535,
      '4': 1,
      '5': 9,
      '10': 'actionssuppressor'
    },
    {
      '1': 'actionssuppressorextensionperiod',
      '3': 183723543,
      '4': 1,
      '5': 5,
      '10': 'actionssuppressorextensionperiod'
    },
    {
      '1': 'actionssuppressorwaitperiod',
      '3': 80227275,
      '4': 1,
      '5': 5,
      '10': 'actionssuppressorwaitperiod'
    },
    {'1': 'alarmactions', '3': 355779506, '4': 3, '5': 9, '10': 'alarmactions'},
    {'1': 'alarmarn', '3': 276019462, '4': 1, '5': 9, '10': 'alarmarn'},
    {
      '1': 'alarmconfigurationupdatedtimestamp',
      '3': 51935494,
      '4': 1,
      '5': 9,
      '10': 'alarmconfigurationupdatedtimestamp'
    },
    {
      '1': 'alarmdescription',
      '3': 30583613,
      '4': 1,
      '5': 9,
      '10': 'alarmdescription'
    },
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {'1': 'alarmrule', '3': 516996973, '4': 1, '5': 9, '10': 'alarmrule'},
    {
      '1': 'insufficientdataactions',
      '3': 498450778,
      '4': 3,
      '5': 9,
      '10': 'insufficientdataactions'
    },
    {'1': 'okactions', '3': 377443763, '4': 3, '5': 9, '10': 'okactions'},
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
    {
      '1': 'statereasondata',
      '3': 262771075,
      '4': 1,
      '5': 9,
      '10': 'statereasondata'
    },
    {
      '1': 'statetransitionedtimestamp',
      '3': 46811093,
      '4': 1,
      '5': 9,
      '10': 'statetransitionedtimestamp'
    },
    {
      '1': 'stateupdatedtimestamp',
      '3': 375406848,
      '4': 1,
      '5': 9,
      '10': 'stateupdatedtimestamp'
    },
    {
      '1': 'statevalue',
      '3': 334526008,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StateValue',
      '10': 'statevalue'
    },
  ],
};

/// Descriptor for `CompositeAlarm`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List compositeAlarmDescriptor = $convert.base64Decode(
    'Cg5Db21wb3NpdGVBbGFybRIpCg5hY3Rpb25zZW5hYmxlZBjmkIw8IAEoCFIOYWN0aW9uc2VuYW'
    'JsZWQSVQoTYWN0aW9uc3N1cHByZXNzZWRieRiwhJC0ASABKA4yHy5tb25pdG9yaW5nLkFjdGlv'
    'bnNTdXBwcmVzc2VkQnlSE2FjdGlvbnNzdXBwcmVzc2VkYnkSOwoXYWN0aW9uc3N1cHByZXNzZW'
    'RyZWFzb24Yo7/2CSABKAlSF2FjdGlvbnNzdXBwcmVzc2VkcmVhc29uEjAKEWFjdGlvbnNzdXBw'
    'cmVzc29yGP/I+OcBIAEoCVIRYWN0aW9uc3N1cHByZXNzb3ISTQogYWN0aW9uc3N1cHByZXNzb3'
    'JleHRlbnNpb25wZXJpb2QYl8zNVyABKAVSIGFjdGlvbnNzdXBwcmVzc29yZXh0ZW5zaW9ucGVy'
    'aW9kEkMKG2FjdGlvbnNzdXBwcmVzc29yd2FpdHBlcmlvZBjL16AmIAEoBVIbYWN0aW9uc3N1cH'
    'ByZXNzb3J3YWl0cGVyaW9kEiYKDGFsYXJtYWN0aW9ucxiyh9OpASADKAlSDGFsYXJtYWN0aW9u'
    'cxIeCghhbGFybWFybhiG8s6DASABKAlSCGFsYXJtYXJuElEKImFsYXJtY29uZmlndXJhdGlvbn'
    'VwZGF0ZWR0aW1lc3RhbXAYhvLhGCABKAlSImFsYXJtY29uZmlndXJhdGlvbnVwZGF0ZWR0aW1l'
    'c3RhbXASLQoQYWxhcm1kZXNjcmlwdGlvbhi91soOIAEoCVIQYWxhcm1kZXNjcmlwdGlvbhIfCg'
    'lhbGFybW5hbWUY4t/lGSABKAlSCWFsYXJtbmFtZRIgCglhbGFybXJ1bGUY7f7C9gEgASgJUglh'
    'bGFybXJ1bGUSPAoXaW5zdWZmaWNpZW50ZGF0YWFjdGlvbnMY2oLX7QEgAygJUhdpbnN1ZmZpY2'
    'llbnRkYXRhYWN0aW9ucxIgCglva2FjdGlvbnMYs6v9swEgAygJUglva2FjdGlvbnMSJAoLc3Rh'
    'dGVyZWFzb24Y89WtswEgASgJUgtzdGF0ZXJlYXNvbhIrCg9zdGF0ZXJlYXNvbmRhdGEYg6OmfS'
    'ABKAlSD3N0YXRlcmVhc29uZGF0YRJBChpzdGF0ZXRyYW5zaXRpb25lZHRpbWVzdGFtcBjVj6kW'
    'IAEoCVIac3RhdGV0cmFuc2l0aW9uZWR0aW1lc3RhbXASOAoVc3RhdGV1cGRhdGVkdGltZXN0YW'
    '1wGICCgbMBIAEoCVIVc3RhdGV1cGRhdGVkdGltZXN0YW1wEjoKCnN0YXRldmFsdWUYuOzBnwEg'
    'ASgOMhYubW9uaXRvcmluZy5TdGF0ZVZhbHVlUgpzdGF0ZXZhbHVl');

@$core.Deprecated('Use concurrentModificationExceptionDescriptor instead')
const ConcurrentModificationException$json = {
  '1': 'ConcurrentModificationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConcurrentModificationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List concurrentModificationExceptionDescriptor =
    $convert.base64Decode(
        'Ch9Db25jdXJyZW50TW9kaWZpY2F0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use conflictExceptionDescriptor instead')
const ConflictException$json = {
  '1': 'ConflictException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictExceptionDescriptor = $convert.base64Decode(
    'ChFDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use dashboardEntryDescriptor instead')
const DashboardEntry$json = {
  '1': 'DashboardEntry',
  '2': [
    {'1': 'dashboardarn', '3': 108051951, '4': 1, '5': 9, '10': 'dashboardarn'},
    {
      '1': 'dashboardname',
      '3': 506599873,
      '4': 1,
      '5': 9,
      '10': 'dashboardname'
    },
    {'1': 'lastmodified', '3': 434048551, '4': 1, '5': 9, '10': 'lastmodified'},
    {'1': 'size', '3': 105352829, '4': 1, '5': 3, '10': 'size'},
  ],
};

/// Descriptor for `DashboardEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dashboardEntryDescriptor = $convert.base64Decode(
    'Cg5EYXNoYm9hcmRFbnRyeRIlCgxkYXNoYm9hcmRhcm4Y7/vCMyABKAlSDGRhc2hib2FyZGFybh'
    'IoCg1kYXNoYm9hcmRuYW1lGMGzyPEBIAEoCVINZGFzaGJvYXJkbmFtZRImCgxsYXN0bW9kaWZp'
    'ZWQYp5z8zgEgASgJUgxsYXN0bW9kaWZpZWQSFQoEc2l6ZRj9nJ4yIAEoA1IEc2l6ZQ==');

@$core.Deprecated('Use dashboardInvalidInputErrorDescriptor instead')
const DashboardInvalidInputError$json = {
  '1': 'DashboardInvalidInputError',
  '2': [
    {
      '1': 'dashboardvalidationmessages',
      '3': 36455465,
      '4': 3,
      '5': 11,
      '6': '.monitoring.DashboardValidationMessage',
      '10': 'dashboardvalidationmessages'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DashboardInvalidInputError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dashboardInvalidInputErrorDescriptor = $convert.base64Decode(
    'ChpEYXNoYm9hcmRJbnZhbGlkSW5wdXRFcnJvchJrChtkYXNoYm9hcmR2YWxpZGF0aW9ubWVzc2'
    'FnZXMYqYixESADKAsyJi5tb25pdG9yaW5nLkRhc2hib2FyZFZhbGlkYXRpb25NZXNzYWdlUhtk'
    'YXNoYm9hcmR2YWxpZGF0aW9ubWVzc2FnZXMSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
    '==');

@$core.Deprecated('Use dashboardNotFoundErrorDescriptor instead')
const DashboardNotFoundError$json = {
  '1': 'DashboardNotFoundError',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DashboardNotFoundError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dashboardNotFoundErrorDescriptor =
    $convert.base64Decode(
        'ChZEYXNoYm9hcmROb3RGb3VuZEVycm9yEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use dashboardValidationMessageDescriptor instead')
const DashboardValidationMessage$json = {
  '1': 'DashboardValidationMessage',
  '2': [
    {'1': 'datapath', '3': 505315415, '4': 1, '5': 9, '10': 'datapath'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DashboardValidationMessage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dashboardValidationMessageDescriptor =
    $convert.base64Decode(
        'ChpEYXNoYm9hcmRWYWxpZGF0aW9uTWVzc2FnZRIeCghkYXRhcGF0aBjXgPrwASABKAlSCGRhdG'
        'FwYXRoEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use datapointDescriptor instead')
const Datapoint$json = {
  '1': 'Datapoint',
  '2': [
    {'1': 'average', '3': 134712777, '4': 1, '5': 1, '10': 'average'},
    {
      '1': 'extendedstatistics',
      '3': 365818700,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Datapoint.ExtendedstatisticsEntry',
      '10': 'extendedstatistics'
    },
    {'1': 'maximum', '3': 43343394, '4': 1, '5': 1, '10': 'maximum'},
    {'1': 'minimum', '3': 438536856, '4': 1, '5': 1, '10': 'minimum'},
    {'1': 'samplecount', '3': 384919839, '4': 1, '5': 1, '10': 'samplecount'},
    {'1': 'sum', '3': 534303305, '4': 1, '5': 1, '10': 'sum'},
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
    {
      '1': 'unit',
      '3': 148989480,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StandardUnit',
      '10': 'unit'
    },
  ],
  '3': [Datapoint_ExtendedstatisticsEntry$json],
};

@$core.Deprecated('Use datapointDescriptor instead')
const Datapoint_ExtendedstatisticsEntry$json = {
  '1': 'ExtendedstatisticsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Datapoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List datapointDescriptor = $convert.base64Decode(
    'CglEYXRhcG9pbnQSGwoHYXZlcmFnZRjJm55AIAEoAVIHYXZlcmFnZRJhChJleHRlbmRlZHN0YX'
    'Rpc3RpY3MYzOa3rgEgAygLMi0ubW9uaXRvcmluZy5EYXRhcG9pbnQuRXh0ZW5kZWRzdGF0aXN0'
    'aWNzRW50cnlSEmV4dGVuZGVkc3RhdGlzdGljcxIbCgdtYXhpbXVtGKK81RQgASgBUgdtYXhpbX'
    'VtEhwKB21pbmltdW0YmJWO0QEgASgBUgdtaW5pbXVtEiQKC3NhbXBsZWNvdW50GJ/SxbcBIAEo'
    'AVILc2FtcGxlY291bnQSFAoDc3VtGMmk4/4BIAEoAVIDc3VtEh8KCXRpbWVzdGFtcBjEw7dNIA'
    'EoCVIJdGltZXN0YW1wEi8KBHVuaXQYqMyFRyABKA4yGC5tb25pdG9yaW5nLlN0YW5kYXJkVW5p'
    'dFIEdW5pdBpFChdFeHRlbmRlZHN0YXRpc3RpY3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCg'
    'V2YWx1ZRgCIAEoAVIFdmFsdWU6AjgB');

@$core.Deprecated('Use deleteAlarmMuteRuleInputDescriptor instead')
const DeleteAlarmMuteRuleInput$json = {
  '1': 'DeleteAlarmMuteRuleInput',
  '2': [
    {
      '1': 'alarmmuterulename',
      '3': 530877511,
      '4': 1,
      '5': 9,
      '10': 'alarmmuterulename'
    },
  ],
};

/// Descriptor for `DeleteAlarmMuteRuleInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAlarmMuteRuleInputDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVBbGFybU11dGVSdWxlSW5wdXQSMAoRYWxhcm1tdXRlcnVsZW5hbWUYx5iS/QEgAS'
        'gJUhFhbGFybW11dGVydWxlbmFtZQ==');

@$core.Deprecated('Use deleteAlarmsInputDescriptor instead')
const DeleteAlarmsInput$json = {
  '1': 'DeleteAlarmsInput',
  '2': [
    {'1': 'alarmnames', '3': 90874583, '4': 3, '5': 9, '10': 'alarmnames'},
  ],
};

/// Descriptor for `DeleteAlarmsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAlarmsInputDescriptor = $convert.base64Decode(
    'ChFEZWxldGVBbGFybXNJbnB1dBIhCgphbGFybW5hbWVzGNfFqisgAygJUgphbGFybW5hbWVz');

@$core.Deprecated('Use deleteAnomalyDetectorInputDescriptor instead')
const DeleteAnomalyDetectorInput$json = {
  '1': 'DeleteAnomalyDetectorInput',
  '2': [
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {
      '1': 'metricmathanomalydetector',
      '3': 439381475,
      '4': 1,
      '5': 11,
      '6': '.monitoring.MetricMathAnomalyDetector',
      '10': 'metricmathanomalydetector'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {
      '1': 'singlemetricanomalydetector',
      '3': 418531933,
      '4': 1,
      '5': 11,
      '6': '.monitoring.SingleMetricAnomalyDetector',
      '10': 'singlemetricanomalydetector'
    },
    {'1': 'stat', '3': 325549752, '4': 1, '5': 9, '10': 'stat'},
  ],
};

/// Descriptor for `DeleteAnomalyDetectorInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAnomalyDetectorInputDescriptor = $convert.base64Decode(
    'ChpEZWxldGVBbm9tYWx5RGV0ZWN0b3JJbnB1dBI5CgpkaW1lbnNpb25zGNGb39wBIAMoCzIVLm'
    '1vbml0b3JpbmcuRGltZW5zaW9uUgpkaW1lbnNpb25zEmcKGW1ldHJpY21hdGhhbm9tYWx5ZGV0'
    'ZWN0b3IY49vB0QEgASgLMiUubW9uaXRvcmluZy5NZXRyaWNNYXRoQW5vbWFseURldGVjdG9yUh'
    'ltZXRyaWNtYXRoYW5vbWFseWRldGVjdG9yEiEKCm1ldHJpY25hbWUY+77aMiABKAlSCm1ldHJp'
    'Y25hbWUSIAoJbmFtZXNwYWNlGMGEuakBIAEoCVIJbmFtZXNwYWNlEm0KG3NpbmdsZW1ldHJpY2'
    'Fub21hbHlkZXRlY3RvchjdlMnHASABKAsyJy5tb25pdG9yaW5nLlNpbmdsZU1ldHJpY0Fub21h'
    'bHlEZXRlY3RvclIbc2luZ2xlbWV0cmljYW5vbWFseWRldGVjdG9yEhYKBHN0YXQYuP2dmwEgAS'
    'gJUgRzdGF0');

@$core.Deprecated('Use deleteAnomalyDetectorOutputDescriptor instead')
const DeleteAnomalyDetectorOutput$json = {
  '1': 'DeleteAnomalyDetectorOutput',
};

/// Descriptor for `DeleteAnomalyDetectorOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAnomalyDetectorOutputDescriptor =
    $convert.base64Decode('ChtEZWxldGVBbm9tYWx5RGV0ZWN0b3JPdXRwdXQ=');

@$core.Deprecated('Use deleteDashboardsInputDescriptor instead')
const DeleteDashboardsInput$json = {
  '1': 'DeleteDashboardsInput',
  '2': [
    {
      '1': 'dashboardnames',
      '3': 467563722,
      '4': 3,
      '5': 9,
      '10': 'dashboardnames'
    },
  ],
};

/// Descriptor for `DeleteDashboardsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDashboardsInputDescriptor = $convert.base64Decode(
    'ChVEZWxldGVEYXNoYm9hcmRzSW5wdXQSKgoOZGFzaGJvYXJkbmFtZXMYyun53gEgAygJUg5kYX'
    'NoYm9hcmRuYW1lcw==');

@$core.Deprecated('Use deleteDashboardsOutputDescriptor instead')
const DeleteDashboardsOutput$json = {
  '1': 'DeleteDashboardsOutput',
};

/// Descriptor for `DeleteDashboardsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDashboardsOutputDescriptor =
    $convert.base64Decode('ChZEZWxldGVEYXNoYm9hcmRzT3V0cHV0');

@$core.Deprecated('Use deleteInsightRulesInputDescriptor instead')
const DeleteInsightRulesInput$json = {
  '1': 'DeleteInsightRulesInput',
  '2': [
    {'1': 'rulenames', '3': 267949170, '4': 3, '5': 9, '10': 'rulenames'},
  ],
};

/// Descriptor for `DeleteInsightRulesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteInsightRulesInputDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVJbnNpZ2h0UnVsZXNJbnB1dBIfCglydWxlbmFtZXMY8qjifyADKAlSCXJ1bGVuYW'
        '1lcw==');

@$core.Deprecated('Use deleteInsightRulesOutputDescriptor instead')
const DeleteInsightRulesOutput$json = {
  '1': 'DeleteInsightRulesOutput',
  '2': [
    {
      '1': 'failures',
      '3': 335467271,
      '4': 3,
      '5': 11,
      '6': '.monitoring.PartialFailure',
      '10': 'failures'
    },
  ],
};

/// Descriptor for `DeleteInsightRulesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteInsightRulesOutputDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVJbnNpZ2h0UnVsZXNPdXRwdXQSOgoIZmFpbHVyZXMYh6b7nwEgAygLMhoubW9uaX'
        'RvcmluZy5QYXJ0aWFsRmFpbHVyZVIIZmFpbHVyZXM=');

@$core.Deprecated('Use deleteMetricStreamInputDescriptor instead')
const DeleteMetricStreamInput$json = {
  '1': 'DeleteMetricStreamInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteMetricStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMetricStreamInputDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVNZXRyaWNTdHJlYW1JbnB1dBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use deleteMetricStreamOutputDescriptor instead')
const DeleteMetricStreamOutput$json = {
  '1': 'DeleteMetricStreamOutput',
};

/// Descriptor for `DeleteMetricStreamOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMetricStreamOutputDescriptor =
    $convert.base64Decode('ChhEZWxldGVNZXRyaWNTdHJlYW1PdXRwdXQ=');

@$core.Deprecated('Use describeAlarmContributorsInputDescriptor instead')
const DescribeAlarmContributorsInput$json = {
  '1': 'DescribeAlarmContributorsInput',
  '2': [
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAlarmContributorsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAlarmContributorsInputDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZUFsYXJtQ29udHJpYnV0b3JzSW5wdXQSHwoJYWxhcm1uYW1lGOLf5RkgASgJUg'
        'lhbGFybW5hbWUSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeAlarmContributorsOutputDescriptor instead')
const DescribeAlarmContributorsOutput$json = {
  '1': 'DescribeAlarmContributorsOutput',
  '2': [
    {
      '1': 'alarmcontributors',
      '3': 265859209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.AlarmContributor',
      '10': 'alarmcontributors'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAlarmContributorsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAlarmContributorsOutputDescriptor =
    $convert.base64Decode(
        'Ch9EZXNjcmliZUFsYXJtQ29udHJpYnV0b3JzT3V0cHV0Ek0KEWFsYXJtY29udHJpYnV0b3JzGI'
        'nh4n4gAygLMhwubW9uaXRvcmluZy5BbGFybUNvbnRyaWJ1dG9yUhFhbGFybWNvbnRyaWJ1dG9y'
        'cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeAlarmHistoryInputDescriptor instead')
const DescribeAlarmHistoryInput$json = {
  '1': 'DescribeAlarmHistoryInput',
  '2': [
    {
      '1': 'alarmcontributorid',
      '3': 305558919,
      '4': 1,
      '5': 9,
      '10': 'alarmcontributorid'
    },
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {
      '1': 'alarmtypes',
      '3': 445751884,
      '4': 3,
      '5': 14,
      '6': '.monitoring.AlarmType',
      '10': 'alarmtypes'
    },
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {
      '1': 'historyitemtype',
      '3': 442439537,
      '4': 1,
      '5': 14,
      '6': '.monitoring.HistoryItemType',
      '10': 'historyitemtype'
    },
    {'1': 'maxrecords', '3': 220314370, '4': 1, '5': 5, '10': 'maxrecords'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'scanby',
      '3': 344142854,
      '4': 1,
      '5': 14,
      '6': '.monitoring.ScanBy',
      '10': 'scanby'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
  ],
};

/// Descriptor for `DescribeAlarmHistoryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAlarmHistoryInputDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZUFsYXJtSGlzdG9yeUlucHV0EjIKEmFsYXJtY29udHJpYnV0b3JpZBiH69mRAS'
    'ABKAlSEmFsYXJtY29udHJpYnV0b3JpZBIfCglhbGFybW5hbWUY4t/lGSABKAlSCWFsYXJtbmFt'
    'ZRI5CgphbGFybXR5cGVzGMzExtQBIAMoDjIVLm1vbml0b3JpbmcuQWxhcm1UeXBlUgphbGFybX'
    'R5cGVzEhsKB2VuZGRhdGUYz7P5JCABKAlSB2VuZGRhdGUSSQoPaGlzdG9yeWl0ZW10eXBlGPGu'
    '/NIBIAEoDjIbLm1vbml0b3JpbmcuSGlzdG9yeUl0ZW1UeXBlUg9oaXN0b3J5aXRlbXR5cGUSIQ'
    'oKbWF4cmVjb3JkcxiC9oZpIAEoBVIKbWF4cmVjb3JkcxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlS'
    'CW5leHR0b2tlbhIuCgZzY2FuYnkYhuiMpAEgASgOMhIubW9uaXRvcmluZy5TY2FuQnlSBnNjYW'
    '5ieRIgCglzdGFydGRhdGUY/Pig1AEgASgJUglzdGFydGRhdGU=');

@$core.Deprecated('Use describeAlarmHistoryOutputDescriptor instead')
const DescribeAlarmHistoryOutput$json = {
  '1': 'DescribeAlarmHistoryOutput',
  '2': [
    {
      '1': 'alarmhistoryitems',
      '3': 152760951,
      '4': 3,
      '5': 11,
      '6': '.monitoring.AlarmHistoryItem',
      '10': 'alarmhistoryitems'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAlarmHistoryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAlarmHistoryOutputDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZUFsYXJtSGlzdG9yeU91dHB1dBJNChFhbGFybWhpc3RvcnlpdGVtcxj35OtIIA'
        'MoCzIcLm1vbml0b3JpbmcuQWxhcm1IaXN0b3J5SXRlbVIRYWxhcm1oaXN0b3J5aXRlbXMSHwoJ'
        'bmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeAlarmsForMetricInputDescriptor instead')
const DescribeAlarmsForMetricInput$json = {
  '1': 'DescribeAlarmsForMetricInput',
  '2': [
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {
      '1': 'extendedstatistic',
      '3': 285620763,
      '4': 1,
      '5': 9,
      '10': 'extendedstatistic'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {
      '1': 'statistic',
      '3': 67293470,
      '4': 1,
      '5': 14,
      '6': '.monitoring.Statistic',
      '10': 'statistic'
    },
    {
      '1': 'unit',
      '3': 148989480,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StandardUnit',
      '10': 'unit'
    },
  ],
};

/// Descriptor for `DescribeAlarmsForMetricInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAlarmsForMetricInputDescriptor = $convert.base64Decode(
    'ChxEZXNjcmliZUFsYXJtc0Zvck1ldHJpY0lucHV0EjkKCmRpbWVuc2lvbnMY0Zvf3AEgAygLMh'
    'UubW9uaXRvcmluZy5EaW1lbnNpb25SCmRpbWVuc2lvbnMSMAoRZXh0ZW5kZWRzdGF0aXN0aWMY'
    'm/SYiAEgASgJUhFleHRlbmRlZHN0YXRpc3RpYxIhCgptZXRyaWNuYW1lGPu+2jIgASgJUgptZX'
    'RyaWNuYW1lEiAKCW5hbWVzcGFjZRjBhLmpASABKAlSCW5hbWVzcGFjZRIZCgZwZXJpb2QYpYiS'
    'OSABKAVSBnBlcmlvZBI2CglzdGF0aXN0aWMYnqKLICABKA4yFS5tb25pdG9yaW5nLlN0YXRpc3'
    'RpY1IJc3RhdGlzdGljEi8KBHVuaXQYqMyFRyABKA4yGC5tb25pdG9yaW5nLlN0YW5kYXJkVW5p'
    'dFIEdW5pdA==');

@$core.Deprecated('Use describeAlarmsForMetricOutputDescriptor instead')
const DescribeAlarmsForMetricOutput$json = {
  '1': 'DescribeAlarmsForMetricOutput',
  '2': [
    {
      '1': 'metricalarms',
      '3': 117806408,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricAlarm',
      '10': 'metricalarms'
    },
  ],
};

/// Descriptor for `DescribeAlarmsForMetricOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAlarmsForMetricOutputDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZUFsYXJtc0Zvck1ldHJpY091dHB1dBI+CgxtZXRyaWNhbGFybXMYyKqWOCADKA'
        'syFy5tb25pdG9yaW5nLk1ldHJpY0FsYXJtUgxtZXRyaWNhbGFybXM=');

@$core.Deprecated('Use describeAlarmsInputDescriptor instead')
const DescribeAlarmsInput$json = {
  '1': 'DescribeAlarmsInput',
  '2': [
    {'1': 'actionprefix', '3': 67679956, '4': 1, '5': 9, '10': 'actionprefix'},
    {
      '1': 'alarmnameprefix',
      '3': 96412406,
      '4': 1,
      '5': 9,
      '10': 'alarmnameprefix'
    },
    {'1': 'alarmnames', '3': 90874583, '4': 3, '5': 9, '10': 'alarmnames'},
    {
      '1': 'alarmtypes',
      '3': 445751884,
      '4': 3,
      '5': 14,
      '6': '.monitoring.AlarmType',
      '10': 'alarmtypes'
    },
    {
      '1': 'childrenofalarmname',
      '3': 204127700,
      '4': 1,
      '5': 9,
      '10': 'childrenofalarmname'
    },
    {'1': 'maxrecords', '3': 220314370, '4': 1, '5': 5, '10': 'maxrecords'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'parentsofalarmname',
      '3': 512433260,
      '4': 1,
      '5': 9,
      '10': 'parentsofalarmname'
    },
    {
      '1': 'statevalue',
      '3': 334526008,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StateValue',
      '10': 'statevalue'
    },
  ],
};

/// Descriptor for `DescribeAlarmsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAlarmsInputDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZUFsYXJtc0lucHV0EiUKDGFjdGlvbnByZWZpeBjU7aIgIAEoCVIMYWN0aW9ucH'
    'JlZml4EisKD2FsYXJtbmFtZXByZWZpeBj2xfwtIAEoCVIPYWxhcm1uYW1lcHJlZml4EiEKCmFs'
    'YXJtbmFtZXMY18WqKyADKAlSCmFsYXJtbmFtZXMSOQoKYWxhcm10eXBlcxjMxMbUASADKA4yFS'
    '5tb25pdG9yaW5nLkFsYXJtVHlwZVIKYWxhcm10eXBlcxIzChNjaGlsZHJlbm9mYWxhcm1uYW1l'
    'GNT7qmEgASgJUhNjaGlsZHJlbm9mYWxhcm1uYW1lEiEKCm1heHJlY29yZHMYgvaGaSABKAVSCm'
    '1heHJlY29yZHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SMgoScGFyZW50c29m'
    'YWxhcm1uYW1lGOy4rPQBIAEoCVIScGFyZW50c29mYWxhcm1uYW1lEjoKCnN0YXRldmFsdWUYuO'
    'zBnwEgASgOMhYubW9uaXRvcmluZy5TdGF0ZVZhbHVlUgpzdGF0ZXZhbHVl');

@$core.Deprecated('Use describeAlarmsOutputDescriptor instead')
const DescribeAlarmsOutput$json = {
  '1': 'DescribeAlarmsOutput',
  '2': [
    {
      '1': 'compositealarms',
      '3': 214076095,
      '4': 3,
      '5': 11,
      '6': '.monitoring.CompositeAlarm',
      '10': 'compositealarms'
    },
    {
      '1': 'metricalarms',
      '3': 117806408,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricAlarm',
      '10': 'metricalarms'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAlarmsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAlarmsOutputDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZUFsYXJtc091dHB1dBJHCg9jb21wb3NpdGVhbGFybXMYv5WKZiADKAsyGi5tb2'
    '5pdG9yaW5nLkNvbXBvc2l0ZUFsYXJtUg9jb21wb3NpdGVhbGFybXMSPgoMbWV0cmljYWxhcm1z'
    'GMiqljggAygLMhcubW9uaXRvcmluZy5NZXRyaWNBbGFybVIMbWV0cmljYWxhcm1zEh8KCW5leH'
    'R0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeAnomalyDetectorsInputDescriptor instead')
const DescribeAnomalyDetectorsInput$json = {
  '1': 'DescribeAnomalyDetectorsInput',
  '2': [
    {
      '1': 'anomalydetectortypes',
      '3': 274334014,
      '4': 3,
      '5': 14,
      '6': '.monitoring.AnomalyDetectorType',
      '10': 'anomalydetectortypes'
    },
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAnomalyDetectorsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAnomalyDetectorsInputDescriptor = $convert.base64Decode(
    'Ch1EZXNjcmliZUFub21hbHlEZXRlY3RvcnNJbnB1dBJXChRhbm9tYWx5ZGV0ZWN0b3J0eXBlcx'
    'i+guiCASADKA4yHy5tb25pdG9yaW5nLkFub21hbHlEZXRlY3RvclR5cGVSFGFub21hbHlkZXRl'
    'Y3RvcnR5cGVzEjkKCmRpbWVuc2lvbnMY0Zvf3AEgAygLMhUubW9uaXRvcmluZy5EaW1lbnNpb2'
    '5SCmRpbWVuc2lvbnMSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSIQoKbWV0'
    'cmljbmFtZRj7vtoyIAEoCVIKbWV0cmljbmFtZRIgCgluYW1lc3BhY2UYwYS5qQEgASgJUgluYW'
    '1lc3BhY2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeAnomalyDetectorsOutputDescriptor instead')
const DescribeAnomalyDetectorsOutput$json = {
  '1': 'DescribeAnomalyDetectorsOutput',
  '2': [
    {
      '1': 'anomalydetectors',
      '3': 399089190,
      '4': 3,
      '5': 11,
      '6': '.monitoring.AnomalyDetector',
      '10': 'anomalydetectors'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAnomalyDetectorsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAnomalyDetectorsOutputDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZUFub21hbHlEZXRlY3RvcnNPdXRwdXQSSwoQYW5vbWFseWRldGVjdG9ycximvK'
        'a+ASADKAsyGy5tb25pdG9yaW5nLkFub21hbHlEZXRlY3RvclIQYW5vbWFseWRldGVjdG9ycxIf'
        'CgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeInsightRulesInputDescriptor instead')
const DescribeInsightRulesInput$json = {
  '1': 'DescribeInsightRulesInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInsightRulesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInsightRulesInputDescriptor =
    $convert.base64Decode(
        'ChlEZXNjcmliZUluc2lnaHRSdWxlc0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYX'
        'hyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeInsightRulesOutputDescriptor instead')
const DescribeInsightRulesOutput$json = {
  '1': 'DescribeInsightRulesOutput',
  '2': [
    {
      '1': 'insightrules',
      '3': 87263475,
      '4': 3,
      '5': 11,
      '6': '.monitoring.InsightRule',
      '10': 'insightrules'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInsightRulesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInsightRulesOutputDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZUluc2lnaHRSdWxlc091dHB1dBI+CgxpbnNpZ2h0cnVsZXMY85HOKSADKAsyFy'
        '5tb25pdG9yaW5nLkluc2lnaHRSdWxlUgxpbnNpZ2h0cnVsZXMSHwoJbmV4dHRva2VuGP6Eumcg'
        'ASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use dimensionDescriptor instead')
const Dimension$json = {
  '1': 'Dimension',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Dimension`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dimensionDescriptor = $convert.base64Decode(
    'CglEaW1lbnNpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIYCgV2YWx1ZRjr8p+KASABKAlSBX'
    'ZhbHVl');

@$core.Deprecated('Use dimensionFilterDescriptor instead')
const DimensionFilter$json = {
  '1': 'DimensionFilter',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `DimensionFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dimensionFilterDescriptor = $convert.base64Decode(
    'Cg9EaW1lbnNpb25GaWx0ZXISFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIYCgV2YWx1ZRjr8p+KAS'
    'ABKAlSBXZhbHVl');

@$core.Deprecated('Use disableAlarmActionsInputDescriptor instead')
const DisableAlarmActionsInput$json = {
  '1': 'DisableAlarmActionsInput',
  '2': [
    {'1': 'alarmnames', '3': 90874583, '4': 3, '5': 9, '10': 'alarmnames'},
  ],
};

/// Descriptor for `DisableAlarmActionsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableAlarmActionsInputDescriptor =
    $convert.base64Decode(
        'ChhEaXNhYmxlQWxhcm1BY3Rpb25zSW5wdXQSIQoKYWxhcm1uYW1lcxjXxaorIAMoCVIKYWxhcm'
        '1uYW1lcw==');

@$core.Deprecated('Use disableInsightRulesInputDescriptor instead')
const DisableInsightRulesInput$json = {
  '1': 'DisableInsightRulesInput',
  '2': [
    {'1': 'rulenames', '3': 267949170, '4': 3, '5': 9, '10': 'rulenames'},
  ],
};

/// Descriptor for `DisableInsightRulesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableInsightRulesInputDescriptor =
    $convert.base64Decode(
        'ChhEaXNhYmxlSW5zaWdodFJ1bGVzSW5wdXQSHwoJcnVsZW5hbWVzGPKo4n8gAygJUglydWxlbm'
        'FtZXM=');

@$core.Deprecated('Use disableInsightRulesOutputDescriptor instead')
const DisableInsightRulesOutput$json = {
  '1': 'DisableInsightRulesOutput',
  '2': [
    {
      '1': 'failures',
      '3': 335467271,
      '4': 3,
      '5': 11,
      '6': '.monitoring.PartialFailure',
      '10': 'failures'
    },
  ],
};

/// Descriptor for `DisableInsightRulesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableInsightRulesOutputDescriptor =
    $convert.base64Decode(
        'ChlEaXNhYmxlSW5zaWdodFJ1bGVzT3V0cHV0EjoKCGZhaWx1cmVzGIem+58BIAMoCzIaLm1vbm'
        'l0b3JpbmcuUGFydGlhbEZhaWx1cmVSCGZhaWx1cmVz');

@$core.Deprecated('Use enableAlarmActionsInputDescriptor instead')
const EnableAlarmActionsInput$json = {
  '1': 'EnableAlarmActionsInput',
  '2': [
    {'1': 'alarmnames', '3': 90874583, '4': 3, '5': 9, '10': 'alarmnames'},
  ],
};

/// Descriptor for `EnableAlarmActionsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableAlarmActionsInputDescriptor =
    $convert.base64Decode(
        'ChdFbmFibGVBbGFybUFjdGlvbnNJbnB1dBIhCgphbGFybW5hbWVzGNfFqisgAygJUgphbGFybW'
        '5hbWVz');

@$core.Deprecated('Use enableInsightRulesInputDescriptor instead')
const EnableInsightRulesInput$json = {
  '1': 'EnableInsightRulesInput',
  '2': [
    {'1': 'rulenames', '3': 267949170, '4': 3, '5': 9, '10': 'rulenames'},
  ],
};

/// Descriptor for `EnableInsightRulesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableInsightRulesInputDescriptor =
    $convert.base64Decode(
        'ChdFbmFibGVJbnNpZ2h0UnVsZXNJbnB1dBIfCglydWxlbmFtZXMY8qjifyADKAlSCXJ1bGVuYW'
        '1lcw==');

@$core.Deprecated('Use enableInsightRulesOutputDescriptor instead')
const EnableInsightRulesOutput$json = {
  '1': 'EnableInsightRulesOutput',
  '2': [
    {
      '1': 'failures',
      '3': 335467271,
      '4': 3,
      '5': 11,
      '6': '.monitoring.PartialFailure',
      '10': 'failures'
    },
  ],
};

/// Descriptor for `EnableInsightRulesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableInsightRulesOutputDescriptor =
    $convert.base64Decode(
        'ChhFbmFibGVJbnNpZ2h0UnVsZXNPdXRwdXQSOgoIZmFpbHVyZXMYh6b7nwEgAygLMhoubW9uaX'
        'RvcmluZy5QYXJ0aWFsRmFpbHVyZVIIZmFpbHVyZXM=');

@$core.Deprecated('Use entityDescriptor instead')
const Entity$json = {
  '1': 'Entity',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Entity.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'keyattributes',
      '3': 462259374,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Entity.KeyattributesEntry',
      '10': 'keyattributes'
    },
  ],
  '3': [Entity_AttributesEntry$json, Entity_KeyattributesEntry$json],
};

@$core.Deprecated('Use entityDescriptor instead')
const Entity_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use entityDescriptor instead')
const Entity_KeyattributesEntry$json = {
  '1': 'KeyattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Entity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityDescriptor = $convert.base64Decode(
    'CgZFbnRpdHkSRQoKYXR0cmlidXRlcxi1qftjIAMoCzIiLm1vbml0b3JpbmcuRW50aXR5LkF0dH'
    'JpYnV0ZXNFbnRyeVIKYXR0cmlidXRlcxJPCg1rZXlhdHRyaWJ1dGVzGK6JttwBIAMoCzIlLm1v'
    'bml0b3JpbmcuRW50aXR5LktleWF0dHJpYnV0ZXNFbnRyeVINa2V5YXR0cmlidXRlcxo9Cg9BdH'
    'RyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4'
    'ARpAChJLZXlhdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKA'
    'lSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use entityMetricDataDescriptor instead')
const EntityMetricData$json = {
  '1': 'EntityMetricData',
  '2': [
    {
      '1': 'entity',
      '3': 10171131,
      '4': 1,
      '5': 11,
      '6': '.monitoring.Entity',
      '10': 'entity'
    },
    {
      '1': 'metricdata',
      '3': 162960562,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricDatum',
      '10': 'metricdata'
    },
  ],
};

/// Descriptor for `EntityMetricData`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List entityMetricDataDescriptor = $convert.base64Decode(
    'ChBFbnRpdHlNZXRyaWNEYXRhEi0KBmVudGl0eRj75ewEIAEoCzISLm1vbml0b3JpbmcuRW50aX'
    'R5UgZlbnRpdHkSOgoKbWV0cmljZGF0YRiyqdpNIAMoCzIXLm1vbml0b3JpbmcuTWV0cmljRGF0'
    'dW1SCm1ldHJpY2RhdGE=');

@$core.Deprecated('Use getAlarmMuteRuleInputDescriptor instead')
const GetAlarmMuteRuleInput$json = {
  '1': 'GetAlarmMuteRuleInput',
  '2': [
    {
      '1': 'alarmmuterulename',
      '3': 530877511,
      '4': 1,
      '5': 9,
      '10': 'alarmmuterulename'
    },
  ],
};

/// Descriptor for `GetAlarmMuteRuleInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAlarmMuteRuleInputDescriptor = $convert.base64Decode(
    'ChVHZXRBbGFybU11dGVSdWxlSW5wdXQSMAoRYWxhcm1tdXRlcnVsZW5hbWUYx5iS/QEgASgJUh'
    'FhbGFybW11dGVydWxlbmFtZQ==');

@$core.Deprecated('Use getAlarmMuteRuleOutputDescriptor instead')
const GetAlarmMuteRuleOutput$json = {
  '1': 'GetAlarmMuteRuleOutput',
  '2': [
    {
      '1': 'alarmmuterulearn',
      '3': 493590237,
      '4': 1,
      '5': 9,
      '10': 'alarmmuterulearn'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'expiredate', '3': 454143207, '4': 1, '5': 9, '10': 'expiredate'},
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {
      '1': 'mutetargets',
      '3': 356717131,
      '4': 1,
      '5': 11,
      '6': '.monitoring.MuteTargets',
      '10': 'mutetargets'
    },
    {'1': 'mutetype', '3': 111213245, '4': 1, '5': 9, '10': 'mutetype'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'rule',
      '3': 475696372,
      '4': 1,
      '5': 11,
      '6': '.monitoring.Rule',
      '10': 'rule'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.monitoring.AlarmMuteRuleStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetAlarmMuteRuleOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAlarmMuteRuleOutputDescriptor = $convert.base64Decode(
    'ChZHZXRBbGFybU11dGVSdWxlT3V0cHV0Ei4KEGFsYXJtbXV0ZXJ1bGVhcm4Y3a2u6wEgASgJUh'
    'BhbGFybW11dGVydWxlYXJuEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIi'
    'CgpleHBpcmVkYXRlGOfZxtgBIAEoCVIKZXhwaXJlZGF0ZRI1ChRsYXN0dXBkYXRlZHRpbWVzdG'
    'FtcBiVy8g/IAEoCVIUbGFzdHVwZGF0ZWR0aW1lc3RhbXASPQoLbXV0ZXRhcmdldHMYy6SMqgEg'
    'ASgLMhcubW9uaXRvcmluZy5NdXRlVGFyZ2V0c1ILbXV0ZXRhcmdldHMSHQoIbXV0ZXR5cGUYvf'
    'WDNSABKAlSCG11dGV0eXBlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSKAoEcnVsZRj0meriASAB'
    'KAsyEC5tb25pdG9yaW5nLlJ1bGVSBHJ1bGUSIAoJc3RhcnRkYXRlGPz4oNQBIAEoCVIJc3Rhcn'
    'RkYXRlEjoKBnN0YXR1cxiQ5PsCIAEoDjIfLm1vbml0b3JpbmcuQWxhcm1NdXRlUnVsZVN0YXR1'
    'c1IGc3RhdHVz');

@$core.Deprecated('Use getDashboardInputDescriptor instead')
const GetDashboardInput$json = {
  '1': 'GetDashboardInput',
  '2': [
    {
      '1': 'dashboardname',
      '3': 506599873,
      '4': 1,
      '5': 9,
      '10': 'dashboardname'
    },
  ],
};

/// Descriptor for `GetDashboardInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDashboardInputDescriptor = $convert.base64Decode(
    'ChFHZXREYXNoYm9hcmRJbnB1dBIoCg1kYXNoYm9hcmRuYW1lGMGzyPEBIAEoCVINZGFzaGJvYX'
    'JkbmFtZQ==');

@$core.Deprecated('Use getDashboardOutputDescriptor instead')
const GetDashboardOutput$json = {
  '1': 'GetDashboardOutput',
  '2': [
    {'1': 'dashboardarn', '3': 108051951, '4': 1, '5': 9, '10': 'dashboardarn'},
    {'1': 'dashboardbody', '3': 3403236, '4': 1, '5': 9, '10': 'dashboardbody'},
    {
      '1': 'dashboardname',
      '3': 506599873,
      '4': 1,
      '5': 9,
      '10': 'dashboardname'
    },
  ],
};

/// Descriptor for `GetDashboardOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDashboardOutputDescriptor = $convert.base64Decode(
    'ChJHZXREYXNoYm9hcmRPdXRwdXQSJQoMZGFzaGJvYXJkYXJuGO/7wjMgASgJUgxkYXNoYm9hcm'
    'Rhcm4SJwoNZGFzaGJvYXJkYm9keRjk288BIAEoCVINZGFzaGJvYXJkYm9keRIoCg1kYXNoYm9h'
    'cmRuYW1lGMGzyPEBIAEoCVINZGFzaGJvYXJkbmFtZQ==');

@$core.Deprecated('Use getInsightRuleReportInputDescriptor instead')
const GetInsightRuleReportInput$json = {
  '1': 'GetInsightRuleReportInput',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {
      '1': 'maxcontributorcount',
      '3': 272640614,
      '4': 1,
      '5': 5,
      '10': 'maxcontributorcount'
    },
    {'1': 'metrics', '3': 436365847, '4': 3, '5': 9, '10': 'metrics'},
    {'1': 'orderby', '3': 353010019, '4': 1, '5': 9, '10': 'orderby'},
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {'1': 'rulename', '3': 214688793, '4': 1, '5': 9, '10': 'rulename'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
  ],
};

/// Descriptor for `GetInsightRuleReportInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInsightRuleReportInputDescriptor = $convert.base64Decode(
    'ChlHZXRJbnNpZ2h0UnVsZVJlcG9ydElucHV0EhsKB2VuZHRpbWUYzO+8HiABKAlSB2VuZHRpbW'
    'USNAoTbWF4Y29udHJpYnV0b3Jjb3VudBjm1ICCASABKAVSE21heGNvbnRyaWJ1dG9yY291bnQS'
    'HAoHbWV0cmljcxiX1InQASADKAlSB21ldHJpY3MSHAoHb3JkZXJieRjjgqqoASABKAlSB29yZG'
    'VyYnkSGQoGcGVyaW9kGKWIkjkgASgFUgZwZXJpb2QSHQoIcnVsZW5hbWUYmcivZiABKAlSCHJ1'
    'bGVuYW1lEiAKCXN0YXJ0dGltZRjvtOWwASABKAlSCXN0YXJ0dGltZQ==');

@$core.Deprecated('Use getInsightRuleReportOutputDescriptor instead')
const GetInsightRuleReportOutput$json = {
  '1': 'GetInsightRuleReportOutput',
  '2': [
    {
      '1': 'aggregatevalue',
      '3': 225313242,
      '4': 1,
      '5': 1,
      '10': 'aggregatevalue'
    },
    {
      '1': 'aggregationstatistic',
      '3': 333233758,
      '4': 1,
      '5': 9,
      '10': 'aggregationstatistic'
    },
    {
      '1': 'approximateuniquecount',
      '3': 55617722,
      '4': 1,
      '5': 3,
      '10': 'approximateuniquecount'
    },
    {
      '1': 'contributors',
      '3': 168168444,
      '4': 3,
      '5': 11,
      '6': '.monitoring.InsightRuleContributor',
      '10': 'contributors'
    },
    {'1': 'keylabels', '3': 16074444, '4': 3, '5': 9, '10': 'keylabels'},
    {
      '1': 'metricdatapoints',
      '3': 117991067,
      '4': 3,
      '5': 11,
      '6': '.monitoring.InsightRuleMetricDatapoint',
      '10': 'metricdatapoints'
    },
  ],
};

/// Descriptor for `GetInsightRuleReportOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInsightRuleReportOutputDescriptor = $convert.base64Decode(
    'ChpHZXRJbnNpZ2h0UnVsZVJlcG9ydE91dHB1dBIpCg5hZ2dyZWdhdGV2YWx1ZRjag7hrIAEoAV'
    'IOYWdncmVnYXRldmFsdWUSNgoUYWdncmVnYXRpb25zdGF0aXN0aWMY3vzyngEgASgJUhRhZ2dy'
    'ZWdhdGlvbnN0YXRpc3RpYxI5ChZhcHByb3hpbWF0ZXVuaXF1ZWNvdW50GLrRwhogASgDUhZhcH'
    'Byb3hpbWF0ZXVuaXF1ZWNvdW50EkkKDGNvbnRyaWJ1dG9ycxj8l5hQIAMoCzIiLm1vbml0b3Jp'
    'bmcuSW5zaWdodFJ1bGVDb250cmlidXRvclIMY29udHJpYnV0b3JzEh8KCWtleWxhYmVscxjMjd'
    'UHIAMoCVIJa2V5bGFiZWxzElUKEG1ldHJpY2RhdGFwb2ludHMYm82hOCADKAsyJi5tb25pdG9y'
    'aW5nLkluc2lnaHRSdWxlTWV0cmljRGF0YXBvaW50UhBtZXRyaWNkYXRhcG9pbnRz');

@$core.Deprecated('Use getMetricDataInputDescriptor instead')
const GetMetricDataInput$json = {
  '1': 'GetMetricDataInput',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {
      '1': 'labeloptions',
      '3': 184949902,
      '4': 1,
      '5': 11,
      '6': '.monitoring.LabelOptions',
      '10': 'labeloptions'
    },
    {
      '1': 'maxdatapoints',
      '3': 97504835,
      '4': 1,
      '5': 5,
      '10': 'maxdatapoints'
    },
    {
      '1': 'metricdataqueries',
      '3': 92209504,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricDataQuery',
      '10': 'metricdataqueries'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'scanby',
      '3': 344142854,
      '4': 1,
      '5': 14,
      '6': '.monitoring.ScanBy',
      '10': 'scanby'
    },
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
  ],
};

/// Descriptor for `GetMetricDataInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMetricDataInputDescriptor = $convert.base64Decode(
    'ChJHZXRNZXRyaWNEYXRhSW5wdXQSGwoHZW5kdGltZRjM77weIAEoCVIHZW5kdGltZRI/CgxsYW'
    'JlbG9wdGlvbnMYjrmYWCABKAsyGC5tb25pdG9yaW5nLkxhYmVsT3B0aW9uc1IMbGFiZWxvcHRp'
    'b25zEicKDW1heGRhdGFwb2ludHMYw5y/LiABKAVSDW1heGRhdGFwb2ludHMSTAoRbWV0cmljZG'
    'F0YXF1ZXJpZXMY4IL8KyADKAsyGy5tb25pdG9yaW5nLk1ldHJpY0RhdGFRdWVyeVIRbWV0cmlj'
    'ZGF0YXF1ZXJpZXMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SLgoGc2NhbmJ5GI'
    'bojKQBIAEoDjISLm1vbml0b3JpbmcuU2NhbkJ5UgZzY2FuYnkSIAoJc3RhcnR0aW1lGO+05bAB'
    'IAEoCVIJc3RhcnR0aW1l');

@$core.Deprecated('Use getMetricDataOutputDescriptor instead')
const GetMetricDataOutput$json = {
  '1': 'GetMetricDataOutput',
  '2': [
    {
      '1': 'messages',
      '3': 409018326,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MessageData',
      '10': 'messages'
    },
    {
      '1': 'metricdataresults',
      '3': 373784842,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricDataResult',
      '10': 'metricdataresults'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `GetMetricDataOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMetricDataOutputDescriptor = $convert.base64Decode(
    'ChNHZXRNZXRyaWNEYXRhT3V0cHV0EjcKCG1lc3NhZ2VzGNa/hMMBIAMoCzIXLm1vbml0b3Jpbm'
    'cuTWVzc2FnZURhdGFSCG1lc3NhZ2VzEk4KEW1ldHJpY2RhdGFyZXN1bHRzGIqCnrIBIAMoCzIc'
    'Lm1vbml0b3JpbmcuTWV0cmljRGF0YVJlc3VsdFIRbWV0cmljZGF0YXJlc3VsdHMSHwoJbmV4dH'
    'Rva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use getMetricStatisticsInputDescriptor instead')
const GetMetricStatisticsInput$json = {
  '1': 'GetMetricStatisticsInput',
  '2': [
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {
      '1': 'extendedstatistics',
      '3': 365818700,
      '4': 3,
      '5': 9,
      '10': 'extendedstatistics'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {
      '1': 'statistics',
      '3': 510636075,
      '4': 3,
      '5': 14,
      '6': '.monitoring.Statistic',
      '10': 'statistics'
    },
    {
      '1': 'unit',
      '3': 148989480,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StandardUnit',
      '10': 'unit'
    },
  ],
};

/// Descriptor for `GetMetricStatisticsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMetricStatisticsInputDescriptor = $convert.base64Decode(
    'ChhHZXRNZXRyaWNTdGF0aXN0aWNzSW5wdXQSOQoKZGltZW5zaW9ucxjRm9/cASADKAsyFS5tb2'
    '5pdG9yaW5nLkRpbWVuc2lvblIKZGltZW5zaW9ucxIbCgdlbmR0aW1lGMzvvB4gASgJUgdlbmR0'
    'aW1lEjIKEmV4dGVuZGVkc3RhdGlzdGljcxjM5reuASADKAlSEmV4dGVuZGVkc3RhdGlzdGljcx'
    'IhCgptZXRyaWNuYW1lGPu+2jIgASgJUgptZXRyaWNuYW1lEiAKCW5hbWVzcGFjZRjBhLmpASAB'
    'KAlSCW5hbWVzcGFjZRIZCgZwZXJpb2QYpYiSOSABKAVSBnBlcmlvZBIgCglzdGFydHRpbWUY77'
    'TlsAEgASgJUglzdGFydHRpbWUSOQoKc3RhdGlzdGljcxir4L7zASADKA4yFS5tb25pdG9yaW5n'
    'LlN0YXRpc3RpY1IKc3RhdGlzdGljcxIvCgR1bml0GKjMhUcgASgOMhgubW9uaXRvcmluZy5TdG'
    'FuZGFyZFVuaXRSBHVuaXQ=');

@$core.Deprecated('Use getMetricStatisticsOutputDescriptor instead')
const GetMetricStatisticsOutput$json = {
  '1': 'GetMetricStatisticsOutput',
  '2': [
    {
      '1': 'datapoints',
      '3': 37512135,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Datapoint',
      '10': 'datapoints'
    },
    {'1': 'label', '3': 516747934, '4': 1, '5': 9, '10': 'label'},
  ],
};

/// Descriptor for `GetMetricStatisticsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMetricStatisticsOutputDescriptor = $convert.base64Decode(
    'ChlHZXRNZXRyaWNTdGF0aXN0aWNzT3V0cHV0EjgKCmRhdGFwb2ludHMYx8fxESADKAsyFS5tb2'
    '5pdG9yaW5nLkRhdGFwb2ludFIKZGF0YXBvaW50cxIYCgVsYWJlbBie5bP2ASABKAlSBWxhYmVs');

@$core.Deprecated('Use getMetricStreamInputDescriptor instead')
const GetMetricStreamInput$json = {
  '1': 'GetMetricStreamInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetMetricStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMetricStreamInputDescriptor =
    $convert.base64Decode(
        'ChRHZXRNZXRyaWNTdHJlYW1JbnB1dBIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use getMetricStreamOutputDescriptor instead')
const GetMetricStreamOutput$json = {
  '1': 'GetMetricStreamOutput',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'excludefilters',
      '3': 426990841,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricStreamFilter',
      '10': 'excludefilters'
    },
    {'1': 'firehosearn', '3': 55494200, '4': 1, '5': 9, '10': 'firehosearn'},
    {
      '1': 'includefilters',
      '3': 90967031,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricStreamFilter',
      '10': 'includefilters'
    },
    {
      '1': 'includelinkedaccountsmetrics',
      '3': 356901952,
      '4': 1,
      '5': 8,
      '10': 'includelinkedaccountsmetrics'
    },
    {
      '1': 'lastupdatedate',
      '3': 65755725,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputformat',
      '3': 519139704,
      '4': 1,
      '5': 14,
      '6': '.monitoring.MetricStreamOutputFormat',
      '10': 'outputformat'
    },
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
    {
      '1': 'statisticsconfigurations',
      '3': 175876178,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricStreamStatisticsConfiguration',
      '10': 'statisticsconfigurations'
    },
  ],
};

/// Descriptor for `GetMetricStreamOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMetricStreamOutputDescriptor = $convert.base64Decode(
    'ChVHZXRNZXRyaWNTdHJlYW1PdXRwdXQSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEiYKDGNyZWF0aW'
    '9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW9uZGF0ZRJKCg5leGNsdWRlZmlsdGVycxj5uc3LASAD'
    'KAsyHi5tb25pdG9yaW5nLk1ldHJpY1N0cmVhbUZpbHRlclIOZXhjbHVkZWZpbHRlcnMSIwoLZm'
    'lyZWhvc2Vhcm4YuIy7GiABKAlSC2ZpcmVob3NlYXJuEkkKDmluY2x1ZGVmaWx0ZXJzGPeXsCsg'
    'AygLMh4ubW9uaXRvcmluZy5NZXRyaWNTdHJlYW1GaWx0ZXJSDmluY2x1ZGVmaWx0ZXJzEkYKHG'
    'luY2x1ZGVsaW5rZWRhY2NvdW50c21ldHJpY3MYwMiXqgEgASgIUhxpbmNsdWRlbGlua2VkYWNj'
    'b3VudHNtZXRyaWNzEikKDmxhc3R1cGRhdGVkYXRlGM20rR8gASgJUg5sYXN0dXBkYXRlZGF0ZR'
    'IVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkwKDG91dHB1dGZvcm1hdBj44sX3ASABKA4yJC5tb25p'
    'dG9yaW5nLk1ldHJpY1N0cmVhbU91dHB1dEZvcm1hdFIMb3V0cHV0Zm9ybWF0EhwKB3JvbGVhcm'
    '4YgfjnmQEgASgJUgdyb2xlYXJuEhgKBXN0YXRlGJfJsu8BIAEoCVIFc3RhdGUSbgoYc3RhdGlz'
    'dGljc2NvbmZpZ3VyYXRpb25zGNLQ7lMgAygLMi8ubW9uaXRvcmluZy5NZXRyaWNTdHJlYW1TdG'
    'F0aXN0aWNzQ29uZmlndXJhdGlvblIYc3RhdGlzdGljc2NvbmZpZ3VyYXRpb25z');

@$core.Deprecated('Use getMetricWidgetImageInputDescriptor instead')
const GetMetricWidgetImageInput$json = {
  '1': 'GetMetricWidgetImageInput',
  '2': [
    {'1': 'metricwidget', '3': 283026730, '4': 1, '5': 9, '10': 'metricwidget'},
    {'1': 'outputformat', '3': 519139704, '4': 1, '5': 9, '10': 'outputformat'},
  ],
};

/// Descriptor for `GetMetricWidgetImageInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMetricWidgetImageInputDescriptor =
    $convert.base64Decode(
        'ChlHZXRNZXRyaWNXaWRnZXRJbWFnZUlucHV0EiYKDG1ldHJpY3dpZGdldBiqyvqGASABKAlSDG'
        '1ldHJpY3dpZGdldBImCgxvdXRwdXRmb3JtYXQY+OLF9wEgASgJUgxvdXRwdXRmb3JtYXQ=');

@$core.Deprecated('Use getMetricWidgetImageOutputDescriptor instead')
const GetMetricWidgetImageOutput$json = {
  '1': 'GetMetricWidgetImageOutput',
  '2': [
    {
      '1': 'metricwidgetimage',
      '3': 484815311,
      '4': 1,
      '5': 12,
      '10': 'metricwidgetimage'
    },
  ],
};

/// Descriptor for `GetMetricWidgetImageOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMetricWidgetImageOutputDescriptor =
    $convert.base64Decode(
        'ChpHZXRNZXRyaWNXaWRnZXRJbWFnZU91dHB1dBIwChFtZXRyaWN3aWRnZXRpbWFnZRjP45bnAS'
        'ABKAxSEW1ldHJpY3dpZGdldGltYWdl');

@$core.Deprecated('Use insightRuleDescriptor instead')
const InsightRule$json = {
  '1': 'InsightRule',
  '2': [
    {
      '1': 'applyontransformedlogs',
      '3': 30966533,
      '4': 1,
      '5': 8,
      '10': 'applyontransformedlogs'
    },
    {'1': 'definition', '3': 222998209, '4': 1, '5': 9, '10': 'definition'},
    {'1': 'managedrule', '3': 144612369, '4': 1, '5': 8, '10': 'managedrule'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'schema', '3': 412122455, '4': 1, '5': 9, '10': 'schema'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
  ],
};

/// Descriptor for `InsightRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insightRuleDescriptor = $convert.base64Decode(
    'CgtJbnNpZ2h0UnVsZRI5ChZhcHBseW9udHJhbnNmb3JtZWRsb2dzGIWG4g4gASgIUhZhcHBseW'
    '9udHJhbnNmb3JtZWRsb2dzEiEKCmRlZmluaXRpb24Ywd2qaiABKAlSCmRlZmluaXRpb24SIwoL'
    'bWFuYWdlZHJ1bGUYkbj6RCABKAhSC21hbmFnZWRydWxlEhUKBG5hbWUYh+aBfyABKAlSBG5hbW'
    'USGgoGc2NoZW1hGNf6wcQBIAEoCVIGc2NoZW1hEhgKBXN0YXRlGJfJsu8BIAEoCVIFc3RhdGU=');

@$core.Deprecated('Use insightRuleContributorDescriptor instead')
const InsightRuleContributor$json = {
  '1': 'InsightRuleContributor',
  '2': [
    {
      '1': 'approximateaggregatevalue',
      '3': 260431894,
      '4': 1,
      '5': 1,
      '10': 'approximateaggregatevalue'
    },
    {
      '1': 'datapoints',
      '3': 37512135,
      '4': 3,
      '5': 11,
      '6': '.monitoring.InsightRuleContributorDatapoint',
      '10': 'datapoints'
    },
    {'1': 'keys', '3': 2831086, '4': 3, '5': 9, '10': 'keys'},
  ],
};

/// Descriptor for `InsightRuleContributor`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insightRuleContributorDescriptor = $convert.base64Decode(
    'ChZJbnNpZ2h0UnVsZUNvbnRyaWJ1dG9yEj8KGWFwcHJveGltYXRlYWdncmVnYXRldmFsdWUYls'
    'CXfCABKAFSGWFwcHJveGltYXRlYWdncmVnYXRldmFsdWUSTgoKZGF0YXBvaW50cxjHx/ERIAMo'
    'CzIrLm1vbml0b3JpbmcuSW5zaWdodFJ1bGVDb250cmlidXRvckRhdGFwb2ludFIKZGF0YXBvaW'
    '50cxIVCgRrZXlzGO7lrAEgAygJUgRrZXlz');

@$core.Deprecated('Use insightRuleContributorDatapointDescriptor instead')
const InsightRuleContributorDatapoint$json = {
  '1': 'InsightRuleContributorDatapoint',
  '2': [
    {
      '1': 'approximatevalue',
      '3': 428761855,
      '4': 1,
      '5': 1,
      '10': 'approximatevalue'
    },
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
  ],
};

/// Descriptor for `InsightRuleContributorDatapoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insightRuleContributorDatapointDescriptor =
    $convert.base64Decode(
        'Ch9JbnNpZ2h0UnVsZUNvbnRyaWJ1dG9yRGF0YXBvaW50Ei4KEGFwcHJveGltYXRldmFsdWUY/8'
        'W5zAEgASgBUhBhcHByb3hpbWF0ZXZhbHVlEh8KCXRpbWVzdGFtcBjEw7dNIAEoCVIJdGltZXN0'
        'YW1w');

@$core.Deprecated('Use insightRuleMetricDatapointDescriptor instead')
const InsightRuleMetricDatapoint$json = {
  '1': 'InsightRuleMetricDatapoint',
  '2': [
    {'1': 'average', '3': 134712777, '4': 1, '5': 1, '10': 'average'},
    {
      '1': 'maxcontributorvalue',
      '3': 389006680,
      '4': 1,
      '5': 1,
      '10': 'maxcontributorvalue'
    },
    {'1': 'maximum', '3': 43343394, '4': 1, '5': 1, '10': 'maximum'},
    {'1': 'minimum', '3': 438536856, '4': 1, '5': 1, '10': 'minimum'},
    {'1': 'samplecount', '3': 384919839, '4': 1, '5': 1, '10': 'samplecount'},
    {'1': 'sum', '3': 534303305, '4': 1, '5': 1, '10': 'sum'},
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
    {
      '1': 'uniquecontributors',
      '3': 504374565,
      '4': 1,
      '5': 1,
      '10': 'uniquecontributors'
    },
  ],
};

/// Descriptor for `InsightRuleMetricDatapoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insightRuleMetricDatapointDescriptor = $convert.base64Decode(
    'ChpJbnNpZ2h0UnVsZU1ldHJpY0RhdGFwb2ludBIbCgdhdmVyYWdlGMmbnkAgASgBUgdhdmVyYW'
    'dlEjQKE21heGNvbnRyaWJ1dG9ydmFsdWUY2Iq/uQEgASgBUhNtYXhjb250cmlidXRvcnZhbHVl'
    'EhsKB21heGltdW0YorzVFCABKAFSB21heGltdW0SHAoHbWluaW11bRiYlY7RASABKAFSB21pbm'
    'ltdW0SJAoLc2FtcGxlY291bnQYn9LFtwEgASgBUgtzYW1wbGVjb3VudBIUCgNzdW0YyaTj/gEg'
    'ASgBUgNzdW0SHwoJdGltZXN0YW1wGMTDt00gASgJUgl0aW1lc3RhbXASMgoSdW5pcXVlY29udH'
    'JpYnV0b3JzGKXKwPABIAEoAVISdW5pcXVlY29udHJpYnV0b3Jz');

@$core.Deprecated('Use internalServiceFaultDescriptor instead')
const InternalServiceFault$json = {
  '1': 'InternalServiceFault',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalServiceFault`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalServiceFaultDescriptor =
    $convert.base64Decode(
        'ChRJbnRlcm5hbFNlcnZpY2VGYXVsdBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidFormatFaultDescriptor instead')
const InvalidFormatFault$json = {
  '1': 'InvalidFormatFault',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidFormatFault`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidFormatFaultDescriptor =
    $convert.base64Decode(
        'ChJJbnZhbGlkRm9ybWF0RmF1bHQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidNextTokenDescriptor instead')
const InvalidNextToken$json = {
  '1': 'InvalidNextToken',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidNextToken`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidNextTokenDescriptor = $convert.base64Decode(
    'ChBJbnZhbGlkTmV4dFRva2VuEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidParameterCombinationExceptionDescriptor instead')
const InvalidParameterCombinationException$json = {
  '1': 'InvalidParameterCombinationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterCombinationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterCombinationExceptionDescriptor =
    $convert.base64Decode(
        'CiRJbnZhbGlkUGFyYW1ldGVyQ29tYmluYXRpb25FeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidParameterValueExceptionDescriptor instead')
const InvalidParameterValueException$json = {
  '1': 'InvalidParameterValueException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterValueException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterValueExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkUGFyYW1ldGVyVmFsdWVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use labelOptionsDescriptor instead')
const LabelOptions$json = {
  '1': 'LabelOptions',
  '2': [
    {'1': 'timezone', '3': 246302531, '4': 1, '5': 9, '10': 'timezone'},
  ],
};

/// Descriptor for `LabelOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List labelOptionsDescriptor = $convert.base64Decode(
    'CgxMYWJlbE9wdGlvbnMSHQoIdGltZXpvbmUYw465dSABKAlSCHRpbWV6b25l');

@$core.Deprecated('Use limitExceededExceptionDescriptor instead')
const LimitExceededException$json = {
  '1': 'LimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChZMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use limitExceededFaultDescriptor instead')
const LimitExceededFault$json = {
  '1': 'LimitExceededFault',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LimitExceededFault`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitExceededFaultDescriptor =
    $convert.base64Decode(
        'ChJMaW1pdEV4Y2VlZGVkRmF1bHQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use listAlarmMuteRulesInputDescriptor instead')
const ListAlarmMuteRulesInput$json = {
  '1': 'ListAlarmMuteRulesInput',
  '2': [
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {'1': 'maxrecords', '3': 220314370, '4': 1, '5': 5, '10': 'maxrecords'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'statuses',
      '3': 374056024,
      '4': 3,
      '5': 14,
      '6': '.monitoring.AlarmMuteRuleStatus',
      '10': 'statuses'
    },
  ],
};

/// Descriptor for `ListAlarmMuteRulesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAlarmMuteRulesInputDescriptor = $convert.base64Decode(
    'ChdMaXN0QWxhcm1NdXRlUnVsZXNJbnB1dBIfCglhbGFybW5hbWUY4t/lGSABKAlSCWFsYXJtbm'
    'FtZRIhCgptYXhyZWNvcmRzGIL2hmkgASgFUgptYXhyZWNvcmRzEh8KCW5leHR0b2tlbhj+hLpn'
    'IAEoCVIJbmV4dHRva2VuEj8KCHN0YXR1c2VzGNjIrrIBIAMoDjIfLm1vbml0b3JpbmcuQWxhcm'
    '1NdXRlUnVsZVN0YXR1c1IIc3RhdHVzZXM=');

@$core.Deprecated('Use listAlarmMuteRulesOutputDescriptor instead')
const ListAlarmMuteRulesOutput$json = {
  '1': 'ListAlarmMuteRulesOutput',
  '2': [
    {
      '1': 'alarmmuterulesummaries',
      '3': 440906836,
      '4': 3,
      '5': 11,
      '6': '.monitoring.AlarmMuteRuleSummary',
      '10': 'alarmmuterulesummaries'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListAlarmMuteRulesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAlarmMuteRulesOutputDescriptor = $convert.base64Decode(
    'ChhMaXN0QWxhcm1NdXRlUnVsZXNPdXRwdXQSXAoWYWxhcm1tdXRlcnVsZXN1bW1hcmllcxjU6J'
    '7SASADKAsyIC5tb25pdG9yaW5nLkFsYXJtTXV0ZVJ1bGVTdW1tYXJ5UhZhbGFybW11dGVydWxl'
    'c3VtbWFyaWVzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listDashboardsInputDescriptor instead')
const ListDashboardsInput$json = {
  '1': 'ListDashboardsInput',
  '2': [
    {
      '1': 'dashboardnameprefix',
      '3': 301015421,
      '4': 1,
      '5': 9,
      '10': 'dashboardnameprefix'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDashboardsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDashboardsInputDescriptor = $convert.base64Decode(
    'ChNMaXN0RGFzaGJvYXJkc0lucHV0EjQKE2Rhc2hib2FyZG5hbWVwcmVmaXgY/cLEjwEgASgJUh'
    'NkYXNoYm9hcmRuYW1lcHJlZml4Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listDashboardsOutputDescriptor instead')
const ListDashboardsOutput$json = {
  '1': 'ListDashboardsOutput',
  '2': [
    {
      '1': 'dashboardentries',
      '3': 281104826,
      '4': 3,
      '5': 11,
      '6': '.monitoring.DashboardEntry',
      '10': 'dashboardentries'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDashboardsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDashboardsOutputDescriptor = $convert.base64Decode(
    'ChRMaXN0RGFzaGJvYXJkc091dHB1dBJKChBkYXNoYm9hcmRlbnRyaWVzGLqjhYYBIAMoCzIaLm'
    '1vbml0b3JpbmcuRGFzaGJvYXJkRW50cnlSEGRhc2hib2FyZGVudHJpZXMSHwoJbmV4dHRva2Vu'
    'GP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listManagedInsightRulesInputDescriptor instead')
const ListManagedInsightRulesInput$json = {
  '1': 'ListManagedInsightRulesInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListManagedInsightRulesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listManagedInsightRulesInputDescriptor =
    $convert.base64Decode(
        'ChxMaXN0TWFuYWdlZEluc2lnaHRSdWxlc0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUg'
        'ptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEiQKC3Jlc291cmNl'
        'YXJuGO3AmbABIAEoCVILcmVzb3VyY2Vhcm4=');

@$core.Deprecated('Use listManagedInsightRulesOutputDescriptor instead')
const ListManagedInsightRulesOutput$json = {
  '1': 'ListManagedInsightRulesOutput',
  '2': [
    {
      '1': 'managedrules',
      '3': 347090906,
      '4': 3,
      '5': 11,
      '6': '.monitoring.ManagedRuleDescription',
      '10': 'managedrules'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListManagedInsightRulesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listManagedInsightRulesOutputDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0TWFuYWdlZEluc2lnaHRSdWxlc091dHB1dBJKCgxtYW5hZ2VkcnVsZXMY2t/ApQEgAy'
        'gLMiIubW9uaXRvcmluZy5NYW5hZ2VkUnVsZURlc2NyaXB0aW9uUgxtYW5hZ2VkcnVsZXMSHwoJ'
        'bmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listMetricStreamsInputDescriptor instead')
const ListMetricStreamsInput$json = {
  '1': 'ListMetricStreamsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListMetricStreamsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMetricStreamsInputDescriptor =
    $convert.base64Decode(
        'ChZMaXN0TWV0cmljU3RyZWFtc0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZX'
        'N1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listMetricStreamsOutputDescriptor instead')
const ListMetricStreamsOutput$json = {
  '1': 'ListMetricStreamsOutput',
  '2': [
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricStreamEntry',
      '10': 'entries'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListMetricStreamsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMetricStreamsOutputDescriptor = $convert.base64Decode(
    'ChdMaXN0TWV0cmljU3RyZWFtc091dHB1dBI7CgdlbnRyaWVzGJTFsuUBIAMoCzIdLm1vbml0b3'
    'JpbmcuTWV0cmljU3RyZWFtRW50cnlSB2VudHJpZXMSHwoJbmV4dHRva2VuGP6EumcgASgJUglu'
    'ZXh0dG9rZW4=');

@$core.Deprecated('Use listMetricsInputDescriptor instead')
const ListMetricsInput$json = {
  '1': 'ListMetricsInput',
  '2': [
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.DimensionFilter',
      '10': 'dimensions'
    },
    {
      '1': 'includelinkedaccounts',
      '3': 228605043,
      '4': 1,
      '5': 8,
      '10': 'includelinkedaccounts'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'owningaccount',
      '3': 339968073,
      '4': 1,
      '5': 9,
      '10': 'owningaccount'
    },
    {
      '1': 'recentlyactive',
      '3': 179617336,
      '4': 1,
      '5': 14,
      '6': '.monitoring.RecentlyActive',
      '10': 'recentlyactive'
    },
  ],
};

/// Descriptor for `ListMetricsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMetricsInputDescriptor = $convert.base64Decode(
    'ChBMaXN0TWV0cmljc0lucHV0Ej8KCmRpbWVuc2lvbnMY0Zvf3AEgAygLMhsubW9uaXRvcmluZy'
    '5EaW1lbnNpb25GaWx0ZXJSCmRpbWVuc2lvbnMSNwoVaW5jbHVkZWxpbmtlZGFjY291bnRzGPP4'
    'gG0gASgIUhVpbmNsdWRlbGlua2VkYWNjb3VudHMSIQoKbWV0cmljbmFtZRj7vtoyIAEoCVIKbW'
    'V0cmljbmFtZRIgCgluYW1lc3BhY2UYwYS5qQEgASgJUgluYW1lc3BhY2USHwoJbmV4dHRva2Vu'
    'GP6EumcgASgJUgluZXh0dG9rZW4SKAoNb3duaW5nYWNjb3VudBjJgI6iASABKAlSDW93bmluZ2'
    'FjY291bnQSRQoOcmVjZW50bHlhY3RpdmUYuPzSVSABKA4yGi5tb25pdG9yaW5nLlJlY2VudGx5'
    'QWN0aXZlUg5yZWNlbnRseWFjdGl2ZQ==');

@$core.Deprecated('Use listMetricsOutputDescriptor instead')
const ListMetricsOutput$json = {
  '1': 'ListMetricsOutput',
  '2': [
    {
      '1': 'metrics',
      '3': 436365847,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Metric',
      '10': 'metrics'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'owningaccounts',
      '3': 21159138,
      '4': 3,
      '5': 9,
      '10': 'owningaccounts'
    },
  ],
};

/// Descriptor for `ListMetricsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMetricsOutputDescriptor = $convert.base64Decode(
    'ChFMaXN0TWV0cmljc091dHB1dBIwCgdtZXRyaWNzGJfUidABIAMoCzISLm1vbml0b3JpbmcuTW'
    'V0cmljUgdtZXRyaWNzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEikKDm93bmlu'
    'Z2FjY291bnRzGOK5iwogAygJUg5vd25pbmdhY2NvdW50cw==');

@$core.Deprecated('Use listTagsForResourceInputDescriptor instead')
const ListTagsForResourceInput$json = {
  '1': 'ListTagsForResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceInputDescriptor =
    $convert.base64Decode(
        'ChhMaXN0VGFnc0ZvclJlc291cmNlSW5wdXQSJAoLcmVzb3VyY2Vhcm4Y7cCZsAEgASgJUgtyZX'
        'NvdXJjZWFybg==');

@$core.Deprecated('Use listTagsForResourceOutputDescriptor instead')
const ListTagsForResourceOutput$json = {
  '1': 'ListTagsForResourceOutput',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceOutputDescriptor =
    $convert.base64Decode(
        'ChlMaXN0VGFnc0ZvclJlc291cmNlT3V0cHV0EicKBHRhZ3MYwcH2tQEgAygLMg8ubW9uaXRvcm'
        'luZy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use managedRuleDescriptor instead')
const ManagedRule$json = {
  '1': 'ManagedRule',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Tag',
      '10': 'tags'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `ManagedRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedRuleDescriptor = $convert.base64Decode(
    'CgtNYW5hZ2VkUnVsZRIkCgtyZXNvdXJjZWFybhjtwJmwASABKAlSC3Jlc291cmNlYXJuEicKBH'
    'RhZ3MYwcH2tQEgAygLMg8ubW9uaXRvcmluZy5UYWdSBHRhZ3MSJgoMdGVtcGxhdGVuYW1lGLGY'
    'keUBIAEoCVIMdGVtcGxhdGVuYW1l');

@$core.Deprecated('Use managedRuleDescriptionDescriptor instead')
const ManagedRuleDescription$json = {
  '1': 'ManagedRuleDescription',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'rulestate',
      '3': 419921189,
      '4': 1,
      '5': 11,
      '6': '.monitoring.ManagedRuleState',
      '10': 'rulestate'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `ManagedRuleDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedRuleDescriptionDescriptor = $convert.base64Decode(
    'ChZNYW5hZ2VkUnVsZURlc2NyaXB0aW9uEiQKC3Jlc291cmNlYXJuGO3AmbABIAEoCVILcmVzb3'
    'VyY2Vhcm4SPgoJcnVsZXN0YXRlGKX6ncgBIAEoCzIcLm1vbml0b3JpbmcuTWFuYWdlZFJ1bGVT'
    'dGF0ZVIJcnVsZXN0YXRlEiYKDHRlbXBsYXRlbmFtZRixmJHlASABKAlSDHRlbXBsYXRlbmFtZQ'
    '==');

@$core.Deprecated('Use managedRuleStateDescriptor instead')
const ManagedRuleState$json = {
  '1': 'ManagedRuleState',
  '2': [
    {'1': 'rulename', '3': 214688793, '4': 1, '5': 9, '10': 'rulename'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
  ],
};

/// Descriptor for `ManagedRuleState`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedRuleStateDescriptor = $convert.base64Decode(
    'ChBNYW5hZ2VkUnVsZVN0YXRlEh0KCHJ1bGVuYW1lGJnIr2YgASgJUghydWxlbmFtZRIYCgVzdG'
    'F0ZRiXybLvASABKAlSBXN0YXRl');

@$core.Deprecated('Use messageDataDescriptor instead')
const MessageData$json = {
  '1': 'MessageData',
  '2': [
    {'1': 'code', '3': 425572629, '4': 1, '5': 9, '10': 'code'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `MessageData`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageDataDescriptor = $convert.base64Decode(
    'CgtNZXNzYWdlRGF0YRIWCgRjb2RlGJXy9soBIAEoCVIEY29kZRIYCgV2YWx1ZRjr8p+KASABKA'
    'lSBXZhbHVl');

@$core.Deprecated('Use metricDescriptor instead')
const Metric$json = {
  '1': 'Metric',
  '2': [
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
  ],
};

/// Descriptor for `Metric`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricDescriptor = $convert.base64Decode(
    'CgZNZXRyaWMSOQoKZGltZW5zaW9ucxjRm9/cASADKAsyFS5tb25pdG9yaW5nLkRpbWVuc2lvbl'
    'IKZGltZW5zaW9ucxIhCgptZXRyaWNuYW1lGPu+2jIgASgJUgptZXRyaWNuYW1lEiAKCW5hbWVz'
    'cGFjZRjBhLmpASABKAlSCW5hbWVzcGFjZQ==');

@$core.Deprecated('Use metricAlarmDescriptor instead')
const MetricAlarm$json = {
  '1': 'MetricAlarm',
  '2': [
    {
      '1': 'actionsenabled',
      '3': 126027878,
      '4': 1,
      '5': 8,
      '10': 'actionsenabled'
    },
    {'1': 'alarmactions', '3': 355779506, '4': 3, '5': 9, '10': 'alarmactions'},
    {'1': 'alarmarn', '3': 276019462, '4': 1, '5': 9, '10': 'alarmarn'},
    {
      '1': 'alarmconfigurationupdatedtimestamp',
      '3': 51935494,
      '4': 1,
      '5': 9,
      '10': 'alarmconfigurationupdatedtimestamp'
    },
    {
      '1': 'alarmdescription',
      '3': 30583613,
      '4': 1,
      '5': 9,
      '10': 'alarmdescription'
    },
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {
      '1': 'comparisonoperator',
      '3': 7059861,
      '4': 1,
      '5': 14,
      '6': '.monitoring.ComparisonOperator',
      '10': 'comparisonoperator'
    },
    {
      '1': 'datapointstoalarm',
      '3': 146378889,
      '4': 1,
      '5': 5,
      '10': 'datapointstoalarm'
    },
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {
      '1': 'evaluatelowsamplecountpercentile',
      '3': 64366323,
      '4': 1,
      '5': 9,
      '10': 'evaluatelowsamplecountpercentile'
    },
    {
      '1': 'evaluationperiods',
      '3': 214794856,
      '4': 1,
      '5': 5,
      '10': 'evaluationperiods'
    },
    {
      '1': 'evaluationstate',
      '3': 338863857,
      '4': 1,
      '5': 14,
      '6': '.monitoring.EvaluationState',
      '10': 'evaluationstate'
    },
    {
      '1': 'extendedstatistic',
      '3': 285620763,
      '4': 1,
      '5': 9,
      '10': 'extendedstatistic'
    },
    {
      '1': 'insufficientdataactions',
      '3': 498450778,
      '4': 3,
      '5': 9,
      '10': 'insufficientdataactions'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {
      '1': 'metrics',
      '3': 436365847,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricDataQuery',
      '10': 'metrics'
    },
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {'1': 'okactions', '3': 377443763, '4': 3, '5': 9, '10': 'okactions'},
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
    {
      '1': 'statereasondata',
      '3': 262771075,
      '4': 1,
      '5': 9,
      '10': 'statereasondata'
    },
    {
      '1': 'statetransitionedtimestamp',
      '3': 46811093,
      '4': 1,
      '5': 9,
      '10': 'statetransitionedtimestamp'
    },
    {
      '1': 'stateupdatedtimestamp',
      '3': 375406848,
      '4': 1,
      '5': 9,
      '10': 'stateupdatedtimestamp'
    },
    {
      '1': 'statevalue',
      '3': 334526008,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StateValue',
      '10': 'statevalue'
    },
    {
      '1': 'statistic',
      '3': 67293470,
      '4': 1,
      '5': 14,
      '6': '.monitoring.Statistic',
      '10': 'statistic'
    },
    {'1': 'threshold', '3': 394884505, '4': 1, '5': 1, '10': 'threshold'},
    {
      '1': 'thresholdmetricid',
      '3': 41342194,
      '4': 1,
      '5': 9,
      '10': 'thresholdmetricid'
    },
    {
      '1': 'treatmissingdata',
      '3': 226418150,
      '4': 1,
      '5': 9,
      '10': 'treatmissingdata'
    },
    {
      '1': 'unit',
      '3': 148989480,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StandardUnit',
      '10': 'unit'
    },
  ],
};

/// Descriptor for `MetricAlarm`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricAlarmDescriptor = $convert.base64Decode(
    'CgtNZXRyaWNBbGFybRIpCg5hY3Rpb25zZW5hYmxlZBjmkIw8IAEoCFIOYWN0aW9uc2VuYWJsZW'
    'QSJgoMYWxhcm1hY3Rpb25zGLKH06kBIAMoCVIMYWxhcm1hY3Rpb25zEh4KCGFsYXJtYXJuGIby'
    'zoMBIAEoCVIIYWxhcm1hcm4SUQoiYWxhcm1jb25maWd1cmF0aW9udXBkYXRlZHRpbWVzdGFtcB'
    'iG8uEYIAEoCVIiYWxhcm1jb25maWd1cmF0aW9udXBkYXRlZHRpbWVzdGFtcBItChBhbGFybWRl'
    'c2NyaXB0aW9uGL3Wyg4gASgJUhBhbGFybWRlc2NyaXB0aW9uEh8KCWFsYXJtbmFtZRji3+UZIA'
    'EoCVIJYWxhcm1uYW1lElEKEmNvbXBhcmlzb25vcGVyYXRvchiV864DIAEoDjIeLm1vbml0b3Jp'
    'bmcuQ29tcGFyaXNvbk9wZXJhdG9yUhJjb21wYXJpc29ub3BlcmF0b3ISLwoRZGF0YXBvaW50c3'
    'RvYWxhcm0YiaHmRSABKAVSEWRhdGFwb2ludHN0b2FsYXJtEjkKCmRpbWVuc2lvbnMY0Zvf3AEg'
    'AygLMhUubW9uaXRvcmluZy5EaW1lbnNpb25SCmRpbWVuc2lvbnMSTQogZXZhbHVhdGVsb3dzYW'
    '1wbGVjb3VudHBlcmNlbnRpbGUY883YHiABKAlSIGV2YWx1YXRlbG93c2FtcGxlY291bnRwZXJj'
    'ZW50aWxlEi8KEWV2YWx1YXRpb25wZXJpb2RzGOiEtmYgASgFUhFldmFsdWF0aW9ucGVyaW9kcx'
    'JJCg9ldmFsdWF0aW9uc3RhdGUY8c3KoQEgASgOMhsubW9uaXRvcmluZy5FdmFsdWF0aW9uU3Rh'
    'dGVSD2V2YWx1YXRpb25zdGF0ZRIwChFleHRlbmRlZHN0YXRpc3RpYxib9JiIASABKAlSEWV4dG'
    'VuZGVkc3RhdGlzdGljEjwKF2luc3VmZmljaWVudGRhdGFhY3Rpb25zGNqC1+0BIAMoCVIXaW5z'
    'dWZmaWNpZW50ZGF0YWFjdGlvbnMSIQoKbWV0cmljbmFtZRj7vtoyIAEoCVIKbWV0cmljbmFtZR'
    'I5CgdtZXRyaWNzGJfUidABIAMoCzIbLm1vbml0b3JpbmcuTWV0cmljRGF0YVF1ZXJ5UgdtZXRy'
    'aWNzEiAKCW5hbWVzcGFjZRjBhLmpASABKAlSCW5hbWVzcGFjZRIgCglva2FjdGlvbnMYs6v9sw'
    'EgAygJUglva2FjdGlvbnMSGQoGcGVyaW9kGKWIkjkgASgFUgZwZXJpb2QSJAoLc3RhdGVyZWFz'
    'b24Y89WtswEgASgJUgtzdGF0ZXJlYXNvbhIrCg9zdGF0ZXJlYXNvbmRhdGEYg6OmfSABKAlSD3'
    'N0YXRlcmVhc29uZGF0YRJBChpzdGF0ZXRyYW5zaXRpb25lZHRpbWVzdGFtcBjVj6kWIAEoCVIa'
    'c3RhdGV0cmFuc2l0aW9uZWR0aW1lc3RhbXASOAoVc3RhdGV1cGRhdGVkdGltZXN0YW1wGICCgb'
    'MBIAEoCVIVc3RhdGV1cGRhdGVkdGltZXN0YW1wEjoKCnN0YXRldmFsdWUYuOzBnwEgASgOMhYu'
    'bW9uaXRvcmluZy5TdGF0ZVZhbHVlUgpzdGF0ZXZhbHVlEjYKCXN0YXRpc3RpYxieoosgIAEoDj'
    'IVLm1vbml0b3JpbmcuU3RhdGlzdGljUglzdGF0aXN0aWMSIAoJdGhyZXNob2xkGJnrpbwBIAEo'
    'AVIJdGhyZXNob2xkEi8KEXRocmVzaG9sZG1ldHJpY2lkGPKp2xMgASgJUhF0aHJlc2hvbGRtZX'
    'RyaWNpZBItChB0cmVhdG1pc3NpbmdkYXRhGOa7+2sgASgJUhB0cmVhdG1pc3NpbmdkYXRhEi8K'
    'BHVuaXQYqMyFRyABKA4yGC5tb25pdG9yaW5nLlN0YW5kYXJkVW5pdFIEdW5pdA==');

@$core.Deprecated('Use metricCharacteristicsDescriptor instead')
const MetricCharacteristics$json = {
  '1': 'MetricCharacteristics',
  '2': [
    {
      '1': 'periodicspikes',
      '3': 454095604,
      '4': 1,
      '5': 8,
      '10': 'periodicspikes'
    },
  ],
};

/// Descriptor for `MetricCharacteristics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricCharacteristicsDescriptor = $convert.base64Decode(
    'ChVNZXRyaWNDaGFyYWN0ZXJpc3RpY3MSKgoOcGVyaW9kaWNzcGlrZXMY9OXD2AEgASgIUg5wZX'
    'Jpb2RpY3NwaWtlcw==');

@$core.Deprecated('Use metricDataQueryDescriptor instead')
const MetricDataQuery$json = {
  '1': 'MetricDataQuery',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {'1': 'expression', '3': 193051916, '4': 1, '5': 9, '10': 'expression'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'label', '3': 516747934, '4': 1, '5': 9, '10': 'label'},
    {
      '1': 'metricstat',
      '3': 167166628,
      '4': 1,
      '5': 11,
      '6': '.monitoring.MetricStat',
      '10': 'metricstat'
    },
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {'1': 'returndata', '3': 512985704, '4': 1, '5': 8, '10': 'returndata'},
  ],
};

/// Descriptor for `MetricDataQuery`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricDataQueryDescriptor = $convert.base64Decode(
    'Cg9NZXRyaWNEYXRhUXVlcnkSHwoJYWNjb3VudGlkGNLBuR8gASgJUglhY2NvdW50aWQSIQoKZX'
    'hwcmVzc2lvbhiM+oZcIAEoCVIKZXhwcmVzc2lvbhISCgJpZBiB8qK3ASABKAlSAmlkEhgKBWxh'
    'YmVsGJ7ls/YBIAEoCVIFbGFiZWwSOQoKbWV0cmljc3RhdBikhdtPIAEoCzIWLm1vbml0b3Jpbm'
    'cuTWV0cmljU3RhdFIKbWV0cmljc3RhdBIZCgZwZXJpb2QYpYiSOSABKAVSBnBlcmlvZBIiCgpy'
    'ZXR1cm5kYXRhGOiUzvQBIAEoCFIKcmV0dXJuZGF0YQ==');

@$core.Deprecated('Use metricDataResultDescriptor instead')
const MetricDataResult$json = {
  '1': 'MetricDataResult',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'label', '3': 516747934, '4': 1, '5': 9, '10': 'label'},
    {
      '1': 'messages',
      '3': 409018326,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MessageData',
      '10': 'messages'
    },
    {
      '1': 'statuscode',
      '3': 303830783,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StatusCode',
      '10': 'statuscode'
    },
    {'1': 'timestamps', '3': 213534737, '4': 3, '5': 9, '10': 'timestamps'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 1, '10': 'values'},
  ],
};

/// Descriptor for `MetricDataResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricDataResultDescriptor = $convert.base64Decode(
    'ChBNZXRyaWNEYXRhUmVzdWx0EhIKAmlkGIHyorcBIAEoCVICaWQSGAoFbGFiZWwYnuWz9gEgAS'
    'gJUgVsYWJlbBI3CghtZXNzYWdlcxjWv4TDASADKAsyFy5tb25pdG9yaW5nLk1lc3NhZ2VEYXRh'
    'UghtZXNzYWdlcxI6CgpzdGF0dXNjb2RlGP+t8JABIAEoDjIWLm1vbml0b3JpbmcuU3RhdHVzQ2'
    '9kZVIKc3RhdHVzY29kZRIhCgp0aW1lc3RhbXBzGJGQ6WUgAygJUgp0aW1lc3RhbXBzEhkKBnZh'
    'bHVlcxjcxLRqIAMoAVIGdmFsdWVz');

@$core.Deprecated('Use metricDatumDescriptor instead')
const MetricDatum$json = {
  '1': 'MetricDatum',
  '2': [
    {'1': 'counts', '3': 113775526, '4': 3, '5': 1, '10': 'counts'},
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {
      '1': 'statisticvalues',
      '3': 92623908,
      '4': 1,
      '5': 11,
      '6': '.monitoring.StatisticSet',
      '10': 'statisticvalues'
    },
    {
      '1': 'storageresolution',
      '3': 194430293,
      '4': 1,
      '5': 5,
      '10': 'storageresolution'
    },
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
    {
      '1': 'unit',
      '3': 148989480,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StandardUnit',
      '10': 'unit'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 1, '10': 'value'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 1, '10': 'values'},
  ],
};

/// Descriptor for `MetricDatum`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricDatumDescriptor = $convert.base64Decode(
    'CgtNZXRyaWNEYXR1bRIZCgZjb3VudHMYpqegNiADKAFSBmNvdW50cxI5CgpkaW1lbnNpb25zGN'
    'Gb39wBIAMoCzIVLm1vbml0b3JpbmcuRGltZW5zaW9uUgpkaW1lbnNpb25zEiEKCm1ldHJpY25h'
    'bWUY+77aMiABKAlSCm1ldHJpY25hbWUSRQoPc3RhdGlzdGljdmFsdWVzGKSolSwgASgLMhgubW'
    '9uaXRvcmluZy5TdGF0aXN0aWNTZXRSD3N0YXRpc3RpY3ZhbHVlcxIvChFzdG9yYWdlcmVzb2x1'
    'dGlvbhjVittcIAEoBVIRc3RvcmFnZXJlc29sdXRpb24SHwoJdGltZXN0YW1wGMTDt00gASgJUg'
    'l0aW1lc3RhbXASLwoEdW5pdBiozIVHIAEoDjIYLm1vbml0b3JpbmcuU3RhbmRhcmRVbml0UgR1'
    'bml0EhgKBXZhbHVlGOvyn4oBIAEoAVIFdmFsdWUSGQoGdmFsdWVzGNzEtGogAygBUgZ2YWx1ZX'
    'M=');

@$core.Deprecated('Use metricMathAnomalyDetectorDescriptor instead')
const MetricMathAnomalyDetector$json = {
  '1': 'MetricMathAnomalyDetector',
  '2': [
    {
      '1': 'metricdataqueries',
      '3': 92209504,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricDataQuery',
      '10': 'metricdataqueries'
    },
  ],
};

/// Descriptor for `MetricMathAnomalyDetector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricMathAnomalyDetectorDescriptor =
    $convert.base64Decode(
        'ChlNZXRyaWNNYXRoQW5vbWFseURldGVjdG9yEkwKEW1ldHJpY2RhdGFxdWVyaWVzGOCC/CsgAy'
        'gLMhsubW9uaXRvcmluZy5NZXRyaWNEYXRhUXVlcnlSEW1ldHJpY2RhdGFxdWVyaWVz');

@$core.Deprecated('Use metricStatDescriptor instead')
const MetricStat$json = {
  '1': 'MetricStat',
  '2': [
    {
      '1': 'metric',
      '3': 386667298,
      '4': 1,
      '5': 11,
      '6': '.monitoring.Metric',
      '10': 'metric'
    },
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {'1': 'stat', '3': 325549752, '4': 1, '5': 9, '10': 'stat'},
    {
      '1': 'unit',
      '3': 148989480,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StandardUnit',
      '10': 'unit'
    },
  ],
};

/// Descriptor for `MetricStat`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricStatDescriptor = $convert.base64Decode(
    'CgpNZXRyaWNTdGF0Ei4KBm1ldHJpYxiiprC4ASABKAsyEi5tb25pdG9yaW5nLk1ldHJpY1IGbW'
    'V0cmljEhkKBnBlcmlvZBiliJI5IAEoBVIGcGVyaW9kEhYKBHN0YXQYuP2dmwEgASgJUgRzdGF0'
    'Ei8KBHVuaXQYqMyFRyABKA4yGC5tb25pdG9yaW5nLlN0YW5kYXJkVW5pdFIEdW5pdA==');

@$core.Deprecated('Use metricStreamEntryDescriptor instead')
const MetricStreamEntry$json = {
  '1': 'MetricStreamEntry',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'firehosearn', '3': 55494200, '4': 1, '5': 9, '10': 'firehosearn'},
    {
      '1': 'lastupdatedate',
      '3': 65755725,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputformat',
      '3': 519139704,
      '4': 1,
      '5': 14,
      '6': '.monitoring.MetricStreamOutputFormat',
      '10': 'outputformat'
    },
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
  ],
};

/// Descriptor for `MetricStreamEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricStreamEntryDescriptor = $convert.base64Decode(
    'ChFNZXRyaWNTdHJlYW1FbnRyeRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SJgoMY3JlYXRpb25kYX'
    'RlGOHYt4kBIAEoCVIMY3JlYXRpb25kYXRlEiMKC2ZpcmVob3NlYXJuGLiMuxogASgJUgtmaXJl'
    'aG9zZWFybhIpCg5sYXN0dXBkYXRlZGF0ZRjNtK0fIAEoCVIObGFzdHVwZGF0ZWRhdGUSFQoEbm'
    'FtZRiH5oF/IAEoCVIEbmFtZRJMCgxvdXRwdXRmb3JtYXQY+OLF9wEgASgOMiQubW9uaXRvcmlu'
    'Zy5NZXRyaWNTdHJlYW1PdXRwdXRGb3JtYXRSDG91dHB1dGZvcm1hdBIYCgVzdGF0ZRiXybLvAS'
    'ABKAlSBXN0YXRl');

@$core.Deprecated('Use metricStreamFilterDescriptor instead')
const MetricStreamFilter$json = {
  '1': 'MetricStreamFilter',
  '2': [
    {'1': 'metricnames', '3': 56384300, '4': 3, '5': 9, '10': 'metricnames'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
  ],
};

/// Descriptor for `MetricStreamFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricStreamFilterDescriptor = $convert.base64Decode(
    'ChJNZXRyaWNTdHJlYW1GaWx0ZXISIwoLbWV0cmljbmFtZXMYrLbxGiADKAlSC21ldHJpY25hbW'
    'VzEiAKCW5hbWVzcGFjZRjBhLmpASABKAlSCW5hbWVzcGFjZQ==');

@$core.Deprecated('Use metricStreamStatisticsConfigurationDescriptor instead')
const MetricStreamStatisticsConfiguration$json = {
  '1': 'MetricStreamStatisticsConfiguration',
  '2': [
    {
      '1': 'additionalstatistics',
      '3': 50823856,
      '4': 3,
      '5': 9,
      '10': 'additionalstatistics'
    },
    {
      '1': 'includemetrics',
      '3': 169287849,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricStreamStatisticsMetric',
      '10': 'includemetrics'
    },
  ],
};

/// Descriptor for `MetricStreamStatisticsConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricStreamStatisticsConfigurationDescriptor =
    $convert.base64Decode(
        'CiNNZXRyaWNTdHJlYW1TdGF0aXN0aWNzQ29uZmlndXJhdGlvbhI1ChRhZGRpdGlvbmFsc3RhdG'
        'lzdGljcxiwhZ4YIAMoCVIUYWRkaXRpb25hbHN0YXRpc3RpY3MSUwoOaW5jbHVkZW1ldHJpY3MY'
        'qcHcUCADKAsyKC5tb25pdG9yaW5nLk1ldHJpY1N0cmVhbVN0YXRpc3RpY3NNZXRyaWNSDmluY2'
        'x1ZGVtZXRyaWNz');

@$core.Deprecated('Use metricStreamStatisticsMetricDescriptor instead')
const MetricStreamStatisticsMetric$json = {
  '1': 'MetricStreamStatisticsMetric',
  '2': [
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
  ],
};

/// Descriptor for `MetricStreamStatisticsMetric`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricStreamStatisticsMetricDescriptor =
    $convert.base64Decode(
        'ChxNZXRyaWNTdHJlYW1TdGF0aXN0aWNzTWV0cmljEiEKCm1ldHJpY25hbWUY+77aMiABKAlSCm'
        '1ldHJpY25hbWUSIAoJbmFtZXNwYWNlGMGEuakBIAEoCVIJbmFtZXNwYWNl');

@$core.Deprecated('Use missingRequiredParameterExceptionDescriptor instead')
const MissingRequiredParameterException$json = {
  '1': 'MissingRequiredParameterException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MissingRequiredParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List missingRequiredParameterExceptionDescriptor =
    $convert.base64Decode(
        'CiFNaXNzaW5nUmVxdWlyZWRQYXJhbWV0ZXJFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use muteTargetsDescriptor instead')
const MuteTargets$json = {
  '1': 'MuteTargets',
  '2': [
    {'1': 'alarmnames', '3': 90874583, '4': 3, '5': 9, '10': 'alarmnames'},
  ],
};

/// Descriptor for `MuteTargets`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List muteTargetsDescriptor = $convert.base64Decode(
    'CgtNdXRlVGFyZ2V0cxIhCgphbGFybW5hbWVzGNfFqisgAygJUgphbGFybW5hbWVz');

@$core.Deprecated('Use partialFailureDescriptor instead')
const PartialFailure$json = {
  '1': 'PartialFailure',
  '2': [
    {
      '1': 'exceptiontype',
      '3': 176478087,
      '4': 1,
      '5': 9,
      '10': 'exceptiontype'
    },
    {'1': 'failurecode', '3': 84707897, '4': 1, '5': 9, '10': 'failurecode'},
    {
      '1': 'failuredescription',
      '3': 249093686,
      '4': 1,
      '5': 9,
      '10': 'failuredescription'
    },
    {
      '1': 'failureresource',
      '3': 315801394,
      '4': 1,
      '5': 9,
      '10': 'failureresource'
    },
  ],
};

/// Descriptor for `PartialFailure`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List partialFailureDescriptor = $convert.base64Decode(
    'Cg5QYXJ0aWFsRmFpbHVyZRInCg1leGNlcHRpb250eXBlGIevk1QgASgJUg1leGNlcHRpb250eX'
    'BlEiMKC2ZhaWx1cmVjb2RlGLmUsiggASgJUgtmYWlsdXJlY29kZRIxChJmYWlsdXJlZGVzY3Jp'
    'cHRpb24YtrzjdiABKAlSEmZhaWx1cmVkZXNjcmlwdGlvbhIsCg9mYWlsdXJlcmVzb3VyY2UYsv'
    '7KlgEgASgJUg9mYWlsdXJlcmVzb3VyY2U=');

@$core.Deprecated('Use putAlarmMuteRuleInputDescriptor instead')
const PutAlarmMuteRuleInput$json = {
  '1': 'PutAlarmMuteRuleInput',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'expiredate', '3': 454143207, '4': 1, '5': 9, '10': 'expiredate'},
    {
      '1': 'mutetargets',
      '3': 356717131,
      '4': 1,
      '5': 11,
      '6': '.monitoring.MuteTargets',
      '10': 'mutetargets'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'rule',
      '3': 475696372,
      '4': 1,
      '5': 11,
      '6': '.monitoring.Rule',
      '10': 'rule'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `PutAlarmMuteRuleInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAlarmMuteRuleInputDescriptor = $convert.base64Decode(
    'ChVQdXRBbGFybU11dGVSdWxlSW5wdXQSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaX'
    'B0aW9uEiIKCmV4cGlyZWRhdGUY59nG2AEgASgJUgpleHBpcmVkYXRlEj0KC211dGV0YXJnZXRz'
    'GMukjKoBIAEoCzIXLm1vbml0b3JpbmcuTXV0ZVRhcmdldHNSC211dGV0YXJnZXRzEhUKBG5hbW'
    'UYh+aBfyABKAlSBG5hbWUSKAoEcnVsZRj0meriASABKAsyEC5tb25pdG9yaW5nLlJ1bGVSBHJ1'
    'bGUSIAoJc3RhcnRkYXRlGPz4oNQBIAEoCVIJc3RhcnRkYXRlEicKBHRhZ3MYwcH2tQEgAygLMg'
    '8ubW9uaXRvcmluZy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use putAnomalyDetectorInputDescriptor instead')
const PutAnomalyDetectorInput$json = {
  '1': 'PutAnomalyDetectorInput',
  '2': [
    {
      '1': 'configuration',
      '3': 442426458,
      '4': 1,
      '5': 11,
      '6': '.monitoring.AnomalyDetectorConfiguration',
      '10': 'configuration'
    },
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {
      '1': 'metriccharacteristics',
      '3': 129404842,
      '4': 1,
      '5': 11,
      '6': '.monitoring.MetricCharacteristics',
      '10': 'metriccharacteristics'
    },
    {
      '1': 'metricmathanomalydetector',
      '3': 439381475,
      '4': 1,
      '5': 11,
      '6': '.monitoring.MetricMathAnomalyDetector',
      '10': 'metricmathanomalydetector'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {
      '1': 'singlemetricanomalydetector',
      '3': 418531933,
      '4': 1,
      '5': 11,
      '6': '.monitoring.SingleMetricAnomalyDetector',
      '10': 'singlemetricanomalydetector'
    },
    {'1': 'stat', '3': 325549752, '4': 1, '5': 9, '10': 'stat'},
  ],
};

/// Descriptor for `PutAnomalyDetectorInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAnomalyDetectorInputDescriptor = $convert.base64Decode(
    'ChdQdXRBbm9tYWx5RGV0ZWN0b3JJbnB1dBJSCg1jb25maWd1cmF0aW9uGNrI+9IBIAEoCzIoLm'
    '1vbml0b3JpbmcuQW5vbWFseURldGVjdG9yQ29uZmlndXJhdGlvblINY29uZmlndXJhdGlvbhI5'
    'CgpkaW1lbnNpb25zGNGb39wBIAMoCzIVLm1vbml0b3JpbmcuRGltZW5zaW9uUgpkaW1lbnNpb2'
    '5zEloKFW1ldHJpY2NoYXJhY3RlcmlzdGljcxiqn9o9IAEoCzIhLm1vbml0b3JpbmcuTWV0cmlj'
    'Q2hhcmFjdGVyaXN0aWNzUhVtZXRyaWNjaGFyYWN0ZXJpc3RpY3MSZwoZbWV0cmljbWF0aGFub2'
    '1hbHlkZXRlY3Rvchjj28HRASABKAsyJS5tb25pdG9yaW5nLk1ldHJpY01hdGhBbm9tYWx5RGV0'
    'ZWN0b3JSGW1ldHJpY21hdGhhbm9tYWx5ZGV0ZWN0b3ISIQoKbWV0cmljbmFtZRj7vtoyIAEoCV'
    'IKbWV0cmljbmFtZRIgCgluYW1lc3BhY2UYwYS5qQEgASgJUgluYW1lc3BhY2USbQobc2luZ2xl'
    'bWV0cmljYW5vbWFseWRldGVjdG9yGN2UyccBIAEoCzInLm1vbml0b3JpbmcuU2luZ2xlTWV0cm'
    'ljQW5vbWFseURldGVjdG9yUhtzaW5nbGVtZXRyaWNhbm9tYWx5ZGV0ZWN0b3ISFgoEc3RhdBi4'
    '/Z2bASABKAlSBHN0YXQ=');

@$core.Deprecated('Use putAnomalyDetectorOutputDescriptor instead')
const PutAnomalyDetectorOutput$json = {
  '1': 'PutAnomalyDetectorOutput',
};

/// Descriptor for `PutAnomalyDetectorOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAnomalyDetectorOutputDescriptor =
    $convert.base64Decode('ChhQdXRBbm9tYWx5RGV0ZWN0b3JPdXRwdXQ=');

@$core.Deprecated('Use putCompositeAlarmInputDescriptor instead')
const PutCompositeAlarmInput$json = {
  '1': 'PutCompositeAlarmInput',
  '2': [
    {
      '1': 'actionsenabled',
      '3': 126027878,
      '4': 1,
      '5': 8,
      '10': 'actionsenabled'
    },
    {
      '1': 'actionssuppressor',
      '3': 486417535,
      '4': 1,
      '5': 9,
      '10': 'actionssuppressor'
    },
    {
      '1': 'actionssuppressorextensionperiod',
      '3': 183723543,
      '4': 1,
      '5': 5,
      '10': 'actionssuppressorextensionperiod'
    },
    {
      '1': 'actionssuppressorwaitperiod',
      '3': 80227275,
      '4': 1,
      '5': 5,
      '10': 'actionssuppressorwaitperiod'
    },
    {'1': 'alarmactions', '3': 355779506, '4': 3, '5': 9, '10': 'alarmactions'},
    {
      '1': 'alarmdescription',
      '3': 30583613,
      '4': 1,
      '5': 9,
      '10': 'alarmdescription'
    },
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {'1': 'alarmrule', '3': 516996973, '4': 1, '5': 9, '10': 'alarmrule'},
    {
      '1': 'insufficientdataactions',
      '3': 498450778,
      '4': 3,
      '5': 9,
      '10': 'insufficientdataactions'
    },
    {'1': 'okactions', '3': 377443763, '4': 3, '5': 9, '10': 'okactions'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `PutCompositeAlarmInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putCompositeAlarmInputDescriptor = $convert.base64Decode(
    'ChZQdXRDb21wb3NpdGVBbGFybUlucHV0EikKDmFjdGlvbnNlbmFibGVkGOaQjDwgASgIUg5hY3'
    'Rpb25zZW5hYmxlZBIwChFhY3Rpb25zc3VwcHJlc3Nvchj/yPjnASABKAlSEWFjdGlvbnNzdXBw'
    'cmVzc29yEk0KIGFjdGlvbnNzdXBwcmVzc29yZXh0ZW5zaW9ucGVyaW9kGJfMzVcgASgFUiBhY3'
    'Rpb25zc3VwcHJlc3NvcmV4dGVuc2lvbnBlcmlvZBJDChthY3Rpb25zc3VwcHJlc3NvcndhaXRw'
    'ZXJpb2QYy9egJiABKAVSG2FjdGlvbnNzdXBwcmVzc29yd2FpdHBlcmlvZBImCgxhbGFybWFjdG'
    'lvbnMYsofTqQEgAygJUgxhbGFybWFjdGlvbnMSLQoQYWxhcm1kZXNjcmlwdGlvbhi91soOIAEo'
    'CVIQYWxhcm1kZXNjcmlwdGlvbhIfCglhbGFybW5hbWUY4t/lGSABKAlSCWFsYXJtbmFtZRIgCg'
    'lhbGFybXJ1bGUY7f7C9gEgASgJUglhbGFybXJ1bGUSPAoXaW5zdWZmaWNpZW50ZGF0YWFjdGlv'
    'bnMY2oLX7QEgAygJUhdpbnN1ZmZpY2llbnRkYXRhYWN0aW9ucxIgCglva2FjdGlvbnMYs6v9sw'
    'EgAygJUglva2FjdGlvbnMSJwoEdGFncxjBwfa1ASADKAsyDy5tb25pdG9yaW5nLlRhZ1IEdGFn'
    'cw==');

@$core.Deprecated('Use putDashboardInputDescriptor instead')
const PutDashboardInput$json = {
  '1': 'PutDashboardInput',
  '2': [
    {'1': 'dashboardbody', '3': 3403236, '4': 1, '5': 9, '10': 'dashboardbody'},
    {
      '1': 'dashboardname',
      '3': 506599873,
      '4': 1,
      '5': 9,
      '10': 'dashboardname'
    },
  ],
};

/// Descriptor for `PutDashboardInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDashboardInputDescriptor = $convert.base64Decode(
    'ChFQdXREYXNoYm9hcmRJbnB1dBInCg1kYXNoYm9hcmRib2R5GOTbzwEgASgJUg1kYXNoYm9hcm'
    'Rib2R5EigKDWRhc2hib2FyZG5hbWUYwbPI8QEgASgJUg1kYXNoYm9hcmRuYW1l');

@$core.Deprecated('Use putDashboardOutputDescriptor instead')
const PutDashboardOutput$json = {
  '1': 'PutDashboardOutput',
  '2': [
    {
      '1': 'dashboardvalidationmessages',
      '3': 515270089,
      '4': 3,
      '5': 11,
      '6': '.monitoring.DashboardValidationMessage',
      '10': 'dashboardvalidationmessages'
    },
  ],
};

/// Descriptor for `PutDashboardOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDashboardOutputDescriptor = $convert.base64Decode(
    'ChJQdXREYXNoYm9hcmRPdXRwdXQSbAobZGFzaGJvYXJkdmFsaWRhdGlvbm1lc3NhZ2VzGMnL2f'
    'UBIAMoCzImLm1vbml0b3JpbmcuRGFzaGJvYXJkVmFsaWRhdGlvbk1lc3NhZ2VSG2Rhc2hib2Fy'
    'ZHZhbGlkYXRpb25tZXNzYWdlcw==');

@$core.Deprecated('Use putInsightRuleInputDescriptor instead')
const PutInsightRuleInput$json = {
  '1': 'PutInsightRuleInput',
  '2': [
    {
      '1': 'applyontransformedlogs',
      '3': 30966533,
      '4': 1,
      '5': 8,
      '10': 'applyontransformedlogs'
    },
    {
      '1': 'ruledefinition',
      '3': 188932559,
      '4': 1,
      '5': 9,
      '10': 'ruledefinition'
    },
    {'1': 'rulename', '3': 214688793, '4': 1, '5': 9, '10': 'rulename'},
    {'1': 'rulestate', '3': 419921189, '4': 1, '5': 9, '10': 'rulestate'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `PutInsightRuleInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putInsightRuleInputDescriptor = $convert.base64Decode(
    'ChNQdXRJbnNpZ2h0UnVsZUlucHV0EjkKFmFwcGx5b250cmFuc2Zvcm1lZGxvZ3MYhYbiDiABKA'
    'hSFmFwcGx5b250cmFuc2Zvcm1lZGxvZ3MSKQoOcnVsZWRlZmluaXRpb24Yz8OLWiABKAlSDnJ1'
    'bGVkZWZpbml0aW9uEh0KCHJ1bGVuYW1lGJnIr2YgASgJUghydWxlbmFtZRIgCglydWxlc3RhdG'
    'UYpfqdyAEgASgJUglydWxlc3RhdGUSJwoEdGFncxjBwfa1ASADKAsyDy5tb25pdG9yaW5nLlRh'
    'Z1IEdGFncw==');

@$core.Deprecated('Use putInsightRuleOutputDescriptor instead')
const PutInsightRuleOutput$json = {
  '1': 'PutInsightRuleOutput',
};

/// Descriptor for `PutInsightRuleOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putInsightRuleOutputDescriptor =
    $convert.base64Decode('ChRQdXRJbnNpZ2h0UnVsZU91dHB1dA==');

@$core.Deprecated('Use putManagedInsightRulesInputDescriptor instead')
const PutManagedInsightRulesInput$json = {
  '1': 'PutManagedInsightRulesInput',
  '2': [
    {
      '1': 'managedrules',
      '3': 347090906,
      '4': 3,
      '5': 11,
      '6': '.monitoring.ManagedRule',
      '10': 'managedrules'
    },
  ],
};

/// Descriptor for `PutManagedInsightRulesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putManagedInsightRulesInputDescriptor =
    $convert.base64Decode(
        'ChtQdXRNYW5hZ2VkSW5zaWdodFJ1bGVzSW5wdXQSPwoMbWFuYWdlZHJ1bGVzGNrfwKUBIAMoCz'
        'IXLm1vbml0b3JpbmcuTWFuYWdlZFJ1bGVSDG1hbmFnZWRydWxlcw==');

@$core.Deprecated('Use putManagedInsightRulesOutputDescriptor instead')
const PutManagedInsightRulesOutput$json = {
  '1': 'PutManagedInsightRulesOutput',
  '2': [
    {
      '1': 'failures',
      '3': 335467271,
      '4': 3,
      '5': 11,
      '6': '.monitoring.PartialFailure',
      '10': 'failures'
    },
  ],
};

/// Descriptor for `PutManagedInsightRulesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putManagedInsightRulesOutputDescriptor =
    $convert.base64Decode(
        'ChxQdXRNYW5hZ2VkSW5zaWdodFJ1bGVzT3V0cHV0EjoKCGZhaWx1cmVzGIem+58BIAMoCzIaLm'
        '1vbml0b3JpbmcuUGFydGlhbEZhaWx1cmVSCGZhaWx1cmVz');

@$core.Deprecated('Use putMetricAlarmInputDescriptor instead')
const PutMetricAlarmInput$json = {
  '1': 'PutMetricAlarmInput',
  '2': [
    {
      '1': 'actionsenabled',
      '3': 126027878,
      '4': 1,
      '5': 8,
      '10': 'actionsenabled'
    },
    {'1': 'alarmactions', '3': 355779506, '4': 3, '5': 9, '10': 'alarmactions'},
    {
      '1': 'alarmdescription',
      '3': 30583613,
      '4': 1,
      '5': 9,
      '10': 'alarmdescription'
    },
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {
      '1': 'comparisonoperator',
      '3': 7059861,
      '4': 1,
      '5': 14,
      '6': '.monitoring.ComparisonOperator',
      '10': 'comparisonoperator'
    },
    {
      '1': 'datapointstoalarm',
      '3': 146378889,
      '4': 1,
      '5': 5,
      '10': 'datapointstoalarm'
    },
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {
      '1': 'evaluatelowsamplecountpercentile',
      '3': 64366323,
      '4': 1,
      '5': 9,
      '10': 'evaluatelowsamplecountpercentile'
    },
    {
      '1': 'evaluationperiods',
      '3': 214794856,
      '4': 1,
      '5': 5,
      '10': 'evaluationperiods'
    },
    {
      '1': 'extendedstatistic',
      '3': 285620763,
      '4': 1,
      '5': 9,
      '10': 'extendedstatistic'
    },
    {
      '1': 'insufficientdataactions',
      '3': 498450778,
      '4': 3,
      '5': 9,
      '10': 'insufficientdataactions'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {
      '1': 'metrics',
      '3': 436365847,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricDataQuery',
      '10': 'metrics'
    },
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {'1': 'okactions', '3': 377443763, '4': 3, '5': 9, '10': 'okactions'},
    {'1': 'period', '3': 119833637, '4': 1, '5': 5, '10': 'period'},
    {
      '1': 'statistic',
      '3': 67293470,
      '4': 1,
      '5': 14,
      '6': '.monitoring.Statistic',
      '10': 'statistic'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Tag',
      '10': 'tags'
    },
    {'1': 'threshold', '3': 394884505, '4': 1, '5': 1, '10': 'threshold'},
    {
      '1': 'thresholdmetricid',
      '3': 41342194,
      '4': 1,
      '5': 9,
      '10': 'thresholdmetricid'
    },
    {
      '1': 'treatmissingdata',
      '3': 226418150,
      '4': 1,
      '5': 9,
      '10': 'treatmissingdata'
    },
    {
      '1': 'unit',
      '3': 148989480,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StandardUnit',
      '10': 'unit'
    },
  ],
};

/// Descriptor for `PutMetricAlarmInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putMetricAlarmInputDescriptor = $convert.base64Decode(
    'ChNQdXRNZXRyaWNBbGFybUlucHV0EikKDmFjdGlvbnNlbmFibGVkGOaQjDwgASgIUg5hY3Rpb2'
    '5zZW5hYmxlZBImCgxhbGFybWFjdGlvbnMYsofTqQEgAygJUgxhbGFybWFjdGlvbnMSLQoQYWxh'
    'cm1kZXNjcmlwdGlvbhi91soOIAEoCVIQYWxhcm1kZXNjcmlwdGlvbhIfCglhbGFybW5hbWUY4t'
    '/lGSABKAlSCWFsYXJtbmFtZRJRChJjb21wYXJpc29ub3BlcmF0b3IYlfOuAyABKA4yHi5tb25p'
    'dG9yaW5nLkNvbXBhcmlzb25PcGVyYXRvclISY29tcGFyaXNvbm9wZXJhdG9yEi8KEWRhdGFwb2'
    'ludHN0b2FsYXJtGImh5kUgASgFUhFkYXRhcG9pbnRzdG9hbGFybRI5CgpkaW1lbnNpb25zGNGb'
    '39wBIAMoCzIVLm1vbml0b3JpbmcuRGltZW5zaW9uUgpkaW1lbnNpb25zEk0KIGV2YWx1YXRlbG'
    '93c2FtcGxlY291bnRwZXJjZW50aWxlGPPN2B4gASgJUiBldmFsdWF0ZWxvd3NhbXBsZWNvdW50'
    'cGVyY2VudGlsZRIvChFldmFsdWF0aW9ucGVyaW9kcxjohLZmIAEoBVIRZXZhbHVhdGlvbnBlcm'
    'lvZHMSMAoRZXh0ZW5kZWRzdGF0aXN0aWMYm/SYiAEgASgJUhFleHRlbmRlZHN0YXRpc3RpYxI8'
    'ChdpbnN1ZmZpY2llbnRkYXRhYWN0aW9ucxjagtftASADKAlSF2luc3VmZmljaWVudGRhdGFhY3'
    'Rpb25zEiEKCm1ldHJpY25hbWUY+77aMiABKAlSCm1ldHJpY25hbWUSOQoHbWV0cmljcxiX1InQ'
    'ASADKAsyGy5tb25pdG9yaW5nLk1ldHJpY0RhdGFRdWVyeVIHbWV0cmljcxIgCgluYW1lc3BhY2'
    'UYwYS5qQEgASgJUgluYW1lc3BhY2USIAoJb2thY3Rpb25zGLOr/bMBIAMoCVIJb2thY3Rpb25z'
    'EhkKBnBlcmlvZBiliJI5IAEoBVIGcGVyaW9kEjYKCXN0YXRpc3RpYxieoosgIAEoDjIVLm1vbm'
    'l0b3JpbmcuU3RhdGlzdGljUglzdGF0aXN0aWMSJwoEdGFncxjBwfa1ASADKAsyDy5tb25pdG9y'
    'aW5nLlRhZ1IEdGFncxIgCgl0aHJlc2hvbGQYmeulvAEgASgBUgl0aHJlc2hvbGQSLwoRdGhyZX'
    'Nob2xkbWV0cmljaWQY8qnbEyABKAlSEXRocmVzaG9sZG1ldHJpY2lkEi0KEHRyZWF0bWlzc2lu'
    'Z2RhdGEY5rv7ayABKAlSEHRyZWF0bWlzc2luZ2RhdGESLwoEdW5pdBiozIVHIAEoDjIYLm1vbm'
    'l0b3JpbmcuU3RhbmRhcmRVbml0UgR1bml0');

@$core.Deprecated('Use putMetricDataInputDescriptor instead')
const PutMetricDataInput$json = {
  '1': 'PutMetricDataInput',
  '2': [
    {
      '1': 'entitymetricdata',
      '3': 305844415,
      '4': 3,
      '5': 11,
      '6': '.monitoring.EntityMetricData',
      '10': 'entitymetricdata'
    },
    {
      '1': 'metricdata',
      '3': 162960562,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricDatum',
      '10': 'metricdata'
    },
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {
      '1': 'strictentityvalidation',
      '3': 496663153,
      '4': 1,
      '5': 8,
      '10': 'strictentityvalidation'
    },
  ],
};

/// Descriptor for `PutMetricDataInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putMetricDataInputDescriptor = $convert.base64Decode(
    'ChJQdXRNZXRyaWNEYXRhSW5wdXQSTAoQZW50aXR5bWV0cmljZGF0YRi/oeuRASADKAsyHC5tb2'
    '5pdG9yaW5nLkVudGl0eU1ldHJpY0RhdGFSEGVudGl0eW1ldHJpY2RhdGESOgoKbWV0cmljZGF0'
    'YRiyqdpNIAMoCzIXLm1vbml0b3JpbmcuTWV0cmljRGF0dW1SCm1ldHJpY2RhdGESIAoJbmFtZX'
    'NwYWNlGMGEuakBIAEoCVIJbmFtZXNwYWNlEjoKFnN0cmljdGVudGl0eXZhbGlkYXRpb24Y8fTp'
    '7AEgASgIUhZzdHJpY3RlbnRpdHl2YWxpZGF0aW9u');

@$core.Deprecated('Use putMetricStreamInputDescriptor instead')
const PutMetricStreamInput$json = {
  '1': 'PutMetricStreamInput',
  '2': [
    {
      '1': 'excludefilters',
      '3': 426990841,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricStreamFilter',
      '10': 'excludefilters'
    },
    {'1': 'firehosearn', '3': 55494200, '4': 1, '5': 9, '10': 'firehosearn'},
    {
      '1': 'includefilters',
      '3': 90967031,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricStreamFilter',
      '10': 'includefilters'
    },
    {
      '1': 'includelinkedaccountsmetrics',
      '3': 356901952,
      '4': 1,
      '5': 8,
      '10': 'includelinkedaccountsmetrics'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputformat',
      '3': 519139704,
      '4': 1,
      '5': 14,
      '6': '.monitoring.MetricStreamOutputFormat',
      '10': 'outputformat'
    },
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'statisticsconfigurations',
      '3': 175876178,
      '4': 3,
      '5': 11,
      '6': '.monitoring.MetricStreamStatisticsConfiguration',
      '10': 'statisticsconfigurations'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `PutMetricStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putMetricStreamInputDescriptor = $convert.base64Decode(
    'ChRQdXRNZXRyaWNTdHJlYW1JbnB1dBJKCg5leGNsdWRlZmlsdGVycxj5uc3LASADKAsyHi5tb2'
    '5pdG9yaW5nLk1ldHJpY1N0cmVhbUZpbHRlclIOZXhjbHVkZWZpbHRlcnMSIwoLZmlyZWhvc2Vh'
    'cm4YuIy7GiABKAlSC2ZpcmVob3NlYXJuEkkKDmluY2x1ZGVmaWx0ZXJzGPeXsCsgAygLMh4ubW'
    '9uaXRvcmluZy5NZXRyaWNTdHJlYW1GaWx0ZXJSDmluY2x1ZGVmaWx0ZXJzEkYKHGluY2x1ZGVs'
    'aW5rZWRhY2NvdW50c21ldHJpY3MYwMiXqgEgASgIUhxpbmNsdWRlbGlua2VkYWNjb3VudHNtZX'
    'RyaWNzEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSTAoMb3V0cHV0Zm9ybWF0GPjixfcBIAEoDjIk'
    'Lm1vbml0b3JpbmcuTWV0cmljU3RyZWFtT3V0cHV0Rm9ybWF0UgxvdXRwdXRmb3JtYXQSHAoHcm'
    '9sZWFybhiB+OeZASABKAlSB3JvbGVhcm4SbgoYc3RhdGlzdGljc2NvbmZpZ3VyYXRpb25zGNLQ'
    '7lMgAygLMi8ubW9uaXRvcmluZy5NZXRyaWNTdHJlYW1TdGF0aXN0aWNzQ29uZmlndXJhdGlvbl'
    'IYc3RhdGlzdGljc2NvbmZpZ3VyYXRpb25zEicKBHRhZ3MYwcH2tQEgAygLMg8ubW9uaXRvcmlu'
    'Zy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use putMetricStreamOutputDescriptor instead')
const PutMetricStreamOutput$json = {
  '1': 'PutMetricStreamOutput',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `PutMetricStreamOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putMetricStreamOutputDescriptor =
    $convert.base64Decode(
        'ChVQdXRNZXRyaWNTdHJlYW1PdXRwdXQSFAoDYXJuGJ2b7b8BIAEoCVIDYXJu');

@$core.Deprecated('Use rangeDescriptor instead')
const Range$json = {
  '1': 'Range',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
  ],
};

/// Descriptor for `Range`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rangeDescriptor = $convert.base64Decode(
    'CgVSYW5nZRIbCgdlbmR0aW1lGMzvvB4gASgJUgdlbmR0aW1lEiAKCXN0YXJ0dGltZRjvtOWwAS'
    'ABKAlSCXN0YXJ0dGltZQ==');

@$core.Deprecated('Use resourceNotFoundDescriptor instead')
const ResourceNotFound$json = {
  '1': 'ResourceNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundDescriptor = $convert.base64Decode(
    'ChBSZXNvdXJjZU5vdEZvdW5kEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor = $convert.base64Decode(
    'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
    'USIgoKcmVzb3VyY2VpZBiRuvH6ASABKAlSCnJlc291cmNlaWQSJgoMcmVzb3VyY2V0eXBlGN6+'
    '2I8BIAEoCVIMcmVzb3VyY2V0eXBl');

@$core.Deprecated('Use ruleDescriptor instead')
const Rule$json = {
  '1': 'Rule',
  '2': [
    {
      '1': 'schedule',
      '3': 66697965,
      '4': 1,
      '5': 11,
      '6': '.monitoring.Schedule',
      '10': 'schedule'
    },
  ],
};

/// Descriptor for `Rule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ruleDescriptor = $convert.base64Decode(
    'CgRSdWxlEjMKCHNjaGVkdWxlGO315h8gASgLMhQubW9uaXRvcmluZy5TY2hlZHVsZVIIc2NoZW'
    'R1bGU=');

@$core.Deprecated('Use scheduleDescriptor instead')
const Schedule$json = {
  '1': 'Schedule',
  '2': [
    {'1': 'duration', '3': 348604718, '4': 1, '5': 9, '10': 'duration'},
    {'1': 'expression', '3': 193051916, '4': 1, '5': 9, '10': 'expression'},
    {'1': 'timezone', '3': 246302531, '4': 1, '5': 9, '10': 'timezone'},
  ],
};

/// Descriptor for `Schedule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduleDescriptor = $convert.base64Decode(
    'CghTY2hlZHVsZRIeCghkdXJhdGlvbhiukp2mASABKAlSCGR1cmF0aW9uEiEKCmV4cHJlc3Npb2'
    '4YjPqGXCABKAlSCmV4cHJlc3Npb24SHQoIdGltZXpvbmUYw465dSABKAlSCHRpbWV6b25l');

@$core.Deprecated('Use setAlarmStateInputDescriptor instead')
const SetAlarmStateInput$json = {
  '1': 'SetAlarmStateInput',
  '2': [
    {'1': 'alarmname', '3': 54095842, '4': 1, '5': 9, '10': 'alarmname'},
    {'1': 'statereason', '3': 376138483, '4': 1, '5': 9, '10': 'statereason'},
    {
      '1': 'statereasondata',
      '3': 262771075,
      '4': 1,
      '5': 9,
      '10': 'statereasondata'
    },
    {
      '1': 'statevalue',
      '3': 334526008,
      '4': 1,
      '5': 14,
      '6': '.monitoring.StateValue',
      '10': 'statevalue'
    },
  ],
};

/// Descriptor for `SetAlarmStateInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setAlarmStateInputDescriptor = $convert.base64Decode(
    'ChJTZXRBbGFybVN0YXRlSW5wdXQSHwoJYWxhcm1uYW1lGOLf5RkgASgJUglhbGFybW5hbWUSJA'
    'oLc3RhdGVyZWFzb24Y89WtswEgASgJUgtzdGF0ZXJlYXNvbhIrCg9zdGF0ZXJlYXNvbmRhdGEY'
    'g6OmfSABKAlSD3N0YXRlcmVhc29uZGF0YRI6CgpzdGF0ZXZhbHVlGLjswZ8BIAEoDjIWLm1vbm'
    'l0b3JpbmcuU3RhdGVWYWx1ZVIKc3RhdGV2YWx1ZQ==');

@$core.Deprecated('Use singleMetricAnomalyDetectorDescriptor instead')
const SingleMetricAnomalyDetector$json = {
  '1': 'SingleMetricAnomalyDetector',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Dimension',
      '10': 'dimensions'
    },
    {'1': 'metricname', '3': 106340219, '4': 1, '5': 9, '10': 'metricname'},
    {'1': 'namespace', '3': 355353153, '4': 1, '5': 9, '10': 'namespace'},
    {'1': 'stat', '3': 325549752, '4': 1, '5': 9, '10': 'stat'},
  ],
};

/// Descriptor for `SingleMetricAnomalyDetector`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List singleMetricAnomalyDetectorDescriptor = $convert.base64Decode(
    'ChtTaW5nbGVNZXRyaWNBbm9tYWx5RGV0ZWN0b3ISHwoJYWNjb3VudGlkGNLBuR8gASgJUglhY2'
    'NvdW50aWQSOQoKZGltZW5zaW9ucxjRm9/cASADKAsyFS5tb25pdG9yaW5nLkRpbWVuc2lvblIK'
    'ZGltZW5zaW9ucxIhCgptZXRyaWNuYW1lGPu+2jIgASgJUgptZXRyaWNuYW1lEiAKCW5hbWVzcG'
    'FjZRjBhLmpASABKAlSCW5hbWVzcGFjZRIWCgRzdGF0GLj9nZsBIAEoCVIEc3RhdA==');

@$core.Deprecated('Use startMetricStreamsInputDescriptor instead')
const StartMetricStreamsInput$json = {
  '1': 'StartMetricStreamsInput',
  '2': [
    {'1': 'names', '3': 324387120, '4': 3, '5': 9, '10': 'names'},
  ],
};

/// Descriptor for `StartMetricStreamsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startMetricStreamsInputDescriptor =
    $convert.base64Decode(
        'ChdTdGFydE1ldHJpY1N0cmVhbXNJbnB1dBIYCgVuYW1lcxiwgteaASADKAlSBW5hbWVz');

@$core.Deprecated('Use startMetricStreamsOutputDescriptor instead')
const StartMetricStreamsOutput$json = {
  '1': 'StartMetricStreamsOutput',
};

/// Descriptor for `StartMetricStreamsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startMetricStreamsOutputDescriptor =
    $convert.base64Decode('ChhTdGFydE1ldHJpY1N0cmVhbXNPdXRwdXQ=');

@$core.Deprecated('Use statisticSetDescriptor instead')
const StatisticSet$json = {
  '1': 'StatisticSet',
  '2': [
    {'1': 'maximum', '3': 43343394, '4': 1, '5': 1, '10': 'maximum'},
    {'1': 'minimum', '3': 438536856, '4': 1, '5': 1, '10': 'minimum'},
    {'1': 'samplecount', '3': 384919839, '4': 1, '5': 1, '10': 'samplecount'},
    {'1': 'sum', '3': 534303305, '4': 1, '5': 1, '10': 'sum'},
  ],
};

/// Descriptor for `StatisticSet`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statisticSetDescriptor = $convert.base64Decode(
    'CgxTdGF0aXN0aWNTZXQSGwoHbWF4aW11bRiivNUUIAEoAVIHbWF4aW11bRIcCgdtaW5pbXVtGJ'
    'iVjtEBIAEoAVIHbWluaW11bRIkCgtzYW1wbGVjb3VudBif0sW3ASABKAFSC3NhbXBsZWNvdW50'
    'EhQKA3N1bRjJpOP+ASABKAFSA3N1bQ==');

@$core.Deprecated('Use stopMetricStreamsInputDescriptor instead')
const StopMetricStreamsInput$json = {
  '1': 'StopMetricStreamsInput',
  '2': [
    {'1': 'names', '3': 324387120, '4': 3, '5': 9, '10': 'names'},
  ],
};

/// Descriptor for `StopMetricStreamsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopMetricStreamsInputDescriptor =
    $convert.base64Decode(
        'ChZTdG9wTWV0cmljU3RyZWFtc0lucHV0EhgKBW5hbWVzGLCC15oBIAMoCVIFbmFtZXM=');

@$core.Deprecated('Use stopMetricStreamsOutputDescriptor instead')
const StopMetricStreamsOutput$json = {
  '1': 'StopMetricStreamsOutput',
};

/// Descriptor for `StopMetricStreamsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopMetricStreamsOutputDescriptor =
    $convert.base64Decode('ChdTdG9wTWV0cmljU3RyZWFtc091dHB1dA==');

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

@$core.Deprecated('Use tagResourceInputDescriptor instead')
const TagResourceInput$json = {
  '1': 'TagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.monitoring.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceInputDescriptor = $convert.base64Decode(
    'ChBUYWdSZXNvdXJjZUlucHV0EiQKC3Jlc291cmNlYXJuGO3AmbABIAEoCVILcmVzb3VyY2Vhcm'
    '4SJwoEdGFncxjBwfa1ASADKAsyDy5tb25pdG9yaW5nLlRhZ1IEdGFncw==');

@$core.Deprecated('Use tagResourceOutputDescriptor instead')
const TagResourceOutput$json = {
  '1': 'TagResourceOutput',
};

/// Descriptor for `TagResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceOutputDescriptor =
    $convert.base64Decode('ChFUYWdSZXNvdXJjZU91dHB1dA==');

@$core.Deprecated('Use untagResourceInputDescriptor instead')
const UntagResourceInput$json = {
  '1': 'UntagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceInputDescriptor = $convert.base64Decode(
    'ChJVbnRhZ1Jlc291cmNlSW5wdXQSJAoLcmVzb3VyY2Vhcm4Y7cCZsAEgASgJUgtyZXNvdXJjZW'
    'FybhIcCgd0YWdrZXlzGPzD85gBIAMoCVIHdGFna2V5cw==');

@$core.Deprecated('Use untagResourceOutputDescriptor instead')
const UntagResourceOutput$json = {
  '1': 'UntagResourceOutput',
};

/// Descriptor for `UntagResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceOutputDescriptor =
    $convert.base64Decode('ChNVbnRhZ1Jlc291cmNlT3V0cHV0');
