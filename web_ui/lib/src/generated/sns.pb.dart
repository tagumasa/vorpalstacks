// This is a generated file - do not edit.
//
// Generated from sns.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'sns.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'sns.pbenum.dart';

class AddPermissionInput extends $pb.GeneratedMessage {
  factory AddPermissionInput({
    $core.String? topicarn,
    $core.Iterable<$core.String>? actionname,
    $core.Iterable<$core.String>? awsaccountid,
    $core.String? label,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    if (actionname != null) result.actionname.addAll(actionname);
    if (awsaccountid != null) result.awsaccountid.addAll(awsaccountid);
    if (label != null) result.label = label;
    return result;
  }

  AddPermissionInput._();

  factory AddPermissionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AddPermissionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddPermissionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..pPS(115541053, _omitFieldNames ? '' : 'actionname')
    ..pPS(370093421, _omitFieldNames ? '' : 'awsaccountid')
    ..aOS(516747934, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddPermissionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddPermissionInput copyWith(void Function(AddPermissionInput) updates) =>
      super.copyWith((message) => updates(message as AddPermissionInput))
          as AddPermissionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddPermissionInput create() => AddPermissionInput._();
  @$core.override
  AddPermissionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AddPermissionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddPermissionInput>(create);
  static AddPermissionInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(115541053)
  $pb.PbList<$core.String> get actionname => $_getList(1);

  @$pb.TagNumber(370093421)
  $pb.PbList<$core.String> get awsaccountid => $_getList(2);

  @$pb.TagNumber(516747934)
  $core.String get label => $_getSZ(3);
  @$pb.TagNumber(516747934)
  set label($core.String value) => $_setString(3, value);
  @$pb.TagNumber(516747934)
  $core.bool hasLabel() => $_has(3);
  @$pb.TagNumber(516747934)
  void clearLabel() => $_clearField(516747934);
}

class AuthorizationErrorException extends $pb.GeneratedMessage {
  factory AuthorizationErrorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  AuthorizationErrorException._();

  factory AuthorizationErrorException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AuthorizationErrorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AuthorizationErrorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthorizationErrorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AuthorizationErrorException copyWith(
          void Function(AuthorizationErrorException) updates) =>
      super.copyWith(
              (message) => updates(message as AuthorizationErrorException))
          as AuthorizationErrorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AuthorizationErrorException create() =>
      AuthorizationErrorException._();
  @$core.override
  AuthorizationErrorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AuthorizationErrorException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AuthorizationErrorException>(create);
  static AuthorizationErrorException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class BatchEntryIdsNotDistinctException extends $pb.GeneratedMessage {
  factory BatchEntryIdsNotDistinctException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  BatchEntryIdsNotDistinctException._();

  factory BatchEntryIdsNotDistinctException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchEntryIdsNotDistinctException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchEntryIdsNotDistinctException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchEntryIdsNotDistinctException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchEntryIdsNotDistinctException copyWith(
          void Function(BatchEntryIdsNotDistinctException) updates) =>
      super.copyWith((message) =>
              updates(message as BatchEntryIdsNotDistinctException))
          as BatchEntryIdsNotDistinctException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchEntryIdsNotDistinctException create() =>
      BatchEntryIdsNotDistinctException._();
  @$core.override
  BatchEntryIdsNotDistinctException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchEntryIdsNotDistinctException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchEntryIdsNotDistinctException>(
          create);
  static BatchEntryIdsNotDistinctException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class BatchRequestTooLongException extends $pb.GeneratedMessage {
  factory BatchRequestTooLongException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  BatchRequestTooLongException._();

  factory BatchRequestTooLongException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchRequestTooLongException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchRequestTooLongException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchRequestTooLongException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchRequestTooLongException copyWith(
          void Function(BatchRequestTooLongException) updates) =>
      super.copyWith(
              (message) => updates(message as BatchRequestTooLongException))
          as BatchRequestTooLongException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchRequestTooLongException create() =>
      BatchRequestTooLongException._();
  @$core.override
  BatchRequestTooLongException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchRequestTooLongException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchRequestTooLongException>(create);
  static BatchRequestTooLongException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
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

class CheckIfPhoneNumberIsOptedOutInput extends $pb.GeneratedMessage {
  factory CheckIfPhoneNumberIsOptedOutInput({
    $core.String? phonenumber,
  }) {
    final result = create();
    if (phonenumber != null) result.phonenumber = phonenumber;
    return result;
  }

  CheckIfPhoneNumberIsOptedOutInput._();

  factory CheckIfPhoneNumberIsOptedOutInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CheckIfPhoneNumberIsOptedOutInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CheckIfPhoneNumberIsOptedOutInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(156377999, _omitFieldNames ? '' : 'phonenumber')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckIfPhoneNumberIsOptedOutInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckIfPhoneNumberIsOptedOutInput copyWith(
          void Function(CheckIfPhoneNumberIsOptedOutInput) updates) =>
      super.copyWith((message) =>
              updates(message as CheckIfPhoneNumberIsOptedOutInput))
          as CheckIfPhoneNumberIsOptedOutInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CheckIfPhoneNumberIsOptedOutInput create() =>
      CheckIfPhoneNumberIsOptedOutInput._();
  @$core.override
  CheckIfPhoneNumberIsOptedOutInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CheckIfPhoneNumberIsOptedOutInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CheckIfPhoneNumberIsOptedOutInput>(
          create);
  static CheckIfPhoneNumberIsOptedOutInput? _defaultInstance;

  @$pb.TagNumber(156377999)
  $core.String get phonenumber => $_getSZ(0);
  @$pb.TagNumber(156377999)
  set phonenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(156377999)
  $core.bool hasPhonenumber() => $_has(0);
  @$pb.TagNumber(156377999)
  void clearPhonenumber() => $_clearField(156377999);
}

class CheckIfPhoneNumberIsOptedOutResponse extends $pb.GeneratedMessage {
  factory CheckIfPhoneNumberIsOptedOutResponse({
    $core.bool? isoptedout,
  }) {
    final result = create();
    if (isoptedout != null) result.isoptedout = isoptedout;
    return result;
  }

  CheckIfPhoneNumberIsOptedOutResponse._();

  factory CheckIfPhoneNumberIsOptedOutResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CheckIfPhoneNumberIsOptedOutResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CheckIfPhoneNumberIsOptedOutResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOB(430412750, _omitFieldNames ? '' : 'isoptedout')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckIfPhoneNumberIsOptedOutResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckIfPhoneNumberIsOptedOutResponse copyWith(
          void Function(CheckIfPhoneNumberIsOptedOutResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CheckIfPhoneNumberIsOptedOutResponse))
          as CheckIfPhoneNumberIsOptedOutResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CheckIfPhoneNumberIsOptedOutResponse create() =>
      CheckIfPhoneNumberIsOptedOutResponse._();
  @$core.override
  CheckIfPhoneNumberIsOptedOutResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CheckIfPhoneNumberIsOptedOutResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CheckIfPhoneNumberIsOptedOutResponse>(create);
  static CheckIfPhoneNumberIsOptedOutResponse? _defaultInstance;

  @$pb.TagNumber(430412750)
  $core.bool get isoptedout => $_getBF(0);
  @$pb.TagNumber(430412750)
  set isoptedout($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(430412750)
  $core.bool hasIsoptedout() => $_has(0);
  @$pb.TagNumber(430412750)
  void clearIsoptedout() => $_clearField(430412750);
}

class ConcurrentAccessException extends $pb.GeneratedMessage {
  factory ConcurrentAccessException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ConcurrentAccessException._();

  factory ConcurrentAccessException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConcurrentAccessException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConcurrentAccessException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConcurrentAccessException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConcurrentAccessException copyWith(
          void Function(ConcurrentAccessException) updates) =>
      super.copyWith((message) => updates(message as ConcurrentAccessException))
          as ConcurrentAccessException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConcurrentAccessException create() => ConcurrentAccessException._();
  @$core.override
  ConcurrentAccessException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConcurrentAccessException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConcurrentAccessException>(create);
  static ConcurrentAccessException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ConfirmSubscriptionInput extends $pb.GeneratedMessage {
  factory ConfirmSubscriptionInput({
    $core.String? topicarn,
    $core.String? token,
    $core.String? authenticateonunsubscribe,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    if (token != null) result.token = token;
    if (authenticateonunsubscribe != null)
      result.authenticateonunsubscribe = authenticateonunsubscribe;
    return result;
  }

  ConfirmSubscriptionInput._();

  factory ConfirmSubscriptionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConfirmSubscriptionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConfirmSubscriptionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..aOS(439704531, _omitFieldNames ? '' : 'token')
    ..aOS(529251911, _omitFieldNames ? '' : 'authenticateonunsubscribe')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConfirmSubscriptionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConfirmSubscriptionInput copyWith(
          void Function(ConfirmSubscriptionInput) updates) =>
      super.copyWith((message) => updates(message as ConfirmSubscriptionInput))
          as ConfirmSubscriptionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConfirmSubscriptionInput create() => ConfirmSubscriptionInput._();
  @$core.override
  ConfirmSubscriptionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConfirmSubscriptionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConfirmSubscriptionInput>(create);
  static ConfirmSubscriptionInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(439704531)
  $core.String get token => $_getSZ(1);
  @$pb.TagNumber(439704531)
  set token($core.String value) => $_setString(1, value);
  @$pb.TagNumber(439704531)
  $core.bool hasToken() => $_has(1);
  @$pb.TagNumber(439704531)
  void clearToken() => $_clearField(439704531);

  @$pb.TagNumber(529251911)
  $core.String get authenticateonunsubscribe => $_getSZ(2);
  @$pb.TagNumber(529251911)
  set authenticateonunsubscribe($core.String value) => $_setString(2, value);
  @$pb.TagNumber(529251911)
  $core.bool hasAuthenticateonunsubscribe() => $_has(2);
  @$pb.TagNumber(529251911)
  void clearAuthenticateonunsubscribe() => $_clearField(529251911);
}

class ConfirmSubscriptionResponse extends $pb.GeneratedMessage {
  factory ConfirmSubscriptionResponse({
    $core.String? subscriptionarn,
  }) {
    final result = create();
    if (subscriptionarn != null) result.subscriptionarn = subscriptionarn;
    return result;
  }

  ConfirmSubscriptionResponse._();

  factory ConfirmSubscriptionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConfirmSubscriptionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConfirmSubscriptionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(279547820, _omitFieldNames ? '' : 'subscriptionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConfirmSubscriptionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConfirmSubscriptionResponse copyWith(
          void Function(ConfirmSubscriptionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ConfirmSubscriptionResponse))
          as ConfirmSubscriptionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConfirmSubscriptionResponse create() =>
      ConfirmSubscriptionResponse._();
  @$core.override
  ConfirmSubscriptionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConfirmSubscriptionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConfirmSubscriptionResponse>(create);
  static ConfirmSubscriptionResponse? _defaultInstance;

  @$pb.TagNumber(279547820)
  $core.String get subscriptionarn => $_getSZ(0);
  @$pb.TagNumber(279547820)
  set subscriptionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(279547820)
  $core.bool hasSubscriptionarn() => $_has(0);
  @$pb.TagNumber(279547820)
  void clearSubscriptionarn() => $_clearField(279547820);
}

class CreateEndpointResponse extends $pb.GeneratedMessage {
  factory CreateEndpointResponse({
    $core.String? endpointarn,
  }) {
    final result = create();
    if (endpointarn != null) result.endpointarn = endpointarn;
    return result;
  }

  CreateEndpointResponse._();

  factory CreateEndpointResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateEndpointResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateEndpointResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(32228660, _omitFieldNames ? '' : 'endpointarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEndpointResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEndpointResponse copyWith(
          void Function(CreateEndpointResponse) updates) =>
      super.copyWith((message) => updates(message as CreateEndpointResponse))
          as CreateEndpointResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateEndpointResponse create() => CreateEndpointResponse._();
  @$core.override
  CreateEndpointResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateEndpointResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateEndpointResponse>(create);
  static CreateEndpointResponse? _defaultInstance;

  @$pb.TagNumber(32228660)
  $core.String get endpointarn => $_getSZ(0);
  @$pb.TagNumber(32228660)
  set endpointarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(32228660)
  $core.bool hasEndpointarn() => $_has(0);
  @$pb.TagNumber(32228660)
  void clearEndpointarn() => $_clearField(32228660);
}

class CreatePlatformApplicationInput extends $pb.GeneratedMessage {
  factory CreatePlatformApplicationInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.String? name,
    $core.String? platform,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    if (name != null) result.name = name;
    if (platform != null) result.platform = platform;
    return result;
  }

  CreatePlatformApplicationInput._();

  factory CreatePlatformApplicationInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePlatformApplicationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePlatformApplicationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'CreatePlatformApplicationInput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(468905683, _omitFieldNames ? '' : 'platform')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePlatformApplicationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePlatformApplicationInput copyWith(
          void Function(CreatePlatformApplicationInput) updates) =>
      super.copyWith(
              (message) => updates(message as CreatePlatformApplicationInput))
          as CreatePlatformApplicationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePlatformApplicationInput create() =>
      CreatePlatformApplicationInput._();
  @$core.override
  CreatePlatformApplicationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePlatformApplicationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePlatformApplicationInput>(create);
  static CreatePlatformApplicationInput? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(468905683)
  $core.String get platform => $_getSZ(2);
  @$pb.TagNumber(468905683)
  set platform($core.String value) => $_setString(2, value);
  @$pb.TagNumber(468905683)
  $core.bool hasPlatform() => $_has(2);
  @$pb.TagNumber(468905683)
  void clearPlatform() => $_clearField(468905683);
}

class CreatePlatformApplicationResponse extends $pb.GeneratedMessage {
  factory CreatePlatformApplicationResponse({
    $core.String? platformapplicationarn,
  }) {
    final result = create();
    if (platformapplicationarn != null)
      result.platformapplicationarn = platformapplicationarn;
    return result;
  }

  CreatePlatformApplicationResponse._();

  factory CreatePlatformApplicationResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePlatformApplicationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePlatformApplicationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(241250568, _omitFieldNames ? '' : 'platformapplicationarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePlatformApplicationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePlatformApplicationResponse copyWith(
          void Function(CreatePlatformApplicationResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreatePlatformApplicationResponse))
          as CreatePlatformApplicationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePlatformApplicationResponse create() =>
      CreatePlatformApplicationResponse._();
  @$core.override
  CreatePlatformApplicationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePlatformApplicationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePlatformApplicationResponse>(
          create);
  static CreatePlatformApplicationResponse? _defaultInstance;

  @$pb.TagNumber(241250568)
  $core.String get platformapplicationarn => $_getSZ(0);
  @$pb.TagNumber(241250568)
  set platformapplicationarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(241250568)
  $core.bool hasPlatformapplicationarn() => $_has(0);
  @$pb.TagNumber(241250568)
  void clearPlatformapplicationarn() => $_clearField(241250568);
}

class CreatePlatformEndpointInput extends $pb.GeneratedMessage {
  factory CreatePlatformEndpointInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.String? platformapplicationarn,
    $core.String? customuserdata,
    $core.String? token,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    if (platformapplicationarn != null)
      result.platformapplicationarn = platformapplicationarn;
    if (customuserdata != null) result.customuserdata = customuserdata;
    if (token != null) result.token = token;
    return result;
  }

  CreatePlatformEndpointInput._();

  factory CreatePlatformEndpointInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePlatformEndpointInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePlatformEndpointInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'CreatePlatformEndpointInput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..aOS(241250568, _omitFieldNames ? '' : 'platformapplicationarn')
    ..aOS(388962962, _omitFieldNames ? '' : 'customuserdata')
    ..aOS(439704531, _omitFieldNames ? '' : 'token')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePlatformEndpointInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePlatformEndpointInput copyWith(
          void Function(CreatePlatformEndpointInput) updates) =>
      super.copyWith(
              (message) => updates(message as CreatePlatformEndpointInput))
          as CreatePlatformEndpointInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePlatformEndpointInput create() =>
      CreatePlatformEndpointInput._();
  @$core.override
  CreatePlatformEndpointInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePlatformEndpointInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePlatformEndpointInput>(create);
  static CreatePlatformEndpointInput? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);

  @$pb.TagNumber(241250568)
  $core.String get platformapplicationarn => $_getSZ(1);
  @$pb.TagNumber(241250568)
  set platformapplicationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(241250568)
  $core.bool hasPlatformapplicationarn() => $_has(1);
  @$pb.TagNumber(241250568)
  void clearPlatformapplicationarn() => $_clearField(241250568);

  @$pb.TagNumber(388962962)
  $core.String get customuserdata => $_getSZ(2);
  @$pb.TagNumber(388962962)
  set customuserdata($core.String value) => $_setString(2, value);
  @$pb.TagNumber(388962962)
  $core.bool hasCustomuserdata() => $_has(2);
  @$pb.TagNumber(388962962)
  void clearCustomuserdata() => $_clearField(388962962);

  @$pb.TagNumber(439704531)
  $core.String get token => $_getSZ(3);
  @$pb.TagNumber(439704531)
  set token($core.String value) => $_setString(3, value);
  @$pb.TagNumber(439704531)
  $core.bool hasToken() => $_has(3);
  @$pb.TagNumber(439704531)
  void clearToken() => $_clearField(439704531);
}

class CreateSMSSandboxPhoneNumberInput extends $pb.GeneratedMessage {
  factory CreateSMSSandboxPhoneNumberInput({
    LanguageCodeString? languagecode,
    $core.String? phonenumber,
  }) {
    final result = create();
    if (languagecode != null) result.languagecode = languagecode;
    if (phonenumber != null) result.phonenumber = phonenumber;
    return result;
  }

  CreateSMSSandboxPhoneNumberInput._();

  factory CreateSMSSandboxPhoneNumberInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateSMSSandboxPhoneNumberInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateSMSSandboxPhoneNumberInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aE<LanguageCodeString>(281903107, _omitFieldNames ? '' : 'languagecode',
        enumValues: LanguageCodeString.values)
    ..aOS(379600239, _omitFieldNames ? '' : 'phonenumber')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSMSSandboxPhoneNumberInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSMSSandboxPhoneNumberInput copyWith(
          void Function(CreateSMSSandboxPhoneNumberInput) updates) =>
      super.copyWith(
              (message) => updates(message as CreateSMSSandboxPhoneNumberInput))
          as CreateSMSSandboxPhoneNumberInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateSMSSandboxPhoneNumberInput create() =>
      CreateSMSSandboxPhoneNumberInput._();
  @$core.override
  CreateSMSSandboxPhoneNumberInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateSMSSandboxPhoneNumberInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateSMSSandboxPhoneNumberInput>(
          create);
  static CreateSMSSandboxPhoneNumberInput? _defaultInstance;

  @$pb.TagNumber(281903107)
  LanguageCodeString get languagecode => $_getN(0);
  @$pb.TagNumber(281903107)
  set languagecode(LanguageCodeString value) => $_setField(281903107, value);
  @$pb.TagNumber(281903107)
  $core.bool hasLanguagecode() => $_has(0);
  @$pb.TagNumber(281903107)
  void clearLanguagecode() => $_clearField(281903107);

  @$pb.TagNumber(379600239)
  $core.String get phonenumber => $_getSZ(1);
  @$pb.TagNumber(379600239)
  set phonenumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(379600239)
  $core.bool hasPhonenumber() => $_has(1);
  @$pb.TagNumber(379600239)
  void clearPhonenumber() => $_clearField(379600239);
}

class CreateSMSSandboxPhoneNumberResult extends $pb.GeneratedMessage {
  factory CreateSMSSandboxPhoneNumberResult() => create();

  CreateSMSSandboxPhoneNumberResult._();

  factory CreateSMSSandboxPhoneNumberResult.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateSMSSandboxPhoneNumberResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateSMSSandboxPhoneNumberResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSMSSandboxPhoneNumberResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSMSSandboxPhoneNumberResult copyWith(
          void Function(CreateSMSSandboxPhoneNumberResult) updates) =>
      super.copyWith((message) =>
              updates(message as CreateSMSSandboxPhoneNumberResult))
          as CreateSMSSandboxPhoneNumberResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateSMSSandboxPhoneNumberResult create() =>
      CreateSMSSandboxPhoneNumberResult._();
  @$core.override
  CreateSMSSandboxPhoneNumberResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateSMSSandboxPhoneNumberResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateSMSSandboxPhoneNumberResult>(
          create);
  static CreateSMSSandboxPhoneNumberResult? _defaultInstance;
}

class CreateTopicInput extends $pb.GeneratedMessage {
  factory CreateTopicInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.String? name,
    $core.Iterable<Tag>? tags,
    $core.String? dataprotectionpolicy,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    if (dataprotectionpolicy != null)
      result.dataprotectionpolicy = dataprotectionpolicy;
    return result;
  }

  CreateTopicInput._();

  factory CreateTopicInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTopicInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTopicInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'CreateTopicInput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(519444573, _omitFieldNames ? '' : 'dataprotectionpolicy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTopicInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTopicInput copyWith(void Function(CreateTopicInput) updates) =>
      super.copyWith((message) => updates(message as CreateTopicInput))
          as CreateTopicInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTopicInput create() => CreateTopicInput._();
  @$core.override
  CreateTopicInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTopicInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTopicInput>(create);
  static CreateTopicInput? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(2);

  @$pb.TagNumber(519444573)
  $core.String get dataprotectionpolicy => $_getSZ(3);
  @$pb.TagNumber(519444573)
  set dataprotectionpolicy($core.String value) => $_setString(3, value);
  @$pb.TagNumber(519444573)
  $core.bool hasDataprotectionpolicy() => $_has(3);
  @$pb.TagNumber(519444573)
  void clearDataprotectionpolicy() => $_clearField(519444573);
}

class CreateTopicResponse extends $pb.GeneratedMessage {
  factory CreateTopicResponse({
    $core.String? topicarn,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    return result;
  }

  CreateTopicResponse._();

  factory CreateTopicResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTopicResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTopicResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTopicResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTopicResponse copyWith(void Function(CreateTopicResponse) updates) =>
      super.copyWith((message) => updates(message as CreateTopicResponse))
          as CreateTopicResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTopicResponse create() => CreateTopicResponse._();
  @$core.override
  CreateTopicResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTopicResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTopicResponse>(create);
  static CreateTopicResponse? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);
}

class DeleteEndpointInput extends $pb.GeneratedMessage {
  factory DeleteEndpointInput({
    $core.String? endpointarn,
  }) {
    final result = create();
    if (endpointarn != null) result.endpointarn = endpointarn;
    return result;
  }

  DeleteEndpointInput._();

  factory DeleteEndpointInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteEndpointInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteEndpointInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(32228660, _omitFieldNames ? '' : 'endpointarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteEndpointInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteEndpointInput copyWith(void Function(DeleteEndpointInput) updates) =>
      super.copyWith((message) => updates(message as DeleteEndpointInput))
          as DeleteEndpointInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteEndpointInput create() => DeleteEndpointInput._();
  @$core.override
  DeleteEndpointInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteEndpointInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteEndpointInput>(create);
  static DeleteEndpointInput? _defaultInstance;

  @$pb.TagNumber(32228660)
  $core.String get endpointarn => $_getSZ(0);
  @$pb.TagNumber(32228660)
  set endpointarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(32228660)
  $core.bool hasEndpointarn() => $_has(0);
  @$pb.TagNumber(32228660)
  void clearEndpointarn() => $_clearField(32228660);
}

class DeletePlatformApplicationInput extends $pb.GeneratedMessage {
  factory DeletePlatformApplicationInput({
    $core.String? platformapplicationarn,
  }) {
    final result = create();
    if (platformapplicationarn != null)
      result.platformapplicationarn = platformapplicationarn;
    return result;
  }

  DeletePlatformApplicationInput._();

  factory DeletePlatformApplicationInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeletePlatformApplicationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeletePlatformApplicationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(241250568, _omitFieldNames ? '' : 'platformapplicationarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePlatformApplicationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePlatformApplicationInput copyWith(
          void Function(DeletePlatformApplicationInput) updates) =>
      super.copyWith(
              (message) => updates(message as DeletePlatformApplicationInput))
          as DeletePlatformApplicationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeletePlatformApplicationInput create() =>
      DeletePlatformApplicationInput._();
  @$core.override
  DeletePlatformApplicationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeletePlatformApplicationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeletePlatformApplicationInput>(create);
  static DeletePlatformApplicationInput? _defaultInstance;

  @$pb.TagNumber(241250568)
  $core.String get platformapplicationarn => $_getSZ(0);
  @$pb.TagNumber(241250568)
  set platformapplicationarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(241250568)
  $core.bool hasPlatformapplicationarn() => $_has(0);
  @$pb.TagNumber(241250568)
  void clearPlatformapplicationarn() => $_clearField(241250568);
}

class DeleteSMSSandboxPhoneNumberInput extends $pb.GeneratedMessage {
  factory DeleteSMSSandboxPhoneNumberInput({
    $core.String? phonenumber,
  }) {
    final result = create();
    if (phonenumber != null) result.phonenumber = phonenumber;
    return result;
  }

  DeleteSMSSandboxPhoneNumberInput._();

  factory DeleteSMSSandboxPhoneNumberInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteSMSSandboxPhoneNumberInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteSMSSandboxPhoneNumberInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(379600239, _omitFieldNames ? '' : 'phonenumber')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSMSSandboxPhoneNumberInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSMSSandboxPhoneNumberInput copyWith(
          void Function(DeleteSMSSandboxPhoneNumberInput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteSMSSandboxPhoneNumberInput))
          as DeleteSMSSandboxPhoneNumberInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteSMSSandboxPhoneNumberInput create() =>
      DeleteSMSSandboxPhoneNumberInput._();
  @$core.override
  DeleteSMSSandboxPhoneNumberInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteSMSSandboxPhoneNumberInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteSMSSandboxPhoneNumberInput>(
          create);
  static DeleteSMSSandboxPhoneNumberInput? _defaultInstance;

  @$pb.TagNumber(379600239)
  $core.String get phonenumber => $_getSZ(0);
  @$pb.TagNumber(379600239)
  set phonenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(379600239)
  $core.bool hasPhonenumber() => $_has(0);
  @$pb.TagNumber(379600239)
  void clearPhonenumber() => $_clearField(379600239);
}

class DeleteSMSSandboxPhoneNumberResult extends $pb.GeneratedMessage {
  factory DeleteSMSSandboxPhoneNumberResult() => create();

  DeleteSMSSandboxPhoneNumberResult._();

  factory DeleteSMSSandboxPhoneNumberResult.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteSMSSandboxPhoneNumberResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteSMSSandboxPhoneNumberResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSMSSandboxPhoneNumberResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSMSSandboxPhoneNumberResult copyWith(
          void Function(DeleteSMSSandboxPhoneNumberResult) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteSMSSandboxPhoneNumberResult))
          as DeleteSMSSandboxPhoneNumberResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteSMSSandboxPhoneNumberResult create() =>
      DeleteSMSSandboxPhoneNumberResult._();
  @$core.override
  DeleteSMSSandboxPhoneNumberResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteSMSSandboxPhoneNumberResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteSMSSandboxPhoneNumberResult>(
          create);
  static DeleteSMSSandboxPhoneNumberResult? _defaultInstance;
}

class DeleteTopicInput extends $pb.GeneratedMessage {
  factory DeleteTopicInput({
    $core.String? topicarn,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    return result;
  }

  DeleteTopicInput._();

  factory DeleteTopicInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTopicInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTopicInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTopicInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTopicInput copyWith(void Function(DeleteTopicInput) updates) =>
      super.copyWith((message) => updates(message as DeleteTopicInput))
          as DeleteTopicInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTopicInput create() => DeleteTopicInput._();
  @$core.override
  DeleteTopicInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTopicInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTopicInput>(create);
  static DeleteTopicInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);
}

class EmptyBatchRequestException extends $pb.GeneratedMessage {
  factory EmptyBatchRequestException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EmptyBatchRequestException._();

  factory EmptyBatchRequestException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EmptyBatchRequestException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EmptyBatchRequestException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EmptyBatchRequestException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EmptyBatchRequestException copyWith(
          void Function(EmptyBatchRequestException) updates) =>
      super.copyWith(
              (message) => updates(message as EmptyBatchRequestException))
          as EmptyBatchRequestException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EmptyBatchRequestException create() => EmptyBatchRequestException._();
  @$core.override
  EmptyBatchRequestException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EmptyBatchRequestException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EmptyBatchRequestException>(create);
  static EmptyBatchRequestException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class Endpoint extends $pb.GeneratedMessage {
  factory Endpoint({
    $core.String? endpointarn,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (endpointarn != null) result.endpointarn = endpointarn;
    if (attributes != null) result.attributes.addEntries(attributes);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(32228660, _omitFieldNames ? '' : 'endpointarn')
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'Endpoint.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
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

  @$pb.TagNumber(32228660)
  $core.String get endpointarn => $_getSZ(0);
  @$pb.TagNumber(32228660)
  set endpointarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(32228660)
  $core.bool hasEndpointarn() => $_has(0);
  @$pb.TagNumber(32228660)
  void clearEndpointarn() => $_clearField(32228660);

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(1);
}

class EndpointDisabledException extends $pb.GeneratedMessage {
  factory EndpointDisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EndpointDisabledException._();

  factory EndpointDisabledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EndpointDisabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EndpointDisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EndpointDisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EndpointDisabledException copyWith(
          void Function(EndpointDisabledException) updates) =>
      super.copyWith((message) => updates(message as EndpointDisabledException))
          as EndpointDisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EndpointDisabledException create() => EndpointDisabledException._();
  @$core.override
  EndpointDisabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EndpointDisabledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EndpointDisabledException>(create);
  static EndpointDisabledException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class FilterPolicyLimitExceededException extends $pb.GeneratedMessage {
  factory FilterPolicyLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  FilterPolicyLimitExceededException._();

  factory FilterPolicyLimitExceededException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilterPolicyLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilterPolicyLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilterPolicyLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilterPolicyLimitExceededException copyWith(
          void Function(FilterPolicyLimitExceededException) updates) =>
      super.copyWith((message) =>
              updates(message as FilterPolicyLimitExceededException))
          as FilterPolicyLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilterPolicyLimitExceededException create() =>
      FilterPolicyLimitExceededException._();
  @$core.override
  FilterPolicyLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilterPolicyLimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilterPolicyLimitExceededException>(
          create);
  static FilterPolicyLimitExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class GetDataProtectionPolicyInput extends $pb.GeneratedMessage {
  factory GetDataProtectionPolicyInput({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  GetDataProtectionPolicyInput._();

  factory GetDataProtectionPolicyInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDataProtectionPolicyInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDataProtectionPolicyInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDataProtectionPolicyInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDataProtectionPolicyInput copyWith(
          void Function(GetDataProtectionPolicyInput) updates) =>
      super.copyWith(
              (message) => updates(message as GetDataProtectionPolicyInput))
          as GetDataProtectionPolicyInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDataProtectionPolicyInput create() =>
      GetDataProtectionPolicyInput._();
  @$core.override
  GetDataProtectionPolicyInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDataProtectionPolicyInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDataProtectionPolicyInput>(create);
  static GetDataProtectionPolicyInput? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class GetDataProtectionPolicyResponse extends $pb.GeneratedMessage {
  factory GetDataProtectionPolicyResponse({
    $core.String? dataprotectionpolicy,
  }) {
    final result = create();
    if (dataprotectionpolicy != null)
      result.dataprotectionpolicy = dataprotectionpolicy;
    return result;
  }

  GetDataProtectionPolicyResponse._();

  factory GetDataProtectionPolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDataProtectionPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDataProtectionPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(519444573, _omitFieldNames ? '' : 'dataprotectionpolicy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDataProtectionPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDataProtectionPolicyResponse copyWith(
          void Function(GetDataProtectionPolicyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetDataProtectionPolicyResponse))
          as GetDataProtectionPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDataProtectionPolicyResponse create() =>
      GetDataProtectionPolicyResponse._();
  @$core.override
  GetDataProtectionPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDataProtectionPolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDataProtectionPolicyResponse>(
          create);
  static GetDataProtectionPolicyResponse? _defaultInstance;

  @$pb.TagNumber(519444573)
  $core.String get dataprotectionpolicy => $_getSZ(0);
  @$pb.TagNumber(519444573)
  set dataprotectionpolicy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(519444573)
  $core.bool hasDataprotectionpolicy() => $_has(0);
  @$pb.TagNumber(519444573)
  void clearDataprotectionpolicy() => $_clearField(519444573);
}

class GetEndpointAttributesInput extends $pb.GeneratedMessage {
  factory GetEndpointAttributesInput({
    $core.String? endpointarn,
  }) {
    final result = create();
    if (endpointarn != null) result.endpointarn = endpointarn;
    return result;
  }

  GetEndpointAttributesInput._();

  factory GetEndpointAttributesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetEndpointAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetEndpointAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(32228660, _omitFieldNames ? '' : 'endpointarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEndpointAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEndpointAttributesInput copyWith(
          void Function(GetEndpointAttributesInput) updates) =>
      super.copyWith(
              (message) => updates(message as GetEndpointAttributesInput))
          as GetEndpointAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetEndpointAttributesInput create() => GetEndpointAttributesInput._();
  @$core.override
  GetEndpointAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetEndpointAttributesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetEndpointAttributesInput>(create);
  static GetEndpointAttributesInput? _defaultInstance;

  @$pb.TagNumber(32228660)
  $core.String get endpointarn => $_getSZ(0);
  @$pb.TagNumber(32228660)
  set endpointarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(32228660)
  $core.bool hasEndpointarn() => $_has(0);
  @$pb.TagNumber(32228660)
  void clearEndpointarn() => $_clearField(32228660);
}

class GetEndpointAttributesResponse extends $pb.GeneratedMessage {
  factory GetEndpointAttributesResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  GetEndpointAttributesResponse._();

  factory GetEndpointAttributesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetEndpointAttributesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetEndpointAttributesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'GetEndpointAttributesResponse.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEndpointAttributesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetEndpointAttributesResponse copyWith(
          void Function(GetEndpointAttributesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetEndpointAttributesResponse))
          as GetEndpointAttributesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetEndpointAttributesResponse create() =>
      GetEndpointAttributesResponse._();
  @$core.override
  GetEndpointAttributesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetEndpointAttributesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetEndpointAttributesResponse>(create);
  static GetEndpointAttributesResponse? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);
}

class GetPlatformApplicationAttributesInput extends $pb.GeneratedMessage {
  factory GetPlatformApplicationAttributesInput({
    $core.String? platformapplicationarn,
  }) {
    final result = create();
    if (platformapplicationarn != null)
      result.platformapplicationarn = platformapplicationarn;
    return result;
  }

  GetPlatformApplicationAttributesInput._();

  factory GetPlatformApplicationAttributesInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPlatformApplicationAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPlatformApplicationAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(241250568, _omitFieldNames ? '' : 'platformapplicationarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPlatformApplicationAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPlatformApplicationAttributesInput copyWith(
          void Function(GetPlatformApplicationAttributesInput) updates) =>
      super.copyWith((message) =>
              updates(message as GetPlatformApplicationAttributesInput))
          as GetPlatformApplicationAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPlatformApplicationAttributesInput create() =>
      GetPlatformApplicationAttributesInput._();
  @$core.override
  GetPlatformApplicationAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPlatformApplicationAttributesInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetPlatformApplicationAttributesInput>(create);
  static GetPlatformApplicationAttributesInput? _defaultInstance;

  @$pb.TagNumber(241250568)
  $core.String get platformapplicationarn => $_getSZ(0);
  @$pb.TagNumber(241250568)
  set platformapplicationarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(241250568)
  $core.bool hasPlatformapplicationarn() => $_has(0);
  @$pb.TagNumber(241250568)
  void clearPlatformapplicationarn() => $_clearField(241250568);
}

class GetPlatformApplicationAttributesResponse extends $pb.GeneratedMessage {
  factory GetPlatformApplicationAttributesResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  GetPlatformApplicationAttributesResponse._();

  factory GetPlatformApplicationAttributesResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPlatformApplicationAttributesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPlatformApplicationAttributesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName:
            'GetPlatformApplicationAttributesResponse.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPlatformApplicationAttributesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPlatformApplicationAttributesResponse copyWith(
          void Function(GetPlatformApplicationAttributesResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetPlatformApplicationAttributesResponse))
          as GetPlatformApplicationAttributesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPlatformApplicationAttributesResponse create() =>
      GetPlatformApplicationAttributesResponse._();
  @$core.override
  GetPlatformApplicationAttributesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPlatformApplicationAttributesResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetPlatformApplicationAttributesResponse>(create);
  static GetPlatformApplicationAttributesResponse? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);
}

class GetSMSAttributesInput extends $pb.GeneratedMessage {
  factory GetSMSAttributesInput({
    $core.Iterable<$core.String>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addAll(attributes);
    return result;
  }

  GetSMSAttributesInput._();

  factory GetSMSAttributesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSMSAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSMSAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..pPS(33545109, _omitFieldNames ? '' : 'attributes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSMSAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSMSAttributesInput copyWith(
          void Function(GetSMSAttributesInput) updates) =>
      super.copyWith((message) => updates(message as GetSMSAttributesInput))
          as GetSMSAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSMSAttributesInput create() => GetSMSAttributesInput._();
  @$core.override
  GetSMSAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSMSAttributesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSMSAttributesInput>(create);
  static GetSMSAttributesInput? _defaultInstance;

  @$pb.TagNumber(33545109)
  $pb.PbList<$core.String> get attributes => $_getList(0);
}

class GetSMSAttributesResponse extends $pb.GeneratedMessage {
  factory GetSMSAttributesResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  GetSMSAttributesResponse._();

  factory GetSMSAttributesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSMSAttributesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSMSAttributesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        33545109, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'GetSMSAttributesResponse.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSMSAttributesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSMSAttributesResponse copyWith(
          void Function(GetSMSAttributesResponse) updates) =>
      super.copyWith((message) => updates(message as GetSMSAttributesResponse))
          as GetSMSAttributesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSMSAttributesResponse create() => GetSMSAttributesResponse._();
  @$core.override
  GetSMSAttributesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSMSAttributesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSMSAttributesResponse>(create);
  static GetSMSAttributesResponse? _defaultInstance;

  @$pb.TagNumber(33545109)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);
}

class GetSMSSandboxAccountStatusInput extends $pb.GeneratedMessage {
  factory GetSMSSandboxAccountStatusInput() => create();

  GetSMSSandboxAccountStatusInput._();

  factory GetSMSSandboxAccountStatusInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSMSSandboxAccountStatusInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSMSSandboxAccountStatusInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSMSSandboxAccountStatusInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSMSSandboxAccountStatusInput copyWith(
          void Function(GetSMSSandboxAccountStatusInput) updates) =>
      super.copyWith(
              (message) => updates(message as GetSMSSandboxAccountStatusInput))
          as GetSMSSandboxAccountStatusInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSMSSandboxAccountStatusInput create() =>
      GetSMSSandboxAccountStatusInput._();
  @$core.override
  GetSMSSandboxAccountStatusInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSMSSandboxAccountStatusInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSMSSandboxAccountStatusInput>(
          create);
  static GetSMSSandboxAccountStatusInput? _defaultInstance;
}

class GetSMSSandboxAccountStatusResult extends $pb.GeneratedMessage {
  factory GetSMSSandboxAccountStatusResult({
    $core.bool? isinsandbox,
  }) {
    final result = create();
    if (isinsandbox != null) result.isinsandbox = isinsandbox;
    return result;
  }

  GetSMSSandboxAccountStatusResult._();

  factory GetSMSSandboxAccountStatusResult.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSMSSandboxAccountStatusResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSMSSandboxAccountStatusResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOB(246264378, _omitFieldNames ? '' : 'isinsandbox')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSMSSandboxAccountStatusResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSMSSandboxAccountStatusResult copyWith(
          void Function(GetSMSSandboxAccountStatusResult) updates) =>
      super.copyWith(
              (message) => updates(message as GetSMSSandboxAccountStatusResult))
          as GetSMSSandboxAccountStatusResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSMSSandboxAccountStatusResult create() =>
      GetSMSSandboxAccountStatusResult._();
  @$core.override
  GetSMSSandboxAccountStatusResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSMSSandboxAccountStatusResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSMSSandboxAccountStatusResult>(
          create);
  static GetSMSSandboxAccountStatusResult? _defaultInstance;

  @$pb.TagNumber(246264378)
  $core.bool get isinsandbox => $_getBF(0);
  @$pb.TagNumber(246264378)
  set isinsandbox($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(246264378)
  $core.bool hasIsinsandbox() => $_has(0);
  @$pb.TagNumber(246264378)
  void clearIsinsandbox() => $_clearField(246264378);
}

class GetSubscriptionAttributesInput extends $pb.GeneratedMessage {
  factory GetSubscriptionAttributesInput({
    $core.String? subscriptionarn,
  }) {
    final result = create();
    if (subscriptionarn != null) result.subscriptionarn = subscriptionarn;
    return result;
  }

  GetSubscriptionAttributesInput._();

  factory GetSubscriptionAttributesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSubscriptionAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSubscriptionAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(279547820, _omitFieldNames ? '' : 'subscriptionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubscriptionAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubscriptionAttributesInput copyWith(
          void Function(GetSubscriptionAttributesInput) updates) =>
      super.copyWith(
              (message) => updates(message as GetSubscriptionAttributesInput))
          as GetSubscriptionAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSubscriptionAttributesInput create() =>
      GetSubscriptionAttributesInput._();
  @$core.override
  GetSubscriptionAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSubscriptionAttributesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSubscriptionAttributesInput>(create);
  static GetSubscriptionAttributesInput? _defaultInstance;

  @$pb.TagNumber(279547820)
  $core.String get subscriptionarn => $_getSZ(0);
  @$pb.TagNumber(279547820)
  set subscriptionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(279547820)
  $core.bool hasSubscriptionarn() => $_has(0);
  @$pb.TagNumber(279547820)
  void clearSubscriptionarn() => $_clearField(279547820);
}

class GetSubscriptionAttributesResponse extends $pb.GeneratedMessage {
  factory GetSubscriptionAttributesResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  GetSubscriptionAttributesResponse._();

  factory GetSubscriptionAttributesResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSubscriptionAttributesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSubscriptionAttributesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'GetSubscriptionAttributesResponse.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubscriptionAttributesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSubscriptionAttributesResponse copyWith(
          void Function(GetSubscriptionAttributesResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetSubscriptionAttributesResponse))
          as GetSubscriptionAttributesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSubscriptionAttributesResponse create() =>
      GetSubscriptionAttributesResponse._();
  @$core.override
  GetSubscriptionAttributesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSubscriptionAttributesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSubscriptionAttributesResponse>(
          create);
  static GetSubscriptionAttributesResponse? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);
}

class GetTopicAttributesInput extends $pb.GeneratedMessage {
  factory GetTopicAttributesInput({
    $core.String? topicarn,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    return result;
  }

  GetTopicAttributesInput._();

  factory GetTopicAttributesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTopicAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTopicAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTopicAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTopicAttributesInput copyWith(
          void Function(GetTopicAttributesInput) updates) =>
      super.copyWith((message) => updates(message as GetTopicAttributesInput))
          as GetTopicAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTopicAttributesInput create() => GetTopicAttributesInput._();
  @$core.override
  GetTopicAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTopicAttributesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTopicAttributesInput>(create);
  static GetTopicAttributesInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);
}

class GetTopicAttributesResponse extends $pb.GeneratedMessage {
  factory GetTopicAttributesResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  GetTopicAttributesResponse._();

  factory GetTopicAttributesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTopicAttributesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTopicAttributesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'GetTopicAttributesResponse.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTopicAttributesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTopicAttributesResponse copyWith(
          void Function(GetTopicAttributesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetTopicAttributesResponse))
          as GetTopicAttributesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTopicAttributesResponse create() => GetTopicAttributesResponse._();
  @$core.override
  GetTopicAttributesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTopicAttributesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTopicAttributesResponse>(create);
  static GetTopicAttributesResponse? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);
}

class InternalErrorException extends $pb.GeneratedMessage {
  factory InternalErrorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InternalErrorException._();

  factory InternalErrorException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InternalErrorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InternalErrorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalErrorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalErrorException copyWith(
          void Function(InternalErrorException) updates) =>
      super.copyWith((message) => updates(message as InternalErrorException))
          as InternalErrorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InternalErrorException create() => InternalErrorException._();
  @$core.override
  InternalErrorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InternalErrorException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InternalErrorException>(create);
  static InternalErrorException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidBatchEntryIdException extends $pb.GeneratedMessage {
  factory InvalidBatchEntryIdException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidBatchEntryIdException._();

  factory InvalidBatchEntryIdException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidBatchEntryIdException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidBatchEntryIdException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidBatchEntryIdException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidBatchEntryIdException copyWith(
          void Function(InvalidBatchEntryIdException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidBatchEntryIdException))
          as InvalidBatchEntryIdException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidBatchEntryIdException create() =>
      InvalidBatchEntryIdException._();
  @$core.override
  InvalidBatchEntryIdException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidBatchEntryIdException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidBatchEntryIdException>(create);
  static InvalidBatchEntryIdException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidParameterException extends $pb.GeneratedMessage {
  factory InvalidParameterException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidParameterException._();

  factory InvalidParameterException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidParameterException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidParameterException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterException copyWith(
          void Function(InvalidParameterException) updates) =>
      super.copyWith((message) => updates(message as InvalidParameterException))
          as InvalidParameterException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidParameterException create() => InvalidParameterException._();
  @$core.override
  InvalidParameterException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidParameterException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidParameterException>(create);
  static InvalidParameterException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidParameterValueException extends $pb.GeneratedMessage {
  factory InvalidParameterValueException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidParameterValueException._();

  factory InvalidParameterValueException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidParameterValueException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidParameterValueException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterValueException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterValueException copyWith(
          void Function(InvalidParameterValueException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidParameterValueException))
          as InvalidParameterValueException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidParameterValueException create() =>
      InvalidParameterValueException._();
  @$core.override
  InvalidParameterValueException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidParameterValueException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidParameterValueException>(create);
  static InvalidParameterValueException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidSecurityException extends $pb.GeneratedMessage {
  factory InvalidSecurityException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidSecurityException._();

  factory InvalidSecurityException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidSecurityException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidSecurityException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSecurityException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSecurityException copyWith(
          void Function(InvalidSecurityException) updates) =>
      super.copyWith((message) => updates(message as InvalidSecurityException))
          as InvalidSecurityException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidSecurityException create() => InvalidSecurityException._();
  @$core.override
  InvalidSecurityException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidSecurityException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidSecurityException>(create);
  static InvalidSecurityException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidStateException extends $pb.GeneratedMessage {
  factory InvalidStateException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidStateException._();

  factory InvalidStateException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidStateException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidStateException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidStateException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidStateException copyWith(
          void Function(InvalidStateException) updates) =>
      super.copyWith((message) => updates(message as InvalidStateException))
          as InvalidStateException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidStateException create() => InvalidStateException._();
  @$core.override
  InvalidStateException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidStateException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidStateException>(create);
  static InvalidStateException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSAccessDeniedException extends $pb.GeneratedMessage {
  factory KMSAccessDeniedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSAccessDeniedException._();

  factory KMSAccessDeniedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSAccessDeniedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSAccessDeniedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSAccessDeniedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSAccessDeniedException copyWith(
          void Function(KMSAccessDeniedException) updates) =>
      super.copyWith((message) => updates(message as KMSAccessDeniedException))
          as KMSAccessDeniedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSAccessDeniedException create() => KMSAccessDeniedException._();
  @$core.override
  KMSAccessDeniedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSAccessDeniedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSAccessDeniedException>(create);
  static KMSAccessDeniedException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSDisabledException extends $pb.GeneratedMessage {
  factory KMSDisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSDisabledException._();

  factory KMSDisabledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSDisabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSDisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSDisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSDisabledException copyWith(void Function(KMSDisabledException) updates) =>
      super.copyWith((message) => updates(message as KMSDisabledException))
          as KMSDisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSDisabledException create() => KMSDisabledException._();
  @$core.override
  KMSDisabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSDisabledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSDisabledException>(create);
  static KMSDisabledException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSInvalidStateException extends $pb.GeneratedMessage {
  factory KMSInvalidStateException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSInvalidStateException._();

  factory KMSInvalidStateException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSInvalidStateException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSInvalidStateException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSInvalidStateException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSInvalidStateException copyWith(
          void Function(KMSInvalidStateException) updates) =>
      super.copyWith((message) => updates(message as KMSInvalidStateException))
          as KMSInvalidStateException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSInvalidStateException create() => KMSInvalidStateException._();
  @$core.override
  KMSInvalidStateException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSInvalidStateException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSInvalidStateException>(create);
  static KMSInvalidStateException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSNotFoundException extends $pb.GeneratedMessage {
  factory KMSNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSNotFoundException._();

  factory KMSNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSNotFoundException copyWith(void Function(KMSNotFoundException) updates) =>
      super.copyWith((message) => updates(message as KMSNotFoundException))
          as KMSNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSNotFoundException create() => KMSNotFoundException._();
  @$core.override
  KMSNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSNotFoundException>(create);
  static KMSNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSOptInRequired extends $pb.GeneratedMessage {
  factory KMSOptInRequired({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSOptInRequired._();

  factory KMSOptInRequired.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSOptInRequired.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSOptInRequired',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSOptInRequired clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSOptInRequired copyWith(void Function(KMSOptInRequired) updates) =>
      super.copyWith((message) => updates(message as KMSOptInRequired))
          as KMSOptInRequired;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSOptInRequired create() => KMSOptInRequired._();
  @$core.override
  KMSOptInRequired createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSOptInRequired getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSOptInRequired>(create);
  static KMSOptInRequired? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSThrottlingException extends $pb.GeneratedMessage {
  factory KMSThrottlingException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSThrottlingException._();

  factory KMSThrottlingException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSThrottlingException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSThrottlingException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSThrottlingException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSThrottlingException copyWith(
          void Function(KMSThrottlingException) updates) =>
      super.copyWith((message) => updates(message as KMSThrottlingException))
          as KMSThrottlingException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSThrottlingException create() => KMSThrottlingException._();
  @$core.override
  KMSThrottlingException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSThrottlingException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSThrottlingException>(create);
  static KMSThrottlingException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ListEndpointsByPlatformApplicationInput extends $pb.GeneratedMessage {
  factory ListEndpointsByPlatformApplicationInput({
    $core.String? nexttoken,
    $core.String? platformapplicationarn,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (platformapplicationarn != null)
      result.platformapplicationarn = platformapplicationarn;
    return result;
  }

  ListEndpointsByPlatformApplicationInput._();

  factory ListEndpointsByPlatformApplicationInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEndpointsByPlatformApplicationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEndpointsByPlatformApplicationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(241250568, _omitFieldNames ? '' : 'platformapplicationarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEndpointsByPlatformApplicationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEndpointsByPlatformApplicationInput copyWith(
          void Function(ListEndpointsByPlatformApplicationInput) updates) =>
      super.copyWith((message) =>
              updates(message as ListEndpointsByPlatformApplicationInput))
          as ListEndpointsByPlatformApplicationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEndpointsByPlatformApplicationInput create() =>
      ListEndpointsByPlatformApplicationInput._();
  @$core.override
  ListEndpointsByPlatformApplicationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEndpointsByPlatformApplicationInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListEndpointsByPlatformApplicationInput>(create);
  static ListEndpointsByPlatformApplicationInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(241250568)
  $core.String get platformapplicationarn => $_getSZ(1);
  @$pb.TagNumber(241250568)
  set platformapplicationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(241250568)
  $core.bool hasPlatformapplicationarn() => $_has(1);
  @$pb.TagNumber(241250568)
  void clearPlatformapplicationarn() => $_clearField(241250568);
}

class ListEndpointsByPlatformApplicationResponse extends $pb.GeneratedMessage {
  factory ListEndpointsByPlatformApplicationResponse({
    $core.Iterable<Endpoint>? endpoints,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (endpoints != null) result.endpoints.addAll(endpoints);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListEndpointsByPlatformApplicationResponse._();

  factory ListEndpointsByPlatformApplicationResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEndpointsByPlatformApplicationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEndpointsByPlatformApplicationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..pPM<Endpoint>(16210494, _omitFieldNames ? '' : 'endpoints',
        subBuilder: Endpoint.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEndpointsByPlatformApplicationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEndpointsByPlatformApplicationResponse copyWith(
          void Function(ListEndpointsByPlatformApplicationResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListEndpointsByPlatformApplicationResponse))
          as ListEndpointsByPlatformApplicationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEndpointsByPlatformApplicationResponse create() =>
      ListEndpointsByPlatformApplicationResponse._();
  @$core.override
  ListEndpointsByPlatformApplicationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEndpointsByPlatformApplicationResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListEndpointsByPlatformApplicationResponse>(create);
  static ListEndpointsByPlatformApplicationResponse? _defaultInstance;

  @$pb.TagNumber(16210494)
  $pb.PbList<Endpoint> get endpoints => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListOriginationNumbersRequest extends $pb.GeneratedMessage {
  factory ListOriginationNumbersRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListOriginationNumbersRequest._();

  factory ListOriginationNumbersRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListOriginationNumbersRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListOriginationNumbersRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListOriginationNumbersRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListOriginationNumbersRequest copyWith(
          void Function(ListOriginationNumbersRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListOriginationNumbersRequest))
          as ListOriginationNumbersRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListOriginationNumbersRequest create() =>
      ListOriginationNumbersRequest._();
  @$core.override
  ListOriginationNumbersRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListOriginationNumbersRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListOriginationNumbersRequest>(create);
  static ListOriginationNumbersRequest? _defaultInstance;

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

class ListOriginationNumbersResult extends $pb.GeneratedMessage {
  factory ListOriginationNumbersResult({
    $core.String? nexttoken,
    $core.Iterable<PhoneNumberInformation>? phonenumbers,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (phonenumbers != null) result.phonenumbers.addAll(phonenumbers);
    return result;
  }

  ListOriginationNumbersResult._();

  factory ListOriginationNumbersResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListOriginationNumbersResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListOriginationNumbersResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<PhoneNumberInformation>(
        457192616, _omitFieldNames ? '' : 'phonenumbers',
        subBuilder: PhoneNumberInformation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListOriginationNumbersResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListOriginationNumbersResult copyWith(
          void Function(ListOriginationNumbersResult) updates) =>
      super.copyWith(
              (message) => updates(message as ListOriginationNumbersResult))
          as ListOriginationNumbersResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListOriginationNumbersResult create() =>
      ListOriginationNumbersResult._();
  @$core.override
  ListOriginationNumbersResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListOriginationNumbersResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListOriginationNumbersResult>(create);
  static ListOriginationNumbersResult? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(457192616)
  $pb.PbList<PhoneNumberInformation> get phonenumbers => $_getList(1);
}

class ListPhoneNumbersOptedOutInput extends $pb.GeneratedMessage {
  factory ListPhoneNumbersOptedOutInput({
    $core.String? nexttoken,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListPhoneNumbersOptedOutInput._();

  factory ListPhoneNumbersOptedOutInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPhoneNumbersOptedOutInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPhoneNumbersOptedOutInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPhoneNumbersOptedOutInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPhoneNumbersOptedOutInput copyWith(
          void Function(ListPhoneNumbersOptedOutInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListPhoneNumbersOptedOutInput))
          as ListPhoneNumbersOptedOutInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPhoneNumbersOptedOutInput create() =>
      ListPhoneNumbersOptedOutInput._();
  @$core.override
  ListPhoneNumbersOptedOutInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPhoneNumbersOptedOutInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPhoneNumbersOptedOutInput>(create);
  static ListPhoneNumbersOptedOutInput? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);
}

class ListPhoneNumbersOptedOutResponse extends $pb.GeneratedMessage {
  factory ListPhoneNumbersOptedOutResponse({
    $core.String? nexttoken,
    $core.Iterable<$core.String>? phonenumbers,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (phonenumbers != null) result.phonenumbers.addAll(phonenumbers);
    return result;
  }

  ListPhoneNumbersOptedOutResponse._();

  factory ListPhoneNumbersOptedOutResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPhoneNumbersOptedOutResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPhoneNumbersOptedOutResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(115833246, _omitFieldNames ? '' : 'nexttoken')
    ..pPS(156149576, _omitFieldNames ? '' : 'phonenumbers')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPhoneNumbersOptedOutResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPhoneNumbersOptedOutResponse copyWith(
          void Function(ListPhoneNumbersOptedOutResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListPhoneNumbersOptedOutResponse))
          as ListPhoneNumbersOptedOutResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPhoneNumbersOptedOutResponse create() =>
      ListPhoneNumbersOptedOutResponse._();
  @$core.override
  ListPhoneNumbersOptedOutResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPhoneNumbersOptedOutResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPhoneNumbersOptedOutResponse>(
          create);
  static ListPhoneNumbersOptedOutResponse? _defaultInstance;

  @$pb.TagNumber(115833246)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(115833246)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115833246)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(115833246)
  void clearNexttoken() => $_clearField(115833246);

  @$pb.TagNumber(156149576)
  $pb.PbList<$core.String> get phonenumbers => $_getList(1);
}

class ListPlatformApplicationsInput extends $pb.GeneratedMessage {
  factory ListPlatformApplicationsInput({
    $core.String? nexttoken,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListPlatformApplicationsInput._();

  factory ListPlatformApplicationsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPlatformApplicationsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPlatformApplicationsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPlatformApplicationsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPlatformApplicationsInput copyWith(
          void Function(ListPlatformApplicationsInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListPlatformApplicationsInput))
          as ListPlatformApplicationsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPlatformApplicationsInput create() =>
      ListPlatformApplicationsInput._();
  @$core.override
  ListPlatformApplicationsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPlatformApplicationsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPlatformApplicationsInput>(create);
  static ListPlatformApplicationsInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListPlatformApplicationsResponse extends $pb.GeneratedMessage {
  factory ListPlatformApplicationsResponse({
    $core.Iterable<PlatformApplication>? platformapplications,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (platformapplications != null)
      result.platformapplications.addAll(platformapplications);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListPlatformApplicationsResponse._();

  factory ListPlatformApplicationsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPlatformApplicationsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPlatformApplicationsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..pPM<PlatformApplication>(
        209735978, _omitFieldNames ? '' : 'platformapplications',
        subBuilder: PlatformApplication.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPlatformApplicationsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPlatformApplicationsResponse copyWith(
          void Function(ListPlatformApplicationsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListPlatformApplicationsResponse))
          as ListPlatformApplicationsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPlatformApplicationsResponse create() =>
      ListPlatformApplicationsResponse._();
  @$core.override
  ListPlatformApplicationsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPlatformApplicationsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPlatformApplicationsResponse>(
          create);
  static ListPlatformApplicationsResponse? _defaultInstance;

  @$pb.TagNumber(209735978)
  $pb.PbList<PlatformApplication> get platformapplications => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListSMSSandboxPhoneNumbersInput extends $pb.GeneratedMessage {
  factory ListSMSSandboxPhoneNumbersInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListSMSSandboxPhoneNumbersInput._();

  factory ListSMSSandboxPhoneNumbersInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSMSSandboxPhoneNumbersInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSMSSandboxPhoneNumbersInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSMSSandboxPhoneNumbersInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSMSSandboxPhoneNumbersInput copyWith(
          void Function(ListSMSSandboxPhoneNumbersInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListSMSSandboxPhoneNumbersInput))
          as ListSMSSandboxPhoneNumbersInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSMSSandboxPhoneNumbersInput create() =>
      ListSMSSandboxPhoneNumbersInput._();
  @$core.override
  ListSMSSandboxPhoneNumbersInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSMSSandboxPhoneNumbersInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSMSSandboxPhoneNumbersInput>(
          create);
  static ListSMSSandboxPhoneNumbersInput? _defaultInstance;

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

class ListSMSSandboxPhoneNumbersResult extends $pb.GeneratedMessage {
  factory ListSMSSandboxPhoneNumbersResult({
    $core.String? nexttoken,
    $core.Iterable<SMSSandboxPhoneNumber>? phonenumbers,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (phonenumbers != null) result.phonenumbers.addAll(phonenumbers);
    return result;
  }

  ListSMSSandboxPhoneNumbersResult._();

  factory ListSMSSandboxPhoneNumbersResult.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSMSSandboxPhoneNumbersResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSMSSandboxPhoneNumbersResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<SMSSandboxPhoneNumber>(
        457192616, _omitFieldNames ? '' : 'phonenumbers',
        subBuilder: SMSSandboxPhoneNumber.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSMSSandboxPhoneNumbersResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSMSSandboxPhoneNumbersResult copyWith(
          void Function(ListSMSSandboxPhoneNumbersResult) updates) =>
      super.copyWith(
              (message) => updates(message as ListSMSSandboxPhoneNumbersResult))
          as ListSMSSandboxPhoneNumbersResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSMSSandboxPhoneNumbersResult create() =>
      ListSMSSandboxPhoneNumbersResult._();
  @$core.override
  ListSMSSandboxPhoneNumbersResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSMSSandboxPhoneNumbersResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSMSSandboxPhoneNumbersResult>(
          create);
  static ListSMSSandboxPhoneNumbersResult? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(457192616)
  $pb.PbList<SMSSandboxPhoneNumber> get phonenumbers => $_getList(1);
}

class ListSubscriptionsByTopicInput extends $pb.GeneratedMessage {
  factory ListSubscriptionsByTopicInput({
    $core.String? topicarn,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListSubscriptionsByTopicInput._();

  factory ListSubscriptionsByTopicInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSubscriptionsByTopicInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSubscriptionsByTopicInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscriptionsByTopicInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscriptionsByTopicInput copyWith(
          void Function(ListSubscriptionsByTopicInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListSubscriptionsByTopicInput))
          as ListSubscriptionsByTopicInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSubscriptionsByTopicInput create() =>
      ListSubscriptionsByTopicInput._();
  @$core.override
  ListSubscriptionsByTopicInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSubscriptionsByTopicInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSubscriptionsByTopicInput>(create);
  static ListSubscriptionsByTopicInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListSubscriptionsByTopicResponse extends $pb.GeneratedMessage {
  factory ListSubscriptionsByTopicResponse({
    $core.Iterable<Subscription>? subscriptions,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (subscriptions != null) result.subscriptions.addAll(subscriptions);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListSubscriptionsByTopicResponse._();

  factory ListSubscriptionsByTopicResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSubscriptionsByTopicResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSubscriptionsByTopicResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..pPM<Subscription>(169711430, _omitFieldNames ? '' : 'subscriptions',
        subBuilder: Subscription.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscriptionsByTopicResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscriptionsByTopicResponse copyWith(
          void Function(ListSubscriptionsByTopicResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListSubscriptionsByTopicResponse))
          as ListSubscriptionsByTopicResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSubscriptionsByTopicResponse create() =>
      ListSubscriptionsByTopicResponse._();
  @$core.override
  ListSubscriptionsByTopicResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSubscriptionsByTopicResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSubscriptionsByTopicResponse>(
          create);
  static ListSubscriptionsByTopicResponse? _defaultInstance;

  @$pb.TagNumber(169711430)
  $pb.PbList<Subscription> get subscriptions => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListSubscriptionsInput extends $pb.GeneratedMessage {
  factory ListSubscriptionsInput({
    $core.String? nexttoken,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListSubscriptionsInput._();

  factory ListSubscriptionsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSubscriptionsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSubscriptionsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscriptionsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscriptionsInput copyWith(
          void Function(ListSubscriptionsInput) updates) =>
      super.copyWith((message) => updates(message as ListSubscriptionsInput))
          as ListSubscriptionsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSubscriptionsInput create() => ListSubscriptionsInput._();
  @$core.override
  ListSubscriptionsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSubscriptionsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSubscriptionsInput>(create);
  static ListSubscriptionsInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListSubscriptionsResponse extends $pb.GeneratedMessage {
  factory ListSubscriptionsResponse({
    $core.Iterable<Subscription>? subscriptions,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (subscriptions != null) result.subscriptions.addAll(subscriptions);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListSubscriptionsResponse._();

  factory ListSubscriptionsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSubscriptionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSubscriptionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..pPM<Subscription>(169711430, _omitFieldNames ? '' : 'subscriptions',
        subBuilder: Subscription.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscriptionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscriptionsResponse copyWith(
          void Function(ListSubscriptionsResponse) updates) =>
      super.copyWith((message) => updates(message as ListSubscriptionsResponse))
          as ListSubscriptionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSubscriptionsResponse create() => ListSubscriptionsResponse._();
  @$core.override
  ListSubscriptionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSubscriptionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSubscriptionsResponse>(create);
  static ListSubscriptionsResponse? _defaultInstance;

  @$pb.TagNumber(169711430)
  $pb.PbList<Subscription> get subscriptions => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
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

class ListTopicsInput extends $pb.GeneratedMessage {
  factory ListTopicsInput({
    $core.String? nexttoken,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListTopicsInput._();

  factory ListTopicsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTopicsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTopicsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTopicsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTopicsInput copyWith(void Function(ListTopicsInput) updates) =>
      super.copyWith((message) => updates(message as ListTopicsInput))
          as ListTopicsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTopicsInput create() => ListTopicsInput._();
  @$core.override
  ListTopicsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTopicsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTopicsInput>(create);
  static ListTopicsInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListTopicsResponse extends $pb.GeneratedMessage {
  factory ListTopicsResponse({
    $core.String? nexttoken,
    $core.Iterable<Topic>? topics,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (topics != null) result.topics.addAll(topics);
    return result;
  }

  ListTopicsResponse._();

  factory ListTopicsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTopicsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTopicsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Topic>(219850038, _omitFieldNames ? '' : 'topics',
        subBuilder: Topic.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTopicsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTopicsResponse copyWith(void Function(ListTopicsResponse) updates) =>
      super.copyWith((message) => updates(message as ListTopicsResponse))
          as ListTopicsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTopicsResponse create() => ListTopicsResponse._();
  @$core.override
  ListTopicsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTopicsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTopicsResponse>(create);
  static ListTopicsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(219850038)
  $pb.PbList<Topic> get topics => $_getList(1);
}

class MessageAttributeValue extends $pb.GeneratedMessage {
  factory MessageAttributeValue({
    $core.String? datatype,
    $core.String? stringvalue,
    $core.List<$core.int>? binaryvalue,
  }) {
    final result = create();
    if (datatype != null) result.datatype = datatype;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(67988590, _omitFieldNames ? '' : 'datatype')
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

  @$pb.TagNumber(184416138)
  $core.String get stringvalue => $_getSZ(1);
  @$pb.TagNumber(184416138)
  set stringvalue($core.String value) => $_setString(1, value);
  @$pb.TagNumber(184416138)
  $core.bool hasStringvalue() => $_has(1);
  @$pb.TagNumber(184416138)
  void clearStringvalue() => $_clearField(184416138);

  @$pb.TagNumber(255476278)
  $core.List<$core.int> get binaryvalue => $_getN(2);
  @$pb.TagNumber(255476278)
  set binaryvalue($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(255476278)
  $core.bool hasBinaryvalue() => $_has(2);
  @$pb.TagNumber(255476278)
  void clearBinaryvalue() => $_clearField(255476278);
}

class NotFoundException extends $pb.GeneratedMessage {
  factory NotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NotFoundException._();

  factory NotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotFoundException copyWith(void Function(NotFoundException) updates) =>
      super.copyWith((message) => updates(message as NotFoundException))
          as NotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NotFoundException create() => NotFoundException._();
  @$core.override
  NotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NotFoundException>(create);
  static NotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class OptInPhoneNumberInput extends $pb.GeneratedMessage {
  factory OptInPhoneNumberInput({
    $core.String? phonenumber,
  }) {
    final result = create();
    if (phonenumber != null) result.phonenumber = phonenumber;
    return result;
  }

  OptInPhoneNumberInput._();

  factory OptInPhoneNumberInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OptInPhoneNumberInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OptInPhoneNumberInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(156377999, _omitFieldNames ? '' : 'phonenumber')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OptInPhoneNumberInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OptInPhoneNumberInput copyWith(
          void Function(OptInPhoneNumberInput) updates) =>
      super.copyWith((message) => updates(message as OptInPhoneNumberInput))
          as OptInPhoneNumberInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OptInPhoneNumberInput create() => OptInPhoneNumberInput._();
  @$core.override
  OptInPhoneNumberInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OptInPhoneNumberInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OptInPhoneNumberInput>(create);
  static OptInPhoneNumberInput? _defaultInstance;

  @$pb.TagNumber(156377999)
  $core.String get phonenumber => $_getSZ(0);
  @$pb.TagNumber(156377999)
  set phonenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(156377999)
  $core.bool hasPhonenumber() => $_has(0);
  @$pb.TagNumber(156377999)
  void clearPhonenumber() => $_clearField(156377999);
}

class OptInPhoneNumberResponse extends $pb.GeneratedMessage {
  factory OptInPhoneNumberResponse() => create();

  OptInPhoneNumberResponse._();

  factory OptInPhoneNumberResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OptInPhoneNumberResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OptInPhoneNumberResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OptInPhoneNumberResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OptInPhoneNumberResponse copyWith(
          void Function(OptInPhoneNumberResponse) updates) =>
      super.copyWith((message) => updates(message as OptInPhoneNumberResponse))
          as OptInPhoneNumberResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OptInPhoneNumberResponse create() => OptInPhoneNumberResponse._();
  @$core.override
  OptInPhoneNumberResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OptInPhoneNumberResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OptInPhoneNumberResponse>(create);
  static OptInPhoneNumberResponse? _defaultInstance;
}

class OptedOutException extends $pb.GeneratedMessage {
  factory OptedOutException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  OptedOutException._();

  factory OptedOutException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OptedOutException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OptedOutException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OptedOutException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OptedOutException copyWith(void Function(OptedOutException) updates) =>
      super.copyWith((message) => updates(message as OptedOutException))
          as OptedOutException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OptedOutException create() => OptedOutException._();
  @$core.override
  OptedOutException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OptedOutException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OptedOutException>(create);
  static OptedOutException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PhoneNumberInformation extends $pb.GeneratedMessage {
  factory PhoneNumberInformation({
    $core.String? status,
    $core.Iterable<NumberCapability>? numbercapabilities,
    RouteType? routetype,
    $core.String? createdat,
    $core.String? iso2countrycode,
    $core.String? phonenumber,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (numbercapabilities != null)
      result.numbercapabilities.addAll(numbercapabilities);
    if (routetype != null) result.routetype = routetype;
    if (createdat != null) result.createdat = createdat;
    if (iso2countrycode != null) result.iso2countrycode = iso2countrycode;
    if (phonenumber != null) result.phonenumber = phonenumber;
    return result;
  }

  PhoneNumberInformation._();

  factory PhoneNumberInformation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PhoneNumberInformation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PhoneNumberInformation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(6222352, _omitFieldNames ? '' : 'status')
    ..pc<NumberCapability>(54004711,
        _omitFieldNames ? '' : 'numbercapabilities', $pb.PbFieldType.KE,
        valueOf: NumberCapability.valueOf,
        enumValues: NumberCapability.values,
        defaultEnumValue: NumberCapability.NUMBER_CAPABILITY_SMS)
    ..aE<RouteType>(170172127, _omitFieldNames ? '' : 'routetype',
        enumValues: RouteType.values)
    ..aOS(258192751, _omitFieldNames ? '' : 'createdat')
    ..aOS(283246908, _omitFieldNames ? '' : 'iso2countrycode')
    ..aOS(379600239, _omitFieldNames ? '' : 'phonenumber')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PhoneNumberInformation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PhoneNumberInformation copyWith(
          void Function(PhoneNumberInformation) updates) =>
      super.copyWith((message) => updates(message as PhoneNumberInformation))
          as PhoneNumberInformation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PhoneNumberInformation create() => PhoneNumberInformation._();
  @$core.override
  PhoneNumberInformation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PhoneNumberInformation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PhoneNumberInformation>(create);
  static PhoneNumberInformation? _defaultInstance;

  @$pb.TagNumber(6222352)
  $core.String get status => $_getSZ(0);
  @$pb.TagNumber(6222352)
  set status($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(54004711)
  $pb.PbList<NumberCapability> get numbercapabilities => $_getList(1);

  @$pb.TagNumber(170172127)
  RouteType get routetype => $_getN(2);
  @$pb.TagNumber(170172127)
  set routetype(RouteType value) => $_setField(170172127, value);
  @$pb.TagNumber(170172127)
  $core.bool hasRoutetype() => $_has(2);
  @$pb.TagNumber(170172127)
  void clearRoutetype() => $_clearField(170172127);

  @$pb.TagNumber(258192751)
  $core.String get createdat => $_getSZ(3);
  @$pb.TagNumber(258192751)
  set createdat($core.String value) => $_setString(3, value);
  @$pb.TagNumber(258192751)
  $core.bool hasCreatedat() => $_has(3);
  @$pb.TagNumber(258192751)
  void clearCreatedat() => $_clearField(258192751);

  @$pb.TagNumber(283246908)
  $core.String get iso2countrycode => $_getSZ(4);
  @$pb.TagNumber(283246908)
  set iso2countrycode($core.String value) => $_setString(4, value);
  @$pb.TagNumber(283246908)
  $core.bool hasIso2countrycode() => $_has(4);
  @$pb.TagNumber(283246908)
  void clearIso2countrycode() => $_clearField(283246908);

  @$pb.TagNumber(379600239)
  $core.String get phonenumber => $_getSZ(5);
  @$pb.TagNumber(379600239)
  set phonenumber($core.String value) => $_setString(5, value);
  @$pb.TagNumber(379600239)
  $core.bool hasPhonenumber() => $_has(5);
  @$pb.TagNumber(379600239)
  void clearPhonenumber() => $_clearField(379600239);
}

class PlatformApplication extends $pb.GeneratedMessage {
  factory PlatformApplication({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.String? platformapplicationarn,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    if (platformapplicationarn != null)
      result.platformapplicationarn = platformapplicationarn;
    return result;
  }

  PlatformApplication._();

  factory PlatformApplication.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PlatformApplication.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PlatformApplication',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'PlatformApplication.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..aOS(241250568, _omitFieldNames ? '' : 'platformapplicationarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlatformApplication clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlatformApplication copyWith(void Function(PlatformApplication) updates) =>
      super.copyWith((message) => updates(message as PlatformApplication))
          as PlatformApplication;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PlatformApplication create() => PlatformApplication._();
  @$core.override
  PlatformApplication createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PlatformApplication getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PlatformApplication>(create);
  static PlatformApplication? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);

  @$pb.TagNumber(241250568)
  $core.String get platformapplicationarn => $_getSZ(1);
  @$pb.TagNumber(241250568)
  set platformapplicationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(241250568)
  $core.bool hasPlatformapplicationarn() => $_has(1);
  @$pb.TagNumber(241250568)
  void clearPlatformapplicationarn() => $_clearField(241250568);
}

class PlatformApplicationDisabledException extends $pb.GeneratedMessage {
  factory PlatformApplicationDisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PlatformApplicationDisabledException._();

  factory PlatformApplicationDisabledException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PlatformApplicationDisabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PlatformApplicationDisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlatformApplicationDisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlatformApplicationDisabledException copyWith(
          void Function(PlatformApplicationDisabledException) updates) =>
      super.copyWith((message) =>
              updates(message as PlatformApplicationDisabledException))
          as PlatformApplicationDisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PlatformApplicationDisabledException create() =>
      PlatformApplicationDisabledException._();
  @$core.override
  PlatformApplicationDisabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PlatformApplicationDisabledException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          PlatformApplicationDisabledException>(create);
  static PlatformApplicationDisabledException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PublishBatchInput extends $pb.GeneratedMessage {
  factory PublishBatchInput({
    $core.String? topicarn,
    $core.Iterable<PublishBatchRequestEntry>? publishbatchrequestentries,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    if (publishbatchrequestentries != null)
      result.publishbatchrequestentries.addAll(publishbatchrequestentries);
    return result;
  }

  PublishBatchInput._();

  factory PublishBatchInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishBatchInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishBatchInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..pPM<PublishBatchRequestEntry>(
        327245576, _omitFieldNames ? '' : 'publishbatchrequestentries',
        subBuilder: PublishBatchRequestEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishBatchInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishBatchInput copyWith(void Function(PublishBatchInput) updates) =>
      super.copyWith((message) => updates(message as PublishBatchInput))
          as PublishBatchInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishBatchInput create() => PublishBatchInput._();
  @$core.override
  PublishBatchInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishBatchInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishBatchInput>(create);
  static PublishBatchInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(327245576)
  $pb.PbList<PublishBatchRequestEntry> get publishbatchrequestentries =>
      $_getList(1);
}

class PublishBatchRequestEntry extends $pb.GeneratedMessage {
  factory PublishBatchRequestEntry({
    $core.String? subject,
    $core.Iterable<$core.MapEntry<$core.String, MessageAttributeValue>>?
        messageattributes,
    $core.String? message,
    $core.String? messagededuplicationid,
    $core.String? id,
    $core.String? messagestructure,
    $core.String? messagegroupid,
  }) {
    final result = create();
    if (subject != null) result.subject = subject;
    if (messageattributes != null)
      result.messageattributes.addEntries(messageattributes);
    if (message != null) result.message = message;
    if (messagededuplicationid != null)
      result.messagededuplicationid = messagededuplicationid;
    if (id != null) result.id = id;
    if (messagestructure != null) result.messagestructure = messagestructure;
    if (messagegroupid != null) result.messagegroupid = messagegroupid;
    return result;
  }

  PublishBatchRequestEntry._();

  factory PublishBatchRequestEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishBatchRequestEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishBatchRequestEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(7939312, _omitFieldNames ? '' : 'subject')
    ..m<$core.String, MessageAttributeValue>(
        56443766, _omitFieldNames ? '' : 'messageattributes',
        entryClassName: 'PublishBatchRequestEntry.MessageattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MessageAttributeValue.create,
        valueDefaultOrMaker: MessageAttributeValue.getDefault,
        packageName: const $pb.PackageName('sns'))
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(379560665, _omitFieldNames ? '' : 'messagededuplicationid')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(402672330, _omitFieldNames ? '' : 'messagestructure')
    ..aOS(419537435, _omitFieldNames ? '' : 'messagegroupid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishBatchRequestEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishBatchRequestEntry copyWith(
          void Function(PublishBatchRequestEntry) updates) =>
      super.copyWith((message) => updates(message as PublishBatchRequestEntry))
          as PublishBatchRequestEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishBatchRequestEntry create() => PublishBatchRequestEntry._();
  @$core.override
  PublishBatchRequestEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishBatchRequestEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishBatchRequestEntry>(create);
  static PublishBatchRequestEntry? _defaultInstance;

  @$pb.TagNumber(7939312)
  $core.String get subject => $_getSZ(0);
  @$pb.TagNumber(7939312)
  set subject($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7939312)
  $core.bool hasSubject() => $_has(0);
  @$pb.TagNumber(7939312)
  void clearSubject() => $_clearField(7939312);

  @$pb.TagNumber(56443766)
  $pb.PbMap<$core.String, MessageAttributeValue> get messageattributes =>
      $_getMap(1);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(2);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(2, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(2);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(379560665)
  $core.String get messagededuplicationid => $_getSZ(3);
  @$pb.TagNumber(379560665)
  set messagededuplicationid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(379560665)
  $core.bool hasMessagededuplicationid() => $_has(3);
  @$pb.TagNumber(379560665)
  void clearMessagededuplicationid() => $_clearField(379560665);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(402672330)
  $core.String get messagestructure => $_getSZ(5);
  @$pb.TagNumber(402672330)
  set messagestructure($core.String value) => $_setString(5, value);
  @$pb.TagNumber(402672330)
  $core.bool hasMessagestructure() => $_has(5);
  @$pb.TagNumber(402672330)
  void clearMessagestructure() => $_clearField(402672330);

  @$pb.TagNumber(419537435)
  $core.String get messagegroupid => $_getSZ(6);
  @$pb.TagNumber(419537435)
  set messagegroupid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(419537435)
  $core.bool hasMessagegroupid() => $_has(6);
  @$pb.TagNumber(419537435)
  void clearMessagegroupid() => $_clearField(419537435);
}

class PublishBatchResponse extends $pb.GeneratedMessage {
  factory PublishBatchResponse({
    $core.Iterable<BatchResultErrorEntry>? failed,
    $core.Iterable<PublishBatchResultEntry>? successful,
  }) {
    final result = create();
    if (failed != null) result.failed.addAll(failed);
    if (successful != null) result.successful.addAll(successful);
    return result;
  }

  PublishBatchResponse._();

  factory PublishBatchResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishBatchResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishBatchResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..pPM<BatchResultErrorEntry>(360301525, _omitFieldNames ? '' : 'failed',
        subBuilder: BatchResultErrorEntry.create)
    ..pPM<PublishBatchResultEntry>(
        412818844, _omitFieldNames ? '' : 'successful',
        subBuilder: PublishBatchResultEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishBatchResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishBatchResponse copyWith(void Function(PublishBatchResponse) updates) =>
      super.copyWith((message) => updates(message as PublishBatchResponse))
          as PublishBatchResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishBatchResponse create() => PublishBatchResponse._();
  @$core.override
  PublishBatchResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishBatchResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishBatchResponse>(create);
  static PublishBatchResponse? _defaultInstance;

  @$pb.TagNumber(360301525)
  $pb.PbList<BatchResultErrorEntry> get failed => $_getList(0);

  @$pb.TagNumber(412818844)
  $pb.PbList<PublishBatchResultEntry> get successful => $_getList(1);
}

class PublishBatchResultEntry extends $pb.GeneratedMessage {
  factory PublishBatchResultEntry({
    $core.String? sequencenumber,
    $core.String? messageid,
    $core.String? id,
  }) {
    final result = create();
    if (sequencenumber != null) result.sequencenumber = sequencenumber;
    if (messageid != null) result.messageid = messageid;
    if (id != null) result.id = id;
    return result;
  }

  PublishBatchResultEntry._();

  factory PublishBatchResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishBatchResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishBatchResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(98094362, _omitFieldNames ? '' : 'sequencenumber')
    ..aOS(360526634, _omitFieldNames ? '' : 'messageid')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishBatchResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishBatchResultEntry copyWith(
          void Function(PublishBatchResultEntry) updates) =>
      super.copyWith((message) => updates(message as PublishBatchResultEntry))
          as PublishBatchResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishBatchResultEntry create() => PublishBatchResultEntry._();
  @$core.override
  PublishBatchResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishBatchResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishBatchResultEntry>(create);
  static PublishBatchResultEntry? _defaultInstance;

  @$pb.TagNumber(98094362)
  $core.String get sequencenumber => $_getSZ(0);
  @$pb.TagNumber(98094362)
  set sequencenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(98094362)
  $core.bool hasSequencenumber() => $_has(0);
  @$pb.TagNumber(98094362)
  void clearSequencenumber() => $_clearField(98094362);

  @$pb.TagNumber(360526634)
  $core.String get messageid => $_getSZ(1);
  @$pb.TagNumber(360526634)
  set messageid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(360526634)
  $core.bool hasMessageid() => $_has(1);
  @$pb.TagNumber(360526634)
  void clearMessageid() => $_clearField(360526634);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class PublishInput extends $pb.GeneratedMessage {
  factory PublishInput({
    $core.String? subject,
    $core.String? topicarn,
    $core.Iterable<$core.MapEntry<$core.String, MessageAttributeValue>>?
        messageattributes,
    $core.String? targetarn,
    $core.String? message,
    $core.String? messagededuplicationid,
    $core.String? phonenumber,
    $core.String? messagestructure,
    $core.String? messagegroupid,
  }) {
    final result = create();
    if (subject != null) result.subject = subject;
    if (topicarn != null) result.topicarn = topicarn;
    if (messageattributes != null)
      result.messageattributes.addEntries(messageattributes);
    if (targetarn != null) result.targetarn = targetarn;
    if (message != null) result.message = message;
    if (messagededuplicationid != null)
      result.messagededuplicationid = messagededuplicationid;
    if (phonenumber != null) result.phonenumber = phonenumber;
    if (messagestructure != null) result.messagestructure = messagestructure;
    if (messagegroupid != null) result.messagegroupid = messagegroupid;
    return result;
  }

  PublishInput._();

  factory PublishInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(7939312, _omitFieldNames ? '' : 'subject')
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..m<$core.String, MessageAttributeValue>(
        56443766, _omitFieldNames ? '' : 'messageattributes',
        entryClassName: 'PublishInput.MessageattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MessageAttributeValue.create,
        valueDefaultOrMaker: MessageAttributeValue.getDefault,
        packageName: const $pb.PackageName('sns'))
    ..aOS(217664144, _omitFieldNames ? '' : 'targetarn')
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(379560665, _omitFieldNames ? '' : 'messagededuplicationid')
    ..aOS(379600239, _omitFieldNames ? '' : 'phonenumber')
    ..aOS(402672330, _omitFieldNames ? '' : 'messagestructure')
    ..aOS(419537435, _omitFieldNames ? '' : 'messagegroupid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishInput copyWith(void Function(PublishInput) updates) =>
      super.copyWith((message) => updates(message as PublishInput))
          as PublishInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishInput create() => PublishInput._();
  @$core.override
  PublishInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishInput>(create);
  static PublishInput? _defaultInstance;

  @$pb.TagNumber(7939312)
  $core.String get subject => $_getSZ(0);
  @$pb.TagNumber(7939312)
  set subject($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7939312)
  $core.bool hasSubject() => $_has(0);
  @$pb.TagNumber(7939312)
  void clearSubject() => $_clearField(7939312);

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(1);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(1);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(56443766)
  $pb.PbMap<$core.String, MessageAttributeValue> get messageattributes =>
      $_getMap(2);

  @$pb.TagNumber(217664144)
  $core.String get targetarn => $_getSZ(3);
  @$pb.TagNumber(217664144)
  set targetarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(217664144)
  $core.bool hasTargetarn() => $_has(3);
  @$pb.TagNumber(217664144)
  void clearTargetarn() => $_clearField(217664144);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(4);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(4, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(4);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(379560665)
  $core.String get messagededuplicationid => $_getSZ(5);
  @$pb.TagNumber(379560665)
  set messagededuplicationid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(379560665)
  $core.bool hasMessagededuplicationid() => $_has(5);
  @$pb.TagNumber(379560665)
  void clearMessagededuplicationid() => $_clearField(379560665);

  @$pb.TagNumber(379600239)
  $core.String get phonenumber => $_getSZ(6);
  @$pb.TagNumber(379600239)
  set phonenumber($core.String value) => $_setString(6, value);
  @$pb.TagNumber(379600239)
  $core.bool hasPhonenumber() => $_has(6);
  @$pb.TagNumber(379600239)
  void clearPhonenumber() => $_clearField(379600239);

  @$pb.TagNumber(402672330)
  $core.String get messagestructure => $_getSZ(7);
  @$pb.TagNumber(402672330)
  set messagestructure($core.String value) => $_setString(7, value);
  @$pb.TagNumber(402672330)
  $core.bool hasMessagestructure() => $_has(7);
  @$pb.TagNumber(402672330)
  void clearMessagestructure() => $_clearField(402672330);

  @$pb.TagNumber(419537435)
  $core.String get messagegroupid => $_getSZ(8);
  @$pb.TagNumber(419537435)
  set messagegroupid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(419537435)
  $core.bool hasMessagegroupid() => $_has(8);
  @$pb.TagNumber(419537435)
  void clearMessagegroupid() => $_clearField(419537435);
}

class PublishResponse extends $pb.GeneratedMessage {
  factory PublishResponse({
    $core.String? sequencenumber,
    $core.String? messageid,
  }) {
    final result = create();
    if (sequencenumber != null) result.sequencenumber = sequencenumber;
    if (messageid != null) result.messageid = messageid;
    return result;
  }

  PublishResponse._();

  factory PublishResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublishResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublishResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(98094362, _omitFieldNames ? '' : 'sequencenumber')
    ..aOS(360526634, _omitFieldNames ? '' : 'messageid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublishResponse copyWith(void Function(PublishResponse) updates) =>
      super.copyWith((message) => updates(message as PublishResponse))
          as PublishResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublishResponse create() => PublishResponse._();
  @$core.override
  PublishResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublishResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublishResponse>(create);
  static PublishResponse? _defaultInstance;

  @$pb.TagNumber(98094362)
  $core.String get sequencenumber => $_getSZ(0);
  @$pb.TagNumber(98094362)
  set sequencenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(98094362)
  $core.bool hasSequencenumber() => $_has(0);
  @$pb.TagNumber(98094362)
  void clearSequencenumber() => $_clearField(98094362);

  @$pb.TagNumber(360526634)
  $core.String get messageid => $_getSZ(1);
  @$pb.TagNumber(360526634)
  set messageid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(360526634)
  $core.bool hasMessageid() => $_has(1);
  @$pb.TagNumber(360526634)
  void clearMessageid() => $_clearField(360526634);
}

class PutDataProtectionPolicyInput extends $pb.GeneratedMessage {
  factory PutDataProtectionPolicyInput({
    $core.String? resourcearn,
    $core.String? dataprotectionpolicy,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (dataprotectionpolicy != null)
      result.dataprotectionpolicy = dataprotectionpolicy;
    return result;
  }

  PutDataProtectionPolicyInput._();

  factory PutDataProtectionPolicyInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutDataProtectionPolicyInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutDataProtectionPolicyInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(519444573, _omitFieldNames ? '' : 'dataprotectionpolicy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutDataProtectionPolicyInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutDataProtectionPolicyInput copyWith(
          void Function(PutDataProtectionPolicyInput) updates) =>
      super.copyWith(
              (message) => updates(message as PutDataProtectionPolicyInput))
          as PutDataProtectionPolicyInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutDataProtectionPolicyInput create() =>
      PutDataProtectionPolicyInput._();
  @$core.override
  PutDataProtectionPolicyInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutDataProtectionPolicyInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutDataProtectionPolicyInput>(create);
  static PutDataProtectionPolicyInput? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

  @$pb.TagNumber(519444573)
  $core.String get dataprotectionpolicy => $_getSZ(1);
  @$pb.TagNumber(519444573)
  set dataprotectionpolicy($core.String value) => $_setString(1, value);
  @$pb.TagNumber(519444573)
  $core.bool hasDataprotectionpolicy() => $_has(1);
  @$pb.TagNumber(519444573)
  void clearDataprotectionpolicy() => $_clearField(519444573);
}

class RemovePermissionInput extends $pb.GeneratedMessage {
  factory RemovePermissionInput({
    $core.String? topicarn,
    $core.String? label,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    if (label != null) result.label = label;
    return result;
  }

  RemovePermissionInput._();

  factory RemovePermissionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemovePermissionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemovePermissionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..aOS(516747934, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemovePermissionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemovePermissionInput copyWith(
          void Function(RemovePermissionInput) updates) =>
      super.copyWith((message) => updates(message as RemovePermissionInput))
          as RemovePermissionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemovePermissionInput create() => RemovePermissionInput._();
  @$core.override
  RemovePermissionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemovePermissionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemovePermissionInput>(create);
  static RemovePermissionInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(516747934)
  $core.String get label => $_getSZ(1);
  @$pb.TagNumber(516747934)
  set label($core.String value) => $_setString(1, value);
  @$pb.TagNumber(516747934)
  $core.bool hasLabel() => $_has(1);
  @$pb.TagNumber(516747934)
  void clearLabel() => $_clearField(516747934);
}

class ReplayLimitExceededException extends $pb.GeneratedMessage {
  factory ReplayLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ReplayLimitExceededException._();

  factory ReplayLimitExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplayLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplayLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplayLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplayLimitExceededException copyWith(
          void Function(ReplayLimitExceededException) updates) =>
      super.copyWith(
              (message) => updates(message as ReplayLimitExceededException))
          as ReplayLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplayLimitExceededException create() =>
      ReplayLimitExceededException._();
  @$core.override
  ReplayLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplayLimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplayLimitExceededException>(create);
  static ReplayLimitExceededException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
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

class SMSSandboxPhoneNumber extends $pb.GeneratedMessage {
  factory SMSSandboxPhoneNumber({
    SMSSandboxPhoneNumberVerificationStatus? status,
    $core.String? phonenumber,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (phonenumber != null) result.phonenumber = phonenumber;
    return result;
  }

  SMSSandboxPhoneNumber._();

  factory SMSSandboxPhoneNumber.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SMSSandboxPhoneNumber.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SMSSandboxPhoneNumber',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aE<SMSSandboxPhoneNumberVerificationStatus>(
        6222352, _omitFieldNames ? '' : 'status',
        enumValues: SMSSandboxPhoneNumberVerificationStatus.values)
    ..aOS(379600239, _omitFieldNames ? '' : 'phonenumber')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SMSSandboxPhoneNumber clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SMSSandboxPhoneNumber copyWith(
          void Function(SMSSandboxPhoneNumber) updates) =>
      super.copyWith((message) => updates(message as SMSSandboxPhoneNumber))
          as SMSSandboxPhoneNumber;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SMSSandboxPhoneNumber create() => SMSSandboxPhoneNumber._();
  @$core.override
  SMSSandboxPhoneNumber createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SMSSandboxPhoneNumber getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SMSSandboxPhoneNumber>(create);
  static SMSSandboxPhoneNumber? _defaultInstance;

  @$pb.TagNumber(6222352)
  SMSSandboxPhoneNumberVerificationStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(SMSSandboxPhoneNumberVerificationStatus value) =>
      $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(379600239)
  $core.String get phonenumber => $_getSZ(1);
  @$pb.TagNumber(379600239)
  set phonenumber($core.String value) => $_setString(1, value);
  @$pb.TagNumber(379600239)
  $core.bool hasPhonenumber() => $_has(1);
  @$pb.TagNumber(379600239)
  void clearPhonenumber() => $_clearField(379600239);
}

class SetEndpointAttributesInput extends $pb.GeneratedMessage {
  factory SetEndpointAttributesInput({
    $core.String? endpointarn,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (endpointarn != null) result.endpointarn = endpointarn;
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  SetEndpointAttributesInput._();

  factory SetEndpointAttributesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetEndpointAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetEndpointAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(32228660, _omitFieldNames ? '' : 'endpointarn')
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'SetEndpointAttributesInput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetEndpointAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetEndpointAttributesInput copyWith(
          void Function(SetEndpointAttributesInput) updates) =>
      super.copyWith(
              (message) => updates(message as SetEndpointAttributesInput))
          as SetEndpointAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetEndpointAttributesInput create() => SetEndpointAttributesInput._();
  @$core.override
  SetEndpointAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetEndpointAttributesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetEndpointAttributesInput>(create);
  static SetEndpointAttributesInput? _defaultInstance;

  @$pb.TagNumber(32228660)
  $core.String get endpointarn => $_getSZ(0);
  @$pb.TagNumber(32228660)
  set endpointarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(32228660)
  $core.bool hasEndpointarn() => $_has(0);
  @$pb.TagNumber(32228660)
  void clearEndpointarn() => $_clearField(32228660);

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(1);
}

class SetPlatformApplicationAttributesInput extends $pb.GeneratedMessage {
  factory SetPlatformApplicationAttributesInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.String? platformapplicationarn,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    if (platformapplicationarn != null)
      result.platformapplicationarn = platformapplicationarn;
    return result;
  }

  SetPlatformApplicationAttributesInput._();

  factory SetPlatformApplicationAttributesInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetPlatformApplicationAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetPlatformApplicationAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'SetPlatformApplicationAttributesInput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..aOS(241250568, _omitFieldNames ? '' : 'platformapplicationarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetPlatformApplicationAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetPlatformApplicationAttributesInput copyWith(
          void Function(SetPlatformApplicationAttributesInput) updates) =>
      super.copyWith((message) =>
              updates(message as SetPlatformApplicationAttributesInput))
          as SetPlatformApplicationAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetPlatformApplicationAttributesInput create() =>
      SetPlatformApplicationAttributesInput._();
  @$core.override
  SetPlatformApplicationAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetPlatformApplicationAttributesInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          SetPlatformApplicationAttributesInput>(create);
  static SetPlatformApplicationAttributesInput? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);

  @$pb.TagNumber(241250568)
  $core.String get platformapplicationarn => $_getSZ(1);
  @$pb.TagNumber(241250568)
  set platformapplicationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(241250568)
  $core.bool hasPlatformapplicationarn() => $_has(1);
  @$pb.TagNumber(241250568)
  void clearPlatformapplicationarn() => $_clearField(241250568);
}

class SetSMSAttributesInput extends $pb.GeneratedMessage {
  factory SetSMSAttributesInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  SetSMSAttributesInput._();

  factory SetSMSAttributesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetSMSAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetSMSAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        33545109, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'SetSMSAttributesInput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetSMSAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetSMSAttributesInput copyWith(
          void Function(SetSMSAttributesInput) updates) =>
      super.copyWith((message) => updates(message as SetSMSAttributesInput))
          as SetSMSAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetSMSAttributesInput create() => SetSMSAttributesInput._();
  @$core.override
  SetSMSAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetSMSAttributesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetSMSAttributesInput>(create);
  static SetSMSAttributesInput? _defaultInstance;

  @$pb.TagNumber(33545109)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);
}

class SetSMSAttributesResponse extends $pb.GeneratedMessage {
  factory SetSMSAttributesResponse() => create();

  SetSMSAttributesResponse._();

  factory SetSMSAttributesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetSMSAttributesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetSMSAttributesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetSMSAttributesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetSMSAttributesResponse copyWith(
          void Function(SetSMSAttributesResponse) updates) =>
      super.copyWith((message) => updates(message as SetSMSAttributesResponse))
          as SetSMSAttributesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetSMSAttributesResponse create() => SetSMSAttributesResponse._();
  @$core.override
  SetSMSAttributesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetSMSAttributesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetSMSAttributesResponse>(create);
  static SetSMSAttributesResponse? _defaultInstance;
}

class SetSubscriptionAttributesInput extends $pb.GeneratedMessage {
  factory SetSubscriptionAttributesInput({
    $core.String? attributevalue,
    $core.String? subscriptionarn,
    $core.String? attributename,
  }) {
    final result = create();
    if (attributevalue != null) result.attributevalue = attributevalue;
    if (subscriptionarn != null) result.subscriptionarn = subscriptionarn;
    if (attributename != null) result.attributename = attributename;
    return result;
  }

  SetSubscriptionAttributesInput._();

  factory SetSubscriptionAttributesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetSubscriptionAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetSubscriptionAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(96769221, _omitFieldNames ? '' : 'attributevalue')
    ..aOS(279547820, _omitFieldNames ? '' : 'subscriptionarn')
    ..aOS(352717485, _omitFieldNames ? '' : 'attributename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetSubscriptionAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetSubscriptionAttributesInput copyWith(
          void Function(SetSubscriptionAttributesInput) updates) =>
      super.copyWith(
              (message) => updates(message as SetSubscriptionAttributesInput))
          as SetSubscriptionAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetSubscriptionAttributesInput create() =>
      SetSubscriptionAttributesInput._();
  @$core.override
  SetSubscriptionAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetSubscriptionAttributesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetSubscriptionAttributesInput>(create);
  static SetSubscriptionAttributesInput? _defaultInstance;

  @$pb.TagNumber(96769221)
  $core.String get attributevalue => $_getSZ(0);
  @$pb.TagNumber(96769221)
  set attributevalue($core.String value) => $_setString(0, value);
  @$pb.TagNumber(96769221)
  $core.bool hasAttributevalue() => $_has(0);
  @$pb.TagNumber(96769221)
  void clearAttributevalue() => $_clearField(96769221);

  @$pb.TagNumber(279547820)
  $core.String get subscriptionarn => $_getSZ(1);
  @$pb.TagNumber(279547820)
  set subscriptionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(279547820)
  $core.bool hasSubscriptionarn() => $_has(1);
  @$pb.TagNumber(279547820)
  void clearSubscriptionarn() => $_clearField(279547820);

  @$pb.TagNumber(352717485)
  $core.String get attributename => $_getSZ(2);
  @$pb.TagNumber(352717485)
  set attributename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(352717485)
  $core.bool hasAttributename() => $_has(2);
  @$pb.TagNumber(352717485)
  void clearAttributename() => $_clearField(352717485);
}

class SetTopicAttributesInput extends $pb.GeneratedMessage {
  factory SetTopicAttributesInput({
    $core.String? topicarn,
    $core.String? attributevalue,
    $core.String? attributename,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    if (attributevalue != null) result.attributevalue = attributevalue;
    if (attributename != null) result.attributename = attributename;
    return result;
  }

  SetTopicAttributesInput._();

  factory SetTopicAttributesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetTopicAttributesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetTopicAttributesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..aOS(96769221, _omitFieldNames ? '' : 'attributevalue')
    ..aOS(352717485, _omitFieldNames ? '' : 'attributename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetTopicAttributesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetTopicAttributesInput copyWith(
          void Function(SetTopicAttributesInput) updates) =>
      super.copyWith((message) => updates(message as SetTopicAttributesInput))
          as SetTopicAttributesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetTopicAttributesInput create() => SetTopicAttributesInput._();
  @$core.override
  SetTopicAttributesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetTopicAttributesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetTopicAttributesInput>(create);
  static SetTopicAttributesInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(96769221)
  $core.String get attributevalue => $_getSZ(1);
  @$pb.TagNumber(96769221)
  set attributevalue($core.String value) => $_setString(1, value);
  @$pb.TagNumber(96769221)
  $core.bool hasAttributevalue() => $_has(1);
  @$pb.TagNumber(96769221)
  void clearAttributevalue() => $_clearField(96769221);

  @$pb.TagNumber(352717485)
  $core.String get attributename => $_getSZ(2);
  @$pb.TagNumber(352717485)
  set attributename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(352717485)
  $core.bool hasAttributename() => $_has(2);
  @$pb.TagNumber(352717485)
  void clearAttributename() => $_clearField(352717485);
}

class StaleTagException extends $pb.GeneratedMessage {
  factory StaleTagException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  StaleTagException._();

  factory StaleTagException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StaleTagException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StaleTagException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StaleTagException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StaleTagException copyWith(void Function(StaleTagException) updates) =>
      super.copyWith((message) => updates(message as StaleTagException))
          as StaleTagException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StaleTagException create() => StaleTagException._();
  @$core.override
  StaleTagException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StaleTagException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StaleTagException>(create);
  static StaleTagException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class SubscribeInput extends $pb.GeneratedMessage {
  factory SubscribeInput({
    $core.String? topicarn,
    $core.bool? returnsubscriptionarn,
    $core.String? endpoint,
    $core.String? protocol,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    if (returnsubscriptionarn != null)
      result.returnsubscriptionarn = returnsubscriptionarn;
    if (endpoint != null) result.endpoint = endpoint;
    if (protocol != null) result.protocol = protocol;
    if (attributes != null) result.attributes.addEntries(attributes);
    return result;
  }

  SubscribeInput._();

  factory SubscribeInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubscribeInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..aOB(78754574, _omitFieldNames ? '' : 'returnsubscriptionarn')
    ..aOS(132634269, _omitFieldNames ? '' : 'endpoint')
    ..aOS(173534166, _omitFieldNames ? '' : 'protocol')
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'SubscribeInput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('sns'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeInput copyWith(void Function(SubscribeInput) updates) =>
      super.copyWith((message) => updates(message as SubscribeInput))
          as SubscribeInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubscribeInput create() => SubscribeInput._();
  @$core.override
  SubscribeInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubscribeInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribeInput>(create);
  static SubscribeInput? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(78754574)
  $core.bool get returnsubscriptionarn => $_getBF(1);
  @$pb.TagNumber(78754574)
  set returnsubscriptionarn($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(78754574)
  $core.bool hasReturnsubscriptionarn() => $_has(1);
  @$pb.TagNumber(78754574)
  void clearReturnsubscriptionarn() => $_clearField(78754574);

  @$pb.TagNumber(132634269)
  $core.String get endpoint => $_getSZ(2);
  @$pb.TagNumber(132634269)
  set endpoint($core.String value) => $_setString(2, value);
  @$pb.TagNumber(132634269)
  $core.bool hasEndpoint() => $_has(2);
  @$pb.TagNumber(132634269)
  void clearEndpoint() => $_clearField(132634269);

  @$pb.TagNumber(173534166)
  $core.String get protocol => $_getSZ(3);
  @$pb.TagNumber(173534166)
  set protocol($core.String value) => $_setString(3, value);
  @$pb.TagNumber(173534166)
  $core.bool hasProtocol() => $_has(3);
  @$pb.TagNumber(173534166)
  void clearProtocol() => $_clearField(173534166);

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(4);
}

class SubscribeResponse extends $pb.GeneratedMessage {
  factory SubscribeResponse({
    $core.String? subscriptionarn,
  }) {
    final result = create();
    if (subscriptionarn != null) result.subscriptionarn = subscriptionarn;
    return result;
  }

  SubscribeResponse._();

  factory SubscribeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubscribeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(279547820, _omitFieldNames ? '' : 'subscriptionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribeResponse copyWith(void Function(SubscribeResponse) updates) =>
      super.copyWith((message) => updates(message as SubscribeResponse))
          as SubscribeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubscribeResponse create() => SubscribeResponse._();
  @$core.override
  SubscribeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubscribeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribeResponse>(create);
  static SubscribeResponse? _defaultInstance;

  @$pb.TagNumber(279547820)
  $core.String get subscriptionarn => $_getSZ(0);
  @$pb.TagNumber(279547820)
  set subscriptionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(279547820)
  $core.bool hasSubscriptionarn() => $_has(0);
  @$pb.TagNumber(279547820)
  void clearSubscriptionarn() => $_clearField(279547820);
}

class Subscription extends $pb.GeneratedMessage {
  factory Subscription({
    $core.String? topicarn,
    $core.String? endpoint,
    $core.String? protocol,
    $core.String? subscriptionarn,
    $core.String? owner,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    if (endpoint != null) result.endpoint = endpoint;
    if (protocol != null) result.protocol = protocol;
    if (subscriptionarn != null) result.subscriptionarn = subscriptionarn;
    if (owner != null) result.owner = owner;
    return result;
  }

  Subscription._();

  factory Subscription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Subscription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Subscription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..aOS(132634269, _omitFieldNames ? '' : 'endpoint')
    ..aOS(173534166, _omitFieldNames ? '' : 'protocol')
    ..aOS(279547820, _omitFieldNames ? '' : 'subscriptionarn')
    ..aOS(455261813, _omitFieldNames ? '' : 'owner')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Subscription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Subscription copyWith(void Function(Subscription) updates) =>
      super.copyWith((message) => updates(message as Subscription))
          as Subscription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Subscription create() => Subscription._();
  @$core.override
  Subscription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Subscription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Subscription>(create);
  static Subscription? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);

  @$pb.TagNumber(132634269)
  $core.String get endpoint => $_getSZ(1);
  @$pb.TagNumber(132634269)
  set endpoint($core.String value) => $_setString(1, value);
  @$pb.TagNumber(132634269)
  $core.bool hasEndpoint() => $_has(1);
  @$pb.TagNumber(132634269)
  void clearEndpoint() => $_clearField(132634269);

  @$pb.TagNumber(173534166)
  $core.String get protocol => $_getSZ(2);
  @$pb.TagNumber(173534166)
  set protocol($core.String value) => $_setString(2, value);
  @$pb.TagNumber(173534166)
  $core.bool hasProtocol() => $_has(2);
  @$pb.TagNumber(173534166)
  void clearProtocol() => $_clearField(173534166);

  @$pb.TagNumber(279547820)
  $core.String get subscriptionarn => $_getSZ(3);
  @$pb.TagNumber(279547820)
  set subscriptionarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(279547820)
  $core.bool hasSubscriptionarn() => $_has(3);
  @$pb.TagNumber(279547820)
  void clearSubscriptionarn() => $_clearField(279547820);

  @$pb.TagNumber(455261813)
  $core.String get owner => $_getSZ(4);
  @$pb.TagNumber(455261813)
  set owner($core.String value) => $_setString(4, value);
  @$pb.TagNumber(455261813)
  $core.bool hasOwner() => $_has(4);
  @$pb.TagNumber(455261813)
  void clearOwner() => $_clearField(455261813);
}

class SubscriptionLimitExceededException extends $pb.GeneratedMessage {
  factory SubscriptionLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  SubscriptionLimitExceededException._();

  factory SubscriptionLimitExceededException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubscriptionLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscriptionLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscriptionLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscriptionLimitExceededException copyWith(
          void Function(SubscriptionLimitExceededException) updates) =>
      super.copyWith((message) =>
              updates(message as SubscriptionLimitExceededException))
          as SubscriptionLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubscriptionLimitExceededException create() =>
      SubscriptionLimitExceededException._();
  @$core.override
  SubscriptionLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubscriptionLimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscriptionLimitExceededException>(
          create);
  static SubscriptionLimitExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
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

class TagLimitExceededException extends $pb.GeneratedMessage {
  factory TagLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TagLimitExceededException._();

  factory TagLimitExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagLimitExceededException copyWith(
          void Function(TagLimitExceededException) updates) =>
      super.copyWith((message) => updates(message as TagLimitExceededException))
          as TagLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagLimitExceededException create() => TagLimitExceededException._();
  @$core.override
  TagLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagLimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagLimitExceededException>(create);
  static TagLimitExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TagPolicyException extends $pb.GeneratedMessage {
  factory TagPolicyException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TagPolicyException._();

  factory TagPolicyException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagPolicyException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagPolicyException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagPolicyException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagPolicyException copyWith(void Function(TagPolicyException) updates) =>
      super.copyWith((message) => updates(message as TagPolicyException))
          as TagPolicyException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagPolicyException create() => TagPolicyException._();
  @$core.override
  TagPolicyException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagPolicyException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagPolicyException>(create);
  static TagPolicyException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
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

class ThrottledException extends $pb.GeneratedMessage {
  factory ThrottledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ThrottledException._();

  factory ThrottledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ThrottledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ThrottledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ThrottledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ThrottledException copyWith(void Function(ThrottledException) updates) =>
      super.copyWith((message) => updates(message as ThrottledException))
          as ThrottledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ThrottledException create() => ThrottledException._();
  @$core.override
  ThrottledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ThrottledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ThrottledException>(create);
  static ThrottledException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TooManyEntriesInBatchRequestException extends $pb.GeneratedMessage {
  factory TooManyEntriesInBatchRequestException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyEntriesInBatchRequestException._();

  factory TooManyEntriesInBatchRequestException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyEntriesInBatchRequestException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyEntriesInBatchRequestException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyEntriesInBatchRequestException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyEntriesInBatchRequestException copyWith(
          void Function(TooManyEntriesInBatchRequestException) updates) =>
      super.copyWith((message) =>
              updates(message as TooManyEntriesInBatchRequestException))
          as TooManyEntriesInBatchRequestException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyEntriesInBatchRequestException create() =>
      TooManyEntriesInBatchRequestException._();
  @$core.override
  TooManyEntriesInBatchRequestException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyEntriesInBatchRequestException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          TooManyEntriesInBatchRequestException>(create);
  static TooManyEntriesInBatchRequestException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class Topic extends $pb.GeneratedMessage {
  factory Topic({
    $core.String? topicarn,
  }) {
    final result = create();
    if (topicarn != null) result.topicarn = topicarn;
    return result;
  }

  Topic._();

  factory Topic.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Topic.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Topic',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(30652956, _omitFieldNames ? '' : 'topicarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Topic clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Topic copyWith(void Function(Topic) updates) =>
      super.copyWith((message) => updates(message as Topic)) as Topic;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Topic create() => Topic._();
  @$core.override
  Topic createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Topic getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Topic>(create);
  static Topic? _defaultInstance;

  @$pb.TagNumber(30652956)
  $core.String get topicarn => $_getSZ(0);
  @$pb.TagNumber(30652956)
  set topicarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30652956)
  $core.bool hasTopicarn() => $_has(0);
  @$pb.TagNumber(30652956)
  void clearTopicarn() => $_clearField(30652956);
}

class TopicLimitExceededException extends $pb.GeneratedMessage {
  factory TopicLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TopicLimitExceededException._();

  factory TopicLimitExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TopicLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TopicLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TopicLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TopicLimitExceededException copyWith(
          void Function(TopicLimitExceededException) updates) =>
      super.copyWith(
              (message) => updates(message as TopicLimitExceededException))
          as TopicLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TopicLimitExceededException create() =>
      TopicLimitExceededException._();
  @$core.override
  TopicLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TopicLimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TopicLimitExceededException>(create);
  static TopicLimitExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class UnsubscribeInput extends $pb.GeneratedMessage {
  factory UnsubscribeInput({
    $core.String? subscriptionarn,
  }) {
    final result = create();
    if (subscriptionarn != null) result.subscriptionarn = subscriptionarn;
    return result;
  }

  UnsubscribeInput._();

  factory UnsubscribeInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnsubscribeInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnsubscribeInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(279547820, _omitFieldNames ? '' : 'subscriptionarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnsubscribeInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnsubscribeInput copyWith(void Function(UnsubscribeInput) updates) =>
      super.copyWith((message) => updates(message as UnsubscribeInput))
          as UnsubscribeInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnsubscribeInput create() => UnsubscribeInput._();
  @$core.override
  UnsubscribeInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnsubscribeInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnsubscribeInput>(create);
  static UnsubscribeInput? _defaultInstance;

  @$pb.TagNumber(279547820)
  $core.String get subscriptionarn => $_getSZ(0);
  @$pb.TagNumber(279547820)
  set subscriptionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(279547820)
  $core.bool hasSubscriptionarn() => $_has(0);
  @$pb.TagNumber(279547820)
  void clearSubscriptionarn() => $_clearField(279547820);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
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

class UserErrorException extends $pb.GeneratedMessage {
  factory UserErrorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  UserErrorException._();

  factory UserErrorException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UserErrorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UserErrorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UserErrorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UserErrorException copyWith(void Function(UserErrorException) updates) =>
      super.copyWith((message) => updates(message as UserErrorException))
          as UserErrorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UserErrorException create() => UserErrorException._();
  @$core.override
  UserErrorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UserErrorException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UserErrorException>(create);
  static UserErrorException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
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

class VerificationException extends $pb.GeneratedMessage {
  factory VerificationException({
    $core.String? status,
    $core.String? message,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (message != null) result.message = message;
    return result;
  }

  VerificationException._();

  factory VerificationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VerificationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VerificationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(6222352, _omitFieldNames ? '' : 'status')
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerificationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerificationException copyWith(
          void Function(VerificationException) updates) =>
      super.copyWith((message) => updates(message as VerificationException))
          as VerificationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerificationException create() => VerificationException._();
  @$core.override
  VerificationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VerificationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VerificationException>(create);
  static VerificationException? _defaultInstance;

  @$pb.TagNumber(6222352)
  $core.String get status => $_getSZ(0);
  @$pb.TagNumber(6222352)
  set status($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class VerifySMSSandboxPhoneNumberInput extends $pb.GeneratedMessage {
  factory VerifySMSSandboxPhoneNumberInput({
    $core.String? phonenumber,
    $core.String? onetimepassword,
  }) {
    final result = create();
    if (phonenumber != null) result.phonenumber = phonenumber;
    if (onetimepassword != null) result.onetimepassword = onetimepassword;
    return result;
  }

  VerifySMSSandboxPhoneNumberInput._();

  factory VerifySMSSandboxPhoneNumberInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VerifySMSSandboxPhoneNumberInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VerifySMSSandboxPhoneNumberInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..aOS(379600239, _omitFieldNames ? '' : 'phonenumber')
    ..aOS(417020900, _omitFieldNames ? '' : 'onetimepassword')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifySMSSandboxPhoneNumberInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifySMSSandboxPhoneNumberInput copyWith(
          void Function(VerifySMSSandboxPhoneNumberInput) updates) =>
      super.copyWith(
              (message) => updates(message as VerifySMSSandboxPhoneNumberInput))
          as VerifySMSSandboxPhoneNumberInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerifySMSSandboxPhoneNumberInput create() =>
      VerifySMSSandboxPhoneNumberInput._();
  @$core.override
  VerifySMSSandboxPhoneNumberInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VerifySMSSandboxPhoneNumberInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VerifySMSSandboxPhoneNumberInput>(
          create);
  static VerifySMSSandboxPhoneNumberInput? _defaultInstance;

  @$pb.TagNumber(379600239)
  $core.String get phonenumber => $_getSZ(0);
  @$pb.TagNumber(379600239)
  set phonenumber($core.String value) => $_setString(0, value);
  @$pb.TagNumber(379600239)
  $core.bool hasPhonenumber() => $_has(0);
  @$pb.TagNumber(379600239)
  void clearPhonenumber() => $_clearField(379600239);

  @$pb.TagNumber(417020900)
  $core.String get onetimepassword => $_getSZ(1);
  @$pb.TagNumber(417020900)
  set onetimepassword($core.String value) => $_setString(1, value);
  @$pb.TagNumber(417020900)
  $core.bool hasOnetimepassword() => $_has(1);
  @$pb.TagNumber(417020900)
  void clearOnetimepassword() => $_clearField(417020900);
}

class VerifySMSSandboxPhoneNumberResult extends $pb.GeneratedMessage {
  factory VerifySMSSandboxPhoneNumberResult() => create();

  VerifySMSSandboxPhoneNumberResult._();

  factory VerifySMSSandboxPhoneNumberResult.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VerifySMSSandboxPhoneNumberResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VerifySMSSandboxPhoneNumberResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'sns'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifySMSSandboxPhoneNumberResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifySMSSandboxPhoneNumberResult copyWith(
          void Function(VerifySMSSandboxPhoneNumberResult) updates) =>
      super.copyWith((message) =>
              updates(message as VerifySMSSandboxPhoneNumberResult))
          as VerifySMSSandboxPhoneNumberResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerifySMSSandboxPhoneNumberResult create() =>
      VerifySMSSandboxPhoneNumberResult._();
  @$core.override
  VerifySMSSandboxPhoneNumberResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VerifySMSSandboxPhoneNumberResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VerifySMSSandboxPhoneNumberResult>(
          create);
  static VerifySMSSandboxPhoneNumberResult? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
