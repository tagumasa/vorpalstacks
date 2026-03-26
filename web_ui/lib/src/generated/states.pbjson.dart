// This is a generated file - do not edit.
//
// Generated from states.proto.

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

@$core.Deprecated('Use encryptionTypeDescriptor instead')
const EncryptionType$json = {
  '1': 'EncryptionType',
  '2': [
    {'1': 'ENCRYPTION_TYPE_CUSTOMER_MANAGED_KMS_KEY', '2': 0},
    {'1': 'ENCRYPTION_TYPE_AWS_OWNED_KEY', '2': 1},
  ],
};

/// Descriptor for `EncryptionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List encryptionTypeDescriptor = $convert.base64Decode(
    'Cg5FbmNyeXB0aW9uVHlwZRIsCihFTkNSWVBUSU9OX1RZUEVfQ1VTVE9NRVJfTUFOQUdFRF9LTV'
    'NfS0VZEAASIQodRU5DUllQVElPTl9UWVBFX0FXU19PV05FRF9LRVkQAQ==');

@$core.Deprecated('Use executionRedriveFilterDescriptor instead')
const ExecutionRedriveFilter$json = {
  '1': 'ExecutionRedriveFilter',
  '2': [
    {'1': 'EXECUTION_REDRIVE_FILTER_REDRIVEN', '2': 0},
    {'1': 'EXECUTION_REDRIVE_FILTER_NOT_REDRIVEN', '2': 1},
  ],
};

/// Descriptor for `ExecutionRedriveFilter`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executionRedriveFilterDescriptor = $convert.base64Decode(
    'ChZFeGVjdXRpb25SZWRyaXZlRmlsdGVyEiUKIUVYRUNVVElPTl9SRURSSVZFX0ZJTFRFUl9SRU'
    'RSSVZFThAAEikKJUVYRUNVVElPTl9SRURSSVZFX0ZJTFRFUl9OT1RfUkVEUklWRU4QAQ==');

@$core.Deprecated('Use executionRedriveStatusDescriptor instead')
const ExecutionRedriveStatus$json = {
  '1': 'ExecutionRedriveStatus',
  '2': [
    {'1': 'EXECUTION_REDRIVE_STATUS_REDRIVABLE_BY_MAP_RUN', '2': 0},
    {'1': 'EXECUTION_REDRIVE_STATUS_NOT_REDRIVABLE', '2': 1},
    {'1': 'EXECUTION_REDRIVE_STATUS_REDRIVABLE', '2': 2},
  ],
};

/// Descriptor for `ExecutionRedriveStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executionRedriveStatusDescriptor = $convert.base64Decode(
    'ChZFeGVjdXRpb25SZWRyaXZlU3RhdHVzEjIKLkVYRUNVVElPTl9SRURSSVZFX1NUQVRVU19SRU'
    'RSSVZBQkxFX0JZX01BUF9SVU4QABIrCidFWEVDVVRJT05fUkVEUklWRV9TVEFUVVNfTk9UX1JF'
    'RFJJVkFCTEUQARInCiNFWEVDVVRJT05fUkVEUklWRV9TVEFUVVNfUkVEUklWQUJMRRAC');

@$core.Deprecated('Use executionStatusDescriptor instead')
const ExecutionStatus$json = {
  '1': 'ExecutionStatus',
  '2': [
    {'1': 'EXECUTION_STATUS_ABORTED', '2': 0},
    {'1': 'EXECUTION_STATUS_RUNNING', '2': 1},
    {'1': 'EXECUTION_STATUS_TIMED_OUT', '2': 2},
    {'1': 'EXECUTION_STATUS_SUCCEEDED', '2': 3},
    {'1': 'EXECUTION_STATUS_PENDING_REDRIVE', '2': 4},
    {'1': 'EXECUTION_STATUS_FAILED', '2': 5},
  ],
};

/// Descriptor for `ExecutionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executionStatusDescriptor = $convert.base64Decode(
    'Cg9FeGVjdXRpb25TdGF0dXMSHAoYRVhFQ1VUSU9OX1NUQVRVU19BQk9SVEVEEAASHAoYRVhFQ1'
    'VUSU9OX1NUQVRVU19SVU5OSU5HEAESHgoaRVhFQ1VUSU9OX1NUQVRVU19USU1FRF9PVVQQAhIe'
    'ChpFWEVDVVRJT05fU1RBVFVTX1NVQ0NFRURFRBADEiQKIEVYRUNVVElPTl9TVEFUVVNfUEVORE'
    'lOR19SRURSSVZFEAQSGwoXRVhFQ1VUSU9OX1NUQVRVU19GQUlMRUQQBQ==');

