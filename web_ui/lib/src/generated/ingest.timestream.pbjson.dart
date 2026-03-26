// This is a generated file - do not edit.
//
// Generated from ingest.timestream.proto.

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

@$core.Deprecated('Use batchLoadDataFormatDescriptor instead')
const BatchLoadDataFormat$json = {
  '1': 'BatchLoadDataFormat',
  '2': [
    {'1': 'BATCH_LOAD_DATA_FORMAT_CSV', '2': 0},
  ],
};

/// Descriptor for `BatchLoadDataFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List batchLoadDataFormatDescriptor = $convert.base64Decode(
    'ChNCYXRjaExvYWREYXRhRm9ybWF0Eh4KGkJBVENIX0xPQURfREFUQV9GT1JNQVRfQ1NWEAA=');

@$core.Deprecated('Use batchLoadStatusDescriptor instead')
const BatchLoadStatus$json = {
  '1': 'BatchLoadStatus',
  '2': [
    {'1': 'BATCH_LOAD_STATUS_PENDING_RESUME', '2': 0},
    {'1': 'BATCH_LOAD_STATUS_IN_PROGRESS', '2': 1},
    {'1': 'BATCH_LOAD_STATUS_SUCCEEDED', '2': 2},
    {'1': 'BATCH_LOAD_STATUS_PROGRESS_STOPPED', '2': 3},
    {'1': 'BATCH_LOAD_STATUS_CREATED', '2': 4},
    {'1': 'BATCH_LOAD_STATUS_FAILED', '2': 5},
  ],
};

/// Descriptor for `BatchLoadStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List batchLoadStatusDescriptor = $convert.base64Decode(
    'Cg9CYXRjaExvYWRTdGF0dXMSJAogQkFUQ0hfTE9BRF9TVEFUVVNfUEVORElOR19SRVNVTUUQAB'
    'IhCh1CQVRDSF9MT0FEX1NUQVRVU19JTl9QUk9HUkVTUxABEh8KG0JBVENIX0xPQURfU1RBVFVT'
    'X1NVQ0NFRURFRBACEiYKIkJBVENIX0xPQURfU1RBVFVTX1BST0dSRVNTX1NUT1BQRUQQAxIdCh'
    'lCQVRDSF9MT0FEX1NUQVRVU19DUkVBVEVEEAQSHAoYQkFUQ0hfTE9BRF9TVEFUVVNfRkFJTEVE'
    'EAU=');

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

@$core.Deprecated('Use measureValueTypeDescriptor instead')
const MeasureValueType$json = {
  '1': 'MeasureValueType',
  '2': [
    {'1': 'MEASURE_VALUE_TYPE_BIGINT', '2': 0},
    {'1': 'MEASURE_VALUE_TYPE_MULTI', '2': 1},
    {'1': 'MEASURE_VALUE_TYPE_VARCHAR', '2': 2},
    {'1': 'MEASURE_VALUE_TYPE_BOOLEAN', '2': 3},
    {'1': 'MEASURE_VALUE_TYPE_DOUBLE', '2': 4},
    {'1': 'MEASURE_VALUE_TYPE_TIMESTAMP', '2': 5},
  ],
};

/// Descriptor for `MeasureValueType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List measureValueTypeDescriptor = $convert.base64Decode(
    'ChBNZWFzdXJlVmFsdWVUeXBlEh0KGU1FQVNVUkVfVkFMVUVfVFlQRV9CSUdJTlQQABIcChhNRU'
    'FTVVJFX1ZBTFVFX1RZUEVfTVVMVEkQARIeChpNRUFTVVJFX1ZBTFVFX1RZUEVfVkFSQ0hBUhAC'
    'Eh4KGk1FQVNVUkVfVkFMVUVfVFlQRV9CT09MRUFOEAMSHQoZTUVBU1VSRV9WQUxVRV9UWVBFX0'
    'RPVUJMRRAEEiAKHE1FQVNVUkVfVkFMVUVfVFlQRV9USU1FU1RBTVAQBQ==');

@$core.Deprecated('Use partitionKeyEnforcementLevelDescriptor instead')
const PartitionKeyEnforcementLevel$json = {
  '1': 'PartitionKeyEnforcementLevel',
  '2': [
    {'1': 'PARTITION_KEY_ENFORCEMENT_LEVEL_OPTIONAL', '2': 0},
    {'1': 'PARTITION_KEY_ENFORCEMENT_LEVEL_REQUIRED', '2': 1},
  ],
};

/// Descriptor for `PartitionKeyEnforcementLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List partitionKeyEnforcementLevelDescriptor =
    $convert.base64Decode(
        'ChxQYXJ0aXRpb25LZXlFbmZvcmNlbWVudExldmVsEiwKKFBBUlRJVElPTl9LRVlfRU5GT1JDRU'
        '1FTlRfTEVWRUxfT1BUSU9OQUwQABIsCihQQVJUSVRJT05fS0VZX0VORk9SQ0VNRU5UX0xFVkVM'
        'X1JFUVVJUkVEEAE=');

@$core.Deprecated('Use partitionKeyTypeDescriptor instead')
const PartitionKeyType$json = {
  '1': 'PartitionKeyType',
  '2': [
    {'1': 'PARTITION_KEY_TYPE_DIMENSION', '2': 0},
    {'1': 'PARTITION_KEY_TYPE_MEASURE', '2': 1},
  ],
};

/// Descriptor for `PartitionKeyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List partitionKeyTypeDescriptor = $convert.base64Decode(
    'ChBQYXJ0aXRpb25LZXlUeXBlEiAKHFBBUlRJVElPTl9LRVlfVFlQRV9ESU1FTlNJT04QABIeCh'
    'pQQVJUSVRJT05fS0VZX1RZUEVfTUVBU1VSRRAB');

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

@$core.Deprecated('Use tableStatusDescriptor instead')
const TableStatus$json = {
  '1': 'TableStatus',
  '2': [
    {'1': 'TABLE_STATUS_RESTORING', '2': 0},
    {'1': 'TABLE_STATUS_ACTIVE', '2': 1},
    {'1': 'TABLE_STATUS_DELETING', '2': 2},
  ],
};

/// Descriptor for `TableStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tableStatusDescriptor = $convert.base64Decode(
    'CgtUYWJsZVN0YXR1cxIaChZUQUJMRV9TVEFUVVNfUkVTVE9SSU5HEAASFwoTVEFCTEVfU1RBVF'
    'VTX0FDVElWRRABEhkKFVRBQkxFX1NUQVRVU19ERUxFVElORxAC');

@$core.Deprecated('Use timeUnitDescriptor instead')
const TimeUnit$json = {
  '1': 'TimeUnit',
  '2': [
    {'1': 'TIME_UNIT_MILLISECONDS', '2': 0},
    {'1': 'TIME_UNIT_SECONDS', '2': 1},
    {'1': 'TIME_UNIT_MICROSECONDS', '2': 2},
    {'1': 'TIME_UNIT_NANOSECONDS', '2': 3},
  ],
};

/// Descriptor for `TimeUnit`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List timeUnitDescriptor = $convert.base64Decode(
    'CghUaW1lVW5pdBIaChZUSU1FX1VOSVRfTUlMTElTRUNPTkRTEAASFQoRVElNRV9VTklUX1NFQ0'
    '9ORFMQARIaChZUSU1FX1VOSVRfTUlDUk9TRUNPTkRTEAISGQoVVElNRV9VTklUX05BTk9TRUNP'
    'TkRTEAM=');

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

@$core.Deprecated('Use batchLoadProgressReportDescriptor instead')
const BatchLoadProgressReport$json = {
  '1': 'BatchLoadProgressReport',
  '2': [
    {'1': 'bytesmetered', '3': 232712245, '4': 1, '5': 3, '10': 'bytesmetered'},
    {'1': 'filefailures', '3': 44779969, '4': 1, '5': 3, '10': 'filefailures'},
    {
      '1': 'parsefailures',
      '3': 330475280,
      '4': 1,
      '5': 3,
      '10': 'parsefailures'
    },
    {
      '1': 'recordingestionfailures',
      '3': 262421066,
      '4': 1,
      '5': 3,
      '10': 'recordingestionfailures'
    },
    {
      '1': 'recordsingested',
      '3': 159925469,
      '4': 1,
      '5': 3,
      '10': 'recordsingested'
    },
    {
      '1': 'recordsprocessed',
      '3': 27133172,
      '4': 1,
      '5': 3,
      '10': 'recordsprocessed'
    },
  ],
};

