// This is a generated file - do not edit.
//
// Generated from athena.proto.

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

@$core.Deprecated('Use authenticationTypeDescriptor instead')
const AuthenticationType$json = {
  '1': 'AuthenticationType',
  '2': [
    {'1': 'AUTHENTICATION_TYPE_DIRECTORY_IDENTITY', '2': 0},
  ],
};

/// Descriptor for `AuthenticationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List authenticationTypeDescriptor = $convert.base64Decode(
    'ChJBdXRoZW50aWNhdGlvblR5cGUSKgomQVVUSEVOVElDQVRJT05fVFlQRV9ESVJFQ1RPUllfSU'
    'RFTlRJVFkQAA==');

@$core.Deprecated('Use calculationExecutionStateDescriptor instead')
const CalculationExecutionState$json = {
  '1': 'CalculationExecutionState',
  '2': [
    {'1': 'CALCULATION_EXECUTION_STATE_QUEUED', '2': 0},
    {'1': 'CALCULATION_EXECUTION_STATE_RUNNING', '2': 1},
    {'1': 'CALCULATION_EXECUTION_STATE_CANCELED', '2': 2},
    {'1': 'CALCULATION_EXECUTION_STATE_CANCELING', '2': 3},
    {'1': 'CALCULATION_EXECUTION_STATE_CREATING', '2': 4},
    {'1': 'CALCULATION_EXECUTION_STATE_COMPLETED', '2': 5},
    {'1': 'CALCULATION_EXECUTION_STATE_CREATED', '2': 6},
    {'1': 'CALCULATION_EXECUTION_STATE_FAILED', '2': 7},
  ],
};

/// Descriptor for `CalculationExecutionState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List calculationExecutionStateDescriptor = $convert.base64Decode(
    'ChlDYWxjdWxhdGlvbkV4ZWN1dGlvblN0YXRlEiYKIkNBTENVTEFUSU9OX0VYRUNVVElPTl9TVE'
    'FURV9RVUVVRUQQABInCiNDQUxDVUxBVElPTl9FWEVDVVRJT05fU1RBVEVfUlVOTklORxABEigK'
    'JENBTENVTEFUSU9OX0VYRUNVVElPTl9TVEFURV9DQU5DRUxFRBACEikKJUNBTENVTEFUSU9OX0'
    'VYRUNVVElPTl9TVEFURV9DQU5DRUxJTkcQAxIoCiRDQUxDVUxBVElPTl9FWEVDVVRJT05fU1RB'
    'VEVfQ1JFQVRJTkcQBBIpCiVDQUxDVUxBVElPTl9FWEVDVVRJT05fU1RBVEVfQ09NUExFVEVEEA'
    'USJwojQ0FMQ1VMQVRJT05fRVhFQ1VUSU9OX1NUQVRFX0NSRUFURUQQBhImCiJDQUxDVUxBVElP'
    'Tl9FWEVDVVRJT05fU1RBVEVfRkFJTEVEEAc=');

@$core.Deprecated('Use capacityAllocationStatusDescriptor instead')
const CapacityAllocationStatus$json = {
  '1': 'CapacityAllocationStatus',
  '2': [
    {'1': 'CAPACITY_ALLOCATION_STATUS_PENDING', '2': 0},
    {'1': 'CAPACITY_ALLOCATION_STATUS_SUCCEEDED', '2': 1},
    {'1': 'CAPACITY_ALLOCATION_STATUS_FAILED', '2': 2},
  ],
};

/// Descriptor for `CapacityAllocationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List capacityAllocationStatusDescriptor = $convert.base64Decode(
    'ChhDYXBhY2l0eUFsbG9jYXRpb25TdGF0dXMSJgoiQ0FQQUNJVFlfQUxMT0NBVElPTl9TVEFUVV'
    'NfUEVORElORxAAEigKJENBUEFDSVRZX0FMTE9DQVRJT05fU1RBVFVTX1NVQ0NFRURFRBABEiUK'
    'IUNBUEFDSVRZX0FMTE9DQVRJT05fU1RBVFVTX0ZBSUxFRBAC');

@$core.Deprecated('Use capacityReservationStatusDescriptor instead')
const CapacityReservationStatus$json = {
  '1': 'CapacityReservationStatus',
  '2': [
    {'1': 'CAPACITY_RESERVATION_STATUS_PENDING', '2': 0},
    {'1': 'CAPACITY_RESERVATION_STATUS_UPDATE_PENDING', '2': 1},
    {'1': 'CAPACITY_RESERVATION_STATUS_CANCELLED', '2': 2},
    {'1': 'CAPACITY_RESERVATION_STATUS_ACTIVE', '2': 3},
    {'1': 'CAPACITY_RESERVATION_STATUS_CANCELLING', '2': 4},
    {'1': 'CAPACITY_RESERVATION_STATUS_FAILED', '2': 5},
  ],
};

/// Descriptor for `CapacityReservationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List capacityReservationStatusDescriptor = $convert.base64Decode(
    'ChlDYXBhY2l0eVJlc2VydmF0aW9uU3RhdHVzEicKI0NBUEFDSVRZX1JFU0VSVkFUSU9OX1NUQV'
    'RVU19QRU5ESU5HEAASLgoqQ0FQQUNJVFlfUkVTRVJWQVRJT05fU1RBVFVTX1VQREFURV9QRU5E'
    'SU5HEAESKQolQ0FQQUNJVFlfUkVTRVJWQVRJT05fU1RBVFVTX0NBTkNFTExFRBACEiYKIkNBUE'
    'FDSVRZX1JFU0VSVkFUSU9OX1NUQVRVU19BQ1RJVkUQAxIqCiZDQVBBQ0lUWV9SRVNFUlZBVElP'
    'Tl9TVEFUVVNfQ0FOQ0VMTElORxAEEiYKIkNBUEFDSVRZX1JFU0VSVkFUSU9OX1NUQVRVU19GQU'
    'lMRUQQBQ==');

@$core.Deprecated('Use columnNullableDescriptor instead')
const ColumnNullable$json = {
  '1': 'ColumnNullable',
  '2': [
    {'1': 'COLUMN_NULLABLE_NOT_NULL', '2': 0},
    {'1': 'COLUMN_NULLABLE_UNKNOWN', '2': 1},
    {'1': 'COLUMN_NULLABLE_NULLABLE', '2': 2},
  ],
};

/// Descriptor for `ColumnNullable`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List columnNullableDescriptor = $convert.base64Decode(
    'Cg5Db2x1bW5OdWxsYWJsZRIcChhDT0xVTU5fTlVMTEFCTEVfTk9UX05VTEwQABIbChdDT0xVTU'
    '5fTlVMTEFCTEVfVU5LTk9XThABEhwKGENPTFVNTl9OVUxMQUJMRV9OVUxMQUJMRRAC');

@$core.Deprecated('Use connectionTypeDescriptor instead')
const ConnectionType$json = {
  '1': 'ConnectionType',
  '2': [
    {'1': 'CONNECTION_TYPE_DB2AS400', '2': 0},
    {'1': 'CONNECTION_TYPE_DB2', '2': 1},
    {'1': 'CONNECTION_TYPE_SQLSERVER', '2': 2},
    {'1': 'CONNECTION_TYPE_REDSHIFT', '2': 3},
    {'1': 'CONNECTION_TYPE_MYSQL', '2': 4},
    {'1': 'CONNECTION_TYPE_OPENSEARCH', '2': 5},
    {'1': 'CONNECTION_TYPE_DYNAMODB', '2': 6},
    {'1': 'CONNECTION_TYPE_CMDB', '2': 7},
    {'1': 'CONNECTION_TYPE_DATALAKEGEN2', '2': 8},
    {'1': 'CONNECTION_TYPE_DOCUMENTDB', '2': 9},
    {'1': 'CONNECTION_TYPE_GOOGLECLOUDSTORAGE', '2': 10},
    {'1': 'CONNECTION_TYPE_SNOWFLAKE', '2': 11},
    {'1': 'CONNECTION_TYPE_SAPHANA', '2': 12},
    {'1': 'CONNECTION_TYPE_TIMESTREAM', '2': 13},
    {'1': 'CONNECTION_TYPE_ORACLE', '2': 14},
    {'1': 'CONNECTION_TYPE_BIGQUERY', '2': 15},
    {'1': 'CONNECTION_TYPE_TPCDS', '2': 16},
    {'1': 'CONNECTION_TYPE_SYNAPSE', '2': 17},
    {'1': 'CONNECTION_TYPE_HBASE', '2': 18},
    {'1': 'CONNECTION_TYPE_POSTGRESQL', '2': 19},
  ],
};

/// Descriptor for `ConnectionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionTypeDescriptor = $convert.base64Decode(
    'Cg5Db25uZWN0aW9uVHlwZRIcChhDT05ORUNUSU9OX1RZUEVfREIyQVM0MDAQABIXChNDT05ORU'
    'NUSU9OX1RZUEVfREIyEAESHQoZQ09OTkVDVElPTl9UWVBFX1NRTFNFUlZFUhACEhwKGENPTk5F'
    'Q1RJT05fVFlQRV9SRURTSElGVBADEhkKFUNPTk5FQ1RJT05fVFlQRV9NWVNRTBAEEh4KGkNPTk'
    '5FQ1RJT05fVFlQRV9PUEVOU0VBUkNIEAUSHAoYQ09OTkVDVElPTl9UWVBFX0RZTkFNT0RCEAYS'
    'GAoUQ09OTkVDVElPTl9UWVBFX0NNREIQBxIgChxDT05ORUNUSU9OX1RZUEVfREFUQUxBS0VHRU'
    '4yEAgSHgoaQ09OTkVDVElPTl9UWVBFX0RPQ1VNRU5UREIQCRImCiJDT05ORUNUSU9OX1RZUEVf'
    'R09PR0xFQ0xPVURTVE9SQUdFEAoSHQoZQ09OTkVDVElPTl9UWVBFX1NOT1dGTEFLRRALEhsKF0'
    'NPTk5FQ1RJT05fVFlQRV9TQVBIQU5BEAwSHgoaQ09OTkVDVElPTl9UWVBFX1RJTUVTVFJFQU0Q'
    'DRIaChZDT05ORUNUSU9OX1RZUEVfT1JBQ0xFEA4SHAoYQ09OTkVDVElPTl9UWVBFX0JJR1FVRV'
    'JZEA8SGQoVQ09OTkVDVElPTl9UWVBFX1RQQ0RTEBASGwoXQ09OTkVDVElPTl9UWVBFX1NZTkFQ'
    'U0UQERIZChVDT05ORUNUSU9OX1RZUEVfSEJBU0UQEhIeChpDT05ORUNUSU9OX1RZUEVfUE9TVE'
    'dSRVNRTBAT');

@$core.Deprecated('Use dataCatalogStatusDescriptor instead')
const DataCatalogStatus$json = {
  '1': 'DataCatalogStatus',
  '2': [
    {'1': 'DATA_CATALOG_STATUS_DELETE_FAILED', '2': 0},
    {'1': 'DATA_CATALOG_STATUS_DELETE_COMPLETE', '2': 1},
    {'1': 'DATA_CATALOG_STATUS_CREATE_COMPLETE', '2': 2},
    {'1': 'DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_FAILED', '2': 3},
    {'1': 'DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_IN_PROGRESS', '2': 4},
    {'1': 'DATA_CATALOG_STATUS_CREATE_IN_PROGRESS', '2': 5},
    {'1': 'DATA_CATALOG_STATUS_CREATE_FAILED', '2': 6},
    {'1': 'DATA_CATALOG_STATUS_DELETE_IN_PROGRESS', '2': 7},
    {'1': 'DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_COMPLETE', '2': 8},
  ],
};

/// Descriptor for `DataCatalogStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dataCatalogStatusDescriptor = $convert.base64Decode(
    'ChFEYXRhQ2F0YWxvZ1N0YXR1cxIlCiFEQVRBX0NBVEFMT0dfU1RBVFVTX0RFTEVURV9GQUlMRU'
    'QQABInCiNEQVRBX0NBVEFMT0dfU1RBVFVTX0RFTEVURV9DT01QTEVURRABEicKI0RBVEFfQ0FU'
    'QUxPR19TVEFUVVNfQ1JFQVRFX0NPTVBMRVRFEAISNAowREFUQV9DQVRBTE9HX1NUQVRVU19DUk'
    'VBVEVfRkFJTEVEX0NMRUFOVVBfRkFJTEVEEAMSOQo1REFUQV9DQVRBTE9HX1NUQVRVU19DUkVB'
    'VEVfRkFJTEVEX0NMRUFOVVBfSU5fUFJPR1JFU1MQBBIqCiZEQVRBX0NBVEFMT0dfU1RBVFVTX0'
    'NSRUFURV9JTl9QUk9HUkVTUxAFEiUKIURBVEFfQ0FUQUxPR19TVEFUVVNfQ1JFQVRFX0ZBSUxF'
    'RBAGEioKJkRBVEFfQ0FUQUxPR19TVEFUVVNfREVMRVRFX0lOX1BST0dSRVNTEAcSNgoyREFUQV'
    '9DQVRBTE9HX1NUQVRVU19DUkVBVEVfRkFJTEVEX0NMRUFOVVBfQ09NUExFVEUQCA==');

@$core.Deprecated('Use dataCatalogTypeDescriptor instead')
const DataCatalogType$json = {
  '1': 'DataCatalogType',
  '2': [
    {'1': 'DATA_CATALOG_TYPE_FEDERATED', '2': 0},
    {'1': 'DATA_CATALOG_TYPE_GLUE', '2': 1},
    {'1': 'DATA_CATALOG_TYPE_LAMBDA', '2': 2},
    {'1': 'DATA_CATALOG_TYPE_HIVE', '2': 3},
  ],
};

/// Descriptor for `DataCatalogType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dataCatalogTypeDescriptor = $convert.base64Decode(
    'Cg9EYXRhQ2F0YWxvZ1R5cGUSHwobREFUQV9DQVRBTE9HX1RZUEVfRkVERVJBVEVEEAASGgoWRE'
    'FUQV9DQVRBTE9HX1RZUEVfR0xVRRABEhwKGERBVEFfQ0FUQUxPR19UWVBFX0xBTUJEQRACEhoK'
    'FkRBVEFfQ0FUQUxPR19UWVBFX0hJVkUQAw==');

@$core.Deprecated('Use encryptionOptionDescriptor instead')
const EncryptionOption$json = {
  '1': 'EncryptionOption',
  '2': [
    {'1': 'ENCRYPTION_OPTION_SSE_S3', '2': 0},
    {'1': 'ENCRYPTION_OPTION_CSE_KMS', '2': 1},
    {'1': 'ENCRYPTION_OPTION_SSE_KMS', '2': 2},
  ],
};

/// Descriptor for `EncryptionOption`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List encryptionOptionDescriptor = $convert.base64Decode(
    'ChBFbmNyeXB0aW9uT3B0aW9uEhwKGEVOQ1JZUFRJT05fT1BUSU9OX1NTRV9TMxAAEh0KGUVOQ1'
    'JZUFRJT05fT1BUSU9OX0NTRV9LTVMQARIdChlFTkNSWVBUSU9OX09QVElPTl9TU0VfS01TEAI=');

@$core.Deprecated('Use executorStateDescriptor instead')
const ExecutorState$json = {
  '1': 'ExecutorState',
  '2': [
    {'1': 'EXECUTOR_STATE_TERMINATING', '2': 0},
    {'1': 'EXECUTOR_STATE_TERMINATED', '2': 1},
    {'1': 'EXECUTOR_STATE_REGISTERED', '2': 2},
    {'1': 'EXECUTOR_STATE_CREATING', '2': 3},
    {'1': 'EXECUTOR_STATE_CREATED', '2': 4},
    {'1': 'EXECUTOR_STATE_FAILED', '2': 5},
  ],
};

/// Descriptor for `ExecutorState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executorStateDescriptor = $convert.base64Decode(
    'Cg1FeGVjdXRvclN0YXRlEh4KGkVYRUNVVE9SX1NUQVRFX1RFUk1JTkFUSU5HEAASHQoZRVhFQ1'
    'VUT1JfU1RBVEVfVEVSTUlOQVRFRBABEh0KGUVYRUNVVE9SX1NUQVRFX1JFR0lTVEVSRUQQAhIb'
    'ChdFWEVDVVRPUl9TVEFURV9DUkVBVElORxADEhoKFkVYRUNVVE9SX1NUQVRFX0NSRUFURUQQBB'
    'IZChVFWEVDVVRPUl9TVEFURV9GQUlMRUQQBQ==');

@$core.Deprecated('Use executorTypeDescriptor instead')
const ExecutorType$json = {
  '1': 'ExecutorType',
  '2': [
    {'1': 'EXECUTOR_TYPE_GATEWAY', '2': 0},
    {'1': 'EXECUTOR_TYPE_WORKER', '2': 1},
    {'1': 'EXECUTOR_TYPE_COORDINATOR', '2': 2},
  ],
};

/// Descriptor for `ExecutorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executorTypeDescriptor = $convert.base64Decode(
    'CgxFeGVjdXRvclR5cGUSGQoVRVhFQ1VUT1JfVFlQRV9HQVRFV0FZEAASGAoURVhFQ1VUT1JfVF'
    'lQRV9XT1JLRVIQARIdChlFWEVDVVRPUl9UWVBFX0NPT1JESU5BVE9SEAI=');

@$core.Deprecated('Use notebookTypeDescriptor instead')
const NotebookType$json = {
  '1': 'NotebookType',
  '2': [
    {'1': 'NOTEBOOK_TYPE_IPYNB', '2': 0},
  ],
};

/// Descriptor for `NotebookType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List notebookTypeDescriptor = $convert
    .base64Decode('CgxOb3RlYm9va1R5cGUSFwoTTk9URUJPT0tfVFlQRV9JUFlOQhAA');

@$core.Deprecated('Use queryExecutionStateDescriptor instead')
const QueryExecutionState$json = {
  '1': 'QueryExecutionState',
  '2': [
    {'1': 'QUERY_EXECUTION_STATE_QUEUED', '2': 0},
    {'1': 'QUERY_EXECUTION_STATE_RUNNING', '2': 1},
    {'1': 'QUERY_EXECUTION_STATE_SUCCEEDED', '2': 2},
    {'1': 'QUERY_EXECUTION_STATE_CANCELLED', '2': 3},
    {'1': 'QUERY_EXECUTION_STATE_FAILED', '2': 4},
  ],
};

/// Descriptor for `QueryExecutionState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queryExecutionStateDescriptor = $convert.base64Decode(
    'ChNRdWVyeUV4ZWN1dGlvblN0YXRlEiAKHFFVRVJZX0VYRUNVVElPTl9TVEFURV9RVUVVRUQQAB'
    'IhCh1RVUVSWV9FWEVDVVRJT05fU1RBVEVfUlVOTklORxABEiMKH1FVRVJZX0VYRUNVVElPTl9T'
    'VEFURV9TVUNDRUVERUQQAhIjCh9RVUVSWV9FWEVDVVRJT05fU1RBVEVfQ0FOQ0VMTEVEEAMSIA'
    'ocUVVFUllfRVhFQ1VUSU9OX1NUQVRFX0ZBSUxFRBAE');

@$core.Deprecated('Use queryResultTypeDescriptor instead')
const QueryResultType$json = {
  '1': 'QueryResultType',
  '2': [
    {'1': 'QUERY_RESULT_TYPE_DATA_MANIFEST', '2': 0},
    {'1': 'QUERY_RESULT_TYPE_DATA_ROWS', '2': 1},
  ],
};

/// Descriptor for `QueryResultType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queryResultTypeDescriptor = $convert.base64Decode(
    'Cg9RdWVyeVJlc3VsdFR5cGUSIwofUVVFUllfUkVTVUxUX1RZUEVfREFUQV9NQU5JRkVTVBAAEh'
    '8KG1FVRVJZX1JFU1VMVF9UWVBFX0RBVEFfUk9XUxAB');

@$core.Deprecated('Use s3AclOptionDescriptor instead')
const S3AclOption$json = {
  '1': 'S3AclOption',
  '2': [
    {'1': 'S3_ACL_OPTION_BUCKET_OWNER_FULL_CONTROL', '2': 0},
  ],
};

/// Descriptor for `S3AclOption`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List s3AclOptionDescriptor = $convert.base64Decode(
    'CgtTM0FjbE9wdGlvbhIrCidTM19BQ0xfT1BUSU9OX0JVQ0tFVF9PV05FUl9GVUxMX0NPTlRST0'
    'wQAA==');

@$core.Deprecated('Use sessionStateDescriptor instead')
const SessionState$json = {
  '1': 'SessionState',
  '2': [
    {'1': 'SESSION_STATE_TERMINATING', '2': 0},
    {'1': 'SESSION_STATE_TERMINATED', '2': 1},
    {'1': 'SESSION_STATE_DEGRADED', '2': 2},
    {'1': 'SESSION_STATE_BUSY', '2': 3},
    {'1': 'SESSION_STATE_IDLE', '2': 4},
    {'1': 'SESSION_STATE_CREATING', '2': 5},
    {'1': 'SESSION_STATE_CREATED', '2': 6},
    {'1': 'SESSION_STATE_FAILED', '2': 7},
  ],
};

/// Descriptor for `SessionState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sessionStateDescriptor = $convert.base64Decode(
    'CgxTZXNzaW9uU3RhdGUSHQoZU0VTU0lPTl9TVEFURV9URVJNSU5BVElORxAAEhwKGFNFU1NJT0'
    '5fU1RBVEVfVEVSTUlOQVRFRBABEhoKFlNFU1NJT05fU1RBVEVfREVHUkFERUQQAhIWChJTRVNT'
    'SU9OX1NUQVRFX0JVU1kQAxIWChJTRVNTSU9OX1NUQVRFX0lETEUQBBIaChZTRVNTSU9OX1NUQV'
    'RFX0NSRUFUSU5HEAUSGQoVU0VTU0lPTl9TVEFURV9DUkVBVEVEEAYSGAoUU0VTU0lPTl9TVEFU'
    'RV9GQUlMRUQQBw==');

@$core.Deprecated('Use statementTypeDescriptor instead')
const StatementType$json = {
  '1': 'StatementType',
  '2': [
    {'1': 'STATEMENT_TYPE_UTILITY', '2': 0},
    {'1': 'STATEMENT_TYPE_DDL', '2': 1},
    {'1': 'STATEMENT_TYPE_DML', '2': 2},
  ],
};

/// Descriptor for `StatementType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List statementTypeDescriptor = $convert.base64Decode(
    'Cg1TdGF0ZW1lbnRUeXBlEhoKFlNUQVRFTUVOVF9UWVBFX1VUSUxJVFkQABIWChJTVEFURU1FTl'
    'RfVFlQRV9EREwQARIWChJTVEFURU1FTlRfVFlQRV9ETUwQAg==');

@$core.Deprecated('Use throttleReasonDescriptor instead')
const ThrottleReason$json = {
  '1': 'ThrottleReason',
  '2': [
    {'1': 'THROTTLE_REASON_CONCURRENT_QUERY_LIMIT_EXCEEDED', '2': 0},
  ],
};

/// Descriptor for `ThrottleReason`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List throttleReasonDescriptor = $convert.base64Decode(
    'Cg5UaHJvdHRsZVJlYXNvbhIzCi9USFJPVFRMRV9SRUFTT05fQ09OQ1VSUkVOVF9RVUVSWV9MSU'
    '1JVF9FWENFRURFRBAA');

@$core.Deprecated('Use workGroupStateDescriptor instead')
const WorkGroupState$json = {
  '1': 'WorkGroupState',
  '2': [
    {'1': 'WORK_GROUP_STATE_DISABLED', '2': 0},
    {'1': 'WORK_GROUP_STATE_ENABLED', '2': 1},
  ],
};

/// Descriptor for `WorkGroupState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List workGroupStateDescriptor = $convert.base64Decode(
    'Cg5Xb3JrR3JvdXBTdGF0ZRIdChlXT1JLX0dST1VQX1NUQVRFX0RJU0FCTEVEEAASHAoYV09SS1'
    '9HUk9VUF9TVEFURV9FTkFCTEVEEAE=');

@$core.Deprecated('Use aclConfigurationDescriptor instead')
const AclConfiguration$json = {
  '1': 'AclConfiguration',
  '2': [
    {
      '1': 's3acloption',
      '3': 97327001,
      '4': 1,
      '5': 14,
      '6': '.athena.S3AclOption',
      '10': 's3acloption'
    },
  ],
};

/// Descriptor for `AclConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aclConfigurationDescriptor = $convert.base64Decode(
    'ChBBY2xDb25maWd1cmF0aW9uEjgKC3MzYWNsb3B0aW9uGJmvtC4gASgOMhMuYXRoZW5hLlMzQW'
    'NsT3B0aW9uUgtzM2FjbG9wdGlvbg==');

@$core.Deprecated('Use applicationDPUSizesDescriptor instead')
const ApplicationDPUSizes$json = {
  '1': 'ApplicationDPUSizes',
  '2': [
    {
      '1': 'applicationruntimeid',
      '3': 300478599,
      '4': 1,
      '5': 9,
      '10': 'applicationruntimeid'
    },
    {
      '1': 'supporteddpusizes',
      '3': 227874103,
      '4': 3,
      '5': 5,
      '10': 'supporteddpusizes'
    },
  ],
};

/// Descriptor for `ApplicationDPUSizes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List applicationDPUSizesDescriptor = $convert.base64Decode(
    'ChNBcHBsaWNhdGlvbkRQVVNpemVzEjYKFGFwcGxpY2F0aW9ucnVudGltZWlkGIfho48BIAEoCV'
    'IUYXBwbGljYXRpb25ydW50aW1laWQSLwoRc3VwcG9ydGVkZHB1c2l6ZXMYt6rUbCADKAVSEXN1'
    'cHBvcnRlZGRwdXNpemVz');

@$core.Deprecated('Use athenaErrorDescriptor instead')
const AthenaError$json = {
  '1': 'AthenaError',
  '2': [
    {
      '1': 'errorcategory',
      '3': 315958414,
      '4': 1,
      '5': 5,
      '10': 'errorcategory'
    },
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'errortype', '3': 398848954, '4': 1, '5': 5, '10': 'errortype'},
    {'1': 'retryable', '3': 83386186, '4': 1, '5': 8, '10': 'retryable'},
  ],
};

/// Descriptor for `AthenaError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List athenaErrorDescriptor = $convert.base64Decode(
    'CgtBdGhlbmFFcnJvchIoCg1lcnJvcmNhdGVnb3J5GI7J1JYBIAEoBVINZXJyb3JjYXRlZ29yeR'
    'ImCgxlcnJvcm1lc3NhZ2UYqYqr9wEgASgJUgxlcnJvcm1lc3NhZ2USIAoJZXJyb3J0eXBlGLrn'
    'l74BIAEoBVIJZXJyb3J0eXBlEh8KCXJldHJ5YWJsZRjKvuEnIAEoCFIJcmV0cnlhYmxl');

@$core.Deprecated('Use batchGetNamedQueryInputDescriptor instead')
const BatchGetNamedQueryInput$json = {
  '1': 'BatchGetNamedQueryInput',
  '2': [
    {'1': 'namedqueryids', '3': 6092797, '4': 3, '5': 9, '10': 'namedqueryids'},
  ],
};

/// Descriptor for `BatchGetNamedQueryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetNamedQueryInputDescriptor =
    $convert.base64Decode(
        'ChdCYXRjaEdldE5hbWVkUXVlcnlJbnB1dBInCg1uYW1lZHF1ZXJ5aWRzGP3v8wIgAygJUg1uYW'
        '1lZHF1ZXJ5aWRz');

@$core.Deprecated('Use batchGetNamedQueryOutputDescriptor instead')
const BatchGetNamedQueryOutput$json = {
  '1': 'BatchGetNamedQueryOutput',
  '2': [
    {
      '1': 'namedqueries',
      '3': 384985627,
      '4': 3,
      '5': 11,
      '6': '.athena.NamedQuery',
      '10': 'namedqueries'
    },
    {
      '1': 'unprocessednamedqueryids',
      '3': 529659590,
      '4': 3,
      '5': 11,
      '6': '.athena.UnprocessedNamedQueryId',
      '10': 'unprocessednamedqueryids'
    },
  ],
};

/// Descriptor for `BatchGetNamedQueryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetNamedQueryOutputDescriptor = $convert.base64Decode(
    'ChhCYXRjaEdldE5hbWVkUXVlcnlPdXRwdXQSOgoMbmFtZWRxdWVyaWVzGJvUybcBIAMoCzISLm'
    'F0aGVuYS5OYW1lZFF1ZXJ5UgxuYW1lZHF1ZXJpZXMSXwoYdW5wcm9jZXNzZWRuYW1lZHF1ZXJ5'
    'aWRzGMbtx/wBIAMoCzIfLmF0aGVuYS5VbnByb2Nlc3NlZE5hbWVkUXVlcnlJZFIYdW5wcm9jZX'
    'NzZWRuYW1lZHF1ZXJ5aWRz');