@$core.Deprecated('Use historyEventTypeDescriptor instead')
const HistoryEventType$json = {
  '1': 'HistoryEventType',
  '2': [
    {'1': 'HISTORY_EVENT_TYPE_TASKSUBMITTED', '2': 0},
    {'1': 'HISTORY_EVENT_TYPE_MAPRUNSUCCEEDED', '2': 1},
    {'1': 'HISTORY_EVENT_TYPE_WAITSTATEEXITED', '2': 2},
    {'1': 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSCHEDULEFAILED', '2': 3},
    {'1': 'HISTORY_EVENT_TYPE_MAPSTATESTARTED', '2': 4},
    {'1': 'HISTORY_EVENT_TYPE_TASKSTARTED', '2': 5},
    {'1': 'HISTORY_EVENT_TYPE_MAPITERATIONFAILED', '2': 6},
    {'1': 'HISTORY_EVENT_TYPE_PASSSTATEEXITED', '2': 7},
    {'1': 'HISTORY_EVENT_TYPE_SUCCEEDSTATEEXITED', '2': 8},
    {'1': 'HISTORY_EVENT_TYPE_ACTIVITYSTARTED', '2': 9},
    {'1': 'HISTORY_EVENT_TYPE_EXECUTIONREDRIVEN', '2': 10},
    {'1': 'HISTORY_EVENT_TYPE_EXECUTIONABORTED', '2': 11},
    {'1': 'HISTORY_EVENT_TYPE_EVALUATIONFAILED', '2': 12},
    {'1': 'HISTORY_EVENT_TYPE_MAPSTATEFAILED', '2': 13},
    {'1': 'HISTORY_EVENT_TYPE_TASKSCHEDULED', '2': 14},
    {'1': 'HISTORY_EVENT_TYPE_MAPRUNREDRIVEN', '2': 15},
    {'1': 'HISTORY_EVENT_TYPE_PARALLELSTATEFAILED', '2': 16},
    {'1': 'HISTORY_EVENT_TYPE_MAPITERATIONSTARTED', '2': 17},
    {'1': 'HISTORY_EVENT_TYPE_ACTIVITYSCHEDULED', '2': 18},
    {'1': 'HISTORY_EVENT_TYPE_EXECUTIONTIMEDOUT', '2': 19},
    {'1': 'HISTORY_EVENT_TYPE_ACTIVITYFAILED', '2': 20},
    {'1': 'HISTORY_EVENT_TYPE_TASKSUBMITFAILED', '2': 21},
    {'1': 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSTARTED', '2': 22},
    {'1': 'HISTORY_EVENT_TYPE_TASKSUCCEEDED', '2': 23},
    {'1': 'HISTORY_EVENT_TYPE_MAPITERATIONSUCCEEDED', '2': 24},
    {'1': 'HISTORY_EVENT_TYPE_ACTIVITYSUCCEEDED', '2': 25},
    {'1': 'HISTORY_EVENT_TYPE_FAILSTATEENTERED', '2': 26},
    {'1': 'HISTORY_EVENT_TYPE_EXECUTIONFAILED', '2': 27},
    {'1': 'HISTORY_EVENT_TYPE_MAPSTATESUCCEEDED', '2': 28},
    {'1': 'HISTORY_EVENT_TYPE_MAPRUNFAILED', '2': 29},
    {'1': 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSCHEDULED', '2': 30},
    {'1': 'HISTORY_EVENT_TYPE_PARALLELSTATEEXITED', '2': 31},
    {'1': 'HISTORY_EVENT_TYPE_EXECUTIONSTARTED', '2': 32},
    {'1': 'HISTORY_EVENT_TYPE_PARALLELSTATEENTERED', '2': 33},
    {'1': 'HISTORY_EVENT_TYPE_TASKSTATEENTERED', '2': 34},
    {'1': 'HISTORY_EVENT_TYPE_CHOICESTATEENTERED', '2': 35},
    {'1': 'HISTORY_EVENT_TYPE_PASSSTATEENTERED', '2': 36},
    {'1': 'HISTORY_EVENT_TYPE_WAITSTATEABORTED', '2': 37},
    {'1': 'HISTORY_EVENT_TYPE_PARALLELSTATEABORTED', '2': 38},
    {'1': 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONTIMEDOUT', '2': 39},
    {'1': 'HISTORY_EVENT_TYPE_MAPRUNABORTED', '2': 40},
    {'1': 'HISTORY_EVENT_TYPE_MAPSTATEABORTED', '2': 41},
    {'1': 'HISTORY_EVENT_TYPE_ACTIVITYTIMEDOUT', '2': 42},
    {'1': 'HISTORY_EVENT_TYPE_SUCCEEDSTATEENTERED', '2': 43},
    {'1': 'HISTORY_EVENT_TYPE_ACTIVITYSCHEDULEFAILED', '2': 44},
    {'1': 'HISTORY_EVENT_TYPE_MAPRUNSTARTED', '2': 45},
    {'1': 'HISTORY_EVENT_TYPE_MAPSTATEENTERED', '2': 46},
    {'1': 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONFAILED', '2': 47},
    {'1': 'HISTORY_EVENT_TYPE_TASKFAILED', '2': 48},
    {'1': 'HISTORY_EVENT_TYPE_CHOICESTATEEXITED', '2': 49},
    {'1': 'HISTORY_EVENT_TYPE_EXECUTIONSUCCEEDED', '2': 50},
    {'1': 'HISTORY_EVENT_TYPE_WAITSTATEENTERED', '2': 51},
    {'1': 'HISTORY_EVENT_TYPE_TASKTIMEDOUT', '2': 52},
    {'1': 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSUCCEEDED', '2': 53},
    {'1': 'HISTORY_EVENT_TYPE_MAPITERATIONABORTED', '2': 54},
    {'1': 'HISTORY_EVENT_TYPE_PARALLELSTATESTARTED', '2': 55},
    {'1': 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSTARTFAILED', '2': 56},
    {'1': 'HISTORY_EVENT_TYPE_TASKSTARTFAILED', '2': 57},
    {'1': 'HISTORY_EVENT_TYPE_PARALLELSTATESUCCEEDED', '2': 58},
    {'1': 'HISTORY_EVENT_TYPE_TASKSTATEEXITED', '2': 59},
    {'1': 'HISTORY_EVENT_TYPE_MAPSTATEEXITED', '2': 60},
    {'1': 'HISTORY_EVENT_TYPE_TASKSTATEABORTED', '2': 61},
  ],
};

/// Descriptor for `HistoryEventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List historyEventTypeDescriptor = $convert.base64Decode(
    'ChBIaXN0b3J5RXZlbnRUeXBlEiQKIEhJU1RPUllfRVZFTlRfVFlQRV9UQVNLU1VCTUlUVEVEEA'
    'ASJgoiSElTVE9SWV9FVkVOVF9UWVBFX01BUFJVTlNVQ0NFRURFRBABEiYKIkhJU1RPUllfRVZF'
    'TlRfVFlQRV9XQUlUU1RBVEVFWElURUQQAhIzCi9ISVNUT1JZX0VWRU5UX1RZUEVfTEFNQkRBRl'
    'VOQ1RJT05TQ0hFRFVMRUZBSUxFRBADEiYKIkhJU1RPUllfRVZFTlRfVFlQRV9NQVBTVEFURVNU'
    'QVJURUQQBBIiCh5ISVNUT1JZX0VWRU5UX1RZUEVfVEFTS1NUQVJURUQQBRIpCiVISVNUT1JZX0'
    'VWRU5UX1RZUEVfTUFQSVRFUkFUSU9ORkFJTEVEEAYSJgoiSElTVE9SWV9FVkVOVF9UWVBFX1BB'
    'U1NTVEFURUVYSVRFRBAHEikKJUhJU1RPUllfRVZFTlRfVFlQRV9TVUNDRUVEU1RBVEVFWElURU'
    'QQCBImCiJISVNUT1JZX0VWRU5UX1RZUEVfQUNUSVZJVFlTVEFSVEVEEAkSKAokSElTVE9SWV9F'
    'VkVOVF9UWVBFX0VYRUNVVElPTlJFRFJJVkVOEAoSJwojSElTVE9SWV9FVkVOVF9UWVBFX0VYRU'
    'NVVElPTkFCT1JURUQQCxInCiNISVNUT1JZX0VWRU5UX1RZUEVfRVZBTFVBVElPTkZBSUxFRBAM'
    'EiUKIUhJU1RPUllfRVZFTlRfVFlQRV9NQVBTVEFURUZBSUxFRBANEiQKIEhJU1RPUllfRVZFTl'
    'RfVFlQRV9UQVNLU0NIRURVTEVEEA4SJQohSElTVE9SWV9FVkVOVF9UWVBFX01BUFJVTlJFRFJJ'
    'VkVOEA8SKgomSElTVE9SWV9FVkVOVF9UWVBFX1BBUkFMTEVMU1RBVEVGQUlMRUQQEBIqCiZISV'
    'NUT1JZX0VWRU5UX1RZUEVfTUFQSVRFUkFUSU9OU1RBUlRFRBAREigKJEhJU1RPUllfRVZFTlRf'
    'VFlQRV9BQ1RJVklUWVNDSEVEVUxFRBASEigKJEhJU1RPUllfRVZFTlRfVFlQRV9FWEVDVVRJT0'
    '5USU1FRE9VVBATEiUKIUhJU1RPUllfRVZFTlRfVFlQRV9BQ1RJVklUWUZBSUxFRBAUEicKI0hJ'
    'U1RPUllfRVZFTlRfVFlQRV9UQVNLU1VCTUlURkFJTEVEEBUSLAooSElTVE9SWV9FVkVOVF9UWV'
    'BFX0xBTUJEQUZVTkNUSU9OU1RBUlRFRBAWEiQKIEhJU1RPUllfRVZFTlRfVFlQRV9UQVNLU1VD'
    'Q0VFREVEEBcSLAooSElTVE9SWV9FVkVOVF9UWVBFX01BUElURVJBVElPTlNVQ0NFRURFRBAYEi'
    'gKJEhJU1RPUllfRVZFTlRfVFlQRV9BQ1RJVklUWVNVQ0NFRURFRBAZEicKI0hJU1RPUllfRVZF'
    'TlRfVFlQRV9GQUlMU1RBVEVFTlRFUkVEEBoSJgoiSElTVE9SWV9FVkVOVF9UWVBFX0VYRUNVVE'
    'lPTkZBSUxFRBAbEigKJEhJU1RPUllfRVZFTlRfVFlQRV9NQVBTVEFURVNVQ0NFRURFRBAcEiMK'
    'H0hJU1RPUllfRVZFTlRfVFlQRV9NQVBSVU5GQUlMRUQQHRIuCipISVNUT1JZX0VWRU5UX1RZUE'
    'VfTEFNQkRBRlVOQ1RJT05TQ0hFRFVMRUQQHhIqCiZISVNUT1JZX0VWRU5UX1RZUEVfUEFSQUxM'
    'RUxTVEFURUVYSVRFRBAfEicKI0hJU1RPUllfRVZFTlRfVFlQRV9FWEVDVVRJT05TVEFSVEVEEC'
    'ASKwonSElTVE9SWV9FVkVOVF9UWVBFX1BBUkFMTEVMU1RBVEVFTlRFUkVEECESJwojSElTVE9S'
    'WV9FVkVOVF9UWVBFX1RBU0tTVEFURUVOVEVSRUQQIhIpCiVISVNUT1JZX0VWRU5UX1RZUEVfQ0'
    'hPSUNFU1RBVEVFTlRFUkVEECMSJwojSElTVE9SWV9FVkVOVF9UWVBFX1BBU1NTVEFURUVOVEVS'
    'RUQQJBInCiNISVNUT1JZX0VWRU5UX1RZUEVfV0FJVFNUQVRFQUJPUlRFRBAlEisKJ0hJU1RPUl'
    'lfRVZFTlRfVFlQRV9QQVJBTExFTFNUQVRFQUJPUlRFRBAmEi0KKUhJU1RPUllfRVZFTlRfVFlQ'
    'RV9MQU1CREFGVU5DVElPTlRJTUVET1VUECcSJAogSElTVE9SWV9FVkVOVF9UWVBFX01BUFJVTk'
    'FCT1JURUQQKBImCiJISVNUT1JZX0VWRU5UX1RZUEVfTUFQU1RBVEVBQk9SVEVEECkSJwojSElT'
    'VE9SWV9FVkVOVF9UWVBFX0FDVElWSVRZVElNRURPVVQQKhIqCiZISVNUT1JZX0VWRU5UX1RZUE'
    'VfU1VDQ0VFRFNUQVRFRU5URVJFRBArEi0KKUhJU1RPUllfRVZFTlRfVFlQRV9BQ1RJVklUWVND'
    'SEVEVUxFRkFJTEVEECwSJAogSElTVE9SWV9FVkVOVF9UWVBFX01BUFJVTlNUQVJURUQQLRImCi'
    'JISVNUT1JZX0VWRU5UX1RZUEVfTUFQU1RBVEVFTlRFUkVEEC4SKwonSElTVE9SWV9FVkVOVF9U'
    'WVBFX0xBTUJEQUZVTkNUSU9ORkFJTEVEEC8SIQodSElTVE9SWV9FVkVOVF9UWVBFX1RBU0tGQU'
    'lMRUQQMBIoCiRISVNUT1JZX0VWRU5UX1RZUEVfQ0hPSUNFU1RBVEVFWElURUQQMRIpCiVISVNU'
    'T1JZX0VWRU5UX1RZUEVfRVhFQ1VUSU9OU1VDQ0VFREVEEDISJwojSElTVE9SWV9FVkVOVF9UWV'
    'BFX1dBSVRTVEFURUVOVEVSRUQQMxIjCh9ISVNUT1JZX0VWRU5UX1RZUEVfVEFTS1RJTUVET1VU'
    'EDQSLgoqSElTVE9SWV9FVkVOVF9UWVBFX0xBTUJEQUZVTkNUSU9OU1VDQ0VFREVEEDUSKgomSE'
    'lTVE9SWV9FVkVOVF9UWVBFX01BUElURVJBVElPTkFCT1JURUQQNhIrCidISVNUT1JZX0VWRU5U'
    'X1RZUEVfUEFSQUxMRUxTVEFURVNUQVJURUQQNxIwCixISVNUT1JZX0VWRU5UX1RZUEVfTEFNQk'
    'RBRlVOQ1RJT05TVEFSVEZBSUxFRBA4EiYKIkhJU1RPUllfRVZFTlRfVFlQRV9UQVNLU1RBUlRG'
    'QUlMRUQQORItCilISVNUT1JZX0VWRU5UX1RZUEVfUEFSQUxMRUxTVEFURVNVQ0NFRURFRBA6Ei'
    'YKIkhJU1RPUllfRVZFTlRfVFlQRV9UQVNLU1RBVEVFWElURUQQOxIlCiFISVNUT1JZX0VWRU5U'
    'X1RZUEVfTUFQU1RBVEVFWElURUQQPBInCiNISVNUT1JZX0VWRU5UX1RZUEVfVEFTS1NUQVRFQU'
    'JPUlRFRBA9');

@$core.Deprecated('Use includedDataDescriptor instead')
const IncludedData$json = {
  '1': 'IncludedData',
  '2': [
    {'1': 'INCLUDED_DATA_METADATA_ONLY', '2': 0},
    {'1': 'INCLUDED_DATA_ALL_DATA', '2': 1},
  ],
};

/// Descriptor for `IncludedData`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List includedDataDescriptor = $convert.base64Decode(
    'CgxJbmNsdWRlZERhdGESHwobSU5DTFVERURfREFUQV9NRVRBREFUQV9PTkxZEAASGgoWSU5DTF'
    'VERURfREFUQV9BTExfREFUQRAB');

@$core.Deprecated('Use inspectionLevelDescriptor instead')
const InspectionLevel$json = {
  '1': 'InspectionLevel',
  '2': [
    {'1': 'INSPECTION_LEVEL_TRACE', '2': 0},
    {'1': 'INSPECTION_LEVEL_DEBUG', '2': 1},
    {'1': 'INSPECTION_LEVEL_INFO', '2': 2},
  ],
};

/// Descriptor for `InspectionLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inspectionLevelDescriptor = $convert.base64Decode(
    'Cg9JbnNwZWN0aW9uTGV2ZWwSGgoWSU5TUEVDVElPTl9MRVZFTF9UUkFDRRAAEhoKFklOU1BFQ1'
    'RJT05fTEVWRUxfREVCVUcQARIZChVJTlNQRUNUSU9OX0xFVkVMX0lORk8QAg==');

@$core.Deprecated('Use kmsKeyStateDescriptor instead')
const KmsKeyState$json = {
  '1': 'KmsKeyState',
  '2': [
    {'1': 'KMS_KEY_STATE_DISABLED', '2': 0},
    {'1': 'KMS_KEY_STATE_PENDING_DELETION', '2': 1},
    {'1': 'KMS_KEY_STATE_UNAVAILABLE', '2': 2},
    {'1': 'KMS_KEY_STATE_PENDING_IMPORT', '2': 3},
    {'1': 'KMS_KEY_STATE_CREATING', '2': 4},
  ],
};

/// Descriptor for `KmsKeyState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List kmsKeyStateDescriptor = $convert.base64Decode(
    'CgtLbXNLZXlTdGF0ZRIaChZLTVNfS0VZX1NUQVRFX0RJU0FCTEVEEAASIgoeS01TX0tFWV9TVE'
    'FURV9QRU5ESU5HX0RFTEVUSU9OEAESHQoZS01TX0tFWV9TVEFURV9VTkFWQUlMQUJMRRACEiAK'
    'HEtNU19LRVlfU1RBVEVfUEVORElOR19JTVBPUlQQAxIaChZLTVNfS0VZX1NUQVRFX0NSRUFUSU'
    '5HEAQ=');

@$core.Deprecated('Use logLevelDescriptor instead')
const LogLevel$json = {
  '1': 'LogLevel',
  '2': [
    {'1': 'LOG_LEVEL_FATAL', '2': 0},
    {'1': 'LOG_LEVEL_OFF', '2': 1},
    {'1': 'LOG_LEVEL_ALL', '2': 2},
    {'1': 'LOG_LEVEL_ERROR', '2': 3},
  ],
};

/// Descriptor for `LogLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List logLevelDescriptor = $convert.base64Decode(
    'CghMb2dMZXZlbBITCg9MT0dfTEVWRUxfRkFUQUwQABIRCg1MT0dfTEVWRUxfT0ZGEAESEQoNTE'
    '9HX0xFVkVMX0FMTBACEhMKD0xPR19MRVZFTF9FUlJPUhAD');

@$core.Deprecated('Use mapRunStatusDescriptor instead')
const MapRunStatus$json = {
  '1': 'MapRunStatus',
  '2': [
    {'1': 'MAP_RUN_STATUS_ABORTED', '2': 0},
    {'1': 'MAP_RUN_STATUS_RUNNING', '2': 1},
    {'1': 'MAP_RUN_STATUS_SUCCEEDED', '2': 2},
    {'1': 'MAP_RUN_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `MapRunStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mapRunStatusDescriptor = $convert.base64Decode(
    'CgxNYXBSdW5TdGF0dXMSGgoWTUFQX1JVTl9TVEFUVVNfQUJPUlRFRBAAEhoKFk1BUF9SVU5fU1'
    'RBVFVTX1JVTk5JTkcQARIcChhNQVBfUlVOX1NUQVRVU19TVUNDRUVERUQQAhIZChVNQVBfUlVO'
    'X1NUQVRVU19GQUlMRUQQAw==');

@$core.Deprecated('Use mockResponseValidationModeDescriptor instead')
const MockResponseValidationMode$json = {
  '1': 'MockResponseValidationMode',
  '2': [
    {'1': 'MOCK_RESPONSE_VALIDATION_MODE_PRESENT', '2': 0},
    {'1': 'MOCK_RESPONSE_VALIDATION_MODE_NONE', '2': 1},
    {'1': 'MOCK_RESPONSE_VALIDATION_MODE_STRICT', '2': 2},
  ],
};

/// Descriptor for `MockResponseValidationMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mockResponseValidationModeDescriptor =
    $convert.base64Decode(
        'ChpNb2NrUmVzcG9uc2VWYWxpZGF0aW9uTW9kZRIpCiVNT0NLX1JFU1BPTlNFX1ZBTElEQVRJT0'
        '5fTU9ERV9QUkVTRU5UEAASJgoiTU9DS19SRVNQT05TRV9WQUxJREFUSU9OX01PREVfTk9ORRAB'
        'EigKJE1PQ0tfUkVTUE9OU0VfVkFMSURBVElPTl9NT0RFX1NUUklDVBAC');

@$core.Deprecated('Use stateMachineStatusDescriptor instead')
const StateMachineStatus$json = {
  '1': 'StateMachineStatus',
  '2': [
    {'1': 'STATE_MACHINE_STATUS_ACTIVE', '2': 0},
    {'1': 'STATE_MACHINE_STATUS_DELETING', '2': 1},
  ],
};

/// Descriptor for `StateMachineStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stateMachineStatusDescriptor = $convert.base64Decode(
    'ChJTdGF0ZU1hY2hpbmVTdGF0dXMSHwobU1RBVEVfTUFDSElORV9TVEFUVVNfQUNUSVZFEAASIQ'
    'odU1RBVEVfTUFDSElORV9TVEFUVVNfREVMRVRJTkcQAQ==');

@$core.Deprecated('Use stateMachineTypeDescriptor instead')
const StateMachineType$json = {
  '1': 'StateMachineType',
  '2': [
    {'1': 'STATE_MACHINE_TYPE_EXPRESS', '2': 0},
    {'1': 'STATE_MACHINE_TYPE_STANDARD', '2': 1},
  ],
};

/// Descriptor for `StateMachineType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stateMachineTypeDescriptor = $convert.base64Decode(
    'ChBTdGF0ZU1hY2hpbmVUeXBlEh4KGlNUQVRFX01BQ0hJTkVfVFlQRV9FWFBSRVNTEAASHwobU1'
    'RBVEVfTUFDSElORV9UWVBFX1NUQU5EQVJEEAE=');

@$core.Deprecated('Use syncExecutionStatusDescriptor instead')
const SyncExecutionStatus$json = {
  '1': 'SyncExecutionStatus',
  '2': [
    {'1': 'SYNC_EXECUTION_STATUS_TIMED_OUT', '2': 0},
    {'1': 'SYNC_EXECUTION_STATUS_SUCCEEDED', '2': 1},
    {'1': 'SYNC_EXECUTION_STATUS_FAILED', '2': 2},
  ],
};

/// Descriptor for `SyncExecutionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List syncExecutionStatusDescriptor = $convert.base64Decode(
    'ChNTeW5jRXhlY3V0aW9uU3RhdHVzEiMKH1NZTkNfRVhFQ1VUSU9OX1NUQVRVU19USU1FRF9PVV'
    'QQABIjCh9TWU5DX0VYRUNVVElPTl9TVEFUVVNfU1VDQ0VFREVEEAESIAocU1lOQ19FWEVDVVRJ'
    'T05fU1RBVFVTX0ZBSUxFRBAC');

@$core.Deprecated('Use testExecutionStatusDescriptor instead')
const TestExecutionStatus$json = {
  '1': 'TestExecutionStatus',
  '2': [
    {'1': 'TEST_EXECUTION_STATUS_RETRIABLE', '2': 0},
    {'1': 'TEST_EXECUTION_STATUS_CAUGHT_ERROR', '2': 1},
    {'1': 'TEST_EXECUTION_STATUS_SUCCEEDED', '2': 2},
    {'1': 'TEST_EXECUTION_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `TestExecutionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List testExecutionStatusDescriptor = $convert.base64Decode(
    'ChNUZXN0RXhlY3V0aW9uU3RhdHVzEiMKH1RFU1RfRVhFQ1VUSU9OX1NUQVRVU19SRVRSSUFCTE'
    'UQABImCiJURVNUX0VYRUNVVElPTl9TVEFUVVNfQ0FVR0hUX0VSUk9SEAESIwofVEVTVF9FWEVD'
    'VVRJT05fU1RBVFVTX1NVQ0NFRURFRBACEiAKHFRFU1RfRVhFQ1VUSU9OX1NUQVRVU19GQUlMRU'
    'QQAw==');

@$core.Deprecated(
    'Use validateStateMachineDefinitionResultCodeDescriptor instead')
const ValidateStateMachineDefinitionResultCode$json = {
  '1': 'ValidateStateMachineDefinitionResultCode',
  '2': [
    {'1': 'VALIDATE_STATE_MACHINE_DEFINITION_RESULT_CODE_OK', '2': 0},
    {'1': 'VALIDATE_STATE_MACHINE_DEFINITION_RESULT_CODE_FAIL', '2': 1},
  ],
};

/// Descriptor for `ValidateStateMachineDefinitionResultCode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List validateStateMachineDefinitionResultCodeDescriptor =
    $convert.base64Decode(
        'CihWYWxpZGF0ZVN0YXRlTWFjaGluZURlZmluaXRpb25SZXN1bHRDb2RlEjQKMFZBTElEQVRFX1'
        'NUQVRFX01BQ0hJTkVfREVGSU5JVElPTl9SRVNVTFRfQ09ERV9PSxAAEjYKMlZBTElEQVRFX1NU'
        'QVRFX01BQ0hJTkVfREVGSU5JVElPTl9SRVNVTFRfQ09ERV9GQUlMEAE=');

@$core
    .Deprecated('Use validateStateMachineDefinitionSeverityDescriptor instead')
const ValidateStateMachineDefinitionSeverity$json = {
  '1': 'ValidateStateMachineDefinitionSeverity',
  '2': [
    {'1': 'VALIDATE_STATE_MACHINE_DEFINITION_SEVERITY_WARNING', '2': 0},
    {'1': 'VALIDATE_STATE_MACHINE_DEFINITION_SEVERITY_ERROR', '2': 1},
  ],
};

/// Descriptor for `ValidateStateMachineDefinitionSeverity`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List validateStateMachineDefinitionSeverityDescriptor =
    $convert.base64Decode(
        'CiZWYWxpZGF0ZVN0YXRlTWFjaGluZURlZmluaXRpb25TZXZlcml0eRI2CjJWQUxJREFURV9TVE'
        'FURV9NQUNISU5FX0RFRklOSVRJT05fU0VWRVJJVFlfV0FSTklORxAAEjQKMFZBTElEQVRFX1NU'
        'QVRFX01BQ0hJTkVfREVGSU5JVElPTl9TRVZFUklUWV9FUlJPUhAB');

@$core.Deprecated('Use validationExceptionReasonDescriptor instead')
const ValidationExceptionReason$json = {
  '1': 'ValidationExceptionReason',
  '2': [
    {
      '1': 'VALIDATION_EXCEPTION_REASON_API_DOES_NOT_SUPPORT_LABELED_ARNS',
      '2': 0
    },
    {'1': 'VALIDATION_EXCEPTION_REASON_INVALID_ROUTING_CONFIGURATION', '2': 1},
    {
      '1': 'VALIDATION_EXCEPTION_REASON_CANNOT_UPDATE_COMPLETED_MAP_RUN',
      '2': 2
    },
    {'1': 'VALIDATION_EXCEPTION_REASON_MISSING_REQUIRED_PARAMETER', '2': 3},
  ],
};

/// Descriptor for `ValidationExceptionReason`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List validationExceptionReasonDescriptor = $convert.base64Decode(
    'ChlWYWxpZGF0aW9uRXhjZXB0aW9uUmVhc29uEkEKPVZBTElEQVRJT05fRVhDRVBUSU9OX1JFQV'
    'NPTl9BUElfRE9FU19OT1RfU1VQUE9SVF9MQUJFTEVEX0FSTlMQABI9CjlWQUxJREFUSU9OX0VY'
    'Q0VQVElPTl9SRUFTT05fSU5WQUxJRF9ST1VUSU5HX0NPTkZJR1VSQVRJT04QARI/CjtWQUxJRE'
    'FUSU9OX0VYQ0VQVElPTl9SRUFTT05fQ0FOTk9UX1VQREFURV9DT01QTEVURURfTUFQX1JVThAC'
    'EjoKNlZBTElEQVRJT05fRVhDRVBUSU9OX1JFQVNPTl9NSVNTSU5HX1JFUVVJUkVEX1BBUkFNRV'
    'RFUhAD');

@$core.Deprecated('Use activityAlreadyExistsDescriptor instead')
const ActivityAlreadyExists$json = {
  '1': 'ActivityAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ActivityAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityAlreadyExistsDescriptor = $convert.base64Decode(
    'ChVBY3Rpdml0eUFscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use activityDoesNotExistDescriptor instead')
const ActivityDoesNotExist$json = {
  '1': 'ActivityDoesNotExist',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ActivityDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityDoesNotExistDescriptor =
    $convert.base64Decode(
        'ChRBY3Rpdml0eURvZXNOb3RFeGlzdBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use activityFailedEventDetailsDescriptor instead')
const ActivityFailedEventDetails$json = {
  '1': 'ActivityFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `ActivityFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChpBY3Rpdml0eUZhaWxlZEV2ZW50RGV0YWlscxIXCgVjYXVzZRihpLtFIAEoCVIFY2F1c2USFw'
        'oFZXJyb3IY0o7GDCABKAlSBWVycm9y');

@$core.Deprecated('Use activityLimitExceededDescriptor instead')
const ActivityLimitExceeded$json = {
  '1': 'ActivityLimitExceeded',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ActivityLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityLimitExceededDescriptor = $convert.base64Decode(
    'ChVBY3Rpdml0eUxpbWl0RXhjZWVkZWQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use activityListItemDescriptor instead')
const ActivityListItem$json = {
  '1': 'ActivityListItem',
  '2': [
    {'1': 'activityarn', '3': 327279492, '4': 1, '5': 9, '10': 'activityarn'},
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `ActivityListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityListItemDescriptor = $convert.base64Decode(
    'ChBBY3Rpdml0eUxpc3RJdGVtEiQKC2FjdGl2aXR5YXJuGITHh5wBIAEoCVILYWN0aXZpdHlhcm'
    '4SJQoMY3JlYXRpb25kYXRlGIHO0XEgASgJUgxjcmVhdGlvbmRhdGUSFQoEbmFtZRjn++ZpIAEo'
    'CVIEbmFtZQ==');

@$core.Deprecated('Use activityScheduleFailedEventDetailsDescriptor instead')
const ActivityScheduleFailedEventDetails$json = {
  '1': 'ActivityScheduleFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `ActivityScheduleFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityScheduleFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'CiJBY3Rpdml0eVNjaGVkdWxlRmFpbGVkRXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0UgASgJUg'
        'VjYXVzZRIXCgVlcnJvchjSjsYMIAEoCVIFZXJyb3I=');

@$core.Deprecated('Use activityScheduledEventDetailsDescriptor instead')
const ActivityScheduledEventDetails$json = {
  '1': 'ActivityScheduledEventDetails',
  '2': [
    {
      '1': 'heartbeatinseconds',
      '3': 125718754,
      '4': 1,
      '5': 3,
      '10': 'heartbeatinseconds'
    },
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'inputdetails',
      '3': 452625788,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'inputdetails'
    },
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {
      '1': 'timeoutinseconds',
      '3': 472710197,
      '4': 1,
      '5': 3,
      '10': 'timeoutinseconds'
    },
  ],
};

/// Descriptor for `ActivityScheduledEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityScheduledEventDetailsDescriptor = $convert.base64Decode(
    'Ch1BY3Rpdml0eVNjaGVkdWxlZEV2ZW50RGV0YWlscxIxChJoZWFydGJlYXRpbnNlY29uZHMY4q'
    'H5OyABKANSEmhlYXJ0YmVhdGluc2Vjb25kcxIYCgVpbnB1dBj83uHOASABKAlSBWlucHV0ElAK'
    'DGlucHV0ZGV0YWlscxj8iurXASABKAsyKC5zdGF0ZXMuSGlzdG9yeUV2ZW50RXhlY3V0aW9uRG'
    'F0YURldGFpbHNSDGlucHV0ZGV0YWlscxIdCghyZXNvdXJjZRj2//1OIAEoCVIIcmVzb3VyY2US'
    'LgoQdGltZW91dGluc2Vjb25kcxi1+LPhASABKANSEHRpbWVvdXRpbnNlY29uZHM=');

@$core.Deprecated('Use activityStartedEventDetailsDescriptor instead')
const ActivityStartedEventDetails$json = {
  '1': 'ActivityStartedEventDetails',
  '2': [
    {'1': 'workername', '3': 526761781, '4': 1, '5': 9, '10': 'workername'},
  ],
};

/// Descriptor for `ActivityStartedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityStartedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChtBY3Rpdml0eVN0YXJ0ZWRFdmVudERldGFpbHMSIgoKd29ya2VybmFtZRi1/pb7ASABKAlSCn'
        'dvcmtlcm5hbWU=');

@$core.Deprecated('Use activitySucceededEventDetailsDescriptor instead')
const ActivitySucceededEventDetails$json = {
  '1': 'ActivitySucceededEventDetails',
  '2': [
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputdetails',
      '3': 393734643,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'outputdetails'
    },
  ],
};

/// Descriptor for `ActivitySucceededEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activitySucceededEventDetailsDescriptor =
    $convert.base64Decode(
        'Ch1BY3Rpdml0eVN1Y2NlZWRlZEV2ZW50RGV0YWlscxIaCgZvdXRwdXQYhZ6lzQEgASgJUgZvdX'
        'RwdXQSUgoNb3V0cHV0ZGV0YWlscxjz09+7ASABKAsyKC5zdGF0ZXMuSGlzdG9yeUV2ZW50RXhl'
        'Y3V0aW9uRGF0YURldGFpbHNSDW91dHB1dGRldGFpbHM=');

@$core.Deprecated('Use activityTimedOutEventDetailsDescriptor instead')
const ActivityTimedOutEventDetails$json = {
  '1': 'ActivityTimedOutEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `ActivityTimedOutEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityTimedOutEventDetailsDescriptor =
    $convert.base64Decode(
        'ChxBY3Rpdml0eVRpbWVkT3V0RXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0UgASgJUgVjYXVzZR'
        'IXCgVlcnJvchjSjsYMIAEoCVIFZXJyb3I=');

@$core.Deprecated('Use activityWorkerLimitExceededDescriptor instead')
const ActivityWorkerLimitExceeded$json = {
  '1': 'ActivityWorkerLimitExceeded',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ActivityWorkerLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activityWorkerLimitExceededDescriptor =
    $convert.base64Decode(
        'ChtBY3Rpdml0eVdvcmtlckxpbWl0RXhjZWVkZWQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use assignedVariablesDetailsDescriptor instead')
const AssignedVariablesDetails$json = {
  '1': 'AssignedVariablesDetails',
  '2': [
    {'1': 'truncated', '3': 407490858, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `AssignedVariablesDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assignedVariablesDetailsDescriptor =
    $convert.base64Decode(
        'ChhBc3NpZ25lZFZhcmlhYmxlc0RldGFpbHMSIAoJdHJ1bmNhdGVkGKqip8IBIAEoCFIJdHJ1bm'
        'NhdGVk');

@$core.Deprecated('Use billingDetailsDescriptor instead')
const BillingDetails$json = {
  '1': 'BillingDetails',
  '2': [
    {
      '1': 'billeddurationinmilliseconds',
      '3': 65131593,
      '4': 1,
      '5': 3,
      '10': 'billeddurationinmilliseconds'
    },
    {
      '1': 'billedmemoryusedinmb',
      '3': 459756940,
      '4': 1,
      '5': 3,
      '10': 'billedmemoryusedinmb'
    },
  ],
};

/// Descriptor for `BillingDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List billingDetailsDescriptor = $convert.base64Decode(
    'Cg5CaWxsaW5nRGV0YWlscxJFChxiaWxsZWRkdXJhdGlvbmlubWlsbGlzZWNvbmRzGMmohx8gAS'
    'gDUhxiaWxsZWRkdXJhdGlvbmlubWlsbGlzZWNvbmRzEjYKFGJpbGxlZG1lbW9yeXVzZWRpbm1i'
    'GIyrndsBIAEoA1IUYmlsbGVkbWVtb3J5dXNlZGlubWI=');

@$core.Deprecated('Use cloudWatchEventsExecutionDataDetailsDescriptor instead')
const CloudWatchEventsExecutionDataDetails$json = {
  '1': 'CloudWatchEventsExecutionDataDetails',
  '2': [
    {'1': 'included', '3': 517673658, '4': 1, '5': 8, '10': 'included'},
  ],
};

/// Descriptor for `CloudWatchEventsExecutionDataDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudWatchEventsExecutionDataDetailsDescriptor =
    $convert.base64Decode(
        'CiRDbG91ZFdhdGNoRXZlbnRzRXhlY3V0aW9uRGF0YURldGFpbHMSHgoIaW5jbHVkZWQYuqXs9g'
        'EgASgIUghpbmNsdWRlZA==');

@$core.Deprecated('Use cloudWatchLogsLogGroupDescriptor instead')
const CloudWatchLogsLogGroup$json = {
  '1': 'CloudWatchLogsLogGroup',
  '2': [
    {'1': 'loggrouparn', '3': 6742512, '4': 1, '5': 9, '10': 'loggrouparn'},
  ],
};

/// Descriptor for `CloudWatchLogsLogGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudWatchLogsLogGroupDescriptor =
    $convert.base64Decode(
        'ChZDbG91ZFdhdGNoTG9nc0xvZ0dyb3VwEiMKC2xvZ2dyb3VwYXJuGPDDmwMgASgJUgtsb2dncm'
        '91cGFybg==');

@$core.Deprecated('Use conflictExceptionDescriptor instead')
const ConflictException$json = {
  '1': 'ConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictExceptionDescriptor = $convert.base64Decode(
    'ChFDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use createActivityInputDescriptor instead')
const CreateActivityInput$json = {
  '1': 'CreateActivityInput',
  '2': [
    {
      '1': 'encryptionconfiguration',
      '3': 167857431,
      '4': 1,
      '5': 11,
      '6': '.states.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.states.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateActivityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createActivityInputDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVBY3Rpdml0eUlucHV0ElwKF2VuY3J5cHRpb25jb25maWd1cmF0aW9uGJeahVAgAS'
    'gLMh8uc3RhdGVzLkVuY3J5cHRpb25Db25maWd1cmF0aW9uUhdlbmNyeXB0aW9uY29uZmlndXJh'
    'dGlvbhIVCgRuYW1lGOf75mkgASgJUgRuYW1lEiMKBHRhZ3MYodfboAEgAygLMgsuc3RhdGVzLl'
    'RhZ1IEdGFncw==');

@$core.Deprecated('Use createActivityOutputDescriptor instead')
const CreateActivityOutput$json = {
  '1': 'CreateActivityOutput',
  '2': [
    {'1': 'activityarn', '3': 327279492, '4': 1, '5': 9, '10': 'activityarn'},
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
  ],
};

/// Descriptor for `CreateActivityOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createActivityOutputDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVBY3Rpdml0eU91dHB1dBIkCgthY3Rpdml0eWFybhiEx4ecASABKAlSC2FjdGl2aX'
    'R5YXJuEiUKDGNyZWF0aW9uZGF0ZRiBztFxIAEoCVIMY3JlYXRpb25kYXRl');

@$core.Deprecated('Use createStateMachineAliasInputDescriptor instead')
const CreateStateMachineAliasInput$json = {
  '1': 'CreateStateMachineAliasInput',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'routingconfiguration',
      '3': 372891510,
      '4': 3,
      '5': 11,
      '6': '.states.RoutingConfigurationListItem',
      '10': 'routingconfiguration'
    },
  ],
};

/// Descriptor for `CreateStateMachineAliasInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createStateMachineAliasInputDescriptor = $convert.base64Decode(
    'ChxDcmVhdGVTdGF0ZU1hY2hpbmVBbGlhc0lucHV0EiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCV'
    'ILZGVzY3JpcHRpb24SFQoEbmFtZRjn++ZpIAEoCVIEbmFtZRJcChRyb3V0aW5nY29uZmlndXJh'
    'dGlvbhj2vuexASADKAsyJC5zdGF0ZXMuUm91dGluZ0NvbmZpZ3VyYXRpb25MaXN0SXRlbVIUcm'
    '91dGluZ2NvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use createStateMachineAliasOutputDescriptor instead')
const CreateStateMachineAliasOutput$json = {
  '1': 'CreateStateMachineAliasOutput',
  '2': [
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
  ],
};

/// Descriptor for `CreateStateMachineAliasOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createStateMachineAliasOutputDescriptor =
    $convert.base64Decode(
        'Ch1DcmVhdGVTdGF0ZU1hY2hpbmVBbGlhc091dHB1dBIlCgxjcmVhdGlvbmRhdGUYgc7RcSABKA'
        'lSDGNyZWF0aW9uZGF0ZRI2ChRzdGF0ZW1hY2hpbmVhbGlhc2FybhiR1PH8ASABKAlSFHN0YXRl'
        'bWFjaGluZWFsaWFzYXJu');

@$core.Deprecated('Use createStateMachineInputDescriptor instead')
const CreateStateMachineInput$json = {
  '1': 'CreateStateMachineInput',
  '2': [
    {'1': 'definition', '3': 68443297, '4': 1, '5': 9, '10': 'definition'},
    {
      '1': 'encryptionconfiguration',
      '3': 167857431,
      '4': 1,
      '5': 11,
      '6': '.states.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {
      '1': 'loggingconfiguration',
      '3': 420811605,
      '4': 1,
      '5': 11,
      '6': '.states.LoggingConfiguration',
      '10': 'loggingconfiguration'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'publish', '3': 264247305, '4': 1, '5': 8, '10': 'publish'},
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.states.Tag',
      '10': 'tags'
    },
    {
      '1': 'tracingconfiguration',
      '3': 491315910,
      '4': 1,
      '5': 11,
      '6': '.states.TracingConfiguration',
      '10': 'tracingconfiguration'
    },
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.states.StateMachineType',
      '10': 'type'
    },
    {
      '1': 'versiondescription',
      '3': 434714300,
      '4': 1,
      '5': 9,
      '10': 'versiondescription'
    },
  ],
};

/// Descriptor for `CreateStateMachineInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createStateMachineInputDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVTdGF0ZU1hY2hpbmVJbnB1dBIhCgpkZWZpbml0aW9uGKG50SAgASgJUgpkZWZpbm'
    'l0aW9uElwKF2VuY3J5cHRpb25jb25maWd1cmF0aW9uGJeahVAgASgLMh8uc3RhdGVzLkVuY3J5'
    'cHRpb25Db25maWd1cmF0aW9uUhdlbmNyeXB0aW9uY29uZmlndXJhdGlvbhJUChRsb2dnaW5nY2'
    '9uZmlndXJhdGlvbhjVptTIASABKAsyHC5zdGF0ZXMuTG9nZ2luZ0NvbmZpZ3VyYXRpb25SFGxv'
    'Z2dpbmdjb25maWd1cmF0aW9uEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSGwoHcHVibGlzaBiJsI'
    'B+IAEoCFIHcHVibGlzaBIbCgdyb2xlYXJuGKGXiVEgASgJUgdyb2xlYXJuEiMKBHRhZ3MYodfb'
    'oAEgAygLMgsuc3RhdGVzLlRhZ1IEdGFncxJUChR0cmFjaW5nY29uZmlndXJhdGlvbhjGxaPqAS'
    'ABKAsyHC5zdGF0ZXMuVHJhY2luZ0NvbmZpZ3VyYXRpb25SFHRyYWNpbmdjb25maWd1cmF0aW9u'
    'EjAKBHR5cGUYzuKfiQEgASgOMhguc3RhdGVzLlN0YXRlTWFjaGluZVR5cGVSBHR5cGUSMgoSdm'
    'Vyc2lvbmRlc2NyaXB0aW9uGLztpM8BIAEoCVISdmVyc2lvbmRlc2NyaXB0aW9u');

@$core.Deprecated('Use createStateMachineOutputDescriptor instead')
const CreateStateMachineOutput$json = {
  '1': 'CreateStateMachineOutput',
  '2': [
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
  ],
};

/// Descriptor for `CreateStateMachineOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createStateMachineOutputDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVTdGF0ZU1hY2hpbmVPdXRwdXQSJQoMY3JlYXRpb25kYXRlGIHO0XEgASgJUgxjcm'
    'VhdGlvbmRhdGUSLAoPc3RhdGVtYWNoaW5lYXJuGPO7xrsBIAEoCVIPc3RhdGVtYWNoaW5lYXJu'
    'EjkKFnN0YXRlbWFjaGluZXZlcnNpb25hcm4Y+YWvISABKAlSFnN0YXRlbWFjaGluZXZlcnNpb2'
    '5hcm4=');

@$core.Deprecated('Use deleteActivityInputDescriptor instead')
const DeleteActivityInput$json = {
  '1': 'DeleteActivityInput',
  '2': [
    {'1': 'activityarn', '3': 327279492, '4': 1, '5': 9, '10': 'activityarn'},
  ],
};

/// Descriptor for `DeleteActivityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteActivityInputDescriptor = $convert.base64Decode(
    'ChNEZWxldGVBY3Rpdml0eUlucHV0EiQKC2FjdGl2aXR5YXJuGITHh5wBIAEoCVILYWN0aXZpdH'
    'lhcm4=');

@$core.Deprecated('Use deleteActivityOutputDescriptor instead')
const DeleteActivityOutput$json = {
  '1': 'DeleteActivityOutput',
};

/// Descriptor for `DeleteActivityOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteActivityOutputDescriptor =
    $convert.base64Decode('ChREZWxldGVBY3Rpdml0eU91dHB1dA==');

@$core.Deprecated('Use deleteStateMachineAliasInputDescriptor instead')
const DeleteStateMachineAliasInput$json = {
  '1': 'DeleteStateMachineAliasInput',
  '2': [
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
  ],
};

/// Descriptor for `DeleteStateMachineAliasInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStateMachineAliasInputDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVTdGF0ZU1hY2hpbmVBbGlhc0lucHV0EjYKFHN0YXRlbWFjaGluZWFsaWFzYXJuGJ'
        'HU8fwBIAEoCVIUc3RhdGVtYWNoaW5lYWxpYXNhcm4=');

@$core.Deprecated('Use deleteStateMachineAliasOutputDescriptor instead')
const DeleteStateMachineAliasOutput$json = {
  '1': 'DeleteStateMachineAliasOutput',
};

/// Descriptor for `DeleteStateMachineAliasOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStateMachineAliasOutputDescriptor =
    $convert.base64Decode('Ch1EZWxldGVTdGF0ZU1hY2hpbmVBbGlhc091dHB1dA==');

@$core.Deprecated('Use deleteStateMachineInputDescriptor instead')
const DeleteStateMachineInput$json = {
  '1': 'DeleteStateMachineInput',
  '2': [
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
  ],
};

/// Descriptor for `DeleteStateMachineInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStateMachineInputDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVTdGF0ZU1hY2hpbmVJbnB1dBIsCg9zdGF0ZW1hY2hpbmVhcm4Y87vGuwEgASgJUg'
        '9zdGF0ZW1hY2hpbmVhcm4=');

@$core.Deprecated('Use deleteStateMachineOutputDescriptor instead')
const DeleteStateMachineOutput$json = {
  '1': 'DeleteStateMachineOutput',
};

/// Descriptor for `DeleteStateMachineOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStateMachineOutputDescriptor =
    $convert.base64Decode('ChhEZWxldGVTdGF0ZU1hY2hpbmVPdXRwdXQ=');

@$core.Deprecated('Use deleteStateMachineVersionInputDescriptor instead')
const DeleteStateMachineVersionInput$json = {
  '1': 'DeleteStateMachineVersionInput',
  '2': [
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
  ],
};

