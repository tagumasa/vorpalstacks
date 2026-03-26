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

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'states.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'states.pbenum.dart';

class ActivityAlreadyExists extends $pb.GeneratedMessage {
  factory ActivityAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ActivityAlreadyExists._();

  factory ActivityAlreadyExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityAlreadyExists copyWith(
          void Function(ActivityAlreadyExists) updates) =>
      super.copyWith((message) => updates(message as ActivityAlreadyExists))
          as ActivityAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityAlreadyExists create() => ActivityAlreadyExists._();
  @$core.override
  ActivityAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityAlreadyExists>(create);
  static ActivityAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ActivityDoesNotExist extends $pb.GeneratedMessage {
  factory ActivityDoesNotExist({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ActivityDoesNotExist._();

  factory ActivityDoesNotExist.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityDoesNotExist.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityDoesNotExist',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityDoesNotExist clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityDoesNotExist copyWith(void Function(ActivityDoesNotExist) updates) =>
      super.copyWith((message) => updates(message as ActivityDoesNotExist))
          as ActivityDoesNotExist;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityDoesNotExist create() => ActivityDoesNotExist._();
  @$core.override
  ActivityDoesNotExist createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityDoesNotExist getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityDoesNotExist>(create);
  static ActivityDoesNotExist? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ActivityFailedEventDetails extends $pb.GeneratedMessage {
  factory ActivityFailedEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  ActivityFailedEventDetails._();

  factory ActivityFailedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityFailedEventDetails copyWith(
          void Function(ActivityFailedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ActivityFailedEventDetails))
          as ActivityFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityFailedEventDetails create() => ActivityFailedEventDetails._();
  @$core.override
  ActivityFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityFailedEventDetails>(create);
  static ActivityFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class ActivityLimitExceeded extends $pb.GeneratedMessage {
  factory ActivityLimitExceeded({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ActivityLimitExceeded._();

  factory ActivityLimitExceeded.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityLimitExceeded.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityLimitExceeded',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityLimitExceeded clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityLimitExceeded copyWith(
          void Function(ActivityLimitExceeded) updates) =>
      super.copyWith((message) => updates(message as ActivityLimitExceeded))
          as ActivityLimitExceeded;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityLimitExceeded create() => ActivityLimitExceeded._();
  @$core.override
  ActivityLimitExceeded createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityLimitExceeded getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityLimitExceeded>(create);
  static ActivityLimitExceeded? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ActivityListItem extends $pb.GeneratedMessage {
  factory ActivityListItem({
    $core.String? name,
    $core.String? creationdate,
    $core.String? activityarn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (activityarn != null) result.activityarn = activityarn;
    return result;
  }

  ActivityListItem._();

  factory ActivityListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aOS(327279492, _omitFieldNames ? '' : 'activityarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityListItem copyWith(void Function(ActivityListItem) updates) =>
      super.copyWith((message) => updates(message as ActivityListItem))
          as ActivityListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityListItem create() => ActivityListItem._();
  @$core.override
  ActivityListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityListItem>(create);
  static ActivityListItem? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(1);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(1);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(327279492)
  $core.String get activityarn => $_getSZ(2);
  @$pb.TagNumber(327279492)
  set activityarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(327279492)
  $core.bool hasActivityarn() => $_has(2);
  @$pb.TagNumber(327279492)
  void clearActivityarn() => $_clearField(327279492);
}

class ActivityScheduleFailedEventDetails extends $pb.GeneratedMessage {
  factory ActivityScheduleFailedEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  ActivityScheduleFailedEventDetails._();

  factory ActivityScheduleFailedEventDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityScheduleFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityScheduleFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityScheduleFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityScheduleFailedEventDetails copyWith(
          void Function(ActivityScheduleFailedEventDetails) updates) =>
      super.copyWith((message) =>
              updates(message as ActivityScheduleFailedEventDetails))
          as ActivityScheduleFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityScheduleFailedEventDetails create() =>
      ActivityScheduleFailedEventDetails._();
  @$core.override
  ActivityScheduleFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityScheduleFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityScheduleFailedEventDetails>(
          create);
  static ActivityScheduleFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class ActivityScheduledEventDetails extends $pb.GeneratedMessage {
  factory ActivityScheduledEventDetails({
    $fixnum.Int64? heartbeatinseconds,
    $core.String? resource,
    $core.String? input,
    HistoryEventExecutionDataDetails? inputdetails,
    $fixnum.Int64? timeoutinseconds,
  }) {
    final result = create();
    if (heartbeatinseconds != null)
      result.heartbeatinseconds = heartbeatinseconds;
    if (resource != null) result.resource = resource;
    if (input != null) result.input = input;
    if (inputdetails != null) result.inputdetails = inputdetails;
    if (timeoutinseconds != null) result.timeoutinseconds = timeoutinseconds;
    return result;
  }

  ActivityScheduledEventDetails._();

  factory ActivityScheduledEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityScheduledEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityScheduledEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aInt64(125718754, _omitFieldNames ? '' : 'heartbeatinseconds')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..aOM<HistoryEventExecutionDataDetails>(
        452625788, _omitFieldNames ? '' : 'inputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..aInt64(472710197, _omitFieldNames ? '' : 'timeoutinseconds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityScheduledEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityScheduledEventDetails copyWith(
          void Function(ActivityScheduledEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ActivityScheduledEventDetails))
          as ActivityScheduledEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityScheduledEventDetails create() =>
      ActivityScheduledEventDetails._();
  @$core.override
  ActivityScheduledEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityScheduledEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityScheduledEventDetails>(create);
  static ActivityScheduledEventDetails? _defaultInstance;

  @$pb.TagNumber(125718754)
  $fixnum.Int64 get heartbeatinseconds => $_getI64(0);
  @$pb.TagNumber(125718754)
  set heartbeatinseconds($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(125718754)
  $core.bool hasHeartbeatinseconds() => $_has(0);
  @$pb.TagNumber(125718754)
  void clearHeartbeatinseconds() => $_clearField(125718754);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(1);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(1, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(1);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(2);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(2, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(2);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);

  @$pb.TagNumber(452625788)
  HistoryEventExecutionDataDetails get inputdetails => $_getN(3);
  @$pb.TagNumber(452625788)
  set inputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(452625788, value);
  @$pb.TagNumber(452625788)
  $core.bool hasInputdetails() => $_has(3);
  @$pb.TagNumber(452625788)
  void clearInputdetails() => $_clearField(452625788);
  @$pb.TagNumber(452625788)
  HistoryEventExecutionDataDetails ensureInputdetails() => $_ensure(3);

  @$pb.TagNumber(472710197)
  $fixnum.Int64 get timeoutinseconds => $_getI64(4);
  @$pb.TagNumber(472710197)
  set timeoutinseconds($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(472710197)
  $core.bool hasTimeoutinseconds() => $_has(4);
  @$pb.TagNumber(472710197)
  void clearTimeoutinseconds() => $_clearField(472710197);
}

class ActivityStartedEventDetails extends $pb.GeneratedMessage {
  factory ActivityStartedEventDetails({
    $core.String? workername,
  }) {
    final result = create();
    if (workername != null) result.workername = workername;
    return result;
  }

  ActivityStartedEventDetails._();

  factory ActivityStartedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityStartedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityStartedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(526761781, _omitFieldNames ? '' : 'workername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityStartedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityStartedEventDetails copyWith(
          void Function(ActivityStartedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ActivityStartedEventDetails))
          as ActivityStartedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityStartedEventDetails create() =>
      ActivityStartedEventDetails._();
  @$core.override
  ActivityStartedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityStartedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityStartedEventDetails>(create);
  static ActivityStartedEventDetails? _defaultInstance;

  @$pb.TagNumber(526761781)
  $core.String get workername => $_getSZ(0);
  @$pb.TagNumber(526761781)
  set workername($core.String value) => $_setString(0, value);
  @$pb.TagNumber(526761781)
  $core.bool hasWorkername() => $_has(0);
  @$pb.TagNumber(526761781)
  void clearWorkername() => $_clearField(526761781);
}

class ActivitySucceededEventDetails extends $pb.GeneratedMessage {
  factory ActivitySucceededEventDetails({
    HistoryEventExecutionDataDetails? outputdetails,
    $core.String? output,
  }) {
    final result = create();
    if (outputdetails != null) result.outputdetails = outputdetails;
    if (output != null) result.output = output;
    return result;
  }

  ActivitySucceededEventDetails._();

  factory ActivitySucceededEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivitySucceededEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivitySucceededEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOM<HistoryEventExecutionDataDetails>(
        393734643, _omitFieldNames ? '' : 'outputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivitySucceededEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivitySucceededEventDetails copyWith(
          void Function(ActivitySucceededEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ActivitySucceededEventDetails))
          as ActivitySucceededEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivitySucceededEventDetails create() =>
      ActivitySucceededEventDetails._();
  @$core.override
  ActivitySucceededEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivitySucceededEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivitySucceededEventDetails>(create);
  static ActivitySucceededEventDetails? _defaultInstance;

  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails get outputdetails => $_getN(0);
  @$pb.TagNumber(393734643)
  set outputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(393734643, value);
  @$pb.TagNumber(393734643)
  $core.bool hasOutputdetails() => $_has(0);
  @$pb.TagNumber(393734643)
  void clearOutputdetails() => $_clearField(393734643);
  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails ensureOutputdetails() => $_ensure(0);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(1);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(1, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(1);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);
}

class ActivityTimedOutEventDetails extends $pb.GeneratedMessage {
  factory ActivityTimedOutEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  ActivityTimedOutEventDetails._();

  factory ActivityTimedOutEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityTimedOutEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityTimedOutEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityTimedOutEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityTimedOutEventDetails copyWith(
          void Function(ActivityTimedOutEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ActivityTimedOutEventDetails))
          as ActivityTimedOutEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityTimedOutEventDetails create() =>
      ActivityTimedOutEventDetails._();
  @$core.override
  ActivityTimedOutEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityTimedOutEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityTimedOutEventDetails>(create);
  static ActivityTimedOutEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class ActivityWorkerLimitExceeded extends $pb.GeneratedMessage {
  factory ActivityWorkerLimitExceeded({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ActivityWorkerLimitExceeded._();

  factory ActivityWorkerLimitExceeded.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivityWorkerLimitExceeded.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivityWorkerLimitExceeded',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityWorkerLimitExceeded clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivityWorkerLimitExceeded copyWith(
          void Function(ActivityWorkerLimitExceeded) updates) =>
      super.copyWith(
              (message) => updates(message as ActivityWorkerLimitExceeded))
          as ActivityWorkerLimitExceeded;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivityWorkerLimitExceeded create() =>
      ActivityWorkerLimitExceeded._();
  @$core.override
  ActivityWorkerLimitExceeded createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivityWorkerLimitExceeded getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivityWorkerLimitExceeded>(create);
  static ActivityWorkerLimitExceeded? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class AssignedVariablesDetails extends $pb.GeneratedMessage {
  factory AssignedVariablesDetails({
    $core.bool? truncated,
  }) {
    final result = create();
    if (truncated != null) result.truncated = truncated;
    return result;
  }

  AssignedVariablesDetails._();

  factory AssignedVariablesDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssignedVariablesDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssignedVariablesDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOB(407490858, _omitFieldNames ? '' : 'truncated')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssignedVariablesDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssignedVariablesDetails copyWith(
          void Function(AssignedVariablesDetails) updates) =>
      super.copyWith((message) => updates(message as AssignedVariablesDetails))
          as AssignedVariablesDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssignedVariablesDetails create() => AssignedVariablesDetails._();
  @$core.override
  AssignedVariablesDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssignedVariablesDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssignedVariablesDetails>(create);
  static AssignedVariablesDetails? _defaultInstance;

  @$pb.TagNumber(407490858)
  $core.bool get truncated => $_getBF(0);
  @$pb.TagNumber(407490858)
  set truncated($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(407490858)
  $core.bool hasTruncated() => $_has(0);
  @$pb.TagNumber(407490858)
  void clearTruncated() => $_clearField(407490858);
}

class BillingDetails extends $pb.GeneratedMessage {
  factory BillingDetails({
    $fixnum.Int64? billeddurationinmilliseconds,
    $fixnum.Int64? billedmemoryusedinmb,
  }) {
    final result = create();
    if (billeddurationinmilliseconds != null)
      result.billeddurationinmilliseconds = billeddurationinmilliseconds;
    if (billedmemoryusedinmb != null)
      result.billedmemoryusedinmb = billedmemoryusedinmb;
    return result;
  }

  BillingDetails._();

  factory BillingDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BillingDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aInt64(65131593, _omitFieldNames ? '' : 'billeddurationinmilliseconds')
    ..aInt64(459756940, _omitFieldNames ? '' : 'billedmemoryusedinmb')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingDetails copyWith(void Function(BillingDetails) updates) =>
      super.copyWith((message) => updates(message as BillingDetails))
          as BillingDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BillingDetails create() => BillingDetails._();
  @$core.override
  BillingDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BillingDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingDetails>(create);
  static BillingDetails? _defaultInstance;

  @$pb.TagNumber(65131593)
  $fixnum.Int64 get billeddurationinmilliseconds => $_getI64(0);
  @$pb.TagNumber(65131593)
  set billeddurationinmilliseconds($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(65131593)
  $core.bool hasBilleddurationinmilliseconds() => $_has(0);
  @$pb.TagNumber(65131593)
  void clearBilleddurationinmilliseconds() => $_clearField(65131593);

  @$pb.TagNumber(459756940)
  $fixnum.Int64 get billedmemoryusedinmb => $_getI64(1);
  @$pb.TagNumber(459756940)
  set billedmemoryusedinmb($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(459756940)
  $core.bool hasBilledmemoryusedinmb() => $_has(1);
  @$pb.TagNumber(459756940)
  void clearBilledmemoryusedinmb() => $_clearField(459756940);
}

class CloudWatchEventsExecutionDataDetails extends $pb.GeneratedMessage {
  factory CloudWatchEventsExecutionDataDetails({
    $core.bool? included,
  }) {
    final result = create();
    if (included != null) result.included = included;
    return result;
  }

  CloudWatchEventsExecutionDataDetails._();

  factory CloudWatchEventsExecutionDataDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudWatchEventsExecutionDataDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudWatchEventsExecutionDataDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOB(517673658, _omitFieldNames ? '' : 'included')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchEventsExecutionDataDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchEventsExecutionDataDetails copyWith(
          void Function(CloudWatchEventsExecutionDataDetails) updates) =>
      super.copyWith((message) =>
              updates(message as CloudWatchEventsExecutionDataDetails))
          as CloudWatchEventsExecutionDataDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudWatchEventsExecutionDataDetails create() =>
      CloudWatchEventsExecutionDataDetails._();
  @$core.override
  CloudWatchEventsExecutionDataDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudWatchEventsExecutionDataDetails getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CloudWatchEventsExecutionDataDetails>(create);
  static CloudWatchEventsExecutionDataDetails? _defaultInstance;

  @$pb.TagNumber(517673658)
  $core.bool get included => $_getBF(0);
  @$pb.TagNumber(517673658)
  set included($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(517673658)
  $core.bool hasIncluded() => $_has(0);
  @$pb.TagNumber(517673658)
  void clearIncluded() => $_clearField(517673658);
}

class CloudWatchLogsLogGroup extends $pb.GeneratedMessage {
  factory CloudWatchLogsLogGroup({
    $core.String? loggrouparn,
  }) {
    final result = create();
    if (loggrouparn != null) result.loggrouparn = loggrouparn;
    return result;
  }

  CloudWatchLogsLogGroup._();

  factory CloudWatchLogsLogGroup.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudWatchLogsLogGroup.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudWatchLogsLogGroup',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(6742512, _omitFieldNames ? '' : 'loggrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchLogsLogGroup clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchLogsLogGroup copyWith(
          void Function(CloudWatchLogsLogGroup) updates) =>
      super.copyWith((message) => updates(message as CloudWatchLogsLogGroup))
          as CloudWatchLogsLogGroup;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudWatchLogsLogGroup create() => CloudWatchLogsLogGroup._();
  @$core.override
  CloudWatchLogsLogGroup createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudWatchLogsLogGroup getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloudWatchLogsLogGroup>(create);
  static CloudWatchLogsLogGroup? _defaultInstance;

  @$pb.TagNumber(6742512)
  $core.String get loggrouparn => $_getSZ(0);
  @$pb.TagNumber(6742512)
  set loggrouparn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6742512)
  $core.bool hasLoggrouparn() => $_has(0);
  @$pb.TagNumber(6742512)
  void clearLoggrouparn() => $_clearField(6742512);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CreateActivityInput extends $pb.GeneratedMessage {
  factory CreateActivityInput({
    EncryptionConfiguration? encryptionconfiguration,
    $core.String? name,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  CreateActivityInput._();

  factory CreateActivityInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateActivityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateActivityInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOM<EncryptionConfiguration>(
        167857431, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(337046433, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateActivityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateActivityInput copyWith(void Function(CreateActivityInput) updates) =>
      super.copyWith((message) => updates(message as CreateActivityInput))
          as CreateActivityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateActivityInput create() => CreateActivityInput._();
  @$core.override
  CreateActivityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateActivityInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateActivityInput>(create);
  static CreateActivityInput? _defaultInstance;

  @$pb.TagNumber(167857431)
  EncryptionConfiguration get encryptionconfiguration => $_getN(0);
  @$pb.TagNumber(167857431)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(167857431, value);
  @$pb.TagNumber(167857431)
  $core.bool hasEncryptionconfiguration() => $_has(0);
  @$pb.TagNumber(167857431)
  void clearEncryptionconfiguration() => $_clearField(167857431);
  @$pb.TagNumber(167857431)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(0);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(337046433)
  $pb.PbList<Tag> get tags => $_getList(2);
}

class CreateActivityOutput extends $pb.GeneratedMessage {
  factory CreateActivityOutput({
    $core.String? creationdate,
    $core.String? activityarn,
  }) {
    final result = create();
    if (creationdate != null) result.creationdate = creationdate;
    if (activityarn != null) result.activityarn = activityarn;
    return result;
  }

  CreateActivityOutput._();

  factory CreateActivityOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateActivityOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateActivityOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aOS(327279492, _omitFieldNames ? '' : 'activityarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateActivityOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateActivityOutput copyWith(void Function(CreateActivityOutput) updates) =>
      super.copyWith((message) => updates(message as CreateActivityOutput))
          as CreateActivityOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateActivityOutput create() => CreateActivityOutput._();
  @$core.override
  CreateActivityOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateActivityOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateActivityOutput>(create);
  static CreateActivityOutput? _defaultInstance;

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(0);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(0);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(327279492)
  $core.String get activityarn => $_getSZ(1);
  @$pb.TagNumber(327279492)
  set activityarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(327279492)
  $core.bool hasActivityarn() => $_has(1);
  @$pb.TagNumber(327279492)
  void clearActivityarn() => $_clearField(327279492);
}

class CreateStateMachineAliasInput extends $pb.GeneratedMessage {
  factory CreateStateMachineAliasInput({
    $core.String? name,
    $core.String? description,
    $core.Iterable<RoutingConfigurationListItem>? routingconfiguration,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    if (routingconfiguration != null)
      result.routingconfiguration.addAll(routingconfiguration);
    return result;
  }

  CreateStateMachineAliasInput._();

  factory CreateStateMachineAliasInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateStateMachineAliasInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateStateMachineAliasInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..pPM<RoutingConfigurationListItem>(
        372891510, _omitFieldNames ? '' : 'routingconfiguration',
        subBuilder: RoutingConfigurationListItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStateMachineAliasInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStateMachineAliasInput copyWith(
          void Function(CreateStateMachineAliasInput) updates) =>
      super.copyWith(
              (message) => updates(message as CreateStateMachineAliasInput))
          as CreateStateMachineAliasInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateStateMachineAliasInput create() =>
      CreateStateMachineAliasInput._();
  @$core.override
  CreateStateMachineAliasInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateStateMachineAliasInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateStateMachineAliasInput>(create);
  static CreateStateMachineAliasInput? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(372891510)
  $pb.PbList<RoutingConfigurationListItem> get routingconfiguration =>
      $_getList(2);
}

class CreateStateMachineAliasOutput extends $pb.GeneratedMessage {
  factory CreateStateMachineAliasOutput({
    $core.String? creationdate,
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (creationdate != null) result.creationdate = creationdate;
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  CreateStateMachineAliasOutput._();

  factory CreateStateMachineAliasOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateStateMachineAliasOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateStateMachineAliasOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStateMachineAliasOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStateMachineAliasOutput copyWith(
          void Function(CreateStateMachineAliasOutput) updates) =>
      super.copyWith(
              (message) => updates(message as CreateStateMachineAliasOutput))
          as CreateStateMachineAliasOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateStateMachineAliasOutput create() =>
      CreateStateMachineAliasOutput._();
  @$core.override
  CreateStateMachineAliasOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateStateMachineAliasOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateStateMachineAliasOutput>(create);
  static CreateStateMachineAliasOutput? _defaultInstance;

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(0);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(0);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(1);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(1);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class CreateStateMachineInput extends $pb.GeneratedMessage {
  factory CreateStateMachineInput({
    $core.String? definition,
    EncryptionConfiguration? encryptionconfiguration,
    $core.String? rolearn,
    $core.String? name,
    $core.bool? publish,
    StateMachineType? type,
    $core.Iterable<Tag>? tags,
    LoggingConfiguration? loggingconfiguration,
    $core.String? versiondescription,
    TracingConfiguration? tracingconfiguration,
  }) {
    final result = create();
    if (definition != null) result.definition = definition;
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (rolearn != null) result.rolearn = rolearn;
    if (name != null) result.name = name;
    if (publish != null) result.publish = publish;
    if (type != null) result.type = type;
    if (tags != null) result.tags.addAll(tags);
    if (loggingconfiguration != null)
      result.loggingconfiguration = loggingconfiguration;
    if (versiondescription != null)
      result.versiondescription = versiondescription;
    if (tracingconfiguration != null)
      result.tracingconfiguration = tracingconfiguration;
    return result;
  }

  CreateStateMachineInput._();

  factory CreateStateMachineInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateStateMachineInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateStateMachineInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(68443297, _omitFieldNames ? '' : 'definition')
    ..aOM<EncryptionConfiguration>(
        167857431, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..aOS(170019745, _omitFieldNames ? '' : 'rolearn')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOB(264247305, _omitFieldNames ? '' : 'publish')
    ..aE<StateMachineType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: StateMachineType.values)
    ..pPM<Tag>(337046433, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOM<LoggingConfiguration>(
        420811605, _omitFieldNames ? '' : 'loggingconfiguration',
        subBuilder: LoggingConfiguration.create)
    ..aOS(434714300, _omitFieldNames ? '' : 'versiondescription')
    ..aOM<TracingConfiguration>(
        491315910, _omitFieldNames ? '' : 'tracingconfiguration',
        subBuilder: TracingConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStateMachineInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStateMachineInput copyWith(
          void Function(CreateStateMachineInput) updates) =>
      super.copyWith((message) => updates(message as CreateStateMachineInput))
          as CreateStateMachineInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateStateMachineInput create() => CreateStateMachineInput._();
  @$core.override
  CreateStateMachineInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateStateMachineInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateStateMachineInput>(create);
  static CreateStateMachineInput? _defaultInstance;

  @$pb.TagNumber(68443297)
  $core.String get definition => $_getSZ(0);
  @$pb.TagNumber(68443297)
  set definition($core.String value) => $_setString(0, value);
  @$pb.TagNumber(68443297)
  $core.bool hasDefinition() => $_has(0);
  @$pb.TagNumber(68443297)
  void clearDefinition() => $_clearField(68443297);

  @$pb.TagNumber(167857431)
  EncryptionConfiguration get encryptionconfiguration => $_getN(1);
  @$pb.TagNumber(167857431)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(167857431, value);
  @$pb.TagNumber(167857431)
  $core.bool hasEncryptionconfiguration() => $_has(1);
  @$pb.TagNumber(167857431)
  void clearEncryptionconfiguration() => $_clearField(167857431);
  @$pb.TagNumber(167857431)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(1);

  @$pb.TagNumber(170019745)
  $core.String get rolearn => $_getSZ(2);
  @$pb.TagNumber(170019745)
  set rolearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(170019745)
  $core.bool hasRolearn() => $_has(2);
  @$pb.TagNumber(170019745)
  void clearRolearn() => $_clearField(170019745);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(264247305)
  $core.bool get publish => $_getBF(4);
  @$pb.TagNumber(264247305)
  set publish($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(264247305)
  $core.bool hasPublish() => $_has(4);
  @$pb.TagNumber(264247305)
  void clearPublish() => $_clearField(264247305);

  @$pb.TagNumber(287830350)
  StateMachineType get type => $_getN(5);
  @$pb.TagNumber(287830350)
  set type(StateMachineType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(5);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(337046433)
  $pb.PbList<Tag> get tags => $_getList(6);

  @$pb.TagNumber(420811605)
  LoggingConfiguration get loggingconfiguration => $_getN(7);
  @$pb.TagNumber(420811605)
  set loggingconfiguration(LoggingConfiguration value) =>
      $_setField(420811605, value);
  @$pb.TagNumber(420811605)
  $core.bool hasLoggingconfiguration() => $_has(7);
  @$pb.TagNumber(420811605)
  void clearLoggingconfiguration() => $_clearField(420811605);
  @$pb.TagNumber(420811605)
  LoggingConfiguration ensureLoggingconfiguration() => $_ensure(7);

  @$pb.TagNumber(434714300)
  $core.String get versiondescription => $_getSZ(8);
  @$pb.TagNumber(434714300)
  set versiondescription($core.String value) => $_setString(8, value);
  @$pb.TagNumber(434714300)
  $core.bool hasVersiondescription() => $_has(8);
  @$pb.TagNumber(434714300)
  void clearVersiondescription() => $_clearField(434714300);

  @$pb.TagNumber(491315910)
  TracingConfiguration get tracingconfiguration => $_getN(9);
  @$pb.TagNumber(491315910)
  set tracingconfiguration(TracingConfiguration value) =>
      $_setField(491315910, value);
  @$pb.TagNumber(491315910)
  $core.bool hasTracingconfiguration() => $_has(9);
  @$pb.TagNumber(491315910)
  void clearTracingconfiguration() => $_clearField(491315910);
  @$pb.TagNumber(491315910)
  TracingConfiguration ensureTracingconfiguration() => $_ensure(9);
}

class CreateStateMachineOutput extends $pb.GeneratedMessage {
  factory CreateStateMachineOutput({
    $core.String? statemachineversionarn,
    $core.String? creationdate,
    $core.String? statemachinearn,
  }) {
    final result = create();
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    if (creationdate != null) result.creationdate = creationdate;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    return result;
  }

  CreateStateMachineOutput._();

  factory CreateStateMachineOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateStateMachineOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateStateMachineOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStateMachineOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStateMachineOutput copyWith(
          void Function(CreateStateMachineOutput) updates) =>
      super.copyWith((message) => updates(message as CreateStateMachineOutput))
          as CreateStateMachineOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateStateMachineOutput create() => CreateStateMachineOutput._();
  @$core.override
  CreateStateMachineOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateStateMachineOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateStateMachineOutput>(create);
  static CreateStateMachineOutput? _defaultInstance;

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(0);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(0);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(1);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(1);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(2);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(2);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);
}

class DeleteActivityInput extends $pb.GeneratedMessage {
  factory DeleteActivityInput({
    $core.String? activityarn,
  }) {
    final result = create();
    if (activityarn != null) result.activityarn = activityarn;
    return result;
  }

  DeleteActivityInput._();

  factory DeleteActivityInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteActivityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteActivityInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(327279492, _omitFieldNames ? '' : 'activityarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteActivityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteActivityInput copyWith(void Function(DeleteActivityInput) updates) =>
      super.copyWith((message) => updates(message as DeleteActivityInput))
          as DeleteActivityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteActivityInput create() => DeleteActivityInput._();
  @$core.override
  DeleteActivityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteActivityInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteActivityInput>(create);
  static DeleteActivityInput? _defaultInstance;

  @$pb.TagNumber(327279492)
  $core.String get activityarn => $_getSZ(0);
  @$pb.TagNumber(327279492)
  set activityarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(327279492)
  $core.bool hasActivityarn() => $_has(0);
  @$pb.TagNumber(327279492)
  void clearActivityarn() => $_clearField(327279492);
}

class DeleteActivityOutput extends $pb.GeneratedMessage {
  factory DeleteActivityOutput() => create();

  DeleteActivityOutput._();

  factory DeleteActivityOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteActivityOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteActivityOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteActivityOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteActivityOutput copyWith(void Function(DeleteActivityOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteActivityOutput))
          as DeleteActivityOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteActivityOutput create() => DeleteActivityOutput._();
  @$core.override
  DeleteActivityOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteActivityOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteActivityOutput>(create);
  static DeleteActivityOutput? _defaultInstance;
}

class DeleteStateMachineAliasInput extends $pb.GeneratedMessage {
  factory DeleteStateMachineAliasInput({
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  DeleteStateMachineAliasInput._();

  factory DeleteStateMachineAliasInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteStateMachineAliasInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteStateMachineAliasInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineAliasInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineAliasInput copyWith(
          void Function(DeleteStateMachineAliasInput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteStateMachineAliasInput))
          as DeleteStateMachineAliasInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineAliasInput create() =>
      DeleteStateMachineAliasInput._();
  @$core.override
  DeleteStateMachineAliasInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineAliasInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteStateMachineAliasInput>(create);
  static DeleteStateMachineAliasInput? _defaultInstance;

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(0);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(0);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class DeleteStateMachineAliasOutput extends $pb.GeneratedMessage {
  factory DeleteStateMachineAliasOutput() => create();

  DeleteStateMachineAliasOutput._();

  factory DeleteStateMachineAliasOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteStateMachineAliasOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteStateMachineAliasOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineAliasOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineAliasOutput copyWith(
          void Function(DeleteStateMachineAliasOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteStateMachineAliasOutput))
          as DeleteStateMachineAliasOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineAliasOutput create() =>
      DeleteStateMachineAliasOutput._();
  @$core.override
  DeleteStateMachineAliasOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineAliasOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteStateMachineAliasOutput>(create);
  static DeleteStateMachineAliasOutput? _defaultInstance;
}

class DeleteStateMachineInput extends $pb.GeneratedMessage {
  factory DeleteStateMachineInput({
    $core.String? statemachinearn,
  }) {
    final result = create();
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    return result;
  }

  DeleteStateMachineInput._();

  factory DeleteStateMachineInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteStateMachineInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteStateMachineInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineInput copyWith(
          void Function(DeleteStateMachineInput) updates) =>
      super.copyWith((message) => updates(message as DeleteStateMachineInput))
          as DeleteStateMachineInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineInput create() => DeleteStateMachineInput._();
  @$core.override
  DeleteStateMachineInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteStateMachineInput>(create);
  static DeleteStateMachineInput? _defaultInstance;

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(0);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(0);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);
}

class DeleteStateMachineOutput extends $pb.GeneratedMessage {
  factory DeleteStateMachineOutput() => create();

  DeleteStateMachineOutput._();

  factory DeleteStateMachineOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteStateMachineOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteStateMachineOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineOutput copyWith(
          void Function(DeleteStateMachineOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteStateMachineOutput))
          as DeleteStateMachineOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineOutput create() => DeleteStateMachineOutput._();
  @$core.override
  DeleteStateMachineOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteStateMachineOutput>(create);
  static DeleteStateMachineOutput? _defaultInstance;
}

class DeleteStateMachineVersionInput extends $pb.GeneratedMessage {
  factory DeleteStateMachineVersionInput({
    $core.String? statemachineversionarn,
  }) {
    final result = create();
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    return result;
  }

  DeleteStateMachineVersionInput._();

  factory DeleteStateMachineVersionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteStateMachineVersionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteStateMachineVersionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineVersionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineVersionInput copyWith(
          void Function(DeleteStateMachineVersionInput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteStateMachineVersionInput))
          as DeleteStateMachineVersionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineVersionInput create() =>
      DeleteStateMachineVersionInput._();
  @$core.override
  DeleteStateMachineVersionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineVersionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteStateMachineVersionInput>(create);
  static DeleteStateMachineVersionInput? _defaultInstance;

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(0);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(0);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);
}

class DeleteStateMachineVersionOutput extends $pb.GeneratedMessage {
  factory DeleteStateMachineVersionOutput() => create();

  DeleteStateMachineVersionOutput._();

  factory DeleteStateMachineVersionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteStateMachineVersionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteStateMachineVersionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineVersionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStateMachineVersionOutput copyWith(
          void Function(DeleteStateMachineVersionOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteStateMachineVersionOutput))
          as DeleteStateMachineVersionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineVersionOutput create() =>
      DeleteStateMachineVersionOutput._();
  @$core.override
  DeleteStateMachineVersionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteStateMachineVersionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteStateMachineVersionOutput>(
          create);
  static DeleteStateMachineVersionOutput? _defaultInstance;
}

class DescribeActivityInput extends $pb.GeneratedMessage {
  factory DescribeActivityInput({
    $core.String? activityarn,
  }) {
    final result = create();
    if (activityarn != null) result.activityarn = activityarn;
    return result;
  }

  DescribeActivityInput._();

  factory DescribeActivityInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeActivityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeActivityInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(327279492, _omitFieldNames ? '' : 'activityarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeActivityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeActivityInput copyWith(
          void Function(DescribeActivityInput) updates) =>
      super.copyWith((message) => updates(message as DescribeActivityInput))
          as DescribeActivityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeActivityInput create() => DescribeActivityInput._();
  @$core.override
  DescribeActivityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeActivityInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeActivityInput>(create);
  static DescribeActivityInput? _defaultInstance;

  @$pb.TagNumber(327279492)
  $core.String get activityarn => $_getSZ(0);
  @$pb.TagNumber(327279492)
  set activityarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(327279492)
  $core.bool hasActivityarn() => $_has(0);
  @$pb.TagNumber(327279492)
  void clearActivityarn() => $_clearField(327279492);
}

class DescribeActivityOutput extends $pb.GeneratedMessage {
  factory DescribeActivityOutput({
    EncryptionConfiguration? encryptionconfiguration,
    $core.String? name,
    $core.String? creationdate,
    $core.String? activityarn,
  }) {
    final result = create();
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (activityarn != null) result.activityarn = activityarn;
    return result;
  }

  DescribeActivityOutput._();

  factory DescribeActivityOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeActivityOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeActivityOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOM<EncryptionConfiguration>(
        167857431, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aOS(327279492, _omitFieldNames ? '' : 'activityarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeActivityOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeActivityOutput copyWith(
          void Function(DescribeActivityOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeActivityOutput))
          as DescribeActivityOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeActivityOutput create() => DescribeActivityOutput._();
  @$core.override
  DescribeActivityOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeActivityOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeActivityOutput>(create);
  static DescribeActivityOutput? _defaultInstance;

  @$pb.TagNumber(167857431)
  EncryptionConfiguration get encryptionconfiguration => $_getN(0);
  @$pb.TagNumber(167857431)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(167857431, value);
  @$pb.TagNumber(167857431)
  $core.bool hasEncryptionconfiguration() => $_has(0);
  @$pb.TagNumber(167857431)
  void clearEncryptionconfiguration() => $_clearField(167857431);
  @$pb.TagNumber(167857431)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(0);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(2);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(2);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(327279492)
  $core.String get activityarn => $_getSZ(3);
  @$pb.TagNumber(327279492)
  set activityarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(327279492)
  $core.bool hasActivityarn() => $_has(3);
  @$pb.TagNumber(327279492)
  void clearActivityarn() => $_clearField(327279492);
}

class DescribeExecutionInput extends $pb.GeneratedMessage {
  factory DescribeExecutionInput({
    IncludedData? includeddata,
    $core.String? executionarn,
  }) {
    final result = create();
    if (includeddata != null) result.includeddata = includeddata;
    if (executionarn != null) result.executionarn = executionarn;
    return result;
  }

  DescribeExecutionInput._();

  factory DescribeExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aE<IncludedData>(109719114, _omitFieldNames ? '' : 'includeddata',
        enumValues: IncludedData.values)
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeExecutionInput copyWith(
          void Function(DescribeExecutionInput) updates) =>
      super.copyWith((message) => updates(message as DescribeExecutionInput))
          as DescribeExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeExecutionInput create() => DescribeExecutionInput._();
  @$core.override
  DescribeExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeExecutionInput>(create);
  static DescribeExecutionInput? _defaultInstance;

  @$pb.TagNumber(109719114)
  IncludedData get includeddata => $_getN(0);
  @$pb.TagNumber(109719114)
  set includeddata(IncludedData value) => $_setField(109719114, value);
  @$pb.TagNumber(109719114)
  $core.bool hasIncludeddata() => $_has(0);
  @$pb.TagNumber(109719114)
  void clearIncludeddata() => $_clearField(109719114);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(1);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(1);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);
}

class DescribeExecutionOutput extends $pb.GeneratedMessage {
  factory DescribeExecutionOutput({
    $core.String? maprunarn,
    $core.String? error,
    $core.String? statemachineversionarn,
    $core.String? cause,
    $core.String? redrivedate,
    $core.String? stopdate,
    $core.String? traceheader,
    $core.String? name,
    ExecutionRedriveStatus? redrivestatus,
    $core.String? executionarn,
    $core.String? redrivestatusreason,
    $core.String? startdate,
    $core.String? statemachinearn,
    CloudWatchEventsExecutionDataDetails? outputdetails,
    $core.String? output,
    $core.String? input,
    ExecutionStatus? status,
    CloudWatchEventsExecutionDataDetails? inputdetails,
    $core.int? redrivecount,
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    if (error != null) result.error = error;
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    if (cause != null) result.cause = cause;
    if (redrivedate != null) result.redrivedate = redrivedate;
    if (stopdate != null) result.stopdate = stopdate;
    if (traceheader != null) result.traceheader = traceheader;
    if (name != null) result.name = name;
    if (redrivestatus != null) result.redrivestatus = redrivestatus;
    if (executionarn != null) result.executionarn = executionarn;
    if (redrivestatusreason != null)
      result.redrivestatusreason = redrivestatusreason;
    if (startdate != null) result.startdate = startdate;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (outputdetails != null) result.outputdetails = outputdetails;
    if (output != null) result.output = output;
    if (input != null) result.input = input;
    if (status != null) result.status = status;
    if (inputdetails != null) result.inputdetails = inputdetails;
    if (redrivecount != null) result.redrivecount = redrivecount;
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  DescribeExecutionOutput._();

  factory DescribeExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(152812125, _omitFieldNames ? '' : 'redrivedate')
    ..aOS(180697434, _omitFieldNames ? '' : 'stopdate')
    ..aOS(219960864, _omitFieldNames ? '' : 'traceheader')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aE<ExecutionRedriveStatus>(
        247102059, _omitFieldNames ? '' : 'redrivestatus',
        enumValues: ExecutionRedriveStatus.values)
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..aOS(339085215, _omitFieldNames ? '' : 'redrivestatusreason')
    ..aOS(364840732, _omitFieldNames ? '' : 'startdate')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aOM<CloudWatchEventsExecutionDataDetails>(
        393734643, _omitFieldNames ? '' : 'outputdetails',
        subBuilder: CloudWatchEventsExecutionDataDetails.create)
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..aE<ExecutionStatus>(441153520, _omitFieldNames ? '' : 'status',
        enumValues: ExecutionStatus.values)
    ..aOM<CloudWatchEventsExecutionDataDetails>(
        452625788, _omitFieldNames ? '' : 'inputdetails',
        subBuilder: CloudWatchEventsExecutionDataDetails.create)
    ..aI(473458696, _omitFieldNames ? '' : 'redrivecount')
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeExecutionOutput copyWith(
          void Function(DescribeExecutionOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeExecutionOutput))
          as DescribeExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeExecutionOutput create() => DescribeExecutionOutput._();
  @$core.override
  DescribeExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeExecutionOutput>(create);
  static DescribeExecutionOutput? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(1);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(1, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(1);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(2);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(2);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(3);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(3, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(3);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(152812125)
  $core.String get redrivedate => $_getSZ(4);
  @$pb.TagNumber(152812125)
  set redrivedate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(152812125)
  $core.bool hasRedrivedate() => $_has(4);
  @$pb.TagNumber(152812125)
  void clearRedrivedate() => $_clearField(152812125);

  @$pb.TagNumber(180697434)
  $core.String get stopdate => $_getSZ(5);
  @$pb.TagNumber(180697434)
  set stopdate($core.String value) => $_setString(5, value);
  @$pb.TagNumber(180697434)
  $core.bool hasStopdate() => $_has(5);
  @$pb.TagNumber(180697434)
  void clearStopdate() => $_clearField(180697434);

  @$pb.TagNumber(219960864)
  $core.String get traceheader => $_getSZ(6);
  @$pb.TagNumber(219960864)
  set traceheader($core.String value) => $_setString(6, value);
  @$pb.TagNumber(219960864)
  $core.bool hasTraceheader() => $_has(6);
  @$pb.TagNumber(219960864)
  void clearTraceheader() => $_clearField(219960864);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(7);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(7, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(7);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(247102059)
  ExecutionRedriveStatus get redrivestatus => $_getN(8);
  @$pb.TagNumber(247102059)
  set redrivestatus(ExecutionRedriveStatus value) =>
      $_setField(247102059, value);
  @$pb.TagNumber(247102059)
  $core.bool hasRedrivestatus() => $_has(8);
  @$pb.TagNumber(247102059)
  void clearRedrivestatus() => $_clearField(247102059);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(9);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(9);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);

  @$pb.TagNumber(339085215)
  $core.String get redrivestatusreason => $_getSZ(10);
  @$pb.TagNumber(339085215)
  set redrivestatusreason($core.String value) => $_setString(10, value);
  @$pb.TagNumber(339085215)
  $core.bool hasRedrivestatusreason() => $_has(10);
  @$pb.TagNumber(339085215)
  void clearRedrivestatusreason() => $_clearField(339085215);

  @$pb.TagNumber(364840732)
  $core.String get startdate => $_getSZ(11);
  @$pb.TagNumber(364840732)
  set startdate($core.String value) => $_setString(11, value);
  @$pb.TagNumber(364840732)
  $core.bool hasStartdate() => $_has(11);
  @$pb.TagNumber(364840732)
  void clearStartdate() => $_clearField(364840732);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(12);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(12, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(12);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(393734643)
  CloudWatchEventsExecutionDataDetails get outputdetails => $_getN(13);
  @$pb.TagNumber(393734643)
  set outputdetails(CloudWatchEventsExecutionDataDetails value) =>
      $_setField(393734643, value);
  @$pb.TagNumber(393734643)
  $core.bool hasOutputdetails() => $_has(13);
  @$pb.TagNumber(393734643)
  void clearOutputdetails() => $_clearField(393734643);
  @$pb.TagNumber(393734643)
  CloudWatchEventsExecutionDataDetails ensureOutputdetails() => $_ensure(13);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(14);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(14, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(14);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(15);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(15, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(15);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);

  @$pb.TagNumber(441153520)
  ExecutionStatus get status => $_getN(16);
  @$pb.TagNumber(441153520)
  set status(ExecutionStatus value) => $_setField(441153520, value);
  @$pb.TagNumber(441153520)
  $core.bool hasStatus() => $_has(16);
  @$pb.TagNumber(441153520)
  void clearStatus() => $_clearField(441153520);

  @$pb.TagNumber(452625788)
  CloudWatchEventsExecutionDataDetails get inputdetails => $_getN(17);
  @$pb.TagNumber(452625788)
  set inputdetails(CloudWatchEventsExecutionDataDetails value) =>
      $_setField(452625788, value);
  @$pb.TagNumber(452625788)
  $core.bool hasInputdetails() => $_has(17);
  @$pb.TagNumber(452625788)
  void clearInputdetails() => $_clearField(452625788);
  @$pb.TagNumber(452625788)
  CloudWatchEventsExecutionDataDetails ensureInputdetails() => $_ensure(17);

  @$pb.TagNumber(473458696)
  $core.int get redrivecount => $_getIZ(18);
  @$pb.TagNumber(473458696)
  set redrivecount($core.int value) => $_setSignedInt32(18, value);
  @$pb.TagNumber(473458696)
  $core.bool hasRedrivecount() => $_has(18);
  @$pb.TagNumber(473458696)
  void clearRedrivecount() => $_clearField(473458696);

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(19);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(19, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(19);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class DescribeMapRunInput extends $pb.GeneratedMessage {
  factory DescribeMapRunInput({
    $core.String? maprunarn,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    return result;
  }

  DescribeMapRunInput._();

  factory DescribeMapRunInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeMapRunInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeMapRunInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeMapRunInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeMapRunInput copyWith(void Function(DescribeMapRunInput) updates) =>
      super.copyWith((message) => updates(message as DescribeMapRunInput))
          as DescribeMapRunInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeMapRunInput create() => DescribeMapRunInput._();
  @$core.override
  DescribeMapRunInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeMapRunInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeMapRunInput>(create);
  static DescribeMapRunInput? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);
}

class DescribeMapRunOutput extends $pb.GeneratedMessage {
  factory DescribeMapRunOutput({
    $core.String? maprunarn,
    $fixnum.Int64? toleratedfailurecount,
    $core.int? maxconcurrency,
    $core.double? toleratedfailurepercentage,
    $core.String? redrivedate,
    $core.String? stopdate,
    $core.String? executionarn,
    MapRunExecutionCounts? executioncounts,
    MapRunItemCounts? itemcounts,
    $core.String? startdate,
    MapRunStatus? status,
    $core.int? redrivecount,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    if (toleratedfailurecount != null)
      result.toleratedfailurecount = toleratedfailurecount;
    if (maxconcurrency != null) result.maxconcurrency = maxconcurrency;
    if (toleratedfailurepercentage != null)
      result.toleratedfailurepercentage = toleratedfailurepercentage;
    if (redrivedate != null) result.redrivedate = redrivedate;
    if (stopdate != null) result.stopdate = stopdate;
    if (executionarn != null) result.executionarn = executionarn;
    if (executioncounts != null) result.executioncounts = executioncounts;
    if (itemcounts != null) result.itemcounts = itemcounts;
    if (startdate != null) result.startdate = startdate;
    if (status != null) result.status = status;
    if (redrivecount != null) result.redrivecount = redrivecount;
    return result;
  }

  DescribeMapRunOutput._();

  factory DescribeMapRunOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeMapRunOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeMapRunOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..aInt64(41834811, _omitFieldNames ? '' : 'toleratedfailurecount')
    ..aI(100901405, _omitFieldNames ? '' : 'maxconcurrency')
    ..aD(116496164, _omitFieldNames ? '' : 'toleratedfailurepercentage',
        fieldType: $pb.PbFieldType.OF)
    ..aOS(152812125, _omitFieldNames ? '' : 'redrivedate')
    ..aOS(180697434, _omitFieldNames ? '' : 'stopdate')
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..aOM<MapRunExecutionCounts>(
        317260758, _omitFieldNames ? '' : 'executioncounts',
        subBuilder: MapRunExecutionCounts.create)
    ..aOM<MapRunItemCounts>(334548339, _omitFieldNames ? '' : 'itemcounts',
        subBuilder: MapRunItemCounts.create)
    ..aOS(364840732, _omitFieldNames ? '' : 'startdate')
    ..aE<MapRunStatus>(441153520, _omitFieldNames ? '' : 'status',
        enumValues: MapRunStatus.values)
    ..aI(473458696, _omitFieldNames ? '' : 'redrivecount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeMapRunOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeMapRunOutput copyWith(void Function(DescribeMapRunOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeMapRunOutput))
          as DescribeMapRunOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeMapRunOutput create() => DescribeMapRunOutput._();
  @$core.override
  DescribeMapRunOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeMapRunOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeMapRunOutput>(create);
  static DescribeMapRunOutput? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);

  @$pb.TagNumber(41834811)
  $fixnum.Int64 get toleratedfailurecount => $_getI64(1);
  @$pb.TagNumber(41834811)
  set toleratedfailurecount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(41834811)
  $core.bool hasToleratedfailurecount() => $_has(1);
  @$pb.TagNumber(41834811)
  void clearToleratedfailurecount() => $_clearField(41834811);

  @$pb.TagNumber(100901405)
  $core.int get maxconcurrency => $_getIZ(2);
  @$pb.TagNumber(100901405)
  set maxconcurrency($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(100901405)
  $core.bool hasMaxconcurrency() => $_has(2);
  @$pb.TagNumber(100901405)
  void clearMaxconcurrency() => $_clearField(100901405);

  @$pb.TagNumber(116496164)
  $core.double get toleratedfailurepercentage => $_getN(3);
  @$pb.TagNumber(116496164)
  set toleratedfailurepercentage($core.double value) => $_setFloat(3, value);
  @$pb.TagNumber(116496164)
  $core.bool hasToleratedfailurepercentage() => $_has(3);
  @$pb.TagNumber(116496164)
  void clearToleratedfailurepercentage() => $_clearField(116496164);

  @$pb.TagNumber(152812125)
  $core.String get redrivedate => $_getSZ(4);
  @$pb.TagNumber(152812125)
  set redrivedate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(152812125)
  $core.bool hasRedrivedate() => $_has(4);
  @$pb.TagNumber(152812125)
  void clearRedrivedate() => $_clearField(152812125);

  @$pb.TagNumber(180697434)
  $core.String get stopdate => $_getSZ(5);
  @$pb.TagNumber(180697434)
  set stopdate($core.String value) => $_setString(5, value);
  @$pb.TagNumber(180697434)
  $core.bool hasStopdate() => $_has(5);
  @$pb.TagNumber(180697434)
  void clearStopdate() => $_clearField(180697434);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(6);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(6);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);

  @$pb.TagNumber(317260758)
  MapRunExecutionCounts get executioncounts => $_getN(7);
  @$pb.TagNumber(317260758)
  set executioncounts(MapRunExecutionCounts value) =>
      $_setField(317260758, value);
  @$pb.TagNumber(317260758)
  $core.bool hasExecutioncounts() => $_has(7);
  @$pb.TagNumber(317260758)
  void clearExecutioncounts() => $_clearField(317260758);
  @$pb.TagNumber(317260758)
  MapRunExecutionCounts ensureExecutioncounts() => $_ensure(7);

  @$pb.TagNumber(334548339)
  MapRunItemCounts get itemcounts => $_getN(8);
  @$pb.TagNumber(334548339)
  set itemcounts(MapRunItemCounts value) => $_setField(334548339, value);
  @$pb.TagNumber(334548339)
  $core.bool hasItemcounts() => $_has(8);
  @$pb.TagNumber(334548339)
  void clearItemcounts() => $_clearField(334548339);
  @$pb.TagNumber(334548339)
  MapRunItemCounts ensureItemcounts() => $_ensure(8);

  @$pb.TagNumber(364840732)
  $core.String get startdate => $_getSZ(9);
  @$pb.TagNumber(364840732)
  set startdate($core.String value) => $_setString(9, value);
  @$pb.TagNumber(364840732)
  $core.bool hasStartdate() => $_has(9);
  @$pb.TagNumber(364840732)
  void clearStartdate() => $_clearField(364840732);

  @$pb.TagNumber(441153520)
  MapRunStatus get status => $_getN(10);
  @$pb.TagNumber(441153520)
  set status(MapRunStatus value) => $_setField(441153520, value);
  @$pb.TagNumber(441153520)
  $core.bool hasStatus() => $_has(10);
  @$pb.TagNumber(441153520)
  void clearStatus() => $_clearField(441153520);

  @$pb.TagNumber(473458696)
  $core.int get redrivecount => $_getIZ(11);
  @$pb.TagNumber(473458696)
  set redrivecount($core.int value) => $_setSignedInt32(11, value);
  @$pb.TagNumber(473458696)
  $core.bool hasRedrivecount() => $_has(11);
  @$pb.TagNumber(473458696)
  void clearRedrivecount() => $_clearField(473458696);
}

class DescribeStateMachineAliasInput extends $pb.GeneratedMessage {
  factory DescribeStateMachineAliasInput({
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  DescribeStateMachineAliasInput._();

  factory DescribeStateMachineAliasInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStateMachineAliasInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStateMachineAliasInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineAliasInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineAliasInput copyWith(
          void Function(DescribeStateMachineAliasInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeStateMachineAliasInput))
          as DescribeStateMachineAliasInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineAliasInput create() =>
      DescribeStateMachineAliasInput._();
  @$core.override
  DescribeStateMachineAliasInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineAliasInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStateMachineAliasInput>(create);
  static DescribeStateMachineAliasInput? _defaultInstance;

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(0);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(0);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class DescribeStateMachineAliasOutput extends $pb.GeneratedMessage {
  factory DescribeStateMachineAliasOutput({
    $core.String? name,
    $core.String? creationdate,
    $core.String? description,
    $core.Iterable<RoutingConfigurationListItem>? routingconfiguration,
    $core.String? updatedate,
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (description != null) result.description = description;
    if (routingconfiguration != null)
      result.routingconfiguration.addAll(routingconfiguration);
    if (updatedate != null) result.updatedate = updatedate;
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  DescribeStateMachineAliasOutput._();

  factory DescribeStateMachineAliasOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStateMachineAliasOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStateMachineAliasOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..pPM<RoutingConfigurationListItem>(
        372891510, _omitFieldNames ? '' : 'routingconfiguration',
        subBuilder: RoutingConfigurationListItem.create)
    ..aOS(510552561, _omitFieldNames ? '' : 'updatedate')
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineAliasOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineAliasOutput copyWith(
          void Function(DescribeStateMachineAliasOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeStateMachineAliasOutput))
          as DescribeStateMachineAliasOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineAliasOutput create() =>
      DescribeStateMachineAliasOutput._();
  @$core.override
  DescribeStateMachineAliasOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineAliasOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStateMachineAliasOutput>(
          create);
  static DescribeStateMachineAliasOutput? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(1);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(1);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(372891510)
  $pb.PbList<RoutingConfigurationListItem> get routingconfiguration =>
      $_getList(3);

  @$pb.TagNumber(510552561)
  $core.String get updatedate => $_getSZ(4);
  @$pb.TagNumber(510552561)
  set updatedate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(510552561)
  $core.bool hasUpdatedate() => $_has(4);
  @$pb.TagNumber(510552561)
  void clearUpdatedate() => $_clearField(510552561);

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(5);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(5);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class DescribeStateMachineForExecutionInput extends $pb.GeneratedMessage {
  factory DescribeStateMachineForExecutionInput({
    IncludedData? includeddata,
    $core.String? executionarn,
  }) {
    final result = create();
    if (includeddata != null) result.includeddata = includeddata;
    if (executionarn != null) result.executionarn = executionarn;
    return result;
  }

  DescribeStateMachineForExecutionInput._();

  factory DescribeStateMachineForExecutionInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStateMachineForExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStateMachineForExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aE<IncludedData>(109719114, _omitFieldNames ? '' : 'includeddata',
        enumValues: IncludedData.values)
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineForExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineForExecutionInput copyWith(
          void Function(DescribeStateMachineForExecutionInput) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeStateMachineForExecutionInput))
          as DescribeStateMachineForExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineForExecutionInput create() =>
      DescribeStateMachineForExecutionInput._();
  @$core.override
  DescribeStateMachineForExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineForExecutionInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DescribeStateMachineForExecutionInput>(create);
  static DescribeStateMachineForExecutionInput? _defaultInstance;

  @$pb.TagNumber(109719114)
  IncludedData get includeddata => $_getN(0);
  @$pb.TagNumber(109719114)
  set includeddata(IncludedData value) => $_setField(109719114, value);
  @$pb.TagNumber(109719114)
  $core.bool hasIncludeddata() => $_has(0);
  @$pb.TagNumber(109719114)
  void clearIncludeddata() => $_clearField(109719114);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(1);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(1);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);
}

class DescribeStateMachineForExecutionOutput extends $pb.GeneratedMessage {
  factory DescribeStateMachineForExecutionOutput({
    $core.String? maprunarn,
    $core.String? definition,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        variablereferences,
    EncryptionConfiguration? encryptionconfiguration,
    $core.String? rolearn,
    $core.String? name,
    $core.String? revisionid,
    $core.String? label,
    $core.String? statemachinearn,
    LoggingConfiguration? loggingconfiguration,
    TracingConfiguration? tracingconfiguration,
    $core.String? updatedate,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    if (definition != null) result.definition = definition;
    if (variablereferences != null)
      result.variablereferences.addEntries(variablereferences);
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (rolearn != null) result.rolearn = rolearn;
    if (name != null) result.name = name;
    if (revisionid != null) result.revisionid = revisionid;
    if (label != null) result.label = label;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (loggingconfiguration != null)
      result.loggingconfiguration = loggingconfiguration;
    if (tracingconfiguration != null)
      result.tracingconfiguration = tracingconfiguration;
    if (updatedate != null) result.updatedate = updatedate;
    return result;
  }

  DescribeStateMachineForExecutionOutput._();

  factory DescribeStateMachineForExecutionOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStateMachineForExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStateMachineForExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..aOS(68443297, _omitFieldNames ? '' : 'definition')
    ..m<$core.String, $core.String>(
        150942252, _omitFieldNames ? '' : 'variablereferences',
        entryClassName:
            'DescribeStateMachineForExecutionOutput.VariablereferencesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('states'))
    ..aOM<EncryptionConfiguration>(
        167857431, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..aOS(170019745, _omitFieldNames ? '' : 'rolearn')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(369170086, _omitFieldNames ? '' : 'revisionid')
    ..aOS(379000830, _omitFieldNames ? '' : 'label')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aOM<LoggingConfiguration>(
        420811605, _omitFieldNames ? '' : 'loggingconfiguration',
        subBuilder: LoggingConfiguration.create)
    ..aOM<TracingConfiguration>(
        491315910, _omitFieldNames ? '' : 'tracingconfiguration',
        subBuilder: TracingConfiguration.create)
    ..aOS(510552561, _omitFieldNames ? '' : 'updatedate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineForExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineForExecutionOutput copyWith(
          void Function(DescribeStateMachineForExecutionOutput) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeStateMachineForExecutionOutput))
          as DescribeStateMachineForExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineForExecutionOutput create() =>
      DescribeStateMachineForExecutionOutput._();
  @$core.override
  DescribeStateMachineForExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineForExecutionOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DescribeStateMachineForExecutionOutput>(create);
  static DescribeStateMachineForExecutionOutput? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);

  @$pb.TagNumber(68443297)
  $core.String get definition => $_getSZ(1);
  @$pb.TagNumber(68443297)
  set definition($core.String value) => $_setString(1, value);
  @$pb.TagNumber(68443297)
  $core.bool hasDefinition() => $_has(1);
  @$pb.TagNumber(68443297)
  void clearDefinition() => $_clearField(68443297);

  @$pb.TagNumber(150942252)
  $pb.PbMap<$core.String, $core.String> get variablereferences => $_getMap(2);

  @$pb.TagNumber(167857431)
  EncryptionConfiguration get encryptionconfiguration => $_getN(3);
  @$pb.TagNumber(167857431)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(167857431, value);
  @$pb.TagNumber(167857431)
  $core.bool hasEncryptionconfiguration() => $_has(3);
  @$pb.TagNumber(167857431)
  void clearEncryptionconfiguration() => $_clearField(167857431);
  @$pb.TagNumber(167857431)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(3);

  @$pb.TagNumber(170019745)
  $core.String get rolearn => $_getSZ(4);
  @$pb.TagNumber(170019745)
  set rolearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(170019745)
  $core.bool hasRolearn() => $_has(4);
  @$pb.TagNumber(170019745)
  void clearRolearn() => $_clearField(170019745);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(369170086)
  $core.String get revisionid => $_getSZ(6);
  @$pb.TagNumber(369170086)
  set revisionid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(369170086)
  $core.bool hasRevisionid() => $_has(6);
  @$pb.TagNumber(369170086)
  void clearRevisionid() => $_clearField(369170086);

  @$pb.TagNumber(379000830)
  $core.String get label => $_getSZ(7);
  @$pb.TagNumber(379000830)
  set label($core.String value) => $_setString(7, value);
  @$pb.TagNumber(379000830)
  $core.bool hasLabel() => $_has(7);
  @$pb.TagNumber(379000830)
  void clearLabel() => $_clearField(379000830);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(8);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(8);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(420811605)
  LoggingConfiguration get loggingconfiguration => $_getN(9);
  @$pb.TagNumber(420811605)
  set loggingconfiguration(LoggingConfiguration value) =>
      $_setField(420811605, value);
  @$pb.TagNumber(420811605)
  $core.bool hasLoggingconfiguration() => $_has(9);
  @$pb.TagNumber(420811605)
  void clearLoggingconfiguration() => $_clearField(420811605);
  @$pb.TagNumber(420811605)
  LoggingConfiguration ensureLoggingconfiguration() => $_ensure(9);

  @$pb.TagNumber(491315910)
  TracingConfiguration get tracingconfiguration => $_getN(10);
  @$pb.TagNumber(491315910)
  set tracingconfiguration(TracingConfiguration value) =>
      $_setField(491315910, value);
  @$pb.TagNumber(491315910)
  $core.bool hasTracingconfiguration() => $_has(10);
  @$pb.TagNumber(491315910)
  void clearTracingconfiguration() => $_clearField(491315910);
  @$pb.TagNumber(491315910)
  TracingConfiguration ensureTracingconfiguration() => $_ensure(10);

  @$pb.TagNumber(510552561)
  $core.String get updatedate => $_getSZ(11);
  @$pb.TagNumber(510552561)
  set updatedate($core.String value) => $_setString(11, value);
  @$pb.TagNumber(510552561)
  $core.bool hasUpdatedate() => $_has(11);
  @$pb.TagNumber(510552561)
  void clearUpdatedate() => $_clearField(510552561);
}

class DescribeStateMachineInput extends $pb.GeneratedMessage {
  factory DescribeStateMachineInput({
    IncludedData? includeddata,
    $core.String? statemachinearn,
  }) {
    final result = create();
    if (includeddata != null) result.includeddata = includeddata;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    return result;
  }

  DescribeStateMachineInput._();

  factory DescribeStateMachineInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStateMachineInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStateMachineInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aE<IncludedData>(109719114, _omitFieldNames ? '' : 'includeddata',
        enumValues: IncludedData.values)
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineInput copyWith(
          void Function(DescribeStateMachineInput) updates) =>
      super.copyWith((message) => updates(message as DescribeStateMachineInput))
          as DescribeStateMachineInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineInput create() => DescribeStateMachineInput._();
  @$core.override
  DescribeStateMachineInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStateMachineInput>(create);
  static DescribeStateMachineInput? _defaultInstance;

  @$pb.TagNumber(109719114)
  IncludedData get includeddata => $_getN(0);
  @$pb.TagNumber(109719114)
  set includeddata(IncludedData value) => $_setField(109719114, value);
  @$pb.TagNumber(109719114)
  $core.bool hasIncludeddata() => $_has(0);
  @$pb.TagNumber(109719114)
  void clearIncludeddata() => $_clearField(109719114);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(1);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(1);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);
}

class DescribeStateMachineOutput extends $pb.GeneratedMessage {
  factory DescribeStateMachineOutput({
    $core.String? definition,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        variablereferences,
    EncryptionConfiguration? encryptionconfiguration,
    $core.String? rolearn,
    $core.String? name,
    $core.String? creationdate,
    StateMachineType? type,
    $core.String? description,
    $core.String? revisionid,
    $core.String? label,
    $core.String? statemachinearn,
    LoggingConfiguration? loggingconfiguration,
    StateMachineStatus? status,
    TracingConfiguration? tracingconfiguration,
  }) {
    final result = create();
    if (definition != null) result.definition = definition;
    if (variablereferences != null)
      result.variablereferences.addEntries(variablereferences);
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (rolearn != null) result.rolearn = rolearn;
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (type != null) result.type = type;
    if (description != null) result.description = description;
    if (revisionid != null) result.revisionid = revisionid;
    if (label != null) result.label = label;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (loggingconfiguration != null)
      result.loggingconfiguration = loggingconfiguration;
    if (status != null) result.status = status;
    if (tracingconfiguration != null)
      result.tracingconfiguration = tracingconfiguration;
    return result;
  }

  DescribeStateMachineOutput._();

  factory DescribeStateMachineOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeStateMachineOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeStateMachineOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(68443297, _omitFieldNames ? '' : 'definition')
    ..m<$core.String, $core.String>(
        150942252, _omitFieldNames ? '' : 'variablereferences',
        entryClassName: 'DescribeStateMachineOutput.VariablereferencesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('states'))
    ..aOM<EncryptionConfiguration>(
        167857431, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..aOS(170019745, _omitFieldNames ? '' : 'rolearn')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aE<StateMachineType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: StateMachineType.values)
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(369170086, _omitFieldNames ? '' : 'revisionid')
    ..aOS(379000830, _omitFieldNames ? '' : 'label')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aOM<LoggingConfiguration>(
        420811605, _omitFieldNames ? '' : 'loggingconfiguration',
        subBuilder: LoggingConfiguration.create)
    ..aE<StateMachineStatus>(441153520, _omitFieldNames ? '' : 'status',
        enumValues: StateMachineStatus.values)
    ..aOM<TracingConfiguration>(
        491315910, _omitFieldNames ? '' : 'tracingconfiguration',
        subBuilder: TracingConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeStateMachineOutput copyWith(
          void Function(DescribeStateMachineOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeStateMachineOutput))
          as DescribeStateMachineOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineOutput create() => DescribeStateMachineOutput._();
  @$core.override
  DescribeStateMachineOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeStateMachineOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeStateMachineOutput>(create);
  static DescribeStateMachineOutput? _defaultInstance;

  @$pb.TagNumber(68443297)
  $core.String get definition => $_getSZ(0);
  @$pb.TagNumber(68443297)
  set definition($core.String value) => $_setString(0, value);
  @$pb.TagNumber(68443297)
  $core.bool hasDefinition() => $_has(0);
  @$pb.TagNumber(68443297)
  void clearDefinition() => $_clearField(68443297);

  @$pb.TagNumber(150942252)
  $pb.PbMap<$core.String, $core.String> get variablereferences => $_getMap(1);

  @$pb.TagNumber(167857431)
  EncryptionConfiguration get encryptionconfiguration => $_getN(2);
  @$pb.TagNumber(167857431)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(167857431, value);
  @$pb.TagNumber(167857431)
  $core.bool hasEncryptionconfiguration() => $_has(2);
  @$pb.TagNumber(167857431)
  void clearEncryptionconfiguration() => $_clearField(167857431);
  @$pb.TagNumber(167857431)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(2);

  @$pb.TagNumber(170019745)
  $core.String get rolearn => $_getSZ(3);
  @$pb.TagNumber(170019745)
  set rolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(170019745)
  $core.bool hasRolearn() => $_has(3);
  @$pb.TagNumber(170019745)
  void clearRolearn() => $_clearField(170019745);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(5);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(5, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(5);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(287830350)
  StateMachineType get type => $_getN(6);
  @$pb.TagNumber(287830350)
  set type(StateMachineType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(6);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(7);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(7, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(7);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(369170086)
  $core.String get revisionid => $_getSZ(8);
  @$pb.TagNumber(369170086)
  set revisionid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(369170086)
  $core.bool hasRevisionid() => $_has(8);
  @$pb.TagNumber(369170086)
  void clearRevisionid() => $_clearField(369170086);

  @$pb.TagNumber(379000830)
  $core.String get label => $_getSZ(9);
  @$pb.TagNumber(379000830)
  set label($core.String value) => $_setString(9, value);
  @$pb.TagNumber(379000830)
  $core.bool hasLabel() => $_has(9);
  @$pb.TagNumber(379000830)
  void clearLabel() => $_clearField(379000830);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(10);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(10, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(10);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(420811605)
  LoggingConfiguration get loggingconfiguration => $_getN(11);
  @$pb.TagNumber(420811605)
  set loggingconfiguration(LoggingConfiguration value) =>
      $_setField(420811605, value);
  @$pb.TagNumber(420811605)
  $core.bool hasLoggingconfiguration() => $_has(11);
  @$pb.TagNumber(420811605)
  void clearLoggingconfiguration() => $_clearField(420811605);
  @$pb.TagNumber(420811605)
  LoggingConfiguration ensureLoggingconfiguration() => $_ensure(11);

  @$pb.TagNumber(441153520)
  StateMachineStatus get status => $_getN(12);
  @$pb.TagNumber(441153520)
  set status(StateMachineStatus value) => $_setField(441153520, value);
  @$pb.TagNumber(441153520)
  $core.bool hasStatus() => $_has(12);
  @$pb.TagNumber(441153520)
  void clearStatus() => $_clearField(441153520);

  @$pb.TagNumber(491315910)
  TracingConfiguration get tracingconfiguration => $_getN(13);
  @$pb.TagNumber(491315910)
  set tracingconfiguration(TracingConfiguration value) =>
      $_setField(491315910, value);
  @$pb.TagNumber(491315910)
  $core.bool hasTracingconfiguration() => $_has(13);
  @$pb.TagNumber(491315910)
  void clearTracingconfiguration() => $_clearField(491315910);
  @$pb.TagNumber(491315910)
  TracingConfiguration ensureTracingconfiguration() => $_ensure(13);
}

class EncryptionConfiguration extends $pb.GeneratedMessage {
  factory EncryptionConfiguration({
    EncryptionType? type,
    $core.int? kmsdatakeyreuseperiodseconds,
    $core.String? kmskeyid,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (kmsdatakeyreuseperiodseconds != null)
      result.kmsdatakeyreuseperiodseconds = kmsdatakeyreuseperiodseconds;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    return result;
  }

  EncryptionConfiguration._();

  factory EncryptionConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EncryptionConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EncryptionConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aE<EncryptionType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: EncryptionType.values)
    ..aI(440747764, _omitFieldNames ? '' : 'kmsdatakeyreuseperiodseconds')
    ..aOS(510698477, _omitFieldNames ? '' : 'kmskeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptionConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptionConfiguration copyWith(
          void Function(EncryptionConfiguration) updates) =>
      super.copyWith((message) => updates(message as EncryptionConfiguration))
          as EncryptionConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EncryptionConfiguration create() => EncryptionConfiguration._();
  @$core.override
  EncryptionConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EncryptionConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EncryptionConfiguration>(create);
  static EncryptionConfiguration? _defaultInstance;

  @$pb.TagNumber(287830350)
  EncryptionType get type => $_getN(0);
  @$pb.TagNumber(287830350)
  set type(EncryptionType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(440747764)
  $core.int get kmsdatakeyreuseperiodseconds => $_getIZ(1);
  @$pb.TagNumber(440747764)
  set kmsdatakeyreuseperiodseconds($core.int value) =>
      $_setSignedInt32(1, value);
  @$pb.TagNumber(440747764)
  $core.bool hasKmsdatakeyreuseperiodseconds() => $_has(1);
  @$pb.TagNumber(440747764)
  void clearKmsdatakeyreuseperiodseconds() => $_clearField(440747764);

  @$pb.TagNumber(510698477)
  $core.String get kmskeyid => $_getSZ(2);
  @$pb.TagNumber(510698477)
  set kmskeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(510698477)
  $core.bool hasKmskeyid() => $_has(2);
  @$pb.TagNumber(510698477)
  void clearKmskeyid() => $_clearField(510698477);
}

class EvaluationFailedEventDetails extends $pb.GeneratedMessage {
  factory EvaluationFailedEventDetails({
    $core.String? error,
    $core.String? cause,
    $core.String? location,
    $core.String? state,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    if (location != null) result.location = location;
    if (state != null) result.state = state;
    return result;
  }

  EvaluationFailedEventDetails._();

  factory EvaluationFailedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EvaluationFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EvaluationFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(200649127, _omitFieldNames ? '' : 'location')
    ..aOS(405877495, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EvaluationFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EvaluationFailedEventDetails copyWith(
          void Function(EvaluationFailedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as EvaluationFailedEventDetails))
          as EvaluationFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EvaluationFailedEventDetails create() =>
      EvaluationFailedEventDetails._();
  @$core.override
  EvaluationFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EvaluationFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EvaluationFailedEventDetails>(create);
  static EvaluationFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(200649127)
  $core.String get location => $_getSZ(2);
  @$pb.TagNumber(200649127)
  set location($core.String value) => $_setString(2, value);
  @$pb.TagNumber(200649127)
  $core.bool hasLocation() => $_has(2);
  @$pb.TagNumber(200649127)
  void clearLocation() => $_clearField(200649127);

  @$pb.TagNumber(405877495)
  $core.String get state => $_getSZ(3);
  @$pb.TagNumber(405877495)
  set state($core.String value) => $_setString(3, value);
  @$pb.TagNumber(405877495)
  $core.bool hasState() => $_has(3);
  @$pb.TagNumber(405877495)
  void clearState() => $_clearField(405877495);
}

class ExecutionAbortedEventDetails extends $pb.GeneratedMessage {
  factory ExecutionAbortedEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  ExecutionAbortedEventDetails._();

  factory ExecutionAbortedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionAbortedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionAbortedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionAbortedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionAbortedEventDetails copyWith(
          void Function(ExecutionAbortedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ExecutionAbortedEventDetails))
          as ExecutionAbortedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionAbortedEventDetails create() =>
      ExecutionAbortedEventDetails._();
  @$core.override
  ExecutionAbortedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionAbortedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionAbortedEventDetails>(create);
  static ExecutionAbortedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class ExecutionAlreadyExists extends $pb.GeneratedMessage {
  factory ExecutionAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExecutionAlreadyExists._();

  factory ExecutionAlreadyExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionAlreadyExists copyWith(
          void Function(ExecutionAlreadyExists) updates) =>
      super.copyWith((message) => updates(message as ExecutionAlreadyExists))
          as ExecutionAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionAlreadyExists create() => ExecutionAlreadyExists._();
  @$core.override
  ExecutionAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionAlreadyExists>(create);
  static ExecutionAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExecutionDoesNotExist extends $pb.GeneratedMessage {
  factory ExecutionDoesNotExist({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExecutionDoesNotExist._();

  factory ExecutionDoesNotExist.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionDoesNotExist.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionDoesNotExist',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionDoesNotExist clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionDoesNotExist copyWith(
          void Function(ExecutionDoesNotExist) updates) =>
      super.copyWith((message) => updates(message as ExecutionDoesNotExist))
          as ExecutionDoesNotExist;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionDoesNotExist create() => ExecutionDoesNotExist._();
  @$core.override
  ExecutionDoesNotExist createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionDoesNotExist getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionDoesNotExist>(create);
  static ExecutionDoesNotExist? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExecutionFailedEventDetails extends $pb.GeneratedMessage {
  factory ExecutionFailedEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  ExecutionFailedEventDetails._();

  factory ExecutionFailedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionFailedEventDetails copyWith(
          void Function(ExecutionFailedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ExecutionFailedEventDetails))
          as ExecutionFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionFailedEventDetails create() =>
      ExecutionFailedEventDetails._();
  @$core.override
  ExecutionFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionFailedEventDetails>(create);
  static ExecutionFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class ExecutionLimitExceeded extends $pb.GeneratedMessage {
  factory ExecutionLimitExceeded({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExecutionLimitExceeded._();

  factory ExecutionLimitExceeded.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionLimitExceeded.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionLimitExceeded',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionLimitExceeded clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionLimitExceeded copyWith(
          void Function(ExecutionLimitExceeded) updates) =>
      super.copyWith((message) => updates(message as ExecutionLimitExceeded))
          as ExecutionLimitExceeded;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionLimitExceeded create() => ExecutionLimitExceeded._();
  @$core.override
  ExecutionLimitExceeded createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionLimitExceeded getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionLimitExceeded>(create);
  static ExecutionLimitExceeded? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExecutionListItem extends $pb.GeneratedMessage {
  factory ExecutionListItem({
    $core.String? maprunarn,
    $core.String? statemachineversionarn,
    $core.String? redrivedate,
    $core.String? stopdate,
    $core.String? name,
    $core.String? executionarn,
    $core.int? itemcount,
    $core.String? startdate,
    $core.String? statemachinearn,
    ExecutionStatus? status,
    $core.int? redrivecount,
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    if (redrivedate != null) result.redrivedate = redrivedate;
    if (stopdate != null) result.stopdate = stopdate;
    if (name != null) result.name = name;
    if (executionarn != null) result.executionarn = executionarn;
    if (itemcount != null) result.itemcount = itemcount;
    if (startdate != null) result.startdate = startdate;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (status != null) result.status = status;
    if (redrivecount != null) result.redrivecount = redrivecount;
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  ExecutionListItem._();

  factory ExecutionListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..aOS(152812125, _omitFieldNames ? '' : 'redrivedate')
    ..aOS(180697434, _omitFieldNames ? '' : 'stopdate')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..aI(349613174, _omitFieldNames ? '' : 'itemcount')
    ..aOS(364840732, _omitFieldNames ? '' : 'startdate')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aE<ExecutionStatus>(441153520, _omitFieldNames ? '' : 'status',
        enumValues: ExecutionStatus.values)
    ..aI(473458696, _omitFieldNames ? '' : 'redrivecount')
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionListItem copyWith(void Function(ExecutionListItem) updates) =>
      super.copyWith((message) => updates(message as ExecutionListItem))
          as ExecutionListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionListItem create() => ExecutionListItem._();
  @$core.override
  ExecutionListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionListItem>(create);
  static ExecutionListItem? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(1);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(1);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);

  @$pb.TagNumber(152812125)
  $core.String get redrivedate => $_getSZ(2);
  @$pb.TagNumber(152812125)
  set redrivedate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(152812125)
  $core.bool hasRedrivedate() => $_has(2);
  @$pb.TagNumber(152812125)
  void clearRedrivedate() => $_clearField(152812125);

  @$pb.TagNumber(180697434)
  $core.String get stopdate => $_getSZ(3);
  @$pb.TagNumber(180697434)
  set stopdate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(180697434)
  $core.bool hasStopdate() => $_has(3);
  @$pb.TagNumber(180697434)
  void clearStopdate() => $_clearField(180697434);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(5);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(5);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);

  @$pb.TagNumber(349613174)
  $core.int get itemcount => $_getIZ(6);
  @$pb.TagNumber(349613174)
  set itemcount($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(349613174)
  $core.bool hasItemcount() => $_has(6);
  @$pb.TagNumber(349613174)
  void clearItemcount() => $_clearField(349613174);

  @$pb.TagNumber(364840732)
  $core.String get startdate => $_getSZ(7);
  @$pb.TagNumber(364840732)
  set startdate($core.String value) => $_setString(7, value);
  @$pb.TagNumber(364840732)
  $core.bool hasStartdate() => $_has(7);
  @$pb.TagNumber(364840732)
  void clearStartdate() => $_clearField(364840732);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(8);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(8);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(441153520)
  ExecutionStatus get status => $_getN(9);
  @$pb.TagNumber(441153520)
  set status(ExecutionStatus value) => $_setField(441153520, value);
  @$pb.TagNumber(441153520)
  $core.bool hasStatus() => $_has(9);
  @$pb.TagNumber(441153520)
  void clearStatus() => $_clearField(441153520);

  @$pb.TagNumber(473458696)
  $core.int get redrivecount => $_getIZ(10);
  @$pb.TagNumber(473458696)
  set redrivecount($core.int value) => $_setSignedInt32(10, value);
  @$pb.TagNumber(473458696)
  $core.bool hasRedrivecount() => $_has(10);
  @$pb.TagNumber(473458696)
  void clearRedrivecount() => $_clearField(473458696);

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(11);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(11, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(11);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class ExecutionNotRedrivable extends $pb.GeneratedMessage {
  factory ExecutionNotRedrivable({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExecutionNotRedrivable._();

  factory ExecutionNotRedrivable.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionNotRedrivable.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionNotRedrivable',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionNotRedrivable clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionNotRedrivable copyWith(
          void Function(ExecutionNotRedrivable) updates) =>
      super.copyWith((message) => updates(message as ExecutionNotRedrivable))
          as ExecutionNotRedrivable;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionNotRedrivable create() => ExecutionNotRedrivable._();
  @$core.override
  ExecutionNotRedrivable createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionNotRedrivable getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionNotRedrivable>(create);
  static ExecutionNotRedrivable? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExecutionRedrivenEventDetails extends $pb.GeneratedMessage {
  factory ExecutionRedrivenEventDetails({
    $core.int? redrivecount,
  }) {
    final result = create();
    if (redrivecount != null) result.redrivecount = redrivecount;
    return result;
  }

  ExecutionRedrivenEventDetails._();

  factory ExecutionRedrivenEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionRedrivenEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionRedrivenEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aI(473458696, _omitFieldNames ? '' : 'redrivecount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionRedrivenEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionRedrivenEventDetails copyWith(
          void Function(ExecutionRedrivenEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ExecutionRedrivenEventDetails))
          as ExecutionRedrivenEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionRedrivenEventDetails create() =>
      ExecutionRedrivenEventDetails._();
  @$core.override
  ExecutionRedrivenEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionRedrivenEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionRedrivenEventDetails>(create);
  static ExecutionRedrivenEventDetails? _defaultInstance;

  @$pb.TagNumber(473458696)
  $core.int get redrivecount => $_getIZ(0);
  @$pb.TagNumber(473458696)
  set redrivecount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(473458696)
  $core.bool hasRedrivecount() => $_has(0);
  @$pb.TagNumber(473458696)
  void clearRedrivecount() => $_clearField(473458696);
}

class ExecutionStartedEventDetails extends $pb.GeneratedMessage {
  factory ExecutionStartedEventDetails({
    $core.String? statemachineversionarn,
    $core.String? rolearn,
    $core.String? input,
    HistoryEventExecutionDataDetails? inputdetails,
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    if (rolearn != null) result.rolearn = rolearn;
    if (input != null) result.input = input;
    if (inputdetails != null) result.inputdetails = inputdetails;
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  ExecutionStartedEventDetails._();

  factory ExecutionStartedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionStartedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionStartedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..aOS(170019745, _omitFieldNames ? '' : 'rolearn')
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..aOM<HistoryEventExecutionDataDetails>(
        452625788, _omitFieldNames ? '' : 'inputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionStartedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionStartedEventDetails copyWith(
          void Function(ExecutionStartedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ExecutionStartedEventDetails))
          as ExecutionStartedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionStartedEventDetails create() =>
      ExecutionStartedEventDetails._();
  @$core.override
  ExecutionStartedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionStartedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionStartedEventDetails>(create);
  static ExecutionStartedEventDetails? _defaultInstance;

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(0);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(0);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);

  @$pb.TagNumber(170019745)
  $core.String get rolearn => $_getSZ(1);
  @$pb.TagNumber(170019745)
  set rolearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(170019745)
  $core.bool hasRolearn() => $_has(1);
  @$pb.TagNumber(170019745)
  void clearRolearn() => $_clearField(170019745);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(2);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(2, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(2);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);

  @$pb.TagNumber(452625788)
  HistoryEventExecutionDataDetails get inputdetails => $_getN(3);
  @$pb.TagNumber(452625788)
  set inputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(452625788, value);
  @$pb.TagNumber(452625788)
  $core.bool hasInputdetails() => $_has(3);
  @$pb.TagNumber(452625788)
  void clearInputdetails() => $_clearField(452625788);
  @$pb.TagNumber(452625788)
  HistoryEventExecutionDataDetails ensureInputdetails() => $_ensure(3);

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(4);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(4);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class ExecutionSucceededEventDetails extends $pb.GeneratedMessage {
  factory ExecutionSucceededEventDetails({
    HistoryEventExecutionDataDetails? outputdetails,
    $core.String? output,
  }) {
    final result = create();
    if (outputdetails != null) result.outputdetails = outputdetails;
    if (output != null) result.output = output;
    return result;
  }

  ExecutionSucceededEventDetails._();

  factory ExecutionSucceededEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionSucceededEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionSucceededEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOM<HistoryEventExecutionDataDetails>(
        393734643, _omitFieldNames ? '' : 'outputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionSucceededEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionSucceededEventDetails copyWith(
          void Function(ExecutionSucceededEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ExecutionSucceededEventDetails))
          as ExecutionSucceededEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionSucceededEventDetails create() =>
      ExecutionSucceededEventDetails._();
  @$core.override
  ExecutionSucceededEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionSucceededEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionSucceededEventDetails>(create);
  static ExecutionSucceededEventDetails? _defaultInstance;

  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails get outputdetails => $_getN(0);
  @$pb.TagNumber(393734643)
  set outputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(393734643, value);
  @$pb.TagNumber(393734643)
  $core.bool hasOutputdetails() => $_has(0);
  @$pb.TagNumber(393734643)
  void clearOutputdetails() => $_clearField(393734643);
  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails ensureOutputdetails() => $_ensure(0);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(1);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(1, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(1);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);
}

class ExecutionTimedOutEventDetails extends $pb.GeneratedMessage {
  factory ExecutionTimedOutEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  ExecutionTimedOutEventDetails._();

  factory ExecutionTimedOutEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecutionTimedOutEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecutionTimedOutEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionTimedOutEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecutionTimedOutEventDetails copyWith(
          void Function(ExecutionTimedOutEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as ExecutionTimedOutEventDetails))
          as ExecutionTimedOutEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecutionTimedOutEventDetails create() =>
      ExecutionTimedOutEventDetails._();
  @$core.override
  ExecutionTimedOutEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecutionTimedOutEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecutionTimedOutEventDetails>(create);
  static ExecutionTimedOutEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class GetActivityTaskInput extends $pb.GeneratedMessage {
  factory GetActivityTaskInput({
    $core.String? activityarn,
    $core.String? workername,
  }) {
    final result = create();
    if (activityarn != null) result.activityarn = activityarn;
    if (workername != null) result.workername = workername;
    return result;
  }

  GetActivityTaskInput._();

  factory GetActivityTaskInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetActivityTaskInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetActivityTaskInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(327279492, _omitFieldNames ? '' : 'activityarn')
    ..aOS(526761781, _omitFieldNames ? '' : 'workername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetActivityTaskInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetActivityTaskInput copyWith(void Function(GetActivityTaskInput) updates) =>
      super.copyWith((message) => updates(message as GetActivityTaskInput))
          as GetActivityTaskInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetActivityTaskInput create() => GetActivityTaskInput._();
  @$core.override
  GetActivityTaskInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetActivityTaskInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetActivityTaskInput>(create);
  static GetActivityTaskInput? _defaultInstance;

  @$pb.TagNumber(327279492)
  $core.String get activityarn => $_getSZ(0);
  @$pb.TagNumber(327279492)
  set activityarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(327279492)
  $core.bool hasActivityarn() => $_has(0);
  @$pb.TagNumber(327279492)
  void clearActivityarn() => $_clearField(327279492);

  @$pb.TagNumber(526761781)
  $core.String get workername => $_getSZ(1);
  @$pb.TagNumber(526761781)
  set workername($core.String value) => $_setString(1, value);
  @$pb.TagNumber(526761781)
  $core.bool hasWorkername() => $_has(1);
  @$pb.TagNumber(526761781)
  void clearWorkername() => $_clearField(526761781);
}

class GetActivityTaskOutput extends $pb.GeneratedMessage {
  factory GetActivityTaskOutput({
    $core.String? input,
    $core.String? tasktoken,
  }) {
    final result = create();
    if (input != null) result.input = input;
    if (tasktoken != null) result.tasktoken = tasktoken;
    return result;
  }

  GetActivityTaskOutput._();

  factory GetActivityTaskOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetActivityTaskOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetActivityTaskOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..aOS(525325834, _omitFieldNames ? '' : 'tasktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetActivityTaskOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetActivityTaskOutput copyWith(
          void Function(GetActivityTaskOutput) updates) =>
      super.copyWith((message) => updates(message as GetActivityTaskOutput))
          as GetActivityTaskOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetActivityTaskOutput create() => GetActivityTaskOutput._();
  @$core.override
  GetActivityTaskOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetActivityTaskOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetActivityTaskOutput>(create);
  static GetActivityTaskOutput? _defaultInstance;

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(0);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(0, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(0);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);

  @$pb.TagNumber(525325834)
  $core.String get tasktoken => $_getSZ(1);
  @$pb.TagNumber(525325834)
  set tasktoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(525325834)
  $core.bool hasTasktoken() => $_has(1);
  @$pb.TagNumber(525325834)
  void clearTasktoken() => $_clearField(525325834);
}

class GetExecutionHistoryInput extends $pb.GeneratedMessage {
  factory GetExecutionHistoryInput({
    $core.String? nexttoken,
    $core.bool? includeexecutiondata,
    $core.String? executionarn,
    $core.bool? reverseorder,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (includeexecutiondata != null)
      result.includeexecutiondata = includeexecutiondata;
    if (executionarn != null) result.executionarn = executionarn;
    if (reverseorder != null) result.reverseorder = reverseorder;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  GetExecutionHistoryInput._();

  factory GetExecutionHistoryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetExecutionHistoryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetExecutionHistoryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..aOB(203899608, _omitFieldNames ? '' : 'includeexecutiondata')
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..aOB(364411768, _omitFieldNames ? '' : 'reverseorder')
    ..aI(465170002, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetExecutionHistoryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetExecutionHistoryInput copyWith(
          void Function(GetExecutionHistoryInput) updates) =>
      super.copyWith((message) => updates(message as GetExecutionHistoryInput))
          as GetExecutionHistoryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetExecutionHistoryInput create() => GetExecutionHistoryInput._();
  @$core.override
  GetExecutionHistoryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetExecutionHistoryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetExecutionHistoryInput>(create);
  static GetExecutionHistoryInput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(203899608)
  $core.bool get includeexecutiondata => $_getBF(1);
  @$pb.TagNumber(203899608)
  set includeexecutiondata($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(203899608)
  $core.bool hasIncludeexecutiondata() => $_has(1);
  @$pb.TagNumber(203899608)
  void clearIncludeexecutiondata() => $_clearField(203899608);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(2);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(2);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);

  @$pb.TagNumber(364411768)
  $core.bool get reverseorder => $_getBF(3);
  @$pb.TagNumber(364411768)
  set reverseorder($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(364411768)
  $core.bool hasReverseorder() => $_has(3);
  @$pb.TagNumber(364411768)
  void clearReverseorder() => $_clearField(364411768);

  @$pb.TagNumber(465170002)
  $core.int get maxresults => $_getIZ(4);
  @$pb.TagNumber(465170002)
  set maxresults($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(465170002)
  $core.bool hasMaxresults() => $_has(4);
  @$pb.TagNumber(465170002)
  void clearMaxresults() => $_clearField(465170002);
}

class GetExecutionHistoryOutput extends $pb.GeneratedMessage {
  factory GetExecutionHistoryOutput({
    $core.String? nexttoken,
    $core.Iterable<HistoryEvent>? events,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (events != null) result.events.addAll(events);
    return result;
  }

  GetExecutionHistoryOutput._();

  factory GetExecutionHistoryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetExecutionHistoryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetExecutionHistoryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<HistoryEvent>(316203909, _omitFieldNames ? '' : 'events',
        subBuilder: HistoryEvent.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetExecutionHistoryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetExecutionHistoryOutput copyWith(
          void Function(GetExecutionHistoryOutput) updates) =>
      super.copyWith((message) => updates(message as GetExecutionHistoryOutput))
          as GetExecutionHistoryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetExecutionHistoryOutput create() => GetExecutionHistoryOutput._();
  @$core.override
  GetExecutionHistoryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetExecutionHistoryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetExecutionHistoryOutput>(create);
  static GetExecutionHistoryOutput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(316203909)
  $pb.PbList<HistoryEvent> get events => $_getList(1);
}

class HistoryEvent extends $pb.GeneratedMessage {
  factory HistoryEvent({
    TaskSucceededEventDetails? tasksucceededeventdetails,
    ActivityStartedEventDetails? activitystartedeventdetails,
    LambdaFunctionScheduledEventDetails? lambdafunctionscheduledeventdetails,
    MapRunFailedEventDetails? maprunfailedeventdetails,
    StateExitedEventDetails? stateexitedeventdetails,
    LambdaFunctionFailedEventDetails? lambdafunctionfailedeventdetails,
    ActivityScheduledEventDetails? activityscheduledeventdetails,
    MapIterationEventDetails? mapiterationsucceededeventdetails,
    ExecutionTimedOutEventDetails? executiontimedouteventdetails,
    ExecutionRedrivenEventDetails? executionredriveneventdetails,
    MapRunStartedEventDetails? maprunstartedeventdetails,
    TaskStartedEventDetails? taskstartedeventdetails,
    ActivityScheduleFailedEventDetails? activityschedulefailedeventdetails,
    TaskStartFailedEventDetails? taskstartfailedeventdetails,
    TaskScheduledEventDetails? taskscheduledeventdetails,
    $fixnum.Int64? previouseventid,
    TaskSubmittedEventDetails? tasksubmittedeventdetails,
    ActivityFailedEventDetails? activityfailedeventdetails,
    HistoryEventType? type,
    MapIterationEventDetails? mapiterationstartedeventdetails,
    ExecutionFailedEventDetails? executionfailedeventdetails,
    TaskFailedEventDetails? taskfailedeventdetails,
    $core.String? timestamp,
    TaskTimedOutEventDetails? tasktimedouteventdetails,
    LambdaFunctionSucceededEventDetails? lambdafunctionsucceededeventdetails,
    MapIterationEventDetails? mapiterationfailedeventdetails,
    ActivityTimedOutEventDetails? activitytimedouteventdetails,
    ExecutionSucceededEventDetails? executionsucceededeventdetails,
    ExecutionStartedEventDetails? executionstartedeventdetails,
    ActivitySucceededEventDetails? activitysucceededeventdetails,
    $fixnum.Int64? id,
    LambdaFunctionTimedOutEventDetails? lambdafunctiontimedouteventdetails,
    LambdaFunctionScheduleFailedEventDetails?
        lambdafunctionschedulefailedeventdetails,
    MapStateStartedEventDetails? mapstatestartedeventdetails,
    MapRunRedrivenEventDetails? maprunredriveneventdetails,
    MapIterationEventDetails? mapiterationabortedeventdetails,
    StateEnteredEventDetails? stateenteredeventdetails,
    TaskSubmitFailedEventDetails? tasksubmitfailedeventdetails,
    EvaluationFailedEventDetails? evaluationfailedeventdetails,
    ExecutionAbortedEventDetails? executionabortedeventdetails,
    LambdaFunctionStartFailedEventDetails?
        lambdafunctionstartfailedeventdetails,
  }) {
    final result = create();
    if (tasksucceededeventdetails != null)
      result.tasksucceededeventdetails = tasksucceededeventdetails;
    if (activitystartedeventdetails != null)
      result.activitystartedeventdetails = activitystartedeventdetails;
    if (lambdafunctionscheduledeventdetails != null)
      result.lambdafunctionscheduledeventdetails =
          lambdafunctionscheduledeventdetails;
    if (maprunfailedeventdetails != null)
      result.maprunfailedeventdetails = maprunfailedeventdetails;
    if (stateexitedeventdetails != null)
      result.stateexitedeventdetails = stateexitedeventdetails;
    if (lambdafunctionfailedeventdetails != null)
      result.lambdafunctionfailedeventdetails =
          lambdafunctionfailedeventdetails;
    if (activityscheduledeventdetails != null)
      result.activityscheduledeventdetails = activityscheduledeventdetails;
    if (mapiterationsucceededeventdetails != null)
      result.mapiterationsucceededeventdetails =
          mapiterationsucceededeventdetails;
    if (executiontimedouteventdetails != null)
      result.executiontimedouteventdetails = executiontimedouteventdetails;
    if (executionredriveneventdetails != null)
      result.executionredriveneventdetails = executionredriveneventdetails;
    if (maprunstartedeventdetails != null)
      result.maprunstartedeventdetails = maprunstartedeventdetails;
    if (taskstartedeventdetails != null)
      result.taskstartedeventdetails = taskstartedeventdetails;
    if (activityschedulefailedeventdetails != null)
      result.activityschedulefailedeventdetails =
          activityschedulefailedeventdetails;
    if (taskstartfailedeventdetails != null)
      result.taskstartfailedeventdetails = taskstartfailedeventdetails;
    if (taskscheduledeventdetails != null)
      result.taskscheduledeventdetails = taskscheduledeventdetails;
    if (previouseventid != null) result.previouseventid = previouseventid;
    if (tasksubmittedeventdetails != null)
      result.tasksubmittedeventdetails = tasksubmittedeventdetails;
    if (activityfailedeventdetails != null)
      result.activityfailedeventdetails = activityfailedeventdetails;
    if (type != null) result.type = type;
    if (mapiterationstartedeventdetails != null)
      result.mapiterationstartedeventdetails = mapiterationstartedeventdetails;
    if (executionfailedeventdetails != null)
      result.executionfailedeventdetails = executionfailedeventdetails;
    if (taskfailedeventdetails != null)
      result.taskfailedeventdetails = taskfailedeventdetails;
    if (timestamp != null) result.timestamp = timestamp;
    if (tasktimedouteventdetails != null)
      result.tasktimedouteventdetails = tasktimedouteventdetails;
    if (lambdafunctionsucceededeventdetails != null)
      result.lambdafunctionsucceededeventdetails =
          lambdafunctionsucceededeventdetails;
    if (mapiterationfailedeventdetails != null)
      result.mapiterationfailedeventdetails = mapiterationfailedeventdetails;
    if (activitytimedouteventdetails != null)
      result.activitytimedouteventdetails = activitytimedouteventdetails;
    if (executionsucceededeventdetails != null)
      result.executionsucceededeventdetails = executionsucceededeventdetails;
    if (executionstartedeventdetails != null)
      result.executionstartedeventdetails = executionstartedeventdetails;
    if (activitysucceededeventdetails != null)
      result.activitysucceededeventdetails = activitysucceededeventdetails;
    if (id != null) result.id = id;
    if (lambdafunctiontimedouteventdetails != null)
      result.lambdafunctiontimedouteventdetails =
          lambdafunctiontimedouteventdetails;
    if (lambdafunctionschedulefailedeventdetails != null)
      result.lambdafunctionschedulefailedeventdetails =
          lambdafunctionschedulefailedeventdetails;
    if (mapstatestartedeventdetails != null)
      result.mapstatestartedeventdetails = mapstatestartedeventdetails;
    if (maprunredriveneventdetails != null)
      result.maprunredriveneventdetails = maprunredriveneventdetails;
    if (mapiterationabortedeventdetails != null)
      result.mapiterationabortedeventdetails = mapiterationabortedeventdetails;
    if (stateenteredeventdetails != null)
      result.stateenteredeventdetails = stateenteredeventdetails;
    if (tasksubmitfailedeventdetails != null)
      result.tasksubmitfailedeventdetails = tasksubmitfailedeventdetails;
    if (evaluationfailedeventdetails != null)
      result.evaluationfailedeventdetails = evaluationfailedeventdetails;
    if (executionabortedeventdetails != null)
      result.executionabortedeventdetails = executionabortedeventdetails;
    if (lambdafunctionstartfailedeventdetails != null)
      result.lambdafunctionstartfailedeventdetails =
          lambdafunctionstartfailedeventdetails;
    return result;
  }

  HistoryEvent._();

  factory HistoryEvent.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HistoryEvent.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HistoryEvent',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOM<TaskSucceededEventDetails>(
        20918556, _omitFieldNames ? '' : 'tasksucceededeventdetails',
        subBuilder: TaskSucceededEventDetails.create)
    ..aOM<ActivityStartedEventDetails>(
        25723004, _omitFieldNames ? '' : 'activitystartedeventdetails',
        subBuilder: ActivityStartedEventDetails.create)
    ..aOM<LambdaFunctionScheduledEventDetails>(
        27704864, _omitFieldNames ? '' : 'lambdafunctionscheduledeventdetails',
        subBuilder: LambdaFunctionScheduledEventDetails.create)
    ..aOM<MapRunFailedEventDetails>(
        30147502, _omitFieldNames ? '' : 'maprunfailedeventdetails',
        subBuilder: MapRunFailedEventDetails.create)
    ..aOM<StateExitedEventDetails>(
        30960782, _omitFieldNames ? '' : 'stateexitedeventdetails',
        subBuilder: StateExitedEventDetails.create)
    ..aOM<LambdaFunctionFailedEventDetails>(
        31352398, _omitFieldNames ? '' : 'lambdafunctionfailedeventdetails',
        subBuilder: LambdaFunctionFailedEventDetails.create)
    ..aOM<ActivityScheduledEventDetails>(
        68173374, _omitFieldNames ? '' : 'activityscheduledeventdetails',
        subBuilder: ActivityScheduledEventDetails.create)
    ..aOM<MapIterationEventDetails>(
        85930040, _omitFieldNames ? '' : 'mapiterationsucceededeventdetails',
        subBuilder: MapIterationEventDetails.create)
    ..aOM<ExecutionTimedOutEventDetails>(
        139990699, _omitFieldNames ? '' : 'executiontimedouteventdetails',
        subBuilder: ExecutionTimedOutEventDetails.create)
    ..aOM<ExecutionRedrivenEventDetails>(
        157425989, _omitFieldNames ? '' : 'executionredriveneventdetails',
        subBuilder: ExecutionRedrivenEventDetails.create)
    ..aOM<MapRunStartedEventDetails>(
        183119902, _omitFieldNames ? '' : 'maprunstartedeventdetails',
        subBuilder: MapRunStartedEventDetails.create)
    ..aOM<TaskStartedEventDetails>(
        192171260, _omitFieldNames ? '' : 'taskstartedeventdetails',
        subBuilder: TaskStartedEventDetails.create)
    ..aOM<ActivityScheduleFailedEventDetails>(
        215835903, _omitFieldNames ? '' : 'activityschedulefailedeventdetails',
        subBuilder: ActivityScheduleFailedEventDetails.create)
    ..aOM<TaskStartFailedEventDetails>(
        237771688, _omitFieldNames ? '' : 'taskstartfailedeventdetails',
        subBuilder: TaskStartFailedEventDetails.create)
    ..aOM<TaskScheduledEventDetails>(
        238821054, _omitFieldNames ? '' : 'taskscheduledeventdetails',
        subBuilder: TaskScheduledEventDetails.create)
    ..aInt64(261811790, _omitFieldNames ? '' : 'previouseventid')
    ..aOM<TaskSubmittedEventDetails>(
        269619140, _omitFieldNames ? '' : 'tasksubmittedeventdetails',
        subBuilder: TaskSubmittedEventDetails.create)
    ..aOM<ActivityFailedEventDetails>(
        274703772, _omitFieldNames ? '' : 'activityfailedeventdetails',
        subBuilder: ActivityFailedEventDetails.create)
    ..aE<HistoryEventType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: HistoryEventType.values)
    ..aOM<MapIterationEventDetails>(
        294963144, _omitFieldNames ? '' : 'mapiterationstartedeventdetails',
        subBuilder: MapIterationEventDetails.create)
    ..aOM<ExecutionFailedEventDetails>(
        297327899, _omitFieldNames ? '' : 'executionfailedeventdetails',
        subBuilder: ExecutionFailedEventDetails.create)
    ..aOM<TaskFailedEventDetails>(
        309753628, _omitFieldNames ? '' : 'taskfailedeventdetails',
        subBuilder: TaskFailedEventDetails.create)
    ..aOS(310629668, _omitFieldNames ? '' : 'timestamp')
    ..aOM<TaskTimedOutEventDetails>(
        330698464, _omitFieldNames ? '' : 'tasktimedouteventdetails',
        subBuilder: TaskTimedOutEventDetails.create)
    ..aOM<LambdaFunctionSucceededEventDetails>(
        347724058, _omitFieldNames ? '' : 'lambdafunctionsucceededeventdetails',
        subBuilder: LambdaFunctionSucceededEventDetails.create)
    ..aOM<MapIterationEventDetails>(
        356461616, _omitFieldNames ? '' : 'mapiterationfailedeventdetails',
        subBuilder: MapIterationEventDetails.create)
    ..aOM<ActivityTimedOutEventDetails>(
        360915296, _omitFieldNames ? '' : 'activitytimedouteventdetails',
        subBuilder: ActivityTimedOutEventDetails.create)
    ..aOM<ExecutionSucceededEventDetails>(
        375748869, _omitFieldNames ? '' : 'executionsucceededeventdetails',
        subBuilder: ExecutionSucceededEventDetails.create)
    ..aOM<ExecutionStartedEventDetails>(
        382892969, _omitFieldNames ? '' : 'executionstartedeventdetails',
        subBuilder: ExecutionStartedEventDetails.create)
    ..aOM<ActivitySucceededEventDetails>(
        387141788, _omitFieldNames ? '' : 'activitysucceededeventdetails',
        subBuilder: ActivitySucceededEventDetails.create)
    ..aInt64(389573345, _omitFieldNames ? '' : 'id')
    ..aOM<LambdaFunctionTimedOutEventDetails>(
        423888938, _omitFieldNames ? '' : 'lambdafunctiontimedouteventdetails',
        subBuilder: LambdaFunctionTimedOutEventDetails.create)
    ..aOM<LambdaFunctionScheduleFailedEventDetails>(449961745,
        _omitFieldNames ? '' : 'lambdafunctionschedulefailedeventdetails',
        subBuilder: LambdaFunctionScheduleFailedEventDetails.create)
    ..aOM<MapStateStartedEventDetails>(
        463022434, _omitFieldNames ? '' : 'mapstatestartedeventdetails',
        subBuilder: MapStateStartedEventDetails.create)
    ..aOM<MapRunRedrivenEventDetails>(
        465066264, _omitFieldNames ? '' : 'maprunredriveneventdetails',
        subBuilder: MapRunRedrivenEventDetails.create)
    ..aOM<MapIterationEventDetails>(
        465966882, _omitFieldNames ? '' : 'mapiterationabortedeventdetails',
        subBuilder: MapIterationEventDetails.create)
    ..aOM<StateEnteredEventDetails>(
        479306524, _omitFieldNames ? '' : 'stateenteredeventdetails',
        subBuilder: StateEnteredEventDetails.create)
    ..aOM<TaskSubmitFailedEventDetails>(
        489315884, _omitFieldNames ? '' : 'tasksubmitfailedeventdetails',
        subBuilder: TaskSubmitFailedEventDetails.create)
    ..aOM<EvaluationFailedEventDetails>(
        490113981, _omitFieldNames ? '' : 'evaluationfailedeventdetails',
        subBuilder: EvaluationFailedEventDetails.create)
    ..aOM<ExecutionAbortedEventDetails>(
        507567407, _omitFieldNames ? '' : 'executionabortedeventdetails',
        subBuilder: ExecutionAbortedEventDetails.create)
    ..aOM<LambdaFunctionStartFailedEventDetails>(508021246,
        _omitFieldNames ? '' : 'lambdafunctionstartfailedeventdetails',
        subBuilder: LambdaFunctionStartFailedEventDetails.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HistoryEvent clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HistoryEvent copyWith(void Function(HistoryEvent) updates) =>
      super.copyWith((message) => updates(message as HistoryEvent))
          as HistoryEvent;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HistoryEvent create() => HistoryEvent._();
  @$core.override
  HistoryEvent createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HistoryEvent getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HistoryEvent>(create);
  static HistoryEvent? _defaultInstance;

  @$pb.TagNumber(20918556)
  TaskSucceededEventDetails get tasksucceededeventdetails => $_getN(0);
  @$pb.TagNumber(20918556)
  set tasksucceededeventdetails(TaskSucceededEventDetails value) =>
      $_setField(20918556, value);
  @$pb.TagNumber(20918556)
  $core.bool hasTasksucceededeventdetails() => $_has(0);
  @$pb.TagNumber(20918556)
  void clearTasksucceededeventdetails() => $_clearField(20918556);
  @$pb.TagNumber(20918556)
  TaskSucceededEventDetails ensureTasksucceededeventdetails() => $_ensure(0);

  @$pb.TagNumber(25723004)
  ActivityStartedEventDetails get activitystartedeventdetails => $_getN(1);
  @$pb.TagNumber(25723004)
  set activitystartedeventdetails(ActivityStartedEventDetails value) =>
      $_setField(25723004, value);
  @$pb.TagNumber(25723004)
  $core.bool hasActivitystartedeventdetails() => $_has(1);
  @$pb.TagNumber(25723004)
  void clearActivitystartedeventdetails() => $_clearField(25723004);
  @$pb.TagNumber(25723004)
  ActivityStartedEventDetails ensureActivitystartedeventdetails() =>
      $_ensure(1);

  @$pb.TagNumber(27704864)
  LambdaFunctionScheduledEventDetails get lambdafunctionscheduledeventdetails =>
      $_getN(2);
  @$pb.TagNumber(27704864)
  set lambdafunctionscheduledeventdetails(
          LambdaFunctionScheduledEventDetails value) =>
      $_setField(27704864, value);
  @$pb.TagNumber(27704864)
  $core.bool hasLambdafunctionscheduledeventdetails() => $_has(2);
  @$pb.TagNumber(27704864)
  void clearLambdafunctionscheduledeventdetails() => $_clearField(27704864);
  @$pb.TagNumber(27704864)
  LambdaFunctionScheduledEventDetails
      ensureLambdafunctionscheduledeventdetails() => $_ensure(2);

  @$pb.TagNumber(30147502)
  MapRunFailedEventDetails get maprunfailedeventdetails => $_getN(3);
  @$pb.TagNumber(30147502)
  set maprunfailedeventdetails(MapRunFailedEventDetails value) =>
      $_setField(30147502, value);
  @$pb.TagNumber(30147502)
  $core.bool hasMaprunfailedeventdetails() => $_has(3);
  @$pb.TagNumber(30147502)
  void clearMaprunfailedeventdetails() => $_clearField(30147502);
  @$pb.TagNumber(30147502)
  MapRunFailedEventDetails ensureMaprunfailedeventdetails() => $_ensure(3);

  @$pb.TagNumber(30960782)
  StateExitedEventDetails get stateexitedeventdetails => $_getN(4);
  @$pb.TagNumber(30960782)
  set stateexitedeventdetails(StateExitedEventDetails value) =>
      $_setField(30960782, value);
  @$pb.TagNumber(30960782)
  $core.bool hasStateexitedeventdetails() => $_has(4);
  @$pb.TagNumber(30960782)
  void clearStateexitedeventdetails() => $_clearField(30960782);
  @$pb.TagNumber(30960782)
  StateExitedEventDetails ensureStateexitedeventdetails() => $_ensure(4);

  @$pb.TagNumber(31352398)
  LambdaFunctionFailedEventDetails get lambdafunctionfailedeventdetails =>
      $_getN(5);
  @$pb.TagNumber(31352398)
  set lambdafunctionfailedeventdetails(
          LambdaFunctionFailedEventDetails value) =>
      $_setField(31352398, value);
  @$pb.TagNumber(31352398)
  $core.bool hasLambdafunctionfailedeventdetails() => $_has(5);
  @$pb.TagNumber(31352398)
  void clearLambdafunctionfailedeventdetails() => $_clearField(31352398);
  @$pb.TagNumber(31352398)
  LambdaFunctionFailedEventDetails ensureLambdafunctionfailedeventdetails() =>
      $_ensure(5);

  @$pb.TagNumber(68173374)
  ActivityScheduledEventDetails get activityscheduledeventdetails => $_getN(6);
  @$pb.TagNumber(68173374)
  set activityscheduledeventdetails(ActivityScheduledEventDetails value) =>
      $_setField(68173374, value);
  @$pb.TagNumber(68173374)
  $core.bool hasActivityscheduledeventdetails() => $_has(6);
  @$pb.TagNumber(68173374)
  void clearActivityscheduledeventdetails() => $_clearField(68173374);
  @$pb.TagNumber(68173374)
  ActivityScheduledEventDetails ensureActivityscheduledeventdetails() =>
      $_ensure(6);

  @$pb.TagNumber(85930040)
  MapIterationEventDetails get mapiterationsucceededeventdetails => $_getN(7);
  @$pb.TagNumber(85930040)
  set mapiterationsucceededeventdetails(MapIterationEventDetails value) =>
      $_setField(85930040, value);
  @$pb.TagNumber(85930040)
  $core.bool hasMapiterationsucceededeventdetails() => $_has(7);
  @$pb.TagNumber(85930040)
  void clearMapiterationsucceededeventdetails() => $_clearField(85930040);
  @$pb.TagNumber(85930040)
  MapIterationEventDetails ensureMapiterationsucceededeventdetails() =>
      $_ensure(7);

  @$pb.TagNumber(139990699)
  ExecutionTimedOutEventDetails get executiontimedouteventdetails => $_getN(8);
  @$pb.TagNumber(139990699)
  set executiontimedouteventdetails(ExecutionTimedOutEventDetails value) =>
      $_setField(139990699, value);
  @$pb.TagNumber(139990699)
  $core.bool hasExecutiontimedouteventdetails() => $_has(8);
  @$pb.TagNumber(139990699)
  void clearExecutiontimedouteventdetails() => $_clearField(139990699);
  @$pb.TagNumber(139990699)
  ExecutionTimedOutEventDetails ensureExecutiontimedouteventdetails() =>
      $_ensure(8);

  @$pb.TagNumber(157425989)
  ExecutionRedrivenEventDetails get executionredriveneventdetails => $_getN(9);
  @$pb.TagNumber(157425989)
  set executionredriveneventdetails(ExecutionRedrivenEventDetails value) =>
      $_setField(157425989, value);
  @$pb.TagNumber(157425989)
  $core.bool hasExecutionredriveneventdetails() => $_has(9);
  @$pb.TagNumber(157425989)
  void clearExecutionredriveneventdetails() => $_clearField(157425989);
  @$pb.TagNumber(157425989)
  ExecutionRedrivenEventDetails ensureExecutionredriveneventdetails() =>
      $_ensure(9);

  @$pb.TagNumber(183119902)
  MapRunStartedEventDetails get maprunstartedeventdetails => $_getN(10);
  @$pb.TagNumber(183119902)
  set maprunstartedeventdetails(MapRunStartedEventDetails value) =>
      $_setField(183119902, value);
  @$pb.TagNumber(183119902)
  $core.bool hasMaprunstartedeventdetails() => $_has(10);
  @$pb.TagNumber(183119902)
  void clearMaprunstartedeventdetails() => $_clearField(183119902);
  @$pb.TagNumber(183119902)
  MapRunStartedEventDetails ensureMaprunstartedeventdetails() => $_ensure(10);

  @$pb.TagNumber(192171260)
  TaskStartedEventDetails get taskstartedeventdetails => $_getN(11);
  @$pb.TagNumber(192171260)
  set taskstartedeventdetails(TaskStartedEventDetails value) =>
      $_setField(192171260, value);
  @$pb.TagNumber(192171260)
  $core.bool hasTaskstartedeventdetails() => $_has(11);
  @$pb.TagNumber(192171260)
  void clearTaskstartedeventdetails() => $_clearField(192171260);
  @$pb.TagNumber(192171260)
  TaskStartedEventDetails ensureTaskstartedeventdetails() => $_ensure(11);

  @$pb.TagNumber(215835903)
  ActivityScheduleFailedEventDetails get activityschedulefailedeventdetails =>
      $_getN(12);
  @$pb.TagNumber(215835903)
  set activityschedulefailedeventdetails(
          ActivityScheduleFailedEventDetails value) =>
      $_setField(215835903, value);
  @$pb.TagNumber(215835903)
  $core.bool hasActivityschedulefailedeventdetails() => $_has(12);
  @$pb.TagNumber(215835903)
  void clearActivityschedulefailedeventdetails() => $_clearField(215835903);
  @$pb.TagNumber(215835903)
  ActivityScheduleFailedEventDetails
      ensureActivityschedulefailedeventdetails() => $_ensure(12);

  @$pb.TagNumber(237771688)
  TaskStartFailedEventDetails get taskstartfailedeventdetails => $_getN(13);
  @$pb.TagNumber(237771688)
  set taskstartfailedeventdetails(TaskStartFailedEventDetails value) =>
      $_setField(237771688, value);
  @$pb.TagNumber(237771688)
  $core.bool hasTaskstartfailedeventdetails() => $_has(13);
  @$pb.TagNumber(237771688)
  void clearTaskstartfailedeventdetails() => $_clearField(237771688);
  @$pb.TagNumber(237771688)
  TaskStartFailedEventDetails ensureTaskstartfailedeventdetails() =>
      $_ensure(13);

  @$pb.TagNumber(238821054)
  TaskScheduledEventDetails get taskscheduledeventdetails => $_getN(14);
  @$pb.TagNumber(238821054)
  set taskscheduledeventdetails(TaskScheduledEventDetails value) =>
      $_setField(238821054, value);
  @$pb.TagNumber(238821054)
  $core.bool hasTaskscheduledeventdetails() => $_has(14);
  @$pb.TagNumber(238821054)
  void clearTaskscheduledeventdetails() => $_clearField(238821054);
  @$pb.TagNumber(238821054)
  TaskScheduledEventDetails ensureTaskscheduledeventdetails() => $_ensure(14);

  @$pb.TagNumber(261811790)
  $fixnum.Int64 get previouseventid => $_getI64(15);
  @$pb.TagNumber(261811790)
  set previouseventid($fixnum.Int64 value) => $_setInt64(15, value);
  @$pb.TagNumber(261811790)
  $core.bool hasPreviouseventid() => $_has(15);
  @$pb.TagNumber(261811790)
  void clearPreviouseventid() => $_clearField(261811790);

  @$pb.TagNumber(269619140)
  TaskSubmittedEventDetails get tasksubmittedeventdetails => $_getN(16);
  @$pb.TagNumber(269619140)
  set tasksubmittedeventdetails(TaskSubmittedEventDetails value) =>
      $_setField(269619140, value);
  @$pb.TagNumber(269619140)
  $core.bool hasTasksubmittedeventdetails() => $_has(16);
  @$pb.TagNumber(269619140)
  void clearTasksubmittedeventdetails() => $_clearField(269619140);
  @$pb.TagNumber(269619140)
  TaskSubmittedEventDetails ensureTasksubmittedeventdetails() => $_ensure(16);

  @$pb.TagNumber(274703772)
  ActivityFailedEventDetails get activityfailedeventdetails => $_getN(17);
  @$pb.TagNumber(274703772)
  set activityfailedeventdetails(ActivityFailedEventDetails value) =>
      $_setField(274703772, value);
  @$pb.TagNumber(274703772)
  $core.bool hasActivityfailedeventdetails() => $_has(17);
  @$pb.TagNumber(274703772)
  void clearActivityfailedeventdetails() => $_clearField(274703772);
  @$pb.TagNumber(274703772)
  ActivityFailedEventDetails ensureActivityfailedeventdetails() => $_ensure(17);

  @$pb.TagNumber(287830350)
  HistoryEventType get type => $_getN(18);
  @$pb.TagNumber(287830350)
  set type(HistoryEventType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(18);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(294963144)
  MapIterationEventDetails get mapiterationstartedeventdetails => $_getN(19);
  @$pb.TagNumber(294963144)
  set mapiterationstartedeventdetails(MapIterationEventDetails value) =>
      $_setField(294963144, value);
  @$pb.TagNumber(294963144)
  $core.bool hasMapiterationstartedeventdetails() => $_has(19);
  @$pb.TagNumber(294963144)
  void clearMapiterationstartedeventdetails() => $_clearField(294963144);
  @$pb.TagNumber(294963144)
  MapIterationEventDetails ensureMapiterationstartedeventdetails() =>
      $_ensure(19);

  @$pb.TagNumber(297327899)
  ExecutionFailedEventDetails get executionfailedeventdetails => $_getN(20);
  @$pb.TagNumber(297327899)
  set executionfailedeventdetails(ExecutionFailedEventDetails value) =>
      $_setField(297327899, value);
  @$pb.TagNumber(297327899)
  $core.bool hasExecutionfailedeventdetails() => $_has(20);
  @$pb.TagNumber(297327899)
  void clearExecutionfailedeventdetails() => $_clearField(297327899);
  @$pb.TagNumber(297327899)
  ExecutionFailedEventDetails ensureExecutionfailedeventdetails() =>
      $_ensure(20);

  @$pb.TagNumber(309753628)
  TaskFailedEventDetails get taskfailedeventdetails => $_getN(21);
  @$pb.TagNumber(309753628)
  set taskfailedeventdetails(TaskFailedEventDetails value) =>
      $_setField(309753628, value);
  @$pb.TagNumber(309753628)
  $core.bool hasTaskfailedeventdetails() => $_has(21);
  @$pb.TagNumber(309753628)
  void clearTaskfailedeventdetails() => $_clearField(309753628);
  @$pb.TagNumber(309753628)
  TaskFailedEventDetails ensureTaskfailedeventdetails() => $_ensure(21);

  @$pb.TagNumber(310629668)
  $core.String get timestamp => $_getSZ(22);
  @$pb.TagNumber(310629668)
  set timestamp($core.String value) => $_setString(22, value);
  @$pb.TagNumber(310629668)
  $core.bool hasTimestamp() => $_has(22);
  @$pb.TagNumber(310629668)
  void clearTimestamp() => $_clearField(310629668);

  @$pb.TagNumber(330698464)
  TaskTimedOutEventDetails get tasktimedouteventdetails => $_getN(23);
  @$pb.TagNumber(330698464)
  set tasktimedouteventdetails(TaskTimedOutEventDetails value) =>
      $_setField(330698464, value);
  @$pb.TagNumber(330698464)
  $core.bool hasTasktimedouteventdetails() => $_has(23);
  @$pb.TagNumber(330698464)
  void clearTasktimedouteventdetails() => $_clearField(330698464);
  @$pb.TagNumber(330698464)
  TaskTimedOutEventDetails ensureTasktimedouteventdetails() => $_ensure(23);

  @$pb.TagNumber(347724058)
  LambdaFunctionSucceededEventDetails get lambdafunctionsucceededeventdetails =>
      $_getN(24);
  @$pb.TagNumber(347724058)
  set lambdafunctionsucceededeventdetails(
          LambdaFunctionSucceededEventDetails value) =>
      $_setField(347724058, value);
  @$pb.TagNumber(347724058)
  $core.bool hasLambdafunctionsucceededeventdetails() => $_has(24);
  @$pb.TagNumber(347724058)
  void clearLambdafunctionsucceededeventdetails() => $_clearField(347724058);
  @$pb.TagNumber(347724058)
  LambdaFunctionSucceededEventDetails
      ensureLambdafunctionsucceededeventdetails() => $_ensure(24);

  @$pb.TagNumber(356461616)
  MapIterationEventDetails get mapiterationfailedeventdetails => $_getN(25);
  @$pb.TagNumber(356461616)
  set mapiterationfailedeventdetails(MapIterationEventDetails value) =>
      $_setField(356461616, value);
  @$pb.TagNumber(356461616)
  $core.bool hasMapiterationfailedeventdetails() => $_has(25);
  @$pb.TagNumber(356461616)
  void clearMapiterationfailedeventdetails() => $_clearField(356461616);
  @$pb.TagNumber(356461616)
  MapIterationEventDetails ensureMapiterationfailedeventdetails() =>
      $_ensure(25);

  @$pb.TagNumber(360915296)
  ActivityTimedOutEventDetails get activitytimedouteventdetails => $_getN(26);
  @$pb.TagNumber(360915296)
  set activitytimedouteventdetails(ActivityTimedOutEventDetails value) =>
      $_setField(360915296, value);
  @$pb.TagNumber(360915296)
  $core.bool hasActivitytimedouteventdetails() => $_has(26);
  @$pb.TagNumber(360915296)
  void clearActivitytimedouteventdetails() => $_clearField(360915296);
  @$pb.TagNumber(360915296)
  ActivityTimedOutEventDetails ensureActivitytimedouteventdetails() =>
      $_ensure(26);

  @$pb.TagNumber(375748869)
  ExecutionSucceededEventDetails get executionsucceededeventdetails =>
      $_getN(27);
  @$pb.TagNumber(375748869)
  set executionsucceededeventdetails(ExecutionSucceededEventDetails value) =>
      $_setField(375748869, value);
  @$pb.TagNumber(375748869)
  $core.bool hasExecutionsucceededeventdetails() => $_has(27);
  @$pb.TagNumber(375748869)
  void clearExecutionsucceededeventdetails() => $_clearField(375748869);
  @$pb.TagNumber(375748869)
  ExecutionSucceededEventDetails ensureExecutionsucceededeventdetails() =>
      $_ensure(27);

  @$pb.TagNumber(382892969)
  ExecutionStartedEventDetails get executionstartedeventdetails => $_getN(28);
  @$pb.TagNumber(382892969)
  set executionstartedeventdetails(ExecutionStartedEventDetails value) =>
      $_setField(382892969, value);
  @$pb.TagNumber(382892969)
  $core.bool hasExecutionstartedeventdetails() => $_has(28);
  @$pb.TagNumber(382892969)
  void clearExecutionstartedeventdetails() => $_clearField(382892969);
  @$pb.TagNumber(382892969)
  ExecutionStartedEventDetails ensureExecutionstartedeventdetails() =>
      $_ensure(28);

  @$pb.TagNumber(387141788)
  ActivitySucceededEventDetails get activitysucceededeventdetails => $_getN(29);
  @$pb.TagNumber(387141788)
  set activitysucceededeventdetails(ActivitySucceededEventDetails value) =>
      $_setField(387141788, value);
  @$pb.TagNumber(387141788)
  $core.bool hasActivitysucceededeventdetails() => $_has(29);
  @$pb.TagNumber(387141788)
  void clearActivitysucceededeventdetails() => $_clearField(387141788);
  @$pb.TagNumber(387141788)
  ActivitySucceededEventDetails ensureActivitysucceededeventdetails() =>
      $_ensure(29);

  @$pb.TagNumber(389573345)
  $fixnum.Int64 get id => $_getI64(30);
  @$pb.TagNumber(389573345)
  set id($fixnum.Int64 value) => $_setInt64(30, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(30);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);

  @$pb.TagNumber(423888938)
  LambdaFunctionTimedOutEventDetails get lambdafunctiontimedouteventdetails =>
      $_getN(31);
  @$pb.TagNumber(423888938)
  set lambdafunctiontimedouteventdetails(
          LambdaFunctionTimedOutEventDetails value) =>
      $_setField(423888938, value);
  @$pb.TagNumber(423888938)
  $core.bool hasLambdafunctiontimedouteventdetails() => $_has(31);
  @$pb.TagNumber(423888938)
  void clearLambdafunctiontimedouteventdetails() => $_clearField(423888938);
  @$pb.TagNumber(423888938)
  LambdaFunctionTimedOutEventDetails
      ensureLambdafunctiontimedouteventdetails() => $_ensure(31);

  @$pb.TagNumber(449961745)
  LambdaFunctionScheduleFailedEventDetails
      get lambdafunctionschedulefailedeventdetails => $_getN(32);
  @$pb.TagNumber(449961745)
  set lambdafunctionschedulefailedeventdetails(
          LambdaFunctionScheduleFailedEventDetails value) =>
      $_setField(449961745, value);
  @$pb.TagNumber(449961745)
  $core.bool hasLambdafunctionschedulefailedeventdetails() => $_has(32);
  @$pb.TagNumber(449961745)
  void clearLambdafunctionschedulefailedeventdetails() =>
      $_clearField(449961745);
  @$pb.TagNumber(449961745)
  LambdaFunctionScheduleFailedEventDetails
      ensureLambdafunctionschedulefailedeventdetails() => $_ensure(32);

  @$pb.TagNumber(463022434)
  MapStateStartedEventDetails get mapstatestartedeventdetails => $_getN(33);
  @$pb.TagNumber(463022434)
  set mapstatestartedeventdetails(MapStateStartedEventDetails value) =>
      $_setField(463022434, value);
  @$pb.TagNumber(463022434)
  $core.bool hasMapstatestartedeventdetails() => $_has(33);
  @$pb.TagNumber(463022434)
  void clearMapstatestartedeventdetails() => $_clearField(463022434);
  @$pb.TagNumber(463022434)
  MapStateStartedEventDetails ensureMapstatestartedeventdetails() =>
      $_ensure(33);

  @$pb.TagNumber(465066264)
  MapRunRedrivenEventDetails get maprunredriveneventdetails => $_getN(34);
  @$pb.TagNumber(465066264)
  set maprunredriveneventdetails(MapRunRedrivenEventDetails value) =>
      $_setField(465066264, value);
  @$pb.TagNumber(465066264)
  $core.bool hasMaprunredriveneventdetails() => $_has(34);
  @$pb.TagNumber(465066264)
  void clearMaprunredriveneventdetails() => $_clearField(465066264);
  @$pb.TagNumber(465066264)
  MapRunRedrivenEventDetails ensureMaprunredriveneventdetails() => $_ensure(34);

  @$pb.TagNumber(465966882)
  MapIterationEventDetails get mapiterationabortedeventdetails => $_getN(35);
  @$pb.TagNumber(465966882)
  set mapiterationabortedeventdetails(MapIterationEventDetails value) =>
      $_setField(465966882, value);
  @$pb.TagNumber(465966882)
  $core.bool hasMapiterationabortedeventdetails() => $_has(35);
  @$pb.TagNumber(465966882)
  void clearMapiterationabortedeventdetails() => $_clearField(465966882);
  @$pb.TagNumber(465966882)
  MapIterationEventDetails ensureMapiterationabortedeventdetails() =>
      $_ensure(35);

  @$pb.TagNumber(479306524)
  StateEnteredEventDetails get stateenteredeventdetails => $_getN(36);
  @$pb.TagNumber(479306524)
  set stateenteredeventdetails(StateEnteredEventDetails value) =>
      $_setField(479306524, value);
  @$pb.TagNumber(479306524)
  $core.bool hasStateenteredeventdetails() => $_has(36);
  @$pb.TagNumber(479306524)
  void clearStateenteredeventdetails() => $_clearField(479306524);
  @$pb.TagNumber(479306524)
  StateEnteredEventDetails ensureStateenteredeventdetails() => $_ensure(36);

  @$pb.TagNumber(489315884)
  TaskSubmitFailedEventDetails get tasksubmitfailedeventdetails => $_getN(37);
  @$pb.TagNumber(489315884)
  set tasksubmitfailedeventdetails(TaskSubmitFailedEventDetails value) =>
      $_setField(489315884, value);
  @$pb.TagNumber(489315884)
  $core.bool hasTasksubmitfailedeventdetails() => $_has(37);
  @$pb.TagNumber(489315884)
  void clearTasksubmitfailedeventdetails() => $_clearField(489315884);
  @$pb.TagNumber(489315884)
  TaskSubmitFailedEventDetails ensureTasksubmitfailedeventdetails() =>
      $_ensure(37);

  @$pb.TagNumber(490113981)
  EvaluationFailedEventDetails get evaluationfailedeventdetails => $_getN(38);
  @$pb.TagNumber(490113981)
  set evaluationfailedeventdetails(EvaluationFailedEventDetails value) =>
      $_setField(490113981, value);
  @$pb.TagNumber(490113981)
  $core.bool hasEvaluationfailedeventdetails() => $_has(38);
  @$pb.TagNumber(490113981)
  void clearEvaluationfailedeventdetails() => $_clearField(490113981);
  @$pb.TagNumber(490113981)
  EvaluationFailedEventDetails ensureEvaluationfailedeventdetails() =>
      $_ensure(38);

  @$pb.TagNumber(507567407)
  ExecutionAbortedEventDetails get executionabortedeventdetails => $_getN(39);
  @$pb.TagNumber(507567407)
  set executionabortedeventdetails(ExecutionAbortedEventDetails value) =>
      $_setField(507567407, value);
  @$pb.TagNumber(507567407)
  $core.bool hasExecutionabortedeventdetails() => $_has(39);
  @$pb.TagNumber(507567407)
  void clearExecutionabortedeventdetails() => $_clearField(507567407);
  @$pb.TagNumber(507567407)
  ExecutionAbortedEventDetails ensureExecutionabortedeventdetails() =>
      $_ensure(39);

  @$pb.TagNumber(508021246)
  LambdaFunctionStartFailedEventDetails
      get lambdafunctionstartfailedeventdetails => $_getN(40);
  @$pb.TagNumber(508021246)
  set lambdafunctionstartfailedeventdetails(
          LambdaFunctionStartFailedEventDetails value) =>
      $_setField(508021246, value);
  @$pb.TagNumber(508021246)
  $core.bool hasLambdafunctionstartfailedeventdetails() => $_has(40);
  @$pb.TagNumber(508021246)
  void clearLambdafunctionstartfailedeventdetails() => $_clearField(508021246);
  @$pb.TagNumber(508021246)
  LambdaFunctionStartFailedEventDetails
      ensureLambdafunctionstartfailedeventdetails() => $_ensure(40);
}

class HistoryEventExecutionDataDetails extends $pb.GeneratedMessage {
  factory HistoryEventExecutionDataDetails({
    $core.bool? truncated,
  }) {
    final result = create();
    if (truncated != null) result.truncated = truncated;
    return result;
  }

  HistoryEventExecutionDataDetails._();

  factory HistoryEventExecutionDataDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HistoryEventExecutionDataDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HistoryEventExecutionDataDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOB(407490858, _omitFieldNames ? '' : 'truncated')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HistoryEventExecutionDataDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HistoryEventExecutionDataDetails copyWith(
          void Function(HistoryEventExecutionDataDetails) updates) =>
      super.copyWith(
              (message) => updates(message as HistoryEventExecutionDataDetails))
          as HistoryEventExecutionDataDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HistoryEventExecutionDataDetails create() =>
      HistoryEventExecutionDataDetails._();
  @$core.override
  HistoryEventExecutionDataDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HistoryEventExecutionDataDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HistoryEventExecutionDataDetails>(
          create);
  static HistoryEventExecutionDataDetails? _defaultInstance;

  @$pb.TagNumber(407490858)
  $core.bool get truncated => $_getBF(0);
  @$pb.TagNumber(407490858)
  set truncated($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(407490858)
  $core.bool hasTruncated() => $_has(0);
  @$pb.TagNumber(407490858)
  void clearTruncated() => $_clearField(407490858);
}

class InspectionData extends $pb.GeneratedMessage {
  factory InspectionData({
    $core.int? toleratedfailurecount,
    $core.int? maxconcurrency,
    $core.double? toleratedfailurepercentage,
    $core.String? variables,
    $core.String? result,
    $core.String? afteritembatcher,
    InspectionErrorDetails? errordetails,
    $core.String? afteritemspointer,
    $core.String? afterparameters,
    $core.String? afterinputpath,
    $core.String? afterarguments,
    $core.String? afteritemselector,
    InspectionDataResponse? response,
    $core.String? input,
    $core.String? afterresultselector,
    $core.String? afterresultpath,
    $core.String? afteritemspath,
    InspectionDataRequest? request,
  }) {
    final result$ = create();
    if (toleratedfailurecount != null)
      result$.toleratedfailurecount = toleratedfailurecount;
    if (maxconcurrency != null) result$.maxconcurrency = maxconcurrency;
    if (toleratedfailurepercentage != null)
      result$.toleratedfailurepercentage = toleratedfailurepercentage;
    if (variables != null) result$.variables = variables;
    if (result != null) result$.result = result;
    if (afteritembatcher != null) result$.afteritembatcher = afteritembatcher;
    if (errordetails != null) result$.errordetails = errordetails;
    if (afteritemspointer != null)
      result$.afteritemspointer = afteritemspointer;
    if (afterparameters != null) result$.afterparameters = afterparameters;
    if (afterinputpath != null) result$.afterinputpath = afterinputpath;
    if (afterarguments != null) result$.afterarguments = afterarguments;
    if (afteritemselector != null)
      result$.afteritemselector = afteritemselector;
    if (response != null) result$.response = response;
    if (input != null) result$.input = input;
    if (afterresultselector != null)
      result$.afterresultselector = afterresultselector;
    if (afterresultpath != null) result$.afterresultpath = afterresultpath;
    if (afteritemspath != null) result$.afteritemspath = afteritemspath;
    if (request != null) result$.request = request;
    return result$;
  }

  InspectionData._();

  factory InspectionData.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InspectionData.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InspectionData',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aI(41834811, _omitFieldNames ? '' : 'toleratedfailurecount')
    ..aI(100901405, _omitFieldNames ? '' : 'maxconcurrency')
    ..aD(116496164, _omitFieldNames ? '' : 'toleratedfailurepercentage',
        fieldType: $pb.PbFieldType.OF)
    ..aOS(162226883, _omitFieldNames ? '' : 'variables')
    ..aOS(171406885, _omitFieldNames ? '' : 'result')
    ..aOS(181123354, _omitFieldNames ? '' : 'afteritembatcher')
    ..aOM<InspectionErrorDetails>(
        192899050, _omitFieldNames ? '' : 'errordetails',
        subBuilder: InspectionErrorDetails.create)
    ..aOS(194882729, _omitFieldNames ? '' : 'afteritemspointer')
    ..aOS(328015918, _omitFieldNames ? '' : 'afterparameters')
    ..aOS(355404745, _omitFieldNames ? '' : 'afterinputpath')
    ..aOS(365960236, _omitFieldNames ? '' : 'afterarguments')
    ..aOS(394258120, _omitFieldNames ? '' : 'afteritemselector')
    ..aOM<InspectionDataResponse>(425574879, _omitFieldNames ? '' : 'response',
        subBuilder: InspectionDataResponse.create)
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..aOS(443240414, _omitFieldNames ? '' : 'afterresultselector')
    ..aOS(480099136, _omitFieldNames ? '' : 'afterresultpath')
    ..aOS(491034053, _omitFieldNames ? '' : 'afteritemspath')
    ..aOM<InspectionDataRequest>(514460083, _omitFieldNames ? '' : 'request',
        subBuilder: InspectionDataRequest.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InspectionData clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InspectionData copyWith(void Function(InspectionData) updates) =>
      super.copyWith((message) => updates(message as InspectionData))
          as InspectionData;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InspectionData create() => InspectionData._();
  @$core.override
  InspectionData createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InspectionData getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InspectionData>(create);
  static InspectionData? _defaultInstance;

  @$pb.TagNumber(41834811)
  $core.int get toleratedfailurecount => $_getIZ(0);
  @$pb.TagNumber(41834811)
  set toleratedfailurecount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(41834811)
  $core.bool hasToleratedfailurecount() => $_has(0);
  @$pb.TagNumber(41834811)
  void clearToleratedfailurecount() => $_clearField(41834811);

  @$pb.TagNumber(100901405)
  $core.int get maxconcurrency => $_getIZ(1);
  @$pb.TagNumber(100901405)
  set maxconcurrency($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(100901405)
  $core.bool hasMaxconcurrency() => $_has(1);
  @$pb.TagNumber(100901405)
  void clearMaxconcurrency() => $_clearField(100901405);

  @$pb.TagNumber(116496164)
  $core.double get toleratedfailurepercentage => $_getN(2);
  @$pb.TagNumber(116496164)
  set toleratedfailurepercentage($core.double value) => $_setFloat(2, value);
  @$pb.TagNumber(116496164)
  $core.bool hasToleratedfailurepercentage() => $_has(2);
  @$pb.TagNumber(116496164)
  void clearToleratedfailurepercentage() => $_clearField(116496164);

  @$pb.TagNumber(162226883)
  $core.String get variables => $_getSZ(3);
  @$pb.TagNumber(162226883)
  set variables($core.String value) => $_setString(3, value);
  @$pb.TagNumber(162226883)
  $core.bool hasVariables() => $_has(3);
  @$pb.TagNumber(162226883)
  void clearVariables() => $_clearField(162226883);

  @$pb.TagNumber(171406885)
  $core.String get result => $_getSZ(4);
  @$pb.TagNumber(171406885)
  set result($core.String value) => $_setString(4, value);
  @$pb.TagNumber(171406885)
  $core.bool hasResult() => $_has(4);
  @$pb.TagNumber(171406885)
  void clearResult() => $_clearField(171406885);

  @$pb.TagNumber(181123354)
  $core.String get afteritembatcher => $_getSZ(5);
  @$pb.TagNumber(181123354)
  set afteritembatcher($core.String value) => $_setString(5, value);
  @$pb.TagNumber(181123354)
  $core.bool hasAfteritembatcher() => $_has(5);
  @$pb.TagNumber(181123354)
  void clearAfteritembatcher() => $_clearField(181123354);

  @$pb.TagNumber(192899050)
  InspectionErrorDetails get errordetails => $_getN(6);
  @$pb.TagNumber(192899050)
  set errordetails(InspectionErrorDetails value) =>
      $_setField(192899050, value);
  @$pb.TagNumber(192899050)
  $core.bool hasErrordetails() => $_has(6);
  @$pb.TagNumber(192899050)
  void clearErrordetails() => $_clearField(192899050);
  @$pb.TagNumber(192899050)
  InspectionErrorDetails ensureErrordetails() => $_ensure(6);

  @$pb.TagNumber(194882729)
  $core.String get afteritemspointer => $_getSZ(7);
  @$pb.TagNumber(194882729)
  set afteritemspointer($core.String value) => $_setString(7, value);
  @$pb.TagNumber(194882729)
  $core.bool hasAfteritemspointer() => $_has(7);
  @$pb.TagNumber(194882729)
  void clearAfteritemspointer() => $_clearField(194882729);

  @$pb.TagNumber(328015918)
  $core.String get afterparameters => $_getSZ(8);
  @$pb.TagNumber(328015918)
  set afterparameters($core.String value) => $_setString(8, value);
  @$pb.TagNumber(328015918)
  $core.bool hasAfterparameters() => $_has(8);
  @$pb.TagNumber(328015918)
  void clearAfterparameters() => $_clearField(328015918);

  @$pb.TagNumber(355404745)
  $core.String get afterinputpath => $_getSZ(9);
  @$pb.TagNumber(355404745)
  set afterinputpath($core.String value) => $_setString(9, value);
  @$pb.TagNumber(355404745)
  $core.bool hasAfterinputpath() => $_has(9);
  @$pb.TagNumber(355404745)
  void clearAfterinputpath() => $_clearField(355404745);

  @$pb.TagNumber(365960236)
  $core.String get afterarguments => $_getSZ(10);
  @$pb.TagNumber(365960236)
  set afterarguments($core.String value) => $_setString(10, value);
  @$pb.TagNumber(365960236)
  $core.bool hasAfterarguments() => $_has(10);
  @$pb.TagNumber(365960236)
  void clearAfterarguments() => $_clearField(365960236);

  @$pb.TagNumber(394258120)
  $core.String get afteritemselector => $_getSZ(11);
  @$pb.TagNumber(394258120)
  set afteritemselector($core.String value) => $_setString(11, value);
  @$pb.TagNumber(394258120)
  $core.bool hasAfteritemselector() => $_has(11);
  @$pb.TagNumber(394258120)
  void clearAfteritemselector() => $_clearField(394258120);

  @$pb.TagNumber(425574879)
  InspectionDataResponse get response => $_getN(12);
  @$pb.TagNumber(425574879)
  set response(InspectionDataResponse value) => $_setField(425574879, value);
  @$pb.TagNumber(425574879)
  $core.bool hasResponse() => $_has(12);
  @$pb.TagNumber(425574879)
  void clearResponse() => $_clearField(425574879);
  @$pb.TagNumber(425574879)
  InspectionDataResponse ensureResponse() => $_ensure(12);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(13);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(13, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(13);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);

  @$pb.TagNumber(443240414)
  $core.String get afterresultselector => $_getSZ(14);
  @$pb.TagNumber(443240414)
  set afterresultselector($core.String value) => $_setString(14, value);
  @$pb.TagNumber(443240414)
  $core.bool hasAfterresultselector() => $_has(14);
  @$pb.TagNumber(443240414)
  void clearAfterresultselector() => $_clearField(443240414);

  @$pb.TagNumber(480099136)
  $core.String get afterresultpath => $_getSZ(15);
  @$pb.TagNumber(480099136)
  set afterresultpath($core.String value) => $_setString(15, value);
  @$pb.TagNumber(480099136)
  $core.bool hasAfterresultpath() => $_has(15);
  @$pb.TagNumber(480099136)
  void clearAfterresultpath() => $_clearField(480099136);

  @$pb.TagNumber(491034053)
  $core.String get afteritemspath => $_getSZ(16);
  @$pb.TagNumber(491034053)
  set afteritemspath($core.String value) => $_setString(16, value);
  @$pb.TagNumber(491034053)
  $core.bool hasAfteritemspath() => $_has(16);
  @$pb.TagNumber(491034053)
  void clearAfteritemspath() => $_clearField(491034053);

  @$pb.TagNumber(514460083)
  InspectionDataRequest get request => $_getN(17);
  @$pb.TagNumber(514460083)
  set request(InspectionDataRequest value) => $_setField(514460083, value);
  @$pb.TagNumber(514460083)
  $core.bool hasRequest() => $_has(17);
  @$pb.TagNumber(514460083)
  void clearRequest() => $_clearField(514460083);
  @$pb.TagNumber(514460083)
  InspectionDataRequest ensureRequest() => $_ensure(17);
}

class InspectionDataRequest extends $pb.GeneratedMessage {
  factory InspectionDataRequest({
    $core.String? method,
    $core.String? url,
    $core.String? headers,
    $core.String? protocol,
    $core.String? body,
  }) {
    final result = create();
    if (method != null) result.method = method;
    if (url != null) result.url = url;
    if (headers != null) result.headers = headers;
    if (protocol != null) result.protocol = protocol;
    if (body != null) result.body = body;
    return result;
  }

  InspectionDataRequest._();

  factory InspectionDataRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InspectionDataRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InspectionDataRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(189134641, _omitFieldNames ? '' : 'method')
    ..aOS(311381023, _omitFieldNames ? '' : 'url')
    ..aOS(375773674, _omitFieldNames ? '' : 'headers')
    ..aOS(455607734, _omitFieldNames ? '' : 'protocol')
    ..aOS(464157046, _omitFieldNames ? '' : 'body')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InspectionDataRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InspectionDataRequest copyWith(
          void Function(InspectionDataRequest) updates) =>
      super.copyWith((message) => updates(message as InspectionDataRequest))
          as InspectionDataRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InspectionDataRequest create() => InspectionDataRequest._();
  @$core.override
  InspectionDataRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InspectionDataRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InspectionDataRequest>(create);
  static InspectionDataRequest? _defaultInstance;

  @$pb.TagNumber(189134641)
  $core.String get method => $_getSZ(0);
  @$pb.TagNumber(189134641)
  set method($core.String value) => $_setString(0, value);
  @$pb.TagNumber(189134641)
  $core.bool hasMethod() => $_has(0);
  @$pb.TagNumber(189134641)
  void clearMethod() => $_clearField(189134641);

  @$pb.TagNumber(311381023)
  $core.String get url => $_getSZ(1);
  @$pb.TagNumber(311381023)
  set url($core.String value) => $_setString(1, value);
  @$pb.TagNumber(311381023)
  $core.bool hasUrl() => $_has(1);
  @$pb.TagNumber(311381023)
  void clearUrl() => $_clearField(311381023);

  @$pb.TagNumber(375773674)
  $core.String get headers => $_getSZ(2);
  @$pb.TagNumber(375773674)
  set headers($core.String value) => $_setString(2, value);
  @$pb.TagNumber(375773674)
  $core.bool hasHeaders() => $_has(2);
  @$pb.TagNumber(375773674)
  void clearHeaders() => $_clearField(375773674);

  @$pb.TagNumber(455607734)
  $core.String get protocol => $_getSZ(3);
  @$pb.TagNumber(455607734)
  set protocol($core.String value) => $_setString(3, value);
  @$pb.TagNumber(455607734)
  $core.bool hasProtocol() => $_has(3);
  @$pb.TagNumber(455607734)
  void clearProtocol() => $_clearField(455607734);

  @$pb.TagNumber(464157046)
  $core.String get body => $_getSZ(4);
  @$pb.TagNumber(464157046)
  set body($core.String value) => $_setString(4, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(4);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);
}

class InspectionDataResponse extends $pb.GeneratedMessage {
  factory InspectionDataResponse({
    $core.String? statuscode,
    $core.String? headers,
    $core.String? protocol,
    $core.String? body,
    $core.String? statusmessage,
  }) {
    final result = create();
    if (statuscode != null) result.statuscode = statuscode;
    if (headers != null) result.headers = headers;
    if (protocol != null) result.protocol = protocol;
    if (body != null) result.body = body;
    if (statusmessage != null) result.statusmessage = statusmessage;
    return result;
  }

  InspectionDataResponse._();

  factory InspectionDataResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InspectionDataResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InspectionDataResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(375773674, _omitFieldNames ? '' : 'headers')
    ..aOS(455607734, _omitFieldNames ? '' : 'protocol')
    ..aOS(464157046, _omitFieldNames ? '' : 'body')
    ..aOS(474462255, _omitFieldNames ? '' : 'statusmessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InspectionDataResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InspectionDataResponse copyWith(
          void Function(InspectionDataResponse) updates) =>
      super.copyWith((message) => updates(message as InspectionDataResponse))
          as InspectionDataResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InspectionDataResponse create() => InspectionDataResponse._();
  @$core.override
  InspectionDataResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InspectionDataResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InspectionDataResponse>(create);
  static InspectionDataResponse? _defaultInstance;

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(0);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(0);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(375773674)
  $core.String get headers => $_getSZ(1);
  @$pb.TagNumber(375773674)
  set headers($core.String value) => $_setString(1, value);
  @$pb.TagNumber(375773674)
  $core.bool hasHeaders() => $_has(1);
  @$pb.TagNumber(375773674)
  void clearHeaders() => $_clearField(375773674);

  @$pb.TagNumber(455607734)
  $core.String get protocol => $_getSZ(2);
  @$pb.TagNumber(455607734)
  set protocol($core.String value) => $_setString(2, value);
  @$pb.TagNumber(455607734)
  $core.bool hasProtocol() => $_has(2);
  @$pb.TagNumber(455607734)
  void clearProtocol() => $_clearField(455607734);

  @$pb.TagNumber(464157046)
  $core.String get body => $_getSZ(3);
  @$pb.TagNumber(464157046)
  set body($core.String value) => $_setString(3, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(3);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);

  @$pb.TagNumber(474462255)
  $core.String get statusmessage => $_getSZ(4);
  @$pb.TagNumber(474462255)
  set statusmessage($core.String value) => $_setString(4, value);
  @$pb.TagNumber(474462255)
  $core.bool hasStatusmessage() => $_has(4);
  @$pb.TagNumber(474462255)
  void clearStatusmessage() => $_clearField(474462255);
}

class InspectionErrorDetails extends $pb.GeneratedMessage {
  factory InspectionErrorDetails({
    $core.int? retryindex,
    $core.int? retrybackoffintervalseconds,
    $core.int? catchindex,
  }) {
    final result = create();
    if (retryindex != null) result.retryindex = retryindex;
    if (retrybackoffintervalseconds != null)
      result.retrybackoffintervalseconds = retrybackoffintervalseconds;
    if (catchindex != null) result.catchindex = catchindex;
    return result;
  }

  InspectionErrorDetails._();

  factory InspectionErrorDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InspectionErrorDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InspectionErrorDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aI(6599888, _omitFieldNames ? '' : 'retryindex')
    ..aI(137754128, _omitFieldNames ? '' : 'retrybackoffintervalseconds')
    ..aI(290365641, _omitFieldNames ? '' : 'catchindex')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InspectionErrorDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InspectionErrorDetails copyWith(
          void Function(InspectionErrorDetails) updates) =>
      super.copyWith((message) => updates(message as InspectionErrorDetails))
          as InspectionErrorDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InspectionErrorDetails create() => InspectionErrorDetails._();
  @$core.override
  InspectionErrorDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InspectionErrorDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InspectionErrorDetails>(create);
  static InspectionErrorDetails? _defaultInstance;

  @$pb.TagNumber(6599888)
  $core.int get retryindex => $_getIZ(0);
  @$pb.TagNumber(6599888)
  set retryindex($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(6599888)
  $core.bool hasRetryindex() => $_has(0);
  @$pb.TagNumber(6599888)
  void clearRetryindex() => $_clearField(6599888);

  @$pb.TagNumber(137754128)
  $core.int get retrybackoffintervalseconds => $_getIZ(1);
  @$pb.TagNumber(137754128)
  set retrybackoffintervalseconds($core.int value) =>
      $_setSignedInt32(1, value);
  @$pb.TagNumber(137754128)
  $core.bool hasRetrybackoffintervalseconds() => $_has(1);
  @$pb.TagNumber(137754128)
  void clearRetrybackoffintervalseconds() => $_clearField(137754128);

  @$pb.TagNumber(290365641)
  $core.int get catchindex => $_getIZ(2);
  @$pb.TagNumber(290365641)
  set catchindex($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(290365641)
  $core.bool hasCatchindex() => $_has(2);
  @$pb.TagNumber(290365641)
  void clearCatchindex() => $_clearField(290365641);
}

class InvalidArn extends $pb.GeneratedMessage {
  factory InvalidArn({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidArn._();

  factory InvalidArn.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidArn.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidArn',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArn clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArn copyWith(void Function(InvalidArn) updates) =>
      super.copyWith((message) => updates(message as InvalidArn)) as InvalidArn;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidArn create() => InvalidArn._();
  @$core.override
  InvalidArn createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidArn getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidArn>(create);
  static InvalidArn? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidDefinition extends $pb.GeneratedMessage {
  factory InvalidDefinition({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidDefinition._();

  factory InvalidDefinition.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidDefinition.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidDefinition',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidDefinition clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidDefinition copyWith(void Function(InvalidDefinition) updates) =>
      super.copyWith((message) => updates(message as InvalidDefinition))
          as InvalidDefinition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidDefinition create() => InvalidDefinition._();
  @$core.override
  InvalidDefinition createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidDefinition getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidDefinition>(create);
  static InvalidDefinition? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidEncryptionConfiguration extends $pb.GeneratedMessage {
  factory InvalidEncryptionConfiguration({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidEncryptionConfiguration._();

  factory InvalidEncryptionConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidEncryptionConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidEncryptionConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEncryptionConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEncryptionConfiguration copyWith(
          void Function(InvalidEncryptionConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidEncryptionConfiguration))
          as InvalidEncryptionConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidEncryptionConfiguration create() =>
      InvalidEncryptionConfiguration._();
  @$core.override
  InvalidEncryptionConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidEncryptionConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidEncryptionConfiguration>(create);
  static InvalidEncryptionConfiguration? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidExecutionInput extends $pb.GeneratedMessage {
  factory InvalidExecutionInput({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidExecutionInput._();

  factory InvalidExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidExecutionInput copyWith(
          void Function(InvalidExecutionInput) updates) =>
      super.copyWith((message) => updates(message as InvalidExecutionInput))
          as InvalidExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidExecutionInput create() => InvalidExecutionInput._();
  @$core.override
  InvalidExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidExecutionInput>(create);
  static InvalidExecutionInput? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidLoggingConfiguration extends $pb.GeneratedMessage {
  factory InvalidLoggingConfiguration({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidLoggingConfiguration._();

  factory InvalidLoggingConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidLoggingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidLoggingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidLoggingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidLoggingConfiguration copyWith(
          void Function(InvalidLoggingConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidLoggingConfiguration))
          as InvalidLoggingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidLoggingConfiguration create() =>
      InvalidLoggingConfiguration._();
  @$core.override
  InvalidLoggingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidLoggingConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidLoggingConfiguration>(create);
  static InvalidLoggingConfiguration? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidName extends $pb.GeneratedMessage {
  factory InvalidName({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidName._();

  factory InvalidName.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidName.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidName',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidName clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidName copyWith(void Function(InvalidName) updates) =>
      super.copyWith((message) => updates(message as InvalidName))
          as InvalidName;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidName create() => InvalidName._();
  @$core.override
  InvalidName createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidName getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidName>(create);
  static InvalidName? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidOutput extends $pb.GeneratedMessage {
  factory InvalidOutput({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidOutput._();

  factory InvalidOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidOutput copyWith(void Function(InvalidOutput) updates) =>
      super.copyWith((message) => updates(message as InvalidOutput))
          as InvalidOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidOutput create() => InvalidOutput._();
  @$core.override
  InvalidOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidOutput>(create);
  static InvalidOutput? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidToken extends $pb.GeneratedMessage {
  factory InvalidToken({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidToken._();

  factory InvalidToken.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidToken.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidToken',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidToken clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidToken copyWith(void Function(InvalidToken) updates) =>
      super.copyWith((message) => updates(message as InvalidToken))
          as InvalidToken;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidToken create() => InvalidToken._();
  @$core.override
  InvalidToken createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidToken getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidToken>(create);
  static InvalidToken? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidTracingConfiguration extends $pb.GeneratedMessage {
  factory InvalidTracingConfiguration({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidTracingConfiguration._();

  factory InvalidTracingConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidTracingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidTracingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTracingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTracingConfiguration copyWith(
          void Function(InvalidTracingConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidTracingConfiguration))
          as InvalidTracingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidTracingConfiguration create() =>
      InvalidTracingConfiguration._();
  @$core.override
  InvalidTracingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidTracingConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidTracingConfiguration>(create);
  static InvalidTracingConfiguration? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsAccessDeniedException extends $pb.GeneratedMessage {
  factory KmsAccessDeniedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsAccessDeniedException._();

  factory KmsAccessDeniedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsAccessDeniedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsAccessDeniedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsAccessDeniedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsAccessDeniedException copyWith(
          void Function(KmsAccessDeniedException) updates) =>
      super.copyWith((message) => updates(message as KmsAccessDeniedException))
          as KmsAccessDeniedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsAccessDeniedException create() => KmsAccessDeniedException._();
  @$core.override
  KmsAccessDeniedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsAccessDeniedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsAccessDeniedException>(create);
  static KmsAccessDeniedException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsInvalidStateException extends $pb.GeneratedMessage {
  factory KmsInvalidStateException({
    $core.String? message,
    KmsKeyState? kmskeystate,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (kmskeystate != null) result.kmskeystate = kmskeystate;
    return result;
  }

  KmsInvalidStateException._();

  factory KmsInvalidStateException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsInvalidStateException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsInvalidStateException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..aE<KmsKeyState>(485859435, _omitFieldNames ? '' : 'kmskeystate',
        enumValues: KmsKeyState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsInvalidStateException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsInvalidStateException copyWith(
          void Function(KmsInvalidStateException) updates) =>
      super.copyWith((message) => updates(message as KmsInvalidStateException))
          as KmsInvalidStateException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsInvalidStateException create() => KmsInvalidStateException._();
  @$core.override
  KmsInvalidStateException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsInvalidStateException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsInvalidStateException>(create);
  static KmsInvalidStateException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(485859435)
  KmsKeyState get kmskeystate => $_getN(1);
  @$pb.TagNumber(485859435)
  set kmskeystate(KmsKeyState value) => $_setField(485859435, value);
  @$pb.TagNumber(485859435)
  $core.bool hasKmskeystate() => $_has(1);
  @$pb.TagNumber(485859435)
  void clearKmskeystate() => $_clearField(485859435);
}

class KmsThrottlingException extends $pb.GeneratedMessage {
  factory KmsThrottlingException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsThrottlingException._();

  factory KmsThrottlingException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsThrottlingException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsThrottlingException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsThrottlingException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsThrottlingException copyWith(
          void Function(KmsThrottlingException) updates) =>
      super.copyWith((message) => updates(message as KmsThrottlingException))
          as KmsThrottlingException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsThrottlingException create() => KmsThrottlingException._();
  @$core.override
  KmsThrottlingException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsThrottlingException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsThrottlingException>(create);
  static KmsThrottlingException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class LambdaFunctionFailedEventDetails extends $pb.GeneratedMessage {
  factory LambdaFunctionFailedEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  LambdaFunctionFailedEventDetails._();

  factory LambdaFunctionFailedEventDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LambdaFunctionFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LambdaFunctionFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionFailedEventDetails copyWith(
          void Function(LambdaFunctionFailedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as LambdaFunctionFailedEventDetails))
          as LambdaFunctionFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LambdaFunctionFailedEventDetails create() =>
      LambdaFunctionFailedEventDetails._();
  @$core.override
  LambdaFunctionFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LambdaFunctionFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LambdaFunctionFailedEventDetails>(
          create);
  static LambdaFunctionFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class LambdaFunctionScheduleFailedEventDetails extends $pb.GeneratedMessage {
  factory LambdaFunctionScheduleFailedEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  LambdaFunctionScheduleFailedEventDetails._();

  factory LambdaFunctionScheduleFailedEventDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LambdaFunctionScheduleFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LambdaFunctionScheduleFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionScheduleFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionScheduleFailedEventDetails copyWith(
          void Function(LambdaFunctionScheduleFailedEventDetails) updates) =>
      super.copyWith((message) =>
              updates(message as LambdaFunctionScheduleFailedEventDetails))
          as LambdaFunctionScheduleFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LambdaFunctionScheduleFailedEventDetails create() =>
      LambdaFunctionScheduleFailedEventDetails._();
  @$core.override
  LambdaFunctionScheduleFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LambdaFunctionScheduleFailedEventDetails getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          LambdaFunctionScheduleFailedEventDetails>(create);
  static LambdaFunctionScheduleFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class LambdaFunctionScheduledEventDetails extends $pb.GeneratedMessage {
  factory LambdaFunctionScheduledEventDetails({
    $core.String? resource,
    TaskCredentials? taskcredentials,
    $core.String? input,
    HistoryEventExecutionDataDetails? inputdetails,
    $fixnum.Int64? timeoutinseconds,
  }) {
    final result = create();
    if (resource != null) result.resource = resource;
    if (taskcredentials != null) result.taskcredentials = taskcredentials;
    if (input != null) result.input = input;
    if (inputdetails != null) result.inputdetails = inputdetails;
    if (timeoutinseconds != null) result.timeoutinseconds = timeoutinseconds;
    return result;
  }

  LambdaFunctionScheduledEventDetails._();

  factory LambdaFunctionScheduledEventDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LambdaFunctionScheduledEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LambdaFunctionScheduledEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..aOM<TaskCredentials>(257843259, _omitFieldNames ? '' : 'taskcredentials',
        subBuilder: TaskCredentials.create)
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..aOM<HistoryEventExecutionDataDetails>(
        452625788, _omitFieldNames ? '' : 'inputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..aInt64(472710197, _omitFieldNames ? '' : 'timeoutinseconds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionScheduledEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionScheduledEventDetails copyWith(
          void Function(LambdaFunctionScheduledEventDetails) updates) =>
      super.copyWith((message) =>
              updates(message as LambdaFunctionScheduledEventDetails))
          as LambdaFunctionScheduledEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LambdaFunctionScheduledEventDetails create() =>
      LambdaFunctionScheduledEventDetails._();
  @$core.override
  LambdaFunctionScheduledEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LambdaFunctionScheduledEventDetails getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          LambdaFunctionScheduledEventDetails>(create);
  static LambdaFunctionScheduledEventDetails? _defaultInstance;

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(0);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(0, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(0);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);

  @$pb.TagNumber(257843259)
  TaskCredentials get taskcredentials => $_getN(1);
  @$pb.TagNumber(257843259)
  set taskcredentials(TaskCredentials value) => $_setField(257843259, value);
  @$pb.TagNumber(257843259)
  $core.bool hasTaskcredentials() => $_has(1);
  @$pb.TagNumber(257843259)
  void clearTaskcredentials() => $_clearField(257843259);
  @$pb.TagNumber(257843259)
  TaskCredentials ensureTaskcredentials() => $_ensure(1);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(2);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(2, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(2);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);

  @$pb.TagNumber(452625788)
  HistoryEventExecutionDataDetails get inputdetails => $_getN(3);
  @$pb.TagNumber(452625788)
  set inputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(452625788, value);
  @$pb.TagNumber(452625788)
  $core.bool hasInputdetails() => $_has(3);
  @$pb.TagNumber(452625788)
  void clearInputdetails() => $_clearField(452625788);
  @$pb.TagNumber(452625788)
  HistoryEventExecutionDataDetails ensureInputdetails() => $_ensure(3);

  @$pb.TagNumber(472710197)
  $fixnum.Int64 get timeoutinseconds => $_getI64(4);
  @$pb.TagNumber(472710197)
  set timeoutinseconds($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(472710197)
  $core.bool hasTimeoutinseconds() => $_has(4);
  @$pb.TagNumber(472710197)
  void clearTimeoutinseconds() => $_clearField(472710197);
}

class LambdaFunctionStartFailedEventDetails extends $pb.GeneratedMessage {
  factory LambdaFunctionStartFailedEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  LambdaFunctionStartFailedEventDetails._();

  factory LambdaFunctionStartFailedEventDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LambdaFunctionStartFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LambdaFunctionStartFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionStartFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionStartFailedEventDetails copyWith(
          void Function(LambdaFunctionStartFailedEventDetails) updates) =>
      super.copyWith((message) =>
              updates(message as LambdaFunctionStartFailedEventDetails))
          as LambdaFunctionStartFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LambdaFunctionStartFailedEventDetails create() =>
      LambdaFunctionStartFailedEventDetails._();
  @$core.override
  LambdaFunctionStartFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LambdaFunctionStartFailedEventDetails getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          LambdaFunctionStartFailedEventDetails>(create);
  static LambdaFunctionStartFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class LambdaFunctionSucceededEventDetails extends $pb.GeneratedMessage {
  factory LambdaFunctionSucceededEventDetails({
    HistoryEventExecutionDataDetails? outputdetails,
    $core.String? output,
  }) {
    final result = create();
    if (outputdetails != null) result.outputdetails = outputdetails;
    if (output != null) result.output = output;
    return result;
  }

  LambdaFunctionSucceededEventDetails._();

  factory LambdaFunctionSucceededEventDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LambdaFunctionSucceededEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LambdaFunctionSucceededEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOM<HistoryEventExecutionDataDetails>(
        393734643, _omitFieldNames ? '' : 'outputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionSucceededEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionSucceededEventDetails copyWith(
          void Function(LambdaFunctionSucceededEventDetails) updates) =>
      super.copyWith((message) =>
              updates(message as LambdaFunctionSucceededEventDetails))
          as LambdaFunctionSucceededEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LambdaFunctionSucceededEventDetails create() =>
      LambdaFunctionSucceededEventDetails._();
  @$core.override
  LambdaFunctionSucceededEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LambdaFunctionSucceededEventDetails getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          LambdaFunctionSucceededEventDetails>(create);
  static LambdaFunctionSucceededEventDetails? _defaultInstance;

  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails get outputdetails => $_getN(0);
  @$pb.TagNumber(393734643)
  set outputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(393734643, value);
  @$pb.TagNumber(393734643)
  $core.bool hasOutputdetails() => $_has(0);
  @$pb.TagNumber(393734643)
  void clearOutputdetails() => $_clearField(393734643);
  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails ensureOutputdetails() => $_ensure(0);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(1);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(1, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(1);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);
}

class LambdaFunctionTimedOutEventDetails extends $pb.GeneratedMessage {
  factory LambdaFunctionTimedOutEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  LambdaFunctionTimedOutEventDetails._();

  factory LambdaFunctionTimedOutEventDetails.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LambdaFunctionTimedOutEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LambdaFunctionTimedOutEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionTimedOutEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LambdaFunctionTimedOutEventDetails copyWith(
          void Function(LambdaFunctionTimedOutEventDetails) updates) =>
      super.copyWith((message) =>
              updates(message as LambdaFunctionTimedOutEventDetails))
          as LambdaFunctionTimedOutEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LambdaFunctionTimedOutEventDetails create() =>
      LambdaFunctionTimedOutEventDetails._();
  @$core.override
  LambdaFunctionTimedOutEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LambdaFunctionTimedOutEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LambdaFunctionTimedOutEventDetails>(
          create);
  static LambdaFunctionTimedOutEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class ListActivitiesInput extends $pb.GeneratedMessage {
  factory ListActivitiesInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListActivitiesInput._();

  factory ListActivitiesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListActivitiesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListActivitiesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..aI(465170002, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListActivitiesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListActivitiesInput copyWith(void Function(ListActivitiesInput) updates) =>
      super.copyWith((message) => updates(message as ListActivitiesInput))
          as ListActivitiesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListActivitiesInput create() => ListActivitiesInput._();
  @$core.override
  ListActivitiesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListActivitiesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListActivitiesInput>(create);
  static ListActivitiesInput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(465170002)
  $core.int get maxresults => $_getIZ(1);
  @$pb.TagNumber(465170002)
  set maxresults($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(465170002)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(465170002)
  void clearMaxresults() => $_clearField(465170002);
}

class ListActivitiesOutput extends $pb.GeneratedMessage {
  factory ListActivitiesOutput({
    $core.String? nexttoken,
    $core.Iterable<ActivityListItem>? activities,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (activities != null) result.activities.addAll(activities);
    return result;
  }

  ListActivitiesOutput._();

  factory ListActivitiesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListActivitiesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListActivitiesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ActivityListItem>(164950681, _omitFieldNames ? '' : 'activities',
        subBuilder: ActivityListItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListActivitiesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListActivitiesOutput copyWith(void Function(ListActivitiesOutput) updates) =>
      super.copyWith((message) => updates(message as ListActivitiesOutput))
          as ListActivitiesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListActivitiesOutput create() => ListActivitiesOutput._();
  @$core.override
  ListActivitiesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListActivitiesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListActivitiesOutput>(create);
  static ListActivitiesOutput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(164950681)
  $pb.PbList<ActivityListItem> get activities => $_getList(1);
}

class ListExecutionsInput extends $pb.GeneratedMessage {
  factory ListExecutionsInput({
    $core.String? maprunarn,
    ExecutionStatus? statusfilter,
    $core.String? nexttoken,
    $core.String? statemachinearn,
    ExecutionRedriveFilter? redrivefilter,
    $core.int? maxresults,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    if (statusfilter != null) result.statusfilter = statusfilter;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (redrivefilter != null) result.redrivefilter = redrivefilter;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListExecutionsInput._();

  factory ListExecutionsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListExecutionsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListExecutionsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..aE<ExecutionStatus>(86045418, _omitFieldNames ? '' : 'statusfilter',
        enumValues: ExecutionStatus.values)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aE<ExecutionRedriveFilter>(
        437804251, _omitFieldNames ? '' : 'redrivefilter',
        enumValues: ExecutionRedriveFilter.values)
    ..aI(465170002, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExecutionsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExecutionsInput copyWith(void Function(ListExecutionsInput) updates) =>
      super.copyWith((message) => updates(message as ListExecutionsInput))
          as ListExecutionsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListExecutionsInput create() => ListExecutionsInput._();
  @$core.override
  ListExecutionsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListExecutionsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListExecutionsInput>(create);
  static ListExecutionsInput? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);

  @$pb.TagNumber(86045418)
  ExecutionStatus get statusfilter => $_getN(1);
  @$pb.TagNumber(86045418)
  set statusfilter(ExecutionStatus value) => $_setField(86045418, value);
  @$pb.TagNumber(86045418)
  $core.bool hasStatusfilter() => $_has(1);
  @$pb.TagNumber(86045418)
  void clearStatusfilter() => $_clearField(86045418);

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(3);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(3);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(437804251)
  ExecutionRedriveFilter get redrivefilter => $_getN(4);
  @$pb.TagNumber(437804251)
  set redrivefilter(ExecutionRedriveFilter value) =>
      $_setField(437804251, value);
  @$pb.TagNumber(437804251)
  $core.bool hasRedrivefilter() => $_has(4);
  @$pb.TagNumber(437804251)
  void clearRedrivefilter() => $_clearField(437804251);

  @$pb.TagNumber(465170002)
  $core.int get maxresults => $_getIZ(5);
  @$pb.TagNumber(465170002)
  set maxresults($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(465170002)
  $core.bool hasMaxresults() => $_has(5);
  @$pb.TagNumber(465170002)
  void clearMaxresults() => $_clearField(465170002);
}

class ListExecutionsOutput extends $pb.GeneratedMessage {
  factory ListExecutionsOutput({
    $core.Iterable<ExecutionListItem>? executions,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (executions != null) result.executions.addAll(executions);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListExecutionsOutput._();

  factory ListExecutionsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListExecutionsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListExecutionsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..pPM<ExecutionListItem>(113395923, _omitFieldNames ? '' : 'executions',
        subBuilder: ExecutionListItem.create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExecutionsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExecutionsOutput copyWith(void Function(ListExecutionsOutput) updates) =>
      super.copyWith((message) => updates(message as ListExecutionsOutput))
          as ListExecutionsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListExecutionsOutput create() => ListExecutionsOutput._();
  @$core.override
  ListExecutionsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListExecutionsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListExecutionsOutput>(create);
  static ListExecutionsOutput? _defaultInstance;

  @$pb.TagNumber(113395923)
  $pb.PbList<ExecutionListItem> get executions => $_getList(0);

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);
}

class ListMapRunsInput extends $pb.GeneratedMessage {
  factory ListMapRunsInput({
    $core.String? nexttoken,
    $core.String? executionarn,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (executionarn != null) result.executionarn = executionarn;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListMapRunsInput._();

  factory ListMapRunsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMapRunsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMapRunsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..aI(465170002, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMapRunsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMapRunsInput copyWith(void Function(ListMapRunsInput) updates) =>
      super.copyWith((message) => updates(message as ListMapRunsInput))
          as ListMapRunsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMapRunsInput create() => ListMapRunsInput._();
  @$core.override
  ListMapRunsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMapRunsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMapRunsInput>(create);
  static ListMapRunsInput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(1);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(1);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);

  @$pb.TagNumber(465170002)
  $core.int get maxresults => $_getIZ(2);
  @$pb.TagNumber(465170002)
  set maxresults($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(465170002)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(465170002)
  void clearMaxresults() => $_clearField(465170002);
}

class ListMapRunsOutput extends $pb.GeneratedMessage {
  factory ListMapRunsOutput({
    $core.Iterable<MapRunListItem>? mapruns,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (mapruns != null) result.mapruns.addAll(mapruns);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListMapRunsOutput._();

  factory ListMapRunsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMapRunsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMapRunsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..pPM<MapRunListItem>(34690200, _omitFieldNames ? '' : 'mapruns',
        subBuilder: MapRunListItem.create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMapRunsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMapRunsOutput copyWith(void Function(ListMapRunsOutput) updates) =>
      super.copyWith((message) => updates(message as ListMapRunsOutput))
          as ListMapRunsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMapRunsOutput create() => ListMapRunsOutput._();
  @$core.override
  ListMapRunsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMapRunsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMapRunsOutput>(create);
  static ListMapRunsOutput? _defaultInstance;

  @$pb.TagNumber(34690200)
  $pb.PbList<MapRunListItem> get mapruns => $_getList(0);

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);
}

class ListStateMachineAliasesInput extends $pb.GeneratedMessage {
  factory ListStateMachineAliasesInput({
    $core.String? nexttoken,
    $core.String? statemachinearn,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListStateMachineAliasesInput._();

  factory ListStateMachineAliasesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStateMachineAliasesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStateMachineAliasesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aI(465170002, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachineAliasesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachineAliasesInput copyWith(
          void Function(ListStateMachineAliasesInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListStateMachineAliasesInput))
          as ListStateMachineAliasesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStateMachineAliasesInput create() =>
      ListStateMachineAliasesInput._();
  @$core.override
  ListStateMachineAliasesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStateMachineAliasesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStateMachineAliasesInput>(create);
  static ListStateMachineAliasesInput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(1);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(1);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(465170002)
  $core.int get maxresults => $_getIZ(2);
  @$pb.TagNumber(465170002)
  set maxresults($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(465170002)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(465170002)
  void clearMaxresults() => $_clearField(465170002);
}

class ListStateMachineAliasesOutput extends $pb.GeneratedMessage {
  factory ListStateMachineAliasesOutput({
    $core.String? nexttoken,
    $core.Iterable<StateMachineAliasListItem>? statemachinealiases,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (statemachinealiases != null)
      result.statemachinealiases.addAll(statemachinealiases);
    return result;
  }

  ListStateMachineAliasesOutput._();

  factory ListStateMachineAliasesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStateMachineAliasesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStateMachineAliasesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<StateMachineAliasListItem>(
        452532502, _omitFieldNames ? '' : 'statemachinealiases',
        subBuilder: StateMachineAliasListItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachineAliasesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachineAliasesOutput copyWith(
          void Function(ListStateMachineAliasesOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ListStateMachineAliasesOutput))
          as ListStateMachineAliasesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStateMachineAliasesOutput create() =>
      ListStateMachineAliasesOutput._();
  @$core.override
  ListStateMachineAliasesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStateMachineAliasesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStateMachineAliasesOutput>(create);
  static ListStateMachineAliasesOutput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(452532502)
  $pb.PbList<StateMachineAliasListItem> get statemachinealiases => $_getList(1);
}

class ListStateMachineVersionsInput extends $pb.GeneratedMessage {
  factory ListStateMachineVersionsInput({
    $core.String? nexttoken,
    $core.String? statemachinearn,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListStateMachineVersionsInput._();

  factory ListStateMachineVersionsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStateMachineVersionsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStateMachineVersionsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aI(465170002, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachineVersionsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachineVersionsInput copyWith(
          void Function(ListStateMachineVersionsInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListStateMachineVersionsInput))
          as ListStateMachineVersionsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStateMachineVersionsInput create() =>
      ListStateMachineVersionsInput._();
  @$core.override
  ListStateMachineVersionsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStateMachineVersionsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStateMachineVersionsInput>(create);
  static ListStateMachineVersionsInput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(1);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(1);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(465170002)
  $core.int get maxresults => $_getIZ(2);
  @$pb.TagNumber(465170002)
  set maxresults($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(465170002)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(465170002)
  void clearMaxresults() => $_clearField(465170002);
}

class ListStateMachineVersionsOutput extends $pb.GeneratedMessage {
  factory ListStateMachineVersionsOutput({
    $core.String? nexttoken,
    $core.Iterable<StateMachineVersionListItem>? statemachineversions,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (statemachineversions != null)
      result.statemachineversions.addAll(statemachineversions);
    return result;
  }

  ListStateMachineVersionsOutput._();

  factory ListStateMachineVersionsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStateMachineVersionsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStateMachineVersionsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<StateMachineVersionListItem>(
        133258663, _omitFieldNames ? '' : 'statemachineversions',
        subBuilder: StateMachineVersionListItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachineVersionsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachineVersionsOutput copyWith(
          void Function(ListStateMachineVersionsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ListStateMachineVersionsOutput))
          as ListStateMachineVersionsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStateMachineVersionsOutput create() =>
      ListStateMachineVersionsOutput._();
  @$core.override
  ListStateMachineVersionsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStateMachineVersionsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStateMachineVersionsOutput>(create);
  static ListStateMachineVersionsOutput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(133258663)
  $pb.PbList<StateMachineVersionListItem> get statemachineversions =>
      $_getList(1);
}

class ListStateMachinesInput extends $pb.GeneratedMessage {
  factory ListStateMachinesInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListStateMachinesInput._();

  factory ListStateMachinesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStateMachinesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStateMachinesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..aI(465170002, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachinesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachinesInput copyWith(
          void Function(ListStateMachinesInput) updates) =>
      super.copyWith((message) => updates(message as ListStateMachinesInput))
          as ListStateMachinesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStateMachinesInput create() => ListStateMachinesInput._();
  @$core.override
  ListStateMachinesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStateMachinesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStateMachinesInput>(create);
  static ListStateMachinesInput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(465170002)
  $core.int get maxresults => $_getIZ(1);
  @$pb.TagNumber(465170002)
  set maxresults($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(465170002)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(465170002)
  void clearMaxresults() => $_clearField(465170002);
}

class ListStateMachinesOutput extends $pb.GeneratedMessage {
  factory ListStateMachinesOutput({
    $core.String? nexttoken,
    $core.Iterable<StateMachineListItem>? statemachines,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (statemachines != null) result.statemachines.addAll(statemachines);
    return result;
  }

  ListStateMachinesOutput._();

  factory ListStateMachinesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListStateMachinesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListStateMachinesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<StateMachineListItem>(
        432113525, _omitFieldNames ? '' : 'statemachines',
        subBuilder: StateMachineListItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachinesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListStateMachinesOutput copyWith(
          void Function(ListStateMachinesOutput) updates) =>
      super.copyWith((message) => updates(message as ListStateMachinesOutput))
          as ListStateMachinesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListStateMachinesOutput create() => ListStateMachinesOutput._();
  @$core.override
  ListStateMachinesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListStateMachinesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListStateMachinesOutput>(create);
  static ListStateMachinesOutput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(432113525)
  $pb.PbList<StateMachineListItem> get statemachines => $_getList(1);
}

class ListTagsForResourceInput extends $pb.GeneratedMessage {
  factory ListTagsForResourceInput({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  ListTagsForResourceInput._();

  factory ListTagsForResourceInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForResourceInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForResourceInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(67806797, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceInput copyWith(
          void Function(ListTagsForResourceInput) updates) =>
      super.copyWith((message) => updates(message as ListTagsForResourceInput))
          as ListTagsForResourceInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceInput create() => ListTagsForResourceInput._();
  @$core.override
  ListTagsForResourceInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForResourceInput>(create);
  static ListTagsForResourceInput? _defaultInstance;

  @$pb.TagNumber(67806797)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(67806797)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(67806797)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(67806797)
  void clearResourcearn() => $_clearField(67806797);
}

class ListTagsForResourceOutput extends $pb.GeneratedMessage {
  factory ListTagsForResourceOutput({
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  ListTagsForResourceOutput._();

  factory ListTagsForResourceOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForResourceOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForResourceOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..pPM<Tag>(337046433, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceOutput copyWith(
          void Function(ListTagsForResourceOutput) updates) =>
      super.copyWith((message) => updates(message as ListTagsForResourceOutput))
          as ListTagsForResourceOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceOutput create() => ListTagsForResourceOutput._();
  @$core.override
  ListTagsForResourceOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForResourceOutput>(create);
  static ListTagsForResourceOutput? _defaultInstance;

  @$pb.TagNumber(337046433)
  $pb.PbList<Tag> get tags => $_getList(0);
}

class LogDestination extends $pb.GeneratedMessage {
  factory LogDestination({
    CloudWatchLogsLogGroup? cloudwatchlogsloggroup,
  }) {
    final result = create();
    if (cloudwatchlogsloggroup != null)
      result.cloudwatchlogsloggroup = cloudwatchlogsloggroup;
    return result;
  }

  LogDestination._();

  factory LogDestination.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LogDestination.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LogDestination',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOM<CloudWatchLogsLogGroup>(
        517373608, _omitFieldNames ? '' : 'cloudwatchlogsloggroup',
        subBuilder: CloudWatchLogsLogGroup.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LogDestination clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LogDestination copyWith(void Function(LogDestination) updates) =>
      super.copyWith((message) => updates(message as LogDestination))
          as LogDestination;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LogDestination create() => LogDestination._();
  @$core.override
  LogDestination createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LogDestination getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LogDestination>(create);
  static LogDestination? _defaultInstance;

  @$pb.TagNumber(517373608)
  CloudWatchLogsLogGroup get cloudwatchlogsloggroup => $_getN(0);
  @$pb.TagNumber(517373608)
  set cloudwatchlogsloggroup(CloudWatchLogsLogGroup value) =>
      $_setField(517373608, value);
  @$pb.TagNumber(517373608)
  $core.bool hasCloudwatchlogsloggroup() => $_has(0);
  @$pb.TagNumber(517373608)
  void clearCloudwatchlogsloggroup() => $_clearField(517373608);
  @$pb.TagNumber(517373608)
  CloudWatchLogsLogGroup ensureCloudwatchlogsloggroup() => $_ensure(0);
}

class LoggingConfiguration extends $pb.GeneratedMessage {
  factory LoggingConfiguration({
    $core.Iterable<LogDestination>? destinations,
    $core.bool? includeexecutiondata,
    LogLevel? level,
  }) {
    final result = create();
    if (destinations != null) result.destinations.addAll(destinations);
    if (includeexecutiondata != null)
      result.includeexecutiondata = includeexecutiondata;
    if (level != null) result.level = level;
    return result;
  }

  LoggingConfiguration._();

  factory LoggingConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LoggingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LoggingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..pPM<LogDestination>(1617189, _omitFieldNames ? '' : 'destinations',
        subBuilder: LogDestination.create)
    ..aOB(203899608, _omitFieldNames ? '' : 'includeexecutiondata')
    ..aE<LogLevel>(463071198, _omitFieldNames ? '' : 'level',
        enumValues: LogLevel.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoggingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoggingConfiguration copyWith(void Function(LoggingConfiguration) updates) =>
      super.copyWith((message) => updates(message as LoggingConfiguration))
          as LoggingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LoggingConfiguration create() => LoggingConfiguration._();
  @$core.override
  LoggingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LoggingConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LoggingConfiguration>(create);
  static LoggingConfiguration? _defaultInstance;

  @$pb.TagNumber(1617189)
  $pb.PbList<LogDestination> get destinations => $_getList(0);

  @$pb.TagNumber(203899608)
  $core.bool get includeexecutiondata => $_getBF(1);
  @$pb.TagNumber(203899608)
  set includeexecutiondata($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(203899608)
  $core.bool hasIncludeexecutiondata() => $_has(1);
  @$pb.TagNumber(203899608)
  void clearIncludeexecutiondata() => $_clearField(203899608);

  @$pb.TagNumber(463071198)
  LogLevel get level => $_getN(2);
  @$pb.TagNumber(463071198)
  set level(LogLevel value) => $_setField(463071198, value);
  @$pb.TagNumber(463071198)
  $core.bool hasLevel() => $_has(2);
  @$pb.TagNumber(463071198)
  void clearLevel() => $_clearField(463071198);
}

class MapIterationEventDetails extends $pb.GeneratedMessage {
  factory MapIterationEventDetails({
    $core.int? index,
    $core.String? name,
  }) {
    final result = create();
    if (index != null) result.index = index;
    if (name != null) result.name = name;
    return result;
  }

  MapIterationEventDetails._();

  factory MapIterationEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MapIterationEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MapIterationEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aI(151693740, _omitFieldNames ? '' : 'index')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapIterationEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapIterationEventDetails copyWith(
          void Function(MapIterationEventDetails) updates) =>
      super.copyWith((message) => updates(message as MapIterationEventDetails))
          as MapIterationEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MapIterationEventDetails create() => MapIterationEventDetails._();
  @$core.override
  MapIterationEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MapIterationEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MapIterationEventDetails>(create);
  static MapIterationEventDetails? _defaultInstance;

  @$pb.TagNumber(151693740)
  $core.int get index => $_getIZ(0);
  @$pb.TagNumber(151693740)
  set index($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(151693740)
  $core.bool hasIndex() => $_has(0);
  @$pb.TagNumber(151693740)
  void clearIndex() => $_clearField(151693740);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);
}

class MapRunExecutionCounts extends $pb.GeneratedMessage {
  factory MapRunExecutionCounts({
    $fixnum.Int64? failed,
    $fixnum.Int64? total,
    $fixnum.Int64? resultswritten,
    $fixnum.Int64? failuresnotredrivable,
    $fixnum.Int64? succeeded,
    $fixnum.Int64? timedout,
    $fixnum.Int64? running,
    $fixnum.Int64? aborted,
    $fixnum.Int64? pendingredrive,
    $fixnum.Int64? pending,
  }) {
    final result = create();
    if (failed != null) result.failed = failed;
    if (total != null) result.total = total;
    if (resultswritten != null) result.resultswritten = resultswritten;
    if (failuresnotredrivable != null)
      result.failuresnotredrivable = failuresnotredrivable;
    if (succeeded != null) result.succeeded = succeeded;
    if (timedout != null) result.timedout = timedout;
    if (running != null) result.running = running;
    if (aborted != null) result.aborted = aborted;
    if (pendingredrive != null) result.pendingredrive = pendingredrive;
    if (pending != null) result.pending = pending;
    return result;
  }

  MapRunExecutionCounts._();

  factory MapRunExecutionCounts.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MapRunExecutionCounts.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MapRunExecutionCounts',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aInt64(11325365, _omitFieldNames ? '' : 'failed')
    ..aInt64(80777982, _omitFieldNames ? '' : 'total')
    ..aInt64(88286449, _omitFieldNames ? '' : 'resultswritten')
    ..aInt64(119492388, _omitFieldNames ? '' : 'failuresnotredrivable')
    ..aInt64(186019827, _omitFieldNames ? '' : 'succeeded')
    ..aInt64(335029733, _omitFieldNames ? '' : 'timedout')
    ..aInt64(343848781, _omitFieldNames ? '' : 'running')
    ..aInt64(381485777, _omitFieldNames ? '' : 'aborted')
    ..aInt64(405645200, _omitFieldNames ? '' : 'pendingredrive')
    ..aInt64(431980349, _omitFieldNames ? '' : 'pending')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunExecutionCounts clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunExecutionCounts copyWith(
          void Function(MapRunExecutionCounts) updates) =>
      super.copyWith((message) => updates(message as MapRunExecutionCounts))
          as MapRunExecutionCounts;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MapRunExecutionCounts create() => MapRunExecutionCounts._();
  @$core.override
  MapRunExecutionCounts createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MapRunExecutionCounts getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MapRunExecutionCounts>(create);
  static MapRunExecutionCounts? _defaultInstance;

  @$pb.TagNumber(11325365)
  $fixnum.Int64 get failed => $_getI64(0);
  @$pb.TagNumber(11325365)
  set failed($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(11325365)
  $core.bool hasFailed() => $_has(0);
  @$pb.TagNumber(11325365)
  void clearFailed() => $_clearField(11325365);

  @$pb.TagNumber(80777982)
  $fixnum.Int64 get total => $_getI64(1);
  @$pb.TagNumber(80777982)
  set total($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(80777982)
  $core.bool hasTotal() => $_has(1);
  @$pb.TagNumber(80777982)
  void clearTotal() => $_clearField(80777982);

  @$pb.TagNumber(88286449)
  $fixnum.Int64 get resultswritten => $_getI64(2);
  @$pb.TagNumber(88286449)
  set resultswritten($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(88286449)
  $core.bool hasResultswritten() => $_has(2);
  @$pb.TagNumber(88286449)
  void clearResultswritten() => $_clearField(88286449);

  @$pb.TagNumber(119492388)
  $fixnum.Int64 get failuresnotredrivable => $_getI64(3);
  @$pb.TagNumber(119492388)
  set failuresnotredrivable($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(119492388)
  $core.bool hasFailuresnotredrivable() => $_has(3);
  @$pb.TagNumber(119492388)
  void clearFailuresnotredrivable() => $_clearField(119492388);

  @$pb.TagNumber(186019827)
  $fixnum.Int64 get succeeded => $_getI64(4);
  @$pb.TagNumber(186019827)
  set succeeded($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(186019827)
  $core.bool hasSucceeded() => $_has(4);
  @$pb.TagNumber(186019827)
  void clearSucceeded() => $_clearField(186019827);

  @$pb.TagNumber(335029733)
  $fixnum.Int64 get timedout => $_getI64(5);
  @$pb.TagNumber(335029733)
  set timedout($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(335029733)
  $core.bool hasTimedout() => $_has(5);
  @$pb.TagNumber(335029733)
  void clearTimedout() => $_clearField(335029733);

  @$pb.TagNumber(343848781)
  $fixnum.Int64 get running => $_getI64(6);
  @$pb.TagNumber(343848781)
  set running($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(343848781)
  $core.bool hasRunning() => $_has(6);
  @$pb.TagNumber(343848781)
  void clearRunning() => $_clearField(343848781);

  @$pb.TagNumber(381485777)
  $fixnum.Int64 get aborted => $_getI64(7);
  @$pb.TagNumber(381485777)
  set aborted($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(381485777)
  $core.bool hasAborted() => $_has(7);
  @$pb.TagNumber(381485777)
  void clearAborted() => $_clearField(381485777);

  @$pb.TagNumber(405645200)
  $fixnum.Int64 get pendingredrive => $_getI64(8);
  @$pb.TagNumber(405645200)
  set pendingredrive($fixnum.Int64 value) => $_setInt64(8, value);
  @$pb.TagNumber(405645200)
  $core.bool hasPendingredrive() => $_has(8);
  @$pb.TagNumber(405645200)
  void clearPendingredrive() => $_clearField(405645200);

  @$pb.TagNumber(431980349)
  $fixnum.Int64 get pending => $_getI64(9);
  @$pb.TagNumber(431980349)
  set pending($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(431980349)
  $core.bool hasPending() => $_has(9);
  @$pb.TagNumber(431980349)
  void clearPending() => $_clearField(431980349);
}

class MapRunFailedEventDetails extends $pb.GeneratedMessage {
  factory MapRunFailedEventDetails({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  MapRunFailedEventDetails._();

  factory MapRunFailedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MapRunFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MapRunFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunFailedEventDetails copyWith(
          void Function(MapRunFailedEventDetails) updates) =>
      super.copyWith((message) => updates(message as MapRunFailedEventDetails))
          as MapRunFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MapRunFailedEventDetails create() => MapRunFailedEventDetails._();
  @$core.override
  MapRunFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MapRunFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MapRunFailedEventDetails>(create);
  static MapRunFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class MapRunItemCounts extends $pb.GeneratedMessage {
  factory MapRunItemCounts({
    $fixnum.Int64? failed,
    $fixnum.Int64? total,
    $fixnum.Int64? resultswritten,
    $fixnum.Int64? failuresnotredrivable,
    $fixnum.Int64? succeeded,
    $fixnum.Int64? timedout,
    $fixnum.Int64? running,
    $fixnum.Int64? aborted,
    $fixnum.Int64? pendingredrive,
    $fixnum.Int64? pending,
  }) {
    final result = create();
    if (failed != null) result.failed = failed;
    if (total != null) result.total = total;
    if (resultswritten != null) result.resultswritten = resultswritten;
    if (failuresnotredrivable != null)
      result.failuresnotredrivable = failuresnotredrivable;
    if (succeeded != null) result.succeeded = succeeded;
    if (timedout != null) result.timedout = timedout;
    if (running != null) result.running = running;
    if (aborted != null) result.aborted = aborted;
    if (pendingredrive != null) result.pendingredrive = pendingredrive;
    if (pending != null) result.pending = pending;
    return result;
  }

  MapRunItemCounts._();

  factory MapRunItemCounts.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MapRunItemCounts.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MapRunItemCounts',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aInt64(11325365, _omitFieldNames ? '' : 'failed')
    ..aInt64(80777982, _omitFieldNames ? '' : 'total')
    ..aInt64(88286449, _omitFieldNames ? '' : 'resultswritten')
    ..aInt64(119492388, _omitFieldNames ? '' : 'failuresnotredrivable')
    ..aInt64(186019827, _omitFieldNames ? '' : 'succeeded')
    ..aInt64(335029733, _omitFieldNames ? '' : 'timedout')
    ..aInt64(343848781, _omitFieldNames ? '' : 'running')
    ..aInt64(381485777, _omitFieldNames ? '' : 'aborted')
    ..aInt64(405645200, _omitFieldNames ? '' : 'pendingredrive')
    ..aInt64(431980349, _omitFieldNames ? '' : 'pending')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunItemCounts clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunItemCounts copyWith(void Function(MapRunItemCounts) updates) =>
      super.copyWith((message) => updates(message as MapRunItemCounts))
          as MapRunItemCounts;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MapRunItemCounts create() => MapRunItemCounts._();
  @$core.override
  MapRunItemCounts createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MapRunItemCounts getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MapRunItemCounts>(create);
  static MapRunItemCounts? _defaultInstance;

  @$pb.TagNumber(11325365)
  $fixnum.Int64 get failed => $_getI64(0);
  @$pb.TagNumber(11325365)
  set failed($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(11325365)
  $core.bool hasFailed() => $_has(0);
  @$pb.TagNumber(11325365)
  void clearFailed() => $_clearField(11325365);

  @$pb.TagNumber(80777982)
  $fixnum.Int64 get total => $_getI64(1);
  @$pb.TagNumber(80777982)
  set total($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(80777982)
  $core.bool hasTotal() => $_has(1);
  @$pb.TagNumber(80777982)
  void clearTotal() => $_clearField(80777982);

  @$pb.TagNumber(88286449)
  $fixnum.Int64 get resultswritten => $_getI64(2);
  @$pb.TagNumber(88286449)
  set resultswritten($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(88286449)
  $core.bool hasResultswritten() => $_has(2);
  @$pb.TagNumber(88286449)
  void clearResultswritten() => $_clearField(88286449);

  @$pb.TagNumber(119492388)
  $fixnum.Int64 get failuresnotredrivable => $_getI64(3);
  @$pb.TagNumber(119492388)
  set failuresnotredrivable($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(119492388)
  $core.bool hasFailuresnotredrivable() => $_has(3);
  @$pb.TagNumber(119492388)
  void clearFailuresnotredrivable() => $_clearField(119492388);

  @$pb.TagNumber(186019827)
  $fixnum.Int64 get succeeded => $_getI64(4);
  @$pb.TagNumber(186019827)
  set succeeded($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(186019827)
  $core.bool hasSucceeded() => $_has(4);
  @$pb.TagNumber(186019827)
  void clearSucceeded() => $_clearField(186019827);

  @$pb.TagNumber(335029733)
  $fixnum.Int64 get timedout => $_getI64(5);
  @$pb.TagNumber(335029733)
  set timedout($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(335029733)
  $core.bool hasTimedout() => $_has(5);
  @$pb.TagNumber(335029733)
  void clearTimedout() => $_clearField(335029733);

  @$pb.TagNumber(343848781)
  $fixnum.Int64 get running => $_getI64(6);
  @$pb.TagNumber(343848781)
  set running($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(343848781)
  $core.bool hasRunning() => $_has(6);
  @$pb.TagNumber(343848781)
  void clearRunning() => $_clearField(343848781);

  @$pb.TagNumber(381485777)
  $fixnum.Int64 get aborted => $_getI64(7);
  @$pb.TagNumber(381485777)
  set aborted($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(381485777)
  $core.bool hasAborted() => $_has(7);
  @$pb.TagNumber(381485777)
  void clearAborted() => $_clearField(381485777);

  @$pb.TagNumber(405645200)
  $fixnum.Int64 get pendingredrive => $_getI64(8);
  @$pb.TagNumber(405645200)
  set pendingredrive($fixnum.Int64 value) => $_setInt64(8, value);
  @$pb.TagNumber(405645200)
  $core.bool hasPendingredrive() => $_has(8);
  @$pb.TagNumber(405645200)
  void clearPendingredrive() => $_clearField(405645200);

  @$pb.TagNumber(431980349)
  $fixnum.Int64 get pending => $_getI64(9);
  @$pb.TagNumber(431980349)
  set pending($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(431980349)
  $core.bool hasPending() => $_has(9);
  @$pb.TagNumber(431980349)
  void clearPending() => $_clearField(431980349);
}

class MapRunListItem extends $pb.GeneratedMessage {
  factory MapRunListItem({
    $core.String? maprunarn,
    $core.String? stopdate,
    $core.String? executionarn,
    $core.String? startdate,
    $core.String? statemachinearn,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    if (stopdate != null) result.stopdate = stopdate;
    if (executionarn != null) result.executionarn = executionarn;
    if (startdate != null) result.startdate = startdate;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    return result;
  }

  MapRunListItem._();

  factory MapRunListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MapRunListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MapRunListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..aOS(180697434, _omitFieldNames ? '' : 'stopdate')
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..aOS(364840732, _omitFieldNames ? '' : 'startdate')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunListItem copyWith(void Function(MapRunListItem) updates) =>
      super.copyWith((message) => updates(message as MapRunListItem))
          as MapRunListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MapRunListItem create() => MapRunListItem._();
  @$core.override
  MapRunListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MapRunListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MapRunListItem>(create);
  static MapRunListItem? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);

  @$pb.TagNumber(180697434)
  $core.String get stopdate => $_getSZ(1);
  @$pb.TagNumber(180697434)
  set stopdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(180697434)
  $core.bool hasStopdate() => $_has(1);
  @$pb.TagNumber(180697434)
  void clearStopdate() => $_clearField(180697434);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(2);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(2);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);

  @$pb.TagNumber(364840732)
  $core.String get startdate => $_getSZ(3);
  @$pb.TagNumber(364840732)
  set startdate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(364840732)
  $core.bool hasStartdate() => $_has(3);
  @$pb.TagNumber(364840732)
  void clearStartdate() => $_clearField(364840732);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(4);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(4);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);
}

class MapRunRedrivenEventDetails extends $pb.GeneratedMessage {
  factory MapRunRedrivenEventDetails({
    $core.String? maprunarn,
    $core.int? redrivecount,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    if (redrivecount != null) result.redrivecount = redrivecount;
    return result;
  }

  MapRunRedrivenEventDetails._();

  factory MapRunRedrivenEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MapRunRedrivenEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MapRunRedrivenEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..aI(473458696, _omitFieldNames ? '' : 'redrivecount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunRedrivenEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunRedrivenEventDetails copyWith(
          void Function(MapRunRedrivenEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as MapRunRedrivenEventDetails))
          as MapRunRedrivenEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MapRunRedrivenEventDetails create() => MapRunRedrivenEventDetails._();
  @$core.override
  MapRunRedrivenEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MapRunRedrivenEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MapRunRedrivenEventDetails>(create);
  static MapRunRedrivenEventDetails? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);

  @$pb.TagNumber(473458696)
  $core.int get redrivecount => $_getIZ(1);
  @$pb.TagNumber(473458696)
  set redrivecount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(473458696)
  $core.bool hasRedrivecount() => $_has(1);
  @$pb.TagNumber(473458696)
  void clearRedrivecount() => $_clearField(473458696);
}

class MapRunStartedEventDetails extends $pb.GeneratedMessage {
  factory MapRunStartedEventDetails({
    $core.String? maprunarn,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    return result;
  }

  MapRunStartedEventDetails._();

  factory MapRunStartedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MapRunStartedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MapRunStartedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunStartedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapRunStartedEventDetails copyWith(
          void Function(MapRunStartedEventDetails) updates) =>
      super.copyWith((message) => updates(message as MapRunStartedEventDetails))
          as MapRunStartedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MapRunStartedEventDetails create() => MapRunStartedEventDetails._();
  @$core.override
  MapRunStartedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MapRunStartedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MapRunStartedEventDetails>(create);
  static MapRunStartedEventDetails? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);
}

class MapStateStartedEventDetails extends $pb.GeneratedMessage {
  factory MapStateStartedEventDetails({
    $core.int? length,
  }) {
    final result = create();
    if (length != null) result.length = length;
    return result;
  }

  MapStateStartedEventDetails._();

  factory MapStateStartedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MapStateStartedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MapStateStartedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aI(63976982, _omitFieldNames ? '' : 'length')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapStateStartedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MapStateStartedEventDetails copyWith(
          void Function(MapStateStartedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as MapStateStartedEventDetails))
          as MapStateStartedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MapStateStartedEventDetails create() =>
      MapStateStartedEventDetails._();
  @$core.override
  MapStateStartedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MapStateStartedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MapStateStartedEventDetails>(create);
  static MapStateStartedEventDetails? _defaultInstance;

  @$pb.TagNumber(63976982)
  $core.int get length => $_getIZ(0);
  @$pb.TagNumber(63976982)
  set length($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(63976982)
  $core.bool hasLength() => $_has(0);
  @$pb.TagNumber(63976982)
  void clearLength() => $_clearField(63976982);
}

class MissingRequiredParameter extends $pb.GeneratedMessage {
  factory MissingRequiredParameter({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  MissingRequiredParameter._();

  factory MissingRequiredParameter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MissingRequiredParameter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MissingRequiredParameter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MissingRequiredParameter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MissingRequiredParameter copyWith(
          void Function(MissingRequiredParameter) updates) =>
      super.copyWith((message) => updates(message as MissingRequiredParameter))
          as MissingRequiredParameter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MissingRequiredParameter create() => MissingRequiredParameter._();
  @$core.override
  MissingRequiredParameter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MissingRequiredParameter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MissingRequiredParameter>(create);
  static MissingRequiredParameter? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class MockErrorOutput extends $pb.GeneratedMessage {
  factory MockErrorOutput({
    $core.String? error,
    $core.String? cause,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    return result;
  }

  MockErrorOutput._();

  factory MockErrorOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MockErrorOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MockErrorOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MockErrorOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MockErrorOutput copyWith(void Function(MockErrorOutput) updates) =>
      super.copyWith((message) => updates(message as MockErrorOutput))
          as MockErrorOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MockErrorOutput create() => MockErrorOutput._();
  @$core.override
  MockErrorOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MockErrorOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MockErrorOutput>(create);
  static MockErrorOutput? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);
}

class MockInput extends $pb.GeneratedMessage {
  factory MockInput({
    $core.String? result,
    MockErrorOutput? erroroutput,
    MockResponseValidationMode? fieldvalidationmode,
  }) {
    final result$ = create();
    if (result != null) result$.result = result;
    if (erroroutput != null) result$.erroroutput = erroroutput;
    if (fieldvalidationmode != null)
      result$.fieldvalidationmode = fieldvalidationmode;
    return result$;
  }

  MockInput._();

  factory MockInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MockInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MockInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(171406885, _omitFieldNames ? '' : 'result')
    ..aOM<MockErrorOutput>(453597273, _omitFieldNames ? '' : 'erroroutput',
        subBuilder: MockErrorOutput.create)
    ..aE<MockResponseValidationMode>(
        519416556, _omitFieldNames ? '' : 'fieldvalidationmode',
        enumValues: MockResponseValidationMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MockInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MockInput copyWith(void Function(MockInput) updates) =>
      super.copyWith((message) => updates(message as MockInput)) as MockInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MockInput create() => MockInput._();
  @$core.override
  MockInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MockInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MockInput>(create);
  static MockInput? _defaultInstance;

  @$pb.TagNumber(171406885)
  $core.String get result => $_getSZ(0);
  @$pb.TagNumber(171406885)
  set result($core.String value) => $_setString(0, value);
  @$pb.TagNumber(171406885)
  $core.bool hasResult() => $_has(0);
  @$pb.TagNumber(171406885)
  void clearResult() => $_clearField(171406885);

  @$pb.TagNumber(453597273)
  MockErrorOutput get erroroutput => $_getN(1);
  @$pb.TagNumber(453597273)
  set erroroutput(MockErrorOutput value) => $_setField(453597273, value);
  @$pb.TagNumber(453597273)
  $core.bool hasErroroutput() => $_has(1);
  @$pb.TagNumber(453597273)
  void clearErroroutput() => $_clearField(453597273);
  @$pb.TagNumber(453597273)
  MockErrorOutput ensureErroroutput() => $_ensure(1);

  @$pb.TagNumber(519416556)
  MockResponseValidationMode get fieldvalidationmode => $_getN(2);
  @$pb.TagNumber(519416556)
  set fieldvalidationmode(MockResponseValidationMode value) =>
      $_setField(519416556, value);
  @$pb.TagNumber(519416556)
  $core.bool hasFieldvalidationmode() => $_has(2);
  @$pb.TagNumber(519416556)
  void clearFieldvalidationmode() => $_clearField(519416556);
}

class PublishStateMachineVersionInput extends $pb.GeneratedMessage {
  factory PublishStateMachineVersionInput({
    $core.String? description,
    $core.String? revisionid,
    $core.String? statemachinearn,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (revisionid != null) result.revisionid = revisionid;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    return result;
  }

  PublishStateMachineVersionInput._();

  factory PublishStateMachineVersionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishStateMachineVersionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishStateMachineVersionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(369170086, _omitFieldNames ? '' : 'revisionid')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishStateMachineVersionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishStateMachineVersionInput copyWith(
          void Function(PublishStateMachineVersionInput) updates) =>
      super.copyWith(
              (message) => updates(message as PublishStateMachineVersionInput))
          as PublishStateMachineVersionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishStateMachineVersionInput create() =>
      PublishStateMachineVersionInput._();
  @$core.override
  PublishStateMachineVersionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishStateMachineVersionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishStateMachineVersionInput>(
          create);
  static PublishStateMachineVersionInput? _defaultInstance;

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(369170086)
  $core.String get revisionid => $_getSZ(1);
  @$pb.TagNumber(369170086)
  set revisionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(369170086)
  $core.bool hasRevisionid() => $_has(1);
  @$pb.TagNumber(369170086)
  void clearRevisionid() => $_clearField(369170086);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(2);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(2);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);
}

class PublishStateMachineVersionOutput extends $pb.GeneratedMessage {
  factory PublishStateMachineVersionOutput({
    $core.String? statemachineversionarn,
    $core.String? creationdate,
  }) {
    final result = create();
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    if (creationdate != null) result.creationdate = creationdate;
    return result;
  }

  PublishStateMachineVersionOutput._();

  factory PublishStateMachineVersionOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishStateMachineVersionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishStateMachineVersionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishStateMachineVersionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishStateMachineVersionOutput copyWith(
          void Function(PublishStateMachineVersionOutput) updates) =>
      super.copyWith(
              (message) => updates(message as PublishStateMachineVersionOutput))
          as PublishStateMachineVersionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishStateMachineVersionOutput create() =>
      PublishStateMachineVersionOutput._();
  @$core.override
  PublishStateMachineVersionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishStateMachineVersionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishStateMachineVersionOutput>(
          create);
  static PublishStateMachineVersionOutput? _defaultInstance;

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(0);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(0);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(1);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(1);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);
}

class RedriveExecutionInput extends $pb.GeneratedMessage {
  factory RedriveExecutionInput({
    $core.String? clienttoken,
    $core.String? executionarn,
  }) {
    final result = create();
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (executionarn != null) result.executionarn = executionarn;
    return result;
  }

  RedriveExecutionInput._();

  factory RedriveExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RedriveExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RedriveExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(272531820, _omitFieldNames ? '' : 'clienttoken')
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedriveExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedriveExecutionInput copyWith(
          void Function(RedriveExecutionInput) updates) =>
      super.copyWith((message) => updates(message as RedriveExecutionInput))
          as RedriveExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RedriveExecutionInput create() => RedriveExecutionInput._();
  @$core.override
  RedriveExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RedriveExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RedriveExecutionInput>(create);
  static RedriveExecutionInput? _defaultInstance;

  @$pb.TagNumber(272531820)
  $core.String get clienttoken => $_getSZ(0);
  @$pb.TagNumber(272531820)
  set clienttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272531820)
  $core.bool hasClienttoken() => $_has(0);
  @$pb.TagNumber(272531820)
  void clearClienttoken() => $_clearField(272531820);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(1);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(1);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);
}

class RedriveExecutionOutput extends $pb.GeneratedMessage {
  factory RedriveExecutionOutput({
    $core.String? redrivedate,
  }) {
    final result = create();
    if (redrivedate != null) result.redrivedate = redrivedate;
    return result;
  }

  RedriveExecutionOutput._();

  factory RedriveExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RedriveExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RedriveExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(152812125, _omitFieldNames ? '' : 'redrivedate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedriveExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedriveExecutionOutput copyWith(
          void Function(RedriveExecutionOutput) updates) =>
      super.copyWith((message) => updates(message as RedriveExecutionOutput))
          as RedriveExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RedriveExecutionOutput create() => RedriveExecutionOutput._();
  @$core.override
  RedriveExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RedriveExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RedriveExecutionOutput>(create);
  static RedriveExecutionOutput? _defaultInstance;

  @$pb.TagNumber(152812125)
  $core.String get redrivedate => $_getSZ(0);
  @$pb.TagNumber(152812125)
  set redrivedate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(152812125)
  $core.bool hasRedrivedate() => $_has(0);
  @$pb.TagNumber(152812125)
  void clearRedrivedate() => $_clearField(152812125);
}

class ResourceNotFound extends $pb.GeneratedMessage {
  factory ResourceNotFound({
    $core.String? resourcename,
    $core.String? message,
  }) {
    final result = create();
    if (resourcename != null) result.resourcename = resourcename;
    if (message != null) result.message = message;
    return result;
  }

  ResourceNotFound._();

  factory ResourceNotFound.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceNotFound.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceNotFound',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(17776375, _omitFieldNames ? '' : 'resourcename')
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceNotFound clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceNotFound copyWith(void Function(ResourceNotFound) updates) =>
      super.copyWith((message) => updates(message as ResourceNotFound))
          as ResourceNotFound;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceNotFound create() => ResourceNotFound._();
  @$core.override
  ResourceNotFound createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceNotFound getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceNotFound>(create);
  static ResourceNotFound? _defaultInstance;

  @$pb.TagNumber(17776375)
  $core.String get resourcename => $_getSZ(0);
  @$pb.TagNumber(17776375)
  set resourcename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(17776375)
  $core.bool hasResourcename() => $_has(0);
  @$pb.TagNumber(17776375)
  void clearResourcename() => $_clearField(17776375);

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class RoutingConfigurationListItem extends $pb.GeneratedMessage {
  factory RoutingConfigurationListItem({
    $core.String? statemachineversionarn,
    $core.int? weight,
  }) {
    final result = create();
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    if (weight != null) result.weight = weight;
    return result;
  }

  RoutingConfigurationListItem._();

  factory RoutingConfigurationListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RoutingConfigurationListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RoutingConfigurationListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..aI(278961850, _omitFieldNames ? '' : 'weight')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RoutingConfigurationListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RoutingConfigurationListItem copyWith(
          void Function(RoutingConfigurationListItem) updates) =>
      super.copyWith(
              (message) => updates(message as RoutingConfigurationListItem))
          as RoutingConfigurationListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RoutingConfigurationListItem create() =>
      RoutingConfigurationListItem._();
  @$core.override
  RoutingConfigurationListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RoutingConfigurationListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RoutingConfigurationListItem>(create);
  static RoutingConfigurationListItem? _defaultInstance;

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(0);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(0);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);

  @$pb.TagNumber(278961850)
  $core.int get weight => $_getIZ(1);
  @$pb.TagNumber(278961850)
  set weight($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(278961850)
  $core.bool hasWeight() => $_has(1);
  @$pb.TagNumber(278961850)
  void clearWeight() => $_clearField(278961850);
}

class SendTaskFailureInput extends $pb.GeneratedMessage {
  factory SendTaskFailureInput({
    $core.String? error,
    $core.String? cause,
    $core.String? tasktoken,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    if (tasktoken != null) result.tasktoken = tasktoken;
    return result;
  }

  SendTaskFailureInput._();

  factory SendTaskFailureInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendTaskFailureInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendTaskFailureInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(525325834, _omitFieldNames ? '' : 'tasktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskFailureInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskFailureInput copyWith(void Function(SendTaskFailureInput) updates) =>
      super.copyWith((message) => updates(message as SendTaskFailureInput))
          as SendTaskFailureInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendTaskFailureInput create() => SendTaskFailureInput._();
  @$core.override
  SendTaskFailureInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendTaskFailureInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendTaskFailureInput>(create);
  static SendTaskFailureInput? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(525325834)
  $core.String get tasktoken => $_getSZ(2);
  @$pb.TagNumber(525325834)
  set tasktoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(525325834)
  $core.bool hasTasktoken() => $_has(2);
  @$pb.TagNumber(525325834)
  void clearTasktoken() => $_clearField(525325834);
}

class SendTaskFailureOutput extends $pb.GeneratedMessage {
  factory SendTaskFailureOutput() => create();

  SendTaskFailureOutput._();

  factory SendTaskFailureOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendTaskFailureOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendTaskFailureOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskFailureOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskFailureOutput copyWith(
          void Function(SendTaskFailureOutput) updates) =>
      super.copyWith((message) => updates(message as SendTaskFailureOutput))
          as SendTaskFailureOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendTaskFailureOutput create() => SendTaskFailureOutput._();
  @$core.override
  SendTaskFailureOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendTaskFailureOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendTaskFailureOutput>(create);
  static SendTaskFailureOutput? _defaultInstance;
}

class SendTaskHeartbeatInput extends $pb.GeneratedMessage {
  factory SendTaskHeartbeatInput({
    $core.String? tasktoken,
  }) {
    final result = create();
    if (tasktoken != null) result.tasktoken = tasktoken;
    return result;
  }

  SendTaskHeartbeatInput._();

  factory SendTaskHeartbeatInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendTaskHeartbeatInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendTaskHeartbeatInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(525325834, _omitFieldNames ? '' : 'tasktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskHeartbeatInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskHeartbeatInput copyWith(
          void Function(SendTaskHeartbeatInput) updates) =>
      super.copyWith((message) => updates(message as SendTaskHeartbeatInput))
          as SendTaskHeartbeatInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendTaskHeartbeatInput create() => SendTaskHeartbeatInput._();
  @$core.override
  SendTaskHeartbeatInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendTaskHeartbeatInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendTaskHeartbeatInput>(create);
  static SendTaskHeartbeatInput? _defaultInstance;

  @$pb.TagNumber(525325834)
  $core.String get tasktoken => $_getSZ(0);
  @$pb.TagNumber(525325834)
  set tasktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(525325834)
  $core.bool hasTasktoken() => $_has(0);
  @$pb.TagNumber(525325834)
  void clearTasktoken() => $_clearField(525325834);
}

class SendTaskHeartbeatOutput extends $pb.GeneratedMessage {
  factory SendTaskHeartbeatOutput() => create();

  SendTaskHeartbeatOutput._();

  factory SendTaskHeartbeatOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendTaskHeartbeatOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendTaskHeartbeatOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskHeartbeatOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskHeartbeatOutput copyWith(
          void Function(SendTaskHeartbeatOutput) updates) =>
      super.copyWith((message) => updates(message as SendTaskHeartbeatOutput))
          as SendTaskHeartbeatOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendTaskHeartbeatOutput create() => SendTaskHeartbeatOutput._();
  @$core.override
  SendTaskHeartbeatOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendTaskHeartbeatOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendTaskHeartbeatOutput>(create);
  static SendTaskHeartbeatOutput? _defaultInstance;
}

class SendTaskSuccessInput extends $pb.GeneratedMessage {
  factory SendTaskSuccessInput({
    $core.String? output,
    $core.String? tasktoken,
  }) {
    final result = create();
    if (output != null) result.output = output;
    if (tasktoken != null) result.tasktoken = tasktoken;
    return result;
  }

  SendTaskSuccessInput._();

  factory SendTaskSuccessInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendTaskSuccessInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendTaskSuccessInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..aOS(525325834, _omitFieldNames ? '' : 'tasktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskSuccessInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskSuccessInput copyWith(void Function(SendTaskSuccessInput) updates) =>
      super.copyWith((message) => updates(message as SendTaskSuccessInput))
          as SendTaskSuccessInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendTaskSuccessInput create() => SendTaskSuccessInput._();
  @$core.override
  SendTaskSuccessInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendTaskSuccessInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendTaskSuccessInput>(create);
  static SendTaskSuccessInput? _defaultInstance;

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(0);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(0, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(0);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);

  @$pb.TagNumber(525325834)
  $core.String get tasktoken => $_getSZ(1);
  @$pb.TagNumber(525325834)
  set tasktoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(525325834)
  $core.bool hasTasktoken() => $_has(1);
  @$pb.TagNumber(525325834)
  void clearTasktoken() => $_clearField(525325834);
}

class SendTaskSuccessOutput extends $pb.GeneratedMessage {
  factory SendTaskSuccessOutput() => create();

  SendTaskSuccessOutput._();

  factory SendTaskSuccessOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendTaskSuccessOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendTaskSuccessOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskSuccessOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendTaskSuccessOutput copyWith(
          void Function(SendTaskSuccessOutput) updates) =>
      super.copyWith((message) => updates(message as SendTaskSuccessOutput))
          as SendTaskSuccessOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendTaskSuccessOutput create() => SendTaskSuccessOutput._();
  @$core.override
  SendTaskSuccessOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendTaskSuccessOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendTaskSuccessOutput>(create);
  static SendTaskSuccessOutput? _defaultInstance;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class StartExecutionInput extends $pb.GeneratedMessage {
  factory StartExecutionInput({
    $core.String? traceheader,
    $core.String? name,
    $core.String? statemachinearn,
    $core.String? input,
  }) {
    final result = create();
    if (traceheader != null) result.traceheader = traceheader;
    if (name != null) result.name = name;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (input != null) result.input = input;
    return result;
  }

  StartExecutionInput._();

  factory StartExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(219960864, _omitFieldNames ? '' : 'traceheader')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartExecutionInput copyWith(void Function(StartExecutionInput) updates) =>
      super.copyWith((message) => updates(message as StartExecutionInput))
          as StartExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartExecutionInput create() => StartExecutionInput._();
  @$core.override
  StartExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartExecutionInput>(create);
  static StartExecutionInput? _defaultInstance;

  @$pb.TagNumber(219960864)
  $core.String get traceheader => $_getSZ(0);
  @$pb.TagNumber(219960864)
  set traceheader($core.String value) => $_setString(0, value);
  @$pb.TagNumber(219960864)
  $core.bool hasTraceheader() => $_has(0);
  @$pb.TagNumber(219960864)
  void clearTraceheader() => $_clearField(219960864);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(2);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(2);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(3);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(3, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(3);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);
}

class StartExecutionOutput extends $pb.GeneratedMessage {
  factory StartExecutionOutput({
    $core.String? executionarn,
    $core.String? startdate,
  }) {
    final result = create();
    if (executionarn != null) result.executionarn = executionarn;
    if (startdate != null) result.startdate = startdate;
    return result;
  }

  StartExecutionOutput._();

  factory StartExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..aOS(364840732, _omitFieldNames ? '' : 'startdate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartExecutionOutput copyWith(void Function(StartExecutionOutput) updates) =>
      super.copyWith((message) => updates(message as StartExecutionOutput))
          as StartExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartExecutionOutput create() => StartExecutionOutput._();
  @$core.override
  StartExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartExecutionOutput>(create);
  static StartExecutionOutput? _defaultInstance;

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(0);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(0);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);

  @$pb.TagNumber(364840732)
  $core.String get startdate => $_getSZ(1);
  @$pb.TagNumber(364840732)
  set startdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364840732)
  $core.bool hasStartdate() => $_has(1);
  @$pb.TagNumber(364840732)
  void clearStartdate() => $_clearField(364840732);
}

class StartSyncExecutionInput extends $pb.GeneratedMessage {
  factory StartSyncExecutionInput({
    IncludedData? includeddata,
    $core.String? traceheader,
    $core.String? name,
    $core.String? statemachinearn,
    $core.String? input,
  }) {
    final result = create();
    if (includeddata != null) result.includeddata = includeddata;
    if (traceheader != null) result.traceheader = traceheader;
    if (name != null) result.name = name;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (input != null) result.input = input;
    return result;
  }

  StartSyncExecutionInput._();

  factory StartSyncExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartSyncExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartSyncExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aE<IncludedData>(109719114, _omitFieldNames ? '' : 'includeddata',
        enumValues: IncludedData.values)
    ..aOS(219960864, _omitFieldNames ? '' : 'traceheader')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartSyncExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartSyncExecutionInput copyWith(
          void Function(StartSyncExecutionInput) updates) =>
      super.copyWith((message) => updates(message as StartSyncExecutionInput))
          as StartSyncExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartSyncExecutionInput create() => StartSyncExecutionInput._();
  @$core.override
  StartSyncExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartSyncExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartSyncExecutionInput>(create);
  static StartSyncExecutionInput? _defaultInstance;

  @$pb.TagNumber(109719114)
  IncludedData get includeddata => $_getN(0);
  @$pb.TagNumber(109719114)
  set includeddata(IncludedData value) => $_setField(109719114, value);
  @$pb.TagNumber(109719114)
  $core.bool hasIncludeddata() => $_has(0);
  @$pb.TagNumber(109719114)
  void clearIncludeddata() => $_clearField(109719114);

  @$pb.TagNumber(219960864)
  $core.String get traceheader => $_getSZ(1);
  @$pb.TagNumber(219960864)
  set traceheader($core.String value) => $_setString(1, value);
  @$pb.TagNumber(219960864)
  $core.bool hasTraceheader() => $_has(1);
  @$pb.TagNumber(219960864)
  void clearTraceheader() => $_clearField(219960864);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(3);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(3);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(4);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(4, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(4);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);
}

class StartSyncExecutionOutput extends $pb.GeneratedMessage {
  factory StartSyncExecutionOutput({
    $core.String? error,
    $core.String? cause,
    $core.String? stopdate,
    $core.String? traceheader,
    $core.String? name,
    BillingDetails? billingdetails,
    $core.String? executionarn,
    $core.String? startdate,
    $core.String? statemachinearn,
    CloudWatchEventsExecutionDataDetails? outputdetails,
    $core.String? output,
    $core.String? input,
    SyncExecutionStatus? status,
    CloudWatchEventsExecutionDataDetails? inputdetails,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    if (stopdate != null) result.stopdate = stopdate;
    if (traceheader != null) result.traceheader = traceheader;
    if (name != null) result.name = name;
    if (billingdetails != null) result.billingdetails = billingdetails;
    if (executionarn != null) result.executionarn = executionarn;
    if (startdate != null) result.startdate = startdate;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (outputdetails != null) result.outputdetails = outputdetails;
    if (output != null) result.output = output;
    if (input != null) result.input = input;
    if (status != null) result.status = status;
    if (inputdetails != null) result.inputdetails = inputdetails;
    return result;
  }

  StartSyncExecutionOutput._();

  factory StartSyncExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartSyncExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartSyncExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(180697434, _omitFieldNames ? '' : 'stopdate')
    ..aOS(219960864, _omitFieldNames ? '' : 'traceheader')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOM<BillingDetails>(270001723, _omitFieldNames ? '' : 'billingdetails',
        subBuilder: BillingDetails.create)
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..aOS(364840732, _omitFieldNames ? '' : 'startdate')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aOM<CloudWatchEventsExecutionDataDetails>(
        393734643, _omitFieldNames ? '' : 'outputdetails',
        subBuilder: CloudWatchEventsExecutionDataDetails.create)
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..aE<SyncExecutionStatus>(441153520, _omitFieldNames ? '' : 'status',
        enumValues: SyncExecutionStatus.values)
    ..aOM<CloudWatchEventsExecutionDataDetails>(
        452625788, _omitFieldNames ? '' : 'inputdetails',
        subBuilder: CloudWatchEventsExecutionDataDetails.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartSyncExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartSyncExecutionOutput copyWith(
          void Function(StartSyncExecutionOutput) updates) =>
      super.copyWith((message) => updates(message as StartSyncExecutionOutput))
          as StartSyncExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartSyncExecutionOutput create() => StartSyncExecutionOutput._();
  @$core.override
  StartSyncExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartSyncExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartSyncExecutionOutput>(create);
  static StartSyncExecutionOutput? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(180697434)
  $core.String get stopdate => $_getSZ(2);
  @$pb.TagNumber(180697434)
  set stopdate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(180697434)
  $core.bool hasStopdate() => $_has(2);
  @$pb.TagNumber(180697434)
  void clearStopdate() => $_clearField(180697434);

  @$pb.TagNumber(219960864)
  $core.String get traceheader => $_getSZ(3);
  @$pb.TagNumber(219960864)
  set traceheader($core.String value) => $_setString(3, value);
  @$pb.TagNumber(219960864)
  $core.bool hasTraceheader() => $_has(3);
  @$pb.TagNumber(219960864)
  void clearTraceheader() => $_clearField(219960864);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(270001723)
  BillingDetails get billingdetails => $_getN(5);
  @$pb.TagNumber(270001723)
  set billingdetails(BillingDetails value) => $_setField(270001723, value);
  @$pb.TagNumber(270001723)
  $core.bool hasBillingdetails() => $_has(5);
  @$pb.TagNumber(270001723)
  void clearBillingdetails() => $_clearField(270001723);
  @$pb.TagNumber(270001723)
  BillingDetails ensureBillingdetails() => $_ensure(5);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(6);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(6);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);

  @$pb.TagNumber(364840732)
  $core.String get startdate => $_getSZ(7);
  @$pb.TagNumber(364840732)
  set startdate($core.String value) => $_setString(7, value);
  @$pb.TagNumber(364840732)
  $core.bool hasStartdate() => $_has(7);
  @$pb.TagNumber(364840732)
  void clearStartdate() => $_clearField(364840732);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(8);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(8);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(393734643)
  CloudWatchEventsExecutionDataDetails get outputdetails => $_getN(9);
  @$pb.TagNumber(393734643)
  set outputdetails(CloudWatchEventsExecutionDataDetails value) =>
      $_setField(393734643, value);
  @$pb.TagNumber(393734643)
  $core.bool hasOutputdetails() => $_has(9);
  @$pb.TagNumber(393734643)
  void clearOutputdetails() => $_clearField(393734643);
  @$pb.TagNumber(393734643)
  CloudWatchEventsExecutionDataDetails ensureOutputdetails() => $_ensure(9);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(10);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(10, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(10);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(11);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(11, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(11);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);

  @$pb.TagNumber(441153520)
  SyncExecutionStatus get status => $_getN(12);
  @$pb.TagNumber(441153520)
  set status(SyncExecutionStatus value) => $_setField(441153520, value);
  @$pb.TagNumber(441153520)
  $core.bool hasStatus() => $_has(12);
  @$pb.TagNumber(441153520)
  void clearStatus() => $_clearField(441153520);

  @$pb.TagNumber(452625788)
  CloudWatchEventsExecutionDataDetails get inputdetails => $_getN(13);
  @$pb.TagNumber(452625788)
  set inputdetails(CloudWatchEventsExecutionDataDetails value) =>
      $_setField(452625788, value);
  @$pb.TagNumber(452625788)
  $core.bool hasInputdetails() => $_has(13);
  @$pb.TagNumber(452625788)
  void clearInputdetails() => $_clearField(452625788);
  @$pb.TagNumber(452625788)
  CloudWatchEventsExecutionDataDetails ensureInputdetails() => $_ensure(13);
}

class StateEnteredEventDetails extends $pb.GeneratedMessage {
  factory StateEnteredEventDetails({
    $core.String? name,
    $core.String? input,
    HistoryEventExecutionDataDetails? inputdetails,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (input != null) result.input = input;
    if (inputdetails != null) result.inputdetails = inputdetails;
    return result;
  }

  StateEnteredEventDetails._();

  factory StateEnteredEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateEnteredEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateEnteredEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..aOM<HistoryEventExecutionDataDetails>(
        452625788, _omitFieldNames ? '' : 'inputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateEnteredEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateEnteredEventDetails copyWith(
          void Function(StateEnteredEventDetails) updates) =>
      super.copyWith((message) => updates(message as StateEnteredEventDetails))
          as StateEnteredEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateEnteredEventDetails create() => StateEnteredEventDetails._();
  @$core.override
  StateEnteredEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateEnteredEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateEnteredEventDetails>(create);
  static StateEnteredEventDetails? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(1);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(1, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(1);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);

  @$pb.TagNumber(452625788)
  HistoryEventExecutionDataDetails get inputdetails => $_getN(2);
  @$pb.TagNumber(452625788)
  set inputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(452625788, value);
  @$pb.TagNumber(452625788)
  $core.bool hasInputdetails() => $_has(2);
  @$pb.TagNumber(452625788)
  void clearInputdetails() => $_clearField(452625788);
  @$pb.TagNumber(452625788)
  HistoryEventExecutionDataDetails ensureInputdetails() => $_ensure(2);
}

class StateExitedEventDetails extends $pb.GeneratedMessage {
  factory StateExitedEventDetails({
    $core.String? name,
    HistoryEventExecutionDataDetails? outputdetails,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        assignedvariables,
    $core.String? output,
    AssignedVariablesDetails? assignedvariablesdetails,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (outputdetails != null) result.outputdetails = outputdetails;
    if (assignedvariables != null)
      result.assignedvariables.addEntries(assignedvariables);
    if (output != null) result.output = output;
    if (assignedvariablesdetails != null)
      result.assignedvariablesdetails = assignedvariablesdetails;
    return result;
  }

  StateExitedEventDetails._();

  factory StateExitedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateExitedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateExitedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOM<HistoryEventExecutionDataDetails>(
        393734643, _omitFieldNames ? '' : 'outputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..m<$core.String, $core.String>(
        411875019, _omitFieldNames ? '' : 'assignedvariables',
        entryClassName: 'StateExitedEventDetails.AssignedvariablesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('states'))
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..aOM<AssignedVariablesDetails>(
        509183609, _omitFieldNames ? '' : 'assignedvariablesdetails',
        subBuilder: AssignedVariablesDetails.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateExitedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateExitedEventDetails copyWith(
          void Function(StateExitedEventDetails) updates) =>
      super.copyWith((message) => updates(message as StateExitedEventDetails))
          as StateExitedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateExitedEventDetails create() => StateExitedEventDetails._();
  @$core.override
  StateExitedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateExitedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateExitedEventDetails>(create);
  static StateExitedEventDetails? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails get outputdetails => $_getN(1);
  @$pb.TagNumber(393734643)
  set outputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(393734643, value);
  @$pb.TagNumber(393734643)
  $core.bool hasOutputdetails() => $_has(1);
  @$pb.TagNumber(393734643)
  void clearOutputdetails() => $_clearField(393734643);
  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails ensureOutputdetails() => $_ensure(1);

  @$pb.TagNumber(411875019)
  $pb.PbMap<$core.String, $core.String> get assignedvariables => $_getMap(2);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(3);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(3, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(3);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);

  @$pb.TagNumber(509183609)
  AssignedVariablesDetails get assignedvariablesdetails => $_getN(4);
  @$pb.TagNumber(509183609)
  set assignedvariablesdetails(AssignedVariablesDetails value) =>
      $_setField(509183609, value);
  @$pb.TagNumber(509183609)
  $core.bool hasAssignedvariablesdetails() => $_has(4);
  @$pb.TagNumber(509183609)
  void clearAssignedvariablesdetails() => $_clearField(509183609);
  @$pb.TagNumber(509183609)
  AssignedVariablesDetails ensureAssignedvariablesdetails() => $_ensure(4);
}

class StateMachineAliasListItem extends $pb.GeneratedMessage {
  factory StateMachineAliasListItem({
    $core.String? creationdate,
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (creationdate != null) result.creationdate = creationdate;
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  StateMachineAliasListItem._();

  factory StateMachineAliasListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateMachineAliasListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateMachineAliasListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineAliasListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineAliasListItem copyWith(
          void Function(StateMachineAliasListItem) updates) =>
      super.copyWith((message) => updates(message as StateMachineAliasListItem))
          as StateMachineAliasListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateMachineAliasListItem create() => StateMachineAliasListItem._();
  @$core.override
  StateMachineAliasListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateMachineAliasListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateMachineAliasListItem>(create);
  static StateMachineAliasListItem? _defaultInstance;

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(0);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(0);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(1);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(1);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class StateMachineAlreadyExists extends $pb.GeneratedMessage {
  factory StateMachineAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  StateMachineAlreadyExists._();

  factory StateMachineAlreadyExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateMachineAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateMachineAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineAlreadyExists copyWith(
          void Function(StateMachineAlreadyExists) updates) =>
      super.copyWith((message) => updates(message as StateMachineAlreadyExists))
          as StateMachineAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateMachineAlreadyExists create() => StateMachineAlreadyExists._();
  @$core.override
  StateMachineAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateMachineAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateMachineAlreadyExists>(create);
  static StateMachineAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class StateMachineDeleting extends $pb.GeneratedMessage {
  factory StateMachineDeleting({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  StateMachineDeleting._();

  factory StateMachineDeleting.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateMachineDeleting.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateMachineDeleting',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineDeleting clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineDeleting copyWith(void Function(StateMachineDeleting) updates) =>
      super.copyWith((message) => updates(message as StateMachineDeleting))
          as StateMachineDeleting;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateMachineDeleting create() => StateMachineDeleting._();
  @$core.override
  StateMachineDeleting createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateMachineDeleting getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateMachineDeleting>(create);
  static StateMachineDeleting? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class StateMachineDoesNotExist extends $pb.GeneratedMessage {
  factory StateMachineDoesNotExist({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  StateMachineDoesNotExist._();

  factory StateMachineDoesNotExist.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateMachineDoesNotExist.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateMachineDoesNotExist',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineDoesNotExist clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineDoesNotExist copyWith(
          void Function(StateMachineDoesNotExist) updates) =>
      super.copyWith((message) => updates(message as StateMachineDoesNotExist))
          as StateMachineDoesNotExist;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateMachineDoesNotExist create() => StateMachineDoesNotExist._();
  @$core.override
  StateMachineDoesNotExist createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateMachineDoesNotExist getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateMachineDoesNotExist>(create);
  static StateMachineDoesNotExist? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class StateMachineLimitExceeded extends $pb.GeneratedMessage {
  factory StateMachineLimitExceeded({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  StateMachineLimitExceeded._();

  factory StateMachineLimitExceeded.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateMachineLimitExceeded.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateMachineLimitExceeded',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineLimitExceeded clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineLimitExceeded copyWith(
          void Function(StateMachineLimitExceeded) updates) =>
      super.copyWith((message) => updates(message as StateMachineLimitExceeded))
          as StateMachineLimitExceeded;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateMachineLimitExceeded create() => StateMachineLimitExceeded._();
  @$core.override
  StateMachineLimitExceeded createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateMachineLimitExceeded getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateMachineLimitExceeded>(create);
  static StateMachineLimitExceeded? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class StateMachineListItem extends $pb.GeneratedMessage {
  factory StateMachineListItem({
    $core.String? name,
    $core.String? creationdate,
    StateMachineType? type,
    $core.String? statemachinearn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (type != null) result.type = type;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    return result;
  }

  StateMachineListItem._();

  factory StateMachineListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateMachineListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateMachineListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..aE<StateMachineType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: StateMachineType.values)
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineListItem copyWith(void Function(StateMachineListItem) updates) =>
      super.copyWith((message) => updates(message as StateMachineListItem))
          as StateMachineListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateMachineListItem create() => StateMachineListItem._();
  @$core.override
  StateMachineListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateMachineListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateMachineListItem>(create);
  static StateMachineListItem? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(1);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(1);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);

  @$pb.TagNumber(287830350)
  StateMachineType get type => $_getN(2);
  @$pb.TagNumber(287830350)
  set type(StateMachineType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(3);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(3);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);
}

class StateMachineTypeNotSupported extends $pb.GeneratedMessage {
  factory StateMachineTypeNotSupported({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  StateMachineTypeNotSupported._();

  factory StateMachineTypeNotSupported.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateMachineTypeNotSupported.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateMachineTypeNotSupported',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineTypeNotSupported clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineTypeNotSupported copyWith(
          void Function(StateMachineTypeNotSupported) updates) =>
      super.copyWith(
              (message) => updates(message as StateMachineTypeNotSupported))
          as StateMachineTypeNotSupported;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateMachineTypeNotSupported create() =>
      StateMachineTypeNotSupported._();
  @$core.override
  StateMachineTypeNotSupported createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateMachineTypeNotSupported getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateMachineTypeNotSupported>(create);
  static StateMachineTypeNotSupported? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class StateMachineVersionListItem extends $pb.GeneratedMessage {
  factory StateMachineVersionListItem({
    $core.String? statemachineversionarn,
    $core.String? creationdate,
  }) {
    final result = create();
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    if (creationdate != null) result.creationdate = creationdate;
    return result;
  }

  StateMachineVersionListItem._();

  factory StateMachineVersionListItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StateMachineVersionListItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StateMachineVersionListItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..aOS(238315265, _omitFieldNames ? '' : 'creationdate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineVersionListItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StateMachineVersionListItem copyWith(
          void Function(StateMachineVersionListItem) updates) =>
      super.copyWith(
              (message) => updates(message as StateMachineVersionListItem))
          as StateMachineVersionListItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StateMachineVersionListItem create() =>
      StateMachineVersionListItem._();
  @$core.override
  StateMachineVersionListItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StateMachineVersionListItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StateMachineVersionListItem>(create);
  static StateMachineVersionListItem? _defaultInstance;

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(0);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(0);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);

  @$pb.TagNumber(238315265)
  $core.String get creationdate => $_getSZ(1);
  @$pb.TagNumber(238315265)
  set creationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(238315265)
  $core.bool hasCreationdate() => $_has(1);
  @$pb.TagNumber(238315265)
  void clearCreationdate() => $_clearField(238315265);
}

class StopExecutionInput extends $pb.GeneratedMessage {
  factory StopExecutionInput({
    $core.String? error,
    $core.String? cause,
    $core.String? executionarn,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    if (executionarn != null) result.executionarn = executionarn;
    return result;
  }

  StopExecutionInput._();

  factory StopExecutionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopExecutionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopExecutionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(314526573, _omitFieldNames ? '' : 'executionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopExecutionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopExecutionInput copyWith(void Function(StopExecutionInput) updates) =>
      super.copyWith((message) => updates(message as StopExecutionInput))
          as StopExecutionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopExecutionInput create() => StopExecutionInput._();
  @$core.override
  StopExecutionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopExecutionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopExecutionInput>(create);
  static StopExecutionInput? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(1);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(1, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(1);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(314526573)
  $core.String get executionarn => $_getSZ(2);
  @$pb.TagNumber(314526573)
  set executionarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(314526573)
  $core.bool hasExecutionarn() => $_has(2);
  @$pb.TagNumber(314526573)
  void clearExecutionarn() => $_clearField(314526573);
}

class StopExecutionOutput extends $pb.GeneratedMessage {
  factory StopExecutionOutput({
    $core.String? stopdate,
  }) {
    final result = create();
    if (stopdate != null) result.stopdate = stopdate;
    return result;
  }

  StopExecutionOutput._();

  factory StopExecutionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopExecutionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopExecutionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(180697434, _omitFieldNames ? '' : 'stopdate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopExecutionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopExecutionOutput copyWith(void Function(StopExecutionOutput) updates) =>
      super.copyWith((message) => updates(message as StopExecutionOutput))
          as StopExecutionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopExecutionOutput create() => StopExecutionOutput._();
  @$core.override
  StopExecutionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopExecutionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopExecutionOutput>(create);
  static StopExecutionOutput? _defaultInstance;

  @$pb.TagNumber(180697434)
  $core.String get stopdate => $_getSZ(0);
  @$pb.TagNumber(180697434)
  set stopdate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(180697434)
  $core.bool hasStopdate() => $_has(0);
  @$pb.TagNumber(180697434)
  void clearStopdate() => $_clearField(180697434);
}

class Tag extends $pb.GeneratedMessage {
  factory Tag({
    $core.String? value,
    $core.String? key,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (key != null) result.key = key;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(39769035, _omitFieldNames ? '' : 'value')
    ..aOS(135645293, _omitFieldNames ? '' : 'key')
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

  @$pb.TagNumber(39769035)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(39769035)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39769035)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(39769035)
  void clearValue() => $_clearField(39769035);

  @$pb.TagNumber(135645293)
  $core.String get key => $_getSZ(1);
  @$pb.TagNumber(135645293)
  set key($core.String value) => $_setString(1, value);
  @$pb.TagNumber(135645293)
  $core.bool hasKey() => $_has(1);
  @$pb.TagNumber(135645293)
  void clearKey() => $_clearField(135645293);
}

class TagResourceInput extends $pb.GeneratedMessage {
  factory TagResourceInput({
    $core.String? resourcearn,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  TagResourceInput._();

  factory TagResourceInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagResourceInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagResourceInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(67806797, _omitFieldNames ? '' : 'resourcearn')
    ..pPM<Tag>(337046433, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceInput copyWith(void Function(TagResourceInput) updates) =>
      super.copyWith((message) => updates(message as TagResourceInput))
          as TagResourceInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagResourceInput create() => TagResourceInput._();
  @$core.override
  TagResourceInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagResourceInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagResourceInput>(create);
  static TagResourceInput? _defaultInstance;

  @$pb.TagNumber(67806797)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(67806797)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(67806797)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(67806797)
  void clearResourcearn() => $_clearField(67806797);

  @$pb.TagNumber(337046433)
  $pb.PbList<Tag> get tags => $_getList(1);
}

class TagResourceOutput extends $pb.GeneratedMessage {
  factory TagResourceOutput() => create();

  TagResourceOutput._();

  factory TagResourceOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagResourceOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagResourceOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceOutput copyWith(void Function(TagResourceOutput) updates) =>
      super.copyWith((message) => updates(message as TagResourceOutput))
          as TagResourceOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagResourceOutput create() => TagResourceOutput._();
  @$core.override
  TagResourceOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagResourceOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagResourceOutput>(create);
  static TagResourceOutput? _defaultInstance;
}

class TaskCredentials extends $pb.GeneratedMessage {
  factory TaskCredentials({
    $core.String? rolearn,
  }) {
    final result = create();
    if (rolearn != null) result.rolearn = rolearn;
    return result;
  }

  TaskCredentials._();

  factory TaskCredentials.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskCredentials.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskCredentials',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(170019745, _omitFieldNames ? '' : 'rolearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskCredentials clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskCredentials copyWith(void Function(TaskCredentials) updates) =>
      super.copyWith((message) => updates(message as TaskCredentials))
          as TaskCredentials;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskCredentials create() => TaskCredentials._();
  @$core.override
  TaskCredentials createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskCredentials getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskCredentials>(create);
  static TaskCredentials? _defaultInstance;

  @$pb.TagNumber(170019745)
  $core.String get rolearn => $_getSZ(0);
  @$pb.TagNumber(170019745)
  set rolearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(170019745)
  $core.bool hasRolearn() => $_has(0);
  @$pb.TagNumber(170019745)
  void clearRolearn() => $_clearField(170019745);
}

class TaskDoesNotExist extends $pb.GeneratedMessage {
  factory TaskDoesNotExist({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TaskDoesNotExist._();

  factory TaskDoesNotExist.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskDoesNotExist.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskDoesNotExist',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskDoesNotExist clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskDoesNotExist copyWith(void Function(TaskDoesNotExist) updates) =>
      super.copyWith((message) => updates(message as TaskDoesNotExist))
          as TaskDoesNotExist;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskDoesNotExist create() => TaskDoesNotExist._();
  @$core.override
  TaskDoesNotExist createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskDoesNotExist getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskDoesNotExist>(create);
  static TaskDoesNotExist? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TaskFailedEventDetails extends $pb.GeneratedMessage {
  factory TaskFailedEventDetails({
    $core.String? resourcetype,
    $core.String? error,
    $core.String? cause,
    $core.String? resource,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    if (resource != null) result.resource = resource;
    return result;
  }

  TaskFailedEventDetails._();

  factory TaskFailedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(7604990, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskFailedEventDetails copyWith(
          void Function(TaskFailedEventDetails) updates) =>
      super.copyWith((message) => updates(message as TaskFailedEventDetails))
          as TaskFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskFailedEventDetails create() => TaskFailedEventDetails._();
  @$core.override
  TaskFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskFailedEventDetails>(create);
  static TaskFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(7604990)
  $core.String get resourcetype => $_getSZ(0);
  @$pb.TagNumber(7604990)
  set resourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7604990)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(7604990)
  void clearResourcetype() => $_clearField(7604990);

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(1);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(1, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(1);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(2);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(2, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(2);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(3);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(3, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(3);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);
}

class TaskScheduledEventDetails extends $pb.GeneratedMessage {
  factory TaskScheduledEventDetails({
    $core.String? resourcetype,
    $core.String? region,
    $fixnum.Int64? heartbeatinseconds,
    $core.String? parameters,
    $core.String? resource,
    TaskCredentials? taskcredentials,
    $fixnum.Int64? timeoutinseconds,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (region != null) result.region = region;
    if (heartbeatinseconds != null)
      result.heartbeatinseconds = heartbeatinseconds;
    if (parameters != null) result.parameters = parameters;
    if (resource != null) result.resource = resource;
    if (taskcredentials != null) result.taskcredentials = taskcredentials;
    if (timeoutinseconds != null) result.timeoutinseconds = timeoutinseconds;
    return result;
  }

  TaskScheduledEventDetails._();

  factory TaskScheduledEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskScheduledEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskScheduledEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(7604990, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(52100734, _omitFieldNames ? '' : 'region')
    ..aInt64(125718754, _omitFieldNames ? '' : 'heartbeatinseconds')
    ..aOS(145043162, _omitFieldNames ? '' : 'parameters')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..aOM<TaskCredentials>(257843259, _omitFieldNames ? '' : 'taskcredentials',
        subBuilder: TaskCredentials.create)
    ..aInt64(472710197, _omitFieldNames ? '' : 'timeoutinseconds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskScheduledEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskScheduledEventDetails copyWith(
          void Function(TaskScheduledEventDetails) updates) =>
      super.copyWith((message) => updates(message as TaskScheduledEventDetails))
          as TaskScheduledEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskScheduledEventDetails create() => TaskScheduledEventDetails._();
  @$core.override
  TaskScheduledEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskScheduledEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskScheduledEventDetails>(create);
  static TaskScheduledEventDetails? _defaultInstance;

  @$pb.TagNumber(7604990)
  $core.String get resourcetype => $_getSZ(0);
  @$pb.TagNumber(7604990)
  set resourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7604990)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(7604990)
  void clearResourcetype() => $_clearField(7604990);

  @$pb.TagNumber(52100734)
  $core.String get region => $_getSZ(1);
  @$pb.TagNumber(52100734)
  set region($core.String value) => $_setString(1, value);
  @$pb.TagNumber(52100734)
  $core.bool hasRegion() => $_has(1);
  @$pb.TagNumber(52100734)
  void clearRegion() => $_clearField(52100734);

  @$pb.TagNumber(125718754)
  $fixnum.Int64 get heartbeatinseconds => $_getI64(2);
  @$pb.TagNumber(125718754)
  set heartbeatinseconds($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(125718754)
  $core.bool hasHeartbeatinseconds() => $_has(2);
  @$pb.TagNumber(125718754)
  void clearHeartbeatinseconds() => $_clearField(125718754);

  @$pb.TagNumber(145043162)
  $core.String get parameters => $_getSZ(3);
  @$pb.TagNumber(145043162)
  set parameters($core.String value) => $_setString(3, value);
  @$pb.TagNumber(145043162)
  $core.bool hasParameters() => $_has(3);
  @$pb.TagNumber(145043162)
  void clearParameters() => $_clearField(145043162);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(4);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(4, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(4);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);

  @$pb.TagNumber(257843259)
  TaskCredentials get taskcredentials => $_getN(5);
  @$pb.TagNumber(257843259)
  set taskcredentials(TaskCredentials value) => $_setField(257843259, value);
  @$pb.TagNumber(257843259)
  $core.bool hasTaskcredentials() => $_has(5);
  @$pb.TagNumber(257843259)
  void clearTaskcredentials() => $_clearField(257843259);
  @$pb.TagNumber(257843259)
  TaskCredentials ensureTaskcredentials() => $_ensure(5);

  @$pb.TagNumber(472710197)
  $fixnum.Int64 get timeoutinseconds => $_getI64(6);
  @$pb.TagNumber(472710197)
  set timeoutinseconds($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(472710197)
  $core.bool hasTimeoutinseconds() => $_has(6);
  @$pb.TagNumber(472710197)
  void clearTimeoutinseconds() => $_clearField(472710197);
}

class TaskStartFailedEventDetails extends $pb.GeneratedMessage {
  factory TaskStartFailedEventDetails({
    $core.String? resourcetype,
    $core.String? error,
    $core.String? cause,
    $core.String? resource,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    if (resource != null) result.resource = resource;
    return result;
  }

  TaskStartFailedEventDetails._();

  factory TaskStartFailedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskStartFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskStartFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(7604990, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskStartFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskStartFailedEventDetails copyWith(
          void Function(TaskStartFailedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as TaskStartFailedEventDetails))
          as TaskStartFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskStartFailedEventDetails create() =>
      TaskStartFailedEventDetails._();
  @$core.override
  TaskStartFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskStartFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskStartFailedEventDetails>(create);
  static TaskStartFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(7604990)
  $core.String get resourcetype => $_getSZ(0);
  @$pb.TagNumber(7604990)
  set resourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7604990)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(7604990)
  void clearResourcetype() => $_clearField(7604990);

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(1);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(1, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(1);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(2);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(2, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(2);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(3);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(3, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(3);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);
}

class TaskStartedEventDetails extends $pb.GeneratedMessage {
  factory TaskStartedEventDetails({
    $core.String? resourcetype,
    $core.String? resource,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (resource != null) result.resource = resource;
    return result;
  }

  TaskStartedEventDetails._();

  factory TaskStartedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskStartedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskStartedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(7604990, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskStartedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskStartedEventDetails copyWith(
          void Function(TaskStartedEventDetails) updates) =>
      super.copyWith((message) => updates(message as TaskStartedEventDetails))
          as TaskStartedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskStartedEventDetails create() => TaskStartedEventDetails._();
  @$core.override
  TaskStartedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskStartedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskStartedEventDetails>(create);
  static TaskStartedEventDetails? _defaultInstance;

  @$pb.TagNumber(7604990)
  $core.String get resourcetype => $_getSZ(0);
  @$pb.TagNumber(7604990)
  set resourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7604990)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(7604990)
  void clearResourcetype() => $_clearField(7604990);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(1);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(1, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(1);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);
}

class TaskSubmitFailedEventDetails extends $pb.GeneratedMessage {
  factory TaskSubmitFailedEventDetails({
    $core.String? resourcetype,
    $core.String? error,
    $core.String? cause,
    $core.String? resource,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    if (resource != null) result.resource = resource;
    return result;
  }

  TaskSubmitFailedEventDetails._();

  factory TaskSubmitFailedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskSubmitFailedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskSubmitFailedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(7604990, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskSubmitFailedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskSubmitFailedEventDetails copyWith(
          void Function(TaskSubmitFailedEventDetails) updates) =>
      super.copyWith(
              (message) => updates(message as TaskSubmitFailedEventDetails))
          as TaskSubmitFailedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskSubmitFailedEventDetails create() =>
      TaskSubmitFailedEventDetails._();
  @$core.override
  TaskSubmitFailedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskSubmitFailedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskSubmitFailedEventDetails>(create);
  static TaskSubmitFailedEventDetails? _defaultInstance;

  @$pb.TagNumber(7604990)
  $core.String get resourcetype => $_getSZ(0);
  @$pb.TagNumber(7604990)
  set resourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7604990)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(7604990)
  void clearResourcetype() => $_clearField(7604990);

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(1);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(1, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(1);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(2);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(2, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(2);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(3);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(3, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(3);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);
}

class TaskSubmittedEventDetails extends $pb.GeneratedMessage {
  factory TaskSubmittedEventDetails({
    $core.String? resourcetype,
    $core.String? resource,
    HistoryEventExecutionDataDetails? outputdetails,
    $core.String? output,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (resource != null) result.resource = resource;
    if (outputdetails != null) result.outputdetails = outputdetails;
    if (output != null) result.output = output;
    return result;
  }

  TaskSubmittedEventDetails._();

  factory TaskSubmittedEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskSubmittedEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskSubmittedEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(7604990, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..aOM<HistoryEventExecutionDataDetails>(
        393734643, _omitFieldNames ? '' : 'outputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskSubmittedEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskSubmittedEventDetails copyWith(
          void Function(TaskSubmittedEventDetails) updates) =>
      super.copyWith((message) => updates(message as TaskSubmittedEventDetails))
          as TaskSubmittedEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskSubmittedEventDetails create() => TaskSubmittedEventDetails._();
  @$core.override
  TaskSubmittedEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskSubmittedEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskSubmittedEventDetails>(create);
  static TaskSubmittedEventDetails? _defaultInstance;

  @$pb.TagNumber(7604990)
  $core.String get resourcetype => $_getSZ(0);
  @$pb.TagNumber(7604990)
  set resourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7604990)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(7604990)
  void clearResourcetype() => $_clearField(7604990);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(1);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(1, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(1);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);

  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails get outputdetails => $_getN(2);
  @$pb.TagNumber(393734643)
  set outputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(393734643, value);
  @$pb.TagNumber(393734643)
  $core.bool hasOutputdetails() => $_has(2);
  @$pb.TagNumber(393734643)
  void clearOutputdetails() => $_clearField(393734643);
  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails ensureOutputdetails() => $_ensure(2);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(3);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(3, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(3);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);
}

class TaskSucceededEventDetails extends $pb.GeneratedMessage {
  factory TaskSucceededEventDetails({
    $core.String? resourcetype,
    $core.String? resource,
    HistoryEventExecutionDataDetails? outputdetails,
    $core.String? output,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (resource != null) result.resource = resource;
    if (outputdetails != null) result.outputdetails = outputdetails;
    if (output != null) result.output = output;
    return result;
  }

  TaskSucceededEventDetails._();

  factory TaskSucceededEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskSucceededEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskSucceededEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(7604990, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..aOM<HistoryEventExecutionDataDetails>(
        393734643, _omitFieldNames ? '' : 'outputdetails',
        subBuilder: HistoryEventExecutionDataDetails.create)
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskSucceededEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskSucceededEventDetails copyWith(
          void Function(TaskSucceededEventDetails) updates) =>
      super.copyWith((message) => updates(message as TaskSucceededEventDetails))
          as TaskSucceededEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskSucceededEventDetails create() => TaskSucceededEventDetails._();
  @$core.override
  TaskSucceededEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskSucceededEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskSucceededEventDetails>(create);
  static TaskSucceededEventDetails? _defaultInstance;

  @$pb.TagNumber(7604990)
  $core.String get resourcetype => $_getSZ(0);
  @$pb.TagNumber(7604990)
  set resourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7604990)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(7604990)
  void clearResourcetype() => $_clearField(7604990);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(1);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(1, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(1);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);

  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails get outputdetails => $_getN(2);
  @$pb.TagNumber(393734643)
  set outputdetails(HistoryEventExecutionDataDetails value) =>
      $_setField(393734643, value);
  @$pb.TagNumber(393734643)
  $core.bool hasOutputdetails() => $_has(2);
  @$pb.TagNumber(393734643)
  void clearOutputdetails() => $_clearField(393734643);
  @$pb.TagNumber(393734643)
  HistoryEventExecutionDataDetails ensureOutputdetails() => $_ensure(2);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(3);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(3, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(3);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);
}

class TaskTimedOut extends $pb.GeneratedMessage {
  factory TaskTimedOut({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TaskTimedOut._();

  factory TaskTimedOut.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskTimedOut.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskTimedOut',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskTimedOut clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskTimedOut copyWith(void Function(TaskTimedOut) updates) =>
      super.copyWith((message) => updates(message as TaskTimedOut))
          as TaskTimedOut;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskTimedOut create() => TaskTimedOut._();
  @$core.override
  TaskTimedOut createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskTimedOut getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskTimedOut>(create);
  static TaskTimedOut? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TaskTimedOutEventDetails extends $pb.GeneratedMessage {
  factory TaskTimedOutEventDetails({
    $core.String? resourcetype,
    $core.String? error,
    $core.String? cause,
    $core.String? resource,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (error != null) result.error = error;
    if (cause != null) result.cause = cause;
    if (resource != null) result.resource = resource;
    return result;
  }

  TaskTimedOutEventDetails._();

  factory TaskTimedOutEventDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TaskTimedOutEventDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TaskTimedOutEventDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(7604990, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskTimedOutEventDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TaskTimedOutEventDetails copyWith(
          void Function(TaskTimedOutEventDetails) updates) =>
      super.copyWith((message) => updates(message as TaskTimedOutEventDetails))
          as TaskTimedOutEventDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TaskTimedOutEventDetails create() => TaskTimedOutEventDetails._();
  @$core.override
  TaskTimedOutEventDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TaskTimedOutEventDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TaskTimedOutEventDetails>(create);
  static TaskTimedOutEventDetails? _defaultInstance;

  @$pb.TagNumber(7604990)
  $core.String get resourcetype => $_getSZ(0);
  @$pb.TagNumber(7604990)
  set resourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7604990)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(7604990)
  void clearResourcetype() => $_clearField(7604990);

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(1);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(1, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(1);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(2);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(2, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(2);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(3);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(3, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(3);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);
}

class TestStateConfiguration extends $pb.GeneratedMessage {
  factory TestStateConfiguration({
    $core.int? retrierretrycount,
    $core.String? errorcausedbystate,
    $core.int? mapiterationfailurecount,
    $core.String? mapitemreaderdata,
  }) {
    final result = create();
    if (retrierretrycount != null) result.retrierretrycount = retrierretrycount;
    if (errorcausedbystate != null)
      result.errorcausedbystate = errorcausedbystate;
    if (mapiterationfailurecount != null)
      result.mapiterationfailurecount = mapiterationfailurecount;
    if (mapitemreaderdata != null) result.mapitemreaderdata = mapitemreaderdata;
    return result;
  }

  TestStateConfiguration._();

  factory TestStateConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestStateConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestStateConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aI(275735648, _omitFieldNames ? '' : 'retrierretrycount')
    ..aOS(313873103, _omitFieldNames ? '' : 'errorcausedbystate')
    ..aI(419587762, _omitFieldNames ? '' : 'mapiterationfailurecount')
    ..aOS(503260646, _omitFieldNames ? '' : 'mapitemreaderdata')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestStateConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestStateConfiguration copyWith(
          void Function(TestStateConfiguration) updates) =>
      super.copyWith((message) => updates(message as TestStateConfiguration))
          as TestStateConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestStateConfiguration create() => TestStateConfiguration._();
  @$core.override
  TestStateConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestStateConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestStateConfiguration>(create);
  static TestStateConfiguration? _defaultInstance;

  @$pb.TagNumber(275735648)
  $core.int get retrierretrycount => $_getIZ(0);
  @$pb.TagNumber(275735648)
  set retrierretrycount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(275735648)
  $core.bool hasRetrierretrycount() => $_has(0);
  @$pb.TagNumber(275735648)
  void clearRetrierretrycount() => $_clearField(275735648);

  @$pb.TagNumber(313873103)
  $core.String get errorcausedbystate => $_getSZ(1);
  @$pb.TagNumber(313873103)
  set errorcausedbystate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(313873103)
  $core.bool hasErrorcausedbystate() => $_has(1);
  @$pb.TagNumber(313873103)
  void clearErrorcausedbystate() => $_clearField(313873103);

  @$pb.TagNumber(419587762)
  $core.int get mapiterationfailurecount => $_getIZ(2);
  @$pb.TagNumber(419587762)
  set mapiterationfailurecount($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(419587762)
  $core.bool hasMapiterationfailurecount() => $_has(2);
  @$pb.TagNumber(419587762)
  void clearMapiterationfailurecount() => $_clearField(419587762);

  @$pb.TagNumber(503260646)
  $core.String get mapitemreaderdata => $_getSZ(3);
  @$pb.TagNumber(503260646)
  set mapitemreaderdata($core.String value) => $_setString(3, value);
  @$pb.TagNumber(503260646)
  $core.bool hasMapitemreaderdata() => $_has(3);
  @$pb.TagNumber(503260646)
  void clearMapitemreaderdata() => $_clearField(503260646);
}

class TestStateInput extends $pb.GeneratedMessage {
  factory TestStateInput({
    TestStateConfiguration? stateconfiguration,
    $core.String? definition,
    $core.String? variables,
    $core.String? rolearn,
    $core.String? context,
    MockInput? mock,
    $core.String? statename,
    InspectionLevel? inspectionlevel,
    $core.bool? revealsecrets,
    $core.String? input,
  }) {
    final result = create();
    if (stateconfiguration != null)
      result.stateconfiguration = stateconfiguration;
    if (definition != null) result.definition = definition;
    if (variables != null) result.variables = variables;
    if (rolearn != null) result.rolearn = rolearn;
    if (context != null) result.context = context;
    if (mock != null) result.mock = mock;
    if (statename != null) result.statename = statename;
    if (inspectionlevel != null) result.inspectionlevel = inspectionlevel;
    if (revealsecrets != null) result.revealsecrets = revealsecrets;
    if (input != null) result.input = input;
    return result;
  }

  TestStateInput._();

  factory TestStateInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestStateInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestStateInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOM<TestStateConfiguration>(
        17002877, _omitFieldNames ? '' : 'stateconfiguration',
        subBuilder: TestStateConfiguration.create)
    ..aOS(68443297, _omitFieldNames ? '' : 'definition')
    ..aOS(162226883, _omitFieldNames ? '' : 'variables')
    ..aOS(170019745, _omitFieldNames ? '' : 'rolearn')
    ..aOS(210178173, _omitFieldNames ? '' : 'context')
    ..aOM<MockInput>(242883628, _omitFieldNames ? '' : 'mock',
        subBuilder: MockInput.create)
    ..aOS(270657590, _omitFieldNames ? '' : 'statename')
    ..aE<InspectionLevel>(277169476, _omitFieldNames ? '' : 'inspectionlevel',
        enumValues: InspectionLevel.values)
    ..aOB(351839742, _omitFieldNames ? '' : 'revealsecrets')
    ..aOS(433614716, _omitFieldNames ? '' : 'input')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestStateInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestStateInput copyWith(void Function(TestStateInput) updates) =>
      super.copyWith((message) => updates(message as TestStateInput))
          as TestStateInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestStateInput create() => TestStateInput._();
  @$core.override
  TestStateInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestStateInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestStateInput>(create);
  static TestStateInput? _defaultInstance;

  @$pb.TagNumber(17002877)
  TestStateConfiguration get stateconfiguration => $_getN(0);
  @$pb.TagNumber(17002877)
  set stateconfiguration(TestStateConfiguration value) =>
      $_setField(17002877, value);
  @$pb.TagNumber(17002877)
  $core.bool hasStateconfiguration() => $_has(0);
  @$pb.TagNumber(17002877)
  void clearStateconfiguration() => $_clearField(17002877);
  @$pb.TagNumber(17002877)
  TestStateConfiguration ensureStateconfiguration() => $_ensure(0);

  @$pb.TagNumber(68443297)
  $core.String get definition => $_getSZ(1);
  @$pb.TagNumber(68443297)
  set definition($core.String value) => $_setString(1, value);
  @$pb.TagNumber(68443297)
  $core.bool hasDefinition() => $_has(1);
  @$pb.TagNumber(68443297)
  void clearDefinition() => $_clearField(68443297);

  @$pb.TagNumber(162226883)
  $core.String get variables => $_getSZ(2);
  @$pb.TagNumber(162226883)
  set variables($core.String value) => $_setString(2, value);
  @$pb.TagNumber(162226883)
  $core.bool hasVariables() => $_has(2);
  @$pb.TagNumber(162226883)
  void clearVariables() => $_clearField(162226883);

  @$pb.TagNumber(170019745)
  $core.String get rolearn => $_getSZ(3);
  @$pb.TagNumber(170019745)
  set rolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(170019745)
  $core.bool hasRolearn() => $_has(3);
  @$pb.TagNumber(170019745)
  void clearRolearn() => $_clearField(170019745);

  @$pb.TagNumber(210178173)
  $core.String get context => $_getSZ(4);
  @$pb.TagNumber(210178173)
  set context($core.String value) => $_setString(4, value);
  @$pb.TagNumber(210178173)
  $core.bool hasContext() => $_has(4);
  @$pb.TagNumber(210178173)
  void clearContext() => $_clearField(210178173);

  @$pb.TagNumber(242883628)
  MockInput get mock => $_getN(5);
  @$pb.TagNumber(242883628)
  set mock(MockInput value) => $_setField(242883628, value);
  @$pb.TagNumber(242883628)
  $core.bool hasMock() => $_has(5);
  @$pb.TagNumber(242883628)
  void clearMock() => $_clearField(242883628);
  @$pb.TagNumber(242883628)
  MockInput ensureMock() => $_ensure(5);

  @$pb.TagNumber(270657590)
  $core.String get statename => $_getSZ(6);
  @$pb.TagNumber(270657590)
  set statename($core.String value) => $_setString(6, value);
  @$pb.TagNumber(270657590)
  $core.bool hasStatename() => $_has(6);
  @$pb.TagNumber(270657590)
  void clearStatename() => $_clearField(270657590);

  @$pb.TagNumber(277169476)
  InspectionLevel get inspectionlevel => $_getN(7);
  @$pb.TagNumber(277169476)
  set inspectionlevel(InspectionLevel value) => $_setField(277169476, value);
  @$pb.TagNumber(277169476)
  $core.bool hasInspectionlevel() => $_has(7);
  @$pb.TagNumber(277169476)
  void clearInspectionlevel() => $_clearField(277169476);

  @$pb.TagNumber(351839742)
  $core.bool get revealsecrets => $_getBF(8);
  @$pb.TagNumber(351839742)
  set revealsecrets($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(351839742)
  $core.bool hasRevealsecrets() => $_has(8);
  @$pb.TagNumber(351839742)
  void clearRevealsecrets() => $_clearField(351839742);

  @$pb.TagNumber(433614716)
  $core.String get input => $_getSZ(9);
  @$pb.TagNumber(433614716)
  set input($core.String value) => $_setString(9, value);
  @$pb.TagNumber(433614716)
  $core.bool hasInput() => $_has(9);
  @$pb.TagNumber(433614716)
  void clearInput() => $_clearField(433614716);
}

class TestStateOutput extends $pb.GeneratedMessage {
  factory TestStateOutput({
    $core.String? error,
    InspectionData? inspectiondata,
    $core.String? cause,
    $core.String? output,
    TestExecutionStatus? status,
    $core.String? nextstate,
  }) {
    final result = create();
    if (error != null) result.error = error;
    if (inspectiondata != null) result.inspectiondata = inspectiondata;
    if (cause != null) result.cause = cause;
    if (output != null) result.output = output;
    if (status != null) result.status = status;
    if (nextstate != null) result.nextstate = nextstate;
    return result;
  }

  TestStateOutput._();

  factory TestStateOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestStateOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestStateOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(26314578, _omitFieldNames ? '' : 'error')
    ..aOM<InspectionData>(113762044, _omitFieldNames ? '' : 'inspectiondata',
        subBuilder: InspectionData.create)
    ..aOS(145674785, _omitFieldNames ? '' : 'cause')
    ..aOS(430526213, _omitFieldNames ? '' : 'output')
    ..aE<TestExecutionStatus>(441153520, _omitFieldNames ? '' : 'status',
        enumValues: TestExecutionStatus.values)
    ..aOS(525594702, _omitFieldNames ? '' : 'nextstate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestStateOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestStateOutput copyWith(void Function(TestStateOutput) updates) =>
      super.copyWith((message) => updates(message as TestStateOutput))
          as TestStateOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestStateOutput create() => TestStateOutput._();
  @$core.override
  TestStateOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestStateOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestStateOutput>(create);
  static TestStateOutput? _defaultInstance;

  @$pb.TagNumber(26314578)
  $core.String get error => $_getSZ(0);
  @$pb.TagNumber(26314578)
  set error($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26314578)
  $core.bool hasError() => $_has(0);
  @$pb.TagNumber(26314578)
  void clearError() => $_clearField(26314578);

  @$pb.TagNumber(113762044)
  InspectionData get inspectiondata => $_getN(1);
  @$pb.TagNumber(113762044)
  set inspectiondata(InspectionData value) => $_setField(113762044, value);
  @$pb.TagNumber(113762044)
  $core.bool hasInspectiondata() => $_has(1);
  @$pb.TagNumber(113762044)
  void clearInspectiondata() => $_clearField(113762044);
  @$pb.TagNumber(113762044)
  InspectionData ensureInspectiondata() => $_ensure(1);

  @$pb.TagNumber(145674785)
  $core.String get cause => $_getSZ(2);
  @$pb.TagNumber(145674785)
  set cause($core.String value) => $_setString(2, value);
  @$pb.TagNumber(145674785)
  $core.bool hasCause() => $_has(2);
  @$pb.TagNumber(145674785)
  void clearCause() => $_clearField(145674785);

  @$pb.TagNumber(430526213)
  $core.String get output => $_getSZ(3);
  @$pb.TagNumber(430526213)
  set output($core.String value) => $_setString(3, value);
  @$pb.TagNumber(430526213)
  $core.bool hasOutput() => $_has(3);
  @$pb.TagNumber(430526213)
  void clearOutput() => $_clearField(430526213);

  @$pb.TagNumber(441153520)
  TestExecutionStatus get status => $_getN(4);
  @$pb.TagNumber(441153520)
  set status(TestExecutionStatus value) => $_setField(441153520, value);
  @$pb.TagNumber(441153520)
  $core.bool hasStatus() => $_has(4);
  @$pb.TagNumber(441153520)
  void clearStatus() => $_clearField(441153520);

  @$pb.TagNumber(525594702)
  $core.String get nextstate => $_getSZ(5);
  @$pb.TagNumber(525594702)
  set nextstate($core.String value) => $_setString(5, value);
  @$pb.TagNumber(525594702)
  $core.bool hasNextstate() => $_has(5);
  @$pb.TagNumber(525594702)
  void clearNextstate() => $_clearField(525594702);
}

class TooManyTags extends $pb.GeneratedMessage {
  factory TooManyTags({
    $core.String? resourcename,
    $core.String? message,
  }) {
    final result = create();
    if (resourcename != null) result.resourcename = resourcename;
    if (message != null) result.message = message;
    return result;
  }

  TooManyTags._();

  factory TooManyTags.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyTags.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyTags',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(17776375, _omitFieldNames ? '' : 'resourcename')
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTags clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTags copyWith(void Function(TooManyTags) updates) =>
      super.copyWith((message) => updates(message as TooManyTags))
          as TooManyTags;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyTags create() => TooManyTags._();
  @$core.override
  TooManyTags createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyTags getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyTags>(create);
  static TooManyTags? _defaultInstance;

  @$pb.TagNumber(17776375)
  $core.String get resourcename => $_getSZ(0);
  @$pb.TagNumber(17776375)
  set resourcename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(17776375)
  $core.bool hasResourcename() => $_has(0);
  @$pb.TagNumber(17776375)
  void clearResourcename() => $_clearField(17776375);

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TracingConfiguration extends $pb.GeneratedMessage {
  factory TracingConfiguration({
    $core.bool? enabled,
  }) {
    final result = create();
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  TracingConfiguration._();

  factory TracingConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TracingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TracingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOB(49525663, _omitFieldNames ? '' : 'enabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TracingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TracingConfiguration copyWith(void Function(TracingConfiguration) updates) =>
      super.copyWith((message) => updates(message as TracingConfiguration))
          as TracingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TracingConfiguration create() => TracingConfiguration._();
  @$core.override
  TracingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TracingConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TracingConfiguration>(create);
  static TracingConfiguration? _defaultInstance;

  @$pb.TagNumber(49525663)
  $core.bool get enabled => $_getBF(0);
  @$pb.TagNumber(49525663)
  set enabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(49525663)
  $core.bool hasEnabled() => $_has(0);
  @$pb.TagNumber(49525663)
  void clearEnabled() => $_clearField(49525663);
}

class UntagResourceInput extends $pb.GeneratedMessage {
  factory UntagResourceInput({
    $core.String? resourcearn,
    $core.Iterable<$core.String>? tagkeys,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
    return result;
  }

  UntagResourceInput._();

  factory UntagResourceInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UntagResourceInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UntagResourceInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(67806797, _omitFieldNames ? '' : 'resourcearn')
    ..pPS(78811036, _omitFieldNames ? '' : 'tagkeys')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceInput copyWith(void Function(UntagResourceInput) updates) =>
      super.copyWith((message) => updates(message as UntagResourceInput))
          as UntagResourceInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UntagResourceInput create() => UntagResourceInput._();
  @$core.override
  UntagResourceInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UntagResourceInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UntagResourceInput>(create);
  static UntagResourceInput? _defaultInstance;

  @$pb.TagNumber(67806797)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(67806797)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(67806797)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(67806797)
  void clearResourcearn() => $_clearField(67806797);

  @$pb.TagNumber(78811036)
  $pb.PbList<$core.String> get tagkeys => $_getList(1);
}

class UntagResourceOutput extends $pb.GeneratedMessage {
  factory UntagResourceOutput() => create();

  UntagResourceOutput._();

  factory UntagResourceOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UntagResourceOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UntagResourceOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceOutput copyWith(void Function(UntagResourceOutput) updates) =>
      super.copyWith((message) => updates(message as UntagResourceOutput))
          as UntagResourceOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UntagResourceOutput create() => UntagResourceOutput._();
  @$core.override
  UntagResourceOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UntagResourceOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UntagResourceOutput>(create);
  static UntagResourceOutput? _defaultInstance;
}

class UpdateMapRunInput extends $pb.GeneratedMessage {
  factory UpdateMapRunInput({
    $core.String? maprunarn,
    $fixnum.Int64? toleratedfailurecount,
    $core.int? maxconcurrency,
    $core.double? toleratedfailurepercentage,
  }) {
    final result = create();
    if (maprunarn != null) result.maprunarn = maprunarn;
    if (toleratedfailurecount != null)
      result.toleratedfailurecount = toleratedfailurecount;
    if (maxconcurrency != null) result.maxconcurrency = maxconcurrency;
    if (toleratedfailurepercentage != null)
      result.toleratedfailurepercentage = toleratedfailurepercentage;
    return result;
  }

  UpdateMapRunInput._();

  factory UpdateMapRunInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateMapRunInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateMapRunInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(18199994, _omitFieldNames ? '' : 'maprunarn')
    ..aInt64(41834811, _omitFieldNames ? '' : 'toleratedfailurecount')
    ..aI(100901405, _omitFieldNames ? '' : 'maxconcurrency')
    ..aD(116496164, _omitFieldNames ? '' : 'toleratedfailurepercentage',
        fieldType: $pb.PbFieldType.OF)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMapRunInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMapRunInput copyWith(void Function(UpdateMapRunInput) updates) =>
      super.copyWith((message) => updates(message as UpdateMapRunInput))
          as UpdateMapRunInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateMapRunInput create() => UpdateMapRunInput._();
  @$core.override
  UpdateMapRunInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateMapRunInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateMapRunInput>(create);
  static UpdateMapRunInput? _defaultInstance;

  @$pb.TagNumber(18199994)
  $core.String get maprunarn => $_getSZ(0);
  @$pb.TagNumber(18199994)
  set maprunarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(18199994)
  $core.bool hasMaprunarn() => $_has(0);
  @$pb.TagNumber(18199994)
  void clearMaprunarn() => $_clearField(18199994);

  @$pb.TagNumber(41834811)
  $fixnum.Int64 get toleratedfailurecount => $_getI64(1);
  @$pb.TagNumber(41834811)
  set toleratedfailurecount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(41834811)
  $core.bool hasToleratedfailurecount() => $_has(1);
  @$pb.TagNumber(41834811)
  void clearToleratedfailurecount() => $_clearField(41834811);

  @$pb.TagNumber(100901405)
  $core.int get maxconcurrency => $_getIZ(2);
  @$pb.TagNumber(100901405)
  set maxconcurrency($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(100901405)
  $core.bool hasMaxconcurrency() => $_has(2);
  @$pb.TagNumber(100901405)
  void clearMaxconcurrency() => $_clearField(100901405);

  @$pb.TagNumber(116496164)
  $core.double get toleratedfailurepercentage => $_getN(3);
  @$pb.TagNumber(116496164)
  set toleratedfailurepercentage($core.double value) => $_setFloat(3, value);
  @$pb.TagNumber(116496164)
  $core.bool hasToleratedfailurepercentage() => $_has(3);
  @$pb.TagNumber(116496164)
  void clearToleratedfailurepercentage() => $_clearField(116496164);
}

class UpdateMapRunOutput extends $pb.GeneratedMessage {
  factory UpdateMapRunOutput() => create();

  UpdateMapRunOutput._();

  factory UpdateMapRunOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateMapRunOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateMapRunOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMapRunOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMapRunOutput copyWith(void Function(UpdateMapRunOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateMapRunOutput))
          as UpdateMapRunOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateMapRunOutput create() => UpdateMapRunOutput._();
  @$core.override
  UpdateMapRunOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateMapRunOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateMapRunOutput>(create);
  static UpdateMapRunOutput? _defaultInstance;
}

class UpdateStateMachineAliasInput extends $pb.GeneratedMessage {
  factory UpdateStateMachineAliasInput({
    $core.String? description,
    $core.Iterable<RoutingConfigurationListItem>? routingconfiguration,
    $core.String? statemachinealiasarn,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (routingconfiguration != null)
      result.routingconfiguration.addAll(routingconfiguration);
    if (statemachinealiasarn != null)
      result.statemachinealiasarn = statemachinealiasarn;
    return result;
  }

  UpdateStateMachineAliasInput._();

  factory UpdateStateMachineAliasInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateStateMachineAliasInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateStateMachineAliasInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..pPM<RoutingConfigurationListItem>(
        372891510, _omitFieldNames ? '' : 'routingconfiguration',
        subBuilder: RoutingConfigurationListItem.create)
    ..aOS(530344465, _omitFieldNames ? '' : 'statemachinealiasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStateMachineAliasInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStateMachineAliasInput copyWith(
          void Function(UpdateStateMachineAliasInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateStateMachineAliasInput))
          as UpdateStateMachineAliasInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateStateMachineAliasInput create() =>
      UpdateStateMachineAliasInput._();
  @$core.override
  UpdateStateMachineAliasInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateStateMachineAliasInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateStateMachineAliasInput>(create);
  static UpdateStateMachineAliasInput? _defaultInstance;

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(372891510)
  $pb.PbList<RoutingConfigurationListItem> get routingconfiguration =>
      $_getList(1);

  @$pb.TagNumber(530344465)
  $core.String get statemachinealiasarn => $_getSZ(2);
  @$pb.TagNumber(530344465)
  set statemachinealiasarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(530344465)
  $core.bool hasStatemachinealiasarn() => $_has(2);
  @$pb.TagNumber(530344465)
  void clearStatemachinealiasarn() => $_clearField(530344465);
}

class UpdateStateMachineAliasOutput extends $pb.GeneratedMessage {
  factory UpdateStateMachineAliasOutput({
    $core.String? updatedate,
  }) {
    final result = create();
    if (updatedate != null) result.updatedate = updatedate;
    return result;
  }

  UpdateStateMachineAliasOutput._();

  factory UpdateStateMachineAliasOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateStateMachineAliasOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateStateMachineAliasOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(510552561, _omitFieldNames ? '' : 'updatedate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStateMachineAliasOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStateMachineAliasOutput copyWith(
          void Function(UpdateStateMachineAliasOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateStateMachineAliasOutput))
          as UpdateStateMachineAliasOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateStateMachineAliasOutput create() =>
      UpdateStateMachineAliasOutput._();
  @$core.override
  UpdateStateMachineAliasOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateStateMachineAliasOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateStateMachineAliasOutput>(create);
  static UpdateStateMachineAliasOutput? _defaultInstance;

  @$pb.TagNumber(510552561)
  $core.String get updatedate => $_getSZ(0);
  @$pb.TagNumber(510552561)
  set updatedate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(510552561)
  $core.bool hasUpdatedate() => $_has(0);
  @$pb.TagNumber(510552561)
  void clearUpdatedate() => $_clearField(510552561);
}

class UpdateStateMachineInput extends $pb.GeneratedMessage {
  factory UpdateStateMachineInput({
    $core.String? definition,
    EncryptionConfiguration? encryptionconfiguration,
    $core.String? rolearn,
    $core.bool? publish,
    $core.String? statemachinearn,
    LoggingConfiguration? loggingconfiguration,
    $core.String? versiondescription,
    TracingConfiguration? tracingconfiguration,
  }) {
    final result = create();
    if (definition != null) result.definition = definition;
    if (encryptionconfiguration != null)
      result.encryptionconfiguration = encryptionconfiguration;
    if (rolearn != null) result.rolearn = rolearn;
    if (publish != null) result.publish = publish;
    if (statemachinearn != null) result.statemachinearn = statemachinearn;
    if (loggingconfiguration != null)
      result.loggingconfiguration = loggingconfiguration;
    if (versiondescription != null)
      result.versiondescription = versiondescription;
    if (tracingconfiguration != null)
      result.tracingconfiguration = tracingconfiguration;
    return result;
  }

  UpdateStateMachineInput._();

  factory UpdateStateMachineInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateStateMachineInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateStateMachineInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(68443297, _omitFieldNames ? '' : 'definition')
    ..aOM<EncryptionConfiguration>(
        167857431, _omitFieldNames ? '' : 'encryptionconfiguration',
        subBuilder: EncryptionConfiguration.create)
    ..aOS(170019745, _omitFieldNames ? '' : 'rolearn')
    ..aOB(264247305, _omitFieldNames ? '' : 'publish')
    ..aOS(393321971, _omitFieldNames ? '' : 'statemachinearn')
    ..aOM<LoggingConfiguration>(
        420811605, _omitFieldNames ? '' : 'loggingconfiguration',
        subBuilder: LoggingConfiguration.create)
    ..aOS(434714300, _omitFieldNames ? '' : 'versiondescription')
    ..aOM<TracingConfiguration>(
        491315910, _omitFieldNames ? '' : 'tracingconfiguration',
        subBuilder: TracingConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStateMachineInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStateMachineInput copyWith(
          void Function(UpdateStateMachineInput) updates) =>
      super.copyWith((message) => updates(message as UpdateStateMachineInput))
          as UpdateStateMachineInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateStateMachineInput create() => UpdateStateMachineInput._();
  @$core.override
  UpdateStateMachineInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateStateMachineInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateStateMachineInput>(create);
  static UpdateStateMachineInput? _defaultInstance;

  @$pb.TagNumber(68443297)
  $core.String get definition => $_getSZ(0);
  @$pb.TagNumber(68443297)
  set definition($core.String value) => $_setString(0, value);
  @$pb.TagNumber(68443297)
  $core.bool hasDefinition() => $_has(0);
  @$pb.TagNumber(68443297)
  void clearDefinition() => $_clearField(68443297);

  @$pb.TagNumber(167857431)
  EncryptionConfiguration get encryptionconfiguration => $_getN(1);
  @$pb.TagNumber(167857431)
  set encryptionconfiguration(EncryptionConfiguration value) =>
      $_setField(167857431, value);
  @$pb.TagNumber(167857431)
  $core.bool hasEncryptionconfiguration() => $_has(1);
  @$pb.TagNumber(167857431)
  void clearEncryptionconfiguration() => $_clearField(167857431);
  @$pb.TagNumber(167857431)
  EncryptionConfiguration ensureEncryptionconfiguration() => $_ensure(1);

  @$pb.TagNumber(170019745)
  $core.String get rolearn => $_getSZ(2);
  @$pb.TagNumber(170019745)
  set rolearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(170019745)
  $core.bool hasRolearn() => $_has(2);
  @$pb.TagNumber(170019745)
  void clearRolearn() => $_clearField(170019745);

  @$pb.TagNumber(264247305)
  $core.bool get publish => $_getBF(3);
  @$pb.TagNumber(264247305)
  set publish($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(264247305)
  $core.bool hasPublish() => $_has(3);
  @$pb.TagNumber(264247305)
  void clearPublish() => $_clearField(264247305);

  @$pb.TagNumber(393321971)
  $core.String get statemachinearn => $_getSZ(4);
  @$pb.TagNumber(393321971)
  set statemachinearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(393321971)
  $core.bool hasStatemachinearn() => $_has(4);
  @$pb.TagNumber(393321971)
  void clearStatemachinearn() => $_clearField(393321971);

  @$pb.TagNumber(420811605)
  LoggingConfiguration get loggingconfiguration => $_getN(5);
  @$pb.TagNumber(420811605)
  set loggingconfiguration(LoggingConfiguration value) =>
      $_setField(420811605, value);
  @$pb.TagNumber(420811605)
  $core.bool hasLoggingconfiguration() => $_has(5);
  @$pb.TagNumber(420811605)
  void clearLoggingconfiguration() => $_clearField(420811605);
  @$pb.TagNumber(420811605)
  LoggingConfiguration ensureLoggingconfiguration() => $_ensure(5);

  @$pb.TagNumber(434714300)
  $core.String get versiondescription => $_getSZ(6);
  @$pb.TagNumber(434714300)
  set versiondescription($core.String value) => $_setString(6, value);
  @$pb.TagNumber(434714300)
  $core.bool hasVersiondescription() => $_has(6);
  @$pb.TagNumber(434714300)
  void clearVersiondescription() => $_clearField(434714300);

  @$pb.TagNumber(491315910)
  TracingConfiguration get tracingconfiguration => $_getN(7);
  @$pb.TagNumber(491315910)
  set tracingconfiguration(TracingConfiguration value) =>
      $_setField(491315910, value);
  @$pb.TagNumber(491315910)
  $core.bool hasTracingconfiguration() => $_has(7);
  @$pb.TagNumber(491315910)
  void clearTracingconfiguration() => $_clearField(491315910);
  @$pb.TagNumber(491315910)
  TracingConfiguration ensureTracingconfiguration() => $_ensure(7);
}

class UpdateStateMachineOutput extends $pb.GeneratedMessage {
  factory UpdateStateMachineOutput({
    $core.String? statemachineversionarn,
    $core.String? revisionid,
    $core.String? updatedate,
  }) {
    final result = create();
    if (statemachineversionarn != null)
      result.statemachineversionarn = statemachineversionarn;
    if (revisionid != null) result.revisionid = revisionid;
    if (updatedate != null) result.updatedate = updatedate;
    return result;
  }

  UpdateStateMachineOutput._();

  factory UpdateStateMachineOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateStateMachineOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateStateMachineOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(69976825, _omitFieldNames ? '' : 'statemachineversionarn')
    ..aOS(369170086, _omitFieldNames ? '' : 'revisionid')
    ..aOS(510552561, _omitFieldNames ? '' : 'updatedate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStateMachineOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStateMachineOutput copyWith(
          void Function(UpdateStateMachineOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateStateMachineOutput))
          as UpdateStateMachineOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateStateMachineOutput create() => UpdateStateMachineOutput._();
  @$core.override
  UpdateStateMachineOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateStateMachineOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateStateMachineOutput>(create);
  static UpdateStateMachineOutput? _defaultInstance;

  @$pb.TagNumber(69976825)
  $core.String get statemachineversionarn => $_getSZ(0);
  @$pb.TagNumber(69976825)
  set statemachineversionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(69976825)
  $core.bool hasStatemachineversionarn() => $_has(0);
  @$pb.TagNumber(69976825)
  void clearStatemachineversionarn() => $_clearField(69976825);

  @$pb.TagNumber(369170086)
  $core.String get revisionid => $_getSZ(1);
  @$pb.TagNumber(369170086)
  set revisionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(369170086)
  $core.bool hasRevisionid() => $_has(1);
  @$pb.TagNumber(369170086)
  void clearRevisionid() => $_clearField(369170086);

  @$pb.TagNumber(510552561)
  $core.String get updatedate => $_getSZ(2);
  @$pb.TagNumber(510552561)
  set updatedate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(510552561)
  $core.bool hasUpdatedate() => $_has(2);
  @$pb.TagNumber(510552561)
  void clearUpdatedate() => $_clearField(510552561);
}

class ValidateStateMachineDefinitionDiagnostic extends $pb.GeneratedMessage {
  factory ValidateStateMachineDefinitionDiagnostic({
    $core.String? message,
    $core.String? location,
    ValidateStateMachineDefinitionSeverity? severity,
    $core.String? code,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (location != null) result.location = location;
    if (severity != null) result.severity = severity;
    if (code != null) result.code = code;
    return result;
  }

  ValidateStateMachineDefinitionDiagnostic._();

  factory ValidateStateMachineDefinitionDiagnostic.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ValidateStateMachineDefinitionDiagnostic.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ValidateStateMachineDefinitionDiagnostic',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..aOS(200649127, _omitFieldNames ? '' : 'location')
    ..aE<ValidateStateMachineDefinitionSeverity>(
        268193715, _omitFieldNames ? '' : 'severity',
        enumValues: ValidateStateMachineDefinitionSeverity.values)
    ..aOS(422669557, _omitFieldNames ? '' : 'code')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateStateMachineDefinitionDiagnostic clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateStateMachineDefinitionDiagnostic copyWith(
          void Function(ValidateStateMachineDefinitionDiagnostic) updates) =>
      super.copyWith((message) =>
              updates(message as ValidateStateMachineDefinitionDiagnostic))
          as ValidateStateMachineDefinitionDiagnostic;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidateStateMachineDefinitionDiagnostic create() =>
      ValidateStateMachineDefinitionDiagnostic._();
  @$core.override
  ValidateStateMachineDefinitionDiagnostic createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ValidateStateMachineDefinitionDiagnostic getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ValidateStateMachineDefinitionDiagnostic>(create);
  static ValidateStateMachineDefinitionDiagnostic? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(200649127)
  $core.String get location => $_getSZ(1);
  @$pb.TagNumber(200649127)
  set location($core.String value) => $_setString(1, value);
  @$pb.TagNumber(200649127)
  $core.bool hasLocation() => $_has(1);
  @$pb.TagNumber(200649127)
  void clearLocation() => $_clearField(200649127);

  @$pb.TagNumber(268193715)
  ValidateStateMachineDefinitionSeverity get severity => $_getN(2);
  @$pb.TagNumber(268193715)
  set severity(ValidateStateMachineDefinitionSeverity value) =>
      $_setField(268193715, value);
  @$pb.TagNumber(268193715)
  $core.bool hasSeverity() => $_has(2);
  @$pb.TagNumber(268193715)
  void clearSeverity() => $_clearField(268193715);

  @$pb.TagNumber(422669557)
  $core.String get code => $_getSZ(3);
  @$pb.TagNumber(422669557)
  set code($core.String value) => $_setString(3, value);
  @$pb.TagNumber(422669557)
  $core.bool hasCode() => $_has(3);
  @$pb.TagNumber(422669557)
  void clearCode() => $_clearField(422669557);
}

class ValidateStateMachineDefinitionInput extends $pb.GeneratedMessage {
  factory ValidateStateMachineDefinitionInput({
    $core.String? definition,
    ValidateStateMachineDefinitionSeverity? severity,
    StateMachineType? type,
    $core.int? maxresults,
  }) {
    final result = create();
    if (definition != null) result.definition = definition;
    if (severity != null) result.severity = severity;
    if (type != null) result.type = type;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ValidateStateMachineDefinitionInput._();

  factory ValidateStateMachineDefinitionInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ValidateStateMachineDefinitionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ValidateStateMachineDefinitionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(68443297, _omitFieldNames ? '' : 'definition')
    ..aE<ValidateStateMachineDefinitionSeverity>(
        268193715, _omitFieldNames ? '' : 'severity',
        enumValues: ValidateStateMachineDefinitionSeverity.values)
    ..aE<StateMachineType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: StateMachineType.values)
    ..aI(465170002, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateStateMachineDefinitionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateStateMachineDefinitionInput copyWith(
          void Function(ValidateStateMachineDefinitionInput) updates) =>
      super.copyWith((message) =>
              updates(message as ValidateStateMachineDefinitionInput))
          as ValidateStateMachineDefinitionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidateStateMachineDefinitionInput create() =>
      ValidateStateMachineDefinitionInput._();
  @$core.override
  ValidateStateMachineDefinitionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ValidateStateMachineDefinitionInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ValidateStateMachineDefinitionInput>(create);
  static ValidateStateMachineDefinitionInput? _defaultInstance;

  @$pb.TagNumber(68443297)
  $core.String get definition => $_getSZ(0);
  @$pb.TagNumber(68443297)
  set definition($core.String value) => $_setString(0, value);
  @$pb.TagNumber(68443297)
  $core.bool hasDefinition() => $_has(0);
  @$pb.TagNumber(68443297)
  void clearDefinition() => $_clearField(68443297);

  @$pb.TagNumber(268193715)
  ValidateStateMachineDefinitionSeverity get severity => $_getN(1);
  @$pb.TagNumber(268193715)
  set severity(ValidateStateMachineDefinitionSeverity value) =>
      $_setField(268193715, value);
  @$pb.TagNumber(268193715)
  $core.bool hasSeverity() => $_has(1);
  @$pb.TagNumber(268193715)
  void clearSeverity() => $_clearField(268193715);

  @$pb.TagNumber(287830350)
  StateMachineType get type => $_getN(2);
  @$pb.TagNumber(287830350)
  set type(StateMachineType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(465170002)
  $core.int get maxresults => $_getIZ(3);
  @$pb.TagNumber(465170002)
  set maxresults($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(465170002)
  $core.bool hasMaxresults() => $_has(3);
  @$pb.TagNumber(465170002)
  void clearMaxresults() => $_clearField(465170002);
}

class ValidateStateMachineDefinitionOutput extends $pb.GeneratedMessage {
  factory ValidateStateMachineDefinitionOutput({
    $core.Iterable<ValidateStateMachineDefinitionDiagnostic>? diagnostics,
    ValidateStateMachineDefinitionResultCode? result,
    $core.bool? truncated,
  }) {
    final result$ = create();
    if (diagnostics != null) result$.diagnostics.addAll(diagnostics);
    if (result != null) result$.result = result;
    if (truncated != null) result$.truncated = truncated;
    return result$;
  }

  ValidateStateMachineDefinitionOutput._();

  factory ValidateStateMachineDefinitionOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ValidateStateMachineDefinitionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ValidateStateMachineDefinitionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..pPM<ValidateStateMachineDefinitionDiagnostic>(
        159421432, _omitFieldNames ? '' : 'diagnostics',
        subBuilder: ValidateStateMachineDefinitionDiagnostic.create)
    ..aE<ValidateStateMachineDefinitionResultCode>(
        171406885, _omitFieldNames ? '' : 'result',
        enumValues: ValidateStateMachineDefinitionResultCode.values)
    ..aOB(407490858, _omitFieldNames ? '' : 'truncated')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateStateMachineDefinitionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateStateMachineDefinitionOutput copyWith(
          void Function(ValidateStateMachineDefinitionOutput) updates) =>
      super.copyWith((message) =>
              updates(message as ValidateStateMachineDefinitionOutput))
          as ValidateStateMachineDefinitionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidateStateMachineDefinitionOutput create() =>
      ValidateStateMachineDefinitionOutput._();
  @$core.override
  ValidateStateMachineDefinitionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ValidateStateMachineDefinitionOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ValidateStateMachineDefinitionOutput>(create);
  static ValidateStateMachineDefinitionOutput? _defaultInstance;

  @$pb.TagNumber(159421432)
  $pb.PbList<ValidateStateMachineDefinitionDiagnostic> get diagnostics =>
      $_getList(0);

  @$pb.TagNumber(171406885)
  ValidateStateMachineDefinitionResultCode get result => $_getN(1);
  @$pb.TagNumber(171406885)
  set result(ValidateStateMachineDefinitionResultCode value) =>
      $_setField(171406885, value);
  @$pb.TagNumber(171406885)
  $core.bool hasResult() => $_has(1);
  @$pb.TagNumber(171406885)
  void clearResult() => $_clearField(171406885);

  @$pb.TagNumber(407490858)
  $core.bool get truncated => $_getBF(2);
  @$pb.TagNumber(407490858)
  set truncated($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(407490858)
  $core.bool hasTruncated() => $_has(2);
  @$pb.TagNumber(407490858)
  void clearTruncated() => $_clearField(407490858);
}

class ValidationException extends $pb.GeneratedMessage {
  factory ValidationException({
    $core.String? message,
    ValidationExceptionReason? reason,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (reason != null) result.reason = reason;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'states'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..aE<ValidationExceptionReason>(413359642, _omitFieldNames ? '' : 'reason',
        enumValues: ValidationExceptionReason.values)
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(413359642)
  ValidationExceptionReason get reason => $_getN(1);
  @$pb.TagNumber(413359642)
  set reason(ValidationExceptionReason value) => $_setField(413359642, value);
  @$pb.TagNumber(413359642)
  $core.bool hasReason() => $_has(1);
  @$pb.TagNumber(413359642)
  void clearReason() => $_clearField(413359642);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
