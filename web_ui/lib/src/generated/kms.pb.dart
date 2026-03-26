// This is a generated file - do not edit.
//
// Generated from kms.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'kms.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'kms.pbenum.dart';

class AliasListEntry extends $pb.GeneratedMessage {
  factory AliasListEntry({
    $core.String? lastupdateddate,
    $core.String? creationdate,
    $core.String? aliasname,
    $core.String? targetkeyid,
    $core.String? aliasarn,
  }) {
    final result = create();
    if (lastupdateddate != null) result.lastupdateddate = lastupdateddate;
    if (creationdate != null) result.creationdate = creationdate;
    if (aliasname != null) result.aliasname = aliasname;
    if (targetkeyid != null) result.targetkeyid = targetkeyid;
    if (aliasarn != null) result.aliasarn = aliasarn;
    return result;
  }

  AliasListEntry._();

  factory AliasListEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AliasListEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AliasListEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(166338449, _omitFieldNames ? '' : 'lastupdateddate')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aOS(313250709, _omitFieldNames ? '' : 'aliasname')
    ..aOS(406196123, _omitFieldNames ? '' : 'targetkeyid')
    ..aOS(461101595, _omitFieldNames ? '' : 'aliasarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AliasListEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AliasListEntry copyWith(void Function(AliasListEntry) updates) =>
      super.copyWith((message) => updates(message as AliasListEntry))
          as AliasListEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AliasListEntry create() => AliasListEntry._();
  @$core.override
  AliasListEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AliasListEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AliasListEntry>(create);
  static AliasListEntry? _defaultInstance;

  @$pb.TagNumber(166338449)
  $core.String get lastupdateddate => $_getSZ(0);
  @$pb.TagNumber(166338449)
  set lastupdateddate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(166338449)
  $core.bool hasLastupdateddate() => $_has(0);
  @$pb.TagNumber(166338449)
  void clearLastupdateddate() => $_clearField(166338449);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(1);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(1);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(313250709)
  $core.String get aliasname => $_getSZ(2);
  @$pb.TagNumber(313250709)
  set aliasname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(313250709)
  $core.bool hasAliasname() => $_has(2);
  @$pb.TagNumber(313250709)
  void clearAliasname() => $_clearField(313250709);

  @$pb.TagNumber(406196123)
  $core.String get targetkeyid => $_getSZ(3);
  @$pb.TagNumber(406196123)
  set targetkeyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(406196123)
  $core.bool hasTargetkeyid() => $_has(3);
  @$pb.TagNumber(406196123)
  void clearTargetkeyid() => $_clearField(406196123);

  @$pb.TagNumber(461101595)
  $core.String get aliasarn => $_getSZ(4);
  @$pb.TagNumber(461101595)
  set aliasarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(461101595)
  $core.bool hasAliasarn() => $_has(4);
  @$pb.TagNumber(461101595)
  void clearAliasarn() => $_clearField(461101595);
}

class AlreadyExistsException extends $pb.GeneratedMessage {
  factory AlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  AlreadyExistsException._();

