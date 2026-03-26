// This is a generated file - do not edit.
//
// Generated from athena.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AuthenticationType extends $pb.ProtobufEnum {
  static const AuthenticationType AUTHENTICATION_TYPE_DIRECTORY_IDENTITY =
      AuthenticationType._(
          0, _omitEnumNames ? '' : 'AUTHENTICATION_TYPE_DIRECTORY_IDENTITY');

  static const $core.List<AuthenticationType> values = <AuthenticationType>[
    AUTHENTICATION_TYPE_DIRECTORY_IDENTITY,
  ];

  static final $core.List<AuthenticationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static AuthenticationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AuthenticationType._(super.value, super.name);
}

class CalculationExecutionState extends $pb.ProtobufEnum {
  static const CalculationExecutionState CALCULATION_EXECUTION_STATE_QUEUED =
      CalculationExecutionState._(
          0, _omitEnumNames ? '' : 'CALCULATION_EXECUTION_STATE_QUEUED');
  static const CalculationExecutionState CALCULATION_EXECUTION_STATE_RUNNING =
      CalculationExecutionState._(
          1, _omitEnumNames ? '' : 'CALCULATION_EXECUTION_STATE_RUNNING');
  static const CalculationExecutionState CALCULATION_EXECUTION_STATE_CANCELED =
      CalculationExecutionState._(
          2, _omitEnumNames ? '' : 'CALCULATION_EXECUTION_STATE_CANCELED');
  static const CalculationExecutionState CALCULATION_EXECUTION_STATE_CANCELING =
      CalculationExecutionState._(
          3, _omitEnumNames ? '' : 'CALCULATION_EXECUTION_STATE_CANCELING');
  static const CalculationExecutionState CALCULATION_EXECUTION_STATE_CREATING =
      CalculationExecutionState._(
          4, _omitEnumNames ? '' : 'CALCULATION_EXECUTION_STATE_CREATING');
  static const CalculationExecutionState CALCULATION_EXECUTION_STATE_COMPLETED =
      CalculationExecutionState._(
          5, _omitEnumNames ? '' : 'CALCULATION_EXECUTION_STATE_COMPLETED');
  static const CalculationExecutionState CALCULATION_EXECUTION_STATE_CREATED =
      CalculationExecutionState._(
          6, _omitEnumNames ? '' : 'CALCULATION_EXECUTION_STATE_CREATED');
  static const CalculationExecutionState CALCULATION_EXECUTION_STATE_FAILED =
      CalculationExecutionState._(
          7, _omitEnumNames ? '' : 'CALCULATION_EXECUTION_STATE_FAILED');

  static const $core.List<CalculationExecutionState> values =
      <CalculationExecutionState>[
    CALCULATION_EXECUTION_STATE_QUEUED,
    CALCULATION_EXECUTION_STATE_RUNNING,
    CALCULATION_EXECUTION_STATE_CANCELED,
    CALCULATION_EXECUTION_STATE_CANCELING,
    CALCULATION_EXECUTION_STATE_CREATING,
    CALCULATION_EXECUTION_STATE_COMPLETED,
    CALCULATION_EXECUTION_STATE_CREATED,
    CALCULATION_EXECUTION_STATE_FAILED,
  ];

  static final $core.List<CalculationExecutionState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static CalculationExecutionState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CalculationExecutionState._(super.value, super.name);
}

class CapacityAllocationStatus extends $pb.ProtobufEnum {
  static const CapacityAllocationStatus CAPACITY_ALLOCATION_STATUS_PENDING =
      CapacityAllocationStatus._(
          0, _omitEnumNames ? '' : 'CAPACITY_ALLOCATION_STATUS_PENDING');
  static const CapacityAllocationStatus CAPACITY_ALLOCATION_STATUS_SUCCEEDED =
      CapacityAllocationStatus._(
          1, _omitEnumNames ? '' : 'CAPACITY_ALLOCATION_STATUS_SUCCEEDED');
  static const CapacityAllocationStatus CAPACITY_ALLOCATION_STATUS_FAILED =
      CapacityAllocationStatus._(
          2, _omitEnumNames ? '' : 'CAPACITY_ALLOCATION_STATUS_FAILED');