/// Descriptor for `BatchLoadProgressReport`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchLoadProgressReportDescriptor = $convert.base64Decode(
    'ChdCYXRjaExvYWRQcm9ncmVzc1JlcG9ydBIlCgxieXRlc21ldGVyZWQYtdD7biABKANSDGJ5dG'
    'VzbWV0ZXJlZBIlCgxmaWxlZmFpbHVyZXMYwZOtFSABKANSDGZpbGVmYWlsdXJlcxIoCg1wYXJz'
    'ZWZhaWx1cmVzGJDOyp0BIAEoA1INcGFyc2VmYWlsdXJlcxI7ChdyZWNvcmRpbmdlc3Rpb25mYW'
    'lsdXJlcxjK9JB9IAEoA1IXcmVjb3JkaW5nZXN0aW9uZmFpbHVyZXMSKwoPcmVjb3Jkc2luZ2Vz'
    'dGVkGN2JoUwgASgDUg9yZWNvcmRzaW5nZXN0ZWQSLQoQcmVjb3Jkc3Byb2Nlc3NlZBj0ifgMIA'
    'EoA1IQcmVjb3Jkc3Byb2Nlc3NlZA==');

@$core.Deprecated('Use batchLoadTaskDescriptor instead')
const BatchLoadTask$json = {
  '1': 'BatchLoadTask',
  '2': [
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {
      '1': 'lastupdatedtime',
      '3': 177046166,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'resumableuntil',
      '3': 340504066,
      '4': 1,
      '5': 9,
      '10': 'resumableuntil'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {'1': 'taskid', '3': 18532514, '4': 1, '5': 9, '10': 'taskid'},
    {
      '1': 'taskstatus',
      '3': 448718071,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.BatchLoadStatus',
      '10': 'taskstatus'
    },
  ],
};

/// Descriptor for `BatchLoadTask`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchLoadTaskDescriptor = $convert.base64Decode(
    'Cg1CYXRjaExvYWRUYXNrEiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3JlYXRpb250aW1lEi'
    'UKDGRhdGFiYXNlbmFtZRjcstkqIAEoCVIMZGF0YWJhc2VuYW1lEisKD2xhc3R1cGRhdGVkdGlt'
    'ZRiWhbZUIAEoCVIPbGFzdHVwZGF0ZWR0aW1lEioKDnJlc3VtYWJsZXVudGlsGILcrqIBIAEoCV'
    'IOcmVzdW1hYmxldW50aWwSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1lEhkKBnRh'
    'c2tpZBiikesIIAEoCVIGdGFza2lkEkYKCnRhc2tzdGF0dXMY98n71QEgASgOMiIuaW5nZXN0Ln'
    'RpbWVzdHJlYW0uQmF0Y2hMb2FkU3RhdHVzUgp0YXNrc3RhdHVz');

@$core.Deprecated('Use batchLoadTaskDescriptionDescriptor instead')
const BatchLoadTaskDescription$json = {
  '1': 'BatchLoadTaskDescription',
  '2': [
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'datamodelconfiguration',
      '3': 195075337,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.DataModelConfiguration',
      '10': 'datamodelconfiguration'
    },
    {
      '1': 'datasourceconfiguration',
      '3': 256337787,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.DataSourceConfiguration',
      '10': 'datasourceconfiguration'
    },
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'lastupdatedtime',
      '3': 177046166,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'progressreport',
      '3': 345169253,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.BatchLoadProgressReport',
      '10': 'progressreport'
    },
    {
      '1': 'recordversion',
      '3': 282309505,
      '4': 1,
      '5': 3,
      '10': 'recordversion'
    },
    {
      '1': 'reportconfiguration',
      '3': 410368276,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.ReportConfiguration',
      '10': 'reportconfiguration'
    },
    {
      '1': 'resumableuntil',
      '3': 340504066,
      '4': 1,
      '5': 9,
      '10': 'resumableuntil'
    },
    {
      '1': 'targetdatabasename',
      '3': 454828599,
      '4': 1,
      '5': 9,
      '10': 'targetdatabasename'
    },
    {
      '1': 'targettablename',
      '3': 298767720,
      '4': 1,
      '5': 9,
      '10': 'targettablename'
    },
    {'1': 'taskid', '3': 18532514, '4': 1, '5': 9, '10': 'taskid'},
    {
      '1': 'taskstatus',
      '3': 448718071,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.BatchLoadStatus',
      '10': 'taskstatus'
    },
  ],
};

/// Descriptor for `BatchLoadTaskDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchLoadTaskDescriptionDescriptor = $convert.base64Decode(
    'ChhCYXRjaExvYWRUYXNrRGVzY3JpcHRpb24SJQoMY3JlYXRpb250aW1lGObPqjEgASgJUgxjcm'
    'VhdGlvbnRpbWUSZAoWZGF0YW1vZGVsY29uZmlndXJhdGlvbhiJuoJdIAEoCzIpLmluZ2VzdC50'
    'aW1lc3RyZWFtLkRhdGFNb2RlbENvbmZpZ3VyYXRpb25SFmRhdGFtb2RlbGNvbmZpZ3VyYXRpb2'
    '4SZwoXZGF0YXNvdXJjZWNvbmZpZ3VyYXRpb24Y+86deiABKAsyKi5pbmdlc3QudGltZXN0cmVh'
    'bS5EYXRhU291cmNlQ29uZmlndXJhdGlvblIXZGF0YXNvdXJjZWNvbmZpZ3VyYXRpb24SJgoMZX'
    'Jyb3JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3JtZXNzYWdlEisKD2xhc3R1cGRhdGVkdGltZRiW'
    'hbZUIAEoCVIPbGFzdHVwZGF0ZWR0aW1lElYKDnByb2dyZXNzcmVwb3J0GOW6y6QBIAEoCzIqLm'
    'luZ2VzdC50aW1lc3RyZWFtLkJhdGNoTG9hZFByb2dyZXNzUmVwb3J0Ug5wcm9ncmVzc3JlcG9y'
    'dBIoCg1yZWNvcmR2ZXJzaW9uGIHnzoYBIAEoA1INcmVjb3JkdmVyc2lvbhJcChNyZXBvcnRjb2'
    '5maWd1cmF0aW9uGJTy1sMBIAEoCzImLmluZ2VzdC50aW1lc3RyZWFtLlJlcG9ydENvbmZpZ3Vy'
    'YXRpb25SE3JlcG9ydGNvbmZpZ3VyYXRpb24SKgoOcmVzdW1hYmxldW50aWwYgtyuogEgASgJUg'
    '5yZXN1bWFibGV1bnRpbBIyChJ0YXJnZXRkYXRhYmFzZW5hbWUYt8Tw2AEgASgJUhJ0YXJnZXRk'
    'YXRhYmFzZW5hbWUSLAoPdGFyZ2V0dGFibGVuYW1lGOiqu44BIAEoCVIPdGFyZ2V0dGFibGVuYW'
    '1lEhkKBnRhc2tpZBiikesIIAEoCVIGdGFza2lkEkYKCnRhc2tzdGF0dXMY98n71QEgASgOMiIu'
    'aW5nZXN0LnRpbWVzdHJlYW0uQmF0Y2hMb2FkU3RhdHVzUgp0YXNrc3RhdHVz');

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

@$core.Deprecated('Use createBatchLoadTaskRequestDescriptor instead')
const CreateBatchLoadTaskRequest$json = {
  '1': 'CreateBatchLoadTaskRequest',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'datamodelconfiguration',
      '3': 195075337,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.DataModelConfiguration',
      '10': 'datamodelconfiguration'
    },
    {
      '1': 'datasourceconfiguration',
      '3': 256337787,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.DataSourceConfiguration',
      '10': 'datasourceconfiguration'
    },
    {
      '1': 'recordversion',
      '3': 282309505,
      '4': 1,
      '5': 3,
      '10': 'recordversion'
    },
    {
      '1': 'reportconfiguration',
      '3': 410368276,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.ReportConfiguration',
      '10': 'reportconfiguration'
    },
    {
      '1': 'targetdatabasename',
      '3': 454828599,
      '4': 1,
      '5': 9,
      '10': 'targetdatabasename'
    },
    {
      '1': 'targettablename',
      '3': 298767720,
      '4': 1,
      '5': 9,
      '10': 'targettablename'
    },
  ],
};

