// This is a generated file - do not edit.
//
// Generated from query.timestream.proto.

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

@$core.Deprecated('Use computeModeDescriptor instead')
const ComputeMode$json = {
  '1': 'ComputeMode',
  '2': [
    {'1': 'COMPUTE_MODE_PROVISIONED', '2': 0},
    {'1': 'COMPUTE_MODE_ON_DEMAND', '2': 1},
  ],
};

/// Descriptor for `ComputeMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List computeModeDescriptor = $convert.base64Decode(
    'CgtDb21wdXRlTW9kZRIcChhDT01QVVRFX01PREVfUFJPVklTSU9ORUQQABIaChZDT01QVVRFX0'
    '1PREVfT05fREVNQU5EEAE=');

@$core.Deprecated('Use dimensionValueTypeDescriptor instead')
const DimensionValueType$json = {
  '1': 'DimensionValueType',
  '2': [
    {'1': 'DIMENSION_VALUE_TYPE_VARCHAR', '2': 0},
  ],
};

/// Descriptor for `DimensionValueType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dimensionValueTypeDescriptor = $convert.base64Decode(
    'ChJEaW1lbnNpb25WYWx1ZVR5cGUSIAocRElNRU5TSU9OX1ZBTFVFX1RZUEVfVkFSQ0hBUhAA');

@$core.Deprecated('Use lastUpdateStatusDescriptor instead')
const LastUpdateStatus$json = {
  '1': 'LastUpdateStatus',
  '2': [
    {'1': 'LAST_UPDATE_STATUS_PENDING', '2': 0},
    {'1': 'LAST_UPDATE_STATUS_SUCCEEDED', '2': 1},
    {'1': 'LAST_UPDATE_STATUS_FAILED', '2': 2},
  ],
};

/// Descriptor for `LastUpdateStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List lastUpdateStatusDescriptor = $convert.base64Decode(
    'ChBMYXN0VXBkYXRlU3RhdHVzEh4KGkxBU1RfVVBEQVRFX1NUQVRVU19QRU5ESU5HEAASIAocTE'
    'FTVF9VUERBVEVfU1RBVFVTX1NVQ0NFRURFRBABEh0KGUxBU1RfVVBEQVRFX1NUQVRVU19GQUlM'
    'RUQQAg==');

@$core.Deprecated('Use measureValueTypeDescriptor instead')
const MeasureValueType$json = {
  '1': 'MeasureValueType',
  '2': [
    {'1': 'MEASURE_VALUE_TYPE_BIGINT', '2': 0},
    {'1': 'MEASURE_VALUE_TYPE_MULTI', '2': 1},
    {'1': 'MEASURE_VALUE_TYPE_VARCHAR', '2': 2},
    {'1': 'MEASURE_VALUE_TYPE_BOOLEAN', '2': 3},
    {'1': 'MEASURE_VALUE_TYPE_DOUBLE', '2': 4},
  ],
};

/// Descriptor for `MeasureValueType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List measureValueTypeDescriptor = $convert.base64Decode(
    'ChBNZWFzdXJlVmFsdWVUeXBlEh0KGU1FQVNVUkVfVkFMVUVfVFlQRV9CSUdJTlQQABIcChhNRU'
    'FTVVJFX1ZBTFVFX1RZUEVfTVVMVEkQARIeChpNRUFTVVJFX1ZBTFVFX1RZUEVfVkFSQ0hBUhAC'
    'Eh4KGk1FQVNVUkVfVkFMVUVfVFlQRV9CT09MRUFOEAMSHQoZTUVBU1VSRV9WQUxVRV9UWVBFX0'
    'RPVUJMRRAE');

@$core.Deprecated('Use queryInsightsModeDescriptor instead')
const QueryInsightsMode$json = {
  '1': 'QueryInsightsMode',
  '2': [
    {'1': 'QUERY_INSIGHTS_MODE_DISABLED', '2': 0},
    {'1': 'QUERY_INSIGHTS_MODE_ENABLED_WITH_RATE_CONTROL', '2': 1},
  ],
};

/// Descriptor for `QueryInsightsMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queryInsightsModeDescriptor = $convert.base64Decode(
    'ChFRdWVyeUluc2lnaHRzTW9kZRIgChxRVUVSWV9JTlNJR0hUU19NT0RFX0RJU0FCTEVEEAASMQ'
    'otUVVFUllfSU5TSUdIVFNfTU9ERV9FTkFCTEVEX1dJVEhfUkFURV9DT05UUk9MEAE=');

@$core.Deprecated('Use queryPricingModelDescriptor instead')
const QueryPricingModel$json = {
  '1': 'QueryPricingModel',
  '2': [
    {'1': 'QUERY_PRICING_MODEL_COMPUTE_UNITS', '2': 0},
    {'1': 'QUERY_PRICING_MODEL_BYTES_SCANNED', '2': 1},
  ],
};

/// Descriptor for `QueryPricingModel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queryPricingModelDescriptor = $convert.base64Decode(
    'ChFRdWVyeVByaWNpbmdNb2RlbBIlCiFRVUVSWV9QUklDSU5HX01PREVMX0NPTVBVVEVfVU5JVF'
    'MQABIlCiFRVUVSWV9QUklDSU5HX01PREVMX0JZVEVTX1NDQU5ORUQQAQ==');

@$core.Deprecated('Use s3EncryptionOptionDescriptor instead')
const S3EncryptionOption$json = {
  '1': 'S3EncryptionOption',
  '2': [
    {'1': 'S3_ENCRYPTION_OPTION_SSE_S3', '2': 0},
    {'1': 'S3_ENCRYPTION_OPTION_SSE_KMS', '2': 1},
  ],
};

/// Descriptor for `S3EncryptionOption`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List s3EncryptionOptionDescriptor = $convert.base64Decode(
    'ChJTM0VuY3J5cHRpb25PcHRpb24SHwobUzNfRU5DUllQVElPTl9PUFRJT05fU1NFX1MzEAASIA'
    'ocUzNfRU5DUllQVElPTl9PUFRJT05fU1NFX0tNUxAB');

@$core.Deprecated('Use scalarMeasureValueTypeDescriptor instead')
const ScalarMeasureValueType$json = {
  '1': 'ScalarMeasureValueType',
  '2': [
    {'1': 'SCALAR_MEASURE_VALUE_TYPE_BIGINT', '2': 0},
    {'1': 'SCALAR_MEASURE_VALUE_TYPE_VARCHAR', '2': 1},
    {'1': 'SCALAR_MEASURE_VALUE_TYPE_BOOLEAN', '2': 2},
    {'1': 'SCALAR_MEASURE_VALUE_TYPE_DOUBLE', '2': 3},
    {'1': 'SCALAR_MEASURE_VALUE_TYPE_TIMESTAMP', '2': 4},
  ],
};

/// Descriptor for `ScalarMeasureValueType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scalarMeasureValueTypeDescriptor = $convert.base64Decode(
    'ChZTY2FsYXJNZWFzdXJlVmFsdWVUeXBlEiQKIFNDQUxBUl9NRUFTVVJFX1ZBTFVFX1RZUEVfQk'
    'lHSU5UEAASJQohU0NBTEFSX01FQVNVUkVfVkFMVUVfVFlQRV9WQVJDSEFSEAESJQohU0NBTEFS'
    'X01FQVNVUkVfVkFMVUVfVFlQRV9CT09MRUFOEAISJAogU0NBTEFSX01FQVNVUkVfVkFMVUVfVF'
    'lQRV9ET1VCTEUQAxInCiNTQ0FMQVJfTUVBU1VSRV9WQUxVRV9UWVBFX1RJTUVTVEFNUBAE');

@$core.Deprecated('Use scalarTypeDescriptor instead')
const ScalarType$json = {
  '1': 'ScalarType',
  '2': [
    {'1': 'SCALAR_TYPE_BIGINT', '2': 0},
    {'1': 'SCALAR_TYPE_UNKNOWN', '2': 1},
    {'1': 'SCALAR_TYPE_INTERVAL_DAY_TO_SECOND', '2': 2},
    {'1': 'SCALAR_TYPE_INTEGER', '2': 3},
    {'1': 'SCALAR_TYPE_TIME', '2': 4},
    {'1': 'SCALAR_TYPE_VARCHAR', '2': 5},
    {'1': 'SCALAR_TYPE_BOOLEAN', '2': 6},
    {'1': 'SCALAR_TYPE_DOUBLE', '2': 7},
    {'1': 'SCALAR_TYPE_TIMESTAMP', '2': 8},
    {'1': 'SCALAR_TYPE_DATE', '2': 9},
    {'1': 'SCALAR_TYPE_INTERVAL_YEAR_TO_MONTH', '2': 10},
  ],
};

/// Descriptor for `ScalarType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scalarTypeDescriptor = $convert.base64Decode(
    'CgpTY2FsYXJUeXBlEhYKElNDQUxBUl9UWVBFX0JJR0lOVBAAEhcKE1NDQUxBUl9UWVBFX1VOS0'
    '5PV04QARImCiJTQ0FMQVJfVFlQRV9JTlRFUlZBTF9EQVlfVE9fU0VDT05EEAISFwoTU0NBTEFS'
    'X1RZUEVfSU5URUdFUhADEhQKEFNDQUxBUl9UWVBFX1RJTUUQBBIXChNTQ0FMQVJfVFlQRV9WQV'
    'JDSEFSEAUSFwoTU0NBTEFSX1RZUEVfQk9PTEVBThAGEhYKElNDQUxBUl9UWVBFX0RPVUJMRRAH'
    'EhkKFVNDQUxBUl9UWVBFX1RJTUVTVEFNUBAIEhQKEFNDQUxBUl9UWVBFX0RBVEUQCRImCiJTQ0'
    'FMQVJfVFlQRV9JTlRFUlZBTF9ZRUFSX1RPX01PTlRIEAo=');

@$core.Deprecated('Use scheduledQueryInsightsModeDescriptor instead')
const ScheduledQueryInsightsMode$json = {
  '1': 'ScheduledQueryInsightsMode',
  '2': [
    {'1': 'SCHEDULED_QUERY_INSIGHTS_MODE_DISABLED', '2': 0},
    {'1': 'SCHEDULED_QUERY_INSIGHTS_MODE_ENABLED_WITH_RATE_CONTROL', '2': 1},
  ],
};

/// Descriptor for `ScheduledQueryInsightsMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scheduledQueryInsightsModeDescriptor =
    $convert.base64Decode(
        'ChpTY2hlZHVsZWRRdWVyeUluc2lnaHRzTW9kZRIqCiZTQ0hFRFVMRURfUVVFUllfSU5TSUdIVF'
        'NfTU9ERV9ESVNBQkxFRBAAEjsKN1NDSEVEVUxFRF9RVUVSWV9JTlNJR0hUU19NT0RFX0VOQUJM'
        'RURfV0lUSF9SQVRFX0NPTlRST0wQAQ==');

@$core.Deprecated('Use scheduledQueryRunStatusDescriptor instead')
const ScheduledQueryRunStatus$json = {
  '1': 'ScheduledQueryRunStatus',
  '2': [
    {'1': 'SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_SUCCESS', '2': 0},
    {'1': 'SCHEDULED_QUERY_RUN_STATUS_MANUAL_TRIGGER_FAILURE', '2': 1},
    {'1': 'SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_FAILURE', '2': 2},
    {'1': 'SCHEDULED_QUERY_RUN_STATUS_AUTO_TRIGGER_SUCCESS', '2': 3},
  ],
};