  static const $core.List<CapacityAllocationStatus> values =
      <CapacityAllocationStatus>[
    CAPACITY_ALLOCATION_STATUS_PENDING,
    CAPACITY_ALLOCATION_STATUS_SUCCEEDED,
    CAPACITY_ALLOCATION_STATUS_FAILED,
  ];

  static final $core.List<CapacityAllocationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static CapacityAllocationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CapacityAllocationStatus._(super.value, super.name);
}

class CapacityReservationStatus extends $pb.ProtobufEnum {
  static const CapacityReservationStatus CAPACITY_RESERVATION_STATUS_PENDING =
      CapacityReservationStatus._(
          0, _omitEnumNames ? '' : 'CAPACITY_RESERVATION_STATUS_PENDING');
  static const CapacityReservationStatus
      CAPACITY_RESERVATION_STATUS_UPDATE_PENDING = CapacityReservationStatus._(
          1,
          _omitEnumNames ? '' : 'CAPACITY_RESERVATION_STATUS_UPDATE_PENDING');
  static const CapacityReservationStatus CAPACITY_RESERVATION_STATUS_CANCELLED =
      CapacityReservationStatus._(
          2, _omitEnumNames ? '' : 'CAPACITY_RESERVATION_STATUS_CANCELLED');
  static const CapacityReservationStatus CAPACITY_RESERVATION_STATUS_ACTIVE =
      CapacityReservationStatus._(
          3, _omitEnumNames ? '' : 'CAPACITY_RESERVATION_STATUS_ACTIVE');
  static const CapacityReservationStatus
      CAPACITY_RESERVATION_STATUS_CANCELLING = CapacityReservationStatus._(
          4, _omitEnumNames ? '' : 'CAPACITY_RESERVATION_STATUS_CANCELLING');
  static const CapacityReservationStatus CAPACITY_RESERVATION_STATUS_FAILED =
      CapacityReservationStatus._(
          5, _omitEnumNames ? '' : 'CAPACITY_RESERVATION_STATUS_FAILED');

  static const $core.List<CapacityReservationStatus> values =
      <CapacityReservationStatus>[
    CAPACITY_RESERVATION_STATUS_PENDING,
    CAPACITY_RESERVATION_STATUS_UPDATE_PENDING,
    CAPACITY_RESERVATION_STATUS_CANCELLED,
    CAPACITY_RESERVATION_STATUS_ACTIVE,
    CAPACITY_RESERVATION_STATUS_CANCELLING,
    CAPACITY_RESERVATION_STATUS_FAILED,
  ];

  static final $core.List<CapacityReservationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static CapacityReservationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CapacityReservationStatus._(super.value, super.name);
}

class ColumnNullable extends $pb.ProtobufEnum {
  static const ColumnNullable COLUMN_NULLABLE_NOT_NULL =
      ColumnNullable._(0, _omitEnumNames ? '' : 'COLUMN_NULLABLE_NOT_NULL');
  static const ColumnNullable COLUMN_NULLABLE_UNKNOWN =
      ColumnNullable._(1, _omitEnumNames ? '' : 'COLUMN_NULLABLE_UNKNOWN');
  static const ColumnNullable COLUMN_NULLABLE_NULLABLE =
      ColumnNullable._(2, _omitEnumNames ? '' : 'COLUMN_NULLABLE_NULLABLE');

  static const $core.List<ColumnNullable> values = <ColumnNullable>[
    COLUMN_NULLABLE_NOT_NULL,
    COLUMN_NULLABLE_UNKNOWN,
    COLUMN_NULLABLE_NULLABLE,
  ];

  static final $core.List<ColumnNullable?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ColumnNullable? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ColumnNullable._(super.value, super.name);
}

