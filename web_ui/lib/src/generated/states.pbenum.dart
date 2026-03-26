// This is a generated file - do not edit.
//
// Generated from states.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class EncryptionType extends $pb.ProtobufEnum {
  static const EncryptionType ENCRYPTION_TYPE_CUSTOMER_MANAGED_KMS_KEY =
      EncryptionType._(
          0, _omitEnumNames ? '' : 'ENCRYPTION_TYPE_CUSTOMER_MANAGED_KMS_KEY');
  static const EncryptionType ENCRYPTION_TYPE_AWS_OWNED_KEY = EncryptionType._(
      1, _omitEnumNames ? '' : 'ENCRYPTION_TYPE_AWS_OWNED_KEY');

  static const $core.List<EncryptionType> values = <EncryptionType>[
    ENCRYPTION_TYPE_CUSTOMER_MANAGED_KMS_KEY,
    ENCRYPTION_TYPE_AWS_OWNED_KEY,
  ];

  static final $core.List<EncryptionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static EncryptionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EncryptionType._(super.value, super.name);
}

class ExecutionRedriveFilter extends $pb.ProtobufEnum {
  static const ExecutionRedriveFilter EXECUTION_REDRIVE_FILTER_REDRIVEN =
      ExecutionRedriveFilter._(
          0, _omitEnumNames ? '' : 'EXECUTION_REDRIVE_FILTER_REDRIVEN');
  static const ExecutionRedriveFilter EXECUTION_REDRIVE_FILTER_NOT_REDRIVEN =
      ExecutionRedriveFilter._(
          1, _omitEnumNames ? '' : 'EXECUTION_REDRIVE_FILTER_NOT_REDRIVEN');

  static const $core.List<ExecutionRedriveFilter> values =
      <ExecutionRedriveFilter>[
    EXECUTION_REDRIVE_FILTER_REDRIVEN,
    EXECUTION_REDRIVE_FILTER_NOT_REDRIVEN,
  ];

  static final $core.List<ExecutionRedriveFilter?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExecutionRedriveFilter? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutionRedriveFilter._(super.value, super.name);
}

class ExecutionRedriveStatus extends $pb.ProtobufEnum {
  static const ExecutionRedriveStatus
      EXECUTION_REDRIVE_STATUS_REDRIVABLE_BY_MAP_RUN = ExecutionRedriveStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'EXECUTION_REDRIVE_STATUS_REDRIVABLE_BY_MAP_RUN');
  static const ExecutionRedriveStatus EXECUTION_REDRIVE_STATUS_NOT_REDRIVABLE =
      ExecutionRedriveStatus._(
          1, _omitEnumNames ? '' : 'EXECUTION_REDRIVE_STATUS_NOT_REDRIVABLE');
  static const ExecutionRedriveStatus EXECUTION_REDRIVE_STATUS_REDRIVABLE =
      ExecutionRedriveStatus._(
          2, _omitEnumNames ? '' : 'EXECUTION_REDRIVE_STATUS_REDRIVABLE');

  static const $core.List<ExecutionRedriveStatus> values =
      <ExecutionRedriveStatus>[
    EXECUTION_REDRIVE_STATUS_REDRIVABLE_BY_MAP_RUN,
    EXECUTION_REDRIVE_STATUS_NOT_REDRIVABLE,
    EXECUTION_REDRIVE_STATUS_REDRIVABLE,
  ];

  static final $core.List<ExecutionRedriveStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ExecutionRedriveStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutionRedriveStatus._(super.value, super.name);
}

class ExecutionStatus extends $pb.ProtobufEnum {
  static const ExecutionStatus EXECUTION_STATUS_ABORTED =
      ExecutionStatus._(0, _omitEnumNames ? '' : 'EXECUTION_STATUS_ABORTED');
  static const ExecutionStatus EXECUTION_STATUS_RUNNING =
      ExecutionStatus._(1, _omitEnumNames ? '' : 'EXECUTION_STATUS_RUNNING');
  static const ExecutionStatus EXECUTION_STATUS_TIMED_OUT =
      ExecutionStatus._(2, _omitEnumNames ? '' : 'EXECUTION_STATUS_TIMED_OUT');
  static const ExecutionStatus EXECUTION_STATUS_SUCCEEDED =
      ExecutionStatus._(3, _omitEnumNames ? '' : 'EXECUTION_STATUS_SUCCEEDED');
  static const ExecutionStatus EXECUTION_STATUS_PENDING_REDRIVE =
      ExecutionStatus._(
          4, _omitEnumNames ? '' : 'EXECUTION_STATUS_PENDING_REDRIVE');
  static const ExecutionStatus EXECUTION_STATUS_FAILED =
      ExecutionStatus._(5, _omitEnumNames ? '' : 'EXECUTION_STATUS_FAILED');