/// Descriptor for `DeleteStateMachineVersionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStateMachineVersionInputDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVTdGF0ZU1hY2hpbmVWZXJzaW9uSW5wdXQSOQoWc3RhdGVtYWNoaW5ldmVyc2lvbm'
        'Fybhj5ha8hIAEoCVIWc3RhdGVtYWNoaW5ldmVyc2lvbmFybg==');

@$core.Deprecated('Use deleteStateMachineVersionOutputDescriptor instead')
const DeleteStateMachineVersionOutput$json = {
  '1': 'DeleteStateMachineVersionOutput',
};

/// Descriptor for `DeleteStateMachineVersionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStateMachineVersionOutputDescriptor =
    $convert.base64Decode('Ch9EZWxldGVTdGF0ZU1hY2hpbmVWZXJzaW9uT3V0cHV0');

@$core.Deprecated('Use describeActivityInputDescriptor instead')
const DescribeActivityInput$json = {
  '1': 'DescribeActivityInput',
  '2': [
    {'1': 'activityarn', '3': 327279492, '4': 1, '5': 9, '10': 'activityarn'},
  ],
};

/// Descriptor for `DescribeActivityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeActivityInputDescriptor = $convert.base64Decode(
    'ChVEZXNjcmliZUFjdGl2aXR5SW5wdXQSJAoLYWN0aXZpdHlhcm4YhMeHnAEgASgJUgthY3Rpdm'
    'l0eWFybg==');

@$core.Deprecated('Use describeActivityOutputDescriptor instead')
const DescribeActivityOutput$json = {
  '1': 'DescribeActivityOutput',
  '2': [
    {'1': 'activityarn', '3': 327279492, '4': 1, '5': 9, '10': 'activityarn'},
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'encryptionconfiguration',
      '3': 167857431,
      '4': 1,
      '5': 11,
      '6': '.states.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeActivityOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeActivityOutputDescriptor = $convert.base64Decode(
    'ChZEZXNjcmliZUFjdGl2aXR5T3V0cHV0EiQKC2FjdGl2aXR5YXJuGITHh5wBIAEoCVILYWN0aX'
    'ZpdHlhcm4SJQoMY3JlYXRpb25kYXRlGIHO0XEgASgJUgxjcmVhdGlvbmRhdGUSXAoXZW5jcnlw'
    'dGlvbmNvbmZpZ3VyYXRpb24Yl5qFUCABKAsyHy5zdGF0ZXMuRW5jcnlwdGlvbkNvbmZpZ3VyYX'
    'Rpb25SF2VuY3J5cHRpb25jb25maWd1cmF0aW9uEhUKBG5hbWUY5/vmaSABKAlSBG5hbWU=');

@$core.Deprecated('Use describeExecutionInputDescriptor instead')
const DescribeExecutionInput$json = {
  '1': 'DescribeExecutionInput',
  '2': [
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {
      '1': 'includeddata',
      '3': 109719114,
      '4': 1,
      '5': 14,
      '6': '.states.IncludedData',
      '10': 'includeddata'
    },
  ],
};

/// Descriptor for `DescribeExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeExecutionInputDescriptor = $convert.base64Decode(
    'ChZEZXNjcmliZUV4ZWN1dGlvbklucHV0EiYKDGV4ZWN1dGlvbmFybhjtlv2VASABKAlSDGV4ZW'
    'N1dGlvbmFybhI7CgxpbmNsdWRlZGRhdGEYytyoNCABKA4yFC5zdGF0ZXMuSW5jbHVkZWREYXRh'
    'UgxpbmNsdWRlZGRhdGE=');

@$core.Deprecated('Use describeExecutionOutputDescriptor instead')
const DescribeExecutionOutput$json = {
  '1': 'DescribeExecutionOutput',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'inputdetails',
      '3': 452625788,
      '4': 1,
      '5': 11,
      '6': '.states.CloudWatchEventsExecutionDataDetails',
      '10': 'inputdetails'
    },
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputdetails',
      '3': 393734643,
      '4': 1,
      '5': 11,
      '6': '.states.CloudWatchEventsExecutionDataDetails',
      '10': 'outputdetails'
    },
    {'1': 'redrivecount', '3': 473458696, '4': 1, '5': 5, '10': 'redrivecount'},
    {'1': 'redrivedate', '3': 152812125, '4': 1, '5': 9, '10': 'redrivedate'},
    {
      '1': 'redrivestatus',
      '3': 247102059,
      '4': 1,
      '5': 14,
      '6': '.states.ExecutionRedriveStatus',
      '10': 'redrivestatus'
    },
    {
      '1': 'redrivestatusreason',
      '3': 339085215,
      '4': 1,
      '5': 9,
      '10': 'redrivestatusreason'
    },
    {'1': 'startdate', '3': 364840732, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.states.ExecutionStatus',
      '10': 'status'
    },
    {'1': 'stopdate', '3': 180697434, '4': 1, '5': 9, '10': 'stopdate'},
    {'1': 'traceheader', '3': 219960864, '4': 1, '5': 9, '10': 'traceheader'},
  ],
};

/// Descriptor for `DescribeExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeExecutionOutputDescriptor = $convert.base64Decode(
    'ChdEZXNjcmliZUV4ZWN1dGlvbk91dHB1dBIXCgVjYXVzZRihpLtFIAEoCVIFY2F1c2USFwoFZX'
    'Jyb3IY0o7GDCABKAlSBWVycm9yEiYKDGV4ZWN1dGlvbmFybhjtlv2VASABKAlSDGV4ZWN1dGlv'
    'bmFybhIYCgVpbnB1dBj83uHOASABKAlSBWlucHV0ElQKDGlucHV0ZGV0YWlscxj8iurXASABKA'
    'syLC5zdGF0ZXMuQ2xvdWRXYXRjaEV2ZW50c0V4ZWN1dGlvbkRhdGFEZXRhaWxzUgxpbnB1dGRl'
    'dGFpbHMSHwoJbWFwcnVuYXJuGLrr1gggASgJUgltYXBydW5hcm4SFQoEbmFtZRjn++ZpIAEoCV'
    'IEbmFtZRIaCgZvdXRwdXQYhZ6lzQEgASgJUgZvdXRwdXQSVgoNb3V0cHV0ZGV0YWlscxjz09+7'
    'ASABKAsyLC5zdGF0ZXMuQ2xvdWRXYXRjaEV2ZW50c0V4ZWN1dGlvbkRhdGFEZXRhaWxzUg1vdX'
    'RwdXRkZXRhaWxzEiYKDHJlZHJpdmVjb3VudBiI0OHhASABKAVSDHJlZHJpdmVjb3VudBIjCgty'
    'ZWRyaXZlZGF0ZRjd9O5IIAEoCVILcmVkcml2ZWRhdGUSRwoNcmVkcml2ZXN0YXR1cxjr9Ol1IA'
    'EoDjIeLnN0YXRlcy5FeGVjdXRpb25SZWRyaXZlU3RhdHVzUg1yZWRyaXZlc3RhdHVzEjQKE3Jl'
    'ZHJpdmVzdGF0dXNyZWFzb24Yn4/YoQEgASgJUhNyZWRyaXZlc3RhdHVzcmVhc29uEiAKCXN0YX'
    'J0ZGF0ZRicjvytASABKAlSCXN0YXJ0ZGF0ZRI2ChRzdGF0ZW1hY2hpbmVhbGlhc2FybhiR1PH8'
    'ASABKAlSFHN0YXRlbWFjaGluZWFsaWFzYXJuEiwKD3N0YXRlbWFjaGluZWFybhjzu8a7ASABKA'
    'lSD3N0YXRlbWFjaGluZWFybhI5ChZzdGF0ZW1hY2hpbmV2ZXJzaW9uYXJuGPmFryEgASgJUhZz'
    'dGF0ZW1hY2hpbmV2ZXJzaW9uYXJuEjMKBnN0YXR1cxjw763SASABKA4yFy5zdGF0ZXMuRXhlY3'
    'V0aW9uU3RhdHVzUgZzdGF0dXMSHQoIc3RvcGRhdGUY2vKUViABKAlSCHN0b3BkYXRlEiMKC3Ry'
    'YWNlaGVhZGVyGKCs8WggASgJUgt0cmFjZWhlYWRlcg==');

@$core.Deprecated('Use describeMapRunInputDescriptor instead')
const DescribeMapRunInput$json = {
  '1': 'DescribeMapRunInput',
  '2': [
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
  ],
};

/// Descriptor for `DescribeMapRunInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMapRunInputDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZU1hcFJ1bklucHV0Eh8KCW1hcHJ1bmFybhi669YIIAEoCVIJbWFwcnVuYXJu');

@$core.Deprecated('Use describeMapRunOutputDescriptor instead')
const DescribeMapRunOutput$json = {
  '1': 'DescribeMapRunOutput',
  '2': [
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {
      '1': 'executioncounts',
      '3': 317260758,
      '4': 1,
      '5': 11,
      '6': '.states.MapRunExecutionCounts',
      '10': 'executioncounts'
    },
    {
      '1': 'itemcounts',
      '3': 334548339,
      '4': 1,
      '5': 11,
      '6': '.states.MapRunItemCounts',
      '10': 'itemcounts'
    },
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
    {
      '1': 'maxconcurrency',
      '3': 100901405,
      '4': 1,
      '5': 5,
      '10': 'maxconcurrency'
    },
    {'1': 'redrivecount', '3': 473458696, '4': 1, '5': 5, '10': 'redrivecount'},
    {'1': 'redrivedate', '3': 152812125, '4': 1, '5': 9, '10': 'redrivedate'},
    {'1': 'startdate', '3': 364840732, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.states.MapRunStatus',
      '10': 'status'
    },
    {'1': 'stopdate', '3': 180697434, '4': 1, '5': 9, '10': 'stopdate'},
    {
      '1': 'toleratedfailurecount',
      '3': 41834811,
      '4': 1,
      '5': 3,
      '10': 'toleratedfailurecount'
    },
    {
      '1': 'toleratedfailurepercentage',
      '3': 116496164,
      '4': 1,
      '5': 2,
      '10': 'toleratedfailurepercentage'
    },
  ],
};

/// Descriptor for `DescribeMapRunOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMapRunOutputDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZU1hcFJ1bk91dHB1dBImCgxleGVjdXRpb25hcm4Y7Zb9lQEgASgJUgxleGVjdX'
    'Rpb25hcm4SSwoPZXhlY3V0aW9uY291bnRzGNaHpJcBIAEoCzIdLnN0YXRlcy5NYXBSdW5FeGVj'
    'dXRpb25Db3VudHNSD2V4ZWN1dGlvbmNvdW50cxI8CgppdGVtY291bnRzGPOaw58BIAEoCzIYLn'
    'N0YXRlcy5NYXBSdW5JdGVtQ291bnRzUgppdGVtY291bnRzEh8KCW1hcHJ1bmFybhi669YIIAEo'
    'CVIJbWFwcnVuYXJuEikKDm1heGNvbmN1cnJlbmN5GJ3EjjAgASgFUg5tYXhjb25jdXJyZW5jeR'
    'ImCgxyZWRyaXZlY291bnQYiNDh4QEgASgFUgxyZWRyaXZlY291bnQSIwoLcmVkcml2ZWRhdGUY'
    '3fTuSCABKAlSC3JlZHJpdmVkYXRlEiAKCXN0YXJ0ZGF0ZRicjvytASABKAlSCXN0YXJ0ZGF0ZR'
    'IwCgZzdGF0dXMY8O+t0gEgASgOMhQuc3RhdGVzLk1hcFJ1blN0YXR1c1IGc3RhdHVzEh0KCHN0'
    'b3BkYXRlGNrylFYgASgJUghzdG9wZGF0ZRI3ChV0b2xlcmF0ZWRmYWlsdXJlY291bnQYu7L5Ey'
    'ABKANSFXRvbGVyYXRlZGZhaWx1cmVjb3VudBJBChp0b2xlcmF0ZWRmYWlsdXJlcGVyY2VudGFn'
    'ZRikrsY3IAEoAlIadG9sZXJhdGVkZmFpbHVyZXBlcmNlbnRhZ2U=');

@$core.Deprecated('Use describeStateMachineAliasInputDescriptor instead')
const DescribeStateMachineAliasInput$json = {
  '1': 'DescribeStateMachineAliasInput',
  '2': [
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
  ],
};

/// Descriptor for `DescribeStateMachineAliasInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStateMachineAliasInputDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZVN0YXRlTWFjaGluZUFsaWFzSW5wdXQSNgoUc3RhdGVtYWNoaW5lYWxpYXNhcm'
        '4YkdTx/AEgASgJUhRzdGF0ZW1hY2hpbmVhbGlhc2Fybg==');

@$core.Deprecated('Use describeStateMachineAliasOutputDescriptor instead')
const DescribeStateMachineAliasOutput$json = {
  '1': 'DescribeStateMachineAliasOutput',
  '2': [
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'routingconfiguration',
      '3': 372891510,
      '4': 3,
      '5': 11,
      '6': '.states.RoutingConfigurationListItem',
      '10': 'routingconfiguration'
    },
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
    {'1': 'updatedate', '3': 510552561, '4': 1, '5': 9, '10': 'updatedate'},
  ],
};

/// Descriptor for `DescribeStateMachineAliasOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStateMachineAliasOutputDescriptor = $convert.base64Decode(
    'Ch9EZXNjcmliZVN0YXRlTWFjaGluZUFsaWFzT3V0cHV0EiUKDGNyZWF0aW9uZGF0ZRiBztFxIA'
    'EoCVIMY3JlYXRpb25kYXRlEiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24S'
    'FQoEbmFtZRjn++ZpIAEoCVIEbmFtZRJcChRyb3V0aW5nY29uZmlndXJhdGlvbhj2vuexASADKA'
    'syJC5zdGF0ZXMuUm91dGluZ0NvbmZpZ3VyYXRpb25MaXN0SXRlbVIUcm91dGluZ2NvbmZpZ3Vy'
    'YXRpb24SNgoUc3RhdGVtYWNoaW5lYWxpYXNhcm4YkdTx/AEgASgJUhRzdGF0ZW1hY2hpbmVhbG'
    'lhc2FybhIiCgp1cGRhdGVkYXRlGPHTufMBIAEoCVIKdXBkYXRlZGF0ZQ==');

@$core.Deprecated('Use describeStateMachineForExecutionInputDescriptor instead')
const DescribeStateMachineForExecutionInput$json = {
  '1': 'DescribeStateMachineForExecutionInput',
  '2': [
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {
      '1': 'includeddata',
      '3': 109719114,
      '4': 1,
      '5': 14,
      '6': '.states.IncludedData',
      '10': 'includeddata'
    },
  ],
};

/// Descriptor for `DescribeStateMachineForExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStateMachineForExecutionInputDescriptor =
    $convert.base64Decode(
        'CiVEZXNjcmliZVN0YXRlTWFjaGluZUZvckV4ZWN1dGlvbklucHV0EiYKDGV4ZWN1dGlvbmFybh'
        'jtlv2VASABKAlSDGV4ZWN1dGlvbmFybhI7CgxpbmNsdWRlZGRhdGEYytyoNCABKA4yFC5zdGF0'
        'ZXMuSW5jbHVkZWREYXRhUgxpbmNsdWRlZGRhdGE=');

@$core
    .Deprecated('Use describeStateMachineForExecutionOutputDescriptor instead')
const DescribeStateMachineForExecutionOutput$json = {
  '1': 'DescribeStateMachineForExecutionOutput',
  '2': [
    {'1': 'definition', '3': 68443297, '4': 1, '5': 9, '10': 'definition'},
    {
      '1': 'encryptionconfiguration',
      '3': 167857431,
      '4': 1,
      '5': 11,
      '6': '.states.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {'1': 'label', '3': 379000830, '4': 1, '5': 9, '10': 'label'},
    {
      '1': 'loggingconfiguration',
      '3': 420811605,
      '4': 1,
      '5': 11,
      '6': '.states.LoggingConfiguration',
      '10': 'loggingconfiguration'
    },
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'revisionid', '3': 369170086, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'tracingconfiguration',
      '3': 491315910,
      '4': 1,
      '5': 11,
      '6': '.states.TracingConfiguration',
      '10': 'tracingconfiguration'
    },
    {'1': 'updatedate', '3': 510552561, '4': 1, '5': 9, '10': 'updatedate'},
    {
      '1': 'variablereferences',
      '3': 150942252,
      '4': 3,
      '5': 11,
      '6':
          '.states.DescribeStateMachineForExecutionOutput.VariablereferencesEntry',
      '10': 'variablereferences'
    },
  ],
  '3': [DescribeStateMachineForExecutionOutput_VariablereferencesEntry$json],
};

@$core
    .Deprecated('Use describeStateMachineForExecutionOutputDescriptor instead')
const DescribeStateMachineForExecutionOutput_VariablereferencesEntry$json = {
  '1': 'VariablereferencesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DescribeStateMachineForExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStateMachineForExecutionOutputDescriptor = $convert.base64Decode(
    'CiZEZXNjcmliZVN0YXRlTWFjaGluZUZvckV4ZWN1dGlvbk91dHB1dBIhCgpkZWZpbml0aW9uGK'
    'G50SAgASgJUgpkZWZpbml0aW9uElwKF2VuY3J5cHRpb25jb25maWd1cmF0aW9uGJeahVAgASgL'
    'Mh8uc3RhdGVzLkVuY3J5cHRpb25Db25maWd1cmF0aW9uUhdlbmNyeXB0aW9uY29uZmlndXJhdG'
    'lvbhIYCgVsYWJlbBj+r9y0ASABKAlSBWxhYmVsElQKFGxvZ2dpbmdjb25maWd1cmF0aW9uGNWm'
    '1MgBIAEoCzIcLnN0YXRlcy5Mb2dnaW5nQ29uZmlndXJhdGlvblIUbG9nZ2luZ2NvbmZpZ3VyYX'
    'Rpb24SHwoJbWFwcnVuYXJuGLrr1gggASgJUgltYXBydW5hcm4SFQoEbmFtZRjn++ZpIAEoCVIE'
    'bmFtZRIiCgpyZXZpc2lvbmlkGKathLABIAEoCVIKcmV2aXNpb25pZBIbCgdyb2xlYXJuGKGXiV'
    'EgASgJUgdyb2xlYXJuEiwKD3N0YXRlbWFjaGluZWFybhjzu8a7ASABKAlSD3N0YXRlbWFjaGlu'
    'ZWFybhJUChR0cmFjaW5nY29uZmlndXJhdGlvbhjGxaPqASABKAsyHC5zdGF0ZXMuVHJhY2luZ0'
    'NvbmZpZ3VyYXRpb25SFHRyYWNpbmdjb25maWd1cmF0aW9uEiIKCnVwZGF0ZWRhdGUY8dO58wEg'
    'ASgJUgp1cGRhdGVkYXRlEnkKEnZhcmlhYmxlcmVmZXJlbmNlcxis5PxHIAMoCzJGLnN0YXRlcy'
    '5EZXNjcmliZVN0YXRlTWFjaGluZUZvckV4ZWN1dGlvbk91dHB1dC5WYXJpYWJsZXJlZmVyZW5j'
    'ZXNFbnRyeVISdmFyaWFibGVyZWZlcmVuY2VzGkUKF1ZhcmlhYmxlcmVmZXJlbmNlc0VudHJ5Eh'
    'AKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use describeStateMachineInputDescriptor instead')
const DescribeStateMachineInput$json = {
  '1': 'DescribeStateMachineInput',
  '2': [
    {
      '1': 'includeddata',
      '3': 109719114,
      '4': 1,
      '5': 14,
      '6': '.states.IncludedData',
      '10': 'includeddata'
    },
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
  ],
};

/// Descriptor for `DescribeStateMachineInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStateMachineInputDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZVN0YXRlTWFjaGluZUlucHV0EjsKDGluY2x1ZGVkZGF0YRjK3Kg0IAEoDjIULn'
    'N0YXRlcy5JbmNsdWRlZERhdGFSDGluY2x1ZGVkZGF0YRIsCg9zdGF0ZW1hY2hpbmVhcm4Y87vG'
    'uwEgASgJUg9zdGF0ZW1hY2hpbmVhcm4=');

@$core.Deprecated('Use describeStateMachineOutputDescriptor instead')
const DescribeStateMachineOutput$json = {
  '1': 'DescribeStateMachineOutput',
  '2': [
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'definition', '3': 68443297, '4': 1, '5': 9, '10': 'definition'},
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'encryptionconfiguration',
      '3': 167857431,
      '4': 1,
      '5': 11,
      '6': '.states.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {'1': 'label', '3': 379000830, '4': 1, '5': 9, '10': 'label'},
    {
      '1': 'loggingconfiguration',
      '3': 420811605,
      '4': 1,
      '5': 11,
      '6': '.states.LoggingConfiguration',
      '10': 'loggingconfiguration'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'revisionid', '3': 369170086, '4': 1, '5': 9, '10': 'revisionid'},
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.states.StateMachineStatus',
      '10': 'status'
    },
    {
      '1': 'tracingconfiguration',
      '3': 491315910,
      '4': 1,
      '5': 11,
      '6': '.states.TracingConfiguration',
      '10': 'tracingconfiguration'
    },
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.states.StateMachineType',
      '10': 'type'
    },
    {
      '1': 'variablereferences',
      '3': 150942252,
      '4': 3,
      '5': 11,
      '6': '.states.DescribeStateMachineOutput.VariablereferencesEntry',
      '10': 'variablereferences'
    },
  ],
  '3': [DescribeStateMachineOutput_VariablereferencesEntry$json],
};

@$core.Deprecated('Use describeStateMachineOutputDescriptor instead')
const DescribeStateMachineOutput_VariablereferencesEntry$json = {
  '1': 'VariablereferencesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DescribeStateMachineOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStateMachineOutputDescriptor = $convert.base64Decode(
    'ChpEZXNjcmliZVN0YXRlTWFjaGluZU91dHB1dBIlCgxjcmVhdGlvbmRhdGUYgc7RcSABKAlSDG'
    'NyZWF0aW9uZGF0ZRIhCgpkZWZpbml0aW9uGKG50SAgASgJUgpkZWZpbml0aW9uEiQKC2Rlc2Ny'
    'aXB0aW9uGOr2vKMBIAEoCVILZGVzY3JpcHRpb24SXAoXZW5jcnlwdGlvbmNvbmZpZ3VyYXRpb2'
    '4Yl5qFUCABKAsyHy5zdGF0ZXMuRW5jcnlwdGlvbkNvbmZpZ3VyYXRpb25SF2VuY3J5cHRpb25j'
    'b25maWd1cmF0aW9uEhgKBWxhYmVsGP6v3LQBIAEoCVIFbGFiZWwSVAoUbG9nZ2luZ2NvbmZpZ3'
    'VyYXRpb24Y1abUyAEgASgLMhwuc3RhdGVzLkxvZ2dpbmdDb25maWd1cmF0aW9uUhRsb2dnaW5n'
    'Y29uZmlndXJhdGlvbhIVCgRuYW1lGOf75mkgASgJUgRuYW1lEiIKCnJldmlzaW9uaWQYpq2EsA'
    'EgASgJUgpyZXZpc2lvbmlkEhsKB3JvbGVhcm4YoZeJUSABKAlSB3JvbGVhcm4SLAoPc3RhdGVt'
    'YWNoaW5lYXJuGPO7xrsBIAEoCVIPc3RhdGVtYWNoaW5lYXJuEjYKBnN0YXR1cxjw763SASABKA'
    '4yGi5zdGF0ZXMuU3RhdGVNYWNoaW5lU3RhdHVzUgZzdGF0dXMSVAoUdHJhY2luZ2NvbmZpZ3Vy'
    'YXRpb24YxsWj6gEgASgLMhwuc3RhdGVzLlRyYWNpbmdDb25maWd1cmF0aW9uUhR0cmFjaW5nY2'
    '9uZmlndXJhdGlvbhIwCgR0eXBlGM7in4kBIAEoDjIYLnN0YXRlcy5TdGF0ZU1hY2hpbmVUeXBl'
    'UgR0eXBlEm0KEnZhcmlhYmxlcmVmZXJlbmNlcxis5PxHIAMoCzI6LnN0YXRlcy5EZXNjcmliZV'
    'N0YXRlTWFjaGluZU91dHB1dC5WYXJpYWJsZXJlZmVyZW5jZXNFbnRyeVISdmFyaWFibGVyZWZl'
    'cmVuY2VzGkUKF1ZhcmlhYmxlcmVmZXJlbmNlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBX'
    'ZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use encryptionConfigurationDescriptor instead')