class ConnectionType extends $pb.ProtobufEnum {
  static const ConnectionType CONNECTION_TYPE_DB2AS400 =
      ConnectionType._(0, _omitEnumNames ? '' : 'CONNECTION_TYPE_DB2AS400');
  static const ConnectionType CONNECTION_TYPE_DB2 =
      ConnectionType._(1, _omitEnumNames ? '' : 'CONNECTION_TYPE_DB2');
  static const ConnectionType CONNECTION_TYPE_SQLSERVER =
      ConnectionType._(2, _omitEnumNames ? '' : 'CONNECTION_TYPE_SQLSERVER');
  static const ConnectionType CONNECTION_TYPE_REDSHIFT =
      ConnectionType._(3, _omitEnumNames ? '' : 'CONNECTION_TYPE_REDSHIFT');
  static const ConnectionType CONNECTION_TYPE_MYSQL =
      ConnectionType._(4, _omitEnumNames ? '' : 'CONNECTION_TYPE_MYSQL');
  static const ConnectionType CONNECTION_TYPE_OPENSEARCH =
      ConnectionType._(5, _omitEnumNames ? '' : 'CONNECTION_TYPE_OPENSEARCH');
  static const ConnectionType CONNECTION_TYPE_DYNAMODB =
      ConnectionType._(6, _omitEnumNames ? '' : 'CONNECTION_TYPE_DYNAMODB');
  static const ConnectionType CONNECTION_TYPE_CMDB =
      ConnectionType._(7, _omitEnumNames ? '' : 'CONNECTION_TYPE_CMDB');
  static const ConnectionType CONNECTION_TYPE_DATALAKEGEN2 =
      ConnectionType._(8, _omitEnumNames ? '' : 'CONNECTION_TYPE_DATALAKEGEN2');
  static const ConnectionType CONNECTION_TYPE_DOCUMENTDB =
      ConnectionType._(9, _omitEnumNames ? '' : 'CONNECTION_TYPE_DOCUMENTDB');
  static const ConnectionType CONNECTION_TYPE_GOOGLECLOUDSTORAGE =
      ConnectionType._(
          10, _omitEnumNames ? '' : 'CONNECTION_TYPE_GOOGLECLOUDSTORAGE');
  static const ConnectionType CONNECTION_TYPE_SNOWFLAKE =
      ConnectionType._(11, _omitEnumNames ? '' : 'CONNECTION_TYPE_SNOWFLAKE');
  static const ConnectionType CONNECTION_TYPE_SAPHANA =
      ConnectionType._(12, _omitEnumNames ? '' : 'CONNECTION_TYPE_SAPHANA');
  static const ConnectionType CONNECTION_TYPE_TIMESTREAM =
      ConnectionType._(13, _omitEnumNames ? '' : 'CONNECTION_TYPE_TIMESTREAM');
  static const ConnectionType CONNECTION_TYPE_ORACLE =
      ConnectionType._(14, _omitEnumNames ? '' : 'CONNECTION_TYPE_ORACLE');
  static const ConnectionType CONNECTION_TYPE_BIGQUERY =
      ConnectionType._(15, _omitEnumNames ? '' : 'CONNECTION_TYPE_BIGQUERY');
  static const ConnectionType CONNECTION_TYPE_TPCDS =
      ConnectionType._(16, _omitEnumNames ? '' : 'CONNECTION_TYPE_TPCDS');
  static const ConnectionType CONNECTION_TYPE_SYNAPSE =
      ConnectionType._(17, _omitEnumNames ? '' : 'CONNECTION_TYPE_SYNAPSE');
  static const ConnectionType CONNECTION_TYPE_HBASE =
      ConnectionType._(18, _omitEnumNames ? '' : 'CONNECTION_TYPE_HBASE');
  static const ConnectionType CONNECTION_TYPE_POSTGRESQL =
      ConnectionType._(19, _omitEnumNames ? '' : 'CONNECTION_TYPE_POSTGRESQL');