  static const $core.List<ExecutionStatus> values = <ExecutionStatus>[
    EXECUTION_STATUS_ABORTED,
    EXECUTION_STATUS_RUNNING,
    EXECUTION_STATUS_TIMED_OUT,
    EXECUTION_STATUS_SUCCEEDED,
    EXECUTION_STATUS_PENDING_REDRIVE,
    EXECUTION_STATUS_FAILED,
  ];

  static final $core.List<ExecutionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ExecutionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutionStatus._(super.value, super.name);
}

class HistoryEventType extends $pb.ProtobufEnum {
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSUBMITTED =
      HistoryEventType._(
          0, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSUBMITTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPRUNSUCCEEDED =
      HistoryEventType._(
          1, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPRUNSUCCEEDED');
  static const HistoryEventType HISTORY_EVENT_TYPE_WAITSTATEEXITED =
      HistoryEventType._(
          2, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_WAITSTATEEXITED');
  static const HistoryEventType
      HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSCHEDULEFAILED = HistoryEventType._(
          3,
          _omitEnumNames
              ? ''
              : 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSCHEDULEFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPSTATESTARTED =
      HistoryEventType._(
          4, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPSTATESTARTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSTARTED =
      HistoryEventType._(
          5, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSTARTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPITERATIONFAILED =
      HistoryEventType._(
          6, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPITERATIONFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_PASSSTATEEXITED =
      HistoryEventType._(
          7, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_PASSSTATEEXITED');
  static const HistoryEventType HISTORY_EVENT_TYPE_SUCCEEDSTATEEXITED =
      HistoryEventType._(
          8, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_SUCCEEDSTATEEXITED');
  static const HistoryEventType HISTORY_EVENT_TYPE_ACTIVITYSTARTED =
      HistoryEventType._(
          9, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_ACTIVITYSTARTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_EXECUTIONREDRIVEN =
      HistoryEventType._(
          10, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_EXECUTIONREDRIVEN');
  static const HistoryEventType HISTORY_EVENT_TYPE_EXECUTIONABORTED =
      HistoryEventType._(
          11, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_EXECUTIONABORTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_EVALUATIONFAILED =
      HistoryEventType._(
          12, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_EVALUATIONFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPSTATEFAILED =
      HistoryEventType._(
          13, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPSTATEFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSCHEDULED =
      HistoryEventType._(
          14, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSCHEDULED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPRUNREDRIVEN =
      HistoryEventType._(
          15, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPRUNREDRIVEN');
  static const HistoryEventType HISTORY_EVENT_TYPE_PARALLELSTATEFAILED =
      HistoryEventType._(
          16, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_PARALLELSTATEFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPITERATIONSTARTED =
      HistoryEventType._(
          17, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPITERATIONSTARTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_ACTIVITYSCHEDULED =
      HistoryEventType._(
          18, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_ACTIVITYSCHEDULED');
  static const HistoryEventType HISTORY_EVENT_TYPE_EXECUTIONTIMEDOUT =
      HistoryEventType._(
          19, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_EXECUTIONTIMEDOUT');
  static const HistoryEventType HISTORY_EVENT_TYPE_ACTIVITYFAILED =
      HistoryEventType._(
          20, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_ACTIVITYFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSUBMITFAILED =
      HistoryEventType._(
          21, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSUBMITFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSTARTED =
      HistoryEventType._(
          22, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSTARTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSUCCEEDED =
      HistoryEventType._(
          23, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSUCCEEDED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPITERATIONSUCCEEDED =
      HistoryEventType._(
          24, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPITERATIONSUCCEEDED');
  static const HistoryEventType HISTORY_EVENT_TYPE_ACTIVITYSUCCEEDED =
      HistoryEventType._(
          25, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_ACTIVITYSUCCEEDED');
  static const HistoryEventType HISTORY_EVENT_TYPE_FAILSTATEENTERED =
      HistoryEventType._(
          26, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_FAILSTATEENTERED');
  static const HistoryEventType HISTORY_EVENT_TYPE_EXECUTIONFAILED =
      HistoryEventType._(
          27, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_EXECUTIONFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPSTATESUCCEEDED =
      HistoryEventType._(
          28, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPSTATESUCCEEDED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPRUNFAILED =
      HistoryEventType._(
          29, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPRUNFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSCHEDULED =
      HistoryEventType._(30,
          _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSCHEDULED');
  static const HistoryEventType HISTORY_EVENT_TYPE_PARALLELSTATEEXITED =
      HistoryEventType._(
          31, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_PARALLELSTATEEXITED');
  static const HistoryEventType HISTORY_EVENT_TYPE_EXECUTIONSTARTED =
      HistoryEventType._(
          32, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_EXECUTIONSTARTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_PARALLELSTATEENTERED =
      HistoryEventType._(
          33, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_PARALLELSTATEENTERED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSTATEENTERED =
      HistoryEventType._(
          34, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSTATEENTERED');
  static const HistoryEventType HISTORY_EVENT_TYPE_CHOICESTATEENTERED =
      HistoryEventType._(
          35, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_CHOICESTATEENTERED');
  static const HistoryEventType HISTORY_EVENT_TYPE_PASSSTATEENTERED =
      HistoryEventType._(
          36, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_PASSSTATEENTERED');
  static const HistoryEventType HISTORY_EVENT_TYPE_WAITSTATEABORTED =
      HistoryEventType._(
          37, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_WAITSTATEABORTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_PARALLELSTATEABORTED =
      HistoryEventType._(
          38, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_PARALLELSTATEABORTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_LAMBDAFUNCTIONTIMEDOUT =
      HistoryEventType._(39,
          _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONTIMEDOUT');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPRUNABORTED =
      HistoryEventType._(
          40, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPRUNABORTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPSTATEABORTED =
      HistoryEventType._(
          41, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPSTATEABORTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_ACTIVITYTIMEDOUT =
      HistoryEventType._(
          42, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_ACTIVITYTIMEDOUT');
  static const HistoryEventType HISTORY_EVENT_TYPE_SUCCEEDSTATEENTERED =
      HistoryEventType._(
          43, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_SUCCEEDSTATEENTERED');
  static const HistoryEventType HISTORY_EVENT_TYPE_ACTIVITYSCHEDULEFAILED =
      HistoryEventType._(44,
          _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_ACTIVITYSCHEDULEFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPRUNSTARTED =
      HistoryEventType._(
          45, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPRUNSTARTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPSTATEENTERED =
      HistoryEventType._(
          46, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPSTATEENTERED');
  static const HistoryEventType HISTORY_EVENT_TYPE_LAMBDAFUNCTIONFAILED =
      HistoryEventType._(
          47, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKFAILED =
      HistoryEventType._(
          48, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_CHOICESTATEEXITED =
      HistoryEventType._(
          49, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_CHOICESTATEEXITED');
  static const HistoryEventType HISTORY_EVENT_TYPE_EXECUTIONSUCCEEDED =
      HistoryEventType._(
          50, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_EXECUTIONSUCCEEDED');
  static const HistoryEventType HISTORY_EVENT_TYPE_WAITSTATEENTERED =
      HistoryEventType._(
          51, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_WAITSTATEENTERED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKTIMEDOUT =
      HistoryEventType._(
          52, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKTIMEDOUT');
  static const HistoryEventType HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSUCCEEDED =
      HistoryEventType._(53,
          _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSUCCEEDED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPITERATIONABORTED =
      HistoryEventType._(
          54, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPITERATIONABORTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_PARALLELSTATESTARTED =
      HistoryEventType._(
          55, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_PARALLELSTATESTARTED');
  static const HistoryEventType HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSTARTFAILED =
      HistoryEventType._(56,
          _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSTARTFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSTARTFAILED =
      HistoryEventType._(
          57, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSTARTFAILED');
  static const HistoryEventType HISTORY_EVENT_TYPE_PARALLELSTATESUCCEEDED =
      HistoryEventType._(58,
          _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_PARALLELSTATESUCCEEDED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSTATEEXITED =
      HistoryEventType._(
          59, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSTATEEXITED');
  static const HistoryEventType HISTORY_EVENT_TYPE_MAPSTATEEXITED =
      HistoryEventType._(
          60, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_MAPSTATEEXITED');
  static const HistoryEventType HISTORY_EVENT_TYPE_TASKSTATEABORTED =
      HistoryEventType._(
          61, _omitEnumNames ? '' : 'HISTORY_EVENT_TYPE_TASKSTATEABORTED');

  static const $core.List<HistoryEventType> values = <HistoryEventType>[
    HISTORY_EVENT_TYPE_TASKSUBMITTED,
    HISTORY_EVENT_TYPE_MAPRUNSUCCEEDED,
    HISTORY_EVENT_TYPE_WAITSTATEEXITED,
    HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSCHEDULEFAILED,
    HISTORY_EVENT_TYPE_MAPSTATESTARTED,
    HISTORY_EVENT_TYPE_TASKSTARTED,
    HISTORY_EVENT_TYPE_MAPITERATIONFAILED,
    HISTORY_EVENT_TYPE_PASSSTATEEXITED,
    HISTORY_EVENT_TYPE_SUCCEEDSTATEEXITED,
    HISTORY_EVENT_TYPE_ACTIVITYSTARTED,
    HISTORY_EVENT_TYPE_EXECUTIONREDRIVEN,
    HISTORY_EVENT_TYPE_EXECUTIONABORTED,
    HISTORY_EVENT_TYPE_EVALUATIONFAILED,
    HISTORY_EVENT_TYPE_MAPSTATEFAILED,
    HISTORY_EVENT_TYPE_TASKSCHEDULED,
    HISTORY_EVENT_TYPE_MAPRUNREDRIVEN,
    HISTORY_EVENT_TYPE_PARALLELSTATEFAILED,
    HISTORY_EVENT_TYPE_MAPITERATIONSTARTED,
    HISTORY_EVENT_TYPE_ACTIVITYSCHEDULED,
    HISTORY_EVENT_TYPE_EXECUTIONTIMEDOUT,
    HISTORY_EVENT_TYPE_ACTIVITYFAILED,
    HISTORY_EVENT_TYPE_TASKSUBMITFAILED,
    HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSTARTED,
    HISTORY_EVENT_TYPE_TASKSUCCEEDED,
    HISTORY_EVENT_TYPE_MAPITERATIONSUCCEEDED,
    HISTORY_EVENT_TYPE_ACTIVITYSUCCEEDED,
    HISTORY_EVENT_TYPE_FAILSTATEENTERED,
    HISTORY_EVENT_TYPE_EXECUTIONFAILED,
    HISTORY_EVENT_TYPE_MAPSTATESUCCEEDED,
    HISTORY_EVENT_TYPE_MAPRUNFAILED,
    HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSCHEDULED,
    HISTORY_EVENT_TYPE_PARALLELSTATEEXITED,
    HISTORY_EVENT_TYPE_EXECUTIONSTARTED,
    HISTORY_EVENT_TYPE_PARALLELSTATEENTERED,
    HISTORY_EVENT_TYPE_TASKSTATEENTERED,
    HISTORY_EVENT_TYPE_CHOICESTATEENTERED,
    HISTORY_EVENT_TYPE_PASSSTATEENTERED,
    HISTORY_EVENT_TYPE_WAITSTATEABORTED,
    HISTORY_EVENT_TYPE_PARALLELSTATEABORTED,
    HISTORY_EVENT_TYPE_LAMBDAFUNCTIONTIMEDOUT,
    HISTORY_EVENT_TYPE_MAPRUNABORTED,
    HISTORY_EVENT_TYPE_MAPSTATEABORTED,
    HISTORY_EVENT_TYPE_ACTIVITYTIMEDOUT,
    HISTORY_EVENT_TYPE_SUCCEEDSTATEENTERED,
    HISTORY_EVENT_TYPE_ACTIVITYSCHEDULEFAILED,
    HISTORY_EVENT_TYPE_MAPRUNSTARTED,
    HISTORY_EVENT_TYPE_MAPSTATEENTERED,
    HISTORY_EVENT_TYPE_LAMBDAFUNCTIONFAILED,
    HISTORY_EVENT_TYPE_TASKFAILED,
    HISTORY_EVENT_TYPE_CHOICESTATEEXITED,
    HISTORY_EVENT_TYPE_EXECUTIONSUCCEEDED,
    HISTORY_EVENT_TYPE_WAITSTATEENTERED,
    HISTORY_EVENT_TYPE_TASKTIMEDOUT,
    HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSUCCEEDED,
    HISTORY_EVENT_TYPE_MAPITERATIONABORTED,
    HISTORY_EVENT_TYPE_PARALLELSTATESTARTED,
    HISTORY_EVENT_TYPE_LAMBDAFUNCTIONSTARTFAILED,
    HISTORY_EVENT_TYPE_TASKSTARTFAILED,
    HISTORY_EVENT_TYPE_PARALLELSTATESUCCEEDED,
    HISTORY_EVENT_TYPE_TASKSTATEEXITED,
    HISTORY_EVENT_TYPE_MAPSTATEEXITED,
    HISTORY_EVENT_TYPE_TASKSTATEABORTED,
  ];

  static final $core.List<HistoryEventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 61);
  static HistoryEventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HistoryEventType._(super.value, super.name);
}

class IncludedData extends $pb.ProtobufEnum {
  static const IncludedData INCLUDED_DATA_METADATA_ONLY =
      IncludedData._(0, _omitEnumNames ? '' : 'INCLUDED_DATA_METADATA_ONLY');
  static const IncludedData INCLUDED_DATA_ALL_DATA =
      IncludedData._(1, _omitEnumNames ? '' : 'INCLUDED_DATA_ALL_DATA');

  static const $core.List<IncludedData> values = <IncludedData>[
    INCLUDED_DATA_METADATA_ONLY,
    INCLUDED_DATA_ALL_DATA,
  ];

  static final $core.List<IncludedData?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IncludedData? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IncludedData._(super.value, super.name);
}

class InspectionLevel extends $pb.ProtobufEnum {
  static const InspectionLevel INSPECTION_LEVEL_TRACE =
      InspectionLevel._(0, _omitEnumNames ? '' : 'INSPECTION_LEVEL_TRACE');
  static const InspectionLevel INSPECTION_LEVEL_DEBUG =
      InspectionLevel._(1, _omitEnumNames ? '' : 'INSPECTION_LEVEL_DEBUG');
  static const InspectionLevel INSPECTION_LEVEL_INFO =
      InspectionLevel._(2, _omitEnumNames ? '' : 'INSPECTION_LEVEL_INFO');

  static const $core.List<InspectionLevel> values = <InspectionLevel>[
    INSPECTION_LEVEL_TRACE,
    INSPECTION_LEVEL_DEBUG,
    INSPECTION_LEVEL_INFO,
  ];

  static final $core.List<InspectionLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static InspectionLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InspectionLevel._(super.value, super.name);
}

class KmsKeyState extends $pb.ProtobufEnum {
  static const KmsKeyState KMS_KEY_STATE_DISABLED =
      KmsKeyState._(0, _omitEnumNames ? '' : 'KMS_KEY_STATE_DISABLED');
  static const KmsKeyState KMS_KEY_STATE_PENDING_DELETION =
      KmsKeyState._(1, _omitEnumNames ? '' : 'KMS_KEY_STATE_PENDING_DELETION');
  static const KmsKeyState KMS_KEY_STATE_UNAVAILABLE =
      KmsKeyState._(2, _omitEnumNames ? '' : 'KMS_KEY_STATE_UNAVAILABLE');
  static const KmsKeyState KMS_KEY_STATE_PENDING_IMPORT =
      KmsKeyState._(3, _omitEnumNames ? '' : 'KMS_KEY_STATE_PENDING_IMPORT');
  static const KmsKeyState KMS_KEY_STATE_CREATING =
      KmsKeyState._(4, _omitEnumNames ? '' : 'KMS_KEY_STATE_CREATING');

  static const $core.List<KmsKeyState> values = <KmsKeyState>[
    KMS_KEY_STATE_DISABLED,
    KMS_KEY_STATE_PENDING_DELETION,
    KMS_KEY_STATE_UNAVAILABLE,
    KMS_KEY_STATE_PENDING_IMPORT,
    KMS_KEY_STATE_CREATING,
  ];

  static final $core.List<KmsKeyState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static KmsKeyState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KmsKeyState._(super.value, super.name);
}

class LogLevel extends $pb.ProtobufEnum {
  static const LogLevel LOG_LEVEL_FATAL =
      LogLevel._(0, _omitEnumNames ? '' : 'LOG_LEVEL_FATAL');
  static const LogLevel LOG_LEVEL_OFF =
      LogLevel._(1, _omitEnumNames ? '' : 'LOG_LEVEL_OFF');
  static const LogLevel LOG_LEVEL_ALL =
      LogLevel._(2, _omitEnumNames ? '' : 'LOG_LEVEL_ALL');
  static const LogLevel LOG_LEVEL_ERROR =
      LogLevel._(3, _omitEnumNames ? '' : 'LOG_LEVEL_ERROR');

  static const $core.List<LogLevel> values = <LogLevel>[
    LOG_LEVEL_FATAL,
    LOG_LEVEL_OFF,
    LOG_LEVEL_ALL,
    LOG_LEVEL_ERROR,
  ];

  static final $core.List<LogLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static LogLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LogLevel._(super.value, super.name);
}

class MapRunStatus extends $pb.ProtobufEnum {
  static const MapRunStatus MAP_RUN_STATUS_ABORTED =
      MapRunStatus._(0, _omitEnumNames ? '' : 'MAP_RUN_STATUS_ABORTED');
  static const MapRunStatus MAP_RUN_STATUS_RUNNING =
      MapRunStatus._(1, _omitEnumNames ? '' : 'MAP_RUN_STATUS_RUNNING');
  static const MapRunStatus MAP_RUN_STATUS_SUCCEEDED =
      MapRunStatus._(2, _omitEnumNames ? '' : 'MAP_RUN_STATUS_SUCCEEDED');
  static const MapRunStatus MAP_RUN_STATUS_FAILED =
      MapRunStatus._(3, _omitEnumNames ? '' : 'MAP_RUN_STATUS_FAILED');

  static const $core.List<MapRunStatus> values = <MapRunStatus>[
    MAP_RUN_STATUS_ABORTED,
    MAP_RUN_STATUS_RUNNING,
    MAP_RUN_STATUS_SUCCEEDED,
    MAP_RUN_STATUS_FAILED,
  ];

  static final $core.List<MapRunStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static MapRunStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MapRunStatus._(super.value, super.name);
}

class MockResponseValidationMode extends $pb.ProtobufEnum {
  static const MockResponseValidationMode
      MOCK_RESPONSE_VALIDATION_MODE_PRESENT = MockResponseValidationMode._(
          0, _omitEnumNames ? '' : 'MOCK_RESPONSE_VALIDATION_MODE_PRESENT');
  static const MockResponseValidationMode MOCK_RESPONSE_VALIDATION_MODE_NONE =
      MockResponseValidationMode._(
          1, _omitEnumNames ? '' : 'MOCK_RESPONSE_VALIDATION_MODE_NONE');
  static const MockResponseValidationMode MOCK_RESPONSE_VALIDATION_MODE_STRICT =
      MockResponseValidationMode._(
          2, _omitEnumNames ? '' : 'MOCK_RESPONSE_VALIDATION_MODE_STRICT');

  static const $core.List<MockResponseValidationMode> values =
      <MockResponseValidationMode>[
    MOCK_RESPONSE_VALIDATION_MODE_PRESENT,
    MOCK_RESPONSE_VALIDATION_MODE_NONE,
    MOCK_RESPONSE_VALIDATION_MODE_STRICT,
  ];

  static final $core.List<MockResponseValidationMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static MockResponseValidationMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MockResponseValidationMode._(super.value, super.name);
}

class StateMachineStatus extends $pb.ProtobufEnum {
  static const StateMachineStatus STATE_MACHINE_STATUS_ACTIVE =
      StateMachineStatus._(
          0, _omitEnumNames ? '' : 'STATE_MACHINE_STATUS_ACTIVE');
  static const StateMachineStatus STATE_MACHINE_STATUS_DELETING =
      StateMachineStatus._(
          1, _omitEnumNames ? '' : 'STATE_MACHINE_STATUS_DELETING');

  static const $core.List<StateMachineStatus> values = <StateMachineStatus>[
    STATE_MACHINE_STATUS_ACTIVE,
    STATE_MACHINE_STATUS_DELETING,
  ];

  static final $core.List<StateMachineStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static StateMachineStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StateMachineStatus._(super.value, super.name);
}

class StateMachineType extends $pb.ProtobufEnum {
  static const StateMachineType STATE_MACHINE_TYPE_EXPRESS =
      StateMachineType._(0, _omitEnumNames ? '' : 'STATE_MACHINE_TYPE_EXPRESS');
  static const StateMachineType STATE_MACHINE_TYPE_STANDARD =
      StateMachineType._(
          1, _omitEnumNames ? '' : 'STATE_MACHINE_TYPE_STANDARD');

  static const $core.List<StateMachineType> values = <StateMachineType>[
    STATE_MACHINE_TYPE_EXPRESS,
    STATE_MACHINE_TYPE_STANDARD,
  ];

  static final $core.List<StateMachineType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static StateMachineType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StateMachineType._(super.value, super.name);
}

class SyncExecutionStatus extends $pb.ProtobufEnum {
  static const SyncExecutionStatus SYNC_EXECUTION_STATUS_TIMED_OUT =
      SyncExecutionStatus._(
          0, _omitEnumNames ? '' : 'SYNC_EXECUTION_STATUS_TIMED_OUT');
  static const SyncExecutionStatus SYNC_EXECUTION_STATUS_SUCCEEDED =
      SyncExecutionStatus._(
          1, _omitEnumNames ? '' : 'SYNC_EXECUTION_STATUS_SUCCEEDED');
  static const SyncExecutionStatus SYNC_EXECUTION_STATUS_FAILED =
      SyncExecutionStatus._(
          2, _omitEnumNames ? '' : 'SYNC_EXECUTION_STATUS_FAILED');

  static const $core.List<SyncExecutionStatus> values = <SyncExecutionStatus>[
    SYNC_EXECUTION_STATUS_TIMED_OUT,
    SYNC_EXECUTION_STATUS_SUCCEEDED,
    SYNC_EXECUTION_STATUS_FAILED,
  ];

  static final $core.List<SyncExecutionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static SyncExecutionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SyncExecutionStatus._(super.value, super.name);
}

class TestExecutionStatus extends $pb.ProtobufEnum {
  static const TestExecutionStatus TEST_EXECUTION_STATUS_RETRIABLE =
      TestExecutionStatus._(
          0, _omitEnumNames ? '' : 'TEST_EXECUTION_STATUS_RETRIABLE');
  static const TestExecutionStatus TEST_EXECUTION_STATUS_CAUGHT_ERROR =
      TestExecutionStatus._(
          1, _omitEnumNames ? '' : 'TEST_EXECUTION_STATUS_CAUGHT_ERROR');
  static const TestExecutionStatus TEST_EXECUTION_STATUS_SUCCEEDED =
      TestExecutionStatus._(
          2, _omitEnumNames ? '' : 'TEST_EXECUTION_STATUS_SUCCEEDED');
  static const TestExecutionStatus TEST_EXECUTION_STATUS_FAILED =
      TestExecutionStatus._(
          3, _omitEnumNames ? '' : 'TEST_EXECUTION_STATUS_FAILED');

  static const $core.List<TestExecutionStatus> values = <TestExecutionStatus>[
    TEST_EXECUTION_STATUS_RETRIABLE,
    TEST_EXECUTION_STATUS_CAUGHT_ERROR,
    TEST_EXECUTION_STATUS_SUCCEEDED,
    TEST_EXECUTION_STATUS_FAILED,
  ];

  static final $core.List<TestExecutionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static TestExecutionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TestExecutionStatus._(super.value, super.name);
}

class ValidateStateMachineDefinitionResultCode extends $pb.ProtobufEnum {
  static const ValidateStateMachineDefinitionResultCode
      VALIDATE_STATE_MACHINE_DEFINITION_RESULT_CODE_OK =
      ValidateStateMachineDefinitionResultCode._(
          0,
          _omitEnumNames
              ? ''
              : 'VALIDATE_STATE_MACHINE_DEFINITION_RESULT_CODE_OK');
  static const ValidateStateMachineDefinitionResultCode
      VALIDATE_STATE_MACHINE_DEFINITION_RESULT_CODE_FAIL =
      ValidateStateMachineDefinitionResultCode._(
          1,
          _omitEnumNames
              ? ''
              : 'VALIDATE_STATE_MACHINE_DEFINITION_RESULT_CODE_FAIL');

  static const $core.List<ValidateStateMachineDefinitionResultCode> values =
      <ValidateStateMachineDefinitionResultCode>[
    VALIDATE_STATE_MACHINE_DEFINITION_RESULT_CODE_OK,
    VALIDATE_STATE_MACHINE_DEFINITION_RESULT_CODE_FAIL,
  ];

  static final $core.List<ValidateStateMachineDefinitionResultCode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ValidateStateMachineDefinitionResultCode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ValidateStateMachineDefinitionResultCode._(super.value, super.name);
}

class ValidateStateMachineDefinitionSeverity extends $pb.ProtobufEnum {
  static const ValidateStateMachineDefinitionSeverity
      VALIDATE_STATE_MACHINE_DEFINITION_SEVERITY_WARNING =
      ValidateStateMachineDefinitionSeverity._(
          0,
          _omitEnumNames
              ? ''
              : 'VALIDATE_STATE_MACHINE_DEFINITION_SEVERITY_WARNING');
  static const ValidateStateMachineDefinitionSeverity
      VALIDATE_STATE_MACHINE_DEFINITION_SEVERITY_ERROR =
      ValidateStateMachineDefinitionSeverity._(
          1,
          _omitEnumNames
              ? ''
              : 'VALIDATE_STATE_MACHINE_DEFINITION_SEVERITY_ERROR');

  static const $core.List<ValidateStateMachineDefinitionSeverity> values =
      <ValidateStateMachineDefinitionSeverity>[
    VALIDATE_STATE_MACHINE_DEFINITION_SEVERITY_WARNING,
    VALIDATE_STATE_MACHINE_DEFINITION_SEVERITY_ERROR,
  ];

  static final $core.List<ValidateStateMachineDefinitionSeverity?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ValidateStateMachineDefinitionSeverity? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ValidateStateMachineDefinitionSeverity._(super.value, super.name);
}

class ValidationExceptionReason extends $pb.ProtobufEnum {
  static const ValidationExceptionReason
      VALIDATION_EXCEPTION_REASON_API_DOES_NOT_SUPPORT_LABELED_ARNS =
      ValidationExceptionReason._(
          0,
          _omitEnumNames
              ? ''
              : 'VALIDATION_EXCEPTION_REASON_API_DOES_NOT_SUPPORT_LABELED_ARNS');
  static const ValidationExceptionReason
      VALIDATION_EXCEPTION_REASON_INVALID_ROUTING_CONFIGURATION =
      ValidationExceptionReason._(
          1,
          _omitEnumNames
              ? ''
              : 'VALIDATION_EXCEPTION_REASON_INVALID_ROUTING_CONFIGURATION');
  static const ValidationExceptionReason
      VALIDATION_EXCEPTION_REASON_CANNOT_UPDATE_COMPLETED_MAP_RUN =
      ValidationExceptionReason._(
          2,
          _omitEnumNames
              ? ''
              : 'VALIDATION_EXCEPTION_REASON_CANNOT_UPDATE_COMPLETED_MAP_RUN');
  static const ValidationExceptionReason
      VALIDATION_EXCEPTION_REASON_MISSING_REQUIRED_PARAMETER =
      ValidationExceptionReason._(
          3,
          _omitEnumNames
              ? ''
              : 'VALIDATION_EXCEPTION_REASON_MISSING_REQUIRED_PARAMETER');

  static const $core.List<ValidationExceptionReason> values =
      <ValidationExceptionReason>[
    VALIDATION_EXCEPTION_REASON_API_DOES_NOT_SUPPORT_LABELED_ARNS,
    VALIDATION_EXCEPTION_REASON_INVALID_ROUTING_CONFIGURATION,
    VALIDATION_EXCEPTION_REASON_CANNOT_UPDATE_COMPLETED_MAP_RUN,
    VALIDATION_EXCEPTION_REASON_MISSING_REQUIRED_PARAMETER,
  ];

  static final $core.List<ValidationExceptionReason?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ValidationExceptionReason? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ValidationExceptionReason._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