const EncryptionConfiguration$json = {
  '1': 'EncryptionConfiguration',
  '2': [
    {
      '1': 'kmsdatakeyreuseperiodseconds',
      '3': 440747764,
      '4': 1,
      '5': 5,
      '10': 'kmsdatakeyreuseperiodseconds'
    },
    {'1': 'kmskeyid', '3': 510698477, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.states.EncryptionType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `EncryptionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptionConfigurationDescriptor = $convert.base64Decode(
    'ChdFbmNyeXB0aW9uQ29uZmlndXJhdGlvbhJGChxrbXNkYXRha2V5cmV1c2VwZXJpb2RzZWNvbm'
    'RzGPSNldIBIAEoBVIca21zZGF0YWtleXJldXNlcGVyaW9kc2Vjb25kcxIeCghrbXNrZXlpZBjt'
    'x8LzASABKAlSCGttc2tleWlkEi4KBHR5cGUYzuKfiQEgASgOMhYuc3RhdGVzLkVuY3J5cHRpb2'
    '5UeXBlUgR0eXBl');

@$core.Deprecated('Use evaluationFailedEventDetailsDescriptor instead')
const EvaluationFailedEventDetails$json = {
  '1': 'EvaluationFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'location', '3': 200649127, '4': 1, '5': 9, '10': 'location'},
    {'1': 'state', '3': 405877495, '4': 1, '5': 9, '10': 'state'},
  ],
};

/// Descriptor for `EvaluationFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List evaluationFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChxFdmFsdWF0aW9uRmFpbGVkRXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0UgASgJUgVjYXVzZR'
        'IXCgVlcnJvchjSjsYMIAEoCVIFZXJyb3ISHQoIbG9jYXRpb24Yp9PWXyABKAlSCGxvY2F0aW9u'
        'EhgKBXN0YXRlGPflxMEBIAEoCVIFc3RhdGU=');

@$core.Deprecated('Use executionAbortedEventDetailsDescriptor instead')
const ExecutionAbortedEventDetails$json = {
  '1': 'ExecutionAbortedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `ExecutionAbortedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionAbortedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChxFeGVjdXRpb25BYm9ydGVkRXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0UgASgJUgVjYXVzZR'
        'IXCgVlcnJvchjSjsYMIAEoCVIFZXJyb3I=');

@$core.Deprecated('Use executionAlreadyExistsDescriptor instead')
const ExecutionAlreadyExists$json = {
  '1': 'ExecutionAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExecutionAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChZFeGVjdXRpb25BbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use executionDoesNotExistDescriptor instead')
const ExecutionDoesNotExist$json = {
  '1': 'ExecutionDoesNotExist',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExecutionDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionDoesNotExistDescriptor = $convert.base64Decode(
    'ChVFeGVjdXRpb25Eb2VzTm90RXhpc3QSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use executionFailedEventDetailsDescriptor instead')
const ExecutionFailedEventDetails$json = {
  '1': 'ExecutionFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `ExecutionFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChtFeGVjdXRpb25GYWlsZWRFdmVudERldGFpbHMSFwoFY2F1c2UYoaS7RSABKAlSBWNhdXNlEh'
        'cKBWVycm9yGNKOxgwgASgJUgVlcnJvcg==');

@$core.Deprecated('Use executionLimitExceededDescriptor instead')
const ExecutionLimitExceeded$json = {
  '1': 'ExecutionLimitExceeded',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExecutionLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionLimitExceededDescriptor =
    $convert.base64Decode(
        'ChZFeGVjdXRpb25MaW1pdEV4Y2VlZGVkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use executionListItemDescriptor instead')
const ExecutionListItem$json = {
  '1': 'ExecutionListItem',
  '2': [
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {'1': 'itemcount', '3': 349613174, '4': 1, '5': 5, '10': 'itemcount'},
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'redrivecount', '3': 473458696, '4': 1, '5': 5, '10': 'redrivecount'},
    {'1': 'redrivedate', '3': 152812125, '4': 1, '5': 9, '10': 'redrivedate'},
    {'1': 'startdate', '3': 364840732, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.states.ExecutionStatus',
      '10': 'status'
    },
    {'1': 'stopdate', '3': 180697434, '4': 1, '5': 9, '10': 'stopdate'},
  ],
};

/// Descriptor for `ExecutionListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionListItemDescriptor = $convert.base64Decode(
    'ChFFeGVjdXRpb25MaXN0SXRlbRImCgxleGVjdXRpb25hcm4Y7Zb9lQEgASgJUgxleGVjdXRpb2'
    '5hcm4SIAoJaXRlbWNvdW50GPbY2qYBIAEoBVIJaXRlbWNvdW50Eh8KCW1hcHJ1bmFybhi669YI'
    'IAEoCVIJbWFwcnVuYXJuEhUKBG5hbWUY5/vmaSABKAlSBG5hbWUSJgoMcmVkcml2ZWNvdW50GI'
    'jQ4eEBIAEoBVIMcmVkcml2ZWNvdW50EiMKC3JlZHJpdmVkYXRlGN307kggASgJUgtyZWRyaXZl'
    'ZGF0ZRIgCglzdGFydGRhdGUYnI78rQEgASgJUglzdGFydGRhdGUSNgoUc3RhdGVtYWNoaW5lYW'
    'xpYXNhcm4YkdTx/AEgASgJUhRzdGF0ZW1hY2hpbmVhbGlhc2FybhIsCg9zdGF0ZW1hY2hpbmVh'
    'cm4Y87vGuwEgASgJUg9zdGF0ZW1hY2hpbmVhcm4SOQoWc3RhdGVtYWNoaW5ldmVyc2lvbmFybh'
    'j5ha8hIAEoCVIWc3RhdGVtYWNoaW5ldmVyc2lvbmFybhIzCgZzdGF0dXMY8O+t0gEgASgOMhcu'
    'c3RhdGVzLkV4ZWN1dGlvblN0YXR1c1IGc3RhdHVzEh0KCHN0b3BkYXRlGNrylFYgASgJUghzdG'
    '9wZGF0ZQ==');

@$core.Deprecated('Use executionNotRedrivableDescriptor instead')
const ExecutionNotRedrivable$json = {
  '1': 'ExecutionNotRedrivable',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExecutionNotRedrivable`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionNotRedrivableDescriptor =
    $convert.base64Decode(
        'ChZFeGVjdXRpb25Ob3RSZWRyaXZhYmxlEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use executionRedrivenEventDetailsDescriptor instead')
const ExecutionRedrivenEventDetails$json = {
  '1': 'ExecutionRedrivenEventDetails',
  '2': [
    {'1': 'redrivecount', '3': 473458696, '4': 1, '5': 5, '10': 'redrivecount'},
  ],
};

/// Descriptor for `ExecutionRedrivenEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionRedrivenEventDetailsDescriptor =
    $convert.base64Decode(
        'Ch1FeGVjdXRpb25SZWRyaXZlbkV2ZW50RGV0YWlscxImCgxyZWRyaXZlY291bnQYiNDh4QEgAS'
        'gFUgxyZWRyaXZlY291bnQ=');

@$core.Deprecated('Use executionStartedEventDetailsDescriptor instead')
const ExecutionStartedEventDetails$json = {
  '1': 'ExecutionStartedEventDetails',
  '2': [
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'inputdetails',
      '3': 452625788,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'inputdetails'
    },
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
  ],
};

/// Descriptor for `ExecutionStartedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionStartedEventDetailsDescriptor = $convert.base64Decode(
    'ChxFeGVjdXRpb25TdGFydGVkRXZlbnREZXRhaWxzEhgKBWlucHV0GPze4c4BIAEoCVIFaW5wdX'
    'QSUAoMaW5wdXRkZXRhaWxzGPyK6tcBIAEoCzIoLnN0YXRlcy5IaXN0b3J5RXZlbnRFeGVjdXRp'
    'b25EYXRhRGV0YWlsc1IMaW5wdXRkZXRhaWxzEhsKB3JvbGVhcm4YoZeJUSABKAlSB3JvbGVhcm'
    '4SNgoUc3RhdGVtYWNoaW5lYWxpYXNhcm4YkdTx/AEgASgJUhRzdGF0ZW1hY2hpbmVhbGlhc2Fy'
    'bhI5ChZzdGF0ZW1hY2hpbmV2ZXJzaW9uYXJuGPmFryEgASgJUhZzdGF0ZW1hY2hpbmV2ZXJzaW'
    '9uYXJu');

@$core.Deprecated('Use executionSucceededEventDetailsDescriptor instead')
const ExecutionSucceededEventDetails$json = {
  '1': 'ExecutionSucceededEventDetails',
  '2': [
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputdetails',
      '3': 393734643,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'outputdetails'
    },
  ],
};

/// Descriptor for `ExecutionSucceededEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionSucceededEventDetailsDescriptor =
    $convert.base64Decode(
        'Ch5FeGVjdXRpb25TdWNjZWVkZWRFdmVudERldGFpbHMSGgoGb3V0cHV0GIWepc0BIAEoCVIGb3'
        'V0cHV0ElIKDW91dHB1dGRldGFpbHMY89PfuwEgASgLMiguc3RhdGVzLkhpc3RvcnlFdmVudEV4'
        'ZWN1dGlvbkRhdGFEZXRhaWxzUg1vdXRwdXRkZXRhaWxz');

@$core.Deprecated('Use executionTimedOutEventDetailsDescriptor instead')
const ExecutionTimedOutEventDetails$json = {
  '1': 'ExecutionTimedOutEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `ExecutionTimedOutEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionTimedOutEventDetailsDescriptor =
    $convert.base64Decode(
        'Ch1FeGVjdXRpb25UaW1lZE91dEV2ZW50RGV0YWlscxIXCgVjYXVzZRihpLtFIAEoCVIFY2F1c2'
        'USFwoFZXJyb3IY0o7GDCABKAlSBWVycm9y');

@$core.Deprecated('Use getActivityTaskInputDescriptor instead')
const GetActivityTaskInput$json = {
  '1': 'GetActivityTaskInput',
  '2': [
    {'1': 'activityarn', '3': 327279492, '4': 1, '5': 9, '10': 'activityarn'},
    {'1': 'workername', '3': 526761781, '4': 1, '5': 9, '10': 'workername'},
  ],
};

/// Descriptor for `GetActivityTaskInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getActivityTaskInputDescriptor = $convert.base64Decode(
    'ChRHZXRBY3Rpdml0eVRhc2tJbnB1dBIkCgthY3Rpdml0eWFybhiEx4ecASABKAlSC2FjdGl2aX'
    'R5YXJuEiIKCndvcmtlcm5hbWUYtf6W+wEgASgJUgp3b3JrZXJuYW1l');

@$core.Deprecated('Use getActivityTaskOutputDescriptor instead')
const GetActivityTaskOutput$json = {
  '1': 'GetActivityTaskOutput',
  '2': [
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {'1': 'tasktoken', '3': 525325834, '4': 1, '5': 9, '10': 'tasktoken'},
  ],
};

/// Descriptor for `GetActivityTaskOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getActivityTaskOutputDescriptor = $convert.base64Decode(
    'ChVHZXRBY3Rpdml0eVRhc2tPdXRwdXQSGAoFaW5wdXQY/N7hzgEgASgJUgVpbnB1dBIgCgl0YX'
    'NrdG9rZW4Yiqy/+gEgASgJUgl0YXNrdG9rZW4=');

@$core.Deprecated('Use getExecutionHistoryInputDescriptor instead')
const GetExecutionHistoryInput$json = {
  '1': 'GetExecutionHistoryInput',
  '2': [
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {
      '1': 'includeexecutiondata',
      '3': 203899608,
      '4': 1,
      '5': 8,
      '10': 'includeexecutiondata'
    },
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'reverseorder', '3': 364411768, '4': 1, '5': 8, '10': 'reverseorder'},
  ],
};

/// Descriptor for `GetExecutionHistoryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getExecutionHistoryInputDescriptor = $convert.base64Decode(
    'ChhHZXRFeGVjdXRpb25IaXN0b3J5SW5wdXQSJgoMZXhlY3V0aW9uYXJuGO2W/ZUBIAEoCVIMZX'
    'hlY3V0aW9uYXJuEjUKFGluY2x1ZGVleGVjdXRpb25kYXRhGNiFnWEgASgIUhRpbmNsdWRlZXhl'
    'Y3V0aW9uZGF0YRIiCgptYXhyZXN1bHRzGNLc590BIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG'
    '9rZW4YnvOdNyABKAlSCW5leHR0b2tlbhImCgxyZXZlcnNlb3JkZXIY+PbhrQEgASgIUgxyZXZl'
    'cnNlb3JkZXI=');