/// Descriptor for `ScheduledQueryRunStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scheduledQueryRunStatusDescriptor = $convert.base64Decode(
    'ChdTY2hlZHVsZWRRdWVyeVJ1blN0YXR1cxI1CjFTQ0hFRFVMRURfUVVFUllfUlVOX1NUQVRVU1'
    '9NQU5VQUxfVFJJR0dFUl9TVUNDRVNTEAASNQoxU0NIRURVTEVEX1FVRVJZX1JVTl9TVEFUVVNf'
    'TUFOVUFMX1RSSUdHRVJfRkFJTFVSRRABEjMKL1NDSEVEVUxFRF9RVUVSWV9SVU5fU1RBVFVTX0'
    'FVVE9fVFJJR0dFUl9GQUlMVVJFEAISMwovU0NIRURVTEVEX1FVRVJZX1JVTl9TVEFUVVNfQVVU'
    'T19UUklHR0VSX1NVQ0NFU1MQAw==');

@$core.Deprecated('Use scheduledQueryStateDescriptor instead')
const ScheduledQueryState$json = {
  '1': 'ScheduledQueryState',
  '2': [
    {'1': 'SCHEDULED_QUERY_STATE_DISABLED', '2': 0},
    {'1': 'SCHEDULED_QUERY_STATE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `ScheduledQueryState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scheduledQueryStateDescriptor = $convert.base64Decode(
    'ChNTY2hlZHVsZWRRdWVyeVN0YXRlEiIKHlNDSEVEVUxFRF9RVUVSWV9TVEFURV9ESVNBQkxFRB'
    'AAEiEKHVNDSEVEVUxFRF9RVUVSWV9TVEFURV9FTkFCTEVEEAE=');

@$core.Deprecated('Use accessDeniedExceptionDescriptor instead')
const AccessDeniedException$json = {
  '1': 'AccessDeniedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDeniedExceptionDescriptor = $convert.base64Decode(
    'ChVBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated(
    'Use accountSettingsNotificationConfigurationDescriptor instead')
const AccountSettingsNotificationConfiguration$json = {
  '1': 'AccountSettingsNotificationConfiguration',
  '2': [
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'snsconfiguration',
      '3': 6844682,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.SnsConfiguration',
      '10': 'snsconfiguration'
    },
  ],
};

/// Descriptor for `AccountSettingsNotificationConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountSettingsNotificationConfigurationDescriptor =
    $convert.base64Decode(
        'CihBY2NvdW50U2V0dGluZ3NOb3RpZmljYXRpb25Db25maWd1cmF0aW9uEhwKB3JvbGVhcm4Ygf'
        'jnmQEgASgJUgdyb2xlYXJuElEKEHNuc2NvbmZpZ3VyYXRpb24YiuKhAyABKAsyIi5xdWVyeS50'
        'aW1lc3RyZWFtLlNuc0NvbmZpZ3VyYXRpb25SEHNuc2NvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use cancelQueryRequestDescriptor instead')
const CancelQueryRequest$json = {
  '1': 'CancelQueryRequest',
  '2': [
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
  ],
};

/// Descriptor for `CancelQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelQueryRequestDescriptor =
    $convert.base64Decode(
        'ChJDYW5jZWxRdWVyeVJlcXVlc3QSGwoHcXVlcnlpZBjv8OY0IAEoCVIHcXVlcnlpZA==');

@$core.Deprecated('Use cancelQueryResponseDescriptor instead')
const CancelQueryResponse$json = {
  '1': 'CancelQueryResponse',
  '2': [
    {
      '1': 'cancellationmessage',
      '3': 226613662,
      '4': 1,
      '5': 9,
      '10': 'cancellationmessage'
    },
  ],
};

/// Descriptor for `CancelQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelQueryResponseDescriptor = $convert.base64Decode(
    'ChNDYW5jZWxRdWVyeVJlc3BvbnNlEjMKE2NhbmNlbGxhdGlvbm1lc3NhZ2UYnrOHbCABKAlSE2'
    'NhbmNlbGxhdGlvbm1lc3NhZ2U=');

@$core.Deprecated('Use columnInfoDescriptor instead')
const ColumnInfo$json = {
  '1': 'ColumnInfo',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.Type',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ColumnInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List columnInfoDescriptor = $convert.base64Decode(
    'CgpDb2x1bW5JbmZvEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSLgoEdHlwZRjuoNeKASABKAsyFi'
    '5xdWVyeS50aW1lc3RyZWFtLlR5cGVSBHR5cGU=');

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

@$core.Deprecated('Use createScheduledQueryRequestDescriptor instead')
const CreateScheduledQueryRequest$json = {
  '1': 'CreateScheduledQueryRequest',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'errorreportconfiguration',
      '3': 222039776,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ErrorReportConfiguration',
      '10': 'errorreportconfiguration'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'notificationconfiguration',
      '3': 290208045,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.NotificationConfiguration',
      '10': 'notificationconfiguration'
    },
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
    {
      '1': 'scheduleconfiguration',
      '3': 526075047,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ScheduleConfiguration',
      '10': 'scheduleconfiguration'
    },
    {
      '1': 'scheduledqueryexecutionrolearn',
      '3': 284244182,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryexecutionrolearn'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.Tag',
      '10': 'tags'
    },
    {
      '1': 'targetconfiguration',
      '3': 499845483,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.TargetConfiguration',
      '10': 'targetconfiguration'
    },
  ],
};

/// Descriptor for `CreateScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createScheduledQueryRequestDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVTY2hlZHVsZWRRdWVyeVJlcXVlc3QSIwoLY2xpZW50dG9rZW4YzPu7QSABKAlSC2'
    'NsaWVudHRva2VuEmkKGGVycm9ycmVwb3J0Y29uZmlndXJhdGlvbhjgnfBpIAEoCzIqLnF1ZXJ5'
    'LnRpbWVzdHJlYW0uRXJyb3JSZXBvcnRDb25maWd1cmF0aW9uUhhlcnJvcnJlcG9ydGNvbmZpZ3'
    'VyYXRpb24SHQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEhUKBG5hbWUYh+aBfyABKAlS'
    'BG5hbWUSbQoZbm90aWZpY2F0aW9uY29uZmlndXJhdGlvbhit8rCKASABKAsyKy5xdWVyeS50aW'
    '1lc3RyZWFtLk5vdGlmaWNhdGlvbkNvbmZpZ3VyYXRpb25SGW5vdGlmaWNhdGlvbmNvbmZpZ3Vy'
    'YXRpb24SJAoLcXVlcnlzdHJpbmcY58rvzwEgASgJUgtxdWVyeXN0cmluZxJhChVzY2hlZHVsZW'
    'NvbmZpZ3VyYXRpb24Yp4nt+gEgASgLMicucXVlcnkudGltZXN0cmVhbS5TY2hlZHVsZUNvbmZp'
    'Z3VyYXRpb25SFXNjaGVkdWxlY29uZmlndXJhdGlvbhJKCh5zY2hlZHVsZWRxdWVyeWV4ZWN1dG'
    'lvbnJvbGVhcm4Y1vHEhwEgASgJUh5zY2hlZHVsZWRxdWVyeWV4ZWN1dGlvbnJvbGVhcm4SLQoE'
    'dGFncxjBwfa1ASADKAsyFS5xdWVyeS50aW1lc3RyZWFtLlRhZ1IEdGFncxJbChN0YXJnZXRjb2'
    '5maWd1cmF0aW9uGOuSrO4BIAEoCzIlLnF1ZXJ5LnRpbWVzdHJlYW0uVGFyZ2V0Q29uZmlndXJh'
    'dGlvblITdGFyZ2V0Y29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use createScheduledQueryResponseDescriptor instead')
const CreateScheduledQueryResponse$json = {
  '1': 'CreateScheduledQueryResponse',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `CreateScheduledQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createScheduledQueryResponseDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVTY2hlZHVsZWRRdWVyeVJlc3BvbnNlEhQKA2Fybhidm+2/ASABKAlSA2Fybg==');

@$core.Deprecated('Use datumDescriptor instead')
const Datum$json = {
  '1': 'Datum',
  '2': [
    {
      '1': 'arrayvalue',
      '3': 20393608,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.Datum',
      '10': 'arrayvalue'
    },
    {'1': 'nullvalue', '3': 440981694, '4': 1, '5': 8, '10': 'nullvalue'},
    {
      '1': 'rowvalue',
      '3': 530552345,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.Row',
      '10': 'rowvalue'
    },
    {'1': 'scalarvalue', '3': 45700451, '4': 1, '5': 9, '10': 'scalarvalue'},
    {
      '1': 'timeseriesvalue',
      '3': 468590991,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.TimeSeriesDataPoint',
      '10': 'timeseriesvalue'
    },
  ],
};

/// Descriptor for `Datum`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List datumDescriptor = $convert.base64Decode(
    'CgVEYXR1bRI6CgphcnJheXZhbHVlGIjd3AkgAygLMhcucXVlcnkudGltZXN0cmVhbS5EYXR1bV'
    'IKYXJyYXl2YWx1ZRIgCgludWxsdmFsdWUYvrGj0gEgASgIUgludWxsdmFsdWUSNQoIcm93dmFs'
    'dWUYmaz+/AEgASgLMhUucXVlcnkudGltZXN0cmVhbS5Sb3dSCHJvd3ZhbHVlEiMKC3NjYWxhcn'
    'ZhbHVlGOOq5RUgASgJUgtzY2FsYXJ2YWx1ZRJTCg90aW1lc2VyaWVzdmFsdWUYj8O43wEgAygL'
    'MiUucXVlcnkudGltZXN0cmVhbS5UaW1lU2VyaWVzRGF0YVBvaW50Ug90aW1lc2VyaWVzdmFsdW'
    'U=');

@$core.Deprecated('Use deleteScheduledQueryRequestDescriptor instead')
const DeleteScheduledQueryRequest$json = {
  '1': 'DeleteScheduledQueryRequest',
  '2': [
    {
      '1': 'scheduledqueryarn',
      '3': 234602964,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
  ],
};

/// Descriptor for `DeleteScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteScheduledQueryRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVTY2hlZHVsZWRRdWVyeVJlcXVlc3QSLwoRc2NoZWR1bGVkcXVlcnlhcm4Y1IPvby'
        'ABKAlSEXNjaGVkdWxlZHF1ZXJ5YXJu');

@$core.Deprecated('Use describeAccountSettingsRequestDescriptor instead')
const DescribeAccountSettingsRequest$json = {
  '1': 'DescribeAccountSettingsRequest',
};

/// Descriptor for `DescribeAccountSettingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAccountSettingsRequestDescriptor =
    $convert.base64Decode('Ch5EZXNjcmliZUFjY291bnRTZXR0aW5nc1JlcXVlc3Q=');

@$core.Deprecated('Use describeAccountSettingsResponseDescriptor instead')
const DescribeAccountSettingsResponse$json = {
  '1': 'DescribeAccountSettingsResponse',
  '2': [
    {'1': 'maxquerytcu', '3': 285295370, '4': 1, '5': 5, '10': 'maxquerytcu'},
    {
      '1': 'querycompute',
      '3': 358277613,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryComputeResponse',
      '10': 'querycompute'
    },
    {
      '1': 'querypricingmodel',
      '3': 453472445,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.QueryPricingModel',
      '10': 'querypricingmodel'
    },
  ],
};

