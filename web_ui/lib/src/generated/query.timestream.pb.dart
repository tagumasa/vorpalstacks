// This is a generated file - do not edit.
//
// Generated from query.timestream.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'query.timestream.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'query.timestream.pbenum.dart';

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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class AccountSettingsNotificationConfiguration extends $pb.GeneratedMessage {
  factory AccountSettingsNotificationConfiguration({
    SnsConfiguration? snsconfiguration,
    $core.String? rolearn,
  }) {
    final result = create();
    if (snsconfiguration != null) result.snsconfiguration = snsconfiguration;
    if (rolearn != null) result.rolearn = rolearn;
    return result;
  }

  AccountSettingsNotificationConfiguration._();

  factory AccountSettingsNotificationConfiguration.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccountSettingsNotificationConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountSettingsNotificationConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<SnsConfiguration>(6844682, _omitFieldNames ? '' : 'snsconfiguration',
        subBuilder: SnsConfiguration.create)
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountSettingsNotificationConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountSettingsNotificationConfiguration copyWith(
          void Function(AccountSettingsNotificationConfiguration) updates) =>
      super.copyWith((message) =>
              updates(message as AccountSettingsNotificationConfiguration))
          as AccountSettingsNotificationConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccountSettingsNotificationConfiguration create() =>
      AccountSettingsNotificationConfiguration._();
  @$core.override
  AccountSettingsNotificationConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccountSettingsNotificationConfiguration getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          AccountSettingsNotificationConfiguration>(create);
  static AccountSettingsNotificationConfiguration? _defaultInstance;

  @$pb.TagNumber(6844682)
  SnsConfiguration get snsconfiguration => $_getN(0);
  @$pb.TagNumber(6844682)
  set snsconfiguration(SnsConfiguration value) => $_setField(6844682, value);
  @$pb.TagNumber(6844682)
  $core.bool hasSnsconfiguration() => $_has(0);
  @$pb.TagNumber(6844682)
  void clearSnsconfiguration() => $_clearField(6844682);
  @$pb.TagNumber(6844682)
  SnsConfiguration ensureSnsconfiguration() => $_ensure(0);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(1);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(1);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);
}

class CancelQueryRequest extends $pb.GeneratedMessage {
  factory CancelQueryRequest({
    $core.String? queryid,
  }) {
    final result = create();
    if (queryid != null) result.queryid = queryid;
    return result;
  }

  CancelQueryRequest._();

  factory CancelQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelQueryRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelQueryRequest copyWith(void Function(CancelQueryRequest) updates) =>
      super.copyWith((message) => updates(message as CancelQueryRequest))
          as CancelQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelQueryRequest create() => CancelQueryRequest._();
  @$core.override
  CancelQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelQueryRequest>(create);
  static CancelQueryRequest? _defaultInstance;

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(0);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(0);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);
}

class CancelQueryResponse extends $pb.GeneratedMessage {
  factory CancelQueryResponse({
    $core.String? cancellationmessage,
  }) {
    final result = create();
    if (cancellationmessage != null)
      result.cancellationmessage = cancellationmessage;
    return result;
  }

  CancelQueryResponse._();