@$core.Deprecated('Use getExecutionHistoryOutputDescriptor instead')
const GetExecutionHistoryOutput$json = {
  '1': 'GetExecutionHistoryOutput',
  '2': [
    {
      '1': 'events',
      '3': 316203909,
      '4': 3,
      '5': 11,
      '6': '.states.HistoryEvent',
      '10': 'events'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `GetExecutionHistoryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getExecutionHistoryOutputDescriptor = $convert.base64Decode(
    'ChlHZXRFeGVjdXRpb25IaXN0b3J5T3V0cHV0EjAKBmV2ZW50cxiFx+OWASADKAsyFC5zdGF0ZX'
    'MuSGlzdG9yeUV2ZW50UgZldmVudHMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use historyEventDescriptor instead')
const HistoryEvent$json = {
  '1': 'HistoryEvent',
  '2': [
    {
      '1': 'activityfailedeventdetails',
      '3': 274703772,
      '4': 1,
      '5': 11,
      '6': '.states.ActivityFailedEventDetails',
      '10': 'activityfailedeventdetails'
    },
    {
      '1': 'activityschedulefailedeventdetails',
      '3': 215835903,
      '4': 1,
      '5': 11,
      '6': '.states.ActivityScheduleFailedEventDetails',
      '10': 'activityschedulefailedeventdetails'
    },
    {
      '1': 'activityscheduledeventdetails',
      '3': 68173374,
      '4': 1,
      '5': 11,
      '6': '.states.ActivityScheduledEventDetails',
      '10': 'activityscheduledeventdetails'
    },
    {
      '1': 'activitystartedeventdetails',
      '3': 25723004,
      '4': 1,
      '5': 11,
      '6': '.states.ActivityStartedEventDetails',
      '10': 'activitystartedeventdetails'
    },
    {
      '1': 'activitysucceededeventdetails',
      '3': 387141788,
      '4': 1,
      '5': 11,
      '6': '.states.ActivitySucceededEventDetails',
      '10': 'activitysucceededeventdetails'
    },
    {
      '1': 'activitytimedouteventdetails',
      '3': 360915296,
      '4': 1,
      '5': 11,
      '6': '.states.ActivityTimedOutEventDetails',
      '10': 'activitytimedouteventdetails'
    },
    {
      '1': 'evaluationfailedeventdetails',
      '3': 490113981,
      '4': 1,
      '5': 11,
      '6': '.states.EvaluationFailedEventDetails',
      '10': 'evaluationfailedeventdetails'
    },
    {
      '1': 'executionabortedeventdetails',
      '3': 507567407,
      '4': 1,
      '5': 11,
      '6': '.states.ExecutionAbortedEventDetails',
      '10': 'executionabortedeventdetails'
    },
    {
      '1': 'executionfailedeventdetails',
      '3': 297327899,
      '4': 1,
      '5': 11,
      '6': '.states.ExecutionFailedEventDetails',
      '10': 'executionfailedeventdetails'
    },
    {
      '1': 'executionredriveneventdetails',
      '3': 157425989,
      '4': 1,
      '5': 11,
      '6': '.states.ExecutionRedrivenEventDetails',
      '10': 'executionredriveneventdetails'
    },
    {
      '1': 'executionstartedeventdetails',
      '3': 382892969,
      '4': 1,
      '5': 11,
      '6': '.states.ExecutionStartedEventDetails',
      '10': 'executionstartedeventdetails'
    },
    {
      '1': 'executionsucceededeventdetails',
      '3': 375748869,
      '4': 1,
      '5': 11,
      '6': '.states.ExecutionSucceededEventDetails',
      '10': 'executionsucceededeventdetails'
    },
    {
      '1': 'executiontimedouteventdetails',
      '3': 139990699,
      '4': 1,
      '5': 11,
      '6': '.states.ExecutionTimedOutEventDetails',
      '10': 'executiontimedouteventdetails'
    },
    {'1': 'id', '3': 389573345, '4': 1, '5': 3, '10': 'id'},
    {
      '1': 'lambdafunctionfailedeventdetails',
      '3': 31352398,
      '4': 1,
      '5': 11,
      '6': '.states.LambdaFunctionFailedEventDetails',
      '10': 'lambdafunctionfailedeventdetails'
    },
    {
      '1': 'lambdafunctionschedulefailedeventdetails',
      '3': 449961745,
      '4': 1,
      '5': 11,
      '6': '.states.LambdaFunctionScheduleFailedEventDetails',
      '10': 'lambdafunctionschedulefailedeventdetails'
    },
    {
      '1': 'lambdafunctionscheduledeventdetails',
      '3': 27704864,
      '4': 1,
      '5': 11,
      '6': '.states.LambdaFunctionScheduledEventDetails',
      '10': 'lambdafunctionscheduledeventdetails'
    },
    {
      '1': 'lambdafunctionstartfailedeventdetails',
      '3': 508021246,
      '4': 1,
      '5': 11,
      '6': '.states.LambdaFunctionStartFailedEventDetails',
      '10': 'lambdafunctionstartfailedeventdetails'
    },
    {
      '1': 'lambdafunctionsucceededeventdetails',
      '3': 347724058,
      '4': 1,
      '5': 11,
      '6': '.states.LambdaFunctionSucceededEventDetails',
      '10': 'lambdafunctionsucceededeventdetails'
    },
    {
      '1': 'lambdafunctiontimedouteventdetails',
      '3': 423888938,
      '4': 1,
      '5': 11,
      '6': '.states.LambdaFunctionTimedOutEventDetails',
      '10': 'lambdafunctiontimedouteventdetails'
    },
    {
      '1': 'mapiterationabortedeventdetails',
      '3': 465966882,
      '4': 1,
      '5': 11,
      '6': '.states.MapIterationEventDetails',
      '10': 'mapiterationabortedeventdetails'
    },
    {
      '1': 'mapiterationfailedeventdetails',
      '3': 356461616,
      '4': 1,
      '5': 11,
      '6': '.states.MapIterationEventDetails',
      '10': 'mapiterationfailedeventdetails'
    },
    {
      '1': 'mapiterationstartedeventdetails',
      '3': 294963144,
      '4': 1,
      '5': 11,
      '6': '.states.MapIterationEventDetails',
      '10': 'mapiterationstartedeventdetails'
    },
    {
      '1': 'mapiterationsucceededeventdetails',
      '3': 85930040,
      '4': 1,
      '5': 11,
      '6': '.states.MapIterationEventDetails',
      '10': 'mapiterationsucceededeventdetails'
    },
    {
      '1': 'maprunfailedeventdetails',
      '3': 30147502,
      '4': 1,
      '5': 11,
      '6': '.states.MapRunFailedEventDetails',
      '10': 'maprunfailedeventdetails'
    },
    {
      '1': 'maprunredriveneventdetails',
      '3': 465066264,
      '4': 1,
      '5': 11,
      '6': '.states.MapRunRedrivenEventDetails',
      '10': 'maprunredriveneventdetails'
    },
    {
      '1': 'maprunstartedeventdetails',
      '3': 183119902,
      '4': 1,
      '5': 11,
      '6': '.states.MapRunStartedEventDetails',
      '10': 'maprunstartedeventdetails'
    },
    {
      '1': 'mapstatestartedeventdetails',
      '3': 463022434,
      '4': 1,
      '5': 11,
      '6': '.states.MapStateStartedEventDetails',
      '10': 'mapstatestartedeventdetails'
    },
    {
      '1': 'previouseventid',
      '3': 261811790,
      '4': 1,
      '5': 3,
      '10': 'previouseventid'
    },
    {
      '1': 'stateenteredeventdetails',
      '3': 479306524,
      '4': 1,
      '5': 11,
      '6': '.states.StateEnteredEventDetails',
      '10': 'stateenteredeventdetails'
    },
    {
      '1': 'stateexitedeventdetails',
      '3': 30960782,
      '4': 1,
      '5': 11,
      '6': '.states.StateExitedEventDetails',
      '10': 'stateexitedeventdetails'
    },
    {
      '1': 'taskfailedeventdetails',
      '3': 309753628,
      '4': 1,
      '5': 11,
      '6': '.states.TaskFailedEventDetails',
      '10': 'taskfailedeventdetails'
    },
    {
      '1': 'taskscheduledeventdetails',
      '3': 238821054,
      '4': 1,
      '5': 11,
      '6': '.states.TaskScheduledEventDetails',
      '10': 'taskscheduledeventdetails'
    },
    {
      '1': 'taskstartfailedeventdetails',
      '3': 237771688,
      '4': 1,
      '5': 11,
      '6': '.states.TaskStartFailedEventDetails',
      '10': 'taskstartfailedeventdetails'
    },
    {
      '1': 'taskstartedeventdetails',
      '3': 192171260,
      '4': 1,
      '5': 11,
      '6': '.states.TaskStartedEventDetails',
      '10': 'taskstartedeventdetails'
    },
    {
      '1': 'tasksubmitfailedeventdetails',
      '3': 489315884,
      '4': 1,
      '5': 11,
      '6': '.states.TaskSubmitFailedEventDetails',
      '10': 'tasksubmitfailedeventdetails'
    },
    {
      '1': 'tasksubmittedeventdetails',
      '3': 269619140,
      '4': 1,
      '5': 11,
      '6': '.states.TaskSubmittedEventDetails',
      '10': 'tasksubmittedeventdetails'
    },
    {
      '1': 'tasksucceededeventdetails',
      '3': 20918556,
      '4': 1,
      '5': 11,
      '6': '.states.TaskSucceededEventDetails',
      '10': 'tasksucceededeventdetails'
    },
    {
      '1': 'tasktimedouteventdetails',
      '3': 330698464,
      '4': 1,
      '5': 11,
      '6': '.states.TaskTimedOutEventDetails',
      '10': 'tasktimedouteventdetails'
    },
    {'1': 'timestamp', '3': 310629668, '4': 1, '5': 9, '10': 'timestamp'},
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.states.HistoryEventType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `HistoryEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List historyEventDescriptor = $convert.base64Decode(
    'CgxIaXN0b3J5RXZlbnQSZgoaYWN0aXZpdHlmYWlsZWRldmVudGRldGFpbHMYnMv+ggEgASgLMi'
    'Iuc3RhdGVzLkFjdGl2aXR5RmFpbGVkRXZlbnREZXRhaWxzUhphY3Rpdml0eWZhaWxlZGV2ZW50'
    'ZGV0YWlscxJ9CiJhY3Rpdml0eXNjaGVkdWxlZmFpbGVkZXZlbnRkZXRhaWxzGP/J9WYgASgLMi'
    'ouc3RhdGVzLkFjdGl2aXR5U2NoZWR1bGVGYWlsZWRFdmVudERldGFpbHNSImFjdGl2aXR5c2No'
    'ZWR1bGVmYWlsZWRldmVudGRldGFpbHMSbgodYWN0aXZpdHlzY2hlZHVsZWRldmVudGRldGFpbH'
    'MYvvzAICABKAsyJS5zdGF0ZXMuQWN0aXZpdHlTY2hlZHVsZWRFdmVudERldGFpbHNSHWFjdGl2'
    'aXR5c2NoZWR1bGVkZXZlbnRkZXRhaWxzEmgKG2FjdGl2aXR5c3RhcnRlZGV2ZW50ZGV0YWlscx'
    'j8gKIMIAEoCzIjLnN0YXRlcy5BY3Rpdml0eVN0YXJ0ZWRFdmVudERldGFpbHNSG2FjdGl2aXR5'
    'c3RhcnRlZGV2ZW50ZGV0YWlscxJvCh1hY3Rpdml0eXN1Y2NlZWRlZGV2ZW50ZGV0YWlscxicoc'
    '24ASABKAsyJS5zdGF0ZXMuQWN0aXZpdHlTdWNjZWVkZWRFdmVudERldGFpbHNSHWFjdGl2aXR5'
    'c3VjY2VlZGVkZXZlbnRkZXRhaWxzEmwKHGFjdGl2aXR5dGltZWRvdXRldmVudGRldGFpbHMY4M'
    'KMrAEgASgLMiQuc3RhdGVzLkFjdGl2aXR5VGltZWRPdXRFdmVudERldGFpbHNSHGFjdGl2aXR5'
    'dGltZWRvdXRldmVudGRldGFpbHMSbAocZXZhbHVhdGlvbmZhaWxlZGV2ZW50ZGV0YWlscxi9l9'
    'rpASABKAsyJC5zdGF0ZXMuRXZhbHVhdGlvbkZhaWxlZEV2ZW50RGV0YWlsc1IcZXZhbHVhdGlv'
    'bmZhaWxlZGV2ZW50ZGV0YWlscxJsChxleGVjdXRpb25hYm9ydGVkZXZlbnRkZXRhaWxzGK+6g/'
    'IBIAEoCzIkLnN0YXRlcy5FeGVjdXRpb25BYm9ydGVkRXZlbnREZXRhaWxzUhxleGVjdXRpb25h'
    'Ym9ydGVkZXZlbnRkZXRhaWxzEmkKG2V4ZWN1dGlvbmZhaWxlZGV2ZW50ZGV0YWlscxibuuONAS'
    'ABKAsyIy5zdGF0ZXMuRXhlY3V0aW9uRmFpbGVkRXZlbnREZXRhaWxzUhtleGVjdXRpb25mYWls'
    'ZWRldmVudGRldGFpbHMSbgodZXhlY3V0aW9ucmVkcml2ZW5ldmVudGRldGFpbHMYxcKISyABKA'
    'syJS5zdGF0ZXMuRXhlY3V0aW9uUmVkcml2ZW5FdmVudERldGFpbHNSHWV4ZWN1dGlvbnJlZHJp'
    'dmVuZXZlbnRkZXRhaWxzEmwKHGV4ZWN1dGlvbnN0YXJ0ZWRldmVudGRldGFpbHMYqffJtgEgAS'
    'gLMiQuc3RhdGVzLkV4ZWN1dGlvblN0YXJ0ZWRFdmVudERldGFpbHNSHGV4ZWN1dGlvbnN0YXJ0'
    'ZWRldmVudGRldGFpbHMScgoeZXhlY3V0aW9uc3VjY2VlZGVkZXZlbnRkZXRhaWxzGIXylbMBIA'
    'EoCzImLnN0YXRlcy5FeGVjdXRpb25TdWNjZWVkZWRFdmVudERldGFpbHNSHmV4ZWN1dGlvbnN1'
    'Y2NlZWRlZGV2ZW50ZGV0YWlscxJuCh1leGVjdXRpb250aW1lZG91dGV2ZW50ZGV0YWlscxirre'
    'BCIAEoCzIlLnN0YXRlcy5FeGVjdXRpb25UaW1lZE91dEV2ZW50RGV0YWlsc1IdZXhlY3V0aW9u'
    'dGltZWRvdXRldmVudGRldGFpbHMSEgoCaWQY4dXhuQEgASgDUgJpZBJ3CiBsYW1iZGFmdW5jdG'
    'lvbmZhaWxlZGV2ZW50ZGV0YWlscxjOzPkOIAEoCzIoLnN0YXRlcy5MYW1iZGFGdW5jdGlvbkZh'
    'aWxlZEV2ZW50RGV0YWlsc1IgbGFtYmRhZnVuY3Rpb25mYWlsZWRldmVudGRldGFpbHMSkAEKKG'
    'xhbWJkYWZ1bmN0aW9uc2NoZWR1bGVmYWlsZWRldmVudGRldGFpbHMYkb7H1gEgASgLMjAuc3Rh'
    'dGVzLkxhbWJkYUZ1bmN0aW9uU2NoZWR1bGVGYWlsZWRFdmVudERldGFpbHNSKGxhbWJkYWZ1bm'
    'N0aW9uc2NoZWR1bGVmYWlsZWRldmVudGRldGFpbHMSgAEKI2xhbWJkYWZ1bmN0aW9uc2NoZWR1'
    'bGVkZXZlbnRkZXRhaWxzGKD8mg0gASgLMisuc3RhdGVzLkxhbWJkYUZ1bmN0aW9uU2NoZWR1bG'
    'VkRXZlbnREZXRhaWxzUiNsYW1iZGFmdW5jdGlvbnNjaGVkdWxlZGV2ZW50ZGV0YWlscxKHAQol'
    'bGFtYmRhZnVuY3Rpb25zdGFydGZhaWxlZGV2ZW50ZGV0YWlscxj+k5/yASABKAsyLS5zdGF0ZX'
    'MuTGFtYmRhRnVuY3Rpb25TdGFydEZhaWxlZEV2ZW50RGV0YWlsc1IlbGFtYmRhZnVuY3Rpb25z'
    'dGFydGZhaWxlZGV2ZW50ZGV0YWlscxKBAQojbGFtYmRhZnVuY3Rpb25zdWNjZWVkZWRldmVudG'
    'RldGFpbHMYmrLnpQEgASgLMisuc3RhdGVzLkxhbWJkYUZ1bmN0aW9uU3VjY2VlZGVkRXZlbnRE'
    'ZXRhaWxzUiNsYW1iZGFmdW5jdGlvbnN1Y2NlZWRlZGV2ZW50ZGV0YWlscxJ+CiJsYW1iZGFmdW'
    '5jdGlvbnRpbWVkb3V0ZXZlbnRkZXRhaWxzGKqQkMoBIAEoCzIqLnN0YXRlcy5MYW1iZGFGdW5j'
    'dGlvblRpbWVkT3V0RXZlbnREZXRhaWxzUiJsYW1iZGFmdW5jdGlvbnRpbWVkb3V0ZXZlbnRkZX'
    'RhaWxzEm4KH21hcGl0ZXJhdGlvbmFib3J0ZWRldmVudGRldGFpbHMYoq6Y3gEgASgLMiAuc3Rh'
    'dGVzLk1hcEl0ZXJhdGlvbkV2ZW50RGV0YWlsc1IfbWFwaXRlcmF0aW9uYWJvcnRlZGV2ZW50ZG'
    'V0YWlscxJsCh5tYXBpdGVyYXRpb25mYWlsZWRldmVudGRldGFpbHMYsNj8qQEgASgLMiAuc3Rh'
    'dGVzLk1hcEl0ZXJhdGlvbkV2ZW50RGV0YWlsc1IebWFwaXRlcmF0aW9uZmFpbGVkZXZlbnRkZX'
    'RhaWxzEm4KH21hcGl0ZXJhdGlvbnN0YXJ0ZWRldmVudGRldGFpbHMYyI/TjAEgASgLMiAuc3Rh'
    'dGVzLk1hcEl0ZXJhdGlvbkV2ZW50RGV0YWlsc1IfbWFwaXRlcmF0aW9uc3RhcnRlZGV2ZW50ZG'
    'V0YWlscxJxCiFtYXBpdGVyYXRpb25zdWNjZWVkZWRldmVudGRldGFpbHMYuOD8KCABKAsyIC5z'
    'dGF0ZXMuTWFwSXRlcmF0aW9uRXZlbnREZXRhaWxzUiFtYXBpdGVyYXRpb25zdWNjZWVkZWRldm'
    'VudGRldGFpbHMSXwoYbWFwcnVuZmFpbGVkZXZlbnRkZXRhaWxzGK6HsA4gASgLMiAuc3RhdGVz'
    'Lk1hcFJ1bkZhaWxlZEV2ZW50RGV0YWlsc1IYbWFwcnVuZmFpbGVkZXZlbnRkZXRhaWxzEmYKGm'
    '1hcHJ1bnJlZHJpdmVuZXZlbnRkZXRhaWxzGJiy4d0BIAEoCzIiLnN0YXRlcy5NYXBSdW5SZWRy'
    'aXZlbkV2ZW50RGV0YWlsc1IabWFwcnVucmVkcml2ZW5ldmVudGRldGFpbHMSYgoZbWFwcnVuc3'
    'RhcnRlZGV2ZW50ZGV0YWlscxie4KhXIAEoCzIhLnN0YXRlcy5NYXBSdW5TdGFydGVkRXZlbnRE'
    'ZXRhaWxzUhltYXBydW5zdGFydGVkZXZlbnRkZXRhaWxzEmkKG21hcHN0YXRlc3RhcnRlZGV2ZW'
    '50ZGV0YWlscxji0uTcASABKAsyIy5zdGF0ZXMuTWFwU3RhdGVTdGFydGVkRXZlbnREZXRhaWxz'
    'UhttYXBzdGF0ZXN0YXJ0ZWRldmVudGRldGFpbHMSKwoPcHJldmlvdXNldmVudGlkGM7c63wgAS'
    'gDUg9wcmV2aW91c2V2ZW50aWQSYAoYc3RhdGVlbnRlcmVkZXZlbnRkZXRhaWxzGJzGxuQBIAEo'
    'CzIgLnN0YXRlcy5TdGF0ZUVudGVyZWRFdmVudERldGFpbHNSGHN0YXRlZW50ZXJlZGV2ZW50ZG'
    'V0YWlscxJcChdzdGF0ZWV4aXRlZGV2ZW50ZGV0YWlscxiO2eEOIAEoCzIfLnN0YXRlcy5TdGF0'
    'ZUV4aXRlZEV2ZW50RGV0YWlsc1IXc3RhdGVleGl0ZWRldmVudGRldGFpbHMSWgoWdGFza2ZhaW'
    'xlZGV2ZW50ZGV0YWlscxic7tmTASABKAsyHi5zdGF0ZXMuVGFza0ZhaWxlZEV2ZW50RGV0YWls'
    'c1IWdGFza2ZhaWxlZGV2ZW50ZGV0YWlscxJiChl0YXNrc2NoZWR1bGVkZXZlbnRkZXRhaWxzGL'
    '698HEgASgLMiEuc3RhdGVzLlRhc2tTY2hlZHVsZWRFdmVudERldGFpbHNSGXRhc2tzY2hlZHVs'
    'ZWRldmVudGRldGFpbHMSaAobdGFza3N0YXJ0ZmFpbGVkZXZlbnRkZXRhaWxzGKi3sHEgASgLMi'
    'Muc3RhdGVzLlRhc2tTdGFydEZhaWxlZEV2ZW50RGV0YWlsc1IbdGFza3N0YXJ0ZmFpbGVkZXZl'
    'bnRkZXRhaWxzElwKF3Rhc2tzdGFydGVkZXZlbnRkZXRhaWxzGPyZ0VsgASgLMh8uc3RhdGVzLl'
    'Rhc2tTdGFydGVkRXZlbnREZXRhaWxzUhd0YXNrc3RhcnRlZGV2ZW50ZGV0YWlscxJsChx0YXNr'
    'c3VibWl0ZmFpbGVkZXZlbnRkZXRhaWxzGKy8qekBIAEoCzIkLnN0YXRlcy5UYXNrU3VibWl0Rm'
    'FpbGVkRXZlbnREZXRhaWxzUhx0YXNrc3VibWl0ZmFpbGVkZXZlbnRkZXRhaWxzEmMKGXRhc2tz'
    'dWJtaXR0ZWRldmVudGRldGFpbHMYxJ/IgAEgASgLMiEuc3RhdGVzLlRhc2tTdWJtaXR0ZWRFdm'
    'VudERldGFpbHNSGXRhc2tzdWJtaXR0ZWRldmVudGRldGFpbHMSYgoZdGFza3N1Y2NlZWRlZGV2'
    'ZW50ZGV0YWlscxic4vwJIAEoCzIhLnN0YXRlcy5UYXNrU3VjY2VlZGVkRXZlbnREZXRhaWxzUh'
    'l0YXNrc3VjY2VlZGVkZXZlbnRkZXRhaWxzEmAKGHRhc2t0aW1lZG91dGV2ZW50ZGV0YWlscxjg'
    'ndidASABKAsyIC5zdGF0ZXMuVGFza1RpbWVkT3V0RXZlbnREZXRhaWxzUhh0YXNrdGltZWRvdX'
    'RldmVudGRldGFpbHMSIAoJdGltZXN0YW1wGKSqj5QBIAEoCVIJdGltZXN0YW1wEjAKBHR5cGUY'
    'zuKfiQEgASgOMhguc3RhdGVzLkhpc3RvcnlFdmVudFR5cGVSBHR5cGU=');

@$core.Deprecated('Use historyEventExecutionDataDetailsDescriptor instead')
const HistoryEventExecutionDataDetails$json = {
  '1': 'HistoryEventExecutionDataDetails',
  '2': [
    {'1': 'truncated', '3': 407490858, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `HistoryEventExecutionDataDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List historyEventExecutionDataDetailsDescriptor =
    $convert.base64Decode(
        'CiBIaXN0b3J5RXZlbnRFeGVjdXRpb25EYXRhRGV0YWlscxIgCgl0cnVuY2F0ZWQYqqKnwgEgAS'
        'gIUgl0cnVuY2F0ZWQ=');

@$core.Deprecated('Use inspectionDataDescriptor instead')
const InspectionData$json = {
  '1': 'InspectionData',
  '2': [
    {
      '1': 'afterarguments',
      '3': 365960236,
      '4': 1,
      '5': 9,
      '10': 'afterarguments'
    },
    {
      '1': 'afterinputpath',
      '3': 355404745,
      '4': 1,
      '5': 9,
      '10': 'afterinputpath'
    },
    {
      '1': 'afteritembatcher',
      '3': 181123354,
      '4': 1,
      '5': 9,
      '10': 'afteritembatcher'
    },
    {
      '1': 'afteritemselector',
      '3': 394258120,
      '4': 1,
      '5': 9,
      '10': 'afteritemselector'
    },
    {
      '1': 'afteritemspath',
      '3': 491034053,
      '4': 1,
      '5': 9,
      '10': 'afteritemspath'
    },
    {
      '1': 'afteritemspointer',
      '3': 194882729,
      '4': 1,
      '5': 9,
      '10': 'afteritemspointer'
    },
    {
      '1': 'afterparameters',
      '3': 328015918,
      '4': 1,
      '5': 9,
      '10': 'afterparameters'
    },
    {
      '1': 'afterresultpath',
      '3': 480099136,
      '4': 1,
      '5': 9,
      '10': 'afterresultpath'
    },
    {
      '1': 'afterresultselector',
      '3': 443240414,
      '4': 1,
      '5': 9,
      '10': 'afterresultselector'
    },
    {
      '1': 'errordetails',
      '3': 192899050,
      '4': 1,
      '5': 11,
      '6': '.states.InspectionErrorDetails',
      '10': 'errordetails'
    },
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'maxconcurrency',
      '3': 100901405,
      '4': 1,
      '5': 5,
      '10': 'maxconcurrency'
    },
    {
      '1': 'request',
      '3': 514460083,
      '4': 1,
      '5': 11,
      '6': '.states.InspectionDataRequest',
      '10': 'request'
    },
    {
      '1': 'response',
      '3': 425574879,
      '4': 1,
      '5': 11,
      '6': '.states.InspectionDataResponse',
      '10': 'response'
    },
    {'1': 'result', '3': 171406885, '4': 1, '5': 9, '10': 'result'},
    {
      '1': 'toleratedfailurecount',
      '3': 41834811,
      '4': 1,
      '5': 5,
      '10': 'toleratedfailurecount'
    },
    {
      '1': 'toleratedfailurepercentage',
      '3': 116496164,
      '4': 1,
      '5': 2,
      '10': 'toleratedfailurepercentage'
    },
    {'1': 'variables', '3': 162226883, '4': 1, '5': 9, '10': 'variables'},
  ],
};

/// Descriptor for `InspectionData`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inspectionDataDescriptor = $convert.base64Decode(
    'Cg5JbnNwZWN0aW9uRGF0YRIqCg5hZnRlcmFyZ3VtZW50cxisuMCuASABKAlSDmFmdGVyYXJndW'
    '1lbnRzEioKDmFmdGVyaW5wdXRwYXRoGMmXvKkBIAEoCVIOYWZ0ZXJpbnB1dHBhdGgSLQoQYWZ0'
    'ZXJpdGVtYmF0Y2hlchia8q5WIAEoCVIQYWZ0ZXJpdGVtYmF0Y2hlchIwChFhZnRlcml0ZW1zZW'
    'xlY3RvchjIzf+7ASABKAlSEWFmdGVyaXRlbXNlbGVjdG9yEioKDmFmdGVyaXRlbXNwYXRoGMWr'
    'kuoBIAEoCVIOYWZ0ZXJpdGVtc3BhdGgSLwoRYWZ0ZXJpdGVtc3BvaW50ZXIYqdn2XCABKAlSEW'
    'FmdGVyaXRlbXNwb2ludGVyEiwKD2FmdGVycGFyYW1ldGVycxiuwLScASABKAlSD2FmdGVycGFy'
    'YW1ldGVycxIsCg9hZnRlcnJlc3VsdHBhdGgYwPb25AEgASgJUg9hZnRlcnJlc3VsdHBhdGgSNA'
    'oTYWZ0ZXJyZXN1bHRzZWxlY3Rvchjen63TASABKAlSE2FmdGVycmVzdWx0c2VsZWN0b3ISRQoM'
    'ZXJyb3JkZXRhaWxzGOrP/VsgASgLMh4uc3RhdGVzLkluc3BlY3Rpb25FcnJvckRldGFpbHNSDG'
    'Vycm9yZGV0YWlscxIYCgVpbnB1dBj83uHOASABKAlSBWlucHV0EikKDm1heGNvbmN1cnJlbmN5'
    'GJ3EjjAgASgFUg5tYXhjb25jdXJyZW5jeRI7CgdyZXF1ZXN0GLOTqPUBIAEoCzIdLnN0YXRlcy'
    '5JbnNwZWN0aW9uRGF0YVJlcXVlc3RSB3JlcXVlc3QSPgoIcmVzcG9uc2UY34P3ygEgASgLMh4u'
    'c3RhdGVzLkluc3BlY3Rpb25EYXRhUmVzcG9uc2VSCHJlc3BvbnNlEhkKBnJlc3VsdBil7N1RIA'
    'EoCVIGcmVzdWx0EjcKFXRvbGVyYXRlZGZhaWx1cmVjb3VudBi7svkTIAEoBVIVdG9sZXJhdGVk'
    'ZmFpbHVyZWNvdW50EkEKGnRvbGVyYXRlZGZhaWx1cmVwZXJjZW50YWdlGKSuxjcgASgCUhp0b2'
    'xlcmF0ZWRmYWlsdXJlcGVyY2VudGFnZRIfCgl2YXJpYWJsZXMYw8WtTSABKAlSCXZhcmlhYmxl'
    'cw==');

@$core.Deprecated('Use inspectionDataRequestDescriptor instead')
const InspectionDataRequest$json = {
  '1': 'InspectionDataRequest',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 9, '10': 'body'},
    {'1': 'headers', '3': 375773674, '4': 1, '5': 9, '10': 'headers'},
    {'1': 'method', '3': 189134641, '4': 1, '5': 9, '10': 'method'},
    {'1': 'protocol', '3': 455607734, '4': 1, '5': 9, '10': 'protocol'},
    {'1': 'url', '3': 311381023, '4': 1, '5': 9, '10': 'url'},
  ],
};

/// Descriptor for `InspectionDataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inspectionDataRequestDescriptor = $convert.base64Decode(
    'ChVJbnNwZWN0aW9uRGF0YVJlcXVlc3QSFgoEYm9keRj28qndASABKAlSBGJvZHkSHAoHaGVhZG'
    'Vycxjqs5ezASABKAlSB2hlYWRlcnMSGQoGbWV0aG9kGLHul1ogASgJUgZtZXRob2QSHgoIcHJv'
    'dG9jb2wYtoug2QEgASgJUghwcm90b2NvbBIUCgN1cmwYn5i9lAEgASgJUgN1cmw=');

@$core.Deprecated('Use inspectionDataResponseDescriptor instead')
const InspectionDataResponse$json = {
  '1': 'InspectionDataResponse',
  '2': [
    {'1': 'body', '3': 464157046, '4': 1, '5': 9, '10': 'body'},
    {'1': 'headers', '3': 375773674, '4': 1, '5': 9, '10': 'headers'},
    {'1': 'protocol', '3': 455607734, '4': 1, '5': 9, '10': 'protocol'},
    {'1': 'statuscode', '3': 299352223, '4': 1, '5': 9, '10': 'statuscode'},
    {
      '1': 'statusmessage',
      '3': 474462255,
      '4': 1,
      '5': 9,
      '10': 'statusmessage'
    },
  ],
};

/// Descriptor for `InspectionDataResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inspectionDataResponseDescriptor = $convert.base64Decode(
    'ChZJbnNwZWN0aW9uRGF0YVJlc3BvbnNlEhYKBGJvZHkY9vKp3QEgASgJUgRib2R5EhwKB2hlYW'
    'RlcnMY6rOXswEgASgJUgdoZWFkZXJzEh4KCHByb3RvY29sGLaLoNkBIAEoCVIIcHJvdG9jb2wS'
    'IgoKc3RhdHVzY29kZRifgd+OASABKAlSCnN0YXR1c2NvZGUSKAoNc3RhdHVzbWVzc2FnZRiv8J'
    '7iASABKAlSDXN0YXR1c21lc3NhZ2U=');

@$core.Deprecated('Use inspectionErrorDetailsDescriptor instead')
const InspectionErrorDetails$json = {
  '1': 'InspectionErrorDetails',
  '2': [
    {'1': 'catchindex', '3': 290365641, '4': 1, '5': 5, '10': 'catchindex'},
    {
      '1': 'retrybackoffintervalseconds',
      '3': 137754128,
      '4': 1,
      '5': 5,
      '10': 'retrybackoffintervalseconds'
    },
    {'1': 'retryindex', '3': 6599888, '4': 1, '5': 5, '10': 'retryindex'},
  ],
};

/// Descriptor for `InspectionErrorDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inspectionErrorDetailsDescriptor = $convert.base64Decode(
    'ChZJbnNwZWN0aW9uRXJyb3JEZXRhaWxzEiIKCmNhdGNoaW5kZXgYycG6igEgASgFUgpjYXRjaG'
    'luZGV4EkMKG3JldHJ5YmFja29mZmludGVydmFsc2Vjb25kcxiQ7NdBIAEoBVIbcmV0cnliYWNr'
    'b2ZmaW50ZXJ2YWxzZWNvbmRzEiEKCnJldHJ5aW5kZXgY0OmSAyABKAVSCnJldHJ5aW5kZXg=');

@$core.Deprecated('Use invalidArnDescriptor instead')
const InvalidArn$json = {
  '1': 'InvalidArn',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidArn`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidArnDescriptor = $convert
    .base64Decode('CgpJbnZhbGlkQXJuEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidDefinitionDescriptor instead')
const InvalidDefinition$json = {
  '1': 'InvalidDefinition',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDefinition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDefinitionDescriptor = $convert.base64Decode(
    'ChFJbnZhbGlkRGVmaW5pdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidEncryptionConfigurationDescriptor instead')
const InvalidEncryptionConfiguration$json = {
  '1': 'InvalidEncryptionConfiguration',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEncryptionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEncryptionConfigurationDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkRW5jcnlwdGlvbkNvbmZpZ3VyYXRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use invalidExecutionInputDescriptor instead')
const InvalidExecutionInput$json = {
  '1': 'InvalidExecutionInput',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidExecutionInputDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkRXhlY3V0aW9uSW5wdXQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidLoggingConfigurationDescriptor instead')
const InvalidLoggingConfiguration$json = {
  '1': 'InvalidLoggingConfiguration',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidLoggingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidLoggingConfigurationDescriptor =
    $convert.base64Decode(
        'ChtJbnZhbGlkTG9nZ2luZ0NvbmZpZ3VyYXRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use invalidNameDescriptor instead')
const InvalidName$json = {
  '1': 'InvalidName',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidName`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidNameDescriptor = $convert
    .base64Decode('CgtJbnZhbGlkTmFtZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidOutputDescriptor instead')
const InvalidOutput$json = {
  '1': 'InvalidOutput',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOutputDescriptor = $convert.base64Decode(
    'Cg1JbnZhbGlkT3V0cHV0EhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidTokenDescriptor instead')
const InvalidToken$json = {
  '1': 'InvalidToken',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidToken`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTokenDescriptor = $convert.base64Decode(
    'CgxJbnZhbGlkVG9rZW4SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidTracingConfigurationDescriptor instead')
const InvalidTracingConfiguration$json = {
  '1': 'InvalidTracingConfiguration',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTracingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTracingConfigurationDescriptor =
    $convert.base64Decode(
        'ChtJbnZhbGlkVHJhY2luZ0NvbmZpZ3VyYXRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use kmsAccessDeniedExceptionDescriptor instead')
const KmsAccessDeniedException$json = {
  '1': 'KmsAccessDeniedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsAccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsAccessDeniedExceptionDescriptor =
    $convert.base64Decode(
        'ChhLbXNBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use kmsInvalidStateExceptionDescriptor instead')
const KmsInvalidStateException$json = {
  '1': 'KmsInvalidStateException',
  '2': [
    {
      '1': 'kmskeystate',
      '3': 485859435,
      '4': 1,
      '5': 14,
      '6': '.states.KmsKeyState',
      '10': 'kmskeystate'
    },
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsInvalidStateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsInvalidStateExceptionDescriptor = $convert.base64Decode(
    'ChhLbXNJbnZhbGlkU3RhdGVFeGNlcHRpb24SOQoLa21za2V5c3RhdGUY68DW5wEgASgOMhMuc3'
    'RhdGVzLkttc0tleVN0YXRlUgtrbXNrZXlzdGF0ZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNz'
    'YWdl');

@$core.Deprecated('Use kmsThrottlingExceptionDescriptor instead')
const KmsThrottlingException$json = {
  '1': 'KmsThrottlingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsThrottlingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsThrottlingExceptionDescriptor =
    $convert.base64Decode(
        'ChZLbXNUaHJvdHRsaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use lambdaFunctionFailedEventDetailsDescriptor instead')
const LambdaFunctionFailedEventDetails$json = {
  '1': 'LambdaFunctionFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `LambdaFunctionFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'CiBMYW1iZGFGdW5jdGlvbkZhaWxlZEV2ZW50RGV0YWlscxIXCgVjYXVzZRihpLtFIAEoCVIFY2'
        'F1c2USFwoFZXJyb3IY0o7GDCABKAlSBWVycm9y');

@$core.Deprecated(
    'Use lambdaFunctionScheduleFailedEventDetailsDescriptor instead')
const LambdaFunctionScheduleFailedEventDetails$json = {
  '1': 'LambdaFunctionScheduleFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `LambdaFunctionScheduleFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionScheduleFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'CihMYW1iZGFGdW5jdGlvblNjaGVkdWxlRmFpbGVkRXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0'
        'UgASgJUgVjYXVzZRIXCgVlcnJvchjSjsYMIAEoCVIFZXJyb3I=');

@$core.Deprecated('Use lambdaFunctionScheduledEventDetailsDescriptor instead')
const LambdaFunctionScheduledEventDetails$json = {
  '1': 'LambdaFunctionScheduledEventDetails',
  '2': [
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'inputdetails',
      '3': 452625788,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'inputdetails'
    },
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {
      '1': 'taskcredentials',
      '3': 257843259,
      '4': 1,
      '5': 11,
      '6': '.states.TaskCredentials',
      '10': 'taskcredentials'
    },
    {
      '1': 'timeoutinseconds',
      '3': 472710197,
      '4': 1,
      '5': 3,
      '10': 'timeoutinseconds'
    },
  ],
};

/// Descriptor for `LambdaFunctionScheduledEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionScheduledEventDetailsDescriptor = $convert.base64Decode(
    'CiNMYW1iZGFGdW5jdGlvblNjaGVkdWxlZEV2ZW50RGV0YWlscxIYCgVpbnB1dBj83uHOASABKA'
    'lSBWlucHV0ElAKDGlucHV0ZGV0YWlscxj8iurXASABKAsyKC5zdGF0ZXMuSGlzdG9yeUV2ZW50'
    'RXhlY3V0aW9uRGF0YURldGFpbHNSDGlucHV0ZGV0YWlscxIdCghyZXNvdXJjZRj2//1OIAEoCV'
    'IIcmVzb3VyY2USRAoPdGFza2NyZWRlbnRpYWxzGLvA+XogASgLMhcuc3RhdGVzLlRhc2tDcmVk'
    'ZW50aWFsc1IPdGFza2NyZWRlbnRpYWxzEi4KEHRpbWVvdXRpbnNlY29uZHMYtfiz4QEgASgDUh'
    'B0aW1lb3V0aW5zZWNvbmRz');

@$core.Deprecated('Use lambdaFunctionStartFailedEventDetailsDescriptor instead')
const LambdaFunctionStartFailedEventDetails$json = {
  '1': 'LambdaFunctionStartFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `LambdaFunctionStartFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionStartFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'CiVMYW1iZGFGdW5jdGlvblN0YXJ0RmFpbGVkRXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0UgAS'
        'gJUgVjYXVzZRIXCgVlcnJvchjSjsYMIAEoCVIFZXJyb3I=');

@$core.Deprecated('Use lambdaFunctionSucceededEventDetailsDescriptor instead')
const LambdaFunctionSucceededEventDetails$json = {
  '1': 'LambdaFunctionSucceededEventDetails',
  '2': [
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputdetails',
      '3': 393734643,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'outputdetails'
    },
  ],
};

/// Descriptor for `LambdaFunctionSucceededEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionSucceededEventDetailsDescriptor =
    $convert.base64Decode(
        'CiNMYW1iZGFGdW5jdGlvblN1Y2NlZWRlZEV2ZW50RGV0YWlscxIaCgZvdXRwdXQYhZ6lzQEgAS'
        'gJUgZvdXRwdXQSUgoNb3V0cHV0ZGV0YWlscxjz09+7ASABKAsyKC5zdGF0ZXMuSGlzdG9yeUV2'
        'ZW50RXhlY3V0aW9uRGF0YURldGFpbHNSDW91dHB1dGRldGFpbHM=');

@$core.Deprecated('Use lambdaFunctionTimedOutEventDetailsDescriptor instead')
const LambdaFunctionTimedOutEventDetails$json = {
  '1': 'LambdaFunctionTimedOutEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `LambdaFunctionTimedOutEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lambdaFunctionTimedOutEventDetailsDescriptor =
    $convert.base64Decode(
        'CiJMYW1iZGFGdW5jdGlvblRpbWVkT3V0RXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0UgASgJUg'
        'VjYXVzZRIXCgVlcnJvchjSjsYMIAEoCVIFZXJyb3I=');

@$core.Deprecated('Use listActivitiesInputDescriptor instead')
const ListActivitiesInput$json = {
  '1': 'ListActivitiesInput',
  '2': [
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListActivitiesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listActivitiesInputDescriptor = $convert.base64Decode(
    'ChNMaXN0QWN0aXZpdGllc0lucHV0EiIKCm1heHJlc3VsdHMY0tzn3QEgASgFUgptYXhyZXN1bH'
    'RzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listActivitiesOutputDescriptor instead')
const ListActivitiesOutput$json = {
  '1': 'ListActivitiesOutput',
  '2': [
    {
      '1': 'activities',
      '3': 164950681,
      '4': 3,
      '5': 11,
      '6': '.states.ActivityListItem',
      '10': 'activities'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListActivitiesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listActivitiesOutputDescriptor = $convert.base64Decode(
    'ChRMaXN0QWN0aXZpdGllc091dHB1dBI7CgphY3Rpdml0aWVzGJnl004gAygLMhguc3RhdGVzLk'
    'FjdGl2aXR5TGlzdEl0ZW1SCmFjdGl2aXRpZXMSHwoJbmV4dHRva2VuGJ7znTcgASgJUgluZXh0'
    'dG9rZW4=');

@$core.Deprecated('Use listExecutionsInputDescriptor instead')
const ListExecutionsInput$json = {
  '1': 'ListExecutionsInput',
  '2': [
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'redrivefilter',
      '3': 437804251,
      '4': 1,
      '5': 14,
      '6': '.states.ExecutionRedriveFilter',
      '10': 'redrivefilter'
    },
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'statusfilter',
      '3': 86045418,
      '4': 1,
      '5': 14,
      '6': '.states.ExecutionStatus',
      '10': 'statusfilter'
    },
  ],
};

/// Descriptor for `ListExecutionsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listExecutionsInputDescriptor = $convert.base64Decode(
    'ChNMaXN0RXhlY3V0aW9uc0lucHV0Eh8KCW1hcHJ1bmFybhi669YIIAEoCVIJbWFwcnVuYXJuEi'
    'IKCm1heHJlc3VsdHMY0tzn3QEgASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhie8503IAEo'
    'CVIJbmV4dHRva2VuEkgKDXJlZHJpdmVmaWx0ZXIY27nh0AEgASgOMh4uc3RhdGVzLkV4ZWN1dG'
    'lvblJlZHJpdmVGaWx0ZXJSDXJlZHJpdmVmaWx0ZXISLAoPc3RhdGVtYWNoaW5lYXJuGPO7xrsB'
    'IAEoCVIPc3RhdGVtYWNoaW5lYXJuEj4KDHN0YXR1c2ZpbHRlchjq5YMpIAEoDjIXLnN0YXRlcy'
    '5FeGVjdXRpb25TdGF0dXNSDHN0YXR1c2ZpbHRlcg==');

@$core.Deprecated('Use listExecutionsOutputDescriptor instead')
const ListExecutionsOutput$json = {
  '1': 'ListExecutionsOutput',
  '2': [
    {
      '1': 'executions',
      '3': 113395923,
      '4': 3,
      '5': 11,
      '6': '.states.ExecutionListItem',
      '10': 'executions'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListExecutionsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listExecutionsOutputDescriptor = $convert.base64Decode(
    'ChRMaXN0RXhlY3V0aW9uc091dHB1dBI8CgpleGVjdXRpb25zGNORiTYgAygLMhkuc3RhdGVzLk'
    'V4ZWN1dGlvbkxpc3RJdGVtUgpleGVjdXRpb25zEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4'
    'dHRva2Vu');

@$core.Deprecated('Use listMapRunsInputDescriptor instead')
const ListMapRunsInput$json = {
  '1': 'ListMapRunsInput',
  '2': [
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListMapRunsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMapRunsInputDescriptor = $convert.base64Decode(
    'ChBMaXN0TWFwUnVuc0lucHV0EiYKDGV4ZWN1dGlvbmFybhjtlv2VASABKAlSDGV4ZWN1dGlvbm'
    'FybhIiCgptYXhyZXN1bHRzGNLc590BIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4YnvOd'
    'NyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listMapRunsOutputDescriptor instead')
const ListMapRunsOutput$json = {
  '1': 'ListMapRunsOutput',
  '2': [
    {
      '1': 'mapruns',
      '3': 34690200,
      '4': 3,
      '5': 11,
      '6': '.states.MapRunListItem',
      '10': 'mapruns'
    },
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListMapRunsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMapRunsOutputDescriptor = $convert.base64Decode(
    'ChFMaXN0TWFwUnVuc091dHB1dBIzCgdtYXBydW5zGJipxRAgAygLMhYuc3RhdGVzLk1hcFJ1bk'
    'xpc3RJdGVtUgdtYXBydW5zEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listStateMachineAliasesInputDescriptor instead')
const ListStateMachineAliasesInput$json = {
  '1': 'ListStateMachineAliasesInput',
  '2': [
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
  ],
};

/// Descriptor for `ListStateMachineAliasesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStateMachineAliasesInputDescriptor =
    $convert.base64Decode(
        'ChxMaXN0U3RhdGVNYWNoaW5lQWxpYXNlc0lucHV0EiIKCm1heHJlc3VsdHMY0tzn3QEgASgFUg'
        'ptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2VuEiwKD3N0YXRlbWFj'
        'aGluZWFybhjzu8a7ASABKAlSD3N0YXRlbWFjaGluZWFybg==');

@$core.Deprecated('Use listStateMachineAliasesOutputDescriptor instead')
const ListStateMachineAliasesOutput$json = {
  '1': 'ListStateMachineAliasesOutput',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'statemachinealiases',
      '3': 452532502,
      '4': 3,
      '5': 11,
      '6': '.states.StateMachineAliasListItem',
      '10': 'statemachinealiases'
    },
  ],
};

/// Descriptor for `ListStateMachineAliasesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStateMachineAliasesOutputDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0U3RhdGVNYWNoaW5lQWxpYXNlc091dHB1dBIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW'
        '5leHR0b2tlbhJXChNzdGF0ZW1hY2hpbmVhbGlhc2VzGJay5NcBIAMoCzIhLnN0YXRlcy5TdGF0'
        'ZU1hY2hpbmVBbGlhc0xpc3RJdGVtUhNzdGF0ZW1hY2hpbmVhbGlhc2Vz');

@$core.Deprecated('Use listStateMachineVersionsInputDescriptor instead')
const ListStateMachineVersionsInput$json = {
  '1': 'ListStateMachineVersionsInput',
  '2': [
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
  ],
};

/// Descriptor for `ListStateMachineVersionsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStateMachineVersionsInputDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0U3RhdGVNYWNoaW5lVmVyc2lvbnNJbnB1dBIiCgptYXhyZXN1bHRzGNLc590BIAEoBV'
        'IKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2tlbhIsCg9zdGF0ZW1h'
        'Y2hpbmVhcm4Y87vGuwEgASgJUg9zdGF0ZW1hY2hpbmVhcm4=');

@$core.Deprecated('Use listStateMachineVersionsOutputDescriptor instead')
const ListStateMachineVersionsOutput$json = {
  '1': 'ListStateMachineVersionsOutput',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'statemachineversions',
      '3': 133258663,
      '4': 3,
      '5': 11,
      '6': '.states.StateMachineVersionListItem',
      '10': 'statemachineversions'
    },
  ],
};

/// Descriptor for `ListStateMachineVersionsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStateMachineVersionsOutputDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0U3RhdGVNYWNoaW5lVmVyc2lvbnNPdXRwdXQSHwoJbmV4dHRva2VuGJ7znTcgASgJUg'
        'luZXh0dG9rZW4SWgoUc3RhdGVtYWNoaW5ldmVyc2lvbnMYp7vFPyADKAsyIy5zdGF0ZXMuU3Rh'
        'dGVNYWNoaW5lVmVyc2lvbkxpc3RJdGVtUhRzdGF0ZW1hY2hpbmV2ZXJzaW9ucw==');

@$core.Deprecated('Use listStateMachinesInputDescriptor instead')
const ListStateMachinesInput$json = {
  '1': 'ListStateMachinesInput',
  '2': [
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListStateMachinesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStateMachinesInputDescriptor =
    $convert.base64Decode(
        'ChZMaXN0U3RhdGVNYWNoaW5lc0lucHV0EiIKCm1heHJlc3VsdHMY0tzn3QEgASgFUgptYXhyZX'
        'N1bHRzEh8KCW5leHR0b2tlbhie8503IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listStateMachinesOutputDescriptor instead')
const ListStateMachinesOutput$json = {
  '1': 'ListStateMachinesOutput',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'statemachines',
      '3': 432113525,
      '4': 3,
      '5': 11,
      '6': '.states.StateMachineListItem',
      '10': 'statemachines'
    },
  ],
};

/// Descriptor for `ListStateMachinesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStateMachinesOutputDescriptor = $convert.base64Decode(
    'ChdMaXN0U3RhdGVNYWNoaW5lc091dHB1dBIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW5leHR0b2'
    'tlbhJGCg1zdGF0ZW1hY2hpbmVzGPWOhs4BIAMoCzIcLnN0YXRlcy5TdGF0ZU1hY2hpbmVMaXN0'
    'SXRlbVINc3RhdGVtYWNoaW5lcw==');

@$core.Deprecated('Use listTagsForResourceInputDescriptor instead')
const ListTagsForResourceInput$json = {
  '1': 'ListTagsForResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceInputDescriptor =
    $convert.base64Decode(
        'ChhMaXN0VGFnc0ZvclJlc291cmNlSW5wdXQSIwoLcmVzb3VyY2Vhcm4YzcyqICABKAlSC3Jlc2'
        '91cmNlYXJu');

@$core.Deprecated('Use listTagsForResourceOutputDescriptor instead')
const ListTagsForResourceOutput$json = {
  '1': 'ListTagsForResourceOutput',
  '2': [
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.states.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceOutputDescriptor =
    $convert.base64Decode(
        'ChlMaXN0VGFnc0ZvclJlc291cmNlT3V0cHV0EiMKBHRhZ3MYodfboAEgAygLMgsuc3RhdGVzLl'
        'RhZ1IEdGFncw==');

@$core.Deprecated('Use logDestinationDescriptor instead')
const LogDestination$json = {
  '1': 'LogDestination',
  '2': [
    {
      '1': 'cloudwatchlogsloggroup',
      '3': 517373608,
      '4': 1,
      '5': 11,
      '6': '.states.CloudWatchLogsLogGroup',
      '10': 'cloudwatchlogsloggroup'
    },
  ],
};

/// Descriptor for `LogDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logDestinationDescriptor = $convert.base64Decode(
    'Cg5Mb2dEZXN0aW5hdGlvbhJaChZjbG91ZHdhdGNobG9nc2xvZ2dyb3VwGKj92fYBIAEoCzIeLn'
    'N0YXRlcy5DbG91ZFdhdGNoTG9nc0xvZ0dyb3VwUhZjbG91ZHdhdGNobG9nc2xvZ2dyb3Vw');

@$core.Deprecated('Use loggingConfigurationDescriptor instead')
const LoggingConfiguration$json = {
  '1': 'LoggingConfiguration',
  '2': [
    {
      '1': 'destinations',
      '3': 1617189,
      '4': 3,
      '5': 11,
      '6': '.states.LogDestination',
      '10': 'destinations'
    },
    {
      '1': 'includeexecutiondata',
      '3': 203899608,
      '4': 1,
      '5': 8,
      '10': 'includeexecutiondata'
    },
    {
      '1': 'level',
      '3': 463071198,
      '4': 1,
      '5': 14,
      '6': '.states.LogLevel',
      '10': 'level'
    },
  ],
};

/// Descriptor for `LoggingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loggingConfigurationDescriptor = $convert.base64Decode(
    'ChRMb2dnaW5nQ29uZmlndXJhdGlvbhI8CgxkZXN0aW5hdGlvbnMYpdpiIAMoCzIWLnN0YXRlcy'
    '5Mb2dEZXN0aW5hdGlvblIMZGVzdGluYXRpb25zEjUKFGluY2x1ZGVleGVjdXRpb25kYXRhGNiF'
    'nWEgASgIUhRpbmNsdWRlZXhlY3V0aW9uZGF0YRIqCgVsZXZlbBjez+fcASABKA4yEC5zdGF0ZX'
    'MuTG9nTGV2ZWxSBWxldmVs');

@$core.Deprecated('Use mapIterationEventDetailsDescriptor instead')
const MapIterationEventDetails$json = {
  '1': 'MapIterationEventDetails',
  '2': [
    {'1': 'index', '3': 151693740, '4': 1, '5': 5, '10': 'index'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `MapIterationEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mapIterationEventDetailsDescriptor =
    $convert.base64Decode(
        'ChhNYXBJdGVyYXRpb25FdmVudERldGFpbHMSFwoFaW5kZXgYrNOqSCABKAVSBWluZGV4EhUKBG'
        '5hbWUY5/vmaSABKAlSBG5hbWU=');

@$core.Deprecated('Use mapRunExecutionCountsDescriptor instead')
const MapRunExecutionCounts$json = {
  '1': 'MapRunExecutionCounts',
  '2': [
    {'1': 'aborted', '3': 381485777, '4': 1, '5': 3, '10': 'aborted'},
    {'1': 'failed', '3': 11325365, '4': 1, '5': 3, '10': 'failed'},
    {
      '1': 'failuresnotredrivable',
      '3': 119492388,
      '4': 1,
      '5': 3,
      '10': 'failuresnotredrivable'
    },
    {'1': 'pending', '3': 431980349, '4': 1, '5': 3, '10': 'pending'},
    {
      '1': 'pendingredrive',
      '3': 405645200,
      '4': 1,
      '5': 3,
      '10': 'pendingredrive'
    },
    {
      '1': 'resultswritten',
      '3': 88286449,
      '4': 1,
      '5': 3,
      '10': 'resultswritten'
    },
    {'1': 'running', '3': 343848781, '4': 1, '5': 3, '10': 'running'},
    {'1': 'succeeded', '3': 186019827, '4': 1, '5': 3, '10': 'succeeded'},
    {'1': 'timedout', '3': 335029733, '4': 1, '5': 3, '10': 'timedout'},
    {'1': 'total', '3': 80777982, '4': 1, '5': 3, '10': 'total'},
  ],
};

/// Descriptor for `MapRunExecutionCounts`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mapRunExecutionCountsDescriptor = $convert.base64Decode(
    'ChVNYXBSdW5FeGVjdXRpb25Db3VudHMSHAoHYWJvcnRlZBjRhfS1ASABKANSB2Fib3J0ZWQSGQ'
    'oGZmFpbGVkGLWfswUgASgDUgZmYWlsZWQSNwoVZmFpbHVyZXNub3RyZWRyaXZhYmxlGKSe/Tgg'
    'ASgDUhVmYWlsdXJlc25vdHJlZHJpdmFibGUSHAoHcGVuZGluZxi9/v3NASABKANSB3BlbmRpbm'
    'cSKgoOcGVuZGluZ3JlZHJpdmUYkM+2wQEgASgDUg5wZW5kaW5ncmVkcml2ZRIpCg5yZXN1bHRz'
    'd3JpdHRlbhjxyYwqIAEoA1IOcmVzdWx0c3dyaXR0ZW4SHAoHcnVubmluZxjN7vqjASABKANSB3'
    'J1bm5pbmcSHwoJc3VjY2VlZGVkGPPf2VggASgDUglzdWNjZWVkZWQSHgoIdGltZWRvdXQY5cvg'
    'nwEgASgDUgh0aW1lZG91dBIXCgV0b3RhbBj+pcImIAEoA1IFdG90YWw=');

@$core.Deprecated('Use mapRunFailedEventDetailsDescriptor instead')
const MapRunFailedEventDetails$json = {
  '1': 'MapRunFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `MapRunFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mapRunFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChhNYXBSdW5GYWlsZWRFdmVudERldGFpbHMSFwoFY2F1c2UYoaS7RSABKAlSBWNhdXNlEhcKBW'
        'Vycm9yGNKOxgwgASgJUgVlcnJvcg==');

@$core.Deprecated('Use mapRunItemCountsDescriptor instead')
const MapRunItemCounts$json = {
  '1': 'MapRunItemCounts',
  '2': [
    {'1': 'aborted', '3': 381485777, '4': 1, '5': 3, '10': 'aborted'},
    {'1': 'failed', '3': 11325365, '4': 1, '5': 3, '10': 'failed'},
    {
      '1': 'failuresnotredrivable',
      '3': 119492388,
      '4': 1,
      '5': 3,
      '10': 'failuresnotredrivable'
    },
    {'1': 'pending', '3': 431980349, '4': 1, '5': 3, '10': 'pending'},
    {
      '1': 'pendingredrive',
      '3': 405645200,
      '4': 1,
      '5': 3,
      '10': 'pendingredrive'
    },
    {
      '1': 'resultswritten',
      '3': 88286449,
      '4': 1,
      '5': 3,
      '10': 'resultswritten'
    },
    {'1': 'running', '3': 343848781, '4': 1, '5': 3, '10': 'running'},
    {'1': 'succeeded', '3': 186019827, '4': 1, '5': 3, '10': 'succeeded'},
    {'1': 'timedout', '3': 335029733, '4': 1, '5': 3, '10': 'timedout'},
    {'1': 'total', '3': 80777982, '4': 1, '5': 3, '10': 'total'},
  ],
};

/// Descriptor for `MapRunItemCounts`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mapRunItemCountsDescriptor = $convert.base64Decode(
    'ChBNYXBSdW5JdGVtQ291bnRzEhwKB2Fib3J0ZWQY0YX0tQEgASgDUgdhYm9ydGVkEhkKBmZhaW'
    'xlZBi1n7MFIAEoA1IGZmFpbGVkEjcKFWZhaWx1cmVzbm90cmVkcml2YWJsZRiknv04IAEoA1IV'
    'ZmFpbHVyZXNub3RyZWRyaXZhYmxlEhwKB3BlbmRpbmcYvf79zQEgASgDUgdwZW5kaW5nEioKDn'
    'BlbmRpbmdyZWRyaXZlGJDPtsEBIAEoA1IOcGVuZGluZ3JlZHJpdmUSKQoOcmVzdWx0c3dyaXR0'
    'ZW4Y8cmMKiABKANSDnJlc3VsdHN3cml0dGVuEhwKB3J1bm5pbmcYze76owEgASgDUgdydW5uaW'
    '5nEh8KCXN1Y2NlZWRlZBjz39lYIAEoA1IJc3VjY2VlZGVkEh4KCHRpbWVkb3V0GOXL4J8BIAEo'
    'A1IIdGltZWRvdXQSFwoFdG90YWwY/qXCJiABKANSBXRvdGFs');

@$core.Deprecated('Use mapRunListItemDescriptor instead')
const MapRunListItem$json = {
  '1': 'MapRunListItem',
  '2': [
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
    {'1': 'startdate', '3': 364840732, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {'1': 'stopdate', '3': 180697434, '4': 1, '5': 9, '10': 'stopdate'},
  ],
};

/// Descriptor for `MapRunListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mapRunListItemDescriptor = $convert.base64Decode(
    'Cg5NYXBSdW5MaXN0SXRlbRImCgxleGVjdXRpb25hcm4Y7Zb9lQEgASgJUgxleGVjdXRpb25hcm'
    '4SHwoJbWFwcnVuYXJuGLrr1gggASgJUgltYXBydW5hcm4SIAoJc3RhcnRkYXRlGJyO/K0BIAEo'
    'CVIJc3RhcnRkYXRlEiwKD3N0YXRlbWFjaGluZWFybhjzu8a7ASABKAlSD3N0YXRlbWFjaGluZW'
    'FybhIdCghzdG9wZGF0ZRja8pRWIAEoCVIIc3RvcGRhdGU=');

@$core.Deprecated('Use mapRunRedrivenEventDetailsDescriptor instead')
const MapRunRedrivenEventDetails$json = {
  '1': 'MapRunRedrivenEventDetails',
  '2': [
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
    {'1': 'redrivecount', '3': 473458696, '4': 1, '5': 5, '10': 'redrivecount'},
  ],
};

/// Descriptor for `MapRunRedrivenEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mapRunRedrivenEventDetailsDescriptor =
    $convert.base64Decode(
        'ChpNYXBSdW5SZWRyaXZlbkV2ZW50RGV0YWlscxIfCgltYXBydW5hcm4YuuvWCCABKAlSCW1hcH'
        'J1bmFybhImCgxyZWRyaXZlY291bnQYiNDh4QEgASgFUgxyZWRyaXZlY291bnQ=');

@$core.Deprecated('Use mapRunStartedEventDetailsDescriptor instead')
const MapRunStartedEventDetails$json = {
  '1': 'MapRunStartedEventDetails',
  '2': [
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
  ],
};

/// Descriptor for `MapRunStartedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mapRunStartedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChlNYXBSdW5TdGFydGVkRXZlbnREZXRhaWxzEh8KCW1hcHJ1bmFybhi669YIIAEoCVIJbWFwcn'
        'VuYXJu');

@$core.Deprecated('Use mapStateStartedEventDetailsDescriptor instead')
const MapStateStartedEventDetails$json = {
  '1': 'MapStateStartedEventDetails',
  '2': [
    {'1': 'length', '3': 63976982, '4': 1, '5': 5, '10': 'length'},
  ],
};

/// Descriptor for `MapStateStartedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mapStateStartedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChtNYXBTdGF0ZVN0YXJ0ZWRFdmVudERldGFpbHMSGQoGbGVuZ3RoGJbswB4gASgFUgZsZW5ndG'
        'g=');

@$core.Deprecated('Use missingRequiredParameterDescriptor instead')
const MissingRequiredParameter$json = {
  '1': 'MissingRequiredParameter',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MissingRequiredParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List missingRequiredParameterDescriptor =
    $convert.base64Decode(
        'ChhNaXNzaW5nUmVxdWlyZWRQYXJhbWV0ZXISGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use mockErrorOutputDescriptor instead')
const MockErrorOutput$json = {
  '1': 'MockErrorOutput',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
  ],
};

/// Descriptor for `MockErrorOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mockErrorOutputDescriptor = $convert.base64Decode(
    'Cg9Nb2NrRXJyb3JPdXRwdXQSFwoFY2F1c2UYoaS7RSABKAlSBWNhdXNlEhcKBWVycm9yGNKOxg'
    'wgASgJUgVlcnJvcg==');

@$core.Deprecated('Use mockInputDescriptor instead')
const MockInput$json = {
  '1': 'MockInput',
  '2': [
    {
      '1': 'erroroutput',
      '3': 453597273,
      '4': 1,
      '5': 11,
      '6': '.states.MockErrorOutput',
      '10': 'erroroutput'
    },
    {
      '1': 'fieldvalidationmode',
      '3': 519416556,
      '4': 1,
      '5': 14,
      '6': '.states.MockResponseValidationMode',
      '10': 'fieldvalidationmode'
    },
    {'1': 'result', '3': 171406885, '4': 1, '5': 9, '10': 'result'},
  ],
};

/// Descriptor for `MockInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mockInputDescriptor = $convert.base64Decode(
    'CglNb2NrSW5wdXQSPQoLZXJyb3JvdXRwdXQY2bCl2AEgASgLMhcuc3RhdGVzLk1vY2tFcnJvck'
    '91dHB1dFILZXJyb3JvdXRwdXQSWAoTZmllbGR2YWxpZGF0aW9ubW9kZRjs1db3ASABKA4yIi5z'
    'dGF0ZXMuTW9ja1Jlc3BvbnNlVmFsaWRhdGlvbk1vZGVSE2ZpZWxkdmFsaWRhdGlvbm1vZGUSGQ'
    'oGcmVzdWx0GKXs3VEgASgJUgZyZXN1bHQ=');

@$core.Deprecated('Use publishStateMachineVersionInputDescriptor instead')
const PublishStateMachineVersionInput$json = {
  '1': 'PublishStateMachineVersionInput',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {'1': 'revisionid', '3': 369170086, '4': 1, '5': 9, '10': 'revisionid'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
  ],
};

/// Descriptor for `PublishStateMachineVersionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishStateMachineVersionInputDescriptor =
    $convert.base64Decode(
        'Ch9QdWJsaXNoU3RhdGVNYWNoaW5lVmVyc2lvbklucHV0EiQKC2Rlc2NyaXB0aW9uGOr2vKMBIA'
        'EoCVILZGVzY3JpcHRpb24SIgoKcmV2aXNpb25pZBimrYSwASABKAlSCnJldmlzaW9uaWQSLAoP'
        'c3RhdGVtYWNoaW5lYXJuGPO7xrsBIAEoCVIPc3RhdGVtYWNoaW5lYXJu');

@$core.Deprecated('Use publishStateMachineVersionOutputDescriptor instead')
const PublishStateMachineVersionOutput$json = {
  '1': 'PublishStateMachineVersionOutput',
  '2': [
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
  ],
};

/// Descriptor for `PublishStateMachineVersionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishStateMachineVersionOutputDescriptor =
    $convert.base64Decode(
        'CiBQdWJsaXNoU3RhdGVNYWNoaW5lVmVyc2lvbk91dHB1dBIlCgxjcmVhdGlvbmRhdGUYgc7RcS'
        'ABKAlSDGNyZWF0aW9uZGF0ZRI5ChZzdGF0ZW1hY2hpbmV2ZXJzaW9uYXJuGPmFryEgASgJUhZz'
        'dGF0ZW1hY2hpbmV2ZXJzaW9uYXJu');

@$core.Deprecated('Use redriveExecutionInputDescriptor instead')
const RedriveExecutionInput$json = {
  '1': 'RedriveExecutionInput',
  '2': [
    {'1': 'clienttoken', '3': 272531820, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
  ],
};

/// Descriptor for `RedriveExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redriveExecutionInputDescriptor = $convert.base64Decode(
    'ChVSZWRyaXZlRXhlY3V0aW9uSW5wdXQSJAoLY2xpZW50dG9rZW4Y7IL6gQEgASgJUgtjbGllbn'
    'R0b2tlbhImCgxleGVjdXRpb25hcm4Y7Zb9lQEgASgJUgxleGVjdXRpb25hcm4=');

@$core.Deprecated('Use redriveExecutionOutputDescriptor instead')
const RedriveExecutionOutput$json = {
  '1': 'RedriveExecutionOutput',
  '2': [
    {'1': 'redrivedate', '3': 152812125, '4': 1, '5': 9, '10': 'redrivedate'},
  ],
};

/// Descriptor for `RedriveExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List redriveExecutionOutputDescriptor =
    $convert.base64Decode(
        'ChZSZWRyaXZlRXhlY3V0aW9uT3V0cHV0EiMKC3JlZHJpdmVkYXRlGN307kggASgJUgtyZWRyaX'
        'ZlZGF0ZQ==');

@$core.Deprecated('Use resourceNotFoundDescriptor instead')
const ResourceNotFound$json = {
  '1': 'ResourceNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'resourcename', '3': 17776375, '4': 1, '5': 9, '10': 'resourcename'},
  ],
};

/// Descriptor for `ResourceNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundDescriptor = $convert.base64Decode(
    'ChBSZXNvdXJjZU5vdEZvdW5kEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2USJQoMcmVzb3'
    'VyY2VuYW1lGPf9vAggASgJUgxyZXNvdXJjZW5hbWU=');

@$core.Deprecated('Use routingConfigurationListItemDescriptor instead')
const RoutingConfigurationListItem$json = {
  '1': 'RoutingConfigurationListItem',
  '2': [
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
    {'1': 'weight', '3': 278961850, '4': 1, '5': 5, '10': 'weight'},
  ],
};

/// Descriptor for `RoutingConfigurationListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List routingConfigurationListItemDescriptor =
    $convert.base64Decode(
        'ChxSb3V0aW5nQ29uZmlndXJhdGlvbkxpc3RJdGVtEjkKFnN0YXRlbWFjaGluZXZlcnNpb25hcm'
        '4Y+YWvISABKAlSFnN0YXRlbWFjaGluZXZlcnNpb25hcm4SGgoGd2VpZ2h0GLq9goUBIAEoBVIG'
        'd2VpZ2h0');

@$core.Deprecated('Use sendTaskFailureInputDescriptor instead')
const SendTaskFailureInput$json = {
  '1': 'SendTaskFailureInput',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'tasktoken', '3': 525325834, '4': 1, '5': 9, '10': 'tasktoken'},
  ],
};