/// Descriptor for `CreateBatchLoadTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBatchLoadTaskRequestDescriptor = $convert.base64Decode(
    'ChpDcmVhdGVCYXRjaExvYWRUYXNrUmVxdWVzdBIjCgtjbGllbnR0b2tlbhjM+7tBIAEoCVILY2'
    'xpZW50dG9rZW4SZAoWZGF0YW1vZGVsY29uZmlndXJhdGlvbhiJuoJdIAEoCzIpLmluZ2VzdC50'
    'aW1lc3RyZWFtLkRhdGFNb2RlbENvbmZpZ3VyYXRpb25SFmRhdGFtb2RlbGNvbmZpZ3VyYXRpb2'
    '4SZwoXZGF0YXNvdXJjZWNvbmZpZ3VyYXRpb24Y+86deiABKAsyKi5pbmdlc3QudGltZXN0cmVh'
    'bS5EYXRhU291cmNlQ29uZmlndXJhdGlvblIXZGF0YXNvdXJjZWNvbmZpZ3VyYXRpb24SKAoNcm'
    'Vjb3JkdmVyc2lvbhiB586GASABKANSDXJlY29yZHZlcnNpb24SXAoTcmVwb3J0Y29uZmlndXJh'
    'dGlvbhiU8tbDASABKAsyJi5pbmdlc3QudGltZXN0cmVhbS5SZXBvcnRDb25maWd1cmF0aW9uUh'
    'NyZXBvcnRjb25maWd1cmF0aW9uEjIKEnRhcmdldGRhdGFiYXNlbmFtZRi3xPDYASABKAlSEnRh'
    'cmdldGRhdGFiYXNlbmFtZRIsCg90YXJnZXR0YWJsZW5hbWUY6Kq7jgEgASgJUg90YXJnZXR0YW'
    'JsZW5hbWU=');

@$core.Deprecated('Use createBatchLoadTaskResponseDescriptor instead')
const CreateBatchLoadTaskResponse$json = {
  '1': 'CreateBatchLoadTaskResponse',
  '2': [
    {'1': 'taskid', '3': 18532514, '4': 1, '5': 9, '10': 'taskid'},
  ],
};

/// Descriptor for `CreateBatchLoadTaskResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createBatchLoadTaskResponseDescriptor =
    $convert.base64Decode(
        'ChtDcmVhdGVCYXRjaExvYWRUYXNrUmVzcG9uc2USGQoGdGFza2lkGKKR6wggASgJUgZ0YXNraW'
        'Q=');

@$core.Deprecated('Use createDatabaseRequestDescriptor instead')
const CreateDatabaseRequest$json = {
  '1': 'CreateDatabaseRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateDatabaseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDatabaseRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVEYXRhYmFzZVJlcXVlc3QSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYm'
    'FzZW5hbWUSHQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEi4KBHRhZ3MYwcH2tQEgAygL'
    'MhYuaW5nZXN0LnRpbWVzdHJlYW0uVGFnUgR0YWdz');

@$core.Deprecated('Use createDatabaseResponseDescriptor instead')
const CreateDatabaseResponse$json = {
  '1': 'CreateDatabaseResponse',
  '2': [
    {
      '1': 'database',
      '3': 278147289,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Database',
      '10': 'database'
    },
  ],
};

/// Descriptor for `CreateDatabaseResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDatabaseResponseDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVEYXRhYmFzZVJlc3BvbnNlEjsKCGRhdGFiYXNlGNnh0IQBIAEoCzIbLmluZ2VzdC'
        '50aW1lc3RyZWFtLkRhdGFiYXNlUghkYXRhYmFzZQ==');

@$core.Deprecated('Use createTableRequestDescriptor instead')
const CreateTableRequest$json = {
  '1': 'CreateTableRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {
      '1': 'magneticstorewriteproperties',
      '3': 127209171,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.MagneticStoreWriteProperties',
      '10': 'magneticstorewriteproperties'
    },
    {
      '1': 'retentionproperties',
      '3': 242841241,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.RetentionProperties',
      '10': 'retentionproperties'
    },
    {
      '1': 'schema',
      '3': 412122455,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Schema',
      '10': 'schema'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateTableRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTableRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVUYWJsZVJlcXVlc3QSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYmFzZW'
    '5hbWUSdgocbWFnbmV0aWNzdG9yZXdyaXRlcHJvcGVydGllcxjTndQ8IAEoCzIvLmluZ2VzdC50'
    'aW1lc3RyZWFtLk1hZ25ldGljU3RvcmVXcml0ZVByb3BlcnRpZXNSHG1hZ25ldGljc3RvcmV3cm'
    'l0ZXByb3BlcnRpZXMSWwoTcmV0ZW50aW9ucHJvcGVydGllcxiZ7eVzIAEoCzImLmluZ2VzdC50'
    'aW1lc3RyZWFtLlJldGVudGlvblByb3BlcnRpZXNSE3JldGVudGlvbnByb3BlcnRpZXMSNQoGc2'
    'NoZW1hGNf6wcQBIAEoCzIZLmluZ2VzdC50aW1lc3RyZWFtLlNjaGVtYVIGc2NoZW1hEiAKCXRh'
    'YmxlbmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZRIuCgR0YWdzGMHB9rUBIAMoCzIWLmluZ2VzdC'
    '50aW1lc3RyZWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createTableResponseDescriptor instead')
const CreateTableResponse$json = {
  '1': 'CreateTableResponse',
  '2': [
    {
      '1': 'table',
      '3': 386722688,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Table',
      '10': 'table'
    },
  ],
};

/// Descriptor for `CreateTableResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTableResponseDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVUYWJsZVJlc3BvbnNlEjIKBXRhYmxlGIDXs7gBIAEoCzIYLmluZ2VzdC50aW1lc3'
    'RyZWFtLlRhYmxlUgV0YWJsZQ==');

@$core.Deprecated('Use csvConfigurationDescriptor instead')
const CsvConfiguration$json = {
  '1': 'CsvConfiguration',
  '2': [
    {
      '1': 'columnseparator',
      '3': 8188515,
      '4': 1,
      '5': 9,
      '10': 'columnseparator'
    },
    {'1': 'escapechar', '3': 210340939, '4': 1, '5': 9, '10': 'escapechar'},
    {'1': 'nullvalue', '3': 440981694, '4': 1, '5': 9, '10': 'nullvalue'},
    {'1': 'quotechar', '3': 242375728, '4': 1, '5': 9, '10': 'quotechar'},
    {
      '1': 'trimwhitespace',
      '3': 386966069,
      '4': 1,
      '5': 8,
      '10': 'trimwhitespace'
    },
  ],
};

/// Descriptor for `CsvConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List csvConfigurationDescriptor = $convert.base64Decode(
    'ChBDc3ZDb25maWd1cmF0aW9uEisKD2NvbHVtbnNlcGFyYXRvchjj5PMDIAEoCVIPY29sdW1uc2'
    'VwYXJhdG9yEiEKCmVzY2FwZWNoYXIYy5imZCABKAlSCmVzY2FwZWNoYXISIAoJbnVsbHZhbHVl'
    'GL6xo9IBIAEoCVIJbnVsbHZhbHVlEh8KCXF1b3RlY2hhchiwuMlzIAEoCVIJcXVvdGVjaGFyEi'
    'oKDnRyaW13aGl0ZXNwYWNlGLXEwrgBIAEoCFIOdHJpbXdoaXRlc3BhY2U=');

@$core.Deprecated('Use dataModelDescriptor instead')
const DataModel$json = {
  '1': 'DataModel',
  '2': [
    {
      '1': 'dimensionmappings',
      '3': 224443741,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.DimensionMapping',
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
      '6': '.ingest.timestream.MixedMeasureMapping',
      '10': 'mixedmeasuremappings'
    },
    {
      '1': 'multimeasuremappings',
      '3': 501736394,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.MultiMeasureMappings',
      '10': 'multimeasuremappings'
    },
    {'1': 'timecolumn', '3': 519449367, '4': 1, '5': 9, '10': 'timecolumn'},
    {
      '1': 'timeunit',
      '3': 181686379,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.TimeUnit',
      '10': 'timeunit'
    },
  ],
};

/// Descriptor for `DataModel`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataModelDescriptor = $convert.base64Decode(
    'CglEYXRhTW9kZWwSVAoRZGltZW5zaW9ubWFwcGluZ3MY3fqCayADKAsyIy5pbmdlc3QudGltZX'
    'N0cmVhbS5EaW1lbnNpb25NYXBwaW5nUhFkaW1lbnNpb25tYXBwaW5ncxIvChFtZWFzdXJlbmFt'
    'ZWNvbHVtbhjPy7w1IAEoCVIRbWVhc3VyZW5hbWVjb2x1bW4SXgoUbWl4ZWRtZWFzdXJlbWFwcG'
    'luZ3MYwMjm+AEgAygLMiYuaW5nZXN0LnRpbWVzdHJlYW0uTWl4ZWRNZWFzdXJlTWFwcGluZ1IU'
    'bWl4ZWRtZWFzdXJlbWFwcGluZ3MSXwoUbXVsdGltZWFzdXJlbWFwcGluZ3MYysef7wEgASgLMi'
    'cuaW5nZXN0LnRpbWVzdHJlYW0uTXVsdGlNZWFzdXJlTWFwcGluZ3NSFG11bHRpbWVhc3VyZW1h'
    'cHBpbmdzEiIKCnRpbWVjb2x1bW4Yl9bY9wEgASgJUgp0aW1lY29sdW1uEjoKCHRpbWV1bml0GO'
    'ug0VYgASgOMhsuaW5nZXN0LnRpbWVzdHJlYW0uVGltZVVuaXRSCHRpbWV1bml0');