  factory AlreadyExistsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlreadyExistsException copyWith(
          void Function(AlreadyExistsException) updates) =>
      super.copyWith((message) => updates(message as AlreadyExistsException))
          as AlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AlreadyExistsException create() => AlreadyExistsException._();
  @$core.override
  AlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AlreadyExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AlreadyExistsException>(create);
  static AlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CancelKeyDeletionRequest extends $pb.GeneratedMessage {
  factory CancelKeyDeletionRequest({
    $core.String? keyid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  CancelKeyDeletionRequest._();

  factory CancelKeyDeletionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelKeyDeletionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelKeyDeletionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelKeyDeletionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelKeyDeletionRequest copyWith(
          void Function(CancelKeyDeletionRequest) updates) =>
      super.copyWith((message) => updates(message as CancelKeyDeletionRequest))
          as CancelKeyDeletionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelKeyDeletionRequest create() => CancelKeyDeletionRequest._();
  @$core.override
  CancelKeyDeletionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelKeyDeletionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelKeyDeletionRequest>(create);
  static CancelKeyDeletionRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class CancelKeyDeletionResponse extends $pb.GeneratedMessage {
  factory CancelKeyDeletionResponse({
    $core.String? keyid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  CancelKeyDeletionResponse._();

  factory CancelKeyDeletionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelKeyDeletionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelKeyDeletionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelKeyDeletionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelKeyDeletionResponse copyWith(
          void Function(CancelKeyDeletionResponse) updates) =>
      super.copyWith((message) => updates(message as CancelKeyDeletionResponse))
          as CancelKeyDeletionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelKeyDeletionResponse create() => CancelKeyDeletionResponse._();
  @$core.override
  CancelKeyDeletionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelKeyDeletionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelKeyDeletionResponse>(create);
  static CancelKeyDeletionResponse? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class CloudHsmClusterInUseException extends $pb.GeneratedMessage {
  factory CloudHsmClusterInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudHsmClusterInUseException._();

  factory CloudHsmClusterInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudHsmClusterInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudHsmClusterInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterInUseException copyWith(
          void Function(CloudHsmClusterInUseException) updates) =>
      super.copyWith(
              (message) => updates(message as CloudHsmClusterInUseException))
          as CloudHsmClusterInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterInUseException create() =>
      CloudHsmClusterInUseException._();
  @$core.override
  CloudHsmClusterInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloudHsmClusterInUseException>(create);
  static CloudHsmClusterInUseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CloudHsmClusterInvalidConfigurationException
    extends $pb.GeneratedMessage {
  factory CloudHsmClusterInvalidConfigurationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudHsmClusterInvalidConfigurationException._();

  factory CloudHsmClusterInvalidConfigurationException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudHsmClusterInvalidConfigurationException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudHsmClusterInvalidConfigurationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterInvalidConfigurationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterInvalidConfigurationException copyWith(
          void Function(CloudHsmClusterInvalidConfigurationException)
              updates) =>
      super.copyWith((message) =>
              updates(message as CloudHsmClusterInvalidConfigurationException))
          as CloudHsmClusterInvalidConfigurationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterInvalidConfigurationException create() =>
      CloudHsmClusterInvalidConfigurationException._();
  @$core.override
  CloudHsmClusterInvalidConfigurationException createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterInvalidConfigurationException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CloudHsmClusterInvalidConfigurationException>(create);
  static CloudHsmClusterInvalidConfigurationException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CloudHsmClusterNotActiveException extends $pb.GeneratedMessage {
  factory CloudHsmClusterNotActiveException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudHsmClusterNotActiveException._();

  factory CloudHsmClusterNotActiveException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudHsmClusterNotActiveException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudHsmClusterNotActiveException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterNotActiveException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterNotActiveException copyWith(
          void Function(CloudHsmClusterNotActiveException) updates) =>
      super.copyWith((message) =>
              updates(message as CloudHsmClusterNotActiveException))
          as CloudHsmClusterNotActiveException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterNotActiveException create() =>
      CloudHsmClusterNotActiveException._();
  @$core.override
  CloudHsmClusterNotActiveException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterNotActiveException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloudHsmClusterNotActiveException>(
          create);
  static CloudHsmClusterNotActiveException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CloudHsmClusterNotFoundException extends $pb.GeneratedMessage {
  factory CloudHsmClusterNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudHsmClusterNotFoundException._();

  factory CloudHsmClusterNotFoundException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudHsmClusterNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudHsmClusterNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterNotFoundException copyWith(
          void Function(CloudHsmClusterNotFoundException) updates) =>
      super.copyWith(
              (message) => updates(message as CloudHsmClusterNotFoundException))
          as CloudHsmClusterNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterNotFoundException create() =>
      CloudHsmClusterNotFoundException._();
  @$core.override
  CloudHsmClusterNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloudHsmClusterNotFoundException>(
          create);
  static CloudHsmClusterNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CloudHsmClusterNotRelatedException extends $pb.GeneratedMessage {
  factory CloudHsmClusterNotRelatedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CloudHsmClusterNotRelatedException._();

  factory CloudHsmClusterNotRelatedException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudHsmClusterNotRelatedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudHsmClusterNotRelatedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterNotRelatedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudHsmClusterNotRelatedException copyWith(
          void Function(CloudHsmClusterNotRelatedException) updates) =>
      super.copyWith((message) =>
              updates(message as CloudHsmClusterNotRelatedException))
          as CloudHsmClusterNotRelatedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterNotRelatedException create() =>
      CloudHsmClusterNotRelatedException._();
  @$core.override
  CloudHsmClusterNotRelatedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudHsmClusterNotRelatedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloudHsmClusterNotRelatedException>(
          create);
  static CloudHsmClusterNotRelatedException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
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

class ConnectCustomKeyStoreRequest extends $pb.GeneratedMessage {
  factory ConnectCustomKeyStoreRequest({
    $core.String? customkeystoreid,
  }) {
    final result = create();
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    return result;
  }

  ConnectCustomKeyStoreRequest._();

  factory ConnectCustomKeyStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectCustomKeyStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectCustomKeyStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectCustomKeyStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectCustomKeyStoreRequest copyWith(
          void Function(ConnectCustomKeyStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ConnectCustomKeyStoreRequest))
          as ConnectCustomKeyStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectCustomKeyStoreRequest create() =>
      ConnectCustomKeyStoreRequest._();
  @$core.override
  ConnectCustomKeyStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectCustomKeyStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectCustomKeyStoreRequest>(create);
  static ConnectCustomKeyStoreRequest? _defaultInstance;

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(0);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(0);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);
}

class ConnectCustomKeyStoreResponse extends $pb.GeneratedMessage {
  factory ConnectCustomKeyStoreResponse() => create();

  ConnectCustomKeyStoreResponse._();

  factory ConnectCustomKeyStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectCustomKeyStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectCustomKeyStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectCustomKeyStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectCustomKeyStoreResponse copyWith(
          void Function(ConnectCustomKeyStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ConnectCustomKeyStoreResponse))
          as ConnectCustomKeyStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectCustomKeyStoreResponse create() =>
      ConnectCustomKeyStoreResponse._();
  @$core.override
  ConnectCustomKeyStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectCustomKeyStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectCustomKeyStoreResponse>(create);
  static ConnectCustomKeyStoreResponse? _defaultInstance;
}

class CreateAliasRequest extends $pb.GeneratedMessage {
  factory CreateAliasRequest({
    $core.String? aliasname,
    $core.String? targetkeyid,
  }) {
    final result = create();
    if (aliasname != null) result.aliasname = aliasname;
    if (targetkeyid != null) result.targetkeyid = targetkeyid;
    return result;
  }

  CreateAliasRequest._();

  factory CreateAliasRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateAliasRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateAliasRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(313250709, _omitFieldNames ? '' : 'aliasname')
    ..aOS(406196123, _omitFieldNames ? '' : 'targetkeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAliasRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAliasRequest copyWith(void Function(CreateAliasRequest) updates) =>
      super.copyWith((message) => updates(message as CreateAliasRequest))
          as CreateAliasRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateAliasRequest create() => CreateAliasRequest._();
  @$core.override
  CreateAliasRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateAliasRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateAliasRequest>(create);
  static CreateAliasRequest? _defaultInstance;

  @$pb.TagNumber(313250709)
  $core.String get aliasname => $_getSZ(0);
  @$pb.TagNumber(313250709)
  set aliasname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(313250709)
  $core.bool hasAliasname() => $_has(0);
  @$pb.TagNumber(313250709)
  void clearAliasname() => $_clearField(313250709);

  @$pb.TagNumber(406196123)
  $core.String get targetkeyid => $_getSZ(1);
  @$pb.TagNumber(406196123)
  set targetkeyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(406196123)
  $core.bool hasTargetkeyid() => $_has(1);
  @$pb.TagNumber(406196123)
  void clearTargetkeyid() => $_clearField(406196123);
}

class CreateCustomKeyStoreRequest extends $pb.GeneratedMessage {
  factory CreateCustomKeyStoreRequest({
    $core.String? trustanchorcertificate,
    $core.String? xksproxyvpcendpointserviceowner,
    $core.String? cloudhsmclusterid,
    $core.String? customkeystorename,
    $core.String? xksproxyuriendpoint,
    XksProxyConnectivityType? xksproxyconnectivity,
    XksProxyAuthenticationCredentialType? xksproxyauthenticationcredential,
    $core.String? xksproxyvpcendpointservicename,
    $core.String? keystorepassword,
    CustomKeyStoreType? customkeystoretype,
    $core.String? xksproxyuripath,
  }) {
    final result = create();
    if (trustanchorcertificate != null)
      result.trustanchorcertificate = trustanchorcertificate;
    if (xksproxyvpcendpointserviceowner != null)
      result.xksproxyvpcendpointserviceowner = xksproxyvpcendpointserviceowner;
    if (cloudhsmclusterid != null) result.cloudhsmclusterid = cloudhsmclusterid;
    if (customkeystorename != null)
      result.customkeystorename = customkeystorename;
    if (xksproxyuriendpoint != null)
      result.xksproxyuriendpoint = xksproxyuriendpoint;
    if (xksproxyconnectivity != null)
      result.xksproxyconnectivity = xksproxyconnectivity;
    if (xksproxyauthenticationcredential != null)
      result.xksproxyauthenticationcredential =
          xksproxyauthenticationcredential;
    if (xksproxyvpcendpointservicename != null)
      result.xksproxyvpcendpointservicename = xksproxyvpcendpointservicename;
    if (keystorepassword != null) result.keystorepassword = keystorepassword;
    if (customkeystoretype != null)
      result.customkeystoretype = customkeystoretype;
    if (xksproxyuripath != null) result.xksproxyuripath = xksproxyuripath;
    return result;
  }

  CreateCustomKeyStoreRequest._();

  factory CreateCustomKeyStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateCustomKeyStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateCustomKeyStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(48354588, _omitFieldNames ? '' : 'trustanchorcertificate')
    ..aOS(55249590, _omitFieldNames ? '' : 'xksproxyvpcendpointserviceowner')
    ..aOS(56498754, _omitFieldNames ? '' : 'cloudhsmclusterid')
    ..aOS(170278046, _omitFieldNames ? '' : 'customkeystorename')
    ..aOS(273255559, _omitFieldNames ? '' : 'xksproxyuriendpoint')
    ..aE<XksProxyConnectivityType>(
        298569161, _omitFieldNames ? '' : 'xksproxyconnectivity',
        enumValues: XksProxyConnectivityType.values)
    ..aOM<XksProxyAuthenticationCredentialType>(
        350418199, _omitFieldNames ? '' : 'xksproxyauthenticationcredential',
        subBuilder: XksProxyAuthenticationCredentialType.create)
    ..aOS(372786130, _omitFieldNames ? '' : 'xksproxyvpcendpointservicename')
    ..aOS(403136353, _omitFieldNames ? '' : 'keystorepassword')
    ..aE<CustomKeyStoreType>(
        415647103, _omitFieldNames ? '' : 'customkeystoretype',
        enumValues: CustomKeyStoreType.values)
    ..aOS(436753509, _omitFieldNames ? '' : 'xksproxyuripath')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCustomKeyStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCustomKeyStoreRequest copyWith(
          void Function(CreateCustomKeyStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateCustomKeyStoreRequest))
          as CreateCustomKeyStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateCustomKeyStoreRequest create() =>
      CreateCustomKeyStoreRequest._();
  @$core.override
  CreateCustomKeyStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateCustomKeyStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateCustomKeyStoreRequest>(create);
  static CreateCustomKeyStoreRequest? _defaultInstance;

  @$pb.TagNumber(48354588)
  $core.String get trustanchorcertificate => $_getSZ(0);
  @$pb.TagNumber(48354588)
  set trustanchorcertificate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(48354588)
  $core.bool hasTrustanchorcertificate() => $_has(0);
  @$pb.TagNumber(48354588)
  void clearTrustanchorcertificate() => $_clearField(48354588);

  @$pb.TagNumber(55249590)
  $core.String get xksproxyvpcendpointserviceowner => $_getSZ(1);
  @$pb.TagNumber(55249590)
  set xksproxyvpcendpointserviceowner($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(55249590)
  $core.bool hasXksproxyvpcendpointserviceowner() => $_has(1);
  @$pb.TagNumber(55249590)
  void clearXksproxyvpcendpointserviceowner() => $_clearField(55249590);

  @$pb.TagNumber(56498754)
  $core.String get cloudhsmclusterid => $_getSZ(2);
  @$pb.TagNumber(56498754)
  set cloudhsmclusterid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(56498754)
  $core.bool hasCloudhsmclusterid() => $_has(2);
  @$pb.TagNumber(56498754)
  void clearCloudhsmclusterid() => $_clearField(56498754);

  @$pb.TagNumber(170278046)
  $core.String get customkeystorename => $_getSZ(3);
  @$pb.TagNumber(170278046)
  set customkeystorename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(170278046)
  $core.bool hasCustomkeystorename() => $_has(3);
  @$pb.TagNumber(170278046)
  void clearCustomkeystorename() => $_clearField(170278046);

  @$pb.TagNumber(273255559)
  $core.String get xksproxyuriendpoint => $_getSZ(4);
  @$pb.TagNumber(273255559)
  set xksproxyuriendpoint($core.String value) => $_setString(4, value);
  @$pb.TagNumber(273255559)
  $core.bool hasXksproxyuriendpoint() => $_has(4);
  @$pb.TagNumber(273255559)
  void clearXksproxyuriendpoint() => $_clearField(273255559);

  @$pb.TagNumber(298569161)
  XksProxyConnectivityType get xksproxyconnectivity => $_getN(5);
  @$pb.TagNumber(298569161)
  set xksproxyconnectivity(XksProxyConnectivityType value) =>
      $_setField(298569161, value);
  @$pb.TagNumber(298569161)
  $core.bool hasXksproxyconnectivity() => $_has(5);
  @$pb.TagNumber(298569161)
  void clearXksproxyconnectivity() => $_clearField(298569161);

  @$pb.TagNumber(350418199)
  XksProxyAuthenticationCredentialType get xksproxyauthenticationcredential =>
      $_getN(6);
  @$pb.TagNumber(350418199)
  set xksproxyauthenticationcredential(
          XksProxyAuthenticationCredentialType value) =>
      $_setField(350418199, value);
  @$pb.TagNumber(350418199)
  $core.bool hasXksproxyauthenticationcredential() => $_has(6);
  @$pb.TagNumber(350418199)
  void clearXksproxyauthenticationcredential() => $_clearField(350418199);
  @$pb.TagNumber(350418199)
  XksProxyAuthenticationCredentialType
      ensureXksproxyauthenticationcredential() => $_ensure(6);

  @$pb.TagNumber(372786130)
  $core.String get xksproxyvpcendpointservicename => $_getSZ(7);
  @$pb.TagNumber(372786130)
  set xksproxyvpcendpointservicename($core.String value) =>
      $_setString(7, value);
  @$pb.TagNumber(372786130)
  $core.bool hasXksproxyvpcendpointservicename() => $_has(7);
  @$pb.TagNumber(372786130)
  void clearXksproxyvpcendpointservicename() => $_clearField(372786130);

  @$pb.TagNumber(403136353)
  $core.String get keystorepassword => $_getSZ(8);
  @$pb.TagNumber(403136353)
  set keystorepassword($core.String value) => $_setString(8, value);
  @$pb.TagNumber(403136353)
  $core.bool hasKeystorepassword() => $_has(8);
  @$pb.TagNumber(403136353)
  void clearKeystorepassword() => $_clearField(403136353);

  @$pb.TagNumber(415647103)
  CustomKeyStoreType get customkeystoretype => $_getN(9);
  @$pb.TagNumber(415647103)
  set customkeystoretype(CustomKeyStoreType value) =>
      $_setField(415647103, value);
  @$pb.TagNumber(415647103)
  $core.bool hasCustomkeystoretype() => $_has(9);
  @$pb.TagNumber(415647103)
  void clearCustomkeystoretype() => $_clearField(415647103);

  @$pb.TagNumber(436753509)
  $core.String get xksproxyuripath => $_getSZ(10);
  @$pb.TagNumber(436753509)
  set xksproxyuripath($core.String value) => $_setString(10, value);
  @$pb.TagNumber(436753509)
  $core.bool hasXksproxyuripath() => $_has(10);
  @$pb.TagNumber(436753509)
  void clearXksproxyuripath() => $_clearField(436753509);
}

class CreateCustomKeyStoreResponse extends $pb.GeneratedMessage {
  factory CreateCustomKeyStoreResponse({
    $core.String? customkeystoreid,
  }) {
    final result = create();
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    return result;
  }

  CreateCustomKeyStoreResponse._();

  factory CreateCustomKeyStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateCustomKeyStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateCustomKeyStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCustomKeyStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCustomKeyStoreResponse copyWith(
          void Function(CreateCustomKeyStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateCustomKeyStoreResponse))
          as CreateCustomKeyStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateCustomKeyStoreResponse create() =>
      CreateCustomKeyStoreResponse._();
  @$core.override
  CreateCustomKeyStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateCustomKeyStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateCustomKeyStoreResponse>(create);
  static CreateCustomKeyStoreResponse? _defaultInstance;

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(0);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(0);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);
}

class CreateGrantRequest extends $pb.GeneratedMessage {
  factory CreateGrantRequest({
    $core.String? retiringprincipal,
    $core.bool? dryrun,
    $core.Iterable<GrantOperation>? operations,
    $core.String? granteeprincipal,
    $core.String? name,
    $core.String? keyid,
    GrantConstraints? constraints,
    $core.Iterable<$core.String>? granttokens,
  }) {
    final result = create();
    if (retiringprincipal != null) result.retiringprincipal = retiringprincipal;
    if (dryrun != null) result.dryrun = dryrun;
    if (operations != null) result.operations.addAll(operations);
    if (granteeprincipal != null) result.granteeprincipal = granteeprincipal;
    if (name != null) result.name = name;
    if (keyid != null) result.keyid = keyid;
    if (constraints != null) result.constraints = constraints;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    return result;
  }

  CreateGrantRequest._();

  factory CreateGrantRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateGrantRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateGrantRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(49541086, _omitFieldNames ? '' : 'retiringprincipal')
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..pc<GrantOperation>(
        126776656, _omitFieldNames ? '' : 'operations', $pb.PbFieldType.KE,
        valueOf: GrantOperation.valueOf,
        enumValues: GrantOperation.values,
        defaultEnumValue: GrantOperation.GRANT_OPERATION_GENERATEDATAKEYPAIR)
    ..aOS(234727364, _omitFieldNames ? '' : 'granteeprincipal')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOM<GrantConstraints>(302297388, _omitFieldNames ? '' : 'constraints',
        subBuilder: GrantConstraints.create)
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGrantRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGrantRequest copyWith(void Function(CreateGrantRequest) updates) =>
      super.copyWith((message) => updates(message as CreateGrantRequest))
          as CreateGrantRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateGrantRequest create() => CreateGrantRequest._();
  @$core.override
  CreateGrantRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateGrantRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateGrantRequest>(create);
  static CreateGrantRequest? _defaultInstance;

  @$pb.TagNumber(49541086)
  $core.String get retiringprincipal => $_getSZ(0);
  @$pb.TagNumber(49541086)
  set retiringprincipal($core.String value) => $_setString(0, value);
  @$pb.TagNumber(49541086)
  $core.bool hasRetiringprincipal() => $_has(0);
  @$pb.TagNumber(49541086)
  void clearRetiringprincipal() => $_clearField(49541086);

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(1);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(1);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(126776656)
  $pb.PbList<GrantOperation> get operations => $_getList(2);

  @$pb.TagNumber(234727364)
  $core.String get granteeprincipal => $_getSZ(3);
  @$pb.TagNumber(234727364)
  set granteeprincipal($core.String value) => $_setString(3, value);
  @$pb.TagNumber(234727364)
  $core.bool hasGranteeprincipal() => $_has(3);
  @$pb.TagNumber(234727364)
  void clearGranteeprincipal() => $_clearField(234727364);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(5);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(5);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(302297388)
  GrantConstraints get constraints => $_getN(6);
  @$pb.TagNumber(302297388)
  set constraints(GrantConstraints value) => $_setField(302297388, value);
  @$pb.TagNumber(302297388)
  $core.bool hasConstraints() => $_has(6);
  @$pb.TagNumber(302297388)
  void clearConstraints() => $_clearField(302297388);
  @$pb.TagNumber(302297388)
  GrantConstraints ensureConstraints() => $_ensure(6);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(7);
}

class CreateGrantResponse extends $pb.GeneratedMessage {
  factory CreateGrantResponse({
    $core.String? grantid,
    $core.String? granttoken,
  }) {
    final result = create();
    if (grantid != null) result.grantid = grantid;
    if (granttoken != null) result.granttoken = granttoken;
    return result;
  }

  CreateGrantResponse._();

  factory CreateGrantResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateGrantResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateGrantResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(66852281, _omitFieldNames ? '' : 'grantid')
    ..aOS(137683547, _omitFieldNames ? '' : 'granttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGrantResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGrantResponse copyWith(void Function(CreateGrantResponse) updates) =>
      super.copyWith((message) => updates(message as CreateGrantResponse))
          as CreateGrantResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateGrantResponse create() => CreateGrantResponse._();
  @$core.override
  CreateGrantResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateGrantResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateGrantResponse>(create);
  static CreateGrantResponse? _defaultInstance;

  @$pb.TagNumber(66852281)
  $core.String get grantid => $_getSZ(0);
  @$pb.TagNumber(66852281)
  set grantid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(66852281)
  $core.bool hasGrantid() => $_has(0);
  @$pb.TagNumber(66852281)
  void clearGrantid() => $_clearField(66852281);

  @$pb.TagNumber(137683547)
  $core.String get granttoken => $_getSZ(1);
  @$pb.TagNumber(137683547)
  set granttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(137683547)
  $core.bool hasGranttoken() => $_has(1);
  @$pb.TagNumber(137683547)
  void clearGranttoken() => $_clearField(137683547);
}

class CreateKeyRequest extends $pb.GeneratedMessage {
  factory CreateKeyRequest({
    $core.String? xkskeyid,
    $core.String? customkeystoreid,
    $core.String? description,
    KeySpec? keyspec,
    $core.bool? bypasspolicylockoutsafetycheck,
    KeyUsageType? keyusage,
    $core.Iterable<Tag>? tags,
    $core.bool? multiregion,
    $core.String? policy,
    CustomerMasterKeySpec? customermasterkeyspec,
    OriginType? origin,
  }) {
    final result = create();
    if (xkskeyid != null) result.xkskeyid = xkskeyid;
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    if (description != null) result.description = description;
    if (keyspec != null) result.keyspec = keyspec;
    if (bypasspolicylockoutsafetycheck != null)
      result.bypasspolicylockoutsafetycheck = bypasspolicylockoutsafetycheck;
    if (keyusage != null) result.keyusage = keyusage;
    if (tags != null) result.tags.addAll(tags);
    if (multiregion != null) result.multiregion = multiregion;
    if (policy != null) result.policy = policy;
    if (customermasterkeyspec != null)
      result.customermasterkeyspec = customermasterkeyspec;
    if (origin != null) result.origin = origin;
    return result;
  }

  CreateKeyRequest._();

  factory CreateKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(12647506, _omitFieldNames ? '' : 'xkskeyid')
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aE<KeySpec>(138220928, _omitFieldNames ? '' : 'keyspec',
        enumValues: KeySpec.values)
    ..aOB(177450851, _omitFieldNames ? '' : 'bypasspolicylockoutsafetycheck')
    ..aE<KeyUsageType>(357216772, _omitFieldNames ? '' : 'keyusage',
        enumValues: KeyUsageType.values)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOB(405769103, _omitFieldNames ? '' : 'multiregion')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..aE<CustomerMasterKeySpec>(
        472470930, _omitFieldNames ? '' : 'customermasterkeyspec',
        enumValues: CustomerMasterKeySpec.values)
    ..aE<OriginType>(529732720, _omitFieldNames ? '' : 'origin',
        enumValues: OriginType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateKeyRequest copyWith(void Function(CreateKeyRequest) updates) =>
      super.copyWith((message) => updates(message as CreateKeyRequest))
          as CreateKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateKeyRequest create() => CreateKeyRequest._();
  @$core.override
  CreateKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateKeyRequest>(create);
  static CreateKeyRequest? _defaultInstance;

  @$pb.TagNumber(12647506)
  $core.String get xkskeyid => $_getSZ(0);
  @$pb.TagNumber(12647506)
  set xkskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(12647506)
  $core.bool hasXkskeyid() => $_has(0);
  @$pb.TagNumber(12647506)
  void clearXkskeyid() => $_clearField(12647506);

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(1);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(1);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(138220928)
  KeySpec get keyspec => $_getN(3);
  @$pb.TagNumber(138220928)
  set keyspec(KeySpec value) => $_setField(138220928, value);
  @$pb.TagNumber(138220928)
  $core.bool hasKeyspec() => $_has(3);
  @$pb.TagNumber(138220928)
  void clearKeyspec() => $_clearField(138220928);

  @$pb.TagNumber(177450851)
  $core.bool get bypasspolicylockoutsafetycheck => $_getBF(4);
  @$pb.TagNumber(177450851)
  set bypasspolicylockoutsafetycheck($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(177450851)
  $core.bool hasBypasspolicylockoutsafetycheck() => $_has(4);
  @$pb.TagNumber(177450851)
  void clearBypasspolicylockoutsafetycheck() => $_clearField(177450851);

  @$pb.TagNumber(357216772)
  KeyUsageType get keyusage => $_getN(5);
  @$pb.TagNumber(357216772)
  set keyusage(KeyUsageType value) => $_setField(357216772, value);
  @$pb.TagNumber(357216772)
  $core.bool hasKeyusage() => $_has(5);
  @$pb.TagNumber(357216772)
  void clearKeyusage() => $_clearField(357216772);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(6);

  @$pb.TagNumber(405769103)
  $core.bool get multiregion => $_getBF(7);
  @$pb.TagNumber(405769103)
  set multiregion($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(405769103)
  $core.bool hasMultiregion() => $_has(7);
  @$pb.TagNumber(405769103)
  void clearMultiregion() => $_clearField(405769103);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(8);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(8, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(8);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);

  @$pb.TagNumber(472470930)
  CustomerMasterKeySpec get customermasterkeyspec => $_getN(9);
  @$pb.TagNumber(472470930)
  set customermasterkeyspec(CustomerMasterKeySpec value) =>
      $_setField(472470930, value);
  @$pb.TagNumber(472470930)
  $core.bool hasCustomermasterkeyspec() => $_has(9);
  @$pb.TagNumber(472470930)
  void clearCustomermasterkeyspec() => $_clearField(472470930);

  @$pb.TagNumber(529732720)
  OriginType get origin => $_getN(10);
  @$pb.TagNumber(529732720)
  set origin(OriginType value) => $_setField(529732720, value);
  @$pb.TagNumber(529732720)
  $core.bool hasOrigin() => $_has(10);
  @$pb.TagNumber(529732720)
  void clearOrigin() => $_clearField(529732720);
}

class CreateKeyResponse extends $pb.GeneratedMessage {
  factory CreateKeyResponse({
    KeyMetadata? keymetadata,
  }) {
    final result = create();
    if (keymetadata != null) result.keymetadata = keymetadata;
    return result;
  }

  CreateKeyResponse._();

  factory CreateKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOM<KeyMetadata>(149520202, _omitFieldNames ? '' : 'keymetadata',
        subBuilder: KeyMetadata.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateKeyResponse copyWith(void Function(CreateKeyResponse) updates) =>
      super.copyWith((message) => updates(message as CreateKeyResponse))
          as CreateKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateKeyResponse create() => CreateKeyResponse._();
  @$core.override
  CreateKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateKeyResponse>(create);
  static CreateKeyResponse? _defaultInstance;

  @$pb.TagNumber(149520202)
  KeyMetadata get keymetadata => $_getN(0);
  @$pb.TagNumber(149520202)
  set keymetadata(KeyMetadata value) => $_setField(149520202, value);
  @$pb.TagNumber(149520202)
  $core.bool hasKeymetadata() => $_has(0);
  @$pb.TagNumber(149520202)
  void clearKeymetadata() => $_clearField(149520202);
  @$pb.TagNumber(149520202)
  KeyMetadata ensureKeymetadata() => $_ensure(0);
}

class CustomKeyStoreHasCMKsException extends $pb.GeneratedMessage {
  factory CustomKeyStoreHasCMKsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CustomKeyStoreHasCMKsException._();

  factory CustomKeyStoreHasCMKsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomKeyStoreHasCMKsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomKeyStoreHasCMKsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoreHasCMKsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoreHasCMKsException copyWith(
          void Function(CustomKeyStoreHasCMKsException) updates) =>
      super.copyWith(
              (message) => updates(message as CustomKeyStoreHasCMKsException))
          as CustomKeyStoreHasCMKsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomKeyStoreHasCMKsException create() =>
      CustomKeyStoreHasCMKsException._();
  @$core.override
  CustomKeyStoreHasCMKsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomKeyStoreHasCMKsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CustomKeyStoreHasCMKsException>(create);
  static CustomKeyStoreHasCMKsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CustomKeyStoreInvalidStateException extends $pb.GeneratedMessage {
  factory CustomKeyStoreInvalidStateException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CustomKeyStoreInvalidStateException._();

  factory CustomKeyStoreInvalidStateException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomKeyStoreInvalidStateException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomKeyStoreInvalidStateException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoreInvalidStateException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoreInvalidStateException copyWith(
          void Function(CustomKeyStoreInvalidStateException) updates) =>
      super.copyWith((message) =>
              updates(message as CustomKeyStoreInvalidStateException))
          as CustomKeyStoreInvalidStateException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomKeyStoreInvalidStateException create() =>
      CustomKeyStoreInvalidStateException._();
  @$core.override
  CustomKeyStoreInvalidStateException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomKeyStoreInvalidStateException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CustomKeyStoreInvalidStateException>(create);
  static CustomKeyStoreInvalidStateException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CustomKeyStoreNameInUseException extends $pb.GeneratedMessage {
  factory CustomKeyStoreNameInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CustomKeyStoreNameInUseException._();

  factory CustomKeyStoreNameInUseException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomKeyStoreNameInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomKeyStoreNameInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoreNameInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoreNameInUseException copyWith(
          void Function(CustomKeyStoreNameInUseException) updates) =>
      super.copyWith(
              (message) => updates(message as CustomKeyStoreNameInUseException))
          as CustomKeyStoreNameInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomKeyStoreNameInUseException create() =>
      CustomKeyStoreNameInUseException._();
  @$core.override
  CustomKeyStoreNameInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomKeyStoreNameInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CustomKeyStoreNameInUseException>(
          create);
  static CustomKeyStoreNameInUseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CustomKeyStoreNotFoundException extends $pb.GeneratedMessage {
  factory CustomKeyStoreNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CustomKeyStoreNotFoundException._();

  factory CustomKeyStoreNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomKeyStoreNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomKeyStoreNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoreNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoreNotFoundException copyWith(
          void Function(CustomKeyStoreNotFoundException) updates) =>
      super.copyWith(
              (message) => updates(message as CustomKeyStoreNotFoundException))
          as CustomKeyStoreNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomKeyStoreNotFoundException create() =>
      CustomKeyStoreNotFoundException._();
  @$core.override
  CustomKeyStoreNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomKeyStoreNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CustomKeyStoreNotFoundException>(
          create);
  static CustomKeyStoreNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CustomKeyStoresListEntry extends $pb.GeneratedMessage {
  factory CustomKeyStoresListEntry({
    $core.String? trustanchorcertificate,
    $core.String? cloudhsmclusterid,
    $core.String? customkeystoreid,
    $core.String? customkeystorename,
    $core.String? creationdate,
    ConnectionErrorCodeType? connectionerrorcode,
    XksProxyConfigurationType? xksproxyconfiguration,
    ConnectionStateType? connectionstate,
    CustomKeyStoreType? customkeystoretype,
  }) {
    final result = create();
    if (trustanchorcertificate != null)
      result.trustanchorcertificate = trustanchorcertificate;
    if (cloudhsmclusterid != null) result.cloudhsmclusterid = cloudhsmclusterid;
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    if (customkeystorename != null)
      result.customkeystorename = customkeystorename;
    if (creationdate != null) result.creationdate = creationdate;
    if (connectionerrorcode != null)
      result.connectionerrorcode = connectionerrorcode;
    if (xksproxyconfiguration != null)
      result.xksproxyconfiguration = xksproxyconfiguration;
    if (connectionstate != null) result.connectionstate = connectionstate;
    if (customkeystoretype != null)
      result.customkeystoretype = customkeystoretype;
    return result;
  }

  CustomKeyStoresListEntry._();

  factory CustomKeyStoresListEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomKeyStoresListEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomKeyStoresListEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(48354588, _omitFieldNames ? '' : 'trustanchorcertificate')
    ..aOS(56498754, _omitFieldNames ? '' : 'cloudhsmclusterid')
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..aOS(170278046, _omitFieldNames ? '' : 'customkeystorename')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aE<ConnectionErrorCodeType>(
        324951101, _omitFieldNames ? '' : 'connectionerrorcode',
        enumValues: ConnectionErrorCodeType.values)
    ..aOM<XksProxyConfigurationType>(
        349047828, _omitFieldNames ? '' : 'xksproxyconfiguration',
        subBuilder: XksProxyConfigurationType.create)
    ..aE<ConnectionStateType>(
        404323675, _omitFieldNames ? '' : 'connectionstate',
        enumValues: ConnectionStateType.values)
    ..aE<CustomKeyStoreType>(
        415647103, _omitFieldNames ? '' : 'customkeystoretype',
        enumValues: CustomKeyStoreType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoresListEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomKeyStoresListEntry copyWith(
          void Function(CustomKeyStoresListEntry) updates) =>
      super.copyWith((message) => updates(message as CustomKeyStoresListEntry))
          as CustomKeyStoresListEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomKeyStoresListEntry create() => CustomKeyStoresListEntry._();
  @$core.override
  CustomKeyStoresListEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomKeyStoresListEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CustomKeyStoresListEntry>(create);
  static CustomKeyStoresListEntry? _defaultInstance;

  @$pb.TagNumber(48354588)
  $core.String get trustanchorcertificate => $_getSZ(0);
  @$pb.TagNumber(48354588)
  set trustanchorcertificate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(48354588)
  $core.bool hasTrustanchorcertificate() => $_has(0);
  @$pb.TagNumber(48354588)
  void clearTrustanchorcertificate() => $_clearField(48354588);

  @$pb.TagNumber(56498754)
  $core.String get cloudhsmclusterid => $_getSZ(1);
  @$pb.TagNumber(56498754)
  set cloudhsmclusterid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(56498754)
  $core.bool hasCloudhsmclusterid() => $_has(1);
  @$pb.TagNumber(56498754)
  void clearCloudhsmclusterid() => $_clearField(56498754);

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(2);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(2);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);

  @$pb.TagNumber(170278046)
  $core.String get customkeystorename => $_getSZ(3);
  @$pb.TagNumber(170278046)
  set customkeystorename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(170278046)
  $core.bool hasCustomkeystorename() => $_has(3);
  @$pb.TagNumber(170278046)
  void clearCustomkeystorename() => $_clearField(170278046);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(4);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(4);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(324951101)
  ConnectionErrorCodeType get connectionerrorcode => $_getN(5);
  @$pb.TagNumber(324951101)
  set connectionerrorcode(ConnectionErrorCodeType value) =>
      $_setField(324951101, value);
  @$pb.TagNumber(324951101)
  $core.bool hasConnectionerrorcode() => $_has(5);
  @$pb.TagNumber(324951101)
  void clearConnectionerrorcode() => $_clearField(324951101);

  @$pb.TagNumber(349047828)
  XksProxyConfigurationType get xksproxyconfiguration => $_getN(6);
  @$pb.TagNumber(349047828)
  set xksproxyconfiguration(XksProxyConfigurationType value) =>
      $_setField(349047828, value);
  @$pb.TagNumber(349047828)
  $core.bool hasXksproxyconfiguration() => $_has(6);
  @$pb.TagNumber(349047828)
  void clearXksproxyconfiguration() => $_clearField(349047828);
  @$pb.TagNumber(349047828)
  XksProxyConfigurationType ensureXksproxyconfiguration() => $_ensure(6);

  @$pb.TagNumber(404323675)
  ConnectionStateType get connectionstate => $_getN(7);
  @$pb.TagNumber(404323675)
  set connectionstate(ConnectionStateType value) =>
      $_setField(404323675, value);
  @$pb.TagNumber(404323675)
  $core.bool hasConnectionstate() => $_has(7);
  @$pb.TagNumber(404323675)
  void clearConnectionstate() => $_clearField(404323675);

  @$pb.TagNumber(415647103)
  CustomKeyStoreType get customkeystoretype => $_getN(8);
  @$pb.TagNumber(415647103)
  set customkeystoretype(CustomKeyStoreType value) =>
      $_setField(415647103, value);
  @$pb.TagNumber(415647103)
  $core.bool hasCustomkeystoretype() => $_has(8);
  @$pb.TagNumber(415647103)
  void clearCustomkeystoretype() => $_clearField(415647103);
}

class DecryptRequest extends $pb.GeneratedMessage {
  factory DecryptRequest({
    $core.Iterable<DryRunModifierType>? dryrunmodifiers,
    $core.bool? dryrun,
    EncryptionAlgorithmSpec? encryptionalgorithm,
    $core.String? keyid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        encryptioncontext,
    $core.List<$core.int>? ciphertextblob,
    $core.Iterable<$core.String>? granttokens,
    RecipientInfo? recipient,
  }) {
    final result = create();
    if (dryrunmodifiers != null) result.dryrunmodifiers.addAll(dryrunmodifiers);
    if (dryrun != null) result.dryrun = dryrun;
    if (encryptionalgorithm != null)
      result.encryptionalgorithm = encryptionalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (encryptioncontext != null)
      result.encryptioncontext.addEntries(encryptioncontext);
    if (ciphertextblob != null) result.ciphertextblob = ciphertextblob;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    if (recipient != null) result.recipient = recipient;
    return result;
  }

  DecryptRequest._();

  factory DecryptRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DecryptRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DecryptRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..pc<DryRunModifierType>(
        24346424, _omitFieldNames ? '' : 'dryrunmodifiers', $pb.PbFieldType.KE,
        valueOf: DryRunModifierType.valueOf,
        enumValues: DryRunModifierType.values,
        defaultEnumValue:
            DryRunModifierType.DRY_RUN_MODIFIER_TYPE_IGNORE_CIPHERTEXT)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aE<EncryptionAlgorithmSpec>(
        203224586, _omitFieldNames ? '' : 'encryptionalgorithm',
        enumValues: EncryptionAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..m<$core.String, $core.String>(
        286249536, _omitFieldNames ? '' : 'encryptioncontext',
        entryClassName: 'DecryptRequest.EncryptioncontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..a<$core.List<$core.int>>(
        338198183, _omitFieldNames ? '' : 'ciphertextblob', $pb.PbFieldType.OY)
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..aOM<RecipientInfo>(445981721, _omitFieldNames ? '' : 'recipient',
        subBuilder: RecipientInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecryptRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecryptRequest copyWith(void Function(DecryptRequest) updates) =>
      super.copyWith((message) => updates(message as DecryptRequest))
          as DecryptRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DecryptRequest create() => DecryptRequest._();
  @$core.override
  DecryptRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DecryptRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DecryptRequest>(create);
  static DecryptRequest? _defaultInstance;

  @$pb.TagNumber(24346424)
  $pb.PbList<DryRunModifierType> get dryrunmodifiers => $_getList(0);

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(1);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(1);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(203224586)
  EncryptionAlgorithmSpec get encryptionalgorithm => $_getN(2);
  @$pb.TagNumber(203224586)
  set encryptionalgorithm(EncryptionAlgorithmSpec value) =>
      $_setField(203224586, value);
  @$pb.TagNumber(203224586)
  $core.bool hasEncryptionalgorithm() => $_has(2);
  @$pb.TagNumber(203224586)
  void clearEncryptionalgorithm() => $_clearField(203224586);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(286249536)
  $pb.PbMap<$core.String, $core.String> get encryptioncontext => $_getMap(4);

  @$pb.TagNumber(338198183)
  $core.List<$core.int> get ciphertextblob => $_getN(5);
  @$pb.TagNumber(338198183)
  set ciphertextblob($core.List<$core.int> value) => $_setBytes(5, value);
  @$pb.TagNumber(338198183)
  $core.bool hasCiphertextblob() => $_has(5);
  @$pb.TagNumber(338198183)
  void clearCiphertextblob() => $_clearField(338198183);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(6);

  @$pb.TagNumber(445981721)
  RecipientInfo get recipient => $_getN(7);
  @$pb.TagNumber(445981721)
  set recipient(RecipientInfo value) => $_setField(445981721, value);
  @$pb.TagNumber(445981721)
  $core.bool hasRecipient() => $_has(7);
  @$pb.TagNumber(445981721)
  void clearRecipient() => $_clearField(445981721);
  @$pb.TagNumber(445981721)
  RecipientInfo ensureRecipient() => $_ensure(7);
}

class DecryptResponse extends $pb.GeneratedMessage {
  factory DecryptResponse({
    $core.List<$core.int>? ciphertextforrecipient,
    $core.List<$core.int>? plaintext,
    $core.String? keymaterialid,
    EncryptionAlgorithmSpec? encryptionalgorithm,
    $core.String? keyid,
  }) {
    final result = create();
    if (ciphertextforrecipient != null)
      result.ciphertextforrecipient = ciphertextforrecipient;
    if (plaintext != null) result.plaintext = plaintext;
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (encryptionalgorithm != null)
      result.encryptionalgorithm = encryptionalgorithm;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  DecryptResponse._();

  factory DecryptResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DecryptResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DecryptResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(75299212,
        _omitFieldNames ? '' : 'ciphertextforrecipient', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(
        88342721, _omitFieldNames ? '' : 'plaintext', $pb.PbFieldType.OY)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..aE<EncryptionAlgorithmSpec>(
        203224586, _omitFieldNames ? '' : 'encryptionalgorithm',
        enumValues: EncryptionAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecryptResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecryptResponse copyWith(void Function(DecryptResponse) updates) =>
      super.copyWith((message) => updates(message as DecryptResponse))
          as DecryptResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DecryptResponse create() => DecryptResponse._();
  @$core.override
  DecryptResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DecryptResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DecryptResponse>(create);
  static DecryptResponse? _defaultInstance;

  @$pb.TagNumber(75299212)
  $core.List<$core.int> get ciphertextforrecipient => $_getN(0);
  @$pb.TagNumber(75299212)
  set ciphertextforrecipient($core.List<$core.int> value) =>
      $_setBytes(0, value);
  @$pb.TagNumber(75299212)
  $core.bool hasCiphertextforrecipient() => $_has(0);
  @$pb.TagNumber(75299212)
  void clearCiphertextforrecipient() => $_clearField(75299212);

  @$pb.TagNumber(88342721)
  $core.List<$core.int> get plaintext => $_getN(1);
  @$pb.TagNumber(88342721)
  set plaintext($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(88342721)
  $core.bool hasPlaintext() => $_has(1);
  @$pb.TagNumber(88342721)
  void clearPlaintext() => $_clearField(88342721);

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(2);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(2);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(203224586)
  EncryptionAlgorithmSpec get encryptionalgorithm => $_getN(3);
  @$pb.TagNumber(203224586)
  set encryptionalgorithm(EncryptionAlgorithmSpec value) =>
      $_setField(203224586, value);
  @$pb.TagNumber(203224586)
  $core.bool hasEncryptionalgorithm() => $_has(3);
  @$pb.TagNumber(203224586)
  void clearEncryptionalgorithm() => $_clearField(203224586);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(4);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(4);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class DeleteAliasRequest extends $pb.GeneratedMessage {
  factory DeleteAliasRequest({
    $core.String? aliasname,
  }) {
    final result = create();
    if (aliasname != null) result.aliasname = aliasname;
    return result;
  }

  DeleteAliasRequest._();

  factory DeleteAliasRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteAliasRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAliasRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(313250709, _omitFieldNames ? '' : 'aliasname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAliasRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAliasRequest copyWith(void Function(DeleteAliasRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteAliasRequest))
          as DeleteAliasRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAliasRequest create() => DeleteAliasRequest._();
  @$core.override
  DeleteAliasRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteAliasRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteAliasRequest>(create);
  static DeleteAliasRequest? _defaultInstance;

  @$pb.TagNumber(313250709)
  $core.String get aliasname => $_getSZ(0);
  @$pb.TagNumber(313250709)
  set aliasname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(313250709)
  $core.bool hasAliasname() => $_has(0);
  @$pb.TagNumber(313250709)
  void clearAliasname() => $_clearField(313250709);
}

class DeleteCustomKeyStoreRequest extends $pb.GeneratedMessage {
  factory DeleteCustomKeyStoreRequest({
    $core.String? customkeystoreid,
  }) {
    final result = create();
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    return result;
  }

  DeleteCustomKeyStoreRequest._();

  factory DeleteCustomKeyStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteCustomKeyStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteCustomKeyStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCustomKeyStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCustomKeyStoreRequest copyWith(
          void Function(DeleteCustomKeyStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteCustomKeyStoreRequest))
          as DeleteCustomKeyStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteCustomKeyStoreRequest create() =>
      DeleteCustomKeyStoreRequest._();
  @$core.override
  DeleteCustomKeyStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteCustomKeyStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteCustomKeyStoreRequest>(create);
  static DeleteCustomKeyStoreRequest? _defaultInstance;

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(0);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(0);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);
}

class DeleteCustomKeyStoreResponse extends $pb.GeneratedMessage {
  factory DeleteCustomKeyStoreResponse() => create();

  DeleteCustomKeyStoreResponse._();

  factory DeleteCustomKeyStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteCustomKeyStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteCustomKeyStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCustomKeyStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCustomKeyStoreResponse copyWith(
          void Function(DeleteCustomKeyStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteCustomKeyStoreResponse))
          as DeleteCustomKeyStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteCustomKeyStoreResponse create() =>
      DeleteCustomKeyStoreResponse._();
  @$core.override
  DeleteCustomKeyStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteCustomKeyStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteCustomKeyStoreResponse>(create);
  static DeleteCustomKeyStoreResponse? _defaultInstance;
}

class DeleteImportedKeyMaterialRequest extends $pb.GeneratedMessage {
  factory DeleteImportedKeyMaterialRequest({
    $core.String? keymaterialid,
    $core.String? keyid,
  }) {
    final result = create();
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  DeleteImportedKeyMaterialRequest._();

  factory DeleteImportedKeyMaterialRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteImportedKeyMaterialRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteImportedKeyMaterialRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteImportedKeyMaterialRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteImportedKeyMaterialRequest copyWith(
          void Function(DeleteImportedKeyMaterialRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteImportedKeyMaterialRequest))
          as DeleteImportedKeyMaterialRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteImportedKeyMaterialRequest create() =>
      DeleteImportedKeyMaterialRequest._();
  @$core.override
  DeleteImportedKeyMaterialRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteImportedKeyMaterialRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteImportedKeyMaterialRequest>(
          create);
  static DeleteImportedKeyMaterialRequest? _defaultInstance;

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(0);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(0);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class DeleteImportedKeyMaterialResponse extends $pb.GeneratedMessage {
  factory DeleteImportedKeyMaterialResponse({
    $core.String? keymaterialid,
    $core.String? keyid,
  }) {
    final result = create();
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  DeleteImportedKeyMaterialResponse._();

  factory DeleteImportedKeyMaterialResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteImportedKeyMaterialResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteImportedKeyMaterialResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteImportedKeyMaterialResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteImportedKeyMaterialResponse copyWith(
          void Function(DeleteImportedKeyMaterialResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteImportedKeyMaterialResponse))
          as DeleteImportedKeyMaterialResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteImportedKeyMaterialResponse create() =>
      DeleteImportedKeyMaterialResponse._();
  @$core.override
  DeleteImportedKeyMaterialResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteImportedKeyMaterialResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteImportedKeyMaterialResponse>(
          create);
  static DeleteImportedKeyMaterialResponse? _defaultInstance;

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(0);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(0);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class DependencyTimeoutException extends $pb.GeneratedMessage {
  factory DependencyTimeoutException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DependencyTimeoutException._();

  factory DependencyTimeoutException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DependencyTimeoutException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DependencyTimeoutException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DependencyTimeoutException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DependencyTimeoutException copyWith(
          void Function(DependencyTimeoutException) updates) =>
      super.copyWith(
              (message) => updates(message as DependencyTimeoutException))
          as DependencyTimeoutException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DependencyTimeoutException create() => DependencyTimeoutException._();
  @$core.override
  DependencyTimeoutException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DependencyTimeoutException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DependencyTimeoutException>(create);
  static DependencyTimeoutException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DeriveSharedSecretRequest extends $pb.GeneratedMessage {
  factory DeriveSharedSecretRequest({
    $core.bool? dryrun,
    KeyAgreementAlgorithmSpec? keyagreementalgorithm,
    $core.List<$core.int>? publickey,
    $core.String? keyid,
    $core.Iterable<$core.String>? granttokens,
    RecipientInfo? recipient,
  }) {
    final result = create();
    if (dryrun != null) result.dryrun = dryrun;
    if (keyagreementalgorithm != null)
      result.keyagreementalgorithm = keyagreementalgorithm;
    if (publickey != null) result.publickey = publickey;
    if (keyid != null) result.keyid = keyid;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    if (recipient != null) result.recipient = recipient;
    return result;
  }

  DeriveSharedSecretRequest._();

  factory DeriveSharedSecretRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeriveSharedSecretRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeriveSharedSecretRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aE<KeyAgreementAlgorithmSpec>(
        99147702, _omitFieldNames ? '' : 'keyagreementalgorithm',
        enumValues: KeyAgreementAlgorithmSpec.values)
    ..a<$core.List<$core.int>>(
        167335776, _omitFieldNames ? '' : 'publickey', $pb.PbFieldType.OY)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..aOM<RecipientInfo>(445981721, _omitFieldNames ? '' : 'recipient',
        subBuilder: RecipientInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeriveSharedSecretRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeriveSharedSecretRequest copyWith(
          void Function(DeriveSharedSecretRequest) updates) =>
      super.copyWith((message) => updates(message as DeriveSharedSecretRequest))
          as DeriveSharedSecretRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeriveSharedSecretRequest create() => DeriveSharedSecretRequest._();
  @$core.override
  DeriveSharedSecretRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeriveSharedSecretRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeriveSharedSecretRequest>(create);
  static DeriveSharedSecretRequest? _defaultInstance;

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(0);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(0);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(99147702)
  KeyAgreementAlgorithmSpec get keyagreementalgorithm => $_getN(1);
  @$pb.TagNumber(99147702)
  set keyagreementalgorithm(KeyAgreementAlgorithmSpec value) =>
      $_setField(99147702, value);
  @$pb.TagNumber(99147702)
  $core.bool hasKeyagreementalgorithm() => $_has(1);
  @$pb.TagNumber(99147702)
  void clearKeyagreementalgorithm() => $_clearField(99147702);

  @$pb.TagNumber(167335776)
  $core.List<$core.int> get publickey => $_getN(2);
  @$pb.TagNumber(167335776)
  set publickey($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(167335776)
  $core.bool hasPublickey() => $_has(2);
  @$pb.TagNumber(167335776)
  void clearPublickey() => $_clearField(167335776);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(4);

  @$pb.TagNumber(445981721)
  RecipientInfo get recipient => $_getN(5);
  @$pb.TagNumber(445981721)
  set recipient(RecipientInfo value) => $_setField(445981721, value);
  @$pb.TagNumber(445981721)
  $core.bool hasRecipient() => $_has(5);
  @$pb.TagNumber(445981721)
  void clearRecipient() => $_clearField(445981721);
  @$pb.TagNumber(445981721)
  RecipientInfo ensureRecipient() => $_ensure(5);
}

class DeriveSharedSecretResponse extends $pb.GeneratedMessage {
  factory DeriveSharedSecretResponse({
    OriginType? keyorigin,
    $core.List<$core.int>? ciphertextforrecipient,
    KeyAgreementAlgorithmSpec? keyagreementalgorithm,
    $core.String? keyid,
    $core.List<$core.int>? sharedsecret,
  }) {
    final result = create();
    if (keyorigin != null) result.keyorigin = keyorigin;
    if (ciphertextforrecipient != null)
      result.ciphertextforrecipient = ciphertextforrecipient;
    if (keyagreementalgorithm != null)
      result.keyagreementalgorithm = keyagreementalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (sharedsecret != null) result.sharedsecret = sharedsecret;
    return result;
  }

  DeriveSharedSecretResponse._();

  factory DeriveSharedSecretResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeriveSharedSecretResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeriveSharedSecretResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<OriginType>(50912311, _omitFieldNames ? '' : 'keyorigin',
        enumValues: OriginType.values)
    ..a<$core.List<$core.int>>(75299212,
        _omitFieldNames ? '' : 'ciphertextforrecipient', $pb.PbFieldType.OY)
    ..aE<KeyAgreementAlgorithmSpec>(
        99147702, _omitFieldNames ? '' : 'keyagreementalgorithm',
        enumValues: KeyAgreementAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..a<$core.List<$core.int>>(
        382876889, _omitFieldNames ? '' : 'sharedsecret', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeriveSharedSecretResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeriveSharedSecretResponse copyWith(
          void Function(DeriveSharedSecretResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeriveSharedSecretResponse))
          as DeriveSharedSecretResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeriveSharedSecretResponse create() => DeriveSharedSecretResponse._();
  @$core.override
  DeriveSharedSecretResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeriveSharedSecretResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeriveSharedSecretResponse>(create);
  static DeriveSharedSecretResponse? _defaultInstance;

  @$pb.TagNumber(50912311)
  OriginType get keyorigin => $_getN(0);
  @$pb.TagNumber(50912311)
  set keyorigin(OriginType value) => $_setField(50912311, value);
  @$pb.TagNumber(50912311)
  $core.bool hasKeyorigin() => $_has(0);
  @$pb.TagNumber(50912311)
  void clearKeyorigin() => $_clearField(50912311);

  @$pb.TagNumber(75299212)
  $core.List<$core.int> get ciphertextforrecipient => $_getN(1);
  @$pb.TagNumber(75299212)
  set ciphertextforrecipient($core.List<$core.int> value) =>
      $_setBytes(1, value);
  @$pb.TagNumber(75299212)
  $core.bool hasCiphertextforrecipient() => $_has(1);
  @$pb.TagNumber(75299212)
  void clearCiphertextforrecipient() => $_clearField(75299212);

  @$pb.TagNumber(99147702)
  KeyAgreementAlgorithmSpec get keyagreementalgorithm => $_getN(2);
  @$pb.TagNumber(99147702)
  set keyagreementalgorithm(KeyAgreementAlgorithmSpec value) =>
      $_setField(99147702, value);
  @$pb.TagNumber(99147702)
  $core.bool hasKeyagreementalgorithm() => $_has(2);
  @$pb.TagNumber(99147702)
  void clearKeyagreementalgorithm() => $_clearField(99147702);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(382876889)
  $core.List<$core.int> get sharedsecret => $_getN(4);
  @$pb.TagNumber(382876889)
  set sharedsecret($core.List<$core.int> value) => $_setBytes(4, value);
  @$pb.TagNumber(382876889)
  $core.bool hasSharedsecret() => $_has(4);
  @$pb.TagNumber(382876889)
  void clearSharedsecret() => $_clearField(382876889);
}

class DescribeCustomKeyStoresRequest extends $pb.GeneratedMessage {
  factory DescribeCustomKeyStoresRequest({
    $core.String? customkeystoreid,
    $core.String? marker,
    $core.String? customkeystorename,
    $core.int? limit,
  }) {
    final result = create();
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    if (marker != null) result.marker = marker;
    if (customkeystorename != null)
      result.customkeystorename = customkeystorename;
    if (limit != null) result.limit = limit;
    return result;
  }

  DescribeCustomKeyStoresRequest._();

  factory DescribeCustomKeyStoresRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeCustomKeyStoresRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeCustomKeyStoresRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(170278046, _omitFieldNames ? '' : 'customkeystorename')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeCustomKeyStoresRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeCustomKeyStoresRequest copyWith(
          void Function(DescribeCustomKeyStoresRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeCustomKeyStoresRequest))
          as DescribeCustomKeyStoresRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeCustomKeyStoresRequest create() =>
      DescribeCustomKeyStoresRequest._();
  @$core.override
  DescribeCustomKeyStoresRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeCustomKeyStoresRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeCustomKeyStoresRequest>(create);
  static DescribeCustomKeyStoresRequest? _defaultInstance;

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(0);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(0);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(1);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(1);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(170278046)
  $core.String get customkeystorename => $_getSZ(2);
  @$pb.TagNumber(170278046)
  set customkeystorename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(170278046)
  $core.bool hasCustomkeystorename() => $_has(2);
  @$pb.TagNumber(170278046)
  void clearCustomkeystorename() => $_clearField(170278046);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class DescribeCustomKeyStoresResponse extends $pb.GeneratedMessage {
  factory DescribeCustomKeyStoresResponse({
    $core.bool? truncated,
    $core.Iterable<CustomKeyStoresListEntry>? customkeystores,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (truncated != null) result.truncated = truncated;
    if (customkeystores != null) result.customkeystores.addAll(customkeystores);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  DescribeCustomKeyStoresResponse._();

  factory DescribeCustomKeyStoresResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeCustomKeyStoresResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeCustomKeyStoresResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(152451018, _omitFieldNames ? '' : 'truncated')
    ..pPM<CustomKeyStoresListEntry>(
        200763800, _omitFieldNames ? '' : 'customkeystores',
        subBuilder: CustomKeyStoresListEntry.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeCustomKeyStoresResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeCustomKeyStoresResponse copyWith(
          void Function(DescribeCustomKeyStoresResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeCustomKeyStoresResponse))
          as DescribeCustomKeyStoresResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeCustomKeyStoresResponse create() =>
      DescribeCustomKeyStoresResponse._();
  @$core.override
  DescribeCustomKeyStoresResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeCustomKeyStoresResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeCustomKeyStoresResponse>(
          create);
  static DescribeCustomKeyStoresResponse? _defaultInstance;

  @$pb.TagNumber(152451018)
  $core.bool get truncated => $_getBF(0);
  @$pb.TagNumber(152451018)
  set truncated($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(152451018)
  $core.bool hasTruncated() => $_has(0);
  @$pb.TagNumber(152451018)
  void clearTruncated() => $_clearField(152451018);

  @$pb.TagNumber(200763800)
  $pb.PbList<CustomKeyStoresListEntry> get customkeystores => $_getList(1);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class DescribeKeyRequest extends $pb.GeneratedMessage {
  factory DescribeKeyRequest({
    $core.String? keyid,
    $core.Iterable<$core.String>? granttokens,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    return result;
  }

  DescribeKeyRequest._();

  factory DescribeKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeKeyRequest copyWith(void Function(DescribeKeyRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeKeyRequest))
          as DescribeKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeKeyRequest create() => DescribeKeyRequest._();
  @$core.override
  DescribeKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeKeyRequest>(create);
  static DescribeKeyRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(1);
}

class DescribeKeyResponse extends $pb.GeneratedMessage {
  factory DescribeKeyResponse({
    KeyMetadata? keymetadata,
  }) {
    final result = create();
    if (keymetadata != null) result.keymetadata = keymetadata;
    return result;
  }

  DescribeKeyResponse._();

  factory DescribeKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOM<KeyMetadata>(149520202, _omitFieldNames ? '' : 'keymetadata',
        subBuilder: KeyMetadata.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeKeyResponse copyWith(void Function(DescribeKeyResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeKeyResponse))
          as DescribeKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeKeyResponse create() => DescribeKeyResponse._();
  @$core.override
  DescribeKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeKeyResponse>(create);
  static DescribeKeyResponse? _defaultInstance;

  @$pb.TagNumber(149520202)
  KeyMetadata get keymetadata => $_getN(0);
  @$pb.TagNumber(149520202)
  set keymetadata(KeyMetadata value) => $_setField(149520202, value);
  @$pb.TagNumber(149520202)
  $core.bool hasKeymetadata() => $_has(0);
  @$pb.TagNumber(149520202)
  void clearKeymetadata() => $_clearField(149520202);
  @$pb.TagNumber(149520202)
  KeyMetadata ensureKeymetadata() => $_ensure(0);
}

class DisableKeyRequest extends $pb.GeneratedMessage {
  factory DisableKeyRequest({
    $core.String? keyid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  DisableKeyRequest._();

  factory DisableKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableKeyRequest copyWith(void Function(DisableKeyRequest) updates) =>
      super.copyWith((message) => updates(message as DisableKeyRequest))
          as DisableKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableKeyRequest create() => DisableKeyRequest._();
  @$core.override
  DisableKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableKeyRequest>(create);
  static DisableKeyRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class DisableKeyRotationRequest extends $pb.GeneratedMessage {
  factory DisableKeyRotationRequest({
    $core.String? keyid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  DisableKeyRotationRequest._();

  factory DisableKeyRotationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableKeyRotationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableKeyRotationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableKeyRotationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableKeyRotationRequest copyWith(
          void Function(DisableKeyRotationRequest) updates) =>
      super.copyWith((message) => updates(message as DisableKeyRotationRequest))
          as DisableKeyRotationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableKeyRotationRequest create() => DisableKeyRotationRequest._();
  @$core.override
  DisableKeyRotationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableKeyRotationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableKeyRotationRequest>(create);
  static DisableKeyRotationRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class DisabledException extends $pb.GeneratedMessage {
  factory DisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DisabledException._();

  factory DisabledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisabledException copyWith(void Function(DisabledException) updates) =>
      super.copyWith((message) => updates(message as DisabledException))
          as DisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisabledException create() => DisabledException._();
  @$core.override
  DisabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisabledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisabledException>(create);
  static DisabledException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DisconnectCustomKeyStoreRequest extends $pb.GeneratedMessage {
  factory DisconnectCustomKeyStoreRequest({
    $core.String? customkeystoreid,
  }) {
    final result = create();
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    return result;
  }

  DisconnectCustomKeyStoreRequest._();

  factory DisconnectCustomKeyStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisconnectCustomKeyStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisconnectCustomKeyStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisconnectCustomKeyStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisconnectCustomKeyStoreRequest copyWith(
          void Function(DisconnectCustomKeyStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DisconnectCustomKeyStoreRequest))
          as DisconnectCustomKeyStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisconnectCustomKeyStoreRequest create() =>
      DisconnectCustomKeyStoreRequest._();
  @$core.override
  DisconnectCustomKeyStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisconnectCustomKeyStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisconnectCustomKeyStoreRequest>(
          create);
  static DisconnectCustomKeyStoreRequest? _defaultInstance;

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(0);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(0);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);
}

class DisconnectCustomKeyStoreResponse extends $pb.GeneratedMessage {
  factory DisconnectCustomKeyStoreResponse() => create();

  DisconnectCustomKeyStoreResponse._();

  factory DisconnectCustomKeyStoreResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisconnectCustomKeyStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisconnectCustomKeyStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisconnectCustomKeyStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisconnectCustomKeyStoreResponse copyWith(
          void Function(DisconnectCustomKeyStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DisconnectCustomKeyStoreResponse))
          as DisconnectCustomKeyStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisconnectCustomKeyStoreResponse create() =>
      DisconnectCustomKeyStoreResponse._();
  @$core.override
  DisconnectCustomKeyStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisconnectCustomKeyStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisconnectCustomKeyStoreResponse>(
          create);
  static DisconnectCustomKeyStoreResponse? _defaultInstance;
}

class DryRunOperationException extends $pb.GeneratedMessage {
  factory DryRunOperationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DryRunOperationException._();

  factory DryRunOperationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DryRunOperationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DryRunOperationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DryRunOperationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DryRunOperationException copyWith(
          void Function(DryRunOperationException) updates) =>
      super.copyWith((message) => updates(message as DryRunOperationException))
          as DryRunOperationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DryRunOperationException create() => DryRunOperationException._();
  @$core.override
  DryRunOperationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DryRunOperationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DryRunOperationException>(create);
  static DryRunOperationException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class EnableKeyRequest extends $pb.GeneratedMessage {
  factory EnableKeyRequest({
    $core.String? keyid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  EnableKeyRequest._();

  factory EnableKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableKeyRequest copyWith(void Function(EnableKeyRequest) updates) =>
      super.copyWith((message) => updates(message as EnableKeyRequest))
          as EnableKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableKeyRequest create() => EnableKeyRequest._();
  @$core.override
  EnableKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableKeyRequest>(create);
  static EnableKeyRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class EnableKeyRotationRequest extends $pb.GeneratedMessage {
  factory EnableKeyRotationRequest({
    $core.int? rotationperiodindays,
    $core.String? keyid,
  }) {
    final result = create();
    if (rotationperiodindays != null)
      result.rotationperiodindays = rotationperiodindays;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  EnableKeyRotationRequest._();

  factory EnableKeyRotationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableKeyRotationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableKeyRotationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aI(118357231, _omitFieldNames ? '' : 'rotationperiodindays')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableKeyRotationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableKeyRotationRequest copyWith(
          void Function(EnableKeyRotationRequest) updates) =>
      super.copyWith((message) => updates(message as EnableKeyRotationRequest))
          as EnableKeyRotationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableKeyRotationRequest create() => EnableKeyRotationRequest._();
  @$core.override
  EnableKeyRotationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableKeyRotationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableKeyRotationRequest>(create);
  static EnableKeyRotationRequest? _defaultInstance;

  @$pb.TagNumber(118357231)
  $core.int get rotationperiodindays => $_getIZ(0);
  @$pb.TagNumber(118357231)
  set rotationperiodindays($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(118357231)
  $core.bool hasRotationperiodindays() => $_has(0);
  @$pb.TagNumber(118357231)
  void clearRotationperiodindays() => $_clearField(118357231);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class EncryptRequest extends $pb.GeneratedMessage {
  factory EncryptRequest({
    $core.List<$core.int>? plaintext,
    $core.bool? dryrun,
    EncryptionAlgorithmSpec? encryptionalgorithm,
    $core.String? keyid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        encryptioncontext,
    $core.Iterable<$core.String>? granttokens,
  }) {
    final result = create();
    if (plaintext != null) result.plaintext = plaintext;
    if (dryrun != null) result.dryrun = dryrun;
    if (encryptionalgorithm != null)
      result.encryptionalgorithm = encryptionalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (encryptioncontext != null)
      result.encryptioncontext.addEntries(encryptioncontext);
    if (granttokens != null) result.granttokens.addAll(granttokens);
    return result;
  }

  EncryptRequest._();

  factory EncryptRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EncryptRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EncryptRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(
        88342721, _omitFieldNames ? '' : 'plaintext', $pb.PbFieldType.OY)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aE<EncryptionAlgorithmSpec>(
        203224586, _omitFieldNames ? '' : 'encryptionalgorithm',
        enumValues: EncryptionAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..m<$core.String, $core.String>(
        286249536, _omitFieldNames ? '' : 'encryptioncontext',
        entryClassName: 'EncryptRequest.EncryptioncontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptRequest copyWith(void Function(EncryptRequest) updates) =>
      super.copyWith((message) => updates(message as EncryptRequest))
          as EncryptRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EncryptRequest create() => EncryptRequest._();
  @$core.override
  EncryptRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EncryptRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EncryptRequest>(create);
  static EncryptRequest? _defaultInstance;

  @$pb.TagNumber(88342721)
  $core.List<$core.int> get plaintext => $_getN(0);
  @$pb.TagNumber(88342721)
  set plaintext($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(88342721)
  $core.bool hasPlaintext() => $_has(0);
  @$pb.TagNumber(88342721)
  void clearPlaintext() => $_clearField(88342721);

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(1);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(1);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(203224586)
  EncryptionAlgorithmSpec get encryptionalgorithm => $_getN(2);
  @$pb.TagNumber(203224586)
  set encryptionalgorithm(EncryptionAlgorithmSpec value) =>
      $_setField(203224586, value);
  @$pb.TagNumber(203224586)
  $core.bool hasEncryptionalgorithm() => $_has(2);
  @$pb.TagNumber(203224586)
  void clearEncryptionalgorithm() => $_clearField(203224586);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(286249536)
  $pb.PbMap<$core.String, $core.String> get encryptioncontext => $_getMap(4);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(5);
}

class EncryptResponse extends $pb.GeneratedMessage {
  factory EncryptResponse({
    EncryptionAlgorithmSpec? encryptionalgorithm,
    $core.String? keyid,
    $core.List<$core.int>? ciphertextblob,
  }) {
    final result = create();
    if (encryptionalgorithm != null)
      result.encryptionalgorithm = encryptionalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (ciphertextblob != null) result.ciphertextblob = ciphertextblob;
    return result;
  }

  EncryptResponse._();

  factory EncryptResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EncryptResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EncryptResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<EncryptionAlgorithmSpec>(
        203224586, _omitFieldNames ? '' : 'encryptionalgorithm',
        enumValues: EncryptionAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..a<$core.List<$core.int>>(
        338198183, _omitFieldNames ? '' : 'ciphertextblob', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptResponse copyWith(void Function(EncryptResponse) updates) =>
      super.copyWith((message) => updates(message as EncryptResponse))
          as EncryptResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EncryptResponse create() => EncryptResponse._();
  @$core.override
  EncryptResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EncryptResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EncryptResponse>(create);
  static EncryptResponse? _defaultInstance;

  @$pb.TagNumber(203224586)
  EncryptionAlgorithmSpec get encryptionalgorithm => $_getN(0);
  @$pb.TagNumber(203224586)
  set encryptionalgorithm(EncryptionAlgorithmSpec value) =>
      $_setField(203224586, value);
  @$pb.TagNumber(203224586)
  $core.bool hasEncryptionalgorithm() => $_has(0);
  @$pb.TagNumber(203224586)
  void clearEncryptionalgorithm() => $_clearField(203224586);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(338198183)
  $core.List<$core.int> get ciphertextblob => $_getN(2);
  @$pb.TagNumber(338198183)
  set ciphertextblob($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(338198183)
  $core.bool hasCiphertextblob() => $_has(2);
  @$pb.TagNumber(338198183)
  void clearCiphertextblob() => $_clearField(338198183);
}

class ExpiredImportTokenException extends $pb.GeneratedMessage {
  factory ExpiredImportTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExpiredImportTokenException._();

  factory ExpiredImportTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExpiredImportTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExpiredImportTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredImportTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpiredImportTokenException copyWith(
          void Function(ExpiredImportTokenException) updates) =>
      super.copyWith(
              (message) => updates(message as ExpiredImportTokenException))
          as ExpiredImportTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExpiredImportTokenException create() =>
      ExpiredImportTokenException._();
  @$core.override
  ExpiredImportTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExpiredImportTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExpiredImportTokenException>(create);
  static ExpiredImportTokenException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class GenerateDataKeyPairRequest extends $pb.GeneratedMessage {
  factory GenerateDataKeyPairRequest({
    $core.bool? dryrun,
    DataKeyPairSpec? keypairspec,
    $core.String? keyid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        encryptioncontext,
    $core.Iterable<$core.String>? granttokens,
    RecipientInfo? recipient,
  }) {
    final result = create();
    if (dryrun != null) result.dryrun = dryrun;
    if (keypairspec != null) result.keypairspec = keypairspec;
    if (keyid != null) result.keyid = keyid;
    if (encryptioncontext != null)
      result.encryptioncontext.addEntries(encryptioncontext);
    if (granttokens != null) result.granttokens.addAll(granttokens);
    if (recipient != null) result.recipient = recipient;
    return result;
  }

  GenerateDataKeyPairRequest._();

  factory GenerateDataKeyPairRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateDataKeyPairRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateDataKeyPairRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aE<DataKeyPairSpec>(142696380, _omitFieldNames ? '' : 'keypairspec',
        enumValues: DataKeyPairSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..m<$core.String, $core.String>(
        286249536, _omitFieldNames ? '' : 'encryptioncontext',
        entryClassName: 'GenerateDataKeyPairRequest.EncryptioncontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..aOM<RecipientInfo>(445981721, _omitFieldNames ? '' : 'recipient',
        subBuilder: RecipientInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyPairRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyPairRequest copyWith(
          void Function(GenerateDataKeyPairRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GenerateDataKeyPairRequest))
          as GenerateDataKeyPairRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyPairRequest create() => GenerateDataKeyPairRequest._();
  @$core.override
  GenerateDataKeyPairRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyPairRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateDataKeyPairRequest>(create);
  static GenerateDataKeyPairRequest? _defaultInstance;

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(0);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(0);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(142696380)
  DataKeyPairSpec get keypairspec => $_getN(1);
  @$pb.TagNumber(142696380)
  set keypairspec(DataKeyPairSpec value) => $_setField(142696380, value);
  @$pb.TagNumber(142696380)
  $core.bool hasKeypairspec() => $_has(1);
  @$pb.TagNumber(142696380)
  void clearKeypairspec() => $_clearField(142696380);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(286249536)
  $pb.PbMap<$core.String, $core.String> get encryptioncontext => $_getMap(3);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(4);

  @$pb.TagNumber(445981721)
  RecipientInfo get recipient => $_getN(5);
  @$pb.TagNumber(445981721)
  set recipient(RecipientInfo value) => $_setField(445981721, value);
  @$pb.TagNumber(445981721)
  $core.bool hasRecipient() => $_has(5);
  @$pb.TagNumber(445981721)
  void clearRecipient() => $_clearField(445981721);
  @$pb.TagNumber(445981721)
  RecipientInfo ensureRecipient() => $_ensure(5);
}

class GenerateDataKeyPairResponse extends $pb.GeneratedMessage {
  factory GenerateDataKeyPairResponse({
    $core.List<$core.int>? ciphertextforrecipient,
    DataKeyPairSpec? keypairspec,
    $core.String? keymaterialid,
    $core.List<$core.int>? publickey,
    $core.String? keyid,
    $core.List<$core.int>? privatekeyciphertextblob,
    $core.List<$core.int>? privatekeyplaintext,
  }) {
    final result = create();
    if (ciphertextforrecipient != null)
      result.ciphertextforrecipient = ciphertextforrecipient;
    if (keypairspec != null) result.keypairspec = keypairspec;
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (publickey != null) result.publickey = publickey;
    if (keyid != null) result.keyid = keyid;
    if (privatekeyciphertextblob != null)
      result.privatekeyciphertextblob = privatekeyciphertextblob;
    if (privatekeyplaintext != null)
      result.privatekeyplaintext = privatekeyplaintext;
    return result;
  }

  GenerateDataKeyPairResponse._();

  factory GenerateDataKeyPairResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateDataKeyPairResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateDataKeyPairResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(75299212,
        _omitFieldNames ? '' : 'ciphertextforrecipient', $pb.PbFieldType.OY)
    ..aE<DataKeyPairSpec>(142696380, _omitFieldNames ? '' : 'keypairspec',
        enumValues: DataKeyPairSpec.values)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..a<$core.List<$core.int>>(
        167335776, _omitFieldNames ? '' : 'publickey', $pb.PbFieldType.OY)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..a<$core.List<$core.int>>(295238401,
        _omitFieldNames ? '' : 'privatekeyciphertextblob', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(422534247,
        _omitFieldNames ? '' : 'privatekeyplaintext', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyPairResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyPairResponse copyWith(
          void Function(GenerateDataKeyPairResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GenerateDataKeyPairResponse))
          as GenerateDataKeyPairResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyPairResponse create() =>
      GenerateDataKeyPairResponse._();
  @$core.override
  GenerateDataKeyPairResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyPairResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateDataKeyPairResponse>(create);
  static GenerateDataKeyPairResponse? _defaultInstance;

  @$pb.TagNumber(75299212)
  $core.List<$core.int> get ciphertextforrecipient => $_getN(0);
  @$pb.TagNumber(75299212)
  set ciphertextforrecipient($core.List<$core.int> value) =>
      $_setBytes(0, value);
  @$pb.TagNumber(75299212)
  $core.bool hasCiphertextforrecipient() => $_has(0);
  @$pb.TagNumber(75299212)
  void clearCiphertextforrecipient() => $_clearField(75299212);

  @$pb.TagNumber(142696380)
  DataKeyPairSpec get keypairspec => $_getN(1);
  @$pb.TagNumber(142696380)
  set keypairspec(DataKeyPairSpec value) => $_setField(142696380, value);
  @$pb.TagNumber(142696380)
  $core.bool hasKeypairspec() => $_has(1);
  @$pb.TagNumber(142696380)
  void clearKeypairspec() => $_clearField(142696380);

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(2);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(2);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(167335776)
  $core.List<$core.int> get publickey => $_getN(3);
  @$pb.TagNumber(167335776)
  set publickey($core.List<$core.int> value) => $_setBytes(3, value);
  @$pb.TagNumber(167335776)
  $core.bool hasPublickey() => $_has(3);
  @$pb.TagNumber(167335776)
  void clearPublickey() => $_clearField(167335776);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(4);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(4);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(295238401)
  $core.List<$core.int> get privatekeyciphertextblob => $_getN(5);
  @$pb.TagNumber(295238401)
  set privatekeyciphertextblob($core.List<$core.int> value) =>
      $_setBytes(5, value);
  @$pb.TagNumber(295238401)
  $core.bool hasPrivatekeyciphertextblob() => $_has(5);
  @$pb.TagNumber(295238401)
  void clearPrivatekeyciphertextblob() => $_clearField(295238401);

  @$pb.TagNumber(422534247)
  $core.List<$core.int> get privatekeyplaintext => $_getN(6);
  @$pb.TagNumber(422534247)
  set privatekeyplaintext($core.List<$core.int> value) => $_setBytes(6, value);
  @$pb.TagNumber(422534247)
  $core.bool hasPrivatekeyplaintext() => $_has(6);
  @$pb.TagNumber(422534247)
  void clearPrivatekeyplaintext() => $_clearField(422534247);
}

class GenerateDataKeyPairWithoutPlaintextRequest extends $pb.GeneratedMessage {
  factory GenerateDataKeyPairWithoutPlaintextRequest({
    $core.bool? dryrun,
    DataKeyPairSpec? keypairspec,
    $core.String? keyid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        encryptioncontext,
    $core.Iterable<$core.String>? granttokens,
  }) {
    final result = create();
    if (dryrun != null) result.dryrun = dryrun;
    if (keypairspec != null) result.keypairspec = keypairspec;
    if (keyid != null) result.keyid = keyid;
    if (encryptioncontext != null)
      result.encryptioncontext.addEntries(encryptioncontext);
    if (granttokens != null) result.granttokens.addAll(granttokens);
    return result;
  }

  GenerateDataKeyPairWithoutPlaintextRequest._();

  factory GenerateDataKeyPairWithoutPlaintextRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateDataKeyPairWithoutPlaintextRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateDataKeyPairWithoutPlaintextRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aE<DataKeyPairSpec>(142696380, _omitFieldNames ? '' : 'keypairspec',
        enumValues: DataKeyPairSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..m<$core.String, $core.String>(
        286249536, _omitFieldNames ? '' : 'encryptioncontext',
        entryClassName:
            'GenerateDataKeyPairWithoutPlaintextRequest.EncryptioncontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyPairWithoutPlaintextRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyPairWithoutPlaintextRequest copyWith(
          void Function(GenerateDataKeyPairWithoutPlaintextRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GenerateDataKeyPairWithoutPlaintextRequest))
          as GenerateDataKeyPairWithoutPlaintextRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyPairWithoutPlaintextRequest create() =>
      GenerateDataKeyPairWithoutPlaintextRequest._();
  @$core.override
  GenerateDataKeyPairWithoutPlaintextRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyPairWithoutPlaintextRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GenerateDataKeyPairWithoutPlaintextRequest>(create);
  static GenerateDataKeyPairWithoutPlaintextRequest? _defaultInstance;

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(0);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(0);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(142696380)
  DataKeyPairSpec get keypairspec => $_getN(1);
  @$pb.TagNumber(142696380)
  set keypairspec(DataKeyPairSpec value) => $_setField(142696380, value);
  @$pb.TagNumber(142696380)
  $core.bool hasKeypairspec() => $_has(1);
  @$pb.TagNumber(142696380)
  void clearKeypairspec() => $_clearField(142696380);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(286249536)
  $pb.PbMap<$core.String, $core.String> get encryptioncontext => $_getMap(3);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(4);
}

class GenerateDataKeyPairWithoutPlaintextResponse extends $pb.GeneratedMessage {
  factory GenerateDataKeyPairWithoutPlaintextResponse({
    DataKeyPairSpec? keypairspec,
    $core.String? keymaterialid,
    $core.List<$core.int>? publickey,
    $core.String? keyid,
    $core.List<$core.int>? privatekeyciphertextblob,
  }) {
    final result = create();
    if (keypairspec != null) result.keypairspec = keypairspec;
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (publickey != null) result.publickey = publickey;
    if (keyid != null) result.keyid = keyid;
    if (privatekeyciphertextblob != null)
      result.privatekeyciphertextblob = privatekeyciphertextblob;
    return result;
  }

  GenerateDataKeyPairWithoutPlaintextResponse._();

  factory GenerateDataKeyPairWithoutPlaintextResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateDataKeyPairWithoutPlaintextResponse.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateDataKeyPairWithoutPlaintextResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<DataKeyPairSpec>(142696380, _omitFieldNames ? '' : 'keypairspec',
        enumValues: DataKeyPairSpec.values)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..a<$core.List<$core.int>>(
        167335776, _omitFieldNames ? '' : 'publickey', $pb.PbFieldType.OY)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..a<$core.List<$core.int>>(295238401,
        _omitFieldNames ? '' : 'privatekeyciphertextblob', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyPairWithoutPlaintextResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyPairWithoutPlaintextResponse copyWith(
          void Function(GenerateDataKeyPairWithoutPlaintextResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GenerateDataKeyPairWithoutPlaintextResponse))
          as GenerateDataKeyPairWithoutPlaintextResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyPairWithoutPlaintextResponse create() =>
      GenerateDataKeyPairWithoutPlaintextResponse._();
  @$core.override
  GenerateDataKeyPairWithoutPlaintextResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyPairWithoutPlaintextResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GenerateDataKeyPairWithoutPlaintextResponse>(create);
  static GenerateDataKeyPairWithoutPlaintextResponse? _defaultInstance;

  @$pb.TagNumber(142696380)
  DataKeyPairSpec get keypairspec => $_getN(0);
  @$pb.TagNumber(142696380)
  set keypairspec(DataKeyPairSpec value) => $_setField(142696380, value);
  @$pb.TagNumber(142696380)
  $core.bool hasKeypairspec() => $_has(0);
  @$pb.TagNumber(142696380)
  void clearKeypairspec() => $_clearField(142696380);

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(1);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(1);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(167335776)
  $core.List<$core.int> get publickey => $_getN(2);
  @$pb.TagNumber(167335776)
  set publickey($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(167335776)
  $core.bool hasPublickey() => $_has(2);
  @$pb.TagNumber(167335776)
  void clearPublickey() => $_clearField(167335776);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(295238401)
  $core.List<$core.int> get privatekeyciphertextblob => $_getN(4);
  @$pb.TagNumber(295238401)
  set privatekeyciphertextblob($core.List<$core.int> value) =>
      $_setBytes(4, value);
  @$pb.TagNumber(295238401)
  $core.bool hasPrivatekeyciphertextblob() => $_has(4);
  @$pb.TagNumber(295238401)
  void clearPrivatekeyciphertextblob() => $_clearField(295238401);
}

class GenerateDataKeyRequest extends $pb.GeneratedMessage {
  factory GenerateDataKeyRequest({
    $core.bool? dryrun,
    DataKeySpec? keyspec,
    $core.String? keyid,
    $core.int? numberofbytes,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        encryptioncontext,
    $core.Iterable<$core.String>? granttokens,
    RecipientInfo? recipient,
  }) {
    final result = create();
    if (dryrun != null) result.dryrun = dryrun;
    if (keyspec != null) result.keyspec = keyspec;
    if (keyid != null) result.keyid = keyid;
    if (numberofbytes != null) result.numberofbytes = numberofbytes;
    if (encryptioncontext != null)
      result.encryptioncontext.addEntries(encryptioncontext);
    if (granttokens != null) result.granttokens.addAll(granttokens);
    if (recipient != null) result.recipient = recipient;
    return result;
  }

  GenerateDataKeyRequest._();

  factory GenerateDataKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateDataKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateDataKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aE<DataKeySpec>(138220928, _omitFieldNames ? '' : 'keyspec',
        enumValues: DataKeySpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(277086201, _omitFieldNames ? '' : 'numberofbytes')
    ..m<$core.String, $core.String>(
        286249536, _omitFieldNames ? '' : 'encryptioncontext',
        entryClassName: 'GenerateDataKeyRequest.EncryptioncontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..aOM<RecipientInfo>(445981721, _omitFieldNames ? '' : 'recipient',
        subBuilder: RecipientInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyRequest copyWith(
          void Function(GenerateDataKeyRequest) updates) =>
      super.copyWith((message) => updates(message as GenerateDataKeyRequest))
          as GenerateDataKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyRequest create() => GenerateDataKeyRequest._();
  @$core.override
  GenerateDataKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateDataKeyRequest>(create);
  static GenerateDataKeyRequest? _defaultInstance;

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(0);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(0);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(138220928)
  DataKeySpec get keyspec => $_getN(1);
  @$pb.TagNumber(138220928)
  set keyspec(DataKeySpec value) => $_setField(138220928, value);
  @$pb.TagNumber(138220928)
  $core.bool hasKeyspec() => $_has(1);
  @$pb.TagNumber(138220928)
  void clearKeyspec() => $_clearField(138220928);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(277086201)
  $core.int get numberofbytes => $_getIZ(3);
  @$pb.TagNumber(277086201)
  set numberofbytes($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(277086201)
  $core.bool hasNumberofbytes() => $_has(3);
  @$pb.TagNumber(277086201)
  void clearNumberofbytes() => $_clearField(277086201);

  @$pb.TagNumber(286249536)
  $pb.PbMap<$core.String, $core.String> get encryptioncontext => $_getMap(4);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(5);

  @$pb.TagNumber(445981721)
  RecipientInfo get recipient => $_getN(6);
  @$pb.TagNumber(445981721)
  set recipient(RecipientInfo value) => $_setField(445981721, value);
  @$pb.TagNumber(445981721)
  $core.bool hasRecipient() => $_has(6);
  @$pb.TagNumber(445981721)
  void clearRecipient() => $_clearField(445981721);
  @$pb.TagNumber(445981721)
  RecipientInfo ensureRecipient() => $_ensure(6);
}

class GenerateDataKeyResponse extends $pb.GeneratedMessage {
  factory GenerateDataKeyResponse({
    $core.List<$core.int>? ciphertextforrecipient,
    $core.List<$core.int>? plaintext,
    $core.String? keymaterialid,
    $core.String? keyid,
    $core.List<$core.int>? ciphertextblob,
  }) {
    final result = create();
    if (ciphertextforrecipient != null)
      result.ciphertextforrecipient = ciphertextforrecipient;
    if (plaintext != null) result.plaintext = plaintext;
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (keyid != null) result.keyid = keyid;
    if (ciphertextblob != null) result.ciphertextblob = ciphertextblob;
    return result;
  }

  GenerateDataKeyResponse._();

  factory GenerateDataKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateDataKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateDataKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(75299212,
        _omitFieldNames ? '' : 'ciphertextforrecipient', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(
        88342721, _omitFieldNames ? '' : 'plaintext', $pb.PbFieldType.OY)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..a<$core.List<$core.int>>(
        338198183, _omitFieldNames ? '' : 'ciphertextblob', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyResponse copyWith(
          void Function(GenerateDataKeyResponse) updates) =>
      super.copyWith((message) => updates(message as GenerateDataKeyResponse))
          as GenerateDataKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyResponse create() => GenerateDataKeyResponse._();
  @$core.override
  GenerateDataKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateDataKeyResponse>(create);
  static GenerateDataKeyResponse? _defaultInstance;

  @$pb.TagNumber(75299212)
  $core.List<$core.int> get ciphertextforrecipient => $_getN(0);
  @$pb.TagNumber(75299212)
  set ciphertextforrecipient($core.List<$core.int> value) =>
      $_setBytes(0, value);
  @$pb.TagNumber(75299212)
  $core.bool hasCiphertextforrecipient() => $_has(0);
  @$pb.TagNumber(75299212)
  void clearCiphertextforrecipient() => $_clearField(75299212);

  @$pb.TagNumber(88342721)
  $core.List<$core.int> get plaintext => $_getN(1);
  @$pb.TagNumber(88342721)
  set plaintext($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(88342721)
  $core.bool hasPlaintext() => $_has(1);
  @$pb.TagNumber(88342721)
  void clearPlaintext() => $_clearField(88342721);

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(2);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(2);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(338198183)
  $core.List<$core.int> get ciphertextblob => $_getN(4);
  @$pb.TagNumber(338198183)
  set ciphertextblob($core.List<$core.int> value) => $_setBytes(4, value);
  @$pb.TagNumber(338198183)
  $core.bool hasCiphertextblob() => $_has(4);
  @$pb.TagNumber(338198183)
  void clearCiphertextblob() => $_clearField(338198183);
}

class GenerateDataKeyWithoutPlaintextRequest extends $pb.GeneratedMessage {
  factory GenerateDataKeyWithoutPlaintextRequest({
    $core.bool? dryrun,
    DataKeySpec? keyspec,
    $core.String? keyid,
    $core.int? numberofbytes,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        encryptioncontext,
    $core.Iterable<$core.String>? granttokens,
  }) {
    final result = create();
    if (dryrun != null) result.dryrun = dryrun;
    if (keyspec != null) result.keyspec = keyspec;
    if (keyid != null) result.keyid = keyid;
    if (numberofbytes != null) result.numberofbytes = numberofbytes;
    if (encryptioncontext != null)
      result.encryptioncontext.addEntries(encryptioncontext);
    if (granttokens != null) result.granttokens.addAll(granttokens);
    return result;
  }

  GenerateDataKeyWithoutPlaintextRequest._();

  factory GenerateDataKeyWithoutPlaintextRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateDataKeyWithoutPlaintextRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateDataKeyWithoutPlaintextRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aE<DataKeySpec>(138220928, _omitFieldNames ? '' : 'keyspec',
        enumValues: DataKeySpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(277086201, _omitFieldNames ? '' : 'numberofbytes')
    ..m<$core.String, $core.String>(
        286249536, _omitFieldNames ? '' : 'encryptioncontext',
        entryClassName:
            'GenerateDataKeyWithoutPlaintextRequest.EncryptioncontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyWithoutPlaintextRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyWithoutPlaintextRequest copyWith(
          void Function(GenerateDataKeyWithoutPlaintextRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GenerateDataKeyWithoutPlaintextRequest))
          as GenerateDataKeyWithoutPlaintextRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyWithoutPlaintextRequest create() =>
      GenerateDataKeyWithoutPlaintextRequest._();
  @$core.override
  GenerateDataKeyWithoutPlaintextRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyWithoutPlaintextRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GenerateDataKeyWithoutPlaintextRequest>(create);
  static GenerateDataKeyWithoutPlaintextRequest? _defaultInstance;

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(0);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(0);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(138220928)
  DataKeySpec get keyspec => $_getN(1);
  @$pb.TagNumber(138220928)
  set keyspec(DataKeySpec value) => $_setField(138220928, value);
  @$pb.TagNumber(138220928)
  $core.bool hasKeyspec() => $_has(1);
  @$pb.TagNumber(138220928)
  void clearKeyspec() => $_clearField(138220928);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(277086201)
  $core.int get numberofbytes => $_getIZ(3);
  @$pb.TagNumber(277086201)
  set numberofbytes($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(277086201)
  $core.bool hasNumberofbytes() => $_has(3);
  @$pb.TagNumber(277086201)
  void clearNumberofbytes() => $_clearField(277086201);

  @$pb.TagNumber(286249536)
  $pb.PbMap<$core.String, $core.String> get encryptioncontext => $_getMap(4);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(5);
}

class GenerateDataKeyWithoutPlaintextResponse extends $pb.GeneratedMessage {
  factory GenerateDataKeyWithoutPlaintextResponse({
    $core.String? keymaterialid,
    $core.String? keyid,
    $core.List<$core.int>? ciphertextblob,
  }) {
    final result = create();
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (keyid != null) result.keyid = keyid;
    if (ciphertextblob != null) result.ciphertextblob = ciphertextblob;
    return result;
  }

  GenerateDataKeyWithoutPlaintextResponse._();

  factory GenerateDataKeyWithoutPlaintextResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateDataKeyWithoutPlaintextResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateDataKeyWithoutPlaintextResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..a<$core.List<$core.int>>(
        338198183, _omitFieldNames ? '' : 'ciphertextblob', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyWithoutPlaintextResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateDataKeyWithoutPlaintextResponse copyWith(
          void Function(GenerateDataKeyWithoutPlaintextResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GenerateDataKeyWithoutPlaintextResponse))
          as GenerateDataKeyWithoutPlaintextResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyWithoutPlaintextResponse create() =>
      GenerateDataKeyWithoutPlaintextResponse._();
  @$core.override
  GenerateDataKeyWithoutPlaintextResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateDataKeyWithoutPlaintextResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GenerateDataKeyWithoutPlaintextResponse>(create);
  static GenerateDataKeyWithoutPlaintextResponse? _defaultInstance;

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(0);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(0);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(338198183)
  $core.List<$core.int> get ciphertextblob => $_getN(2);
  @$pb.TagNumber(338198183)
  set ciphertextblob($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(338198183)
  $core.bool hasCiphertextblob() => $_has(2);
  @$pb.TagNumber(338198183)
  void clearCiphertextblob() => $_clearField(338198183);
}

class GenerateMacRequest extends $pb.GeneratedMessage {
  factory GenerateMacRequest({
    $core.bool? dryrun,
    $core.List<$core.int>? message,
    MacAlgorithmSpec? macalgorithm,
    $core.String? keyid,
    $core.Iterable<$core.String>? granttokens,
  }) {
    final result = create();
    if (dryrun != null) result.dryrun = dryrun;
    if (message != null) result.message = message;
    if (macalgorithm != null) result.macalgorithm = macalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    return result;
  }

  GenerateMacRequest._();

  factory GenerateMacRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateMacRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateMacRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..a<$core.List<$core.int>>(
        235854213, _omitFieldNames ? '' : 'message', $pb.PbFieldType.OY)
    ..aE<MacAlgorithmSpec>(253519878, _omitFieldNames ? '' : 'macalgorithm',
        enumValues: MacAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateMacRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateMacRequest copyWith(void Function(GenerateMacRequest) updates) =>
      super.copyWith((message) => updates(message as GenerateMacRequest))
          as GenerateMacRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateMacRequest create() => GenerateMacRequest._();
  @$core.override
  GenerateMacRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateMacRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateMacRequest>(create);
  static GenerateMacRequest? _defaultInstance;

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(0);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(0);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(235854213)
  $core.List<$core.int> get message => $_getN(1);
  @$pb.TagNumber(235854213)
  set message($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(253519878)
  MacAlgorithmSpec get macalgorithm => $_getN(2);
  @$pb.TagNumber(253519878)
  set macalgorithm(MacAlgorithmSpec value) => $_setField(253519878, value);
  @$pb.TagNumber(253519878)
  $core.bool hasMacalgorithm() => $_has(2);
  @$pb.TagNumber(253519878)
  void clearMacalgorithm() => $_clearField(253519878);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(4);
}

class GenerateMacResponse extends $pb.GeneratedMessage {
  factory GenerateMacResponse({
    MacAlgorithmSpec? macalgorithm,
    $core.String? keyid,
    $core.List<$core.int>? mac,
  }) {
    final result = create();
    if (macalgorithm != null) result.macalgorithm = macalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (mac != null) result.mac = mac;
    return result;
  }

  GenerateMacResponse._();

  factory GenerateMacResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateMacResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateMacResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<MacAlgorithmSpec>(253519878, _omitFieldNames ? '' : 'macalgorithm',
        enumValues: MacAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..a<$core.List<$core.int>>(
        296223945, _omitFieldNames ? '' : 'mac', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateMacResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateMacResponse copyWith(void Function(GenerateMacResponse) updates) =>
      super.copyWith((message) => updates(message as GenerateMacResponse))
          as GenerateMacResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateMacResponse create() => GenerateMacResponse._();
  @$core.override
  GenerateMacResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateMacResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateMacResponse>(create);
  static GenerateMacResponse? _defaultInstance;

  @$pb.TagNumber(253519878)
  MacAlgorithmSpec get macalgorithm => $_getN(0);
  @$pb.TagNumber(253519878)
  set macalgorithm(MacAlgorithmSpec value) => $_setField(253519878, value);
  @$pb.TagNumber(253519878)
  $core.bool hasMacalgorithm() => $_has(0);
  @$pb.TagNumber(253519878)
  void clearMacalgorithm() => $_clearField(253519878);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(296223945)
  $core.List<$core.int> get mac => $_getN(2);
  @$pb.TagNumber(296223945)
  set mac($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(296223945)
  $core.bool hasMac() => $_has(2);
  @$pb.TagNumber(296223945)
  void clearMac() => $_clearField(296223945);
}

class GenerateRandomRequest extends $pb.GeneratedMessage {
  factory GenerateRandomRequest({
    $core.String? customkeystoreid,
    $core.int? numberofbytes,
    RecipientInfo? recipient,
  }) {
    final result = create();
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    if (numberofbytes != null) result.numberofbytes = numberofbytes;
    if (recipient != null) result.recipient = recipient;
    return result;
  }

  GenerateRandomRequest._();

  factory GenerateRandomRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateRandomRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateRandomRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..aI(277086201, _omitFieldNames ? '' : 'numberofbytes')
    ..aOM<RecipientInfo>(445981721, _omitFieldNames ? '' : 'recipient',
        subBuilder: RecipientInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateRandomRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateRandomRequest copyWith(
          void Function(GenerateRandomRequest) updates) =>
      super.copyWith((message) => updates(message as GenerateRandomRequest))
          as GenerateRandomRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateRandomRequest create() => GenerateRandomRequest._();
  @$core.override
  GenerateRandomRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateRandomRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateRandomRequest>(create);
  static GenerateRandomRequest? _defaultInstance;

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(0);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(0);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);

  @$pb.TagNumber(277086201)
  $core.int get numberofbytes => $_getIZ(1);
  @$pb.TagNumber(277086201)
  set numberofbytes($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(277086201)
  $core.bool hasNumberofbytes() => $_has(1);
  @$pb.TagNumber(277086201)
  void clearNumberofbytes() => $_clearField(277086201);

  @$pb.TagNumber(445981721)
  RecipientInfo get recipient => $_getN(2);
  @$pb.TagNumber(445981721)
  set recipient(RecipientInfo value) => $_setField(445981721, value);
  @$pb.TagNumber(445981721)
  $core.bool hasRecipient() => $_has(2);
  @$pb.TagNumber(445981721)
  void clearRecipient() => $_clearField(445981721);
  @$pb.TagNumber(445981721)
  RecipientInfo ensureRecipient() => $_ensure(2);
}

class GenerateRandomResponse extends $pb.GeneratedMessage {
  factory GenerateRandomResponse({
    $core.List<$core.int>? ciphertextforrecipient,
    $core.List<$core.int>? plaintext,
  }) {
    final result = create();
    if (ciphertextforrecipient != null)
      result.ciphertextforrecipient = ciphertextforrecipient;
    if (plaintext != null) result.plaintext = plaintext;
    return result;
  }

  GenerateRandomResponse._();

  factory GenerateRandomResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateRandomResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateRandomResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(75299212,
        _omitFieldNames ? '' : 'ciphertextforrecipient', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(
        88342721, _omitFieldNames ? '' : 'plaintext', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateRandomResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateRandomResponse copyWith(
          void Function(GenerateRandomResponse) updates) =>
      super.copyWith((message) => updates(message as GenerateRandomResponse))
          as GenerateRandomResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateRandomResponse create() => GenerateRandomResponse._();
  @$core.override
  GenerateRandomResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateRandomResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateRandomResponse>(create);
  static GenerateRandomResponse? _defaultInstance;

  @$pb.TagNumber(75299212)
  $core.List<$core.int> get ciphertextforrecipient => $_getN(0);
  @$pb.TagNumber(75299212)
  set ciphertextforrecipient($core.List<$core.int> value) =>
      $_setBytes(0, value);
  @$pb.TagNumber(75299212)
  $core.bool hasCiphertextforrecipient() => $_has(0);
  @$pb.TagNumber(75299212)
  void clearCiphertextforrecipient() => $_clearField(75299212);

  @$pb.TagNumber(88342721)
  $core.List<$core.int> get plaintext => $_getN(1);
  @$pb.TagNumber(88342721)
  set plaintext($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(88342721)
  $core.bool hasPlaintext() => $_has(1);
  @$pb.TagNumber(88342721)
  void clearPlaintext() => $_clearField(88342721);
}

class GetKeyPolicyRequest extends $pb.GeneratedMessage {
  factory GetKeyPolicyRequest({
    $core.String? policyname,
    $core.String? keyid,
  }) {
    final result = create();
    if (policyname != null) result.policyname = policyname;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  GetKeyPolicyRequest._();

  factory GetKeyPolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetKeyPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetKeyPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(266468029, _omitFieldNames ? '' : 'policyname')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetKeyPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetKeyPolicyRequest copyWith(void Function(GetKeyPolicyRequest) updates) =>
      super.copyWith((message) => updates(message as GetKeyPolicyRequest))
          as GetKeyPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetKeyPolicyRequest create() => GetKeyPolicyRequest._();
  @$core.override
  GetKeyPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetKeyPolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetKeyPolicyRequest>(create);
  static GetKeyPolicyRequest? _defaultInstance;

  @$pb.TagNumber(266468029)
  $core.String get policyname => $_getSZ(0);
  @$pb.TagNumber(266468029)
  set policyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266468029)
  $core.bool hasPolicyname() => $_has(0);
  @$pb.TagNumber(266468029)
  void clearPolicyname() => $_clearField(266468029);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class GetKeyPolicyResponse extends $pb.GeneratedMessage {
  factory GetKeyPolicyResponse({
    $core.String? policyname,
    $core.String? policy,
  }) {
    final result = create();
    if (policyname != null) result.policyname = policyname;
    if (policy != null) result.policy = policy;
    return result;
  }

  GetKeyPolicyResponse._();

  factory GetKeyPolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetKeyPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetKeyPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(266468029, _omitFieldNames ? '' : 'policyname')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetKeyPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetKeyPolicyResponse copyWith(void Function(GetKeyPolicyResponse) updates) =>
      super.copyWith((message) => updates(message as GetKeyPolicyResponse))
          as GetKeyPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetKeyPolicyResponse create() => GetKeyPolicyResponse._();
  @$core.override
  GetKeyPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetKeyPolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetKeyPolicyResponse>(create);
  static GetKeyPolicyResponse? _defaultInstance;

  @$pb.TagNumber(266468029)
  $core.String get policyname => $_getSZ(0);
  @$pb.TagNumber(266468029)
  set policyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266468029)
  $core.bool hasPolicyname() => $_has(0);
  @$pb.TagNumber(266468029)
  void clearPolicyname() => $_clearField(266468029);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(1);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(1, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(1);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class GetKeyRotationStatusRequest extends $pb.GeneratedMessage {
  factory GetKeyRotationStatusRequest({
    $core.String? keyid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  GetKeyRotationStatusRequest._();

  factory GetKeyRotationStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetKeyRotationStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetKeyRotationStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetKeyRotationStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetKeyRotationStatusRequest copyWith(
          void Function(GetKeyRotationStatusRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetKeyRotationStatusRequest))
          as GetKeyRotationStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetKeyRotationStatusRequest create() =>
      GetKeyRotationStatusRequest._();
  @$core.override
  GetKeyRotationStatusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetKeyRotationStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetKeyRotationStatusRequest>(create);
  static GetKeyRotationStatusRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class GetKeyRotationStatusResponse extends $pb.GeneratedMessage {
  factory GetKeyRotationStatusResponse({
    $core.int? rotationperiodindays,
    $core.String? nextrotationdate,
    $core.String? keyid,
    $core.String? ondemandrotationstartdate,
    $core.bool? keyrotationenabled,
  }) {
    final result = create();
    if (rotationperiodindays != null)
      result.rotationperiodindays = rotationperiodindays;
    if (nextrotationdate != null) result.nextrotationdate = nextrotationdate;
    if (keyid != null) result.keyid = keyid;
    if (ondemandrotationstartdate != null)
      result.ondemandrotationstartdate = ondemandrotationstartdate;
    if (keyrotationenabled != null)
      result.keyrotationenabled = keyrotationenabled;
    return result;
  }

  GetKeyRotationStatusResponse._();

  factory GetKeyRotationStatusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetKeyRotationStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetKeyRotationStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aI(118357231, _omitFieldNames ? '' : 'rotationperiodindays')
    ..aOS(192035355, _omitFieldNames ? '' : 'nextrotationdate')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(360279652, _omitFieldNames ? '' : 'ondemandrotationstartdate')
    ..aOB(525956616, _omitFieldNames ? '' : 'keyrotationenabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetKeyRotationStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetKeyRotationStatusResponse copyWith(
          void Function(GetKeyRotationStatusResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetKeyRotationStatusResponse))
          as GetKeyRotationStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetKeyRotationStatusResponse create() =>
      GetKeyRotationStatusResponse._();
  @$core.override
  GetKeyRotationStatusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetKeyRotationStatusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetKeyRotationStatusResponse>(create);
  static GetKeyRotationStatusResponse? _defaultInstance;

  @$pb.TagNumber(118357231)
  $core.int get rotationperiodindays => $_getIZ(0);
  @$pb.TagNumber(118357231)
  set rotationperiodindays($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(118357231)
  $core.bool hasRotationperiodindays() => $_has(0);
  @$pb.TagNumber(118357231)
  void clearRotationperiodindays() => $_clearField(118357231);

  @$pb.TagNumber(192035355)
  $core.String get nextrotationdate => $_getSZ(1);
  @$pb.TagNumber(192035355)
  set nextrotationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(192035355)
  $core.bool hasNextrotationdate() => $_has(1);
  @$pb.TagNumber(192035355)
  void clearNextrotationdate() => $_clearField(192035355);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(360279652)
  $core.String get ondemandrotationstartdate => $_getSZ(3);
  @$pb.TagNumber(360279652)
  set ondemandrotationstartdate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(360279652)
  $core.bool hasOndemandrotationstartdate() => $_has(3);
  @$pb.TagNumber(360279652)
  void clearOndemandrotationstartdate() => $_clearField(360279652);

  @$pb.TagNumber(525956616)
  $core.bool get keyrotationenabled => $_getBF(4);
  @$pb.TagNumber(525956616)
  set keyrotationenabled($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(525956616)
  $core.bool hasKeyrotationenabled() => $_has(4);
  @$pb.TagNumber(525956616)
  void clearKeyrotationenabled() => $_clearField(525956616);
}

class GetParametersForImportRequest extends $pb.GeneratedMessage {
  factory GetParametersForImportRequest({
    AlgorithmSpec? wrappingalgorithm,
    $core.String? keyid,
    WrappingKeySpec? wrappingkeyspec,
  }) {
    final result = create();
    if (wrappingalgorithm != null) result.wrappingalgorithm = wrappingalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (wrappingkeyspec != null) result.wrappingkeyspec = wrappingkeyspec;
    return result;
  }

  GetParametersForImportRequest._();

  factory GetParametersForImportRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetParametersForImportRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetParametersForImportRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<AlgorithmSpec>(164083247, _omitFieldNames ? '' : 'wrappingalgorithm',
        enumValues: AlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aE<WrappingKeySpec>(521228924, _omitFieldNames ? '' : 'wrappingkeyspec',
        enumValues: WrappingKeySpec.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetParametersForImportRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetParametersForImportRequest copyWith(
          void Function(GetParametersForImportRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetParametersForImportRequest))
          as GetParametersForImportRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetParametersForImportRequest create() =>
      GetParametersForImportRequest._();
  @$core.override
  GetParametersForImportRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetParametersForImportRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetParametersForImportRequest>(create);
  static GetParametersForImportRequest? _defaultInstance;

  @$pb.TagNumber(164083247)
  AlgorithmSpec get wrappingalgorithm => $_getN(0);
  @$pb.TagNumber(164083247)
  set wrappingalgorithm(AlgorithmSpec value) => $_setField(164083247, value);
  @$pb.TagNumber(164083247)
  $core.bool hasWrappingalgorithm() => $_has(0);
  @$pb.TagNumber(164083247)
  void clearWrappingalgorithm() => $_clearField(164083247);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(521228924)
  WrappingKeySpec get wrappingkeyspec => $_getN(2);
  @$pb.TagNumber(521228924)
  set wrappingkeyspec(WrappingKeySpec value) => $_setField(521228924, value);
  @$pb.TagNumber(521228924)
  $core.bool hasWrappingkeyspec() => $_has(2);
  @$pb.TagNumber(521228924)
  void clearWrappingkeyspec() => $_clearField(521228924);
}

class GetParametersForImportResponse extends $pb.GeneratedMessage {
  factory GetParametersForImportResponse({
    $core.List<$core.int>? publickey,
    $core.String? keyid,
    $core.String? parametersvalidto,
    $core.List<$core.int>? importtoken,
  }) {
    final result = create();
    if (publickey != null) result.publickey = publickey;
    if (keyid != null) result.keyid = keyid;
    if (parametersvalidto != null) result.parametersvalidto = parametersvalidto;
    if (importtoken != null) result.importtoken = importtoken;
    return result;
  }

  GetParametersForImportResponse._();

  factory GetParametersForImportResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetParametersForImportResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetParametersForImportResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(
        167335776, _omitFieldNames ? '' : 'publickey', $pb.PbFieldType.OY)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(394913717, _omitFieldNames ? '' : 'parametersvalidto')
    ..a<$core.List<$core.int>>(
        461726162, _omitFieldNames ? '' : 'importtoken', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetParametersForImportResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetParametersForImportResponse copyWith(
          void Function(GetParametersForImportResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetParametersForImportResponse))
          as GetParametersForImportResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetParametersForImportResponse create() =>
      GetParametersForImportResponse._();
  @$core.override
  GetParametersForImportResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetParametersForImportResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetParametersForImportResponse>(create);
  static GetParametersForImportResponse? _defaultInstance;

  @$pb.TagNumber(167335776)
  $core.List<$core.int> get publickey => $_getN(0);
  @$pb.TagNumber(167335776)
  set publickey($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(167335776)
  $core.bool hasPublickey() => $_has(0);
  @$pb.TagNumber(167335776)
  void clearPublickey() => $_clearField(167335776);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(394913717)
  $core.String get parametersvalidto => $_getSZ(2);
  @$pb.TagNumber(394913717)
  set parametersvalidto($core.String value) => $_setString(2, value);
  @$pb.TagNumber(394913717)
  $core.bool hasParametersvalidto() => $_has(2);
  @$pb.TagNumber(394913717)
  void clearParametersvalidto() => $_clearField(394913717);

  @$pb.TagNumber(461726162)
  $core.List<$core.int> get importtoken => $_getN(3);
  @$pb.TagNumber(461726162)
  set importtoken($core.List<$core.int> value) => $_setBytes(3, value);
  @$pb.TagNumber(461726162)
  $core.bool hasImporttoken() => $_has(3);
  @$pb.TagNumber(461726162)
  void clearImporttoken() => $_clearField(461726162);
}

class GetPublicKeyRequest extends $pb.GeneratedMessage {
  factory GetPublicKeyRequest({
    $core.String? keyid,
    $core.Iterable<$core.String>? granttokens,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    return result;
  }

  GetPublicKeyRequest._();

  factory GetPublicKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPublicKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPublicKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPublicKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPublicKeyRequest copyWith(void Function(GetPublicKeyRequest) updates) =>
      super.copyWith((message) => updates(message as GetPublicKeyRequest))
          as GetPublicKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPublicKeyRequest create() => GetPublicKeyRequest._();
  @$core.override
  GetPublicKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPublicKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPublicKeyRequest>(create);
  static GetPublicKeyRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(1);
}

class GetPublicKeyResponse extends $pb.GeneratedMessage {
  factory GetPublicKeyResponse({
    KeySpec? keyspec,
    $core.List<$core.int>? publickey,
    $core.Iterable<EncryptionAlgorithmSpec>? encryptionalgorithms,
    $core.String? keyid,
    $core.Iterable<KeyAgreementAlgorithmSpec>? keyagreementalgorithms,
    KeyUsageType? keyusage,
    CustomerMasterKeySpec? customermasterkeyspec,
    $core.Iterable<SigningAlgorithmSpec>? signingalgorithms,
  }) {
    final result = create();
    if (keyspec != null) result.keyspec = keyspec;
    if (publickey != null) result.publickey = publickey;
    if (encryptionalgorithms != null)
      result.encryptionalgorithms.addAll(encryptionalgorithms);
    if (keyid != null) result.keyid = keyid;
    if (keyagreementalgorithms != null)
      result.keyagreementalgorithms.addAll(keyagreementalgorithms);
    if (keyusage != null) result.keyusage = keyusage;
    if (customermasterkeyspec != null)
      result.customermasterkeyspec = customermasterkeyspec;
    if (signingalgorithms != null)
      result.signingalgorithms.addAll(signingalgorithms);
    return result;
  }

  GetPublicKeyResponse._();

  factory GetPublicKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPublicKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPublicKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<KeySpec>(138220928, _omitFieldNames ? '' : 'keyspec',
        enumValues: KeySpec.values)
    ..a<$core.List<$core.int>>(
        167335776, _omitFieldNames ? '' : 'publickey', $pb.PbFieldType.OY)
    ..pc<EncryptionAlgorithmSpec>(194511375,
        _omitFieldNames ? '' : 'encryptionalgorithms', $pb.PbFieldType.KE,
        valueOf: EncryptionAlgorithmSpec.valueOf,
        enumValues: EncryptionAlgorithmSpec.values,
        defaultEnumValue:
            EncryptionAlgorithmSpec.ENCRYPTION_ALGORITHM_SPEC_SM2PKE)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pc<KeyAgreementAlgorithmSpec>(328746163,
        _omitFieldNames ? '' : 'keyagreementalgorithms', $pb.PbFieldType.KE,
        valueOf: KeyAgreementAlgorithmSpec.valueOf,
        enumValues: KeyAgreementAlgorithmSpec.values,
        defaultEnumValue:
            KeyAgreementAlgorithmSpec.KEY_AGREEMENT_ALGORITHM_SPEC_ECDH)
    ..aE<KeyUsageType>(357216772, _omitFieldNames ? '' : 'keyusage',
        enumValues: KeyUsageType.values)
    ..aE<CustomerMasterKeySpec>(
        472470930, _omitFieldNames ? '' : 'customermasterkeyspec',
        enumValues: CustomerMasterKeySpec.values)
    ..pc<SigningAlgorithmSpec>(508241975,
        _omitFieldNames ? '' : 'signingalgorithms', $pb.PbFieldType.KE,
        valueOf: SigningAlgorithmSpec.valueOf,
        enumValues: SigningAlgorithmSpec.values,
        defaultEnumValue:
            SigningAlgorithmSpec.SIGNING_ALGORITHM_SPEC_ECDSA_SHA_384)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPublicKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPublicKeyResponse copyWith(void Function(GetPublicKeyResponse) updates) =>
      super.copyWith((message) => updates(message as GetPublicKeyResponse))
          as GetPublicKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPublicKeyResponse create() => GetPublicKeyResponse._();
  @$core.override
  GetPublicKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPublicKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPublicKeyResponse>(create);
  static GetPublicKeyResponse? _defaultInstance;

  @$pb.TagNumber(138220928)
  KeySpec get keyspec => $_getN(0);
  @$pb.TagNumber(138220928)
  set keyspec(KeySpec value) => $_setField(138220928, value);
  @$pb.TagNumber(138220928)
  $core.bool hasKeyspec() => $_has(0);
  @$pb.TagNumber(138220928)
  void clearKeyspec() => $_clearField(138220928);

  @$pb.TagNumber(167335776)
  $core.List<$core.int> get publickey => $_getN(1);
  @$pb.TagNumber(167335776)
  set publickey($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(167335776)
  $core.bool hasPublickey() => $_has(1);
  @$pb.TagNumber(167335776)
  void clearPublickey() => $_clearField(167335776);

  @$pb.TagNumber(194511375)
  $pb.PbList<EncryptionAlgorithmSpec> get encryptionalgorithms => $_getList(2);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(328746163)
  $pb.PbList<KeyAgreementAlgorithmSpec> get keyagreementalgorithms =>
      $_getList(4);

  @$pb.TagNumber(357216772)
  KeyUsageType get keyusage => $_getN(5);
  @$pb.TagNumber(357216772)
  set keyusage(KeyUsageType value) => $_setField(357216772, value);
  @$pb.TagNumber(357216772)
  $core.bool hasKeyusage() => $_has(5);
  @$pb.TagNumber(357216772)
  void clearKeyusage() => $_clearField(357216772);

  @$pb.TagNumber(472470930)
  CustomerMasterKeySpec get customermasterkeyspec => $_getN(6);
  @$pb.TagNumber(472470930)
  set customermasterkeyspec(CustomerMasterKeySpec value) =>
      $_setField(472470930, value);
  @$pb.TagNumber(472470930)
  $core.bool hasCustomermasterkeyspec() => $_has(6);
  @$pb.TagNumber(472470930)
  void clearCustomermasterkeyspec() => $_clearField(472470930);

  @$pb.TagNumber(508241975)
  $pb.PbList<SigningAlgorithmSpec> get signingalgorithms => $_getList(7);
}

class GrantConstraints extends $pb.GeneratedMessage {
  factory GrantConstraints({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        encryptioncontextequals,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        encryptioncontextsubset,
  }) {
    final result = create();
    if (encryptioncontextequals != null)
      result.encryptioncontextequals.addEntries(encryptioncontextequals);
    if (encryptioncontextsubset != null)
      result.encryptioncontextsubset.addEntries(encryptioncontextsubset);
    return result;
  }

  GrantConstraints._();

  factory GrantConstraints.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GrantConstraints.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GrantConstraints',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        68403171, _omitFieldNames ? '' : 'encryptioncontextequals',
        entryClassName: 'GrantConstraints.EncryptioncontextequalsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..m<$core.String, $core.String>(
        72310514, _omitFieldNames ? '' : 'encryptioncontextsubset',
        entryClassName: 'GrantConstraints.EncryptioncontextsubsetEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GrantConstraints clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GrantConstraints copyWith(void Function(GrantConstraints) updates) =>
      super.copyWith((message) => updates(message as GrantConstraints))
          as GrantConstraints;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GrantConstraints create() => GrantConstraints._();
  @$core.override
  GrantConstraints createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GrantConstraints getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GrantConstraints>(create);
  static GrantConstraints? _defaultInstance;

  @$pb.TagNumber(68403171)
  $pb.PbMap<$core.String, $core.String> get encryptioncontextequals =>
      $_getMap(0);

  @$pb.TagNumber(72310514)
  $pb.PbMap<$core.String, $core.String> get encryptioncontextsubset =>
      $_getMap(1);
}

class GrantListEntry extends $pb.GeneratedMessage {
  factory GrantListEntry({
    $core.String? issuingaccount,
    $core.String? retiringprincipal,
    $core.String? grantid,
    $core.Iterable<GrantOperation>? operations,
    $core.String? granteeprincipal,
    $core.String? name,
    $core.String? keyid,
    $core.String? creationdate,
    GrantConstraints? constraints,
  }) {
    final result = create();
    if (issuingaccount != null) result.issuingaccount = issuingaccount;
    if (retiringprincipal != null) result.retiringprincipal = retiringprincipal;
    if (grantid != null) result.grantid = grantid;
    if (operations != null) result.operations.addAll(operations);
    if (granteeprincipal != null) result.granteeprincipal = granteeprincipal;
    if (name != null) result.name = name;
    if (keyid != null) result.keyid = keyid;
    if (creationdate != null) result.creationdate = creationdate;
    if (constraints != null) result.constraints = constraints;
    return result;
  }

  GrantListEntry._();

  factory GrantListEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GrantListEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GrantListEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(47662575, _omitFieldNames ? '' : 'issuingaccount')
    ..aOS(49541086, _omitFieldNames ? '' : 'retiringprincipal')
    ..aOS(66852281, _omitFieldNames ? '' : 'grantid')
    ..pc<GrantOperation>(
        126776656, _omitFieldNames ? '' : 'operations', $pb.PbFieldType.KE,
        valueOf: GrantOperation.valueOf,
        enumValues: GrantOperation.values,
        defaultEnumValue: GrantOperation.GRANT_OPERATION_GENERATEDATAKEYPAIR)
    ..aOS(234727364, _omitFieldNames ? '' : 'granteeprincipal')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aOM<GrantConstraints>(302297388, _omitFieldNames ? '' : 'constraints',
        subBuilder: GrantConstraints.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GrantListEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GrantListEntry copyWith(void Function(GrantListEntry) updates) =>
      super.copyWith((message) => updates(message as GrantListEntry))
          as GrantListEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GrantListEntry create() => GrantListEntry._();
  @$core.override
  GrantListEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GrantListEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GrantListEntry>(create);
  static GrantListEntry? _defaultInstance;

  @$pb.TagNumber(47662575)
  $core.String get issuingaccount => $_getSZ(0);
  @$pb.TagNumber(47662575)
  set issuingaccount($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47662575)
  $core.bool hasIssuingaccount() => $_has(0);
  @$pb.TagNumber(47662575)
  void clearIssuingaccount() => $_clearField(47662575);

  @$pb.TagNumber(49541086)
  $core.String get retiringprincipal => $_getSZ(1);
  @$pb.TagNumber(49541086)
  set retiringprincipal($core.String value) => $_setString(1, value);
  @$pb.TagNumber(49541086)
  $core.bool hasRetiringprincipal() => $_has(1);
  @$pb.TagNumber(49541086)
  void clearRetiringprincipal() => $_clearField(49541086);

  @$pb.TagNumber(66852281)
  $core.String get grantid => $_getSZ(2);
  @$pb.TagNumber(66852281)
  set grantid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(66852281)
  $core.bool hasGrantid() => $_has(2);
  @$pb.TagNumber(66852281)
  void clearGrantid() => $_clearField(66852281);

  @$pb.TagNumber(126776656)
  $pb.PbList<GrantOperation> get operations => $_getList(3);

  @$pb.TagNumber(234727364)
  $core.String get granteeprincipal => $_getSZ(4);
  @$pb.TagNumber(234727364)
  set granteeprincipal($core.String value) => $_setString(4, value);
  @$pb.TagNumber(234727364)
  $core.bool hasGranteeprincipal() => $_has(4);
  @$pb.TagNumber(234727364)
  void clearGranteeprincipal() => $_clearField(234727364);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(6);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(6);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(7);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(7, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(7);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(302297388)
  GrantConstraints get constraints => $_getN(8);
  @$pb.TagNumber(302297388)
  set constraints(GrantConstraints value) => $_setField(302297388, value);
  @$pb.TagNumber(302297388)
  $core.bool hasConstraints() => $_has(8);
  @$pb.TagNumber(302297388)
  void clearConstraints() => $_clearField(302297388);
  @$pb.TagNumber(302297388)
  GrantConstraints ensureConstraints() => $_ensure(8);
}

class ImportKeyMaterialRequest extends $pb.GeneratedMessage {
  factory ImportKeyMaterialRequest({
    $core.List<$core.int>? encryptedkeymaterial,
    ExpirationModelType? expirationmodel,
    $core.String? keymaterialid,
    $core.String? keyid,
    $core.String? keymaterialdescription,
    ImportType? importtype,
    $core.List<$core.int>? importtoken,
    $core.String? validto,
  }) {
    final result = create();
    if (encryptedkeymaterial != null)
      result.encryptedkeymaterial = encryptedkeymaterial;
    if (expirationmodel != null) result.expirationmodel = expirationmodel;
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (keyid != null) result.keyid = keyid;
    if (keymaterialdescription != null)
      result.keymaterialdescription = keymaterialdescription;
    if (importtype != null) result.importtype = importtype;
    if (importtoken != null) result.importtoken = importtoken;
    if (validto != null) result.validto = validto;
    return result;
  }

  ImportKeyMaterialRequest._();

  factory ImportKeyMaterialRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportKeyMaterialRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportKeyMaterialRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(3185968,
        _omitFieldNames ? '' : 'encryptedkeymaterial', $pb.PbFieldType.OY)
    ..aE<ExpirationModelType>(
        113933558, _omitFieldNames ? '' : 'expirationmodel',
        enumValues: ExpirationModelType.values)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(277153546, _omitFieldNames ? '' : 'keymaterialdescription')
    ..aE<ImportType>(331980349, _omitFieldNames ? '' : 'importtype',
        enumValues: ImportType.values)
    ..a<$core.List<$core.int>>(
        461726162, _omitFieldNames ? '' : 'importtoken', $pb.PbFieldType.OY)
    ..aOS(522718673, _omitFieldNames ? '' : 'validto')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportKeyMaterialRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportKeyMaterialRequest copyWith(
          void Function(ImportKeyMaterialRequest) updates) =>
      super.copyWith((message) => updates(message as ImportKeyMaterialRequest))
          as ImportKeyMaterialRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportKeyMaterialRequest create() => ImportKeyMaterialRequest._();
  @$core.override
  ImportKeyMaterialRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportKeyMaterialRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportKeyMaterialRequest>(create);
  static ImportKeyMaterialRequest? _defaultInstance;

  @$pb.TagNumber(3185968)
  $core.List<$core.int> get encryptedkeymaterial => $_getN(0);
  @$pb.TagNumber(3185968)
  set encryptedkeymaterial($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(3185968)
  $core.bool hasEncryptedkeymaterial() => $_has(0);
  @$pb.TagNumber(3185968)
  void clearEncryptedkeymaterial() => $_clearField(3185968);

  @$pb.TagNumber(113933558)
  ExpirationModelType get expirationmodel => $_getN(1);
  @$pb.TagNumber(113933558)
  set expirationmodel(ExpirationModelType value) =>
      $_setField(113933558, value);
  @$pb.TagNumber(113933558)
  $core.bool hasExpirationmodel() => $_has(1);
  @$pb.TagNumber(113933558)
  void clearExpirationmodel() => $_clearField(113933558);

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(2);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(2);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(277153546)
  $core.String get keymaterialdescription => $_getSZ(4);
  @$pb.TagNumber(277153546)
  set keymaterialdescription($core.String value) => $_setString(4, value);
  @$pb.TagNumber(277153546)
  $core.bool hasKeymaterialdescription() => $_has(4);
  @$pb.TagNumber(277153546)
  void clearKeymaterialdescription() => $_clearField(277153546);

  @$pb.TagNumber(331980349)
  ImportType get importtype => $_getN(5);
  @$pb.TagNumber(331980349)
  set importtype(ImportType value) => $_setField(331980349, value);
  @$pb.TagNumber(331980349)
  $core.bool hasImporttype() => $_has(5);
  @$pb.TagNumber(331980349)
  void clearImporttype() => $_clearField(331980349);

  @$pb.TagNumber(461726162)
  $core.List<$core.int> get importtoken => $_getN(6);
  @$pb.TagNumber(461726162)
  set importtoken($core.List<$core.int> value) => $_setBytes(6, value);
  @$pb.TagNumber(461726162)
  $core.bool hasImporttoken() => $_has(6);
  @$pb.TagNumber(461726162)
  void clearImporttoken() => $_clearField(461726162);

  @$pb.TagNumber(522718673)
  $core.String get validto => $_getSZ(7);
  @$pb.TagNumber(522718673)
  set validto($core.String value) => $_setString(7, value);
  @$pb.TagNumber(522718673)
  $core.bool hasValidto() => $_has(7);
  @$pb.TagNumber(522718673)
  void clearValidto() => $_clearField(522718673);
}

class ImportKeyMaterialResponse extends $pb.GeneratedMessage {
  factory ImportKeyMaterialResponse({
    $core.String? keymaterialid,
    $core.String? keyid,
  }) {
    final result = create();
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  ImportKeyMaterialResponse._();

  factory ImportKeyMaterialResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportKeyMaterialResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportKeyMaterialResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportKeyMaterialResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportKeyMaterialResponse copyWith(
          void Function(ImportKeyMaterialResponse) updates) =>
      super.copyWith((message) => updates(message as ImportKeyMaterialResponse))
          as ImportKeyMaterialResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportKeyMaterialResponse create() => ImportKeyMaterialResponse._();
  @$core.override
  ImportKeyMaterialResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportKeyMaterialResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportKeyMaterialResponse>(create);
  static ImportKeyMaterialResponse? _defaultInstance;

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(0);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(0);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class IncorrectKeyException extends $pb.GeneratedMessage {
  factory IncorrectKeyException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IncorrectKeyException._();

  factory IncorrectKeyException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IncorrectKeyException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IncorrectKeyException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncorrectKeyException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncorrectKeyException copyWith(
          void Function(IncorrectKeyException) updates) =>
      super.copyWith((message) => updates(message as IncorrectKeyException))
          as IncorrectKeyException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IncorrectKeyException create() => IncorrectKeyException._();
  @$core.override
  IncorrectKeyException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IncorrectKeyException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IncorrectKeyException>(create);
  static IncorrectKeyException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class IncorrectKeyMaterialException extends $pb.GeneratedMessage {
  factory IncorrectKeyMaterialException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IncorrectKeyMaterialException._();

  factory IncorrectKeyMaterialException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IncorrectKeyMaterialException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IncorrectKeyMaterialException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncorrectKeyMaterialException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncorrectKeyMaterialException copyWith(
          void Function(IncorrectKeyMaterialException) updates) =>
      super.copyWith(
              (message) => updates(message as IncorrectKeyMaterialException))
          as IncorrectKeyMaterialException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IncorrectKeyMaterialException create() =>
      IncorrectKeyMaterialException._();
  @$core.override
  IncorrectKeyMaterialException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IncorrectKeyMaterialException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IncorrectKeyMaterialException>(create);
  static IncorrectKeyMaterialException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class IncorrectTrustAnchorException extends $pb.GeneratedMessage {
  factory IncorrectTrustAnchorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IncorrectTrustAnchorException._();

  factory IncorrectTrustAnchorException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IncorrectTrustAnchorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IncorrectTrustAnchorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncorrectTrustAnchorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncorrectTrustAnchorException copyWith(
          void Function(IncorrectTrustAnchorException) updates) =>
      super.copyWith(
              (message) => updates(message as IncorrectTrustAnchorException))
          as IncorrectTrustAnchorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IncorrectTrustAnchorException create() =>
      IncorrectTrustAnchorException._();
  @$core.override
  IncorrectTrustAnchorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IncorrectTrustAnchorException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IncorrectTrustAnchorException>(create);
  static IncorrectTrustAnchorException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidAliasNameException extends $pb.GeneratedMessage {
  factory InvalidAliasNameException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidAliasNameException._();

  factory InvalidAliasNameException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidAliasNameException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidAliasNameException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAliasNameException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidAliasNameException copyWith(
          void Function(InvalidAliasNameException) updates) =>
      super.copyWith((message) => updates(message as InvalidAliasNameException))
          as InvalidAliasNameException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidAliasNameException create() => InvalidAliasNameException._();
  @$core.override
  InvalidAliasNameException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidAliasNameException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidAliasNameException>(create);
  static InvalidAliasNameException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidArnException extends $pb.GeneratedMessage {
  factory InvalidArnException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidArnException._();

  factory InvalidArnException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidArnException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidArnException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArnException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArnException copyWith(void Function(InvalidArnException) updates) =>
      super.copyWith((message) => updates(message as InvalidArnException))
          as InvalidArnException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidArnException create() => InvalidArnException._();
  @$core.override
  InvalidArnException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidArnException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidArnException>(create);
  static InvalidArnException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidCiphertextException extends $pb.GeneratedMessage {
  factory InvalidCiphertextException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidCiphertextException._();

  factory InvalidCiphertextException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidCiphertextException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidCiphertextException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidCiphertextException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidCiphertextException copyWith(
          void Function(InvalidCiphertextException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidCiphertextException))
          as InvalidCiphertextException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidCiphertextException create() => InvalidCiphertextException._();
  @$core.override
  InvalidCiphertextException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidCiphertextException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidCiphertextException>(create);
  static InvalidCiphertextException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidGrantIdException extends $pb.GeneratedMessage {
  factory InvalidGrantIdException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidGrantIdException._();

  factory InvalidGrantIdException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidGrantIdException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidGrantIdException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidGrantIdException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidGrantIdException copyWith(
          void Function(InvalidGrantIdException) updates) =>
      super.copyWith((message) => updates(message as InvalidGrantIdException))
          as InvalidGrantIdException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidGrantIdException create() => InvalidGrantIdException._();
  @$core.override
  InvalidGrantIdException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidGrantIdException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidGrantIdException>(create);
  static InvalidGrantIdException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidGrantTokenException extends $pb.GeneratedMessage {
  factory InvalidGrantTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidGrantTokenException._();

  factory InvalidGrantTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidGrantTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidGrantTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidGrantTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidGrantTokenException copyWith(
          void Function(InvalidGrantTokenException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidGrantTokenException))
          as InvalidGrantTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidGrantTokenException create() => InvalidGrantTokenException._();
  @$core.override
  InvalidGrantTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidGrantTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidGrantTokenException>(create);
  static InvalidGrantTokenException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidImportTokenException extends $pb.GeneratedMessage {
  factory InvalidImportTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidImportTokenException._();

  factory InvalidImportTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidImportTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidImportTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidImportTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidImportTokenException copyWith(
          void Function(InvalidImportTokenException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidImportTokenException))
          as InvalidImportTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidImportTokenException create() =>
      InvalidImportTokenException._();
  @$core.override
  InvalidImportTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidImportTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidImportTokenException>(create);
  static InvalidImportTokenException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidKeyUsageException extends $pb.GeneratedMessage {
  factory InvalidKeyUsageException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidKeyUsageException._();

  factory InvalidKeyUsageException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidKeyUsageException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidKeyUsageException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKeyUsageException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKeyUsageException copyWith(
          void Function(InvalidKeyUsageException) updates) =>
      super.copyWith((message) => updates(message as InvalidKeyUsageException))
          as InvalidKeyUsageException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidKeyUsageException create() => InvalidKeyUsageException._();
  @$core.override
  InvalidKeyUsageException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidKeyUsageException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidKeyUsageException>(create);
  static InvalidKeyUsageException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidMarkerException extends $pb.GeneratedMessage {
  factory InvalidMarkerException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidMarkerException._();

  factory InvalidMarkerException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidMarkerException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidMarkerException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidMarkerException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidMarkerException copyWith(
          void Function(InvalidMarkerException) updates) =>
      super.copyWith((message) => updates(message as InvalidMarkerException))
          as InvalidMarkerException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidMarkerException create() => InvalidMarkerException._();
  @$core.override
  InvalidMarkerException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidMarkerException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidMarkerException>(create);
  static InvalidMarkerException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSInternalException extends $pb.GeneratedMessage {
  factory KMSInternalException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSInternalException._();

  factory KMSInternalException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSInternalException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSInternalException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSInternalException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSInternalException copyWith(void Function(KMSInternalException) updates) =>
      super.copyWith((message) => updates(message as KMSInternalException))
          as KMSInternalException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSInternalException create() => KMSInternalException._();
  @$core.override
  KMSInternalException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSInternalException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSInternalException>(create);
  static KMSInternalException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSInvalidMacException extends $pb.GeneratedMessage {
  factory KMSInvalidMacException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSInvalidMacException._();

  factory KMSInvalidMacException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSInvalidMacException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSInvalidMacException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSInvalidMacException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSInvalidMacException copyWith(
          void Function(KMSInvalidMacException) updates) =>
      super.copyWith((message) => updates(message as KMSInvalidMacException))
          as KMSInvalidMacException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSInvalidMacException create() => KMSInvalidMacException._();
  @$core.override
  KMSInvalidMacException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSInvalidMacException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSInvalidMacException>(create);
  static KMSInvalidMacException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KMSInvalidSignatureException extends $pb.GeneratedMessage {
  factory KMSInvalidSignatureException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KMSInvalidSignatureException._();

  factory KMSInvalidSignatureException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KMSInvalidSignatureException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KMSInvalidSignatureException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSInvalidSignatureException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KMSInvalidSignatureException copyWith(
          void Function(KMSInvalidSignatureException) updates) =>
      super.copyWith(
              (message) => updates(message as KMSInvalidSignatureException))
          as KMSInvalidSignatureException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KMSInvalidSignatureException create() =>
      KMSInvalidSignatureException._();
  @$core.override
  KMSInvalidSignatureException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KMSInvalidSignatureException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KMSInvalidSignatureException>(create);
  static KMSInvalidSignatureException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
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

class KeyListEntry extends $pb.GeneratedMessage {
  factory KeyListEntry({
    $core.String? keyid,
    $core.String? keyarn,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    if (keyarn != null) result.keyarn = keyarn;
    return result;
  }

  KeyListEntry._();

  factory KeyListEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeyListEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeyListEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(418055012, _omitFieldNames ? '' : 'keyarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeyListEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeyListEntry copyWith(void Function(KeyListEntry) updates) =>
      super.copyWith((message) => updates(message as KeyListEntry))
          as KeyListEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeyListEntry create() => KeyListEntry._();
  @$core.override
  KeyListEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeyListEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeyListEntry>(create);
  static KeyListEntry? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(418055012)
  $core.String get keyarn => $_getSZ(1);
  @$pb.TagNumber(418055012)
  set keyarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(418055012)
  $core.bool hasKeyarn() => $_has(1);
  @$pb.TagNumber(418055012)
  void clearKeyarn() => $_clearField(418055012);
}

class KeyMetadata extends $pb.GeneratedMessage {
  factory KeyMetadata({
    KeyState? keystate,
    $core.String? cloudhsmclusterid,
    $core.String? customkeystoreid,
    ExpirationModelType? expirationmodel,
    $core.String? description,
    KeySpec? keyspec,
    $core.String? currentkeymaterialid,
    $core.Iterable<EncryptionAlgorithmSpec>? encryptionalgorithms,
    $core.String? keyid,
    $core.String? creationdate,
    $core.Iterable<KeyAgreementAlgorithmSpec>? keyagreementalgorithms,
    $core.String? deletiondate,
    KeyUsageType? keyusage,
    XksKeyConfigurationType? xkskeyconfiguration,
    $core.String? awsaccountid,
    $core.String? arn,
    $core.bool? multiregion,
    KeyManagerType? keymanager,
    MultiRegionConfiguration? multiregionconfiguration,
    CustomerMasterKeySpec? customermasterkeyspec,
    $core.bool? enabled,
    $core.int? pendingdeletionwindowindays,
    $core.Iterable<SigningAlgorithmSpec>? signingalgorithms,
    $core.String? validto,
    OriginType? origin,
    $core.Iterable<MacAlgorithmSpec>? macalgorithms,
  }) {
    final result = create();
    if (keystate != null) result.keystate = keystate;
    if (cloudhsmclusterid != null) result.cloudhsmclusterid = cloudhsmclusterid;
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    if (expirationmodel != null) result.expirationmodel = expirationmodel;
    if (description != null) result.description = description;
    if (keyspec != null) result.keyspec = keyspec;
    if (currentkeymaterialid != null)
      result.currentkeymaterialid = currentkeymaterialid;
    if (encryptionalgorithms != null)
      result.encryptionalgorithms.addAll(encryptionalgorithms);
    if (keyid != null) result.keyid = keyid;
    if (creationdate != null) result.creationdate = creationdate;
    if (keyagreementalgorithms != null)
      result.keyagreementalgorithms.addAll(keyagreementalgorithms);
    if (deletiondate != null) result.deletiondate = deletiondate;
    if (keyusage != null) result.keyusage = keyusage;
    if (xkskeyconfiguration != null)
      result.xkskeyconfiguration = xkskeyconfiguration;
    if (awsaccountid != null) result.awsaccountid = awsaccountid;
    if (arn != null) result.arn = arn;
    if (multiregion != null) result.multiregion = multiregion;
    if (keymanager != null) result.keymanager = keymanager;
    if (multiregionconfiguration != null)
      result.multiregionconfiguration = multiregionconfiguration;
    if (customermasterkeyspec != null)
      result.customermasterkeyspec = customermasterkeyspec;
    if (enabled != null) result.enabled = enabled;
    if (pendingdeletionwindowindays != null)
      result.pendingdeletionwindowindays = pendingdeletionwindowindays;
    if (signingalgorithms != null)
      result.signingalgorithms.addAll(signingalgorithms);
    if (validto != null) result.validto = validto;
    if (origin != null) result.origin = origin;
    if (macalgorithms != null) result.macalgorithms.addAll(macalgorithms);
    return result;
  }

  KeyMetadata._();

  factory KeyMetadata.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeyMetadata.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeyMetadata',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<KeyState>(27894226, _omitFieldNames ? '' : 'keystate',
        enumValues: KeyState.values)
    ..aOS(56498754, _omitFieldNames ? '' : 'cloudhsmclusterid')
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..aE<ExpirationModelType>(
        113933558, _omitFieldNames ? '' : 'expirationmodel',
        enumValues: ExpirationModelType.values)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aE<KeySpec>(138220928, _omitFieldNames ? '' : 'keyspec',
        enumValues: KeySpec.values)
    ..aOS(183721586, _omitFieldNames ? '' : 'currentkeymaterialid')
    ..pc<EncryptionAlgorithmSpec>(194511375,
        _omitFieldNames ? '' : 'encryptionalgorithms', $pb.PbFieldType.KE,
        valueOf: EncryptionAlgorithmSpec.valueOf,
        enumValues: EncryptionAlgorithmSpec.values,
        defaultEnumValue:
            EncryptionAlgorithmSpec.ENCRYPTION_ALGORITHM_SPEC_SM2PKE)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..pc<KeyAgreementAlgorithmSpec>(328746163,
        _omitFieldNames ? '' : 'keyagreementalgorithms', $pb.PbFieldType.KE,
        valueOf: KeyAgreementAlgorithmSpec.valueOf,
        enumValues: KeyAgreementAlgorithmSpec.values,
        defaultEnumValue:
            KeyAgreementAlgorithmSpec.KEY_AGREEMENT_ALGORITHM_SPEC_ECDH)
    ..aOS(347845564, _omitFieldNames ? '' : 'deletiondate')
    ..aE<KeyUsageType>(357216772, _omitFieldNames ? '' : 'keyusage',
        enumValues: KeyUsageType.values)
    ..aOM<XksKeyConfigurationType>(
        359766455, _omitFieldNames ? '' : 'xkskeyconfiguration',
        subBuilder: XksKeyConfigurationType.create)
    ..aOS(370093421, _omitFieldNames ? '' : 'awsaccountid')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOB(405769103, _omitFieldNames ? '' : 'multiregion')
    ..aE<KeyManagerType>(432733972, _omitFieldNames ? '' : 'keymanager',
        enumValues: KeyManagerType.values)
    ..aOM<MultiRegionConfiguration>(
        460103813, _omitFieldNames ? '' : 'multiregionconfiguration',
        subBuilder: MultiRegionConfiguration.create)
    ..aE<CustomerMasterKeySpec>(
        472470930, _omitFieldNames ? '' : 'customermasterkeyspec',
        enumValues: CustomerMasterKeySpec.values)
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..aI(480447795, _omitFieldNames ? '' : 'pendingdeletionwindowindays')
    ..pc<SigningAlgorithmSpec>(508241975,
        _omitFieldNames ? '' : 'signingalgorithms', $pb.PbFieldType.KE,
        valueOf: SigningAlgorithmSpec.valueOf,
        enumValues: SigningAlgorithmSpec.values,
        defaultEnumValue:
            SigningAlgorithmSpec.SIGNING_ALGORITHM_SPEC_ECDSA_SHA_384)
    ..aOS(522718673, _omitFieldNames ? '' : 'validto')
    ..aE<OriginType>(529732720, _omitFieldNames ? '' : 'origin',
        enumValues: OriginType.values)
    ..pc<MacAlgorithmSpec>(
        532181443, _omitFieldNames ? '' : 'macalgorithms', $pb.PbFieldType.KE,
        valueOf: MacAlgorithmSpec.valueOf,
        enumValues: MacAlgorithmSpec.values,
        defaultEnumValue: MacAlgorithmSpec.MAC_ALGORITHM_SPEC_HMAC_SHA_512)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeyMetadata clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeyMetadata copyWith(void Function(KeyMetadata) updates) =>
      super.copyWith((message) => updates(message as KeyMetadata))
          as KeyMetadata;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeyMetadata create() => KeyMetadata._();
  @$core.override
  KeyMetadata createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeyMetadata getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeyMetadata>(create);
  static KeyMetadata? _defaultInstance;

  @$pb.TagNumber(27894226)
  KeyState get keystate => $_getN(0);
  @$pb.TagNumber(27894226)
  set keystate(KeyState value) => $_setField(27894226, value);
  @$pb.TagNumber(27894226)
  $core.bool hasKeystate() => $_has(0);
  @$pb.TagNumber(27894226)
  void clearKeystate() => $_clearField(27894226);

  @$pb.TagNumber(56498754)
  $core.String get cloudhsmclusterid => $_getSZ(1);
  @$pb.TagNumber(56498754)
  set cloudhsmclusterid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(56498754)
  $core.bool hasCloudhsmclusterid() => $_has(1);
  @$pb.TagNumber(56498754)
  void clearCloudhsmclusterid() => $_clearField(56498754);

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(2);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(2);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);

  @$pb.TagNumber(113933558)
  ExpirationModelType get expirationmodel => $_getN(3);
  @$pb.TagNumber(113933558)
  set expirationmodel(ExpirationModelType value) =>
      $_setField(113933558, value);
  @$pb.TagNumber(113933558)
  $core.bool hasExpirationmodel() => $_has(3);
  @$pb.TagNumber(113933558)
  void clearExpirationmodel() => $_clearField(113933558);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(138220928)
  KeySpec get keyspec => $_getN(5);
  @$pb.TagNumber(138220928)
  set keyspec(KeySpec value) => $_setField(138220928, value);
  @$pb.TagNumber(138220928)
  $core.bool hasKeyspec() => $_has(5);
  @$pb.TagNumber(138220928)
  void clearKeyspec() => $_clearField(138220928);

  @$pb.TagNumber(183721586)
  $core.String get currentkeymaterialid => $_getSZ(6);
  @$pb.TagNumber(183721586)
  set currentkeymaterialid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(183721586)
  $core.bool hasCurrentkeymaterialid() => $_has(6);
  @$pb.TagNumber(183721586)
  void clearCurrentkeymaterialid() => $_clearField(183721586);

  @$pb.TagNumber(194511375)
  $pb.PbList<EncryptionAlgorithmSpec> get encryptionalgorithms => $_getList(7);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(8);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(8);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(9);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(9, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(9);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(328746163)
  $pb.PbList<KeyAgreementAlgorithmSpec> get keyagreementalgorithms =>
      $_getList(10);

  @$pb.TagNumber(347845564)
  $core.String get deletiondate => $_getSZ(11);
  @$pb.TagNumber(347845564)
  set deletiondate($core.String value) => $_setString(11, value);
  @$pb.TagNumber(347845564)
  $core.bool hasDeletiondate() => $_has(11);
  @$pb.TagNumber(347845564)
  void clearDeletiondate() => $_clearField(347845564);

  @$pb.TagNumber(357216772)
  KeyUsageType get keyusage => $_getN(12);
  @$pb.TagNumber(357216772)
  set keyusage(KeyUsageType value) => $_setField(357216772, value);
  @$pb.TagNumber(357216772)
  $core.bool hasKeyusage() => $_has(12);
  @$pb.TagNumber(357216772)
  void clearKeyusage() => $_clearField(357216772);

  @$pb.TagNumber(359766455)
  XksKeyConfigurationType get xkskeyconfiguration => $_getN(13);
  @$pb.TagNumber(359766455)
  set xkskeyconfiguration(XksKeyConfigurationType value) =>
      $_setField(359766455, value);
  @$pb.TagNumber(359766455)
  $core.bool hasXkskeyconfiguration() => $_has(13);
  @$pb.TagNumber(359766455)
  void clearXkskeyconfiguration() => $_clearField(359766455);
  @$pb.TagNumber(359766455)
  XksKeyConfigurationType ensureXkskeyconfiguration() => $_ensure(13);

  @$pb.TagNumber(370093421)
  $core.String get awsaccountid => $_getSZ(14);
  @$pb.TagNumber(370093421)
  set awsaccountid($core.String value) => $_setString(14, value);
  @$pb.TagNumber(370093421)
  $core.bool hasAwsaccountid() => $_has(14);
  @$pb.TagNumber(370093421)
  void clearAwsaccountid() => $_clearField(370093421);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(15);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(15, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(15);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(405769103)
  $core.bool get multiregion => $_getBF(16);
  @$pb.TagNumber(405769103)
  set multiregion($core.bool value) => $_setBool(16, value);
  @$pb.TagNumber(405769103)
  $core.bool hasMultiregion() => $_has(16);
  @$pb.TagNumber(405769103)
  void clearMultiregion() => $_clearField(405769103);

  @$pb.TagNumber(432733972)
  KeyManagerType get keymanager => $_getN(17);
  @$pb.TagNumber(432733972)
  set keymanager(KeyManagerType value) => $_setField(432733972, value);
  @$pb.TagNumber(432733972)
  $core.bool hasKeymanager() => $_has(17);
  @$pb.TagNumber(432733972)
  void clearKeymanager() => $_clearField(432733972);

  @$pb.TagNumber(460103813)
  MultiRegionConfiguration get multiregionconfiguration => $_getN(18);
  @$pb.TagNumber(460103813)
  set multiregionconfiguration(MultiRegionConfiguration value) =>
      $_setField(460103813, value);
  @$pb.TagNumber(460103813)
  $core.bool hasMultiregionconfiguration() => $_has(18);
  @$pb.TagNumber(460103813)
  void clearMultiregionconfiguration() => $_clearField(460103813);
  @$pb.TagNumber(460103813)
  MultiRegionConfiguration ensureMultiregionconfiguration() => $_ensure(18);

  @$pb.TagNumber(472470930)
  CustomerMasterKeySpec get customermasterkeyspec => $_getN(19);
  @$pb.TagNumber(472470930)
  set customermasterkeyspec(CustomerMasterKeySpec value) =>
      $_setField(472470930, value);
  @$pb.TagNumber(472470930)
  $core.bool hasCustomermasterkeyspec() => $_has(19);
  @$pb.TagNumber(472470930)
  void clearCustomermasterkeyspec() => $_clearField(472470930);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(20);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(20, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(20);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);

  @$pb.TagNumber(480447795)
  $core.int get pendingdeletionwindowindays => $_getIZ(21);
  @$pb.TagNumber(480447795)
  set pendingdeletionwindowindays($core.int value) =>
      $_setSignedInt32(21, value);
  @$pb.TagNumber(480447795)
  $core.bool hasPendingdeletionwindowindays() => $_has(21);
  @$pb.TagNumber(480447795)
  void clearPendingdeletionwindowindays() => $_clearField(480447795);

  @$pb.TagNumber(508241975)
  $pb.PbList<SigningAlgorithmSpec> get signingalgorithms => $_getList(22);

  @$pb.TagNumber(522718673)
  $core.String get validto => $_getSZ(23);
  @$pb.TagNumber(522718673)
  set validto($core.String value) => $_setString(23, value);
  @$pb.TagNumber(522718673)
  $core.bool hasValidto() => $_has(23);
  @$pb.TagNumber(522718673)
  void clearValidto() => $_clearField(522718673);

  @$pb.TagNumber(529732720)
  OriginType get origin => $_getN(24);
  @$pb.TagNumber(529732720)
  set origin(OriginType value) => $_setField(529732720, value);
  @$pb.TagNumber(529732720)
  $core.bool hasOrigin() => $_has(24);
  @$pb.TagNumber(529732720)
  void clearOrigin() => $_clearField(529732720);

  @$pb.TagNumber(532181443)
  $pb.PbList<MacAlgorithmSpec> get macalgorithms => $_getList(25);
}

class KeyUnavailableException extends $pb.GeneratedMessage {
  factory KeyUnavailableException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KeyUnavailableException._();

  factory KeyUnavailableException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeyUnavailableException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeyUnavailableException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeyUnavailableException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeyUnavailableException copyWith(
          void Function(KeyUnavailableException) updates) =>
      super.copyWith((message) => updates(message as KeyUnavailableException))
          as KeyUnavailableException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeyUnavailableException create() => KeyUnavailableException._();
  @$core.override
  KeyUnavailableException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeyUnavailableException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeyUnavailableException>(create);
  static KeyUnavailableException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class LimitExceededException extends $pb.GeneratedMessage {
  factory LimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  LimitExceededException._();

  factory LimitExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LimitExceededException copyWith(
          void Function(LimitExceededException) updates) =>
      super.copyWith((message) => updates(message as LimitExceededException))
          as LimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LimitExceededException create() => LimitExceededException._();
  @$core.override
  LimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LimitExceededException>(create);
  static LimitExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ListAliasesRequest extends $pb.GeneratedMessage {
  factory ListAliasesRequest({
    $core.String? marker,
    $core.String? keyid,
    $core.int? limit,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (keyid != null) result.keyid = keyid;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListAliasesRequest._();

  factory ListAliasesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAliasesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAliasesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAliasesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAliasesRequest copyWith(void Function(ListAliasesRequest) updates) =>
      super.copyWith((message) => updates(message as ListAliasesRequest))
          as ListAliasesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAliasesRequest create() => ListAliasesRequest._();
  @$core.override
  ListAliasesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAliasesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListAliasesRequest>(create);
  static ListAliasesRequest? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListAliasesResponse extends $pb.GeneratedMessage {
  factory ListAliasesResponse({
    $core.bool? truncated,
    $core.Iterable<AliasListEntry>? aliases,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (truncated != null) result.truncated = truncated;
    if (aliases != null) result.aliases.addAll(aliases);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListAliasesResponse._();

  factory ListAliasesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAliasesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAliasesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(152451018, _omitFieldNames ? '' : 'truncated')
    ..pPM<AliasListEntry>(476693696, _omitFieldNames ? '' : 'aliases',
        subBuilder: AliasListEntry.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAliasesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAliasesResponse copyWith(void Function(ListAliasesResponse) updates) =>
      super.copyWith((message) => updates(message as ListAliasesResponse))
          as ListAliasesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAliasesResponse create() => ListAliasesResponse._();
  @$core.override
  ListAliasesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAliasesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListAliasesResponse>(create);
  static ListAliasesResponse? _defaultInstance;

  @$pb.TagNumber(152451018)
  $core.bool get truncated => $_getBF(0);
  @$pb.TagNumber(152451018)
  set truncated($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(152451018)
  $core.bool hasTruncated() => $_has(0);
  @$pb.TagNumber(152451018)
  void clearTruncated() => $_clearField(152451018);

  @$pb.TagNumber(476693696)
  $pb.PbList<AliasListEntry> get aliases => $_getList(1);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListGrantsRequest extends $pb.GeneratedMessage {
  factory ListGrantsRequest({
    $core.String? grantid,
    $core.String? marker,
    $core.String? granteeprincipal,
    $core.String? keyid,
    $core.int? limit,
  }) {
    final result = create();
    if (grantid != null) result.grantid = grantid;
    if (marker != null) result.marker = marker;
    if (granteeprincipal != null) result.granteeprincipal = granteeprincipal;
    if (keyid != null) result.keyid = keyid;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListGrantsRequest._();

  factory ListGrantsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListGrantsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListGrantsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(66852281, _omitFieldNames ? '' : 'grantid')
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(234727364, _omitFieldNames ? '' : 'granteeprincipal')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGrantsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGrantsRequest copyWith(void Function(ListGrantsRequest) updates) =>
      super.copyWith((message) => updates(message as ListGrantsRequest))
          as ListGrantsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListGrantsRequest create() => ListGrantsRequest._();
  @$core.override
  ListGrantsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListGrantsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListGrantsRequest>(create);
  static ListGrantsRequest? _defaultInstance;

  @$pb.TagNumber(66852281)
  $core.String get grantid => $_getSZ(0);
  @$pb.TagNumber(66852281)
  set grantid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(66852281)
  $core.bool hasGrantid() => $_has(0);
  @$pb.TagNumber(66852281)
  void clearGrantid() => $_clearField(66852281);

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(1);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(1);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(234727364)
  $core.String get granteeprincipal => $_getSZ(2);
  @$pb.TagNumber(234727364)
  set granteeprincipal($core.String value) => $_setString(2, value);
  @$pb.TagNumber(234727364)
  $core.bool hasGranteeprincipal() => $_has(2);
  @$pb.TagNumber(234727364)
  void clearGranteeprincipal() => $_clearField(234727364);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(4);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(4);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListGrantsResponse extends $pb.GeneratedMessage {
  factory ListGrantsResponse({
    $core.bool? truncated,
    $core.Iterable<GrantListEntry>? grants,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (truncated != null) result.truncated = truncated;
    if (grants != null) result.grants.addAll(grants);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListGrantsResponse._();

  factory ListGrantsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListGrantsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListGrantsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(152451018, _omitFieldNames ? '' : 'truncated')
    ..pPM<GrantListEntry>(226910555, _omitFieldNames ? '' : 'grants',
        subBuilder: GrantListEntry.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGrantsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGrantsResponse copyWith(void Function(ListGrantsResponse) updates) =>
      super.copyWith((message) => updates(message as ListGrantsResponse))
          as ListGrantsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListGrantsResponse create() => ListGrantsResponse._();
  @$core.override
  ListGrantsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListGrantsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListGrantsResponse>(create);
  static ListGrantsResponse? _defaultInstance;

  @$pb.TagNumber(152451018)
  $core.bool get truncated => $_getBF(0);
  @$pb.TagNumber(152451018)
  set truncated($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(152451018)
  $core.bool hasTruncated() => $_has(0);
  @$pb.TagNumber(152451018)
  void clearTruncated() => $_clearField(152451018);

  @$pb.TagNumber(226910555)
  $pb.PbList<GrantListEntry> get grants => $_getList(1);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListKeyPoliciesRequest extends $pb.GeneratedMessage {
  factory ListKeyPoliciesRequest({
    $core.String? marker,
    $core.String? keyid,
    $core.int? limit,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (keyid != null) result.keyid = keyid;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListKeyPoliciesRequest._();

  factory ListKeyPoliciesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListKeyPoliciesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListKeyPoliciesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeyPoliciesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeyPoliciesRequest copyWith(
          void Function(ListKeyPoliciesRequest) updates) =>
      super.copyWith((message) => updates(message as ListKeyPoliciesRequest))
          as ListKeyPoliciesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListKeyPoliciesRequest create() => ListKeyPoliciesRequest._();
  @$core.override
  ListKeyPoliciesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListKeyPoliciesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListKeyPoliciesRequest>(create);
  static ListKeyPoliciesRequest? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListKeyPoliciesResponse extends $pb.GeneratedMessage {
  factory ListKeyPoliciesResponse({
    $core.bool? truncated,
    $core.Iterable<$core.String>? policynames,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (truncated != null) result.truncated = truncated;
    if (policynames != null) result.policynames.addAll(policynames);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListKeyPoliciesResponse._();

  factory ListKeyPoliciesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListKeyPoliciesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListKeyPoliciesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(152451018, _omitFieldNames ? '' : 'truncated')
    ..pPS(264098782, _omitFieldNames ? '' : 'policynames')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeyPoliciesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeyPoliciesResponse copyWith(
          void Function(ListKeyPoliciesResponse) updates) =>
      super.copyWith((message) => updates(message as ListKeyPoliciesResponse))
          as ListKeyPoliciesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListKeyPoliciesResponse create() => ListKeyPoliciesResponse._();
  @$core.override
  ListKeyPoliciesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListKeyPoliciesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListKeyPoliciesResponse>(create);
  static ListKeyPoliciesResponse? _defaultInstance;

  @$pb.TagNumber(152451018)
  $core.bool get truncated => $_getBF(0);
  @$pb.TagNumber(152451018)
  set truncated($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(152451018)
  $core.bool hasTruncated() => $_has(0);
  @$pb.TagNumber(152451018)
  void clearTruncated() => $_clearField(152451018);

  @$pb.TagNumber(264098782)
  $pb.PbList<$core.String> get policynames => $_getList(1);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListKeyRotationsRequest extends $pb.GeneratedMessage {
  factory ListKeyRotationsRequest({
    $core.String? marker,
    $core.String? keyid,
    $core.int? limit,
    IncludeKeyMaterial? includekeymaterial,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (keyid != null) result.keyid = keyid;
    if (limit != null) result.limit = limit;
    if (includekeymaterial != null)
      result.includekeymaterial = includekeymaterial;
    return result;
  }

  ListKeyRotationsRequest._();

  factory ListKeyRotationsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListKeyRotationsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListKeyRotationsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aE<IncludeKeyMaterial>(
        531559428, _omitFieldNames ? '' : 'includekeymaterial',
        enumValues: IncludeKeyMaterial.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeyRotationsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeyRotationsRequest copyWith(
          void Function(ListKeyRotationsRequest) updates) =>
      super.copyWith((message) => updates(message as ListKeyRotationsRequest))
          as ListKeyRotationsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListKeyRotationsRequest create() => ListKeyRotationsRequest._();
  @$core.override
  ListKeyRotationsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListKeyRotationsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListKeyRotationsRequest>(create);
  static ListKeyRotationsRequest? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531559428)
  IncludeKeyMaterial get includekeymaterial => $_getN(3);
  @$pb.TagNumber(531559428)
  set includekeymaterial(IncludeKeyMaterial value) =>
      $_setField(531559428, value);
  @$pb.TagNumber(531559428)
  $core.bool hasIncludekeymaterial() => $_has(3);
  @$pb.TagNumber(531559428)
  void clearIncludekeymaterial() => $_clearField(531559428);
}

class ListKeyRotationsResponse extends $pb.GeneratedMessage {
  factory ListKeyRotationsResponse({
    $core.Iterable<RotationsListEntry>? rotations,
    $core.bool? truncated,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (rotations != null) result.rotations.addAll(rotations);
    if (truncated != null) result.truncated = truncated;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListKeyRotationsResponse._();

  factory ListKeyRotationsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListKeyRotationsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListKeyRotationsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..pPM<RotationsListEntry>(24209381, _omitFieldNames ? '' : 'rotations',
        subBuilder: RotationsListEntry.create)
    ..aOB(152451018, _omitFieldNames ? '' : 'truncated')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeyRotationsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeyRotationsResponse copyWith(
          void Function(ListKeyRotationsResponse) updates) =>
      super.copyWith((message) => updates(message as ListKeyRotationsResponse))
          as ListKeyRotationsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListKeyRotationsResponse create() => ListKeyRotationsResponse._();
  @$core.override
  ListKeyRotationsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListKeyRotationsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListKeyRotationsResponse>(create);
  static ListKeyRotationsResponse? _defaultInstance;

  @$pb.TagNumber(24209381)
  $pb.PbList<RotationsListEntry> get rotations => $_getList(0);

  @$pb.TagNumber(152451018)
  $core.bool get truncated => $_getBF(1);
  @$pb.TagNumber(152451018)
  set truncated($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(152451018)
  $core.bool hasTruncated() => $_has(1);
  @$pb.TagNumber(152451018)
  void clearTruncated() => $_clearField(152451018);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListKeysRequest extends $pb.GeneratedMessage {
  factory ListKeysRequest({
    $core.String? marker,
    $core.int? limit,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListKeysRequest._();

  factory ListKeysRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListKeysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListKeysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeysRequest copyWith(void Function(ListKeysRequest) updates) =>
      super.copyWith((message) => updates(message as ListKeysRequest))
          as ListKeysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListKeysRequest create() => ListKeysRequest._();
  @$core.override
  ListKeysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListKeysRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListKeysRequest>(create);
  static ListKeysRequest? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListKeysResponse extends $pb.GeneratedMessage {
  factory ListKeysResponse({
    $core.Iterable<KeyListEntry>? keys,
    $core.bool? truncated,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (keys != null) result.keys.addAll(keys);
    if (truncated != null) result.truncated = truncated;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListKeysResponse._();

  factory ListKeysResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListKeysResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListKeysResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..pPM<KeyListEntry>(2831086, _omitFieldNames ? '' : 'keys',
        subBuilder: KeyListEntry.create)
    ..aOB(152451018, _omitFieldNames ? '' : 'truncated')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeysResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListKeysResponse copyWith(void Function(ListKeysResponse) updates) =>
      super.copyWith((message) => updates(message as ListKeysResponse))
          as ListKeysResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListKeysResponse create() => ListKeysResponse._();
  @$core.override
  ListKeysResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListKeysResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListKeysResponse>(create);
  static ListKeysResponse? _defaultInstance;

  @$pb.TagNumber(2831086)
  $pb.PbList<KeyListEntry> get keys => $_getList(0);

  @$pb.TagNumber(152451018)
  $core.bool get truncated => $_getBF(1);
  @$pb.TagNumber(152451018)
  set truncated($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(152451018)
  $core.bool hasTruncated() => $_has(1);
  @$pb.TagNumber(152451018)
  void clearTruncated() => $_clearField(152451018);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListResourceTagsRequest extends $pb.GeneratedMessage {
  factory ListResourceTagsRequest({
    $core.String? marker,
    $core.String? keyid,
    $core.int? limit,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (keyid != null) result.keyid = keyid;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListResourceTagsRequest._();

  factory ListResourceTagsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListResourceTagsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListResourceTagsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourceTagsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourceTagsRequest copyWith(
          void Function(ListResourceTagsRequest) updates) =>
      super.copyWith((message) => updates(message as ListResourceTagsRequest))
          as ListResourceTagsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListResourceTagsRequest create() => ListResourceTagsRequest._();
  @$core.override
  ListResourceTagsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListResourceTagsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListResourceTagsRequest>(create);
  static ListResourceTagsRequest? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListResourceTagsResponse extends $pb.GeneratedMessage {
  factory ListResourceTagsResponse({
    $core.bool? truncated,
    $core.Iterable<Tag>? tags,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (truncated != null) result.truncated = truncated;
    if (tags != null) result.tags.addAll(tags);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListResourceTagsResponse._();

  factory ListResourceTagsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListResourceTagsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListResourceTagsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(152451018, _omitFieldNames ? '' : 'truncated')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourceTagsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourceTagsResponse copyWith(
          void Function(ListResourceTagsResponse) updates) =>
      super.copyWith((message) => updates(message as ListResourceTagsResponse))
          as ListResourceTagsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListResourceTagsResponse create() => ListResourceTagsResponse._();
  @$core.override
  ListResourceTagsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListResourceTagsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListResourceTagsResponse>(create);
  static ListResourceTagsResponse? _defaultInstance;

  @$pb.TagNumber(152451018)
  $core.bool get truncated => $_getBF(0);
  @$pb.TagNumber(152451018)
  set truncated($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(152451018)
  $core.bool hasTruncated() => $_has(0);
  @$pb.TagNumber(152451018)
  void clearTruncated() => $_clearField(152451018);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListRetirableGrantsRequest extends $pb.GeneratedMessage {
  factory ListRetirableGrantsRequest({
    $core.String? retiringprincipal,
    $core.String? marker,
    $core.int? limit,
  }) {
    final result = create();
    if (retiringprincipal != null) result.retiringprincipal = retiringprincipal;
    if (marker != null) result.marker = marker;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListRetirableGrantsRequest._();

  factory ListRetirableGrantsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRetirableGrantsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRetirableGrantsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(49541086, _omitFieldNames ? '' : 'retiringprincipal')
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRetirableGrantsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRetirableGrantsRequest copyWith(
          void Function(ListRetirableGrantsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListRetirableGrantsRequest))
          as ListRetirableGrantsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRetirableGrantsRequest create() => ListRetirableGrantsRequest._();
  @$core.override
  ListRetirableGrantsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRetirableGrantsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRetirableGrantsRequest>(create);
  static ListRetirableGrantsRequest? _defaultInstance;

  @$pb.TagNumber(49541086)
  $core.String get retiringprincipal => $_getSZ(0);
  @$pb.TagNumber(49541086)
  set retiringprincipal($core.String value) => $_setString(0, value);
  @$pb.TagNumber(49541086)
  $core.bool hasRetiringprincipal() => $_has(0);
  @$pb.TagNumber(49541086)
  void clearRetiringprincipal() => $_clearField(49541086);

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(1);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(1);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class MalformedPolicyDocumentException extends $pb.GeneratedMessage {
  factory MalformedPolicyDocumentException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  MalformedPolicyDocumentException._();

  factory MalformedPolicyDocumentException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MalformedPolicyDocumentException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MalformedPolicyDocumentException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MalformedPolicyDocumentException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MalformedPolicyDocumentException copyWith(
          void Function(MalformedPolicyDocumentException) updates) =>
      super.copyWith(
              (message) => updates(message as MalformedPolicyDocumentException))
          as MalformedPolicyDocumentException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MalformedPolicyDocumentException create() =>
      MalformedPolicyDocumentException._();
  @$core.override
  MalformedPolicyDocumentException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MalformedPolicyDocumentException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MalformedPolicyDocumentException>(
          create);
  static MalformedPolicyDocumentException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class MultiRegionConfiguration extends $pb.GeneratedMessage {
  factory MultiRegionConfiguration({
    MultiRegionKey? primarykey,
    $core.Iterable<MultiRegionKey>? replicakeys,
    MultiRegionKeyType? multiregionkeytype,
  }) {
    final result = create();
    if (primarykey != null) result.primarykey = primarykey;
    if (replicakeys != null) result.replicakeys.addAll(replicakeys);
    if (multiregionkeytype != null)
      result.multiregionkeytype = multiregionkeytype;
    return result;
  }

  MultiRegionConfiguration._();

  factory MultiRegionConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MultiRegionConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MultiRegionConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOM<MultiRegionKey>(174019825, _omitFieldNames ? '' : 'primarykey',
        subBuilder: MultiRegionKey.create)
    ..pPM<MultiRegionKey>(266801624, _omitFieldNames ? '' : 'replicakeys',
        subBuilder: MultiRegionKey.create)
    ..aE<MultiRegionKeyType>(
        483927198, _omitFieldNames ? '' : 'multiregionkeytype',
        enumValues: MultiRegionKeyType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MultiRegionConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MultiRegionConfiguration copyWith(
          void Function(MultiRegionConfiguration) updates) =>
      super.copyWith((message) => updates(message as MultiRegionConfiguration))
          as MultiRegionConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MultiRegionConfiguration create() => MultiRegionConfiguration._();
  @$core.override
  MultiRegionConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MultiRegionConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MultiRegionConfiguration>(create);
  static MultiRegionConfiguration? _defaultInstance;

  @$pb.TagNumber(174019825)
  MultiRegionKey get primarykey => $_getN(0);
  @$pb.TagNumber(174019825)
  set primarykey(MultiRegionKey value) => $_setField(174019825, value);
  @$pb.TagNumber(174019825)
  $core.bool hasPrimarykey() => $_has(0);
  @$pb.TagNumber(174019825)
  void clearPrimarykey() => $_clearField(174019825);
  @$pb.TagNumber(174019825)
  MultiRegionKey ensurePrimarykey() => $_ensure(0);

  @$pb.TagNumber(266801624)
  $pb.PbList<MultiRegionKey> get replicakeys => $_getList(1);

  @$pb.TagNumber(483927198)
  MultiRegionKeyType get multiregionkeytype => $_getN(2);
  @$pb.TagNumber(483927198)
  set multiregionkeytype(MultiRegionKeyType value) =>
      $_setField(483927198, value);
  @$pb.TagNumber(483927198)
  $core.bool hasMultiregionkeytype() => $_has(2);
  @$pb.TagNumber(483927198)
  void clearMultiregionkeytype() => $_clearField(483927198);
}

class MultiRegionKey extends $pb.GeneratedMessage {
  factory MultiRegionKey({
    $core.String? region,
    $core.String? arn,
  }) {
    final result = create();
    if (region != null) result.region = region;
    if (arn != null) result.arn = arn;
    return result;
  }

  MultiRegionKey._();

  factory MultiRegionKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MultiRegionKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MultiRegionKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(154040478, _omitFieldNames ? '' : 'region')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MultiRegionKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MultiRegionKey copyWith(void Function(MultiRegionKey) updates) =>
      super.copyWith((message) => updates(message as MultiRegionKey))
          as MultiRegionKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MultiRegionKey create() => MultiRegionKey._();
  @$core.override
  MultiRegionKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MultiRegionKey getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MultiRegionKey>(create);
  static MultiRegionKey? _defaultInstance;

  @$pb.TagNumber(154040478)
  $core.String get region => $_getSZ(0);
  @$pb.TagNumber(154040478)
  set region($core.String value) => $_setString(0, value);
  @$pb.TagNumber(154040478)
  $core.bool hasRegion() => $_has(0);
  @$pb.TagNumber(154040478)
  void clearRegion() => $_clearField(154040478);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
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

class PutKeyPolicyRequest extends $pb.GeneratedMessage {
  factory PutKeyPolicyRequest({
    $core.bool? bypasspolicylockoutsafetycheck,
    $core.String? policyname,
    $core.String? keyid,
    $core.String? policy,
  }) {
    final result = create();
    if (bypasspolicylockoutsafetycheck != null)
      result.bypasspolicylockoutsafetycheck = bypasspolicylockoutsafetycheck;
    if (policyname != null) result.policyname = policyname;
    if (keyid != null) result.keyid = keyid;
    if (policy != null) result.policy = policy;
    return result;
  }

  PutKeyPolicyRequest._();

  factory PutKeyPolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutKeyPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutKeyPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(177450851, _omitFieldNames ? '' : 'bypasspolicylockoutsafetycheck')
    ..aOS(266468029, _omitFieldNames ? '' : 'policyname')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutKeyPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutKeyPolicyRequest copyWith(void Function(PutKeyPolicyRequest) updates) =>
      super.copyWith((message) => updates(message as PutKeyPolicyRequest))
          as PutKeyPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutKeyPolicyRequest create() => PutKeyPolicyRequest._();
  @$core.override
  PutKeyPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutKeyPolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutKeyPolicyRequest>(create);
  static PutKeyPolicyRequest? _defaultInstance;

  @$pb.TagNumber(177450851)
  $core.bool get bypasspolicylockoutsafetycheck => $_getBF(0);
  @$pb.TagNumber(177450851)
  set bypasspolicylockoutsafetycheck($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(177450851)
  $core.bool hasBypasspolicylockoutsafetycheck() => $_has(0);
  @$pb.TagNumber(177450851)
  void clearBypasspolicylockoutsafetycheck() => $_clearField(177450851);

  @$pb.TagNumber(266468029)
  $core.String get policyname => $_getSZ(1);
  @$pb.TagNumber(266468029)
  set policyname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266468029)
  $core.bool hasPolicyname() => $_has(1);
  @$pb.TagNumber(266468029)
  void clearPolicyname() => $_clearField(266468029);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(3);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(3, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(3);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class ReEncryptRequest extends $pb.GeneratedMessage {
  factory ReEncryptRequest({
    EncryptionAlgorithmSpec? destinationencryptionalgorithm,
    $core.Iterable<DryRunModifierType>? dryrunmodifiers,
    $core.bool? dryrun,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        destinationencryptioncontext,
    $core.String? sourcekeyid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        sourceencryptioncontext,
    EncryptionAlgorithmSpec? sourceencryptionalgorithm,
    $core.List<$core.int>? ciphertextblob,
    $core.Iterable<$core.String>? granttokens,
    $core.String? destinationkeyid,
  }) {
    final result = create();
    if (destinationencryptionalgorithm != null)
      result.destinationencryptionalgorithm = destinationencryptionalgorithm;
    if (dryrunmodifiers != null) result.dryrunmodifiers.addAll(dryrunmodifiers);
    if (dryrun != null) result.dryrun = dryrun;
    if (destinationencryptioncontext != null)
      result.destinationencryptioncontext
          .addEntries(destinationencryptioncontext);
    if (sourcekeyid != null) result.sourcekeyid = sourcekeyid;
    if (sourceencryptioncontext != null)
      result.sourceencryptioncontext.addEntries(sourceencryptioncontext);
    if (sourceencryptionalgorithm != null)
      result.sourceencryptionalgorithm = sourceencryptionalgorithm;
    if (ciphertextblob != null) result.ciphertextblob = ciphertextblob;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    if (destinationkeyid != null) result.destinationkeyid = destinationkeyid;
    return result;
  }

  ReEncryptRequest._();

  factory ReEncryptRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReEncryptRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReEncryptRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<EncryptionAlgorithmSpec>(
        3500944, _omitFieldNames ? '' : 'destinationencryptionalgorithm',
        enumValues: EncryptionAlgorithmSpec.values)
    ..pc<DryRunModifierType>(
        24346424, _omitFieldNames ? '' : 'dryrunmodifiers', $pb.PbFieldType.KE,
        valueOf: DryRunModifierType.valueOf,
        enumValues: DryRunModifierType.values,
        defaultEnumValue:
            DryRunModifierType.DRY_RUN_MODIFIER_TYPE_IGNORE_CIPHERTEXT)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..m<$core.String, $core.String>(
        116435710, _omitFieldNames ? '' : 'destinationencryptioncontext',
        entryClassName: 'ReEncryptRequest.DestinationencryptioncontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..aOS(137105771, _omitFieldNames ? '' : 'sourcekeyid')
    ..m<$core.String, $core.String>(
        203785001, _omitFieldNames ? '' : 'sourceencryptioncontext',
        entryClassName: 'ReEncryptRequest.SourceencryptioncontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('kms'))
    ..aE<EncryptionAlgorithmSpec>(
        283215847, _omitFieldNames ? '' : 'sourceencryptionalgorithm',
        enumValues: EncryptionAlgorithmSpec.values)
    ..a<$core.List<$core.int>>(
        338198183, _omitFieldNames ? '' : 'ciphertextblob', $pb.PbFieldType.OY)
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..aOS(396219964, _omitFieldNames ? '' : 'destinationkeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReEncryptRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReEncryptRequest copyWith(void Function(ReEncryptRequest) updates) =>
      super.copyWith((message) => updates(message as ReEncryptRequest))
          as ReEncryptRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReEncryptRequest create() => ReEncryptRequest._();
  @$core.override
  ReEncryptRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReEncryptRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReEncryptRequest>(create);
  static ReEncryptRequest? _defaultInstance;

  @$pb.TagNumber(3500944)
  EncryptionAlgorithmSpec get destinationencryptionalgorithm => $_getN(0);
  @$pb.TagNumber(3500944)
  set destinationencryptionalgorithm(EncryptionAlgorithmSpec value) =>
      $_setField(3500944, value);
  @$pb.TagNumber(3500944)
  $core.bool hasDestinationencryptionalgorithm() => $_has(0);
  @$pb.TagNumber(3500944)
  void clearDestinationencryptionalgorithm() => $_clearField(3500944);

  @$pb.TagNumber(24346424)
  $pb.PbList<DryRunModifierType> get dryrunmodifiers => $_getList(1);

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(2);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(2);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(116435710)
  $pb.PbMap<$core.String, $core.String> get destinationencryptioncontext =>
      $_getMap(3);

  @$pb.TagNumber(137105771)
  $core.String get sourcekeyid => $_getSZ(4);
  @$pb.TagNumber(137105771)
  set sourcekeyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(137105771)
  $core.bool hasSourcekeyid() => $_has(4);
  @$pb.TagNumber(137105771)
  void clearSourcekeyid() => $_clearField(137105771);

  @$pb.TagNumber(203785001)
  $pb.PbMap<$core.String, $core.String> get sourceencryptioncontext =>
      $_getMap(5);

  @$pb.TagNumber(283215847)
  EncryptionAlgorithmSpec get sourceencryptionalgorithm => $_getN(6);
  @$pb.TagNumber(283215847)
  set sourceencryptionalgorithm(EncryptionAlgorithmSpec value) =>
      $_setField(283215847, value);
  @$pb.TagNumber(283215847)
  $core.bool hasSourceencryptionalgorithm() => $_has(6);
  @$pb.TagNumber(283215847)
  void clearSourceencryptionalgorithm() => $_clearField(283215847);

  @$pb.TagNumber(338198183)
  $core.List<$core.int> get ciphertextblob => $_getN(7);
  @$pb.TagNumber(338198183)
  set ciphertextblob($core.List<$core.int> value) => $_setBytes(7, value);
  @$pb.TagNumber(338198183)
  $core.bool hasCiphertextblob() => $_has(7);
  @$pb.TagNumber(338198183)
  void clearCiphertextblob() => $_clearField(338198183);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(8);

  @$pb.TagNumber(396219964)
  $core.String get destinationkeyid => $_getSZ(9);
  @$pb.TagNumber(396219964)
  set destinationkeyid($core.String value) => $_setString(9, value);
  @$pb.TagNumber(396219964)
  $core.bool hasDestinationkeyid() => $_has(9);
  @$pb.TagNumber(396219964)
  void clearDestinationkeyid() => $_clearField(396219964);
}

class ReEncryptResponse extends $pb.GeneratedMessage {
  factory ReEncryptResponse({
    EncryptionAlgorithmSpec? destinationencryptionalgorithm,
    $core.String? sourcekeymaterialid,
    $core.String? sourcekeyid,
    $core.String? keyid,
    EncryptionAlgorithmSpec? sourceencryptionalgorithm,
    $core.List<$core.int>? ciphertextblob,
    $core.String? destinationkeymaterialid,
  }) {
    final result = create();
    if (destinationencryptionalgorithm != null)
      result.destinationencryptionalgorithm = destinationencryptionalgorithm;
    if (sourcekeymaterialid != null)
      result.sourcekeymaterialid = sourcekeymaterialid;
    if (sourcekeyid != null) result.sourcekeyid = sourcekeyid;
    if (keyid != null) result.keyid = keyid;
    if (sourceencryptionalgorithm != null)
      result.sourceencryptionalgorithm = sourceencryptionalgorithm;
    if (ciphertextblob != null) result.ciphertextblob = ciphertextblob;
    if (destinationkeymaterialid != null)
      result.destinationkeymaterialid = destinationkeymaterialid;
    return result;
  }

  ReEncryptResponse._();

  factory ReEncryptResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReEncryptResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReEncryptResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<EncryptionAlgorithmSpec>(
        3500944, _omitFieldNames ? '' : 'destinationencryptionalgorithm',
        enumValues: EncryptionAlgorithmSpec.values)
    ..aOS(34789220, _omitFieldNames ? '' : 'sourcekeymaterialid')
    ..aOS(137105771, _omitFieldNames ? '' : 'sourcekeyid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aE<EncryptionAlgorithmSpec>(
        283215847, _omitFieldNames ? '' : 'sourceencryptionalgorithm',
        enumValues: EncryptionAlgorithmSpec.values)
    ..a<$core.List<$core.int>>(
        338198183, _omitFieldNames ? '' : 'ciphertextblob', $pb.PbFieldType.OY)
    ..aOS(469834039, _omitFieldNames ? '' : 'destinationkeymaterialid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReEncryptResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReEncryptResponse copyWith(void Function(ReEncryptResponse) updates) =>
      super.copyWith((message) => updates(message as ReEncryptResponse))
          as ReEncryptResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReEncryptResponse create() => ReEncryptResponse._();
  @$core.override
  ReEncryptResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReEncryptResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReEncryptResponse>(create);
  static ReEncryptResponse? _defaultInstance;

  @$pb.TagNumber(3500944)
  EncryptionAlgorithmSpec get destinationencryptionalgorithm => $_getN(0);
  @$pb.TagNumber(3500944)
  set destinationencryptionalgorithm(EncryptionAlgorithmSpec value) =>
      $_setField(3500944, value);
  @$pb.TagNumber(3500944)
  $core.bool hasDestinationencryptionalgorithm() => $_has(0);
  @$pb.TagNumber(3500944)
  void clearDestinationencryptionalgorithm() => $_clearField(3500944);

  @$pb.TagNumber(34789220)
  $core.String get sourcekeymaterialid => $_getSZ(1);
  @$pb.TagNumber(34789220)
  set sourcekeymaterialid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(34789220)
  $core.bool hasSourcekeymaterialid() => $_has(1);
  @$pb.TagNumber(34789220)
  void clearSourcekeymaterialid() => $_clearField(34789220);

  @$pb.TagNumber(137105771)
  $core.String get sourcekeyid => $_getSZ(2);
  @$pb.TagNumber(137105771)
  set sourcekeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(137105771)
  $core.bool hasSourcekeyid() => $_has(2);
  @$pb.TagNumber(137105771)
  void clearSourcekeyid() => $_clearField(137105771);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(283215847)
  EncryptionAlgorithmSpec get sourceencryptionalgorithm => $_getN(4);
  @$pb.TagNumber(283215847)
  set sourceencryptionalgorithm(EncryptionAlgorithmSpec value) =>
      $_setField(283215847, value);
  @$pb.TagNumber(283215847)
  $core.bool hasSourceencryptionalgorithm() => $_has(4);
  @$pb.TagNumber(283215847)
  void clearSourceencryptionalgorithm() => $_clearField(283215847);

  @$pb.TagNumber(338198183)
  $core.List<$core.int> get ciphertextblob => $_getN(5);
  @$pb.TagNumber(338198183)
  set ciphertextblob($core.List<$core.int> value) => $_setBytes(5, value);
  @$pb.TagNumber(338198183)
  $core.bool hasCiphertextblob() => $_has(5);
  @$pb.TagNumber(338198183)
  void clearCiphertextblob() => $_clearField(338198183);

  @$pb.TagNumber(469834039)
  $core.String get destinationkeymaterialid => $_getSZ(6);
  @$pb.TagNumber(469834039)
  set destinationkeymaterialid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(469834039)
  $core.bool hasDestinationkeymaterialid() => $_has(6);
  @$pb.TagNumber(469834039)
  void clearDestinationkeymaterialid() => $_clearField(469834039);
}

class RecipientInfo extends $pb.GeneratedMessage {
  factory RecipientInfo({
    $core.List<$core.int>? attestationdocument,
    KeyEncryptionMechanism? keyencryptionalgorithm,
  }) {
    final result = create();
    if (attestationdocument != null)
      result.attestationdocument = attestationdocument;
    if (keyencryptionalgorithm != null)
      result.keyencryptionalgorithm = keyencryptionalgorithm;
    return result;
  }

  RecipientInfo._();

  factory RecipientInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RecipientInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RecipientInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(217786849,
        _omitFieldNames ? '' : 'attestationdocument', $pb.PbFieldType.OY)
    ..aE<KeyEncryptionMechanism>(
        478234803, _omitFieldNames ? '' : 'keyencryptionalgorithm',
        enumValues: KeyEncryptionMechanism.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecipientInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RecipientInfo copyWith(void Function(RecipientInfo) updates) =>
      super.copyWith((message) => updates(message as RecipientInfo))
          as RecipientInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RecipientInfo create() => RecipientInfo._();
  @$core.override
  RecipientInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RecipientInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RecipientInfo>(create);
  static RecipientInfo? _defaultInstance;

  @$pb.TagNumber(217786849)
  $core.List<$core.int> get attestationdocument => $_getN(0);
  @$pb.TagNumber(217786849)
  set attestationdocument($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(217786849)
  $core.bool hasAttestationdocument() => $_has(0);
  @$pb.TagNumber(217786849)
  void clearAttestationdocument() => $_clearField(217786849);

  @$pb.TagNumber(478234803)
  KeyEncryptionMechanism get keyencryptionalgorithm => $_getN(1);
  @$pb.TagNumber(478234803)
  set keyencryptionalgorithm(KeyEncryptionMechanism value) =>
      $_setField(478234803, value);
  @$pb.TagNumber(478234803)
  $core.bool hasKeyencryptionalgorithm() => $_has(1);
  @$pb.TagNumber(478234803)
  void clearKeyencryptionalgorithm() => $_clearField(478234803);
}

class ReplicateKeyRequest extends $pb.GeneratedMessage {
  factory ReplicateKeyRequest({
    $core.String? description,
    $core.String? replicaregion,
    $core.bool? bypasspolicylockoutsafetycheck,
    $core.String? keyid,
    $core.Iterable<Tag>? tags,
    $core.String? policy,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (replicaregion != null) result.replicaregion = replicaregion;
    if (bypasspolicylockoutsafetycheck != null)
      result.bypasspolicylockoutsafetycheck = bypasspolicylockoutsafetycheck;
    if (keyid != null) result.keyid = keyid;
    if (tags != null) result.tags.addAll(tags);
    if (policy != null) result.policy = policy;
    return result;
  }

  ReplicateKeyRequest._();

  factory ReplicateKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicateKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicateKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(160061584, _omitFieldNames ? '' : 'replicaregion')
    ..aOB(177450851, _omitFieldNames ? '' : 'bypasspolicylockoutsafetycheck')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicateKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicateKeyRequest copyWith(void Function(ReplicateKeyRequest) updates) =>
      super.copyWith((message) => updates(message as ReplicateKeyRequest))
          as ReplicateKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicateKeyRequest create() => ReplicateKeyRequest._();
  @$core.override
  ReplicateKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicateKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicateKeyRequest>(create);
  static ReplicateKeyRequest? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(160061584)
  $core.String get replicaregion => $_getSZ(1);
  @$pb.TagNumber(160061584)
  set replicaregion($core.String value) => $_setString(1, value);
  @$pb.TagNumber(160061584)
  $core.bool hasReplicaregion() => $_has(1);
  @$pb.TagNumber(160061584)
  void clearReplicaregion() => $_clearField(160061584);

  @$pb.TagNumber(177450851)
  $core.bool get bypasspolicylockoutsafetycheck => $_getBF(2);
  @$pb.TagNumber(177450851)
  set bypasspolicylockoutsafetycheck($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(177450851)
  $core.bool hasBypasspolicylockoutsafetycheck() => $_has(2);
  @$pb.TagNumber(177450851)
  void clearBypasspolicylockoutsafetycheck() => $_clearField(177450851);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(4);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(5);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(5, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(5);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class ReplicateKeyResponse extends $pb.GeneratedMessage {
  factory ReplicateKeyResponse({
    $core.Iterable<Tag>? replicatags,
    KeyMetadata? replicakeymetadata,
    $core.String? replicapolicy,
  }) {
    final result = create();
    if (replicatags != null) result.replicatags.addAll(replicatags);
    if (replicakeymetadata != null)
      result.replicakeymetadata = replicakeymetadata;
    if (replicapolicy != null) result.replicapolicy = replicapolicy;
    return result;
  }

  ReplicateKeyResponse._();

  factory ReplicateKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicateKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicateKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..pPM<Tag>(17934651, _omitFieldNames ? '' : 'replicatags',
        subBuilder: Tag.create)
    ..aOM<KeyMetadata>(68328236, _omitFieldNames ? '' : 'replicakeymetadata',
        subBuilder: KeyMetadata.create)
    ..aOS(279018266, _omitFieldNames ? '' : 'replicapolicy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicateKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicateKeyResponse copyWith(void Function(ReplicateKeyResponse) updates) =>
      super.copyWith((message) => updates(message as ReplicateKeyResponse))
          as ReplicateKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicateKeyResponse create() => ReplicateKeyResponse._();
  @$core.override
  ReplicateKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicateKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicateKeyResponse>(create);
  static ReplicateKeyResponse? _defaultInstance;

  @$pb.TagNumber(17934651)
  $pb.PbList<Tag> get replicatags => $_getList(0);

  @$pb.TagNumber(68328236)
  KeyMetadata get replicakeymetadata => $_getN(1);
  @$pb.TagNumber(68328236)
  set replicakeymetadata(KeyMetadata value) => $_setField(68328236, value);
  @$pb.TagNumber(68328236)
  $core.bool hasReplicakeymetadata() => $_has(1);
  @$pb.TagNumber(68328236)
  void clearReplicakeymetadata() => $_clearField(68328236);
  @$pb.TagNumber(68328236)
  KeyMetadata ensureReplicakeymetadata() => $_ensure(1);

  @$pb.TagNumber(279018266)
  $core.String get replicapolicy => $_getSZ(2);
  @$pb.TagNumber(279018266)
  set replicapolicy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(279018266)
  $core.bool hasReplicapolicy() => $_has(2);
  @$pb.TagNumber(279018266)
  void clearReplicapolicy() => $_clearField(279018266);
}

class RetireGrantRequest extends $pb.GeneratedMessage {
  factory RetireGrantRequest({
    $core.String? grantid,
    $core.bool? dryrun,
    $core.String? granttoken,
    $core.String? keyid,
  }) {
    final result = create();
    if (grantid != null) result.grantid = grantid;
    if (dryrun != null) result.dryrun = dryrun;
    if (granttoken != null) result.granttoken = granttoken;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  RetireGrantRequest._();

  factory RetireGrantRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RetireGrantRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RetireGrantRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(66852281, _omitFieldNames ? '' : 'grantid')
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aOS(137683547, _omitFieldNames ? '' : 'granttoken')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RetireGrantRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RetireGrantRequest copyWith(void Function(RetireGrantRequest) updates) =>
      super.copyWith((message) => updates(message as RetireGrantRequest))
          as RetireGrantRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RetireGrantRequest create() => RetireGrantRequest._();
  @$core.override
  RetireGrantRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RetireGrantRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RetireGrantRequest>(create);
  static RetireGrantRequest? _defaultInstance;

  @$pb.TagNumber(66852281)
  $core.String get grantid => $_getSZ(0);
  @$pb.TagNumber(66852281)
  set grantid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(66852281)
  $core.bool hasGrantid() => $_has(0);
  @$pb.TagNumber(66852281)
  void clearGrantid() => $_clearField(66852281);

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(1);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(1);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(137683547)
  $core.String get granttoken => $_getSZ(2);
  @$pb.TagNumber(137683547)
  set granttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(137683547)
  $core.bool hasGranttoken() => $_has(2);
  @$pb.TagNumber(137683547)
  void clearGranttoken() => $_clearField(137683547);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class RevokeGrantRequest extends $pb.GeneratedMessage {
  factory RevokeGrantRequest({
    $core.String? grantid,
    $core.bool? dryrun,
    $core.String? keyid,
  }) {
    final result = create();
    if (grantid != null) result.grantid = grantid;
    if (dryrun != null) result.dryrun = dryrun;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  RevokeGrantRequest._();

  factory RevokeGrantRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RevokeGrantRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RevokeGrantRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(66852281, _omitFieldNames ? '' : 'grantid')
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RevokeGrantRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RevokeGrantRequest copyWith(void Function(RevokeGrantRequest) updates) =>
      super.copyWith((message) => updates(message as RevokeGrantRequest))
          as RevokeGrantRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RevokeGrantRequest create() => RevokeGrantRequest._();
  @$core.override
  RevokeGrantRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RevokeGrantRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RevokeGrantRequest>(create);
  static RevokeGrantRequest? _defaultInstance;

  @$pb.TagNumber(66852281)
  $core.String get grantid => $_getSZ(0);
  @$pb.TagNumber(66852281)
  set grantid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(66852281)
  $core.bool hasGrantid() => $_has(0);
  @$pb.TagNumber(66852281)
  void clearGrantid() => $_clearField(66852281);

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(1);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(1);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class RotateKeyOnDemandRequest extends $pb.GeneratedMessage {
  factory RotateKeyOnDemandRequest({
    $core.String? keyid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  RotateKeyOnDemandRequest._();

  factory RotateKeyOnDemandRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RotateKeyOnDemandRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RotateKeyOnDemandRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotateKeyOnDemandRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotateKeyOnDemandRequest copyWith(
          void Function(RotateKeyOnDemandRequest) updates) =>
      super.copyWith((message) => updates(message as RotateKeyOnDemandRequest))
          as RotateKeyOnDemandRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RotateKeyOnDemandRequest create() => RotateKeyOnDemandRequest._();
  @$core.override
  RotateKeyOnDemandRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RotateKeyOnDemandRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RotateKeyOnDemandRequest>(create);
  static RotateKeyOnDemandRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class RotateKeyOnDemandResponse extends $pb.GeneratedMessage {
  factory RotateKeyOnDemandResponse({
    $core.String? keyid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  RotateKeyOnDemandResponse._();

  factory RotateKeyOnDemandResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RotateKeyOnDemandResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RotateKeyOnDemandResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotateKeyOnDemandResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotateKeyOnDemandResponse copyWith(
          void Function(RotateKeyOnDemandResponse) updates) =>
      super.copyWith((message) => updates(message as RotateKeyOnDemandResponse))
          as RotateKeyOnDemandResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RotateKeyOnDemandResponse create() => RotateKeyOnDemandResponse._();
  @$core.override
  RotateKeyOnDemandResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RotateKeyOnDemandResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RotateKeyOnDemandResponse>(create);
  static RotateKeyOnDemandResponse? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class RotationsListEntry extends $pb.GeneratedMessage {
  factory RotationsListEntry({
    ImportState? importstate,
    ExpirationModelType? expirationmodel,
    RotationType? rotationtype,
    $core.String? keymaterialid,
    $core.String? keyid,
    $core.String? keymaterialdescription,
    KeyMaterialState? keymaterialstate,
    $core.String? validto,
    $core.String? rotationdate,
  }) {
    final result = create();
    if (importstate != null) result.importstate = importstate;
    if (expirationmodel != null) result.expirationmodel = expirationmodel;
    if (rotationtype != null) result.rotationtype = rotationtype;
    if (keymaterialid != null) result.keymaterialid = keymaterialid;
    if (keyid != null) result.keyid = keyid;
    if (keymaterialdescription != null)
      result.keymaterialdescription = keymaterialdescription;
    if (keymaterialstate != null) result.keymaterialstate = keymaterialstate;
    if (validto != null) result.validto = validto;
    if (rotationdate != null) result.rotationdate = rotationdate;
    return result;
  }

  RotationsListEntry._();

  factory RotationsListEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RotationsListEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RotationsListEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<ImportState>(32548970, _omitFieldNames ? '' : 'importstate',
        enumValues: ImportState.values)
    ..aE<ExpirationModelType>(
        113933558, _omitFieldNames ? '' : 'expirationmodel',
        enumValues: ExpirationModelType.values)
    ..aE<RotationType>(122951592, _omitFieldNames ? '' : 'rotationtype',
        enumValues: RotationType.values)
    ..aOS(147011585, _omitFieldNames ? '' : 'keymaterialid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(277153546, _omitFieldNames ? '' : 'keymaterialdescription')
    ..aE<KeyMaterialState>(431806871, _omitFieldNames ? '' : 'keymaterialstate',
        enumValues: KeyMaterialState.values)
    ..aOS(522718673, _omitFieldNames ? '' : 'validto')
    ..aOS(529238652, _omitFieldNames ? '' : 'rotationdate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotationsListEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotationsListEntry copyWith(void Function(RotationsListEntry) updates) =>
      super.copyWith((message) => updates(message as RotationsListEntry))
          as RotationsListEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RotationsListEntry create() => RotationsListEntry._();
  @$core.override
  RotationsListEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RotationsListEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RotationsListEntry>(create);
  static RotationsListEntry? _defaultInstance;

  @$pb.TagNumber(32548970)
  ImportState get importstate => $_getN(0);
  @$pb.TagNumber(32548970)
  set importstate(ImportState value) => $_setField(32548970, value);
  @$pb.TagNumber(32548970)
  $core.bool hasImportstate() => $_has(0);
  @$pb.TagNumber(32548970)
  void clearImportstate() => $_clearField(32548970);

  @$pb.TagNumber(113933558)
  ExpirationModelType get expirationmodel => $_getN(1);
  @$pb.TagNumber(113933558)
  set expirationmodel(ExpirationModelType value) =>
      $_setField(113933558, value);
  @$pb.TagNumber(113933558)
  $core.bool hasExpirationmodel() => $_has(1);
  @$pb.TagNumber(113933558)
  void clearExpirationmodel() => $_clearField(113933558);

  @$pb.TagNumber(122951592)
  RotationType get rotationtype => $_getN(2);
  @$pb.TagNumber(122951592)
  set rotationtype(RotationType value) => $_setField(122951592, value);
  @$pb.TagNumber(122951592)
  $core.bool hasRotationtype() => $_has(2);
  @$pb.TagNumber(122951592)
  void clearRotationtype() => $_clearField(122951592);

  @$pb.TagNumber(147011585)
  $core.String get keymaterialid => $_getSZ(3);
  @$pb.TagNumber(147011585)
  set keymaterialid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(147011585)
  $core.bool hasKeymaterialid() => $_has(3);
  @$pb.TagNumber(147011585)
  void clearKeymaterialid() => $_clearField(147011585);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(4);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(4);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(277153546)
  $core.String get keymaterialdescription => $_getSZ(5);
  @$pb.TagNumber(277153546)
  set keymaterialdescription($core.String value) => $_setString(5, value);
  @$pb.TagNumber(277153546)
  $core.bool hasKeymaterialdescription() => $_has(5);
  @$pb.TagNumber(277153546)
  void clearKeymaterialdescription() => $_clearField(277153546);

  @$pb.TagNumber(431806871)
  KeyMaterialState get keymaterialstate => $_getN(6);
  @$pb.TagNumber(431806871)
  set keymaterialstate(KeyMaterialState value) => $_setField(431806871, value);
  @$pb.TagNumber(431806871)
  $core.bool hasKeymaterialstate() => $_has(6);
  @$pb.TagNumber(431806871)
  void clearKeymaterialstate() => $_clearField(431806871);

  @$pb.TagNumber(522718673)
  $core.String get validto => $_getSZ(7);
  @$pb.TagNumber(522718673)
  set validto($core.String value) => $_setString(7, value);
  @$pb.TagNumber(522718673)
  $core.bool hasValidto() => $_has(7);
  @$pb.TagNumber(522718673)
  void clearValidto() => $_clearField(522718673);

  @$pb.TagNumber(529238652)
  $core.String get rotationdate => $_getSZ(8);
  @$pb.TagNumber(529238652)
  set rotationdate($core.String value) => $_setString(8, value);
  @$pb.TagNumber(529238652)
  $core.bool hasRotationdate() => $_has(8);
  @$pb.TagNumber(529238652)
  void clearRotationdate() => $_clearField(529238652);
}

class ScheduleKeyDeletionRequest extends $pb.GeneratedMessage {
  factory ScheduleKeyDeletionRequest({
    $core.String? keyid,
    $core.int? pendingwindowindays,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    if (pendingwindowindays != null)
      result.pendingwindowindays = pendingwindowindays;
    return result;
  }

  ScheduleKeyDeletionRequest._();

  factory ScheduleKeyDeletionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduleKeyDeletionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduleKeyDeletionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aI(532945081, _omitFieldNames ? '' : 'pendingwindowindays')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleKeyDeletionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleKeyDeletionRequest copyWith(
          void Function(ScheduleKeyDeletionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ScheduleKeyDeletionRequest))
          as ScheduleKeyDeletionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduleKeyDeletionRequest create() => ScheduleKeyDeletionRequest._();
  @$core.override
  ScheduleKeyDeletionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduleKeyDeletionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduleKeyDeletionRequest>(create);
  static ScheduleKeyDeletionRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(532945081)
  $core.int get pendingwindowindays => $_getIZ(1);
  @$pb.TagNumber(532945081)
  set pendingwindowindays($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(532945081)
  $core.bool hasPendingwindowindays() => $_has(1);
  @$pb.TagNumber(532945081)
  void clearPendingwindowindays() => $_clearField(532945081);
}

class ScheduleKeyDeletionResponse extends $pb.GeneratedMessage {
  factory ScheduleKeyDeletionResponse({
    KeyState? keystate,
    $core.String? keyid,
    $core.String? deletiondate,
    $core.int? pendingwindowindays,
  }) {
    final result = create();
    if (keystate != null) result.keystate = keystate;
    if (keyid != null) result.keyid = keyid;
    if (deletiondate != null) result.deletiondate = deletiondate;
    if (pendingwindowindays != null)
      result.pendingwindowindays = pendingwindowindays;
    return result;
  }

  ScheduleKeyDeletionResponse._();

  factory ScheduleKeyDeletionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduleKeyDeletionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduleKeyDeletionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<KeyState>(27894226, _omitFieldNames ? '' : 'keystate',
        enumValues: KeyState.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(347845564, _omitFieldNames ? '' : 'deletiondate')
    ..aI(532945081, _omitFieldNames ? '' : 'pendingwindowindays')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleKeyDeletionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleKeyDeletionResponse copyWith(
          void Function(ScheduleKeyDeletionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ScheduleKeyDeletionResponse))
          as ScheduleKeyDeletionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduleKeyDeletionResponse create() =>
      ScheduleKeyDeletionResponse._();
  @$core.override
  ScheduleKeyDeletionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduleKeyDeletionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduleKeyDeletionResponse>(create);
  static ScheduleKeyDeletionResponse? _defaultInstance;

  @$pb.TagNumber(27894226)
  KeyState get keystate => $_getN(0);
  @$pb.TagNumber(27894226)
  set keystate(KeyState value) => $_setField(27894226, value);
  @$pb.TagNumber(27894226)
  $core.bool hasKeystate() => $_has(0);
  @$pb.TagNumber(27894226)
  void clearKeystate() => $_clearField(27894226);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(347845564)
  $core.String get deletiondate => $_getSZ(2);
  @$pb.TagNumber(347845564)
  set deletiondate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(347845564)
  $core.bool hasDeletiondate() => $_has(2);
  @$pb.TagNumber(347845564)
  void clearDeletiondate() => $_clearField(347845564);

  @$pb.TagNumber(532945081)
  $core.int get pendingwindowindays => $_getIZ(3);
  @$pb.TagNumber(532945081)
  set pendingwindowindays($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(532945081)
  $core.bool hasPendingwindowindays() => $_has(3);
  @$pb.TagNumber(532945081)
  void clearPendingwindowindays() => $_clearField(532945081);
}

class SignRequest extends $pb.GeneratedMessage {
  factory SignRequest({
    $core.bool? dryrun,
    $core.List<$core.int>? message,
    $core.String? keyid,
    $core.Iterable<$core.String>? granttokens,
    MessageType? messagetype,
    SigningAlgorithmSpec? signingalgorithm,
  }) {
    final result = create();
    if (dryrun != null) result.dryrun = dryrun;
    if (message != null) result.message = message;
    if (keyid != null) result.keyid = keyid;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    if (messagetype != null) result.messagetype = messagetype;
    if (signingalgorithm != null) result.signingalgorithm = signingalgorithm;
    return result;
  }

  SignRequest._();

  factory SignRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SignRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SignRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..a<$core.List<$core.int>>(
        235854213, _omitFieldNames ? '' : 'message', $pb.PbFieldType.OY)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..aE<MessageType>(452676013, _omitFieldNames ? '' : 'messagetype',
        enumValues: MessageType.values)
    ..aE<SigningAlgorithmSpec>(
        488091842, _omitFieldNames ? '' : 'signingalgorithm',
        enumValues: SigningAlgorithmSpec.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignRequest copyWith(void Function(SignRequest) updates) =>
      super.copyWith((message) => updates(message as SignRequest))
          as SignRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignRequest create() => SignRequest._();
  @$core.override
  SignRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SignRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SignRequest>(create);
  static SignRequest? _defaultInstance;

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(0);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(0);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(235854213)
  $core.List<$core.int> get message => $_getN(1);
  @$pb.TagNumber(235854213)
  set message($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(3);

  @$pb.TagNumber(452676013)
  MessageType get messagetype => $_getN(4);
  @$pb.TagNumber(452676013)
  set messagetype(MessageType value) => $_setField(452676013, value);
  @$pb.TagNumber(452676013)
  $core.bool hasMessagetype() => $_has(4);
  @$pb.TagNumber(452676013)
  void clearMessagetype() => $_clearField(452676013);

  @$pb.TagNumber(488091842)
  SigningAlgorithmSpec get signingalgorithm => $_getN(5);
  @$pb.TagNumber(488091842)
  set signingalgorithm(SigningAlgorithmSpec value) =>
      $_setField(488091842, value);
  @$pb.TagNumber(488091842)
  $core.bool hasSigningalgorithm() => $_has(5);
  @$pb.TagNumber(488091842)
  void clearSigningalgorithm() => $_clearField(488091842);
}

class SignResponse extends $pb.GeneratedMessage {
  factory SignResponse({
    $core.List<$core.int>? signature,
    $core.String? keyid,
    SigningAlgorithmSpec? signingalgorithm,
  }) {
    final result = create();
    if (signature != null) result.signature = signature;
    if (keyid != null) result.keyid = keyid;
    if (signingalgorithm != null) result.signingalgorithm = signingalgorithm;
    return result;
  }

  SignResponse._();

  factory SignResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SignResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SignResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(
        4785422, _omitFieldNames ? '' : 'signature', $pb.PbFieldType.OY)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aE<SigningAlgorithmSpec>(
        488091842, _omitFieldNames ? '' : 'signingalgorithm',
        enumValues: SigningAlgorithmSpec.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignResponse copyWith(void Function(SignResponse) updates) =>
      super.copyWith((message) => updates(message as SignResponse))
          as SignResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignResponse create() => SignResponse._();
  @$core.override
  SignResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SignResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SignResponse>(create);
  static SignResponse? _defaultInstance;

  @$pb.TagNumber(4785422)
  $core.List<$core.int> get signature => $_getN(0);
  @$pb.TagNumber(4785422)
  set signature($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(4785422)
  $core.bool hasSignature() => $_has(0);
  @$pb.TagNumber(4785422)
  void clearSignature() => $_clearField(4785422);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(488091842)
  SigningAlgorithmSpec get signingalgorithm => $_getN(2);
  @$pb.TagNumber(488091842)
  set signingalgorithm(SigningAlgorithmSpec value) =>
      $_setField(488091842, value);
  @$pb.TagNumber(488091842)
  $core.bool hasSigningalgorithm() => $_has(2);
  @$pb.TagNumber(488091842)
  void clearSigningalgorithm() => $_clearField(488091842);
}

class Tag extends $pb.GeneratedMessage {
  factory Tag({
    $core.String? tagkey,
    $core.String? tagvalue,
  }) {
    final result = create();
    if (tagkey != null) result.tagkey = tagkey;
    if (tagvalue != null) result.tagvalue = tagvalue;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(175603339, _omitFieldNames ? '' : 'tagkey')
    ..aOS(487739505, _omitFieldNames ? '' : 'tagvalue')
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

  @$pb.TagNumber(175603339)
  $core.String get tagkey => $_getSZ(0);
  @$pb.TagNumber(175603339)
  set tagkey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(175603339)
  $core.bool hasTagkey() => $_has(0);
  @$pb.TagNumber(175603339)
  void clearTagkey() => $_clearField(175603339);

  @$pb.TagNumber(487739505)
  $core.String get tagvalue => $_getSZ(1);
  @$pb.TagNumber(487739505)
  set tagvalue($core.String value) => $_setString(1, value);
  @$pb.TagNumber(487739505)
  $core.bool hasTagvalue() => $_has(1);
  @$pb.TagNumber(487739505)
  void clearTagvalue() => $_clearField(487739505);
}

class TagException extends $pb.GeneratedMessage {
  factory TagException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TagException._();

  factory TagException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagException copyWith(void Function(TagException) updates) =>
      super.copyWith((message) => updates(message as TagException))
          as TagException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagException create() => TagException._();
  @$core.override
  TagException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagException>(create);
  static TagException? _defaultInstance;

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
    $core.String? keyid,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
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

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);
}

class UnsupportedOperationException extends $pb.GeneratedMessage {
  factory UnsupportedOperationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  UnsupportedOperationException._();

  factory UnsupportedOperationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnsupportedOperationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnsupportedOperationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnsupportedOperationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnsupportedOperationException copyWith(
          void Function(UnsupportedOperationException) updates) =>
      super.copyWith(
              (message) => updates(message as UnsupportedOperationException))
          as UnsupportedOperationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnsupportedOperationException create() =>
      UnsupportedOperationException._();
  @$core.override
  UnsupportedOperationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnsupportedOperationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnsupportedOperationException>(create);
  static UnsupportedOperationException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class UntagResourceRequest extends $pb.GeneratedMessage {
  factory UntagResourceRequest({
    $core.String? keyid,
    $core.Iterable<$core.String>? tagkeys,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
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

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(320659964)
  $pb.PbList<$core.String> get tagkeys => $_getList(1);
}

class UpdateAliasRequest extends $pb.GeneratedMessage {
  factory UpdateAliasRequest({
    $core.String? aliasname,
    $core.String? targetkeyid,
  }) {
    final result = create();
    if (aliasname != null) result.aliasname = aliasname;
    if (targetkeyid != null) result.targetkeyid = targetkeyid;
    return result;
  }

  UpdateAliasRequest._();

  factory UpdateAliasRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateAliasRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAliasRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(313250709, _omitFieldNames ? '' : 'aliasname')
    ..aOS(406196123, _omitFieldNames ? '' : 'targetkeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAliasRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAliasRequest copyWith(void Function(UpdateAliasRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateAliasRequest))
          as UpdateAliasRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAliasRequest create() => UpdateAliasRequest._();
  @$core.override
  UpdateAliasRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateAliasRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateAliasRequest>(create);
  static UpdateAliasRequest? _defaultInstance;

  @$pb.TagNumber(313250709)
  $core.String get aliasname => $_getSZ(0);
  @$pb.TagNumber(313250709)
  set aliasname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(313250709)
  $core.bool hasAliasname() => $_has(0);
  @$pb.TagNumber(313250709)
  void clearAliasname() => $_clearField(313250709);

  @$pb.TagNumber(406196123)
  $core.String get targetkeyid => $_getSZ(1);
  @$pb.TagNumber(406196123)
  set targetkeyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(406196123)
  $core.bool hasTargetkeyid() => $_has(1);
  @$pb.TagNumber(406196123)
  void clearTargetkeyid() => $_clearField(406196123);
}

class UpdateCustomKeyStoreRequest extends $pb.GeneratedMessage {
  factory UpdateCustomKeyStoreRequest({
    $core.String? newcustomkeystorename,
    $core.String? xksproxyvpcendpointserviceowner,
    $core.String? cloudhsmclusterid,
    $core.String? customkeystoreid,
    $core.String? xksproxyuriendpoint,
    XksProxyConnectivityType? xksproxyconnectivity,
    XksProxyAuthenticationCredentialType? xksproxyauthenticationcredential,
    $core.String? xksproxyvpcendpointservicename,
    $core.String? keystorepassword,
    $core.String? xksproxyuripath,
  }) {
    final result = create();
    if (newcustomkeystorename != null)
      result.newcustomkeystorename = newcustomkeystorename;
    if (xksproxyvpcendpointserviceowner != null)
      result.xksproxyvpcendpointserviceowner = xksproxyvpcendpointserviceowner;
    if (cloudhsmclusterid != null) result.cloudhsmclusterid = cloudhsmclusterid;
    if (customkeystoreid != null) result.customkeystoreid = customkeystoreid;
    if (xksproxyuriendpoint != null)
      result.xksproxyuriendpoint = xksproxyuriendpoint;
    if (xksproxyconnectivity != null)
      result.xksproxyconnectivity = xksproxyconnectivity;
    if (xksproxyauthenticationcredential != null)
      result.xksproxyauthenticationcredential =
          xksproxyauthenticationcredential;
    if (xksproxyvpcendpointservicename != null)
      result.xksproxyvpcendpointservicename = xksproxyvpcendpointservicename;
    if (keystorepassword != null) result.keystorepassword = keystorepassword;
    if (xksproxyuripath != null) result.xksproxyuripath = xksproxyuripath;
    return result;
  }

  UpdateCustomKeyStoreRequest._();

  factory UpdateCustomKeyStoreRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateCustomKeyStoreRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateCustomKeyStoreRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(10936866, _omitFieldNames ? '' : 'newcustomkeystorename')
    ..aOS(55249590, _omitFieldNames ? '' : 'xksproxyvpcendpointserviceowner')
    ..aOS(56498754, _omitFieldNames ? '' : 'cloudhsmclusterid')
    ..aOS(88348228, _omitFieldNames ? '' : 'customkeystoreid')
    ..aOS(273255559, _omitFieldNames ? '' : 'xksproxyuriendpoint')
    ..aE<XksProxyConnectivityType>(
        298569161, _omitFieldNames ? '' : 'xksproxyconnectivity',
        enumValues: XksProxyConnectivityType.values)
    ..aOM<XksProxyAuthenticationCredentialType>(
        350418199, _omitFieldNames ? '' : 'xksproxyauthenticationcredential',
        subBuilder: XksProxyAuthenticationCredentialType.create)
    ..aOS(372786130, _omitFieldNames ? '' : 'xksproxyvpcendpointservicename')
    ..aOS(403136353, _omitFieldNames ? '' : 'keystorepassword')
    ..aOS(436753509, _omitFieldNames ? '' : 'xksproxyuripath')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCustomKeyStoreRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCustomKeyStoreRequest copyWith(
          void Function(UpdateCustomKeyStoreRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateCustomKeyStoreRequest))
          as UpdateCustomKeyStoreRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateCustomKeyStoreRequest create() =>
      UpdateCustomKeyStoreRequest._();
  @$core.override
  UpdateCustomKeyStoreRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateCustomKeyStoreRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateCustomKeyStoreRequest>(create);
  static UpdateCustomKeyStoreRequest? _defaultInstance;

  @$pb.TagNumber(10936866)
  $core.String get newcustomkeystorename => $_getSZ(0);
  @$pb.TagNumber(10936866)
  set newcustomkeystorename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(10936866)
  $core.bool hasNewcustomkeystorename() => $_has(0);
  @$pb.TagNumber(10936866)
  void clearNewcustomkeystorename() => $_clearField(10936866);

  @$pb.TagNumber(55249590)
  $core.String get xksproxyvpcendpointserviceowner => $_getSZ(1);
  @$pb.TagNumber(55249590)
  set xksproxyvpcendpointserviceowner($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(55249590)
  $core.bool hasXksproxyvpcendpointserviceowner() => $_has(1);
  @$pb.TagNumber(55249590)
  void clearXksproxyvpcendpointserviceowner() => $_clearField(55249590);

  @$pb.TagNumber(56498754)
  $core.String get cloudhsmclusterid => $_getSZ(2);
  @$pb.TagNumber(56498754)
  set cloudhsmclusterid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(56498754)
  $core.bool hasCloudhsmclusterid() => $_has(2);
  @$pb.TagNumber(56498754)
  void clearCloudhsmclusterid() => $_clearField(56498754);

  @$pb.TagNumber(88348228)
  $core.String get customkeystoreid => $_getSZ(3);
  @$pb.TagNumber(88348228)
  set customkeystoreid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(88348228)
  $core.bool hasCustomkeystoreid() => $_has(3);
  @$pb.TagNumber(88348228)
  void clearCustomkeystoreid() => $_clearField(88348228);

  @$pb.TagNumber(273255559)
  $core.String get xksproxyuriendpoint => $_getSZ(4);
  @$pb.TagNumber(273255559)
  set xksproxyuriendpoint($core.String value) => $_setString(4, value);
  @$pb.TagNumber(273255559)
  $core.bool hasXksproxyuriendpoint() => $_has(4);
  @$pb.TagNumber(273255559)
  void clearXksproxyuriendpoint() => $_clearField(273255559);

  @$pb.TagNumber(298569161)
  XksProxyConnectivityType get xksproxyconnectivity => $_getN(5);
  @$pb.TagNumber(298569161)
  set xksproxyconnectivity(XksProxyConnectivityType value) =>
      $_setField(298569161, value);
  @$pb.TagNumber(298569161)
  $core.bool hasXksproxyconnectivity() => $_has(5);
  @$pb.TagNumber(298569161)
  void clearXksproxyconnectivity() => $_clearField(298569161);

  @$pb.TagNumber(350418199)
  XksProxyAuthenticationCredentialType get xksproxyauthenticationcredential =>
      $_getN(6);
  @$pb.TagNumber(350418199)
  set xksproxyauthenticationcredential(
          XksProxyAuthenticationCredentialType value) =>
      $_setField(350418199, value);
  @$pb.TagNumber(350418199)
  $core.bool hasXksproxyauthenticationcredential() => $_has(6);
  @$pb.TagNumber(350418199)
  void clearXksproxyauthenticationcredential() => $_clearField(350418199);
  @$pb.TagNumber(350418199)
  XksProxyAuthenticationCredentialType
      ensureXksproxyauthenticationcredential() => $_ensure(6);

  @$pb.TagNumber(372786130)
  $core.String get xksproxyvpcendpointservicename => $_getSZ(7);
  @$pb.TagNumber(372786130)
  set xksproxyvpcendpointservicename($core.String value) =>
      $_setString(7, value);
  @$pb.TagNumber(372786130)
  $core.bool hasXksproxyvpcendpointservicename() => $_has(7);
  @$pb.TagNumber(372786130)
  void clearXksproxyvpcendpointservicename() => $_clearField(372786130);

  @$pb.TagNumber(403136353)
  $core.String get keystorepassword => $_getSZ(8);
  @$pb.TagNumber(403136353)
  set keystorepassword($core.String value) => $_setString(8, value);
  @$pb.TagNumber(403136353)
  $core.bool hasKeystorepassword() => $_has(8);
  @$pb.TagNumber(403136353)
  void clearKeystorepassword() => $_clearField(403136353);

  @$pb.TagNumber(436753509)
  $core.String get xksproxyuripath => $_getSZ(9);
  @$pb.TagNumber(436753509)
  set xksproxyuripath($core.String value) => $_setString(9, value);
  @$pb.TagNumber(436753509)
  $core.bool hasXksproxyuripath() => $_has(9);
  @$pb.TagNumber(436753509)
  void clearXksproxyuripath() => $_clearField(436753509);
}

class UpdateCustomKeyStoreResponse extends $pb.GeneratedMessage {
  factory UpdateCustomKeyStoreResponse() => create();

  UpdateCustomKeyStoreResponse._();

  factory UpdateCustomKeyStoreResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateCustomKeyStoreResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateCustomKeyStoreResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCustomKeyStoreResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateCustomKeyStoreResponse copyWith(
          void Function(UpdateCustomKeyStoreResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateCustomKeyStoreResponse))
          as UpdateCustomKeyStoreResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateCustomKeyStoreResponse create() =>
      UpdateCustomKeyStoreResponse._();
  @$core.override
  UpdateCustomKeyStoreResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateCustomKeyStoreResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateCustomKeyStoreResponse>(create);
  static UpdateCustomKeyStoreResponse? _defaultInstance;
}

class UpdateKeyDescriptionRequest extends $pb.GeneratedMessage {
  factory UpdateKeyDescriptionRequest({
    $core.String? description,
    $core.String? keyid,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  UpdateKeyDescriptionRequest._();

  factory UpdateKeyDescriptionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateKeyDescriptionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateKeyDescriptionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateKeyDescriptionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateKeyDescriptionRequest copyWith(
          void Function(UpdateKeyDescriptionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateKeyDescriptionRequest))
          as UpdateKeyDescriptionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateKeyDescriptionRequest create() =>
      UpdateKeyDescriptionRequest._();
  @$core.override
  UpdateKeyDescriptionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateKeyDescriptionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateKeyDescriptionRequest>(create);
  static UpdateKeyDescriptionRequest? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);
}

class UpdatePrimaryRegionRequest extends $pb.GeneratedMessage {
  factory UpdatePrimaryRegionRequest({
    $core.String? keyid,
    $core.String? primaryregion,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    if (primaryregion != null) result.primaryregion = primaryregion;
    return result;
  }

  UpdatePrimaryRegionRequest._();

  factory UpdatePrimaryRegionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdatePrimaryRegionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdatePrimaryRegionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOS(480901186, _omitFieldNames ? '' : 'primaryregion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdatePrimaryRegionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdatePrimaryRegionRequest copyWith(
          void Function(UpdatePrimaryRegionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdatePrimaryRegionRequest))
          as UpdatePrimaryRegionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdatePrimaryRegionRequest create() => UpdatePrimaryRegionRequest._();
  @$core.override
  UpdatePrimaryRegionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdatePrimaryRegionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdatePrimaryRegionRequest>(create);
  static UpdatePrimaryRegionRequest? _defaultInstance;

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(480901186)
  $core.String get primaryregion => $_getSZ(1);
  @$pb.TagNumber(480901186)
  set primaryregion($core.String value) => $_setString(1, value);
  @$pb.TagNumber(480901186)
  $core.bool hasPrimaryregion() => $_has(1);
  @$pb.TagNumber(480901186)
  void clearPrimaryregion() => $_clearField(480901186);
}

class VerifyMacRequest extends $pb.GeneratedMessage {
  factory VerifyMacRequest({
    $core.bool? dryrun,
    $core.List<$core.int>? message,
    MacAlgorithmSpec? macalgorithm,
    $core.String? keyid,
    $core.List<$core.int>? mac,
    $core.Iterable<$core.String>? granttokens,
  }) {
    final result = create();
    if (dryrun != null) result.dryrun = dryrun;
    if (message != null) result.message = message;
    if (macalgorithm != null) result.macalgorithm = macalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (mac != null) result.mac = mac;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    return result;
  }

  VerifyMacRequest._();

  factory VerifyMacRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VerifyMacRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VerifyMacRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..a<$core.List<$core.int>>(
        235854213, _omitFieldNames ? '' : 'message', $pb.PbFieldType.OY)
    ..aE<MacAlgorithmSpec>(253519878, _omitFieldNames ? '' : 'macalgorithm',
        enumValues: MacAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..a<$core.List<$core.int>>(
        296223945, _omitFieldNames ? '' : 'mac', $pb.PbFieldType.OY)
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifyMacRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifyMacRequest copyWith(void Function(VerifyMacRequest) updates) =>
      super.copyWith((message) => updates(message as VerifyMacRequest))
          as VerifyMacRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerifyMacRequest create() => VerifyMacRequest._();
  @$core.override
  VerifyMacRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VerifyMacRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VerifyMacRequest>(create);
  static VerifyMacRequest? _defaultInstance;

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(0);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(0);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(235854213)
  $core.List<$core.int> get message => $_getN(1);
  @$pb.TagNumber(235854213)
  set message($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(253519878)
  MacAlgorithmSpec get macalgorithm => $_getN(2);
  @$pb.TagNumber(253519878)
  set macalgorithm(MacAlgorithmSpec value) => $_setField(253519878, value);
  @$pb.TagNumber(253519878)
  $core.bool hasMacalgorithm() => $_has(2);
  @$pb.TagNumber(253519878)
  void clearMacalgorithm() => $_clearField(253519878);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(296223945)
  $core.List<$core.int> get mac => $_getN(4);
  @$pb.TagNumber(296223945)
  set mac($core.List<$core.int> value) => $_setBytes(4, value);
  @$pb.TagNumber(296223945)
  $core.bool hasMac() => $_has(4);
  @$pb.TagNumber(296223945)
  void clearMac() => $_clearField(296223945);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(5);
}

class VerifyMacResponse extends $pb.GeneratedMessage {
  factory VerifyMacResponse({
    MacAlgorithmSpec? macalgorithm,
    $core.String? keyid,
    $core.bool? macvalid,
  }) {
    final result = create();
    if (macalgorithm != null) result.macalgorithm = macalgorithm;
    if (keyid != null) result.keyid = keyid;
    if (macvalid != null) result.macvalid = macvalid;
    return result;
  }

  VerifyMacResponse._();

  factory VerifyMacResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VerifyMacResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VerifyMacResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aE<MacAlgorithmSpec>(253519878, _omitFieldNames ? '' : 'macalgorithm',
        enumValues: MacAlgorithmSpec.values)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aOB(482746075, _omitFieldNames ? '' : 'macvalid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifyMacResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifyMacResponse copyWith(void Function(VerifyMacResponse) updates) =>
      super.copyWith((message) => updates(message as VerifyMacResponse))
          as VerifyMacResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerifyMacResponse create() => VerifyMacResponse._();
  @$core.override
  VerifyMacResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VerifyMacResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VerifyMacResponse>(create);
  static VerifyMacResponse? _defaultInstance;

  @$pb.TagNumber(253519878)
  MacAlgorithmSpec get macalgorithm => $_getN(0);
  @$pb.TagNumber(253519878)
  set macalgorithm(MacAlgorithmSpec value) => $_setField(253519878, value);
  @$pb.TagNumber(253519878)
  $core.bool hasMacalgorithm() => $_has(0);
  @$pb.TagNumber(253519878)
  void clearMacalgorithm() => $_clearField(253519878);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(482746075)
  $core.bool get macvalid => $_getBF(2);
  @$pb.TagNumber(482746075)
  set macvalid($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(482746075)
  $core.bool hasMacvalid() => $_has(2);
  @$pb.TagNumber(482746075)
  void clearMacvalid() => $_clearField(482746075);
}

class VerifyRequest extends $pb.GeneratedMessage {
  factory VerifyRequest({
    $core.List<$core.int>? signature,
    $core.bool? dryrun,
    $core.List<$core.int>? message,
    $core.String? keyid,
    $core.Iterable<$core.String>? granttokens,
    MessageType? messagetype,
    SigningAlgorithmSpec? signingalgorithm,
  }) {
    final result = create();
    if (signature != null) result.signature = signature;
    if (dryrun != null) result.dryrun = dryrun;
    if (message != null) result.message = message;
    if (keyid != null) result.keyid = keyid;
    if (granttokens != null) result.granttokens.addAll(granttokens);
    if (messagetype != null) result.messagetype = messagetype;
    if (signingalgorithm != null) result.signingalgorithm = signingalgorithm;
    return result;
  }

  VerifyRequest._();

  factory VerifyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VerifyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VerifyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(
        4785422, _omitFieldNames ? '' : 'signature', $pb.PbFieldType.OY)
    ..aOB(92204984, _omitFieldNames ? '' : 'dryrun')
    ..a<$core.List<$core.int>>(
        235854213, _omitFieldNames ? '' : 'message', $pb.PbFieldType.OY)
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..pPS(339740300, _omitFieldNames ? '' : 'granttokens')
    ..aE<MessageType>(452676013, _omitFieldNames ? '' : 'messagetype',
        enumValues: MessageType.values)
    ..aE<SigningAlgorithmSpec>(
        488091842, _omitFieldNames ? '' : 'signingalgorithm',
        enumValues: SigningAlgorithmSpec.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifyRequest copyWith(void Function(VerifyRequest) updates) =>
      super.copyWith((message) => updates(message as VerifyRequest))
          as VerifyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerifyRequest create() => VerifyRequest._();
  @$core.override
  VerifyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VerifyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VerifyRequest>(create);
  static VerifyRequest? _defaultInstance;

  @$pb.TagNumber(4785422)
  $core.List<$core.int> get signature => $_getN(0);
  @$pb.TagNumber(4785422)
  set signature($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(4785422)
  $core.bool hasSignature() => $_has(0);
  @$pb.TagNumber(4785422)
  void clearSignature() => $_clearField(4785422);

  @$pb.TagNumber(92204984)
  $core.bool get dryrun => $_getBF(1);
  @$pb.TagNumber(92204984)
  set dryrun($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(92204984)
  $core.bool hasDryrun() => $_has(1);
  @$pb.TagNumber(92204984)
  void clearDryrun() => $_clearField(92204984);

  @$pb.TagNumber(235854213)
  $core.List<$core.int> get message => $_getN(2);
  @$pb.TagNumber(235854213)
  set message($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(2);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(3);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(3);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(339740300)
  $pb.PbList<$core.String> get granttokens => $_getList(4);

  @$pb.TagNumber(452676013)
  MessageType get messagetype => $_getN(5);
  @$pb.TagNumber(452676013)
  set messagetype(MessageType value) => $_setField(452676013, value);
  @$pb.TagNumber(452676013)
  $core.bool hasMessagetype() => $_has(5);
  @$pb.TagNumber(452676013)
  void clearMessagetype() => $_clearField(452676013);

  @$pb.TagNumber(488091842)
  SigningAlgorithmSpec get signingalgorithm => $_getN(6);
  @$pb.TagNumber(488091842)
  set signingalgorithm(SigningAlgorithmSpec value) =>
      $_setField(488091842, value);
  @$pb.TagNumber(488091842)
  $core.bool hasSigningalgorithm() => $_has(6);
  @$pb.TagNumber(488091842)
  void clearSigningalgorithm() => $_clearField(488091842);
}

class VerifyResponse extends $pb.GeneratedMessage {
  factory VerifyResponse({
    $core.bool? signaturevalid,
    $core.String? keyid,
    SigningAlgorithmSpec? signingalgorithm,
  }) {
    final result = create();
    if (signaturevalid != null) result.signaturevalid = signaturevalid;
    if (keyid != null) result.keyid = keyid;
    if (signingalgorithm != null) result.signingalgorithm = signingalgorithm;
    return result;
  }

  VerifyResponse._();

  factory VerifyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VerifyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VerifyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOB(272180330, _omitFieldNames ? '' : 'signaturevalid')
    ..aOS(275906594, _omitFieldNames ? '' : 'keyid')
    ..aE<SigningAlgorithmSpec>(
        488091842, _omitFieldNames ? '' : 'signingalgorithm',
        enumValues: SigningAlgorithmSpec.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VerifyResponse copyWith(void Function(VerifyResponse) updates) =>
      super.copyWith((message) => updates(message as VerifyResponse))
          as VerifyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerifyResponse create() => VerifyResponse._();
  @$core.override
  VerifyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VerifyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VerifyResponse>(create);
  static VerifyResponse? _defaultInstance;

  @$pb.TagNumber(272180330)
  $core.bool get signaturevalid => $_getBF(0);
  @$pb.TagNumber(272180330)
  set signaturevalid($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(272180330)
  $core.bool hasSignaturevalid() => $_has(0);
  @$pb.TagNumber(272180330)
  void clearSignaturevalid() => $_clearField(272180330);

  @$pb.TagNumber(275906594)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(275906594)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275906594)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(275906594)
  void clearKeyid() => $_clearField(275906594);

  @$pb.TagNumber(488091842)
  SigningAlgorithmSpec get signingalgorithm => $_getN(2);
  @$pb.TagNumber(488091842)
  set signingalgorithm(SigningAlgorithmSpec value) =>
      $_setField(488091842, value);
  @$pb.TagNumber(488091842)
  $core.bool hasSigningalgorithm() => $_has(2);
  @$pb.TagNumber(488091842)
  void clearSigningalgorithm() => $_clearField(488091842);
}

class XksKeyAlreadyInUseException extends $pb.GeneratedMessage {
  factory XksKeyAlreadyInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksKeyAlreadyInUseException._();

  factory XksKeyAlreadyInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksKeyAlreadyInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksKeyAlreadyInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksKeyAlreadyInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksKeyAlreadyInUseException copyWith(
          void Function(XksKeyAlreadyInUseException) updates) =>
      super.copyWith(
              (message) => updates(message as XksKeyAlreadyInUseException))
          as XksKeyAlreadyInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksKeyAlreadyInUseException create() =>
      XksKeyAlreadyInUseException._();
  @$core.override
  XksKeyAlreadyInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksKeyAlreadyInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XksKeyAlreadyInUseException>(create);
  static XksKeyAlreadyInUseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksKeyConfigurationType extends $pb.GeneratedMessage {
  factory XksKeyConfigurationType({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  XksKeyConfigurationType._();

  factory XksKeyConfigurationType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksKeyConfigurationType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksKeyConfigurationType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksKeyConfigurationType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksKeyConfigurationType copyWith(
          void Function(XksKeyConfigurationType) updates) =>
      super.copyWith((message) => updates(message as XksKeyConfigurationType))
          as XksKeyConfigurationType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksKeyConfigurationType create() => XksKeyConfigurationType._();
  @$core.override
  XksKeyConfigurationType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksKeyConfigurationType getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XksKeyConfigurationType>(create);
  static XksKeyConfigurationType? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class XksKeyInvalidConfigurationException extends $pb.GeneratedMessage {
  factory XksKeyInvalidConfigurationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksKeyInvalidConfigurationException._();

  factory XksKeyInvalidConfigurationException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksKeyInvalidConfigurationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksKeyInvalidConfigurationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksKeyInvalidConfigurationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksKeyInvalidConfigurationException copyWith(
          void Function(XksKeyInvalidConfigurationException) updates) =>
      super.copyWith((message) =>
              updates(message as XksKeyInvalidConfigurationException))
          as XksKeyInvalidConfigurationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksKeyInvalidConfigurationException create() =>
      XksKeyInvalidConfigurationException._();
  @$core.override
  XksKeyInvalidConfigurationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksKeyInvalidConfigurationException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          XksKeyInvalidConfigurationException>(create);
  static XksKeyInvalidConfigurationException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksKeyNotFoundException extends $pb.GeneratedMessage {
  factory XksKeyNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksKeyNotFoundException._();

  factory XksKeyNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksKeyNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksKeyNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksKeyNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksKeyNotFoundException copyWith(
          void Function(XksKeyNotFoundException) updates) =>
      super.copyWith((message) => updates(message as XksKeyNotFoundException))
          as XksKeyNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksKeyNotFoundException create() => XksKeyNotFoundException._();
  @$core.override
  XksKeyNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksKeyNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XksKeyNotFoundException>(create);
  static XksKeyNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyAuthenticationCredentialType extends $pb.GeneratedMessage {
  factory XksProxyAuthenticationCredentialType({
    $core.String? rawsecretaccesskey,
    $core.String? accesskeyid,
  }) {
    final result = create();
    if (rawsecretaccesskey != null)
      result.rawsecretaccesskey = rawsecretaccesskey;
    if (accesskeyid != null) result.accesskeyid = accesskeyid;
    return result;
  }

  XksProxyAuthenticationCredentialType._();

  factory XksProxyAuthenticationCredentialType.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyAuthenticationCredentialType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyAuthenticationCredentialType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(68137571, _omitFieldNames ? '' : 'rawsecretaccesskey')
    ..aOS(453893024, _omitFieldNames ? '' : 'accesskeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyAuthenticationCredentialType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyAuthenticationCredentialType copyWith(
          void Function(XksProxyAuthenticationCredentialType) updates) =>
      super.copyWith((message) =>
              updates(message as XksProxyAuthenticationCredentialType))
          as XksProxyAuthenticationCredentialType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyAuthenticationCredentialType create() =>
      XksProxyAuthenticationCredentialType._();
  @$core.override
  XksProxyAuthenticationCredentialType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyAuthenticationCredentialType getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          XksProxyAuthenticationCredentialType>(create);
  static XksProxyAuthenticationCredentialType? _defaultInstance;

  @$pb.TagNumber(68137571)
  $core.String get rawsecretaccesskey => $_getSZ(0);
  @$pb.TagNumber(68137571)
  set rawsecretaccesskey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(68137571)
  $core.bool hasRawsecretaccesskey() => $_has(0);
  @$pb.TagNumber(68137571)
  void clearRawsecretaccesskey() => $_clearField(68137571);

  @$pb.TagNumber(453893024)
  $core.String get accesskeyid => $_getSZ(1);
  @$pb.TagNumber(453893024)
  set accesskeyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(453893024)
  $core.bool hasAccesskeyid() => $_has(1);
  @$pb.TagNumber(453893024)
  void clearAccesskeyid() => $_clearField(453893024);
}

class XksProxyConfigurationType extends $pb.GeneratedMessage {
  factory XksProxyConfigurationType({
    $core.String? uriendpoint,
    XksProxyConnectivityType? connectivity,
    $core.String? vpcendpointservicename,
    $core.String? uripath,
    $core.String? vpcendpointserviceowner,
    $core.String? accesskeyid,
  }) {
    final result = create();
    if (uriendpoint != null) result.uriendpoint = uriendpoint;
    if (connectivity != null) result.connectivity = connectivity;
    if (vpcendpointservicename != null)
      result.vpcendpointservicename = vpcendpointservicename;
    if (uripath != null) result.uripath = uripath;
    if (vpcendpointserviceowner != null)
      result.vpcendpointserviceowner = vpcendpointserviceowner;
    if (accesskeyid != null) result.accesskeyid = accesskeyid;
    return result;
  }

  XksProxyConfigurationType._();

  factory XksProxyConfigurationType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyConfigurationType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyConfigurationType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(79142005, _omitFieldNames ? '' : 'uriendpoint')
    ..aE<XksProxyConnectivityType>(
        210638147, _omitFieldNames ? '' : 'connectivity',
        enumValues: XksProxyConnectivityType.values)
    ..aOS(269882444, _omitFieldNames ? '' : 'vpcendpointservicename')
    ..aOS(288340351, _omitFieldNames ? '' : 'uripath')
    ..aOS(298819456, _omitFieldNames ? '' : 'vpcendpointserviceowner')
    ..aOS(453893024, _omitFieldNames ? '' : 'accesskeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyConfigurationType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyConfigurationType copyWith(
          void Function(XksProxyConfigurationType) updates) =>
      super.copyWith((message) => updates(message as XksProxyConfigurationType))
          as XksProxyConfigurationType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyConfigurationType create() => XksProxyConfigurationType._();
  @$core.override
  XksProxyConfigurationType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyConfigurationType getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XksProxyConfigurationType>(create);
  static XksProxyConfigurationType? _defaultInstance;

  @$pb.TagNumber(79142005)
  $core.String get uriendpoint => $_getSZ(0);
  @$pb.TagNumber(79142005)
  set uriendpoint($core.String value) => $_setString(0, value);
  @$pb.TagNumber(79142005)
  $core.bool hasUriendpoint() => $_has(0);
  @$pb.TagNumber(79142005)
  void clearUriendpoint() => $_clearField(79142005);

  @$pb.TagNumber(210638147)
  XksProxyConnectivityType get connectivity => $_getN(1);
  @$pb.TagNumber(210638147)
  set connectivity(XksProxyConnectivityType value) =>
      $_setField(210638147, value);
  @$pb.TagNumber(210638147)
  $core.bool hasConnectivity() => $_has(1);
  @$pb.TagNumber(210638147)
  void clearConnectivity() => $_clearField(210638147);

  @$pb.TagNumber(269882444)
  $core.String get vpcendpointservicename => $_getSZ(2);
  @$pb.TagNumber(269882444)
  set vpcendpointservicename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(269882444)
  $core.bool hasVpcendpointservicename() => $_has(2);
  @$pb.TagNumber(269882444)
  void clearVpcendpointservicename() => $_clearField(269882444);

  @$pb.TagNumber(288340351)
  $core.String get uripath => $_getSZ(3);
  @$pb.TagNumber(288340351)
  set uripath($core.String value) => $_setString(3, value);
  @$pb.TagNumber(288340351)
  $core.bool hasUripath() => $_has(3);
  @$pb.TagNumber(288340351)
  void clearUripath() => $_clearField(288340351);

  @$pb.TagNumber(298819456)
  $core.String get vpcendpointserviceowner => $_getSZ(4);
  @$pb.TagNumber(298819456)
  set vpcendpointserviceowner($core.String value) => $_setString(4, value);
  @$pb.TagNumber(298819456)
  $core.bool hasVpcendpointserviceowner() => $_has(4);
  @$pb.TagNumber(298819456)
  void clearVpcendpointserviceowner() => $_clearField(298819456);

  @$pb.TagNumber(453893024)
  $core.String get accesskeyid => $_getSZ(5);
  @$pb.TagNumber(453893024)
  set accesskeyid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(453893024)
  $core.bool hasAccesskeyid() => $_has(5);
  @$pb.TagNumber(453893024)
  void clearAccesskeyid() => $_clearField(453893024);
}

class XksProxyIncorrectAuthenticationCredentialException
    extends $pb.GeneratedMessage {
  factory XksProxyIncorrectAuthenticationCredentialException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyIncorrectAuthenticationCredentialException._();

  factory XksProxyIncorrectAuthenticationCredentialException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyIncorrectAuthenticationCredentialException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames
          ? ''
          : 'XksProxyIncorrectAuthenticationCredentialException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyIncorrectAuthenticationCredentialException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyIncorrectAuthenticationCredentialException copyWith(
          void Function(XksProxyIncorrectAuthenticationCredentialException)
              updates) =>
      super.copyWith((message) => updates(
              message as XksProxyIncorrectAuthenticationCredentialException))
          as XksProxyIncorrectAuthenticationCredentialException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyIncorrectAuthenticationCredentialException create() =>
      XksProxyIncorrectAuthenticationCredentialException._();
  @$core.override
  XksProxyIncorrectAuthenticationCredentialException createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static XksProxyIncorrectAuthenticationCredentialException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          XksProxyIncorrectAuthenticationCredentialException>(create);
  static XksProxyIncorrectAuthenticationCredentialException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyInvalidConfigurationException extends $pb.GeneratedMessage {
  factory XksProxyInvalidConfigurationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyInvalidConfigurationException._();

  factory XksProxyInvalidConfigurationException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyInvalidConfigurationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyInvalidConfigurationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyInvalidConfigurationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyInvalidConfigurationException copyWith(
          void Function(XksProxyInvalidConfigurationException) updates) =>
      super.copyWith((message) =>
              updates(message as XksProxyInvalidConfigurationException))
          as XksProxyInvalidConfigurationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyInvalidConfigurationException create() =>
      XksProxyInvalidConfigurationException._();
  @$core.override
  XksProxyInvalidConfigurationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyInvalidConfigurationException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          XksProxyInvalidConfigurationException>(create);
  static XksProxyInvalidConfigurationException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyInvalidResponseException extends $pb.GeneratedMessage {
  factory XksProxyInvalidResponseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyInvalidResponseException._();

  factory XksProxyInvalidResponseException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyInvalidResponseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyInvalidResponseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyInvalidResponseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyInvalidResponseException copyWith(
          void Function(XksProxyInvalidResponseException) updates) =>
      super.copyWith(
              (message) => updates(message as XksProxyInvalidResponseException))
          as XksProxyInvalidResponseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyInvalidResponseException create() =>
      XksProxyInvalidResponseException._();
  @$core.override
  XksProxyInvalidResponseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyInvalidResponseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XksProxyInvalidResponseException>(
          create);
  static XksProxyInvalidResponseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyUriEndpointInUseException extends $pb.GeneratedMessage {
  factory XksProxyUriEndpointInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyUriEndpointInUseException._();

  factory XksProxyUriEndpointInUseException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyUriEndpointInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyUriEndpointInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyUriEndpointInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyUriEndpointInUseException copyWith(
          void Function(XksProxyUriEndpointInUseException) updates) =>
      super.copyWith((message) =>
              updates(message as XksProxyUriEndpointInUseException))
          as XksProxyUriEndpointInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyUriEndpointInUseException create() =>
      XksProxyUriEndpointInUseException._();
  @$core.override
  XksProxyUriEndpointInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyUriEndpointInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XksProxyUriEndpointInUseException>(
          create);
  static XksProxyUriEndpointInUseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyUriInUseException extends $pb.GeneratedMessage {
  factory XksProxyUriInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyUriInUseException._();

  factory XksProxyUriInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyUriInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyUriInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyUriInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyUriInUseException copyWith(
          void Function(XksProxyUriInUseException) updates) =>
      super.copyWith((message) => updates(message as XksProxyUriInUseException))
          as XksProxyUriInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyUriInUseException create() => XksProxyUriInUseException._();
  @$core.override
  XksProxyUriInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyUriInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XksProxyUriInUseException>(create);
  static XksProxyUriInUseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyUriUnreachableException extends $pb.GeneratedMessage {
  factory XksProxyUriUnreachableException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyUriUnreachableException._();

  factory XksProxyUriUnreachableException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyUriUnreachableException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyUriUnreachableException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyUriUnreachableException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyUriUnreachableException copyWith(
          void Function(XksProxyUriUnreachableException) updates) =>
      super.copyWith(
              (message) => updates(message as XksProxyUriUnreachableException))
          as XksProxyUriUnreachableException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyUriUnreachableException create() =>
      XksProxyUriUnreachableException._();
  @$core.override
  XksProxyUriUnreachableException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyUriUnreachableException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XksProxyUriUnreachableException>(
          create);
  static XksProxyUriUnreachableException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyVpcEndpointServiceInUseException extends $pb.GeneratedMessage {
  factory XksProxyVpcEndpointServiceInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyVpcEndpointServiceInUseException._();

  factory XksProxyVpcEndpointServiceInUseException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyVpcEndpointServiceInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyVpcEndpointServiceInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyVpcEndpointServiceInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyVpcEndpointServiceInUseException copyWith(
          void Function(XksProxyVpcEndpointServiceInUseException) updates) =>
      super.copyWith((message) =>
              updates(message as XksProxyVpcEndpointServiceInUseException))
          as XksProxyVpcEndpointServiceInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyVpcEndpointServiceInUseException create() =>
      XksProxyVpcEndpointServiceInUseException._();
  @$core.override
  XksProxyVpcEndpointServiceInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyVpcEndpointServiceInUseException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          XksProxyVpcEndpointServiceInUseException>(create);
  static XksProxyVpcEndpointServiceInUseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyVpcEndpointServiceInvalidConfigurationException
    extends $pb.GeneratedMessage {
  factory XksProxyVpcEndpointServiceInvalidConfigurationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyVpcEndpointServiceInvalidConfigurationException._();

  factory XksProxyVpcEndpointServiceInvalidConfigurationException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyVpcEndpointServiceInvalidConfigurationException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames
          ? ''
          : 'XksProxyVpcEndpointServiceInvalidConfigurationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyVpcEndpointServiceInvalidConfigurationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyVpcEndpointServiceInvalidConfigurationException copyWith(
          void Function(XksProxyVpcEndpointServiceInvalidConfigurationException)
              updates) =>
      super.copyWith((message) => updates(message
              as XksProxyVpcEndpointServiceInvalidConfigurationException))
          as XksProxyVpcEndpointServiceInvalidConfigurationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyVpcEndpointServiceInvalidConfigurationException create() =>
      XksProxyVpcEndpointServiceInvalidConfigurationException._();
  @$core.override
  XksProxyVpcEndpointServiceInvalidConfigurationException
      createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyVpcEndpointServiceInvalidConfigurationException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          XksProxyVpcEndpointServiceInvalidConfigurationException>(create);
  static XksProxyVpcEndpointServiceInvalidConfigurationException?
      _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class XksProxyVpcEndpointServiceNotFoundException extends $pb.GeneratedMessage {
  factory XksProxyVpcEndpointServiceNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  XksProxyVpcEndpointServiceNotFoundException._();

  factory XksProxyVpcEndpointServiceNotFoundException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XksProxyVpcEndpointServiceNotFoundException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XksProxyVpcEndpointServiceNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'kms'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyVpcEndpointServiceNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XksProxyVpcEndpointServiceNotFoundException copyWith(
          void Function(XksProxyVpcEndpointServiceNotFoundException) updates) =>
      super.copyWith((message) =>
              updates(message as XksProxyVpcEndpointServiceNotFoundException))
          as XksProxyVpcEndpointServiceNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XksProxyVpcEndpointServiceNotFoundException create() =>
      XksProxyVpcEndpointServiceNotFoundException._();
  @$core.override
  XksProxyVpcEndpointServiceNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XksProxyVpcEndpointServiceNotFoundException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          XksProxyVpcEndpointServiceNotFoundException>(create);
  static XksProxyVpcEndpointServiceNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