/// Descriptor for `DescribeAccountSettingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAccountSettingsResponseDescriptor = $convert.base64Decode(
    'Ch9EZXNjcmliZUFjY291bnRTZXR0aW5nc1Jlc3BvbnNlEiQKC21heHF1ZXJ5dGN1GIqGhYgBIA'
    'EoBVILbWF4cXVlcnl0Y3USTgoMcXVlcnljb21wdXRlGO3D66oBIAEoCzImLnF1ZXJ5LnRpbWVz'
    'dHJlYW0uUXVlcnlDb21wdXRlUmVzcG9uc2VSDHF1ZXJ5Y29tcHV0ZRJVChFxdWVyeXByaWNpbm'
    'dtb2RlbBi94Z3YASABKA4yIy5xdWVyeS50aW1lc3RyZWFtLlF1ZXJ5UHJpY2luZ01vZGVsUhFx'
    'dWVyeXByaWNpbmdtb2RlbA==');

@$core.Deprecated('Use describeEndpointsRequestDescriptor instead')
const DescribeEndpointsRequest$json = {
  '1': 'DescribeEndpointsRequest',
};

/// Descriptor for `DescribeEndpointsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEndpointsRequestDescriptor =
    $convert.base64Decode('ChhEZXNjcmliZUVuZHBvaW50c1JlcXVlc3Q=');

@$core.Deprecated('Use describeEndpointsResponseDescriptor instead')
const DescribeEndpointsResponse$json = {
  '1': 'DescribeEndpointsResponse',
  '2': [
    {
      '1': 'endpoints',
      '3': 16210494,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.Endpoint',
      '10': 'endpoints'
    },
  ],
};

/// Descriptor for `DescribeEndpointsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEndpointsResponseDescriptor =
    $convert.base64Decode(
        'ChlEZXNjcmliZUVuZHBvaW50c1Jlc3BvbnNlEjsKCWVuZHBvaW50cxi+tN0HIAMoCzIaLnF1ZX'
        'J5LnRpbWVzdHJlYW0uRW5kcG9pbnRSCWVuZHBvaW50cw==');

@$core.Deprecated('Use describeScheduledQueryRequestDescriptor instead')
const DescribeScheduledQueryRequest$json = {
  '1': 'DescribeScheduledQueryRequest',
  '2': [
    {
      '1': 'scheduledqueryarn',
      '3': 234602964,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
  ],
};

/// Descriptor for `DescribeScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeScheduledQueryRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZVNjaGVkdWxlZFF1ZXJ5UmVxdWVzdBIvChFzY2hlZHVsZWRxdWVyeWFybhjUg+'
        '9vIAEoCVIRc2NoZWR1bGVkcXVlcnlhcm4=');

@$core.Deprecated('Use describeScheduledQueryResponseDescriptor instead')
const DescribeScheduledQueryResponse$json = {
  '1': 'DescribeScheduledQueryResponse',
  '2': [
    {
      '1': 'scheduledquery',
      '3': 338643709,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ScheduledQueryDescription',
      '10': 'scheduledquery'
    },
  ],
};

/// Descriptor for `DescribeScheduledQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeScheduledQueryResponseDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZVNjaGVkdWxlZFF1ZXJ5UmVzcG9uc2USVwoOc2NoZWR1bGVkcXVlcnkY/ZW9oQ'
        'EgASgLMisucXVlcnkudGltZXN0cmVhbS5TY2hlZHVsZWRRdWVyeURlc2NyaXB0aW9uUg5zY2hl'
        'ZHVsZWRxdWVyeQ==');

@$core.Deprecated('Use dimensionMappingDescriptor instead')
const DimensionMapping$json = {
  '1': 'DimensionMapping',
  '2': [
    {
      '1': 'dimensionvaluetype',
      '3': 267417961,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.DimensionValueType',
      '10': 'dimensionvaluetype'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DimensionMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dimensionMappingDescriptor = $convert.base64Decode(
    'ChBEaW1lbnNpb25NYXBwaW5nElcKEmRpbWVuc2lvbnZhbHVldHlwZRjp8sF/IAEoDjIkLnF1ZX'
    'J5LnRpbWVzdHJlYW0uRGltZW5zaW9uVmFsdWVUeXBlUhJkaW1lbnNpb252YWx1ZXR5cGUSFQoE'
    'bmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use endpointDescriptor instead')
const Endpoint$json = {
  '1': 'Endpoint',
  '2': [
    {'1': 'address', '3': 268787956, '4': 1, '5': 9, '10': 'address'},
    {
      '1': 'cacheperiodinminutes',
      '3': 71526449,
      '4': 1,
      '5': 3,
      '10': 'cacheperiodinminutes'
    },
  ],
};

/// Descriptor for `Endpoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endpointDescriptor = $convert.base64Decode(
    'CghFbmRwb2ludBIcCgdhZGRyZXNzGPTBlYABIAEoCVIHYWRkcmVzcxI1ChRjYWNoZXBlcmlvZG'
    'lubWludXRlcxix0I0iIAEoA1IUY2FjaGVwZXJpb2Rpbm1pbnV0ZXM=');

@$core.Deprecated('Use errorReportConfigurationDescriptor instead')
const ErrorReportConfiguration$json = {
  '1': 'ErrorReportConfiguration',
  '2': [
    {
      '1': 's3configuration',
      '3': 27828476,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.S3Configuration',
      '10': 's3configuration'
    },
  ],
};

/// Descriptor for `ErrorReportConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorReportConfigurationDescriptor =
    $convert.base64Decode(
        'ChhFcnJvclJlcG9ydENvbmZpZ3VyYXRpb24STgoPczNjb25maWd1cmF0aW9uGPzBog0gASgLMi'
        'EucXVlcnkudGltZXN0cmVhbS5TM0NvbmZpZ3VyYXRpb25SD3MzY29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use errorReportLocationDescriptor instead')
const ErrorReportLocation$json = {
  '1': 'ErrorReportLocation',
  '2': [
    {
      '1': 's3reportlocation',
      '3': 162597615,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.S3ReportLocation',
      '10': 's3reportlocation'
    },
  ],
};

/// Descriptor for `ErrorReportLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorReportLocationDescriptor = $convert.base64Decode(
    'ChNFcnJvclJlcG9ydExvY2F0aW9uElEKEHMzcmVwb3J0bG9jYXRpb24Y75XETSABKAsyIi5xdW'
    'VyeS50aW1lc3RyZWFtLlMzUmVwb3J0TG9jYXRpb25SEHMzcmVwb3J0bG9jYXRpb24=');

@$core.Deprecated('Use executeScheduledQueryRequestDescriptor instead')
const ExecuteScheduledQueryRequest$json = {
  '1': 'ExecuteScheduledQueryRequest',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'invocationtime',
      '3': 331845291,
      '4': 1,
      '5': 9,
      '10': 'invocationtime'
    },
    {
      '1': 'queryinsights',
      '3': 105458863,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ScheduledQueryInsights',
      '10': 'queryinsights'
    },
    {
      '1': 'scheduledqueryarn',
      '3': 234602964,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
  ],
};

/// Descriptor for `ExecuteScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executeScheduledQueryRequestDescriptor = $convert.base64Decode(
    'ChxFeGVjdXRlU2NoZWR1bGVkUXVlcnlSZXF1ZXN0EiMKC2NsaWVudHRva2VuGMz7u0EgASgJUg'
    'tjbGllbnR0b2tlbhIqCg5pbnZvY2F0aW9udGltZRirnZ6eASABKAlSDmludm9jYXRpb250aW1l'
    'ElEKDXF1ZXJ5aW5zaWdodHMYr9mkMiABKAsyKC5xdWVyeS50aW1lc3RyZWFtLlNjaGVkdWxlZF'
    'F1ZXJ5SW5zaWdodHNSDXF1ZXJ5aW5zaWdodHMSLwoRc2NoZWR1bGVkcXVlcnlhcm4Y1IPvbyAB'
    'KAlSEXNjaGVkdWxlZHF1ZXJ5YXJu');

@$core.Deprecated('Use executionStatsDescriptor instead')
const ExecutionStats$json = {
  '1': 'ExecutionStats',
  '2': [
    {'1': 'bytesmetered', '3': 232712245, '4': 1, '5': 3, '10': 'bytesmetered'},
    {
      '1': 'cumulativebytesscanned',
      '3': 413524176,
      '4': 1,
      '5': 3,
      '10': 'cumulativebytesscanned'
    },
    {'1': 'datawrites', '3': 406476830, '4': 1, '5': 3, '10': 'datawrites'},
    {
      '1': 'executiontimeinmillis',
      '3': 447596178,
      '4': 1,
      '5': 3,
      '10': 'executiontimeinmillis'
    },
    {
      '1': 'queryresultrows',
      '3': 240264704,
      '4': 1,
      '5': 3,
      '10': 'queryresultrows'
    },
    {
      '1': 'recordsingested',
      '3': 159925469,
      '4': 1,
      '5': 3,
      '10': 'recordsingested'
    },
  ],
};

/// Descriptor for `ExecutionStats`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionStatsDescriptor = $convert.base64Decode(
    'Cg5FeGVjdXRpb25TdGF0cxIlCgxieXRlc21ldGVyZWQYtdD7biABKANSDGJ5dGVzbWV0ZXJlZB'
    'I6ChZjdW11bGF0aXZlYnl0ZXNzY2FubmVkGNDBl8UBIAEoA1IWY3VtdWxhdGl2ZWJ5dGVzc2Nh'
    'bm5lZBIiCgpkYXRhd3JpdGVzGJ6w6cEBIAEoA1IKZGF0YXdyaXRlcxI4ChVleGVjdXRpb250aW'
    '1laW5taWxsaXMYko231QEgASgDUhVleGVjdXRpb250aW1laW5taWxsaXMSKwoPcXVlcnlyZXN1'
    'bHRyb3dzGIDMyHIgASgDUg9xdWVyeXJlc3VsdHJvd3MSKwoPcmVjb3Jkc2luZ2VzdGVkGN2JoU'
    'wgASgDUg9yZWNvcmRzaW5nZXN0ZWQ=');

@$core.Deprecated('Use internalServerExceptionDescriptor instead')
const InternalServerException$json = {
  '1': 'InternalServerException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalServerException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalServerExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbnRlcm5hbFNlcnZlckV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidEndpointExceptionDescriptor instead')
const InvalidEndpointException$json = {
  '1': 'InvalidEndpointException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidEndpointException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidEndpointExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkRW5kcG9pbnRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use lastUpdateDescriptor instead')
const LastUpdate$json = {
  '1': 'LastUpdate',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.LastUpdateStatus',
      '10': 'status'
    },
    {
      '1': 'statusmessage',
      '3': 72590095,
      '4': 1,
      '5': 9,
      '10': 'statusmessage'
    },
    {
      '1': 'targetquerytcu',
      '3': 183880621,
      '4': 1,
      '5': 5,
      '10': 'targetquerytcu'
    },
  ],
};

/// Descriptor for `LastUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lastUpdateDescriptor = $convert.base64Decode(
    'CgpMYXN0VXBkYXRlEj0KBnN0YXR1cxiQ5PsCIAEoDjIiLnF1ZXJ5LnRpbWVzdHJlYW0uTGFzdF'
    'VwZGF0ZVN0YXR1c1IGc3RhdHVzEicKDXN0YXR1c21lc3NhZ2UYj8bOIiABKAlSDXN0YXR1c21l'
    'c3NhZ2USKQoOdGFyZ2V0cXVlcnl0Y3UYrZfXVyABKAVSDnRhcmdldHF1ZXJ5dGN1');