  static const $core.List<ConnectionType> values = <ConnectionType>[
    CONNECTION_TYPE_DB2AS400,
    CONNECTION_TYPE_DB2,
    CONNECTION_TYPE_SQLSERVER,
    CONNECTION_TYPE_REDSHIFT,
    CONNECTION_TYPE_MYSQL,
    CONNECTION_TYPE_OPENSEARCH,
    CONNECTION_TYPE_DYNAMODB,
    CONNECTION_TYPE_CMDB,
    CONNECTION_TYPE_DATALAKEGEN2,
    CONNECTION_TYPE_DOCUMENTDB,
    CONNECTION_TYPE_GOOGLECLOUDSTORAGE,
    CONNECTION_TYPE_SNOWFLAKE,
    CONNECTION_TYPE_SAPHANA,
    CONNECTION_TYPE_TIMESTREAM,
    CONNECTION_TYPE_ORACLE,
    CONNECTION_TYPE_BIGQUERY,
    CONNECTION_TYPE_TPCDS,
    CONNECTION_TYPE_SYNAPSE,
    CONNECTION_TYPE_HBASE,
    CONNECTION_TYPE_POSTGRESQL,
  ];

  static final $core.List<ConnectionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 19);
  static ConnectionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionType._(super.value, super.name);
}

class DataCatalogStatus extends $pb.ProtobufEnum {
  static const DataCatalogStatus DATA_CATALOG_STATUS_DELETE_FAILED =
      DataCatalogStatus._(
          0, _omitEnumNames ? '' : 'DATA_CATALOG_STATUS_DELETE_FAILED');
  static const DataCatalogStatus DATA_CATALOG_STATUS_DELETE_COMPLETE =
      DataCatalogStatus._(
          1, _omitEnumNames ? '' : 'DATA_CATALOG_STATUS_DELETE_COMPLETE');
  static const DataCatalogStatus DATA_CATALOG_STATUS_CREATE_COMPLETE =
      DataCatalogStatus._(
          2, _omitEnumNames ? '' : 'DATA_CATALOG_STATUS_CREATE_COMPLETE');
  static const DataCatalogStatus
      DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_FAILED = DataCatalogStatus._(
          3,
          _omitEnumNames
              ? ''
              : 'DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_FAILED');
  static const DataCatalogStatus
      DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_IN_PROGRESS =
      DataCatalogStatus._(
          4,
          _omitEnumNames
              ? ''
              : 'DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_IN_PROGRESS');
  static const DataCatalogStatus DATA_CATALOG_STATUS_CREATE_IN_PROGRESS =
      DataCatalogStatus._(
          5, _omitEnumNames ? '' : 'DATA_CATALOG_STATUS_CREATE_IN_PROGRESS');
  static const DataCatalogStatus DATA_CATALOG_STATUS_CREATE_FAILED =
      DataCatalogStatus._(
          6, _omitEnumNames ? '' : 'DATA_CATALOG_STATUS_CREATE_FAILED');
  static const DataCatalogStatus DATA_CATALOG_STATUS_DELETE_IN_PROGRESS =
      DataCatalogStatus._(
          7, _omitEnumNames ? '' : 'DATA_CATALOG_STATUS_DELETE_IN_PROGRESS');
  static const DataCatalogStatus
      DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_COMPLETE = DataCatalogStatus._(
          8,
          _omitEnumNames
              ? ''
              : 'DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_COMPLETE');

  static const $core.List<DataCatalogStatus> values = <DataCatalogStatus>[
    DATA_CATALOG_STATUS_DELETE_FAILED,
    DATA_CATALOG_STATUS_DELETE_COMPLETE,
    DATA_CATALOG_STATUS_CREATE_COMPLETE,
    DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_FAILED,
    DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_IN_PROGRESS,
    DATA_CATALOG_STATUS_CREATE_IN_PROGRESS,
    DATA_CATALOG_STATUS_CREATE_FAILED,
    DATA_CATALOG_STATUS_DELETE_IN_PROGRESS,
    DATA_CATALOG_STATUS_CREATE_FAILED_CLEANUP_COMPLETE,
  ];

  static final $core.List<DataCatalogStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 8);
  static DataCatalogStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DataCatalogStatus._(super.value, super.name);
}