/// Descriptor for `SendTaskFailureInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTaskFailureInputDescriptor = $convert.base64Decode(
    'ChRTZW5kVGFza0ZhaWx1cmVJbnB1dBIXCgVjYXVzZRihpLtFIAEoCVIFY2F1c2USFwoFZXJyb3'
    'IY0o7GDCABKAlSBWVycm9yEiAKCXRhc2t0b2tlbhiKrL/6ASABKAlSCXRhc2t0b2tlbg==');

@$core.Deprecated('Use sendTaskFailureOutputDescriptor instead')
const SendTaskFailureOutput$json = {
  '1': 'SendTaskFailureOutput',
};

/// Descriptor for `SendTaskFailureOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTaskFailureOutputDescriptor =
    $convert.base64Decode('ChVTZW5kVGFza0ZhaWx1cmVPdXRwdXQ=');

@$core.Deprecated('Use sendTaskHeartbeatInputDescriptor instead')
const SendTaskHeartbeatInput$json = {
  '1': 'SendTaskHeartbeatInput',
  '2': [
    {'1': 'tasktoken', '3': 525325834, '4': 1, '5': 9, '10': 'tasktoken'},
  ],
};

/// Descriptor for `SendTaskHeartbeatInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTaskHeartbeatInputDescriptor =
    $convert.base64Decode(
        'ChZTZW5kVGFza0hlYXJ0YmVhdElucHV0EiAKCXRhc2t0b2tlbhiKrL/6ASABKAlSCXRhc2t0b2'
        'tlbg==');

