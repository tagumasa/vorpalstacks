// This is a generated file - do not edit.
//
// Generated from ingest.timestream.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'ingest.timestream.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'ingest.timestream.pbenum.dart';

class AccessDeniedException extends $pb.GeneratedMessage {
  factory AccessDeniedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  AccessDeniedException._();

  factory AccessDeniedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccessDeniedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccessDeniedException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccessDeniedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccessDeniedException copyWith(
          void Function(AccessDeniedException) updates) =>
      super.copyWith((message) => updates(message as AccessDeniedException))
          as AccessDeniedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccessDeniedException create() => AccessDeniedException._();
  @$core.override
  AccessDeniedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccessDeniedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccessDeniedException>(create);
  static AccessDeniedException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class BatchLoadProgressReport extends $pb.GeneratedMessage {
  factory BatchLoadProgressReport({
    $fixnum.Int64? recordsprocessed,
    $fixnum.Int64? filefailures,
    $fixnum.Int64? recordsingested,
    $fixnum.Int64? bytesmetered,
    $fixnum.Int64? recordingestionfailures,
    $fixnum.Int64? parsefailures,
  }) {
    final result = create();
    if (recordsprocessed != null) result.recordsprocessed = recordsprocessed;
    if (filefailures != null) result.filefailures = filefailures;
    if (recordsingested != null) result.recordsingested = recordsingested;
    if (bytesmetered != null) result.bytesmetered = bytesmetered;
    if (recordingestionfailures != null)
      result.recordingestionfailures = recordingestionfailures;
    if (parsefailures != null) result.parsefailures = parsefailures;
    return result;
  }

  BatchLoadProgressReport._();

  factory BatchLoadProgressReport.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchLoadProgressReport.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchLoadProgressReport',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aInt64(27133172, _omitFieldNames ? '' : 'recordsprocessed')
    ..aInt64(44779969, _omitFieldNames ? '' : 'filefailures')
    ..aInt64(159925469, _omitFieldNames ? '' : 'recordsingested')
    ..aInt64(232712245, _omitFieldNames ? '' : 'bytesmetered')
    ..aInt64(262421066, _omitFieldNames ? '' : 'recordingestionfailures')
    ..aInt64(330475280, _omitFieldNames ? '' : 'parsefailures')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchLoadProgressReport clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchLoadProgressReport copyWith(
          void Function(BatchLoadProgressReport) updates) =>
      super.copyWith((message) => updates(message as BatchLoadProgressReport))
          as BatchLoadProgressReport;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchLoadProgressReport create() => BatchLoadProgressReport._();
  @$core.override
  BatchLoadProgressReport createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchLoadProgressReport getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchLoadProgressReport>(create);
  static BatchLoadProgressReport? _defaultInstance;

  @$pb.TagNumber(27133172)
  $fixnum.Int64 get recordsprocessed => $_getI64(0);
  @$pb.TagNumber(27133172)
  set recordsprocessed($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(27133172)
  $core.bool hasRecordsprocessed() => $_has(0);
  @$pb.TagNumber(27133172)
  void clearRecordsprocessed() => $_clearField(27133172);

  @$pb.TagNumber(44779969)
  $fixnum.Int64 get filefailures => $_getI64(1);
  @$pb.TagNumber(44779969)
  set filefailures($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(44779969)
  $core.bool hasFilefailures() => $_has(1);
  @$pb.TagNumber(44779969)
  void clearFilefailures() => $_clearField(44779969);

  @$pb.TagNumber(159925469)
  $fixnum.Int64 get recordsingested => $_getI64(2);
  @$pb.TagNumber(159925469)
  set recordsingested($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(159925469)
  $core.bool hasRecordsingested() => $_has(2);
  @$pb.TagNumber(159925469)
  void clearRecordsingested() => $_clearField(159925469);

  @$pb.TagNumber(232712245)
  $fixnum.Int64 get bytesmetered => $_getI64(3);
  @$pb.TagNumber(232712245)
  set bytesmetered($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(232712245)
  $core.bool hasBytesmetered() => $_has(3);
  @$pb.TagNumber(232712245)
  void clearBytesmetered() => $_clearField(232712245);

  @$pb.TagNumber(262421066)
  $fixnum.Int64 get recordingestionfailures => $_getI64(4);
  @$pb.TagNumber(262421066)
  set recordingestionfailures($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(262421066)
  $core.bool hasRecordingestionfailures() => $_has(4);
  @$pb.TagNumber(262421066)
  void clearRecordingestionfailures() => $_clearField(262421066);

  @$pb.TagNumber(330475280)
  $fixnum.Int64 get parsefailures => $_getI64(5);
  @$pb.TagNumber(330475280)
  set parsefailures($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(330475280)
  $core.bool hasParsefailures() => $_has(5);
  @$pb.TagNumber(330475280)
  void clearParsefailures() => $_clearField(330475280);
}

class BatchLoadTask extends $pb.GeneratedMessage {
  factory BatchLoadTask({
    $core.String? taskid,
    $core.String? databasename,
    $core.String? creationtime,
    $core.String? lastupdatedtime,
    $core.String? tablename,
    $core.String? resumableuntil,
    BatchLoadStatus? taskstatus,
  }) {
    final result = create();
    if (taskid != null) result.taskid = taskid;
    if (databasename != null) result.databasename = databasename;
    if (creationtime != null) result.creationtime = creationtime;
    if (lastupdatedtime != null) result.lastupdatedtime = lastupdatedtime;
    if (tablename != null) result.tablename = tablename;
    if (resumableuntil != null) result.resumableuntil = resumableuntil;
    if (taskstatus != null) result.taskstatus = taskstatus;
    return result;
  }

  BatchLoadTask._();

  factory BatchLoadTask.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchLoadTask.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchLoadTask',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(18532514, _omitFieldNames ? '' : 'taskid')
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(177046166, _omitFieldNames ? '' : 'lastupdatedtime')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(340504066, _omitFieldNames ? '' : 'resumableuntil')
    ..aE<BatchLoadStatus>(448718071, _omitFieldNames ? '' : 'taskstatus',
        enumValues: BatchLoadStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchLoadTask clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchLoadTask copyWith(void Function(BatchLoadTask) updates) =>
      super.copyWith((message) => updates(message as BatchLoadTask))
          as BatchLoadTask;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchLoadTask create() => BatchLoadTask._();
  @$core.override
  BatchLoadTask createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchLoadTask getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchLoadTask>(create);
  static BatchLoadTask? _defaultInstance;

  @$pb.TagNumber(18532514)
  $core.String get taskid => $_getSZ(0);
  @$pb.TagNumber(18532514)
  set taskid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18532514)
  $core.bool hasTaskid() => $_has(0);
  @$pb.TagNumber(18532514)
  void clearTaskid() => $_clearField(18532514);

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(1);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(1);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(177046166)
  $core.String get lastupdatedtime => $_getSZ(3);
  @$pb.TagNumber(177046166)
  set lastupdatedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(177046166)
  $core.bool hasLastupdatedtime() => $_has(3);
  @$pb.TagNumber(177046166)
  void clearLastupdatedtime() => $_clearField(177046166);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(4);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(4, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(4);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(340504066)
  $core.String get resumableuntil => $_getSZ(5);
  @$pb.TagNumber(340504066)
  set resumableuntil($core.String value) => $_setString(5, value);
  @$pb.TagNumber(340504066)
  $core.bool hasResumableuntil() => $_has(5);
  @$pb.TagNumber(340504066)
  void clearResumableuntil() => $_clearField(340504066);

  @$pb.TagNumber(448718071)
  BatchLoadStatus get taskstatus => $_getN(6);
  @$pb.TagNumber(448718071)
  set taskstatus(BatchLoadStatus value) => $_setField(448718071, value);
  @$pb.TagNumber(448718071)
  $core.bool hasTaskstatus() => $_has(6);
  @$pb.TagNumber(448718071)
  void clearTaskstatus() => $_clearField(448718071);
}

class BatchLoadTaskDescription extends $pb.GeneratedMessage {
  factory BatchLoadTaskDescription({
    $core.String? taskid,
    $core.String? creationtime,
    $core.String? lastupdatedtime,
    DataModelConfiguration? datamodelconfiguration,
    DataSourceConfiguration? datasourceconfiguration,
    $fixnum.Int64? recordversion,
    $core.String? targettablename,
    $core.String? resumableuntil,
    BatchLoadProgressReport? progressreport,
    ReportConfiguration? reportconfiguration,
    BatchLoadStatus? taskstatus,
    $core.String? targetdatabasename,
    $core.String? errormessage,
  }) {
    final result = create();
    if (taskid != null) result.taskid = taskid;
    if (creationtime != null) result.creationtime = creationtime;
    if (lastupdatedtime != null) result.lastupdatedtime = lastupdatedtime;
    if (datamodelconfiguration != null)
      result.datamodelconfiguration = datamodelconfiguration;
    if (datasourceconfiguration != null)
      result.datasourceconfiguration = datasourceconfiguration;
    if (recordversion != null) result.recordversion = recordversion;
    if (targettablename != null) result.targettablename = targettablename;
    if (resumableuntil != null) result.resumableuntil = resumableuntil;
    if (progressreport != null) result.progressreport = progressreport;
    if (reportconfiguration != null)
      result.reportconfiguration = reportconfiguration;
    if (taskstatus != null) result.taskstatus = taskstatus;
    if (targetdatabasename != null)
      result.targetdatabasename = targetdatabasename;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  BatchLoadTaskDescription._();

  factory BatchLoadTaskDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchLoadTaskDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchLoadTaskDescription',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(18532514, _omitFieldNames ? '' : 'taskid')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(177046166, _omitFieldNames ? '' : 'lastupdatedtime')
    ..aOM<DataModelConfiguration>(
        195075337, _omitFieldNames ? '' : 'datamodelconfiguration',
        subBuilder: DataModelConfiguration.create)
    ..aOM<DataSourceConfiguration>(
        256337787, _omitFieldNames ? '' : 'datasourceconfiguration',
        subBuilder: DataSourceConfiguration.create)
    ..aInt64(282309505, _omitFieldNames ? '' : 'recordversion')
    ..aOS(298767720, _omitFieldNames ? '' : 'targettablename')
    ..aOS(340504066, _omitFieldNames ? '' : 'resumableuntil')
    ..aOM<BatchLoadProgressReport>(
        345169253, _omitFieldNames ? '' : 'progressreport',
        subBuilder: BatchLoadProgressReport.create)
    ..aOM<ReportConfiguration>(
        410368276, _omitFieldNames ? '' : 'reportconfiguration',
        subBuilder: ReportConfiguration.create)
    ..aE<BatchLoadStatus>(448718071, _omitFieldNames ? '' : 'taskstatus',
        enumValues: BatchLoadStatus.values)
    ..aOS(454828599, _omitFieldNames ? '' : 'targetdatabasename')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchLoadTaskDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchLoadTaskDescription copyWith(
          void Function(BatchLoadTaskDescription) updates) =>
      super.copyWith((message) => updates(message as BatchLoadTaskDescription))
          as BatchLoadTaskDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchLoadTaskDescription create() => BatchLoadTaskDescription._();
  @$core.override
  BatchLoadTaskDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchLoadTaskDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchLoadTaskDescription>(create);
  static BatchLoadTaskDescription? _defaultInstance;

  @$pb.TagNumber(18532514)
  $core.String get taskid => $_getSZ(0);
  @$pb.TagNumber(18532514)
  set taskid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18532514)
  $core.bool hasTaskid() => $_has(0);
  @$pb.TagNumber(18532514)
  void clearTaskid() => $_clearField(18532514);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(177046166)
  $core.String get lastupdatedtime => $_getSZ(2);
  @$pb.TagNumber(177046166)
  set lastupdatedtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(177046166)
  $core.bool hasLastupdatedtime() => $_has(2);
  @$pb.TagNumber(177046166)
  void clearLastupdatedtime() => $_clearField(177046166);

  @$pb.TagNumber(195075337)
  DataModelConfiguration get datamodelconfiguration => $_getN(3);
  @$pb.TagNumber(195075337)
  set datamodelconfiguration(DataModelConfiguration value) =>
      $_setField(195075337, value);
  @$pb.TagNumber(195075337)
  $core.bool hasDatamodelconfiguration() => $_has(3);
  @$pb.TagNumber(195075337)
  void clearDatamodelconfiguration() => $_clearField(195075337);
  @$pb.TagNumber(195075337)
  DataModelConfiguration ensureDatamodelconfiguration() => $_ensure(3);

  @$pb.TagNumber(256337787)
  DataSourceConfiguration get datasourceconfiguration => $_getN(4);
  @$pb.TagNumber(256337787)
  set datasourceconfiguration(DataSourceConfiguration value) =>
      $_setField(256337787, value);
  @$pb.TagNumber(256337787)
  $core.bool hasDatasourceconfiguration() => $_has(4);
  @$pb.TagNumber(256337787)
  void clearDatasourceconfiguration() => $_clearField(256337787);
  @$pb.TagNumber(256337787)
  DataSourceConfiguration ensureDatasourceconfiguration() => $_ensure(4);

  @$pb.TagNumber(282309505)
  $fixnum.Int64 get recordversion => $_getI64(5);
  @$pb.TagNumber(282309505)
  set recordversion($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(282309505)
  $core.bool hasRecordversion() => $_has(5);
  @$pb.TagNumber(282309505)
  void clearRecordversion() => $_clearField(282309505);

  @$pb.TagNumber(298767720)
  $core.String get targettablename => $_getSZ(6);
  @$pb.TagNumber(298767720)
  set targettablename($core.String value) => $_setString(6, value);
  @$pb.TagNumber(298767720)
  $core.bool hasTargettablename() => $_has(6);
  @$pb.TagNumber(298767720)
  void clearTargettablename() => $_clearField(298767720);

  @$pb.TagNumber(340504066)
  $core.String get resumableuntil => $_getSZ(7);
  @$pb.TagNumber(340504066)
  set resumableuntil($core.String value) => $_setString(7, value);
  @$pb.TagNumber(340504066)
  $core.bool hasResumableuntil() => $_has(7);
  @$pb.TagNumber(340504066)
  void clearResumableuntil() => $_clearField(340504066);

  @$pb.TagNumber(345169253)
  BatchLoadProgressReport get progressreport => $_getN(8);
  @$pb.TagNumber(345169253)
  set progressreport(BatchLoadProgressReport value) =>
      $_setField(345169253, value);
  @$pb.TagNumber(345169253)
  $core.bool hasProgressreport() => $_has(8);
  @$pb.TagNumber(345169253)
  void clearProgressreport() => $_clearField(345169253);
  @$pb.TagNumber(345169253)
  BatchLoadProgressReport ensureProgressreport() => $_ensure(8);

  @$pb.TagNumber(410368276)
  ReportConfiguration get reportconfiguration => $_getN(9);
  @$pb.TagNumber(410368276)
  set reportconfiguration(ReportConfiguration value) =>
      $_setField(410368276, value);
  @$pb.TagNumber(410368276)
  $core.bool hasReportconfiguration() => $_has(9);
  @$pb.TagNumber(410368276)
  void clearReportconfiguration() => $_clearField(410368276);
  @$pb.TagNumber(410368276)
  ReportConfiguration ensureReportconfiguration() => $_ensure(9);

  @$pb.TagNumber(448718071)
  BatchLoadStatus get taskstatus => $_getN(10);
  @$pb.TagNumber(448718071)
  set taskstatus(BatchLoadStatus value) => $_setField(448718071, value);
  @$pb.TagNumber(448718071)
  $core.bool hasTaskstatus() => $_has(10);
  @$pb.TagNumber(448718071)
  void clearTaskstatus() => $_clearField(448718071);

  @$pb.TagNumber(454828599)
  $core.String get targetdatabasename => $_getSZ(11);
  @$pb.TagNumber(454828599)
  set targetdatabasename($core.String value) => $_setString(11, value);
  @$pb.TagNumber(454828599)
  $core.bool hasTargetdatabasename() => $_has(11);
  @$pb.TagNumber(454828599)
  void clearTargetdatabasename() => $_clearField(454828599);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(12);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(12, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(12);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class ConflictException extends $pb.GeneratedMessage {
  factory ConflictException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ConflictException._();

  factory ConflictException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConflictException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConflictException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConflictException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConflictException copyWith(void Function(ConflictException) updates) =>
      super.copyWith((message) => updates(message as ConflictException))
          as ConflictException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConflictException create() => ConflictException._();
  @$core.override
  ConflictException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConflictException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConflictException>(create);
  static ConflictException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CreateBatchLoadTaskRequest extends $pb.GeneratedMessage {
  factory CreateBatchLoadTaskRequest({
    $core.String? clienttoken,
    DataModelConfiguration? datamodelconfiguration,
    DataSourceConfiguration? datasourceconfiguration,
    $fixnum.Int64? recordversion,
    $core.String? targettablename,
    ReportConfiguration? reportconfiguration,
    $core.String? targetdatabasename,
  }) {
    final result = create();
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (datamodelconfiguration != null)
      result.datamodelconfiguration = datamodelconfiguration;
    if (datasourceconfiguration != null)
      result.datasourceconfiguration = datasourceconfiguration;
    if (recordversion != null) result.recordversion = recordversion;
    if (targettablename != null) result.targettablename = targettablename;
    if (reportconfiguration != null)
      result.reportconfiguration = reportconfiguration;
    if (targetdatabasename != null)
      result.targetdatabasename = targetdatabasename;
    return result;
  }

  CreateBatchLoadTaskRequest._();

  factory CreateBatchLoadTaskRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateBatchLoadTaskRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateBatchLoadTaskRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOM<DataModelConfiguration>(
        195075337, _omitFieldNames ? '' : 'datamodelconfiguration',
        subBuilder: DataModelConfiguration.create)
    ..aOM<DataSourceConfiguration>(
        256337787, _omitFieldNames ? '' : 'datasourceconfiguration',
        subBuilder: DataSourceConfiguration.create)
    ..aInt64(282309505, _omitFieldNames ? '' : 'recordversion')
    ..aOS(298767720, _omitFieldNames ? '' : 'targettablename')
    ..aOM<ReportConfiguration>(
        410368276, _omitFieldNames ? '' : 'reportconfiguration',
        subBuilder: ReportConfiguration.create)
    ..aOS(454828599, _omitFieldNames ? '' : 'targetdatabasename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBatchLoadTaskRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBatchLoadTaskRequest copyWith(
          void Function(CreateBatchLoadTaskRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateBatchLoadTaskRequest))
          as CreateBatchLoadTaskRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateBatchLoadTaskRequest create() => CreateBatchLoadTaskRequest._();
  @$core.override
  CreateBatchLoadTaskRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateBatchLoadTaskRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateBatchLoadTaskRequest>(create);
  static CreateBatchLoadTaskRequest? _defaultInstance;

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(0);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(0);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(195075337)
  DataModelConfiguration get datamodelconfiguration => $_getN(1);
  @$pb.TagNumber(195075337)
  set datamodelconfiguration(DataModelConfiguration value) =>
      $_setField(195075337, value);
  @$pb.TagNumber(195075337)
  $core.bool hasDatamodelconfiguration() => $_has(1);
  @$pb.TagNumber(195075337)
  void clearDatamodelconfiguration() => $_clearField(195075337);
  @$pb.TagNumber(195075337)
  DataModelConfiguration ensureDatamodelconfiguration() => $_ensure(1);

  @$pb.TagNumber(256337787)
  DataSourceConfiguration get datasourceconfiguration => $_getN(2);
  @$pb.TagNumber(256337787)
  set datasourceconfiguration(DataSourceConfiguration value) =>
      $_setField(256337787, value);
  @$pb.TagNumber(256337787)
  $core.bool hasDatasourceconfiguration() => $_has(2);
  @$pb.TagNumber(256337787)
  void clearDatasourceconfiguration() => $_clearField(256337787);
  @$pb.TagNumber(256337787)
  DataSourceConfiguration ensureDatasourceconfiguration() => $_ensure(2);

  @$pb.TagNumber(282309505)
  $fixnum.Int64 get recordversion => $_getI64(3);
  @$pb.TagNumber(282309505)
  set recordversion($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(282309505)
  $core.bool hasRecordversion() => $_has(3);
  @$pb.TagNumber(282309505)
  void clearRecordversion() => $_clearField(282309505);

  @$pb.TagNumber(298767720)
  $core.String get targettablename => $_getSZ(4);
  @$pb.TagNumber(298767720)
  set targettablename($core.String value) => $_setString(4, value);
  @$pb.TagNumber(298767720)
  $core.bool hasTargettablename() => $_has(4);
  @$pb.TagNumber(298767720)
  void clearTargettablename() => $_clearField(298767720);

  @$pb.TagNumber(410368276)
  ReportConfiguration get reportconfiguration => $_getN(5);
  @$pb.TagNumber(410368276)
  set reportconfiguration(ReportConfiguration value) =>
      $_setField(410368276, value);
  @$pb.TagNumber(410368276)
  $core.bool hasReportconfiguration() => $_has(5);
  @$pb.TagNumber(410368276)
  void clearReportconfiguration() => $_clearField(410368276);
  @$pb.TagNumber(410368276)
  ReportConfiguration ensureReportconfiguration() => $_ensure(5);

  @$pb.TagNumber(454828599)
  $core.String get targetdatabasename => $_getSZ(6);
  @$pb.TagNumber(454828599)
  set targetdatabasename($core.String value) => $_setString(6, value);
  @$pb.TagNumber(454828599)
  $core.bool hasTargetdatabasename() => $_has(6);
  @$pb.TagNumber(454828599)
  void clearTargetdatabasename() => $_clearField(454828599);
}

class CreateBatchLoadTaskResponse extends $pb.GeneratedMessage {
  factory CreateBatchLoadTaskResponse({
    $core.String? taskid,
  }) {
    final result = create();
    if (taskid != null) result.taskid = taskid;
    return result;
  }

  CreateBatchLoadTaskResponse._();

  factory CreateBatchLoadTaskResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateBatchLoadTaskResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateBatchLoadTaskResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(18532514, _omitFieldNames ? '' : 'taskid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBatchLoadTaskResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBatchLoadTaskResponse copyWith(
          void Function(CreateBatchLoadTaskResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateBatchLoadTaskResponse))
          as CreateBatchLoadTaskResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateBatchLoadTaskResponse create() =>
      CreateBatchLoadTaskResponse._();
  @$core.override
  CreateBatchLoadTaskResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateBatchLoadTaskResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateBatchLoadTaskResponse>(create);
  static CreateBatchLoadTaskResponse? _defaultInstance;

  @$pb.TagNumber(18532514)
  $core.String get taskid => $_getSZ(0);
  @$pb.TagNumber(18532514)
  set taskid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18532514)
  $core.bool hasTaskid() => $_has(0);
  @$pb.TagNumber(18532514)
  void clearTaskid() => $_clearField(18532514);
}

class CreateDatabaseRequest extends $pb.GeneratedMessage {
  factory CreateDatabaseRequest({
    $core.String? kmskeyid,
    $core.String? databasename,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (databasename != null) result.databasename = databasename;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  CreateDatabaseRequest._();

  factory CreateDatabaseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDatabaseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDatabaseRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDatabaseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDatabaseRequest copyWith(
          void Function(CreateDatabaseRequest) updates) =>
      super.copyWith((message) => updates(message as CreateDatabaseRequest))
          as CreateDatabaseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDatabaseRequest create() => CreateDatabaseRequest._();
  @$core.override
  CreateDatabaseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDatabaseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDatabaseRequest>(create);
  static CreateDatabaseRequest? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(1);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(1);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(2);
}

class CreateDatabaseResponse extends $pb.GeneratedMessage {
  factory CreateDatabaseResponse({
    Database? database,
  }) {
    final result = create();
    if (database != null) result.database = database;
    return result;
  }

  CreateDatabaseResponse._();

  factory CreateDatabaseResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDatabaseResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDatabaseResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<Database>(278147289, _omitFieldNames ? '' : 'database',
        subBuilder: Database.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDatabaseResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDatabaseResponse copyWith(
          void Function(CreateDatabaseResponse) updates) =>
      super.copyWith((message) => updates(message as CreateDatabaseResponse))
          as CreateDatabaseResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDatabaseResponse create() => CreateDatabaseResponse._();
  @$core.override
  CreateDatabaseResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDatabaseResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDatabaseResponse>(create);
  static CreateDatabaseResponse? _defaultInstance;

  @$pb.TagNumber(278147289)
  Database get database => $_getN(0);
  @$pb.TagNumber(278147289)
  set database(Database value) => $_setField(278147289, value);
  @$pb.TagNumber(278147289)
  $core.bool hasDatabase() => $_has(0);
  @$pb.TagNumber(278147289)
  void clearDatabase() => $_clearField(278147289);
  @$pb.TagNumber(278147289)
  Database ensureDatabase() => $_ensure(0);
}

class CreateTableRequest extends $pb.GeneratedMessage {
  factory CreateTableRequest({
    $core.String? databasename,
    MagneticStoreWriteProperties? magneticstorewriteproperties,
    RetentionProperties? retentionproperties,
    $core.String? tablename,
    $core.Iterable<Tag>? tags,
    Schema? schema,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (magneticstorewriteproperties != null)
      result.magneticstorewriteproperties = magneticstorewriteproperties;
    if (retentionproperties != null)
      result.retentionproperties = retentionproperties;
    if (tablename != null) result.tablename = tablename;
    if (tags != null) result.tags.addAll(tags);
    if (schema != null) result.schema = schema;
    return result;
  }

  CreateTableRequest._();

  factory CreateTableRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTableRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTableRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOM<MagneticStoreWriteProperties>(
        127209171, _omitFieldNames ? '' : 'magneticstorewriteproperties',
        subBuilder: MagneticStoreWriteProperties.create)
    ..aOM<RetentionProperties>(
        242841241, _omitFieldNames ? '' : 'retentionproperties',
        subBuilder: RetentionProperties.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOM<Schema>(412122455, _omitFieldNames ? '' : 'schema',
        subBuilder: Schema.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTableRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTableRequest copyWith(void Function(CreateTableRequest) updates) =>
      super.copyWith((message) => updates(message as CreateTableRequest))
          as CreateTableRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTableRequest create() => CreateTableRequest._();
  @$core.override
  CreateTableRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTableRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTableRequest>(create);
  static CreateTableRequest? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(127209171)
  MagneticStoreWriteProperties get magneticstorewriteproperties => $_getN(1);
  @$pb.TagNumber(127209171)
  set magneticstorewriteproperties(MagneticStoreWriteProperties value) =>
      $_setField(127209171, value);
  @$pb.TagNumber(127209171)
  $core.bool hasMagneticstorewriteproperties() => $_has(1);
  @$pb.TagNumber(127209171)
  void clearMagneticstorewriteproperties() => $_clearField(127209171);
  @$pb.TagNumber(127209171)
  MagneticStoreWriteProperties ensureMagneticstorewriteproperties() =>
      $_ensure(1);

  @$pb.TagNumber(242841241)
  RetentionProperties get retentionproperties => $_getN(2);
  @$pb.TagNumber(242841241)
  set retentionproperties(RetentionProperties value) =>
      $_setField(242841241, value);
  @$pb.TagNumber(242841241)
  $core.bool hasRetentionproperties() => $_has(2);
  @$pb.TagNumber(242841241)
  void clearRetentionproperties() => $_clearField(242841241);
  @$pb.TagNumber(242841241)
  RetentionProperties ensureRetentionproperties() => $_ensure(2);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(4);

  @$pb.TagNumber(412122455)
  Schema get schema => $_getN(5);
  @$pb.TagNumber(412122455)
  set schema(Schema value) => $_setField(412122455, value);
  @$pb.TagNumber(412122455)
  $core.bool hasSchema() => $_has(5);
  @$pb.TagNumber(412122455)
  void clearSchema() => $_clearField(412122455);
  @$pb.TagNumber(412122455)
  Schema ensureSchema() => $_ensure(5);
}

class CreateTableResponse extends $pb.GeneratedMessage {
  factory CreateTableResponse({
    Table? table,
  }) {
    final result = create();
    if (table != null) result.table = table;
    return result;
  }

  CreateTableResponse._();

  factory CreateTableResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTableResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTableResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<Table>(386722688, _omitFieldNames ? '' : 'table',
        subBuilder: Table.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTableResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTableResponse copyWith(void Function(CreateTableResponse) updates) =>
      super.copyWith((message) => updates(message as CreateTableResponse))
          as CreateTableResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTableResponse create() => CreateTableResponse._();
  @$core.override
  CreateTableResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTableResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTableResponse>(create);
  static CreateTableResponse? _defaultInstance;

  @$pb.TagNumber(386722688)
  Table get table => $_getN(0);
  @$pb.TagNumber(386722688)
  set table(Table value) => $_setField(386722688, value);
  @$pb.TagNumber(386722688)
  $core.bool hasTable() => $_has(0);
  @$pb.TagNumber(386722688)
  void clearTable() => $_clearField(386722688);
  @$pb.TagNumber(386722688)
  Table ensureTable() => $_ensure(0);
}

class CsvConfiguration extends $pb.GeneratedMessage {
  factory CsvConfiguration({
    $core.String? columnseparator,
    $core.String? escapechar,
    $core.String? quotechar,
    $core.bool? trimwhitespace,
    $core.String? nullvalue,
  }) {
    final result = create();
    if (columnseparator != null) result.columnseparator = columnseparator;
    if (escapechar != null) result.escapechar = escapechar;
    if (quotechar != null) result.quotechar = quotechar;
    if (trimwhitespace != null) result.trimwhitespace = trimwhitespace;
    if (nullvalue != null) result.nullvalue = nullvalue;
    return result;
  }

  CsvConfiguration._();

  factory CsvConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CsvConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CsvConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(8188515, _omitFieldNames ? '' : 'columnseparator')
    ..aOS(210340939, _omitFieldNames ? '' : 'escapechar')
    ..aOS(242375728, _omitFieldNames ? '' : 'quotechar')
    ..aOB(386966069, _omitFieldNames ? '' : 'trimwhitespace')
    ..aOS(440981694, _omitFieldNames ? '' : 'nullvalue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CsvConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CsvConfiguration copyWith(void Function(CsvConfiguration) updates) =>
      super.copyWith((message) => updates(message as CsvConfiguration))
          as CsvConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CsvConfiguration create() => CsvConfiguration._();
  @$core.override
  CsvConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CsvConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CsvConfiguration>(create);
  static CsvConfiguration? _defaultInstance;

  @$pb.TagNumber(8188515)
  $core.String get columnseparator => $_getSZ(0);
  @$pb.TagNumber(8188515)
  set columnseparator($core.String value) => $_setString(0, value);
  @$pb.TagNumber(8188515)
  $core.bool hasColumnseparator() => $_has(0);
  @$pb.TagNumber(8188515)
  void clearColumnseparator() => $_clearField(8188515);

  @$pb.TagNumber(210340939)
  $core.String get escapechar => $_getSZ(1);
  @$pb.TagNumber(210340939)
  set escapechar($core.String value) => $_setString(1, value);
  @$pb.TagNumber(210340939)
  $core.bool hasEscapechar() => $_has(1);
  @$pb.TagNumber(210340939)
  void clearEscapechar() => $_clearField(210340939);

  @$pb.TagNumber(242375728)
  $core.String get quotechar => $_getSZ(2);
  @$pb.TagNumber(242375728)
  set quotechar($core.String value) => $_setString(2, value);
  @$pb.TagNumber(242375728)
  $core.bool hasQuotechar() => $_has(2);
  @$pb.TagNumber(242375728)
  void clearQuotechar() => $_clearField(242375728);

  @$pb.TagNumber(386966069)
  $core.bool get trimwhitespace => $_getBF(3);
  @$pb.TagNumber(386966069)
  set trimwhitespace($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(386966069)
  $core.bool hasTrimwhitespace() => $_has(3);
  @$pb.TagNumber(386966069)
  void clearTrimwhitespace() => $_clearField(386966069);

  @$pb.TagNumber(440981694)
  $core.String get nullvalue => $_getSZ(4);
  @$pb.TagNumber(440981694)
  set nullvalue($core.String value) => $_setString(4, value);
  @$pb.TagNumber(440981694)
  $core.bool hasNullvalue() => $_has(4);
  @$pb.TagNumber(440981694)
  void clearNullvalue() => $_clearField(440981694);
}

class DataModel extends $pb.GeneratedMessage {
  factory DataModel({
    $core.String? measurenamecolumn,
    TimeUnit? timeunit,
    $core.Iterable<DimensionMapping>? dimensionmappings,
    MultiMeasureMappings? multimeasuremappings,
    $core.String? timecolumn,
    $core.Iterable<MixedMeasureMapping>? mixedmeasuremappings,
  }) {
    final result = create();
    if (measurenamecolumn != null) result.measurenamecolumn = measurenamecolumn;
    if (timeunit != null) result.timeunit = timeunit;
    if (dimensionmappings != null)
      result.dimensionmappings.addAll(dimensionmappings);
    if (multimeasuremappings != null)
      result.multimeasuremappings = multimeasuremappings;
    if (timecolumn != null) result.timecolumn = timecolumn;
    if (mixedmeasuremappings != null)
      result.mixedmeasuremappings.addAll(mixedmeasuremappings);
    return result;
  }

  DataModel._();

  factory DataModel.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataModel.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataModel',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(112141775, _omitFieldNames ? '' : 'measurenamecolumn')
    ..aE<TimeUnit>(181686379, _omitFieldNames ? '' : 'timeunit',
        enumValues: TimeUnit.values)
    ..pPM<DimensionMapping>(
        224443741, _omitFieldNames ? '' : 'dimensionmappings',
        subBuilder: DimensionMapping.create)
    ..aOM<MultiMeasureMappings>(
        501736394, _omitFieldNames ? '' : 'multimeasuremappings',
        subBuilder: MultiMeasureMappings.create)
    ..aOS(519449367, _omitFieldNames ? '' : 'timecolumn')
    ..pPM<MixedMeasureMapping>(
        521774144, _omitFieldNames ? '' : 'mixedmeasuremappings',
        subBuilder: MixedMeasureMapping.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataModel clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataModel copyWith(void Function(DataModel) updates) =>
      super.copyWith((message) => updates(message as DataModel)) as DataModel;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataModel create() => DataModel._();
  @$core.override
  DataModel createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataModel getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DataModel>(create);
  static DataModel? _defaultInstance;

  @$pb.TagNumber(112141775)
  $core.String get measurenamecolumn => $_getSZ(0);
  @$pb.TagNumber(112141775)
  set measurenamecolumn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112141775)
  $core.bool hasMeasurenamecolumn() => $_has(0);
  @$pb.TagNumber(112141775)
  void clearMeasurenamecolumn() => $_clearField(112141775);

  @$pb.TagNumber(181686379)
  TimeUnit get timeunit => $_getN(1);
  @$pb.TagNumber(181686379)
  set timeunit(TimeUnit value) => $_setField(181686379, value);
  @$pb.TagNumber(181686379)
  $core.bool hasTimeunit() => $_has(1);
  @$pb.TagNumber(181686379)
  void clearTimeunit() => $_clearField(181686379);

  @$pb.TagNumber(224443741)
  $pb.PbList<DimensionMapping> get dimensionmappings => $_getList(2);

  @$pb.TagNumber(501736394)
  MultiMeasureMappings get multimeasuremappings => $_getN(3);
  @$pb.TagNumber(501736394)
  set multimeasuremappings(MultiMeasureMappings value) =>
      $_setField(501736394, value);
  @$pb.TagNumber(501736394)
  $core.bool hasMultimeasuremappings() => $_has(3);
  @$pb.TagNumber(501736394)
  void clearMultimeasuremappings() => $_clearField(501736394);
  @$pb.TagNumber(501736394)
  MultiMeasureMappings ensureMultimeasuremappings() => $_ensure(3);

  @$pb.TagNumber(519449367)
  $core.String get timecolumn => $_getSZ(4);
  @$pb.TagNumber(519449367)
  set timecolumn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(519449367)
  $core.bool hasTimecolumn() => $_has(4);
  @$pb.TagNumber(519449367)
  void clearTimecolumn() => $_clearField(519449367);

  @$pb.TagNumber(521774144)
  $pb.PbList<MixedMeasureMapping> get mixedmeasuremappings => $_getList(5);
}

class DataModelConfiguration extends $pb.GeneratedMessage {
  factory DataModelConfiguration({
    DataModelS3Configuration? datamodels3configuration,
    DataModel? datamodel,
  }) {
    final result = create();
    if (datamodels3configuration != null)
      result.datamodels3configuration = datamodels3configuration;
    if (datamodel != null) result.datamodel = datamodel;
    return result;
  }

  DataModelConfiguration._();

  factory DataModelConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataModelConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataModelConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<DataModelS3Configuration>(
        175841835, _omitFieldNames ? '' : 'datamodels3configuration',
        subBuilder: DataModelS3Configuration.create)
    ..aOM<DataModel>(428899579, _omitFieldNames ? '' : 'datamodel',
        subBuilder: DataModel.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataModelConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataModelConfiguration copyWith(
          void Function(DataModelConfiguration) updates) =>
      super.copyWith((message) => updates(message as DataModelConfiguration))
          as DataModelConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataModelConfiguration create() => DataModelConfiguration._();
  @$core.override
  DataModelConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataModelConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataModelConfiguration>(create);
  static DataModelConfiguration? _defaultInstance;

  @$pb.TagNumber(175841835)
  DataModelS3Configuration get datamodels3configuration => $_getN(0);
  @$pb.TagNumber(175841835)
  set datamodels3configuration(DataModelS3Configuration value) =>
      $_setField(175841835, value);
  @$pb.TagNumber(175841835)
  $core.bool hasDatamodels3configuration() => $_has(0);
  @$pb.TagNumber(175841835)
  void clearDatamodels3configuration() => $_clearField(175841835);
  @$pb.TagNumber(175841835)
  DataModelS3Configuration ensureDatamodels3configuration() => $_ensure(0);

  @$pb.TagNumber(428899579)
  DataModel get datamodel => $_getN(1);
  @$pb.TagNumber(428899579)
  set datamodel(DataModel value) => $_setField(428899579, value);
  @$pb.TagNumber(428899579)
  $core.bool hasDatamodel() => $_has(1);
  @$pb.TagNumber(428899579)
  void clearDatamodel() => $_clearField(428899579);
  @$pb.TagNumber(428899579)
  DataModel ensureDatamodel() => $_ensure(1);
}

class DataModelS3Configuration extends $pb.GeneratedMessage {
  factory DataModelS3Configuration({
    $core.String? bucketname,
    $core.String? objectkey,
  }) {
    final result = create();
    if (bucketname != null) result.bucketname = bucketname;
    if (objectkey != null) result.objectkey = objectkey;
    return result;
  }

  DataModelS3Configuration._();

  factory DataModelS3Configuration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataModelS3Configuration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataModelS3Configuration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(208117045, _omitFieldNames ? '' : 'bucketname')
    ..aOS(335986226, _omitFieldNames ? '' : 'objectkey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataModelS3Configuration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataModelS3Configuration copyWith(
          void Function(DataModelS3Configuration) updates) =>
      super.copyWith((message) => updates(message as DataModelS3Configuration))
          as DataModelS3Configuration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataModelS3Configuration create() => DataModelS3Configuration._();
  @$core.override
  DataModelS3Configuration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataModelS3Configuration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataModelS3Configuration>(create);
  static DataModelS3Configuration? _defaultInstance;

  @$pb.TagNumber(208117045)
  $core.String get bucketname => $_getSZ(0);
  @$pb.TagNumber(208117045)
  set bucketname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(208117045)
  $core.bool hasBucketname() => $_has(0);
  @$pb.TagNumber(208117045)
  void clearBucketname() => $_clearField(208117045);

  @$pb.TagNumber(335986226)
  $core.String get objectkey => $_getSZ(1);
  @$pb.TagNumber(335986226)
  set objectkey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(335986226)
  $core.bool hasObjectkey() => $_has(1);
  @$pb.TagNumber(335986226)
  void clearObjectkey() => $_clearField(335986226);
}

class DataSourceConfiguration extends $pb.GeneratedMessage {
  factory DataSourceConfiguration({
    BatchLoadDataFormat? dataformat,
    CsvConfiguration? csvconfiguration,
    DataSourceS3Configuration? datasources3configuration,
  }) {
    final result = create();
    if (dataformat != null) result.dataformat = dataformat;
    if (csvconfiguration != null) result.csvconfiguration = csvconfiguration;
    if (datasources3configuration != null)
      result.datasources3configuration = datasources3configuration;
    return result;
  }

  DataSourceConfiguration._();

  factory DataSourceConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataSourceConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataSourceConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aE<BatchLoadDataFormat>(89652083, _omitFieldNames ? '' : 'dataformat',
        enumValues: BatchLoadDataFormat.values)
    ..aOM<CsvConfiguration>(
        264008448, _omitFieldNames ? '' : 'csvconfiguration',
        subBuilder: CsvConfiguration.create)
    ..aOM<DataSourceS3Configuration>(
        345023213, _omitFieldNames ? '' : 'datasources3configuration',
        subBuilder: DataSourceS3Configuration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataSourceConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataSourceConfiguration copyWith(
          void Function(DataSourceConfiguration) updates) =>
      super.copyWith((message) => updates(message as DataSourceConfiguration))
          as DataSourceConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataSourceConfiguration create() => DataSourceConfiguration._();
  @$core.override
  DataSourceConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataSourceConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataSourceConfiguration>(create);
  static DataSourceConfiguration? _defaultInstance;

  @$pb.TagNumber(89652083)
  BatchLoadDataFormat get dataformat => $_getN(0);
  @$pb.TagNumber(89652083)
  set dataformat(BatchLoadDataFormat value) => $_setField(89652083, value);
  @$pb.TagNumber(89652083)
  $core.bool hasDataformat() => $_has(0);
  @$pb.TagNumber(89652083)
  void clearDataformat() => $_clearField(89652083);

  @$pb.TagNumber(264008448)
  CsvConfiguration get csvconfiguration => $_getN(1);
  @$pb.TagNumber(264008448)
  set csvconfiguration(CsvConfiguration value) => $_setField(264008448, value);
  @$pb.TagNumber(264008448)
  $core.bool hasCsvconfiguration() => $_has(1);
  @$pb.TagNumber(264008448)
  void clearCsvconfiguration() => $_clearField(264008448);
  @$pb.TagNumber(264008448)
  CsvConfiguration ensureCsvconfiguration() => $_ensure(1);

  @$pb.TagNumber(345023213)
  DataSourceS3Configuration get datasources3configuration => $_getN(2);
  @$pb.TagNumber(345023213)
  set datasources3configuration(DataSourceS3Configuration value) =>
      $_setField(345023213, value);
  @$pb.TagNumber(345023213)
  $core.bool hasDatasources3configuration() => $_has(2);
  @$pb.TagNumber(345023213)
  void clearDatasources3configuration() => $_clearField(345023213);
  @$pb.TagNumber(345023213)
  DataSourceS3Configuration ensureDatasources3configuration() => $_ensure(2);
}

class DataSourceS3Configuration extends $pb.GeneratedMessage {
  factory DataSourceS3Configuration({
    $core.String? objectkeyprefix,
    $core.String? bucketname,
  }) {
    final result = create();
    if (objectkeyprefix != null) result.objectkeyprefix = objectkeyprefix;
    if (bucketname != null) result.bucketname = bucketname;
    return result;
  }

  DataSourceS3Configuration._();

  factory DataSourceS3Configuration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataSourceS3Configuration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataSourceS3Configuration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(132617574, _omitFieldNames ? '' : 'objectkeyprefix')
    ..aOS(208117045, _omitFieldNames ? '' : 'bucketname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataSourceS3Configuration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataSourceS3Configuration copyWith(
          void Function(DataSourceS3Configuration) updates) =>
      super.copyWith((message) => updates(message as DataSourceS3Configuration))
          as DataSourceS3Configuration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataSourceS3Configuration create() => DataSourceS3Configuration._();
  @$core.override
  DataSourceS3Configuration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataSourceS3Configuration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataSourceS3Configuration>(create);
  static DataSourceS3Configuration? _defaultInstance;

  @$pb.TagNumber(132617574)
  $core.String get objectkeyprefix => $_getSZ(0);
  @$pb.TagNumber(132617574)
  set objectkeyprefix($core.String value) => $_setString(0, value);
  @$pb.TagNumber(132617574)
  $core.bool hasObjectkeyprefix() => $_has(0);
  @$pb.TagNumber(132617574)
  void clearObjectkeyprefix() => $_clearField(132617574);

  @$pb.TagNumber(208117045)
  $core.String get bucketname => $_getSZ(1);
  @$pb.TagNumber(208117045)
  set bucketname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(208117045)
  $core.bool hasBucketname() => $_has(1);
  @$pb.TagNumber(208117045)
  void clearBucketname() => $_clearField(208117045);
}

class Database extends $pb.GeneratedMessage {
  factory Database({
    $core.String? kmskeyid,
    $core.String? databasename,
    $core.String? creationtime,
    $core.String? lastupdatedtime,
    $fixnum.Int64? tablecount,
    $core.String? arn,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (databasename != null) result.databasename = databasename;
    if (creationtime != null) result.creationtime = creationtime;
    if (lastupdatedtime != null) result.lastupdatedtime = lastupdatedtime;
    if (tablecount != null) result.tablecount = tablecount;
    if (arn != null) result.arn = arn;
    return result;
  }

  Database._();

  factory Database.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Database.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Database',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(177046166, _omitFieldNames ? '' : 'lastupdatedtime')
    ..aInt64(323868967, _omitFieldNames ? '' : 'tablecount')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Database clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Database copyWith(void Function(Database) updates) =>
      super.copyWith((message) => updates(message as Database)) as Database;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Database create() => Database._();
  @$core.override
  Database createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Database getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Database>(create);
  static Database? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(1);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(1);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(177046166)
  $core.String get lastupdatedtime => $_getSZ(3);
  @$pb.TagNumber(177046166)
  set lastupdatedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(177046166)
  $core.bool hasLastupdatedtime() => $_has(3);
  @$pb.TagNumber(177046166)
  void clearLastupdatedtime() => $_clearField(177046166);

  @$pb.TagNumber(323868967)
  $fixnum.Int64 get tablecount => $_getI64(4);
  @$pb.TagNumber(323868967)
  set tablecount($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(323868967)
  $core.bool hasTablecount() => $_has(4);
  @$pb.TagNumber(323868967)
  void clearTablecount() => $_clearField(323868967);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(5);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(5);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class DeleteDatabaseRequest extends $pb.GeneratedMessage {
  factory DeleteDatabaseRequest({
    $core.String? databasename,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    return result;
  }

  DeleteDatabaseRequest._();

  factory DeleteDatabaseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDatabaseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDatabaseRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDatabaseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDatabaseRequest copyWith(
          void Function(DeleteDatabaseRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteDatabaseRequest))
          as DeleteDatabaseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDatabaseRequest create() => DeleteDatabaseRequest._();
  @$core.override
  DeleteDatabaseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDatabaseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDatabaseRequest>(create);
  static DeleteDatabaseRequest? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);
}

class DeleteTableRequest extends $pb.GeneratedMessage {
  factory DeleteTableRequest({
    $core.String? databasename,
    $core.String? tablename,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DeleteTableRequest._();

  factory DeleteTableRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTableRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTableRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTableRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTableRequest copyWith(void Function(DeleteTableRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteTableRequest))
          as DeleteTableRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTableRequest create() => DeleteTableRequest._();
  @$core.override
  DeleteTableRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTableRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTableRequest>(create);
  static DeleteTableRequest? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeBatchLoadTaskRequest extends $pb.GeneratedMessage {
  factory DescribeBatchLoadTaskRequest({
    $core.String? taskid,
  }) {
    final result = create();
    if (taskid != null) result.taskid = taskid;
    return result;
  }

  DescribeBatchLoadTaskRequest._();

  factory DescribeBatchLoadTaskRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeBatchLoadTaskRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeBatchLoadTaskRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(18532514, _omitFieldNames ? '' : 'taskid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeBatchLoadTaskRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeBatchLoadTaskRequest copyWith(
          void Function(DescribeBatchLoadTaskRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeBatchLoadTaskRequest))
          as DescribeBatchLoadTaskRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeBatchLoadTaskRequest create() =>
      DescribeBatchLoadTaskRequest._();
  @$core.override
  DescribeBatchLoadTaskRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeBatchLoadTaskRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeBatchLoadTaskRequest>(create);
  static DescribeBatchLoadTaskRequest? _defaultInstance;

  @$pb.TagNumber(18532514)
  $core.String get taskid => $_getSZ(0);
  @$pb.TagNumber(18532514)
  set taskid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18532514)
  $core.bool hasTaskid() => $_has(0);
  @$pb.TagNumber(18532514)
  void clearTaskid() => $_clearField(18532514);
}

class DescribeBatchLoadTaskResponse extends $pb.GeneratedMessage {
  factory DescribeBatchLoadTaskResponse({
    BatchLoadTaskDescription? batchloadtaskdescription,
  }) {
    final result = create();
    if (batchloadtaskdescription != null)
      result.batchloadtaskdescription = batchloadtaskdescription;
    return result;
  }

  DescribeBatchLoadTaskResponse._();

  factory DescribeBatchLoadTaskResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeBatchLoadTaskResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeBatchLoadTaskResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<BatchLoadTaskDescription>(
        368256877, _omitFieldNames ? '' : 'batchloadtaskdescription',
        subBuilder: BatchLoadTaskDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeBatchLoadTaskResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeBatchLoadTaskResponse copyWith(
          void Function(DescribeBatchLoadTaskResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeBatchLoadTaskResponse))
          as DescribeBatchLoadTaskResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeBatchLoadTaskResponse create() =>
      DescribeBatchLoadTaskResponse._();
  @$core.override
  DescribeBatchLoadTaskResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeBatchLoadTaskResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeBatchLoadTaskResponse>(create);
  static DescribeBatchLoadTaskResponse? _defaultInstance;

  @$pb.TagNumber(368256877)
  BatchLoadTaskDescription get batchloadtaskdescription => $_getN(0);
  @$pb.TagNumber(368256877)
  set batchloadtaskdescription(BatchLoadTaskDescription value) =>
      $_setField(368256877, value);
  @$pb.TagNumber(368256877)
  $core.bool hasBatchloadtaskdescription() => $_has(0);
  @$pb.TagNumber(368256877)
  void clearBatchloadtaskdescription() => $_clearField(368256877);
  @$pb.TagNumber(368256877)
  BatchLoadTaskDescription ensureBatchloadtaskdescription() => $_ensure(0);
}

class DescribeDatabaseRequest extends $pb.GeneratedMessage {
  factory DescribeDatabaseRequest({
    $core.String? databasename,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    return result;
  }

  DescribeDatabaseRequest._();

  factory DescribeDatabaseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeDatabaseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeDatabaseRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeDatabaseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeDatabaseRequest copyWith(
          void Function(DescribeDatabaseRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeDatabaseRequest))
          as DescribeDatabaseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeDatabaseRequest create() => DescribeDatabaseRequest._();
  @$core.override
  DescribeDatabaseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeDatabaseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeDatabaseRequest>(create);
  static DescribeDatabaseRequest? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);
}

class DescribeDatabaseResponse extends $pb.GeneratedMessage {
  factory DescribeDatabaseResponse({
    Database? database,
  }) {
    final result = create();
    if (database != null) result.database = database;
    return result;
  }

  DescribeDatabaseResponse._();

  factory DescribeDatabaseResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeDatabaseResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeDatabaseResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<Database>(278147289, _omitFieldNames ? '' : 'database',
        subBuilder: Database.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeDatabaseResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeDatabaseResponse copyWith(
          void Function(DescribeDatabaseResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeDatabaseResponse))
          as DescribeDatabaseResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeDatabaseResponse create() => DescribeDatabaseResponse._();
  @$core.override
  DescribeDatabaseResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeDatabaseResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeDatabaseResponse>(create);
  static DescribeDatabaseResponse? _defaultInstance;

  @$pb.TagNumber(278147289)
  Database get database => $_getN(0);
  @$pb.TagNumber(278147289)
  set database(Database value) => $_setField(278147289, value);
  @$pb.TagNumber(278147289)
  $core.bool hasDatabase() => $_has(0);
  @$pb.TagNumber(278147289)
  void clearDatabase() => $_clearField(278147289);
  @$pb.TagNumber(278147289)
  Database ensureDatabase() => $_ensure(0);
}

class DescribeEndpointsRequest extends $pb.GeneratedMessage {
  factory DescribeEndpointsRequest() => create();

  DescribeEndpointsRequest._();

  factory DescribeEndpointsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeEndpointsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeEndpointsRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEndpointsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEndpointsRequest copyWith(
          void Function(DescribeEndpointsRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeEndpointsRequest))
          as DescribeEndpointsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeEndpointsRequest create() => DescribeEndpointsRequest._();
  @$core.override
  DescribeEndpointsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeEndpointsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeEndpointsRequest>(create);
  static DescribeEndpointsRequest? _defaultInstance;
}

class DescribeEndpointsResponse extends $pb.GeneratedMessage {
  factory DescribeEndpointsResponse({
    $core.Iterable<Endpoint>? endpoints,
  }) {
    final result = create();
    if (endpoints != null) result.endpoints.addAll(endpoints);
    return result;
  }

  DescribeEndpointsResponse._();

  factory DescribeEndpointsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeEndpointsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeEndpointsResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..pPM<Endpoint>(16210494, _omitFieldNames ? '' : 'endpoints',
        subBuilder: Endpoint.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEndpointsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEndpointsResponse copyWith(
          void Function(DescribeEndpointsResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeEndpointsResponse))
          as DescribeEndpointsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeEndpointsResponse create() => DescribeEndpointsResponse._();
  @$core.override
  DescribeEndpointsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeEndpointsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeEndpointsResponse>(create);
  static DescribeEndpointsResponse? _defaultInstance;

  @$pb.TagNumber(16210494)
  $pb.PbList<Endpoint> get endpoints => $_getList(0);
}

class DescribeTableRequest extends $pb.GeneratedMessage {
  factory DescribeTableRequest({
    $core.String? databasename,
    $core.String? tablename,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DescribeTableRequest._();

  factory DescribeTableRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTableRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTableRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableRequest copyWith(void Function(DescribeTableRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeTableRequest))
          as DescribeTableRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTableRequest create() => DescribeTableRequest._();
  @$core.override
  DescribeTableRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTableRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeTableRequest>(create);
  static DescribeTableRequest? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeTableResponse extends $pb.GeneratedMessage {
  factory DescribeTableResponse({
    Table? table,
  }) {
    final result = create();
    if (table != null) result.table = table;
    return result;
  }

  DescribeTableResponse._();

  factory DescribeTableResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTableResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTableResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<Table>(386722688, _omitFieldNames ? '' : 'table',
        subBuilder: Table.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableResponse copyWith(
          void Function(DescribeTableResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeTableResponse))
          as DescribeTableResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTableResponse create() => DescribeTableResponse._();
  @$core.override
  DescribeTableResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTableResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeTableResponse>(create);
  static DescribeTableResponse? _defaultInstance;

  @$pb.TagNumber(386722688)
  Table get table => $_getN(0);
  @$pb.TagNumber(386722688)
  set table(Table value) => $_setField(386722688, value);
  @$pb.TagNumber(386722688)
  $core.bool hasTable() => $_has(0);
  @$pb.TagNumber(386722688)
  void clearTable() => $_clearField(386722688);
  @$pb.TagNumber(386722688)
  Table ensureTable() => $_ensure(0);
}

class Dimension extends $pb.GeneratedMessage {
  factory Dimension({
    $core.String? name,
    DimensionValueType? dimensionvaluetype,
    $core.String? value,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (dimensionvaluetype != null)
      result.dimensionvaluetype = dimensionvaluetype;
    if (value != null) result.value = value;
    return result;
  }

  Dimension._();

  factory Dimension.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Dimension.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Dimension',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<DimensionValueType>(
        267417961, _omitFieldNames ? '' : 'dimensionvaluetype',
        enumValues: DimensionValueType.values)
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Dimension clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Dimension copyWith(void Function(Dimension) updates) =>
      super.copyWith((message) => updates(message as Dimension)) as Dimension;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Dimension create() => Dimension._();
  @$core.override
  Dimension createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Dimension getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Dimension>(create);
  static Dimension? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(267417961)
  DimensionValueType get dimensionvaluetype => $_getN(1);
  @$pb.TagNumber(267417961)
  set dimensionvaluetype(DimensionValueType value) =>
      $_setField(267417961, value);
  @$pb.TagNumber(267417961)
  $core.bool hasDimensionvaluetype() => $_has(1);
  @$pb.TagNumber(267417961)
  void clearDimensionvaluetype() => $_clearField(267417961);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(2);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(2, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(2);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class DimensionMapping extends $pb.GeneratedMessage {
  factory DimensionMapping({
    $core.String? sourcecolumn,
    $core.String? destinationcolumn,
  }) {
    final result = create();
    if (sourcecolumn != null) result.sourcecolumn = sourcecolumn;
    if (destinationcolumn != null) result.destinationcolumn = destinationcolumn;
    return result;
  }

  DimensionMapping._();

  factory DimensionMapping.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DimensionMapping.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DimensionMapping',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(219947651, _omitFieldNames ? '' : 'sourcecolumn')
    ..aOS(338930478, _omitFieldNames ? '' : 'destinationcolumn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DimensionMapping clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DimensionMapping copyWith(void Function(DimensionMapping) updates) =>
      super.copyWith((message) => updates(message as DimensionMapping))
          as DimensionMapping;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DimensionMapping create() => DimensionMapping._();
  @$core.override
  DimensionMapping createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DimensionMapping getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DimensionMapping>(create);
  static DimensionMapping? _defaultInstance;

  @$pb.TagNumber(219947651)
  $core.String get sourcecolumn => $_getSZ(0);
  @$pb.TagNumber(219947651)
  set sourcecolumn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(219947651)
  $core.bool hasSourcecolumn() => $_has(0);
  @$pb.TagNumber(219947651)
  void clearSourcecolumn() => $_clearField(219947651);

  @$pb.TagNumber(338930478)
  $core.String get destinationcolumn => $_getSZ(1);
  @$pb.TagNumber(338930478)
  set destinationcolumn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(338930478)
  $core.bool hasDestinationcolumn() => $_has(1);
  @$pb.TagNumber(338930478)
  void clearDestinationcolumn() => $_clearField(338930478);
}

class Endpoint extends $pb.GeneratedMessage {
  factory Endpoint({
    $fixnum.Int64? cacheperiodinminutes,
    $core.String? address,
  }) {
    final result = create();
    if (cacheperiodinminutes != null)
      result.cacheperiodinminutes = cacheperiodinminutes;
    if (address != null) result.address = address;
    return result;
  }

  Endpoint._();

  factory Endpoint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Endpoint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Endpoint',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aInt64(71526449, _omitFieldNames ? '' : 'cacheperiodinminutes')
    ..aOS(268787956, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Endpoint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Endpoint copyWith(void Function(Endpoint) updates) =>
      super.copyWith((message) => updates(message as Endpoint)) as Endpoint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Endpoint create() => Endpoint._();
  @$core.override
  Endpoint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Endpoint getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Endpoint>(create);
  static Endpoint? _defaultInstance;

  @$pb.TagNumber(71526449)
  $fixnum.Int64 get cacheperiodinminutes => $_getI64(0);
  @$pb.TagNumber(71526449)
  set cacheperiodinminutes($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(71526449)
  $core.bool hasCacheperiodinminutes() => $_has(0);
  @$pb.TagNumber(71526449)
  void clearCacheperiodinminutes() => $_clearField(71526449);

  @$pb.TagNumber(268787956)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(268787956)
  set address($core.String value) => $_setString(1, value);
  @$pb.TagNumber(268787956)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(268787956)
  void clearAddress() => $_clearField(268787956);
}

class InternalServerException extends $pb.GeneratedMessage {
  factory InternalServerException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InternalServerException._();

  factory InternalServerException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InternalServerException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InternalServerException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServerException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServerException copyWith(
          void Function(InternalServerException) updates) =>
      super.copyWith((message) => updates(message as InternalServerException))
          as InternalServerException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InternalServerException create() => InternalServerException._();
  @$core.override
  InternalServerException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InternalServerException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InternalServerException>(create);
  static InternalServerException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidEndpointException extends $pb.GeneratedMessage {
  factory InvalidEndpointException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidEndpointException._();

  factory InvalidEndpointException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidEndpointException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidEndpointException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEndpointException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEndpointException copyWith(
          void Function(InvalidEndpointException) updates) =>
      super.copyWith((message) => updates(message as InvalidEndpointException))
          as InvalidEndpointException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidEndpointException create() => InvalidEndpointException._();
  @$core.override
  InvalidEndpointException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidEndpointException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidEndpointException>(create);
  static InvalidEndpointException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ListBatchLoadTasksRequest extends $pb.GeneratedMessage {
  factory ListBatchLoadTasksRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
    BatchLoadStatus? taskstatus,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (taskstatus != null) result.taskstatus = taskstatus;
    return result;
  }

  ListBatchLoadTasksRequest._();

  factory ListBatchLoadTasksRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListBatchLoadTasksRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListBatchLoadTasksRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aE<BatchLoadStatus>(448718071, _omitFieldNames ? '' : 'taskstatus',
        enumValues: BatchLoadStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListBatchLoadTasksRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListBatchLoadTasksRequest copyWith(
          void Function(ListBatchLoadTasksRequest) updates) =>
      super.copyWith((message) => updates(message as ListBatchLoadTasksRequest))
          as ListBatchLoadTasksRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListBatchLoadTasksRequest create() => ListBatchLoadTasksRequest._();
  @$core.override
  ListBatchLoadTasksRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListBatchLoadTasksRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListBatchLoadTasksRequest>(create);
  static ListBatchLoadTasksRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(1);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(448718071)
  BatchLoadStatus get taskstatus => $_getN(2);
  @$pb.TagNumber(448718071)
  set taskstatus(BatchLoadStatus value) => $_setField(448718071, value);
  @$pb.TagNumber(448718071)
  $core.bool hasTaskstatus() => $_has(2);
  @$pb.TagNumber(448718071)
  void clearTaskstatus() => $_clearField(448718071);
}

class ListBatchLoadTasksResponse extends $pb.GeneratedMessage {
  factory ListBatchLoadTasksResponse({
    $core.String? nexttoken,
    $core.Iterable<BatchLoadTask>? batchloadtasks,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (batchloadtasks != null) result.batchloadtasks.addAll(batchloadtasks);
    return result;
  }

  ListBatchLoadTasksResponse._();

  factory ListBatchLoadTasksResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListBatchLoadTasksResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListBatchLoadTasksResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<BatchLoadTask>(292572092, _omitFieldNames ? '' : 'batchloadtasks',
        subBuilder: BatchLoadTask.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListBatchLoadTasksResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListBatchLoadTasksResponse copyWith(
          void Function(ListBatchLoadTasksResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListBatchLoadTasksResponse))
          as ListBatchLoadTasksResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListBatchLoadTasksResponse create() => ListBatchLoadTasksResponse._();
  @$core.override
  ListBatchLoadTasksResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListBatchLoadTasksResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListBatchLoadTasksResponse>(create);
  static ListBatchLoadTasksResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(292572092)
  $pb.PbList<BatchLoadTask> get batchloadtasks => $_getList(1);
}

class ListDatabasesRequest extends $pb.GeneratedMessage {
  factory ListDatabasesRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListDatabasesRequest._();

  factory ListDatabasesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDatabasesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDatabasesRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDatabasesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDatabasesRequest copyWith(void Function(ListDatabasesRequest) updates) =>
      super.copyWith((message) => updates(message as ListDatabasesRequest))
          as ListDatabasesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDatabasesRequest create() => ListDatabasesRequest._();
  @$core.override
  ListDatabasesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDatabasesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDatabasesRequest>(create);
  static ListDatabasesRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(1);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListDatabasesResponse extends $pb.GeneratedMessage {
  factory ListDatabasesResponse({
    $core.Iterable<Database>? databases,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (databases != null) result.databases.addAll(databases);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListDatabasesResponse._();

  factory ListDatabasesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDatabasesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDatabasesResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..pPM<Database>(71867698, _omitFieldNames ? '' : 'databases',
        subBuilder: Database.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDatabasesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDatabasesResponse copyWith(
          void Function(ListDatabasesResponse) updates) =>
      super.copyWith((message) => updates(message as ListDatabasesResponse))
          as ListDatabasesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDatabasesResponse create() => ListDatabasesResponse._();
  @$core.override
  ListDatabasesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDatabasesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDatabasesResponse>(create);
  static ListDatabasesResponse? _defaultInstance;

  @$pb.TagNumber(71867698)
  $pb.PbList<Database> get databases => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListTablesRequest extends $pb.GeneratedMessage {
  factory ListTablesRequest({
    $core.String? databasename,
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListTablesRequest._();

  factory ListTablesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTablesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTablesRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTablesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTablesRequest copyWith(void Function(ListTablesRequest) updates) =>
      super.copyWith((message) => updates(message as ListTablesRequest))
          as ListTablesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTablesRequest create() => ListTablesRequest._();
  @$core.override
  ListTablesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTablesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTablesRequest>(create);
  static ListTablesRequest? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(2);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListTablesResponse extends $pb.GeneratedMessage {
  factory ListTablesResponse({
    $core.String? nexttoken,
    $core.Iterable<Table>? tables,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (tables != null) result.tables.addAll(tables);
    return result;
  }

  ListTablesResponse._();

  factory ListTablesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTablesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTablesResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Table>(357958629, _omitFieldNames ? '' : 'tables',
        subBuilder: Table.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTablesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTablesResponse copyWith(void Function(ListTablesResponse) updates) =>
      super.copyWith((message) => updates(message as ListTablesResponse))
          as ListTablesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTablesResponse create() => ListTablesResponse._();
  @$core.override
  ListTablesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTablesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTablesResponse>(create);
  static ListTablesResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(357958629)
  $pb.PbList<Table> get tables => $_getList(1);
}

class ListTagsForResourceRequest extends $pb.GeneratedMessage {
  factory ListTagsForResourceRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  ListTagsForResourceRequest._();

  factory ListTagsForResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForResourceRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceRequest copyWith(
          void Function(ListTagsForResourceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListTagsForResourceRequest))
          as ListTagsForResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceRequest create() => ListTagsForResourceRequest._();
  @$core.override
  ListTagsForResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForResourceRequest>(create);
  static ListTagsForResourceRequest? _defaultInstance;

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);
}

class ListTagsForResourceResponse extends $pb.GeneratedMessage {
  factory ListTagsForResourceResponse({
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  ListTagsForResourceResponse._();

  factory ListTagsForResourceResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForResourceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForResourceResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceResponse copyWith(
          void Function(ListTagsForResourceResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListTagsForResourceResponse))
          as ListTagsForResourceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceResponse create() =>
      ListTagsForResourceResponse._();
  @$core.override
  ListTagsForResourceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForResourceResponse>(create);
  static ListTagsForResourceResponse? _defaultInstance;

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(0);
}

class MagneticStoreRejectedDataLocation extends $pb.GeneratedMessage {
  factory MagneticStoreRejectedDataLocation({
    S3Configuration? s3configuration,
  }) {
    final result = create();
    if (s3configuration != null) result.s3configuration = s3configuration;
    return result;
  }

  MagneticStoreRejectedDataLocation._();

  factory MagneticStoreRejectedDataLocation.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MagneticStoreRejectedDataLocation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MagneticStoreRejectedDataLocation',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<S3Configuration>(27828476, _omitFieldNames ? '' : 's3configuration',
        subBuilder: S3Configuration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MagneticStoreRejectedDataLocation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MagneticStoreRejectedDataLocation copyWith(
          void Function(MagneticStoreRejectedDataLocation) updates) =>
      super.copyWith((message) =>
              updates(message as MagneticStoreRejectedDataLocation))
          as MagneticStoreRejectedDataLocation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MagneticStoreRejectedDataLocation create() =>
      MagneticStoreRejectedDataLocation._();
  @$core.override
  MagneticStoreRejectedDataLocation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MagneticStoreRejectedDataLocation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MagneticStoreRejectedDataLocation>(
          create);
  static MagneticStoreRejectedDataLocation? _defaultInstance;

  @$pb.TagNumber(27828476)
  S3Configuration get s3configuration => $_getN(0);
  @$pb.TagNumber(27828476)
  set s3configuration(S3Configuration value) => $_setField(27828476, value);
  @$pb.TagNumber(27828476)
  $core.bool hasS3configuration() => $_has(0);
  @$pb.TagNumber(27828476)
  void clearS3configuration() => $_clearField(27828476);
  @$pb.TagNumber(27828476)
  S3Configuration ensureS3configuration() => $_ensure(0);
}

class MagneticStoreWriteProperties extends $pb.GeneratedMessage {
  factory MagneticStoreWriteProperties({
    MagneticStoreRejectedDataLocation? magneticstorerejecteddatalocation,
    $core.bool? enablemagneticstorewrites,
  }) {
    final result = create();
    if (magneticstorerejecteddatalocation != null)
      result.magneticstorerejecteddatalocation =
          magneticstorerejecteddatalocation;
    if (enablemagneticstorewrites != null)
      result.enablemagneticstorewrites = enablemagneticstorewrites;
    return result;
  }

  MagneticStoreWriteProperties._();

  factory MagneticStoreWriteProperties.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MagneticStoreWriteProperties.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MagneticStoreWriteProperties',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<MagneticStoreRejectedDataLocation>(
        47660222, _omitFieldNames ? '' : 'magneticstorerejecteddatalocation',
        subBuilder: MagneticStoreRejectedDataLocation.create)
    ..aOB(489945306, _omitFieldNames ? '' : 'enablemagneticstorewrites')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MagneticStoreWriteProperties clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MagneticStoreWriteProperties copyWith(
          void Function(MagneticStoreWriteProperties) updates) =>
      super.copyWith(
              (message) => updates(message as MagneticStoreWriteProperties))
          as MagneticStoreWriteProperties;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MagneticStoreWriteProperties create() =>
      MagneticStoreWriteProperties._();
  @$core.override
  MagneticStoreWriteProperties createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MagneticStoreWriteProperties getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MagneticStoreWriteProperties>(create);
  static MagneticStoreWriteProperties? _defaultInstance;

  @$pb.TagNumber(47660222)
  MagneticStoreRejectedDataLocation get magneticstorerejecteddatalocation =>
      $_getN(0);
  @$pb.TagNumber(47660222)
  set magneticstorerejecteddatalocation(
          MagneticStoreRejectedDataLocation value) =>
      $_setField(47660222, value);
  @$pb.TagNumber(47660222)
  $core.bool hasMagneticstorerejecteddatalocation() => $_has(0);
  @$pb.TagNumber(47660222)
  void clearMagneticstorerejecteddatalocation() => $_clearField(47660222);
  @$pb.TagNumber(47660222)
  MagneticStoreRejectedDataLocation ensureMagneticstorerejecteddatalocation() =>
      $_ensure(0);

  @$pb.TagNumber(489945306)
  $core.bool get enablemagneticstorewrites => $_getBF(1);
  @$pb.TagNumber(489945306)
  set enablemagneticstorewrites($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(489945306)
  $core.bool hasEnablemagneticstorewrites() => $_has(1);
  @$pb.TagNumber(489945306)
  void clearEnablemagneticstorewrites() => $_clearField(489945306);
}

class MeasureValue extends $pb.GeneratedMessage {
  factory MeasureValue({
    $core.String? name,
    $core.String? value,
    MeasureValueType? type,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (value != null) result.value = value;
    if (type != null) result.type = type;
    return result;
  }

  MeasureValue._();

  factory MeasureValue.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MeasureValue.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MeasureValue',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..aE<MeasureValueType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: MeasureValueType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MeasureValue clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MeasureValue copyWith(void Function(MeasureValue) updates) =>
      super.copyWith((message) => updates(message as MeasureValue))
          as MeasureValue;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MeasureValue create() => MeasureValue._();
  @$core.override
  MeasureValue createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MeasureValue getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MeasureValue>(create);
  static MeasureValue? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(1, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(290836590)
  MeasureValueType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(MeasureValueType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class MixedMeasureMapping extends $pb.GeneratedMessage {
  factory MixedMeasureMapping({
    $core.String? sourcecolumn,
    $core.Iterable<MultiMeasureAttributeMapping>? multimeasureattributemappings,
    $core.String? measurename,
    MeasureValueType? measurevaluetype,
    $core.String? targetmeasurename,
  }) {
    final result = create();
    if (sourcecolumn != null) result.sourcecolumn = sourcecolumn;
    if (multimeasureattributemappings != null)
      result.multimeasureattributemappings
          .addAll(multimeasureattributemappings);
    if (measurename != null) result.measurename = measurename;
    if (measurevaluetype != null) result.measurevaluetype = measurevaluetype;
    if (targetmeasurename != null) result.targetmeasurename = targetmeasurename;
    return result;
  }

  MixedMeasureMapping._();

  factory MixedMeasureMapping.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MixedMeasureMapping.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MixedMeasureMapping',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(219947651, _omitFieldNames ? '' : 'sourcecolumn')
    ..pPM<MultiMeasureAttributeMapping>(
        311133918, _omitFieldNames ? '' : 'multimeasureattributemappings',
        subBuilder: MultiMeasureAttributeMapping.create)
    ..aOS(426079069, _omitFieldNames ? '' : 'measurename')
    ..aE<MeasureValueType>(466683165, _omitFieldNames ? '' : 'measurevaluetype',
        enumValues: MeasureValueType.values)
    ..aOS(469508316, _omitFieldNames ? '' : 'targetmeasurename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MixedMeasureMapping clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MixedMeasureMapping copyWith(void Function(MixedMeasureMapping) updates) =>
      super.copyWith((message) => updates(message as MixedMeasureMapping))
          as MixedMeasureMapping;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MixedMeasureMapping create() => MixedMeasureMapping._();
  @$core.override
  MixedMeasureMapping createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MixedMeasureMapping getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MixedMeasureMapping>(create);
  static MixedMeasureMapping? _defaultInstance;

  @$pb.TagNumber(219947651)
  $core.String get sourcecolumn => $_getSZ(0);
  @$pb.TagNumber(219947651)
  set sourcecolumn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(219947651)
  $core.bool hasSourcecolumn() => $_has(0);
  @$pb.TagNumber(219947651)
  void clearSourcecolumn() => $_clearField(219947651);

  @$pb.TagNumber(311133918)
  $pb.PbList<MultiMeasureAttributeMapping> get multimeasureattributemappings =>
      $_getList(1);

  @$pb.TagNumber(426079069)
  $core.String get measurename => $_getSZ(2);
  @$pb.TagNumber(426079069)
  set measurename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(426079069)
  $core.bool hasMeasurename() => $_has(2);
  @$pb.TagNumber(426079069)
  void clearMeasurename() => $_clearField(426079069);

  @$pb.TagNumber(466683165)
  MeasureValueType get measurevaluetype => $_getN(3);
  @$pb.TagNumber(466683165)
  set measurevaluetype(MeasureValueType value) => $_setField(466683165, value);
  @$pb.TagNumber(466683165)
  $core.bool hasMeasurevaluetype() => $_has(3);
  @$pb.TagNumber(466683165)
  void clearMeasurevaluetype() => $_clearField(466683165);

  @$pb.TagNumber(469508316)
  $core.String get targetmeasurename => $_getSZ(4);
  @$pb.TagNumber(469508316)
  set targetmeasurename($core.String value) => $_setString(4, value);
  @$pb.TagNumber(469508316)
  $core.bool hasTargetmeasurename() => $_has(4);
  @$pb.TagNumber(469508316)
  void clearTargetmeasurename() => $_clearField(469508316);
}

class MultiMeasureAttributeMapping extends $pb.GeneratedMessage {
  factory MultiMeasureAttributeMapping({
    $core.String? sourcecolumn,
    $core.String? targetmultimeasureattributename,
    ScalarMeasureValueType? measurevaluetype,
  }) {
    final result = create();
    if (sourcecolumn != null) result.sourcecolumn = sourcecolumn;
    if (targetmultimeasureattributename != null)
      result.targetmultimeasureattributename = targetmultimeasureattributename;
    if (measurevaluetype != null) result.measurevaluetype = measurevaluetype;
    return result;
  }

  MultiMeasureAttributeMapping._();

  factory MultiMeasureAttributeMapping.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MultiMeasureAttributeMapping.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MultiMeasureAttributeMapping',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(219947651, _omitFieldNames ? '' : 'sourcecolumn')
    ..aOS(415623663, _omitFieldNames ? '' : 'targetmultimeasureattributename')
    ..aE<ScalarMeasureValueType>(
        466683165, _omitFieldNames ? '' : 'measurevaluetype',
        enumValues: ScalarMeasureValueType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MultiMeasureAttributeMapping clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MultiMeasureAttributeMapping copyWith(
          void Function(MultiMeasureAttributeMapping) updates) =>
      super.copyWith(
              (message) => updates(message as MultiMeasureAttributeMapping))
          as MultiMeasureAttributeMapping;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MultiMeasureAttributeMapping create() =>
      MultiMeasureAttributeMapping._();
  @$core.override
  MultiMeasureAttributeMapping createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MultiMeasureAttributeMapping getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MultiMeasureAttributeMapping>(create);
  static MultiMeasureAttributeMapping? _defaultInstance;

  @$pb.TagNumber(219947651)
  $core.String get sourcecolumn => $_getSZ(0);
  @$pb.TagNumber(219947651)
  set sourcecolumn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(219947651)
  $core.bool hasSourcecolumn() => $_has(0);
  @$pb.TagNumber(219947651)
  void clearSourcecolumn() => $_clearField(219947651);

  @$pb.TagNumber(415623663)
  $core.String get targetmultimeasureattributename => $_getSZ(1);
  @$pb.TagNumber(415623663)
  set targetmultimeasureattributename($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(415623663)
  $core.bool hasTargetmultimeasureattributename() => $_has(1);
  @$pb.TagNumber(415623663)
  void clearTargetmultimeasureattributename() => $_clearField(415623663);

  @$pb.TagNumber(466683165)
  ScalarMeasureValueType get measurevaluetype => $_getN(2);
  @$pb.TagNumber(466683165)
  set measurevaluetype(ScalarMeasureValueType value) =>
      $_setField(466683165, value);
  @$pb.TagNumber(466683165)
  $core.bool hasMeasurevaluetype() => $_has(2);
  @$pb.TagNumber(466683165)
  void clearMeasurevaluetype() => $_clearField(466683165);
}

class MultiMeasureMappings extends $pb.GeneratedMessage {
  factory MultiMeasureMappings({
    $core.String? targetmultimeasurename,
    $core.Iterable<MultiMeasureAttributeMapping>? multimeasureattributemappings,
  }) {
    final result = create();
    if (targetmultimeasurename != null)
      result.targetmultimeasurename = targetmultimeasurename;
    if (multimeasureattributemappings != null)
      result.multimeasureattributemappings
          .addAll(multimeasureattributemappings);
    return result;
  }

  MultiMeasureMappings._();

  factory MultiMeasureMappings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MultiMeasureMappings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MultiMeasureMappings',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(27420053, _omitFieldNames ? '' : 'targetmultimeasurename')
    ..pPM<MultiMeasureAttributeMapping>(
        311133918, _omitFieldNames ? '' : 'multimeasureattributemappings',
        subBuilder: MultiMeasureAttributeMapping.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MultiMeasureMappings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MultiMeasureMappings copyWith(void Function(MultiMeasureMappings) updates) =>
      super.copyWith((message) => updates(message as MultiMeasureMappings))
          as MultiMeasureMappings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MultiMeasureMappings create() => MultiMeasureMappings._();
  @$core.override
  MultiMeasureMappings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MultiMeasureMappings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MultiMeasureMappings>(create);
  static MultiMeasureMappings? _defaultInstance;

  @$pb.TagNumber(27420053)
  $core.String get targetmultimeasurename => $_getSZ(0);
  @$pb.TagNumber(27420053)
  set targetmultimeasurename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(27420053)
  $core.bool hasTargetmultimeasurename() => $_has(0);
  @$pb.TagNumber(27420053)
  void clearTargetmultimeasurename() => $_clearField(27420053);

  @$pb.TagNumber(311133918)
  $pb.PbList<MultiMeasureAttributeMapping> get multimeasureattributemappings =>
      $_getList(1);
}

class PartitionKey extends $pb.GeneratedMessage {
  factory PartitionKey({
    PartitionKeyEnforcementLevel? enforcementinrecord,
    $core.String? name,
    PartitionKeyType? type,
  }) {
    final result = create();
    if (enforcementinrecord != null)
      result.enforcementinrecord = enforcementinrecord;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    return result;
  }

  PartitionKey._();

  factory PartitionKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PartitionKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PartitionKey',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aE<PartitionKeyEnforcementLevel>(
        121680624, _omitFieldNames ? '' : 'enforcementinrecord',
        enumValues: PartitionKeyEnforcementLevel.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<PartitionKeyType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: PartitionKeyType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartitionKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartitionKey copyWith(void Function(PartitionKey) updates) =>
      super.copyWith((message) => updates(message as PartitionKey))
          as PartitionKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PartitionKey create() => PartitionKey._();
  @$core.override
  PartitionKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PartitionKey getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PartitionKey>(create);
  static PartitionKey? _defaultInstance;

  @$pb.TagNumber(121680624)
  PartitionKeyEnforcementLevel get enforcementinrecord => $_getN(0);
  @$pb.TagNumber(121680624)
  set enforcementinrecord(PartitionKeyEnforcementLevel value) =>
      $_setField(121680624, value);
  @$pb.TagNumber(121680624)
  $core.bool hasEnforcementinrecord() => $_has(0);
  @$pb.TagNumber(121680624)
  void clearEnforcementinrecord() => $_clearField(121680624);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  PartitionKeyType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(PartitionKeyType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class Record extends $pb.GeneratedMessage {
  factory Record({
    $core.Iterable<MeasureValue>? measurevalues,
    TimeUnit? timeunit,
    $core.String? measurevalue,
    $core.String? measurename,
    $core.Iterable<Dimension>? dimensions,
    MeasureValueType? measurevaluetype,
    $fixnum.Int64? version,
    $core.String? time,
  }) {
    final result = create();
    if (measurevalues != null) result.measurevalues.addAll(measurevalues);
    if (timeunit != null) result.timeunit = timeunit;
    if (measurevalue != null) result.measurevalue = measurevalue;
    if (measurename != null) result.measurename = measurename;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    if (measurevaluetype != null) result.measurevaluetype = measurevaluetype;
    if (version != null) result.version = version;
    if (time != null) result.time = time;
    return result;
  }

  Record._();

  factory Record.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Record.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Record',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..pPM<MeasureValue>(126050982, _omitFieldNames ? '' : 'measurevalues',
        subBuilder: MeasureValue.create)
    ..aE<TimeUnit>(181686379, _omitFieldNames ? '' : 'timeunit',
        enumValues: TimeUnit.values)
    ..aOS(407670165, _omitFieldNames ? '' : 'measurevalue')
    ..aOS(426079069, _omitFieldNames ? '' : 'measurename')
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..aE<MeasureValueType>(466683165, _omitFieldNames ? '' : 'measurevaluetype',
        enumValues: MeasureValueType.values)
    ..aInt64(500028728, _omitFieldNames ? '' : 'version')
    ..aOS(535094277, _omitFieldNames ? '' : 'time')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Record clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Record copyWith(void Function(Record) updates) =>
      super.copyWith((message) => updates(message as Record)) as Record;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Record create() => Record._();
  @$core.override
  Record createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Record getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Record>(create);
  static Record? _defaultInstance;

  @$pb.TagNumber(126050982)
  $pb.PbList<MeasureValue> get measurevalues => $_getList(0);

  @$pb.TagNumber(181686379)
  TimeUnit get timeunit => $_getN(1);
  @$pb.TagNumber(181686379)
  set timeunit(TimeUnit value) => $_setField(181686379, value);
  @$pb.TagNumber(181686379)
  $core.bool hasTimeunit() => $_has(1);
  @$pb.TagNumber(181686379)
  void clearTimeunit() => $_clearField(181686379);

  @$pb.TagNumber(407670165)
  $core.String get measurevalue => $_getSZ(2);
  @$pb.TagNumber(407670165)
  set measurevalue($core.String value) => $_setString(2, value);
  @$pb.TagNumber(407670165)
  $core.bool hasMeasurevalue() => $_has(2);
  @$pb.TagNumber(407670165)
  void clearMeasurevalue() => $_clearField(407670165);

  @$pb.TagNumber(426079069)
  $core.String get measurename => $_getSZ(3);
  @$pb.TagNumber(426079069)
  set measurename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(426079069)
  $core.bool hasMeasurename() => $_has(3);
  @$pb.TagNumber(426079069)
  void clearMeasurename() => $_clearField(426079069);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(4);

  @$pb.TagNumber(466683165)
  MeasureValueType get measurevaluetype => $_getN(5);
  @$pb.TagNumber(466683165)
  set measurevaluetype(MeasureValueType value) => $_setField(466683165, value);
  @$pb.TagNumber(466683165)
  $core.bool hasMeasurevaluetype() => $_has(5);
  @$pb.TagNumber(466683165)
  void clearMeasurevaluetype() => $_clearField(466683165);

  @$pb.TagNumber(500028728)
  $fixnum.Int64 get version => $_getI64(6);
  @$pb.TagNumber(500028728)
  set version($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(6);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);

  @$pb.TagNumber(535094277)
  $core.String get time => $_getSZ(7);
  @$pb.TagNumber(535094277)
  set time($core.String value) => $_setString(7, value);
  @$pb.TagNumber(535094277)
  $core.bool hasTime() => $_has(7);
  @$pb.TagNumber(535094277)
  void clearTime() => $_clearField(535094277);
}

class RecordsIngested extends $pb.GeneratedMessage {
  factory RecordsIngested({
    $core.int? magneticstore,
    $core.int? total,
    $core.int? memorystore,
  }) {
    final result = create();
    if (magneticstore != null) result.magneticstore = magneticstore;
    if (total != null) result.total = total;
    if (memorystore != null) result.memorystore = memorystore;
    return result;
  }

  RecordsIngested._();

  factory RecordsIngested.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecordsIngested.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecordsIngested',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aI(200805165, _omitFieldNames ? '' : 'magneticstore')
    ..aI(218525086, _omitFieldNames ? '' : 'total')
    ..aI(323607360, _omitFieldNames ? '' : 'memorystore')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecordsIngested clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecordsIngested copyWith(void Function(RecordsIngested) updates) =>
      super.copyWith((message) => updates(message as RecordsIngested))
          as RecordsIngested;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecordsIngested create() => RecordsIngested._();
  @$core.override
  RecordsIngested createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecordsIngested getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecordsIngested>(create);
  static RecordsIngested? _defaultInstance;

  @$pb.TagNumber(200805165)
  $core.int get magneticstore => $_getIZ(0);
  @$pb.TagNumber(200805165)
  set magneticstore($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(200805165)
  $core.bool hasMagneticstore() => $_has(0);
  @$pb.TagNumber(200805165)
  void clearMagneticstore() => $_clearField(200805165);

  @$pb.TagNumber(218525086)
  $core.int get total => $_getIZ(1);
  @$pb.TagNumber(218525086)
  set total($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(218525086)
  $core.bool hasTotal() => $_has(1);
  @$pb.TagNumber(218525086)
  void clearTotal() => $_clearField(218525086);

  @$pb.TagNumber(323607360)
  $core.int get memorystore => $_getIZ(2);
  @$pb.TagNumber(323607360)
  set memorystore($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(323607360)
  $core.bool hasMemorystore() => $_has(2);
  @$pb.TagNumber(323607360)
  void clearMemorystore() => $_clearField(323607360);
}

class RejectedRecord extends $pb.GeneratedMessage {
  factory RejectedRecord({
    $core.String? reason,
    $core.int? recordindex,
    $fixnum.Int64? existingversion,
  }) {
    final result = create();
    if (reason != null) result.reason = reason;
    if (recordindex != null) result.recordindex = recordindex;
    if (existingversion != null) result.existingversion = existingversion;
    return result;
  }

  RejectedRecord._();

  factory RejectedRecord.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RejectedRecord.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RejectedRecord',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(20005178, _omitFieldNames ? '' : 'reason')
    ..aI(221086889, _omitFieldNames ? '' : 'recordindex')
    ..aInt64(457856209, _omitFieldNames ? '' : 'existingversion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RejectedRecord clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RejectedRecord copyWith(void Function(RejectedRecord) updates) =>
      super.copyWith((message) => updates(message as RejectedRecord))
          as RejectedRecord;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RejectedRecord create() => RejectedRecord._();
  @$core.override
  RejectedRecord createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RejectedRecord getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RejectedRecord>(create);
  static RejectedRecord? _defaultInstance;

  @$pb.TagNumber(20005178)
  $core.String get reason => $_getSZ(0);
  @$pb.TagNumber(20005178)
  set reason($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20005178)
  $core.bool hasReason() => $_has(0);
  @$pb.TagNumber(20005178)
  void clearReason() => $_clearField(20005178);

  @$pb.TagNumber(221086889)
  $core.int get recordindex => $_getIZ(1);
  @$pb.TagNumber(221086889)
  set recordindex($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(221086889)
  $core.bool hasRecordindex() => $_has(1);
  @$pb.TagNumber(221086889)
  void clearRecordindex() => $_clearField(221086889);

  @$pb.TagNumber(457856209)
  $fixnum.Int64 get existingversion => $_getI64(2);
  @$pb.TagNumber(457856209)
  set existingversion($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(457856209)
  $core.bool hasExistingversion() => $_has(2);
  @$pb.TagNumber(457856209)
  void clearExistingversion() => $_clearField(457856209);
}

class RejectedRecordsException extends $pb.GeneratedMessage {
  factory RejectedRecordsException({
    $core.String? message,
    $core.Iterable<RejectedRecord>? rejectedrecords,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (rejectedrecords != null) result.rejectedrecords.addAll(rejectedrecords);
    return result;
  }

  RejectedRecordsException._();

  factory RejectedRecordsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RejectedRecordsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RejectedRecordsException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..pPM<RejectedRecord>(374345864, _omitFieldNames ? '' : 'rejectedrecords',
        subBuilder: RejectedRecord.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RejectedRecordsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RejectedRecordsException copyWith(
          void Function(RejectedRecordsException) updates) =>
      super.copyWith((message) => updates(message as RejectedRecordsException))
          as RejectedRecordsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RejectedRecordsException create() => RejectedRecordsException._();
  @$core.override
  RejectedRecordsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RejectedRecordsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RejectedRecordsException>(create);
  static RejectedRecordsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(374345864)
  $pb.PbList<RejectedRecord> get rejectedrecords => $_getList(1);
}

class ReportConfiguration extends $pb.GeneratedMessage {
  factory ReportConfiguration({
    ReportS3Configuration? reports3configuration,
  }) {
    final result = create();
    if (reports3configuration != null)
      result.reports3configuration = reports3configuration;
    return result;
  }

  ReportConfiguration._();

  factory ReportConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReportConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReportConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<ReportS3Configuration>(
        463435614, _omitFieldNames ? '' : 'reports3configuration',
        subBuilder: ReportS3Configuration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReportConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReportConfiguration copyWith(void Function(ReportConfiguration) updates) =>
      super.copyWith((message) => updates(message as ReportConfiguration))
          as ReportConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReportConfiguration create() => ReportConfiguration._();
  @$core.override
  ReportConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReportConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReportConfiguration>(create);
  static ReportConfiguration? _defaultInstance;

  @$pb.TagNumber(463435614)
  ReportS3Configuration get reports3configuration => $_getN(0);
  @$pb.TagNumber(463435614)
  set reports3configuration(ReportS3Configuration value) =>
      $_setField(463435614, value);
  @$pb.TagNumber(463435614)
  $core.bool hasReports3configuration() => $_has(0);
  @$pb.TagNumber(463435614)
  void clearReports3configuration() => $_clearField(463435614);
  @$pb.TagNumber(463435614)
  ReportS3Configuration ensureReports3configuration() => $_ensure(0);
}

class ReportS3Configuration extends $pb.GeneratedMessage {
  factory ReportS3Configuration({
    $core.String? kmskeyid,
    $core.String? objectkeyprefix,
    S3EncryptionOption? encryptionoption,
    $core.String? bucketname,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (objectkeyprefix != null) result.objectkeyprefix = objectkeyprefix;
    if (encryptionoption != null) result.encryptionoption = encryptionoption;
    if (bucketname != null) result.bucketname = bucketname;
    return result;
  }

  ReportS3Configuration._();

  factory ReportS3Configuration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReportS3Configuration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReportS3Configuration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(132617574, _omitFieldNames ? '' : 'objectkeyprefix')
    ..aE<S3EncryptionOption>(
        160833062, _omitFieldNames ? '' : 'encryptionoption',
        enumValues: S3EncryptionOption.values)
    ..aOS(208117045, _omitFieldNames ? '' : 'bucketname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReportS3Configuration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReportS3Configuration copyWith(
          void Function(ReportS3Configuration) updates) =>
      super.copyWith((message) => updates(message as ReportS3Configuration))
          as ReportS3Configuration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReportS3Configuration create() => ReportS3Configuration._();
  @$core.override
  ReportS3Configuration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReportS3Configuration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReportS3Configuration>(create);
  static ReportS3Configuration? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(132617574)
  $core.String get objectkeyprefix => $_getSZ(1);
  @$pb.TagNumber(132617574)
  set objectkeyprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(132617574)
  $core.bool hasObjectkeyprefix() => $_has(1);
  @$pb.TagNumber(132617574)
  void clearObjectkeyprefix() => $_clearField(132617574);

  @$pb.TagNumber(160833062)
  S3EncryptionOption get encryptionoption => $_getN(2);
  @$pb.TagNumber(160833062)
  set encryptionoption(S3EncryptionOption value) =>
      $_setField(160833062, value);
  @$pb.TagNumber(160833062)
  $core.bool hasEncryptionoption() => $_has(2);
  @$pb.TagNumber(160833062)
  void clearEncryptionoption() => $_clearField(160833062);

  @$pb.TagNumber(208117045)
  $core.String get bucketname => $_getSZ(3);
  @$pb.TagNumber(208117045)
  set bucketname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(208117045)
  $core.bool hasBucketname() => $_has(3);
  @$pb.TagNumber(208117045)
  void clearBucketname() => $_clearField(208117045);
}

class ResourceNotFoundException extends $pb.GeneratedMessage {
  factory ResourceNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourceNotFoundException._();

  factory ResourceNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceNotFoundException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceNotFoundException copyWith(
          void Function(ResourceNotFoundException) updates) =>
      super.copyWith((message) => updates(message as ResourceNotFoundException))
          as ResourceNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceNotFoundException create() => ResourceNotFoundException._();
  @$core.override
  ResourceNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceNotFoundException>(create);
  static ResourceNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ResumeBatchLoadTaskRequest extends $pb.GeneratedMessage {
  factory ResumeBatchLoadTaskRequest({
    $core.String? taskid,
  }) {
    final result = create();
    if (taskid != null) result.taskid = taskid;
    return result;
  }

  ResumeBatchLoadTaskRequest._();

  factory ResumeBatchLoadTaskRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResumeBatchLoadTaskRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResumeBatchLoadTaskRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(18532514, _omitFieldNames ? '' : 'taskid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResumeBatchLoadTaskRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResumeBatchLoadTaskRequest copyWith(
          void Function(ResumeBatchLoadTaskRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ResumeBatchLoadTaskRequest))
          as ResumeBatchLoadTaskRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResumeBatchLoadTaskRequest create() => ResumeBatchLoadTaskRequest._();
  @$core.override
  ResumeBatchLoadTaskRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResumeBatchLoadTaskRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResumeBatchLoadTaskRequest>(create);
  static ResumeBatchLoadTaskRequest? _defaultInstance;

  @$pb.TagNumber(18532514)
  $core.String get taskid => $_getSZ(0);
  @$pb.TagNumber(18532514)
  set taskid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18532514)
  $core.bool hasTaskid() => $_has(0);
  @$pb.TagNumber(18532514)
  void clearTaskid() => $_clearField(18532514);
}

class ResumeBatchLoadTaskResponse extends $pb.GeneratedMessage {
  factory ResumeBatchLoadTaskResponse() => create();

  ResumeBatchLoadTaskResponse._();

  factory ResumeBatchLoadTaskResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResumeBatchLoadTaskResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResumeBatchLoadTaskResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResumeBatchLoadTaskResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResumeBatchLoadTaskResponse copyWith(
          void Function(ResumeBatchLoadTaskResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ResumeBatchLoadTaskResponse))
          as ResumeBatchLoadTaskResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResumeBatchLoadTaskResponse create() =>
      ResumeBatchLoadTaskResponse._();
  @$core.override
  ResumeBatchLoadTaskResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResumeBatchLoadTaskResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResumeBatchLoadTaskResponse>(create);
  static ResumeBatchLoadTaskResponse? _defaultInstance;
}

class RetentionProperties extends $pb.GeneratedMessage {
  factory RetentionProperties({
    $fixnum.Int64? magneticstoreretentionperiodindays,
    $fixnum.Int64? memorystoreretentionperiodinhours,
  }) {
    final result = create();
    if (magneticstoreretentionperiodindays != null)
      result.magneticstoreretentionperiodindays =
          magneticstoreretentionperiodindays;
    if (memorystoreretentionperiodinhours != null)
      result.memorystoreretentionperiodinhours =
          memorystoreretentionperiodinhours;
    return result;
  }

  RetentionProperties._();

  factory RetentionProperties.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RetentionProperties.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RetentionProperties',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aInt64(
        65073472, _omitFieldNames ? '' : 'magneticstoreretentionperiodindays')
    ..aInt64(
        136905329, _omitFieldNames ? '' : 'memorystoreretentionperiodinhours')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RetentionProperties clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RetentionProperties copyWith(void Function(RetentionProperties) updates) =>
      super.copyWith((message) => updates(message as RetentionProperties))
          as RetentionProperties;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RetentionProperties create() => RetentionProperties._();
  @$core.override
  RetentionProperties createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RetentionProperties getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RetentionProperties>(create);
  static RetentionProperties? _defaultInstance;

  @$pb.TagNumber(65073472)
  $fixnum.Int64 get magneticstoreretentionperiodindays => $_getI64(0);
  @$pb.TagNumber(65073472)
  set magneticstoreretentionperiodindays($fixnum.Int64 value) =>
      $_setInt64(0, value);
  @$pb.TagNumber(65073472)
  $core.bool hasMagneticstoreretentionperiodindays() => $_has(0);
  @$pb.TagNumber(65073472)
  void clearMagneticstoreretentionperiodindays() => $_clearField(65073472);

  @$pb.TagNumber(136905329)
  $fixnum.Int64 get memorystoreretentionperiodinhours => $_getI64(1);
  @$pb.TagNumber(136905329)
  set memorystoreretentionperiodinhours($fixnum.Int64 value) =>
      $_setInt64(1, value);
  @$pb.TagNumber(136905329)
  $core.bool hasMemorystoreretentionperiodinhours() => $_has(1);
  @$pb.TagNumber(136905329)
  void clearMemorystoreretentionperiodinhours() => $_clearField(136905329);
}

class S3Configuration extends $pb.GeneratedMessage {
  factory S3Configuration({
    $core.String? kmskeyid,
    $core.String? objectkeyprefix,
    S3EncryptionOption? encryptionoption,
    $core.String? bucketname,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (objectkeyprefix != null) result.objectkeyprefix = objectkeyprefix;
    if (encryptionoption != null) result.encryptionoption = encryptionoption;
    if (bucketname != null) result.bucketname = bucketname;
    return result;
  }

  S3Configuration._();

  factory S3Configuration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory S3Configuration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'S3Configuration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(132617574, _omitFieldNames ? '' : 'objectkeyprefix')
    ..aE<S3EncryptionOption>(
        160833062, _omitFieldNames ? '' : 'encryptionoption',
        enumValues: S3EncryptionOption.values)
    ..aOS(208117045, _omitFieldNames ? '' : 'bucketname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3Configuration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3Configuration copyWith(void Function(S3Configuration) updates) =>
      super.copyWith((message) => updates(message as S3Configuration))
          as S3Configuration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static S3Configuration create() => S3Configuration._();
  @$core.override
  S3Configuration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static S3Configuration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<S3Configuration>(create);
  static S3Configuration? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(132617574)
  $core.String get objectkeyprefix => $_getSZ(1);
  @$pb.TagNumber(132617574)
  set objectkeyprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(132617574)
  $core.bool hasObjectkeyprefix() => $_has(1);
  @$pb.TagNumber(132617574)
  void clearObjectkeyprefix() => $_clearField(132617574);

  @$pb.TagNumber(160833062)
  S3EncryptionOption get encryptionoption => $_getN(2);
  @$pb.TagNumber(160833062)
  set encryptionoption(S3EncryptionOption value) =>
      $_setField(160833062, value);
  @$pb.TagNumber(160833062)
  $core.bool hasEncryptionoption() => $_has(2);
  @$pb.TagNumber(160833062)
  void clearEncryptionoption() => $_clearField(160833062);

  @$pb.TagNumber(208117045)
  $core.String get bucketname => $_getSZ(3);
  @$pb.TagNumber(208117045)
  set bucketname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(208117045)
  $core.bool hasBucketname() => $_has(3);
  @$pb.TagNumber(208117045)
  void clearBucketname() => $_clearField(208117045);
}

class Schema extends $pb.GeneratedMessage {
  factory Schema({
    $core.Iterable<PartitionKey>? compositepartitionkey,
  }) {
    final result = create();
    if (compositepartitionkey != null)
      result.compositepartitionkey.addAll(compositepartitionkey);
    return result;
  }

  Schema._();

  factory Schema.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Schema.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Schema',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..pPM<PartitionKey>(
        338450538, _omitFieldNames ? '' : 'compositepartitionkey',
        subBuilder: PartitionKey.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Schema clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Schema copyWith(void Function(Schema) updates) =>
      super.copyWith((message) => updates(message as Schema)) as Schema;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Schema create() => Schema._();
  @$core.override
  Schema createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Schema getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Schema>(create);
  static Schema? _defaultInstance;

  @$pb.TagNumber(338450538)
  $pb.PbList<PartitionKey> get compositepartitionkey => $_getList(0);
}

class ServiceQuotaExceededException extends $pb.GeneratedMessage {
  factory ServiceQuotaExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ServiceQuotaExceededException._();

  factory ServiceQuotaExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ServiceQuotaExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ServiceQuotaExceededException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServiceQuotaExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServiceQuotaExceededException copyWith(
          void Function(ServiceQuotaExceededException) updates) =>
      super.copyWith(
              (message) => updates(message as ServiceQuotaExceededException))
          as ServiceQuotaExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ServiceQuotaExceededException create() =>
      ServiceQuotaExceededException._();
  @$core.override
  ServiceQuotaExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ServiceQuotaExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ServiceQuotaExceededException>(create);
  static ServiceQuotaExceededException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class Table extends $pb.GeneratedMessage {
  factory Table({
    $core.String? databasename,
    $core.String? creationtime,
    MagneticStoreWriteProperties? magneticstorewriteproperties,
    $core.String? lastupdatedtime,
    TableStatus? tablestatus,
    RetentionProperties? retentionproperties,
    $core.String? tablename,
    $core.String? arn,
    Schema? schema,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (creationtime != null) result.creationtime = creationtime;
    if (magneticstorewriteproperties != null)
      result.magneticstorewriteproperties = magneticstorewriteproperties;
    if (lastupdatedtime != null) result.lastupdatedtime = lastupdatedtime;
    if (tablestatus != null) result.tablestatus = tablestatus;
    if (retentionproperties != null)
      result.retentionproperties = retentionproperties;
    if (tablename != null) result.tablename = tablename;
    if (arn != null) result.arn = arn;
    if (schema != null) result.schema = schema;
    return result;
  }

  Table._();

  factory Table.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Table.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Table',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOM<MagneticStoreWriteProperties>(
        127209171, _omitFieldNames ? '' : 'magneticstorewriteproperties',
        subBuilder: MagneticStoreWriteProperties.create)
    ..aOS(177046166, _omitFieldNames ? '' : 'lastupdatedtime')
    ..aE<TableStatus>(207908810, _omitFieldNames ? '' : 'tablestatus',
        enumValues: TableStatus.values)
    ..aOM<RetentionProperties>(
        242841241, _omitFieldNames ? '' : 'retentionproperties',
        subBuilder: RetentionProperties.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOM<Schema>(412122455, _omitFieldNames ? '' : 'schema',
        subBuilder: Schema.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Table clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Table copyWith(void Function(Table) updates) =>
      super.copyWith((message) => updates(message as Table)) as Table;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Table create() => Table._();
  @$core.override
  Table createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Table getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Table>(create);
  static Table? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(127209171)
  MagneticStoreWriteProperties get magneticstorewriteproperties => $_getN(2);
  @$pb.TagNumber(127209171)
  set magneticstorewriteproperties(MagneticStoreWriteProperties value) =>
      $_setField(127209171, value);
  @$pb.TagNumber(127209171)
  $core.bool hasMagneticstorewriteproperties() => $_has(2);
  @$pb.TagNumber(127209171)
  void clearMagneticstorewriteproperties() => $_clearField(127209171);
  @$pb.TagNumber(127209171)
  MagneticStoreWriteProperties ensureMagneticstorewriteproperties() =>
      $_ensure(2);

  @$pb.TagNumber(177046166)
  $core.String get lastupdatedtime => $_getSZ(3);
  @$pb.TagNumber(177046166)
  set lastupdatedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(177046166)
  $core.bool hasLastupdatedtime() => $_has(3);
  @$pb.TagNumber(177046166)
  void clearLastupdatedtime() => $_clearField(177046166);

  @$pb.TagNumber(207908810)
  TableStatus get tablestatus => $_getN(4);
  @$pb.TagNumber(207908810)
  set tablestatus(TableStatus value) => $_setField(207908810, value);
  @$pb.TagNumber(207908810)
  $core.bool hasTablestatus() => $_has(4);
  @$pb.TagNumber(207908810)
  void clearTablestatus() => $_clearField(207908810);

  @$pb.TagNumber(242841241)
  RetentionProperties get retentionproperties => $_getN(5);
  @$pb.TagNumber(242841241)
  set retentionproperties(RetentionProperties value) =>
      $_setField(242841241, value);
  @$pb.TagNumber(242841241)
  $core.bool hasRetentionproperties() => $_has(5);
  @$pb.TagNumber(242841241)
  void clearRetentionproperties() => $_clearField(242841241);
  @$pb.TagNumber(242841241)
  RetentionProperties ensureRetentionproperties() => $_ensure(5);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(6);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(6, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(6);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(7);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(7, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(7);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(412122455)
  Schema get schema => $_getN(8);
  @$pb.TagNumber(412122455)
  set schema(Schema value) => $_setField(412122455, value);
  @$pb.TagNumber(412122455)
  $core.bool hasSchema() => $_has(8);
  @$pb.TagNumber(412122455)
  void clearSchema() => $_clearField(412122455);
  @$pb.TagNumber(412122455)
  Schema ensureSchema() => $_ensure(8);
}

class Tag extends $pb.GeneratedMessage {
  factory Tag({
    $core.String? key,
    $core.String? value,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    return result;
  }

  Tag._();

  factory Tag.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Tag.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Tag',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Tag clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Tag copyWith(void Function(Tag) updates) =>
      super.copyWith((message) => updates(message as Tag)) as Tag;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Tag create() => Tag._();
  @$core.override
  Tag createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Tag getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Tag>(create);
  static Tag? _defaultInstance;

  @$pb.TagNumber(219859213)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(219859213)
  set key($core.String value) => $_setString(0, value);
  @$pb.TagNumber(219859213)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(219859213)
  void clearKey() => $_clearField(219859213);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(1, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class TagResourceRequest extends $pb.GeneratedMessage {
  factory TagResourceRequest({
    $core.String? resourcearn,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  TagResourceRequest._();

  factory TagResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagResourceRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceRequest copyWith(void Function(TagResourceRequest) updates) =>
      super.copyWith((message) => updates(message as TagResourceRequest))
          as TagResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagResourceRequest create() => TagResourceRequest._();
  @$core.override
  TagResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagResourceRequest>(create);
  static TagResourceRequest? _defaultInstance;

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);
}

class TagResourceResponse extends $pb.GeneratedMessage {
  factory TagResourceResponse() => create();

  TagResourceResponse._();

  factory TagResourceResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagResourceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagResourceResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceResponse copyWith(void Function(TagResourceResponse) updates) =>
      super.copyWith((message) => updates(message as TagResourceResponse))
          as TagResourceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagResourceResponse create() => TagResourceResponse._();
  @$core.override
  TagResourceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagResourceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagResourceResponse>(create);
  static TagResourceResponse? _defaultInstance;
}

class ThrottlingException extends $pb.GeneratedMessage {
  factory ThrottlingException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ThrottlingException._();

  factory ThrottlingException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ThrottlingException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ThrottlingException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ThrottlingException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ThrottlingException copyWith(void Function(ThrottlingException) updates) =>
      super.copyWith((message) => updates(message as ThrottlingException))
          as ThrottlingException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ThrottlingException create() => ThrottlingException._();
  @$core.override
  ThrottlingException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ThrottlingException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ThrottlingException>(create);
  static ThrottlingException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class UntagResourceRequest extends $pb.GeneratedMessage {
  factory UntagResourceRequest({
    $core.Iterable<$core.String>? tagkeys,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  UntagResourceRequest._();

  factory UntagResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UntagResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UntagResourceRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceRequest copyWith(void Function(UntagResourceRequest) updates) =>
      super.copyWith((message) => updates(message as UntagResourceRequest))
          as UntagResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UntagResourceRequest create() => UntagResourceRequest._();
  @$core.override
  UntagResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UntagResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UntagResourceRequest>(create);
  static UntagResourceRequest? _defaultInstance;

  @$pb.TagNumber(320659964)
  $pb.PbList<$core.String> get tagkeys => $_getList(0);

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);
}

class UntagResourceResponse extends $pb.GeneratedMessage {
  factory UntagResourceResponse() => create();

  UntagResourceResponse._();

  factory UntagResourceResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UntagResourceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UntagResourceResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceResponse copyWith(
          void Function(UntagResourceResponse) updates) =>
      super.copyWith((message) => updates(message as UntagResourceResponse))
          as UntagResourceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UntagResourceResponse create() => UntagResourceResponse._();
  @$core.override
  UntagResourceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UntagResourceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UntagResourceResponse>(create);
  static UntagResourceResponse? _defaultInstance;
}

class UpdateDatabaseRequest extends $pb.GeneratedMessage {
  factory UpdateDatabaseRequest({
    $core.String? kmskeyid,
    $core.String? databasename,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (databasename != null) result.databasename = databasename;
    return result;
  }

  UpdateDatabaseRequest._();

  factory UpdateDatabaseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDatabaseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDatabaseRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDatabaseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDatabaseRequest copyWith(
          void Function(UpdateDatabaseRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateDatabaseRequest))
          as UpdateDatabaseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDatabaseRequest create() => UpdateDatabaseRequest._();
  @$core.override
  UpdateDatabaseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDatabaseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDatabaseRequest>(create);
  static UpdateDatabaseRequest? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(1);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(1);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);
}

class UpdateDatabaseResponse extends $pb.GeneratedMessage {
  factory UpdateDatabaseResponse({
    Database? database,
  }) {
    final result = create();
    if (database != null) result.database = database;
    return result;
  }

  UpdateDatabaseResponse._();

  factory UpdateDatabaseResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDatabaseResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDatabaseResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<Database>(278147289, _omitFieldNames ? '' : 'database',
        subBuilder: Database.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDatabaseResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDatabaseResponse copyWith(
          void Function(UpdateDatabaseResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateDatabaseResponse))
          as UpdateDatabaseResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDatabaseResponse create() => UpdateDatabaseResponse._();
  @$core.override
  UpdateDatabaseResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDatabaseResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDatabaseResponse>(create);
  static UpdateDatabaseResponse? _defaultInstance;

  @$pb.TagNumber(278147289)
  Database get database => $_getN(0);
  @$pb.TagNumber(278147289)
  set database(Database value) => $_setField(278147289, value);
  @$pb.TagNumber(278147289)
  $core.bool hasDatabase() => $_has(0);
  @$pb.TagNumber(278147289)
  void clearDatabase() => $_clearField(278147289);
  @$pb.TagNumber(278147289)
  Database ensureDatabase() => $_ensure(0);
}

class UpdateTableRequest extends $pb.GeneratedMessage {
  factory UpdateTableRequest({
    $core.String? databasename,
    MagneticStoreWriteProperties? magneticstorewriteproperties,
    RetentionProperties? retentionproperties,
    $core.String? tablename,
    Schema? schema,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (magneticstorewriteproperties != null)
      result.magneticstorewriteproperties = magneticstorewriteproperties;
    if (retentionproperties != null)
      result.retentionproperties = retentionproperties;
    if (tablename != null) result.tablename = tablename;
    if (schema != null) result.schema = schema;
    return result;
  }

  UpdateTableRequest._();

  factory UpdateTableRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTableRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTableRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOM<MagneticStoreWriteProperties>(
        127209171, _omitFieldNames ? '' : 'magneticstorewriteproperties',
        subBuilder: MagneticStoreWriteProperties.create)
    ..aOM<RetentionProperties>(
        242841241, _omitFieldNames ? '' : 'retentionproperties',
        subBuilder: RetentionProperties.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<Schema>(412122455, _omitFieldNames ? '' : 'schema',
        subBuilder: Schema.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableRequest copyWith(void Function(UpdateTableRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateTableRequest))
          as UpdateTableRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTableRequest create() => UpdateTableRequest._();
  @$core.override
  UpdateTableRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTableRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTableRequest>(create);
  static UpdateTableRequest? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(127209171)
  MagneticStoreWriteProperties get magneticstorewriteproperties => $_getN(1);
  @$pb.TagNumber(127209171)
  set magneticstorewriteproperties(MagneticStoreWriteProperties value) =>
      $_setField(127209171, value);
  @$pb.TagNumber(127209171)
  $core.bool hasMagneticstorewriteproperties() => $_has(1);
  @$pb.TagNumber(127209171)
  void clearMagneticstorewriteproperties() => $_clearField(127209171);
  @$pb.TagNumber(127209171)
  MagneticStoreWriteProperties ensureMagneticstorewriteproperties() =>
      $_ensure(1);

  @$pb.TagNumber(242841241)
  RetentionProperties get retentionproperties => $_getN(2);
  @$pb.TagNumber(242841241)
  set retentionproperties(RetentionProperties value) =>
      $_setField(242841241, value);
  @$pb.TagNumber(242841241)
  $core.bool hasRetentionproperties() => $_has(2);
  @$pb.TagNumber(242841241)
  void clearRetentionproperties() => $_clearField(242841241);
  @$pb.TagNumber(242841241)
  RetentionProperties ensureRetentionproperties() => $_ensure(2);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(412122455)
  Schema get schema => $_getN(4);
  @$pb.TagNumber(412122455)
  set schema(Schema value) => $_setField(412122455, value);
  @$pb.TagNumber(412122455)
  $core.bool hasSchema() => $_has(4);
  @$pb.TagNumber(412122455)
  void clearSchema() => $_clearField(412122455);
  @$pb.TagNumber(412122455)
  Schema ensureSchema() => $_ensure(4);
}

class UpdateTableResponse extends $pb.GeneratedMessage {
  factory UpdateTableResponse({
    Table? table,
  }) {
    final result = create();
    if (table != null) result.table = table;
    return result;
  }

  UpdateTableResponse._();

  factory UpdateTableResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTableResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTableResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<Table>(386722688, _omitFieldNames ? '' : 'table',
        subBuilder: Table.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableResponse copyWith(void Function(UpdateTableResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateTableResponse))
          as UpdateTableResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTableResponse create() => UpdateTableResponse._();
  @$core.override
  UpdateTableResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTableResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTableResponse>(create);
  static UpdateTableResponse? _defaultInstance;

  @$pb.TagNumber(386722688)
  Table get table => $_getN(0);
  @$pb.TagNumber(386722688)
  set table(Table value) => $_setField(386722688, value);
  @$pb.TagNumber(386722688)
  $core.bool hasTable() => $_has(0);
  @$pb.TagNumber(386722688)
  void clearTable() => $_clearField(386722688);
  @$pb.TagNumber(386722688)
  Table ensureTable() => $_ensure(0);
}

class ValidationException extends $pb.GeneratedMessage {
  factory ValidationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ValidationException._();

  factory ValidationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ValidationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ValidationException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidationException copyWith(void Function(ValidationException) updates) =>
      super.copyWith((message) => updates(message as ValidationException))
          as ValidationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidationException create() => ValidationException._();
  @$core.override
  ValidationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ValidationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ValidationException>(create);
  static ValidationException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WriteRecordsRequest extends $pb.GeneratedMessage {
  factory WriteRecordsRequest({
    $core.String? databasename,
    Record? commonattributes,
    $core.String? tablename,
    $core.Iterable<Record>? records,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (commonattributes != null) result.commonattributes = commonattributes;
    if (tablename != null) result.tablename = tablename;
    if (records != null) result.records.addAll(records);
    return result;
  }

  WriteRecordsRequest._();

  factory WriteRecordsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WriteRecordsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WriteRecordsRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOM<Record>(203721350, _omitFieldNames ? '' : 'commonattributes',
        subBuilder: Record.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..pPM<Record>(423557454, _omitFieldNames ? '' : 'records',
        subBuilder: Record.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WriteRecordsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WriteRecordsRequest copyWith(void Function(WriteRecordsRequest) updates) =>
      super.copyWith((message) => updates(message as WriteRecordsRequest))
          as WriteRecordsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WriteRecordsRequest create() => WriteRecordsRequest._();
  @$core.override
  WriteRecordsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WriteRecordsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WriteRecordsRequest>(create);
  static WriteRecordsRequest? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(203721350)
  Record get commonattributes => $_getN(1);
  @$pb.TagNumber(203721350)
  set commonattributes(Record value) => $_setField(203721350, value);
  @$pb.TagNumber(203721350)
  $core.bool hasCommonattributes() => $_has(1);
  @$pb.TagNumber(203721350)
  void clearCommonattributes() => $_clearField(203721350);
  @$pb.TagNumber(203721350)
  Record ensureCommonattributes() => $_ensure(1);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(2);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(2);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(423557454)
  $pb.PbList<Record> get records => $_getList(3);
}

class WriteRecordsResponse extends $pb.GeneratedMessage {
  factory WriteRecordsResponse({
    RecordsIngested? recordsingested,
  }) {
    final result = create();
    if (recordsingested != null) result.recordsingested = recordsingested;
    return result;
  }

  WriteRecordsResponse._();

  factory WriteRecordsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WriteRecordsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WriteRecordsResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'ingest.timestream'),
      createEmptyInstance: create)
    ..aOM<RecordsIngested>(159925469, _omitFieldNames ? '' : 'recordsingested',
        subBuilder: RecordsIngested.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WriteRecordsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WriteRecordsResponse copyWith(void Function(WriteRecordsResponse) updates) =>
      super.copyWith((message) => updates(message as WriteRecordsResponse))
          as WriteRecordsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WriteRecordsResponse create() => WriteRecordsResponse._();
  @$core.override
  WriteRecordsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WriteRecordsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WriteRecordsResponse>(create);
  static WriteRecordsResponse? _defaultInstance;

  @$pb.TagNumber(159925469)
  RecordsIngested get recordsingested => $_getN(0);
  @$pb.TagNumber(159925469)
  set recordsingested(RecordsIngested value) => $_setField(159925469, value);
  @$pb.TagNumber(159925469)
  $core.bool hasRecordsingested() => $_has(0);
  @$pb.TagNumber(159925469)
  void clearRecordsingested() => $_clearField(159925469);
  @$pb.TagNumber(159925469)
  RecordsIngested ensureRecordsingested() => $_ensure(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