  factory CancelQueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelQueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelQueryResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(226613662, _omitFieldNames ? '' : 'cancellationmessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelQueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelQueryResponse copyWith(void Function(CancelQueryResponse) updates) =>
      super.copyWith((message) => updates(message as CancelQueryResponse))
          as CancelQueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelQueryResponse create() => CancelQueryResponse._();
  @$core.override
  CancelQueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelQueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelQueryResponse>(create);
  static CancelQueryResponse? _defaultInstance;

  @$pb.TagNumber(226613662)
  $core.String get cancellationmessage => $_getSZ(0);
  @$pb.TagNumber(226613662)
  set cancellationmessage($core.String value) => $_setString(0, value);
  @$pb.TagNumber(226613662)
  $core.bool hasCancellationmessage() => $_has(0);
  @$pb.TagNumber(226613662)
  void clearCancellationmessage() => $_clearField(226613662);
}

class ColumnInfo extends $pb.GeneratedMessage {
  factory ColumnInfo({
    $core.String? name,
    Type? type,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    return result;
  }

  ColumnInfo._();

  factory ColumnInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ColumnInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ColumnInfo',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<Type>(290836590, _omitFieldNames ? '' : 'type',
        subBuilder: Type.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ColumnInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ColumnInfo copyWith(void Function(ColumnInfo) updates) =>
      super.copyWith((message) => updates(message as ColumnInfo)) as ColumnInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ColumnInfo create() => ColumnInfo._();
  @$core.override
  ColumnInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ColumnInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ColumnInfo>(create);
  static ColumnInfo? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  Type get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(Type value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
  @$pb.TagNumber(290836590)
  Type ensureType() => $_ensure(1);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class CreateScheduledQueryRequest extends $pb.GeneratedMessage {
  factory CreateScheduledQueryRequest({
    $core.String? kmskeyid,
    $core.String? clienttoken,
    ErrorReportConfiguration? errorreportconfiguration,
    $core.String? name,
    $core.String? scheduledqueryexecutionrolearn,
    NotificationConfiguration? notificationconfiguration,
    $core.Iterable<Tag>? tags,
    $core.String? querystring,
    TargetConfiguration? targetconfiguration,
    ScheduleConfiguration? scheduleconfiguration,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (errorreportconfiguration != null)
      result.errorreportconfiguration = errorreportconfiguration;
    if (name != null) result.name = name;
    if (scheduledqueryexecutionrolearn != null)
      result.scheduledqueryexecutionrolearn = scheduledqueryexecutionrolearn;
    if (notificationconfiguration != null)
      result.notificationconfiguration = notificationconfiguration;
    if (tags != null) result.tags.addAll(tags);
    if (querystring != null) result.querystring = querystring;
    if (targetconfiguration != null)
      result.targetconfiguration = targetconfiguration;
    if (scheduleconfiguration != null)
      result.scheduleconfiguration = scheduleconfiguration;
    return result;
  }

  CreateScheduledQueryRequest._();

  factory CreateScheduledQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateScheduledQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateScheduledQueryRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOM<ErrorReportConfiguration>(
        222039776, _omitFieldNames ? '' : 'errorreportconfiguration',
        subBuilder: ErrorReportConfiguration.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(284244182, _omitFieldNames ? '' : 'scheduledqueryexecutionrolearn')
    ..aOM<NotificationConfiguration>(
        290208045, _omitFieldNames ? '' : 'notificationconfiguration',
        subBuilder: NotificationConfiguration.create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..aOM<TargetConfiguration>(
        499845483, _omitFieldNames ? '' : 'targetconfiguration',
        subBuilder: TargetConfiguration.create)
    ..aOM<ScheduleConfiguration>(
        526075047, _omitFieldNames ? '' : 'scheduleconfiguration',
        subBuilder: ScheduleConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduledQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduledQueryRequest copyWith(
          void Function(CreateScheduledQueryRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateScheduledQueryRequest))
          as CreateScheduledQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateScheduledQueryRequest create() =>
      CreateScheduledQueryRequest._();
  @$core.override
  CreateScheduledQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateScheduledQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateScheduledQueryRequest>(create);
  static CreateScheduledQueryRequest? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(1);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(1);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(222039776)
  ErrorReportConfiguration get errorreportconfiguration => $_getN(2);
  @$pb.TagNumber(222039776)
  set errorreportconfiguration(ErrorReportConfiguration value) =>
      $_setField(222039776, value);
  @$pb.TagNumber(222039776)
  $core.bool hasErrorreportconfiguration() => $_has(2);
  @$pb.TagNumber(222039776)
  void clearErrorreportconfiguration() => $_clearField(222039776);
  @$pb.TagNumber(222039776)
  ErrorReportConfiguration ensureErrorreportconfiguration() => $_ensure(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(284244182)
  $core.String get scheduledqueryexecutionrolearn => $_getSZ(4);
  @$pb.TagNumber(284244182)
  set scheduledqueryexecutionrolearn($core.String value) =>
      $_setString(4, value);
  @$pb.TagNumber(284244182)
  $core.bool hasScheduledqueryexecutionrolearn() => $_has(4);
  @$pb.TagNumber(284244182)
  void clearScheduledqueryexecutionrolearn() => $_clearField(284244182);

  @$pb.TagNumber(290208045)
  NotificationConfiguration get notificationconfiguration => $_getN(5);
  @$pb.TagNumber(290208045)
  set notificationconfiguration(NotificationConfiguration value) =>
      $_setField(290208045, value);
  @$pb.TagNumber(290208045)
  $core.bool hasNotificationconfiguration() => $_has(5);
  @$pb.TagNumber(290208045)
  void clearNotificationconfiguration() => $_clearField(290208045);
  @$pb.TagNumber(290208045)
  NotificationConfiguration ensureNotificationconfiguration() => $_ensure(5);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(6);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(7);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(7, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(7);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);

  @$pb.TagNumber(499845483)
  TargetConfiguration get targetconfiguration => $_getN(8);
  @$pb.TagNumber(499845483)
  set targetconfiguration(TargetConfiguration value) =>
      $_setField(499845483, value);
  @$pb.TagNumber(499845483)
  $core.bool hasTargetconfiguration() => $_has(8);
  @$pb.TagNumber(499845483)
  void clearTargetconfiguration() => $_clearField(499845483);
  @$pb.TagNumber(499845483)
  TargetConfiguration ensureTargetconfiguration() => $_ensure(8);

  @$pb.TagNumber(526075047)
  ScheduleConfiguration get scheduleconfiguration => $_getN(9);
  @$pb.TagNumber(526075047)
  set scheduleconfiguration(ScheduleConfiguration value) =>
      $_setField(526075047, value);
  @$pb.TagNumber(526075047)
  $core.bool hasScheduleconfiguration() => $_has(9);
  @$pb.TagNumber(526075047)
  void clearScheduleconfiguration() => $_clearField(526075047);
  @$pb.TagNumber(526075047)
  ScheduleConfiguration ensureScheduleconfiguration() => $_ensure(9);
}

class CreateScheduledQueryResponse extends $pb.GeneratedMessage {
  factory CreateScheduledQueryResponse({
    $core.String? arn,
  }) {
    final result = create();
    if (arn != null) result.arn = arn;
    return result;
  }

  CreateScheduledQueryResponse._();

  factory CreateScheduledQueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateScheduledQueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateScheduledQueryResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduledQueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduledQueryResponse copyWith(
          void Function(CreateScheduledQueryResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateScheduledQueryResponse))
          as CreateScheduledQueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateScheduledQueryResponse create() =>
      CreateScheduledQueryResponse._();
  @$core.override
  CreateScheduledQueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateScheduledQueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateScheduledQueryResponse>(create);
  static CreateScheduledQueryResponse? _defaultInstance;

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(0);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(0);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class Datum extends $pb.GeneratedMessage {
  factory Datum({
    $core.Iterable<Datum>? arrayvalue,
    $core.String? scalarvalue,
    $core.bool? nullvalue,
    $core.Iterable<TimeSeriesDataPoint>? timeseriesvalue,
    Row? rowvalue,
  }) {
    final result = create();
    if (arrayvalue != null) result.arrayvalue.addAll(arrayvalue);
    if (scalarvalue != null) result.scalarvalue = scalarvalue;
    if (nullvalue != null) result.nullvalue = nullvalue;
    if (timeseriesvalue != null) result.timeseriesvalue.addAll(timeseriesvalue);
    if (rowvalue != null) result.rowvalue = rowvalue;
    return result;
  }

  Datum._();

  factory Datum.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Datum.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Datum',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..pPM<Datum>(20393608, _omitFieldNames ? '' : 'arrayvalue',
        subBuilder: Datum.create)
    ..aOS(45700451, _omitFieldNames ? '' : 'scalarvalue')
    ..aOB(440981694, _omitFieldNames ? '' : 'nullvalue')
    ..pPM<TimeSeriesDataPoint>(
        468590991, _omitFieldNames ? '' : 'timeseriesvalue',
        subBuilder: TimeSeriesDataPoint.create)
    ..aOM<Row>(530552345, _omitFieldNames ? '' : 'rowvalue',
        subBuilder: Row.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Datum clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Datum copyWith(void Function(Datum) updates) =>
      super.copyWith((message) => updates(message as Datum)) as Datum;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Datum create() => Datum._();
  @$core.override
  Datum createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Datum getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Datum>(create);
  static Datum? _defaultInstance;

  @$pb.TagNumber(20393608)
  $pb.PbList<Datum> get arrayvalue => $_getList(0);

  @$pb.TagNumber(45700451)
  $core.String get scalarvalue => $_getSZ(1);
  @$pb.TagNumber(45700451)
  set scalarvalue($core.String value) => $_setString(1, value);
  @$pb.TagNumber(45700451)
  $core.bool hasScalarvalue() => $_has(1);
  @$pb.TagNumber(45700451)
  void clearScalarvalue() => $_clearField(45700451);

  @$pb.TagNumber(440981694)
  $core.bool get nullvalue => $_getBF(2);
  @$pb.TagNumber(440981694)
  set nullvalue($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(440981694)
  $core.bool hasNullvalue() => $_has(2);
  @$pb.TagNumber(440981694)
  void clearNullvalue() => $_clearField(440981694);

  @$pb.TagNumber(468590991)
  $pb.PbList<TimeSeriesDataPoint> get timeseriesvalue => $_getList(3);

  @$pb.TagNumber(530552345)
  Row get rowvalue => $_getN(4);
  @$pb.TagNumber(530552345)
  set rowvalue(Row value) => $_setField(530552345, value);
  @$pb.TagNumber(530552345)
  $core.bool hasRowvalue() => $_has(4);
  @$pb.TagNumber(530552345)
  void clearRowvalue() => $_clearField(530552345);
  @$pb.TagNumber(530552345)
  Row ensureRowvalue() => $_ensure(4);
}

class DeleteScheduledQueryRequest extends $pb.GeneratedMessage {
  factory DeleteScheduledQueryRequest({
    $core.String? scheduledqueryarn,
  }) {
    final result = create();
    if (scheduledqueryarn != null) result.scheduledqueryarn = scheduledqueryarn;
    return result;
  }

  DeleteScheduledQueryRequest._();

  factory DeleteScheduledQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteScheduledQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteScheduledQueryRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(234602964, _omitFieldNames ? '' : 'scheduledqueryarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduledQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduledQueryRequest copyWith(
          void Function(DeleteScheduledQueryRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteScheduledQueryRequest))
          as DeleteScheduledQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteScheduledQueryRequest create() =>
      DeleteScheduledQueryRequest._();
  @$core.override
  DeleteScheduledQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteScheduledQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteScheduledQueryRequest>(create);
  static DeleteScheduledQueryRequest? _defaultInstance;

  @$pb.TagNumber(234602964)
  $core.String get scheduledqueryarn => $_getSZ(0);
  @$pb.TagNumber(234602964)
  set scheduledqueryarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234602964)
  $core.bool hasScheduledqueryarn() => $_has(0);
  @$pb.TagNumber(234602964)
  void clearScheduledqueryarn() => $_clearField(234602964);
}

class DescribeAccountSettingsRequest extends $pb.GeneratedMessage {
  factory DescribeAccountSettingsRequest() => create();

  DescribeAccountSettingsRequest._();

  factory DescribeAccountSettingsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAccountSettingsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAccountSettingsRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAccountSettingsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAccountSettingsRequest copyWith(
          void Function(DescribeAccountSettingsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAccountSettingsRequest))
          as DescribeAccountSettingsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAccountSettingsRequest create() =>
      DescribeAccountSettingsRequest._();
  @$core.override
  DescribeAccountSettingsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAccountSettingsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAccountSettingsRequest>(create);
  static DescribeAccountSettingsRequest? _defaultInstance;
}

class DescribeAccountSettingsResponse extends $pb.GeneratedMessage {
  factory DescribeAccountSettingsResponse({
    $core.int? maxquerytcu,
    QueryComputeResponse? querycompute,
    QueryPricingModel? querypricingmodel,
  }) {
    final result = create();
    if (maxquerytcu != null) result.maxquerytcu = maxquerytcu;
    if (querycompute != null) result.querycompute = querycompute;
    if (querypricingmodel != null) result.querypricingmodel = querypricingmodel;
    return result;
  }

  DescribeAccountSettingsResponse._();

  factory DescribeAccountSettingsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAccountSettingsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAccountSettingsResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aI(285295370, _omitFieldNames ? '' : 'maxquerytcu')
    ..aOM<QueryComputeResponse>(
        358277613, _omitFieldNames ? '' : 'querycompute',
        subBuilder: QueryComputeResponse.create)
    ..aE<QueryPricingModel>(
        453472445, _omitFieldNames ? '' : 'querypricingmodel',
        enumValues: QueryPricingModel.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAccountSettingsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAccountSettingsResponse copyWith(
          void Function(DescribeAccountSettingsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAccountSettingsResponse))
          as DescribeAccountSettingsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAccountSettingsResponse create() =>
      DescribeAccountSettingsResponse._();
  @$core.override
  DescribeAccountSettingsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAccountSettingsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAccountSettingsResponse>(
          create);
  static DescribeAccountSettingsResponse? _defaultInstance;

  @$pb.TagNumber(285295370)
  $core.int get maxquerytcu => $_getIZ(0);
  @$pb.TagNumber(285295370)
  set maxquerytcu($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(285295370)
  $core.bool hasMaxquerytcu() => $_has(0);
  @$pb.TagNumber(285295370)
  void clearMaxquerytcu() => $_clearField(285295370);

  @$pb.TagNumber(358277613)
  QueryComputeResponse get querycompute => $_getN(1);
  @$pb.TagNumber(358277613)
  set querycompute(QueryComputeResponse value) => $_setField(358277613, value);
  @$pb.TagNumber(358277613)
  $core.bool hasQuerycompute() => $_has(1);
  @$pb.TagNumber(358277613)
  void clearQuerycompute() => $_clearField(358277613);
  @$pb.TagNumber(358277613)
  QueryComputeResponse ensureQuerycompute() => $_ensure(1);

  @$pb.TagNumber(453472445)
  QueryPricingModel get querypricingmodel => $_getN(2);
  @$pb.TagNumber(453472445)
  set querypricingmodel(QueryPricingModel value) =>
      $_setField(453472445, value);
  @$pb.TagNumber(453472445)
  $core.bool hasQuerypricingmodel() => $_has(2);
  @$pb.TagNumber(453472445)
  void clearQuerypricingmodel() => $_clearField(453472445);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class DescribeScheduledQueryRequest extends $pb.GeneratedMessage {
  factory DescribeScheduledQueryRequest({
    $core.String? scheduledqueryarn,
  }) {
    final result = create();
    if (scheduledqueryarn != null) result.scheduledqueryarn = scheduledqueryarn;
    return result;
  }

  DescribeScheduledQueryRequest._();

  factory DescribeScheduledQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeScheduledQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeScheduledQueryRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(234602964, _omitFieldNames ? '' : 'scheduledqueryarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeScheduledQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeScheduledQueryRequest copyWith(
          void Function(DescribeScheduledQueryRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeScheduledQueryRequest))
          as DescribeScheduledQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeScheduledQueryRequest create() =>
      DescribeScheduledQueryRequest._();
  @$core.override
  DescribeScheduledQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeScheduledQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeScheduledQueryRequest>(create);
  static DescribeScheduledQueryRequest? _defaultInstance;

  @$pb.TagNumber(234602964)
  $core.String get scheduledqueryarn => $_getSZ(0);
  @$pb.TagNumber(234602964)
  set scheduledqueryarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234602964)
  $core.bool hasScheduledqueryarn() => $_has(0);
  @$pb.TagNumber(234602964)
  void clearScheduledqueryarn() => $_clearField(234602964);
}

class DescribeScheduledQueryResponse extends $pb.GeneratedMessage {
  factory DescribeScheduledQueryResponse({
    ScheduledQueryDescription? scheduledquery,
  }) {
    final result = create();
    if (scheduledquery != null) result.scheduledquery = scheduledquery;
    return result;
  }

  DescribeScheduledQueryResponse._();

  factory DescribeScheduledQueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeScheduledQueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeScheduledQueryResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<ScheduledQueryDescription>(
        338643709, _omitFieldNames ? '' : 'scheduledquery',
        subBuilder: ScheduledQueryDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeScheduledQueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeScheduledQueryResponse copyWith(
          void Function(DescribeScheduledQueryResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeScheduledQueryResponse))
          as DescribeScheduledQueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeScheduledQueryResponse create() =>
      DescribeScheduledQueryResponse._();
  @$core.override
  DescribeScheduledQueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeScheduledQueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeScheduledQueryResponse>(create);
  static DescribeScheduledQueryResponse? _defaultInstance;

  @$pb.TagNumber(338643709)
  ScheduledQueryDescription get scheduledquery => $_getN(0);
  @$pb.TagNumber(338643709)
  set scheduledquery(ScheduledQueryDescription value) =>
      $_setField(338643709, value);
  @$pb.TagNumber(338643709)
  $core.bool hasScheduledquery() => $_has(0);
  @$pb.TagNumber(338643709)
  void clearScheduledquery() => $_clearField(338643709);
  @$pb.TagNumber(338643709)
  ScheduledQueryDescription ensureScheduledquery() => $_ensure(0);
}

class DimensionMapping extends $pb.GeneratedMessage {
  factory DimensionMapping({
    $core.String? name,
    DimensionValueType? dimensionvaluetype,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (dimensionvaluetype != null)
      result.dimensionvaluetype = dimensionvaluetype;
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<DimensionValueType>(
        267417961, _omitFieldNames ? '' : 'dimensionvaluetype',
        enumValues: DimensionValueType.values)
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class ErrorReportConfiguration extends $pb.GeneratedMessage {
  factory ErrorReportConfiguration({
    S3Configuration? s3configuration,
  }) {
    final result = create();
    if (s3configuration != null) result.s3configuration = s3configuration;
    return result;
  }

  ErrorReportConfiguration._();

  factory ErrorReportConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ErrorReportConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ErrorReportConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<S3Configuration>(27828476, _omitFieldNames ? '' : 's3configuration',
        subBuilder: S3Configuration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ErrorReportConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ErrorReportConfiguration copyWith(
          void Function(ErrorReportConfiguration) updates) =>
      super.copyWith((message) => updates(message as ErrorReportConfiguration))
          as ErrorReportConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ErrorReportConfiguration create() => ErrorReportConfiguration._();
  @$core.override
  ErrorReportConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ErrorReportConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ErrorReportConfiguration>(create);
  static ErrorReportConfiguration? _defaultInstance;

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

class ErrorReportLocation extends $pb.GeneratedMessage {
  factory ErrorReportLocation({
    S3ReportLocation? s3reportlocation,
  }) {
    final result = create();
    if (s3reportlocation != null) result.s3reportlocation = s3reportlocation;
    return result;
  }

  ErrorReportLocation._();

  factory ErrorReportLocation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ErrorReportLocation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ErrorReportLocation',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<S3ReportLocation>(
        162597615, _omitFieldNames ? '' : 's3reportlocation',
        subBuilder: S3ReportLocation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ErrorReportLocation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ErrorReportLocation copyWith(void Function(ErrorReportLocation) updates) =>
      super.copyWith((message) => updates(message as ErrorReportLocation))
          as ErrorReportLocation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ErrorReportLocation create() => ErrorReportLocation._();
  @$core.override
  ErrorReportLocation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ErrorReportLocation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ErrorReportLocation>(create);
  static ErrorReportLocation? _defaultInstance;

  @$pb.TagNumber(162597615)
  S3ReportLocation get s3reportlocation => $_getN(0);
  @$pb.TagNumber(162597615)
  set s3reportlocation(S3ReportLocation value) => $_setField(162597615, value);
  @$pb.TagNumber(162597615)
  $core.bool hasS3reportlocation() => $_has(0);
  @$pb.TagNumber(162597615)
  void clearS3reportlocation() => $_clearField(162597615);
  @$pb.TagNumber(162597615)
  S3ReportLocation ensureS3reportlocation() => $_ensure(0);
}

class ExecuteScheduledQueryRequest extends $pb.GeneratedMessage {
  factory ExecuteScheduledQueryRequest({
    ScheduledQueryInsights? queryinsights,
    $core.String? clienttoken,
    $core.String? scheduledqueryarn,
    $core.String? invocationtime,
  }) {
    final result = create();
    if (queryinsights != null) result.queryinsights = queryinsights;
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (scheduledqueryarn != null) result.scheduledqueryarn = scheduledqueryarn;
    if (invocationtime != null) result.invocationtime = invocationtime;
    return result;
  }

  ExecuteScheduledQueryRequest._();

  factory ExecuteScheduledQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecuteScheduledQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecuteScheduledQueryRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<ScheduledQueryInsights>(
        105458863, _omitFieldNames ? '' : 'queryinsights',
        subBuilder: ScheduledQueryInsights.create)
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOS(234602964, _omitFieldNames ? '' : 'scheduledqueryarn')
    ..aOS(331845291, _omitFieldNames ? '' : 'invocationtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteScheduledQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteScheduledQueryRequest copyWith(
          void Function(ExecuteScheduledQueryRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ExecuteScheduledQueryRequest))
          as ExecuteScheduledQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecuteScheduledQueryRequest create() =>
      ExecuteScheduledQueryRequest._();
  @$core.override
  ExecuteScheduledQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecuteScheduledQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecuteScheduledQueryRequest>(create);
  static ExecuteScheduledQueryRequest? _defaultInstance;

  @$pb.TagNumber(105458863)
  ScheduledQueryInsights get queryinsights => $_getN(0);
  @$pb.TagNumber(105458863)
  set queryinsights(ScheduledQueryInsights value) =>
      $_setField(105458863, value);
  @$pb.TagNumber(105458863)
  $core.bool hasQueryinsights() => $_has(0);
  @$pb.TagNumber(105458863)
  void clearQueryinsights() => $_clearField(105458863);
  @$pb.TagNumber(105458863)
  ScheduledQueryInsights ensureQueryinsights() => $_ensure(0);

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(1);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(1);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(234602964)
  $core.String get scheduledqueryarn => $_getSZ(2);
  @$pb.TagNumber(234602964)
  set scheduledqueryarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(234602964)
  $core.bool hasScheduledqueryarn() => $_has(2);
  @$pb.TagNumber(234602964)
  void clearScheduledqueryarn() => $_clearField(234602964);

  @$pb.TagNumber(331845291)
  $core.String get invocationtime => $_getSZ(3);
  @$pb.TagNumber(331845291)
  set invocationtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(331845291)
  $core.bool hasInvocationtime() => $_has(3);
  @$pb.TagNumber(331845291)
  void clearInvocationtime() => $_clearField(331845291);
}

class ExecutionStats extends $pb.GeneratedMessage {
  factory ExecutionStats({
    $fixnum.Int64? recordsingested,
    $fixnum.Int64? bytesmetered,
    $fixnum.Int64? queryresultrows,
    $fixnum.Int64? datawrites,
    $fixnum.Int64? cumulativebytesscanned,
    $fixnum.Int64? executiontimeinmillis,
  }) {
    final result = create();
    if (recordsingested != null) result.recordsingested = recordsingested;
    if (bytesmetered != null) result.bytesmetered = bytesmetered;
    if (queryresultrows != null) result.queryresultrows = queryresultrows;
    if (datawrites != null) result.datawrites = datawrites;
    if (cumulativebytesscanned != null)
      result.cumulativebytesscanned = cumulativebytesscanned;
    if (executiontimeinmillis != null)
      result.executiontimeinmillis = executiontimeinmillis;
    return result;
  }

  ExecutionStats._();

  factory ExecutionStats.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionStats.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionStats',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aInt64(159925469, _omitFieldNames ? '' : 'recordsingested')
    ..aInt64(232712245, _omitFieldNames ? '' : 'bytesmetered')
    ..aInt64(240264704, _omitFieldNames ? '' : 'queryresultrows')
    ..aInt64(406476830, _omitFieldNames ? '' : 'datawrites')
    ..aInt64(413524176, _omitFieldNames ? '' : 'cumulativebytesscanned')
    ..aInt64(447596178, _omitFieldNames ? '' : 'executiontimeinmillis')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionStats clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionStats copyWith(void Function(ExecutionStats) updates) =>
      super.copyWith((message) => updates(message as ExecutionStats))
          as ExecutionStats;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionStats create() => ExecutionStats._();
  @$core.override
  ExecutionStats createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionStats getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionStats>(create);
  static ExecutionStats? _defaultInstance;

  @$pb.TagNumber(159925469)
  $fixnum.Int64 get recordsingested => $_getI64(0);
  @$pb.TagNumber(159925469)
  set recordsingested($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(159925469)
  $core.bool hasRecordsingested() => $_has(0);
  @$pb.TagNumber(159925469)
  void clearRecordsingested() => $_clearField(159925469);

  @$pb.TagNumber(232712245)
  $fixnum.Int64 get bytesmetered => $_getI64(1);
  @$pb.TagNumber(232712245)
  set bytesmetered($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(232712245)
  $core.bool hasBytesmetered() => $_has(1);
  @$pb.TagNumber(232712245)
  void clearBytesmetered() => $_clearField(232712245);

  @$pb.TagNumber(240264704)
  $fixnum.Int64 get queryresultrows => $_getI64(2);
  @$pb.TagNumber(240264704)
  set queryresultrows($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(240264704)
  $core.bool hasQueryresultrows() => $_has(2);
  @$pb.TagNumber(240264704)
  void clearQueryresultrows() => $_clearField(240264704);

  @$pb.TagNumber(406476830)
  $fixnum.Int64 get datawrites => $_getI64(3);
  @$pb.TagNumber(406476830)
  set datawrites($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(406476830)
  $core.bool hasDatawrites() => $_has(3);
  @$pb.TagNumber(406476830)
  void clearDatawrites() => $_clearField(406476830);

  @$pb.TagNumber(413524176)
  $fixnum.Int64 get cumulativebytesscanned => $_getI64(4);
  @$pb.TagNumber(413524176)
  set cumulativebytesscanned($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(413524176)
  $core.bool hasCumulativebytesscanned() => $_has(4);
  @$pb.TagNumber(413524176)
  void clearCumulativebytesscanned() => $_clearField(413524176);

  @$pb.TagNumber(447596178)
  $fixnum.Int64 get executiontimeinmillis => $_getI64(5);
  @$pb.TagNumber(447596178)
  set executiontimeinmillis($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(447596178)
  $core.bool hasExecutiontimeinmillis() => $_has(5);
  @$pb.TagNumber(447596178)
  void clearExecutiontimeinmillis() => $_clearField(447596178);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class LastUpdate extends $pb.GeneratedMessage {
  factory LastUpdate({
    LastUpdateStatus? status,
    $core.String? statusmessage,
    $core.int? targetquerytcu,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (statusmessage != null) result.statusmessage = statusmessage;
    if (targetquerytcu != null) result.targetquerytcu = targetquerytcu;
    return result;
  }

  LastUpdate._();

  factory LastUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LastUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LastUpdate',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aE<LastUpdateStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: LastUpdateStatus.values)
    ..aOS(72590095, _omitFieldNames ? '' : 'statusmessage')
    ..aI(183880621, _omitFieldNames ? '' : 'targetquerytcu')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LastUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LastUpdate copyWith(void Function(LastUpdate) updates) =>
      super.copyWith((message) => updates(message as LastUpdate)) as LastUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LastUpdate create() => LastUpdate._();
  @$core.override
  LastUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LastUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LastUpdate>(create);
  static LastUpdate? _defaultInstance;

  @$pb.TagNumber(6222352)
  LastUpdateStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(LastUpdateStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(72590095)
  $core.String get statusmessage => $_getSZ(1);
  @$pb.TagNumber(72590095)
  set statusmessage($core.String value) => $_setString(1, value);
  @$pb.TagNumber(72590095)
  $core.bool hasStatusmessage() => $_has(1);
  @$pb.TagNumber(72590095)
  void clearStatusmessage() => $_clearField(72590095);

  @$pb.TagNumber(183880621)
  $core.int get targetquerytcu => $_getIZ(2);
  @$pb.TagNumber(183880621)
  set targetquerytcu($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(183880621)
  $core.bool hasTargetquerytcu() => $_has(2);
  @$pb.TagNumber(183880621)
  void clearTargetquerytcu() => $_clearField(183880621);
}

class ListScheduledQueriesRequest extends $pb.GeneratedMessage {
  factory ListScheduledQueriesRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListScheduledQueriesRequest._();

  factory ListScheduledQueriesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListScheduledQueriesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListScheduledQueriesRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListScheduledQueriesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListScheduledQueriesRequest copyWith(
          void Function(ListScheduledQueriesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListScheduledQueriesRequest))
          as ListScheduledQueriesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListScheduledQueriesRequest create() =>
      ListScheduledQueriesRequest._();
  @$core.override
  ListScheduledQueriesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListScheduledQueriesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListScheduledQueriesRequest>(create);
  static ListScheduledQueriesRequest? _defaultInstance;

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

class ListScheduledQueriesResponse extends $pb.GeneratedMessage {
  factory ListScheduledQueriesResponse({
    $core.String? nexttoken,
    $core.Iterable<ScheduledQuery>? scheduledqueries,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (scheduledqueries != null)
      result.scheduledqueries.addAll(scheduledqueries);
    return result;
  }

  ListScheduledQueriesResponse._();

  factory ListScheduledQueriesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListScheduledQueriesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListScheduledQueriesResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ScheduledQuery>(458789865, _omitFieldNames ? '' : 'scheduledqueries',
        subBuilder: ScheduledQuery.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListScheduledQueriesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListScheduledQueriesResponse copyWith(
          void Function(ListScheduledQueriesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListScheduledQueriesResponse))
          as ListScheduledQueriesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListScheduledQueriesResponse create() =>
      ListScheduledQueriesResponse._();
  @$core.override
  ListScheduledQueriesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListScheduledQueriesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListScheduledQueriesResponse>(create);
  static ListScheduledQueriesResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(458789865)
  $pb.PbList<ScheduledQuery> get scheduledqueries => $_getList(1);
}

class ListTagsForResourceRequest extends $pb.GeneratedMessage {
  factory ListTagsForResourceRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
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

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(2);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(2);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);
}

class ListTagsForResourceResponse extends $pb.GeneratedMessage {
  factory ListTagsForResourceResponse({
    $core.String? nexttoken,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
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

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class NotificationConfiguration extends $pb.GeneratedMessage {
  factory NotificationConfiguration({
    SnsConfiguration? snsconfiguration,
  }) {
    final result = create();
    if (snsconfiguration != null) result.snsconfiguration = snsconfiguration;
    return result;
  }

  NotificationConfiguration._();

  factory NotificationConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NotificationConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NotificationConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<SnsConfiguration>(6844682, _omitFieldNames ? '' : 'snsconfiguration',
        subBuilder: SnsConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotificationConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotificationConfiguration copyWith(
          void Function(NotificationConfiguration) updates) =>
      super.copyWith((message) => updates(message as NotificationConfiguration))
          as NotificationConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NotificationConfiguration create() => NotificationConfiguration._();
  @$core.override
  NotificationConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NotificationConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NotificationConfiguration>(create);
  static NotificationConfiguration? _defaultInstance;

  @$pb.TagNumber(6844682)
  SnsConfiguration get snsconfiguration => $_getN(0);
  @$pb.TagNumber(6844682)
  set snsconfiguration(SnsConfiguration value) => $_setField(6844682, value);
  @$pb.TagNumber(6844682)
  $core.bool hasSnsconfiguration() => $_has(0);
  @$pb.TagNumber(6844682)
  void clearSnsconfiguration() => $_clearField(6844682);
  @$pb.TagNumber(6844682)
  SnsConfiguration ensureSnsconfiguration() => $_ensure(0);
}

class ParameterMapping extends $pb.GeneratedMessage {
  factory ParameterMapping({
    $core.String? name,
    Type? type,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    return result;
  }

  ParameterMapping._();

  factory ParameterMapping.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ParameterMapping.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ParameterMapping',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<Type>(290836590, _omitFieldNames ? '' : 'type',
        subBuilder: Type.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ParameterMapping clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ParameterMapping copyWith(void Function(ParameterMapping) updates) =>
      super.copyWith((message) => updates(message as ParameterMapping))
          as ParameterMapping;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ParameterMapping create() => ParameterMapping._();
  @$core.override
  ParameterMapping createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ParameterMapping getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ParameterMapping>(create);
  static ParameterMapping? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  Type get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(Type value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
  @$pb.TagNumber(290836590)
  Type ensureType() => $_ensure(1);
}

class PrepareQueryRequest extends $pb.GeneratedMessage {
  factory PrepareQueryRequest({
    $core.String? querystring,
    $core.bool? validateonly,
  }) {
    final result = create();
    if (querystring != null) result.querystring = querystring;
    if (validateonly != null) result.validateonly = validateonly;
    return result;
  }

  PrepareQueryRequest._();

  factory PrepareQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PrepareQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PrepareQueryRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..aOB(453040966, _omitFieldNames ? '' : 'validateonly')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PrepareQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PrepareQueryRequest copyWith(void Function(PrepareQueryRequest) updates) =>
      super.copyWith((message) => updates(message as PrepareQueryRequest))
          as PrepareQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PrepareQueryRequest create() => PrepareQueryRequest._();
  @$core.override
  PrepareQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PrepareQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PrepareQueryRequest>(create);
  static PrepareQueryRequest? _defaultInstance;

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(0);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(0, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(0);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);

  @$pb.TagNumber(453040966)
  $core.bool get validateonly => $_getBF(1);
  @$pb.TagNumber(453040966)
  set validateonly($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(453040966)
  $core.bool hasValidateonly() => $_has(1);
  @$pb.TagNumber(453040966)
  void clearValidateonly() => $_clearField(453040966);
}

class PrepareQueryResponse extends $pb.GeneratedMessage {
  factory PrepareQueryResponse({
    $core.Iterable<SelectColumn>? columns,
    $core.String? querystring,
    $core.Iterable<ParameterMapping>? parameters,
  }) {
    final result = create();
    if (columns != null) result.columns.addAll(columns);
    if (querystring != null) result.querystring = querystring;
    if (parameters != null) result.parameters.addAll(parameters);
    return result;
  }

  PrepareQueryResponse._();

  factory PrepareQueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PrepareQueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PrepareQueryResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..pPM<SelectColumn>(169177053, _omitFieldNames ? '' : 'columns',
        subBuilder: SelectColumn.create)
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..pPM<ParameterMapping>(494900218, _omitFieldNames ? '' : 'parameters',
        subBuilder: ParameterMapping.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PrepareQueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PrepareQueryResponse copyWith(void Function(PrepareQueryResponse) updates) =>
      super.copyWith((message) => updates(message as PrepareQueryResponse))
          as PrepareQueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PrepareQueryResponse create() => PrepareQueryResponse._();
  @$core.override
  PrepareQueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PrepareQueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PrepareQueryResponse>(create);
  static PrepareQueryResponse? _defaultInstance;

  @$pb.TagNumber(169177053)
  $pb.PbList<SelectColumn> get columns => $_getList(0);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(1);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(1, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(1);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);

  @$pb.TagNumber(494900218)
  $pb.PbList<ParameterMapping> get parameters => $_getList(2);
}

class ProvisionedCapacityRequest extends $pb.GeneratedMessage {
  factory ProvisionedCapacityRequest({
    $core.int? targetquerytcu,
    AccountSettingsNotificationConfiguration? notificationconfiguration,
  }) {
    final result = create();
    if (targetquerytcu != null) result.targetquerytcu = targetquerytcu;
    if (notificationconfiguration != null)
      result.notificationconfiguration = notificationconfiguration;
    return result;
  }

  ProvisionedCapacityRequest._();

  factory ProvisionedCapacityRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProvisionedCapacityRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProvisionedCapacityRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aI(183880621, _omitFieldNames ? '' : 'targetquerytcu')
    ..aOM<AccountSettingsNotificationConfiguration>(
        290208045, _omitFieldNames ? '' : 'notificationconfiguration',
        subBuilder: AccountSettingsNotificationConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedCapacityRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedCapacityRequest copyWith(
          void Function(ProvisionedCapacityRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ProvisionedCapacityRequest))
          as ProvisionedCapacityRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProvisionedCapacityRequest create() => ProvisionedCapacityRequest._();
  @$core.override
  ProvisionedCapacityRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProvisionedCapacityRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProvisionedCapacityRequest>(create);
  static ProvisionedCapacityRequest? _defaultInstance;

  @$pb.TagNumber(183880621)
  $core.int get targetquerytcu => $_getIZ(0);
  @$pb.TagNumber(183880621)
  set targetquerytcu($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(183880621)
  $core.bool hasTargetquerytcu() => $_has(0);
  @$pb.TagNumber(183880621)
  void clearTargetquerytcu() => $_clearField(183880621);

  @$pb.TagNumber(290208045)
  AccountSettingsNotificationConfiguration get notificationconfiguration =>
      $_getN(1);
  @$pb.TagNumber(290208045)
  set notificationconfiguration(
          AccountSettingsNotificationConfiguration value) =>
      $_setField(290208045, value);
  @$pb.TagNumber(290208045)
  $core.bool hasNotificationconfiguration() => $_has(1);
  @$pb.TagNumber(290208045)
  void clearNotificationconfiguration() => $_clearField(290208045);
  @$pb.TagNumber(290208045)
  AccountSettingsNotificationConfiguration ensureNotificationconfiguration() =>
      $_ensure(1);
}

class ProvisionedCapacityResponse extends $pb.GeneratedMessage {
  factory ProvisionedCapacityResponse({
    $core.int? activequerytcu,
    AccountSettingsNotificationConfiguration? notificationconfiguration,
    LastUpdate? lastupdate,
  }) {
    final result = create();
    if (activequerytcu != null) result.activequerytcu = activequerytcu;
    if (notificationconfiguration != null)
      result.notificationconfiguration = notificationconfiguration;
    if (lastupdate != null) result.lastupdate = lastupdate;
    return result;
  }

  ProvisionedCapacityResponse._();

  factory ProvisionedCapacityResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProvisionedCapacityResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProvisionedCapacityResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aI(136246688, _omitFieldNames ? '' : 'activequerytcu')
    ..aOM<AccountSettingsNotificationConfiguration>(
        290208045, _omitFieldNames ? '' : 'notificationconfiguration',
        subBuilder: AccountSettingsNotificationConfiguration.create)
    ..aOM<LastUpdate>(331125817, _omitFieldNames ? '' : 'lastupdate',
        subBuilder: LastUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedCapacityResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedCapacityResponse copyWith(
          void Function(ProvisionedCapacityResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ProvisionedCapacityResponse))
          as ProvisionedCapacityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProvisionedCapacityResponse create() =>
      ProvisionedCapacityResponse._();
  @$core.override
  ProvisionedCapacityResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProvisionedCapacityResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProvisionedCapacityResponse>(create);
  static ProvisionedCapacityResponse? _defaultInstance;

  @$pb.TagNumber(136246688)
  $core.int get activequerytcu => $_getIZ(0);
  @$pb.TagNumber(136246688)
  set activequerytcu($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(136246688)
  $core.bool hasActivequerytcu() => $_has(0);
  @$pb.TagNumber(136246688)
  void clearActivequerytcu() => $_clearField(136246688);

  @$pb.TagNumber(290208045)
  AccountSettingsNotificationConfiguration get notificationconfiguration =>
      $_getN(1);
  @$pb.TagNumber(290208045)
  set notificationconfiguration(
          AccountSettingsNotificationConfiguration value) =>
      $_setField(290208045, value);
  @$pb.TagNumber(290208045)
  $core.bool hasNotificationconfiguration() => $_has(1);
  @$pb.TagNumber(290208045)
  void clearNotificationconfiguration() => $_clearField(290208045);
  @$pb.TagNumber(290208045)
  AccountSettingsNotificationConfiguration ensureNotificationconfiguration() =>
      $_ensure(1);

  @$pb.TagNumber(331125817)
  LastUpdate get lastupdate => $_getN(2);
  @$pb.TagNumber(331125817)
  set lastupdate(LastUpdate value) => $_setField(331125817, value);
  @$pb.TagNumber(331125817)
  $core.bool hasLastupdate() => $_has(2);
  @$pb.TagNumber(331125817)
  void clearLastupdate() => $_clearField(331125817);
  @$pb.TagNumber(331125817)
  LastUpdate ensureLastupdate() => $_ensure(2);
}

class QueryComputeRequest extends $pb.GeneratedMessage {
  factory QueryComputeRequest({
    ProvisionedCapacityRequest? provisionedcapacity,
    ComputeMode? computemode,
  }) {
    final result = create();
    if (provisionedcapacity != null)
      result.provisionedcapacity = provisionedcapacity;
    if (computemode != null) result.computemode = computemode;
    return result;
  }

  QueryComputeRequest._();

  factory QueryComputeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryComputeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryComputeRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<ProvisionedCapacityRequest>(
        181400702, _omitFieldNames ? '' : 'provisionedcapacity',
        subBuilder: ProvisionedCapacityRequest.create)
    ..aE<ComputeMode>(205227458, _omitFieldNames ? '' : 'computemode',
        enumValues: ComputeMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryComputeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryComputeRequest copyWith(void Function(QueryComputeRequest) updates) =>
      super.copyWith((message) => updates(message as QueryComputeRequest))
          as QueryComputeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryComputeRequest create() => QueryComputeRequest._();
  @$core.override
  QueryComputeRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryComputeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryComputeRequest>(create);
  static QueryComputeRequest? _defaultInstance;

  @$pb.TagNumber(181400702)
  ProvisionedCapacityRequest get provisionedcapacity => $_getN(0);
  @$pb.TagNumber(181400702)
  set provisionedcapacity(ProvisionedCapacityRequest value) =>
      $_setField(181400702, value);
  @$pb.TagNumber(181400702)
  $core.bool hasProvisionedcapacity() => $_has(0);
  @$pb.TagNumber(181400702)
  void clearProvisionedcapacity() => $_clearField(181400702);
  @$pb.TagNumber(181400702)
  ProvisionedCapacityRequest ensureProvisionedcapacity() => $_ensure(0);

  @$pb.TagNumber(205227458)
  ComputeMode get computemode => $_getN(1);
  @$pb.TagNumber(205227458)
  set computemode(ComputeMode value) => $_setField(205227458, value);
  @$pb.TagNumber(205227458)
  $core.bool hasComputemode() => $_has(1);
  @$pb.TagNumber(205227458)
  void clearComputemode() => $_clearField(205227458);
}

class QueryComputeResponse extends $pb.GeneratedMessage {
  factory QueryComputeResponse({
    ProvisionedCapacityResponse? provisionedcapacity,
    ComputeMode? computemode,
  }) {
    final result = create();
    if (provisionedcapacity != null)
      result.provisionedcapacity = provisionedcapacity;
    if (computemode != null) result.computemode = computemode;
    return result;
  }

  QueryComputeResponse._();

  factory QueryComputeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryComputeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryComputeResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<ProvisionedCapacityResponse>(
        181400702, _omitFieldNames ? '' : 'provisionedcapacity',
        subBuilder: ProvisionedCapacityResponse.create)
    ..aE<ComputeMode>(205227458, _omitFieldNames ? '' : 'computemode',
        enumValues: ComputeMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryComputeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryComputeResponse copyWith(void Function(QueryComputeResponse) updates) =>
      super.copyWith((message) => updates(message as QueryComputeResponse))
          as QueryComputeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryComputeResponse create() => QueryComputeResponse._();
  @$core.override
  QueryComputeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryComputeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryComputeResponse>(create);
  static QueryComputeResponse? _defaultInstance;

  @$pb.TagNumber(181400702)
  ProvisionedCapacityResponse get provisionedcapacity => $_getN(0);
  @$pb.TagNumber(181400702)
  set provisionedcapacity(ProvisionedCapacityResponse value) =>
      $_setField(181400702, value);
  @$pb.TagNumber(181400702)
  $core.bool hasProvisionedcapacity() => $_has(0);
  @$pb.TagNumber(181400702)
  void clearProvisionedcapacity() => $_clearField(181400702);
  @$pb.TagNumber(181400702)
  ProvisionedCapacityResponse ensureProvisionedcapacity() => $_ensure(0);

  @$pb.TagNumber(205227458)
  ComputeMode get computemode => $_getN(1);
  @$pb.TagNumber(205227458)
  set computemode(ComputeMode value) => $_setField(205227458, value);
  @$pb.TagNumber(205227458)
  $core.bool hasComputemode() => $_has(1);
  @$pb.TagNumber(205227458)
  void clearComputemode() => $_clearField(205227458);
}

class QueryExecutionException extends $pb.GeneratedMessage {
  factory QueryExecutionException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  QueryExecutionException._();

  factory QueryExecutionException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryExecutionException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryExecutionException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecutionException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryExecutionException copyWith(
          void Function(QueryExecutionException) updates) =>
      super.copyWith((message) => updates(message as QueryExecutionException))
          as QueryExecutionException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryExecutionException create() => QueryExecutionException._();
  @$core.override
  QueryExecutionException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryExecutionException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryExecutionException>(create);
  static QueryExecutionException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class QueryInsights extends $pb.GeneratedMessage {
  factory QueryInsights({
    QueryInsightsMode? mode,
  }) {
    final result = create();
    if (mode != null) result.mode = mode;
    return result;
  }

  QueryInsights._();

  factory QueryInsights.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryInsights.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryInsights',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aE<QueryInsightsMode>(323909427, _omitFieldNames ? '' : 'mode',
        enumValues: QueryInsightsMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryInsights clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryInsights copyWith(void Function(QueryInsights) updates) =>
      super.copyWith((message) => updates(message as QueryInsights))
          as QueryInsights;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryInsights create() => QueryInsights._();
  @$core.override
  QueryInsights createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryInsights getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryInsights>(create);
  static QueryInsights? _defaultInstance;

  @$pb.TagNumber(323909427)
  QueryInsightsMode get mode => $_getN(0);
  @$pb.TagNumber(323909427)
  set mode(QueryInsightsMode value) => $_setField(323909427, value);
  @$pb.TagNumber(323909427)
  $core.bool hasMode() => $_has(0);
  @$pb.TagNumber(323909427)
  void clearMode() => $_clearField(323909427);
}

class QueryInsightsResponse extends $pb.GeneratedMessage {
  factory QueryInsightsResponse({
    $fixnum.Int64? unloadwrittenrows,
    $fixnum.Int64? querytablecount,
    $fixnum.Int64? outputrows,
    QuerySpatialCoverage? queryspatialcoverage,
    QueryTemporalRange? querytemporalrange,
    $fixnum.Int64? unloadwrittenbytes,
    $fixnum.Int64? outputbytes,
    $fixnum.Int64? unloadpartitioncount,
  }) {
    final result = create();
    if (unloadwrittenrows != null) result.unloadwrittenrows = unloadwrittenrows;
    if (querytablecount != null) result.querytablecount = querytablecount;
    if (outputrows != null) result.outputrows = outputrows;
    if (queryspatialcoverage != null)
      result.queryspatialcoverage = queryspatialcoverage;
    if (querytemporalrange != null)
      result.querytemporalrange = querytemporalrange;
    if (unloadwrittenbytes != null)
      result.unloadwrittenbytes = unloadwrittenbytes;
    if (outputbytes != null) result.outputbytes = outputbytes;
    if (unloadpartitioncount != null)
      result.unloadpartitioncount = unloadpartitioncount;
    return result;
  }

  QueryInsightsResponse._();

  factory QueryInsightsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryInsightsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryInsightsResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aInt64(20403789, _omitFieldNames ? '' : 'unloadwrittenrows')
    ..aInt64(28893653, _omitFieldNames ? '' : 'querytablecount')
    ..aInt64(138873322, _omitFieldNames ? '' : 'outputrows')
    ..aOM<QuerySpatialCoverage>(
        200825688, _omitFieldNames ? '' : 'queryspatialcoverage',
        subBuilder: QuerySpatialCoverage.create)
    ..aOM<QueryTemporalRange>(
        200957543, _omitFieldNames ? '' : 'querytemporalrange',
        subBuilder: QueryTemporalRange.create)
    ..aInt64(247315473, _omitFieldNames ? '' : 'unloadwrittenbytes')
    ..aInt64(318971400, _omitFieldNames ? '' : 'outputbytes')
    ..aInt64(467201664, _omitFieldNames ? '' : 'unloadpartitioncount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryInsightsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryInsightsResponse copyWith(
          void Function(QueryInsightsResponse) updates) =>
      super.copyWith((message) => updates(message as QueryInsightsResponse))
          as QueryInsightsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryInsightsResponse create() => QueryInsightsResponse._();
  @$core.override
  QueryInsightsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryInsightsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryInsightsResponse>(create);
  static QueryInsightsResponse? _defaultInstance;

  @$pb.TagNumber(20403789)
  $fixnum.Int64 get unloadwrittenrows => $_getI64(0);
  @$pb.TagNumber(20403789)
  set unloadwrittenrows($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(20403789)
  $core.bool hasUnloadwrittenrows() => $_has(0);
  @$pb.TagNumber(20403789)
  void clearUnloadwrittenrows() => $_clearField(20403789);

  @$pb.TagNumber(28893653)
  $fixnum.Int64 get querytablecount => $_getI64(1);
  @$pb.TagNumber(28893653)
  set querytablecount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(28893653)
  $core.bool hasQuerytablecount() => $_has(1);
  @$pb.TagNumber(28893653)
  void clearQuerytablecount() => $_clearField(28893653);

  @$pb.TagNumber(138873322)
  $fixnum.Int64 get outputrows => $_getI64(2);
  @$pb.TagNumber(138873322)
  set outputrows($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(138873322)
  $core.bool hasOutputrows() => $_has(2);
  @$pb.TagNumber(138873322)
  void clearOutputrows() => $_clearField(138873322);

  @$pb.TagNumber(200825688)
  QuerySpatialCoverage get queryspatialcoverage => $_getN(3);
  @$pb.TagNumber(200825688)
  set queryspatialcoverage(QuerySpatialCoverage value) =>
      $_setField(200825688, value);
  @$pb.TagNumber(200825688)
  $core.bool hasQueryspatialcoverage() => $_has(3);
  @$pb.TagNumber(200825688)
  void clearQueryspatialcoverage() => $_clearField(200825688);
  @$pb.TagNumber(200825688)
  QuerySpatialCoverage ensureQueryspatialcoverage() => $_ensure(3);

  @$pb.TagNumber(200957543)
  QueryTemporalRange get querytemporalrange => $_getN(4);
  @$pb.TagNumber(200957543)
  set querytemporalrange(QueryTemporalRange value) =>
      $_setField(200957543, value);
  @$pb.TagNumber(200957543)
  $core.bool hasQuerytemporalrange() => $_has(4);
  @$pb.TagNumber(200957543)
  void clearQuerytemporalrange() => $_clearField(200957543);
  @$pb.TagNumber(200957543)
  QueryTemporalRange ensureQuerytemporalrange() => $_ensure(4);

  @$pb.TagNumber(247315473)
  $fixnum.Int64 get unloadwrittenbytes => $_getI64(5);
  @$pb.TagNumber(247315473)
  set unloadwrittenbytes($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(247315473)
  $core.bool hasUnloadwrittenbytes() => $_has(5);
  @$pb.TagNumber(247315473)
  void clearUnloadwrittenbytes() => $_clearField(247315473);

  @$pb.TagNumber(318971400)
  $fixnum.Int64 get outputbytes => $_getI64(6);
  @$pb.TagNumber(318971400)
  set outputbytes($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(318971400)
  $core.bool hasOutputbytes() => $_has(6);
  @$pb.TagNumber(318971400)
  void clearOutputbytes() => $_clearField(318971400);

  @$pb.TagNumber(467201664)
  $fixnum.Int64 get unloadpartitioncount => $_getI64(7);
  @$pb.TagNumber(467201664)
  set unloadpartitioncount($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(467201664)
  $core.bool hasUnloadpartitioncount() => $_has(7);
  @$pb.TagNumber(467201664)
  void clearUnloadpartitioncount() => $_clearField(467201664);
}

class QueryRequest extends $pb.GeneratedMessage {
  factory QueryRequest({
    QueryInsights? queryinsights,
    $core.String? clienttoken,
    $core.String? nexttoken,
    $core.int? maxrows,
    $core.String? querystring,
  }) {
    final result = create();
    if (queryinsights != null) result.queryinsights = queryinsights;
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxrows != null) result.maxrows = maxrows;
    if (querystring != null) result.querystring = querystring;
    return result;
  }

  QueryRequest._();

  factory QueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<QueryInsights>(105458863, _omitFieldNames ? '' : 'queryinsights',
        subBuilder: QueryInsights.create)
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(251920525, _omitFieldNames ? '' : 'maxrows')
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryRequest copyWith(void Function(QueryRequest) updates) =>
      super.copyWith((message) => updates(message as QueryRequest))
          as QueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryRequest create() => QueryRequest._();
  @$core.override
  QueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryRequest>(create);
  static QueryRequest? _defaultInstance;

  @$pb.TagNumber(105458863)
  QueryInsights get queryinsights => $_getN(0);
  @$pb.TagNumber(105458863)
  set queryinsights(QueryInsights value) => $_setField(105458863, value);
  @$pb.TagNumber(105458863)
  $core.bool hasQueryinsights() => $_has(0);
  @$pb.TagNumber(105458863)
  void clearQueryinsights() => $_clearField(105458863);
  @$pb.TagNumber(105458863)
  QueryInsights ensureQueryinsights() => $_ensure(0);

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(1);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(1);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(251920525)
  $core.int get maxrows => $_getIZ(3);
  @$pb.TagNumber(251920525)
  set maxrows($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(251920525)
  $core.bool hasMaxrows() => $_has(3);
  @$pb.TagNumber(251920525)
  void clearMaxrows() => $_clearField(251920525);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(4);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(4, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(4);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);
}

class QueryResponse extends $pb.GeneratedMessage {
  factory QueryResponse({
    $core.String? queryid,
    $core.Iterable<Row>? rows,
    $core.String? nexttoken,
    QueryInsightsResponse? queryinsightsresponse,
    $core.Iterable<ColumnInfo>? columninfo,
    QueryStatus? querystatus,
  }) {
    final result = create();
    if (queryid != null) result.queryid = queryid;
    if (rows != null) result.rows.addAll(rows);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (queryinsightsresponse != null)
      result.queryinsightsresponse = queryinsightsresponse;
    if (columninfo != null) result.columninfo.addAll(columninfo);
    if (querystatus != null) result.querystatus = querystatus;
    return result;
  }

  QueryResponse._();

  factory QueryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(110737519, _omitFieldNames ? '' : 'queryid')
    ..pPM<Row>(174428857, _omitFieldNames ? '' : 'rows', subBuilder: Row.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOM<QueryInsightsResponse>(
        354278130, _omitFieldNames ? '' : 'queryinsightsresponse',
        subBuilder: QueryInsightsResponse.create)
    ..pPM<ColumnInfo>(364742404, _omitFieldNames ? '' : 'columninfo',
        subBuilder: ColumnInfo.create)
    ..aOM<QueryStatus>(367016406, _omitFieldNames ? '' : 'querystatus',
        subBuilder: QueryStatus.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryResponse copyWith(void Function(QueryResponse) updates) =>
      super.copyWith((message) => updates(message as QueryResponse))
          as QueryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryResponse create() => QueryResponse._();
  @$core.override
  QueryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryResponse>(create);
  static QueryResponse? _defaultInstance;

  @$pb.TagNumber(110737519)
  $core.String get queryid => $_getSZ(0);
  @$pb.TagNumber(110737519)
  set queryid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(110737519)
  $core.bool hasQueryid() => $_has(0);
  @$pb.TagNumber(110737519)
  void clearQueryid() => $_clearField(110737519);

  @$pb.TagNumber(174428857)
  $pb.PbList<Row> get rows => $_getList(1);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(354278130)
  QueryInsightsResponse get queryinsightsresponse => $_getN(3);
  @$pb.TagNumber(354278130)
  set queryinsightsresponse(QueryInsightsResponse value) =>
      $_setField(354278130, value);
  @$pb.TagNumber(354278130)
  $core.bool hasQueryinsightsresponse() => $_has(3);
  @$pb.TagNumber(354278130)
  void clearQueryinsightsresponse() => $_clearField(354278130);
  @$pb.TagNumber(354278130)
  QueryInsightsResponse ensureQueryinsightsresponse() => $_ensure(3);

  @$pb.TagNumber(364742404)
  $pb.PbList<ColumnInfo> get columninfo => $_getList(4);

  @$pb.TagNumber(367016406)
  QueryStatus get querystatus => $_getN(5);
  @$pb.TagNumber(367016406)
  set querystatus(QueryStatus value) => $_setField(367016406, value);
  @$pb.TagNumber(367016406)
  $core.bool hasQuerystatus() => $_has(5);
  @$pb.TagNumber(367016406)
  void clearQuerystatus() => $_clearField(367016406);
  @$pb.TagNumber(367016406)
  QueryStatus ensureQuerystatus() => $_ensure(5);
}

class QuerySpatialCoverage extends $pb.GeneratedMessage {
  factory QuerySpatialCoverage({
    QuerySpatialCoverageMax? max,
  }) {
    final result = create();
    if (max != null) result.max = max;
    return result;
  }

  QuerySpatialCoverage._();

  factory QuerySpatialCoverage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuerySpatialCoverage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuerySpatialCoverage',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<QuerySpatialCoverageMax>(480764858, _omitFieldNames ? '' : 'max',
        subBuilder: QuerySpatialCoverageMax.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuerySpatialCoverage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuerySpatialCoverage copyWith(void Function(QuerySpatialCoverage) updates) =>
      super.copyWith((message) => updates(message as QuerySpatialCoverage))
          as QuerySpatialCoverage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuerySpatialCoverage create() => QuerySpatialCoverage._();
  @$core.override
  QuerySpatialCoverage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuerySpatialCoverage getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuerySpatialCoverage>(create);
  static QuerySpatialCoverage? _defaultInstance;

  @$pb.TagNumber(480764858)
  QuerySpatialCoverageMax get max => $_getN(0);
  @$pb.TagNumber(480764858)
  set max(QuerySpatialCoverageMax value) => $_setField(480764858, value);
  @$pb.TagNumber(480764858)
  $core.bool hasMax() => $_has(0);
  @$pb.TagNumber(480764858)
  void clearMax() => $_clearField(480764858);
  @$pb.TagNumber(480764858)
  QuerySpatialCoverageMax ensureMax() => $_ensure(0);
}

class QuerySpatialCoverageMax extends $pb.GeneratedMessage {
  factory QuerySpatialCoverageMax({
    $core.double? value,
    $core.Iterable<$core.String>? partitionkey,
    $core.String? tablearn,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (partitionkey != null) result.partitionkey.addAll(partitionkey);
    if (tablearn != null) result.tablearn = tablearn;
    return result;
  }

  QuerySpatialCoverageMax._();

  factory QuerySpatialCoverageMax.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuerySpatialCoverageMax.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuerySpatialCoverageMax',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aD(289929579, _omitFieldNames ? '' : 'value')
    ..pPS(379379617, _omitFieldNames ? '' : 'partitionkey')
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuerySpatialCoverageMax clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuerySpatialCoverageMax copyWith(
          void Function(QuerySpatialCoverageMax) updates) =>
      super.copyWith((message) => updates(message as QuerySpatialCoverageMax))
          as QuerySpatialCoverageMax;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuerySpatialCoverageMax create() => QuerySpatialCoverageMax._();
  @$core.override
  QuerySpatialCoverageMax createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuerySpatialCoverageMax getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuerySpatialCoverageMax>(create);
  static QuerySpatialCoverageMax? _defaultInstance;

  @$pb.TagNumber(289929579)
  $core.double get value => $_getN(0);
  @$pb.TagNumber(289929579)
  set value($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(379379617)
  $pb.PbList<$core.String> get partitionkey => $_getList(1);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(2);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(2);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);
}

class QueryStatus extends $pb.GeneratedMessage {
  factory QueryStatus({
    $fixnum.Int64? cumulativebytesmetered,
    $core.double? progresspercentage,
    $fixnum.Int64? cumulativebytesscanned,
  }) {
    final result = create();
    if (cumulativebytesmetered != null)
      result.cumulativebytesmetered = cumulativebytesmetered;
    if (progresspercentage != null)
      result.progresspercentage = progresspercentage;
    if (cumulativebytesscanned != null)
      result.cumulativebytesscanned = cumulativebytesscanned;
    return result;
  }

  QueryStatus._();

  factory QueryStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryStatus',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aInt64(196882576, _omitFieldNames ? '' : 'cumulativebytesmetered')
    ..aD(211894727, _omitFieldNames ? '' : 'progresspercentage')
    ..aInt64(413524176, _omitFieldNames ? '' : 'cumulativebytesscanned')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryStatus copyWith(void Function(QueryStatus) updates) =>
      super.copyWith((message) => updates(message as QueryStatus))
          as QueryStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryStatus create() => QueryStatus._();
  @$core.override
  QueryStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryStatus>(create);
  static QueryStatus? _defaultInstance;

  @$pb.TagNumber(196882576)
  $fixnum.Int64 get cumulativebytesmetered => $_getI64(0);
  @$pb.TagNumber(196882576)
  set cumulativebytesmetered($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(196882576)
  $core.bool hasCumulativebytesmetered() => $_has(0);
  @$pb.TagNumber(196882576)
  void clearCumulativebytesmetered() => $_clearField(196882576);

  @$pb.TagNumber(211894727)
  $core.double get progresspercentage => $_getN(1);
  @$pb.TagNumber(211894727)
  set progresspercentage($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(211894727)
  $core.bool hasProgresspercentage() => $_has(1);
  @$pb.TagNumber(211894727)
  void clearProgresspercentage() => $_clearField(211894727);

  @$pb.TagNumber(413524176)
  $fixnum.Int64 get cumulativebytesscanned => $_getI64(2);
  @$pb.TagNumber(413524176)
  set cumulativebytesscanned($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(413524176)
  $core.bool hasCumulativebytesscanned() => $_has(2);
  @$pb.TagNumber(413524176)
  void clearCumulativebytesscanned() => $_clearField(413524176);
}

class QueryTemporalRange extends $pb.GeneratedMessage {
  factory QueryTemporalRange({
    QueryTemporalRangeMax? max,
  }) {
    final result = create();
    if (max != null) result.max = max;
    return result;
  }

  QueryTemporalRange._();

  factory QueryTemporalRange.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryTemporalRange.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryTemporalRange',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<QueryTemporalRangeMax>(480764858, _omitFieldNames ? '' : 'max',
        subBuilder: QueryTemporalRangeMax.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTemporalRange clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTemporalRange copyWith(void Function(QueryTemporalRange) updates) =>
      super.copyWith((message) => updates(message as QueryTemporalRange))
          as QueryTemporalRange;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTemporalRange create() => QueryTemporalRange._();
  @$core.override
  QueryTemporalRange createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryTemporalRange getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryTemporalRange>(create);
  static QueryTemporalRange? _defaultInstance;

  @$pb.TagNumber(480764858)
  QueryTemporalRangeMax get max => $_getN(0);
  @$pb.TagNumber(480764858)
  set max(QueryTemporalRangeMax value) => $_setField(480764858, value);
  @$pb.TagNumber(480764858)
  $core.bool hasMax() => $_has(0);
  @$pb.TagNumber(480764858)
  void clearMax() => $_clearField(480764858);
  @$pb.TagNumber(480764858)
  QueryTemporalRangeMax ensureMax() => $_ensure(0);
}

class QueryTemporalRangeMax extends $pb.GeneratedMessage {
  factory QueryTemporalRangeMax({
    $fixnum.Int64? value,
    $core.String? tablearn,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (tablearn != null) result.tablearn = tablearn;
    return result;
  }

  QueryTemporalRangeMax._();

  factory QueryTemporalRangeMax.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryTemporalRangeMax.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryTemporalRangeMax',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aInt64(289929579, _omitFieldNames ? '' : 'value')
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTemporalRangeMax clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTemporalRangeMax copyWith(
          void Function(QueryTemporalRangeMax) updates) =>
      super.copyWith((message) => updates(message as QueryTemporalRangeMax))
          as QueryTemporalRangeMax;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTemporalRangeMax create() => QueryTemporalRangeMax._();
  @$core.override
  QueryTemporalRangeMax createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryTemporalRangeMax getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryTemporalRangeMax>(create);
  static QueryTemporalRangeMax? _defaultInstance;

  @$pb.TagNumber(289929579)
  $fixnum.Int64 get value => $_getI64(0);
  @$pb.TagNumber(289929579)
  set value($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(1);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(1);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);
}

class ResourceNotFoundException extends $pb.GeneratedMessage {
  factory ResourceNotFoundException({
    $core.String? scheduledqueryarn,
    $core.String? message,
  }) {
    final result = create();
    if (scheduledqueryarn != null) result.scheduledqueryarn = scheduledqueryarn;
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(234602964, _omitFieldNames ? '' : 'scheduledqueryarn')
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

  @$pb.TagNumber(234602964)
  $core.String get scheduledqueryarn => $_getSZ(0);
  @$pb.TagNumber(234602964)
  set scheduledqueryarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234602964)
  $core.bool hasScheduledqueryarn() => $_has(0);
  @$pb.TagNumber(234602964)
  void clearScheduledqueryarn() => $_clearField(234602964);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class Row extends $pb.GeneratedMessage {
  factory Row({
    $core.Iterable<Datum>? data,
  }) {
    final result = create();
    if (data != null) result.data.addAll(data);
    return result;
  }

  Row._();

  factory Row.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Row.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Row',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..pPM<Datum>(525498822, _omitFieldNames ? '' : 'data',
        subBuilder: Datum.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Row clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Row copyWith(void Function(Row) updates) =>
      super.copyWith((message) => updates(message as Row)) as Row;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Row create() => Row._();
  @$core.override
  Row createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Row getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Row>(create);
  static Row? _defaultInstance;

  @$pb.TagNumber(525498822)
  $pb.PbList<Datum> get data => $_getList(0);
}

class S3Configuration extends $pb.GeneratedMessage {
  factory S3Configuration({
    $core.String? objectkeyprefix,
    S3EncryptionOption? encryptionoption,
    $core.String? bucketname,
  }) {
    final result = create();
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(132617574)
  $core.String get objectkeyprefix => $_getSZ(0);
  @$pb.TagNumber(132617574)
  set objectkeyprefix($core.String value) => $_setString(0, value);
  @$pb.TagNumber(132617574)
  $core.bool hasObjectkeyprefix() => $_has(0);
  @$pb.TagNumber(132617574)
  void clearObjectkeyprefix() => $_clearField(132617574);

  @$pb.TagNumber(160833062)
  S3EncryptionOption get encryptionoption => $_getN(1);
  @$pb.TagNumber(160833062)
  set encryptionoption(S3EncryptionOption value) =>
      $_setField(160833062, value);
  @$pb.TagNumber(160833062)
  $core.bool hasEncryptionoption() => $_has(1);
  @$pb.TagNumber(160833062)
  void clearEncryptionoption() => $_clearField(160833062);

  @$pb.TagNumber(208117045)
  $core.String get bucketname => $_getSZ(2);
  @$pb.TagNumber(208117045)
  set bucketname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(208117045)
  $core.bool hasBucketname() => $_has(2);
  @$pb.TagNumber(208117045)
  void clearBucketname() => $_clearField(208117045);
}

class S3ReportLocation extends $pb.GeneratedMessage {
  factory S3ReportLocation({
    $core.String? bucketname,
    $core.String? objectkey,
  }) {
    final result = create();
    if (bucketname != null) result.bucketname = bucketname;
    if (objectkey != null) result.objectkey = objectkey;
    return result;
  }

  S3ReportLocation._();

  factory S3ReportLocation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory S3ReportLocation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'S3ReportLocation',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(208117045, _omitFieldNames ? '' : 'bucketname')
    ..aOS(335986226, _omitFieldNames ? '' : 'objectkey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3ReportLocation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3ReportLocation copyWith(void Function(S3ReportLocation) updates) =>
      super.copyWith((message) => updates(message as S3ReportLocation))
          as S3ReportLocation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static S3ReportLocation create() => S3ReportLocation._();
  @$core.override
  S3ReportLocation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static S3ReportLocation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<S3ReportLocation>(create);
  static S3ReportLocation? _defaultInstance;

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

class ScheduleConfiguration extends $pb.GeneratedMessage {
  factory ScheduleConfiguration({
    $core.String? scheduleexpression,
  }) {
    final result = create();
    if (scheduleexpression != null)
      result.scheduleexpression = scheduleexpression;
    return result;
  }

  ScheduleConfiguration._();

  factory ScheduleConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduleConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduleConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(446089471, _omitFieldNames ? '' : 'scheduleexpression')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleConfiguration copyWith(
          void Function(ScheduleConfiguration) updates) =>
      super.copyWith((message) => updates(message as ScheduleConfiguration))
          as ScheduleConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduleConfiguration create() => ScheduleConfiguration._();
  @$core.override
  ScheduleConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduleConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduleConfiguration>(create);
  static ScheduleConfiguration? _defaultInstance;

  @$pb.TagNumber(446089471)
  $core.String get scheduleexpression => $_getSZ(0);
  @$pb.TagNumber(446089471)
  set scheduleexpression($core.String value) => $_setString(0, value);
  @$pb.TagNumber(446089471)
  $core.bool hasScheduleexpression() => $_has(0);
  @$pb.TagNumber(446089471)
  void clearScheduleexpression() => $_clearField(446089471);
}

class ScheduledQuery extends $pb.GeneratedMessage {
  factory ScheduledQuery({
    $core.String? previousinvocationtime,
    $core.String? creationtime,
    TargetDestination? targetdestination,
    ErrorReportConfiguration? errorreportconfiguration,
    $core.String? name,
    $core.String? arn,
    $core.String? nextinvocationtime,
    ScheduledQueryRunStatus? lastrunstatus,
    ScheduledQueryState? state,
  }) {
    final result = create();
    if (previousinvocationtime != null)
      result.previousinvocationtime = previousinvocationtime;
    if (creationtime != null) result.creationtime = creationtime;
    if (targetdestination != null) result.targetdestination = targetdestination;
    if (errorreportconfiguration != null)
      result.errorreportconfiguration = errorreportconfiguration;
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    if (nextinvocationtime != null)
      result.nextinvocationtime = nextinvocationtime;
    if (lastrunstatus != null) result.lastrunstatus = lastrunstatus;
    if (state != null) result.state = state;
    return result;
  }

  ScheduledQuery._();

  factory ScheduledQuery.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduledQuery.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduledQuery',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(7530344, _omitFieldNames ? '' : 'previousinvocationtime')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOM<TargetDestination>(
        121667401, _omitFieldNames ? '' : 'targetdestination',
        subBuilder: TargetDestination.create)
    ..aOM<ErrorReportConfiguration>(
        222039776, _omitFieldNames ? '' : 'errorreportconfiguration',
        subBuilder: ErrorReportConfiguration.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(424223272, _omitFieldNames ? '' : 'nextinvocationtime')
    ..aE<ScheduledQueryRunStatus>(
        441976361, _omitFieldNames ? '' : 'lastrunstatus',
        enumValues: ScheduledQueryRunStatus.values)
    ..aE<ScheduledQueryState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ScheduledQueryState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQuery clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQuery copyWith(void Function(ScheduledQuery) updates) =>
      super.copyWith((message) => updates(message as ScheduledQuery))
          as ScheduledQuery;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduledQuery create() => ScheduledQuery._();
  @$core.override
  ScheduledQuery createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduledQuery getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduledQuery>(create);
  static ScheduledQuery? _defaultInstance;

  @$pb.TagNumber(7530344)
  $core.String get previousinvocationtime => $_getSZ(0);
  @$pb.TagNumber(7530344)
  set previousinvocationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7530344)
  $core.bool hasPreviousinvocationtime() => $_has(0);
  @$pb.TagNumber(7530344)
  void clearPreviousinvocationtime() => $_clearField(7530344);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(121667401)
  TargetDestination get targetdestination => $_getN(2);
  @$pb.TagNumber(121667401)
  set targetdestination(TargetDestination value) =>
      $_setField(121667401, value);
  @$pb.TagNumber(121667401)
  $core.bool hasTargetdestination() => $_has(2);
  @$pb.TagNumber(121667401)
  void clearTargetdestination() => $_clearField(121667401);
  @$pb.TagNumber(121667401)
  TargetDestination ensureTargetdestination() => $_ensure(2);

  @$pb.TagNumber(222039776)
  ErrorReportConfiguration get errorreportconfiguration => $_getN(3);
  @$pb.TagNumber(222039776)
  set errorreportconfiguration(ErrorReportConfiguration value) =>
      $_setField(222039776, value);
  @$pb.TagNumber(222039776)
  $core.bool hasErrorreportconfiguration() => $_has(3);
  @$pb.TagNumber(222039776)
  void clearErrorreportconfiguration() => $_clearField(222039776);
  @$pb.TagNumber(222039776)
  ErrorReportConfiguration ensureErrorreportconfiguration() => $_ensure(3);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(5);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(5);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(424223272)
  $core.String get nextinvocationtime => $_getSZ(6);
  @$pb.TagNumber(424223272)
  set nextinvocationtime($core.String value) => $_setString(6, value);
  @$pb.TagNumber(424223272)
  $core.bool hasNextinvocationtime() => $_has(6);
  @$pb.TagNumber(424223272)
  void clearNextinvocationtime() => $_clearField(424223272);

  @$pb.TagNumber(441976361)
  ScheduledQueryRunStatus get lastrunstatus => $_getN(7);
  @$pb.TagNumber(441976361)
  set lastrunstatus(ScheduledQueryRunStatus value) =>
      $_setField(441976361, value);
  @$pb.TagNumber(441976361)
  $core.bool hasLastrunstatus() => $_has(7);
  @$pb.TagNumber(441976361)
  void clearLastrunstatus() => $_clearField(441976361);

  @$pb.TagNumber(502047895)
  ScheduledQueryState get state => $_getN(8);
  @$pb.TagNumber(502047895)
  set state(ScheduledQueryState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(8);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class ScheduledQueryDescription extends $pb.GeneratedMessage {
  factory ScheduledQueryDescription({
    $core.String? previousinvocationtime,
    $core.String? kmskeyid,
    $core.String? creationtime,
    ErrorReportConfiguration? errorreportconfiguration,
    ScheduledQueryRunSummary? lastrunsummary,
    $core.String? name,
    $core.String? scheduledqueryexecutionrolearn,
    NotificationConfiguration? notificationconfiguration,
    $core.Iterable<ScheduledQueryRunSummary>? recentlyfailedruns,
    $core.String? arn,
    $core.String? nextinvocationtime,
    $core.String? querystring,
    TargetConfiguration? targetconfiguration,
    ScheduledQueryState? state,
    ScheduleConfiguration? scheduleconfiguration,
  }) {
    final result = create();
    if (previousinvocationtime != null)
      result.previousinvocationtime = previousinvocationtime;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (creationtime != null) result.creationtime = creationtime;
    if (errorreportconfiguration != null)
      result.errorreportconfiguration = errorreportconfiguration;
    if (lastrunsummary != null) result.lastrunsummary = lastrunsummary;
    if (name != null) result.name = name;
    if (scheduledqueryexecutionrolearn != null)
      result.scheduledqueryexecutionrolearn = scheduledqueryexecutionrolearn;
    if (notificationconfiguration != null)
      result.notificationconfiguration = notificationconfiguration;
    if (recentlyfailedruns != null)
      result.recentlyfailedruns.addAll(recentlyfailedruns);
    if (arn != null) result.arn = arn;
    if (nextinvocationtime != null)
      result.nextinvocationtime = nextinvocationtime;
    if (querystring != null) result.querystring = querystring;
    if (targetconfiguration != null)
      result.targetconfiguration = targetconfiguration;
    if (state != null) result.state = state;
    if (scheduleconfiguration != null)
      result.scheduleconfiguration = scheduleconfiguration;
    return result;
  }

  ScheduledQueryDescription._();

  factory ScheduledQueryDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduledQueryDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduledQueryDescription',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(7530344, _omitFieldNames ? '' : 'previousinvocationtime')
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOM<ErrorReportConfiguration>(
        222039776, _omitFieldNames ? '' : 'errorreportconfiguration',
        subBuilder: ErrorReportConfiguration.create)
    ..aOM<ScheduledQueryRunSummary>(
        238724003, _omitFieldNames ? '' : 'lastrunsummary',
        subBuilder: ScheduledQueryRunSummary.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(284244182, _omitFieldNames ? '' : 'scheduledqueryexecutionrolearn')
    ..aOM<NotificationConfiguration>(
        290208045, _omitFieldNames ? '' : 'notificationconfiguration',
        subBuilder: NotificationConfiguration.create)
    ..pPM<ScheduledQueryRunSummary>(
        365830319, _omitFieldNames ? '' : 'recentlyfailedruns',
        subBuilder: ScheduledQueryRunSummary.create)
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(424223272, _omitFieldNames ? '' : 'nextinvocationtime')
    ..aOS(435938663, _omitFieldNames ? '' : 'querystring')
    ..aOM<TargetConfiguration>(
        499845483, _omitFieldNames ? '' : 'targetconfiguration',
        subBuilder: TargetConfiguration.create)
    ..aE<ScheduledQueryState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ScheduledQueryState.values)
    ..aOM<ScheduleConfiguration>(
        526075047, _omitFieldNames ? '' : 'scheduleconfiguration',
        subBuilder: ScheduleConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQueryDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQueryDescription copyWith(
          void Function(ScheduledQueryDescription) updates) =>
      super.copyWith((message) => updates(message as ScheduledQueryDescription))
          as ScheduledQueryDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduledQueryDescription create() => ScheduledQueryDescription._();
  @$core.override
  ScheduledQueryDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduledQueryDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduledQueryDescription>(create);
  static ScheduledQueryDescription? _defaultInstance;

  @$pb.TagNumber(7530344)
  $core.String get previousinvocationtime => $_getSZ(0);
  @$pb.TagNumber(7530344)
  set previousinvocationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7530344)
  $core.bool hasPreviousinvocationtime() => $_has(0);
  @$pb.TagNumber(7530344)
  void clearPreviousinvocationtime() => $_clearField(7530344);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(1);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(1);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(222039776)
  ErrorReportConfiguration get errorreportconfiguration => $_getN(3);
  @$pb.TagNumber(222039776)
  set errorreportconfiguration(ErrorReportConfiguration value) =>
      $_setField(222039776, value);
  @$pb.TagNumber(222039776)
  $core.bool hasErrorreportconfiguration() => $_has(3);
  @$pb.TagNumber(222039776)
  void clearErrorreportconfiguration() => $_clearField(222039776);
  @$pb.TagNumber(222039776)
  ErrorReportConfiguration ensureErrorreportconfiguration() => $_ensure(3);

  @$pb.TagNumber(238724003)
  ScheduledQueryRunSummary get lastrunsummary => $_getN(4);
  @$pb.TagNumber(238724003)
  set lastrunsummary(ScheduledQueryRunSummary value) =>
      $_setField(238724003, value);
  @$pb.TagNumber(238724003)
  $core.bool hasLastrunsummary() => $_has(4);
  @$pb.TagNumber(238724003)
  void clearLastrunsummary() => $_clearField(238724003);
  @$pb.TagNumber(238724003)
  ScheduledQueryRunSummary ensureLastrunsummary() => $_ensure(4);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(284244182)
  $core.String get scheduledqueryexecutionrolearn => $_getSZ(6);
  @$pb.TagNumber(284244182)
  set scheduledqueryexecutionrolearn($core.String value) =>
      $_setString(6, value);
  @$pb.TagNumber(284244182)
  $core.bool hasScheduledqueryexecutionrolearn() => $_has(6);
  @$pb.TagNumber(284244182)
  void clearScheduledqueryexecutionrolearn() => $_clearField(284244182);

  @$pb.TagNumber(290208045)
  NotificationConfiguration get notificationconfiguration => $_getN(7);
  @$pb.TagNumber(290208045)
  set notificationconfiguration(NotificationConfiguration value) =>
      $_setField(290208045, value);
  @$pb.TagNumber(290208045)
  $core.bool hasNotificationconfiguration() => $_has(7);
  @$pb.TagNumber(290208045)
  void clearNotificationconfiguration() => $_clearField(290208045);
  @$pb.TagNumber(290208045)
  NotificationConfiguration ensureNotificationconfiguration() => $_ensure(7);

  @$pb.TagNumber(365830319)
  $pb.PbList<ScheduledQueryRunSummary> get recentlyfailedruns => $_getList(8);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(9);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(9);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(424223272)
  $core.String get nextinvocationtime => $_getSZ(10);
  @$pb.TagNumber(424223272)
  set nextinvocationtime($core.String value) => $_setString(10, value);
  @$pb.TagNumber(424223272)
  $core.bool hasNextinvocationtime() => $_has(10);
  @$pb.TagNumber(424223272)
  void clearNextinvocationtime() => $_clearField(424223272);

  @$pb.TagNumber(435938663)
  $core.String get querystring => $_getSZ(11);
  @$pb.TagNumber(435938663)
  set querystring($core.String value) => $_setString(11, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(11);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);

  @$pb.TagNumber(499845483)
  TargetConfiguration get targetconfiguration => $_getN(12);
  @$pb.TagNumber(499845483)
  set targetconfiguration(TargetConfiguration value) =>
      $_setField(499845483, value);
  @$pb.TagNumber(499845483)
  $core.bool hasTargetconfiguration() => $_has(12);
  @$pb.TagNumber(499845483)
  void clearTargetconfiguration() => $_clearField(499845483);
  @$pb.TagNumber(499845483)
  TargetConfiguration ensureTargetconfiguration() => $_ensure(12);

  @$pb.TagNumber(502047895)
  ScheduledQueryState get state => $_getN(13);
  @$pb.TagNumber(502047895)
  set state(ScheduledQueryState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(13);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);

  @$pb.TagNumber(526075047)
  ScheduleConfiguration get scheduleconfiguration => $_getN(14);
  @$pb.TagNumber(526075047)
  set scheduleconfiguration(ScheduleConfiguration value) =>
      $_setField(526075047, value);
  @$pb.TagNumber(526075047)
  $core.bool hasScheduleconfiguration() => $_has(14);
  @$pb.TagNumber(526075047)
  void clearScheduleconfiguration() => $_clearField(526075047);
  @$pb.TagNumber(526075047)
  ScheduleConfiguration ensureScheduleconfiguration() => $_ensure(14);
}

class ScheduledQueryInsights extends $pb.GeneratedMessage {
  factory ScheduledQueryInsights({
    ScheduledQueryInsightsMode? mode,
  }) {
    final result = create();
    if (mode != null) result.mode = mode;
    return result;
  }

  ScheduledQueryInsights._();

  factory ScheduledQueryInsights.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduledQueryInsights.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduledQueryInsights',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aE<ScheduledQueryInsightsMode>(323909427, _omitFieldNames ? '' : 'mode',
        enumValues: ScheduledQueryInsightsMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQueryInsights clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQueryInsights copyWith(
          void Function(ScheduledQueryInsights) updates) =>
      super.copyWith((message) => updates(message as ScheduledQueryInsights))
          as ScheduledQueryInsights;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduledQueryInsights create() => ScheduledQueryInsights._();
  @$core.override
  ScheduledQueryInsights createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduledQueryInsights getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduledQueryInsights>(create);
  static ScheduledQueryInsights? _defaultInstance;

  @$pb.TagNumber(323909427)
  ScheduledQueryInsightsMode get mode => $_getN(0);
  @$pb.TagNumber(323909427)
  set mode(ScheduledQueryInsightsMode value) => $_setField(323909427, value);
  @$pb.TagNumber(323909427)
  $core.bool hasMode() => $_has(0);
  @$pb.TagNumber(323909427)
  void clearMode() => $_clearField(323909427);
}

class ScheduledQueryInsightsResponse extends $pb.GeneratedMessage {
  factory ScheduledQueryInsightsResponse({
    $fixnum.Int64? querytablecount,
    $fixnum.Int64? outputrows,
    QuerySpatialCoverage? queryspatialcoverage,
    QueryTemporalRange? querytemporalrange,
    $fixnum.Int64? outputbytes,
  }) {
    final result = create();
    if (querytablecount != null) result.querytablecount = querytablecount;
    if (outputrows != null) result.outputrows = outputrows;
    if (queryspatialcoverage != null)
      result.queryspatialcoverage = queryspatialcoverage;
    if (querytemporalrange != null)
      result.querytemporalrange = querytemporalrange;
    if (outputbytes != null) result.outputbytes = outputbytes;
    return result;
  }

  ScheduledQueryInsightsResponse._();

  factory ScheduledQueryInsightsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduledQueryInsightsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduledQueryInsightsResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aInt64(28893653, _omitFieldNames ? '' : 'querytablecount')
    ..aInt64(138873322, _omitFieldNames ? '' : 'outputrows')
    ..aOM<QuerySpatialCoverage>(
        200825688, _omitFieldNames ? '' : 'queryspatialcoverage',
        subBuilder: QuerySpatialCoverage.create)
    ..aOM<QueryTemporalRange>(
        200957543, _omitFieldNames ? '' : 'querytemporalrange',
        subBuilder: QueryTemporalRange.create)
    ..aInt64(318971400, _omitFieldNames ? '' : 'outputbytes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQueryInsightsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQueryInsightsResponse copyWith(
          void Function(ScheduledQueryInsightsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ScheduledQueryInsightsResponse))
          as ScheduledQueryInsightsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduledQueryInsightsResponse create() =>
      ScheduledQueryInsightsResponse._();
  @$core.override
  ScheduledQueryInsightsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduledQueryInsightsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduledQueryInsightsResponse>(create);
  static ScheduledQueryInsightsResponse? _defaultInstance;

  @$pb.TagNumber(28893653)
  $fixnum.Int64 get querytablecount => $_getI64(0);
  @$pb.TagNumber(28893653)
  set querytablecount($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(28893653)
  $core.bool hasQuerytablecount() => $_has(0);
  @$pb.TagNumber(28893653)
  void clearQuerytablecount() => $_clearField(28893653);

  @$pb.TagNumber(138873322)
  $fixnum.Int64 get outputrows => $_getI64(1);
  @$pb.TagNumber(138873322)
  set outputrows($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(138873322)
  $core.bool hasOutputrows() => $_has(1);
  @$pb.TagNumber(138873322)
  void clearOutputrows() => $_clearField(138873322);

  @$pb.TagNumber(200825688)
  QuerySpatialCoverage get queryspatialcoverage => $_getN(2);
  @$pb.TagNumber(200825688)
  set queryspatialcoverage(QuerySpatialCoverage value) =>
      $_setField(200825688, value);
  @$pb.TagNumber(200825688)
  $core.bool hasQueryspatialcoverage() => $_has(2);
  @$pb.TagNumber(200825688)
  void clearQueryspatialcoverage() => $_clearField(200825688);
  @$pb.TagNumber(200825688)
  QuerySpatialCoverage ensureQueryspatialcoverage() => $_ensure(2);

  @$pb.TagNumber(200957543)
  QueryTemporalRange get querytemporalrange => $_getN(3);
  @$pb.TagNumber(200957543)
  set querytemporalrange(QueryTemporalRange value) =>
      $_setField(200957543, value);
  @$pb.TagNumber(200957543)
  $core.bool hasQuerytemporalrange() => $_has(3);
  @$pb.TagNumber(200957543)
  void clearQuerytemporalrange() => $_clearField(200957543);
  @$pb.TagNumber(200957543)
  QueryTemporalRange ensureQuerytemporalrange() => $_ensure(3);

  @$pb.TagNumber(318971400)
  $fixnum.Int64 get outputbytes => $_getI64(4);
  @$pb.TagNumber(318971400)
  set outputbytes($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(318971400)
  $core.bool hasOutputbytes() => $_has(4);
  @$pb.TagNumber(318971400)
  void clearOutputbytes() => $_clearField(318971400);
}

class ScheduledQueryRunSummary extends $pb.GeneratedMessage {
  factory ScheduledQueryRunSummary({
    ErrorReportLocation? errorreportlocation,
    ExecutionStats? executionstats,
    $core.String? failurereason,
    $core.String? triggertime,
    ScheduledQueryRunStatus? runstatus,
    $core.String? invocationtime,
    ScheduledQueryInsightsResponse? queryinsightsresponse,
  }) {
    final result = create();
    if (errorreportlocation != null)
      result.errorreportlocation = errorreportlocation;
    if (executionstats != null) result.executionstats = executionstats;
    if (failurereason != null) result.failurereason = failurereason;
    if (triggertime != null) result.triggertime = triggertime;
    if (runstatus != null) result.runstatus = runstatus;
    if (invocationtime != null) result.invocationtime = invocationtime;
    if (queryinsightsresponse != null)
      result.queryinsightsresponse = queryinsightsresponse;
    return result;
  }

  ScheduledQueryRunSummary._();

  factory ScheduledQueryRunSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduledQueryRunSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduledQueryRunSummary',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<ErrorReportLocation>(
        117286905, _omitFieldNames ? '' : 'errorreportlocation',
        subBuilder: ErrorReportLocation.create)
    ..aOM<ExecutionStats>(220056093, _omitFieldNames ? '' : 'executionstats',
        subBuilder: ExecutionStats.create)
    ..aOS(232322142, _omitFieldNames ? '' : 'failurereason')
    ..aOS(268796699, _omitFieldNames ? '' : 'triggertime')
    ..aE<ScheduledQueryRunStatus>(293822805, _omitFieldNames ? '' : 'runstatus',
        enumValues: ScheduledQueryRunStatus.values)
    ..aOS(331845291, _omitFieldNames ? '' : 'invocationtime')
    ..aOM<ScheduledQueryInsightsResponse>(
        354278130, _omitFieldNames ? '' : 'queryinsightsresponse',
        subBuilder: ScheduledQueryInsightsResponse.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQueryRunSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduledQueryRunSummary copyWith(
          void Function(ScheduledQueryRunSummary) updates) =>
      super.copyWith((message) => updates(message as ScheduledQueryRunSummary))
          as ScheduledQueryRunSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduledQueryRunSummary create() => ScheduledQueryRunSummary._();
  @$core.override
  ScheduledQueryRunSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduledQueryRunSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduledQueryRunSummary>(create);
  static ScheduledQueryRunSummary? _defaultInstance;

  @$pb.TagNumber(117286905)
  ErrorReportLocation get errorreportlocation => $_getN(0);
  @$pb.TagNumber(117286905)
  set errorreportlocation(ErrorReportLocation value) =>
      $_setField(117286905, value);
  @$pb.TagNumber(117286905)
  $core.bool hasErrorreportlocation() => $_has(0);
  @$pb.TagNumber(117286905)
  void clearErrorreportlocation() => $_clearField(117286905);
  @$pb.TagNumber(117286905)
  ErrorReportLocation ensureErrorreportlocation() => $_ensure(0);

  @$pb.TagNumber(220056093)
  ExecutionStats get executionstats => $_getN(1);
  @$pb.TagNumber(220056093)
  set executionstats(ExecutionStats value) => $_setField(220056093, value);
  @$pb.TagNumber(220056093)
  $core.bool hasExecutionstats() => $_has(1);
  @$pb.TagNumber(220056093)
  void clearExecutionstats() => $_clearField(220056093);
  @$pb.TagNumber(220056093)
  ExecutionStats ensureExecutionstats() => $_ensure(1);

  @$pb.TagNumber(232322142)
  $core.String get failurereason => $_getSZ(2);
  @$pb.TagNumber(232322142)
  set failurereason($core.String value) => $_setString(2, value);
  @$pb.TagNumber(232322142)
  $core.bool hasFailurereason() => $_has(2);
  @$pb.TagNumber(232322142)
  void clearFailurereason() => $_clearField(232322142);

  @$pb.TagNumber(268796699)
  $core.String get triggertime => $_getSZ(3);
  @$pb.TagNumber(268796699)
  set triggertime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(268796699)
  $core.bool hasTriggertime() => $_has(3);
  @$pb.TagNumber(268796699)
  void clearTriggertime() => $_clearField(268796699);

  @$pb.TagNumber(293822805)
  ScheduledQueryRunStatus get runstatus => $_getN(4);
  @$pb.TagNumber(293822805)
  set runstatus(ScheduledQueryRunStatus value) => $_setField(293822805, value);
  @$pb.TagNumber(293822805)
  $core.bool hasRunstatus() => $_has(4);
  @$pb.TagNumber(293822805)
  void clearRunstatus() => $_clearField(293822805);

  @$pb.TagNumber(331845291)
  $core.String get invocationtime => $_getSZ(5);
  @$pb.TagNumber(331845291)
  set invocationtime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(331845291)
  $core.bool hasInvocationtime() => $_has(5);
  @$pb.TagNumber(331845291)
  void clearInvocationtime() => $_clearField(331845291);

  @$pb.TagNumber(354278130)
  ScheduledQueryInsightsResponse get queryinsightsresponse => $_getN(6);
  @$pb.TagNumber(354278130)
  set queryinsightsresponse(ScheduledQueryInsightsResponse value) =>
      $_setField(354278130, value);
  @$pb.TagNumber(354278130)
  $core.bool hasQueryinsightsresponse() => $_has(6);
  @$pb.TagNumber(354278130)
  void clearQueryinsightsresponse() => $_clearField(354278130);
  @$pb.TagNumber(354278130)
  ScheduledQueryInsightsResponse ensureQueryinsightsresponse() => $_ensure(6);
}

class SelectColumn extends $pb.GeneratedMessage {
  factory SelectColumn({
    $core.String? databasename,
    $core.bool? aliased,
    $core.String? name,
    $core.String? tablename,
    Type? type,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (aliased != null) result.aliased = aliased;
    if (name != null) result.name = name;
    if (tablename != null) result.tablename = tablename;
    if (type != null) result.type = type;
    return result;
  }

  SelectColumn._();

  factory SelectColumn.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SelectColumn.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SelectColumn',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOB(157931831, _omitFieldNames ? '' : 'aliased')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<Type>(290836590, _omitFieldNames ? '' : 'type',
        subBuilder: Type.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SelectColumn clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SelectColumn copyWith(void Function(SelectColumn) updates) =>
      super.copyWith((message) => updates(message as SelectColumn))
          as SelectColumn;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SelectColumn create() => SelectColumn._();
  @$core.override
  SelectColumn createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SelectColumn getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SelectColumn>(create);
  static SelectColumn? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(157931831)
  $core.bool get aliased => $_getBF(1);
  @$pb.TagNumber(157931831)
  set aliased($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(157931831)
  $core.bool hasAliased() => $_has(1);
  @$pb.TagNumber(157931831)
  void clearAliased() => $_clearField(157931831);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(290836590)
  Type get type => $_getN(4);
  @$pb.TagNumber(290836590)
  set type(Type value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(4);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
  @$pb.TagNumber(290836590)
  Type ensureType() => $_ensure(4);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class SnsConfiguration extends $pb.GeneratedMessage {
  factory SnsConfiguration({
    $core.String? topicarn,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    return result;
  }

  SnsConfiguration._();

  factory SnsConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SnsConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SnsConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SnsConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SnsConfiguration copyWith(void Function(SnsConfiguration) updates) =>
      super.copyWith((message) => updates(message as SnsConfiguration))
          as SnsConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SnsConfiguration create() => SnsConfiguration._();
  @$core.override
  SnsConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SnsConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SnsConfiguration>(create);
  static SnsConfiguration? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class TargetConfiguration extends $pb.GeneratedMessage {
  factory TargetConfiguration({
    TimestreamConfiguration? timestreamconfiguration,
  }) {
    final result = create();
    if (timestreamconfiguration != null)
      result.timestreamconfiguration = timestreamconfiguration;
    return result;
  }

  TargetConfiguration._();

  factory TargetConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TargetConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TargetConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<TimestreamConfiguration>(
        30086935, _omitFieldNames ? '' : 'timestreamconfiguration',
        subBuilder: TimestreamConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TargetConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TargetConfiguration copyWith(void Function(TargetConfiguration) updates) =>
      super.copyWith((message) => updates(message as TargetConfiguration))
          as TargetConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TargetConfiguration create() => TargetConfiguration._();
  @$core.override
  TargetConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TargetConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TargetConfiguration>(create);
  static TargetConfiguration? _defaultInstance;

  @$pb.TagNumber(30086935)
  TimestreamConfiguration get timestreamconfiguration => $_getN(0);
  @$pb.TagNumber(30086935)
  set timestreamconfiguration(TimestreamConfiguration value) =>
      $_setField(30086935, value);
  @$pb.TagNumber(30086935)
  $core.bool hasTimestreamconfiguration() => $_has(0);
  @$pb.TagNumber(30086935)
  void clearTimestreamconfiguration() => $_clearField(30086935);
  @$pb.TagNumber(30086935)
  TimestreamConfiguration ensureTimestreamconfiguration() => $_ensure(0);
}

class TargetDestination extends $pb.GeneratedMessage {
  factory TargetDestination({
    TimestreamDestination? timestreamdestination,
  }) {
    final result = create();
    if (timestreamdestination != null)
      result.timestreamdestination = timestreamdestination;
    return result;
  }

  TargetDestination._();

  factory TargetDestination.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TargetDestination.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TargetDestination',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<TimestreamDestination>(
        259349965, _omitFieldNames ? '' : 'timestreamdestination',
        subBuilder: TimestreamDestination.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TargetDestination clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TargetDestination copyWith(void Function(TargetDestination) updates) =>
      super.copyWith((message) => updates(message as TargetDestination))
          as TargetDestination;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TargetDestination create() => TargetDestination._();
  @$core.override
  TargetDestination createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TargetDestination getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TargetDestination>(create);
  static TargetDestination? _defaultInstance;

  @$pb.TagNumber(259349965)
  TimestreamDestination get timestreamdestination => $_getN(0);
  @$pb.TagNumber(259349965)
  set timestreamdestination(TimestreamDestination value) =>
      $_setField(259349965, value);
  @$pb.TagNumber(259349965)
  $core.bool hasTimestreamdestination() => $_has(0);
  @$pb.TagNumber(259349965)
  void clearTimestreamdestination() => $_clearField(259349965);
  @$pb.TagNumber(259349965)
  TimestreamDestination ensureTimestreamdestination() => $_ensure(0);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class TimeSeriesDataPoint extends $pb.GeneratedMessage {
  factory TimeSeriesDataPoint({
    Datum? value,
    $core.String? time,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (time != null) result.time = time;
    return result;
  }

  TimeSeriesDataPoint._();

  factory TimeSeriesDataPoint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TimeSeriesDataPoint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TimeSeriesDataPoint',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOM<Datum>(289929579, _omitFieldNames ? '' : 'value',
        subBuilder: Datum.create)
    ..aOS(535094277, _omitFieldNames ? '' : 'time')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeSeriesDataPoint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeSeriesDataPoint copyWith(void Function(TimeSeriesDataPoint) updates) =>
      super.copyWith((message) => updates(message as TimeSeriesDataPoint))
          as TimeSeriesDataPoint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TimeSeriesDataPoint create() => TimeSeriesDataPoint._();
  @$core.override
  TimeSeriesDataPoint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TimeSeriesDataPoint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TimeSeriesDataPoint>(create);
  static TimeSeriesDataPoint? _defaultInstance;

  @$pb.TagNumber(289929579)
  Datum get value => $_getN(0);
  @$pb.TagNumber(289929579)
  set value(Datum value) => $_setField(289929579, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
  @$pb.TagNumber(289929579)
  Datum ensureValue() => $_ensure(0);

  @$pb.TagNumber(535094277)
  $core.String get time => $_getSZ(1);
  @$pb.TagNumber(535094277)
  set time($core.String value) => $_setString(1, value);
  @$pb.TagNumber(535094277)
  $core.bool hasTime() => $_has(1);
  @$pb.TagNumber(535094277)
  void clearTime() => $_clearField(535094277);
}

class TimestreamConfiguration extends $pb.GeneratedMessage {
  factory TimestreamConfiguration({
    $core.String? databasename,
    $core.String? measurenamecolumn,
    $core.Iterable<DimensionMapping>? dimensionmappings,
    $core.String? tablename,
    MultiMeasureMappings? multimeasuremappings,
    $core.String? timecolumn,
    $core.Iterable<MixedMeasureMapping>? mixedmeasuremappings,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (measurenamecolumn != null) result.measurenamecolumn = measurenamecolumn;
    if (dimensionmappings != null)
      result.dimensionmappings.addAll(dimensionmappings);
    if (tablename != null) result.tablename = tablename;
    if (multimeasuremappings != null)
      result.multimeasuremappings = multimeasuremappings;
    if (timecolumn != null) result.timecolumn = timecolumn;
    if (mixedmeasuremappings != null)
      result.mixedmeasuremappings.addAll(mixedmeasuremappings);
    return result;
  }

  TimestreamConfiguration._();

  factory TimestreamConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TimestreamConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TimestreamConfiguration',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(112141775, _omitFieldNames ? '' : 'measurenamecolumn')
    ..pPM<DimensionMapping>(
        224443741, _omitFieldNames ? '' : 'dimensionmappings',
        subBuilder: DimensionMapping.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<MultiMeasureMappings>(
        501736394, _omitFieldNames ? '' : 'multimeasuremappings',
        subBuilder: MultiMeasureMappings.create)
    ..aOS(519449367, _omitFieldNames ? '' : 'timecolumn')
    ..pPM<MixedMeasureMapping>(
        521774144, _omitFieldNames ? '' : 'mixedmeasuremappings',
        subBuilder: MixedMeasureMapping.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimestreamConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimestreamConfiguration copyWith(
          void Function(TimestreamConfiguration) updates) =>
      super.copyWith((message) => updates(message as TimestreamConfiguration))
          as TimestreamConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TimestreamConfiguration create() => TimestreamConfiguration._();
  @$core.override
  TimestreamConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TimestreamConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TimestreamConfiguration>(create);
  static TimestreamConfiguration? _defaultInstance;

  @$pb.TagNumber(89545052)
  $core.String get databasename => $_getSZ(0);
  @$pb.TagNumber(89545052)
  set databasename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89545052)
  $core.bool hasDatabasename() => $_has(0);
  @$pb.TagNumber(89545052)
  void clearDatabasename() => $_clearField(89545052);

  @$pb.TagNumber(112141775)
  $core.String get measurenamecolumn => $_getSZ(1);
  @$pb.TagNumber(112141775)
  set measurenamecolumn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(112141775)
  $core.bool hasMeasurenamecolumn() => $_has(1);
  @$pb.TagNumber(112141775)
  void clearMeasurenamecolumn() => $_clearField(112141775);

  @$pb.TagNumber(224443741)
  $pb.PbList<DimensionMapping> get dimensionmappings => $_getList(2);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(501736394)
  MultiMeasureMappings get multimeasuremappings => $_getN(4);
  @$pb.TagNumber(501736394)
  set multimeasuremappings(MultiMeasureMappings value) =>
      $_setField(501736394, value);
  @$pb.TagNumber(501736394)
  $core.bool hasMultimeasuremappings() => $_has(4);
  @$pb.TagNumber(501736394)
  void clearMultimeasuremappings() => $_clearField(501736394);
  @$pb.TagNumber(501736394)
  MultiMeasureMappings ensureMultimeasuremappings() => $_ensure(4);

  @$pb.TagNumber(519449367)
  $core.String get timecolumn => $_getSZ(5);
  @$pb.TagNumber(519449367)
  set timecolumn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(519449367)
  $core.bool hasTimecolumn() => $_has(5);
  @$pb.TagNumber(519449367)
  void clearTimecolumn() => $_clearField(519449367);

  @$pb.TagNumber(521774144)
  $pb.PbList<MixedMeasureMapping> get mixedmeasuremappings => $_getList(6);
}

class TimestreamDestination extends $pb.GeneratedMessage {
  factory TimestreamDestination({
    $core.String? databasename,
    $core.String? tablename,
  }) {
    final result = create();
    if (databasename != null) result.databasename = databasename;
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  TimestreamDestination._();

  factory TimestreamDestination.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TimestreamDestination.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TimestreamDestination',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(89545052, _omitFieldNames ? '' : 'databasename')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimestreamDestination clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimestreamDestination copyWith(
          void Function(TimestreamDestination) updates) =>
      super.copyWith((message) => updates(message as TimestreamDestination))
          as TimestreamDestination;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TimestreamDestination create() => TimestreamDestination._();
  @$core.override
  TimestreamDestination createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TimestreamDestination getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TimestreamDestination>(create);
  static TimestreamDestination? _defaultInstance;

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

class Type extends $pb.GeneratedMessage {
  factory Type({
    $core.Iterable<ColumnInfo>? rowcolumninfo,
    ScalarType? scalartype,
    ColumnInfo? timeseriesmeasurevaluecolumninfo,
    ColumnInfo? arraycolumninfo,
  }) {
    final result = create();
    if (rowcolumninfo != null) result.rowcolumninfo.addAll(rowcolumninfo);
    if (scalartype != null) result.scalartype = scalartype;
    if (timeseriesmeasurevaluecolumninfo != null)
      result.timeseriesmeasurevaluecolumninfo =
          timeseriesmeasurevaluecolumninfo;
    if (arraycolumninfo != null) result.arraycolumninfo = arraycolumninfo;
    return result;
  }

  Type._();

  factory Type.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Type.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Type',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..pPM<ColumnInfo>(37626042, _omitFieldNames ? '' : 'rowcolumninfo',
        subBuilder: ColumnInfo.create)
    ..aE<ScalarType>(316685622, _omitFieldNames ? '' : 'scalartype',
        enumValues: ScalarType.values)
    ..aOM<ColumnInfo>(
        324829411, _omitFieldNames ? '' : 'timeseriesmeasurevaluecolumninfo',
        subBuilder: ColumnInfo.create)
    ..aOM<ColumnInfo>(476744353, _omitFieldNames ? '' : 'arraycolumninfo',
        subBuilder: ColumnInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Type clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Type copyWith(void Function(Type) updates) =>
      super.copyWith((message) => updates(message as Type)) as Type;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Type create() => Type._();
  @$core.override
  Type createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Type getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Type>(create);
  static Type? _defaultInstance;

  @$pb.TagNumber(37626042)
  $pb.PbList<ColumnInfo> get rowcolumninfo => $_getList(0);

  @$pb.TagNumber(316685622)
  ScalarType get scalartype => $_getN(1);
  @$pb.TagNumber(316685622)
  set scalartype(ScalarType value) => $_setField(316685622, value);
  @$pb.TagNumber(316685622)
  $core.bool hasScalartype() => $_has(1);
  @$pb.TagNumber(316685622)
  void clearScalartype() => $_clearField(316685622);

  @$pb.TagNumber(324829411)
  ColumnInfo get timeseriesmeasurevaluecolumninfo => $_getN(2);
  @$pb.TagNumber(324829411)
  set timeseriesmeasurevaluecolumninfo(ColumnInfo value) =>
      $_setField(324829411, value);
  @$pb.TagNumber(324829411)
  $core.bool hasTimeseriesmeasurevaluecolumninfo() => $_has(2);
  @$pb.TagNumber(324829411)
  void clearTimeseriesmeasurevaluecolumninfo() => $_clearField(324829411);
  @$pb.TagNumber(324829411)
  ColumnInfo ensureTimeseriesmeasurevaluecolumninfo() => $_ensure(2);

  @$pb.TagNumber(476744353)
  ColumnInfo get arraycolumninfo => $_getN(3);
  @$pb.TagNumber(476744353)
  set arraycolumninfo(ColumnInfo value) => $_setField(476744353, value);
  @$pb.TagNumber(476744353)
  $core.bool hasArraycolumninfo() => $_has(3);
  @$pb.TagNumber(476744353)
  void clearArraycolumninfo() => $_clearField(476744353);
  @$pb.TagNumber(476744353)
  ColumnInfo ensureArraycolumninfo() => $_ensure(3);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

class UpdateAccountSettingsRequest extends $pb.GeneratedMessage {
  factory UpdateAccountSettingsRequest({
    $core.int? maxquerytcu,
    QueryComputeRequest? querycompute,
    QueryPricingModel? querypricingmodel,
  }) {
    final result = create();
    if (maxquerytcu != null) result.maxquerytcu = maxquerytcu;
    if (querycompute != null) result.querycompute = querycompute;
    if (querypricingmodel != null) result.querypricingmodel = querypricingmodel;
    return result;
  }

  UpdateAccountSettingsRequest._();

  factory UpdateAccountSettingsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateAccountSettingsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAccountSettingsRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aI(285295370, _omitFieldNames ? '' : 'maxquerytcu')
    ..aOM<QueryComputeRequest>(358277613, _omitFieldNames ? '' : 'querycompute',
        subBuilder: QueryComputeRequest.create)
    ..aE<QueryPricingModel>(
        453472445, _omitFieldNames ? '' : 'querypricingmodel',
        enumValues: QueryPricingModel.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountSettingsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountSettingsRequest copyWith(
          void Function(UpdateAccountSettingsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateAccountSettingsRequest))
          as UpdateAccountSettingsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAccountSettingsRequest create() =>
      UpdateAccountSettingsRequest._();
  @$core.override
  UpdateAccountSettingsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateAccountSettingsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateAccountSettingsRequest>(create);
  static UpdateAccountSettingsRequest? _defaultInstance;

  @$pb.TagNumber(285295370)
  $core.int get maxquerytcu => $_getIZ(0);
  @$pb.TagNumber(285295370)
  set maxquerytcu($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(285295370)
  $core.bool hasMaxquerytcu() => $_has(0);
  @$pb.TagNumber(285295370)
  void clearMaxquerytcu() => $_clearField(285295370);

  @$pb.TagNumber(358277613)
  QueryComputeRequest get querycompute => $_getN(1);
  @$pb.TagNumber(358277613)
  set querycompute(QueryComputeRequest value) => $_setField(358277613, value);
  @$pb.TagNumber(358277613)
  $core.bool hasQuerycompute() => $_has(1);
  @$pb.TagNumber(358277613)
  void clearQuerycompute() => $_clearField(358277613);
  @$pb.TagNumber(358277613)
  QueryComputeRequest ensureQuerycompute() => $_ensure(1);

  @$pb.TagNumber(453472445)
  QueryPricingModel get querypricingmodel => $_getN(2);
  @$pb.TagNumber(453472445)
  set querypricingmodel(QueryPricingModel value) =>
      $_setField(453472445, value);
  @$pb.TagNumber(453472445)
  $core.bool hasQuerypricingmodel() => $_has(2);
  @$pb.TagNumber(453472445)
  void clearQuerypricingmodel() => $_clearField(453472445);
}

class UpdateAccountSettingsResponse extends $pb.GeneratedMessage {
  factory UpdateAccountSettingsResponse({
    $core.int? maxquerytcu,
    QueryComputeResponse? querycompute,
    QueryPricingModel? querypricingmodel,
  }) {
    final result = create();
    if (maxquerytcu != null) result.maxquerytcu = maxquerytcu;
    if (querycompute != null) result.querycompute = querycompute;
    if (querypricingmodel != null) result.querypricingmodel = querypricingmodel;
    return result;
  }

  UpdateAccountSettingsResponse._();

  factory UpdateAccountSettingsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateAccountSettingsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAccountSettingsResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aI(285295370, _omitFieldNames ? '' : 'maxquerytcu')
    ..aOM<QueryComputeResponse>(
        358277613, _omitFieldNames ? '' : 'querycompute',
        subBuilder: QueryComputeResponse.create)
    ..aE<QueryPricingModel>(
        453472445, _omitFieldNames ? '' : 'querypricingmodel',
        enumValues: QueryPricingModel.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountSettingsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountSettingsResponse copyWith(
          void Function(UpdateAccountSettingsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateAccountSettingsResponse))
          as UpdateAccountSettingsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAccountSettingsResponse create() =>
      UpdateAccountSettingsResponse._();
  @$core.override
  UpdateAccountSettingsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateAccountSettingsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateAccountSettingsResponse>(create);
  static UpdateAccountSettingsResponse? _defaultInstance;

  @$pb.TagNumber(285295370)
  $core.int get maxquerytcu => $_getIZ(0);
  @$pb.TagNumber(285295370)
  set maxquerytcu($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(285295370)
  $core.bool hasMaxquerytcu() => $_has(0);
  @$pb.TagNumber(285295370)
  void clearMaxquerytcu() => $_clearField(285295370);

  @$pb.TagNumber(358277613)
  QueryComputeResponse get querycompute => $_getN(1);
  @$pb.TagNumber(358277613)
  set querycompute(QueryComputeResponse value) => $_setField(358277613, value);
  @$pb.TagNumber(358277613)
  $core.bool hasQuerycompute() => $_has(1);
  @$pb.TagNumber(358277613)
  void clearQuerycompute() => $_clearField(358277613);
  @$pb.TagNumber(358277613)
  QueryComputeResponse ensureQuerycompute() => $_ensure(1);

  @$pb.TagNumber(453472445)
  QueryPricingModel get querypricingmodel => $_getN(2);
  @$pb.TagNumber(453472445)
  set querypricingmodel(QueryPricingModel value) =>
      $_setField(453472445, value);
  @$pb.TagNumber(453472445)
  $core.bool hasQuerypricingmodel() => $_has(2);
  @$pb.TagNumber(453472445)
  void clearQuerypricingmodel() => $_clearField(453472445);
}

class UpdateScheduledQueryRequest extends $pb.GeneratedMessage {
  factory UpdateScheduledQueryRequest({
    $core.String? scheduledqueryarn,
    ScheduledQueryState? state,
  }) {
    final result = create();
    if (scheduledqueryarn != null) result.scheduledqueryarn = scheduledqueryarn;
    if (state != null) result.state = state;
    return result;
  }

  UpdateScheduledQueryRequest._();

  factory UpdateScheduledQueryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateScheduledQueryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateScheduledQueryRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
      createEmptyInstance: create)
    ..aOS(234602964, _omitFieldNames ? '' : 'scheduledqueryarn')
    ..aE<ScheduledQueryState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ScheduledQueryState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateScheduledQueryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateScheduledQueryRequest copyWith(
          void Function(UpdateScheduledQueryRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateScheduledQueryRequest))
          as UpdateScheduledQueryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateScheduledQueryRequest create() =>
      UpdateScheduledQueryRequest._();
  @$core.override
  UpdateScheduledQueryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateScheduledQueryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateScheduledQueryRequest>(create);
  static UpdateScheduledQueryRequest? _defaultInstance;

  @$pb.TagNumber(234602964)
  $core.String get scheduledqueryarn => $_getSZ(0);
  @$pb.TagNumber(234602964)
  set scheduledqueryarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234602964)
  $core.bool hasScheduledqueryarn() => $_has(0);
  @$pb.TagNumber(234602964)
  void clearScheduledqueryarn() => $_clearField(234602964);

  @$pb.TagNumber(502047895)
  ScheduledQueryState get state => $_getN(1);
  @$pb.TagNumber(502047895)
  set state(ScheduledQueryState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(1);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'query.timestream'),
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

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