@$core.Deprecated('Use dataModelConfigurationDescriptor instead')
const DataModelConfiguration$json = {
  '1': 'DataModelConfiguration',
  '2': [
    {
      '1': 'datamodel',
      '3': 428899579,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.DataModel',
      '10': 'datamodel'
    },
    {
      '1': 'datamodels3configuration',
      '3': 175841835,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.DataModelS3Configuration',
      '10': 'datamodels3configuration'
    },
  ],
};

/// Descriptor for `DataModelConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataModelConfigurationDescriptor = $convert.base64Decode(
    'ChZEYXRhTW9kZWxDb25maWd1cmF0aW9uEj4KCWRhdGFtb2RlbBj7+cHMASABKAsyHC5pbmdlc3'
    'QudGltZXN0cmVhbS5EYXRhTW9kZWxSCWRhdGFtb2RlbBJqChhkYXRhbW9kZWxzM2NvbmZpZ3Vy'
    'YXRpb24Yq8TsUyABKAsyKy5pbmdlc3QudGltZXN0cmVhbS5EYXRhTW9kZWxTM0NvbmZpZ3VyYX'
    'Rpb25SGGRhdGFtb2RlbHMzY29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use dataModelS3ConfigurationDescriptor instead')
const DataModelS3Configuration$json = {
  '1': 'DataModelS3Configuration',
  '2': [
    {'1': 'bucketname', '3': 208117045, '4': 1, '5': 9, '10': 'bucketname'},
    {'1': 'objectkey', '3': 335986226, '4': 1, '5': 9, '10': 'objectkey'},
  ],
};

/// Descriptor for `DataModelS3Configuration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataModelS3ConfigurationDescriptor =
    $convert.base64Decode(
        'ChhEYXRhTW9kZWxTM0NvbmZpZ3VyYXRpb24SIQoKYnVja2V0bmFtZRi1up5jIAEoCVIKYnVja2'
        'V0bmFtZRIgCglvYmplY3RrZXkYsvyaoAEgASgJUglvYmplY3RrZXk=');

@$core.Deprecated('Use dataSourceConfigurationDescriptor instead')
const DataSourceConfiguration$json = {
  '1': 'DataSourceConfiguration',
  '2': [
    {
      '1': 'csvconfiguration',
      '3': 264008448,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.CsvConfiguration',
      '10': 'csvconfiguration'
    },
    {
      '1': 'dataformat',
      '3': 89652083,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.BatchLoadDataFormat',
      '10': 'dataformat'
    },
    {
      '1': 'datasources3configuration',
      '3': 345023213,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.DataSourceS3Configuration',
      '10': 'datasources3configuration'
    },
  ],
};

/// Descriptor for `DataSourceConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataSourceConfigurationDescriptor = $convert.base64Decode(
    'ChdEYXRhU291cmNlQ29uZmlndXJhdGlvbhJSChBjc3Zjb25maWd1cmF0aW9uGIDm8X0gASgLMi'
    'MuaW5nZXN0LnRpbWVzdHJlYW0uQ3N2Q29uZmlndXJhdGlvblIQY3N2Y29uZmlndXJhdGlvbhJJ'
    'CgpkYXRhZm9ybWF0GPP23yogASgOMiYuaW5nZXN0LnRpbWVzdHJlYW0uQmF0Y2hMb2FkRGF0YU'
    'Zvcm1hdFIKZGF0YWZvcm1hdBJuChlkYXRhc291cmNlczNjb25maWd1cmF0aW9uGO3FwqQBIAEo'
    'CzIsLmluZ2VzdC50aW1lc3RyZWFtLkRhdGFTb3VyY2VTM0NvbmZpZ3VyYXRpb25SGWRhdGFzb3'
    'VyY2VzM2NvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use dataSourceS3ConfigurationDescriptor instead')
const DataSourceS3Configuration$json = {
  '1': 'DataSourceS3Configuration',
  '2': [
    {'1': 'bucketname', '3': 208117045, '4': 1, '5': 9, '10': 'bucketname'},
    {
      '1': 'objectkeyprefix',
      '3': 132617574,
      '4': 1,
      '5': 9,
      '10': 'objectkeyprefix'
    },
  ],
};

/// Descriptor for `DataSourceS3Configuration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataSourceS3ConfigurationDescriptor =
    $convert.base64Decode(
        'ChlEYXRhU291cmNlUzNDb25maWd1cmF0aW9uEiEKCmJ1Y2tldG5hbWUYtbqeYyABKAlSCmJ1Y2'
        'tldG5hbWUSKwoPb2JqZWN0a2V5cHJlZml4GOaqnj8gASgJUg9vYmplY3RrZXlwcmVmaXg=');

@$core.Deprecated('Use databaseDescriptor instead')
const Database$json = {
  '1': 'Database',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'lastupdatedtime',
      '3': 177046166,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtime'
    },
    {'1': 'tablecount', '3': 323868967, '4': 1, '5': 3, '10': 'tablecount'},
  ],
};

/// Descriptor for `Database`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List databaseDescriptor = $convert.base64Decode(
    'CghEYXRhYmFzZRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SJQoMY3JlYXRpb250aW1lGObPqjEgAS'
    'gJUgxjcmVhdGlvbnRpbWUSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYmFzZW5hbWUS'
    'HQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlkEisKD2xhc3R1cGRhdGVkdGltZRiWhbZUIA'
    'EoCVIPbGFzdHVwZGF0ZWR0aW1lEiIKCnRhYmxlY291bnQYp7K3mgEgASgDUgp0YWJsZWNvdW50');

@$core.Deprecated('Use deleteDatabaseRequestDescriptor instead')
const DeleteDatabaseRequest$json = {
  '1': 'DeleteDatabaseRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
  ],
};

/// Descriptor for `DeleteDatabaseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDatabaseRequestDescriptor = $convert.base64Decode(
    'ChVEZWxldGVEYXRhYmFzZVJlcXVlc3QSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYm'
    'FzZW5hbWU=');

@$core.Deprecated('Use deleteTableRequestDescriptor instead')
const DeleteTableRequest$json = {
  '1': 'DeleteTableRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DeleteTableRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTableRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVUYWJsZVJlcXVlc3QSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYmFzZW'
    '5hbWUSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1l');

@$core.Deprecated('Use describeBatchLoadTaskRequestDescriptor instead')
const DescribeBatchLoadTaskRequest$json = {
  '1': 'DescribeBatchLoadTaskRequest',
  '2': [
    {'1': 'taskid', '3': 18532514, '4': 1, '5': 9, '10': 'taskid'},
  ],
};

/// Descriptor for `DescribeBatchLoadTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeBatchLoadTaskRequestDescriptor =
    $convert.base64Decode(
        'ChxEZXNjcmliZUJhdGNoTG9hZFRhc2tSZXF1ZXN0EhkKBnRhc2tpZBiikesIIAEoCVIGdGFza2'
        'lk');

@$core.Deprecated('Use describeBatchLoadTaskResponseDescriptor instead')
const DescribeBatchLoadTaskResponse$json = {
  '1': 'DescribeBatchLoadTaskResponse',
  '2': [
    {
      '1': 'batchloadtaskdescription',
      '3': 368256877,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.BatchLoadTaskDescription',
      '10': 'batchloadtaskdescription'
    },
  ],
};

/// Descriptor for `DescribeBatchLoadTaskResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeBatchLoadTaskResponseDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZUJhdGNoTG9hZFRhc2tSZXNwb25zZRJrChhiYXRjaGxvYWR0YXNrZGVzY3JpcH'
        'Rpb24Y7c7MrwEgASgLMisuaW5nZXN0LnRpbWVzdHJlYW0uQmF0Y2hMb2FkVGFza0Rlc2NyaXB0'
        'aW9uUhhiYXRjaGxvYWR0YXNrZGVzY3JpcHRpb24=');

@$core.Deprecated('Use describeDatabaseRequestDescriptor instead')
const DescribeDatabaseRequest$json = {
  '1': 'DescribeDatabaseRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
  ],
};

/// Descriptor for `DescribeDatabaseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDatabaseRequestDescriptor =
    $convert.base64Decode(
        'ChdEZXNjcmliZURhdGFiYXNlUmVxdWVzdBIlCgxkYXRhYmFzZW5hbWUY3LLZKiABKAlSDGRhdG'
        'FiYXNlbmFtZQ==');