@$core.Deprecated('Use batchGetPreparedStatementInputDescriptor instead')
const BatchGetPreparedStatementInput$json = {
  '1': 'BatchGetPreparedStatementInput',
  '2': [
    {
      '1': 'preparedstatementnames',
      '3': 352693824,
      '4': 3,
      '5': 9,
      '10': 'preparedstatementnames'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `BatchGetPreparedStatementInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetPreparedStatementInputDescriptor =
    $convert.base64Decode(
        'Ch5CYXRjaEdldFByZXBhcmVkU3RhdGVtZW50SW5wdXQSOgoWcHJlcGFyZWRzdGF0ZW1lbnRuYW'
        '1lcxjA3JaoASADKAlSFnByZXBhcmVkc3RhdGVtZW50bmFtZXMSIAoJd29ya2dyb3VwGIStofEB'
        'IAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use batchGetPreparedStatementOutputDescriptor instead')
const BatchGetPreparedStatementOutput$json = {
  '1': 'BatchGetPreparedStatementOutput',
  '2': [
    {
      '1': 'preparedstatements',
      '3': 526923667,
      '4': 3,
      '5': 11,
      '6': '.athena.PreparedStatement',
      '10': 'preparedstatements'
    },
    {
      '1': 'unprocessedpreparedstatementnames',
      '3': 526426761,
      '4': 3,
      '5': 11,
      '6': '.athena.UnprocessedPreparedStatementName',
      '10': 'unprocessedpreparedstatementnames'
    },
  ],
};

/// Descriptor for `BatchGetPreparedStatementOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetPreparedStatementOutputDescriptor = $convert.base64Decode(
    'Ch9CYXRjaEdldFByZXBhcmVkU3RhdGVtZW50T3V0cHV0Ek0KEnByZXBhcmVkc3RhdGVtZW50cx'
    'iT76D7ASADKAsyGS5hdGhlbmEuUHJlcGFyZWRTdGF0ZW1lbnRSEnByZXBhcmVkc3RhdGVtZW50'
    'cxJ6CiF1bnByb2Nlc3NlZHByZXBhcmVkc3RhdGVtZW50bmFtZXMYicWC+wEgAygLMiguYXRoZW'
    '5hLlVucHJvY2Vzc2VkUHJlcGFyZWRTdGF0ZW1lbnROYW1lUiF1bnByb2Nlc3NlZHByZXBhcmVk'
    'c3RhdGVtZW50bmFtZXM=');

@$core.Deprecated('Use batchGetQueryExecutionInputDescriptor instead')
const BatchGetQueryExecutionInput$json = {
  '1': 'BatchGetQueryExecutionInput',
  '2': [
    {
      '1': 'queryexecutionids',
      '3': 493941192,
      '4': 3,
      '5': 9,
      '10': 'queryexecutionids'
    },
  ],
};

/// Descriptor for `BatchGetQueryExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetQueryExecutionInputDescriptor =
    $convert.base64Decode(
        'ChtCYXRjaEdldFF1ZXJ5RXhlY3V0aW9uSW5wdXQSMAoRcXVlcnlleGVjdXRpb25pZHMYyOPD6w'
        'EgAygJUhFxdWVyeWV4ZWN1dGlvbmlkcw==');

@$core.Deprecated('Use batchGetQueryExecutionOutputDescriptor instead')
const BatchGetQueryExecutionOutput$json = {
  '1': 'BatchGetQueryExecutionOutput',
  '2': [
    {
      '1': 'queryexecutions',
      '3': 91405361,
      '4': 3,
      '5': 11,
      '6': '.athena.QueryExecution',
      '10': 'queryexecutions'
    },
    {
      '1': 'unprocessedqueryexecutionids',
      '3': 85288387,
      '4': 3,
      '5': 11,
      '6': '.athena.UnprocessedQueryExecutionId',
      '10': 'unprocessedqueryexecutionids'
    },
  ],
};

/// Descriptor for `BatchGetQueryExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetQueryExecutionOutputDescriptor = $convert.base64Decode(
    'ChxCYXRjaEdldFF1ZXJ5RXhlY3V0aW9uT3V0cHV0EkMKD3F1ZXJ5ZXhlY3V0aW9ucxix+MorIA'
    'MoCzIWLmF0aGVuYS5RdWVyeUV4ZWN1dGlvblIPcXVlcnlleGVjdXRpb25zEmoKHHVucHJvY2Vz'
    'c2VkcXVlcnlleGVjdXRpb25pZHMYw8vVKCADKAsyIy5hdGhlbmEuVW5wcm9jZXNzZWRRdWVyeU'
    'V4ZWN1dGlvbklkUhx1bnByb2Nlc3NlZHF1ZXJ5ZXhlY3V0aW9uaWRz');

@$core.Deprecated('Use calculationConfigurationDescriptor instead')
const CalculationConfiguration$json = {
  '1': 'CalculationConfiguration',
  '2': [
    {'1': 'codeblock', '3': 23945838, '4': 1, '5': 9, '10': 'codeblock'},
  ],
};

/// Descriptor for `CalculationConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculationConfigurationDescriptor =
    $convert.base64Decode(
        'ChhDYWxjdWxhdGlvbkNvbmZpZ3VyYXRpb24SHwoJY29kZWJsb2NrGO7EtQsgASgJUgljb2RlYm'
        'xvY2s=');

@$core.Deprecated('Use calculationResultDescriptor instead')
const CalculationResult$json = {
  '1': 'CalculationResult',
  '2': [
    {'1': 'results3uri', '3': 120920277, '4': 1, '5': 9, '10': 'results3uri'},
    {'1': 'resulttype', '3': 261078637, '4': 1, '5': 9, '10': 'resulttype'},
    {
      '1': 'stderrors3uri',
      '3': 227945541,
      '4': 1,
      '5': 9,
      '10': 'stderrors3uri'
    },
    {'1': 'stdouts3uri', '3': 100477551, '4': 1, '5': 9, '10': 'stdouts3uri'},
  ],
};

/// Descriptor for `CalculationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculationResultDescriptor = $convert.base64Decode(
    'ChFDYWxjdWxhdGlvblJlc3VsdBIjCgtyZXN1bHRzM3VyaRjVsdQ5IAEoCVILcmVzdWx0czN1cm'
    'kSIQoKcmVzdWx0dHlwZRjt/L58IAEoCVIKcmVzdWx0dHlwZRInCg1zdGRlcnJvcnMzdXJpGMXY'
    '2GwgASgJUg1zdGRlcnJvcnMzdXJpEiMKC3N0ZG91dHMzdXJpGO/U9C8gASgJUgtzdGRvdXRzM3'
    'VyaQ==');

@$core.Deprecated('Use calculationStatisticsDescriptor instead')
const CalculationStatistics$json = {
  '1': 'CalculationStatistics',
  '2': [
    {
      '1': 'dpuexecutioninmillis',
      '3': 174857936,
      '4': 1,
      '5': 3,
      '10': 'dpuexecutioninmillis'
    },
    {'1': 'progress', '3': 439787879, '4': 1, '5': 9, '10': 'progress'},
  ],
};

/// Descriptor for `CalculationStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculationStatisticsDescriptor = $convert.base64Decode(
    'ChVDYWxjdWxhdGlvblN0YXRpc3RpY3MSNQoUZHB1ZXhlY3V0aW9uaW5taWxsaXMY0L2wUyABKA'
    'NSFGRwdWV4ZWN1dGlvbmlubWlsbGlzEh4KCHByb2dyZXNzGOfC2tEBIAEoCVIIcHJvZ3Jlc3M=');

@$core.Deprecated('Use calculationStatusDescriptor instead')
const CalculationStatus$json = {
  '1': 'CalculationStatus',
  '2': [
    {
      '1': 'completiondatetime',
      '3': 175822779,
      '4': 1,
      '5': 9,
      '10': 'completiondatetime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.CalculationExecutionState',
      '10': 'state'
    },
    {
      '1': 'statechangereason',
      '3': 228940439,
      '4': 1,
      '5': 9,
      '10': 'statechangereason'
    },
    {
      '1': 'submissiondatetime',
      '3': 449650437,
      '4': 1,
      '5': 9,
      '10': 'submissiondatetime'
    },
  ],
};

/// Descriptor for `CalculationStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculationStatusDescriptor = $convert.base64Decode(
    'ChFDYWxjdWxhdGlvblN0YXR1cxIxChJjb21wbGV0aW9uZGF0ZXRpbWUYu6/rUyABKAlSEmNvbX'
    'BsZXRpb25kYXRldGltZRI7CgVzdGF0ZRiXybLvASABKA4yIS5hdGhlbmEuQ2FsY3VsYXRpb25F'
    'eGVjdXRpb25TdGF0ZVIFc3RhdGUSLwoRc3RhdGVjaGFuZ2VyZWFzb24Yl7WVbSABKAlSEXN0YX'
    'RlY2hhbmdlcmVhc29uEjIKEnN1Ym1pc3Npb25kYXRldGltZRiFvrTWASABKAlSEnN1Ym1pc3Np'
    'b25kYXRldGltZQ==');

@$core.Deprecated('Use calculationSummaryDescriptor instead')
const CalculationSummary$json = {
  '1': 'CalculationSummary',
  '2': [
    {
      '1': 'calculationexecutionid',
      '3': 80028050,
      '4': 1,
      '5': 9,
      '10': 'calculationexecutionid'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.athena.CalculationStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `CalculationSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculationSummaryDescriptor = $convert.base64Decode(
    'ChJDYWxjdWxhdGlvblN1bW1hcnkSOQoWY2FsY3VsYXRpb25leGVjdXRpb25pZBiSw5QmIAEoCV'
    'IWY2FsY3VsYXRpb25leGVjdXRpb25pZBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3Jp'
    'cHRpb24SNAoGc3RhdHVzGJDk+wIgASgLMhkuYXRoZW5hLkNhbGN1bGF0aW9uU3RhdHVzUgZzdG'
    'F0dXM=');

@$core.Deprecated('Use cancelCapacityReservationInputDescriptor instead')
const CancelCapacityReservationInput$json = {
  '1': 'CancelCapacityReservationInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CancelCapacityReservationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelCapacityReservationInputDescriptor =
    $convert.base64Decode(
        'Ch5DYW5jZWxDYXBhY2l0eVJlc2VydmF0aW9uSW5wdXQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ'
        '==');

@$core.Deprecated('Use cancelCapacityReservationOutputDescriptor instead')
const CancelCapacityReservationOutput$json = {
  '1': 'CancelCapacityReservationOutput',
};

/// Descriptor for `CancelCapacityReservationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelCapacityReservationOutputDescriptor =
    $convert.base64Decode('Ch9DYW5jZWxDYXBhY2l0eVJlc2VydmF0aW9uT3V0cHV0');

@$core.Deprecated('Use capacityAllocationDescriptor instead')
const CapacityAllocation$json = {
  '1': 'CapacityAllocation',
  '2': [
    {
      '1': 'requestcompletiontime',
      '3': 324037826,
      '4': 1,
      '5': 9,
      '10': 'requestcompletiontime'
    },
    {'1': 'requesttime', '3': 507812148, '4': 1, '5': 9, '10': 'requesttime'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.athena.CapacityAllocationStatus',
      '10': 'status'
    },
    {
      '1': 'statusmessage',
      '3': 72590095,
      '4': 1,
      '5': 9,
      '10': 'statusmessage'
    },
  ],
};

/// Descriptor for `CapacityAllocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityAllocationDescriptor = $convert.base64Decode(
    'ChJDYXBhY2l0eUFsbG9jYXRpb24SOAoVcmVxdWVzdGNvbXBsZXRpb250aW1lGMLZwZoBIAEoCV'
    'IVcmVxdWVzdGNvbXBsZXRpb250aW1lEiQKC3JlcXVlc3R0aW1lGLSykvIBIAEoCVILcmVxdWVz'
    'dHRpbWUSOwoGc3RhdHVzGJDk+wIgASgOMiAuYXRoZW5hLkNhcGFjaXR5QWxsb2NhdGlvblN0YX'
    'R1c1IGc3RhdHVzEicKDXN0YXR1c21lc3NhZ2UYj8bOIiABKAlSDXN0YXR1c21lc3NhZ2U=');

@$core.Deprecated('Use capacityAssignmentDescriptor instead')
const CapacityAssignment$json = {
  '1': 'CapacityAssignment',
  '2': [
    {
      '1': 'workgroupnames',
      '3': 486049026,
      '4': 3,
      '5': 9,
      '10': 'workgroupnames'
    },
  ],
};

/// Descriptor for `CapacityAssignment`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityAssignmentDescriptor = $convert.base64Decode(
    'ChJDYXBhY2l0eUFzc2lnbm1lbnQSKgoOd29ya2dyb3VwbmFtZXMYgori5wEgAygJUg53b3JrZ3'
    'JvdXBuYW1lcw==');

@$core.Deprecated('Use capacityAssignmentConfigurationDescriptor instead')
const CapacityAssignmentConfiguration$json = {
  '1': 'CapacityAssignmentConfiguration',
  '2': [
    {
      '1': 'capacityassignments',
      '3': 345772294,
      '4': 3,
      '5': 11,
      '6': '.athena.CapacityAssignment',
      '10': 'capacityassignments'
    },
    {
      '1': 'capacityreservationname',
      '3': 327567687,
      '4': 1,
      '5': 9,
      '10': 'capacityreservationname'
    },
  ],
};

/// Descriptor for `CapacityAssignmentConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityAssignmentConfigurationDescriptor =
    $convert.base64Decode(
        'Ch9DYXBhY2l0eUFzc2lnbm1lbnRDb25maWd1cmF0aW9uElAKE2NhcGFjaXR5YXNzaWdubWVudH'
        'MYhqLwpAEgAygLMhouYXRoZW5hLkNhcGFjaXR5QXNzaWdubWVudFITY2FwYWNpdHlhc3NpZ25t'
        'ZW50cxI8ChdjYXBhY2l0eXJlc2VydmF0aW9ubmFtZRjHkpmcASABKAlSF2NhcGFjaXR5cmVzZX'
        'J2YXRpb25uYW1l');

@$core.Deprecated('Use capacityReservationDescriptor instead')
const CapacityReservation$json = {
  '1': 'CapacityReservation',
  '2': [
    {
      '1': 'allocateddpus',
      '3': 252958879,
      '4': 1,
      '5': 5,
      '10': 'allocateddpus'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastallocation',
      '3': 274654476,
      '4': 1,
      '5': 11,
      '6': '.athena.CapacityAllocation',
      '10': 'lastallocation'
    },
    {
      '1': 'lastsuccessfulallocationtime',
      '3': 383368641,
      '4': 1,
      '5': 9,
      '10': 'lastsuccessfulallocationtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.athena.CapacityReservationStatus',
      '10': 'status'
    },
    {'1': 'targetdpus', '3': 367520745, '4': 1, '5': 5, '10': 'targetdpus'},
  ],
};

/// Descriptor for `CapacityReservation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityReservationDescriptor = $convert.base64Decode(
    'ChNDYXBhY2l0eVJlc2VydmF0aW9uEicKDWFsbG9jYXRlZGRwdXMYn7HPeCABKAVSDWFsbG9jYX'
    'RlZGRwdXMSJQoMY3JlYXRpb250aW1lGObPqjEgASgJUgxjcmVhdGlvbnRpbWUSRgoObGFzdGFs'
    'bG9jYXRpb24YjMr7ggEgASgLMhouYXRoZW5hLkNhcGFjaXR5QWxsb2NhdGlvblIObGFzdGFsbG'
    '9jYXRpb24SRgocbGFzdHN1Y2Nlc3NmdWxhbGxvY2F0aW9udGltZRjB++a2ASABKAlSHGxhc3Rz'
    'dWNjZXNzZnVsYWxsb2NhdGlvbnRpbWUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRI8CgZzdGF0dX'
    'MYkOT7AiABKA4yIS5hdGhlbmEuQ2FwYWNpdHlSZXNlcnZhdGlvblN0YXR1c1IGc3RhdHVzEiIK'
    'CnRhcmdldGRwdXMY6defrwEgASgFUgp0YXJnZXRkcHVz');

@$core.Deprecated('Use classificationDescriptor instead')
const Classification$json = {
  '1': 'Classification',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'properties',
      '3': 29886973,
      '4': 3,
      '5': 11,
      '6': '.athena.Classification.PropertiesEntry',
      '10': 'properties'
    },
  ],
  '3': [Classification_PropertiesEntry$json],
};

@$core.Deprecated('Use classificationDescriptor instead')
const Classification_PropertiesEntry$json = {
  '1': 'PropertiesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Classification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List classificationDescriptor = $convert.base64Decode(
    'Cg5DbGFzc2lmaWNhdGlvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkkKCnByb3BlcnRpZXMY/Z'
    'OgDiADKAsyJi5hdGhlbmEuQ2xhc3NpZmljYXRpb24uUHJvcGVydGllc0VudHJ5Ugpwcm9wZXJ0'
    'aWVzGj0KD1Byb3BlcnRpZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCV'
    'IFdmFsdWU6AjgB');

@$core.Deprecated('Use cloudWatchLoggingConfigurationDescriptor instead')
const CloudWatchLoggingConfiguration$json = {
  '1': 'CloudWatchLoggingConfiguration',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'loggroup', '3': 148580073, '4': 1, '5': 9, '10': 'loggroup'},
    {
      '1': 'logstreamnameprefix',
      '3': 437213671,
      '4': 1,
      '5': 9,
      '10': 'logstreamnameprefix'
    },
    {
      '1': 'logtypes',
      '3': 491693055,
      '4': 3,
      '5': 11,
      '6': '.athena.CloudWatchLoggingConfiguration.LogtypesEntry',
      '10': 'logtypes'
    },
  ],
  '3': [CloudWatchLoggingConfiguration_LogtypesEntry$json],
};

@$core.Deprecated('Use cloudWatchLoggingConfigurationDescriptor instead')
const CloudWatchLoggingConfiguration_LogtypesEntry$json = {
  '1': 'LogtypesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CloudWatchLoggingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudWatchLoggingConfigurationDescriptor = $convert.base64Decode(
    'Ch5DbG91ZFdhdGNoTG9nZ2luZ0NvbmZpZ3VyYXRpb24SHAoHZW5hYmxlZBi/yJvkASABKAhSB2'
    'VuYWJsZWQSHQoIbG9nZ3JvdXAY6c3sRiABKAlSCGxvZ2dyb3VwEjQKE2xvZ3N0cmVhbW5hbWVw'
    'cmVmaXgY57O90AEgASgJUhNsb2dzdHJlYW1uYW1lcHJlZml4ElQKCGxvZ3R5cGVzGP/HuuoBIA'
    'MoCzI0LmF0aGVuYS5DbG91ZFdhdGNoTG9nZ2luZ0NvbmZpZ3VyYXRpb24uTG9ndHlwZXNFbnRy'
    'eVIIbG9ndHlwZXMaOwoNTG9ndHlwZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZR'
    'gCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use columnDescriptor instead')