@$core.Deprecated('Use listScheduledQueriesRequestDescriptor instead')
const ListScheduledQueriesRequest$json = {
  '1': 'ListScheduledQueriesRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListScheduledQueriesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listScheduledQueriesRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0U2NoZWR1bGVkUXVlcmllc1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm'
        '1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listScheduledQueriesResponseDescriptor instead')
const ListScheduledQueriesResponse$json = {
  '1': 'ListScheduledQueriesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'scheduledqueries',
      '3': 458789865,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.ScheduledQuery',
      '10': 'scheduledqueries'
    },
  ],
};

/// Descriptor for `ListScheduledQueriesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listScheduledQueriesResponseDescriptor =
    $convert.base64Decode(
        'ChxMaXN0U2NoZWR1bGVkUXVlcmllc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbm'
        'V4dHRva2VuElAKEHNjaGVkdWxlZHF1ZXJpZXMY6afi2gEgAygLMiAucXVlcnkudGltZXN0cmVh'
        'bS5TY2hlZHVsZWRRdWVyeVIQc2NoZWR1bGVkcXVlcmllcw==');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
        'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIkCgtyZXNvdXJjZWFy'
        'bhjtwJmwASABKAlSC3Jlc291cmNlYXJu');

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse$json = {
  '1': 'ListTagsForResourceResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4SLQoEdGFncxjBwfa1ASADKAsyFS5xdWVyeS50aW1lc3RyZWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use mixedMeasureMappingDescriptor instead')
const MixedMeasureMapping$json = {
  '1': 'MixedMeasureMapping',
  '2': [
    {'1': 'measurename', '3': 426079069, '4': 1, '5': 9, '10': 'measurename'},
    {
      '1': 'measurevaluetype',
      '3': 466683165,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.MeasureValueType',
      '10': 'measurevaluetype'
    },
    {
      '1': 'multimeasureattributemappings',
      '3': 311133918,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.MultiMeasureAttributeMapping',
      '10': 'multimeasureattributemappings'
    },
    {'1': 'sourcecolumn', '3': 219947651, '4': 1, '5': 9, '10': 'sourcecolumn'},
    {
      '1': 'targetmeasurename',
      '3': 469508316,
      '4': 1,
      '5': 9,
      '10': 'targetmeasurename'
    },
  ],
};

/// Descriptor for `MixedMeasureMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mixedMeasureMappingDescriptor = $convert.base64Decode(
    'ChNNaXhlZE1lYXN1cmVNYXBwaW5nEiQKC21lYXN1cmVuYW1lGN3mlcsBIAEoCVILbWVhc3VyZW'
    '5hbWUSUgoQbWVhc3VyZXZhbHVldHlwZRidisTeASABKA4yIi5xdWVyeS50aW1lc3RyZWFtLk1l'
    'YXN1cmVWYWx1ZVR5cGVSEG1lYXN1cmV2YWx1ZXR5cGUSeAodbXVsdGltZWFzdXJlYXR0cmlidX'
    'RlbWFwcGluZ3MY3o2ulAEgAygLMi4ucXVlcnkudGltZXN0cmVhbS5NdWx0aU1lYXN1cmVBdHRy'
    'aWJ1dGVNYXBwaW5nUh1tdWx0aW1lYXN1cmVhdHRyaWJ1dGVtYXBwaW5ncxIlCgxzb3VyY2Vjb2'
    'x1bW4Yg8XwaCABKAlSDHNvdXJjZWNvbHVtbhIwChF0YXJnZXRtZWFzdXJlbmFtZRjcwfDfASAB'
    'KAlSEXRhcmdldG1lYXN1cmVuYW1l');

@$core.Deprecated('Use multiMeasureAttributeMappingDescriptor instead')
const MultiMeasureAttributeMapping$json = {
  '1': 'MultiMeasureAttributeMapping',
  '2': [
    {
      '1': 'measurevaluetype',
      '3': 466683165,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ScalarMeasureValueType',
      '10': 'measurevaluetype'
    },
    {'1': 'sourcecolumn', '3': 219947651, '4': 1, '5': 9, '10': 'sourcecolumn'},
    {
      '1': 'targetmultimeasureattributename',
      '3': 415623663,
      '4': 1,
      '5': 9,
      '10': 'targetmultimeasureattributename'
    },
  ],
};

/// Descriptor for `MultiMeasureAttributeMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List multiMeasureAttributeMappingDescriptor = $convert.base64Decode(
    'ChxNdWx0aU1lYXN1cmVBdHRyaWJ1dGVNYXBwaW5nElgKEG1lYXN1cmV2YWx1ZXR5cGUYnYrE3g'
    'EgASgOMigucXVlcnkudGltZXN0cmVhbS5TY2FsYXJNZWFzdXJlVmFsdWVUeXBlUhBtZWFzdXJl'
    'dmFsdWV0eXBlEiUKDHNvdXJjZWNvbHVtbhiDxfBoIAEoCVIMc291cmNlY29sdW1uEkwKH3Rhcm'
    'dldG11bHRpbWVhc3VyZWF0dHJpYnV0ZW5hbWUY79OXxgEgASgJUh90YXJnZXRtdWx0aW1lYXN1'
    'cmVhdHRyaWJ1dGVuYW1l');

@$core.Deprecated('Use multiMeasureMappingsDescriptor instead')
const MultiMeasureMappings$json = {
  '1': 'MultiMeasureMappings',
  '2': [
    {
      '1': 'multimeasureattributemappings',
      '3': 311133918,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.MultiMeasureAttributeMapping',
      '10': 'multimeasureattributemappings'
    },
    {
      '1': 'targetmultimeasurename',
      '3': 27420053,
      '4': 1,
      '5': 9,
      '10': 'targetmultimeasurename'
    },
  ],
};

/// Descriptor for `MultiMeasureMappings`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List multiMeasureMappingsDescriptor = $convert.base64Decode(
    'ChRNdWx0aU1lYXN1cmVNYXBwaW5ncxJ4Ch1tdWx0aW1lYXN1cmVhdHRyaWJ1dGVtYXBwaW5ncx'
    'jeja6UASADKAsyLi5xdWVyeS50aW1lc3RyZWFtLk11bHRpTWVhc3VyZUF0dHJpYnV0ZU1hcHBp'
    'bmdSHW11bHRpbWVhc3VyZWF0dHJpYnV0ZW1hcHBpbmdzEjkKFnRhcmdldG11bHRpbWVhc3VyZW'
    '5hbWUYlcuJDSABKAlSFnRhcmdldG11bHRpbWVhc3VyZW5hbWU=');

@$core.Deprecated('Use notificationConfigurationDescriptor instead')
const NotificationConfiguration$json = {
  '1': 'NotificationConfiguration',
  '2': [
    {
      '1': 'snsconfiguration',
      '3': 6844682,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.SnsConfiguration',
      '10': 'snsconfiguration'
    },
  ],
};

/// Descriptor for `NotificationConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notificationConfigurationDescriptor = $convert.base64Decode(
    'ChlOb3RpZmljYXRpb25Db25maWd1cmF0aW9uElEKEHNuc2NvbmZpZ3VyYXRpb24YiuKhAyABKA'
    'syIi5xdWVyeS50aW1lc3RyZWFtLlNuc0NvbmZpZ3VyYXRpb25SEHNuc2NvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use parameterMappingDescriptor instead')
const ParameterMapping$json = {
  '1': 'ParameterMapping',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.Type',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ParameterMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterMappingDescriptor = $convert.base64Decode(
    'ChBQYXJhbWV0ZXJNYXBwaW5nEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSLgoEdHlwZRjuoNeKAS'
    'ABKAsyFi5xdWVyeS50aW1lc3RyZWFtLlR5cGVSBHR5cGU=');

@$core.Deprecated('Use prepareQueryRequestDescriptor instead')
const PrepareQueryRequest$json = {
  '1': 'PrepareQueryRequest',
  '2': [
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
    {'1': 'validateonly', '3': 453040966, '4': 1, '5': 8, '10': 'validateonly'},
  ],
};

/// Descriptor for `PrepareQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List prepareQueryRequestDescriptor = $convert.base64Decode(
    'ChNQcmVwYXJlUXVlcnlSZXF1ZXN0EiQKC3F1ZXJ5c3RyaW5nGOfK788BIAEoCVILcXVlcnlzdH'
    'JpbmcSJgoMdmFsaWRhdGVvbmx5GMa2g9gBIAEoCFIMdmFsaWRhdGVvbmx5');

@$core.Deprecated('Use prepareQueryResponseDescriptor instead')
const PrepareQueryResponse$json = {
  '1': 'PrepareQueryResponse',
  '2': [
    {
      '1': 'columns',
      '3': 169177053,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.SelectColumn',
      '10': 'columns'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.ParameterMapping',
      '10': 'parameters'
    },
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
  ],
};

/// Descriptor for `PrepareQueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List prepareQueryResponseDescriptor = $convert.base64Decode(
    'ChRQcmVwYXJlUXVlcnlSZXNwb25zZRI7Cgdjb2x1bW5zGN3f1VAgAygLMh4ucXVlcnkudGltZX'
    'N0cmVhbS5TZWxlY3RDb2x1bW5SB2NvbHVtbnMSRgoKcGFyYW1ldGVycxj6p/7rASADKAsyIi5x'
    'dWVyeS50aW1lc3RyZWFtLlBhcmFtZXRlck1hcHBpbmdSCnBhcmFtZXRlcnMSJAoLcXVlcnlzdH'
    'JpbmcY58rvzwEgASgJUgtxdWVyeXN0cmluZw==');

@$core.Deprecated('Use provisionedCapacityRequestDescriptor instead')
const ProvisionedCapacityRequest$json = {
  '1': 'ProvisionedCapacityRequest',
  '2': [
    {
      '1': 'notificationconfiguration',
      '3': 290208045,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.AccountSettingsNotificationConfiguration',
      '10': 'notificationconfiguration'
    },
    {
      '1': 'targetquerytcu',
      '3': 183880621,
      '4': 1,
      '5': 5,
      '10': 'targetquerytcu'
    },
  ],
};

/// Descriptor for `ProvisionedCapacityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedCapacityRequestDescriptor = $convert.base64Decode(
    'ChpQcm92aXNpb25lZENhcGFjaXR5UmVxdWVzdBJ8Chlub3RpZmljYXRpb25jb25maWd1cmF0aW'
    '9uGK3ysIoBIAEoCzI6LnF1ZXJ5LnRpbWVzdHJlYW0uQWNjb3VudFNldHRpbmdzTm90aWZpY2F0'
    'aW9uQ29uZmlndXJhdGlvblIZbm90aWZpY2F0aW9uY29uZmlndXJhdGlvbhIpCg50YXJnZXRxdW'
    'VyeXRjdRitl9dXIAEoBVIOdGFyZ2V0cXVlcnl0Y3U=');

@$core.Deprecated('Use provisionedCapacityResponseDescriptor instead')
const ProvisionedCapacityResponse$json = {
  '1': 'ProvisionedCapacityResponse',
  '2': [
    {
      '1': 'activequerytcu',
      '3': 136246688,
      '4': 1,
      '5': 5,
      '10': 'activequerytcu'
    },
    {
      '1': 'lastupdate',
      '3': 331125817,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.LastUpdate',
      '10': 'lastupdate'
    },
    {
      '1': 'notificationconfiguration',
      '3': 290208045,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.AccountSettingsNotificationConfiguration',
      '10': 'notificationconfiguration'
    },
  ],
};