class DataCatalogType extends $pb.ProtobufEnum {
  static const DataCatalogType DATA_CATALOG_TYPE_FEDERATED =
      DataCatalogType._(0, _omitEnumNames ? '' : 'DATA_CATALOG_TYPE_FEDERATED');
  static const DataCatalogType DATA_CATALOG_TYPE_GLUE =
      DataCatalogType._(1, _omitEnumNames ? '' : 'DATA_CATALOG_TYPE_GLUE');
  static const DataCatalogType DATA_CATALOG_TYPE_LAMBDA =
      DataCatalogType._(2, _omitEnumNames ? '' : 'DATA_CATALOG_TYPE_LAMBDA');
  static const DataCatalogType DATA_CATALOG_TYPE_HIVE =
      DataCatalogType._(3, _omitEnumNames ? '' : 'DATA_CATALOG_TYPE_HIVE');

  static const $core.List<DataCatalogType> values = <DataCatalogType>[
    DATA_CATALOG_TYPE_FEDERATED,
    DATA_CATALOG_TYPE_GLUE,
    DATA_CATALOG_TYPE_LAMBDA,
    DATA_CATALOG_TYPE_HIVE,
  ];

  static final $core.List<DataCatalogType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static DataCatalogType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DataCatalogType._(super.value, super.name);
}

class EncryptionOption extends $pb.ProtobufEnum {
  static const EncryptionOption ENCRYPTION_OPTION_SSE_S3 =
      EncryptionOption._(0, _omitEnumNames ? '' : 'ENCRYPTION_OPTION_SSE_S3');
  static const EncryptionOption ENCRYPTION_OPTION_CSE_KMS =
      EncryptionOption._(1, _omitEnumNames ? '' : 'ENCRYPTION_OPTION_CSE_KMS');
  static const EncryptionOption ENCRYPTION_OPTION_SSE_KMS =
      EncryptionOption._(2, _omitEnumNames ? '' : 'ENCRYPTION_OPTION_SSE_KMS');

  static const $core.List<EncryptionOption> values = <EncryptionOption>[
    ENCRYPTION_OPTION_SSE_S3,
    ENCRYPTION_OPTION_CSE_KMS,
    ENCRYPTION_OPTION_SSE_KMS,
  ];

  static final $core.List<EncryptionOption?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EncryptionOption? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EncryptionOption._(super.value, super.name);
}

class ExecutorState extends $pb.ProtobufEnum {
  static const ExecutorState EXECUTOR_STATE_TERMINATING =
      ExecutorState._(0, _omitEnumNames ? '' : 'EXECUTOR_STATE_TERMINATING');
  static const ExecutorState EXECUTOR_STATE_TERMINATED =
      ExecutorState._(1, _omitEnumNames ? '' : 'EXECUTOR_STATE_TERMINATED');
  static const ExecutorState EXECUTOR_STATE_REGISTERED =
      ExecutorState._(2, _omitEnumNames ? '' : 'EXECUTOR_STATE_REGISTERED');
  static const ExecutorState EXECUTOR_STATE_CREATING =
      ExecutorState._(3, _omitEnumNames ? '' : 'EXECUTOR_STATE_CREATING');
  static const ExecutorState EXECUTOR_STATE_CREATED =
      ExecutorState._(4, _omitEnumNames ? '' : 'EXECUTOR_STATE_CREATED');
  static const ExecutorState EXECUTOR_STATE_FAILED =
      ExecutorState._(5, _omitEnumNames ? '' : 'EXECUTOR_STATE_FAILED');

  static const $core.List<ExecutorState> values = <ExecutorState>[
    EXECUTOR_STATE_TERMINATING,
    EXECUTOR_STATE_TERMINATED,
    EXECUTOR_STATE_REGISTERED,
    EXECUTOR_STATE_CREATING,
    EXECUTOR_STATE_CREATED,
    EXECUTOR_STATE_FAILED,
  ];

  static final $core.List<ExecutorState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ExecutorState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutorState._(super.value, super.name);
}

class ExecutorType extends $pb.ProtobufEnum {
  static const ExecutorType EXECUTOR_TYPE_GATEWAY =
      ExecutorType._(0, _omitEnumNames ? '' : 'EXECUTOR_TYPE_GATEWAY');
  static const ExecutorType EXECUTOR_TYPE_WORKER =
      ExecutorType._(1, _omitEnumNames ? '' : 'EXECUTOR_TYPE_WORKER');
  static const ExecutorType EXECUTOR_TYPE_COORDINATOR =
      ExecutorType._(2, _omitEnumNames ? '' : 'EXECUTOR_TYPE_COORDINATOR');