const Column$json = {
  '1': 'Column',
  '2': [
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `Column`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List columnDescriptor = $convert.base64Decode(
    'CgZDb2x1bW4SHAoHY29tbWVudBj/v77CASABKAlSB2NvbW1lbnQSFQoEbmFtZRiH5oF/IAEoCV'
    'IEbmFtZRIWCgR0eXBlGO6g14oBIAEoCVIEdHlwZQ==');

@$core.Deprecated('Use columnInfoDescriptor instead')
const ColumnInfo$json = {
  '1': 'ColumnInfo',
  '2': [
    {
      '1': 'casesensitive',
      '3': 258546956,
      '4': 1,
      '5': 8,
      '10': 'casesensitive'
    },
    {'1': 'catalogname', '3': 518825212, '4': 1, '5': 9, '10': 'catalogname'},
    {'1': 'label', '3': 516747934, '4': 1, '5': 9, '10': 'label'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nullable',
      '3': 373261405,
      '4': 1,
      '5': 14,
      '6': '.athena.ColumnNullable',
      '10': 'nullable'
    },
    {'1': 'precision', '3': 110022584, '4': 1, '5': 5, '10': 'precision'},
    {'1': 'scale', '3': 139628050, '4': 1, '5': 5, '10': 'scale'},
    {'1': 'schemaname', '3': 443785942, '4': 1, '5': 9, '10': 'schemaname'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {'1': 'type', '3': 290836590, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `ColumnInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List columnInfoDescriptor = $convert.base64Decode(
    'CgpDb2x1bW5JbmZvEicKDWNhc2VzZW5zaXRpdmUYjLqkeyABKAhSDWNhc2VzZW5zaXRpdmUSJA'
    'oLY2F0YWxvZ25hbWUY/Mmy9wEgASgJUgtjYXRhbG9nbmFtZRIYCgVsYWJlbBie5bP2ASABKAlS'
    'BWxhYmVsEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSNgoIbnVsbGFibGUY3Yj+sQEgASgOMhYuYX'
    'RoZW5hLkNvbHVtbk51bGxhYmxlUghudWxsYWJsZRIfCglwcmVjaXNpb24YuJ+7NCABKAVSCXBy'
    'ZWNpc2lvbhIXCgVzY2FsZRiSnMpCIAEoBVIFc2NhbGUSIgoKc2NoZW1hbmFtZRjWxc7TASABKA'
    'lSCnNjaGVtYW5hbWUSIAoJdGFibGVuYW1lGN3k2oEBIAEoCVIJdGFibGVuYW1lEhYKBHR5cGUY'
    '7qDXigEgASgJUgR0eXBl');

@$core.Deprecated('Use createCapacityReservationInputDescriptor instead')
const CreateCapacityReservationInput$json = {
  '1': 'CreateCapacityReservationInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.athena.Tag',
      '10': 'tags'
    },
    {'1': 'targetdpus', '3': 367520745, '4': 1, '5': 5, '10': 'targetdpus'},
  ],
};

/// Descriptor for `CreateCapacityReservationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCapacityReservationInputDescriptor =
    $convert.base64Decode(
        'Ch5DcmVhdGVDYXBhY2l0eVJlc2VydmF0aW9uSW5wdXQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZR'
        'IjCgR0YWdzGMHB9rUBIAMoCzILLmF0aGVuYS5UYWdSBHRhZ3MSIgoKdGFyZ2V0ZHB1cxjp15+v'
        'ASABKAVSCnRhcmdldGRwdXM=');

@$core.Deprecated('Use createCapacityReservationOutputDescriptor instead')
const CreateCapacityReservationOutput$json = {
  '1': 'CreateCapacityReservationOutput',
};

/// Descriptor for `CreateCapacityReservationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCapacityReservationOutputDescriptor =
    $convert.base64Decode('Ch9DcmVhdGVDYXBhY2l0eVJlc2VydmF0aW9uT3V0cHV0');

@$core.Deprecated('Use createDataCatalogInputDescriptor instead')
const CreateDataCatalogInput$json = {
  '1': 'CreateDataCatalogInput',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.athena.CreateDataCatalogInput.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.athena.Tag',
      '10': 'tags'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.athena.DataCatalogType',
      '10': 'type'
    },
  ],
  '3': [CreateDataCatalogInput_ParametersEntry$json],
};

@$core.Deprecated('Use createDataCatalogInputDescriptor instead')
const CreateDataCatalogInput_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateDataCatalogInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDataCatalogInputDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVEYXRhQ2F0YWxvZ0lucHV0EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcm'
    'lwdGlvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lElIKCnBhcmFtZXRlcnMY+qf+6wEgAygLMi4u'
    'YXRoZW5hLkNyZWF0ZURhdGFDYXRhbG9nSW5wdXQuUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZX'
    'JzEiMKBHRhZ3MYwcH2tQEgAygLMgsuYXRoZW5hLlRhZ1IEdGFncxIvCgR0eXBlGO6g14oBIAEo'
    'DjIXLmF0aGVuYS5EYXRhQ2F0YWxvZ1R5cGVSBHR5cGUaPQoPUGFyYW1ldGVyc0VudHJ5EhAKA2'
    'tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createDataCatalogOutputDescriptor instead')
const CreateDataCatalogOutput$json = {
  '1': 'CreateDataCatalogOutput',
  '2': [
    {
      '1': 'datacatalog',
      '3': 209694649,
      '4': 1,
      '5': 11,
      '6': '.athena.DataCatalog',
      '10': 'datacatalog'
    },
  ],
};

/// Descriptor for `CreateDataCatalogOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDataCatalogOutputDescriptor =
    $convert.base64Decode(
        'ChdDcmVhdGVEYXRhQ2F0YWxvZ091dHB1dBI4CgtkYXRhY2F0YWxvZxi53/5jIAEoCzITLmF0aG'
        'VuYS5EYXRhQ2F0YWxvZ1ILZGF0YWNhdGFsb2c=');

@$core.Deprecated('Use createNamedQueryInputDescriptor instead')
const CreateNamedQueryInput$json = {
  '1': 'CreateNamedQueryInput',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {'1': 'database', '3': 278147289, '4': 1, '5': 9, '10': 'database'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `CreateNamedQueryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createNamedQueryInputDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVOYW1lZFF1ZXJ5SW5wdXQSMgoSY2xpZW50cmVxdWVzdHRva2VuGPHvotkBIAEoCV'
    'ISY2xpZW50cmVxdWVzdHRva2VuEh4KCGRhdGFiYXNlGNnh0IQBIAEoCVIIZGF0YWJhc2USIwoL'
    'ZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbW'
    'USJAoLcXVlcnlzdHJpbmcY58rvzwEgASgJUgtxdWVyeXN0cmluZxIgCgl3b3JrZ3JvdXAYhK2h'
    '8QEgASgJUgl3b3JrZ3JvdXA=');

@$core.Deprecated('Use createNamedQueryOutputDescriptor instead')
const CreateNamedQueryOutput$json = {
  '1': 'CreateNamedQueryOutput',
  '2': [
    {'1': 'namedqueryid', '3': 330896872, '4': 1, '5': 9, '10': 'namedqueryid'},
  ],
};

/// Descriptor for `CreateNamedQueryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createNamedQueryOutputDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVOYW1lZFF1ZXJ5T3V0cHV0EiYKDG5hbWVkcXVlcnlpZBjoq+SdASABKAlSDG5hbW'
        'VkcXVlcnlpZA==');

@$core.Deprecated('Use createNotebookInputDescriptor instead')
const CreateNotebookInput$json = {
  '1': 'CreateNotebookInput',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `CreateNotebookInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createNotebookInputDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVOb3RlYm9va0lucHV0EjIKEmNsaWVudHJlcXVlc3R0b2tlbhjx76LZASABKAlSEm'
    'NsaWVudHJlcXVlc3R0b2tlbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEiAKCXdvcmtncm91cBiE'
    'raHxASABKAlSCXdvcmtncm91cA==');

@$core.Deprecated('Use createNotebookOutputDescriptor instead')
const CreateNotebookOutput$json = {
  '1': 'CreateNotebookOutput',
  '2': [
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
  ],
};

/// Descriptor for `CreateNotebookOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createNotebookOutputDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVOb3RlYm9va091dHB1dBIhCgpub3RlYm9va2lkGN60lUsgASgJUgpub3RlYm9va2'
    'lk');

@$core.Deprecated('Use createPreparedStatementInputDescriptor instead')
const CreatePreparedStatementInput$json = {
  '1': 'CreatePreparedStatementInput',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'querystatement',
      '3': 340852217,
      '4': 1,
      '5': 9,
      '10': 'querystatement'
    },
    {
      '1': 'statementname',
      '3': 23047926,
      '4': 1,
      '5': 9,
      '10': 'statementname'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `CreatePreparedStatementInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPreparedStatementInputDescriptor = $convert.base64Decode(
    'ChxDcmVhdGVQcmVwYXJlZFN0YXRlbWVudElucHV0EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUg'
    'tkZXNjcmlwdGlvbhIqCg5xdWVyeXN0YXRlbWVudBj5+8OiASABKAlSDnF1ZXJ5c3RhdGVtZW50'
    'EicKDXN0YXRlbWVudG5hbWUY9t3+CiABKAlSDXN0YXRlbWVudG5hbWUSIAoJd29ya2dyb3VwGI'
    'StofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use createPreparedStatementOutputDescriptor instead')
const CreatePreparedStatementOutput$json = {
  '1': 'CreatePreparedStatementOutput',
};

/// Descriptor for `CreatePreparedStatementOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPreparedStatementOutputDescriptor =
    $convert.base64Decode('Ch1DcmVhdGVQcmVwYXJlZFN0YXRlbWVudE91dHB1dA==');

@$core.Deprecated('Use createPresignedNotebookUrlRequestDescriptor instead')
const CreatePresignedNotebookUrlRequest$json = {
  '1': 'CreatePresignedNotebookUrlRequest',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `CreatePresignedNotebookUrlRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPresignedNotebookUrlRequestDescriptor =
    $convert.base64Decode(
        'CiFDcmVhdGVQcmVzaWduZWROb3RlYm9va1VybFJlcXVlc3QSHwoJc2Vzc2lvbmlkGLuE5QkgAS'
        'gJUglzZXNzaW9uaWQ=');

@$core.Deprecated('Use createPresignedNotebookUrlResponseDescriptor instead')
const CreatePresignedNotebookUrlResponse$json = {
  '1': 'CreatePresignedNotebookUrlResponse',
  '2': [
    {'1': 'authtoken', '3': 349771493, '4': 1, '5': 9, '10': 'authtoken'},
    {
      '1': 'authtokenexpirationtime',
      '3': 406065689,
      '4': 1,
      '5': 3,
      '10': 'authtokenexpirationtime'
    },
    {'1': 'notebookurl', '3': 454352434, '4': 1, '5': 9, '10': 'notebookurl'},
  ],
};

/// Descriptor for `CreatePresignedNotebookUrlResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPresignedNotebookUrlResponseDescriptor =
    $convert.base64Decode(
        'CiJDcmVhdGVQcmVzaWduZWROb3RlYm9va1VybFJlc3BvbnNlEiAKCWF1dGh0b2tlbhjlreSmAS'
        'ABKAlSCWF1dGh0b2tlbhI8ChdhdXRodG9rZW5leHBpcmF0aW9udGltZRiZpNDBASABKANSF2F1'
        'dGh0b2tlbmV4cGlyYXRpb250aW1lEiQKC25vdGVib29rdXJsGLK809gBIAEoCVILbm90ZWJvb2'
        't1cmw=');

@$core.Deprecated('Use createWorkGroupInputDescriptor instead')
const CreateWorkGroupInput$json = {
  '1': 'CreateWorkGroupInput',
  '2': [
    {
      '1': 'configuration',
      '3': 442426458,
      '4': 1,
      '5': 11,
      '6': '.athena.WorkGroupConfiguration',
      '10': 'configuration'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.athena.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateWorkGroupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWorkGroupInputDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVXb3JrR3JvdXBJbnB1dBJICg1jb25maWd1cmF0aW9uGNrI+9IBIAEoCzIeLmF0aG'
    'VuYS5Xb3JrR3JvdXBDb25maWd1cmF0aW9uUg1jb25maWd1cmF0aW9uEiMKC2Rlc2NyaXB0aW9u'
    'GIr0+TYgASgJUgtkZXNjcmlwdGlvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEiMKBHRhZ3MYwc'
    'H2tQEgAygLMgsuYXRoZW5hLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createWorkGroupOutputDescriptor instead')
const CreateWorkGroupOutput$json = {
  '1': 'CreateWorkGroupOutput',
};

/// Descriptor for `CreateWorkGroupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWorkGroupOutputDescriptor =
    $convert.base64Decode('ChVDcmVhdGVXb3JrR3JvdXBPdXRwdXQ=');

@$core
    .Deprecated('Use customerContentEncryptionConfigurationDescriptor instead')
const CustomerContentEncryptionConfiguration$json = {
  '1': 'CustomerContentEncryptionConfiguration',
  '2': [
    {'1': 'kmskey', '3': 114561194, '4': 1, '5': 9, '10': 'kmskey'},
  ],
};

/// Descriptor for `CustomerContentEncryptionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customerContentEncryptionConfigurationDescriptor =
    $convert.base64Decode(
        'CiZDdXN0b21lckNvbnRlbnRFbmNyeXB0aW9uQ29uZmlndXJhdGlvbhIZCgZrbXNrZXkYqqHQNi'
        'ABKAlSBmttc2tleQ==');

@$core.Deprecated('Use dataCatalogDescriptor instead')
const DataCatalog$json = {
  '1': 'DataCatalog',
  '2': [
    {
      '1': 'connectiontype',
      '3': 489507282,
      '4': 1,
      '5': 14,
      '6': '.athena.ConnectionType',
      '10': 'connectiontype'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'error', '3': 328047858, '4': 1, '5': 9, '10': 'error'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.athena.DataCatalog.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.athena.DataCatalogStatus',
      '10': 'status'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.athena.DataCatalogType',
      '10': 'type'
    },
  ],
  '3': [DataCatalog_ParametersEntry$json],
};

@$core.Deprecated('Use dataCatalogDescriptor instead')
const DataCatalog_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DataCatalog`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataCatalogDescriptor = $convert.base64Decode(
    'CgtEYXRhQ2F0YWxvZxJCCg5jb25uZWN0aW9udHlwZRjSk7XpASABKA4yFi5hdGhlbmEuQ29ubm'
    'VjdGlvblR5cGVSDmNvbm5lY3Rpb250eXBlEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNj'
    'cmlwdGlvbhIYCgVlcnJvchjyubacASABKAlSBWVycm9yEhUKBG5hbWUYh+aBfyABKAlSBG5hbW'
    'USRwoKcGFyYW1ldGVycxj6p/7rASADKAsyIy5hdGhlbmEuRGF0YUNhdGFsb2cuUGFyYW1ldGVy'
    'c0VudHJ5UgpwYXJhbWV0ZXJzEjQKBnN0YXR1cxiQ5PsCIAEoDjIZLmF0aGVuYS5EYXRhQ2F0YW'
    'xvZ1N0YXR1c1IGc3RhdHVzEi8KBHR5cGUY7qDXigEgASgOMhcuYXRoZW5hLkRhdGFDYXRhbG9n'
    'VHlwZVIEdHlwZRo9Cg9QYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdW'
    'UYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use dataCatalogSummaryDescriptor instead')
const DataCatalogSummary$json = {
  '1': 'DataCatalogSummary',
  '2': [
    {'1': 'catalogname', '3': 518825212, '4': 1, '5': 9, '10': 'catalogname'},
    {
      '1': 'connectiontype',
      '3': 489507282,
      '4': 1,
      '5': 14,
      '6': '.athena.ConnectionType',
      '10': 'connectiontype'
    },
    {'1': 'error', '3': 328047858, '4': 1, '5': 9, '10': 'error'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.athena.DataCatalogStatus',
      '10': 'status'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.athena.DataCatalogType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `DataCatalogSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataCatalogSummaryDescriptor = $convert.base64Decode(
    'ChJEYXRhQ2F0YWxvZ1N1bW1hcnkSJAoLY2F0YWxvZ25hbWUY/Mmy9wEgASgJUgtjYXRhbG9nbm'
    'FtZRJCCg5jb25uZWN0aW9udHlwZRjSk7XpASABKA4yFi5hdGhlbmEuQ29ubmVjdGlvblR5cGVS'
    'DmNvbm5lY3Rpb250eXBlEhgKBWVycm9yGPK5tpwBIAEoCVIFZXJyb3ISNAoGc3RhdHVzGJDk+w'
    'IgASgOMhkuYXRoZW5hLkRhdGFDYXRhbG9nU3RhdHVzUgZzdGF0dXMSLwoEdHlwZRjuoNeKASAB'
    'KA4yFy5hdGhlbmEuRGF0YUNhdGFsb2dUeXBlUgR0eXBl');

@$core.Deprecated('Use databaseDescriptor instead')
const Database$json = {
  '1': 'Database',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.athena.Database.ParametersEntry',
      '10': 'parameters'
    },
  ],
  '3': [Database_ParametersEntry$json],
};

@$core.Deprecated('Use databaseDescriptor instead')
const Database_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Database`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List databaseDescriptor = $convert.base64Decode(
    'CghEYXRhYmFzZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SFQoEbmFtZR'
    'iH5oF/IAEoCVIEbmFtZRJECgpwYXJhbWV0ZXJzGPqn/usBIAMoCzIgLmF0aGVuYS5EYXRhYmFz'
    'ZS5QYXJhbWV0ZXJzRW50cnlSCnBhcmFtZXRlcnMaPQoPUGFyYW1ldGVyc0VudHJ5EhAKA2tleR'
    'gBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use datumDescriptor instead')
const Datum$json = {
  '1': 'Datum',
  '2': [
    {'1': 'varcharvalue', '3': 286740796, '4': 1, '5': 9, '10': 'varcharvalue'},
  ],
};

/// Descriptor for `Datum`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List datumDescriptor = $convert.base64Decode(
    'CgVEYXR1bRImCgx2YXJjaGFydmFsdWUYvKLdiAEgASgJUgx2YXJjaGFydmFsdWU=');

@$core.Deprecated('Use deleteCapacityReservationInputDescriptor instead')
const DeleteCapacityReservationInput$json = {
  '1': 'DeleteCapacityReservationInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteCapacityReservationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCapacityReservationInputDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVDYXBhY2l0eVJlc2VydmF0aW9uSW5wdXQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ'
        '==');

@$core.Deprecated('Use deleteCapacityReservationOutputDescriptor instead')
const DeleteCapacityReservationOutput$json = {
  '1': 'DeleteCapacityReservationOutput',
};

/// Descriptor for `DeleteCapacityReservationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCapacityReservationOutputDescriptor =
    $convert.base64Decode('Ch9EZWxldGVDYXBhY2l0eVJlc2VydmF0aW9uT3V0cHV0');

@$core.Deprecated('Use deleteDataCatalogInputDescriptor instead')
const DeleteDataCatalogInput$json = {
  '1': 'DeleteDataCatalogInput',
  '2': [
    {
      '1': 'deletecatalogonly',
      '3': 310631780,
      '4': 1,
      '5': 8,
      '10': 'deletecatalogonly'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteDataCatalogInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDataCatalogInputDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVEYXRhQ2F0YWxvZ0lucHV0EjAKEWRlbGV0ZWNhdGFsb2dvbmx5GOS6j5QBIAEoCF'
        'IRZGVsZXRlY2F0YWxvZ29ubHkSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use deleteDataCatalogOutputDescriptor instead')
const DeleteDataCatalogOutput$json = {
  '1': 'DeleteDataCatalogOutput',
  '2': [
    {
      '1': 'datacatalog',
      '3': 209694649,
      '4': 1,
      '5': 11,
      '6': '.athena.DataCatalog',
      '10': 'datacatalog'
    },
  ],
};

/// Descriptor for `DeleteDataCatalogOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDataCatalogOutputDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVEYXRhQ2F0YWxvZ091dHB1dBI4CgtkYXRhY2F0YWxvZxi53/5jIAEoCzITLmF0aG'
        'VuYS5EYXRhQ2F0YWxvZ1ILZGF0YWNhdGFsb2c=');

@$core.Deprecated('Use deleteNamedQueryInputDescriptor instead')
const DeleteNamedQueryInput$json = {
  '1': 'DeleteNamedQueryInput',
  '2': [
    {'1': 'namedqueryid', '3': 330896872, '4': 1, '5': 9, '10': 'namedqueryid'},
  ],
};

/// Descriptor for `DeleteNamedQueryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteNamedQueryInputDescriptor = $convert.base64Decode(
    'ChVEZWxldGVOYW1lZFF1ZXJ5SW5wdXQSJgoMbmFtZWRxdWVyeWlkGOir5J0BIAEoCVIMbmFtZW'
    'RxdWVyeWlk');

@$core.Deprecated('Use deleteNamedQueryOutputDescriptor instead')
const DeleteNamedQueryOutput$json = {
  '1': 'DeleteNamedQueryOutput',
};

/// Descriptor for `DeleteNamedQueryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteNamedQueryOutputDescriptor =
    $convert.base64Decode('ChZEZWxldGVOYW1lZFF1ZXJ5T3V0cHV0');

@$core.Deprecated('Use deleteNotebookInputDescriptor instead')
const DeleteNotebookInput$json = {
  '1': 'DeleteNotebookInput',
  '2': [
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
  ],
};

/// Descriptor for `DeleteNotebookInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteNotebookInputDescriptor = $convert.base64Decode(
    'ChNEZWxldGVOb3RlYm9va0lucHV0EiEKCm5vdGVib29raWQY3rSVSyABKAlSCm5vdGVib29raW'
    'Q=');

@$core.Deprecated('Use deleteNotebookOutputDescriptor instead')
const DeleteNotebookOutput$json = {
  '1': 'DeleteNotebookOutput',
};

/// Descriptor for `DeleteNotebookOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteNotebookOutputDescriptor =
    $convert.base64Decode('ChREZWxldGVOb3RlYm9va091dHB1dA==');

@$core.Deprecated('Use deletePreparedStatementInputDescriptor instead')
const DeletePreparedStatementInput$json = {
  '1': 'DeletePreparedStatementInput',
  '2': [
    {
      '1': 'statementname',
      '3': 23047926,
      '4': 1,
      '5': 9,
      '10': 'statementname'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `DeletePreparedStatementInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePreparedStatementInputDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVQcmVwYXJlZFN0YXRlbWVudElucHV0EicKDXN0YXRlbWVudG5hbWUY9t3+CiABKA'
        'lSDXN0YXRlbWVudG5hbWUSIAoJd29ya2dyb3VwGIStofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use deletePreparedStatementOutputDescriptor instead')
const DeletePreparedStatementOutput$json = {
  '1': 'DeletePreparedStatementOutput',
};

/// Descriptor for `DeletePreparedStatementOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePreparedStatementOutputDescriptor =
    $convert.base64Decode('Ch1EZWxldGVQcmVwYXJlZFN0YXRlbWVudE91dHB1dA==');

@$core.Deprecated('Use deleteWorkGroupInputDescriptor instead')
const DeleteWorkGroupInput$json = {
  '1': 'DeleteWorkGroupInput',
  '2': [
    {
      '1': 'recursivedeleteoption',
      '3': 398167566,
      '4': 1,
      '5': 8,
      '10': 'recursivedeleteoption'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `DeleteWorkGroupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteWorkGroupInputDescriptor = $convert.base64Decode(
    'ChREZWxldGVXb3JrR3JvdXBJbnB1dBI4ChVyZWN1cnNpdmVkZWxldGVvcHRpb24YjpzuvQEgAS'
    'gIUhVyZWN1cnNpdmVkZWxldGVvcHRpb24SIAoJd29ya2dyb3VwGIStofEBIAEoCVIJd29ya2dy'
    'b3Vw');

@$core.Deprecated('Use deleteWorkGroupOutputDescriptor instead')
const DeleteWorkGroupOutput$json = {
  '1': 'DeleteWorkGroupOutput',
};

/// Descriptor for `DeleteWorkGroupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteWorkGroupOutputDescriptor =
    $convert.base64Decode('ChVEZWxldGVXb3JrR3JvdXBPdXRwdXQ=');

@$core.Deprecated('Use encryptionConfigurationDescriptor instead')
const EncryptionConfiguration$json = {
  '1': 'EncryptionConfiguration',
  '2': [
    {
      '1': 'encryptionoption',
      '3': 160833062,
      '4': 1,
      '5': 14,
      '6': '.athena.EncryptionOption',
      '10': 'encryptionoption'
    },
    {'1': 'kmskey', '3': 114561194, '4': 1, '5': 9, '10': 'kmskey'},
  ],
};

/// Descriptor for `EncryptionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptionConfigurationDescriptor = $convert.base64Decode(
    'ChdFbmNyeXB0aW9uQ29uZmlndXJhdGlvbhJHChBlbmNyeXB0aW9ub3B0aW9uGKa82EwgASgOMh'
    'guYXRoZW5hLkVuY3J5cHRpb25PcHRpb25SEGVuY3J5cHRpb25vcHRpb24SGQoGa21za2V5GKqh'
    '0DYgASgJUgZrbXNrZXk=');

@$core.Deprecated('Use engineConfigurationDescriptor instead')
const EngineConfiguration$json = {
  '1': 'EngineConfiguration',
  '2': [
    {
      '1': 'additionalconfigs',
      '3': 341218716,
      '4': 3,
      '5': 11,
      '6': '.athena.EngineConfiguration.AdditionalconfigsEntry',
      '10': 'additionalconfigs'
    },
    {
      '1': 'classifications',
      '3': 129400407,
      '4': 3,
      '5': 11,
      '6': '.athena.Classification',
      '10': 'classifications'
    },
    {
      '1': 'coordinatordpusize',
      '3': 158247866,
      '4': 1,
      '5': 5,
      '10': 'coordinatordpusize'
    },
    {
      '1': 'defaultexecutordpusize',
      '3': 42210020,
      '4': 1,
      '5': 5,
      '10': 'defaultexecutordpusize'
    },
    {
      '1': 'maxconcurrentdpus',
      '3': 349976345,
      '4': 1,
      '5': 5,
      '10': 'maxconcurrentdpus'
    },
    {
      '1': 'sparkproperties',
      '3': 157244376,
      '4': 3,
      '5': 11,
      '6': '.athena.EngineConfiguration.SparkpropertiesEntry',
      '10': 'sparkproperties'
    },
  ],
  '3': [
    EngineConfiguration_AdditionalconfigsEntry$json,
    EngineConfiguration_SparkpropertiesEntry$json
  ],
};

@$core.Deprecated('Use engineConfigurationDescriptor instead')
const EngineConfiguration_AdditionalconfigsEntry$json = {
  '1': 'AdditionalconfigsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use engineConfigurationDescriptor instead')
const EngineConfiguration_SparkpropertiesEntry$json = {
  '1': 'SparkpropertiesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `EngineConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List engineConfigurationDescriptor = $convert.base64Decode(
    'ChNFbmdpbmVDb25maWd1cmF0aW9uEmQKEWFkZGl0aW9uYWxjb25maWdzGJyr2qIBIAMoCzIyLm'
    'F0aGVuYS5FbmdpbmVDb25maWd1cmF0aW9uLkFkZGl0aW9uYWxjb25maWdzRW50cnlSEWFkZGl0'
    'aW9uYWxjb25maWdzEkMKD2NsYXNzaWZpY2F0aW9ucxjX/Nk9IAMoCzIWLmF0aGVuYS5DbGFzc2'
    'lmaWNhdGlvblIPY2xhc3NpZmljYXRpb25zEjEKEmNvb3JkaW5hdG9yZHB1c2l6ZRi617pLIAEo'
    'BVISY29vcmRpbmF0b3JkcHVzaXplEjkKFmRlZmF1bHRleGVjdXRvcmRwdXNpemUY5KWQFCABKA'
    'VSFmRlZmF1bHRleGVjdXRvcmRwdXNpemUSMAoRbWF4Y29uY3VycmVudGRwdXMYme7wpgEgASgF'
    'UhFtYXhjb25jdXJyZW50ZHB1cxJdCg9zcGFya3Byb3BlcnRpZXMY2Lf9SiADKAsyMC5hdGhlbm'
    'EuRW5naW5lQ29uZmlndXJhdGlvbi5TcGFya3Byb3BlcnRpZXNFbnRyeVIPc3Bhcmtwcm9wZXJ0'
    'aWVzGkQKFkFkZGl0aW9uYWxjb25maWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdW'
    'UYAiABKAlSBXZhbHVlOgI4ARpCChRTcGFya3Byb3BlcnRpZXNFbnRyeRIQCgNrZXkYASABKAlS'
    'A2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use engineVersionDescriptor instead')
const EngineVersion$json = {
  '1': 'EngineVersion',
  '2': [
    {
      '1': 'effectiveengineversion',
      '3': 382949365,
      '4': 1,
      '5': 9,
      '10': 'effectiveengineversion'
    },
    {
      '1': 'selectedengineversion',
      '3': 199609827,
      '4': 1,
      '5': 9,
      '10': 'selectedengineversion'
    },
  ],
};

/// Descriptor for `EngineVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List engineVersionDescriptor = $convert.base64Decode(
    'Cg1FbmdpbmVWZXJzaW9uEjoKFmVmZmVjdGl2ZWVuZ2luZXZlcnNpb24Y9a/NtgEgASgJUhZlZm'
    'ZlY3RpdmVlbmdpbmV2ZXJzaW9uEjcKFXNlbGVjdGVkZW5naW5ldmVyc2lvbhjjm5dfIAEoCVIV'
    'c2VsZWN0ZWRlbmdpbmV2ZXJzaW9u');

@$core.Deprecated('Use executorsSummaryDescriptor instead')
const ExecutorsSummary$json = {
  '1': 'ExecutorsSummary',
  '2': [
    {'1': 'executorid', '3': 307566450, '4': 1, '5': 9, '10': 'executorid'},
    {'1': 'executorsize', '3': 220298710, '4': 1, '5': 3, '10': 'executorsize'},
    {
      '1': 'executorstate',
      '3': 187122914,
      '4': 1,
      '5': 14,
      '6': '.athena.ExecutorState',
      '10': 'executorstate'
    },
    {
      '1': 'executortype',
      '3': 213715221,
      '4': 1,
      '5': 14,
      '6': '.athena.ExecutorType',
      '10': 'executortype'
    },
    {
      '1': 'startdatetime',
      '3': 88518355,
      '4': 1,
      '5': 3,
      '10': 'startdatetime'
    },
    {
      '1': 'terminationdatetime',
      '3': 253850405,
      '4': 1,
      '5': 3,
      '10': 'terminationdatetime'
    },
  ],
};

/// Descriptor for `ExecutorsSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executorsSummaryDescriptor = $convert.base64Decode(
    'ChBFeGVjdXRvcnNTdW1tYXJ5EiIKCmV4ZWN1dG9yaWQY8q7UkgEgASgJUgpleGVjdXRvcmlkEi'
    'UKDGV4ZWN1dG9yc2l6ZRjW+4VpIAEoA1IMZXhlY3V0b3JzaXplEj4KDWV4ZWN1dG9yc3RhdGUY'
    '4omdWSABKA4yFS5hdGhlbmEuRXhlY3V0b3JTdGF0ZVINZXhlY3V0b3JzdGF0ZRI7CgxleGVjdX'
    'RvcnR5cGUYlZL0ZSABKA4yFC5hdGhlbmEuRXhlY3V0b3JUeXBlUgxleGVjdXRvcnR5cGUSJwoN'
    'c3RhcnRkYXRldGltZRjT3ZoqIAEoA1INc3RhcnRkYXRldGltZRIzChN0ZXJtaW5hdGlvbmRhdG'
    'V0aW1lGKXmhXkgASgDUhN0ZXJtaW5hdGlvbmRhdGV0aW1l');

@$core.Deprecated('Use exportNotebookInputDescriptor instead')
const ExportNotebookInput$json = {
  '1': 'ExportNotebookInput',
  '2': [
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
  ],
};

/// Descriptor for `ExportNotebookInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportNotebookInputDescriptor = $convert.base64Decode(
    'ChNFeHBvcnROb3RlYm9va0lucHV0EiEKCm5vdGVib29raWQY3rSVSyABKAlSCm5vdGVib29raW'
    'Q=');

@$core.Deprecated('Use exportNotebookOutputDescriptor instead')
const ExportNotebookOutput$json = {
  '1': 'ExportNotebookOutput',
  '2': [
    {
      '1': 'notebookmetadata',
      '3': 77536390,
      '4': 1,
      '5': 11,
      '6': '.athena.NotebookMetadata',
      '10': 'notebookmetadata'
    },
    {'1': 'payload', '3': 6526790, '4': 1, '5': 9, '10': 'payload'},
  ],
};

/// Descriptor for `ExportNotebookOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportNotebookOutputDescriptor = $convert.base64Decode(
    'ChRFeHBvcnROb3RlYm9va091dHB1dBJHChBub3RlYm9va21ldGFkYXRhGIa5/CQgASgLMhguYX'
    'RoZW5hLk5vdGVib29rTWV0YWRhdGFSEG5vdGVib29rbWV0YWRhdGESGwoHcGF5bG9hZBjGro4D'
    'IAEoCVIHcGF5bG9hZA==');

@$core.Deprecated('Use filterDefinitionDescriptor instead')
const FilterDefinition$json = {
  '1': 'FilterDefinition',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `FilterDefinition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterDefinitionDescriptor = $convert
    .base64Decode('ChBGaWx0ZXJEZWZpbml0aW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use getCalculationExecutionCodeRequestDescriptor instead')
const GetCalculationExecutionCodeRequest$json = {
  '1': 'GetCalculationExecutionCodeRequest',
  '2': [
    {
      '1': 'calculationexecutionid',
      '3': 80028050,
      '4': 1,
      '5': 9,
      '10': 'calculationexecutionid'
    },
  ],
};

/// Descriptor for `GetCalculationExecutionCodeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalculationExecutionCodeRequestDescriptor =
    $convert.base64Decode(
        'CiJHZXRDYWxjdWxhdGlvbkV4ZWN1dGlvbkNvZGVSZXF1ZXN0EjkKFmNhbGN1bGF0aW9uZXhlY3'
        'V0aW9uaWQYksOUJiABKAlSFmNhbGN1bGF0aW9uZXhlY3V0aW9uaWQ=');

@$core.Deprecated('Use getCalculationExecutionCodeResponseDescriptor instead')
const GetCalculationExecutionCodeResponse$json = {
  '1': 'GetCalculationExecutionCodeResponse',
  '2': [
    {'1': 'codeblock', '3': 23945838, '4': 1, '5': 9, '10': 'codeblock'},
  ],
};

/// Descriptor for `GetCalculationExecutionCodeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalculationExecutionCodeResponseDescriptor =
    $convert.base64Decode(
        'CiNHZXRDYWxjdWxhdGlvbkV4ZWN1dGlvbkNvZGVSZXNwb25zZRIfCgljb2RlYmxvY2sY7sS1Cy'
        'ABKAlSCWNvZGVibG9jaw==');

@$core.Deprecated('Use getCalculationExecutionRequestDescriptor instead')
const GetCalculationExecutionRequest$json = {
  '1': 'GetCalculationExecutionRequest',
  '2': [
    {
      '1': 'calculationexecutionid',
      '3': 80028050,
      '4': 1,
      '5': 9,
      '10': 'calculationexecutionid'
    },
  ],
};

/// Descriptor for `GetCalculationExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalculationExecutionRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXRDYWxjdWxhdGlvbkV4ZWN1dGlvblJlcXVlc3QSOQoWY2FsY3VsYXRpb25leGVjdXRpb2'
        '5pZBiSw5QmIAEoCVIWY2FsY3VsYXRpb25leGVjdXRpb25pZA==');

@$core.Deprecated('Use getCalculationExecutionResponseDescriptor instead')
const GetCalculationExecutionResponse$json = {
  '1': 'GetCalculationExecutionResponse',
  '2': [
    {
      '1': 'calculationexecutionid',
      '3': 80028050,
      '4': 1,
      '5': 9,
      '10': 'calculationexecutionid'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'result',
      '3': 273346629,
      '4': 1,
      '5': 11,
      '6': '.athena.CalculationResult',
      '10': 'result'
    },
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {
      '1': 'statistics',
      '3': 510636075,
      '4': 1,
      '5': 11,
      '6': '.athena.CalculationStatistics',
      '10': 'statistics'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.athena.CalculationStatus',
      '10': 'status'
    },
    {
      '1': 'workingdirectory',
      '3': 478970252,
      '4': 1,
      '5': 9,
      '10': 'workingdirectory'
    },
  ],
};

/// Descriptor for `GetCalculationExecutionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalculationExecutionResponseDescriptor = $convert.base64Decode(
    'Ch9HZXRDYWxjdWxhdGlvbkV4ZWN1dGlvblJlc3BvbnNlEjkKFmNhbGN1bGF0aW9uZXhlY3V0aW'
    '9uaWQYksOUJiABKAlSFmNhbGN1bGF0aW9uZXhlY3V0aW9uaWQSIwoLZGVzY3JpcHRpb24YivT5'
    'NiABKAlSC2Rlc2NyaXB0aW9uEjUKBnJlc3VsdBjF4KuCASABKAsyGS5hdGhlbmEuQ2FsY3VsYX'
    'Rpb25SZXN1bHRSBnJlc3VsdBIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb25pZBJBCgpz'
    'dGF0aXN0aWNzGKvgvvMBIAEoCzIdLmF0aGVuYS5DYWxjdWxhdGlvblN0YXRpc3RpY3NSCnN0YX'
    'Rpc3RpY3MSNAoGc3RhdHVzGJDk+wIgASgLMhkuYXRoZW5hLkNhbGN1bGF0aW9uU3RhdHVzUgZz'
    'dGF0dXMSLgoQd29ya2luZ2RpcmVjdG9yeRiMg7LkASABKAlSEHdvcmtpbmdkaXJlY3Rvcnk=');

@$core.Deprecated('Use getCalculationExecutionStatusRequestDescriptor instead')
const GetCalculationExecutionStatusRequest$json = {
  '1': 'GetCalculationExecutionStatusRequest',
  '2': [
    {
      '1': 'calculationexecutionid',
      '3': 80028050,
      '4': 1,
      '5': 9,
      '10': 'calculationexecutionid'
    },
  ],
};

/// Descriptor for `GetCalculationExecutionStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalculationExecutionStatusRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRDYWxjdWxhdGlvbkV4ZWN1dGlvblN0YXR1c1JlcXVlc3QSOQoWY2FsY3VsYXRpb25leG'
        'VjdXRpb25pZBiSw5QmIAEoCVIWY2FsY3VsYXRpb25leGVjdXRpb25pZA==');

@$core.Deprecated('Use getCalculationExecutionStatusResponseDescriptor instead')
const GetCalculationExecutionStatusResponse$json = {
  '1': 'GetCalculationExecutionStatusResponse',
  '2': [
    {
      '1': 'statistics',
      '3': 510636075,
      '4': 1,
      '5': 11,
      '6': '.athena.CalculationStatistics',
      '10': 'statistics'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.athena.CalculationStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetCalculationExecutionStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalculationExecutionStatusResponseDescriptor =
    $convert.base64Decode(
        'CiVHZXRDYWxjdWxhdGlvbkV4ZWN1dGlvblN0YXR1c1Jlc3BvbnNlEkEKCnN0YXRpc3RpY3MYq+'
        'C+8wEgASgLMh0uYXRoZW5hLkNhbGN1bGF0aW9uU3RhdGlzdGljc1IKc3RhdGlzdGljcxI0CgZz'
        'dGF0dXMYkOT7AiABKAsyGS5hdGhlbmEuQ2FsY3VsYXRpb25TdGF0dXNSBnN0YXR1cw==');

@$core
    .Deprecated('Use getCapacityAssignmentConfigurationInputDescriptor instead')
const GetCapacityAssignmentConfigurationInput$json = {
  '1': 'GetCapacityAssignmentConfigurationInput',
  '2': [
    {
      '1': 'capacityreservationname',
      '3': 327567687,
      '4': 1,
      '5': 9,
      '10': 'capacityreservationname'
    },
  ],
};

/// Descriptor for `GetCapacityAssignmentConfigurationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCapacityAssignmentConfigurationInputDescriptor =
    $convert.base64Decode(
        'CidHZXRDYXBhY2l0eUFzc2lnbm1lbnRDb25maWd1cmF0aW9uSW5wdXQSPAoXY2FwYWNpdHlyZX'
        'NlcnZhdGlvbm5hbWUYx5KZnAEgASgJUhdjYXBhY2l0eXJlc2VydmF0aW9ubmFtZQ==');

@$core.Deprecated(
    'Use getCapacityAssignmentConfigurationOutputDescriptor instead')
const GetCapacityAssignmentConfigurationOutput$json = {
  '1': 'GetCapacityAssignmentConfigurationOutput',
  '2': [
    {
      '1': 'capacityassignmentconfiguration',
      '3': 9812735,
      '4': 1,
      '5': 11,
      '6': '.athena.CapacityAssignmentConfiguration',
      '10': 'capacityassignmentconfiguration'
    },
  ],
};

/// Descriptor for `GetCapacityAssignmentConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCapacityAssignmentConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CihHZXRDYXBhY2l0eUFzc2lnbm1lbnRDb25maWd1cmF0aW9uT3V0cHV0EnQKH2NhcGFjaXR5YX'
        'NzaWdubWVudGNvbmZpZ3VyYXRpb24Y//XWBCABKAsyJy5hdGhlbmEuQ2FwYWNpdHlBc3NpZ25t'
        'ZW50Q29uZmlndXJhdGlvblIfY2FwYWNpdHlhc3NpZ25tZW50Y29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use getCapacityReservationInputDescriptor instead')
const GetCapacityReservationInput$json = {
  '1': 'GetCapacityReservationInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetCapacityReservationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCapacityReservationInputDescriptor =
    $convert.base64Decode(
        'ChtHZXRDYXBhY2l0eVJlc2VydmF0aW9uSW5wdXQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use getCapacityReservationOutputDescriptor instead')
const GetCapacityReservationOutput$json = {
  '1': 'GetCapacityReservationOutput',
  '2': [
    {
      '1': 'capacityreservation',
      '3': 451869958,
      '4': 1,
      '5': 11,
      '6': '.athena.CapacityReservation',
      '10': 'capacityreservation'
    },
  ],
};

/// Descriptor for `GetCapacityReservationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCapacityReservationOutputDescriptor =
    $convert.base64Decode(
        'ChxHZXRDYXBhY2l0eVJlc2VydmF0aW9uT3V0cHV0ElEKE2NhcGFjaXR5cmVzZXJ2YXRpb24Yhv'
        'q71wEgASgLMhsuYXRoZW5hLkNhcGFjaXR5UmVzZXJ2YXRpb25SE2NhcGFjaXR5cmVzZXJ2YXRp'
        'b24=');

@$core.Deprecated('Use getDataCatalogInputDescriptor instead')
const GetDataCatalogInput$json = {
  '1': 'GetDataCatalogInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `GetDataCatalogInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDataCatalogInputDescriptor = $convert.base64Decode(
    'ChNHZXREYXRhQ2F0YWxvZ0lucHV0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSIAoJd29ya2dyb3'
    'VwGIStofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use getDataCatalogOutputDescriptor instead')
const GetDataCatalogOutput$json = {
  '1': 'GetDataCatalogOutput',
  '2': [
    {
      '1': 'datacatalog',
      '3': 209694649,
      '4': 1,
      '5': 11,
      '6': '.athena.DataCatalog',
      '10': 'datacatalog'
    },
  ],
};

/// Descriptor for `GetDataCatalogOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDataCatalogOutputDescriptor = $convert.base64Decode(
    'ChRHZXREYXRhQ2F0YWxvZ091dHB1dBI4CgtkYXRhY2F0YWxvZxi53/5jIAEoCzITLmF0aGVuYS'
    '5EYXRhQ2F0YWxvZ1ILZGF0YWNhdGFsb2c=');

@$core.Deprecated('Use getDatabaseInputDescriptor instead')
const GetDatabaseInput$json = {
  '1': 'GetDatabaseInput',
  '2': [
    {'1': 'catalogname', '3': 518825212, '4': 1, '5': 9, '10': 'catalogname'},
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `GetDatabaseInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDatabaseInputDescriptor = $convert.base64Decode(
    'ChBHZXREYXRhYmFzZUlucHV0EiQKC2NhdGFsb2duYW1lGPzJsvcBIAEoCVILY2F0YWxvZ25hbW'
    'USJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYmFzZW5hbWUSIAoJd29ya2dyb3VwGISt'
    'ofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use getDatabaseOutputDescriptor instead')
const GetDatabaseOutput$json = {
  '1': 'GetDatabaseOutput',
  '2': [
    {
      '1': 'database',
      '3': 278147289,
      '4': 1,
      '5': 11,
      '6': '.athena.Database',
      '10': 'database'
    },
  ],
};

/// Descriptor for `GetDatabaseOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDatabaseOutputDescriptor = $convert.base64Decode(
    'ChFHZXREYXRhYmFzZU91dHB1dBIwCghkYXRhYmFzZRjZ4dCEASABKAsyEC5hdGhlbmEuRGF0YW'
    'Jhc2VSCGRhdGFiYXNl');

@$core.Deprecated('Use getNamedQueryInputDescriptor instead')
const GetNamedQueryInput$json = {
  '1': 'GetNamedQueryInput',
  '2': [
    {'1': 'namedqueryid', '3': 330896872, '4': 1, '5': 9, '10': 'namedqueryid'},
  ],
};

/// Descriptor for `GetNamedQueryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNamedQueryInputDescriptor = $convert.base64Decode(
    'ChJHZXROYW1lZFF1ZXJ5SW5wdXQSJgoMbmFtZWRxdWVyeWlkGOir5J0BIAEoCVIMbmFtZWRxdW'
    'VyeWlk');

@$core.Deprecated('Use getNamedQueryOutputDescriptor instead')
const GetNamedQueryOutput$json = {
  '1': 'GetNamedQueryOutput',
  '2': [
    {
      '1': 'namedquery',
      '3': 170671691,
      '4': 1,
      '5': 11,
      '6': '.athena.NamedQuery',
      '10': 'namedquery'
    },
  ],
};

/// Descriptor for `GetNamedQueryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNamedQueryOutputDescriptor = $convert.base64Decode(
    'ChNHZXROYW1lZFF1ZXJ5T3V0cHV0EjUKCm5hbWVkcXVlcnkYy/ywUSABKAsyEi5hdGhlbmEuTm'
    'FtZWRRdWVyeVIKbmFtZWRxdWVyeQ==');

@$core.Deprecated('Use getNotebookMetadataInputDescriptor instead')
const GetNotebookMetadataInput$json = {
  '1': 'GetNotebookMetadataInput',
  '2': [
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
  ],
};

/// Descriptor for `GetNotebookMetadataInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNotebookMetadataInputDescriptor =
    $convert.base64Decode(
        'ChhHZXROb3RlYm9va01ldGFkYXRhSW5wdXQSIQoKbm90ZWJvb2tpZBjetJVLIAEoCVIKbm90ZW'
        'Jvb2tpZA==');

@$core.Deprecated('Use getNotebookMetadataOutputDescriptor instead')
const GetNotebookMetadataOutput$json = {
  '1': 'GetNotebookMetadataOutput',
  '2': [
    {
      '1': 'notebookmetadata',
      '3': 77536390,
      '4': 1,
      '5': 11,
      '6': '.athena.NotebookMetadata',
      '10': 'notebookmetadata'
    },
  ],
};

/// Descriptor for `GetNotebookMetadataOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNotebookMetadataOutputDescriptor =
    $convert.base64Decode(
        'ChlHZXROb3RlYm9va01ldGFkYXRhT3V0cHV0EkcKEG5vdGVib29rbWV0YWRhdGEYhrn8JCABKA'
        'syGC5hdGhlbmEuTm90ZWJvb2tNZXRhZGF0YVIQbm90ZWJvb2ttZXRhZGF0YQ==');

@$core.Deprecated('Use getPreparedStatementInputDescriptor instead')
const GetPreparedStatementInput$json = {
  '1': 'GetPreparedStatementInput',
  '2': [
    {
      '1': 'statementname',
      '3': 23047926,
      '4': 1,
      '5': 9,
      '10': 'statementname'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `GetPreparedStatementInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPreparedStatementInputDescriptor =
    $convert.base64Decode(
        'ChlHZXRQcmVwYXJlZFN0YXRlbWVudElucHV0EicKDXN0YXRlbWVudG5hbWUY9t3+CiABKAlSDX'
        'N0YXRlbWVudG5hbWUSIAoJd29ya2dyb3VwGIStofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use getPreparedStatementOutputDescriptor instead')
const GetPreparedStatementOutput$json = {
  '1': 'GetPreparedStatementOutput',
  '2': [
    {
      '1': 'preparedstatement',
      '3': 164916502,
      '4': 1,
      '5': 11,
      '6': '.athena.PreparedStatement',
      '10': 'preparedstatement'
    },
  ],
};

/// Descriptor for `GetPreparedStatementOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPreparedStatementOutputDescriptor =
    $convert.base64Decode(
        'ChpHZXRQcmVwYXJlZFN0YXRlbWVudE91dHB1dBJKChFwcmVwYXJlZHN0YXRlbWVudBiW2tFOIA'
        'EoCzIZLmF0aGVuYS5QcmVwYXJlZFN0YXRlbWVudFIRcHJlcGFyZWRzdGF0ZW1lbnQ=');

@$core.Deprecated('Use getQueryExecutionInputDescriptor instead')
const GetQueryExecutionInput$json = {
  '1': 'GetQueryExecutionInput',
  '2': [
    {
      '1': 'queryexecutionid',
      '3': 467615503,
      '4': 1,
      '5': 9,
      '10': 'queryexecutionid'
    },
  ],
};

/// Descriptor for `GetQueryExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryExecutionInputDescriptor =
    $convert.base64Decode(
        'ChZHZXRRdWVyeUV4ZWN1dGlvbklucHV0Ei4KEHF1ZXJ5ZXhlY3V0aW9uaWQYj/783gEgASgJUh'
        'BxdWVyeWV4ZWN1dGlvbmlk');

@$core.Deprecated('Use getQueryExecutionOutputDescriptor instead')
const GetQueryExecutionOutput$json = {
  '1': 'GetQueryExecutionOutput',
  '2': [
    {
      '1': 'queryexecution',
      '3': 62173540,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryExecution',
      '10': 'queryexecution'
    },
  ],
};

/// Descriptor for `GetQueryExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryExecutionOutputDescriptor =
    $convert.base64Decode(
        'ChdHZXRRdWVyeUV4ZWN1dGlvbk91dHB1dBJBCg5xdWVyeWV4ZWN1dGlvbhjk4tIdIAEoCzIWLm'
        'F0aGVuYS5RdWVyeUV4ZWN1dGlvblIOcXVlcnlleGVjdXRpb24=');

@$core.Deprecated('Use getQueryResultsInputDescriptor instead')
const GetQueryResultsInput$json = {
  '1': 'GetQueryResultsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'queryexecutionid',
      '3': 467615503,
      '4': 1,
      '5': 9,
      '10': 'queryexecutionid'
    },
    {
      '1': 'queryresulttype',
      '3': 3724043,
      '4': 1,
      '5': 14,
      '6': '.athena.QueryResultType',
      '10': 'queryresulttype'
    },
  ],
};

/// Descriptor for `GetQueryResultsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryResultsInputDescriptor = $convert.base64Decode(
    'ChRHZXRRdWVyeVJlc3VsdHNJbnB1dBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdW'
    'x0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIuChBxdWVyeWV4ZWN1dGlvbmlk'
    'GI/+/N4BIAEoCVIQcXVlcnlleGVjdXRpb25pZBJECg9xdWVyeXJlc3VsdHR5cGUYi6bjASABKA'
    '4yFy5hdGhlbmEuUXVlcnlSZXN1bHRUeXBlUg9xdWVyeXJlc3VsdHR5cGU=');

@$core.Deprecated('Use getQueryResultsOutputDescriptor instead')
const GetQueryResultsOutput$json = {
  '1': 'GetQueryResultsOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resultset',
      '3': 146844365,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultSet',
      '10': 'resultset'
    },
    {'1': 'updatecount', '3': 312282732, '4': 1, '5': 3, '10': 'updatecount'},
  ],
};

/// Descriptor for `GetQueryResultsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryResultsOutputDescriptor = $convert.base64Decode(
    'ChVHZXRRdWVyeVJlc3VsdHNPdXRwdXQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW'
    '4SMgoJcmVzdWx0c2V0GM3VgkYgASgLMhEuYXRoZW5hLlJlc3VsdFNldFIJcmVzdWx0c2V0EiQK'
    'C3VwZGF0ZWNvdW50GOyc9JQBIAEoA1ILdXBkYXRlY291bnQ=');

@$core.Deprecated('Use getQueryRuntimeStatisticsInputDescriptor instead')
const GetQueryRuntimeStatisticsInput$json = {
  '1': 'GetQueryRuntimeStatisticsInput',
  '2': [
    {
      '1': 'queryexecutionid',
      '3': 467615503,
      '4': 1,
      '5': 9,
      '10': 'queryexecutionid'
    },
  ],
};

/// Descriptor for `GetQueryRuntimeStatisticsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryRuntimeStatisticsInputDescriptor =
    $convert.base64Decode(
        'Ch5HZXRRdWVyeVJ1bnRpbWVTdGF0aXN0aWNzSW5wdXQSLgoQcXVlcnlleGVjdXRpb25pZBiP/v'
        'zeASABKAlSEHF1ZXJ5ZXhlY3V0aW9uaWQ=');

@$core.Deprecated('Use getQueryRuntimeStatisticsOutputDescriptor instead')
const GetQueryRuntimeStatisticsOutput$json = {
  '1': 'GetQueryRuntimeStatisticsOutput',
  '2': [
    {
      '1': 'queryruntimestatistics',
      '3': 157051307,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryRuntimeStatistics',
      '10': 'queryruntimestatistics'
    },
  ],
};

/// Descriptor for `GetQueryRuntimeStatisticsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueryRuntimeStatisticsOutputDescriptor =
    $convert.base64Decode(
        'Ch9HZXRRdWVyeVJ1bnRpbWVTdGF0aXN0aWNzT3V0cHV0ElkKFnF1ZXJ5cnVudGltZXN0YXRpc3'
        'RpY3MYq9PxSiABKAsyHi5hdGhlbmEuUXVlcnlSdW50aW1lU3RhdGlzdGljc1IWcXVlcnlydW50'
        'aW1lc3RhdGlzdGljcw==');

@$core.Deprecated('Use getResourceDashboardRequestDescriptor instead')
const GetResourceDashboardRequest$json = {
  '1': 'GetResourceDashboardRequest',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `GetResourceDashboardRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourceDashboardRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRSZXNvdXJjZURhc2hib2FyZFJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4Y7cCZsAEgASgJUg'
        'tyZXNvdXJjZWFybg==');

@$core.Deprecated('Use getResourceDashboardResponseDescriptor instead')
const GetResourceDashboardResponse$json = {
  '1': 'GetResourceDashboardResponse',
  '2': [
    {'1': 'url', '3': 354018239, '4': 1, '5': 9, '10': 'url'},
  ],
};

/// Descriptor for `GetResourceDashboardResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourceDashboardResponseDescriptor =
    $convert.base64Decode(
        'ChxHZXRSZXNvdXJjZURhc2hib2FyZFJlc3BvbnNlEhQKA3VybBi/x+eoASABKAlSA3VybA==');

@$core.Deprecated('Use getSessionEndpointRequestDescriptor instead')
const GetSessionEndpointRequest$json = {
  '1': 'GetSessionEndpointRequest',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `GetSessionEndpointRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionEndpointRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRTZXNzaW9uRW5kcG9pbnRSZXF1ZXN0Eh8KCXNlc3Npb25pZBi7hOUJIAEoCVIJc2Vzc2'
        'lvbmlk');

@$core.Deprecated('Use getSessionEndpointResponseDescriptor instead')
const GetSessionEndpointResponse$json = {
  '1': 'GetSessionEndpointResponse',
  '2': [
    {'1': 'authtoken', '3': 349771493, '4': 1, '5': 9, '10': 'authtoken'},
    {
      '1': 'authtokenexpirationtime',
      '3': 406065689,
      '4': 1,
      '5': 9,
      '10': 'authtokenexpirationtime'
    },
    {'1': 'endpointurl', '3': 31787414, '4': 1, '5': 9, '10': 'endpointurl'},
  ],
};

/// Descriptor for `GetSessionEndpointResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionEndpointResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRTZXNzaW9uRW5kcG9pbnRSZXNwb25zZRIgCglhdXRodG9rZW4Y5a3kpgEgASgJUglhdX'
        'RodG9rZW4SPAoXYXV0aHRva2VuZXhwaXJhdGlvbnRpbWUYmaTQwQEgASgJUhdhdXRodG9rZW5l'
        'eHBpcmF0aW9udGltZRIjCgtlbmRwb2ludHVybBiWk5QPIAEoCVILZW5kcG9pbnR1cmw=');

@$core.Deprecated('Use getSessionRequestDescriptor instead')
const GetSessionRequest$json = {
  '1': 'GetSessionRequest',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `GetSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionRequestDescriptor = $convert.base64Decode(
    'ChFHZXRTZXNzaW9uUmVxdWVzdBIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb25pZA==');

@$core.Deprecated('Use getSessionResponseDescriptor instead')
const GetSessionResponse$json = {
  '1': 'GetSessionResponse',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'engineconfiguration',
      '3': 341629412,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineConfiguration',
      '10': 'engineconfiguration'
    },
    {
      '1': 'engineversion',
      '3': 44953462,
      '4': 1,
      '5': 9,
      '10': 'engineversion'
    },
    {
      '1': 'monitoringconfiguration',
      '3': 364891928,
      '4': 1,
      '5': 11,
      '6': '.athena.MonitoringConfiguration',
      '10': 'monitoringconfiguration'
    },
    {
      '1': 'notebookversion',
      '3': 528689837,
      '4': 1,
      '5': 9,
      '10': 'notebookversion'
    },
    {
      '1': 'sessionconfiguration',
      '3': 211592776,
      '4': 1,
      '5': 11,
      '6': '.athena.SessionConfiguration',
      '10': 'sessionconfiguration'
    },
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {
      '1': 'statistics',
      '3': 510636075,
      '4': 1,
      '5': 11,
      '6': '.athena.SessionStatistics',
      '10': 'statistics'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.athena.SessionStatus',
      '10': 'status'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `GetSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionResponseDescriptor = $convert.base64Decode(
    'ChJHZXRTZXNzaW9uUmVzcG9uc2USIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW'
    '9uElEKE2VuZ2luZWNvbmZpZ3VyYXRpb24Y5LPzogEgASgLMhsuYXRoZW5hLkVuZ2luZUNvbmZp'
    'Z3VyYXRpb25SE2VuZ2luZWNvbmZpZ3VyYXRpb24SJwoNZW5naW5ldmVyc2lvbhj23rcVIAEoCV'
    'INZW5naW5ldmVyc2lvbhJdChdtb25pdG9yaW5nY29uZmlndXJhdGlvbhiYnv+tASABKAsyHy5h'
    'dGhlbmEuTW9uaXRvcmluZ0NvbmZpZ3VyYXRpb25SF21vbml0b3Jpbmdjb25maWd1cmF0aW9uEi'
    'wKD25vdGVib29rdmVyc2lvbhit1Yz8ASABKAlSD25vdGVib29rdmVyc2lvbhJTChRzZXNzaW9u'
    'Y29uZmlndXJhdGlvbhjIzPJkIAEoCzIcLmF0aGVuYS5TZXNzaW9uQ29uZmlndXJhdGlvblIUc2'
    'Vzc2lvbmNvbmZpZ3VyYXRpb24SHwoJc2Vzc2lvbmlkGLuE5QkgASgJUglzZXNzaW9uaWQSPQoK'
    'c3RhdGlzdGljcxir4L7zASABKAsyGS5hdGhlbmEuU2Vzc2lvblN0YXRpc3RpY3NSCnN0YXRpc3'
    'RpY3MSMAoGc3RhdHVzGJDk+wIgASgLMhUuYXRoZW5hLlNlc3Npb25TdGF0dXNSBnN0YXR1cxIg'
    'Cgl3b3JrZ3JvdXAYhK2h8QEgASgJUgl3b3JrZ3JvdXA=');

@$core.Deprecated('Use getSessionStatusRequestDescriptor instead')
const GetSessionStatusRequest$json = {
  '1': 'GetSessionStatusRequest',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `GetSessionStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionStatusRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRTZXNzaW9uU3RhdHVzUmVxdWVzdBIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb2'
        '5pZA==');

@$core.Deprecated('Use getSessionStatusResponseDescriptor instead')
const GetSessionStatusResponse$json = {
  '1': 'GetSessionStatusResponse',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.athena.SessionStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetSessionStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionStatusResponseDescriptor = $convert.base64Decode(
    'ChhHZXRTZXNzaW9uU3RhdHVzUmVzcG9uc2USHwoJc2Vzc2lvbmlkGLuE5QkgASgJUglzZXNzaW'
    '9uaWQSMAoGc3RhdHVzGJDk+wIgASgLMhUuYXRoZW5hLlNlc3Npb25TdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use getTableMetadataInputDescriptor instead')
const GetTableMetadataInput$json = {
  '1': 'GetTableMetadataInput',
  '2': [
    {'1': 'catalogname', '3': 518825212, '4': 1, '5': 9, '10': 'catalogname'},
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'tablename', '3': 272020061, '4': 1, '5': 9, '10': 'tablename'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `GetTableMetadataInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTableMetadataInputDescriptor = $convert.base64Decode(
    'ChVHZXRUYWJsZU1ldGFkYXRhSW5wdXQSJAoLY2F0YWxvZ25hbWUY/Mmy9wEgASgJUgtjYXRhbG'
    '9nbmFtZRIlCgxkYXRhYmFzZW5hbWUY3LLZKiABKAlSDGRhdGFiYXNlbmFtZRIgCgl0YWJsZW5h'
    'bWUY3eTagQEgASgJUgl0YWJsZW5hbWUSIAoJd29ya2dyb3VwGIStofEBIAEoCVIJd29ya2dyb3'
    'Vw');

@$core.Deprecated('Use getTableMetadataOutputDescriptor instead')
const GetTableMetadataOutput$json = {
  '1': 'GetTableMetadataOutput',
  '2': [
    {
      '1': 'tablemetadata',
      '3': 393071395,
      '4': 1,
      '5': 11,
      '6': '.athena.TableMetadata',
      '10': 'tablemetadata'
    },
  ],
};

/// Descriptor for `GetTableMetadataOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTableMetadataOutputDescriptor =
    $convert.base64Decode(
        'ChZHZXRUYWJsZU1ldGFkYXRhT3V0cHV0Ej8KDXRhYmxlbWV0YWRhdGEYo5a3uwEgASgLMhUuYX'
        'RoZW5hLlRhYmxlTWV0YWRhdGFSDXRhYmxlbWV0YWRhdGE=');

@$core.Deprecated('Use getWorkGroupInputDescriptor instead')
const GetWorkGroupInput$json = {
  '1': 'GetWorkGroupInput',
  '2': [
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `GetWorkGroupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWorkGroupInputDescriptor = $convert.base64Decode(
    'ChFHZXRXb3JrR3JvdXBJbnB1dBIgCgl3b3JrZ3JvdXAYhK2h8QEgASgJUgl3b3JrZ3JvdXA=');

@$core.Deprecated('Use getWorkGroupOutputDescriptor instead')
const GetWorkGroupOutput$json = {
  '1': 'GetWorkGroupOutput',
  '2': [
    {
      '1': 'workgroup',
      '3': 505960068,
      '4': 1,
      '5': 11,
      '6': '.athena.WorkGroup',
      '10': 'workgroup'
    },
  ],
};

/// Descriptor for `GetWorkGroupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWorkGroupOutputDescriptor = $convert.base64Decode(
    'ChJHZXRXb3JrR3JvdXBPdXRwdXQSMwoJd29ya2dyb3VwGIStofEBIAEoCzIRLmF0aGVuYS5Xb3'
    'JrR3JvdXBSCXdvcmtncm91cA==');

@$core.Deprecated('Use identityCenterConfigurationDescriptor instead')
const IdentityCenterConfiguration$json = {
  '1': 'IdentityCenterConfiguration',
  '2': [
    {
      '1': 'enableidentitycenter',
      '3': 418515238,
      '4': 1,
      '5': 8,
      '10': 'enableidentitycenter'
    },
    {
      '1': 'identitycenterinstancearn',
      '3': 469575873,
      '4': 1,
      '5': 9,
      '10': 'identitycenterinstancearn'
    },
  ],
};

/// Descriptor for `IdentityCenterConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List identityCenterConfigurationDescriptor =
    $convert.base64Decode(
        'ChtJZGVudGl0eUNlbnRlckNvbmZpZ3VyYXRpb24SNgoUZW5hYmxlaWRlbnRpdHljZW50ZXIYpp'
        'LIxwEgASgIUhRlbmFibGVpZGVudGl0eWNlbnRlchJAChlpZGVudGl0eWNlbnRlcmluc3RhbmNl'
        'YXJuGMHR9N8BIAEoCVIZaWRlbnRpdHljZW50ZXJpbnN0YW5jZWFybg==');

@$core.Deprecated('Use importNotebookInputDescriptor instead')
const ImportNotebookInput$json = {
  '1': 'ImportNotebookInput',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'notebooks3locationuri',
      '3': 460333982,
      '4': 1,
      '5': 9,
      '10': 'notebooks3locationuri'
    },
    {'1': 'payload', '3': 6526790, '4': 1, '5': 9, '10': 'payload'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.athena.NotebookType',
      '10': 'type'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ImportNotebookInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importNotebookInputDescriptor = $convert.base64Decode(
    'ChNJbXBvcnROb3RlYm9va0lucHV0EjIKEmNsaWVudHJlcXVlc3R0b2tlbhjx76LZASABKAlSEm'
    'NsaWVudHJlcXVlc3R0b2tlbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEjgKFW5vdGVib29rczNs'
    'b2NhdGlvbnVyaRiex8DbASABKAlSFW5vdGVib29rczNsb2NhdGlvbnVyaRIbCgdwYXlsb2FkGM'
    'aujgMgASgJUgdwYXlsb2FkEiwKBHR5cGUY7qDXigEgASgOMhQuYXRoZW5hLk5vdGVib29rVHlw'
    'ZVIEdHlwZRIgCgl3b3JrZ3JvdXAYhK2h8QEgASgJUgl3b3JrZ3JvdXA=');

@$core.Deprecated('Use importNotebookOutputDescriptor instead')
const ImportNotebookOutput$json = {
  '1': 'ImportNotebookOutput',
  '2': [
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
  ],
};

/// Descriptor for `ImportNotebookOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importNotebookOutputDescriptor = $convert.base64Decode(
    'ChRJbXBvcnROb3RlYm9va091dHB1dBIhCgpub3RlYm9va2lkGN60lUsgASgJUgpub3RlYm9va2'
    'lk');

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

@$core.Deprecated('Use invalidRequestExceptionDescriptor instead')
const InvalidRequestException$json = {
  '1': 'InvalidRequestException',
  '2': [
    {
      '1': 'athenaerrorcode',
      '3': 335153222,
      '4': 1,
      '5': 9,
      '10': 'athenaerrorcode'
    },
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidRequestException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRequestExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbnZhbGlkUmVxdWVzdEV4Y2VwdGlvbhIsCg9hdGhlbmFlcnJvcmNvZGUYxpDonwEgASgJUg'
        '9hdGhlbmFlcnJvcmNvZGUSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use listApplicationDPUSizesInputDescriptor instead')
const ListApplicationDPUSizesInput$json = {
  '1': 'ListApplicationDPUSizesInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListApplicationDPUSizesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listApplicationDPUSizesInputDescriptor =
    $convert.base64Decode(
        'ChxMaXN0QXBwbGljYXRpb25EUFVTaXplc0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUg'
        'ptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listApplicationDPUSizesOutputDescriptor instead')
const ListApplicationDPUSizesOutput$json = {
  '1': 'ListApplicationDPUSizesOutput',
  '2': [
    {
      '1': 'applicationdpusizes',
      '3': 60251851,
      '4': 3,
      '5': 11,
      '6': '.athena.ApplicationDPUSizes',
      '10': 'applicationdpusizes'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListApplicationDPUSizesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listApplicationDPUSizesOutputDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0QXBwbGljYXRpb25EUFVTaXplc091dHB1dBJQChNhcHBsaWNhdGlvbmRwdXNpemVzGM'
        'u93RwgAygLMhsuYXRoZW5hLkFwcGxpY2F0aW9uRFBVU2l6ZXNSE2FwcGxpY2F0aW9uZHB1c2l6'
        'ZXMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listCalculationExecutionsRequestDescriptor instead')
const ListCalculationExecutionsRequest$json = {
  '1': 'ListCalculationExecutionsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {
      '1': 'statefilter',
      '3': 184693297,
      '4': 1,
      '5': 14,
      '6': '.athena.CalculationExecutionState',
      '10': 'statefilter'
    },
  ],
};

/// Descriptor for `ListCalculationExecutionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCalculationExecutionsRequestDescriptor =
    $convert.base64Decode(
        'CiBMaXN0Q2FsY3VsYXRpb25FeGVjdXRpb25zUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIA'
        'EoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIfCglzZXNz'
        'aW9uaWQYu4TlCSABKAlSCXNlc3Npb25pZBJGCgtzdGF0ZWZpbHRlchix5IhYIAEoDjIhLmF0aG'
        'VuYS5DYWxjdWxhdGlvbkV4ZWN1dGlvblN0YXRlUgtzdGF0ZWZpbHRlcg==');

@$core.Deprecated('Use listCalculationExecutionsResponseDescriptor instead')
const ListCalculationExecutionsResponse$json = {
  '1': 'ListCalculationExecutionsResponse',
  '2': [
    {
      '1': 'calculations',
      '3': 335028832,
      '4': 3,
      '5': 11,
      '6': '.athena.CalculationSummary',
      '10': 'calculations'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCalculationExecutionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCalculationExecutionsResponseDescriptor =
    $convert.base64Decode(
        'CiFMaXN0Q2FsY3VsYXRpb25FeGVjdXRpb25zUmVzcG9uc2USQgoMY2FsY3VsYXRpb25zGODE4J'
        '8BIAMoCzIaLmF0aGVuYS5DYWxjdWxhdGlvblN1bW1hcnlSDGNhbGN1bGF0aW9ucxIfCgluZXh0'
        'dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listCapacityReservationsInputDescriptor instead')
const ListCapacityReservationsInput$json = {
  '1': 'ListCapacityReservationsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCapacityReservationsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCapacityReservationsInputDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q2FwYWNpdHlSZXNlcnZhdGlvbnNJbnB1dBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBV'
        'IKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listCapacityReservationsOutputDescriptor instead')
const ListCapacityReservationsOutput$json = {
  '1': 'ListCapacityReservationsOutput',
  '2': [
    {
      '1': 'capacityreservations',
      '3': 473497795,
      '4': 3,
      '5': 11,
      '6': '.athena.CapacityReservation',
      '10': 'capacityreservations'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCapacityReservationsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCapacityReservationsOutputDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0Q2FwYWNpdHlSZXNlcnZhdGlvbnNPdXRwdXQSUwoUY2FwYWNpdHlyZXNlcnZhdGlvbn'
        'MYw4Hk4QEgAygLMhsuYXRoZW5hLkNhcGFjaXR5UmVzZXJ2YXRpb25SFGNhcGFjaXR5cmVzZXJ2'
        'YXRpb25zEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listDataCatalogsInputDescriptor instead')
const ListDataCatalogsInput$json = {
  '1': 'ListDataCatalogsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ListDataCatalogsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDataCatalogsInputDescriptor = $convert.base64Decode(
    'ChVMaXN0RGF0YUNhdGFsb2dzSW5wdXQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3'
    'VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIAoJd29ya2dyb3VwGIStofEB'
    'IAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use listDataCatalogsOutputDescriptor instead')
const ListDataCatalogsOutput$json = {
  '1': 'ListDataCatalogsOutput',
  '2': [
    {
      '1': 'datacatalogssummary',
      '3': 499180880,
      '4': 3,
      '5': 11,
      '6': '.athena.DataCatalogSummary',
      '10': 'datacatalogssummary'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDataCatalogsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDataCatalogsOutputDescriptor = $convert.base64Decode(
    'ChZMaXN0RGF0YUNhdGFsb2dzT3V0cHV0ElAKE2RhdGFjYXRhbG9nc3N1bW1hcnkY0MqD7gEgAy'
    'gLMhouYXRoZW5hLkRhdGFDYXRhbG9nU3VtbWFyeVITZGF0YWNhdGFsb2dzc3VtbWFyeRIfCglu'
    'ZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listDatabasesInputDescriptor instead')
const ListDatabasesInput$json = {
  '1': 'ListDatabasesInput',
  '2': [
    {'1': 'catalogname', '3': 518825212, '4': 1, '5': 9, '10': 'catalogname'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ListDatabasesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDatabasesInputDescriptor = $convert.base64Decode(
    'ChJMaXN0RGF0YWJhc2VzSW5wdXQSJAoLY2F0YWxvZ25hbWUY/Mmy9wEgASgJUgtjYXRhbG9nbm'
    'FtZRIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6'
    'ZyABKAlSCW5leHR0b2tlbhIgCgl3b3JrZ3JvdXAYhK2h8QEgASgJUgl3b3JrZ3JvdXA=');

@$core.Deprecated('Use listDatabasesOutputDescriptor instead')
const ListDatabasesOutput$json = {
  '1': 'ListDatabasesOutput',
  '2': [
    {
      '1': 'databaselist',
      '3': 393720897,
      '4': 3,
      '5': 11,
      '6': '.athena.Database',
      '10': 'databaselist'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDatabasesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDatabasesOutputDescriptor = $convert.base64Decode(
    'ChNMaXN0RGF0YWJhc2VzT3V0cHV0EjgKDGRhdGFiYXNlbGlzdBjB6N67ASADKAsyEC5hdGhlbm'
    'EuRGF0YWJhc2VSDGRhdGFiYXNlbGlzdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tl'
    'bg==');

@$core.Deprecated('Use listEngineVersionsInputDescriptor instead')
const ListEngineVersionsInput$json = {
  '1': 'ListEngineVersionsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEngineVersionsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEngineVersionsInputDescriptor =
    $convert.base64Decode(
        'ChdMaXN0RW5naW5lVmVyc2lvbnNJbnB1dBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cm'
        'VzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listEngineVersionsOutputDescriptor instead')
const ListEngineVersionsOutput$json = {
  '1': 'ListEngineVersionsOutput',
  '2': [
    {
      '1': 'engineversions',
      '3': 500123251,
      '4': 3,
      '5': 11,
      '6': '.athena.EngineVersion',
      '10': 'engineversions'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEngineVersionsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEngineVersionsOutputDescriptor = $convert.base64Decode(
    'ChhMaXN0RW5naW5lVmVyc2lvbnNPdXRwdXQSQQoOZW5naW5ldmVyc2lvbnMY84y97gEgAygLMh'
    'UuYXRoZW5hLkVuZ2luZVZlcnNpb25SDmVuZ2luZXZlcnNpb25zEh8KCW5leHR0b2tlbhj+hLpn'
    'IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listExecutorsRequestDescriptor instead')
const ListExecutorsRequest$json = {
  '1': 'ListExecutorsRequest',
  '2': [
    {
      '1': 'executorstatefilter',
      '3': 208448852,
      '4': 1,
      '5': 14,
      '6': '.athena.ExecutorState',
      '10': 'executorstatefilter'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `ListExecutorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listExecutorsRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0RXhlY3V0b3JzUmVxdWVzdBJKChNleGVjdXRvcnN0YXRlZmlsdGVyGNTasmMgASgOMh'
    'UuYXRoZW5hLkV4ZWN1dG9yU3RhdGVSE2V4ZWN1dG9yc3RhdGVmaWx0ZXISIgoKbWF4cmVzdWx0'
    'cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW'
    '4SHwoJc2Vzc2lvbmlkGLuE5QkgASgJUglzZXNzaW9uaWQ=');

@$core.Deprecated('Use listExecutorsResponseDescriptor instead')
const ListExecutorsResponse$json = {
  '1': 'ListExecutorsResponse',
  '2': [
    {
      '1': 'executorssummary',
      '3': 450697516,
      '4': 3,
      '5': 11,
      '6': '.athena.ExecutorsSummary',
      '10': 'executorssummary'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `ListExecutorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listExecutorsResponseDescriptor = $convert.base64Decode(
    'ChVMaXN0RXhlY3V0b3JzUmVzcG9uc2USSAoQZXhlY3V0b3Jzc3VtbWFyeRissvTWASADKAsyGC'
    '5hdGhlbmEuRXhlY3V0b3JzU3VtbWFyeVIQZXhlY3V0b3Jzc3VtbWFyeRIfCgluZXh0dG9rZW4Y'
    '/oS6ZyABKAlSCW5leHR0b2tlbhIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb25pZA==');

@$core.Deprecated('Use listNamedQueriesInputDescriptor instead')
const ListNamedQueriesInput$json = {
  '1': 'ListNamedQueriesInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ListNamedQueriesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNamedQueriesInputDescriptor = $convert.base64Decode(
    'ChVMaXN0TmFtZWRRdWVyaWVzSW5wdXQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3'
    'VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIAoJd29ya2dyb3VwGIStofEB'
    'IAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use listNamedQueriesOutputDescriptor instead')
const ListNamedQueriesOutput$json = {
  '1': 'ListNamedQueriesOutput',
  '2': [
    {'1': 'namedqueryids', '3': 6092797, '4': 3, '5': 9, '10': 'namedqueryids'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListNamedQueriesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNamedQueriesOutputDescriptor =
    $convert.base64Decode(
        'ChZMaXN0TmFtZWRRdWVyaWVzT3V0cHV0EicKDW5hbWVkcXVlcnlpZHMY/e/zAiADKAlSDW5hbW'
        'VkcXVlcnlpZHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listNotebookMetadataInputDescriptor instead')
const ListNotebookMetadataInput$json = {
  '1': 'ListNotebookMetadataInput',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 1,
      '5': 11,
      '6': '.athena.FilterDefinition',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ListNotebookMetadataInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNotebookMetadataInputDescriptor = $convert.base64Decode(
    'ChlMaXN0Tm90ZWJvb2tNZXRhZGF0YUlucHV0EjUKB2ZpbHRlcnMY7c3qWSABKAsyGC5hdGhlbm'
    'EuRmlsdGVyRGVmaW5pdGlvblIHZmlsdGVycxIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4'
    'cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIgCgl3b3JrZ3JvdXAYhK'
    '2h8QEgASgJUgl3b3JrZ3JvdXA=');

@$core.Deprecated('Use listNotebookMetadataOutputDescriptor instead')
const ListNotebookMetadataOutput$json = {
  '1': 'ListNotebookMetadataOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'notebookmetadatalist',
      '3': 319238242,
      '4': 3,
      '5': 11,
      '6': '.athena.NotebookMetadata',
      '10': 'notebookmetadatalist'
    },
  ],
};

/// Descriptor for `ListNotebookMetadataOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNotebookMetadataOutputDescriptor =
    $convert.base64Decode(
        'ChpMaXN0Tm90ZWJvb2tNZXRhZGF0YU91dHB1dBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leH'
        'R0b2tlbhJQChRub3RlYm9va21ldGFkYXRhbGlzdBji4JyYASADKAsyGC5hdGhlbmEuTm90ZWJv'
        'b2tNZXRhZGF0YVIUbm90ZWJvb2ttZXRhZGF0YWxpc3Q=');

@$core.Deprecated('Use listNotebookSessionsRequestDescriptor instead')
const ListNotebookSessionsRequest$json = {
  '1': 'ListNotebookSessionsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
  ],
};

/// Descriptor for `ListNotebookSessionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNotebookSessionsRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0Tm90ZWJvb2tTZXNzaW9uc1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm'
        '1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIQoKbm90ZWJvb2tp'
        'ZBjetJVLIAEoCVIKbm90ZWJvb2tpZA==');

@$core.Deprecated('Use listNotebookSessionsResponseDescriptor instead')
const ListNotebookSessionsResponse$json = {
  '1': 'ListNotebookSessionsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'notebooksessionslist',
      '3': 247821802,
      '4': 3,
      '5': 11,
      '6': '.athena.NotebookSessionSummary',
      '10': 'notebooksessionslist'
    },
  ],
};

/// Descriptor for `ListNotebookSessionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNotebookSessionsResponseDescriptor =
    $convert.base64Decode(
        'ChxMaXN0Tm90ZWJvb2tTZXNzaW9uc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbm'
        'V4dHRva2VuElUKFG5vdGVib29rc2Vzc2lvbnNsaXN0GOrrlXYgAygLMh4uYXRoZW5hLk5vdGVi'
        'b29rU2Vzc2lvblN1bW1hcnlSFG5vdGVib29rc2Vzc2lvbnNsaXN0');

@$core.Deprecated('Use listPreparedStatementsInputDescriptor instead')
const ListPreparedStatementsInput$json = {
  '1': 'ListPreparedStatementsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ListPreparedStatementsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPreparedStatementsInputDescriptor =
    $convert.base64Decode(
        'ChtMaXN0UHJlcGFyZWRTdGF0ZW1lbnRzSW5wdXQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm'
        '1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIAoJd29ya2dyb3Vw'
        'GIStofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use listPreparedStatementsOutputDescriptor instead')
const ListPreparedStatementsOutput$json = {
  '1': 'ListPreparedStatementsOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'preparedstatements',
      '3': 526923667,
      '4': 3,
      '5': 11,
      '6': '.athena.PreparedStatementSummary',
      '10': 'preparedstatements'
    },
  ],
};

/// Descriptor for `ListPreparedStatementsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPreparedStatementsOutputDescriptor =
    $convert.base64Decode(
        'ChxMaXN0UHJlcGFyZWRTdGF0ZW1lbnRzT3V0cHV0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbm'
        'V4dHRva2VuElQKEnByZXBhcmVkc3RhdGVtZW50cxiT76D7ASADKAsyIC5hdGhlbmEuUHJlcGFy'
        'ZWRTdGF0ZW1lbnRTdW1tYXJ5UhJwcmVwYXJlZHN0YXRlbWVudHM=');

@$core.Deprecated('Use listQueryExecutionsInputDescriptor instead')
const ListQueryExecutionsInput$json = {
  '1': 'ListQueryExecutionsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ListQueryExecutionsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueryExecutionsInputDescriptor = $convert.base64Decode(
    'ChhMaXN0UXVlcnlFeGVjdXRpb25zSW5wdXQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heH'
    'Jlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIAoJd29ya2dyb3VwGISt'
    'ofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use listQueryExecutionsOutputDescriptor instead')
const ListQueryExecutionsOutput$json = {
  '1': 'ListQueryExecutionsOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'queryexecutionids',
      '3': 493941192,
      '4': 3,
      '5': 9,
      '10': 'queryexecutionids'
    },
  ],
};

/// Descriptor for `ListQueryExecutionsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueryExecutionsOutputDescriptor = $convert.base64Decode(
    'ChlMaXN0UXVlcnlFeGVjdXRpb25zT3V0cHV0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dH'
    'Rva2VuEjAKEXF1ZXJ5ZXhlY3V0aW9uaWRzGMjjw+sBIAMoCVIRcXVlcnlleGVjdXRpb25pZHM=');

@$core.Deprecated('Use listSessionsRequestDescriptor instead')
const ListSessionsRequest$json = {
  '1': 'ListSessionsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'statefilter',
      '3': 184693297,
      '4': 1,
      '5': 14,
      '6': '.athena.SessionState',
      '10': 'statefilter'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ListSessionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSessionsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0U2Vzc2lvbnNSZXF1ZXN0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bH'
    'RzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEjkKC3N0YXRlZmlsdGVyGLHkiFgg'
    'ASgOMhQuYXRoZW5hLlNlc3Npb25TdGF0ZVILc3RhdGVmaWx0ZXISIAoJd29ya2dyb3VwGIStof'
    'EBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use listSessionsResponseDescriptor instead')
const ListSessionsResponse$json = {
  '1': 'ListSessionsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'sessions',
      '3': 379226861,
      '4': 3,
      '5': 11,
      '6': '.athena.SessionSummary',
      '10': 'sessions'
    },
  ],
};

/// Descriptor for `ListSessionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSessionsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0U2Vzc2lvbnNSZXNwb25zZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbh'
    'I2CghzZXNzaW9ucxjtleq0ASADKAsyFi5hdGhlbmEuU2Vzc2lvblN1bW1hcnlSCHNlc3Npb25z');

@$core.Deprecated('Use listTableMetadataInputDescriptor instead')
const ListTableMetadataInput$json = {
  '1': 'ListTableMetadataInput',
  '2': [
    {'1': 'catalogname', '3': 518825212, '4': 1, '5': 9, '10': 'catalogname'},
    {'1': 'databasename', '3': 89545052, '4': 1, '5': 9, '10': 'databasename'},
    {'1': 'expression', '3': 193051916, '4': 1, '5': 9, '10': 'expression'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `ListTableMetadataInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTableMetadataInputDescriptor = $convert.base64Decode(
    'ChZMaXN0VGFibGVNZXRhZGF0YUlucHV0EiQKC2NhdGFsb2duYW1lGPzJsvcBIAEoCVILY2F0YW'
    'xvZ25hbWUSJQoMZGF0YWJhc2VuYW1lGNyy2SogASgJUgxkYXRhYmFzZW5hbWUSIQoKZXhwcmVz'
    'c2lvbhiM+oZcIAEoCVIKZXhwcmVzc2lvbhIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cm'
    'VzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIgCgl3b3JrZ3JvdXAYhK2h'
    '8QEgASgJUgl3b3JrZ3JvdXA=');

@$core.Deprecated('Use listTableMetadataOutputDescriptor instead')
const ListTableMetadataOutput$json = {
  '1': 'ListTableMetadataOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'tablemetadatalist',
      '3': 532092415,
      '4': 3,
      '5': 11,
      '6': '.athena.TableMetadata',
      '10': 'tablemetadatalist'
    },
  ],
};

/// Descriptor for `ListTableMetadataOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTableMetadataOutputDescriptor = $convert.base64Decode(
    'ChdMaXN0VGFibGVNZXRhZGF0YU91dHB1dBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2'
    'tlbhJHChF0YWJsZW1ldGFkYXRhbGlzdBj/q9z9ASADKAsyFS5hdGhlbmEuVGFibGVNZXRhZGF0'
    'YVIRdGFibGVtZXRhZGF0YWxpc3Q=');

@$core.Deprecated('Use listTagsForResourceInputDescriptor instead')
const ListTagsForResourceInput$json = {
  '1': 'ListTagsForResourceInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceInputDescriptor = $convert.base64Decode(
    'ChhMaXN0VGFnc0ZvclJlc291cmNlSW5wdXQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heH'
    'Jlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SJAoLcmVzb3VyY2Vhcm4Y'
    '7cCZsAEgASgJUgtyZXNvdXJjZWFybg==');

@$core.Deprecated('Use listTagsForResourceOutputDescriptor instead')
const ListTagsForResourceOutput$json = {
  '1': 'ListTagsForResourceOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.athena.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceOutputDescriptor =
    $convert.base64Decode(
        'ChlMaXN0VGFnc0ZvclJlc291cmNlT3V0cHV0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dH'
        'Rva2VuEiMKBHRhZ3MYwcH2tQEgAygLMgsuYXRoZW5hLlRhZ1IEdGFncw==');

@$core.Deprecated('Use listWorkGroupsInputDescriptor instead')
const ListWorkGroupsInput$json = {
  '1': 'ListWorkGroupsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListWorkGroupsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWorkGroupsInputDescriptor = $convert.base64Decode(
    'ChNMaXN0V29ya0dyb3Vwc0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bH'
    'RzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listWorkGroupsOutputDescriptor instead')
const ListWorkGroupsOutput$json = {
  '1': 'ListWorkGroupsOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'workgroups',
      '3': 159439825,
      '4': 3,
      '5': 11,
      '6': '.athena.WorkGroupSummary',
      '10': 'workgroups'
    },
  ],
};

/// Descriptor for `ListWorkGroupsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWorkGroupsOutputDescriptor = $convert.base64Decode(
    'ChRMaXN0V29ya0dyb3Vwc091dHB1dBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbh'
    'I7Cgp3b3JrZ3JvdXBzGNG3g0wgAygLMhguYXRoZW5hLldvcmtHcm91cFN1bW1hcnlSCndvcmtn'
    'cm91cHM=');

@$core.Deprecated('Use managedLoggingConfigurationDescriptor instead')
const ManagedLoggingConfiguration$json = {
  '1': 'ManagedLoggingConfiguration',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'kmskey', '3': 114561194, '4': 1, '5': 9, '10': 'kmskey'},
  ],
};

/// Descriptor for `ManagedLoggingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedLoggingConfigurationDescriptor =
    $convert.base64Decode(
        'ChtNYW5hZ2VkTG9nZ2luZ0NvbmZpZ3VyYXRpb24SHAoHZW5hYmxlZBi/yJvkASABKAhSB2VuYW'
        'JsZWQSGQoGa21za2V5GKqh0DYgASgJUgZrbXNrZXk=');

@$core.Deprecated('Use managedQueryResultsConfigurationDescriptor instead')
const ManagedQueryResultsConfiguration$json = {
  '1': 'ManagedQueryResultsConfiguration',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.athena.ManagedQueryResultsEncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
  ],
};

/// Descriptor for `ManagedQueryResultsConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedQueryResultsConfigurationDescriptor =
    $convert.base64Decode(
        'CiBNYW5hZ2VkUXVlcnlSZXN1bHRzQ29uZmlndXJhdGlvbhIcCgdlbmFibGVkGL/Im+QBIAEoCF'
        'IHZW5hYmxlZBJvChdlbmNyeXB0aW9uY29uZmlndXJhdGlvbhj3xtNrIAEoCzIyLmF0aGVuYS5N'
        'YW5hZ2VkUXVlcnlSZXN1bHRzRW5jcnlwdGlvbkNvbmZpZ3VyYXRpb25SF2VuY3J5cHRpb25jb2'
        '5maWd1cmF0aW9u');

@$core
    .Deprecated('Use managedQueryResultsConfigurationUpdatesDescriptor instead')
const ManagedQueryResultsConfigurationUpdates$json = {
  '1': 'ManagedQueryResultsConfigurationUpdates',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.athena.ManagedQueryResultsEncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {
      '1': 'removeencryptionconfiguration',
      '3': 294004335,
      '4': 1,
      '5': 8,
      '10': 'removeencryptionconfiguration'
    },
  ],
};

/// Descriptor for `ManagedQueryResultsConfigurationUpdates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List managedQueryResultsConfigurationUpdatesDescriptor =
    $convert.base64Decode(
        'CidNYW5hZ2VkUXVlcnlSZXN1bHRzQ29uZmlndXJhdGlvblVwZGF0ZXMSHAoHZW5hYmxlZBi/yJ'
        'vkASABKAhSB2VuYWJsZWQSbwoXZW5jcnlwdGlvbmNvbmZpZ3VyYXRpb24Y98bTayABKAsyMi5h'
        'dGhlbmEuTWFuYWdlZFF1ZXJ5UmVzdWx0c0VuY3J5cHRpb25Db25maWd1cmF0aW9uUhdlbmNyeX'
        'B0aW9uY29uZmlndXJhdGlvbhJICh1yZW1vdmVlbmNyeXB0aW9uY29uZmlndXJhdGlvbhjvzJiM'
        'ASABKAhSHXJlbW92ZWVuY3J5cHRpb25jb25maWd1cmF0aW9u');

@$core.Deprecated(
    'Use managedQueryResultsEncryptionConfigurationDescriptor instead')
const ManagedQueryResultsEncryptionConfiguration$json = {
  '1': 'ManagedQueryResultsEncryptionConfiguration',
  '2': [
    {'1': 'kmskey', '3': 114561194, '4': 1, '5': 9, '10': 'kmskey'},
  ],
};

/// Descriptor for `ManagedQueryResultsEncryptionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    managedQueryResultsEncryptionConfigurationDescriptor =
    $convert.base64Decode(
        'CipNYW5hZ2VkUXVlcnlSZXN1bHRzRW5jcnlwdGlvbkNvbmZpZ3VyYXRpb24SGQoGa21za2V5GK'
        'qh0DYgASgJUgZrbXNrZXk=');

@$core.Deprecated('Use metadataExceptionDescriptor instead')
const MetadataException$json = {
  '1': 'MetadataException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MetadataException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataExceptionDescriptor = $convert.base64Decode(
    'ChFNZXRhZGF0YUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use monitoringConfigurationDescriptor instead')
const MonitoringConfiguration$json = {
  '1': 'MonitoringConfiguration',
  '2': [
    {
      '1': 'cloudwatchloggingconfiguration',
      '3': 60618733,
      '4': 1,
      '5': 11,
      '6': '.athena.CloudWatchLoggingConfiguration',
      '10': 'cloudwatchloggingconfiguration'
    },
    {
      '1': 'managedloggingconfiguration',
      '3': 270845728,
      '4': 1,
      '5': 11,
      '6': '.athena.ManagedLoggingConfiguration',
      '10': 'managedloggingconfiguration'
    },
    {
      '1': 's3loggingconfiguration',
      '3': 15691207,
      '4': 1,
      '5': 11,
      '6': '.athena.S3LoggingConfiguration',
      '10': 's3loggingconfiguration'
    },
  ],
};

/// Descriptor for `MonitoringConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List monitoringConfigurationDescriptor = $convert.base64Decode(
    'ChdNb25pdG9yaW5nQ29uZmlndXJhdGlvbhJxCh5jbG91ZHdhdGNobG9nZ2luZ2NvbmZpZ3VyYX'
    'Rpb24Y7e/zHCABKAsyJi5hdGhlbmEuQ2xvdWRXYXRjaExvZ2dpbmdDb25maWd1cmF0aW9uUh5j'
    'bG91ZHdhdGNobG9nZ2luZ2NvbmZpZ3VyYXRpb24SaQobbWFuYWdlZGxvZ2dpbmdjb25maWd1cm'
    'F0aW9uGKCOk4EBIAEoCzIjLmF0aGVuYS5NYW5hZ2VkTG9nZ2luZ0NvbmZpZ3VyYXRpb25SG21h'
    'bmFnZWRsb2dnaW5nY29uZmlndXJhdGlvbhJZChZzM2xvZ2dpbmdjb25maWd1cmF0aW9uGMfbvQ'
    'cgASgLMh4uYXRoZW5hLlMzTG9nZ2luZ0NvbmZpZ3VyYXRpb25SFnMzbG9nZ2luZ2NvbmZpZ3Vy'
    'YXRpb24=');

@$core.Deprecated('Use namedQueryDescriptor instead')
const NamedQuery$json = {
  '1': 'NamedQuery',
  '2': [
    {'1': 'database', '3': 278147289, '4': 1, '5': 9, '10': 'database'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'namedqueryid', '3': 330896872, '4': 1, '5': 9, '10': 'namedqueryid'},
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `NamedQuery`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List namedQueryDescriptor = $convert.base64Decode(
    'CgpOYW1lZFF1ZXJ5Eh4KCGRhdGFiYXNlGNnh0IQBIAEoCVIIZGF0YWJhc2USIwoLZGVzY3JpcH'
    'Rpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSJgoMbmFt'
    'ZWRxdWVyeWlkGOir5J0BIAEoCVIMbmFtZWRxdWVyeWlkEiQKC3F1ZXJ5c3RyaW5nGOfK788BIA'
    'EoCVILcXVlcnlzdHJpbmcSIAoJd29ya2dyb3VwGIStofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use notebookMetadataDescriptor instead')
const NotebookMetadata$json = {
  '1': 'NotebookMetadata',
  '2': [
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.athena.NotebookType',
      '10': 'type'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `NotebookMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notebookMetadataDescriptor = $convert.base64Decode(
    'ChBOb3RlYm9va01ldGFkYXRhEiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3JlYXRpb250aW'
    '1lEi0KEGxhc3Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWUSFQoEbmFt'
    'ZRiH5oF/IAEoCVIEbmFtZRIhCgpub3RlYm9va2lkGN60lUsgASgJUgpub3RlYm9va2lkEiwKBH'
    'R5cGUY7qDXigEgASgOMhQuYXRoZW5hLk5vdGVib29rVHlwZVIEdHlwZRIgCgl3b3JrZ3JvdXAY'
    'hK2h8QEgASgJUgl3b3JrZ3JvdXA=');

@$core.Deprecated('Use notebookSessionSummaryDescriptor instead')
const NotebookSessionSummary$json = {
  '1': 'NotebookSessionSummary',
  '2': [
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `NotebookSessionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notebookSessionSummaryDescriptor =
    $convert.base64Decode(
        'ChZOb3RlYm9va1Nlc3Npb25TdW1tYXJ5EiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3JlYX'
        'Rpb250aW1lEh8KCXNlc3Npb25pZBi7hOUJIAEoCVIJc2Vzc2lvbmlk');

@$core.Deprecated('Use preparedStatementDescriptor instead')
const PreparedStatement$json = {
  '1': 'PreparedStatement',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'querystatement',
      '3': 340852217,
      '4': 1,
      '5': 9,
      '10': 'querystatement'
    },
    {
      '1': 'statementname',
      '3': 23047926,
      '4': 1,
      '5': 9,
      '10': 'statementname'
    },
    {
      '1': 'workgroupname',
      '3': 258526185,
      '4': 1,
      '5': 9,
      '10': 'workgroupname'
    },
  ],
};

/// Descriptor for `PreparedStatement`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List preparedStatementDescriptor = $convert.base64Decode(
    'ChFQcmVwYXJlZFN0YXRlbWVudBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb2'
    '4SLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZRIqCg5xdWVy'
    'eXN0YXRlbWVudBj5+8OiASABKAlSDnF1ZXJ5c3RhdGVtZW50EicKDXN0YXRlbWVudG5hbWUY9t'
    '3+CiABKAlSDXN0YXRlbWVudG5hbWUSJwoNd29ya2dyb3VwbmFtZRjpl6N7IAEoCVINd29ya2dy'
    'b3VwbmFtZQ==');

@$core.Deprecated('Use preparedStatementSummaryDescriptor instead')
const PreparedStatementSummary$json = {
  '1': 'PreparedStatementSummary',
  '2': [
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'statementname',
      '3': 23047926,
      '4': 1,
      '5': 9,
      '10': 'statementname'
    },
  ],
};

/// Descriptor for `PreparedStatementSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List preparedStatementSummaryDescriptor = $convert.base64Decode(
    'ChhQcmVwYXJlZFN0YXRlbWVudFN1bW1hcnkSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCV'
    'IQbGFzdG1vZGlmaWVkdGltZRInCg1zdGF0ZW1lbnRuYW1lGPbd/gogASgJUg1zdGF0ZW1lbnRu'
    'YW1l');

@$core
    .Deprecated('Use putCapacityAssignmentConfigurationInputDescriptor instead')
const PutCapacityAssignmentConfigurationInput$json = {
  '1': 'PutCapacityAssignmentConfigurationInput',
  '2': [
    {
      '1': 'capacityassignments',
      '3': 345772294,
      '4': 3,
      '5': 11,
      '6': '.athena.CapacityAssignment',
      '10': 'capacityassignments'
    },
    {
      '1': 'capacityreservationname',
      '3': 327567687,
      '4': 1,
      '5': 9,
      '10': 'capacityreservationname'
    },
  ],
};

/// Descriptor for `PutCapacityAssignmentConfigurationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putCapacityAssignmentConfigurationInputDescriptor =
    $convert.base64Decode(
        'CidQdXRDYXBhY2l0eUFzc2lnbm1lbnRDb25maWd1cmF0aW9uSW5wdXQSUAoTY2FwYWNpdHlhc3'
        'NpZ25tZW50cxiGovCkASADKAsyGi5hdGhlbmEuQ2FwYWNpdHlBc3NpZ25tZW50UhNjYXBhY2l0'
        'eWFzc2lnbm1lbnRzEjwKF2NhcGFjaXR5cmVzZXJ2YXRpb25uYW1lGMeSmZwBIAEoCVIXY2FwYW'
        'NpdHlyZXNlcnZhdGlvbm5hbWU=');

@$core.Deprecated(
    'Use putCapacityAssignmentConfigurationOutputDescriptor instead')
const PutCapacityAssignmentConfigurationOutput$json = {
  '1': 'PutCapacityAssignmentConfigurationOutput',
};

/// Descriptor for `PutCapacityAssignmentConfigurationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putCapacityAssignmentConfigurationOutputDescriptor =
    $convert.base64Decode(
        'CihQdXRDYXBhY2l0eUFzc2lnbm1lbnRDb25maWd1cmF0aW9uT3V0cHV0');

@$core.Deprecated('Use queryExecutionDescriptor instead')
const QueryExecution$json = {
  '1': 'QueryExecution',
  '2': [
    {
      '1': 'engineversion',
      '3': 44953462,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineVersion',
      '10': 'engineversion'
    },
    {
      '1': 'executionparameters',
      '3': 527591242,
      '4': 3,
      '5': 9,
      '10': 'executionparameters'
    },
    {
      '1': 'managedqueryresultsconfiguration',
      '3': 159215683,
      '4': 1,
      '5': 11,
      '6': '.athena.ManagedQueryResultsConfiguration',
      '10': 'managedqueryresultsconfiguration'
    },
    {'1': 'query', '3': 512354180, '4': 1, '5': 9, '10': 'query'},
    {
      '1': 'queryexecutioncontext',
      '3': 139302379,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryExecutionContext',
      '10': 'queryexecutioncontext'
    },
    {
      '1': 'queryexecutionid',
      '3': 467615503,
      '4': 1,
      '5': 9,
      '10': 'queryexecutionid'
    },
    {
      '1': 'queryresultss3accessgrantsconfiguration',
      '3': 264403639,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryResultsS3AccessGrantsConfiguration',
      '10': 'queryresultss3accessgrantsconfiguration'
    },
    {
      '1': 'resultconfiguration',
      '3': 183031503,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultConfiguration',
      '10': 'resultconfiguration'
    },
    {
      '1': 'resultreuseconfiguration',
      '3': 482796543,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultReuseConfiguration',
      '10': 'resultreuseconfiguration'
    },
    {
      '1': 'statementtype',
      '3': 286463655,
      '4': 1,
      '5': 14,
      '6': '.athena.StatementType',
      '10': 'statementtype'
    },
    {
      '1': 'statistics',
      '3': 510636075,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryExecutionStatistics',
      '10': 'statistics'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryExecutionStatus',
      '10': 'status'
    },
    {
      '1': 'substatementtype',
      '3': 504257527,
      '4': 1,
      '5': 9,
      '10': 'substatementtype'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `QueryExecution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryExecutionDescriptor = $convert.base64Decode(
    'Cg5RdWVyeUV4ZWN1dGlvbhI+Cg1lbmdpbmV2ZXJzaW9uGPbetxUgASgLMhUuYXRoZW5hLkVuZ2'
    'luZVZlcnNpb25SDWVuZ2luZXZlcnNpb24SNAoTZXhlY3V0aW9ucGFyYW1ldGVycxjKzsn7ASAD'
    'KAlSE2V4ZWN1dGlvbnBhcmFtZXRlcnMSdwogbWFuYWdlZHF1ZXJ5cmVzdWx0c2NvbmZpZ3VyYX'
    'Rpb24Yw+D1SyABKAsyKC5hdGhlbmEuTWFuYWdlZFF1ZXJ5UmVzdWx0c0NvbmZpZ3VyYXRpb25S'
    'IG1hbmFnZWRxdWVyeXJlc3VsdHNjb25maWd1cmF0aW9uEhgKBXF1ZXJ5GITPp/QBIAEoCVIFcX'
    'VlcnkSVgoVcXVlcnlleGVjdXRpb25jb250ZXh0GOurtkIgASgLMh0uYXRoZW5hLlF1ZXJ5RXhl'
    'Y3V0aW9uQ29udGV4dFIVcXVlcnlleGVjdXRpb25jb250ZXh0Ei4KEHF1ZXJ5ZXhlY3V0aW9uaW'
    'QYj/783gEgASgJUhBxdWVyeWV4ZWN1dGlvbmlkEowBCidxdWVyeXJlc3VsdHNzM2FjY2Vzc2dy'
    'YW50c2NvbmZpZ3VyYXRpb24Yt/WJfiABKAsyLy5hdGhlbmEuUXVlcnlSZXN1bHRzUzNBY2Nlc3'
    'NHcmFudHNDb25maWd1cmF0aW9uUidxdWVyeXJlc3VsdHNzM2FjY2Vzc2dyYW50c2NvbmZpZ3Vy'
    'YXRpb24SUAoTcmVzdWx0Y29uZmlndXJhdGlvbhjPraNXIAEoCzIbLmF0aGVuYS5SZXN1bHRDb2'
    '5maWd1cmF0aW9uUhNyZXN1bHRjb25maWd1cmF0aW9uEmAKGHJlc3VsdHJldXNlY29uZmlndXJh'
    'dGlvbhj/x5vmASABKAsyIC5hdGhlbmEuUmVzdWx0UmV1c2VDb25maWd1cmF0aW9uUhhyZXN1bH'
    'RyZXVzZWNvbmZpZ3VyYXRpb24SPwoNc3RhdGVtZW50dHlwZRinrcyIASABKA4yFS5hdGhlbmEu'
    'U3RhdGVtZW50VHlwZVINc3RhdGVtZW50dHlwZRJECgpzdGF0aXN0aWNzGKvgvvMBIAEoCzIgLm'
    'F0aGVuYS5RdWVyeUV4ZWN1dGlvblN0YXRpc3RpY3NSCnN0YXRpc3RpY3MSNwoGc3RhdHVzGJDk'
    '+wIgASgLMhwuYXRoZW5hLlF1ZXJ5RXhlY3V0aW9uU3RhdHVzUgZzdGF0dXMSLgoQc3Vic3RhdG'
    'VtZW50dHlwZRj3t7nwASABKAlSEHN1YnN0YXRlbWVudHR5cGUSIAoJd29ya2dyb3VwGIStofEB'
    'IAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use queryExecutionContextDescriptor instead')
const QueryExecutionContext$json = {
  '1': 'QueryExecutionContext',
  '2': [
    {'1': 'catalog', '3': 111213497, '4': 1, '5': 9, '10': 'catalog'},
    {'1': 'database', '3': 278147289, '4': 1, '5': 9, '10': 'database'},
  ],
};

/// Descriptor for `QueryExecutionContext`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryExecutionContextDescriptor = $convert.base64Decode(
    'ChVRdWVyeUV4ZWN1dGlvbkNvbnRleHQSGwoHY2F0YWxvZxi594M1IAEoCVIHY2F0YWxvZxIeCg'
    'hkYXRhYmFzZRjZ4dCEASABKAlSCGRhdGFiYXNl');

@$core.Deprecated('Use queryExecutionStatisticsDescriptor instead')
const QueryExecutionStatistics$json = {
  '1': 'QueryExecutionStatistics',
  '2': [
    {
      '1': 'datamanifestlocation',
      '3': 329970068,
      '4': 1,
      '5': 9,
      '10': 'datamanifestlocation'
    },
    {
      '1': 'datascannedinbytes',
      '3': 125540604,
      '4': 1,
      '5': 3,
      '10': 'datascannedinbytes'
    },
    {'1': 'dpucount', '3': 296377724, '4': 1, '5': 1, '10': 'dpucount'},
    {
      '1': 'engineexecutiontimeinmillis',
      '3': 228680528,
      '4': 1,
      '5': 3,
      '10': 'engineexecutiontimeinmillis'
    },
    {
      '1': 'queryplanningtimeinmillis',
      '3': 482562295,
      '4': 1,
      '5': 3,
      '10': 'queryplanningtimeinmillis'
    },
    {
      '1': 'queryqueuetimeinmillis',
      '3': 221814191,
      '4': 1,
      '5': 3,
      '10': 'queryqueuetimeinmillis'
    },
    {
      '1': 'resultreuseinformation',
      '3': 261773827,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultReuseInformation',
      '10': 'resultreuseinformation'
    },
    {
      '1': 'servicepreprocessingtimeinmillis',
      '3': 185806899,
      '4': 1,
      '5': 3,
      '10': 'servicepreprocessingtimeinmillis'
    },
    {
      '1': 'serviceprocessingtimeinmillis',
      '3': 81556436,
      '4': 1,
      '5': 3,
      '10': 'serviceprocessingtimeinmillis'
    },
    {
      '1': 'totalexecutiontimeinmillis',
      '3': 251249322,
      '4': 1,
      '5': 3,
      '10': 'totalexecutiontimeinmillis'
    },
  ],
};

/// Descriptor for `QueryExecutionStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryExecutionStatisticsDescriptor = $convert.base64Decode(
    'ChhRdWVyeUV4ZWN1dGlvblN0YXRpc3RpY3MSNgoUZGF0YW1hbmlmZXN0bG9jYXRpb24YlOOrnQ'
    'EgASgJUhRkYXRhbWFuaWZlc3Rsb2NhdGlvbhIxChJkYXRhc2Nhbm5lZGluYnl0ZXMY/LHuOyAB'
    'KANSEmRhdGFzY2FubmVkaW5ieXRlcxIeCghkcHVjb3VudBj8uqmNASABKAFSCGRwdWNvdW50Ek'
    'MKG2VuZ2luZWV4ZWN1dGlvbnRpbWVpbm1pbGxpcxjQxoVtIAEoA1IbZW5naW5lZXhlY3V0aW9u'
    'dGltZWlubWlsbGlzEkAKGXF1ZXJ5cGxhbm5pbmd0aW1laW5taWxsaXMY96GN5gEgASgDUhlxdW'
    'VyeXBsYW5uaW5ndGltZWlubWlsbGlzEjkKFnF1ZXJ5cXVldWV0aW1laW5taWxsaXMYr7viaSAB'
    'KANSFnF1ZXJ5cXVldWV0aW1laW5taWxsaXMSWQoWcmVzdWx0cmV1c2VpbmZvcm1hdGlvbhiDtO'
    'l8IAEoCzIeLmF0aGVuYS5SZXN1bHRSZXVzZUluZm9ybWF0aW9uUhZyZXN1bHRyZXVzZWluZm9y'
    'bWF0aW9uEk0KIHNlcnZpY2VwcmVwcm9jZXNzaW5ndGltZWlubWlsbGlzGLPgzFggASgDUiBzZX'
    'J2aWNlcHJlcHJvY2Vzc2luZ3RpbWVpbm1pbGxpcxJHCh1zZXJ2aWNlcHJvY2Vzc2luZ3RpbWVp'
    'bm1pbGxpcxjU5/EmIAEoA1Idc2VydmljZXByb2Nlc3Npbmd0aW1laW5taWxsaXMSQQoadG90YW'
    'xleGVjdXRpb250aW1laW5taWxsaXMYqoXndyABKANSGnRvdGFsZXhlY3V0aW9udGltZWlubWls'
    'bGlz');

@$core.Deprecated('Use queryExecutionStatusDescriptor instead')
const QueryExecutionStatus$json = {
  '1': 'QueryExecutionStatus',
  '2': [
    {
      '1': 'athenaerror',
      '3': 82066385,
      '4': 1,
      '5': 11,
      '6': '.athena.AthenaError',
      '10': 'athenaerror'
    },
    {
      '1': 'completiondatetime',
      '3': 175822779,
      '4': 1,
      '5': 9,
      '10': 'completiondatetime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.QueryExecutionState',
      '10': 'state'
    },
    {
      '1': 'statechangereason',
      '3': 228940439,
      '4': 1,
      '5': 9,
      '10': 'statechangereason'
    },
    {
      '1': 'submissiondatetime',
      '3': 449650437,
      '4': 1,
      '5': 9,
      '10': 'submissiondatetime'
    },
  ],
};

/// Descriptor for `QueryExecutionStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryExecutionStatusDescriptor = $convert.base64Decode(
    'ChRRdWVyeUV4ZWN1dGlvblN0YXR1cxI4CgthdGhlbmFlcnJvchjR95AnIAEoCzITLmF0aGVuYS'
    '5BdGhlbmFFcnJvclILYXRoZW5hZXJyb3ISMQoSY29tcGxldGlvbmRhdGV0aW1lGLuv61MgASgJ'
    'UhJjb21wbGV0aW9uZGF0ZXRpbWUSNQoFc3RhdGUYl8my7wEgASgOMhsuYXRoZW5hLlF1ZXJ5RX'
    'hlY3V0aW9uU3RhdGVSBXN0YXRlEi8KEXN0YXRlY2hhbmdlcmVhc29uGJe1lW0gASgJUhFzdGF0'
    'ZWNoYW5nZXJlYXNvbhIyChJzdWJtaXNzaW9uZGF0ZXRpbWUYhb601gEgASgJUhJzdWJtaXNzaW'
    '9uZGF0ZXRpbWU=');

@$core
    .Deprecated('Use queryResultsS3AccessGrantsConfigurationDescriptor instead')
const QueryResultsS3AccessGrantsConfiguration$json = {
  '1': 'QueryResultsS3AccessGrantsConfiguration',
  '2': [
    {
      '1': 'authenticationtype',
      '3': 277200874,
      '4': 1,
      '5': 14,
      '6': '.athena.AuthenticationType',
      '10': 'authenticationtype'
    },
    {
      '1': 'createuserlevelprefix',
      '3': 493230469,
      '4': 1,
      '5': 8,
      '10': 'createuserlevelprefix'
    },
    {
      '1': 'enables3accessgrants',
      '3': 96318870,
      '4': 1,
      '5': 8,
      '10': 'enables3accessgrants'
    },
  ],
};

/// Descriptor for `QueryResultsS3AccessGrantsConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryResultsS3AccessGrantsConfigurationDescriptor =
    $convert.base64Decode(
        'CidRdWVyeVJlc3VsdHNTM0FjY2Vzc0dyYW50c0NvbmZpZ3VyYXRpb24STgoSYXV0aGVudGljYX'
        'Rpb250eXBlGOr/loQBIAEoDjIaLmF0aGVuYS5BdXRoZW50aWNhdGlvblR5cGVSEmF1dGhlbnRp'
        'Y2F0aW9udHlwZRI4ChVjcmVhdGV1c2VybGV2ZWxwcmVmaXgYhbOY6wEgASgIUhVjcmVhdGV1c2'
        'VybGV2ZWxwcmVmaXgSNQoUZW5hYmxlczNhY2Nlc3NncmFudHMYluv2LSABKAhSFGVuYWJsZXMz'
        'YWNjZXNzZ3JhbnRz');

@$core.Deprecated('Use queryRuntimeStatisticsDescriptor instead')
const QueryRuntimeStatistics$json = {
  '1': 'QueryRuntimeStatistics',
  '2': [
    {
      '1': 'outputstage',
      '3': 61483955,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryStage',
      '10': 'outputstage'
    },
    {
      '1': 'rows',
      '3': 174428857,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryRuntimeStatisticsRows',
      '10': 'rows'
    },
    {
      '1': 'timeline',
      '3': 282559479,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryRuntimeStatisticsTimeline',
      '10': 'timeline'
    },
  ],
};

/// Descriptor for `QueryRuntimeStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryRuntimeStatisticsDescriptor = $convert.base64Decode(
    'ChZRdWVyeVJ1bnRpbWVTdGF0aXN0aWNzEjcKC291dHB1dHN0YWdlGLPXqB0gASgLMhIuYXRoZW'
    '5hLlF1ZXJ5U3RhZ2VSC291dHB1dHN0YWdlEjkKBHJvd3MYuaWWUyABKAsyIi5hdGhlbmEuUXVl'
    'cnlSdW50aW1lU3RhdGlzdGljc1Jvd3NSBHJvd3MSRgoIdGltZWxpbmUY94fehgEgASgLMiYuYX'
    'RoZW5hLlF1ZXJ5UnVudGltZVN0YXRpc3RpY3NUaW1lbGluZVIIdGltZWxpbmU=');

@$core.Deprecated('Use queryRuntimeStatisticsRowsDescriptor instead')
const QueryRuntimeStatisticsRows$json = {
  '1': 'QueryRuntimeStatisticsRows',
  '2': [
    {'1': 'inputbytes', '3': 431684783, '4': 1, '5': 3, '10': 'inputbytes'},
    {'1': 'inputrows', '3': 30551127, '4': 1, '5': 3, '10': 'inputrows'},
    {'1': 'outputbytes', '3': 318971400, '4': 1, '5': 3, '10': 'outputbytes'},
    {'1': 'outputrows', '3': 138873322, '4': 1, '5': 3, '10': 'outputrows'},
  ],
};

/// Descriptor for `QueryRuntimeStatisticsRows`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryRuntimeStatisticsRowsDescriptor = $convert.base64Decode(
    'ChpRdWVyeVJ1bnRpbWVTdGF0aXN0aWNzUm93cxIiCgppbnB1dGJ5dGVzGK/5680BIAEoA1IKaW'
    '5wdXRieXRlcxIfCglpbnB1dHJvd3MY19jIDiABKANSCWlucHV0cm93cxIkCgtvdXRwdXRieXRl'
    'cxiIvIyYASABKANSC291dHB1dGJ5dGVzEiEKCm91dHB1dHJvd3MY6pOcQiABKANSCm91dHB1dH'
    'Jvd3M=');

@$core.Deprecated('Use queryRuntimeStatisticsTimelineDescriptor instead')
const QueryRuntimeStatisticsTimeline$json = {
  '1': 'QueryRuntimeStatisticsTimeline',
  '2': [
    {
      '1': 'engineexecutiontimeinmillis',
      '3': 228680528,
      '4': 1,
      '5': 3,
      '10': 'engineexecutiontimeinmillis'
    },
    {
      '1': 'queryplanningtimeinmillis',
      '3': 482562295,
      '4': 1,
      '5': 3,
      '10': 'queryplanningtimeinmillis'
    },
    {
      '1': 'queryqueuetimeinmillis',
      '3': 221814191,
      '4': 1,
      '5': 3,
      '10': 'queryqueuetimeinmillis'
    },
    {
      '1': 'servicepreprocessingtimeinmillis',
      '3': 185806899,
      '4': 1,
      '5': 3,
      '10': 'servicepreprocessingtimeinmillis'
    },
    {
      '1': 'serviceprocessingtimeinmillis',
      '3': 81556436,
      '4': 1,
      '5': 3,
      '10': 'serviceprocessingtimeinmillis'
    },
    {
      '1': 'totalexecutiontimeinmillis',
      '3': 251249322,
      '4': 1,
      '5': 3,
      '10': 'totalexecutiontimeinmillis'
    },
  ],
};

/// Descriptor for `QueryRuntimeStatisticsTimeline`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryRuntimeStatisticsTimelineDescriptor = $convert.base64Decode(
    'Ch5RdWVyeVJ1bnRpbWVTdGF0aXN0aWNzVGltZWxpbmUSQwobZW5naW5lZXhlY3V0aW9udGltZW'
    'lubWlsbGlzGNDGhW0gASgDUhtlbmdpbmVleGVjdXRpb250aW1laW5taWxsaXMSQAoZcXVlcnlw'
    'bGFubmluZ3RpbWVpbm1pbGxpcxj3oY3mASABKANSGXF1ZXJ5cGxhbm5pbmd0aW1laW5taWxsaX'
    'MSOQoWcXVlcnlxdWV1ZXRpbWVpbm1pbGxpcxivu+JpIAEoA1IWcXVlcnlxdWV1ZXRpbWVpbm1p'
    'bGxpcxJNCiBzZXJ2aWNlcHJlcHJvY2Vzc2luZ3RpbWVpbm1pbGxpcxiz4MxYIAEoA1Igc2Vydm'
    'ljZXByZXByb2Nlc3Npbmd0aW1laW5taWxsaXMSRwodc2VydmljZXByb2Nlc3Npbmd0aW1laW5t'
    'aWxsaXMY1OfxJiABKANSHXNlcnZpY2Vwcm9jZXNzaW5ndGltZWlubWlsbGlzEkEKGnRvdGFsZX'
    'hlY3V0aW9udGltZWlubWlsbGlzGKqF53cgASgDUhp0b3RhbGV4ZWN1dGlvbnRpbWVpbm1pbGxp'
    'cw==');

@$core.Deprecated('Use queryStageDescriptor instead')
const QueryStage$json = {
  '1': 'QueryStage',
  '2': [
    {
      '1': 'executiontime',
      '3': 379716053,
      '4': 1,
      '5': 3,
      '10': 'executiontime'
    },
    {'1': 'inputbytes', '3': 431684783, '4': 1, '5': 3, '10': 'inputbytes'},
    {'1': 'inputrows', '3': 30551127, '4': 1, '5': 3, '10': 'inputrows'},
    {'1': 'outputbytes', '3': 318971400, '4': 1, '5': 3, '10': 'outputbytes'},
    {'1': 'outputrows', '3': 138873322, '4': 1, '5': 3, '10': 'outputrows'},
    {
      '1': 'querystageplan',
      '3': 36313545,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryStagePlanNode',
      '10': 'querystageplan'
    },
    {'1': 'stageid', '3': 328165497, '4': 1, '5': 3, '10': 'stageid'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
    {
      '1': 'substages',
      '3': 114732779,
      '4': 3,
      '5': 11,
      '6': '.athena.QueryStage',
      '10': 'substages'
    },
  ],
};

/// Descriptor for `QueryStage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryStageDescriptor = $convert.base64Decode(
    'CgpRdWVyeVN0YWdlEigKDWV4ZWN1dGlvbnRpbWUY1YOItQEgASgDUg1leGVjdXRpb250aW1lEi'
    'IKCmlucHV0Ynl0ZXMYr/nrzQEgASgDUgppbnB1dGJ5dGVzEh8KCWlucHV0cm93cxjX2MgOIAEo'
    'A1IJaW5wdXRyb3dzEiQKC291dHB1dGJ5dGVzGIi8jJgBIAEoA1ILb3V0cHV0Ynl0ZXMSIQoKb3'
    'V0cHV0cm93cxjqk5xCIAEoA1IKb3V0cHV0cm93cxJFCg5xdWVyeXN0YWdlcGxhbhjJs6gRIAEo'
    'CzIaLmF0aGVuYS5RdWVyeVN0YWdlUGxhbk5vZGVSDnF1ZXJ5c3RhZ2VwbGFuEhwKB3N0YWdlaW'
    'QY+dC9nAEgASgDUgdzdGFnZWlkEhgKBXN0YXRlGJfJsu8BIAEoCVIFc3RhdGUSMwoJc3Vic3Rh'
    'Z2VzGOvd2jYgAygLMhIuYXRoZW5hLlF1ZXJ5U3RhZ2VSCXN1YnN0YWdlcw==');

@$core.Deprecated('Use queryStagePlanNodeDescriptor instead')
const QueryStagePlanNode$json = {
  '1': 'QueryStagePlanNode',
  '2': [
    {
      '1': 'children',
      '3': 188567027,
      '4': 3,
      '5': 11,
      '6': '.athena.QueryStagePlanNode',
      '10': 'children'
    },
    {'1': 'identifier', '3': 41865311, '4': 1, '5': 9, '10': 'identifier'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'remotesources',
      '3': 143743504,
      '4': 3,
      '5': 9,
      '10': 'remotesources'
    },
  ],
};

/// Descriptor for `QueryStagePlanNode`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryStagePlanNodeDescriptor = $convert.base64Decode(
    'ChJRdWVyeVN0YWdlUGxhbk5vZGUSOQoIY2hpbGRyZW4Y85v1WSADKAsyGi5hdGhlbmEuUXVlcn'
    'lTdGFnZVBsYW5Ob2RlUghjaGlsZHJlbhIhCgppZGVudGlmaWVyGN+g+xMgASgJUgppZGVudGlm'
    'aWVyEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSJwoNcmVtb3Rlc291cmNlcxiQtMVEIAMoCVINcm'
    'Vtb3Rlc291cmNlcw==');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'resourcename', '3': 269834071, '4': 1, '5': 9, '10': 'resourcename'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'USJgoMcmVzb3VyY2VuYW1lGNeu1YABIAEoCVIMcmVzb3VyY2VuYW1l');

@$core.Deprecated('Use resultConfigurationDescriptor instead')
const ResultConfiguration$json = {
  '1': 'ResultConfiguration',
  '2': [
    {
      '1': 'aclconfiguration',
      '3': 9179900,
      '4': 1,
      '5': 11,
      '6': '.athena.AclConfiguration',
      '10': 'aclconfiguration'
    },
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.athena.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'outputlocation',
      '3': 67991028,
      '4': 1,
      '5': 9,
      '10': 'outputlocation'
    },
  ],
};

/// Descriptor for `ResultConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultConfigurationDescriptor = $convert.base64Decode(
    'ChNSZXN1bHRDb25maWd1cmF0aW9uEkcKEGFjbGNvbmZpZ3VyYXRpb24Y/KWwBCABKAsyGC5hdG'
    'hlbmEuQWNsQ29uZmlndXJhdGlvblIQYWNsY29uZmlndXJhdGlvbhJcChdlbmNyeXB0aW9uY29u'
    'ZmlndXJhdGlvbhj3xtNrIAEoCzIfLmF0aGVuYS5FbmNyeXB0aW9uQ29uZmlndXJhdGlvblIXZW'
    '5jcnlwdGlvbmNvbmZpZ3VyYXRpb24SMwoTZXhwZWN0ZWRidWNrZXRvd25lchin3fw+IAEoCVIT'
    'ZXhwZWN0ZWRidWNrZXRvd25lchIpCg5vdXRwdXRsb2NhdGlvbhj067UgIAEoCVIOb3V0cHV0bG'
    '9jYXRpb24=');

@$core.Deprecated('Use resultConfigurationUpdatesDescriptor instead')
const ResultConfigurationUpdates$json = {
  '1': 'ResultConfigurationUpdates',
  '2': [
    {
      '1': 'aclconfiguration',
      '3': 9179900,
      '4': 1,
      '5': 11,
      '6': '.athena.AclConfiguration',
      '10': 'aclconfiguration'
    },
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.athena.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {
      '1': 'expectedbucketowner',
      '3': 132066983,
      '4': 1,
      '5': 9,
      '10': 'expectedbucketowner'
    },
    {
      '1': 'outputlocation',
      '3': 67991028,
      '4': 1,
      '5': 9,
      '10': 'outputlocation'
    },
    {
      '1': 'removeaclconfiguration',
      '3': 378505796,
      '4': 1,
      '5': 8,
      '10': 'removeaclconfiguration'
    },
    {
      '1': 'removeencryptionconfiguration',
      '3': 294004335,
      '4': 1,
      '5': 8,
      '10': 'removeencryptionconfiguration'
    },
    {
      '1': 'removeexpectedbucketowner',
      '3': 198850607,
      '4': 1,
      '5': 8,
      '10': 'removeexpectedbucketowner'
    },
    {
      '1': 'removeoutputlocation',
      '3': 197075308,
      '4': 1,
      '5': 8,
      '10': 'removeoutputlocation'
    },
  ],
};

/// Descriptor for `ResultConfigurationUpdates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultConfigurationUpdatesDescriptor = $convert.base64Decode(
    'ChpSZXN1bHRDb25maWd1cmF0aW9uVXBkYXRlcxJHChBhY2xjb25maWd1cmF0aW9uGPylsAQgAS'
    'gLMhguYXRoZW5hLkFjbENvbmZpZ3VyYXRpb25SEGFjbGNvbmZpZ3VyYXRpb24SXAoXZW5jcnlw'
    'dGlvbmNvbmZpZ3VyYXRpb24Y98bTayABKAsyHy5hdGhlbmEuRW5jcnlwdGlvbkNvbmZpZ3VyYX'
    'Rpb25SF2VuY3J5cHRpb25jb25maWd1cmF0aW9uEjMKE2V4cGVjdGVkYnVja2V0b3duZXIYp938'
    'PiABKAlSE2V4cGVjdGVkYnVja2V0b3duZXISKQoOb3V0cHV0bG9jYXRpb24Y9Ou1ICABKAlSDm'
    '91dHB1dGxvY2F0aW9uEjoKFnJlbW92ZWFjbGNvbmZpZ3VyYXRpb24YxJS+tAEgASgIUhZyZW1v'
    'dmVhY2xjb25maWd1cmF0aW9uEkgKHXJlbW92ZWVuY3J5cHRpb25jb25maWd1cmF0aW9uGO/MmI'
    'wBIAEoCFIdcmVtb3ZlZW5jcnlwdGlvbmNvbmZpZ3VyYXRpb24SPwoZcmVtb3ZlZXhwZWN0ZWRi'
    'dWNrZXRvd25lchiv8OheIAEoCFIZcmVtb3ZlZXhwZWN0ZWRidWNrZXRvd25lchI1ChRyZW1vdm'
    'VvdXRwdXRsb2NhdGlvbhjswvxdIAEoCFIUcmVtb3Zlb3V0cHV0bG9jYXRpb24=');

@$core.Deprecated('Use resultReuseByAgeConfigurationDescriptor instead')
const ResultReuseByAgeConfiguration$json = {
  '1': 'ResultReuseByAgeConfiguration',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'maxageinminutes',
      '3': 211566877,
      '4': 1,
      '5': 5,
      '10': 'maxageinminutes'
    },
  ],
};

/// Descriptor for `ResultReuseByAgeConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultReuseByAgeConfigurationDescriptor =
    $convert.base64Decode(
        'Ch1SZXN1bHRSZXVzZUJ5QWdlQ29uZmlndXJhdGlvbhIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW'
        '5hYmxlZBIrCg9tYXhhZ2Vpbm1pbnV0ZXMYnYLxZCABKAVSD21heGFnZWlubWludXRlcw==');

@$core.Deprecated('Use resultReuseConfigurationDescriptor instead')
const ResultReuseConfiguration$json = {
  '1': 'ResultReuseConfiguration',
  '2': [
    {
      '1': 'resultreusebyageconfiguration',
      '3': 13481299,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultReuseByAgeConfiguration',
      '10': 'resultreusebyageconfiguration'
    },
  ],
};

/// Descriptor for `ResultReuseConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultReuseConfigurationDescriptor = $convert.base64Decode(
    'ChhSZXN1bHRSZXVzZUNvbmZpZ3VyYXRpb24SbgodcmVzdWx0cmV1c2VieWFnZWNvbmZpZ3VyYX'
    'Rpb24Y0+q2BiABKAsyJS5hdGhlbmEuUmVzdWx0UmV1c2VCeUFnZUNvbmZpZ3VyYXRpb25SHXJl'
    'c3VsdHJldXNlYnlhZ2Vjb25maWd1cmF0aW9u');

@$core.Deprecated('Use resultReuseInformationDescriptor instead')
const ResultReuseInformation$json = {
  '1': 'ResultReuseInformation',
  '2': [
    {
      '1': 'reusedpreviousresult',
      '3': 448068450,
      '4': 1,
      '5': 8,
      '10': 'reusedpreviousresult'
    },
  ],
};

/// Descriptor for `ResultReuseInformation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultReuseInformationDescriptor =
    $convert.base64Decode(
        'ChZSZXN1bHRSZXVzZUluZm9ybWF0aW9uEjYKFHJldXNlZHByZXZpb3VzcmVzdWx0GOL209UBIA'
        'EoCFIUcmV1c2VkcHJldmlvdXNyZXN1bHQ=');

@$core.Deprecated('Use resultSetDescriptor instead')
const ResultSet$json = {
  '1': 'ResultSet',
  '2': [
    {
      '1': 'resultsetmetadata',
      '3': 317769226,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultSetMetadata',
      '10': 'resultsetmetadata'
    },
    {
      '1': 'rows',
      '3': 174428857,
      '4': 3,
      '5': 11,
      '6': '.athena.Row',
      '10': 'rows'
    },
  ],
};

/// Descriptor for `ResultSet`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultSetDescriptor = $convert.base64Decode(
    'CglSZXN1bHRTZXQSSwoRcmVzdWx0c2V0bWV0YWRhdGEYiozDlwEgASgLMhkuYXRoZW5hLlJlc3'
    'VsdFNldE1ldGFkYXRhUhFyZXN1bHRzZXRtZXRhZGF0YRIiCgRyb3dzGLmlllMgAygLMgsuYXRo'
    'ZW5hLlJvd1IEcm93cw==');

@$core.Deprecated('Use resultSetMetadataDescriptor instead')
const ResultSetMetadata$json = {
  '1': 'ResultSetMetadata',
  '2': [
    {
      '1': 'columninfo',
      '3': 364742404,
      '4': 3,
      '5': 11,
      '6': '.athena.ColumnInfo',
      '10': 'columninfo'
    },
  ],
};

/// Descriptor for `ResultSetMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultSetMetadataDescriptor = $convert.base64Decode(
    'ChFSZXN1bHRTZXRNZXRhZGF0YRI2Cgpjb2x1bW5pbmZvGISO9q0BIAMoCzISLmF0aGVuYS5Db2'
    'x1bW5JbmZvUgpjb2x1bW5pbmZv');

@$core.Deprecated('Use rowDescriptor instead')
const Row$json = {
  '1': 'Row',
  '2': [
    {
      '1': 'data',
      '3': 525498822,
      '4': 3,
      '5': 11,
      '6': '.athena.Datum',
      '10': 'data'
    },
  ],
};

/// Descriptor for `Row`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rowDescriptor = $convert.base64Decode(
    'CgNSb3cSJQoEZGF0YRjG88n6ASADKAsyDS5hdGhlbmEuRGF0dW1SBGRhdGE=');

@$core.Deprecated('Use s3LoggingConfigurationDescriptor instead')
const S3LoggingConfiguration$json = {
  '1': 'S3LoggingConfiguration',
  '2': [
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'kmskey', '3': 114561194, '4': 1, '5': 9, '10': 'kmskey'},
    {'1': 'loglocation', '3': 192028619, '4': 1, '5': 9, '10': 'loglocation'},
  ],
};

/// Descriptor for `S3LoggingConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3LoggingConfigurationDescriptor = $convert.base64Decode(
    'ChZTM0xvZ2dpbmdDb25maWd1cmF0aW9uEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVkEh'
    'kKBmttc2tleRiqodA2IAEoCVIGa21za2V5EiMKC2xvZ2xvY2F0aW9uGMu/yFsgASgJUgtsb2ds'
    'b2NhdGlvbg==');

@$core.Deprecated('Use sessionAlreadyExistsExceptionDescriptor instead')
const SessionAlreadyExistsException$json = {
  '1': 'SessionAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SessionAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'Ch1TZXNzaW9uQWxyZWFkeUV4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use sessionConfigurationDescriptor instead')
const SessionConfiguration$json = {
  '1': 'SessionConfiguration',
  '2': [
    {
      '1': 'encryptionconfiguration',
      '3': 225764215,
      '4': 1,
      '5': 11,
      '6': '.athena.EncryptionConfiguration',
      '10': 'encryptionconfiguration'
    },
    {
      '1': 'executionrole',
      '3': 253307658,
      '4': 1,
      '5': 9,
      '10': 'executionrole'
    },
    {
      '1': 'idletimeoutseconds',
      '3': 61279260,
      '4': 1,
      '5': 3,
      '10': 'idletimeoutseconds'
    },
    {
      '1': 'sessionidletimeoutinminutes',
      '3': 515304989,
      '4': 1,
      '5': 5,
      '10': 'sessionidletimeoutinminutes'
    },
    {
      '1': 'workingdirectory',
      '3': 478970252,
      '4': 1,
      '5': 9,
      '10': 'workingdirectory'
    },
  ],
};

/// Descriptor for `SessionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionConfigurationDescriptor = $convert.base64Decode(
    'ChRTZXNzaW9uQ29uZmlndXJhdGlvbhJcChdlbmNyeXB0aW9uY29uZmlndXJhdGlvbhj3xtNrIA'
    'EoCzIfLmF0aGVuYS5FbmNyeXB0aW9uQ29uZmlndXJhdGlvblIXZW5jcnlwdGlvbmNvbmZpZ3Vy'
    'YXRpb24SJwoNZXhlY3V0aW9ucm9sZRiK1uR4IAEoCVINZXhlY3V0aW9ucm9sZRIxChJpZGxldG'
    'ltZW91dHNlY29uZHMYnJicHSABKANSEmlkbGV0aW1lb3V0c2Vjb25kcxJEChtzZXNzaW9uaWRs'
    'ZXRpbWVvdXRpbm1pbnV0ZXMYndzb9QEgASgFUhtzZXNzaW9uaWRsZXRpbWVvdXRpbm1pbnV0ZX'
    'MSLgoQd29ya2luZ2RpcmVjdG9yeRiMg7LkASABKAlSEHdvcmtpbmdkaXJlY3Rvcnk=');

@$core.Deprecated('Use sessionStatisticsDescriptor instead')
const SessionStatistics$json = {
  '1': 'SessionStatistics',
  '2': [
    {
      '1': 'dpuexecutioninmillis',
      '3': 174857936,
      '4': 1,
      '5': 3,
      '10': 'dpuexecutioninmillis'
    },
  ],
};

/// Descriptor for `SessionStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionStatisticsDescriptor = $convert.base64Decode(
    'ChFTZXNzaW9uU3RhdGlzdGljcxI1ChRkcHVleGVjdXRpb25pbm1pbGxpcxjQvbBTIAEoA1IUZH'
    'B1ZXhlY3V0aW9uaW5taWxsaXM=');

@$core.Deprecated('Use sessionStatusDescriptor instead')
const SessionStatus$json = {
  '1': 'SessionStatus',
  '2': [
    {'1': 'enddatetime', '3': 372000488, '4': 1, '5': 9, '10': 'enddatetime'},
    {
      '1': 'idlesincedatetime',
      '3': 454742375,
      '4': 1,
      '5': 9,
      '10': 'idlesincedatetime'
    },
    {
      '1': 'lastmodifieddatetime',
      '3': 406490164,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddatetime'
    },
    {
      '1': 'startdatetime',
      '3': 88518355,
      '4': 1,
      '5': 9,
      '10': 'startdatetime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.SessionState',
      '10': 'state'
    },
    {
      '1': 'statechangereason',
      '3': 228940439,
      '4': 1,
      '5': 9,
      '10': 'statechangereason'
    },
  ],
};

/// Descriptor for `SessionStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionStatusDescriptor = $convert.base64Decode(
    'Cg1TZXNzaW9uU3RhdHVzEiQKC2VuZGRhdGV0aW1lGOiNsbEBIAEoCVILZW5kZGF0ZXRpbWUSMA'
    'oRaWRsZXNpbmNlZGF0ZXRpbWUY56Lr2AEgASgJUhFpZGxlc2luY2VkYXRldGltZRI2ChRsYXN0'
    'bW9kaWZpZWRkYXRldGltZRi0mOrBASABKAlSFGxhc3Rtb2RpZmllZGRhdGV0aW1lEicKDXN0YX'
    'J0ZGF0ZXRpbWUY092aKiABKAlSDXN0YXJ0ZGF0ZXRpbWUSLgoFc3RhdGUYl8my7wEgASgOMhQu'
    'YXRoZW5hLlNlc3Npb25TdGF0ZVIFc3RhdGUSLwoRc3RhdGVjaGFuZ2VyZWFzb24Yl7WVbSABKA'
    'lSEXN0YXRlY2hhbmdlcmVhc29u');

@$core.Deprecated('Use sessionSummaryDescriptor instead')
const SessionSummary$json = {
  '1': 'SessionSummary',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'engineversion',
      '3': 44953462,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineVersion',
      '10': 'engineversion'
    },
    {
      '1': 'notebookversion',
      '3': 528689837,
      '4': 1,
      '5': 9,
      '10': 'notebookversion'
    },
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.athena.SessionStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `SessionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionSummaryDescriptor = $convert.base64Decode(
    'Cg5TZXNzaW9uU3VtbWFyeRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SPg'
    'oNZW5naW5ldmVyc2lvbhj23rcVIAEoCzIVLmF0aGVuYS5FbmdpbmVWZXJzaW9uUg1lbmdpbmV2'
    'ZXJzaW9uEiwKD25vdGVib29rdmVyc2lvbhit1Yz8ASABKAlSD25vdGVib29rdmVyc2lvbhIfCg'
    'lzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb25pZBIwCgZzdGF0dXMYkOT7AiABKAsyFS5hdGhl'
    'bmEuU2Vzc2lvblN0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use startCalculationExecutionRequestDescriptor instead')
const StartCalculationExecutionRequest$json = {
  '1': 'StartCalculationExecutionRequest',
  '2': [
    {
      '1': 'calculationconfiguration',
      '3': 459676797,
      '4': 1,
      '5': 11,
      '6': '.athena.CalculationConfiguration',
      '10': 'calculationconfiguration'
    },
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {'1': 'codeblock', '3': 23945838, '4': 1, '5': 9, '10': 'codeblock'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `StartCalculationExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startCalculationExecutionRequestDescriptor = $convert.base64Decode(
    'CiBTdGFydENhbGN1bGF0aW9uRXhlY3V0aW9uUmVxdWVzdBJgChhjYWxjdWxhdGlvbmNvbmZpZ3'
    'VyYXRpb24Y/biY2wEgASgLMiAuYXRoZW5hLkNhbGN1bGF0aW9uQ29uZmlndXJhdGlvblIYY2Fs'
    'Y3VsYXRpb25jb25maWd1cmF0aW9uEjIKEmNsaWVudHJlcXVlc3R0b2tlbhjx76LZASABKAlSEm'
    'NsaWVudHJlcXVlc3R0b2tlbhIfCgljb2RlYmxvY2sY7sS1CyABKAlSCWNvZGVibG9jaxIjCgtk'
    'ZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SHwoJc2Vzc2lvbmlkGLuE5QkgASgJUg'
    'lzZXNzaW9uaWQ=');

@$core.Deprecated('Use startCalculationExecutionResponseDescriptor instead')
const StartCalculationExecutionResponse$json = {
  '1': 'StartCalculationExecutionResponse',
  '2': [
    {
      '1': 'calculationexecutionid',
      '3': 80028050,
      '4': 1,
      '5': 9,
      '10': 'calculationexecutionid'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.CalculationExecutionState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `StartCalculationExecutionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startCalculationExecutionResponseDescriptor =
    $convert.base64Decode(
        'CiFTdGFydENhbGN1bGF0aW9uRXhlY3V0aW9uUmVzcG9uc2USOQoWY2FsY3VsYXRpb25leGVjdX'
        'Rpb25pZBiSw5QmIAEoCVIWY2FsY3VsYXRpb25leGVjdXRpb25pZBI7CgVzdGF0ZRiXybLvASAB'
        'KA4yIS5hdGhlbmEuQ2FsY3VsYXRpb25FeGVjdXRpb25TdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use startQueryExecutionInputDescriptor instead')
const StartQueryExecutionInput$json = {
  '1': 'StartQueryExecutionInput',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {
      '1': 'engineconfiguration',
      '3': 341629412,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineConfiguration',
      '10': 'engineconfiguration'
    },
    {
      '1': 'executionparameters',
      '3': 527591242,
      '4': 3,
      '5': 9,
      '10': 'executionparameters'
    },
    {
      '1': 'queryexecutioncontext',
      '3': 139302379,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryExecutionContext',
      '10': 'queryexecutioncontext'
    },
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
    {
      '1': 'resultconfiguration',
      '3': 183031503,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultConfiguration',
      '10': 'resultconfiguration'
    },
    {
      '1': 'resultreuseconfiguration',
      '3': 482796543,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultReuseConfiguration',
      '10': 'resultreuseconfiguration'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `StartQueryExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startQueryExecutionInputDescriptor = $convert.base64Decode(
    'ChhTdGFydFF1ZXJ5RXhlY3V0aW9uSW5wdXQSMgoSY2xpZW50cmVxdWVzdHRva2VuGPHvotkBIA'
    'EoCVISY2xpZW50cmVxdWVzdHRva2VuElEKE2VuZ2luZWNvbmZpZ3VyYXRpb24Y5LPzogEgASgL'
    'MhsuYXRoZW5hLkVuZ2luZUNvbmZpZ3VyYXRpb25SE2VuZ2luZWNvbmZpZ3VyYXRpb24SNAoTZX'
    'hlY3V0aW9ucGFyYW1ldGVycxjKzsn7ASADKAlSE2V4ZWN1dGlvbnBhcmFtZXRlcnMSVgoVcXVl'
    'cnlleGVjdXRpb25jb250ZXh0GOurtkIgASgLMh0uYXRoZW5hLlF1ZXJ5RXhlY3V0aW9uQ29udG'
    'V4dFIVcXVlcnlleGVjdXRpb25jb250ZXh0EiQKC3F1ZXJ5c3RyaW5nGOfK788BIAEoCVILcXVl'
    'cnlzdHJpbmcSUAoTcmVzdWx0Y29uZmlndXJhdGlvbhjPraNXIAEoCzIbLmF0aGVuYS5SZXN1bH'
    'RDb25maWd1cmF0aW9uUhNyZXN1bHRjb25maWd1cmF0aW9uEmAKGHJlc3VsdHJldXNlY29uZmln'
    'dXJhdGlvbhj/x5vmASABKAsyIC5hdGhlbmEuUmVzdWx0UmV1c2VDb25maWd1cmF0aW9uUhhyZX'
    'N1bHRyZXVzZWNvbmZpZ3VyYXRpb24SIAoJd29ya2dyb3VwGIStofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use startQueryExecutionOutputDescriptor instead')
const StartQueryExecutionOutput$json = {
  '1': 'StartQueryExecutionOutput',
  '2': [
    {
      '1': 'queryexecutionid',
      '3': 467615503,
      '4': 1,
      '5': 9,
      '10': 'queryexecutionid'
    },
  ],
};

/// Descriptor for `StartQueryExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startQueryExecutionOutputDescriptor =
    $convert.base64Decode(
        'ChlTdGFydFF1ZXJ5RXhlY3V0aW9uT3V0cHV0Ei4KEHF1ZXJ5ZXhlY3V0aW9uaWQYj/783gEgAS'
        'gJUhBxdWVyeWV4ZWN1dGlvbmlk');

@$core.Deprecated('Use startSessionRequestDescriptor instead')
const StartSessionRequest$json = {
  '1': 'StartSessionRequest',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {
      '1': 'copyworkgrouptags',
      '3': 429355506,
      '4': 1,
      '5': 8,
      '10': 'copyworkgrouptags'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'engineconfiguration',
      '3': 341629412,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineConfiguration',
      '10': 'engineconfiguration'
    },
    {
      '1': 'executionrole',
      '3': 253307658,
      '4': 1,
      '5': 9,
      '10': 'executionrole'
    },
    {
      '1': 'monitoringconfiguration',
      '3': 364891928,
      '4': 1,
      '5': 11,
      '6': '.athena.MonitoringConfiguration',
      '10': 'monitoringconfiguration'
    },
    {
      '1': 'notebookversion',
      '3': 528689837,
      '4': 1,
      '5': 9,
      '10': 'notebookversion'
    },
    {
      '1': 'sessionidletimeoutinminutes',
      '3': 515304989,
      '4': 1,
      '5': 5,
      '10': 'sessionidletimeoutinminutes'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.athena.Tag',
      '10': 'tags'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `StartSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startSessionRequestDescriptor = $convert.base64Decode(
    'ChNTdGFydFNlc3Npb25SZXF1ZXN0EjIKEmNsaWVudHJlcXVlc3R0b2tlbhjx76LZASABKAlSEm'
    'NsaWVudHJlcXVlc3R0b2tlbhIwChFjb3B5d29ya2dyb3VwdGFncxjy493MASABKAhSEWNvcHl3'
    'b3JrZ3JvdXB0YWdzEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhJRChNlbm'
    'dpbmVjb25maWd1cmF0aW9uGOSz86IBIAEoCzIbLmF0aGVuYS5FbmdpbmVDb25maWd1cmF0aW9u'
    'UhNlbmdpbmVjb25maWd1cmF0aW9uEicKDWV4ZWN1dGlvbnJvbGUYitbkeCABKAlSDWV4ZWN1dG'
    'lvbnJvbGUSXQoXbW9uaXRvcmluZ2NvbmZpZ3VyYXRpb24YmJ7/rQEgASgLMh8uYXRoZW5hLk1v'
    'bml0b3JpbmdDb25maWd1cmF0aW9uUhdtb25pdG9yaW5nY29uZmlndXJhdGlvbhIsCg9ub3RlYm'
    '9va3ZlcnNpb24YrdWM/AEgASgJUg9ub3RlYm9va3ZlcnNpb24SRAobc2Vzc2lvbmlkbGV0aW1l'
    'b3V0aW5taW51dGVzGJ3c2/UBIAEoBVIbc2Vzc2lvbmlkbGV0aW1lb3V0aW5taW51dGVzEiMKBH'
    'RhZ3MYwcH2tQEgAygLMgsuYXRoZW5hLlRhZ1IEdGFncxIgCgl3b3JrZ3JvdXAYhK2h8QEgASgJ'
    'Ugl3b3JrZ3JvdXA=');

@$core.Deprecated('Use startSessionResponseDescriptor instead')
const StartSessionResponse$json = {
  '1': 'StartSessionResponse',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.SessionState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `StartSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startSessionResponseDescriptor = $convert.base64Decode(
    'ChRTdGFydFNlc3Npb25SZXNwb25zZRIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb25pZB'
    'IuCgVzdGF0ZRiXybLvASABKA4yFC5hdGhlbmEuU2Vzc2lvblN0YXRlUgVzdGF0ZQ==');

@$core.Deprecated('Use stopCalculationExecutionRequestDescriptor instead')
const StopCalculationExecutionRequest$json = {
  '1': 'StopCalculationExecutionRequest',
  '2': [
    {
      '1': 'calculationexecutionid',
      '3': 80028050,
      '4': 1,
      '5': 9,
      '10': 'calculationexecutionid'
    },
  ],
};

/// Descriptor for `StopCalculationExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopCalculationExecutionRequestDescriptor =
    $convert.base64Decode(
        'Ch9TdG9wQ2FsY3VsYXRpb25FeGVjdXRpb25SZXF1ZXN0EjkKFmNhbGN1bGF0aW9uZXhlY3V0aW'
        '9uaWQYksOUJiABKAlSFmNhbGN1bGF0aW9uZXhlY3V0aW9uaWQ=');

@$core.Deprecated('Use stopCalculationExecutionResponseDescriptor instead')
const StopCalculationExecutionResponse$json = {
  '1': 'StopCalculationExecutionResponse',
  '2': [
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.CalculationExecutionState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `StopCalculationExecutionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopCalculationExecutionResponseDescriptor =
    $convert.base64Decode(
        'CiBTdG9wQ2FsY3VsYXRpb25FeGVjdXRpb25SZXNwb25zZRI7CgVzdGF0ZRiXybLvASABKA4yIS'
        '5hdGhlbmEuQ2FsY3VsYXRpb25FeGVjdXRpb25TdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use stopQueryExecutionInputDescriptor instead')
const StopQueryExecutionInput$json = {
  '1': 'StopQueryExecutionInput',
  '2': [
    {
      '1': 'queryexecutionid',
      '3': 467615503,
      '4': 1,
      '5': 9,
      '10': 'queryexecutionid'
    },
  ],
};

/// Descriptor for `StopQueryExecutionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopQueryExecutionInputDescriptor =
    $convert.base64Decode(
        'ChdTdG9wUXVlcnlFeGVjdXRpb25JbnB1dBIuChBxdWVyeWV4ZWN1dGlvbmlkGI/+/N4BIAEoCV'
        'IQcXVlcnlleGVjdXRpb25pZA==');

@$core.Deprecated('Use stopQueryExecutionOutputDescriptor instead')
const StopQueryExecutionOutput$json = {
  '1': 'StopQueryExecutionOutput',
};

/// Descriptor for `StopQueryExecutionOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopQueryExecutionOutputDescriptor =
    $convert.base64Decode('ChhTdG9wUXVlcnlFeGVjdXRpb25PdXRwdXQ=');

@$core.Deprecated('Use tableMetadataDescriptor instead')
const TableMetadata$json = {
  '1': 'TableMetadata',
  '2': [
    {
      '1': 'columns',
      '3': 169177053,
      '4': 3,
      '5': 11,
      '6': '.athena.Column',
      '10': 'columns'
    },
    {'1': 'createtime', '3': 490895933, '4': 1, '5': 9, '10': 'createtime'},
    {
      '1': 'lastaccesstime',
      '3': 516574551,
      '4': 1,
      '5': 9,
      '10': 'lastaccesstime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.athena.TableMetadata.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'partitionkeys',
      '3': 200562986,
      '4': 3,
      '5': 11,
      '6': '.athena.Column',
      '10': 'partitionkeys'
    },
    {'1': 'tabletype', '3': 476171176, '4': 1, '5': 9, '10': 'tabletype'},
  ],
  '3': [TableMetadata_ParametersEntry$json],
};

@$core.Deprecated('Use tableMetadataDescriptor instead')
const TableMetadata_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TableMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tableMetadataDescriptor = $convert.base64Decode(
    'Cg1UYWJsZU1ldGFkYXRhEisKB2NvbHVtbnMY3d/VUCADKAsyDi5hdGhlbmEuQ29sdW1uUgdjb2'
    'x1bW5zEiIKCmNyZWF0ZXRpbWUYvfSJ6gEgASgJUgpjcmVhdGV0aW1lEioKDmxhc3RhY2Nlc3N0'
    'aW1lGNeaqfYBIAEoCVIObGFzdGFjY2Vzc3RpbWUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRJJCg'
    'pwYXJhbWV0ZXJzGPqn/usBIAMoCzIlLmF0aGVuYS5UYWJsZU1ldGFkYXRhLlBhcmFtZXRlcnNF'
    'bnRyeVIKcGFyYW1ldGVycxI3Cg1wYXJ0aXRpb25rZXlzGKqy0V8gAygLMg4uYXRoZW5hLkNvbH'
    'VtblINcGFydGl0aW9ua2V5cxIgCgl0YWJsZXR5cGUYqJeH4wEgASgJUgl0YWJsZXR5cGUaPQoP'
    'UGFyYW1ldGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZT'
    'oCOAE=');

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
      '6': '.athena.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceInputDescriptor = $convert.base64Decode(
    'ChBUYWdSZXNvdXJjZUlucHV0EiQKC3Jlc291cmNlYXJuGO3AmbABIAEoCVILcmVzb3VyY2Vhcm'
    '4SIwoEdGFncxjBwfa1ASADKAsyCy5hdGhlbmEuVGFnUgR0YWdz');

@$core.Deprecated('Use tagResourceOutputDescriptor instead')
const TagResourceOutput$json = {
  '1': 'TagResourceOutput',
};

/// Descriptor for `TagResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceOutputDescriptor =
    $convert.base64Decode('ChFUYWdSZXNvdXJjZU91dHB1dA==');

@$core.Deprecated('Use terminateSessionRequestDescriptor instead')
const TerminateSessionRequest$json = {
  '1': 'TerminateSessionRequest',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `TerminateSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List terminateSessionRequestDescriptor =
    $convert.base64Decode(
        'ChdUZXJtaW5hdGVTZXNzaW9uUmVxdWVzdBIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb2'
        '5pZA==');

@$core.Deprecated('Use terminateSessionResponseDescriptor instead')
const TerminateSessionResponse$json = {
  '1': 'TerminateSessionResponse',
  '2': [
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.SessionState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `TerminateSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List terminateSessionResponseDescriptor =
    $convert.base64Decode(
        'ChhUZXJtaW5hdGVTZXNzaW9uUmVzcG9uc2USLgoFc3RhdGUYl8my7wEgASgOMhQuYXRoZW5hLl'
        'Nlc3Npb25TdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use tooManyRequestsExceptionDescriptor instead')
const TooManyRequestsException$json = {
  '1': 'TooManyRequestsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'reason',
      '3': 20005178,
      '4': 1,
      '5': 14,
      '6': '.athena.ThrottleReason',
      '10': 'reason'
    },
  ],
};

/// Descriptor for `TooManyRequestsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyRequestsExceptionDescriptor =
    $convert.base64Decode(
        'ChhUb29NYW55UmVxdWVzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZR'
        'IxCgZyZWFzb24YuoLFCSABKA4yFi5hdGhlbmEuVGhyb3R0bGVSZWFzb25SBnJlYXNvbg==');

@$core.Deprecated('Use unprocessedNamedQueryIdDescriptor instead')
const UnprocessedNamedQueryId$json = {
  '1': 'UnprocessedNamedQueryId',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {'1': 'namedqueryid', '3': 330896872, '4': 1, '5': 9, '10': 'namedqueryid'},
  ],
};

/// Descriptor for `UnprocessedNamedQueryId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unprocessedNamedQueryIdDescriptor = $convert.base64Decode(
    'ChdVbnByb2Nlc3NlZE5hbWVkUXVlcnlJZBIfCgllcnJvcmNvZGUYmdbDECABKAlSCWVycm9yY2'
    '9kZRImCgxlcnJvcm1lc3NhZ2UYqYqr9wEgASgJUgxlcnJvcm1lc3NhZ2USJgoMbmFtZWRxdWVy'
    'eWlkGOir5J0BIAEoCVIMbmFtZWRxdWVyeWlk');

@$core.Deprecated('Use unprocessedPreparedStatementNameDescriptor instead')
const UnprocessedPreparedStatementName$json = {
  '1': 'UnprocessedPreparedStatementName',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'statementname',
      '3': 23047926,
      '4': 1,
      '5': 9,
      '10': 'statementname'
    },
  ],
};

/// Descriptor for `UnprocessedPreparedStatementName`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unprocessedPreparedStatementNameDescriptor =
    $convert.base64Decode(
        'CiBVbnByb2Nlc3NlZFByZXBhcmVkU3RhdGVtZW50TmFtZRIfCgllcnJvcmNvZGUYmdbDECABKA'
        'lSCWVycm9yY29kZRImCgxlcnJvcm1lc3NhZ2UYqYqr9wEgASgJUgxlcnJvcm1lc3NhZ2USJwoN'
        'c3RhdGVtZW50bmFtZRj23f4KIAEoCVINc3RhdGVtZW50bmFtZQ==');

@$core.Deprecated('Use unprocessedQueryExecutionIdDescriptor instead')
const UnprocessedQueryExecutionId$json = {
  '1': 'UnprocessedQueryExecutionId',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'queryexecutionid',
      '3': 467615503,
      '4': 1,
      '5': 9,
      '10': 'queryexecutionid'
    },
  ],
};