@$core.Deprecated('Use describeDatabaseResponseDescriptor instead')
const DescribeDatabaseResponse$json = {
  '1': 'DescribeDatabaseResponse',
  '2': [
    {
      '1': 'database',
      '3': 278147289,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Database',
      '10': 'database'
    },
  ],
};

/// Descriptor for `DescribeDatabaseResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDatabaseResponseDescriptor =
    $convert.base64Decode(
        'ChhEZXNjcmliZURhdGFiYXNlUmVzcG9uc2USOwoIZGF0YWJhc2UY2eHQhAEgASgLMhsuaW5nZX'
        'N0LnRpbWVzdHJlYW0uRGF0YWJhc2VSCGRhdGFiYXNl');

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
      '6': '.ingest.timestream.Endpoint',
      '10': 'endpoints'
    },
  ],
};

/// Descriptor for `DescribeEndpointsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeEndpointsResponseDescriptor =
    $convert.base64Decode(
        'ChlEZXNjcmliZUVuZHBvaW50c1Jlc3BvbnNlEjwKCWVuZHBvaW50cxi+tN0HIAMoCzIbLmluZ2'
        'VzdC50aW1lc3RyZWFtLkVuZHBvaW50UgllbmRwb2ludHM=');

@$core.Deprecated('Use describeTableRequestDescriptor instead')
const DescribeTableRequest$json = {
  '1': 'DescribeTableRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `DescribeTableRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTableRequestDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZVRhYmxlUmVxdWVzdBIlCgxkYXRhYmFzZW5hbWUY3LLZKiABKAlSDGRhdGFiYX'
    'NlbmFtZRIgCgl0YWJsZW5hbWUY3eTagQEgASgJUgl0YWJsZW5hbWU=');

@$core.Deprecated('Use describeTableResponseDescriptor instead')
const DescribeTableResponse$json = {
  '1': 'DescribeTableResponse',
  '2': [
    {
      '1': 'table',
      '3': 386722688,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Table',
      '10': 'table'
    },
  ],
};

/// Descriptor for `DescribeTableResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeTableResponseDescriptor = $convert.base64Decode(
    'ChVEZXNjcmliZVRhYmxlUmVzcG9uc2USMgoFdGFibGUYgNezuAEgASgLMhguaW5nZXN0LnRpbW'
    'VzdHJlYW0uVGFibGVSBXRhYmxl');

@$core.Deprecated('Use dimensionDescriptor instead')
const Dimension$json = {
  '1': 'Dimension',
  '2': [
    {
      '1': 'dimensionvaluetype',
      '3': 267417961,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.DimensionValueType',
      '10': 'dimensionvaluetype'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Dimension`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dimensionDescriptor = $convert.base64Decode(
    'CglEaW1lbnNpb24SWAoSZGltZW5zaW9udmFsdWV0eXBlGOnywX8gASgOMiUuaW5nZXN0LnRpbW'
    'VzdHJlYW0uRGltZW5zaW9uVmFsdWVUeXBlUhJkaW1lbnNpb252YWx1ZXR5cGUSFQoEbmFtZRiH'
    '5oF/IAEoCVIEbmFtZRIYCgV2YWx1ZRjr8p+KASABKAlSBXZhbHVl');

@$core.Deprecated('Use dimensionMappingDescriptor instead')
const DimensionMapping$json = {
  '1': 'DimensionMapping',
  '2': [
    {
      '1': 'destinationcolumn',
      '3': 338930478,
      '4': 1,
      '5': 9,
      '10': 'destinationcolumn'
    },
    {'1': 'sourcecolumn', '3': 219947651, '4': 1, '5': 9, '10': 'sourcecolumn'},
  ],
};

/// Descriptor for `DimensionMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dimensionMappingDescriptor = $convert.base64Decode(
    'ChBEaW1lbnNpb25NYXBwaW5nEjAKEWRlc3RpbmF0aW9uY29sdW1uGK7WzqEBIAEoCVIRZGVzdG'
    'luYXRpb25jb2x1bW4SJQoMc291cmNlY29sdW1uGIPF8GggASgJUgxzb3VyY2Vjb2x1bW4=');

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

@$core.Deprecated('Use listBatchLoadTasksRequestDescriptor instead')
const ListBatchLoadTasksRequest$json = {
  '1': 'ListBatchLoadTasksRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'taskstatus',
      '3': 448718071,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.BatchLoadStatus',
      '10': 'taskstatus'
    },
  ],
};

/// Descriptor for `ListBatchLoadTasksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBatchLoadTasksRequestDescriptor = $convert.base64Decode(
    'ChlMaXN0QmF0Y2hMb2FkVGFza3NSZXF1ZXN0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYX'
    'hyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEkYKCnRhc2tzdGF0dXMY'
    '98n71QEgASgOMiIuaW5nZXN0LnRpbWVzdHJlYW0uQmF0Y2hMb2FkU3RhdHVzUgp0YXNrc3RhdH'
    'Vz');