/// Descriptor for `ProvisionedCapacityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedCapacityResponseDescriptor = $convert.base64Decode(
    'ChtQcm92aXNpb25lZENhcGFjaXR5UmVzcG9uc2USKQoOYWN0aXZlcXVlcnl0Y3UYoOv7QCABKA'
    'VSDmFjdGl2ZXF1ZXJ5dGN1EkAKCmxhc3R1cGRhdGUYuajynQEgASgLMhwucXVlcnkudGltZXN0'
    'cmVhbS5MYXN0VXBkYXRlUgpsYXN0dXBkYXRlEnwKGW5vdGlmaWNhdGlvbmNvbmZpZ3VyYXRpb2'
    '4YrfKwigEgASgLMjoucXVlcnkudGltZXN0cmVhbS5BY2NvdW50U2V0dGluZ3NOb3RpZmljYXRp'
    'b25Db25maWd1cmF0aW9uUhlub3RpZmljYXRpb25jb25maWd1cmF0aW9u');

@$core.Deprecated('Use queryComputeRequestDescriptor instead')
const QueryComputeRequest$json = {
  '1': 'QueryComputeRequest',
  '2': [
    {
      '1': 'computemode',
      '3': 205227458,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ComputeMode',
      '10': 'computemode'
    },
    {
      '1': 'provisionedcapacity',
      '3': 181400702,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ProvisionedCapacityRequest',
      '10': 'provisionedcapacity'
    },
  ],
};

/// Descriptor for `QueryComputeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryComputeRequestDescriptor = $convert.base64Decode(
    'ChNRdWVyeUNvbXB1dGVSZXF1ZXN0EkIKC2NvbXB1dGVtb2RlGMKL7mEgASgOMh0ucXVlcnkudG'
    'ltZXN0cmVhbS5Db21wdXRlTW9kZVILY29tcHV0ZW1vZGUSYQoTcHJvdmlzaW9uZWRjYXBhY2l0'
    'eRj+6L9WIAEoCzIsLnF1ZXJ5LnRpbWVzdHJlYW0uUHJvdmlzaW9uZWRDYXBhY2l0eVJlcXVlc3'
    'RSE3Byb3Zpc2lvbmVkY2FwYWNpdHk=');

@$core.Deprecated('Use queryComputeResponseDescriptor instead')
const QueryComputeResponse$json = {
  '1': 'QueryComputeResponse',
  '2': [
    {
      '1': 'computemode',
      '3': 205227458,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ComputeMode',
      '10': 'computemode'
    },
    {
      '1': 'provisionedcapacity',
      '3': 181400702,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ProvisionedCapacityResponse',
      '10': 'provisionedcapacity'
    },
  ],
};

/// Descriptor for `QueryComputeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryComputeResponseDescriptor = $convert.base64Decode(
    'ChRRdWVyeUNvbXB1dGVSZXNwb25zZRJCCgtjb21wdXRlbW9kZRjCi+5hIAEoDjIdLnF1ZXJ5Ln'
    'RpbWVzdHJlYW0uQ29tcHV0ZU1vZGVSC2NvbXB1dGVtb2RlEmIKE3Byb3Zpc2lvbmVkY2FwYWNp'
    'dHkY/ui/ViABKAsyLS5xdWVyeS50aW1lc3RyZWFtLlByb3Zpc2lvbmVkQ2FwYWNpdHlSZXNwb2'
    '5zZVITcHJvdmlzaW9uZWRjYXBhY2l0eQ==');

@$core.Deprecated('Use queryExecutionExceptionDescriptor instead')
const QueryExecutionException$json = {
  '1': 'QueryExecutionException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `QueryExecutionException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryExecutionExceptionDescriptor =
    $convert.base64Decode(
        'ChdRdWVyeUV4ZWN1dGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use queryInsightsDescriptor instead')
const QueryInsights$json = {
  '1': 'QueryInsights',
  '2': [
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.QueryInsightsMode',
      '10': 'mode'
    },
  ],
};

/// Descriptor for `QueryInsights`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryInsightsDescriptor = $convert.base64Decode(
    'Cg1RdWVyeUluc2lnaHRzEjsKBG1vZGUYs+65mgEgASgOMiMucXVlcnkudGltZXN0cmVhbS5RdW'
    'VyeUluc2lnaHRzTW9kZVIEbW9kZQ==');

@$core.Deprecated('Use queryInsightsResponseDescriptor instead')
const QueryInsightsResponse$json = {
  '1': 'QueryInsightsResponse',
  '2': [
    {'1': 'outputbytes', '3': 318971400, '4': 1, '5': 3, '10': 'outputbytes'},
    {'1': 'outputrows', '3': 138873322, '4': 1, '5': 3, '10': 'outputrows'},
    {
      '1': 'queryspatialcoverage',
      '3': 200825688,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QuerySpatialCoverage',
      '10': 'queryspatialcoverage'
    },
    {
      '1': 'querytablecount',
      '3': 28893653,
      '4': 1,
      '5': 3,
      '10': 'querytablecount'
    },
    {
      '1': 'querytemporalrange',
      '3': 200957543,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryTemporalRange',
      '10': 'querytemporalrange'
    },
    {
      '1': 'unloadpartitioncount',
      '3': 467201664,
      '4': 1,
      '5': 3,
      '10': 'unloadpartitioncount'
    },
    {
      '1': 'unloadwrittenbytes',
      '3': 247315473,
      '4': 1,
      '5': 3,
      '10': 'unloadwrittenbytes'
    },
    {
      '1': 'unloadwrittenrows',
      '3': 20403789,
      '4': 1,
      '5': 3,
      '10': 'unloadwrittenrows'
    },
  ],
};

/// Descriptor for `QueryInsightsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryInsightsResponseDescriptor = $convert.base64Decode(
    'ChVRdWVyeUluc2lnaHRzUmVzcG9uc2USJAoLb3V0cHV0Ynl0ZXMYiLyMmAEgASgDUgtvdXRwdX'
    'RieXRlcxIhCgpvdXRwdXRyb3dzGOqTnEIgASgDUgpvdXRwdXRyb3dzEl0KFHF1ZXJ5c3BhdGlh'
    'bGNvdmVyYWdlGNi24V8gASgLMiYucXVlcnkudGltZXN0cmVhbS5RdWVyeVNwYXRpYWxDb3Zlcm'
    'FnZVIUcXVlcnlzcGF0aWFsY292ZXJhZ2USKwoPcXVlcnl0YWJsZWNvdW50GNXD4w0gASgDUg9x'
    'dWVyeXRhYmxlY291bnQSVwoScXVlcnl0ZW1wb3JhbHJhbmdlGOe86V8gASgLMiQucXVlcnkudG'
    'ltZXN0cmVhbS5RdWVyeVRlbXBvcmFsUmFuZ2VSEnF1ZXJ5dGVtcG9yYWxyYW5nZRI2ChR1bmxv'
    'YWRwYXJ0aXRpb25jb3VudBiA3ePeASABKANSFHVubG9hZHBhcnRpdGlvbmNvdW50EjEKEnVubG'
    '9hZHdyaXR0ZW5ieXRlcxiR+PZ1IAEoA1ISdW5sb2Fkd3JpdHRlbmJ5dGVzEi8KEXVubG9hZHdy'
    'aXR0ZW5yb3dzGM2s3QkgASgDUhF1bmxvYWR3cml0dGVucm93cw==');

@$core.Deprecated('Use queryRequestDescriptor instead')
const QueryRequest$json = {
  '1': 'QueryRequest',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'maxrows', '3': 251920525, '4': 1, '5': 5, '10': 'maxrows'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'queryinsights',
      '3': 105458863,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryInsights',
      '10': 'queryinsights'
    },
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
  ],
};

/// Descriptor for `QueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryRequestDescriptor = $convert.base64Decode(
    'CgxRdWVyeVJlcXVlc3QSIwoLY2xpZW50dG9rZW4YzPu7QSABKAlSC2NsaWVudHRva2VuEhsKB2'
    '1heHJvd3MYjYGQeCABKAVSB21heHJvd3MSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9r'
    'ZW4SSAoNcXVlcnlpbnNpZ2h0cxiv2aQyIAEoCzIfLnF1ZXJ5LnRpbWVzdHJlYW0uUXVlcnlJbn'
    'NpZ2h0c1INcXVlcnlpbnNpZ2h0cxIkCgtxdWVyeXN0cmluZxjnyu/PASABKAlSC3F1ZXJ5c3Ry'
    'aW5n');

@$core.Deprecated('Use queryResponseDescriptor instead')
const QueryResponse$json = {
  '1': 'QueryResponse',
  '2': [
    {
      '1': 'columninfo',
      '3': 364742404,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.ColumnInfo',
      '10': 'columninfo'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'queryid', '3': 110737519, '4': 1, '5': 9, '10': 'queryid'},
    {
      '1': 'queryinsightsresponse',
      '3': 354278130,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryInsightsResponse',
      '10': 'queryinsightsresponse'
    },
    {
      '1': 'querystatus',
      '3': 367016406,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryStatus',
      '10': 'querystatus'
    },
    {
      '1': 'rows',
      '3': 174428857,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.Row',
      '10': 'rows'
    },
  ],
};

/// Descriptor for `QueryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryResponseDescriptor = $convert.base64Decode(
    'Cg1RdWVyeVJlc3BvbnNlEkAKCmNvbHVtbmluZm8YhI72rQEgAygLMhwucXVlcnkudGltZXN0cm'
    'VhbS5Db2x1bW5JbmZvUgpjb2x1bW5pbmZvEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRv'
    'a2VuEhsKB3F1ZXJ5aWQY7/DmNCABKAlSB3F1ZXJ5aWQSYQoVcXVlcnlpbnNpZ2h0c3Jlc3Bvbn'
    'NlGPK196gBIAEoCzInLnF1ZXJ5LnRpbWVzdHJlYW0uUXVlcnlJbnNpZ2h0c1Jlc3BvbnNlUhVx'
    'dWVyeWluc2lnaHRzcmVzcG9uc2USQwoLcXVlcnlzdGF0dXMY1vOArwEgASgLMh0ucXVlcnkudG'
    'ltZXN0cmVhbS5RdWVyeVN0YXR1c1ILcXVlcnlzdGF0dXMSLAoEcm93cxi5pZZTIAMoCzIVLnF1'
    'ZXJ5LnRpbWVzdHJlYW0uUm93UgRyb3dz');

@$core.Deprecated('Use querySpatialCoverageDescriptor instead')
const QuerySpatialCoverage$json = {
  '1': 'QuerySpatialCoverage',
  '2': [
    {
      '1': 'max',
      '3': 480764858,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QuerySpatialCoverageMax',
      '10': 'max'
    },
  ],
};

/// Descriptor for `QuerySpatialCoverage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List querySpatialCoverageDescriptor = $convert.base64Decode(
    'ChRRdWVyeVNwYXRpYWxDb3ZlcmFnZRI/CgNtYXgYusef5QEgASgLMikucXVlcnkudGltZXN0cm'
    'VhbS5RdWVyeVNwYXRpYWxDb3ZlcmFnZU1heFIDbWF4');

@$core.Deprecated('Use querySpatialCoverageMaxDescriptor instead')
const QuerySpatialCoverageMax$json = {
  '1': 'QuerySpatialCoverageMax',
  '2': [
    {'1': 'partitionkey', '3': 379379617, '4': 3, '5': 9, '10': 'partitionkey'},
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 1, '10': 'value'},
  ],
};