/// Descriptor for `UnprocessedQueryExecutionId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unprocessedQueryExecutionIdDescriptor =
    $convert.base64Decode(
        'ChtVbnByb2Nlc3NlZFF1ZXJ5RXhlY3V0aW9uSWQSHwoJZXJyb3Jjb2RlGJnWwxAgASgJUgllcn'
        'JvcmNvZGUSJgoMZXJyb3JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3JtZXNzYWdlEi4KEHF1ZXJ5'
        'ZXhlY3V0aW9uaWQYj/783gEgASgJUhBxdWVyeWV4ZWN1dGlvbmlk');

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

@$core.Deprecated('Use updateCapacityReservationInputDescriptor instead')
const UpdateCapacityReservationInput$json = {
  '1': 'UpdateCapacityReservationInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'targetdpus', '3': 367520745, '4': 1, '5': 5, '10': 'targetdpus'},
  ],
};

/// Descriptor for `UpdateCapacityReservationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCapacityReservationInputDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVDYXBhY2l0eVJlc2VydmF0aW9uSW5wdXQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZR'
        'IiCgp0YXJnZXRkcHVzGOnXn68BIAEoBVIKdGFyZ2V0ZHB1cw==');

@$core.Deprecated('Use updateCapacityReservationOutputDescriptor instead')
const UpdateCapacityReservationOutput$json = {
  '1': 'UpdateCapacityReservationOutput',
};