@$core.Deprecated('Use listBatchLoadTasksResponseDescriptor instead')
const ListBatchLoadTasksResponse$json = {
  '1': 'ListBatchLoadTasksResponse',
  '2': [
    {
      '1': 'batchloadtasks',
      '3': 292572092,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.BatchLoadTask',
      '10': 'batchloadtasks'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListBatchLoadTasksResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listBatchLoadTasksResponseDescriptor =
    $convert.base64Decode(
        'ChpMaXN0QmF0Y2hMb2FkVGFza3NSZXNwb25zZRJMCg5iYXRjaGxvYWR0YXNrcxi8l8GLASADKA'
        'syIC5pbmdlc3QudGltZXN0cmVhbS5CYXRjaExvYWRUYXNrUg5iYXRjaGxvYWR0YXNrcxIfCglu'
        'ZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listDatabasesRequestDescriptor instead')
const ListDatabasesRequest$json = {
  '1': 'ListDatabasesRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDatabasesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDatabasesRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0RGF0YWJhc2VzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdW'
    'x0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listDatabasesResponseDescriptor instead')
const ListDatabasesResponse$json = {
  '1': 'ListDatabasesResponse',
  '2': [
    {
      '1': 'databases',
      '3': 71867698,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.Database',
      '10': 'databases'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDatabasesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDatabasesResponseDescriptor = $convert.base64Decode(
    'ChVMaXN0RGF0YWJhc2VzUmVzcG9uc2USPAoJZGF0YWJhc2VzGLK6oiIgAygLMhsuaW5nZXN0Ln'
    'RpbWVzdHJlYW0uRGF0YWJhc2VSCWRhdGFiYXNlcxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5l'
    'eHR0b2tlbg==');

@$core.Deprecated('Use listTablesRequestDescriptor instead')
const ListTablesRequest$json = {
  '1': 'ListTablesRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListTablesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTablesRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0VGFibGVzUmVxdWVzdBIlCgxkYXRhYmFzZW5hbWUY3LLZKiABKAlSDGRhdGFiYXNlbm'
    'FtZRIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6'
    'ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listTablesResponseDescriptor instead')
const ListTablesResponse$json = {
  '1': 'ListTablesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'tables',
      '3': 357958629,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.Table',
      '10': 'tables'
    },
  ],
};

/// Descriptor for `ListTablesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTablesResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0VGFibGVzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SNA'
    'oGdGFibGVzGOWH2KoBIAMoCzIYLmluZ2VzdC50aW1lc3RyZWFtLlRhYmxlUgZ0YWJsZXM=');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIkCgtyZXNvdXJjZWFybhjtwJmwASABKAlSC3'
        'Jlc291cmNlYXJu');

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse$json = {
  '1': 'ListTagsForResourceResponse',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USLgoEdGFncxjBwfa1ASADKAsyFi5pbmdlc3'
        'QudGltZXN0cmVhbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use magneticStoreRejectedDataLocationDescriptor instead')
const MagneticStoreRejectedDataLocation$json = {
  '1': 'MagneticStoreRejectedDataLocation',
  '2': [
    {
      '1': 's3configuration',
      '3': 27828476,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.S3Configuration',
      '10': 's3configuration'
    },
  ],
};

/// Descriptor for `MagneticStoreRejectedDataLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List magneticStoreRejectedDataLocationDescriptor =
    $convert.base64Decode(
        'CiFNYWduZXRpY1N0b3JlUmVqZWN0ZWREYXRhTG9jYXRpb24STwoPczNjb25maWd1cmF0aW9uGP'
        'zBog0gASgLMiIuaW5nZXN0LnRpbWVzdHJlYW0uUzNDb25maWd1cmF0aW9uUg9zM2NvbmZpZ3Vy'
        'YXRpb24=');

@$core.Deprecated('Use magneticStoreWritePropertiesDescriptor instead')
const MagneticStoreWriteProperties$json = {
  '1': 'MagneticStoreWriteProperties',
  '2': [
    {
      '1': 'enablemagneticstorewrites',
      '3': 489945306,
      '4': 1,
      '5': 8,
      '10': 'enablemagneticstorewrites'
    },
    {
      '1': 'magneticstorerejecteddatalocation',
      '3': 47660222,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.MagneticStoreRejectedDataLocation',
      '10': 'magneticstorerejecteddatalocation'
    },
  ],
};

/// Descriptor for `MagneticStoreWriteProperties`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List magneticStoreWritePropertiesDescriptor = $convert.base64Decode(
    'ChxNYWduZXRpY1N0b3JlV3JpdGVQcm9wZXJ0aWVzEkAKGWVuYWJsZW1hZ25ldGljc3RvcmV3cm'
    'l0ZXMY2vHP6QEgASgIUhllbmFibGVtYWduZXRpY3N0b3Jld3JpdGVzEoUBCiFtYWduZXRpY3N0'
    'b3JlcmVqZWN0ZWRkYXRhbG9jYXRpb24YvvncFiABKAsyNC5pbmdlc3QudGltZXN0cmVhbS5NYW'
    'duZXRpY1N0b3JlUmVqZWN0ZWREYXRhTG9jYXRpb25SIW1hZ25ldGljc3RvcmVyZWplY3RlZGRh'
    'dGFsb2NhdGlvbg==');

@$core.Deprecated('Use measureValueDescriptor instead')
const MeasureValue$json = {
  '1': 'MeasureValue',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.MeasureValueType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `MeasureValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List measureValueDescriptor = $convert.base64Decode(
    'CgxNZWFzdXJlVmFsdWUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRI7CgR0eXBlGO6g14oBIAEoDj'
    'IjLmluZ2VzdC50aW1lc3RyZWFtLk1lYXN1cmVWYWx1ZVR5cGVSBHR5cGUSGAoFdmFsdWUY6/Kf'
    'igEgASgJUgV2YWx1ZQ==');

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
      '6': '.ingest.timestream.MeasureValueType',
      '10': 'measurevaluetype'
    },
    {
      '1': 'multimeasureattributemappings',
      '3': 311133918,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.MultiMeasureAttributeMapping',
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
    '5hbWUSUwoQbWVhc3VyZXZhbHVldHlwZRidisTeASABKA4yIy5pbmdlc3QudGltZXN0cmVhbS5N'
    'ZWFzdXJlVmFsdWVUeXBlUhBtZWFzdXJldmFsdWV0eXBlEnkKHW11bHRpbWVhc3VyZWF0dHJpYn'
    'V0ZW1hcHBpbmdzGN6NrpQBIAMoCzIvLmluZ2VzdC50aW1lc3RyZWFtLk11bHRpTWVhc3VyZUF0'
    'dHJpYnV0ZU1hcHBpbmdSHW11bHRpbWVhc3VyZWF0dHJpYnV0ZW1hcHBpbmdzEiUKDHNvdXJjZW'
    'NvbHVtbhiDxfBoIAEoCVIMc291cmNlY29sdW1uEjAKEXRhcmdldG1lYXN1cmVuYW1lGNzB8N8B'
    'IAEoCVIRdGFyZ2V0bWVhc3VyZW5hbWU=');

@$core.Deprecated('Use multiMeasureAttributeMappingDescriptor instead')
const MultiMeasureAttributeMapping$json = {
  '1': 'MultiMeasureAttributeMapping',
  '2': [
    {
      '1': 'measurevaluetype',
      '3': 466683165,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.ScalarMeasureValueType',
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
    'ChxNdWx0aU1lYXN1cmVBdHRyaWJ1dGVNYXBwaW5nElkKEG1lYXN1cmV2YWx1ZXR5cGUYnYrE3g'
    'EgASgOMikuaW5nZXN0LnRpbWVzdHJlYW0uU2NhbGFyTWVhc3VyZVZhbHVlVHlwZVIQbWVhc3Vy'
    'ZXZhbHVldHlwZRIlCgxzb3VyY2Vjb2x1bW4Yg8XwaCABKAlSDHNvdXJjZWNvbHVtbhJMCh90YX'
    'JnZXRtdWx0aW1lYXN1cmVhdHRyaWJ1dGVuYW1lGO/Tl8YBIAEoCVIfdGFyZ2V0bXVsdGltZWFz'
    'dXJlYXR0cmlidXRlbmFtZQ==');

@$core.Deprecated('Use multiMeasureMappingsDescriptor instead')
const MultiMeasureMappings$json = {
  '1': 'MultiMeasureMappings',
  '2': [
    {
      '1': 'multimeasureattributemappings',
      '3': 311133918,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.MultiMeasureAttributeMapping',
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
    'ChRNdWx0aU1lYXN1cmVNYXBwaW5ncxJ5Ch1tdWx0aW1lYXN1cmVhdHRyaWJ1dGVtYXBwaW5ncx'
    'jeja6UASADKAsyLy5pbmdlc3QudGltZXN0cmVhbS5NdWx0aU1lYXN1cmVBdHRyaWJ1dGVNYXBw'
    'aW5nUh1tdWx0aW1lYXN1cmVhdHRyaWJ1dGVtYXBwaW5ncxI5ChZ0YXJnZXRtdWx0aW1lYXN1cm'
    'VuYW1lGJXLiQ0gASgJUhZ0YXJnZXRtdWx0aW1lYXN1cmVuYW1l');

@$core.Deprecated('Use partitionKeyDescriptor instead')
const PartitionKey$json = {
  '1': 'PartitionKey',
  '2': [
    {
      '1': 'enforcementinrecord',
      '3': 121680624,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.PartitionKeyEnforcementLevel',
      '10': 'enforcementinrecord'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.PartitionKeyType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `PartitionKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List partitionKeyDescriptor = $convert.base64Decode(
    'CgxQYXJ0aXRpb25LZXkSZAoTZW5mb3JjZW1lbnRpbnJlY29yZBjw5YI6IAEoDjIvLmluZ2VzdC'
    '50aW1lc3RyZWFtLlBhcnRpdGlvbktleUVuZm9yY2VtZW50TGV2ZWxSE2VuZm9yY2VtZW50aW5y'
    'ZWNvcmQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRI7CgR0eXBlGO6g14oBIAEoDjIjLmluZ2VzdC'
    '50aW1lc3RyZWFtLlBhcnRpdGlvbktleVR5cGVSBHR5cGU=');

@$core.Deprecated('Use recordDescriptor instead')
const Record$json = {
  '1': 'Record',
  '2': [
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.Dimension',
      '10': 'dimensions'
    },
    {'1': 'measurename', '3': 426079069, '4': 1, '5': 9, '10': 'measurename'},
    {'1': 'measurevalue', '3': 407670165, '4': 1, '5': 9, '10': 'measurevalue'},
    {
      '1': 'measurevaluetype',
      '3': 466683165,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.MeasureValueType',
      '10': 'measurevaluetype'
    },
    {
      '1': 'measurevalues',
      '3': 126050982,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.MeasureValue',
      '10': 'measurevalues'
    },
    {'1': 'time', '3': 535094277, '4': 1, '5': 9, '10': 'time'},
    {
      '1': 'timeunit',
      '3': 181686379,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.TimeUnit',
      '10': 'timeunit'
    },
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `Record`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recordDescriptor = $convert.base64Decode(
    'CgZSZWNvcmQSQAoKZGltZW5zaW9ucxjRm9/cASADKAsyHC5pbmdlc3QudGltZXN0cmVhbS5EaW'
    '1lbnNpb25SCmRpbWVuc2lvbnMSJAoLbWVhc3VyZW5hbWUY3eaVywEgASgJUgttZWFzdXJlbmFt'
    'ZRImCgxtZWFzdXJldmFsdWUYlZuywgEgASgJUgxtZWFzdXJldmFsdWUSUwoQbWVhc3VyZXZhbH'
    'VldHlwZRidisTeASABKA4yIy5pbmdlc3QudGltZXN0cmVhbS5NZWFzdXJlVmFsdWVUeXBlUhBt'
    'ZWFzdXJldmFsdWV0eXBlEkgKDW1lYXN1cmV2YWx1ZXMYpsWNPCADKAsyHy5pbmdlc3QudGltZX'
    'N0cmVhbS5NZWFzdXJlVmFsdWVSDW1lYXN1cmV2YWx1ZXMSFgoEdGltZRiFyJP/ASABKAlSBHRp'
    'bWUSOgoIdGltZXVuaXQY66DRViABKA4yGy5pbmdlc3QudGltZXN0cmVhbS5UaW1lVW5pdFIIdG'
    'ltZXVuaXQSHAoHdmVyc2lvbhi4qrfuASABKANSB3ZlcnNpb24=');

@$core.Deprecated('Use recordsIngestedDescriptor instead')
const RecordsIngested$json = {
  '1': 'RecordsIngested',
  '2': [
    {
      '1': 'magneticstore',
      '3': 200805165,
      '4': 1,
      '5': 5,
      '10': 'magneticstore'
    },
    {'1': 'memorystore', '3': 323607360, '4': 1, '5': 5, '10': 'memorystore'},
    {'1': 'total', '3': 218525086, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `RecordsIngested`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recordsIngestedDescriptor = $convert.base64Decode(
    'Cg9SZWNvcmRzSW5nZXN0ZWQSJwoNbWFnbmV0aWNzdG9yZRitluBfIAEoBVINbWFnbmV0aWNzdG'
    '9yZRIkCgttZW1vcnlzdG9yZRjAtqeaASABKAVSC21lbW9yeXN0b3JlEhcKBXRvdGFsGJ7bmWgg'
    'ASgFUgV0b3RhbA==');

@$core.Deprecated('Use rejectedRecordDescriptor instead')
const RejectedRecord$json = {
  '1': 'RejectedRecord',
  '2': [
    {
      '1': 'existingversion',
      '3': 457856209,
      '4': 1,
      '5': 3,
      '10': 'existingversion'
    },
    {'1': 'reason', '3': 20005178, '4': 1, '5': 9, '10': 'reason'},
    {'1': 'recordindex', '3': 221086889, '4': 1, '5': 5, '10': 'recordindex'},
  ],
};

/// Descriptor for `RejectedRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rejectedRecordDescriptor = $convert.base64Decode(
    'Cg5SZWplY3RlZFJlY29yZBIsCg9leGlzdGluZ3ZlcnNpb24Y0amp2gEgASgDUg9leGlzdGluZ3'
    'ZlcnNpb24SGQoGcmVhc29uGLqCxQkgASgJUgZyZWFzb24SIwoLcmVjb3JkaW5kZXgYqYm2aSAB'
    'KAVSC3JlY29yZGluZGV4');

@$core.Deprecated('Use rejectedRecordsExceptionDescriptor instead')
const RejectedRecordsException$json = {
  '1': 'RejectedRecordsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'rejectedrecords',
      '3': 374345864,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.RejectedRecord',
      '10': 'rejectedrecords'
    },
  ],
};

/// Descriptor for `RejectedRecordsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rejectedRecordsExceptionDescriptor = $convert.base64Decode(
    'ChhSZWplY3RlZFJlY29yZHNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
    'JPCg9yZWplY3RlZHJlY29yZHMYiKHAsgEgAygLMiEuaW5nZXN0LnRpbWVzdHJlYW0uUmVqZWN0'
    'ZWRSZWNvcmRSD3JlamVjdGVkcmVjb3Jkcw==');

@$core.Deprecated('Use reportConfigurationDescriptor instead')
const ReportConfiguration$json = {
  '1': 'ReportConfiguration',
  '2': [
    {
      '1': 'reports3configuration',
      '3': 463435614,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.ReportS3Configuration',
      '10': 'reports3configuration'
    },
  ],
};

/// Descriptor for `ReportConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reportConfigurationDescriptor = $convert.base64Decode(
    'ChNSZXBvcnRDb25maWd1cmF0aW9uEmIKFXJlcG9ydHMzY29uZmlndXJhdGlvbhje7v3cASABKA'
    'syKC5pbmdlc3QudGltZXN0cmVhbS5SZXBvcnRTM0NvbmZpZ3VyYXRpb25SFXJlcG9ydHMzY29u'
    'ZmlndXJhdGlvbg==');

@$core.Deprecated('Use reportS3ConfigurationDescriptor instead')
const ReportS3Configuration$json = {
  '1': 'ReportS3Configuration',
  '2': [
    {'1': 'bucketname', '3': 208117045, '4': 1, '5': 9, '10': 'bucketname'},
    {
      '1': 'encryptionoption',
      '3': 160833062,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.S3EncryptionOption',
      '10': 'encryptionoption'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
    {
      '1': 'objectkeyprefix',
      '3': 132617574,
      '4': 1,
      '5': 9,
      '10': 'objectkeyprefix'
    },
  ],
};

/// Descriptor for `ReportS3Configuration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reportS3ConfigurationDescriptor = $convert.base64Decode(
    'ChVSZXBvcnRTM0NvbmZpZ3VyYXRpb24SIQoKYnVja2V0bmFtZRi1up5jIAEoCVIKYnVja2V0bm'
    'FtZRJUChBlbmNyeXB0aW9ub3B0aW9uGKa82EwgASgOMiUuaW5nZXN0LnRpbWVzdHJlYW0uUzNF'
    'bmNyeXB0aW9uT3B0aW9uUhBlbmNyeXB0aW9ub3B0aW9uEh0KCGttc2tleWlkGI3JlxYgASgJUg'
    'hrbXNrZXlpZBIrCg9vYmplY3RrZXlwcmVmaXgY5qqePyABKAlSD29iamVjdGtleXByZWZpeA==');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use resumeBatchLoadTaskRequestDescriptor instead')
const ResumeBatchLoadTaskRequest$json = {
  '1': 'ResumeBatchLoadTaskRequest',
  '2': [
    {'1': 'taskid', '3': 18532514, '4': 1, '5': 9, '10': 'taskid'},
  ],
};

/// Descriptor for `ResumeBatchLoadTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resumeBatchLoadTaskRequestDescriptor =
    $convert.base64Decode(
        'ChpSZXN1bWVCYXRjaExvYWRUYXNrUmVxdWVzdBIZCgZ0YXNraWQYopHrCCABKAlSBnRhc2tpZA'
        '==');

@$core.Deprecated('Use resumeBatchLoadTaskResponseDescriptor instead')
const ResumeBatchLoadTaskResponse$json = {
  '1': 'ResumeBatchLoadTaskResponse',
};

/// Descriptor for `ResumeBatchLoadTaskResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resumeBatchLoadTaskResponseDescriptor =
    $convert.base64Decode('ChtSZXN1bWVCYXRjaExvYWRUYXNrUmVzcG9uc2U=');

@$core.Deprecated('Use retentionPropertiesDescriptor instead')
const RetentionProperties$json = {
  '1': 'RetentionProperties',
  '2': [
    {
      '1': 'magneticstoreretentionperiodindays',
      '3': 65073472,
      '4': 1,
      '5': 3,
      '10': 'magneticstoreretentionperiodindays'
    },
    {
      '1': 'memorystoreretentionperiodinhours',
      '3': 136905329,
      '4': 1,
      '5': 3,
      '10': 'memorystoreretentionperiodinhours'
    },
  ],
};

/// Descriptor for `RetentionProperties`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List retentionPropertiesDescriptor = $convert.base64Decode(
    'ChNSZXRlbnRpb25Qcm9wZXJ0aWVzElEKIm1hZ25ldGljc3RvcmVyZXRlbnRpb25wZXJpb2Rpbm'
    'RheXMYwOKDHyABKANSIm1hZ25ldGljc3RvcmVyZXRlbnRpb25wZXJpb2RpbmRheXMSTwohbWVt'
    'b3J5c3RvcmVyZXRlbnRpb25wZXJpb2RpbmhvdXJzGPGEpEEgASgDUiFtZW1vcnlzdG9yZXJldG'
    'VudGlvbnBlcmlvZGluaG91cnM=');

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
      '6': '.ingest.timestream.S3EncryptionOption',
      '10': 'encryptionoption'
    },
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
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
    'Cg9TM0NvbmZpZ3VyYXRpb24SIQoKYnVja2V0bmFtZRi1up5jIAEoCVIKYnVja2V0bmFtZRJUCh'
    'BlbmNyeXB0aW9ub3B0aW9uGKa82EwgASgOMiUuaW5nZXN0LnRpbWVzdHJlYW0uUzNFbmNyeXB0'
    'aW9uT3B0aW9uUhBlbmNyeXB0aW9ub3B0aW9uEh0KCGttc2tleWlkGI3JlxYgASgJUghrbXNrZX'
    'lpZBIrCg9vYmplY3RrZXlwcmVmaXgY5qqePyABKAlSD29iamVjdGtleXByZWZpeA==');

