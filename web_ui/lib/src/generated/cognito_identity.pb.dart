// This is a generated file - do not edit.
//
// Generated from cognito_identity.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'cognito_identity.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'cognito_identity.pbenum.dart';

class CognitoIdentityProvider extends $pb.GeneratedMessage {
  factory CognitoIdentityProvider({
    $core.bool? serversidetokencheck,
    $core.String? clientid,
    $core.String? providername,
  }) {
    final result = create();
    if (serversidetokencheck != null)
      result.serversidetokencheck = serversidetokencheck;
    if (clientid != null) result.clientid = clientid;
    if (providername != null) result.providername = providername;
    return result;
  }

  CognitoIdentityProvider._();

  factory CognitoIdentityProvider.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CognitoIdentityProvider.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CognitoIdentityProvider',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOB(291427543, _omitFieldNames ? '' : 'serversidetokencheck')
    ..aOS(448902180, _omitFieldNames ? '' : 'clientid')
    ..aOS(485101816, _omitFieldNames ? '' : 'providername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CognitoIdentityProvider clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CognitoIdentityProvider copyWith(
          void Function(CognitoIdentityProvider) updates) =>
      super.copyWith((message) => updates(message as CognitoIdentityProvider))
          as CognitoIdentityProvider;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CognitoIdentityProvider create() => CognitoIdentityProvider._();
  @$core.override
  CognitoIdentityProvider createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CognitoIdentityProvider getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CognitoIdentityProvider>(create);
  static CognitoIdentityProvider? _defaultInstance;

  @$pb.TagNumber(291427543)
  $core.bool get serversidetokencheck => $_getBF(0);
  @$pb.TagNumber(291427543)
  set serversidetokencheck($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(291427543)
  $core.bool hasServersidetokencheck() => $_has(0);
  @$pb.TagNumber(291427543)
  void clearServersidetokencheck() => $_clearField(291427543);

  @$pb.TagNumber(448902180)
  $core.String get clientid => $_getSZ(1);
  @$pb.TagNumber(448902180)
  set clientid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(448902180)
  $core.bool hasClientid() => $_has(1);
  @$pb.TagNumber(448902180)
  void clearClientid() => $_clearField(448902180);

  @$pb.TagNumber(485101816)
  $core.String get providername => $_getSZ(2);
  @$pb.TagNumber(485101816)
  set providername($core.String value) => $_setString(2, value);
  @$pb.TagNumber(485101816)
  $core.bool hasProvidername() => $_has(2);
  @$pb.TagNumber(485101816)
  void clearProvidername() => $_clearField(485101816);
}

class ConcurrentModificationException extends $pb.GeneratedMessage {
  factory ConcurrentModificationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ConcurrentModificationException._();