/// Descriptor for `UpdateCapacityReservationOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCapacityReservationOutputDescriptor =
    $convert.base64Decode('Ch9VcGRhdGVDYXBhY2l0eVJlc2VydmF0aW9uT3V0cHV0');

@$core.Deprecated('Use updateDataCatalogInputDescriptor instead')
const UpdateDataCatalogInput$json = {
  '1': 'UpdateDataCatalogInput',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.athena.UpdateDataCatalogInput.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.athena.DataCatalogType',
      '10': 'type'
    },
  ],
  '3': [UpdateDataCatalogInput_ParametersEntry$json],
};

@$core.Deprecated('Use updateDataCatalogInputDescriptor instead')
const UpdateDataCatalogInput_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateDataCatalogInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDataCatalogInputDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVEYXRhQ2F0YWxvZ0lucHV0EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcm'
    'lwdGlvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lElIKCnBhcmFtZXRlcnMY+qf+6wEgAygLMi4u'
    'YXRoZW5hLlVwZGF0ZURhdGFDYXRhbG9nSW5wdXQuUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZX'
    'JzEi8KBHR5cGUY7qDXigEgASgOMhcuYXRoZW5hLkRhdGFDYXRhbG9nVHlwZVIEdHlwZRo9Cg9Q'
    'YXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOg'
    'I4AQ==');