  static const $core.List<ExecutorType> values = <ExecutorType>[
    EXECUTOR_TYPE_GATEWAY,
    EXECUTOR_TYPE_WORKER,
    EXECUTOR_TYPE_COORDINATOR,
  ];

  static final $core.List<ExecutorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ExecutorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutorType._(super.value, super.name);
}

class NotebookType extends $pb.ProtobufEnum {
  static const NotebookType NOTEBOOK_TYPE_IPYNB =
      NotebookType._(0, _omitEnumNames ? '' : 'NOTEBOOK_TYPE_IPYNB');

  static const $core.List<NotebookType> values = <NotebookType>[
    NOTEBOOK_TYPE_IPYNB,
  ];

  static final $core.List<NotebookType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static NotebookType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NotebookType._(super.value, super.name);
}

class QueryExecutionState extends $pb.ProtobufEnum {
  static const QueryExecutionState QUERY_EXECUTION_STATE_QUEUED =
      QueryExecutionState._(
          0, _omitEnumNames ? '' : 'QUERY_EXECUTION_STATE_QUEUED');
  static const QueryExecutionState QUERY_EXECUTION_STATE_RUNNING =
      QueryExecutionState._(
          1, _omitEnumNames ? '' : 'QUERY_EXECUTION_STATE_RUNNING');
  static const QueryExecutionState QUERY_EXECUTION_STATE_SUCCEEDED =
      QueryExecutionState._(
          2, _omitEnumNames ? '' : 'QUERY_EXECUTION_STATE_SUCCEEDED');
  static const QueryExecutionState QUERY_EXECUTION_STATE_CANCELLED =
      QueryExecutionState._(
          3, _omitEnumNames ? '' : 'QUERY_EXECUTION_STATE_CANCELLED');
  static const QueryExecutionState QUERY_EXECUTION_STATE_FAILED =
      QueryExecutionState._(
          4, _omitEnumNames ? '' : 'QUERY_EXECUTION_STATE_FAILED');

  static const $core.List<QueryExecutionState> values = <QueryExecutionState>[
    QUERY_EXECUTION_STATE_QUEUED,
    QUERY_EXECUTION_STATE_RUNNING,
    QUERY_EXECUTION_STATE_SUCCEEDED,
    QUERY_EXECUTION_STATE_CANCELLED,
    QUERY_EXECUTION_STATE_FAILED,
  ];

  static final $core.List<QueryExecutionState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static QueryExecutionState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueryExecutionState._(super.value, super.name);
}

class QueryResultType extends $pb.ProtobufEnum {
  static const QueryResultType QUERY_RESULT_TYPE_DATA_MANIFEST =
      QueryResultType._(
          0, _omitEnumNames ? '' : 'QUERY_RESULT_TYPE_DATA_MANIFEST');
  static const QueryResultType QUERY_RESULT_TYPE_DATA_ROWS =
      QueryResultType._(1, _omitEnumNames ? '' : 'QUERY_RESULT_TYPE_DATA_ROWS');

  static const $core.List<QueryResultType> values = <QueryResultType>[
    QUERY_RESULT_TYPE_DATA_MANIFEST,
    QUERY_RESULT_TYPE_DATA_ROWS,
  ];

  static final $core.List<QueryResultType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static QueryResultType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueryResultType._(super.value, super.name);
}

class S3AclOption extends $pb.ProtobufEnum {
  static const S3AclOption S3_ACL_OPTION_BUCKET_OWNER_FULL_CONTROL =
      S3AclOption._(
          0, _omitEnumNames ? '' : 'S3_ACL_OPTION_BUCKET_OWNER_FULL_CONTROL');

  static const $core.List<S3AclOption> values = <S3AclOption>[
    S3_ACL_OPTION_BUCKET_OWNER_FULL_CONTROL,
  ];

  static final $core.List<S3AclOption?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static S3AclOption? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const S3AclOption._(super.value, super.name);
}