/// Descriptor for `QuerySpatialCoverageMax`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List querySpatialCoverageMaxDescriptor = $convert.base64Decode(
    'ChdRdWVyeVNwYXRpYWxDb3ZlcmFnZU1heBImCgxwYXJ0aXRpb25rZXkYob/ztAEgAygJUgxwYX'
    'J0aXRpb25rZXkSHgoIdGFibGVhcm4Y44DrzQEgASgJUgh0YWJsZWFybhIYCgV2YWx1ZRjr8p+K'
    'ASABKAFSBXZhbHVl');

@$core.Deprecated('Use queryStatusDescriptor instead')
const QueryStatus$json = {
  '1': 'QueryStatus',
  '2': [
    {
      '1': 'cumulativebytesmetered',
      '3': 196882576,
      '4': 1,
      '5': 3,
      '10': 'cumulativebytesmetered'
    },
    {
      '1': 'cumulativebytesscanned',
      '3': 413524176,
      '4': 1,
      '5': 3,
      '10': 'cumulativebytesscanned'
    },
    {
      '1': 'progresspercentage',
      '3': 211894727,
      '4': 1,
      '5': 1,
      '10': 'progresspercentage'
    },
  ],
};

/// Descriptor for `QueryStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryStatusDescriptor = $convert.base64Decode(
    'CgtRdWVyeVN0YXR1cxI5ChZjdW11bGF0aXZlYnl0ZXNtZXRlcmVkGJDh8F0gASgDUhZjdW11bG'
    'F0aXZlYnl0ZXNtZXRlcmVkEjoKFmN1bXVsYXRpdmVieXRlc3NjYW5uZWQY0MGXxQEgASgDUhZj'
    'dW11bGF0aXZlYnl0ZXNzY2FubmVkEjEKEnByb2dyZXNzcGVyY2VudGFnZRjHg4VlIAEoAVIScH'
    'JvZ3Jlc3NwZXJjZW50YWdl');

@$core.Deprecated('Use queryTemporalRangeDescriptor instead')
const QueryTemporalRange$json = {
  '1': 'QueryTemporalRange',
  '2': [
    {
      '1': 'max',
      '3': 480764858,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryTemporalRangeMax',
      '10': 'max'
    },
  ],
};

/// Descriptor for `QueryTemporalRange`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTemporalRangeDescriptor = $convert.base64Decode(
    'ChJRdWVyeVRlbXBvcmFsUmFuZ2USPQoDbWF4GLrHn+UBIAEoCzInLnF1ZXJ5LnRpbWVzdHJlYW'
    '0uUXVlcnlUZW1wb3JhbFJhbmdlTWF4UgNtYXg=');

@$core.Deprecated('Use queryTemporalRangeMaxDescriptor instead')
const QueryTemporalRangeMax$json = {
  '1': 'QueryTemporalRangeMax',
  '2': [
    {'1': 'tablearn', '3': 431669347, '4': 1, '5': 9, '10': 'tablearn'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 3, '10': 'value'},
  ],
};

/// Descriptor for `QueryTemporalRangeMax`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTemporalRangeMaxDescriptor = $convert.base64Decode(
    'ChVRdWVyeVRlbXBvcmFsUmFuZ2VNYXgSHgoIdGFibGVhcm4Y44DrzQEgASgJUgh0YWJsZWFybh'
    'IYCgV2YWx1ZRjr8p+KASABKANSBXZhbHVl');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'scheduledqueryarn',
      '3': 234602964,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'USLwoRc2NoZWR1bGVkcXVlcnlhcm4Y1IPvbyABKAlSEXNjaGVkdWxlZHF1ZXJ5YXJu');

@$core.Deprecated('Use rowDescriptor instead')
const Row$json = {
  '1': 'Row',
  '2': [
    {
      '1': 'data',
      '3': 525498822,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.Datum',
      '10': 'data'
    },
  ],
};

/// Descriptor for `Row`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rowDescriptor = $convert.base64Decode(
    'CgNSb3cSLwoEZGF0YRjG88n6ASADKAsyFy5xdWVyeS50aW1lc3RyZWFtLkRhdHVtUgRkYXRh');

@$core.Deprecated('Use s3ConfigurationDescriptor instead')
const S3Configuration$json = {
  '1': 'S3Configuration',
  '2': [
    {'1': 'bucketname', '3': 208117045, '4': 1, '5': 9, '10': 'bucketname'},
    {
      '1': 'encryptionoption',
      '3': 160833062,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.S3EncryptionOption',
      '10': 'encryptionoption'
    },
    {
      '1': 'objectkeyprefix',
      '3': 132617574,
      '4': 1,
      '5': 9,
      '10': 'objectkeyprefix'
    },
  ],
};

/// Descriptor for `S3Configuration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3ConfigurationDescriptor = $convert.base64Decode(
    'Cg9TM0NvbmZpZ3VyYXRpb24SIQoKYnVja2V0bmFtZRi1up5jIAEoCVIKYnVja2V0bmFtZRJTCh'
    'BlbmNyeXB0aW9ub3B0aW9uGKa82EwgASgOMiQucXVlcnkudGltZXN0cmVhbS5TM0VuY3J5cHRp'
    'b25PcHRpb25SEGVuY3J5cHRpb25vcHRpb24SKwoPb2JqZWN0a2V5cHJlZml4GOaqnj8gASgJUg'
    '9vYmplY3RrZXlwcmVmaXg=');

@$core.Deprecated('Use s3ReportLocationDescriptor instead')
const S3ReportLocation$json = {
  '1': 'S3ReportLocation',
  '2': [
    {'1': 'bucketname', '3': 208117045, '4': 1, '5': 9, '10': 'bucketname'},
    {'1': 'objectkey', '3': 335986226, '4': 1, '5': 9, '10': 'objectkey'},
  ],
};

/// Descriptor for `S3ReportLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3ReportLocationDescriptor = $convert.base64Decode(
    'ChBTM1JlcG9ydExvY2F0aW9uEiEKCmJ1Y2tldG5hbWUYtbqeYyABKAlSCmJ1Y2tldG5hbWUSIA'
    'oJb2JqZWN0a2V5GLL8mqABIAEoCVIJb2JqZWN0a2V5');

@$core.Deprecated('Use scheduleConfigurationDescriptor instead')
const ScheduleConfiguration$json = {
  '1': 'ScheduleConfiguration',
  '2': [
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
  ],
};

/// Descriptor for `ScheduleConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduleConfigurationDescriptor = $convert.base64Decode(
    'ChVTY2hlZHVsZUNvbmZpZ3VyYXRpb24SMgoSc2NoZWR1bGVleHByZXNzaW9uGP+R29QBIAEoCV'
    'ISc2NoZWR1bGVleHByZXNzaW9u');

@$core.Deprecated('Use scheduledQueryDescriptor instead')
const ScheduledQuery$json = {
  '1': 'ScheduledQuery',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'errorreportconfiguration',
      '3': 222039776,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ErrorReportConfiguration',
      '10': 'errorreportconfiguration'
    },
    {
      '1': 'lastrunstatus',
      '3': 441976361,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ScheduledQueryRunStatus',
      '10': 'lastrunstatus'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nextinvocationtime',
      '3': 424223272,
      '4': 1,
      '5': 9,
      '10': 'nextinvocationtime'
    },
    {
      '1': 'previousinvocationtime',
      '3': 7530344,
      '4': 1,
      '5': 9,
      '10': 'previousinvocationtime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ScheduledQueryState',
      '10': 'state'
    },
    {
      '1': 'targetdestination',
      '3': 121667401,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.TargetDestination',
      '10': 'targetdestination'
    },
  ],
};

/// Descriptor for `ScheduledQuery`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduledQueryDescriptor = $convert.base64Decode(
    'Cg5TY2hlZHVsZWRRdWVyeRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SJQoMY3JlYXRpb250aW1lGO'
    'bPqjEgASgJUgxjcmVhdGlvbnRpbWUSaQoYZXJyb3JyZXBvcnRjb25maWd1cmF0aW9uGOCd8Gkg'
    'ASgLMioucXVlcnkudGltZXN0cmVhbS5FcnJvclJlcG9ydENvbmZpZ3VyYXRpb25SGGVycm9ycm'
    'Vwb3J0Y29uZmlndXJhdGlvbhJTCg1sYXN0cnVuc3RhdHVzGKmM4NIBIAEoDjIpLnF1ZXJ5LnRp'
    'bWVzdHJlYW0uU2NoZWR1bGVkUXVlcnlSdW5TdGF0dXNSDWxhc3RydW5zdGF0dXMSFQoEbmFtZR'
    'iH5oF/IAEoCVIEbmFtZRIyChJuZXh0aW52b2NhdGlvbnRpbWUYqMSkygEgASgJUhJuZXh0aW52'
    'b2NhdGlvbnRpbWUSOQoWcHJldmlvdXNpbnZvY2F0aW9udGltZRjozssDIAEoCVIWcHJldmlvdX'
    'NpbnZvY2F0aW9udGltZRI/CgVzdGF0ZRiXybLvASABKA4yJS5xdWVyeS50aW1lc3RyZWFtLlNj'
    'aGVkdWxlZFF1ZXJ5U3RhdGVSBXN0YXRlElQKEXRhcmdldGRlc3RpbmF0aW9uGMn+gTogASgLMi'
    'MucXVlcnkudGltZXN0cmVhbS5UYXJnZXREZXN0aW5hdGlvblIRdGFyZ2V0ZGVzdGluYXRpb24=');

@$core.Deprecated('Use scheduledQueryDescriptionDescriptor instead')
const ScheduledQueryDescription$json = {
  '1': 'ScheduledQueryDescription',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'errorreportconfiguration',
      '3': 222039776,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ErrorReportConfiguration',
      '10': 'errorreportconfiguration'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'lastrunsummary',
      '3': 238724003,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ScheduledQueryRunSummary',
      '10': 'lastrunsummary'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nextinvocationtime',
      '3': 424223272,
      '4': 1,
      '5': 9,
      '10': 'nextinvocationtime'
    },
    {
      '1': 'notificationconfiguration',
      '3': 290208045,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.NotificationConfiguration',
      '10': 'notificationconfiguration'
    },
    {
      '1': 'previousinvocationtime',
      '3': 7530344,
      '4': 1,
      '5': 9,
      '10': 'previousinvocationtime'
    },
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
    {
      '1': 'recentlyfailedruns',
      '3': 365830319,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.ScheduledQueryRunSummary',
      '10': 'recentlyfailedruns'
    },
    {
      '1': 'scheduleconfiguration',
      '3': 526075047,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ScheduleConfiguration',
      '10': 'scheduleconfiguration'
    },
    {
      '1': 'scheduledqueryexecutionrolearn',
      '3': 284244182,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryexecutionrolearn'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ScheduledQueryState',
      '10': 'state'
    },
    {
      '1': 'targetconfiguration',
      '3': 499845483,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.TargetConfiguration',
      '10': 'targetconfiguration'
    },
  ],
};

