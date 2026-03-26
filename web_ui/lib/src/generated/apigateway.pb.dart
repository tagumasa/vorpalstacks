// This is a generated file - do not edit.
//
// Generated from apigateway.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'apigateway.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'apigateway.pbenum.dart';

class AccessLogSettings extends $pb.GeneratedMessage {
  factory AccessLogSettings({
    $core.String? destinationarn,
    $core.String? format,
  }) {
    final result = create();
    if (destinationarn != null) result.destinationarn = destinationarn;
    if (format != null) result.format = format;
    return result;
  }

  AccessLogSettings._();

  factory AccessLogSettings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccessLogSettings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccessLogSettings',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(427601315, _omitFieldNames ? '' : 'destinationarn')
    ..aOS(429753683, _omitFieldNames ? '' : 'format')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccessLogSettings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccessLogSettings copyWith(void Function(AccessLogSettings) updates) =>
      super.copyWith((message) => updates(message as AccessLogSettings))
          as AccessLogSettings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccessLogSettings create() => AccessLogSettings._();
  @$core.override
  AccessLogSettings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccessLogSettings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccessLogSettings>(create);
  static AccessLogSettings? _defaultInstance;

  @$pb.TagNumber(427601315)
  $core.String get destinationarn => $_getSZ(0);
  @$pb.TagNumber(427601315)
  set destinationarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(427601315)
  $core.bool hasDestinationarn() => $_has(0);
  @$pb.TagNumber(427601315)
  void clearDestinationarn() => $_clearField(427601315);

  @$pb.TagNumber(429753683)
  $core.String get format => $_getSZ(1);
  @$pb.TagNumber(429753683)
  set format($core.String value) => $_setString(1, value);
  @$pb.TagNumber(429753683)
  $core.bool hasFormat() => $_has(1);
  @$pb.TagNumber(429753683)
  void clearFormat() => $_clearField(429753683);
}

class Account extends $pb.GeneratedMessage {
  factory Account({
    $core.String? cloudwatchrolearn,
    $core.String? apikeyversion,
    ThrottleSettings? throttlesettings,
    $core.Iterable<$core.String>? features,
  }) {
    final result = create();
    if (cloudwatchrolearn != null) result.cloudwatchrolearn = cloudwatchrolearn;
    if (apikeyversion != null) result.apikeyversion = apikeyversion;
    if (throttlesettings != null) result.throttlesettings = throttlesettings;
    if (features != null) result.features.addAll(features);
    return result;
  }

  Account._();

  factory Account.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Account.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Account',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(144858329, _omitFieldNames ? '' : 'cloudwatchrolearn')
    ..aOS(148957539, _omitFieldNames ? '' : 'apikeyversion')
    ..aOM<ThrottleSettings>(
        165000097, _omitFieldNames ? '' : 'throttlesettings',
        subBuilder: ThrottleSettings.create)
    ..pPS(528712651, _omitFieldNames ? '' : 'features')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Account clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Account copyWith(void Function(Account) updates) =>
      super.copyWith((message) => updates(message as Account)) as Account;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Account create() => Account._();
  @$core.override
  Account createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Account getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Account>(create);
  static Account? _defaultInstance;

  @$pb.TagNumber(144858329)
  $core.String get cloudwatchrolearn => $_getSZ(0);
  @$pb.TagNumber(144858329)
  set cloudwatchrolearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(144858329)
  $core.bool hasCloudwatchrolearn() => $_has(0);
  @$pb.TagNumber(144858329)
  void clearCloudwatchrolearn() => $_clearField(144858329);

  @$pb.TagNumber(148957539)
  $core.String get apikeyversion => $_getSZ(1);
  @$pb.TagNumber(148957539)
  set apikeyversion($core.String value) => $_setString(1, value);
  @$pb.TagNumber(148957539)
  $core.bool hasApikeyversion() => $_has(1);
  @$pb.TagNumber(148957539)
  void clearApikeyversion() => $_clearField(148957539);

  @$pb.TagNumber(165000097)
  ThrottleSettings get throttlesettings => $_getN(2);
  @$pb.TagNumber(165000097)
  set throttlesettings(ThrottleSettings value) => $_setField(165000097, value);
  @$pb.TagNumber(165000097)
  $core.bool hasThrottlesettings() => $_has(2);
  @$pb.TagNumber(165000097)
  void clearThrottlesettings() => $_clearField(165000097);
  @$pb.TagNumber(165000097)
  ThrottleSettings ensureThrottlesettings() => $_ensure(2);

  @$pb.TagNumber(528712651)
  $pb.PbList<$core.String> get features => $_getList(3);
}

class ApiKey extends $pb.GeneratedMessage {
  factory ApiKey({
    $core.String? value,
    $core.bool? enabled,
    $core.String? createddate,
    $core.String? name,
    $core.String? customerid,
    $core.Iterable<$core.String>? stagekeys,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    $core.String? id,
    $core.String? lastupdateddate,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (enabled != null) result.enabled = enabled;
    if (createddate != null) result.createddate = createddate;
    if (name != null) result.name = name;
    if (customerid != null) result.customerid = customerid;
    if (stagekeys != null) result.stagekeys.addAll(stagekeys);
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (id != null) result.id = id;
    if (lastupdateddate != null) result.lastupdateddate = lastupdateddate;
    return result;
  }

  ApiKey._();

  factory ApiKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ApiKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApiKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(39769035, _omitFieldNames ? '' : 'value')
    ..aOB(49525663, _omitFieldNames ? '' : 'enabled')
    ..aOS(53061200, _omitFieldNames ? '' : 'createddate')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(227830269, _omitFieldNames ? '' : 'customerid')
    ..pPS(287991830, _omitFieldNames ? '' : 'stagekeys')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'ApiKey.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..aOS(448453361, _omitFieldNames ? '' : 'lastupdateddate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiKey copyWith(void Function(ApiKey) updates) =>
      super.copyWith((message) => updates(message as ApiKey)) as ApiKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ApiKey create() => ApiKey._();
  @$core.override
  ApiKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ApiKey getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ApiKey>(create);
  static ApiKey? _defaultInstance;

  @$pb.TagNumber(39769035)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(39769035)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39769035)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(39769035)
  void clearValue() => $_clearField(39769035);

  @$pb.TagNumber(49525663)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(49525663)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(49525663)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(49525663)
  void clearEnabled() => $_clearField(49525663);

  @$pb.TagNumber(53061200)
  $core.String get createddate => $_getSZ(2);
  @$pb.TagNumber(53061200)
  set createddate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(53061200)
  $core.bool hasCreateddate() => $_has(2);
  @$pb.TagNumber(53061200)
  void clearCreateddate() => $_clearField(53061200);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(227830269)
  $core.String get customerid => $_getSZ(4);
  @$pb.TagNumber(227830269)
  set customerid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(227830269)
  $core.bool hasCustomerid() => $_has(4);
  @$pb.TagNumber(227830269)
  void clearCustomerid() => $_clearField(227830269);

  @$pb.TagNumber(287991830)
  $pb.PbList<$core.String> get stagekeys => $_getList(5);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(6);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(7);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(7, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(7);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(8);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(8, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(8);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);

  @$pb.TagNumber(448453361)
  $core.String get lastupdateddate => $_getSZ(9);
  @$pb.TagNumber(448453361)
  set lastupdateddate($core.String value) => $_setString(9, value);
  @$pb.TagNumber(448453361)
  $core.bool hasLastupdateddate() => $_has(9);
  @$pb.TagNumber(448453361)
  void clearLastupdateddate() => $_clearField(448453361);
}

class ApiKeyIds extends $pb.GeneratedMessage {
  factory ApiKeyIds({
    $core.Iterable<$core.String>? ids,
    $core.Iterable<$core.String>? warnings,
  }) {
    final result = create();
    if (ids != null) result.ids.addAll(ids);
    if (warnings != null) result.warnings.addAll(warnings);
    return result;
  }

  ApiKeyIds._();

  factory ApiKeyIds.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ApiKeyIds.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApiKeyIds',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(13616490, _omitFieldNames ? '' : 'ids')
    ..pPS(185617301, _omitFieldNames ? '' : 'warnings')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiKeyIds clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiKeyIds copyWith(void Function(ApiKeyIds) updates) =>
      super.copyWith((message) => updates(message as ApiKeyIds)) as ApiKeyIds;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ApiKeyIds create() => ApiKeyIds._();
  @$core.override
  ApiKeyIds createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ApiKeyIds getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ApiKeyIds>(create);
  static ApiKeyIds? _defaultInstance;

  @$pb.TagNumber(13616490)
  $pb.PbList<$core.String> get ids => $_getList(0);

  @$pb.TagNumber(185617301)
  $pb.PbList<$core.String> get warnings => $_getList(1);
}

class ApiKeys extends $pb.GeneratedMessage {
  factory ApiKeys({
    $core.Iterable<$core.String>? warnings,
    $core.String? position,
    $core.Iterable<ApiKey>? items,
  }) {
    final result = create();
    if (warnings != null) result.warnings.addAll(warnings);
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  ApiKeys._();

  factory ApiKeys.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ApiKeys.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApiKeys',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(185617301, _omitFieldNames ? '' : 'warnings')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<ApiKey>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: ApiKey.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiKeys clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiKeys copyWith(void Function(ApiKeys) updates) =>
      super.copyWith((message) => updates(message as ApiKeys)) as ApiKeys;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ApiKeys create() => ApiKeys._();
  @$core.override
  ApiKeys createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ApiKeys getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ApiKeys>(create);
  static ApiKeys? _defaultInstance;

  @$pb.TagNumber(185617301)
  $pb.PbList<$core.String> get warnings => $_getList(0);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<ApiKey> get items => $_getList(2);
}

class ApiStage extends $pb.GeneratedMessage {
  factory ApiStage({
    $core.String? apiid,
    $core.String? stage,
    $core.Iterable<$core.MapEntry<$core.String, ThrottleSettings>>? throttle,
  }) {
    final result = create();
    if (apiid != null) result.apiid = apiid;
    if (stage != null) result.stage = stage;
    if (throttle != null) result.throttle.addEntries(throttle);
    return result;
  }

  ApiStage._();

  factory ApiStage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ApiStage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApiStage',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(113380971, _omitFieldNames ? '' : 'apiid')
    ..aOS(140155438, _omitFieldNames ? '' : 'stage')
    ..m<$core.String, ThrottleSettings>(
        395260638, _omitFieldNames ? '' : 'throttle',
        entryClassName: 'ApiStage.ThrottleEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: ThrottleSettings.create,
        valueDefaultOrMaker: ThrottleSettings.getDefault,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiStage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiStage copyWith(void Function(ApiStage) updates) =>
      super.copyWith((message) => updates(message as ApiStage)) as ApiStage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ApiStage create() => ApiStage._();
  @$core.override
  ApiStage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ApiStage getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ApiStage>(create);
  static ApiStage? _defaultInstance;

  @$pb.TagNumber(113380971)
  $core.String get apiid => $_getSZ(0);
  @$pb.TagNumber(113380971)
  set apiid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(113380971)
  $core.bool hasApiid() => $_has(0);
  @$pb.TagNumber(113380971)
  void clearApiid() => $_clearField(113380971);

  @$pb.TagNumber(140155438)
  $core.String get stage => $_getSZ(1);
  @$pb.TagNumber(140155438)
  set stage($core.String value) => $_setString(1, value);
  @$pb.TagNumber(140155438)
  $core.bool hasStage() => $_has(1);
  @$pb.TagNumber(140155438)
  void clearStage() => $_clearField(140155438);

  @$pb.TagNumber(395260638)
  $pb.PbMap<$core.String, ThrottleSettings> get throttle => $_getMap(2);
}

class Authorizer extends $pb.GeneratedMessage {
  factory Authorizer({
    $core.int? authorizerresultttlinseconds,
    $core.String? authtype,
    $core.String? name,
    $core.String? identityvalidationexpression,
    $core.String? authorizercredentials,
    $core.String? identitysource,
    AuthorizerType? type,
    $core.Iterable<$core.String>? providerarns,
    $core.String? id,
    $core.String? authorizeruri,
  }) {
    final result = create();
    if (authorizerresultttlinseconds != null)
      result.authorizerresultttlinseconds = authorizerresultttlinseconds;
    if (authtype != null) result.authtype = authtype;
    if (name != null) result.name = name;
    if (identityvalidationexpression != null)
      result.identityvalidationexpression = identityvalidationexpression;
    if (authorizercredentials != null)
      result.authorizercredentials = authorizercredentials;
    if (identitysource != null) result.identitysource = identitysource;
    if (type != null) result.type = type;
    if (providerarns != null) result.providerarns.addAll(providerarns);
    if (id != null) result.id = id;
    if (authorizeruri != null) result.authorizeruri = authorizeruri;
    return result;
  }

  Authorizer._();

  factory Authorizer.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Authorizer.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Authorizer',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(135440208, _omitFieldNames ? '' : 'authorizerresultttlinseconds')
    ..aOS(162773848, _omitFieldNames ? '' : 'authtype')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(227211199, _omitFieldNames ? '' : 'identityvalidationexpression')
    ..aOS(233575233, _omitFieldNames ? '' : 'authorizercredentials')
    ..aOS(285615231, _omitFieldNames ? '' : 'identitysource')
    ..aE<AuthorizerType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: AuthorizerType.values)
    ..pPS(301486689, _omitFieldNames ? '' : 'providerarns')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..aOS(525146137, _omitFieldNames ? '' : 'authorizeruri')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Authorizer clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Authorizer copyWith(void Function(Authorizer) updates) =>
      super.copyWith((message) => updates(message as Authorizer)) as Authorizer;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Authorizer create() => Authorizer._();
  @$core.override
  Authorizer createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Authorizer getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Authorizer>(create);
  static Authorizer? _defaultInstance;

  @$pb.TagNumber(135440208)
  $core.int get authorizerresultttlinseconds => $_getIZ(0);
  @$pb.TagNumber(135440208)
  set authorizerresultttlinseconds($core.int value) =>
      $_setSignedInt32(0, value);
  @$pb.TagNumber(135440208)
  $core.bool hasAuthorizerresultttlinseconds() => $_has(0);
  @$pb.TagNumber(135440208)
  void clearAuthorizerresultttlinseconds() => $_clearField(135440208);

  @$pb.TagNumber(162773848)
  $core.String get authtype => $_getSZ(1);
  @$pb.TagNumber(162773848)
  set authtype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(162773848)
  $core.bool hasAuthtype() => $_has(1);
  @$pb.TagNumber(162773848)
  void clearAuthtype() => $_clearField(162773848);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(227211199)
  $core.String get identityvalidationexpression => $_getSZ(3);
  @$pb.TagNumber(227211199)
  set identityvalidationexpression($core.String value) => $_setString(3, value);
  @$pb.TagNumber(227211199)
  $core.bool hasIdentityvalidationexpression() => $_has(3);
  @$pb.TagNumber(227211199)
  void clearIdentityvalidationexpression() => $_clearField(227211199);

  @$pb.TagNumber(233575233)
  $core.String get authorizercredentials => $_getSZ(4);
  @$pb.TagNumber(233575233)
  set authorizercredentials($core.String value) => $_setString(4, value);
  @$pb.TagNumber(233575233)
  $core.bool hasAuthorizercredentials() => $_has(4);
  @$pb.TagNumber(233575233)
  void clearAuthorizercredentials() => $_clearField(233575233);

  @$pb.TagNumber(285615231)
  $core.String get identitysource => $_getSZ(5);
  @$pb.TagNumber(285615231)
  set identitysource($core.String value) => $_setString(5, value);
  @$pb.TagNumber(285615231)
  $core.bool hasIdentitysource() => $_has(5);
  @$pb.TagNumber(285615231)
  void clearIdentitysource() => $_clearField(285615231);

  @$pb.TagNumber(287830350)
  AuthorizerType get type => $_getN(6);
  @$pb.TagNumber(287830350)
  set type(AuthorizerType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(6);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(301486689)
  $pb.PbList<$core.String> get providerarns => $_getList(7);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(8);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(8, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(8);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);

  @$pb.TagNumber(525146137)
  $core.String get authorizeruri => $_getSZ(9);
  @$pb.TagNumber(525146137)
  set authorizeruri($core.String value) => $_setString(9, value);
  @$pb.TagNumber(525146137)
  $core.bool hasAuthorizeruri() => $_has(9);
  @$pb.TagNumber(525146137)
  void clearAuthorizeruri() => $_clearField(525146137);
}

class Authorizers extends $pb.GeneratedMessage {
  factory Authorizers({
    $core.String? position,
    $core.Iterable<Authorizer>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  Authorizers._();

  factory Authorizers.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Authorizers.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Authorizers',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<Authorizer>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: Authorizer.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Authorizers clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Authorizers copyWith(void Function(Authorizers) updates) =>
      super.copyWith((message) => updates(message as Authorizers))
          as Authorizers;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Authorizers create() => Authorizers._();
  @$core.override
  Authorizers createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Authorizers getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Authorizers>(create);
  static Authorizers? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<Authorizer> get items => $_getList(1);
}

class BadRequestException extends $pb.GeneratedMessage {
  factory BadRequestException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  BadRequestException._();

  factory BadRequestException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BadRequestException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BadRequestException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BadRequestException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BadRequestException copyWith(void Function(BadRequestException) updates) =>
      super.copyWith((message) => updates(message as BadRequestException))
          as BadRequestException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BadRequestException create() => BadRequestException._();
  @$core.override
  BadRequestException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BadRequestException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BadRequestException>(create);
  static BadRequestException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class BasePathMapping extends $pb.GeneratedMessage {
  factory BasePathMapping({
    $core.String? stage,
    $core.String? basepath,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stage != null) result.stage = stage;
    if (basepath != null) result.basepath = basepath;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  BasePathMapping._();

  factory BasePathMapping.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BasePathMapping.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BasePathMapping',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(140155438, _omitFieldNames ? '' : 'stage')
    ..aOS(267528880, _omitFieldNames ? '' : 'basepath')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BasePathMapping clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BasePathMapping copyWith(void Function(BasePathMapping) updates) =>
      super.copyWith((message) => updates(message as BasePathMapping))
          as BasePathMapping;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BasePathMapping create() => BasePathMapping._();
  @$core.override
  BasePathMapping createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BasePathMapping getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BasePathMapping>(create);
  static BasePathMapping? _defaultInstance;

  @$pb.TagNumber(140155438)
  $core.String get stage => $_getSZ(0);
  @$pb.TagNumber(140155438)
  set stage($core.String value) => $_setString(0, value);
  @$pb.TagNumber(140155438)
  $core.bool hasStage() => $_has(0);
  @$pb.TagNumber(140155438)
  void clearStage() => $_clearField(140155438);

  @$pb.TagNumber(267528880)
  $core.String get basepath => $_getSZ(1);
  @$pb.TagNumber(267528880)
  set basepath($core.String value) => $_setString(1, value);
  @$pb.TagNumber(267528880)
  $core.bool hasBasepath() => $_has(1);
  @$pb.TagNumber(267528880)
  void clearBasepath() => $_clearField(267528880);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class BasePathMappings extends $pb.GeneratedMessage {
  factory BasePathMappings({
    $core.String? position,
    $core.Iterable<BasePathMapping>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  BasePathMappings._();

  factory BasePathMappings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BasePathMappings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BasePathMappings',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<BasePathMapping>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: BasePathMapping.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BasePathMappings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BasePathMappings copyWith(void Function(BasePathMappings) updates) =>
      super.copyWith((message) => updates(message as BasePathMappings))
          as BasePathMappings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BasePathMappings create() => BasePathMappings._();
  @$core.override
  BasePathMappings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BasePathMappings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BasePathMappings>(create);
  static BasePathMappings? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<BasePathMapping> get items => $_getList(1);
}

class CanarySettings extends $pb.GeneratedMessage {
  factory CanarySettings({
    $core.double? percenttraffic,
    $core.bool? usestagecache,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        stagevariableoverrides,
    $core.String? deploymentid,
  }) {
    final result = create();
    if (percenttraffic != null) result.percenttraffic = percenttraffic;
    if (usestagecache != null) result.usestagecache = usestagecache;
    if (stagevariableoverrides != null)
      result.stagevariableoverrides.addEntries(stagevariableoverrides);
    if (deploymentid != null) result.deploymentid = deploymentid;
    return result;
  }

  CanarySettings._();

  factory CanarySettings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CanarySettings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CanarySettings',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aD(147177864, _omitFieldNames ? '' : 'percenttraffic')
    ..aOB(179841697, _omitFieldNames ? '' : 'usestagecache')
    ..m<$core.String, $core.String>(
        221124259, _omitFieldNames ? '' : 'stagevariableoverrides',
        entryClassName: 'CanarySettings.StagevariableoverridesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(439369188, _omitFieldNames ? '' : 'deploymentid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CanarySettings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CanarySettings copyWith(void Function(CanarySettings) updates) =>
      super.copyWith((message) => updates(message as CanarySettings))
          as CanarySettings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CanarySettings create() => CanarySettings._();
  @$core.override
  CanarySettings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CanarySettings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CanarySettings>(create);
  static CanarySettings? _defaultInstance;

  @$pb.TagNumber(147177864)
  $core.double get percenttraffic => $_getN(0);
  @$pb.TagNumber(147177864)
  set percenttraffic($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(147177864)
  $core.bool hasPercenttraffic() => $_has(0);
  @$pb.TagNumber(147177864)
  void clearPercenttraffic() => $_clearField(147177864);

  @$pb.TagNumber(179841697)
  $core.bool get usestagecache => $_getBF(1);
  @$pb.TagNumber(179841697)
  set usestagecache($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(179841697)
  $core.bool hasUsestagecache() => $_has(1);
  @$pb.TagNumber(179841697)
  void clearUsestagecache() => $_clearField(179841697);

  @$pb.TagNumber(221124259)
  $pb.PbMap<$core.String, $core.String> get stagevariableoverrides =>
      $_getMap(2);

  @$pb.TagNumber(439369188)
  $core.String get deploymentid => $_getSZ(3);
  @$pb.TagNumber(439369188)
  set deploymentid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(439369188)
  $core.bool hasDeploymentid() => $_has(3);
  @$pb.TagNumber(439369188)
  void clearDeploymentid() => $_clearField(439369188);
}

class ClientCertificate extends $pb.GeneratedMessage {
  factory ClientCertificate({
    $core.String? createddate,
    $core.String? clientcertificateid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    $core.String? expirationdate,
    $core.String? pemencodedcertificate,
  }) {
    final result = create();
    if (createddate != null) result.createddate = createddate;
    if (clientcertificateid != null)
      result.clientcertificateid = clientcertificateid;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (expirationdate != null) result.expirationdate = expirationdate;
    if (pemencodedcertificate != null)
      result.pemencodedcertificate = pemencodedcertificate;
    return result;
  }

  ClientCertificate._();

  factory ClientCertificate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ClientCertificate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ClientCertificate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(53061200, _omitFieldNames ? '' : 'createddate')
    ..aOS(276222909, _omitFieldNames ? '' : 'clientcertificateid')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'ClientCertificate.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(391033245, _omitFieldNames ? '' : 'expirationdate')
    ..aOS(449484817, _omitFieldNames ? '' : 'pemencodedcertificate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientCertificate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientCertificate copyWith(void Function(ClientCertificate) updates) =>
      super.copyWith((message) => updates(message as ClientCertificate))
          as ClientCertificate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ClientCertificate create() => ClientCertificate._();
  @$core.override
  ClientCertificate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ClientCertificate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ClientCertificate>(create);
  static ClientCertificate? _defaultInstance;

  @$pb.TagNumber(53061200)
  $core.String get createddate => $_getSZ(0);
  @$pb.TagNumber(53061200)
  set createddate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(53061200)
  $core.bool hasCreateddate() => $_has(0);
  @$pb.TagNumber(53061200)
  void clearCreateddate() => $_clearField(53061200);

  @$pb.TagNumber(276222909)
  $core.String get clientcertificateid => $_getSZ(1);
  @$pb.TagNumber(276222909)
  set clientcertificateid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(276222909)
  $core.bool hasClientcertificateid() => $_has(1);
  @$pb.TagNumber(276222909)
  void clearClientcertificateid() => $_clearField(276222909);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(2);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(391033245)
  $core.String get expirationdate => $_getSZ(4);
  @$pb.TagNumber(391033245)
  set expirationdate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(391033245)
  $core.bool hasExpirationdate() => $_has(4);
  @$pb.TagNumber(391033245)
  void clearExpirationdate() => $_clearField(391033245);

  @$pb.TagNumber(449484817)
  $core.String get pemencodedcertificate => $_getSZ(5);
  @$pb.TagNumber(449484817)
  set pemencodedcertificate($core.String value) => $_setString(5, value);
  @$pb.TagNumber(449484817)
  $core.bool hasPemencodedcertificate() => $_has(5);
  @$pb.TagNumber(449484817)
  void clearPemencodedcertificate() => $_clearField(449484817);
}

class ClientCertificates extends $pb.GeneratedMessage {
  factory ClientCertificates({
    $core.String? position,
    $core.Iterable<ClientCertificate>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  ClientCertificates._();

  factory ClientCertificates.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ClientCertificates.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ClientCertificates',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<ClientCertificate>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: ClientCertificate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientCertificates clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientCertificates copyWith(void Function(ClientCertificates) updates) =>
      super.copyWith((message) => updates(message as ClientCertificates))
          as ClientCertificates;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ClientCertificates create() => ClientCertificates._();
  @$core.override
  ClientCertificates createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ClientCertificates getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ClientCertificates>(create);
  static ClientCertificates? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<ClientCertificate> get items => $_getList(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
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

class CreateApiKeyRequest extends $pb.GeneratedMessage {
  factory CreateApiKeyRequest({
    $core.String? value,
    $core.bool? enabled,
    $core.bool? generatedistinctid,
    $core.String? name,
    $core.String? customerid,
    $core.Iterable<StageKey>? stagekeys,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (enabled != null) result.enabled = enabled;
    if (generatedistinctid != null)
      result.generatedistinctid = generatedistinctid;
    if (name != null) result.name = name;
    if (customerid != null) result.customerid = customerid;
    if (stagekeys != null) result.stagekeys.addAll(stagekeys);
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    return result;
  }

  CreateApiKeyRequest._();

  factory CreateApiKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateApiKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateApiKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(39769035, _omitFieldNames ? '' : 'value')
    ..aOB(49525663, _omitFieldNames ? '' : 'enabled')
    ..aOB(99833588, _omitFieldNames ? '' : 'generatedistinctid')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(227830269, _omitFieldNames ? '' : 'customerid')
    ..pPM<StageKey>(287991830, _omitFieldNames ? '' : 'stagekeys',
        subBuilder: StageKey.create)
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateApiKeyRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateApiKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateApiKeyRequest copyWith(void Function(CreateApiKeyRequest) updates) =>
      super.copyWith((message) => updates(message as CreateApiKeyRequest))
          as CreateApiKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateApiKeyRequest create() => CreateApiKeyRequest._();
  @$core.override
  CreateApiKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateApiKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateApiKeyRequest>(create);
  static CreateApiKeyRequest? _defaultInstance;

  @$pb.TagNumber(39769035)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(39769035)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39769035)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(39769035)
  void clearValue() => $_clearField(39769035);

  @$pb.TagNumber(49525663)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(49525663)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(49525663)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(49525663)
  void clearEnabled() => $_clearField(49525663);

  @$pb.TagNumber(99833588)
  $core.bool get generatedistinctid => $_getBF(2);
  @$pb.TagNumber(99833588)
  set generatedistinctid($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(99833588)
  $core.bool hasGeneratedistinctid() => $_has(2);
  @$pb.TagNumber(99833588)
  void clearGeneratedistinctid() => $_clearField(99833588);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(227830269)
  $core.String get customerid => $_getSZ(4);
  @$pb.TagNumber(227830269)
  set customerid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(227830269)
  $core.bool hasCustomerid() => $_has(4);
  @$pb.TagNumber(227830269)
  void clearCustomerid() => $_clearField(227830269);

  @$pb.TagNumber(287991830)
  $pb.PbList<StageKey> get stagekeys => $_getList(5);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(6);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(7);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(7, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(7);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);
}

class CreateAuthorizerRequest extends $pb.GeneratedMessage {
  factory CreateAuthorizerRequest({
    $core.int? authorizerresultttlinseconds,
    $core.String? authtype,
    $core.String? name,
    $core.String? identityvalidationexpression,
    $core.String? authorizercredentials,
    $core.String? identitysource,
    AuthorizerType? type,
    $core.Iterable<$core.String>? providerarns,
    $core.String? restapiid,
    $core.String? authorizeruri,
  }) {
    final result = create();
    if (authorizerresultttlinseconds != null)
      result.authorizerresultttlinseconds = authorizerresultttlinseconds;
    if (authtype != null) result.authtype = authtype;
    if (name != null) result.name = name;
    if (identityvalidationexpression != null)
      result.identityvalidationexpression = identityvalidationexpression;
    if (authorizercredentials != null)
      result.authorizercredentials = authorizercredentials;
    if (identitysource != null) result.identitysource = identitysource;
    if (type != null) result.type = type;
    if (providerarns != null) result.providerarns.addAll(providerarns);
    if (restapiid != null) result.restapiid = restapiid;
    if (authorizeruri != null) result.authorizeruri = authorizeruri;
    return result;
  }

  CreateAuthorizerRequest._();

  factory CreateAuthorizerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateAuthorizerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateAuthorizerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(135440208, _omitFieldNames ? '' : 'authorizerresultttlinseconds')
    ..aOS(162773848, _omitFieldNames ? '' : 'authtype')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(227211199, _omitFieldNames ? '' : 'identityvalidationexpression')
    ..aOS(233575233, _omitFieldNames ? '' : 'authorizercredentials')
    ..aOS(285615231, _omitFieldNames ? '' : 'identitysource')
    ..aE<AuthorizerType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: AuthorizerType.values)
    ..pPS(301486689, _omitFieldNames ? '' : 'providerarns')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(525146137, _omitFieldNames ? '' : 'authorizeruri')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAuthorizerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAuthorizerRequest copyWith(
          void Function(CreateAuthorizerRequest) updates) =>
      super.copyWith((message) => updates(message as CreateAuthorizerRequest))
          as CreateAuthorizerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateAuthorizerRequest create() => CreateAuthorizerRequest._();
  @$core.override
  CreateAuthorizerRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateAuthorizerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateAuthorizerRequest>(create);
  static CreateAuthorizerRequest? _defaultInstance;

  @$pb.TagNumber(135440208)
  $core.int get authorizerresultttlinseconds => $_getIZ(0);
  @$pb.TagNumber(135440208)
  set authorizerresultttlinseconds($core.int value) =>
      $_setSignedInt32(0, value);
  @$pb.TagNumber(135440208)
  $core.bool hasAuthorizerresultttlinseconds() => $_has(0);
  @$pb.TagNumber(135440208)
  void clearAuthorizerresultttlinseconds() => $_clearField(135440208);

  @$pb.TagNumber(162773848)
  $core.String get authtype => $_getSZ(1);
  @$pb.TagNumber(162773848)
  set authtype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(162773848)
  $core.bool hasAuthtype() => $_has(1);
  @$pb.TagNumber(162773848)
  void clearAuthtype() => $_clearField(162773848);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(227211199)
  $core.String get identityvalidationexpression => $_getSZ(3);
  @$pb.TagNumber(227211199)
  set identityvalidationexpression($core.String value) => $_setString(3, value);
  @$pb.TagNumber(227211199)
  $core.bool hasIdentityvalidationexpression() => $_has(3);
  @$pb.TagNumber(227211199)
  void clearIdentityvalidationexpression() => $_clearField(227211199);

  @$pb.TagNumber(233575233)
  $core.String get authorizercredentials => $_getSZ(4);
  @$pb.TagNumber(233575233)
  set authorizercredentials($core.String value) => $_setString(4, value);
  @$pb.TagNumber(233575233)
  $core.bool hasAuthorizercredentials() => $_has(4);
  @$pb.TagNumber(233575233)
  void clearAuthorizercredentials() => $_clearField(233575233);

  @$pb.TagNumber(285615231)
  $core.String get identitysource => $_getSZ(5);
  @$pb.TagNumber(285615231)
  set identitysource($core.String value) => $_setString(5, value);
  @$pb.TagNumber(285615231)
  $core.bool hasIdentitysource() => $_has(5);
  @$pb.TagNumber(285615231)
  void clearIdentitysource() => $_clearField(285615231);

  @$pb.TagNumber(287830350)
  AuthorizerType get type => $_getN(6);
  @$pb.TagNumber(287830350)
  set type(AuthorizerType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(6);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(301486689)
  $pb.PbList<$core.String> get providerarns => $_getList(7);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(8);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(8);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(525146137)
  $core.String get authorizeruri => $_getSZ(9);
  @$pb.TagNumber(525146137)
  set authorizeruri($core.String value) => $_setString(9, value);
  @$pb.TagNumber(525146137)
  $core.bool hasAuthorizeruri() => $_has(9);
  @$pb.TagNumber(525146137)
  void clearAuthorizeruri() => $_clearField(525146137);
}

class CreateBasePathMappingRequest extends $pb.GeneratedMessage {
  factory CreateBasePathMappingRequest({
    $core.String? stage,
    $core.String? basepath,
    $core.String? domainnameid,
    $core.String? restapiid,
    $core.String? domainname,
  }) {
    final result = create();
    if (stage != null) result.stage = stage;
    if (basepath != null) result.basepath = basepath;
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (restapiid != null) result.restapiid = restapiid;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  CreateBasePathMappingRequest._();

  factory CreateBasePathMappingRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateBasePathMappingRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateBasePathMappingRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(140155438, _omitFieldNames ? '' : 'stage')
    ..aOS(267528880, _omitFieldNames ? '' : 'basepath')
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBasePathMappingRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBasePathMappingRequest copyWith(
          void Function(CreateBasePathMappingRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateBasePathMappingRequest))
          as CreateBasePathMappingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateBasePathMappingRequest create() =>
      CreateBasePathMappingRequest._();
  @$core.override
  CreateBasePathMappingRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateBasePathMappingRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateBasePathMappingRequest>(create);
  static CreateBasePathMappingRequest? _defaultInstance;

  @$pb.TagNumber(140155438)
  $core.String get stage => $_getSZ(0);
  @$pb.TagNumber(140155438)
  set stage($core.String value) => $_setString(0, value);
  @$pb.TagNumber(140155438)
  $core.bool hasStage() => $_has(0);
  @$pb.TagNumber(140155438)
  void clearStage() => $_clearField(140155438);

  @$pb.TagNumber(267528880)
  $core.String get basepath => $_getSZ(1);
  @$pb.TagNumber(267528880)
  set basepath($core.String value) => $_setString(1, value);
  @$pb.TagNumber(267528880)
  $core.bool hasBasepath() => $_has(1);
  @$pb.TagNumber(267528880)
  void clearBasepath() => $_clearField(267528880);

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(2);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(2);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(4);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(4);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);
}

class CreateDeploymentRequest extends $pb.GeneratedMessage {
  factory CreateDeploymentRequest({
    $core.String? stagename,
    $core.bool? cacheclusterenabled,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? variables,
    CacheClusterSize? cacheclustersize,
    DeploymentCanarySettings? canarysettings,
    $core.String? description,
    $core.String? restapiid,
    $core.bool? tracingenabled,
    $core.String? stagedescription,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (cacheclusterenabled != null)
      result.cacheclusterenabled = cacheclusterenabled;
    if (variables != null) result.variables.addEntries(variables);
    if (cacheclustersize != null) result.cacheclustersize = cacheclustersize;
    if (canarysettings != null) result.canarysettings = canarysettings;
    if (description != null) result.description = description;
    if (restapiid != null) result.restapiid = restapiid;
    if (tracingenabled != null) result.tracingenabled = tracingenabled;
    if (stagedescription != null) result.stagedescription = stagedescription;
    return result;
  }

  CreateDeploymentRequest._();

  factory CreateDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..aOB(63967991, _omitFieldNames ? '' : 'cacheclusterenabled')
    ..m<$core.String, $core.String>(
        162226883, _omitFieldNames ? '' : 'variables',
        entryClassName: 'CreateDeploymentRequest.VariablesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aE<CacheClusterSize>(232189861, _omitFieldNames ? '' : 'cacheclustersize',
        enumValues: CacheClusterSize.values)
    ..aOM<DeploymentCanarySettings>(
        285544261, _omitFieldNames ? '' : 'canarysettings',
        subBuilder: DeploymentCanarySettings.create)
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOB(390995731, _omitFieldNames ? '' : 'tracingenabled')
    ..aOS(496169986, _omitFieldNames ? '' : 'stagedescription')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDeploymentRequest copyWith(
          void Function(CreateDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as CreateDeploymentRequest))
          as CreateDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDeploymentRequest create() => CreateDeploymentRequest._();
  @$core.override
  CreateDeploymentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDeploymentRequest>(create);
  static CreateDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(63967991)
  $core.bool get cacheclusterenabled => $_getBF(1);
  @$pb.TagNumber(63967991)
  set cacheclusterenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(63967991)
  $core.bool hasCacheclusterenabled() => $_has(1);
  @$pb.TagNumber(63967991)
  void clearCacheclusterenabled() => $_clearField(63967991);

  @$pb.TagNumber(162226883)
  $pb.PbMap<$core.String, $core.String> get variables => $_getMap(2);

  @$pb.TagNumber(232189861)
  CacheClusterSize get cacheclustersize => $_getN(3);
  @$pb.TagNumber(232189861)
  set cacheclustersize(CacheClusterSize value) => $_setField(232189861, value);
  @$pb.TagNumber(232189861)
  $core.bool hasCacheclustersize() => $_has(3);
  @$pb.TagNumber(232189861)
  void clearCacheclustersize() => $_clearField(232189861);

  @$pb.TagNumber(285544261)
  DeploymentCanarySettings get canarysettings => $_getN(4);
  @$pb.TagNumber(285544261)
  set canarysettings(DeploymentCanarySettings value) =>
      $_setField(285544261, value);
  @$pb.TagNumber(285544261)
  $core.bool hasCanarysettings() => $_has(4);
  @$pb.TagNumber(285544261)
  void clearCanarysettings() => $_clearField(285544261);
  @$pb.TagNumber(285544261)
  DeploymentCanarySettings ensureCanarysettings() => $_ensure(4);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(5);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(5, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(5);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(6);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(6);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(390995731)
  $core.bool get tracingenabled => $_getBF(7);
  @$pb.TagNumber(390995731)
  set tracingenabled($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(390995731)
  $core.bool hasTracingenabled() => $_has(7);
  @$pb.TagNumber(390995731)
  void clearTracingenabled() => $_clearField(390995731);

  @$pb.TagNumber(496169986)
  $core.String get stagedescription => $_getSZ(8);
  @$pb.TagNumber(496169986)
  set stagedescription($core.String value) => $_setString(8, value);
  @$pb.TagNumber(496169986)
  $core.bool hasStagedescription() => $_has(8);
  @$pb.TagNumber(496169986)
  void clearStagedescription() => $_clearField(496169986);
}

class CreateDocumentationPartRequest extends $pb.GeneratedMessage {
  factory CreateDocumentationPartRequest({
    DocumentationPartLocation? location,
    $core.String? properties,
    $core.String? restapiid,
  }) {
    final result = create();
    if (location != null) result.location = location;
    if (properties != null) result.properties = properties;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  CreateDocumentationPartRequest._();

  factory CreateDocumentationPartRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDocumentationPartRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDocumentationPartRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOM<DocumentationPartLocation>(
        200649127, _omitFieldNames ? '' : 'location',
        subBuilder: DocumentationPartLocation.create)
    ..aOS(299789533, _omitFieldNames ? '' : 'properties')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDocumentationPartRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDocumentationPartRequest copyWith(
          void Function(CreateDocumentationPartRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateDocumentationPartRequest))
          as CreateDocumentationPartRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDocumentationPartRequest create() =>
      CreateDocumentationPartRequest._();
  @$core.override
  CreateDocumentationPartRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDocumentationPartRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDocumentationPartRequest>(create);
  static CreateDocumentationPartRequest? _defaultInstance;

  @$pb.TagNumber(200649127)
  DocumentationPartLocation get location => $_getN(0);
  @$pb.TagNumber(200649127)
  set location(DocumentationPartLocation value) => $_setField(200649127, value);
  @$pb.TagNumber(200649127)
  $core.bool hasLocation() => $_has(0);
  @$pb.TagNumber(200649127)
  void clearLocation() => $_clearField(200649127);
  @$pb.TagNumber(200649127)
  DocumentationPartLocation ensureLocation() => $_ensure(0);

  @$pb.TagNumber(299789533)
  $core.String get properties => $_getSZ(1);
  @$pb.TagNumber(299789533)
  set properties($core.String value) => $_setString(1, value);
  @$pb.TagNumber(299789533)
  $core.bool hasProperties() => $_has(1);
  @$pb.TagNumber(299789533)
  void clearProperties() => $_clearField(299789533);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class CreateDocumentationVersionRequest extends $pb.GeneratedMessage {
  factory CreateDocumentationVersionRequest({
    $core.String? stagename,
    $core.String? documentationversion,
    $core.String? description,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (documentationversion != null)
      result.documentationversion = documentationversion;
    if (description != null) result.description = description;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  CreateDocumentationVersionRequest._();

  factory CreateDocumentationVersionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDocumentationVersionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDocumentationVersionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..aOS(167009804, _omitFieldNames ? '' : 'documentationversion')
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDocumentationVersionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDocumentationVersionRequest copyWith(
          void Function(CreateDocumentationVersionRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreateDocumentationVersionRequest))
          as CreateDocumentationVersionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDocumentationVersionRequest create() =>
      CreateDocumentationVersionRequest._();
  @$core.override
  CreateDocumentationVersionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDocumentationVersionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDocumentationVersionRequest>(
          create);
  static CreateDocumentationVersionRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(167009804)
  $core.String get documentationversion => $_getSZ(1);
  @$pb.TagNumber(167009804)
  set documentationversion($core.String value) => $_setString(1, value);
  @$pb.TagNumber(167009804)
  $core.bool hasDocumentationversion() => $_has(1);
  @$pb.TagNumber(167009804)
  void clearDocumentationversion() => $_clearField(167009804);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class CreateDomainNameAccessAssociationRequest extends $pb.GeneratedMessage {
  factory CreateDomainNameAccessAssociationRequest({
    AccessAssociationSourceType? accessassociationsourcetype,
    $core.String? domainnamearn,
    $core.String? accessassociationsource,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
  }) {
    final result = create();
    if (accessassociationsourcetype != null)
      result.accessassociationsourcetype = accessassociationsourcetype;
    if (domainnamearn != null) result.domainnamearn = domainnamearn;
    if (accessassociationsource != null)
      result.accessassociationsource = accessassociationsource;
    if (tags != null) result.tags.addEntries(tags);
    return result;
  }

  CreateDomainNameAccessAssociationRequest._();

  factory CreateDomainNameAccessAssociationRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDomainNameAccessAssociationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDomainNameAccessAssociationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aE<AccessAssociationSourceType>(
        176397628, _omitFieldNames ? '' : 'accessassociationsourcetype',
        enumValues: AccessAssociationSourceType.values)
    ..aOS(244019094, _omitFieldNames ? '' : 'domainnamearn')
    ..aOS(328257828, _omitFieldNames ? '' : 'accessassociationsource')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateDomainNameAccessAssociationRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDomainNameAccessAssociationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDomainNameAccessAssociationRequest copyWith(
          void Function(CreateDomainNameAccessAssociationRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreateDomainNameAccessAssociationRequest))
          as CreateDomainNameAccessAssociationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDomainNameAccessAssociationRequest create() =>
      CreateDomainNameAccessAssociationRequest._();
  @$core.override
  CreateDomainNameAccessAssociationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDomainNameAccessAssociationRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateDomainNameAccessAssociationRequest>(create);
  static CreateDomainNameAccessAssociationRequest? _defaultInstance;

  @$pb.TagNumber(176397628)
  AccessAssociationSourceType get accessassociationsourcetype => $_getN(0);
  @$pb.TagNumber(176397628)
  set accessassociationsourcetype(AccessAssociationSourceType value) =>
      $_setField(176397628, value);
  @$pb.TagNumber(176397628)
  $core.bool hasAccessassociationsourcetype() => $_has(0);
  @$pb.TagNumber(176397628)
  void clearAccessassociationsourcetype() => $_clearField(176397628);

  @$pb.TagNumber(244019094)
  $core.String get domainnamearn => $_getSZ(1);
  @$pb.TagNumber(244019094)
  set domainnamearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(244019094)
  $core.bool hasDomainnamearn() => $_has(1);
  @$pb.TagNumber(244019094)
  void clearDomainnamearn() => $_clearField(244019094);

  @$pb.TagNumber(328257828)
  $core.String get accessassociationsource => $_getSZ(2);
  @$pb.TagNumber(328257828)
  set accessassociationsource($core.String value) => $_setString(2, value);
  @$pb.TagNumber(328257828)
  $core.bool hasAccessassociationsource() => $_has(2);
  @$pb.TagNumber(328257828)
  void clearAccessassociationsource() => $_clearField(328257828);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(3);
}

class CreateDomainNameRequest extends $pb.GeneratedMessage {
  factory CreateDomainNameRequest({
    $core.String? certificatechain,
    MutualTlsAuthenticationInput? mutualtlsauthentication,
    $core.String? ownershipverificationcertificatearn,
    $core.String? regionalcertificatearn,
    $core.String? certificatename,
    $core.String? certificatebody,
    $core.String? policy,
    $core.String? certificateprivatekey,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    EndpointAccessMode? endpointaccessmode,
    $core.String? domainname,
    $core.String? certificatearn,
    $core.String? regionalcertificatename,
    EndpointConfiguration? endpointconfiguration,
    SecurityPolicy? securitypolicy,
    RoutingMode? routingmode,
  }) {
    final result = create();
    if (certificatechain != null) result.certificatechain = certificatechain;
    if (mutualtlsauthentication != null)
      result.mutualtlsauthentication = mutualtlsauthentication;
    if (ownershipverificationcertificatearn != null)
      result.ownershipverificationcertificatearn =
          ownershipverificationcertificatearn;
    if (regionalcertificatearn != null)
      result.regionalcertificatearn = regionalcertificatearn;
    if (certificatename != null) result.certificatename = certificatename;
    if (certificatebody != null) result.certificatebody = certificatebody;
    if (policy != null) result.policy = policy;
    if (certificateprivatekey != null)
      result.certificateprivatekey = certificateprivatekey;
    if (tags != null) result.tags.addEntries(tags);
    if (endpointaccessmode != null)
      result.endpointaccessmode = endpointaccessmode;
    if (domainname != null) result.domainname = domainname;
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (regionalcertificatename != null)
      result.regionalcertificatename = regionalcertificatename;
    if (endpointconfiguration != null)
      result.endpointconfiguration = endpointconfiguration;
    if (securitypolicy != null) result.securitypolicy = securitypolicy;
    if (routingmode != null) result.routingmode = routingmode;
    return result;
  }

  CreateDomainNameRequest._();

  factory CreateDomainNameRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateDomainNameRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateDomainNameRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(50547450, _omitFieldNames ? '' : 'certificatechain')
    ..aOM<MutualTlsAuthenticationInput>(
        99462043, _omitFieldNames ? '' : 'mutualtlsauthentication',
        subBuilder: MutualTlsAuthenticationInput.create)
    ..aOS(
        100550548, _omitFieldNames ? '' : 'ownershipverificationcertificatearn')
    ..aOS(137066579, _omitFieldNames ? '' : 'regionalcertificatearn')
    ..aOS(140276948, _omitFieldNames ? '' : 'certificatename')
    ..aOS(183933821, _omitFieldNames ? '' : 'certificatebody')
    ..aOS(247528064, _omitFieldNames ? '' : 'policy')
    ..aOS(277544067, _omitFieldNames ? '' : 'certificateprivatekey')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateDomainNameRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aE<EndpointAccessMode>(
        356705630, _omitFieldNames ? '' : 'endpointaccessmode',
        enumValues: EndpointAccessMode.values)
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..aOS(425831704, _omitFieldNames ? '' : 'certificatearn')
    ..aOS(431716941, _omitFieldNames ? '' : 'regionalcertificatename')
    ..aOM<EndpointConfiguration>(
        487543735, _omitFieldNames ? '' : 'endpointconfiguration',
        subBuilder: EndpointConfiguration.create)
    ..aE<SecurityPolicy>(491792990, _omitFieldNames ? '' : 'securitypolicy',
        enumValues: SecurityPolicy.values)
    ..aE<RoutingMode>(506342119, _omitFieldNames ? '' : 'routingmode',
        enumValues: RoutingMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDomainNameRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateDomainNameRequest copyWith(
          void Function(CreateDomainNameRequest) updates) =>
      super.copyWith((message) => updates(message as CreateDomainNameRequest))
          as CreateDomainNameRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateDomainNameRequest create() => CreateDomainNameRequest._();
  @$core.override
  CreateDomainNameRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateDomainNameRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateDomainNameRequest>(create);
  static CreateDomainNameRequest? _defaultInstance;

  @$pb.TagNumber(50547450)
  $core.String get certificatechain => $_getSZ(0);
  @$pb.TagNumber(50547450)
  set certificatechain($core.String value) => $_setString(0, value);
  @$pb.TagNumber(50547450)
  $core.bool hasCertificatechain() => $_has(0);
  @$pb.TagNumber(50547450)
  void clearCertificatechain() => $_clearField(50547450);

  @$pb.TagNumber(99462043)
  MutualTlsAuthenticationInput get mutualtlsauthentication => $_getN(1);
  @$pb.TagNumber(99462043)
  set mutualtlsauthentication(MutualTlsAuthenticationInput value) =>
      $_setField(99462043, value);
  @$pb.TagNumber(99462043)
  $core.bool hasMutualtlsauthentication() => $_has(1);
  @$pb.TagNumber(99462043)
  void clearMutualtlsauthentication() => $_clearField(99462043);
  @$pb.TagNumber(99462043)
  MutualTlsAuthenticationInput ensureMutualtlsauthentication() => $_ensure(1);

  @$pb.TagNumber(100550548)
  $core.String get ownershipverificationcertificatearn => $_getSZ(2);
  @$pb.TagNumber(100550548)
  set ownershipverificationcertificatearn($core.String value) =>
      $_setString(2, value);
  @$pb.TagNumber(100550548)
  $core.bool hasOwnershipverificationcertificatearn() => $_has(2);
  @$pb.TagNumber(100550548)
  void clearOwnershipverificationcertificatearn() => $_clearField(100550548);

  @$pb.TagNumber(137066579)
  $core.String get regionalcertificatearn => $_getSZ(3);
  @$pb.TagNumber(137066579)
  set regionalcertificatearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(137066579)
  $core.bool hasRegionalcertificatearn() => $_has(3);
  @$pb.TagNumber(137066579)
  void clearRegionalcertificatearn() => $_clearField(137066579);

  @$pb.TagNumber(140276948)
  $core.String get certificatename => $_getSZ(4);
  @$pb.TagNumber(140276948)
  set certificatename($core.String value) => $_setString(4, value);
  @$pb.TagNumber(140276948)
  $core.bool hasCertificatename() => $_has(4);
  @$pb.TagNumber(140276948)
  void clearCertificatename() => $_clearField(140276948);

  @$pb.TagNumber(183933821)
  $core.String get certificatebody => $_getSZ(5);
  @$pb.TagNumber(183933821)
  set certificatebody($core.String value) => $_setString(5, value);
  @$pb.TagNumber(183933821)
  $core.bool hasCertificatebody() => $_has(5);
  @$pb.TagNumber(183933821)
  void clearCertificatebody() => $_clearField(183933821);

  @$pb.TagNumber(247528064)
  $core.String get policy => $_getSZ(6);
  @$pb.TagNumber(247528064)
  set policy($core.String value) => $_setString(6, value);
  @$pb.TagNumber(247528064)
  $core.bool hasPolicy() => $_has(6);
  @$pb.TagNumber(247528064)
  void clearPolicy() => $_clearField(247528064);

  @$pb.TagNumber(277544067)
  $core.String get certificateprivatekey => $_getSZ(7);
  @$pb.TagNumber(277544067)
  set certificateprivatekey($core.String value) => $_setString(7, value);
  @$pb.TagNumber(277544067)
  $core.bool hasCertificateprivatekey() => $_has(7);
  @$pb.TagNumber(277544067)
  void clearCertificateprivatekey() => $_clearField(277544067);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(8);

  @$pb.TagNumber(356705630)
  EndpointAccessMode get endpointaccessmode => $_getN(9);
  @$pb.TagNumber(356705630)
  set endpointaccessmode(EndpointAccessMode value) =>
      $_setField(356705630, value);
  @$pb.TagNumber(356705630)
  $core.bool hasEndpointaccessmode() => $_has(9);
  @$pb.TagNumber(356705630)
  void clearEndpointaccessmode() => $_clearField(356705630);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(10);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(10, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(10);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);

  @$pb.TagNumber(425831704)
  $core.String get certificatearn => $_getSZ(11);
  @$pb.TagNumber(425831704)
  set certificatearn($core.String value) => $_setString(11, value);
  @$pb.TagNumber(425831704)
  $core.bool hasCertificatearn() => $_has(11);
  @$pb.TagNumber(425831704)
  void clearCertificatearn() => $_clearField(425831704);

  @$pb.TagNumber(431716941)
  $core.String get regionalcertificatename => $_getSZ(12);
  @$pb.TagNumber(431716941)
  set regionalcertificatename($core.String value) => $_setString(12, value);
  @$pb.TagNumber(431716941)
  $core.bool hasRegionalcertificatename() => $_has(12);
  @$pb.TagNumber(431716941)
  void clearRegionalcertificatename() => $_clearField(431716941);

  @$pb.TagNumber(487543735)
  EndpointConfiguration get endpointconfiguration => $_getN(13);
  @$pb.TagNumber(487543735)
  set endpointconfiguration(EndpointConfiguration value) =>
      $_setField(487543735, value);
  @$pb.TagNumber(487543735)
  $core.bool hasEndpointconfiguration() => $_has(13);
  @$pb.TagNumber(487543735)
  void clearEndpointconfiguration() => $_clearField(487543735);
  @$pb.TagNumber(487543735)
  EndpointConfiguration ensureEndpointconfiguration() => $_ensure(13);

  @$pb.TagNumber(491792990)
  SecurityPolicy get securitypolicy => $_getN(14);
  @$pb.TagNumber(491792990)
  set securitypolicy(SecurityPolicy value) => $_setField(491792990, value);
  @$pb.TagNumber(491792990)
  $core.bool hasSecuritypolicy() => $_has(14);
  @$pb.TagNumber(491792990)
  void clearSecuritypolicy() => $_clearField(491792990);

  @$pb.TagNumber(506342119)
  RoutingMode get routingmode => $_getN(15);
  @$pb.TagNumber(506342119)
  set routingmode(RoutingMode value) => $_setField(506342119, value);
  @$pb.TagNumber(506342119)
  $core.bool hasRoutingmode() => $_has(15);
  @$pb.TagNumber(506342119)
  void clearRoutingmode() => $_clearField(506342119);
}

class CreateModelRequest extends $pb.GeneratedMessage {
  factory CreateModelRequest({
    $core.String? name,
    $core.String? contenttype,
    $core.String? schema,
    $core.String? description,
    $core.String? restapiid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (contenttype != null) result.contenttype = contenttype;
    if (schema != null) result.schema = schema;
    if (description != null) result.description = description;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  CreateModelRequest._();

  factory CreateModelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateModelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateModelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(281764659, _omitFieldNames ? '' : 'contenttype')
    ..aOS(310182711, _omitFieldNames ? '' : 'schema')
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateModelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateModelRequest copyWith(void Function(CreateModelRequest) updates) =>
      super.copyWith((message) => updates(message as CreateModelRequest))
          as CreateModelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateModelRequest create() => CreateModelRequest._();
  @$core.override
  CreateModelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateModelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateModelRequest>(create);
  static CreateModelRequest? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(281764659)
  $core.String get contenttype => $_getSZ(1);
  @$pb.TagNumber(281764659)
  set contenttype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(281764659)
  $core.bool hasContenttype() => $_has(1);
  @$pb.TagNumber(281764659)
  void clearContenttype() => $_clearField(281764659);

  @$pb.TagNumber(310182711)
  $core.String get schema => $_getSZ(2);
  @$pb.TagNumber(310182711)
  set schema($core.String value) => $_setString(2, value);
  @$pb.TagNumber(310182711)
  $core.bool hasSchema() => $_has(2);
  @$pb.TagNumber(310182711)
  void clearSchema() => $_clearField(310182711);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(4);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(4);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class CreateRequestValidatorRequest extends $pb.GeneratedMessage {
  factory CreateRequestValidatorRequest({
    $core.String? name,
    $core.String? restapiid,
    $core.bool? validaterequestbody,
    $core.bool? validaterequestparameters,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (restapiid != null) result.restapiid = restapiid;
    if (validaterequestbody != null)
      result.validaterequestbody = validaterequestbody;
    if (validaterequestparameters != null)
      result.validaterequestparameters = validaterequestparameters;
    return result;
  }

  CreateRequestValidatorRequest._();

  factory CreateRequestValidatorRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRequestValidatorRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRequestValidatorRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOB(397505841, _omitFieldNames ? '' : 'validaterequestbody')
    ..aOB(464035801, _omitFieldNames ? '' : 'validaterequestparameters')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRequestValidatorRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRequestValidatorRequest copyWith(
          void Function(CreateRequestValidatorRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateRequestValidatorRequest))
          as CreateRequestValidatorRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRequestValidatorRequest create() =>
      CreateRequestValidatorRequest._();
  @$core.override
  CreateRequestValidatorRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRequestValidatorRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRequestValidatorRequest>(create);
  static CreateRequestValidatorRequest? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(397505841)
  $core.bool get validaterequestbody => $_getBF(2);
  @$pb.TagNumber(397505841)
  set validaterequestbody($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(397505841)
  $core.bool hasValidaterequestbody() => $_has(2);
  @$pb.TagNumber(397505841)
  void clearValidaterequestbody() => $_clearField(397505841);

  @$pb.TagNumber(464035801)
  $core.bool get validaterequestparameters => $_getBF(3);
  @$pb.TagNumber(464035801)
  set validaterequestparameters($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(464035801)
  $core.bool hasValidaterequestparameters() => $_has(3);
  @$pb.TagNumber(464035801)
  void clearValidaterequestparameters() => $_clearField(464035801);
}

class CreateResourceRequest extends $pb.GeneratedMessage {
  factory CreateResourceRequest({
    $core.String? parentid,
    $core.String? restapiid,
    $core.String? pathpart,
  }) {
    final result = create();
    if (parentid != null) result.parentid = parentid;
    if (restapiid != null) result.restapiid = restapiid;
    if (pathpart != null) result.pathpart = pathpart;
    return result;
  }

  CreateResourceRequest._();

  factory CreateResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateResourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(106050857, _omitFieldNames ? '' : 'parentid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(487915984, _omitFieldNames ? '' : 'pathpart')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateResourceRequest copyWith(
          void Function(CreateResourceRequest) updates) =>
      super.copyWith((message) => updates(message as CreateResourceRequest))
          as CreateResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateResourceRequest create() => CreateResourceRequest._();
  @$core.override
  CreateResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateResourceRequest>(create);
  static CreateResourceRequest? _defaultInstance;

  @$pb.TagNumber(106050857)
  $core.String get parentid => $_getSZ(0);
  @$pb.TagNumber(106050857)
  set parentid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106050857)
  $core.bool hasParentid() => $_has(0);
  @$pb.TagNumber(106050857)
  void clearParentid() => $_clearField(106050857);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(487915984)
  $core.String get pathpart => $_getSZ(2);
  @$pb.TagNumber(487915984)
  set pathpart($core.String value) => $_setString(2, value);
  @$pb.TagNumber(487915984)
  $core.bool hasPathpart() => $_has(2);
  @$pb.TagNumber(487915984)
  void clearPathpart() => $_clearField(487915984);
}

class CreateRestApiRequest extends $pb.GeneratedMessage {
  factory CreateRestApiRequest({
    $core.String? version,
    ApiKeySourceType? apikeysource,
    $core.bool? disableexecuteapiendpoint,
    $core.String? name,
    $core.String? policy,
    $core.int? minimumcompressionsize,
    $core.String? clonefrom,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    EndpointAccessMode? endpointaccessmode,
    $core.Iterable<$core.String>? binarymediatypes,
    EndpointConfiguration? endpointconfiguration,
    SecurityPolicy? securitypolicy,
  }) {
    final result = create();
    if (version != null) result.version = version;
    if (apikeysource != null) result.apikeysource = apikeysource;
    if (disableexecuteapiendpoint != null)
      result.disableexecuteapiendpoint = disableexecuteapiendpoint;
    if (name != null) result.name = name;
    if (policy != null) result.policy = policy;
    if (minimumcompressionsize != null)
      result.minimumcompressionsize = minimumcompressionsize;
    if (clonefrom != null) result.clonefrom = clonefrom;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (endpointaccessmode != null)
      result.endpointaccessmode = endpointaccessmode;
    if (binarymediatypes != null)
      result.binarymediatypes.addAll(binarymediatypes);
    if (endpointconfiguration != null)
      result.endpointconfiguration = endpointconfiguration;
    if (securitypolicy != null) result.securitypolicy = securitypolicy;
    return result;
  }

  CreateRestApiRequest._();

  factory CreateRestApiRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRestApiRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRestApiRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(108113560, _omitFieldNames ? '' : 'version')
    ..aE<ApiKeySourceType>(108531220, _omitFieldNames ? '' : 'apikeysource',
        enumValues: ApiKeySourceType.values)
    ..aOB(148140696, _omitFieldNames ? '' : 'disableexecuteapiendpoint')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(247528064, _omitFieldNames ? '' : 'policy')
    ..aI(254902719, _omitFieldNames ? '' : 'minimumcompressionsize')
    ..aOS(263376551, _omitFieldNames ? '' : 'clonefrom')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateRestApiRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aE<EndpointAccessMode>(
        356705630, _omitFieldNames ? '' : 'endpointaccessmode',
        enumValues: EndpointAccessMode.values)
    ..pPS(406416146, _omitFieldNames ? '' : 'binarymediatypes')
    ..aOM<EndpointConfiguration>(
        487543735, _omitFieldNames ? '' : 'endpointconfiguration',
        subBuilder: EndpointConfiguration.create)
    ..aE<SecurityPolicy>(491792990, _omitFieldNames ? '' : 'securitypolicy',
        enumValues: SecurityPolicy.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRestApiRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRestApiRequest copyWith(void Function(CreateRestApiRequest) updates) =>
      super.copyWith((message) => updates(message as CreateRestApiRequest))
          as CreateRestApiRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRestApiRequest create() => CreateRestApiRequest._();
  @$core.override
  CreateRestApiRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRestApiRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRestApiRequest>(create);
  static CreateRestApiRequest? _defaultInstance;

  @$pb.TagNumber(108113560)
  $core.String get version => $_getSZ(0);
  @$pb.TagNumber(108113560)
  set version($core.String value) => $_setString(0, value);
  @$pb.TagNumber(108113560)
  $core.bool hasVersion() => $_has(0);
  @$pb.TagNumber(108113560)
  void clearVersion() => $_clearField(108113560);

  @$pb.TagNumber(108531220)
  ApiKeySourceType get apikeysource => $_getN(1);
  @$pb.TagNumber(108531220)
  set apikeysource(ApiKeySourceType value) => $_setField(108531220, value);
  @$pb.TagNumber(108531220)
  $core.bool hasApikeysource() => $_has(1);
  @$pb.TagNumber(108531220)
  void clearApikeysource() => $_clearField(108531220);

  @$pb.TagNumber(148140696)
  $core.bool get disableexecuteapiendpoint => $_getBF(2);
  @$pb.TagNumber(148140696)
  set disableexecuteapiendpoint($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(148140696)
  $core.bool hasDisableexecuteapiendpoint() => $_has(2);
  @$pb.TagNumber(148140696)
  void clearDisableexecuteapiendpoint() => $_clearField(148140696);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(247528064)
  $core.String get policy => $_getSZ(4);
  @$pb.TagNumber(247528064)
  set policy($core.String value) => $_setString(4, value);
  @$pb.TagNumber(247528064)
  $core.bool hasPolicy() => $_has(4);
  @$pb.TagNumber(247528064)
  void clearPolicy() => $_clearField(247528064);

  @$pb.TagNumber(254902719)
  $core.int get minimumcompressionsize => $_getIZ(5);
  @$pb.TagNumber(254902719)
  set minimumcompressionsize($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(254902719)
  $core.bool hasMinimumcompressionsize() => $_has(5);
  @$pb.TagNumber(254902719)
  void clearMinimumcompressionsize() => $_clearField(254902719);

  @$pb.TagNumber(263376551)
  $core.String get clonefrom => $_getSZ(6);
  @$pb.TagNumber(263376551)
  set clonefrom($core.String value) => $_setString(6, value);
  @$pb.TagNumber(263376551)
  $core.bool hasClonefrom() => $_has(6);
  @$pb.TagNumber(263376551)
  void clearClonefrom() => $_clearField(263376551);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(7);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(8);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(8, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(8);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(356705630)
  EndpointAccessMode get endpointaccessmode => $_getN(9);
  @$pb.TagNumber(356705630)
  set endpointaccessmode(EndpointAccessMode value) =>
      $_setField(356705630, value);
  @$pb.TagNumber(356705630)
  $core.bool hasEndpointaccessmode() => $_has(9);
  @$pb.TagNumber(356705630)
  void clearEndpointaccessmode() => $_clearField(356705630);

  @$pb.TagNumber(406416146)
  $pb.PbList<$core.String> get binarymediatypes => $_getList(10);

  @$pb.TagNumber(487543735)
  EndpointConfiguration get endpointconfiguration => $_getN(11);
  @$pb.TagNumber(487543735)
  set endpointconfiguration(EndpointConfiguration value) =>
      $_setField(487543735, value);
  @$pb.TagNumber(487543735)
  $core.bool hasEndpointconfiguration() => $_has(11);
  @$pb.TagNumber(487543735)
  void clearEndpointconfiguration() => $_clearField(487543735);
  @$pb.TagNumber(487543735)
  EndpointConfiguration ensureEndpointconfiguration() => $_ensure(11);

  @$pb.TagNumber(491792990)
  SecurityPolicy get securitypolicy => $_getN(12);
  @$pb.TagNumber(491792990)
  set securitypolicy(SecurityPolicy value) => $_setField(491792990, value);
  @$pb.TagNumber(491792990)
  $core.bool hasSecuritypolicy() => $_has(12);
  @$pb.TagNumber(491792990)
  void clearSecuritypolicy() => $_clearField(491792990);
}

class CreateStageRequest extends $pb.GeneratedMessage {
  factory CreateStageRequest({
    $core.String? stagename,
    $core.bool? cacheclusterenabled,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? variables,
    $core.String? documentationversion,
    CacheClusterSize? cacheclustersize,
    CanarySettings? canarysettings,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    $core.String? restapiid,
    $core.bool? tracingenabled,
    $core.String? deploymentid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (cacheclusterenabled != null)
      result.cacheclusterenabled = cacheclusterenabled;
    if (variables != null) result.variables.addEntries(variables);
    if (documentationversion != null)
      result.documentationversion = documentationversion;
    if (cacheclustersize != null) result.cacheclustersize = cacheclustersize;
    if (canarysettings != null) result.canarysettings = canarysettings;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (restapiid != null) result.restapiid = restapiid;
    if (tracingenabled != null) result.tracingenabled = tracingenabled;
    if (deploymentid != null) result.deploymentid = deploymentid;
    return result;
  }

  CreateStageRequest._();

  factory CreateStageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateStageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateStageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..aOB(63967991, _omitFieldNames ? '' : 'cacheclusterenabled')
    ..m<$core.String, $core.String>(
        162226883, _omitFieldNames ? '' : 'variables',
        entryClassName: 'CreateStageRequest.VariablesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(167009804, _omitFieldNames ? '' : 'documentationversion')
    ..aE<CacheClusterSize>(232189861, _omitFieldNames ? '' : 'cacheclustersize',
        enumValues: CacheClusterSize.values)
    ..aOM<CanarySettings>(285544261, _omitFieldNames ? '' : 'canarysettings',
        subBuilder: CanarySettings.create)
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateStageRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOB(390995731, _omitFieldNames ? '' : 'tracingenabled')
    ..aOS(439369188, _omitFieldNames ? '' : 'deploymentid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateStageRequest copyWith(void Function(CreateStageRequest) updates) =>
      super.copyWith((message) => updates(message as CreateStageRequest))
          as CreateStageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateStageRequest create() => CreateStageRequest._();
  @$core.override
  CreateStageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateStageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateStageRequest>(create);
  static CreateStageRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(63967991)
  $core.bool get cacheclusterenabled => $_getBF(1);
  @$pb.TagNumber(63967991)
  set cacheclusterenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(63967991)
  $core.bool hasCacheclusterenabled() => $_has(1);
  @$pb.TagNumber(63967991)
  void clearCacheclusterenabled() => $_clearField(63967991);

  @$pb.TagNumber(162226883)
  $pb.PbMap<$core.String, $core.String> get variables => $_getMap(2);

  @$pb.TagNumber(167009804)
  $core.String get documentationversion => $_getSZ(3);
  @$pb.TagNumber(167009804)
  set documentationversion($core.String value) => $_setString(3, value);
  @$pb.TagNumber(167009804)
  $core.bool hasDocumentationversion() => $_has(3);
  @$pb.TagNumber(167009804)
  void clearDocumentationversion() => $_clearField(167009804);

  @$pb.TagNumber(232189861)
  CacheClusterSize get cacheclustersize => $_getN(4);
  @$pb.TagNumber(232189861)
  set cacheclustersize(CacheClusterSize value) => $_setField(232189861, value);
  @$pb.TagNumber(232189861)
  $core.bool hasCacheclustersize() => $_has(4);
  @$pb.TagNumber(232189861)
  void clearCacheclustersize() => $_clearField(232189861);

  @$pb.TagNumber(285544261)
  CanarySettings get canarysettings => $_getN(5);
  @$pb.TagNumber(285544261)
  set canarysettings(CanarySettings value) => $_setField(285544261, value);
  @$pb.TagNumber(285544261)
  $core.bool hasCanarysettings() => $_has(5);
  @$pb.TagNumber(285544261)
  void clearCanarysettings() => $_clearField(285544261);
  @$pb.TagNumber(285544261)
  CanarySettings ensureCanarysettings() => $_ensure(5);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(6);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(7);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(7, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(7);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(8);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(8);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(390995731)
  $core.bool get tracingenabled => $_getBF(9);
  @$pb.TagNumber(390995731)
  set tracingenabled($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(390995731)
  $core.bool hasTracingenabled() => $_has(9);
  @$pb.TagNumber(390995731)
  void clearTracingenabled() => $_clearField(390995731);

  @$pb.TagNumber(439369188)
  $core.String get deploymentid => $_getSZ(10);
  @$pb.TagNumber(439369188)
  set deploymentid($core.String value) => $_setString(10, value);
  @$pb.TagNumber(439369188)
  $core.bool hasDeploymentid() => $_has(10);
  @$pb.TagNumber(439369188)
  void clearDeploymentid() => $_clearField(439369188);
}

class CreateUsagePlanKeyRequest extends $pb.GeneratedMessage {
  factory CreateUsagePlanKeyRequest({
    $core.String? keytype,
    $core.String? keyid,
    $core.String? usageplanid,
  }) {
    final result = create();
    if (keytype != null) result.keytype = keytype;
    if (keyid != null) result.keyid = keyid;
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  CreateUsagePlanKeyRequest._();

  factory CreateUsagePlanKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateUsagePlanKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateUsagePlanKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(107017349, _omitFieldNames ? '' : 'keytype')
    ..aOS(479913282, _omitFieldNames ? '' : 'keyid')
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateUsagePlanKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateUsagePlanKeyRequest copyWith(
          void Function(CreateUsagePlanKeyRequest) updates) =>
      super.copyWith((message) => updates(message as CreateUsagePlanKeyRequest))
          as CreateUsagePlanKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateUsagePlanKeyRequest create() => CreateUsagePlanKeyRequest._();
  @$core.override
  CreateUsagePlanKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateUsagePlanKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateUsagePlanKeyRequest>(create);
  static CreateUsagePlanKeyRequest? _defaultInstance;

  @$pb.TagNumber(107017349)
  $core.String get keytype => $_getSZ(0);
  @$pb.TagNumber(107017349)
  set keytype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(107017349)
  $core.bool hasKeytype() => $_has(0);
  @$pb.TagNumber(107017349)
  void clearKeytype() => $_clearField(107017349);

  @$pb.TagNumber(479913282)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(479913282)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(479913282)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(479913282)
  void clearKeyid() => $_clearField(479913282);

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(2);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(2);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class CreateUsagePlanRequest extends $pb.GeneratedMessage {
  factory CreateUsagePlanRequest({
    $core.Iterable<ApiStage>? apistages,
    $core.String? name,
    QuotaSettings? quota,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    ThrottleSettings? throttle,
  }) {
    final result = create();
    if (apistages != null) result.apistages.addAll(apistages);
    if (name != null) result.name = name;
    if (quota != null) result.quota = quota;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (throttle != null) result.throttle = throttle;
    return result;
  }

  CreateUsagePlanRequest._();

  factory CreateUsagePlanRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateUsagePlanRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateUsagePlanRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<ApiStage>(64558449, _omitFieldNames ? '' : 'apistages',
        subBuilder: ApiStage.create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOM<QuotaSettings>(243824012, _omitFieldNames ? '' : 'quota',
        subBuilder: QuotaSettings.create)
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateUsagePlanRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOM<ThrottleSettings>(395260638, _omitFieldNames ? '' : 'throttle',
        subBuilder: ThrottleSettings.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateUsagePlanRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateUsagePlanRequest copyWith(
          void Function(CreateUsagePlanRequest) updates) =>
      super.copyWith((message) => updates(message as CreateUsagePlanRequest))
          as CreateUsagePlanRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateUsagePlanRequest create() => CreateUsagePlanRequest._();
  @$core.override
  CreateUsagePlanRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateUsagePlanRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateUsagePlanRequest>(create);
  static CreateUsagePlanRequest? _defaultInstance;

  @$pb.TagNumber(64558449)
  $pb.PbList<ApiStage> get apistages => $_getList(0);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(243824012)
  QuotaSettings get quota => $_getN(2);
  @$pb.TagNumber(243824012)
  set quota(QuotaSettings value) => $_setField(243824012, value);
  @$pb.TagNumber(243824012)
  $core.bool hasQuota() => $_has(2);
  @$pb.TagNumber(243824012)
  void clearQuota() => $_clearField(243824012);
  @$pb.TagNumber(243824012)
  QuotaSettings ensureQuota() => $_ensure(2);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(3);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(395260638)
  ThrottleSettings get throttle => $_getN(5);
  @$pb.TagNumber(395260638)
  set throttle(ThrottleSettings value) => $_setField(395260638, value);
  @$pb.TagNumber(395260638)
  $core.bool hasThrottle() => $_has(5);
  @$pb.TagNumber(395260638)
  void clearThrottle() => $_clearField(395260638);
  @$pb.TagNumber(395260638)
  ThrottleSettings ensureThrottle() => $_ensure(5);
}

class CreateVpcLinkRequest extends $pb.GeneratedMessage {
  factory CreateVpcLinkRequest({
    $core.Iterable<$core.String>? targetarns,
    $core.String? name,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
  }) {
    final result = create();
    if (targetarns != null) result.targetarns.addAll(targetarns);
    if (name != null) result.name = name;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    return result;
  }

  CreateVpcLinkRequest._();

  factory CreateVpcLinkRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateVpcLinkRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateVpcLinkRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(46319317, _omitFieldNames ? '' : 'targetarns')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'CreateVpcLinkRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateVpcLinkRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateVpcLinkRequest copyWith(void Function(CreateVpcLinkRequest) updates) =>
      super.copyWith((message) => updates(message as CreateVpcLinkRequest))
          as CreateVpcLinkRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateVpcLinkRequest create() => CreateVpcLinkRequest._();
  @$core.override
  CreateVpcLinkRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateVpcLinkRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateVpcLinkRequest>(create);
  static CreateVpcLinkRequest? _defaultInstance;

  @$pb.TagNumber(46319317)
  $pb.PbList<$core.String> get targetarns => $_getList(0);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(2);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);
}

class DeleteApiKeyRequest extends $pb.GeneratedMessage {
  factory DeleteApiKeyRequest({
    $core.String? apikey,
  }) {
    final result = create();
    if (apikey != null) result.apikey = apikey;
    return result;
  }

  DeleteApiKeyRequest._();

  factory DeleteApiKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteApiKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteApiKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(490490655, _omitFieldNames ? '' : 'apikey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteApiKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteApiKeyRequest copyWith(void Function(DeleteApiKeyRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteApiKeyRequest))
          as DeleteApiKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteApiKeyRequest create() => DeleteApiKeyRequest._();
  @$core.override
  DeleteApiKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteApiKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteApiKeyRequest>(create);
  static DeleteApiKeyRequest? _defaultInstance;

  @$pb.TagNumber(490490655)
  $core.String get apikey => $_getSZ(0);
  @$pb.TagNumber(490490655)
  set apikey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(490490655)
  $core.bool hasApikey() => $_has(0);
  @$pb.TagNumber(490490655)
  void clearApikey() => $_clearField(490490655);
}

class DeleteAuthorizerRequest extends $pb.GeneratedMessage {
  factory DeleteAuthorizerRequest({
    $core.String? authorizerid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (authorizerid != null) result.authorizerid = authorizerid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteAuthorizerRequest._();

  factory DeleteAuthorizerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteAuthorizerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAuthorizerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(111773148, _omitFieldNames ? '' : 'authorizerid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAuthorizerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAuthorizerRequest copyWith(
          void Function(DeleteAuthorizerRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteAuthorizerRequest))
          as DeleteAuthorizerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAuthorizerRequest create() => DeleteAuthorizerRequest._();
  @$core.override
  DeleteAuthorizerRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteAuthorizerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteAuthorizerRequest>(create);
  static DeleteAuthorizerRequest? _defaultInstance;

  @$pb.TagNumber(111773148)
  $core.String get authorizerid => $_getSZ(0);
  @$pb.TagNumber(111773148)
  set authorizerid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(111773148)
  $core.bool hasAuthorizerid() => $_has(0);
  @$pb.TagNumber(111773148)
  void clearAuthorizerid() => $_clearField(111773148);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteBasePathMappingRequest extends $pb.GeneratedMessage {
  factory DeleteBasePathMappingRequest({
    $core.String? basepath,
    $core.String? domainnameid,
    $core.String? domainname,
  }) {
    final result = create();
    if (basepath != null) result.basepath = basepath;
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  DeleteBasePathMappingRequest._();

  factory DeleteBasePathMappingRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteBasePathMappingRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteBasePathMappingRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(267528880, _omitFieldNames ? '' : 'basepath')
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteBasePathMappingRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteBasePathMappingRequest copyWith(
          void Function(DeleteBasePathMappingRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteBasePathMappingRequest))
          as DeleteBasePathMappingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteBasePathMappingRequest create() =>
      DeleteBasePathMappingRequest._();
  @$core.override
  DeleteBasePathMappingRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteBasePathMappingRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteBasePathMappingRequest>(create);
  static DeleteBasePathMappingRequest? _defaultInstance;

  @$pb.TagNumber(267528880)
  $core.String get basepath => $_getSZ(0);
  @$pb.TagNumber(267528880)
  set basepath($core.String value) => $_setString(0, value);
  @$pb.TagNumber(267528880)
  $core.bool hasBasepath() => $_has(0);
  @$pb.TagNumber(267528880)
  void clearBasepath() => $_clearField(267528880);

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(1);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(1);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(2);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(2);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);
}

class DeleteClientCertificateRequest extends $pb.GeneratedMessage {
  factory DeleteClientCertificateRequest({
    $core.String? clientcertificateid,
  }) {
    final result = create();
    if (clientcertificateid != null)
      result.clientcertificateid = clientcertificateid;
    return result;
  }

  DeleteClientCertificateRequest._();

  factory DeleteClientCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteClientCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteClientCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(276222909, _omitFieldNames ? '' : 'clientcertificateid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteClientCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteClientCertificateRequest copyWith(
          void Function(DeleteClientCertificateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteClientCertificateRequest))
          as DeleteClientCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteClientCertificateRequest create() =>
      DeleteClientCertificateRequest._();
  @$core.override
  DeleteClientCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteClientCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteClientCertificateRequest>(create);
  static DeleteClientCertificateRequest? _defaultInstance;

  @$pb.TagNumber(276222909)
  $core.String get clientcertificateid => $_getSZ(0);
  @$pb.TagNumber(276222909)
  set clientcertificateid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(276222909)
  $core.bool hasClientcertificateid() => $_has(0);
  @$pb.TagNumber(276222909)
  void clearClientcertificateid() => $_clearField(276222909);
}

class DeleteDeploymentRequest extends $pb.GeneratedMessage {
  factory DeleteDeploymentRequest({
    $core.String? restapiid,
    $core.String? deploymentid,
  }) {
    final result = create();
    if (restapiid != null) result.restapiid = restapiid;
    if (deploymentid != null) result.deploymentid = deploymentid;
    return result;
  }

  DeleteDeploymentRequest._();

  factory DeleteDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(439369188, _omitFieldNames ? '' : 'deploymentid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDeploymentRequest copyWith(
          void Function(DeleteDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteDeploymentRequest))
          as DeleteDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDeploymentRequest create() => DeleteDeploymentRequest._();
  @$core.override
  DeleteDeploymentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDeploymentRequest>(create);
  static DeleteDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(0);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(0);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(439369188)
  $core.String get deploymentid => $_getSZ(1);
  @$pb.TagNumber(439369188)
  set deploymentid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(439369188)
  $core.bool hasDeploymentid() => $_has(1);
  @$pb.TagNumber(439369188)
  void clearDeploymentid() => $_clearField(439369188);
}

class DeleteDocumentationPartRequest extends $pb.GeneratedMessage {
  factory DeleteDocumentationPartRequest({
    $core.String? documentationpartid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (documentationpartid != null)
      result.documentationpartid = documentationpartid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteDocumentationPartRequest._();

  factory DeleteDocumentationPartRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDocumentationPartRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDocumentationPartRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(286552774, _omitFieldNames ? '' : 'documentationpartid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDocumentationPartRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDocumentationPartRequest copyWith(
          void Function(DeleteDocumentationPartRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteDocumentationPartRequest))
          as DeleteDocumentationPartRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDocumentationPartRequest create() =>
      DeleteDocumentationPartRequest._();
  @$core.override
  DeleteDocumentationPartRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDocumentationPartRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDocumentationPartRequest>(create);
  static DeleteDocumentationPartRequest? _defaultInstance;

  @$pb.TagNumber(286552774)
  $core.String get documentationpartid => $_getSZ(0);
  @$pb.TagNumber(286552774)
  set documentationpartid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(286552774)
  $core.bool hasDocumentationpartid() => $_has(0);
  @$pb.TagNumber(286552774)
  void clearDocumentationpartid() => $_clearField(286552774);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteDocumentationVersionRequest extends $pb.GeneratedMessage {
  factory DeleteDocumentationVersionRequest({
    $core.String? documentationversion,
    $core.String? restapiid,
  }) {
    final result = create();
    if (documentationversion != null)
      result.documentationversion = documentationversion;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteDocumentationVersionRequest._();

  factory DeleteDocumentationVersionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDocumentationVersionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDocumentationVersionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(167009804, _omitFieldNames ? '' : 'documentationversion')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDocumentationVersionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDocumentationVersionRequest copyWith(
          void Function(DeleteDocumentationVersionRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteDocumentationVersionRequest))
          as DeleteDocumentationVersionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDocumentationVersionRequest create() =>
      DeleteDocumentationVersionRequest._();
  @$core.override
  DeleteDocumentationVersionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDocumentationVersionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDocumentationVersionRequest>(
          create);
  static DeleteDocumentationVersionRequest? _defaultInstance;

  @$pb.TagNumber(167009804)
  $core.String get documentationversion => $_getSZ(0);
  @$pb.TagNumber(167009804)
  set documentationversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(167009804)
  $core.bool hasDocumentationversion() => $_has(0);
  @$pb.TagNumber(167009804)
  void clearDocumentationversion() => $_clearField(167009804);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteDomainNameAccessAssociationRequest extends $pb.GeneratedMessage {
  factory DeleteDomainNameAccessAssociationRequest({
    $core.String? domainnameaccessassociationarn,
  }) {
    final result = create();
    if (domainnameaccessassociationarn != null)
      result.domainnameaccessassociationarn = domainnameaccessassociationarn;
    return result;
  }

  DeleteDomainNameAccessAssociationRequest._();

  factory DeleteDomainNameAccessAssociationRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDomainNameAccessAssociationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDomainNameAccessAssociationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(281017927, _omitFieldNames ? '' : 'domainnameaccessassociationarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDomainNameAccessAssociationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDomainNameAccessAssociationRequest copyWith(
          void Function(DeleteDomainNameAccessAssociationRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteDomainNameAccessAssociationRequest))
          as DeleteDomainNameAccessAssociationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDomainNameAccessAssociationRequest create() =>
      DeleteDomainNameAccessAssociationRequest._();
  @$core.override
  DeleteDomainNameAccessAssociationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDomainNameAccessAssociationRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeleteDomainNameAccessAssociationRequest>(create);
  static DeleteDomainNameAccessAssociationRequest? _defaultInstance;

  @$pb.TagNumber(281017927)
  $core.String get domainnameaccessassociationarn => $_getSZ(0);
  @$pb.TagNumber(281017927)
  set domainnameaccessassociationarn($core.String value) =>
      $_setString(0, value);
  @$pb.TagNumber(281017927)
  $core.bool hasDomainnameaccessassociationarn() => $_has(0);
  @$pb.TagNumber(281017927)
  void clearDomainnameaccessassociationarn() => $_clearField(281017927);
}

class DeleteDomainNameRequest extends $pb.GeneratedMessage {
  factory DeleteDomainNameRequest({
    $core.String? domainnameid,
    $core.String? domainname,
  }) {
    final result = create();
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  DeleteDomainNameRequest._();

  factory DeleteDomainNameRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDomainNameRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDomainNameRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDomainNameRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDomainNameRequest copyWith(
          void Function(DeleteDomainNameRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteDomainNameRequest))
          as DeleteDomainNameRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDomainNameRequest create() => DeleteDomainNameRequest._();
  @$core.override
  DeleteDomainNameRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDomainNameRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDomainNameRequest>(create);
  static DeleteDomainNameRequest? _defaultInstance;

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(0);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(0);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(1);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(1);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);
}

class DeleteGatewayResponseRequest extends $pb.GeneratedMessage {
  factory DeleteGatewayResponseRequest({
    GatewayResponseType? responsetype,
    $core.String? restapiid,
  }) {
    final result = create();
    if (responsetype != null) result.responsetype = responsetype;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteGatewayResponseRequest._();

  factory DeleteGatewayResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteGatewayResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteGatewayResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aE<GatewayResponseType>(377935935, _omitFieldNames ? '' : 'responsetype',
        enumValues: GatewayResponseType.values)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGatewayResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGatewayResponseRequest copyWith(
          void Function(DeleteGatewayResponseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteGatewayResponseRequest))
          as DeleteGatewayResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteGatewayResponseRequest create() =>
      DeleteGatewayResponseRequest._();
  @$core.override
  DeleteGatewayResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteGatewayResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteGatewayResponseRequest>(create);
  static DeleteGatewayResponseRequest? _defaultInstance;

  @$pb.TagNumber(377935935)
  GatewayResponseType get responsetype => $_getN(0);
  @$pb.TagNumber(377935935)
  set responsetype(GatewayResponseType value) => $_setField(377935935, value);
  @$pb.TagNumber(377935935)
  $core.bool hasResponsetype() => $_has(0);
  @$pb.TagNumber(377935935)
  void clearResponsetype() => $_clearField(377935935);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteIntegrationRequest extends $pb.GeneratedMessage {
  factory DeleteIntegrationRequest({
    $core.String? httpmethod,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteIntegrationRequest._();

  factory DeleteIntegrationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteIntegrationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteIntegrationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIntegrationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIntegrationRequest copyWith(
          void Function(DeleteIntegrationRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteIntegrationRequest))
          as DeleteIntegrationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteIntegrationRequest create() => DeleteIntegrationRequest._();
  @$core.override
  DeleteIntegrationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteIntegrationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteIntegrationRequest>(create);
  static DeleteIntegrationRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteIntegrationResponseRequest extends $pb.GeneratedMessage {
  factory DeleteIntegrationResponseRequest({
    $core.String? httpmethod,
    $core.String? statuscode,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (statuscode != null) result.statuscode = statuscode;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteIntegrationResponseRequest._();

  factory DeleteIntegrationResponseRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteIntegrationResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteIntegrationResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIntegrationResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIntegrationResponseRequest copyWith(
          void Function(DeleteIntegrationResponseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteIntegrationResponseRequest))
          as DeleteIntegrationResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteIntegrationResponseRequest create() =>
      DeleteIntegrationResponseRequest._();
  @$core.override
  DeleteIntegrationResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteIntegrationResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteIntegrationResponseRequest>(
          create);
  static DeleteIntegrationResponseRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(1);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(1);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(2);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(2);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteMethodRequest extends $pb.GeneratedMessage {
  factory DeleteMethodRequest({
    $core.String? httpmethod,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteMethodRequest._();

  factory DeleteMethodRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMethodRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMethodRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMethodRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMethodRequest copyWith(void Function(DeleteMethodRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteMethodRequest))
          as DeleteMethodRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMethodRequest create() => DeleteMethodRequest._();
  @$core.override
  DeleteMethodRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMethodRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMethodRequest>(create);
  static DeleteMethodRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteMethodResponseRequest extends $pb.GeneratedMessage {
  factory DeleteMethodResponseRequest({
    $core.String? httpmethod,
    $core.String? statuscode,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (statuscode != null) result.statuscode = statuscode;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteMethodResponseRequest._();

  factory DeleteMethodResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMethodResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMethodResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMethodResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMethodResponseRequest copyWith(
          void Function(DeleteMethodResponseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteMethodResponseRequest))
          as DeleteMethodResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMethodResponseRequest create() =>
      DeleteMethodResponseRequest._();
  @$core.override
  DeleteMethodResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMethodResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMethodResponseRequest>(create);
  static DeleteMethodResponseRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(1);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(1);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(2);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(2);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteModelRequest extends $pb.GeneratedMessage {
  factory DeleteModelRequest({
    $core.String? modelname,
    $core.String? restapiid,
  }) {
    final result = create();
    if (modelname != null) result.modelname = modelname;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteModelRequest._();

  factory DeleteModelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteModelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteModelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(176835330, _omitFieldNames ? '' : 'modelname')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteModelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteModelRequest copyWith(void Function(DeleteModelRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteModelRequest))
          as DeleteModelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteModelRequest create() => DeleteModelRequest._();
  @$core.override
  DeleteModelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteModelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteModelRequest>(create);
  static DeleteModelRequest? _defaultInstance;

  @$pb.TagNumber(176835330)
  $core.String get modelname => $_getSZ(0);
  @$pb.TagNumber(176835330)
  set modelname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(176835330)
  $core.bool hasModelname() => $_has(0);
  @$pb.TagNumber(176835330)
  void clearModelname() => $_clearField(176835330);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteRequestValidatorRequest extends $pb.GeneratedMessage {
  factory DeleteRequestValidatorRequest({
    $core.String? restapiid,
    $core.String? requestvalidatorid,
  }) {
    final result = create();
    if (restapiid != null) result.restapiid = restapiid;
    if (requestvalidatorid != null)
      result.requestvalidatorid = requestvalidatorid;
    return result;
  }

  DeleteRequestValidatorRequest._();

  factory DeleteRequestValidatorRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRequestValidatorRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRequestValidatorRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(517546134, _omitFieldNames ? '' : 'requestvalidatorid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRequestValidatorRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRequestValidatorRequest copyWith(
          void Function(DeleteRequestValidatorRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteRequestValidatorRequest))
          as DeleteRequestValidatorRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRequestValidatorRequest create() =>
      DeleteRequestValidatorRequest._();
  @$core.override
  DeleteRequestValidatorRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRequestValidatorRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRequestValidatorRequest>(create);
  static DeleteRequestValidatorRequest? _defaultInstance;

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(0);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(0);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(517546134)
  $core.String get requestvalidatorid => $_getSZ(1);
  @$pb.TagNumber(517546134)
  set requestvalidatorid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(517546134)
  $core.bool hasRequestvalidatorid() => $_has(1);
  @$pb.TagNumber(517546134)
  void clearRequestvalidatorid() => $_clearField(517546134);
}

class DeleteResourceRequest extends $pb.GeneratedMessage {
  factory DeleteResourceRequest({
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteResourceRequest._();

  factory DeleteResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteResourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourceRequest copyWith(
          void Function(DeleteResourceRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteResourceRequest))
          as DeleteResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteResourceRequest create() => DeleteResourceRequest._();
  @$core.override
  DeleteResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteResourceRequest>(create);
  static DeleteResourceRequest? _defaultInstance;

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(0);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(0);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteRestApiRequest extends $pb.GeneratedMessage {
  factory DeleteRestApiRequest({
    $core.String? restapiid,
  }) {
    final result = create();
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteRestApiRequest._();

  factory DeleteRestApiRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRestApiRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRestApiRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRestApiRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRestApiRequest copyWith(void Function(DeleteRestApiRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteRestApiRequest))
          as DeleteRestApiRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRestApiRequest create() => DeleteRestApiRequest._();
  @$core.override
  DeleteRestApiRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRestApiRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRestApiRequest>(create);
  static DeleteRestApiRequest? _defaultInstance;

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(0);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(0);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteStageRequest extends $pb.GeneratedMessage {
  factory DeleteStageRequest({
    $core.String? stagename,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  DeleteStageRequest._();

  factory DeleteStageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteStageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteStageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteStageRequest copyWith(void Function(DeleteStageRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteStageRequest))
          as DeleteStageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteStageRequest create() => DeleteStageRequest._();
  @$core.override
  DeleteStageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteStageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteStageRequest>(create);
  static DeleteStageRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class DeleteUsagePlanKeyRequest extends $pb.GeneratedMessage {
  factory DeleteUsagePlanKeyRequest({
    $core.String? keyid,
    $core.String? usageplanid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  DeleteUsagePlanKeyRequest._();

  factory DeleteUsagePlanKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteUsagePlanKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteUsagePlanKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(479913282, _omitFieldNames ? '' : 'keyid')
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteUsagePlanKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteUsagePlanKeyRequest copyWith(
          void Function(DeleteUsagePlanKeyRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteUsagePlanKeyRequest))
          as DeleteUsagePlanKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteUsagePlanKeyRequest create() => DeleteUsagePlanKeyRequest._();
  @$core.override
  DeleteUsagePlanKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteUsagePlanKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteUsagePlanKeyRequest>(create);
  static DeleteUsagePlanKeyRequest? _defaultInstance;

  @$pb.TagNumber(479913282)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(479913282)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(479913282)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(479913282)
  void clearKeyid() => $_clearField(479913282);

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(1);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(1);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class DeleteUsagePlanRequest extends $pb.GeneratedMessage {
  factory DeleteUsagePlanRequest({
    $core.String? usageplanid,
  }) {
    final result = create();
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  DeleteUsagePlanRequest._();

  factory DeleteUsagePlanRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteUsagePlanRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteUsagePlanRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteUsagePlanRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteUsagePlanRequest copyWith(
          void Function(DeleteUsagePlanRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteUsagePlanRequest))
          as DeleteUsagePlanRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteUsagePlanRequest create() => DeleteUsagePlanRequest._();
  @$core.override
  DeleteUsagePlanRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteUsagePlanRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteUsagePlanRequest>(create);
  static DeleteUsagePlanRequest? _defaultInstance;

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(0);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(0);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class DeleteVpcLinkRequest extends $pb.GeneratedMessage {
  factory DeleteVpcLinkRequest({
    $core.String? vpclinkid,
  }) {
    final result = create();
    if (vpclinkid != null) result.vpclinkid = vpclinkid;
    return result;
  }

  DeleteVpcLinkRequest._();

  factory DeleteVpcLinkRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteVpcLinkRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteVpcLinkRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(27515846, _omitFieldNames ? '' : 'vpclinkid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteVpcLinkRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteVpcLinkRequest copyWith(void Function(DeleteVpcLinkRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteVpcLinkRequest))
          as DeleteVpcLinkRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteVpcLinkRequest create() => DeleteVpcLinkRequest._();
  @$core.override
  DeleteVpcLinkRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteVpcLinkRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteVpcLinkRequest>(create);
  static DeleteVpcLinkRequest? _defaultInstance;

  @$pb.TagNumber(27515846)
  $core.String get vpclinkid => $_getSZ(0);
  @$pb.TagNumber(27515846)
  set vpclinkid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(27515846)
  $core.bool hasVpclinkid() => $_has(0);
  @$pb.TagNumber(27515846)
  void clearVpclinkid() => $_clearField(27515846);
}

class Deployment extends $pb.GeneratedMessage {
  factory Deployment({
    $core.String? createddate,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? apisummary,
    $core.String? description,
    $core.String? id,
  }) {
    final result = create();
    if (createddate != null) result.createddate = createddate;
    if (apisummary != null) result.apisummary.addEntries(apisummary);
    if (description != null) result.description = description;
    if (id != null) result.id = id;
    return result;
  }

  Deployment._();

  factory Deployment.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Deployment.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Deployment',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(53061200, _omitFieldNames ? '' : 'createddate')
    ..m<$core.String, $core.String>(
        159675170, _omitFieldNames ? '' : 'apisummary',
        entryClassName: 'Deployment.ApisummaryEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Deployment clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Deployment copyWith(void Function(Deployment) updates) =>
      super.copyWith((message) => updates(message as Deployment)) as Deployment;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Deployment create() => Deployment._();
  @$core.override
  Deployment createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Deployment getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Deployment>(create);
  static Deployment? _defaultInstance;

  @$pb.TagNumber(53061200)
  $core.String get createddate => $_getSZ(0);
  @$pb.TagNumber(53061200)
  set createddate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(53061200)
  $core.bool hasCreateddate() => $_has(0);
  @$pb.TagNumber(53061200)
  void clearCreateddate() => $_clearField(53061200);

  @$pb.TagNumber(159675170)
  $pb.PbMap<$core.String, $core.String> get apisummary => $_getMap(1);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);
}

class DeploymentCanarySettings extends $pb.GeneratedMessage {
  factory DeploymentCanarySettings({
    $core.double? percenttraffic,
    $core.bool? usestagecache,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        stagevariableoverrides,
  }) {
    final result = create();
    if (percenttraffic != null) result.percenttraffic = percenttraffic;
    if (usestagecache != null) result.usestagecache = usestagecache;
    if (stagevariableoverrides != null)
      result.stagevariableoverrides.addEntries(stagevariableoverrides);
    return result;
  }

  DeploymentCanarySettings._();

  factory DeploymentCanarySettings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeploymentCanarySettings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeploymentCanarySettings',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aD(147177864, _omitFieldNames ? '' : 'percenttraffic')
    ..aOB(179841697, _omitFieldNames ? '' : 'usestagecache')
    ..m<$core.String, $core.String>(
        221124259, _omitFieldNames ? '' : 'stagevariableoverrides',
        entryClassName: 'DeploymentCanarySettings.StagevariableoverridesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeploymentCanarySettings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeploymentCanarySettings copyWith(
          void Function(DeploymentCanarySettings) updates) =>
      super.copyWith((message) => updates(message as DeploymentCanarySettings))
          as DeploymentCanarySettings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeploymentCanarySettings create() => DeploymentCanarySettings._();
  @$core.override
  DeploymentCanarySettings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeploymentCanarySettings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeploymentCanarySettings>(create);
  static DeploymentCanarySettings? _defaultInstance;

  @$pb.TagNumber(147177864)
  $core.double get percenttraffic => $_getN(0);
  @$pb.TagNumber(147177864)
  set percenttraffic($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(147177864)
  $core.bool hasPercenttraffic() => $_has(0);
  @$pb.TagNumber(147177864)
  void clearPercenttraffic() => $_clearField(147177864);

  @$pb.TagNumber(179841697)
  $core.bool get usestagecache => $_getBF(1);
  @$pb.TagNumber(179841697)
  set usestagecache($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(179841697)
  $core.bool hasUsestagecache() => $_has(1);
  @$pb.TagNumber(179841697)
  void clearUsestagecache() => $_clearField(179841697);

  @$pb.TagNumber(221124259)
  $pb.PbMap<$core.String, $core.String> get stagevariableoverrides =>
      $_getMap(2);
}

class Deployments extends $pb.GeneratedMessage {
  factory Deployments({
    $core.String? position,
    $core.Iterable<Deployment>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  Deployments._();

  factory Deployments.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Deployments.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Deployments',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<Deployment>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: Deployment.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Deployments clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Deployments copyWith(void Function(Deployments) updates) =>
      super.copyWith((message) => updates(message as Deployments))
          as Deployments;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Deployments create() => Deployments._();
  @$core.override
  Deployments createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Deployments getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Deployments>(create);
  static Deployments? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<Deployment> get items => $_getList(1);
}

class DocumentationPart extends $pb.GeneratedMessage {
  factory DocumentationPart({
    DocumentationPartLocation? location,
    $core.String? properties,
    $core.String? id,
  }) {
    final result = create();
    if (location != null) result.location = location;
    if (properties != null) result.properties = properties;
    if (id != null) result.id = id;
    return result;
  }

  DocumentationPart._();

  factory DocumentationPart.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DocumentationPart.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DocumentationPart',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOM<DocumentationPartLocation>(
        200649127, _omitFieldNames ? '' : 'location',
        subBuilder: DocumentationPartLocation.create)
    ..aOS(299789533, _omitFieldNames ? '' : 'properties')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationPart clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationPart copyWith(void Function(DocumentationPart) updates) =>
      super.copyWith((message) => updates(message as DocumentationPart))
          as DocumentationPart;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DocumentationPart create() => DocumentationPart._();
  @$core.override
  DocumentationPart createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DocumentationPart getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DocumentationPart>(create);
  static DocumentationPart? _defaultInstance;

  @$pb.TagNumber(200649127)
  DocumentationPartLocation get location => $_getN(0);
  @$pb.TagNumber(200649127)
  set location(DocumentationPartLocation value) => $_setField(200649127, value);
  @$pb.TagNumber(200649127)
  $core.bool hasLocation() => $_has(0);
  @$pb.TagNumber(200649127)
  void clearLocation() => $_clearField(200649127);
  @$pb.TagNumber(200649127)
  DocumentationPartLocation ensureLocation() => $_ensure(0);

  @$pb.TagNumber(299789533)
  $core.String get properties => $_getSZ(1);
  @$pb.TagNumber(299789533)
  set properties($core.String value) => $_setString(1, value);
  @$pb.TagNumber(299789533)
  $core.bool hasProperties() => $_has(1);
  @$pb.TagNumber(299789533)
  void clearProperties() => $_clearField(299789533);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);
}

class DocumentationPartIds extends $pb.GeneratedMessage {
  factory DocumentationPartIds({
    $core.Iterable<$core.String>? ids,
    $core.Iterable<$core.String>? warnings,
  }) {
    final result = create();
    if (ids != null) result.ids.addAll(ids);
    if (warnings != null) result.warnings.addAll(warnings);
    return result;
  }

  DocumentationPartIds._();

  factory DocumentationPartIds.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DocumentationPartIds.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DocumentationPartIds',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(13616490, _omitFieldNames ? '' : 'ids')
    ..pPS(185617301, _omitFieldNames ? '' : 'warnings')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationPartIds clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationPartIds copyWith(void Function(DocumentationPartIds) updates) =>
      super.copyWith((message) => updates(message as DocumentationPartIds))
          as DocumentationPartIds;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DocumentationPartIds create() => DocumentationPartIds._();
  @$core.override
  DocumentationPartIds createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DocumentationPartIds getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DocumentationPartIds>(create);
  static DocumentationPartIds? _defaultInstance;

  @$pb.TagNumber(13616490)
  $pb.PbList<$core.String> get ids => $_getList(0);

  @$pb.TagNumber(185617301)
  $pb.PbList<$core.String> get warnings => $_getList(1);
}

class DocumentationPartLocation extends $pb.GeneratedMessage {
  factory DocumentationPartLocation({
    $core.String? path,
    $core.String? method,
    $core.String? name,
    DocumentationPartType? type,
    $core.String? statuscode,
  }) {
    final result = create();
    if (path != null) result.path = path;
    if (method != null) result.method = method;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (statuscode != null) result.statuscode = statuscode;
    return result;
  }

  DocumentationPartLocation._();

  factory DocumentationPartLocation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DocumentationPartLocation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DocumentationPartLocation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(75975991, _omitFieldNames ? '' : 'path')
    ..aOS(189134641, _omitFieldNames ? '' : 'method')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aE<DocumentationPartType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: DocumentationPartType.values)
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationPartLocation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationPartLocation copyWith(
          void Function(DocumentationPartLocation) updates) =>
      super.copyWith((message) => updates(message as DocumentationPartLocation))
          as DocumentationPartLocation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DocumentationPartLocation create() => DocumentationPartLocation._();
  @$core.override
  DocumentationPartLocation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DocumentationPartLocation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DocumentationPartLocation>(create);
  static DocumentationPartLocation? _defaultInstance;

  @$pb.TagNumber(75975991)
  $core.String get path => $_getSZ(0);
  @$pb.TagNumber(75975991)
  set path($core.String value) => $_setString(0, value);
  @$pb.TagNumber(75975991)
  $core.bool hasPath() => $_has(0);
  @$pb.TagNumber(75975991)
  void clearPath() => $_clearField(75975991);

  @$pb.TagNumber(189134641)
  $core.String get method => $_getSZ(1);
  @$pb.TagNumber(189134641)
  set method($core.String value) => $_setString(1, value);
  @$pb.TagNumber(189134641)
  $core.bool hasMethod() => $_has(1);
  @$pb.TagNumber(189134641)
  void clearMethod() => $_clearField(189134641);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(287830350)
  DocumentationPartType get type => $_getN(3);
  @$pb.TagNumber(287830350)
  set type(DocumentationPartType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(3);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(4);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(4, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(4);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);
}

class DocumentationParts extends $pb.GeneratedMessage {
  factory DocumentationParts({
    $core.String? position,
    $core.Iterable<DocumentationPart>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  DocumentationParts._();

  factory DocumentationParts.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DocumentationParts.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DocumentationParts',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<DocumentationPart>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: DocumentationPart.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationParts clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationParts copyWith(void Function(DocumentationParts) updates) =>
      super.copyWith((message) => updates(message as DocumentationParts))
          as DocumentationParts;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DocumentationParts create() => DocumentationParts._();
  @$core.override
  DocumentationParts createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DocumentationParts getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DocumentationParts>(create);
  static DocumentationParts? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<DocumentationPart> get items => $_getList(1);
}

class DocumentationVersion extends $pb.GeneratedMessage {
  factory DocumentationVersion({
    $core.String? createddate,
    $core.String? version,
    $core.String? description,
  }) {
    final result = create();
    if (createddate != null) result.createddate = createddate;
    if (version != null) result.version = version;
    if (description != null) result.description = description;
    return result;
  }

  DocumentationVersion._();

  factory DocumentationVersion.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DocumentationVersion.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DocumentationVersion',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(53061200, _omitFieldNames ? '' : 'createddate')
    ..aOS(108113560, _omitFieldNames ? '' : 'version')
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationVersion clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationVersion copyWith(void Function(DocumentationVersion) updates) =>
      super.copyWith((message) => updates(message as DocumentationVersion))
          as DocumentationVersion;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DocumentationVersion create() => DocumentationVersion._();
  @$core.override
  DocumentationVersion createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DocumentationVersion getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DocumentationVersion>(create);
  static DocumentationVersion? _defaultInstance;

  @$pb.TagNumber(53061200)
  $core.String get createddate => $_getSZ(0);
  @$pb.TagNumber(53061200)
  set createddate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(53061200)
  $core.bool hasCreateddate() => $_has(0);
  @$pb.TagNumber(53061200)
  void clearCreateddate() => $_clearField(53061200);

  @$pb.TagNumber(108113560)
  $core.String get version => $_getSZ(1);
  @$pb.TagNumber(108113560)
  set version($core.String value) => $_setString(1, value);
  @$pb.TagNumber(108113560)
  $core.bool hasVersion() => $_has(1);
  @$pb.TagNumber(108113560)
  void clearVersion() => $_clearField(108113560);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);
}

class DocumentationVersions extends $pb.GeneratedMessage {
  factory DocumentationVersions({
    $core.String? position,
    $core.Iterable<DocumentationVersion>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  DocumentationVersions._();

  factory DocumentationVersions.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DocumentationVersions.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DocumentationVersions',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<DocumentationVersion>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: DocumentationVersion.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationVersions clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DocumentationVersions copyWith(
          void Function(DocumentationVersions) updates) =>
      super.copyWith((message) => updates(message as DocumentationVersions))
          as DocumentationVersions;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DocumentationVersions create() => DocumentationVersions._();
  @$core.override
  DocumentationVersions createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DocumentationVersions getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DocumentationVersions>(create);
  static DocumentationVersions? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<DocumentationVersion> get items => $_getList(1);
}

class DomainName extends $pb.GeneratedMessage {
  factory DomainName({
    $core.String? managementpolicy,
    MutualTlsAuthentication? mutualtlsauthentication,
    $core.String? ownershipverificationcertificatearn,
    $core.String? regionalcertificatearn,
    $core.String? certificatename,
    $core.String? distributionhostedzoneid,
    $core.String? regionaldomainname,
    $core.String? domainnamearn,
    $core.String? policy,
    $core.String? distributiondomainname,
    DomainNameStatus? domainnamestatus,
    $core.String? domainnameid,
    $core.String? domainnamestatusmessage,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    EndpointAccessMode? endpointaccessmode,
    $core.String? domainname,
    $core.String? certificatearn,
    $core.String? regionalcertificatename,
    $core.String? regionalhostedzoneid,
    EndpointConfiguration? endpointconfiguration,
    SecurityPolicy? securitypolicy,
    $core.String? certificateuploaddate,
    RoutingMode? routingmode,
  }) {
    final result = create();
    if (managementpolicy != null) result.managementpolicy = managementpolicy;
    if (mutualtlsauthentication != null)
      result.mutualtlsauthentication = mutualtlsauthentication;
    if (ownershipverificationcertificatearn != null)
      result.ownershipverificationcertificatearn =
          ownershipverificationcertificatearn;
    if (regionalcertificatearn != null)
      result.regionalcertificatearn = regionalcertificatearn;
    if (certificatename != null) result.certificatename = certificatename;
    if (distributionhostedzoneid != null)
      result.distributionhostedzoneid = distributionhostedzoneid;
    if (regionaldomainname != null)
      result.regionaldomainname = regionaldomainname;
    if (domainnamearn != null) result.domainnamearn = domainnamearn;
    if (policy != null) result.policy = policy;
    if (distributiondomainname != null)
      result.distributiondomainname = distributiondomainname;
    if (domainnamestatus != null) result.domainnamestatus = domainnamestatus;
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (domainnamestatusmessage != null)
      result.domainnamestatusmessage = domainnamestatusmessage;
    if (tags != null) result.tags.addEntries(tags);
    if (endpointaccessmode != null)
      result.endpointaccessmode = endpointaccessmode;
    if (domainname != null) result.domainname = domainname;
    if (certificatearn != null) result.certificatearn = certificatearn;
    if (regionalcertificatename != null)
      result.regionalcertificatename = regionalcertificatename;
    if (regionalhostedzoneid != null)
      result.regionalhostedzoneid = regionalhostedzoneid;
    if (endpointconfiguration != null)
      result.endpointconfiguration = endpointconfiguration;
    if (securitypolicy != null) result.securitypolicy = securitypolicy;
    if (certificateuploaddate != null)
      result.certificateuploaddate = certificateuploaddate;
    if (routingmode != null) result.routingmode = routingmode;
    return result;
  }

  DomainName._();

  factory DomainName.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DomainName.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DomainName',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(90103737, _omitFieldNames ? '' : 'managementpolicy')
    ..aOM<MutualTlsAuthentication>(
        99462043, _omitFieldNames ? '' : 'mutualtlsauthentication',
        subBuilder: MutualTlsAuthentication.create)
    ..aOS(
        100550548, _omitFieldNames ? '' : 'ownershipverificationcertificatearn')
    ..aOS(137066579, _omitFieldNames ? '' : 'regionalcertificatearn')
    ..aOS(140276948, _omitFieldNames ? '' : 'certificatename')
    ..aOS(185867208, _omitFieldNames ? '' : 'distributionhostedzoneid')
    ..aOS(198256560, _omitFieldNames ? '' : 'regionaldomainname')
    ..aOS(244019094, _omitFieldNames ? '' : 'domainnamearn')
    ..aOS(247528064, _omitFieldNames ? '' : 'policy')
    ..aOS(266229213, _omitFieldNames ? '' : 'distributiondomainname')
    ..aE<DomainNameStatus>(275080629, _omitFieldNames ? '' : 'domainnamestatus',
        enumValues: DomainNameStatus.values)
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aOS(335045044, _omitFieldNames ? '' : 'domainnamestatusmessage')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'DomainName.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aE<EndpointAccessMode>(
        356705630, _omitFieldNames ? '' : 'endpointaccessmode',
        enumValues: EndpointAccessMode.values)
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..aOS(425831704, _omitFieldNames ? '' : 'certificatearn')
    ..aOS(431716941, _omitFieldNames ? '' : 'regionalcertificatename')
    ..aOS(467276949, _omitFieldNames ? '' : 'regionalhostedzoneid')
    ..aOM<EndpointConfiguration>(
        487543735, _omitFieldNames ? '' : 'endpointconfiguration',
        subBuilder: EndpointConfiguration.create)
    ..aE<SecurityPolicy>(491792990, _omitFieldNames ? '' : 'securitypolicy',
        enumValues: SecurityPolicy.values)
    ..aOS(504466814, _omitFieldNames ? '' : 'certificateuploaddate')
    ..aE<RoutingMode>(506342119, _omitFieldNames ? '' : 'routingmode',
        enumValues: RoutingMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainName clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainName copyWith(void Function(DomainName) updates) =>
      super.copyWith((message) => updates(message as DomainName)) as DomainName;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DomainName create() => DomainName._();
  @$core.override
  DomainName createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DomainName getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DomainName>(create);
  static DomainName? _defaultInstance;

  @$pb.TagNumber(90103737)
  $core.String get managementpolicy => $_getSZ(0);
  @$pb.TagNumber(90103737)
  set managementpolicy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(90103737)
  $core.bool hasManagementpolicy() => $_has(0);
  @$pb.TagNumber(90103737)
  void clearManagementpolicy() => $_clearField(90103737);

  @$pb.TagNumber(99462043)
  MutualTlsAuthentication get mutualtlsauthentication => $_getN(1);
  @$pb.TagNumber(99462043)
  set mutualtlsauthentication(MutualTlsAuthentication value) =>
      $_setField(99462043, value);
  @$pb.TagNumber(99462043)
  $core.bool hasMutualtlsauthentication() => $_has(1);
  @$pb.TagNumber(99462043)
  void clearMutualtlsauthentication() => $_clearField(99462043);
  @$pb.TagNumber(99462043)
  MutualTlsAuthentication ensureMutualtlsauthentication() => $_ensure(1);

  @$pb.TagNumber(100550548)
  $core.String get ownershipverificationcertificatearn => $_getSZ(2);
  @$pb.TagNumber(100550548)
  set ownershipverificationcertificatearn($core.String value) =>
      $_setString(2, value);
  @$pb.TagNumber(100550548)
  $core.bool hasOwnershipverificationcertificatearn() => $_has(2);
  @$pb.TagNumber(100550548)
  void clearOwnershipverificationcertificatearn() => $_clearField(100550548);

  @$pb.TagNumber(137066579)
  $core.String get regionalcertificatearn => $_getSZ(3);
  @$pb.TagNumber(137066579)
  set regionalcertificatearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(137066579)
  $core.bool hasRegionalcertificatearn() => $_has(3);
  @$pb.TagNumber(137066579)
  void clearRegionalcertificatearn() => $_clearField(137066579);

  @$pb.TagNumber(140276948)
  $core.String get certificatename => $_getSZ(4);
  @$pb.TagNumber(140276948)
  set certificatename($core.String value) => $_setString(4, value);
  @$pb.TagNumber(140276948)
  $core.bool hasCertificatename() => $_has(4);
  @$pb.TagNumber(140276948)
  void clearCertificatename() => $_clearField(140276948);

  @$pb.TagNumber(185867208)
  $core.String get distributionhostedzoneid => $_getSZ(5);
  @$pb.TagNumber(185867208)
  set distributionhostedzoneid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(185867208)
  $core.bool hasDistributionhostedzoneid() => $_has(5);
  @$pb.TagNumber(185867208)
  void clearDistributionhostedzoneid() => $_clearField(185867208);

  @$pb.TagNumber(198256560)
  $core.String get regionaldomainname => $_getSZ(6);
  @$pb.TagNumber(198256560)
  set regionaldomainname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(198256560)
  $core.bool hasRegionaldomainname() => $_has(6);
  @$pb.TagNumber(198256560)
  void clearRegionaldomainname() => $_clearField(198256560);

  @$pb.TagNumber(244019094)
  $core.String get domainnamearn => $_getSZ(7);
  @$pb.TagNumber(244019094)
  set domainnamearn($core.String value) => $_setString(7, value);
  @$pb.TagNumber(244019094)
  $core.bool hasDomainnamearn() => $_has(7);
  @$pb.TagNumber(244019094)
  void clearDomainnamearn() => $_clearField(244019094);

  @$pb.TagNumber(247528064)
  $core.String get policy => $_getSZ(8);
  @$pb.TagNumber(247528064)
  set policy($core.String value) => $_setString(8, value);
  @$pb.TagNumber(247528064)
  $core.bool hasPolicy() => $_has(8);
  @$pb.TagNumber(247528064)
  void clearPolicy() => $_clearField(247528064);

  @$pb.TagNumber(266229213)
  $core.String get distributiondomainname => $_getSZ(9);
  @$pb.TagNumber(266229213)
  set distributiondomainname($core.String value) => $_setString(9, value);
  @$pb.TagNumber(266229213)
  $core.bool hasDistributiondomainname() => $_has(9);
  @$pb.TagNumber(266229213)
  void clearDistributiondomainname() => $_clearField(266229213);

  @$pb.TagNumber(275080629)
  DomainNameStatus get domainnamestatus => $_getN(10);
  @$pb.TagNumber(275080629)
  set domainnamestatus(DomainNameStatus value) => $_setField(275080629, value);
  @$pb.TagNumber(275080629)
  $core.bool hasDomainnamestatus() => $_has(10);
  @$pb.TagNumber(275080629)
  void clearDomainnamestatus() => $_clearField(275080629);

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(11);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(11, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(11);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(335045044)
  $core.String get domainnamestatusmessage => $_getSZ(12);
  @$pb.TagNumber(335045044)
  set domainnamestatusmessage($core.String value) => $_setString(12, value);
  @$pb.TagNumber(335045044)
  $core.bool hasDomainnamestatusmessage() => $_has(12);
  @$pb.TagNumber(335045044)
  void clearDomainnamestatusmessage() => $_clearField(335045044);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(13);

  @$pb.TagNumber(356705630)
  EndpointAccessMode get endpointaccessmode => $_getN(14);
  @$pb.TagNumber(356705630)
  set endpointaccessmode(EndpointAccessMode value) =>
      $_setField(356705630, value);
  @$pb.TagNumber(356705630)
  $core.bool hasEndpointaccessmode() => $_has(14);
  @$pb.TagNumber(356705630)
  void clearEndpointaccessmode() => $_clearField(356705630);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(15);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(15, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(15);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);

  @$pb.TagNumber(425831704)
  $core.String get certificatearn => $_getSZ(16);
  @$pb.TagNumber(425831704)
  set certificatearn($core.String value) => $_setString(16, value);
  @$pb.TagNumber(425831704)
  $core.bool hasCertificatearn() => $_has(16);
  @$pb.TagNumber(425831704)
  void clearCertificatearn() => $_clearField(425831704);

  @$pb.TagNumber(431716941)
  $core.String get regionalcertificatename => $_getSZ(17);
  @$pb.TagNumber(431716941)
  set regionalcertificatename($core.String value) => $_setString(17, value);
  @$pb.TagNumber(431716941)
  $core.bool hasRegionalcertificatename() => $_has(17);
  @$pb.TagNumber(431716941)
  void clearRegionalcertificatename() => $_clearField(431716941);

  @$pb.TagNumber(467276949)
  $core.String get regionalhostedzoneid => $_getSZ(18);
  @$pb.TagNumber(467276949)
  set regionalhostedzoneid($core.String value) => $_setString(18, value);
  @$pb.TagNumber(467276949)
  $core.bool hasRegionalhostedzoneid() => $_has(18);
  @$pb.TagNumber(467276949)
  void clearRegionalhostedzoneid() => $_clearField(467276949);

  @$pb.TagNumber(487543735)
  EndpointConfiguration get endpointconfiguration => $_getN(19);
  @$pb.TagNumber(487543735)
  set endpointconfiguration(EndpointConfiguration value) =>
      $_setField(487543735, value);
  @$pb.TagNumber(487543735)
  $core.bool hasEndpointconfiguration() => $_has(19);
  @$pb.TagNumber(487543735)
  void clearEndpointconfiguration() => $_clearField(487543735);
  @$pb.TagNumber(487543735)
  EndpointConfiguration ensureEndpointconfiguration() => $_ensure(19);

  @$pb.TagNumber(491792990)
  SecurityPolicy get securitypolicy => $_getN(20);
  @$pb.TagNumber(491792990)
  set securitypolicy(SecurityPolicy value) => $_setField(491792990, value);
  @$pb.TagNumber(491792990)
  $core.bool hasSecuritypolicy() => $_has(20);
  @$pb.TagNumber(491792990)
  void clearSecuritypolicy() => $_clearField(491792990);

  @$pb.TagNumber(504466814)
  $core.String get certificateuploaddate => $_getSZ(21);
  @$pb.TagNumber(504466814)
  set certificateuploaddate($core.String value) => $_setString(21, value);
  @$pb.TagNumber(504466814)
  $core.bool hasCertificateuploaddate() => $_has(21);
  @$pb.TagNumber(504466814)
  void clearCertificateuploaddate() => $_clearField(504466814);

  @$pb.TagNumber(506342119)
  RoutingMode get routingmode => $_getN(22);
  @$pb.TagNumber(506342119)
  set routingmode(RoutingMode value) => $_setField(506342119, value);
  @$pb.TagNumber(506342119)
  $core.bool hasRoutingmode() => $_has(22);
  @$pb.TagNumber(506342119)
  void clearRoutingmode() => $_clearField(506342119);
}

class DomainNameAccessAssociation extends $pb.GeneratedMessage {
  factory DomainNameAccessAssociation({
    AccessAssociationSourceType? accessassociationsourcetype,
    $core.String? domainnamearn,
    $core.String? domainnameaccessassociationarn,
    $core.String? accessassociationsource,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
  }) {
    final result = create();
    if (accessassociationsourcetype != null)
      result.accessassociationsourcetype = accessassociationsourcetype;
    if (domainnamearn != null) result.domainnamearn = domainnamearn;
    if (domainnameaccessassociationarn != null)
      result.domainnameaccessassociationarn = domainnameaccessassociationarn;
    if (accessassociationsource != null)
      result.accessassociationsource = accessassociationsource;
    if (tags != null) result.tags.addEntries(tags);
    return result;
  }

  DomainNameAccessAssociation._();

  factory DomainNameAccessAssociation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DomainNameAccessAssociation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DomainNameAccessAssociation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aE<AccessAssociationSourceType>(
        176397628, _omitFieldNames ? '' : 'accessassociationsourcetype',
        enumValues: AccessAssociationSourceType.values)
    ..aOS(244019094, _omitFieldNames ? '' : 'domainnamearn')
    ..aOS(281017927, _omitFieldNames ? '' : 'domainnameaccessassociationarn')
    ..aOS(328257828, _omitFieldNames ? '' : 'accessassociationsource')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'DomainNameAccessAssociation.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainNameAccessAssociation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainNameAccessAssociation copyWith(
          void Function(DomainNameAccessAssociation) updates) =>
      super.copyWith(
              (message) => updates(message as DomainNameAccessAssociation))
          as DomainNameAccessAssociation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DomainNameAccessAssociation create() =>
      DomainNameAccessAssociation._();
  @$core.override
  DomainNameAccessAssociation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DomainNameAccessAssociation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DomainNameAccessAssociation>(create);
  static DomainNameAccessAssociation? _defaultInstance;

  @$pb.TagNumber(176397628)
  AccessAssociationSourceType get accessassociationsourcetype => $_getN(0);
  @$pb.TagNumber(176397628)
  set accessassociationsourcetype(AccessAssociationSourceType value) =>
      $_setField(176397628, value);
  @$pb.TagNumber(176397628)
  $core.bool hasAccessassociationsourcetype() => $_has(0);
  @$pb.TagNumber(176397628)
  void clearAccessassociationsourcetype() => $_clearField(176397628);

  @$pb.TagNumber(244019094)
  $core.String get domainnamearn => $_getSZ(1);
  @$pb.TagNumber(244019094)
  set domainnamearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(244019094)
  $core.bool hasDomainnamearn() => $_has(1);
  @$pb.TagNumber(244019094)
  void clearDomainnamearn() => $_clearField(244019094);

  @$pb.TagNumber(281017927)
  $core.String get domainnameaccessassociationarn => $_getSZ(2);
  @$pb.TagNumber(281017927)
  set domainnameaccessassociationarn($core.String value) =>
      $_setString(2, value);
  @$pb.TagNumber(281017927)
  $core.bool hasDomainnameaccessassociationarn() => $_has(2);
  @$pb.TagNumber(281017927)
  void clearDomainnameaccessassociationarn() => $_clearField(281017927);

  @$pb.TagNumber(328257828)
  $core.String get accessassociationsource => $_getSZ(3);
  @$pb.TagNumber(328257828)
  set accessassociationsource($core.String value) => $_setString(3, value);
  @$pb.TagNumber(328257828)
  $core.bool hasAccessassociationsource() => $_has(3);
  @$pb.TagNumber(328257828)
  void clearAccessassociationsource() => $_clearField(328257828);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(4);
}

class DomainNameAccessAssociations extends $pb.GeneratedMessage {
  factory DomainNameAccessAssociations({
    $core.String? position,
    $core.Iterable<DomainNameAccessAssociation>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  DomainNameAccessAssociations._();

  factory DomainNameAccessAssociations.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DomainNameAccessAssociations.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DomainNameAccessAssociations',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<DomainNameAccessAssociation>(
        444150672, _omitFieldNames ? '' : 'items',
        subBuilder: DomainNameAccessAssociation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainNameAccessAssociations clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainNameAccessAssociations copyWith(
          void Function(DomainNameAccessAssociations) updates) =>
      super.copyWith(
              (message) => updates(message as DomainNameAccessAssociations))
          as DomainNameAccessAssociations;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DomainNameAccessAssociations create() =>
      DomainNameAccessAssociations._();
  @$core.override
  DomainNameAccessAssociations createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DomainNameAccessAssociations getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DomainNameAccessAssociations>(create);
  static DomainNameAccessAssociations? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<DomainNameAccessAssociation> get items => $_getList(1);
}

class DomainNames extends $pb.GeneratedMessage {
  factory DomainNames({
    $core.String? position,
    $core.Iterable<DomainName>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  DomainNames._();

  factory DomainNames.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DomainNames.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DomainNames',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<DomainName>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: DomainName.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainNames clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DomainNames copyWith(void Function(DomainNames) updates) =>
      super.copyWith((message) => updates(message as DomainNames))
          as DomainNames;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DomainNames create() => DomainNames._();
  @$core.override
  DomainNames createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DomainNames getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DomainNames>(create);
  static DomainNames? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<DomainName> get items => $_getList(1);
}

class EndpointConfiguration extends $pb.GeneratedMessage {
  factory EndpointConfiguration({
    $core.Iterable<$core.String>? vpcendpointids,
    IpAddressType? ipaddresstype,
    $core.Iterable<EndpointType>? types,
  }) {
    final result = create();
    if (vpcendpointids != null) result.vpcendpointids.addAll(vpcendpointids);
    if (ipaddresstype != null) result.ipaddresstype = ipaddresstype;
    if (types != null) result.types.addAll(types);
    return result;
  }

  EndpointConfiguration._();

  factory EndpointConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EndpointConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EndpointConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(135085374, _omitFieldNames ? '' : 'vpcendpointids')
    ..aE<IpAddressType>(393863813, _omitFieldNames ? '' : 'ipaddresstype',
        enumValues: IpAddressType.values)
    ..pc<EndpointType>(
        534824091, _omitFieldNames ? '' : 'types', $pb.PbFieldType.KE,
        valueOf: EndpointType.valueOf,
        enumValues: EndpointType.values,
        defaultEnumValue: EndpointType.ENDPOINT_TYPE_PRIVATE)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EndpointConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EndpointConfiguration copyWith(
          void Function(EndpointConfiguration) updates) =>
      super.copyWith((message) => updates(message as EndpointConfiguration))
          as EndpointConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EndpointConfiguration create() => EndpointConfiguration._();
  @$core.override
  EndpointConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EndpointConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EndpointConfiguration>(create);
  static EndpointConfiguration? _defaultInstance;

  @$pb.TagNumber(135085374)
  $pb.PbList<$core.String> get vpcendpointids => $_getList(0);

  @$pb.TagNumber(393863813)
  IpAddressType get ipaddresstype => $_getN(1);
  @$pb.TagNumber(393863813)
  set ipaddresstype(IpAddressType value) => $_setField(393863813, value);
  @$pb.TagNumber(393863813)
  $core.bool hasIpaddresstype() => $_has(1);
  @$pb.TagNumber(393863813)
  void clearIpaddresstype() => $_clearField(393863813);

  @$pb.TagNumber(534824091)
  $pb.PbList<EndpointType> get types => $_getList(2);
}

class ExportResponse extends $pb.GeneratedMessage {
  factory ExportResponse({
    $core.String? contenttype,
    $core.String? contentdisposition,
    $core.List<$core.int>? body,
  }) {
    final result = create();
    if (contenttype != null) result.contenttype = contenttype;
    if (contentdisposition != null)
      result.contentdisposition = contentdisposition;
    if (body != null) result.body = body;
    return result;
  }

  ExportResponse._();

  factory ExportResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(281764659, _omitFieldNames ? '' : 'contenttype')
    ..aOS(375146466, _omitFieldNames ? '' : 'contentdisposition')
    ..a<$core.List<$core.int>>(
        464157046, _omitFieldNames ? '' : 'body', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportResponse copyWith(void Function(ExportResponse) updates) =>
      super.copyWith((message) => updates(message as ExportResponse))
          as ExportResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportResponse create() => ExportResponse._();
  @$core.override
  ExportResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportResponse>(create);
  static ExportResponse? _defaultInstance;

  @$pb.TagNumber(281764659)
  $core.String get contenttype => $_getSZ(0);
  @$pb.TagNumber(281764659)
  set contenttype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(281764659)
  $core.bool hasContenttype() => $_has(0);
  @$pb.TagNumber(281764659)
  void clearContenttype() => $_clearField(281764659);

  @$pb.TagNumber(375146466)
  $core.String get contentdisposition => $_getSZ(1);
  @$pb.TagNumber(375146466)
  set contentdisposition($core.String value) => $_setString(1, value);
  @$pb.TagNumber(375146466)
  $core.bool hasContentdisposition() => $_has(1);
  @$pb.TagNumber(375146466)
  void clearContentdisposition() => $_clearField(375146466);

  @$pb.TagNumber(464157046)
  $core.List<$core.int> get body => $_getN(2);
  @$pb.TagNumber(464157046)
  set body($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(2);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);
}

class FlushStageAuthorizersCacheRequest extends $pb.GeneratedMessage {
  factory FlushStageAuthorizersCacheRequest({
    $core.String? stagename,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  FlushStageAuthorizersCacheRequest._();

  factory FlushStageAuthorizersCacheRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FlushStageAuthorizersCacheRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FlushStageAuthorizersCacheRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FlushStageAuthorizersCacheRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FlushStageAuthorizersCacheRequest copyWith(
          void Function(FlushStageAuthorizersCacheRequest) updates) =>
      super.copyWith((message) =>
              updates(message as FlushStageAuthorizersCacheRequest))
          as FlushStageAuthorizersCacheRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FlushStageAuthorizersCacheRequest create() =>
      FlushStageAuthorizersCacheRequest._();
  @$core.override
  FlushStageAuthorizersCacheRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FlushStageAuthorizersCacheRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FlushStageAuthorizersCacheRequest>(
          create);
  static FlushStageAuthorizersCacheRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class FlushStageCacheRequest extends $pb.GeneratedMessage {
  factory FlushStageCacheRequest({
    $core.String? stagename,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  FlushStageCacheRequest._();

  factory FlushStageCacheRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FlushStageCacheRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FlushStageCacheRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FlushStageCacheRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FlushStageCacheRequest copyWith(
          void Function(FlushStageCacheRequest) updates) =>
      super.copyWith((message) => updates(message as FlushStageCacheRequest))
          as FlushStageCacheRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FlushStageCacheRequest create() => FlushStageCacheRequest._();
  @$core.override
  FlushStageCacheRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FlushStageCacheRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FlushStageCacheRequest>(create);
  static FlushStageCacheRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GatewayResponse extends $pb.GeneratedMessage {
  factory GatewayResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        responseparameters,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        responsetemplates,
    $core.bool? defaultresponse,
    $core.String? statuscode,
    GatewayResponseType? responsetype,
  }) {
    final result = create();
    if (responseparameters != null)
      result.responseparameters.addEntries(responseparameters);
    if (responsetemplates != null)
      result.responsetemplates.addEntries(responsetemplates);
    if (defaultresponse != null) result.defaultresponse = defaultresponse;
    if (statuscode != null) result.statuscode = statuscode;
    if (responsetype != null) result.responsetype = responsetype;
    return result;
  }

  GatewayResponse._();

  factory GatewayResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GatewayResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GatewayResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        64271839, _omitFieldNames ? '' : 'responseparameters',
        entryClassName: 'GatewayResponse.ResponseparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..m<$core.String, $core.String>(
        107376570, _omitFieldNames ? '' : 'responsetemplates',
        entryClassName: 'GatewayResponse.ResponsetemplatesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOB(164690914, _omitFieldNames ? '' : 'defaultresponse')
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aE<GatewayResponseType>(377935935, _omitFieldNames ? '' : 'responsetype',
        enumValues: GatewayResponseType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GatewayResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GatewayResponse copyWith(void Function(GatewayResponse) updates) =>
      super.copyWith((message) => updates(message as GatewayResponse))
          as GatewayResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GatewayResponse create() => GatewayResponse._();
  @$core.override
  GatewayResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GatewayResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GatewayResponse>(create);
  static GatewayResponse? _defaultInstance;

  @$pb.TagNumber(64271839)
  $pb.PbMap<$core.String, $core.String> get responseparameters => $_getMap(0);

  @$pb.TagNumber(107376570)
  $pb.PbMap<$core.String, $core.String> get responsetemplates => $_getMap(1);

  @$pb.TagNumber(164690914)
  $core.bool get defaultresponse => $_getBF(2);
  @$pb.TagNumber(164690914)
  set defaultresponse($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(164690914)
  $core.bool hasDefaultresponse() => $_has(2);
  @$pb.TagNumber(164690914)
  void clearDefaultresponse() => $_clearField(164690914);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(3);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(3, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(3);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(377935935)
  GatewayResponseType get responsetype => $_getN(4);
  @$pb.TagNumber(377935935)
  set responsetype(GatewayResponseType value) => $_setField(377935935, value);
  @$pb.TagNumber(377935935)
  $core.bool hasResponsetype() => $_has(4);
  @$pb.TagNumber(377935935)
  void clearResponsetype() => $_clearField(377935935);
}

class GatewayResponses extends $pb.GeneratedMessage {
  factory GatewayResponses({
    $core.String? position,
    $core.Iterable<GatewayResponse>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  GatewayResponses._();

  factory GatewayResponses.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GatewayResponses.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GatewayResponses',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<GatewayResponse>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: GatewayResponse.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GatewayResponses clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GatewayResponses copyWith(void Function(GatewayResponses) updates) =>
      super.copyWith((message) => updates(message as GatewayResponses))
          as GatewayResponses;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GatewayResponses create() => GatewayResponses._();
  @$core.override
  GatewayResponses createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GatewayResponses getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GatewayResponses>(create);
  static GatewayResponses? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<GatewayResponse> get items => $_getList(1);
}

class GenerateClientCertificateRequest extends $pb.GeneratedMessage {
  factory GenerateClientCertificateRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
  }) {
    final result = create();
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    return result;
  }

  GenerateClientCertificateRequest._();

  factory GenerateClientCertificateRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateClientCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateClientCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'GenerateClientCertificateRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateClientCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateClientCertificateRequest copyWith(
          void Function(GenerateClientCertificateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GenerateClientCertificateRequest))
          as GenerateClientCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateClientCertificateRequest create() =>
      GenerateClientCertificateRequest._();
  @$core.override
  GenerateClientCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateClientCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateClientCertificateRequest>(
          create);
  static GenerateClientCertificateRequest? _defaultInstance;

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(0);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);
}

class GetAccountRequest extends $pb.GeneratedMessage {
  factory GetAccountRequest() => create();

  GetAccountRequest._();

  factory GetAccountRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAccountRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAccountRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountRequest copyWith(void Function(GetAccountRequest) updates) =>
      super.copyWith((message) => updates(message as GetAccountRequest))
          as GetAccountRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccountRequest create() => GetAccountRequest._();
  @$core.override
  GetAccountRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAccountRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAccountRequest>(create);
  static GetAccountRequest? _defaultInstance;
}

class GetApiKeyRequest extends $pb.GeneratedMessage {
  factory GetApiKeyRequest({
    $core.bool? includevalue,
    $core.String? apikey,
  }) {
    final result = create();
    if (includevalue != null) result.includevalue = includevalue;
    if (apikey != null) result.apikey = apikey;
    return result;
  }

  GetApiKeyRequest._();

  factory GetApiKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetApiKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetApiKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOB(462860701, _omitFieldNames ? '' : 'includevalue')
    ..aOS(490490655, _omitFieldNames ? '' : 'apikey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetApiKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetApiKeyRequest copyWith(void Function(GetApiKeyRequest) updates) =>
      super.copyWith((message) => updates(message as GetApiKeyRequest))
          as GetApiKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetApiKeyRequest create() => GetApiKeyRequest._();
  @$core.override
  GetApiKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetApiKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetApiKeyRequest>(create);
  static GetApiKeyRequest? _defaultInstance;

  @$pb.TagNumber(462860701)
  $core.bool get includevalue => $_getBF(0);
  @$pb.TagNumber(462860701)
  set includevalue($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(462860701)
  $core.bool hasIncludevalue() => $_has(0);
  @$pb.TagNumber(462860701)
  void clearIncludevalue() => $_clearField(462860701);

  @$pb.TagNumber(490490655)
  $core.String get apikey => $_getSZ(1);
  @$pb.TagNumber(490490655)
  set apikey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(490490655)
  $core.bool hasApikey() => $_has(1);
  @$pb.TagNumber(490490655)
  void clearApikey() => $_clearField(490490655);
}

class GetApiKeysRequest extends $pb.GeneratedMessage {
  factory GetApiKeysRequest({
    $core.String? namequery,
    $core.String? customerid,
    $core.int? limit,
    $core.String? position,
    $core.bool? includevalues,
  }) {
    final result = create();
    if (namequery != null) result.namequery = namequery;
    if (customerid != null) result.customerid = customerid;
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (includevalues != null) result.includevalues = includevalues;
    return result;
  }

  GetApiKeysRequest._();

  factory GetApiKeysRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetApiKeysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetApiKeysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(51018795, _omitFieldNames ? '' : 'namequery')
    ..aOS(227830269, _omitFieldNames ? '' : 'customerid')
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOB(490347326, _omitFieldNames ? '' : 'includevalues')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetApiKeysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetApiKeysRequest copyWith(void Function(GetApiKeysRequest) updates) =>
      super.copyWith((message) => updates(message as GetApiKeysRequest))
          as GetApiKeysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetApiKeysRequest create() => GetApiKeysRequest._();
  @$core.override
  GetApiKeysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetApiKeysRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetApiKeysRequest>(create);
  static GetApiKeysRequest? _defaultInstance;

  @$pb.TagNumber(51018795)
  $core.String get namequery => $_getSZ(0);
  @$pb.TagNumber(51018795)
  set namequery($core.String value) => $_setString(0, value);
  @$pb.TagNumber(51018795)
  $core.bool hasNamequery() => $_has(0);
  @$pb.TagNumber(51018795)
  void clearNamequery() => $_clearField(51018795);

  @$pb.TagNumber(227830269)
  $core.String get customerid => $_getSZ(1);
  @$pb.TagNumber(227830269)
  set customerid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(227830269)
  $core.bool hasCustomerid() => $_has(1);
  @$pb.TagNumber(227830269)
  void clearCustomerid() => $_clearField(227830269);

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(3);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(3, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(3);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(490347326)
  $core.bool get includevalues => $_getBF(4);
  @$pb.TagNumber(490347326)
  set includevalues($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(490347326)
  $core.bool hasIncludevalues() => $_has(4);
  @$pb.TagNumber(490347326)
  void clearIncludevalues() => $_clearField(490347326);
}

class GetAuthorizerRequest extends $pb.GeneratedMessage {
  factory GetAuthorizerRequest({
    $core.String? authorizerid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (authorizerid != null) result.authorizerid = authorizerid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetAuthorizerRequest._();

  factory GetAuthorizerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAuthorizerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAuthorizerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(111773148, _omitFieldNames ? '' : 'authorizerid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAuthorizerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAuthorizerRequest copyWith(void Function(GetAuthorizerRequest) updates) =>
      super.copyWith((message) => updates(message as GetAuthorizerRequest))
          as GetAuthorizerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAuthorizerRequest create() => GetAuthorizerRequest._();
  @$core.override
  GetAuthorizerRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAuthorizerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAuthorizerRequest>(create);
  static GetAuthorizerRequest? _defaultInstance;

  @$pb.TagNumber(111773148)
  $core.String get authorizerid => $_getSZ(0);
  @$pb.TagNumber(111773148)
  set authorizerid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(111773148)
  $core.bool hasAuthorizerid() => $_has(0);
  @$pb.TagNumber(111773148)
  void clearAuthorizerid() => $_clearField(111773148);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetAuthorizersRequest extends $pb.GeneratedMessage {
  factory GetAuthorizersRequest({
    $core.int? limit,
    $core.String? position,
    $core.String? restapiid,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetAuthorizersRequest._();

  factory GetAuthorizersRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAuthorizersRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAuthorizersRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAuthorizersRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAuthorizersRequest copyWith(
          void Function(GetAuthorizersRequest) updates) =>
      super.copyWith((message) => updates(message as GetAuthorizersRequest))
          as GetAuthorizersRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAuthorizersRequest create() => GetAuthorizersRequest._();
  @$core.override
  GetAuthorizersRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAuthorizersRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAuthorizersRequest>(create);
  static GetAuthorizersRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetBasePathMappingRequest extends $pb.GeneratedMessage {
  factory GetBasePathMappingRequest({
    $core.String? basepath,
    $core.String? domainnameid,
    $core.String? domainname,
  }) {
    final result = create();
    if (basepath != null) result.basepath = basepath;
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  GetBasePathMappingRequest._();

  factory GetBasePathMappingRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBasePathMappingRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBasePathMappingRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(267528880, _omitFieldNames ? '' : 'basepath')
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBasePathMappingRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBasePathMappingRequest copyWith(
          void Function(GetBasePathMappingRequest) updates) =>
      super.copyWith((message) => updates(message as GetBasePathMappingRequest))
          as GetBasePathMappingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBasePathMappingRequest create() => GetBasePathMappingRequest._();
  @$core.override
  GetBasePathMappingRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBasePathMappingRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBasePathMappingRequest>(create);
  static GetBasePathMappingRequest? _defaultInstance;

  @$pb.TagNumber(267528880)
  $core.String get basepath => $_getSZ(0);
  @$pb.TagNumber(267528880)
  set basepath($core.String value) => $_setString(0, value);
  @$pb.TagNumber(267528880)
  $core.bool hasBasepath() => $_has(0);
  @$pb.TagNumber(267528880)
  void clearBasepath() => $_clearField(267528880);

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(1);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(1);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(2);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(2);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);
}

class GetBasePathMappingsRequest extends $pb.GeneratedMessage {
  factory GetBasePathMappingsRequest({
    $core.String? domainnameid,
    $core.int? limit,
    $core.String? position,
    $core.String? domainname,
  }) {
    final result = create();
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  GetBasePathMappingsRequest._();

  factory GetBasePathMappingsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBasePathMappingsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBasePathMappingsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBasePathMappingsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBasePathMappingsRequest copyWith(
          void Function(GetBasePathMappingsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetBasePathMappingsRequest))
          as GetBasePathMappingsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBasePathMappingsRequest create() => GetBasePathMappingsRequest._();
  @$core.override
  GetBasePathMappingsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBasePathMappingsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBasePathMappingsRequest>(create);
  static GetBasePathMappingsRequest? _defaultInstance;

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(0);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(0);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(2);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(2, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(2);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(3);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(3);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);
}

class GetClientCertificateRequest extends $pb.GeneratedMessage {
  factory GetClientCertificateRequest({
    $core.String? clientcertificateid,
  }) {
    final result = create();
    if (clientcertificateid != null)
      result.clientcertificateid = clientcertificateid;
    return result;
  }

  GetClientCertificateRequest._();

  factory GetClientCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetClientCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetClientCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(276222909, _omitFieldNames ? '' : 'clientcertificateid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetClientCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetClientCertificateRequest copyWith(
          void Function(GetClientCertificateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetClientCertificateRequest))
          as GetClientCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetClientCertificateRequest create() =>
      GetClientCertificateRequest._();
  @$core.override
  GetClientCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetClientCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetClientCertificateRequest>(create);
  static GetClientCertificateRequest? _defaultInstance;

  @$pb.TagNumber(276222909)
  $core.String get clientcertificateid => $_getSZ(0);
  @$pb.TagNumber(276222909)
  set clientcertificateid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(276222909)
  $core.bool hasClientcertificateid() => $_has(0);
  @$pb.TagNumber(276222909)
  void clearClientcertificateid() => $_clearField(276222909);
}

class GetClientCertificatesRequest extends $pb.GeneratedMessage {
  factory GetClientCertificatesRequest({
    $core.int? limit,
    $core.String? position,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    return result;
  }

  GetClientCertificatesRequest._();

  factory GetClientCertificatesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetClientCertificatesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetClientCertificatesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetClientCertificatesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetClientCertificatesRequest copyWith(
          void Function(GetClientCertificatesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetClientCertificatesRequest))
          as GetClientCertificatesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetClientCertificatesRequest create() =>
      GetClientCertificatesRequest._();
  @$core.override
  GetClientCertificatesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetClientCertificatesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetClientCertificatesRequest>(create);
  static GetClientCertificatesRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);
}

class GetDeploymentRequest extends $pb.GeneratedMessage {
  factory GetDeploymentRequest({
    $core.Iterable<$core.String>? embed,
    $core.String? restapiid,
    $core.String? deploymentid,
  }) {
    final result = create();
    if (embed != null) result.embed.addAll(embed);
    if (restapiid != null) result.restapiid = restapiid;
    if (deploymentid != null) result.deploymentid = deploymentid;
    return result;
  }

  GetDeploymentRequest._();

  factory GetDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(136029775, _omitFieldNames ? '' : 'embed')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(439369188, _omitFieldNames ? '' : 'deploymentid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDeploymentRequest copyWith(void Function(GetDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as GetDeploymentRequest))
          as GetDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDeploymentRequest create() => GetDeploymentRequest._();
  @$core.override
  GetDeploymentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDeploymentRequest>(create);
  static GetDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(136029775)
  $pb.PbList<$core.String> get embed => $_getList(0);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(439369188)
  $core.String get deploymentid => $_getSZ(2);
  @$pb.TagNumber(439369188)
  set deploymentid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(439369188)
  $core.bool hasDeploymentid() => $_has(2);
  @$pb.TagNumber(439369188)
  void clearDeploymentid() => $_clearField(439369188);
}

class GetDeploymentsRequest extends $pb.GeneratedMessage {
  factory GetDeploymentsRequest({
    $core.int? limit,
    $core.String? position,
    $core.String? restapiid,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetDeploymentsRequest._();

  factory GetDeploymentsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDeploymentsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDeploymentsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDeploymentsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDeploymentsRequest copyWith(
          void Function(GetDeploymentsRequest) updates) =>
      super.copyWith((message) => updates(message as GetDeploymentsRequest))
          as GetDeploymentsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDeploymentsRequest create() => GetDeploymentsRequest._();
  @$core.override
  GetDeploymentsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDeploymentsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDeploymentsRequest>(create);
  static GetDeploymentsRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetDocumentationPartRequest extends $pb.GeneratedMessage {
  factory GetDocumentationPartRequest({
    $core.String? documentationpartid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (documentationpartid != null)
      result.documentationpartid = documentationpartid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetDocumentationPartRequest._();

  factory GetDocumentationPartRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDocumentationPartRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDocumentationPartRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(286552774, _omitFieldNames ? '' : 'documentationpartid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDocumentationPartRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDocumentationPartRequest copyWith(
          void Function(GetDocumentationPartRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetDocumentationPartRequest))
          as GetDocumentationPartRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDocumentationPartRequest create() =>
      GetDocumentationPartRequest._();
  @$core.override
  GetDocumentationPartRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDocumentationPartRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDocumentationPartRequest>(create);
  static GetDocumentationPartRequest? _defaultInstance;

  @$pb.TagNumber(286552774)
  $core.String get documentationpartid => $_getSZ(0);
  @$pb.TagNumber(286552774)
  set documentationpartid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(286552774)
  $core.bool hasDocumentationpartid() => $_has(0);
  @$pb.TagNumber(286552774)
  void clearDocumentationpartid() => $_clearField(286552774);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetDocumentationPartsRequest extends $pb.GeneratedMessage {
  factory GetDocumentationPartsRequest({
    $core.String? namequery,
    $core.String? path,
    DocumentationPartType? type,
    $core.int? limit,
    $core.String? position,
    $core.String? restapiid,
    LocationStatusType? locationstatus,
  }) {
    final result = create();
    if (namequery != null) result.namequery = namequery;
    if (path != null) result.path = path;
    if (type != null) result.type = type;
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (restapiid != null) result.restapiid = restapiid;
    if (locationstatus != null) result.locationstatus = locationstatus;
    return result;
  }

  GetDocumentationPartsRequest._();

  factory GetDocumentationPartsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDocumentationPartsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDocumentationPartsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(51018795, _omitFieldNames ? '' : 'namequery')
    ..aOS(75975991, _omitFieldNames ? '' : 'path')
    ..aE<DocumentationPartType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: DocumentationPartType.values)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aE<LocationStatusType>(532215305, _omitFieldNames ? '' : 'locationstatus',
        enumValues: LocationStatusType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDocumentationPartsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDocumentationPartsRequest copyWith(
          void Function(GetDocumentationPartsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetDocumentationPartsRequest))
          as GetDocumentationPartsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDocumentationPartsRequest create() =>
      GetDocumentationPartsRequest._();
  @$core.override
  GetDocumentationPartsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDocumentationPartsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDocumentationPartsRequest>(create);
  static GetDocumentationPartsRequest? _defaultInstance;

  @$pb.TagNumber(51018795)
  $core.String get namequery => $_getSZ(0);
  @$pb.TagNumber(51018795)
  set namequery($core.String value) => $_setString(0, value);
  @$pb.TagNumber(51018795)
  $core.bool hasNamequery() => $_has(0);
  @$pb.TagNumber(51018795)
  void clearNamequery() => $_clearField(51018795);

  @$pb.TagNumber(75975991)
  $core.String get path => $_getSZ(1);
  @$pb.TagNumber(75975991)
  set path($core.String value) => $_setString(1, value);
  @$pb.TagNumber(75975991)
  $core.bool hasPath() => $_has(1);
  @$pb.TagNumber(75975991)
  void clearPath() => $_clearField(75975991);

  @$pb.TagNumber(287830350)
  DocumentationPartType get type => $_getN(2);
  @$pb.TagNumber(287830350)
  set type(DocumentationPartType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(4);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(4, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(4);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(5);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(5);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(532215305)
  LocationStatusType get locationstatus => $_getN(6);
  @$pb.TagNumber(532215305)
  set locationstatus(LocationStatusType value) => $_setField(532215305, value);
  @$pb.TagNumber(532215305)
  $core.bool hasLocationstatus() => $_has(6);
  @$pb.TagNumber(532215305)
  void clearLocationstatus() => $_clearField(532215305);
}

class GetDocumentationVersionRequest extends $pb.GeneratedMessage {
  factory GetDocumentationVersionRequest({
    $core.String? documentationversion,
    $core.String? restapiid,
  }) {
    final result = create();
    if (documentationversion != null)
      result.documentationversion = documentationversion;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetDocumentationVersionRequest._();

  factory GetDocumentationVersionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDocumentationVersionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDocumentationVersionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(167009804, _omitFieldNames ? '' : 'documentationversion')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDocumentationVersionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDocumentationVersionRequest copyWith(
          void Function(GetDocumentationVersionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetDocumentationVersionRequest))
          as GetDocumentationVersionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDocumentationVersionRequest create() =>
      GetDocumentationVersionRequest._();
  @$core.override
  GetDocumentationVersionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDocumentationVersionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDocumentationVersionRequest>(create);
  static GetDocumentationVersionRequest? _defaultInstance;

  @$pb.TagNumber(167009804)
  $core.String get documentationversion => $_getSZ(0);
  @$pb.TagNumber(167009804)
  set documentationversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(167009804)
  $core.bool hasDocumentationversion() => $_has(0);
  @$pb.TagNumber(167009804)
  void clearDocumentationversion() => $_clearField(167009804);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetDocumentationVersionsRequest extends $pb.GeneratedMessage {
  factory GetDocumentationVersionsRequest({
    $core.int? limit,
    $core.String? position,
    $core.String? restapiid,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetDocumentationVersionsRequest._();

  factory GetDocumentationVersionsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDocumentationVersionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDocumentationVersionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDocumentationVersionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDocumentationVersionsRequest copyWith(
          void Function(GetDocumentationVersionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetDocumentationVersionsRequest))
          as GetDocumentationVersionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDocumentationVersionsRequest create() =>
      GetDocumentationVersionsRequest._();
  @$core.override
  GetDocumentationVersionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDocumentationVersionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDocumentationVersionsRequest>(
          create);
  static GetDocumentationVersionsRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetDomainNameAccessAssociationsRequest extends $pb.GeneratedMessage {
  factory GetDomainNameAccessAssociationsRequest({
    ResourceOwner? resourceowner,
    $core.int? limit,
    $core.String? position,
  }) {
    final result = create();
    if (resourceowner != null) result.resourceowner = resourceowner;
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    return result;
  }

  GetDomainNameAccessAssociationsRequest._();

  factory GetDomainNameAccessAssociationsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDomainNameAccessAssociationsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDomainNameAccessAssociationsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aE<ResourceOwner>(259175301, _omitFieldNames ? '' : 'resourceowner',
        enumValues: ResourceOwner.values)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDomainNameAccessAssociationsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDomainNameAccessAssociationsRequest copyWith(
          void Function(GetDomainNameAccessAssociationsRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetDomainNameAccessAssociationsRequest))
          as GetDomainNameAccessAssociationsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDomainNameAccessAssociationsRequest create() =>
      GetDomainNameAccessAssociationsRequest._();
  @$core.override
  GetDomainNameAccessAssociationsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDomainNameAccessAssociationsRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetDomainNameAccessAssociationsRequest>(create);
  static GetDomainNameAccessAssociationsRequest? _defaultInstance;

  @$pb.TagNumber(259175301)
  ResourceOwner get resourceowner => $_getN(0);
  @$pb.TagNumber(259175301)
  set resourceowner(ResourceOwner value) => $_setField(259175301, value);
  @$pb.TagNumber(259175301)
  $core.bool hasResourceowner() => $_has(0);
  @$pb.TagNumber(259175301)
  void clearResourceowner() => $_clearField(259175301);

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(2);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(2, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(2);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);
}

class GetDomainNameRequest extends $pb.GeneratedMessage {
  factory GetDomainNameRequest({
    $core.String? domainnameid,
    $core.String? domainname,
  }) {
    final result = create();
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  GetDomainNameRequest._();

  factory GetDomainNameRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDomainNameRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDomainNameRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDomainNameRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDomainNameRequest copyWith(void Function(GetDomainNameRequest) updates) =>
      super.copyWith((message) => updates(message as GetDomainNameRequest))
          as GetDomainNameRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDomainNameRequest create() => GetDomainNameRequest._();
  @$core.override
  GetDomainNameRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDomainNameRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDomainNameRequest>(create);
  static GetDomainNameRequest? _defaultInstance;

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(0);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(0);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(1);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(1);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);
}

class GetDomainNamesRequest extends $pb.GeneratedMessage {
  factory GetDomainNamesRequest({
    ResourceOwner? resourceowner,
    $core.int? limit,
    $core.String? position,
  }) {
    final result = create();
    if (resourceowner != null) result.resourceowner = resourceowner;
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    return result;
  }

  GetDomainNamesRequest._();

  factory GetDomainNamesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDomainNamesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDomainNamesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aE<ResourceOwner>(259175301, _omitFieldNames ? '' : 'resourceowner',
        enumValues: ResourceOwner.values)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDomainNamesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDomainNamesRequest copyWith(
          void Function(GetDomainNamesRequest) updates) =>
      super.copyWith((message) => updates(message as GetDomainNamesRequest))
          as GetDomainNamesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDomainNamesRequest create() => GetDomainNamesRequest._();
  @$core.override
  GetDomainNamesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDomainNamesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDomainNamesRequest>(create);
  static GetDomainNamesRequest? _defaultInstance;

  @$pb.TagNumber(259175301)
  ResourceOwner get resourceowner => $_getN(0);
  @$pb.TagNumber(259175301)
  set resourceowner(ResourceOwner value) => $_setField(259175301, value);
  @$pb.TagNumber(259175301)
  $core.bool hasResourceowner() => $_has(0);
  @$pb.TagNumber(259175301)
  void clearResourceowner() => $_clearField(259175301);

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(2);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(2, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(2);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);
}

class GetExportRequest extends $pb.GeneratedMessage {
  factory GetExportRequest({
    $core.String? stagename,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
    $core.String? accepts,
    $core.String? exporttype,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (parameters != null) result.parameters.addEntries(parameters);
    if (accepts != null) result.accepts = accepts;
    if (exporttype != null) result.exporttype = exporttype;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetExportRequest._();

  factory GetExportRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetExportRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetExportRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..m<$core.String, $core.String>(
        145043162, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'GetExportRequest.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(192079791, _omitFieldNames ? '' : 'accepts')
    ..aOS(243495788, _omitFieldNames ? '' : 'exporttype')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetExportRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetExportRequest copyWith(void Function(GetExportRequest) updates) =>
      super.copyWith((message) => updates(message as GetExportRequest))
          as GetExportRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetExportRequest create() => GetExportRequest._();
  @$core.override
  GetExportRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetExportRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetExportRequest>(create);
  static GetExportRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(145043162)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(1);

  @$pb.TagNumber(192079791)
  $core.String get accepts => $_getSZ(2);
  @$pb.TagNumber(192079791)
  set accepts($core.String value) => $_setString(2, value);
  @$pb.TagNumber(192079791)
  $core.bool hasAccepts() => $_has(2);
  @$pb.TagNumber(192079791)
  void clearAccepts() => $_clearField(192079791);

  @$pb.TagNumber(243495788)
  $core.String get exporttype => $_getSZ(3);
  @$pb.TagNumber(243495788)
  set exporttype($core.String value) => $_setString(3, value);
  @$pb.TagNumber(243495788)
  $core.bool hasExporttype() => $_has(3);
  @$pb.TagNumber(243495788)
  void clearExporttype() => $_clearField(243495788);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(4);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(4);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetGatewayResponseRequest extends $pb.GeneratedMessage {
  factory GetGatewayResponseRequest({
    GatewayResponseType? responsetype,
    $core.String? restapiid,
  }) {
    final result = create();
    if (responsetype != null) result.responsetype = responsetype;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetGatewayResponseRequest._();

  factory GetGatewayResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetGatewayResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetGatewayResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aE<GatewayResponseType>(377935935, _omitFieldNames ? '' : 'responsetype',
        enumValues: GatewayResponseType.values)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGatewayResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGatewayResponseRequest copyWith(
          void Function(GetGatewayResponseRequest) updates) =>
      super.copyWith((message) => updates(message as GetGatewayResponseRequest))
          as GetGatewayResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetGatewayResponseRequest create() => GetGatewayResponseRequest._();
  @$core.override
  GetGatewayResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetGatewayResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetGatewayResponseRequest>(create);
  static GetGatewayResponseRequest? _defaultInstance;

  @$pb.TagNumber(377935935)
  GatewayResponseType get responsetype => $_getN(0);
  @$pb.TagNumber(377935935)
  set responsetype(GatewayResponseType value) => $_setField(377935935, value);
  @$pb.TagNumber(377935935)
  $core.bool hasResponsetype() => $_has(0);
  @$pb.TagNumber(377935935)
  void clearResponsetype() => $_clearField(377935935);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetGatewayResponsesRequest extends $pb.GeneratedMessage {
  factory GetGatewayResponsesRequest({
    $core.int? limit,
    $core.String? position,
    $core.String? restapiid,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetGatewayResponsesRequest._();

  factory GetGatewayResponsesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetGatewayResponsesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetGatewayResponsesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGatewayResponsesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGatewayResponsesRequest copyWith(
          void Function(GetGatewayResponsesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetGatewayResponsesRequest))
          as GetGatewayResponsesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetGatewayResponsesRequest create() => GetGatewayResponsesRequest._();
  @$core.override
  GetGatewayResponsesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetGatewayResponsesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetGatewayResponsesRequest>(create);
  static GetGatewayResponsesRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetIntegrationRequest extends $pb.GeneratedMessage {
  factory GetIntegrationRequest({
    $core.String? httpmethod,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetIntegrationRequest._();

  factory GetIntegrationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetIntegrationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetIntegrationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIntegrationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIntegrationRequest copyWith(
          void Function(GetIntegrationRequest) updates) =>
      super.copyWith((message) => updates(message as GetIntegrationRequest))
          as GetIntegrationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetIntegrationRequest create() => GetIntegrationRequest._();
  @$core.override
  GetIntegrationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetIntegrationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetIntegrationRequest>(create);
  static GetIntegrationRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetIntegrationResponseRequest extends $pb.GeneratedMessage {
  factory GetIntegrationResponseRequest({
    $core.String? httpmethod,
    $core.String? statuscode,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (statuscode != null) result.statuscode = statuscode;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetIntegrationResponseRequest._();

  factory GetIntegrationResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetIntegrationResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetIntegrationResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIntegrationResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIntegrationResponseRequest copyWith(
          void Function(GetIntegrationResponseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetIntegrationResponseRequest))
          as GetIntegrationResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetIntegrationResponseRequest create() =>
      GetIntegrationResponseRequest._();
  @$core.override
  GetIntegrationResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetIntegrationResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetIntegrationResponseRequest>(create);
  static GetIntegrationResponseRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(1);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(1);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(2);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(2);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetMethodRequest extends $pb.GeneratedMessage {
  factory GetMethodRequest({
    $core.String? httpmethod,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetMethodRequest._();

  factory GetMethodRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMethodRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMethodRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMethodRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMethodRequest copyWith(void Function(GetMethodRequest) updates) =>
      super.copyWith((message) => updates(message as GetMethodRequest))
          as GetMethodRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMethodRequest create() => GetMethodRequest._();
  @$core.override
  GetMethodRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMethodRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMethodRequest>(create);
  static GetMethodRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetMethodResponseRequest extends $pb.GeneratedMessage {
  factory GetMethodResponseRequest({
    $core.String? httpmethod,
    $core.String? statuscode,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (statuscode != null) result.statuscode = statuscode;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetMethodResponseRequest._();

  factory GetMethodResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMethodResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMethodResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMethodResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMethodResponseRequest copyWith(
          void Function(GetMethodResponseRequest) updates) =>
      super.copyWith((message) => updates(message as GetMethodResponseRequest))
          as GetMethodResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMethodResponseRequest create() => GetMethodResponseRequest._();
  @$core.override
  GetMethodResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMethodResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMethodResponseRequest>(create);
  static GetMethodResponseRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(1);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(1);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(2);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(2);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetModelRequest extends $pb.GeneratedMessage {
  factory GetModelRequest({
    $core.bool? flatten,
    $core.String? modelname,
    $core.String? restapiid,
  }) {
    final result = create();
    if (flatten != null) result.flatten = flatten;
    if (modelname != null) result.modelname = modelname;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetModelRequest._();

  factory GetModelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetModelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetModelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOB(100024266, _omitFieldNames ? '' : 'flatten')
    ..aOS(176835330, _omitFieldNames ? '' : 'modelname')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetModelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetModelRequest copyWith(void Function(GetModelRequest) updates) =>
      super.copyWith((message) => updates(message as GetModelRequest))
          as GetModelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetModelRequest create() => GetModelRequest._();
  @$core.override
  GetModelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetModelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetModelRequest>(create);
  static GetModelRequest? _defaultInstance;

  @$pb.TagNumber(100024266)
  $core.bool get flatten => $_getBF(0);
  @$pb.TagNumber(100024266)
  set flatten($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(100024266)
  $core.bool hasFlatten() => $_has(0);
  @$pb.TagNumber(100024266)
  void clearFlatten() => $_clearField(100024266);

  @$pb.TagNumber(176835330)
  $core.String get modelname => $_getSZ(1);
  @$pb.TagNumber(176835330)
  set modelname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(176835330)
  $core.bool hasModelname() => $_has(1);
  @$pb.TagNumber(176835330)
  void clearModelname() => $_clearField(176835330);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetModelTemplateRequest extends $pb.GeneratedMessage {
  factory GetModelTemplateRequest({
    $core.String? modelname,
    $core.String? restapiid,
  }) {
    final result = create();
    if (modelname != null) result.modelname = modelname;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetModelTemplateRequest._();

  factory GetModelTemplateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetModelTemplateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetModelTemplateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(176835330, _omitFieldNames ? '' : 'modelname')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetModelTemplateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetModelTemplateRequest copyWith(
          void Function(GetModelTemplateRequest) updates) =>
      super.copyWith((message) => updates(message as GetModelTemplateRequest))
          as GetModelTemplateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetModelTemplateRequest create() => GetModelTemplateRequest._();
  @$core.override
  GetModelTemplateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetModelTemplateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetModelTemplateRequest>(create);
  static GetModelTemplateRequest? _defaultInstance;

  @$pb.TagNumber(176835330)
  $core.String get modelname => $_getSZ(0);
  @$pb.TagNumber(176835330)
  set modelname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(176835330)
  $core.bool hasModelname() => $_has(0);
  @$pb.TagNumber(176835330)
  void clearModelname() => $_clearField(176835330);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetModelsRequest extends $pb.GeneratedMessage {
  factory GetModelsRequest({
    $core.int? limit,
    $core.String? position,
    $core.String? restapiid,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetModelsRequest._();

  factory GetModelsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetModelsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetModelsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetModelsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetModelsRequest copyWith(void Function(GetModelsRequest) updates) =>
      super.copyWith((message) => updates(message as GetModelsRequest))
          as GetModelsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetModelsRequest create() => GetModelsRequest._();
  @$core.override
  GetModelsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetModelsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetModelsRequest>(create);
  static GetModelsRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetRequestValidatorRequest extends $pb.GeneratedMessage {
  factory GetRequestValidatorRequest({
    $core.String? restapiid,
    $core.String? requestvalidatorid,
  }) {
    final result = create();
    if (restapiid != null) result.restapiid = restapiid;
    if (requestvalidatorid != null)
      result.requestvalidatorid = requestvalidatorid;
    return result;
  }

  GetRequestValidatorRequest._();

  factory GetRequestValidatorRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRequestValidatorRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRequestValidatorRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(517546134, _omitFieldNames ? '' : 'requestvalidatorid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRequestValidatorRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRequestValidatorRequest copyWith(
          void Function(GetRequestValidatorRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetRequestValidatorRequest))
          as GetRequestValidatorRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRequestValidatorRequest create() => GetRequestValidatorRequest._();
  @$core.override
  GetRequestValidatorRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRequestValidatorRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRequestValidatorRequest>(create);
  static GetRequestValidatorRequest? _defaultInstance;

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(0);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(0);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(517546134)
  $core.String get requestvalidatorid => $_getSZ(1);
  @$pb.TagNumber(517546134)
  set requestvalidatorid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(517546134)
  $core.bool hasRequestvalidatorid() => $_has(1);
  @$pb.TagNumber(517546134)
  void clearRequestvalidatorid() => $_clearField(517546134);
}

class GetRequestValidatorsRequest extends $pb.GeneratedMessage {
  factory GetRequestValidatorsRequest({
    $core.int? limit,
    $core.String? position,
    $core.String? restapiid,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetRequestValidatorsRequest._();

  factory GetRequestValidatorsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRequestValidatorsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRequestValidatorsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRequestValidatorsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRequestValidatorsRequest copyWith(
          void Function(GetRequestValidatorsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetRequestValidatorsRequest))
          as GetRequestValidatorsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRequestValidatorsRequest create() =>
      GetRequestValidatorsRequest._();
  @$core.override
  GetRequestValidatorsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRequestValidatorsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRequestValidatorsRequest>(create);
  static GetRequestValidatorsRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetResourceRequest extends $pb.GeneratedMessage {
  factory GetResourceRequest({
    $core.Iterable<$core.String>? embed,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (embed != null) result.embed.addAll(embed);
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetResourceRequest._();

  factory GetResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(136029775, _omitFieldNames ? '' : 'embed')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourceRequest copyWith(void Function(GetResourceRequest) updates) =>
      super.copyWith((message) => updates(message as GetResourceRequest))
          as GetResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourceRequest create() => GetResourceRequest._();
  @$core.override
  GetResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourceRequest>(create);
  static GetResourceRequest? _defaultInstance;

  @$pb.TagNumber(136029775)
  $pb.PbList<$core.String> get embed => $_getList(0);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetResourcesRequest extends $pb.GeneratedMessage {
  factory GetResourcesRequest({
    $core.Iterable<$core.String>? embed,
    $core.int? limit,
    $core.String? position,
    $core.String? restapiid,
  }) {
    final result = create();
    if (embed != null) result.embed.addAll(embed);
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetResourcesRequest._();

  factory GetResourcesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourcesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourcesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(136029775, _omitFieldNames ? '' : 'embed')
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcesRequest copyWith(void Function(GetResourcesRequest) updates) =>
      super.copyWith((message) => updates(message as GetResourcesRequest))
          as GetResourcesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourcesRequest create() => GetResourcesRequest._();
  @$core.override
  GetResourcesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourcesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourcesRequest>(create);
  static GetResourcesRequest? _defaultInstance;

  @$pb.TagNumber(136029775)
  $pb.PbList<$core.String> get embed => $_getList(0);

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(2);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(2, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(2);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetRestApiRequest extends $pb.GeneratedMessage {
  factory GetRestApiRequest({
    $core.String? restapiid,
  }) {
    final result = create();
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetRestApiRequest._();

  factory GetRestApiRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRestApiRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRestApiRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRestApiRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRestApiRequest copyWith(void Function(GetRestApiRequest) updates) =>
      super.copyWith((message) => updates(message as GetRestApiRequest))
          as GetRestApiRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRestApiRequest create() => GetRestApiRequest._();
  @$core.override
  GetRestApiRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRestApiRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRestApiRequest>(create);
  static GetRestApiRequest? _defaultInstance;

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(0);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(0);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetRestApisRequest extends $pb.GeneratedMessage {
  factory GetRestApisRequest({
    $core.int? limit,
    $core.String? position,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    return result;
  }

  GetRestApisRequest._();

  factory GetRestApisRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRestApisRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRestApisRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRestApisRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRestApisRequest copyWith(void Function(GetRestApisRequest) updates) =>
      super.copyWith((message) => updates(message as GetRestApisRequest))
          as GetRestApisRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRestApisRequest create() => GetRestApisRequest._();
  @$core.override
  GetRestApisRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRestApisRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRestApisRequest>(create);
  static GetRestApisRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);
}

class GetSdkRequest extends $pb.GeneratedMessage {
  factory GetSdkRequest({
    $core.String? stagename,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
    $core.String? sdktype,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (parameters != null) result.parameters.addEntries(parameters);
    if (sdktype != null) result.sdktype = sdktype;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetSdkRequest._();

  factory GetSdkRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSdkRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSdkRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..m<$core.String, $core.String>(
        145043162, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'GetSdkRequest.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(146213956, _omitFieldNames ? '' : 'sdktype')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSdkRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSdkRequest copyWith(void Function(GetSdkRequest) updates) =>
      super.copyWith((message) => updates(message as GetSdkRequest))
          as GetSdkRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSdkRequest create() => GetSdkRequest._();
  @$core.override
  GetSdkRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSdkRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSdkRequest>(create);
  static GetSdkRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(145043162)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(1);

  @$pb.TagNumber(146213956)
  $core.String get sdktype => $_getSZ(2);
  @$pb.TagNumber(146213956)
  set sdktype($core.String value) => $_setString(2, value);
  @$pb.TagNumber(146213956)
  $core.bool hasSdktype() => $_has(2);
  @$pb.TagNumber(146213956)
  void clearSdktype() => $_clearField(146213956);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetSdkTypeRequest extends $pb.GeneratedMessage {
  factory GetSdkTypeRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  GetSdkTypeRequest._();

  factory GetSdkTypeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSdkTypeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSdkTypeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSdkTypeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSdkTypeRequest copyWith(void Function(GetSdkTypeRequest) updates) =>
      super.copyWith((message) => updates(message as GetSdkTypeRequest))
          as GetSdkTypeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSdkTypeRequest create() => GetSdkTypeRequest._();
  @$core.override
  GetSdkTypeRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSdkTypeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSdkTypeRequest>(create);
  static GetSdkTypeRequest? _defaultInstance;

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);
}

class GetSdkTypesRequest extends $pb.GeneratedMessage {
  factory GetSdkTypesRequest({
    $core.int? limit,
    $core.String? position,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    return result;
  }

  GetSdkTypesRequest._();

  factory GetSdkTypesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSdkTypesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSdkTypesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSdkTypesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSdkTypesRequest copyWith(void Function(GetSdkTypesRequest) updates) =>
      super.copyWith((message) => updates(message as GetSdkTypesRequest))
          as GetSdkTypesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSdkTypesRequest create() => GetSdkTypesRequest._();
  @$core.override
  GetSdkTypesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSdkTypesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSdkTypesRequest>(create);
  static GetSdkTypesRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);
}

class GetStageRequest extends $pb.GeneratedMessage {
  factory GetStageRequest({
    $core.String? stagename,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  GetStageRequest._();

  factory GetStageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetStageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetStageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetStageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetStageRequest copyWith(void Function(GetStageRequest) updates) =>
      super.copyWith((message) => updates(message as GetStageRequest))
          as GetStageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetStageRequest create() => GetStageRequest._();
  @$core.override
  GetStageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetStageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetStageRequest>(create);
  static GetStageRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class GetStagesRequest extends $pb.GeneratedMessage {
  factory GetStagesRequest({
    $core.String? restapiid,
    $core.String? deploymentid,
  }) {
    final result = create();
    if (restapiid != null) result.restapiid = restapiid;
    if (deploymentid != null) result.deploymentid = deploymentid;
    return result;
  }

  GetStagesRequest._();

  factory GetStagesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetStagesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetStagesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(439369188, _omitFieldNames ? '' : 'deploymentid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetStagesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetStagesRequest copyWith(void Function(GetStagesRequest) updates) =>
      super.copyWith((message) => updates(message as GetStagesRequest))
          as GetStagesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetStagesRequest create() => GetStagesRequest._();
  @$core.override
  GetStagesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetStagesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetStagesRequest>(create);
  static GetStagesRequest? _defaultInstance;

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(0);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(0);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(439369188)
  $core.String get deploymentid => $_getSZ(1);
  @$pb.TagNumber(439369188)
  set deploymentid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(439369188)
  $core.bool hasDeploymentid() => $_has(1);
  @$pb.TagNumber(439369188)
  void clearDeploymentid() => $_clearField(439369188);
}

class GetTagsRequest extends $pb.GeneratedMessage {
  factory GetTagsRequest({
    $core.String? resourcearn,
    $core.int? limit,
    $core.String? position,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    return result;
  }

  GetTagsRequest._();

  factory GetTagsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTagsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTagsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(67806797, _omitFieldNames ? '' : 'resourcearn')
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTagsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTagsRequest copyWith(void Function(GetTagsRequest) updates) =>
      super.copyWith((message) => updates(message as GetTagsRequest))
          as GetTagsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTagsRequest create() => GetTagsRequest._();
  @$core.override
  GetTagsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTagsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTagsRequest>(create);
  static GetTagsRequest? _defaultInstance;

  @$pb.TagNumber(67806797)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(67806797)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(67806797)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(67806797)
  void clearResourcearn() => $_clearField(67806797);

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(2);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(2, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(2);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);
}

class GetUsagePlanKeyRequest extends $pb.GeneratedMessage {
  factory GetUsagePlanKeyRequest({
    $core.String? keyid,
    $core.String? usageplanid,
  }) {
    final result = create();
    if (keyid != null) result.keyid = keyid;
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  GetUsagePlanKeyRequest._();

  factory GetUsagePlanKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetUsagePlanKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetUsagePlanKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(479913282, _omitFieldNames ? '' : 'keyid')
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsagePlanKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsagePlanKeyRequest copyWith(
          void Function(GetUsagePlanKeyRequest) updates) =>
      super.copyWith((message) => updates(message as GetUsagePlanKeyRequest))
          as GetUsagePlanKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetUsagePlanKeyRequest create() => GetUsagePlanKeyRequest._();
  @$core.override
  GetUsagePlanKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetUsagePlanKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetUsagePlanKeyRequest>(create);
  static GetUsagePlanKeyRequest? _defaultInstance;

  @$pb.TagNumber(479913282)
  $core.String get keyid => $_getSZ(0);
  @$pb.TagNumber(479913282)
  set keyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(479913282)
  $core.bool hasKeyid() => $_has(0);
  @$pb.TagNumber(479913282)
  void clearKeyid() => $_clearField(479913282);

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(1);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(1);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class GetUsagePlanKeysRequest extends $pb.GeneratedMessage {
  factory GetUsagePlanKeysRequest({
    $core.String? namequery,
    $core.int? limit,
    $core.String? position,
    $core.String? usageplanid,
  }) {
    final result = create();
    if (namequery != null) result.namequery = namequery;
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  GetUsagePlanKeysRequest._();

  factory GetUsagePlanKeysRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetUsagePlanKeysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetUsagePlanKeysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(51018795, _omitFieldNames ? '' : 'namequery')
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsagePlanKeysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsagePlanKeysRequest copyWith(
          void Function(GetUsagePlanKeysRequest) updates) =>
      super.copyWith((message) => updates(message as GetUsagePlanKeysRequest))
          as GetUsagePlanKeysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetUsagePlanKeysRequest create() => GetUsagePlanKeysRequest._();
  @$core.override
  GetUsagePlanKeysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetUsagePlanKeysRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetUsagePlanKeysRequest>(create);
  static GetUsagePlanKeysRequest? _defaultInstance;

  @$pb.TagNumber(51018795)
  $core.String get namequery => $_getSZ(0);
  @$pb.TagNumber(51018795)
  set namequery($core.String value) => $_setString(0, value);
  @$pb.TagNumber(51018795)
  $core.bool hasNamequery() => $_has(0);
  @$pb.TagNumber(51018795)
  void clearNamequery() => $_clearField(51018795);

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(2);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(2, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(2);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(3);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(3);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class GetUsagePlanRequest extends $pb.GeneratedMessage {
  factory GetUsagePlanRequest({
    $core.String? usageplanid,
  }) {
    final result = create();
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  GetUsagePlanRequest._();

  factory GetUsagePlanRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetUsagePlanRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetUsagePlanRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsagePlanRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsagePlanRequest copyWith(void Function(GetUsagePlanRequest) updates) =>
      super.copyWith((message) => updates(message as GetUsagePlanRequest))
          as GetUsagePlanRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetUsagePlanRequest create() => GetUsagePlanRequest._();
  @$core.override
  GetUsagePlanRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetUsagePlanRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetUsagePlanRequest>(create);
  static GetUsagePlanRequest? _defaultInstance;

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(0);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(0);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class GetUsagePlansRequest extends $pb.GeneratedMessage {
  factory GetUsagePlansRequest({
    $core.int? limit,
    $core.String? position,
    $core.String? keyid,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (keyid != null) result.keyid = keyid;
    return result;
  }

  GetUsagePlansRequest._();

  factory GetUsagePlansRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetUsagePlansRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetUsagePlansRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(479913282, _omitFieldNames ? '' : 'keyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsagePlansRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsagePlansRequest copyWith(void Function(GetUsagePlansRequest) updates) =>
      super.copyWith((message) => updates(message as GetUsagePlansRequest))
          as GetUsagePlansRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetUsagePlansRequest create() => GetUsagePlansRequest._();
  @$core.override
  GetUsagePlansRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetUsagePlansRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetUsagePlansRequest>(create);
  static GetUsagePlansRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(479913282)
  $core.String get keyid => $_getSZ(2);
  @$pb.TagNumber(479913282)
  set keyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(479913282)
  $core.bool hasKeyid() => $_has(2);
  @$pb.TagNumber(479913282)
  void clearKeyid() => $_clearField(479913282);
}

class GetUsageRequest extends $pb.GeneratedMessage {
  factory GetUsageRequest({
    $core.int? limit,
    $core.String? position,
    $core.String? startdate,
    $core.String? enddate,
    $core.String? keyid,
    $core.String? usageplanid,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    if (startdate != null) result.startdate = startdate;
    if (enddate != null) result.enddate = enddate;
    if (keyid != null) result.keyid = keyid;
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  GetUsageRequest._();

  factory GetUsageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetUsageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetUsageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(364840732, _omitFieldNames ? '' : 'startdate')
    ..aOS(384831215, _omitFieldNames ? '' : 'enddate')
    ..aOS(479913282, _omitFieldNames ? '' : 'keyid')
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetUsageRequest copyWith(void Function(GetUsageRequest) updates) =>
      super.copyWith((message) => updates(message as GetUsageRequest))
          as GetUsageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetUsageRequest create() => GetUsageRequest._();
  @$core.override
  GetUsageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetUsageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetUsageRequest>(create);
  static GetUsageRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(364840732)
  $core.String get startdate => $_getSZ(2);
  @$pb.TagNumber(364840732)
  set startdate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(364840732)
  $core.bool hasStartdate() => $_has(2);
  @$pb.TagNumber(364840732)
  void clearStartdate() => $_clearField(364840732);

  @$pb.TagNumber(384831215)
  $core.String get enddate => $_getSZ(3);
  @$pb.TagNumber(384831215)
  set enddate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384831215)
  $core.bool hasEnddate() => $_has(3);
  @$pb.TagNumber(384831215)
  void clearEnddate() => $_clearField(384831215);

  @$pb.TagNumber(479913282)
  $core.String get keyid => $_getSZ(4);
  @$pb.TagNumber(479913282)
  set keyid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(479913282)
  $core.bool hasKeyid() => $_has(4);
  @$pb.TagNumber(479913282)
  void clearKeyid() => $_clearField(479913282);

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(5);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(5);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class GetVpcLinkRequest extends $pb.GeneratedMessage {
  factory GetVpcLinkRequest({
    $core.String? vpclinkid,
  }) {
    final result = create();
    if (vpclinkid != null) result.vpclinkid = vpclinkid;
    return result;
  }

  GetVpcLinkRequest._();

  factory GetVpcLinkRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetVpcLinkRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetVpcLinkRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(27515846, _omitFieldNames ? '' : 'vpclinkid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetVpcLinkRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetVpcLinkRequest copyWith(void Function(GetVpcLinkRequest) updates) =>
      super.copyWith((message) => updates(message as GetVpcLinkRequest))
          as GetVpcLinkRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetVpcLinkRequest create() => GetVpcLinkRequest._();
  @$core.override
  GetVpcLinkRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetVpcLinkRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetVpcLinkRequest>(create);
  static GetVpcLinkRequest? _defaultInstance;

  @$pb.TagNumber(27515846)
  $core.String get vpclinkid => $_getSZ(0);
  @$pb.TagNumber(27515846)
  set vpclinkid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(27515846)
  $core.bool hasVpclinkid() => $_has(0);
  @$pb.TagNumber(27515846)
  void clearVpclinkid() => $_clearField(27515846);
}

class GetVpcLinksRequest extends $pb.GeneratedMessage {
  factory GetVpcLinksRequest({
    $core.int? limit,
    $core.String? position,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (position != null) result.position = position;
    return result;
  }

  GetVpcLinksRequest._();

  factory GetVpcLinksRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetVpcLinksRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetVpcLinksRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetVpcLinksRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetVpcLinksRequest copyWith(void Function(GetVpcLinksRequest) updates) =>
      super.copyWith((message) => updates(message as GetVpcLinksRequest))
          as GetVpcLinksRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetVpcLinksRequest create() => GetVpcLinksRequest._();
  @$core.override
  GetVpcLinksRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetVpcLinksRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetVpcLinksRequest>(create);
  static GetVpcLinksRequest? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(1);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(1, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(1);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);
}

class ImportApiKeysRequest extends $pb.GeneratedMessage {
  factory ImportApiKeysRequest({
    ApiKeysFormat? format,
    $core.bool? failonwarnings,
    $core.List<$core.int>? body,
  }) {
    final result = create();
    if (format != null) result.format = format;
    if (failonwarnings != null) result.failonwarnings = failonwarnings;
    if (body != null) result.body = body;
    return result;
  }

  ImportApiKeysRequest._();

  factory ImportApiKeysRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportApiKeysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportApiKeysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aE<ApiKeysFormat>(429753683, _omitFieldNames ? '' : 'format',
        enumValues: ApiKeysFormat.values)
    ..aOB(434363958, _omitFieldNames ? '' : 'failonwarnings')
    ..a<$core.List<$core.int>>(
        464157046, _omitFieldNames ? '' : 'body', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportApiKeysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportApiKeysRequest copyWith(void Function(ImportApiKeysRequest) updates) =>
      super.copyWith((message) => updates(message as ImportApiKeysRequest))
          as ImportApiKeysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportApiKeysRequest create() => ImportApiKeysRequest._();
  @$core.override
  ImportApiKeysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportApiKeysRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportApiKeysRequest>(create);
  static ImportApiKeysRequest? _defaultInstance;

  @$pb.TagNumber(429753683)
  ApiKeysFormat get format => $_getN(0);
  @$pb.TagNumber(429753683)
  set format(ApiKeysFormat value) => $_setField(429753683, value);
  @$pb.TagNumber(429753683)
  $core.bool hasFormat() => $_has(0);
  @$pb.TagNumber(429753683)
  void clearFormat() => $_clearField(429753683);

  @$pb.TagNumber(434363958)
  $core.bool get failonwarnings => $_getBF(1);
  @$pb.TagNumber(434363958)
  set failonwarnings($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(434363958)
  $core.bool hasFailonwarnings() => $_has(1);
  @$pb.TagNumber(434363958)
  void clearFailonwarnings() => $_clearField(434363958);

  @$pb.TagNumber(464157046)
  $core.List<$core.int> get body => $_getN(2);
  @$pb.TagNumber(464157046)
  set body($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(2);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);
}

class ImportDocumentationPartsRequest extends $pb.GeneratedMessage {
  factory ImportDocumentationPartsRequest({
    PutMode? mode,
    $core.String? restapiid,
    $core.bool? failonwarnings,
    $core.List<$core.int>? body,
  }) {
    final result = create();
    if (mode != null) result.mode = mode;
    if (restapiid != null) result.restapiid = restapiid;
    if (failonwarnings != null) result.failonwarnings = failonwarnings;
    if (body != null) result.body = body;
    return result;
  }

  ImportDocumentationPartsRequest._();

  factory ImportDocumentationPartsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportDocumentationPartsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportDocumentationPartsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aE<PutMode>(208592915, _omitFieldNames ? '' : 'mode',
        enumValues: PutMode.values)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOB(434363958, _omitFieldNames ? '' : 'failonwarnings')
    ..a<$core.List<$core.int>>(
        464157046, _omitFieldNames ? '' : 'body', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportDocumentationPartsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportDocumentationPartsRequest copyWith(
          void Function(ImportDocumentationPartsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ImportDocumentationPartsRequest))
          as ImportDocumentationPartsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportDocumentationPartsRequest create() =>
      ImportDocumentationPartsRequest._();
  @$core.override
  ImportDocumentationPartsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportDocumentationPartsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportDocumentationPartsRequest>(
          create);
  static ImportDocumentationPartsRequest? _defaultInstance;

  @$pb.TagNumber(208592915)
  PutMode get mode => $_getN(0);
  @$pb.TagNumber(208592915)
  set mode(PutMode value) => $_setField(208592915, value);
  @$pb.TagNumber(208592915)
  $core.bool hasMode() => $_has(0);
  @$pb.TagNumber(208592915)
  void clearMode() => $_clearField(208592915);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(434363958)
  $core.bool get failonwarnings => $_getBF(2);
  @$pb.TagNumber(434363958)
  set failonwarnings($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(434363958)
  $core.bool hasFailonwarnings() => $_has(2);
  @$pb.TagNumber(434363958)
  void clearFailonwarnings() => $_clearField(434363958);

  @$pb.TagNumber(464157046)
  $core.List<$core.int> get body => $_getN(3);
  @$pb.TagNumber(464157046)
  set body($core.List<$core.int> value) => $_setBytes(3, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(3);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);
}

class ImportRestApiRequest extends $pb.GeneratedMessage {
  factory ImportRestApiRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
    $core.bool? failonwarnings,
    $core.List<$core.int>? body,
  }) {
    final result = create();
    if (parameters != null) result.parameters.addEntries(parameters);
    if (failonwarnings != null) result.failonwarnings = failonwarnings;
    if (body != null) result.body = body;
    return result;
  }

  ImportRestApiRequest._();

  factory ImportRestApiRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportRestApiRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportRestApiRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        145043162, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'ImportRestApiRequest.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOB(434363958, _omitFieldNames ? '' : 'failonwarnings')
    ..a<$core.List<$core.int>>(
        464157046, _omitFieldNames ? '' : 'body', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportRestApiRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportRestApiRequest copyWith(void Function(ImportRestApiRequest) updates) =>
      super.copyWith((message) => updates(message as ImportRestApiRequest))
          as ImportRestApiRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportRestApiRequest create() => ImportRestApiRequest._();
  @$core.override
  ImportRestApiRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportRestApiRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportRestApiRequest>(create);
  static ImportRestApiRequest? _defaultInstance;

  @$pb.TagNumber(145043162)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(0);

  @$pb.TagNumber(434363958)
  $core.bool get failonwarnings => $_getBF(1);
  @$pb.TagNumber(434363958)
  set failonwarnings($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(434363958)
  $core.bool hasFailonwarnings() => $_has(1);
  @$pb.TagNumber(434363958)
  void clearFailonwarnings() => $_clearField(434363958);

  @$pb.TagNumber(464157046)
  $core.List<$core.int> get body => $_getN(2);
  @$pb.TagNumber(464157046)
  set body($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(2);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);
}

class Integration extends $pb.GeneratedMessage {
  factory Integration({
    $core.String? integrationtarget,
    $core.String? cachenamespace,
    TlsConfig? tlsconfig,
    $core.String? httpmethod,
    $core.String? credentials,
    IntegrationType? type,
    $core.String? passthroughbehavior,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        requesttemplates,
    ConnectionType? connectiontype,
    $core.int? timeoutinmillis,
    $core.Iterable<$core.MapEntry<$core.String, IntegrationResponse>>?
        integrationresponses,
    $core.String? uri,
    $core.String? connectionid,
    ResponseTransferMode? responsetransfermode,
    $core.Iterable<$core.String>? cachekeyparameters,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        requestparameters,
    ContentHandlingStrategy? contenthandling,
  }) {
    final result = create();
    if (integrationtarget != null) result.integrationtarget = integrationtarget;
    if (cachenamespace != null) result.cachenamespace = cachenamespace;
    if (tlsconfig != null) result.tlsconfig = tlsconfig;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (credentials != null) result.credentials = credentials;
    if (type != null) result.type = type;
    if (passthroughbehavior != null)
      result.passthroughbehavior = passthroughbehavior;
    if (requesttemplates != null)
      result.requesttemplates.addEntries(requesttemplates);
    if (connectiontype != null) result.connectiontype = connectiontype;
    if (timeoutinmillis != null) result.timeoutinmillis = timeoutinmillis;
    if (integrationresponses != null)
      result.integrationresponses.addEntries(integrationresponses);
    if (uri != null) result.uri = uri;
    if (connectionid != null) result.connectionid = connectionid;
    if (responsetransfermode != null)
      result.responsetransfermode = responsetransfermode;
    if (cachekeyparameters != null)
      result.cachekeyparameters.addAll(cachekeyparameters);
    if (requestparameters != null)
      result.requestparameters.addEntries(requestparameters);
    if (contenthandling != null) result.contenthandling = contenthandling;
    return result;
  }

  Integration._();

  factory Integration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Integration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Integration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(17646705, _omitFieldNames ? '' : 'integrationtarget')
    ..aOS(85102753, _omitFieldNames ? '' : 'cachenamespace')
    ..aOM<TlsConfig>(108946693, _omitFieldNames ? '' : 'tlsconfig',
        subBuilder: TlsConfig.create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(150838226, _omitFieldNames ? '' : 'credentials')
    ..aE<IntegrationType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: IntegrationType.values)
    ..aOS(310796908, _omitFieldNames ? '' : 'passthroughbehavior')
    ..m<$core.String, $core.String>(
        333512166, _omitFieldNames ? '' : 'requesttemplates',
        entryClassName: 'Integration.RequesttemplatesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aE<ConnectionType>(336253170, _omitFieldNames ? '' : 'connectiontype',
        enumValues: ConnectionType.values)
    ..aI(378229126, _omitFieldNames ? '' : 'timeoutinmillis')
    ..m<$core.String, IntegrationResponse>(
        386580464, _omitFieldNames ? '' : 'integrationresponses',
        entryClassName: 'Integration.IntegrationresponsesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: IntegrationResponse.create,
        valueDefaultOrMaker: IntegrationResponse.getDefault,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(395269118, _omitFieldNames ? '' : 'uri')
    ..aOS(450027965, _omitFieldNames ? '' : 'connectionid')
    ..aE<ResponseTransferMode>(
        458910787, _omitFieldNames ? '' : 'responsetransfermode',
        enumValues: ResponseTransferMode.values)
    ..pPS(481441313, _omitFieldNames ? '' : 'cachekeyparameters')
    ..m<$core.String, $core.String>(
        523499939, _omitFieldNames ? '' : 'requestparameters',
        entryClassName: 'Integration.RequestparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aE<ContentHandlingStrategy>(
        533182832, _omitFieldNames ? '' : 'contenthandling',
        enumValues: ContentHandlingStrategy.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Integration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Integration copyWith(void Function(Integration) updates) =>
      super.copyWith((message) => updates(message as Integration))
          as Integration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Integration create() => Integration._();
  @$core.override
  Integration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Integration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Integration>(create);
  static Integration? _defaultInstance;

  @$pb.TagNumber(17646705)
  $core.String get integrationtarget => $_getSZ(0);
  @$pb.TagNumber(17646705)
  set integrationtarget($core.String value) => $_setString(0, value);
  @$pb.TagNumber(17646705)
  $core.bool hasIntegrationtarget() => $_has(0);
  @$pb.TagNumber(17646705)
  void clearIntegrationtarget() => $_clearField(17646705);

  @$pb.TagNumber(85102753)
  $core.String get cachenamespace => $_getSZ(1);
  @$pb.TagNumber(85102753)
  set cachenamespace($core.String value) => $_setString(1, value);
  @$pb.TagNumber(85102753)
  $core.bool hasCachenamespace() => $_has(1);
  @$pb.TagNumber(85102753)
  void clearCachenamespace() => $_clearField(85102753);

  @$pb.TagNumber(108946693)
  TlsConfig get tlsconfig => $_getN(2);
  @$pb.TagNumber(108946693)
  set tlsconfig(TlsConfig value) => $_setField(108946693, value);
  @$pb.TagNumber(108946693)
  $core.bool hasTlsconfig() => $_has(2);
  @$pb.TagNumber(108946693)
  void clearTlsconfig() => $_clearField(108946693);
  @$pb.TagNumber(108946693)
  TlsConfig ensureTlsconfig() => $_ensure(2);

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(3);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(3);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(150838226)
  $core.String get credentials => $_getSZ(4);
  @$pb.TagNumber(150838226)
  set credentials($core.String value) => $_setString(4, value);
  @$pb.TagNumber(150838226)
  $core.bool hasCredentials() => $_has(4);
  @$pb.TagNumber(150838226)
  void clearCredentials() => $_clearField(150838226);

  @$pb.TagNumber(287830350)
  IntegrationType get type => $_getN(5);
  @$pb.TagNumber(287830350)
  set type(IntegrationType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(5);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(310796908)
  $core.String get passthroughbehavior => $_getSZ(6);
  @$pb.TagNumber(310796908)
  set passthroughbehavior($core.String value) => $_setString(6, value);
  @$pb.TagNumber(310796908)
  $core.bool hasPassthroughbehavior() => $_has(6);
  @$pb.TagNumber(310796908)
  void clearPassthroughbehavior() => $_clearField(310796908);

  @$pb.TagNumber(333512166)
  $pb.PbMap<$core.String, $core.String> get requesttemplates => $_getMap(7);

  @$pb.TagNumber(336253170)
  ConnectionType get connectiontype => $_getN(8);
  @$pb.TagNumber(336253170)
  set connectiontype(ConnectionType value) => $_setField(336253170, value);
  @$pb.TagNumber(336253170)
  $core.bool hasConnectiontype() => $_has(8);
  @$pb.TagNumber(336253170)
  void clearConnectiontype() => $_clearField(336253170);

  @$pb.TagNumber(378229126)
  $core.int get timeoutinmillis => $_getIZ(9);
  @$pb.TagNumber(378229126)
  set timeoutinmillis($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(378229126)
  $core.bool hasTimeoutinmillis() => $_has(9);
  @$pb.TagNumber(378229126)
  void clearTimeoutinmillis() => $_clearField(378229126);

  @$pb.TagNumber(386580464)
  $pb.PbMap<$core.String, IntegrationResponse> get integrationresponses =>
      $_getMap(10);

  @$pb.TagNumber(395269118)
  $core.String get uri => $_getSZ(11);
  @$pb.TagNumber(395269118)
  set uri($core.String value) => $_setString(11, value);
  @$pb.TagNumber(395269118)
  $core.bool hasUri() => $_has(11);
  @$pb.TagNumber(395269118)
  void clearUri() => $_clearField(395269118);

  @$pb.TagNumber(450027965)
  $core.String get connectionid => $_getSZ(12);
  @$pb.TagNumber(450027965)
  set connectionid($core.String value) => $_setString(12, value);
  @$pb.TagNumber(450027965)
  $core.bool hasConnectionid() => $_has(12);
  @$pb.TagNumber(450027965)
  void clearConnectionid() => $_clearField(450027965);

  @$pb.TagNumber(458910787)
  ResponseTransferMode get responsetransfermode => $_getN(13);
  @$pb.TagNumber(458910787)
  set responsetransfermode(ResponseTransferMode value) =>
      $_setField(458910787, value);
  @$pb.TagNumber(458910787)
  $core.bool hasResponsetransfermode() => $_has(13);
  @$pb.TagNumber(458910787)
  void clearResponsetransfermode() => $_clearField(458910787);

  @$pb.TagNumber(481441313)
  $pb.PbList<$core.String> get cachekeyparameters => $_getList(14);

  @$pb.TagNumber(523499939)
  $pb.PbMap<$core.String, $core.String> get requestparameters => $_getMap(15);

  @$pb.TagNumber(533182832)
  ContentHandlingStrategy get contenthandling => $_getN(16);
  @$pb.TagNumber(533182832)
  set contenthandling(ContentHandlingStrategy value) =>
      $_setField(533182832, value);
  @$pb.TagNumber(533182832)
  $core.bool hasContenthandling() => $_has(16);
  @$pb.TagNumber(533182832)
  void clearContenthandling() => $_clearField(533182832);
}

class IntegrationResponse extends $pb.GeneratedMessage {
  factory IntegrationResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        responseparameters,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        responsetemplates,
    $core.String? statuscode,
    $core.String? selectionpattern,
    ContentHandlingStrategy? contenthandling,
  }) {
    final result = create();
    if (responseparameters != null)
      result.responseparameters.addEntries(responseparameters);
    if (responsetemplates != null)
      result.responsetemplates.addEntries(responsetemplates);
    if (statuscode != null) result.statuscode = statuscode;
    if (selectionpattern != null) result.selectionpattern = selectionpattern;
    if (contenthandling != null) result.contenthandling = contenthandling;
    return result;
  }

  IntegrationResponse._();

  factory IntegrationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IntegrationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IntegrationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        64271839, _omitFieldNames ? '' : 'responseparameters',
        entryClassName: 'IntegrationResponse.ResponseparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..m<$core.String, $core.String>(
        107376570, _omitFieldNames ? '' : 'responsetemplates',
        entryClassName: 'IntegrationResponse.ResponsetemplatesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(470634042, _omitFieldNames ? '' : 'selectionpattern')
    ..aE<ContentHandlingStrategy>(
        533182832, _omitFieldNames ? '' : 'contenthandling',
        enumValues: ContentHandlingStrategy.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IntegrationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IntegrationResponse copyWith(void Function(IntegrationResponse) updates) =>
      super.copyWith((message) => updates(message as IntegrationResponse))
          as IntegrationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IntegrationResponse create() => IntegrationResponse._();
  @$core.override
  IntegrationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IntegrationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IntegrationResponse>(create);
  static IntegrationResponse? _defaultInstance;

  @$pb.TagNumber(64271839)
  $pb.PbMap<$core.String, $core.String> get responseparameters => $_getMap(0);

  @$pb.TagNumber(107376570)
  $pb.PbMap<$core.String, $core.String> get responsetemplates => $_getMap(1);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(2);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(2);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(470634042)
  $core.String get selectionpattern => $_getSZ(3);
  @$pb.TagNumber(470634042)
  set selectionpattern($core.String value) => $_setString(3, value);
  @$pb.TagNumber(470634042)
  $core.bool hasSelectionpattern() => $_has(3);
  @$pb.TagNumber(470634042)
  void clearSelectionpattern() => $_clearField(470634042);

  @$pb.TagNumber(533182832)
  ContentHandlingStrategy get contenthandling => $_getN(4);
  @$pb.TagNumber(533182832)
  set contenthandling(ContentHandlingStrategy value) =>
      $_setField(533182832, value);
  @$pb.TagNumber(533182832)
  $core.bool hasContenthandling() => $_has(4);
  @$pb.TagNumber(533182832)
  void clearContenthandling() => $_clearField(533182832);
}

class LimitExceededException extends $pb.GeneratedMessage {
  factory LimitExceededException({
    $core.String? message,
    $core.String? retryafterseconds,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (retryafterseconds != null) result.retryafterseconds = retryafterseconds;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..aOS(436039555, _omitFieldNames ? '' : 'retryafterseconds')
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

  @$pb.TagNumber(436039555)
  $core.String get retryafterseconds => $_getSZ(1);
  @$pb.TagNumber(436039555)
  set retryafterseconds($core.String value) => $_setString(1, value);
  @$pb.TagNumber(436039555)
  $core.bool hasRetryafterseconds() => $_has(1);
  @$pb.TagNumber(436039555)
  void clearRetryafterseconds() => $_clearField(436039555);
}

class Method extends $pb.GeneratedMessage {
  factory Method({
    $core.String? authorizerid,
    $core.String? httpmethod,
    $core.String? operationname,
    $core.Iterable<$core.MapEntry<$core.String, MethodResponse>>?
        methodresponses,
    $core.String? authorizationtype,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? requestmodels,
    $core.Iterable<$core.String>? authorizationscopes,
    $core.bool? apikeyrequired,
    $core.String? requestvalidatorid,
    Integration? methodintegration,
    $core.Iterable<$core.MapEntry<$core.String, $core.bool>>? requestparameters,
  }) {
    final result = create();
    if (authorizerid != null) result.authorizerid = authorizerid;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (operationname != null) result.operationname = operationname;
    if (methodresponses != null)
      result.methodresponses.addEntries(methodresponses);
    if (authorizationtype != null) result.authorizationtype = authorizationtype;
    if (requestmodels != null) result.requestmodels.addEntries(requestmodels);
    if (authorizationscopes != null)
      result.authorizationscopes.addAll(authorizationscopes);
    if (apikeyrequired != null) result.apikeyrequired = apikeyrequired;
    if (requestvalidatorid != null)
      result.requestvalidatorid = requestvalidatorid;
    if (methodintegration != null) result.methodintegration = methodintegration;
    if (requestparameters != null)
      result.requestparameters.addEntries(requestparameters);
    return result;
  }

  Method._();

  factory Method.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Method.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Method',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(111773148, _omitFieldNames ? '' : 'authorizerid')
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(178909574, _omitFieldNames ? '' : 'operationname')
    ..m<$core.String, MethodResponse>(
        231818421, _omitFieldNames ? '' : 'methodresponses',
        entryClassName: 'Method.MethodresponsesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MethodResponse.create,
        valueDefaultOrMaker: MethodResponse.getDefault,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(387986911, _omitFieldNames ? '' : 'authorizationtype')
    ..m<$core.String, $core.String>(
        397252853, _omitFieldNames ? '' : 'requestmodels',
        entryClassName: 'Method.RequestmodelsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..pPS(423149932, _omitFieldNames ? '' : 'authorizationscopes')
    ..aOB(435360152, _omitFieldNames ? '' : 'apikeyrequired')
    ..aOS(517546134, _omitFieldNames ? '' : 'requestvalidatorid')
    ..aOM<Integration>(518245059, _omitFieldNames ? '' : 'methodintegration',
        subBuilder: Integration.create)
    ..m<$core.String, $core.bool>(
        523499939, _omitFieldNames ? '' : 'requestparameters',
        entryClassName: 'Method.RequestparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OB,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Method clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Method copyWith(void Function(Method) updates) =>
      super.copyWith((message) => updates(message as Method)) as Method;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Method create() => Method._();
  @$core.override
  Method createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Method getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Method>(create);
  static Method? _defaultInstance;

  @$pb.TagNumber(111773148)
  $core.String get authorizerid => $_getSZ(0);
  @$pb.TagNumber(111773148)
  set authorizerid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(111773148)
  $core.bool hasAuthorizerid() => $_has(0);
  @$pb.TagNumber(111773148)
  void clearAuthorizerid() => $_clearField(111773148);

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(1);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(1);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(178909574)
  $core.String get operationname => $_getSZ(2);
  @$pb.TagNumber(178909574)
  set operationname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(178909574)
  $core.bool hasOperationname() => $_has(2);
  @$pb.TagNumber(178909574)
  void clearOperationname() => $_clearField(178909574);

  @$pb.TagNumber(231818421)
  $pb.PbMap<$core.String, MethodResponse> get methodresponses => $_getMap(3);

  @$pb.TagNumber(387986911)
  $core.String get authorizationtype => $_getSZ(4);
  @$pb.TagNumber(387986911)
  set authorizationtype($core.String value) => $_setString(4, value);
  @$pb.TagNumber(387986911)
  $core.bool hasAuthorizationtype() => $_has(4);
  @$pb.TagNumber(387986911)
  void clearAuthorizationtype() => $_clearField(387986911);

  @$pb.TagNumber(397252853)
  $pb.PbMap<$core.String, $core.String> get requestmodels => $_getMap(5);

  @$pb.TagNumber(423149932)
  $pb.PbList<$core.String> get authorizationscopes => $_getList(6);

  @$pb.TagNumber(435360152)
  $core.bool get apikeyrequired => $_getBF(7);
  @$pb.TagNumber(435360152)
  set apikeyrequired($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(435360152)
  $core.bool hasApikeyrequired() => $_has(7);
  @$pb.TagNumber(435360152)
  void clearApikeyrequired() => $_clearField(435360152);

  @$pb.TagNumber(517546134)
  $core.String get requestvalidatorid => $_getSZ(8);
  @$pb.TagNumber(517546134)
  set requestvalidatorid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(517546134)
  $core.bool hasRequestvalidatorid() => $_has(8);
  @$pb.TagNumber(517546134)
  void clearRequestvalidatorid() => $_clearField(517546134);

  @$pb.TagNumber(518245059)
  Integration get methodintegration => $_getN(9);
  @$pb.TagNumber(518245059)
  set methodintegration(Integration value) => $_setField(518245059, value);
  @$pb.TagNumber(518245059)
  $core.bool hasMethodintegration() => $_has(9);
  @$pb.TagNumber(518245059)
  void clearMethodintegration() => $_clearField(518245059);
  @$pb.TagNumber(518245059)
  Integration ensureMethodintegration() => $_ensure(9);

  @$pb.TagNumber(523499939)
  $pb.PbMap<$core.String, $core.bool> get requestparameters => $_getMap(10);
}

class MethodResponse extends $pb.GeneratedMessage {
  factory MethodResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.bool>>?
        responseparameters,
    $core.String? statuscode,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? responsemodels,
  }) {
    final result = create();
    if (responseparameters != null)
      result.responseparameters.addEntries(responseparameters);
    if (statuscode != null) result.statuscode = statuscode;
    if (responsemodels != null)
      result.responsemodels.addEntries(responsemodels);
    return result;
  }

  MethodResponse._();

  factory MethodResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MethodResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MethodResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.bool>(
        64271839, _omitFieldNames ? '' : 'responseparameters',
        entryClassName: 'MethodResponse.ResponseparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OB,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..m<$core.String, $core.String>(
        356574313, _omitFieldNames ? '' : 'responsemodels',
        entryClassName: 'MethodResponse.ResponsemodelsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MethodResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MethodResponse copyWith(void Function(MethodResponse) updates) =>
      super.copyWith((message) => updates(message as MethodResponse))
          as MethodResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MethodResponse create() => MethodResponse._();
  @$core.override
  MethodResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MethodResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MethodResponse>(create);
  static MethodResponse? _defaultInstance;

  @$pb.TagNumber(64271839)
  $pb.PbMap<$core.String, $core.bool> get responseparameters => $_getMap(0);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(1);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(1);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(356574313)
  $pb.PbMap<$core.String, $core.String> get responsemodels => $_getMap(2);
}

class MethodSetting extends $pb.GeneratedMessage {
  factory MethodSetting({
    $core.String? logginglevel,
    $core.int? cachettlinseconds,
    $core.bool? metricsenabled,
    $core.bool? datatraceenabled,
    $core.double? throttlingratelimit,
    $core.int? throttlingburstlimit,
    UnauthorizedCacheControlHeaderStrategy?
        unauthorizedcachecontrolheaderstrategy,
    $core.bool? cachedataencrypted,
    $core.bool? cachingenabled,
    $core.bool? requireauthorizationforcachecontrol,
  }) {
    final result = create();
    if (logginglevel != null) result.logginglevel = logginglevel;
    if (cachettlinseconds != null) result.cachettlinseconds = cachettlinseconds;
    if (metricsenabled != null) result.metricsenabled = metricsenabled;
    if (datatraceenabled != null) result.datatraceenabled = datatraceenabled;
    if (throttlingratelimit != null)
      result.throttlingratelimit = throttlingratelimit;
    if (throttlingburstlimit != null)
      result.throttlingburstlimit = throttlingburstlimit;
    if (unauthorizedcachecontrolheaderstrategy != null)
      result.unauthorizedcachecontrolheaderstrategy =
          unauthorizedcachecontrolheaderstrategy;
    if (cachedataencrypted != null)
      result.cachedataencrypted = cachedataencrypted;
    if (cachingenabled != null) result.cachingenabled = cachingenabled;
    if (requireauthorizationforcachecontrol != null)
      result.requireauthorizationforcachecontrol =
          requireauthorizationforcachecontrol;
    return result;
  }

  MethodSetting._();

  factory MethodSetting.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MethodSetting.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MethodSetting',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(59396637, _omitFieldNames ? '' : 'logginglevel')
    ..aI(79996982, _omitFieldNames ? '' : 'cachettlinseconds')
    ..aOB(142460292, _omitFieldNames ? '' : 'metricsenabled')
    ..aOB(363519852, _omitFieldNames ? '' : 'datatraceenabled')
    ..aD(371718088, _omitFieldNames ? '' : 'throttlingratelimit')
    ..aI(402901688, _omitFieldNames ? '' : 'throttlingburstlimit')
    ..aE<UnauthorizedCacheControlHeaderStrategy>(476741277,
        _omitFieldNames ? '' : 'unauthorizedcachecontrolheaderstrategy',
        enumValues: UnauthorizedCacheControlHeaderStrategy.values)
    ..aOB(481075060, _omitFieldNames ? '' : 'cachedataencrypted')
    ..aOB(489524028, _omitFieldNames ? '' : 'cachingenabled')
    ..aOB(
        529394912, _omitFieldNames ? '' : 'requireauthorizationforcachecontrol')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MethodSetting clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MethodSetting copyWith(void Function(MethodSetting) updates) =>
      super.copyWith((message) => updates(message as MethodSetting))
          as MethodSetting;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MethodSetting create() => MethodSetting._();
  @$core.override
  MethodSetting createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MethodSetting getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MethodSetting>(create);
  static MethodSetting? _defaultInstance;

  @$pb.TagNumber(59396637)
  $core.String get logginglevel => $_getSZ(0);
  @$pb.TagNumber(59396637)
  set logginglevel($core.String value) => $_setString(0, value);
  @$pb.TagNumber(59396637)
  $core.bool hasLogginglevel() => $_has(0);
  @$pb.TagNumber(59396637)
  void clearLogginglevel() => $_clearField(59396637);

  @$pb.TagNumber(79996982)
  $core.int get cachettlinseconds => $_getIZ(1);
  @$pb.TagNumber(79996982)
  set cachettlinseconds($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(79996982)
  $core.bool hasCachettlinseconds() => $_has(1);
  @$pb.TagNumber(79996982)
  void clearCachettlinseconds() => $_clearField(79996982);

  @$pb.TagNumber(142460292)
  $core.bool get metricsenabled => $_getBF(2);
  @$pb.TagNumber(142460292)
  set metricsenabled($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(142460292)
  $core.bool hasMetricsenabled() => $_has(2);
  @$pb.TagNumber(142460292)
  void clearMetricsenabled() => $_clearField(142460292);

  @$pb.TagNumber(363519852)
  $core.bool get datatraceenabled => $_getBF(3);
  @$pb.TagNumber(363519852)
  set datatraceenabled($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(363519852)
  $core.bool hasDatatraceenabled() => $_has(3);
  @$pb.TagNumber(363519852)
  void clearDatatraceenabled() => $_clearField(363519852);

  @$pb.TagNumber(371718088)
  $core.double get throttlingratelimit => $_getN(4);
  @$pb.TagNumber(371718088)
  set throttlingratelimit($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(371718088)
  $core.bool hasThrottlingratelimit() => $_has(4);
  @$pb.TagNumber(371718088)
  void clearThrottlingratelimit() => $_clearField(371718088);

  @$pb.TagNumber(402901688)
  $core.int get throttlingburstlimit => $_getIZ(5);
  @$pb.TagNumber(402901688)
  set throttlingburstlimit($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(402901688)
  $core.bool hasThrottlingburstlimit() => $_has(5);
  @$pb.TagNumber(402901688)
  void clearThrottlingburstlimit() => $_clearField(402901688);

  @$pb.TagNumber(476741277)
  UnauthorizedCacheControlHeaderStrategy
      get unauthorizedcachecontrolheaderstrategy => $_getN(6);
  @$pb.TagNumber(476741277)
  set unauthorizedcachecontrolheaderstrategy(
          UnauthorizedCacheControlHeaderStrategy value) =>
      $_setField(476741277, value);
  @$pb.TagNumber(476741277)
  $core.bool hasUnauthorizedcachecontrolheaderstrategy() => $_has(6);
  @$pb.TagNumber(476741277)
  void clearUnauthorizedcachecontrolheaderstrategy() => $_clearField(476741277);

  @$pb.TagNumber(481075060)
  $core.bool get cachedataencrypted => $_getBF(7);
  @$pb.TagNumber(481075060)
  set cachedataencrypted($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(481075060)
  $core.bool hasCachedataencrypted() => $_has(7);
  @$pb.TagNumber(481075060)
  void clearCachedataencrypted() => $_clearField(481075060);

  @$pb.TagNumber(489524028)
  $core.bool get cachingenabled => $_getBF(8);
  @$pb.TagNumber(489524028)
  set cachingenabled($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(489524028)
  $core.bool hasCachingenabled() => $_has(8);
  @$pb.TagNumber(489524028)
  void clearCachingenabled() => $_clearField(489524028);

  @$pb.TagNumber(529394912)
  $core.bool get requireauthorizationforcachecontrol => $_getBF(9);
  @$pb.TagNumber(529394912)
  set requireauthorizationforcachecontrol($core.bool value) =>
      $_setBool(9, value);
  @$pb.TagNumber(529394912)
  $core.bool hasRequireauthorizationforcachecontrol() => $_has(9);
  @$pb.TagNumber(529394912)
  void clearRequireauthorizationforcachecontrol() => $_clearField(529394912);
}

class MethodSnapshot extends $pb.GeneratedMessage {
  factory MethodSnapshot({
    $core.String? authorizationtype,
    $core.bool? apikeyrequired,
  }) {
    final result = create();
    if (authorizationtype != null) result.authorizationtype = authorizationtype;
    if (apikeyrequired != null) result.apikeyrequired = apikeyrequired;
    return result;
  }

  MethodSnapshot._();

  factory MethodSnapshot.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MethodSnapshot.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MethodSnapshot',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(387986911, _omitFieldNames ? '' : 'authorizationtype')
    ..aOB(435360152, _omitFieldNames ? '' : 'apikeyrequired')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MethodSnapshot clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MethodSnapshot copyWith(void Function(MethodSnapshot) updates) =>
      super.copyWith((message) => updates(message as MethodSnapshot))
          as MethodSnapshot;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MethodSnapshot create() => MethodSnapshot._();
  @$core.override
  MethodSnapshot createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MethodSnapshot getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MethodSnapshot>(create);
  static MethodSnapshot? _defaultInstance;

  @$pb.TagNumber(387986911)
  $core.String get authorizationtype => $_getSZ(0);
  @$pb.TagNumber(387986911)
  set authorizationtype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(387986911)
  $core.bool hasAuthorizationtype() => $_has(0);
  @$pb.TagNumber(387986911)
  void clearAuthorizationtype() => $_clearField(387986911);

  @$pb.TagNumber(435360152)
  $core.bool get apikeyrequired => $_getBF(1);
  @$pb.TagNumber(435360152)
  set apikeyrequired($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(435360152)
  $core.bool hasApikeyrequired() => $_has(1);
  @$pb.TagNumber(435360152)
  void clearApikeyrequired() => $_clearField(435360152);
}

class Model extends $pb.GeneratedMessage {
  factory Model({
    $core.String? name,
    $core.String? contenttype,
    $core.String? schema,
    $core.String? description,
    $core.String? id,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (contenttype != null) result.contenttype = contenttype;
    if (schema != null) result.schema = schema;
    if (description != null) result.description = description;
    if (id != null) result.id = id;
    return result;
  }

  Model._();

  factory Model.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Model.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Model',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(281764659, _omitFieldNames ? '' : 'contenttype')
    ..aOS(310182711, _omitFieldNames ? '' : 'schema')
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Model clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Model copyWith(void Function(Model) updates) =>
      super.copyWith((message) => updates(message as Model)) as Model;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Model create() => Model._();
  @$core.override
  Model createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Model getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Model>(create);
  static Model? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(281764659)
  $core.String get contenttype => $_getSZ(1);
  @$pb.TagNumber(281764659)
  set contenttype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(281764659)
  $core.bool hasContenttype() => $_has(1);
  @$pb.TagNumber(281764659)
  void clearContenttype() => $_clearField(281764659);

  @$pb.TagNumber(310182711)
  $core.String get schema => $_getSZ(2);
  @$pb.TagNumber(310182711)
  set schema($core.String value) => $_setString(2, value);
  @$pb.TagNumber(310182711)
  $core.bool hasSchema() => $_has(2);
  @$pb.TagNumber(310182711)
  void clearSchema() => $_clearField(310182711);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);
}

class Models extends $pb.GeneratedMessage {
  factory Models({
    $core.String? position,
    $core.Iterable<Model>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  Models._();

  factory Models.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Models.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Models',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<Model>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: Model.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Models clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Models copyWith(void Function(Models) updates) =>
      super.copyWith((message) => updates(message as Models)) as Models;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Models create() => Models._();
  @$core.override
  Models createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Models getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Models>(create);
  static Models? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<Model> get items => $_getList(1);
}

class MutualTlsAuthentication extends $pb.GeneratedMessage {
  factory MutualTlsAuthentication({
    $core.String? truststoreversion,
    $core.String? truststoreuri,
    $core.Iterable<$core.String>? truststorewarnings,
  }) {
    final result = create();
    if (truststoreversion != null) result.truststoreversion = truststoreversion;
    if (truststoreuri != null) result.truststoreuri = truststoreuri;
    if (truststorewarnings != null)
      result.truststorewarnings.addAll(truststorewarnings);
    return result;
  }

  MutualTlsAuthentication._();

  factory MutualTlsAuthentication.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MutualTlsAuthentication.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MutualTlsAuthentication',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(119080291, _omitFieldNames ? '' : 'truststoreversion')
    ..aOS(120246545, _omitFieldNames ? '' : 'truststoreuri')
    ..pPS(420536820, _omitFieldNames ? '' : 'truststorewarnings')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MutualTlsAuthentication clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MutualTlsAuthentication copyWith(
          void Function(MutualTlsAuthentication) updates) =>
      super.copyWith((message) => updates(message as MutualTlsAuthentication))
          as MutualTlsAuthentication;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MutualTlsAuthentication create() => MutualTlsAuthentication._();
  @$core.override
  MutualTlsAuthentication createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MutualTlsAuthentication getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MutualTlsAuthentication>(create);
  static MutualTlsAuthentication? _defaultInstance;

  @$pb.TagNumber(119080291)
  $core.String get truststoreversion => $_getSZ(0);
  @$pb.TagNumber(119080291)
  set truststoreversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(119080291)
  $core.bool hasTruststoreversion() => $_has(0);
  @$pb.TagNumber(119080291)
  void clearTruststoreversion() => $_clearField(119080291);

  @$pb.TagNumber(120246545)
  $core.String get truststoreuri => $_getSZ(1);
  @$pb.TagNumber(120246545)
  set truststoreuri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(120246545)
  $core.bool hasTruststoreuri() => $_has(1);
  @$pb.TagNumber(120246545)
  void clearTruststoreuri() => $_clearField(120246545);

  @$pb.TagNumber(420536820)
  $pb.PbList<$core.String> get truststorewarnings => $_getList(2);
}

class MutualTlsAuthenticationInput extends $pb.GeneratedMessage {
  factory MutualTlsAuthenticationInput({
    $core.String? truststoreversion,
    $core.String? truststoreuri,
  }) {
    final result = create();
    if (truststoreversion != null) result.truststoreversion = truststoreversion;
    if (truststoreuri != null) result.truststoreuri = truststoreuri;
    return result;
  }

  MutualTlsAuthenticationInput._();

  factory MutualTlsAuthenticationInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MutualTlsAuthenticationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MutualTlsAuthenticationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(119080291, _omitFieldNames ? '' : 'truststoreversion')
    ..aOS(120246545, _omitFieldNames ? '' : 'truststoreuri')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MutualTlsAuthenticationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MutualTlsAuthenticationInput copyWith(
          void Function(MutualTlsAuthenticationInput) updates) =>
      super.copyWith(
              (message) => updates(message as MutualTlsAuthenticationInput))
          as MutualTlsAuthenticationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MutualTlsAuthenticationInput create() =>
      MutualTlsAuthenticationInput._();
  @$core.override
  MutualTlsAuthenticationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MutualTlsAuthenticationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MutualTlsAuthenticationInput>(create);
  static MutualTlsAuthenticationInput? _defaultInstance;

  @$pb.TagNumber(119080291)
  $core.String get truststoreversion => $_getSZ(0);
  @$pb.TagNumber(119080291)
  set truststoreversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(119080291)
  $core.bool hasTruststoreversion() => $_has(0);
  @$pb.TagNumber(119080291)
  void clearTruststoreversion() => $_clearField(119080291);

  @$pb.TagNumber(120246545)
  $core.String get truststoreuri => $_getSZ(1);
  @$pb.TagNumber(120246545)
  set truststoreuri($core.String value) => $_setString(1, value);
  @$pb.TagNumber(120246545)
  $core.bool hasTruststoreuri() => $_has(1);
  @$pb.TagNumber(120246545)
  void clearTruststoreuri() => $_clearField(120246545);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
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

class PatchOperation extends $pb.GeneratedMessage {
  factory PatchOperation({
    $core.String? value,
    $core.String? path,
    $core.String? from,
    Op? op,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (path != null) result.path = path;
    if (from != null) result.from = from;
    if (op != null) result.op = op;
    return result;
  }

  PatchOperation._();

  factory PatchOperation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PatchOperation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PatchOperation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(39769035, _omitFieldNames ? '' : 'value')
    ..aOS(75975991, _omitFieldNames ? '' : 'path')
    ..aOS(365789302, _omitFieldNames ? '' : 'from')
    ..aE<Op>(523513003, _omitFieldNames ? '' : 'op', enumValues: Op.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PatchOperation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PatchOperation copyWith(void Function(PatchOperation) updates) =>
      super.copyWith((message) => updates(message as PatchOperation))
          as PatchOperation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PatchOperation create() => PatchOperation._();
  @$core.override
  PatchOperation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PatchOperation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PatchOperation>(create);
  static PatchOperation? _defaultInstance;

  @$pb.TagNumber(39769035)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(39769035)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39769035)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(39769035)
  void clearValue() => $_clearField(39769035);

  @$pb.TagNumber(75975991)
  $core.String get path => $_getSZ(1);
  @$pb.TagNumber(75975991)
  set path($core.String value) => $_setString(1, value);
  @$pb.TagNumber(75975991)
  $core.bool hasPath() => $_has(1);
  @$pb.TagNumber(75975991)
  void clearPath() => $_clearField(75975991);

  @$pb.TagNumber(365789302)
  $core.String get from => $_getSZ(2);
  @$pb.TagNumber(365789302)
  set from($core.String value) => $_setString(2, value);
  @$pb.TagNumber(365789302)
  $core.bool hasFrom() => $_has(2);
  @$pb.TagNumber(365789302)
  void clearFrom() => $_clearField(365789302);

  @$pb.TagNumber(523513003)
  Op get op => $_getN(3);
  @$pb.TagNumber(523513003)
  set op(Op value) => $_setField(523513003, value);
  @$pb.TagNumber(523513003)
  $core.bool hasOp() => $_has(3);
  @$pb.TagNumber(523513003)
  void clearOp() => $_clearField(523513003);
}

class PutGatewayResponseRequest extends $pb.GeneratedMessage {
  factory PutGatewayResponseRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        responseparameters,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        responsetemplates,
    $core.String? statuscode,
    GatewayResponseType? responsetype,
    $core.String? restapiid,
  }) {
    final result = create();
    if (responseparameters != null)
      result.responseparameters.addEntries(responseparameters);
    if (responsetemplates != null)
      result.responsetemplates.addEntries(responsetemplates);
    if (statuscode != null) result.statuscode = statuscode;
    if (responsetype != null) result.responsetype = responsetype;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  PutGatewayResponseRequest._();

  factory PutGatewayResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutGatewayResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutGatewayResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        64271839, _omitFieldNames ? '' : 'responseparameters',
        entryClassName: 'PutGatewayResponseRequest.ResponseparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..m<$core.String, $core.String>(
        107376570, _omitFieldNames ? '' : 'responsetemplates',
        entryClassName: 'PutGatewayResponseRequest.ResponsetemplatesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aE<GatewayResponseType>(377935935, _omitFieldNames ? '' : 'responsetype',
        enumValues: GatewayResponseType.values)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutGatewayResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutGatewayResponseRequest copyWith(
          void Function(PutGatewayResponseRequest) updates) =>
      super.copyWith((message) => updates(message as PutGatewayResponseRequest))
          as PutGatewayResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutGatewayResponseRequest create() => PutGatewayResponseRequest._();
  @$core.override
  PutGatewayResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutGatewayResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutGatewayResponseRequest>(create);
  static PutGatewayResponseRequest? _defaultInstance;

  @$pb.TagNumber(64271839)
  $pb.PbMap<$core.String, $core.String> get responseparameters => $_getMap(0);

  @$pb.TagNumber(107376570)
  $pb.PbMap<$core.String, $core.String> get responsetemplates => $_getMap(1);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(2);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(2);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(377935935)
  GatewayResponseType get responsetype => $_getN(3);
  @$pb.TagNumber(377935935)
  set responsetype(GatewayResponseType value) => $_setField(377935935, value);
  @$pb.TagNumber(377935935)
  $core.bool hasResponsetype() => $_has(3);
  @$pb.TagNumber(377935935)
  void clearResponsetype() => $_clearField(377935935);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(4);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(4);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class PutIntegrationRequest extends $pb.GeneratedMessage {
  factory PutIntegrationRequest({
    $core.String? integrationtarget,
    $core.String? cachenamespace,
    TlsConfig? tlsconfig,
    $core.String? httpmethod,
    $core.String? credentials,
    IntegrationType? type,
    $core.String? passthroughbehavior,
    $core.String? resourceid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        requesttemplates,
    ConnectionType? connectiontype,
    $core.String? integrationhttpmethod,
    $core.int? timeoutinmillis,
    $core.String? restapiid,
    $core.String? uri,
    $core.String? connectionid,
    ResponseTransferMode? responsetransfermode,
    $core.Iterable<$core.String>? cachekeyparameters,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        requestparameters,
    ContentHandlingStrategy? contenthandling,
  }) {
    final result = create();
    if (integrationtarget != null) result.integrationtarget = integrationtarget;
    if (cachenamespace != null) result.cachenamespace = cachenamespace;
    if (tlsconfig != null) result.tlsconfig = tlsconfig;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (credentials != null) result.credentials = credentials;
    if (type != null) result.type = type;
    if (passthroughbehavior != null)
      result.passthroughbehavior = passthroughbehavior;
    if (resourceid != null) result.resourceid = resourceid;
    if (requesttemplates != null)
      result.requesttemplates.addEntries(requesttemplates);
    if (connectiontype != null) result.connectiontype = connectiontype;
    if (integrationhttpmethod != null)
      result.integrationhttpmethod = integrationhttpmethod;
    if (timeoutinmillis != null) result.timeoutinmillis = timeoutinmillis;
    if (restapiid != null) result.restapiid = restapiid;
    if (uri != null) result.uri = uri;
    if (connectionid != null) result.connectionid = connectionid;
    if (responsetransfermode != null)
      result.responsetransfermode = responsetransfermode;
    if (cachekeyparameters != null)
      result.cachekeyparameters.addAll(cachekeyparameters);
    if (requestparameters != null)
      result.requestparameters.addEntries(requestparameters);
    if (contenthandling != null) result.contenthandling = contenthandling;
    return result;
  }

  PutIntegrationRequest._();

  factory PutIntegrationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutIntegrationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutIntegrationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(17646705, _omitFieldNames ? '' : 'integrationtarget')
    ..aOS(85102753, _omitFieldNames ? '' : 'cachenamespace')
    ..aOM<TlsConfig>(108946693, _omitFieldNames ? '' : 'tlsconfig',
        subBuilder: TlsConfig.create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(150838226, _omitFieldNames ? '' : 'credentials')
    ..aE<IntegrationType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: IntegrationType.values)
    ..aOS(310796908, _omitFieldNames ? '' : 'passthroughbehavior')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..m<$core.String, $core.String>(
        333512166, _omitFieldNames ? '' : 'requesttemplates',
        entryClassName: 'PutIntegrationRequest.RequesttemplatesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aE<ConnectionType>(336253170, _omitFieldNames ? '' : 'connectiontype',
        enumValues: ConnectionType.values)
    ..aOS(355314073, _omitFieldNames ? '' : 'integrationhttpmethod')
    ..aI(378229126, _omitFieldNames ? '' : 'timeoutinmillis')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(395269118, _omitFieldNames ? '' : 'uri')
    ..aOS(450027965, _omitFieldNames ? '' : 'connectionid')
    ..aE<ResponseTransferMode>(
        458910787, _omitFieldNames ? '' : 'responsetransfermode',
        enumValues: ResponseTransferMode.values)
    ..pPS(481441313, _omitFieldNames ? '' : 'cachekeyparameters')
    ..m<$core.String, $core.String>(
        523499939, _omitFieldNames ? '' : 'requestparameters',
        entryClassName: 'PutIntegrationRequest.RequestparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aE<ContentHandlingStrategy>(
        533182832, _omitFieldNames ? '' : 'contenthandling',
        enumValues: ContentHandlingStrategy.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutIntegrationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutIntegrationRequest copyWith(
          void Function(PutIntegrationRequest) updates) =>
      super.copyWith((message) => updates(message as PutIntegrationRequest))
          as PutIntegrationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutIntegrationRequest create() => PutIntegrationRequest._();
  @$core.override
  PutIntegrationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutIntegrationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutIntegrationRequest>(create);
  static PutIntegrationRequest? _defaultInstance;

  @$pb.TagNumber(17646705)
  $core.String get integrationtarget => $_getSZ(0);
  @$pb.TagNumber(17646705)
  set integrationtarget($core.String value) => $_setString(0, value);
  @$pb.TagNumber(17646705)
  $core.bool hasIntegrationtarget() => $_has(0);
  @$pb.TagNumber(17646705)
  void clearIntegrationtarget() => $_clearField(17646705);

  @$pb.TagNumber(85102753)
  $core.String get cachenamespace => $_getSZ(1);
  @$pb.TagNumber(85102753)
  set cachenamespace($core.String value) => $_setString(1, value);
  @$pb.TagNumber(85102753)
  $core.bool hasCachenamespace() => $_has(1);
  @$pb.TagNumber(85102753)
  void clearCachenamespace() => $_clearField(85102753);

  @$pb.TagNumber(108946693)
  TlsConfig get tlsconfig => $_getN(2);
  @$pb.TagNumber(108946693)
  set tlsconfig(TlsConfig value) => $_setField(108946693, value);
  @$pb.TagNumber(108946693)
  $core.bool hasTlsconfig() => $_has(2);
  @$pb.TagNumber(108946693)
  void clearTlsconfig() => $_clearField(108946693);
  @$pb.TagNumber(108946693)
  TlsConfig ensureTlsconfig() => $_ensure(2);

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(3);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(3);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(150838226)
  $core.String get credentials => $_getSZ(4);
  @$pb.TagNumber(150838226)
  set credentials($core.String value) => $_setString(4, value);
  @$pb.TagNumber(150838226)
  $core.bool hasCredentials() => $_has(4);
  @$pb.TagNumber(150838226)
  void clearCredentials() => $_clearField(150838226);

  @$pb.TagNumber(287830350)
  IntegrationType get type => $_getN(5);
  @$pb.TagNumber(287830350)
  set type(IntegrationType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(5);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(310796908)
  $core.String get passthroughbehavior => $_getSZ(6);
  @$pb.TagNumber(310796908)
  set passthroughbehavior($core.String value) => $_setString(6, value);
  @$pb.TagNumber(310796908)
  $core.bool hasPassthroughbehavior() => $_has(6);
  @$pb.TagNumber(310796908)
  void clearPassthroughbehavior() => $_clearField(310796908);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(7);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(7, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(7);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(333512166)
  $pb.PbMap<$core.String, $core.String> get requesttemplates => $_getMap(8);

  @$pb.TagNumber(336253170)
  ConnectionType get connectiontype => $_getN(9);
  @$pb.TagNumber(336253170)
  set connectiontype(ConnectionType value) => $_setField(336253170, value);
  @$pb.TagNumber(336253170)
  $core.bool hasConnectiontype() => $_has(9);
  @$pb.TagNumber(336253170)
  void clearConnectiontype() => $_clearField(336253170);

  @$pb.TagNumber(355314073)
  $core.String get integrationhttpmethod => $_getSZ(10);
  @$pb.TagNumber(355314073)
  set integrationhttpmethod($core.String value) => $_setString(10, value);
  @$pb.TagNumber(355314073)
  $core.bool hasIntegrationhttpmethod() => $_has(10);
  @$pb.TagNumber(355314073)
  void clearIntegrationhttpmethod() => $_clearField(355314073);

  @$pb.TagNumber(378229126)
  $core.int get timeoutinmillis => $_getIZ(11);
  @$pb.TagNumber(378229126)
  set timeoutinmillis($core.int value) => $_setSignedInt32(11, value);
  @$pb.TagNumber(378229126)
  $core.bool hasTimeoutinmillis() => $_has(11);
  @$pb.TagNumber(378229126)
  void clearTimeoutinmillis() => $_clearField(378229126);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(12);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(12, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(12);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(395269118)
  $core.String get uri => $_getSZ(13);
  @$pb.TagNumber(395269118)
  set uri($core.String value) => $_setString(13, value);
  @$pb.TagNumber(395269118)
  $core.bool hasUri() => $_has(13);
  @$pb.TagNumber(395269118)
  void clearUri() => $_clearField(395269118);

  @$pb.TagNumber(450027965)
  $core.String get connectionid => $_getSZ(14);
  @$pb.TagNumber(450027965)
  set connectionid($core.String value) => $_setString(14, value);
  @$pb.TagNumber(450027965)
  $core.bool hasConnectionid() => $_has(14);
  @$pb.TagNumber(450027965)
  void clearConnectionid() => $_clearField(450027965);

  @$pb.TagNumber(458910787)
  ResponseTransferMode get responsetransfermode => $_getN(15);
  @$pb.TagNumber(458910787)
  set responsetransfermode(ResponseTransferMode value) =>
      $_setField(458910787, value);
  @$pb.TagNumber(458910787)
  $core.bool hasResponsetransfermode() => $_has(15);
  @$pb.TagNumber(458910787)
  void clearResponsetransfermode() => $_clearField(458910787);

  @$pb.TagNumber(481441313)
  $pb.PbList<$core.String> get cachekeyparameters => $_getList(16);

  @$pb.TagNumber(523499939)
  $pb.PbMap<$core.String, $core.String> get requestparameters => $_getMap(17);

  @$pb.TagNumber(533182832)
  ContentHandlingStrategy get contenthandling => $_getN(18);
  @$pb.TagNumber(533182832)
  set contenthandling(ContentHandlingStrategy value) =>
      $_setField(533182832, value);
  @$pb.TagNumber(533182832)
  $core.bool hasContenthandling() => $_has(18);
  @$pb.TagNumber(533182832)
  void clearContenthandling() => $_clearField(533182832);
}

class PutIntegrationResponseRequest extends $pb.GeneratedMessage {
  factory PutIntegrationResponseRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        responseparameters,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        responsetemplates,
    $core.String? httpmethod,
    $core.String? statuscode,
    $core.String? resourceid,
    $core.String? restapiid,
    $core.String? selectionpattern,
    ContentHandlingStrategy? contenthandling,
  }) {
    final result = create();
    if (responseparameters != null)
      result.responseparameters.addEntries(responseparameters);
    if (responsetemplates != null)
      result.responsetemplates.addEntries(responsetemplates);
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (statuscode != null) result.statuscode = statuscode;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    if (selectionpattern != null) result.selectionpattern = selectionpattern;
    if (contenthandling != null) result.contenthandling = contenthandling;
    return result;
  }

  PutIntegrationResponseRequest._();

  factory PutIntegrationResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutIntegrationResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutIntegrationResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        64271839, _omitFieldNames ? '' : 'responseparameters',
        entryClassName: 'PutIntegrationResponseRequest.ResponseparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..m<$core.String, $core.String>(
        107376570, _omitFieldNames ? '' : 'responsetemplates',
        entryClassName: 'PutIntegrationResponseRequest.ResponsetemplatesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(470634042, _omitFieldNames ? '' : 'selectionpattern')
    ..aE<ContentHandlingStrategy>(
        533182832, _omitFieldNames ? '' : 'contenthandling',
        enumValues: ContentHandlingStrategy.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutIntegrationResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutIntegrationResponseRequest copyWith(
          void Function(PutIntegrationResponseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PutIntegrationResponseRequest))
          as PutIntegrationResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutIntegrationResponseRequest create() =>
      PutIntegrationResponseRequest._();
  @$core.override
  PutIntegrationResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutIntegrationResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutIntegrationResponseRequest>(create);
  static PutIntegrationResponseRequest? _defaultInstance;

  @$pb.TagNumber(64271839)
  $pb.PbMap<$core.String, $core.String> get responseparameters => $_getMap(0);

  @$pb.TagNumber(107376570)
  $pb.PbMap<$core.String, $core.String> get responsetemplates => $_getMap(1);

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(2);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(2);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(3);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(3, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(3);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(4);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(4);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(5);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(5);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(470634042)
  $core.String get selectionpattern => $_getSZ(6);
  @$pb.TagNumber(470634042)
  set selectionpattern($core.String value) => $_setString(6, value);
  @$pb.TagNumber(470634042)
  $core.bool hasSelectionpattern() => $_has(6);
  @$pb.TagNumber(470634042)
  void clearSelectionpattern() => $_clearField(470634042);

  @$pb.TagNumber(533182832)
  ContentHandlingStrategy get contenthandling => $_getN(7);
  @$pb.TagNumber(533182832)
  set contenthandling(ContentHandlingStrategy value) =>
      $_setField(533182832, value);
  @$pb.TagNumber(533182832)
  $core.bool hasContenthandling() => $_has(7);
  @$pb.TagNumber(533182832)
  void clearContenthandling() => $_clearField(533182832);
}

class PutMethodRequest extends $pb.GeneratedMessage {
  factory PutMethodRequest({
    $core.String? authorizerid,
    $core.String? httpmethod,
    $core.String? operationname,
    $core.String? resourceid,
    $core.String? restapiid,
    $core.String? authorizationtype,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? requestmodels,
    $core.Iterable<$core.String>? authorizationscopes,
    $core.bool? apikeyrequired,
    $core.String? requestvalidatorid,
    $core.Iterable<$core.MapEntry<$core.String, $core.bool>>? requestparameters,
  }) {
    final result = create();
    if (authorizerid != null) result.authorizerid = authorizerid;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (operationname != null) result.operationname = operationname;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    if (authorizationtype != null) result.authorizationtype = authorizationtype;
    if (requestmodels != null) result.requestmodels.addEntries(requestmodels);
    if (authorizationscopes != null)
      result.authorizationscopes.addAll(authorizationscopes);
    if (apikeyrequired != null) result.apikeyrequired = apikeyrequired;
    if (requestvalidatorid != null)
      result.requestvalidatorid = requestvalidatorid;
    if (requestparameters != null)
      result.requestparameters.addEntries(requestparameters);
    return result;
  }

  PutMethodRequest._();

  factory PutMethodRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutMethodRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutMethodRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(111773148, _omitFieldNames ? '' : 'authorizerid')
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(178909574, _omitFieldNames ? '' : 'operationname')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(387986911, _omitFieldNames ? '' : 'authorizationtype')
    ..m<$core.String, $core.String>(
        397252853, _omitFieldNames ? '' : 'requestmodels',
        entryClassName: 'PutMethodRequest.RequestmodelsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..pPS(423149932, _omitFieldNames ? '' : 'authorizationscopes')
    ..aOB(435360152, _omitFieldNames ? '' : 'apikeyrequired')
    ..aOS(517546134, _omitFieldNames ? '' : 'requestvalidatorid')
    ..m<$core.String, $core.bool>(
        523499939, _omitFieldNames ? '' : 'requestparameters',
        entryClassName: 'PutMethodRequest.RequestparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OB,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMethodRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMethodRequest copyWith(void Function(PutMethodRequest) updates) =>
      super.copyWith((message) => updates(message as PutMethodRequest))
          as PutMethodRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutMethodRequest create() => PutMethodRequest._();
  @$core.override
  PutMethodRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutMethodRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutMethodRequest>(create);
  static PutMethodRequest? _defaultInstance;

  @$pb.TagNumber(111773148)
  $core.String get authorizerid => $_getSZ(0);
  @$pb.TagNumber(111773148)
  set authorizerid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(111773148)
  $core.bool hasAuthorizerid() => $_has(0);
  @$pb.TagNumber(111773148)
  void clearAuthorizerid() => $_clearField(111773148);

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(1);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(1);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(178909574)
  $core.String get operationname => $_getSZ(2);
  @$pb.TagNumber(178909574)
  set operationname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(178909574)
  $core.bool hasOperationname() => $_has(2);
  @$pb.TagNumber(178909574)
  void clearOperationname() => $_clearField(178909574);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(3);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(3);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(4);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(4);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(387986911)
  $core.String get authorizationtype => $_getSZ(5);
  @$pb.TagNumber(387986911)
  set authorizationtype($core.String value) => $_setString(5, value);
  @$pb.TagNumber(387986911)
  $core.bool hasAuthorizationtype() => $_has(5);
  @$pb.TagNumber(387986911)
  void clearAuthorizationtype() => $_clearField(387986911);

  @$pb.TagNumber(397252853)
  $pb.PbMap<$core.String, $core.String> get requestmodels => $_getMap(6);

  @$pb.TagNumber(423149932)
  $pb.PbList<$core.String> get authorizationscopes => $_getList(7);

  @$pb.TagNumber(435360152)
  $core.bool get apikeyrequired => $_getBF(8);
  @$pb.TagNumber(435360152)
  set apikeyrequired($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(435360152)
  $core.bool hasApikeyrequired() => $_has(8);
  @$pb.TagNumber(435360152)
  void clearApikeyrequired() => $_clearField(435360152);

  @$pb.TagNumber(517546134)
  $core.String get requestvalidatorid => $_getSZ(9);
  @$pb.TagNumber(517546134)
  set requestvalidatorid($core.String value) => $_setString(9, value);
  @$pb.TagNumber(517546134)
  $core.bool hasRequestvalidatorid() => $_has(9);
  @$pb.TagNumber(517546134)
  void clearRequestvalidatorid() => $_clearField(517546134);

  @$pb.TagNumber(523499939)
  $pb.PbMap<$core.String, $core.bool> get requestparameters => $_getMap(10);
}

class PutMethodResponseRequest extends $pb.GeneratedMessage {
  factory PutMethodResponseRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.bool>>?
        responseparameters,
    $core.String? httpmethod,
    $core.String? statuscode,
    $core.String? resourceid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? responsemodels,
    $core.String? restapiid,
  }) {
    final result = create();
    if (responseparameters != null)
      result.responseparameters.addEntries(responseparameters);
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (statuscode != null) result.statuscode = statuscode;
    if (resourceid != null) result.resourceid = resourceid;
    if (responsemodels != null)
      result.responsemodels.addEntries(responsemodels);
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  PutMethodResponseRequest._();

  factory PutMethodResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutMethodResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutMethodResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.bool>(
        64271839, _omitFieldNames ? '' : 'responseparameters',
        entryClassName: 'PutMethodResponseRequest.ResponseparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OB,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..m<$core.String, $core.String>(
        356574313, _omitFieldNames ? '' : 'responsemodels',
        entryClassName: 'PutMethodResponseRequest.ResponsemodelsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMethodResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMethodResponseRequest copyWith(
          void Function(PutMethodResponseRequest) updates) =>
      super.copyWith((message) => updates(message as PutMethodResponseRequest))
          as PutMethodResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutMethodResponseRequest create() => PutMethodResponseRequest._();
  @$core.override
  PutMethodResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutMethodResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutMethodResponseRequest>(create);
  static PutMethodResponseRequest? _defaultInstance;

  @$pb.TagNumber(64271839)
  $pb.PbMap<$core.String, $core.bool> get responseparameters => $_getMap(0);

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(1);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(1);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(2);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(2);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(3);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(3);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(356574313)
  $pb.PbMap<$core.String, $core.String> get responsemodels => $_getMap(4);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(5);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(5);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class PutRestApiRequest extends $pb.GeneratedMessage {
  factory PutRestApiRequest({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? parameters,
    PutMode? mode,
    $core.String? restapiid,
    $core.bool? failonwarnings,
    $core.List<$core.int>? body,
  }) {
    final result = create();
    if (parameters != null) result.parameters.addEntries(parameters);
    if (mode != null) result.mode = mode;
    if (restapiid != null) result.restapiid = restapiid;
    if (failonwarnings != null) result.failonwarnings = failonwarnings;
    if (body != null) result.body = body;
    return result;
  }

  PutRestApiRequest._();

  factory PutRestApiRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRestApiRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRestApiRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        145043162, _omitFieldNames ? '' : 'parameters',
        entryClassName: 'PutRestApiRequest.ParametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aE<PutMode>(208592915, _omitFieldNames ? '' : 'mode',
        enumValues: PutMode.values)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOB(434363958, _omitFieldNames ? '' : 'failonwarnings')
    ..a<$core.List<$core.int>>(
        464157046, _omitFieldNames ? '' : 'body', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRestApiRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRestApiRequest copyWith(void Function(PutRestApiRequest) updates) =>
      super.copyWith((message) => updates(message as PutRestApiRequest))
          as PutRestApiRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRestApiRequest create() => PutRestApiRequest._();
  @$core.override
  PutRestApiRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRestApiRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRestApiRequest>(create);
  static PutRestApiRequest? _defaultInstance;

  @$pb.TagNumber(145043162)
  $pb.PbMap<$core.String, $core.String> get parameters => $_getMap(0);

  @$pb.TagNumber(208592915)
  PutMode get mode => $_getN(1);
  @$pb.TagNumber(208592915)
  set mode(PutMode value) => $_setField(208592915, value);
  @$pb.TagNumber(208592915)
  $core.bool hasMode() => $_has(1);
  @$pb.TagNumber(208592915)
  void clearMode() => $_clearField(208592915);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(434363958)
  $core.bool get failonwarnings => $_getBF(3);
  @$pb.TagNumber(434363958)
  set failonwarnings($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(434363958)
  $core.bool hasFailonwarnings() => $_has(3);
  @$pb.TagNumber(434363958)
  void clearFailonwarnings() => $_clearField(434363958);

  @$pb.TagNumber(464157046)
  $core.List<$core.int> get body => $_getN(4);
  @$pb.TagNumber(464157046)
  set body($core.List<$core.int> value) => $_setBytes(4, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(4);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);
}

class QuotaSettings extends $pb.GeneratedMessage {
  factory QuotaSettings({
    $core.int? limit,
    $core.int? offset,
    QuotaPeriodType? period,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (offset != null) result.offset = offset;
    if (period != null) result.period = period;
    return result;
  }

  QuotaSettings._();

  factory QuotaSettings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QuotaSettings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QuotaSettings',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(316332341, _omitFieldNames ? '' : 'limit')
    ..aI(348705739, _omitFieldNames ? '' : 'offset')
    ..aE<QuotaPeriodType>(432621317, _omitFieldNames ? '' : 'period',
        enumValues: QuotaPeriodType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaSettings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QuotaSettings copyWith(void Function(QuotaSettings) updates) =>
      super.copyWith((message) => updates(message as QuotaSettings))
          as QuotaSettings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QuotaSettings create() => QuotaSettings._();
  @$core.override
  QuotaSettings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QuotaSettings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QuotaSettings>(create);
  static QuotaSettings? _defaultInstance;

  @$pb.TagNumber(316332341)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(316332341)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(316332341)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(316332341)
  void clearLimit() => $_clearField(316332341);

  @$pb.TagNumber(348705739)
  $core.int get offset => $_getIZ(1);
  @$pb.TagNumber(348705739)
  set offset($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(348705739)
  $core.bool hasOffset() => $_has(1);
  @$pb.TagNumber(348705739)
  void clearOffset() => $_clearField(348705739);

  @$pb.TagNumber(432621317)
  QuotaPeriodType get period => $_getN(2);
  @$pb.TagNumber(432621317)
  set period(QuotaPeriodType value) => $_setField(432621317, value);
  @$pb.TagNumber(432621317)
  $core.bool hasPeriod() => $_has(2);
  @$pb.TagNumber(432621317)
  void clearPeriod() => $_clearField(432621317);
}

class RejectDomainNameAccessAssociationRequest extends $pb.GeneratedMessage {
  factory RejectDomainNameAccessAssociationRequest({
    $core.String? domainnamearn,
    $core.String? domainnameaccessassociationarn,
  }) {
    final result = create();
    if (domainnamearn != null) result.domainnamearn = domainnamearn;
    if (domainnameaccessassociationarn != null)
      result.domainnameaccessassociationarn = domainnameaccessassociationarn;
    return result;
  }

  RejectDomainNameAccessAssociationRequest._();

  factory RejectDomainNameAccessAssociationRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RejectDomainNameAccessAssociationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RejectDomainNameAccessAssociationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(244019094, _omitFieldNames ? '' : 'domainnamearn')
    ..aOS(281017927, _omitFieldNames ? '' : 'domainnameaccessassociationarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RejectDomainNameAccessAssociationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RejectDomainNameAccessAssociationRequest copyWith(
          void Function(RejectDomainNameAccessAssociationRequest) updates) =>
      super.copyWith((message) =>
              updates(message as RejectDomainNameAccessAssociationRequest))
          as RejectDomainNameAccessAssociationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RejectDomainNameAccessAssociationRequest create() =>
      RejectDomainNameAccessAssociationRequest._();
  @$core.override
  RejectDomainNameAccessAssociationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RejectDomainNameAccessAssociationRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          RejectDomainNameAccessAssociationRequest>(create);
  static RejectDomainNameAccessAssociationRequest? _defaultInstance;

  @$pb.TagNumber(244019094)
  $core.String get domainnamearn => $_getSZ(0);
  @$pb.TagNumber(244019094)
  set domainnamearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(244019094)
  $core.bool hasDomainnamearn() => $_has(0);
  @$pb.TagNumber(244019094)
  void clearDomainnamearn() => $_clearField(244019094);

  @$pb.TagNumber(281017927)
  $core.String get domainnameaccessassociationarn => $_getSZ(1);
  @$pb.TagNumber(281017927)
  set domainnameaccessassociationarn($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(281017927)
  $core.bool hasDomainnameaccessassociationarn() => $_has(1);
  @$pb.TagNumber(281017927)
  void clearDomainnameaccessassociationarn() => $_clearField(281017927);
}

class RequestValidator extends $pb.GeneratedMessage {
  factory RequestValidator({
    $core.String? name,
    $core.String? id,
    $core.bool? validaterequestbody,
    $core.bool? validaterequestparameters,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (validaterequestbody != null)
      result.validaterequestbody = validaterequestbody;
    if (validaterequestparameters != null)
      result.validaterequestparameters = validaterequestparameters;
    return result;
  }

  RequestValidator._();

  factory RequestValidator.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestValidator.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestValidator',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..aOB(397505841, _omitFieldNames ? '' : 'validaterequestbody')
    ..aOB(464035801, _omitFieldNames ? '' : 'validaterequestparameters')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestValidator clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestValidator copyWith(void Function(RequestValidator) updates) =>
      super.copyWith((message) => updates(message as RequestValidator))
          as RequestValidator;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestValidator create() => RequestValidator._();
  @$core.override
  RequestValidator createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestValidator getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestValidator>(create);
  static RequestValidator? _defaultInstance;

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);

  @$pb.TagNumber(397505841)
  $core.bool get validaterequestbody => $_getBF(2);
  @$pb.TagNumber(397505841)
  set validaterequestbody($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(397505841)
  $core.bool hasValidaterequestbody() => $_has(2);
  @$pb.TagNumber(397505841)
  void clearValidaterequestbody() => $_clearField(397505841);

  @$pb.TagNumber(464035801)
  $core.bool get validaterequestparameters => $_getBF(3);
  @$pb.TagNumber(464035801)
  set validaterequestparameters($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(464035801)
  $core.bool hasValidaterequestparameters() => $_has(3);
  @$pb.TagNumber(464035801)
  void clearValidaterequestparameters() => $_clearField(464035801);
}

class RequestValidators extends $pb.GeneratedMessage {
  factory RequestValidators({
    $core.String? position,
    $core.Iterable<RequestValidator>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  RequestValidators._();

  factory RequestValidators.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestValidators.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestValidators',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<RequestValidator>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: RequestValidator.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestValidators clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestValidators copyWith(void Function(RequestValidators) updates) =>
      super.copyWith((message) => updates(message as RequestValidators))
          as RequestValidators;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestValidators create() => RequestValidators._();
  @$core.override
  RequestValidators createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestValidators getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestValidators>(create);
  static RequestValidators? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<RequestValidator> get items => $_getList(1);
}

class Resource extends $pb.GeneratedMessage {
  factory Resource({
    $core.String? path,
    $core.String? parentid,
    $core.Iterable<$core.MapEntry<$core.String, Method>>? resourcemethods,
    $core.String? id,
    $core.String? pathpart,
  }) {
    final result = create();
    if (path != null) result.path = path;
    if (parentid != null) result.parentid = parentid;
    if (resourcemethods != null)
      result.resourcemethods.addEntries(resourcemethods);
    if (id != null) result.id = id;
    if (pathpart != null) result.pathpart = pathpart;
    return result;
  }

  Resource._();

  factory Resource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Resource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Resource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(75975991, _omitFieldNames ? '' : 'path')
    ..aOS(106050857, _omitFieldNames ? '' : 'parentid')
    ..m<$core.String, Method>(
        307700458, _omitFieldNames ? '' : 'resourcemethods',
        entryClassName: 'Resource.ResourcemethodsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: Method.create,
        valueDefaultOrMaker: Method.getDefault,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..aOS(487915984, _omitFieldNames ? '' : 'pathpart')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Resource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Resource copyWith(void Function(Resource) updates) =>
      super.copyWith((message) => updates(message as Resource)) as Resource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Resource create() => Resource._();
  @$core.override
  Resource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Resource getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Resource>(create);
  static Resource? _defaultInstance;

  @$pb.TagNumber(75975991)
  $core.String get path => $_getSZ(0);
  @$pb.TagNumber(75975991)
  set path($core.String value) => $_setString(0, value);
  @$pb.TagNumber(75975991)
  $core.bool hasPath() => $_has(0);
  @$pb.TagNumber(75975991)
  void clearPath() => $_clearField(75975991);

  @$pb.TagNumber(106050857)
  $core.String get parentid => $_getSZ(1);
  @$pb.TagNumber(106050857)
  set parentid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(106050857)
  $core.bool hasParentid() => $_has(1);
  @$pb.TagNumber(106050857)
  void clearParentid() => $_clearField(106050857);

  @$pb.TagNumber(307700458)
  $pb.PbMap<$core.String, Method> get resourcemethods => $_getMap(2);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);

  @$pb.TagNumber(487915984)
  $core.String get pathpart => $_getSZ(4);
  @$pb.TagNumber(487915984)
  set pathpart($core.String value) => $_setString(4, value);
  @$pb.TagNumber(487915984)
  $core.bool hasPathpart() => $_has(4);
  @$pb.TagNumber(487915984)
  void clearPathpart() => $_clearField(487915984);
}

class Resources extends $pb.GeneratedMessage {
  factory Resources({
    $core.String? position,
    $core.Iterable<Resource>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  Resources._();

  factory Resources.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Resources.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Resources',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<Resource>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: Resource.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Resources clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Resources copyWith(void Function(Resources) updates) =>
      super.copyWith((message) => updates(message as Resources)) as Resources;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Resources create() => Resources._();
  @$core.override
  Resources createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Resources getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Resources>(create);
  static Resources? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<Resource> get items => $_getList(1);
}

class RestApi extends $pb.GeneratedMessage {
  factory RestApi({
    $core.String? createddate,
    $core.String? version,
    ApiKeySourceType? apikeysource,
    $core.bool? disableexecuteapiendpoint,
    $core.Iterable<$core.String>? warnings,
    ApiStatus? apistatus,
    $core.String? name,
    $core.String? policy,
    $core.int? minimumcompressionsize,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    $core.String? apistatusmessage,
    EndpointAccessMode? endpointaccessmode,
    $core.String? rootresourceid,
    $core.String? id,
    $core.Iterable<$core.String>? binarymediatypes,
    EndpointConfiguration? endpointconfiguration,
    SecurityPolicy? securitypolicy,
  }) {
    final result = create();
    if (createddate != null) result.createddate = createddate;
    if (version != null) result.version = version;
    if (apikeysource != null) result.apikeysource = apikeysource;
    if (disableexecuteapiendpoint != null)
      result.disableexecuteapiendpoint = disableexecuteapiendpoint;
    if (warnings != null) result.warnings.addAll(warnings);
    if (apistatus != null) result.apistatus = apistatus;
    if (name != null) result.name = name;
    if (policy != null) result.policy = policy;
    if (minimumcompressionsize != null)
      result.minimumcompressionsize = minimumcompressionsize;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (apistatusmessage != null) result.apistatusmessage = apistatusmessage;
    if (endpointaccessmode != null)
      result.endpointaccessmode = endpointaccessmode;
    if (rootresourceid != null) result.rootresourceid = rootresourceid;
    if (id != null) result.id = id;
    if (binarymediatypes != null)
      result.binarymediatypes.addAll(binarymediatypes);
    if (endpointconfiguration != null)
      result.endpointconfiguration = endpointconfiguration;
    if (securitypolicy != null) result.securitypolicy = securitypolicy;
    return result;
  }

  RestApi._();

  factory RestApi.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestApi.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestApi',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(53061200, _omitFieldNames ? '' : 'createddate')
    ..aOS(108113560, _omitFieldNames ? '' : 'version')
    ..aE<ApiKeySourceType>(108531220, _omitFieldNames ? '' : 'apikeysource',
        enumValues: ApiKeySourceType.values)
    ..aOB(148140696, _omitFieldNames ? '' : 'disableexecuteapiendpoint')
    ..pPS(185617301, _omitFieldNames ? '' : 'warnings')
    ..aE<ApiStatus>(200568018, _omitFieldNames ? '' : 'apistatus',
        enumValues: ApiStatus.values)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(247528064, _omitFieldNames ? '' : 'policy')
    ..aI(254902719, _omitFieldNames ? '' : 'minimumcompressionsize')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'RestApi.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(353995209, _omitFieldNames ? '' : 'apistatusmessage')
    ..aE<EndpointAccessMode>(
        356705630, _omitFieldNames ? '' : 'endpointaccessmode',
        enumValues: EndpointAccessMode.values)
    ..aOS(360157585, _omitFieldNames ? '' : 'rootresourceid')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..pPS(406416146, _omitFieldNames ? '' : 'binarymediatypes')
    ..aOM<EndpointConfiguration>(
        487543735, _omitFieldNames ? '' : 'endpointconfiguration',
        subBuilder: EndpointConfiguration.create)
    ..aE<SecurityPolicy>(491792990, _omitFieldNames ? '' : 'securitypolicy',
        enumValues: SecurityPolicy.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestApi clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestApi copyWith(void Function(RestApi) updates) =>
      super.copyWith((message) => updates(message as RestApi)) as RestApi;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestApi create() => RestApi._();
  @$core.override
  RestApi createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestApi getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RestApi>(create);
  static RestApi? _defaultInstance;

  @$pb.TagNumber(53061200)
  $core.String get createddate => $_getSZ(0);
  @$pb.TagNumber(53061200)
  set createddate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(53061200)
  $core.bool hasCreateddate() => $_has(0);
  @$pb.TagNumber(53061200)
  void clearCreateddate() => $_clearField(53061200);

  @$pb.TagNumber(108113560)
  $core.String get version => $_getSZ(1);
  @$pb.TagNumber(108113560)
  set version($core.String value) => $_setString(1, value);
  @$pb.TagNumber(108113560)
  $core.bool hasVersion() => $_has(1);
  @$pb.TagNumber(108113560)
  void clearVersion() => $_clearField(108113560);

  @$pb.TagNumber(108531220)
  ApiKeySourceType get apikeysource => $_getN(2);
  @$pb.TagNumber(108531220)
  set apikeysource(ApiKeySourceType value) => $_setField(108531220, value);
  @$pb.TagNumber(108531220)
  $core.bool hasApikeysource() => $_has(2);
  @$pb.TagNumber(108531220)
  void clearApikeysource() => $_clearField(108531220);

  @$pb.TagNumber(148140696)
  $core.bool get disableexecuteapiendpoint => $_getBF(3);
  @$pb.TagNumber(148140696)
  set disableexecuteapiendpoint($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(148140696)
  $core.bool hasDisableexecuteapiendpoint() => $_has(3);
  @$pb.TagNumber(148140696)
  void clearDisableexecuteapiendpoint() => $_clearField(148140696);

  @$pb.TagNumber(185617301)
  $pb.PbList<$core.String> get warnings => $_getList(4);

  @$pb.TagNumber(200568018)
  ApiStatus get apistatus => $_getN(5);
  @$pb.TagNumber(200568018)
  set apistatus(ApiStatus value) => $_setField(200568018, value);
  @$pb.TagNumber(200568018)
  $core.bool hasApistatus() => $_has(5);
  @$pb.TagNumber(200568018)
  void clearApistatus() => $_clearField(200568018);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(247528064)
  $core.String get policy => $_getSZ(7);
  @$pb.TagNumber(247528064)
  set policy($core.String value) => $_setString(7, value);
  @$pb.TagNumber(247528064)
  $core.bool hasPolicy() => $_has(7);
  @$pb.TagNumber(247528064)
  void clearPolicy() => $_clearField(247528064);

  @$pb.TagNumber(254902719)
  $core.int get minimumcompressionsize => $_getIZ(8);
  @$pb.TagNumber(254902719)
  set minimumcompressionsize($core.int value) => $_setSignedInt32(8, value);
  @$pb.TagNumber(254902719)
  $core.bool hasMinimumcompressionsize() => $_has(8);
  @$pb.TagNumber(254902719)
  void clearMinimumcompressionsize() => $_clearField(254902719);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(9);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(10);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(10, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(10);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(353995209)
  $core.String get apistatusmessage => $_getSZ(11);
  @$pb.TagNumber(353995209)
  set apistatusmessage($core.String value) => $_setString(11, value);
  @$pb.TagNumber(353995209)
  $core.bool hasApistatusmessage() => $_has(11);
  @$pb.TagNumber(353995209)
  void clearApistatusmessage() => $_clearField(353995209);

  @$pb.TagNumber(356705630)
  EndpointAccessMode get endpointaccessmode => $_getN(12);
  @$pb.TagNumber(356705630)
  set endpointaccessmode(EndpointAccessMode value) =>
      $_setField(356705630, value);
  @$pb.TagNumber(356705630)
  $core.bool hasEndpointaccessmode() => $_has(12);
  @$pb.TagNumber(356705630)
  void clearEndpointaccessmode() => $_clearField(356705630);

  @$pb.TagNumber(360157585)
  $core.String get rootresourceid => $_getSZ(13);
  @$pb.TagNumber(360157585)
  set rootresourceid($core.String value) => $_setString(13, value);
  @$pb.TagNumber(360157585)
  $core.bool hasRootresourceid() => $_has(13);
  @$pb.TagNumber(360157585)
  void clearRootresourceid() => $_clearField(360157585);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(14);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(14, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(14);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);

  @$pb.TagNumber(406416146)
  $pb.PbList<$core.String> get binarymediatypes => $_getList(15);

  @$pb.TagNumber(487543735)
  EndpointConfiguration get endpointconfiguration => $_getN(16);
  @$pb.TagNumber(487543735)
  set endpointconfiguration(EndpointConfiguration value) =>
      $_setField(487543735, value);
  @$pb.TagNumber(487543735)
  $core.bool hasEndpointconfiguration() => $_has(16);
  @$pb.TagNumber(487543735)
  void clearEndpointconfiguration() => $_clearField(487543735);
  @$pb.TagNumber(487543735)
  EndpointConfiguration ensureEndpointconfiguration() => $_ensure(16);

  @$pb.TagNumber(491792990)
  SecurityPolicy get securitypolicy => $_getN(17);
  @$pb.TagNumber(491792990)
  set securitypolicy(SecurityPolicy value) => $_setField(491792990, value);
  @$pb.TagNumber(491792990)
  $core.bool hasSecuritypolicy() => $_has(17);
  @$pb.TagNumber(491792990)
  void clearSecuritypolicy() => $_clearField(491792990);
}

class RestApis extends $pb.GeneratedMessage {
  factory RestApis({
    $core.String? position,
    $core.Iterable<RestApi>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  RestApis._();

  factory RestApis.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestApis.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestApis',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<RestApi>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: RestApi.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestApis clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestApis copyWith(void Function(RestApis) updates) =>
      super.copyWith((message) => updates(message as RestApis)) as RestApis;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestApis create() => RestApis._();
  @$core.override
  RestApis createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestApis getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RestApis>(create);
  static RestApis? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<RestApi> get items => $_getList(1);
}

class SdkConfigurationProperty extends $pb.GeneratedMessage {
  factory SdkConfigurationProperty({
    $core.String? friendlyname,
    $core.bool? required,
    $core.String? name,
    $core.String? description,
    $core.String? defaultvalue,
  }) {
    final result = create();
    if (friendlyname != null) result.friendlyname = friendlyname;
    if (required != null) result.required = required;
    if (name != null) result.name = name;
    if (description != null) result.description = description;
    if (defaultvalue != null) result.defaultvalue = defaultvalue;
    return result;
  }

  SdkConfigurationProperty._();

  factory SdkConfigurationProperty.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SdkConfigurationProperty.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SdkConfigurationProperty',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(39239102, _omitFieldNames ? '' : 'friendlyname')
    ..aOB(76318241, _omitFieldNames ? '' : 'required')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(403858624, _omitFieldNames ? '' : 'defaultvalue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SdkConfigurationProperty clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SdkConfigurationProperty copyWith(
          void Function(SdkConfigurationProperty) updates) =>
      super.copyWith((message) => updates(message as SdkConfigurationProperty))
          as SdkConfigurationProperty;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SdkConfigurationProperty create() => SdkConfigurationProperty._();
  @$core.override
  SdkConfigurationProperty createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SdkConfigurationProperty getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SdkConfigurationProperty>(create);
  static SdkConfigurationProperty? _defaultInstance;

  @$pb.TagNumber(39239102)
  $core.String get friendlyname => $_getSZ(0);
  @$pb.TagNumber(39239102)
  set friendlyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39239102)
  $core.bool hasFriendlyname() => $_has(0);
  @$pb.TagNumber(39239102)
  void clearFriendlyname() => $_clearField(39239102);

  @$pb.TagNumber(76318241)
  $core.bool get required => $_getBF(1);
  @$pb.TagNumber(76318241)
  set required($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(76318241)
  $core.bool hasRequired() => $_has(1);
  @$pb.TagNumber(76318241)
  void clearRequired() => $_clearField(76318241);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(403858624)
  $core.String get defaultvalue => $_getSZ(4);
  @$pb.TagNumber(403858624)
  set defaultvalue($core.String value) => $_setString(4, value);
  @$pb.TagNumber(403858624)
  $core.bool hasDefaultvalue() => $_has(4);
  @$pb.TagNumber(403858624)
  void clearDefaultvalue() => $_clearField(403858624);
}

class SdkResponse extends $pb.GeneratedMessage {
  factory SdkResponse({
    $core.String? contenttype,
    $core.String? contentdisposition,
    $core.List<$core.int>? body,
  }) {
    final result = create();
    if (contenttype != null) result.contenttype = contenttype;
    if (contentdisposition != null)
      result.contentdisposition = contentdisposition;
    if (body != null) result.body = body;
    return result;
  }

  SdkResponse._();

  factory SdkResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SdkResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SdkResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(281764659, _omitFieldNames ? '' : 'contenttype')
    ..aOS(375146466, _omitFieldNames ? '' : 'contentdisposition')
    ..a<$core.List<$core.int>>(
        464157046, _omitFieldNames ? '' : 'body', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SdkResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SdkResponse copyWith(void Function(SdkResponse) updates) =>
      super.copyWith((message) => updates(message as SdkResponse))
          as SdkResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SdkResponse create() => SdkResponse._();
  @$core.override
  SdkResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SdkResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SdkResponse>(create);
  static SdkResponse? _defaultInstance;

  @$pb.TagNumber(281764659)
  $core.String get contenttype => $_getSZ(0);
  @$pb.TagNumber(281764659)
  set contenttype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(281764659)
  $core.bool hasContenttype() => $_has(0);
  @$pb.TagNumber(281764659)
  void clearContenttype() => $_clearField(281764659);

  @$pb.TagNumber(375146466)
  $core.String get contentdisposition => $_getSZ(1);
  @$pb.TagNumber(375146466)
  set contentdisposition($core.String value) => $_setString(1, value);
  @$pb.TagNumber(375146466)
  $core.bool hasContentdisposition() => $_has(1);
  @$pb.TagNumber(375146466)
  void clearContentdisposition() => $_clearField(375146466);

  @$pb.TagNumber(464157046)
  $core.List<$core.int> get body => $_getN(2);
  @$pb.TagNumber(464157046)
  set body($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(2);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);
}

class SdkType extends $pb.GeneratedMessage {
  factory SdkType({
    $core.String? friendlyname,
    $core.Iterable<SdkConfigurationProperty>? configurationproperties,
    $core.String? description,
    $core.String? id,
  }) {
    final result = create();
    if (friendlyname != null) result.friendlyname = friendlyname;
    if (configurationproperties != null)
      result.configurationproperties.addAll(configurationproperties);
    if (description != null) result.description = description;
    if (id != null) result.id = id;
    return result;
  }

  SdkType._();

  factory SdkType.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SdkType.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SdkType',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(39239102, _omitFieldNames ? '' : 'friendlyname')
    ..pPM<SdkConfigurationProperty>(
        113650241, _omitFieldNames ? '' : 'configurationproperties',
        subBuilder: SdkConfigurationProperty.create)
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SdkType clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SdkType copyWith(void Function(SdkType) updates) =>
      super.copyWith((message) => updates(message as SdkType)) as SdkType;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SdkType create() => SdkType._();
  @$core.override
  SdkType createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SdkType getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SdkType>(create);
  static SdkType? _defaultInstance;

  @$pb.TagNumber(39239102)
  $core.String get friendlyname => $_getSZ(0);
  @$pb.TagNumber(39239102)
  set friendlyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39239102)
  $core.bool hasFriendlyname() => $_has(0);
  @$pb.TagNumber(39239102)
  void clearFriendlyname() => $_clearField(39239102);

  @$pb.TagNumber(113650241)
  $pb.PbList<SdkConfigurationProperty> get configurationproperties =>
      $_getList(1);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);
}

class SdkTypes extends $pb.GeneratedMessage {
  factory SdkTypes({
    $core.Iterable<SdkType>? items,
  }) {
    final result = create();
    if (items != null) result.items.addAll(items);
    return result;
  }

  SdkTypes._();

  factory SdkTypes.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SdkTypes.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SdkTypes',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<SdkType>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: SdkType.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SdkTypes clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SdkTypes copyWith(void Function(SdkTypes) updates) =>
      super.copyWith((message) => updates(message as SdkTypes)) as SdkTypes;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SdkTypes create() => SdkTypes._();
  @$core.override
  SdkTypes createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SdkTypes getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SdkTypes>(create);
  static SdkTypes? _defaultInstance;

  @$pb.TagNumber(444150672)
  $pb.PbList<SdkType> get items => $_getList(0);
}

class ServiceUnavailableException extends $pb.GeneratedMessage {
  factory ServiceUnavailableException({
    $core.String? message,
    $core.String? retryafterseconds,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (retryafterseconds != null) result.retryafterseconds = retryafterseconds;
    return result;
  }

  ServiceUnavailableException._();

  factory ServiceUnavailableException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ServiceUnavailableException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ServiceUnavailableException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..aOS(436039555, _omitFieldNames ? '' : 'retryafterseconds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServiceUnavailableException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServiceUnavailableException copyWith(
          void Function(ServiceUnavailableException) updates) =>
      super.copyWith(
              (message) => updates(message as ServiceUnavailableException))
          as ServiceUnavailableException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ServiceUnavailableException create() =>
      ServiceUnavailableException._();
  @$core.override
  ServiceUnavailableException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ServiceUnavailableException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ServiceUnavailableException>(create);
  static ServiceUnavailableException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(436039555)
  $core.String get retryafterseconds => $_getSZ(1);
  @$pb.TagNumber(436039555)
  set retryafterseconds($core.String value) => $_setString(1, value);
  @$pb.TagNumber(436039555)
  $core.bool hasRetryafterseconds() => $_has(1);
  @$pb.TagNumber(436039555)
  void clearRetryafterseconds() => $_clearField(436039555);
}

class Stage extends $pb.GeneratedMessage {
  factory Stage({
    $core.String? stagename,
    $core.Iterable<$core.MapEntry<$core.String, MethodSetting>>? methodsettings,
    $core.String? createddate,
    $core.bool? cacheclusterenabled,
    AccessLogSettings? accesslogsettings,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? variables,
    $core.String? documentationversion,
    CacheClusterSize? cacheclustersize,
    $core.String? webaclarn,
    $core.String? clientcertificateid,
    CanarySettings? canarysettings,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    CacheClusterStatus? cacheclusterstatus,
    $core.bool? tracingenabled,
    $core.String? deploymentid,
    $core.String? lastupdateddate,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (methodsettings != null)
      result.methodsettings.addEntries(methodsettings);
    if (createddate != null) result.createddate = createddate;
    if (cacheclusterenabled != null)
      result.cacheclusterenabled = cacheclusterenabled;
    if (accesslogsettings != null) result.accesslogsettings = accesslogsettings;
    if (variables != null) result.variables.addEntries(variables);
    if (documentationversion != null)
      result.documentationversion = documentationversion;
    if (cacheclustersize != null) result.cacheclustersize = cacheclustersize;
    if (webaclarn != null) result.webaclarn = webaclarn;
    if (clientcertificateid != null)
      result.clientcertificateid = clientcertificateid;
    if (canarysettings != null) result.canarysettings = canarysettings;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (cacheclusterstatus != null)
      result.cacheclusterstatus = cacheclusterstatus;
    if (tracingenabled != null) result.tracingenabled = tracingenabled;
    if (deploymentid != null) result.deploymentid = deploymentid;
    if (lastupdateddate != null) result.lastupdateddate = lastupdateddate;
    return result;
  }

  Stage._();

  factory Stage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Stage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Stage',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..m<$core.String, MethodSetting>(
        30387838, _omitFieldNames ? '' : 'methodsettings',
        entryClassName: 'Stage.MethodsettingsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: MethodSetting.create,
        valueDefaultOrMaker: MethodSetting.getDefault,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(53061200, _omitFieldNames ? '' : 'createddate')
    ..aOB(63967991, _omitFieldNames ? '' : 'cacheclusterenabled')
    ..aOM<AccessLogSettings>(
        76216807, _omitFieldNames ? '' : 'accesslogsettings',
        subBuilder: AccessLogSettings.create)
    ..m<$core.String, $core.String>(
        162226883, _omitFieldNames ? '' : 'variables',
        entryClassName: 'Stage.VariablesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(167009804, _omitFieldNames ? '' : 'documentationversion')
    ..aE<CacheClusterSize>(232189861, _omitFieldNames ? '' : 'cacheclustersize',
        enumValues: CacheClusterSize.values)
    ..aOS(243701763, _omitFieldNames ? '' : 'webaclarn')
    ..aOS(276222909, _omitFieldNames ? '' : 'clientcertificateid')
    ..aOM<CanarySettings>(285544261, _omitFieldNames ? '' : 'canarysettings',
        subBuilder: CanarySettings.create)
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'Stage.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aE<CacheClusterStatus>(
        385293784, _omitFieldNames ? '' : 'cacheclusterstatus',
        enumValues: CacheClusterStatus.values)
    ..aOB(390995731, _omitFieldNames ? '' : 'tracingenabled')
    ..aOS(439369188, _omitFieldNames ? '' : 'deploymentid')
    ..aOS(448453361, _omitFieldNames ? '' : 'lastupdateddate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Stage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Stage copyWith(void Function(Stage) updates) =>
      super.copyWith((message) => updates(message as Stage)) as Stage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Stage create() => Stage._();
  @$core.override
  Stage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Stage getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Stage>(create);
  static Stage? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(30387838)
  $pb.PbMap<$core.String, MethodSetting> get methodsettings => $_getMap(1);

  @$pb.TagNumber(53061200)
  $core.String get createddate => $_getSZ(2);
  @$pb.TagNumber(53061200)
  set createddate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(53061200)
  $core.bool hasCreateddate() => $_has(2);
  @$pb.TagNumber(53061200)
  void clearCreateddate() => $_clearField(53061200);

  @$pb.TagNumber(63967991)
  $core.bool get cacheclusterenabled => $_getBF(3);
  @$pb.TagNumber(63967991)
  set cacheclusterenabled($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(63967991)
  $core.bool hasCacheclusterenabled() => $_has(3);
  @$pb.TagNumber(63967991)
  void clearCacheclusterenabled() => $_clearField(63967991);

  @$pb.TagNumber(76216807)
  AccessLogSettings get accesslogsettings => $_getN(4);
  @$pb.TagNumber(76216807)
  set accesslogsettings(AccessLogSettings value) => $_setField(76216807, value);
  @$pb.TagNumber(76216807)
  $core.bool hasAccesslogsettings() => $_has(4);
  @$pb.TagNumber(76216807)
  void clearAccesslogsettings() => $_clearField(76216807);
  @$pb.TagNumber(76216807)
  AccessLogSettings ensureAccesslogsettings() => $_ensure(4);

  @$pb.TagNumber(162226883)
  $pb.PbMap<$core.String, $core.String> get variables => $_getMap(5);

  @$pb.TagNumber(167009804)
  $core.String get documentationversion => $_getSZ(6);
  @$pb.TagNumber(167009804)
  set documentationversion($core.String value) => $_setString(6, value);
  @$pb.TagNumber(167009804)
  $core.bool hasDocumentationversion() => $_has(6);
  @$pb.TagNumber(167009804)
  void clearDocumentationversion() => $_clearField(167009804);

  @$pb.TagNumber(232189861)
  CacheClusterSize get cacheclustersize => $_getN(7);
  @$pb.TagNumber(232189861)
  set cacheclustersize(CacheClusterSize value) => $_setField(232189861, value);
  @$pb.TagNumber(232189861)
  $core.bool hasCacheclustersize() => $_has(7);
  @$pb.TagNumber(232189861)
  void clearCacheclustersize() => $_clearField(232189861);

  @$pb.TagNumber(243701763)
  $core.String get webaclarn => $_getSZ(8);
  @$pb.TagNumber(243701763)
  set webaclarn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(243701763)
  $core.bool hasWebaclarn() => $_has(8);
  @$pb.TagNumber(243701763)
  void clearWebaclarn() => $_clearField(243701763);

  @$pb.TagNumber(276222909)
  $core.String get clientcertificateid => $_getSZ(9);
  @$pb.TagNumber(276222909)
  set clientcertificateid($core.String value) => $_setString(9, value);
  @$pb.TagNumber(276222909)
  $core.bool hasClientcertificateid() => $_has(9);
  @$pb.TagNumber(276222909)
  void clearClientcertificateid() => $_clearField(276222909);

  @$pb.TagNumber(285544261)
  CanarySettings get canarysettings => $_getN(10);
  @$pb.TagNumber(285544261)
  set canarysettings(CanarySettings value) => $_setField(285544261, value);
  @$pb.TagNumber(285544261)
  $core.bool hasCanarysettings() => $_has(10);
  @$pb.TagNumber(285544261)
  void clearCanarysettings() => $_clearField(285544261);
  @$pb.TagNumber(285544261)
  CanarySettings ensureCanarysettings() => $_ensure(10);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(11);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(12);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(12, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(12);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(385293784)
  CacheClusterStatus get cacheclusterstatus => $_getN(13);
  @$pb.TagNumber(385293784)
  set cacheclusterstatus(CacheClusterStatus value) =>
      $_setField(385293784, value);
  @$pb.TagNumber(385293784)
  $core.bool hasCacheclusterstatus() => $_has(13);
  @$pb.TagNumber(385293784)
  void clearCacheclusterstatus() => $_clearField(385293784);

  @$pb.TagNumber(390995731)
  $core.bool get tracingenabled => $_getBF(14);
  @$pb.TagNumber(390995731)
  set tracingenabled($core.bool value) => $_setBool(14, value);
  @$pb.TagNumber(390995731)
  $core.bool hasTracingenabled() => $_has(14);
  @$pb.TagNumber(390995731)
  void clearTracingenabled() => $_clearField(390995731);

  @$pb.TagNumber(439369188)
  $core.String get deploymentid => $_getSZ(15);
  @$pb.TagNumber(439369188)
  set deploymentid($core.String value) => $_setString(15, value);
  @$pb.TagNumber(439369188)
  $core.bool hasDeploymentid() => $_has(15);
  @$pb.TagNumber(439369188)
  void clearDeploymentid() => $_clearField(439369188);

  @$pb.TagNumber(448453361)
  $core.String get lastupdateddate => $_getSZ(16);
  @$pb.TagNumber(448453361)
  set lastupdateddate($core.String value) => $_setString(16, value);
  @$pb.TagNumber(448453361)
  $core.bool hasLastupdateddate() => $_has(16);
  @$pb.TagNumber(448453361)
  void clearLastupdateddate() => $_clearField(448453361);
}

class StageKey extends $pb.GeneratedMessage {
  factory StageKey({
    $core.String? stagename,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  StageKey._();

  factory StageKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StageKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StageKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StageKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StageKey copyWith(void Function(StageKey) updates) =>
      super.copyWith((message) => updates(message as StageKey)) as StageKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StageKey create() => StageKey._();
  @$core.override
  StageKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StageKey getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<StageKey>(create);
  static StageKey? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class Stages extends $pb.GeneratedMessage {
  factory Stages({
    $core.Iterable<Stage>? item,
  }) {
    final result = create();
    if (item != null) result.item.addAll(item);
    return result;
  }

  Stages._();

  factory Stages.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Stages.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Stages',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<Stage>(523776999, _omitFieldNames ? '' : 'item',
        subBuilder: Stage.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Stages clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Stages copyWith(void Function(Stages) updates) =>
      super.copyWith((message) => updates(message as Stages)) as Stages;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Stages create() => Stages._();
  @$core.override
  Stages createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Stages getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Stages>(create);
  static Stages? _defaultInstance;

  @$pb.TagNumber(523776999)
  $pb.PbList<Stage> get item => $_getList(0);
}

class TagResourceRequest extends $pb.GeneratedMessage {
  factory TagResourceRequest({
    $core.String? resourcearn,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (tags != null) result.tags.addEntries(tags);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(67806797, _omitFieldNames ? '' : 'resourcearn')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'TagResourceRequest.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
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

  @$pb.TagNumber(67806797)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(67806797)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(67806797)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(67806797)
  void clearResourcearn() => $_clearField(67806797);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(1);
}

class Tags extends $pb.GeneratedMessage {
  factory Tags({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
  }) {
    final result = create();
    if (tags != null) result.tags.addEntries(tags);
    return result;
  }

  Tags._();

  factory Tags.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Tags.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Tags',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'Tags.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Tags clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Tags copyWith(void Function(Tags) updates) =>
      super.copyWith((message) => updates(message as Tags)) as Tags;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Tags create() => Tags._();
  @$core.override
  Tags createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Tags getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Tags>(create);
  static Tags? _defaultInstance;

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(0);
}

class Template extends $pb.GeneratedMessage {
  factory Template({
    $core.String? value,
  }) {
    final result = create();
    if (value != null) result.value = value;
    return result;
  }

  Template._();

  factory Template.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Template.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Template',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(39769035, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Template clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Template copyWith(void Function(Template) updates) =>
      super.copyWith((message) => updates(message as Template)) as Template;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Template create() => Template._();
  @$core.override
  Template createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Template getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Template>(create);
  static Template? _defaultInstance;

  @$pb.TagNumber(39769035)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(39769035)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39769035)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(39769035)
  void clearValue() => $_clearField(39769035);
}

class TestInvokeAuthorizerRequest extends $pb.GeneratedMessage {
  factory TestInvokeAuthorizerRequest({
    $core.String? authorizerid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        multivalueheaders,
    $core.String? pathwithquerystring,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? stagevariables,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? headers,
    $core.String? restapiid,
    $core.String? body,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        additionalcontext,
  }) {
    final result = create();
    if (authorizerid != null) result.authorizerid = authorizerid;
    if (multivalueheaders != null)
      result.multivalueheaders.addEntries(multivalueheaders);
    if (pathwithquerystring != null)
      result.pathwithquerystring = pathwithquerystring;
    if (stagevariables != null)
      result.stagevariables.addEntries(stagevariables);
    if (headers != null) result.headers.addEntries(headers);
    if (restapiid != null) result.restapiid = restapiid;
    if (body != null) result.body = body;
    if (additionalcontext != null)
      result.additionalcontext.addEntries(additionalcontext);
    return result;
  }

  TestInvokeAuthorizerRequest._();

  factory TestInvokeAuthorizerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestInvokeAuthorizerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestInvokeAuthorizerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(111773148, _omitFieldNames ? '' : 'authorizerid')
    ..m<$core.String, $core.String>(
        142421420, _omitFieldNames ? '' : 'multivalueheaders',
        entryClassName: 'TestInvokeAuthorizerRequest.MultivalueheadersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(159950736, _omitFieldNames ? '' : 'pathwithquerystring')
    ..m<$core.String, $core.String>(
        208524923, _omitFieldNames ? '' : 'stagevariables',
        entryClassName: 'TestInvokeAuthorizerRequest.StagevariablesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..m<$core.String, $core.String>(375773674, _omitFieldNames ? '' : 'headers',
        entryClassName: 'TestInvokeAuthorizerRequest.HeadersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(464157046, _omitFieldNames ? '' : 'body')
    ..m<$core.String, $core.String>(
        471949152, _omitFieldNames ? '' : 'additionalcontext',
        entryClassName: 'TestInvokeAuthorizerRequest.AdditionalcontextEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestInvokeAuthorizerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestInvokeAuthorizerRequest copyWith(
          void Function(TestInvokeAuthorizerRequest) updates) =>
      super.copyWith(
              (message) => updates(message as TestInvokeAuthorizerRequest))
          as TestInvokeAuthorizerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestInvokeAuthorizerRequest create() =>
      TestInvokeAuthorizerRequest._();
  @$core.override
  TestInvokeAuthorizerRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestInvokeAuthorizerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestInvokeAuthorizerRequest>(create);
  static TestInvokeAuthorizerRequest? _defaultInstance;

  @$pb.TagNumber(111773148)
  $core.String get authorizerid => $_getSZ(0);
  @$pb.TagNumber(111773148)
  set authorizerid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(111773148)
  $core.bool hasAuthorizerid() => $_has(0);
  @$pb.TagNumber(111773148)
  void clearAuthorizerid() => $_clearField(111773148);

  @$pb.TagNumber(142421420)
  $pb.PbMap<$core.String, $core.String> get multivalueheaders => $_getMap(1);

  @$pb.TagNumber(159950736)
  $core.String get pathwithquerystring => $_getSZ(2);
  @$pb.TagNumber(159950736)
  set pathwithquerystring($core.String value) => $_setString(2, value);
  @$pb.TagNumber(159950736)
  $core.bool hasPathwithquerystring() => $_has(2);
  @$pb.TagNumber(159950736)
  void clearPathwithquerystring() => $_clearField(159950736);

  @$pb.TagNumber(208524923)
  $pb.PbMap<$core.String, $core.String> get stagevariables => $_getMap(3);

  @$pb.TagNumber(375773674)
  $pb.PbMap<$core.String, $core.String> get headers => $_getMap(4);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(5);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(5);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(464157046)
  $core.String get body => $_getSZ(6);
  @$pb.TagNumber(464157046)
  set body($core.String value) => $_setString(6, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(6);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);

  @$pb.TagNumber(471949152)
  $pb.PbMap<$core.String, $core.String> get additionalcontext => $_getMap(7);
}

class TestInvokeAuthorizerResponse extends $pb.GeneratedMessage {
  factory TestInvokeAuthorizerResponse({
    $core.int? clientstatus,
    $core.String? policy,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? authorization,
    $fixnum.Int64? latency,
    $core.String? principalid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? claims,
    $core.String? log,
  }) {
    final result = create();
    if (clientstatus != null) result.clientstatus = clientstatus;
    if (policy != null) result.policy = policy;
    if (authorization != null) result.authorization.addEntries(authorization);
    if (latency != null) result.latency = latency;
    if (principalid != null) result.principalid = principalid;
    if (claims != null) result.claims.addEntries(claims);
    if (log != null) result.log = log;
    return result;
  }

  TestInvokeAuthorizerResponse._();

  factory TestInvokeAuthorizerResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestInvokeAuthorizerResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestInvokeAuthorizerResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(35642913, _omitFieldNames ? '' : 'clientstatus')
    ..aOS(247528064, _omitFieldNames ? '' : 'policy')
    ..m<$core.String, $core.String>(
        288774079, _omitFieldNames ? '' : 'authorization',
        entryClassName: 'TestInvokeAuthorizerResponse.AuthorizationEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aInt64(318473050, _omitFieldNames ? '' : 'latency')
    ..aOS(350710285, _omitFieldNames ? '' : 'principalid')
    ..m<$core.String, $core.String>(479124501, _omitFieldNames ? '' : 'claims',
        entryClassName: 'TestInvokeAuthorizerResponse.ClaimsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(525422930, _omitFieldNames ? '' : 'log')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestInvokeAuthorizerResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestInvokeAuthorizerResponse copyWith(
          void Function(TestInvokeAuthorizerResponse) updates) =>
      super.copyWith(
              (message) => updates(message as TestInvokeAuthorizerResponse))
          as TestInvokeAuthorizerResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestInvokeAuthorizerResponse create() =>
      TestInvokeAuthorizerResponse._();
  @$core.override
  TestInvokeAuthorizerResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestInvokeAuthorizerResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestInvokeAuthorizerResponse>(create);
  static TestInvokeAuthorizerResponse? _defaultInstance;

  @$pb.TagNumber(35642913)
  $core.int get clientstatus => $_getIZ(0);
  @$pb.TagNumber(35642913)
  set clientstatus($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(35642913)
  $core.bool hasClientstatus() => $_has(0);
  @$pb.TagNumber(35642913)
  void clearClientstatus() => $_clearField(35642913);

  @$pb.TagNumber(247528064)
  $core.String get policy => $_getSZ(1);
  @$pb.TagNumber(247528064)
  set policy($core.String value) => $_setString(1, value);
  @$pb.TagNumber(247528064)
  $core.bool hasPolicy() => $_has(1);
  @$pb.TagNumber(247528064)
  void clearPolicy() => $_clearField(247528064);

  @$pb.TagNumber(288774079)
  $pb.PbMap<$core.String, $core.String> get authorization => $_getMap(2);

  @$pb.TagNumber(318473050)
  $fixnum.Int64 get latency => $_getI64(3);
  @$pb.TagNumber(318473050)
  set latency($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(318473050)
  $core.bool hasLatency() => $_has(3);
  @$pb.TagNumber(318473050)
  void clearLatency() => $_clearField(318473050);

  @$pb.TagNumber(350710285)
  $core.String get principalid => $_getSZ(4);
  @$pb.TagNumber(350710285)
  set principalid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(350710285)
  $core.bool hasPrincipalid() => $_has(4);
  @$pb.TagNumber(350710285)
  void clearPrincipalid() => $_clearField(350710285);

  @$pb.TagNumber(479124501)
  $pb.PbMap<$core.String, $core.String> get claims => $_getMap(5);

  @$pb.TagNumber(525422930)
  $core.String get log => $_getSZ(6);
  @$pb.TagNumber(525422930)
  set log($core.String value) => $_setString(6, value);
  @$pb.TagNumber(525422930)
  $core.bool hasLog() => $_has(6);
  @$pb.TagNumber(525422930)
  void clearLog() => $_clearField(525422930);
}

class TestInvokeMethodRequest extends $pb.GeneratedMessage {
  factory TestInvokeMethodRequest({
    $core.String? httpmethod,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        multivalueheaders,
    $core.String? pathwithquerystring,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? stagevariables,
    $core.String? clientcertificateid,
    $core.String? resourceid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? headers,
    $core.String? restapiid,
    $core.String? body,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (multivalueheaders != null)
      result.multivalueheaders.addEntries(multivalueheaders);
    if (pathwithquerystring != null)
      result.pathwithquerystring = pathwithquerystring;
    if (stagevariables != null)
      result.stagevariables.addEntries(stagevariables);
    if (clientcertificateid != null)
      result.clientcertificateid = clientcertificateid;
    if (resourceid != null) result.resourceid = resourceid;
    if (headers != null) result.headers.addEntries(headers);
    if (restapiid != null) result.restapiid = restapiid;
    if (body != null) result.body = body;
    return result;
  }

  TestInvokeMethodRequest._();

  factory TestInvokeMethodRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestInvokeMethodRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestInvokeMethodRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..m<$core.String, $core.String>(
        142421420, _omitFieldNames ? '' : 'multivalueheaders',
        entryClassName: 'TestInvokeMethodRequest.MultivalueheadersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(159950736, _omitFieldNames ? '' : 'pathwithquerystring')
    ..m<$core.String, $core.String>(
        208524923, _omitFieldNames ? '' : 'stagevariables',
        entryClassName: 'TestInvokeMethodRequest.StagevariablesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(276222909, _omitFieldNames ? '' : 'clientcertificateid')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..m<$core.String, $core.String>(375773674, _omitFieldNames ? '' : 'headers',
        entryClassName: 'TestInvokeMethodRequest.HeadersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(464157046, _omitFieldNames ? '' : 'body')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestInvokeMethodRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestInvokeMethodRequest copyWith(
          void Function(TestInvokeMethodRequest) updates) =>
      super.copyWith((message) => updates(message as TestInvokeMethodRequest))
          as TestInvokeMethodRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestInvokeMethodRequest create() => TestInvokeMethodRequest._();
  @$core.override
  TestInvokeMethodRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestInvokeMethodRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestInvokeMethodRequest>(create);
  static TestInvokeMethodRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(142421420)
  $pb.PbMap<$core.String, $core.String> get multivalueheaders => $_getMap(1);

  @$pb.TagNumber(159950736)
  $core.String get pathwithquerystring => $_getSZ(2);
  @$pb.TagNumber(159950736)
  set pathwithquerystring($core.String value) => $_setString(2, value);
  @$pb.TagNumber(159950736)
  $core.bool hasPathwithquerystring() => $_has(2);
  @$pb.TagNumber(159950736)
  void clearPathwithquerystring() => $_clearField(159950736);

  @$pb.TagNumber(208524923)
  $pb.PbMap<$core.String, $core.String> get stagevariables => $_getMap(3);

  @$pb.TagNumber(276222909)
  $core.String get clientcertificateid => $_getSZ(4);
  @$pb.TagNumber(276222909)
  set clientcertificateid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(276222909)
  $core.bool hasClientcertificateid() => $_has(4);
  @$pb.TagNumber(276222909)
  void clearClientcertificateid() => $_clearField(276222909);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(5);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(5);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(375773674)
  $pb.PbMap<$core.String, $core.String> get headers => $_getMap(6);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(7);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(7, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(7);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(464157046)
  $core.String get body => $_getSZ(8);
  @$pb.TagNumber(464157046)
  set body($core.String value) => $_setString(8, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(8);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);
}

class TestInvokeMethodResponse extends $pb.GeneratedMessage {
  factory TestInvokeMethodResponse({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        multivalueheaders,
    $fixnum.Int64? latency,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? headers,
    $core.int? status,
    $core.String? body,
    $core.String? log,
  }) {
    final result = create();
    if (multivalueheaders != null)
      result.multivalueheaders.addEntries(multivalueheaders);
    if (latency != null) result.latency = latency;
    if (headers != null) result.headers.addEntries(headers);
    if (status != null) result.status = status;
    if (body != null) result.body = body;
    if (log != null) result.log = log;
    return result;
  }

  TestInvokeMethodResponse._();

  factory TestInvokeMethodResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestInvokeMethodResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestInvokeMethodResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        142421420, _omitFieldNames ? '' : 'multivalueheaders',
        entryClassName: 'TestInvokeMethodResponse.MultivalueheadersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aInt64(318473050, _omitFieldNames ? '' : 'latency')
    ..m<$core.String, $core.String>(375773674, _omitFieldNames ? '' : 'headers',
        entryClassName: 'TestInvokeMethodResponse.HeadersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aI(441153520, _omitFieldNames ? '' : 'status')
    ..aOS(464157046, _omitFieldNames ? '' : 'body')
    ..aOS(525422930, _omitFieldNames ? '' : 'log')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestInvokeMethodResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestInvokeMethodResponse copyWith(
          void Function(TestInvokeMethodResponse) updates) =>
      super.copyWith((message) => updates(message as TestInvokeMethodResponse))
          as TestInvokeMethodResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestInvokeMethodResponse create() => TestInvokeMethodResponse._();
  @$core.override
  TestInvokeMethodResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestInvokeMethodResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestInvokeMethodResponse>(create);
  static TestInvokeMethodResponse? _defaultInstance;

  @$pb.TagNumber(142421420)
  $pb.PbMap<$core.String, $core.String> get multivalueheaders => $_getMap(0);

  @$pb.TagNumber(318473050)
  $fixnum.Int64 get latency => $_getI64(1);
  @$pb.TagNumber(318473050)
  set latency($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(318473050)
  $core.bool hasLatency() => $_has(1);
  @$pb.TagNumber(318473050)
  void clearLatency() => $_clearField(318473050);

  @$pb.TagNumber(375773674)
  $pb.PbMap<$core.String, $core.String> get headers => $_getMap(2);

  @$pb.TagNumber(441153520)
  $core.int get status => $_getIZ(3);
  @$pb.TagNumber(441153520)
  set status($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(441153520)
  $core.bool hasStatus() => $_has(3);
  @$pb.TagNumber(441153520)
  void clearStatus() => $_clearField(441153520);

  @$pb.TagNumber(464157046)
  $core.String get body => $_getSZ(4);
  @$pb.TagNumber(464157046)
  set body($core.String value) => $_setString(4, value);
  @$pb.TagNumber(464157046)
  $core.bool hasBody() => $_has(4);
  @$pb.TagNumber(464157046)
  void clearBody() => $_clearField(464157046);

  @$pb.TagNumber(525422930)
  $core.String get log => $_getSZ(5);
  @$pb.TagNumber(525422930)
  set log($core.String value) => $_setString(5, value);
  @$pb.TagNumber(525422930)
  $core.bool hasLog() => $_has(5);
  @$pb.TagNumber(525422930)
  void clearLog() => $_clearField(525422930);
}

class ThrottleSettings extends $pb.GeneratedMessage {
  factory ThrottleSettings({
    $core.int? burstlimit,
    $core.double? ratelimit,
  }) {
    final result = create();
    if (burstlimit != null) result.burstlimit = burstlimit;
    if (ratelimit != null) result.ratelimit = ratelimit;
    return result;
  }

  ThrottleSettings._();

  factory ThrottleSettings.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ThrottleSettings.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ThrottleSettings',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aI(37855041, _omitFieldNames ? '' : 'burstlimit')
    ..aD(505789539, _omitFieldNames ? '' : 'ratelimit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ThrottleSettings clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ThrottleSettings copyWith(void Function(ThrottleSettings) updates) =>
      super.copyWith((message) => updates(message as ThrottleSettings))
          as ThrottleSettings;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ThrottleSettings create() => ThrottleSettings._();
  @$core.override
  ThrottleSettings createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ThrottleSettings getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ThrottleSettings>(create);
  static ThrottleSettings? _defaultInstance;

  @$pb.TagNumber(37855041)
  $core.int get burstlimit => $_getIZ(0);
  @$pb.TagNumber(37855041)
  set burstlimit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(37855041)
  $core.bool hasBurstlimit() => $_has(0);
  @$pb.TagNumber(37855041)
  void clearBurstlimit() => $_clearField(37855041);

  @$pb.TagNumber(505789539)
  $core.double get ratelimit => $_getN(1);
  @$pb.TagNumber(505789539)
  set ratelimit($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(505789539)
  $core.bool hasRatelimit() => $_has(1);
  @$pb.TagNumber(505789539)
  void clearRatelimit() => $_clearField(505789539);
}

class TlsConfig extends $pb.GeneratedMessage {
  factory TlsConfig({
    $core.bool? insecureskipverification,
  }) {
    final result = create();
    if (insecureskipverification != null)
      result.insecureskipverification = insecureskipverification;
    return result;
  }

  TlsConfig._();

  factory TlsConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TlsConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TlsConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOB(78100228, _omitFieldNames ? '' : 'insecureskipverification')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TlsConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TlsConfig copyWith(void Function(TlsConfig) updates) =>
      super.copyWith((message) => updates(message as TlsConfig)) as TlsConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TlsConfig create() => TlsConfig._();
  @$core.override
  TlsConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TlsConfig getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TlsConfig>(create);
  static TlsConfig? _defaultInstance;

  @$pb.TagNumber(78100228)
  $core.bool get insecureskipverification => $_getBF(0);
  @$pb.TagNumber(78100228)
  set insecureskipverification($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(78100228)
  $core.bool hasInsecureskipverification() => $_has(0);
  @$pb.TagNumber(78100228)
  void clearInsecureskipverification() => $_clearField(78100228);
}

class TooManyRequestsException extends $pb.GeneratedMessage {
  factory TooManyRequestsException({
    $core.String? message,
    $core.String? retryafterseconds,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (retryafterseconds != null) result.retryafterseconds = retryafterseconds;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..aOS(436039555, _omitFieldNames ? '' : 'retryafterseconds')
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

  @$pb.TagNumber(436039555)
  $core.String get retryafterseconds => $_getSZ(1);
  @$pb.TagNumber(436039555)
  set retryafterseconds($core.String value) => $_setString(1, value);
  @$pb.TagNumber(436039555)
  $core.bool hasRetryafterseconds() => $_has(1);
  @$pb.TagNumber(436039555)
  void clearRetryafterseconds() => $_clearField(436039555);
}

class UnauthorizedException extends $pb.GeneratedMessage {
  factory UnauthorizedException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  UnauthorizedException._();

  factory UnauthorizedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnauthorizedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnauthorizedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnauthorizedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnauthorizedException copyWith(
          void Function(UnauthorizedException) updates) =>
      super.copyWith((message) => updates(message as UnauthorizedException))
          as UnauthorizedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnauthorizedException create() => UnauthorizedException._();
  @$core.override
  UnauthorizedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnauthorizedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnauthorizedException>(create);
  static UnauthorizedException? _defaultInstance;

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
    $core.String? resourcearn,
    $core.Iterable<$core.String>? tagkeys,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(67806797, _omitFieldNames ? '' : 'resourcearn')
    ..pPS(78811036, _omitFieldNames ? '' : 'tagkeys')
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

class UpdateAccountRequest extends $pb.GeneratedMessage {
  factory UpdateAccountRequest({
    $core.Iterable<PatchOperation>? patchoperations,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    return result;
  }

  UpdateAccountRequest._();

  factory UpdateAccountRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateAccountRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAccountRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAccountRequest copyWith(void Function(UpdateAccountRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateAccountRequest))
          as UpdateAccountRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAccountRequest create() => UpdateAccountRequest._();
  @$core.override
  UpdateAccountRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateAccountRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateAccountRequest>(create);
  static UpdateAccountRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);
}

class UpdateApiKeyRequest extends $pb.GeneratedMessage {
  factory UpdateApiKeyRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? apikey,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (apikey != null) result.apikey = apikey;
    return result;
  }

  UpdateApiKeyRequest._();

  factory UpdateApiKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateApiKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateApiKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(490490655, _omitFieldNames ? '' : 'apikey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateApiKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateApiKeyRequest copyWith(void Function(UpdateApiKeyRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateApiKeyRequest))
          as UpdateApiKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateApiKeyRequest create() => UpdateApiKeyRequest._();
  @$core.override
  UpdateApiKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateApiKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateApiKeyRequest>(create);
  static UpdateApiKeyRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(490490655)
  $core.String get apikey => $_getSZ(1);
  @$pb.TagNumber(490490655)
  set apikey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(490490655)
  $core.bool hasApikey() => $_has(1);
  @$pb.TagNumber(490490655)
  void clearApikey() => $_clearField(490490655);
}

class UpdateAuthorizerRequest extends $pb.GeneratedMessage {
  factory UpdateAuthorizerRequest({
    $core.String? authorizerid,
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? restapiid,
  }) {
    final result = create();
    if (authorizerid != null) result.authorizerid = authorizerid;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateAuthorizerRequest._();

  factory UpdateAuthorizerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateAuthorizerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateAuthorizerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(111773148, _omitFieldNames ? '' : 'authorizerid')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAuthorizerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateAuthorizerRequest copyWith(
          void Function(UpdateAuthorizerRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateAuthorizerRequest))
          as UpdateAuthorizerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateAuthorizerRequest create() => UpdateAuthorizerRequest._();
  @$core.override
  UpdateAuthorizerRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateAuthorizerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateAuthorizerRequest>(create);
  static UpdateAuthorizerRequest? _defaultInstance;

  @$pb.TagNumber(111773148)
  $core.String get authorizerid => $_getSZ(0);
  @$pb.TagNumber(111773148)
  set authorizerid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(111773148)
  $core.bool hasAuthorizerid() => $_has(0);
  @$pb.TagNumber(111773148)
  void clearAuthorizerid() => $_clearField(111773148);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateBasePathMappingRequest extends $pb.GeneratedMessage {
  factory UpdateBasePathMappingRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? basepath,
    $core.String? domainnameid,
    $core.String? domainname,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (basepath != null) result.basepath = basepath;
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  UpdateBasePathMappingRequest._();

  factory UpdateBasePathMappingRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateBasePathMappingRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateBasePathMappingRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(267528880, _omitFieldNames ? '' : 'basepath')
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateBasePathMappingRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateBasePathMappingRequest copyWith(
          void Function(UpdateBasePathMappingRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateBasePathMappingRequest))
          as UpdateBasePathMappingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateBasePathMappingRequest create() =>
      UpdateBasePathMappingRequest._();
  @$core.override
  UpdateBasePathMappingRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateBasePathMappingRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateBasePathMappingRequest>(create);
  static UpdateBasePathMappingRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(267528880)
  $core.String get basepath => $_getSZ(1);
  @$pb.TagNumber(267528880)
  set basepath($core.String value) => $_setString(1, value);
  @$pb.TagNumber(267528880)
  $core.bool hasBasepath() => $_has(1);
  @$pb.TagNumber(267528880)
  void clearBasepath() => $_clearField(267528880);

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(2);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(2);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(3);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(3);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);
}

class UpdateClientCertificateRequest extends $pb.GeneratedMessage {
  factory UpdateClientCertificateRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? clientcertificateid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (clientcertificateid != null)
      result.clientcertificateid = clientcertificateid;
    return result;
  }

  UpdateClientCertificateRequest._();

  factory UpdateClientCertificateRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateClientCertificateRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateClientCertificateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(276222909, _omitFieldNames ? '' : 'clientcertificateid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateClientCertificateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateClientCertificateRequest copyWith(
          void Function(UpdateClientCertificateRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateClientCertificateRequest))
          as UpdateClientCertificateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateClientCertificateRequest create() =>
      UpdateClientCertificateRequest._();
  @$core.override
  UpdateClientCertificateRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateClientCertificateRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateClientCertificateRequest>(create);
  static UpdateClientCertificateRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(276222909)
  $core.String get clientcertificateid => $_getSZ(1);
  @$pb.TagNumber(276222909)
  set clientcertificateid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(276222909)
  $core.bool hasClientcertificateid() => $_has(1);
  @$pb.TagNumber(276222909)
  void clearClientcertificateid() => $_clearField(276222909);
}

class UpdateDeploymentRequest extends $pb.GeneratedMessage {
  factory UpdateDeploymentRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? restapiid,
    $core.String? deploymentid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (restapiid != null) result.restapiid = restapiid;
    if (deploymentid != null) result.deploymentid = deploymentid;
    return result;
  }

  UpdateDeploymentRequest._();

  factory UpdateDeploymentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDeploymentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDeploymentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(439369188, _omitFieldNames ? '' : 'deploymentid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDeploymentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDeploymentRequest copyWith(
          void Function(UpdateDeploymentRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateDeploymentRequest))
          as UpdateDeploymentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDeploymentRequest create() => UpdateDeploymentRequest._();
  @$core.override
  UpdateDeploymentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDeploymentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDeploymentRequest>(create);
  static UpdateDeploymentRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(439369188)
  $core.String get deploymentid => $_getSZ(2);
  @$pb.TagNumber(439369188)
  set deploymentid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(439369188)
  $core.bool hasDeploymentid() => $_has(2);
  @$pb.TagNumber(439369188)
  void clearDeploymentid() => $_clearField(439369188);
}

class UpdateDocumentationPartRequest extends $pb.GeneratedMessage {
  factory UpdateDocumentationPartRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? documentationpartid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (documentationpartid != null)
      result.documentationpartid = documentationpartid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateDocumentationPartRequest._();

  factory UpdateDocumentationPartRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDocumentationPartRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDocumentationPartRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(286552774, _omitFieldNames ? '' : 'documentationpartid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDocumentationPartRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDocumentationPartRequest copyWith(
          void Function(UpdateDocumentationPartRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateDocumentationPartRequest))
          as UpdateDocumentationPartRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDocumentationPartRequest create() =>
      UpdateDocumentationPartRequest._();
  @$core.override
  UpdateDocumentationPartRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDocumentationPartRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDocumentationPartRequest>(create);
  static UpdateDocumentationPartRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(286552774)
  $core.String get documentationpartid => $_getSZ(1);
  @$pb.TagNumber(286552774)
  set documentationpartid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(286552774)
  $core.bool hasDocumentationpartid() => $_has(1);
  @$pb.TagNumber(286552774)
  void clearDocumentationpartid() => $_clearField(286552774);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateDocumentationVersionRequest extends $pb.GeneratedMessage {
  factory UpdateDocumentationVersionRequest({
    $core.String? documentationversion,
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? restapiid,
  }) {
    final result = create();
    if (documentationversion != null)
      result.documentationversion = documentationversion;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateDocumentationVersionRequest._();

  factory UpdateDocumentationVersionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDocumentationVersionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDocumentationVersionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(167009804, _omitFieldNames ? '' : 'documentationversion')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDocumentationVersionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDocumentationVersionRequest copyWith(
          void Function(UpdateDocumentationVersionRequest) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateDocumentationVersionRequest))
          as UpdateDocumentationVersionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDocumentationVersionRequest create() =>
      UpdateDocumentationVersionRequest._();
  @$core.override
  UpdateDocumentationVersionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDocumentationVersionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDocumentationVersionRequest>(
          create);
  static UpdateDocumentationVersionRequest? _defaultInstance;

  @$pb.TagNumber(167009804)
  $core.String get documentationversion => $_getSZ(0);
  @$pb.TagNumber(167009804)
  set documentationversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(167009804)
  $core.bool hasDocumentationversion() => $_has(0);
  @$pb.TagNumber(167009804)
  void clearDocumentationversion() => $_clearField(167009804);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateDomainNameRequest extends $pb.GeneratedMessage {
  factory UpdateDomainNameRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? domainnameid,
    $core.String? domainname,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (domainnameid != null) result.domainnameid = domainnameid;
    if (domainname != null) result.domainname = domainname;
    return result;
  }

  UpdateDomainNameRequest._();

  factory UpdateDomainNameRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateDomainNameRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateDomainNameRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(298270248, _omitFieldNames ? '' : 'domainnameid')
    ..aOS(390326667, _omitFieldNames ? '' : 'domainname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDomainNameRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateDomainNameRequest copyWith(
          void Function(UpdateDomainNameRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateDomainNameRequest))
          as UpdateDomainNameRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateDomainNameRequest create() => UpdateDomainNameRequest._();
  @$core.override
  UpdateDomainNameRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateDomainNameRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateDomainNameRequest>(create);
  static UpdateDomainNameRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(298270248)
  $core.String get domainnameid => $_getSZ(1);
  @$pb.TagNumber(298270248)
  set domainnameid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(298270248)
  $core.bool hasDomainnameid() => $_has(1);
  @$pb.TagNumber(298270248)
  void clearDomainnameid() => $_clearField(298270248);

  @$pb.TagNumber(390326667)
  $core.String get domainname => $_getSZ(2);
  @$pb.TagNumber(390326667)
  set domainname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(390326667)
  $core.bool hasDomainname() => $_has(2);
  @$pb.TagNumber(390326667)
  void clearDomainname() => $_clearField(390326667);
}

class UpdateGatewayResponseRequest extends $pb.GeneratedMessage {
  factory UpdateGatewayResponseRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    GatewayResponseType? responsetype,
    $core.String? restapiid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (responsetype != null) result.responsetype = responsetype;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateGatewayResponseRequest._();

  factory UpdateGatewayResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateGatewayResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateGatewayResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aE<GatewayResponseType>(377935935, _omitFieldNames ? '' : 'responsetype',
        enumValues: GatewayResponseType.values)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGatewayResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGatewayResponseRequest copyWith(
          void Function(UpdateGatewayResponseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateGatewayResponseRequest))
          as UpdateGatewayResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateGatewayResponseRequest create() =>
      UpdateGatewayResponseRequest._();
  @$core.override
  UpdateGatewayResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateGatewayResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateGatewayResponseRequest>(create);
  static UpdateGatewayResponseRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(377935935)
  GatewayResponseType get responsetype => $_getN(1);
  @$pb.TagNumber(377935935)
  set responsetype(GatewayResponseType value) => $_setField(377935935, value);
  @$pb.TagNumber(377935935)
  $core.bool hasResponsetype() => $_has(1);
  @$pb.TagNumber(377935935)
  void clearResponsetype() => $_clearField(377935935);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateIntegrationRequest extends $pb.GeneratedMessage {
  factory UpdateIntegrationRequest({
    $core.String? httpmethod,
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateIntegrationRequest._();

  factory UpdateIntegrationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateIntegrationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateIntegrationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateIntegrationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateIntegrationRequest copyWith(
          void Function(UpdateIntegrationRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateIntegrationRequest))
          as UpdateIntegrationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateIntegrationRequest create() => UpdateIntegrationRequest._();
  @$core.override
  UpdateIntegrationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateIntegrationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateIntegrationRequest>(create);
  static UpdateIntegrationRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(2);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(2);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateIntegrationResponseRequest extends $pb.GeneratedMessage {
  factory UpdateIntegrationResponseRequest({
    $core.String? httpmethod,
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? statuscode,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (statuscode != null) result.statuscode = statuscode;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateIntegrationResponseRequest._();

  factory UpdateIntegrationResponseRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateIntegrationResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateIntegrationResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateIntegrationResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateIntegrationResponseRequest copyWith(
          void Function(UpdateIntegrationResponseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateIntegrationResponseRequest))
          as UpdateIntegrationResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateIntegrationResponseRequest create() =>
      UpdateIntegrationResponseRequest._();
  @$core.override
  UpdateIntegrationResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateIntegrationResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateIntegrationResponseRequest>(
          create);
  static UpdateIntegrationResponseRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(2);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(2);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(3);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(3);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(4);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(4);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateMethodRequest extends $pb.GeneratedMessage {
  factory UpdateMethodRequest({
    $core.String? httpmethod,
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateMethodRequest._();

  factory UpdateMethodRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateMethodRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateMethodRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMethodRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMethodRequest copyWith(void Function(UpdateMethodRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateMethodRequest))
          as UpdateMethodRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateMethodRequest create() => UpdateMethodRequest._();
  @$core.override
  UpdateMethodRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateMethodRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateMethodRequest>(create);
  static UpdateMethodRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(2);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(2);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(3);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(3);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateMethodResponseRequest extends $pb.GeneratedMessage {
  factory UpdateMethodResponseRequest({
    $core.String? httpmethod,
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? statuscode,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (statuscode != null) result.statuscode = statuscode;
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateMethodResponseRequest._();

  factory UpdateMethodResponseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateMethodResponseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateMethodResponseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(115276273, _omitFieldNames ? '' : 'httpmethod')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(299352223, _omitFieldNames ? '' : 'statuscode')
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMethodResponseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateMethodResponseRequest copyWith(
          void Function(UpdateMethodResponseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateMethodResponseRequest))
          as UpdateMethodResponseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateMethodResponseRequest create() =>
      UpdateMethodResponseRequest._();
  @$core.override
  UpdateMethodResponseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateMethodResponseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateMethodResponseRequest>(create);
  static UpdateMethodResponseRequest? _defaultInstance;

  @$pb.TagNumber(115276273)
  $core.String get httpmethod => $_getSZ(0);
  @$pb.TagNumber(115276273)
  set httpmethod($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115276273)
  $core.bool hasHttpmethod() => $_has(0);
  @$pb.TagNumber(115276273)
  void clearHttpmethod() => $_clearField(115276273);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);

  @$pb.TagNumber(299352223)
  $core.String get statuscode => $_getSZ(2);
  @$pb.TagNumber(299352223)
  set statuscode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(299352223)
  $core.bool hasStatuscode() => $_has(2);
  @$pb.TagNumber(299352223)
  void clearStatuscode() => $_clearField(299352223);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(3);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(3);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(4);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(4);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateModelRequest extends $pb.GeneratedMessage {
  factory UpdateModelRequest({
    $core.String? modelname,
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? restapiid,
  }) {
    final result = create();
    if (modelname != null) result.modelname = modelname;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateModelRequest._();

  factory UpdateModelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateModelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateModelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(176835330, _omitFieldNames ? '' : 'modelname')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateModelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateModelRequest copyWith(void Function(UpdateModelRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateModelRequest))
          as UpdateModelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateModelRequest create() => UpdateModelRequest._();
  @$core.override
  UpdateModelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateModelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateModelRequest>(create);
  static UpdateModelRequest? _defaultInstance;

  @$pb.TagNumber(176835330)
  $core.String get modelname => $_getSZ(0);
  @$pb.TagNumber(176835330)
  set modelname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(176835330)
  $core.bool hasModelname() => $_has(0);
  @$pb.TagNumber(176835330)
  void clearModelname() => $_clearField(176835330);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateRequestValidatorRequest extends $pb.GeneratedMessage {
  factory UpdateRequestValidatorRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? restapiid,
    $core.String? requestvalidatorid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (restapiid != null) result.restapiid = restapiid;
    if (requestvalidatorid != null)
      result.requestvalidatorid = requestvalidatorid;
    return result;
  }

  UpdateRequestValidatorRequest._();

  factory UpdateRequestValidatorRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRequestValidatorRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRequestValidatorRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..aOS(517546134, _omitFieldNames ? '' : 'requestvalidatorid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRequestValidatorRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRequestValidatorRequest copyWith(
          void Function(UpdateRequestValidatorRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateRequestValidatorRequest))
          as UpdateRequestValidatorRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRequestValidatorRequest create() =>
      UpdateRequestValidatorRequest._();
  @$core.override
  UpdateRequestValidatorRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRequestValidatorRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRequestValidatorRequest>(create);
  static UpdateRequestValidatorRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);

  @$pb.TagNumber(517546134)
  $core.String get requestvalidatorid => $_getSZ(2);
  @$pb.TagNumber(517546134)
  set requestvalidatorid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(517546134)
  $core.bool hasRequestvalidatorid() => $_has(2);
  @$pb.TagNumber(517546134)
  void clearRequestvalidatorid() => $_clearField(517546134);
}

class UpdateResourceRequest extends $pb.GeneratedMessage {
  factory UpdateResourceRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? resourceid,
    $core.String? restapiid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (resourceid != null) result.resourceid = resourceid;
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateResourceRequest._();

  factory UpdateResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateResourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(318922417, _omitFieldNames ? '' : 'resourceid')
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateResourceRequest copyWith(
          void Function(UpdateResourceRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateResourceRequest))
          as UpdateResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateResourceRequest create() => UpdateResourceRequest._();
  @$core.override
  UpdateResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateResourceRequest>(create);
  static UpdateResourceRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(318922417)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(318922417)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(318922417)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(318922417)
  void clearResourceid() => $_clearField(318922417);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateRestApiRequest extends $pb.GeneratedMessage {
  factory UpdateRestApiRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? restapiid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateRestApiRequest._();

  factory UpdateRestApiRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRestApiRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRestApiRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRestApiRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRestApiRequest copyWith(void Function(UpdateRestApiRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateRestApiRequest))
          as UpdateRestApiRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRestApiRequest create() => UpdateRestApiRequest._();
  @$core.override
  UpdateRestApiRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRestApiRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRestApiRequest>(create);
  static UpdateRestApiRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(1);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(1);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateStageRequest extends $pb.GeneratedMessage {
  factory UpdateStageRequest({
    $core.String? stagename,
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? restapiid,
  }) {
    final result = create();
    if (stagename != null) result.stagename = stagename;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (restapiid != null) result.restapiid = restapiid;
    return result;
  }

  UpdateStageRequest._();

  factory UpdateStageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateStageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateStageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(9563663, _omitFieldNames ? '' : 'stagename')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(383799833, _omitFieldNames ? '' : 'restapiid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateStageRequest copyWith(void Function(UpdateStageRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateStageRequest))
          as UpdateStageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateStageRequest create() => UpdateStageRequest._();
  @$core.override
  UpdateStageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateStageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateStageRequest>(create);
  static UpdateStageRequest? _defaultInstance;

  @$pb.TagNumber(9563663)
  $core.String get stagename => $_getSZ(0);
  @$pb.TagNumber(9563663)
  set stagename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9563663)
  $core.bool hasStagename() => $_has(0);
  @$pb.TagNumber(9563663)
  void clearStagename() => $_clearField(9563663);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);

  @$pb.TagNumber(383799833)
  $core.String get restapiid => $_getSZ(2);
  @$pb.TagNumber(383799833)
  set restapiid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(383799833)
  $core.bool hasRestapiid() => $_has(2);
  @$pb.TagNumber(383799833)
  void clearRestapiid() => $_clearField(383799833);
}

class UpdateUsagePlanRequest extends $pb.GeneratedMessage {
  factory UpdateUsagePlanRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? usageplanid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  UpdateUsagePlanRequest._();

  factory UpdateUsagePlanRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateUsagePlanRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateUsagePlanRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateUsagePlanRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateUsagePlanRequest copyWith(
          void Function(UpdateUsagePlanRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateUsagePlanRequest))
          as UpdateUsagePlanRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateUsagePlanRequest create() => UpdateUsagePlanRequest._();
  @$core.override
  UpdateUsagePlanRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateUsagePlanRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateUsagePlanRequest>(create);
  static UpdateUsagePlanRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(1);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(1);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class UpdateUsageRequest extends $pb.GeneratedMessage {
  factory UpdateUsageRequest({
    $core.Iterable<PatchOperation>? patchoperations,
    $core.String? keyid,
    $core.String? usageplanid,
  }) {
    final result = create();
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    if (keyid != null) result.keyid = keyid;
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  UpdateUsageRequest._();

  factory UpdateUsageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateUsageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateUsageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..aOS(479913282, _omitFieldNames ? '' : 'keyid')
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateUsageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateUsageRequest copyWith(void Function(UpdateUsageRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateUsageRequest))
          as UpdateUsageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateUsageRequest create() => UpdateUsageRequest._();
  @$core.override
  UpdateUsageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateUsageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateUsageRequest>(create);
  static UpdateUsageRequest? _defaultInstance;

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(0);

  @$pb.TagNumber(479913282)
  $core.String get keyid => $_getSZ(1);
  @$pb.TagNumber(479913282)
  set keyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(479913282)
  $core.bool hasKeyid() => $_has(1);
  @$pb.TagNumber(479913282)
  void clearKeyid() => $_clearField(479913282);

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(2);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(2);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class UpdateVpcLinkRequest extends $pb.GeneratedMessage {
  factory UpdateVpcLinkRequest({
    $core.String? vpclinkid,
    $core.Iterable<PatchOperation>? patchoperations,
  }) {
    final result = create();
    if (vpclinkid != null) result.vpclinkid = vpclinkid;
    if (patchoperations != null) result.patchoperations.addAll(patchoperations);
    return result;
  }

  UpdateVpcLinkRequest._();

  factory UpdateVpcLinkRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateVpcLinkRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateVpcLinkRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(27515846, _omitFieldNames ? '' : 'vpclinkid')
    ..pPM<PatchOperation>(201637420, _omitFieldNames ? '' : 'patchoperations',
        subBuilder: PatchOperation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateVpcLinkRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateVpcLinkRequest copyWith(void Function(UpdateVpcLinkRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateVpcLinkRequest))
          as UpdateVpcLinkRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateVpcLinkRequest create() => UpdateVpcLinkRequest._();
  @$core.override
  UpdateVpcLinkRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateVpcLinkRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateVpcLinkRequest>(create);
  static UpdateVpcLinkRequest? _defaultInstance;

  @$pb.TagNumber(27515846)
  $core.String get vpclinkid => $_getSZ(0);
  @$pb.TagNumber(27515846)
  set vpclinkid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(27515846)
  $core.bool hasVpclinkid() => $_has(0);
  @$pb.TagNumber(27515846)
  void clearVpclinkid() => $_clearField(27515846);

  @$pb.TagNumber(201637420)
  $pb.PbList<PatchOperation> get patchoperations => $_getList(1);
}

class Usage extends $pb.GeneratedMessage {
  factory Usage({
    $core.String? position,
    $core.String? startdate,
    $core.String? enddate,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? items,
    $core.String? usageplanid,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (startdate != null) result.startdate = startdate;
    if (enddate != null) result.enddate = enddate;
    if (items != null) result.items.addEntries(items);
    if (usageplanid != null) result.usageplanid = usageplanid;
    return result;
  }

  Usage._();

  factory Usage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Usage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Usage',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..aOS(364840732, _omitFieldNames ? '' : 'startdate')
    ..aOS(384831215, _omitFieldNames ? '' : 'enddate')
    ..m<$core.String, $core.String>(444150672, _omitFieldNames ? '' : 'items',
        entryClassName: 'Usage.ItemsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(509179991, _omitFieldNames ? '' : 'usageplanid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Usage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Usage copyWith(void Function(Usage) updates) =>
      super.copyWith((message) => updates(message as Usage)) as Usage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Usage create() => Usage._();
  @$core.override
  Usage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Usage getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Usage>(create);
  static Usage? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(364840732)
  $core.String get startdate => $_getSZ(1);
  @$pb.TagNumber(364840732)
  set startdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364840732)
  $core.bool hasStartdate() => $_has(1);
  @$pb.TagNumber(364840732)
  void clearStartdate() => $_clearField(364840732);

  @$pb.TagNumber(384831215)
  $core.String get enddate => $_getSZ(2);
  @$pb.TagNumber(384831215)
  set enddate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384831215)
  $core.bool hasEnddate() => $_has(2);
  @$pb.TagNumber(384831215)
  void clearEnddate() => $_clearField(384831215);

  @$pb.TagNumber(444150672)
  $pb.PbMap<$core.String, $core.String> get items => $_getMap(3);

  @$pb.TagNumber(509179991)
  $core.String get usageplanid => $_getSZ(4);
  @$pb.TagNumber(509179991)
  set usageplanid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(509179991)
  $core.bool hasUsageplanid() => $_has(4);
  @$pb.TagNumber(509179991)
  void clearUsageplanid() => $_clearField(509179991);
}

class UsagePlan extends $pb.GeneratedMessage {
  factory UsagePlan({
    $core.Iterable<ApiStage>? apistages,
    $core.String? name,
    QuotaSettings? quota,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    $core.String? id,
    ThrottleSettings? throttle,
    $core.String? productcode,
  }) {
    final result = create();
    if (apistages != null) result.apistages.addAll(apistages);
    if (name != null) result.name = name;
    if (quota != null) result.quota = quota;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (id != null) result.id = id;
    if (throttle != null) result.throttle = throttle;
    if (productcode != null) result.productcode = productcode;
    return result;
  }

  UsagePlan._();

  factory UsagePlan.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UsagePlan.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UsagePlan',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPM<ApiStage>(64558449, _omitFieldNames ? '' : 'apistages',
        subBuilder: ApiStage.create)
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOM<QuotaSettings>(243824012, _omitFieldNames ? '' : 'quota',
        subBuilder: QuotaSettings.create)
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'UsagePlan.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..aOM<ThrottleSettings>(395260638, _omitFieldNames ? '' : 'throttle',
        subBuilder: ThrottleSettings.create)
    ..aOS(533381226, _omitFieldNames ? '' : 'productcode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsagePlan clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsagePlan copyWith(void Function(UsagePlan) updates) =>
      super.copyWith((message) => updates(message as UsagePlan)) as UsagePlan;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UsagePlan create() => UsagePlan._();
  @$core.override
  UsagePlan createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UsagePlan getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UsagePlan>(create);
  static UsagePlan? _defaultInstance;

  @$pb.TagNumber(64558449)
  $pb.PbList<ApiStage> get apistages => $_getList(0);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(243824012)
  QuotaSettings get quota => $_getN(2);
  @$pb.TagNumber(243824012)
  set quota(QuotaSettings value) => $_setField(243824012, value);
  @$pb.TagNumber(243824012)
  $core.bool hasQuota() => $_has(2);
  @$pb.TagNumber(243824012)
  void clearQuota() => $_clearField(243824012);
  @$pb.TagNumber(243824012)
  QuotaSettings ensureQuota() => $_ensure(2);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(3);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(5);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(5, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(5);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);

  @$pb.TagNumber(395260638)
  ThrottleSettings get throttle => $_getN(6);
  @$pb.TagNumber(395260638)
  set throttle(ThrottleSettings value) => $_setField(395260638, value);
  @$pb.TagNumber(395260638)
  $core.bool hasThrottle() => $_has(6);
  @$pb.TagNumber(395260638)
  void clearThrottle() => $_clearField(395260638);
  @$pb.TagNumber(395260638)
  ThrottleSettings ensureThrottle() => $_ensure(6);

  @$pb.TagNumber(533381226)
  $core.String get productcode => $_getSZ(7);
  @$pb.TagNumber(533381226)
  set productcode($core.String value) => $_setString(7, value);
  @$pb.TagNumber(533381226)
  $core.bool hasProductcode() => $_has(7);
  @$pb.TagNumber(533381226)
  void clearProductcode() => $_clearField(533381226);
}

class UsagePlanKey extends $pb.GeneratedMessage {
  factory UsagePlanKey({
    $core.String? value,
    $core.String? name,
    $core.String? type,
    $core.String? id,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (id != null) result.id = id;
    return result;
  }

  UsagePlanKey._();

  factory UsagePlanKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UsagePlanKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UsagePlanKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(39769035, _omitFieldNames ? '' : 'value')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..aOS(287830350, _omitFieldNames ? '' : 'type')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsagePlanKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsagePlanKey copyWith(void Function(UsagePlanKey) updates) =>
      super.copyWith((message) => updates(message as UsagePlanKey))
          as UsagePlanKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UsagePlanKey create() => UsagePlanKey._();
  @$core.override
  UsagePlanKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UsagePlanKey getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UsagePlanKey>(create);
  static UsagePlanKey? _defaultInstance;

  @$pb.TagNumber(39769035)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(39769035)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39769035)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(39769035)
  void clearValue() => $_clearField(39769035);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(287830350)
  $core.String get type => $_getSZ(2);
  @$pb.TagNumber(287830350)
  set type($core.String value) => $_setString(2, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);
}

class UsagePlanKeys extends $pb.GeneratedMessage {
  factory UsagePlanKeys({
    $core.String? position,
    $core.Iterable<UsagePlanKey>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  UsagePlanKeys._();

  factory UsagePlanKeys.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UsagePlanKeys.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UsagePlanKeys',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<UsagePlanKey>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: UsagePlanKey.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsagePlanKeys clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsagePlanKeys copyWith(void Function(UsagePlanKeys) updates) =>
      super.copyWith((message) => updates(message as UsagePlanKeys))
          as UsagePlanKeys;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UsagePlanKeys create() => UsagePlanKeys._();
  @$core.override
  UsagePlanKeys createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UsagePlanKeys getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UsagePlanKeys>(create);
  static UsagePlanKeys? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<UsagePlanKey> get items => $_getList(1);
}

class UsagePlans extends $pb.GeneratedMessage {
  factory UsagePlans({
    $core.String? position,
    $core.Iterable<UsagePlan>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  UsagePlans._();

  factory UsagePlans.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UsagePlans.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UsagePlans',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<UsagePlan>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: UsagePlan.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsagePlans clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsagePlans copyWith(void Function(UsagePlans) updates) =>
      super.copyWith((message) => updates(message as UsagePlans)) as UsagePlans;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UsagePlans create() => UsagePlans._();
  @$core.override
  UsagePlans createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UsagePlans getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UsagePlans>(create);
  static UsagePlans? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<UsagePlan> get items => $_getList(1);
}

class VpcLink extends $pb.GeneratedMessage {
  factory VpcLink({
    $core.Iterable<$core.String>? targetarns,
    $core.String? name,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.String? description,
    $core.String? id,
    VpcLinkStatus? status,
    $core.String? statusmessage,
  }) {
    final result = create();
    if (targetarns != null) result.targetarns.addAll(targetarns);
    if (name != null) result.name = name;
    if (tags != null) result.tags.addEntries(tags);
    if (description != null) result.description = description;
    if (id != null) result.id = id;
    if (status != null) result.status = status;
    if (statusmessage != null) result.statusmessage = statusmessage;
    return result;
  }

  VpcLink._();

  factory VpcLink.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VpcLink.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VpcLink',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..pPS(46319317, _omitFieldNames ? '' : 'targetarns')
    ..aOS(221887975, _omitFieldNames ? '' : 'name')
    ..m<$core.String, $core.String>(337046433, _omitFieldNames ? '' : 'tags',
        entryClassName: 'VpcLink.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('apigateway'))
    ..aOS(342834026, _omitFieldNames ? '' : 'description')
    ..aOS(389573345, _omitFieldNames ? '' : 'id')
    ..aE<VpcLinkStatus>(441153520, _omitFieldNames ? '' : 'status',
        enumValues: VpcLinkStatus.values)
    ..aOS(474462255, _omitFieldNames ? '' : 'statusmessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VpcLink clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VpcLink copyWith(void Function(VpcLink) updates) =>
      super.copyWith((message) => updates(message as VpcLink)) as VpcLink;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VpcLink create() => VpcLink._();
  @$core.override
  VpcLink createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VpcLink getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VpcLink>(create);
  static VpcLink? _defaultInstance;

  @$pb.TagNumber(46319317)
  $pb.PbList<$core.String> get targetarns => $_getList(0);

  @$pb.TagNumber(221887975)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(221887975)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(221887975)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(221887975)
  void clearName() => $_clearField(221887975);

  @$pb.TagNumber(337046433)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(2);

  @$pb.TagNumber(342834026)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(342834026)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(342834026)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(342834026)
  void clearDescription() => $_clearField(342834026);

  @$pb.TagNumber(389573345)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(389573345)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(389573345)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(389573345)
  void clearId() => $_clearField(389573345);

  @$pb.TagNumber(441153520)
  VpcLinkStatus get status => $_getN(5);
  @$pb.TagNumber(441153520)
  set status(VpcLinkStatus value) => $_setField(441153520, value);
  @$pb.TagNumber(441153520)
  $core.bool hasStatus() => $_has(5);
  @$pb.TagNumber(441153520)
  void clearStatus() => $_clearField(441153520);

  @$pb.TagNumber(474462255)
  $core.String get statusmessage => $_getSZ(6);
  @$pb.TagNumber(474462255)
  set statusmessage($core.String value) => $_setString(6, value);
  @$pb.TagNumber(474462255)
  $core.bool hasStatusmessage() => $_has(6);
  @$pb.TagNumber(474462255)
  void clearStatusmessage() => $_clearField(474462255);
}

class VpcLinks extends $pb.GeneratedMessage {
  factory VpcLinks({
    $core.String? position,
    $core.Iterable<VpcLink>? items,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (items != null) result.items.addAll(items);
    return result;
  }

  VpcLinks._();

  factory VpcLinks.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VpcLinks.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VpcLinks',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'apigateway'),
      createEmptyInstance: create)
    ..aOS(323964427, _omitFieldNames ? '' : 'position')
    ..pPM<VpcLink>(444150672, _omitFieldNames ? '' : 'items',
        subBuilder: VpcLink.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VpcLinks clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VpcLinks copyWith(void Function(VpcLinks) updates) =>
      super.copyWith((message) => updates(message as VpcLinks)) as VpcLinks;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VpcLinks create() => VpcLinks._();
  @$core.override
  VpcLinks createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VpcLinks getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VpcLinks>(create);
  static VpcLinks? _defaultInstance;

  @$pb.TagNumber(323964427)
  $core.String get position => $_getSZ(0);
  @$pb.TagNumber(323964427)
  set position($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323964427)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(323964427)
  void clearPosition() => $_clearField(323964427);

  @$pb.TagNumber(444150672)
  $pb.PbList<VpcLink> get items => $_getList(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