@$core.Deprecated('Use updateDataCatalogOutputDescriptor instead')
const UpdateDataCatalogOutput$json = {
  '1': 'UpdateDataCatalogOutput',
};

/// Descriptor for `UpdateDataCatalogOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDataCatalogOutputDescriptor =
    $convert.base64Decode('ChdVcGRhdGVEYXRhQ2F0YWxvZ091dHB1dA==');

@$core.Deprecated('Use updateNamedQueryInputDescriptor instead')
const UpdateNamedQueryInput$json = {
  '1': 'UpdateNamedQueryInput',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'namedqueryid', '3': 330896872, '4': 1, '5': 9, '10': 'namedqueryid'},
    {'1': 'querystring', '3': 435938663, '4': 1, '5': 9, '10': 'querystring'},
  ],
};

/// Descriptor for `UpdateNamedQueryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateNamedQueryInputDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVOYW1lZFF1ZXJ5SW5wdXQSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaX'
    'B0aW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSJgoMbmFtZWRxdWVyeWlkGOir5J0BIAEoCVIM'
    'bmFtZWRxdWVyeWlkEiQKC3F1ZXJ5c3RyaW5nGOfK788BIAEoCVILcXVlcnlzdHJpbmc=');

@$core.Deprecated('Use updateNamedQueryOutputDescriptor instead')
const UpdateNamedQueryOutput$json = {
  '1': 'UpdateNamedQueryOutput',
};

