// This is a generated file - do not edit.
//
// Generated from sqs.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'sqs.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'sqs.pbenum.dart';

class AddPermissionRequest extends $pb.GeneratedMessage {
  factory AddPermissionRequest({
    $core.Iterable<$core.String>? actions,
    $core.Iterable<$core.String>? awsaccountids,
    $core.String? queueurl,
    $core.String? label,
  }) {
    final result = create();
    if (actions != null) result.actions.addAll(actions);
    if (awsaccountids != null) result.awsaccountids.addAll(awsaccountids);
    if (queueurl != null) result.queueurl = queueurl;
    if (label != null) result.label = label;
    return result;
  }

  AddPermissionRequest._();

  factory AddPermissionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AddPermissionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddPermissionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPS(106935557, _omitFieldNames ? '' : 'actions')
    ..pPS(417597390, _omitFieldNames ? '' : 'awsaccountids')
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..aOS(516747934, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddPermissionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddPermissionRequest copyWith(void Function(AddPermissionRequest) updates) =>
      super.copyWith((message) => updates(message as AddPermissionRequest))
          as AddPermissionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddPermissionRequest create() => AddPermissionRequest._();
  @$core.override
  AddPermissionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AddPermissionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddPermissionRequest>(create);
  static AddPermissionRequest? _defaultInstance;

  @$pb.TagNumber(106935557)
  $pb.PbList<$core.String> get actions => $_getList(0);

  @$pb.TagNumber(417597390)
  $pb.PbList<$core.String> get awsaccountids => $_getList(1);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(2);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(2, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(2);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);

  @$pb.TagNumber(516747934)
  $core.String get label => $_getSZ(3);
  @$pb.TagNumber(516747934)
  set label($core.String value) => $_setString(3, value);
  @$pb.TagNumber(516747934)
  $core.bool hasLabel() => $_has(3);
  @$pb.TagNumber(516747934)
  void clearLabel() => $_clearField(516747934);
}

class BatchEntryIdsNotDistinct extends $pb.GeneratedMessage {
  factory BatchEntryIdsNotDistinct({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  BatchEntryIdsNotDistinct._();

  factory BatchEntryIdsNotDistinct.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchEntryIdsNotDistinct.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchEntryIdsNotDistinct',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchEntryIdsNotDistinct clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchEntryIdsNotDistinct copyWith(
          void Function(BatchEntryIdsNotDistinct) updates) =>
      super.copyWith((message) => updates(message as BatchEntryIdsNotDistinct))
          as BatchEntryIdsNotDistinct;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchEntryIdsNotDistinct create() => BatchEntryIdsNotDistinct._();
  @$core.override
  BatchEntryIdsNotDistinct createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchEntryIdsNotDistinct getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchEntryIdsNotDistinct>(create);
  static BatchEntryIdsNotDistinct? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class BatchRequestTooLong extends $pb.GeneratedMessage {
  factory BatchRequestTooLong({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  BatchRequestTooLong._();

  factory BatchRequestTooLong.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchRequestTooLong.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchRequestTooLong',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchRequestTooLong clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchRequestTooLong copyWith(void Function(BatchRequestTooLong) updates) =>
      super.copyWith((message) => updates(message as BatchRequestTooLong))
          as BatchRequestTooLong;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchRequestTooLong create() => BatchRequestTooLong._();
  @$core.override
  BatchRequestTooLong createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchRequestTooLong getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchRequestTooLong>(create);
  static BatchRequestTooLong? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class BatchResultErrorEntry extends $pb.GeneratedMessage {
  factory BatchResultErrorEntry({
    $core.bool? senderfault,
    $core.String? message,
    $core.String? id,
    $core.String? code,
  }) {
    final result = create();
    if (senderfault != null) result.senderfault = senderfault;
    if (message != null) result.message = message;
    if (id != null) result.id = id;
    if (code != null) result.code = code;
    return result;
  }

  BatchResultErrorEntry._();

  factory BatchResultErrorEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchResultErrorEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchResultErrorEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOB(28412929, _omitFieldNames ? '' : 'senderfault')
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(425572629, _omitFieldNames ? '' : 'code')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchResultErrorEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchResultErrorEntry copyWith(
          void Function(BatchResultErrorEntry) updates) =>
      super.copyWith((message) => updates(message as BatchResultErrorEntry))
          as BatchResultErrorEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchResultErrorEntry create() => BatchResultErrorEntry._();
  @$core.override
  BatchResultErrorEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchResultErrorEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchResultErrorEntry>(create);
  static BatchResultErrorEntry? _defaultInstance;

  @$pb.TagNumber(28412929)
  $core.bool get senderfault => $_getBF(0);
  @$pb.TagNumber(28412929)
  set senderfault($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(28412929)
  $core.bool hasSenderfault() => $_has(0);
  @$pb.TagNumber(28412929)
  void clearSenderfault() => $_clearField(28412929);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(425572629)
  $core.String get code => $_getSZ(3);
  @$pb.TagNumber(425572629)
  set code($core.String value) => $_setString(3, value);
  @$pb.TagNumber(425572629)
  $core.bool hasCode() => $_has(3);
  @$pb.TagNumber(425572629)
  void clearCode() => $_clearField(425572629);
}

class CancelMessageMoveTaskRequest extends $pb.GeneratedMessage {
  factory CancelMessageMoveTaskRequest({
    $core.String? taskhandle,
  }) {
    final result = create();
    if (taskhandle != null) result.taskhandle = taskhandle;
    return result;
  }

  CancelMessageMoveTaskRequest._();

  factory CancelMessageMoveTaskRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelMessageMoveTaskRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelMessageMoveTaskRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(190544291, _omitFieldNames ? '' : 'taskhandle')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelMessageMoveTaskRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelMessageMoveTaskRequest copyWith(
          void Function(CancelMessageMoveTaskRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CancelMessageMoveTaskRequest))
          as CancelMessageMoveTaskRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelMessageMoveTaskRequest create() =>
      CancelMessageMoveTaskRequest._();
  @$core.override
  CancelMessageMoveTaskRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelMessageMoveTaskRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelMessageMoveTaskRequest>(create);
  static CancelMessageMoveTaskRequest? _defaultInstance;

  @$pb.TagNumber(190544291)
  $core.String get taskhandle => $_getSZ(0);
  @$pb.TagNumber(190544291)
  set taskhandle($core.String value) => $_setString(0, value);
  @$pb.TagNumber(190544291)
  $core.bool hasTaskhandle() => $_has(0);
  @$pb.TagNumber(190544291)
  void clearTaskhandle() => $_clearField(190544291);
}

class CancelMessageMoveTaskResult extends $pb.GeneratedMessage {
  factory CancelMessageMoveTaskResult({
    $fixnum.Int64? approximatenumberofmessagesmoved,
  }) {
    final result = create();
    if (approximatenumberofmessagesmoved != null)
      result.approximatenumberofmessagesmoved =
          approximatenumberofmessagesmoved;
    return result;
  }

  CancelMessageMoveTaskResult._();

  factory CancelMessageMoveTaskResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelMessageMoveTaskResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelMessageMoveTaskResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aInt64(
        510181657, _omitFieldNames ? '' : 'approximatenumberofmessagesmoved')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelMessageMoveTaskResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelMessageMoveTaskResult copyWith(
          void Function(CancelMessageMoveTaskResult) updates) =>
      super.copyWith(
              (message) => updates(message as CancelMessageMoveTaskResult))
          as CancelMessageMoveTaskResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelMessageMoveTaskResult create() =>
      CancelMessageMoveTaskResult._();
  @$core.override
  CancelMessageMoveTaskResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelMessageMoveTaskResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelMessageMoveTaskResult>(create);
  static CancelMessageMoveTaskResult? _defaultInstance;

  @$pb.TagNumber(510181657)
  $fixnum.Int64 get approximatenumberofmessagesmoved => $_getI64(0);
  @$pb.TagNumber(510181657)
  set approximatenumberofmessagesmoved($fixnum.Int64 value) =>
      $_setInt64(0, value);
  @$pb.TagNumber(510181657)
  $core.bool hasApproximatenumberofmessagesmoved() => $_has(0);
  @$pb.TagNumber(510181657)
  void clearApproximatenumberofmessagesmoved() => $_clearField(510181657);
}

class ChangeMessageVisibilityBatchRequest extends $pb.GeneratedMessage {
  factory ChangeMessageVisibilityBatchRequest({
    $core.Iterable<ChangeMessageVisibilityBatchRequestEntry>? entries,
    $core.String? queueurl,
  }) {
    final result = create();
    if (entries != null) result.entries.addAll(entries);
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  ChangeMessageVisibilityBatchRequest._();

  factory ChangeMessageVisibilityBatchRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeMessageVisibilityBatchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeMessageVisibilityBatchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPM<ChangeMessageVisibilityBatchRequestEntry>(
        481075860, _omitFieldNames ? '' : 'entries',
        subBuilder: ChangeMessageVisibilityBatchRequestEntry.create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityBatchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityBatchRequest copyWith(
          void Function(ChangeMessageVisibilityBatchRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ChangeMessageVisibilityBatchRequest))
          as ChangeMessageVisibilityBatchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityBatchRequest create() =>
      ChangeMessageVisibilityBatchRequest._();
  @$core.override
  ChangeMessageVisibilityBatchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityBatchRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ChangeMessageVisibilityBatchRequest>(create);
  static ChangeMessageVisibilityBatchRequest? _defaultInstance;

  @$pb.TagNumber(481075860)
  $pb.PbList<ChangeMessageVisibilityBatchRequestEntry> get entries =>
      $_getList(0);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(1);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(1);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class ChangeMessageVisibilityBatchRequestEntry extends $pb.GeneratedMessage {
  factory ChangeMessageVisibilityBatchRequestEntry({
    $core.String? receipthandle,
    $core.String? id,
    $core.int? visibilitytimeout,
  }) {
    final result = create();
    if (receipthandle != null) result.receipthandle = receipthandle;
    if (id != null) result.id = id;
    if (visibilitytimeout != null) result.visibilitytimeout = visibilitytimeout;
    return result;
  }

  ChangeMessageVisibilityBatchRequestEntry._();

  factory ChangeMessageVisibilityBatchRequestEntry.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeMessageVisibilityBatchRequestEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeMessageVisibilityBatchRequestEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(134471750, _omitFieldNames ? '' : 'receipthandle')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aI(460820073, _omitFieldNames ? '' : 'visibilitytimeout')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityBatchRequestEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityBatchRequestEntry copyWith(
          void Function(ChangeMessageVisibilityBatchRequestEntry) updates) =>
      super.copyWith((message) =>
              updates(message as ChangeMessageVisibilityBatchRequestEntry))
          as ChangeMessageVisibilityBatchRequestEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityBatchRequestEntry create() =>
      ChangeMessageVisibilityBatchRequestEntry._();
  @$core.override
  ChangeMessageVisibilityBatchRequestEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityBatchRequestEntry getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ChangeMessageVisibilityBatchRequestEntry>(create);
  static ChangeMessageVisibilityBatchRequestEntry? _defaultInstance;

  @$pb.TagNumber(134471750)
  $core.String get receipthandle => $_getSZ(0);
  @$pb.TagNumber(134471750)
  set receipthandle($core.String value) => $_setString(0, value);
  @$pb.TagNumber(134471750)
  $core.bool hasReceipthandle() => $_has(0);
  @$pb.TagNumber(134471750)
  void clearReceipthandle() => $_clearField(134471750);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(460820073)
  $core.int get visibilitytimeout => $_getIZ(2);
  @$pb.TagNumber(460820073)
  set visibilitytimeout($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(460820073)
  $core.bool hasVisibilitytimeout() => $_has(2);
  @$pb.TagNumber(460820073)
  void clearVisibilitytimeout() => $_clearField(460820073);
}

class ChangeMessageVisibilityBatchResult extends $pb.GeneratedMessage {
  factory ChangeMessageVisibilityBatchResult({
    $core.Iterable<BatchResultErrorEntry>? failed,
    $core.Iterable<ChangeMessageVisibilityBatchResultEntry>? successful,
  }) {
    final result = create();
    if (failed != null) result.failed.addAll(failed);
    if (successful != null) result.successful.addAll(successful);
    return result;
  }

  ChangeMessageVisibilityBatchResult._();

  factory ChangeMessageVisibilityBatchResult.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeMessageVisibilityBatchResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeMessageVisibilityBatchResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPM<BatchResultErrorEntry>(360301525, _omitFieldNames ? '' : 'failed',
        subBuilder: BatchResultErrorEntry.create)
    ..pPM<ChangeMessageVisibilityBatchResultEntry>(
        412818844, _omitFieldNames ? '' : 'successful',
        subBuilder: ChangeMessageVisibilityBatchResultEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityBatchResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityBatchResult copyWith(
          void Function(ChangeMessageVisibilityBatchResult) updates) =>
      super.copyWith((message) =>
              updates(message as ChangeMessageVisibilityBatchResult))
          as ChangeMessageVisibilityBatchResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityBatchResult create() =>
      ChangeMessageVisibilityBatchResult._();
  @$core.override
  ChangeMessageVisibilityBatchResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityBatchResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeMessageVisibilityBatchResult>(
          create);
  static ChangeMessageVisibilityBatchResult? _defaultInstance;

  @$pb.TagNumber(360301525)
  $pb.PbList<BatchResultErrorEntry> get failed => $_getList(0);

  @$pb.TagNumber(412818844)
  $pb.PbList<ChangeMessageVisibilityBatchResultEntry> get successful =>
      $_getList(1);
}

class ChangeMessageVisibilityBatchResultEntry extends $pb.GeneratedMessage {
  factory ChangeMessageVisibilityBatchResultEntry({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  ChangeMessageVisibilityBatchResultEntry._();

  factory ChangeMessageVisibilityBatchResultEntry.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeMessageVisibilityBatchResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeMessageVisibilityBatchResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityBatchResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityBatchResultEntry copyWith(
          void Function(ChangeMessageVisibilityBatchResultEntry) updates) =>
      super.copyWith((message) =>
              updates(message as ChangeMessageVisibilityBatchResultEntry))
          as ChangeMessageVisibilityBatchResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityBatchResultEntry create() =>
      ChangeMessageVisibilityBatchResultEntry._();
  @$core.override
  ChangeMessageVisibilityBatchResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityBatchResultEntry getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ChangeMessageVisibilityBatchResultEntry>(create);
  static ChangeMessageVisibilityBatchResultEntry? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class ChangeMessageVisibilityRequest extends $pb.GeneratedMessage {
  factory ChangeMessageVisibilityRequest({
    $core.String? receipthandle,
    $core.int? visibilitytimeout,
    $core.String? queueurl,
  }) {
    final result = create();
    if (receipthandle != null) result.receipthandle = receipthandle;
    if (visibilitytimeout != null) result.visibilitytimeout = visibilitytimeout;
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  ChangeMessageVisibilityRequest._();

  factory ChangeMessageVisibilityRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeMessageVisibilityRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeMessageVisibilityRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(134471750, _omitFieldNames ? '' : 'receipthandle')
    ..aI(460820073, _omitFieldNames ? '' : 'visibilitytimeout')
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeMessageVisibilityRequest copyWith(
          void Function(ChangeMessageVisibilityRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ChangeMessageVisibilityRequest))
          as ChangeMessageVisibilityRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityRequest create() =>
      ChangeMessageVisibilityRequest._();
  @$core.override
  ChangeMessageVisibilityRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeMessageVisibilityRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeMessageVisibilityRequest>(create);
  static ChangeMessageVisibilityRequest? _defaultInstance;

  @$pb.TagNumber(134471750)
  $core.String get receipthandle => $_getSZ(0);
  @$pb.TagNumber(134471750)
  set receipthandle($core.String value) => $_setString(0, value);
  @$pb.TagNumber(134471750)
  $core.bool hasReceipthandle() => $_has(0);
  @$pb.TagNumber(134471750)
  void clearReceipthandle() => $_clearField(134471750);

  @$pb.TagNumber(460820073)
  $core.int get visibilitytimeout => $_getIZ(1);
  @$pb.TagNumber(460820073)
  set visibilitytimeout($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(460820073)
  $core.bool hasVisibilitytimeout() => $_has(1);
  @$pb.TagNumber(460820073)
  void clearVisibilitytimeout() => $_clearField(460820073);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(2);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(2, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(2);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class CreateQueueRequest extends $pb.GeneratedMessage {
  factory CreateQueueRequest({
    $core.String? queuename,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
  }) {
    final result = create();
    if (queuename != null) result.queuename = queuename;
    if (attributes != null) result.attributes.addEntries(attributes);
    if (tags != null) result.tags.addEntries(tags);
    return result;
  }

  CreateQueueRequest._();

  factory CreateQueueRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateQueueRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateQueueRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(154364636, _omitFieldNames ? '' : 'queuename')
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'CreateQueueRequest.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sqs'))
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateQueueRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sqs'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateQueueRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateQueueRequest copyWith(void Function(CreateQueueRequest) updates) =>
      super.copyWith((message) => updates(message as CreateQueueRequest))
          as CreateQueueRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateQueueRequest create() => CreateQueueRequest._();
  @$core.override
  CreateQueueRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateQueueRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateQueueRequest>(create);
  static CreateQueueRequest? _defaultInstance;

  @$pb.TagNumber(154364636)
  $core.String get queuename => $_getSZ(0);
  @$pb.TagNumber(154364636)
  set queuename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(154364636)
  $core.bool hasQueuename() => $_has(0);
  @$pb.TagNumber(154364636)
  void clearQueuename() => $_clearField(154364636);

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(1);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(2);
}

class CreateQueueResult extends $pb.GeneratedMessage {
  factory CreateQueueResult({
    $core.String? queueurl,
  }) {
    final result = create();
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  CreateQueueResult._();

  factory CreateQueueResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateQueueResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateQueueResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateQueueResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateQueueResult copyWith(void Function(CreateQueueResult) updates) =>
      super.copyWith((message) => updates(message as CreateQueueResult))
          as CreateQueueResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateQueueResult create() => CreateQueueResult._();
  @$core.override
  CreateQueueResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateQueueResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateQueueResult>(create);
  static CreateQueueResult? _defaultInstance;

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(0);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(0);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class DeleteMessageBatchRequest extends $pb.GeneratedMessage {
  factory DeleteMessageBatchRequest({
    $core.Iterable<DeleteMessageBatchRequestEntry>? entries,
    $core.String? queueurl,
  }) {
    final result = create();
    if (entries != null) result.entries.addAll(entries);
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  DeleteMessageBatchRequest._();

  factory DeleteMessageBatchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMessageBatchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMessageBatchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPM<DeleteMessageBatchRequestEntry>(
        481075860, _omitFieldNames ? '' : 'entries',
        subBuilder: DeleteMessageBatchRequestEntry.create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageBatchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageBatchRequest copyWith(
          void Function(DeleteMessageBatchRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteMessageBatchRequest))
          as DeleteMessageBatchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMessageBatchRequest create() => DeleteMessageBatchRequest._();
  @$core.override
  DeleteMessageBatchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMessageBatchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMessageBatchRequest>(create);
  static DeleteMessageBatchRequest? _defaultInstance;

  @$pb.TagNumber(481075860)
  $pb.PbList<DeleteMessageBatchRequestEntry> get entries => $_getList(0);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(1);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(1);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class DeleteMessageBatchRequestEntry extends $pb.GeneratedMessage {
  factory DeleteMessageBatchRequestEntry({
    $core.String? receipthandle,
    $core.String? id,
  }) {
    final result = create();
    if (receipthandle != null) result.receipthandle = receipthandle;
    if (id != null) result.id = id;
    return result;
  }

  DeleteMessageBatchRequestEntry._();

  factory DeleteMessageBatchRequestEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMessageBatchRequestEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMessageBatchRequestEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(134471750, _omitFieldNames ? '' : 'receipthandle')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageBatchRequestEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageBatchRequestEntry copyWith(
          void Function(DeleteMessageBatchRequestEntry) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteMessageBatchRequestEntry))
          as DeleteMessageBatchRequestEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMessageBatchRequestEntry create() =>
      DeleteMessageBatchRequestEntry._();
  @$core.override
  DeleteMessageBatchRequestEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMessageBatchRequestEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMessageBatchRequestEntry>(create);
  static DeleteMessageBatchRequestEntry? _defaultInstance;

  @$pb.TagNumber(134471750)
  $core.String get receipthandle => $_getSZ(0);
  @$pb.TagNumber(134471750)
  set receipthandle($core.String value) => $_setString(0, value);
  @$pb.TagNumber(134471750)
  $core.bool hasReceipthandle() => $_has(0);
  @$pb.TagNumber(134471750)
  void clearReceipthandle() => $_clearField(134471750);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteMessageBatchResult extends $pb.GeneratedMessage {
  factory DeleteMessageBatchResult({
    $core.Iterable<BatchResultErrorEntry>? failed,
    $core.Iterable<DeleteMessageBatchResultEntry>? successful,
  }) {
    final result = create();
    if (failed != null) result.failed.addAll(failed);
    if (successful != null) result.successful.addAll(successful);
    return result;
  }

  DeleteMessageBatchResult._();

  factory DeleteMessageBatchResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMessageBatchResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMessageBatchResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPM<BatchResultErrorEntry>(360301525, _omitFieldNames ? '' : 'failed',
        subBuilder: BatchResultErrorEntry.create)
    ..pPM<DeleteMessageBatchResultEntry>(
        412818844, _omitFieldNames ? '' : 'successful',
        subBuilder: DeleteMessageBatchResultEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageBatchResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageBatchResult copyWith(
          void Function(DeleteMessageBatchResult) updates) =>
      super.copyWith((message) => updates(message as DeleteMessageBatchResult))
          as DeleteMessageBatchResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMessageBatchResult create() => DeleteMessageBatchResult._();
  @$core.override
  DeleteMessageBatchResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMessageBatchResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMessageBatchResult>(create);
  static DeleteMessageBatchResult? _defaultInstance;

  @$pb.TagNumber(360301525)
  $pb.PbList<BatchResultErrorEntry> get failed => $_getList(0);

  @$pb.TagNumber(412818844)
  $pb.PbList<DeleteMessageBatchResultEntry> get successful => $_getList(1);
}

class DeleteMessageBatchResultEntry extends $pb.GeneratedMessage {
  factory DeleteMessageBatchResultEntry({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  DeleteMessageBatchResultEntry._();

  factory DeleteMessageBatchResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMessageBatchResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMessageBatchResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageBatchResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageBatchResultEntry copyWith(
          void Function(DeleteMessageBatchResultEntry) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteMessageBatchResultEntry))
          as DeleteMessageBatchResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMessageBatchResultEntry create() =>
      DeleteMessageBatchResultEntry._();
  @$core.override
  DeleteMessageBatchResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMessageBatchResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMessageBatchResultEntry>(create);
  static DeleteMessageBatchResultEntry? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteMessageRequest extends $pb.GeneratedMessage {
  factory DeleteMessageRequest({
    $core.String? receipthandle,
    $core.String? queueurl,
  }) {
    final result = create();
    if (receipthandle != null) result.receipthandle = receipthandle;
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  DeleteMessageRequest._();

  factory DeleteMessageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMessageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMessageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(134471750, _omitFieldNames ? '' : 'receipthandle')
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMessageRequest copyWith(void Function(DeleteMessageRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteMessageRequest))
          as DeleteMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMessageRequest create() => DeleteMessageRequest._();
  @$core.override
  DeleteMessageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMessageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMessageRequest>(create);
  static DeleteMessageRequest? _defaultInstance;

  @$pb.TagNumber(134471750)
  $core.String get receipthandle => $_getSZ(0);
  @$pb.TagNumber(134471750)
  set receipthandle($core.String value) => $_setString(0, value);
  @$pb.TagNumber(134471750)
  $core.bool hasReceipthandle() => $_has(0);
  @$pb.TagNumber(134471750)
  void clearReceipthandle() => $_clearField(134471750);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(1);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(1);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class DeleteQueueRequest extends $pb.GeneratedMessage {
  factory DeleteQueueRequest({
    $core.String? queueurl,
  }) {
    final result = create();
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  DeleteQueueRequest._();

  factory DeleteQueueRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteQueueRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteQueueRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteQueueRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteQueueRequest copyWith(void Function(DeleteQueueRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteQueueRequest))
          as DeleteQueueRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteQueueRequest create() => DeleteQueueRequest._();
  @$core.override
  DeleteQueueRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteQueueRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteQueueRequest>(create);
  static DeleteQueueRequest? _defaultInstance;

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(0);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(0);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class EmptyBatchRequest extends $pb.GeneratedMessage {
  factory EmptyBatchRequest({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EmptyBatchRequest._();

  factory EmptyBatchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EmptyBatchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EmptyBatchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EmptyBatchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EmptyBatchRequest copyWith(void Function(EmptyBatchRequest) updates) =>
      super.copyWith((message) => updates(message as EmptyBatchRequest))
          as EmptyBatchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EmptyBatchRequest create() => EmptyBatchRequest._();
  @$core.override
  EmptyBatchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EmptyBatchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EmptyBatchRequest>(create);
  static EmptyBatchRequest? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class GetQueueAttributesRequest extends $pb.GeneratedMessage {
  factory GetQueueAttributesRequest({
    $core.Iterable<QueueAttributeName>? attributenames,
    $core.String? queueurl,
  }) {
    final result = create();
    if (attributenames != null) result.attributenames.addAll(attributenames);
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  GetQueueAttributesRequest._();

  factory GetQueueAttributesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueueAttributesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueueAttributesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pc<QueueAttributeName>(
        394468622, _omitFieldNames ? '' : 'attributenames', $pb.PbFieldType.KE,
        valueOf: QueueAttributeName.valueOf,
        enumValues: QueueAttributeName.values,
        defaultEnumValue:
            QueueAttributeName.QUEUE_ATTRIBUTE_NAME_MAXIMUMMESSAGESIZE)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueueAttributesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueueAttributesRequest copyWith(
          void Function(GetQueueAttributesRequest) updates) =>
      super.copyWith((message) => updates(message as GetQueueAttributesRequest))
          as GetQueueAttributesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueueAttributesRequest create() => GetQueueAttributesRequest._();
  @$core.override
  GetQueueAttributesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueueAttributesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueueAttributesRequest>(create);
  static GetQueueAttributesRequest? _defaultInstance;

  @$pb.TagNumber(394468622)
  $pb.PbList<QueueAttributeName> get attributenames => $_getList(0);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(1);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(1);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class GetQueueAttributesResult extends $pb.GeneratedMessage {
  factory GetQueueAttributesResult({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  GetQueueAttributesResult._();

  factory GetQueueAttributesResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueueAttributesResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueueAttributesResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'GetQueueAttributesResult.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sqs'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueueAttributesResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueueAttributesResult copyWith(
          void Function(GetQueueAttributesResult) updates) =>
      super.copyWith((message) => updates(message as GetQueueAttributesResult))
          as GetQueueAttributesResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueueAttributesResult create() => GetQueueAttributesResult._();
  @$core.override
  GetQueueAttributesResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueueAttributesResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueueAttributesResult>(create);
  static GetQueueAttributesResult? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);
}

class GetQueueUrlRequest extends $pb.GeneratedMessage {
  factory GetQueueUrlRequest({
    $core.String? queuename,
    $core.String? queueownerawsaccountid,
  }) {
    final result = create();
    if (queuename != null) result.queuename = queuename;
    if (queueownerawsaccountid != null)
      result.queueownerawsaccountid = queueownerawsaccountid;
    return result;
  }

  GetQueueUrlRequest._();

  factory GetQueueUrlRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueueUrlRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueueUrlRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(154364636, _omitFieldNames ? '' : 'queuename')
    ..aOS(248813375, _omitFieldNames ? '' : 'queueownerawsaccountid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueueUrlRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueueUrlRequest copyWith(void Function(GetQueueUrlRequest) updates) =>
      super.copyWith((message) => updates(message as GetQueueUrlRequest))
          as GetQueueUrlRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueueUrlRequest create() => GetQueueUrlRequest._();
  @$core.override
  GetQueueUrlRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueueUrlRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueueUrlRequest>(create);
  static GetQueueUrlRequest? _defaultInstance;

  @$pb.TagNumber(154364636)
  $core.String get queuename => $_getSZ(0);
  @$pb.TagNumber(154364636)
  set queuename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(154364636)
  $core.bool hasQueuename() => $_has(0);
  @$pb.TagNumber(154364636)
  void clearQueuename() => $_clearField(154364636);

  @$pb.TagNumber(248813375)
  $core.String get queueownerawsaccountid => $_getSZ(1);
  @$pb.TagNumber(248813375)
  set queueownerawsaccountid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(248813375)
  $core.bool hasQueueownerawsaccountid() => $_has(1);
  @$pb.TagNumber(248813375)
  void clearQueueownerawsaccountid() => $_clearField(248813375);
}

class GetQueueUrlResult extends $pb.GeneratedMessage {
  factory GetQueueUrlResult({
    $core.String? queueurl,
  }) {
    final result = create();
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  GetQueueUrlResult._();

  factory GetQueueUrlResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueueUrlResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueueUrlResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueueUrlResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueueUrlResult copyWith(void Function(GetQueueUrlResult) updates) =>
      super.copyWith((message) => updates(message as GetQueueUrlResult))
          as GetQueueUrlResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueueUrlResult create() => GetQueueUrlResult._();
  @$core.override
  GetQueueUrlResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueueUrlResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueueUrlResult>(create);
  static GetQueueUrlResult? _defaultInstance;

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(0);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(0);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class InvalidAddress extends $pb.GeneratedMessage {
  factory InvalidAddress({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidAddress._();

  factory InvalidAddress.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidAddress.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidAddress',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAddress clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAddress copyWith(void Function(InvalidAddress) updates) =>
      super.copyWith((message) => updates(message as InvalidAddress))
          as InvalidAddress;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidAddress create() => InvalidAddress._();
  @$core.override
  InvalidAddress createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidAddress getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidAddress>(create);
  static InvalidAddress? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidAttributeName extends $pb.GeneratedMessage {
  factory InvalidAttributeName({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidAttributeName._();

  factory InvalidAttributeName.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidAttributeName.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidAttributeName',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAttributeName clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAttributeName copyWith(void Function(InvalidAttributeName) updates) =>
      super.copyWith((message) => updates(message as InvalidAttributeName))
          as InvalidAttributeName;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidAttributeName create() => InvalidAttributeName._();
  @$core.override
  InvalidAttributeName createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidAttributeName getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidAttributeName>(create);
  static InvalidAttributeName? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidAttributeValue extends $pb.GeneratedMessage {
  factory InvalidAttributeValue({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidAttributeValue._();

  factory InvalidAttributeValue.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidAttributeValue.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidAttributeValue',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAttributeValue clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAttributeValue copyWith(
          void Function(InvalidAttributeValue) updates) =>
      super.copyWith((message) => updates(message as InvalidAttributeValue))
          as InvalidAttributeValue;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidAttributeValue create() => InvalidAttributeValue._();
  @$core.override
  InvalidAttributeValue createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidAttributeValue getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidAttributeValue>(create);
  static InvalidAttributeValue? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidBatchEntryId extends $pb.GeneratedMessage {
  factory InvalidBatchEntryId({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidBatchEntryId._();

  factory InvalidBatchEntryId.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidBatchEntryId.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidBatchEntryId',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidBatchEntryId clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidBatchEntryId copyWith(void Function(InvalidBatchEntryId) updates) =>
      super.copyWith((message) => updates(message as InvalidBatchEntryId))
          as InvalidBatchEntryId;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidBatchEntryId create() => InvalidBatchEntryId._();
  @$core.override
  InvalidBatchEntryId createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidBatchEntryId getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidBatchEntryId>(create);
  static InvalidBatchEntryId? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidIdFormat extends $pb.GeneratedMessage {
  factory InvalidIdFormat() => create();

  InvalidIdFormat._();

  factory InvalidIdFormat.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidIdFormat.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidIdFormat',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidIdFormat clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidIdFormat copyWith(void Function(InvalidIdFormat) updates) =>
      super.copyWith((message) => updates(message as InvalidIdFormat))
          as InvalidIdFormat;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidIdFormat create() => InvalidIdFormat._();
  @$core.override
  InvalidIdFormat createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidIdFormat getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidIdFormat>(create);
  static InvalidIdFormat? _defaultInstance;
}

class InvalidMessageContents extends $pb.GeneratedMessage {
  factory InvalidMessageContents({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidMessageContents._();

  factory InvalidMessageContents.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidMessageContents.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidMessageContents',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidMessageContents clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidMessageContents copyWith(
          void Function(InvalidMessageContents) updates) =>
      super.copyWith((message) => updates(message as InvalidMessageContents))
          as InvalidMessageContents;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidMessageContents create() => InvalidMessageContents._();
  @$core.override
  InvalidMessageContents createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidMessageContents getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidMessageContents>(create);
  static InvalidMessageContents? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidSecurity extends $pb.GeneratedMessage {
  factory InvalidSecurity({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidSecurity._();

  factory InvalidSecurity.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidSecurity.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidSecurity',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSecurity clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSecurity copyWith(void Function(InvalidSecurity) updates) =>
      super.copyWith((message) => updates(message as InvalidSecurity))
          as InvalidSecurity;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidSecurity create() => InvalidSecurity._();
  @$core.override
  InvalidSecurity createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidSecurity getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidSecurity>(create);
  static InvalidSecurity? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsAccessDenied extends $pb.GeneratedMessage {
  factory KmsAccessDenied({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsAccessDenied._();

  factory KmsAccessDenied.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsAccessDenied.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsAccessDenied',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsAccessDenied clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsAccessDenied copyWith(void Function(KmsAccessDenied) updates) =>
      super.copyWith((message) => updates(message as KmsAccessDenied))
          as KmsAccessDenied;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsAccessDenied create() => KmsAccessDenied._();
  @$core.override
  KmsAccessDenied createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsAccessDenied getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsAccessDenied>(create);
  static KmsAccessDenied? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsDisabled extends $pb.GeneratedMessage {
  factory KmsDisabled({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsDisabled._();

  factory KmsDisabled.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsDisabled.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsDisabled',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsDisabled clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsDisabled copyWith(void Function(KmsDisabled) updates) =>
      super.copyWith((message) => updates(message as KmsDisabled))
          as KmsDisabled;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsDisabled create() => KmsDisabled._();
  @$core.override
  KmsDisabled createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsDisabled getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsDisabled>(create);
  static KmsDisabled? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsInvalidKeyUsage extends $pb.GeneratedMessage {
  factory KmsInvalidKeyUsage({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsInvalidKeyUsage._();

  factory KmsInvalidKeyUsage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsInvalidKeyUsage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsInvalidKeyUsage',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsInvalidKeyUsage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsInvalidKeyUsage copyWith(void Function(KmsInvalidKeyUsage) updates) =>
      super.copyWith((message) => updates(message as KmsInvalidKeyUsage))
          as KmsInvalidKeyUsage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsInvalidKeyUsage create() => KmsInvalidKeyUsage._();
  @$core.override
  KmsInvalidKeyUsage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsInvalidKeyUsage getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsInvalidKeyUsage>(create);
  static KmsInvalidKeyUsage? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsInvalidState extends $pb.GeneratedMessage {
  factory KmsInvalidState({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsInvalidState._();

  factory KmsInvalidState.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsInvalidState.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsInvalidState',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsInvalidState clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsInvalidState copyWith(void Function(KmsInvalidState) updates) =>
      super.copyWith((message) => updates(message as KmsInvalidState))
          as KmsInvalidState;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsInvalidState create() => KmsInvalidState._();
  @$core.override
  KmsInvalidState createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsInvalidState getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsInvalidState>(create);
  static KmsInvalidState? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsNotFound extends $pb.GeneratedMessage {
  factory KmsNotFound({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsNotFound._();

  factory KmsNotFound.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsNotFound.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsNotFound',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsNotFound clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsNotFound copyWith(void Function(KmsNotFound) updates) =>
      super.copyWith((message) => updates(message as KmsNotFound))
          as KmsNotFound;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsNotFound create() => KmsNotFound._();
  @$core.override
  KmsNotFound createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsNotFound getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsNotFound>(create);
  static KmsNotFound? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsOptInRequired extends $pb.GeneratedMessage {
  factory KmsOptInRequired({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsOptInRequired._();

  factory KmsOptInRequired.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsOptInRequired.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsOptInRequired',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsOptInRequired clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsOptInRequired copyWith(void Function(KmsOptInRequired) updates) =>
      super.copyWith((message) => updates(message as KmsOptInRequired))
          as KmsOptInRequired;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsOptInRequired create() => KmsOptInRequired._();
  @$core.override
  KmsOptInRequired createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsOptInRequired getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsOptInRequired>(create);
  static KmsOptInRequired? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KmsThrottled extends $pb.GeneratedMessage {
  factory KmsThrottled({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KmsThrottled._();

  factory KmsThrottled.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KmsThrottled.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KmsThrottled',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsThrottled clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KmsThrottled copyWith(void Function(KmsThrottled) updates) =>
      super.copyWith((message) => updates(message as KmsThrottled))
          as KmsThrottled;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KmsThrottled create() => KmsThrottled._();
  @$core.override
  KmsThrottled createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KmsThrottled getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KmsThrottled>(create);
  static KmsThrottled? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ListDeadLetterSourceQueuesRequest extends $pb.GeneratedMessage {
  factory ListDeadLetterSourceQueuesRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? queueurl,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  ListDeadLetterSourceQueuesRequest._();

  factory ListDeadLetterSourceQueuesRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDeadLetterSourceQueuesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDeadLetterSourceQueuesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDeadLetterSourceQueuesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDeadLetterSourceQueuesRequest copyWith(
          void Function(ListDeadLetterSourceQueuesRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ListDeadLetterSourceQueuesRequest))
          as ListDeadLetterSourceQueuesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDeadLetterSourceQueuesRequest create() =>
      ListDeadLetterSourceQueuesRequest._();
  @$core.override
  ListDeadLetterSourceQueuesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDeadLetterSourceQueuesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDeadLetterSourceQueuesRequest>(
          create);
  static ListDeadLetterSourceQueuesRequest? _defaultInstance;

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

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(2);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(2, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(2);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class ListDeadLetterSourceQueuesResult extends $pb.GeneratedMessage {
  factory ListDeadLetterSourceQueuesResult({
    $core.String? nexttoken,
    $core.Iterable<$core.String>? queueurls,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (queueurls != null) result.queueurls.addAll(queueurls);
    return result;
  }

  ListDeadLetterSourceQueuesResult._();

  factory ListDeadLetterSourceQueuesResult.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDeadLetterSourceQueuesResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDeadLetterSourceQueuesResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPS(475536111, _omitFieldNames ? '' : 'queueurls')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDeadLetterSourceQueuesResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDeadLetterSourceQueuesResult copyWith(
          void Function(ListDeadLetterSourceQueuesResult) updates) =>
      super.copyWith(
              (message) => updates(message as ListDeadLetterSourceQueuesResult))
          as ListDeadLetterSourceQueuesResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDeadLetterSourceQueuesResult create() =>
      ListDeadLetterSourceQueuesResult._();
  @$core.override
  ListDeadLetterSourceQueuesResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDeadLetterSourceQueuesResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDeadLetterSourceQueuesResult>(
          create);
  static ListDeadLetterSourceQueuesResult? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(475536111)
  $pb.PbList<$core.String> get queueurls => $_getList(1);
}

class ListMessageMoveTasksRequest extends $pb.GeneratedMessage {
  factory ListMessageMoveTasksRequest({
    $core.int? maxresults,
    $core.String? sourcearn,
  }) {
    final result = create();
    if (maxresults != null) result.maxresults = maxresults;
    if (sourcearn != null) result.sourcearn = sourcearn;
    return result;
  }

  ListMessageMoveTasksRequest._();

  factory ListMessageMoveTasksRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMessageMoveTasksRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMessageMoveTasksRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(439903072, _omitFieldNames ? '' : 'sourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessageMoveTasksRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessageMoveTasksRequest copyWith(
          void Function(ListMessageMoveTasksRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListMessageMoveTasksRequest))
          as ListMessageMoveTasksRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMessageMoveTasksRequest create() =>
      ListMessageMoveTasksRequest._();
  @$core.override
  ListMessageMoveTasksRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMessageMoveTasksRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMessageMoveTasksRequest>(create);
  static ListMessageMoveTasksRequest? _defaultInstance;

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(0);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(0);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(439903072)
  $core.String get sourcearn => $_getSZ(1);
  @$pb.TagNumber(439903072)
  set sourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(439903072)
  $core.bool hasSourcearn() => $_has(1);
  @$pb.TagNumber(439903072)
  void clearSourcearn() => $_clearField(439903072);
}

class ListMessageMoveTasksResult extends $pb.GeneratedMessage {
  factory ListMessageMoveTasksResult({
    $core.Iterable<ListMessageMoveTasksResultEntry>? results,
  }) {
    final result = create();
    if (results != null) result.results.addAll(results);
    return result;
  }

  ListMessageMoveTasksResult._();

  factory ListMessageMoveTasksResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMessageMoveTasksResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMessageMoveTasksResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPM<ListMessageMoveTasksResultEntry>(
        486024854, _omitFieldNames ? '' : 'results',
        subBuilder: ListMessageMoveTasksResultEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessageMoveTasksResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessageMoveTasksResult copyWith(
          void Function(ListMessageMoveTasksResult) updates) =>
      super.copyWith(
              (message) => updates(message as ListMessageMoveTasksResult))
          as ListMessageMoveTasksResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMessageMoveTasksResult create() => ListMessageMoveTasksResult._();
  @$core.override
  ListMessageMoveTasksResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMessageMoveTasksResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMessageMoveTasksResult>(create);
  static ListMessageMoveTasksResult? _defaultInstance;

  @$pb.TagNumber(486024854)
  $pb.PbList<ListMessageMoveTasksResultEntry> get results => $_getList(0);
}

class ListMessageMoveTasksResultEntry extends $pb.GeneratedMessage {
  factory ListMessageMoveTasksResultEntry({
    $core.String? status,
    $fixnum.Int64? approximatenumberofmessagestomove,
    $core.String? taskhandle,
    $core.String? failurereason,
    $core.int? maxnumberofmessagespersecond,
    $core.String? destinationarn,
    $fixnum.Int64? startedtimestamp,
    $core.String? sourcearn,
    $fixnum.Int64? approximatenumberofmessagesmoved,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (approximatenumberofmessagestomove != null)
      result.approximatenumberofmessagestomove =
          approximatenumberofmessagestomove;
    if (taskhandle != null) result.taskhandle = taskhandle;
    if (failurereason != null) result.failurereason = failurereason;
    if (maxnumberofmessagespersecond != null)
      result.maxnumberofmessagespersecond = maxnumberofmessagespersecond;
    if (destinationarn != null) result.destinationarn = destinationarn;
    if (startedtimestamp != null) result.startedtimestamp = startedtimestamp;
    if (sourcearn != null) result.sourcearn = sourcearn;
    if (approximatenumberofmessagesmoved != null)
      result.approximatenumberofmessagesmoved =
          approximatenumberofmessagesmoved;
    return result;
  }

  ListMessageMoveTasksResultEntry._();

  factory ListMessageMoveTasksResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMessageMoveTasksResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMessageMoveTasksResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(6222352, _omitFieldNames ? '' : 'status')
    ..aInt64(
        108945690, _omitFieldNames ? '' : 'approximatenumberofmessagestomove')
    ..aOS(190544291, _omitFieldNames ? '' : 'taskhandle')
    ..aOS(232322142, _omitFieldNames ? '' : 'failurereason')
    ..aI(335779921, _omitFieldNames ? '' : 'maxnumberofmessagespersecond')
    ..aOS(375726595, _omitFieldNames ? '' : 'destinationarn')
    ..aInt64(397447975, _omitFieldNames ? '' : 'startedtimestamp')
    ..aOS(439903072, _omitFieldNames ? '' : 'sourcearn')
    ..aInt64(
        510181657, _omitFieldNames ? '' : 'approximatenumberofmessagesmoved')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessageMoveTasksResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMessageMoveTasksResultEntry copyWith(
          void Function(ListMessageMoveTasksResultEntry) updates) =>
      super.copyWith(
              (message) => updates(message as ListMessageMoveTasksResultEntry))
          as ListMessageMoveTasksResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMessageMoveTasksResultEntry create() =>
      ListMessageMoveTasksResultEntry._();
  @$core.override
  ListMessageMoveTasksResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMessageMoveTasksResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMessageMoveTasksResultEntry>(
          create);
  static ListMessageMoveTasksResultEntry? _defaultInstance;

  @$pb.TagNumber(6222352)
  $core.String get status => $_getSZ(0);
  @$pb.TagNumber(6222352)
  set status($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(108945690)
  $fixnum.Int64 get approximatenumberofmessagestomove => $_getI64(1);
  @$pb.TagNumber(108945690)
  set approximatenumberofmessagestomove($fixnum.Int64 value) =>
      $_setInt64(1, value);
  @$pb.TagNumber(108945690)
  $core.bool hasApproximatenumberofmessagestomove() => $_has(1);
  @$pb.TagNumber(108945690)
  void clearApproximatenumberofmessagestomove() => $_clearField(108945690);

  @$pb.TagNumber(190544291)
  $core.String get taskhandle => $_getSZ(2);
  @$pb.TagNumber(190544291)
  set taskhandle($core.String value) => $_setString(2, value);
  @$pb.TagNumber(190544291)
  $core.bool hasTaskhandle() => $_has(2);
  @$pb.TagNumber(190544291)
  void clearTaskhandle() => $_clearField(190544291);

  @$pb.TagNumber(232322142)
  $core.String get failurereason => $_getSZ(3);
  @$pb.TagNumber(232322142)
  set failurereason($core.String value) => $_setString(3, value);
  @$pb.TagNumber(232322142)
  $core.bool hasFailurereason() => $_has(3);
  @$pb.TagNumber(232322142)
  void clearFailurereason() => $_clearField(232322142);

  @$pb.TagNumber(335779921)
  $core.int get maxnumberofmessagespersecond => $_getIZ(4);
  @$pb.TagNumber(335779921)
  set maxnumberofmessagespersecond($core.int value) =>
      $_setSignedInt32(4, value);
  @$pb.TagNumber(335779921)
  $core.bool hasMaxnumberofmessagespersecond() => $_has(4);
  @$pb.TagNumber(335779921)
  void clearMaxnumberofmessagespersecond() => $_clearField(335779921);

  @$pb.TagNumber(375726595)
  $core.String get destinationarn => $_getSZ(5);
  @$pb.TagNumber(375726595)
  set destinationarn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(375726595)
  $core.bool hasDestinationarn() => $_has(5);
  @$pb.TagNumber(375726595)
  void clearDestinationarn() => $_clearField(375726595);

  @$pb.TagNumber(397447975)
  $fixnum.Int64 get startedtimestamp => $_getI64(6);
  @$pb.TagNumber(397447975)
  set startedtimestamp($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(397447975)
  $core.bool hasStartedtimestamp() => $_has(6);
  @$pb.TagNumber(397447975)
  void clearStartedtimestamp() => $_clearField(397447975);

  @$pb.TagNumber(439903072)
  $core.String get sourcearn => $_getSZ(7);
  @$pb.TagNumber(439903072)
  set sourcearn($core.String value) => $_setString(7, value);
  @$pb.TagNumber(439903072)
  $core.bool hasSourcearn() => $_has(7);
  @$pb.TagNumber(439903072)
  void clearSourcearn() => $_clearField(439903072);

  @$pb.TagNumber(510181657)
  $fixnum.Int64 get approximatenumberofmessagesmoved => $_getI64(8);
  @$pb.TagNumber(510181657)
  set approximatenumberofmessagesmoved($fixnum.Int64 value) =>
      $_setInt64(8, value);
  @$pb.TagNumber(510181657)
  $core.bool hasApproximatenumberofmessagesmoved() => $_has(8);
  @$pb.TagNumber(510181657)
  void clearApproximatenumberofmessagesmoved() => $_clearField(510181657);
}

class ListQueueTagsRequest extends $pb.GeneratedMessage {
  factory ListQueueTagsRequest({
    $core.String? queueurl,
  }) {
    final result = create();
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  ListQueueTagsRequest._();

  factory ListQueueTagsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueueTagsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueueTagsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueueTagsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueueTagsRequest copyWith(void Function(ListQueueTagsRequest) updates) =>
      super.copyWith((message) => updates(message as ListQueueTagsRequest))
          as ListQueueTagsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueueTagsRequest create() => ListQueueTagsRequest._();
  @$core.override
  ListQueueTagsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueueTagsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueueTagsRequest>(create);
  static ListQueueTagsRequest? _defaultInstance;

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(0);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(0);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class ListQueueTagsResult extends $pb.GeneratedMessage {
  factory ListQueueTagsResult({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
  }) {
    final result = create();
    if (tags != null) result.tags.addEntries(tags);
    return result;
  }

  ListQueueTagsResult._();

  factory ListQueueTagsResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueueTagsResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueueTagsResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'ListQueueTagsResult.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sqs'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueueTagsResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueueTagsResult copyWith(void Function(ListQueueTagsResult) updates) =>
      super.copyWith((message) => updates(message as ListQueueTagsResult))
          as ListQueueTagsResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueueTagsResult create() => ListQueueTagsResult._();
  @$core.override
  ListQueueTagsResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueueTagsResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueueTagsResult>(create);
  static ListQueueTagsResult? _defaultInstance;

  @$pb.TagNumber(381526209)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(0);
}

class ListQueuesRequest extends $pb.GeneratedMessage {
  factory ListQueuesRequest({
    $core.String? nexttoken,
    $core.String? queuenameprefix,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (queuenameprefix != null) result.queuenameprefix = queuenameprefix;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListQueuesRequest._();

  factory ListQueuesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueuesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueuesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(269478416, _omitFieldNames ? '' : 'queuenameprefix')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueuesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueuesRequest copyWith(void Function(ListQueuesRequest) updates) =>
      super.copyWith((message) => updates(message as ListQueuesRequest))
          as ListQueuesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueuesRequest create() => ListQueuesRequest._();
  @$core.override
  ListQueuesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueuesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueuesRequest>(create);
  static ListQueuesRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(269478416)
  $core.String get queuenameprefix => $_getSZ(1);
  @$pb.TagNumber(269478416)
  set queuenameprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(269478416)
  $core.bool hasQueuenameprefix() => $_has(1);
  @$pb.TagNumber(269478416)
  void clearQueuenameprefix() => $_clearField(269478416);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(2);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListQueuesResult extends $pb.GeneratedMessage {
  factory ListQueuesResult({
    $core.Iterable<$core.String>? queueurls,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (queueurls != null) result.queueurls.addAll(queueurls);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListQueuesResult._();

  factory ListQueuesResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueuesResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueuesResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPS(62522575, _omitFieldNames ? '' : 'queueurls')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueuesResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueuesResult copyWith(void Function(ListQueuesResult) updates) =>
      super.copyWith((message) => updates(message as ListQueuesResult))
          as ListQueuesResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueuesResult create() => ListQueuesResult._();
  @$core.override
  ListQueuesResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueuesResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueuesResult>(create);
  static ListQueuesResult? _defaultInstance;

  @$pb.TagNumber(62522575)
  $pb.PbList<$core.String> get queueurls => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class Message extends $pb.GeneratedMessage {
  factory Message({
    $core.String? body,
    $core.Iterable<$core.MapEntry<$core.String, MessageAttributeValue>>?
        messageattributes,
    $core.String? receipthandle,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.String? md5ofbody,
    $core.String? messageid,
    $core.String? md5ofmessageattributes,
  }) {
    final result = create();
    if (body != null) result.body = body;
    if (messageattributes != null)
      result.messageattributes.addEntries(messageattributes);
    if (receipthandle != null) result.receipthandle = receipthandle;
    if (attributes != null) result.attributes.addEntries(attributes);
    if (md5ofbody != null) result.md5ofbody = md5ofbody;
    if (messageid != null) result.messageid = messageid;
    if (md5ofmessageattributes != null)
      result.md5ofmessageattributes = md5ofmessageattributes;
    return result;
  }

  Message._();

  factory Message.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Message.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Message',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(42602646, _omitFieldNames ? '' : 'body')
    ..m<$core.String, MessageAttributeValue>(
        56443766, _omitFieldNames ? '' : 'messageattributes',
        entryClassName: 'Message.MessageattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MessageAttributeValue.create,
        valueDefaultOrMaker: MessageAttributeValue.getDefault,
        packageName: const $pb.PackageName('sqs'))
    ..aOS(134471750, _omitFieldNames ? '' : 'receipthandle')
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'Message.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sqs'))
    ..aOS(229335371, _omitFieldNames ? '' : 'md5ofbody')
    ..aOS(360526634, _omitFieldNames ? '' : 'messageid')
    ..aOS(499647077, _omitFieldNames ? '' : 'md5ofmessageattributes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Message clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Message copyWith(void Function(Message) updates) =>
      super.copyWith((message) => updates(message as Message)) as Message;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Message create() => Message._();
  @$core.override
  Message createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Message getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Message>(create);
  static Message? _defaultInstance;

  @$pb.TagNumber(42602646)
  $core.String get body => $_getSZ(0);
  @$pb.TagNumber(42602646)
  set body($core.String value) => $_setString(0, value);
  @$pb.TagNumber(42602646)
  $core.bool hasBody() => $_has(0);
  @$pb.TagNumber(42602646)
  void clearBody() => $_clearField(42602646);

  @$pb.TagNumber(56443766)
  $pb.PbMap<$core.String, MessageAttributeValue> get messageattributes =>
      $_getMap(1);

  @$pb.TagNumber(134471750)
  $core.String get receipthandle => $_getSZ(2);
  @$pb.TagNumber(134471750)
  set receipthandle($core.String value) => $_setString(2, value);
  @$pb.TagNumber(134471750)
  $core.bool hasReceipthandle() => $_has(2);
  @$pb.TagNumber(134471750)
  void clearReceipthandle() => $_clearField(134471750);

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(3);

  @$pb.TagNumber(229335371)
  $core.String get md5ofbody => $_getSZ(4);
  @$pb.TagNumber(229335371)
  set md5ofbody($core.String value) => $_setString(4, value);
  @$pb.TagNumber(229335371)
  $core.bool hasMd5ofbody() => $_has(4);
  @$pb.TagNumber(229335371)
  void clearMd5ofbody() => $_clearField(229335371);

  @$pb.TagNumber(360526634)
  $core.String get messageid => $_getSZ(5);
  @$pb.TagNumber(360526634)
  set messageid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(360526634)
  $core.bool hasMessageid() => $_has(5);
  @$pb.TagNumber(360526634)
  void clearMessageid() => $_clearField(360526634);

  @$pb.TagNumber(499647077)
  $core.String get md5ofmessageattributes => $_getSZ(6);
  @$pb.TagNumber(499647077)
  set md5ofmessageattributes($core.String value) => $_setString(6, value);
  @$pb.TagNumber(499647077)
  $core.bool hasMd5ofmessageattributes() => $_has(6);
  @$pb.TagNumber(499647077)
  void clearMd5ofmessageattributes() => $_clearField(499647077);
}

class MessageAttributeValue extends $pb.GeneratedMessage {
  factory MessageAttributeValue({
    $core.String? datatype,
    $core.Iterable<$core.String>? stringlistvalues,
    $core.Iterable<$core.List<$core.int>>? binarylistvalues,
    $core.String? stringvalue,
    $core.List<$core.int>? binaryvalue,
  }) {
    final result = create();
    if (datatype != null) result.datatype = datatype;
    if (stringlistvalues != null)
      result.stringlistvalues.addAll(stringlistvalues);
    if (binarylistvalues != null)
      result.binarylistvalues.addAll(binarylistvalues);
    if (stringvalue != null) result.stringvalue = stringvalue;
    if (binaryvalue != null) result.binaryvalue = binaryvalue;
    return result;
  }

  MessageAttributeValue._();

  factory MessageAttributeValue.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MessageAttributeValue.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MessageAttributeValue',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(67988590, _omitFieldNames ? '' : 'datatype')
    ..pPS(86527527, _omitFieldNames ? '' : 'stringlistvalues')
    ..p<$core.List<$core.int>>(158345259,
        _omitFieldNames ? '' : 'binarylistvalues', $pb.PbFieldType.PY)
    ..aOS(184416138, _omitFieldNames ? '' : 'stringvalue')
    ..a<$core.List<$core.int>>(
        255476278, _omitFieldNames ? '' : 'binaryvalue', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageAttributeValue clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageAttributeValue copyWith(
          void Function(MessageAttributeValue) updates) =>
      super.copyWith((message) => updates(message as MessageAttributeValue))
          as MessageAttributeValue;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MessageAttributeValue create() => MessageAttributeValue._();
  @$core.override
  MessageAttributeValue createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MessageAttributeValue getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MessageAttributeValue>(create);
  static MessageAttributeValue? _defaultInstance;

  @$pb.TagNumber(67988590)
  $core.String get datatype => $_getSZ(0);
  @$pb.TagNumber(67988590)
  set datatype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(67988590)
  $core.bool hasDatatype() => $_has(0);
  @$pb.TagNumber(67988590)
  void clearDatatype() => $_clearField(67988590);

  @$pb.TagNumber(86527527)
  $pb.PbList<$core.String> get stringlistvalues => $_getList(1);

  @$pb.TagNumber(158345259)
  $pb.PbList<$core.List<$core.int>> get binarylistvalues => $_getList(2);

  @$pb.TagNumber(184416138)
  $core.String get stringvalue => $_getSZ(3);
  @$pb.TagNumber(184416138)
  set stringvalue($core.String value) => $_setString(3, value);
  @$pb.TagNumber(184416138)
  $core.bool hasStringvalue() => $_has(3);
  @$pb.TagNumber(184416138)
  void clearStringvalue() => $_clearField(184416138);

  @$pb.TagNumber(255476278)
  $core.List<$core.int> get binaryvalue => $_getN(4);
  @$pb.TagNumber(255476278)
  set binaryvalue($core.List<$core.int> value) => $_setBytes(4, value);
  @$pb.TagNumber(255476278)
  $core.bool hasBinaryvalue() => $_has(4);
  @$pb.TagNumber(255476278)
  void clearBinaryvalue() => $_clearField(255476278);
}

class MessageNotInflight extends $pb.GeneratedMessage {
  factory MessageNotInflight() => create();

  MessageNotInflight._();

  factory MessageNotInflight.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MessageNotInflight.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MessageNotInflight',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageNotInflight clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageNotInflight copyWith(void Function(MessageNotInflight) updates) =>
      super.copyWith((message) => updates(message as MessageNotInflight))
          as MessageNotInflight;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MessageNotInflight create() => MessageNotInflight._();
  @$core.override
  MessageNotInflight createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MessageNotInflight getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MessageNotInflight>(create);
  static MessageNotInflight? _defaultInstance;
}

class MessageSystemAttributeValue extends $pb.GeneratedMessage {
  factory MessageSystemAttributeValue({
    $core.String? datatype,
    $core.Iterable<$core.String>? stringlistvalues,
    $core.Iterable<$core.List<$core.int>>? binarylistvalues,
    $core.String? stringvalue,
    $core.List<$core.int>? binaryvalue,
  }) {
    final result = create();
    if (datatype != null) result.datatype = datatype;
    if (stringlistvalues != null)
      result.stringlistvalues.addAll(stringlistvalues);
    if (binarylistvalues != null)
      result.binarylistvalues.addAll(binarylistvalues);
    if (stringvalue != null) result.stringvalue = stringvalue;
    if (binaryvalue != null) result.binaryvalue = binaryvalue;
    return result;
  }

  MessageSystemAttributeValue._();

  factory MessageSystemAttributeValue.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MessageSystemAttributeValue.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MessageSystemAttributeValue',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(67988590, _omitFieldNames ? '' : 'datatype')
    ..pPS(86527527, _omitFieldNames ? '' : 'stringlistvalues')
    ..p<$core.List<$core.int>>(158345259,
        _omitFieldNames ? '' : 'binarylistvalues', $pb.PbFieldType.PY)
    ..aOS(184416138, _omitFieldNames ? '' : 'stringvalue')
    ..a<$core.List<$core.int>>(
        255476278, _omitFieldNames ? '' : 'binaryvalue', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageSystemAttributeValue clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageSystemAttributeValue copyWith(
          void Function(MessageSystemAttributeValue) updates) =>
      super.copyWith(
              (message) => updates(message as MessageSystemAttributeValue))
          as MessageSystemAttributeValue;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MessageSystemAttributeValue create() =>
      MessageSystemAttributeValue._();
  @$core.override
  MessageSystemAttributeValue createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MessageSystemAttributeValue getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MessageSystemAttributeValue>(create);
  static MessageSystemAttributeValue? _defaultInstance;

  @$pb.TagNumber(67988590)
  $core.String get datatype => $_getSZ(0);
  @$pb.TagNumber(67988590)
  set datatype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(67988590)
  $core.bool hasDatatype() => $_has(0);
  @$pb.TagNumber(67988590)
  void clearDatatype() => $_clearField(67988590);

  @$pb.TagNumber(86527527)
  $pb.PbList<$core.String> get stringlistvalues => $_getList(1);

  @$pb.TagNumber(158345259)
  $pb.PbList<$core.List<$core.int>> get binarylistvalues => $_getList(2);

  @$pb.TagNumber(184416138)
  $core.String get stringvalue => $_getSZ(3);
  @$pb.TagNumber(184416138)
  set stringvalue($core.String value) => $_setString(3, value);
  @$pb.TagNumber(184416138)
  $core.bool hasStringvalue() => $_has(3);
  @$pb.TagNumber(184416138)
  void clearStringvalue() => $_clearField(184416138);

  @$pb.TagNumber(255476278)
  $core.List<$core.int> get binaryvalue => $_getN(4);
  @$pb.TagNumber(255476278)
  set binaryvalue($core.List<$core.int> value) => $_setBytes(4, value);
  @$pb.TagNumber(255476278)
  $core.bool hasBinaryvalue() => $_has(4);
  @$pb.TagNumber(255476278)
  void clearBinaryvalue() => $_clearField(255476278);
}

class OverLimit extends $pb.GeneratedMessage {
  factory OverLimit({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  OverLimit._();

  factory OverLimit.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OverLimit.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OverLimit',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OverLimit clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OverLimit copyWith(void Function(OverLimit) updates) =>
      super.copyWith((message) => updates(message as OverLimit)) as OverLimit;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OverLimit create() => OverLimit._();
  @$core.override
  OverLimit createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OverLimit getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<OverLimit>(create);
  static OverLimit? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PurgeQueueInProgress extends $pb.GeneratedMessage {
  factory PurgeQueueInProgress({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PurgeQueueInProgress._();

  factory PurgeQueueInProgress.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PurgeQueueInProgress.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PurgeQueueInProgress',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PurgeQueueInProgress clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PurgeQueueInProgress copyWith(void Function(PurgeQueueInProgress) updates) =>
      super.copyWith((message) => updates(message as PurgeQueueInProgress))
          as PurgeQueueInProgress;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PurgeQueueInProgress create() => PurgeQueueInProgress._();
  @$core.override
  PurgeQueueInProgress createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PurgeQueueInProgress getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PurgeQueueInProgress>(create);
  static PurgeQueueInProgress? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PurgeQueueRequest extends $pb.GeneratedMessage {
  factory PurgeQueueRequest({
    $core.String? queueurl,
  }) {
    final result = create();
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  PurgeQueueRequest._();

  factory PurgeQueueRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PurgeQueueRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PurgeQueueRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PurgeQueueRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PurgeQueueRequest copyWith(void Function(PurgeQueueRequest) updates) =>
      super.copyWith((message) => updates(message as PurgeQueueRequest))
          as PurgeQueueRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PurgeQueueRequest create() => PurgeQueueRequest._();
  @$core.override
  PurgeQueueRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PurgeQueueRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PurgeQueueRequest>(create);
  static PurgeQueueRequest? _defaultInstance;

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(0);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(0);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class QueueDeletedRecently extends $pb.GeneratedMessage {
  factory QueueDeletedRecently({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  QueueDeletedRecently._();

  factory QueueDeletedRecently.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueueDeletedRecently.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueueDeletedRecently',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueueDeletedRecently clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueueDeletedRecently copyWith(void Function(QueueDeletedRecently) updates) =>
      super.copyWith((message) => updates(message as QueueDeletedRecently))
          as QueueDeletedRecently;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueueDeletedRecently create() => QueueDeletedRecently._();
  @$core.override
  QueueDeletedRecently createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueueDeletedRecently getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueueDeletedRecently>(create);
  static QueueDeletedRecently? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class QueueDoesNotExist extends $pb.GeneratedMessage {
  factory QueueDoesNotExist({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  QueueDoesNotExist._();

  factory QueueDoesNotExist.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueueDoesNotExist.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueueDoesNotExist',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueueDoesNotExist clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueueDoesNotExist copyWith(void Function(QueueDoesNotExist) updates) =>
      super.copyWith((message) => updates(message as QueueDoesNotExist))
          as QueueDoesNotExist;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueueDoesNotExist create() => QueueDoesNotExist._();
  @$core.override
  QueueDoesNotExist createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueueDoesNotExist getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueueDoesNotExist>(create);
  static QueueDoesNotExist? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class QueueNameExists extends $pb.GeneratedMessage {
  factory QueueNameExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  QueueNameExists._();

  factory QueueNameExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueueNameExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueueNameExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueueNameExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueueNameExists copyWith(void Function(QueueNameExists) updates) =>
      super.copyWith((message) => updates(message as QueueNameExists))
          as QueueNameExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueueNameExists create() => QueueNameExists._();
  @$core.override
  QueueNameExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueueNameExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueueNameExists>(create);
  static QueueNameExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ReceiptHandleIsInvalid extends $pb.GeneratedMessage {
  factory ReceiptHandleIsInvalid({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ReceiptHandleIsInvalid._();

  factory ReceiptHandleIsInvalid.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReceiptHandleIsInvalid.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReceiptHandleIsInvalid',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReceiptHandleIsInvalid clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReceiptHandleIsInvalid copyWith(
          void Function(ReceiptHandleIsInvalid) updates) =>
      super.copyWith((message) => updates(message as ReceiptHandleIsInvalid))
          as ReceiptHandleIsInvalid;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReceiptHandleIsInvalid create() => ReceiptHandleIsInvalid._();
  @$core.override
  ReceiptHandleIsInvalid createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReceiptHandleIsInvalid getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReceiptHandleIsInvalid>(create);
  static ReceiptHandleIsInvalid? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ReceiveMessageRequest extends $pb.GeneratedMessage {
  factory ReceiveMessageRequest({
    $core.Iterable<MessageSystemAttributeName>? messagesystemattributenames,
    $core.int? maxnumberofmessages,
    $core.Iterable<$core.String>? messageattributenames,
    $core.Iterable<QueueAttributeName>? attributenames,
    $core.int? waittimeseconds,
    $core.String? receiverequestattemptid,
    $core.int? visibilitytimeout,
    $core.String? queueurl,
  }) {
    final result = create();
    if (messagesystemattributenames != null)
      result.messagesystemattributenames.addAll(messagesystemattributenames);
    if (maxnumberofmessages != null)
      result.maxnumberofmessages = maxnumberofmessages;
    if (messageattributenames != null)
      result.messageattributenames.addAll(messageattributenames);
    if (attributenames != null) result.attributenames.addAll(attributenames);
    if (waittimeseconds != null) result.waittimeseconds = waittimeseconds;
    if (receiverequestattemptid != null)
      result.receiverequestattemptid = receiverequestattemptid;
    if (visibilitytimeout != null) result.visibilitytimeout = visibilitytimeout;
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  ReceiveMessageRequest._();

  factory ReceiveMessageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReceiveMessageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReceiveMessageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pc<MessageSystemAttributeName>(
        42109014,
        _omitFieldNames ? '' : 'messagesystemattributenames',
        $pb.PbFieldType.KE,
        valueOf: MessageSystemAttributeName.valueOf,
        enumValues: MessageSystemAttributeName.values,
        defaultEnumValue: MessageSystemAttributeName
            .MESSAGE_SYSTEM_ATTRIBUTE_NAME_SEQUENCENUMBER)
    ..aI(139866822, _omitFieldNames ? '' : 'maxnumberofmessages')
    ..pPS(332558373, _omitFieldNames ? '' : 'messageattributenames')
    ..pc<QueueAttributeName>(
        394468622, _omitFieldNames ? '' : 'attributenames', $pb.PbFieldType.KE,
        valueOf: QueueAttributeName.valueOf,
        enumValues: QueueAttributeName.values,
        defaultEnumValue:
            QueueAttributeName.QUEUE_ATTRIBUTE_NAME_MAXIMUMMESSAGESIZE)
    ..aI(398991863, _omitFieldNames ? '' : 'waittimeseconds')
    ..aOS(455135954, _omitFieldNames ? '' : 'receiverequestattemptid')
    ..aI(460820073, _omitFieldNames ? '' : 'visibilitytimeout')
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReceiveMessageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReceiveMessageRequest copyWith(
          void Function(ReceiveMessageRequest) updates) =>
      super.copyWith((message) => updates(message as ReceiveMessageRequest))
          as ReceiveMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReceiveMessageRequest create() => ReceiveMessageRequest._();
  @$core.override
  ReceiveMessageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReceiveMessageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReceiveMessageRequest>(create);
  static ReceiveMessageRequest? _defaultInstance;

  @$pb.TagNumber(42109014)
  $pb.PbList<MessageSystemAttributeName> get messagesystemattributenames =>
      $_getList(0);

  @$pb.TagNumber(139866822)
  $core.int get maxnumberofmessages => $_getIZ(1);
  @$pb.TagNumber(139866822)
  set maxnumberofmessages($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(139866822)
  $core.bool hasMaxnumberofmessages() => $_has(1);
  @$pb.TagNumber(139866822)
  void clearMaxnumberofmessages() => $_clearField(139866822);

  @$pb.TagNumber(332558373)
  $pb.PbList<$core.String> get messageattributenames => $_getList(2);

  @$pb.TagNumber(394468622)
  $pb.PbList<QueueAttributeName> get attributenames => $_getList(3);

  @$pb.TagNumber(398991863)
  $core.int get waittimeseconds => $_getIZ(4);
  @$pb.TagNumber(398991863)
  set waittimeseconds($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(398991863)
  $core.bool hasWaittimeseconds() => $_has(4);
  @$pb.TagNumber(398991863)
  void clearWaittimeseconds() => $_clearField(398991863);

  @$pb.TagNumber(455135954)
  $core.String get receiverequestattemptid => $_getSZ(5);
  @$pb.TagNumber(455135954)
  set receiverequestattemptid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(455135954)
  $core.bool hasReceiverequestattemptid() => $_has(5);
  @$pb.TagNumber(455135954)
  void clearReceiverequestattemptid() => $_clearField(455135954);

  @$pb.TagNumber(460820073)
  $core.int get visibilitytimeout => $_getIZ(6);
  @$pb.TagNumber(460820073)
  set visibilitytimeout($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(460820073)
  $core.bool hasVisibilitytimeout() => $_has(6);
  @$pb.TagNumber(460820073)
  void clearVisibilitytimeout() => $_clearField(460820073);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(7);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(7, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(7);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class ReceiveMessageResult extends $pb.GeneratedMessage {
  factory ReceiveMessageResult({
    $core.Iterable<Message>? messages,
  }) {
    final result = create();
    if (messages != null) result.messages.addAll(messages);
    return result;
  }

  ReceiveMessageResult._();

  factory ReceiveMessageResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReceiveMessageResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReceiveMessageResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPM<Message>(409018326, _omitFieldNames ? '' : 'messages',
        subBuilder: Message.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReceiveMessageResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReceiveMessageResult copyWith(void Function(ReceiveMessageResult) updates) =>
      super.copyWith((message) => updates(message as ReceiveMessageResult))
          as ReceiveMessageResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReceiveMessageResult create() => ReceiveMessageResult._();
  @$core.override
  ReceiveMessageResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReceiveMessageResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReceiveMessageResult>(create);
  static ReceiveMessageResult? _defaultInstance;

  @$pb.TagNumber(409018326)
  $pb.PbList<Message> get messages => $_getList(0);
}

class RemovePermissionRequest extends $pb.GeneratedMessage {
  factory RemovePermissionRequest({
    $core.String? queueurl,
    $core.String? label,
  }) {
    final result = create();
    if (queueurl != null) result.queueurl = queueurl;
    if (label != null) result.label = label;
    return result;
  }

  RemovePermissionRequest._();

  factory RemovePermissionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemovePermissionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemovePermissionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..aOS(516747934, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemovePermissionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemovePermissionRequest copyWith(
          void Function(RemovePermissionRequest) updates) =>
      super.copyWith((message) => updates(message as RemovePermissionRequest))
          as RemovePermissionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemovePermissionRequest create() => RemovePermissionRequest._();
  @$core.override
  RemovePermissionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemovePermissionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemovePermissionRequest>(create);
  static RemovePermissionRequest? _defaultInstance;

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(0);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(0);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);

  @$pb.TagNumber(516747934)
  $core.String get label => $_getSZ(1);
  @$pb.TagNumber(516747934)
  set label($core.String value) => $_setString(1, value);
  @$pb.TagNumber(516747934)
  $core.bool hasLabel() => $_has(1);
  @$pb.TagNumber(516747934)
  void clearLabel() => $_clearField(516747934);
}

class RequestThrottled extends $pb.GeneratedMessage {
  factory RequestThrottled({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  RequestThrottled._();

  factory RequestThrottled.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestThrottled.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestThrottled',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestThrottled clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestThrottled copyWith(void Function(RequestThrottled) updates) =>
      super.copyWith((message) => updates(message as RequestThrottled))
          as RequestThrottled;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestThrottled create() => RequestThrottled._();
  @$core.override
  RequestThrottled createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestThrottled getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestThrottled>(create);
  static RequestThrottled? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class SendMessageBatchRequest extends $pb.GeneratedMessage {
  factory SendMessageBatchRequest({
    $core.Iterable<SendMessageBatchRequestEntry>? entries,
    $core.String? queueurl,
  }) {
    final result = create();
    if (entries != null) result.entries.addAll(entries);
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  SendMessageBatchRequest._();

  factory SendMessageBatchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendMessageBatchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendMessageBatchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPM<SendMessageBatchRequestEntry>(
        481075860, _omitFieldNames ? '' : 'entries',
        subBuilder: SendMessageBatchRequestEntry.create)
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageBatchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageBatchRequest copyWith(
          void Function(SendMessageBatchRequest) updates) =>
      super.copyWith((message) => updates(message as SendMessageBatchRequest))
          as SendMessageBatchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendMessageBatchRequest create() => SendMessageBatchRequest._();
  @$core.override
  SendMessageBatchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendMessageBatchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendMessageBatchRequest>(create);
  static SendMessageBatchRequest? _defaultInstance;

  @$pb.TagNumber(481075860)
  $pb.PbList<SendMessageBatchRequestEntry> get entries => $_getList(0);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(1);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(1);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class SendMessageBatchRequestEntry extends $pb.GeneratedMessage {
  factory SendMessageBatchRequestEntry({
    $core.int? delayseconds,
    $core.Iterable<$core.MapEntry<$core.String, MessageAttributeValue>>?
        messageattributes,
    $core.String? messagebody,
    $core.Iterable<$core.MapEntry<$core.String, MessageSystemAttributeValue>>?
        messagesystemattributes,
    $core.String? messagededuplicationid,
    $core.String? id,
    $core.String? messagegroupid,
  }) {
    final result = create();
    if (delayseconds != null) result.delayseconds = delayseconds;
    if (messageattributes != null)
      result.messageattributes.addEntries(messageattributes);
    if (messagebody != null) result.messagebody = messagebody;
    if (messagesystemattributes != null)
      result.messagesystemattributes.addEntries(messagesystemattributes);
    if (messagededuplicationid != null)
      result.messagededuplicationid = messagededuplicationid;
    if (id != null) result.id = id;
    if (messagegroupid != null) result.messagegroupid = messagegroupid;
    return result;
  }

  SendMessageBatchRequestEntry._();

  factory SendMessageBatchRequestEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendMessageBatchRequestEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendMessageBatchRequestEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aI(48268198, _omitFieldNames ? '' : 'delayseconds')
    ..m<$core.String, MessageAttributeValue>(
        56443766, _omitFieldNames ? '' : 'messageattributes',
        entryClassName: 'SendMessageBatchRequestEntry.MessageattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MessageAttributeValue.create,
        valueDefaultOrMaker: MessageAttributeValue.getDefault,
        packageName: const $pb.PackageName('sqs'))
    ..aOS(56920001, _omitFieldNames ? '' : 'messagebody')
    ..m<$core.String, MessageSystemAttributeValue>(
        194951997, _omitFieldNames ? '' : 'messagesystemattributes',
        entryClassName:
            'SendMessageBatchRequestEntry.MessagesystemattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MessageSystemAttributeValue.create,
        valueDefaultOrMaker: MessageSystemAttributeValue.getDefault,
        packageName: const $pb.PackageName('sqs'))
    ..aOS(379560665, _omitFieldNames ? '' : 'messagededuplicationid')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(419537435, _omitFieldNames ? '' : 'messagegroupid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageBatchRequestEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageBatchRequestEntry copyWith(
          void Function(SendMessageBatchRequestEntry) updates) =>
      super.copyWith(
              (message) => updates(message as SendMessageBatchRequestEntry))
          as SendMessageBatchRequestEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendMessageBatchRequestEntry create() =>
      SendMessageBatchRequestEntry._();
  @$core.override
  SendMessageBatchRequestEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendMessageBatchRequestEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendMessageBatchRequestEntry>(create);
  static SendMessageBatchRequestEntry? _defaultInstance;

  @$pb.TagNumber(48268198)
  $core.int get delayseconds => $_getIZ(0);
  @$pb.TagNumber(48268198)
  set delayseconds($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(48268198)
  $core.bool hasDelayseconds() => $_has(0);
  @$pb.TagNumber(48268198)
  void clearDelayseconds() => $_clearField(48268198);

  @$pb.TagNumber(56443766)
  $pb.PbMap<$core.String, MessageAttributeValue> get messageattributes =>
      $_getMap(1);

  @$pb.TagNumber(56920001)
  $core.String get messagebody => $_getSZ(2);
  @$pb.TagNumber(56920001)
  set messagebody($core.String value) => $_setString(2, value);
  @$pb.TagNumber(56920001)
  $core.bool hasMessagebody() => $_has(2);
  @$pb.TagNumber(56920001)
  void clearMessagebody() => $_clearField(56920001);

  @$pb.TagNumber(194951997)
  $pb.PbMap<$core.String, MessageSystemAttributeValue>
      get messagesystemattributes => $_getMap(3);

  @$pb.TagNumber(379560665)
  $core.String get messagededuplicationid => $_getSZ(4);
  @$pb.TagNumber(379560665)
  set messagededuplicationid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(379560665)
  $core.bool hasMessagededuplicationid() => $_has(4);
  @$pb.TagNumber(379560665)
  void clearMessagededuplicationid() => $_clearField(379560665);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(5);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(5, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(5);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(419537435)
  $core.String get messagegroupid => $_getSZ(6);
  @$pb.TagNumber(419537435)
  set messagegroupid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(419537435)
  $core.bool hasMessagegroupid() => $_has(6);
  @$pb.TagNumber(419537435)
  void clearMessagegroupid() => $_clearField(419537435);
}

class SendMessageBatchResult extends $pb.GeneratedMessage {
  factory SendMessageBatchResult({
    $core.Iterable<BatchResultErrorEntry>? failed,
    $core.Iterable<SendMessageBatchResultEntry>? successful,
  }) {
    final result = create();
    if (failed != null) result.failed.addAll(failed);
    if (successful != null) result.successful.addAll(successful);
    return result;
  }

  SendMessageBatchResult._();

  factory SendMessageBatchResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendMessageBatchResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendMessageBatchResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPM<BatchResultErrorEntry>(360301525, _omitFieldNames ? '' : 'failed',
        subBuilder: BatchResultErrorEntry.create)
    ..pPM<SendMessageBatchResultEntry>(
        412818844, _omitFieldNames ? '' : 'successful',
        subBuilder: SendMessageBatchResultEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageBatchResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageBatchResult copyWith(
          void Function(SendMessageBatchResult) updates) =>
      super.copyWith((message) => updates(message as SendMessageBatchResult))
          as SendMessageBatchResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendMessageBatchResult create() => SendMessageBatchResult._();
  @$core.override
  SendMessageBatchResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendMessageBatchResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendMessageBatchResult>(create);
  static SendMessageBatchResult? _defaultInstance;

  @$pb.TagNumber(360301525)
  $pb.PbList<BatchResultErrorEntry> get failed => $_getList(0);

  @$pb.TagNumber(412818844)
  $pb.PbList<SendMessageBatchResultEntry> get successful => $_getList(1);
}

class SendMessageBatchResultEntry extends $pb.GeneratedMessage {
  factory SendMessageBatchResultEntry({
    $core.String? md5ofmessagebody,
    $core.String? sequencenumber,
    $core.String? md5ofmessagesystemattributes,
    $core.String? messageid,
    $core.String? id,
    $core.String? md5ofmessageattributes,
  }) {
    final result = create();
    if (md5ofmessagebody != null) result.md5ofmessagebody = md5ofmessagebody;
    if (sequencenumber != null) result.sequencenumber = sequencenumber;
    if (md5ofmessagesystemattributes != null)
      result.md5ofmessagesystemattributes = md5ofmessagesystemattributes;
    if (messageid != null) result.messageid = messageid;
    if (id != null) result.id = id;
    if (md5ofmessageattributes != null)
      result.md5ofmessageattributes = md5ofmessageattributes;
    return result;
  }

  SendMessageBatchResultEntry._();

  factory SendMessageBatchResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendMessageBatchResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendMessageBatchResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(28462758, _omitFieldNames ? '' : 'md5ofmessagebody')
    ..aOS(98094362, _omitFieldNames ? '' : 'sequencenumber')
    ..aOS(304512206, _omitFieldNames ? '' : 'md5ofmessagesystemattributes')
    ..aOS(360526634, _omitFieldNames ? '' : 'messageid')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(499647077, _omitFieldNames ? '' : 'md5ofmessageattributes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageBatchResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageBatchResultEntry copyWith(
          void Function(SendMessageBatchResultEntry) updates) =>
      super.copyWith(
              (message) => updates(message as SendMessageBatchResultEntry))
          as SendMessageBatchResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendMessageBatchResultEntry create() =>
      SendMessageBatchResultEntry._();
  @$core.override
  SendMessageBatchResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendMessageBatchResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendMessageBatchResultEntry>(create);
  static SendMessageBatchResultEntry? _defaultInstance;

  @$pb.TagNumber(28462758)
  $core.String get md5ofmessagebody => $_getSZ(0);
  @$pb.TagNumber(28462758)
  set md5ofmessagebody($core.String value) => $_setString(0, value);
  @$pb.TagNumber(28462758)
  $core.bool hasMd5ofmessagebody() => $_has(0);
  @$pb.TagNumber(28462758)
  void clearMd5ofmessagebody() => $_clearField(28462758);

  @$pb.TagNumber(98094362)
  $core.String get sequencenumber => $_getSZ(1);
  @$pb.TagNumber(98094362)
  set sequencenumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(98094362)
  $core.bool hasSequencenumber() => $_has(1);
  @$pb.TagNumber(98094362)
  void clearSequencenumber() => $_clearField(98094362);

  @$pb.TagNumber(304512206)
  $core.String get md5ofmessagesystemattributes => $_getSZ(2);
  @$pb.TagNumber(304512206)
  set md5ofmessagesystemattributes($core.String value) => $_setString(2, value);
  @$pb.TagNumber(304512206)
  $core.bool hasMd5ofmessagesystemattributes() => $_has(2);
  @$pb.TagNumber(304512206)
  void clearMd5ofmessagesystemattributes() => $_clearField(304512206);

  @$pb.TagNumber(360526634)
  $core.String get messageid => $_getSZ(3);
  @$pb.TagNumber(360526634)
  set messageid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(360526634)
  $core.bool hasMessageid() => $_has(3);
  @$pb.TagNumber(360526634)
  void clearMessageid() => $_clearField(360526634);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(499647077)
  $core.String get md5ofmessageattributes => $_getSZ(5);
  @$pb.TagNumber(499647077)
  set md5ofmessageattributes($core.String value) => $_setString(5, value);
  @$pb.TagNumber(499647077)
  $core.bool hasMd5ofmessageattributes() => $_has(5);
  @$pb.TagNumber(499647077)
  void clearMd5ofmessageattributes() => $_clearField(499647077);
}

class SendMessageRequest extends $pb.GeneratedMessage {
  factory SendMessageRequest({
    $core.int? delayseconds,
    $core.Iterable<$core.MapEntry<$core.String, MessageAttributeValue>>?
        messageattributes,
    $core.String? messagebody,
    $core.Iterable<$core.MapEntry<$core.String, MessageSystemAttributeValue>>?
        messagesystemattributes,
    $core.String? messagededuplicationid,
    $core.String? messagegroupid,
    $core.String? queueurl,
  }) {
    final result = create();
    if (delayseconds != null) result.delayseconds = delayseconds;
    if (messageattributes != null)
      result.messageattributes.addEntries(messageattributes);
    if (messagebody != null) result.messagebody = messagebody;
    if (messagesystemattributes != null)
      result.messagesystemattributes.addEntries(messagesystemattributes);
    if (messagededuplicationid != null)
      result.messagededuplicationid = messagededuplicationid;
    if (messagegroupid != null) result.messagegroupid = messagegroupid;
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  SendMessageRequest._();

  factory SendMessageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendMessageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendMessageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aI(48268198, _omitFieldNames ? '' : 'delayseconds')
    ..m<$core.String, MessageAttributeValue>(
        56443766, _omitFieldNames ? '' : 'messageattributes',
        entryClassName: 'SendMessageRequest.MessageattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MessageAttributeValue.create,
        valueDefaultOrMaker: MessageAttributeValue.getDefault,
        packageName: const $pb.PackageName('sqs'))
    ..aOS(56920001, _omitFieldNames ? '' : 'messagebody')
    ..m<$core.String, MessageSystemAttributeValue>(
        194951997, _omitFieldNames ? '' : 'messagesystemattributes',
        entryClassName: 'SendMessageRequest.MessagesystemattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MessageSystemAttributeValue.create,
        valueDefaultOrMaker: MessageSystemAttributeValue.getDefault,
        packageName: const $pb.PackageName('sqs'))
    ..aOS(379560665, _omitFieldNames ? '' : 'messagededuplicationid')
    ..aOS(419537435, _omitFieldNames ? '' : 'messagegroupid')
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageRequest copyWith(void Function(SendMessageRequest) updates) =>
      super.copyWith((message) => updates(message as SendMessageRequest))
          as SendMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendMessageRequest create() => SendMessageRequest._();
  @$core.override
  SendMessageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendMessageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendMessageRequest>(create);
  static SendMessageRequest? _defaultInstance;

  @$pb.TagNumber(48268198)
  $core.int get delayseconds => $_getIZ(0);
  @$pb.TagNumber(48268198)
  set delayseconds($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(48268198)
  $core.bool hasDelayseconds() => $_has(0);
  @$pb.TagNumber(48268198)
  void clearDelayseconds() => $_clearField(48268198);

  @$pb.TagNumber(56443766)
  $pb.PbMap<$core.String, MessageAttributeValue> get messageattributes =>
      $_getMap(1);

  @$pb.TagNumber(56920001)
  $core.String get messagebody => $_getSZ(2);
  @$pb.TagNumber(56920001)
  set messagebody($core.String value) => $_setString(2, value);
  @$pb.TagNumber(56920001)
  $core.bool hasMessagebody() => $_has(2);
  @$pb.TagNumber(56920001)
  void clearMessagebody() => $_clearField(56920001);

  @$pb.TagNumber(194951997)
  $pb.PbMap<$core.String, MessageSystemAttributeValue>
      get messagesystemattributes => $_getMap(3);

  @$pb.TagNumber(379560665)
  $core.String get messagededuplicationid => $_getSZ(4);
  @$pb.TagNumber(379560665)
  set messagededuplicationid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(379560665)
  $core.bool hasMessagededuplicationid() => $_has(4);
  @$pb.TagNumber(379560665)
  void clearMessagededuplicationid() => $_clearField(379560665);

  @$pb.TagNumber(419537435)
  $core.String get messagegroupid => $_getSZ(5);
  @$pb.TagNumber(419537435)
  set messagegroupid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(419537435)
  $core.bool hasMessagegroupid() => $_has(5);
  @$pb.TagNumber(419537435)
  void clearMessagegroupid() => $_clearField(419537435);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(6);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(6, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(6);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class SendMessageResult extends $pb.GeneratedMessage {
  factory SendMessageResult({
    $core.String? md5ofmessagebody,
    $core.String? sequencenumber,
    $core.String? md5ofmessagesystemattributes,
    $core.String? messageid,
    $core.String? md5ofmessageattributes,
  }) {
    final result = create();
    if (md5ofmessagebody != null) result.md5ofmessagebody = md5ofmessagebody;
    if (sequencenumber != null) result.sequencenumber = sequencenumber;
    if (md5ofmessagesystemattributes != null)
      result.md5ofmessagesystemattributes = md5ofmessagesystemattributes;
    if (messageid != null) result.messageid = messageid;
    if (md5ofmessageattributes != null)
      result.md5ofmessageattributes = md5ofmessageattributes;
    return result;
  }

  SendMessageResult._();

  factory SendMessageResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SendMessageResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SendMessageResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(28462758, _omitFieldNames ? '' : 'md5ofmessagebody')
    ..aOS(98094362, _omitFieldNames ? '' : 'sequencenumber')
    ..aOS(304512206, _omitFieldNames ? '' : 'md5ofmessagesystemattributes')
    ..aOS(360526634, _omitFieldNames ? '' : 'messageid')
    ..aOS(499647077, _omitFieldNames ? '' : 'md5ofmessageattributes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SendMessageResult copyWith(void Function(SendMessageResult) updates) =>
      super.copyWith((message) => updates(message as SendMessageResult))
          as SendMessageResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SendMessageResult create() => SendMessageResult._();
  @$core.override
  SendMessageResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SendMessageResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SendMessageResult>(create);
  static SendMessageResult? _defaultInstance;

  @$pb.TagNumber(28462758)
  $core.String get md5ofmessagebody => $_getSZ(0);
  @$pb.TagNumber(28462758)
  set md5ofmessagebody($core.String value) => $_setString(0, value);
  @$pb.TagNumber(28462758)
  $core.bool hasMd5ofmessagebody() => $_has(0);
  @$pb.TagNumber(28462758)
  void clearMd5ofmessagebody() => $_clearField(28462758);

  @$pb.TagNumber(98094362)
  $core.String get sequencenumber => $_getSZ(1);
  @$pb.TagNumber(98094362)
  set sequencenumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(98094362)
  $core.bool hasSequencenumber() => $_has(1);
  @$pb.TagNumber(98094362)
  void clearSequencenumber() => $_clearField(98094362);

  @$pb.TagNumber(304512206)
  $core.String get md5ofmessagesystemattributes => $_getSZ(2);
  @$pb.TagNumber(304512206)
  set md5ofmessagesystemattributes($core.String value) => $_setString(2, value);
  @$pb.TagNumber(304512206)
  $core.bool hasMd5ofmessagesystemattributes() => $_has(2);
  @$pb.TagNumber(304512206)
  void clearMd5ofmessagesystemattributes() => $_clearField(304512206);

  @$pb.TagNumber(360526634)
  $core.String get messageid => $_getSZ(3);
  @$pb.TagNumber(360526634)
  set messageid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(360526634)
  $core.bool hasMessageid() => $_has(3);
  @$pb.TagNumber(360526634)
  void clearMessageid() => $_clearField(360526634);

  @$pb.TagNumber(499647077)
  $core.String get md5ofmessageattributes => $_getSZ(4);
  @$pb.TagNumber(499647077)
  set md5ofmessageattributes($core.String value) => $_setString(4, value);
  @$pb.TagNumber(499647077)
  $core.bool hasMd5ofmessageattributes() => $_has(4);
  @$pb.TagNumber(499647077)
  void clearMd5ofmessageattributes() => $_clearField(499647077);
}

class SetQueueAttributesRequest extends $pb.GeneratedMessage {
  factory SetQueueAttributesRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.String? queueurl,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  SetQueueAttributesRequest._();

  factory SetQueueAttributesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetQueueAttributesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetQueueAttributesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'SetQueueAttributesRequest.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sqs'))
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetQueueAttributesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetQueueAttributesRequest copyWith(
          void Function(SetQueueAttributesRequest) updates) =>
      super.copyWith((message) => updates(message as SetQueueAttributesRequest))
          as SetQueueAttributesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetQueueAttributesRequest create() => SetQueueAttributesRequest._();
  @$core.override
  SetQueueAttributesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetQueueAttributesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetQueueAttributesRequest>(create);
  static SetQueueAttributesRequest? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(1);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(1);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class StartMessageMoveTaskRequest extends $pb.GeneratedMessage {
  factory StartMessageMoveTaskRequest({
    $core.int? maxnumberofmessagespersecond,
    $core.String? destinationarn,
    $core.String? sourcearn,
  }) {
    final result = create();
    if (maxnumberofmessagespersecond != null)
      result.maxnumberofmessagespersecond = maxnumberofmessagespersecond;
    if (destinationarn != null) result.destinationarn = destinationarn;
    if (sourcearn != null) result.sourcearn = sourcearn;
    return result;
  }

  StartMessageMoveTaskRequest._();

  factory StartMessageMoveTaskRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartMessageMoveTaskRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartMessageMoveTaskRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aI(335779921, _omitFieldNames ? '' : 'maxnumberofmessagespersecond')
    ..aOS(375726595, _omitFieldNames ? '' : 'destinationarn')
    ..aOS(439903072, _omitFieldNames ? '' : 'sourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartMessageMoveTaskRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartMessageMoveTaskRequest copyWith(
          void Function(StartMessageMoveTaskRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StartMessageMoveTaskRequest))
          as StartMessageMoveTaskRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartMessageMoveTaskRequest create() =>
      StartMessageMoveTaskRequest._();
  @$core.override
  StartMessageMoveTaskRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartMessageMoveTaskRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartMessageMoveTaskRequest>(create);
  static StartMessageMoveTaskRequest? _defaultInstance;

  @$pb.TagNumber(335779921)
  $core.int get maxnumberofmessagespersecond => $_getIZ(0);
  @$pb.TagNumber(335779921)
  set maxnumberofmessagespersecond($core.int value) =>
      $_setSignedInt32(0, value);
  @$pb.TagNumber(335779921)
  $core.bool hasMaxnumberofmessagespersecond() => $_has(0);
  @$pb.TagNumber(335779921)
  void clearMaxnumberofmessagespersecond() => $_clearField(335779921);

  @$pb.TagNumber(375726595)
  $core.String get destinationarn => $_getSZ(1);
  @$pb.TagNumber(375726595)
  set destinationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(375726595)
  $core.bool hasDestinationarn() => $_has(1);
  @$pb.TagNumber(375726595)
  void clearDestinationarn() => $_clearField(375726595);

  @$pb.TagNumber(439903072)
  $core.String get sourcearn => $_getSZ(2);
  @$pb.TagNumber(439903072)
  set sourcearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(439903072)
  $core.bool hasSourcearn() => $_has(2);
  @$pb.TagNumber(439903072)
  void clearSourcearn() => $_clearField(439903072);
}

class StartMessageMoveTaskResult extends $pb.GeneratedMessage {
  factory StartMessageMoveTaskResult({
    $core.String? taskhandle,
  }) {
    final result = create();
    if (taskhandle != null) result.taskhandle = taskhandle;
    return result;
  }

  StartMessageMoveTaskResult._();

  factory StartMessageMoveTaskResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartMessageMoveTaskResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartMessageMoveTaskResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(190544291, _omitFieldNames ? '' : 'taskhandle')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartMessageMoveTaskResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartMessageMoveTaskResult copyWith(
          void Function(StartMessageMoveTaskResult) updates) =>
      super.copyWith(
              (message) => updates(message as StartMessageMoveTaskResult))
          as StartMessageMoveTaskResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartMessageMoveTaskResult create() => StartMessageMoveTaskResult._();
  @$core.override
  StartMessageMoveTaskResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartMessageMoveTaskResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartMessageMoveTaskResult>(create);
  static StartMessageMoveTaskResult? _defaultInstance;

  @$pb.TagNumber(190544291)
  $core.String get taskhandle => $_getSZ(0);
  @$pb.TagNumber(190544291)
  set taskhandle($core.String value) => $_setString(0, value);
  @$pb.TagNumber(190544291)
  $core.bool hasTaskhandle() => $_has(0);
  @$pb.TagNumber(190544291)
  void clearTaskhandle() => $_clearField(190544291);
}

class TagQueueRequest extends $pb.GeneratedMessage {
  factory TagQueueRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? queueurl,
  }) {
    final result = create();
    if (tags != null) result.tags.addEntries(tags);
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  TagQueueRequest._();

  factory TagQueueRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagQueueRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagQueueRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'TagQueueRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sqs'))
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagQueueRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagQueueRequest copyWith(void Function(TagQueueRequest) updates) =>
      super.copyWith((message) => updates(message as TagQueueRequest))
          as TagQueueRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagQueueRequest create() => TagQueueRequest._();
  @$core.override
  TagQueueRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagQueueRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagQueueRequest>(create);
  static TagQueueRequest? _defaultInstance;

  @$pb.TagNumber(381526209)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(0);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(1);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(1);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

class TooManyEntriesInBatchRequest extends $pb.GeneratedMessage {
  factory TooManyEntriesInBatchRequest({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyEntriesInBatchRequest._();

  factory TooManyEntriesInBatchRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyEntriesInBatchRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyEntriesInBatchRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyEntriesInBatchRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyEntriesInBatchRequest copyWith(
          void Function(TooManyEntriesInBatchRequest) updates) =>
      super.copyWith(
              (message) => updates(message as TooManyEntriesInBatchRequest))
          as TooManyEntriesInBatchRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyEntriesInBatchRequest create() =>
      TooManyEntriesInBatchRequest._();
  @$core.override
  TooManyEntriesInBatchRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyEntriesInBatchRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyEntriesInBatchRequest>(create);
  static TooManyEntriesInBatchRequest? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class UnsupportedOperation extends $pb.GeneratedMessage {
  factory UnsupportedOperation({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  UnsupportedOperation._();

  factory UnsupportedOperation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnsupportedOperation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnsupportedOperation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnsupportedOperation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnsupportedOperation copyWith(void Function(UnsupportedOperation) updates) =>
      super.copyWith((message) => updates(message as UnsupportedOperation))
          as UnsupportedOperation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnsupportedOperation create() => UnsupportedOperation._();
  @$core.override
  UnsupportedOperation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnsupportedOperation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnsupportedOperation>(create);
  static UnsupportedOperation? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class UntagQueueRequest extends $pb.GeneratedMessage {
  factory UntagQueueRequest({
    $core.Iterable<$core.String>? tagkeys,
    $core.String? queueurl,
  }) {
    final result = create();
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
    if (queueurl != null) result.queueurl = queueurl;
    return result;
  }

  UntagQueueRequest._();

  factory UntagQueueRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UntagQueueRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UntagQueueRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sqs'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(510632138, _omitFieldNames ? '' : 'queueurl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagQueueRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagQueueRequest copyWith(void Function(UntagQueueRequest) updates) =>
      super.copyWith((message) => updates(message as UntagQueueRequest))
          as UntagQueueRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UntagQueueRequest create() => UntagQueueRequest._();
  @$core.override
  UntagQueueRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UntagQueueRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UntagQueueRequest>(create);
  static UntagQueueRequest? _defaultInstance;

  @$pb.TagNumber(320659964)
  $pb.PbList<$core.String> get tagkeys => $_getList(0);

  @$pb.TagNumber(510632138)
  $core.String get queueurl => $_getSZ(1);
  @$pb.TagNumber(510632138)
  set queueurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(510632138)
  $core.bool hasQueueurl() => $_has(1);
  @$pb.TagNumber(510632138)
  void clearQueueurl() => $_clearField(510632138);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