@$core.Deprecated('Use schemaDescriptor instead')
const Schema$json = {
  '1': 'Schema',
  '2': [
    {
      '1': 'compositepartitionkey',
      '3': 338450538,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.PartitionKey',
      '10': 'compositepartitionkey'
    },
  ],
};

/// Descriptor for `Schema`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List schemaDescriptor = $convert.base64Decode(
    'CgZTY2hlbWESWQoVY29tcG9zaXRlcGFydGl0aW9ua2V5GOqwsaEBIAMoCzIfLmluZ2VzdC50aW'
    '1lc3RyZWFtLlBhcnRpdGlvbktleVIVY29tcG9zaXRlcGFydGl0aW9ua2V5');

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

@$core.Deprecated('Use tableDescriptor instead')
const Table$json = {
  '1': 'Table',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {
      '1': 'lastupdatedtime',
      '3': 177046166,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtime'
    },
    {
      '1': 'magneticstorewriteproperties',
      '3': 127209171,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.MagneticStoreWriteProperties',
      '10': 'magneticstorewriteproperties'
    },
    {
      '1': 'retentionproperties',
      '3': 242841241,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.RetentionProperties',
      '10': 'retentionproperties'
    },
    {
      '1': 'schema',
      '3': 412122455,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Schema',
      '10': 'schema'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {
      '1': 'tablestatus',
      '3': 207908810,
      '4': 1,
      '5': 14,
      '6': '.ingest.timestream.TableStatus',
      '10': 'tablestatus'
    },
  ],
};