@$core.Deprecated('Use sendTaskHeartbeatOutputDescriptor instead')
const SendTaskHeartbeatOutput$json = {
  '1': 'SendTaskHeartbeatOutput',
};

/// Descriptor for `SendTaskHeartbeatOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTaskHeartbeatOutputDescriptor =
    $convert.base64Decode('ChdTZW5kVGFza0hlYXJ0YmVhdE91dHB1dA==');

@$core.Deprecated('Use sendTaskSuccessInputDescriptor instead')
const SendTaskSuccessInput$json = {
  '1': 'SendTaskSuccessInput',
  '2': [
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {'1': 'tasktoken', '3': 525325834, '4': 1, '5': 9, '10': 'tasktoken'},
  ],
};

/// Descriptor for `SendTaskSuccessInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTaskSuccessInputDescriptor = $convert.base64Decode(
    'ChRTZW5kVGFza1N1Y2Nlc3NJbnB1dBIaCgZvdXRwdXQYhZ6lzQEgASgJUgZvdXRwdXQSIAoJdG'
    'Fza3Rva2VuGIqsv/oBIAEoCVIJdGFza3Rva2Vu');

@$core.Deprecated('Use sendTaskSuccessOutputDescriptor instead')
const SendTaskSuccessOutput$json = {
  '1': 'SendTaskSuccessOutput',
};

/// Descriptor for `SendTaskSuccessOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendTaskSuccessOutputDescriptor =
    $convert.base64Decode('ChVTZW5kVGFza1N1Y2Nlc3NPdXRwdXQ=');

@$core.Deprecated('Use serviceQuotaExceededExceptionDescriptor instead')
const ServiceQuotaExceededException$json = {
  '1': 'ServiceQuotaExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceQuotaExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceQuotaExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch1TZXJ2aWNlUXVvdGFFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use startExecutionInputDescriptor instead')
const StartExecutionInput$json = {
  '1': 'StartExecutionInput',
  '2': [
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {'1': 'traceheader', '3': 219960864, '4': 1, '5': 9, '10': 'traceheader'},
  ],
};

/// Descriptor for `StartExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startExecutionInputDescriptor = $convert.base64Decode(
    'ChNTdGFydEV4ZWN1dGlvbklucHV0EhgKBWlucHV0GPze4c4BIAEoCVIFaW5wdXQSFQoEbmFtZR'
    'jn++ZpIAEoCVIEbmFtZRIsCg9zdGF0ZW1hY2hpbmVhcm4Y87vGuwEgASgJUg9zdGF0ZW1hY2hp'
    'bmVhcm4SIwoLdHJhY2VoZWFkZXIYoKzxaCABKAlSC3RyYWNlaGVhZGVy');

@$core.Deprecated('Use startExecutionOutputDescriptor instead')
const StartExecutionOutput$json = {
  '1': 'StartExecutionOutput',
  '2': [
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {'1': 'startdate', '3': 364840732, '4': 1, '5': 9, '10': 'startdate'},
  ],
};

/// Descriptor for `StartExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startExecutionOutputDescriptor = $convert.base64Decode(
    'ChRTdGFydEV4ZWN1dGlvbk91dHB1dBImCgxleGVjdXRpb25hcm4Y7Zb9lQEgASgJUgxleGVjdX'
    'Rpb25hcm4SIAoJc3RhcnRkYXRlGJyO/K0BIAEoCVIJc3RhcnRkYXRl');

@$core.Deprecated('Use startSyncExecutionInputDescriptor instead')
const StartSyncExecutionInput$json = {
  '1': 'StartSyncExecutionInput',
  '2': [
    {
      '1': 'includeddata',
      '3': 109719114,
      '4': 1,
      '5': 14,
      '6': '.states.IncludedData',
      '10': 'includeddata'
    },
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {'1': 'traceheader', '3': 219960864, '4': 1, '5': 9, '10': 'traceheader'},
  ],
};

/// Descriptor for `StartSyncExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startSyncExecutionInputDescriptor = $convert.base64Decode(
    'ChdTdGFydFN5bmNFeGVjdXRpb25JbnB1dBI7CgxpbmNsdWRlZGRhdGEYytyoNCABKA4yFC5zdG'
    'F0ZXMuSW5jbHVkZWREYXRhUgxpbmNsdWRlZGRhdGESGAoFaW5wdXQY/N7hzgEgASgJUgVpbnB1'
    'dBIVCgRuYW1lGOf75mkgASgJUgRuYW1lEiwKD3N0YXRlbWFjaGluZWFybhjzu8a7ASABKAlSD3'
    'N0YXRlbWFjaGluZWFybhIjCgt0cmFjZWhlYWRlchigrPFoIAEoCVILdHJhY2VoZWFkZXI=');

@$core.Deprecated('Use startSyncExecutionOutputDescriptor instead')
const StartSyncExecutionOutput$json = {
  '1': 'StartSyncExecutionOutput',
  '2': [
    {
      '1': 'billingdetails',
      '3': 270001723,
      '4': 1,
      '5': 11,
      '6': '.states.BillingDetails',
      '10': 'billingdetails'
    },
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'inputdetails',
      '3': 452625788,
      '4': 1,
      '5': 11,
      '6': '.states.CloudWatchEventsExecutionDataDetails',
      '10': 'inputdetails'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputdetails',
      '3': 393734643,
      '4': 1,
      '5': 11,
      '6': '.states.CloudWatchEventsExecutionDataDetails',
      '10': 'outputdetails'
    },
    {'1': 'startdate', '3': 364840732, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.states.SyncExecutionStatus',
      '10': 'status'
    },
    {'1': 'stopdate', '3': 180697434, '4': 1, '5': 9, '10': 'stopdate'},
    {'1': 'traceheader', '3': 219960864, '4': 1, '5': 9, '10': 'traceheader'},
  ],
};

/// Descriptor for `StartSyncExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startSyncExecutionOutputDescriptor = $convert.base64Decode(
    'ChhTdGFydFN5bmNFeGVjdXRpb25PdXRwdXQSQgoOYmlsbGluZ2RldGFpbHMYu8zfgAEgASgLMh'
    'Yuc3RhdGVzLkJpbGxpbmdEZXRhaWxzUg5iaWxsaW5nZGV0YWlscxIXCgVjYXVzZRihpLtFIAEo'
    'CVIFY2F1c2USFwoFZXJyb3IY0o7GDCABKAlSBWVycm9yEiYKDGV4ZWN1dGlvbmFybhjtlv2VAS'
    'ABKAlSDGV4ZWN1dGlvbmFybhIYCgVpbnB1dBj83uHOASABKAlSBWlucHV0ElQKDGlucHV0ZGV0'
    'YWlscxj8iurXASABKAsyLC5zdGF0ZXMuQ2xvdWRXYXRjaEV2ZW50c0V4ZWN1dGlvbkRhdGFEZX'
    'RhaWxzUgxpbnB1dGRldGFpbHMSFQoEbmFtZRjn++ZpIAEoCVIEbmFtZRIaCgZvdXRwdXQYhZ6l'
    'zQEgASgJUgZvdXRwdXQSVgoNb3V0cHV0ZGV0YWlscxjz09+7ASABKAsyLC5zdGF0ZXMuQ2xvdW'
    'RXYXRjaEV2ZW50c0V4ZWN1dGlvbkRhdGFEZXRhaWxzUg1vdXRwdXRkZXRhaWxzEiAKCXN0YXJ0'
    'ZGF0ZRicjvytASABKAlSCXN0YXJ0ZGF0ZRIsCg9zdGF0ZW1hY2hpbmVhcm4Y87vGuwEgASgJUg'
    '9zdGF0ZW1hY2hpbmVhcm4SNwoGc3RhdHVzGPDvrdIBIAEoDjIbLnN0YXRlcy5TeW5jRXhlY3V0'
    'aW9uU3RhdHVzUgZzdGF0dXMSHQoIc3RvcGRhdGUY2vKUViABKAlSCHN0b3BkYXRlEiMKC3RyYW'
    'NlaGVhZGVyGKCs8WggASgJUgt0cmFjZWhlYWRlcg==');

@$core.Deprecated('Use stateEnteredEventDetailsDescriptor instead')
const StateEnteredEventDetails$json = {
  '1': 'StateEnteredEventDetails',
  '2': [
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'inputdetails',
      '3': 452625788,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'inputdetails'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `StateEnteredEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateEnteredEventDetailsDescriptor = $convert.base64Decode(
    'ChhTdGF0ZUVudGVyZWRFdmVudERldGFpbHMSGAoFaW5wdXQY/N7hzgEgASgJUgVpbnB1dBJQCg'
    'xpbnB1dGRldGFpbHMY/Irq1wEgASgLMiguc3RhdGVzLkhpc3RvcnlFdmVudEV4ZWN1dGlvbkRh'
    'dGFEZXRhaWxzUgxpbnB1dGRldGFpbHMSFQoEbmFtZRjn++ZpIAEoCVIEbmFtZQ==');

@$core.Deprecated('Use stateExitedEventDetailsDescriptor instead')
const StateExitedEventDetails$json = {
  '1': 'StateExitedEventDetails',
  '2': [
    {
      '1': 'assignedvariables',
      '3': 411875019,
      '4': 3,
      '5': 11,
      '6': '.states.StateExitedEventDetails.AssignedvariablesEntry',
      '10': 'assignedvariables'
    },
    {
      '1': 'assignedvariablesdetails',
      '3': 509183609,
      '4': 1,
      '5': 11,
      '6': '.states.AssignedVariablesDetails',
      '10': 'assignedvariablesdetails'
    },
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputdetails',
      '3': 393734643,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'outputdetails'
    },
  ],
  '3': [StateExitedEventDetails_AssignedvariablesEntry$json],
};

@$core.Deprecated('Use stateExitedEventDetailsDescriptor instead')
const StateExitedEventDetails_AssignedvariablesEntry$json = {
  '1': 'AssignedvariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `StateExitedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateExitedEventDetailsDescriptor = $convert.base64Decode(
    'ChdTdGF0ZUV4aXRlZEV2ZW50RGV0YWlscxJoChFhc3NpZ25lZHZhcmlhYmxlcxjL7bLEASADKA'
    'syNi5zdGF0ZXMuU3RhdGVFeGl0ZWRFdmVudERldGFpbHMuQXNzaWduZWR2YXJpYWJsZXNFbnRy'
    'eVIRYXNzaWduZWR2YXJpYWJsZXMSYAoYYXNzaWduZWR2YXJpYWJsZXNkZXRhaWxzGPmM5vIBIA'
    'EoCzIgLnN0YXRlcy5Bc3NpZ25lZFZhcmlhYmxlc0RldGFpbHNSGGFzc2lnbmVkdmFyaWFibGVz'
    'ZGV0YWlscxIVCgRuYW1lGOf75mkgASgJUgRuYW1lEhoKBm91dHB1dBiFnqXNASABKAlSBm91dH'
    'B1dBJSCg1vdXRwdXRkZXRhaWxzGPPT37sBIAEoCzIoLnN0YXRlcy5IaXN0b3J5RXZlbnRFeGVj'
    'dXRpb25EYXRhRGV0YWlsc1INb3V0cHV0ZGV0YWlscxpEChZBc3NpZ25lZHZhcmlhYmxlc0VudH'
    'J5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use stateMachineAliasListItemDescriptor instead')
const StateMachineAliasListItem$json = {
  '1': 'StateMachineAliasListItem',
  '2': [
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
  ],
};

/// Descriptor for `StateMachineAliasListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateMachineAliasListItemDescriptor = $convert.base64Decode(
    'ChlTdGF0ZU1hY2hpbmVBbGlhc0xpc3RJdGVtEiUKDGNyZWF0aW9uZGF0ZRiBztFxIAEoCVIMY3'
    'JlYXRpb25kYXRlEjYKFHN0YXRlbWFjaGluZWFsaWFzYXJuGJHU8fwBIAEoCVIUc3RhdGVtYWNo'
    'aW5lYWxpYXNhcm4=');

@$core.Deprecated('Use stateMachineAlreadyExistsDescriptor instead')
const StateMachineAlreadyExists$json = {
  '1': 'StateMachineAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StateMachineAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateMachineAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChlTdGF0ZU1hY2hpbmVBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use stateMachineDeletingDescriptor instead')
const StateMachineDeleting$json = {
  '1': 'StateMachineDeleting',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StateMachineDeleting`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateMachineDeletingDescriptor =
    $convert.base64Decode(
        'ChRTdGF0ZU1hY2hpbmVEZWxldGluZxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use stateMachineDoesNotExistDescriptor instead')
const StateMachineDoesNotExist$json = {
  '1': 'StateMachineDoesNotExist',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StateMachineDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateMachineDoesNotExistDescriptor =
    $convert.base64Decode(
        'ChhTdGF0ZU1hY2hpbmVEb2VzTm90RXhpc3QSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use stateMachineLimitExceededDescriptor instead')
const StateMachineLimitExceeded$json = {
  '1': 'StateMachineLimitExceeded',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StateMachineLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateMachineLimitExceededDescriptor =
    $convert.base64Decode(
        'ChlTdGF0ZU1hY2hpbmVMaW1pdEV4Y2VlZGVkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use stateMachineListItemDescriptor instead')
const StateMachineListItem$json = {
  '1': 'StateMachineListItem',
  '2': [
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'name', '3': 221887975, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.states.StateMachineType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `StateMachineListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateMachineListItemDescriptor = $convert.base64Decode(
    'ChRTdGF0ZU1hY2hpbmVMaXN0SXRlbRIlCgxjcmVhdGlvbmRhdGUYgc7RcSABKAlSDGNyZWF0aW'
    '9uZGF0ZRIVCgRuYW1lGOf75mkgASgJUgRuYW1lEiwKD3N0YXRlbWFjaGluZWFybhjzu8a7ASAB'
    'KAlSD3N0YXRlbWFjaGluZWFybhIwCgR0eXBlGM7in4kBIAEoDjIYLnN0YXRlcy5TdGF0ZU1hY2'
    'hpbmVUeXBlUgR0eXBl');

@$core.Deprecated('Use stateMachineTypeNotSupportedDescriptor instead')
const StateMachineTypeNotSupported$json = {
  '1': 'StateMachineTypeNotSupported',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StateMachineTypeNotSupported`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateMachineTypeNotSupportedDescriptor =
    $convert.base64Decode(
        'ChxTdGF0ZU1hY2hpbmVUeXBlTm90U3VwcG9ydGVkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use stateMachineVersionListItemDescriptor instead')
const StateMachineVersionListItem$json = {
  '1': 'StateMachineVersionListItem',
  '2': [
    {'1': 'creationdate', '3': 238315265, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
  ],
};

/// Descriptor for `StateMachineVersionListItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stateMachineVersionListItemDescriptor =
    $convert.base64Decode(
        'ChtTdGF0ZU1hY2hpbmVWZXJzaW9uTGlzdEl0ZW0SJQoMY3JlYXRpb25kYXRlGIHO0XEgASgJUg'
        'xjcmVhdGlvbmRhdGUSOQoWc3RhdGVtYWNoaW5ldmVyc2lvbmFybhj5ha8hIAEoCVIWc3RhdGVt'
        'YWNoaW5ldmVyc2lvbmFybg==');

@$core.Deprecated('Use stopExecutionInputDescriptor instead')
const StopExecutionInput$json = {
  '1': 'StopExecutionInput',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'executionarn', '3': 314526573, '4': 1, '5': 9, '10': 'executionarn'},
  ],
};

/// Descriptor for `StopExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopExecutionInputDescriptor = $convert.base64Decode(
    'ChJTdG9wRXhlY3V0aW9uSW5wdXQSFwoFY2F1c2UYoaS7RSABKAlSBWNhdXNlEhcKBWVycm9yGN'
    'KOxgwgASgJUgVlcnJvchImCgxleGVjdXRpb25hcm4Y7Zb9lQEgASgJUgxleGVjdXRpb25hcm4=');

@$core.Deprecated('Use stopExecutionOutputDescriptor instead')
const StopExecutionOutput$json = {
  '1': 'StopExecutionOutput',
  '2': [
    {'1': 'stopdate', '3': 180697434, '4': 1, '5': 9, '10': 'stopdate'},
  ],
};

/// Descriptor for `StopExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopExecutionOutputDescriptor = $convert.base64Decode(
    'ChNTdG9wRXhlY3V0aW9uT3V0cHV0Eh0KCHN0b3BkYXRlGNrylFYgASgJUghzdG9wZGF0ZQ==');