/// Descriptor for `ScheduledQueryDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduledQueryDescriptionDescriptor = $convert.base64Decode(
    'ChlTY2hlZHVsZWRRdWVyeURlc2NyaXB0aW9uEhQKA2Fybhidm+2/ASABKAlSA2FybhIlCgxjcm'
    'VhdGlvbnRpbWUY5s+qMSABKAlSDGNyZWF0aW9udGltZRJpChhlcnJvcnJlcG9ydGNvbmZpZ3Vy'
    'YXRpb24Y4J3waSABKAsyKi5xdWVyeS50aW1lc3RyZWFtLkVycm9yUmVwb3J0Q29uZmlndXJhdG'
    'lvblIYZXJyb3JyZXBvcnRjb25maWd1cmF0aW9uEh0KCGttc2tleWlkGI3JlxYgASgJUghrbXNr'
    'ZXlpZBJVCg5sYXN0cnVuc3VtbWFyeRijx+pxIAEoCzIqLnF1ZXJ5LnRpbWVzdHJlYW0uU2NoZW'
    'R1bGVkUXVlcnlSdW5TdW1tYXJ5Ug5sYXN0cnVuc3VtbWFyeRIVCgRuYW1lGIfmgX8gASgJUgRu'
    'YW1lEjIKEm5leHRpbnZvY2F0aW9udGltZRioxKTKASABKAlSEm5leHRpbnZvY2F0aW9udGltZR'
    'JtChlub3RpZmljYXRpb25jb25maWd1cmF0aW9uGK3ysIoBIAEoCzIrLnF1ZXJ5LnRpbWVzdHJl'
    'YW0uTm90aWZpY2F0aW9uQ29uZmlndXJhdGlvblIZbm90aWZpY2F0aW9uY29uZmlndXJhdGlvbh'
    'I5ChZwcmV2aW91c2ludm9jYXRpb250aW1lGOjOywMgASgJUhZwcmV2aW91c2ludm9jYXRpb250'
    'aW1lEiQKC3F1ZXJ5c3RyaW5nGOfK788BIAEoCVILcXVlcnlzdHJpbmcSXgoScmVjZW50bHlmYW'
    'lsZWRydW5zGK/BuK4BIAMoCzIqLnF1ZXJ5LnRpbWVzdHJlYW0uU2NoZWR1bGVkUXVlcnlSdW5T'
    'dW1tYXJ5UhJyZWNlbnRseWZhaWxlZHJ1bnMSYQoVc2NoZWR1bGVjb25maWd1cmF0aW9uGKeJ7f'
    'oBIAEoCzInLnF1ZXJ5LnRpbWVzdHJlYW0uU2NoZWR1bGVDb25maWd1cmF0aW9uUhVzY2hlZHVs'
    'ZWNvbmZpZ3VyYXRpb24SSgoec2NoZWR1bGVkcXVlcnlleGVjdXRpb25yb2xlYXJuGNbxxIcBIA'
    'EoCVIec2NoZWR1bGVkcXVlcnlleGVjdXRpb25yb2xlYXJuEj8KBXN0YXRlGJfJsu8BIAEoDjIl'
    'LnF1ZXJ5LnRpbWVzdHJlYW0uU2NoZWR1bGVkUXVlcnlTdGF0ZVIFc3RhdGUSWwoTdGFyZ2V0Y2'
    '9uZmlndXJhdGlvbhjrkqzuASABKAsyJS5xdWVyeS50aW1lc3RyZWFtLlRhcmdldENvbmZpZ3Vy'
    'YXRpb25SE3RhcmdldGNvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use scheduledQueryInsightsDescriptor instead')
const ScheduledQueryInsights$json = {
  '1': 'ScheduledQueryInsights',
  '2': [
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ScheduledQueryInsightsMode',
      '10': 'mode'
    },
  ],
};

/// Descriptor for `ScheduledQueryInsights`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduledQueryInsightsDescriptor =
    $convert.base64Decode(
        'ChZTY2hlZHVsZWRRdWVyeUluc2lnaHRzEkQKBG1vZGUYs+65mgEgASgOMiwucXVlcnkudGltZX'
        'N0cmVhbS5TY2hlZHVsZWRRdWVyeUluc2lnaHRzTW9kZVIEbW9kZQ==');

@$core.Deprecated('Use scheduledQueryInsightsResponseDescriptor instead')
const ScheduledQueryInsightsResponse$json = {
  '1': 'ScheduledQueryInsightsResponse',
  '2': [
    {'1': 'outputbytes', '3': 318971400, '4': 1, '5': 3, '10': 'outputbytes'},
    {'1': 'outputrows', '3': 138873322, '4': 1, '5': 3, '10': 'outputrows'},
    {
      '1': 'queryspatialcoverage',
      '3': 200825688,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QuerySpatialCoverage',
      '10': 'queryspatialcoverage'
    },
    {
      '1': 'querytablecount',
      '3': 28893653,
      '4': 1,
      '5': 3,
      '10': 'querytablecount'
    },
    {
      '1': 'querytemporalrange',
      '3': 200957543,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryTemporalRange',
      '10': 'querytemporalrange'
    },
  ],
};

/// Descriptor for `ScheduledQueryInsightsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduledQueryInsightsResponseDescriptor = $convert.base64Decode(
    'Ch5TY2hlZHVsZWRRdWVyeUluc2lnaHRzUmVzcG9uc2USJAoLb3V0cHV0Ynl0ZXMYiLyMmAEgAS'
    'gDUgtvdXRwdXRieXRlcxIhCgpvdXRwdXRyb3dzGOqTnEIgASgDUgpvdXRwdXRyb3dzEl0KFHF1'
    'ZXJ5c3BhdGlhbGNvdmVyYWdlGNi24V8gASgLMiYucXVlcnkudGltZXN0cmVhbS5RdWVyeVNwYX'
    'RpYWxDb3ZlcmFnZVIUcXVlcnlzcGF0aWFsY292ZXJhZ2USKwoPcXVlcnl0YWJsZWNvdW50GNXD'
    '4w0gASgDUg9xdWVyeXRhYmxlY291bnQSVwoScXVlcnl0ZW1wb3JhbHJhbmdlGOe86V8gASgLMi'
    'QucXVlcnkudGltZXN0cmVhbS5RdWVyeVRlbXBvcmFsUmFuZ2VSEnF1ZXJ5dGVtcG9yYWxyYW5n'
    'ZQ==');

@$core.Deprecated('Use scheduledQueryRunSummaryDescriptor instead')
const ScheduledQueryRunSummary$json = {
  '1': 'ScheduledQueryRunSummary',
  '2': [
    {
      '1': 'errorreportlocation',
      '3': 117286905,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ErrorReportLocation',
      '10': 'errorreportlocation'
    },
    {
      '1': 'executionstats',
      '3': 220056093,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ExecutionStats',
      '10': 'executionstats'
    },
    {
      '1': 'failurereason',
      '3': 232322142,
      '4': 1,
      '5': 9,
      '10': 'failurereason'
    },
    {
      '1': 'invocationtime',
      '3': 331845291,
      '4': 1,
      '5': 9,
      '10': 'invocationtime'
    },
    {
      '1': 'queryinsightsresponse',
      '3': 354278130,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ScheduledQueryInsightsResponse',
      '10': 'queryinsightsresponse'
    },
    {
      '1': 'runstatus',
      '3': 293822805,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ScheduledQueryRunStatus',
      '10': 'runstatus'
    },
    {'1': 'triggertime', '3': 268796699, '4': 1, '5': 9, '10': 'triggertime'},
  ],
};

/// Descriptor for `ScheduledQueryRunSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduledQueryRunSummaryDescriptor = $convert.base64Decode(
    'ChhTY2hlZHVsZWRRdWVyeVJ1blN1bW1hcnkSWgoTZXJyb3JyZXBvcnRsb2NhdGlvbhj5z/Y3IA'
    'EoCzIlLnF1ZXJ5LnRpbWVzdHJlYW0uRXJyb3JSZXBvcnRMb2NhdGlvblITZXJyb3JyZXBvcnRs'
    'b2NhdGlvbhJLCg5leGVjdXRpb25zdGF0cxidlPdoIAEoCzIgLnF1ZXJ5LnRpbWVzdHJlYW0uRX'
    'hlY3V0aW9uU3RhdHNSDmV4ZWN1dGlvbnN0YXRzEicKDWZhaWx1cmVyZWFzb24Y3ujjbiABKAlS'
    'DWZhaWx1cmVyZWFzb24SKgoOaW52b2NhdGlvbnRpbWUYq52engEgASgJUg5pbnZvY2F0aW9udG'
    'ltZRJqChVxdWVyeWluc2lnaHRzcmVzcG9uc2UY8rX3qAEgASgLMjAucXVlcnkudGltZXN0cmVh'
    'bS5TY2hlZHVsZWRRdWVyeUluc2lnaHRzUmVzcG9uc2VSFXF1ZXJ5aW5zaWdodHNyZXNwb25zZR'
    'JLCglydW5zdGF0dXMY1cKNjAEgASgOMikucXVlcnkudGltZXN0cmVhbS5TY2hlZHVsZWRRdWVy'
    'eVJ1blN0YXR1c1IJcnVuc3RhdHVzEiQKC3RyaWdnZXJ0aW1lGJuGloABIAEoCVILdHJpZ2dlcn'
    'RpbWU=');

@$core.Deprecated('Use selectColumnDescriptor instead')
const SelectColumn$json = {
  '1': 'SelectColumn',
  '2': [
    {'1': 'aliased', '3': 157931831, '4': 1, '5': 8, '10': 'aliased'},
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.Type',
      '10': 'type'
    },
  ],
};

/// Descriptor for `SelectColumn`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectColumnDescriptor = $convert.base64Decode(
    'CgxTZWxlY3RDb2x1bW4SGwoHYWxpYXNlZBi3sqdLIAEoCFIHYWxpYXNlZBIlCgxkYXRhYmFzZW'
    '5hbWUY3LLZKiABKAlSDGRhdGFiYXNlbmFtZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEiAKCXRh'
    'YmxlbmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZRIuCgR0eXBlGO6g14oBIAEoCzIWLnF1ZXJ5Ln'
    'RpbWVzdHJlYW0uVHlwZVIEdHlwZQ==');

@$core.Deprecated('Use serviceQuotaExceededExceptionDescriptor instead')
const ServiceQuotaExceededException$json = {
  '1': 'ServiceQuotaExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceQuotaExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceQuotaExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch1TZXJ2aWNlUXVvdGFFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use snsConfigurationDescriptor instead')
const SnsConfiguration$json = {
  '1': 'SnsConfiguration',
  '2': [
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `SnsConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List snsConfigurationDescriptor = $convert.base64Decode(
    'ChBTbnNDb25maWd1cmF0aW9uEh0KCHRvcGljYXJuGJz0zg4gASgJUgh0b3BpY2Fybg==');

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

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4Y7cCZsAEgASgJUgtyZXNvdXJjZW'
    'FybhItCgR0YWdzGMHB9rUBIAMoCzIVLnF1ZXJ5LnRpbWVzdHJlYW0uVGFnUgR0YWdz');

@$core.Deprecated('Use tagResourceResponseDescriptor instead')
const TagResourceResponse$json = {
  '1': 'TagResourceResponse',
};

/// Descriptor for `TagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceResponseDescriptor =
    $convert.base64Decode('ChNUYWdSZXNvdXJjZVJlc3BvbnNl');

@$core.Deprecated('Use targetConfigurationDescriptor instead')
const TargetConfiguration$json = {
  '1': 'TargetConfiguration',
  '2': [
    {
      '1': 'timestreamconfiguration',
      '3': 30086935,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.TimestreamConfiguration',
      '10': 'timestreamconfiguration'
    },
  ],
};

/// Descriptor for `TargetConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetConfigurationDescriptor = $convert.base64Decode(
    'ChNUYXJnZXRDb25maWd1cmF0aW9uEmYKF3RpbWVzdHJlYW1jb25maWd1cmF0aW9uGJeurA4gAS'
    'gLMikucXVlcnkudGltZXN0cmVhbS5UaW1lc3RyZWFtQ29uZmlndXJhdGlvblIXdGltZXN0cmVh'
    'bWNvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use targetDestinationDescriptor instead')
const TargetDestination$json = {
  '1': 'TargetDestination',
  '2': [
    {
      '1': 'timestreamdestination',
      '3': 259349965,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.TimestreamDestination',
      '10': 'timestreamdestination'
    },
  ],
};

/// Descriptor for `TargetDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetDestinationDescriptor = $convert.base64Decode(
    'ChFUYXJnZXREZXN0aW5hdGlvbhJgChV0aW1lc3RyZWFtZGVzdGluYXRpb24YzbvVeyABKAsyJy'
    '5xdWVyeS50aW1lc3RyZWFtLlRpbWVzdHJlYW1EZXN0aW5hdGlvblIVdGltZXN0cmVhbWRlc3Rp'
    'bmF0aW9u');

