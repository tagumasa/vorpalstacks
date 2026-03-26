// This is a generated file - do not edit.
//
// Generated from secretsmanager.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'secretsmanager.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'secretsmanager.pbenum.dart';

class APIErrorType extends $pb.GeneratedMessage {
  factory APIErrorType({
    $core.String? errorcode,
    $core.String? message,
    $core.String? secretid,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (message != null) result.message = message;
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  APIErrorType._();

  factory APIErrorType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory APIErrorType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'APIErrorType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  APIErrorType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  APIErrorType copyWith(void Function(APIErrorType) updates) =>
      super.copyWith((message) => updates(message as APIErrorType))
          as APIErrorType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static APIErrorType create() => APIErrorType._();
  @$core.override
  APIErrorType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static APIErrorType getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<APIErrorType>(create);
  static APIErrorType? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(2);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(2);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class BatchGetSecretValueRequest extends $pb.GeneratedMessage {
  factory BatchGetSecretValueRequest({
    $core.Iterable<Filter>? filters,
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.Iterable<$core.String>? secretidlist,
  }) {
    final result = create();
    if (filters != null) result.filters.addAll(filters);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (secretidlist != null) result.secretidlist.addAll(secretidlist);
    return result;
  }

  BatchGetSecretValueRequest._();

  factory BatchGetSecretValueRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetSecretValueRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetSecretValueRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..pPM<Filter>(188393197, _omitFieldNames ? '' : 'filters',
        subBuilder: Filter.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..pPS(398967021, _omitFieldNames ? '' : 'secretidlist')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetSecretValueRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetSecretValueRequest copyWith(
          void Function(BatchGetSecretValueRequest) updates) =>
      super.copyWith(
              (message) => updates(message as BatchGetSecretValueRequest))
          as BatchGetSecretValueRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetSecretValueRequest create() => BatchGetSecretValueRequest._();
  @$core.override
  BatchGetSecretValueRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetSecretValueRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetSecretValueRequest>(create);
  static BatchGetSecretValueRequest? _defaultInstance;

  @$pb.TagNumber(188393197)
  $pb.PbList<Filter> get filters => $_getList(0);

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

  @$pb.TagNumber(398967021)
  $pb.PbList<$core.String> get secretidlist => $_getList(3);
}

class BatchGetSecretValueResponse extends $pb.GeneratedMessage {
  factory BatchGetSecretValueResponse({
    $core.Iterable<SecretValueEntry>? secretvalues,
    $core.Iterable<APIErrorType>? errors,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (secretvalues != null) result.secretvalues.addAll(secretvalues);
    if (errors != null) result.errors.addAll(errors);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  BatchGetSecretValueResponse._();

  factory BatchGetSecretValueResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetSecretValueResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetSecretValueResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..pPM<SecretValueEntry>(79551512, _omitFieldNames ? '' : 'secretvalues',
        subBuilder: SecretValueEntry.create)
    ..pPM<APIErrorType>(166551719, _omitFieldNames ? '' : 'errors',
        subBuilder: APIErrorType.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetSecretValueResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetSecretValueResponse copyWith(
          void Function(BatchGetSecretValueResponse) updates) =>
      super.copyWith(
              (message) => updates(message as BatchGetSecretValueResponse))
          as BatchGetSecretValueResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetSecretValueResponse create() =>
      BatchGetSecretValueResponse._();
  @$core.override
  BatchGetSecretValueResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetSecretValueResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetSecretValueResponse>(create);
  static BatchGetSecretValueResponse? _defaultInstance;

  @$pb.TagNumber(79551512)
  $pb.PbList<SecretValueEntry> get secretvalues => $_getList(0);

  @$pb.TagNumber(166551719)
  $pb.PbList<APIErrorType> get errors => $_getList(1);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class CancelRotateSecretRequest extends $pb.GeneratedMessage {
  factory CancelRotateSecretRequest({
    $core.String? secretid,
  }) {
    final result = create();
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  CancelRotateSecretRequest._();

  factory CancelRotateSecretRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelRotateSecretRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelRotateSecretRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelRotateSecretRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelRotateSecretRequest copyWith(
          void Function(CancelRotateSecretRequest) updates) =>
      super.copyWith((message) => updates(message as CancelRotateSecretRequest))
          as CancelRotateSecretRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelRotateSecretRequest create() => CancelRotateSecretRequest._();
  @$core.override
  CancelRotateSecretRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelRotateSecretRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelRotateSecretRequest>(create);
  static CancelRotateSecretRequest? _defaultInstance;

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(0);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(0);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class CancelRotateSecretResponse extends $pb.GeneratedMessage {
  factory CancelRotateSecretResponse({
    $core.String? name,
    $core.String? versionid,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (versionid != null) result.versionid = versionid;
    if (arn != null) result.arn = arn;
    return result;
  }

  CancelRotateSecretResponse._();

  factory CancelRotateSecretResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelRotateSecretResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelRotateSecretResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelRotateSecretResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelRotateSecretResponse copyWith(
          void Function(CancelRotateSecretResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CancelRotateSecretResponse))
          as CancelRotateSecretResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelRotateSecretResponse create() => CancelRotateSecretResponse._();
  @$core.override
  CancelRotateSecretResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelRotateSecretResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelRotateSecretResponse>(create);
  static CancelRotateSecretResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(1);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(1);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class CreateSecretRequest extends $pb.GeneratedMessage {
  factory CreateSecretRequest({
    $core.String? kmskeyid,
    $core.List<$core.int>? secretbinary,
    $core.String? description,
    $core.String? secretstring,
    $core.bool? forceoverwritereplicasecret,
    $core.String? name,
    $core.String? type,
    $core.Iterable<Tag>? tags,
    $core.String? clientrequesttoken,
    $core.Iterable<ReplicaRegionType>? addreplicaregions,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (secretbinary != null) result.secretbinary = secretbinary;
    if (description != null) result.description = description;
    if (secretstring != null) result.secretstring = secretstring;
    if (forceoverwritereplicasecret != null)
      result.forceoverwritereplicasecret = forceoverwritereplicasecret;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (tags != null) result.tags.addAll(tags);
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (addreplicaregions != null)
      result.addreplicaregions.addAll(addreplicaregions);
    return result;
  }

  CreateSecretRequest._();

  factory CreateSecretRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateSecretRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateSecretRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..a<$core.List<$core.int>>(
        94375681, _omitFieldNames ? '' : 'secretbinary', $pb.PbFieldType.OY)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(190782253, _omitFieldNames ? '' : 'secretstring')
    ..aOB(247407324, _omitFieldNames ? '' : 'forceoverwritereplicasecret')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..pPM<ReplicaRegionType>(
        461171870, _omitFieldNames ? '' : 'addreplicaregions',
        subBuilder: ReplicaRegionType.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSecretRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSecretRequest copyWith(void Function(CreateSecretRequest) updates) =>
      super.copyWith((message) => updates(message as CreateSecretRequest))
          as CreateSecretRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateSecretRequest create() => CreateSecretRequest._();
  @$core.override
  CreateSecretRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateSecretRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateSecretRequest>(create);
  static CreateSecretRequest? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(94375681)
  $core.List<$core.int> get secretbinary => $_getN(1);
  @$pb.TagNumber(94375681)
  set secretbinary($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(94375681)
  $core.bool hasSecretbinary() => $_has(1);
  @$pb.TagNumber(94375681)
  void clearSecretbinary() => $_clearField(94375681);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(190782253)
  $core.String get secretstring => $_getSZ(3);
  @$pb.TagNumber(190782253)
  set secretstring($core.String value) => $_setString(3, value);
  @$pb.TagNumber(190782253)
  $core.bool hasSecretstring() => $_has(3);
  @$pb.TagNumber(190782253)
  void clearSecretstring() => $_clearField(190782253);

  @$pb.TagNumber(247407324)
  $core.bool get forceoverwritereplicasecret => $_getBF(4);
  @$pb.TagNumber(247407324)
  set forceoverwritereplicasecret($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(247407324)
  $core.bool hasForceoverwritereplicasecret() => $_has(4);
  @$pb.TagNumber(247407324)
  void clearForceoverwritereplicasecret() => $_clearField(247407324);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  $core.String get type => $_getSZ(6);
  @$pb.TagNumber(290836590)
  set type($core.String value) => $_setString(6, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(6);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(7);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(8);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(8, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(8);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(461171870)
  $pb.PbList<ReplicaRegionType> get addreplicaregions => $_getList(9);
}

class CreateSecretResponse extends $pb.GeneratedMessage {
  factory CreateSecretResponse({
    $core.String? name,
    $core.String? versionid,
    $core.String? arn,
    $core.Iterable<ReplicationStatusType>? replicationstatus,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (versionid != null) result.versionid = versionid;
    if (arn != null) result.arn = arn;
    if (replicationstatus != null)
      result.replicationstatus.addAll(replicationstatus);
    return result;
  }

  CreateSecretResponse._();

  factory CreateSecretResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateSecretResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateSecretResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..pPM<ReplicationStatusType>(
        529093900, _omitFieldNames ? '' : 'replicationstatus',
        subBuilder: ReplicationStatusType.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSecretResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSecretResponse copyWith(void Function(CreateSecretResponse) updates) =>
      super.copyWith((message) => updates(message as CreateSecretResponse))
          as CreateSecretResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateSecretResponse create() => CreateSecretResponse._();
  @$core.override
  CreateSecretResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateSecretResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateSecretResponse>(create);
  static CreateSecretResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(1);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(1);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(529093900)
  $pb.PbList<ReplicationStatusType> get replicationstatus => $_getList(3);
}

class DecryptionFailure extends $pb.GeneratedMessage {
  factory DecryptionFailure({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DecryptionFailure._();

  factory DecryptionFailure.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DecryptionFailure.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DecryptionFailure',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecryptionFailure clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecryptionFailure copyWith(void Function(DecryptionFailure) updates) =>
      super.copyWith((message) => updates(message as DecryptionFailure))
          as DecryptionFailure;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DecryptionFailure create() => DecryptionFailure._();
  @$core.override
  DecryptionFailure createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DecryptionFailure getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DecryptionFailure>(create);
  static DecryptionFailure? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class DeleteResourcePolicyRequest extends $pb.GeneratedMessage {
  factory DeleteResourcePolicyRequest({
    $core.String? secretid,
  }) {
    final result = create();
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  DeleteResourcePolicyRequest._();

  factory DeleteResourcePolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteResourcePolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteResourcePolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourcePolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourcePolicyRequest copyWith(
          void Function(DeleteResourcePolicyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteResourcePolicyRequest))
          as DeleteResourcePolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteResourcePolicyRequest create() =>
      DeleteResourcePolicyRequest._();
  @$core.override
  DeleteResourcePolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteResourcePolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteResourcePolicyRequest>(create);
  static DeleteResourcePolicyRequest? _defaultInstance;

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(0);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(0);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class DeleteResourcePolicyResponse extends $pb.GeneratedMessage {
  factory DeleteResourcePolicyResponse({
    $core.String? name,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    return result;
  }

  DeleteResourcePolicyResponse._();

  factory DeleteResourcePolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteResourcePolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteResourcePolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourcePolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourcePolicyResponse copyWith(
          void Function(DeleteResourcePolicyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteResourcePolicyResponse))
          as DeleteResourcePolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteResourcePolicyResponse create() =>
      DeleteResourcePolicyResponse._();
  @$core.override
  DeleteResourcePolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteResourcePolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteResourcePolicyResponse>(create);
  static DeleteResourcePolicyResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class DeleteSecretRequest extends $pb.GeneratedMessage {
  factory DeleteSecretRequest({
    $fixnum.Int64? recoverywindowindays,
    $core.bool? forcedeletewithoutrecovery,
    $core.String? secretid,
  }) {
    final result = create();
    if (recoverywindowindays != null)
      result.recoverywindowindays = recoverywindowindays;
    if (forcedeletewithoutrecovery != null)
      result.forcedeletewithoutrecovery = forcedeletewithoutrecovery;
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  DeleteSecretRequest._();

  factory DeleteSecretRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteSecretRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteSecretRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aInt64(141399243, _omitFieldNames ? '' : 'recoverywindowindays')
    ..aOB(179882315, _omitFieldNames ? '' : 'forcedeletewithoutrecovery')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSecretRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSecretRequest copyWith(void Function(DeleteSecretRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteSecretRequest))
          as DeleteSecretRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteSecretRequest create() => DeleteSecretRequest._();
  @$core.override
  DeleteSecretRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteSecretRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteSecretRequest>(create);
  static DeleteSecretRequest? _defaultInstance;

  @$pb.TagNumber(141399243)
  $fixnum.Int64 get recoverywindowindays => $_getI64(0);
  @$pb.TagNumber(141399243)
  set recoverywindowindays($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(141399243)
  $core.bool hasRecoverywindowindays() => $_has(0);
  @$pb.TagNumber(141399243)
  void clearRecoverywindowindays() => $_clearField(141399243);

  @$pb.TagNumber(179882315)
  $core.bool get forcedeletewithoutrecovery => $_getBF(1);
  @$pb.TagNumber(179882315)
  set forcedeletewithoutrecovery($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(179882315)
  $core.bool hasForcedeletewithoutrecovery() => $_has(1);
  @$pb.TagNumber(179882315)
  void clearForcedeletewithoutrecovery() => $_clearField(179882315);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(2);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(2);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class DeleteSecretResponse extends $pb.GeneratedMessage {
  factory DeleteSecretResponse({
    $core.String? name,
    $core.String? deletiondate,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (deletiondate != null) result.deletiondate = deletiondate;
    if (arn != null) result.arn = arn;
    return result;
  }

  DeleteSecretResponse._();

  factory DeleteSecretResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteSecretResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteSecretResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(347845564, _omitFieldNames ? '' : 'deletiondate')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSecretResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSecretResponse copyWith(void Function(DeleteSecretResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteSecretResponse))
          as DeleteSecretResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteSecretResponse create() => DeleteSecretResponse._();
  @$core.override
  DeleteSecretResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteSecretResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteSecretResponse>(create);
  static DeleteSecretResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(347845564)
  $core.String get deletiondate => $_getSZ(1);
  @$pb.TagNumber(347845564)
  set deletiondate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(347845564)
  $core.bool hasDeletiondate() => $_has(1);
  @$pb.TagNumber(347845564)
  void clearDeletiondate() => $_clearField(347845564);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class DescribeSecretRequest extends $pb.GeneratedMessage {
  factory DescribeSecretRequest({
    $core.String? secretid,
  }) {
    final result = create();
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  DescribeSecretRequest._();

  factory DescribeSecretRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeSecretRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeSecretRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeSecretRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeSecretRequest copyWith(
          void Function(DescribeSecretRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeSecretRequest))
          as DescribeSecretRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeSecretRequest create() => DescribeSecretRequest._();
  @$core.override
  DescribeSecretRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeSecretRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeSecretRequest>(create);
  static DescribeSecretRequest? _defaultInstance;

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(0);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(0);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class DescribeSecretResponse extends $pb.GeneratedMessage {
  factory DescribeSecretResponse({
    $core.String? kmskeyid,
    $core.Iterable<ExternalSecretRotationMetadataItem>?
        externalsecretrotationmetadata,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        versionidstostages,
    $core.String? description,
    $core.String? nextrotationdate,
    $core.String? lastaccesseddate,
    $core.bool? rotationenabled,
    RotationRulesType? rotationrules,
    $core.String? name,
    $core.String? type,
    $core.String? lastchangeddate,
    $core.String? rotationlambdaarn,
    $core.Iterable<Tag>? tags,
    $core.String? arn,
    $core.String? createddate,
    $core.String? owningservice,
    $core.String? externalsecretrotationrolearn,
    $core.String? primaryregion,
    $core.String? lastrotateddate,
    $core.String? deleteddate,
    $core.Iterable<ReplicationStatusType>? replicationstatus,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (externalsecretrotationmetadata != null)
      result.externalsecretrotationmetadata
          .addAll(externalsecretrotationmetadata);
    if (versionidstostages != null)
      result.versionidstostages.addEntries(versionidstostages);
    if (description != null) result.description = description;
    if (nextrotationdate != null) result.nextrotationdate = nextrotationdate;
    if (lastaccesseddate != null) result.lastaccesseddate = lastaccesseddate;
    if (rotationenabled != null) result.rotationenabled = rotationenabled;
    if (rotationrules != null) result.rotationrules = rotationrules;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (lastchangeddate != null) result.lastchangeddate = lastchangeddate;
    if (rotationlambdaarn != null) result.rotationlambdaarn = rotationlambdaarn;
    if (tags != null) result.tags.addAll(tags);
    if (arn != null) result.arn = arn;
    if (createddate != null) result.createddate = createddate;
    if (owningservice != null) result.owningservice = owningservice;
    if (externalsecretrotationrolearn != null)
      result.externalsecretrotationrolearn = externalsecretrotationrolearn;
    if (primaryregion != null) result.primaryregion = primaryregion;
    if (lastrotateddate != null) result.lastrotateddate = lastrotateddate;
    if (deleteddate != null) result.deleteddate = deleteddate;
    if (replicationstatus != null)
      result.replicationstatus.addAll(replicationstatus);
    return result;
  }

  DescribeSecretResponse._();

  factory DescribeSecretResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeSecretResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeSecretResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..pPM<ExternalSecretRotationMetadataItem>(
        52900542, _omitFieldNames ? '' : 'externalsecretrotationmetadata',
        subBuilder: ExternalSecretRotationMetadataItem.create)
    ..m<$core.String, $core.String>(
        90698314, _omitFieldNames ? '' : 'versionidstostages',
        entryClassName: 'DescribeSecretResponse.VersionidstostagesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('secretsmanager'))
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(192035355, _omitFieldNames ? '' : 'nextrotationdate')
    ..aOS(194418963, _omitFieldNames ? '' : 'lastaccesseddate')
    ..aOB(209507301, _omitFieldNames ? '' : 'rotationenabled')
    ..aOM<RotationRulesType>(259458135, _omitFieldNames ? '' : 'rotationrules',
        subBuilder: RotationRulesType.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
    ..aOS(314015460, _omitFieldNames ? '' : 'lastchangeddate')
    ..aOS(335026080, _omitFieldNames ? '' : 'rotationlambdaarn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..aOS(416929840, _omitFieldNames ? '' : 'createddate')
    ..aOS(462487817, _omitFieldNames ? '' : 'owningservice')
    ..aOS(470712576, _omitFieldNames ? '' : 'externalsecretrotationrolearn')
    ..aOS(480901186, _omitFieldNames ? '' : 'primaryregion')
    ..aOS(501475691, _omitFieldNames ? '' : 'lastrotateddate')
    ..aOS(516314255, _omitFieldNames ? '' : 'deleteddate')
    ..pPM<ReplicationStatusType>(
        529093900, _omitFieldNames ? '' : 'replicationstatus',
        subBuilder: ReplicationStatusType.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeSecretResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeSecretResponse copyWith(
          void Function(DescribeSecretResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeSecretResponse))
          as DescribeSecretResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeSecretResponse create() => DescribeSecretResponse._();
  @$core.override
  DescribeSecretResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeSecretResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeSecretResponse>(create);
  static DescribeSecretResponse? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(52900542)
  $pb.PbList<ExternalSecretRotationMetadataItem>
      get externalsecretrotationmetadata => $_getList(1);

  @$pb.TagNumber(90698314)
  $pb.PbMap<$core.String, $core.String> get versionidstostages => $_getMap(2);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(192035355)
  $core.String get nextrotationdate => $_getSZ(4);
  @$pb.TagNumber(192035355)
  set nextrotationdate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(192035355)
  $core.bool hasNextrotationdate() => $_has(4);
  @$pb.TagNumber(192035355)
  void clearNextrotationdate() => $_clearField(192035355);

  @$pb.TagNumber(194418963)
  $core.String get lastaccesseddate => $_getSZ(5);
  @$pb.TagNumber(194418963)
  set lastaccesseddate($core.String value) => $_setString(5, value);
  @$pb.TagNumber(194418963)
  $core.bool hasLastaccesseddate() => $_has(5);
  @$pb.TagNumber(194418963)
  void clearLastaccesseddate() => $_clearField(194418963);

  @$pb.TagNumber(209507301)
  $core.bool get rotationenabled => $_getBF(6);
  @$pb.TagNumber(209507301)
  set rotationenabled($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(209507301)
  $core.bool hasRotationenabled() => $_has(6);
  @$pb.TagNumber(209507301)
  void clearRotationenabled() => $_clearField(209507301);

  @$pb.TagNumber(259458135)
  RotationRulesType get rotationrules => $_getN(7);
  @$pb.TagNumber(259458135)
  set rotationrules(RotationRulesType value) => $_setField(259458135, value);
  @$pb.TagNumber(259458135)
  $core.bool hasRotationrules() => $_has(7);
  @$pb.TagNumber(259458135)
  void clearRotationrules() => $_clearField(259458135);
  @$pb.TagNumber(259458135)
  RotationRulesType ensureRotationrules() => $_ensure(7);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(8);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(8, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(8);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  $core.String get type => $_getSZ(9);
  @$pb.TagNumber(290836590)
  set type($core.String value) => $_setString(9, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(9);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(314015460)
  $core.String get lastchangeddate => $_getSZ(10);
  @$pb.TagNumber(314015460)
  set lastchangeddate($core.String value) => $_setString(10, value);
  @$pb.TagNumber(314015460)
  $core.bool hasLastchangeddate() => $_has(10);
  @$pb.TagNumber(314015460)
  void clearLastchangeddate() => $_clearField(314015460);

  @$pb.TagNumber(335026080)
  $core.String get rotationlambdaarn => $_getSZ(11);
  @$pb.TagNumber(335026080)
  set rotationlambdaarn($core.String value) => $_setString(11, value);
  @$pb.TagNumber(335026080)
  $core.bool hasRotationlambdaarn() => $_has(11);
  @$pb.TagNumber(335026080)
  void clearRotationlambdaarn() => $_clearField(335026080);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(12);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(13);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(13, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(13);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(416929840)
  $core.String get createddate => $_getSZ(14);
  @$pb.TagNumber(416929840)
  set createddate($core.String value) => $_setString(14, value);
  @$pb.TagNumber(416929840)
  $core.bool hasCreateddate() => $_has(14);
  @$pb.TagNumber(416929840)
  void clearCreateddate() => $_clearField(416929840);

  @$pb.TagNumber(462487817)
  $core.String get owningservice => $_getSZ(15);
  @$pb.TagNumber(462487817)
  set owningservice($core.String value) => $_setString(15, value);
  @$pb.TagNumber(462487817)
  $core.bool hasOwningservice() => $_has(15);
  @$pb.TagNumber(462487817)
  void clearOwningservice() => $_clearField(462487817);

  @$pb.TagNumber(470712576)
  $core.String get externalsecretrotationrolearn => $_getSZ(16);
  @$pb.TagNumber(470712576)
  set externalsecretrotationrolearn($core.String value) =>
      $_setString(16, value);
  @$pb.TagNumber(470712576)
  $core.bool hasExternalsecretrotationrolearn() => $_has(16);
  @$pb.TagNumber(470712576)
  void clearExternalsecretrotationrolearn() => $_clearField(470712576);

  @$pb.TagNumber(480901186)
  $core.String get primaryregion => $_getSZ(17);
  @$pb.TagNumber(480901186)
  set primaryregion($core.String value) => $_setString(17, value);
  @$pb.TagNumber(480901186)
  $core.bool hasPrimaryregion() => $_has(17);
  @$pb.TagNumber(480901186)
  void clearPrimaryregion() => $_clearField(480901186);

  @$pb.TagNumber(501475691)
  $core.String get lastrotateddate => $_getSZ(18);
  @$pb.TagNumber(501475691)
  set lastrotateddate($core.String value) => $_setString(18, value);
  @$pb.TagNumber(501475691)
  $core.bool hasLastrotateddate() => $_has(18);
  @$pb.TagNumber(501475691)
  void clearLastrotateddate() => $_clearField(501475691);

  @$pb.TagNumber(516314255)
  $core.String get deleteddate => $_getSZ(19);
  @$pb.TagNumber(516314255)
  set deleteddate($core.String value) => $_setString(19, value);
  @$pb.TagNumber(516314255)
  $core.bool hasDeleteddate() => $_has(19);
  @$pb.TagNumber(516314255)
  void clearDeleteddate() => $_clearField(516314255);

  @$pb.TagNumber(529093900)
  $pb.PbList<ReplicationStatusType> get replicationstatus => $_getList(20);
}

class EncryptionFailure extends $pb.GeneratedMessage {
  factory EncryptionFailure({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  EncryptionFailure._();

  factory EncryptionFailure.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EncryptionFailure.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EncryptionFailure',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptionFailure clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EncryptionFailure copyWith(void Function(EncryptionFailure) updates) =>
      super.copyWith((message) => updates(message as EncryptionFailure))
          as EncryptionFailure;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EncryptionFailure create() => EncryptionFailure._();
  @$core.override
  EncryptionFailure createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EncryptionFailure getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EncryptionFailure>(create);
  static EncryptionFailure? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ExternalSecretRotationMetadataItem extends $pb.GeneratedMessage {
  factory ExternalSecretRotationMetadataItem({
    $core.String? key,
    $core.String? value,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    return result;
  }

  ExternalSecretRotationMetadataItem._();

  factory ExternalSecretRotationMetadataItem.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExternalSecretRotationMetadataItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExternalSecretRotationMetadataItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExternalSecretRotationMetadataItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExternalSecretRotationMetadataItem copyWith(
          void Function(ExternalSecretRotationMetadataItem) updates) =>
      super.copyWith((message) =>
              updates(message as ExternalSecretRotationMetadataItem))
          as ExternalSecretRotationMetadataItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExternalSecretRotationMetadataItem create() =>
      ExternalSecretRotationMetadataItem._();
  @$core.override
  ExternalSecretRotationMetadataItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExternalSecretRotationMetadataItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExternalSecretRotationMetadataItem>(
          create);
  static ExternalSecretRotationMetadataItem? _defaultInstance;

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

class Filter extends $pb.GeneratedMessage {
  factory Filter({
    FilterNameStringType? key,
    $core.Iterable<$core.String>? values,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (values != null) result.values.addAll(values);
    return result;
  }

  Filter._();

  factory Filter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Filter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Filter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aE<FilterNameStringType>(219859213, _omitFieldNames ? '' : 'key',
        enumValues: FilterNameStringType.values)
    ..pPS(223158876, _omitFieldNames ? '' : 'values')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Filter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Filter copyWith(void Function(Filter) updates) =>
      super.copyWith((message) => updates(message as Filter)) as Filter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Filter create() => Filter._();
  @$core.override
  Filter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Filter getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Filter>(create);
  static Filter? _defaultInstance;

  @$pb.TagNumber(219859213)
  FilterNameStringType get key => $_getN(0);
  @$pb.TagNumber(219859213)
  set key(FilterNameStringType value) => $_setField(219859213, value);
  @$pb.TagNumber(219859213)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(219859213)
  void clearKey() => $_clearField(219859213);

  @$pb.TagNumber(223158876)
  $pb.PbList<$core.String> get values => $_getList(1);
}

class GetRandomPasswordRequest extends $pb.GeneratedMessage {
  factory GetRandomPasswordRequest({
    $core.bool? excludepunctuation,
    $core.bool? requireeachincludedtype,
    $core.bool? excludenumbers,
    $fixnum.Int64? passwordlength,
    $core.bool? includespace,
    $core.bool? excludelowercase,
    $core.String? excludecharacters,
    $core.bool? excludeuppercase,
  }) {
    final result = create();
    if (excludepunctuation != null)
      result.excludepunctuation = excludepunctuation;
    if (requireeachincludedtype != null)
      result.requireeachincludedtype = requireeachincludedtype;
    if (excludenumbers != null) result.excludenumbers = excludenumbers;
    if (passwordlength != null) result.passwordlength = passwordlength;
    if (includespace != null) result.includespace = includespace;
    if (excludelowercase != null) result.excludelowercase = excludelowercase;
    if (excludecharacters != null) result.excludecharacters = excludecharacters;
    if (excludeuppercase != null) result.excludeuppercase = excludeuppercase;
    return result;
  }

  GetRandomPasswordRequest._();

  factory GetRandomPasswordRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRandomPasswordRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRandomPasswordRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOB(78530732, _omitFieldNames ? '' : 'excludepunctuation')
    ..aOB(176215282, _omitFieldNames ? '' : 'requireeachincludedtype')
    ..aOB(216382246, _omitFieldNames ? '' : 'excludenumbers')
    ..aInt64(216611365, _omitFieldNames ? '' : 'passwordlength')
    ..aOB(216842628, _omitFieldNames ? '' : 'includespace')
    ..aOB(225858843, _omitFieldNames ? '' : 'excludelowercase')
    ..aOS(335851390, _omitFieldNames ? '' : 'excludecharacters')
    ..aOB(345594144, _omitFieldNames ? '' : 'excludeuppercase')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRandomPasswordRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRandomPasswordRequest copyWith(
          void Function(GetRandomPasswordRequest) updates) =>
      super.copyWith((message) => updates(message as GetRandomPasswordRequest))
          as GetRandomPasswordRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRandomPasswordRequest create() => GetRandomPasswordRequest._();
  @$core.override
  GetRandomPasswordRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRandomPasswordRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRandomPasswordRequest>(create);
  static GetRandomPasswordRequest? _defaultInstance;

  @$pb.TagNumber(78530732)
  $core.bool get excludepunctuation => $_getBF(0);
  @$pb.TagNumber(78530732)
  set excludepunctuation($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(78530732)
  $core.bool hasExcludepunctuation() => $_has(0);
  @$pb.TagNumber(78530732)
  void clearExcludepunctuation() => $_clearField(78530732);

  @$pb.TagNumber(176215282)
  $core.bool get requireeachincludedtype => $_getBF(1);
  @$pb.TagNumber(176215282)
  set requireeachincludedtype($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(176215282)
  $core.bool hasRequireeachincludedtype() => $_has(1);
  @$pb.TagNumber(176215282)
  void clearRequireeachincludedtype() => $_clearField(176215282);

  @$pb.TagNumber(216382246)
  $core.bool get excludenumbers => $_getBF(2);
  @$pb.TagNumber(216382246)
  set excludenumbers($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(216382246)
  $core.bool hasExcludenumbers() => $_has(2);
  @$pb.TagNumber(216382246)
  void clearExcludenumbers() => $_clearField(216382246);

  @$pb.TagNumber(216611365)
  $fixnum.Int64 get passwordlength => $_getI64(3);
  @$pb.TagNumber(216611365)
  set passwordlength($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(216611365)
  $core.bool hasPasswordlength() => $_has(3);
  @$pb.TagNumber(216611365)
  void clearPasswordlength() => $_clearField(216611365);

  @$pb.TagNumber(216842628)
  $core.bool get includespace => $_getBF(4);
  @$pb.TagNumber(216842628)
  set includespace($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(216842628)
  $core.bool hasIncludespace() => $_has(4);
  @$pb.TagNumber(216842628)
  void clearIncludespace() => $_clearField(216842628);

  @$pb.TagNumber(225858843)
  $core.bool get excludelowercase => $_getBF(5);
  @$pb.TagNumber(225858843)
  set excludelowercase($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(225858843)
  $core.bool hasExcludelowercase() => $_has(5);
  @$pb.TagNumber(225858843)
  void clearExcludelowercase() => $_clearField(225858843);

  @$pb.TagNumber(335851390)
  $core.String get excludecharacters => $_getSZ(6);
  @$pb.TagNumber(335851390)
  set excludecharacters($core.String value) => $_setString(6, value);
  @$pb.TagNumber(335851390)
  $core.bool hasExcludecharacters() => $_has(6);
  @$pb.TagNumber(335851390)
  void clearExcludecharacters() => $_clearField(335851390);

  @$pb.TagNumber(345594144)
  $core.bool get excludeuppercase => $_getBF(7);
  @$pb.TagNumber(345594144)
  set excludeuppercase($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(345594144)
  $core.bool hasExcludeuppercase() => $_has(7);
  @$pb.TagNumber(345594144)
  void clearExcludeuppercase() => $_clearField(345594144);
}

class GetRandomPasswordResponse extends $pb.GeneratedMessage {
  factory GetRandomPasswordResponse({
    $core.String? randompassword,
  }) {
    final result = create();
    if (randompassword != null) result.randompassword = randompassword;
    return result;
  }

  GetRandomPasswordResponse._();

  factory GetRandomPasswordResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRandomPasswordResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRandomPasswordResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(477267856, _omitFieldNames ? '' : 'randompassword')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRandomPasswordResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRandomPasswordResponse copyWith(
          void Function(GetRandomPasswordResponse) updates) =>
      super.copyWith((message) => updates(message as GetRandomPasswordResponse))
          as GetRandomPasswordResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRandomPasswordResponse create() => GetRandomPasswordResponse._();
  @$core.override
  GetRandomPasswordResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRandomPasswordResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRandomPasswordResponse>(create);
  static GetRandomPasswordResponse? _defaultInstance;

  @$pb.TagNumber(477267856)
  $core.String get randompassword => $_getSZ(0);
  @$pb.TagNumber(477267856)
  set randompassword($core.String value) => $_setString(0, value);
  @$pb.TagNumber(477267856)
  $core.bool hasRandompassword() => $_has(0);
  @$pb.TagNumber(477267856)
  void clearRandompassword() => $_clearField(477267856);
}

class GetResourcePolicyRequest extends $pb.GeneratedMessage {
  factory GetResourcePolicyRequest({
    $core.String? secretid,
  }) {
    final result = create();
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  GetResourcePolicyRequest._();

  factory GetResourcePolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourcePolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourcePolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePolicyRequest copyWith(
          void Function(GetResourcePolicyRequest) updates) =>
      super.copyWith((message) => updates(message as GetResourcePolicyRequest))
          as GetResourcePolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourcePolicyRequest create() => GetResourcePolicyRequest._();
  @$core.override
  GetResourcePolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourcePolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourcePolicyRequest>(create);
  static GetResourcePolicyRequest? _defaultInstance;

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(0);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(0);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class GetResourcePolicyResponse extends $pb.GeneratedMessage {
  factory GetResourcePolicyResponse({
    $core.String? resourcepolicy,
    $core.String? name,
    $core.String? arn,
  }) {
    final result = create();
    if (resourcepolicy != null) result.resourcepolicy = resourcepolicy;
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    return result;
  }

  GetResourcePolicyResponse._();

  factory GetResourcePolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourcePolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourcePolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(15747632, _omitFieldNames ? '' : 'resourcepolicy')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePolicyResponse copyWith(
          void Function(GetResourcePolicyResponse) updates) =>
      super.copyWith((message) => updates(message as GetResourcePolicyResponse))
          as GetResourcePolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourcePolicyResponse create() => GetResourcePolicyResponse._();
  @$core.override
  GetResourcePolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourcePolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourcePolicyResponse>(create);
  static GetResourcePolicyResponse? _defaultInstance;

  @$pb.TagNumber(15747632)
  $core.String get resourcepolicy => $_getSZ(0);
  @$pb.TagNumber(15747632)
  set resourcepolicy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(15747632)
  $core.bool hasResourcepolicy() => $_has(0);
  @$pb.TagNumber(15747632)
  void clearResourcepolicy() => $_clearField(15747632);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class GetSecretValueRequest extends $pb.GeneratedMessage {
  factory GetSecretValueRequest({
    $core.String? versionstage,
    $core.String? versionid,
    $core.String? secretid,
  }) {
    final result = create();
    if (versionstage != null) result.versionstage = versionstage;
    if (versionid != null) result.versionid = versionid;
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  GetSecretValueRequest._();

  factory GetSecretValueRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSecretValueRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSecretValueRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(229692340, _omitFieldNames ? '' : 'versionstage')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSecretValueRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSecretValueRequest copyWith(
          void Function(GetSecretValueRequest) updates) =>
      super.copyWith((message) => updates(message as GetSecretValueRequest))
          as GetSecretValueRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSecretValueRequest create() => GetSecretValueRequest._();
  @$core.override
  GetSecretValueRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSecretValueRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSecretValueRequest>(create);
  static GetSecretValueRequest? _defaultInstance;

  @$pb.TagNumber(229692340)
  $core.String get versionstage => $_getSZ(0);
  @$pb.TagNumber(229692340)
  set versionstage($core.String value) => $_setString(0, value);
  @$pb.TagNumber(229692340)
  $core.bool hasVersionstage() => $_has(0);
  @$pb.TagNumber(229692340)
  void clearVersionstage() => $_clearField(229692340);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(1);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(1);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(2);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(2);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class GetSecretValueResponse extends $pb.GeneratedMessage {
  factory GetSecretValueResponse({
    $core.List<$core.int>? secretbinary,
    $core.String? secretstring,
    $core.Iterable<$core.String>? versionstages,
    $core.String? name,
    $core.String? versionid,
    $core.String? arn,
    $core.String? createddate,
  }) {
    final result = create();
    if (secretbinary != null) result.secretbinary = secretbinary;
    if (secretstring != null) result.secretstring = secretstring;
    if (versionstages != null) result.versionstages.addAll(versionstages);
    if (name != null) result.name = name;
    if (versionid != null) result.versionid = versionid;
    if (arn != null) result.arn = arn;
    if (createddate != null) result.createddate = createddate;
    return result;
  }

  GetSecretValueResponse._();

  factory GetSecretValueResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSecretValueResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSecretValueResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(
        94375681, _omitFieldNames ? '' : 'secretbinary', $pb.PbFieldType.OY)
    ..aOS(190782253, _omitFieldNames ? '' : 'secretstring')
    ..pPS(224220993, _omitFieldNames ? '' : 'versionstages')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..aOS(416929840, _omitFieldNames ? '' : 'createddate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSecretValueResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSecretValueResponse copyWith(
          void Function(GetSecretValueResponse) updates) =>
      super.copyWith((message) => updates(message as GetSecretValueResponse))
          as GetSecretValueResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSecretValueResponse create() => GetSecretValueResponse._();
  @$core.override
  GetSecretValueResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSecretValueResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSecretValueResponse>(create);
  static GetSecretValueResponse? _defaultInstance;

  @$pb.TagNumber(94375681)
  $core.List<$core.int> get secretbinary => $_getN(0);
  @$pb.TagNumber(94375681)
  set secretbinary($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(94375681)
  $core.bool hasSecretbinary() => $_has(0);
  @$pb.TagNumber(94375681)
  void clearSecretbinary() => $_clearField(94375681);

  @$pb.TagNumber(190782253)
  $core.String get secretstring => $_getSZ(1);
  @$pb.TagNumber(190782253)
  set secretstring($core.String value) => $_setString(1, value);
  @$pb.TagNumber(190782253)
  $core.bool hasSecretstring() => $_has(1);
  @$pb.TagNumber(190782253)
  void clearSecretstring() => $_clearField(190782253);

  @$pb.TagNumber(224220993)
  $pb.PbList<$core.String> get versionstages => $_getList(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(4);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(4);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(5);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(5);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(416929840)
  $core.String get createddate => $_getSZ(6);
  @$pb.TagNumber(416929840)
  set createddate($core.String value) => $_setString(6, value);
  @$pb.TagNumber(416929840)
  $core.bool hasCreateddate() => $_has(6);
  @$pb.TagNumber(416929840)
  void clearCreateddate() => $_clearField(416929840);
}

class InternalServiceError extends $pb.GeneratedMessage {
  factory InternalServiceError({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InternalServiceError._();

  factory InternalServiceError.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InternalServiceError.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InternalServiceError',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServiceError clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServiceError copyWith(void Function(InternalServiceError) updates) =>
      super.copyWith((message) => updates(message as InternalServiceError))
          as InternalServiceError;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InternalServiceError create() => InternalServiceError._();
  @$core.override
  InternalServiceError createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InternalServiceError getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InternalServiceError>(create);
  static InternalServiceError? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidNextTokenException extends $pb.GeneratedMessage {
  factory InvalidNextTokenException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidNextTokenException._();

  factory InvalidNextTokenException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidNextTokenException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidNextTokenException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidNextTokenException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidNextTokenException copyWith(
          void Function(InvalidNextTokenException) updates) =>
      super.copyWith((message) => updates(message as InvalidNextTokenException))
          as InvalidNextTokenException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidNextTokenException create() => InvalidNextTokenException._();
  @$core.override
  InvalidNextTokenException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidNextTokenException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidNextTokenException>(create);
  static InvalidNextTokenException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidRequestException extends $pb.GeneratedMessage {
  factory InvalidRequestException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidRequestException._();

  factory InvalidRequestException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidRequestException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidRequestException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidRequestException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidRequestException copyWith(
          void Function(InvalidRequestException) updates) =>
      super.copyWith((message) => updates(message as InvalidRequestException))
          as InvalidRequestException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidRequestException create() => InvalidRequestException._();
  @$core.override
  InvalidRequestException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidRequestException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidRequestException>(create);
  static InvalidRequestException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ListSecretVersionIdsRequest extends $pb.GeneratedMessage {
  factory ListSecretVersionIdsRequest({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.bool? includedeprecated,
    $core.String? secretid,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (includedeprecated != null) result.includedeprecated = includedeprecated;
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  ListSecretVersionIdsRequest._();

  factory ListSecretVersionIdsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSecretVersionIdsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSecretVersionIdsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOB(299058751, _omitFieldNames ? '' : 'includedeprecated')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSecretVersionIdsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSecretVersionIdsRequest copyWith(
          void Function(ListSecretVersionIdsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListSecretVersionIdsRequest))
          as ListSecretVersionIdsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSecretVersionIdsRequest create() =>
      ListSecretVersionIdsRequest._();
  @$core.override
  ListSecretVersionIdsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSecretVersionIdsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSecretVersionIdsRequest>(create);
  static ListSecretVersionIdsRequest? _defaultInstance;

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

  @$pb.TagNumber(299058751)
  $core.bool get includedeprecated => $_getBF(2);
  @$pb.TagNumber(299058751)
  set includedeprecated($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(299058751)
  $core.bool hasIncludedeprecated() => $_has(2);
  @$pb.TagNumber(299058751)
  void clearIncludedeprecated() => $_clearField(299058751);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(3);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(3);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class ListSecretVersionIdsResponse extends $pb.GeneratedMessage {
  factory ListSecretVersionIdsResponse({
    $core.String? nexttoken,
    $core.Iterable<SecretVersionsListEntry>? versions,
    $core.String? name,
    $core.String? arn,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (versions != null) result.versions.addAll(versions);
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    return result;
  }

  ListSecretVersionIdsResponse._();

  factory ListSecretVersionIdsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSecretVersionIdsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSecretVersionIdsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<SecretVersionsListEntry>(252099085, _omitFieldNames ? '' : 'versions',
        subBuilder: SecretVersionsListEntry.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSecretVersionIdsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSecretVersionIdsResponse copyWith(
          void Function(ListSecretVersionIdsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListSecretVersionIdsResponse))
          as ListSecretVersionIdsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSecretVersionIdsResponse create() =>
      ListSecretVersionIdsResponse._();
  @$core.override
  ListSecretVersionIdsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSecretVersionIdsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSecretVersionIdsResponse>(create);
  static ListSecretVersionIdsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(252099085)
  $pb.PbList<SecretVersionsListEntry> get versions => $_getList(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(3);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(3);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class ListSecretsRequest extends $pb.GeneratedMessage {
  factory ListSecretsRequest({
    $core.bool? includeplanneddeletion,
    SortByType? sortby,
    $core.Iterable<Filter>? filters,
    $core.String? nexttoken,
    SortOrderType? sortorder,
    $core.int? maxresults,
  }) {
    final result = create();
    if (includeplanneddeletion != null)
      result.includeplanneddeletion = includeplanneddeletion;
    if (sortby != null) result.sortby = sortby;
    if (filters != null) result.filters.addAll(filters);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (sortorder != null) result.sortorder = sortorder;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListSecretsRequest._();

  factory ListSecretsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSecretsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSecretsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOB(64231622, _omitFieldNames ? '' : 'includeplanneddeletion')
    ..aE<SortByType>(186052369, _omitFieldNames ? '' : 'sortby',
        enumValues: SortByType.values)
    ..pPM<Filter>(188393197, _omitFieldNames ? '' : 'filters',
        subBuilder: Filter.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aE<SortOrderType>(274231684, _omitFieldNames ? '' : 'sortorder',
        enumValues: SortOrderType.values)
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSecretsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSecretsRequest copyWith(void Function(ListSecretsRequest) updates) =>
      super.copyWith((message) => updates(message as ListSecretsRequest))
          as ListSecretsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSecretsRequest create() => ListSecretsRequest._();
  @$core.override
  ListSecretsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSecretsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSecretsRequest>(create);
  static ListSecretsRequest? _defaultInstance;

  @$pb.TagNumber(64231622)
  $core.bool get includeplanneddeletion => $_getBF(0);
  @$pb.TagNumber(64231622)
  set includeplanneddeletion($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(64231622)
  $core.bool hasIncludeplanneddeletion() => $_has(0);
  @$pb.TagNumber(64231622)
  void clearIncludeplanneddeletion() => $_clearField(64231622);

  @$pb.TagNumber(186052369)
  SortByType get sortby => $_getN(1);
  @$pb.TagNumber(186052369)
  set sortby(SortByType value) => $_setField(186052369, value);
  @$pb.TagNumber(186052369)
  $core.bool hasSortby() => $_has(1);
  @$pb.TagNumber(186052369)
  void clearSortby() => $_clearField(186052369);

  @$pb.TagNumber(188393197)
  $pb.PbList<Filter> get filters => $_getList(2);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(3);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(3, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(3);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(274231684)
  SortOrderType get sortorder => $_getN(4);
  @$pb.TagNumber(274231684)
  set sortorder(SortOrderType value) => $_setField(274231684, value);
  @$pb.TagNumber(274231684)
  $core.bool hasSortorder() => $_has(4);
  @$pb.TagNumber(274231684)
  void clearSortorder() => $_clearField(274231684);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(5);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(5);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListSecretsResponse extends $pb.GeneratedMessage {
  factory ListSecretsResponse({
    $core.String? nexttoken,
    $core.Iterable<SecretListEntry>? secretlist,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (secretlist != null) result.secretlist.addAll(secretlist);
    return result;
  }

  ListSecretsResponse._();

  factory ListSecretsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSecretsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSecretsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<SecretListEntry>(280320894, _omitFieldNames ? '' : 'secretlist',
        subBuilder: SecretListEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSecretsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSecretsResponse copyWith(void Function(ListSecretsResponse) updates) =>
      super.copyWith((message) => updates(message as ListSecretsResponse))
          as ListSecretsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSecretsResponse create() => ListSecretsResponse._();
  @$core.override
  ListSecretsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSecretsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSecretsResponse>(create);
  static ListSecretsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(280320894)
  $pb.PbList<SecretListEntry> get secretlist => $_getList(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class PreconditionNotMetException extends $pb.GeneratedMessage {
  factory PreconditionNotMetException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PreconditionNotMetException._();

  factory PreconditionNotMetException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PreconditionNotMetException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PreconditionNotMetException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PreconditionNotMetException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PreconditionNotMetException copyWith(
          void Function(PreconditionNotMetException) updates) =>
      super.copyWith(
              (message) => updates(message as PreconditionNotMetException))
          as PreconditionNotMetException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PreconditionNotMetException create() =>
      PreconditionNotMetException._();
  @$core.override
  PreconditionNotMetException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PreconditionNotMetException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PreconditionNotMetException>(create);
  static PreconditionNotMetException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class PublicPolicyException extends $pb.GeneratedMessage {
  factory PublicPolicyException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PublicPolicyException._();

  factory PublicPolicyException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublicPolicyException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublicPolicyException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublicPolicyException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublicPolicyException copyWith(
          void Function(PublicPolicyException) updates) =>
      super.copyWith((message) => updates(message as PublicPolicyException))
          as PublicPolicyException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublicPolicyException create() => PublicPolicyException._();
  @$core.override
  PublicPolicyException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublicPolicyException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublicPolicyException>(create);
  static PublicPolicyException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class PutResourcePolicyRequest extends $pb.GeneratedMessage {
  factory PutResourcePolicyRequest({
    $core.String? resourcepolicy,
    $core.String? secretid,
    $core.bool? blockpublicpolicy,
  }) {
    final result = create();
    if (resourcepolicy != null) result.resourcepolicy = resourcepolicy;
    if (secretid != null) result.secretid = secretid;
    if (blockpublicpolicy != null) result.blockpublicpolicy = blockpublicpolicy;
    return result;
  }

  PutResourcePolicyRequest._();

  factory PutResourcePolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutResourcePolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutResourcePolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(15747632, _omitFieldNames ? '' : 'resourcepolicy')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..aOB(505379326, _omitFieldNames ? '' : 'blockpublicpolicy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutResourcePolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutResourcePolicyRequest copyWith(
          void Function(PutResourcePolicyRequest) updates) =>
      super.copyWith((message) => updates(message as PutResourcePolicyRequest))
          as PutResourcePolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutResourcePolicyRequest create() => PutResourcePolicyRequest._();
  @$core.override
  PutResourcePolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutResourcePolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutResourcePolicyRequest>(create);
  static PutResourcePolicyRequest? _defaultInstance;

  @$pb.TagNumber(15747632)
  $core.String get resourcepolicy => $_getSZ(0);
  @$pb.TagNumber(15747632)
  set resourcepolicy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(15747632)
  $core.bool hasResourcepolicy() => $_has(0);
  @$pb.TagNumber(15747632)
  void clearResourcepolicy() => $_clearField(15747632);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(1);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(1);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);

  @$pb.TagNumber(505379326)
  $core.bool get blockpublicpolicy => $_getBF(2);
  @$pb.TagNumber(505379326)
  set blockpublicpolicy($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(505379326)
  $core.bool hasBlockpublicpolicy() => $_has(2);
  @$pb.TagNumber(505379326)
  void clearBlockpublicpolicy() => $_clearField(505379326);
}

class PutResourcePolicyResponse extends $pb.GeneratedMessage {
  factory PutResourcePolicyResponse({
    $core.String? name,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    return result;
  }

  PutResourcePolicyResponse._();

  factory PutResourcePolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutResourcePolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutResourcePolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutResourcePolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutResourcePolicyResponse copyWith(
          void Function(PutResourcePolicyResponse) updates) =>
      super.copyWith((message) => updates(message as PutResourcePolicyResponse))
          as PutResourcePolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutResourcePolicyResponse create() => PutResourcePolicyResponse._();
  @$core.override
  PutResourcePolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutResourcePolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutResourcePolicyResponse>(create);
  static PutResourcePolicyResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class PutSecretValueRequest extends $pb.GeneratedMessage {
  factory PutSecretValueRequest({
    $core.List<$core.int>? secretbinary,
    $core.String? secretstring,
    $core.Iterable<$core.String>? versionstages,
    $core.String? rotationtoken,
    $core.String? secretid,
    $core.String? clientrequesttoken,
  }) {
    final result = create();
    if (secretbinary != null) result.secretbinary = secretbinary;
    if (secretstring != null) result.secretstring = secretstring;
    if (versionstages != null) result.versionstages.addAll(versionstages);
    if (rotationtoken != null) result.rotationtoken = rotationtoken;
    if (secretid != null) result.secretid = secretid;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    return result;
  }

  PutSecretValueRequest._();

  factory PutSecretValueRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutSecretValueRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutSecretValueRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(
        94375681, _omitFieldNames ? '' : 'secretbinary', $pb.PbFieldType.OY)
    ..aOS(190782253, _omitFieldNames ? '' : 'secretstring')
    ..pPS(224220993, _omitFieldNames ? '' : 'versionstages')
    ..aOS(292175477, _omitFieldNames ? '' : 'rotationtoken')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutSecretValueRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutSecretValueRequest copyWith(
          void Function(PutSecretValueRequest) updates) =>
      super.copyWith((message) => updates(message as PutSecretValueRequest))
          as PutSecretValueRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutSecretValueRequest create() => PutSecretValueRequest._();
  @$core.override
  PutSecretValueRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutSecretValueRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutSecretValueRequest>(create);
  static PutSecretValueRequest? _defaultInstance;

  @$pb.TagNumber(94375681)
  $core.List<$core.int> get secretbinary => $_getN(0);
  @$pb.TagNumber(94375681)
  set secretbinary($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(94375681)
  $core.bool hasSecretbinary() => $_has(0);
  @$pb.TagNumber(94375681)
  void clearSecretbinary() => $_clearField(94375681);

  @$pb.TagNumber(190782253)
  $core.String get secretstring => $_getSZ(1);
  @$pb.TagNumber(190782253)
  set secretstring($core.String value) => $_setString(1, value);
  @$pb.TagNumber(190782253)
  $core.bool hasSecretstring() => $_has(1);
  @$pb.TagNumber(190782253)
  void clearSecretstring() => $_clearField(190782253);

  @$pb.TagNumber(224220993)
  $pb.PbList<$core.String> get versionstages => $_getList(2);

  @$pb.TagNumber(292175477)
  $core.String get rotationtoken => $_getSZ(3);
  @$pb.TagNumber(292175477)
  set rotationtoken($core.String value) => $_setString(3, value);
  @$pb.TagNumber(292175477)
  $core.bool hasRotationtoken() => $_has(3);
  @$pb.TagNumber(292175477)
  void clearRotationtoken() => $_clearField(292175477);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(4);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(4);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(5);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(5, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(5);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);
}

class PutSecretValueResponse extends $pb.GeneratedMessage {
  factory PutSecretValueResponse({
    $core.Iterable<$core.String>? versionstages,
    $core.String? name,
    $core.String? versionid,
    $core.String? arn,
  }) {
    final result = create();
    if (versionstages != null) result.versionstages.addAll(versionstages);
    if (name != null) result.name = name;
    if (versionid != null) result.versionid = versionid;
    if (arn != null) result.arn = arn;
    return result;
  }

  PutSecretValueResponse._();

  factory PutSecretValueResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutSecretValueResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutSecretValueResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..pPS(224220993, _omitFieldNames ? '' : 'versionstages')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutSecretValueResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutSecretValueResponse copyWith(
          void Function(PutSecretValueResponse) updates) =>
      super.copyWith((message) => updates(message as PutSecretValueResponse))
          as PutSecretValueResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutSecretValueResponse create() => PutSecretValueResponse._();
  @$core.override
  PutSecretValueResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutSecretValueResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutSecretValueResponse>(create);
  static PutSecretValueResponse? _defaultInstance;

  @$pb.TagNumber(224220993)
  $pb.PbList<$core.String> get versionstages => $_getList(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(2);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(2);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(3);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(3);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class RemoveRegionsFromReplicationRequest extends $pb.GeneratedMessage {
  factory RemoveRegionsFromReplicationRequest({
    $core.String? secretid,
    $core.Iterable<$core.String>? removereplicaregions,
  }) {
    final result = create();
    if (secretid != null) result.secretid = secretid;
    if (removereplicaregions != null)
      result.removereplicaregions.addAll(removereplicaregions);
    return result;
  }

  RemoveRegionsFromReplicationRequest._();

  factory RemoveRegionsFromReplicationRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveRegionsFromReplicationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveRegionsFromReplicationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..pPS(367753789, _omitFieldNames ? '' : 'removereplicaregions')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveRegionsFromReplicationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveRegionsFromReplicationRequest copyWith(
          void Function(RemoveRegionsFromReplicationRequest) updates) =>
      super.copyWith((message) =>
              updates(message as RemoveRegionsFromReplicationRequest))
          as RemoveRegionsFromReplicationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveRegionsFromReplicationRequest create() =>
      RemoveRegionsFromReplicationRequest._();
  @$core.override
  RemoveRegionsFromReplicationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveRegionsFromReplicationRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          RemoveRegionsFromReplicationRequest>(create);
  static RemoveRegionsFromReplicationRequest? _defaultInstance;

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(0);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(0);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);

  @$pb.TagNumber(367753789)
  $pb.PbList<$core.String> get removereplicaregions => $_getList(1);
}

class RemoveRegionsFromReplicationResponse extends $pb.GeneratedMessage {
  factory RemoveRegionsFromReplicationResponse({
    $core.String? arn,
    $core.Iterable<ReplicationStatusType>? replicationstatus,
  }) {
    final result = create();
    if (arn != null) result.arn = arn;
    if (replicationstatus != null)
      result.replicationstatus.addAll(replicationstatus);
    return result;
  }

  RemoveRegionsFromReplicationResponse._();

  factory RemoveRegionsFromReplicationResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveRegionsFromReplicationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveRegionsFromReplicationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..pPM<ReplicationStatusType>(
        529093900, _omitFieldNames ? '' : 'replicationstatus',
        subBuilder: ReplicationStatusType.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveRegionsFromReplicationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveRegionsFromReplicationResponse copyWith(
          void Function(RemoveRegionsFromReplicationResponse) updates) =>
      super.copyWith((message) =>
              updates(message as RemoveRegionsFromReplicationResponse))
          as RemoveRegionsFromReplicationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveRegionsFromReplicationResponse create() =>
      RemoveRegionsFromReplicationResponse._();
  @$core.override
  RemoveRegionsFromReplicationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveRegionsFromReplicationResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          RemoveRegionsFromReplicationResponse>(create);
  static RemoveRegionsFromReplicationResponse? _defaultInstance;

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(0);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(0);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(529093900)
  $pb.PbList<ReplicationStatusType> get replicationstatus => $_getList(1);
}

class ReplicaRegionType extends $pb.GeneratedMessage {
  factory ReplicaRegionType({
    $core.String? kmskeyid,
    $core.String? region,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (region != null) result.region = region;
    return result;
  }

  ReplicaRegionType._();

  factory ReplicaRegionType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaRegionType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaRegionType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(154040478, _omitFieldNames ? '' : 'region')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaRegionType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaRegionType copyWith(void Function(ReplicaRegionType) updates) =>
      super.copyWith((message) => updates(message as ReplicaRegionType))
          as ReplicaRegionType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaRegionType create() => ReplicaRegionType._();
  @$core.override
  ReplicaRegionType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaRegionType getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaRegionType>(create);
  static ReplicaRegionType? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(154040478)
  $core.String get region => $_getSZ(1);
  @$pb.TagNumber(154040478)
  set region($core.String value) => $_setString(1, value);
  @$pb.TagNumber(154040478)
  $core.bool hasRegion() => $_has(1);
  @$pb.TagNumber(154040478)
  void clearRegion() => $_clearField(154040478);
}

class ReplicateSecretToRegionsRequest extends $pb.GeneratedMessage {
  factory ReplicateSecretToRegionsRequest({
    $core.bool? forceoverwritereplicasecret,
    $core.String? secretid,
    $core.Iterable<ReplicaRegionType>? addreplicaregions,
  }) {
    final result = create();
    if (forceoverwritereplicasecret != null)
      result.forceoverwritereplicasecret = forceoverwritereplicasecret;
    if (secretid != null) result.secretid = secretid;
    if (addreplicaregions != null)
      result.addreplicaregions.addAll(addreplicaregions);
    return result;
  }

  ReplicateSecretToRegionsRequest._();

  factory ReplicateSecretToRegionsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicateSecretToRegionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicateSecretToRegionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOB(247407324, _omitFieldNames ? '' : 'forceoverwritereplicasecret')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..pPM<ReplicaRegionType>(
        461171870, _omitFieldNames ? '' : 'addreplicaregions',
        subBuilder: ReplicaRegionType.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicateSecretToRegionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicateSecretToRegionsRequest copyWith(
          void Function(ReplicateSecretToRegionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ReplicateSecretToRegionsRequest))
          as ReplicateSecretToRegionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicateSecretToRegionsRequest create() =>
      ReplicateSecretToRegionsRequest._();
  @$core.override
  ReplicateSecretToRegionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicateSecretToRegionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicateSecretToRegionsRequest>(
          create);
  static ReplicateSecretToRegionsRequest? _defaultInstance;

  @$pb.TagNumber(247407324)
  $core.bool get forceoverwritereplicasecret => $_getBF(0);
  @$pb.TagNumber(247407324)
  set forceoverwritereplicasecret($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(247407324)
  $core.bool hasForceoverwritereplicasecret() => $_has(0);
  @$pb.TagNumber(247407324)
  void clearForceoverwritereplicasecret() => $_clearField(247407324);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(1);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(1);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);

  @$pb.TagNumber(461171870)
  $pb.PbList<ReplicaRegionType> get addreplicaregions => $_getList(2);
}

class ReplicateSecretToRegionsResponse extends $pb.GeneratedMessage {
  factory ReplicateSecretToRegionsResponse({
    $core.String? arn,
    $core.Iterable<ReplicationStatusType>? replicationstatus,
  }) {
    final result = create();
    if (arn != null) result.arn = arn;
    if (replicationstatus != null)
      result.replicationstatus.addAll(replicationstatus);
    return result;
  }

  ReplicateSecretToRegionsResponse._();

  factory ReplicateSecretToRegionsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicateSecretToRegionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicateSecretToRegionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..pPM<ReplicationStatusType>(
        529093900, _omitFieldNames ? '' : 'replicationstatus',
        subBuilder: ReplicationStatusType.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicateSecretToRegionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicateSecretToRegionsResponse copyWith(
          void Function(ReplicateSecretToRegionsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ReplicateSecretToRegionsResponse))
          as ReplicateSecretToRegionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicateSecretToRegionsResponse create() =>
      ReplicateSecretToRegionsResponse._();
  @$core.override
  ReplicateSecretToRegionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicateSecretToRegionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicateSecretToRegionsResponse>(
          create);
  static ReplicateSecretToRegionsResponse? _defaultInstance;

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(0);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(0);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(529093900)
  $pb.PbList<ReplicationStatusType> get replicationstatus => $_getList(1);
}

class ReplicationStatusType extends $pb.GeneratedMessage {
  factory ReplicationStatusType({
    StatusType? status,
    $core.String? kmskeyid,
    $core.String? statusmessage,
    $core.String? region,
    $core.String? lastaccesseddate,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (statusmessage != null) result.statusmessage = statusmessage;
    if (region != null) result.region = region;
    if (lastaccesseddate != null) result.lastaccesseddate = lastaccesseddate;
    return result;
  }

  ReplicationStatusType._();

  factory ReplicationStatusType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicationStatusType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicationStatusType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aE<StatusType>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: StatusType.values)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..aOS(72590095, _omitFieldNames ? '' : 'statusmessage')
    ..aOS(154040478, _omitFieldNames ? '' : 'region')
    ..aOS(194418963, _omitFieldNames ? '' : 'lastaccesseddate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicationStatusType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicationStatusType copyWith(
          void Function(ReplicationStatusType) updates) =>
      super.copyWith((message) => updates(message as ReplicationStatusType))
          as ReplicationStatusType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicationStatusType create() => ReplicationStatusType._();
  @$core.override
  ReplicationStatusType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicationStatusType getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicationStatusType>(create);
  static ReplicationStatusType? _defaultInstance;

  @$pb.TagNumber(6222352)
  StatusType get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(StatusType value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(1);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(1);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(72590095)
  $core.String get statusmessage => $_getSZ(2);
  @$pb.TagNumber(72590095)
  set statusmessage($core.String value) => $_setString(2, value);
  @$pb.TagNumber(72590095)
  $core.bool hasStatusmessage() => $_has(2);
  @$pb.TagNumber(72590095)
  void clearStatusmessage() => $_clearField(72590095);

  @$pb.TagNumber(154040478)
  $core.String get region => $_getSZ(3);
  @$pb.TagNumber(154040478)
  set region($core.String value) => $_setString(3, value);
  @$pb.TagNumber(154040478)
  $core.bool hasRegion() => $_has(3);
  @$pb.TagNumber(154040478)
  void clearRegion() => $_clearField(154040478);

  @$pb.TagNumber(194418963)
  $core.String get lastaccesseddate => $_getSZ(4);
  @$pb.TagNumber(194418963)
  set lastaccesseddate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(194418963)
  $core.bool hasLastaccesseddate() => $_has(4);
  @$pb.TagNumber(194418963)
  void clearLastaccesseddate() => $_clearField(194418963);
}

class ResourceExistsException extends $pb.GeneratedMessage {
  factory ResourceExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourceExistsException._();

  factory ResourceExistsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceExistsException copyWith(
          void Function(ResourceExistsException) updates) =>
      super.copyWith((message) => updates(message as ResourceExistsException))
          as ResourceExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceExistsException create() => ResourceExistsException._();
  @$core.override
  ResourceExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceExistsException>(create);
  static ResourceExistsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
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

class RestoreSecretRequest extends $pb.GeneratedMessage {
  factory RestoreSecretRequest({
    $core.String? secretid,
  }) {
    final result = create();
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  RestoreSecretRequest._();

  factory RestoreSecretRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreSecretRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreSecretRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreSecretRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreSecretRequest copyWith(void Function(RestoreSecretRequest) updates) =>
      super.copyWith((message) => updates(message as RestoreSecretRequest))
          as RestoreSecretRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreSecretRequest create() => RestoreSecretRequest._();
  @$core.override
  RestoreSecretRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreSecretRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreSecretRequest>(create);
  static RestoreSecretRequest? _defaultInstance;

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(0);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(0);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class RestoreSecretResponse extends $pb.GeneratedMessage {
  factory RestoreSecretResponse({
    $core.String? name,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    return result;
  }

  RestoreSecretResponse._();

  factory RestoreSecretResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreSecretResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreSecretResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreSecretResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreSecretResponse copyWith(
          void Function(RestoreSecretResponse) updates) =>
      super.copyWith((message) => updates(message as RestoreSecretResponse))
          as RestoreSecretResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreSecretResponse create() => RestoreSecretResponse._();
  @$core.override
  RestoreSecretResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreSecretResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreSecretResponse>(create);
  static RestoreSecretResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class RotateSecretRequest extends $pb.GeneratedMessage {
  factory RotateSecretRequest({
    $core.Iterable<ExternalSecretRotationMetadataItem>?
        externalsecretrotationmetadata,
    RotationRulesType? rotationrules,
    $core.bool? rotateimmediately,
    $core.String? rotationlambdaarn,
    $core.String? secretid,
    $core.String? clientrequesttoken,
    $core.String? externalsecretrotationrolearn,
  }) {
    final result = create();
    if (externalsecretrotationmetadata != null)
      result.externalsecretrotationmetadata
          .addAll(externalsecretrotationmetadata);
    if (rotationrules != null) result.rotationrules = rotationrules;
    if (rotateimmediately != null) result.rotateimmediately = rotateimmediately;
    if (rotationlambdaarn != null) result.rotationlambdaarn = rotationlambdaarn;
    if (secretid != null) result.secretid = secretid;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (externalsecretrotationrolearn != null)
      result.externalsecretrotationrolearn = externalsecretrotationrolearn;
    return result;
  }

  RotateSecretRequest._();

  factory RotateSecretRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RotateSecretRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RotateSecretRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..pPM<ExternalSecretRotationMetadataItem>(
        52900542, _omitFieldNames ? '' : 'externalsecretrotationmetadata',
        subBuilder: ExternalSecretRotationMetadataItem.create)
    ..aOM<RotationRulesType>(259458135, _omitFieldNames ? '' : 'rotationrules',
        subBuilder: RotationRulesType.create)
    ..aOB(265384053, _omitFieldNames ? '' : 'rotateimmediately')
    ..aOS(335026080, _omitFieldNames ? '' : 'rotationlambdaarn')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..aOS(470712576, _omitFieldNames ? '' : 'externalsecretrotationrolearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotateSecretRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotateSecretRequest copyWith(void Function(RotateSecretRequest) updates) =>
      super.copyWith((message) => updates(message as RotateSecretRequest))
          as RotateSecretRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RotateSecretRequest create() => RotateSecretRequest._();
  @$core.override
  RotateSecretRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RotateSecretRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RotateSecretRequest>(create);
  static RotateSecretRequest? _defaultInstance;

  @$pb.TagNumber(52900542)
  $pb.PbList<ExternalSecretRotationMetadataItem>
      get externalsecretrotationmetadata => $_getList(0);

  @$pb.TagNumber(259458135)
  RotationRulesType get rotationrules => $_getN(1);
  @$pb.TagNumber(259458135)
  set rotationrules(RotationRulesType value) => $_setField(259458135, value);
  @$pb.TagNumber(259458135)
  $core.bool hasRotationrules() => $_has(1);
  @$pb.TagNumber(259458135)
  void clearRotationrules() => $_clearField(259458135);
  @$pb.TagNumber(259458135)
  RotationRulesType ensureRotationrules() => $_ensure(1);

  @$pb.TagNumber(265384053)
  $core.bool get rotateimmediately => $_getBF(2);
  @$pb.TagNumber(265384053)
  set rotateimmediately($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(265384053)
  $core.bool hasRotateimmediately() => $_has(2);
  @$pb.TagNumber(265384053)
  void clearRotateimmediately() => $_clearField(265384053);

  @$pb.TagNumber(335026080)
  $core.String get rotationlambdaarn => $_getSZ(3);
  @$pb.TagNumber(335026080)
  set rotationlambdaarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(335026080)
  $core.bool hasRotationlambdaarn() => $_has(3);
  @$pb.TagNumber(335026080)
  void clearRotationlambdaarn() => $_clearField(335026080);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(4);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(4);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(5);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(5, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(5);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(470712576)
  $core.String get externalsecretrotationrolearn => $_getSZ(6);
  @$pb.TagNumber(470712576)
  set externalsecretrotationrolearn($core.String value) =>
      $_setString(6, value);
  @$pb.TagNumber(470712576)
  $core.bool hasExternalsecretrotationrolearn() => $_has(6);
  @$pb.TagNumber(470712576)
  void clearExternalsecretrotationrolearn() => $_clearField(470712576);
}

class RotateSecretResponse extends $pb.GeneratedMessage {
  factory RotateSecretResponse({
    $core.String? name,
    $core.String? versionid,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (versionid != null) result.versionid = versionid;
    if (arn != null) result.arn = arn;
    return result;
  }

  RotateSecretResponse._();

  factory RotateSecretResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RotateSecretResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RotateSecretResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotateSecretResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotateSecretResponse copyWith(void Function(RotateSecretResponse) updates) =>
      super.copyWith((message) => updates(message as RotateSecretResponse))
          as RotateSecretResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RotateSecretResponse create() => RotateSecretResponse._();
  @$core.override
  RotateSecretResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RotateSecretResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RotateSecretResponse>(create);
  static RotateSecretResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(1);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(1);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class RotationRulesType extends $pb.GeneratedMessage {
  factory RotationRulesType({
    $core.String? duration,
    $fixnum.Int64? automaticallyafterdays,
    $core.String? scheduleexpression,
  }) {
    final result = create();
    if (duration != null) result.duration = duration;
    if (automaticallyafterdays != null)
      result.automaticallyafterdays = automaticallyafterdays;
    if (scheduleexpression != null)
      result.scheduleexpression = scheduleexpression;
    return result;
  }

  RotationRulesType._();

  factory RotationRulesType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RotationRulesType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RotationRulesType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(348604718, _omitFieldNames ? '' : 'duration')
    ..aInt64(350893940, _omitFieldNames ? '' : 'automaticallyafterdays')
    ..aOS(446089471, _omitFieldNames ? '' : 'scheduleexpression')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotationRulesType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RotationRulesType copyWith(void Function(RotationRulesType) updates) =>
      super.copyWith((message) => updates(message as RotationRulesType))
          as RotationRulesType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RotationRulesType create() => RotationRulesType._();
  @$core.override
  RotationRulesType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RotationRulesType getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RotationRulesType>(create);
  static RotationRulesType? _defaultInstance;

  @$pb.TagNumber(348604718)
  $core.String get duration => $_getSZ(0);
  @$pb.TagNumber(348604718)
  set duration($core.String value) => $_setString(0, value);
  @$pb.TagNumber(348604718)
  $core.bool hasDuration() => $_has(0);
  @$pb.TagNumber(348604718)
  void clearDuration() => $_clearField(348604718);

  @$pb.TagNumber(350893940)
  $fixnum.Int64 get automaticallyafterdays => $_getI64(1);
  @$pb.TagNumber(350893940)
  set automaticallyafterdays($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(350893940)
  $core.bool hasAutomaticallyafterdays() => $_has(1);
  @$pb.TagNumber(350893940)
  void clearAutomaticallyafterdays() => $_clearField(350893940);

  @$pb.TagNumber(446089471)
  $core.String get scheduleexpression => $_getSZ(2);
  @$pb.TagNumber(446089471)
  set scheduleexpression($core.String value) => $_setString(2, value);
  @$pb.TagNumber(446089471)
  $core.bool hasScheduleexpression() => $_has(2);
  @$pb.TagNumber(446089471)
  void clearScheduleexpression() => $_clearField(446089471);
}

class SecretListEntry extends $pb.GeneratedMessage {
  factory SecretListEntry({
    $core.String? kmskeyid,
    $core.Iterable<ExternalSecretRotationMetadataItem>?
        externalsecretrotationmetadata,
    $core.String? description,
    $core.String? nextrotationdate,
    $core.String? lastaccesseddate,
    $core.bool? rotationenabled,
    RotationRulesType? rotationrules,
    $core.String? name,
    $core.String? type,
    $core.String? lastchangeddate,
    $core.String? rotationlambdaarn,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        secretversionstostages,
    $core.Iterable<Tag>? tags,
    $core.String? arn,
    $core.String? createddate,
    $core.String? owningservice,
    $core.String? externalsecretrotationrolearn,
    $core.String? primaryregion,
    $core.String? lastrotateddate,
    $core.String? deleteddate,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (externalsecretrotationmetadata != null)
      result.externalsecretrotationmetadata
          .addAll(externalsecretrotationmetadata);
    if (description != null) result.description = description;
    if (nextrotationdate != null) result.nextrotationdate = nextrotationdate;
    if (lastaccesseddate != null) result.lastaccesseddate = lastaccesseddate;
    if (rotationenabled != null) result.rotationenabled = rotationenabled;
    if (rotationrules != null) result.rotationrules = rotationrules;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (lastchangeddate != null) result.lastchangeddate = lastchangeddate;
    if (rotationlambdaarn != null) result.rotationlambdaarn = rotationlambdaarn;
    if (secretversionstostages != null)
      result.secretversionstostages.addEntries(secretversionstostages);
    if (tags != null) result.tags.addAll(tags);
    if (arn != null) result.arn = arn;
    if (createddate != null) result.createddate = createddate;
    if (owningservice != null) result.owningservice = owningservice;
    if (externalsecretrotationrolearn != null)
      result.externalsecretrotationrolearn = externalsecretrotationrolearn;
    if (primaryregion != null) result.primaryregion = primaryregion;
    if (lastrotateddate != null) result.lastrotateddate = lastrotateddate;
    if (deleteddate != null) result.deleteddate = deleteddate;
    return result;
  }

  SecretListEntry._();

  factory SecretListEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SecretListEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SecretListEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..pPM<ExternalSecretRotationMetadataItem>(
        52900542, _omitFieldNames ? '' : 'externalsecretrotationmetadata',
        subBuilder: ExternalSecretRotationMetadataItem.create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(192035355, _omitFieldNames ? '' : 'nextrotationdate')
    ..aOS(194418963, _omitFieldNames ? '' : 'lastaccesseddate')
    ..aOB(209507301, _omitFieldNames ? '' : 'rotationenabled')
    ..aOM<RotationRulesType>(259458135, _omitFieldNames ? '' : 'rotationrules',
        subBuilder: RotationRulesType.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
    ..aOS(314015460, _omitFieldNames ? '' : 'lastchangeddate')
    ..aOS(335026080, _omitFieldNames ? '' : 'rotationlambdaarn')
    ..m<$core.String, $core.String>(
        356331823, _omitFieldNames ? '' : 'secretversionstostages',
        entryClassName: 'SecretListEntry.SecretversionstostagesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('secretsmanager'))
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..aOS(416929840, _omitFieldNames ? '' : 'createddate')
    ..aOS(462487817, _omitFieldNames ? '' : 'owningservice')
    ..aOS(470712576, _omitFieldNames ? '' : 'externalsecretrotationrolearn')
    ..aOS(480901186, _omitFieldNames ? '' : 'primaryregion')
    ..aOS(501475691, _omitFieldNames ? '' : 'lastrotateddate')
    ..aOS(516314255, _omitFieldNames ? '' : 'deleteddate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SecretListEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SecretListEntry copyWith(void Function(SecretListEntry) updates) =>
      super.copyWith((message) => updates(message as SecretListEntry))
          as SecretListEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SecretListEntry create() => SecretListEntry._();
  @$core.override
  SecretListEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SecretListEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SecretListEntry>(create);
  static SecretListEntry? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(52900542)
  $pb.PbList<ExternalSecretRotationMetadataItem>
      get externalsecretrotationmetadata => $_getList(1);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(192035355)
  $core.String get nextrotationdate => $_getSZ(3);
  @$pb.TagNumber(192035355)
  set nextrotationdate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(192035355)
  $core.bool hasNextrotationdate() => $_has(3);
  @$pb.TagNumber(192035355)
  void clearNextrotationdate() => $_clearField(192035355);

  @$pb.TagNumber(194418963)
  $core.String get lastaccesseddate => $_getSZ(4);
  @$pb.TagNumber(194418963)
  set lastaccesseddate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(194418963)
  $core.bool hasLastaccesseddate() => $_has(4);
  @$pb.TagNumber(194418963)
  void clearLastaccesseddate() => $_clearField(194418963);

  @$pb.TagNumber(209507301)
  $core.bool get rotationenabled => $_getBF(5);
  @$pb.TagNumber(209507301)
  set rotationenabled($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(209507301)
  $core.bool hasRotationenabled() => $_has(5);
  @$pb.TagNumber(209507301)
  void clearRotationenabled() => $_clearField(209507301);

  @$pb.TagNumber(259458135)
  RotationRulesType get rotationrules => $_getN(6);
  @$pb.TagNumber(259458135)
  set rotationrules(RotationRulesType value) => $_setField(259458135, value);
  @$pb.TagNumber(259458135)
  $core.bool hasRotationrules() => $_has(6);
  @$pb.TagNumber(259458135)
  void clearRotationrules() => $_clearField(259458135);
  @$pb.TagNumber(259458135)
  RotationRulesType ensureRotationrules() => $_ensure(6);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(7);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(7, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(7);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  $core.String get type => $_getSZ(8);
  @$pb.TagNumber(290836590)
  set type($core.String value) => $_setString(8, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(8);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(314015460)
  $core.String get lastchangeddate => $_getSZ(9);
  @$pb.TagNumber(314015460)
  set lastchangeddate($core.String value) => $_setString(9, value);
  @$pb.TagNumber(314015460)
  $core.bool hasLastchangeddate() => $_has(9);
  @$pb.TagNumber(314015460)
  void clearLastchangeddate() => $_clearField(314015460);

  @$pb.TagNumber(335026080)
  $core.String get rotationlambdaarn => $_getSZ(10);
  @$pb.TagNumber(335026080)
  set rotationlambdaarn($core.String value) => $_setString(10, value);
  @$pb.TagNumber(335026080)
  $core.bool hasRotationlambdaarn() => $_has(10);
  @$pb.TagNumber(335026080)
  void clearRotationlambdaarn() => $_clearField(335026080);

  @$pb.TagNumber(356331823)
  $pb.PbMap<$core.String, $core.String> get secretversionstostages =>
      $_getMap(11);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(12);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(13);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(13, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(13);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(416929840)
  $core.String get createddate => $_getSZ(14);
  @$pb.TagNumber(416929840)
  set createddate($core.String value) => $_setString(14, value);
  @$pb.TagNumber(416929840)
  $core.bool hasCreateddate() => $_has(14);
  @$pb.TagNumber(416929840)
  void clearCreateddate() => $_clearField(416929840);

  @$pb.TagNumber(462487817)
  $core.String get owningservice => $_getSZ(15);
  @$pb.TagNumber(462487817)
  set owningservice($core.String value) => $_setString(15, value);
  @$pb.TagNumber(462487817)
  $core.bool hasOwningservice() => $_has(15);
  @$pb.TagNumber(462487817)
  void clearOwningservice() => $_clearField(462487817);

  @$pb.TagNumber(470712576)
  $core.String get externalsecretrotationrolearn => $_getSZ(16);
  @$pb.TagNumber(470712576)
  set externalsecretrotationrolearn($core.String value) =>
      $_setString(16, value);
  @$pb.TagNumber(470712576)
  $core.bool hasExternalsecretrotationrolearn() => $_has(16);
  @$pb.TagNumber(470712576)
  void clearExternalsecretrotationrolearn() => $_clearField(470712576);

  @$pb.TagNumber(480901186)
  $core.String get primaryregion => $_getSZ(17);
  @$pb.TagNumber(480901186)
  set primaryregion($core.String value) => $_setString(17, value);
  @$pb.TagNumber(480901186)
  $core.bool hasPrimaryregion() => $_has(17);
  @$pb.TagNumber(480901186)
  void clearPrimaryregion() => $_clearField(480901186);

  @$pb.TagNumber(501475691)
  $core.String get lastrotateddate => $_getSZ(18);
  @$pb.TagNumber(501475691)
  set lastrotateddate($core.String value) => $_setString(18, value);
  @$pb.TagNumber(501475691)
  $core.bool hasLastrotateddate() => $_has(18);
  @$pb.TagNumber(501475691)
  void clearLastrotateddate() => $_clearField(501475691);

  @$pb.TagNumber(516314255)
  $core.String get deleteddate => $_getSZ(19);
  @$pb.TagNumber(516314255)
  set deleteddate($core.String value) => $_setString(19, value);
  @$pb.TagNumber(516314255)
  $core.bool hasDeleteddate() => $_has(19);
  @$pb.TagNumber(516314255)
  void clearDeleteddate() => $_clearField(516314255);
}

class SecretValueEntry extends $pb.GeneratedMessage {
  factory SecretValueEntry({
    $core.List<$core.int>? secretbinary,
    $core.String? secretstring,
    $core.Iterable<$core.String>? versionstages,
    $core.String? name,
    $core.String? versionid,
    $core.String? arn,
    $core.String? createddate,
  }) {
    final result = create();
    if (secretbinary != null) result.secretbinary = secretbinary;
    if (secretstring != null) result.secretstring = secretstring;
    if (versionstages != null) result.versionstages.addAll(versionstages);
    if (name != null) result.name = name;
    if (versionid != null) result.versionid = versionid;
    if (arn != null) result.arn = arn;
    if (createddate != null) result.createddate = createddate;
    return result;
  }

  SecretValueEntry._();

  factory SecretValueEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SecretValueEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SecretValueEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(
        94375681, _omitFieldNames ? '' : 'secretbinary', $pb.PbFieldType.OY)
    ..aOS(190782253, _omitFieldNames ? '' : 'secretstring')
    ..pPS(224220993, _omitFieldNames ? '' : 'versionstages')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..aOS(416929840, _omitFieldNames ? '' : 'createddate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SecretValueEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SecretValueEntry copyWith(void Function(SecretValueEntry) updates) =>
      super.copyWith((message) => updates(message as SecretValueEntry))
          as SecretValueEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SecretValueEntry create() => SecretValueEntry._();
  @$core.override
  SecretValueEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SecretValueEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SecretValueEntry>(create);
  static SecretValueEntry? _defaultInstance;

  @$pb.TagNumber(94375681)
  $core.List<$core.int> get secretbinary => $_getN(0);
  @$pb.TagNumber(94375681)
  set secretbinary($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(94375681)
  $core.bool hasSecretbinary() => $_has(0);
  @$pb.TagNumber(94375681)
  void clearSecretbinary() => $_clearField(94375681);

  @$pb.TagNumber(190782253)
  $core.String get secretstring => $_getSZ(1);
  @$pb.TagNumber(190782253)
  set secretstring($core.String value) => $_setString(1, value);
  @$pb.TagNumber(190782253)
  $core.bool hasSecretstring() => $_has(1);
  @$pb.TagNumber(190782253)
  void clearSecretstring() => $_clearField(190782253);

  @$pb.TagNumber(224220993)
  $pb.PbList<$core.String> get versionstages => $_getList(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(4);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(4);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(5);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(5);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(416929840)
  $core.String get createddate => $_getSZ(6);
  @$pb.TagNumber(416929840)
  set createddate($core.String value) => $_setString(6, value);
  @$pb.TagNumber(416929840)
  $core.bool hasCreateddate() => $_has(6);
  @$pb.TagNumber(416929840)
  void clearCreateddate() => $_clearField(416929840);
}

class SecretVersionsListEntry extends $pb.GeneratedMessage {
  factory SecretVersionsListEntry({
    $core.String? lastaccesseddate,
    $core.Iterable<$core.String>? versionstages,
    $core.String? versionid,
    $core.String? createddate,
    $core.Iterable<$core.String>? kmskeyids,
  }) {
    final result = create();
    if (lastaccesseddate != null) result.lastaccesseddate = lastaccesseddate;
    if (versionstages != null) result.versionstages.addAll(versionstages);
    if (versionid != null) result.versionid = versionid;
    if (createddate != null) result.createddate = createddate;
    if (kmskeyids != null) result.kmskeyids.addAll(kmskeyids);
    return result;
  }

  SecretVersionsListEntry._();

  factory SecretVersionsListEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SecretVersionsListEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SecretVersionsListEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(194418963, _omitFieldNames ? '' : 'lastaccesseddate')
    ..pPS(224220993, _omitFieldNames ? '' : 'versionstages')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(416929840, _omitFieldNames ? '' : 'createddate')
    ..pPS(478641518, _omitFieldNames ? '' : 'kmskeyids')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SecretVersionsListEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SecretVersionsListEntry copyWith(
          void Function(SecretVersionsListEntry) updates) =>
      super.copyWith((message) => updates(message as SecretVersionsListEntry))
          as SecretVersionsListEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SecretVersionsListEntry create() => SecretVersionsListEntry._();
  @$core.override
  SecretVersionsListEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SecretVersionsListEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SecretVersionsListEntry>(create);
  static SecretVersionsListEntry? _defaultInstance;

  @$pb.TagNumber(194418963)
  $core.String get lastaccesseddate => $_getSZ(0);
  @$pb.TagNumber(194418963)
  set lastaccesseddate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(194418963)
  $core.bool hasLastaccesseddate() => $_has(0);
  @$pb.TagNumber(194418963)
  void clearLastaccesseddate() => $_clearField(194418963);

  @$pb.TagNumber(224220993)
  $pb.PbList<$core.String> get versionstages => $_getList(1);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(2);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(2);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(416929840)
  $core.String get createddate => $_getSZ(3);
  @$pb.TagNumber(416929840)
  set createddate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(416929840)
  $core.bool hasCreateddate() => $_has(3);
  @$pb.TagNumber(416929840)
  void clearCreateddate() => $_clearField(416929840);

  @$pb.TagNumber(478641518)
  $pb.PbList<$core.String> get kmskeyids => $_getList(4);
}

class StopReplicationToReplicaRequest extends $pb.GeneratedMessage {
  factory StopReplicationToReplicaRequest({
    $core.String? secretid,
  }) {
    final result = create();
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  StopReplicationToReplicaRequest._();

  factory StopReplicationToReplicaRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopReplicationToReplicaRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopReplicationToReplicaRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopReplicationToReplicaRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopReplicationToReplicaRequest copyWith(
          void Function(StopReplicationToReplicaRequest) updates) =>
      super.copyWith(
              (message) => updates(message as StopReplicationToReplicaRequest))
          as StopReplicationToReplicaRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopReplicationToReplicaRequest create() =>
      StopReplicationToReplicaRequest._();
  @$core.override
  StopReplicationToReplicaRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopReplicationToReplicaRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopReplicationToReplicaRequest>(
          create);
  static StopReplicationToReplicaRequest? _defaultInstance;

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(0);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(0);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class StopReplicationToReplicaResponse extends $pb.GeneratedMessage {
  factory StopReplicationToReplicaResponse({
    $core.String? arn,
  }) {
    final result = create();
    if (arn != null) result.arn = arn;
    return result;
  }

  StopReplicationToReplicaResponse._();

  factory StopReplicationToReplicaResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopReplicationToReplicaResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopReplicationToReplicaResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopReplicationToReplicaResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopReplicationToReplicaResponse copyWith(
          void Function(StopReplicationToReplicaResponse) updates) =>
      super.copyWith(
              (message) => updates(message as StopReplicationToReplicaResponse))
          as StopReplicationToReplicaResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopReplicationToReplicaResponse create() =>
      StopReplicationToReplicaResponse._();
  @$core.override
  StopReplicationToReplicaResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopReplicationToReplicaResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopReplicationToReplicaResponse>(
          create);
  static StopReplicationToReplicaResponse? _defaultInstance;

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(0);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(0);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
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
    $core.String? secretid,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (secretid != null) result.secretid = secretid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
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

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(0);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(0);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);
}

class UntagResourceRequest extends $pb.GeneratedMessage {
  factory UntagResourceRequest({
    $core.Iterable<$core.String>? tagkeys,
    $core.String? secretid,
  }) {
    final result = create();
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
    if (secretid != null) result.secretid = secretid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
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

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(1);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(1);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class UpdateSecretRequest extends $pb.GeneratedMessage {
  factory UpdateSecretRequest({
    $core.String? kmskeyid,
    $core.List<$core.int>? secretbinary,
    $core.String? description,
    $core.String? secretstring,
    $core.String? type,
    $core.String? secretid,
    $core.String? clientrequesttoken,
  }) {
    final result = create();
    if (kmskeyid != null) result.kmskeyid = kmskeyid;
    if (secretbinary != null) result.secretbinary = secretbinary;
    if (description != null) result.description = description;
    if (secretstring != null) result.secretstring = secretstring;
    if (type != null) result.type = type;
    if (secretid != null) result.secretid = secretid;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    return result;
  }

  UpdateSecretRequest._();

  factory UpdateSecretRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateSecretRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateSecretRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(46523533, _omitFieldNames ? '' : 'kmskeyid')
    ..a<$core.List<$core.int>>(
        94375681, _omitFieldNames ? '' : 'secretbinary', $pb.PbFieldType.OY)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(190782253, _omitFieldNames ? '' : 'secretstring')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSecretRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSecretRequest copyWith(void Function(UpdateSecretRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateSecretRequest))
          as UpdateSecretRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateSecretRequest create() => UpdateSecretRequest._();
  @$core.override
  UpdateSecretRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateSecretRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateSecretRequest>(create);
  static UpdateSecretRequest? _defaultInstance;

  @$pb.TagNumber(46523533)
  $core.String get kmskeyid => $_getSZ(0);
  @$pb.TagNumber(46523533)
  set kmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(46523533)
  $core.bool hasKmskeyid() => $_has(0);
  @$pb.TagNumber(46523533)
  void clearKmskeyid() => $_clearField(46523533);

  @$pb.TagNumber(94375681)
  $core.List<$core.int> get secretbinary => $_getN(1);
  @$pb.TagNumber(94375681)
  set secretbinary($core.List<$core.int> value) => $_setBytes(1, value);
  @$pb.TagNumber(94375681)
  $core.bool hasSecretbinary() => $_has(1);
  @$pb.TagNumber(94375681)
  void clearSecretbinary() => $_clearField(94375681);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(190782253)
  $core.String get secretstring => $_getSZ(3);
  @$pb.TagNumber(190782253)
  set secretstring($core.String value) => $_setString(3, value);
  @$pb.TagNumber(190782253)
  $core.bool hasSecretstring() => $_has(3);
  @$pb.TagNumber(190782253)
  void clearSecretstring() => $_clearField(190782253);

  @$pb.TagNumber(290836590)
  $core.String get type => $_getSZ(4);
  @$pb.TagNumber(290836590)
  set type($core.String value) => $_setString(4, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(4);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(5);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(5);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(6);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(6, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(6);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);
}

class UpdateSecretResponse extends $pb.GeneratedMessage {
  factory UpdateSecretResponse({
    $core.String? name,
    $core.String? versionid,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (versionid != null) result.versionid = versionid;
    if (arn != null) result.arn = arn;
    return result;
  }

  UpdateSecretResponse._();

  factory UpdateSecretResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateSecretResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateSecretResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(338063515, _omitFieldNames ? '' : 'versionid')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSecretResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSecretResponse copyWith(void Function(UpdateSecretResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateSecretResponse))
          as UpdateSecretResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateSecretResponse create() => UpdateSecretResponse._();
  @$core.override
  UpdateSecretResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateSecretResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateSecretResponse>(create);
  static UpdateSecretResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(338063515)
  $core.String get versionid => $_getSZ(1);
  @$pb.TagNumber(338063515)
  set versionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(338063515)
  $core.bool hasVersionid() => $_has(1);
  @$pb.TagNumber(338063515)
  void clearVersionid() => $_clearField(338063515);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class UpdateSecretVersionStageRequest extends $pb.GeneratedMessage {
  factory UpdateSecretVersionStageRequest({
    $core.String? removefromversionid,
    $core.String? versionstage,
    $core.String? secretid,
    $core.String? movetoversionid,
  }) {
    final result = create();
    if (removefromversionid != null)
      result.removefromversionid = removefromversionid;
    if (versionstage != null) result.versionstage = versionstage;
    if (secretid != null) result.secretid = secretid;
    if (movetoversionid != null) result.movetoversionid = movetoversionid;
    return result;
  }

  UpdateSecretVersionStageRequest._();

  factory UpdateSecretVersionStageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateSecretVersionStageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateSecretVersionStageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(194704147, _omitFieldNames ? '' : 'removefromversionid')
    ..aOS(229692340, _omitFieldNames ? '' : 'versionstage')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..aOS(509017411, _omitFieldNames ? '' : 'movetoversionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSecretVersionStageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSecretVersionStageRequest copyWith(
          void Function(UpdateSecretVersionStageRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateSecretVersionStageRequest))
          as UpdateSecretVersionStageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateSecretVersionStageRequest create() =>
      UpdateSecretVersionStageRequest._();
  @$core.override
  UpdateSecretVersionStageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateSecretVersionStageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateSecretVersionStageRequest>(
          create);
  static UpdateSecretVersionStageRequest? _defaultInstance;

  @$pb.TagNumber(194704147)
  $core.String get removefromversionid => $_getSZ(0);
  @$pb.TagNumber(194704147)
  set removefromversionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(194704147)
  $core.bool hasRemovefromversionid() => $_has(0);
  @$pb.TagNumber(194704147)
  void clearRemovefromversionid() => $_clearField(194704147);

  @$pb.TagNumber(229692340)
  $core.String get versionstage => $_getSZ(1);
  @$pb.TagNumber(229692340)
  set versionstage($core.String value) => $_setString(1, value);
  @$pb.TagNumber(229692340)
  $core.bool hasVersionstage() => $_has(1);
  @$pb.TagNumber(229692340)
  void clearVersionstage() => $_clearField(229692340);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(2);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(2);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);

  @$pb.TagNumber(509017411)
  $core.String get movetoversionid => $_getSZ(3);
  @$pb.TagNumber(509017411)
  set movetoversionid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(509017411)
  $core.bool hasMovetoversionid() => $_has(3);
  @$pb.TagNumber(509017411)
  void clearMovetoversionid() => $_clearField(509017411);
}

class UpdateSecretVersionStageResponse extends $pb.GeneratedMessage {
  factory UpdateSecretVersionStageResponse({
    $core.String? name,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    return result;
  }

  UpdateSecretVersionStageResponse._();

  factory UpdateSecretVersionStageResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateSecretVersionStageResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateSecretVersionStageResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSecretVersionStageResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSecretVersionStageResponse copyWith(
          void Function(UpdateSecretVersionStageResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateSecretVersionStageResponse))
          as UpdateSecretVersionStageResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateSecretVersionStageResponse create() =>
      UpdateSecretVersionStageResponse._();
  @$core.override
  UpdateSecretVersionStageResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateSecretVersionStageResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateSecretVersionStageResponse>(
          create);
  static UpdateSecretVersionStageResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class ValidateResourcePolicyRequest extends $pb.GeneratedMessage {
  factory ValidateResourcePolicyRequest({
    $core.String? resourcepolicy,
    $core.String? secretid,
  }) {
    final result = create();
    if (resourcepolicy != null) result.resourcepolicy = resourcepolicy;
    if (secretid != null) result.secretid = secretid;
    return result;
  }

  ValidateResourcePolicyRequest._();

  factory ValidateResourcePolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ValidateResourcePolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ValidateResourcePolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(15747632, _omitFieldNames ? '' : 'resourcepolicy')
    ..aOS(341502821, _omitFieldNames ? '' : 'secretid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateResourcePolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateResourcePolicyRequest copyWith(
          void Function(ValidateResourcePolicyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ValidateResourcePolicyRequest))
          as ValidateResourcePolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidateResourcePolicyRequest create() =>
      ValidateResourcePolicyRequest._();
  @$core.override
  ValidateResourcePolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ValidateResourcePolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ValidateResourcePolicyRequest>(create);
  static ValidateResourcePolicyRequest? _defaultInstance;

  @$pb.TagNumber(15747632)
  $core.String get resourcepolicy => $_getSZ(0);
  @$pb.TagNumber(15747632)
  set resourcepolicy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(15747632)
  $core.bool hasResourcepolicy() => $_has(0);
  @$pb.TagNumber(15747632)
  void clearResourcepolicy() => $_clearField(15747632);

  @$pb.TagNumber(341502821)
  $core.String get secretid => $_getSZ(1);
  @$pb.TagNumber(341502821)
  set secretid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(341502821)
  $core.bool hasSecretid() => $_has(1);
  @$pb.TagNumber(341502821)
  void clearSecretid() => $_clearField(341502821);
}

class ValidateResourcePolicyResponse extends $pb.GeneratedMessage {
  factory ValidateResourcePolicyResponse({
    $core.bool? policyvalidationpassed,
    $core.Iterable<ValidationErrorsEntry>? validationerrors,
  }) {
    final result = create();
    if (policyvalidationpassed != null)
      result.policyvalidationpassed = policyvalidationpassed;
    if (validationerrors != null)
      result.validationerrors.addAll(validationerrors);
    return result;
  }

  ValidateResourcePolicyResponse._();

  factory ValidateResourcePolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ValidateResourcePolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ValidateResourcePolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOB(298018737, _omitFieldNames ? '' : 'policyvalidationpassed')
    ..pPM<ValidationErrorsEntry>(
        381330622, _omitFieldNames ? '' : 'validationerrors',
        subBuilder: ValidationErrorsEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateResourcePolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidateResourcePolicyResponse copyWith(
          void Function(ValidateResourcePolicyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ValidateResourcePolicyResponse))
          as ValidateResourcePolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidateResourcePolicyResponse create() =>
      ValidateResourcePolicyResponse._();
  @$core.override
  ValidateResourcePolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ValidateResourcePolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ValidateResourcePolicyResponse>(create);
  static ValidateResourcePolicyResponse? _defaultInstance;

  @$pb.TagNumber(298018737)
  $core.bool get policyvalidationpassed => $_getBF(0);
  @$pb.TagNumber(298018737)
  set policyvalidationpassed($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(298018737)
  $core.bool hasPolicyvalidationpassed() => $_has(0);
  @$pb.TagNumber(298018737)
  void clearPolicyvalidationpassed() => $_clearField(298018737);

  @$pb.TagNumber(381330622)
  $pb.PbList<ValidationErrorsEntry> get validationerrors => $_getList(1);
}

class ValidationErrorsEntry extends $pb.GeneratedMessage {
  factory ValidationErrorsEntry({
    $core.String? checkname,
    $core.String? errormessage,
  }) {
    final result = create();
    if (checkname != null) result.checkname = checkname;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  ValidationErrorsEntry._();

  factory ValidationErrorsEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ValidationErrorsEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ValidationErrorsEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'secretsmanager'),
      createEmptyInstance: create)
    ..aOS(68143813, _omitFieldNames ? '' : 'checkname')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidationErrorsEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidationErrorsEntry copyWith(
          void Function(ValidationErrorsEntry) updates) =>
      super.copyWith((message) => updates(message as ValidationErrorsEntry))
          as ValidationErrorsEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidationErrorsEntry create() => ValidationErrorsEntry._();
  @$core.override
  ValidationErrorsEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ValidationErrorsEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ValidationErrorsEntry>(create);
  static ValidationErrorsEntry? _defaultInstance;

  @$pb.TagNumber(68143813)
  $core.String get checkname => $_getSZ(0);
  @$pb.TagNumber(68143813)
  set checkname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(68143813)
  $core.bool hasCheckname() => $_has(0);
  @$pb.TagNumber(68143813)
  void clearCheckname() => $_clearField(68143813);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(1);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(1, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(1);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