@$core.Deprecated('Use tagDescriptor instead')
const Tag$json = {
  '1': 'Tag',
  '2': [
    {'1': 'key', '3': 135645293, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Tag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagDescriptor = $convert.base64Decode(
    'CgNUYWcSEwoDa2V5GO2Q10AgASgJUgNrZXkSFwoFdmFsdWUYy6f7EiABKAlSBXZhbHVl');

@$core.Deprecated('Use tagResourceInputDescriptor instead')
const TagResourceInput$json = {
  '1': 'TagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.states.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceInputDescriptor = $convert.base64Decode(
    'ChBUYWdSZXNvdXJjZUlucHV0EiMKC3Jlc291cmNlYXJuGM3MqiAgASgJUgtyZXNvdXJjZWFybh'
    'IjCgR0YWdzGKHX26ABIAMoCzILLnN0YXRlcy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use tagResourceOutputDescriptor instead')
const TagResourceOutput$json = {
  '1': 'TagResourceOutput',
};

/// Descriptor for `TagResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceOutputDescriptor =
    $convert.base64Decode('ChFUYWdSZXNvdXJjZU91dHB1dA==');

@$core.Deprecated('Use taskCredentialsDescriptor instead')
const TaskCredentials$json = {
  '1': 'TaskCredentials',
  '2': [
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
  ],
};

/// Descriptor for `TaskCredentials`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskCredentialsDescriptor = $convert.base64Decode(
    'Cg9UYXNrQ3JlZGVudGlhbHMSGwoHcm9sZWFybhihl4lRIAEoCVIHcm9sZWFybg==');

@$core.Deprecated('Use taskDoesNotExistDescriptor instead')
const TaskDoesNotExist$json = {
  '1': 'TaskDoesNotExist',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TaskDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskDoesNotExistDescriptor = $convert.base64Decode(
    'ChBUYXNrRG9lc05vdEV4aXN0EhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use taskFailedEventDetailsDescriptor instead')
const TaskFailedEventDetails$json = {
  '1': 'TaskFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `TaskFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskFailedEventDetailsDescriptor = $convert.base64Decode(
    'ChZUYXNrRmFpbGVkRXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0UgASgJUgVjYXVzZRIXCgVlcn'
    'JvchjSjsYMIAEoCVIFZXJyb3ISHQoIcmVzb3VyY2UY9v/9TiABKAlSCHJlc291cmNlEiUKDHJl'
    'c291cmNldHlwZRj+ldADIAEoCVIMcmVzb3VyY2V0eXBl');

@$core.Deprecated('Use taskScheduledEventDetailsDescriptor instead')
const TaskScheduledEventDetails$json = {
  '1': 'TaskScheduledEventDetails',
  '2': [
    {
      '1': 'heartbeatinseconds',
      '3': 125718754,
      '4': 1,
      '5': 3,
      '10': 'heartbeatinseconds'
    },
    {'1': 'parameters', '3': 145043162, '4': 1, '5': 9, '10': 'parameters'},
    {'1': 'region', '3': 52100734, '4': 1, '5': 9, '10': 'region'},
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
    {
      '1': 'taskcredentials',
      '3': 257843259,
      '4': 1,
      '5': 11,
      '6': '.states.TaskCredentials',
      '10': 'taskcredentials'
    },
    {
      '1': 'timeoutinseconds',
      '3': 472710197,
      '4': 1,
      '5': 3,
      '10': 'timeoutinseconds'
    },
  ],
};

/// Descriptor for `TaskScheduledEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskScheduledEventDetailsDescriptor = $convert.base64Decode(
    'ChlUYXNrU2NoZWR1bGVkRXZlbnREZXRhaWxzEjEKEmhlYXJ0YmVhdGluc2Vjb25kcxjiofk7IA'
    'EoA1ISaGVhcnRiZWF0aW5zZWNvbmRzEiEKCnBhcmFtZXRlcnMY2t2URSABKAlSCnBhcmFtZXRl'
    'cnMSGQoGcmVnaW9uGP786xggASgJUgZyZWdpb24SHQoIcmVzb3VyY2UY9v/9TiABKAlSCHJlc2'
    '91cmNlEiUKDHJlc291cmNldHlwZRj+ldADIAEoCVIMcmVzb3VyY2V0eXBlEkQKD3Rhc2tjcmVk'
    'ZW50aWFscxi7wPl6IAEoCzIXLnN0YXRlcy5UYXNrQ3JlZGVudGlhbHNSD3Rhc2tjcmVkZW50aW'
    'FscxIuChB0aW1lb3V0aW5zZWNvbmRzGLX4s+EBIAEoA1IQdGltZW91dGluc2Vjb25kcw==');

@$core.Deprecated('Use taskStartFailedEventDetailsDescriptor instead')
const TaskStartFailedEventDetails$json = {
  '1': 'TaskStartFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `TaskStartFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskStartFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChtUYXNrU3RhcnRGYWlsZWRFdmVudERldGFpbHMSFwoFY2F1c2UYoaS7RSABKAlSBWNhdXNlEh'
        'cKBWVycm9yGNKOxgwgASgJUgVlcnJvchIdCghyZXNvdXJjZRj2//1OIAEoCVIIcmVzb3VyY2US'
        'JQoMcmVzb3VyY2V0eXBlGP6V0AMgASgJUgxyZXNvdXJjZXR5cGU=');

@$core.Deprecated('Use taskStartedEventDetailsDescriptor instead')
const TaskStartedEventDetails$json = {
  '1': 'TaskStartedEventDetails',
  '2': [
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `TaskStartedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskStartedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChdUYXNrU3RhcnRlZEV2ZW50RGV0YWlscxIdCghyZXNvdXJjZRj2//1OIAEoCVIIcmVzb3VyY2'
        'USJQoMcmVzb3VyY2V0eXBlGP6V0AMgASgJUgxyZXNvdXJjZXR5cGU=');

@$core.Deprecated('Use taskSubmitFailedEventDetailsDescriptor instead')
const TaskSubmitFailedEventDetails$json = {
  '1': 'TaskSubmitFailedEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `TaskSubmitFailedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskSubmitFailedEventDetailsDescriptor =
    $convert.base64Decode(
        'ChxUYXNrU3VibWl0RmFpbGVkRXZlbnREZXRhaWxzEhcKBWNhdXNlGKGku0UgASgJUgVjYXVzZR'
        'IXCgVlcnJvchjSjsYMIAEoCVIFZXJyb3ISHQoIcmVzb3VyY2UY9v/9TiABKAlSCHJlc291cmNl'
        'EiUKDHJlc291cmNldHlwZRj+ldADIAEoCVIMcmVzb3VyY2V0eXBl');

@$core.Deprecated('Use taskSubmittedEventDetailsDescriptor instead')
const TaskSubmittedEventDetails$json = {
  '1': 'TaskSubmittedEventDetails',
  '2': [
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputdetails',
      '3': 393734643,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'outputdetails'
    },
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `TaskSubmittedEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskSubmittedEventDetailsDescriptor = $convert.base64Decode(
    'ChlUYXNrU3VibWl0dGVkRXZlbnREZXRhaWxzEhoKBm91dHB1dBiFnqXNASABKAlSBm91dHB1dB'
    'JSCg1vdXRwdXRkZXRhaWxzGPPT37sBIAEoCzIoLnN0YXRlcy5IaXN0b3J5RXZlbnRFeGVjdXRp'
    'b25EYXRhRGV0YWlsc1INb3V0cHV0ZGV0YWlscxIdCghyZXNvdXJjZRj2//1OIAEoCVIIcmVzb3'
    'VyY2USJQoMcmVzb3VyY2V0eXBlGP6V0AMgASgJUgxyZXNvdXJjZXR5cGU=');

@$core.Deprecated('Use taskSucceededEventDetailsDescriptor instead')
const TaskSucceededEventDetails$json = {
  '1': 'TaskSucceededEventDetails',
  '2': [
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputdetails',
      '3': 393734643,
      '4': 1,
      '5': 11,
      '6': '.states.HistoryEventExecutionDataDetails',
      '10': 'outputdetails'
    },
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `TaskSucceededEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskSucceededEventDetailsDescriptor = $convert.base64Decode(
    'ChlUYXNrU3VjY2VlZGVkRXZlbnREZXRhaWxzEhoKBm91dHB1dBiFnqXNASABKAlSBm91dHB1dB'
    'JSCg1vdXRwdXRkZXRhaWxzGPPT37sBIAEoCzIoLnN0YXRlcy5IaXN0b3J5RXZlbnRFeGVjdXRp'
    'b25EYXRhRGV0YWlsc1INb3V0cHV0ZGV0YWlscxIdCghyZXNvdXJjZRj2//1OIAEoCVIIcmVzb3'
    'VyY2USJQoMcmVzb3VyY2V0eXBlGP6V0AMgASgJUgxyZXNvdXJjZXR5cGU=');

@$core.Deprecated('Use taskTimedOutDescriptor instead')
const TaskTimedOut$json = {
  '1': 'TaskTimedOut',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TaskTimedOut`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskTimedOutDescriptor = $convert.base64Decode(
    'CgxUYXNrVGltZWRPdXQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use taskTimedOutEventDetailsDescriptor instead')
const TaskTimedOutEventDetails$json = {
  '1': 'TaskTimedOutEventDetails',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {'1': 'resource', '3': 165642230, '4': 1, '5': 9, '10': 'resource'},
    {'1': 'resourcetype', '3': 7604990, '4': 1, '5': 9, '10': 'resourcetype'},
  ],
};

/// Descriptor for `TaskTimedOutEventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List taskTimedOutEventDetailsDescriptor = $convert.base64Decode(
    'ChhUYXNrVGltZWRPdXRFdmVudERldGFpbHMSFwoFY2F1c2UYoaS7RSABKAlSBWNhdXNlEhcKBW'
    'Vycm9yGNKOxgwgASgJUgVlcnJvchIdCghyZXNvdXJjZRj2//1OIAEoCVIIcmVzb3VyY2USJQoM'
    'cmVzb3VyY2V0eXBlGP6V0AMgASgJUgxyZXNvdXJjZXR5cGU=');

@$core.Deprecated('Use testStateConfigurationDescriptor instead')
const TestStateConfiguration$json = {
  '1': 'TestStateConfiguration',
  '2': [
    {
      '1': 'errorcausedbystate',
      '3': 313873103,
      '4': 1,
      '5': 9,
      '10': 'errorcausedbystate'
    },
    {
      '1': 'mapitemreaderdata',
      '3': 503260646,
      '4': 1,
      '5': 9,
      '10': 'mapitemreaderdata'
    },
    {
      '1': 'mapiterationfailurecount',
      '3': 419587762,
      '4': 1,
      '5': 5,
      '10': 'mapiterationfailurecount'
    },
    {
      '1': 'retrierretrycount',
      '3': 275735648,
      '4': 1,
      '5': 5,
      '10': 'retrierretrycount'
    },
  ],
};

/// Descriptor for `TestStateConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testStateConfigurationDescriptor = $convert.base64Decode(
    'ChZUZXN0U3RhdGVDb25maWd1cmF0aW9uEjIKEmVycm9yY2F1c2VkYnlzdGF0ZRjPpdWVASABKA'
    'lSEmVycm9yY2F1c2VkYnlzdGF0ZRIwChFtYXBpdGVtcmVhZGVyZGF0YRjmy/zvASABKAlSEW1h'
    'cGl0ZW1yZWFkZXJkYXRhEj4KGG1hcGl0ZXJhdGlvbmZhaWx1cmVjb3VudBiyzYnIASABKAVSGG'
    '1hcGl0ZXJhdGlvbmZhaWx1cmVjb3VudBIwChFyZXRyaWVycmV0cnljb3VudBjgyL2DASABKAVS'
    'EXJldHJpZXJyZXRyeWNvdW50');

@$core.Deprecated('Use testStateInputDescriptor instead')
const TestStateInput$json = {
  '1': 'TestStateInput',
  '2': [
    {'1': 'context', '3': 210178173, '4': 1, '5': 9, '10': 'context'},
    {'1': 'definition', '3': 68443297, '4': 1, '5': 9, '10': 'definition'},
    {'1': 'input', '3': 433614716, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'inspectionlevel',
      '3': 277169476,
      '4': 1,
      '5': 14,
      '6': '.states.InspectionLevel',
      '10': 'inspectionlevel'
    },
    {
      '1': 'mock',
      '3': 242883628,
      '4': 1,
      '5': 11,
      '6': '.states.MockInput',
      '10': 'mock'
    },
    {
      '1': 'revealsecrets',
      '3': 351839742,
      '4': 1,
      '5': 8,
      '10': 'revealsecrets'
    },
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'stateconfiguration',
      '3': 17002877,
      '4': 1,
      '5': 11,
      '6': '.states.TestStateConfiguration',
      '10': 'stateconfiguration'
    },
    {'1': 'statename', '3': 270657590, '4': 1, '5': 9, '10': 'statename'},
    {'1': 'variables', '3': 162226883, '4': 1, '5': 9, '10': 'variables'},
  ],
};

/// Descriptor for `TestStateInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testStateInputDescriptor = $convert.base64Decode(
    'Cg5UZXN0U3RhdGVJbnB1dBIbCgdjb250ZXh0GP2gnGQgASgJUgdjb250ZXh0EiEKCmRlZmluaX'
    'Rpb24YobnRICABKAlSCmRlZmluaXRpb24SGAoFaW5wdXQY/N7hzgEgASgJUgVpbnB1dBJFCg9p'
    'bnNwZWN0aW9ubGV2ZWwYxIqVhAEgASgOMhcuc3RhdGVzLkluc3BlY3Rpb25MZXZlbFIPaW5zcG'
    'VjdGlvbmxldmVsEigKBG1vY2sYrLjocyABKAsyES5zdGF0ZXMuTW9ja0lucHV0UgRtb2NrEigK'
    'DXJldmVhbHNlY3JldHMY/svipwEgASgIUg1yZXZlYWxzZWNyZXRzEhsKB3JvbGVhcm4YoZeJUS'
    'ABKAlSB3JvbGVhcm4SUQoSc3RhdGVjb25maWd1cmF0aW9uGP3ijQggASgLMh4uc3RhdGVzLlRl'
    'c3RTdGF0ZUNvbmZpZ3VyYXRpb25SEnN0YXRlY29uZmlndXJhdGlvbhIgCglzdGF0ZW5hbWUYtt'
    'CHgQEgASgJUglzdGF0ZW5hbWUSHwoJdmFyaWFibGVzGMPFrU0gASgJUgl2YXJpYWJsZXM=');

@$core.Deprecated('Use testStateOutputDescriptor instead')
const TestStateOutput$json = {
  '1': 'TestStateOutput',
  '2': [
    {'1': 'cause', '3': 145674785, '4': 1, '5': 9, '10': 'cause'},
    {'1': 'error', '3': 26314578, '4': 1, '5': 9, '10': 'error'},
    {
      '1': 'inspectiondata',
      '3': 113762044,
      '4': 1,
      '5': 11,
      '6': '.states.InspectionData',
      '10': 'inspectiondata'
    },
    {'1': 'nextstate', '3': 525594702, '4': 1, '5': 9, '10': 'nextstate'},
    {'1': 'output', '3': 430526213, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'status',
      '3': 441153520,
      '4': 1,
      '5': 14,
      '6': '.states.TestExecutionStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `TestStateOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testStateOutputDescriptor = $convert.base64Decode(
    'Cg9UZXN0U3RhdGVPdXRwdXQSFwoFY2F1c2UYoaS7RSABKAlSBWNhdXNlEhcKBWVycm9yGNKOxg'
    'wgASgJUgVlcnJvchJBCg5pbnNwZWN0aW9uZGF0YRj8vZ82IAEoCzIWLnN0YXRlcy5JbnNwZWN0'
    'aW9uRGF0YVIOaW5zcGVjdGlvbmRhdGESIAoJbmV4dHN0YXRlGM7gz/oBIAEoCVIJbmV4dHN0YX'
    'RlEhoKBm91dHB1dBiFnqXNASABKAlSBm91dHB1dBI3CgZzdGF0dXMY8O+t0gEgASgOMhsuc3Rh'
    'dGVzLlRlc3RFeGVjdXRpb25TdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use tooManyTagsDescriptor instead')
const TooManyTags$json = {
  '1': 'TooManyTags',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {'1': 'resourcename', '3': 17776375, '4': 1, '5': 9, '10': 'resourcename'},
  ],
};

/// Descriptor for `TooManyTags`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyTagsDescriptor = $convert.base64Decode(
    'CgtUb29NYW55VGFncxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdlEiUKDHJlc291cmNlbm'
    'FtZRj3/bwIIAEoCVIMcmVzb3VyY2VuYW1l');

@$core.Deprecated('Use tracingConfigurationDescriptor instead')
const TracingConfiguration$json = {
  '1': 'TracingConfiguration',
  '2': [
    {'1': 'enabled', '3': 49525663, '4': 1, '5': 8, '10': 'enabled'},
  ],
};

/// Descriptor for `TracingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tracingConfigurationDescriptor =
    $convert.base64Decode(
        'ChRUcmFjaW5nQ29uZmlndXJhdGlvbhIbCgdlbmFibGVkGJ/nzhcgASgIUgdlbmFibGVk');

@$core.Deprecated('Use untagResourceInputDescriptor instead')
const UntagResourceInput$json = {
  '1': 'UntagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 67806797, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 78811036, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceInputDescriptor = $convert.base64Decode(
    'ChJVbnRhZ1Jlc291cmNlSW5wdXQSIwoLcmVzb3VyY2Vhcm4YzcyqICABKAlSC3Jlc291cmNlYX'
    'JuEhsKB3RhZ2tleXMYnJ/KJSADKAlSB3RhZ2tleXM=');

@$core.Deprecated('Use untagResourceOutputDescriptor instead')
const UntagResourceOutput$json = {
  '1': 'UntagResourceOutput',
};

/// Descriptor for `UntagResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceOutputDescriptor =
    $convert.base64Decode('ChNVbnRhZ1Jlc291cmNlT3V0cHV0');

@$core.Deprecated('Use updateMapRunInputDescriptor instead')
const UpdateMapRunInput$json = {
  '1': 'UpdateMapRunInput',
  '2': [
    {'1': 'maprunarn', '3': 18199994, '4': 1, '5': 9, '10': 'maprunarn'},
    {
      '1': 'maxconcurrency',
      '3': 100901405,
      '4': 1,
      '5': 5,
      '10': 'maxconcurrency'
    },
    {
      '1': 'toleratedfailurecount',
      '3': 41834811,
      '4': 1,
      '5': 3,
      '10': 'toleratedfailurecount'
    },
    {
      '1': 'toleratedfailurepercentage',
      '3': 116496164,
      '4': 1,
      '5': 2,
      '10': 'toleratedfailurepercentage'
    },
  ],
};

/// Descriptor for `UpdateMapRunInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMapRunInputDescriptor = $convert.base64Decode(
    'ChFVcGRhdGVNYXBSdW5JbnB1dBIfCgltYXBydW5hcm4YuuvWCCABKAlSCW1hcHJ1bmFybhIpCg'
    '5tYXhjb25jdXJyZW5jeRidxI4wIAEoBVIObWF4Y29uY3VycmVuY3kSNwoVdG9sZXJhdGVkZmFp'
    'bHVyZWNvdW50GLuy+RMgASgDUhV0b2xlcmF0ZWRmYWlsdXJlY291bnQSQQoadG9sZXJhdGVkZm'
    'FpbHVyZXBlcmNlbnRhZ2UYpK7GNyABKAJSGnRvbGVyYXRlZGZhaWx1cmVwZXJjZW50YWdl');

@$core.Deprecated('Use updateMapRunOutputDescriptor instead')
const UpdateMapRunOutput$json = {
  '1': 'UpdateMapRunOutput',
};

/// Descriptor for `UpdateMapRunOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMapRunOutputDescriptor =
    $convert.base64Decode('ChJVcGRhdGVNYXBSdW5PdXRwdXQ=');

@$core.Deprecated('Use updateStateMachineAliasInputDescriptor instead')
const UpdateStateMachineAliasInput$json = {
  '1': 'UpdateStateMachineAliasInput',
  '2': [
    {'1': 'description', '3': 342834026, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'routingconfiguration',
      '3': 372891510,
      '4': 3,
      '5': 11,
      '6': '.states.RoutingConfigurationListItem',
      '10': 'routingconfiguration'
    },
    {
      '1': 'statemachinealiasarn',
      '3': 530344465,
      '4': 1,
      '5': 9,
      '10': 'statemachinealiasarn'
    },
  ],
};

/// Descriptor for `UpdateStateMachineAliasInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStateMachineAliasInputDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVTdGF0ZU1hY2hpbmVBbGlhc0lucHV0EiQKC2Rlc2NyaXB0aW9uGOr2vKMBIAEoCV'
    'ILZGVzY3JpcHRpb24SXAoUcm91dGluZ2NvbmZpZ3VyYXRpb24Y9r7nsQEgAygLMiQuc3RhdGVz'
    'LlJvdXRpbmdDb25maWd1cmF0aW9uTGlzdEl0ZW1SFHJvdXRpbmdjb25maWd1cmF0aW9uEjYKFH'
    'N0YXRlbWFjaGluZWFsaWFzYXJuGJHU8fwBIAEoCVIUc3RhdGVtYWNoaW5lYWxpYXNhcm4=');

@$core.Deprecated('Use updateStateMachineAliasOutputDescriptor instead')
const UpdateStateMachineAliasOutput$json = {
  '1': 'UpdateStateMachineAliasOutput',
  '2': [
    {'1': 'updatedate', '3': 510552561, '4': 1, '5': 9, '10': 'updatedate'},
  ],
};

/// Descriptor for `UpdateStateMachineAliasOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStateMachineAliasOutputDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVTdGF0ZU1hY2hpbmVBbGlhc091dHB1dBIiCgp1cGRhdGVkYXRlGPHTufMBIAEoCV'
        'IKdXBkYXRlZGF0ZQ==');

@$core.Deprecated('Use updateStateMachineInputDescriptor instead')
const UpdateStateMachineInput$json = {
  '1': 'UpdateStateMachineInput',
  '2': [
    {'1': 'definition', '3': 68443297, '4': 1, '5': 9, '10': 'definition'},
    {
      '1': 'encryptionconfiguration',
      '3': 167857431,
      '4': 1,
      '5': 11,
      '6': '.states.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {
      '1': 'loggingconfiguration',
      '3': 420811605,
      '4': 1,
      '5': 11,
      '6': '.states.LoggingConfiguration',
      '10': 'loggingconfiguration'
    },
    {'1': 'publish', '3': 264247305, '4': 1, '5': 8, '10': 'publish'},
    {'1': 'rolearn', '3': 170019745, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'statemachinearn',
      '3': 393321971,
      '4': 1,
      '5': 9,
      '10': 'statemachinearn'
    },
    {
      '1': 'tracingconfiguration',
      '3': 491315910,
      '4': 1,
      '5': 11,
      '6': '.states.TracingConfiguration',
      '10': 'tracingconfiguration'
    },
    {
      '1': 'versiondescription',
      '3': 434714300,
      '4': 1,
      '5': 9,
      '10': 'versiondescription'
    },
  ],
};

/// Descriptor for `UpdateStateMachineInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStateMachineInputDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVTdGF0ZU1hY2hpbmVJbnB1dBIhCgpkZWZpbml0aW9uGKG50SAgASgJUgpkZWZpbm'
    'l0aW9uElwKF2VuY3J5cHRpb25jb25maWd1cmF0aW9uGJeahVAgASgLMh8uc3RhdGVzLkVuY3J5'
    'cHRpb25Db25maWd1cmF0aW9uUhdlbmNyeXB0aW9uY29uZmlndXJhdGlvbhJUChRsb2dnaW5nY2'
    '9uZmlndXJhdGlvbhjVptTIASABKAsyHC5zdGF0ZXMuTG9nZ2luZ0NvbmZpZ3VyYXRpb25SFGxv'
    'Z2dpbmdjb25maWd1cmF0aW9uEhsKB3B1Ymxpc2gYibCAfiABKAhSB3B1Ymxpc2gSGwoHcm9sZW'
    'Fybhihl4lRIAEoCVIHcm9sZWFybhIsCg9zdGF0ZW1hY2hpbmVhcm4Y87vGuwEgASgJUg9zdGF0'
    'ZW1hY2hpbmVhcm4SVAoUdHJhY2luZ2NvbmZpZ3VyYXRpb24YxsWj6gEgASgLMhwuc3RhdGVzLl'
    'RyYWNpbmdDb25maWd1cmF0aW9uUhR0cmFjaW5nY29uZmlndXJhdGlvbhIyChJ2ZXJzaW9uZGVz'
    'Y3JpcHRpb24YvO2kzwEgASgJUhJ2ZXJzaW9uZGVzY3JpcHRpb24=');

@$core.Deprecated('Use updateStateMachineOutputDescriptor instead')
const UpdateStateMachineOutput$json = {
  '1': 'UpdateStateMachineOutput',
  '2': [
    {'1': 'revisionid', '3': 369170086, '4': 1, '5': 9, '10': 'revisionid'},
    {
      '1': 'statemachineversionarn',
      '3': 69976825,
      '4': 1,
      '5': 9,
      '10': 'statemachineversionarn'
    },
    {'1': 'updatedate', '3': 510552561, '4': 1, '5': 9, '10': 'updatedate'},
  ],
};

/// Descriptor for `UpdateStateMachineOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStateMachineOutputDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVTdGF0ZU1hY2hpbmVPdXRwdXQSIgoKcmV2aXNpb25pZBimrYSwASABKAlSCnJldm'
    'lzaW9uaWQSOQoWc3RhdGVtYWNoaW5ldmVyc2lvbmFybhj5ha8hIAEoCVIWc3RhdGVtYWNoaW5l'
    'dmVyc2lvbmFybhIiCgp1cGRhdGVkYXRlGPHTufMBIAEoCVIKdXBkYXRlZGF0ZQ==');

@$core.Deprecated(
    'Use validateStateMachineDefinitionDiagnosticDescriptor instead')
const ValidateStateMachineDefinitionDiagnostic$json = {
  '1': 'ValidateStateMachineDefinitionDiagnostic',
  '2': [
    {'1': 'code', '3': 422669557, '4': 1, '5': 9, '10': 'code'},
    {'1': 'location', '3': 200649127, '4': 1, '5': 9, '10': 'location'},
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'severity',
      '3': 268193715,
      '4': 1,
      '5': 14,
      '6': '.states.ValidateStateMachineDefinitionSeverity',
      '10': 'severity'
    },
  ],
};

/// Descriptor for `ValidateStateMachineDefinitionDiagnostic`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateStateMachineDefinitionDiagnosticDescriptor =
    $convert.base64Decode(
        'CihWYWxpZGF0ZVN0YXRlTWFjaGluZURlZmluaXRpb25EaWFnbm9zdGljEhYKBGNvZGUY9dnFyQ'
        'EgASgJUgRjb2RlEh0KCGxvY2F0aW9uGKfT1l8gASgJUghsb2NhdGlvbhIbCgdtZXNzYWdlGOWR'
        'yCcgASgJUgdtZXNzYWdlEk0KCHNldmVyaXR5GLOf8X8gASgOMi4uc3RhdGVzLlZhbGlkYXRlU3'
        'RhdGVNYWNoaW5lRGVmaW5pdGlvblNldmVyaXR5UghzZXZlcml0eQ==');

@$core.Deprecated('Use validateStateMachineDefinitionInputDescriptor instead')
const ValidateStateMachineDefinitionInput$json = {
  '1': 'ValidateStateMachineDefinitionInput',
  '2': [
    {'1': 'definition', '3': 68443297, '4': 1, '5': 9, '10': 'definition'},
    {'1': 'maxresults', '3': 465170002, '4': 1, '5': 5, '10': 'maxresults'},
    {
      '1': 'severity',
      '3': 268193715,
      '4': 1,
      '5': 14,
      '6': '.states.ValidateStateMachineDefinitionSeverity',
      '10': 'severity'
    },
    {
      '1': 'type',
      '3': 287830350,
      '4': 1,
      '5': 14,
      '6': '.states.StateMachineType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ValidateStateMachineDefinitionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateStateMachineDefinitionInputDescriptor =
    $convert.base64Decode(
        'CiNWYWxpZGF0ZVN0YXRlTWFjaGluZURlZmluaXRpb25JbnB1dBIhCgpkZWZpbml0aW9uGKG50S'
        'AgASgJUgpkZWZpbml0aW9uEiIKCm1heHJlc3VsdHMY0tzn3QEgASgFUgptYXhyZXN1bHRzEk0K'
        'CHNldmVyaXR5GLOf8X8gASgOMi4uc3RhdGVzLlZhbGlkYXRlU3RhdGVNYWNoaW5lRGVmaW5pdG'
        'lvblNldmVyaXR5UghzZXZlcml0eRIwCgR0eXBlGM7in4kBIAEoDjIYLnN0YXRlcy5TdGF0ZU1h'
        'Y2hpbmVUeXBlUgR0eXBl');

@$core.Deprecated('Use validateStateMachineDefinitionOutputDescriptor instead')
const ValidateStateMachineDefinitionOutput$json = {
  '1': 'ValidateStateMachineDefinitionOutput',
  '2': [
    {
      '1': 'diagnostics',
      '3': 159421432,
      '4': 3,
      '5': 11,
      '6': '.states.ValidateStateMachineDefinitionDiagnostic',
      '10': 'diagnostics'
    },
    {
      '1': 'result',
      '3': 171406885,
      '4': 1,
      '5': 14,
      '6': '.states.ValidateStateMachineDefinitionResultCode',
      '10': 'result'
    },
    {'1': 'truncated', '3': 407490858, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `ValidateStateMachineDefinitionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateStateMachineDefinitionOutputDescriptor =
    $convert.base64Decode(
        'CiRWYWxpZGF0ZVN0YXRlTWFjaGluZURlZmluaXRpb25PdXRwdXQSVQoLZGlhZ25vc3RpY3MY+K'
        'eCTCADKAsyMC5zdGF0ZXMuVmFsaWRhdGVTdGF0ZU1hY2hpbmVEZWZpbml0aW9uRGlhZ25vc3Rp'
        'Y1ILZGlhZ25vc3RpY3MSSwoGcmVzdWx0GKXs3VEgASgOMjAuc3RhdGVzLlZhbGlkYXRlU3RhdG'
        'VNYWNoaW5lRGVmaW5pdGlvblJlc3VsdENvZGVSBnJlc3VsdBIgCgl0cnVuY2F0ZWQYqqKnwgEg'
        'ASgIUgl0cnVuY2F0ZWQ=');

@$core.Deprecated('Use validationExceptionDescriptor instead')
const ValidationException$json = {
  '1': 'ValidationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'reason',
      '3': 413359642,
      '4': 1,
      '5': 14,
      '6': '.states.ValidationExceptionReason',
      '10': 'reason'
    },
  ],
};

/// Descriptor for `ValidationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationExceptionDescriptor = $convert.base64Decode(
    'ChNWYWxpZGF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2USPQoGcm'
    'Vhc29uGJq8jcUBIAEoDjIhLnN0YXRlcy5WYWxpZGF0aW9uRXhjZXB0aW9uUmVhc29uUgZyZWFz'
    'b24=');