class SessionState extends $pb.ProtobufEnum {
  static const SessionState SESSION_STATE_TERMINATING =
      SessionState._(0, _omitEnumNames ? '' : 'SESSION_STATE_TERMINATING');
  static const SessionState SESSION_STATE_TERMINATED =
      SessionState._(1, _omitEnumNames ? '' : 'SESSION_STATE_TERMINATED');
  static const SessionState SESSION_STATE_DEGRADED =
      SessionState._(2, _omitEnumNames ? '' : 'SESSION_STATE_DEGRADED');
  static const SessionState SESSION_STATE_BUSY =
      SessionState._(3, _omitEnumNames ? '' : 'SESSION_STATE_BUSY');
  static const SessionState SESSION_STATE_IDLE =
      SessionState._(4, _omitEnumNames ? '' : 'SESSION_STATE_IDLE');
  static const SessionState SESSION_STATE_CREATING =
      SessionState._(5, _omitEnumNames ? '' : 'SESSION_STATE_CREATING');
  static const SessionState SESSION_STATE_CREATED =
      SessionState._(6, _omitEnumNames ? '' : 'SESSION_STATE_CREATED');
  static const SessionState SESSION_STATE_FAILED =
      SessionState._(7, _omitEnumNames ? '' : 'SESSION_STATE_FAILED');

  static const $core.List<SessionState> values = <SessionState>[
    SESSION_STATE_TERMINATING,
    SESSION_STATE_TERMINATED,
    SESSION_STATE_DEGRADED,
    SESSION_STATE_BUSY,
    SESSION_STATE_IDLE,
    SESSION_STATE_CREATING,
    SESSION_STATE_CREATED,
    SESSION_STATE_FAILED,
  ];

  static final $core.List<SessionState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static SessionState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SessionState._(super.value, super.name);
}

class StatementType extends $pb.ProtobufEnum {
  static const StatementType STATEMENT_TYPE_UTILITY =
      StatementType._(0, _omitEnumNames ? '' : 'STATEMENT_TYPE_UTILITY');
  static const StatementType STATEMENT_TYPE_DDL =
      StatementType._(1, _omitEnumNames ? '' : 'STATEMENT_TYPE_DDL');
  static const StatementType STATEMENT_TYPE_DML =
      StatementType._(2, _omitEnumNames ? '' : 'STATEMENT_TYPE_DML');

  static const $core.List<StatementType> values = <StatementType>[
    STATEMENT_TYPE_UTILITY,
    STATEMENT_TYPE_DDL,
    STATEMENT_TYPE_DML,
  ];

  static final $core.List<StatementType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static StatementType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StatementType._(super.value, super.name);
}

class ThrottleReason extends $pb.ProtobufEnum {
  static const ThrottleReason THROTTLE_REASON_CONCURRENT_QUERY_LIMIT_EXCEEDED =
      ThrottleReason._(
          0,
          _omitEnumNames
              ? ''
              : 'THROTTLE_REASON_CONCURRENT_QUERY_LIMIT_EXCEEDED');

  static const $core.List<ThrottleReason> values = <ThrottleReason>[
    THROTTLE_REASON_CONCURRENT_QUERY_LIMIT_EXCEEDED,
  ];

  static final $core.List<ThrottleReason?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ThrottleReason? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ThrottleReason._(super.value, super.name);
}

class WorkGroupState extends $pb.ProtobufEnum {
  static const WorkGroupState WORK_GROUP_STATE_DISABLED =
      WorkGroupState._(0, _omitEnumNames ? '' : 'WORK_GROUP_STATE_DISABLED');
  static const WorkGroupState WORK_GROUP_STATE_ENABLED =
      WorkGroupState._(1, _omitEnumNames ? '' : 'WORK_GROUP_STATE_ENABLED');

  static const $core.List<WorkGroupState> values = <WorkGroupState>[
    WORK_GROUP_STATE_DISABLED,
    WORK_GROUP_STATE_ENABLED,
  ];

  static final $core.List<WorkGroupState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static WorkGroupState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const WorkGroupState._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