  factory ConcurrentModificationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConcurrentModificationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConcurrentModificationException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConcurrentModificationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConcurrentModificationException copyWith(
          void Function(ConcurrentModificationException) updates) =>
      super.copyWith(
              (message) => updates(message as ConcurrentModificationException))
          as ConcurrentModificationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConcurrentModificationException create() =>
      ConcurrentModificationException._();
  @$core.override
  ConcurrentModificationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConcurrentModificationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConcurrentModificationException>(
          create);
  static ConcurrentModificationException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class CreateIdentityPoolInput extends $pb.GeneratedMessage {
  factory CreateIdentityPoolInput({
    $core.Iterable<$core.String>? samlproviderarns,
    $core.String? identitypoolname,
    $core.bool? allowclassicflow,
    $core.Iterable<CognitoIdentityProvider>? cognitoidentityproviders,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        supportedloginproviders,
    $core.Iterable<$core.String>? openidconnectproviderarns,
    $core.bool? allowunauthenticatedidentities,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        identitypooltags,
    $core.String? developerprovidername,
  }) {
    final result = create();
    if (samlproviderarns != null)
      result.samlproviderarns.addAll(samlproviderarns);
    if (identitypoolname != null) result.identitypoolname = identitypoolname;
    if (allowclassicflow != null) result.allowclassicflow = allowclassicflow;
    if (cognitoidentityproviders != null)
      result.cognitoidentityproviders.addAll(cognitoidentityproviders);
    if (supportedloginproviders != null)
      result.supportedloginproviders.addEntries(supportedloginproviders);
    if (openidconnectproviderarns != null)
      result.openidconnectproviderarns.addAll(openidconnectproviderarns);
    if (allowunauthenticatedidentities != null)
      result.allowunauthenticatedidentities = allowunauthenticatedidentities;
    if (identitypooltags != null)
      result.identitypooltags.addEntries(identitypooltags);
    if (developerprovidername != null)
      result.developerprovidername = developerprovidername;
    return result;
  }

  CreateIdentityPoolInput._();

  factory CreateIdentityPoolInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateIdentityPoolInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateIdentityPoolInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..pPS(24907670, _omitFieldNames ? '' : 'samlproviderarns')
    ..aOS(41275091, _omitFieldNames ? '' : 'identitypoolname')
    ..aOB(101299523, _omitFieldNames ? '' : 'allowclassicflow')
    ..pPM<CognitoIdentityProvider>(
        109610365, _omitFieldNames ? '' : 'cognitoidentityproviders',
        subBuilder: CognitoIdentityProvider.create)
    ..m<$core.String, $core.String>(
        125752617, _omitFieldNames ? '' : 'supportedloginproviders',
        entryClassName: 'CreateIdentityPoolInput.SupportedloginprovidersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..pPS(137395556, _omitFieldNames ? '' : 'openidconnectproviderarns')
    ..aOB(290546287, _omitFieldNames ? '' : 'allowunauthenticatedidentities')
    ..m<$core.String, $core.String>(
        305630405, _omitFieldNames ? '' : 'identitypooltags',
        entryClassName: 'CreateIdentityPoolInput.IdentitypooltagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(517589792, _omitFieldNames ? '' : 'developerprovidername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateIdentityPoolInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateIdentityPoolInput copyWith(
          void Function(CreateIdentityPoolInput) updates) =>
      super.copyWith((message) => updates(message as CreateIdentityPoolInput))
          as CreateIdentityPoolInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateIdentityPoolInput create() => CreateIdentityPoolInput._();
  @$core.override
  CreateIdentityPoolInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateIdentityPoolInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateIdentityPoolInput>(create);
  static CreateIdentityPoolInput? _defaultInstance;

  @$pb.TagNumber(24907670)
  $pb.PbList<$core.String> get samlproviderarns => $_getList(0);

  @$pb.TagNumber(41275091)
  $core.String get identitypoolname => $_getSZ(1);
  @$pb.TagNumber(41275091)
  set identitypoolname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(41275091)
  $core.bool hasIdentitypoolname() => $_has(1);
  @$pb.TagNumber(41275091)
  void clearIdentitypoolname() => $_clearField(41275091);

  @$pb.TagNumber(101299523)
  $core.bool get allowclassicflow => $_getBF(2);
  @$pb.TagNumber(101299523)
  set allowclassicflow($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(101299523)
  $core.bool hasAllowclassicflow() => $_has(2);
  @$pb.TagNumber(101299523)
  void clearAllowclassicflow() => $_clearField(101299523);

  @$pb.TagNumber(109610365)
  $pb.PbList<CognitoIdentityProvider> get cognitoidentityproviders =>
      $_getList(3);

  @$pb.TagNumber(125752617)
  $pb.PbMap<$core.String, $core.String> get supportedloginproviders =>
      $_getMap(4);

  @$pb.TagNumber(137395556)
  $pb.PbList<$core.String> get openidconnectproviderarns => $_getList(5);

  @$pb.TagNumber(290546287)
  $core.bool get allowunauthenticatedidentities => $_getBF(6);
  @$pb.TagNumber(290546287)
  set allowunauthenticatedidentities($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(290546287)
  $core.bool hasAllowunauthenticatedidentities() => $_has(6);
  @$pb.TagNumber(290546287)
  void clearAllowunauthenticatedidentities() => $_clearField(290546287);

  @$pb.TagNumber(305630405)
  $pb.PbMap<$core.String, $core.String> get identitypooltags => $_getMap(7);

  @$pb.TagNumber(517589792)
  $core.String get developerprovidername => $_getSZ(8);
  @$pb.TagNumber(517589792)
  set developerprovidername($core.String value) => $_setString(8, value);
  @$pb.TagNumber(517589792)
  $core.bool hasDeveloperprovidername() => $_has(8);
  @$pb.TagNumber(517589792)
  void clearDeveloperprovidername() => $_clearField(517589792);
}

class Credentials extends $pb.GeneratedMessage {
  factory Credentials({
    $core.String? sessiontoken,
    $core.String? expiration,
    $core.String? accesskeyid,
    $core.String? secretkey,
  }) {
    final result = create();
    if (sessiontoken != null) result.sessiontoken = sessiontoken;
    if (expiration != null) result.expiration = expiration;
    if (accesskeyid != null) result.accesskeyid = accesskeyid;
    if (secretkey != null) result.secretkey = secretkey;
    return result;
  }

  Credentials._();

  factory Credentials.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Credentials.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Credentials',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(211161069, _omitFieldNames ? '' : 'sessiontoken')
    ..aOS(245879945, _omitFieldNames ? '' : 'expiration')
    ..aOS(453893024, _omitFieldNames ? '' : 'accesskeyid')
    ..aOS(465028505, _omitFieldNames ? '' : 'secretkey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Credentials clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Credentials copyWith(void Function(Credentials) updates) =>
      super.copyWith((message) => updates(message as Credentials))
          as Credentials;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Credentials create() => Credentials._();
  @$core.override
  Credentials createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Credentials getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Credentials>(create);
  static Credentials? _defaultInstance;

  @$pb.TagNumber(211161069)
  $core.String get sessiontoken => $_getSZ(0);
  @$pb.TagNumber(211161069)
  set sessiontoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(211161069)
  $core.bool hasSessiontoken() => $_has(0);
  @$pb.TagNumber(211161069)
  void clearSessiontoken() => $_clearField(211161069);

  @$pb.TagNumber(245879945)
  $core.String get expiration => $_getSZ(1);
  @$pb.TagNumber(245879945)
  set expiration($core.String value) => $_setString(1, value);
  @$pb.TagNumber(245879945)
  $core.bool hasExpiration() => $_has(1);
  @$pb.TagNumber(245879945)
  void clearExpiration() => $_clearField(245879945);

  @$pb.TagNumber(453893024)
  $core.String get accesskeyid => $_getSZ(2);
  @$pb.TagNumber(453893024)
  set accesskeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(453893024)
  $core.bool hasAccesskeyid() => $_has(2);
  @$pb.TagNumber(453893024)
  void clearAccesskeyid() => $_clearField(453893024);

  @$pb.TagNumber(465028505)
  $core.String get secretkey => $_getSZ(3);
  @$pb.TagNumber(465028505)
  set secretkey($core.String value) => $_setString(3, value);
  @$pb.TagNumber(465028505)
  $core.bool hasSecretkey() => $_has(3);
  @$pb.TagNumber(465028505)
  void clearSecretkey() => $_clearField(465028505);
}

class DeleteIdentitiesInput extends $pb.GeneratedMessage {
  factory DeleteIdentitiesInput({
    $core.Iterable<$core.String>? identityidstodelete,
  }) {
    final result = create();
    if (identityidstodelete != null)
      result.identityidstodelete.addAll(identityidstodelete);
    return result;
  }

  DeleteIdentitiesInput._();

  factory DeleteIdentitiesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteIdentitiesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteIdentitiesInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..pPS(340357990, _omitFieldNames ? '' : 'identityidstodelete')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIdentitiesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIdentitiesInput copyWith(
          void Function(DeleteIdentitiesInput) updates) =>
      super.copyWith((message) => updates(message as DeleteIdentitiesInput))
          as DeleteIdentitiesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteIdentitiesInput create() => DeleteIdentitiesInput._();
  @$core.override
  DeleteIdentitiesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteIdentitiesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteIdentitiesInput>(create);
  static DeleteIdentitiesInput? _defaultInstance;

  @$pb.TagNumber(340357990)
  $pb.PbList<$core.String> get identityidstodelete => $_getList(0);
}

class DeleteIdentitiesResponse extends $pb.GeneratedMessage {
  factory DeleteIdentitiesResponse({
    $core.Iterable<UnprocessedIdentityId>? unprocessedidentityids,
  }) {
    final result = create();
    if (unprocessedidentityids != null)
      result.unprocessedidentityids.addAll(unprocessedidentityids);
    return result;
  }

  DeleteIdentitiesResponse._();

  factory DeleteIdentitiesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteIdentitiesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteIdentitiesResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..pPM<UnprocessedIdentityId>(
        532881071, _omitFieldNames ? '' : 'unprocessedidentityids',
        subBuilder: UnprocessedIdentityId.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIdentitiesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIdentitiesResponse copyWith(
          void Function(DeleteIdentitiesResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteIdentitiesResponse))
          as DeleteIdentitiesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteIdentitiesResponse create() => DeleteIdentitiesResponse._();
  @$core.override
  DeleteIdentitiesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteIdentitiesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteIdentitiesResponse>(create);
  static DeleteIdentitiesResponse? _defaultInstance;

  @$pb.TagNumber(532881071)
  $pb.PbList<UnprocessedIdentityId> get unprocessedidentityids => $_getList(0);
}

class DeleteIdentityPoolInput extends $pb.GeneratedMessage {
  factory DeleteIdentityPoolInput({
    $core.String? identitypoolid,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    return result;
  }

  DeleteIdentityPoolInput._();

  factory DeleteIdentityPoolInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteIdentityPoolInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteIdentityPoolInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIdentityPoolInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIdentityPoolInput copyWith(
          void Function(DeleteIdentityPoolInput) updates) =>
      super.copyWith((message) => updates(message as DeleteIdentityPoolInput))
          as DeleteIdentityPoolInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteIdentityPoolInput create() => DeleteIdentityPoolInput._();
  @$core.override
  DeleteIdentityPoolInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteIdentityPoolInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteIdentityPoolInput>(create);
  static DeleteIdentityPoolInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);
}

class DescribeIdentityInput extends $pb.GeneratedMessage {
  factory DescribeIdentityInput({
    $core.String? identityid,
  }) {
    final result = create();
    if (identityid != null) result.identityid = identityid;
    return result;
  }

  DescribeIdentityInput._();

  factory DescribeIdentityInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeIdentityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeIdentityInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeIdentityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeIdentityInput copyWith(
          void Function(DescribeIdentityInput) updates) =>
      super.copyWith((message) => updates(message as DescribeIdentityInput))
          as DescribeIdentityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeIdentityInput create() => DescribeIdentityInput._();
  @$core.override
  DescribeIdentityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeIdentityInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeIdentityInput>(create);
  static DescribeIdentityInput? _defaultInstance;

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(0);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(0);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);
}

class DescribeIdentityPoolInput extends $pb.GeneratedMessage {
  factory DescribeIdentityPoolInput({
    $core.String? identitypoolid,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    return result;
  }

  DescribeIdentityPoolInput._();

  factory DescribeIdentityPoolInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeIdentityPoolInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeIdentityPoolInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeIdentityPoolInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeIdentityPoolInput copyWith(
          void Function(DescribeIdentityPoolInput) updates) =>
      super.copyWith((message) => updates(message as DescribeIdentityPoolInput))
          as DescribeIdentityPoolInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeIdentityPoolInput create() => DescribeIdentityPoolInput._();
  @$core.override
  DescribeIdentityPoolInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeIdentityPoolInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeIdentityPoolInput>(create);
  static DescribeIdentityPoolInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);
}

class DeveloperUserAlreadyRegisteredException extends $pb.GeneratedMessage {
  factory DeveloperUserAlreadyRegisteredException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DeveloperUserAlreadyRegisteredException._();

  factory DeveloperUserAlreadyRegisteredException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeveloperUserAlreadyRegisteredException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeveloperUserAlreadyRegisteredException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeveloperUserAlreadyRegisteredException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeveloperUserAlreadyRegisteredException copyWith(
          void Function(DeveloperUserAlreadyRegisteredException) updates) =>
      super.copyWith((message) =>
              updates(message as DeveloperUserAlreadyRegisteredException))
          as DeveloperUserAlreadyRegisteredException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeveloperUserAlreadyRegisteredException create() =>
      DeveloperUserAlreadyRegisteredException._();
  @$core.override
  DeveloperUserAlreadyRegisteredException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeveloperUserAlreadyRegisteredException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeveloperUserAlreadyRegisteredException>(create);
  static DeveloperUserAlreadyRegisteredException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExternalServiceException extends $pb.GeneratedMessage {
  factory ExternalServiceException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExternalServiceException._();

  factory ExternalServiceException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExternalServiceException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExternalServiceException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExternalServiceException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExternalServiceException copyWith(
          void Function(ExternalServiceException) updates) =>
      super.copyWith((message) => updates(message as ExternalServiceException))
          as ExternalServiceException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExternalServiceException create() => ExternalServiceException._();
  @$core.override
  ExternalServiceException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExternalServiceException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExternalServiceException>(create);
  static ExternalServiceException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class GetCredentialsForIdentityInput extends $pb.GeneratedMessage {
  factory GetCredentialsForIdentityInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? logins,
    $core.String? customrolearn,
    $core.String? identityid,
  }) {
    final result = create();
    if (logins != null) result.logins.addEntries(logins);
    if (customrolearn != null) result.customrolearn = customrolearn;
    if (identityid != null) result.identityid = identityid;
    return result;
  }

  GetCredentialsForIdentityInput._();

  factory GetCredentialsForIdentityInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCredentialsForIdentityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCredentialsForIdentityInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(109702772, _omitFieldNames ? '' : 'logins',
        entryClassName: 'GetCredentialsForIdentityInput.LoginsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(146997938, _omitFieldNames ? '' : 'customrolearn')
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCredentialsForIdentityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCredentialsForIdentityInput copyWith(
          void Function(GetCredentialsForIdentityInput) updates) =>
      super.copyWith(
              (message) => updates(message as GetCredentialsForIdentityInput))
          as GetCredentialsForIdentityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCredentialsForIdentityInput create() =>
      GetCredentialsForIdentityInput._();
  @$core.override
  GetCredentialsForIdentityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCredentialsForIdentityInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCredentialsForIdentityInput>(create);
  static GetCredentialsForIdentityInput? _defaultInstance;

  @$pb.TagNumber(109702772)
  $pb.PbMap<$core.String, $core.String> get logins => $_getMap(0);

  @$pb.TagNumber(146997938)
  $core.String get customrolearn => $_getSZ(1);
  @$pb.TagNumber(146997938)
  set customrolearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(146997938)
  $core.bool hasCustomrolearn() => $_has(1);
  @$pb.TagNumber(146997938)
  void clearCustomrolearn() => $_clearField(146997938);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(2);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(2);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);
}

class GetCredentialsForIdentityResponse extends $pb.GeneratedMessage {
  factory GetCredentialsForIdentityResponse({
    $core.String? identityid,
    Credentials? credentials,
  }) {
    final result = create();
    if (identityid != null) result.identityid = identityid;
    if (credentials != null) result.credentials = credentials;
    return result;
  }

  GetCredentialsForIdentityResponse._();

  factory GetCredentialsForIdentityResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCredentialsForIdentityResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCredentialsForIdentityResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..aOM<Credentials>(381914482, _omitFieldNames ? '' : 'credentials',
        subBuilder: Credentials.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCredentialsForIdentityResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCredentialsForIdentityResponse copyWith(
          void Function(GetCredentialsForIdentityResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetCredentialsForIdentityResponse))
          as GetCredentialsForIdentityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCredentialsForIdentityResponse create() =>
      GetCredentialsForIdentityResponse._();
  @$core.override
  GetCredentialsForIdentityResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCredentialsForIdentityResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCredentialsForIdentityResponse>(
          create);
  static GetCredentialsForIdentityResponse? _defaultInstance;

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(0);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(0);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);

  @$pb.TagNumber(381914482)
  Credentials get credentials => $_getN(1);
  @$pb.TagNumber(381914482)
  set credentials(Credentials value) => $_setField(381914482, value);
  @$pb.TagNumber(381914482)
  $core.bool hasCredentials() => $_has(1);
  @$pb.TagNumber(381914482)
  void clearCredentials() => $_clearField(381914482);
  @$pb.TagNumber(381914482)
  Credentials ensureCredentials() => $_ensure(1);
}

class GetIdInput extends $pb.GeneratedMessage {
  factory GetIdInput({
    $core.String? identitypoolid,
    $core.String? accountid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? logins,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (accountid != null) result.accountid = accountid;
    if (logins != null) result.logins.addEntries(logins);
    return result;
  }

  GetIdInput._();

  factory GetIdInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetIdInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetIdInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOS(65954002, _omitFieldNames ? '' : 'accountid')
    ..m<$core.String, $core.String>(109702772, _omitFieldNames ? '' : 'logins',
        entryClassName: 'GetIdInput.LoginsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIdInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIdInput copyWith(void Function(GetIdInput) updates) =>
      super.copyWith((message) => updates(message as GetIdInput)) as GetIdInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetIdInput create() => GetIdInput._();
  @$core.override
  GetIdInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetIdInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetIdInput>(create);
  static GetIdInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(65954002)
  $core.String get accountid => $_getSZ(1);
  @$pb.TagNumber(65954002)
  set accountid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(65954002)
  $core.bool hasAccountid() => $_has(1);
  @$pb.TagNumber(65954002)
  void clearAccountid() => $_clearField(65954002);

  @$pb.TagNumber(109702772)
  $pb.PbMap<$core.String, $core.String> get logins => $_getMap(2);
}

class GetIdResponse extends $pb.GeneratedMessage {
  factory GetIdResponse({
    $core.String? identityid,
  }) {
    final result = create();
    if (identityid != null) result.identityid = identityid;
    return result;
  }

  GetIdResponse._();

  factory GetIdResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetIdResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetIdResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIdResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIdResponse copyWith(void Function(GetIdResponse) updates) =>
      super.copyWith((message) => updates(message as GetIdResponse))
          as GetIdResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetIdResponse create() => GetIdResponse._();
  @$core.override
  GetIdResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetIdResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetIdResponse>(create);
  static GetIdResponse? _defaultInstance;

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(0);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(0);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);
}

class GetIdentityPoolRolesInput extends $pb.GeneratedMessage {
  factory GetIdentityPoolRolesInput({
    $core.String? identitypoolid,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    return result;
  }

  GetIdentityPoolRolesInput._();

  factory GetIdentityPoolRolesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetIdentityPoolRolesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetIdentityPoolRolesInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIdentityPoolRolesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIdentityPoolRolesInput copyWith(
          void Function(GetIdentityPoolRolesInput) updates) =>
      super.copyWith((message) => updates(message as GetIdentityPoolRolesInput))
          as GetIdentityPoolRolesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetIdentityPoolRolesInput create() => GetIdentityPoolRolesInput._();
  @$core.override
  GetIdentityPoolRolesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetIdentityPoolRolesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetIdentityPoolRolesInput>(create);
  static GetIdentityPoolRolesInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);
}

class GetIdentityPoolRolesResponse extends $pb.GeneratedMessage {
  factory GetIdentityPoolRolesResponse({
    $core.String? identitypoolid,
    $core.Iterable<$core.MapEntry<$core.String, RoleMapping>>? rolemappings,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? roles,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (rolemappings != null) result.rolemappings.addEntries(rolemappings);
    if (roles != null) result.roles.addEntries(roles);
    return result;
  }

  GetIdentityPoolRolesResponse._();

  factory GetIdentityPoolRolesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetIdentityPoolRolesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetIdentityPoolRolesResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..m<$core.String, RoleMapping>(
        96154047, _omitFieldNames ? '' : 'rolemappings',
        entryClassName: 'GetIdentityPoolRolesResponse.RolemappingsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: RoleMapping.create,
        valueDefaultOrMaker: RoleMapping.getDefault,
        packageName: const $pb.PackageName('cognito_identity'))
    ..m<$core.String, $core.String>(511168127, _omitFieldNames ? '' : 'roles',
        entryClassName: 'GetIdentityPoolRolesResponse.RolesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIdentityPoolRolesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIdentityPoolRolesResponse copyWith(
          void Function(GetIdentityPoolRolesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetIdentityPoolRolesResponse))
          as GetIdentityPoolRolesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetIdentityPoolRolesResponse create() =>
      GetIdentityPoolRolesResponse._();
  @$core.override
  GetIdentityPoolRolesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetIdentityPoolRolesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetIdentityPoolRolesResponse>(create);
  static GetIdentityPoolRolesResponse? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(96154047)
  $pb.PbMap<$core.String, RoleMapping> get rolemappings => $_getMap(1);

  @$pb.TagNumber(511168127)
  $pb.PbMap<$core.String, $core.String> get roles => $_getMap(2);
}

class GetOpenIdTokenForDeveloperIdentityInput extends $pb.GeneratedMessage {
  factory GetOpenIdTokenForDeveloperIdentityInput({
    $core.String? identitypoolid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? logins,
    $core.String? identityid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? principaltags,
    $fixnum.Int64? tokenduration,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (logins != null) result.logins.addEntries(logins);
    if (identityid != null) result.identityid = identityid;
    if (principaltags != null) result.principaltags.addEntries(principaltags);
    if (tokenduration != null) result.tokenduration = tokenduration;
    return result;
  }

  GetOpenIdTokenForDeveloperIdentityInput._();

  factory GetOpenIdTokenForDeveloperIdentityInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetOpenIdTokenForDeveloperIdentityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetOpenIdTokenForDeveloperIdentityInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..m<$core.String, $core.String>(109702772, _omitFieldNames ? '' : 'logins',
        entryClassName: 'GetOpenIdTokenForDeveloperIdentityInput.LoginsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..m<$core.String, $core.String>(
        346698229, _omitFieldNames ? '' : 'principaltags',
        entryClassName:
            'GetOpenIdTokenForDeveloperIdentityInput.PrincipaltagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aInt64(402772367, _omitFieldNames ? '' : 'tokenduration')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetOpenIdTokenForDeveloperIdentityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetOpenIdTokenForDeveloperIdentityInput copyWith(
          void Function(GetOpenIdTokenForDeveloperIdentityInput) updates) =>
      super.copyWith((message) =>
              updates(message as GetOpenIdTokenForDeveloperIdentityInput))
          as GetOpenIdTokenForDeveloperIdentityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetOpenIdTokenForDeveloperIdentityInput create() =>
      GetOpenIdTokenForDeveloperIdentityInput._();
  @$core.override
  GetOpenIdTokenForDeveloperIdentityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetOpenIdTokenForDeveloperIdentityInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetOpenIdTokenForDeveloperIdentityInput>(create);
  static GetOpenIdTokenForDeveloperIdentityInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(109702772)
  $pb.PbMap<$core.String, $core.String> get logins => $_getMap(1);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(2);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(2);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);

  @$pb.TagNumber(346698229)
  $pb.PbMap<$core.String, $core.String> get principaltags => $_getMap(3);

  @$pb.TagNumber(402772367)
  $fixnum.Int64 get tokenduration => $_getI64(4);
  @$pb.TagNumber(402772367)
  set tokenduration($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(402772367)
  $core.bool hasTokenduration() => $_has(4);
  @$pb.TagNumber(402772367)
  void clearTokenduration() => $_clearField(402772367);
}

class GetOpenIdTokenForDeveloperIdentityResponse extends $pb.GeneratedMessage {
  factory GetOpenIdTokenForDeveloperIdentityResponse({
    $core.String? identityid,
    $core.String? token,
  }) {
    final result = create();
    if (identityid != null) result.identityid = identityid;
    if (token != null) result.token = token;
    return result;
  }

  GetOpenIdTokenForDeveloperIdentityResponse._();

  factory GetOpenIdTokenForDeveloperIdentityResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetOpenIdTokenForDeveloperIdentityResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetOpenIdTokenForDeveloperIdentityResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..aOS(439704531, _omitFieldNames ? '' : 'token')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetOpenIdTokenForDeveloperIdentityResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetOpenIdTokenForDeveloperIdentityResponse copyWith(
          void Function(GetOpenIdTokenForDeveloperIdentityResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetOpenIdTokenForDeveloperIdentityResponse))
          as GetOpenIdTokenForDeveloperIdentityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetOpenIdTokenForDeveloperIdentityResponse create() =>
      GetOpenIdTokenForDeveloperIdentityResponse._();
  @$core.override
  GetOpenIdTokenForDeveloperIdentityResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetOpenIdTokenForDeveloperIdentityResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetOpenIdTokenForDeveloperIdentityResponse>(create);
  static GetOpenIdTokenForDeveloperIdentityResponse? _defaultInstance;

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(0);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(0);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);

  @$pb.TagNumber(439704531)
  $core.String get token => $_getSZ(1);
  @$pb.TagNumber(439704531)
  set token($core.String value) => $_setString(1, value);
  @$pb.TagNumber(439704531)
  $core.bool hasToken() => $_has(1);
  @$pb.TagNumber(439704531)
  void clearToken() => $_clearField(439704531);
}

class GetOpenIdTokenInput extends $pb.GeneratedMessage {
  factory GetOpenIdTokenInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? logins,
    $core.String? identityid,
  }) {
    final result = create();
    if (logins != null) result.logins.addEntries(logins);
    if (identityid != null) result.identityid = identityid;
    return result;
  }

  GetOpenIdTokenInput._();

  factory GetOpenIdTokenInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetOpenIdTokenInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetOpenIdTokenInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(109702772, _omitFieldNames ? '' : 'logins',
        entryClassName: 'GetOpenIdTokenInput.LoginsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetOpenIdTokenInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetOpenIdTokenInput copyWith(void Function(GetOpenIdTokenInput) updates) =>
      super.copyWith((message) => updates(message as GetOpenIdTokenInput))
          as GetOpenIdTokenInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetOpenIdTokenInput create() => GetOpenIdTokenInput._();
  @$core.override
  GetOpenIdTokenInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetOpenIdTokenInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetOpenIdTokenInput>(create);
  static GetOpenIdTokenInput? _defaultInstance;

  @$pb.TagNumber(109702772)
  $pb.PbMap<$core.String, $core.String> get logins => $_getMap(0);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(1);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(1);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);
}

class GetOpenIdTokenResponse extends $pb.GeneratedMessage {
  factory GetOpenIdTokenResponse({
    $core.String? identityid,
    $core.String? token,
  }) {
    final result = create();
    if (identityid != null) result.identityid = identityid;
    if (token != null) result.token = token;
    return result;
  }

  GetOpenIdTokenResponse._();

  factory GetOpenIdTokenResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetOpenIdTokenResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetOpenIdTokenResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..aOS(439704531, _omitFieldNames ? '' : 'token')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetOpenIdTokenResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetOpenIdTokenResponse copyWith(
          void Function(GetOpenIdTokenResponse) updates) =>
      super.copyWith((message) => updates(message as GetOpenIdTokenResponse))
          as GetOpenIdTokenResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetOpenIdTokenResponse create() => GetOpenIdTokenResponse._();
  @$core.override
  GetOpenIdTokenResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetOpenIdTokenResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetOpenIdTokenResponse>(create);
  static GetOpenIdTokenResponse? _defaultInstance;

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(0);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(0);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);

  @$pb.TagNumber(439704531)
  $core.String get token => $_getSZ(1);
  @$pb.TagNumber(439704531)
  set token($core.String value) => $_setString(1, value);
  @$pb.TagNumber(439704531)
  $core.bool hasToken() => $_has(1);
  @$pb.TagNumber(439704531)
  void clearToken() => $_clearField(439704531);
}

class GetPrincipalTagAttributeMapInput extends $pb.GeneratedMessage {
  factory GetPrincipalTagAttributeMapInput({
    $core.String? identitypoolid,
    $core.String? identityprovidername,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (identityprovidername != null)
      result.identityprovidername = identityprovidername;
    return result;
  }

  GetPrincipalTagAttributeMapInput._();

  factory GetPrincipalTagAttributeMapInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPrincipalTagAttributeMapInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPrincipalTagAttributeMapInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOS(419175410, _omitFieldNames ? '' : 'identityprovidername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPrincipalTagAttributeMapInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPrincipalTagAttributeMapInput copyWith(
          void Function(GetPrincipalTagAttributeMapInput) updates) =>
      super.copyWith(
              (message) => updates(message as GetPrincipalTagAttributeMapInput))
          as GetPrincipalTagAttributeMapInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPrincipalTagAttributeMapInput create() =>
      GetPrincipalTagAttributeMapInput._();
  @$core.override
  GetPrincipalTagAttributeMapInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPrincipalTagAttributeMapInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPrincipalTagAttributeMapInput>(
          create);
  static GetPrincipalTagAttributeMapInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(419175410)
  $core.String get identityprovidername => $_getSZ(1);
  @$pb.TagNumber(419175410)
  set identityprovidername($core.String value) => $_setString(1, value);
  @$pb.TagNumber(419175410)
  $core.bool hasIdentityprovidername() => $_has(1);
  @$pb.TagNumber(419175410)
  void clearIdentityprovidername() => $_clearField(419175410);
}

class GetPrincipalTagAttributeMapResponse extends $pb.GeneratedMessage {
  factory GetPrincipalTagAttributeMapResponse({
    $core.String? identitypoolid,
    $core.bool? usedefaults,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? principaltags,
    $core.String? identityprovidername,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (usedefaults != null) result.usedefaults = usedefaults;
    if (principaltags != null) result.principaltags.addEntries(principaltags);
    if (identityprovidername != null)
      result.identityprovidername = identityprovidername;
    return result;
  }

  GetPrincipalTagAttributeMapResponse._();

  factory GetPrincipalTagAttributeMapResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPrincipalTagAttributeMapResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPrincipalTagAttributeMapResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOB(306413999, _omitFieldNames ? '' : 'usedefaults')
    ..m<$core.String, $core.String>(
        346698229, _omitFieldNames ? '' : 'principaltags',
        entryClassName:
            'GetPrincipalTagAttributeMapResponse.PrincipaltagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(419175410, _omitFieldNames ? '' : 'identityprovidername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPrincipalTagAttributeMapResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPrincipalTagAttributeMapResponse copyWith(
          void Function(GetPrincipalTagAttributeMapResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetPrincipalTagAttributeMapResponse))
          as GetPrincipalTagAttributeMapResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPrincipalTagAttributeMapResponse create() =>
      GetPrincipalTagAttributeMapResponse._();
  @$core.override
  GetPrincipalTagAttributeMapResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPrincipalTagAttributeMapResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetPrincipalTagAttributeMapResponse>(create);
  static GetPrincipalTagAttributeMapResponse? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(306413999)
  $core.bool get usedefaults => $_getBF(1);
  @$pb.TagNumber(306413999)
  set usedefaults($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(306413999)
  $core.bool hasUsedefaults() => $_has(1);
  @$pb.TagNumber(306413999)
  void clearUsedefaults() => $_clearField(306413999);

  @$pb.TagNumber(346698229)
  $pb.PbMap<$core.String, $core.String> get principaltags => $_getMap(2);

  @$pb.TagNumber(419175410)
  $core.String get identityprovidername => $_getSZ(3);
  @$pb.TagNumber(419175410)
  set identityprovidername($core.String value) => $_setString(3, value);
  @$pb.TagNumber(419175410)
  $core.bool hasIdentityprovidername() => $_has(3);
  @$pb.TagNumber(419175410)
  void clearIdentityprovidername() => $_clearField(419175410);
}

class IdentityDescription extends $pb.GeneratedMessage {
  factory IdentityDescription({
    $core.String? lastmodifieddate,
    $core.Iterable<$core.String>? logins,
    $core.String? identityid,
    $core.String? creationdate,
  }) {
    final result = create();
    if (lastmodifieddate != null) result.lastmodifieddate = lastmodifieddate;
    if (logins != null) result.logins.addAll(logins);
    if (identityid != null) result.identityid = identityid;
    if (creationdate != null) result.creationdate = creationdate;
    return result;
  }

  IdentityDescription._();

  factory IdentityDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IdentityDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IdentityDescription',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(24249427, _omitFieldNames ? '' : 'lastmodifieddate')
    ..pPS(109702772, _omitFieldNames ? '' : 'logins')
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdentityDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdentityDescription copyWith(void Function(IdentityDescription) updates) =>
      super.copyWith((message) => updates(message as IdentityDescription))
          as IdentityDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IdentityDescription create() => IdentityDescription._();
  @$core.override
  IdentityDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IdentityDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IdentityDescription>(create);
  static IdentityDescription? _defaultInstance;

  @$pb.TagNumber(24249427)
  $core.String get lastmodifieddate => $_getSZ(0);
  @$pb.TagNumber(24249427)
  set lastmodifieddate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(24249427)
  $core.bool hasLastmodifieddate() => $_has(0);
  @$pb.TagNumber(24249427)
  void clearLastmodifieddate() => $_clearField(24249427);

  @$pb.TagNumber(109702772)
  $pb.PbList<$core.String> get logins => $_getList(1);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(2);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(2);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(3);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(3);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);
}

class IdentityPool extends $pb.GeneratedMessage {
  factory IdentityPool({
    $core.String? identitypoolid,
    $core.Iterable<$core.String>? samlproviderarns,
    $core.String? identitypoolname,
    $core.bool? allowclassicflow,
    $core.Iterable<CognitoIdentityProvider>? cognitoidentityproviders,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        supportedloginproviders,
    $core.Iterable<$core.String>? openidconnectproviderarns,
    $core.bool? allowunauthenticatedidentities,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        identitypooltags,
    $core.String? developerprovidername,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (samlproviderarns != null)
      result.samlproviderarns.addAll(samlproviderarns);
    if (identitypoolname != null) result.identitypoolname = identitypoolname;
    if (allowclassicflow != null) result.allowclassicflow = allowclassicflow;
    if (cognitoidentityproviders != null)
      result.cognitoidentityproviders.addAll(cognitoidentityproviders);
    if (supportedloginproviders != null)
      result.supportedloginproviders.addEntries(supportedloginproviders);
    if (openidconnectproviderarns != null)
      result.openidconnectproviderarns.addAll(openidconnectproviderarns);
    if (allowunauthenticatedidentities != null)
      result.allowunauthenticatedidentities = allowunauthenticatedidentities;
    if (identitypooltags != null)
      result.identitypooltags.addEntries(identitypooltags);
    if (developerprovidername != null)
      result.developerprovidername = developerprovidername;
    return result;
  }

  IdentityPool._();

  factory IdentityPool.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IdentityPool.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IdentityPool',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..pPS(24907670, _omitFieldNames ? '' : 'samlproviderarns')
    ..aOS(41275091, _omitFieldNames ? '' : 'identitypoolname')
    ..aOB(101299523, _omitFieldNames ? '' : 'allowclassicflow')
    ..pPM<CognitoIdentityProvider>(
        109610365, _omitFieldNames ? '' : 'cognitoidentityproviders',
        subBuilder: CognitoIdentityProvider.create)
    ..m<$core.String, $core.String>(
        125752617, _omitFieldNames ? '' : 'supportedloginproviders',
        entryClassName: 'IdentityPool.SupportedloginprovidersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..pPS(137395556, _omitFieldNames ? '' : 'openidconnectproviderarns')
    ..aOB(290546287, _omitFieldNames ? '' : 'allowunauthenticatedidentities')
    ..m<$core.String, $core.String>(
        305630405, _omitFieldNames ? '' : 'identitypooltags',
        entryClassName: 'IdentityPool.IdentitypooltagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(517589792, _omitFieldNames ? '' : 'developerprovidername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdentityPool clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdentityPool copyWith(void Function(IdentityPool) updates) =>
      super.copyWith((message) => updates(message as IdentityPool))
          as IdentityPool;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IdentityPool create() => IdentityPool._();
  @$core.override
  IdentityPool createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IdentityPool getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IdentityPool>(create);
  static IdentityPool? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(24907670)
  $pb.PbList<$core.String> get samlproviderarns => $_getList(1);

  @$pb.TagNumber(41275091)
  $core.String get identitypoolname => $_getSZ(2);
  @$pb.TagNumber(41275091)
  set identitypoolname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(41275091)
  $core.bool hasIdentitypoolname() => $_has(2);
  @$pb.TagNumber(41275091)
  void clearIdentitypoolname() => $_clearField(41275091);

  @$pb.TagNumber(101299523)
  $core.bool get allowclassicflow => $_getBF(3);
  @$pb.TagNumber(101299523)
  set allowclassicflow($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(101299523)
  $core.bool hasAllowclassicflow() => $_has(3);
  @$pb.TagNumber(101299523)
  void clearAllowclassicflow() => $_clearField(101299523);

  @$pb.TagNumber(109610365)
  $pb.PbList<CognitoIdentityProvider> get cognitoidentityproviders =>
      $_getList(4);

  @$pb.TagNumber(125752617)
  $pb.PbMap<$core.String, $core.String> get supportedloginproviders =>
      $_getMap(5);

  @$pb.TagNumber(137395556)
  $pb.PbList<$core.String> get openidconnectproviderarns => $_getList(6);

  @$pb.TagNumber(290546287)
  $core.bool get allowunauthenticatedidentities => $_getBF(7);
  @$pb.TagNumber(290546287)
  set allowunauthenticatedidentities($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(290546287)
  $core.bool hasAllowunauthenticatedidentities() => $_has(7);
  @$pb.TagNumber(290546287)
  void clearAllowunauthenticatedidentities() => $_clearField(290546287);

  @$pb.TagNumber(305630405)
  $pb.PbMap<$core.String, $core.String> get identitypooltags => $_getMap(8);

  @$pb.TagNumber(517589792)
  $core.String get developerprovidername => $_getSZ(9);
  @$pb.TagNumber(517589792)
  set developerprovidername($core.String value) => $_setString(9, value);
  @$pb.TagNumber(517589792)
  $core.bool hasDeveloperprovidername() => $_has(9);
  @$pb.TagNumber(517589792)
  void clearDeveloperprovidername() => $_clearField(517589792);
}

class IdentityPoolShortDescription extends $pb.GeneratedMessage {
  factory IdentityPoolShortDescription({
    $core.String? identitypoolid,
    $core.String? identitypoolname,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (identitypoolname != null) result.identitypoolname = identitypoolname;
    return result;
  }

  IdentityPoolShortDescription._();

  factory IdentityPoolShortDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IdentityPoolShortDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IdentityPoolShortDescription',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOS(41275091, _omitFieldNames ? '' : 'identitypoolname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdentityPoolShortDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdentityPoolShortDescription copyWith(
          void Function(IdentityPoolShortDescription) updates) =>
      super.copyWith(
              (message) => updates(message as IdentityPoolShortDescription))
          as IdentityPoolShortDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IdentityPoolShortDescription create() =>
      IdentityPoolShortDescription._();
  @$core.override
  IdentityPoolShortDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IdentityPoolShortDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IdentityPoolShortDescription>(create);
  static IdentityPoolShortDescription? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(41275091)
  $core.String get identitypoolname => $_getSZ(1);
  @$pb.TagNumber(41275091)
  set identitypoolname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(41275091)
  $core.bool hasIdentitypoolname() => $_has(1);
  @$pb.TagNumber(41275091)
  void clearIdentitypoolname() => $_clearField(41275091);
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
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
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

class InvalidIdentityPoolConfigurationException extends $pb.GeneratedMessage {
  factory InvalidIdentityPoolConfigurationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidIdentityPoolConfigurationException._();

  factory InvalidIdentityPoolConfigurationException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidIdentityPoolConfigurationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidIdentityPoolConfigurationException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidIdentityPoolConfigurationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidIdentityPoolConfigurationException copyWith(
          void Function(InvalidIdentityPoolConfigurationException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidIdentityPoolConfigurationException))
          as InvalidIdentityPoolConfigurationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidIdentityPoolConfigurationException create() =>
      InvalidIdentityPoolConfigurationException._();
  @$core.override
  InvalidIdentityPoolConfigurationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidIdentityPoolConfigurationException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidIdentityPoolConfigurationException>(create);
  static InvalidIdentityPoolConfigurationException? _defaultInstance;

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
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
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
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
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

class ListIdentitiesInput extends $pb.GeneratedMessage {
  factory ListIdentitiesInput({
    $core.String? identitypoolid,
    $core.bool? hidedisabled,
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (hidedisabled != null) result.hidedisabled = hidedisabled;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListIdentitiesInput._();

  factory ListIdentitiesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListIdentitiesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListIdentitiesInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOB(189789070, _omitFieldNames ? '' : 'hidedisabled')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIdentitiesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIdentitiesInput copyWith(void Function(ListIdentitiesInput) updates) =>
      super.copyWith((message) => updates(message as ListIdentitiesInput))
          as ListIdentitiesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListIdentitiesInput create() => ListIdentitiesInput._();
  @$core.override
  ListIdentitiesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListIdentitiesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListIdentitiesInput>(create);
  static ListIdentitiesInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(189789070)
  $core.bool get hidedisabled => $_getBF(1);
  @$pb.TagNumber(189789070)
  set hidedisabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(189789070)
  $core.bool hasHidedisabled() => $_has(1);
  @$pb.TagNumber(189789070)
  void clearHidedisabled() => $_clearField(189789070);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(3);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(3);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListIdentitiesResponse extends $pb.GeneratedMessage {
  factory ListIdentitiesResponse({
    $core.String? identitypoolid,
    $core.String? nexttoken,
    $core.Iterable<IdentityDescription>? identities,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (identities != null) result.identities.addAll(identities);
    return result;
  }

  ListIdentitiesResponse._();

  factory ListIdentitiesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListIdentitiesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListIdentitiesResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<IdentityDescription>(452470428, _omitFieldNames ? '' : 'identities',
        subBuilder: IdentityDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIdentitiesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIdentitiesResponse copyWith(
          void Function(ListIdentitiesResponse) updates) =>
      super.copyWith((message) => updates(message as ListIdentitiesResponse))
          as ListIdentitiesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListIdentitiesResponse create() => ListIdentitiesResponse._();
  @$core.override
  ListIdentitiesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListIdentitiesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListIdentitiesResponse>(create);
  static ListIdentitiesResponse? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(452470428)
  $pb.PbList<IdentityDescription> get identities => $_getList(2);
}

class ListIdentityPoolsInput extends $pb.GeneratedMessage {
  factory ListIdentityPoolsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListIdentityPoolsInput._();

  factory ListIdentityPoolsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListIdentityPoolsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListIdentityPoolsInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIdentityPoolsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIdentityPoolsInput copyWith(
          void Function(ListIdentityPoolsInput) updates) =>
      super.copyWith((message) => updates(message as ListIdentityPoolsInput))
          as ListIdentityPoolsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListIdentityPoolsInput create() => ListIdentityPoolsInput._();
  @$core.override
  ListIdentityPoolsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListIdentityPoolsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListIdentityPoolsInput>(create);
  static ListIdentityPoolsInput? _defaultInstance;

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

class ListIdentityPoolsResponse extends $pb.GeneratedMessage {
  factory ListIdentityPoolsResponse({
    $core.Iterable<IdentityPoolShortDescription>? identitypools,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (identitypools != null) result.identitypools.addAll(identitypools);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListIdentityPoolsResponse._();

  factory ListIdentityPoolsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListIdentityPoolsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListIdentityPoolsResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..pPM<IdentityPoolShortDescription>(
        158350303, _omitFieldNames ? '' : 'identitypools',
        subBuilder: IdentityPoolShortDescription.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIdentityPoolsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIdentityPoolsResponse copyWith(
          void Function(ListIdentityPoolsResponse) updates) =>
      super.copyWith((message) => updates(message as ListIdentityPoolsResponse))
          as ListIdentityPoolsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListIdentityPoolsResponse create() => ListIdentityPoolsResponse._();
  @$core.override
  ListIdentityPoolsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListIdentityPoolsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListIdentityPoolsResponse>(create);
  static ListIdentityPoolsResponse? _defaultInstance;

  @$pb.TagNumber(158350303)
  $pb.PbList<IdentityPoolShortDescription> get identitypools => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
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
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
  }) {
    final result = create();
    if (tags != null) result.tags.addEntries(tags);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'ListTagsForResourceResponse.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
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
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(0);
}

class LookupDeveloperIdentityInput extends $pb.GeneratedMessage {
  factory LookupDeveloperIdentityInput({
    $core.String? identitypoolid,
    $core.String? nexttoken,
    $core.String? identityid,
    $core.int? maxresults,
    $core.String? developeruseridentifier,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (identityid != null) result.identityid = identityid;
    if (maxresults != null) result.maxresults = maxresults;
    if (developeruseridentifier != null)
      result.developeruseridentifier = developeruseridentifier;
    return result;
  }

  LookupDeveloperIdentityInput._();

  factory LookupDeveloperIdentityInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LookupDeveloperIdentityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LookupDeveloperIdentityInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(349826618, _omitFieldNames ? '' : 'developeruseridentifier')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupDeveloperIdentityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupDeveloperIdentityInput copyWith(
          void Function(LookupDeveloperIdentityInput) updates) =>
      super.copyWith(
              (message) => updates(message as LookupDeveloperIdentityInput))
          as LookupDeveloperIdentityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LookupDeveloperIdentityInput create() =>
      LookupDeveloperIdentityInput._();
  @$core.override
  LookupDeveloperIdentityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LookupDeveloperIdentityInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LookupDeveloperIdentityInput>(create);
  static LookupDeveloperIdentityInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(2);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(2);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(3);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(3);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(349826618)
  $core.String get developeruseridentifier => $_getSZ(4);
  @$pb.TagNumber(349826618)
  set developeruseridentifier($core.String value) => $_setString(4, value);
  @$pb.TagNumber(349826618)
  $core.bool hasDeveloperuseridentifier() => $_has(4);
  @$pb.TagNumber(349826618)
  void clearDeveloperuseridentifier() => $_clearField(349826618);
}

class LookupDeveloperIdentityResponse extends $pb.GeneratedMessage {
  factory LookupDeveloperIdentityResponse({
    $core.Iterable<$core.String>? developeruseridentifierlist,
    $core.String? nexttoken,
    $core.String? identityid,
  }) {
    final result = create();
    if (developeruseridentifierlist != null)
      result.developeruseridentifierlist.addAll(developeruseridentifierlist);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (identityid != null) result.identityid = identityid;
    return result;
  }

  LookupDeveloperIdentityResponse._();

  factory LookupDeveloperIdentityResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LookupDeveloperIdentityResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LookupDeveloperIdentityResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..pPS(86156134, _omitFieldNames ? '' : 'developeruseridentifierlist')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupDeveloperIdentityResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LookupDeveloperIdentityResponse copyWith(
          void Function(LookupDeveloperIdentityResponse) updates) =>
      super.copyWith(
              (message) => updates(message as LookupDeveloperIdentityResponse))
          as LookupDeveloperIdentityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LookupDeveloperIdentityResponse create() =>
      LookupDeveloperIdentityResponse._();
  @$core.override
  LookupDeveloperIdentityResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LookupDeveloperIdentityResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LookupDeveloperIdentityResponse>(
          create);
  static LookupDeveloperIdentityResponse? _defaultInstance;

  @$pb.TagNumber(86156134)
  $pb.PbList<$core.String> get developeruseridentifierlist => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(2);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(2);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);
}

class MappingRule extends $pb.GeneratedMessage {
  factory MappingRule({
    $core.String? claim,
    MappingRuleMatchType? matchtype,
    $core.String? value,
    $core.String? rolearn,
  }) {
    final result = create();
    if (claim != null) result.claim = claim;
    if (matchtype != null) result.matchtype = matchtype;
    if (value != null) result.value = value;
    if (rolearn != null) result.rolearn = rolearn;
    return result;
  }

  MappingRule._();

  factory MappingRule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MappingRule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MappingRule',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(48235024, _omitFieldNames ? '' : 'claim')
    ..aE<MappingRuleMatchType>(94592295, _omitFieldNames ? '' : 'matchtype',
        enumValues: MappingRuleMatchType.values)
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..aOS(327777153, _omitFieldNames ? '' : 'rolearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MappingRule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MappingRule copyWith(void Function(MappingRule) updates) =>
      super.copyWith((message) => updates(message as MappingRule))
          as MappingRule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MappingRule create() => MappingRule._();
  @$core.override
  MappingRule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MappingRule getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MappingRule>(create);
  static MappingRule? _defaultInstance;

  @$pb.TagNumber(48235024)
  $core.String get claim => $_getSZ(0);
  @$pb.TagNumber(48235024)
  set claim($core.String value) => $_setString(0, value);
  @$pb.TagNumber(48235024)
  $core.bool hasClaim() => $_has(0);
  @$pb.TagNumber(48235024)
  void clearClaim() => $_clearField(48235024);

  @$pb.TagNumber(94592295)
  MappingRuleMatchType get matchtype => $_getN(1);
  @$pb.TagNumber(94592295)
  set matchtype(MappingRuleMatchType value) => $_setField(94592295, value);
  @$pb.TagNumber(94592295)
  $core.bool hasMatchtype() => $_has(1);
  @$pb.TagNumber(94592295)
  void clearMatchtype() => $_clearField(94592295);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(2);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(2, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(2);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(327777153)
  $core.String get rolearn => $_getSZ(3);
  @$pb.TagNumber(327777153)
  set rolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(327777153)
  $core.bool hasRolearn() => $_has(3);
  @$pb.TagNumber(327777153)
  void clearRolearn() => $_clearField(327777153);
}

class MergeDeveloperIdentitiesInput extends $pb.GeneratedMessage {
  factory MergeDeveloperIdentitiesInput({
    $core.String? identitypoolid,
    $core.String? sourceuseridentifier,
    $core.String? destinationuseridentifier,
    $core.String? developerprovidername,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (sourceuseridentifier != null)
      result.sourceuseridentifier = sourceuseridentifier;
    if (destinationuseridentifier != null)
      result.destinationuseridentifier = destinationuseridentifier;
    if (developerprovidername != null)
      result.developerprovidername = developerprovidername;
    return result;
  }

  MergeDeveloperIdentitiesInput._();

  factory MergeDeveloperIdentitiesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MergeDeveloperIdentitiesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MergeDeveloperIdentitiesInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOS(222471153, _omitFieldNames ? '' : 'sourceuseridentifier')
    ..aOS(333657056, _omitFieldNames ? '' : 'destinationuseridentifier')
    ..aOS(517589792, _omitFieldNames ? '' : 'developerprovidername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeDeveloperIdentitiesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeDeveloperIdentitiesInput copyWith(
          void Function(MergeDeveloperIdentitiesInput) updates) =>
      super.copyWith(
              (message) => updates(message as MergeDeveloperIdentitiesInput))
          as MergeDeveloperIdentitiesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MergeDeveloperIdentitiesInput create() =>
      MergeDeveloperIdentitiesInput._();
  @$core.override
  MergeDeveloperIdentitiesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MergeDeveloperIdentitiesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MergeDeveloperIdentitiesInput>(create);
  static MergeDeveloperIdentitiesInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(222471153)
  $core.String get sourceuseridentifier => $_getSZ(1);
  @$pb.TagNumber(222471153)
  set sourceuseridentifier($core.String value) => $_setString(1, value);
  @$pb.TagNumber(222471153)
  $core.bool hasSourceuseridentifier() => $_has(1);
  @$pb.TagNumber(222471153)
  void clearSourceuseridentifier() => $_clearField(222471153);

  @$pb.TagNumber(333657056)
  $core.String get destinationuseridentifier => $_getSZ(2);
  @$pb.TagNumber(333657056)
  set destinationuseridentifier($core.String value) => $_setString(2, value);
  @$pb.TagNumber(333657056)
  $core.bool hasDestinationuseridentifier() => $_has(2);
  @$pb.TagNumber(333657056)
  void clearDestinationuseridentifier() => $_clearField(333657056);

  @$pb.TagNumber(517589792)
  $core.String get developerprovidername => $_getSZ(3);
  @$pb.TagNumber(517589792)
  set developerprovidername($core.String value) => $_setString(3, value);
  @$pb.TagNumber(517589792)
  $core.bool hasDeveloperprovidername() => $_has(3);
  @$pb.TagNumber(517589792)
  void clearDeveloperprovidername() => $_clearField(517589792);
}

class MergeDeveloperIdentitiesResponse extends $pb.GeneratedMessage {
  factory MergeDeveloperIdentitiesResponse({
    $core.String? identityid,
  }) {
    final result = create();
    if (identityid != null) result.identityid = identityid;
    return result;
  }

  MergeDeveloperIdentitiesResponse._();

  factory MergeDeveloperIdentitiesResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MergeDeveloperIdentitiesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MergeDeveloperIdentitiesResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeDeveloperIdentitiesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MergeDeveloperIdentitiesResponse copyWith(
          void Function(MergeDeveloperIdentitiesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as MergeDeveloperIdentitiesResponse))
          as MergeDeveloperIdentitiesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MergeDeveloperIdentitiesResponse create() =>
      MergeDeveloperIdentitiesResponse._();
  @$core.override
  MergeDeveloperIdentitiesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MergeDeveloperIdentitiesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MergeDeveloperIdentitiesResponse>(
          create);
  static MergeDeveloperIdentitiesResponse? _defaultInstance;

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(0);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(0);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);
}

class NotAuthorizedException extends $pb.GeneratedMessage {
  factory NotAuthorizedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NotAuthorizedException._();

  factory NotAuthorizedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NotAuthorizedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NotAuthorizedException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotAuthorizedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotAuthorizedException copyWith(
          void Function(NotAuthorizedException) updates) =>
      super.copyWith((message) => updates(message as NotAuthorizedException))
          as NotAuthorizedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NotAuthorizedException create() => NotAuthorizedException._();
  @$core.override
  NotAuthorizedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NotAuthorizedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NotAuthorizedException>(create);
  static NotAuthorizedException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ResourceConflictException extends $pb.GeneratedMessage {
  factory ResourceConflictException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourceConflictException._();

  factory ResourceConflictException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceConflictException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceConflictException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceConflictException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceConflictException copyWith(
          void Function(ResourceConflictException) updates) =>
      super.copyWith((message) => updates(message as ResourceConflictException))
          as ResourceConflictException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceConflictException create() => ResourceConflictException._();
  @$core.override
  ResourceConflictException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceConflictException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceConflictException>(create);
  static ResourceConflictException? _defaultInstance;

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
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
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

class RoleMapping extends $pb.GeneratedMessage {
  factory RoleMapping({
    RulesConfigurationType? rulesconfiguration,
    RoleMappingType? type,
    AmbiguousRoleResolutionType? ambiguousroleresolution,
  }) {
    final result = create();
    if (rulesconfiguration != null)
      result.rulesconfiguration = rulesconfiguration;
    if (type != null) result.type = type;
    if (ambiguousroleresolution != null)
      result.ambiguousroleresolution = ambiguousroleresolution;
    return result;
  }

  RoleMapping._();

  factory RoleMapping.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RoleMapping.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RoleMapping',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOM<RulesConfigurationType>(
        123571251, _omitFieldNames ? '' : 'rulesconfiguration',
        subBuilder: RulesConfigurationType.create)
    ..aE<RoleMappingType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: RoleMappingType.values)
    ..aE<AmbiguousRoleResolutionType>(
        410304778, _omitFieldNames ? '' : 'ambiguousroleresolution',
        enumValues: AmbiguousRoleResolutionType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RoleMapping clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RoleMapping copyWith(void Function(RoleMapping) updates) =>
      super.copyWith((message) => updates(message as RoleMapping))
          as RoleMapping;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RoleMapping create() => RoleMapping._();
  @$core.override
  RoleMapping createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RoleMapping getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RoleMapping>(create);
  static RoleMapping? _defaultInstance;

  @$pb.TagNumber(123571251)
  RulesConfigurationType get rulesconfiguration => $_getN(0);
  @$pb.TagNumber(123571251)
  set rulesconfiguration(RulesConfigurationType value) =>
      $_setField(123571251, value);
  @$pb.TagNumber(123571251)
  $core.bool hasRulesconfiguration() => $_has(0);
  @$pb.TagNumber(123571251)
  void clearRulesconfiguration() => $_clearField(123571251);
  @$pb.TagNumber(123571251)
  RulesConfigurationType ensureRulesconfiguration() => $_ensure(0);

  @$pb.TagNumber(290836590)
  RoleMappingType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(RoleMappingType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(410304778)
  AmbiguousRoleResolutionType get ambiguousroleresolution => $_getN(2);
  @$pb.TagNumber(410304778)
  set ambiguousroleresolution(AmbiguousRoleResolutionType value) =>
      $_setField(410304778, value);
  @$pb.TagNumber(410304778)
  $core.bool hasAmbiguousroleresolution() => $_has(2);
  @$pb.TagNumber(410304778)
  void clearAmbiguousroleresolution() => $_clearField(410304778);
}

class RulesConfigurationType extends $pb.GeneratedMessage {
  factory RulesConfigurationType({
    $core.Iterable<MappingRule>? rules,
  }) {
    final result = create();
    if (rules != null) result.rules.addAll(rules);
    return result;
  }

  RulesConfigurationType._();

  factory RulesConfigurationType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RulesConfigurationType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RulesConfigurationType',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..pPM<MappingRule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: MappingRule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RulesConfigurationType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RulesConfigurationType copyWith(
          void Function(RulesConfigurationType) updates) =>
      super.copyWith((message) => updates(message as RulesConfigurationType))
          as RulesConfigurationType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RulesConfigurationType create() => RulesConfigurationType._();
  @$core.override
  RulesConfigurationType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RulesConfigurationType getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RulesConfigurationType>(create);
  static RulesConfigurationType? _defaultInstance;

  @$pb.TagNumber(42675585)
  $pb.PbList<MappingRule> get rules => $_getList(0);
}

class SetIdentityPoolRolesInput extends $pb.GeneratedMessage {
  factory SetIdentityPoolRolesInput({
    $core.String? identitypoolid,
    $core.Iterable<$core.MapEntry<$core.String, RoleMapping>>? rolemappings,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? roles,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (rolemappings != null) result.rolemappings.addEntries(rolemappings);
    if (roles != null) result.roles.addEntries(roles);
    return result;
  }

  SetIdentityPoolRolesInput._();

  factory SetIdentityPoolRolesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetIdentityPoolRolesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetIdentityPoolRolesInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..m<$core.String, RoleMapping>(
        96154047, _omitFieldNames ? '' : 'rolemappings',
        entryClassName: 'SetIdentityPoolRolesInput.RolemappingsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: RoleMapping.create,
        valueDefaultOrMaker: RoleMapping.getDefault,
        packageName: const $pb.PackageName('cognito_identity'))
    ..m<$core.String, $core.String>(511168127, _omitFieldNames ? '' : 'roles',
        entryClassName: 'SetIdentityPoolRolesInput.RolesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetIdentityPoolRolesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetIdentityPoolRolesInput copyWith(
          void Function(SetIdentityPoolRolesInput) updates) =>
      super.copyWith((message) => updates(message as SetIdentityPoolRolesInput))
          as SetIdentityPoolRolesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetIdentityPoolRolesInput create() => SetIdentityPoolRolesInput._();
  @$core.override
  SetIdentityPoolRolesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetIdentityPoolRolesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetIdentityPoolRolesInput>(create);
  static SetIdentityPoolRolesInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(96154047)
  $pb.PbMap<$core.String, RoleMapping> get rolemappings => $_getMap(1);

  @$pb.TagNumber(511168127)
  $pb.PbMap<$core.String, $core.String> get roles => $_getMap(2);
}

class SetPrincipalTagAttributeMapInput extends $pb.GeneratedMessage {
  factory SetPrincipalTagAttributeMapInput({
    $core.String? identitypoolid,
    $core.bool? usedefaults,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? principaltags,
    $core.String? identityprovidername,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (usedefaults != null) result.usedefaults = usedefaults;
    if (principaltags != null) result.principaltags.addEntries(principaltags);
    if (identityprovidername != null)
      result.identityprovidername = identityprovidername;
    return result;
  }

  SetPrincipalTagAttributeMapInput._();

  factory SetPrincipalTagAttributeMapInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetPrincipalTagAttributeMapInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetPrincipalTagAttributeMapInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOB(306413999, _omitFieldNames ? '' : 'usedefaults')
    ..m<$core.String, $core.String>(
        346698229, _omitFieldNames ? '' : 'principaltags',
        entryClassName: 'SetPrincipalTagAttributeMapInput.PrincipaltagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(419175410, _omitFieldNames ? '' : 'identityprovidername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetPrincipalTagAttributeMapInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetPrincipalTagAttributeMapInput copyWith(
          void Function(SetPrincipalTagAttributeMapInput) updates) =>
      super.copyWith(
              (message) => updates(message as SetPrincipalTagAttributeMapInput))
          as SetPrincipalTagAttributeMapInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetPrincipalTagAttributeMapInput create() =>
      SetPrincipalTagAttributeMapInput._();
  @$core.override
  SetPrincipalTagAttributeMapInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetPrincipalTagAttributeMapInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetPrincipalTagAttributeMapInput>(
          create);
  static SetPrincipalTagAttributeMapInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(306413999)
  $core.bool get usedefaults => $_getBF(1);
  @$pb.TagNumber(306413999)
  set usedefaults($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(306413999)
  $core.bool hasUsedefaults() => $_has(1);
  @$pb.TagNumber(306413999)
  void clearUsedefaults() => $_clearField(306413999);

  @$pb.TagNumber(346698229)
  $pb.PbMap<$core.String, $core.String> get principaltags => $_getMap(2);

  @$pb.TagNumber(419175410)
  $core.String get identityprovidername => $_getSZ(3);
  @$pb.TagNumber(419175410)
  set identityprovidername($core.String value) => $_setString(3, value);
  @$pb.TagNumber(419175410)
  $core.bool hasIdentityprovidername() => $_has(3);
  @$pb.TagNumber(419175410)
  void clearIdentityprovidername() => $_clearField(419175410);
}

class SetPrincipalTagAttributeMapResponse extends $pb.GeneratedMessage {
  factory SetPrincipalTagAttributeMapResponse({
    $core.String? identitypoolid,
    $core.bool? usedefaults,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? principaltags,
    $core.String? identityprovidername,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (usedefaults != null) result.usedefaults = usedefaults;
    if (principaltags != null) result.principaltags.addEntries(principaltags);
    if (identityprovidername != null)
      result.identityprovidername = identityprovidername;
    return result;
  }

  SetPrincipalTagAttributeMapResponse._();

  factory SetPrincipalTagAttributeMapResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetPrincipalTagAttributeMapResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetPrincipalTagAttributeMapResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOB(306413999, _omitFieldNames ? '' : 'usedefaults')
    ..m<$core.String, $core.String>(
        346698229, _omitFieldNames ? '' : 'principaltags',
        entryClassName:
            'SetPrincipalTagAttributeMapResponse.PrincipaltagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(419175410, _omitFieldNames ? '' : 'identityprovidername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetPrincipalTagAttributeMapResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetPrincipalTagAttributeMapResponse copyWith(
          void Function(SetPrincipalTagAttributeMapResponse) updates) =>
      super.copyWith((message) =>
              updates(message as SetPrincipalTagAttributeMapResponse))
          as SetPrincipalTagAttributeMapResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetPrincipalTagAttributeMapResponse create() =>
      SetPrincipalTagAttributeMapResponse._();
  @$core.override
  SetPrincipalTagAttributeMapResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetPrincipalTagAttributeMapResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          SetPrincipalTagAttributeMapResponse>(create);
  static SetPrincipalTagAttributeMapResponse? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(306413999)
  $core.bool get usedefaults => $_getBF(1);
  @$pb.TagNumber(306413999)
  set usedefaults($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(306413999)
  $core.bool hasUsedefaults() => $_has(1);
  @$pb.TagNumber(306413999)
  void clearUsedefaults() => $_clearField(306413999);

  @$pb.TagNumber(346698229)
  $pb.PbMap<$core.String, $core.String> get principaltags => $_getMap(2);

  @$pb.TagNumber(419175410)
  $core.String get identityprovidername => $_getSZ(3);
  @$pb.TagNumber(419175410)
  set identityprovidername($core.String value) => $_setString(3, value);
  @$pb.TagNumber(419175410)
  $core.bool hasIdentityprovidername() => $_has(3);
  @$pb.TagNumber(419175410)
  void clearIdentityprovidername() => $_clearField(419175410);
}

class TagResourceInput extends $pb.GeneratedMessage {
  factory TagResourceInput({
    $core.String? resourcearn,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (tags != null) result.tags.addEntries(tags);
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
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'TagResourceInput.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

  @$pb.TagNumber(381526209)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(1);
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
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
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

class TooManyRequestsException extends $pb.GeneratedMessage {
  factory TooManyRequestsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyRequestsException._();

  factory TooManyRequestsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyRequestsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyRequestsException',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyRequestsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyRequestsException copyWith(
          void Function(TooManyRequestsException) updates) =>
      super.copyWith((message) => updates(message as TooManyRequestsException))
          as TooManyRequestsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyRequestsException create() => TooManyRequestsException._();
  @$core.override
  TooManyRequestsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyRequestsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyRequestsException>(create);
  static TooManyRequestsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class UnlinkDeveloperIdentityInput extends $pb.GeneratedMessage {
  factory UnlinkDeveloperIdentityInput({
    $core.String? identitypoolid,
    $core.String? identityid,
    $core.String? developeruseridentifier,
    $core.String? developerprovidername,
  }) {
    final result = create();
    if (identitypoolid != null) result.identitypoolid = identitypoolid;
    if (identityid != null) result.identityid = identityid;
    if (developeruseridentifier != null)
      result.developeruseridentifier = developeruseridentifier;
    if (developerprovidername != null)
      result.developerprovidername = developerprovidername;
    return result;
  }

  UnlinkDeveloperIdentityInput._();

  factory UnlinkDeveloperIdentityInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnlinkDeveloperIdentityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnlinkDeveloperIdentityInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aOS(23936765, _omitFieldNames ? '' : 'identitypoolid')
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..aOS(349826618, _omitFieldNames ? '' : 'developeruseridentifier')
    ..aOS(517589792, _omitFieldNames ? '' : 'developerprovidername')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnlinkDeveloperIdentityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnlinkDeveloperIdentityInput copyWith(
          void Function(UnlinkDeveloperIdentityInput) updates) =>
      super.copyWith(
              (message) => updates(message as UnlinkDeveloperIdentityInput))
          as UnlinkDeveloperIdentityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnlinkDeveloperIdentityInput create() =>
      UnlinkDeveloperIdentityInput._();
  @$core.override
  UnlinkDeveloperIdentityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnlinkDeveloperIdentityInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnlinkDeveloperIdentityInput>(create);
  static UnlinkDeveloperIdentityInput? _defaultInstance;

  @$pb.TagNumber(23936765)
  $core.String get identitypoolid => $_getSZ(0);
  @$pb.TagNumber(23936765)
  set identitypoolid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23936765)
  $core.bool hasIdentitypoolid() => $_has(0);
  @$pb.TagNumber(23936765)
  void clearIdentitypoolid() => $_clearField(23936765);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(1);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(1);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);

  @$pb.TagNumber(349826618)
  $core.String get developeruseridentifier => $_getSZ(2);
  @$pb.TagNumber(349826618)
  set developeruseridentifier($core.String value) => $_setString(2, value);
  @$pb.TagNumber(349826618)
  $core.bool hasDeveloperuseridentifier() => $_has(2);
  @$pb.TagNumber(349826618)
  void clearDeveloperuseridentifier() => $_clearField(349826618);

  @$pb.TagNumber(517589792)
  $core.String get developerprovidername => $_getSZ(3);
  @$pb.TagNumber(517589792)
  set developerprovidername($core.String value) => $_setString(3, value);
  @$pb.TagNumber(517589792)
  $core.bool hasDeveloperprovidername() => $_has(3);
  @$pb.TagNumber(517589792)
  void clearDeveloperprovidername() => $_clearField(517589792);
}

class UnlinkIdentityInput extends $pb.GeneratedMessage {
  factory UnlinkIdentityInput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? logins,
    $core.String? identityid,
    $core.Iterable<$core.String>? loginstoremove,
  }) {
    final result = create();
    if (logins != null) result.logins.addEntries(logins);
    if (identityid != null) result.identityid = identityid;
    if (loginstoremove != null) result.loginstoremove.addAll(loginstoremove);
    return result;
  }

  UnlinkIdentityInput._();

  factory UnlinkIdentityInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnlinkIdentityInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnlinkIdentityInput',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(109702772, _omitFieldNames ? '' : 'logins',
        entryClassName: 'UnlinkIdentityInput.LoginsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('cognito_identity'))
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..pPS(490649531, _omitFieldNames ? '' : 'loginstoremove')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnlinkIdentityInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnlinkIdentityInput copyWith(void Function(UnlinkIdentityInput) updates) =>
      super.copyWith((message) => updates(message as UnlinkIdentityInput))
          as UnlinkIdentityInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnlinkIdentityInput create() => UnlinkIdentityInput._();
  @$core.override
  UnlinkIdentityInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnlinkIdentityInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnlinkIdentityInput>(create);
  static UnlinkIdentityInput? _defaultInstance;

  @$pb.TagNumber(109702772)
  $pb.PbMap<$core.String, $core.String> get logins => $_getMap(0);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(1);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(1);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);

  @$pb.TagNumber(490649531)
  $pb.PbList<$core.String> get loginstoremove => $_getList(2);
}

class UnprocessedIdentityId extends $pb.GeneratedMessage {
  factory UnprocessedIdentityId({
    ErrorCode? errorcode,
    $core.String? identityid,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (identityid != null) result.identityid = identityid;
    return result;
  }

  UnprocessedIdentityId._();

  factory UnprocessedIdentityId.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnprocessedIdentityId.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnprocessedIdentityId',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..aE<ErrorCode>(34663193, _omitFieldNames ? '' : 'errorcode',
        enumValues: ErrorCode.values)
    ..aOS(234187223, _omitFieldNames ? '' : 'identityid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnprocessedIdentityId clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnprocessedIdentityId copyWith(
          void Function(UnprocessedIdentityId) updates) =>
      super.copyWith((message) => updates(message as UnprocessedIdentityId))
          as UnprocessedIdentityId;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnprocessedIdentityId create() => UnprocessedIdentityId._();
  @$core.override
  UnprocessedIdentityId createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnprocessedIdentityId getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnprocessedIdentityId>(create);
  static UnprocessedIdentityId? _defaultInstance;

  @$pb.TagNumber(34663193)
  ErrorCode get errorcode => $_getN(0);
  @$pb.TagNumber(34663193)
  set errorcode(ErrorCode value) => $_setField(34663193, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(234187223)
  $core.String get identityid => $_getSZ(1);
  @$pb.TagNumber(234187223)
  set identityid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(234187223)
  $core.bool hasIdentityid() => $_has(1);
  @$pb.TagNumber(234187223)
  void clearIdentityid() => $_clearField(234187223);
}

class UntagResourceInput extends $pb.GeneratedMessage {
  factory UntagResourceInput({
    $core.Iterable<$core.String>? tagkeys,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
    if (resourcearn != null) result.resourcearn = resourcearn;
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
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'cognito_identity'),
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

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