/// Descriptor for `UpdateNamedQueryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateNamedQueryOutputDescriptor =
    $convert.base64Decode('ChZVcGRhdGVOYW1lZFF1ZXJ5T3V0cHV0');

@$core.Deprecated('Use updateNotebookInputDescriptor instead')
const UpdateNotebookInput$json = {
  '1': 'UpdateNotebookInput',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
    {'1': 'payload', '3': 6526790, '4': 1, '5': 9, '10': 'payload'},
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.athena.NotebookType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `UpdateNotebookInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateNotebookInputDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVOb3RlYm9va0lucHV0EjIKEmNsaWVudHJlcXVlc3R0b2tlbhjx76LZASABKAlSEm'
    'NsaWVudHJlcXVlc3R0b2tlbhIhCgpub3RlYm9va2lkGN60lUsgASgJUgpub3RlYm9va2lkEhsK'
    'B3BheWxvYWQYxq6OAyABKAlSB3BheWxvYWQSHwoJc2Vzc2lvbmlkGLuE5QkgASgJUglzZXNzaW'
    '9uaWQSLAoEdHlwZRjuoNeKASABKA4yFC5hdGhlbmEuTm90ZWJvb2tUeXBlUgR0eXBl');

@$core.Deprecated('Use updateNotebookMetadataInputDescriptor instead')
const UpdateNotebookMetadataInput$json = {
  '1': 'UpdateNotebookMetadataInput',
  '2': [
    {
      '1': 'clientrequesttoken',
      '3': 455653361,
      '4': 1,
      '5': 9,
      '10': 'clientrequesttoken'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'notebookid', '3': 157637214, '4': 1, '5': 9, '10': 'notebookid'},
  ],
};

/// Descriptor for `UpdateNotebookMetadataInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateNotebookMetadataInputDescriptor =
    $convert.base64Decode(
        'ChtVcGRhdGVOb3RlYm9va01ldGFkYXRhSW5wdXQSMgoSY2xpZW50cmVxdWVzdHRva2VuGPHvot'
        'kBIAEoCVISY2xpZW50cmVxdWVzdHRva2VuEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSIQoKbm90'
        'ZWJvb2tpZBjetJVLIAEoCVIKbm90ZWJvb2tpZA==');

@$core.Deprecated('Use updateNotebookMetadataOutputDescriptor instead')
const UpdateNotebookMetadataOutput$json = {
  '1': 'UpdateNotebookMetadataOutput',
};

/// Descriptor for `UpdateNotebookMetadataOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateNotebookMetadataOutputDescriptor =
    $convert.base64Decode('ChxVcGRhdGVOb3RlYm9va01ldGFkYXRhT3V0cHV0');

@$core.Deprecated('Use updateNotebookOutputDescriptor instead')
const UpdateNotebookOutput$json = {
  '1': 'UpdateNotebookOutput',
};

/// Descriptor for `UpdateNotebookOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateNotebookOutputDescriptor =
    $convert.base64Decode('ChRVcGRhdGVOb3RlYm9va091dHB1dA==');

@$core.Deprecated('Use updatePreparedStatementInputDescriptor instead')
const UpdatePreparedStatementInput$json = {
  '1': 'UpdatePreparedStatementInput',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'querystatement',
      '3': 340852217,
      '4': 1,
      '5': 9,
      '10': 'querystatement'
    },
    {
      '1': 'statementname',
      '3': 23047926,
      '4': 1,
      '5': 9,
      '10': 'statementname'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `UpdatePreparedStatementInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePreparedStatementInputDescriptor = $convert.base64Decode(
    'ChxVcGRhdGVQcmVwYXJlZFN0YXRlbWVudElucHV0EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUg'
    'tkZXNjcmlwdGlvbhIqCg5xdWVyeXN0YXRlbWVudBj5+8OiASABKAlSDnF1ZXJ5c3RhdGVtZW50'
    'EicKDXN0YXRlbWVudG5hbWUY9t3+CiABKAlSDXN0YXRlbWVudG5hbWUSIAoJd29ya2dyb3VwGI'
    'StofEBIAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use updatePreparedStatementOutputDescriptor instead')
const UpdatePreparedStatementOutput$json = {
  '1': 'UpdatePreparedStatementOutput',
};

/// Descriptor for `UpdatePreparedStatementOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePreparedStatementOutputDescriptor =
    $convert.base64Decode('Ch1VcGRhdGVQcmVwYXJlZFN0YXRlbWVudE91dHB1dA==');

@$core.Deprecated('Use updateWorkGroupInputDescriptor instead')
const UpdateWorkGroupInput$json = {
  '1': 'UpdateWorkGroupInput',
  '2': [
    {
      '1': 'configurationupdates',
      '3': 133706738,
      '4': 1,
      '5': 11,
      '6': '.athena.WorkGroupConfigurationUpdates',
      '10': 'configurationupdates'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.WorkGroupState',
      '10': 'state'
    },
    {'1': 'workgroup', '3': 505960068, '4': 1, '5': 9, '10': 'workgroup'},
  ],
};

/// Descriptor for `UpdateWorkGroupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateWorkGroupInputDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVXb3JrR3JvdXBJbnB1dBJcChRjb25maWd1cmF0aW9udXBkYXRlcxjy5+A/IAEoCz'
    'IlLmF0aGVuYS5Xb3JrR3JvdXBDb25maWd1cmF0aW9uVXBkYXRlc1IUY29uZmlndXJhdGlvbnVw'
    'ZGF0ZXMSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEjAKBXN0YXRlGJfJsu'
    '8BIAEoDjIWLmF0aGVuYS5Xb3JrR3JvdXBTdGF0ZVIFc3RhdGUSIAoJd29ya2dyb3VwGIStofEB'
    'IAEoCVIJd29ya2dyb3Vw');

@$core.Deprecated('Use updateWorkGroupOutputDescriptor instead')
const UpdateWorkGroupOutput$json = {
  '1': 'UpdateWorkGroupOutput',
};

/// Descriptor for `UpdateWorkGroupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateWorkGroupOutputDescriptor =
    $convert.base64Decode('ChVVcGRhdGVXb3JrR3JvdXBPdXRwdXQ=');

@$core.Deprecated('Use workGroupDescriptor instead')
const WorkGroup$json = {
  '1': 'WorkGroup',
  '2': [
    {
      '1': 'configuration',
      '3': 442426458,
      '4': 1,
      '5': 11,
      '6': '.athena.WorkGroupConfiguration',
      '10': 'configuration'
    },
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'identitycenterapplicationarn',
      '3': 338293532,
      '4': 1,
      '5': 9,
      '10': 'identitycenterapplicationarn'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.WorkGroupState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `WorkGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List workGroupDescriptor = $convert.base64Decode(
    'CglXb3JrR3JvdXASSAoNY29uZmlndXJhdGlvbhjayPvSASABKAsyHi5hdGhlbmEuV29ya0dyb3'
    'VwQ29uZmlndXJhdGlvblINY29uZmlndXJhdGlvbhIlCgxjcmVhdGlvbnRpbWUY5s+qMSABKAlS'
    'DGNyZWF0aW9udGltZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SRgocaW'
    'RlbnRpdHljZW50ZXJhcHBsaWNhdGlvbmFybhic5qehASABKAlSHGlkZW50aXR5Y2VudGVyYXBw'
    'bGljYXRpb25hcm4SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIwCgVzdGF0ZRiXybLvASABKA4yFi'
    '5hdGhlbmEuV29ya0dyb3VwU3RhdGVSBXN0YXRl');

@$core.Deprecated('Use workGroupConfigurationDescriptor instead')
const WorkGroupConfiguration$json = {
  '1': 'WorkGroupConfiguration',
  '2': [
    {
      '1': 'additionalconfiguration',
      '3': 389584375,
      '4': 1,
      '5': 9,
      '10': 'additionalconfiguration'
    },
    {
      '1': 'bytesscannedcutoffperquery',
      '3': 265761289,
      '4': 1,
      '5': 3,
      '10': 'bytesscannedcutoffperquery'
    },
    {
      '1': 'customercontentencryptionconfiguration',
      '3': 165213900,
      '4': 1,
      '5': 11,
      '6': '.athena.CustomerContentEncryptionConfiguration',
      '10': 'customercontentencryptionconfiguration'
    },
    {
      '1': 'enableminimumencryptionconfiguration',
      '3': 238637616,
      '4': 1,
      '5': 8,
      '10': 'enableminimumencryptionconfiguration'
    },
    {
      '1': 'enforceworkgroupconfiguration',
      '3': 152602624,
      '4': 1,
      '5': 8,
      '10': 'enforceworkgroupconfiguration'
    },
    {
      '1': 'engineconfiguration',
      '3': 341629412,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineConfiguration',
      '10': 'engineconfiguration'
    },
    {
      '1': 'engineversion',
      '3': 44953462,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineVersion',
      '10': 'engineversion'
    },
    {
      '1': 'executionrole',
      '3': 253307658,
      '4': 1,
      '5': 9,
      '10': 'executionrole'
    },
    {
      '1': 'identitycenterconfiguration',
      '3': 236974917,
      '4': 1,
      '5': 11,
      '6': '.athena.IdentityCenterConfiguration',
      '10': 'identitycenterconfiguration'
    },
    {
      '1': 'managedqueryresultsconfiguration',
      '3': 159215683,
      '4': 1,
      '5': 11,
      '6': '.athena.ManagedQueryResultsConfiguration',
      '10': 'managedqueryresultsconfiguration'
    },
    {
      '1': 'monitoringconfiguration',
      '3': 364891928,
      '4': 1,
      '5': 11,
      '6': '.athena.MonitoringConfiguration',
      '10': 'monitoringconfiguration'
    },
    {
      '1': 'publishcloudwatchmetricsenabled',
      '3': 493112579,
      '4': 1,
      '5': 8,
      '10': 'publishcloudwatchmetricsenabled'
    },
    {
      '1': 'queryresultss3accessgrantsconfiguration',
      '3': 264403639,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryResultsS3AccessGrantsConfiguration',
      '10': 'queryresultss3accessgrantsconfiguration'
    },
    {
      '1': 'requesterpaysenabled',
      '3': 448832780,
      '4': 1,
      '5': 8,
      '10': 'requesterpaysenabled'
    },
    {
      '1': 'resultconfiguration',
      '3': 183031503,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultConfiguration',
      '10': 'resultconfiguration'
    },
  ],
};

/// Descriptor for `WorkGroupConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List workGroupConfigurationDescriptor = $convert.base64Decode(
    'ChZXb3JrR3JvdXBDb25maWd1cmF0aW9uEjwKF2FkZGl0aW9uYWxjb25maWd1cmF0aW9uGPer4r'
    'kBIAEoCVIXYWRkaXRpb25hbGNvbmZpZ3VyYXRpb24SQQoaYnl0ZXNzY2FubmVkY3V0b2ZmcGVy'
    'cXVlcnkYieTcfiABKANSGmJ5dGVzc2Nhbm5lZGN1dG9mZnBlcnF1ZXJ5EokBCiZjdXN0b21lcm'
    'NvbnRlbnRlbmNyeXB0aW9uY29uZmlndXJhdGlvbhjM7eNOIAEoCzIuLmF0aGVuYS5DdXN0b21l'
    'ckNvbnRlbnRFbmNyeXB0aW9uQ29uZmlndXJhdGlvblImY3VzdG9tZXJjb250ZW50ZW5jcnlwdG'
    'lvbmNvbmZpZ3VyYXRpb24SVQokZW5hYmxlbWluaW11bWVuY3J5cHRpb25jb25maWd1cmF0aW9u'
    'GLCk5XEgASgIUiRlbmFibGVtaW5pbXVtZW5jcnlwdGlvbmNvbmZpZ3VyYXRpb24SRwodZW5mb3'
    'JjZXdvcmtncm91cGNvbmZpZ3VyYXRpb24YgJDiSCABKAhSHWVuZm9yY2V3b3JrZ3JvdXBjb25m'
    'aWd1cmF0aW9uElEKE2VuZ2luZWNvbmZpZ3VyYXRpb24Y5LPzogEgASgLMhsuYXRoZW5hLkVuZ2'
    'luZUNvbmZpZ3VyYXRpb25SE2VuZ2luZWNvbmZpZ3VyYXRpb24SPgoNZW5naW5ldmVyc2lvbhj2'
    '3rcVIAEoCzIVLmF0aGVuYS5FbmdpbmVWZXJzaW9uUg1lbmdpbmV2ZXJzaW9uEicKDWV4ZWN1dG'
    'lvbnJvbGUYitbkeCABKAlSDWV4ZWN1dGlvbnJvbGUSaAobaWRlbnRpdHljZW50ZXJjb25maWd1'
    'cmF0aW9uGMXm/3AgASgLMiMuYXRoZW5hLklkZW50aXR5Q2VudGVyQ29uZmlndXJhdGlvblIbaW'
    'RlbnRpdHljZW50ZXJjb25maWd1cmF0aW9uEncKIG1hbmFnZWRxdWVyeXJlc3VsdHNjb25maWd1'
    'cmF0aW9uGMPg9UsgASgLMiguYXRoZW5hLk1hbmFnZWRRdWVyeVJlc3VsdHNDb25maWd1cmF0aW'
    '9uUiBtYW5hZ2VkcXVlcnlyZXN1bHRzY29uZmlndXJhdGlvbhJdChdtb25pdG9yaW5nY29uZmln'
    'dXJhdGlvbhiYnv+tASABKAsyHy5hdGhlbmEuTW9uaXRvcmluZ0NvbmZpZ3VyYXRpb25SF21vbm'
    'l0b3Jpbmdjb25maWd1cmF0aW9uEkwKH3B1Ymxpc2hjbG91ZHdhdGNobWV0cmljc2VuYWJsZWQY'
    'g5qR6wEgASgIUh9wdWJsaXNoY2xvdWR3YXRjaG1ldHJpY3NlbmFibGVkEowBCidxdWVyeXJlc3'
    'VsdHNzM2FjY2Vzc2dyYW50c2NvbmZpZ3VyYXRpb24Yt/WJfiABKAsyLy5hdGhlbmEuUXVlcnlS'
    'ZXN1bHRzUzNBY2Nlc3NHcmFudHNDb25maWd1cmF0aW9uUidxdWVyeXJlc3VsdHNzM2FjY2Vzc2'
    'dyYW50c2NvbmZpZ3VyYXRpb24SNgoUcmVxdWVzdGVycGF5c2VuYWJsZWQYjMqC1gEgASgIUhRy'
    'ZXF1ZXN0ZXJwYXlzZW5hYmxlZBJQChNyZXN1bHRjb25maWd1cmF0aW9uGM+to1cgASgLMhsuYX'
    'RoZW5hLlJlc3VsdENvbmZpZ3VyYXRpb25SE3Jlc3VsdGNvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use workGroupConfigurationUpdatesDescriptor instead')
const WorkGroupConfigurationUpdates$json = {
  '1': 'WorkGroupConfigurationUpdates',
  '2': [
    {
      '1': 'additionalconfiguration',
      '3': 389584375,
      '4': 1,
      '5': 9,
      '10': 'additionalconfiguration'
    },
    {
      '1': 'bytesscannedcutoffperquery',
      '3': 265761289,
      '4': 1,
      '5': 3,
      '10': 'bytesscannedcutoffperquery'
    },
    {
      '1': 'customercontentencryptionconfiguration',
      '3': 165213900,
      '4': 1,
      '5': 11,
      '6': '.athena.CustomerContentEncryptionConfiguration',
      '10': 'customercontentencryptionconfiguration'
    },
    {
      '1': 'enableminimumencryptionconfiguration',
      '3': 238637616,
      '4': 1,
      '5': 8,
      '10': 'enableminimumencryptionconfiguration'
    },
    {
      '1': 'enforceworkgroupconfiguration',
      '3': 152602624,
      '4': 1,
      '5': 8,
      '10': 'enforceworkgroupconfiguration'
    },
    {
      '1': 'engineconfiguration',
      '3': 341629412,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineConfiguration',
      '10': 'engineconfiguration'
    },
    {
      '1': 'engineversion',
      '3': 44953462,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineVersion',
      '10': 'engineversion'
    },
    {
      '1': 'executionrole',
      '3': 253307658,
      '4': 1,
      '5': 9,
      '10': 'executionrole'
    },
    {
      '1': 'managedqueryresultsconfigurationupdates',
      '3': 211421145,
      '4': 1,
      '5': 11,
      '6': '.athena.ManagedQueryResultsConfigurationUpdates',
      '10': 'managedqueryresultsconfigurationupdates'
    },
    {
      '1': 'monitoringconfiguration',
      '3': 364891928,
      '4': 1,
      '5': 11,
      '6': '.athena.MonitoringConfiguration',
      '10': 'monitoringconfiguration'
    },
    {
      '1': 'publishcloudwatchmetricsenabled',
      '3': 493112579,
      '4': 1,
      '5': 8,
      '10': 'publishcloudwatchmetricsenabled'
    },
    {
      '1': 'queryresultss3accessgrantsconfiguration',
      '3': 264403639,
      '4': 1,
      '5': 11,
      '6': '.athena.QueryResultsS3AccessGrantsConfiguration',
      '10': 'queryresultss3accessgrantsconfiguration'
    },
    {
      '1': 'removebytesscannedcutoffperquery',
      '3': 130456497,
      '4': 1,
      '5': 8,
      '10': 'removebytesscannedcutoffperquery'
    },
    {
      '1': 'removecustomercontentencryptionconfiguration',
      '3': 418402884,
      '4': 1,
      '5': 8,
      '10': 'removecustomercontentencryptionconfiguration'
    },
    {
      '1': 'requesterpaysenabled',
      '3': 448832780,
      '4': 1,
      '5': 8,
      '10': 'requesterpaysenabled'
    },
    {
      '1': 'resultconfigurationupdates',
      '3': 97856045,
      '4': 1,
      '5': 11,
      '6': '.athena.ResultConfigurationUpdates',
      '10': 'resultconfigurationupdates'
    },
  ],
};

/// Descriptor for `WorkGroupConfigurationUpdates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List workGroupConfigurationUpdatesDescriptor = $convert.base64Decode(
    'Ch1Xb3JrR3JvdXBDb25maWd1cmF0aW9uVXBkYXRlcxI8ChdhZGRpdGlvbmFsY29uZmlndXJhdG'
    'lvbhj3q+K5ASABKAlSF2FkZGl0aW9uYWxjb25maWd1cmF0aW9uEkEKGmJ5dGVzc2Nhbm5lZGN1'
    'dG9mZnBlcnF1ZXJ5GInk3H4gASgDUhpieXRlc3NjYW5uZWRjdXRvZmZwZXJxdWVyeRKJAQomY3'
    'VzdG9tZXJjb250ZW50ZW5jcnlwdGlvbmNvbmZpZ3VyYXRpb24YzO3jTiABKAsyLi5hdGhlbmEu'
    'Q3VzdG9tZXJDb250ZW50RW5jcnlwdGlvbkNvbmZpZ3VyYXRpb25SJmN1c3RvbWVyY29udGVudG'
    'VuY3J5cHRpb25jb25maWd1cmF0aW9uElUKJGVuYWJsZW1pbmltdW1lbmNyeXB0aW9uY29uZmln'
    'dXJhdGlvbhiwpOVxIAEoCFIkZW5hYmxlbWluaW11bWVuY3J5cHRpb25jb25maWd1cmF0aW9uEk'
    'cKHWVuZm9yY2V3b3JrZ3JvdXBjb25maWd1cmF0aW9uGICQ4kggASgIUh1lbmZvcmNld29ya2dy'
    'b3VwY29uZmlndXJhdGlvbhJRChNlbmdpbmVjb25maWd1cmF0aW9uGOSz86IBIAEoCzIbLmF0aG'
    'VuYS5FbmdpbmVDb25maWd1cmF0aW9uUhNlbmdpbmVjb25maWd1cmF0aW9uEj4KDWVuZ2luZXZl'
    'cnNpb24Y9t63FSABKAsyFS5hdGhlbmEuRW5naW5lVmVyc2lvblINZW5naW5ldmVyc2lvbhInCg'
    '1leGVjdXRpb25yb2xlGIrW5HggASgJUg1leGVjdXRpb25yb2xlEowBCidtYW5hZ2VkcXVlcnly'
    'ZXN1bHRzY29uZmlndXJhdGlvbnVwZGF0ZXMY2Y/oZCABKAsyLy5hdGhlbmEuTWFuYWdlZFF1ZX'
    'J5UmVzdWx0c0NvbmZpZ3VyYXRpb25VcGRhdGVzUidtYW5hZ2VkcXVlcnlyZXN1bHRzY29uZmln'
    'dXJhdGlvbnVwZGF0ZXMSXQoXbW9uaXRvcmluZ2NvbmZpZ3VyYXRpb24YmJ7/rQEgASgLMh8uYX'
    'RoZW5hLk1vbml0b3JpbmdDb25maWd1cmF0aW9uUhdtb25pdG9yaW5nY29uZmlndXJhdGlvbhJM'
    'Ch9wdWJsaXNoY2xvdWR3YXRjaG1ldHJpY3NlbmFibGVkGIOakesBIAEoCFIfcHVibGlzaGNsb3'
    'Vkd2F0Y2htZXRyaWNzZW5hYmxlZBKMAQoncXVlcnlyZXN1bHRzczNhY2Nlc3NncmFudHNjb25m'
    'aWd1cmF0aW9uGLf1iX4gASgLMi8uYXRoZW5hLlF1ZXJ5UmVzdWx0c1MzQWNjZXNzR3JhbnRzQ2'
    '9uZmlndXJhdGlvblIncXVlcnlyZXN1bHRzczNhY2Nlc3NncmFudHNjb25maWd1cmF0aW9uEk0K'
    'IHJlbW92ZWJ5dGVzc2Nhbm5lZGN1dG9mZnBlcnF1ZXJ5GLG3mj4gASgIUiByZW1vdmVieXRlc3'
    'NjYW5uZWRjdXRvZmZwZXJxdWVyeRJmCixyZW1vdmVjdXN0b21lcmNvbnRlbnRlbmNyeXB0aW9u'
    'Y29uZmlndXJhdGlvbhjEpMHHASABKAhSLHJlbW92ZWN1c3RvbWVyY29udGVudGVuY3J5cHRpb2'
    '5jb25maWd1cmF0aW9uEjYKFHJlcXVlc3RlcnBheXNlbmFibGVkGIzKgtYBIAEoCFIUcmVxdWVz'
    'dGVycGF5c2VuYWJsZWQSZQoacmVzdWx0Y29uZmlndXJhdGlvbnVwZGF0ZXMYrdTULiABKAsyIi'
    '5hdGhlbmEuUmVzdWx0Q29uZmlndXJhdGlvblVwZGF0ZXNSGnJlc3VsdGNvbmZpZ3VyYXRpb251'
    'cGRhdGVz');

@$core.Deprecated('Use workGroupSummaryDescriptor instead')
const WorkGroupSummary$json = {
  '1': 'WorkGroupSummary',
  '2': [
    {'1': 'creationtime', '3': 103458790, '4': 1, '5': 9, '10': 'creationtime'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'engineversion',
      '3': 44953462,
      '4': 1,
      '5': 11,
      '6': '.athena.EngineVersion',
      '10': 'engineversion'
    },
    {
      '1': 'identitycenterapplicationarn',
      '3': 338293532,
      '4': 1,
      '5': 9,
      '10': 'identitycenterapplicationarn'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.athena.WorkGroupState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `WorkGroupSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List workGroupSummaryDescriptor = $convert.base64Decode(
    'ChBXb3JrR3JvdXBTdW1tYXJ5EiUKDGNyZWF0aW9udGltZRjmz6oxIAEoCVIMY3JlYXRpb250aW'
    '1lEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhI+Cg1lbmdpbmV2ZXJzaW9u'
    'GPbetxUgASgLMhUuYXRoZW5hLkVuZ2luZVZlcnNpb25SDWVuZ2luZXZlcnNpb24SRgocaWRlbn'
    'RpdHljZW50ZXJhcHBsaWNhdGlvbmFybhic5qehASABKAlSHGlkZW50aXR5Y2VudGVyYXBwbGlj'
    'YXRpb25hcm4SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIwCgVzdGF0ZRiXybLvASABKA4yFi5hdG'
    'hlbmEuV29ya0dyb3VwU3RhdGVSBXN0YXRl');