/// Descriptor for `Table`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableDescriptor = $convert.base64Decode(
    'CgVUYWJsZRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SJQoMY3JlYXRpb250aW1lGObPqjEgASgJUg'
    'xjcmVhdGlvbnRpbWUSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYmFzZW5hbWUSKwoP'
    'bGFzdHVwZGF0ZWR0aW1lGJaFtlQgASgJUg9sYXN0dXBkYXRlZHRpbWUSdgocbWFnbmV0aWNzdG'
    '9yZXdyaXRlcHJvcGVydGllcxjTndQ8IAEoCzIvLmluZ2VzdC50aW1lc3RyZWFtLk1hZ25ldGlj'
    'U3RvcmVXcml0ZVByb3BlcnRpZXNSHG1hZ25ldGljc3RvcmV3cml0ZXByb3BlcnRpZXMSWwoTcm'
    'V0ZW50aW9ucHJvcGVydGllcxiZ7eVzIAEoCzImLmluZ2VzdC50aW1lc3RyZWFtLlJldGVudGlv'
    'blByb3BlcnRpZXNSE3JldGVudGlvbnByb3BlcnRpZXMSNQoGc2NoZW1hGNf6wcQBIAEoCzIZLm'
    'luZ2VzdC50aW1lc3RyZWFtLlNjaGVtYVIGc2NoZW1hEiAKCXRhYmxlbmFtZRjd5NqBASABKAlS'
    'CXRhYmxlbmFtZRJDCgt0YWJsZXN0YXR1cxjK35FjIAEoDjIeLmluZ2VzdC50aW1lc3RyZWFtLl'
    'RhYmxlU3RhdHVzUgt0YWJsZXN0YXR1cw==');

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
      '6': '.ingest.timestream.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4Y7cCZsAEgASgJUgtyZXNvdXJjZW'
    'FybhIuCgR0YWdzGMHB9rUBIAMoCzIWLmluZ2VzdC50aW1lc3RyZWFtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use tagResourceResponseDescriptor instead')
const TagResourceResponse$json = {
  '1': 'TagResourceResponse',
};

/// Descriptor for `TagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceResponseDescriptor =
    $convert.base64Decode('ChNUYWdSZXNvdXJjZVJlc3BvbnNl');

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

@$core.Deprecated('Use updateDatabaseRequestDescriptor instead')
const UpdateDatabaseRequest$json = {
  '1': 'UpdateDatabaseRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'kmskeyid', '3': 46523533, '4': 1, '5': 9, '10': 'kmskeyid'},
  ],
};

/// Descriptor for `UpdateDatabaseRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDatabaseRequestDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVEYXRhYmFzZVJlcXVlc3QSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYm'
    'FzZW5hbWUSHQoIa21za2V5aWQYjcmXFiABKAlSCGttc2tleWlk');

@$core.Deprecated('Use updateDatabaseResponseDescriptor instead')
const UpdateDatabaseResponse$json = {
  '1': 'UpdateDatabaseResponse',
  '2': [
    {
      '1': 'database',
      '3': 278147289,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Database',
      '10': 'database'
    },
  ],
};

/// Descriptor for `UpdateDatabaseResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDatabaseResponseDescriptor =
    $convert.base64Decode(
        'ChZVcGRhdGVEYXRhYmFzZVJlc3BvbnNlEjsKCGRhdGFiYXNlGNnh0IQBIAEoCzIbLmluZ2VzdC'
        '50aW1lc3RyZWFtLkRhdGFiYXNlUghkYXRhYmFzZQ==');

@$core.Deprecated('Use updateTableRequestDescriptor instead')
const UpdateTableRequest$json = {
  '1': 'UpdateTableRequest',
  '2': [
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {
      '1': 'magneticstorewriteproperties',
      '3': 127209171,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.MagneticStoreWriteProperties',
      '10': 'magneticstorewriteproperties'
    },
    {
      '1': 'retentionproperties',
      '3': 242841241,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.RetentionProperties',
      '10': 'retentionproperties'
    },
    {
      '1': 'schema',
      '3': 412122455,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Schema',
      '10': 'schema'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `UpdateTableRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTableRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVUYWJsZVJlcXVlc3QSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYmFzZW'
    '5hbWUSdgocbWFnbmV0aWNzdG9yZXdyaXRlcHJvcGVydGllcxjTndQ8IAEoCzIvLmluZ2VzdC50'
    'aW1lc3RyZWFtLk1hZ25ldGljU3RvcmVXcml0ZVByb3BlcnRpZXNSHG1hZ25ldGljc3RvcmV3cm'
    'l0ZXByb3BlcnRpZXMSWwoTcmV0ZW50aW9ucHJvcGVydGllcxiZ7eVzIAEoCzImLmluZ2VzdC50'
    'aW1lc3RyZWFtLlJldGVudGlvblByb3BlcnRpZXNSE3JldGVudGlvbnByb3BlcnRpZXMSNQoGc2'
    'NoZW1hGNf6wcQBIAEoCzIZLmluZ2VzdC50aW1lc3RyZWFtLlNjaGVtYVIGc2NoZW1hEiAKCXRh'
    'YmxlbmFtZRjd5NqBASABKAlSCXRhYmxlbmFtZQ==');

@$core.Deprecated('Use updateTableResponseDescriptor instead')
const UpdateTableResponse$json = {
  '1': 'UpdateTableResponse',
  '2': [
    {
      '1': 'table',
      '3': 386722688,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Table',
      '10': 'table'
    },
  ],
};

/// Descriptor for `UpdateTableResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateTableResponseDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVUYWJsZVJlc3BvbnNlEjIKBXRhYmxlGIDXs7gBIAEoCzIYLmluZ2VzdC50aW1lc3'
    'RyZWFtLlRhYmxlUgV0YWJsZQ==');

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

@$core.Deprecated('Use writeRecordsRequestDescriptor instead')
const WriteRecordsRequest$json = {
  '1': 'WriteRecordsRequest',
  '2': [
    {
      '1': 'commonattributes',
      '3': 203721350,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.Record',
      '10': 'commonattributes'
    },
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {
      '1': 'records',
      '3': 423557454,
      '4': 3,
      '5': 11,
      '6': '.ingest.timestream.Record',
      '10': 'records'
    },
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
  ],
};

/// Descriptor for `WriteRecordsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List writeRecordsRequestDescriptor = $convert.base64Decode(
    'ChNXcml0ZVJlY29yZHNSZXF1ZXN0EkgKEGNvbW1vbmF0dHJpYnV0ZXMYhpWSYSABKAsyGS5pbm'
    'dlc3QudGltZXN0cmVhbS5SZWNvcmRSEGNvbW1vbmF0dHJpYnV0ZXMSJQoMZGF0YWJhc2VuYW1l'
    'GNyy2SogASgJUgxkYXRhYmFzZW5hbWUSNwoHcmVjb3JkcxjO8vvJASADKAsyGS5pbmdlc3QudG'
    'ltZXN0cmVhbS5SZWNvcmRSB3JlY29yZHMSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVu'
    'YW1l');

@$core.Deprecated('Use writeRecordsResponseDescriptor instead')
const WriteRecordsResponse$json = {
  '1': 'WriteRecordsResponse',
  '2': [
    {
      '1': 'recordsingested',
      '3': 159925469,
      '4': 1,
      '5': 11,
      '6': '.ingest.timestream.RecordsIngested',
      '10': 'recordsingested'
    },
  ],
};

/// Descriptor for `WriteRecordsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List writeRecordsResponseDescriptor = $convert.base64Decode(
    'ChRXcml0ZVJlY29yZHNSZXNwb25zZRJPCg9yZWNvcmRzaW5nZXN0ZWQY3YmhTCABKAsyIi5pbm'
    'dlc3QudGltZXN0cmVhbS5SZWNvcmRzSW5nZXN0ZWRSD3JlY29yZHNpbmdlc3RlZA==');
