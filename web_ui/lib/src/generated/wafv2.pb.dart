// This is a generated file - do not edit.
//
// Generated from wafv2.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'wafv2.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'wafv2.pbenum.dart';

class APIKeySummary extends $pb.GeneratedMessage {
  factory APIKeySummary({
    $core.Iterable<$core.String>? tokendomains,
    $core.String? creationtimestamp,
    $core.String? apikey,
    $core.int? version,
  }) {
    final result = create();
    if (tokendomains != null) result.tokendomains.addAll(tokendomains);
    if (creationtimestamp != null) result.creationtimestamp = creationtimestamp;
    if (apikey != null) result.apikey = apikey;
    if (version != null) result.version = version;
    return result;
  }

  APIKeySummary._();

  factory APIKeySummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory APIKeySummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'APIKeySummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(11638560, _omitFieldNames ? '' : 'tokendomains')
    ..aOS(24480293, _omitFieldNames ? '' : 'creationtimestamp')
    ..aOS(274818239, _omitFieldNames ? '' : 'apikey')
    ..aI(500028728, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  APIKeySummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  APIKeySummary copyWith(void Function(APIKeySummary) updates) =>
      super.copyWith((message) => updates(message as APIKeySummary))
          as APIKeySummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static APIKeySummary create() => APIKeySummary._();
  @$core.override
  APIKeySummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static APIKeySummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<APIKeySummary>(create);
  static APIKeySummary? _defaultInstance;

  @$pb.TagNumber(11638560)
  $pb.PbList<$core.String> get tokendomains => $_getList(0);

  @$pb.TagNumber(24480293)
  $core.String get creationtimestamp => $_getSZ(1);
  @$pb.TagNumber(24480293)
  set creationtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(24480293)
  $core.bool hasCreationtimestamp() => $_has(1);
  @$pb.TagNumber(24480293)
  void clearCreationtimestamp() => $_clearField(24480293);

  @$pb.TagNumber(274818239)
  $core.String get apikey => $_getSZ(2);
  @$pb.TagNumber(274818239)
  set apikey($core.String value) => $_setString(2, value);
  @$pb.TagNumber(274818239)
  $core.bool hasApikey() => $_has(2);
  @$pb.TagNumber(274818239)
  void clearApikey() => $_clearField(274818239);

  @$pb.TagNumber(500028728)
  $core.int get version => $_getIZ(3);
  @$pb.TagNumber(500028728)
  set version($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(3);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);
}

class AWSManagedRulesACFPRuleSet extends $pb.GeneratedMessage {
  factory AWSManagedRulesACFPRuleSet({
    $core.String? registrationpagepath,
    $core.bool? enableregexinpath,
    ResponseInspection? responseinspection,
    $core.String? creationpath,
    RequestInspectionACFP? requestinspection,
  }) {
    final result = create();
    if (registrationpagepath != null)
      result.registrationpagepath = registrationpagepath;
    if (enableregexinpath != null) result.enableregexinpath = enableregexinpath;
    if (responseinspection != null)
      result.responseinspection = responseinspection;
    if (creationpath != null) result.creationpath = creationpath;
    if (requestinspection != null) result.requestinspection = requestinspection;
    return result;
  }

  AWSManagedRulesACFPRuleSet._();

  factory AWSManagedRulesACFPRuleSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AWSManagedRulesACFPRuleSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AWSManagedRulesACFPRuleSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(26308901, _omitFieldNames ? '' : 'registrationpagepath')
    ..aOB(36760922, _omitFieldNames ? '' : 'enableregexinpath')
    ..aOM<ResponseInspection>(
        267610877, _omitFieldNames ? '' : 'responseinspection',
        subBuilder: ResponseInspection.create)
    ..aOS(349715284, _omitFieldNames ? '' : 'creationpath')
    ..aOM<RequestInspectionACFP>(
        389146793, _omitFieldNames ? '' : 'requestinspection',
        subBuilder: RequestInspectionACFP.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AWSManagedRulesACFPRuleSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AWSManagedRulesACFPRuleSet copyWith(
          void Function(AWSManagedRulesACFPRuleSet) updates) =>
      super.copyWith(
              (message) => updates(message as AWSManagedRulesACFPRuleSet))
          as AWSManagedRulesACFPRuleSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AWSManagedRulesACFPRuleSet create() => AWSManagedRulesACFPRuleSet._();
  @$core.override
  AWSManagedRulesACFPRuleSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AWSManagedRulesACFPRuleSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AWSManagedRulesACFPRuleSet>(create);
  static AWSManagedRulesACFPRuleSet? _defaultInstance;

  @$pb.TagNumber(26308901)
  $core.String get registrationpagepath => $_getSZ(0);
  @$pb.TagNumber(26308901)
  set registrationpagepath($core.String value) => $_setString(0, value);
  @$pb.TagNumber(26308901)
  $core.bool hasRegistrationpagepath() => $_has(0);
  @$pb.TagNumber(26308901)
  void clearRegistrationpagepath() => $_clearField(26308901);

  @$pb.TagNumber(36760922)
  $core.bool get enableregexinpath => $_getBF(1);
  @$pb.TagNumber(36760922)
  set enableregexinpath($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(36760922)
  $core.bool hasEnableregexinpath() => $_has(1);
  @$pb.TagNumber(36760922)
  void clearEnableregexinpath() => $_clearField(36760922);

  @$pb.TagNumber(267610877)
  ResponseInspection get responseinspection => $_getN(2);
  @$pb.TagNumber(267610877)
  set responseinspection(ResponseInspection value) =>
      $_setField(267610877, value);
  @$pb.TagNumber(267610877)
  $core.bool hasResponseinspection() => $_has(2);
  @$pb.TagNumber(267610877)
  void clearResponseinspection() => $_clearField(267610877);
  @$pb.TagNumber(267610877)
  ResponseInspection ensureResponseinspection() => $_ensure(2);

  @$pb.TagNumber(349715284)
  $core.String get creationpath => $_getSZ(3);
  @$pb.TagNumber(349715284)
  set creationpath($core.String value) => $_setString(3, value);
  @$pb.TagNumber(349715284)
  $core.bool hasCreationpath() => $_has(3);
  @$pb.TagNumber(349715284)
  void clearCreationpath() => $_clearField(349715284);

  @$pb.TagNumber(389146793)
  RequestInspectionACFP get requestinspection => $_getN(4);
  @$pb.TagNumber(389146793)
  set requestinspection(RequestInspectionACFP value) =>
      $_setField(389146793, value);
  @$pb.TagNumber(389146793)
  $core.bool hasRequestinspection() => $_has(4);
  @$pb.TagNumber(389146793)
  void clearRequestinspection() => $_clearField(389146793);
  @$pb.TagNumber(389146793)
  RequestInspectionACFP ensureRequestinspection() => $_ensure(4);
}

class AWSManagedRulesATPRuleSet extends $pb.GeneratedMessage {
  factory AWSManagedRulesATPRuleSet({
    $core.bool? enableregexinpath,
    $core.String? loginpath,
    ResponseInspection? responseinspection,
    RequestInspection? requestinspection,
  }) {
    final result = create();
    if (enableregexinpath != null) result.enableregexinpath = enableregexinpath;
    if (loginpath != null) result.loginpath = loginpath;
    if (responseinspection != null)
      result.responseinspection = responseinspection;
    if (requestinspection != null) result.requestinspection = requestinspection;
    return result;
  }

  AWSManagedRulesATPRuleSet._();

  factory AWSManagedRulesATPRuleSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AWSManagedRulesATPRuleSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AWSManagedRulesATPRuleSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOB(36760922, _omitFieldNames ? '' : 'enableregexinpath')
    ..aOS(128281874, _omitFieldNames ? '' : 'loginpath')
    ..aOM<ResponseInspection>(
        267610877, _omitFieldNames ? '' : 'responseinspection',
        subBuilder: ResponseInspection.create)
    ..aOM<RequestInspection>(
        389146793, _omitFieldNames ? '' : 'requestinspection',
        subBuilder: RequestInspection.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AWSManagedRulesATPRuleSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AWSManagedRulesATPRuleSet copyWith(
          void Function(AWSManagedRulesATPRuleSet) updates) =>
      super.copyWith((message) => updates(message as AWSManagedRulesATPRuleSet))
          as AWSManagedRulesATPRuleSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AWSManagedRulesATPRuleSet create() => AWSManagedRulesATPRuleSet._();
  @$core.override
  AWSManagedRulesATPRuleSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AWSManagedRulesATPRuleSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AWSManagedRulesATPRuleSet>(create);
  static AWSManagedRulesATPRuleSet? _defaultInstance;

  @$pb.TagNumber(36760922)
  $core.bool get enableregexinpath => $_getBF(0);
  @$pb.TagNumber(36760922)
  set enableregexinpath($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(36760922)
  $core.bool hasEnableregexinpath() => $_has(0);
  @$pb.TagNumber(36760922)
  void clearEnableregexinpath() => $_clearField(36760922);

  @$pb.TagNumber(128281874)
  $core.String get loginpath => $_getSZ(1);
  @$pb.TagNumber(128281874)
  set loginpath($core.String value) => $_setString(1, value);
  @$pb.TagNumber(128281874)
  $core.bool hasLoginpath() => $_has(1);
  @$pb.TagNumber(128281874)
  void clearLoginpath() => $_clearField(128281874);

  @$pb.TagNumber(267610877)
  ResponseInspection get responseinspection => $_getN(2);
  @$pb.TagNumber(267610877)
  set responseinspection(ResponseInspection value) =>
      $_setField(267610877, value);
  @$pb.TagNumber(267610877)
  $core.bool hasResponseinspection() => $_has(2);
  @$pb.TagNumber(267610877)
  void clearResponseinspection() => $_clearField(267610877);
  @$pb.TagNumber(267610877)
  ResponseInspection ensureResponseinspection() => $_ensure(2);

  @$pb.TagNumber(389146793)
  RequestInspection get requestinspection => $_getN(3);
  @$pb.TagNumber(389146793)
  set requestinspection(RequestInspection value) =>
      $_setField(389146793, value);
  @$pb.TagNumber(389146793)
  $core.bool hasRequestinspection() => $_has(3);
  @$pb.TagNumber(389146793)
  void clearRequestinspection() => $_clearField(389146793);
  @$pb.TagNumber(389146793)
  RequestInspection ensureRequestinspection() => $_ensure(3);
}

class AWSManagedRulesAntiDDoSRuleSet extends $pb.GeneratedMessage {
  factory AWSManagedRulesAntiDDoSRuleSet({
    ClientSideActionConfig? clientsideactionconfig,
    SensitivityToAct? sensitivitytoblock,
  }) {
    final result = create();
    if (clientsideactionconfig != null)
      result.clientsideactionconfig = clientsideactionconfig;
    if (sensitivitytoblock != null)
      result.sensitivitytoblock = sensitivitytoblock;
    return result;
  }

  AWSManagedRulesAntiDDoSRuleSet._();

  factory AWSManagedRulesAntiDDoSRuleSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AWSManagedRulesAntiDDoSRuleSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AWSManagedRulesAntiDDoSRuleSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<ClientSideActionConfig>(
        55744822, _omitFieldNames ? '' : 'clientsideactionconfig',
        subBuilder: ClientSideActionConfig.create)
    ..aE<SensitivityToAct>(
        531809347, _omitFieldNames ? '' : 'sensitivitytoblock',
        enumValues: SensitivityToAct.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AWSManagedRulesAntiDDoSRuleSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AWSManagedRulesAntiDDoSRuleSet copyWith(
          void Function(AWSManagedRulesAntiDDoSRuleSet) updates) =>
      super.copyWith(
              (message) => updates(message as AWSManagedRulesAntiDDoSRuleSet))
          as AWSManagedRulesAntiDDoSRuleSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AWSManagedRulesAntiDDoSRuleSet create() =>
      AWSManagedRulesAntiDDoSRuleSet._();
  @$core.override
  AWSManagedRulesAntiDDoSRuleSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AWSManagedRulesAntiDDoSRuleSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AWSManagedRulesAntiDDoSRuleSet>(create);
  static AWSManagedRulesAntiDDoSRuleSet? _defaultInstance;

  @$pb.TagNumber(55744822)
  ClientSideActionConfig get clientsideactionconfig => $_getN(0);
  @$pb.TagNumber(55744822)
  set clientsideactionconfig(ClientSideActionConfig value) =>
      $_setField(55744822, value);
  @$pb.TagNumber(55744822)
  $core.bool hasClientsideactionconfig() => $_has(0);
  @$pb.TagNumber(55744822)
  void clearClientsideactionconfig() => $_clearField(55744822);
  @$pb.TagNumber(55744822)
  ClientSideActionConfig ensureClientsideactionconfig() => $_ensure(0);

  @$pb.TagNumber(531809347)
  SensitivityToAct get sensitivitytoblock => $_getN(1);
  @$pb.TagNumber(531809347)
  set sensitivitytoblock(SensitivityToAct value) =>
      $_setField(531809347, value);
  @$pb.TagNumber(531809347)
  $core.bool hasSensitivitytoblock() => $_has(1);
  @$pb.TagNumber(531809347)
  void clearSensitivitytoblock() => $_clearField(531809347);
}

class AWSManagedRulesBotControlRuleSet extends $pb.GeneratedMessage {
  factory AWSManagedRulesBotControlRuleSet({
    InspectionLevel? inspectionlevel,
    $core.bool? enablemachinelearning,
  }) {
    final result = create();
    if (inspectionlevel != null) result.inspectionlevel = inspectionlevel;
    if (enablemachinelearning != null)
      result.enablemachinelearning = enablemachinelearning;
    return result;
  }

  AWSManagedRulesBotControlRuleSet._();

  factory AWSManagedRulesBotControlRuleSet.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AWSManagedRulesBotControlRuleSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AWSManagedRulesBotControlRuleSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<InspectionLevel>(55668132, _omitFieldNames ? '' : 'inspectionlevel',
        enumValues: InspectionLevel.values)
    ..aOB(99141846, _omitFieldNames ? '' : 'enablemachinelearning')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AWSManagedRulesBotControlRuleSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AWSManagedRulesBotControlRuleSet copyWith(
          void Function(AWSManagedRulesBotControlRuleSet) updates) =>
      super.copyWith(
              (message) => updates(message as AWSManagedRulesBotControlRuleSet))
          as AWSManagedRulesBotControlRuleSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AWSManagedRulesBotControlRuleSet create() =>
      AWSManagedRulesBotControlRuleSet._();
  @$core.override
  AWSManagedRulesBotControlRuleSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AWSManagedRulesBotControlRuleSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AWSManagedRulesBotControlRuleSet>(
          create);
  static AWSManagedRulesBotControlRuleSet? _defaultInstance;

  @$pb.TagNumber(55668132)
  InspectionLevel get inspectionlevel => $_getN(0);
  @$pb.TagNumber(55668132)
  set inspectionlevel(InspectionLevel value) => $_setField(55668132, value);
  @$pb.TagNumber(55668132)
  $core.bool hasInspectionlevel() => $_has(0);
  @$pb.TagNumber(55668132)
  void clearInspectionlevel() => $_clearField(55668132);

  @$pb.TagNumber(99141846)
  $core.bool get enablemachinelearning => $_getBF(1);
  @$pb.TagNumber(99141846)
  set enablemachinelearning($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(99141846)
  $core.bool hasEnablemachinelearning() => $_has(1);
  @$pb.TagNumber(99141846)
  void clearEnablemachinelearning() => $_clearField(99141846);
}

class ActionCondition extends $pb.GeneratedMessage {
  factory ActionCondition({
    ActionValue? action,
  }) {
    final result = create();
    if (action != null) result.action = action;
    return result;
  }

  ActionCondition._();

  factory ActionCondition.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActionCondition.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActionCondition',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<ActionValue>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ActionValue.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActionCondition clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActionCondition copyWith(void Function(ActionCondition) updates) =>
      super.copyWith((message) => updates(message as ActionCondition))
          as ActionCondition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActionCondition create() => ActionCondition._();
  @$core.override
  ActionCondition createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActionCondition getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActionCondition>(create);
  static ActionCondition? _defaultInstance;

  @$pb.TagNumber(175614240)
  ActionValue get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(ActionValue value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
}

class AddressField extends $pb.GeneratedMessage {
  factory AddressField({
    $core.String? identifier,
  }) {
    final result = create();
    if (identifier != null) result.identifier = identifier;
    return result;
  }

  AddressField._();

  factory AddressField.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AddressField.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddressField',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(41865311, _omitFieldNames ? '' : 'identifier')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddressField clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddressField copyWith(void Function(AddressField) updates) =>
      super.copyWith((message) => updates(message as AddressField))
          as AddressField;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddressField create() => AddressField._();
  @$core.override
  AddressField createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AddressField getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddressField>(create);
  static AddressField? _defaultInstance;

  @$pb.TagNumber(41865311)
  $core.String get identifier => $_getSZ(0);
  @$pb.TagNumber(41865311)
  set identifier($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41865311)
  $core.bool hasIdentifier() => $_has(0);
  @$pb.TagNumber(41865311)
  void clearIdentifier() => $_clearField(41865311);
}

class All extends $pb.GeneratedMessage {
  factory All() => create();

  All._();

  factory All.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory All.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'All',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  All clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  All copyWith(void Function(All) updates) =>
      super.copyWith((message) => updates(message as All)) as All;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static All create() => All._();
  @$core.override
  All createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static All getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<All>(create);
  static All? _defaultInstance;
}

class AllQueryArguments extends $pb.GeneratedMessage {
  factory AllQueryArguments() => create();

  AllQueryArguments._();

  factory AllQueryArguments.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AllQueryArguments.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AllQueryArguments',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AllQueryArguments clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AllQueryArguments copyWith(void Function(AllQueryArguments) updates) =>
      super.copyWith((message) => updates(message as AllQueryArguments))
          as AllQueryArguments;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AllQueryArguments create() => AllQueryArguments._();
  @$core.override
  AllQueryArguments createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AllQueryArguments getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AllQueryArguments>(create);
  static AllQueryArguments? _defaultInstance;
}

class AllowAction extends $pb.GeneratedMessage {
  factory AllowAction({
    CustomRequestHandling? customrequesthandling,
  }) {
    final result = create();
    if (customrequesthandling != null)
      result.customrequesthandling = customrequesthandling;
    return result;
  }

  AllowAction._();

  factory AllowAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AllowAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AllowAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<CustomRequestHandling>(
        14845503, _omitFieldNames ? '' : 'customrequesthandling',
        subBuilder: CustomRequestHandling.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AllowAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AllowAction copyWith(void Function(AllowAction) updates) =>
      super.copyWith((message) => updates(message as AllowAction))
          as AllowAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AllowAction create() => AllowAction._();
  @$core.override
  AllowAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AllowAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AllowAction>(create);
  static AllowAction? _defaultInstance;

  @$pb.TagNumber(14845503)
  CustomRequestHandling get customrequesthandling => $_getN(0);
  @$pb.TagNumber(14845503)
  set customrequesthandling(CustomRequestHandling value) =>
      $_setField(14845503, value);
  @$pb.TagNumber(14845503)
  $core.bool hasCustomrequesthandling() => $_has(0);
  @$pb.TagNumber(14845503)
  void clearCustomrequesthandling() => $_clearField(14845503);
  @$pb.TagNumber(14845503)
  CustomRequestHandling ensureCustomrequesthandling() => $_ensure(0);
}

class AndStatement extends $pb.GeneratedMessage {
  factory AndStatement({
    $core.Iterable<Statement>? statements,
  }) {
    final result = create();
    if (statements != null) result.statements.addAll(statements);
    return result;
  }

  AndStatement._();

  factory AndStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AndStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AndStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<Statement>(488352288, _omitFieldNames ? '' : 'statements',
        subBuilder: Statement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AndStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AndStatement copyWith(void Function(AndStatement) updates) =>
      super.copyWith((message) => updates(message as AndStatement))
          as AndStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AndStatement create() => AndStatement._();
  @$core.override
  AndStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AndStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AndStatement>(create);
  static AndStatement? _defaultInstance;

  @$pb.TagNumber(488352288)
  $pb.PbList<Statement> get statements => $_getList(0);
}

class ApplicationAttribute extends $pb.GeneratedMessage {
  factory ApplicationAttribute({
    $core.Iterable<$core.String>? values,
    $core.String? name,
  }) {
    final result = create();
    if (values != null) result.values.addAll(values);
    if (name != null) result.name = name;
    return result;
  }

  ApplicationAttribute._();

  factory ApplicationAttribute.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ApplicationAttribute.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApplicationAttribute',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(223158876, _omitFieldNames ? '' : 'values')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApplicationAttribute clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApplicationAttribute copyWith(void Function(ApplicationAttribute) updates) =>
      super.copyWith((message) => updates(message as ApplicationAttribute))
          as ApplicationAttribute;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ApplicationAttribute create() => ApplicationAttribute._();
  @$core.override
  ApplicationAttribute createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ApplicationAttribute getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ApplicationAttribute>(create);
  static ApplicationAttribute? _defaultInstance;

  @$pb.TagNumber(223158876)
  $pb.PbList<$core.String> get values => $_getList(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ApplicationConfig extends $pb.GeneratedMessage {
  factory ApplicationConfig({
    $core.Iterable<ApplicationAttribute>? attributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addAll(attributes);
    return result;
  }

  ApplicationConfig._();

  factory ApplicationConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ApplicationConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApplicationConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<ApplicationAttribute>(209638581, _omitFieldNames ? '' : 'attributes',
        subBuilder: ApplicationAttribute.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApplicationConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApplicationConfig copyWith(void Function(ApplicationConfig) updates) =>
      super.copyWith((message) => updates(message as ApplicationConfig))
          as ApplicationConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ApplicationConfig create() => ApplicationConfig._();
  @$core.override
  ApplicationConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ApplicationConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ApplicationConfig>(create);
  static ApplicationConfig? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbList<ApplicationAttribute> get attributes => $_getList(0);
}

class AsnMatchStatement extends $pb.GeneratedMessage {
  factory AsnMatchStatement({
    ForwardedIPConfig? forwardedipconfig,
    $core.Iterable<$fixnum.Int64>? asnlist,
  }) {
    final result = create();
    if (forwardedipconfig != null) result.forwardedipconfig = forwardedipconfig;
    if (asnlist != null) result.asnlist.addAll(asnlist);
    return result;
  }

  AsnMatchStatement._();

  factory AsnMatchStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AsnMatchStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AsnMatchStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<ForwardedIPConfig>(
        259846797, _omitFieldNames ? '' : 'forwardedipconfig',
        subBuilder: ForwardedIPConfig.create)
    ..p<$fixnum.Int64>(
        319333670, _omitFieldNames ? '' : 'asnlist', $pb.PbFieldType.K6)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AsnMatchStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AsnMatchStatement copyWith(void Function(AsnMatchStatement) updates) =>
      super.copyWith((message) => updates(message as AsnMatchStatement))
          as AsnMatchStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AsnMatchStatement create() => AsnMatchStatement._();
  @$core.override
  AsnMatchStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AsnMatchStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AsnMatchStatement>(create);
  static AsnMatchStatement? _defaultInstance;

  @$pb.TagNumber(259846797)
  ForwardedIPConfig get forwardedipconfig => $_getN(0);
  @$pb.TagNumber(259846797)
  set forwardedipconfig(ForwardedIPConfig value) =>
      $_setField(259846797, value);
  @$pb.TagNumber(259846797)
  $core.bool hasForwardedipconfig() => $_has(0);
  @$pb.TagNumber(259846797)
  void clearForwardedipconfig() => $_clearField(259846797);
  @$pb.TagNumber(259846797)
  ForwardedIPConfig ensureForwardedipconfig() => $_ensure(0);

  @$pb.TagNumber(319333670)
  $pb.PbList<$fixnum.Int64> get asnlist => $_getList(1);
}

class AssociateWebACLRequest extends $pb.GeneratedMessage {
  factory AssociateWebACLRequest({
    $core.String? webaclarn,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (webaclarn != null) result.webaclarn = webaclarn;
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  AssociateWebACLRequest._();

  factory AssociateWebACLRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssociateWebACLRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssociateWebACLRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(82506659, _omitFieldNames ? '' : 'webaclarn')
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociateWebACLRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociateWebACLRequest copyWith(
          void Function(AssociateWebACLRequest) updates) =>
      super.copyWith((message) => updates(message as AssociateWebACLRequest))
          as AssociateWebACLRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssociateWebACLRequest create() => AssociateWebACLRequest._();
  @$core.override
  AssociateWebACLRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssociateWebACLRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssociateWebACLRequest>(create);
  static AssociateWebACLRequest? _defaultInstance;

  @$pb.TagNumber(82506659)
  $core.String get webaclarn => $_getSZ(0);
  @$pb.TagNumber(82506659)
  set webaclarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82506659)
  $core.bool hasWebaclarn() => $_has(0);
  @$pb.TagNumber(82506659)
  void clearWebaclarn() => $_clearField(82506659);

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class AssociateWebACLResponse extends $pb.GeneratedMessage {
  factory AssociateWebACLResponse() => create();

  AssociateWebACLResponse._();

  factory AssociateWebACLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssociateWebACLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssociateWebACLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociateWebACLResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociateWebACLResponse copyWith(
          void Function(AssociateWebACLResponse) updates) =>
      super.copyWith((message) => updates(message as AssociateWebACLResponse))
          as AssociateWebACLResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssociateWebACLResponse create() => AssociateWebACLResponse._();
  @$core.override
  AssociateWebACLResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssociateWebACLResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssociateWebACLResponse>(create);
  static AssociateWebACLResponse? _defaultInstance;
}

class AssociationConfig extends $pb.GeneratedMessage {
  factory AssociationConfig({
    $core.Iterable<
            $core
            .MapEntry<$core.String, RequestBodyAssociatedResourceTypeConfig>>?
        requestbody,
  }) {
    final result = create();
    if (requestbody != null) result.requestbody.addEntries(requestbody);
    return result;
  }

  AssociationConfig._();

  factory AssociationConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssociationConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssociationConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..m<$core.String, RequestBodyAssociatedResourceTypeConfig>(
        479920951, _omitFieldNames ? '' : 'requestbody',
        entryClassName: 'AssociationConfig.RequestbodyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: RequestBodyAssociatedResourceTypeConfig.create,
        valueDefaultOrMaker: RequestBodyAssociatedResourceTypeConfig.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociationConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociationConfig copyWith(void Function(AssociationConfig) updates) =>
      super.copyWith((message) => updates(message as AssociationConfig))
          as AssociationConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssociationConfig create() => AssociationConfig._();
  @$core.override
  AssociationConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssociationConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssociationConfig>(create);
  static AssociationConfig? _defaultInstance;

  @$pb.TagNumber(479920951)
  $pb.PbMap<$core.String, RequestBodyAssociatedResourceTypeConfig>
      get requestbody => $_getMap(0);
}

class BlockAction extends $pb.GeneratedMessage {
  factory BlockAction({
    CustomResponse? customresponse,
  }) {
    final result = create();
    if (customresponse != null) result.customresponse = customresponse;
    return result;
  }

  BlockAction._();

  factory BlockAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BlockAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BlockAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<CustomResponse>(137441042, _omitFieldNames ? '' : 'customresponse',
        subBuilder: CustomResponse.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BlockAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BlockAction copyWith(void Function(BlockAction) updates) =>
      super.copyWith((message) => updates(message as BlockAction))
          as BlockAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BlockAction create() => BlockAction._();
  @$core.override
  BlockAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BlockAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BlockAction>(create);
  static BlockAction? _defaultInstance;

  @$pb.TagNumber(137441042)
  CustomResponse get customresponse => $_getN(0);
  @$pb.TagNumber(137441042)
  set customresponse(CustomResponse value) => $_setField(137441042, value);
  @$pb.TagNumber(137441042)
  $core.bool hasCustomresponse() => $_has(0);
  @$pb.TagNumber(137441042)
  void clearCustomresponse() => $_clearField(137441042);
  @$pb.TagNumber(137441042)
  CustomResponse ensureCustomresponse() => $_ensure(0);
}

class Body extends $pb.GeneratedMessage {
  factory Body({
    OversizeHandling? oversizehandling,
  }) {
    final result = create();
    if (oversizehandling != null) result.oversizehandling = oversizehandling;
    return result;
  }

  Body._();

  factory Body.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Body.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Body',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<OversizeHandling>(139375132, _omitFieldNames ? '' : 'oversizehandling',
        enumValues: OversizeHandling.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Body clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Body copyWith(void Function(Body) updates) =>
      super.copyWith((message) => updates(message as Body)) as Body;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Body create() => Body._();
  @$core.override
  Body createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Body getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Body>(create);
  static Body? _defaultInstance;

  @$pb.TagNumber(139375132)
  OversizeHandling get oversizehandling => $_getN(0);
  @$pb.TagNumber(139375132)
  set oversizehandling(OversizeHandling value) => $_setField(139375132, value);
  @$pb.TagNumber(139375132)
  $core.bool hasOversizehandling() => $_has(0);
  @$pb.TagNumber(139375132)
  void clearOversizehandling() => $_clearField(139375132);
}

class BotStatistics extends $pb.GeneratedMessage {
  factory BotStatistics({
    $core.String? botname,
    $fixnum.Int64? requestcount,
    $core.double? percentage,
  }) {
    final result = create();
    if (botname != null) result.botname = botname;
    if (requestcount != null) result.requestcount = requestcount;
    if (percentage != null) result.percentage = percentage;
    return result;
  }

  BotStatistics._();

  factory BotStatistics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BotStatistics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BotStatistics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(88172446, _omitFieldNames ? '' : 'botname')
    ..aInt64(288519674, _omitFieldNames ? '' : 'requestcount')
    ..aD(341153238, _omitFieldNames ? '' : 'percentage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BotStatistics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BotStatistics copyWith(void Function(BotStatistics) updates) =>
      super.copyWith((message) => updates(message as BotStatistics))
          as BotStatistics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BotStatistics create() => BotStatistics._();
  @$core.override
  BotStatistics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BotStatistics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BotStatistics>(create);
  static BotStatistics? _defaultInstance;

  @$pb.TagNumber(88172446)
  $core.String get botname => $_getSZ(0);
  @$pb.TagNumber(88172446)
  set botname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88172446)
  $core.bool hasBotname() => $_has(0);
  @$pb.TagNumber(88172446)
  void clearBotname() => $_clearField(88172446);

  @$pb.TagNumber(288519674)
  $fixnum.Int64 get requestcount => $_getI64(1);
  @$pb.TagNumber(288519674)
  set requestcount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(288519674)
  $core.bool hasRequestcount() => $_has(1);
  @$pb.TagNumber(288519674)
  void clearRequestcount() => $_clearField(288519674);

  @$pb.TagNumber(341153238)
  $core.double get percentage => $_getN(2);
  @$pb.TagNumber(341153238)
  set percentage($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(341153238)
  $core.bool hasPercentage() => $_has(2);
  @$pb.TagNumber(341153238)
  void clearPercentage() => $_clearField(341153238);
}

class ByteMatchStatement extends $pb.GeneratedMessage {
  factory ByteMatchStatement({
    PositionalConstraint? positionalconstraint,
    $core.Iterable<TextTransformation>? texttransformations,
    $core.List<$core.int>? searchstring,
    FieldToMatch? fieldtomatch,
  }) {
    final result = create();
    if (positionalconstraint != null)
      result.positionalconstraint = positionalconstraint;
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (searchstring != null) result.searchstring = searchstring;
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    return result;
  }

  ByteMatchStatement._();

  factory ByteMatchStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ByteMatchStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ByteMatchStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<PositionalConstraint>(
        260734525, _omitFieldNames ? '' : 'positionalconstraint',
        enumValues: PositionalConstraint.values)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..a<$core.List<$core.int>>(
        318687365, _omitFieldNames ? '' : 'searchstring', $pb.PbFieldType.OY)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchStatement copyWith(void Function(ByteMatchStatement) updates) =>
      super.copyWith((message) => updates(message as ByteMatchStatement))
          as ByteMatchStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ByteMatchStatement create() => ByteMatchStatement._();
  @$core.override
  ByteMatchStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ByteMatchStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ByteMatchStatement>(create);
  static ByteMatchStatement? _defaultInstance;

  @$pb.TagNumber(260734525)
  PositionalConstraint get positionalconstraint => $_getN(0);
  @$pb.TagNumber(260734525)
  set positionalconstraint(PositionalConstraint value) =>
      $_setField(260734525, value);
  @$pb.TagNumber(260734525)
  $core.bool hasPositionalconstraint() => $_has(0);
  @$pb.TagNumber(260734525)
  void clearPositionalconstraint() => $_clearField(260734525);

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(1);

  @$pb.TagNumber(318687365)
  $core.List<$core.int> get searchstring => $_getN(2);
  @$pb.TagNumber(318687365)
  set searchstring($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(318687365)
  $core.bool hasSearchstring() => $_has(2);
  @$pb.TagNumber(318687365)
  void clearSearchstring() => $_clearField(318687365);

  @$pb.TagNumber(338348372)
  FieldToMatch get fieldtomatch => $_getN(3);
  @$pb.TagNumber(338348372)
  set fieldtomatch(FieldToMatch value) => $_setField(338348372, value);
  @$pb.TagNumber(338348372)
  $core.bool hasFieldtomatch() => $_has(3);
  @$pb.TagNumber(338348372)
  void clearFieldtomatch() => $_clearField(338348372);
  @$pb.TagNumber(338348372)
  FieldToMatch ensureFieldtomatch() => $_ensure(3);
}

class CaptchaAction extends $pb.GeneratedMessage {
  factory CaptchaAction({
    CustomRequestHandling? customrequesthandling,
  }) {
    final result = create();
    if (customrequesthandling != null)
      result.customrequesthandling = customrequesthandling;
    return result;
  }

  CaptchaAction._();

  factory CaptchaAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CaptchaAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CaptchaAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<CustomRequestHandling>(
        14845503, _omitFieldNames ? '' : 'customrequesthandling',
        subBuilder: CustomRequestHandling.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CaptchaAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CaptchaAction copyWith(void Function(CaptchaAction) updates) =>
      super.copyWith((message) => updates(message as CaptchaAction))
          as CaptchaAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CaptchaAction create() => CaptchaAction._();
  @$core.override
  CaptchaAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CaptchaAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CaptchaAction>(create);
  static CaptchaAction? _defaultInstance;

  @$pb.TagNumber(14845503)
  CustomRequestHandling get customrequesthandling => $_getN(0);
  @$pb.TagNumber(14845503)
  set customrequesthandling(CustomRequestHandling value) =>
      $_setField(14845503, value);
  @$pb.TagNumber(14845503)
  $core.bool hasCustomrequesthandling() => $_has(0);
  @$pb.TagNumber(14845503)
  void clearCustomrequesthandling() => $_clearField(14845503);
  @$pb.TagNumber(14845503)
  CustomRequestHandling ensureCustomrequesthandling() => $_ensure(0);
}

class CaptchaConfig extends $pb.GeneratedMessage {
  factory CaptchaConfig({
    ImmunityTimeProperty? immunitytimeproperty,
  }) {
    final result = create();
    if (immunitytimeproperty != null)
      result.immunitytimeproperty = immunitytimeproperty;
    return result;
  }

  CaptchaConfig._();

  factory CaptchaConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CaptchaConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CaptchaConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<ImmunityTimeProperty>(
        28330836, _omitFieldNames ? '' : 'immunitytimeproperty',
        subBuilder: ImmunityTimeProperty.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CaptchaConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CaptchaConfig copyWith(void Function(CaptchaConfig) updates) =>
      super.copyWith((message) => updates(message as CaptchaConfig))
          as CaptchaConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CaptchaConfig create() => CaptchaConfig._();
  @$core.override
  CaptchaConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CaptchaConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CaptchaConfig>(create);
  static CaptchaConfig? _defaultInstance;

  @$pb.TagNumber(28330836)
  ImmunityTimeProperty get immunitytimeproperty => $_getN(0);
  @$pb.TagNumber(28330836)
  set immunitytimeproperty(ImmunityTimeProperty value) =>
      $_setField(28330836, value);
  @$pb.TagNumber(28330836)
  $core.bool hasImmunitytimeproperty() => $_has(0);
  @$pb.TagNumber(28330836)
  void clearImmunitytimeproperty() => $_clearField(28330836);
  @$pb.TagNumber(28330836)
  ImmunityTimeProperty ensureImmunitytimeproperty() => $_ensure(0);
}

class CaptchaResponse extends $pb.GeneratedMessage {
  factory CaptchaResponse({
    FailureReason? failurereason,
    $fixnum.Int64? solvetimestamp,
    $core.int? responsecode,
  }) {
    final result = create();
    if (failurereason != null) result.failurereason = failurereason;
    if (solvetimestamp != null) result.solvetimestamp = solvetimestamp;
    if (responsecode != null) result.responsecode = responsecode;
    return result;
  }

  CaptchaResponse._();

  factory CaptchaResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CaptchaResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CaptchaResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<FailureReason>(232322142, _omitFieldNames ? '' : 'failurereason',
        enumValues: FailureReason.values)
    ..aInt64(433307777, _omitFieldNames ? '' : 'solvetimestamp')
    ..aI(447553700, _omitFieldNames ? '' : 'responsecode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CaptchaResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CaptchaResponse copyWith(void Function(CaptchaResponse) updates) =>
      super.copyWith((message) => updates(message as CaptchaResponse))
          as CaptchaResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CaptchaResponse create() => CaptchaResponse._();
  @$core.override
  CaptchaResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CaptchaResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CaptchaResponse>(create);
  static CaptchaResponse? _defaultInstance;

  @$pb.TagNumber(232322142)
  FailureReason get failurereason => $_getN(0);
  @$pb.TagNumber(232322142)
  set failurereason(FailureReason value) => $_setField(232322142, value);
  @$pb.TagNumber(232322142)
  $core.bool hasFailurereason() => $_has(0);
  @$pb.TagNumber(232322142)
  void clearFailurereason() => $_clearField(232322142);

  @$pb.TagNumber(433307777)
  $fixnum.Int64 get solvetimestamp => $_getI64(1);
  @$pb.TagNumber(433307777)
  set solvetimestamp($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(433307777)
  $core.bool hasSolvetimestamp() => $_has(1);
  @$pb.TagNumber(433307777)
  void clearSolvetimestamp() => $_clearField(433307777);

  @$pb.TagNumber(447553700)
  $core.int get responsecode => $_getIZ(2);
  @$pb.TagNumber(447553700)
  set responsecode($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(447553700)
  $core.bool hasResponsecode() => $_has(2);
  @$pb.TagNumber(447553700)
  void clearResponsecode() => $_clearField(447553700);
}

class ChallengeAction extends $pb.GeneratedMessage {
  factory ChallengeAction({
    CustomRequestHandling? customrequesthandling,
  }) {
    final result = create();
    if (customrequesthandling != null)
      result.customrequesthandling = customrequesthandling;
    return result;
  }

  ChallengeAction._();

  factory ChallengeAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChallengeAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChallengeAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<CustomRequestHandling>(
        14845503, _omitFieldNames ? '' : 'customrequesthandling',
        subBuilder: CustomRequestHandling.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChallengeAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChallengeAction copyWith(void Function(ChallengeAction) updates) =>
      super.copyWith((message) => updates(message as ChallengeAction))
          as ChallengeAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChallengeAction create() => ChallengeAction._();
  @$core.override
  ChallengeAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChallengeAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChallengeAction>(create);
  static ChallengeAction? _defaultInstance;

  @$pb.TagNumber(14845503)
  CustomRequestHandling get customrequesthandling => $_getN(0);
  @$pb.TagNumber(14845503)
  set customrequesthandling(CustomRequestHandling value) =>
      $_setField(14845503, value);
  @$pb.TagNumber(14845503)
  $core.bool hasCustomrequesthandling() => $_has(0);
  @$pb.TagNumber(14845503)
  void clearCustomrequesthandling() => $_clearField(14845503);
  @$pb.TagNumber(14845503)
  CustomRequestHandling ensureCustomrequesthandling() => $_ensure(0);
}

class ChallengeConfig extends $pb.GeneratedMessage {
  factory ChallengeConfig({
    ImmunityTimeProperty? immunitytimeproperty,
  }) {
    final result = create();
    if (immunitytimeproperty != null)
      result.immunitytimeproperty = immunitytimeproperty;
    return result;
  }

  ChallengeConfig._();

  factory ChallengeConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChallengeConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChallengeConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<ImmunityTimeProperty>(
        28330836, _omitFieldNames ? '' : 'immunitytimeproperty',
        subBuilder: ImmunityTimeProperty.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChallengeConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChallengeConfig copyWith(void Function(ChallengeConfig) updates) =>
      super.copyWith((message) => updates(message as ChallengeConfig))
          as ChallengeConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChallengeConfig create() => ChallengeConfig._();
  @$core.override
  ChallengeConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChallengeConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChallengeConfig>(create);
  static ChallengeConfig? _defaultInstance;

  @$pb.TagNumber(28330836)
  ImmunityTimeProperty get immunitytimeproperty => $_getN(0);
  @$pb.TagNumber(28330836)
  set immunitytimeproperty(ImmunityTimeProperty value) =>
      $_setField(28330836, value);
  @$pb.TagNumber(28330836)
  $core.bool hasImmunitytimeproperty() => $_has(0);
  @$pb.TagNumber(28330836)
  void clearImmunitytimeproperty() => $_clearField(28330836);
  @$pb.TagNumber(28330836)
  ImmunityTimeProperty ensureImmunitytimeproperty() => $_ensure(0);
}

class ChallengeResponse extends $pb.GeneratedMessage {
  factory ChallengeResponse({
    FailureReason? failurereason,
    $fixnum.Int64? solvetimestamp,
    $core.int? responsecode,
  }) {
    final result = create();
    if (failurereason != null) result.failurereason = failurereason;
    if (solvetimestamp != null) result.solvetimestamp = solvetimestamp;
    if (responsecode != null) result.responsecode = responsecode;
    return result;
  }

  ChallengeResponse._();

  factory ChallengeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChallengeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChallengeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<FailureReason>(232322142, _omitFieldNames ? '' : 'failurereason',
        enumValues: FailureReason.values)
    ..aInt64(433307777, _omitFieldNames ? '' : 'solvetimestamp')
    ..aI(447553700, _omitFieldNames ? '' : 'responsecode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChallengeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChallengeResponse copyWith(void Function(ChallengeResponse) updates) =>
      super.copyWith((message) => updates(message as ChallengeResponse))
          as ChallengeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChallengeResponse create() => ChallengeResponse._();
  @$core.override
  ChallengeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChallengeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChallengeResponse>(create);
  static ChallengeResponse? _defaultInstance;

  @$pb.TagNumber(232322142)
  FailureReason get failurereason => $_getN(0);
  @$pb.TagNumber(232322142)
  set failurereason(FailureReason value) => $_setField(232322142, value);
  @$pb.TagNumber(232322142)
  $core.bool hasFailurereason() => $_has(0);
  @$pb.TagNumber(232322142)
  void clearFailurereason() => $_clearField(232322142);

  @$pb.TagNumber(433307777)
  $fixnum.Int64 get solvetimestamp => $_getI64(1);
  @$pb.TagNumber(433307777)
  set solvetimestamp($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(433307777)
  $core.bool hasSolvetimestamp() => $_has(1);
  @$pb.TagNumber(433307777)
  void clearSolvetimestamp() => $_clearField(433307777);

  @$pb.TagNumber(447553700)
  $core.int get responsecode => $_getIZ(2);
  @$pb.TagNumber(447553700)
  set responsecode($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(447553700)
  $core.bool hasResponsecode() => $_has(2);
  @$pb.TagNumber(447553700)
  void clearResponsecode() => $_clearField(447553700);
}

class CheckCapacityRequest extends $pb.GeneratedMessage {
  factory CheckCapacityRequest({
    $core.Iterable<Rule>? rules,
    Scope? scope,
  }) {
    final result = create();
    if (rules != null) result.rules.addAll(rules);
    if (scope != null) result.scope = scope;
    return result;
  }

  CheckCapacityRequest._();

  factory CheckCapacityRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CheckCapacityRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CheckCapacityRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<Rule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: Rule.create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckCapacityRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckCapacityRequest copyWith(void Function(CheckCapacityRequest) updates) =>
      super.copyWith((message) => updates(message as CheckCapacityRequest))
          as CheckCapacityRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CheckCapacityRequest create() => CheckCapacityRequest._();
  @$core.override
  CheckCapacityRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CheckCapacityRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CheckCapacityRequest>(create);
  static CheckCapacityRequest? _defaultInstance;

  @$pb.TagNumber(42675585)
  $pb.PbList<Rule> get rules => $_getList(0);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);
}

class CheckCapacityResponse extends $pb.GeneratedMessage {
  factory CheckCapacityResponse({
    $fixnum.Int64? capacity,
  }) {
    final result = create();
    if (capacity != null) result.capacity = capacity;
    return result;
  }

  CheckCapacityResponse._();

  factory CheckCapacityResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CheckCapacityResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CheckCapacityResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aInt64(107253930, _omitFieldNames ? '' : 'capacity')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckCapacityResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CheckCapacityResponse copyWith(
          void Function(CheckCapacityResponse) updates) =>
      super.copyWith((message) => updates(message as CheckCapacityResponse))
          as CheckCapacityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CheckCapacityResponse create() => CheckCapacityResponse._();
  @$core.override
  CheckCapacityResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CheckCapacityResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CheckCapacityResponse>(create);
  static CheckCapacityResponse? _defaultInstance;

  @$pb.TagNumber(107253930)
  $fixnum.Int64 get capacity => $_getI64(0);
  @$pb.TagNumber(107253930)
  set capacity($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(107253930)
  $core.bool hasCapacity() => $_has(0);
  @$pb.TagNumber(107253930)
  void clearCapacity() => $_clearField(107253930);
}

class ClientSideAction extends $pb.GeneratedMessage {
  factory ClientSideAction({
    SensitivityToAct? sensitivity,
    $core.Iterable<Regex>? exempturiregularexpressions,
    UsageOfAction? usageofaction,
  }) {
    final result = create();
    if (sensitivity != null) result.sensitivity = sensitivity;
    if (exempturiregularexpressions != null)
      result.exempturiregularexpressions.addAll(exempturiregularexpressions);
    if (usageofaction != null) result.usageofaction = usageofaction;
    return result;
  }

  ClientSideAction._();

  factory ClientSideAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ClientSideAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ClientSideAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<SensitivityToAct>(30112173, _omitFieldNames ? '' : 'sensitivity',
        enumValues: SensitivityToAct.values)
    ..pPM<Regex>(
        148735428, _omitFieldNames ? '' : 'exempturiregularexpressions',
        subBuilder: Regex.create)
    ..aE<UsageOfAction>(213215420, _omitFieldNames ? '' : 'usageofaction',
        enumValues: UsageOfAction.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientSideAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientSideAction copyWith(void Function(ClientSideAction) updates) =>
      super.copyWith((message) => updates(message as ClientSideAction))
          as ClientSideAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ClientSideAction create() => ClientSideAction._();
  @$core.override
  ClientSideAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ClientSideAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ClientSideAction>(create);
  static ClientSideAction? _defaultInstance;

  @$pb.TagNumber(30112173)
  SensitivityToAct get sensitivity => $_getN(0);
  @$pb.TagNumber(30112173)
  set sensitivity(SensitivityToAct value) => $_setField(30112173, value);
  @$pb.TagNumber(30112173)
  $core.bool hasSensitivity() => $_has(0);
  @$pb.TagNumber(30112173)
  void clearSensitivity() => $_clearField(30112173);

  @$pb.TagNumber(148735428)
  $pb.PbList<Regex> get exempturiregularexpressions => $_getList(1);

  @$pb.TagNumber(213215420)
  UsageOfAction get usageofaction => $_getN(2);
  @$pb.TagNumber(213215420)
  set usageofaction(UsageOfAction value) => $_setField(213215420, value);
  @$pb.TagNumber(213215420)
  $core.bool hasUsageofaction() => $_has(2);
  @$pb.TagNumber(213215420)
  void clearUsageofaction() => $_clearField(213215420);
}

class ClientSideActionConfig extends $pb.GeneratedMessage {
  factory ClientSideActionConfig({
    ClientSideAction? challenge,
  }) {
    final result = create();
    if (challenge != null) result.challenge = challenge;
    return result;
  }

  ClientSideActionConfig._();

  factory ClientSideActionConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ClientSideActionConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ClientSideActionConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<ClientSideAction>(128811743, _omitFieldNames ? '' : 'challenge',
        subBuilder: ClientSideAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientSideActionConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientSideActionConfig copyWith(
          void Function(ClientSideActionConfig) updates) =>
      super.copyWith((message) => updates(message as ClientSideActionConfig))
          as ClientSideActionConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ClientSideActionConfig create() => ClientSideActionConfig._();
  @$core.override
  ClientSideActionConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ClientSideActionConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ClientSideActionConfig>(create);
  static ClientSideActionConfig? _defaultInstance;

  @$pb.TagNumber(128811743)
  ClientSideAction get challenge => $_getN(0);
  @$pb.TagNumber(128811743)
  set challenge(ClientSideAction value) => $_setField(128811743, value);
  @$pb.TagNumber(128811743)
  $core.bool hasChallenge() => $_has(0);
  @$pb.TagNumber(128811743)
  void clearChallenge() => $_clearField(128811743);
  @$pb.TagNumber(128811743)
  ClientSideAction ensureChallenge() => $_ensure(0);
}

class Condition extends $pb.GeneratedMessage {
  factory Condition({
    ActionCondition? actioncondition,
    LabelNameCondition? labelnamecondition,
  }) {
    final result = create();
    if (actioncondition != null) result.actioncondition = actioncondition;
    if (labelnamecondition != null)
      result.labelnamecondition = labelnamecondition;
    return result;
  }

  Condition._();

  factory Condition.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Condition.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Condition',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<ActionCondition>(46185833, _omitFieldNames ? '' : 'actioncondition',
        subBuilder: ActionCondition.create)
    ..aOM<LabelNameCondition>(
        360367980, _omitFieldNames ? '' : 'labelnamecondition',
        subBuilder: LabelNameCondition.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Condition clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Condition copyWith(void Function(Condition) updates) =>
      super.copyWith((message) => updates(message as Condition)) as Condition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Condition create() => Condition._();
  @$core.override
  Condition createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Condition getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Condition>(create);
  static Condition? _defaultInstance;

  @$pb.TagNumber(46185833)
  ActionCondition get actioncondition => $_getN(0);
  @$pb.TagNumber(46185833)
  set actioncondition(ActionCondition value) => $_setField(46185833, value);
  @$pb.TagNumber(46185833)
  $core.bool hasActioncondition() => $_has(0);
  @$pb.TagNumber(46185833)
  void clearActioncondition() => $_clearField(46185833);
  @$pb.TagNumber(46185833)
  ActionCondition ensureActioncondition() => $_ensure(0);

  @$pb.TagNumber(360367980)
  LabelNameCondition get labelnamecondition => $_getN(1);
  @$pb.TagNumber(360367980)
  set labelnamecondition(LabelNameCondition value) =>
      $_setField(360367980, value);
  @$pb.TagNumber(360367980)
  $core.bool hasLabelnamecondition() => $_has(1);
  @$pb.TagNumber(360367980)
  void clearLabelnamecondition() => $_clearField(360367980);
  @$pb.TagNumber(360367980)
  LabelNameCondition ensureLabelnamecondition() => $_ensure(1);
}

class CookieMatchPattern extends $pb.GeneratedMessage {
  factory CookieMatchPattern({
    $core.Iterable<$core.String>? includedcookies,
    $core.Iterable<$core.String>? excludedcookies,
    All? all,
  }) {
    final result = create();
    if (includedcookies != null) result.includedcookies.addAll(includedcookies);
    if (excludedcookies != null) result.excludedcookies.addAll(excludedcookies);
    if (all != null) result.all = all;
    return result;
  }

  CookieMatchPattern._();

  factory CookieMatchPattern.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CookieMatchPattern.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CookieMatchPattern',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(297750153, _omitFieldNames ? '' : 'includedcookies')
    ..pPS(362170031, _omitFieldNames ? '' : 'excludedcookies')
    ..aOM<All>(363848549, _omitFieldNames ? '' : 'all', subBuilder: All.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CookieMatchPattern clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CookieMatchPattern copyWith(void Function(CookieMatchPattern) updates) =>
      super.copyWith((message) => updates(message as CookieMatchPattern))
          as CookieMatchPattern;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CookieMatchPattern create() => CookieMatchPattern._();
  @$core.override
  CookieMatchPattern createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CookieMatchPattern getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CookieMatchPattern>(create);
  static CookieMatchPattern? _defaultInstance;

  @$pb.TagNumber(297750153)
  $pb.PbList<$core.String> get includedcookies => $_getList(0);

  @$pb.TagNumber(362170031)
  $pb.PbList<$core.String> get excludedcookies => $_getList(1);

  @$pb.TagNumber(363848549)
  All get all => $_getN(2);
  @$pb.TagNumber(363848549)
  set all(All value) => $_setField(363848549, value);
  @$pb.TagNumber(363848549)
  $core.bool hasAll() => $_has(2);
  @$pb.TagNumber(363848549)
  void clearAll() => $_clearField(363848549);
  @$pb.TagNumber(363848549)
  All ensureAll() => $_ensure(2);
}

class Cookies extends $pb.GeneratedMessage {
  factory Cookies({
    OversizeHandling? oversizehandling,
    MapMatchScope? matchscope,
    CookieMatchPattern? matchpattern,
  }) {
    final result = create();
    if (oversizehandling != null) result.oversizehandling = oversizehandling;
    if (matchscope != null) result.matchscope = matchscope;
    if (matchpattern != null) result.matchpattern = matchpattern;
    return result;
  }

  Cookies._();

  factory Cookies.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Cookies.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Cookies',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<OversizeHandling>(139375132, _omitFieldNames ? '' : 'oversizehandling',
        enumValues: OversizeHandling.values)
    ..aE<MapMatchScope>(272445459, _omitFieldNames ? '' : 'matchscope',
        enumValues: MapMatchScope.values)
    ..aOM<CookieMatchPattern>(294565637, _omitFieldNames ? '' : 'matchpattern',
        subBuilder: CookieMatchPattern.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Cookies clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Cookies copyWith(void Function(Cookies) updates) =>
      super.copyWith((message) => updates(message as Cookies)) as Cookies;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Cookies create() => Cookies._();
  @$core.override
  Cookies createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Cookies getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Cookies>(create);
  static Cookies? _defaultInstance;

  @$pb.TagNumber(139375132)
  OversizeHandling get oversizehandling => $_getN(0);
  @$pb.TagNumber(139375132)
  set oversizehandling(OversizeHandling value) => $_setField(139375132, value);
  @$pb.TagNumber(139375132)
  $core.bool hasOversizehandling() => $_has(0);
  @$pb.TagNumber(139375132)
  void clearOversizehandling() => $_clearField(139375132);

  @$pb.TagNumber(272445459)
  MapMatchScope get matchscope => $_getN(1);
  @$pb.TagNumber(272445459)
  set matchscope(MapMatchScope value) => $_setField(272445459, value);
  @$pb.TagNumber(272445459)
  $core.bool hasMatchscope() => $_has(1);
  @$pb.TagNumber(272445459)
  void clearMatchscope() => $_clearField(272445459);

  @$pb.TagNumber(294565637)
  CookieMatchPattern get matchpattern => $_getN(2);
  @$pb.TagNumber(294565637)
  set matchpattern(CookieMatchPattern value) => $_setField(294565637, value);
  @$pb.TagNumber(294565637)
  $core.bool hasMatchpattern() => $_has(2);
  @$pb.TagNumber(294565637)
  void clearMatchpattern() => $_clearField(294565637);
  @$pb.TagNumber(294565637)
  CookieMatchPattern ensureMatchpattern() => $_ensure(2);
}

class CountAction extends $pb.GeneratedMessage {
  factory CountAction({
    CustomRequestHandling? customrequesthandling,
  }) {
    final result = create();
    if (customrequesthandling != null)
      result.customrequesthandling = customrequesthandling;
    return result;
  }

  CountAction._();

  factory CountAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CountAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CountAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<CustomRequestHandling>(
        14845503, _omitFieldNames ? '' : 'customrequesthandling',
        subBuilder: CustomRequestHandling.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CountAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CountAction copyWith(void Function(CountAction) updates) =>
      super.copyWith((message) => updates(message as CountAction))
          as CountAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CountAction create() => CountAction._();
  @$core.override
  CountAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CountAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CountAction>(create);
  static CountAction? _defaultInstance;

  @$pb.TagNumber(14845503)
  CustomRequestHandling get customrequesthandling => $_getN(0);
  @$pb.TagNumber(14845503)
  set customrequesthandling(CustomRequestHandling value) =>
      $_setField(14845503, value);
  @$pb.TagNumber(14845503)
  $core.bool hasCustomrequesthandling() => $_has(0);
  @$pb.TagNumber(14845503)
  void clearCustomrequesthandling() => $_clearField(14845503);
  @$pb.TagNumber(14845503)
  CustomRequestHandling ensureCustomrequesthandling() => $_ensure(0);
}

class CreateAPIKeyRequest extends $pb.GeneratedMessage {
  factory CreateAPIKeyRequest({
    $core.Iterable<$core.String>? tokendomains,
    Scope? scope,
  }) {
    final result = create();
    if (tokendomains != null) result.tokendomains.addAll(tokendomains);
    if (scope != null) result.scope = scope;
    return result;
  }

  CreateAPIKeyRequest._();

  factory CreateAPIKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateAPIKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateAPIKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(11638560, _omitFieldNames ? '' : 'tokendomains')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAPIKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAPIKeyRequest copyWith(void Function(CreateAPIKeyRequest) updates) =>
      super.copyWith((message) => updates(message as CreateAPIKeyRequest))
          as CreateAPIKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateAPIKeyRequest create() => CreateAPIKeyRequest._();
  @$core.override
  CreateAPIKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateAPIKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateAPIKeyRequest>(create);
  static CreateAPIKeyRequest? _defaultInstance;

  @$pb.TagNumber(11638560)
  $pb.PbList<$core.String> get tokendomains => $_getList(0);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);
}

class CreateAPIKeyResponse extends $pb.GeneratedMessage {
  factory CreateAPIKeyResponse({
    $core.String? apikey,
  }) {
    final result = create();
    if (apikey != null) result.apikey = apikey;
    return result;
  }

  CreateAPIKeyResponse._();

  factory CreateAPIKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateAPIKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateAPIKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(274818239, _omitFieldNames ? '' : 'apikey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAPIKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateAPIKeyResponse copyWith(void Function(CreateAPIKeyResponse) updates) =>
      super.copyWith((message) => updates(message as CreateAPIKeyResponse))
          as CreateAPIKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateAPIKeyResponse create() => CreateAPIKeyResponse._();
  @$core.override
  CreateAPIKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateAPIKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateAPIKeyResponse>(create);
  static CreateAPIKeyResponse? _defaultInstance;

  @$pb.TagNumber(274818239)
  $core.String get apikey => $_getSZ(0);
  @$pb.TagNumber(274818239)
  set apikey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(274818239)
  $core.bool hasApikey() => $_has(0);
  @$pb.TagNumber(274818239)
  void clearApikey() => $_clearField(274818239);
}

class CreateIPSetRequest extends $pb.GeneratedMessage {
  factory CreateIPSetRequest({
    Scope? scope,
    $core.String? description,
    $core.String? name,
    IPAddressVersion? ipaddressversion,
    $core.Iterable<$core.String>? addresses,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (ipaddressversion != null) result.ipaddressversion = ipaddressversion;
    if (addresses != null) result.addresses.addAll(addresses);
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  CreateIPSetRequest._();

  factory CreateIPSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateIPSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateIPSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<IPAddressVersion>(313363841, _omitFieldNames ? '' : 'ipaddressversion',
        enumValues: IPAddressVersion.values)
    ..pPS(375939972, _omitFieldNames ? '' : 'addresses')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateIPSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateIPSetRequest copyWith(void Function(CreateIPSetRequest) updates) =>
      super.copyWith((message) => updates(message as CreateIPSetRequest))
          as CreateIPSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateIPSetRequest create() => CreateIPSetRequest._();
  @$core.override
  CreateIPSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateIPSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateIPSetRequest>(create);
  static CreateIPSetRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(313363841)
  IPAddressVersion get ipaddressversion => $_getN(3);
  @$pb.TagNumber(313363841)
  set ipaddressversion(IPAddressVersion value) => $_setField(313363841, value);
  @$pb.TagNumber(313363841)
  $core.bool hasIpaddressversion() => $_has(3);
  @$pb.TagNumber(313363841)
  void clearIpaddressversion() => $_clearField(313363841);

  @$pb.TagNumber(375939972)
  $pb.PbList<$core.String> get addresses => $_getList(4);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(5);
}

class CreateIPSetResponse extends $pb.GeneratedMessage {
  factory CreateIPSetResponse({
    IPSetSummary? summary,
  }) {
    final result = create();
    if (summary != null) result.summary = summary;
    return result;
  }

  CreateIPSetResponse._();

  factory CreateIPSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateIPSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateIPSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<IPSetSummary>(21935540, _omitFieldNames ? '' : 'summary',
        subBuilder: IPSetSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateIPSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateIPSetResponse copyWith(void Function(CreateIPSetResponse) updates) =>
      super.copyWith((message) => updates(message as CreateIPSetResponse))
          as CreateIPSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateIPSetResponse create() => CreateIPSetResponse._();
  @$core.override
  CreateIPSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateIPSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateIPSetResponse>(create);
  static CreateIPSetResponse? _defaultInstance;

  @$pb.TagNumber(21935540)
  IPSetSummary get summary => $_getN(0);
  @$pb.TagNumber(21935540)
  set summary(IPSetSummary value) => $_setField(21935540, value);
  @$pb.TagNumber(21935540)
  $core.bool hasSummary() => $_has(0);
  @$pb.TagNumber(21935540)
  void clearSummary() => $_clearField(21935540);
  @$pb.TagNumber(21935540)
  IPSetSummary ensureSummary() => $_ensure(0);
}

class CreateRegexPatternSetRequest extends $pb.GeneratedMessage {
  factory CreateRegexPatternSetRequest({
    Scope? scope,
    $core.String? description,
    $core.Iterable<Regex>? regularexpressionlist,
    $core.String? name,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (description != null) result.description = description;
    if (regularexpressionlist != null)
      result.regularexpressionlist.addAll(regularexpressionlist);
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  CreateRegexPatternSetRequest._();

  factory CreateRegexPatternSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRegexPatternSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRegexPatternSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..pPM<Regex>(123612838, _omitFieldNames ? '' : 'regularexpressionlist',
        subBuilder: Regex.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRegexPatternSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRegexPatternSetRequest copyWith(
          void Function(CreateRegexPatternSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateRegexPatternSetRequest))
          as CreateRegexPatternSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRegexPatternSetRequest create() =>
      CreateRegexPatternSetRequest._();
  @$core.override
  CreateRegexPatternSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRegexPatternSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRegexPatternSetRequest>(create);
  static CreateRegexPatternSetRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(123612838)
  $pb.PbList<Regex> get regularexpressionlist => $_getList(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(4);
}

class CreateRegexPatternSetResponse extends $pb.GeneratedMessage {
  factory CreateRegexPatternSetResponse({
    RegexPatternSetSummary? summary,
  }) {
    final result = create();
    if (summary != null) result.summary = summary;
    return result;
  }

  CreateRegexPatternSetResponse._();

  factory CreateRegexPatternSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRegexPatternSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRegexPatternSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<RegexPatternSetSummary>(21935540, _omitFieldNames ? '' : 'summary',
        subBuilder: RegexPatternSetSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRegexPatternSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRegexPatternSetResponse copyWith(
          void Function(CreateRegexPatternSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateRegexPatternSetResponse))
          as CreateRegexPatternSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRegexPatternSetResponse create() =>
      CreateRegexPatternSetResponse._();
  @$core.override
  CreateRegexPatternSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRegexPatternSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRegexPatternSetResponse>(create);
  static CreateRegexPatternSetResponse? _defaultInstance;

  @$pb.TagNumber(21935540)
  RegexPatternSetSummary get summary => $_getN(0);
  @$pb.TagNumber(21935540)
  set summary(RegexPatternSetSummary value) => $_setField(21935540, value);
  @$pb.TagNumber(21935540)
  $core.bool hasSummary() => $_has(0);
  @$pb.TagNumber(21935540)
  void clearSummary() => $_clearField(21935540);
  @$pb.TagNumber(21935540)
  RegexPatternSetSummary ensureSummary() => $_ensure(0);
}

class CreateRuleGroupRequest extends $pb.GeneratedMessage {
  factory CreateRuleGroupRequest({
    $core.Iterable<Rule>? rules,
    $core.Iterable<$core.MapEntry<$core.String, CustomResponseBody>>?
        customresponsebodies,
    Scope? scope,
    $fixnum.Int64? capacity,
    $core.String? description,
    $core.String? name,
    $core.Iterable<Tag>? tags,
    VisibilityConfig? visibilityconfig,
  }) {
    final result = create();
    if (rules != null) result.rules.addAll(rules);
    if (customresponsebodies != null)
      result.customresponsebodies.addEntries(customresponsebodies);
    if (scope != null) result.scope = scope;
    if (capacity != null) result.capacity = capacity;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    if (visibilityconfig != null) result.visibilityconfig = visibilityconfig;
    return result;
  }

  CreateRuleGroupRequest._();

  factory CreateRuleGroupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRuleGroupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRuleGroupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<Rule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: Rule.create)
    ..m<$core.String, CustomResponseBody>(
        42731774, _omitFieldNames ? '' : 'customresponsebodies',
        entryClassName: 'CreateRuleGroupRequest.CustomresponsebodiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: CustomResponseBody.create,
        valueDefaultOrMaker: CustomResponseBody.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aInt64(107253930, _omitFieldNames ? '' : 'capacity')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOM<VisibilityConfig>(
        507762728, _omitFieldNames ? '' : 'visibilityconfig',
        subBuilder: VisibilityConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRuleGroupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRuleGroupRequest copyWith(
          void Function(CreateRuleGroupRequest) updates) =>
      super.copyWith((message) => updates(message as CreateRuleGroupRequest))
          as CreateRuleGroupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRuleGroupRequest create() => CreateRuleGroupRequest._();
  @$core.override
  CreateRuleGroupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRuleGroupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRuleGroupRequest>(create);
  static CreateRuleGroupRequest? _defaultInstance;

  @$pb.TagNumber(42675585)
  $pb.PbList<Rule> get rules => $_getList(0);

  @$pb.TagNumber(42731774)
  $pb.PbMap<$core.String, CustomResponseBody> get customresponsebodies =>
      $_getMap(1);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(2);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(2);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(107253930)
  $fixnum.Int64 get capacity => $_getI64(3);
  @$pb.TagNumber(107253930)
  set capacity($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(107253930)
  $core.bool hasCapacity() => $_has(3);
  @$pb.TagNumber(107253930)
  void clearCapacity() => $_clearField(107253930);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(6);

  @$pb.TagNumber(507762728)
  VisibilityConfig get visibilityconfig => $_getN(7);
  @$pb.TagNumber(507762728)
  set visibilityconfig(VisibilityConfig value) => $_setField(507762728, value);
  @$pb.TagNumber(507762728)
  $core.bool hasVisibilityconfig() => $_has(7);
  @$pb.TagNumber(507762728)
  void clearVisibilityconfig() => $_clearField(507762728);
  @$pb.TagNumber(507762728)
  VisibilityConfig ensureVisibilityconfig() => $_ensure(7);
}

class CreateRuleGroupResponse extends $pb.GeneratedMessage {
  factory CreateRuleGroupResponse({
    RuleGroupSummary? summary,
  }) {
    final result = create();
    if (summary != null) result.summary = summary;
    return result;
  }

  CreateRuleGroupResponse._();

  factory CreateRuleGroupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRuleGroupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRuleGroupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<RuleGroupSummary>(21935540, _omitFieldNames ? '' : 'summary',
        subBuilder: RuleGroupSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRuleGroupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRuleGroupResponse copyWith(
          void Function(CreateRuleGroupResponse) updates) =>
      super.copyWith((message) => updates(message as CreateRuleGroupResponse))
          as CreateRuleGroupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRuleGroupResponse create() => CreateRuleGroupResponse._();
  @$core.override
  CreateRuleGroupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRuleGroupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRuleGroupResponse>(create);
  static CreateRuleGroupResponse? _defaultInstance;

  @$pb.TagNumber(21935540)
  RuleGroupSummary get summary => $_getN(0);
  @$pb.TagNumber(21935540)
  set summary(RuleGroupSummary value) => $_setField(21935540, value);
  @$pb.TagNumber(21935540)
  $core.bool hasSummary() => $_has(0);
  @$pb.TagNumber(21935540)
  void clearSummary() => $_clearField(21935540);
  @$pb.TagNumber(21935540)
  RuleGroupSummary ensureSummary() => $_ensure(0);
}

class CreateWebACLRequest extends $pb.GeneratedMessage {
  factory CreateWebACLRequest({
    $core.Iterable<$core.String>? tokendomains,
    $core.Iterable<Rule>? rules,
    $core.Iterable<$core.MapEntry<$core.String, CustomResponseBody>>?
        customresponsebodies,
    ChallengeConfig? challengeconfig,
    CaptchaConfig? captchaconfig,
    Scope? scope,
    OnSourceDDoSProtectionConfig? onsourceddosprotectionconfig,
    $core.String? description,
    $core.String? name,
    DefaultAction? defaultaction,
    $core.Iterable<Tag>? tags,
    AssociationConfig? associationconfig,
    DataProtectionConfig? dataprotectionconfig,
    ApplicationConfig? applicationconfig,
    VisibilityConfig? visibilityconfig,
  }) {
    final result = create();
    if (tokendomains != null) result.tokendomains.addAll(tokendomains);
    if (rules != null) result.rules.addAll(rules);
    if (customresponsebodies != null)
      result.customresponsebodies.addEntries(customresponsebodies);
    if (challengeconfig != null) result.challengeconfig = challengeconfig;
    if (captchaconfig != null) result.captchaconfig = captchaconfig;
    if (scope != null) result.scope = scope;
    if (onsourceddosprotectionconfig != null)
      result.onsourceddosprotectionconfig = onsourceddosprotectionconfig;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (defaultaction != null) result.defaultaction = defaultaction;
    if (tags != null) result.tags.addAll(tags);
    if (associationconfig != null) result.associationconfig = associationconfig;
    if (dataprotectionconfig != null)
      result.dataprotectionconfig = dataprotectionconfig;
    if (applicationconfig != null) result.applicationconfig = applicationconfig;
    if (visibilityconfig != null) result.visibilityconfig = visibilityconfig;
    return result;
  }

  CreateWebACLRequest._();

  factory CreateWebACLRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateWebACLRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateWebACLRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(11638560, _omitFieldNames ? '' : 'tokendomains')
    ..pPM<Rule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: Rule.create)
    ..m<$core.String, CustomResponseBody>(
        42731774, _omitFieldNames ? '' : 'customresponsebodies',
        entryClassName: 'CreateWebACLRequest.CustomresponsebodiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: CustomResponseBody.create,
        valueDefaultOrMaker: CustomResponseBody.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..aOM<ChallengeConfig>(48990889, _omitFieldNames ? '' : 'challengeconfig',
        subBuilder: ChallengeConfig.create)
    ..aOM<CaptchaConfig>(60547064, _omitFieldNames ? '' : 'captchaconfig',
        subBuilder: CaptchaConfig.create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOM<OnSourceDDoSProtectionConfig>(
        105229063, _omitFieldNames ? '' : 'onsourceddosprotectionconfig',
        subBuilder: OnSourceDDoSProtectionConfig.create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<DefaultAction>(322663861, _omitFieldNames ? '' : 'defaultaction',
        subBuilder: DefaultAction.create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOM<AssociationConfig>(
        419675691, _omitFieldNames ? '' : 'associationconfig',
        subBuilder: AssociationConfig.create)
    ..aOM<DataProtectionConfig>(
        464792245, _omitFieldNames ? '' : 'dataprotectionconfig',
        subBuilder: DataProtectionConfig.create)
    ..aOM<ApplicationConfig>(
        501131976, _omitFieldNames ? '' : 'applicationconfig',
        subBuilder: ApplicationConfig.create)
    ..aOM<VisibilityConfig>(
        507762728, _omitFieldNames ? '' : 'visibilityconfig',
        subBuilder: VisibilityConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWebACLRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWebACLRequest copyWith(void Function(CreateWebACLRequest) updates) =>
      super.copyWith((message) => updates(message as CreateWebACLRequest))
          as CreateWebACLRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWebACLRequest create() => CreateWebACLRequest._();
  @$core.override
  CreateWebACLRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateWebACLRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateWebACLRequest>(create);
  static CreateWebACLRequest? _defaultInstance;

  @$pb.TagNumber(11638560)
  $pb.PbList<$core.String> get tokendomains => $_getList(0);

  @$pb.TagNumber(42675585)
  $pb.PbList<Rule> get rules => $_getList(1);

  @$pb.TagNumber(42731774)
  $pb.PbMap<$core.String, CustomResponseBody> get customresponsebodies =>
      $_getMap(2);

  @$pb.TagNumber(48990889)
  ChallengeConfig get challengeconfig => $_getN(3);
  @$pb.TagNumber(48990889)
  set challengeconfig(ChallengeConfig value) => $_setField(48990889, value);
  @$pb.TagNumber(48990889)
  $core.bool hasChallengeconfig() => $_has(3);
  @$pb.TagNumber(48990889)
  void clearChallengeconfig() => $_clearField(48990889);
  @$pb.TagNumber(48990889)
  ChallengeConfig ensureChallengeconfig() => $_ensure(3);

  @$pb.TagNumber(60547064)
  CaptchaConfig get captchaconfig => $_getN(4);
  @$pb.TagNumber(60547064)
  set captchaconfig(CaptchaConfig value) => $_setField(60547064, value);
  @$pb.TagNumber(60547064)
  $core.bool hasCaptchaconfig() => $_has(4);
  @$pb.TagNumber(60547064)
  void clearCaptchaconfig() => $_clearField(60547064);
  @$pb.TagNumber(60547064)
  CaptchaConfig ensureCaptchaconfig() => $_ensure(4);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(5);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(5);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(105229063)
  OnSourceDDoSProtectionConfig get onsourceddosprotectionconfig => $_getN(6);
  @$pb.TagNumber(105229063)
  set onsourceddosprotectionconfig(OnSourceDDoSProtectionConfig value) =>
      $_setField(105229063, value);
  @$pb.TagNumber(105229063)
  $core.bool hasOnsourceddosprotectionconfig() => $_has(6);
  @$pb.TagNumber(105229063)
  void clearOnsourceddosprotectionconfig() => $_clearField(105229063);
  @$pb.TagNumber(105229063)
  OnSourceDDoSProtectionConfig ensureOnsourceddosprotectionconfig() =>
      $_ensure(6);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(7);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(7, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(7);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(8);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(8, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(8);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(322663861)
  DefaultAction get defaultaction => $_getN(9);
  @$pb.TagNumber(322663861)
  set defaultaction(DefaultAction value) => $_setField(322663861, value);
  @$pb.TagNumber(322663861)
  $core.bool hasDefaultaction() => $_has(9);
  @$pb.TagNumber(322663861)
  void clearDefaultaction() => $_clearField(322663861);
  @$pb.TagNumber(322663861)
  DefaultAction ensureDefaultaction() => $_ensure(9);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(10);

  @$pb.TagNumber(419675691)
  AssociationConfig get associationconfig => $_getN(11);
  @$pb.TagNumber(419675691)
  set associationconfig(AssociationConfig value) =>
      $_setField(419675691, value);
  @$pb.TagNumber(419675691)
  $core.bool hasAssociationconfig() => $_has(11);
  @$pb.TagNumber(419675691)
  void clearAssociationconfig() => $_clearField(419675691);
  @$pb.TagNumber(419675691)
  AssociationConfig ensureAssociationconfig() => $_ensure(11);

  @$pb.TagNumber(464792245)
  DataProtectionConfig get dataprotectionconfig => $_getN(12);
  @$pb.TagNumber(464792245)
  set dataprotectionconfig(DataProtectionConfig value) =>
      $_setField(464792245, value);
  @$pb.TagNumber(464792245)
  $core.bool hasDataprotectionconfig() => $_has(12);
  @$pb.TagNumber(464792245)
  void clearDataprotectionconfig() => $_clearField(464792245);
  @$pb.TagNumber(464792245)
  DataProtectionConfig ensureDataprotectionconfig() => $_ensure(12);

  @$pb.TagNumber(501131976)
  ApplicationConfig get applicationconfig => $_getN(13);
  @$pb.TagNumber(501131976)
  set applicationconfig(ApplicationConfig value) =>
      $_setField(501131976, value);
  @$pb.TagNumber(501131976)
  $core.bool hasApplicationconfig() => $_has(13);
  @$pb.TagNumber(501131976)
  void clearApplicationconfig() => $_clearField(501131976);
  @$pb.TagNumber(501131976)
  ApplicationConfig ensureApplicationconfig() => $_ensure(13);

  @$pb.TagNumber(507762728)
  VisibilityConfig get visibilityconfig => $_getN(14);
  @$pb.TagNumber(507762728)
  set visibilityconfig(VisibilityConfig value) => $_setField(507762728, value);
  @$pb.TagNumber(507762728)
  $core.bool hasVisibilityconfig() => $_has(14);
  @$pb.TagNumber(507762728)
  void clearVisibilityconfig() => $_clearField(507762728);
  @$pb.TagNumber(507762728)
  VisibilityConfig ensureVisibilityconfig() => $_ensure(14);
}

class CreateWebACLResponse extends $pb.GeneratedMessage {
  factory CreateWebACLResponse({
    WebACLSummary? summary,
  }) {
    final result = create();
    if (summary != null) result.summary = summary;
    return result;
  }

  CreateWebACLResponse._();

  factory CreateWebACLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateWebACLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateWebACLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<WebACLSummary>(21935540, _omitFieldNames ? '' : 'summary',
        subBuilder: WebACLSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWebACLResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWebACLResponse copyWith(void Function(CreateWebACLResponse) updates) =>
      super.copyWith((message) => updates(message as CreateWebACLResponse))
          as CreateWebACLResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWebACLResponse create() => CreateWebACLResponse._();
  @$core.override
  CreateWebACLResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateWebACLResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateWebACLResponse>(create);
  static CreateWebACLResponse? _defaultInstance;

  @$pb.TagNumber(21935540)
  WebACLSummary get summary => $_getN(0);
  @$pb.TagNumber(21935540)
  set summary(WebACLSummary value) => $_setField(21935540, value);
  @$pb.TagNumber(21935540)
  $core.bool hasSummary() => $_has(0);
  @$pb.TagNumber(21935540)
  void clearSummary() => $_clearField(21935540);
  @$pb.TagNumber(21935540)
  WebACLSummary ensureSummary() => $_ensure(0);
}

class CustomHTTPHeader extends $pb.GeneratedMessage {
  factory CustomHTTPHeader({
    $core.String? name,
    $core.String? value,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (value != null) result.value = value;
    return result;
  }

  CustomHTTPHeader._();

  factory CustomHTTPHeader.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomHTTPHeader.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomHTTPHeader',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomHTTPHeader clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomHTTPHeader copyWith(void Function(CustomHTTPHeader) updates) =>
      super.copyWith((message) => updates(message as CustomHTTPHeader))
          as CustomHTTPHeader;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomHTTPHeader create() => CustomHTTPHeader._();
  @$core.override
  CustomHTTPHeader createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomHTTPHeader getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CustomHTTPHeader>(create);
  static CustomHTTPHeader? _defaultInstance;

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
}

class CustomRequestHandling extends $pb.GeneratedMessage {
  factory CustomRequestHandling({
    $core.Iterable<CustomHTTPHeader>? insertheaders,
  }) {
    final result = create();
    if (insertheaders != null) result.insertheaders.addAll(insertheaders);
    return result;
  }

  CustomRequestHandling._();

  factory CustomRequestHandling.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomRequestHandling.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomRequestHandling',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<CustomHTTPHeader>(161154427, _omitFieldNames ? '' : 'insertheaders',
        subBuilder: CustomHTTPHeader.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomRequestHandling clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomRequestHandling copyWith(
          void Function(CustomRequestHandling) updates) =>
      super.copyWith((message) => updates(message as CustomRequestHandling))
          as CustomRequestHandling;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomRequestHandling create() => CustomRequestHandling._();
  @$core.override
  CustomRequestHandling createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomRequestHandling getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CustomRequestHandling>(create);
  static CustomRequestHandling? _defaultInstance;

  @$pb.TagNumber(161154427)
  $pb.PbList<CustomHTTPHeader> get insertheaders => $_getList(0);
}

class CustomResponse extends $pb.GeneratedMessage {
  factory CustomResponse({
    $core.String? customresponsebodykey,
    $core.Iterable<CustomHTTPHeader>? responseheaders,
    $core.int? responsecode,
  }) {
    final result = create();
    if (customresponsebodykey != null)
      result.customresponsebodykey = customresponsebodykey;
    if (responseheaders != null) result.responseheaders.addAll(responseheaders);
    if (responsecode != null) result.responsecode = responsecode;
    return result;
  }

  CustomResponse._();

  factory CustomResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(77361257, _omitFieldNames ? '' : 'customresponsebodykey')
    ..pPM<CustomHTTPHeader>(171647029, _omitFieldNames ? '' : 'responseheaders',
        subBuilder: CustomHTTPHeader.create)
    ..aI(447553700, _omitFieldNames ? '' : 'responsecode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomResponse copyWith(void Function(CustomResponse) updates) =>
      super.copyWith((message) => updates(message as CustomResponse))
          as CustomResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomResponse create() => CustomResponse._();
  @$core.override
  CustomResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CustomResponse>(create);
  static CustomResponse? _defaultInstance;

  @$pb.TagNumber(77361257)
  $core.String get customresponsebodykey => $_getSZ(0);
  @$pb.TagNumber(77361257)
  set customresponsebodykey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(77361257)
  $core.bool hasCustomresponsebodykey() => $_has(0);
  @$pb.TagNumber(77361257)
  void clearCustomresponsebodykey() => $_clearField(77361257);

  @$pb.TagNumber(171647029)
  $pb.PbList<CustomHTTPHeader> get responseheaders => $_getList(1);

  @$pb.TagNumber(447553700)
  $core.int get responsecode => $_getIZ(2);
  @$pb.TagNumber(447553700)
  set responsecode($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(447553700)
  $core.bool hasResponsecode() => $_has(2);
  @$pb.TagNumber(447553700)
  void clearResponsecode() => $_clearField(447553700);
}

class CustomResponseBody extends $pb.GeneratedMessage {
  factory CustomResponseBody({
    $core.String? content,
    ResponseContentType? contenttype,
  }) {
    final result = create();
    if (content != null) result.content = content;
    if (contenttype != null) result.contenttype = contenttype;
    return result;
  }

  CustomResponseBody._();

  factory CustomResponseBody.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CustomResponseBody.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CustomResponseBody',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(23568227, _omitFieldNames ? '' : 'content')
    ..aE<ResponseContentType>(333064851, _omitFieldNames ? '' : 'contenttype',
        enumValues: ResponseContentType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomResponseBody clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CustomResponseBody copyWith(void Function(CustomResponseBody) updates) =>
      super.copyWith((message) => updates(message as CustomResponseBody))
          as CustomResponseBody;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CustomResponseBody create() => CustomResponseBody._();
  @$core.override
  CustomResponseBody createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CustomResponseBody getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CustomResponseBody>(create);
  static CustomResponseBody? _defaultInstance;

  @$pb.TagNumber(23568227)
  $core.String get content => $_getSZ(0);
  @$pb.TagNumber(23568227)
  set content($core.String value) => $_setString(0, value);
  @$pb.TagNumber(23568227)
  $core.bool hasContent() => $_has(0);
  @$pb.TagNumber(23568227)
  void clearContent() => $_clearField(23568227);

  @$pb.TagNumber(333064851)
  ResponseContentType get contenttype => $_getN(1);
  @$pb.TagNumber(333064851)
  set contenttype(ResponseContentType value) => $_setField(333064851, value);
  @$pb.TagNumber(333064851)
  $core.bool hasContenttype() => $_has(1);
  @$pb.TagNumber(333064851)
  void clearContenttype() => $_clearField(333064851);
}

class DataProtection extends $pb.GeneratedMessage {
  factory DataProtection({
    DataProtectionAction? action,
    $core.bool? excluderatebaseddetails,
    FieldToProtect? field_263732488,
    $core.bool? excluderulematchdetails,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (excluderatebaseddetails != null)
      result.excluderatebaseddetails = excluderatebaseddetails;
    if (field_263732488 != null) result.field_263732488 = field_263732488;
    if (excluderulematchdetails != null)
      result.excluderulematchdetails = excluderulematchdetails;
    return result;
  }

  DataProtection._();

  factory DataProtection.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataProtection.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataProtection',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<DataProtectionAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: DataProtectionAction.values)
    ..aOB(254694493, _omitFieldNames ? '' : 'excluderatebaseddetails')
    ..aOM<FieldToProtect>(263732488, _omitFieldNames ? '' : 'field',
        subBuilder: FieldToProtect.create)
    ..aOB(275855327, _omitFieldNames ? '' : 'excluderulematchdetails')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataProtection clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataProtection copyWith(void Function(DataProtection) updates) =>
      super.copyWith((message) => updates(message as DataProtection))
          as DataProtection;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataProtection create() => DataProtection._();
  @$core.override
  DataProtection createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataProtection getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataProtection>(create);
  static DataProtection? _defaultInstance;

  @$pb.TagNumber(175614240)
  DataProtectionAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(DataProtectionAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(254694493)
  $core.bool get excluderatebaseddetails => $_getBF(1);
  @$pb.TagNumber(254694493)
  set excluderatebaseddetails($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(254694493)
  $core.bool hasExcluderatebaseddetails() => $_has(1);
  @$pb.TagNumber(254694493)
  void clearExcluderatebaseddetails() => $_clearField(254694493);

  @$pb.TagNumber(263732488)
  FieldToProtect get field_263732488 => $_getN(2);
  @$pb.TagNumber(263732488)
  set field_263732488(FieldToProtect value) => $_setField(263732488, value);
  @$pb.TagNumber(263732488)
  $core.bool hasField_263732488() => $_has(2);
  @$pb.TagNumber(263732488)
  void clearField_263732488() => $_clearField(263732488);
  @$pb.TagNumber(263732488)
  FieldToProtect ensureField_263732488() => $_ensure(2);

  @$pb.TagNumber(275855327)
  $core.bool get excluderulematchdetails => $_getBF(3);
  @$pb.TagNumber(275855327)
  set excluderulematchdetails($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(275855327)
  $core.bool hasExcluderulematchdetails() => $_has(3);
  @$pb.TagNumber(275855327)
  void clearExcluderulematchdetails() => $_clearField(275855327);
}

class DataProtectionConfig extends $pb.GeneratedMessage {
  factory DataProtectionConfig({
    $core.Iterable<DataProtection>? dataprotections,
  }) {
    final result = create();
    if (dataprotections != null) result.dataprotections.addAll(dataprotections);
    return result;
  }

  DataProtectionConfig._();

  factory DataProtectionConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DataProtectionConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DataProtectionConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<DataProtection>(168389196, _omitFieldNames ? '' : 'dataprotections',
        subBuilder: DataProtection.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataProtectionConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DataProtectionConfig copyWith(void Function(DataProtectionConfig) updates) =>
      super.copyWith((message) => updates(message as DataProtectionConfig))
          as DataProtectionConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DataProtectionConfig create() => DataProtectionConfig._();
  @$core.override
  DataProtectionConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DataProtectionConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DataProtectionConfig>(create);
  static DataProtectionConfig? _defaultInstance;

  @$pb.TagNumber(168389196)
  $pb.PbList<DataProtection> get dataprotections => $_getList(0);
}

class DefaultAction extends $pb.GeneratedMessage {
  factory DefaultAction({
    AllowAction? allow,
    BlockAction? block,
  }) {
    final result = create();
    if (allow != null) result.allow = allow;
    if (block != null) result.block = block;
    return result;
  }

  DefaultAction._();

  factory DefaultAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DefaultAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DefaultAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<AllowAction>(342694355, _omitFieldNames ? '' : 'allow',
        subBuilder: AllowAction.create)
    ..aOM<BlockAction>(487014403, _omitFieldNames ? '' : 'block',
        subBuilder: BlockAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DefaultAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DefaultAction copyWith(void Function(DefaultAction) updates) =>
      super.copyWith((message) => updates(message as DefaultAction))
          as DefaultAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DefaultAction create() => DefaultAction._();
  @$core.override
  DefaultAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DefaultAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DefaultAction>(create);
  static DefaultAction? _defaultInstance;

  @$pb.TagNumber(342694355)
  AllowAction get allow => $_getN(0);
  @$pb.TagNumber(342694355)
  set allow(AllowAction value) => $_setField(342694355, value);
  @$pb.TagNumber(342694355)
  $core.bool hasAllow() => $_has(0);
  @$pb.TagNumber(342694355)
  void clearAllow() => $_clearField(342694355);
  @$pb.TagNumber(342694355)
  AllowAction ensureAllow() => $_ensure(0);

  @$pb.TagNumber(487014403)
  BlockAction get block => $_getN(1);
  @$pb.TagNumber(487014403)
  set block(BlockAction value) => $_setField(487014403, value);
  @$pb.TagNumber(487014403)
  $core.bool hasBlock() => $_has(1);
  @$pb.TagNumber(487014403)
  void clearBlock() => $_clearField(487014403);
  @$pb.TagNumber(487014403)
  BlockAction ensureBlock() => $_ensure(1);
}

class DeleteAPIKeyRequest extends $pb.GeneratedMessage {
  factory DeleteAPIKeyRequest({
    Scope? scope,
    $core.String? apikey,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (apikey != null) result.apikey = apikey;
    return result;
  }

  DeleteAPIKeyRequest._();

  factory DeleteAPIKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteAPIKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAPIKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(274818239, _omitFieldNames ? '' : 'apikey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAPIKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAPIKeyRequest copyWith(void Function(DeleteAPIKeyRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteAPIKeyRequest))
          as DeleteAPIKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAPIKeyRequest create() => DeleteAPIKeyRequest._();
  @$core.override
  DeleteAPIKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteAPIKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteAPIKeyRequest>(create);
  static DeleteAPIKeyRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(274818239)
  $core.String get apikey => $_getSZ(1);
  @$pb.TagNumber(274818239)
  set apikey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(274818239)
  $core.bool hasApikey() => $_has(1);
  @$pb.TagNumber(274818239)
  void clearApikey() => $_clearField(274818239);
}

class DeleteAPIKeyResponse extends $pb.GeneratedMessage {
  factory DeleteAPIKeyResponse() => create();

  DeleteAPIKeyResponse._();

  factory DeleteAPIKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteAPIKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAPIKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAPIKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAPIKeyResponse copyWith(void Function(DeleteAPIKeyResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteAPIKeyResponse))
          as DeleteAPIKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAPIKeyResponse create() => DeleteAPIKeyResponse._();
  @$core.override
  DeleteAPIKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteAPIKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteAPIKeyResponse>(create);
  static DeleteAPIKeyResponse? _defaultInstance;
}

class DeleteFirewallManagerRuleGroupsRequest extends $pb.GeneratedMessage {
  factory DeleteFirewallManagerRuleGroupsRequest({
    $core.String? webaclarn,
    $core.String? webacllocktoken,
  }) {
    final result = create();
    if (webaclarn != null) result.webaclarn = webaclarn;
    if (webacllocktoken != null) result.webacllocktoken = webacllocktoken;
    return result;
  }

  DeleteFirewallManagerRuleGroupsRequest._();

  factory DeleteFirewallManagerRuleGroupsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteFirewallManagerRuleGroupsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteFirewallManagerRuleGroupsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(82506659, _omitFieldNames ? '' : 'webaclarn')
    ..aOS(229622606, _omitFieldNames ? '' : 'webacllocktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteFirewallManagerRuleGroupsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteFirewallManagerRuleGroupsRequest copyWith(
          void Function(DeleteFirewallManagerRuleGroupsRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteFirewallManagerRuleGroupsRequest))
          as DeleteFirewallManagerRuleGroupsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteFirewallManagerRuleGroupsRequest create() =>
      DeleteFirewallManagerRuleGroupsRequest._();
  @$core.override
  DeleteFirewallManagerRuleGroupsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteFirewallManagerRuleGroupsRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeleteFirewallManagerRuleGroupsRequest>(create);
  static DeleteFirewallManagerRuleGroupsRequest? _defaultInstance;

  @$pb.TagNumber(82506659)
  $core.String get webaclarn => $_getSZ(0);
  @$pb.TagNumber(82506659)
  set webaclarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82506659)
  $core.bool hasWebaclarn() => $_has(0);
  @$pb.TagNumber(82506659)
  void clearWebaclarn() => $_clearField(82506659);

  @$pb.TagNumber(229622606)
  $core.String get webacllocktoken => $_getSZ(1);
  @$pb.TagNumber(229622606)
  set webacllocktoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(229622606)
  $core.bool hasWebacllocktoken() => $_has(1);
  @$pb.TagNumber(229622606)
  void clearWebacllocktoken() => $_clearField(229622606);
}

class DeleteFirewallManagerRuleGroupsResponse extends $pb.GeneratedMessage {
  factory DeleteFirewallManagerRuleGroupsResponse({
    $core.String? nextwebacllocktoken,
  }) {
    final result = create();
    if (nextwebacllocktoken != null)
      result.nextwebacllocktoken = nextwebacllocktoken;
    return result;
  }

  DeleteFirewallManagerRuleGroupsResponse._();

  factory DeleteFirewallManagerRuleGroupsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteFirewallManagerRuleGroupsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteFirewallManagerRuleGroupsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(89595355, _omitFieldNames ? '' : 'nextwebacllocktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteFirewallManagerRuleGroupsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteFirewallManagerRuleGroupsResponse copyWith(
          void Function(DeleteFirewallManagerRuleGroupsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteFirewallManagerRuleGroupsResponse))
          as DeleteFirewallManagerRuleGroupsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteFirewallManagerRuleGroupsResponse create() =>
      DeleteFirewallManagerRuleGroupsResponse._();
  @$core.override
  DeleteFirewallManagerRuleGroupsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteFirewallManagerRuleGroupsResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeleteFirewallManagerRuleGroupsResponse>(create);
  static DeleteFirewallManagerRuleGroupsResponse? _defaultInstance;

  @$pb.TagNumber(89595355)
  $core.String get nextwebacllocktoken => $_getSZ(0);
  @$pb.TagNumber(89595355)
  set nextwebacllocktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89595355)
  $core.bool hasNextwebacllocktoken() => $_has(0);
  @$pb.TagNumber(89595355)
  void clearNextwebacllocktoken() => $_clearField(89595355);
}

class DeleteIPSetRequest extends $pb.GeneratedMessage {
  factory DeleteIPSetRequest({
    $core.String? locktoken,
    Scope? scope,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  DeleteIPSetRequest._();

  factory DeleteIPSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteIPSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteIPSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIPSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIPSetRequest copyWith(void Function(DeleteIPSetRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteIPSetRequest))
          as DeleteIPSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteIPSetRequest create() => DeleteIPSetRequest._();
  @$core.override
  DeleteIPSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteIPSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteIPSetRequest>(create);
  static DeleteIPSetRequest? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteIPSetResponse extends $pb.GeneratedMessage {
  factory DeleteIPSetResponse() => create();

  DeleteIPSetResponse._();

  factory DeleteIPSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteIPSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteIPSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIPSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteIPSetResponse copyWith(void Function(DeleteIPSetResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteIPSetResponse))
          as DeleteIPSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteIPSetResponse create() => DeleteIPSetResponse._();
  @$core.override
  DeleteIPSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteIPSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteIPSetResponse>(create);
  static DeleteIPSetResponse? _defaultInstance;
}

class DeleteLoggingConfigurationRequest extends $pb.GeneratedMessage {
  factory DeleteLoggingConfigurationRequest({
    LogType? logtype,
    LogScope? logscope,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (logtype != null) result.logtype = logtype;
    if (logscope != null) result.logscope = logscope;
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  DeleteLoggingConfigurationRequest._();

  factory DeleteLoggingConfigurationRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteLoggingConfigurationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteLoggingConfigurationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<LogType>(116703930, _omitFieldNames ? '' : 'logtype',
        enumValues: LogType.values)
    ..aE<LogScope>(188235840, _omitFieldNames ? '' : 'logscope',
        enumValues: LogScope.values)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteLoggingConfigurationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteLoggingConfigurationRequest copyWith(
          void Function(DeleteLoggingConfigurationRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteLoggingConfigurationRequest))
          as DeleteLoggingConfigurationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteLoggingConfigurationRequest create() =>
      DeleteLoggingConfigurationRequest._();
  @$core.override
  DeleteLoggingConfigurationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteLoggingConfigurationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteLoggingConfigurationRequest>(
          create);
  static DeleteLoggingConfigurationRequest? _defaultInstance;

  @$pb.TagNumber(116703930)
  LogType get logtype => $_getN(0);
  @$pb.TagNumber(116703930)
  set logtype(LogType value) => $_setField(116703930, value);
  @$pb.TagNumber(116703930)
  $core.bool hasLogtype() => $_has(0);
  @$pb.TagNumber(116703930)
  void clearLogtype() => $_clearField(116703930);

  @$pb.TagNumber(188235840)
  LogScope get logscope => $_getN(1);
  @$pb.TagNumber(188235840)
  set logscope(LogScope value) => $_setField(188235840, value);
  @$pb.TagNumber(188235840)
  $core.bool hasLogscope() => $_has(1);
  @$pb.TagNumber(188235840)
  void clearLogscope() => $_clearField(188235840);

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(2);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(2);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class DeleteLoggingConfigurationResponse extends $pb.GeneratedMessage {
  factory DeleteLoggingConfigurationResponse() => create();

  DeleteLoggingConfigurationResponse._();

  factory DeleteLoggingConfigurationResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteLoggingConfigurationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteLoggingConfigurationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteLoggingConfigurationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteLoggingConfigurationResponse copyWith(
          void Function(DeleteLoggingConfigurationResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteLoggingConfigurationResponse))
          as DeleteLoggingConfigurationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteLoggingConfigurationResponse create() =>
      DeleteLoggingConfigurationResponse._();
  @$core.override
  DeleteLoggingConfigurationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteLoggingConfigurationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteLoggingConfigurationResponse>(
          create);
  static DeleteLoggingConfigurationResponse? _defaultInstance;
}

class DeletePermissionPolicyRequest extends $pb.GeneratedMessage {
  factory DeletePermissionPolicyRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  DeletePermissionPolicyRequest._();

  factory DeletePermissionPolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeletePermissionPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeletePermissionPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePermissionPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePermissionPolicyRequest copyWith(
          void Function(DeletePermissionPolicyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeletePermissionPolicyRequest))
          as DeletePermissionPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeletePermissionPolicyRequest create() =>
      DeletePermissionPolicyRequest._();
  @$core.override
  DeletePermissionPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeletePermissionPolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeletePermissionPolicyRequest>(create);
  static DeletePermissionPolicyRequest? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class DeletePermissionPolicyResponse extends $pb.GeneratedMessage {
  factory DeletePermissionPolicyResponse() => create();

  DeletePermissionPolicyResponse._();

  factory DeletePermissionPolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeletePermissionPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeletePermissionPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePermissionPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePermissionPolicyResponse copyWith(
          void Function(DeletePermissionPolicyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeletePermissionPolicyResponse))
          as DeletePermissionPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeletePermissionPolicyResponse create() =>
      DeletePermissionPolicyResponse._();
  @$core.override
  DeletePermissionPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeletePermissionPolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeletePermissionPolicyResponse>(create);
  static DeletePermissionPolicyResponse? _defaultInstance;
}

class DeleteRegexPatternSetRequest extends $pb.GeneratedMessage {
  factory DeleteRegexPatternSetRequest({
    $core.String? locktoken,
    Scope? scope,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  DeleteRegexPatternSetRequest._();

  factory DeleteRegexPatternSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRegexPatternSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRegexPatternSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRegexPatternSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRegexPatternSetRequest copyWith(
          void Function(DeleteRegexPatternSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteRegexPatternSetRequest))
          as DeleteRegexPatternSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRegexPatternSetRequest create() =>
      DeleteRegexPatternSetRequest._();
  @$core.override
  DeleteRegexPatternSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRegexPatternSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRegexPatternSetRequest>(create);
  static DeleteRegexPatternSetRequest? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteRegexPatternSetResponse extends $pb.GeneratedMessage {
  factory DeleteRegexPatternSetResponse() => create();

  DeleteRegexPatternSetResponse._();

  factory DeleteRegexPatternSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRegexPatternSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRegexPatternSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRegexPatternSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRegexPatternSetResponse copyWith(
          void Function(DeleteRegexPatternSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteRegexPatternSetResponse))
          as DeleteRegexPatternSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRegexPatternSetResponse create() =>
      DeleteRegexPatternSetResponse._();
  @$core.override
  DeleteRegexPatternSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRegexPatternSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRegexPatternSetResponse>(create);
  static DeleteRegexPatternSetResponse? _defaultInstance;
}

class DeleteRuleGroupRequest extends $pb.GeneratedMessage {
  factory DeleteRuleGroupRequest({
    $core.String? locktoken,
    Scope? scope,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  DeleteRuleGroupRequest._();

  factory DeleteRuleGroupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRuleGroupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRuleGroupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRuleGroupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRuleGroupRequest copyWith(
          void Function(DeleteRuleGroupRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteRuleGroupRequest))
          as DeleteRuleGroupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRuleGroupRequest create() => DeleteRuleGroupRequest._();
  @$core.override
  DeleteRuleGroupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRuleGroupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRuleGroupRequest>(create);
  static DeleteRuleGroupRequest? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteRuleGroupResponse extends $pb.GeneratedMessage {
  factory DeleteRuleGroupResponse() => create();

  DeleteRuleGroupResponse._();

  factory DeleteRuleGroupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRuleGroupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRuleGroupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRuleGroupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRuleGroupResponse copyWith(
          void Function(DeleteRuleGroupResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteRuleGroupResponse))
          as DeleteRuleGroupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRuleGroupResponse create() => DeleteRuleGroupResponse._();
  @$core.override
  DeleteRuleGroupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRuleGroupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRuleGroupResponse>(create);
  static DeleteRuleGroupResponse? _defaultInstance;
}

class DeleteWebACLRequest extends $pb.GeneratedMessage {
  factory DeleteWebACLRequest({
    $core.String? locktoken,
    Scope? scope,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  DeleteWebACLRequest._();

  factory DeleteWebACLRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteWebACLRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteWebACLRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteWebACLRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteWebACLRequest copyWith(void Function(DeleteWebACLRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteWebACLRequest))
          as DeleteWebACLRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteWebACLRequest create() => DeleteWebACLRequest._();
  @$core.override
  DeleteWebACLRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteWebACLRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteWebACLRequest>(create);
  static DeleteWebACLRequest? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteWebACLResponse extends $pb.GeneratedMessage {
  factory DeleteWebACLResponse() => create();

  DeleteWebACLResponse._();

  factory DeleteWebACLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteWebACLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteWebACLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteWebACLResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteWebACLResponse copyWith(void Function(DeleteWebACLResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteWebACLResponse))
          as DeleteWebACLResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteWebACLResponse create() => DeleteWebACLResponse._();
  @$core.override
  DeleteWebACLResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteWebACLResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteWebACLResponse>(create);
  static DeleteWebACLResponse? _defaultInstance;
}

class DescribeAllManagedProductsRequest extends $pb.GeneratedMessage {
  factory DescribeAllManagedProductsRequest({
    Scope? scope,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    return result;
  }

  DescribeAllManagedProductsRequest._();

  factory DescribeAllManagedProductsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAllManagedProductsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAllManagedProductsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAllManagedProductsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAllManagedProductsRequest copyWith(
          void Function(DescribeAllManagedProductsRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeAllManagedProductsRequest))
          as DescribeAllManagedProductsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAllManagedProductsRequest create() =>
      DescribeAllManagedProductsRequest._();
  @$core.override
  DescribeAllManagedProductsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAllManagedProductsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAllManagedProductsRequest>(
          create);
  static DescribeAllManagedProductsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);
}

class DescribeAllManagedProductsResponse extends $pb.GeneratedMessage {
  factory DescribeAllManagedProductsResponse({
    $core.Iterable<ManagedProductDescriptor>? managedproducts,
  }) {
    final result = create();
    if (managedproducts != null) result.managedproducts.addAll(managedproducts);
    return result;
  }

  DescribeAllManagedProductsResponse._();

  factory DescribeAllManagedProductsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAllManagedProductsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAllManagedProductsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<ManagedProductDescriptor>(
        245391099, _omitFieldNames ? '' : 'managedproducts',
        subBuilder: ManagedProductDescriptor.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAllManagedProductsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAllManagedProductsResponse copyWith(
          void Function(DescribeAllManagedProductsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeAllManagedProductsResponse))
          as DescribeAllManagedProductsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAllManagedProductsResponse create() =>
      DescribeAllManagedProductsResponse._();
  @$core.override
  DescribeAllManagedProductsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAllManagedProductsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAllManagedProductsResponse>(
          create);
  static DescribeAllManagedProductsResponse? _defaultInstance;

  @$pb.TagNumber(245391099)
  $pb.PbList<ManagedProductDescriptor> get managedproducts => $_getList(0);
}

class DescribeManagedProductsByVendorRequest extends $pb.GeneratedMessage {
  factory DescribeManagedProductsByVendorRequest({
    Scope? scope,
    $core.String? vendorname,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (vendorname != null) result.vendorname = vendorname;
    return result;
  }

  DescribeManagedProductsByVendorRequest._();

  factory DescribeManagedProductsByVendorRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeManagedProductsByVendorRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeManagedProductsByVendorRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(159329689, _omitFieldNames ? '' : 'vendorname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeManagedProductsByVendorRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeManagedProductsByVendorRequest copyWith(
          void Function(DescribeManagedProductsByVendorRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeManagedProductsByVendorRequest))
          as DescribeManagedProductsByVendorRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeManagedProductsByVendorRequest create() =>
      DescribeManagedProductsByVendorRequest._();
  @$core.override
  DescribeManagedProductsByVendorRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeManagedProductsByVendorRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DescribeManagedProductsByVendorRequest>(create);
  static DescribeManagedProductsByVendorRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(159329689)
  $core.String get vendorname => $_getSZ(1);
  @$pb.TagNumber(159329689)
  set vendorname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(159329689)
  $core.bool hasVendorname() => $_has(1);
  @$pb.TagNumber(159329689)
  void clearVendorname() => $_clearField(159329689);
}

class DescribeManagedProductsByVendorResponse extends $pb.GeneratedMessage {
  factory DescribeManagedProductsByVendorResponse({
    $core.Iterable<ManagedProductDescriptor>? managedproducts,
  }) {
    final result = create();
    if (managedproducts != null) result.managedproducts.addAll(managedproducts);
    return result;
  }

  DescribeManagedProductsByVendorResponse._();

  factory DescribeManagedProductsByVendorResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeManagedProductsByVendorResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeManagedProductsByVendorResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<ManagedProductDescriptor>(
        245391099, _omitFieldNames ? '' : 'managedproducts',
        subBuilder: ManagedProductDescriptor.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeManagedProductsByVendorResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeManagedProductsByVendorResponse copyWith(
          void Function(DescribeManagedProductsByVendorResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeManagedProductsByVendorResponse))
          as DescribeManagedProductsByVendorResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeManagedProductsByVendorResponse create() =>
      DescribeManagedProductsByVendorResponse._();
  @$core.override
  DescribeManagedProductsByVendorResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeManagedProductsByVendorResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DescribeManagedProductsByVendorResponse>(create);
  static DescribeManagedProductsByVendorResponse? _defaultInstance;

  @$pb.TagNumber(245391099)
  $pb.PbList<ManagedProductDescriptor> get managedproducts => $_getList(0);
}

class DescribeManagedRuleGroupRequest extends $pb.GeneratedMessage {
  factory DescribeManagedRuleGroupRequest({
    Scope? scope,
    $core.String? vendorname,
    $core.String? versionname,
    $core.String? name,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (vendorname != null) result.vendorname = vendorname;
    if (versionname != null) result.versionname = versionname;
    if (name != null) result.name = name;
    return result;
  }

  DescribeManagedRuleGroupRequest._();

  factory DescribeManagedRuleGroupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeManagedRuleGroupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeManagedRuleGroupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(159329689, _omitFieldNames ? '' : 'vendorname')
    ..aOS(227348949, _omitFieldNames ? '' : 'versionname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeManagedRuleGroupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeManagedRuleGroupRequest copyWith(
          void Function(DescribeManagedRuleGroupRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeManagedRuleGroupRequest))
          as DescribeManagedRuleGroupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeManagedRuleGroupRequest create() =>
      DescribeManagedRuleGroupRequest._();
  @$core.override
  DescribeManagedRuleGroupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeManagedRuleGroupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeManagedRuleGroupRequest>(
          create);
  static DescribeManagedRuleGroupRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(159329689)
  $core.String get vendorname => $_getSZ(1);
  @$pb.TagNumber(159329689)
  set vendorname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(159329689)
  $core.bool hasVendorname() => $_has(1);
  @$pb.TagNumber(159329689)
  void clearVendorname() => $_clearField(159329689);

  @$pb.TagNumber(227348949)
  $core.String get versionname => $_getSZ(2);
  @$pb.TagNumber(227348949)
  set versionname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(227348949)
  $core.bool hasVersionname() => $_has(2);
  @$pb.TagNumber(227348949)
  void clearVersionname() => $_clearField(227348949);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DescribeManagedRuleGroupResponse extends $pb.GeneratedMessage {
  factory DescribeManagedRuleGroupResponse({
    $core.String? labelnamespace,
    $core.Iterable<RuleSummary>? rules,
    $core.Iterable<LabelSummary>? consumedlabels,
    $core.Iterable<LabelSummary>? availablelabels,
    $fixnum.Int64? capacity,
    $core.String? versionname,
    $core.String? snstopicarn,
  }) {
    final result = create();
    if (labelnamespace != null) result.labelnamespace = labelnamespace;
    if (rules != null) result.rules.addAll(rules);
    if (consumedlabels != null) result.consumedlabels.addAll(consumedlabels);
    if (availablelabels != null) result.availablelabels.addAll(availablelabels);
    if (capacity != null) result.capacity = capacity;
    if (versionname != null) result.versionname = versionname;
    if (snstopicarn != null) result.snstopicarn = snstopicarn;
    return result;
  }

  DescribeManagedRuleGroupResponse._();

  factory DescribeManagedRuleGroupResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeManagedRuleGroupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeManagedRuleGroupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(39797417, _omitFieldNames ? '' : 'labelnamespace')
    ..pPM<RuleSummary>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: RuleSummary.create)
    ..pPM<LabelSummary>(43813949, _omitFieldNames ? '' : 'consumedlabels',
        subBuilder: LabelSummary.create)
    ..pPM<LabelSummary>(51059984, _omitFieldNames ? '' : 'availablelabels',
        subBuilder: LabelSummary.create)
    ..aInt64(107253930, _omitFieldNames ? '' : 'capacity')
    ..aOS(227348949, _omitFieldNames ? '' : 'versionname')
    ..aOS(374815596, _omitFieldNames ? '' : 'snstopicarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeManagedRuleGroupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeManagedRuleGroupResponse copyWith(
          void Function(DescribeManagedRuleGroupResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeManagedRuleGroupResponse))
          as DescribeManagedRuleGroupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeManagedRuleGroupResponse create() =>
      DescribeManagedRuleGroupResponse._();
  @$core.override
  DescribeManagedRuleGroupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeManagedRuleGroupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeManagedRuleGroupResponse>(
          create);
  static DescribeManagedRuleGroupResponse? _defaultInstance;

  @$pb.TagNumber(39797417)
  $core.String get labelnamespace => $_getSZ(0);
  @$pb.TagNumber(39797417)
  set labelnamespace($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39797417)
  $core.bool hasLabelnamespace() => $_has(0);
  @$pb.TagNumber(39797417)
  void clearLabelnamespace() => $_clearField(39797417);

  @$pb.TagNumber(42675585)
  $pb.PbList<RuleSummary> get rules => $_getList(1);

  @$pb.TagNumber(43813949)
  $pb.PbList<LabelSummary> get consumedlabels => $_getList(2);

  @$pb.TagNumber(51059984)
  $pb.PbList<LabelSummary> get availablelabels => $_getList(3);

  @$pb.TagNumber(107253930)
  $fixnum.Int64 get capacity => $_getI64(4);
  @$pb.TagNumber(107253930)
  set capacity($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(107253930)
  $core.bool hasCapacity() => $_has(4);
  @$pb.TagNumber(107253930)
  void clearCapacity() => $_clearField(107253930);

  @$pb.TagNumber(227348949)
  $core.String get versionname => $_getSZ(5);
  @$pb.TagNumber(227348949)
  set versionname($core.String value) => $_setString(5, value);
  @$pb.TagNumber(227348949)
  $core.bool hasVersionname() => $_has(5);
  @$pb.TagNumber(227348949)
  void clearVersionname() => $_clearField(227348949);

  @$pb.TagNumber(374815596)
  $core.String get snstopicarn => $_getSZ(6);
  @$pb.TagNumber(374815596)
  set snstopicarn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(374815596)
  $core.bool hasSnstopicarn() => $_has(6);
  @$pb.TagNumber(374815596)
  void clearSnstopicarn() => $_clearField(374815596);
}

class DisallowedFeature extends $pb.GeneratedMessage {
  factory DisallowedFeature({
    $core.String? requiredpricingplan,
    $core.String? feature,
  }) {
    final result = create();
    if (requiredpricingplan != null)
      result.requiredpricingplan = requiredpricingplan;
    if (feature != null) result.feature = feature;
    return result;
  }

  DisallowedFeature._();

  factory DisallowedFeature.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisallowedFeature.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisallowedFeature',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(501231534, _omitFieldNames ? '' : 'requiredpricingplan')
    ..aOS(512819934, _omitFieldNames ? '' : 'feature')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisallowedFeature clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisallowedFeature copyWith(void Function(DisallowedFeature) updates) =>
      super.copyWith((message) => updates(message as DisallowedFeature))
          as DisallowedFeature;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisallowedFeature create() => DisallowedFeature._();
  @$core.override
  DisallowedFeature createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisallowedFeature getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisallowedFeature>(create);
  static DisallowedFeature? _defaultInstance;

  @$pb.TagNumber(501231534)
  $core.String get requiredpricingplan => $_getSZ(0);
  @$pb.TagNumber(501231534)
  set requiredpricingplan($core.String value) => $_setString(0, value);
  @$pb.TagNumber(501231534)
  $core.bool hasRequiredpricingplan() => $_has(0);
  @$pb.TagNumber(501231534)
  void clearRequiredpricingplan() => $_clearField(501231534);

  @$pb.TagNumber(512819934)
  $core.String get feature => $_getSZ(1);
  @$pb.TagNumber(512819934)
  set feature($core.String value) => $_setString(1, value);
  @$pb.TagNumber(512819934)
  $core.bool hasFeature() => $_has(1);
  @$pb.TagNumber(512819934)
  void clearFeature() => $_clearField(512819934);
}

class DisassociateWebACLRequest extends $pb.GeneratedMessage {
  factory DisassociateWebACLRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  DisassociateWebACLRequest._();

  factory DisassociateWebACLRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisassociateWebACLRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisassociateWebACLRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisassociateWebACLRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisassociateWebACLRequest copyWith(
          void Function(DisassociateWebACLRequest) updates) =>
      super.copyWith((message) => updates(message as DisassociateWebACLRequest))
          as DisassociateWebACLRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisassociateWebACLRequest create() => DisassociateWebACLRequest._();
  @$core.override
  DisassociateWebACLRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisassociateWebACLRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisassociateWebACLRequest>(create);
  static DisassociateWebACLRequest? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class DisassociateWebACLResponse extends $pb.GeneratedMessage {
  factory DisassociateWebACLResponse() => create();

  DisassociateWebACLResponse._();

  factory DisassociateWebACLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisassociateWebACLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisassociateWebACLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisassociateWebACLResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisassociateWebACLResponse copyWith(
          void Function(DisassociateWebACLResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DisassociateWebACLResponse))
          as DisassociateWebACLResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisassociateWebACLResponse create() => DisassociateWebACLResponse._();
  @$core.override
  DisassociateWebACLResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisassociateWebACLResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisassociateWebACLResponse>(create);
  static DisassociateWebACLResponse? _defaultInstance;
}

class EmailField extends $pb.GeneratedMessage {
  factory EmailField({
    $core.String? identifier,
  }) {
    final result = create();
    if (identifier != null) result.identifier = identifier;
    return result;
  }

  EmailField._();

  factory EmailField.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EmailField.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EmailField',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(41865311, _omitFieldNames ? '' : 'identifier')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EmailField clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EmailField copyWith(void Function(EmailField) updates) =>
      super.copyWith((message) => updates(message as EmailField)) as EmailField;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EmailField create() => EmailField._();
  @$core.override
  EmailField createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EmailField getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EmailField>(create);
  static EmailField? _defaultInstance;

  @$pb.TagNumber(41865311)
  $core.String get identifier => $_getSZ(0);
  @$pb.TagNumber(41865311)
  set identifier($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41865311)
  $core.bool hasIdentifier() => $_has(0);
  @$pb.TagNumber(41865311)
  void clearIdentifier() => $_clearField(41865311);
}

class ExcludedRule extends $pb.GeneratedMessage {
  factory ExcludedRule({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  ExcludedRule._();

  factory ExcludedRule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExcludedRule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExcludedRule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExcludedRule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExcludedRule copyWith(void Function(ExcludedRule) updates) =>
      super.copyWith((message) => updates(message as ExcludedRule))
          as ExcludedRule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExcludedRule create() => ExcludedRule._();
  @$core.override
  ExcludedRule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExcludedRule getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExcludedRule>(create);
  static ExcludedRule? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class FieldToMatch extends $pb.GeneratedMessage {
  factory FieldToMatch({
    UriFragment? urifragment,
    Body? body,
    SingleQueryArgument? singlequeryargument,
    JsonBody? jsonbody,
    AllQueryArguments? allqueryarguments,
    UriPath? uripath,
    Headers? headers,
    JA4Fingerprint? ja4fingerprint,
    Method? method,
    Cookies? cookies,
    HeaderOrder? headerorder,
    SingleHeader? singleheader,
    QueryString? querystring,
    JA3Fingerprint? ja3fingerprint,
  }) {
    final result = create();
    if (urifragment != null) result.urifragment = urifragment;
    if (body != null) result.body = body;
    if (singlequeryargument != null)
      result.singlequeryargument = singlequeryargument;
    if (jsonbody != null) result.jsonbody = jsonbody;
    if (allqueryarguments != null) result.allqueryarguments = allqueryarguments;
    if (uripath != null) result.uripath = uripath;
    if (headers != null) result.headers = headers;
    if (ja4fingerprint != null) result.ja4fingerprint = ja4fingerprint;
    if (method != null) result.method = method;
    if (cookies != null) result.cookies = cookies;
    if (headerorder != null) result.headerorder = headerorder;
    if (singleheader != null) result.singleheader = singleheader;
    if (querystring != null) result.querystring = querystring;
    if (ja3fingerprint != null) result.ja3fingerprint = ja3fingerprint;
    return result;
  }

  FieldToMatch._();

  factory FieldToMatch.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FieldToMatch.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FieldToMatch',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<UriFragment>(11939318, _omitFieldNames ? '' : 'urifragment',
        subBuilder: UriFragment.create)
    ..aOM<Body>(42602646, _omitFieldNames ? '' : 'body',
        subBuilder: Body.create)
    ..aOM<SingleQueryArgument>(
        191839869, _omitFieldNames ? '' : 'singlequeryargument',
        subBuilder: SingleQueryArgument.create)
    ..aOM<JsonBody>(219447052, _omitFieldNames ? '' : 'jsonbody',
        subBuilder: JsonBody.create)
    ..aOM<AllQueryArguments>(
        271813941, _omitFieldNames ? '' : 'allqueryarguments',
        subBuilder: AllQueryArguments.create)
    ..aOM<UriPath>(288340351, _omitFieldNames ? '' : 'uripath',
        subBuilder: UriPath.create)
    ..aOM<Headers>(323967370, _omitFieldNames ? '' : 'headers',
        subBuilder: Headers.create)
    ..aOM<JA4Fingerprint>(352209397, _omitFieldNames ? '' : 'ja4fingerprint',
        subBuilder: JA4Fingerprint.create)
    ..aOM<Method>(413321041, _omitFieldNames ? '' : 'method',
        subBuilder: Method.create)
    ..aOM<Cookies>(418634949, _omitFieldNames ? '' : 'cookies',
        subBuilder: Cookies.create)
    ..aOM<HeaderOrder>(419090489, _omitFieldNames ? '' : 'headerorder',
        subBuilder: HeaderOrder.create)
    ..aOM<SingleHeader>(434029149, _omitFieldNames ? '' : 'singleheader',
        subBuilder: SingleHeader.create)
    ..aOM<QueryString>(435938663, _omitFieldNames ? '' : 'querystring',
        subBuilder: QueryString.create)
    ..aOM<JA3Fingerprint>(521270332, _omitFieldNames ? '' : 'ja3fingerprint',
        subBuilder: JA3Fingerprint.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FieldToMatch clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FieldToMatch copyWith(void Function(FieldToMatch) updates) =>
      super.copyWith((message) => updates(message as FieldToMatch))
          as FieldToMatch;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FieldToMatch create() => FieldToMatch._();
  @$core.override
  FieldToMatch createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FieldToMatch getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FieldToMatch>(create);
  static FieldToMatch? _defaultInstance;

  @$pb.TagNumber(11939318)
  UriFragment get urifragment => $_getN(0);
  @$pb.TagNumber(11939318)
  set urifragment(UriFragment value) => $_setField(11939318, value);
  @$pb.TagNumber(11939318)
  $core.bool hasUrifragment() => $_has(0);
  @$pb.TagNumber(11939318)
  void clearUrifragment() => $_clearField(11939318);
  @$pb.TagNumber(11939318)
  UriFragment ensureUrifragment() => $_ensure(0);

  @$pb.TagNumber(42602646)
  Body get body => $_getN(1);
  @$pb.TagNumber(42602646)
  set body(Body value) => $_setField(42602646, value);
  @$pb.TagNumber(42602646)
  $core.bool hasBody() => $_has(1);
  @$pb.TagNumber(42602646)
  void clearBody() => $_clearField(42602646);
  @$pb.TagNumber(42602646)
  Body ensureBody() => $_ensure(1);

  @$pb.TagNumber(191839869)
  SingleQueryArgument get singlequeryargument => $_getN(2);
  @$pb.TagNumber(191839869)
  set singlequeryargument(SingleQueryArgument value) =>
      $_setField(191839869, value);
  @$pb.TagNumber(191839869)
  $core.bool hasSinglequeryargument() => $_has(2);
  @$pb.TagNumber(191839869)
  void clearSinglequeryargument() => $_clearField(191839869);
  @$pb.TagNumber(191839869)
  SingleQueryArgument ensureSinglequeryargument() => $_ensure(2);

  @$pb.TagNumber(219447052)
  JsonBody get jsonbody => $_getN(3);
  @$pb.TagNumber(219447052)
  set jsonbody(JsonBody value) => $_setField(219447052, value);
  @$pb.TagNumber(219447052)
  $core.bool hasJsonbody() => $_has(3);
  @$pb.TagNumber(219447052)
  void clearJsonbody() => $_clearField(219447052);
  @$pb.TagNumber(219447052)
  JsonBody ensureJsonbody() => $_ensure(3);

  @$pb.TagNumber(271813941)
  AllQueryArguments get allqueryarguments => $_getN(4);
  @$pb.TagNumber(271813941)
  set allqueryarguments(AllQueryArguments value) =>
      $_setField(271813941, value);
  @$pb.TagNumber(271813941)
  $core.bool hasAllqueryarguments() => $_has(4);
  @$pb.TagNumber(271813941)
  void clearAllqueryarguments() => $_clearField(271813941);
  @$pb.TagNumber(271813941)
  AllQueryArguments ensureAllqueryarguments() => $_ensure(4);

  @$pb.TagNumber(288340351)
  UriPath get uripath => $_getN(5);
  @$pb.TagNumber(288340351)
  set uripath(UriPath value) => $_setField(288340351, value);
  @$pb.TagNumber(288340351)
  $core.bool hasUripath() => $_has(5);
  @$pb.TagNumber(288340351)
  void clearUripath() => $_clearField(288340351);
  @$pb.TagNumber(288340351)
  UriPath ensureUripath() => $_ensure(5);

  @$pb.TagNumber(323967370)
  Headers get headers => $_getN(6);
  @$pb.TagNumber(323967370)
  set headers(Headers value) => $_setField(323967370, value);
  @$pb.TagNumber(323967370)
  $core.bool hasHeaders() => $_has(6);
  @$pb.TagNumber(323967370)
  void clearHeaders() => $_clearField(323967370);
  @$pb.TagNumber(323967370)
  Headers ensureHeaders() => $_ensure(6);

  @$pb.TagNumber(352209397)
  JA4Fingerprint get ja4fingerprint => $_getN(7);
  @$pb.TagNumber(352209397)
  set ja4fingerprint(JA4Fingerprint value) => $_setField(352209397, value);
  @$pb.TagNumber(352209397)
  $core.bool hasJa4fingerprint() => $_has(7);
  @$pb.TagNumber(352209397)
  void clearJa4fingerprint() => $_clearField(352209397);
  @$pb.TagNumber(352209397)
  JA4Fingerprint ensureJa4fingerprint() => $_ensure(7);

  @$pb.TagNumber(413321041)
  Method get method => $_getN(8);
  @$pb.TagNumber(413321041)
  set method(Method value) => $_setField(413321041, value);
  @$pb.TagNumber(413321041)
  $core.bool hasMethod() => $_has(8);
  @$pb.TagNumber(413321041)
  void clearMethod() => $_clearField(413321041);
  @$pb.TagNumber(413321041)
  Method ensureMethod() => $_ensure(8);

  @$pb.TagNumber(418634949)
  Cookies get cookies => $_getN(9);
  @$pb.TagNumber(418634949)
  set cookies(Cookies value) => $_setField(418634949, value);
  @$pb.TagNumber(418634949)
  $core.bool hasCookies() => $_has(9);
  @$pb.TagNumber(418634949)
  void clearCookies() => $_clearField(418634949);
  @$pb.TagNumber(418634949)
  Cookies ensureCookies() => $_ensure(9);

  @$pb.TagNumber(419090489)
  HeaderOrder get headerorder => $_getN(10);
  @$pb.TagNumber(419090489)
  set headerorder(HeaderOrder value) => $_setField(419090489, value);
  @$pb.TagNumber(419090489)
  $core.bool hasHeaderorder() => $_has(10);
  @$pb.TagNumber(419090489)
  void clearHeaderorder() => $_clearField(419090489);
  @$pb.TagNumber(419090489)
  HeaderOrder ensureHeaderorder() => $_ensure(10);

  @$pb.TagNumber(434029149)
  SingleHeader get singleheader => $_getN(11);
  @$pb.TagNumber(434029149)
  set singleheader(SingleHeader value) => $_setField(434029149, value);
  @$pb.TagNumber(434029149)
  $core.bool hasSingleheader() => $_has(11);
  @$pb.TagNumber(434029149)
  void clearSingleheader() => $_clearField(434029149);
  @$pb.TagNumber(434029149)
  SingleHeader ensureSingleheader() => $_ensure(11);

  @$pb.TagNumber(435938663)
  QueryString get querystring => $_getN(12);
  @$pb.TagNumber(435938663)
  set querystring(QueryString value) => $_setField(435938663, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(12);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);
  @$pb.TagNumber(435938663)
  QueryString ensureQuerystring() => $_ensure(12);

  @$pb.TagNumber(521270332)
  JA3Fingerprint get ja3fingerprint => $_getN(13);
  @$pb.TagNumber(521270332)
  set ja3fingerprint(JA3Fingerprint value) => $_setField(521270332, value);
  @$pb.TagNumber(521270332)
  $core.bool hasJa3fingerprint() => $_has(13);
  @$pb.TagNumber(521270332)
  void clearJa3fingerprint() => $_clearField(521270332);
  @$pb.TagNumber(521270332)
  JA3Fingerprint ensureJa3fingerprint() => $_ensure(13);
}

class FieldToProtect extends $pb.GeneratedMessage {
  factory FieldToProtect({
    $core.Iterable<$core.String>? fieldkeys,
    FieldToProtectType? fieldtype,
  }) {
    final result = create();
    if (fieldkeys != null) result.fieldkeys.addAll(fieldkeys);
    if (fieldtype != null) result.fieldtype = fieldtype;
    return result;
  }

  FieldToProtect._();

  factory FieldToProtect.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FieldToProtect.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FieldToProtect',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(69138200, _omitFieldNames ? '' : 'fieldkeys')
    ..aE<FieldToProtectType>(363323648, _omitFieldNames ? '' : 'fieldtype',
        enumValues: FieldToProtectType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FieldToProtect clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FieldToProtect copyWith(void Function(FieldToProtect) updates) =>
      super.copyWith((message) => updates(message as FieldToProtect))
          as FieldToProtect;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FieldToProtect create() => FieldToProtect._();
  @$core.override
  FieldToProtect createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FieldToProtect getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FieldToProtect>(create);
  static FieldToProtect? _defaultInstance;

  @$pb.TagNumber(69138200)
  $pb.PbList<$core.String> get fieldkeys => $_getList(0);

  @$pb.TagNumber(363323648)
  FieldToProtectType get fieldtype => $_getN(1);
  @$pb.TagNumber(363323648)
  set fieldtype(FieldToProtectType value) => $_setField(363323648, value);
  @$pb.TagNumber(363323648)
  $core.bool hasFieldtype() => $_has(1);
  @$pb.TagNumber(363323648)
  void clearFieldtype() => $_clearField(363323648);
}

class Filter extends $pb.GeneratedMessage {
  factory Filter({
    FilterRequirement? requirement,
    $core.Iterable<Condition>? conditions,
    FilterBehavior? behavior,
  }) {
    final result = create();
    if (requirement != null) result.requirement = requirement;
    if (conditions != null) result.conditions.addAll(conditions);
    if (behavior != null) result.behavior = behavior;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<FilterRequirement>(213947081, _omitFieldNames ? '' : 'requirement',
        enumValues: FilterRequirement.values)
    ..pPM<Condition>(298559192, _omitFieldNames ? '' : 'conditions',
        subBuilder: Condition.create)
    ..aE<FilterBehavior>(521774152, _omitFieldNames ? '' : 'behavior',
        enumValues: FilterBehavior.values)
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

  @$pb.TagNumber(213947081)
  FilterRequirement get requirement => $_getN(0);
  @$pb.TagNumber(213947081)
  set requirement(FilterRequirement value) => $_setField(213947081, value);
  @$pb.TagNumber(213947081)
  $core.bool hasRequirement() => $_has(0);
  @$pb.TagNumber(213947081)
  void clearRequirement() => $_clearField(213947081);

  @$pb.TagNumber(298559192)
  $pb.PbList<Condition> get conditions => $_getList(1);

  @$pb.TagNumber(521774152)
  FilterBehavior get behavior => $_getN(2);
  @$pb.TagNumber(521774152)
  set behavior(FilterBehavior value) => $_setField(521774152, value);
  @$pb.TagNumber(521774152)
  $core.bool hasBehavior() => $_has(2);
  @$pb.TagNumber(521774152)
  void clearBehavior() => $_clearField(521774152);
}

class FilterSource extends $pb.GeneratedMessage {
  factory FilterSource({
    $core.String? botname,
    $core.String? botcategory,
    $core.String? botorganization,
  }) {
    final result = create();
    if (botname != null) result.botname = botname;
    if (botcategory != null) result.botcategory = botcategory;
    if (botorganization != null) result.botorganization = botorganization;
    return result;
  }

  FilterSource._();

  factory FilterSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FilterSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FilterSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(88172446, _omitFieldNames ? '' : 'botname')
    ..aOS(117673751, _omitFieldNames ? '' : 'botcategory')
    ..aOS(397082802, _omitFieldNames ? '' : 'botorganization')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilterSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FilterSource copyWith(void Function(FilterSource) updates) =>
      super.copyWith((message) => updates(message as FilterSource))
          as FilterSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FilterSource create() => FilterSource._();
  @$core.override
  FilterSource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FilterSource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FilterSource>(create);
  static FilterSource? _defaultInstance;

  @$pb.TagNumber(88172446)
  $core.String get botname => $_getSZ(0);
  @$pb.TagNumber(88172446)
  set botname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88172446)
  $core.bool hasBotname() => $_has(0);
  @$pb.TagNumber(88172446)
  void clearBotname() => $_clearField(88172446);

  @$pb.TagNumber(117673751)
  $core.String get botcategory => $_getSZ(1);
  @$pb.TagNumber(117673751)
  set botcategory($core.String value) => $_setString(1, value);
  @$pb.TagNumber(117673751)
  $core.bool hasBotcategory() => $_has(1);
  @$pb.TagNumber(117673751)
  void clearBotcategory() => $_clearField(117673751);

  @$pb.TagNumber(397082802)
  $core.String get botorganization => $_getSZ(2);
  @$pb.TagNumber(397082802)
  set botorganization($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397082802)
  $core.bool hasBotorganization() => $_has(2);
  @$pb.TagNumber(397082802)
  void clearBotorganization() => $_clearField(397082802);
}

class FirewallManagerRuleGroup extends $pb.GeneratedMessage {
  factory FirewallManagerRuleGroup({
    $core.int? priority,
    FirewallManagerStatement? firewallmanagerstatement,
    $core.String? name,
    VisibilityConfig? visibilityconfig,
    OverrideAction? overrideaction,
  }) {
    final result = create();
    if (priority != null) result.priority = priority;
    if (firewallmanagerstatement != null)
      result.firewallmanagerstatement = firewallmanagerstatement;
    if (name != null) result.name = name;
    if (visibilityconfig != null) result.visibilityconfig = visibilityconfig;
    if (overrideaction != null) result.overrideaction = overrideaction;
    return result;
  }

  FirewallManagerRuleGroup._();

  factory FirewallManagerRuleGroup.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FirewallManagerRuleGroup.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FirewallManagerRuleGroup',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aI(109944618, _omitFieldNames ? '' : 'priority')
    ..aOM<FirewallManagerStatement>(
        172974270, _omitFieldNames ? '' : 'firewallmanagerstatement',
        subBuilder: FirewallManagerStatement.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<VisibilityConfig>(
        507762728, _omitFieldNames ? '' : 'visibilityconfig',
        subBuilder: VisibilityConfig.create)
    ..aOM<OverrideAction>(515842888, _omitFieldNames ? '' : 'overrideaction',
        subBuilder: OverrideAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FirewallManagerRuleGroup clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FirewallManagerRuleGroup copyWith(
          void Function(FirewallManagerRuleGroup) updates) =>
      super.copyWith((message) => updates(message as FirewallManagerRuleGroup))
          as FirewallManagerRuleGroup;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FirewallManagerRuleGroup create() => FirewallManagerRuleGroup._();
  @$core.override
  FirewallManagerRuleGroup createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FirewallManagerRuleGroup getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FirewallManagerRuleGroup>(create);
  static FirewallManagerRuleGroup? _defaultInstance;

  @$pb.TagNumber(109944618)
  $core.int get priority => $_getIZ(0);
  @$pb.TagNumber(109944618)
  set priority($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(109944618)
  $core.bool hasPriority() => $_has(0);
  @$pb.TagNumber(109944618)
  void clearPriority() => $_clearField(109944618);

  @$pb.TagNumber(172974270)
  FirewallManagerStatement get firewallmanagerstatement => $_getN(1);
  @$pb.TagNumber(172974270)
  set firewallmanagerstatement(FirewallManagerStatement value) =>
      $_setField(172974270, value);
  @$pb.TagNumber(172974270)
  $core.bool hasFirewallmanagerstatement() => $_has(1);
  @$pb.TagNumber(172974270)
  void clearFirewallmanagerstatement() => $_clearField(172974270);
  @$pb.TagNumber(172974270)
  FirewallManagerStatement ensureFirewallmanagerstatement() => $_ensure(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(507762728)
  VisibilityConfig get visibilityconfig => $_getN(3);
  @$pb.TagNumber(507762728)
  set visibilityconfig(VisibilityConfig value) => $_setField(507762728, value);
  @$pb.TagNumber(507762728)
  $core.bool hasVisibilityconfig() => $_has(3);
  @$pb.TagNumber(507762728)
  void clearVisibilityconfig() => $_clearField(507762728);
  @$pb.TagNumber(507762728)
  VisibilityConfig ensureVisibilityconfig() => $_ensure(3);

  @$pb.TagNumber(515842888)
  OverrideAction get overrideaction => $_getN(4);
  @$pb.TagNumber(515842888)
  set overrideaction(OverrideAction value) => $_setField(515842888, value);
  @$pb.TagNumber(515842888)
  $core.bool hasOverrideaction() => $_has(4);
  @$pb.TagNumber(515842888)
  void clearOverrideaction() => $_clearField(515842888);
  @$pb.TagNumber(515842888)
  OverrideAction ensureOverrideaction() => $_ensure(4);
}

class FirewallManagerStatement extends $pb.GeneratedMessage {
  factory FirewallManagerStatement({
    RuleGroupReferenceStatement? rulegroupreferencestatement,
    ManagedRuleGroupStatement? managedrulegroupstatement,
  }) {
    final result = create();
    if (rulegroupreferencestatement != null)
      result.rulegroupreferencestatement = rulegroupreferencestatement;
    if (managedrulegroupstatement != null)
      result.managedrulegroupstatement = managedrulegroupstatement;
    return result;
  }

  FirewallManagerStatement._();

  factory FirewallManagerStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FirewallManagerStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FirewallManagerStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<RuleGroupReferenceStatement>(
        460139461, _omitFieldNames ? '' : 'rulegroupreferencestatement',
        subBuilder: RuleGroupReferenceStatement.create)
    ..aOM<ManagedRuleGroupStatement>(
        533120745, _omitFieldNames ? '' : 'managedrulegroupstatement',
        subBuilder: ManagedRuleGroupStatement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FirewallManagerStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FirewallManagerStatement copyWith(
          void Function(FirewallManagerStatement) updates) =>
      super.copyWith((message) => updates(message as FirewallManagerStatement))
          as FirewallManagerStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FirewallManagerStatement create() => FirewallManagerStatement._();
  @$core.override
  FirewallManagerStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FirewallManagerStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FirewallManagerStatement>(create);
  static FirewallManagerStatement? _defaultInstance;

  @$pb.TagNumber(460139461)
  RuleGroupReferenceStatement get rulegroupreferencestatement => $_getN(0);
  @$pb.TagNumber(460139461)
  set rulegroupreferencestatement(RuleGroupReferenceStatement value) =>
      $_setField(460139461, value);
  @$pb.TagNumber(460139461)
  $core.bool hasRulegroupreferencestatement() => $_has(0);
  @$pb.TagNumber(460139461)
  void clearRulegroupreferencestatement() => $_clearField(460139461);
  @$pb.TagNumber(460139461)
  RuleGroupReferenceStatement ensureRulegroupreferencestatement() =>
      $_ensure(0);

  @$pb.TagNumber(533120745)
  ManagedRuleGroupStatement get managedrulegroupstatement => $_getN(1);
  @$pb.TagNumber(533120745)
  set managedrulegroupstatement(ManagedRuleGroupStatement value) =>
      $_setField(533120745, value);
  @$pb.TagNumber(533120745)
  $core.bool hasManagedrulegroupstatement() => $_has(1);
  @$pb.TagNumber(533120745)
  void clearManagedrulegroupstatement() => $_clearField(533120745);
  @$pb.TagNumber(533120745)
  ManagedRuleGroupStatement ensureManagedrulegroupstatement() => $_ensure(1);
}

class ForwardedIPConfig extends $pb.GeneratedMessage {
  factory ForwardedIPConfig({
    $core.String? headername,
    FallbackBehavior? fallbackbehavior,
  }) {
    final result = create();
    if (headername != null) result.headername = headername;
    if (fallbackbehavior != null) result.fallbackbehavior = fallbackbehavior;
    return result;
  }

  ForwardedIPConfig._();

  factory ForwardedIPConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ForwardedIPConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ForwardedIPConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235244324, _omitFieldNames ? '' : 'headername')
    ..aE<FallbackBehavior>(440114542, _omitFieldNames ? '' : 'fallbackbehavior',
        enumValues: FallbackBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ForwardedIPConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ForwardedIPConfig copyWith(void Function(ForwardedIPConfig) updates) =>
      super.copyWith((message) => updates(message as ForwardedIPConfig))
          as ForwardedIPConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ForwardedIPConfig create() => ForwardedIPConfig._();
  @$core.override
  ForwardedIPConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ForwardedIPConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ForwardedIPConfig>(create);
  static ForwardedIPConfig? _defaultInstance;

  @$pb.TagNumber(235244324)
  $core.String get headername => $_getSZ(0);
  @$pb.TagNumber(235244324)
  set headername($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235244324)
  $core.bool hasHeadername() => $_has(0);
  @$pb.TagNumber(235244324)
  void clearHeadername() => $_clearField(235244324);

  @$pb.TagNumber(440114542)
  FallbackBehavior get fallbackbehavior => $_getN(1);
  @$pb.TagNumber(440114542)
  set fallbackbehavior(FallbackBehavior value) => $_setField(440114542, value);
  @$pb.TagNumber(440114542)
  $core.bool hasFallbackbehavior() => $_has(1);
  @$pb.TagNumber(440114542)
  void clearFallbackbehavior() => $_clearField(440114542);
}

class GenerateMobileSdkReleaseUrlRequest extends $pb.GeneratedMessage {
  factory GenerateMobileSdkReleaseUrlRequest({
    $core.String? releaseversion,
    Platform? platform,
  }) {
    final result = create();
    if (releaseversion != null) result.releaseversion = releaseversion;
    if (platform != null) result.platform = platform;
    return result;
  }

  GenerateMobileSdkReleaseUrlRequest._();

  factory GenerateMobileSdkReleaseUrlRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateMobileSdkReleaseUrlRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateMobileSdkReleaseUrlRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(8739171, _omitFieldNames ? '' : 'releaseversion')
    ..aE<Platform>(468905683, _omitFieldNames ? '' : 'platform',
        enumValues: Platform.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateMobileSdkReleaseUrlRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateMobileSdkReleaseUrlRequest copyWith(
          void Function(GenerateMobileSdkReleaseUrlRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GenerateMobileSdkReleaseUrlRequest))
          as GenerateMobileSdkReleaseUrlRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateMobileSdkReleaseUrlRequest create() =>
      GenerateMobileSdkReleaseUrlRequest._();
  @$core.override
  GenerateMobileSdkReleaseUrlRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateMobileSdkReleaseUrlRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GenerateMobileSdkReleaseUrlRequest>(
          create);
  static GenerateMobileSdkReleaseUrlRequest? _defaultInstance;

  @$pb.TagNumber(8739171)
  $core.String get releaseversion => $_getSZ(0);
  @$pb.TagNumber(8739171)
  set releaseversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(8739171)
  $core.bool hasReleaseversion() => $_has(0);
  @$pb.TagNumber(8739171)
  void clearReleaseversion() => $_clearField(8739171);

  @$pb.TagNumber(468905683)
  Platform get platform => $_getN(1);
  @$pb.TagNumber(468905683)
  set platform(Platform value) => $_setField(468905683, value);
  @$pb.TagNumber(468905683)
  $core.bool hasPlatform() => $_has(1);
  @$pb.TagNumber(468905683)
  void clearPlatform() => $_clearField(468905683);
}

class GenerateMobileSdkReleaseUrlResponse extends $pb.GeneratedMessage {
  factory GenerateMobileSdkReleaseUrlResponse({
    $core.String? url,
  }) {
    final result = create();
    if (url != null) result.url = url;
    return result;
  }

  GenerateMobileSdkReleaseUrlResponse._();

  factory GenerateMobileSdkReleaseUrlResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GenerateMobileSdkReleaseUrlResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GenerateMobileSdkReleaseUrlResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(354018239, _omitFieldNames ? '' : 'url')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateMobileSdkReleaseUrlResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GenerateMobileSdkReleaseUrlResponse copyWith(
          void Function(GenerateMobileSdkReleaseUrlResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GenerateMobileSdkReleaseUrlResponse))
          as GenerateMobileSdkReleaseUrlResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GenerateMobileSdkReleaseUrlResponse create() =>
      GenerateMobileSdkReleaseUrlResponse._();
  @$core.override
  GenerateMobileSdkReleaseUrlResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GenerateMobileSdkReleaseUrlResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GenerateMobileSdkReleaseUrlResponse>(create);
  static GenerateMobileSdkReleaseUrlResponse? _defaultInstance;

  @$pb.TagNumber(354018239)
  $core.String get url => $_getSZ(0);
  @$pb.TagNumber(354018239)
  set url($core.String value) => $_setString(0, value);
  @$pb.TagNumber(354018239)
  $core.bool hasUrl() => $_has(0);
  @$pb.TagNumber(354018239)
  void clearUrl() => $_clearField(354018239);
}

class GeoMatchStatement extends $pb.GeneratedMessage {
  factory GeoMatchStatement({
    ForwardedIPConfig? forwardedipconfig,
    $core.Iterable<CountryCode>? countrycodes,
  }) {
    final result = create();
    if (forwardedipconfig != null) result.forwardedipconfig = forwardedipconfig;
    if (countrycodes != null) result.countrycodes.addAll(countrycodes);
    return result;
  }

  GeoMatchStatement._();

  factory GeoMatchStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GeoMatchStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoMatchStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<ForwardedIPConfig>(
        259846797, _omitFieldNames ? '' : 'forwardedipconfig',
        subBuilder: ForwardedIPConfig.create)
    ..pc<CountryCode>(
        334257650, _omitFieldNames ? '' : 'countrycodes', $pb.PbFieldType.KE,
        valueOf: CountryCode.valueOf,
        enumValues: CountryCode.values,
        defaultEnumValue: CountryCode.COUNTRY_CODE_US)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchStatement copyWith(void Function(GeoMatchStatement) updates) =>
      super.copyWith((message) => updates(message as GeoMatchStatement))
          as GeoMatchStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoMatchStatement create() => GeoMatchStatement._();
  @$core.override
  GeoMatchStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GeoMatchStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GeoMatchStatement>(create);
  static GeoMatchStatement? _defaultInstance;

  @$pb.TagNumber(259846797)
  ForwardedIPConfig get forwardedipconfig => $_getN(0);
  @$pb.TagNumber(259846797)
  set forwardedipconfig(ForwardedIPConfig value) =>
      $_setField(259846797, value);
  @$pb.TagNumber(259846797)
  $core.bool hasForwardedipconfig() => $_has(0);
  @$pb.TagNumber(259846797)
  void clearForwardedipconfig() => $_clearField(259846797);
  @$pb.TagNumber(259846797)
  ForwardedIPConfig ensureForwardedipconfig() => $_ensure(0);

  @$pb.TagNumber(334257650)
  $pb.PbList<CountryCode> get countrycodes => $_getList(1);
}

class GetDecryptedAPIKeyRequest extends $pb.GeneratedMessage {
  factory GetDecryptedAPIKeyRequest({
    Scope? scope,
    $core.String? apikey,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (apikey != null) result.apikey = apikey;
    return result;
  }

  GetDecryptedAPIKeyRequest._();

  factory GetDecryptedAPIKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDecryptedAPIKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDecryptedAPIKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(274818239, _omitFieldNames ? '' : 'apikey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDecryptedAPIKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDecryptedAPIKeyRequest copyWith(
          void Function(GetDecryptedAPIKeyRequest) updates) =>
      super.copyWith((message) => updates(message as GetDecryptedAPIKeyRequest))
          as GetDecryptedAPIKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDecryptedAPIKeyRequest create() => GetDecryptedAPIKeyRequest._();
  @$core.override
  GetDecryptedAPIKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDecryptedAPIKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDecryptedAPIKeyRequest>(create);
  static GetDecryptedAPIKeyRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(274818239)
  $core.String get apikey => $_getSZ(1);
  @$pb.TagNumber(274818239)
  set apikey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(274818239)
  $core.bool hasApikey() => $_has(1);
  @$pb.TagNumber(274818239)
  void clearApikey() => $_clearField(274818239);
}

class GetDecryptedAPIKeyResponse extends $pb.GeneratedMessage {
  factory GetDecryptedAPIKeyResponse({
    $core.Iterable<$core.String>? tokendomains,
    $core.String? creationtimestamp,
  }) {
    final result = create();
    if (tokendomains != null) result.tokendomains.addAll(tokendomains);
    if (creationtimestamp != null) result.creationtimestamp = creationtimestamp;
    return result;
  }

  GetDecryptedAPIKeyResponse._();

  factory GetDecryptedAPIKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDecryptedAPIKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDecryptedAPIKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(11638560, _omitFieldNames ? '' : 'tokendomains')
    ..aOS(24480293, _omitFieldNames ? '' : 'creationtimestamp')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDecryptedAPIKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDecryptedAPIKeyResponse copyWith(
          void Function(GetDecryptedAPIKeyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetDecryptedAPIKeyResponse))
          as GetDecryptedAPIKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDecryptedAPIKeyResponse create() => GetDecryptedAPIKeyResponse._();
  @$core.override
  GetDecryptedAPIKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDecryptedAPIKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDecryptedAPIKeyResponse>(create);
  static GetDecryptedAPIKeyResponse? _defaultInstance;

  @$pb.TagNumber(11638560)
  $pb.PbList<$core.String> get tokendomains => $_getList(0);

  @$pb.TagNumber(24480293)
  $core.String get creationtimestamp => $_getSZ(1);
  @$pb.TagNumber(24480293)
  set creationtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(24480293)
  $core.bool hasCreationtimestamp() => $_has(1);
  @$pb.TagNumber(24480293)
  void clearCreationtimestamp() => $_clearField(24480293);
}

class GetIPSetRequest extends $pb.GeneratedMessage {
  factory GetIPSetRequest({
    Scope? scope,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  GetIPSetRequest._();

  factory GetIPSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetIPSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetIPSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIPSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIPSetRequest copyWith(void Function(GetIPSetRequest) updates) =>
      super.copyWith((message) => updates(message as GetIPSetRequest))
          as GetIPSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetIPSetRequest create() => GetIPSetRequest._();
  @$core.override
  GetIPSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetIPSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetIPSetRequest>(create);
  static GetIPSetRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class GetIPSetResponse extends $pb.GeneratedMessage {
  factory GetIPSetResponse({
    $core.String? locktoken,
    IPSet? ipset,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (ipset != null) result.ipset = ipset;
    return result;
  }

  GetIPSetResponse._();

  factory GetIPSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetIPSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetIPSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOM<IPSet>(436412565, _omitFieldNames ? '' : 'ipset',
        subBuilder: IPSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIPSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetIPSetResponse copyWith(void Function(GetIPSetResponse) updates) =>
      super.copyWith((message) => updates(message as GetIPSetResponse))
          as GetIPSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetIPSetResponse create() => GetIPSetResponse._();
  @$core.override
  GetIPSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetIPSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetIPSetResponse>(create);
  static GetIPSetResponse? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(436412565)
  IPSet get ipset => $_getN(1);
  @$pb.TagNumber(436412565)
  set ipset(IPSet value) => $_setField(436412565, value);
  @$pb.TagNumber(436412565)
  $core.bool hasIpset() => $_has(1);
  @$pb.TagNumber(436412565)
  void clearIpset() => $_clearField(436412565);
  @$pb.TagNumber(436412565)
  IPSet ensureIpset() => $_ensure(1);
}

class GetLoggingConfigurationRequest extends $pb.GeneratedMessage {
  factory GetLoggingConfigurationRequest({
    LogType? logtype,
    LogScope? logscope,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (logtype != null) result.logtype = logtype;
    if (logscope != null) result.logscope = logscope;
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  GetLoggingConfigurationRequest._();

  factory GetLoggingConfigurationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetLoggingConfigurationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetLoggingConfigurationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<LogType>(116703930, _omitFieldNames ? '' : 'logtype',
        enumValues: LogType.values)
    ..aE<LogScope>(188235840, _omitFieldNames ? '' : 'logscope',
        enumValues: LogScope.values)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetLoggingConfigurationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetLoggingConfigurationRequest copyWith(
          void Function(GetLoggingConfigurationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetLoggingConfigurationRequest))
          as GetLoggingConfigurationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetLoggingConfigurationRequest create() =>
      GetLoggingConfigurationRequest._();
  @$core.override
  GetLoggingConfigurationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetLoggingConfigurationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetLoggingConfigurationRequest>(create);
  static GetLoggingConfigurationRequest? _defaultInstance;

  @$pb.TagNumber(116703930)
  LogType get logtype => $_getN(0);
  @$pb.TagNumber(116703930)
  set logtype(LogType value) => $_setField(116703930, value);
  @$pb.TagNumber(116703930)
  $core.bool hasLogtype() => $_has(0);
  @$pb.TagNumber(116703930)
  void clearLogtype() => $_clearField(116703930);

  @$pb.TagNumber(188235840)
  LogScope get logscope => $_getN(1);
  @$pb.TagNumber(188235840)
  set logscope(LogScope value) => $_setField(188235840, value);
  @$pb.TagNumber(188235840)
  $core.bool hasLogscope() => $_has(1);
  @$pb.TagNumber(188235840)
  void clearLogscope() => $_clearField(188235840);

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(2);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(2);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class GetLoggingConfigurationResponse extends $pb.GeneratedMessage {
  factory GetLoggingConfigurationResponse({
    LoggingConfiguration? loggingconfiguration,
  }) {
    final result = create();
    if (loggingconfiguration != null)
      result.loggingconfiguration = loggingconfiguration;
    return result;
  }

  GetLoggingConfigurationResponse._();

  factory GetLoggingConfigurationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetLoggingConfigurationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetLoggingConfigurationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<LoggingConfiguration>(
        359027765, _omitFieldNames ? '' : 'loggingconfiguration',
        subBuilder: LoggingConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetLoggingConfigurationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetLoggingConfigurationResponse copyWith(
          void Function(GetLoggingConfigurationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetLoggingConfigurationResponse))
          as GetLoggingConfigurationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetLoggingConfigurationResponse create() =>
      GetLoggingConfigurationResponse._();
  @$core.override
  GetLoggingConfigurationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetLoggingConfigurationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetLoggingConfigurationResponse>(
          create);
  static GetLoggingConfigurationResponse? _defaultInstance;

  @$pb.TagNumber(359027765)
  LoggingConfiguration get loggingconfiguration => $_getN(0);
  @$pb.TagNumber(359027765)
  set loggingconfiguration(LoggingConfiguration value) =>
      $_setField(359027765, value);
  @$pb.TagNumber(359027765)
  $core.bool hasLoggingconfiguration() => $_has(0);
  @$pb.TagNumber(359027765)
  void clearLoggingconfiguration() => $_clearField(359027765);
  @$pb.TagNumber(359027765)
  LoggingConfiguration ensureLoggingconfiguration() => $_ensure(0);
}

class GetManagedRuleSetRequest extends $pb.GeneratedMessage {
  factory GetManagedRuleSetRequest({
    Scope? scope,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  GetManagedRuleSetRequest._();

  factory GetManagedRuleSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetManagedRuleSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetManagedRuleSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetManagedRuleSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetManagedRuleSetRequest copyWith(
          void Function(GetManagedRuleSetRequest) updates) =>
      super.copyWith((message) => updates(message as GetManagedRuleSetRequest))
          as GetManagedRuleSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetManagedRuleSetRequest create() => GetManagedRuleSetRequest._();
  @$core.override
  GetManagedRuleSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetManagedRuleSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetManagedRuleSetRequest>(create);
  static GetManagedRuleSetRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class GetManagedRuleSetResponse extends $pb.GeneratedMessage {
  factory GetManagedRuleSetResponse({
    $core.String? locktoken,
    ManagedRuleSet? managedruleset,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (managedruleset != null) result.managedruleset = managedruleset;
    return result;
  }

  GetManagedRuleSetResponse._();

  factory GetManagedRuleSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetManagedRuleSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetManagedRuleSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOM<ManagedRuleSet>(151429153, _omitFieldNames ? '' : 'managedruleset',
        subBuilder: ManagedRuleSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetManagedRuleSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetManagedRuleSetResponse copyWith(
          void Function(GetManagedRuleSetResponse) updates) =>
      super.copyWith((message) => updates(message as GetManagedRuleSetResponse))
          as GetManagedRuleSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetManagedRuleSetResponse create() => GetManagedRuleSetResponse._();
  @$core.override
  GetManagedRuleSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetManagedRuleSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetManagedRuleSetResponse>(create);
  static GetManagedRuleSetResponse? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(151429153)
  ManagedRuleSet get managedruleset => $_getN(1);
  @$pb.TagNumber(151429153)
  set managedruleset(ManagedRuleSet value) => $_setField(151429153, value);
  @$pb.TagNumber(151429153)
  $core.bool hasManagedruleset() => $_has(1);
  @$pb.TagNumber(151429153)
  void clearManagedruleset() => $_clearField(151429153);
  @$pb.TagNumber(151429153)
  ManagedRuleSet ensureManagedruleset() => $_ensure(1);
}

class GetMobileSdkReleaseRequest extends $pb.GeneratedMessage {
  factory GetMobileSdkReleaseRequest({
    $core.String? releaseversion,
    Platform? platform,
  }) {
    final result = create();
    if (releaseversion != null) result.releaseversion = releaseversion;
    if (platform != null) result.platform = platform;
    return result;
  }

  GetMobileSdkReleaseRequest._();

  factory GetMobileSdkReleaseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMobileSdkReleaseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMobileSdkReleaseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(8739171, _omitFieldNames ? '' : 'releaseversion')
    ..aE<Platform>(468905683, _omitFieldNames ? '' : 'platform',
        enumValues: Platform.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMobileSdkReleaseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMobileSdkReleaseRequest copyWith(
          void Function(GetMobileSdkReleaseRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetMobileSdkReleaseRequest))
          as GetMobileSdkReleaseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMobileSdkReleaseRequest create() => GetMobileSdkReleaseRequest._();
  @$core.override
  GetMobileSdkReleaseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMobileSdkReleaseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMobileSdkReleaseRequest>(create);
  static GetMobileSdkReleaseRequest? _defaultInstance;

  @$pb.TagNumber(8739171)
  $core.String get releaseversion => $_getSZ(0);
  @$pb.TagNumber(8739171)
  set releaseversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(8739171)
  $core.bool hasReleaseversion() => $_has(0);
  @$pb.TagNumber(8739171)
  void clearReleaseversion() => $_clearField(8739171);

  @$pb.TagNumber(468905683)
  Platform get platform => $_getN(1);
  @$pb.TagNumber(468905683)
  set platform(Platform value) => $_setField(468905683, value);
  @$pb.TagNumber(468905683)
  $core.bool hasPlatform() => $_has(1);
  @$pb.TagNumber(468905683)
  void clearPlatform() => $_clearField(468905683);
}

class GetMobileSdkReleaseResponse extends $pb.GeneratedMessage {
  factory GetMobileSdkReleaseResponse({
    MobileSdkRelease? mobilesdkrelease,
  }) {
    final result = create();
    if (mobilesdkrelease != null) result.mobilesdkrelease = mobilesdkrelease;
    return result;
  }

  GetMobileSdkReleaseResponse._();

  factory GetMobileSdkReleaseResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMobileSdkReleaseResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMobileSdkReleaseResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<MobileSdkRelease>(
        520859435, _omitFieldNames ? '' : 'mobilesdkrelease',
        subBuilder: MobileSdkRelease.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMobileSdkReleaseResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMobileSdkReleaseResponse copyWith(
          void Function(GetMobileSdkReleaseResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetMobileSdkReleaseResponse))
          as GetMobileSdkReleaseResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMobileSdkReleaseResponse create() =>
      GetMobileSdkReleaseResponse._();
  @$core.override
  GetMobileSdkReleaseResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMobileSdkReleaseResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMobileSdkReleaseResponse>(create);
  static GetMobileSdkReleaseResponse? _defaultInstance;

  @$pb.TagNumber(520859435)
  MobileSdkRelease get mobilesdkrelease => $_getN(0);
  @$pb.TagNumber(520859435)
  set mobilesdkrelease(MobileSdkRelease value) => $_setField(520859435, value);
  @$pb.TagNumber(520859435)
  $core.bool hasMobilesdkrelease() => $_has(0);
  @$pb.TagNumber(520859435)
  void clearMobilesdkrelease() => $_clearField(520859435);
  @$pb.TagNumber(520859435)
  MobileSdkRelease ensureMobilesdkrelease() => $_ensure(0);
}

class GetPermissionPolicyRequest extends $pb.GeneratedMessage {
  factory GetPermissionPolicyRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  GetPermissionPolicyRequest._();

  factory GetPermissionPolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPermissionPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPermissionPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPermissionPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPermissionPolicyRequest copyWith(
          void Function(GetPermissionPolicyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetPermissionPolicyRequest))
          as GetPermissionPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPermissionPolicyRequest create() => GetPermissionPolicyRequest._();
  @$core.override
  GetPermissionPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPermissionPolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPermissionPolicyRequest>(create);
  static GetPermissionPolicyRequest? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class GetPermissionPolicyResponse extends $pb.GeneratedMessage {
  factory GetPermissionPolicyResponse({
    $core.String? policy,
  }) {
    final result = create();
    if (policy != null) result.policy = policy;
    return result;
  }

  GetPermissionPolicyResponse._();

  factory GetPermissionPolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPermissionPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPermissionPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPermissionPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPermissionPolicyResponse copyWith(
          void Function(GetPermissionPolicyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetPermissionPolicyResponse))
          as GetPermissionPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPermissionPolicyResponse create() =>
      GetPermissionPolicyResponse._();
  @$core.override
  GetPermissionPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPermissionPolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPermissionPolicyResponse>(create);
  static GetPermissionPolicyResponse? _defaultInstance;

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(0);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(0);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class GetRateBasedStatementManagedKeysRequest extends $pb.GeneratedMessage {
  factory GetRateBasedStatementManagedKeysRequest({
    $core.String? rulegrouprulename,
    Scope? scope,
    $core.String? webaclid,
    $core.String? rulename,
    $core.String? webaclname,
  }) {
    final result = create();
    if (rulegrouprulename != null) result.rulegrouprulename = rulegrouprulename;
    if (scope != null) result.scope = scope;
    if (webaclid != null) result.webaclid = webaclid;
    if (rulename != null) result.rulename = rulename;
    if (webaclname != null) result.webaclname = webaclname;
    return result;
  }

  GetRateBasedStatementManagedKeysRequest._();

  factory GetRateBasedStatementManagedKeysRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRateBasedStatementManagedKeysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRateBasedStatementManagedKeysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(20574264, _omitFieldNames ? '' : 'rulegrouprulename')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(161274579, _omitFieldNames ? '' : 'webaclid')
    ..aOS(214688793, _omitFieldNames ? '' : 'rulename')
    ..aOS(261264029, _omitFieldNames ? '' : 'webaclname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedStatementManagedKeysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedStatementManagedKeysRequest copyWith(
          void Function(GetRateBasedStatementManagedKeysRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetRateBasedStatementManagedKeysRequest))
          as GetRateBasedStatementManagedKeysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRateBasedStatementManagedKeysRequest create() =>
      GetRateBasedStatementManagedKeysRequest._();
  @$core.override
  GetRateBasedStatementManagedKeysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRateBasedStatementManagedKeysRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetRateBasedStatementManagedKeysRequest>(create);
  static GetRateBasedStatementManagedKeysRequest? _defaultInstance;

  @$pb.TagNumber(20574264)
  $core.String get rulegrouprulename => $_getSZ(0);
  @$pb.TagNumber(20574264)
  set rulegrouprulename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20574264)
  $core.bool hasRulegrouprulename() => $_has(0);
  @$pb.TagNumber(20574264)
  void clearRulegrouprulename() => $_clearField(20574264);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(161274579)
  $core.String get webaclid => $_getSZ(2);
  @$pb.TagNumber(161274579)
  set webaclid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(161274579)
  $core.bool hasWebaclid() => $_has(2);
  @$pb.TagNumber(161274579)
  void clearWebaclid() => $_clearField(161274579);

  @$pb.TagNumber(214688793)
  $core.String get rulename => $_getSZ(3);
  @$pb.TagNumber(214688793)
  set rulename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(214688793)
  $core.bool hasRulename() => $_has(3);
  @$pb.TagNumber(214688793)
  void clearRulename() => $_clearField(214688793);

  @$pb.TagNumber(261264029)
  $core.String get webaclname => $_getSZ(4);
  @$pb.TagNumber(261264029)
  set webaclname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(261264029)
  $core.bool hasWebaclname() => $_has(4);
  @$pb.TagNumber(261264029)
  void clearWebaclname() => $_clearField(261264029);
}

class GetRateBasedStatementManagedKeysResponse extends $pb.GeneratedMessage {
  factory GetRateBasedStatementManagedKeysResponse({
    RateBasedStatementManagedKeysIPSet? managedkeysipv6,
    RateBasedStatementManagedKeysIPSet? managedkeysipv4,
  }) {
    final result = create();
    if (managedkeysipv6 != null) result.managedkeysipv6 = managedkeysipv6;
    if (managedkeysipv4 != null) result.managedkeysipv4 = managedkeysipv4;
    return result;
  }

  GetRateBasedStatementManagedKeysResponse._();

  factory GetRateBasedStatementManagedKeysResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRateBasedStatementManagedKeysResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRateBasedStatementManagedKeysResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<RateBasedStatementManagedKeysIPSet>(
        137610708, _omitFieldNames ? '' : 'managedkeysipv6',
        subBuilder: RateBasedStatementManagedKeysIPSet.create)
    ..aOM<RateBasedStatementManagedKeysIPSet>(
        171165946, _omitFieldNames ? '' : 'managedkeysipv4',
        subBuilder: RateBasedStatementManagedKeysIPSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedStatementManagedKeysResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedStatementManagedKeysResponse copyWith(
          void Function(GetRateBasedStatementManagedKeysResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetRateBasedStatementManagedKeysResponse))
          as GetRateBasedStatementManagedKeysResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRateBasedStatementManagedKeysResponse create() =>
      GetRateBasedStatementManagedKeysResponse._();
  @$core.override
  GetRateBasedStatementManagedKeysResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRateBasedStatementManagedKeysResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetRateBasedStatementManagedKeysResponse>(create);
  static GetRateBasedStatementManagedKeysResponse? _defaultInstance;

  @$pb.TagNumber(137610708)
  RateBasedStatementManagedKeysIPSet get managedkeysipv6 => $_getN(0);
  @$pb.TagNumber(137610708)
  set managedkeysipv6(RateBasedStatementManagedKeysIPSet value) =>
      $_setField(137610708, value);
  @$pb.TagNumber(137610708)
  $core.bool hasManagedkeysipv6() => $_has(0);
  @$pb.TagNumber(137610708)
  void clearManagedkeysipv6() => $_clearField(137610708);
  @$pb.TagNumber(137610708)
  RateBasedStatementManagedKeysIPSet ensureManagedkeysipv6() => $_ensure(0);

  @$pb.TagNumber(171165946)
  RateBasedStatementManagedKeysIPSet get managedkeysipv4 => $_getN(1);
  @$pb.TagNumber(171165946)
  set managedkeysipv4(RateBasedStatementManagedKeysIPSet value) =>
      $_setField(171165946, value);
  @$pb.TagNumber(171165946)
  $core.bool hasManagedkeysipv4() => $_has(1);
  @$pb.TagNumber(171165946)
  void clearManagedkeysipv4() => $_clearField(171165946);
  @$pb.TagNumber(171165946)
  RateBasedStatementManagedKeysIPSet ensureManagedkeysipv4() => $_ensure(1);
}

class GetRegexPatternSetRequest extends $pb.GeneratedMessage {
  factory GetRegexPatternSetRequest({
    Scope? scope,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  GetRegexPatternSetRequest._();

  factory GetRegexPatternSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRegexPatternSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRegexPatternSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRegexPatternSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRegexPatternSetRequest copyWith(
          void Function(GetRegexPatternSetRequest) updates) =>
      super.copyWith((message) => updates(message as GetRegexPatternSetRequest))
          as GetRegexPatternSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRegexPatternSetRequest create() => GetRegexPatternSetRequest._();
  @$core.override
  GetRegexPatternSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRegexPatternSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRegexPatternSetRequest>(create);
  static GetRegexPatternSetRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class GetRegexPatternSetResponse extends $pb.GeneratedMessage {
  factory GetRegexPatternSetResponse({
    RegexPatternSet? regexpatternset,
    $core.String? locktoken,
  }) {
    final result = create();
    if (regexpatternset != null) result.regexpatternset = regexpatternset;
    if (locktoken != null) result.locktoken = locktoken;
    return result;
  }

  GetRegexPatternSetResponse._();

  factory GetRegexPatternSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRegexPatternSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRegexPatternSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<RegexPatternSet>(9374915, _omitFieldNames ? '' : 'regexpatternset',
        subBuilder: RegexPatternSet.create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRegexPatternSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRegexPatternSetResponse copyWith(
          void Function(GetRegexPatternSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetRegexPatternSetResponse))
          as GetRegexPatternSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRegexPatternSetResponse create() => GetRegexPatternSetResponse._();
  @$core.override
  GetRegexPatternSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRegexPatternSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRegexPatternSetResponse>(create);
  static GetRegexPatternSetResponse? _defaultInstance;

  @$pb.TagNumber(9374915)
  RegexPatternSet get regexpatternset => $_getN(0);
  @$pb.TagNumber(9374915)
  set regexpatternset(RegexPatternSet value) => $_setField(9374915, value);
  @$pb.TagNumber(9374915)
  $core.bool hasRegexpatternset() => $_has(0);
  @$pb.TagNumber(9374915)
  void clearRegexpatternset() => $_clearField(9374915);
  @$pb.TagNumber(9374915)
  RegexPatternSet ensureRegexpatternset() => $_ensure(0);

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(1);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(1);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);
}

class GetRuleGroupRequest extends $pb.GeneratedMessage {
  factory GetRuleGroupRequest({
    Scope? scope,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  GetRuleGroupRequest._();

  factory GetRuleGroupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRuleGroupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRuleGroupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRuleGroupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRuleGroupRequest copyWith(void Function(GetRuleGroupRequest) updates) =>
      super.copyWith((message) => updates(message as GetRuleGroupRequest))
          as GetRuleGroupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRuleGroupRequest create() => GetRuleGroupRequest._();
  @$core.override
  GetRuleGroupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRuleGroupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRuleGroupRequest>(create);
  static GetRuleGroupRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(3);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(3);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class GetRuleGroupResponse extends $pb.GeneratedMessage {
  factory GetRuleGroupResponse({
    $core.String? locktoken,
    RuleGroup? rulegroup,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (rulegroup != null) result.rulegroup = rulegroup;
    return result;
  }

  GetRuleGroupResponse._();

  factory GetRuleGroupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRuleGroupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRuleGroupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOM<RuleGroup>(267398571, _omitFieldNames ? '' : 'rulegroup',
        subBuilder: RuleGroup.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRuleGroupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRuleGroupResponse copyWith(void Function(GetRuleGroupResponse) updates) =>
      super.copyWith((message) => updates(message as GetRuleGroupResponse))
          as GetRuleGroupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRuleGroupResponse create() => GetRuleGroupResponse._();
  @$core.override
  GetRuleGroupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRuleGroupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRuleGroupResponse>(create);
  static GetRuleGroupResponse? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(267398571)
  RuleGroup get rulegroup => $_getN(1);
  @$pb.TagNumber(267398571)
  set rulegroup(RuleGroup value) => $_setField(267398571, value);
  @$pb.TagNumber(267398571)
  $core.bool hasRulegroup() => $_has(1);
  @$pb.TagNumber(267398571)
  void clearRulegroup() => $_clearField(267398571);
  @$pb.TagNumber(267398571)
  RuleGroup ensureRulegroup() => $_ensure(1);
}

class GetSampledRequestsRequest extends $pb.GeneratedMessage {
  factory GetSampledRequestsRequest({
    Scope? scope,
    TimeWindow? timewindow,
    $core.String? rulemetricname,
    $fixnum.Int64? maxitems,
    $core.String? webaclarn,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (timewindow != null) result.timewindow = timewindow;
    if (rulemetricname != null) result.rulemetricname = rulemetricname;
    if (maxitems != null) result.maxitems = maxitems;
    if (webaclarn != null) result.webaclarn = webaclarn;
    return result;
  }

  GetSampledRequestsRequest._();

  factory GetSampledRequestsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSampledRequestsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSampledRequestsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOM<TimeWindow>(140543513, _omitFieldNames ? '' : 'timewindow',
        subBuilder: TimeWindow.create)
    ..aOS(141478429, _omitFieldNames ? '' : 'rulemetricname')
    ..aInt64(506899220, _omitFieldNames ? '' : 'maxitems')
    ..aOS(526110243, _omitFieldNames ? '' : 'webaclarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSampledRequestsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSampledRequestsRequest copyWith(
          void Function(GetSampledRequestsRequest) updates) =>
      super.copyWith((message) => updates(message as GetSampledRequestsRequest))
          as GetSampledRequestsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSampledRequestsRequest create() => GetSampledRequestsRequest._();
  @$core.override
  GetSampledRequestsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSampledRequestsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSampledRequestsRequest>(create);
  static GetSampledRequestsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(140543513)
  TimeWindow get timewindow => $_getN(1);
  @$pb.TagNumber(140543513)
  set timewindow(TimeWindow value) => $_setField(140543513, value);
  @$pb.TagNumber(140543513)
  $core.bool hasTimewindow() => $_has(1);
  @$pb.TagNumber(140543513)
  void clearTimewindow() => $_clearField(140543513);
  @$pb.TagNumber(140543513)
  TimeWindow ensureTimewindow() => $_ensure(1);

  @$pb.TagNumber(141478429)
  $core.String get rulemetricname => $_getSZ(2);
  @$pb.TagNumber(141478429)
  set rulemetricname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(141478429)
  $core.bool hasRulemetricname() => $_has(2);
  @$pb.TagNumber(141478429)
  void clearRulemetricname() => $_clearField(141478429);

  @$pb.TagNumber(506899220)
  $fixnum.Int64 get maxitems => $_getI64(3);
  @$pb.TagNumber(506899220)
  set maxitems($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);

  @$pb.TagNumber(526110243)
  $core.String get webaclarn => $_getSZ(4);
  @$pb.TagNumber(526110243)
  set webaclarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(526110243)
  $core.bool hasWebaclarn() => $_has(4);
  @$pb.TagNumber(526110243)
  void clearWebaclarn() => $_clearField(526110243);
}

class GetSampledRequestsResponse extends $pb.GeneratedMessage {
  factory GetSampledRequestsResponse({
    $core.Iterable<SampledHTTPRequest>? sampledrequests,
    TimeWindow? timewindow,
    $fixnum.Int64? populationsize,
  }) {
    final result = create();
    if (sampledrequests != null) result.sampledrequests.addAll(sampledrequests);
    if (timewindow != null) result.timewindow = timewindow;
    if (populationsize != null) result.populationsize = populationsize;
    return result;
  }

  GetSampledRequestsResponse._();

  factory GetSampledRequestsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSampledRequestsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSampledRequestsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<SampledHTTPRequest>(3233688, _omitFieldNames ? '' : 'sampledrequests',
        subBuilder: SampledHTTPRequest.create)
    ..aOM<TimeWindow>(140543513, _omitFieldNames ? '' : 'timewindow',
        subBuilder: TimeWindow.create)
    ..aInt64(140605436, _omitFieldNames ? '' : 'populationsize')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSampledRequestsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSampledRequestsResponse copyWith(
          void Function(GetSampledRequestsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetSampledRequestsResponse))
          as GetSampledRequestsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSampledRequestsResponse create() => GetSampledRequestsResponse._();
  @$core.override
  GetSampledRequestsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSampledRequestsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSampledRequestsResponse>(create);
  static GetSampledRequestsResponse? _defaultInstance;

  @$pb.TagNumber(3233688)
  $pb.PbList<SampledHTTPRequest> get sampledrequests => $_getList(0);

  @$pb.TagNumber(140543513)
  TimeWindow get timewindow => $_getN(1);
  @$pb.TagNumber(140543513)
  set timewindow(TimeWindow value) => $_setField(140543513, value);
  @$pb.TagNumber(140543513)
  $core.bool hasTimewindow() => $_has(1);
  @$pb.TagNumber(140543513)
  void clearTimewindow() => $_clearField(140543513);
  @$pb.TagNumber(140543513)
  TimeWindow ensureTimewindow() => $_ensure(1);

  @$pb.TagNumber(140605436)
  $fixnum.Int64 get populationsize => $_getI64(2);
  @$pb.TagNumber(140605436)
  set populationsize($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(140605436)
  $core.bool hasPopulationsize() => $_has(2);
  @$pb.TagNumber(140605436)
  void clearPopulationsize() => $_clearField(140605436);
}

class GetTopPathStatisticsByTrafficRequest extends $pb.GeneratedMessage {
  factory GetTopPathStatisticsByTrafficRequest({
    Scope? scope,
    $core.String? botname,
    $core.String? botcategory,
    TimeWindow? timewindow,
    $core.String? uripathprefix,
    $core.int? numberoftoptrafficbotsperpath,
    $core.String? botorganization,
    $core.int? limit,
    $core.String? webaclarn,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (botname != null) result.botname = botname;
    if (botcategory != null) result.botcategory = botcategory;
    if (timewindow != null) result.timewindow = timewindow;
    if (uripathprefix != null) result.uripathprefix = uripathprefix;
    if (numberoftoptrafficbotsperpath != null)
      result.numberoftoptrafficbotsperpath = numberoftoptrafficbotsperpath;
    if (botorganization != null) result.botorganization = botorganization;
    if (limit != null) result.limit = limit;
    if (webaclarn != null) result.webaclarn = webaclarn;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  GetTopPathStatisticsByTrafficRequest._();

  factory GetTopPathStatisticsByTrafficRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTopPathStatisticsByTrafficRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTopPathStatisticsByTrafficRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(88172446, _omitFieldNames ? '' : 'botname')
    ..aOS(117673751, _omitFieldNames ? '' : 'botcategory')
    ..aOM<TimeWindow>(140543513, _omitFieldNames ? '' : 'timewindow',
        subBuilder: TimeWindow.create)
    ..aOS(222177395, _omitFieldNames ? '' : 'uripathprefix')
    ..aI(302037370, _omitFieldNames ? '' : 'numberoftoptrafficbotsperpath')
    ..aOS(397082802, _omitFieldNames ? '' : 'botorganization')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(526110243, _omitFieldNames ? '' : 'webaclarn')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTopPathStatisticsByTrafficRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTopPathStatisticsByTrafficRequest copyWith(
          void Function(GetTopPathStatisticsByTrafficRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetTopPathStatisticsByTrafficRequest))
          as GetTopPathStatisticsByTrafficRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTopPathStatisticsByTrafficRequest create() =>
      GetTopPathStatisticsByTrafficRequest._();
  @$core.override
  GetTopPathStatisticsByTrafficRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTopPathStatisticsByTrafficRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetTopPathStatisticsByTrafficRequest>(create);
  static GetTopPathStatisticsByTrafficRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(88172446)
  $core.String get botname => $_getSZ(1);
  @$pb.TagNumber(88172446)
  set botname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(88172446)
  $core.bool hasBotname() => $_has(1);
  @$pb.TagNumber(88172446)
  void clearBotname() => $_clearField(88172446);

  @$pb.TagNumber(117673751)
  $core.String get botcategory => $_getSZ(2);
  @$pb.TagNumber(117673751)
  set botcategory($core.String value) => $_setString(2, value);
  @$pb.TagNumber(117673751)
  $core.bool hasBotcategory() => $_has(2);
  @$pb.TagNumber(117673751)
  void clearBotcategory() => $_clearField(117673751);

  @$pb.TagNumber(140543513)
  TimeWindow get timewindow => $_getN(3);
  @$pb.TagNumber(140543513)
  set timewindow(TimeWindow value) => $_setField(140543513, value);
  @$pb.TagNumber(140543513)
  $core.bool hasTimewindow() => $_has(3);
  @$pb.TagNumber(140543513)
  void clearTimewindow() => $_clearField(140543513);
  @$pb.TagNumber(140543513)
  TimeWindow ensureTimewindow() => $_ensure(3);

  @$pb.TagNumber(222177395)
  $core.String get uripathprefix => $_getSZ(4);
  @$pb.TagNumber(222177395)
  set uripathprefix($core.String value) => $_setString(4, value);
  @$pb.TagNumber(222177395)
  $core.bool hasUripathprefix() => $_has(4);
  @$pb.TagNumber(222177395)
  void clearUripathprefix() => $_clearField(222177395);

  @$pb.TagNumber(302037370)
  $core.int get numberoftoptrafficbotsperpath => $_getIZ(5);
  @$pb.TagNumber(302037370)
  set numberoftoptrafficbotsperpath($core.int value) =>
      $_setSignedInt32(5, value);
  @$pb.TagNumber(302037370)
  $core.bool hasNumberoftoptrafficbotsperpath() => $_has(5);
  @$pb.TagNumber(302037370)
  void clearNumberoftoptrafficbotsperpath() => $_clearField(302037370);

  @$pb.TagNumber(397082802)
  $core.String get botorganization => $_getSZ(6);
  @$pb.TagNumber(397082802)
  set botorganization($core.String value) => $_setString(6, value);
  @$pb.TagNumber(397082802)
  $core.bool hasBotorganization() => $_has(6);
  @$pb.TagNumber(397082802)
  void clearBotorganization() => $_clearField(397082802);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(7);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(7, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(7);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(526110243)
  $core.String get webaclarn => $_getSZ(8);
  @$pb.TagNumber(526110243)
  set webaclarn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(526110243)
  $core.bool hasWebaclarn() => $_has(8);
  @$pb.TagNumber(526110243)
  void clearWebaclarn() => $_clearField(526110243);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(9);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(9, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(9);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class GetTopPathStatisticsByTrafficResponse extends $pb.GeneratedMessage {
  factory GetTopPathStatisticsByTrafficResponse({
    $core.Iterable<PathStatistics>? pathstatistics,
    $core.Iterable<PathStatistics>? topcategories,
    $fixnum.Int64? totalrequestcount,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (pathstatistics != null) result.pathstatistics.addAll(pathstatistics);
    if (topcategories != null) result.topcategories.addAll(topcategories);
    if (totalrequestcount != null) result.totalrequestcount = totalrequestcount;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  GetTopPathStatisticsByTrafficResponse._();

  factory GetTopPathStatisticsByTrafficResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTopPathStatisticsByTrafficResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTopPathStatisticsByTrafficResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<PathStatistics>(322056386, _omitFieldNames ? '' : 'pathstatistics',
        subBuilder: PathStatistics.create)
    ..pPM<PathStatistics>(469322721, _omitFieldNames ? '' : 'topcategories',
        subBuilder: PathStatistics.create)
    ..aInt64(525689458, _omitFieldNames ? '' : 'totalrequestcount')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTopPathStatisticsByTrafficResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTopPathStatisticsByTrafficResponse copyWith(
          void Function(GetTopPathStatisticsByTrafficResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetTopPathStatisticsByTrafficResponse))
          as GetTopPathStatisticsByTrafficResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTopPathStatisticsByTrafficResponse create() =>
      GetTopPathStatisticsByTrafficResponse._();
  @$core.override
  GetTopPathStatisticsByTrafficResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTopPathStatisticsByTrafficResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetTopPathStatisticsByTrafficResponse>(create);
  static GetTopPathStatisticsByTrafficResponse? _defaultInstance;

  @$pb.TagNumber(322056386)
  $pb.PbList<PathStatistics> get pathstatistics => $_getList(0);

  @$pb.TagNumber(469322721)
  $pb.PbList<PathStatistics> get topcategories => $_getList(1);

  @$pb.TagNumber(525689458)
  $fixnum.Int64 get totalrequestcount => $_getI64(2);
  @$pb.TagNumber(525689458)
  set totalrequestcount($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(525689458)
  $core.bool hasTotalrequestcount() => $_has(2);
  @$pb.TagNumber(525689458)
  void clearTotalrequestcount() => $_clearField(525689458);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(3);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(3, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(3);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class GetWebACLForResourceRequest extends $pb.GeneratedMessage {
  factory GetWebACLForResourceRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  GetWebACLForResourceRequest._();

  factory GetWebACLForResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWebACLForResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWebACLForResourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebACLForResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebACLForResourceRequest copyWith(
          void Function(GetWebACLForResourceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetWebACLForResourceRequest))
          as GetWebACLForResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWebACLForResourceRequest create() =>
      GetWebACLForResourceRequest._();
  @$core.override
  GetWebACLForResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWebACLForResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWebACLForResourceRequest>(create);
  static GetWebACLForResourceRequest? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class GetWebACLForResourceResponse extends $pb.GeneratedMessage {
  factory GetWebACLForResourceResponse({
    WebACL? webacl,
  }) {
    final result = create();
    if (webacl != null) result.webacl = webacl;
    return result;
  }

  GetWebACLForResourceResponse._();

  factory GetWebACLForResourceResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWebACLForResourceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWebACLForResourceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<WebACL>(343373504, _omitFieldNames ? '' : 'webacl',
        subBuilder: WebACL.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebACLForResourceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebACLForResourceResponse copyWith(
          void Function(GetWebACLForResourceResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetWebACLForResourceResponse))
          as GetWebACLForResourceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWebACLForResourceResponse create() =>
      GetWebACLForResourceResponse._();
  @$core.override
  GetWebACLForResourceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWebACLForResourceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWebACLForResourceResponse>(create);
  static GetWebACLForResourceResponse? _defaultInstance;

  @$pb.TagNumber(343373504)
  WebACL get webacl => $_getN(0);
  @$pb.TagNumber(343373504)
  set webacl(WebACL value) => $_setField(343373504, value);
  @$pb.TagNumber(343373504)
  $core.bool hasWebacl() => $_has(0);
  @$pb.TagNumber(343373504)
  void clearWebacl() => $_clearField(343373504);
  @$pb.TagNumber(343373504)
  WebACL ensureWebacl() => $_ensure(0);
}

class GetWebACLRequest extends $pb.GeneratedMessage {
  factory GetWebACLRequest({
    Scope? scope,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  GetWebACLRequest._();

  factory GetWebACLRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWebACLRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWebACLRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebACLRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebACLRequest copyWith(void Function(GetWebACLRequest) updates) =>
      super.copyWith((message) => updates(message as GetWebACLRequest))
          as GetWebACLRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWebACLRequest create() => GetWebACLRequest._();
  @$core.override
  GetWebACLRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWebACLRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWebACLRequest>(create);
  static GetWebACLRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(3);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(3);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class GetWebACLResponse extends $pb.GeneratedMessage {
  factory GetWebACLResponse({
    $core.String? locktoken,
    $core.String? applicationintegrationurl,
    WebACL? webacl,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (applicationintegrationurl != null)
      result.applicationintegrationurl = applicationintegrationurl;
    if (webacl != null) result.webacl = webacl;
    return result;
  }

  GetWebACLResponse._();

  factory GetWebACLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWebACLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWebACLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOS(182743751, _omitFieldNames ? '' : 'applicationintegrationurl')
    ..aOM<WebACL>(343373504, _omitFieldNames ? '' : 'webacl',
        subBuilder: WebACL.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebACLResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWebACLResponse copyWith(void Function(GetWebACLResponse) updates) =>
      super.copyWith((message) => updates(message as GetWebACLResponse))
          as GetWebACLResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWebACLResponse create() => GetWebACLResponse._();
  @$core.override
  GetWebACLResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWebACLResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWebACLResponse>(create);
  static GetWebACLResponse? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(182743751)
  $core.String get applicationintegrationurl => $_getSZ(1);
  @$pb.TagNumber(182743751)
  set applicationintegrationurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(182743751)
  $core.bool hasApplicationintegrationurl() => $_has(1);
  @$pb.TagNumber(182743751)
  void clearApplicationintegrationurl() => $_clearField(182743751);

  @$pb.TagNumber(343373504)
  WebACL get webacl => $_getN(2);
  @$pb.TagNumber(343373504)
  set webacl(WebACL value) => $_setField(343373504, value);
  @$pb.TagNumber(343373504)
  $core.bool hasWebacl() => $_has(2);
  @$pb.TagNumber(343373504)
  void clearWebacl() => $_clearField(343373504);
  @$pb.TagNumber(343373504)
  WebACL ensureWebacl() => $_ensure(2);
}

class HTTPHeader extends $pb.GeneratedMessage {
  factory HTTPHeader({
    $core.String? name,
    $core.String? value,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (value != null) result.value = value;
    return result;
  }

  HTTPHeader._();

  factory HTTPHeader.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HTTPHeader.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HTTPHeader',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HTTPHeader clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HTTPHeader copyWith(void Function(HTTPHeader) updates) =>
      super.copyWith((message) => updates(message as HTTPHeader)) as HTTPHeader;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HTTPHeader create() => HTTPHeader._();
  @$core.override
  HTTPHeader createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HTTPHeader getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HTTPHeader>(create);
  static HTTPHeader? _defaultInstance;

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
}

class HTTPRequest extends $pb.GeneratedMessage {
  factory HTTPRequest({
    $core.String? httpversion,
    $core.String? country,
    $core.String? clientip,
    $core.Iterable<HTTPHeader>? headers,
    $core.String? method,
    $core.String? uri,
  }) {
    final result = create();
    if (httpversion != null) result.httpversion = httpversion;
    if (country != null) result.country = country;
    if (clientip != null) result.clientip = clientip;
    if (headers != null) result.headers.addAll(headers);
    if (method != null) result.method = method;
    if (uri != null) result.uri = uri;
    return result;
  }

  HTTPRequest._();

  factory HTTPRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HTTPRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HTTPRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(12814936, _omitFieldNames ? '' : 'httpversion')
    ..aOS(83164786, _omitFieldNames ? '' : 'country')
    ..aOS(247557856, _omitFieldNames ? '' : 'clientip')
    ..pPM<HTTPHeader>(323967370, _omitFieldNames ? '' : 'headers',
        subBuilder: HTTPHeader.create)
    ..aOS(413321041, _omitFieldNames ? '' : 'method')
    ..aOS(443116318, _omitFieldNames ? '' : 'uri')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HTTPRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HTTPRequest copyWith(void Function(HTTPRequest) updates) =>
      super.copyWith((message) => updates(message as HTTPRequest))
          as HTTPRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HTTPRequest create() => HTTPRequest._();
  @$core.override
  HTTPRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HTTPRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HTTPRequest>(create);
  static HTTPRequest? _defaultInstance;

  @$pb.TagNumber(12814936)
  $core.String get httpversion => $_getSZ(0);
  @$pb.TagNumber(12814936)
  set httpversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(12814936)
  $core.bool hasHttpversion() => $_has(0);
  @$pb.TagNumber(12814936)
  void clearHttpversion() => $_clearField(12814936);

  @$pb.TagNumber(83164786)
  $core.String get country => $_getSZ(1);
  @$pb.TagNumber(83164786)
  set country($core.String value) => $_setString(1, value);
  @$pb.TagNumber(83164786)
  $core.bool hasCountry() => $_has(1);
  @$pb.TagNumber(83164786)
  void clearCountry() => $_clearField(83164786);

  @$pb.TagNumber(247557856)
  $core.String get clientip => $_getSZ(2);
  @$pb.TagNumber(247557856)
  set clientip($core.String value) => $_setString(2, value);
  @$pb.TagNumber(247557856)
  $core.bool hasClientip() => $_has(2);
  @$pb.TagNumber(247557856)
  void clearClientip() => $_clearField(247557856);

  @$pb.TagNumber(323967370)
  $pb.PbList<HTTPHeader> get headers => $_getList(3);

  @$pb.TagNumber(413321041)
  $core.String get method => $_getSZ(4);
  @$pb.TagNumber(413321041)
  set method($core.String value) => $_setString(4, value);
  @$pb.TagNumber(413321041)
  $core.bool hasMethod() => $_has(4);
  @$pb.TagNumber(413321041)
  void clearMethod() => $_clearField(413321041);

  @$pb.TagNumber(443116318)
  $core.String get uri => $_getSZ(5);
  @$pb.TagNumber(443116318)
  set uri($core.String value) => $_setString(5, value);
  @$pb.TagNumber(443116318)
  $core.bool hasUri() => $_has(5);
  @$pb.TagNumber(443116318)
  void clearUri() => $_clearField(443116318);
}

class HeaderMatchPattern extends $pb.GeneratedMessage {
  factory HeaderMatchPattern({
    $core.Iterable<$core.String>? excludedheaders,
    $core.Iterable<$core.String>? includedheaders,
    All? all,
  }) {
    final result = create();
    if (excludedheaders != null) result.excludedheaders.addAll(excludedheaders);
    if (includedheaders != null) result.includedheaders.addAll(includedheaders);
    if (all != null) result.all = all;
    return result;
  }

  HeaderMatchPattern._();

  factory HeaderMatchPattern.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HeaderMatchPattern.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HeaderMatchPattern',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(30227480, _omitFieldNames ? '' : 'excludedheaders')
    ..pPS(97683566, _omitFieldNames ? '' : 'includedheaders')
    ..aOM<All>(363848549, _omitFieldNames ? '' : 'all', subBuilder: All.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HeaderMatchPattern clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HeaderMatchPattern copyWith(void Function(HeaderMatchPattern) updates) =>
      super.copyWith((message) => updates(message as HeaderMatchPattern))
          as HeaderMatchPattern;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HeaderMatchPattern create() => HeaderMatchPattern._();
  @$core.override
  HeaderMatchPattern createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HeaderMatchPattern getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HeaderMatchPattern>(create);
  static HeaderMatchPattern? _defaultInstance;

  @$pb.TagNumber(30227480)
  $pb.PbList<$core.String> get excludedheaders => $_getList(0);

  @$pb.TagNumber(97683566)
  $pb.PbList<$core.String> get includedheaders => $_getList(1);

  @$pb.TagNumber(363848549)
  All get all => $_getN(2);
  @$pb.TagNumber(363848549)
  set all(All value) => $_setField(363848549, value);
  @$pb.TagNumber(363848549)
  $core.bool hasAll() => $_has(2);
  @$pb.TagNumber(363848549)
  void clearAll() => $_clearField(363848549);
  @$pb.TagNumber(363848549)
  All ensureAll() => $_ensure(2);
}

class HeaderOrder extends $pb.GeneratedMessage {
  factory HeaderOrder({
    OversizeHandling? oversizehandling,
  }) {
    final result = create();
    if (oversizehandling != null) result.oversizehandling = oversizehandling;
    return result;
  }

  HeaderOrder._();

  factory HeaderOrder.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HeaderOrder.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HeaderOrder',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<OversizeHandling>(139375132, _omitFieldNames ? '' : 'oversizehandling',
        enumValues: OversizeHandling.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HeaderOrder clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HeaderOrder copyWith(void Function(HeaderOrder) updates) =>
      super.copyWith((message) => updates(message as HeaderOrder))
          as HeaderOrder;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HeaderOrder create() => HeaderOrder._();
  @$core.override
  HeaderOrder createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HeaderOrder getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HeaderOrder>(create);
  static HeaderOrder? _defaultInstance;

  @$pb.TagNumber(139375132)
  OversizeHandling get oversizehandling => $_getN(0);
  @$pb.TagNumber(139375132)
  set oversizehandling(OversizeHandling value) => $_setField(139375132, value);
  @$pb.TagNumber(139375132)
  $core.bool hasOversizehandling() => $_has(0);
  @$pb.TagNumber(139375132)
  void clearOversizehandling() => $_clearField(139375132);
}

class Headers extends $pb.GeneratedMessage {
  factory Headers({
    OversizeHandling? oversizehandling,
    MapMatchScope? matchscope,
    HeaderMatchPattern? matchpattern,
  }) {
    final result = create();
    if (oversizehandling != null) result.oversizehandling = oversizehandling;
    if (matchscope != null) result.matchscope = matchscope;
    if (matchpattern != null) result.matchpattern = matchpattern;
    return result;
  }

  Headers._();

  factory Headers.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Headers.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Headers',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<OversizeHandling>(139375132, _omitFieldNames ? '' : 'oversizehandling',
        enumValues: OversizeHandling.values)
    ..aE<MapMatchScope>(272445459, _omitFieldNames ? '' : 'matchscope',
        enumValues: MapMatchScope.values)
    ..aOM<HeaderMatchPattern>(294565637, _omitFieldNames ? '' : 'matchpattern',
        subBuilder: HeaderMatchPattern.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Headers clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Headers copyWith(void Function(Headers) updates) =>
      super.copyWith((message) => updates(message as Headers)) as Headers;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Headers create() => Headers._();
  @$core.override
  Headers createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Headers getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Headers>(create);
  static Headers? _defaultInstance;

  @$pb.TagNumber(139375132)
  OversizeHandling get oversizehandling => $_getN(0);
  @$pb.TagNumber(139375132)
  set oversizehandling(OversizeHandling value) => $_setField(139375132, value);
  @$pb.TagNumber(139375132)
  $core.bool hasOversizehandling() => $_has(0);
  @$pb.TagNumber(139375132)
  void clearOversizehandling() => $_clearField(139375132);

  @$pb.TagNumber(272445459)
  MapMatchScope get matchscope => $_getN(1);
  @$pb.TagNumber(272445459)
  set matchscope(MapMatchScope value) => $_setField(272445459, value);
  @$pb.TagNumber(272445459)
  $core.bool hasMatchscope() => $_has(1);
  @$pb.TagNumber(272445459)
  void clearMatchscope() => $_clearField(272445459);

  @$pb.TagNumber(294565637)
  HeaderMatchPattern get matchpattern => $_getN(2);
  @$pb.TagNumber(294565637)
  set matchpattern(HeaderMatchPattern value) => $_setField(294565637, value);
  @$pb.TagNumber(294565637)
  $core.bool hasMatchpattern() => $_has(2);
  @$pb.TagNumber(294565637)
  void clearMatchpattern() => $_clearField(294565637);
  @$pb.TagNumber(294565637)
  HeaderMatchPattern ensureMatchpattern() => $_ensure(2);
}

class IPSet extends $pb.GeneratedMessage {
  factory IPSet({
    $core.String? description,
    $core.String? name,
    IPAddressVersion? ipaddressversion,
    $core.Iterable<$core.String>? addresses,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (ipaddressversion != null) result.ipaddressversion = ipaddressversion;
    if (addresses != null) result.addresses.addAll(addresses);
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  IPSet._();

  factory IPSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IPSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IPSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<IPAddressVersion>(313363841, _omitFieldNames ? '' : 'ipaddressversion',
        enumValues: IPAddressVersion.values)
    ..pPS(375939972, _omitFieldNames ? '' : 'addresses')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSet copyWith(void Function(IPSet) updates) =>
      super.copyWith((message) => updates(message as IPSet)) as IPSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IPSet create() => IPSet._();
  @$core.override
  IPSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IPSet getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<IPSet>(create);
  static IPSet? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(313363841)
  IPAddressVersion get ipaddressversion => $_getN(2);
  @$pb.TagNumber(313363841)
  set ipaddressversion(IPAddressVersion value) => $_setField(313363841, value);
  @$pb.TagNumber(313363841)
  $core.bool hasIpaddressversion() => $_has(2);
  @$pb.TagNumber(313363841)
  void clearIpaddressversion() => $_clearField(313363841);

  @$pb.TagNumber(375939972)
  $pb.PbList<$core.String> get addresses => $_getList(3);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(5);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(5);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class IPSetForwardedIPConfig extends $pb.GeneratedMessage {
  factory IPSetForwardedIPConfig({
    ForwardedIPPosition? position,
    $core.String? headername,
    FallbackBehavior? fallbackbehavior,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (headername != null) result.headername = headername;
    if (fallbackbehavior != null) result.fallbackbehavior = fallbackbehavior;
    return result;
  }

  IPSetForwardedIPConfig._();

  factory IPSetForwardedIPConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IPSetForwardedIPConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IPSetForwardedIPConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<ForwardedIPPosition>(41890859, _omitFieldNames ? '' : 'position',
        enumValues: ForwardedIPPosition.values)
    ..aOS(235244324, _omitFieldNames ? '' : 'headername')
    ..aE<FallbackBehavior>(440114542, _omitFieldNames ? '' : 'fallbackbehavior',
        enumValues: FallbackBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetForwardedIPConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetForwardedIPConfig copyWith(
          void Function(IPSetForwardedIPConfig) updates) =>
      super.copyWith((message) => updates(message as IPSetForwardedIPConfig))
          as IPSetForwardedIPConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IPSetForwardedIPConfig create() => IPSetForwardedIPConfig._();
  @$core.override
  IPSetForwardedIPConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IPSetForwardedIPConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IPSetForwardedIPConfig>(create);
  static IPSetForwardedIPConfig? _defaultInstance;

  @$pb.TagNumber(41890859)
  ForwardedIPPosition get position => $_getN(0);
  @$pb.TagNumber(41890859)
  set position(ForwardedIPPosition value) => $_setField(41890859, value);
  @$pb.TagNumber(41890859)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(41890859)
  void clearPosition() => $_clearField(41890859);

  @$pb.TagNumber(235244324)
  $core.String get headername => $_getSZ(1);
  @$pb.TagNumber(235244324)
  set headername($core.String value) => $_setString(1, value);
  @$pb.TagNumber(235244324)
  $core.bool hasHeadername() => $_has(1);
  @$pb.TagNumber(235244324)
  void clearHeadername() => $_clearField(235244324);

  @$pb.TagNumber(440114542)
  FallbackBehavior get fallbackbehavior => $_getN(2);
  @$pb.TagNumber(440114542)
  set fallbackbehavior(FallbackBehavior value) => $_setField(440114542, value);
  @$pb.TagNumber(440114542)
  $core.bool hasFallbackbehavior() => $_has(2);
  @$pb.TagNumber(440114542)
  void clearFallbackbehavior() => $_clearField(440114542);
}

class IPSetReferenceStatement extends $pb.GeneratedMessage {
  factory IPSetReferenceStatement({
    IPSetForwardedIPConfig? ipsetforwardedipconfig,
    $core.String? arn,
  }) {
    final result = create();
    if (ipsetforwardedipconfig != null)
      result.ipsetforwardedipconfig = ipsetforwardedipconfig;
    if (arn != null) result.arn = arn;
    return result;
  }

  IPSetReferenceStatement._();

  factory IPSetReferenceStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IPSetReferenceStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IPSetReferenceStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<IPSetForwardedIPConfig>(
        245440244, _omitFieldNames ? '' : 'ipsetforwardedipconfig',
        subBuilder: IPSetForwardedIPConfig.create)
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetReferenceStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetReferenceStatement copyWith(
          void Function(IPSetReferenceStatement) updates) =>
      super.copyWith((message) => updates(message as IPSetReferenceStatement))
          as IPSetReferenceStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IPSetReferenceStatement create() => IPSetReferenceStatement._();
  @$core.override
  IPSetReferenceStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IPSetReferenceStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IPSetReferenceStatement>(create);
  static IPSetReferenceStatement? _defaultInstance;

  @$pb.TagNumber(245440244)
  IPSetForwardedIPConfig get ipsetforwardedipconfig => $_getN(0);
  @$pb.TagNumber(245440244)
  set ipsetforwardedipconfig(IPSetForwardedIPConfig value) =>
      $_setField(245440244, value);
  @$pb.TagNumber(245440244)
  $core.bool hasIpsetforwardedipconfig() => $_has(0);
  @$pb.TagNumber(245440244)
  void clearIpsetforwardedipconfig() => $_clearField(245440244);
  @$pb.TagNumber(245440244)
  IPSetForwardedIPConfig ensureIpsetforwardedipconfig() => $_ensure(0);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class IPSetSummary extends $pb.GeneratedMessage {
  factory IPSetSummary({
    $core.String? locktoken,
    $core.String? description,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  IPSetSummary._();

  factory IPSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IPSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IPSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetSummary copyWith(void Function(IPSetSummary) updates) =>
      super.copyWith((message) => updates(message as IPSetSummary))
          as IPSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IPSetSummary create() => IPSetSummary._();
  @$core.override
  IPSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IPSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IPSetSummary>(create);
  static IPSetSummary? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class ImmunityTimeProperty extends $pb.GeneratedMessage {
  factory ImmunityTimeProperty({
    $fixnum.Int64? immunitytime,
  }) {
    final result = create();
    if (immunitytime != null) result.immunitytime = immunitytime;
    return result;
  }

  ImmunityTimeProperty._();

  factory ImmunityTimeProperty.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImmunityTimeProperty.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImmunityTimeProperty',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aInt64(293097607, _omitFieldNames ? '' : 'immunitytime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImmunityTimeProperty clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImmunityTimeProperty copyWith(void Function(ImmunityTimeProperty) updates) =>
      super.copyWith((message) => updates(message as ImmunityTimeProperty))
          as ImmunityTimeProperty;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImmunityTimeProperty create() => ImmunityTimeProperty._();
  @$core.override
  ImmunityTimeProperty createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImmunityTimeProperty getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImmunityTimeProperty>(create);
  static ImmunityTimeProperty? _defaultInstance;

  @$pb.TagNumber(293097607)
  $fixnum.Int64 get immunitytime => $_getI64(0);
  @$pb.TagNumber(293097607)
  set immunitytime($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(293097607)
  $core.bool hasImmunitytime() => $_has(0);
  @$pb.TagNumber(293097607)
  void clearImmunitytime() => $_clearField(293097607);
}

class JA3Fingerprint extends $pb.GeneratedMessage {
  factory JA3Fingerprint({
    FallbackBehavior? fallbackbehavior,
  }) {
    final result = create();
    if (fallbackbehavior != null) result.fallbackbehavior = fallbackbehavior;
    return result;
  }

  JA3Fingerprint._();

  factory JA3Fingerprint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory JA3Fingerprint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'JA3Fingerprint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<FallbackBehavior>(440114542, _omitFieldNames ? '' : 'fallbackbehavior',
        enumValues: FallbackBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JA3Fingerprint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JA3Fingerprint copyWith(void Function(JA3Fingerprint) updates) =>
      super.copyWith((message) => updates(message as JA3Fingerprint))
          as JA3Fingerprint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static JA3Fingerprint create() => JA3Fingerprint._();
  @$core.override
  JA3Fingerprint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static JA3Fingerprint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<JA3Fingerprint>(create);
  static JA3Fingerprint? _defaultInstance;

  @$pb.TagNumber(440114542)
  FallbackBehavior get fallbackbehavior => $_getN(0);
  @$pb.TagNumber(440114542)
  set fallbackbehavior(FallbackBehavior value) => $_setField(440114542, value);
  @$pb.TagNumber(440114542)
  $core.bool hasFallbackbehavior() => $_has(0);
  @$pb.TagNumber(440114542)
  void clearFallbackbehavior() => $_clearField(440114542);
}

class JA4Fingerprint extends $pb.GeneratedMessage {
  factory JA4Fingerprint({
    FallbackBehavior? fallbackbehavior,
  }) {
    final result = create();
    if (fallbackbehavior != null) result.fallbackbehavior = fallbackbehavior;
    return result;
  }

  JA4Fingerprint._();

  factory JA4Fingerprint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory JA4Fingerprint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'JA4Fingerprint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<FallbackBehavior>(440114542, _omitFieldNames ? '' : 'fallbackbehavior',
        enumValues: FallbackBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JA4Fingerprint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JA4Fingerprint copyWith(void Function(JA4Fingerprint) updates) =>
      super.copyWith((message) => updates(message as JA4Fingerprint))
          as JA4Fingerprint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static JA4Fingerprint create() => JA4Fingerprint._();
  @$core.override
  JA4Fingerprint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static JA4Fingerprint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<JA4Fingerprint>(create);
  static JA4Fingerprint? _defaultInstance;

  @$pb.TagNumber(440114542)
  FallbackBehavior get fallbackbehavior => $_getN(0);
  @$pb.TagNumber(440114542)
  set fallbackbehavior(FallbackBehavior value) => $_setField(440114542, value);
  @$pb.TagNumber(440114542)
  $core.bool hasFallbackbehavior() => $_has(0);
  @$pb.TagNumber(440114542)
  void clearFallbackbehavior() => $_clearField(440114542);
}

class JsonBody extends $pb.GeneratedMessage {
  factory JsonBody({
    OversizeHandling? oversizehandling,
    JsonMatchScope? matchscope,
    JsonMatchPattern? matchpattern,
    BodyParsingFallbackBehavior? invalidfallbackbehavior,
  }) {
    final result = create();
    if (oversizehandling != null) result.oversizehandling = oversizehandling;
    if (matchscope != null) result.matchscope = matchscope;
    if (matchpattern != null) result.matchpattern = matchpattern;
    if (invalidfallbackbehavior != null)
      result.invalidfallbackbehavior = invalidfallbackbehavior;
    return result;
  }

  JsonBody._();

  factory JsonBody.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory JsonBody.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'JsonBody',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<OversizeHandling>(139375132, _omitFieldNames ? '' : 'oversizehandling',
        enumValues: OversizeHandling.values)
    ..aE<JsonMatchScope>(272445459, _omitFieldNames ? '' : 'matchscope',
        enumValues: JsonMatchScope.values)
    ..aOM<JsonMatchPattern>(294565637, _omitFieldNames ? '' : 'matchpattern',
        subBuilder: JsonMatchPattern.create)
    ..aE<BodyParsingFallbackBehavior>(
        302950641, _omitFieldNames ? '' : 'invalidfallbackbehavior',
        enumValues: BodyParsingFallbackBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JsonBody clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JsonBody copyWith(void Function(JsonBody) updates) =>
      super.copyWith((message) => updates(message as JsonBody)) as JsonBody;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static JsonBody create() => JsonBody._();
  @$core.override
  JsonBody createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static JsonBody getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<JsonBody>(create);
  static JsonBody? _defaultInstance;

  @$pb.TagNumber(139375132)
  OversizeHandling get oversizehandling => $_getN(0);
  @$pb.TagNumber(139375132)
  set oversizehandling(OversizeHandling value) => $_setField(139375132, value);
  @$pb.TagNumber(139375132)
  $core.bool hasOversizehandling() => $_has(0);
  @$pb.TagNumber(139375132)
  void clearOversizehandling() => $_clearField(139375132);

  @$pb.TagNumber(272445459)
  JsonMatchScope get matchscope => $_getN(1);
  @$pb.TagNumber(272445459)
  set matchscope(JsonMatchScope value) => $_setField(272445459, value);
  @$pb.TagNumber(272445459)
  $core.bool hasMatchscope() => $_has(1);
  @$pb.TagNumber(272445459)
  void clearMatchscope() => $_clearField(272445459);

  @$pb.TagNumber(294565637)
  JsonMatchPattern get matchpattern => $_getN(2);
  @$pb.TagNumber(294565637)
  set matchpattern(JsonMatchPattern value) => $_setField(294565637, value);
  @$pb.TagNumber(294565637)
  $core.bool hasMatchpattern() => $_has(2);
  @$pb.TagNumber(294565637)
  void clearMatchpattern() => $_clearField(294565637);
  @$pb.TagNumber(294565637)
  JsonMatchPattern ensureMatchpattern() => $_ensure(2);

  @$pb.TagNumber(302950641)
  BodyParsingFallbackBehavior get invalidfallbackbehavior => $_getN(3);
  @$pb.TagNumber(302950641)
  set invalidfallbackbehavior(BodyParsingFallbackBehavior value) =>
      $_setField(302950641, value);
  @$pb.TagNumber(302950641)
  $core.bool hasInvalidfallbackbehavior() => $_has(3);
  @$pb.TagNumber(302950641)
  void clearInvalidfallbackbehavior() => $_clearField(302950641);
}

class JsonMatchPattern extends $pb.GeneratedMessage {
  factory JsonMatchPattern({
    All? all,
    $core.Iterable<$core.String>? includedpaths,
  }) {
    final result = create();
    if (all != null) result.all = all;
    if (includedpaths != null) result.includedpaths.addAll(includedpaths);
    return result;
  }

  JsonMatchPattern._();

  factory JsonMatchPattern.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory JsonMatchPattern.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'JsonMatchPattern',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<All>(363848549, _omitFieldNames ? '' : 'all', subBuilder: All.create)
    ..pPS(474316228, _omitFieldNames ? '' : 'includedpaths')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JsonMatchPattern clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  JsonMatchPattern copyWith(void Function(JsonMatchPattern) updates) =>
      super.copyWith((message) => updates(message as JsonMatchPattern))
          as JsonMatchPattern;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static JsonMatchPattern create() => JsonMatchPattern._();
  @$core.override
  JsonMatchPattern createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static JsonMatchPattern getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<JsonMatchPattern>(create);
  static JsonMatchPattern? _defaultInstance;

  @$pb.TagNumber(363848549)
  All get all => $_getN(0);
  @$pb.TagNumber(363848549)
  set all(All value) => $_setField(363848549, value);
  @$pb.TagNumber(363848549)
  $core.bool hasAll() => $_has(0);
  @$pb.TagNumber(363848549)
  void clearAll() => $_clearField(363848549);
  @$pb.TagNumber(363848549)
  All ensureAll() => $_ensure(0);

  @$pb.TagNumber(474316228)
  $pb.PbList<$core.String> get includedpaths => $_getList(1);
}

class Label extends $pb.GeneratedMessage {
  factory Label({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  Label._();

  factory Label.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Label.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Label',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Label clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Label copyWith(void Function(Label) updates) =>
      super.copyWith((message) => updates(message as Label)) as Label;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Label create() => Label._();
  @$core.override
  Label createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Label getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Label>(create);
  static Label? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class LabelMatchStatement extends $pb.GeneratedMessage {
  factory LabelMatchStatement({
    LabelMatchScope? scope,
    $core.String? key,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (key != null) result.key = key;
    return result;
  }

  LabelMatchStatement._();

  factory LabelMatchStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LabelMatchStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LabelMatchStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<LabelMatchScope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: LabelMatchScope.values)
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LabelMatchStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LabelMatchStatement copyWith(void Function(LabelMatchStatement) updates) =>
      super.copyWith((message) => updates(message as LabelMatchStatement))
          as LabelMatchStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LabelMatchStatement create() => LabelMatchStatement._();
  @$core.override
  LabelMatchStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LabelMatchStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LabelMatchStatement>(create);
  static LabelMatchStatement? _defaultInstance;

  @$pb.TagNumber(65430924)
  LabelMatchScope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(LabelMatchScope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(219859213)
  $core.String get key => $_getSZ(1);
  @$pb.TagNumber(219859213)
  set key($core.String value) => $_setString(1, value);
  @$pb.TagNumber(219859213)
  $core.bool hasKey() => $_has(1);
  @$pb.TagNumber(219859213)
  void clearKey() => $_clearField(219859213);
}

class LabelNameCondition extends $pb.GeneratedMessage {
  factory LabelNameCondition({
    $core.String? labelname,
  }) {
    final result = create();
    if (labelname != null) result.labelname = labelname;
    return result;
  }

  LabelNameCondition._();

  factory LabelNameCondition.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LabelNameCondition.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LabelNameCondition',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(238708671, _omitFieldNames ? '' : 'labelname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LabelNameCondition clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LabelNameCondition copyWith(void Function(LabelNameCondition) updates) =>
      super.copyWith((message) => updates(message as LabelNameCondition))
          as LabelNameCondition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LabelNameCondition create() => LabelNameCondition._();
  @$core.override
  LabelNameCondition createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LabelNameCondition getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LabelNameCondition>(create);
  static LabelNameCondition? _defaultInstance;

  @$pb.TagNumber(238708671)
  $core.String get labelname => $_getSZ(0);
  @$pb.TagNumber(238708671)
  set labelname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(238708671)
  $core.bool hasLabelname() => $_has(0);
  @$pb.TagNumber(238708671)
  void clearLabelname() => $_clearField(238708671);
}

class LabelSummary extends $pb.GeneratedMessage {
  factory LabelSummary({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  LabelSummary._();

  factory LabelSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LabelSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LabelSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LabelSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LabelSummary copyWith(void Function(LabelSummary) updates) =>
      super.copyWith((message) => updates(message as LabelSummary))
          as LabelSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LabelSummary create() => LabelSummary._();
  @$core.override
  LabelSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LabelSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LabelSummary>(create);
  static LabelSummary? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ListAPIKeysRequest extends $pb.GeneratedMessage {
  factory ListAPIKeysRequest({
    Scope? scope,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListAPIKeysRequest._();

  factory ListAPIKeysRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAPIKeysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAPIKeysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAPIKeysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAPIKeysRequest copyWith(void Function(ListAPIKeysRequest) updates) =>
      super.copyWith((message) => updates(message as ListAPIKeysRequest))
          as ListAPIKeysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAPIKeysRequest create() => ListAPIKeysRequest._();
  @$core.override
  ListAPIKeysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAPIKeysRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListAPIKeysRequest>(create);
  static ListAPIKeysRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListAPIKeysResponse extends $pb.GeneratedMessage {
  factory ListAPIKeysResponse({
    $core.Iterable<APIKeySummary>? apikeysummaries,
    $core.String? applicationintegrationurl,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (apikeysummaries != null) result.apikeysummaries.addAll(apikeysummaries);
    if (applicationintegrationurl != null)
      result.applicationintegrationurl = applicationintegrationurl;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListAPIKeysResponse._();

  factory ListAPIKeysResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAPIKeysResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAPIKeysResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<APIKeySummary>(137828491, _omitFieldNames ? '' : 'apikeysummaries',
        subBuilder: APIKeySummary.create)
    ..aOS(182743751, _omitFieldNames ? '' : 'applicationintegrationurl')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAPIKeysResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAPIKeysResponse copyWith(void Function(ListAPIKeysResponse) updates) =>
      super.copyWith((message) => updates(message as ListAPIKeysResponse))
          as ListAPIKeysResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAPIKeysResponse create() => ListAPIKeysResponse._();
  @$core.override
  ListAPIKeysResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAPIKeysResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListAPIKeysResponse>(create);
  static ListAPIKeysResponse? _defaultInstance;

  @$pb.TagNumber(137828491)
  $pb.PbList<APIKeySummary> get apikeysummaries => $_getList(0);

  @$pb.TagNumber(182743751)
  $core.String get applicationintegrationurl => $_getSZ(1);
  @$pb.TagNumber(182743751)
  set applicationintegrationurl($core.String value) => $_setString(1, value);
  @$pb.TagNumber(182743751)
  $core.bool hasApplicationintegrationurl() => $_has(1);
  @$pb.TagNumber(182743751)
  void clearApplicationintegrationurl() => $_clearField(182743751);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListAvailableManagedRuleGroupVersionsRequest
    extends $pb.GeneratedMessage {
  factory ListAvailableManagedRuleGroupVersionsRequest({
    Scope? scope,
    $core.String? vendorname,
    $core.String? name,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (vendorname != null) result.vendorname = vendorname;
    if (name != null) result.name = name;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListAvailableManagedRuleGroupVersionsRequest._();

  factory ListAvailableManagedRuleGroupVersionsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAvailableManagedRuleGroupVersionsRequest.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAvailableManagedRuleGroupVersionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(159329689, _omitFieldNames ? '' : 'vendorname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAvailableManagedRuleGroupVersionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAvailableManagedRuleGroupVersionsRequest copyWith(
          void Function(ListAvailableManagedRuleGroupVersionsRequest)
              updates) =>
      super.copyWith((message) =>
              updates(message as ListAvailableManagedRuleGroupVersionsRequest))
          as ListAvailableManagedRuleGroupVersionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAvailableManagedRuleGroupVersionsRequest create() =>
      ListAvailableManagedRuleGroupVersionsRequest._();
  @$core.override
  ListAvailableManagedRuleGroupVersionsRequest createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static ListAvailableManagedRuleGroupVersionsRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListAvailableManagedRuleGroupVersionsRequest>(create);
  static ListAvailableManagedRuleGroupVersionsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(159329689)
  $core.String get vendorname => $_getSZ(1);
  @$pb.TagNumber(159329689)
  set vendorname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(159329689)
  $core.bool hasVendorname() => $_has(1);
  @$pb.TagNumber(159329689)
  void clearVendorname() => $_clearField(159329689);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(4);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(4, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(4);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListAvailableManagedRuleGroupVersionsResponse
    extends $pb.GeneratedMessage {
  factory ListAvailableManagedRuleGroupVersionsResponse({
    $core.Iterable<ManagedRuleGroupVersion>? versions,
    $core.String? currentdefaultversion,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (versions != null) result.versions.addAll(versions);
    if (currentdefaultversion != null)
      result.currentdefaultversion = currentdefaultversion;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListAvailableManagedRuleGroupVersionsResponse._();

  factory ListAvailableManagedRuleGroupVersionsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAvailableManagedRuleGroupVersionsResponse.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAvailableManagedRuleGroupVersionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<ManagedRuleGroupVersion>(252099085, _omitFieldNames ? '' : 'versions',
        subBuilder: ManagedRuleGroupVersion.create)
    ..aOS(329624010, _omitFieldNames ? '' : 'currentdefaultversion')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAvailableManagedRuleGroupVersionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAvailableManagedRuleGroupVersionsResponse copyWith(
          void Function(ListAvailableManagedRuleGroupVersionsResponse)
              updates) =>
      super.copyWith((message) =>
              updates(message as ListAvailableManagedRuleGroupVersionsResponse))
          as ListAvailableManagedRuleGroupVersionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAvailableManagedRuleGroupVersionsResponse create() =>
      ListAvailableManagedRuleGroupVersionsResponse._();
  @$core.override
  ListAvailableManagedRuleGroupVersionsResponse createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static ListAvailableManagedRuleGroupVersionsResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListAvailableManagedRuleGroupVersionsResponse>(create);
  static ListAvailableManagedRuleGroupVersionsResponse? _defaultInstance;

  @$pb.TagNumber(252099085)
  $pb.PbList<ManagedRuleGroupVersion> get versions => $_getList(0);

  @$pb.TagNumber(329624010)
  $core.String get currentdefaultversion => $_getSZ(1);
  @$pb.TagNumber(329624010)
  set currentdefaultversion($core.String value) => $_setString(1, value);
  @$pb.TagNumber(329624010)
  $core.bool hasCurrentdefaultversion() => $_has(1);
  @$pb.TagNumber(329624010)
  void clearCurrentdefaultversion() => $_clearField(329624010);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListAvailableManagedRuleGroupsRequest extends $pb.GeneratedMessage {
  factory ListAvailableManagedRuleGroupsRequest({
    Scope? scope,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListAvailableManagedRuleGroupsRequest._();

  factory ListAvailableManagedRuleGroupsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAvailableManagedRuleGroupsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAvailableManagedRuleGroupsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAvailableManagedRuleGroupsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAvailableManagedRuleGroupsRequest copyWith(
          void Function(ListAvailableManagedRuleGroupsRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ListAvailableManagedRuleGroupsRequest))
          as ListAvailableManagedRuleGroupsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAvailableManagedRuleGroupsRequest create() =>
      ListAvailableManagedRuleGroupsRequest._();
  @$core.override
  ListAvailableManagedRuleGroupsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAvailableManagedRuleGroupsRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListAvailableManagedRuleGroupsRequest>(create);
  static ListAvailableManagedRuleGroupsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListAvailableManagedRuleGroupsResponse extends $pb.GeneratedMessage {
  factory ListAvailableManagedRuleGroupsResponse({
    $core.Iterable<ManagedRuleGroupSummary>? managedrulegroups,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (managedrulegroups != null)
      result.managedrulegroups.addAll(managedrulegroups);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListAvailableManagedRuleGroupsResponse._();

  factory ListAvailableManagedRuleGroupsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAvailableManagedRuleGroupsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAvailableManagedRuleGroupsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<ManagedRuleGroupSummary>(
        46574561, _omitFieldNames ? '' : 'managedrulegroups',
        subBuilder: ManagedRuleGroupSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAvailableManagedRuleGroupsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAvailableManagedRuleGroupsResponse copyWith(
          void Function(ListAvailableManagedRuleGroupsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListAvailableManagedRuleGroupsResponse))
          as ListAvailableManagedRuleGroupsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAvailableManagedRuleGroupsResponse create() =>
      ListAvailableManagedRuleGroupsResponse._();
  @$core.override
  ListAvailableManagedRuleGroupsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAvailableManagedRuleGroupsResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListAvailableManagedRuleGroupsResponse>(create);
  static ListAvailableManagedRuleGroupsResponse? _defaultInstance;

  @$pb.TagNumber(46574561)
  $pb.PbList<ManagedRuleGroupSummary> get managedrulegroups => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListIPSetsRequest extends $pb.GeneratedMessage {
  factory ListIPSetsRequest({
    Scope? scope,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListIPSetsRequest._();

  factory ListIPSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListIPSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListIPSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIPSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIPSetsRequest copyWith(void Function(ListIPSetsRequest) updates) =>
      super.copyWith((message) => updates(message as ListIPSetsRequest))
          as ListIPSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListIPSetsRequest create() => ListIPSetsRequest._();
  @$core.override
  ListIPSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListIPSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListIPSetsRequest>(create);
  static ListIPSetsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListIPSetsResponse extends $pb.GeneratedMessage {
  factory ListIPSetsResponse({
    $core.Iterable<IPSetSummary>? ipsets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (ipsets != null) result.ipsets.addAll(ipsets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListIPSetsResponse._();

  factory ListIPSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListIPSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListIPSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<IPSetSummary>(434949030, _omitFieldNames ? '' : 'ipsets',
        subBuilder: IPSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIPSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListIPSetsResponse copyWith(void Function(ListIPSetsResponse) updates) =>
      super.copyWith((message) => updates(message as ListIPSetsResponse))
          as ListIPSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListIPSetsResponse create() => ListIPSetsResponse._();
  @$core.override
  ListIPSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListIPSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListIPSetsResponse>(create);
  static ListIPSetsResponse? _defaultInstance;

  @$pb.TagNumber(434949030)
  $pb.PbList<IPSetSummary> get ipsets => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListLoggingConfigurationsRequest extends $pb.GeneratedMessage {
  factory ListLoggingConfigurationsRequest({
    Scope? scope,
    LogScope? logscope,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (logscope != null) result.logscope = logscope;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListLoggingConfigurationsRequest._();

  factory ListLoggingConfigurationsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListLoggingConfigurationsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListLoggingConfigurationsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aE<LogScope>(188235840, _omitFieldNames ? '' : 'logscope',
        enumValues: LogScope.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListLoggingConfigurationsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListLoggingConfigurationsRequest copyWith(
          void Function(ListLoggingConfigurationsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListLoggingConfigurationsRequest))
          as ListLoggingConfigurationsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListLoggingConfigurationsRequest create() =>
      ListLoggingConfigurationsRequest._();
  @$core.override
  ListLoggingConfigurationsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListLoggingConfigurationsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListLoggingConfigurationsRequest>(
          create);
  static ListLoggingConfigurationsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(188235840)
  LogScope get logscope => $_getN(1);
  @$pb.TagNumber(188235840)
  set logscope(LogScope value) => $_setField(188235840, value);
  @$pb.TagNumber(188235840)
  $core.bool hasLogscope() => $_has(1);
  @$pb.TagNumber(188235840)
  void clearLogscope() => $_clearField(188235840);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(3);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(3, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(3);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListLoggingConfigurationsResponse extends $pb.GeneratedMessage {
  factory ListLoggingConfigurationsResponse({
    $core.Iterable<LoggingConfiguration>? loggingconfigurations,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (loggingconfigurations != null)
      result.loggingconfigurations.addAll(loggingconfigurations);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListLoggingConfigurationsResponse._();

  factory ListLoggingConfigurationsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListLoggingConfigurationsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListLoggingConfigurationsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<LoggingConfiguration>(
        387361734, _omitFieldNames ? '' : 'loggingconfigurations',
        subBuilder: LoggingConfiguration.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListLoggingConfigurationsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListLoggingConfigurationsResponse copyWith(
          void Function(ListLoggingConfigurationsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListLoggingConfigurationsResponse))
          as ListLoggingConfigurationsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListLoggingConfigurationsResponse create() =>
      ListLoggingConfigurationsResponse._();
  @$core.override
  ListLoggingConfigurationsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListLoggingConfigurationsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListLoggingConfigurationsResponse>(
          create);
  static ListLoggingConfigurationsResponse? _defaultInstance;

  @$pb.TagNumber(387361734)
  $pb.PbList<LoggingConfiguration> get loggingconfigurations => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListManagedRuleSetsRequest extends $pb.GeneratedMessage {
  factory ListManagedRuleSetsRequest({
    Scope? scope,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListManagedRuleSetsRequest._();

  factory ListManagedRuleSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListManagedRuleSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListManagedRuleSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListManagedRuleSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListManagedRuleSetsRequest copyWith(
          void Function(ListManagedRuleSetsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListManagedRuleSetsRequest))
          as ListManagedRuleSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListManagedRuleSetsRequest create() => ListManagedRuleSetsRequest._();
  @$core.override
  ListManagedRuleSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListManagedRuleSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListManagedRuleSetsRequest>(create);
  static ListManagedRuleSetsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListManagedRuleSetsResponse extends $pb.GeneratedMessage {
  factory ListManagedRuleSetsResponse({
    $core.Iterable<ManagedRuleSetSummary>? managedrulesets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (managedrulesets != null) result.managedrulesets.addAll(managedrulesets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListManagedRuleSetsResponse._();

  factory ListManagedRuleSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListManagedRuleSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListManagedRuleSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<ManagedRuleSetSummary>(
        141451946, _omitFieldNames ? '' : 'managedrulesets',
        subBuilder: ManagedRuleSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListManagedRuleSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListManagedRuleSetsResponse copyWith(
          void Function(ListManagedRuleSetsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListManagedRuleSetsResponse))
          as ListManagedRuleSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListManagedRuleSetsResponse create() =>
      ListManagedRuleSetsResponse._();
  @$core.override
  ListManagedRuleSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListManagedRuleSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListManagedRuleSetsResponse>(create);
  static ListManagedRuleSetsResponse? _defaultInstance;

  @$pb.TagNumber(141451946)
  $pb.PbList<ManagedRuleSetSummary> get managedrulesets => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListMobileSdkReleasesRequest extends $pb.GeneratedMessage {
  factory ListMobileSdkReleasesRequest({
    $core.int? limit,
    Platform? platform,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (platform != null) result.platform = platform;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListMobileSdkReleasesRequest._();

  factory ListMobileSdkReleasesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMobileSdkReleasesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMobileSdkReleasesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aE<Platform>(468905683, _omitFieldNames ? '' : 'platform',
        enumValues: Platform.values)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMobileSdkReleasesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMobileSdkReleasesRequest copyWith(
          void Function(ListMobileSdkReleasesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListMobileSdkReleasesRequest))
          as ListMobileSdkReleasesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMobileSdkReleasesRequest create() =>
      ListMobileSdkReleasesRequest._();
  @$core.override
  ListMobileSdkReleasesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMobileSdkReleasesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMobileSdkReleasesRequest>(create);
  static ListMobileSdkReleasesRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(468905683)
  Platform get platform => $_getN(1);
  @$pb.TagNumber(468905683)
  set platform(Platform value) => $_setField(468905683, value);
  @$pb.TagNumber(468905683)
  $core.bool hasPlatform() => $_has(1);
  @$pb.TagNumber(468905683)
  void clearPlatform() => $_clearField(468905683);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListMobileSdkReleasesResponse extends $pb.GeneratedMessage {
  factory ListMobileSdkReleasesResponse({
    $core.Iterable<ReleaseSummary>? releasesummaries,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (releasesummaries != null)
      result.releasesummaries.addAll(releasesummaries);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListMobileSdkReleasesResponse._();

  factory ListMobileSdkReleasesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMobileSdkReleasesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMobileSdkReleasesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<ReleaseSummary>(494914987, _omitFieldNames ? '' : 'releasesummaries',
        subBuilder: ReleaseSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMobileSdkReleasesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMobileSdkReleasesResponse copyWith(
          void Function(ListMobileSdkReleasesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListMobileSdkReleasesResponse))
          as ListMobileSdkReleasesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMobileSdkReleasesResponse create() =>
      ListMobileSdkReleasesResponse._();
  @$core.override
  ListMobileSdkReleasesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMobileSdkReleasesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMobileSdkReleasesResponse>(create);
  static ListMobileSdkReleasesResponse? _defaultInstance;

  @$pb.TagNumber(494914987)
  $pb.PbList<ReleaseSummary> get releasesummaries => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListRegexPatternSetsRequest extends $pb.GeneratedMessage {
  factory ListRegexPatternSetsRequest({
    Scope? scope,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRegexPatternSetsRequest._();

  factory ListRegexPatternSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRegexPatternSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRegexPatternSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRegexPatternSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRegexPatternSetsRequest copyWith(
          void Function(ListRegexPatternSetsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListRegexPatternSetsRequest))
          as ListRegexPatternSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRegexPatternSetsRequest create() =>
      ListRegexPatternSetsRequest._();
  @$core.override
  ListRegexPatternSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRegexPatternSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRegexPatternSetsRequest>(create);
  static ListRegexPatternSetsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListRegexPatternSetsResponse extends $pb.GeneratedMessage {
  factory ListRegexPatternSetsResponse({
    $core.Iterable<RegexPatternSetSummary>? regexpatternsets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (regexpatternsets != null)
      result.regexpatternsets.addAll(regexpatternsets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRegexPatternSetsResponse._();

  factory ListRegexPatternSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRegexPatternSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRegexPatternSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<RegexPatternSetSummary>(
        305199780, _omitFieldNames ? '' : 'regexpatternsets',
        subBuilder: RegexPatternSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRegexPatternSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRegexPatternSetsResponse copyWith(
          void Function(ListRegexPatternSetsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListRegexPatternSetsResponse))
          as ListRegexPatternSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRegexPatternSetsResponse create() =>
      ListRegexPatternSetsResponse._();
  @$core.override
  ListRegexPatternSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRegexPatternSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRegexPatternSetsResponse>(create);
  static ListRegexPatternSetsResponse? _defaultInstance;

  @$pb.TagNumber(305199780)
  $pb.PbList<RegexPatternSetSummary> get regexpatternsets => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListResourcesForWebACLRequest extends $pb.GeneratedMessage {
  factory ListResourcesForWebACLRequest({
    $core.String? webaclarn,
    ResourceType? resourcetype,
  }) {
    final result = create();
    if (webaclarn != null) result.webaclarn = webaclarn;
    if (resourcetype != null) result.resourcetype = resourcetype;
    return result;
  }

  ListResourcesForWebACLRequest._();

  factory ListResourcesForWebACLRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListResourcesForWebACLRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListResourcesForWebACLRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(82506659, _omitFieldNames ? '' : 'webaclarn')
    ..aE<ResourceType>(301342558, _omitFieldNames ? '' : 'resourcetype',
        enumValues: ResourceType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourcesForWebACLRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourcesForWebACLRequest copyWith(
          void Function(ListResourcesForWebACLRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListResourcesForWebACLRequest))
          as ListResourcesForWebACLRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListResourcesForWebACLRequest create() =>
      ListResourcesForWebACLRequest._();
  @$core.override
  ListResourcesForWebACLRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListResourcesForWebACLRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListResourcesForWebACLRequest>(create);
  static ListResourcesForWebACLRequest? _defaultInstance;

  @$pb.TagNumber(82506659)
  $core.String get webaclarn => $_getSZ(0);
  @$pb.TagNumber(82506659)
  set webaclarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82506659)
  $core.bool hasWebaclarn() => $_has(0);
  @$pb.TagNumber(82506659)
  void clearWebaclarn() => $_clearField(82506659);

  @$pb.TagNumber(301342558)
  ResourceType get resourcetype => $_getN(1);
  @$pb.TagNumber(301342558)
  set resourcetype(ResourceType value) => $_setField(301342558, value);
  @$pb.TagNumber(301342558)
  $core.bool hasResourcetype() => $_has(1);
  @$pb.TagNumber(301342558)
  void clearResourcetype() => $_clearField(301342558);
}

class ListResourcesForWebACLResponse extends $pb.GeneratedMessage {
  factory ListResourcesForWebACLResponse({
    $core.Iterable<$core.String>? resourcearns,
  }) {
    final result = create();
    if (resourcearns != null) result.resourcearns.addAll(resourcearns);
    return result;
  }

  ListResourcesForWebACLResponse._();

  factory ListResourcesForWebACLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListResourcesForWebACLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListResourcesForWebACLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(222677390, _omitFieldNames ? '' : 'resourcearns')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourcesForWebACLResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourcesForWebACLResponse copyWith(
          void Function(ListResourcesForWebACLResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListResourcesForWebACLResponse))
          as ListResourcesForWebACLResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListResourcesForWebACLResponse create() =>
      ListResourcesForWebACLResponse._();
  @$core.override
  ListResourcesForWebACLResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListResourcesForWebACLResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListResourcesForWebACLResponse>(create);
  static ListResourcesForWebACLResponse? _defaultInstance;

  @$pb.TagNumber(222677390)
  $pb.PbList<$core.String> get resourcearns => $_getList(0);
}

class ListRuleGroupsRequest extends $pb.GeneratedMessage {
  factory ListRuleGroupsRequest({
    Scope? scope,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRuleGroupsRequest._();

  factory ListRuleGroupsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRuleGroupsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRuleGroupsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRuleGroupsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRuleGroupsRequest copyWith(
          void Function(ListRuleGroupsRequest) updates) =>
      super.copyWith((message) => updates(message as ListRuleGroupsRequest))
          as ListRuleGroupsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRuleGroupsRequest create() => ListRuleGroupsRequest._();
  @$core.override
  ListRuleGroupsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRuleGroupsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRuleGroupsRequest>(create);
  static ListRuleGroupsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListRuleGroupsResponse extends $pb.GeneratedMessage {
  factory ListRuleGroupsResponse({
    $core.Iterable<RuleGroupSummary>? rulegroups,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (rulegroups != null) result.rulegroups.addAll(rulegroups);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRuleGroupsResponse._();

  factory ListRuleGroupsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRuleGroupsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRuleGroupsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<RuleGroupSummary>(270019740, _omitFieldNames ? '' : 'rulegroups',
        subBuilder: RuleGroupSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRuleGroupsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRuleGroupsResponse copyWith(
          void Function(ListRuleGroupsResponse) updates) =>
      super.copyWith((message) => updates(message as ListRuleGroupsResponse))
          as ListRuleGroupsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRuleGroupsResponse create() => ListRuleGroupsResponse._();
  @$core.override
  ListRuleGroupsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRuleGroupsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRuleGroupsResponse>(create);
  static ListRuleGroupsResponse? _defaultInstance;

  @$pb.TagNumber(270019740)
  $pb.PbList<RuleGroupSummary> get rulegroups => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListTagsForResourceRequest extends $pb.GeneratedMessage {
  factory ListTagsForResourceRequest({
    $core.String? resourcearn,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
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

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListTagsForResourceResponse extends $pb.GeneratedMessage {
  factory ListTagsForResourceResponse({
    TagInfoForResource? taginfoforresource,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (taginfoforresource != null)
      result.taginfoforresource = taginfoforresource;
    if (nextmarker != null) result.nextmarker = nextmarker;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<TagInfoForResource>(
        8132955, _omitFieldNames ? '' : 'taginfoforresource',
        subBuilder: TagInfoForResource.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
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

  @$pb.TagNumber(8132955)
  TagInfoForResource get taginfoforresource => $_getN(0);
  @$pb.TagNumber(8132955)
  set taginfoforresource(TagInfoForResource value) =>
      $_setField(8132955, value);
  @$pb.TagNumber(8132955)
  $core.bool hasTaginfoforresource() => $_has(0);
  @$pb.TagNumber(8132955)
  void clearTaginfoforresource() => $_clearField(8132955);
  @$pb.TagNumber(8132955)
  TagInfoForResource ensureTaginfoforresource() => $_ensure(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListWebACLsRequest extends $pb.GeneratedMessage {
  factory ListWebACLsRequest({
    Scope? scope,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (scope != null) result.scope = scope;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListWebACLsRequest._();

  factory ListWebACLsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListWebACLsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListWebACLsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWebACLsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWebACLsRequest copyWith(void Function(ListWebACLsRequest) updates) =>
      super.copyWith((message) => updates(message as ListWebACLsRequest))
          as ListWebACLsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListWebACLsRequest create() => ListWebACLsRequest._();
  @$core.override
  ListWebACLsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListWebACLsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListWebACLsRequest>(create);
  static ListWebACLsRequest? _defaultInstance;

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(0);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(0);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(2);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(2);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListWebACLsResponse extends $pb.GeneratedMessage {
  factory ListWebACLsResponse({
    $core.Iterable<WebACLSummary>? webacls,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (webacls != null) result.webacls.addAll(webacls);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListWebACLsResponse._();

  factory ListWebACLsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListWebACLsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListWebACLsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<WebACLSummary>(68158245, _omitFieldNames ? '' : 'webacls',
        subBuilder: WebACLSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWebACLsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWebACLsResponse copyWith(void Function(ListWebACLsResponse) updates) =>
      super.copyWith((message) => updates(message as ListWebACLsResponse))
          as ListWebACLsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListWebACLsResponse create() => ListWebACLsResponse._();
  @$core.override
  ListWebACLsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListWebACLsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListWebACLsResponse>(create);
  static ListWebACLsResponse? _defaultInstance;

  @$pb.TagNumber(68158245)
  $pb.PbList<WebACLSummary> get webacls => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class LoggingConfiguration extends $pb.GeneratedMessage {
  factory LoggingConfiguration({
    $core.Iterable<$core.String>? logdestinationconfigs,
    LoggingFilter? loggingfilter,
    $core.Iterable<FieldToMatch>? redactedfields,
    LogType? logtype,
    LogScope? logscope,
    $core.String? resourcearn,
    $core.bool? managedbyfirewallmanager,
  }) {
    final result = create();
    if (logdestinationconfigs != null)
      result.logdestinationconfigs.addAll(logdestinationconfigs);
    if (loggingfilter != null) result.loggingfilter = loggingfilter;
    if (redactedfields != null) result.redactedfields.addAll(redactedfields);
    if (logtype != null) result.logtype = logtype;
    if (logscope != null) result.logscope = logscope;
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (managedbyfirewallmanager != null)
      result.managedbyfirewallmanager = managedbyfirewallmanager;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(22070207, _omitFieldNames ? '' : 'logdestinationconfigs')
    ..aOM<LoggingFilter>(28441529, _omitFieldNames ? '' : 'loggingfilter',
        subBuilder: LoggingFilter.create)
    ..pPM<FieldToMatch>(48530737, _omitFieldNames ? '' : 'redactedfields',
        subBuilder: FieldToMatch.create)
    ..aE<LogType>(116703930, _omitFieldNames ? '' : 'logtype',
        enumValues: LogType.values)
    ..aE<LogScope>(188235840, _omitFieldNames ? '' : 'logscope',
        enumValues: LogScope.values)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..aOB(415792743, _omitFieldNames ? '' : 'managedbyfirewallmanager')
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

  @$pb.TagNumber(22070207)
  $pb.PbList<$core.String> get logdestinationconfigs => $_getList(0);

  @$pb.TagNumber(28441529)
  LoggingFilter get loggingfilter => $_getN(1);
  @$pb.TagNumber(28441529)
  set loggingfilter(LoggingFilter value) => $_setField(28441529, value);
  @$pb.TagNumber(28441529)
  $core.bool hasLoggingfilter() => $_has(1);
  @$pb.TagNumber(28441529)
  void clearLoggingfilter() => $_clearField(28441529);
  @$pb.TagNumber(28441529)
  LoggingFilter ensureLoggingfilter() => $_ensure(1);

  @$pb.TagNumber(48530737)
  $pb.PbList<FieldToMatch> get redactedfields => $_getList(2);

  @$pb.TagNumber(116703930)
  LogType get logtype => $_getN(3);
  @$pb.TagNumber(116703930)
  set logtype(LogType value) => $_setField(116703930, value);
  @$pb.TagNumber(116703930)
  $core.bool hasLogtype() => $_has(3);
  @$pb.TagNumber(116703930)
  void clearLogtype() => $_clearField(116703930);

  @$pb.TagNumber(188235840)
  LogScope get logscope => $_getN(4);
  @$pb.TagNumber(188235840)
  set logscope(LogScope value) => $_setField(188235840, value);
  @$pb.TagNumber(188235840)
  $core.bool hasLogscope() => $_has(4);
  @$pb.TagNumber(188235840)
  void clearLogscope() => $_clearField(188235840);

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(5);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(5);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

  @$pb.TagNumber(415792743)
  $core.bool get managedbyfirewallmanager => $_getBF(6);
  @$pb.TagNumber(415792743)
  set managedbyfirewallmanager($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(415792743)
  $core.bool hasManagedbyfirewallmanager() => $_has(6);
  @$pb.TagNumber(415792743)
  void clearManagedbyfirewallmanager() => $_clearField(415792743);
}

class LoggingFilter extends $pb.GeneratedMessage {
  factory LoggingFilter({
    $core.Iterable<Filter>? filters,
    FilterBehavior? defaultbehavior,
  }) {
    final result = create();
    if (filters != null) result.filters.addAll(filters);
    if (defaultbehavior != null) result.defaultbehavior = defaultbehavior;
    return result;
  }

  LoggingFilter._();

  factory LoggingFilter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LoggingFilter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LoggingFilter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<Filter>(188393197, _omitFieldNames ? '' : 'filters',
        subBuilder: Filter.create)
    ..aE<FilterBehavior>(294212489, _omitFieldNames ? '' : 'defaultbehavior',
        enumValues: FilterBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoggingFilter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoggingFilter copyWith(void Function(LoggingFilter) updates) =>
      super.copyWith((message) => updates(message as LoggingFilter))
          as LoggingFilter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LoggingFilter create() => LoggingFilter._();
  @$core.override
  LoggingFilter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LoggingFilter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LoggingFilter>(create);
  static LoggingFilter? _defaultInstance;

  @$pb.TagNumber(188393197)
  $pb.PbList<Filter> get filters => $_getList(0);

  @$pb.TagNumber(294212489)
  FilterBehavior get defaultbehavior => $_getN(1);
  @$pb.TagNumber(294212489)
  set defaultbehavior(FilterBehavior value) => $_setField(294212489, value);
  @$pb.TagNumber(294212489)
  $core.bool hasDefaultbehavior() => $_has(1);
  @$pb.TagNumber(294212489)
  void clearDefaultbehavior() => $_clearField(294212489);
}

class ManagedProductDescriptor extends $pb.GeneratedMessage {
  factory ManagedProductDescriptor({
    $core.String? productid,
    $core.String? productdescription,
    $core.String? vendorname,
    $core.String? managedrulesetname,
    $core.String? productlink,
    $core.bool? isversioningsupported,
    $core.bool? isadvancedmanagedruleset,
    $core.String? snstopicarn,
    $core.String? producttitle,
  }) {
    final result = create();
    if (productid != null) result.productid = productid;
    if (productdescription != null)
      result.productdescription = productdescription;
    if (vendorname != null) result.vendorname = vendorname;
    if (managedrulesetname != null)
      result.managedrulesetname = managedrulesetname;
    if (productlink != null) result.productlink = productlink;
    if (isversioningsupported != null)
      result.isversioningsupported = isversioningsupported;
    if (isadvancedmanagedruleset != null)
      result.isadvancedmanagedruleset = isadvancedmanagedruleset;
    if (snstopicarn != null) result.snstopicarn = snstopicarn;
    if (producttitle != null) result.producttitle = producttitle;
    return result;
  }

  ManagedProductDescriptor._();

  factory ManagedProductDescriptor.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedProductDescriptor.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedProductDescriptor',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(53264354, _omitFieldNames ? '' : 'productid')
    ..aOS(78087559, _omitFieldNames ? '' : 'productdescription')
    ..aOS(159329689, _omitFieldNames ? '' : 'vendorname')
    ..aOS(161099140, _omitFieldNames ? '' : 'managedrulesetname')
    ..aOS(191063997, _omitFieldNames ? '' : 'productlink')
    ..aOB(195956662, _omitFieldNames ? '' : 'isversioningsupported')
    ..aOB(288313871, _omitFieldNames ? '' : 'isadvancedmanagedruleset')
    ..aOS(374815596, _omitFieldNames ? '' : 'snstopicarn')
    ..aOS(452337643, _omitFieldNames ? '' : 'producttitle')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedProductDescriptor clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedProductDescriptor copyWith(
          void Function(ManagedProductDescriptor) updates) =>
      super.copyWith((message) => updates(message as ManagedProductDescriptor))
          as ManagedProductDescriptor;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedProductDescriptor create() => ManagedProductDescriptor._();
  @$core.override
  ManagedProductDescriptor createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedProductDescriptor getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedProductDescriptor>(create);
  static ManagedProductDescriptor? _defaultInstance;

  @$pb.TagNumber(53264354)
  $core.String get productid => $_getSZ(0);
  @$pb.TagNumber(53264354)
  set productid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(53264354)
  $core.bool hasProductid() => $_has(0);
  @$pb.TagNumber(53264354)
  void clearProductid() => $_clearField(53264354);

  @$pb.TagNumber(78087559)
  $core.String get productdescription => $_getSZ(1);
  @$pb.TagNumber(78087559)
  set productdescription($core.String value) => $_setString(1, value);
  @$pb.TagNumber(78087559)
  $core.bool hasProductdescription() => $_has(1);
  @$pb.TagNumber(78087559)
  void clearProductdescription() => $_clearField(78087559);

  @$pb.TagNumber(159329689)
  $core.String get vendorname => $_getSZ(2);
  @$pb.TagNumber(159329689)
  set vendorname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(159329689)
  $core.bool hasVendorname() => $_has(2);
  @$pb.TagNumber(159329689)
  void clearVendorname() => $_clearField(159329689);

  @$pb.TagNumber(161099140)
  $core.String get managedrulesetname => $_getSZ(3);
  @$pb.TagNumber(161099140)
  set managedrulesetname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(161099140)
  $core.bool hasManagedrulesetname() => $_has(3);
  @$pb.TagNumber(161099140)
  void clearManagedrulesetname() => $_clearField(161099140);

  @$pb.TagNumber(191063997)
  $core.String get productlink => $_getSZ(4);
  @$pb.TagNumber(191063997)
  set productlink($core.String value) => $_setString(4, value);
  @$pb.TagNumber(191063997)
  $core.bool hasProductlink() => $_has(4);
  @$pb.TagNumber(191063997)
  void clearProductlink() => $_clearField(191063997);

  @$pb.TagNumber(195956662)
  $core.bool get isversioningsupported => $_getBF(5);
  @$pb.TagNumber(195956662)
  set isversioningsupported($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(195956662)
  $core.bool hasIsversioningsupported() => $_has(5);
  @$pb.TagNumber(195956662)
  void clearIsversioningsupported() => $_clearField(195956662);

  @$pb.TagNumber(288313871)
  $core.bool get isadvancedmanagedruleset => $_getBF(6);
  @$pb.TagNumber(288313871)
  set isadvancedmanagedruleset($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(288313871)
  $core.bool hasIsadvancedmanagedruleset() => $_has(6);
  @$pb.TagNumber(288313871)
  void clearIsadvancedmanagedruleset() => $_clearField(288313871);

  @$pb.TagNumber(374815596)
  $core.String get snstopicarn => $_getSZ(7);
  @$pb.TagNumber(374815596)
  set snstopicarn($core.String value) => $_setString(7, value);
  @$pb.TagNumber(374815596)
  $core.bool hasSnstopicarn() => $_has(7);
  @$pb.TagNumber(374815596)
  void clearSnstopicarn() => $_clearField(374815596);

  @$pb.TagNumber(452337643)
  $core.String get producttitle => $_getSZ(8);
  @$pb.TagNumber(452337643)
  set producttitle($core.String value) => $_setString(8, value);
  @$pb.TagNumber(452337643)
  $core.bool hasProducttitle() => $_has(8);
  @$pb.TagNumber(452337643)
  void clearProducttitle() => $_clearField(452337643);
}

class ManagedRuleGroupConfig extends $pb.GeneratedMessage {
  factory ManagedRuleGroupConfig({
    UsernameField? usernamefield,
    $core.String? loginpath,
    AWSManagedRulesATPRuleSet? awsmanagedrulesatpruleset,
    AWSManagedRulesBotControlRuleSet? awsmanagedrulesbotcontrolruleset,
    PasswordField? passwordfield,
    AWSManagedRulesACFPRuleSet? awsmanagedrulesacfpruleset,
    PayloadType? payloadtype,
    AWSManagedRulesAntiDDoSRuleSet? awsmanagedrulesantiddosruleset,
  }) {
    final result = create();
    if (usernamefield != null) result.usernamefield = usernamefield;
    if (loginpath != null) result.loginpath = loginpath;
    if (awsmanagedrulesatpruleset != null)
      result.awsmanagedrulesatpruleset = awsmanagedrulesatpruleset;
    if (awsmanagedrulesbotcontrolruleset != null)
      result.awsmanagedrulesbotcontrolruleset =
          awsmanagedrulesbotcontrolruleset;
    if (passwordfield != null) result.passwordfield = passwordfield;
    if (awsmanagedrulesacfpruleset != null)
      result.awsmanagedrulesacfpruleset = awsmanagedrulesacfpruleset;
    if (payloadtype != null) result.payloadtype = payloadtype;
    if (awsmanagedrulesantiddosruleset != null)
      result.awsmanagedrulesantiddosruleset = awsmanagedrulesantiddosruleset;
    return result;
  }

  ManagedRuleGroupConfig._();

  factory ManagedRuleGroupConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleGroupConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleGroupConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<UsernameField>(125830068, _omitFieldNames ? '' : 'usernamefield',
        subBuilder: UsernameField.create)
    ..aOS(128281874, _omitFieldNames ? '' : 'loginpath')
    ..aOM<AWSManagedRulesATPRuleSet>(
        263010970, _omitFieldNames ? '' : 'awsmanagedrulesatpruleset',
        subBuilder: AWSManagedRulesATPRuleSet.create)
    ..aOM<AWSManagedRulesBotControlRuleSet>(
        283828283, _omitFieldNames ? '' : 'awsmanagedrulesbotcontrolruleset',
        subBuilder: AWSManagedRulesBotControlRuleSet.create)
    ..aOM<PasswordField>(318147221, _omitFieldNames ? '' : 'passwordfield',
        subBuilder: PasswordField.create)
    ..aOM<AWSManagedRulesACFPRuleSet>(
        356597463, _omitFieldNames ? '' : 'awsmanagedrulesacfpruleset',
        subBuilder: AWSManagedRulesACFPRuleSet.create)
    ..aE<PayloadType>(510845422, _omitFieldNames ? '' : 'payloadtype',
        enumValues: PayloadType.values)
    ..aOM<AWSManagedRulesAntiDDoSRuleSet>(
        517768335, _omitFieldNames ? '' : 'awsmanagedrulesantiddosruleset',
        subBuilder: AWSManagedRulesAntiDDoSRuleSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleGroupConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleGroupConfig copyWith(
          void Function(ManagedRuleGroupConfig) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleGroupConfig))
          as ManagedRuleGroupConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleGroupConfig create() => ManagedRuleGroupConfig._();
  @$core.override
  ManagedRuleGroupConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleGroupConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleGroupConfig>(create);
  static ManagedRuleGroupConfig? _defaultInstance;

  @$pb.TagNumber(125830068)
  UsernameField get usernamefield => $_getN(0);
  @$pb.TagNumber(125830068)
  set usernamefield(UsernameField value) => $_setField(125830068, value);
  @$pb.TagNumber(125830068)
  $core.bool hasUsernamefield() => $_has(0);
  @$pb.TagNumber(125830068)
  void clearUsernamefield() => $_clearField(125830068);
  @$pb.TagNumber(125830068)
  UsernameField ensureUsernamefield() => $_ensure(0);

  @$pb.TagNumber(128281874)
  $core.String get loginpath => $_getSZ(1);
  @$pb.TagNumber(128281874)
  set loginpath($core.String value) => $_setString(1, value);
  @$pb.TagNumber(128281874)
  $core.bool hasLoginpath() => $_has(1);
  @$pb.TagNumber(128281874)
  void clearLoginpath() => $_clearField(128281874);

  @$pb.TagNumber(263010970)
  AWSManagedRulesATPRuleSet get awsmanagedrulesatpruleset => $_getN(2);
  @$pb.TagNumber(263010970)
  set awsmanagedrulesatpruleset(AWSManagedRulesATPRuleSet value) =>
      $_setField(263010970, value);
  @$pb.TagNumber(263010970)
  $core.bool hasAwsmanagedrulesatpruleset() => $_has(2);
  @$pb.TagNumber(263010970)
  void clearAwsmanagedrulesatpruleset() => $_clearField(263010970);
  @$pb.TagNumber(263010970)
  AWSManagedRulesATPRuleSet ensureAwsmanagedrulesatpruleset() => $_ensure(2);

  @$pb.TagNumber(283828283)
  AWSManagedRulesBotControlRuleSet get awsmanagedrulesbotcontrolruleset =>
      $_getN(3);
  @$pb.TagNumber(283828283)
  set awsmanagedrulesbotcontrolruleset(
          AWSManagedRulesBotControlRuleSet value) =>
      $_setField(283828283, value);
  @$pb.TagNumber(283828283)
  $core.bool hasAwsmanagedrulesbotcontrolruleset() => $_has(3);
  @$pb.TagNumber(283828283)
  void clearAwsmanagedrulesbotcontrolruleset() => $_clearField(283828283);
  @$pb.TagNumber(283828283)
  AWSManagedRulesBotControlRuleSet ensureAwsmanagedrulesbotcontrolruleset() =>
      $_ensure(3);

  @$pb.TagNumber(318147221)
  PasswordField get passwordfield => $_getN(4);
  @$pb.TagNumber(318147221)
  set passwordfield(PasswordField value) => $_setField(318147221, value);
  @$pb.TagNumber(318147221)
  $core.bool hasPasswordfield() => $_has(4);
  @$pb.TagNumber(318147221)
  void clearPasswordfield() => $_clearField(318147221);
  @$pb.TagNumber(318147221)
  PasswordField ensurePasswordfield() => $_ensure(4);

  @$pb.TagNumber(356597463)
  AWSManagedRulesACFPRuleSet get awsmanagedrulesacfpruleset => $_getN(5);
  @$pb.TagNumber(356597463)
  set awsmanagedrulesacfpruleset(AWSManagedRulesACFPRuleSet value) =>
      $_setField(356597463, value);
  @$pb.TagNumber(356597463)
  $core.bool hasAwsmanagedrulesacfpruleset() => $_has(5);
  @$pb.TagNumber(356597463)
  void clearAwsmanagedrulesacfpruleset() => $_clearField(356597463);
  @$pb.TagNumber(356597463)
  AWSManagedRulesACFPRuleSet ensureAwsmanagedrulesacfpruleset() => $_ensure(5);

  @$pb.TagNumber(510845422)
  PayloadType get payloadtype => $_getN(6);
  @$pb.TagNumber(510845422)
  set payloadtype(PayloadType value) => $_setField(510845422, value);
  @$pb.TagNumber(510845422)
  $core.bool hasPayloadtype() => $_has(6);
  @$pb.TagNumber(510845422)
  void clearPayloadtype() => $_clearField(510845422);

  @$pb.TagNumber(517768335)
  AWSManagedRulesAntiDDoSRuleSet get awsmanagedrulesantiddosruleset =>
      $_getN(7);
  @$pb.TagNumber(517768335)
  set awsmanagedrulesantiddosruleset(AWSManagedRulesAntiDDoSRuleSet value) =>
      $_setField(517768335, value);
  @$pb.TagNumber(517768335)
  $core.bool hasAwsmanagedrulesantiddosruleset() => $_has(7);
  @$pb.TagNumber(517768335)
  void clearAwsmanagedrulesantiddosruleset() => $_clearField(517768335);
  @$pb.TagNumber(517768335)
  AWSManagedRulesAntiDDoSRuleSet ensureAwsmanagedrulesantiddosruleset() =>
      $_ensure(7);
}

class ManagedRuleGroupStatement extends $pb.GeneratedMessage {
  factory ManagedRuleGroupStatement({
    Statement? scopedownstatement,
    $core.Iterable<ExcludedRule>? excludedrules,
    $core.String? vendorname,
    $core.Iterable<ManagedRuleGroupConfig>? managedrulegroupconfigs,
    $core.String? name,
    $core.Iterable<RuleActionOverride>? ruleactionoverrides,
    $core.String? version,
  }) {
    final result = create();
    if (scopedownstatement != null)
      result.scopedownstatement = scopedownstatement;
    if (excludedrules != null) result.excludedrules.addAll(excludedrules);
    if (vendorname != null) result.vendorname = vendorname;
    if (managedrulegroupconfigs != null)
      result.managedrulegroupconfigs.addAll(managedrulegroupconfigs);
    if (name != null) result.name = name;
    if (ruleactionoverrides != null)
      result.ruleactionoverrides.addAll(ruleactionoverrides);
    if (version != null) result.version = version;
    return result;
  }

  ManagedRuleGroupStatement._();

  factory ManagedRuleGroupStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleGroupStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleGroupStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<Statement>(116108605, _omitFieldNames ? '' : 'scopedownstatement',
        subBuilder: Statement.create)
    ..pPM<ExcludedRule>(129929071, _omitFieldNames ? '' : 'excludedrules',
        subBuilder: ExcludedRule.create)
    ..aOS(159329689, _omitFieldNames ? '' : 'vendorname')
    ..pPM<ManagedRuleGroupConfig>(
        222138791, _omitFieldNames ? '' : 'managedrulegroupconfigs',
        subBuilder: ManagedRuleGroupConfig.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<RuleActionOverride>(
        354236935, _omitFieldNames ? '' : 'ruleactionoverrides',
        subBuilder: RuleActionOverride.create)
    ..aOS(500028728, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleGroupStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleGroupStatement copyWith(
          void Function(ManagedRuleGroupStatement) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleGroupStatement))
          as ManagedRuleGroupStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleGroupStatement create() => ManagedRuleGroupStatement._();
  @$core.override
  ManagedRuleGroupStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleGroupStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleGroupStatement>(create);
  static ManagedRuleGroupStatement? _defaultInstance;

  @$pb.TagNumber(116108605)
  Statement get scopedownstatement => $_getN(0);
  @$pb.TagNumber(116108605)
  set scopedownstatement(Statement value) => $_setField(116108605, value);
  @$pb.TagNumber(116108605)
  $core.bool hasScopedownstatement() => $_has(0);
  @$pb.TagNumber(116108605)
  void clearScopedownstatement() => $_clearField(116108605);
  @$pb.TagNumber(116108605)
  Statement ensureScopedownstatement() => $_ensure(0);

  @$pb.TagNumber(129929071)
  $pb.PbList<ExcludedRule> get excludedrules => $_getList(1);

  @$pb.TagNumber(159329689)
  $core.String get vendorname => $_getSZ(2);
  @$pb.TagNumber(159329689)
  set vendorname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(159329689)
  $core.bool hasVendorname() => $_has(2);
  @$pb.TagNumber(159329689)
  void clearVendorname() => $_clearField(159329689);

  @$pb.TagNumber(222138791)
  $pb.PbList<ManagedRuleGroupConfig> get managedrulegroupconfigs =>
      $_getList(3);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(354236935)
  $pb.PbList<RuleActionOverride> get ruleactionoverrides => $_getList(5);

  @$pb.TagNumber(500028728)
  $core.String get version => $_getSZ(6);
  @$pb.TagNumber(500028728)
  set version($core.String value) => $_setString(6, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(6);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);
}

class ManagedRuleGroupSummary extends $pb.GeneratedMessage {
  factory ManagedRuleGroupSummary({
    $core.String? description,
    $core.String? vendorname,
    $core.bool? versioningsupported,
    $core.String? name,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (vendorname != null) result.vendorname = vendorname;
    if (versioningsupported != null)
      result.versioningsupported = versioningsupported;
    if (name != null) result.name = name;
    return result;
  }

  ManagedRuleGroupSummary._();

  factory ManagedRuleGroupSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleGroupSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleGroupSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(159329689, _omitFieldNames ? '' : 'vendorname')
    ..aOB(166873158, _omitFieldNames ? '' : 'versioningsupported')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleGroupSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleGroupSummary copyWith(
          void Function(ManagedRuleGroupSummary) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleGroupSummary))
          as ManagedRuleGroupSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleGroupSummary create() => ManagedRuleGroupSummary._();
  @$core.override
  ManagedRuleGroupSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleGroupSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleGroupSummary>(create);
  static ManagedRuleGroupSummary? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(159329689)
  $core.String get vendorname => $_getSZ(1);
  @$pb.TagNumber(159329689)
  set vendorname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(159329689)
  $core.bool hasVendorname() => $_has(1);
  @$pb.TagNumber(159329689)
  void clearVendorname() => $_clearField(159329689);

  @$pb.TagNumber(166873158)
  $core.bool get versioningsupported => $_getBF(2);
  @$pb.TagNumber(166873158)
  set versioningsupported($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(166873158)
  $core.bool hasVersioningsupported() => $_has(2);
  @$pb.TagNumber(166873158)
  void clearVersioningsupported() => $_clearField(166873158);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ManagedRuleGroupVersion extends $pb.GeneratedMessage {
  factory ManagedRuleGroupVersion({
    $core.String? lastupdatetimestamp,
    $core.String? name,
  }) {
    final result = create();
    if (lastupdatetimestamp != null)
      result.lastupdatetimestamp = lastupdatetimestamp;
    if (name != null) result.name = name;
    return result;
  }

  ManagedRuleGroupVersion._();

  factory ManagedRuleGroupVersion.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleGroupVersion.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleGroupVersion',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(118818825, _omitFieldNames ? '' : 'lastupdatetimestamp')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleGroupVersion clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleGroupVersion copyWith(
          void Function(ManagedRuleGroupVersion) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleGroupVersion))
          as ManagedRuleGroupVersion;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleGroupVersion create() => ManagedRuleGroupVersion._();
  @$core.override
  ManagedRuleGroupVersion createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleGroupVersion getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleGroupVersion>(create);
  static ManagedRuleGroupVersion? _defaultInstance;

  @$pb.TagNumber(118818825)
  $core.String get lastupdatetimestamp => $_getSZ(0);
  @$pb.TagNumber(118818825)
  set lastupdatetimestamp($core.String value) => $_setString(0, value);
  @$pb.TagNumber(118818825)
  $core.bool hasLastupdatetimestamp() => $_has(0);
  @$pb.TagNumber(118818825)
  void clearLastupdatetimestamp() => $_clearField(118818825);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ManagedRuleSet extends $pb.GeneratedMessage {
  factory ManagedRuleSet({
    $core.String? labelnamespace,
    $core.String? recommendedversion,
    $core.String? description,
    $core.Iterable<$core.MapEntry<$core.String, ManagedRuleSetVersion>>?
        publishedversions,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (labelnamespace != null) result.labelnamespace = labelnamespace;
    if (recommendedversion != null)
      result.recommendedversion = recommendedversion;
    if (description != null) result.description = description;
    if (publishedversions != null)
      result.publishedversions.addEntries(publishedversions);
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  ManagedRuleSet._();

  factory ManagedRuleSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(39797417, _omitFieldNames ? '' : 'labelnamespace')
    ..aOS(83103137, _omitFieldNames ? '' : 'recommendedversion')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..m<$core.String, ManagedRuleSetVersion>(
        168260833, _omitFieldNames ? '' : 'publishedversions',
        entryClassName: 'ManagedRuleSet.PublishedversionsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: ManagedRuleSetVersion.create,
        valueDefaultOrMaker: ManagedRuleSetVersion.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleSet copyWith(void Function(ManagedRuleSet) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleSet))
          as ManagedRuleSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleSet create() => ManagedRuleSet._();
  @$core.override
  ManagedRuleSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleSet>(create);
  static ManagedRuleSet? _defaultInstance;

  @$pb.TagNumber(39797417)
  $core.String get labelnamespace => $_getSZ(0);
  @$pb.TagNumber(39797417)
  set labelnamespace($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39797417)
  $core.bool hasLabelnamespace() => $_has(0);
  @$pb.TagNumber(39797417)
  void clearLabelnamespace() => $_clearField(39797417);

  @$pb.TagNumber(83103137)
  $core.String get recommendedversion => $_getSZ(1);
  @$pb.TagNumber(83103137)
  set recommendedversion($core.String value) => $_setString(1, value);
  @$pb.TagNumber(83103137)
  $core.bool hasRecommendedversion() => $_has(1);
  @$pb.TagNumber(83103137)
  void clearRecommendedversion() => $_clearField(83103137);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(168260833)
  $pb.PbMap<$core.String, ManagedRuleSetVersion> get publishedversions =>
      $_getMap(3);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(5);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(5, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(5);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(6);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(6);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class ManagedRuleSetSummary extends $pb.GeneratedMessage {
  factory ManagedRuleSetSummary({
    $core.String? labelnamespace,
    $core.String? locktoken,
    $core.String? description,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (labelnamespace != null) result.labelnamespace = labelnamespace;
    if (locktoken != null) result.locktoken = locktoken;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  ManagedRuleSetSummary._();

  factory ManagedRuleSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(39797417, _omitFieldNames ? '' : 'labelnamespace')
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleSetSummary copyWith(
          void Function(ManagedRuleSetSummary) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleSetSummary))
          as ManagedRuleSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleSetSummary create() => ManagedRuleSetSummary._();
  @$core.override
  ManagedRuleSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleSetSummary>(create);
  static ManagedRuleSetSummary? _defaultInstance;

  @$pb.TagNumber(39797417)
  $core.String get labelnamespace => $_getSZ(0);
  @$pb.TagNumber(39797417)
  set labelnamespace($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39797417)
  $core.bool hasLabelnamespace() => $_has(0);
  @$pb.TagNumber(39797417)
  void clearLabelnamespace() => $_clearField(39797417);

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(1);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(1);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(5);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(5);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class ManagedRuleSetVersion extends $pb.GeneratedMessage {
  factory ManagedRuleSetVersion({
    $fixnum.Int64? capacity,
    $core.String? lastupdatetimestamp,
    $core.int? forecastedlifetime,
    $core.String? associatedrulegrouparn,
    $core.String? expirytimestamp,
    $core.String? publishtimestamp,
  }) {
    final result = create();
    if (capacity != null) result.capacity = capacity;
    if (lastupdatetimestamp != null)
      result.lastupdatetimestamp = lastupdatetimestamp;
    if (forecastedlifetime != null)
      result.forecastedlifetime = forecastedlifetime;
    if (associatedrulegrouparn != null)
      result.associatedrulegrouparn = associatedrulegrouparn;
    if (expirytimestamp != null) result.expirytimestamp = expirytimestamp;
    if (publishtimestamp != null) result.publishtimestamp = publishtimestamp;
    return result;
  }

  ManagedRuleSetVersion._();

  factory ManagedRuleSetVersion.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleSetVersion.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleSetVersion',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aInt64(107253930, _omitFieldNames ? '' : 'capacity')
    ..aOS(118818825, _omitFieldNames ? '' : 'lastupdatetimestamp')
    ..aI(149179979, _omitFieldNames ? '' : 'forecastedlifetime')
    ..aOS(416522796, _omitFieldNames ? '' : 'associatedrulegrouparn')
    ..aOS(460460551, _omitFieldNames ? '' : 'expirytimestamp')
    ..aOS(477892057, _omitFieldNames ? '' : 'publishtimestamp')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleSetVersion clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleSetVersion copyWith(
          void Function(ManagedRuleSetVersion) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleSetVersion))
          as ManagedRuleSetVersion;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleSetVersion create() => ManagedRuleSetVersion._();
  @$core.override
  ManagedRuleSetVersion createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleSetVersion getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleSetVersion>(create);
  static ManagedRuleSetVersion? _defaultInstance;

  @$pb.TagNumber(107253930)
  $fixnum.Int64 get capacity => $_getI64(0);
  @$pb.TagNumber(107253930)
  set capacity($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(107253930)
  $core.bool hasCapacity() => $_has(0);
  @$pb.TagNumber(107253930)
  void clearCapacity() => $_clearField(107253930);

  @$pb.TagNumber(118818825)
  $core.String get lastupdatetimestamp => $_getSZ(1);
  @$pb.TagNumber(118818825)
  set lastupdatetimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(118818825)
  $core.bool hasLastupdatetimestamp() => $_has(1);
  @$pb.TagNumber(118818825)
  void clearLastupdatetimestamp() => $_clearField(118818825);

  @$pb.TagNumber(149179979)
  $core.int get forecastedlifetime => $_getIZ(2);
  @$pb.TagNumber(149179979)
  set forecastedlifetime($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(149179979)
  $core.bool hasForecastedlifetime() => $_has(2);
  @$pb.TagNumber(149179979)
  void clearForecastedlifetime() => $_clearField(149179979);

  @$pb.TagNumber(416522796)
  $core.String get associatedrulegrouparn => $_getSZ(3);
  @$pb.TagNumber(416522796)
  set associatedrulegrouparn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(416522796)
  $core.bool hasAssociatedrulegrouparn() => $_has(3);
  @$pb.TagNumber(416522796)
  void clearAssociatedrulegrouparn() => $_clearField(416522796);

  @$pb.TagNumber(460460551)
  $core.String get expirytimestamp => $_getSZ(4);
  @$pb.TagNumber(460460551)
  set expirytimestamp($core.String value) => $_setString(4, value);
  @$pb.TagNumber(460460551)
  $core.bool hasExpirytimestamp() => $_has(4);
  @$pb.TagNumber(460460551)
  void clearExpirytimestamp() => $_clearField(460460551);

  @$pb.TagNumber(477892057)
  $core.String get publishtimestamp => $_getSZ(5);
  @$pb.TagNumber(477892057)
  set publishtimestamp($core.String value) => $_setString(5, value);
  @$pb.TagNumber(477892057)
  $core.bool hasPublishtimestamp() => $_has(5);
  @$pb.TagNumber(477892057)
  void clearPublishtimestamp() => $_clearField(477892057);
}

class Method extends $pb.GeneratedMessage {
  factory Method() => create();

  Method._();

  factory Method.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Method.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Method',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
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
}

class MobileSdkRelease extends $pb.GeneratedMessage {
  factory MobileSdkRelease({
    $core.String? releaseversion,
    $core.String? releasenotes,
    $core.String? timestamp,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (releaseversion != null) result.releaseversion = releaseversion;
    if (releasenotes != null) result.releasenotes = releasenotes;
    if (timestamp != null) result.timestamp = timestamp;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  MobileSdkRelease._();

  factory MobileSdkRelease.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MobileSdkRelease.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MobileSdkRelease',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(8739171, _omitFieldNames ? '' : 'releaseversion')
    ..aOS(98800172, _omitFieldNames ? '' : 'releasenotes')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MobileSdkRelease clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MobileSdkRelease copyWith(void Function(MobileSdkRelease) updates) =>
      super.copyWith((message) => updates(message as MobileSdkRelease))
          as MobileSdkRelease;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MobileSdkRelease create() => MobileSdkRelease._();
  @$core.override
  MobileSdkRelease createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MobileSdkRelease getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MobileSdkRelease>(create);
  static MobileSdkRelease? _defaultInstance;

  @$pb.TagNumber(8739171)
  $core.String get releaseversion => $_getSZ(0);
  @$pb.TagNumber(8739171)
  set releaseversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(8739171)
  $core.bool hasReleaseversion() => $_has(0);
  @$pb.TagNumber(8739171)
  void clearReleaseversion() => $_clearField(8739171);

  @$pb.TagNumber(98800172)
  $core.String get releasenotes => $_getSZ(1);
  @$pb.TagNumber(98800172)
  set releasenotes($core.String value) => $_setString(1, value);
  @$pb.TagNumber(98800172)
  $core.bool hasReleasenotes() => $_has(1);
  @$pb.TagNumber(98800172)
  void clearReleasenotes() => $_clearField(98800172);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(2);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(2);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(3);
}

class NoneAction extends $pb.GeneratedMessage {
  factory NoneAction() => create();

  NoneAction._();

  factory NoneAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoneAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoneAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoneAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoneAction copyWith(void Function(NoneAction) updates) =>
      super.copyWith((message) => updates(message as NoneAction)) as NoneAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoneAction create() => NoneAction._();
  @$core.override
  NoneAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoneAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoneAction>(create);
  static NoneAction? _defaultInstance;
}

class NotStatement extends $pb.GeneratedMessage {
  factory NotStatement({
    Statement? statement,
  }) {
    final result = create();
    if (statement != null) result.statement = statement;
    return result;
  }

  NotStatement._();

  factory NotStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NotStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NotStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<Statement>(248790199, _omitFieldNames ? '' : 'statement',
        subBuilder: Statement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NotStatement copyWith(void Function(NotStatement) updates) =>
      super.copyWith((message) => updates(message as NotStatement))
          as NotStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NotStatement create() => NotStatement._();
  @$core.override
  NotStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NotStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NotStatement>(create);
  static NotStatement? _defaultInstance;

  @$pb.TagNumber(248790199)
  Statement get statement => $_getN(0);
  @$pb.TagNumber(248790199)
  set statement(Statement value) => $_setField(248790199, value);
  @$pb.TagNumber(248790199)
  $core.bool hasStatement() => $_has(0);
  @$pb.TagNumber(248790199)
  void clearStatement() => $_clearField(248790199);
  @$pb.TagNumber(248790199)
  Statement ensureStatement() => $_ensure(0);
}

class OnSourceDDoSProtectionConfig extends $pb.GeneratedMessage {
  factory OnSourceDDoSProtectionConfig({
    LowReputationMode? alblowreputationmode,
  }) {
    final result = create();
    if (alblowreputationmode != null)
      result.alblowreputationmode = alblowreputationmode;
    return result;
  }

  OnSourceDDoSProtectionConfig._();

  factory OnSourceDDoSProtectionConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OnSourceDDoSProtectionConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OnSourceDDoSProtectionConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<LowReputationMode>(
        236884769, _omitFieldNames ? '' : 'alblowreputationmode',
        enumValues: LowReputationMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OnSourceDDoSProtectionConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OnSourceDDoSProtectionConfig copyWith(
          void Function(OnSourceDDoSProtectionConfig) updates) =>
      super.copyWith(
              (message) => updates(message as OnSourceDDoSProtectionConfig))
          as OnSourceDDoSProtectionConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OnSourceDDoSProtectionConfig create() =>
      OnSourceDDoSProtectionConfig._();
  @$core.override
  OnSourceDDoSProtectionConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OnSourceDDoSProtectionConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OnSourceDDoSProtectionConfig>(create);
  static OnSourceDDoSProtectionConfig? _defaultInstance;

  @$pb.TagNumber(236884769)
  LowReputationMode get alblowreputationmode => $_getN(0);
  @$pb.TagNumber(236884769)
  set alblowreputationmode(LowReputationMode value) =>
      $_setField(236884769, value);
  @$pb.TagNumber(236884769)
  $core.bool hasAlblowreputationmode() => $_has(0);
  @$pb.TagNumber(236884769)
  void clearAlblowreputationmode() => $_clearField(236884769);
}

class OrStatement extends $pb.GeneratedMessage {
  factory OrStatement({
    $core.Iterable<Statement>? statements,
  }) {
    final result = create();
    if (statements != null) result.statements.addAll(statements);
    return result;
  }

  OrStatement._();

  factory OrStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OrStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OrStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<Statement>(488352288, _omitFieldNames ? '' : 'statements',
        subBuilder: Statement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OrStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OrStatement copyWith(void Function(OrStatement) updates) =>
      super.copyWith((message) => updates(message as OrStatement))
          as OrStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OrStatement create() => OrStatement._();
  @$core.override
  OrStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OrStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OrStatement>(create);
  static OrStatement? _defaultInstance;

  @$pb.TagNumber(488352288)
  $pb.PbList<Statement> get statements => $_getList(0);
}

class OverrideAction extends $pb.GeneratedMessage {
  factory OverrideAction({
    CountAction? count,
    NoneAction? none,
  }) {
    final result = create();
    if (count != null) result.count = count;
    if (none != null) result.none = none;
    return result;
  }

  OverrideAction._();

  factory OverrideAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OverrideAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OverrideAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<CountAction>(31963285, _omitFieldNames ? '' : 'count',
        subBuilder: CountAction.create)
    ..aOM<NoneAction>(273676284, _omitFieldNames ? '' : 'none',
        subBuilder: NoneAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OverrideAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OverrideAction copyWith(void Function(OverrideAction) updates) =>
      super.copyWith((message) => updates(message as OverrideAction))
          as OverrideAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OverrideAction create() => OverrideAction._();
  @$core.override
  OverrideAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OverrideAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OverrideAction>(create);
  static OverrideAction? _defaultInstance;

  @$pb.TagNumber(31963285)
  CountAction get count => $_getN(0);
  @$pb.TagNumber(31963285)
  set count(CountAction value) => $_setField(31963285, value);
  @$pb.TagNumber(31963285)
  $core.bool hasCount() => $_has(0);
  @$pb.TagNumber(31963285)
  void clearCount() => $_clearField(31963285);
  @$pb.TagNumber(31963285)
  CountAction ensureCount() => $_ensure(0);

  @$pb.TagNumber(273676284)
  NoneAction get none => $_getN(1);
  @$pb.TagNumber(273676284)
  set none(NoneAction value) => $_setField(273676284, value);
  @$pb.TagNumber(273676284)
  $core.bool hasNone() => $_has(1);
  @$pb.TagNumber(273676284)
  void clearNone() => $_clearField(273676284);
  @$pb.TagNumber(273676284)
  NoneAction ensureNone() => $_ensure(1);
}

class PasswordField extends $pb.GeneratedMessage {
  factory PasswordField({
    $core.String? identifier,
  }) {
    final result = create();
    if (identifier != null) result.identifier = identifier;
    return result;
  }

  PasswordField._();

  factory PasswordField.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PasswordField.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PasswordField',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(41865311, _omitFieldNames ? '' : 'identifier')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PasswordField clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PasswordField copyWith(void Function(PasswordField) updates) =>
      super.copyWith((message) => updates(message as PasswordField))
          as PasswordField;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PasswordField create() => PasswordField._();
  @$core.override
  PasswordField createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PasswordField getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PasswordField>(create);
  static PasswordField? _defaultInstance;

  @$pb.TagNumber(41865311)
  $core.String get identifier => $_getSZ(0);
  @$pb.TagNumber(41865311)
  set identifier($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41865311)
  $core.bool hasIdentifier() => $_has(0);
  @$pb.TagNumber(41865311)
  void clearIdentifier() => $_clearField(41865311);
}

class PathStatistics extends $pb.GeneratedMessage {
  factory PathStatistics({
    FilterSource? source,
    $core.String? path,
    $fixnum.Int64? requestcount,
    $core.double? percentage,
    $core.Iterable<BotStatistics>? topbots,
  }) {
    final result = create();
    if (source != null) result.source = source;
    if (path != null) result.path = path;
    if (requestcount != null) result.requestcount = requestcount;
    if (percentage != null) result.percentage = percentage;
    if (topbots != null) result.topbots.addAll(topbots);
    return result;
  }

  PathStatistics._();

  factory PathStatistics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PathStatistics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PathStatistics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<FilterSource>(31630329, _omitFieldNames ? '' : 'source',
        subBuilder: FilterSource.create)
    ..aOS(191292503, _omitFieldNames ? '' : 'path')
    ..aInt64(288519674, _omitFieldNames ? '' : 'requestcount')
    ..aD(341153238, _omitFieldNames ? '' : 'percentage')
    ..pPM<BotStatistics>(521235415, _omitFieldNames ? '' : 'topbots',
        subBuilder: BotStatistics.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PathStatistics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PathStatistics copyWith(void Function(PathStatistics) updates) =>
      super.copyWith((message) => updates(message as PathStatistics))
          as PathStatistics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PathStatistics create() => PathStatistics._();
  @$core.override
  PathStatistics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PathStatistics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PathStatistics>(create);
  static PathStatistics? _defaultInstance;

  @$pb.TagNumber(31630329)
  FilterSource get source => $_getN(0);
  @$pb.TagNumber(31630329)
  set source(FilterSource value) => $_setField(31630329, value);
  @$pb.TagNumber(31630329)
  $core.bool hasSource() => $_has(0);
  @$pb.TagNumber(31630329)
  void clearSource() => $_clearField(31630329);
  @$pb.TagNumber(31630329)
  FilterSource ensureSource() => $_ensure(0);

  @$pb.TagNumber(191292503)
  $core.String get path => $_getSZ(1);
  @$pb.TagNumber(191292503)
  set path($core.String value) => $_setString(1, value);
  @$pb.TagNumber(191292503)
  $core.bool hasPath() => $_has(1);
  @$pb.TagNumber(191292503)
  void clearPath() => $_clearField(191292503);

  @$pb.TagNumber(288519674)
  $fixnum.Int64 get requestcount => $_getI64(2);
  @$pb.TagNumber(288519674)
  set requestcount($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(288519674)
  $core.bool hasRequestcount() => $_has(2);
  @$pb.TagNumber(288519674)
  void clearRequestcount() => $_clearField(288519674);

  @$pb.TagNumber(341153238)
  $core.double get percentage => $_getN(3);
  @$pb.TagNumber(341153238)
  set percentage($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(341153238)
  $core.bool hasPercentage() => $_has(3);
  @$pb.TagNumber(341153238)
  void clearPercentage() => $_clearField(341153238);

  @$pb.TagNumber(521235415)
  $pb.PbList<BotStatistics> get topbots => $_getList(4);
}

class PhoneNumberField extends $pb.GeneratedMessage {
  factory PhoneNumberField({
    $core.String? identifier,
  }) {
    final result = create();
    if (identifier != null) result.identifier = identifier;
    return result;
  }

  PhoneNumberField._();

  factory PhoneNumberField.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PhoneNumberField.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PhoneNumberField',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(41865311, _omitFieldNames ? '' : 'identifier')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PhoneNumberField clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PhoneNumberField copyWith(void Function(PhoneNumberField) updates) =>
      super.copyWith((message) => updates(message as PhoneNumberField))
          as PhoneNumberField;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PhoneNumberField create() => PhoneNumberField._();
  @$core.override
  PhoneNumberField createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PhoneNumberField getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PhoneNumberField>(create);
  static PhoneNumberField? _defaultInstance;

  @$pb.TagNumber(41865311)
  $core.String get identifier => $_getSZ(0);
  @$pb.TagNumber(41865311)
  set identifier($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41865311)
  $core.bool hasIdentifier() => $_has(0);
  @$pb.TagNumber(41865311)
  void clearIdentifier() => $_clearField(41865311);
}

class PutLoggingConfigurationRequest extends $pb.GeneratedMessage {
  factory PutLoggingConfigurationRequest({
    LoggingConfiguration? loggingconfiguration,
  }) {
    final result = create();
    if (loggingconfiguration != null)
      result.loggingconfiguration = loggingconfiguration;
    return result;
  }

  PutLoggingConfigurationRequest._();

  factory PutLoggingConfigurationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutLoggingConfigurationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutLoggingConfigurationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<LoggingConfiguration>(
        359027765, _omitFieldNames ? '' : 'loggingconfiguration',
        subBuilder: LoggingConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutLoggingConfigurationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutLoggingConfigurationRequest copyWith(
          void Function(PutLoggingConfigurationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PutLoggingConfigurationRequest))
          as PutLoggingConfigurationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutLoggingConfigurationRequest create() =>
      PutLoggingConfigurationRequest._();
  @$core.override
  PutLoggingConfigurationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutLoggingConfigurationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutLoggingConfigurationRequest>(create);
  static PutLoggingConfigurationRequest? _defaultInstance;

  @$pb.TagNumber(359027765)
  LoggingConfiguration get loggingconfiguration => $_getN(0);
  @$pb.TagNumber(359027765)
  set loggingconfiguration(LoggingConfiguration value) =>
      $_setField(359027765, value);
  @$pb.TagNumber(359027765)
  $core.bool hasLoggingconfiguration() => $_has(0);
  @$pb.TagNumber(359027765)
  void clearLoggingconfiguration() => $_clearField(359027765);
  @$pb.TagNumber(359027765)
  LoggingConfiguration ensureLoggingconfiguration() => $_ensure(0);
}

class PutLoggingConfigurationResponse extends $pb.GeneratedMessage {
  factory PutLoggingConfigurationResponse({
    LoggingConfiguration? loggingconfiguration,
  }) {
    final result = create();
    if (loggingconfiguration != null)
      result.loggingconfiguration = loggingconfiguration;
    return result;
  }

  PutLoggingConfigurationResponse._();

  factory PutLoggingConfigurationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutLoggingConfigurationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutLoggingConfigurationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<LoggingConfiguration>(
        359027765, _omitFieldNames ? '' : 'loggingconfiguration',
        subBuilder: LoggingConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutLoggingConfigurationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutLoggingConfigurationResponse copyWith(
          void Function(PutLoggingConfigurationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PutLoggingConfigurationResponse))
          as PutLoggingConfigurationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutLoggingConfigurationResponse create() =>
      PutLoggingConfigurationResponse._();
  @$core.override
  PutLoggingConfigurationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutLoggingConfigurationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutLoggingConfigurationResponse>(
          create);
  static PutLoggingConfigurationResponse? _defaultInstance;

  @$pb.TagNumber(359027765)
  LoggingConfiguration get loggingconfiguration => $_getN(0);
  @$pb.TagNumber(359027765)
  set loggingconfiguration(LoggingConfiguration value) =>
      $_setField(359027765, value);
  @$pb.TagNumber(359027765)
  $core.bool hasLoggingconfiguration() => $_has(0);
  @$pb.TagNumber(359027765)
  void clearLoggingconfiguration() => $_clearField(359027765);
  @$pb.TagNumber(359027765)
  LoggingConfiguration ensureLoggingconfiguration() => $_ensure(0);
}

class PutManagedRuleSetVersionsRequest extends $pb.GeneratedMessage {
  factory PutManagedRuleSetVersionsRequest({
    $core.String? locktoken,
    $core.Iterable<$core.MapEntry<$core.String, VersionToPublish>>?
        versionstopublish,
    Scope? scope,
    $core.String? recommendedversion,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (versionstopublish != null)
      result.versionstopublish.addEntries(versionstopublish);
    if (scope != null) result.scope = scope;
    if (recommendedversion != null)
      result.recommendedversion = recommendedversion;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  PutManagedRuleSetVersionsRequest._();

  factory PutManagedRuleSetVersionsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutManagedRuleSetVersionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutManagedRuleSetVersionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..m<$core.String, VersionToPublish>(
        54645961, _omitFieldNames ? '' : 'versionstopublish',
        entryClassName:
            'PutManagedRuleSetVersionsRequest.VersionstopublishEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: VersionToPublish.create,
        valueDefaultOrMaker: VersionToPublish.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(83103137, _omitFieldNames ? '' : 'recommendedversion')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutManagedRuleSetVersionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutManagedRuleSetVersionsRequest copyWith(
          void Function(PutManagedRuleSetVersionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PutManagedRuleSetVersionsRequest))
          as PutManagedRuleSetVersionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutManagedRuleSetVersionsRequest create() =>
      PutManagedRuleSetVersionsRequest._();
  @$core.override
  PutManagedRuleSetVersionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutManagedRuleSetVersionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutManagedRuleSetVersionsRequest>(
          create);
  static PutManagedRuleSetVersionsRequest? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(54645961)
  $pb.PbMap<$core.String, VersionToPublish> get versionstopublish =>
      $_getMap(1);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(2);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(2);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(83103137)
  $core.String get recommendedversion => $_getSZ(3);
  @$pb.TagNumber(83103137)
  set recommendedversion($core.String value) => $_setString(3, value);
  @$pb.TagNumber(83103137)
  $core.bool hasRecommendedversion() => $_has(3);
  @$pb.TagNumber(83103137)
  void clearRecommendedversion() => $_clearField(83103137);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(5);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(5, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(5);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class PutManagedRuleSetVersionsResponse extends $pb.GeneratedMessage {
  factory PutManagedRuleSetVersionsResponse({
    $core.String? nextlocktoken,
  }) {
    final result = create();
    if (nextlocktoken != null) result.nextlocktoken = nextlocktoken;
    return result;
  }

  PutManagedRuleSetVersionsResponse._();

  factory PutManagedRuleSetVersionsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutManagedRuleSetVersionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutManagedRuleSetVersionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(473728437, _omitFieldNames ? '' : 'nextlocktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutManagedRuleSetVersionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutManagedRuleSetVersionsResponse copyWith(
          void Function(PutManagedRuleSetVersionsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as PutManagedRuleSetVersionsResponse))
          as PutManagedRuleSetVersionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutManagedRuleSetVersionsResponse create() =>
      PutManagedRuleSetVersionsResponse._();
  @$core.override
  PutManagedRuleSetVersionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutManagedRuleSetVersionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutManagedRuleSetVersionsResponse>(
          create);
  static PutManagedRuleSetVersionsResponse? _defaultInstance;

  @$pb.TagNumber(473728437)
  $core.String get nextlocktoken => $_getSZ(0);
  @$pb.TagNumber(473728437)
  set nextlocktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(473728437)
  $core.bool hasNextlocktoken() => $_has(0);
  @$pb.TagNumber(473728437)
  void clearNextlocktoken() => $_clearField(473728437);
}

class PutPermissionPolicyRequest extends $pb.GeneratedMessage {
  factory PutPermissionPolicyRequest({
    $core.String? resourcearn,
    $core.String? policy,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (policy != null) result.policy = policy;
    return result;
  }

  PutPermissionPolicyRequest._();

  factory PutPermissionPolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutPermissionPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutPermissionPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPermissionPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPermissionPolicyRequest copyWith(
          void Function(PutPermissionPolicyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as PutPermissionPolicyRequest))
          as PutPermissionPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutPermissionPolicyRequest create() => PutPermissionPolicyRequest._();
  @$core.override
  PutPermissionPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutPermissionPolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutPermissionPolicyRequest>(create);
  static PutPermissionPolicyRequest? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(1);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(1, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(1);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class PutPermissionPolicyResponse extends $pb.GeneratedMessage {
  factory PutPermissionPolicyResponse() => create();

  PutPermissionPolicyResponse._();

  factory PutPermissionPolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutPermissionPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutPermissionPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPermissionPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPermissionPolicyResponse copyWith(
          void Function(PutPermissionPolicyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as PutPermissionPolicyResponse))
          as PutPermissionPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutPermissionPolicyResponse create() =>
      PutPermissionPolicyResponse._();
  @$core.override
  PutPermissionPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutPermissionPolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutPermissionPolicyResponse>(create);
  static PutPermissionPolicyResponse? _defaultInstance;
}

class QueryString extends $pb.GeneratedMessage {
  factory QueryString() => create();

  QueryString._();

  factory QueryString.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryString.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryString',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryString clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryString copyWith(void Function(QueryString) updates) =>
      super.copyWith((message) => updates(message as QueryString))
          as QueryString;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryString create() => QueryString._();
  @$core.override
  QueryString createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryString getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryString>(create);
  static QueryString? _defaultInstance;
}

class RateBasedStatement extends $pb.GeneratedMessage {
  factory RateBasedStatement({
    Statement? scopedownstatement,
    RateBasedStatementAggregateKeyType? aggregatekeytype,
    ForwardedIPConfig? forwardedipconfig,
    $fixnum.Int64? evaluationwindowsec,
    $fixnum.Int64? limit,
    $core.Iterable<RateBasedStatementCustomKey>? customkeys,
  }) {
    final result = create();
    if (scopedownstatement != null)
      result.scopedownstatement = scopedownstatement;
    if (aggregatekeytype != null) result.aggregatekeytype = aggregatekeytype;
    if (forwardedipconfig != null) result.forwardedipconfig = forwardedipconfig;
    if (evaluationwindowsec != null)
      result.evaluationwindowsec = evaluationwindowsec;
    if (limit != null) result.limit = limit;
    if (customkeys != null) result.customkeys.addAll(customkeys);
    return result;
  }

  RateBasedStatement._();

  factory RateBasedStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateBasedStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateBasedStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<Statement>(116108605, _omitFieldNames ? '' : 'scopedownstatement',
        subBuilder: Statement.create)
    ..aE<RateBasedStatementAggregateKeyType>(
        200283184, _omitFieldNames ? '' : 'aggregatekeytype',
        enumValues: RateBasedStatementAggregateKeyType.values)
    ..aOM<ForwardedIPConfig>(
        259846797, _omitFieldNames ? '' : 'forwardedipconfig',
        subBuilder: ForwardedIPConfig.create)
    ..aInt64(333939569, _omitFieldNames ? '' : 'evaluationwindowsec')
    ..aInt64(412502741, _omitFieldNames ? '' : 'limit')
    ..pPM<RateBasedStatementCustomKey>(
        429675987, _omitFieldNames ? '' : 'customkeys',
        subBuilder: RateBasedStatementCustomKey.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateBasedStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateBasedStatement copyWith(void Function(RateBasedStatement) updates) =>
      super.copyWith((message) => updates(message as RateBasedStatement))
          as RateBasedStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateBasedStatement create() => RateBasedStatement._();
  @$core.override
  RateBasedStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateBasedStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateBasedStatement>(create);
  static RateBasedStatement? _defaultInstance;

  @$pb.TagNumber(116108605)
  Statement get scopedownstatement => $_getN(0);
  @$pb.TagNumber(116108605)
  set scopedownstatement(Statement value) => $_setField(116108605, value);
  @$pb.TagNumber(116108605)
  $core.bool hasScopedownstatement() => $_has(0);
  @$pb.TagNumber(116108605)
  void clearScopedownstatement() => $_clearField(116108605);
  @$pb.TagNumber(116108605)
  Statement ensureScopedownstatement() => $_ensure(0);

  @$pb.TagNumber(200283184)
  RateBasedStatementAggregateKeyType get aggregatekeytype => $_getN(1);
  @$pb.TagNumber(200283184)
  set aggregatekeytype(RateBasedStatementAggregateKeyType value) =>
      $_setField(200283184, value);
  @$pb.TagNumber(200283184)
  $core.bool hasAggregatekeytype() => $_has(1);
  @$pb.TagNumber(200283184)
  void clearAggregatekeytype() => $_clearField(200283184);

  @$pb.TagNumber(259846797)
  ForwardedIPConfig get forwardedipconfig => $_getN(2);
  @$pb.TagNumber(259846797)
  set forwardedipconfig(ForwardedIPConfig value) =>
      $_setField(259846797, value);
  @$pb.TagNumber(259846797)
  $core.bool hasForwardedipconfig() => $_has(2);
  @$pb.TagNumber(259846797)
  void clearForwardedipconfig() => $_clearField(259846797);
  @$pb.TagNumber(259846797)
  ForwardedIPConfig ensureForwardedipconfig() => $_ensure(2);

  @$pb.TagNumber(333939569)
  $fixnum.Int64 get evaluationwindowsec => $_getI64(3);
  @$pb.TagNumber(333939569)
  set evaluationwindowsec($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(333939569)
  $core.bool hasEvaluationwindowsec() => $_has(3);
  @$pb.TagNumber(333939569)
  void clearEvaluationwindowsec() => $_clearField(333939569);

  @$pb.TagNumber(412502741)
  $fixnum.Int64 get limit => $_getI64(4);
  @$pb.TagNumber(412502741)
  set limit($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(4);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(429675987)
  $pb.PbList<RateBasedStatementCustomKey> get customkeys => $_getList(5);
}

class RateBasedStatementCustomKey extends $pb.GeneratedMessage {
  factory RateBasedStatementCustomKey({
    RateLimitLabelNamespace? labelnamespace,
    RateLimitIP? ip,
    RateLimitQueryArgument? queryargument,
    RateLimitUriPath? uripath,
    RateLimitHeader? header,
    RateLimitJA4Fingerprint? ja4fingerprint,
    RateLimitHTTPMethod? httpmethod,
    RateLimitAsn? asn,
    RateLimitQueryString? querystring,
    RateLimitForwardedIP? forwardedip,
    RateLimitCookie? cookie,
    RateLimitJA3Fingerprint? ja3fingerprint,
  }) {
    final result = create();
    if (labelnamespace != null) result.labelnamespace = labelnamespace;
    if (ip != null) result.ip = ip;
    if (queryargument != null) result.queryargument = queryargument;
    if (uripath != null) result.uripath = uripath;
    if (header != null) result.header = header;
    if (ja4fingerprint != null) result.ja4fingerprint = ja4fingerprint;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (asn != null) result.asn = asn;
    if (querystring != null) result.querystring = querystring;
    if (forwardedip != null) result.forwardedip = forwardedip;
    if (cookie != null) result.cookie = cookie;
    if (ja3fingerprint != null) result.ja3fingerprint = ja3fingerprint;
    return result;
  }

  RateBasedStatementCustomKey._();

  factory RateBasedStatementCustomKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateBasedStatementCustomKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateBasedStatementCustomKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<RateLimitLabelNamespace>(
        39797417, _omitFieldNames ? '' : 'labelnamespace',
        subBuilder: RateLimitLabelNamespace.create)
    ..aOM<RateLimitIP>(183044829, _omitFieldNames ? '' : 'ip',
        subBuilder: RateLimitIP.create)
    ..aOM<RateLimitQueryArgument>(
        277971873, _omitFieldNames ? '' : 'queryargument',
        subBuilder: RateLimitQueryArgument.create)
    ..aOM<RateLimitUriPath>(288340351, _omitFieldNames ? '' : 'uripath',
        subBuilder: RateLimitUriPath.create)
    ..aOM<RateLimitHeader>(290429313, _omitFieldNames ? '' : 'header',
        subBuilder: RateLimitHeader.create)
    ..aOM<RateLimitJA4Fingerprint>(
        352209397, _omitFieldNames ? '' : 'ja4fingerprint',
        subBuilder: RateLimitJA4Fingerprint.create)
    ..aOM<RateLimitHTTPMethod>(368543473, _omitFieldNames ? '' : 'httpmethod',
        subBuilder: RateLimitHTTPMethod.create)
    ..aOM<RateLimitAsn>(430799034, _omitFieldNames ? '' : 'asn',
        subBuilder: RateLimitAsn.create)
    ..aOM<RateLimitQueryString>(435938663, _omitFieldNames ? '' : 'querystring',
        subBuilder: RateLimitQueryString.create)
    ..aOM<RateLimitForwardedIP>(478980035, _omitFieldNames ? '' : 'forwardedip',
        subBuilder: RateLimitForwardedIP.create)
    ..aOM<RateLimitCookie>(498776800, _omitFieldNames ? '' : 'cookie',
        subBuilder: RateLimitCookie.create)
    ..aOM<RateLimitJA3Fingerprint>(
        521270332, _omitFieldNames ? '' : 'ja3fingerprint',
        subBuilder: RateLimitJA3Fingerprint.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateBasedStatementCustomKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateBasedStatementCustomKey copyWith(
          void Function(RateBasedStatementCustomKey) updates) =>
      super.copyWith(
              (message) => updates(message as RateBasedStatementCustomKey))
          as RateBasedStatementCustomKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateBasedStatementCustomKey create() =>
      RateBasedStatementCustomKey._();
  @$core.override
  RateBasedStatementCustomKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateBasedStatementCustomKey getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateBasedStatementCustomKey>(create);
  static RateBasedStatementCustomKey? _defaultInstance;

  @$pb.TagNumber(39797417)
  RateLimitLabelNamespace get labelnamespace => $_getN(0);
  @$pb.TagNumber(39797417)
  set labelnamespace(RateLimitLabelNamespace value) =>
      $_setField(39797417, value);
  @$pb.TagNumber(39797417)
  $core.bool hasLabelnamespace() => $_has(0);
  @$pb.TagNumber(39797417)
  void clearLabelnamespace() => $_clearField(39797417);
  @$pb.TagNumber(39797417)
  RateLimitLabelNamespace ensureLabelnamespace() => $_ensure(0);

  @$pb.TagNumber(183044829)
  RateLimitIP get ip => $_getN(1);
  @$pb.TagNumber(183044829)
  set ip(RateLimitIP value) => $_setField(183044829, value);
  @$pb.TagNumber(183044829)
  $core.bool hasIp() => $_has(1);
  @$pb.TagNumber(183044829)
  void clearIp() => $_clearField(183044829);
  @$pb.TagNumber(183044829)
  RateLimitIP ensureIp() => $_ensure(1);

  @$pb.TagNumber(277971873)
  RateLimitQueryArgument get queryargument => $_getN(2);
  @$pb.TagNumber(277971873)
  set queryargument(RateLimitQueryArgument value) =>
      $_setField(277971873, value);
  @$pb.TagNumber(277971873)
  $core.bool hasQueryargument() => $_has(2);
  @$pb.TagNumber(277971873)
  void clearQueryargument() => $_clearField(277971873);
  @$pb.TagNumber(277971873)
  RateLimitQueryArgument ensureQueryargument() => $_ensure(2);

  @$pb.TagNumber(288340351)
  RateLimitUriPath get uripath => $_getN(3);
  @$pb.TagNumber(288340351)
  set uripath(RateLimitUriPath value) => $_setField(288340351, value);
  @$pb.TagNumber(288340351)
  $core.bool hasUripath() => $_has(3);
  @$pb.TagNumber(288340351)
  void clearUripath() => $_clearField(288340351);
  @$pb.TagNumber(288340351)
  RateLimitUriPath ensureUripath() => $_ensure(3);

  @$pb.TagNumber(290429313)
  RateLimitHeader get header => $_getN(4);
  @$pb.TagNumber(290429313)
  set header(RateLimitHeader value) => $_setField(290429313, value);
  @$pb.TagNumber(290429313)
  $core.bool hasHeader() => $_has(4);
  @$pb.TagNumber(290429313)
  void clearHeader() => $_clearField(290429313);
  @$pb.TagNumber(290429313)
  RateLimitHeader ensureHeader() => $_ensure(4);

  @$pb.TagNumber(352209397)
  RateLimitJA4Fingerprint get ja4fingerprint => $_getN(5);
  @$pb.TagNumber(352209397)
  set ja4fingerprint(RateLimitJA4Fingerprint value) =>
      $_setField(352209397, value);
  @$pb.TagNumber(352209397)
  $core.bool hasJa4fingerprint() => $_has(5);
  @$pb.TagNumber(352209397)
  void clearJa4fingerprint() => $_clearField(352209397);
  @$pb.TagNumber(352209397)
  RateLimitJA4Fingerprint ensureJa4fingerprint() => $_ensure(5);

  @$pb.TagNumber(368543473)
  RateLimitHTTPMethod get httpmethod => $_getN(6);
  @$pb.TagNumber(368543473)
  set httpmethod(RateLimitHTTPMethod value) => $_setField(368543473, value);
  @$pb.TagNumber(368543473)
  $core.bool hasHttpmethod() => $_has(6);
  @$pb.TagNumber(368543473)
  void clearHttpmethod() => $_clearField(368543473);
  @$pb.TagNumber(368543473)
  RateLimitHTTPMethod ensureHttpmethod() => $_ensure(6);

  @$pb.TagNumber(430799034)
  RateLimitAsn get asn => $_getN(7);
  @$pb.TagNumber(430799034)
  set asn(RateLimitAsn value) => $_setField(430799034, value);
  @$pb.TagNumber(430799034)
  $core.bool hasAsn() => $_has(7);
  @$pb.TagNumber(430799034)
  void clearAsn() => $_clearField(430799034);
  @$pb.TagNumber(430799034)
  RateLimitAsn ensureAsn() => $_ensure(7);

  @$pb.TagNumber(435938663)
  RateLimitQueryString get querystring => $_getN(8);
  @$pb.TagNumber(435938663)
  set querystring(RateLimitQueryString value) => $_setField(435938663, value);
  @$pb.TagNumber(435938663)
  $core.bool hasQuerystring() => $_has(8);
  @$pb.TagNumber(435938663)
  void clearQuerystring() => $_clearField(435938663);
  @$pb.TagNumber(435938663)
  RateLimitQueryString ensureQuerystring() => $_ensure(8);

  @$pb.TagNumber(478980035)
  RateLimitForwardedIP get forwardedip => $_getN(9);
  @$pb.TagNumber(478980035)
  set forwardedip(RateLimitForwardedIP value) => $_setField(478980035, value);
  @$pb.TagNumber(478980035)
  $core.bool hasForwardedip() => $_has(9);
  @$pb.TagNumber(478980035)
  void clearForwardedip() => $_clearField(478980035);
  @$pb.TagNumber(478980035)
  RateLimitForwardedIP ensureForwardedip() => $_ensure(9);

  @$pb.TagNumber(498776800)
  RateLimitCookie get cookie => $_getN(10);
  @$pb.TagNumber(498776800)
  set cookie(RateLimitCookie value) => $_setField(498776800, value);
  @$pb.TagNumber(498776800)
  $core.bool hasCookie() => $_has(10);
  @$pb.TagNumber(498776800)
  void clearCookie() => $_clearField(498776800);
  @$pb.TagNumber(498776800)
  RateLimitCookie ensureCookie() => $_ensure(10);

  @$pb.TagNumber(521270332)
  RateLimitJA3Fingerprint get ja3fingerprint => $_getN(11);
  @$pb.TagNumber(521270332)
  set ja3fingerprint(RateLimitJA3Fingerprint value) =>
      $_setField(521270332, value);
  @$pb.TagNumber(521270332)
  $core.bool hasJa3fingerprint() => $_has(11);
  @$pb.TagNumber(521270332)
  void clearJa3fingerprint() => $_clearField(521270332);
  @$pb.TagNumber(521270332)
  RateLimitJA3Fingerprint ensureJa3fingerprint() => $_ensure(11);
}

class RateBasedStatementManagedKeysIPSet extends $pb.GeneratedMessage {
  factory RateBasedStatementManagedKeysIPSet({
    IPAddressVersion? ipaddressversion,
    $core.Iterable<$core.String>? addresses,
  }) {
    final result = create();
    if (ipaddressversion != null) result.ipaddressversion = ipaddressversion;
    if (addresses != null) result.addresses.addAll(addresses);
    return result;
  }

  RateBasedStatementManagedKeysIPSet._();

  factory RateBasedStatementManagedKeysIPSet.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateBasedStatementManagedKeysIPSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateBasedStatementManagedKeysIPSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<IPAddressVersion>(313363841, _omitFieldNames ? '' : 'ipaddressversion',
        enumValues: IPAddressVersion.values)
    ..pPS(375939972, _omitFieldNames ? '' : 'addresses')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateBasedStatementManagedKeysIPSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateBasedStatementManagedKeysIPSet copyWith(
          void Function(RateBasedStatementManagedKeysIPSet) updates) =>
      super.copyWith((message) =>
              updates(message as RateBasedStatementManagedKeysIPSet))
          as RateBasedStatementManagedKeysIPSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateBasedStatementManagedKeysIPSet create() =>
      RateBasedStatementManagedKeysIPSet._();
  @$core.override
  RateBasedStatementManagedKeysIPSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateBasedStatementManagedKeysIPSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateBasedStatementManagedKeysIPSet>(
          create);
  static RateBasedStatementManagedKeysIPSet? _defaultInstance;

  @$pb.TagNumber(313363841)
  IPAddressVersion get ipaddressversion => $_getN(0);
  @$pb.TagNumber(313363841)
  set ipaddressversion(IPAddressVersion value) => $_setField(313363841, value);
  @$pb.TagNumber(313363841)
  $core.bool hasIpaddressversion() => $_has(0);
  @$pb.TagNumber(313363841)
  void clearIpaddressversion() => $_clearField(313363841);

  @$pb.TagNumber(375939972)
  $pb.PbList<$core.String> get addresses => $_getList(1);
}

class RateLimitAsn extends $pb.GeneratedMessage {
  factory RateLimitAsn() => create();

  RateLimitAsn._();

  factory RateLimitAsn.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitAsn.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitAsn',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitAsn clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitAsn copyWith(void Function(RateLimitAsn) updates) =>
      super.copyWith((message) => updates(message as RateLimitAsn))
          as RateLimitAsn;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitAsn create() => RateLimitAsn._();
  @$core.override
  RateLimitAsn createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitAsn getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitAsn>(create);
  static RateLimitAsn? _defaultInstance;
}

class RateLimitCookie extends $pb.GeneratedMessage {
  factory RateLimitCookie({
    $core.Iterable<TextTransformation>? texttransformations,
    $core.String? name,
  }) {
    final result = create();
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (name != null) result.name = name;
    return result;
  }

  RateLimitCookie._();

  factory RateLimitCookie.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitCookie.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitCookie',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitCookie clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitCookie copyWith(void Function(RateLimitCookie) updates) =>
      super.copyWith((message) => updates(message as RateLimitCookie))
          as RateLimitCookie;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitCookie create() => RateLimitCookie._();
  @$core.override
  RateLimitCookie createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitCookie getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitCookie>(create);
  static RateLimitCookie? _defaultInstance;

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class RateLimitForwardedIP extends $pb.GeneratedMessage {
  factory RateLimitForwardedIP() => create();

  RateLimitForwardedIP._();

  factory RateLimitForwardedIP.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitForwardedIP.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitForwardedIP',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitForwardedIP clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitForwardedIP copyWith(void Function(RateLimitForwardedIP) updates) =>
      super.copyWith((message) => updates(message as RateLimitForwardedIP))
          as RateLimitForwardedIP;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitForwardedIP create() => RateLimitForwardedIP._();
  @$core.override
  RateLimitForwardedIP createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitForwardedIP getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitForwardedIP>(create);
  static RateLimitForwardedIP? _defaultInstance;
}

class RateLimitHTTPMethod extends $pb.GeneratedMessage {
  factory RateLimitHTTPMethod() => create();

  RateLimitHTTPMethod._();

  factory RateLimitHTTPMethod.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitHTTPMethod.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitHTTPMethod',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitHTTPMethod clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitHTTPMethod copyWith(void Function(RateLimitHTTPMethod) updates) =>
      super.copyWith((message) => updates(message as RateLimitHTTPMethod))
          as RateLimitHTTPMethod;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitHTTPMethod create() => RateLimitHTTPMethod._();
  @$core.override
  RateLimitHTTPMethod createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitHTTPMethod getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitHTTPMethod>(create);
  static RateLimitHTTPMethod? _defaultInstance;
}

class RateLimitHeader extends $pb.GeneratedMessage {
  factory RateLimitHeader({
    $core.Iterable<TextTransformation>? texttransformations,
    $core.String? name,
  }) {
    final result = create();
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (name != null) result.name = name;
    return result;
  }

  RateLimitHeader._();

  factory RateLimitHeader.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitHeader.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitHeader',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitHeader clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitHeader copyWith(void Function(RateLimitHeader) updates) =>
      super.copyWith((message) => updates(message as RateLimitHeader))
          as RateLimitHeader;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitHeader create() => RateLimitHeader._();
  @$core.override
  RateLimitHeader createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitHeader getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitHeader>(create);
  static RateLimitHeader? _defaultInstance;

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class RateLimitIP extends $pb.GeneratedMessage {
  factory RateLimitIP() => create();

  RateLimitIP._();

  factory RateLimitIP.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitIP.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitIP',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitIP clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitIP copyWith(void Function(RateLimitIP) updates) =>
      super.copyWith((message) => updates(message as RateLimitIP))
          as RateLimitIP;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitIP create() => RateLimitIP._();
  @$core.override
  RateLimitIP createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitIP getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitIP>(create);
  static RateLimitIP? _defaultInstance;
}

class RateLimitJA3Fingerprint extends $pb.GeneratedMessage {
  factory RateLimitJA3Fingerprint({
    FallbackBehavior? fallbackbehavior,
  }) {
    final result = create();
    if (fallbackbehavior != null) result.fallbackbehavior = fallbackbehavior;
    return result;
  }

  RateLimitJA3Fingerprint._();

  factory RateLimitJA3Fingerprint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitJA3Fingerprint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitJA3Fingerprint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<FallbackBehavior>(440114542, _omitFieldNames ? '' : 'fallbackbehavior',
        enumValues: FallbackBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitJA3Fingerprint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitJA3Fingerprint copyWith(
          void Function(RateLimitJA3Fingerprint) updates) =>
      super.copyWith((message) => updates(message as RateLimitJA3Fingerprint))
          as RateLimitJA3Fingerprint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitJA3Fingerprint create() => RateLimitJA3Fingerprint._();
  @$core.override
  RateLimitJA3Fingerprint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitJA3Fingerprint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitJA3Fingerprint>(create);
  static RateLimitJA3Fingerprint? _defaultInstance;

  @$pb.TagNumber(440114542)
  FallbackBehavior get fallbackbehavior => $_getN(0);
  @$pb.TagNumber(440114542)
  set fallbackbehavior(FallbackBehavior value) => $_setField(440114542, value);
  @$pb.TagNumber(440114542)
  $core.bool hasFallbackbehavior() => $_has(0);
  @$pb.TagNumber(440114542)
  void clearFallbackbehavior() => $_clearField(440114542);
}

class RateLimitJA4Fingerprint extends $pb.GeneratedMessage {
  factory RateLimitJA4Fingerprint({
    FallbackBehavior? fallbackbehavior,
  }) {
    final result = create();
    if (fallbackbehavior != null) result.fallbackbehavior = fallbackbehavior;
    return result;
  }

  RateLimitJA4Fingerprint._();

  factory RateLimitJA4Fingerprint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitJA4Fingerprint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitJA4Fingerprint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<FallbackBehavior>(440114542, _omitFieldNames ? '' : 'fallbackbehavior',
        enumValues: FallbackBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitJA4Fingerprint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitJA4Fingerprint copyWith(
          void Function(RateLimitJA4Fingerprint) updates) =>
      super.copyWith((message) => updates(message as RateLimitJA4Fingerprint))
          as RateLimitJA4Fingerprint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitJA4Fingerprint create() => RateLimitJA4Fingerprint._();
  @$core.override
  RateLimitJA4Fingerprint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitJA4Fingerprint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitJA4Fingerprint>(create);
  static RateLimitJA4Fingerprint? _defaultInstance;

  @$pb.TagNumber(440114542)
  FallbackBehavior get fallbackbehavior => $_getN(0);
  @$pb.TagNumber(440114542)
  set fallbackbehavior(FallbackBehavior value) => $_setField(440114542, value);
  @$pb.TagNumber(440114542)
  $core.bool hasFallbackbehavior() => $_has(0);
  @$pb.TagNumber(440114542)
  void clearFallbackbehavior() => $_clearField(440114542);
}

class RateLimitLabelNamespace extends $pb.GeneratedMessage {
  factory RateLimitLabelNamespace({
    $core.String? namespace,
  }) {
    final result = create();
    if (namespace != null) result.namespace = namespace;
    return result;
  }

  RateLimitLabelNamespace._();

  factory RateLimitLabelNamespace.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitLabelNamespace.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitLabelNamespace',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitLabelNamespace clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitLabelNamespace copyWith(
          void Function(RateLimitLabelNamespace) updates) =>
      super.copyWith((message) => updates(message as RateLimitLabelNamespace))
          as RateLimitLabelNamespace;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitLabelNamespace create() => RateLimitLabelNamespace._();
  @$core.override
  RateLimitLabelNamespace createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitLabelNamespace getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitLabelNamespace>(create);
  static RateLimitLabelNamespace? _defaultInstance;

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(0);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(0, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(0);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);
}

class RateLimitQueryArgument extends $pb.GeneratedMessage {
  factory RateLimitQueryArgument({
    $core.Iterable<TextTransformation>? texttransformations,
    $core.String? name,
  }) {
    final result = create();
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (name != null) result.name = name;
    return result;
  }

  RateLimitQueryArgument._();

  factory RateLimitQueryArgument.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitQueryArgument.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitQueryArgument',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitQueryArgument clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitQueryArgument copyWith(
          void Function(RateLimitQueryArgument) updates) =>
      super.copyWith((message) => updates(message as RateLimitQueryArgument))
          as RateLimitQueryArgument;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitQueryArgument create() => RateLimitQueryArgument._();
  @$core.override
  RateLimitQueryArgument createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitQueryArgument getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitQueryArgument>(create);
  static RateLimitQueryArgument? _defaultInstance;

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class RateLimitQueryString extends $pb.GeneratedMessage {
  factory RateLimitQueryString({
    $core.Iterable<TextTransformation>? texttransformations,
  }) {
    final result = create();
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    return result;
  }

  RateLimitQueryString._();

  factory RateLimitQueryString.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitQueryString.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitQueryString',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitQueryString clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitQueryString copyWith(void Function(RateLimitQueryString) updates) =>
      super.copyWith((message) => updates(message as RateLimitQueryString))
          as RateLimitQueryString;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitQueryString create() => RateLimitQueryString._();
  @$core.override
  RateLimitQueryString createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitQueryString getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitQueryString>(create);
  static RateLimitQueryString? _defaultInstance;

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(0);
}

class RateLimitUriPath extends $pb.GeneratedMessage {
  factory RateLimitUriPath({
    $core.Iterable<TextTransformation>? texttransformations,
  }) {
    final result = create();
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    return result;
  }

  RateLimitUriPath._();

  factory RateLimitUriPath.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateLimitUriPath.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateLimitUriPath',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitUriPath clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateLimitUriPath copyWith(void Function(RateLimitUriPath) updates) =>
      super.copyWith((message) => updates(message as RateLimitUriPath))
          as RateLimitUriPath;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateLimitUriPath create() => RateLimitUriPath._();
  @$core.override
  RateLimitUriPath createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateLimitUriPath getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateLimitUriPath>(create);
  static RateLimitUriPath? _defaultInstance;

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(0);
}

class Regex extends $pb.GeneratedMessage {
  factory Regex({
    $core.String? regexstring,
  }) {
    final result = create();
    if (regexstring != null) result.regexstring = regexstring;
    return result;
  }

  Regex._();

  factory Regex.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Regex.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Regex',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(138533058, _omitFieldNames ? '' : 'regexstring')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Regex clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Regex copyWith(void Function(Regex) updates) =>
      super.copyWith((message) => updates(message as Regex)) as Regex;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Regex create() => Regex._();
  @$core.override
  Regex createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Regex getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Regex>(create);
  static Regex? _defaultInstance;

  @$pb.TagNumber(138533058)
  $core.String get regexstring => $_getSZ(0);
  @$pb.TagNumber(138533058)
  set regexstring($core.String value) => $_setString(0, value);
  @$pb.TagNumber(138533058)
  $core.bool hasRegexstring() => $_has(0);
  @$pb.TagNumber(138533058)
  void clearRegexstring() => $_clearField(138533058);
}

class RegexMatchStatement extends $pb.GeneratedMessage {
  factory RegexMatchStatement({
    $core.String? regexstring,
    $core.Iterable<TextTransformation>? texttransformations,
    FieldToMatch? fieldtomatch,
  }) {
    final result = create();
    if (regexstring != null) result.regexstring = regexstring;
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    return result;
  }

  RegexMatchStatement._();

  factory RegexMatchStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexMatchStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexMatchStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(138533058, _omitFieldNames ? '' : 'regexstring')
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchStatement copyWith(void Function(RegexMatchStatement) updates) =>
      super.copyWith((message) => updates(message as RegexMatchStatement))
          as RegexMatchStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexMatchStatement create() => RegexMatchStatement._();
  @$core.override
  RegexMatchStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexMatchStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexMatchStatement>(create);
  static RegexMatchStatement? _defaultInstance;

  @$pb.TagNumber(138533058)
  $core.String get regexstring => $_getSZ(0);
  @$pb.TagNumber(138533058)
  set regexstring($core.String value) => $_setString(0, value);
  @$pb.TagNumber(138533058)
  $core.bool hasRegexstring() => $_has(0);
  @$pb.TagNumber(138533058)
  void clearRegexstring() => $_clearField(138533058);

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(1);

  @$pb.TagNumber(338348372)
  FieldToMatch get fieldtomatch => $_getN(2);
  @$pb.TagNumber(338348372)
  set fieldtomatch(FieldToMatch value) => $_setField(338348372, value);
  @$pb.TagNumber(338348372)
  $core.bool hasFieldtomatch() => $_has(2);
  @$pb.TagNumber(338348372)
  void clearFieldtomatch() => $_clearField(338348372);
  @$pb.TagNumber(338348372)
  FieldToMatch ensureFieldtomatch() => $_ensure(2);
}

class RegexPatternSet extends $pb.GeneratedMessage {
  factory RegexPatternSet({
    $core.String? description,
    $core.Iterable<Regex>? regularexpressionlist,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (regularexpressionlist != null)
      result.regularexpressionlist.addAll(regularexpressionlist);
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  RegexPatternSet._();

  factory RegexPatternSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexPatternSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexPatternSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..pPM<Regex>(123612838, _omitFieldNames ? '' : 'regularexpressionlist',
        subBuilder: Regex.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexPatternSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexPatternSet copyWith(void Function(RegexPatternSet) updates) =>
      super.copyWith((message) => updates(message as RegexPatternSet))
          as RegexPatternSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexPatternSet create() => RegexPatternSet._();
  @$core.override
  RegexPatternSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexPatternSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexPatternSet>(create);
  static RegexPatternSet? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(123612838)
  $pb.PbList<Regex> get regularexpressionlist => $_getList(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class RegexPatternSetReferenceStatement extends $pb.GeneratedMessage {
  factory RegexPatternSetReferenceStatement({
    $core.Iterable<TextTransformation>? texttransformations,
    FieldToMatch? fieldtomatch,
    $core.String? arn,
  }) {
    final result = create();
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    if (arn != null) result.arn = arn;
    return result;
  }

  RegexPatternSetReferenceStatement._();

  factory RegexPatternSetReferenceStatement.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexPatternSetReferenceStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexPatternSetReferenceStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexPatternSetReferenceStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexPatternSetReferenceStatement copyWith(
          void Function(RegexPatternSetReferenceStatement) updates) =>
      super.copyWith((message) =>
              updates(message as RegexPatternSetReferenceStatement))
          as RegexPatternSetReferenceStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexPatternSetReferenceStatement create() =>
      RegexPatternSetReferenceStatement._();
  @$core.override
  RegexPatternSetReferenceStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexPatternSetReferenceStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexPatternSetReferenceStatement>(
          create);
  static RegexPatternSetReferenceStatement? _defaultInstance;

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(0);

  @$pb.TagNumber(338348372)
  FieldToMatch get fieldtomatch => $_getN(1);
  @$pb.TagNumber(338348372)
  set fieldtomatch(FieldToMatch value) => $_setField(338348372, value);
  @$pb.TagNumber(338348372)
  $core.bool hasFieldtomatch() => $_has(1);
  @$pb.TagNumber(338348372)
  void clearFieldtomatch() => $_clearField(338348372);
  @$pb.TagNumber(338348372)
  FieldToMatch ensureFieldtomatch() => $_ensure(1);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class RegexPatternSetSummary extends $pb.GeneratedMessage {
  factory RegexPatternSetSummary({
    $core.String? locktoken,
    $core.String? description,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  RegexPatternSetSummary._();

  factory RegexPatternSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexPatternSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexPatternSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexPatternSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexPatternSetSummary copyWith(
          void Function(RegexPatternSetSummary) updates) =>
      super.copyWith((message) => updates(message as RegexPatternSetSummary))
          as RegexPatternSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexPatternSetSummary create() => RegexPatternSetSummary._();
  @$core.override
  RegexPatternSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexPatternSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexPatternSetSummary>(create);
  static RegexPatternSetSummary? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class ReleaseSummary extends $pb.GeneratedMessage {
  factory ReleaseSummary({
    $core.String? releaseversion,
    $core.String? timestamp,
  }) {
    final result = create();
    if (releaseversion != null) result.releaseversion = releaseversion;
    if (timestamp != null) result.timestamp = timestamp;
    return result;
  }

  ReleaseSummary._();

  factory ReleaseSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReleaseSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReleaseSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(8739171, _omitFieldNames ? '' : 'releaseversion')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReleaseSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReleaseSummary copyWith(void Function(ReleaseSummary) updates) =>
      super.copyWith((message) => updates(message as ReleaseSummary))
          as ReleaseSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReleaseSummary create() => ReleaseSummary._();
  @$core.override
  ReleaseSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReleaseSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReleaseSummary>(create);
  static ReleaseSummary? _defaultInstance;

  @$pb.TagNumber(8739171)
  $core.String get releaseversion => $_getSZ(0);
  @$pb.TagNumber(8739171)
  set releaseversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(8739171)
  $core.bool hasReleaseversion() => $_has(0);
  @$pb.TagNumber(8739171)
  void clearReleaseversion() => $_clearField(8739171);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(1);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(1);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);
}

class RequestBodyAssociatedResourceTypeConfig extends $pb.GeneratedMessage {
  factory RequestBodyAssociatedResourceTypeConfig({
    SizeInspectionLimit? defaultsizeinspectionlimit,
  }) {
    final result = create();
    if (defaultsizeinspectionlimit != null)
      result.defaultsizeinspectionlimit = defaultsizeinspectionlimit;
    return result;
  }

  RequestBodyAssociatedResourceTypeConfig._();

  factory RequestBodyAssociatedResourceTypeConfig.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestBodyAssociatedResourceTypeConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestBodyAssociatedResourceTypeConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<SizeInspectionLimit>(
        344217929, _omitFieldNames ? '' : 'defaultsizeinspectionlimit',
        enumValues: SizeInspectionLimit.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestBodyAssociatedResourceTypeConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestBodyAssociatedResourceTypeConfig copyWith(
          void Function(RequestBodyAssociatedResourceTypeConfig) updates) =>
      super.copyWith((message) =>
              updates(message as RequestBodyAssociatedResourceTypeConfig))
          as RequestBodyAssociatedResourceTypeConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestBodyAssociatedResourceTypeConfig create() =>
      RequestBodyAssociatedResourceTypeConfig._();
  @$core.override
  RequestBodyAssociatedResourceTypeConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestBodyAssociatedResourceTypeConfig getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          RequestBodyAssociatedResourceTypeConfig>(create);
  static RequestBodyAssociatedResourceTypeConfig? _defaultInstance;

  @$pb.TagNumber(344217929)
  SizeInspectionLimit get defaultsizeinspectionlimit => $_getN(0);
  @$pb.TagNumber(344217929)
  set defaultsizeinspectionlimit(SizeInspectionLimit value) =>
      $_setField(344217929, value);
  @$pb.TagNumber(344217929)
  $core.bool hasDefaultsizeinspectionlimit() => $_has(0);
  @$pb.TagNumber(344217929)
  void clearDefaultsizeinspectionlimit() => $_clearField(344217929);
}

class RequestInspection extends $pb.GeneratedMessage {
  factory RequestInspection({
    UsernameField? usernamefield,
    PasswordField? passwordfield,
    PayloadType? payloadtype,
  }) {
    final result = create();
    if (usernamefield != null) result.usernamefield = usernamefield;
    if (passwordfield != null) result.passwordfield = passwordfield;
    if (payloadtype != null) result.payloadtype = payloadtype;
    return result;
  }

  RequestInspection._();

  factory RequestInspection.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestInspection.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestInspection',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<UsernameField>(125830068, _omitFieldNames ? '' : 'usernamefield',
        subBuilder: UsernameField.create)
    ..aOM<PasswordField>(318147221, _omitFieldNames ? '' : 'passwordfield',
        subBuilder: PasswordField.create)
    ..aE<PayloadType>(510845422, _omitFieldNames ? '' : 'payloadtype',
        enumValues: PayloadType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestInspection clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestInspection copyWith(void Function(RequestInspection) updates) =>
      super.copyWith((message) => updates(message as RequestInspection))
          as RequestInspection;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestInspection create() => RequestInspection._();
  @$core.override
  RequestInspection createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestInspection getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestInspection>(create);
  static RequestInspection? _defaultInstance;

  @$pb.TagNumber(125830068)
  UsernameField get usernamefield => $_getN(0);
  @$pb.TagNumber(125830068)
  set usernamefield(UsernameField value) => $_setField(125830068, value);
  @$pb.TagNumber(125830068)
  $core.bool hasUsernamefield() => $_has(0);
  @$pb.TagNumber(125830068)
  void clearUsernamefield() => $_clearField(125830068);
  @$pb.TagNumber(125830068)
  UsernameField ensureUsernamefield() => $_ensure(0);

  @$pb.TagNumber(318147221)
  PasswordField get passwordfield => $_getN(1);
  @$pb.TagNumber(318147221)
  set passwordfield(PasswordField value) => $_setField(318147221, value);
  @$pb.TagNumber(318147221)
  $core.bool hasPasswordfield() => $_has(1);
  @$pb.TagNumber(318147221)
  void clearPasswordfield() => $_clearField(318147221);
  @$pb.TagNumber(318147221)
  PasswordField ensurePasswordfield() => $_ensure(1);

  @$pb.TagNumber(510845422)
  PayloadType get payloadtype => $_getN(2);
  @$pb.TagNumber(510845422)
  set payloadtype(PayloadType value) => $_setField(510845422, value);
  @$pb.TagNumber(510845422)
  $core.bool hasPayloadtype() => $_has(2);
  @$pb.TagNumber(510845422)
  void clearPayloadtype() => $_clearField(510845422);
}

class RequestInspectionACFP extends $pb.GeneratedMessage {
  factory RequestInspectionACFP({
    EmailField? emailfield,
    $core.Iterable<PhoneNumberField>? phonenumberfields,
    UsernameField? usernamefield,
    PasswordField? passwordfield,
    $core.Iterable<AddressField>? addressfields,
    PayloadType? payloadtype,
  }) {
    final result = create();
    if (emailfield != null) result.emailfield = emailfield;
    if (phonenumberfields != null)
      result.phonenumberfields.addAll(phonenumberfields);
    if (usernamefield != null) result.usernamefield = usernamefield;
    if (passwordfield != null) result.passwordfield = passwordfield;
    if (addressfields != null) result.addressfields.addAll(addressfields);
    if (payloadtype != null) result.payloadtype = payloadtype;
    return result;
  }

  RequestInspectionACFP._();

  factory RequestInspectionACFP.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestInspectionACFP.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestInspectionACFP',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<EmailField>(23404762, _omitFieldNames ? '' : 'emailfield',
        subBuilder: EmailField.create)
    ..pPM<PhoneNumberField>(
        113611148, _omitFieldNames ? '' : 'phonenumberfields',
        subBuilder: PhoneNumberField.create)
    ..aOM<UsernameField>(125830068, _omitFieldNames ? '' : 'usernamefield',
        subBuilder: UsernameField.create)
    ..aOM<PasswordField>(318147221, _omitFieldNames ? '' : 'passwordfield',
        subBuilder: PasswordField.create)
    ..pPM<AddressField>(364143675, _omitFieldNames ? '' : 'addressfields',
        subBuilder: AddressField.create)
    ..aE<PayloadType>(510845422, _omitFieldNames ? '' : 'payloadtype',
        enumValues: PayloadType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestInspectionACFP clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestInspectionACFP copyWith(
          void Function(RequestInspectionACFP) updates) =>
      super.copyWith((message) => updates(message as RequestInspectionACFP))
          as RequestInspectionACFP;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestInspectionACFP create() => RequestInspectionACFP._();
  @$core.override
  RequestInspectionACFP createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestInspectionACFP getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestInspectionACFP>(create);
  static RequestInspectionACFP? _defaultInstance;

  @$pb.TagNumber(23404762)
  EmailField get emailfield => $_getN(0);
  @$pb.TagNumber(23404762)
  set emailfield(EmailField value) => $_setField(23404762, value);
  @$pb.TagNumber(23404762)
  $core.bool hasEmailfield() => $_has(0);
  @$pb.TagNumber(23404762)
  void clearEmailfield() => $_clearField(23404762);
  @$pb.TagNumber(23404762)
  EmailField ensureEmailfield() => $_ensure(0);

  @$pb.TagNumber(113611148)
  $pb.PbList<PhoneNumberField> get phonenumberfields => $_getList(1);

  @$pb.TagNumber(125830068)
  UsernameField get usernamefield => $_getN(2);
  @$pb.TagNumber(125830068)
  set usernamefield(UsernameField value) => $_setField(125830068, value);
  @$pb.TagNumber(125830068)
  $core.bool hasUsernamefield() => $_has(2);
  @$pb.TagNumber(125830068)
  void clearUsernamefield() => $_clearField(125830068);
  @$pb.TagNumber(125830068)
  UsernameField ensureUsernamefield() => $_ensure(2);

  @$pb.TagNumber(318147221)
  PasswordField get passwordfield => $_getN(3);
  @$pb.TagNumber(318147221)
  set passwordfield(PasswordField value) => $_setField(318147221, value);
  @$pb.TagNumber(318147221)
  $core.bool hasPasswordfield() => $_has(3);
  @$pb.TagNumber(318147221)
  void clearPasswordfield() => $_clearField(318147221);
  @$pb.TagNumber(318147221)
  PasswordField ensurePasswordfield() => $_ensure(3);

  @$pb.TagNumber(364143675)
  $pb.PbList<AddressField> get addressfields => $_getList(4);

  @$pb.TagNumber(510845422)
  PayloadType get payloadtype => $_getN(5);
  @$pb.TagNumber(510845422)
  set payloadtype(PayloadType value) => $_setField(510845422, value);
  @$pb.TagNumber(510845422)
  $core.bool hasPayloadtype() => $_has(5);
  @$pb.TagNumber(510845422)
  void clearPayloadtype() => $_clearField(510845422);
}

class ResponseInspection extends $pb.GeneratedMessage {
  factory ResponseInspection({
    ResponseInspectionBodyContains? bodycontains,
    ResponseInspectionHeader? header,
    ResponseInspectionStatusCode? statuscode,
    ResponseInspectionJson? json,
  }) {
    final result = create();
    if (bodycontains != null) result.bodycontains = bodycontains;
    if (header != null) result.header = header;
    if (statuscode != null) result.statuscode = statuscode;
    if (json != null) result.json = json;
    return result;
  }

  ResponseInspection._();

  factory ResponseInspection.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResponseInspection.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResponseInspection',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<ResponseInspectionBodyContains>(
        118115049, _omitFieldNames ? '' : 'bodycontains',
        subBuilder: ResponseInspectionBodyContains.create)
    ..aOM<ResponseInspectionHeader>(290429313, _omitFieldNames ? '' : 'header',
        subBuilder: ResponseInspectionHeader.create)
    ..aOM<ResponseInspectionStatusCode>(
        303830783, _omitFieldNames ? '' : 'statuscode',
        subBuilder: ResponseInspectionStatusCode.create)
    ..aOM<ResponseInspectionJson>(495111268, _omitFieldNames ? '' : 'json',
        subBuilder: ResponseInspectionJson.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspection clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspection copyWith(void Function(ResponseInspection) updates) =>
      super.copyWith((message) => updates(message as ResponseInspection))
          as ResponseInspection;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResponseInspection create() => ResponseInspection._();
  @$core.override
  ResponseInspection createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResponseInspection getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResponseInspection>(create);
  static ResponseInspection? _defaultInstance;

  @$pb.TagNumber(118115049)
  ResponseInspectionBodyContains get bodycontains => $_getN(0);
  @$pb.TagNumber(118115049)
  set bodycontains(ResponseInspectionBodyContains value) =>
      $_setField(118115049, value);
  @$pb.TagNumber(118115049)
  $core.bool hasBodycontains() => $_has(0);
  @$pb.TagNumber(118115049)
  void clearBodycontains() => $_clearField(118115049);
  @$pb.TagNumber(118115049)
  ResponseInspectionBodyContains ensureBodycontains() => $_ensure(0);

  @$pb.TagNumber(290429313)
  ResponseInspectionHeader get header => $_getN(1);
  @$pb.TagNumber(290429313)
  set header(ResponseInspectionHeader value) => $_setField(290429313, value);
  @$pb.TagNumber(290429313)
  $core.bool hasHeader() => $_has(1);
  @$pb.TagNumber(290429313)
  void clearHeader() => $_clearField(290429313);
  @$pb.TagNumber(290429313)
  ResponseInspectionHeader ensureHeader() => $_ensure(1);

  @$pb.TagNumber(303830783)
  ResponseInspectionStatusCode get statuscode => $_getN(2);
  @$pb.TagNumber(303830783)
  set statuscode(ResponseInspectionStatusCode value) =>
      $_setField(303830783, value);
  @$pb.TagNumber(303830783)
  $core.bool hasStatuscode() => $_has(2);
  @$pb.TagNumber(303830783)
  void clearStatuscode() => $_clearField(303830783);
  @$pb.TagNumber(303830783)
  ResponseInspectionStatusCode ensureStatuscode() => $_ensure(2);

  @$pb.TagNumber(495111268)
  ResponseInspectionJson get json => $_getN(3);
  @$pb.TagNumber(495111268)
  set json(ResponseInspectionJson value) => $_setField(495111268, value);
  @$pb.TagNumber(495111268)
  $core.bool hasJson() => $_has(3);
  @$pb.TagNumber(495111268)
  void clearJson() => $_clearField(495111268);
  @$pb.TagNumber(495111268)
  ResponseInspectionJson ensureJson() => $_ensure(3);
}

class ResponseInspectionBodyContains extends $pb.GeneratedMessage {
  factory ResponseInspectionBodyContains({
    $core.Iterable<$core.String>? successstrings,
    $core.Iterable<$core.String>? failurestrings,
  }) {
    final result = create();
    if (successstrings != null) result.successstrings.addAll(successstrings);
    if (failurestrings != null) result.failurestrings.addAll(failurestrings);
    return result;
  }

  ResponseInspectionBodyContains._();

  factory ResponseInspectionBodyContains.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResponseInspectionBodyContains.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResponseInspectionBodyContains',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(206563583, _omitFieldNames ? '' : 'successstrings')
    ..pPS(244829422, _omitFieldNames ? '' : 'failurestrings')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspectionBodyContains clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspectionBodyContains copyWith(
          void Function(ResponseInspectionBodyContains) updates) =>
      super.copyWith(
              (message) => updates(message as ResponseInspectionBodyContains))
          as ResponseInspectionBodyContains;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResponseInspectionBodyContains create() =>
      ResponseInspectionBodyContains._();
  @$core.override
  ResponseInspectionBodyContains createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResponseInspectionBodyContains getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResponseInspectionBodyContains>(create);
  static ResponseInspectionBodyContains? _defaultInstance;

  @$pb.TagNumber(206563583)
  $pb.PbList<$core.String> get successstrings => $_getList(0);

  @$pb.TagNumber(244829422)
  $pb.PbList<$core.String> get failurestrings => $_getList(1);
}

class ResponseInspectionHeader extends $pb.GeneratedMessage {
  factory ResponseInspectionHeader({
    $core.Iterable<$core.String>? failurevalues,
    $core.String? name,
    $core.Iterable<$core.String>? successvalues,
  }) {
    final result = create();
    if (failurevalues != null) result.failurevalues.addAll(failurevalues);
    if (name != null) result.name = name;
    if (successvalues != null) result.successvalues.addAll(successvalues);
    return result;
  }

  ResponseInspectionHeader._();

  factory ResponseInspectionHeader.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResponseInspectionHeader.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResponseInspectionHeader',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(227897080, _omitFieldNames ? '' : 'failurevalues')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPS(363881655, _omitFieldNames ? '' : 'successvalues')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspectionHeader clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspectionHeader copyWith(
          void Function(ResponseInspectionHeader) updates) =>
      super.copyWith((message) => updates(message as ResponseInspectionHeader))
          as ResponseInspectionHeader;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResponseInspectionHeader create() => ResponseInspectionHeader._();
  @$core.override
  ResponseInspectionHeader createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResponseInspectionHeader getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResponseInspectionHeader>(create);
  static ResponseInspectionHeader? _defaultInstance;

  @$pb.TagNumber(227897080)
  $pb.PbList<$core.String> get failurevalues => $_getList(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(363881655)
  $pb.PbList<$core.String> get successvalues => $_getList(2);
}

class ResponseInspectionJson extends $pb.GeneratedMessage {
  factory ResponseInspectionJson({
    $core.String? identifier,
    $core.Iterable<$core.String>? failurevalues,
    $core.Iterable<$core.String>? successvalues,
  }) {
    final result = create();
    if (identifier != null) result.identifier = identifier;
    if (failurevalues != null) result.failurevalues.addAll(failurevalues);
    if (successvalues != null) result.successvalues.addAll(successvalues);
    return result;
  }

  ResponseInspectionJson._();

  factory ResponseInspectionJson.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResponseInspectionJson.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResponseInspectionJson',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(41865311, _omitFieldNames ? '' : 'identifier')
    ..pPS(227897080, _omitFieldNames ? '' : 'failurevalues')
    ..pPS(363881655, _omitFieldNames ? '' : 'successvalues')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspectionJson clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspectionJson copyWith(
          void Function(ResponseInspectionJson) updates) =>
      super.copyWith((message) => updates(message as ResponseInspectionJson))
          as ResponseInspectionJson;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResponseInspectionJson create() => ResponseInspectionJson._();
  @$core.override
  ResponseInspectionJson createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResponseInspectionJson getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResponseInspectionJson>(create);
  static ResponseInspectionJson? _defaultInstance;

  @$pb.TagNumber(41865311)
  $core.String get identifier => $_getSZ(0);
  @$pb.TagNumber(41865311)
  set identifier($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41865311)
  $core.bool hasIdentifier() => $_has(0);
  @$pb.TagNumber(41865311)
  void clearIdentifier() => $_clearField(41865311);

  @$pb.TagNumber(227897080)
  $pb.PbList<$core.String> get failurevalues => $_getList(1);

  @$pb.TagNumber(363881655)
  $pb.PbList<$core.String> get successvalues => $_getList(2);
}

class ResponseInspectionStatusCode extends $pb.GeneratedMessage {
  factory ResponseInspectionStatusCode({
    $core.Iterable<$core.int>? successcodes,
    $core.Iterable<$core.int>? failurecodes,
  }) {
    final result = create();
    if (successcodes != null) result.successcodes.addAll(successcodes);
    if (failurecodes != null) result.failurecodes.addAll(failurecodes);
    return result;
  }

  ResponseInspectionStatusCode._();

  factory ResponseInspectionStatusCode.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResponseInspectionStatusCode.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResponseInspectionStatusCode',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..p<$core.int>(
        123139107, _omitFieldNames ? '' : 'successcodes', $pb.PbFieldType.K3)
    ..p<$core.int>(
        498971666, _omitFieldNames ? '' : 'failurecodes', $pb.PbFieldType.K3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspectionStatusCode clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResponseInspectionStatusCode copyWith(
          void Function(ResponseInspectionStatusCode) updates) =>
      super.copyWith(
              (message) => updates(message as ResponseInspectionStatusCode))
          as ResponseInspectionStatusCode;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResponseInspectionStatusCode create() =>
      ResponseInspectionStatusCode._();
  @$core.override
  ResponseInspectionStatusCode createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResponseInspectionStatusCode getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResponseInspectionStatusCode>(create);
  static ResponseInspectionStatusCode? _defaultInstance;

  @$pb.TagNumber(123139107)
  $pb.PbList<$core.int> get successcodes => $_getList(0);

  @$pb.TagNumber(498971666)
  $pb.PbList<$core.int> get failurecodes => $_getList(1);
}

class Rule extends $pb.GeneratedMessage {
  factory Rule({
    $core.Iterable<Label>? rulelabels,
    ChallengeConfig? challengeconfig,
    CaptchaConfig? captchaconfig,
    $core.int? priority,
    RuleAction? action,
    Statement? statement,
    $core.String? name,
    VisibilityConfig? visibilityconfig,
    OverrideAction? overrideaction,
  }) {
    final result = create();
    if (rulelabels != null) result.rulelabels.addAll(rulelabels);
    if (challengeconfig != null) result.challengeconfig = challengeconfig;
    if (captchaconfig != null) result.captchaconfig = captchaconfig;
    if (priority != null) result.priority = priority;
    if (action != null) result.action = action;
    if (statement != null) result.statement = statement;
    if (name != null) result.name = name;
    if (visibilityconfig != null) result.visibilityconfig = visibilityconfig;
    if (overrideaction != null) result.overrideaction = overrideaction;
    return result;
  }

  Rule._();

  factory Rule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Rule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Rule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<Label>(14670713, _omitFieldNames ? '' : 'rulelabels',
        subBuilder: Label.create)
    ..aOM<ChallengeConfig>(48990889, _omitFieldNames ? '' : 'challengeconfig',
        subBuilder: ChallengeConfig.create)
    ..aOM<CaptchaConfig>(60547064, _omitFieldNames ? '' : 'captchaconfig',
        subBuilder: CaptchaConfig.create)
    ..aI(109944618, _omitFieldNames ? '' : 'priority')
    ..aOM<RuleAction>(175614240, _omitFieldNames ? '' : 'action',
        subBuilder: RuleAction.create)
    ..aOM<Statement>(248790199, _omitFieldNames ? '' : 'statement',
        subBuilder: Statement.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<VisibilityConfig>(
        507762728, _omitFieldNames ? '' : 'visibilityconfig',
        subBuilder: VisibilityConfig.create)
    ..aOM<OverrideAction>(515842888, _omitFieldNames ? '' : 'overrideaction',
        subBuilder: OverrideAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Rule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Rule copyWith(void Function(Rule) updates) =>
      super.copyWith((message) => updates(message as Rule)) as Rule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Rule create() => Rule._();
  @$core.override
  Rule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Rule getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Rule>(create);
  static Rule? _defaultInstance;

  @$pb.TagNumber(14670713)
  $pb.PbList<Label> get rulelabels => $_getList(0);

  @$pb.TagNumber(48990889)
  ChallengeConfig get challengeconfig => $_getN(1);
  @$pb.TagNumber(48990889)
  set challengeconfig(ChallengeConfig value) => $_setField(48990889, value);
  @$pb.TagNumber(48990889)
  $core.bool hasChallengeconfig() => $_has(1);
  @$pb.TagNumber(48990889)
  void clearChallengeconfig() => $_clearField(48990889);
  @$pb.TagNumber(48990889)
  ChallengeConfig ensureChallengeconfig() => $_ensure(1);

  @$pb.TagNumber(60547064)
  CaptchaConfig get captchaconfig => $_getN(2);
  @$pb.TagNumber(60547064)
  set captchaconfig(CaptchaConfig value) => $_setField(60547064, value);
  @$pb.TagNumber(60547064)
  $core.bool hasCaptchaconfig() => $_has(2);
  @$pb.TagNumber(60547064)
  void clearCaptchaconfig() => $_clearField(60547064);
  @$pb.TagNumber(60547064)
  CaptchaConfig ensureCaptchaconfig() => $_ensure(2);

  @$pb.TagNumber(109944618)
  $core.int get priority => $_getIZ(3);
  @$pb.TagNumber(109944618)
  set priority($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(109944618)
  $core.bool hasPriority() => $_has(3);
  @$pb.TagNumber(109944618)
  void clearPriority() => $_clearField(109944618);

  @$pb.TagNumber(175614240)
  RuleAction get action => $_getN(4);
  @$pb.TagNumber(175614240)
  set action(RuleAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(4);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
  @$pb.TagNumber(175614240)
  RuleAction ensureAction() => $_ensure(4);

  @$pb.TagNumber(248790199)
  Statement get statement => $_getN(5);
  @$pb.TagNumber(248790199)
  set statement(Statement value) => $_setField(248790199, value);
  @$pb.TagNumber(248790199)
  $core.bool hasStatement() => $_has(5);
  @$pb.TagNumber(248790199)
  void clearStatement() => $_clearField(248790199);
  @$pb.TagNumber(248790199)
  Statement ensureStatement() => $_ensure(5);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(507762728)
  VisibilityConfig get visibilityconfig => $_getN(7);
  @$pb.TagNumber(507762728)
  set visibilityconfig(VisibilityConfig value) => $_setField(507762728, value);
  @$pb.TagNumber(507762728)
  $core.bool hasVisibilityconfig() => $_has(7);
  @$pb.TagNumber(507762728)
  void clearVisibilityconfig() => $_clearField(507762728);
  @$pb.TagNumber(507762728)
  VisibilityConfig ensureVisibilityconfig() => $_ensure(7);

  @$pb.TagNumber(515842888)
  OverrideAction get overrideaction => $_getN(8);
  @$pb.TagNumber(515842888)
  set overrideaction(OverrideAction value) => $_setField(515842888, value);
  @$pb.TagNumber(515842888)
  $core.bool hasOverrideaction() => $_has(8);
  @$pb.TagNumber(515842888)
  void clearOverrideaction() => $_clearField(515842888);
  @$pb.TagNumber(515842888)
  OverrideAction ensureOverrideaction() => $_ensure(8);
}

class RuleAction extends $pb.GeneratedMessage {
  factory RuleAction({
    CountAction? count,
    ChallengeAction? challenge,
    AllowAction? allow,
    BlockAction? block,
    CaptchaAction? captcha,
  }) {
    final result = create();
    if (count != null) result.count = count;
    if (challenge != null) result.challenge = challenge;
    if (allow != null) result.allow = allow;
    if (block != null) result.block = block;
    if (captcha != null) result.captcha = captcha;
    return result;
  }

  RuleAction._();

  factory RuleAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RuleAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RuleAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<CountAction>(31963285, _omitFieldNames ? '' : 'count',
        subBuilder: CountAction.create)
    ..aOM<ChallengeAction>(128811743, _omitFieldNames ? '' : 'challenge',
        subBuilder: ChallengeAction.create)
    ..aOM<AllowAction>(342694355, _omitFieldNames ? '' : 'allow',
        subBuilder: AllowAction.create)
    ..aOM<BlockAction>(487014403, _omitFieldNames ? '' : 'block',
        subBuilder: BlockAction.create)
    ..aOM<CaptchaAction>(520084182, _omitFieldNames ? '' : 'captcha',
        subBuilder: CaptchaAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleAction copyWith(void Function(RuleAction) updates) =>
      super.copyWith((message) => updates(message as RuleAction)) as RuleAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RuleAction create() => RuleAction._();
  @$core.override
  RuleAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RuleAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RuleAction>(create);
  static RuleAction? _defaultInstance;

  @$pb.TagNumber(31963285)
  CountAction get count => $_getN(0);
  @$pb.TagNumber(31963285)
  set count(CountAction value) => $_setField(31963285, value);
  @$pb.TagNumber(31963285)
  $core.bool hasCount() => $_has(0);
  @$pb.TagNumber(31963285)
  void clearCount() => $_clearField(31963285);
  @$pb.TagNumber(31963285)
  CountAction ensureCount() => $_ensure(0);

  @$pb.TagNumber(128811743)
  ChallengeAction get challenge => $_getN(1);
  @$pb.TagNumber(128811743)
  set challenge(ChallengeAction value) => $_setField(128811743, value);
  @$pb.TagNumber(128811743)
  $core.bool hasChallenge() => $_has(1);
  @$pb.TagNumber(128811743)
  void clearChallenge() => $_clearField(128811743);
  @$pb.TagNumber(128811743)
  ChallengeAction ensureChallenge() => $_ensure(1);

  @$pb.TagNumber(342694355)
  AllowAction get allow => $_getN(2);
  @$pb.TagNumber(342694355)
  set allow(AllowAction value) => $_setField(342694355, value);
  @$pb.TagNumber(342694355)
  $core.bool hasAllow() => $_has(2);
  @$pb.TagNumber(342694355)
  void clearAllow() => $_clearField(342694355);
  @$pb.TagNumber(342694355)
  AllowAction ensureAllow() => $_ensure(2);

  @$pb.TagNumber(487014403)
  BlockAction get block => $_getN(3);
  @$pb.TagNumber(487014403)
  set block(BlockAction value) => $_setField(487014403, value);
  @$pb.TagNumber(487014403)
  $core.bool hasBlock() => $_has(3);
  @$pb.TagNumber(487014403)
  void clearBlock() => $_clearField(487014403);
  @$pb.TagNumber(487014403)
  BlockAction ensureBlock() => $_ensure(3);

  @$pb.TagNumber(520084182)
  CaptchaAction get captcha => $_getN(4);
  @$pb.TagNumber(520084182)
  set captcha(CaptchaAction value) => $_setField(520084182, value);
  @$pb.TagNumber(520084182)
  $core.bool hasCaptcha() => $_has(4);
  @$pb.TagNumber(520084182)
  void clearCaptcha() => $_clearField(520084182);
  @$pb.TagNumber(520084182)
  CaptchaAction ensureCaptcha() => $_ensure(4);
}

class RuleActionOverride extends $pb.GeneratedMessage {
  factory RuleActionOverride({
    $core.String? name,
    RuleAction? actiontouse,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (actiontouse != null) result.actiontouse = actiontouse;
    return result;
  }

  RuleActionOverride._();

  factory RuleActionOverride.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RuleActionOverride.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RuleActionOverride',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<RuleAction>(470813150, _omitFieldNames ? '' : 'actiontouse',
        subBuilder: RuleAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleActionOverride clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleActionOverride copyWith(void Function(RuleActionOverride) updates) =>
      super.copyWith((message) => updates(message as RuleActionOverride))
          as RuleActionOverride;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RuleActionOverride create() => RuleActionOverride._();
  @$core.override
  RuleActionOverride createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RuleActionOverride getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RuleActionOverride>(create);
  static RuleActionOverride? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(470813150)
  RuleAction get actiontouse => $_getN(1);
  @$pb.TagNumber(470813150)
  set actiontouse(RuleAction value) => $_setField(470813150, value);
  @$pb.TagNumber(470813150)
  $core.bool hasActiontouse() => $_has(1);
  @$pb.TagNumber(470813150)
  void clearActiontouse() => $_clearField(470813150);
  @$pb.TagNumber(470813150)
  RuleAction ensureActiontouse() => $_ensure(1);
}

class RuleGroup extends $pb.GeneratedMessage {
  factory RuleGroup({
    $core.String? labelnamespace,
    $core.Iterable<Rule>? rules,
    $core.Iterable<$core.MapEntry<$core.String, CustomResponseBody>>?
        customresponsebodies,
    $core.Iterable<LabelSummary>? consumedlabels,
    $core.Iterable<LabelSummary>? availablelabels,
    $fixnum.Int64? capacity,
    $core.String? description,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
    VisibilityConfig? visibilityconfig,
  }) {
    final result = create();
    if (labelnamespace != null) result.labelnamespace = labelnamespace;
    if (rules != null) result.rules.addAll(rules);
    if (customresponsebodies != null)
      result.customresponsebodies.addEntries(customresponsebodies);
    if (consumedlabels != null) result.consumedlabels.addAll(consumedlabels);
    if (availablelabels != null) result.availablelabels.addAll(availablelabels);
    if (capacity != null) result.capacity = capacity;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    if (visibilityconfig != null) result.visibilityconfig = visibilityconfig;
    return result;
  }

  RuleGroup._();

  factory RuleGroup.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RuleGroup.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RuleGroup',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(39797417, _omitFieldNames ? '' : 'labelnamespace')
    ..pPM<Rule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: Rule.create)
    ..m<$core.String, CustomResponseBody>(
        42731774, _omitFieldNames ? '' : 'customresponsebodies',
        entryClassName: 'RuleGroup.CustomresponsebodiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: CustomResponseBody.create,
        valueDefaultOrMaker: CustomResponseBody.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..pPM<LabelSummary>(43813949, _omitFieldNames ? '' : 'consumedlabels',
        subBuilder: LabelSummary.create)
    ..pPM<LabelSummary>(51059984, _omitFieldNames ? '' : 'availablelabels',
        subBuilder: LabelSummary.create)
    ..aInt64(107253930, _omitFieldNames ? '' : 'capacity')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..aOM<VisibilityConfig>(
        507762728, _omitFieldNames ? '' : 'visibilityconfig',
        subBuilder: VisibilityConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleGroup clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleGroup copyWith(void Function(RuleGroup) updates) =>
      super.copyWith((message) => updates(message as RuleGroup)) as RuleGroup;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RuleGroup create() => RuleGroup._();
  @$core.override
  RuleGroup createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RuleGroup getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RuleGroup>(create);
  static RuleGroup? _defaultInstance;

  @$pb.TagNumber(39797417)
  $core.String get labelnamespace => $_getSZ(0);
  @$pb.TagNumber(39797417)
  set labelnamespace($core.String value) => $_setString(0, value);
  @$pb.TagNumber(39797417)
  $core.bool hasLabelnamespace() => $_has(0);
  @$pb.TagNumber(39797417)
  void clearLabelnamespace() => $_clearField(39797417);

  @$pb.TagNumber(42675585)
  $pb.PbList<Rule> get rules => $_getList(1);

  @$pb.TagNumber(42731774)
  $pb.PbMap<$core.String, CustomResponseBody> get customresponsebodies =>
      $_getMap(2);

  @$pb.TagNumber(43813949)
  $pb.PbList<LabelSummary> get consumedlabels => $_getList(3);

  @$pb.TagNumber(51059984)
  $pb.PbList<LabelSummary> get availablelabels => $_getList(4);

  @$pb.TagNumber(107253930)
  $fixnum.Int64 get capacity => $_getI64(5);
  @$pb.TagNumber(107253930)
  set capacity($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(107253930)
  $core.bool hasCapacity() => $_has(5);
  @$pb.TagNumber(107253930)
  void clearCapacity() => $_clearField(107253930);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(6);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(6, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(6);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(7);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(7, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(7);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(8);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(8, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(8);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(9);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(9);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(507762728)
  VisibilityConfig get visibilityconfig => $_getN(10);
  @$pb.TagNumber(507762728)
  set visibilityconfig(VisibilityConfig value) => $_setField(507762728, value);
  @$pb.TagNumber(507762728)
  $core.bool hasVisibilityconfig() => $_has(10);
  @$pb.TagNumber(507762728)
  void clearVisibilityconfig() => $_clearField(507762728);
  @$pb.TagNumber(507762728)
  VisibilityConfig ensureVisibilityconfig() => $_ensure(10);
}

class RuleGroupReferenceStatement extends $pb.GeneratedMessage {
  factory RuleGroupReferenceStatement({
    $core.Iterable<ExcludedRule>? excludedrules,
    $core.Iterable<RuleActionOverride>? ruleactionoverrides,
    $core.String? arn,
  }) {
    final result = create();
    if (excludedrules != null) result.excludedrules.addAll(excludedrules);
    if (ruleactionoverrides != null)
      result.ruleactionoverrides.addAll(ruleactionoverrides);
    if (arn != null) result.arn = arn;
    return result;
  }

  RuleGroupReferenceStatement._();

  factory RuleGroupReferenceStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RuleGroupReferenceStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RuleGroupReferenceStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<ExcludedRule>(129929071, _omitFieldNames ? '' : 'excludedrules',
        subBuilder: ExcludedRule.create)
    ..pPM<RuleActionOverride>(
        354236935, _omitFieldNames ? '' : 'ruleactionoverrides',
        subBuilder: RuleActionOverride.create)
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleGroupReferenceStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleGroupReferenceStatement copyWith(
          void Function(RuleGroupReferenceStatement) updates) =>
      super.copyWith(
              (message) => updates(message as RuleGroupReferenceStatement))
          as RuleGroupReferenceStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RuleGroupReferenceStatement create() =>
      RuleGroupReferenceStatement._();
  @$core.override
  RuleGroupReferenceStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RuleGroupReferenceStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RuleGroupReferenceStatement>(create);
  static RuleGroupReferenceStatement? _defaultInstance;

  @$pb.TagNumber(129929071)
  $pb.PbList<ExcludedRule> get excludedrules => $_getList(0);

  @$pb.TagNumber(354236935)
  $pb.PbList<RuleActionOverride> get ruleactionoverrides => $_getList(1);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class RuleGroupSummary extends $pb.GeneratedMessage {
  factory RuleGroupSummary({
    $core.String? locktoken,
    $core.String? description,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  RuleGroupSummary._();

  factory RuleGroupSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RuleGroupSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RuleGroupSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleGroupSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleGroupSummary copyWith(void Function(RuleGroupSummary) updates) =>
      super.copyWith((message) => updates(message as RuleGroupSummary))
          as RuleGroupSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RuleGroupSummary create() => RuleGroupSummary._();
  @$core.override
  RuleGroupSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RuleGroupSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RuleGroupSummary>(create);
  static RuleGroupSummary? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class RuleSummary extends $pb.GeneratedMessage {
  factory RuleSummary({
    RuleAction? action,
    $core.String? name,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (name != null) result.name = name;
    return result;
  }

  RuleSummary._();

  factory RuleSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RuleSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RuleSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<RuleAction>(175614240, _omitFieldNames ? '' : 'action',
        subBuilder: RuleAction.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleSummary copyWith(void Function(RuleSummary) updates) =>
      super.copyWith((message) => updates(message as RuleSummary))
          as RuleSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RuleSummary create() => RuleSummary._();
  @$core.override
  RuleSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RuleSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RuleSummary>(create);
  static RuleSummary? _defaultInstance;

  @$pb.TagNumber(175614240)
  RuleAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(RuleAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
  @$pb.TagNumber(175614240)
  RuleAction ensureAction() => $_ensure(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class SampledHTTPRequest extends $pb.GeneratedMessage {
  factory SampledHTTPRequest({
    $core.String? overriddenaction,
    HTTPRequest? request,
    $core.int? responsecodesent,
    $core.String? timestamp,
    $core.String? action,
    $core.Iterable<Label>? labels,
    ChallengeResponse? challengeresponse,
    $core.String? rulenamewithinrulegroup,
    $core.Iterable<HTTPHeader>? requestheadersinserted,
    $fixnum.Int64? weight,
    CaptchaResponse? captcharesponse,
  }) {
    final result = create();
    if (overriddenaction != null) result.overriddenaction = overriddenaction;
    if (request != null) result.request = request;
    if (responsecodesent != null) result.responsecodesent = responsecodesent;
    if (timestamp != null) result.timestamp = timestamp;
    if (action != null) result.action = action;
    if (labels != null) result.labels.addAll(labels);
    if (challengeresponse != null) result.challengeresponse = challengeresponse;
    if (rulenamewithinrulegroup != null)
      result.rulenamewithinrulegroup = rulenamewithinrulegroup;
    if (requestheadersinserted != null)
      result.requestheadersinserted.addAll(requestheadersinserted);
    if (weight != null) result.weight = weight;
    if (captcharesponse != null) result.captcharesponse = captcharesponse;
    return result;
  }

  SampledHTTPRequest._();

  factory SampledHTTPRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SampledHTTPRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SampledHTTPRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(35244204, _omitFieldNames ? '' : 'overriddenaction')
    ..aOM<HTTPRequest>(38093139, _omitFieldNames ? '' : 'request',
        subBuilder: HTTPRequest.create)
    ..aI(108347942, _omitFieldNames ? '' : 'responsecodesent')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aOS(175614240, _omitFieldNames ? '' : 'action')
    ..pPM<Label>(178416811, _omitFieldNames ? '' : 'labels',
        subBuilder: Label.create)
    ..aOM<ChallengeResponse>(
        268501730, _omitFieldNames ? '' : 'challengeresponse',
        subBuilder: ChallengeResponse.create)
    ..aOS(317544521, _omitFieldNames ? '' : 'rulenamewithinrulegroup')
    ..pPM<HTTPHeader>(
        367751765, _omitFieldNames ? '' : 'requestheadersinserted',
        subBuilder: HTTPHeader.create)
    ..aInt64(422581466, _omitFieldNames ? '' : 'weight')
    ..aOM<CaptchaResponse>(484053583, _omitFieldNames ? '' : 'captcharesponse',
        subBuilder: CaptchaResponse.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SampledHTTPRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SampledHTTPRequest copyWith(void Function(SampledHTTPRequest) updates) =>
      super.copyWith((message) => updates(message as SampledHTTPRequest))
          as SampledHTTPRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SampledHTTPRequest create() => SampledHTTPRequest._();
  @$core.override
  SampledHTTPRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SampledHTTPRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SampledHTTPRequest>(create);
  static SampledHTTPRequest? _defaultInstance;

  @$pb.TagNumber(35244204)
  $core.String get overriddenaction => $_getSZ(0);
  @$pb.TagNumber(35244204)
  set overriddenaction($core.String value) => $_setString(0, value);
  @$pb.TagNumber(35244204)
  $core.bool hasOverriddenaction() => $_has(0);
  @$pb.TagNumber(35244204)
  void clearOverriddenaction() => $_clearField(35244204);

  @$pb.TagNumber(38093139)
  HTTPRequest get request => $_getN(1);
  @$pb.TagNumber(38093139)
  set request(HTTPRequest value) => $_setField(38093139, value);
  @$pb.TagNumber(38093139)
  $core.bool hasRequest() => $_has(1);
  @$pb.TagNumber(38093139)
  void clearRequest() => $_clearField(38093139);
  @$pb.TagNumber(38093139)
  HTTPRequest ensureRequest() => $_ensure(1);

  @$pb.TagNumber(108347942)
  $core.int get responsecodesent => $_getIZ(2);
  @$pb.TagNumber(108347942)
  set responsecodesent($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(108347942)
  $core.bool hasResponsecodesent() => $_has(2);
  @$pb.TagNumber(108347942)
  void clearResponsecodesent() => $_clearField(108347942);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(3);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(3);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(175614240)
  $core.String get action => $_getSZ(4);
  @$pb.TagNumber(175614240)
  set action($core.String value) => $_setString(4, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(4);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(178416811)
  $pb.PbList<Label> get labels => $_getList(5);

  @$pb.TagNumber(268501730)
  ChallengeResponse get challengeresponse => $_getN(6);
  @$pb.TagNumber(268501730)
  set challengeresponse(ChallengeResponse value) =>
      $_setField(268501730, value);
  @$pb.TagNumber(268501730)
  $core.bool hasChallengeresponse() => $_has(6);
  @$pb.TagNumber(268501730)
  void clearChallengeresponse() => $_clearField(268501730);
  @$pb.TagNumber(268501730)
  ChallengeResponse ensureChallengeresponse() => $_ensure(6);

  @$pb.TagNumber(317544521)
  $core.String get rulenamewithinrulegroup => $_getSZ(7);
  @$pb.TagNumber(317544521)
  set rulenamewithinrulegroup($core.String value) => $_setString(7, value);
  @$pb.TagNumber(317544521)
  $core.bool hasRulenamewithinrulegroup() => $_has(7);
  @$pb.TagNumber(317544521)
  void clearRulenamewithinrulegroup() => $_clearField(317544521);

  @$pb.TagNumber(367751765)
  $pb.PbList<HTTPHeader> get requestheadersinserted => $_getList(8);

  @$pb.TagNumber(422581466)
  $fixnum.Int64 get weight => $_getI64(9);
  @$pb.TagNumber(422581466)
  set weight($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(422581466)
  $core.bool hasWeight() => $_has(9);
  @$pb.TagNumber(422581466)
  void clearWeight() => $_clearField(422581466);

  @$pb.TagNumber(484053583)
  CaptchaResponse get captcharesponse => $_getN(10);
  @$pb.TagNumber(484053583)
  set captcharesponse(CaptchaResponse value) => $_setField(484053583, value);
  @$pb.TagNumber(484053583)
  $core.bool hasCaptcharesponse() => $_has(10);
  @$pb.TagNumber(484053583)
  void clearCaptcharesponse() => $_clearField(484053583);
  @$pb.TagNumber(484053583)
  CaptchaResponse ensureCaptcharesponse() => $_ensure(10);
}

class SingleHeader extends $pb.GeneratedMessage {
  factory SingleHeader({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  SingleHeader._();

  factory SingleHeader.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SingleHeader.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SingleHeader',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SingleHeader clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SingleHeader copyWith(void Function(SingleHeader) updates) =>
      super.copyWith((message) => updates(message as SingleHeader))
          as SingleHeader;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SingleHeader create() => SingleHeader._();
  @$core.override
  SingleHeader createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SingleHeader getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SingleHeader>(create);
  static SingleHeader? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class SingleQueryArgument extends $pb.GeneratedMessage {
  factory SingleQueryArgument({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  SingleQueryArgument._();

  factory SingleQueryArgument.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SingleQueryArgument.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SingleQueryArgument',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SingleQueryArgument clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SingleQueryArgument copyWith(void Function(SingleQueryArgument) updates) =>
      super.copyWith((message) => updates(message as SingleQueryArgument))
          as SingleQueryArgument;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SingleQueryArgument create() => SingleQueryArgument._();
  @$core.override
  SingleQueryArgument createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SingleQueryArgument getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SingleQueryArgument>(create);
  static SingleQueryArgument? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class SizeConstraintStatement extends $pb.GeneratedMessage {
  factory SizeConstraintStatement({
    ComparisonOperator? comparisonoperator,
    $fixnum.Int64? size,
    $core.Iterable<TextTransformation>? texttransformations,
    FieldToMatch? fieldtomatch,
  }) {
    final result = create();
    if (comparisonoperator != null)
      result.comparisonoperator = comparisonoperator;
    if (size != null) result.size = size;
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    return result;
  }

  SizeConstraintStatement._();

  factory SizeConstraintStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SizeConstraintStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SizeConstraintStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<ComparisonOperator>(
        7059861, _omitFieldNames ? '' : 'comparisonoperator',
        enumValues: ComparisonOperator.values)
    ..aInt64(105352829, _omitFieldNames ? '' : 'size')
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraintStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraintStatement copyWith(
          void Function(SizeConstraintStatement) updates) =>
      super.copyWith((message) => updates(message as SizeConstraintStatement))
          as SizeConstraintStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SizeConstraintStatement create() => SizeConstraintStatement._();
  @$core.override
  SizeConstraintStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SizeConstraintStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SizeConstraintStatement>(create);
  static SizeConstraintStatement? _defaultInstance;

  @$pb.TagNumber(7059861)
  ComparisonOperator get comparisonoperator => $_getN(0);
  @$pb.TagNumber(7059861)
  set comparisonoperator(ComparisonOperator value) =>
      $_setField(7059861, value);
  @$pb.TagNumber(7059861)
  $core.bool hasComparisonoperator() => $_has(0);
  @$pb.TagNumber(7059861)
  void clearComparisonoperator() => $_clearField(7059861);

  @$pb.TagNumber(105352829)
  $fixnum.Int64 get size => $_getI64(1);
  @$pb.TagNumber(105352829)
  set size($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(105352829)
  $core.bool hasSize() => $_has(1);
  @$pb.TagNumber(105352829)
  void clearSize() => $_clearField(105352829);

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(2);

  @$pb.TagNumber(338348372)
  FieldToMatch get fieldtomatch => $_getN(3);
  @$pb.TagNumber(338348372)
  set fieldtomatch(FieldToMatch value) => $_setField(338348372, value);
  @$pb.TagNumber(338348372)
  $core.bool hasFieldtomatch() => $_has(3);
  @$pb.TagNumber(338348372)
  void clearFieldtomatch() => $_clearField(338348372);
  @$pb.TagNumber(338348372)
  FieldToMatch ensureFieldtomatch() => $_ensure(3);
}

class SqliMatchStatement extends $pb.GeneratedMessage {
  factory SqliMatchStatement({
    SensitivityLevel? sensitivitylevel,
    $core.Iterable<TextTransformation>? texttransformations,
    FieldToMatch? fieldtomatch,
  }) {
    final result = create();
    if (sensitivitylevel != null) result.sensitivitylevel = sensitivitylevel;
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    return result;
  }

  SqliMatchStatement._();

  factory SqliMatchStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SqliMatchStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SqliMatchStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<SensitivityLevel>(12020595, _omitFieldNames ? '' : 'sensitivitylevel',
        enumValues: SensitivityLevel.values)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqliMatchStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqliMatchStatement copyWith(void Function(SqliMatchStatement) updates) =>
      super.copyWith((message) => updates(message as SqliMatchStatement))
          as SqliMatchStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SqliMatchStatement create() => SqliMatchStatement._();
  @$core.override
  SqliMatchStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SqliMatchStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SqliMatchStatement>(create);
  static SqliMatchStatement? _defaultInstance;

  @$pb.TagNumber(12020595)
  SensitivityLevel get sensitivitylevel => $_getN(0);
  @$pb.TagNumber(12020595)
  set sensitivitylevel(SensitivityLevel value) => $_setField(12020595, value);
  @$pb.TagNumber(12020595)
  $core.bool hasSensitivitylevel() => $_has(0);
  @$pb.TagNumber(12020595)
  void clearSensitivitylevel() => $_clearField(12020595);

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(1);

  @$pb.TagNumber(338348372)
  FieldToMatch get fieldtomatch => $_getN(2);
  @$pb.TagNumber(338348372)
  set fieldtomatch(FieldToMatch value) => $_setField(338348372, value);
  @$pb.TagNumber(338348372)
  $core.bool hasFieldtomatch() => $_has(2);
  @$pb.TagNumber(338348372)
  void clearFieldtomatch() => $_clearField(338348372);
  @$pb.TagNumber(338348372)
  FieldToMatch ensureFieldtomatch() => $_ensure(2);
}

class Statement extends $pb.GeneratedMessage {
  factory Statement({
    LabelMatchStatement? labelmatchstatement,
    RegexMatchStatement? regexmatchstatement,
    SqliMatchStatement? sqlimatchstatement,
    ByteMatchStatement? bytematchstatement,
    SizeConstraintStatement? sizeconstraintstatement,
    GeoMatchStatement? geomatchstatement,
    OrStatement? orstatement,
    AndStatement? andstatement,
    RateBasedStatement? ratebasedstatement,
    XssMatchStatement? xssmatchstatement,
    RegexPatternSetReferenceStatement? regexpatternsetreferencestatement,
    IPSetReferenceStatement? ipsetreferencestatement,
    AsnMatchStatement? asnmatchstatement,
    RuleGroupReferenceStatement? rulegroupreferencestatement,
    NotStatement? notstatement,
    ManagedRuleGroupStatement? managedrulegroupstatement,
  }) {
    final result = create();
    if (labelmatchstatement != null)
      result.labelmatchstatement = labelmatchstatement;
    if (regexmatchstatement != null)
      result.regexmatchstatement = regexmatchstatement;
    if (sqlimatchstatement != null)
      result.sqlimatchstatement = sqlimatchstatement;
    if (bytematchstatement != null)
      result.bytematchstatement = bytematchstatement;
    if (sizeconstraintstatement != null)
      result.sizeconstraintstatement = sizeconstraintstatement;
    if (geomatchstatement != null) result.geomatchstatement = geomatchstatement;
    if (orstatement != null) result.orstatement = orstatement;
    if (andstatement != null) result.andstatement = andstatement;
    if (ratebasedstatement != null)
      result.ratebasedstatement = ratebasedstatement;
    if (xssmatchstatement != null) result.xssmatchstatement = xssmatchstatement;
    if (regexpatternsetreferencestatement != null)
      result.regexpatternsetreferencestatement =
          regexpatternsetreferencestatement;
    if (ipsetreferencestatement != null)
      result.ipsetreferencestatement = ipsetreferencestatement;
    if (asnmatchstatement != null) result.asnmatchstatement = asnmatchstatement;
    if (rulegroupreferencestatement != null)
      result.rulegroupreferencestatement = rulegroupreferencestatement;
    if (notstatement != null) result.notstatement = notstatement;
    if (managedrulegroupstatement != null)
      result.managedrulegroupstatement = managedrulegroupstatement;
    return result;
  }

  Statement._();

  factory Statement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Statement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Statement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOM<LabelMatchStatement>(
        46165448, _omitFieldNames ? '' : 'labelmatchstatement',
        subBuilder: LabelMatchStatement.create)
    ..aOM<RegexMatchStatement>(
        47948763, _omitFieldNames ? '' : 'regexmatchstatement',
        subBuilder: RegexMatchStatement.create)
    ..aOM<SqliMatchStatement>(
        53921561, _omitFieldNames ? '' : 'sqlimatchstatement',
        subBuilder: SqliMatchStatement.create)
    ..aOM<ByteMatchStatement>(
        215080438, _omitFieldNames ? '' : 'bytematchstatement',
        subBuilder: ByteMatchStatement.create)
    ..aOM<SizeConstraintStatement>(
        244031361, _omitFieldNames ? '' : 'sizeconstraintstatement',
        subBuilder: SizeConstraintStatement.create)
    ..aOM<GeoMatchStatement>(
        275553657, _omitFieldNames ? '' : 'geomatchstatement',
        subBuilder: GeoMatchStatement.create)
    ..aOM<OrStatement>(281998178, _omitFieldNames ? '' : 'orstatement',
        subBuilder: OrStatement.create)
    ..aOM<AndStatement>(345032416, _omitFieldNames ? '' : 'andstatement',
        subBuilder: AndStatement.create)
    ..aOM<RateBasedStatement>(
        357469284, _omitFieldNames ? '' : 'ratebasedstatement',
        subBuilder: RateBasedStatement.create)
    ..aOM<XssMatchStatement>(
        383220104, _omitFieldNames ? '' : 'xssmatchstatement',
        subBuilder: XssMatchStatement.create)
    ..aOM<RegexPatternSetReferenceStatement>(
        409481069, _omitFieldNames ? '' : 'regexpatternsetreferencestatement',
        subBuilder: RegexPatternSetReferenceStatement.create)
    ..aOM<IPSetReferenceStatement>(
        419813595, _omitFieldNames ? '' : 'ipsetreferencestatement',
        subBuilder: IPSetReferenceStatement.create)
    ..aOM<AsnMatchStatement>(
        432340332, _omitFieldNames ? '' : 'asnmatchstatement',
        subBuilder: AsnMatchStatement.create)
    ..aOM<RuleGroupReferenceStatement>(
        460139461, _omitFieldNames ? '' : 'rulegroupreferencestatement',
        subBuilder: RuleGroupReferenceStatement.create)
    ..aOM<NotStatement>(474254724, _omitFieldNames ? '' : 'notstatement',
        subBuilder: NotStatement.create)
    ..aOM<ManagedRuleGroupStatement>(
        533120745, _omitFieldNames ? '' : 'managedrulegroupstatement',
        subBuilder: ManagedRuleGroupStatement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Statement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Statement copyWith(void Function(Statement) updates) =>
      super.copyWith((message) => updates(message as Statement)) as Statement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Statement create() => Statement._();
  @$core.override
  Statement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Statement getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Statement>(create);
  static Statement? _defaultInstance;

  @$pb.TagNumber(46165448)
  LabelMatchStatement get labelmatchstatement => $_getN(0);
  @$pb.TagNumber(46165448)
  set labelmatchstatement(LabelMatchStatement value) =>
      $_setField(46165448, value);
  @$pb.TagNumber(46165448)
  $core.bool hasLabelmatchstatement() => $_has(0);
  @$pb.TagNumber(46165448)
  void clearLabelmatchstatement() => $_clearField(46165448);
  @$pb.TagNumber(46165448)
  LabelMatchStatement ensureLabelmatchstatement() => $_ensure(0);

  @$pb.TagNumber(47948763)
  RegexMatchStatement get regexmatchstatement => $_getN(1);
  @$pb.TagNumber(47948763)
  set regexmatchstatement(RegexMatchStatement value) =>
      $_setField(47948763, value);
  @$pb.TagNumber(47948763)
  $core.bool hasRegexmatchstatement() => $_has(1);
  @$pb.TagNumber(47948763)
  void clearRegexmatchstatement() => $_clearField(47948763);
  @$pb.TagNumber(47948763)
  RegexMatchStatement ensureRegexmatchstatement() => $_ensure(1);

  @$pb.TagNumber(53921561)
  SqliMatchStatement get sqlimatchstatement => $_getN(2);
  @$pb.TagNumber(53921561)
  set sqlimatchstatement(SqliMatchStatement value) =>
      $_setField(53921561, value);
  @$pb.TagNumber(53921561)
  $core.bool hasSqlimatchstatement() => $_has(2);
  @$pb.TagNumber(53921561)
  void clearSqlimatchstatement() => $_clearField(53921561);
  @$pb.TagNumber(53921561)
  SqliMatchStatement ensureSqlimatchstatement() => $_ensure(2);

  @$pb.TagNumber(215080438)
  ByteMatchStatement get bytematchstatement => $_getN(3);
  @$pb.TagNumber(215080438)
  set bytematchstatement(ByteMatchStatement value) =>
      $_setField(215080438, value);
  @$pb.TagNumber(215080438)
  $core.bool hasBytematchstatement() => $_has(3);
  @$pb.TagNumber(215080438)
  void clearBytematchstatement() => $_clearField(215080438);
  @$pb.TagNumber(215080438)
  ByteMatchStatement ensureBytematchstatement() => $_ensure(3);

  @$pb.TagNumber(244031361)
  SizeConstraintStatement get sizeconstraintstatement => $_getN(4);
  @$pb.TagNumber(244031361)
  set sizeconstraintstatement(SizeConstraintStatement value) =>
      $_setField(244031361, value);
  @$pb.TagNumber(244031361)
  $core.bool hasSizeconstraintstatement() => $_has(4);
  @$pb.TagNumber(244031361)
  void clearSizeconstraintstatement() => $_clearField(244031361);
  @$pb.TagNumber(244031361)
  SizeConstraintStatement ensureSizeconstraintstatement() => $_ensure(4);

  @$pb.TagNumber(275553657)
  GeoMatchStatement get geomatchstatement => $_getN(5);
  @$pb.TagNumber(275553657)
  set geomatchstatement(GeoMatchStatement value) =>
      $_setField(275553657, value);
  @$pb.TagNumber(275553657)
  $core.bool hasGeomatchstatement() => $_has(5);
  @$pb.TagNumber(275553657)
  void clearGeomatchstatement() => $_clearField(275553657);
  @$pb.TagNumber(275553657)
  GeoMatchStatement ensureGeomatchstatement() => $_ensure(5);

  @$pb.TagNumber(281998178)
  OrStatement get orstatement => $_getN(6);
  @$pb.TagNumber(281998178)
  set orstatement(OrStatement value) => $_setField(281998178, value);
  @$pb.TagNumber(281998178)
  $core.bool hasOrstatement() => $_has(6);
  @$pb.TagNumber(281998178)
  void clearOrstatement() => $_clearField(281998178);
  @$pb.TagNumber(281998178)
  OrStatement ensureOrstatement() => $_ensure(6);

  @$pb.TagNumber(345032416)
  AndStatement get andstatement => $_getN(7);
  @$pb.TagNumber(345032416)
  set andstatement(AndStatement value) => $_setField(345032416, value);
  @$pb.TagNumber(345032416)
  $core.bool hasAndstatement() => $_has(7);
  @$pb.TagNumber(345032416)
  void clearAndstatement() => $_clearField(345032416);
  @$pb.TagNumber(345032416)
  AndStatement ensureAndstatement() => $_ensure(7);

  @$pb.TagNumber(357469284)
  RateBasedStatement get ratebasedstatement => $_getN(8);
  @$pb.TagNumber(357469284)
  set ratebasedstatement(RateBasedStatement value) =>
      $_setField(357469284, value);
  @$pb.TagNumber(357469284)
  $core.bool hasRatebasedstatement() => $_has(8);
  @$pb.TagNumber(357469284)
  void clearRatebasedstatement() => $_clearField(357469284);
  @$pb.TagNumber(357469284)
  RateBasedStatement ensureRatebasedstatement() => $_ensure(8);

  @$pb.TagNumber(383220104)
  XssMatchStatement get xssmatchstatement => $_getN(9);
  @$pb.TagNumber(383220104)
  set xssmatchstatement(XssMatchStatement value) =>
      $_setField(383220104, value);
  @$pb.TagNumber(383220104)
  $core.bool hasXssmatchstatement() => $_has(9);
  @$pb.TagNumber(383220104)
  void clearXssmatchstatement() => $_clearField(383220104);
  @$pb.TagNumber(383220104)
  XssMatchStatement ensureXssmatchstatement() => $_ensure(9);

  @$pb.TagNumber(409481069)
  RegexPatternSetReferenceStatement get regexpatternsetreferencestatement =>
      $_getN(10);
  @$pb.TagNumber(409481069)
  set regexpatternsetreferencestatement(
          RegexPatternSetReferenceStatement value) =>
      $_setField(409481069, value);
  @$pb.TagNumber(409481069)
  $core.bool hasRegexpatternsetreferencestatement() => $_has(10);
  @$pb.TagNumber(409481069)
  void clearRegexpatternsetreferencestatement() => $_clearField(409481069);
  @$pb.TagNumber(409481069)
  RegexPatternSetReferenceStatement ensureRegexpatternsetreferencestatement() =>
      $_ensure(10);

  @$pb.TagNumber(419813595)
  IPSetReferenceStatement get ipsetreferencestatement => $_getN(11);
  @$pb.TagNumber(419813595)
  set ipsetreferencestatement(IPSetReferenceStatement value) =>
      $_setField(419813595, value);
  @$pb.TagNumber(419813595)
  $core.bool hasIpsetreferencestatement() => $_has(11);
  @$pb.TagNumber(419813595)
  void clearIpsetreferencestatement() => $_clearField(419813595);
  @$pb.TagNumber(419813595)
  IPSetReferenceStatement ensureIpsetreferencestatement() => $_ensure(11);

  @$pb.TagNumber(432340332)
  AsnMatchStatement get asnmatchstatement => $_getN(12);
  @$pb.TagNumber(432340332)
  set asnmatchstatement(AsnMatchStatement value) =>
      $_setField(432340332, value);
  @$pb.TagNumber(432340332)
  $core.bool hasAsnmatchstatement() => $_has(12);
  @$pb.TagNumber(432340332)
  void clearAsnmatchstatement() => $_clearField(432340332);
  @$pb.TagNumber(432340332)
  AsnMatchStatement ensureAsnmatchstatement() => $_ensure(12);

  @$pb.TagNumber(460139461)
  RuleGroupReferenceStatement get rulegroupreferencestatement => $_getN(13);
  @$pb.TagNumber(460139461)
  set rulegroupreferencestatement(RuleGroupReferenceStatement value) =>
      $_setField(460139461, value);
  @$pb.TagNumber(460139461)
  $core.bool hasRulegroupreferencestatement() => $_has(13);
  @$pb.TagNumber(460139461)
  void clearRulegroupreferencestatement() => $_clearField(460139461);
  @$pb.TagNumber(460139461)
  RuleGroupReferenceStatement ensureRulegroupreferencestatement() =>
      $_ensure(13);

  @$pb.TagNumber(474254724)
  NotStatement get notstatement => $_getN(14);
  @$pb.TagNumber(474254724)
  set notstatement(NotStatement value) => $_setField(474254724, value);
  @$pb.TagNumber(474254724)
  $core.bool hasNotstatement() => $_has(14);
  @$pb.TagNumber(474254724)
  void clearNotstatement() => $_clearField(474254724);
  @$pb.TagNumber(474254724)
  NotStatement ensureNotstatement() => $_ensure(14);

  @$pb.TagNumber(533120745)
  ManagedRuleGroupStatement get managedrulegroupstatement => $_getN(15);
  @$pb.TagNumber(533120745)
  set managedrulegroupstatement(ManagedRuleGroupStatement value) =>
      $_setField(533120745, value);
  @$pb.TagNumber(533120745)
  $core.bool hasManagedrulegroupstatement() => $_has(15);
  @$pb.TagNumber(533120745)
  void clearManagedrulegroupstatement() => $_clearField(533120745);
  @$pb.TagNumber(533120745)
  ManagedRuleGroupStatement ensureManagedrulegroupstatement() => $_ensure(15);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
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

class TagInfoForResource extends $pb.GeneratedMessage {
  factory TagInfoForResource({
    $core.String? resourcearn,
    $core.Iterable<Tag>? taglist,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (taglist != null) result.taglist.addAll(taglist);
    return result;
  }

  TagInfoForResource._();

  factory TagInfoForResource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagInfoForResource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagInfoForResource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..pPM<Tag>(429416860, _omitFieldNames ? '' : 'taglist',
        subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagInfoForResource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagInfoForResource copyWith(void Function(TagInfoForResource) updates) =>
      super.copyWith((message) => updates(message as TagInfoForResource))
          as TagInfoForResource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagInfoForResource create() => TagInfoForResource._();
  @$core.override
  TagInfoForResource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagInfoForResource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagInfoForResource>(create);
  static TagInfoForResource? _defaultInstance;

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(429416860)
  $pb.PbList<Tag> get taglist => $_getList(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
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

class TextTransformation extends $pb.GeneratedMessage {
  factory TextTransformation({
    $core.int? priority,
    TextTransformationType? type,
  }) {
    final result = create();
    if (priority != null) result.priority = priority;
    if (type != null) result.type = type;
    return result;
  }

  TextTransformation._();

  factory TextTransformation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TextTransformation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TextTransformation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aI(109944618, _omitFieldNames ? '' : 'priority')
    ..aE<TextTransformationType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: TextTransformationType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TextTransformation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TextTransformation copyWith(void Function(TextTransformation) updates) =>
      super.copyWith((message) => updates(message as TextTransformation))
          as TextTransformation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TextTransformation create() => TextTransformation._();
  @$core.override
  TextTransformation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TextTransformation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TextTransformation>(create);
  static TextTransformation? _defaultInstance;

  @$pb.TagNumber(109944618)
  $core.int get priority => $_getIZ(0);
  @$pb.TagNumber(109944618)
  set priority($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(109944618)
  $core.bool hasPriority() => $_has(0);
  @$pb.TagNumber(109944618)
  void clearPriority() => $_clearField(109944618);

  @$pb.TagNumber(290836590)
  TextTransformationType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(TextTransformationType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class TimeWindow extends $pb.GeneratedMessage {
  factory TimeWindow({
    $core.String? endtime,
    $core.String? starttime,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (starttime != null) result.starttime = starttime;
    return result;
  }

  TimeWindow._();

  factory TimeWindow.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TimeWindow.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TimeWindow',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeWindow clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeWindow copyWith(void Function(TimeWindow) updates) =>
      super.copyWith((message) => updates(message as TimeWindow)) as TimeWindow;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TimeWindow create() => TimeWindow._();
  @$core.override
  TimeWindow createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TimeWindow getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TimeWindow>(create);
  static TimeWindow? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(1);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(1);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
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

class UpdateIPSetRequest extends $pb.GeneratedMessage {
  factory UpdateIPSetRequest({
    $core.String? locktoken,
    Scope? scope,
    $core.String? description,
    $core.String? name,
    $core.Iterable<$core.String>? addresses,
    $core.String? id,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (scope != null) result.scope = scope;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (addresses != null) result.addresses.addAll(addresses);
    if (id != null) result.id = id;
    return result;
  }

  UpdateIPSetRequest._();

  factory UpdateIPSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateIPSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateIPSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPS(375939972, _omitFieldNames ? '' : 'addresses')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateIPSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateIPSetRequest copyWith(void Function(UpdateIPSetRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateIPSetRequest))
          as UpdateIPSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateIPSetRequest create() => UpdateIPSetRequest._();
  @$core.override
  UpdateIPSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateIPSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateIPSetRequest>(create);
  static UpdateIPSetRequest? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(375939972)
  $pb.PbList<$core.String> get addresses => $_getList(4);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(5);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(5, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(5);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class UpdateIPSetResponse extends $pb.GeneratedMessage {
  factory UpdateIPSetResponse({
    $core.String? nextlocktoken,
  }) {
    final result = create();
    if (nextlocktoken != null) result.nextlocktoken = nextlocktoken;
    return result;
  }

  UpdateIPSetResponse._();

  factory UpdateIPSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateIPSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateIPSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(473728437, _omitFieldNames ? '' : 'nextlocktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateIPSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateIPSetResponse copyWith(void Function(UpdateIPSetResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateIPSetResponse))
          as UpdateIPSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateIPSetResponse create() => UpdateIPSetResponse._();
  @$core.override
  UpdateIPSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateIPSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateIPSetResponse>(create);
  static UpdateIPSetResponse? _defaultInstance;

  @$pb.TagNumber(473728437)
  $core.String get nextlocktoken => $_getSZ(0);
  @$pb.TagNumber(473728437)
  set nextlocktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(473728437)
  $core.bool hasNextlocktoken() => $_has(0);
  @$pb.TagNumber(473728437)
  void clearNextlocktoken() => $_clearField(473728437);
}

class UpdateManagedRuleSetVersionExpiryDateRequest
    extends $pb.GeneratedMessage {
  factory UpdateManagedRuleSetVersionExpiryDateRequest({
    $core.String? locktoken,
    Scope? scope,
    $core.String? name,
    $core.String? versiontoexpire,
    $core.String? id,
    $core.String? expirytimestamp,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (scope != null) result.scope = scope;
    if (name != null) result.name = name;
    if (versiontoexpire != null) result.versiontoexpire = versiontoexpire;
    if (id != null) result.id = id;
    if (expirytimestamp != null) result.expirytimestamp = expirytimestamp;
    return result;
  }

  UpdateManagedRuleSetVersionExpiryDateRequest._();

  factory UpdateManagedRuleSetVersionExpiryDateRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateManagedRuleSetVersionExpiryDateRequest.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateManagedRuleSetVersionExpiryDateRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(272895286, _omitFieldNames ? '' : 'versiontoexpire')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(460460551, _omitFieldNames ? '' : 'expirytimestamp')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateManagedRuleSetVersionExpiryDateRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateManagedRuleSetVersionExpiryDateRequest copyWith(
          void Function(UpdateManagedRuleSetVersionExpiryDateRequest)
              updates) =>
      super.copyWith((message) =>
              updates(message as UpdateManagedRuleSetVersionExpiryDateRequest))
          as UpdateManagedRuleSetVersionExpiryDateRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateManagedRuleSetVersionExpiryDateRequest create() =>
      UpdateManagedRuleSetVersionExpiryDateRequest._();
  @$core.override
  UpdateManagedRuleSetVersionExpiryDateRequest createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static UpdateManagedRuleSetVersionExpiryDateRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateManagedRuleSetVersionExpiryDateRequest>(create);
  static UpdateManagedRuleSetVersionExpiryDateRequest? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(272895286)
  $core.String get versiontoexpire => $_getSZ(3);
  @$pb.TagNumber(272895286)
  set versiontoexpire($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272895286)
  $core.bool hasVersiontoexpire() => $_has(3);
  @$pb.TagNumber(272895286)
  void clearVersiontoexpire() => $_clearField(272895286);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(460460551)
  $core.String get expirytimestamp => $_getSZ(5);
  @$pb.TagNumber(460460551)
  set expirytimestamp($core.String value) => $_setString(5, value);
  @$pb.TagNumber(460460551)
  $core.bool hasExpirytimestamp() => $_has(5);
  @$pb.TagNumber(460460551)
  void clearExpirytimestamp() => $_clearField(460460551);
}

class UpdateManagedRuleSetVersionExpiryDateResponse
    extends $pb.GeneratedMessage {
  factory UpdateManagedRuleSetVersionExpiryDateResponse({
    $core.String? expiringversion,
    $core.String? expirytimestamp,
    $core.String? nextlocktoken,
  }) {
    final result = create();
    if (expiringversion != null) result.expiringversion = expiringversion;
    if (expirytimestamp != null) result.expirytimestamp = expirytimestamp;
    if (nextlocktoken != null) result.nextlocktoken = nextlocktoken;
    return result;
  }

  UpdateManagedRuleSetVersionExpiryDateResponse._();

  factory UpdateManagedRuleSetVersionExpiryDateResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateManagedRuleSetVersionExpiryDateResponse.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateManagedRuleSetVersionExpiryDateResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(282072504, _omitFieldNames ? '' : 'expiringversion')
    ..aOS(460460551, _omitFieldNames ? '' : 'expirytimestamp')
    ..aOS(473728437, _omitFieldNames ? '' : 'nextlocktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateManagedRuleSetVersionExpiryDateResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateManagedRuleSetVersionExpiryDateResponse copyWith(
          void Function(UpdateManagedRuleSetVersionExpiryDateResponse)
              updates) =>
      super.copyWith((message) =>
              updates(message as UpdateManagedRuleSetVersionExpiryDateResponse))
          as UpdateManagedRuleSetVersionExpiryDateResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateManagedRuleSetVersionExpiryDateResponse create() =>
      UpdateManagedRuleSetVersionExpiryDateResponse._();
  @$core.override
  UpdateManagedRuleSetVersionExpiryDateResponse createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static UpdateManagedRuleSetVersionExpiryDateResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateManagedRuleSetVersionExpiryDateResponse>(create);
  static UpdateManagedRuleSetVersionExpiryDateResponse? _defaultInstance;

  @$pb.TagNumber(282072504)
  $core.String get expiringversion => $_getSZ(0);
  @$pb.TagNumber(282072504)
  set expiringversion($core.String value) => $_setString(0, value);
  @$pb.TagNumber(282072504)
  $core.bool hasExpiringversion() => $_has(0);
  @$pb.TagNumber(282072504)
  void clearExpiringversion() => $_clearField(282072504);

  @$pb.TagNumber(460460551)
  $core.String get expirytimestamp => $_getSZ(1);
  @$pb.TagNumber(460460551)
  set expirytimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(460460551)
  $core.bool hasExpirytimestamp() => $_has(1);
  @$pb.TagNumber(460460551)
  void clearExpirytimestamp() => $_clearField(460460551);

  @$pb.TagNumber(473728437)
  $core.String get nextlocktoken => $_getSZ(2);
  @$pb.TagNumber(473728437)
  set nextlocktoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(473728437)
  $core.bool hasNextlocktoken() => $_has(2);
  @$pb.TagNumber(473728437)
  void clearNextlocktoken() => $_clearField(473728437);
}

class UpdateRegexPatternSetRequest extends $pb.GeneratedMessage {
  factory UpdateRegexPatternSetRequest({
    $core.String? locktoken,
    Scope? scope,
    $core.String? description,
    $core.Iterable<Regex>? regularexpressionlist,
    $core.String? name,
    $core.String? id,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (scope != null) result.scope = scope;
    if (description != null) result.description = description;
    if (regularexpressionlist != null)
      result.regularexpressionlist.addAll(regularexpressionlist);
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    return result;
  }

  UpdateRegexPatternSetRequest._();

  factory UpdateRegexPatternSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRegexPatternSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRegexPatternSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..pPM<Regex>(123612838, _omitFieldNames ? '' : 'regularexpressionlist',
        subBuilder: Regex.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRegexPatternSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRegexPatternSetRequest copyWith(
          void Function(UpdateRegexPatternSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateRegexPatternSetRequest))
          as UpdateRegexPatternSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRegexPatternSetRequest create() =>
      UpdateRegexPatternSetRequest._();
  @$core.override
  UpdateRegexPatternSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRegexPatternSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRegexPatternSetRequest>(create);
  static UpdateRegexPatternSetRequest? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(1);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(1);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(123612838)
  $pb.PbList<Regex> get regularexpressionlist => $_getList(3);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(5);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(5, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(5);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class UpdateRegexPatternSetResponse extends $pb.GeneratedMessage {
  factory UpdateRegexPatternSetResponse({
    $core.String? nextlocktoken,
  }) {
    final result = create();
    if (nextlocktoken != null) result.nextlocktoken = nextlocktoken;
    return result;
  }

  UpdateRegexPatternSetResponse._();

  factory UpdateRegexPatternSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRegexPatternSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRegexPatternSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(473728437, _omitFieldNames ? '' : 'nextlocktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRegexPatternSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRegexPatternSetResponse copyWith(
          void Function(UpdateRegexPatternSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateRegexPatternSetResponse))
          as UpdateRegexPatternSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRegexPatternSetResponse create() =>
      UpdateRegexPatternSetResponse._();
  @$core.override
  UpdateRegexPatternSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRegexPatternSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRegexPatternSetResponse>(create);
  static UpdateRegexPatternSetResponse? _defaultInstance;

  @$pb.TagNumber(473728437)
  $core.String get nextlocktoken => $_getSZ(0);
  @$pb.TagNumber(473728437)
  set nextlocktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(473728437)
  $core.bool hasNextlocktoken() => $_has(0);
  @$pb.TagNumber(473728437)
  void clearNextlocktoken() => $_clearField(473728437);
}

class UpdateRuleGroupRequest extends $pb.GeneratedMessage {
  factory UpdateRuleGroupRequest({
    $core.Iterable<Rule>? rules,
    $core.Iterable<$core.MapEntry<$core.String, CustomResponseBody>>?
        customresponsebodies,
    $core.String? locktoken,
    Scope? scope,
    $core.String? description,
    $core.String? name,
    $core.String? id,
    VisibilityConfig? visibilityconfig,
  }) {
    final result = create();
    if (rules != null) result.rules.addAll(rules);
    if (customresponsebodies != null)
      result.customresponsebodies.addEntries(customresponsebodies);
    if (locktoken != null) result.locktoken = locktoken;
    if (scope != null) result.scope = scope;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (visibilityconfig != null) result.visibilityconfig = visibilityconfig;
    return result;
  }

  UpdateRuleGroupRequest._();

  factory UpdateRuleGroupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRuleGroupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRuleGroupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<Rule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: Rule.create)
    ..m<$core.String, CustomResponseBody>(
        42731774, _omitFieldNames ? '' : 'customresponsebodies',
        entryClassName: 'UpdateRuleGroupRequest.CustomresponsebodiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: CustomResponseBody.create,
        valueDefaultOrMaker: CustomResponseBody.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOM<VisibilityConfig>(
        507762728, _omitFieldNames ? '' : 'visibilityconfig',
        subBuilder: VisibilityConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRuleGroupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRuleGroupRequest copyWith(
          void Function(UpdateRuleGroupRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateRuleGroupRequest))
          as UpdateRuleGroupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRuleGroupRequest create() => UpdateRuleGroupRequest._();
  @$core.override
  UpdateRuleGroupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRuleGroupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRuleGroupRequest>(create);
  static UpdateRuleGroupRequest? _defaultInstance;

  @$pb.TagNumber(42675585)
  $pb.PbList<Rule> get rules => $_getList(0);

  @$pb.TagNumber(42731774)
  $pb.PbMap<$core.String, CustomResponseBody> get customresponsebodies =>
      $_getMap(1);

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(2);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(2);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(3);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(3);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(6);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(6, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(6);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(507762728)
  VisibilityConfig get visibilityconfig => $_getN(7);
  @$pb.TagNumber(507762728)
  set visibilityconfig(VisibilityConfig value) => $_setField(507762728, value);
  @$pb.TagNumber(507762728)
  $core.bool hasVisibilityconfig() => $_has(7);
  @$pb.TagNumber(507762728)
  void clearVisibilityconfig() => $_clearField(507762728);
  @$pb.TagNumber(507762728)
  VisibilityConfig ensureVisibilityconfig() => $_ensure(7);
}

class UpdateRuleGroupResponse extends $pb.GeneratedMessage {
  factory UpdateRuleGroupResponse({
    $core.String? nextlocktoken,
  }) {
    final result = create();
    if (nextlocktoken != null) result.nextlocktoken = nextlocktoken;
    return result;
  }

  UpdateRuleGroupResponse._();

  factory UpdateRuleGroupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRuleGroupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRuleGroupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(473728437, _omitFieldNames ? '' : 'nextlocktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRuleGroupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRuleGroupResponse copyWith(
          void Function(UpdateRuleGroupResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateRuleGroupResponse))
          as UpdateRuleGroupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRuleGroupResponse create() => UpdateRuleGroupResponse._();
  @$core.override
  UpdateRuleGroupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRuleGroupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRuleGroupResponse>(create);
  static UpdateRuleGroupResponse? _defaultInstance;

  @$pb.TagNumber(473728437)
  $core.String get nextlocktoken => $_getSZ(0);
  @$pb.TagNumber(473728437)
  set nextlocktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(473728437)
  $core.bool hasNextlocktoken() => $_has(0);
  @$pb.TagNumber(473728437)
  void clearNextlocktoken() => $_clearField(473728437);
}

class UpdateWebACLRequest extends $pb.GeneratedMessage {
  factory UpdateWebACLRequest({
    $core.Iterable<$core.String>? tokendomains,
    $core.Iterable<Rule>? rules,
    $core.Iterable<$core.MapEntry<$core.String, CustomResponseBody>>?
        customresponsebodies,
    $core.String? locktoken,
    ChallengeConfig? challengeconfig,
    CaptchaConfig? captchaconfig,
    Scope? scope,
    OnSourceDDoSProtectionConfig? onsourceddosprotectionconfig,
    $core.String? description,
    $core.String? name,
    DefaultAction? defaultaction,
    $core.String? id,
    AssociationConfig? associationconfig,
    DataProtectionConfig? dataprotectionconfig,
    ApplicationConfig? applicationconfig,
    VisibilityConfig? visibilityconfig,
  }) {
    final result = create();
    if (tokendomains != null) result.tokendomains.addAll(tokendomains);
    if (rules != null) result.rules.addAll(rules);
    if (customresponsebodies != null)
      result.customresponsebodies.addEntries(customresponsebodies);
    if (locktoken != null) result.locktoken = locktoken;
    if (challengeconfig != null) result.challengeconfig = challengeconfig;
    if (captchaconfig != null) result.captchaconfig = captchaconfig;
    if (scope != null) result.scope = scope;
    if (onsourceddosprotectionconfig != null)
      result.onsourceddosprotectionconfig = onsourceddosprotectionconfig;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (defaultaction != null) result.defaultaction = defaultaction;
    if (id != null) result.id = id;
    if (associationconfig != null) result.associationconfig = associationconfig;
    if (dataprotectionconfig != null)
      result.dataprotectionconfig = dataprotectionconfig;
    if (applicationconfig != null) result.applicationconfig = applicationconfig;
    if (visibilityconfig != null) result.visibilityconfig = visibilityconfig;
    return result;
  }

  UpdateWebACLRequest._();

  factory UpdateWebACLRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateWebACLRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateWebACLRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(11638560, _omitFieldNames ? '' : 'tokendomains')
    ..pPM<Rule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: Rule.create)
    ..m<$core.String, CustomResponseBody>(
        42731774, _omitFieldNames ? '' : 'customresponsebodies',
        entryClassName: 'UpdateWebACLRequest.CustomresponsebodiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: CustomResponseBody.create,
        valueDefaultOrMaker: CustomResponseBody.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOM<ChallengeConfig>(48990889, _omitFieldNames ? '' : 'challengeconfig',
        subBuilder: ChallengeConfig.create)
    ..aOM<CaptchaConfig>(60547064, _omitFieldNames ? '' : 'captchaconfig',
        subBuilder: CaptchaConfig.create)
    ..aE<Scope>(65430924, _omitFieldNames ? '' : 'scope',
        enumValues: Scope.values)
    ..aOM<OnSourceDDoSProtectionConfig>(
        105229063, _omitFieldNames ? '' : 'onsourceddosprotectionconfig',
        subBuilder: OnSourceDDoSProtectionConfig.create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<DefaultAction>(322663861, _omitFieldNames ? '' : 'defaultaction',
        subBuilder: DefaultAction.create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOM<AssociationConfig>(
        419675691, _omitFieldNames ? '' : 'associationconfig',
        subBuilder: AssociationConfig.create)
    ..aOM<DataProtectionConfig>(
        464792245, _omitFieldNames ? '' : 'dataprotectionconfig',
        subBuilder: DataProtectionConfig.create)
    ..aOM<ApplicationConfig>(
        501131976, _omitFieldNames ? '' : 'applicationconfig',
        subBuilder: ApplicationConfig.create)
    ..aOM<VisibilityConfig>(
        507762728, _omitFieldNames ? '' : 'visibilityconfig',
        subBuilder: VisibilityConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateWebACLRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateWebACLRequest copyWith(void Function(UpdateWebACLRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateWebACLRequest))
          as UpdateWebACLRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateWebACLRequest create() => UpdateWebACLRequest._();
  @$core.override
  UpdateWebACLRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateWebACLRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateWebACLRequest>(create);
  static UpdateWebACLRequest? _defaultInstance;

  @$pb.TagNumber(11638560)
  $pb.PbList<$core.String> get tokendomains => $_getList(0);

  @$pb.TagNumber(42675585)
  $pb.PbList<Rule> get rules => $_getList(1);

  @$pb.TagNumber(42731774)
  $pb.PbMap<$core.String, CustomResponseBody> get customresponsebodies =>
      $_getMap(2);

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(3);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(3, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(3);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(48990889)
  ChallengeConfig get challengeconfig => $_getN(4);
  @$pb.TagNumber(48990889)
  set challengeconfig(ChallengeConfig value) => $_setField(48990889, value);
  @$pb.TagNumber(48990889)
  $core.bool hasChallengeconfig() => $_has(4);
  @$pb.TagNumber(48990889)
  void clearChallengeconfig() => $_clearField(48990889);
  @$pb.TagNumber(48990889)
  ChallengeConfig ensureChallengeconfig() => $_ensure(4);

  @$pb.TagNumber(60547064)
  CaptchaConfig get captchaconfig => $_getN(5);
  @$pb.TagNumber(60547064)
  set captchaconfig(CaptchaConfig value) => $_setField(60547064, value);
  @$pb.TagNumber(60547064)
  $core.bool hasCaptchaconfig() => $_has(5);
  @$pb.TagNumber(60547064)
  void clearCaptchaconfig() => $_clearField(60547064);
  @$pb.TagNumber(60547064)
  CaptchaConfig ensureCaptchaconfig() => $_ensure(5);

  @$pb.TagNumber(65430924)
  Scope get scope => $_getN(6);
  @$pb.TagNumber(65430924)
  set scope(Scope value) => $_setField(65430924, value);
  @$pb.TagNumber(65430924)
  $core.bool hasScope() => $_has(6);
  @$pb.TagNumber(65430924)
  void clearScope() => $_clearField(65430924);

  @$pb.TagNumber(105229063)
  OnSourceDDoSProtectionConfig get onsourceddosprotectionconfig => $_getN(7);
  @$pb.TagNumber(105229063)
  set onsourceddosprotectionconfig(OnSourceDDoSProtectionConfig value) =>
      $_setField(105229063, value);
  @$pb.TagNumber(105229063)
  $core.bool hasOnsourceddosprotectionconfig() => $_has(7);
  @$pb.TagNumber(105229063)
  void clearOnsourceddosprotectionconfig() => $_clearField(105229063);
  @$pb.TagNumber(105229063)
  OnSourceDDoSProtectionConfig ensureOnsourceddosprotectionconfig() =>
      $_ensure(7);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(8);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(8, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(8);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(9);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(9, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(9);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(322663861)
  DefaultAction get defaultaction => $_getN(10);
  @$pb.TagNumber(322663861)
  set defaultaction(DefaultAction value) => $_setField(322663861, value);
  @$pb.TagNumber(322663861)
  $core.bool hasDefaultaction() => $_has(10);
  @$pb.TagNumber(322663861)
  void clearDefaultaction() => $_clearField(322663861);
  @$pb.TagNumber(322663861)
  DefaultAction ensureDefaultaction() => $_ensure(10);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(11);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(11, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(11);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(419675691)
  AssociationConfig get associationconfig => $_getN(12);
  @$pb.TagNumber(419675691)
  set associationconfig(AssociationConfig value) =>
      $_setField(419675691, value);
  @$pb.TagNumber(419675691)
  $core.bool hasAssociationconfig() => $_has(12);
  @$pb.TagNumber(419675691)
  void clearAssociationconfig() => $_clearField(419675691);
  @$pb.TagNumber(419675691)
  AssociationConfig ensureAssociationconfig() => $_ensure(12);

  @$pb.TagNumber(464792245)
  DataProtectionConfig get dataprotectionconfig => $_getN(13);
  @$pb.TagNumber(464792245)
  set dataprotectionconfig(DataProtectionConfig value) =>
      $_setField(464792245, value);
  @$pb.TagNumber(464792245)
  $core.bool hasDataprotectionconfig() => $_has(13);
  @$pb.TagNumber(464792245)
  void clearDataprotectionconfig() => $_clearField(464792245);
  @$pb.TagNumber(464792245)
  DataProtectionConfig ensureDataprotectionconfig() => $_ensure(13);

  @$pb.TagNumber(501131976)
  ApplicationConfig get applicationconfig => $_getN(14);
  @$pb.TagNumber(501131976)
  set applicationconfig(ApplicationConfig value) =>
      $_setField(501131976, value);
  @$pb.TagNumber(501131976)
  $core.bool hasApplicationconfig() => $_has(14);
  @$pb.TagNumber(501131976)
  void clearApplicationconfig() => $_clearField(501131976);
  @$pb.TagNumber(501131976)
  ApplicationConfig ensureApplicationconfig() => $_ensure(14);

  @$pb.TagNumber(507762728)
  VisibilityConfig get visibilityconfig => $_getN(15);
  @$pb.TagNumber(507762728)
  set visibilityconfig(VisibilityConfig value) => $_setField(507762728, value);
  @$pb.TagNumber(507762728)
  $core.bool hasVisibilityconfig() => $_has(15);
  @$pb.TagNumber(507762728)
  void clearVisibilityconfig() => $_clearField(507762728);
  @$pb.TagNumber(507762728)
  VisibilityConfig ensureVisibilityconfig() => $_ensure(15);
}

class UpdateWebACLResponse extends $pb.GeneratedMessage {
  factory UpdateWebACLResponse({
    $core.String? nextlocktoken,
  }) {
    final result = create();
    if (nextlocktoken != null) result.nextlocktoken = nextlocktoken;
    return result;
  }

  UpdateWebACLResponse._();

  factory UpdateWebACLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateWebACLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateWebACLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(473728437, _omitFieldNames ? '' : 'nextlocktoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateWebACLResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateWebACLResponse copyWith(void Function(UpdateWebACLResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateWebACLResponse))
          as UpdateWebACLResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateWebACLResponse create() => UpdateWebACLResponse._();
  @$core.override
  UpdateWebACLResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateWebACLResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateWebACLResponse>(create);
  static UpdateWebACLResponse? _defaultInstance;

  @$pb.TagNumber(473728437)
  $core.String get nextlocktoken => $_getSZ(0);
  @$pb.TagNumber(473728437)
  set nextlocktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(473728437)
  $core.bool hasNextlocktoken() => $_has(0);
  @$pb.TagNumber(473728437)
  void clearNextlocktoken() => $_clearField(473728437);
}

class UriFragment extends $pb.GeneratedMessage {
  factory UriFragment({
    FallbackBehavior? fallbackbehavior,
  }) {
    final result = create();
    if (fallbackbehavior != null) result.fallbackbehavior = fallbackbehavior;
    return result;
  }

  UriFragment._();

  factory UriFragment.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UriFragment.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UriFragment',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aE<FallbackBehavior>(440114542, _omitFieldNames ? '' : 'fallbackbehavior',
        enumValues: FallbackBehavior.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UriFragment clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UriFragment copyWith(void Function(UriFragment) updates) =>
      super.copyWith((message) => updates(message as UriFragment))
          as UriFragment;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UriFragment create() => UriFragment._();
  @$core.override
  UriFragment createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UriFragment getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UriFragment>(create);
  static UriFragment? _defaultInstance;

  @$pb.TagNumber(440114542)
  FallbackBehavior get fallbackbehavior => $_getN(0);
  @$pb.TagNumber(440114542)
  set fallbackbehavior(FallbackBehavior value) => $_setField(440114542, value);
  @$pb.TagNumber(440114542)
  $core.bool hasFallbackbehavior() => $_has(0);
  @$pb.TagNumber(440114542)
  void clearFallbackbehavior() => $_clearField(440114542);
}

class UriPath extends $pb.GeneratedMessage {
  factory UriPath() => create();

  UriPath._();

  factory UriPath.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UriPath.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UriPath',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UriPath clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UriPath copyWith(void Function(UriPath) updates) =>
      super.copyWith((message) => updates(message as UriPath)) as UriPath;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UriPath create() => UriPath._();
  @$core.override
  UriPath createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UriPath getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UriPath>(create);
  static UriPath? _defaultInstance;
}

class UsernameField extends $pb.GeneratedMessage {
  factory UsernameField({
    $core.String? identifier,
  }) {
    final result = create();
    if (identifier != null) result.identifier = identifier;
    return result;
  }

  UsernameField._();

  factory UsernameField.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UsernameField.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UsernameField',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(41865311, _omitFieldNames ? '' : 'identifier')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsernameField clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UsernameField copyWith(void Function(UsernameField) updates) =>
      super.copyWith((message) => updates(message as UsernameField))
          as UsernameField;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UsernameField create() => UsernameField._();
  @$core.override
  UsernameField createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UsernameField getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UsernameField>(create);
  static UsernameField? _defaultInstance;

  @$pb.TagNumber(41865311)
  $core.String get identifier => $_getSZ(0);
  @$pb.TagNumber(41865311)
  set identifier($core.String value) => $_setString(0, value);
  @$pb.TagNumber(41865311)
  $core.bool hasIdentifier() => $_has(0);
  @$pb.TagNumber(41865311)
  void clearIdentifier() => $_clearField(41865311);
}

class VersionToPublish extends $pb.GeneratedMessage {
  factory VersionToPublish({
    $core.int? forecastedlifetime,
    $core.String? associatedrulegrouparn,
  }) {
    final result = create();
    if (forecastedlifetime != null)
      result.forecastedlifetime = forecastedlifetime;
    if (associatedrulegrouparn != null)
      result.associatedrulegrouparn = associatedrulegrouparn;
    return result;
  }

  VersionToPublish._();

  factory VersionToPublish.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VersionToPublish.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VersionToPublish',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aI(149179979, _omitFieldNames ? '' : 'forecastedlifetime')
    ..aOS(416522796, _omitFieldNames ? '' : 'associatedrulegrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VersionToPublish clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VersionToPublish copyWith(void Function(VersionToPublish) updates) =>
      super.copyWith((message) => updates(message as VersionToPublish))
          as VersionToPublish;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VersionToPublish create() => VersionToPublish._();
  @$core.override
  VersionToPublish createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VersionToPublish getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VersionToPublish>(create);
  static VersionToPublish? _defaultInstance;

  @$pb.TagNumber(149179979)
  $core.int get forecastedlifetime => $_getIZ(0);
  @$pb.TagNumber(149179979)
  set forecastedlifetime($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(149179979)
  $core.bool hasForecastedlifetime() => $_has(0);
  @$pb.TagNumber(149179979)
  void clearForecastedlifetime() => $_clearField(149179979);

  @$pb.TagNumber(416522796)
  $core.String get associatedrulegrouparn => $_getSZ(1);
  @$pb.TagNumber(416522796)
  set associatedrulegrouparn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(416522796)
  $core.bool hasAssociatedrulegrouparn() => $_has(1);
  @$pb.TagNumber(416522796)
  void clearAssociatedrulegrouparn() => $_clearField(416522796);
}

class VisibilityConfig extends $pb.GeneratedMessage {
  factory VisibilityConfig({
    $core.String? metricname,
    $core.bool? cloudwatchmetricsenabled,
    $core.bool? sampledrequestsenabled,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (cloudwatchmetricsenabled != null)
      result.cloudwatchmetricsenabled = cloudwatchmetricsenabled;
    if (sampledrequestsenabled != null)
      result.sampledrequestsenabled = sampledrequestsenabled;
    return result;
  }

  VisibilityConfig._();

  factory VisibilityConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VisibilityConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VisibilityConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOB(387050620, _omitFieldNames ? '' : 'cloudwatchmetricsenabled')
    ..aOB(393409933, _omitFieldNames ? '' : 'sampledrequestsenabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VisibilityConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VisibilityConfig copyWith(void Function(VisibilityConfig) updates) =>
      super.copyWith((message) => updates(message as VisibilityConfig))
          as VisibilityConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VisibilityConfig create() => VisibilityConfig._();
  @$core.override
  VisibilityConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VisibilityConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VisibilityConfig>(create);
  static VisibilityConfig? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(387050620)
  $core.bool get cloudwatchmetricsenabled => $_getBF(1);
  @$pb.TagNumber(387050620)
  set cloudwatchmetricsenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(387050620)
  $core.bool hasCloudwatchmetricsenabled() => $_has(1);
  @$pb.TagNumber(387050620)
  void clearCloudwatchmetricsenabled() => $_clearField(387050620);

  @$pb.TagNumber(393409933)
  $core.bool get sampledrequestsenabled => $_getBF(2);
  @$pb.TagNumber(393409933)
  set sampledrequestsenabled($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(393409933)
  $core.bool hasSampledrequestsenabled() => $_has(2);
  @$pb.TagNumber(393409933)
  void clearSampledrequestsenabled() => $_clearField(393409933);
}

class WAFAssociatedItemException extends $pb.GeneratedMessage {
  factory WAFAssociatedItemException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFAssociatedItemException._();

  factory WAFAssociatedItemException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFAssociatedItemException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFAssociatedItemException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFAssociatedItemException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFAssociatedItemException copyWith(
          void Function(WAFAssociatedItemException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFAssociatedItemException))
          as WAFAssociatedItemException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFAssociatedItemException create() => WAFAssociatedItemException._();
  @$core.override
  WAFAssociatedItemException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFAssociatedItemException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFAssociatedItemException>(create);
  static WAFAssociatedItemException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFConfigurationWarningException extends $pb.GeneratedMessage {
  factory WAFConfigurationWarningException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFConfigurationWarningException._();

  factory WAFConfigurationWarningException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFConfigurationWarningException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFConfigurationWarningException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFConfigurationWarningException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFConfigurationWarningException copyWith(
          void Function(WAFConfigurationWarningException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFConfigurationWarningException))
          as WAFConfigurationWarningException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFConfigurationWarningException create() =>
      WAFConfigurationWarningException._();
  @$core.override
  WAFConfigurationWarningException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFConfigurationWarningException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFConfigurationWarningException>(
          create);
  static WAFConfigurationWarningException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFDuplicateItemException extends $pb.GeneratedMessage {
  factory WAFDuplicateItemException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFDuplicateItemException._();

  factory WAFDuplicateItemException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFDuplicateItemException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFDuplicateItemException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFDuplicateItemException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFDuplicateItemException copyWith(
          void Function(WAFDuplicateItemException) updates) =>
      super.copyWith((message) => updates(message as WAFDuplicateItemException))
          as WAFDuplicateItemException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFDuplicateItemException create() => WAFDuplicateItemException._();
  @$core.override
  WAFDuplicateItemException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFDuplicateItemException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFDuplicateItemException>(create);
  static WAFDuplicateItemException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFExpiredManagedRuleGroupVersionException extends $pb.GeneratedMessage {
  factory WAFExpiredManagedRuleGroupVersionException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFExpiredManagedRuleGroupVersionException._();

  factory WAFExpiredManagedRuleGroupVersionException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFExpiredManagedRuleGroupVersionException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFExpiredManagedRuleGroupVersionException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFExpiredManagedRuleGroupVersionException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFExpiredManagedRuleGroupVersionException copyWith(
          void Function(WAFExpiredManagedRuleGroupVersionException) updates) =>
      super.copyWith((message) =>
              updates(message as WAFExpiredManagedRuleGroupVersionException))
          as WAFExpiredManagedRuleGroupVersionException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFExpiredManagedRuleGroupVersionException create() =>
      WAFExpiredManagedRuleGroupVersionException._();
  @$core.override
  WAFExpiredManagedRuleGroupVersionException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFExpiredManagedRuleGroupVersionException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          WAFExpiredManagedRuleGroupVersionException>(create);
  static WAFExpiredManagedRuleGroupVersionException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFFeatureNotIncludedInPricingPlanException extends $pb.GeneratedMessage {
  factory WAFFeatureNotIncludedInPricingPlanException({
    $core.String? message,
    $core.Iterable<DisallowedFeature>? disallowedfeatures,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (disallowedfeatures != null)
      result.disallowedfeatures.addAll(disallowedfeatures);
    return result;
  }

  WAFFeatureNotIncludedInPricingPlanException._();

  factory WAFFeatureNotIncludedInPricingPlanException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFFeatureNotIncludedInPricingPlanException.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFFeatureNotIncludedInPricingPlanException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..pPM<DisallowedFeature>(
        445316575, _omitFieldNames ? '' : 'disallowedfeatures',
        subBuilder: DisallowedFeature.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFFeatureNotIncludedInPricingPlanException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFFeatureNotIncludedInPricingPlanException copyWith(
          void Function(WAFFeatureNotIncludedInPricingPlanException) updates) =>
      super.copyWith((message) =>
              updates(message as WAFFeatureNotIncludedInPricingPlanException))
          as WAFFeatureNotIncludedInPricingPlanException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFFeatureNotIncludedInPricingPlanException create() =>
      WAFFeatureNotIncludedInPricingPlanException._();
  @$core.override
  WAFFeatureNotIncludedInPricingPlanException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFFeatureNotIncludedInPricingPlanException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          WAFFeatureNotIncludedInPricingPlanException>(create);
  static WAFFeatureNotIncludedInPricingPlanException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(445316575)
  $pb.PbList<DisallowedFeature> get disallowedfeatures => $_getList(1);
}

class WAFInternalErrorException extends $pb.GeneratedMessage {
  factory WAFInternalErrorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFInternalErrorException._();

  factory WAFInternalErrorException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFInternalErrorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFInternalErrorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInternalErrorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInternalErrorException copyWith(
          void Function(WAFInternalErrorException) updates) =>
      super.copyWith((message) => updates(message as WAFInternalErrorException))
          as WAFInternalErrorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFInternalErrorException create() => WAFInternalErrorException._();
  @$core.override
  WAFInternalErrorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFInternalErrorException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFInternalErrorException>(create);
  static WAFInternalErrorException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFInvalidOperationException extends $pb.GeneratedMessage {
  factory WAFInvalidOperationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFInvalidOperationException._();

  factory WAFInvalidOperationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFInvalidOperationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFInvalidOperationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidOperationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidOperationException copyWith(
          void Function(WAFInvalidOperationException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFInvalidOperationException))
          as WAFInvalidOperationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFInvalidOperationException create() =>
      WAFInvalidOperationException._();
  @$core.override
  WAFInvalidOperationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFInvalidOperationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFInvalidOperationException>(create);
  static WAFInvalidOperationException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFInvalidParameterException extends $pb.GeneratedMessage {
  factory WAFInvalidParameterException({
    $core.String? reason,
    $core.String? message,
    ParameterExceptionField? field_263732488,
    $core.String? parameter,
  }) {
    final result = create();
    if (reason != null) result.reason = reason;
    if (message != null) result.message = message;
    if (field_263732488 != null) result.field_263732488 = field_263732488;
    if (parameter != null) result.parameter = parameter;
    return result;
  }

  WAFInvalidParameterException._();

  factory WAFInvalidParameterException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFInvalidParameterException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFInvalidParameterException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(20005178, _omitFieldNames ? '' : 'reason')
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..aE<ParameterExceptionField>(263732488, _omitFieldNames ? '' : 'field',
        enumValues: ParameterExceptionField.values)
    ..aOS(407419825, _omitFieldNames ? '' : 'parameter')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidParameterException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidParameterException copyWith(
          void Function(WAFInvalidParameterException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFInvalidParameterException))
          as WAFInvalidParameterException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFInvalidParameterException create() =>
      WAFInvalidParameterException._();
  @$core.override
  WAFInvalidParameterException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFInvalidParameterException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFInvalidParameterException>(create);
  static WAFInvalidParameterException? _defaultInstance;

  @$pb.TagNumber(20005178)
  $core.String get reason => $_getSZ(0);
  @$pb.TagNumber(20005178)
  set reason($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20005178)
  $core.bool hasReason() => $_has(0);
  @$pb.TagNumber(20005178)
  void clearReason() => $_clearField(20005178);

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(263732488)
  ParameterExceptionField get field_263732488 => $_getN(2);
  @$pb.TagNumber(263732488)
  set field_263732488(ParameterExceptionField value) =>
      $_setField(263732488, value);
  @$pb.TagNumber(263732488)
  $core.bool hasField_263732488() => $_has(2);
  @$pb.TagNumber(263732488)
  void clearField_263732488() => $_clearField(263732488);

  @$pb.TagNumber(407419825)
  $core.String get parameter => $_getSZ(3);
  @$pb.TagNumber(407419825)
  set parameter($core.String value) => $_setString(3, value);
  @$pb.TagNumber(407419825)
  $core.bool hasParameter() => $_has(3);
  @$pb.TagNumber(407419825)
  void clearParameter() => $_clearField(407419825);
}

class WAFInvalidPermissionPolicyException extends $pb.GeneratedMessage {
  factory WAFInvalidPermissionPolicyException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFInvalidPermissionPolicyException._();

  factory WAFInvalidPermissionPolicyException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFInvalidPermissionPolicyException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFInvalidPermissionPolicyException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidPermissionPolicyException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidPermissionPolicyException copyWith(
          void Function(WAFInvalidPermissionPolicyException) updates) =>
      super.copyWith((message) =>
              updates(message as WAFInvalidPermissionPolicyException))
          as WAFInvalidPermissionPolicyException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFInvalidPermissionPolicyException create() =>
      WAFInvalidPermissionPolicyException._();
  @$core.override
  WAFInvalidPermissionPolicyException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFInvalidPermissionPolicyException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          WAFInvalidPermissionPolicyException>(create);
  static WAFInvalidPermissionPolicyException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFInvalidResourceException extends $pb.GeneratedMessage {
  factory WAFInvalidResourceException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFInvalidResourceException._();

  factory WAFInvalidResourceException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFInvalidResourceException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFInvalidResourceException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidResourceException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidResourceException copyWith(
          void Function(WAFInvalidResourceException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFInvalidResourceException))
          as WAFInvalidResourceException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFInvalidResourceException create() =>
      WAFInvalidResourceException._();
  @$core.override
  WAFInvalidResourceException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFInvalidResourceException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFInvalidResourceException>(create);
  static WAFInvalidResourceException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFLimitsExceededException extends $pb.GeneratedMessage {
  factory WAFLimitsExceededException({
    $core.String? sourcetype,
    $core.String? message,
  }) {
    final result = create();
    if (sourcetype != null) result.sourcetype = sourcetype;
    if (message != null) result.message = message;
    return result;
  }

  WAFLimitsExceededException._();

  factory WAFLimitsExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFLimitsExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFLimitsExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(195731217, _omitFieldNames ? '' : 'sourcetype')
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFLimitsExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFLimitsExceededException copyWith(
          void Function(WAFLimitsExceededException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFLimitsExceededException))
          as WAFLimitsExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFLimitsExceededException create() => WAFLimitsExceededException._();
  @$core.override
  WAFLimitsExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFLimitsExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFLimitsExceededException>(create);
  static WAFLimitsExceededException? _defaultInstance;

  @$pb.TagNumber(195731217)
  $core.String get sourcetype => $_getSZ(0);
  @$pb.TagNumber(195731217)
  set sourcetype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(195731217)
  $core.bool hasSourcetype() => $_has(0);
  @$pb.TagNumber(195731217)
  void clearSourcetype() => $_clearField(195731217);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFLogDestinationPermissionIssueException extends $pb.GeneratedMessage {
  factory WAFLogDestinationPermissionIssueException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFLogDestinationPermissionIssueException._();

  factory WAFLogDestinationPermissionIssueException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFLogDestinationPermissionIssueException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFLogDestinationPermissionIssueException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFLogDestinationPermissionIssueException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFLogDestinationPermissionIssueException copyWith(
          void Function(WAFLogDestinationPermissionIssueException) updates) =>
      super.copyWith((message) =>
              updates(message as WAFLogDestinationPermissionIssueException))
          as WAFLogDestinationPermissionIssueException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFLogDestinationPermissionIssueException create() =>
      WAFLogDestinationPermissionIssueException._();
  @$core.override
  WAFLogDestinationPermissionIssueException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFLogDestinationPermissionIssueException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          WAFLogDestinationPermissionIssueException>(create);
  static WAFLogDestinationPermissionIssueException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFNonexistentItemException extends $pb.GeneratedMessage {
  factory WAFNonexistentItemException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFNonexistentItemException._();

  factory WAFNonexistentItemException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFNonexistentItemException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFNonexistentItemException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFNonexistentItemException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFNonexistentItemException copyWith(
          void Function(WAFNonexistentItemException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFNonexistentItemException))
          as WAFNonexistentItemException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFNonexistentItemException create() =>
      WAFNonexistentItemException._();
  @$core.override
  WAFNonexistentItemException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFNonexistentItemException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFNonexistentItemException>(create);
  static WAFNonexistentItemException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFOptimisticLockException extends $pb.GeneratedMessage {
  factory WAFOptimisticLockException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFOptimisticLockException._();

  factory WAFOptimisticLockException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFOptimisticLockException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFOptimisticLockException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFOptimisticLockException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFOptimisticLockException copyWith(
          void Function(WAFOptimisticLockException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFOptimisticLockException))
          as WAFOptimisticLockException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFOptimisticLockException create() => WAFOptimisticLockException._();
  @$core.override
  WAFOptimisticLockException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFOptimisticLockException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFOptimisticLockException>(create);
  static WAFOptimisticLockException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFServiceLinkedRoleErrorException extends $pb.GeneratedMessage {
  factory WAFServiceLinkedRoleErrorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFServiceLinkedRoleErrorException._();

  factory WAFServiceLinkedRoleErrorException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFServiceLinkedRoleErrorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFServiceLinkedRoleErrorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFServiceLinkedRoleErrorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFServiceLinkedRoleErrorException copyWith(
          void Function(WAFServiceLinkedRoleErrorException) updates) =>
      super.copyWith((message) =>
              updates(message as WAFServiceLinkedRoleErrorException))
          as WAFServiceLinkedRoleErrorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFServiceLinkedRoleErrorException create() =>
      WAFServiceLinkedRoleErrorException._();
  @$core.override
  WAFServiceLinkedRoleErrorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFServiceLinkedRoleErrorException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFServiceLinkedRoleErrorException>(
          create);
  static WAFServiceLinkedRoleErrorException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFSubscriptionNotFoundException extends $pb.GeneratedMessage {
  factory WAFSubscriptionNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFSubscriptionNotFoundException._();

  factory WAFSubscriptionNotFoundException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFSubscriptionNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFSubscriptionNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFSubscriptionNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFSubscriptionNotFoundException copyWith(
          void Function(WAFSubscriptionNotFoundException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFSubscriptionNotFoundException))
          as WAFSubscriptionNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFSubscriptionNotFoundException create() =>
      WAFSubscriptionNotFoundException._();
  @$core.override
  WAFSubscriptionNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFSubscriptionNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFSubscriptionNotFoundException>(
          create);
  static WAFSubscriptionNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFTagOperationException extends $pb.GeneratedMessage {
  factory WAFTagOperationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFTagOperationException._();

  factory WAFTagOperationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFTagOperationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFTagOperationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFTagOperationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFTagOperationException copyWith(
          void Function(WAFTagOperationException) updates) =>
      super.copyWith((message) => updates(message as WAFTagOperationException))
          as WAFTagOperationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFTagOperationException create() => WAFTagOperationException._();
  @$core.override
  WAFTagOperationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFTagOperationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFTagOperationException>(create);
  static WAFTagOperationException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFTagOperationInternalErrorException extends $pb.GeneratedMessage {
  factory WAFTagOperationInternalErrorException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFTagOperationInternalErrorException._();

  factory WAFTagOperationInternalErrorException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFTagOperationInternalErrorException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFTagOperationInternalErrorException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFTagOperationInternalErrorException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFTagOperationInternalErrorException copyWith(
          void Function(WAFTagOperationInternalErrorException) updates) =>
      super.copyWith((message) =>
              updates(message as WAFTagOperationInternalErrorException))
          as WAFTagOperationInternalErrorException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFTagOperationInternalErrorException create() =>
      WAFTagOperationInternalErrorException._();
  @$core.override
  WAFTagOperationInternalErrorException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFTagOperationInternalErrorException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          WAFTagOperationInternalErrorException>(create);
  static WAFTagOperationInternalErrorException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFUnavailableEntityException extends $pb.GeneratedMessage {
  factory WAFUnavailableEntityException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFUnavailableEntityException._();

  factory WAFUnavailableEntityException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFUnavailableEntityException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFUnavailableEntityException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFUnavailableEntityException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFUnavailableEntityException copyWith(
          void Function(WAFUnavailableEntityException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFUnavailableEntityException))
          as WAFUnavailableEntityException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFUnavailableEntityException create() =>
      WAFUnavailableEntityException._();
  @$core.override
  WAFUnavailableEntityException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFUnavailableEntityException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFUnavailableEntityException>(create);
  static WAFUnavailableEntityException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WAFUnsupportedAggregateKeyTypeException extends $pb.GeneratedMessage {
  factory WAFUnsupportedAggregateKeyTypeException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFUnsupportedAggregateKeyTypeException._();

  factory WAFUnsupportedAggregateKeyTypeException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFUnsupportedAggregateKeyTypeException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFUnsupportedAggregateKeyTypeException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFUnsupportedAggregateKeyTypeException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFUnsupportedAggregateKeyTypeException copyWith(
          void Function(WAFUnsupportedAggregateKeyTypeException) updates) =>
      super.copyWith((message) =>
              updates(message as WAFUnsupportedAggregateKeyTypeException))
          as WAFUnsupportedAggregateKeyTypeException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFUnsupportedAggregateKeyTypeException create() =>
      WAFUnsupportedAggregateKeyTypeException._();
  @$core.override
  WAFUnsupportedAggregateKeyTypeException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFUnsupportedAggregateKeyTypeException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          WAFUnsupportedAggregateKeyTypeException>(create);
  static WAFUnsupportedAggregateKeyTypeException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class WebACL extends $pb.GeneratedMessage {
  factory WebACL({
    $core.Iterable<$core.String>? tokendomains,
    $core.String? labelnamespace,
    $core.Iterable<Rule>? rules,
    $core.Iterable<$core.MapEntry<$core.String, CustomResponseBody>>?
        customresponsebodies,
    ChallengeConfig? challengeconfig,
    CaptchaConfig? captchaconfig,
    OnSourceDDoSProtectionConfig? onsourceddosprotectionconfig,
    $fixnum.Int64? capacity,
    $core.String? description,
    $core.String? name,
    $core.bool? retrofittedbyfirewallmanager,
    DefaultAction? defaultaction,
    $core.Iterable<FirewallManagerRuleGroup>?
        preprocessfirewallmanagerrulegroups,
    $core.String? id,
    $core.String? arn,
    $core.bool? managedbyfirewallmanager,
    AssociationConfig? associationconfig,
    DataProtectionConfig? dataprotectionconfig,
    $core.Iterable<FirewallManagerRuleGroup>?
        postprocessfirewallmanagerrulegroups,
    ApplicationConfig? applicationconfig,
    VisibilityConfig? visibilityconfig,
  }) {
    final result = create();
    if (tokendomains != null) result.tokendomains.addAll(tokendomains);
    if (labelnamespace != null) result.labelnamespace = labelnamespace;
    if (rules != null) result.rules.addAll(rules);
    if (customresponsebodies != null)
      result.customresponsebodies.addEntries(customresponsebodies);
    if (challengeconfig != null) result.challengeconfig = challengeconfig;
    if (captchaconfig != null) result.captchaconfig = captchaconfig;
    if (onsourceddosprotectionconfig != null)
      result.onsourceddosprotectionconfig = onsourceddosprotectionconfig;
    if (capacity != null) result.capacity = capacity;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (retrofittedbyfirewallmanager != null)
      result.retrofittedbyfirewallmanager = retrofittedbyfirewallmanager;
    if (defaultaction != null) result.defaultaction = defaultaction;
    if (preprocessfirewallmanagerrulegroups != null)
      result.preprocessfirewallmanagerrulegroups
          .addAll(preprocessfirewallmanagerrulegroups);
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    if (managedbyfirewallmanager != null)
      result.managedbyfirewallmanager = managedbyfirewallmanager;
    if (associationconfig != null) result.associationconfig = associationconfig;
    if (dataprotectionconfig != null)
      result.dataprotectionconfig = dataprotectionconfig;
    if (postprocessfirewallmanagerrulegroups != null)
      result.postprocessfirewallmanagerrulegroups
          .addAll(postprocessfirewallmanagerrulegroups);
    if (applicationconfig != null) result.applicationconfig = applicationconfig;
    if (visibilityconfig != null) result.visibilityconfig = visibilityconfig;
    return result;
  }

  WebACL._();

  factory WebACL.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WebACL.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WebACL',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPS(11638560, _omitFieldNames ? '' : 'tokendomains')
    ..aOS(39797417, _omitFieldNames ? '' : 'labelnamespace')
    ..pPM<Rule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: Rule.create)
    ..m<$core.String, CustomResponseBody>(
        42731774, _omitFieldNames ? '' : 'customresponsebodies',
        entryClassName: 'WebACL.CustomresponsebodiesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: CustomResponseBody.create,
        valueDefaultOrMaker: CustomResponseBody.getDefault,
        packageName: const $pb.PackageName('wafv2'))
    ..aOM<ChallengeConfig>(48990889, _omitFieldNames ? '' : 'challengeconfig',
        subBuilder: ChallengeConfig.create)
    ..aOM<CaptchaConfig>(60547064, _omitFieldNames ? '' : 'captchaconfig',
        subBuilder: CaptchaConfig.create)
    ..aOM<OnSourceDDoSProtectionConfig>(
        105229063, _omitFieldNames ? '' : 'onsourceddosprotectionconfig',
        subBuilder: OnSourceDDoSProtectionConfig.create)
    ..aInt64(107253930, _omitFieldNames ? '' : 'capacity')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOB(279052502, _omitFieldNames ? '' : 'retrofittedbyfirewallmanager')
    ..aOM<DefaultAction>(322663861, _omitFieldNames ? '' : 'defaultaction',
        subBuilder: DefaultAction.create)
    ..pPM<FirewallManagerRuleGroup>(
        331512685, _omitFieldNames ? '' : 'preprocessfirewallmanagerrulegroups',
        subBuilder: FirewallManagerRuleGroup.create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..aOB(415792743, _omitFieldNames ? '' : 'managedbyfirewallmanager')
    ..aOM<AssociationConfig>(
        419675691, _omitFieldNames ? '' : 'associationconfig',
        subBuilder: AssociationConfig.create)
    ..aOM<DataProtectionConfig>(
        464792245, _omitFieldNames ? '' : 'dataprotectionconfig',
        subBuilder: DataProtectionConfig.create)
    ..pPM<FirewallManagerRuleGroup>(494837446,
        _omitFieldNames ? '' : 'postprocessfirewallmanagerrulegroups',
        subBuilder: FirewallManagerRuleGroup.create)
    ..aOM<ApplicationConfig>(
        501131976, _omitFieldNames ? '' : 'applicationconfig',
        subBuilder: ApplicationConfig.create)
    ..aOM<VisibilityConfig>(
        507762728, _omitFieldNames ? '' : 'visibilityconfig',
        subBuilder: VisibilityConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WebACL clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WebACL copyWith(void Function(WebACL) updates) =>
      super.copyWith((message) => updates(message as WebACL)) as WebACL;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WebACL create() => WebACL._();
  @$core.override
  WebACL createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WebACL getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<WebACL>(create);
  static WebACL? _defaultInstance;

  @$pb.TagNumber(11638560)
  $pb.PbList<$core.String> get tokendomains => $_getList(0);

  @$pb.TagNumber(39797417)
  $core.String get labelnamespace => $_getSZ(1);
  @$pb.TagNumber(39797417)
  set labelnamespace($core.String value) => $_setString(1, value);
  @$pb.TagNumber(39797417)
  $core.bool hasLabelnamespace() => $_has(1);
  @$pb.TagNumber(39797417)
  void clearLabelnamespace() => $_clearField(39797417);

  @$pb.TagNumber(42675585)
  $pb.PbList<Rule> get rules => $_getList(2);

  @$pb.TagNumber(42731774)
  $pb.PbMap<$core.String, CustomResponseBody> get customresponsebodies =>
      $_getMap(3);

  @$pb.TagNumber(48990889)
  ChallengeConfig get challengeconfig => $_getN(4);
  @$pb.TagNumber(48990889)
  set challengeconfig(ChallengeConfig value) => $_setField(48990889, value);
  @$pb.TagNumber(48990889)
  $core.bool hasChallengeconfig() => $_has(4);
  @$pb.TagNumber(48990889)
  void clearChallengeconfig() => $_clearField(48990889);
  @$pb.TagNumber(48990889)
  ChallengeConfig ensureChallengeconfig() => $_ensure(4);

  @$pb.TagNumber(60547064)
  CaptchaConfig get captchaconfig => $_getN(5);
  @$pb.TagNumber(60547064)
  set captchaconfig(CaptchaConfig value) => $_setField(60547064, value);
  @$pb.TagNumber(60547064)
  $core.bool hasCaptchaconfig() => $_has(5);
  @$pb.TagNumber(60547064)
  void clearCaptchaconfig() => $_clearField(60547064);
  @$pb.TagNumber(60547064)
  CaptchaConfig ensureCaptchaconfig() => $_ensure(5);

  @$pb.TagNumber(105229063)
  OnSourceDDoSProtectionConfig get onsourceddosprotectionconfig => $_getN(6);
  @$pb.TagNumber(105229063)
  set onsourceddosprotectionconfig(OnSourceDDoSProtectionConfig value) =>
      $_setField(105229063, value);
  @$pb.TagNumber(105229063)
  $core.bool hasOnsourceddosprotectionconfig() => $_has(6);
  @$pb.TagNumber(105229063)
  void clearOnsourceddosprotectionconfig() => $_clearField(105229063);
  @$pb.TagNumber(105229063)
  OnSourceDDoSProtectionConfig ensureOnsourceddosprotectionconfig() =>
      $_ensure(6);

  @$pb.TagNumber(107253930)
  $fixnum.Int64 get capacity => $_getI64(7);
  @$pb.TagNumber(107253930)
  set capacity($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(107253930)
  $core.bool hasCapacity() => $_has(7);
  @$pb.TagNumber(107253930)
  void clearCapacity() => $_clearField(107253930);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(8);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(8, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(8);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(9);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(9, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(9);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(279052502)
  $core.bool get retrofittedbyfirewallmanager => $_getBF(10);
  @$pb.TagNumber(279052502)
  set retrofittedbyfirewallmanager($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(279052502)
  $core.bool hasRetrofittedbyfirewallmanager() => $_has(10);
  @$pb.TagNumber(279052502)
  void clearRetrofittedbyfirewallmanager() => $_clearField(279052502);

  @$pb.TagNumber(322663861)
  DefaultAction get defaultaction => $_getN(11);
  @$pb.TagNumber(322663861)
  set defaultaction(DefaultAction value) => $_setField(322663861, value);
  @$pb.TagNumber(322663861)
  $core.bool hasDefaultaction() => $_has(11);
  @$pb.TagNumber(322663861)
  void clearDefaultaction() => $_clearField(322663861);
  @$pb.TagNumber(322663861)
  DefaultAction ensureDefaultaction() => $_ensure(11);

  @$pb.TagNumber(331512685)
  $pb.PbList<FirewallManagerRuleGroup>
      get preprocessfirewallmanagerrulegroups => $_getList(12);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(13);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(13, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(13);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(14);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(14, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(14);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);

  @$pb.TagNumber(415792743)
  $core.bool get managedbyfirewallmanager => $_getBF(15);
  @$pb.TagNumber(415792743)
  set managedbyfirewallmanager($core.bool value) => $_setBool(15, value);
  @$pb.TagNumber(415792743)
  $core.bool hasManagedbyfirewallmanager() => $_has(15);
  @$pb.TagNumber(415792743)
  void clearManagedbyfirewallmanager() => $_clearField(415792743);

  @$pb.TagNumber(419675691)
  AssociationConfig get associationconfig => $_getN(16);
  @$pb.TagNumber(419675691)
  set associationconfig(AssociationConfig value) =>
      $_setField(419675691, value);
  @$pb.TagNumber(419675691)
  $core.bool hasAssociationconfig() => $_has(16);
  @$pb.TagNumber(419675691)
  void clearAssociationconfig() => $_clearField(419675691);
  @$pb.TagNumber(419675691)
  AssociationConfig ensureAssociationconfig() => $_ensure(16);

  @$pb.TagNumber(464792245)
  DataProtectionConfig get dataprotectionconfig => $_getN(17);
  @$pb.TagNumber(464792245)
  set dataprotectionconfig(DataProtectionConfig value) =>
      $_setField(464792245, value);
  @$pb.TagNumber(464792245)
  $core.bool hasDataprotectionconfig() => $_has(17);
  @$pb.TagNumber(464792245)
  void clearDataprotectionconfig() => $_clearField(464792245);
  @$pb.TagNumber(464792245)
  DataProtectionConfig ensureDataprotectionconfig() => $_ensure(17);

  @$pb.TagNumber(494837446)
  $pb.PbList<FirewallManagerRuleGroup>
      get postprocessfirewallmanagerrulegroups => $_getList(18);

  @$pb.TagNumber(501131976)
  ApplicationConfig get applicationconfig => $_getN(19);
  @$pb.TagNumber(501131976)
  set applicationconfig(ApplicationConfig value) =>
      $_setField(501131976, value);
  @$pb.TagNumber(501131976)
  $core.bool hasApplicationconfig() => $_has(19);
  @$pb.TagNumber(501131976)
  void clearApplicationconfig() => $_clearField(501131976);
  @$pb.TagNumber(501131976)
  ApplicationConfig ensureApplicationconfig() => $_ensure(19);

  @$pb.TagNumber(507762728)
  VisibilityConfig get visibilityconfig => $_getN(20);
  @$pb.TagNumber(507762728)
  set visibilityconfig(VisibilityConfig value) => $_setField(507762728, value);
  @$pb.TagNumber(507762728)
  $core.bool hasVisibilityconfig() => $_has(20);
  @$pb.TagNumber(507762728)
  void clearVisibilityconfig() => $_clearField(507762728);
  @$pb.TagNumber(507762728)
  VisibilityConfig ensureVisibilityconfig() => $_ensure(20);
}

class WebACLSummary extends $pb.GeneratedMessage {
  factory WebACLSummary({
    $core.String? locktoken,
    $core.String? description,
    $core.String? name,
    $core.String? id,
    $core.String? arn,
  }) {
    final result = create();
    if (locktoken != null) result.locktoken = locktoken;
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    return result;
  }

  WebACLSummary._();

  factory WebACLSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WebACLSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WebACLSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..aOS(47056792, _omitFieldNames ? '' : 'locktoken')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(397135389, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WebACLSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WebACLSummary copyWith(void Function(WebACLSummary) updates) =>
      super.copyWith((message) => updates(message as WebACLSummary))
          as WebACLSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WebACLSummary create() => WebACLSummary._();
  @$core.override
  WebACLSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WebACLSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WebACLSummary>(create);
  static WebACLSummary? _defaultInstance;

  @$pb.TagNumber(47056792)
  $core.String get locktoken => $_getSZ(0);
  @$pb.TagNumber(47056792)
  set locktoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(47056792)
  $core.bool hasLocktoken() => $_has(0);
  @$pb.TagNumber(47056792)
  void clearLocktoken() => $_clearField(47056792);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(397135389)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(397135389)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(397135389)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(397135389)
  void clearArn() => $_clearField(397135389);
}

class XssMatchStatement extends $pb.GeneratedMessage {
  factory XssMatchStatement({
    $core.Iterable<TextTransformation>? texttransformations,
    FieldToMatch? fieldtomatch,
  }) {
    final result = create();
    if (texttransformations != null)
      result.texttransformations.addAll(texttransformations);
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    return result;
  }

  XssMatchStatement._();

  factory XssMatchStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XssMatchStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XssMatchStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'wafv2'),
      createEmptyInstance: create)
    ..pPM<TextTransformation>(
        261837309, _omitFieldNames ? '' : 'texttransformations',
        subBuilder: TextTransformation.create)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchStatement copyWith(void Function(XssMatchStatement) updates) =>
      super.copyWith((message) => updates(message as XssMatchStatement))
          as XssMatchStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XssMatchStatement create() => XssMatchStatement._();
  @$core.override
  XssMatchStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XssMatchStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XssMatchStatement>(create);
  static XssMatchStatement? _defaultInstance;

  @$pb.TagNumber(261837309)
  $pb.PbList<TextTransformation> get texttransformations => $_getList(0);

  @$pb.TagNumber(338348372)
  FieldToMatch get fieldtomatch => $_getN(1);
  @$pb.TagNumber(338348372)
  set fieldtomatch(FieldToMatch value) => $_setField(338348372, value);
  @$pb.TagNumber(338348372)
  $core.bool hasFieldtomatch() => $_has(1);
  @$pb.TagNumber(338348372)
  void clearFieldtomatch() => $_clearField(338348372);
  @$pb.TagNumber(338348372)
  FieldToMatch ensureFieldtomatch() => $_ensure(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