@$core.Deprecated('Use throttlingExceptionDescriptor instead')
const ThrottlingException$json = {
  '1': 'ThrottlingException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ThrottlingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttlingExceptionDescriptor =
    $convert.base64Decode(
        'ChNUaHJvdHRsaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use timeSeriesDataPointDescriptor instead')
const TimeSeriesDataPoint$json = {
  '1': 'TimeSeriesDataPoint',
  '2': [
    {'1': 'time', '3': 535094277, '4': 1, '5': 9, '10': 'time'},
    {
      '1': 'value',
      '3': 289929579,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.Datum',
      '10': 'value'
    },
  ],
};

/// Descriptor for `TimeSeriesDataPoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List timeSeriesDataPointDescriptor = $convert.base64Decode(
    'ChNUaW1lU2VyaWVzRGF0YVBvaW50EhYKBHRpbWUYhciT/wEgASgJUgR0aW1lEjEKBXZhbHVlGO'
    'vyn4oBIAEoCzIXLnF1ZXJ5LnRpbWVzdHJlYW0uRGF0dW1SBXZhbHVl');

@$core.Deprecated('Use timestreamConfigurationDescriptor instead')
const TimestreamConfiguration$json = {
  '1': 'TimestreamConfiguration',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {
      '1': 'dimensionmappings',
      '3': 224443741,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.DimensionMapping',
      '10': 'dimensionmappings'
    },
    {
      '1': 'measurenamecolumn',
      '3': 112141775,
      '4': 1,
      '5': 9,
      '10': 'measurenamecolumn'
    },
    {
      '1': 'mixedmeasuremappings',
      '3': 521774144,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.MixedMeasureMapping',
      '10': 'mixedmeasuremappings'
    },
    {
      '1': 'multimeasuremappings',
      '3': 501736394,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.MultiMeasureMappings',
      '10': 'multimeasuremappings'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {'1': 'timecolumn', '3': 519449367, '4': 1, '5': 9, '10': 'timecolumn'},
  ],
};

/// Descriptor for `TimestreamConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List timestreamConfigurationDescriptor = $convert.base64Decode(
    'ChdUaW1lc3RyZWFtQ29uZmlndXJhdGlvbhIlCgxkYXRhYmFzZW5hbWUY3LLZKiABKAlSDGRhdG'
    'FiYXNlbmFtZRJTChFkaW1lbnNpb25tYXBwaW5ncxjd+oJrIAMoCzIiLnF1ZXJ5LnRpbWVzdHJl'
    'YW0uRGltZW5zaW9uTWFwcGluZ1IRZGltZW5zaW9ubWFwcGluZ3MSLwoRbWVhc3VyZW5hbWVjb2'
    'x1bW4Yz8u8NSABKAlSEW1lYXN1cmVuYW1lY29sdW1uEl0KFG1peGVkbWVhc3VyZW1hcHBpbmdz'
    'GMDI5vgBIAMoCzIlLnF1ZXJ5LnRpbWVzdHJlYW0uTWl4ZWRNZWFzdXJlTWFwcGluZ1IUbWl4ZW'
    'RtZWFzdXJlbWFwcGluZ3MSXgoUbXVsdGltZWFzdXJlbWFwcGluZ3MYysef7wEgASgLMiYucXVl'
    'cnkudGltZXN0cmVhbS5NdWx0aU1lYXN1cmVNYXBwaW5nc1IUbXVsdGltZWFzdXJlbWFwcGluZ3'
    'MSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1lEiIKCnRpbWVjb2x1bW4Yl9bY9wEg'
    'ASgJUgp0aW1lY29sdW1u');

@$core.Deprecated('Use timestreamDestinationDescriptor instead')
const TimestreamDestination$json = {
  '1': 'TimestreamDestination',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `TimestreamDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List timestreamDestinationDescriptor = $convert.base64Decode(
    'ChVUaW1lc3RyZWFtRGVzdGluYXRpb24SJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYm'
    'FzZW5hbWUSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use typeDescriptor instead')
const Type$json = {
  '1': 'Type',
  '2': [
    {
      '1': 'arraycolumninfo',
      '3': 476744353,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ColumnInfo',
      '10': 'arraycolumninfo'
    },
    {
      '1': 'rowcolumninfo',
      '3': 37626042,
      '4': 3,
      '5': 11,
      '6': '.query.timestream.ColumnInfo',
      '10': 'rowcolumninfo'
    },
    {
      '1': 'scalartype',
      '3': 316685622,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ScalarType',
      '10': 'scalartype'
    },
    {
      '1': 'timeseriesmeasurevaluecolumninfo',
      '3': 324829411,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.ColumnInfo',
      '10': 'timeseriesmeasurevaluecolumninfo'
    },
  ],
};

/// Descriptor for `Type`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List typeDescriptor = $convert.base64Decode(
    'CgRUeXBlEkoKD2FycmF5Y29sdW1uaW5mbxihlarjASABKAsyHC5xdWVyeS50aW1lc3RyZWFtLk'
    'NvbHVtbkluZm9SD2FycmF5Y29sdW1uaW5mbxJFCg1yb3djb2x1bW5pbmZvGLrB+BEgAygLMhwu'
    'cXVlcnkudGltZXN0cmVhbS5Db2x1bW5JbmZvUg1yb3djb2x1bW5pbmZvEkAKCnNjYWxhcnR5cG'
    'UYtvqAlwEgASgOMhwucXVlcnkudGltZXN0cmVhbS5TY2FsYXJUeXBlUgpzY2FsYXJ0eXBlEmwK'
    'IHRpbWVzZXJpZXNtZWFzdXJldmFsdWVjb2x1bW5pbmZvGOOB8poBIAEoCzIcLnF1ZXJ5LnRpbW'
    'VzdHJlYW0uQ29sdW1uSW5mb1IgdGltZXNlcmllc21lYXN1cmV2YWx1ZWNvbHVtbmluZm8=');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIkCgtyZXNvdXJjZWFybhjtwJmwASABKAlSC3Jlc291cm'
    'NlYXJuEhwKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use untagResourceResponseDescriptor instead')
const UntagResourceResponse$json = {
  '1': 'UntagResourceResponse',
};

/// Descriptor for `UntagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceResponseDescriptor =
    $convert.base64Decode('ChVVbnRhZ1Jlc291cmNlUmVzcG9uc2U=');

@$core.Deprecated('Use updateAccountSettingsRequestDescriptor instead')
const UpdateAccountSettingsRequest$json = {
  '1': 'UpdateAccountSettingsRequest',
  '2': [
    {'1': 'maxquerytcu', '3': 285295370, '4': 1, '5': 5, '10': 'maxquerytcu'},
    {
      '1': 'querycompute',
      '3': 358277613,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryComputeRequest',
      '10': 'querycompute'
    },
    {
      '1': 'querypricingmodel',
      '3': 453472445,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.QueryPricingModel',
      '10': 'querypricingmodel'
    },
  ],
};

/// Descriptor for `UpdateAccountSettingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAccountSettingsRequestDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVBY2NvdW50U2V0dGluZ3NSZXF1ZXN0EiQKC21heHF1ZXJ5dGN1GIqGhYgBIAEoBV'
    'ILbWF4cXVlcnl0Y3USTQoMcXVlcnljb21wdXRlGO3D66oBIAEoCzIlLnF1ZXJ5LnRpbWVzdHJl'
    'YW0uUXVlcnlDb21wdXRlUmVxdWVzdFIMcXVlcnljb21wdXRlElUKEXF1ZXJ5cHJpY2luZ21vZG'
    'VsGL3hndgBIAEoDjIjLnF1ZXJ5LnRpbWVzdHJlYW0uUXVlcnlQcmljaW5nTW9kZWxSEXF1ZXJ5'
    'cHJpY2luZ21vZGVs');

@$core.Deprecated('Use updateAccountSettingsResponseDescriptor instead')
const UpdateAccountSettingsResponse$json = {
  '1': 'UpdateAccountSettingsResponse',
  '2': [
    {'1': 'maxquerytcu', '3': 285295370, '4': 1, '5': 5, '10': 'maxquerytcu'},
    {
      '1': 'querycompute',
      '3': 358277613,
      '4': 1,
      '5': 11,
      '6': '.query.timestream.QueryComputeResponse',
      '10': 'querycompute'
    },
    {
      '1': 'querypricingmodel',
      '3': 453472445,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.QueryPricingModel',
      '10': 'querypricingmodel'
    },
  ],
};

/// Descriptor for `UpdateAccountSettingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAccountSettingsResponseDescriptor = $convert.base64Decode(
    'Ch1VcGRhdGVBY2NvdW50U2V0dGluZ3NSZXNwb25zZRIkCgttYXhxdWVyeXRjdRiKhoWIASABKA'
    'VSC21heHF1ZXJ5dGN1Ek4KDHF1ZXJ5Y29tcHV0ZRjtw+uqASABKAsyJi5xdWVyeS50aW1lc3Ry'
    'ZWFtLlF1ZXJ5Q29tcHV0ZVJlc3BvbnNlUgxxdWVyeWNvbXB1dGUSVQoRcXVlcnlwcmljaW5nbW'
    '9kZWwYveGd2AEgASgOMiMucXVlcnkudGltZXN0cmVhbS5RdWVyeVByaWNpbmdNb2RlbFIRcXVl'
    'cnlwcmljaW5nbW9kZWw=');

@$core.Deprecated('Use updateScheduledQueryRequestDescriptor instead')
const UpdateScheduledQueryRequest$json = {
  '1': 'UpdateScheduledQueryRequest',
  '2': [
    {
      '1': 'scheduledqueryarn',
      '3': 234602964,
      '4': 1,
      '5': 9,
      '10': 'scheduledqueryarn'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.query.timestream.ScheduledQueryState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `UpdateScheduledQueryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateScheduledQueryRequestDescriptor =
    $convert.base64Decode(
        'ChtVcGRhdGVTY2hlZHVsZWRRdWVyeVJlcXVlc3QSLwoRc2NoZWR1bGVkcXVlcnlhcm4Y1IPvby'
        'ABKAlSEXNjaGVkdWxlZHF1ZXJ5YXJuEj8KBXN0YXRlGJfJsu8BIAEoDjIlLnF1ZXJ5LnRpbWVz'
        'dHJlYW0uU2NoZWR1bGVkUXVlcnlTdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use validationExceptionDescriptor instead')
const ValidationException$json = {
  '1': 'ValidationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ValidationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationExceptionDescriptor =
    $convert.base64Decode(
        'ChNWYWxpZGF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');
