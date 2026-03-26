// This is a generated file - do not edit.
//
// Generated from dynamodb.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'dynamodb.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'dynamodb.pbenum.dart';

class ArchivalSummary extends $pb.GeneratedMessage {
  factory ArchivalSummary({
    $core.String? archivalbackuparn,
    $core.String? archivalreason,
    $core.String? archivaldatetime,
  }) {
    final result = create();
    if (archivalbackuparn != null) result.archivalbackuparn = archivalbackuparn;
    if (archivalreason != null) result.archivalreason = archivalreason;
    if (archivaldatetime != null) result.archivaldatetime = archivaldatetime;
    return result;
  }

  ArchivalSummary._();

  factory ArchivalSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ArchivalSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ArchivalSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(95197293, _omitFieldNames ? '' : 'archivalbackuparn')
    ..aOS(132420188, _omitFieldNames ? '' : 'archivalreason')
    ..aOS(424460035, _omitFieldNames ? '' : 'archivaldatetime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ArchivalSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ArchivalSummary copyWith(void Function(ArchivalSummary) updates) =>
      super.copyWith((message) => updates(message as ArchivalSummary))
          as ArchivalSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ArchivalSummary create() => ArchivalSummary._();
  @$core.override
  ArchivalSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ArchivalSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ArchivalSummary>(create);
  static ArchivalSummary? _defaultInstance;

  @$pb.TagNumber(95197293)
  $core.String get archivalbackuparn => $_getSZ(0);
  @$pb.TagNumber(95197293)
  set archivalbackuparn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(95197293)
  $core.bool hasArchivalbackuparn() => $_has(0);
  @$pb.TagNumber(95197293)
  void clearArchivalbackuparn() => $_clearField(95197293);

  @$pb.TagNumber(132420188)
  $core.String get archivalreason => $_getSZ(1);
  @$pb.TagNumber(132420188)
  set archivalreason($core.String value) => $_setString(1, value);
  @$pb.TagNumber(132420188)
  $core.bool hasArchivalreason() => $_has(1);
  @$pb.TagNumber(132420188)
  void clearArchivalreason() => $_clearField(132420188);

  @$pb.TagNumber(424460035)
  $core.String get archivaldatetime => $_getSZ(2);
  @$pb.TagNumber(424460035)
  set archivaldatetime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(424460035)
  $core.bool hasArchivaldatetime() => $_has(2);
  @$pb.TagNumber(424460035)
  void clearArchivaldatetime() => $_clearField(424460035);
}

class AttributeDefinition extends $pb.GeneratedMessage {
  factory AttributeDefinition({
    ScalarAttributeType? attributetype,
    $core.String? attributename,
  }) {
    final result = create();
    if (attributetype != null) result.attributetype = attributetype;
    if (attributename != null) result.attributename = attributename;
    return result;
  }

  AttributeDefinition._();

  factory AttributeDefinition.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AttributeDefinition.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AttributeDefinition',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ScalarAttributeType>(54612120, _omitFieldNames ? '' : 'attributetype',
        enumValues: ScalarAttributeType.values)
    ..aOS(352717485, _omitFieldNames ? '' : 'attributename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributeDefinition clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributeDefinition copyWith(void Function(AttributeDefinition) updates) =>
      super.copyWith((message) => updates(message as AttributeDefinition))
          as AttributeDefinition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AttributeDefinition create() => AttributeDefinition._();
  @$core.override
  AttributeDefinition createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AttributeDefinition getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AttributeDefinition>(create);
  static AttributeDefinition? _defaultInstance;

  @$pb.TagNumber(54612120)
  ScalarAttributeType get attributetype => $_getN(0);
  @$pb.TagNumber(54612120)
  set attributetype(ScalarAttributeType value) => $_setField(54612120, value);
  @$pb.TagNumber(54612120)
  $core.bool hasAttributetype() => $_has(0);
  @$pb.TagNumber(54612120)
  void clearAttributetype() => $_clearField(54612120);

  @$pb.TagNumber(352717485)
  $core.String get attributename => $_getSZ(1);
  @$pb.TagNumber(352717485)
  set attributename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(352717485)
  $core.bool hasAttributename() => $_has(1);
  @$pb.TagNumber(352717485)
  void clearAttributename() => $_clearField(352717485);
}

class AttributeValue extends $pb.GeneratedMessage {
  factory AttributeValue({
    $core.Iterable<$core.String>? ns,
    $core.Iterable<$core.String>? ss,
    $core.List<$core.int>? b,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? m,
    $core.Iterable<AttributeValue>? l,
    $core.String? n,
    $core.Iterable<$core.List<$core.int>>? bs,
    $core.bool? bool_307144766,
    $core.String? s,
    $core.bool? null_426761765,
  }) {
    final result = create();
    if (ns != null) result.ns.addAll(ns);
    if (ss != null) result.ss.addAll(ss);
    if (b != null) result.b = b;
    if (m != null) result.m.addEntries(m);
    if (l != null) result.l.addAll(l);
    if (n != null) result.n = n;
    if (bs != null) result.bs.addAll(bs);
    if (bool_307144766 != null) result.bool_307144766 = bool_307144766;
    if (s != null) result.s = s;
    if (null_426761765 != null) result.null_426761765 = null_426761765;
    return result;
  }

  AttributeValue._();

  factory AttributeValue.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AttributeValue.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AttributeValue',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPS(98983847, _omitFieldNames ? '' : 'ns')
    ..pPS(100834020, _omitFieldNames ? '' : 'ss')
    ..a<$core.List<$core.int>>(
        118225798, _omitFieldNames ? '' : 'b', $pb.PbFieldType.OY)
    ..m<$core.String, AttributeValue>(135003417, _omitFieldNames ? '' : 'm',
        entryClassName: 'AttributeValue.MEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..pPM<AttributeValue>(151781036, _omitFieldNames ? '' : 'l',
        subBuilder: AttributeValue.create)
    ..aOS(185336274, _omitFieldNames ? '' : 'n')
    ..p<$core.List<$core.int>>(
        232616419, _omitFieldNames ? '' : 'bs', $pb.PbFieldType.PY)
    ..aOB(307144766, _omitFieldNames ? '' : 'bool')
    ..aOS(369890083, _omitFieldNames ? '' : 's')
    ..aOB(426761765, _omitFieldNames ? '' : 'null')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributeValue clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributeValue copyWith(void Function(AttributeValue) updates) =>
      super.copyWith((message) => updates(message as AttributeValue))
          as AttributeValue;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AttributeValue create() => AttributeValue._();
  @$core.override
  AttributeValue createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AttributeValue getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AttributeValue>(create);
  static AttributeValue? _defaultInstance;

  @$pb.TagNumber(98983847)
  $pb.PbList<$core.String> get ns => $_getList(0);

  @$pb.TagNumber(100834020)
  $pb.PbList<$core.String> get ss => $_getList(1);

  @$pb.TagNumber(118225798)
  $core.List<$core.int> get b => $_getN(2);
  @$pb.TagNumber(118225798)
  set b($core.List<$core.int> value) => $_setBytes(2, value);
  @$pb.TagNumber(118225798)
  $core.bool hasB() => $_has(2);
  @$pb.TagNumber(118225798)
  void clearB() => $_clearField(118225798);

  @$pb.TagNumber(135003417)
  $pb.PbMap<$core.String, AttributeValue> get m => $_getMap(3);

  @$pb.TagNumber(151781036)
  $pb.PbList<AttributeValue> get l => $_getList(4);

  @$pb.TagNumber(185336274)
  $core.String get n => $_getSZ(5);
  @$pb.TagNumber(185336274)
  set n($core.String value) => $_setString(5, value);
  @$pb.TagNumber(185336274)
  $core.bool hasN() => $_has(5);
  @$pb.TagNumber(185336274)
  void clearN() => $_clearField(185336274);

  @$pb.TagNumber(232616419)
  $pb.PbList<$core.List<$core.int>> get bs => $_getList(6);

  @$pb.TagNumber(307144766)
  $core.bool get bool_307144766 => $_getBF(7);
  @$pb.TagNumber(307144766)
  set bool_307144766($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(307144766)
  $core.bool hasBool_307144766() => $_has(7);
  @$pb.TagNumber(307144766)
  void clearBool_307144766() => $_clearField(307144766);

  @$pb.TagNumber(369890083)
  $core.String get s => $_getSZ(8);
  @$pb.TagNumber(369890083)
  set s($core.String value) => $_setString(8, value);
  @$pb.TagNumber(369890083)
  $core.bool hasS() => $_has(8);
  @$pb.TagNumber(369890083)
  void clearS() => $_clearField(369890083);

  @$pb.TagNumber(426761765)
  $core.bool get null_426761765 => $_getBF(9);
  @$pb.TagNumber(426761765)
  set null_426761765($core.bool value) => $_setBool(9, value);
  @$pb.TagNumber(426761765)
  $core.bool hasNull_426761765() => $_has(9);
  @$pb.TagNumber(426761765)
  void clearNull_426761765() => $_clearField(426761765);
}

class AttributeValueUpdate extends $pb.GeneratedMessage {
  factory AttributeValueUpdate({
    AttributeAction? action,
    AttributeValue? value,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (value != null) result.value = value;
    return result;
  }

  AttributeValueUpdate._();

  factory AttributeValueUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AttributeValueUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AttributeValueUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<AttributeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: AttributeAction.values)
    ..aOM<AttributeValue>(289929579, _omitFieldNames ? '' : 'value',
        subBuilder: AttributeValue.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributeValueUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AttributeValueUpdate copyWith(void Function(AttributeValueUpdate) updates) =>
      super.copyWith((message) => updates(message as AttributeValueUpdate))
          as AttributeValueUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AttributeValueUpdate create() => AttributeValueUpdate._();
  @$core.override
  AttributeValueUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AttributeValueUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AttributeValueUpdate>(create);
  static AttributeValueUpdate? _defaultInstance;

  @$pb.TagNumber(175614240)
  AttributeAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(AttributeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(289929579)
  AttributeValue get value => $_getN(1);
  @$pb.TagNumber(289929579)
  set value(AttributeValue value) => $_setField(289929579, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
  @$pb.TagNumber(289929579)
  AttributeValue ensureValue() => $_ensure(1);
}

class AutoScalingPolicyDescription extends $pb.GeneratedMessage {
  factory AutoScalingPolicyDescription({
    $core.String? policyname,
    AutoScalingTargetTrackingScalingPolicyConfigurationDescription?
        targettrackingscalingpolicyconfiguration,
  }) {
    final result = create();
    if (policyname != null) result.policyname = policyname;
    if (targettrackingscalingpolicyconfiguration != null)
      result.targettrackingscalingpolicyconfiguration =
          targettrackingscalingpolicyconfiguration;
    return result;
  }

  AutoScalingPolicyDescription._();

  factory AutoScalingPolicyDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AutoScalingPolicyDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoScalingPolicyDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(266468029, _omitFieldNames ? '' : 'policyname')
    ..aOM<AutoScalingTargetTrackingScalingPolicyConfigurationDescription>(
        474521147,
        _omitFieldNames ? '' : 'targettrackingscalingpolicyconfiguration',
        subBuilder:
            AutoScalingTargetTrackingScalingPolicyConfigurationDescription
                .create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingPolicyDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingPolicyDescription copyWith(
          void Function(AutoScalingPolicyDescription) updates) =>
      super.copyWith(
              (message) => updates(message as AutoScalingPolicyDescription))
          as AutoScalingPolicyDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AutoScalingPolicyDescription create() =>
      AutoScalingPolicyDescription._();
  @$core.override
  AutoScalingPolicyDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AutoScalingPolicyDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoScalingPolicyDescription>(create);
  static AutoScalingPolicyDescription? _defaultInstance;

  @$pb.TagNumber(266468029)
  $core.String get policyname => $_getSZ(0);
  @$pb.TagNumber(266468029)
  set policyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266468029)
  $core.bool hasPolicyname() => $_has(0);
  @$pb.TagNumber(266468029)
  void clearPolicyname() => $_clearField(266468029);

  @$pb.TagNumber(474521147)
  AutoScalingTargetTrackingScalingPolicyConfigurationDescription
      get targettrackingscalingpolicyconfiguration => $_getN(1);
  @$pb.TagNumber(474521147)
  set targettrackingscalingpolicyconfiguration(
          AutoScalingTargetTrackingScalingPolicyConfigurationDescription
              value) =>
      $_setField(474521147, value);
  @$pb.TagNumber(474521147)
  $core.bool hasTargettrackingscalingpolicyconfiguration() => $_has(1);
  @$pb.TagNumber(474521147)
  void clearTargettrackingscalingpolicyconfiguration() =>
      $_clearField(474521147);
  @$pb.TagNumber(474521147)
  AutoScalingTargetTrackingScalingPolicyConfigurationDescription
      ensureTargettrackingscalingpolicyconfiguration() => $_ensure(1);
}

class AutoScalingPolicyUpdate extends $pb.GeneratedMessage {
  factory AutoScalingPolicyUpdate({
    $core.String? policyname,
    AutoScalingTargetTrackingScalingPolicyConfigurationUpdate?
        targettrackingscalingpolicyconfiguration,
  }) {
    final result = create();
    if (policyname != null) result.policyname = policyname;
    if (targettrackingscalingpolicyconfiguration != null)
      result.targettrackingscalingpolicyconfiguration =
          targettrackingscalingpolicyconfiguration;
    return result;
  }

  AutoScalingPolicyUpdate._();

  factory AutoScalingPolicyUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AutoScalingPolicyUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoScalingPolicyUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(266468029, _omitFieldNames ? '' : 'policyname')
    ..aOM<AutoScalingTargetTrackingScalingPolicyConfigurationUpdate>(474521147,
        _omitFieldNames ? '' : 'targettrackingscalingpolicyconfiguration',
        subBuilder:
            AutoScalingTargetTrackingScalingPolicyConfigurationUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingPolicyUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingPolicyUpdate copyWith(
          void Function(AutoScalingPolicyUpdate) updates) =>
      super.copyWith((message) => updates(message as AutoScalingPolicyUpdate))
          as AutoScalingPolicyUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AutoScalingPolicyUpdate create() => AutoScalingPolicyUpdate._();
  @$core.override
  AutoScalingPolicyUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AutoScalingPolicyUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoScalingPolicyUpdate>(create);
  static AutoScalingPolicyUpdate? _defaultInstance;

  @$pb.TagNumber(266468029)
  $core.String get policyname => $_getSZ(0);
  @$pb.TagNumber(266468029)
  set policyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266468029)
  $core.bool hasPolicyname() => $_has(0);
  @$pb.TagNumber(266468029)
  void clearPolicyname() => $_clearField(266468029);

  @$pb.TagNumber(474521147)
  AutoScalingTargetTrackingScalingPolicyConfigurationUpdate
      get targettrackingscalingpolicyconfiguration => $_getN(1);
  @$pb.TagNumber(474521147)
  set targettrackingscalingpolicyconfiguration(
          AutoScalingTargetTrackingScalingPolicyConfigurationUpdate value) =>
      $_setField(474521147, value);
  @$pb.TagNumber(474521147)
  $core.bool hasTargettrackingscalingpolicyconfiguration() => $_has(1);
  @$pb.TagNumber(474521147)
  void clearTargettrackingscalingpolicyconfiguration() =>
      $_clearField(474521147);
  @$pb.TagNumber(474521147)
  AutoScalingTargetTrackingScalingPolicyConfigurationUpdate
      ensureTargettrackingscalingpolicyconfiguration() => $_ensure(1);
}

class AutoScalingSettingsDescription extends $pb.GeneratedMessage {
  factory AutoScalingSettingsDescription({
    $core.bool? autoscalingdisabled,
    $fixnum.Int64? minimumunits,
    $core.Iterable<AutoScalingPolicyDescription>? scalingpolicies,
    $core.String? autoscalingrolearn,
    $fixnum.Int64? maximumunits,
  }) {
    final result = create();
    if (autoscalingdisabled != null)
      result.autoscalingdisabled = autoscalingdisabled;
    if (minimumunits != null) result.minimumunits = minimumunits;
    if (scalingpolicies != null) result.scalingpolicies.addAll(scalingpolicies);
    if (autoscalingrolearn != null)
      result.autoscalingrolearn = autoscalingrolearn;
    if (maximumunits != null) result.maximumunits = maximumunits;
    return result;
  }

  AutoScalingSettingsDescription._();

  factory AutoScalingSettingsDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AutoScalingSettingsDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoScalingSettingsDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOB(183835736, _omitFieldNames ? '' : 'autoscalingdisabled')
    ..aInt64(190148659, _omitFieldNames ? '' : 'minimumunits')
    ..pPM<AutoScalingPolicyDescription>(
        289494257, _omitFieldNames ? '' : 'scalingpolicies',
        subBuilder: AutoScalingPolicyDescription.create)
    ..aOS(348085703, _omitFieldNames ? '' : 'autoscalingrolearn')
    ..aInt64(501867481, _omitFieldNames ? '' : 'maximumunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingSettingsDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingSettingsDescription copyWith(
          void Function(AutoScalingSettingsDescription) updates) =>
      super.copyWith(
              (message) => updates(message as AutoScalingSettingsDescription))
          as AutoScalingSettingsDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AutoScalingSettingsDescription create() =>
      AutoScalingSettingsDescription._();
  @$core.override
  AutoScalingSettingsDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AutoScalingSettingsDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoScalingSettingsDescription>(create);
  static AutoScalingSettingsDescription? _defaultInstance;

  @$pb.TagNumber(183835736)
  $core.bool get autoscalingdisabled => $_getBF(0);
  @$pb.TagNumber(183835736)
  set autoscalingdisabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(183835736)
  $core.bool hasAutoscalingdisabled() => $_has(0);
  @$pb.TagNumber(183835736)
  void clearAutoscalingdisabled() => $_clearField(183835736);

  @$pb.TagNumber(190148659)
  $fixnum.Int64 get minimumunits => $_getI64(1);
  @$pb.TagNumber(190148659)
  set minimumunits($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(190148659)
  $core.bool hasMinimumunits() => $_has(1);
  @$pb.TagNumber(190148659)
  void clearMinimumunits() => $_clearField(190148659);

  @$pb.TagNumber(289494257)
  $pb.PbList<AutoScalingPolicyDescription> get scalingpolicies => $_getList(2);

  @$pb.TagNumber(348085703)
  $core.String get autoscalingrolearn => $_getSZ(3);
  @$pb.TagNumber(348085703)
  set autoscalingrolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(348085703)
  $core.bool hasAutoscalingrolearn() => $_has(3);
  @$pb.TagNumber(348085703)
  void clearAutoscalingrolearn() => $_clearField(348085703);

  @$pb.TagNumber(501867481)
  $fixnum.Int64 get maximumunits => $_getI64(4);
  @$pb.TagNumber(501867481)
  set maximumunits($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(501867481)
  $core.bool hasMaximumunits() => $_has(4);
  @$pb.TagNumber(501867481)
  void clearMaximumunits() => $_clearField(501867481);
}

class AutoScalingSettingsUpdate extends $pb.GeneratedMessage {
  factory AutoScalingSettingsUpdate({
    $core.bool? autoscalingdisabled,
    AutoScalingPolicyUpdate? scalingpolicyupdate,
    $fixnum.Int64? minimumunits,
    $core.String? autoscalingrolearn,
    $fixnum.Int64? maximumunits,
  }) {
    final result = create();
    if (autoscalingdisabled != null)
      result.autoscalingdisabled = autoscalingdisabled;
    if (scalingpolicyupdate != null)
      result.scalingpolicyupdate = scalingpolicyupdate;
    if (minimumunits != null) result.minimumunits = minimumunits;
    if (autoscalingrolearn != null)
      result.autoscalingrolearn = autoscalingrolearn;
    if (maximumunits != null) result.maximumunits = maximumunits;
    return result;
  }

  AutoScalingSettingsUpdate._();

  factory AutoScalingSettingsUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AutoScalingSettingsUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AutoScalingSettingsUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOB(183835736, _omitFieldNames ? '' : 'autoscalingdisabled')
    ..aOM<AutoScalingPolicyUpdate>(
        189139846, _omitFieldNames ? '' : 'scalingpolicyupdate',
        subBuilder: AutoScalingPolicyUpdate.create)
    ..aInt64(190148659, _omitFieldNames ? '' : 'minimumunits')
    ..aOS(348085703, _omitFieldNames ? '' : 'autoscalingrolearn')
    ..aInt64(501867481, _omitFieldNames ? '' : 'maximumunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingSettingsUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingSettingsUpdate copyWith(
          void Function(AutoScalingSettingsUpdate) updates) =>
      super.copyWith((message) => updates(message as AutoScalingSettingsUpdate))
          as AutoScalingSettingsUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AutoScalingSettingsUpdate create() => AutoScalingSettingsUpdate._();
  @$core.override
  AutoScalingSettingsUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AutoScalingSettingsUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AutoScalingSettingsUpdate>(create);
  static AutoScalingSettingsUpdate? _defaultInstance;

  @$pb.TagNumber(183835736)
  $core.bool get autoscalingdisabled => $_getBF(0);
  @$pb.TagNumber(183835736)
  set autoscalingdisabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(183835736)
  $core.bool hasAutoscalingdisabled() => $_has(0);
  @$pb.TagNumber(183835736)
  void clearAutoscalingdisabled() => $_clearField(183835736);

  @$pb.TagNumber(189139846)
  AutoScalingPolicyUpdate get scalingpolicyupdate => $_getN(1);
  @$pb.TagNumber(189139846)
  set scalingpolicyupdate(AutoScalingPolicyUpdate value) =>
      $_setField(189139846, value);
  @$pb.TagNumber(189139846)
  $core.bool hasScalingpolicyupdate() => $_has(1);
  @$pb.TagNumber(189139846)
  void clearScalingpolicyupdate() => $_clearField(189139846);
  @$pb.TagNumber(189139846)
  AutoScalingPolicyUpdate ensureScalingpolicyupdate() => $_ensure(1);

  @$pb.TagNumber(190148659)
  $fixnum.Int64 get minimumunits => $_getI64(2);
  @$pb.TagNumber(190148659)
  set minimumunits($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(190148659)
  $core.bool hasMinimumunits() => $_has(2);
  @$pb.TagNumber(190148659)
  void clearMinimumunits() => $_clearField(190148659);

  @$pb.TagNumber(348085703)
  $core.String get autoscalingrolearn => $_getSZ(3);
  @$pb.TagNumber(348085703)
  set autoscalingrolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(348085703)
  $core.bool hasAutoscalingrolearn() => $_has(3);
  @$pb.TagNumber(348085703)
  void clearAutoscalingrolearn() => $_clearField(348085703);

  @$pb.TagNumber(501867481)
  $fixnum.Int64 get maximumunits => $_getI64(4);
  @$pb.TagNumber(501867481)
  set maximumunits($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(501867481)
  $core.bool hasMaximumunits() => $_has(4);
  @$pb.TagNumber(501867481)
  void clearMaximumunits() => $_clearField(501867481);
}

class AutoScalingTargetTrackingScalingPolicyConfigurationDescription
    extends $pb.GeneratedMessage {
  factory AutoScalingTargetTrackingScalingPolicyConfigurationDescription({
    $core.int? scaleincooldown,
    $core.double? targetvalue,
    $core.int? scaleoutcooldown,
    $core.bool? disablescalein,
  }) {
    final result = create();
    if (scaleincooldown != null) result.scaleincooldown = scaleincooldown;
    if (targetvalue != null) result.targetvalue = targetvalue;
    if (scaleoutcooldown != null) result.scaleoutcooldown = scaleoutcooldown;
    if (disablescalein != null) result.disablescalein = disablescalein;
    return result;
  }

  AutoScalingTargetTrackingScalingPolicyConfigurationDescription._();

  factory AutoScalingTargetTrackingScalingPolicyConfigurationDescription.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AutoScalingTargetTrackingScalingPolicyConfigurationDescription.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames
          ? ''
          : 'AutoScalingTargetTrackingScalingPolicyConfigurationDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aI(45925876, _omitFieldNames ? '' : 'scaleincooldown')
    ..aD(118247738, _omitFieldNames ? '' : 'targetvalue')
    ..aI(316995887, _omitFieldNames ? '' : 'scaleoutcooldown')
    ..aOB(464986377, _omitFieldNames ? '' : 'disablescalein')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingTargetTrackingScalingPolicyConfigurationDescription clone() =>
      deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingTargetTrackingScalingPolicyConfigurationDescription copyWith(
          void Function(
                  AutoScalingTargetTrackingScalingPolicyConfigurationDescription)
              updates) =>
      super.copyWith((message) => updates(message
              as AutoScalingTargetTrackingScalingPolicyConfigurationDescription))
          as AutoScalingTargetTrackingScalingPolicyConfigurationDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AutoScalingTargetTrackingScalingPolicyConfigurationDescription
      create() =>
          AutoScalingTargetTrackingScalingPolicyConfigurationDescription._();
  @$core.override
  AutoScalingTargetTrackingScalingPolicyConfigurationDescription
      createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AutoScalingTargetTrackingScalingPolicyConfigurationDescription
      getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
              AutoScalingTargetTrackingScalingPolicyConfigurationDescription>(
          create);
  static AutoScalingTargetTrackingScalingPolicyConfigurationDescription?
      _defaultInstance;

  @$pb.TagNumber(45925876)
  $core.int get scaleincooldown => $_getIZ(0);
  @$pb.TagNumber(45925876)
  set scaleincooldown($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(45925876)
  $core.bool hasScaleincooldown() => $_has(0);
  @$pb.TagNumber(45925876)
  void clearScaleincooldown() => $_clearField(45925876);

  @$pb.TagNumber(118247738)
  $core.double get targetvalue => $_getN(1);
  @$pb.TagNumber(118247738)
  set targetvalue($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(118247738)
  $core.bool hasTargetvalue() => $_has(1);
  @$pb.TagNumber(118247738)
  void clearTargetvalue() => $_clearField(118247738);

  @$pb.TagNumber(316995887)
  $core.int get scaleoutcooldown => $_getIZ(2);
  @$pb.TagNumber(316995887)
  set scaleoutcooldown($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(316995887)
  $core.bool hasScaleoutcooldown() => $_has(2);
  @$pb.TagNumber(316995887)
  void clearScaleoutcooldown() => $_clearField(316995887);

  @$pb.TagNumber(464986377)
  $core.bool get disablescalein => $_getBF(3);
  @$pb.TagNumber(464986377)
  set disablescalein($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(464986377)
  $core.bool hasDisablescalein() => $_has(3);
  @$pb.TagNumber(464986377)
  void clearDisablescalein() => $_clearField(464986377);
}

class AutoScalingTargetTrackingScalingPolicyConfigurationUpdate
    extends $pb.GeneratedMessage {
  factory AutoScalingTargetTrackingScalingPolicyConfigurationUpdate({
    $core.int? scaleincooldown,
    $core.double? targetvalue,
    $core.int? scaleoutcooldown,
    $core.bool? disablescalein,
  }) {
    final result = create();
    if (scaleincooldown != null) result.scaleincooldown = scaleincooldown;
    if (targetvalue != null) result.targetvalue = targetvalue;
    if (scaleoutcooldown != null) result.scaleoutcooldown = scaleoutcooldown;
    if (disablescalein != null) result.disablescalein = disablescalein;
    return result;
  }

  AutoScalingTargetTrackingScalingPolicyConfigurationUpdate._();

  factory AutoScalingTargetTrackingScalingPolicyConfigurationUpdate.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AutoScalingTargetTrackingScalingPolicyConfigurationUpdate.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames
          ? ''
          : 'AutoScalingTargetTrackingScalingPolicyConfigurationUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aI(45925876, _omitFieldNames ? '' : 'scaleincooldown')
    ..aD(118247738, _omitFieldNames ? '' : 'targetvalue')
    ..aI(316995887, _omitFieldNames ? '' : 'scaleoutcooldown')
    ..aOB(464986377, _omitFieldNames ? '' : 'disablescalein')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingTargetTrackingScalingPolicyConfigurationUpdate clone() =>
      deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AutoScalingTargetTrackingScalingPolicyConfigurationUpdate copyWith(
          void Function(
                  AutoScalingTargetTrackingScalingPolicyConfigurationUpdate)
              updates) =>
      super.copyWith((message) => updates(message
              as AutoScalingTargetTrackingScalingPolicyConfigurationUpdate))
          as AutoScalingTargetTrackingScalingPolicyConfigurationUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AutoScalingTargetTrackingScalingPolicyConfigurationUpdate create() =>
      AutoScalingTargetTrackingScalingPolicyConfigurationUpdate._();
  @$core.override
  AutoScalingTargetTrackingScalingPolicyConfigurationUpdate
      createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AutoScalingTargetTrackingScalingPolicyConfigurationUpdate
      getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          AutoScalingTargetTrackingScalingPolicyConfigurationUpdate>(create);
  static AutoScalingTargetTrackingScalingPolicyConfigurationUpdate?
      _defaultInstance;

  @$pb.TagNumber(45925876)
  $core.int get scaleincooldown => $_getIZ(0);
  @$pb.TagNumber(45925876)
  set scaleincooldown($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(45925876)
  $core.bool hasScaleincooldown() => $_has(0);
  @$pb.TagNumber(45925876)
  void clearScaleincooldown() => $_clearField(45925876);

  @$pb.TagNumber(118247738)
  $core.double get targetvalue => $_getN(1);
  @$pb.TagNumber(118247738)
  set targetvalue($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(118247738)
  $core.bool hasTargetvalue() => $_has(1);
  @$pb.TagNumber(118247738)
  void clearTargetvalue() => $_clearField(118247738);

  @$pb.TagNumber(316995887)
  $core.int get scaleoutcooldown => $_getIZ(2);
  @$pb.TagNumber(316995887)
  set scaleoutcooldown($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(316995887)
  $core.bool hasScaleoutcooldown() => $_has(2);
  @$pb.TagNumber(316995887)
  void clearScaleoutcooldown() => $_clearField(316995887);

  @$pb.TagNumber(464986377)
  $core.bool get disablescalein => $_getBF(3);
  @$pb.TagNumber(464986377)
  set disablescalein($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(464986377)
  $core.bool hasDisablescalein() => $_has(3);
  @$pb.TagNumber(464986377)
  void clearDisablescalein() => $_clearField(464986377);
}

class BackupDescription extends $pb.GeneratedMessage {
  factory BackupDescription({
    SourceTableFeatureDetails? sourcetablefeaturedetails,
    BackupDetails? backupdetails,
    SourceTableDetails? sourcetabledetails,
  }) {
    final result = create();
    if (sourcetablefeaturedetails != null)
      result.sourcetablefeaturedetails = sourcetablefeaturedetails;
    if (backupdetails != null) result.backupdetails = backupdetails;
    if (sourcetabledetails != null)
      result.sourcetabledetails = sourcetabledetails;
    return result;
  }

  BackupDescription._();

  factory BackupDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BackupDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BackupDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<SourceTableFeatureDetails>(
        21598361, _omitFieldNames ? '' : 'sourcetablefeaturedetails',
        subBuilder: SourceTableFeatureDetails.create)
    ..aOM<BackupDetails>(383357432, _omitFieldNames ? '' : 'backupdetails',
        subBuilder: BackupDetails.create)
    ..aOM<SourceTableDetails>(
        507937979, _omitFieldNames ? '' : 'sourcetabledetails',
        subBuilder: SourceTableDetails.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupDescription copyWith(void Function(BackupDescription) updates) =>
      super.copyWith((message) => updates(message as BackupDescription))
          as BackupDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BackupDescription create() => BackupDescription._();
  @$core.override
  BackupDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BackupDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BackupDescription>(create);
  static BackupDescription? _defaultInstance;

  @$pb.TagNumber(21598361)
  SourceTableFeatureDetails get sourcetablefeaturedetails => $_getN(0);
  @$pb.TagNumber(21598361)
  set sourcetablefeaturedetails(SourceTableFeatureDetails value) =>
      $_setField(21598361, value);
  @$pb.TagNumber(21598361)
  $core.bool hasSourcetablefeaturedetails() => $_has(0);
  @$pb.TagNumber(21598361)
  void clearSourcetablefeaturedetails() => $_clearField(21598361);
  @$pb.TagNumber(21598361)
  SourceTableFeatureDetails ensureSourcetablefeaturedetails() => $_ensure(0);

  @$pb.TagNumber(383357432)
  BackupDetails get backupdetails => $_getN(1);
  @$pb.TagNumber(383357432)
  set backupdetails(BackupDetails value) => $_setField(383357432, value);
  @$pb.TagNumber(383357432)
  $core.bool hasBackupdetails() => $_has(1);
  @$pb.TagNumber(383357432)
  void clearBackupdetails() => $_clearField(383357432);
  @$pb.TagNumber(383357432)
  BackupDetails ensureBackupdetails() => $_ensure(1);

  @$pb.TagNumber(507937979)
  SourceTableDetails get sourcetabledetails => $_getN(2);
  @$pb.TagNumber(507937979)
  set sourcetabledetails(SourceTableDetails value) =>
      $_setField(507937979, value);
  @$pb.TagNumber(507937979)
  $core.bool hasSourcetabledetails() => $_has(2);
  @$pb.TagNumber(507937979)
  void clearSourcetabledetails() => $_clearField(507937979);
  @$pb.TagNumber(507937979)
  SourceTableDetails ensureSourcetabledetails() => $_ensure(2);
}

class BackupDetails extends $pb.GeneratedMessage {
  factory BackupDetails({
    BackupType? backuptype,
    $fixnum.Int64? backupsizebytes,
    $core.String? backupcreationdatetime,
    $core.String? backuparn,
    BackupStatus? backupstatus,
    $core.String? backupname,
    $core.String? backupexpirydatetime,
  }) {
    final result = create();
    if (backuptype != null) result.backuptype = backuptype;
    if (backupsizebytes != null) result.backupsizebytes = backupsizebytes;
    if (backupcreationdatetime != null)
      result.backupcreationdatetime = backupcreationdatetime;
    if (backuparn != null) result.backuparn = backuparn;
    if (backupstatus != null) result.backupstatus = backupstatus;
    if (backupname != null) result.backupname = backupname;
    if (backupexpirydatetime != null)
      result.backupexpirydatetime = backupexpirydatetime;
    return result;
  }

  BackupDetails._();

  factory BackupDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BackupDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BackupDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<BackupType>(134973992, _omitFieldNames ? '' : 'backuptype',
        enumValues: BackupType.values)
    ..aInt64(147336318, _omitFieldNames ? '' : 'backupsizebytes')
    ..aOS(368022476, _omitFieldNames ? '' : 'backupcreationdatetime')
    ..aOS(370874339, _omitFieldNames ? '' : 'backuparn')
    ..aE<BackupStatus>(382505546, _omitFieldNames ? '' : 'backupstatus',
        enumValues: BackupStatus.values)
    ..aOS(467693789, _omitFieldNames ? '' : 'backupname')
    ..aOS(471291762, _omitFieldNames ? '' : 'backupexpirydatetime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupDetails copyWith(void Function(BackupDetails) updates) =>
      super.copyWith((message) => updates(message as BackupDetails))
          as BackupDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BackupDetails create() => BackupDetails._();
  @$core.override
  BackupDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BackupDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BackupDetails>(create);
  static BackupDetails? _defaultInstance;

  @$pb.TagNumber(134973992)
  BackupType get backuptype => $_getN(0);
  @$pb.TagNumber(134973992)
  set backuptype(BackupType value) => $_setField(134973992, value);
  @$pb.TagNumber(134973992)
  $core.bool hasBackuptype() => $_has(0);
  @$pb.TagNumber(134973992)
  void clearBackuptype() => $_clearField(134973992);

  @$pb.TagNumber(147336318)
  $fixnum.Int64 get backupsizebytes => $_getI64(1);
  @$pb.TagNumber(147336318)
  set backupsizebytes($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(147336318)
  $core.bool hasBackupsizebytes() => $_has(1);
  @$pb.TagNumber(147336318)
  void clearBackupsizebytes() => $_clearField(147336318);

  @$pb.TagNumber(368022476)
  $core.String get backupcreationdatetime => $_getSZ(2);
  @$pb.TagNumber(368022476)
  set backupcreationdatetime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(368022476)
  $core.bool hasBackupcreationdatetime() => $_has(2);
  @$pb.TagNumber(368022476)
  void clearBackupcreationdatetime() => $_clearField(368022476);

  @$pb.TagNumber(370874339)
  $core.String get backuparn => $_getSZ(3);
  @$pb.TagNumber(370874339)
  set backuparn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(370874339)
  $core.bool hasBackuparn() => $_has(3);
  @$pb.TagNumber(370874339)
  void clearBackuparn() => $_clearField(370874339);

  @$pb.TagNumber(382505546)
  BackupStatus get backupstatus => $_getN(4);
  @$pb.TagNumber(382505546)
  set backupstatus(BackupStatus value) => $_setField(382505546, value);
  @$pb.TagNumber(382505546)
  $core.bool hasBackupstatus() => $_has(4);
  @$pb.TagNumber(382505546)
  void clearBackupstatus() => $_clearField(382505546);

  @$pb.TagNumber(467693789)
  $core.String get backupname => $_getSZ(5);
  @$pb.TagNumber(467693789)
  set backupname($core.String value) => $_setString(5, value);
  @$pb.TagNumber(467693789)
  $core.bool hasBackupname() => $_has(5);
  @$pb.TagNumber(467693789)
  void clearBackupname() => $_clearField(467693789);

  @$pb.TagNumber(471291762)
  $core.String get backupexpirydatetime => $_getSZ(6);
  @$pb.TagNumber(471291762)
  set backupexpirydatetime($core.String value) => $_setString(6, value);
  @$pb.TagNumber(471291762)
  $core.bool hasBackupexpirydatetime() => $_has(6);
  @$pb.TagNumber(471291762)
  void clearBackupexpirydatetime() => $_clearField(471291762);
}

class BackupInUseException extends $pb.GeneratedMessage {
  factory BackupInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  BackupInUseException._();

  factory BackupInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BackupInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BackupInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupInUseException copyWith(void Function(BackupInUseException) updates) =>
      super.copyWith((message) => updates(message as BackupInUseException))
          as BackupInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BackupInUseException create() => BackupInUseException._();
  @$core.override
  BackupInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BackupInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BackupInUseException>(create);
  static BackupInUseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class BackupNotFoundException extends $pb.GeneratedMessage {
  factory BackupNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  BackupNotFoundException._();

  factory BackupNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BackupNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BackupNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupNotFoundException copyWith(
          void Function(BackupNotFoundException) updates) =>
      super.copyWith((message) => updates(message as BackupNotFoundException))
          as BackupNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BackupNotFoundException create() => BackupNotFoundException._();
  @$core.override
  BackupNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BackupNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BackupNotFoundException>(create);
  static BackupNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class BackupSummary extends $pb.GeneratedMessage {
  factory BackupSummary({
    BackupType? backuptype,
    $fixnum.Int64? backupsizebytes,
    $core.String? tablename,
    $core.String? backupcreationdatetime,
    $core.String? backuparn,
    BackupStatus? backupstatus,
    $core.String? tablearn,
    $core.String? tableid,
    $core.String? backupname,
    $core.String? backupexpirydatetime,
  }) {
    final result = create();
    if (backuptype != null) result.backuptype = backuptype;
    if (backupsizebytes != null) result.backupsizebytes = backupsizebytes;
    if (tablename != null) result.tablename = tablename;
    if (backupcreationdatetime != null)
      result.backupcreationdatetime = backupcreationdatetime;
    if (backuparn != null) result.backuparn = backuparn;
    if (backupstatus != null) result.backupstatus = backupstatus;
    if (tablearn != null) result.tablearn = tablearn;
    if (tableid != null) result.tableid = tableid;
    if (backupname != null) result.backupname = backupname;
    if (backupexpirydatetime != null)
      result.backupexpirydatetime = backupexpirydatetime;
    return result;
  }

  BackupSummary._();

  factory BackupSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BackupSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BackupSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<BackupType>(134973992, _omitFieldNames ? '' : 'backuptype',
        enumValues: BackupType.values)
    ..aInt64(147336318, _omitFieldNames ? '' : 'backupsizebytes')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(368022476, _omitFieldNames ? '' : 'backupcreationdatetime')
    ..aOS(370874339, _omitFieldNames ? '' : 'backuparn')
    ..aE<BackupStatus>(382505546, _omitFieldNames ? '' : 'backupstatus',
        enumValues: BackupStatus.values)
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..aOS(449893011, _omitFieldNames ? '' : 'tableid')
    ..aOS(467693789, _omitFieldNames ? '' : 'backupname')
    ..aOS(471291762, _omitFieldNames ? '' : 'backupexpirydatetime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BackupSummary copyWith(void Function(BackupSummary) updates) =>
      super.copyWith((message) => updates(message as BackupSummary))
          as BackupSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BackupSummary create() => BackupSummary._();
  @$core.override
  BackupSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BackupSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BackupSummary>(create);
  static BackupSummary? _defaultInstance;

  @$pb.TagNumber(134973992)
  BackupType get backuptype => $_getN(0);
  @$pb.TagNumber(134973992)
  set backuptype(BackupType value) => $_setField(134973992, value);
  @$pb.TagNumber(134973992)
  $core.bool hasBackuptype() => $_has(0);
  @$pb.TagNumber(134973992)
  void clearBackuptype() => $_clearField(134973992);

  @$pb.TagNumber(147336318)
  $fixnum.Int64 get backupsizebytes => $_getI64(1);
  @$pb.TagNumber(147336318)
  set backupsizebytes($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(147336318)
  $core.bool hasBackupsizebytes() => $_has(1);
  @$pb.TagNumber(147336318)
  void clearBackupsizebytes() => $_clearField(147336318);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(2);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(2);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(368022476)
  $core.String get backupcreationdatetime => $_getSZ(3);
  @$pb.TagNumber(368022476)
  set backupcreationdatetime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(368022476)
  $core.bool hasBackupcreationdatetime() => $_has(3);
  @$pb.TagNumber(368022476)
  void clearBackupcreationdatetime() => $_clearField(368022476);

  @$pb.TagNumber(370874339)
  $core.String get backuparn => $_getSZ(4);
  @$pb.TagNumber(370874339)
  set backuparn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(370874339)
  $core.bool hasBackuparn() => $_has(4);
  @$pb.TagNumber(370874339)
  void clearBackuparn() => $_clearField(370874339);

  @$pb.TagNumber(382505546)
  BackupStatus get backupstatus => $_getN(5);
  @$pb.TagNumber(382505546)
  set backupstatus(BackupStatus value) => $_setField(382505546, value);
  @$pb.TagNumber(382505546)
  $core.bool hasBackupstatus() => $_has(5);
  @$pb.TagNumber(382505546)
  void clearBackupstatus() => $_clearField(382505546);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(6);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(6);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);

  @$pb.TagNumber(449893011)
  $core.String get tableid => $_getSZ(7);
  @$pb.TagNumber(449893011)
  set tableid($core.String value) => $_setString(7, value);
  @$pb.TagNumber(449893011)
  $core.bool hasTableid() => $_has(7);
  @$pb.TagNumber(449893011)
  void clearTableid() => $_clearField(449893011);

  @$pb.TagNumber(467693789)
  $core.String get backupname => $_getSZ(8);
  @$pb.TagNumber(467693789)
  set backupname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(467693789)
  $core.bool hasBackupname() => $_has(8);
  @$pb.TagNumber(467693789)
  void clearBackupname() => $_clearField(467693789);

  @$pb.TagNumber(471291762)
  $core.String get backupexpirydatetime => $_getSZ(9);
  @$pb.TagNumber(471291762)
  set backupexpirydatetime($core.String value) => $_setString(9, value);
  @$pb.TagNumber(471291762)
  $core.bool hasBackupexpirydatetime() => $_has(9);
  @$pb.TagNumber(471291762)
  void clearBackupexpirydatetime() => $_clearField(471291762);
}

class BatchExecuteStatementInput extends $pb.GeneratedMessage {
  factory BatchExecuteStatementInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<BatchStatementRequest>? statements,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (statements != null) result.statements.addAll(statements);
    return result;
  }

  BatchExecuteStatementInput._();

  factory BatchExecuteStatementInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchExecuteStatementInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchExecuteStatementInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..pPM<BatchStatementRequest>(488352288, _omitFieldNames ? '' : 'statements',
        subBuilder: BatchStatementRequest.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchExecuteStatementInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchExecuteStatementInput copyWith(
          void Function(BatchExecuteStatementInput) updates) =>
      super.copyWith(
              (message) => updates(message as BatchExecuteStatementInput))
          as BatchExecuteStatementInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchExecuteStatementInput create() => BatchExecuteStatementInput._();
  @$core.override
  BatchExecuteStatementInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchExecuteStatementInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchExecuteStatementInput>(create);
  static BatchExecuteStatementInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(488352288)
  $pb.PbList<BatchStatementRequest> get statements => $_getList(1);
}

class BatchExecuteStatementOutput extends $pb.GeneratedMessage {
  factory BatchExecuteStatementOutput({
    $core.Iterable<BatchStatementResponse>? responses,
    $core.Iterable<ConsumedCapacity>? consumedcapacity,
  }) {
    final result = create();
    if (responses != null) result.responses.addAll(responses);
    if (consumedcapacity != null)
      result.consumedcapacity.addAll(consumedcapacity);
    return result;
  }

  BatchExecuteStatementOutput._();

  factory BatchExecuteStatementOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchExecuteStatementOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchExecuteStatementOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<BatchStatementResponse>(114852856, _omitFieldNames ? '' : 'responses',
        subBuilder: BatchStatementResponse.create)
    ..pPM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchExecuteStatementOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchExecuteStatementOutput copyWith(
          void Function(BatchExecuteStatementOutput) updates) =>
      super.copyWith(
              (message) => updates(message as BatchExecuteStatementOutput))
          as BatchExecuteStatementOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchExecuteStatementOutput create() =>
      BatchExecuteStatementOutput._();
  @$core.override
  BatchExecuteStatementOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchExecuteStatementOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchExecuteStatementOutput>(create);
  static BatchExecuteStatementOutput? _defaultInstance;

  @$pb.TagNumber(114852856)
  $pb.PbList<BatchStatementResponse> get responses => $_getList(0);

  @$pb.TagNumber(449336620)
  $pb.PbList<ConsumedCapacity> get consumedcapacity => $_getList(1);
}

class BatchGetItemInput extends $pb.GeneratedMessage {
  factory BatchGetItemInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<$core.MapEntry<$core.String, KeysAndAttributes>>?
        requestitems,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (requestitems != null) result.requestitems.addEntries(requestitems);
    return result;
  }

  BatchGetItemInput._();

  factory BatchGetItemInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetItemInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetItemInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..m<$core.String, KeysAndAttributes>(
        247720687, _omitFieldNames ? '' : 'requestitems',
        entryClassName: 'BatchGetItemInput.RequestitemsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: KeysAndAttributes.create,
        valueDefaultOrMaker: KeysAndAttributes.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetItemInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetItemInput copyWith(void Function(BatchGetItemInput) updates) =>
      super.copyWith((message) => updates(message as BatchGetItemInput))
          as BatchGetItemInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetItemInput create() => BatchGetItemInput._();
  @$core.override
  BatchGetItemInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetItemInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetItemInput>(create);
  static BatchGetItemInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(247720687)
  $pb.PbMap<$core.String, KeysAndAttributes> get requestitems => $_getMap(1);
}

class BatchGetItemOutput extends $pb.GeneratedMessage {
  factory BatchGetItemOutput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? responses,
    $core.Iterable<$core.MapEntry<$core.String, KeysAndAttributes>>?
        unprocessedkeys,
    $core.Iterable<ConsumedCapacity>? consumedcapacity,
  }) {
    final result = create();
    if (responses != null) result.responses.addEntries(responses);
    if (unprocessedkeys != null)
      result.unprocessedkeys.addEntries(unprocessedkeys);
    if (consumedcapacity != null)
      result.consumedcapacity.addAll(consumedcapacity);
    return result;
  }

  BatchGetItemOutput._();

  factory BatchGetItemOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchGetItemOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchGetItemOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        114852856, _omitFieldNames ? '' : 'responses',
        entryClassName: 'BatchGetItemOutput.ResponsesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, KeysAndAttributes>(
        123999275, _omitFieldNames ? '' : 'unprocessedkeys',
        entryClassName: 'BatchGetItemOutput.UnprocessedkeysEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: KeysAndAttributes.create,
        valueDefaultOrMaker: KeysAndAttributes.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..pPM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetItemOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchGetItemOutput copyWith(void Function(BatchGetItemOutput) updates) =>
      super.copyWith((message) => updates(message as BatchGetItemOutput))
          as BatchGetItemOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchGetItemOutput create() => BatchGetItemOutput._();
  @$core.override
  BatchGetItemOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchGetItemOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchGetItemOutput>(create);
  static BatchGetItemOutput? _defaultInstance;

  @$pb.TagNumber(114852856)
  $pb.PbMap<$core.String, $core.String> get responses => $_getMap(0);

  @$pb.TagNumber(123999275)
  $pb.PbMap<$core.String, KeysAndAttributes> get unprocessedkeys => $_getMap(1);

  @$pb.TagNumber(449336620)
  $pb.PbList<ConsumedCapacity> get consumedcapacity => $_getList(2);
}

class BatchStatementError extends $pb.GeneratedMessage {
  factory BatchStatementError({
    $core.String? message,
    BatchStatementErrorCodeEnum? code,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (code != null) result.code = code;
    if (item != null) result.item.addEntries(item);
    return result;
  }

  BatchStatementError._();

  factory BatchStatementError.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchStatementError.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchStatementError',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aE<BatchStatementErrorCodeEnum>(425572629, _omitFieldNames ? '' : 'code',
        enumValues: BatchStatementErrorCodeEnum.values)
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'BatchStatementError.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchStatementError clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchStatementError copyWith(void Function(BatchStatementError) updates) =>
      super.copyWith((message) => updates(message as BatchStatementError))
          as BatchStatementError;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchStatementError create() => BatchStatementError._();
  @$core.override
  BatchStatementError createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchStatementError getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchStatementError>(create);
  static BatchStatementError? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(425572629)
  BatchStatementErrorCodeEnum get code => $_getN(1);
  @$pb.TagNumber(425572629)
  set code(BatchStatementErrorCodeEnum value) => $_setField(425572629, value);
  @$pb.TagNumber(425572629)
  $core.bool hasCode() => $_has(1);
  @$pb.TagNumber(425572629)
  void clearCode() => $_clearField(425572629);

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(2);
}

class BatchStatementRequest extends $pb.GeneratedMessage {
  factory BatchStatementRequest({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    $core.String? statement,
    $core.Iterable<AttributeValue>? parameters,
    $core.bool? consistentread,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (statement != null) result.statement = statement;
    if (parameters != null) result.parameters.addAll(parameters);
    if (consistentread != null) result.consistentread = consistentread;
    return result;
  }

  BatchStatementRequest._();

  factory BatchStatementRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchStatementRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchStatementRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..aOS(248790199, _omitFieldNames ? '' : 'statement')
    ..pPM<AttributeValue>(494900218, _omitFieldNames ? '' : 'parameters',
        subBuilder: AttributeValue.create)
    ..aOB(531556994, _omitFieldNames ? '' : 'consistentread')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchStatementRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchStatementRequest copyWith(
          void Function(BatchStatementRequest) updates) =>
      super.copyWith((message) => updates(message as BatchStatementRequest))
          as BatchStatementRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchStatementRequest create() => BatchStatementRequest._();
  @$core.override
  BatchStatementRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchStatementRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchStatementRequest>(create);
  static BatchStatementRequest? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(248790199)
  $core.String get statement => $_getSZ(1);
  @$pb.TagNumber(248790199)
  set statement($core.String value) => $_setString(1, value);
  @$pb.TagNumber(248790199)
  $core.bool hasStatement() => $_has(1);
  @$pb.TagNumber(248790199)
  void clearStatement() => $_clearField(248790199);

  @$pb.TagNumber(494900218)
  $pb.PbList<AttributeValue> get parameters => $_getList(2);

  @$pb.TagNumber(531556994)
  $core.bool get consistentread => $_getBF(3);
  @$pb.TagNumber(531556994)
  set consistentread($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(531556994)
  $core.bool hasConsistentread() => $_has(3);
  @$pb.TagNumber(531556994)
  void clearConsistentread() => $_clearField(531556994);
}

class BatchStatementResponse extends $pb.GeneratedMessage {
  factory BatchStatementResponse({
    $core.String? tablename,
    BatchStatementError? error,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    if (error != null) result.error = error;
    if (item != null) result.item.addEntries(item);
    return result;
  }

  BatchStatementResponse._();

  factory BatchStatementResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchStatementResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchStatementResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<BatchStatementError>(328047858, _omitFieldNames ? '' : 'error',
        subBuilder: BatchStatementError.create)
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'BatchStatementResponse.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchStatementResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchStatementResponse copyWith(
          void Function(BatchStatementResponse) updates) =>
      super.copyWith((message) => updates(message as BatchStatementResponse))
          as BatchStatementResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchStatementResponse create() => BatchStatementResponse._();
  @$core.override
  BatchStatementResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchStatementResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchStatementResponse>(create);
  static BatchStatementResponse? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(328047858)
  BatchStatementError get error => $_getN(1);
  @$pb.TagNumber(328047858)
  set error(BatchStatementError value) => $_setField(328047858, value);
  @$pb.TagNumber(328047858)
  $core.bool hasError() => $_has(1);
  @$pb.TagNumber(328047858)
  void clearError() => $_clearField(328047858);
  @$pb.TagNumber(328047858)
  BatchStatementError ensureError() => $_ensure(1);

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(2);
}

class BatchWriteItemInput extends $pb.GeneratedMessage {
  factory BatchWriteItemInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? requestitems,
    ReturnItemCollectionMetrics? returnitemcollectionmetrics,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (requestitems != null) result.requestitems.addEntries(requestitems);
    if (returnitemcollectionmetrics != null)
      result.returnitemcollectionmetrics = returnitemcollectionmetrics;
    return result;
  }

  BatchWriteItemInput._();

  factory BatchWriteItemInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchWriteItemInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchWriteItemInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..m<$core.String, $core.String>(
        247720687, _omitFieldNames ? '' : 'requestitems',
        entryClassName: 'BatchWriteItemInput.RequestitemsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<ReturnItemCollectionMetrics>(
        255507354, _omitFieldNames ? '' : 'returnitemcollectionmetrics',
        enumValues: ReturnItemCollectionMetrics.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchWriteItemInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchWriteItemInput copyWith(void Function(BatchWriteItemInput) updates) =>
      super.copyWith((message) => updates(message as BatchWriteItemInput))
          as BatchWriteItemInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchWriteItemInput create() => BatchWriteItemInput._();
  @$core.override
  BatchWriteItemInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchWriteItemInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchWriteItemInput>(create);
  static BatchWriteItemInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(247720687)
  $pb.PbMap<$core.String, $core.String> get requestitems => $_getMap(1);

  @$pb.TagNumber(255507354)
  ReturnItemCollectionMetrics get returnitemcollectionmetrics => $_getN(2);
  @$pb.TagNumber(255507354)
  set returnitemcollectionmetrics(ReturnItemCollectionMetrics value) =>
      $_setField(255507354, value);
  @$pb.TagNumber(255507354)
  $core.bool hasReturnitemcollectionmetrics() => $_has(2);
  @$pb.TagNumber(255507354)
  void clearReturnitemcollectionmetrics() => $_clearField(255507354);
}

class BatchWriteItemOutput extends $pb.GeneratedMessage {
  factory BatchWriteItemOutput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        unprocesseditems,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        itemcollectionmetrics,
    $core.Iterable<ConsumedCapacity>? consumedcapacity,
  }) {
    final result = create();
    if (unprocesseditems != null)
      result.unprocesseditems.addEntries(unprocesseditems);
    if (itemcollectionmetrics != null)
      result.itemcollectionmetrics.addEntries(itemcollectionmetrics);
    if (consumedcapacity != null)
      result.consumedcapacity.addAll(consumedcapacity);
    return result;
  }

  BatchWriteItemOutput._();

  factory BatchWriteItemOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchWriteItemOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchWriteItemOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        173269603, _omitFieldNames ? '' : 'unprocesseditems',
        entryClassName: 'BatchWriteItemOutput.UnprocesseditemsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, $core.String>(
        185317452, _omitFieldNames ? '' : 'itemcollectionmetrics',
        entryClassName: 'BatchWriteItemOutput.ItemcollectionmetricsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..pPM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchWriteItemOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchWriteItemOutput copyWith(void Function(BatchWriteItemOutput) updates) =>
      super.copyWith((message) => updates(message as BatchWriteItemOutput))
          as BatchWriteItemOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchWriteItemOutput create() => BatchWriteItemOutput._();
  @$core.override
  BatchWriteItemOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchWriteItemOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchWriteItemOutput>(create);
  static BatchWriteItemOutput? _defaultInstance;

  @$pb.TagNumber(173269603)
  $pb.PbMap<$core.String, $core.String> get unprocesseditems => $_getMap(0);

  @$pb.TagNumber(185317452)
  $pb.PbMap<$core.String, $core.String> get itemcollectionmetrics =>
      $_getMap(1);

  @$pb.TagNumber(449336620)
  $pb.PbList<ConsumedCapacity> get consumedcapacity => $_getList(2);
}

class BillingModeSummary extends $pb.GeneratedMessage {
  factory BillingModeSummary({
    $core.String? lastupdatetopayperrequestdatetime,
    BillingMode? billingmode,
  }) {
    final result = create();
    if (lastupdatetopayperrequestdatetime != null)
      result.lastupdatetopayperrequestdatetime =
          lastupdatetopayperrequestdatetime;
    if (billingmode != null) result.billingmode = billingmode;
    return result;
  }

  BillingModeSummary._();

  factory BillingModeSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BillingModeSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BillingModeSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(40315649, _omitFieldNames ? '' : 'lastupdatetopayperrequestdatetime')
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingModeSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BillingModeSummary copyWith(void Function(BillingModeSummary) updates) =>
      super.copyWith((message) => updates(message as BillingModeSummary))
          as BillingModeSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BillingModeSummary create() => BillingModeSummary._();
  @$core.override
  BillingModeSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BillingModeSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BillingModeSummary>(create);
  static BillingModeSummary? _defaultInstance;

  @$pb.TagNumber(40315649)
  $core.String get lastupdatetopayperrequestdatetime => $_getSZ(0);
  @$pb.TagNumber(40315649)
  set lastupdatetopayperrequestdatetime($core.String value) =>
      $_setString(0, value);
  @$pb.TagNumber(40315649)
  $core.bool hasLastupdatetopayperrequestdatetime() => $_has(0);
  @$pb.TagNumber(40315649)
  void clearLastupdatetopayperrequestdatetime() => $_clearField(40315649);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(1);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(1);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);
}

class CancellationReason extends $pb.GeneratedMessage {
  factory CancellationReason({
    $core.String? message,
    $core.String? code,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (code != null) result.code = code;
    if (item != null) result.item.addEntries(item);
    return result;
  }

  CancellationReason._();

  factory CancellationReason.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancellationReason.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancellationReason',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(425572629, _omitFieldNames ? '' : 'code')
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'CancellationReason.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancellationReason clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancellationReason copyWith(void Function(CancellationReason) updates) =>
      super.copyWith((message) => updates(message as CancellationReason))
          as CancellationReason;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancellationReason create() => CancellationReason._();
  @$core.override
  CancellationReason createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancellationReason getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancellationReason>(create);
  static CancellationReason? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(425572629)
  $core.String get code => $_getSZ(1);
  @$pb.TagNumber(425572629)
  set code($core.String value) => $_setString(1, value);
  @$pb.TagNumber(425572629)
  $core.bool hasCode() => $_has(1);
  @$pb.TagNumber(425572629)
  void clearCode() => $_clearField(425572629);

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(2);
}

class Capacity extends $pb.GeneratedMessage {
  factory Capacity({
    $core.double? writecapacityunits,
    $core.double? readcapacityunits,
    $core.double? capacityunits,
  }) {
    final result = create();
    if (writecapacityunits != null)
      result.writecapacityunits = writecapacityunits;
    if (readcapacityunits != null) result.readcapacityunits = readcapacityunits;
    if (capacityunits != null) result.capacityunits = capacityunits;
    return result;
  }

  Capacity._();

  factory Capacity.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Capacity.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Capacity',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aD(26938032, _omitFieldNames ? '' : 'writecapacityunits')
    ..aD(43945489, _omitFieldNames ? '' : 'readcapacityunits')
    ..aD(521922929, _omitFieldNames ? '' : 'capacityunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Capacity clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Capacity copyWith(void Function(Capacity) updates) =>
      super.copyWith((message) => updates(message as Capacity)) as Capacity;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Capacity create() => Capacity._();
  @$core.override
  Capacity createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Capacity getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Capacity>(create);
  static Capacity? _defaultInstance;

  @$pb.TagNumber(26938032)
  $core.double get writecapacityunits => $_getN(0);
  @$pb.TagNumber(26938032)
  set writecapacityunits($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(26938032)
  $core.bool hasWritecapacityunits() => $_has(0);
  @$pb.TagNumber(26938032)
  void clearWritecapacityunits() => $_clearField(26938032);

  @$pb.TagNumber(43945489)
  $core.double get readcapacityunits => $_getN(1);
  @$pb.TagNumber(43945489)
  set readcapacityunits($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(43945489)
  $core.bool hasReadcapacityunits() => $_has(1);
  @$pb.TagNumber(43945489)
  void clearReadcapacityunits() => $_clearField(43945489);

  @$pb.TagNumber(521922929)
  $core.double get capacityunits => $_getN(2);
  @$pb.TagNumber(521922929)
  set capacityunits($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(521922929)
  $core.bool hasCapacityunits() => $_has(2);
  @$pb.TagNumber(521922929)
  void clearCapacityunits() => $_clearField(521922929);
}

class Condition extends $pb.GeneratedMessage {
  factory Condition({
    ComparisonOperator? comparisonoperator,
    $core.Iterable<AttributeValue>? attributevaluelist,
  }) {
    final result = create();
    if (comparisonoperator != null)
      result.comparisonoperator = comparisonoperator;
    if (attributevaluelist != null)
      result.attributevaluelist.addAll(attributevaluelist);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ComparisonOperator>(
        7059861, _omitFieldNames ? '' : 'comparisonoperator',
        enumValues: ComparisonOperator.values)
    ..pPM<AttributeValue>(
        157424013, _omitFieldNames ? '' : 'attributevaluelist',
        subBuilder: AttributeValue.create)
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

  @$pb.TagNumber(7059861)
  ComparisonOperator get comparisonoperator => $_getN(0);
  @$pb.TagNumber(7059861)
  set comparisonoperator(ComparisonOperator value) =>
      $_setField(7059861, value);
  @$pb.TagNumber(7059861)
  $core.bool hasComparisonoperator() => $_has(0);
  @$pb.TagNumber(7059861)
  void clearComparisonoperator() => $_clearField(7059861);

  @$pb.TagNumber(157424013)
  $pb.PbList<AttributeValue> get attributevaluelist => $_getList(1);
}

class ConditionCheck extends $pb.GeneratedMessage {
  factory ConditionCheck({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? key,
    $core.String? tablename,
    $core.String? conditionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (key != null) result.key.addEntries(key);
    if (tablename != null) result.tablename = tablename;
    if (conditionexpression != null)
      result.conditionexpression = conditionexpression;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    return result;
  }

  ConditionCheck._();

  factory ConditionCheck.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConditionCheck.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConditionCheck',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'ConditionCheck.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, AttributeValue>(219859213, _omitFieldNames ? '' : 'key',
        entryClassName: 'ConditionCheck.KeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(409657405, _omitFieldNames ? '' : 'conditionexpression')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'ConditionCheck.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConditionCheck clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConditionCheck copyWith(void Function(ConditionCheck) updates) =>
      super.copyWith((message) => updates(message as ConditionCheck))
          as ConditionCheck;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConditionCheck create() => ConditionCheck._();
  @$core.override
  ConditionCheck createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConditionCheck getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConditionCheck>(create);
  static ConditionCheck? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(1);

  @$pb.TagNumber(219859213)
  $pb.PbMap<$core.String, AttributeValue> get key => $_getMap(2);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(409657405)
  $core.String get conditionexpression => $_getSZ(4);
  @$pb.TagNumber(409657405)
  set conditionexpression($core.String value) => $_setString(4, value);
  @$pb.TagNumber(409657405)
  $core.bool hasConditionexpression() => $_has(4);
  @$pb.TagNumber(409657405)
  void clearConditionexpression() => $_clearField(409657405);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(5);
}

class ConditionalCheckFailedException extends $pb.GeneratedMessage {
  factory ConditionalCheckFailedException({
    $core.String? message,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (item != null) result.item.addEntries(item);
    return result;
  }

  ConditionalCheckFailedException._();

  factory ConditionalCheckFailedException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConditionalCheckFailedException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConditionalCheckFailedException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'ConditionalCheckFailedException.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConditionalCheckFailedException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConditionalCheckFailedException copyWith(
          void Function(ConditionalCheckFailedException) updates) =>
      super.copyWith(
              (message) => updates(message as ConditionalCheckFailedException))
          as ConditionalCheckFailedException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConditionalCheckFailedException create() =>
      ConditionalCheckFailedException._();
  @$core.override
  ConditionalCheckFailedException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConditionalCheckFailedException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConditionalCheckFailedException>(
          create);
  static ConditionalCheckFailedException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(1);
}

class ConsumedCapacity extends $pb.GeneratedMessage {
  factory ConsumedCapacity({
    $core.double? writecapacityunits,
    $core.double? readcapacityunits,
    $core.String? tablename,
    $core.Iterable<$core.MapEntry<$core.String, Capacity>>?
        localsecondaryindexes,
    Capacity? table,
    $core.Iterable<$core.MapEntry<$core.String, Capacity>>?
        globalsecondaryindexes,
    $core.double? capacityunits,
  }) {
    final result = create();
    if (writecapacityunits != null)
      result.writecapacityunits = writecapacityunits;
    if (readcapacityunits != null) result.readcapacityunits = readcapacityunits;
    if (tablename != null) result.tablename = tablename;
    if (localsecondaryindexes != null)
      result.localsecondaryindexes.addEntries(localsecondaryindexes);
    if (table != null) result.table = table;
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addEntries(globalsecondaryindexes);
    if (capacityunits != null) result.capacityunits = capacityunits;
    return result;
  }

  ConsumedCapacity._();

  factory ConsumedCapacity.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConsumedCapacity.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConsumedCapacity',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aD(26938032, _omitFieldNames ? '' : 'writecapacityunits')
    ..aD(43945489, _omitFieldNames ? '' : 'readcapacityunits')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..m<$core.String, Capacity>(
        362339959, _omitFieldNames ? '' : 'localsecondaryindexes',
        entryClassName: 'ConsumedCapacity.LocalsecondaryindexesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: Capacity.create,
        valueDefaultOrMaker: Capacity.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOM<Capacity>(386722688, _omitFieldNames ? '' : 'table',
        subBuilder: Capacity.create)
    ..m<$core.String, Capacity>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        entryClassName: 'ConsumedCapacity.GlobalsecondaryindexesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: Capacity.create,
        valueDefaultOrMaker: Capacity.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aD(521922929, _omitFieldNames ? '' : 'capacityunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConsumedCapacity clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConsumedCapacity copyWith(void Function(ConsumedCapacity) updates) =>
      super.copyWith((message) => updates(message as ConsumedCapacity))
          as ConsumedCapacity;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConsumedCapacity create() => ConsumedCapacity._();
  @$core.override
  ConsumedCapacity createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConsumedCapacity getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConsumedCapacity>(create);
  static ConsumedCapacity? _defaultInstance;

  @$pb.TagNumber(26938032)
  $core.double get writecapacityunits => $_getN(0);
  @$pb.TagNumber(26938032)
  set writecapacityunits($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(26938032)
  $core.bool hasWritecapacityunits() => $_has(0);
  @$pb.TagNumber(26938032)
  void clearWritecapacityunits() => $_clearField(26938032);

  @$pb.TagNumber(43945489)
  $core.double get readcapacityunits => $_getN(1);
  @$pb.TagNumber(43945489)
  set readcapacityunits($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(43945489)
  $core.bool hasReadcapacityunits() => $_has(1);
  @$pb.TagNumber(43945489)
  void clearReadcapacityunits() => $_clearField(43945489);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(2);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(2);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(362339959)
  $pb.PbMap<$core.String, Capacity> get localsecondaryindexes => $_getMap(3);

  @$pb.TagNumber(386722688)
  Capacity get table => $_getN(4);
  @$pb.TagNumber(386722688)
  set table(Capacity value) => $_setField(386722688, value);
  @$pb.TagNumber(386722688)
  $core.bool hasTable() => $_has(4);
  @$pb.TagNumber(386722688)
  void clearTable() => $_clearField(386722688);
  @$pb.TagNumber(386722688)
  Capacity ensureTable() => $_ensure(4);

  @$pb.TagNumber(409156905)
  $pb.PbMap<$core.String, Capacity> get globalsecondaryindexes => $_getMap(5);

  @$pb.TagNumber(521922929)
  $core.double get capacityunits => $_getN(6);
  @$pb.TagNumber(521922929)
  set capacityunits($core.double value) => $_setDouble(6, value);
  @$pb.TagNumber(521922929)
  $core.bool hasCapacityunits() => $_has(6);
  @$pb.TagNumber(521922929)
  void clearCapacityunits() => $_clearField(521922929);
}

class ContinuousBackupsDescription extends $pb.GeneratedMessage {
  factory ContinuousBackupsDescription({
    PointInTimeRecoveryDescription? pointintimerecoverydescription,
    ContinuousBackupsStatus? continuousbackupsstatus,
  }) {
    final result = create();
    if (pointintimerecoverydescription != null)
      result.pointintimerecoverydescription = pointintimerecoverydescription;
    if (continuousbackupsstatus != null)
      result.continuousbackupsstatus = continuousbackupsstatus;
    return result;
  }

  ContinuousBackupsDescription._();

  factory ContinuousBackupsDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ContinuousBackupsDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ContinuousBackupsDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<PointInTimeRecoveryDescription>(
        319300211, _omitFieldNames ? '' : 'pointintimerecoverydescription',
        subBuilder: PointInTimeRecoveryDescription.create)
    ..aE<ContinuousBackupsStatus>(
        501975808, _omitFieldNames ? '' : 'continuousbackupsstatus',
        enumValues: ContinuousBackupsStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContinuousBackupsDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContinuousBackupsDescription copyWith(
          void Function(ContinuousBackupsDescription) updates) =>
      super.copyWith(
              (message) => updates(message as ContinuousBackupsDescription))
          as ContinuousBackupsDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ContinuousBackupsDescription create() =>
      ContinuousBackupsDescription._();
  @$core.override
  ContinuousBackupsDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ContinuousBackupsDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ContinuousBackupsDescription>(create);
  static ContinuousBackupsDescription? _defaultInstance;

  @$pb.TagNumber(319300211)
  PointInTimeRecoveryDescription get pointintimerecoverydescription =>
      $_getN(0);
  @$pb.TagNumber(319300211)
  set pointintimerecoverydescription(PointInTimeRecoveryDescription value) =>
      $_setField(319300211, value);
  @$pb.TagNumber(319300211)
  $core.bool hasPointintimerecoverydescription() => $_has(0);
  @$pb.TagNumber(319300211)
  void clearPointintimerecoverydescription() => $_clearField(319300211);
  @$pb.TagNumber(319300211)
  PointInTimeRecoveryDescription ensurePointintimerecoverydescription() =>
      $_ensure(0);

  @$pb.TagNumber(501975808)
  ContinuousBackupsStatus get continuousbackupsstatus => $_getN(1);
  @$pb.TagNumber(501975808)
  set continuousbackupsstatus(ContinuousBackupsStatus value) =>
      $_setField(501975808, value);
  @$pb.TagNumber(501975808)
  $core.bool hasContinuousbackupsstatus() => $_has(1);
  @$pb.TagNumber(501975808)
  void clearContinuousbackupsstatus() => $_clearField(501975808);
}

class ContinuousBackupsUnavailableException extends $pb.GeneratedMessage {
  factory ContinuousBackupsUnavailableException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ContinuousBackupsUnavailableException._();

  factory ContinuousBackupsUnavailableException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ContinuousBackupsUnavailableException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ContinuousBackupsUnavailableException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContinuousBackupsUnavailableException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContinuousBackupsUnavailableException copyWith(
          void Function(ContinuousBackupsUnavailableException) updates) =>
      super.copyWith((message) =>
              updates(message as ContinuousBackupsUnavailableException))
          as ContinuousBackupsUnavailableException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ContinuousBackupsUnavailableException create() =>
      ContinuousBackupsUnavailableException._();
  @$core.override
  ContinuousBackupsUnavailableException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ContinuousBackupsUnavailableException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ContinuousBackupsUnavailableException>(create);
  static ContinuousBackupsUnavailableException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ContributorInsightsSummary extends $pb.GeneratedMessage {
  factory ContributorInsightsSummary({
    ContributorInsightsMode? contributorinsightsmode,
    $core.String? indexname,
    $core.String? tablename,
    ContributorInsightsStatus? contributorinsightsstatus,
  }) {
    final result = create();
    if (contributorinsightsmode != null)
      result.contributorinsightsmode = contributorinsightsmode;
    if (indexname != null) result.indexname = indexname;
    if (tablename != null) result.tablename = tablename;
    if (contributorinsightsstatus != null)
      result.contributorinsightsstatus = contributorinsightsstatus;
    return result;
  }

  ContributorInsightsSummary._();

  factory ContributorInsightsSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ContributorInsightsSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ContributorInsightsSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ContributorInsightsMode>(
        86700161, _omitFieldNames ? '' : 'contributorinsightsmode',
        enumValues: ContributorInsightsMode.values)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aE<ContributorInsightsStatus>(
        363347282, _omitFieldNames ? '' : 'contributorinsightsstatus',
        enumValues: ContributorInsightsStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContributorInsightsSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContributorInsightsSummary copyWith(
          void Function(ContributorInsightsSummary) updates) =>
      super.copyWith(
              (message) => updates(message as ContributorInsightsSummary))
          as ContributorInsightsSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ContributorInsightsSummary create() => ContributorInsightsSummary._();
  @$core.override
  ContributorInsightsSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ContributorInsightsSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ContributorInsightsSummary>(create);
  static ContributorInsightsSummary? _defaultInstance;

  @$pb.TagNumber(86700161)
  ContributorInsightsMode get contributorinsightsmode => $_getN(0);
  @$pb.TagNumber(86700161)
  set contributorinsightsmode(ContributorInsightsMode value) =>
      $_setField(86700161, value);
  @$pb.TagNumber(86700161)
  $core.bool hasContributorinsightsmode() => $_has(0);
  @$pb.TagNumber(86700161)
  void clearContributorinsightsmode() => $_clearField(86700161);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(2);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(2);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(363347282)
  ContributorInsightsStatus get contributorinsightsstatus => $_getN(3);
  @$pb.TagNumber(363347282)
  set contributorinsightsstatus(ContributorInsightsStatus value) =>
      $_setField(363347282, value);
  @$pb.TagNumber(363347282)
  $core.bool hasContributorinsightsstatus() => $_has(3);
  @$pb.TagNumber(363347282)
  void clearContributorinsightsstatus() => $_clearField(363347282);
}

class CreateBackupInput extends $pb.GeneratedMessage {
  factory CreateBackupInput({
    $core.String? tablename,
    $core.String? backupname,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    if (backupname != null) result.backupname = backupname;
    return result;
  }

  CreateBackupInput._();

  factory CreateBackupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateBackupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateBackupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(467693789, _omitFieldNames ? '' : 'backupname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBackupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBackupInput copyWith(void Function(CreateBackupInput) updates) =>
      super.copyWith((message) => updates(message as CreateBackupInput))
          as CreateBackupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateBackupInput create() => CreateBackupInput._();
  @$core.override
  CreateBackupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateBackupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateBackupInput>(create);
  static CreateBackupInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(467693789)
  $core.String get backupname => $_getSZ(1);
  @$pb.TagNumber(467693789)
  set backupname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(467693789)
  $core.bool hasBackupname() => $_has(1);
  @$pb.TagNumber(467693789)
  void clearBackupname() => $_clearField(467693789);
}

class CreateBackupOutput extends $pb.GeneratedMessage {
  factory CreateBackupOutput({
    BackupDetails? backupdetails,
  }) {
    final result = create();
    if (backupdetails != null) result.backupdetails = backupdetails;
    return result;
  }

  CreateBackupOutput._();

  factory CreateBackupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateBackupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateBackupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<BackupDetails>(383357432, _omitFieldNames ? '' : 'backupdetails',
        subBuilder: BackupDetails.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBackupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateBackupOutput copyWith(void Function(CreateBackupOutput) updates) =>
      super.copyWith((message) => updates(message as CreateBackupOutput))
          as CreateBackupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateBackupOutput create() => CreateBackupOutput._();
  @$core.override
  CreateBackupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateBackupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateBackupOutput>(create);
  static CreateBackupOutput? _defaultInstance;

  @$pb.TagNumber(383357432)
  BackupDetails get backupdetails => $_getN(0);
  @$pb.TagNumber(383357432)
  set backupdetails(BackupDetails value) => $_setField(383357432, value);
  @$pb.TagNumber(383357432)
  $core.bool hasBackupdetails() => $_has(0);
  @$pb.TagNumber(383357432)
  void clearBackupdetails() => $_clearField(383357432);
  @$pb.TagNumber(383357432)
  BackupDetails ensureBackupdetails() => $_ensure(0);
}

class CreateGlobalSecondaryIndexAction extends $pb.GeneratedMessage {
  factory CreateGlobalSecondaryIndexAction({
    ProvisionedThroughput? provisionedthroughput,
    $core.String? indexname,
    Projection? projection,
    WarmThroughput? warmthroughput,
    $core.Iterable<KeySchemaElement>? keyschema,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (indexname != null) result.indexname = indexname;
    if (projection != null) result.projection = projection;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  CreateGlobalSecondaryIndexAction._();

  factory CreateGlobalSecondaryIndexAction.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateGlobalSecondaryIndexAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateGlobalSecondaryIndexAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughput>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughput.create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<Projection>(105045921, _omitFieldNames ? '' : 'projection',
        subBuilder: Projection.create)
    ..aOM<WarmThroughput>(290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: WarmThroughput.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGlobalSecondaryIndexAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGlobalSecondaryIndexAction copyWith(
          void Function(CreateGlobalSecondaryIndexAction) updates) =>
      super.copyWith(
              (message) => updates(message as CreateGlobalSecondaryIndexAction))
          as CreateGlobalSecondaryIndexAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateGlobalSecondaryIndexAction create() =>
      CreateGlobalSecondaryIndexAction._();
  @$core.override
  CreateGlobalSecondaryIndexAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateGlobalSecondaryIndexAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateGlobalSecondaryIndexAction>(
          create);
  static CreateGlobalSecondaryIndexAction? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughput get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughput value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughput ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(105045921)
  Projection get projection => $_getN(2);
  @$pb.TagNumber(105045921)
  set projection(Projection value) => $_setField(105045921, value);
  @$pb.TagNumber(105045921)
  $core.bool hasProjection() => $_has(2);
  @$pb.TagNumber(105045921)
  void clearProjection() => $_clearField(105045921);
  @$pb.TagNumber(105045921)
  Projection ensureProjection() => $_ensure(2);

  @$pb.TagNumber(290598659)
  WarmThroughput get warmthroughput => $_getN(3);
  @$pb.TagNumber(290598659)
  set warmthroughput(WarmThroughput value) => $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(3);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  WarmThroughput ensureWarmthroughput() => $_ensure(3);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(4);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(5);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(5);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(5);
}

class CreateGlobalTableInput extends $pb.GeneratedMessage {
  factory CreateGlobalTableInput({
    $core.Iterable<Replica>? replicationgroup,
    $core.String? globaltablename,
  }) {
    final result = create();
    if (replicationgroup != null)
      result.replicationgroup.addAll(replicationgroup);
    if (globaltablename != null) result.globaltablename = globaltablename;
    return result;
  }

  CreateGlobalTableInput._();

  factory CreateGlobalTableInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateGlobalTableInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateGlobalTableInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<Replica>(190970785, _omitFieldNames ? '' : 'replicationgroup',
        subBuilder: Replica.create)
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGlobalTableInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGlobalTableInput copyWith(
          void Function(CreateGlobalTableInput) updates) =>
      super.copyWith((message) => updates(message as CreateGlobalTableInput))
          as CreateGlobalTableInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateGlobalTableInput create() => CreateGlobalTableInput._();
  @$core.override
  CreateGlobalTableInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateGlobalTableInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateGlobalTableInput>(create);
  static CreateGlobalTableInput? _defaultInstance;

  @$pb.TagNumber(190970785)
  $pb.PbList<Replica> get replicationgroup => $_getList(0);

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(1);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(1);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);
}

class CreateGlobalTableOutput extends $pb.GeneratedMessage {
  factory CreateGlobalTableOutput({
    GlobalTableDescription? globaltabledescription,
  }) {
    final result = create();
    if (globaltabledescription != null)
      result.globaltabledescription = globaltabledescription;
    return result;
  }

  CreateGlobalTableOutput._();

  factory CreateGlobalTableOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateGlobalTableOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateGlobalTableOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<GlobalTableDescription>(
        342116405, _omitFieldNames ? '' : 'globaltabledescription',
        subBuilder: GlobalTableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGlobalTableOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGlobalTableOutput copyWith(
          void Function(CreateGlobalTableOutput) updates) =>
      super.copyWith((message) => updates(message as CreateGlobalTableOutput))
          as CreateGlobalTableOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateGlobalTableOutput create() => CreateGlobalTableOutput._();
  @$core.override
  CreateGlobalTableOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateGlobalTableOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateGlobalTableOutput>(create);
  static CreateGlobalTableOutput? _defaultInstance;

  @$pb.TagNumber(342116405)
  GlobalTableDescription get globaltabledescription => $_getN(0);
  @$pb.TagNumber(342116405)
  set globaltabledescription(GlobalTableDescription value) =>
      $_setField(342116405, value);
  @$pb.TagNumber(342116405)
  $core.bool hasGlobaltabledescription() => $_has(0);
  @$pb.TagNumber(342116405)
  void clearGlobaltabledescription() => $_clearField(342116405);
  @$pb.TagNumber(342116405)
  GlobalTableDescription ensureGlobaltabledescription() => $_ensure(0);
}

class CreateGlobalTableWitnessGroupMemberAction extends $pb.GeneratedMessage {
  factory CreateGlobalTableWitnessGroupMemberAction({
    $core.String? regionname,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    return result;
  }

  CreateGlobalTableWitnessGroupMemberAction._();

  factory CreateGlobalTableWitnessGroupMemberAction.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateGlobalTableWitnessGroupMemberAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateGlobalTableWitnessGroupMemberAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGlobalTableWitnessGroupMemberAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGlobalTableWitnessGroupMemberAction copyWith(
          void Function(CreateGlobalTableWitnessGroupMemberAction) updates) =>
      super.copyWith((message) =>
              updates(message as CreateGlobalTableWitnessGroupMemberAction))
          as CreateGlobalTableWitnessGroupMemberAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateGlobalTableWitnessGroupMemberAction create() =>
      CreateGlobalTableWitnessGroupMemberAction._();
  @$core.override
  CreateGlobalTableWitnessGroupMemberAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateGlobalTableWitnessGroupMemberAction getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateGlobalTableWitnessGroupMemberAction>(create);
  static CreateGlobalTableWitnessGroupMemberAction? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);
}

class CreateReplicaAction extends $pb.GeneratedMessage {
  factory CreateReplicaAction({
    $core.String? regionname,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    return result;
  }

  CreateReplicaAction._();

  factory CreateReplicaAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateReplicaAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateReplicaAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateReplicaAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateReplicaAction copyWith(void Function(CreateReplicaAction) updates) =>
      super.copyWith((message) => updates(message as CreateReplicaAction))
          as CreateReplicaAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateReplicaAction create() => CreateReplicaAction._();
  @$core.override
  CreateReplicaAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateReplicaAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateReplicaAction>(create);
  static CreateReplicaAction? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);
}

class CreateReplicationGroupMemberAction extends $pb.GeneratedMessage {
  factory CreateReplicationGroupMemberAction({
    $core.String? regionname,
    OnDemandThroughputOverride? ondemandthroughputoverride,
    $core.Iterable<ReplicaGlobalSecondaryIndex>? globalsecondaryindexes,
    ProvisionedThroughputOverride? provisionedthroughputoverride,
    TableClass? tableclassoverride,
    $core.String? kmsmasterkeyid,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    if (ondemandthroughputoverride != null)
      result.ondemandthroughputoverride = ondemandthroughputoverride;
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addAll(globalsecondaryindexes);
    if (provisionedthroughputoverride != null)
      result.provisionedthroughputoverride = provisionedthroughputoverride;
    if (tableclassoverride != null)
      result.tableclassoverride = tableclassoverride;
    if (kmsmasterkeyid != null) result.kmsmasterkeyid = kmsmasterkeyid;
    return result;
  }

  CreateReplicationGroupMemberAction._();

  factory CreateReplicationGroupMemberAction.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateReplicationGroupMemberAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateReplicationGroupMemberAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..aOM<OnDemandThroughputOverride>(
        317165234, _omitFieldNames ? '' : 'ondemandthroughputoverride',
        subBuilder: OnDemandThroughputOverride.create)
    ..pPM<ReplicaGlobalSecondaryIndex>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        subBuilder: ReplicaGlobalSecondaryIndex.create)
    ..aOM<ProvisionedThroughputOverride>(
        413332116, _omitFieldNames ? '' : 'provisionedthroughputoverride',
        subBuilder: ProvisionedThroughputOverride.create)
    ..aE<TableClass>(415569842, _omitFieldNames ? '' : 'tableclassoverride',
        enumValues: TableClass.values)
    ..aOS(521459443, _omitFieldNames ? '' : 'kmsmasterkeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateReplicationGroupMemberAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateReplicationGroupMemberAction copyWith(
          void Function(CreateReplicationGroupMemberAction) updates) =>
      super.copyWith((message) =>
              updates(message as CreateReplicationGroupMemberAction))
          as CreateReplicationGroupMemberAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateReplicationGroupMemberAction create() =>
      CreateReplicationGroupMemberAction._();
  @$core.override
  CreateReplicationGroupMemberAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateReplicationGroupMemberAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateReplicationGroupMemberAction>(
          create);
  static CreateReplicationGroupMemberAction? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride get ondemandthroughputoverride => $_getN(1);
  @$pb.TagNumber(317165234)
  set ondemandthroughputoverride(OnDemandThroughputOverride value) =>
      $_setField(317165234, value);
  @$pb.TagNumber(317165234)
  $core.bool hasOndemandthroughputoverride() => $_has(1);
  @$pb.TagNumber(317165234)
  void clearOndemandthroughputoverride() => $_clearField(317165234);
  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride ensureOndemandthroughputoverride() => $_ensure(1);

  @$pb.TagNumber(409156905)
  $pb.PbList<ReplicaGlobalSecondaryIndex> get globalsecondaryindexes =>
      $_getList(2);

  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride get provisionedthroughputoverride => $_getN(3);
  @$pb.TagNumber(413332116)
  set provisionedthroughputoverride(ProvisionedThroughputOverride value) =>
      $_setField(413332116, value);
  @$pb.TagNumber(413332116)
  $core.bool hasProvisionedthroughputoverride() => $_has(3);
  @$pb.TagNumber(413332116)
  void clearProvisionedthroughputoverride() => $_clearField(413332116);
  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride ensureProvisionedthroughputoverride() =>
      $_ensure(3);

  @$pb.TagNumber(415569842)
  TableClass get tableclassoverride => $_getN(4);
  @$pb.TagNumber(415569842)
  set tableclassoverride(TableClass value) => $_setField(415569842, value);
  @$pb.TagNumber(415569842)
  $core.bool hasTableclassoverride() => $_has(4);
  @$pb.TagNumber(415569842)
  void clearTableclassoverride() => $_clearField(415569842);

  @$pb.TagNumber(521459443)
  $core.String get kmsmasterkeyid => $_getSZ(5);
  @$pb.TagNumber(521459443)
  set kmsmasterkeyid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(521459443)
  $core.bool hasKmsmasterkeyid() => $_has(5);
  @$pb.TagNumber(521459443)
  void clearKmsmasterkeyid() => $_clearField(521459443);
}

class CreateTableInput extends $pb.GeneratedMessage {
  factory CreateTableInput({
    ProvisionedThroughput? provisionedthroughput,
    GlobalTableSettingsReplicationMode? globaltablesettingsreplicationmode,
    $core.String? resourcepolicy,
    SSESpecification? ssespecification,
    BillingMode? billingmode,
    $core.bool? deletionprotectionenabled,
    $core.String? tablename,
    WarmThroughput? warmthroughput,
    $core.Iterable<KeySchemaElement>? keyschema,
    TableClass? tableclass,
    $core.Iterable<LocalSecondaryIndex>? localsecondaryindexes,
    $core.Iterable<Tag>? tags,
    StreamSpecification? streamspecification,
    $core.Iterable<GlobalSecondaryIndex>? globalsecondaryindexes,
    $core.Iterable<AttributeDefinition>? attributedefinitions,
    $core.String? globaltablesourcearn,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (globaltablesettingsreplicationmode != null)
      result.globaltablesettingsreplicationmode =
          globaltablesettingsreplicationmode;
    if (resourcepolicy != null) result.resourcepolicy = resourcepolicy;
    if (ssespecification != null) result.ssespecification = ssespecification;
    if (billingmode != null) result.billingmode = billingmode;
    if (deletionprotectionenabled != null)
      result.deletionprotectionenabled = deletionprotectionenabled;
    if (tablename != null) result.tablename = tablename;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (tableclass != null) result.tableclass = tableclass;
    if (localsecondaryindexes != null)
      result.localsecondaryindexes.addAll(localsecondaryindexes);
    if (tags != null) result.tags.addAll(tags);
    if (streamspecification != null)
      result.streamspecification = streamspecification;
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addAll(globalsecondaryindexes);
    if (attributedefinitions != null)
      result.attributedefinitions.addAll(attributedefinitions);
    if (globaltablesourcearn != null)
      result.globaltablesourcearn = globaltablesourcearn;
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  CreateTableInput._();

  factory CreateTableInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTableInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTableInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughput>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughput.create)
    ..aE<GlobalTableSettingsReplicationMode>(
        10446577, _omitFieldNames ? '' : 'globaltablesettingsreplicationmode',
        enumValues: GlobalTableSettingsReplicationMode.values)
    ..aOS(15747632, _omitFieldNames ? '' : 'resourcepolicy')
    ..aOM<SSESpecification>(31692444, _omitFieldNames ? '' : 'ssespecification',
        subBuilder: SSESpecification.create)
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aOB(259418450, _omitFieldNames ? '' : 'deletionprotectionenabled')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<WarmThroughput>(290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: WarmThroughput.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..aE<TableClass>(342890498, _omitFieldNames ? '' : 'tableclass',
        enumValues: TableClass.values)
    ..pPM<LocalSecondaryIndex>(
        362339959, _omitFieldNames ? '' : 'localsecondaryindexes',
        subBuilder: LocalSecondaryIndex.create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOM<StreamSpecification>(
        403922627, _omitFieldNames ? '' : 'streamspecification',
        subBuilder: StreamSpecification.create)
    ..pPM<GlobalSecondaryIndex>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        subBuilder: GlobalSecondaryIndex.create)
    ..pPM<AttributeDefinition>(
        414687108, _omitFieldNames ? '' : 'attributedefinitions',
        subBuilder: AttributeDefinition.create)
    ..aOS(443614787, _omitFieldNames ? '' : 'globaltablesourcearn')
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTableInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTableInput copyWith(void Function(CreateTableInput) updates) =>
      super.copyWith((message) => updates(message as CreateTableInput))
          as CreateTableInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTableInput create() => CreateTableInput._();
  @$core.override
  CreateTableInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTableInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTableInput>(create);
  static CreateTableInput? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughput get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughput value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughput ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(10446577)
  GlobalTableSettingsReplicationMode get globaltablesettingsreplicationmode =>
      $_getN(1);
  @$pb.TagNumber(10446577)
  set globaltablesettingsreplicationmode(
          GlobalTableSettingsReplicationMode value) =>
      $_setField(10446577, value);
  @$pb.TagNumber(10446577)
  $core.bool hasGlobaltablesettingsreplicationmode() => $_has(1);
  @$pb.TagNumber(10446577)
  void clearGlobaltablesettingsreplicationmode() => $_clearField(10446577);

  @$pb.TagNumber(15747632)
  $core.String get resourcepolicy => $_getSZ(2);
  @$pb.TagNumber(15747632)
  set resourcepolicy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(15747632)
  $core.bool hasResourcepolicy() => $_has(2);
  @$pb.TagNumber(15747632)
  void clearResourcepolicy() => $_clearField(15747632);

  @$pb.TagNumber(31692444)
  SSESpecification get ssespecification => $_getN(3);
  @$pb.TagNumber(31692444)
  set ssespecification(SSESpecification value) => $_setField(31692444, value);
  @$pb.TagNumber(31692444)
  $core.bool hasSsespecification() => $_has(3);
  @$pb.TagNumber(31692444)
  void clearSsespecification() => $_clearField(31692444);
  @$pb.TagNumber(31692444)
  SSESpecification ensureSsespecification() => $_ensure(3);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(4);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(4);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(259418450)
  $core.bool get deletionprotectionenabled => $_getBF(5);
  @$pb.TagNumber(259418450)
  set deletionprotectionenabled($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(259418450)
  $core.bool hasDeletionprotectionenabled() => $_has(5);
  @$pb.TagNumber(259418450)
  void clearDeletionprotectionenabled() => $_clearField(259418450);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(6);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(6, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(6);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(290598659)
  WarmThroughput get warmthroughput => $_getN(7);
  @$pb.TagNumber(290598659)
  set warmthroughput(WarmThroughput value) => $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(7);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  WarmThroughput ensureWarmthroughput() => $_ensure(7);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(8);

  @$pb.TagNumber(342890498)
  TableClass get tableclass => $_getN(9);
  @$pb.TagNumber(342890498)
  set tableclass(TableClass value) => $_setField(342890498, value);
  @$pb.TagNumber(342890498)
  $core.bool hasTableclass() => $_has(9);
  @$pb.TagNumber(342890498)
  void clearTableclass() => $_clearField(342890498);

  @$pb.TagNumber(362339959)
  $pb.PbList<LocalSecondaryIndex> get localsecondaryindexes => $_getList(10);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(11);

  @$pb.TagNumber(403922627)
  StreamSpecification get streamspecification => $_getN(12);
  @$pb.TagNumber(403922627)
  set streamspecification(StreamSpecification value) =>
      $_setField(403922627, value);
  @$pb.TagNumber(403922627)
  $core.bool hasStreamspecification() => $_has(12);
  @$pb.TagNumber(403922627)
  void clearStreamspecification() => $_clearField(403922627);
  @$pb.TagNumber(403922627)
  StreamSpecification ensureStreamspecification() => $_ensure(12);

  @$pb.TagNumber(409156905)
  $pb.PbList<GlobalSecondaryIndex> get globalsecondaryindexes => $_getList(13);

  @$pb.TagNumber(414687108)
  $pb.PbList<AttributeDefinition> get attributedefinitions => $_getList(14);

  @$pb.TagNumber(443614787)
  $core.String get globaltablesourcearn => $_getSZ(15);
  @$pb.TagNumber(443614787)
  set globaltablesourcearn($core.String value) => $_setString(15, value);
  @$pb.TagNumber(443614787)
  $core.bool hasGlobaltablesourcearn() => $_has(15);
  @$pb.TagNumber(443614787)
  void clearGlobaltablesourcearn() => $_clearField(443614787);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(16);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(16);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(16);
}

class CreateTableOutput extends $pb.GeneratedMessage {
  factory CreateTableOutput({
    TableDescription? tabledescription,
  }) {
    final result = create();
    if (tabledescription != null) result.tabledescription = tabledescription;
    return result;
  }

  CreateTableOutput._();

  factory CreateTableOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTableOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTableOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TableDescription>(19962388, _omitFieldNames ? '' : 'tabledescription',
        subBuilder: TableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTableOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTableOutput copyWith(void Function(CreateTableOutput) updates) =>
      super.copyWith((message) => updates(message as CreateTableOutput))
          as CreateTableOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTableOutput create() => CreateTableOutput._();
  @$core.override
  CreateTableOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTableOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTableOutput>(create);
  static CreateTableOutput? _defaultInstance;

  @$pb.TagNumber(19962388)
  TableDescription get tabledescription => $_getN(0);
  @$pb.TagNumber(19962388)
  set tabledescription(TableDescription value) => $_setField(19962388, value);
  @$pb.TagNumber(19962388)
  $core.bool hasTabledescription() => $_has(0);
  @$pb.TagNumber(19962388)
  void clearTabledescription() => $_clearField(19962388);
  @$pb.TagNumber(19962388)
  TableDescription ensureTabledescription() => $_ensure(0);
}

class CsvOptions extends $pb.GeneratedMessage {
  factory CsvOptions({
    $core.String? delimiter,
    $core.Iterable<$core.String>? headerlist,
  }) {
    final result = create();
    if (delimiter != null) result.delimiter = delimiter;
    if (headerlist != null) result.headerlist.addAll(headerlist);
    return result;
  }

  CsvOptions._();

  factory CsvOptions.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CsvOptions.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CsvOptions',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(302132379, _omitFieldNames ? '' : 'delimiter')
    ..pPS(424018665, _omitFieldNames ? '' : 'headerlist')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CsvOptions clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CsvOptions copyWith(void Function(CsvOptions) updates) =>
      super.copyWith((message) => updates(message as CsvOptions)) as CsvOptions;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CsvOptions create() => CsvOptions._();
  @$core.override
  CsvOptions createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CsvOptions getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CsvOptions>(create);
  static CsvOptions? _defaultInstance;

  @$pb.TagNumber(302132379)
  $core.String get delimiter => $_getSZ(0);
  @$pb.TagNumber(302132379)
  set delimiter($core.String value) => $_setString(0, value);
  @$pb.TagNumber(302132379)
  $core.bool hasDelimiter() => $_has(0);
  @$pb.TagNumber(302132379)
  void clearDelimiter() => $_clearField(302132379);

  @$pb.TagNumber(424018665)
  $pb.PbList<$core.String> get headerlist => $_getList(1);
}

class Delete extends $pb.GeneratedMessage {
  factory Delete({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? key,
    $core.String? tablename,
    $core.String? conditionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (key != null) result.key.addEntries(key);
    if (tablename != null) result.tablename = tablename;
    if (conditionexpression != null)
      result.conditionexpression = conditionexpression;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    return result;
  }

  Delete._();

  factory Delete.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Delete.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Delete',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'Delete.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, AttributeValue>(219859213, _omitFieldNames ? '' : 'key',
        entryClassName: 'Delete.KeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(409657405, _omitFieldNames ? '' : 'conditionexpression')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'Delete.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Delete clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Delete copyWith(void Function(Delete) updates) =>
      super.copyWith((message) => updates(message as Delete)) as Delete;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Delete create() => Delete._();
  @$core.override
  Delete createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Delete getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Delete>(create);
  static Delete? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(1);

  @$pb.TagNumber(219859213)
  $pb.PbMap<$core.String, AttributeValue> get key => $_getMap(2);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(409657405)
  $core.String get conditionexpression => $_getSZ(4);
  @$pb.TagNumber(409657405)
  set conditionexpression($core.String value) => $_setString(4, value);
  @$pb.TagNumber(409657405)
  $core.bool hasConditionexpression() => $_has(4);
  @$pb.TagNumber(409657405)
  void clearConditionexpression() => $_clearField(409657405);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(5);
}

class DeleteBackupInput extends $pb.GeneratedMessage {
  factory DeleteBackupInput({
    $core.String? backuparn,
  }) {
    final result = create();
    if (backuparn != null) result.backuparn = backuparn;
    return result;
  }

  DeleteBackupInput._();

  factory DeleteBackupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteBackupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteBackupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(370874339, _omitFieldNames ? '' : 'backuparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteBackupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteBackupInput copyWith(void Function(DeleteBackupInput) updates) =>
      super.copyWith((message) => updates(message as DeleteBackupInput))
          as DeleteBackupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteBackupInput create() => DeleteBackupInput._();
  @$core.override
  DeleteBackupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteBackupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteBackupInput>(create);
  static DeleteBackupInput? _defaultInstance;

  @$pb.TagNumber(370874339)
  $core.String get backuparn => $_getSZ(0);
  @$pb.TagNumber(370874339)
  set backuparn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(370874339)
  $core.bool hasBackuparn() => $_has(0);
  @$pb.TagNumber(370874339)
  void clearBackuparn() => $_clearField(370874339);
}

class DeleteBackupOutput extends $pb.GeneratedMessage {
  factory DeleteBackupOutput({
    BackupDescription? backupdescription,
  }) {
    final result = create();
    if (backupdescription != null) result.backupdescription = backupdescription;
    return result;
  }

  DeleteBackupOutput._();

  factory DeleteBackupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteBackupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteBackupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<BackupDescription>(
        294936980, _omitFieldNames ? '' : 'backupdescription',
        subBuilder: BackupDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteBackupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteBackupOutput copyWith(void Function(DeleteBackupOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteBackupOutput))
          as DeleteBackupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteBackupOutput create() => DeleteBackupOutput._();
  @$core.override
  DeleteBackupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteBackupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteBackupOutput>(create);
  static DeleteBackupOutput? _defaultInstance;

  @$pb.TagNumber(294936980)
  BackupDescription get backupdescription => $_getN(0);
  @$pb.TagNumber(294936980)
  set backupdescription(BackupDescription value) =>
      $_setField(294936980, value);
  @$pb.TagNumber(294936980)
  $core.bool hasBackupdescription() => $_has(0);
  @$pb.TagNumber(294936980)
  void clearBackupdescription() => $_clearField(294936980);
  @$pb.TagNumber(294936980)
  BackupDescription ensureBackupdescription() => $_ensure(0);
}

class DeleteGlobalSecondaryIndexAction extends $pb.GeneratedMessage {
  factory DeleteGlobalSecondaryIndexAction({
    $core.String? indexname,
  }) {
    final result = create();
    if (indexname != null) result.indexname = indexname;
    return result;
  }

  DeleteGlobalSecondaryIndexAction._();

  factory DeleteGlobalSecondaryIndexAction.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteGlobalSecondaryIndexAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteGlobalSecondaryIndexAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGlobalSecondaryIndexAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGlobalSecondaryIndexAction copyWith(
          void Function(DeleteGlobalSecondaryIndexAction) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteGlobalSecondaryIndexAction))
          as DeleteGlobalSecondaryIndexAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteGlobalSecondaryIndexAction create() =>
      DeleteGlobalSecondaryIndexAction._();
  @$core.override
  DeleteGlobalSecondaryIndexAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteGlobalSecondaryIndexAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteGlobalSecondaryIndexAction>(
          create);
  static DeleteGlobalSecondaryIndexAction? _defaultInstance;

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(0);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(0);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);
}

class DeleteGlobalTableWitnessGroupMemberAction extends $pb.GeneratedMessage {
  factory DeleteGlobalTableWitnessGroupMemberAction({
    $core.String? regionname,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    return result;
  }

  DeleteGlobalTableWitnessGroupMemberAction._();

  factory DeleteGlobalTableWitnessGroupMemberAction.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteGlobalTableWitnessGroupMemberAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteGlobalTableWitnessGroupMemberAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGlobalTableWitnessGroupMemberAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGlobalTableWitnessGroupMemberAction copyWith(
          void Function(DeleteGlobalTableWitnessGroupMemberAction) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteGlobalTableWitnessGroupMemberAction))
          as DeleteGlobalTableWitnessGroupMemberAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteGlobalTableWitnessGroupMemberAction create() =>
      DeleteGlobalTableWitnessGroupMemberAction._();
  @$core.override
  DeleteGlobalTableWitnessGroupMemberAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteGlobalTableWitnessGroupMemberAction getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeleteGlobalTableWitnessGroupMemberAction>(create);
  static DeleteGlobalTableWitnessGroupMemberAction? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);
}

class DeleteItemInput extends $pb.GeneratedMessage {
  factory DeleteItemInput({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<$core.MapEntry<$core.String, ExpectedAttributeValue>>?
        expected,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    ConditionalOperator? conditionaloperator,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? key,
    ReturnItemCollectionMetrics? returnitemcollectionmetrics,
    $core.String? tablename,
    ReturnValue? returnvalues,
    $core.String? conditionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (expected != null) result.expected.addEntries(expected);
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (conditionaloperator != null)
      result.conditionaloperator = conditionaloperator;
    if (key != null) result.key.addEntries(key);
    if (returnitemcollectionmetrics != null)
      result.returnitemcollectionmetrics = returnitemcollectionmetrics;
    if (tablename != null) result.tablename = tablename;
    if (returnvalues != null) result.returnvalues = returnvalues;
    if (conditionexpression != null)
      result.conditionexpression = conditionexpression;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    return result;
  }

  DeleteItemInput._();

  factory DeleteItemInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteItemInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteItemInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..m<$core.String, ExpectedAttributeValue>(
        106557946, _omitFieldNames ? '' : 'expected',
        entryClassName: 'DeleteItemInput.ExpectedEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: ExpectedAttributeValue.create,
        valueDefaultOrMaker: ExpectedAttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'DeleteItemInput.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<ConditionalOperator>(
        172066260, _omitFieldNames ? '' : 'conditionaloperator',
        enumValues: ConditionalOperator.values)
    ..m<$core.String, AttributeValue>(219859213, _omitFieldNames ? '' : 'key',
        entryClassName: 'DeleteItemInput.KeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<ReturnItemCollectionMetrics>(
        255507354, _omitFieldNames ? '' : 'returnitemcollectionmetrics',
        enumValues: ReturnItemCollectionMetrics.values)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aE<ReturnValue>(402960198, _omitFieldNames ? '' : 'returnvalues',
        enumValues: ReturnValue.values)
    ..aOS(409657405, _omitFieldNames ? '' : 'conditionexpression')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'DeleteItemInput.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteItemInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteItemInput copyWith(void Function(DeleteItemInput) updates) =>
      super.copyWith((message) => updates(message as DeleteItemInput))
          as DeleteItemInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteItemInput create() => DeleteItemInput._();
  @$core.override
  DeleteItemInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteItemInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteItemInput>(create);
  static DeleteItemInput? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(1);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(1);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(106557946)
  $pb.PbMap<$core.String, ExpectedAttributeValue> get expected => $_getMap(2);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(3);

  @$pb.TagNumber(172066260)
  ConditionalOperator get conditionaloperator => $_getN(4);
  @$pb.TagNumber(172066260)
  set conditionaloperator(ConditionalOperator value) =>
      $_setField(172066260, value);
  @$pb.TagNumber(172066260)
  $core.bool hasConditionaloperator() => $_has(4);
  @$pb.TagNumber(172066260)
  void clearConditionaloperator() => $_clearField(172066260);

  @$pb.TagNumber(219859213)
  $pb.PbMap<$core.String, AttributeValue> get key => $_getMap(5);

  @$pb.TagNumber(255507354)
  ReturnItemCollectionMetrics get returnitemcollectionmetrics => $_getN(6);
  @$pb.TagNumber(255507354)
  set returnitemcollectionmetrics(ReturnItemCollectionMetrics value) =>
      $_setField(255507354, value);
  @$pb.TagNumber(255507354)
  $core.bool hasReturnitemcollectionmetrics() => $_has(6);
  @$pb.TagNumber(255507354)
  void clearReturnitemcollectionmetrics() => $_clearField(255507354);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(7);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(7, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(7);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(402960198)
  ReturnValue get returnvalues => $_getN(8);
  @$pb.TagNumber(402960198)
  set returnvalues(ReturnValue value) => $_setField(402960198, value);
  @$pb.TagNumber(402960198)
  $core.bool hasReturnvalues() => $_has(8);
  @$pb.TagNumber(402960198)
  void clearReturnvalues() => $_clearField(402960198);

  @$pb.TagNumber(409657405)
  $core.String get conditionexpression => $_getSZ(9);
  @$pb.TagNumber(409657405)
  set conditionexpression($core.String value) => $_setString(9, value);
  @$pb.TagNumber(409657405)
  $core.bool hasConditionexpression() => $_has(9);
  @$pb.TagNumber(409657405)
  void clearConditionexpression() => $_clearField(409657405);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(10);
}

class DeleteItemOutput extends $pb.GeneratedMessage {
  factory DeleteItemOutput({
    ItemCollectionMetrics? itemcollectionmetrics,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? attributes,
    ConsumedCapacity? consumedcapacity,
  }) {
    final result = create();
    if (itemcollectionmetrics != null)
      result.itemcollectionmetrics = itemcollectionmetrics;
    if (attributes != null) result.attributes.addEntries(attributes);
    if (consumedcapacity != null) result.consumedcapacity = consumedcapacity;
    return result;
  }

  DeleteItemOutput._();

  factory DeleteItemOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteItemOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteItemOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ItemCollectionMetrics>(
        185317452, _omitFieldNames ? '' : 'itemcollectionmetrics',
        subBuilder: ItemCollectionMetrics.create)
    ..m<$core.String, AttributeValue>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'DeleteItemOutput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteItemOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteItemOutput copyWith(void Function(DeleteItemOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteItemOutput))
          as DeleteItemOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteItemOutput create() => DeleteItemOutput._();
  @$core.override
  DeleteItemOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteItemOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteItemOutput>(create);
  static DeleteItemOutput? _defaultInstance;

  @$pb.TagNumber(185317452)
  ItemCollectionMetrics get itemcollectionmetrics => $_getN(0);
  @$pb.TagNumber(185317452)
  set itemcollectionmetrics(ItemCollectionMetrics value) =>
      $_setField(185317452, value);
  @$pb.TagNumber(185317452)
  $core.bool hasItemcollectionmetrics() => $_has(0);
  @$pb.TagNumber(185317452)
  void clearItemcollectionmetrics() => $_clearField(185317452);
  @$pb.TagNumber(185317452)
  ItemCollectionMetrics ensureItemcollectionmetrics() => $_ensure(0);

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, AttributeValue> get attributes => $_getMap(1);

  @$pb.TagNumber(449336620)
  ConsumedCapacity get consumedcapacity => $_getN(2);
  @$pb.TagNumber(449336620)
  set consumedcapacity(ConsumedCapacity value) => $_setField(449336620, value);
  @$pb.TagNumber(449336620)
  $core.bool hasConsumedcapacity() => $_has(2);
  @$pb.TagNumber(449336620)
  void clearConsumedcapacity() => $_clearField(449336620);
  @$pb.TagNumber(449336620)
  ConsumedCapacity ensureConsumedcapacity() => $_ensure(2);
}

class DeleteReplicaAction extends $pb.GeneratedMessage {
  factory DeleteReplicaAction({
    $core.String? regionname,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    return result;
  }

  DeleteReplicaAction._();

  factory DeleteReplicaAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteReplicaAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteReplicaAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteReplicaAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteReplicaAction copyWith(void Function(DeleteReplicaAction) updates) =>
      super.copyWith((message) => updates(message as DeleteReplicaAction))
          as DeleteReplicaAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteReplicaAction create() => DeleteReplicaAction._();
  @$core.override
  DeleteReplicaAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteReplicaAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteReplicaAction>(create);
  static DeleteReplicaAction? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);
}

class DeleteReplicationGroupMemberAction extends $pb.GeneratedMessage {
  factory DeleteReplicationGroupMemberAction({
    $core.String? regionname,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    return result;
  }

  DeleteReplicationGroupMemberAction._();

  factory DeleteReplicationGroupMemberAction.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteReplicationGroupMemberAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteReplicationGroupMemberAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteReplicationGroupMemberAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteReplicationGroupMemberAction copyWith(
          void Function(DeleteReplicationGroupMemberAction) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteReplicationGroupMemberAction))
          as DeleteReplicationGroupMemberAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteReplicationGroupMemberAction create() =>
      DeleteReplicationGroupMemberAction._();
  @$core.override
  DeleteReplicationGroupMemberAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteReplicationGroupMemberAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteReplicationGroupMemberAction>(
          create);
  static DeleteReplicationGroupMemberAction? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);
}

class DeleteRequest extends $pb.GeneratedMessage {
  factory DeleteRequest({
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? key,
  }) {
    final result = create();
    if (key != null) result.key.addEntries(key);
    return result;
  }

  DeleteRequest._();

  factory DeleteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, AttributeValue>(219859213, _omitFieldNames ? '' : 'key',
        entryClassName: 'DeleteRequest.KeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRequest copyWith(void Function(DeleteRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteRequest))
          as DeleteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRequest create() => DeleteRequest._();
  @$core.override
  DeleteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRequest>(create);
  static DeleteRequest? _defaultInstance;

  @$pb.TagNumber(219859213)
  $pb.PbMap<$core.String, AttributeValue> get key => $_getMap(0);
}

class DeleteResourcePolicyInput extends $pb.GeneratedMessage {
  factory DeleteResourcePolicyInput({
    $core.String? resourcearn,
    $core.String? expectedrevisionid,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (expectedrevisionid != null)
      result.expectedrevisionid = expectedrevisionid;
    return result;
  }

  DeleteResourcePolicyInput._();

  factory DeleteResourcePolicyInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteResourcePolicyInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteResourcePolicyInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(456463970, _omitFieldNames ? '' : 'expectedrevisionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourcePolicyInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourcePolicyInput copyWith(
          void Function(DeleteResourcePolicyInput) updates) =>
      super.copyWith((message) => updates(message as DeleteResourcePolicyInput))
          as DeleteResourcePolicyInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteResourcePolicyInput create() => DeleteResourcePolicyInput._();
  @$core.override
  DeleteResourcePolicyInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteResourcePolicyInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteResourcePolicyInput>(create);
  static DeleteResourcePolicyInput? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

  @$pb.TagNumber(456463970)
  $core.String get expectedrevisionid => $_getSZ(1);
  @$pb.TagNumber(456463970)
  set expectedrevisionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(456463970)
  $core.bool hasExpectedrevisionid() => $_has(1);
  @$pb.TagNumber(456463970)
  void clearExpectedrevisionid() => $_clearField(456463970);
}

class DeleteResourcePolicyOutput extends $pb.GeneratedMessage {
  factory DeleteResourcePolicyOutput({
    $core.String? revisionid,
  }) {
    final result = create();
    if (revisionid != null) result.revisionid = revisionid;
    return result;
  }

  DeleteResourcePolicyOutput._();

  factory DeleteResourcePolicyOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteResourcePolicyOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteResourcePolicyOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(499618182, _omitFieldNames ? '' : 'revisionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourcePolicyOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteResourcePolicyOutput copyWith(
          void Function(DeleteResourcePolicyOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteResourcePolicyOutput))
          as DeleteResourcePolicyOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteResourcePolicyOutput create() => DeleteResourcePolicyOutput._();
  @$core.override
  DeleteResourcePolicyOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteResourcePolicyOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteResourcePolicyOutput>(create);
  static DeleteResourcePolicyOutput? _defaultInstance;

  @$pb.TagNumber(499618182)
  $core.String get revisionid => $_getSZ(0);
  @$pb.TagNumber(499618182)
  set revisionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(499618182)
  $core.bool hasRevisionid() => $_has(0);
  @$pb.TagNumber(499618182)
  void clearRevisionid() => $_clearField(499618182);
}

class DeleteTableInput extends $pb.GeneratedMessage {
  factory DeleteTableInput({
    $core.String? tablename,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DeleteTableInput._();

  factory DeleteTableInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTableInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTableInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTableInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTableInput copyWith(void Function(DeleteTableInput) updates) =>
      super.copyWith((message) => updates(message as DeleteTableInput))
          as DeleteTableInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTableInput create() => DeleteTableInput._();
  @$core.override
  DeleteTableInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTableInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTableInput>(create);
  static DeleteTableInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DeleteTableOutput extends $pb.GeneratedMessage {
  factory DeleteTableOutput({
    TableDescription? tabledescription,
  }) {
    final result = create();
    if (tabledescription != null) result.tabledescription = tabledescription;
    return result;
  }

  DeleteTableOutput._();

  factory DeleteTableOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTableOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTableOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TableDescription>(19962388, _omitFieldNames ? '' : 'tabledescription',
        subBuilder: TableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTableOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTableOutput copyWith(void Function(DeleteTableOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteTableOutput))
          as DeleteTableOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTableOutput create() => DeleteTableOutput._();
  @$core.override
  DeleteTableOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTableOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTableOutput>(create);
  static DeleteTableOutput? _defaultInstance;

  @$pb.TagNumber(19962388)
  TableDescription get tabledescription => $_getN(0);
  @$pb.TagNumber(19962388)
  set tabledescription(TableDescription value) => $_setField(19962388, value);
  @$pb.TagNumber(19962388)
  $core.bool hasTabledescription() => $_has(0);
  @$pb.TagNumber(19962388)
  void clearTabledescription() => $_clearField(19962388);
  @$pb.TagNumber(19962388)
  TableDescription ensureTabledescription() => $_ensure(0);
}

class DescribeBackupInput extends $pb.GeneratedMessage {
  factory DescribeBackupInput({
    $core.String? backuparn,
  }) {
    final result = create();
    if (backuparn != null) result.backuparn = backuparn;
    return result;
  }

  DescribeBackupInput._();

  factory DescribeBackupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeBackupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeBackupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(370874339, _omitFieldNames ? '' : 'backuparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeBackupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeBackupInput copyWith(void Function(DescribeBackupInput) updates) =>
      super.copyWith((message) => updates(message as DescribeBackupInput))
          as DescribeBackupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeBackupInput create() => DescribeBackupInput._();
  @$core.override
  DescribeBackupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeBackupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeBackupInput>(create);
  static DescribeBackupInput? _defaultInstance;

  @$pb.TagNumber(370874339)
  $core.String get backuparn => $_getSZ(0);
  @$pb.TagNumber(370874339)
  set backuparn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(370874339)
  $core.bool hasBackuparn() => $_has(0);
  @$pb.TagNumber(370874339)
  void clearBackuparn() => $_clearField(370874339);
}

class DescribeBackupOutput extends $pb.GeneratedMessage {
  factory DescribeBackupOutput({
    BackupDescription? backupdescription,
  }) {
    final result = create();
    if (backupdescription != null) result.backupdescription = backupdescription;
    return result;
  }

  DescribeBackupOutput._();

  factory DescribeBackupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeBackupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeBackupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<BackupDescription>(
        294936980, _omitFieldNames ? '' : 'backupdescription',
        subBuilder: BackupDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeBackupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeBackupOutput copyWith(void Function(DescribeBackupOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeBackupOutput))
          as DescribeBackupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeBackupOutput create() => DescribeBackupOutput._();
  @$core.override
  DescribeBackupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeBackupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeBackupOutput>(create);
  static DescribeBackupOutput? _defaultInstance;

  @$pb.TagNumber(294936980)
  BackupDescription get backupdescription => $_getN(0);
  @$pb.TagNumber(294936980)
  set backupdescription(BackupDescription value) =>
      $_setField(294936980, value);
  @$pb.TagNumber(294936980)
  $core.bool hasBackupdescription() => $_has(0);
  @$pb.TagNumber(294936980)
  void clearBackupdescription() => $_clearField(294936980);
  @$pb.TagNumber(294936980)
  BackupDescription ensureBackupdescription() => $_ensure(0);
}

class DescribeContinuousBackupsInput extends $pb.GeneratedMessage {
  factory DescribeContinuousBackupsInput({
    $core.String? tablename,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DescribeContinuousBackupsInput._();

  factory DescribeContinuousBackupsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeContinuousBackupsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeContinuousBackupsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeContinuousBackupsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeContinuousBackupsInput copyWith(
          void Function(DescribeContinuousBackupsInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeContinuousBackupsInput))
          as DescribeContinuousBackupsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeContinuousBackupsInput create() =>
      DescribeContinuousBackupsInput._();
  @$core.override
  DescribeContinuousBackupsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeContinuousBackupsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeContinuousBackupsInput>(create);
  static DescribeContinuousBackupsInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeContinuousBackupsOutput extends $pb.GeneratedMessage {
  factory DescribeContinuousBackupsOutput({
    ContinuousBackupsDescription? continuousbackupsdescription,
  }) {
    final result = create();
    if (continuousbackupsdescription != null)
      result.continuousbackupsdescription = continuousbackupsdescription;
    return result;
  }

  DescribeContinuousBackupsOutput._();

  factory DescribeContinuousBackupsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeContinuousBackupsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeContinuousBackupsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ContinuousBackupsDescription>(
        446030650, _omitFieldNames ? '' : 'continuousbackupsdescription',
        subBuilder: ContinuousBackupsDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeContinuousBackupsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeContinuousBackupsOutput copyWith(
          void Function(DescribeContinuousBackupsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeContinuousBackupsOutput))
          as DescribeContinuousBackupsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeContinuousBackupsOutput create() =>
      DescribeContinuousBackupsOutput._();
  @$core.override
  DescribeContinuousBackupsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeContinuousBackupsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeContinuousBackupsOutput>(
          create);
  static DescribeContinuousBackupsOutput? _defaultInstance;

  @$pb.TagNumber(446030650)
  ContinuousBackupsDescription get continuousbackupsdescription => $_getN(0);
  @$pb.TagNumber(446030650)
  set continuousbackupsdescription(ContinuousBackupsDescription value) =>
      $_setField(446030650, value);
  @$pb.TagNumber(446030650)
  $core.bool hasContinuousbackupsdescription() => $_has(0);
  @$pb.TagNumber(446030650)
  void clearContinuousbackupsdescription() => $_clearField(446030650);
  @$pb.TagNumber(446030650)
  ContinuousBackupsDescription ensureContinuousbackupsdescription() =>
      $_ensure(0);
}

class DescribeContributorInsightsInput extends $pb.GeneratedMessage {
  factory DescribeContributorInsightsInput({
    $core.String? indexname,
    $core.String? tablename,
  }) {
    final result = create();
    if (indexname != null) result.indexname = indexname;
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DescribeContributorInsightsInput._();

  factory DescribeContributorInsightsInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeContributorInsightsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeContributorInsightsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeContributorInsightsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeContributorInsightsInput copyWith(
          void Function(DescribeContributorInsightsInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeContributorInsightsInput))
          as DescribeContributorInsightsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeContributorInsightsInput create() =>
      DescribeContributorInsightsInput._();
  @$core.override
  DescribeContributorInsightsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeContributorInsightsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeContributorInsightsInput>(
          create);
  static DescribeContributorInsightsInput? _defaultInstance;

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(0);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(0);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeContributorInsightsOutput extends $pb.GeneratedMessage {
  factory DescribeContributorInsightsOutput({
    ContributorInsightsMode? contributorinsightsmode,
    $core.String? indexname,
    $core.Iterable<$core.String>? contributorinsightsrulelist,
    $core.String? tablename,
    FailureException? failureexception,
    ContributorInsightsStatus? contributorinsightsstatus,
    $core.String? lastupdatedatetime,
  }) {
    final result = create();
    if (contributorinsightsmode != null)
      result.contributorinsightsmode = contributorinsightsmode;
    if (indexname != null) result.indexname = indexname;
    if (contributorinsightsrulelist != null)
      result.contributorinsightsrulelist.addAll(contributorinsightsrulelist);
    if (tablename != null) result.tablename = tablename;
    if (failureexception != null) result.failureexception = failureexception;
    if (contributorinsightsstatus != null)
      result.contributorinsightsstatus = contributorinsightsstatus;
    if (lastupdatedatetime != null)
      result.lastupdatedatetime = lastupdatedatetime;
    return result;
  }

  DescribeContributorInsightsOutput._();

  factory DescribeContributorInsightsOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeContributorInsightsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeContributorInsightsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ContributorInsightsMode>(
        86700161, _omitFieldNames ? '' : 'contributorinsightsmode',
        enumValues: ContributorInsightsMode.values)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..pPS(140475206, _omitFieldNames ? '' : 'contributorinsightsrulelist')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<FailureException>(
        284668899, _omitFieldNames ? '' : 'failureexception',
        subBuilder: FailureException.create)
    ..aE<ContributorInsightsStatus>(
        363347282, _omitFieldNames ? '' : 'contributorinsightsstatus',
        enumValues: ContributorInsightsStatus.values)
    ..aOS(452274318, _omitFieldNames ? '' : 'lastupdatedatetime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeContributorInsightsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeContributorInsightsOutput copyWith(
          void Function(DescribeContributorInsightsOutput) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeContributorInsightsOutput))
          as DescribeContributorInsightsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeContributorInsightsOutput create() =>
      DescribeContributorInsightsOutput._();
  @$core.override
  DescribeContributorInsightsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeContributorInsightsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeContributorInsightsOutput>(
          create);
  static DescribeContributorInsightsOutput? _defaultInstance;

  @$pb.TagNumber(86700161)
  ContributorInsightsMode get contributorinsightsmode => $_getN(0);
  @$pb.TagNumber(86700161)
  set contributorinsightsmode(ContributorInsightsMode value) =>
      $_setField(86700161, value);
  @$pb.TagNumber(86700161)
  $core.bool hasContributorinsightsmode() => $_has(0);
  @$pb.TagNumber(86700161)
  void clearContributorinsightsmode() => $_clearField(86700161);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(140475206)
  $pb.PbList<$core.String> get contributorinsightsrulelist => $_getList(2);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(284668899)
  FailureException get failureexception => $_getN(4);
  @$pb.TagNumber(284668899)
  set failureexception(FailureException value) => $_setField(284668899, value);
  @$pb.TagNumber(284668899)
  $core.bool hasFailureexception() => $_has(4);
  @$pb.TagNumber(284668899)
  void clearFailureexception() => $_clearField(284668899);
  @$pb.TagNumber(284668899)
  FailureException ensureFailureexception() => $_ensure(4);

  @$pb.TagNumber(363347282)
  ContributorInsightsStatus get contributorinsightsstatus => $_getN(5);
  @$pb.TagNumber(363347282)
  set contributorinsightsstatus(ContributorInsightsStatus value) =>
      $_setField(363347282, value);
  @$pb.TagNumber(363347282)
  $core.bool hasContributorinsightsstatus() => $_has(5);
  @$pb.TagNumber(363347282)
  void clearContributorinsightsstatus() => $_clearField(363347282);

  @$pb.TagNumber(452274318)
  $core.String get lastupdatedatetime => $_getSZ(6);
  @$pb.TagNumber(452274318)
  set lastupdatedatetime($core.String value) => $_setString(6, value);
  @$pb.TagNumber(452274318)
  $core.bool hasLastupdatedatetime() => $_has(6);
  @$pb.TagNumber(452274318)
  void clearLastupdatedatetime() => $_clearField(452274318);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
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

class DescribeExportInput extends $pb.GeneratedMessage {
  factory DescribeExportInput({
    $core.String? exportarn,
  }) {
    final result = create();
    if (exportarn != null) result.exportarn = exportarn;
    return result;
  }

  DescribeExportInput._();

  factory DescribeExportInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeExportInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeExportInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(3661287, _omitFieldNames ? '' : 'exportarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeExportInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeExportInput copyWith(void Function(DescribeExportInput) updates) =>
      super.copyWith((message) => updates(message as DescribeExportInput))
          as DescribeExportInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeExportInput create() => DescribeExportInput._();
  @$core.override
  DescribeExportInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeExportInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeExportInput>(create);
  static DescribeExportInput? _defaultInstance;

  @$pb.TagNumber(3661287)
  $core.String get exportarn => $_getSZ(0);
  @$pb.TagNumber(3661287)
  set exportarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(3661287)
  $core.bool hasExportarn() => $_has(0);
  @$pb.TagNumber(3661287)
  void clearExportarn() => $_clearField(3661287);
}

class DescribeExportOutput extends $pb.GeneratedMessage {
  factory DescribeExportOutput({
    ExportDescription? exportdescription,
  }) {
    final result = create();
    if (exportdescription != null) result.exportdescription = exportdescription;
    return result;
  }

  DescribeExportOutput._();

  factory DescribeExportOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeExportOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeExportOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ExportDescription>(
        14399944, _omitFieldNames ? '' : 'exportdescription',
        subBuilder: ExportDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeExportOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeExportOutput copyWith(void Function(DescribeExportOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeExportOutput))
          as DescribeExportOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeExportOutput create() => DescribeExportOutput._();
  @$core.override
  DescribeExportOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeExportOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeExportOutput>(create);
  static DescribeExportOutput? _defaultInstance;

  @$pb.TagNumber(14399944)
  ExportDescription get exportdescription => $_getN(0);
  @$pb.TagNumber(14399944)
  set exportdescription(ExportDescription value) => $_setField(14399944, value);
  @$pb.TagNumber(14399944)
  $core.bool hasExportdescription() => $_has(0);
  @$pb.TagNumber(14399944)
  void clearExportdescription() => $_clearField(14399944);
  @$pb.TagNumber(14399944)
  ExportDescription ensureExportdescription() => $_ensure(0);
}

class DescribeGlobalTableInput extends $pb.GeneratedMessage {
  factory DescribeGlobalTableInput({
    $core.String? globaltablename,
  }) {
    final result = create();
    if (globaltablename != null) result.globaltablename = globaltablename;
    return result;
  }

  DescribeGlobalTableInput._();

  factory DescribeGlobalTableInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeGlobalTableInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeGlobalTableInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeGlobalTableInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeGlobalTableInput copyWith(
          void Function(DescribeGlobalTableInput) updates) =>
      super.copyWith((message) => updates(message as DescribeGlobalTableInput))
          as DescribeGlobalTableInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeGlobalTableInput create() => DescribeGlobalTableInput._();
  @$core.override
  DescribeGlobalTableInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeGlobalTableInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeGlobalTableInput>(create);
  static DescribeGlobalTableInput? _defaultInstance;

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(0);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(0);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);
}

class DescribeGlobalTableOutput extends $pb.GeneratedMessage {
  factory DescribeGlobalTableOutput({
    GlobalTableDescription? globaltabledescription,
  }) {
    final result = create();
    if (globaltabledescription != null)
      result.globaltabledescription = globaltabledescription;
    return result;
  }

  DescribeGlobalTableOutput._();

  factory DescribeGlobalTableOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeGlobalTableOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeGlobalTableOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<GlobalTableDescription>(
        342116405, _omitFieldNames ? '' : 'globaltabledescription',
        subBuilder: GlobalTableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeGlobalTableOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeGlobalTableOutput copyWith(
          void Function(DescribeGlobalTableOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeGlobalTableOutput))
          as DescribeGlobalTableOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeGlobalTableOutput create() => DescribeGlobalTableOutput._();
  @$core.override
  DescribeGlobalTableOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeGlobalTableOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeGlobalTableOutput>(create);
  static DescribeGlobalTableOutput? _defaultInstance;

  @$pb.TagNumber(342116405)
  GlobalTableDescription get globaltabledescription => $_getN(0);
  @$pb.TagNumber(342116405)
  set globaltabledescription(GlobalTableDescription value) =>
      $_setField(342116405, value);
  @$pb.TagNumber(342116405)
  $core.bool hasGlobaltabledescription() => $_has(0);
  @$pb.TagNumber(342116405)
  void clearGlobaltabledescription() => $_clearField(342116405);
  @$pb.TagNumber(342116405)
  GlobalTableDescription ensureGlobaltabledescription() => $_ensure(0);
}

class DescribeGlobalTableSettingsInput extends $pb.GeneratedMessage {
  factory DescribeGlobalTableSettingsInput({
    $core.String? globaltablename,
  }) {
    final result = create();
    if (globaltablename != null) result.globaltablename = globaltablename;
    return result;
  }

  DescribeGlobalTableSettingsInput._();

  factory DescribeGlobalTableSettingsInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeGlobalTableSettingsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeGlobalTableSettingsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeGlobalTableSettingsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeGlobalTableSettingsInput copyWith(
          void Function(DescribeGlobalTableSettingsInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeGlobalTableSettingsInput))
          as DescribeGlobalTableSettingsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeGlobalTableSettingsInput create() =>
      DescribeGlobalTableSettingsInput._();
  @$core.override
  DescribeGlobalTableSettingsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeGlobalTableSettingsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeGlobalTableSettingsInput>(
          create);
  static DescribeGlobalTableSettingsInput? _defaultInstance;

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(0);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(0);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);
}

class DescribeGlobalTableSettingsOutput extends $pb.GeneratedMessage {
  factory DescribeGlobalTableSettingsOutput({
    $core.String? globaltablename,
    $core.Iterable<ReplicaSettingsDescription>? replicasettings,
  }) {
    final result = create();
    if (globaltablename != null) result.globaltablename = globaltablename;
    if (replicasettings != null) result.replicasettings.addAll(replicasettings);
    return result;
  }

  DescribeGlobalTableSettingsOutput._();

  factory DescribeGlobalTableSettingsOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeGlobalTableSettingsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeGlobalTableSettingsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..pPM<ReplicaSettingsDescription>(
        288882107, _omitFieldNames ? '' : 'replicasettings',
        subBuilder: ReplicaSettingsDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeGlobalTableSettingsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeGlobalTableSettingsOutput copyWith(
          void Function(DescribeGlobalTableSettingsOutput) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeGlobalTableSettingsOutput))
          as DescribeGlobalTableSettingsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeGlobalTableSettingsOutput create() =>
      DescribeGlobalTableSettingsOutput._();
  @$core.override
  DescribeGlobalTableSettingsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeGlobalTableSettingsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeGlobalTableSettingsOutput>(
          create);
  static DescribeGlobalTableSettingsOutput? _defaultInstance;

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(0);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(0);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);

  @$pb.TagNumber(288882107)
  $pb.PbList<ReplicaSettingsDescription> get replicasettings => $_getList(1);
}

class DescribeImportInput extends $pb.GeneratedMessage {
  factory DescribeImportInput({
    $core.String? importarn,
  }) {
    final result = create();
    if (importarn != null) result.importarn = importarn;
    return result;
  }

  DescribeImportInput._();

  factory DescribeImportInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeImportInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeImportInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(444379628, _omitFieldNames ? '' : 'importarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeImportInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeImportInput copyWith(void Function(DescribeImportInput) updates) =>
      super.copyWith((message) => updates(message as DescribeImportInput))
          as DescribeImportInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeImportInput create() => DescribeImportInput._();
  @$core.override
  DescribeImportInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeImportInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeImportInput>(create);
  static DescribeImportInput? _defaultInstance;

  @$pb.TagNumber(444379628)
  $core.String get importarn => $_getSZ(0);
  @$pb.TagNumber(444379628)
  set importarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(444379628)
  $core.bool hasImportarn() => $_has(0);
  @$pb.TagNumber(444379628)
  void clearImportarn() => $_clearField(444379628);
}

class DescribeImportOutput extends $pb.GeneratedMessage {
  factory DescribeImportOutput({
    ImportTableDescription? importtabledescription,
  }) {
    final result = create();
    if (importtabledescription != null)
      result.importtabledescription = importtabledescription;
    return result;
  }

  DescribeImportOutput._();

  factory DescribeImportOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeImportOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeImportOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ImportTableDescription>(
        407746467, _omitFieldNames ? '' : 'importtabledescription',
        subBuilder: ImportTableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeImportOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeImportOutput copyWith(void Function(DescribeImportOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeImportOutput))
          as DescribeImportOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeImportOutput create() => DescribeImportOutput._();
  @$core.override
  DescribeImportOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeImportOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeImportOutput>(create);
  static DescribeImportOutput? _defaultInstance;

  @$pb.TagNumber(407746467)
  ImportTableDescription get importtabledescription => $_getN(0);
  @$pb.TagNumber(407746467)
  set importtabledescription(ImportTableDescription value) =>
      $_setField(407746467, value);
  @$pb.TagNumber(407746467)
  $core.bool hasImporttabledescription() => $_has(0);
  @$pb.TagNumber(407746467)
  void clearImporttabledescription() => $_clearField(407746467);
  @$pb.TagNumber(407746467)
  ImportTableDescription ensureImporttabledescription() => $_ensure(0);
}

class DescribeKinesisStreamingDestinationInput extends $pb.GeneratedMessage {
  factory DescribeKinesisStreamingDestinationInput({
    $core.String? tablename,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DescribeKinesisStreamingDestinationInput._();

  factory DescribeKinesisStreamingDestinationInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeKinesisStreamingDestinationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeKinesisStreamingDestinationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeKinesisStreamingDestinationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeKinesisStreamingDestinationInput copyWith(
          void Function(DescribeKinesisStreamingDestinationInput) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeKinesisStreamingDestinationInput))
          as DescribeKinesisStreamingDestinationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeKinesisStreamingDestinationInput create() =>
      DescribeKinesisStreamingDestinationInput._();
  @$core.override
  DescribeKinesisStreamingDestinationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeKinesisStreamingDestinationInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DescribeKinesisStreamingDestinationInput>(create);
  static DescribeKinesisStreamingDestinationInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeKinesisStreamingDestinationOutput extends $pb.GeneratedMessage {
  factory DescribeKinesisStreamingDestinationOutput({
    $core.Iterable<KinesisDataStreamDestination>? kinesisdatastreamdestinations,
    $core.String? tablename,
  }) {
    final result = create();
    if (kinesisdatastreamdestinations != null)
      result.kinesisdatastreamdestinations
          .addAll(kinesisdatastreamdestinations);
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DescribeKinesisStreamingDestinationOutput._();

  factory DescribeKinesisStreamingDestinationOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeKinesisStreamingDestinationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeKinesisStreamingDestinationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<KinesisDataStreamDestination>(
        4645129, _omitFieldNames ? '' : 'kinesisdatastreamdestinations',
        subBuilder: KinesisDataStreamDestination.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeKinesisStreamingDestinationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeKinesisStreamingDestinationOutput copyWith(
          void Function(DescribeKinesisStreamingDestinationOutput) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeKinesisStreamingDestinationOutput))
          as DescribeKinesisStreamingDestinationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeKinesisStreamingDestinationOutput create() =>
      DescribeKinesisStreamingDestinationOutput._();
  @$core.override
  DescribeKinesisStreamingDestinationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeKinesisStreamingDestinationOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DescribeKinesisStreamingDestinationOutput>(create);
  static DescribeKinesisStreamingDestinationOutput? _defaultInstance;

  @$pb.TagNumber(4645129)
  $pb.PbList<KinesisDataStreamDestination> get kinesisdatastreamdestinations =>
      $_getList(0);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeLimitsInput extends $pb.GeneratedMessage {
  factory DescribeLimitsInput() => create();

  DescribeLimitsInput._();

  factory DescribeLimitsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeLimitsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeLimitsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeLimitsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeLimitsInput copyWith(void Function(DescribeLimitsInput) updates) =>
      super.copyWith((message) => updates(message as DescribeLimitsInput))
          as DescribeLimitsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeLimitsInput create() => DescribeLimitsInput._();
  @$core.override
  DescribeLimitsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeLimitsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeLimitsInput>(create);
  static DescribeLimitsInput? _defaultInstance;
}

class DescribeLimitsOutput extends $pb.GeneratedMessage {
  factory DescribeLimitsOutput({
    $fixnum.Int64? tablemaxreadcapacityunits,
    $fixnum.Int64? accountmaxreadcapacityunits,
    $fixnum.Int64? tablemaxwritecapacityunits,
    $fixnum.Int64? accountmaxwritecapacityunits,
  }) {
    final result = create();
    if (tablemaxreadcapacityunits != null)
      result.tablemaxreadcapacityunits = tablemaxreadcapacityunits;
    if (accountmaxreadcapacityunits != null)
      result.accountmaxreadcapacityunits = accountmaxreadcapacityunits;
    if (tablemaxwritecapacityunits != null)
      result.tablemaxwritecapacityunits = tablemaxwritecapacityunits;
    if (accountmaxwritecapacityunits != null)
      result.accountmaxwritecapacityunits = accountmaxwritecapacityunits;
    return result;
  }

  DescribeLimitsOutput._();

  factory DescribeLimitsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeLimitsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeLimitsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(165993307, _omitFieldNames ? '' : 'tablemaxreadcapacityunits')
    ..aInt64(318006242, _omitFieldNames ? '' : 'accountmaxreadcapacityunits')
    ..aInt64(444127818, _omitFieldNames ? '' : 'tablemaxwritecapacityunits')
    ..aInt64(526094897, _omitFieldNames ? '' : 'accountmaxwritecapacityunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeLimitsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeLimitsOutput copyWith(void Function(DescribeLimitsOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeLimitsOutput))
          as DescribeLimitsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeLimitsOutput create() => DescribeLimitsOutput._();
  @$core.override
  DescribeLimitsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeLimitsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeLimitsOutput>(create);
  static DescribeLimitsOutput? _defaultInstance;

  @$pb.TagNumber(165993307)
  $fixnum.Int64 get tablemaxreadcapacityunits => $_getI64(0);
  @$pb.TagNumber(165993307)
  set tablemaxreadcapacityunits($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(165993307)
  $core.bool hasTablemaxreadcapacityunits() => $_has(0);
  @$pb.TagNumber(165993307)
  void clearTablemaxreadcapacityunits() => $_clearField(165993307);

  @$pb.TagNumber(318006242)
  $fixnum.Int64 get accountmaxreadcapacityunits => $_getI64(1);
  @$pb.TagNumber(318006242)
  set accountmaxreadcapacityunits($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(318006242)
  $core.bool hasAccountmaxreadcapacityunits() => $_has(1);
  @$pb.TagNumber(318006242)
  void clearAccountmaxreadcapacityunits() => $_clearField(318006242);

  @$pb.TagNumber(444127818)
  $fixnum.Int64 get tablemaxwritecapacityunits => $_getI64(2);
  @$pb.TagNumber(444127818)
  set tablemaxwritecapacityunits($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(444127818)
  $core.bool hasTablemaxwritecapacityunits() => $_has(2);
  @$pb.TagNumber(444127818)
  void clearTablemaxwritecapacityunits() => $_clearField(444127818);

  @$pb.TagNumber(526094897)
  $fixnum.Int64 get accountmaxwritecapacityunits => $_getI64(3);
  @$pb.TagNumber(526094897)
  set accountmaxwritecapacityunits($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(526094897)
  $core.bool hasAccountmaxwritecapacityunits() => $_has(3);
  @$pb.TagNumber(526094897)
  void clearAccountmaxwritecapacityunits() => $_clearField(526094897);
}

class DescribeTableInput extends $pb.GeneratedMessage {
  factory DescribeTableInput({
    $core.String? tablename,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DescribeTableInput._();

  factory DescribeTableInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTableInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTableInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableInput copyWith(void Function(DescribeTableInput) updates) =>
      super.copyWith((message) => updates(message as DescribeTableInput))
          as DescribeTableInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTableInput create() => DescribeTableInput._();
  @$core.override
  DescribeTableInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTableInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeTableInput>(create);
  static DescribeTableInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeTableOutput extends $pb.GeneratedMessage {
  factory DescribeTableOutput({
    TableDescription? table,
  }) {
    final result = create();
    if (table != null) result.table = table;
    return result;
  }

  DescribeTableOutput._();

  factory DescribeTableOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTableOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTableOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TableDescription>(386722688, _omitFieldNames ? '' : 'table',
        subBuilder: TableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableOutput copyWith(void Function(DescribeTableOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeTableOutput))
          as DescribeTableOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTableOutput create() => DescribeTableOutput._();
  @$core.override
  DescribeTableOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTableOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeTableOutput>(create);
  static DescribeTableOutput? _defaultInstance;

  @$pb.TagNumber(386722688)
  TableDescription get table => $_getN(0);
  @$pb.TagNumber(386722688)
  set table(TableDescription value) => $_setField(386722688, value);
  @$pb.TagNumber(386722688)
  $core.bool hasTable() => $_has(0);
  @$pb.TagNumber(386722688)
  void clearTable() => $_clearField(386722688);
  @$pb.TagNumber(386722688)
  TableDescription ensureTable() => $_ensure(0);
}

class DescribeTableReplicaAutoScalingInput extends $pb.GeneratedMessage {
  factory DescribeTableReplicaAutoScalingInput({
    $core.String? tablename,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DescribeTableReplicaAutoScalingInput._();

  factory DescribeTableReplicaAutoScalingInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTableReplicaAutoScalingInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTableReplicaAutoScalingInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableReplicaAutoScalingInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableReplicaAutoScalingInput copyWith(
          void Function(DescribeTableReplicaAutoScalingInput) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeTableReplicaAutoScalingInput))
          as DescribeTableReplicaAutoScalingInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTableReplicaAutoScalingInput create() =>
      DescribeTableReplicaAutoScalingInput._();
  @$core.override
  DescribeTableReplicaAutoScalingInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTableReplicaAutoScalingInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DescribeTableReplicaAutoScalingInput>(create);
  static DescribeTableReplicaAutoScalingInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeTableReplicaAutoScalingOutput extends $pb.GeneratedMessage {
  factory DescribeTableReplicaAutoScalingOutput({
    TableAutoScalingDescription? tableautoscalingdescription,
  }) {
    final result = create();
    if (tableautoscalingdescription != null)
      result.tableautoscalingdescription = tableautoscalingdescription;
    return result;
  }

  DescribeTableReplicaAutoScalingOutput._();

  factory DescribeTableReplicaAutoScalingOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTableReplicaAutoScalingOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTableReplicaAutoScalingOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TableAutoScalingDescription>(
        490422806, _omitFieldNames ? '' : 'tableautoscalingdescription',
        subBuilder: TableAutoScalingDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableReplicaAutoScalingOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTableReplicaAutoScalingOutput copyWith(
          void Function(DescribeTableReplicaAutoScalingOutput) updates) =>
      super.copyWith((message) =>
              updates(message as DescribeTableReplicaAutoScalingOutput))
          as DescribeTableReplicaAutoScalingOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTableReplicaAutoScalingOutput create() =>
      DescribeTableReplicaAutoScalingOutput._();
  @$core.override
  DescribeTableReplicaAutoScalingOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTableReplicaAutoScalingOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DescribeTableReplicaAutoScalingOutput>(create);
  static DescribeTableReplicaAutoScalingOutput? _defaultInstance;

  @$pb.TagNumber(490422806)
  TableAutoScalingDescription get tableautoscalingdescription => $_getN(0);
  @$pb.TagNumber(490422806)
  set tableautoscalingdescription(TableAutoScalingDescription value) =>
      $_setField(490422806, value);
  @$pb.TagNumber(490422806)
  $core.bool hasTableautoscalingdescription() => $_has(0);
  @$pb.TagNumber(490422806)
  void clearTableautoscalingdescription() => $_clearField(490422806);
  @$pb.TagNumber(490422806)
  TableAutoScalingDescription ensureTableautoscalingdescription() =>
      $_ensure(0);
}

class DescribeTimeToLiveInput extends $pb.GeneratedMessage {
  factory DescribeTimeToLiveInput({
    $core.String? tablename,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  DescribeTimeToLiveInput._();

  factory DescribeTimeToLiveInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTimeToLiveInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTimeToLiveInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTimeToLiveInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTimeToLiveInput copyWith(
          void Function(DescribeTimeToLiveInput) updates) =>
      super.copyWith((message) => updates(message as DescribeTimeToLiveInput))
          as DescribeTimeToLiveInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTimeToLiveInput create() => DescribeTimeToLiveInput._();
  @$core.override
  DescribeTimeToLiveInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTimeToLiveInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeTimeToLiveInput>(create);
  static DescribeTimeToLiveInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class DescribeTimeToLiveOutput extends $pb.GeneratedMessage {
  factory DescribeTimeToLiveOutput({
    TimeToLiveDescription? timetolivedescription,
  }) {
    final result = create();
    if (timetolivedescription != null)
      result.timetolivedescription = timetolivedescription;
    return result;
  }

  DescribeTimeToLiveOutput._();

  factory DescribeTimeToLiveOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeTimeToLiveOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeTimeToLiveOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TimeToLiveDescription>(
        367590876, _omitFieldNames ? '' : 'timetolivedescription',
        subBuilder: TimeToLiveDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTimeToLiveOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeTimeToLiveOutput copyWith(
          void Function(DescribeTimeToLiveOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeTimeToLiveOutput))
          as DescribeTimeToLiveOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeTimeToLiveOutput create() => DescribeTimeToLiveOutput._();
  @$core.override
  DescribeTimeToLiveOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeTimeToLiveOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeTimeToLiveOutput>(create);
  static DescribeTimeToLiveOutput? _defaultInstance;

  @$pb.TagNumber(367590876)
  TimeToLiveDescription get timetolivedescription => $_getN(0);
  @$pb.TagNumber(367590876)
  set timetolivedescription(TimeToLiveDescription value) =>
      $_setField(367590876, value);
  @$pb.TagNumber(367590876)
  $core.bool hasTimetolivedescription() => $_has(0);
  @$pb.TagNumber(367590876)
  void clearTimetolivedescription() => $_clearField(367590876);
  @$pb.TagNumber(367590876)
  TimeToLiveDescription ensureTimetolivedescription() => $_ensure(0);
}

class DuplicateItemException extends $pb.GeneratedMessage {
  factory DuplicateItemException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DuplicateItemException._();

  factory DuplicateItemException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DuplicateItemException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DuplicateItemException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DuplicateItemException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DuplicateItemException copyWith(
          void Function(DuplicateItemException) updates) =>
      super.copyWith((message) => updates(message as DuplicateItemException))
          as DuplicateItemException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DuplicateItemException create() => DuplicateItemException._();
  @$core.override
  DuplicateItemException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DuplicateItemException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DuplicateItemException>(create);
  static DuplicateItemException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class EnableKinesisStreamingConfiguration extends $pb.GeneratedMessage {
  factory EnableKinesisStreamingConfiguration({
    ApproximateCreationDateTimePrecision? approximatecreationdatetimeprecision,
  }) {
    final result = create();
    if (approximatecreationdatetimeprecision != null)
      result.approximatecreationdatetimeprecision =
          approximatecreationdatetimeprecision;
    return result;
  }

  EnableKinesisStreamingConfiguration._();

  factory EnableKinesisStreamingConfiguration.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableKinesisStreamingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableKinesisStreamingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ApproximateCreationDateTimePrecision>(392293352,
        _omitFieldNames ? '' : 'approximatecreationdatetimeprecision',
        enumValues: ApproximateCreationDateTimePrecision.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableKinesisStreamingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableKinesisStreamingConfiguration copyWith(
          void Function(EnableKinesisStreamingConfiguration) updates) =>
      super.copyWith((message) =>
              updates(message as EnableKinesisStreamingConfiguration))
          as EnableKinesisStreamingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableKinesisStreamingConfiguration create() =>
      EnableKinesisStreamingConfiguration._();
  @$core.override
  EnableKinesisStreamingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableKinesisStreamingConfiguration getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          EnableKinesisStreamingConfiguration>(create);
  static EnableKinesisStreamingConfiguration? _defaultInstance;

  @$pb.TagNumber(392293352)
  ApproximateCreationDateTimePrecision
      get approximatecreationdatetimeprecision => $_getN(0);
  @$pb.TagNumber(392293352)
  set approximatecreationdatetimeprecision(
          ApproximateCreationDateTimePrecision value) =>
      $_setField(392293352, value);
  @$pb.TagNumber(392293352)
  $core.bool hasApproximatecreationdatetimeprecision() => $_has(0);
  @$pb.TagNumber(392293352)
  void clearApproximatecreationdatetimeprecision() => $_clearField(392293352);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
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

class ExecuteStatementInput extends $pb.GeneratedMessage {
  factory ExecuteStatementInput({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.String? nexttoken,
    $core.String? statement,
    $core.int? limit,
    $core.Iterable<AttributeValue>? parameters,
    $core.bool? consistentread,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (statement != null) result.statement = statement;
    if (limit != null) result.limit = limit;
    if (parameters != null) result.parameters.addAll(parameters);
    if (consistentread != null) result.consistentread = consistentread;
    return result;
  }

  ExecuteStatementInput._();

  factory ExecuteStatementInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecuteStatementInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecuteStatementInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(248790199, _omitFieldNames ? '' : 'statement')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..pPM<AttributeValue>(494900218, _omitFieldNames ? '' : 'parameters',
        subBuilder: AttributeValue.create)
    ..aOB(531556994, _omitFieldNames ? '' : 'consistentread')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteStatementInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteStatementInput copyWith(
          void Function(ExecuteStatementInput) updates) =>
      super.copyWith((message) => updates(message as ExecuteStatementInput))
          as ExecuteStatementInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecuteStatementInput create() => ExecuteStatementInput._();
  @$core.override
  ExecuteStatementInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecuteStatementInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecuteStatementInput>(create);
  static ExecuteStatementInput? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(1);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(1);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(248790199)
  $core.String get statement => $_getSZ(3);
  @$pb.TagNumber(248790199)
  set statement($core.String value) => $_setString(3, value);
  @$pb.TagNumber(248790199)
  $core.bool hasStatement() => $_has(3);
  @$pb.TagNumber(248790199)
  void clearStatement() => $_clearField(248790199);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(4);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(4);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(494900218)
  $pb.PbList<AttributeValue> get parameters => $_getList(5);

  @$pb.TagNumber(531556994)
  $core.bool get consistentread => $_getBF(6);
  @$pb.TagNumber(531556994)
  set consistentread($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(531556994)
  $core.bool hasConsistentread() => $_has(6);
  @$pb.TagNumber(531556994)
  void clearConsistentread() => $_clearField(531556994);
}

class ExecuteStatementOutput extends $pb.GeneratedMessage {
  factory ExecuteStatementOutput({
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? items,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        lastevaluatedkey,
    $core.String? nexttoken,
    ConsumedCapacity? consumedcapacity,
  }) {
    final result = create();
    if (items != null) result.items.addEntries(items);
    if (lastevaluatedkey != null)
      result.lastevaluatedkey.addEntries(lastevaluatedkey);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (consumedcapacity != null) result.consumedcapacity = consumedcapacity;
    return result;
  }

  ExecuteStatementOutput._();

  factory ExecuteStatementOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecuteStatementOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecuteStatementOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, AttributeValue>(3553328, _omitFieldNames ? '' : 'items',
        entryClassName: 'ExecuteStatementOutput.ItemsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, AttributeValue>(
        54319830, _omitFieldNames ? '' : 'lastevaluatedkey',
        entryClassName: 'ExecuteStatementOutput.LastevaluatedkeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteStatementOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteStatementOutput copyWith(
          void Function(ExecuteStatementOutput) updates) =>
      super.copyWith((message) => updates(message as ExecuteStatementOutput))
          as ExecuteStatementOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecuteStatementOutput create() => ExecuteStatementOutput._();
  @$core.override
  ExecuteStatementOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecuteStatementOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecuteStatementOutput>(create);
  static ExecuteStatementOutput? _defaultInstance;

  @$pb.TagNumber(3553328)
  $pb.PbMap<$core.String, AttributeValue> get items => $_getMap(0);

  @$pb.TagNumber(54319830)
  $pb.PbMap<$core.String, AttributeValue> get lastevaluatedkey => $_getMap(1);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(449336620)
  ConsumedCapacity get consumedcapacity => $_getN(3);
  @$pb.TagNumber(449336620)
  set consumedcapacity(ConsumedCapacity value) => $_setField(449336620, value);
  @$pb.TagNumber(449336620)
  $core.bool hasConsumedcapacity() => $_has(3);
  @$pb.TagNumber(449336620)
  void clearConsumedcapacity() => $_clearField(449336620);
  @$pb.TagNumber(449336620)
  ConsumedCapacity ensureConsumedcapacity() => $_ensure(3);
}

class ExecuteTransactionInput extends $pb.GeneratedMessage {
  factory ExecuteTransactionInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<ParameterizedStatement>? transactstatements,
    $core.String? clientrequesttoken,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (transactstatements != null)
      result.transactstatements.addAll(transactstatements);
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    return result;
  }

  ExecuteTransactionInput._();

  factory ExecuteTransactionInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecuteTransactionInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecuteTransactionInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..pPM<ParameterizedStatement>(
        446038782, _omitFieldNames ? '' : 'transactstatements',
        subBuilder: ParameterizedStatement.create)
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteTransactionInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteTransactionInput copyWith(
          void Function(ExecuteTransactionInput) updates) =>
      super.copyWith((message) => updates(message as ExecuteTransactionInput))
          as ExecuteTransactionInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecuteTransactionInput create() => ExecuteTransactionInput._();
  @$core.override
  ExecuteTransactionInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecuteTransactionInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecuteTransactionInput>(create);
  static ExecuteTransactionInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(446038782)
  $pb.PbList<ParameterizedStatement> get transactstatements => $_getList(1);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(2);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(2);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);
}

class ExecuteTransactionOutput extends $pb.GeneratedMessage {
  factory ExecuteTransactionOutput({
    $core.Iterable<ItemResponse>? responses,
    $core.Iterable<ConsumedCapacity>? consumedcapacity,
  }) {
    final result = create();
    if (responses != null) result.responses.addAll(responses);
    if (consumedcapacity != null)
      result.consumedcapacity.addAll(consumedcapacity);
    return result;
  }

  ExecuteTransactionOutput._();

  factory ExecuteTransactionOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExecuteTransactionOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExecuteTransactionOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<ItemResponse>(114852856, _omitFieldNames ? '' : 'responses',
        subBuilder: ItemResponse.create)
    ..pPM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteTransactionOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExecuteTransactionOutput copyWith(
          void Function(ExecuteTransactionOutput) updates) =>
      super.copyWith((message) => updates(message as ExecuteTransactionOutput))
          as ExecuteTransactionOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExecuteTransactionOutput create() => ExecuteTransactionOutput._();
  @$core.override
  ExecuteTransactionOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExecuteTransactionOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExecuteTransactionOutput>(create);
  static ExecuteTransactionOutput? _defaultInstance;

  @$pb.TagNumber(114852856)
  $pb.PbList<ItemResponse> get responses => $_getList(0);

  @$pb.TagNumber(449336620)
  $pb.PbList<ConsumedCapacity> get consumedcapacity => $_getList(1);
}

class ExpectedAttributeValue extends $pb.GeneratedMessage {
  factory ExpectedAttributeValue({
    ComparisonOperator? comparisonoperator,
    $core.Iterable<AttributeValue>? attributevaluelist,
    $core.bool? exists,
    AttributeValue? value,
  }) {
    final result = create();
    if (comparisonoperator != null)
      result.comparisonoperator = comparisonoperator;
    if (attributevaluelist != null)
      result.attributevaluelist.addAll(attributevaluelist);
    if (exists != null) result.exists = exists;
    if (value != null) result.value = value;
    return result;
  }

  ExpectedAttributeValue._();

  factory ExpectedAttributeValue.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExpectedAttributeValue.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExpectedAttributeValue',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ComparisonOperator>(
        7059861, _omitFieldNames ? '' : 'comparisonoperator',
        enumValues: ComparisonOperator.values)
    ..pPM<AttributeValue>(
        157424013, _omitFieldNames ? '' : 'attributevaluelist',
        subBuilder: AttributeValue.create)
    ..aOB(265084382, _omitFieldNames ? '' : 'exists')
    ..aOM<AttributeValue>(289929579, _omitFieldNames ? '' : 'value',
        subBuilder: AttributeValue.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpectedAttributeValue clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExpectedAttributeValue copyWith(
          void Function(ExpectedAttributeValue) updates) =>
      super.copyWith((message) => updates(message as ExpectedAttributeValue))
          as ExpectedAttributeValue;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExpectedAttributeValue create() => ExpectedAttributeValue._();
  @$core.override
  ExpectedAttributeValue createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExpectedAttributeValue getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExpectedAttributeValue>(create);
  static ExpectedAttributeValue? _defaultInstance;

  @$pb.TagNumber(7059861)
  ComparisonOperator get comparisonoperator => $_getN(0);
  @$pb.TagNumber(7059861)
  set comparisonoperator(ComparisonOperator value) =>
      $_setField(7059861, value);
  @$pb.TagNumber(7059861)
  $core.bool hasComparisonoperator() => $_has(0);
  @$pb.TagNumber(7059861)
  void clearComparisonoperator() => $_clearField(7059861);

  @$pb.TagNumber(157424013)
  $pb.PbList<AttributeValue> get attributevaluelist => $_getList(1);

  @$pb.TagNumber(265084382)
  $core.bool get exists => $_getBF(2);
  @$pb.TagNumber(265084382)
  set exists($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(265084382)
  $core.bool hasExists() => $_has(2);
  @$pb.TagNumber(265084382)
  void clearExists() => $_clearField(265084382);

  @$pb.TagNumber(289929579)
  AttributeValue get value => $_getN(3);
  @$pb.TagNumber(289929579)
  set value(AttributeValue value) => $_setField(289929579, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(3);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
  @$pb.TagNumber(289929579)
  AttributeValue ensureValue() => $_ensure(3);
}

class ExportConflictException extends $pb.GeneratedMessage {
  factory ExportConflictException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExportConflictException._();

  factory ExportConflictException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportConflictException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportConflictException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportConflictException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportConflictException copyWith(
          void Function(ExportConflictException) updates) =>
      super.copyWith((message) => updates(message as ExportConflictException))
          as ExportConflictException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportConflictException create() => ExportConflictException._();
  @$core.override
  ExportConflictException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportConflictException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportConflictException>(create);
  static ExportConflictException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExportDescription extends $pb.GeneratedMessage {
  factory ExportDescription({
    $core.String? exportarn,
    $core.String? s3ssekmskeyid,
    IncrementalExportSpecification? incrementalexportspecification,
    $core.String? s3prefix,
    $fixnum.Int64? itemcount,
    $core.String? exporttime,
    $core.String? endtime,
    $core.String? failurecode,
    $core.String? s3bucket,
    $core.String? clienttoken,
    ExportType? exporttype,
    $core.String? exportmanifest,
    S3SseAlgorithm? s3ssealgorithm,
    ExportFormat? exportformat,
    $core.String? s3bucketowner,
    $core.String? failuremessage,
    $core.String? starttime,
    $core.String? tablearn,
    $fixnum.Int64? billedsizebytes,
    $core.String? tableid,
    ExportStatus? exportstatus,
  }) {
    final result = create();
    if (exportarn != null) result.exportarn = exportarn;
    if (s3ssekmskeyid != null) result.s3ssekmskeyid = s3ssekmskeyid;
    if (incrementalexportspecification != null)
      result.incrementalexportspecification = incrementalexportspecification;
    if (s3prefix != null) result.s3prefix = s3prefix;
    if (itemcount != null) result.itemcount = itemcount;
    if (exporttime != null) result.exporttime = exporttime;
    if (endtime != null) result.endtime = endtime;
    if (failurecode != null) result.failurecode = failurecode;
    if (s3bucket != null) result.s3bucket = s3bucket;
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (exporttype != null) result.exporttype = exporttype;
    if (exportmanifest != null) result.exportmanifest = exportmanifest;
    if (s3ssealgorithm != null) result.s3ssealgorithm = s3ssealgorithm;
    if (exportformat != null) result.exportformat = exportformat;
    if (s3bucketowner != null) result.s3bucketowner = s3bucketowner;
    if (failuremessage != null) result.failuremessage = failuremessage;
    if (starttime != null) result.starttime = starttime;
    if (tablearn != null) result.tablearn = tablearn;
    if (billedsizebytes != null) result.billedsizebytes = billedsizebytes;
    if (tableid != null) result.tableid = tableid;
    if (exportstatus != null) result.exportstatus = exportstatus;
    return result;
  }

  ExportDescription._();

  factory ExportDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(3661287, _omitFieldNames ? '' : 'exportarn')
    ..aOS(9386110, _omitFieldNames ? '' : 's3ssekmskeyid')
    ..aOM<IncrementalExportSpecification>(
        20792295, _omitFieldNames ? '' : 'incrementalexportspecification',
        subBuilder: IncrementalExportSpecification.create)
    ..aOS(21529336, _omitFieldNames ? '' : 's3prefix')
    ..aInt64(26280022, _omitFieldNames ? '' : 'itemcount')
    ..aOS(28335083, _omitFieldNames ? '' : 'exporttime')
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aOS(84707897, _omitFieldNames ? '' : 'failurecode')
    ..aOS(114031434, _omitFieldNames ? '' : 's3bucket')
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aE<ExportType>(189377484, _omitFieldNames ? '' : 'exporttype',
        enumValues: ExportType.values)
    ..aOS(196735863, _omitFieldNames ? '' : 'exportmanifest')
    ..aE<S3SseAlgorithm>(213032706, _omitFieldNames ? '' : 's3ssealgorithm',
        enumValues: S3SseAlgorithm.values)
    ..aE<ExportFormat>(300799837, _omitFieldNames ? '' : 'exportformat',
        enumValues: ExportFormat.values)
    ..aOS(351576129, _omitFieldNames ? '' : 's3bucketowner')
    ..aOS(353556937, _omitFieldNames ? '' : 'failuremessage')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..aInt64(443315368, _omitFieldNames ? '' : 'billedsizebytes')
    ..aOS(449893011, _omitFieldNames ? '' : 'tableid')
    ..aE<ExportStatus>(459702918, _omitFieldNames ? '' : 'exportstatus',
        enumValues: ExportStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportDescription copyWith(void Function(ExportDescription) updates) =>
      super.copyWith((message) => updates(message as ExportDescription))
          as ExportDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportDescription create() => ExportDescription._();
  @$core.override
  ExportDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportDescription>(create);
  static ExportDescription? _defaultInstance;

  @$pb.TagNumber(3661287)
  $core.String get exportarn => $_getSZ(0);
  @$pb.TagNumber(3661287)
  set exportarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(3661287)
  $core.bool hasExportarn() => $_has(0);
  @$pb.TagNumber(3661287)
  void clearExportarn() => $_clearField(3661287);

  @$pb.TagNumber(9386110)
  $core.String get s3ssekmskeyid => $_getSZ(1);
  @$pb.TagNumber(9386110)
  set s3ssekmskeyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(9386110)
  $core.bool hasS3ssekmskeyid() => $_has(1);
  @$pb.TagNumber(9386110)
  void clearS3ssekmskeyid() => $_clearField(9386110);

  @$pb.TagNumber(20792295)
  IncrementalExportSpecification get incrementalexportspecification =>
      $_getN(2);
  @$pb.TagNumber(20792295)
  set incrementalexportspecification(IncrementalExportSpecification value) =>
      $_setField(20792295, value);
  @$pb.TagNumber(20792295)
  $core.bool hasIncrementalexportspecification() => $_has(2);
  @$pb.TagNumber(20792295)
  void clearIncrementalexportspecification() => $_clearField(20792295);
  @$pb.TagNumber(20792295)
  IncrementalExportSpecification ensureIncrementalexportspecification() =>
      $_ensure(2);

  @$pb.TagNumber(21529336)
  $core.String get s3prefix => $_getSZ(3);
  @$pb.TagNumber(21529336)
  set s3prefix($core.String value) => $_setString(3, value);
  @$pb.TagNumber(21529336)
  $core.bool hasS3prefix() => $_has(3);
  @$pb.TagNumber(21529336)
  void clearS3prefix() => $_clearField(21529336);

  @$pb.TagNumber(26280022)
  $fixnum.Int64 get itemcount => $_getI64(4);
  @$pb.TagNumber(26280022)
  set itemcount($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(26280022)
  $core.bool hasItemcount() => $_has(4);
  @$pb.TagNumber(26280022)
  void clearItemcount() => $_clearField(26280022);

  @$pb.TagNumber(28335083)
  $core.String get exporttime => $_getSZ(5);
  @$pb.TagNumber(28335083)
  set exporttime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(28335083)
  $core.bool hasExporttime() => $_has(5);
  @$pb.TagNumber(28335083)
  void clearExporttime() => $_clearField(28335083);

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(6);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(6, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(6);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(84707897)
  $core.String get failurecode => $_getSZ(7);
  @$pb.TagNumber(84707897)
  set failurecode($core.String value) => $_setString(7, value);
  @$pb.TagNumber(84707897)
  $core.bool hasFailurecode() => $_has(7);
  @$pb.TagNumber(84707897)
  void clearFailurecode() => $_clearField(84707897);

  @$pb.TagNumber(114031434)
  $core.String get s3bucket => $_getSZ(8);
  @$pb.TagNumber(114031434)
  set s3bucket($core.String value) => $_setString(8, value);
  @$pb.TagNumber(114031434)
  $core.bool hasS3bucket() => $_has(8);
  @$pb.TagNumber(114031434)
  void clearS3bucket() => $_clearField(114031434);

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(9);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(9, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(9);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(189377484)
  ExportType get exporttype => $_getN(10);
  @$pb.TagNumber(189377484)
  set exporttype(ExportType value) => $_setField(189377484, value);
  @$pb.TagNumber(189377484)
  $core.bool hasExporttype() => $_has(10);
  @$pb.TagNumber(189377484)
  void clearExporttype() => $_clearField(189377484);

  @$pb.TagNumber(196735863)
  $core.String get exportmanifest => $_getSZ(11);
  @$pb.TagNumber(196735863)
  set exportmanifest($core.String value) => $_setString(11, value);
  @$pb.TagNumber(196735863)
  $core.bool hasExportmanifest() => $_has(11);
  @$pb.TagNumber(196735863)
  void clearExportmanifest() => $_clearField(196735863);

  @$pb.TagNumber(213032706)
  S3SseAlgorithm get s3ssealgorithm => $_getN(12);
  @$pb.TagNumber(213032706)
  set s3ssealgorithm(S3SseAlgorithm value) => $_setField(213032706, value);
  @$pb.TagNumber(213032706)
  $core.bool hasS3ssealgorithm() => $_has(12);
  @$pb.TagNumber(213032706)
  void clearS3ssealgorithm() => $_clearField(213032706);

  @$pb.TagNumber(300799837)
  ExportFormat get exportformat => $_getN(13);
  @$pb.TagNumber(300799837)
  set exportformat(ExportFormat value) => $_setField(300799837, value);
  @$pb.TagNumber(300799837)
  $core.bool hasExportformat() => $_has(13);
  @$pb.TagNumber(300799837)
  void clearExportformat() => $_clearField(300799837);

  @$pb.TagNumber(351576129)
  $core.String get s3bucketowner => $_getSZ(14);
  @$pb.TagNumber(351576129)
  set s3bucketowner($core.String value) => $_setString(14, value);
  @$pb.TagNumber(351576129)
  $core.bool hasS3bucketowner() => $_has(14);
  @$pb.TagNumber(351576129)
  void clearS3bucketowner() => $_clearField(351576129);

  @$pb.TagNumber(353556937)
  $core.String get failuremessage => $_getSZ(15);
  @$pb.TagNumber(353556937)
  set failuremessage($core.String value) => $_setString(15, value);
  @$pb.TagNumber(353556937)
  $core.bool hasFailuremessage() => $_has(15);
  @$pb.TagNumber(353556937)
  void clearFailuremessage() => $_clearField(353556937);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(16);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(16, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(16);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(17);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(17, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(17);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);

  @$pb.TagNumber(443315368)
  $fixnum.Int64 get billedsizebytes => $_getI64(18);
  @$pb.TagNumber(443315368)
  set billedsizebytes($fixnum.Int64 value) => $_setInt64(18, value);
  @$pb.TagNumber(443315368)
  $core.bool hasBilledsizebytes() => $_has(18);
  @$pb.TagNumber(443315368)
  void clearBilledsizebytes() => $_clearField(443315368);

  @$pb.TagNumber(449893011)
  $core.String get tableid => $_getSZ(19);
  @$pb.TagNumber(449893011)
  set tableid($core.String value) => $_setString(19, value);
  @$pb.TagNumber(449893011)
  $core.bool hasTableid() => $_has(19);
  @$pb.TagNumber(449893011)
  void clearTableid() => $_clearField(449893011);

  @$pb.TagNumber(459702918)
  ExportStatus get exportstatus => $_getN(20);
  @$pb.TagNumber(459702918)
  set exportstatus(ExportStatus value) => $_setField(459702918, value);
  @$pb.TagNumber(459702918)
  $core.bool hasExportstatus() => $_has(20);
  @$pb.TagNumber(459702918)
  void clearExportstatus() => $_clearField(459702918);
}

class ExportNotFoundException extends $pb.GeneratedMessage {
  factory ExportNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ExportNotFoundException._();

  factory ExportNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportNotFoundException copyWith(
          void Function(ExportNotFoundException) updates) =>
      super.copyWith((message) => updates(message as ExportNotFoundException))
          as ExportNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportNotFoundException create() => ExportNotFoundException._();
  @$core.override
  ExportNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportNotFoundException>(create);
  static ExportNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ExportSummary extends $pb.GeneratedMessage {
  factory ExportSummary({
    $core.String? exportarn,
    ExportType? exporttype,
    ExportStatus? exportstatus,
  }) {
    final result = create();
    if (exportarn != null) result.exportarn = exportarn;
    if (exporttype != null) result.exporttype = exporttype;
    if (exportstatus != null) result.exportstatus = exportstatus;
    return result;
  }

  ExportSummary._();

  factory ExportSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(3661287, _omitFieldNames ? '' : 'exportarn')
    ..aE<ExportType>(189377484, _omitFieldNames ? '' : 'exporttype',
        enumValues: ExportType.values)
    ..aE<ExportStatus>(459702918, _omitFieldNames ? '' : 'exportstatus',
        enumValues: ExportStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportSummary copyWith(void Function(ExportSummary) updates) =>
      super.copyWith((message) => updates(message as ExportSummary))
          as ExportSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportSummary create() => ExportSummary._();
  @$core.override
  ExportSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportSummary>(create);
  static ExportSummary? _defaultInstance;

  @$pb.TagNumber(3661287)
  $core.String get exportarn => $_getSZ(0);
  @$pb.TagNumber(3661287)
  set exportarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(3661287)
  $core.bool hasExportarn() => $_has(0);
  @$pb.TagNumber(3661287)
  void clearExportarn() => $_clearField(3661287);

  @$pb.TagNumber(189377484)
  ExportType get exporttype => $_getN(1);
  @$pb.TagNumber(189377484)
  set exporttype(ExportType value) => $_setField(189377484, value);
  @$pb.TagNumber(189377484)
  $core.bool hasExporttype() => $_has(1);
  @$pb.TagNumber(189377484)
  void clearExporttype() => $_clearField(189377484);

  @$pb.TagNumber(459702918)
  ExportStatus get exportstatus => $_getN(2);
  @$pb.TagNumber(459702918)
  set exportstatus(ExportStatus value) => $_setField(459702918, value);
  @$pb.TagNumber(459702918)
  $core.bool hasExportstatus() => $_has(2);
  @$pb.TagNumber(459702918)
  void clearExportstatus() => $_clearField(459702918);
}

class ExportTableToPointInTimeInput extends $pb.GeneratedMessage {
  factory ExportTableToPointInTimeInput({
    $core.String? s3ssekmskeyid,
    IncrementalExportSpecification? incrementalexportspecification,
    $core.String? s3prefix,
    $core.String? exporttime,
    $core.String? s3bucket,
    $core.String? clienttoken,
    ExportType? exporttype,
    S3SseAlgorithm? s3ssealgorithm,
    ExportFormat? exportformat,
    $core.String? s3bucketowner,
    $core.String? tablearn,
  }) {
    final result = create();
    if (s3ssekmskeyid != null) result.s3ssekmskeyid = s3ssekmskeyid;
    if (incrementalexportspecification != null)
      result.incrementalexportspecification = incrementalexportspecification;
    if (s3prefix != null) result.s3prefix = s3prefix;
    if (exporttime != null) result.exporttime = exporttime;
    if (s3bucket != null) result.s3bucket = s3bucket;
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (exporttype != null) result.exporttype = exporttype;
    if (s3ssealgorithm != null) result.s3ssealgorithm = s3ssealgorithm;
    if (exportformat != null) result.exportformat = exportformat;
    if (s3bucketowner != null) result.s3bucketowner = s3bucketowner;
    if (tablearn != null) result.tablearn = tablearn;
    return result;
  }

  ExportTableToPointInTimeInput._();

  factory ExportTableToPointInTimeInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportTableToPointInTimeInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportTableToPointInTimeInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(9386110, _omitFieldNames ? '' : 's3ssekmskeyid')
    ..aOM<IncrementalExportSpecification>(
        20792295, _omitFieldNames ? '' : 'incrementalexportspecification',
        subBuilder: IncrementalExportSpecification.create)
    ..aOS(21529336, _omitFieldNames ? '' : 's3prefix')
    ..aOS(28335083, _omitFieldNames ? '' : 'exporttime')
    ..aOS(114031434, _omitFieldNames ? '' : 's3bucket')
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aE<ExportType>(189377484, _omitFieldNames ? '' : 'exporttype',
        enumValues: ExportType.values)
    ..aE<S3SseAlgorithm>(213032706, _omitFieldNames ? '' : 's3ssealgorithm',
        enumValues: S3SseAlgorithm.values)
    ..aE<ExportFormat>(300799837, _omitFieldNames ? '' : 'exportformat',
        enumValues: ExportFormat.values)
    ..aOS(351576129, _omitFieldNames ? '' : 's3bucketowner')
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportTableToPointInTimeInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportTableToPointInTimeInput copyWith(
          void Function(ExportTableToPointInTimeInput) updates) =>
      super.copyWith(
              (message) => updates(message as ExportTableToPointInTimeInput))
          as ExportTableToPointInTimeInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportTableToPointInTimeInput create() =>
      ExportTableToPointInTimeInput._();
  @$core.override
  ExportTableToPointInTimeInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportTableToPointInTimeInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportTableToPointInTimeInput>(create);
  static ExportTableToPointInTimeInput? _defaultInstance;

  @$pb.TagNumber(9386110)
  $core.String get s3ssekmskeyid => $_getSZ(0);
  @$pb.TagNumber(9386110)
  set s3ssekmskeyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(9386110)
  $core.bool hasS3ssekmskeyid() => $_has(0);
  @$pb.TagNumber(9386110)
  void clearS3ssekmskeyid() => $_clearField(9386110);

  @$pb.TagNumber(20792295)
  IncrementalExportSpecification get incrementalexportspecification =>
      $_getN(1);
  @$pb.TagNumber(20792295)
  set incrementalexportspecification(IncrementalExportSpecification value) =>
      $_setField(20792295, value);
  @$pb.TagNumber(20792295)
  $core.bool hasIncrementalexportspecification() => $_has(1);
  @$pb.TagNumber(20792295)
  void clearIncrementalexportspecification() => $_clearField(20792295);
  @$pb.TagNumber(20792295)
  IncrementalExportSpecification ensureIncrementalexportspecification() =>
      $_ensure(1);

  @$pb.TagNumber(21529336)
  $core.String get s3prefix => $_getSZ(2);
  @$pb.TagNumber(21529336)
  set s3prefix($core.String value) => $_setString(2, value);
  @$pb.TagNumber(21529336)
  $core.bool hasS3prefix() => $_has(2);
  @$pb.TagNumber(21529336)
  void clearS3prefix() => $_clearField(21529336);

  @$pb.TagNumber(28335083)
  $core.String get exporttime => $_getSZ(3);
  @$pb.TagNumber(28335083)
  set exporttime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(28335083)
  $core.bool hasExporttime() => $_has(3);
  @$pb.TagNumber(28335083)
  void clearExporttime() => $_clearField(28335083);

  @$pb.TagNumber(114031434)
  $core.String get s3bucket => $_getSZ(4);
  @$pb.TagNumber(114031434)
  set s3bucket($core.String value) => $_setString(4, value);
  @$pb.TagNumber(114031434)
  $core.bool hasS3bucket() => $_has(4);
  @$pb.TagNumber(114031434)
  void clearS3bucket() => $_clearField(114031434);

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(5);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(5, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(5);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(189377484)
  ExportType get exporttype => $_getN(6);
  @$pb.TagNumber(189377484)
  set exporttype(ExportType value) => $_setField(189377484, value);
  @$pb.TagNumber(189377484)
  $core.bool hasExporttype() => $_has(6);
  @$pb.TagNumber(189377484)
  void clearExporttype() => $_clearField(189377484);

  @$pb.TagNumber(213032706)
  S3SseAlgorithm get s3ssealgorithm => $_getN(7);
  @$pb.TagNumber(213032706)
  set s3ssealgorithm(S3SseAlgorithm value) => $_setField(213032706, value);
  @$pb.TagNumber(213032706)
  $core.bool hasS3ssealgorithm() => $_has(7);
  @$pb.TagNumber(213032706)
  void clearS3ssealgorithm() => $_clearField(213032706);

  @$pb.TagNumber(300799837)
  ExportFormat get exportformat => $_getN(8);
  @$pb.TagNumber(300799837)
  set exportformat(ExportFormat value) => $_setField(300799837, value);
  @$pb.TagNumber(300799837)
  $core.bool hasExportformat() => $_has(8);
  @$pb.TagNumber(300799837)
  void clearExportformat() => $_clearField(300799837);

  @$pb.TagNumber(351576129)
  $core.String get s3bucketowner => $_getSZ(9);
  @$pb.TagNumber(351576129)
  set s3bucketowner($core.String value) => $_setString(9, value);
  @$pb.TagNumber(351576129)
  $core.bool hasS3bucketowner() => $_has(9);
  @$pb.TagNumber(351576129)
  void clearS3bucketowner() => $_clearField(351576129);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(10);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(10, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(10);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);
}

class ExportTableToPointInTimeOutput extends $pb.GeneratedMessage {
  factory ExportTableToPointInTimeOutput({
    ExportDescription? exportdescription,
  }) {
    final result = create();
    if (exportdescription != null) result.exportdescription = exportdescription;
    return result;
  }

  ExportTableToPointInTimeOutput._();

  factory ExportTableToPointInTimeOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ExportTableToPointInTimeOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ExportTableToPointInTimeOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ExportDescription>(
        14399944, _omitFieldNames ? '' : 'exportdescription',
        subBuilder: ExportDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportTableToPointInTimeOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ExportTableToPointInTimeOutput copyWith(
          void Function(ExportTableToPointInTimeOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ExportTableToPointInTimeOutput))
          as ExportTableToPointInTimeOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ExportTableToPointInTimeOutput create() =>
      ExportTableToPointInTimeOutput._();
  @$core.override
  ExportTableToPointInTimeOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ExportTableToPointInTimeOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ExportTableToPointInTimeOutput>(create);
  static ExportTableToPointInTimeOutput? _defaultInstance;

  @$pb.TagNumber(14399944)
  ExportDescription get exportdescription => $_getN(0);
  @$pb.TagNumber(14399944)
  set exportdescription(ExportDescription value) => $_setField(14399944, value);
  @$pb.TagNumber(14399944)
  $core.bool hasExportdescription() => $_has(0);
  @$pb.TagNumber(14399944)
  void clearExportdescription() => $_clearField(14399944);
  @$pb.TagNumber(14399944)
  ExportDescription ensureExportdescription() => $_ensure(0);
}

class FailureException extends $pb.GeneratedMessage {
  factory FailureException({
    $core.String? exceptiondescription,
    $core.String? exceptionname,
  }) {
    final result = create();
    if (exceptiondescription != null)
      result.exceptiondescription = exceptiondescription;
    if (exceptionname != null) result.exceptionname = exceptionname;
    return result;
  }

  FailureException._();

  factory FailureException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FailureException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FailureException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(75755113, _omitFieldNames ? '' : 'exceptiondescription')
    ..aOS(449933270, _omitFieldNames ? '' : 'exceptionname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FailureException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FailureException copyWith(void Function(FailureException) updates) =>
      super.copyWith((message) => updates(message as FailureException))
          as FailureException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FailureException create() => FailureException._();
  @$core.override
  FailureException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FailureException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FailureException>(create);
  static FailureException? _defaultInstance;

  @$pb.TagNumber(75755113)
  $core.String get exceptiondescription => $_getSZ(0);
  @$pb.TagNumber(75755113)
  set exceptiondescription($core.String value) => $_setString(0, value);
  @$pb.TagNumber(75755113)
  $core.bool hasExceptiondescription() => $_has(0);
  @$pb.TagNumber(75755113)
  void clearExceptiondescription() => $_clearField(75755113);

  @$pb.TagNumber(449933270)
  $core.String get exceptionname => $_getSZ(1);
  @$pb.TagNumber(449933270)
  set exceptionname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(449933270)
  $core.bool hasExceptionname() => $_has(1);
  @$pb.TagNumber(449933270)
  void clearExceptionname() => $_clearField(449933270);
}

class Get extends $pb.GeneratedMessage {
  factory Get({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.String? projectionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? key,
    $core.String? tablename,
  }) {
    final result = create();
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (projectionexpression != null)
      result.projectionexpression = projectionexpression;
    if (key != null) result.key.addEntries(key);
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  Get._();

  factory Get.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Get.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Get',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'Get.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(150730243, _omitFieldNames ? '' : 'projectionexpression')
    ..m<$core.String, AttributeValue>(219859213, _omitFieldNames ? '' : 'key',
        entryClassName: 'Get.KeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Get clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Get copyWith(void Function(Get) updates) =>
      super.copyWith((message) => updates(message as Get)) as Get;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Get create() => Get._();
  @$core.override
  Get createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Get getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Get>(create);
  static Get? _defaultInstance;

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(0);

  @$pb.TagNumber(150730243)
  $core.String get projectionexpression => $_getSZ(1);
  @$pb.TagNumber(150730243)
  set projectionexpression($core.String value) => $_setString(1, value);
  @$pb.TagNumber(150730243)
  $core.bool hasProjectionexpression() => $_has(1);
  @$pb.TagNumber(150730243)
  void clearProjectionexpression() => $_clearField(150730243);

  @$pb.TagNumber(219859213)
  $pb.PbMap<$core.String, AttributeValue> get key => $_getMap(2);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class GetItemInput extends $pb.GeneratedMessage {
  factory GetItemInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.String? projectionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? key,
    $core.String? tablename,
    $core.Iterable<$core.String>? attributestoget,
    $core.bool? consistentread,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (projectionexpression != null)
      result.projectionexpression = projectionexpression;
    if (key != null) result.key.addEntries(key);
    if (tablename != null) result.tablename = tablename;
    if (attributestoget != null) result.attributestoget.addAll(attributestoget);
    if (consistentread != null) result.consistentread = consistentread;
    return result;
  }

  GetItemInput._();

  factory GetItemInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetItemInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetItemInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'GetItemInput.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(150730243, _omitFieldNames ? '' : 'projectionexpression')
    ..m<$core.String, AttributeValue>(219859213, _omitFieldNames ? '' : 'key',
        entryClassName: 'GetItemInput.KeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..pPS(311382592, _omitFieldNames ? '' : 'attributestoget')
    ..aOB(531556994, _omitFieldNames ? '' : 'consistentread')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetItemInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetItemInput copyWith(void Function(GetItemInput) updates) =>
      super.copyWith((message) => updates(message as GetItemInput))
          as GetItemInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetItemInput create() => GetItemInput._();
  @$core.override
  GetItemInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetItemInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetItemInput>(create);
  static GetItemInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(1);

  @$pb.TagNumber(150730243)
  $core.String get projectionexpression => $_getSZ(2);
  @$pb.TagNumber(150730243)
  set projectionexpression($core.String value) => $_setString(2, value);
  @$pb.TagNumber(150730243)
  $core.bool hasProjectionexpression() => $_has(2);
  @$pb.TagNumber(150730243)
  void clearProjectionexpression() => $_clearField(150730243);

  @$pb.TagNumber(219859213)
  $pb.PbMap<$core.String, AttributeValue> get key => $_getMap(3);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(4);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(4, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(4);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(311382592)
  $pb.PbList<$core.String> get attributestoget => $_getList(5);

  @$pb.TagNumber(531556994)
  $core.bool get consistentread => $_getBF(6);
  @$pb.TagNumber(531556994)
  set consistentread($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(531556994)
  $core.bool hasConsistentread() => $_has(6);
  @$pb.TagNumber(531556994)
  void clearConsistentread() => $_clearField(531556994);
}

class GetItemOutput extends $pb.GeneratedMessage {
  factory GetItemOutput({
    ConsumedCapacity? consumedcapacity,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (consumedcapacity != null) result.consumedcapacity = consumedcapacity;
    if (item != null) result.item.addEntries(item);
    return result;
  }

  GetItemOutput._();

  factory GetItemOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetItemOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetItemOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'GetItemOutput.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetItemOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetItemOutput copyWith(void Function(GetItemOutput) updates) =>
      super.copyWith((message) => updates(message as GetItemOutput))
          as GetItemOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetItemOutput create() => GetItemOutput._();
  @$core.override
  GetItemOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetItemOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetItemOutput>(create);
  static GetItemOutput? _defaultInstance;

  @$pb.TagNumber(449336620)
  ConsumedCapacity get consumedcapacity => $_getN(0);
  @$pb.TagNumber(449336620)
  set consumedcapacity(ConsumedCapacity value) => $_setField(449336620, value);
  @$pb.TagNumber(449336620)
  $core.bool hasConsumedcapacity() => $_has(0);
  @$pb.TagNumber(449336620)
  void clearConsumedcapacity() => $_clearField(449336620);
  @$pb.TagNumber(449336620)
  ConsumedCapacity ensureConsumedcapacity() => $_ensure(0);

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(1);
}

class GetResourcePolicyInput extends $pb.GeneratedMessage {
  factory GetResourcePolicyInput({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  GetResourcePolicyInput._();

  factory GetResourcePolicyInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourcePolicyInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourcePolicyInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePolicyInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePolicyInput copyWith(
          void Function(GetResourcePolicyInput) updates) =>
      super.copyWith((message) => updates(message as GetResourcePolicyInput))
          as GetResourcePolicyInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourcePolicyInput create() => GetResourcePolicyInput._();
  @$core.override
  GetResourcePolicyInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourcePolicyInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourcePolicyInput>(create);
  static GetResourcePolicyInput? _defaultInstance;

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class GetResourcePolicyOutput extends $pb.GeneratedMessage {
  factory GetResourcePolicyOutput({
    $core.String? policy,
    $core.String? revisionid,
  }) {
    final result = create();
    if (policy != null) result.policy = policy;
    if (revisionid != null) result.revisionid = revisionid;
    return result;
  }

  GetResourcePolicyOutput._();

  factory GetResourcePolicyOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetResourcePolicyOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetResourcePolicyOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..aOS(499618182, _omitFieldNames ? '' : 'revisionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePolicyOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetResourcePolicyOutput copyWith(
          void Function(GetResourcePolicyOutput) updates) =>
      super.copyWith((message) => updates(message as GetResourcePolicyOutput))
          as GetResourcePolicyOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetResourcePolicyOutput create() => GetResourcePolicyOutput._();
  @$core.override
  GetResourcePolicyOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetResourcePolicyOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetResourcePolicyOutput>(create);
  static GetResourcePolicyOutput? _defaultInstance;

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(0);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(0, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(0);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);

  @$pb.TagNumber(499618182)
  $core.String get revisionid => $_getSZ(1);
  @$pb.TagNumber(499618182)
  set revisionid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(499618182)
  $core.bool hasRevisionid() => $_has(1);
  @$pb.TagNumber(499618182)
  void clearRevisionid() => $_clearField(499618182);
}

class GlobalSecondaryIndex extends $pb.GeneratedMessage {
  factory GlobalSecondaryIndex({
    ProvisionedThroughput? provisionedthroughput,
    $core.String? indexname,
    Projection? projection,
    WarmThroughput? warmthroughput,
    $core.Iterable<KeySchemaElement>? keyschema,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (indexname != null) result.indexname = indexname;
    if (projection != null) result.projection = projection;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  GlobalSecondaryIndex._();

  factory GlobalSecondaryIndex.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalSecondaryIndex.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalSecondaryIndex',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughput>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughput.create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<Projection>(105045921, _omitFieldNames ? '' : 'projection',
        subBuilder: Projection.create)
    ..aOM<WarmThroughput>(290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: WarmThroughput.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndex clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndex copyWith(void Function(GlobalSecondaryIndex) updates) =>
      super.copyWith((message) => updates(message as GlobalSecondaryIndex))
          as GlobalSecondaryIndex;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndex create() => GlobalSecondaryIndex._();
  @$core.override
  GlobalSecondaryIndex createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndex getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalSecondaryIndex>(create);
  static GlobalSecondaryIndex? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughput get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughput value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughput ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(105045921)
  Projection get projection => $_getN(2);
  @$pb.TagNumber(105045921)
  set projection(Projection value) => $_setField(105045921, value);
  @$pb.TagNumber(105045921)
  $core.bool hasProjection() => $_has(2);
  @$pb.TagNumber(105045921)
  void clearProjection() => $_clearField(105045921);
  @$pb.TagNumber(105045921)
  Projection ensureProjection() => $_ensure(2);

  @$pb.TagNumber(290598659)
  WarmThroughput get warmthroughput => $_getN(3);
  @$pb.TagNumber(290598659)
  set warmthroughput(WarmThroughput value) => $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(3);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  WarmThroughput ensureWarmthroughput() => $_ensure(3);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(4);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(5);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(5);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(5);
}

class GlobalSecondaryIndexAutoScalingUpdate extends $pb.GeneratedMessage {
  factory GlobalSecondaryIndexAutoScalingUpdate({
    $core.String? indexname,
    AutoScalingSettingsUpdate? provisionedwritecapacityautoscalingupdate,
  }) {
    final result = create();
    if (indexname != null) result.indexname = indexname;
    if (provisionedwritecapacityautoscalingupdate != null)
      result.provisionedwritecapacityautoscalingupdate =
          provisionedwritecapacityautoscalingupdate;
    return result;
  }

  GlobalSecondaryIndexAutoScalingUpdate._();

  factory GlobalSecondaryIndexAutoScalingUpdate.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalSecondaryIndexAutoScalingUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalSecondaryIndexAutoScalingUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<AutoScalingSettingsUpdate>(316618294,
        _omitFieldNames ? '' : 'provisionedwritecapacityautoscalingupdate',
        subBuilder: AutoScalingSettingsUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexAutoScalingUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexAutoScalingUpdate copyWith(
          void Function(GlobalSecondaryIndexAutoScalingUpdate) updates) =>
      super.copyWith((message) =>
              updates(message as GlobalSecondaryIndexAutoScalingUpdate))
          as GlobalSecondaryIndexAutoScalingUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexAutoScalingUpdate create() =>
      GlobalSecondaryIndexAutoScalingUpdate._();
  @$core.override
  GlobalSecondaryIndexAutoScalingUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexAutoScalingUpdate getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GlobalSecondaryIndexAutoScalingUpdate>(create);
  static GlobalSecondaryIndexAutoScalingUpdate? _defaultInstance;

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(0);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(0);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(316618294)
  AutoScalingSettingsUpdate get provisionedwritecapacityautoscalingupdate =>
      $_getN(1);
  @$pb.TagNumber(316618294)
  set provisionedwritecapacityautoscalingupdate(
          AutoScalingSettingsUpdate value) =>
      $_setField(316618294, value);
  @$pb.TagNumber(316618294)
  $core.bool hasProvisionedwritecapacityautoscalingupdate() => $_has(1);
  @$pb.TagNumber(316618294)
  void clearProvisionedwritecapacityautoscalingupdate() =>
      $_clearField(316618294);
  @$pb.TagNumber(316618294)
  AutoScalingSettingsUpdate ensureProvisionedwritecapacityautoscalingupdate() =>
      $_ensure(1);
}

class GlobalSecondaryIndexDescription extends $pb.GeneratedMessage {
  factory GlobalSecondaryIndexDescription({
    ProvisionedThroughputDescription? provisionedthroughput,
    $fixnum.Int64? itemcount,
    $core.String? indexname,
    Projection? projection,
    $core.bool? backfilling,
    GlobalSecondaryIndexWarmThroughputDescription? warmthroughput,
    $core.Iterable<KeySchemaElement>? keyschema,
    IndexStatus? indexstatus,
    $core.String? indexarn,
    $fixnum.Int64? indexsizebytes,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (itemcount != null) result.itemcount = itemcount;
    if (indexname != null) result.indexname = indexname;
    if (projection != null) result.projection = projection;
    if (backfilling != null) result.backfilling = backfilling;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (indexstatus != null) result.indexstatus = indexstatus;
    if (indexarn != null) result.indexarn = indexarn;
    if (indexsizebytes != null) result.indexsizebytes = indexsizebytes;
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  GlobalSecondaryIndexDescription._();

  factory GlobalSecondaryIndexDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalSecondaryIndexDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalSecondaryIndexDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughputDescription>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughputDescription.create)
    ..aInt64(26280022, _omitFieldNames ? '' : 'itemcount')
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<Projection>(105045921, _omitFieldNames ? '' : 'projection',
        subBuilder: Projection.create)
    ..aOB(251413370, _omitFieldNames ? '' : 'backfilling')
    ..aOM<GlobalSecondaryIndexWarmThroughputDescription>(
        290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: GlobalSecondaryIndexWarmThroughputDescription.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..aE<IndexStatus>(364436830, _omitFieldNames ? '' : 'indexstatus',
        enumValues: IndexStatus.values)
    ..aOS(374335615, _omitFieldNames ? '' : 'indexarn')
    ..aInt64(395738346, _omitFieldNames ? '' : 'indexsizebytes')
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexDescription copyWith(
          void Function(GlobalSecondaryIndexDescription) updates) =>
      super.copyWith(
              (message) => updates(message as GlobalSecondaryIndexDescription))
          as GlobalSecondaryIndexDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexDescription create() =>
      GlobalSecondaryIndexDescription._();
  @$core.override
  GlobalSecondaryIndexDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalSecondaryIndexDescription>(
          create);
  static GlobalSecondaryIndexDescription? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughputDescription get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughputDescription value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughputDescription ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(26280022)
  $fixnum.Int64 get itemcount => $_getI64(1);
  @$pb.TagNumber(26280022)
  set itemcount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(26280022)
  $core.bool hasItemcount() => $_has(1);
  @$pb.TagNumber(26280022)
  void clearItemcount() => $_clearField(26280022);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(2);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(2);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(105045921)
  Projection get projection => $_getN(3);
  @$pb.TagNumber(105045921)
  set projection(Projection value) => $_setField(105045921, value);
  @$pb.TagNumber(105045921)
  $core.bool hasProjection() => $_has(3);
  @$pb.TagNumber(105045921)
  void clearProjection() => $_clearField(105045921);
  @$pb.TagNumber(105045921)
  Projection ensureProjection() => $_ensure(3);

  @$pb.TagNumber(251413370)
  $core.bool get backfilling => $_getBF(4);
  @$pb.TagNumber(251413370)
  set backfilling($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(251413370)
  $core.bool hasBackfilling() => $_has(4);
  @$pb.TagNumber(251413370)
  void clearBackfilling() => $_clearField(251413370);

  @$pb.TagNumber(290598659)
  GlobalSecondaryIndexWarmThroughputDescription get warmthroughput => $_getN(5);
  @$pb.TagNumber(290598659)
  set warmthroughput(GlobalSecondaryIndexWarmThroughputDescription value) =>
      $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(5);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  GlobalSecondaryIndexWarmThroughputDescription ensureWarmthroughput() =>
      $_ensure(5);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(6);

  @$pb.TagNumber(364436830)
  IndexStatus get indexstatus => $_getN(7);
  @$pb.TagNumber(364436830)
  set indexstatus(IndexStatus value) => $_setField(364436830, value);
  @$pb.TagNumber(364436830)
  $core.bool hasIndexstatus() => $_has(7);
  @$pb.TagNumber(364436830)
  void clearIndexstatus() => $_clearField(364436830);

  @$pb.TagNumber(374335615)
  $core.String get indexarn => $_getSZ(8);
  @$pb.TagNumber(374335615)
  set indexarn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(374335615)
  $core.bool hasIndexarn() => $_has(8);
  @$pb.TagNumber(374335615)
  void clearIndexarn() => $_clearField(374335615);

  @$pb.TagNumber(395738346)
  $fixnum.Int64 get indexsizebytes => $_getI64(9);
  @$pb.TagNumber(395738346)
  set indexsizebytes($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(395738346)
  $core.bool hasIndexsizebytes() => $_has(9);
  @$pb.TagNumber(395738346)
  void clearIndexsizebytes() => $_clearField(395738346);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(10);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(10);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(10);
}

class GlobalSecondaryIndexInfo extends $pb.GeneratedMessage {
  factory GlobalSecondaryIndexInfo({
    ProvisionedThroughput? provisionedthroughput,
    $core.String? indexname,
    Projection? projection,
    $core.Iterable<KeySchemaElement>? keyschema,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (indexname != null) result.indexname = indexname;
    if (projection != null) result.projection = projection;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  GlobalSecondaryIndexInfo._();

  factory GlobalSecondaryIndexInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalSecondaryIndexInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalSecondaryIndexInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughput>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughput.create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<Projection>(105045921, _omitFieldNames ? '' : 'projection',
        subBuilder: Projection.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexInfo copyWith(
          void Function(GlobalSecondaryIndexInfo) updates) =>
      super.copyWith((message) => updates(message as GlobalSecondaryIndexInfo))
          as GlobalSecondaryIndexInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexInfo create() => GlobalSecondaryIndexInfo._();
  @$core.override
  GlobalSecondaryIndexInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalSecondaryIndexInfo>(create);
  static GlobalSecondaryIndexInfo? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughput get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughput value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughput ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(105045921)
  Projection get projection => $_getN(2);
  @$pb.TagNumber(105045921)
  set projection(Projection value) => $_setField(105045921, value);
  @$pb.TagNumber(105045921)
  $core.bool hasProjection() => $_has(2);
  @$pb.TagNumber(105045921)
  void clearProjection() => $_clearField(105045921);
  @$pb.TagNumber(105045921)
  Projection ensureProjection() => $_ensure(2);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(3);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(4);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(4);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(4);
}

class GlobalSecondaryIndexUpdate extends $pb.GeneratedMessage {
  factory GlobalSecondaryIndexUpdate({
    UpdateGlobalSecondaryIndexAction? update,
    DeleteGlobalSecondaryIndexAction? delete,
    CreateGlobalSecondaryIndexAction? create_420340862,
  }) {
    final result = create();
    if (update != null) result.update = update;
    if (delete != null) result.delete = delete;
    if (create_420340862 != null) result.create_420340862 = create_420340862;
    return result;
  }

  GlobalSecondaryIndexUpdate._();

  factory GlobalSecondaryIndexUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalSecondaryIndexUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalSecondaryIndexUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<UpdateGlobalSecondaryIndexAction>(
        237178517, _omitFieldNames ? '' : 'update',
        subBuilder: UpdateGlobalSecondaryIndexAction.create)
    ..aOM<DeleteGlobalSecondaryIndexAction>(
        395831915, _omitFieldNames ? '' : 'delete',
        subBuilder: DeleteGlobalSecondaryIndexAction.create)
    ..aOM<CreateGlobalSecondaryIndexAction>(
        420340862, _omitFieldNames ? '' : 'create',
        subBuilder: CreateGlobalSecondaryIndexAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexUpdate copyWith(
          void Function(GlobalSecondaryIndexUpdate) updates) =>
      super.copyWith(
              (message) => updates(message as GlobalSecondaryIndexUpdate))
          as GlobalSecondaryIndexUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexUpdate create() => GlobalSecondaryIndexUpdate._();
  @$core.override
  GlobalSecondaryIndexUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalSecondaryIndexUpdate>(create);
  static GlobalSecondaryIndexUpdate? _defaultInstance;

  @$pb.TagNumber(237178517)
  UpdateGlobalSecondaryIndexAction get update => $_getN(0);
  @$pb.TagNumber(237178517)
  set update(UpdateGlobalSecondaryIndexAction value) =>
      $_setField(237178517, value);
  @$pb.TagNumber(237178517)
  $core.bool hasUpdate() => $_has(0);
  @$pb.TagNumber(237178517)
  void clearUpdate() => $_clearField(237178517);
  @$pb.TagNumber(237178517)
  UpdateGlobalSecondaryIndexAction ensureUpdate() => $_ensure(0);

  @$pb.TagNumber(395831915)
  DeleteGlobalSecondaryIndexAction get delete => $_getN(1);
  @$pb.TagNumber(395831915)
  set delete(DeleteGlobalSecondaryIndexAction value) =>
      $_setField(395831915, value);
  @$pb.TagNumber(395831915)
  $core.bool hasDelete() => $_has(1);
  @$pb.TagNumber(395831915)
  void clearDelete() => $_clearField(395831915);
  @$pb.TagNumber(395831915)
  DeleteGlobalSecondaryIndexAction ensureDelete() => $_ensure(1);

  @$pb.TagNumber(420340862)
  CreateGlobalSecondaryIndexAction get create_420340862 => $_getN(2);
  @$pb.TagNumber(420340862)
  set create_420340862(CreateGlobalSecondaryIndexAction value) =>
      $_setField(420340862, value);
  @$pb.TagNumber(420340862)
  $core.bool hasCreate_420340862() => $_has(2);
  @$pb.TagNumber(420340862)
  void clearCreate_420340862() => $_clearField(420340862);
  @$pb.TagNumber(420340862)
  CreateGlobalSecondaryIndexAction ensureCreate_420340862() => $_ensure(2);
}

class GlobalSecondaryIndexWarmThroughputDescription
    extends $pb.GeneratedMessage {
  factory GlobalSecondaryIndexWarmThroughputDescription({
    IndexStatus? status,
    $fixnum.Int64? readunitspersecond,
    $fixnum.Int64? writeunitspersecond,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (readunitspersecond != null)
      result.readunitspersecond = readunitspersecond;
    if (writeunitspersecond != null)
      result.writeunitspersecond = writeunitspersecond;
    return result;
  }

  GlobalSecondaryIndexWarmThroughputDescription._();

  factory GlobalSecondaryIndexWarmThroughputDescription.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalSecondaryIndexWarmThroughputDescription.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalSecondaryIndexWarmThroughputDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<IndexStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: IndexStatus.values)
    ..aInt64(11400732, _omitFieldNames ? '' : 'readunitspersecond')
    ..aInt64(339770127, _omitFieldNames ? '' : 'writeunitspersecond')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexWarmThroughputDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalSecondaryIndexWarmThroughputDescription copyWith(
          void Function(GlobalSecondaryIndexWarmThroughputDescription)
              updates) =>
      super.copyWith((message) =>
              updates(message as GlobalSecondaryIndexWarmThroughputDescription))
          as GlobalSecondaryIndexWarmThroughputDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexWarmThroughputDescription create() =>
      GlobalSecondaryIndexWarmThroughputDescription._();
  @$core.override
  GlobalSecondaryIndexWarmThroughputDescription createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static GlobalSecondaryIndexWarmThroughputDescription getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GlobalSecondaryIndexWarmThroughputDescription>(create);
  static GlobalSecondaryIndexWarmThroughputDescription? _defaultInstance;

  @$pb.TagNumber(6222352)
  IndexStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(IndexStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(11400732)
  $fixnum.Int64 get readunitspersecond => $_getI64(1);
  @$pb.TagNumber(11400732)
  set readunitspersecond($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(11400732)
  $core.bool hasReadunitspersecond() => $_has(1);
  @$pb.TagNumber(11400732)
  void clearReadunitspersecond() => $_clearField(11400732);

  @$pb.TagNumber(339770127)
  $fixnum.Int64 get writeunitspersecond => $_getI64(2);
  @$pb.TagNumber(339770127)
  set writeunitspersecond($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(339770127)
  $core.bool hasWriteunitspersecond() => $_has(2);
  @$pb.TagNumber(339770127)
  void clearWriteunitspersecond() => $_clearField(339770127);
}

class GlobalTable extends $pb.GeneratedMessage {
  factory GlobalTable({
    $core.Iterable<Replica>? replicationgroup,
    $core.String? globaltablename,
  }) {
    final result = create();
    if (replicationgroup != null)
      result.replicationgroup.addAll(replicationgroup);
    if (globaltablename != null) result.globaltablename = globaltablename;
    return result;
  }

  GlobalTable._();

  factory GlobalTable.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalTable.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalTable',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<Replica>(190970785, _omitFieldNames ? '' : 'replicationgroup',
        subBuilder: Replica.create)
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTable clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTable copyWith(void Function(GlobalTable) updates) =>
      super.copyWith((message) => updates(message as GlobalTable))
          as GlobalTable;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalTable create() => GlobalTable._();
  @$core.override
  GlobalTable createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalTable getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalTable>(create);
  static GlobalTable? _defaultInstance;

  @$pb.TagNumber(190970785)
  $pb.PbList<Replica> get replicationgroup => $_getList(0);

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(1);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(1);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);
}

class GlobalTableAlreadyExistsException extends $pb.GeneratedMessage {
  factory GlobalTableAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  GlobalTableAlreadyExistsException._();

  factory GlobalTableAlreadyExistsException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalTableAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalTableAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableAlreadyExistsException copyWith(
          void Function(GlobalTableAlreadyExistsException) updates) =>
      super.copyWith((message) =>
              updates(message as GlobalTableAlreadyExistsException))
          as GlobalTableAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalTableAlreadyExistsException create() =>
      GlobalTableAlreadyExistsException._();
  @$core.override
  GlobalTableAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalTableAlreadyExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalTableAlreadyExistsException>(
          create);
  static GlobalTableAlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class GlobalTableDescription extends $pb.GeneratedMessage {
  factory GlobalTableDescription({
    GlobalTableStatus? globaltablestatus,
    $core.String? creationdatetime,
    $core.Iterable<ReplicaDescription>? replicationgroup,
    $core.String? globaltablearn,
    $core.String? globaltablename,
  }) {
    final result = create();
    if (globaltablestatus != null) result.globaltablestatus = globaltablestatus;
    if (creationdatetime != null) result.creationdatetime = creationdatetime;
    if (replicationgroup != null)
      result.replicationgroup.addAll(replicationgroup);
    if (globaltablearn != null) result.globaltablearn = globaltablearn;
    if (globaltablename != null) result.globaltablename = globaltablename;
    return result;
  }

  GlobalTableDescription._();

  factory GlobalTableDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalTableDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalTableDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<GlobalTableStatus>(8833293, _omitFieldNames ? '' : 'globaltablestatus',
        enumValues: GlobalTableStatus.values)
    ..aOS(48904698, _omitFieldNames ? '' : 'creationdatetime')
    ..pPM<ReplicaDescription>(
        190970785, _omitFieldNames ? '' : 'replicationgroup',
        subBuilder: ReplicaDescription.create)
    ..aOS(262302830, _omitFieldNames ? '' : 'globaltablearn')
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableDescription copyWith(
          void Function(GlobalTableDescription) updates) =>
      super.copyWith((message) => updates(message as GlobalTableDescription))
          as GlobalTableDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalTableDescription create() => GlobalTableDescription._();
  @$core.override
  GlobalTableDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalTableDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalTableDescription>(create);
  static GlobalTableDescription? _defaultInstance;

  @$pb.TagNumber(8833293)
  GlobalTableStatus get globaltablestatus => $_getN(0);
  @$pb.TagNumber(8833293)
  set globaltablestatus(GlobalTableStatus value) => $_setField(8833293, value);
  @$pb.TagNumber(8833293)
  $core.bool hasGlobaltablestatus() => $_has(0);
  @$pb.TagNumber(8833293)
  void clearGlobaltablestatus() => $_clearField(8833293);

  @$pb.TagNumber(48904698)
  $core.String get creationdatetime => $_getSZ(1);
  @$pb.TagNumber(48904698)
  set creationdatetime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(48904698)
  $core.bool hasCreationdatetime() => $_has(1);
  @$pb.TagNumber(48904698)
  void clearCreationdatetime() => $_clearField(48904698);

  @$pb.TagNumber(190970785)
  $pb.PbList<ReplicaDescription> get replicationgroup => $_getList(2);

  @$pb.TagNumber(262302830)
  $core.String get globaltablearn => $_getSZ(3);
  @$pb.TagNumber(262302830)
  set globaltablearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(262302830)
  $core.bool hasGlobaltablearn() => $_has(3);
  @$pb.TagNumber(262302830)
  void clearGlobaltablearn() => $_clearField(262302830);

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(4);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(4, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(4);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);
}

class GlobalTableGlobalSecondaryIndexSettingsUpdate
    extends $pb.GeneratedMessage {
  factory GlobalTableGlobalSecondaryIndexSettingsUpdate({
    $core.String? indexname,
    $fixnum.Int64? provisionedwritecapacityunits,
    AutoScalingSettingsUpdate?
        provisionedwritecapacityautoscalingsettingsupdate,
  }) {
    final result = create();
    if (indexname != null) result.indexname = indexname;
    if (provisionedwritecapacityunits != null)
      result.provisionedwritecapacityunits = provisionedwritecapacityunits;
    if (provisionedwritecapacityautoscalingsettingsupdate != null)
      result.provisionedwritecapacityautoscalingsettingsupdate =
          provisionedwritecapacityautoscalingsettingsupdate;
    return result;
  }

  GlobalTableGlobalSecondaryIndexSettingsUpdate._();

  factory GlobalTableGlobalSecondaryIndexSettingsUpdate.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalTableGlobalSecondaryIndexSettingsUpdate.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalTableGlobalSecondaryIndexSettingsUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aInt64(225881684, _omitFieldNames ? '' : 'provisionedwritecapacityunits')
    ..aOM<AutoScalingSettingsUpdate>(
        302140761,
        _omitFieldNames
            ? ''
            : 'provisionedwritecapacityautoscalingsettingsupdate',
        subBuilder: AutoScalingSettingsUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableGlobalSecondaryIndexSettingsUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableGlobalSecondaryIndexSettingsUpdate copyWith(
          void Function(GlobalTableGlobalSecondaryIndexSettingsUpdate)
              updates) =>
      super.copyWith((message) =>
              updates(message as GlobalTableGlobalSecondaryIndexSettingsUpdate))
          as GlobalTableGlobalSecondaryIndexSettingsUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalTableGlobalSecondaryIndexSettingsUpdate create() =>
      GlobalTableGlobalSecondaryIndexSettingsUpdate._();
  @$core.override
  GlobalTableGlobalSecondaryIndexSettingsUpdate createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static GlobalTableGlobalSecondaryIndexSettingsUpdate getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GlobalTableGlobalSecondaryIndexSettingsUpdate>(create);
  static GlobalTableGlobalSecondaryIndexSettingsUpdate? _defaultInstance;

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(0);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(0);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(225881684)
  $fixnum.Int64 get provisionedwritecapacityunits => $_getI64(1);
  @$pb.TagNumber(225881684)
  set provisionedwritecapacityunits($fixnum.Int64 value) =>
      $_setInt64(1, value);
  @$pb.TagNumber(225881684)
  $core.bool hasProvisionedwritecapacityunits() => $_has(1);
  @$pb.TagNumber(225881684)
  void clearProvisionedwritecapacityunits() => $_clearField(225881684);

  @$pb.TagNumber(302140761)
  AutoScalingSettingsUpdate
      get provisionedwritecapacityautoscalingsettingsupdate => $_getN(2);
  @$pb.TagNumber(302140761)
  set provisionedwritecapacityautoscalingsettingsupdate(
          AutoScalingSettingsUpdate value) =>
      $_setField(302140761, value);
  @$pb.TagNumber(302140761)
  $core.bool hasProvisionedwritecapacityautoscalingsettingsupdate() => $_has(2);
  @$pb.TagNumber(302140761)
  void clearProvisionedwritecapacityautoscalingsettingsupdate() =>
      $_clearField(302140761);
  @$pb.TagNumber(302140761)
  AutoScalingSettingsUpdate
      ensureProvisionedwritecapacityautoscalingsettingsupdate() => $_ensure(2);
}

class GlobalTableNotFoundException extends $pb.GeneratedMessage {
  factory GlobalTableNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  GlobalTableNotFoundException._();

  factory GlobalTableNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalTableNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalTableNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableNotFoundException copyWith(
          void Function(GlobalTableNotFoundException) updates) =>
      super.copyWith(
              (message) => updates(message as GlobalTableNotFoundException))
          as GlobalTableNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalTableNotFoundException create() =>
      GlobalTableNotFoundException._();
  @$core.override
  GlobalTableNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalTableNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalTableNotFoundException>(create);
  static GlobalTableNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class GlobalTableWitnessDescription extends $pb.GeneratedMessage {
  factory GlobalTableWitnessDescription({
    $core.String? regionname,
    WitnessStatus? witnessstatus,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    if (witnessstatus != null) result.witnessstatus = witnessstatus;
    return result;
  }

  GlobalTableWitnessDescription._();

  factory GlobalTableWitnessDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalTableWitnessDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalTableWitnessDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..aE<WitnessStatus>(303352887, _omitFieldNames ? '' : 'witnessstatus',
        enumValues: WitnessStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableWitnessDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableWitnessDescription copyWith(
          void Function(GlobalTableWitnessDescription) updates) =>
      super.copyWith(
              (message) => updates(message as GlobalTableWitnessDescription))
          as GlobalTableWitnessDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalTableWitnessDescription create() =>
      GlobalTableWitnessDescription._();
  @$core.override
  GlobalTableWitnessDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalTableWitnessDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalTableWitnessDescription>(create);
  static GlobalTableWitnessDescription? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(303352887)
  WitnessStatus get witnessstatus => $_getN(1);
  @$pb.TagNumber(303352887)
  set witnessstatus(WitnessStatus value) => $_setField(303352887, value);
  @$pb.TagNumber(303352887)
  $core.bool hasWitnessstatus() => $_has(1);
  @$pb.TagNumber(303352887)
  void clearWitnessstatus() => $_clearField(303352887);
}

class GlobalTableWitnessGroupUpdate extends $pb.GeneratedMessage {
  factory GlobalTableWitnessGroupUpdate({
    DeleteGlobalTableWitnessGroupMemberAction? delete,
    CreateGlobalTableWitnessGroupMemberAction? create_420340862,
  }) {
    final result = create();
    if (delete != null) result.delete = delete;
    if (create_420340862 != null) result.create_420340862 = create_420340862;
    return result;
  }

  GlobalTableWitnessGroupUpdate._();

  factory GlobalTableWitnessGroupUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GlobalTableWitnessGroupUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GlobalTableWitnessGroupUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<DeleteGlobalTableWitnessGroupMemberAction>(
        395831915, _omitFieldNames ? '' : 'delete',
        subBuilder: DeleteGlobalTableWitnessGroupMemberAction.create)
    ..aOM<CreateGlobalTableWitnessGroupMemberAction>(
        420340862, _omitFieldNames ? '' : 'create',
        subBuilder: CreateGlobalTableWitnessGroupMemberAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableWitnessGroupUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GlobalTableWitnessGroupUpdate copyWith(
          void Function(GlobalTableWitnessGroupUpdate) updates) =>
      super.copyWith(
              (message) => updates(message as GlobalTableWitnessGroupUpdate))
          as GlobalTableWitnessGroupUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GlobalTableWitnessGroupUpdate create() =>
      GlobalTableWitnessGroupUpdate._();
  @$core.override
  GlobalTableWitnessGroupUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GlobalTableWitnessGroupUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GlobalTableWitnessGroupUpdate>(create);
  static GlobalTableWitnessGroupUpdate? _defaultInstance;

  @$pb.TagNumber(395831915)
  DeleteGlobalTableWitnessGroupMemberAction get delete => $_getN(0);
  @$pb.TagNumber(395831915)
  set delete(DeleteGlobalTableWitnessGroupMemberAction value) =>
      $_setField(395831915, value);
  @$pb.TagNumber(395831915)
  $core.bool hasDelete() => $_has(0);
  @$pb.TagNumber(395831915)
  void clearDelete() => $_clearField(395831915);
  @$pb.TagNumber(395831915)
  DeleteGlobalTableWitnessGroupMemberAction ensureDelete() => $_ensure(0);

  @$pb.TagNumber(420340862)
  CreateGlobalTableWitnessGroupMemberAction get create_420340862 => $_getN(1);
  @$pb.TagNumber(420340862)
  set create_420340862(CreateGlobalTableWitnessGroupMemberAction value) =>
      $_setField(420340862, value);
  @$pb.TagNumber(420340862)
  $core.bool hasCreate_420340862() => $_has(1);
  @$pb.TagNumber(420340862)
  void clearCreate_420340862() => $_clearField(420340862);
  @$pb.TagNumber(420340862)
  CreateGlobalTableWitnessGroupMemberAction ensureCreate_420340862() =>
      $_ensure(1);
}

class IdempotentParameterMismatchException extends $pb.GeneratedMessage {
  factory IdempotentParameterMismatchException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IdempotentParameterMismatchException._();

  factory IdempotentParameterMismatchException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IdempotentParameterMismatchException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IdempotentParameterMismatchException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdempotentParameterMismatchException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IdempotentParameterMismatchException copyWith(
          void Function(IdempotentParameterMismatchException) updates) =>
      super.copyWith((message) =>
              updates(message as IdempotentParameterMismatchException))
          as IdempotentParameterMismatchException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IdempotentParameterMismatchException create() =>
      IdempotentParameterMismatchException._();
  @$core.override
  IdempotentParameterMismatchException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IdempotentParameterMismatchException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          IdempotentParameterMismatchException>(create);
  static IdempotentParameterMismatchException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ImportConflictException extends $pb.GeneratedMessage {
  factory ImportConflictException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ImportConflictException._();

  factory ImportConflictException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportConflictException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportConflictException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportConflictException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportConflictException copyWith(
          void Function(ImportConflictException) updates) =>
      super.copyWith((message) => updates(message as ImportConflictException))
          as ImportConflictException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportConflictException create() => ImportConflictException._();
  @$core.override
  ImportConflictException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportConflictException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportConflictException>(create);
  static ImportConflictException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ImportNotFoundException extends $pb.GeneratedMessage {
  factory ImportNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ImportNotFoundException._();

  factory ImportNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportNotFoundException copyWith(
          void Function(ImportNotFoundException) updates) =>
      super.copyWith((message) => updates(message as ImportNotFoundException))
          as ImportNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportNotFoundException create() => ImportNotFoundException._();
  @$core.override
  ImportNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportNotFoundException>(create);
  static ImportNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ImportSummary extends $pb.GeneratedMessage {
  factory ImportSummary({
    $core.String? endtime,
    ImportStatus? importstatus,
    $core.String? cloudwatchloggrouparn,
    S3BucketSource? s3bucketsource,
    $core.String? starttime,
    InputFormat? inputformat,
    $core.String? tablearn,
    $core.String? importarn,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (importstatus != null) result.importstatus = importstatus;
    if (cloudwatchloggrouparn != null)
      result.cloudwatchloggrouparn = cloudwatchloggrouparn;
    if (s3bucketsource != null) result.s3bucketsource = s3bucketsource;
    if (starttime != null) result.starttime = starttime;
    if (inputformat != null) result.inputformat = inputformat;
    if (tablearn != null) result.tablearn = tablearn;
    if (importarn != null) result.importarn = importarn;
    return result;
  }

  ImportSummary._();

  factory ImportSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aE<ImportStatus>(129077631, _omitFieldNames ? '' : 'importstatus',
        enumValues: ImportStatus.values)
    ..aOS(171042008, _omitFieldNames ? '' : 'cloudwatchloggrouparn')
    ..aOM<S3BucketSource>(202310037, _omitFieldNames ? '' : 's3bucketsource',
        subBuilder: S3BucketSource.create)
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..aE<InputFormat>(405664101, _omitFieldNames ? '' : 'inputformat',
        enumValues: InputFormat.values)
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..aOS(444379628, _omitFieldNames ? '' : 'importarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportSummary copyWith(void Function(ImportSummary) updates) =>
      super.copyWith((message) => updates(message as ImportSummary))
          as ImportSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportSummary create() => ImportSummary._();
  @$core.override
  ImportSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportSummary>(create);
  static ImportSummary? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(129077631)
  ImportStatus get importstatus => $_getN(1);
  @$pb.TagNumber(129077631)
  set importstatus(ImportStatus value) => $_setField(129077631, value);
  @$pb.TagNumber(129077631)
  $core.bool hasImportstatus() => $_has(1);
  @$pb.TagNumber(129077631)
  void clearImportstatus() => $_clearField(129077631);

  @$pb.TagNumber(171042008)
  $core.String get cloudwatchloggrouparn => $_getSZ(2);
  @$pb.TagNumber(171042008)
  set cloudwatchloggrouparn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(171042008)
  $core.bool hasCloudwatchloggrouparn() => $_has(2);
  @$pb.TagNumber(171042008)
  void clearCloudwatchloggrouparn() => $_clearField(171042008);

  @$pb.TagNumber(202310037)
  S3BucketSource get s3bucketsource => $_getN(3);
  @$pb.TagNumber(202310037)
  set s3bucketsource(S3BucketSource value) => $_setField(202310037, value);
  @$pb.TagNumber(202310037)
  $core.bool hasS3bucketsource() => $_has(3);
  @$pb.TagNumber(202310037)
  void clearS3bucketsource() => $_clearField(202310037);
  @$pb.TagNumber(202310037)
  S3BucketSource ensureS3bucketsource() => $_ensure(3);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(4);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(4, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(4);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);

  @$pb.TagNumber(405664101)
  InputFormat get inputformat => $_getN(5);
  @$pb.TagNumber(405664101)
  set inputformat(InputFormat value) => $_setField(405664101, value);
  @$pb.TagNumber(405664101)
  $core.bool hasInputformat() => $_has(5);
  @$pb.TagNumber(405664101)
  void clearInputformat() => $_clearField(405664101);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(6);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(6);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);

  @$pb.TagNumber(444379628)
  $core.String get importarn => $_getSZ(7);
  @$pb.TagNumber(444379628)
  set importarn($core.String value) => $_setString(7, value);
  @$pb.TagNumber(444379628)
  $core.bool hasImportarn() => $_has(7);
  @$pb.TagNumber(444379628)
  void clearImportarn() => $_clearField(444379628);
}

class ImportTableDescription extends $pb.GeneratedMessage {
  factory ImportTableDescription({
    $core.String? endtime,
    $fixnum.Int64? processedsizebytes,
    $core.String? failurecode,
    ImportStatus? importstatus,
    $core.String? clienttoken,
    $core.String? cloudwatchloggrouparn,
    S3BucketSource? s3bucketsource,
    $fixnum.Int64? importeditemcount,
    InputFormatOptions? inputformatoptions,
    $fixnum.Int64? errorcount,
    $fixnum.Int64? processeditemcount,
    $core.String? failuremessage,
    $core.String? starttime,
    InputCompressionType? inputcompressiontype,
    InputFormat? inputformat,
    $core.String? tablearn,
    $core.String? importarn,
    $core.String? tableid,
    TableCreationParameters? tablecreationparameters,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (processedsizebytes != null)
      result.processedsizebytes = processedsizebytes;
    if (failurecode != null) result.failurecode = failurecode;
    if (importstatus != null) result.importstatus = importstatus;
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (cloudwatchloggrouparn != null)
      result.cloudwatchloggrouparn = cloudwatchloggrouparn;
    if (s3bucketsource != null) result.s3bucketsource = s3bucketsource;
    if (importeditemcount != null) result.importeditemcount = importeditemcount;
    if (inputformatoptions != null)
      result.inputformatoptions = inputformatoptions;
    if (errorcount != null) result.errorcount = errorcount;
    if (processeditemcount != null)
      result.processeditemcount = processeditemcount;
    if (failuremessage != null) result.failuremessage = failuremessage;
    if (starttime != null) result.starttime = starttime;
    if (inputcompressiontype != null)
      result.inputcompressiontype = inputcompressiontype;
    if (inputformat != null) result.inputformat = inputformat;
    if (tablearn != null) result.tablearn = tablearn;
    if (importarn != null) result.importarn = importarn;
    if (tableid != null) result.tableid = tableid;
    if (tablecreationparameters != null)
      result.tablecreationparameters = tablecreationparameters;
    return result;
  }

  ImportTableDescription._();

  factory ImportTableDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportTableDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportTableDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aInt64(72287562, _omitFieldNames ? '' : 'processedsizebytes')
    ..aOS(84707897, _omitFieldNames ? '' : 'failurecode')
    ..aE<ImportStatus>(129077631, _omitFieldNames ? '' : 'importstatus',
        enumValues: ImportStatus.values)
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOS(171042008, _omitFieldNames ? '' : 'cloudwatchloggrouparn')
    ..aOM<S3BucketSource>(202310037, _omitFieldNames ? '' : 's3bucketsource',
        subBuilder: S3BucketSource.create)
    ..aInt64(202622198, _omitFieldNames ? '' : 'importeditemcount')
    ..aOM<InputFormatOptions>(
        249201403, _omitFieldNames ? '' : 'inputformatoptions',
        subBuilder: InputFormatOptions.create)
    ..aInt64(311137001, _omitFieldNames ? '' : 'errorcount')
    ..aInt64(340163668, _omitFieldNames ? '' : 'processeditemcount')
    ..aOS(353556937, _omitFieldNames ? '' : 'failuremessage')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..aE<InputCompressionType>(
        392699396, _omitFieldNames ? '' : 'inputcompressiontype',
        enumValues: InputCompressionType.values)
    ..aE<InputFormat>(405664101, _omitFieldNames ? '' : 'inputformat',
        enumValues: InputFormat.values)
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..aOS(444379628, _omitFieldNames ? '' : 'importarn')
    ..aOS(449893011, _omitFieldNames ? '' : 'tableid')
    ..aOM<TableCreationParameters>(
        487419839, _omitFieldNames ? '' : 'tablecreationparameters',
        subBuilder: TableCreationParameters.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportTableDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportTableDescription copyWith(
          void Function(ImportTableDescription) updates) =>
      super.copyWith((message) => updates(message as ImportTableDescription))
          as ImportTableDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportTableDescription create() => ImportTableDescription._();
  @$core.override
  ImportTableDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportTableDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportTableDescription>(create);
  static ImportTableDescription? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(72287562)
  $fixnum.Int64 get processedsizebytes => $_getI64(1);
  @$pb.TagNumber(72287562)
  set processedsizebytes($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(72287562)
  $core.bool hasProcessedsizebytes() => $_has(1);
  @$pb.TagNumber(72287562)
  void clearProcessedsizebytes() => $_clearField(72287562);

  @$pb.TagNumber(84707897)
  $core.String get failurecode => $_getSZ(2);
  @$pb.TagNumber(84707897)
  set failurecode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(84707897)
  $core.bool hasFailurecode() => $_has(2);
  @$pb.TagNumber(84707897)
  void clearFailurecode() => $_clearField(84707897);

  @$pb.TagNumber(129077631)
  ImportStatus get importstatus => $_getN(3);
  @$pb.TagNumber(129077631)
  set importstatus(ImportStatus value) => $_setField(129077631, value);
  @$pb.TagNumber(129077631)
  $core.bool hasImportstatus() => $_has(3);
  @$pb.TagNumber(129077631)
  void clearImportstatus() => $_clearField(129077631);

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(4);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(4);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(171042008)
  $core.String get cloudwatchloggrouparn => $_getSZ(5);
  @$pb.TagNumber(171042008)
  set cloudwatchloggrouparn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(171042008)
  $core.bool hasCloudwatchloggrouparn() => $_has(5);
  @$pb.TagNumber(171042008)
  void clearCloudwatchloggrouparn() => $_clearField(171042008);

  @$pb.TagNumber(202310037)
  S3BucketSource get s3bucketsource => $_getN(6);
  @$pb.TagNumber(202310037)
  set s3bucketsource(S3BucketSource value) => $_setField(202310037, value);
  @$pb.TagNumber(202310037)
  $core.bool hasS3bucketsource() => $_has(6);
  @$pb.TagNumber(202310037)
  void clearS3bucketsource() => $_clearField(202310037);
  @$pb.TagNumber(202310037)
  S3BucketSource ensureS3bucketsource() => $_ensure(6);

  @$pb.TagNumber(202622198)
  $fixnum.Int64 get importeditemcount => $_getI64(7);
  @$pb.TagNumber(202622198)
  set importeditemcount($fixnum.Int64 value) => $_setInt64(7, value);
  @$pb.TagNumber(202622198)
  $core.bool hasImporteditemcount() => $_has(7);
  @$pb.TagNumber(202622198)
  void clearImporteditemcount() => $_clearField(202622198);

  @$pb.TagNumber(249201403)
  InputFormatOptions get inputformatoptions => $_getN(8);
  @$pb.TagNumber(249201403)
  set inputformatoptions(InputFormatOptions value) =>
      $_setField(249201403, value);
  @$pb.TagNumber(249201403)
  $core.bool hasInputformatoptions() => $_has(8);
  @$pb.TagNumber(249201403)
  void clearInputformatoptions() => $_clearField(249201403);
  @$pb.TagNumber(249201403)
  InputFormatOptions ensureInputformatoptions() => $_ensure(8);

  @$pb.TagNumber(311137001)
  $fixnum.Int64 get errorcount => $_getI64(9);
  @$pb.TagNumber(311137001)
  set errorcount($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(311137001)
  $core.bool hasErrorcount() => $_has(9);
  @$pb.TagNumber(311137001)
  void clearErrorcount() => $_clearField(311137001);

  @$pb.TagNumber(340163668)
  $fixnum.Int64 get processeditemcount => $_getI64(10);
  @$pb.TagNumber(340163668)
  set processeditemcount($fixnum.Int64 value) => $_setInt64(10, value);
  @$pb.TagNumber(340163668)
  $core.bool hasProcesseditemcount() => $_has(10);
  @$pb.TagNumber(340163668)
  void clearProcesseditemcount() => $_clearField(340163668);

  @$pb.TagNumber(353556937)
  $core.String get failuremessage => $_getSZ(11);
  @$pb.TagNumber(353556937)
  set failuremessage($core.String value) => $_setString(11, value);
  @$pb.TagNumber(353556937)
  $core.bool hasFailuremessage() => $_has(11);
  @$pb.TagNumber(353556937)
  void clearFailuremessage() => $_clearField(353556937);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(12);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(12, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(12);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);

  @$pb.TagNumber(392699396)
  InputCompressionType get inputcompressiontype => $_getN(13);
  @$pb.TagNumber(392699396)
  set inputcompressiontype(InputCompressionType value) =>
      $_setField(392699396, value);
  @$pb.TagNumber(392699396)
  $core.bool hasInputcompressiontype() => $_has(13);
  @$pb.TagNumber(392699396)
  void clearInputcompressiontype() => $_clearField(392699396);

  @$pb.TagNumber(405664101)
  InputFormat get inputformat => $_getN(14);
  @$pb.TagNumber(405664101)
  set inputformat(InputFormat value) => $_setField(405664101, value);
  @$pb.TagNumber(405664101)
  $core.bool hasInputformat() => $_has(14);
  @$pb.TagNumber(405664101)
  void clearInputformat() => $_clearField(405664101);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(15);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(15, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(15);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);

  @$pb.TagNumber(444379628)
  $core.String get importarn => $_getSZ(16);
  @$pb.TagNumber(444379628)
  set importarn($core.String value) => $_setString(16, value);
  @$pb.TagNumber(444379628)
  $core.bool hasImportarn() => $_has(16);
  @$pb.TagNumber(444379628)
  void clearImportarn() => $_clearField(444379628);

  @$pb.TagNumber(449893011)
  $core.String get tableid => $_getSZ(17);
  @$pb.TagNumber(449893011)
  set tableid($core.String value) => $_setString(17, value);
  @$pb.TagNumber(449893011)
  $core.bool hasTableid() => $_has(17);
  @$pb.TagNumber(449893011)
  void clearTableid() => $_clearField(449893011);

  @$pb.TagNumber(487419839)
  TableCreationParameters get tablecreationparameters => $_getN(18);
  @$pb.TagNumber(487419839)
  set tablecreationparameters(TableCreationParameters value) =>
      $_setField(487419839, value);
  @$pb.TagNumber(487419839)
  $core.bool hasTablecreationparameters() => $_has(18);
  @$pb.TagNumber(487419839)
  void clearTablecreationparameters() => $_clearField(487419839);
  @$pb.TagNumber(487419839)
  TableCreationParameters ensureTablecreationparameters() => $_ensure(18);
}

class ImportTableInput extends $pb.GeneratedMessage {
  factory ImportTableInput({
    $core.String? clienttoken,
    S3BucketSource? s3bucketsource,
    InputFormatOptions? inputformatoptions,
    InputCompressionType? inputcompressiontype,
    InputFormat? inputformat,
    TableCreationParameters? tablecreationparameters,
  }) {
    final result = create();
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (s3bucketsource != null) result.s3bucketsource = s3bucketsource;
    if (inputformatoptions != null)
      result.inputformatoptions = inputformatoptions;
    if (inputcompressiontype != null)
      result.inputcompressiontype = inputcompressiontype;
    if (inputformat != null) result.inputformat = inputformat;
    if (tablecreationparameters != null)
      result.tablecreationparameters = tablecreationparameters;
    return result;
  }

  ImportTableInput._();

  factory ImportTableInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportTableInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportTableInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOM<S3BucketSource>(202310037, _omitFieldNames ? '' : 's3bucketsource',
        subBuilder: S3BucketSource.create)
    ..aOM<InputFormatOptions>(
        249201403, _omitFieldNames ? '' : 'inputformatoptions',
        subBuilder: InputFormatOptions.create)
    ..aE<InputCompressionType>(
        392699396, _omitFieldNames ? '' : 'inputcompressiontype',
        enumValues: InputCompressionType.values)
    ..aE<InputFormat>(405664101, _omitFieldNames ? '' : 'inputformat',
        enumValues: InputFormat.values)
    ..aOM<TableCreationParameters>(
        487419839, _omitFieldNames ? '' : 'tablecreationparameters',
        subBuilder: TableCreationParameters.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportTableInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportTableInput copyWith(void Function(ImportTableInput) updates) =>
      super.copyWith((message) => updates(message as ImportTableInput))
          as ImportTableInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportTableInput create() => ImportTableInput._();
  @$core.override
  ImportTableInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportTableInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportTableInput>(create);
  static ImportTableInput? _defaultInstance;

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(0);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(0);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(202310037)
  S3BucketSource get s3bucketsource => $_getN(1);
  @$pb.TagNumber(202310037)
  set s3bucketsource(S3BucketSource value) => $_setField(202310037, value);
  @$pb.TagNumber(202310037)
  $core.bool hasS3bucketsource() => $_has(1);
  @$pb.TagNumber(202310037)
  void clearS3bucketsource() => $_clearField(202310037);
  @$pb.TagNumber(202310037)
  S3BucketSource ensureS3bucketsource() => $_ensure(1);

  @$pb.TagNumber(249201403)
  InputFormatOptions get inputformatoptions => $_getN(2);
  @$pb.TagNumber(249201403)
  set inputformatoptions(InputFormatOptions value) =>
      $_setField(249201403, value);
  @$pb.TagNumber(249201403)
  $core.bool hasInputformatoptions() => $_has(2);
  @$pb.TagNumber(249201403)
  void clearInputformatoptions() => $_clearField(249201403);
  @$pb.TagNumber(249201403)
  InputFormatOptions ensureInputformatoptions() => $_ensure(2);

  @$pb.TagNumber(392699396)
  InputCompressionType get inputcompressiontype => $_getN(3);
  @$pb.TagNumber(392699396)
  set inputcompressiontype(InputCompressionType value) =>
      $_setField(392699396, value);
  @$pb.TagNumber(392699396)
  $core.bool hasInputcompressiontype() => $_has(3);
  @$pb.TagNumber(392699396)
  void clearInputcompressiontype() => $_clearField(392699396);

  @$pb.TagNumber(405664101)
  InputFormat get inputformat => $_getN(4);
  @$pb.TagNumber(405664101)
  set inputformat(InputFormat value) => $_setField(405664101, value);
  @$pb.TagNumber(405664101)
  $core.bool hasInputformat() => $_has(4);
  @$pb.TagNumber(405664101)
  void clearInputformat() => $_clearField(405664101);

  @$pb.TagNumber(487419839)
  TableCreationParameters get tablecreationparameters => $_getN(5);
  @$pb.TagNumber(487419839)
  set tablecreationparameters(TableCreationParameters value) =>
      $_setField(487419839, value);
  @$pb.TagNumber(487419839)
  $core.bool hasTablecreationparameters() => $_has(5);
  @$pb.TagNumber(487419839)
  void clearTablecreationparameters() => $_clearField(487419839);
  @$pb.TagNumber(487419839)
  TableCreationParameters ensureTablecreationparameters() => $_ensure(5);
}

class ImportTableOutput extends $pb.GeneratedMessage {
  factory ImportTableOutput({
    ImportTableDescription? importtabledescription,
  }) {
    final result = create();
    if (importtabledescription != null)
      result.importtabledescription = importtabledescription;
    return result;
  }

  ImportTableOutput._();

  factory ImportTableOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ImportTableOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ImportTableOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ImportTableDescription>(
        407746467, _omitFieldNames ? '' : 'importtabledescription',
        subBuilder: ImportTableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportTableOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ImportTableOutput copyWith(void Function(ImportTableOutput) updates) =>
      super.copyWith((message) => updates(message as ImportTableOutput))
          as ImportTableOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ImportTableOutput create() => ImportTableOutput._();
  @$core.override
  ImportTableOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ImportTableOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ImportTableOutput>(create);
  static ImportTableOutput? _defaultInstance;

  @$pb.TagNumber(407746467)
  ImportTableDescription get importtabledescription => $_getN(0);
  @$pb.TagNumber(407746467)
  set importtabledescription(ImportTableDescription value) =>
      $_setField(407746467, value);
  @$pb.TagNumber(407746467)
  $core.bool hasImporttabledescription() => $_has(0);
  @$pb.TagNumber(407746467)
  void clearImporttabledescription() => $_clearField(407746467);
  @$pb.TagNumber(407746467)
  ImportTableDescription ensureImporttabledescription() => $_ensure(0);
}

class IncrementalExportSpecification extends $pb.GeneratedMessage {
  factory IncrementalExportSpecification({
    ExportViewType? exportviewtype,
    $core.String? exporttotime,
    $core.String? exportfromtime,
  }) {
    final result = create();
    if (exportviewtype != null) result.exportviewtype = exportviewtype;
    if (exporttotime != null) result.exporttotime = exporttotime;
    if (exportfromtime != null) result.exportfromtime = exportfromtime;
    return result;
  }

  IncrementalExportSpecification._();

  factory IncrementalExportSpecification.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IncrementalExportSpecification.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IncrementalExportSpecification',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ExportViewType>(26765639, _omitFieldNames ? '' : 'exportviewtype',
        enumValues: ExportViewType.values)
    ..aOS(242369368, _omitFieldNames ? '' : 'exporttotime')
    ..aOS(466248415, _omitFieldNames ? '' : 'exportfromtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncrementalExportSpecification clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncrementalExportSpecification copyWith(
          void Function(IncrementalExportSpecification) updates) =>
      super.copyWith(
              (message) => updates(message as IncrementalExportSpecification))
          as IncrementalExportSpecification;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IncrementalExportSpecification create() =>
      IncrementalExportSpecification._();
  @$core.override
  IncrementalExportSpecification createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IncrementalExportSpecification getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IncrementalExportSpecification>(create);
  static IncrementalExportSpecification? _defaultInstance;

  @$pb.TagNumber(26765639)
  ExportViewType get exportviewtype => $_getN(0);
  @$pb.TagNumber(26765639)
  set exportviewtype(ExportViewType value) => $_setField(26765639, value);
  @$pb.TagNumber(26765639)
  $core.bool hasExportviewtype() => $_has(0);
  @$pb.TagNumber(26765639)
  void clearExportviewtype() => $_clearField(26765639);

  @$pb.TagNumber(242369368)
  $core.String get exporttotime => $_getSZ(1);
  @$pb.TagNumber(242369368)
  set exporttotime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(242369368)
  $core.bool hasExporttotime() => $_has(1);
  @$pb.TagNumber(242369368)
  void clearExporttotime() => $_clearField(242369368);

  @$pb.TagNumber(466248415)
  $core.String get exportfromtime => $_getSZ(2);
  @$pb.TagNumber(466248415)
  set exportfromtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(466248415)
  $core.bool hasExportfromtime() => $_has(2);
  @$pb.TagNumber(466248415)
  void clearExportfromtime() => $_clearField(466248415);
}

class IndexNotFoundException extends $pb.GeneratedMessage {
  factory IndexNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IndexNotFoundException._();

  factory IndexNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IndexNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IndexNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IndexNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IndexNotFoundException copyWith(
          void Function(IndexNotFoundException) updates) =>
      super.copyWith((message) => updates(message as IndexNotFoundException))
          as IndexNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IndexNotFoundException create() => IndexNotFoundException._();
  @$core.override
  IndexNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IndexNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IndexNotFoundException>(create);
  static IndexNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InputFormatOptions extends $pb.GeneratedMessage {
  factory InputFormatOptions({
    CsvOptions? csv,
  }) {
    final result = create();
    if (csv != null) result.csv = csv;
    return result;
  }

  InputFormatOptions._();

  factory InputFormatOptions.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InputFormatOptions.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InputFormatOptions',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<CsvOptions>(450056448, _omitFieldNames ? '' : 'csv',
        subBuilder: CsvOptions.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InputFormatOptions clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InputFormatOptions copyWith(void Function(InputFormatOptions) updates) =>
      super.copyWith((message) => updates(message as InputFormatOptions))
          as InputFormatOptions;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InputFormatOptions create() => InputFormatOptions._();
  @$core.override
  InputFormatOptions createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InputFormatOptions getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InputFormatOptions>(create);
  static InputFormatOptions? _defaultInstance;

  @$pb.TagNumber(450056448)
  CsvOptions get csv => $_getN(0);
  @$pb.TagNumber(450056448)
  set csv(CsvOptions value) => $_setField(450056448, value);
  @$pb.TagNumber(450056448)
  $core.bool hasCsv() => $_has(0);
  @$pb.TagNumber(450056448)
  void clearCsv() => $_clearField(450056448);
  @$pb.TagNumber(450056448)
  CsvOptions ensureCsv() => $_ensure(0);
}

class InternalServerError extends $pb.GeneratedMessage {
  factory InternalServerError({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InternalServerError._();

  factory InternalServerError.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InternalServerError.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InternalServerError',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServerError clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServerError copyWith(void Function(InternalServerError) updates) =>
      super.copyWith((message) => updates(message as InternalServerError))
          as InternalServerError;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InternalServerError create() => InternalServerError._();
  @$core.override
  InternalServerError createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InternalServerError getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InternalServerError>(create);
  static InternalServerError? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
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

class InvalidExportTimeException extends $pb.GeneratedMessage {
  factory InvalidExportTimeException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidExportTimeException._();

  factory InvalidExportTimeException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidExportTimeException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidExportTimeException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidExportTimeException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidExportTimeException copyWith(
          void Function(InvalidExportTimeException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidExportTimeException))
          as InvalidExportTimeException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidExportTimeException create() => InvalidExportTimeException._();
  @$core.override
  InvalidExportTimeException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidExportTimeException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidExportTimeException>(create);
  static InvalidExportTimeException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidRestoreTimeException extends $pb.GeneratedMessage {
  factory InvalidRestoreTimeException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidRestoreTimeException._();

  factory InvalidRestoreTimeException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidRestoreTimeException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidRestoreTimeException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidRestoreTimeException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidRestoreTimeException copyWith(
          void Function(InvalidRestoreTimeException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidRestoreTimeException))
          as InvalidRestoreTimeException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidRestoreTimeException create() =>
      InvalidRestoreTimeException._();
  @$core.override
  InvalidRestoreTimeException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidRestoreTimeException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidRestoreTimeException>(create);
  static InvalidRestoreTimeException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ItemCollectionMetrics extends $pb.GeneratedMessage {
  factory ItemCollectionMetrics({
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        itemcollectionkey,
    $core.Iterable<$core.double>? sizeestimaterangegb,
  }) {
    final result = create();
    if (itemcollectionkey != null)
      result.itemcollectionkey.addEntries(itemcollectionkey);
    if (sizeestimaterangegb != null)
      result.sizeestimaterangegb.addAll(sizeestimaterangegb);
    return result;
  }

  ItemCollectionMetrics._();

  factory ItemCollectionMetrics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ItemCollectionMetrics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ItemCollectionMetrics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, AttributeValue>(
        6645926, _omitFieldNames ? '' : 'itemcollectionkey',
        entryClassName: 'ItemCollectionMetrics.ItemcollectionkeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..p<$core.double>(271707895, _omitFieldNames ? '' : 'sizeestimaterangegb',
        $pb.PbFieldType.KD)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ItemCollectionMetrics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ItemCollectionMetrics copyWith(
          void Function(ItemCollectionMetrics) updates) =>
      super.copyWith((message) => updates(message as ItemCollectionMetrics))
          as ItemCollectionMetrics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ItemCollectionMetrics create() => ItemCollectionMetrics._();
  @$core.override
  ItemCollectionMetrics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ItemCollectionMetrics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ItemCollectionMetrics>(create);
  static ItemCollectionMetrics? _defaultInstance;

  @$pb.TagNumber(6645926)
  $pb.PbMap<$core.String, AttributeValue> get itemcollectionkey => $_getMap(0);

  @$pb.TagNumber(271707895)
  $pb.PbList<$core.double> get sizeestimaterangegb => $_getList(1);
}

class ItemCollectionSizeLimitExceededException extends $pb.GeneratedMessage {
  factory ItemCollectionSizeLimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ItemCollectionSizeLimitExceededException._();

  factory ItemCollectionSizeLimitExceededException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ItemCollectionSizeLimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ItemCollectionSizeLimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ItemCollectionSizeLimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ItemCollectionSizeLimitExceededException copyWith(
          void Function(ItemCollectionSizeLimitExceededException) updates) =>
      super.copyWith((message) =>
              updates(message as ItemCollectionSizeLimitExceededException))
          as ItemCollectionSizeLimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ItemCollectionSizeLimitExceededException create() =>
      ItemCollectionSizeLimitExceededException._();
  @$core.override
  ItemCollectionSizeLimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ItemCollectionSizeLimitExceededException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ItemCollectionSizeLimitExceededException>(create);
  static ItemCollectionSizeLimitExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ItemResponse extends $pb.GeneratedMessage {
  factory ItemResponse({
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (item != null) result.item.addEntries(item);
    return result;
  }

  ItemResponse._();

  factory ItemResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ItemResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ItemResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'ItemResponse.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ItemResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ItemResponse copyWith(void Function(ItemResponse) updates) =>
      super.copyWith((message) => updates(message as ItemResponse))
          as ItemResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ItemResponse create() => ItemResponse._();
  @$core.override
  ItemResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ItemResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ItemResponse>(create);
  static ItemResponse? _defaultInstance;

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(0);
}

class KeySchemaElement extends $pb.GeneratedMessage {
  factory KeySchemaElement({
    KeyType? keytype,
    $core.String? attributename,
  }) {
    final result = create();
    if (keytype != null) result.keytype = keytype;
    if (attributename != null) result.attributename = attributename;
    return result;
  }

  KeySchemaElement._();

  factory KeySchemaElement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeySchemaElement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeySchemaElement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<KeyType>(5029221, _omitFieldNames ? '' : 'keytype',
        enumValues: KeyType.values)
    ..aOS(352717485, _omitFieldNames ? '' : 'attributename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySchemaElement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySchemaElement copyWith(void Function(KeySchemaElement) updates) =>
      super.copyWith((message) => updates(message as KeySchemaElement))
          as KeySchemaElement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeySchemaElement create() => KeySchemaElement._();
  @$core.override
  KeySchemaElement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeySchemaElement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeySchemaElement>(create);
  static KeySchemaElement? _defaultInstance;

  @$pb.TagNumber(5029221)
  KeyType get keytype => $_getN(0);
  @$pb.TagNumber(5029221)
  set keytype(KeyType value) => $_setField(5029221, value);
  @$pb.TagNumber(5029221)
  $core.bool hasKeytype() => $_has(0);
  @$pb.TagNumber(5029221)
  void clearKeytype() => $_clearField(5029221);

  @$pb.TagNumber(352717485)
  $core.String get attributename => $_getSZ(1);
  @$pb.TagNumber(352717485)
  set attributename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(352717485)
  $core.bool hasAttributename() => $_has(1);
  @$pb.TagNumber(352717485)
  void clearAttributename() => $_clearField(352717485);
}

class KeysAndAttributes extends $pb.GeneratedMessage {
  factory KeysAndAttributes({
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? keys,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.String? projectionexpression,
    $core.Iterable<$core.String>? attributestoget,
    $core.bool? consistentread,
  }) {
    final result = create();
    if (keys != null) result.keys.addEntries(keys);
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (projectionexpression != null)
      result.projectionexpression = projectionexpression;
    if (attributestoget != null) result.attributestoget.addAll(attributestoget);
    if (consistentread != null) result.consistentread = consistentread;
    return result;
  }

  KeysAndAttributes._();

  factory KeysAndAttributes.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeysAndAttributes.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeysAndAttributes',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, AttributeValue>(2831086, _omitFieldNames ? '' : 'keys',
        entryClassName: 'KeysAndAttributes.KeysEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'KeysAndAttributes.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(150730243, _omitFieldNames ? '' : 'projectionexpression')
    ..pPS(311382592, _omitFieldNames ? '' : 'attributestoget')
    ..aOB(531556994, _omitFieldNames ? '' : 'consistentread')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeysAndAttributes clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeysAndAttributes copyWith(void Function(KeysAndAttributes) updates) =>
      super.copyWith((message) => updates(message as KeysAndAttributes))
          as KeysAndAttributes;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeysAndAttributes create() => KeysAndAttributes._();
  @$core.override
  KeysAndAttributes createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeysAndAttributes getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeysAndAttributes>(create);
  static KeysAndAttributes? _defaultInstance;

  @$pb.TagNumber(2831086)
  $pb.PbMap<$core.String, AttributeValue> get keys => $_getMap(0);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(1);

  @$pb.TagNumber(150730243)
  $core.String get projectionexpression => $_getSZ(2);
  @$pb.TagNumber(150730243)
  set projectionexpression($core.String value) => $_setString(2, value);
  @$pb.TagNumber(150730243)
  $core.bool hasProjectionexpression() => $_has(2);
  @$pb.TagNumber(150730243)
  void clearProjectionexpression() => $_clearField(150730243);

  @$pb.TagNumber(311382592)
  $pb.PbList<$core.String> get attributestoget => $_getList(3);

  @$pb.TagNumber(531556994)
  $core.bool get consistentread => $_getBF(4);
  @$pb.TagNumber(531556994)
  set consistentread($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(531556994)
  $core.bool hasConsistentread() => $_has(4);
  @$pb.TagNumber(531556994)
  void clearConsistentread() => $_clearField(531556994);
}

class KinesisDataStreamDestination extends $pb.GeneratedMessage {
  factory KinesisDataStreamDestination({
    DestinationStatus? destinationstatus,
    ApproximateCreationDateTimePrecision? approximatecreationdatetimeprecision,
    $core.String? destinationstatusdescription,
    $core.String? streamarn,
  }) {
    final result = create();
    if (destinationstatus != null) result.destinationstatus = destinationstatus;
    if (approximatecreationdatetimeprecision != null)
      result.approximatecreationdatetimeprecision =
          approximatecreationdatetimeprecision;
    if (destinationstatusdescription != null)
      result.destinationstatusdescription = destinationstatusdescription;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  KinesisDataStreamDestination._();

  factory KinesisDataStreamDestination.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KinesisDataStreamDestination.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KinesisDataStreamDestination',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<DestinationStatus>(
        381248234, _omitFieldNames ? '' : 'destinationstatus',
        enumValues: DestinationStatus.values)
    ..aE<ApproximateCreationDateTimePrecision>(392293352,
        _omitFieldNames ? '' : 'approximatecreationdatetimeprecision',
        enumValues: ApproximateCreationDateTimePrecision.values)
    ..aOS(499573086, _omitFieldNames ? '' : 'destinationstatusdescription')
    ..aOS(513423709, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KinesisDataStreamDestination clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KinesisDataStreamDestination copyWith(
          void Function(KinesisDataStreamDestination) updates) =>
      super.copyWith(
              (message) => updates(message as KinesisDataStreamDestination))
          as KinesisDataStreamDestination;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KinesisDataStreamDestination create() =>
      KinesisDataStreamDestination._();
  @$core.override
  KinesisDataStreamDestination createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KinesisDataStreamDestination getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KinesisDataStreamDestination>(create);
  static KinesisDataStreamDestination? _defaultInstance;

  @$pb.TagNumber(381248234)
  DestinationStatus get destinationstatus => $_getN(0);
  @$pb.TagNumber(381248234)
  set destinationstatus(DestinationStatus value) =>
      $_setField(381248234, value);
  @$pb.TagNumber(381248234)
  $core.bool hasDestinationstatus() => $_has(0);
  @$pb.TagNumber(381248234)
  void clearDestinationstatus() => $_clearField(381248234);

  @$pb.TagNumber(392293352)
  ApproximateCreationDateTimePrecision
      get approximatecreationdatetimeprecision => $_getN(1);
  @$pb.TagNumber(392293352)
  set approximatecreationdatetimeprecision(
          ApproximateCreationDateTimePrecision value) =>
      $_setField(392293352, value);
  @$pb.TagNumber(392293352)
  $core.bool hasApproximatecreationdatetimeprecision() => $_has(1);
  @$pb.TagNumber(392293352)
  void clearApproximatecreationdatetimeprecision() => $_clearField(392293352);

  @$pb.TagNumber(499573086)
  $core.String get destinationstatusdescription => $_getSZ(2);
  @$pb.TagNumber(499573086)
  set destinationstatusdescription($core.String value) => $_setString(2, value);
  @$pb.TagNumber(499573086)
  $core.bool hasDestinationstatusdescription() => $_has(2);
  @$pb.TagNumber(499573086)
  void clearDestinationstatusdescription() => $_clearField(499573086);

  @$pb.TagNumber(513423709)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(513423709)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(513423709)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(513423709)
  void clearStreamarn() => $_clearField(513423709);
}

class KinesisStreamingDestinationInput extends $pb.GeneratedMessage {
  factory KinesisStreamingDestinationInput({
    $core.String? tablename,
    EnableKinesisStreamingConfiguration? enablekinesisstreamingconfiguration,
    $core.String? streamarn,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    if (enablekinesisstreamingconfiguration != null)
      result.enablekinesisstreamingconfiguration =
          enablekinesisstreamingconfiguration;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  KinesisStreamingDestinationInput._();

  factory KinesisStreamingDestinationInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KinesisStreamingDestinationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KinesisStreamingDestinationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<EnableKinesisStreamingConfiguration>(
        300307843, _omitFieldNames ? '' : 'enablekinesisstreamingconfiguration',
        subBuilder: EnableKinesisStreamingConfiguration.create)
    ..aOS(513423709, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KinesisStreamingDestinationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KinesisStreamingDestinationInput copyWith(
          void Function(KinesisStreamingDestinationInput) updates) =>
      super.copyWith(
              (message) => updates(message as KinesisStreamingDestinationInput))
          as KinesisStreamingDestinationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KinesisStreamingDestinationInput create() =>
      KinesisStreamingDestinationInput._();
  @$core.override
  KinesisStreamingDestinationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KinesisStreamingDestinationInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KinesisStreamingDestinationInput>(
          create);
  static KinesisStreamingDestinationInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(300307843)
  EnableKinesisStreamingConfiguration get enablekinesisstreamingconfiguration =>
      $_getN(1);
  @$pb.TagNumber(300307843)
  set enablekinesisstreamingconfiguration(
          EnableKinesisStreamingConfiguration value) =>
      $_setField(300307843, value);
  @$pb.TagNumber(300307843)
  $core.bool hasEnablekinesisstreamingconfiguration() => $_has(1);
  @$pb.TagNumber(300307843)
  void clearEnablekinesisstreamingconfiguration() => $_clearField(300307843);
  @$pb.TagNumber(300307843)
  EnableKinesisStreamingConfiguration
      ensureEnablekinesisstreamingconfiguration() => $_ensure(1);

  @$pb.TagNumber(513423709)
  $core.String get streamarn => $_getSZ(2);
  @$pb.TagNumber(513423709)
  set streamarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(513423709)
  $core.bool hasStreamarn() => $_has(2);
  @$pb.TagNumber(513423709)
  void clearStreamarn() => $_clearField(513423709);
}

class KinesisStreamingDestinationOutput extends $pb.GeneratedMessage {
  factory KinesisStreamingDestinationOutput({
    $core.String? tablename,
    EnableKinesisStreamingConfiguration? enablekinesisstreamingconfiguration,
    DestinationStatus? destinationstatus,
    $core.String? streamarn,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    if (enablekinesisstreamingconfiguration != null)
      result.enablekinesisstreamingconfiguration =
          enablekinesisstreamingconfiguration;
    if (destinationstatus != null) result.destinationstatus = destinationstatus;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  KinesisStreamingDestinationOutput._();

  factory KinesisStreamingDestinationOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KinesisStreamingDestinationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KinesisStreamingDestinationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<EnableKinesisStreamingConfiguration>(
        300307843, _omitFieldNames ? '' : 'enablekinesisstreamingconfiguration',
        subBuilder: EnableKinesisStreamingConfiguration.create)
    ..aE<DestinationStatus>(
        381248234, _omitFieldNames ? '' : 'destinationstatus',
        enumValues: DestinationStatus.values)
    ..aOS(513423709, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KinesisStreamingDestinationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KinesisStreamingDestinationOutput copyWith(
          void Function(KinesisStreamingDestinationOutput) updates) =>
      super.copyWith((message) =>
              updates(message as KinesisStreamingDestinationOutput))
          as KinesisStreamingDestinationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KinesisStreamingDestinationOutput create() =>
      KinesisStreamingDestinationOutput._();
  @$core.override
  KinesisStreamingDestinationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KinesisStreamingDestinationOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KinesisStreamingDestinationOutput>(
          create);
  static KinesisStreamingDestinationOutput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(300307843)
  EnableKinesisStreamingConfiguration get enablekinesisstreamingconfiguration =>
      $_getN(1);
  @$pb.TagNumber(300307843)
  set enablekinesisstreamingconfiguration(
          EnableKinesisStreamingConfiguration value) =>
      $_setField(300307843, value);
  @$pb.TagNumber(300307843)
  $core.bool hasEnablekinesisstreamingconfiguration() => $_has(1);
  @$pb.TagNumber(300307843)
  void clearEnablekinesisstreamingconfiguration() => $_clearField(300307843);
  @$pb.TagNumber(300307843)
  EnableKinesisStreamingConfiguration
      ensureEnablekinesisstreamingconfiguration() => $_ensure(1);

  @$pb.TagNumber(381248234)
  DestinationStatus get destinationstatus => $_getN(2);
  @$pb.TagNumber(381248234)
  set destinationstatus(DestinationStatus value) =>
      $_setField(381248234, value);
  @$pb.TagNumber(381248234)
  $core.bool hasDestinationstatus() => $_has(2);
  @$pb.TagNumber(381248234)
  void clearDestinationstatus() => $_clearField(381248234);

  @$pb.TagNumber(513423709)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(513423709)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(513423709)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(513423709)
  void clearStreamarn() => $_clearField(513423709);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
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

class ListBackupsInput extends $pb.GeneratedMessage {
  factory ListBackupsInput({
    BackupTypeFilter? backuptype,
    $core.String? tablename,
    $core.String? exclusivestartbackuparn,
    $core.int? limit,
    $core.String? timerangeupperbound,
    $core.String? timerangelowerbound,
  }) {
    final result = create();
    if (backuptype != null) result.backuptype = backuptype;
    if (tablename != null) result.tablename = tablename;
    if (exclusivestartbackuparn != null)
      result.exclusivestartbackuparn = exclusivestartbackuparn;
    if (limit != null) result.limit = limit;
    if (timerangeupperbound != null)
      result.timerangeupperbound = timerangeupperbound;
    if (timerangelowerbound != null)
      result.timerangelowerbound = timerangelowerbound;
    return result;
  }

  ListBackupsInput._();

  factory ListBackupsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListBackupsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListBackupsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<BackupTypeFilter>(134973992, _omitFieldNames ? '' : 'backuptype',
        enumValues: BackupTypeFilter.values)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(381696959, _omitFieldNames ? '' : 'exclusivestartbackuparn')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(420889434, _omitFieldNames ? '' : 'timerangeupperbound')
    ..aOS(506236543, _omitFieldNames ? '' : 'timerangelowerbound')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListBackupsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListBackupsInput copyWith(void Function(ListBackupsInput) updates) =>
      super.copyWith((message) => updates(message as ListBackupsInput))
          as ListBackupsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListBackupsInput create() => ListBackupsInput._();
  @$core.override
  ListBackupsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListBackupsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListBackupsInput>(create);
  static ListBackupsInput? _defaultInstance;

  @$pb.TagNumber(134973992)
  BackupTypeFilter get backuptype => $_getN(0);
  @$pb.TagNumber(134973992)
  set backuptype(BackupTypeFilter value) => $_setField(134973992, value);
  @$pb.TagNumber(134973992)
  $core.bool hasBackuptype() => $_has(0);
  @$pb.TagNumber(134973992)
  void clearBackuptype() => $_clearField(134973992);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(381696959)
  $core.String get exclusivestartbackuparn => $_getSZ(2);
  @$pb.TagNumber(381696959)
  set exclusivestartbackuparn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(381696959)
  $core.bool hasExclusivestartbackuparn() => $_has(2);
  @$pb.TagNumber(381696959)
  void clearExclusivestartbackuparn() => $_clearField(381696959);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(420889434)
  $core.String get timerangeupperbound => $_getSZ(4);
  @$pb.TagNumber(420889434)
  set timerangeupperbound($core.String value) => $_setString(4, value);
  @$pb.TagNumber(420889434)
  $core.bool hasTimerangeupperbound() => $_has(4);
  @$pb.TagNumber(420889434)
  void clearTimerangeupperbound() => $_clearField(420889434);

  @$pb.TagNumber(506236543)
  $core.String get timerangelowerbound => $_getSZ(5);
  @$pb.TagNumber(506236543)
  set timerangelowerbound($core.String value) => $_setString(5, value);
  @$pb.TagNumber(506236543)
  $core.bool hasTimerangelowerbound() => $_has(5);
  @$pb.TagNumber(506236543)
  void clearTimerangelowerbound() => $_clearField(506236543);
}

class ListBackupsOutput extends $pb.GeneratedMessage {
  factory ListBackupsOutput({
    $core.String? lastevaluatedbackuparn,
    $core.Iterable<BackupSummary>? backupsummaries,
  }) {
    final result = create();
    if (lastevaluatedbackuparn != null)
      result.lastevaluatedbackuparn = lastevaluatedbackuparn;
    if (backupsummaries != null) result.backupsummaries.addAll(backupsummaries);
    return result;
  }

  ListBackupsOutput._();

  factory ListBackupsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListBackupsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListBackupsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(254744240, _omitFieldNames ? '' : 'lastevaluatedbackuparn')
    ..pPM<BackupSummary>(375274198, _omitFieldNames ? '' : 'backupsummaries',
        subBuilder: BackupSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListBackupsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListBackupsOutput copyWith(void Function(ListBackupsOutput) updates) =>
      super.copyWith((message) => updates(message as ListBackupsOutput))
          as ListBackupsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListBackupsOutput create() => ListBackupsOutput._();
  @$core.override
  ListBackupsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListBackupsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListBackupsOutput>(create);
  static ListBackupsOutput? _defaultInstance;

  @$pb.TagNumber(254744240)
  $core.String get lastevaluatedbackuparn => $_getSZ(0);
  @$pb.TagNumber(254744240)
  set lastevaluatedbackuparn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(254744240)
  $core.bool hasLastevaluatedbackuparn() => $_has(0);
  @$pb.TagNumber(254744240)
  void clearLastevaluatedbackuparn() => $_clearField(254744240);

  @$pb.TagNumber(375274198)
  $pb.PbList<BackupSummary> get backupsummaries => $_getList(1);
}

class ListContributorInsightsInput extends $pb.GeneratedMessage {
  factory ListContributorInsightsInput({
    $core.String? nexttoken,
    $core.String? tablename,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (tablename != null) result.tablename = tablename;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListContributorInsightsInput._();

  factory ListContributorInsightsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListContributorInsightsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListContributorInsightsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListContributorInsightsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListContributorInsightsInput copyWith(
          void Function(ListContributorInsightsInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListContributorInsightsInput))
          as ListContributorInsightsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListContributorInsightsInput create() =>
      ListContributorInsightsInput._();
  @$core.override
  ListContributorInsightsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListContributorInsightsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListContributorInsightsInput>(create);
  static ListContributorInsightsInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(2);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListContributorInsightsOutput extends $pb.GeneratedMessage {
  factory ListContributorInsightsOutput({
    $core.String? nexttoken,
    $core.Iterable<ContributorInsightsSummary>? contributorinsightssummaries,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (contributorinsightssummaries != null)
      result.contributorinsightssummaries.addAll(contributorinsightssummaries);
    return result;
  }

  ListContributorInsightsOutput._();

  factory ListContributorInsightsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListContributorInsightsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListContributorInsightsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ContributorInsightsSummary>(
        228077470, _omitFieldNames ? '' : 'contributorinsightssummaries',
        subBuilder: ContributorInsightsSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListContributorInsightsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListContributorInsightsOutput copyWith(
          void Function(ListContributorInsightsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ListContributorInsightsOutput))
          as ListContributorInsightsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListContributorInsightsOutput create() =>
      ListContributorInsightsOutput._();
  @$core.override
  ListContributorInsightsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListContributorInsightsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListContributorInsightsOutput>(create);
  static ListContributorInsightsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(228077470)
  $pb.PbList<ContributorInsightsSummary> get contributorinsightssummaries =>
      $_getList(1);
}

class ListExportsInput extends $pb.GeneratedMessage {
  factory ListExportsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? tablearn,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (tablearn != null) result.tablearn = tablearn;
    return result;
  }

  ListExportsInput._();

  factory ListExportsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListExportsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListExportsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExportsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExportsInput copyWith(void Function(ListExportsInput) updates) =>
      super.copyWith((message) => updates(message as ListExportsInput))
          as ListExportsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListExportsInput create() => ListExportsInput._();
  @$core.override
  ListExportsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListExportsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListExportsInput>(create);
  static ListExportsInput? _defaultInstance;

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

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(2);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(2);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);
}

class ListExportsOutput extends $pb.GeneratedMessage {
  factory ListExportsOutput({
    $core.String? nexttoken,
    $core.Iterable<ExportSummary>? exportsummaries,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (exportsummaries != null) result.exportsummaries.addAll(exportsummaries);
    return result;
  }

  ListExportsOutput._();

  factory ListExportsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListExportsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListExportsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ExportSummary>(378369890, _omitFieldNames ? '' : 'exportsummaries',
        subBuilder: ExportSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExportsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListExportsOutput copyWith(void Function(ListExportsOutput) updates) =>
      super.copyWith((message) => updates(message as ListExportsOutput))
          as ListExportsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListExportsOutput create() => ListExportsOutput._();
  @$core.override
  ListExportsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListExportsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListExportsOutput>(create);
  static ListExportsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(378369890)
  $pb.PbList<ExportSummary> get exportsummaries => $_getList(1);
}

class ListGlobalTablesInput extends $pb.GeneratedMessage {
  factory ListGlobalTablesInput({
    $core.String? exclusivestartglobaltablename,
    $core.String? regionname,
    $core.int? limit,
  }) {
    final result = create();
    if (exclusivestartglobaltablename != null)
      result.exclusivestartglobaltablename = exclusivestartglobaltablename;
    if (regionname != null) result.regionname = regionname;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListGlobalTablesInput._();

  factory ListGlobalTablesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListGlobalTablesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListGlobalTablesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(89274086, _omitFieldNames ? '' : 'exclusivestartglobaltablename')
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGlobalTablesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGlobalTablesInput copyWith(
          void Function(ListGlobalTablesInput) updates) =>
      super.copyWith((message) => updates(message as ListGlobalTablesInput))
          as ListGlobalTablesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListGlobalTablesInput create() => ListGlobalTablesInput._();
  @$core.override
  ListGlobalTablesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListGlobalTablesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListGlobalTablesInput>(create);
  static ListGlobalTablesInput? _defaultInstance;

  @$pb.TagNumber(89274086)
  $core.String get exclusivestartglobaltablename => $_getSZ(0);
  @$pb.TagNumber(89274086)
  set exclusivestartglobaltablename($core.String value) =>
      $_setString(0, value);
  @$pb.TagNumber(89274086)
  $core.bool hasExclusivestartglobaltablename() => $_has(0);
  @$pb.TagNumber(89274086)
  void clearExclusivestartglobaltablename() => $_clearField(89274086);

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(1);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(1);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListGlobalTablesOutput extends $pb.GeneratedMessage {
  factory ListGlobalTablesOutput({
    $core.String? lastevaluatedglobaltablename,
    $core.Iterable<GlobalTable>? globaltables,
  }) {
    final result = create();
    if (lastevaluatedglobaltablename != null)
      result.lastevaluatedglobaltablename = lastevaluatedglobaltablename;
    if (globaltables != null) result.globaltables.addAll(globaltables);
    return result;
  }

  ListGlobalTablesOutput._();

  factory ListGlobalTablesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListGlobalTablesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListGlobalTablesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(157875353, _omitFieldNames ? '' : 'lastevaluatedglobaltablename')
    ..pPM<GlobalTable>(213108516, _omitFieldNames ? '' : 'globaltables',
        subBuilder: GlobalTable.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGlobalTablesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGlobalTablesOutput copyWith(
          void Function(ListGlobalTablesOutput) updates) =>
      super.copyWith((message) => updates(message as ListGlobalTablesOutput))
          as ListGlobalTablesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListGlobalTablesOutput create() => ListGlobalTablesOutput._();
  @$core.override
  ListGlobalTablesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListGlobalTablesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListGlobalTablesOutput>(create);
  static ListGlobalTablesOutput? _defaultInstance;

  @$pb.TagNumber(157875353)
  $core.String get lastevaluatedglobaltablename => $_getSZ(0);
  @$pb.TagNumber(157875353)
  set lastevaluatedglobaltablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(157875353)
  $core.bool hasLastevaluatedglobaltablename() => $_has(0);
  @$pb.TagNumber(157875353)
  void clearLastevaluatedglobaltablename() => $_clearField(157875353);

  @$pb.TagNumber(213108516)
  $pb.PbList<GlobalTable> get globaltables => $_getList(1);
}

class ListImportsInput extends $pb.GeneratedMessage {
  factory ListImportsInput({
    $core.String? nexttoken,
    $core.String? tablearn,
    $core.int? pagesize,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (tablearn != null) result.tablearn = tablearn;
    if (pagesize != null) result.pagesize = pagesize;
    return result;
  }

  ListImportsInput._();

  factory ListImportsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListImportsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListImportsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..aI(438340024, _omitFieldNames ? '' : 'pagesize')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportsInput copyWith(void Function(ListImportsInput) updates) =>
      super.copyWith((message) => updates(message as ListImportsInput))
          as ListImportsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListImportsInput create() => ListImportsInput._();
  @$core.override
  ListImportsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListImportsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListImportsInput>(create);
  static ListImportsInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(1);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(1);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);

  @$pb.TagNumber(438340024)
  $core.int get pagesize => $_getIZ(2);
  @$pb.TagNumber(438340024)
  set pagesize($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(438340024)
  $core.bool hasPagesize() => $_has(2);
  @$pb.TagNumber(438340024)
  void clearPagesize() => $_clearField(438340024);
}

class ListImportsOutput extends $pb.GeneratedMessage {
  factory ListImportsOutput({
    $core.Iterable<ImportSummary>? importsummarylist,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (importsummarylist != null)
      result.importsummarylist.addAll(importsummarylist);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListImportsOutput._();

  factory ListImportsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListImportsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListImportsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<ImportSummary>(19020549, _omitFieldNames ? '' : 'importsummarylist',
        subBuilder: ImportSummary.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListImportsOutput copyWith(void Function(ListImportsOutput) updates) =>
      super.copyWith((message) => updates(message as ListImportsOutput))
          as ListImportsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListImportsOutput create() => ListImportsOutput._();
  @$core.override
  ListImportsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListImportsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListImportsOutput>(create);
  static ListImportsOutput? _defaultInstance;

  @$pb.TagNumber(19020549)
  $pb.PbList<ImportSummary> get importsummarylist => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListTablesInput extends $pb.GeneratedMessage {
  factory ListTablesInput({
    $core.String? exclusivestarttablename,
    $core.int? limit,
  }) {
    final result = create();
    if (exclusivestarttablename != null)
      result.exclusivestarttablename = exclusivestarttablename;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListTablesInput._();

  factory ListTablesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTablesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTablesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(4554457, _omitFieldNames ? '' : 'exclusivestarttablename')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTablesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTablesInput copyWith(void Function(ListTablesInput) updates) =>
      super.copyWith((message) => updates(message as ListTablesInput))
          as ListTablesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTablesInput create() => ListTablesInput._();
  @$core.override
  ListTablesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTablesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTablesInput>(create);
  static ListTablesInput? _defaultInstance;

  @$pb.TagNumber(4554457)
  $core.String get exclusivestarttablename => $_getSZ(0);
  @$pb.TagNumber(4554457)
  set exclusivestarttablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(4554457)
  $core.bool hasExclusivestarttablename() => $_has(0);
  @$pb.TagNumber(4554457)
  void clearExclusivestarttablename() => $_clearField(4554457);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListTablesOutput extends $pb.GeneratedMessage {
  factory ListTablesOutput({
    $core.Iterable<$core.String>? tablenames,
    $core.String? lastevaluatedtablename,
  }) {
    final result = create();
    if (tablenames != null) result.tablenames.addAll(tablenames);
    if (lastevaluatedtablename != null)
      result.lastevaluatedtablename = lastevaluatedtablename;
    return result;
  }

  ListTablesOutput._();

  factory ListTablesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTablesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTablesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPS(354058238, _omitFieldNames ? '' : 'tablenames')
    ..aOS(436880490, _omitFieldNames ? '' : 'lastevaluatedtablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTablesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTablesOutput copyWith(void Function(ListTablesOutput) updates) =>
      super.copyWith((message) => updates(message as ListTablesOutput))
          as ListTablesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTablesOutput create() => ListTablesOutput._();
  @$core.override
  ListTablesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTablesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTablesOutput>(create);
  static ListTablesOutput? _defaultInstance;

  @$pb.TagNumber(354058238)
  $pb.PbList<$core.String> get tablenames => $_getList(0);

  @$pb.TagNumber(436880490)
  $core.String get lastevaluatedtablename => $_getSZ(1);
  @$pb.TagNumber(436880490)
  set lastevaluatedtablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(436880490)
  $core.bool hasLastevaluatedtablename() => $_has(1);
  @$pb.TagNumber(436880490)
  void clearLastevaluatedtablename() => $_clearField(436880490);
}

class ListTagsOfResourceInput extends $pb.GeneratedMessage {
  factory ListTagsOfResourceInput({
    $core.String? nexttoken,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  ListTagsOfResourceInput._();

  factory ListTagsOfResourceInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsOfResourceInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsOfResourceInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsOfResourceInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsOfResourceInput copyWith(
          void Function(ListTagsOfResourceInput) updates) =>
      super.copyWith((message) => updates(message as ListTagsOfResourceInput))
          as ListTagsOfResourceInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsOfResourceInput create() => ListTagsOfResourceInput._();
  @$core.override
  ListTagsOfResourceInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsOfResourceInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsOfResourceInput>(create);
  static ListTagsOfResourceInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class ListTagsOfResourceOutput extends $pb.GeneratedMessage {
  factory ListTagsOfResourceOutput({
    $core.String? nexttoken,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  ListTagsOfResourceOutput._();

  factory ListTagsOfResourceOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsOfResourceOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsOfResourceOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsOfResourceOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsOfResourceOutput copyWith(
          void Function(ListTagsOfResourceOutput) updates) =>
      super.copyWith((message) => updates(message as ListTagsOfResourceOutput))
          as ListTagsOfResourceOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsOfResourceOutput create() => ListTagsOfResourceOutput._();
  @$core.override
  ListTagsOfResourceOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsOfResourceOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsOfResourceOutput>(create);
  static ListTagsOfResourceOutput? _defaultInstance;

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

class LocalSecondaryIndex extends $pb.GeneratedMessage {
  factory LocalSecondaryIndex({
    $core.String? indexname,
    Projection? projection,
    $core.Iterable<KeySchemaElement>? keyschema,
  }) {
    final result = create();
    if (indexname != null) result.indexname = indexname;
    if (projection != null) result.projection = projection;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    return result;
  }

  LocalSecondaryIndex._();

  factory LocalSecondaryIndex.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocalSecondaryIndex.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocalSecondaryIndex',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<Projection>(105045921, _omitFieldNames ? '' : 'projection',
        subBuilder: Projection.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocalSecondaryIndex clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocalSecondaryIndex copyWith(void Function(LocalSecondaryIndex) updates) =>
      super.copyWith((message) => updates(message as LocalSecondaryIndex))
          as LocalSecondaryIndex;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LocalSecondaryIndex create() => LocalSecondaryIndex._();
  @$core.override
  LocalSecondaryIndex createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocalSecondaryIndex getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocalSecondaryIndex>(create);
  static LocalSecondaryIndex? _defaultInstance;

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(0);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(0);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(105045921)
  Projection get projection => $_getN(1);
  @$pb.TagNumber(105045921)
  set projection(Projection value) => $_setField(105045921, value);
  @$pb.TagNumber(105045921)
  $core.bool hasProjection() => $_has(1);
  @$pb.TagNumber(105045921)
  void clearProjection() => $_clearField(105045921);
  @$pb.TagNumber(105045921)
  Projection ensureProjection() => $_ensure(1);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(2);
}

class LocalSecondaryIndexDescription extends $pb.GeneratedMessage {
  factory LocalSecondaryIndexDescription({
    $fixnum.Int64? itemcount,
    $core.String? indexname,
    Projection? projection,
    $core.Iterable<KeySchemaElement>? keyschema,
    $core.String? indexarn,
    $fixnum.Int64? indexsizebytes,
  }) {
    final result = create();
    if (itemcount != null) result.itemcount = itemcount;
    if (indexname != null) result.indexname = indexname;
    if (projection != null) result.projection = projection;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (indexarn != null) result.indexarn = indexarn;
    if (indexsizebytes != null) result.indexsizebytes = indexsizebytes;
    return result;
  }

  LocalSecondaryIndexDescription._();

  factory LocalSecondaryIndexDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocalSecondaryIndexDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocalSecondaryIndexDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(26280022, _omitFieldNames ? '' : 'itemcount')
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<Projection>(105045921, _omitFieldNames ? '' : 'projection',
        subBuilder: Projection.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..aOS(374335615, _omitFieldNames ? '' : 'indexarn')
    ..aInt64(395738346, _omitFieldNames ? '' : 'indexsizebytes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocalSecondaryIndexDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocalSecondaryIndexDescription copyWith(
          void Function(LocalSecondaryIndexDescription) updates) =>
      super.copyWith(
              (message) => updates(message as LocalSecondaryIndexDescription))
          as LocalSecondaryIndexDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LocalSecondaryIndexDescription create() =>
      LocalSecondaryIndexDescription._();
  @$core.override
  LocalSecondaryIndexDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocalSecondaryIndexDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocalSecondaryIndexDescription>(create);
  static LocalSecondaryIndexDescription? _defaultInstance;

  @$pb.TagNumber(26280022)
  $fixnum.Int64 get itemcount => $_getI64(0);
  @$pb.TagNumber(26280022)
  set itemcount($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(26280022)
  $core.bool hasItemcount() => $_has(0);
  @$pb.TagNumber(26280022)
  void clearItemcount() => $_clearField(26280022);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(105045921)
  Projection get projection => $_getN(2);
  @$pb.TagNumber(105045921)
  set projection(Projection value) => $_setField(105045921, value);
  @$pb.TagNumber(105045921)
  $core.bool hasProjection() => $_has(2);
  @$pb.TagNumber(105045921)
  void clearProjection() => $_clearField(105045921);
  @$pb.TagNumber(105045921)
  Projection ensureProjection() => $_ensure(2);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(3);

  @$pb.TagNumber(374335615)
  $core.String get indexarn => $_getSZ(4);
  @$pb.TagNumber(374335615)
  set indexarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(374335615)
  $core.bool hasIndexarn() => $_has(4);
  @$pb.TagNumber(374335615)
  void clearIndexarn() => $_clearField(374335615);

  @$pb.TagNumber(395738346)
  $fixnum.Int64 get indexsizebytes => $_getI64(5);
  @$pb.TagNumber(395738346)
  set indexsizebytes($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(395738346)
  $core.bool hasIndexsizebytes() => $_has(5);
  @$pb.TagNumber(395738346)
  void clearIndexsizebytes() => $_clearField(395738346);
}

class LocalSecondaryIndexInfo extends $pb.GeneratedMessage {
  factory LocalSecondaryIndexInfo({
    $core.String? indexname,
    Projection? projection,
    $core.Iterable<KeySchemaElement>? keyschema,
  }) {
    final result = create();
    if (indexname != null) result.indexname = indexname;
    if (projection != null) result.projection = projection;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    return result;
  }

  LocalSecondaryIndexInfo._();

  factory LocalSecondaryIndexInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocalSecondaryIndexInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocalSecondaryIndexInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<Projection>(105045921, _omitFieldNames ? '' : 'projection',
        subBuilder: Projection.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocalSecondaryIndexInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocalSecondaryIndexInfo copyWith(
          void Function(LocalSecondaryIndexInfo) updates) =>
      super.copyWith((message) => updates(message as LocalSecondaryIndexInfo))
          as LocalSecondaryIndexInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LocalSecondaryIndexInfo create() => LocalSecondaryIndexInfo._();
  @$core.override
  LocalSecondaryIndexInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocalSecondaryIndexInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocalSecondaryIndexInfo>(create);
  static LocalSecondaryIndexInfo? _defaultInstance;

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(0);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(0);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(105045921)
  Projection get projection => $_getN(1);
  @$pb.TagNumber(105045921)
  set projection(Projection value) => $_setField(105045921, value);
  @$pb.TagNumber(105045921)
  $core.bool hasProjection() => $_has(1);
  @$pb.TagNumber(105045921)
  void clearProjection() => $_clearField(105045921);
  @$pb.TagNumber(105045921)
  Projection ensureProjection() => $_ensure(1);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(2);
}

class OnDemandThroughput extends $pb.GeneratedMessage {
  factory OnDemandThroughput({
    $fixnum.Int64? maxwriterequestunits,
    $fixnum.Int64? maxreadrequestunits,
  }) {
    final result = create();
    if (maxwriterequestunits != null)
      result.maxwriterequestunits = maxwriterequestunits;
    if (maxreadrequestunits != null)
      result.maxreadrequestunits = maxreadrequestunits;
    return result;
  }

  OnDemandThroughput._();

  factory OnDemandThroughput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OnDemandThroughput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OnDemandThroughput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(99152185, _omitFieldNames ? '' : 'maxwriterequestunits')
    ..aInt64(361996322, _omitFieldNames ? '' : 'maxreadrequestunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OnDemandThroughput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OnDemandThroughput copyWith(void Function(OnDemandThroughput) updates) =>
      super.copyWith((message) => updates(message as OnDemandThroughput))
          as OnDemandThroughput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OnDemandThroughput create() => OnDemandThroughput._();
  @$core.override
  OnDemandThroughput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OnDemandThroughput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OnDemandThroughput>(create);
  static OnDemandThroughput? _defaultInstance;

  @$pb.TagNumber(99152185)
  $fixnum.Int64 get maxwriterequestunits => $_getI64(0);
  @$pb.TagNumber(99152185)
  set maxwriterequestunits($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(99152185)
  $core.bool hasMaxwriterequestunits() => $_has(0);
  @$pb.TagNumber(99152185)
  void clearMaxwriterequestunits() => $_clearField(99152185);

  @$pb.TagNumber(361996322)
  $fixnum.Int64 get maxreadrequestunits => $_getI64(1);
  @$pb.TagNumber(361996322)
  set maxreadrequestunits($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(361996322)
  $core.bool hasMaxreadrequestunits() => $_has(1);
  @$pb.TagNumber(361996322)
  void clearMaxreadrequestunits() => $_clearField(361996322);
}

class OnDemandThroughputOverride extends $pb.GeneratedMessage {
  factory OnDemandThroughputOverride({
    $fixnum.Int64? maxreadrequestunits,
  }) {
    final result = create();
    if (maxreadrequestunits != null)
      result.maxreadrequestunits = maxreadrequestunits;
    return result;
  }

  OnDemandThroughputOverride._();

  factory OnDemandThroughputOverride.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OnDemandThroughputOverride.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OnDemandThroughputOverride',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(361996322, _omitFieldNames ? '' : 'maxreadrequestunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OnDemandThroughputOverride clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OnDemandThroughputOverride copyWith(
          void Function(OnDemandThroughputOverride) updates) =>
      super.copyWith(
              (message) => updates(message as OnDemandThroughputOverride))
          as OnDemandThroughputOverride;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OnDemandThroughputOverride create() => OnDemandThroughputOverride._();
  @$core.override
  OnDemandThroughputOverride createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OnDemandThroughputOverride getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OnDemandThroughputOverride>(create);
  static OnDemandThroughputOverride? _defaultInstance;

  @$pb.TagNumber(361996322)
  $fixnum.Int64 get maxreadrequestunits => $_getI64(0);
  @$pb.TagNumber(361996322)
  set maxreadrequestunits($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(361996322)
  $core.bool hasMaxreadrequestunits() => $_has(0);
  @$pb.TagNumber(361996322)
  void clearMaxreadrequestunits() => $_clearField(361996322);
}

class ParameterizedStatement extends $pb.GeneratedMessage {
  factory ParameterizedStatement({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    $core.String? statement,
    $core.Iterable<AttributeValue>? parameters,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (statement != null) result.statement = statement;
    if (parameters != null) result.parameters.addAll(parameters);
    return result;
  }

  ParameterizedStatement._();

  factory ParameterizedStatement.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ParameterizedStatement.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ParameterizedStatement',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..aOS(248790199, _omitFieldNames ? '' : 'statement')
    ..pPM<AttributeValue>(494900218, _omitFieldNames ? '' : 'parameters',
        subBuilder: AttributeValue.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ParameterizedStatement clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ParameterizedStatement copyWith(
          void Function(ParameterizedStatement) updates) =>
      super.copyWith((message) => updates(message as ParameterizedStatement))
          as ParameterizedStatement;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ParameterizedStatement create() => ParameterizedStatement._();
  @$core.override
  ParameterizedStatement createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ParameterizedStatement getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ParameterizedStatement>(create);
  static ParameterizedStatement? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(248790199)
  $core.String get statement => $_getSZ(1);
  @$pb.TagNumber(248790199)
  set statement($core.String value) => $_setString(1, value);
  @$pb.TagNumber(248790199)
  $core.bool hasStatement() => $_has(1);
  @$pb.TagNumber(248790199)
  void clearStatement() => $_clearField(248790199);

  @$pb.TagNumber(494900218)
  $pb.PbList<AttributeValue> get parameters => $_getList(2);
}

class PointInTimeRecoveryDescription extends $pb.GeneratedMessage {
  factory PointInTimeRecoveryDescription({
    $core.int? recoveryperiodindays,
    $core.String? latestrestorabledatetime,
    $core.String? earliestrestorabledatetime,
    PointInTimeRecoveryStatus? pointintimerecoverystatus,
  }) {
    final result = create();
    if (recoveryperiodindays != null)
      result.recoveryperiodindays = recoveryperiodindays;
    if (latestrestorabledatetime != null)
      result.latestrestorabledatetime = latestrestorabledatetime;
    if (earliestrestorabledatetime != null)
      result.earliestrestorabledatetime = earliestrestorabledatetime;
    if (pointintimerecoverystatus != null)
      result.pointintimerecoverystatus = pointintimerecoverystatus;
    return result;
  }

  PointInTimeRecoveryDescription._();

  factory PointInTimeRecoveryDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PointInTimeRecoveryDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PointInTimeRecoveryDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aI(67453592, _omitFieldNames ? '' : 'recoveryperiodindays')
    ..aOS(131030757, _omitFieldNames ? '' : 'latestrestorabledatetime')
    ..aOS(440478879, _omitFieldNames ? '' : 'earliestrestorabledatetime')
    ..aE<PointInTimeRecoveryStatus>(
        510838075, _omitFieldNames ? '' : 'pointintimerecoverystatus',
        enumValues: PointInTimeRecoveryStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PointInTimeRecoveryDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PointInTimeRecoveryDescription copyWith(
          void Function(PointInTimeRecoveryDescription) updates) =>
      super.copyWith(
              (message) => updates(message as PointInTimeRecoveryDescription))
          as PointInTimeRecoveryDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PointInTimeRecoveryDescription create() =>
      PointInTimeRecoveryDescription._();
  @$core.override
  PointInTimeRecoveryDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PointInTimeRecoveryDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PointInTimeRecoveryDescription>(create);
  static PointInTimeRecoveryDescription? _defaultInstance;

  @$pb.TagNumber(67453592)
  $core.int get recoveryperiodindays => $_getIZ(0);
  @$pb.TagNumber(67453592)
  set recoveryperiodindays($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(67453592)
  $core.bool hasRecoveryperiodindays() => $_has(0);
  @$pb.TagNumber(67453592)
  void clearRecoveryperiodindays() => $_clearField(67453592);

  @$pb.TagNumber(131030757)
  $core.String get latestrestorabledatetime => $_getSZ(1);
  @$pb.TagNumber(131030757)
  set latestrestorabledatetime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(131030757)
  $core.bool hasLatestrestorabledatetime() => $_has(1);
  @$pb.TagNumber(131030757)
  void clearLatestrestorabledatetime() => $_clearField(131030757);

  @$pb.TagNumber(440478879)
  $core.String get earliestrestorabledatetime => $_getSZ(2);
  @$pb.TagNumber(440478879)
  set earliestrestorabledatetime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(440478879)
  $core.bool hasEarliestrestorabledatetime() => $_has(2);
  @$pb.TagNumber(440478879)
  void clearEarliestrestorabledatetime() => $_clearField(440478879);

  @$pb.TagNumber(510838075)
  PointInTimeRecoveryStatus get pointintimerecoverystatus => $_getN(3);
  @$pb.TagNumber(510838075)
  set pointintimerecoverystatus(PointInTimeRecoveryStatus value) =>
      $_setField(510838075, value);
  @$pb.TagNumber(510838075)
  $core.bool hasPointintimerecoverystatus() => $_has(3);
  @$pb.TagNumber(510838075)
  void clearPointintimerecoverystatus() => $_clearField(510838075);
}

class PointInTimeRecoverySpecification extends $pb.GeneratedMessage {
  factory PointInTimeRecoverySpecification({
    $core.int? recoveryperiodindays,
    $core.bool? pointintimerecoveryenabled,
  }) {
    final result = create();
    if (recoveryperiodindays != null)
      result.recoveryperiodindays = recoveryperiodindays;
    if (pointintimerecoveryenabled != null)
      result.pointintimerecoveryenabled = pointintimerecoveryenabled;
    return result;
  }

  PointInTimeRecoverySpecification._();

  factory PointInTimeRecoverySpecification.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PointInTimeRecoverySpecification.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PointInTimeRecoverySpecification',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aI(67453592, _omitFieldNames ? '' : 'recoveryperiodindays')
    ..aOB(294176170, _omitFieldNames ? '' : 'pointintimerecoveryenabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PointInTimeRecoverySpecification clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PointInTimeRecoverySpecification copyWith(
          void Function(PointInTimeRecoverySpecification) updates) =>
      super.copyWith(
              (message) => updates(message as PointInTimeRecoverySpecification))
          as PointInTimeRecoverySpecification;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PointInTimeRecoverySpecification create() =>
      PointInTimeRecoverySpecification._();
  @$core.override
  PointInTimeRecoverySpecification createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PointInTimeRecoverySpecification getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PointInTimeRecoverySpecification>(
          create);
  static PointInTimeRecoverySpecification? _defaultInstance;

  @$pb.TagNumber(67453592)
  $core.int get recoveryperiodindays => $_getIZ(0);
  @$pb.TagNumber(67453592)
  set recoveryperiodindays($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(67453592)
  $core.bool hasRecoveryperiodindays() => $_has(0);
  @$pb.TagNumber(67453592)
  void clearRecoveryperiodindays() => $_clearField(67453592);

  @$pb.TagNumber(294176170)
  $core.bool get pointintimerecoveryenabled => $_getBF(1);
  @$pb.TagNumber(294176170)
  set pointintimerecoveryenabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(294176170)
  $core.bool hasPointintimerecoveryenabled() => $_has(1);
  @$pb.TagNumber(294176170)
  void clearPointintimerecoveryenabled() => $_clearField(294176170);
}

class PointInTimeRecoveryUnavailableException extends $pb.GeneratedMessage {
  factory PointInTimeRecoveryUnavailableException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PointInTimeRecoveryUnavailableException._();

  factory PointInTimeRecoveryUnavailableException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PointInTimeRecoveryUnavailableException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PointInTimeRecoveryUnavailableException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PointInTimeRecoveryUnavailableException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PointInTimeRecoveryUnavailableException copyWith(
          void Function(PointInTimeRecoveryUnavailableException) updates) =>
      super.copyWith((message) =>
              updates(message as PointInTimeRecoveryUnavailableException))
          as PointInTimeRecoveryUnavailableException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PointInTimeRecoveryUnavailableException create() =>
      PointInTimeRecoveryUnavailableException._();
  @$core.override
  PointInTimeRecoveryUnavailableException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PointInTimeRecoveryUnavailableException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          PointInTimeRecoveryUnavailableException>(create);
  static PointInTimeRecoveryUnavailableException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PolicyNotFoundException extends $pb.GeneratedMessage {
  factory PolicyNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PolicyNotFoundException._();

  factory PolicyNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PolicyNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PolicyNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PolicyNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PolicyNotFoundException copyWith(
          void Function(PolicyNotFoundException) updates) =>
      super.copyWith((message) => updates(message as PolicyNotFoundException))
          as PolicyNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PolicyNotFoundException create() => PolicyNotFoundException._();
  @$core.override
  PolicyNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PolicyNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PolicyNotFoundException>(create);
  static PolicyNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class Projection extends $pb.GeneratedMessage {
  factory Projection({
    ProjectionType? projectiontype,
    $core.Iterable<$core.String>? nonkeyattributes,
  }) {
    final result = create();
    if (projectiontype != null) result.projectiontype = projectiontype;
    if (nonkeyattributes != null)
      result.nonkeyattributes.addAll(nonkeyattributes);
    return result;
  }

  Projection._();

  factory Projection.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Projection.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Projection',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ProjectionType>(120720617, _omitFieldNames ? '' : 'projectiontype',
        enumValues: ProjectionType.values)
    ..pPS(312245447, _omitFieldNames ? '' : 'nonkeyattributes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Projection clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Projection copyWith(void Function(Projection) updates) =>
      super.copyWith((message) => updates(message as Projection)) as Projection;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Projection create() => Projection._();
  @$core.override
  Projection createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Projection getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Projection>(create);
  static Projection? _defaultInstance;

  @$pb.TagNumber(120720617)
  ProjectionType get projectiontype => $_getN(0);
  @$pb.TagNumber(120720617)
  set projectiontype(ProjectionType value) => $_setField(120720617, value);
  @$pb.TagNumber(120720617)
  $core.bool hasProjectiontype() => $_has(0);
  @$pb.TagNumber(120720617)
  void clearProjectiontype() => $_clearField(120720617);

  @$pb.TagNumber(312245447)
  $pb.PbList<$core.String> get nonkeyattributes => $_getList(1);
}

class ProvisionedThroughput extends $pb.GeneratedMessage {
  factory ProvisionedThroughput({
    $fixnum.Int64? writecapacityunits,
    $fixnum.Int64? readcapacityunits,
  }) {
    final result = create();
    if (writecapacityunits != null)
      result.writecapacityunits = writecapacityunits;
    if (readcapacityunits != null) result.readcapacityunits = readcapacityunits;
    return result;
  }

  ProvisionedThroughput._();

  factory ProvisionedThroughput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProvisionedThroughput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProvisionedThroughput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(26938032, _omitFieldNames ? '' : 'writecapacityunits')
    ..aInt64(43945489, _omitFieldNames ? '' : 'readcapacityunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedThroughput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedThroughput copyWith(
          void Function(ProvisionedThroughput) updates) =>
      super.copyWith((message) => updates(message as ProvisionedThroughput))
          as ProvisionedThroughput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProvisionedThroughput create() => ProvisionedThroughput._();
  @$core.override
  ProvisionedThroughput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProvisionedThroughput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProvisionedThroughput>(create);
  static ProvisionedThroughput? _defaultInstance;

  @$pb.TagNumber(26938032)
  $fixnum.Int64 get writecapacityunits => $_getI64(0);
  @$pb.TagNumber(26938032)
  set writecapacityunits($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(26938032)
  $core.bool hasWritecapacityunits() => $_has(0);
  @$pb.TagNumber(26938032)
  void clearWritecapacityunits() => $_clearField(26938032);

  @$pb.TagNumber(43945489)
  $fixnum.Int64 get readcapacityunits => $_getI64(1);
  @$pb.TagNumber(43945489)
  set readcapacityunits($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(43945489)
  $core.bool hasReadcapacityunits() => $_has(1);
  @$pb.TagNumber(43945489)
  void clearReadcapacityunits() => $_clearField(43945489);
}

class ProvisionedThroughputDescription extends $pb.GeneratedMessage {
  factory ProvisionedThroughputDescription({
    $fixnum.Int64? writecapacityunits,
    $fixnum.Int64? numberofdecreasestoday,
    $fixnum.Int64? readcapacityunits,
    $core.String? lastincreasedatetime,
    $core.String? lastdecreasedatetime,
  }) {
    final result = create();
    if (writecapacityunits != null)
      result.writecapacityunits = writecapacityunits;
    if (numberofdecreasestoday != null)
      result.numberofdecreasestoday = numberofdecreasestoday;
    if (readcapacityunits != null) result.readcapacityunits = readcapacityunits;
    if (lastincreasedatetime != null)
      result.lastincreasedatetime = lastincreasedatetime;
    if (lastdecreasedatetime != null)
      result.lastdecreasedatetime = lastdecreasedatetime;
    return result;
  }

  ProvisionedThroughputDescription._();

  factory ProvisionedThroughputDescription.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProvisionedThroughputDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProvisionedThroughputDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(26938032, _omitFieldNames ? '' : 'writecapacityunits')
    ..aInt64(34313772, _omitFieldNames ? '' : 'numberofdecreasestoday')
    ..aInt64(43945489, _omitFieldNames ? '' : 'readcapacityunits')
    ..aOS(490688181, _omitFieldNames ? '' : 'lastincreasedatetime')
    ..aOS(494947041, _omitFieldNames ? '' : 'lastdecreasedatetime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedThroughputDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedThroughputDescription copyWith(
          void Function(ProvisionedThroughputDescription) updates) =>
      super.copyWith(
              (message) => updates(message as ProvisionedThroughputDescription))
          as ProvisionedThroughputDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProvisionedThroughputDescription create() =>
      ProvisionedThroughputDescription._();
  @$core.override
  ProvisionedThroughputDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProvisionedThroughputDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProvisionedThroughputDescription>(
          create);
  static ProvisionedThroughputDescription? _defaultInstance;

  @$pb.TagNumber(26938032)
  $fixnum.Int64 get writecapacityunits => $_getI64(0);
  @$pb.TagNumber(26938032)
  set writecapacityunits($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(26938032)
  $core.bool hasWritecapacityunits() => $_has(0);
  @$pb.TagNumber(26938032)
  void clearWritecapacityunits() => $_clearField(26938032);

  @$pb.TagNumber(34313772)
  $fixnum.Int64 get numberofdecreasestoday => $_getI64(1);
  @$pb.TagNumber(34313772)
  set numberofdecreasestoday($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(34313772)
  $core.bool hasNumberofdecreasestoday() => $_has(1);
  @$pb.TagNumber(34313772)
  void clearNumberofdecreasestoday() => $_clearField(34313772);

  @$pb.TagNumber(43945489)
  $fixnum.Int64 get readcapacityunits => $_getI64(2);
  @$pb.TagNumber(43945489)
  set readcapacityunits($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(43945489)
  $core.bool hasReadcapacityunits() => $_has(2);
  @$pb.TagNumber(43945489)
  void clearReadcapacityunits() => $_clearField(43945489);

  @$pb.TagNumber(490688181)
  $core.String get lastincreasedatetime => $_getSZ(3);
  @$pb.TagNumber(490688181)
  set lastincreasedatetime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(490688181)
  $core.bool hasLastincreasedatetime() => $_has(3);
  @$pb.TagNumber(490688181)
  void clearLastincreasedatetime() => $_clearField(490688181);

  @$pb.TagNumber(494947041)
  $core.String get lastdecreasedatetime => $_getSZ(4);
  @$pb.TagNumber(494947041)
  set lastdecreasedatetime($core.String value) => $_setString(4, value);
  @$pb.TagNumber(494947041)
  $core.bool hasLastdecreasedatetime() => $_has(4);
  @$pb.TagNumber(494947041)
  void clearLastdecreasedatetime() => $_clearField(494947041);
}

class ProvisionedThroughputExceededException extends $pb.GeneratedMessage {
  factory ProvisionedThroughputExceededException({
    $core.String? message,
    $core.Iterable<ThrottlingReason>? throttlingreasons,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (throttlingreasons != null)
      result.throttlingreasons.addAll(throttlingreasons);
    return result;
  }

  ProvisionedThroughputExceededException._();

  factory ProvisionedThroughputExceededException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProvisionedThroughputExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProvisionedThroughputExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..pPM<ThrottlingReason>(
        284436852, _omitFieldNames ? '' : 'throttlingreasons',
        subBuilder: ThrottlingReason.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedThroughputExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedThroughputExceededException copyWith(
          void Function(ProvisionedThroughputExceededException) updates) =>
      super.copyWith((message) =>
              updates(message as ProvisionedThroughputExceededException))
          as ProvisionedThroughputExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProvisionedThroughputExceededException create() =>
      ProvisionedThroughputExceededException._();
  @$core.override
  ProvisionedThroughputExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProvisionedThroughputExceededException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ProvisionedThroughputExceededException>(create);
  static ProvisionedThroughputExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(284436852)
  $pb.PbList<ThrottlingReason> get throttlingreasons => $_getList(1);
}

class ProvisionedThroughputOverride extends $pb.GeneratedMessage {
  factory ProvisionedThroughputOverride({
    $fixnum.Int64? readcapacityunits,
  }) {
    final result = create();
    if (readcapacityunits != null) result.readcapacityunits = readcapacityunits;
    return result;
  }

  ProvisionedThroughputOverride._();

  factory ProvisionedThroughputOverride.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProvisionedThroughputOverride.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProvisionedThroughputOverride',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(43945489, _omitFieldNames ? '' : 'readcapacityunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedThroughputOverride clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProvisionedThroughputOverride copyWith(
          void Function(ProvisionedThroughputOverride) updates) =>
      super.copyWith(
              (message) => updates(message as ProvisionedThroughputOverride))
          as ProvisionedThroughputOverride;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProvisionedThroughputOverride create() =>
      ProvisionedThroughputOverride._();
  @$core.override
  ProvisionedThroughputOverride createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProvisionedThroughputOverride getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProvisionedThroughputOverride>(create);
  static ProvisionedThroughputOverride? _defaultInstance;

  @$pb.TagNumber(43945489)
  $fixnum.Int64 get readcapacityunits => $_getI64(0);
  @$pb.TagNumber(43945489)
  set readcapacityunits($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(43945489)
  $core.bool hasReadcapacityunits() => $_has(0);
  @$pb.TagNumber(43945489)
  void clearReadcapacityunits() => $_clearField(43945489);
}

class Put extends $pb.GeneratedMessage {
  factory Put({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.String? tablename,
    $core.String? conditionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (tablename != null) result.tablename = tablename;
    if (conditionexpression != null)
      result.conditionexpression = conditionexpression;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    if (item != null) result.item.addEntries(item);
    return result;
  }

  Put._();

  factory Put.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Put.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Put',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'Put.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(409657405, _omitFieldNames ? '' : 'conditionexpression')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'Put.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'Put.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Put clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Put copyWith(void Function(Put) updates) =>
      super.copyWith((message) => updates(message as Put)) as Put;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Put create() => Put._();
  @$core.override
  Put createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Put getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Put>(create);
  static Put? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(1);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(2);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(2);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(409657405)
  $core.String get conditionexpression => $_getSZ(3);
  @$pb.TagNumber(409657405)
  set conditionexpression($core.String value) => $_setString(3, value);
  @$pb.TagNumber(409657405)
  $core.bool hasConditionexpression() => $_has(3);
  @$pb.TagNumber(409657405)
  void clearConditionexpression() => $_clearField(409657405);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(4);

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(5);
}

class PutItemInput extends $pb.GeneratedMessage {
  factory PutItemInput({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<$core.MapEntry<$core.String, ExpectedAttributeValue>>?
        expected,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    ConditionalOperator? conditionaloperator,
    ReturnItemCollectionMetrics? returnitemcollectionmetrics,
    $core.String? tablename,
    ReturnValue? returnvalues,
    $core.String? conditionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (expected != null) result.expected.addEntries(expected);
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (conditionaloperator != null)
      result.conditionaloperator = conditionaloperator;
    if (returnitemcollectionmetrics != null)
      result.returnitemcollectionmetrics = returnitemcollectionmetrics;
    if (tablename != null) result.tablename = tablename;
    if (returnvalues != null) result.returnvalues = returnvalues;
    if (conditionexpression != null)
      result.conditionexpression = conditionexpression;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    if (item != null) result.item.addEntries(item);
    return result;
  }

  PutItemInput._();

  factory PutItemInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutItemInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutItemInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..m<$core.String, ExpectedAttributeValue>(
        106557946, _omitFieldNames ? '' : 'expected',
        entryClassName: 'PutItemInput.ExpectedEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: ExpectedAttributeValue.create,
        valueDefaultOrMaker: ExpectedAttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'PutItemInput.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<ConditionalOperator>(
        172066260, _omitFieldNames ? '' : 'conditionaloperator',
        enumValues: ConditionalOperator.values)
    ..aE<ReturnItemCollectionMetrics>(
        255507354, _omitFieldNames ? '' : 'returnitemcollectionmetrics',
        enumValues: ReturnItemCollectionMetrics.values)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aE<ReturnValue>(402960198, _omitFieldNames ? '' : 'returnvalues',
        enumValues: ReturnValue.values)
    ..aOS(409657405, _omitFieldNames ? '' : 'conditionexpression')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'PutItemInput.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'PutItemInput.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutItemInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutItemInput copyWith(void Function(PutItemInput) updates) =>
      super.copyWith((message) => updates(message as PutItemInput))
          as PutItemInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutItemInput create() => PutItemInput._();
  @$core.override
  PutItemInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutItemInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutItemInput>(create);
  static PutItemInput? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(1);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(1);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(106557946)
  $pb.PbMap<$core.String, ExpectedAttributeValue> get expected => $_getMap(2);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(3);

  @$pb.TagNumber(172066260)
  ConditionalOperator get conditionaloperator => $_getN(4);
  @$pb.TagNumber(172066260)
  set conditionaloperator(ConditionalOperator value) =>
      $_setField(172066260, value);
  @$pb.TagNumber(172066260)
  $core.bool hasConditionaloperator() => $_has(4);
  @$pb.TagNumber(172066260)
  void clearConditionaloperator() => $_clearField(172066260);

  @$pb.TagNumber(255507354)
  ReturnItemCollectionMetrics get returnitemcollectionmetrics => $_getN(5);
  @$pb.TagNumber(255507354)
  set returnitemcollectionmetrics(ReturnItemCollectionMetrics value) =>
      $_setField(255507354, value);
  @$pb.TagNumber(255507354)
  $core.bool hasReturnitemcollectionmetrics() => $_has(5);
  @$pb.TagNumber(255507354)
  void clearReturnitemcollectionmetrics() => $_clearField(255507354);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(6);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(6, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(6);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(402960198)
  ReturnValue get returnvalues => $_getN(7);
  @$pb.TagNumber(402960198)
  set returnvalues(ReturnValue value) => $_setField(402960198, value);
  @$pb.TagNumber(402960198)
  $core.bool hasReturnvalues() => $_has(7);
  @$pb.TagNumber(402960198)
  void clearReturnvalues() => $_clearField(402960198);

  @$pb.TagNumber(409657405)
  $core.String get conditionexpression => $_getSZ(8);
  @$pb.TagNumber(409657405)
  set conditionexpression($core.String value) => $_setString(8, value);
  @$pb.TagNumber(409657405)
  $core.bool hasConditionexpression() => $_has(8);
  @$pb.TagNumber(409657405)
  void clearConditionexpression() => $_clearField(409657405);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(9);

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(10);
}

class PutItemOutput extends $pb.GeneratedMessage {
  factory PutItemOutput({
    ItemCollectionMetrics? itemcollectionmetrics,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? attributes,
    ConsumedCapacity? consumedcapacity,
  }) {
    final result = create();
    if (itemcollectionmetrics != null)
      result.itemcollectionmetrics = itemcollectionmetrics;
    if (attributes != null) result.attributes.addEntries(attributes);
    if (consumedcapacity != null) result.consumedcapacity = consumedcapacity;
    return result;
  }

  PutItemOutput._();

  factory PutItemOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutItemOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutItemOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ItemCollectionMetrics>(
        185317452, _omitFieldNames ? '' : 'itemcollectionmetrics',
        subBuilder: ItemCollectionMetrics.create)
    ..m<$core.String, AttributeValue>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'PutItemOutput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutItemOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutItemOutput copyWith(void Function(PutItemOutput) updates) =>
      super.copyWith((message) => updates(message as PutItemOutput))
          as PutItemOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutItemOutput create() => PutItemOutput._();
  @$core.override
  PutItemOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutItemOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutItemOutput>(create);
  static PutItemOutput? _defaultInstance;

  @$pb.TagNumber(185317452)
  ItemCollectionMetrics get itemcollectionmetrics => $_getN(0);
  @$pb.TagNumber(185317452)
  set itemcollectionmetrics(ItemCollectionMetrics value) =>
      $_setField(185317452, value);
  @$pb.TagNumber(185317452)
  $core.bool hasItemcollectionmetrics() => $_has(0);
  @$pb.TagNumber(185317452)
  void clearItemcollectionmetrics() => $_clearField(185317452);
  @$pb.TagNumber(185317452)
  ItemCollectionMetrics ensureItemcollectionmetrics() => $_ensure(0);

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, AttributeValue> get attributes => $_getMap(1);

  @$pb.TagNumber(449336620)
  ConsumedCapacity get consumedcapacity => $_getN(2);
  @$pb.TagNumber(449336620)
  set consumedcapacity(ConsumedCapacity value) => $_setField(449336620, value);
  @$pb.TagNumber(449336620)
  $core.bool hasConsumedcapacity() => $_has(2);
  @$pb.TagNumber(449336620)
  void clearConsumedcapacity() => $_clearField(449336620);
  @$pb.TagNumber(449336620)
  ConsumedCapacity ensureConsumedcapacity() => $_ensure(2);
}

class PutRequest extends $pb.GeneratedMessage {
  factory PutRequest({
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? item,
  }) {
    final result = create();
    if (item != null) result.item.addEntries(item);
    return result;
  }

  PutRequest._();

  factory PutRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, AttributeValue>(526680071, _omitFieldNames ? '' : 'item',
        entryClassName: 'PutRequest.ItemEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRequest copyWith(void Function(PutRequest) updates) =>
      super.copyWith((message) => updates(message as PutRequest)) as PutRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRequest create() => PutRequest._();
  @$core.override
  PutRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRequest>(create);
  static PutRequest? _defaultInstance;

  @$pb.TagNumber(526680071)
  $pb.PbMap<$core.String, AttributeValue> get item => $_getMap(0);
}

class PutResourcePolicyInput extends $pb.GeneratedMessage {
  factory PutResourcePolicyInput({
    $core.bool? confirmremoveselfresourceaccess,
    $core.String? resourcearn,
    $core.String? expectedrevisionid,
    $core.String? policy,
  }) {
    final result = create();
    if (confirmremoveselfresourceaccess != null)
      result.confirmremoveselfresourceaccess = confirmremoveselfresourceaccess;
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (expectedrevisionid != null)
      result.expectedrevisionid = expectedrevisionid;
    if (policy != null) result.policy = policy;
    return result;
  }

  PutResourcePolicyInput._();

  factory PutResourcePolicyInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutResourcePolicyInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutResourcePolicyInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOB(88224484, _omitFieldNames ? '' : 'confirmremoveselfresourceaccess')
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..aOS(456463970, _omitFieldNames ? '' : 'expectedrevisionid')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutResourcePolicyInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutResourcePolicyInput copyWith(
          void Function(PutResourcePolicyInput) updates) =>
      super.copyWith((message) => updates(message as PutResourcePolicyInput))
          as PutResourcePolicyInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutResourcePolicyInput create() => PutResourcePolicyInput._();
  @$core.override
  PutResourcePolicyInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutResourcePolicyInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutResourcePolicyInput>(create);
  static PutResourcePolicyInput? _defaultInstance;

  @$pb.TagNumber(88224484)
  $core.bool get confirmremoveselfresourceaccess => $_getBF(0);
  @$pb.TagNumber(88224484)
  set confirmremoveselfresourceaccess($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(88224484)
  $core.bool hasConfirmremoveselfresourceaccess() => $_has(0);
  @$pb.TagNumber(88224484)
  void clearConfirmremoveselfresourceaccess() => $_clearField(88224484);

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);

  @$pb.TagNumber(456463970)
  $core.String get expectedrevisionid => $_getSZ(2);
  @$pb.TagNumber(456463970)
  set expectedrevisionid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(456463970)
  $core.bool hasExpectedrevisionid() => $_has(2);
  @$pb.TagNumber(456463970)
  void clearExpectedrevisionid() => $_clearField(456463970);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(3);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(3, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(3);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class PutResourcePolicyOutput extends $pb.GeneratedMessage {
  factory PutResourcePolicyOutput({
    $core.String? revisionid,
  }) {
    final result = create();
    if (revisionid != null) result.revisionid = revisionid;
    return result;
  }

  PutResourcePolicyOutput._();

  factory PutResourcePolicyOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutResourcePolicyOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutResourcePolicyOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(499618182, _omitFieldNames ? '' : 'revisionid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutResourcePolicyOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutResourcePolicyOutput copyWith(
          void Function(PutResourcePolicyOutput) updates) =>
      super.copyWith((message) => updates(message as PutResourcePolicyOutput))
          as PutResourcePolicyOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutResourcePolicyOutput create() => PutResourcePolicyOutput._();
  @$core.override
  PutResourcePolicyOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutResourcePolicyOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutResourcePolicyOutput>(create);
  static PutResourcePolicyOutput? _defaultInstance;

  @$pb.TagNumber(499618182)
  $core.String get revisionid => $_getSZ(0);
  @$pb.TagNumber(499618182)
  set revisionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(499618182)
  $core.bool hasRevisionid() => $_has(0);
  @$pb.TagNumber(499618182)
  void clearRevisionid() => $_clearField(499618182);
}

class QueryInput extends $pb.GeneratedMessage {
  factory QueryInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.String? filterexpression,
    $core.String? indexname,
    $core.Iterable<$core.MapEntry<$core.String, Condition>>? queryfilter,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.String? projectionexpression,
    ConditionalOperator? conditionaloperator,
    $core.String? keyconditionexpression,
    $core.String? tablename,
    $core.Iterable<$core.MapEntry<$core.String, Condition>>? keyconditions,
    $core.Iterable<$core.String>? attributestoget,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        exclusivestartkey,
    $core.bool? scanindexforward,
    $core.int? limit,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
    Select? select,
    $core.bool? consistentread,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (filterexpression != null) result.filterexpression = filterexpression;
    if (indexname != null) result.indexname = indexname;
    if (queryfilter != null) result.queryfilter.addEntries(queryfilter);
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (projectionexpression != null)
      result.projectionexpression = projectionexpression;
    if (conditionaloperator != null)
      result.conditionaloperator = conditionaloperator;
    if (keyconditionexpression != null)
      result.keyconditionexpression = keyconditionexpression;
    if (tablename != null) result.tablename = tablename;
    if (keyconditions != null) result.keyconditions.addEntries(keyconditions);
    if (attributestoget != null) result.attributestoget.addAll(attributestoget);
    if (exclusivestartkey != null)
      result.exclusivestartkey.addEntries(exclusivestartkey);
    if (scanindexforward != null) result.scanindexforward = scanindexforward;
    if (limit != null) result.limit = limit;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    if (select != null) result.select = select;
    if (consistentread != null) result.consistentread = consistentread;
    return result;
  }

  QueryInput._();

  factory QueryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..aOS(68019610, _omitFieldNames ? '' : 'filterexpression')
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..m<$core.String, Condition>(
        107880942, _omitFieldNames ? '' : 'queryfilter',
        entryClassName: 'QueryInput.QueryfilterEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: Condition.create,
        valueDefaultOrMaker: Condition.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'QueryInput.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(150730243, _omitFieldNames ? '' : 'projectionexpression')
    ..aE<ConditionalOperator>(
        172066260, _omitFieldNames ? '' : 'conditionaloperator',
        enumValues: ConditionalOperator.values)
    ..aOS(218579280, _omitFieldNames ? '' : 'keyconditionexpression')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..m<$core.String, Condition>(
        307034879, _omitFieldNames ? '' : 'keyconditions',
        entryClassName: 'QueryInput.KeyconditionsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: Condition.create,
        valueDefaultOrMaker: Condition.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..pPS(311382592, _omitFieldNames ? '' : 'attributestoget')
    ..m<$core.String, AttributeValue>(
        348137297, _omitFieldNames ? '' : 'exclusivestartkey',
        entryClassName: 'QueryInput.ExclusivestartkeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOB(360066954, _omitFieldNames ? '' : 'scanindexforward')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'QueryInput.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<Select>(512305998, _omitFieldNames ? '' : 'select',
        enumValues: Select.values)
    ..aOB(531556994, _omitFieldNames ? '' : 'consistentread')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryInput copyWith(void Function(QueryInput) updates) =>
      super.copyWith((message) => updates(message as QueryInput)) as QueryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryInput create() => QueryInput._();
  @$core.override
  QueryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryInput>(create);
  static QueryInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(68019610)
  $core.String get filterexpression => $_getSZ(1);
  @$pb.TagNumber(68019610)
  set filterexpression($core.String value) => $_setString(1, value);
  @$pb.TagNumber(68019610)
  $core.bool hasFilterexpression() => $_has(1);
  @$pb.TagNumber(68019610)
  void clearFilterexpression() => $_clearField(68019610);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(2);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(2);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(107880942)
  $pb.PbMap<$core.String, Condition> get queryfilter => $_getMap(3);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(4);

  @$pb.TagNumber(150730243)
  $core.String get projectionexpression => $_getSZ(5);
  @$pb.TagNumber(150730243)
  set projectionexpression($core.String value) => $_setString(5, value);
  @$pb.TagNumber(150730243)
  $core.bool hasProjectionexpression() => $_has(5);
  @$pb.TagNumber(150730243)
  void clearProjectionexpression() => $_clearField(150730243);

  @$pb.TagNumber(172066260)
  ConditionalOperator get conditionaloperator => $_getN(6);
  @$pb.TagNumber(172066260)
  set conditionaloperator(ConditionalOperator value) =>
      $_setField(172066260, value);
  @$pb.TagNumber(172066260)
  $core.bool hasConditionaloperator() => $_has(6);
  @$pb.TagNumber(172066260)
  void clearConditionaloperator() => $_clearField(172066260);

  @$pb.TagNumber(218579280)
  $core.String get keyconditionexpression => $_getSZ(7);
  @$pb.TagNumber(218579280)
  set keyconditionexpression($core.String value) => $_setString(7, value);
  @$pb.TagNumber(218579280)
  $core.bool hasKeyconditionexpression() => $_has(7);
  @$pb.TagNumber(218579280)
  void clearKeyconditionexpression() => $_clearField(218579280);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(8);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(8, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(8);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(307034879)
  $pb.PbMap<$core.String, Condition> get keyconditions => $_getMap(9);

  @$pb.TagNumber(311382592)
  $pb.PbList<$core.String> get attributestoget => $_getList(10);

  @$pb.TagNumber(348137297)
  $pb.PbMap<$core.String, AttributeValue> get exclusivestartkey => $_getMap(11);

  @$pb.TagNumber(360066954)
  $core.bool get scanindexforward => $_getBF(12);
  @$pb.TagNumber(360066954)
  set scanindexforward($core.bool value) => $_setBool(12, value);
  @$pb.TagNumber(360066954)
  $core.bool hasScanindexforward() => $_has(12);
  @$pb.TagNumber(360066954)
  void clearScanindexforward() => $_clearField(360066954);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(13);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(13, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(13);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(14);

  @$pb.TagNumber(512305998)
  Select get select => $_getN(15);
  @$pb.TagNumber(512305998)
  set select(Select value) => $_setField(512305998, value);
  @$pb.TagNumber(512305998)
  $core.bool hasSelect() => $_has(15);
  @$pb.TagNumber(512305998)
  void clearSelect() => $_clearField(512305998);

  @$pb.TagNumber(531556994)
  $core.bool get consistentread => $_getBF(16);
  @$pb.TagNumber(531556994)
  set consistentread($core.bool value) => $_setBool(16, value);
  @$pb.TagNumber(531556994)
  $core.bool hasConsistentread() => $_has(16);
  @$pb.TagNumber(531556994)
  void clearConsistentread() => $_clearField(531556994);
}

class QueryOutput extends $pb.GeneratedMessage {
  factory QueryOutput({
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? items,
    $core.int? count,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        lastevaluatedkey,
    ConsumedCapacity? consumedcapacity,
    $core.int? scannedcount,
  }) {
    final result = create();
    if (items != null) result.items.addEntries(items);
    if (count != null) result.count = count;
    if (lastevaluatedkey != null)
      result.lastevaluatedkey.addEntries(lastevaluatedkey);
    if (consumedcapacity != null) result.consumedcapacity = consumedcapacity;
    if (scannedcount != null) result.scannedcount = scannedcount;
    return result;
  }

  QueryOutput._();

  factory QueryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, AttributeValue>(3553328, _omitFieldNames ? '' : 'items',
        entryClassName: 'QueryOutput.ItemsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aI(31963285, _omitFieldNames ? '' : 'count')
    ..m<$core.String, AttributeValue>(
        54319830, _omitFieldNames ? '' : 'lastevaluatedkey',
        entryClassName: 'QueryOutput.LastevaluatedkeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..aI(531161315, _omitFieldNames ? '' : 'scannedcount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryOutput copyWith(void Function(QueryOutput) updates) =>
      super.copyWith((message) => updates(message as QueryOutput))
          as QueryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryOutput create() => QueryOutput._();
  @$core.override
  QueryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryOutput>(create);
  static QueryOutput? _defaultInstance;

  @$pb.TagNumber(3553328)
  $pb.PbMap<$core.String, AttributeValue> get items => $_getMap(0);

  @$pb.TagNumber(31963285)
  $core.int get count => $_getIZ(1);
  @$pb.TagNumber(31963285)
  set count($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(31963285)
  $core.bool hasCount() => $_has(1);
  @$pb.TagNumber(31963285)
  void clearCount() => $_clearField(31963285);

  @$pb.TagNumber(54319830)
  $pb.PbMap<$core.String, AttributeValue> get lastevaluatedkey => $_getMap(2);

  @$pb.TagNumber(449336620)
  ConsumedCapacity get consumedcapacity => $_getN(3);
  @$pb.TagNumber(449336620)
  set consumedcapacity(ConsumedCapacity value) => $_setField(449336620, value);
  @$pb.TagNumber(449336620)
  $core.bool hasConsumedcapacity() => $_has(3);
  @$pb.TagNumber(449336620)
  void clearConsumedcapacity() => $_clearField(449336620);
  @$pb.TagNumber(449336620)
  ConsumedCapacity ensureConsumedcapacity() => $_ensure(3);

  @$pb.TagNumber(531161315)
  $core.int get scannedcount => $_getIZ(4);
  @$pb.TagNumber(531161315)
  set scannedcount($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(531161315)
  $core.bool hasScannedcount() => $_has(4);
  @$pb.TagNumber(531161315)
  void clearScannedcount() => $_clearField(531161315);
}

class Replica extends $pb.GeneratedMessage {
  factory Replica({
    $core.String? regionname,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    return result;
  }

  Replica._();

  factory Replica.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Replica.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Replica',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Replica clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Replica copyWith(void Function(Replica) updates) =>
      super.copyWith((message) => updates(message as Replica)) as Replica;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Replica create() => Replica._();
  @$core.override
  Replica createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Replica getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Replica>(create);
  static Replica? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);
}

class ReplicaAlreadyExistsException extends $pb.GeneratedMessage {
  factory ReplicaAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ReplicaAlreadyExistsException._();

  factory ReplicaAlreadyExistsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaAlreadyExistsException copyWith(
          void Function(ReplicaAlreadyExistsException) updates) =>
      super.copyWith(
              (message) => updates(message as ReplicaAlreadyExistsException))
          as ReplicaAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaAlreadyExistsException create() =>
      ReplicaAlreadyExistsException._();
  @$core.override
  ReplicaAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaAlreadyExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaAlreadyExistsException>(create);
  static ReplicaAlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ReplicaAutoScalingDescription extends $pb.GeneratedMessage {
  factory ReplicaAutoScalingDescription({
    $core.String? regionname,
    AutoScalingSettingsDescription?
        replicaprovisionedreadcapacityautoscalingsettings,
    AutoScalingSettingsDescription?
        replicaprovisionedwritecapacityautoscalingsettings,
    $core.Iterable<ReplicaGlobalSecondaryIndexAutoScalingDescription>?
        globalsecondaryindexes,
    ReplicaStatus? replicastatus,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    if (replicaprovisionedreadcapacityautoscalingsettings != null)
      result.replicaprovisionedreadcapacityautoscalingsettings =
          replicaprovisionedreadcapacityautoscalingsettings;
    if (replicaprovisionedwritecapacityautoscalingsettings != null)
      result.replicaprovisionedwritecapacityautoscalingsettings =
          replicaprovisionedwritecapacityautoscalingsettings;
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addAll(globalsecondaryindexes);
    if (replicastatus != null) result.replicastatus = replicastatus;
    return result;
  }

  ReplicaAutoScalingDescription._();

  factory ReplicaAutoScalingDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaAutoScalingDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaAutoScalingDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..aOM<AutoScalingSettingsDescription>(
        125131885,
        _omitFieldNames
            ? ''
            : 'replicaprovisionedreadcapacityautoscalingsettings',
        subBuilder: AutoScalingSettingsDescription.create)
    ..aOM<AutoScalingSettingsDescription>(
        293841796,
        _omitFieldNames
            ? ''
            : 'replicaprovisionedwritecapacityautoscalingsettings',
        subBuilder: AutoScalingSettingsDescription.create)
    ..pPM<ReplicaGlobalSecondaryIndexAutoScalingDescription>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        subBuilder: ReplicaGlobalSecondaryIndexAutoScalingDescription.create)
    ..aE<ReplicaStatus>(466739730, _omitFieldNames ? '' : 'replicastatus',
        enumValues: ReplicaStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaAutoScalingDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaAutoScalingDescription copyWith(
          void Function(ReplicaAutoScalingDescription) updates) =>
      super.copyWith(
              (message) => updates(message as ReplicaAutoScalingDescription))
          as ReplicaAutoScalingDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaAutoScalingDescription create() =>
      ReplicaAutoScalingDescription._();
  @$core.override
  ReplicaAutoScalingDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaAutoScalingDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaAutoScalingDescription>(create);
  static ReplicaAutoScalingDescription? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(125131885)
  AutoScalingSettingsDescription
      get replicaprovisionedreadcapacityautoscalingsettings => $_getN(1);
  @$pb.TagNumber(125131885)
  set replicaprovisionedreadcapacityautoscalingsettings(
          AutoScalingSettingsDescription value) =>
      $_setField(125131885, value);
  @$pb.TagNumber(125131885)
  $core.bool hasReplicaprovisionedreadcapacityautoscalingsettings() => $_has(1);
  @$pb.TagNumber(125131885)
  void clearReplicaprovisionedreadcapacityautoscalingsettings() =>
      $_clearField(125131885);
  @$pb.TagNumber(125131885)
  AutoScalingSettingsDescription
      ensureReplicaprovisionedreadcapacityautoscalingsettings() => $_ensure(1);

  @$pb.TagNumber(293841796)
  AutoScalingSettingsDescription
      get replicaprovisionedwritecapacityautoscalingsettings => $_getN(2);
  @$pb.TagNumber(293841796)
  set replicaprovisionedwritecapacityautoscalingsettings(
          AutoScalingSettingsDescription value) =>
      $_setField(293841796, value);
  @$pb.TagNumber(293841796)
  $core.bool hasReplicaprovisionedwritecapacityautoscalingsettings() =>
      $_has(2);
  @$pb.TagNumber(293841796)
  void clearReplicaprovisionedwritecapacityautoscalingsettings() =>
      $_clearField(293841796);
  @$pb.TagNumber(293841796)
  AutoScalingSettingsDescription
      ensureReplicaprovisionedwritecapacityautoscalingsettings() => $_ensure(2);

  @$pb.TagNumber(409156905)
  $pb.PbList<ReplicaGlobalSecondaryIndexAutoScalingDescription>
      get globalsecondaryindexes => $_getList(3);

  @$pb.TagNumber(466739730)
  ReplicaStatus get replicastatus => $_getN(4);
  @$pb.TagNumber(466739730)
  set replicastatus(ReplicaStatus value) => $_setField(466739730, value);
  @$pb.TagNumber(466739730)
  $core.bool hasReplicastatus() => $_has(4);
  @$pb.TagNumber(466739730)
  void clearReplicastatus() => $_clearField(466739730);
}

class ReplicaAutoScalingUpdate extends $pb.GeneratedMessage {
  factory ReplicaAutoScalingUpdate({
    $core.Iterable<ReplicaGlobalSecondaryIndexAutoScalingUpdate>?
        replicaglobalsecondaryindexupdates,
    $core.String? regionname,
    AutoScalingSettingsUpdate? replicaprovisionedreadcapacityautoscalingupdate,
  }) {
    final result = create();
    if (replicaglobalsecondaryindexupdates != null)
      result.replicaglobalsecondaryindexupdates
          .addAll(replicaglobalsecondaryindexupdates);
    if (regionname != null) result.regionname = regionname;
    if (replicaprovisionedreadcapacityautoscalingupdate != null)
      result.replicaprovisionedreadcapacityautoscalingupdate =
          replicaprovisionedreadcapacityautoscalingupdate;
    return result;
  }

  ReplicaAutoScalingUpdate._();

  factory ReplicaAutoScalingUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaAutoScalingUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaAutoScalingUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<ReplicaGlobalSecondaryIndexAutoScalingUpdate>(
        16441201, _omitFieldNames ? '' : 'replicaglobalsecondaryindexupdates',
        subBuilder: ReplicaGlobalSecondaryIndexAutoScalingUpdate.create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..aOM<AutoScalingSettingsUpdate>(
        262279073,
        _omitFieldNames
            ? ''
            : 'replicaprovisionedreadcapacityautoscalingupdate',
        subBuilder: AutoScalingSettingsUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaAutoScalingUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaAutoScalingUpdate copyWith(
          void Function(ReplicaAutoScalingUpdate) updates) =>
      super.copyWith((message) => updates(message as ReplicaAutoScalingUpdate))
          as ReplicaAutoScalingUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaAutoScalingUpdate create() => ReplicaAutoScalingUpdate._();
  @$core.override
  ReplicaAutoScalingUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaAutoScalingUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaAutoScalingUpdate>(create);
  static ReplicaAutoScalingUpdate? _defaultInstance;

  @$pb.TagNumber(16441201)
  $pb.PbList<ReplicaGlobalSecondaryIndexAutoScalingUpdate>
      get replicaglobalsecondaryindexupdates => $_getList(0);

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(1);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(1);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(262279073)
  AutoScalingSettingsUpdate
      get replicaprovisionedreadcapacityautoscalingupdate => $_getN(2);
  @$pb.TagNumber(262279073)
  set replicaprovisionedreadcapacityautoscalingupdate(
          AutoScalingSettingsUpdate value) =>
      $_setField(262279073, value);
  @$pb.TagNumber(262279073)
  $core.bool hasReplicaprovisionedreadcapacityautoscalingupdate() => $_has(2);
  @$pb.TagNumber(262279073)
  void clearReplicaprovisionedreadcapacityautoscalingupdate() =>
      $_clearField(262279073);
  @$pb.TagNumber(262279073)
  AutoScalingSettingsUpdate
      ensureReplicaprovisionedreadcapacityautoscalingupdate() => $_ensure(2);
}

class ReplicaDescription extends $pb.GeneratedMessage {
  factory ReplicaDescription({
    GlobalTableSettingsReplicationMode? globaltablesettingsreplicationmode,
    $core.String? replicainaccessibledatetime,
    $core.String? regionname,
    TableWarmThroughputDescription? warmthroughput,
    OnDemandThroughputOverride? ondemandthroughputoverride,
    $core.Iterable<ReplicaGlobalSecondaryIndexDescription>?
        globalsecondaryindexes,
    ProvisionedThroughputOverride? provisionedthroughputoverride,
    $core.String? replicastatusdescription,
    $core.String? replicastatuspercentprogress,
    ReplicaStatus? replicastatus,
    TableClassSummary? replicatableclasssummary,
    $core.String? kmsmasterkeyid,
  }) {
    final result = create();
    if (globaltablesettingsreplicationmode != null)
      result.globaltablesettingsreplicationmode =
          globaltablesettingsreplicationmode;
    if (replicainaccessibledatetime != null)
      result.replicainaccessibledatetime = replicainaccessibledatetime;
    if (regionname != null) result.regionname = regionname;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (ondemandthroughputoverride != null)
      result.ondemandthroughputoverride = ondemandthroughputoverride;
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addAll(globalsecondaryindexes);
    if (provisionedthroughputoverride != null)
      result.provisionedthroughputoverride = provisionedthroughputoverride;
    if (replicastatusdescription != null)
      result.replicastatusdescription = replicastatusdescription;
    if (replicastatuspercentprogress != null)
      result.replicastatuspercentprogress = replicastatuspercentprogress;
    if (replicastatus != null) result.replicastatus = replicastatus;
    if (replicatableclasssummary != null)
      result.replicatableclasssummary = replicatableclasssummary;
    if (kmsmasterkeyid != null) result.kmsmasterkeyid = kmsmasterkeyid;
    return result;
  }

  ReplicaDescription._();

  factory ReplicaDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<GlobalTableSettingsReplicationMode>(
        10446577, _omitFieldNames ? '' : 'globaltablesettingsreplicationmode',
        enumValues: GlobalTableSettingsReplicationMode.values)
    ..aOS(20060608, _omitFieldNames ? '' : 'replicainaccessibledatetime')
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..aOM<TableWarmThroughputDescription>(
        290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: TableWarmThroughputDescription.create)
    ..aOM<OnDemandThroughputOverride>(
        317165234, _omitFieldNames ? '' : 'ondemandthroughputoverride',
        subBuilder: OnDemandThroughputOverride.create)
    ..pPM<ReplicaGlobalSecondaryIndexDescription>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        subBuilder: ReplicaGlobalSecondaryIndexDescription.create)
    ..aOM<ProvisionedThroughputOverride>(
        413332116, _omitFieldNames ? '' : 'provisionedthroughputoverride',
        subBuilder: ProvisionedThroughputOverride.create)
    ..aOS(416683638, _omitFieldNames ? '' : 'replicastatusdescription')
    ..aOS(456382130, _omitFieldNames ? '' : 'replicastatuspercentprogress')
    ..aE<ReplicaStatus>(466739730, _omitFieldNames ? '' : 'replicastatus',
        enumValues: ReplicaStatus.values)
    ..aOM<TableClassSummary>(
        519483614, _omitFieldNames ? '' : 'replicatableclasssummary',
        subBuilder: TableClassSummary.create)
    ..aOS(521459443, _omitFieldNames ? '' : 'kmsmasterkeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaDescription copyWith(void Function(ReplicaDescription) updates) =>
      super.copyWith((message) => updates(message as ReplicaDescription))
          as ReplicaDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaDescription create() => ReplicaDescription._();
  @$core.override
  ReplicaDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaDescription>(create);
  static ReplicaDescription? _defaultInstance;

  @$pb.TagNumber(10446577)
  GlobalTableSettingsReplicationMode get globaltablesettingsreplicationmode =>
      $_getN(0);
  @$pb.TagNumber(10446577)
  set globaltablesettingsreplicationmode(
          GlobalTableSettingsReplicationMode value) =>
      $_setField(10446577, value);
  @$pb.TagNumber(10446577)
  $core.bool hasGlobaltablesettingsreplicationmode() => $_has(0);
  @$pb.TagNumber(10446577)
  void clearGlobaltablesettingsreplicationmode() => $_clearField(10446577);

  @$pb.TagNumber(20060608)
  $core.String get replicainaccessibledatetime => $_getSZ(1);
  @$pb.TagNumber(20060608)
  set replicainaccessibledatetime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(20060608)
  $core.bool hasReplicainaccessibledatetime() => $_has(1);
  @$pb.TagNumber(20060608)
  void clearReplicainaccessibledatetime() => $_clearField(20060608);

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(2);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(2);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(290598659)
  TableWarmThroughputDescription get warmthroughput => $_getN(3);
  @$pb.TagNumber(290598659)
  set warmthroughput(TableWarmThroughputDescription value) =>
      $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(3);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  TableWarmThroughputDescription ensureWarmthroughput() => $_ensure(3);

  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride get ondemandthroughputoverride => $_getN(4);
  @$pb.TagNumber(317165234)
  set ondemandthroughputoverride(OnDemandThroughputOverride value) =>
      $_setField(317165234, value);
  @$pb.TagNumber(317165234)
  $core.bool hasOndemandthroughputoverride() => $_has(4);
  @$pb.TagNumber(317165234)
  void clearOndemandthroughputoverride() => $_clearField(317165234);
  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride ensureOndemandthroughputoverride() => $_ensure(4);

  @$pb.TagNumber(409156905)
  $pb.PbList<ReplicaGlobalSecondaryIndexDescription>
      get globalsecondaryindexes => $_getList(5);

  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride get provisionedthroughputoverride => $_getN(6);
  @$pb.TagNumber(413332116)
  set provisionedthroughputoverride(ProvisionedThroughputOverride value) =>
      $_setField(413332116, value);
  @$pb.TagNumber(413332116)
  $core.bool hasProvisionedthroughputoverride() => $_has(6);
  @$pb.TagNumber(413332116)
  void clearProvisionedthroughputoverride() => $_clearField(413332116);
  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride ensureProvisionedthroughputoverride() =>
      $_ensure(6);

  @$pb.TagNumber(416683638)
  $core.String get replicastatusdescription => $_getSZ(7);
  @$pb.TagNumber(416683638)
  set replicastatusdescription($core.String value) => $_setString(7, value);
  @$pb.TagNumber(416683638)
  $core.bool hasReplicastatusdescription() => $_has(7);
  @$pb.TagNumber(416683638)
  void clearReplicastatusdescription() => $_clearField(416683638);

  @$pb.TagNumber(456382130)
  $core.String get replicastatuspercentprogress => $_getSZ(8);
  @$pb.TagNumber(456382130)
  set replicastatuspercentprogress($core.String value) => $_setString(8, value);
  @$pb.TagNumber(456382130)
  $core.bool hasReplicastatuspercentprogress() => $_has(8);
  @$pb.TagNumber(456382130)
  void clearReplicastatuspercentprogress() => $_clearField(456382130);

  @$pb.TagNumber(466739730)
  ReplicaStatus get replicastatus => $_getN(9);
  @$pb.TagNumber(466739730)
  set replicastatus(ReplicaStatus value) => $_setField(466739730, value);
  @$pb.TagNumber(466739730)
  $core.bool hasReplicastatus() => $_has(9);
  @$pb.TagNumber(466739730)
  void clearReplicastatus() => $_clearField(466739730);

  @$pb.TagNumber(519483614)
  TableClassSummary get replicatableclasssummary => $_getN(10);
  @$pb.TagNumber(519483614)
  set replicatableclasssummary(TableClassSummary value) =>
      $_setField(519483614, value);
  @$pb.TagNumber(519483614)
  $core.bool hasReplicatableclasssummary() => $_has(10);
  @$pb.TagNumber(519483614)
  void clearReplicatableclasssummary() => $_clearField(519483614);
  @$pb.TagNumber(519483614)
  TableClassSummary ensureReplicatableclasssummary() => $_ensure(10);

  @$pb.TagNumber(521459443)
  $core.String get kmsmasterkeyid => $_getSZ(11);
  @$pb.TagNumber(521459443)
  set kmsmasterkeyid($core.String value) => $_setString(11, value);
  @$pb.TagNumber(521459443)
  $core.bool hasKmsmasterkeyid() => $_has(11);
  @$pb.TagNumber(521459443)
  void clearKmsmasterkeyid() => $_clearField(521459443);
}

class ReplicaGlobalSecondaryIndex extends $pb.GeneratedMessage {
  factory ReplicaGlobalSecondaryIndex({
    $core.String? indexname,
    OnDemandThroughputOverride? ondemandthroughputoverride,
    ProvisionedThroughputOverride? provisionedthroughputoverride,
  }) {
    final result = create();
    if (indexname != null) result.indexname = indexname;
    if (ondemandthroughputoverride != null)
      result.ondemandthroughputoverride = ondemandthroughputoverride;
    if (provisionedthroughputoverride != null)
      result.provisionedthroughputoverride = provisionedthroughputoverride;
    return result;
  }

  ReplicaGlobalSecondaryIndex._();

  factory ReplicaGlobalSecondaryIndex.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaGlobalSecondaryIndex.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaGlobalSecondaryIndex',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<OnDemandThroughputOverride>(
        317165234, _omitFieldNames ? '' : 'ondemandthroughputoverride',
        subBuilder: OnDemandThroughputOverride.create)
    ..aOM<ProvisionedThroughputOverride>(
        413332116, _omitFieldNames ? '' : 'provisionedthroughputoverride',
        subBuilder: ProvisionedThroughputOverride.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndex clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndex copyWith(
          void Function(ReplicaGlobalSecondaryIndex) updates) =>
      super.copyWith(
              (message) => updates(message as ReplicaGlobalSecondaryIndex))
          as ReplicaGlobalSecondaryIndex;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndex create() =>
      ReplicaGlobalSecondaryIndex._();
  @$core.override
  ReplicaGlobalSecondaryIndex createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndex getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaGlobalSecondaryIndex>(create);
  static ReplicaGlobalSecondaryIndex? _defaultInstance;

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(0);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(0);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride get ondemandthroughputoverride => $_getN(1);
  @$pb.TagNumber(317165234)
  set ondemandthroughputoverride(OnDemandThroughputOverride value) =>
      $_setField(317165234, value);
  @$pb.TagNumber(317165234)
  $core.bool hasOndemandthroughputoverride() => $_has(1);
  @$pb.TagNumber(317165234)
  void clearOndemandthroughputoverride() => $_clearField(317165234);
  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride ensureOndemandthroughputoverride() => $_ensure(1);

  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride get provisionedthroughputoverride => $_getN(2);
  @$pb.TagNumber(413332116)
  set provisionedthroughputoverride(ProvisionedThroughputOverride value) =>
      $_setField(413332116, value);
  @$pb.TagNumber(413332116)
  $core.bool hasProvisionedthroughputoverride() => $_has(2);
  @$pb.TagNumber(413332116)
  void clearProvisionedthroughputoverride() => $_clearField(413332116);
  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride ensureProvisionedthroughputoverride() =>
      $_ensure(2);
}

class ReplicaGlobalSecondaryIndexAutoScalingDescription
    extends $pb.GeneratedMessage {
  factory ReplicaGlobalSecondaryIndexAutoScalingDescription({
    AutoScalingSettingsDescription? provisionedreadcapacityautoscalingsettings,
    $core.String? indexname,
    AutoScalingSettingsDescription? provisionedwritecapacityautoscalingsettings,
    IndexStatus? indexstatus,
  }) {
    final result = create();
    if (provisionedreadcapacityautoscalingsettings != null)
      result.provisionedreadcapacityautoscalingsettings =
          provisionedreadcapacityautoscalingsettings;
    if (indexname != null) result.indexname = indexname;
    if (provisionedwritecapacityautoscalingsettings != null)
      result.provisionedwritecapacityautoscalingsettings =
          provisionedwritecapacityautoscalingsettings;
    if (indexstatus != null) result.indexstatus = indexstatus;
    return result;
  }

  ReplicaGlobalSecondaryIndexAutoScalingDescription._();

  factory ReplicaGlobalSecondaryIndexAutoScalingDescription.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaGlobalSecondaryIndexAutoScalingDescription.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames
          ? ''
          : 'ReplicaGlobalSecondaryIndexAutoScalingDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<AutoScalingSettingsDescription>(85617667,
        _omitFieldNames ? '' : 'provisionedreadcapacityautoscalingsettings',
        subBuilder: AutoScalingSettingsDescription.create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<AutoScalingSettingsDescription>(168409114,
        _omitFieldNames ? '' : 'provisionedwritecapacityautoscalingsettings',
        subBuilder: AutoScalingSettingsDescription.create)
    ..aE<IndexStatus>(364436830, _omitFieldNames ? '' : 'indexstatus',
        enumValues: IndexStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexAutoScalingDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexAutoScalingDescription copyWith(
          void Function(ReplicaGlobalSecondaryIndexAutoScalingDescription)
              updates) =>
      super.copyWith((message) => updates(
              message as ReplicaGlobalSecondaryIndexAutoScalingDescription))
          as ReplicaGlobalSecondaryIndexAutoScalingDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexAutoScalingDescription create() =>
      ReplicaGlobalSecondaryIndexAutoScalingDescription._();
  @$core.override
  ReplicaGlobalSecondaryIndexAutoScalingDescription createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexAutoScalingDescription getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ReplicaGlobalSecondaryIndexAutoScalingDescription>(create);
  static ReplicaGlobalSecondaryIndexAutoScalingDescription? _defaultInstance;

  @$pb.TagNumber(85617667)
  AutoScalingSettingsDescription
      get provisionedreadcapacityautoscalingsettings => $_getN(0);
  @$pb.TagNumber(85617667)
  set provisionedreadcapacityautoscalingsettings(
          AutoScalingSettingsDescription value) =>
      $_setField(85617667, value);
  @$pb.TagNumber(85617667)
  $core.bool hasProvisionedreadcapacityautoscalingsettings() => $_has(0);
  @$pb.TagNumber(85617667)
  void clearProvisionedreadcapacityautoscalingsettings() =>
      $_clearField(85617667);
  @$pb.TagNumber(85617667)
  AutoScalingSettingsDescription
      ensureProvisionedreadcapacityautoscalingsettings() => $_ensure(0);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(168409114)
  AutoScalingSettingsDescription
      get provisionedwritecapacityautoscalingsettings => $_getN(2);
  @$pb.TagNumber(168409114)
  set provisionedwritecapacityautoscalingsettings(
          AutoScalingSettingsDescription value) =>
      $_setField(168409114, value);
  @$pb.TagNumber(168409114)
  $core.bool hasProvisionedwritecapacityautoscalingsettings() => $_has(2);
  @$pb.TagNumber(168409114)
  void clearProvisionedwritecapacityautoscalingsettings() =>
      $_clearField(168409114);
  @$pb.TagNumber(168409114)
  AutoScalingSettingsDescription
      ensureProvisionedwritecapacityautoscalingsettings() => $_ensure(2);

  @$pb.TagNumber(364436830)
  IndexStatus get indexstatus => $_getN(3);
  @$pb.TagNumber(364436830)
  set indexstatus(IndexStatus value) => $_setField(364436830, value);
  @$pb.TagNumber(364436830)
  $core.bool hasIndexstatus() => $_has(3);
  @$pb.TagNumber(364436830)
  void clearIndexstatus() => $_clearField(364436830);
}

class ReplicaGlobalSecondaryIndexAutoScalingUpdate
    extends $pb.GeneratedMessage {
  factory ReplicaGlobalSecondaryIndexAutoScalingUpdate({
    AutoScalingSettingsUpdate? provisionedreadcapacityautoscalingupdate,
    $core.String? indexname,
  }) {
    final result = create();
    if (provisionedreadcapacityautoscalingupdate != null)
      result.provisionedreadcapacityautoscalingupdate =
          provisionedreadcapacityautoscalingupdate;
    if (indexname != null) result.indexname = indexname;
    return result;
  }

  ReplicaGlobalSecondaryIndexAutoScalingUpdate._();

  factory ReplicaGlobalSecondaryIndexAutoScalingUpdate.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaGlobalSecondaryIndexAutoScalingUpdate.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaGlobalSecondaryIndexAutoScalingUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<AutoScalingSettingsUpdate>(22287551,
        _omitFieldNames ? '' : 'provisionedreadcapacityautoscalingupdate',
        subBuilder: AutoScalingSettingsUpdate.create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexAutoScalingUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexAutoScalingUpdate copyWith(
          void Function(ReplicaGlobalSecondaryIndexAutoScalingUpdate)
              updates) =>
      super.copyWith((message) =>
              updates(message as ReplicaGlobalSecondaryIndexAutoScalingUpdate))
          as ReplicaGlobalSecondaryIndexAutoScalingUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexAutoScalingUpdate create() =>
      ReplicaGlobalSecondaryIndexAutoScalingUpdate._();
  @$core.override
  ReplicaGlobalSecondaryIndexAutoScalingUpdate createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexAutoScalingUpdate getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ReplicaGlobalSecondaryIndexAutoScalingUpdate>(create);
  static ReplicaGlobalSecondaryIndexAutoScalingUpdate? _defaultInstance;

  @$pb.TagNumber(22287551)
  AutoScalingSettingsUpdate get provisionedreadcapacityautoscalingupdate =>
      $_getN(0);
  @$pb.TagNumber(22287551)
  set provisionedreadcapacityautoscalingupdate(
          AutoScalingSettingsUpdate value) =>
      $_setField(22287551, value);
  @$pb.TagNumber(22287551)
  $core.bool hasProvisionedreadcapacityautoscalingupdate() => $_has(0);
  @$pb.TagNumber(22287551)
  void clearProvisionedreadcapacityautoscalingupdate() =>
      $_clearField(22287551);
  @$pb.TagNumber(22287551)
  AutoScalingSettingsUpdate ensureProvisionedreadcapacityautoscalingupdate() =>
      $_ensure(0);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);
}

class ReplicaGlobalSecondaryIndexDescription extends $pb.GeneratedMessage {
  factory ReplicaGlobalSecondaryIndexDescription({
    $core.String? indexname,
    GlobalSecondaryIndexWarmThroughputDescription? warmthroughput,
    OnDemandThroughputOverride? ondemandthroughputoverride,
    ProvisionedThroughputOverride? provisionedthroughputoverride,
  }) {
    final result = create();
    if (indexname != null) result.indexname = indexname;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (ondemandthroughputoverride != null)
      result.ondemandthroughputoverride = ondemandthroughputoverride;
    if (provisionedthroughputoverride != null)
      result.provisionedthroughputoverride = provisionedthroughputoverride;
    return result;
  }

  ReplicaGlobalSecondaryIndexDescription._();

  factory ReplicaGlobalSecondaryIndexDescription.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaGlobalSecondaryIndexDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaGlobalSecondaryIndexDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<GlobalSecondaryIndexWarmThroughputDescription>(
        290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: GlobalSecondaryIndexWarmThroughputDescription.create)
    ..aOM<OnDemandThroughputOverride>(
        317165234, _omitFieldNames ? '' : 'ondemandthroughputoverride',
        subBuilder: OnDemandThroughputOverride.create)
    ..aOM<ProvisionedThroughputOverride>(
        413332116, _omitFieldNames ? '' : 'provisionedthroughputoverride',
        subBuilder: ProvisionedThroughputOverride.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexDescription copyWith(
          void Function(ReplicaGlobalSecondaryIndexDescription) updates) =>
      super.copyWith((message) =>
              updates(message as ReplicaGlobalSecondaryIndexDescription))
          as ReplicaGlobalSecondaryIndexDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexDescription create() =>
      ReplicaGlobalSecondaryIndexDescription._();
  @$core.override
  ReplicaGlobalSecondaryIndexDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexDescription getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ReplicaGlobalSecondaryIndexDescription>(create);
  static ReplicaGlobalSecondaryIndexDescription? _defaultInstance;

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(0);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(0);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(290598659)
  GlobalSecondaryIndexWarmThroughputDescription get warmthroughput => $_getN(1);
  @$pb.TagNumber(290598659)
  set warmthroughput(GlobalSecondaryIndexWarmThroughputDescription value) =>
      $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(1);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  GlobalSecondaryIndexWarmThroughputDescription ensureWarmthroughput() =>
      $_ensure(1);

  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride get ondemandthroughputoverride => $_getN(2);
  @$pb.TagNumber(317165234)
  set ondemandthroughputoverride(OnDemandThroughputOverride value) =>
      $_setField(317165234, value);
  @$pb.TagNumber(317165234)
  $core.bool hasOndemandthroughputoverride() => $_has(2);
  @$pb.TagNumber(317165234)
  void clearOndemandthroughputoverride() => $_clearField(317165234);
  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride ensureOndemandthroughputoverride() => $_ensure(2);

  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride get provisionedthroughputoverride => $_getN(3);
  @$pb.TagNumber(413332116)
  set provisionedthroughputoverride(ProvisionedThroughputOverride value) =>
      $_setField(413332116, value);
  @$pb.TagNumber(413332116)
  $core.bool hasProvisionedthroughputoverride() => $_has(3);
  @$pb.TagNumber(413332116)
  void clearProvisionedthroughputoverride() => $_clearField(413332116);
  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride ensureProvisionedthroughputoverride() =>
      $_ensure(3);
}

class ReplicaGlobalSecondaryIndexSettingsDescription
    extends $pb.GeneratedMessage {
  factory ReplicaGlobalSecondaryIndexSettingsDescription({
    AutoScalingSettingsDescription? provisionedreadcapacityautoscalingsettings,
    $core.String? indexname,
    AutoScalingSettingsDescription? provisionedwritecapacityautoscalingsettings,
    $fixnum.Int64? provisionedwritecapacityunits,
    $fixnum.Int64? provisionedreadcapacityunits,
    IndexStatus? indexstatus,
  }) {
    final result = create();
    if (provisionedreadcapacityautoscalingsettings != null)
      result.provisionedreadcapacityautoscalingsettings =
          provisionedreadcapacityautoscalingsettings;
    if (indexname != null) result.indexname = indexname;
    if (provisionedwritecapacityautoscalingsettings != null)
      result.provisionedwritecapacityautoscalingsettings =
          provisionedwritecapacityautoscalingsettings;
    if (provisionedwritecapacityunits != null)
      result.provisionedwritecapacityunits = provisionedwritecapacityunits;
    if (provisionedreadcapacityunits != null)
      result.provisionedreadcapacityunits = provisionedreadcapacityunits;
    if (indexstatus != null) result.indexstatus = indexstatus;
    return result;
  }

  ReplicaGlobalSecondaryIndexSettingsDescription._();

  factory ReplicaGlobalSecondaryIndexSettingsDescription.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaGlobalSecondaryIndexSettingsDescription.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaGlobalSecondaryIndexSettingsDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<AutoScalingSettingsDescription>(85617667,
        _omitFieldNames ? '' : 'provisionedreadcapacityautoscalingsettings',
        subBuilder: AutoScalingSettingsDescription.create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<AutoScalingSettingsDescription>(168409114,
        _omitFieldNames ? '' : 'provisionedwritecapacityautoscalingsettings',
        subBuilder: AutoScalingSettingsDescription.create)
    ..aInt64(225881684, _omitFieldNames ? '' : 'provisionedwritecapacityunits')
    ..aInt64(350750021, _omitFieldNames ? '' : 'provisionedreadcapacityunits')
    ..aE<IndexStatus>(364436830, _omitFieldNames ? '' : 'indexstatus',
        enumValues: IndexStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexSettingsDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexSettingsDescription copyWith(
          void Function(ReplicaGlobalSecondaryIndexSettingsDescription)
              updates) =>
      super.copyWith((message) => updates(
              message as ReplicaGlobalSecondaryIndexSettingsDescription))
          as ReplicaGlobalSecondaryIndexSettingsDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexSettingsDescription create() =>
      ReplicaGlobalSecondaryIndexSettingsDescription._();
  @$core.override
  ReplicaGlobalSecondaryIndexSettingsDescription createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexSettingsDescription getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ReplicaGlobalSecondaryIndexSettingsDescription>(create);
  static ReplicaGlobalSecondaryIndexSettingsDescription? _defaultInstance;

  @$pb.TagNumber(85617667)
  AutoScalingSettingsDescription
      get provisionedreadcapacityautoscalingsettings => $_getN(0);
  @$pb.TagNumber(85617667)
  set provisionedreadcapacityautoscalingsettings(
          AutoScalingSettingsDescription value) =>
      $_setField(85617667, value);
  @$pb.TagNumber(85617667)
  $core.bool hasProvisionedreadcapacityautoscalingsettings() => $_has(0);
  @$pb.TagNumber(85617667)
  void clearProvisionedreadcapacityautoscalingsettings() =>
      $_clearField(85617667);
  @$pb.TagNumber(85617667)
  AutoScalingSettingsDescription
      ensureProvisionedreadcapacityautoscalingsettings() => $_ensure(0);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(168409114)
  AutoScalingSettingsDescription
      get provisionedwritecapacityautoscalingsettings => $_getN(2);
  @$pb.TagNumber(168409114)
  set provisionedwritecapacityautoscalingsettings(
          AutoScalingSettingsDescription value) =>
      $_setField(168409114, value);
  @$pb.TagNumber(168409114)
  $core.bool hasProvisionedwritecapacityautoscalingsettings() => $_has(2);
  @$pb.TagNumber(168409114)
  void clearProvisionedwritecapacityautoscalingsettings() =>
      $_clearField(168409114);
  @$pb.TagNumber(168409114)
  AutoScalingSettingsDescription
      ensureProvisionedwritecapacityautoscalingsettings() => $_ensure(2);

  @$pb.TagNumber(225881684)
  $fixnum.Int64 get provisionedwritecapacityunits => $_getI64(3);
  @$pb.TagNumber(225881684)
  set provisionedwritecapacityunits($fixnum.Int64 value) =>
      $_setInt64(3, value);
  @$pb.TagNumber(225881684)
  $core.bool hasProvisionedwritecapacityunits() => $_has(3);
  @$pb.TagNumber(225881684)
  void clearProvisionedwritecapacityunits() => $_clearField(225881684);

  @$pb.TagNumber(350750021)
  $fixnum.Int64 get provisionedreadcapacityunits => $_getI64(4);
  @$pb.TagNumber(350750021)
  set provisionedreadcapacityunits($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(350750021)
  $core.bool hasProvisionedreadcapacityunits() => $_has(4);
  @$pb.TagNumber(350750021)
  void clearProvisionedreadcapacityunits() => $_clearField(350750021);

  @$pb.TagNumber(364436830)
  IndexStatus get indexstatus => $_getN(5);
  @$pb.TagNumber(364436830)
  set indexstatus(IndexStatus value) => $_setField(364436830, value);
  @$pb.TagNumber(364436830)
  $core.bool hasIndexstatus() => $_has(5);
  @$pb.TagNumber(364436830)
  void clearIndexstatus() => $_clearField(364436830);
}

class ReplicaGlobalSecondaryIndexSettingsUpdate extends $pb.GeneratedMessage {
  factory ReplicaGlobalSecondaryIndexSettingsUpdate({
    AutoScalingSettingsUpdate? provisionedreadcapacityautoscalingsettingsupdate,
    $core.String? indexname,
    $fixnum.Int64? provisionedreadcapacityunits,
  }) {
    final result = create();
    if (provisionedreadcapacityautoscalingsettingsupdate != null)
      result.provisionedreadcapacityautoscalingsettingsupdate =
          provisionedreadcapacityautoscalingsettingsupdate;
    if (indexname != null) result.indexname = indexname;
    if (provisionedreadcapacityunits != null)
      result.provisionedreadcapacityunits = provisionedreadcapacityunits;
    return result;
  }

  ReplicaGlobalSecondaryIndexSettingsUpdate._();

  factory ReplicaGlobalSecondaryIndexSettingsUpdate.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaGlobalSecondaryIndexSettingsUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaGlobalSecondaryIndexSettingsUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<AutoScalingSettingsUpdate>(
        57419900,
        _omitFieldNames
            ? ''
            : 'provisionedreadcapacityautoscalingsettingsupdate',
        subBuilder: AutoScalingSettingsUpdate.create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aInt64(350750021, _omitFieldNames ? '' : 'provisionedreadcapacityunits')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexSettingsUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaGlobalSecondaryIndexSettingsUpdate copyWith(
          void Function(ReplicaGlobalSecondaryIndexSettingsUpdate) updates) =>
      super.copyWith((message) =>
              updates(message as ReplicaGlobalSecondaryIndexSettingsUpdate))
          as ReplicaGlobalSecondaryIndexSettingsUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexSettingsUpdate create() =>
      ReplicaGlobalSecondaryIndexSettingsUpdate._();
  @$core.override
  ReplicaGlobalSecondaryIndexSettingsUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaGlobalSecondaryIndexSettingsUpdate getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ReplicaGlobalSecondaryIndexSettingsUpdate>(create);
  static ReplicaGlobalSecondaryIndexSettingsUpdate? _defaultInstance;

  @$pb.TagNumber(57419900)
  AutoScalingSettingsUpdate
      get provisionedreadcapacityautoscalingsettingsupdate => $_getN(0);
  @$pb.TagNumber(57419900)
  set provisionedreadcapacityautoscalingsettingsupdate(
          AutoScalingSettingsUpdate value) =>
      $_setField(57419900, value);
  @$pb.TagNumber(57419900)
  $core.bool hasProvisionedreadcapacityautoscalingsettingsupdate() => $_has(0);
  @$pb.TagNumber(57419900)
  void clearProvisionedreadcapacityautoscalingsettingsupdate() =>
      $_clearField(57419900);
  @$pb.TagNumber(57419900)
  AutoScalingSettingsUpdate
      ensureProvisionedreadcapacityautoscalingsettingsupdate() => $_ensure(0);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(350750021)
  $fixnum.Int64 get provisionedreadcapacityunits => $_getI64(2);
  @$pb.TagNumber(350750021)
  set provisionedreadcapacityunits($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(350750021)
  $core.bool hasProvisionedreadcapacityunits() => $_has(2);
  @$pb.TagNumber(350750021)
  void clearProvisionedreadcapacityunits() => $_clearField(350750021);
}

class ReplicaNotFoundException extends $pb.GeneratedMessage {
  factory ReplicaNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ReplicaNotFoundException._();

  factory ReplicaNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaNotFoundException copyWith(
          void Function(ReplicaNotFoundException) updates) =>
      super.copyWith((message) => updates(message as ReplicaNotFoundException))
          as ReplicaNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaNotFoundException create() => ReplicaNotFoundException._();
  @$core.override
  ReplicaNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaNotFoundException>(create);
  static ReplicaNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ReplicaSettingsDescription extends $pb.GeneratedMessage {
  factory ReplicaSettingsDescription({
    $fixnum.Int64? replicaprovisionedreadcapacityunits,
    $core.String? regionname,
    AutoScalingSettingsDescription?
        replicaprovisionedreadcapacityautoscalingsettings,
    AutoScalingSettingsDescription?
        replicaprovisionedwritecapacityautoscalingsettings,
    $fixnum.Int64? replicaprovisionedwritecapacityunits,
    BillingModeSummary? replicabillingmodesummary,
    ReplicaStatus? replicastatus,
    TableClassSummary? replicatableclasssummary,
    $core.Iterable<ReplicaGlobalSecondaryIndexSettingsDescription>?
        replicaglobalsecondaryindexsettings,
  }) {
    final result = create();
    if (replicaprovisionedreadcapacityunits != null)
      result.replicaprovisionedreadcapacityunits =
          replicaprovisionedreadcapacityunits;
    if (regionname != null) result.regionname = regionname;
    if (replicaprovisionedreadcapacityautoscalingsettings != null)
      result.replicaprovisionedreadcapacityautoscalingsettings =
          replicaprovisionedreadcapacityautoscalingsettings;
    if (replicaprovisionedwritecapacityautoscalingsettings != null)
      result.replicaprovisionedwritecapacityautoscalingsettings =
          replicaprovisionedwritecapacityautoscalingsettings;
    if (replicaprovisionedwritecapacityunits != null)
      result.replicaprovisionedwritecapacityunits =
          replicaprovisionedwritecapacityunits;
    if (replicabillingmodesummary != null)
      result.replicabillingmodesummary = replicabillingmodesummary;
    if (replicastatus != null) result.replicastatus = replicastatus;
    if (replicatableclasssummary != null)
      result.replicatableclasssummary = replicatableclasssummary;
    if (replicaglobalsecondaryindexsettings != null)
      result.replicaglobalsecondaryindexsettings
          .addAll(replicaglobalsecondaryindexsettings);
    return result;
  }

  ReplicaSettingsDescription._();

  factory ReplicaSettingsDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaSettingsDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaSettingsDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(
        82081083, _omitFieldNames ? '' : 'replicaprovisionedreadcapacityunits')
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..aOM<AutoScalingSettingsDescription>(
        125131885,
        _omitFieldNames
            ? ''
            : 'replicaprovisionedreadcapacityautoscalingsettings',
        subBuilder: AutoScalingSettingsDescription.create)
    ..aOM<AutoScalingSettingsDescription>(
        293841796,
        _omitFieldNames
            ? ''
            : 'replicaprovisionedwritecapacityautoscalingsettings',
        subBuilder: AutoScalingSettingsDescription.create)
    ..aInt64(356738858,
        _omitFieldNames ? '' : 'replicaprovisionedwritecapacityunits')
    ..aOM<BillingModeSummary>(
        456333796, _omitFieldNames ? '' : 'replicabillingmodesummary',
        subBuilder: BillingModeSummary.create)
    ..aE<ReplicaStatus>(466739730, _omitFieldNames ? '' : 'replicastatus',
        enumValues: ReplicaStatus.values)
    ..aOM<TableClassSummary>(
        519483614, _omitFieldNames ? '' : 'replicatableclasssummary',
        subBuilder: TableClassSummary.create)
    ..pPM<ReplicaGlobalSecondaryIndexSettingsDescription>(
        521234416, _omitFieldNames ? '' : 'replicaglobalsecondaryindexsettings',
        subBuilder: ReplicaGlobalSecondaryIndexSettingsDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaSettingsDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaSettingsDescription copyWith(
          void Function(ReplicaSettingsDescription) updates) =>
      super.copyWith(
              (message) => updates(message as ReplicaSettingsDescription))
          as ReplicaSettingsDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaSettingsDescription create() => ReplicaSettingsDescription._();
  @$core.override
  ReplicaSettingsDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaSettingsDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaSettingsDescription>(create);
  static ReplicaSettingsDescription? _defaultInstance;

  @$pb.TagNumber(82081083)
  $fixnum.Int64 get replicaprovisionedreadcapacityunits => $_getI64(0);
  @$pb.TagNumber(82081083)
  set replicaprovisionedreadcapacityunits($fixnum.Int64 value) =>
      $_setInt64(0, value);
  @$pb.TagNumber(82081083)
  $core.bool hasReplicaprovisionedreadcapacityunits() => $_has(0);
  @$pb.TagNumber(82081083)
  void clearReplicaprovisionedreadcapacityunits() => $_clearField(82081083);

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(1);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(1);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(125131885)
  AutoScalingSettingsDescription
      get replicaprovisionedreadcapacityautoscalingsettings => $_getN(2);
  @$pb.TagNumber(125131885)
  set replicaprovisionedreadcapacityautoscalingsettings(
          AutoScalingSettingsDescription value) =>
      $_setField(125131885, value);
  @$pb.TagNumber(125131885)
  $core.bool hasReplicaprovisionedreadcapacityautoscalingsettings() => $_has(2);
  @$pb.TagNumber(125131885)
  void clearReplicaprovisionedreadcapacityautoscalingsettings() =>
      $_clearField(125131885);
  @$pb.TagNumber(125131885)
  AutoScalingSettingsDescription
      ensureReplicaprovisionedreadcapacityautoscalingsettings() => $_ensure(2);

  @$pb.TagNumber(293841796)
  AutoScalingSettingsDescription
      get replicaprovisionedwritecapacityautoscalingsettings => $_getN(3);
  @$pb.TagNumber(293841796)
  set replicaprovisionedwritecapacityautoscalingsettings(
          AutoScalingSettingsDescription value) =>
      $_setField(293841796, value);
  @$pb.TagNumber(293841796)
  $core.bool hasReplicaprovisionedwritecapacityautoscalingsettings() =>
      $_has(3);
  @$pb.TagNumber(293841796)
  void clearReplicaprovisionedwritecapacityautoscalingsettings() =>
      $_clearField(293841796);
  @$pb.TagNumber(293841796)
  AutoScalingSettingsDescription
      ensureReplicaprovisionedwritecapacityautoscalingsettings() => $_ensure(3);

  @$pb.TagNumber(356738858)
  $fixnum.Int64 get replicaprovisionedwritecapacityunits => $_getI64(4);
  @$pb.TagNumber(356738858)
  set replicaprovisionedwritecapacityunits($fixnum.Int64 value) =>
      $_setInt64(4, value);
  @$pb.TagNumber(356738858)
  $core.bool hasReplicaprovisionedwritecapacityunits() => $_has(4);
  @$pb.TagNumber(356738858)
  void clearReplicaprovisionedwritecapacityunits() => $_clearField(356738858);

  @$pb.TagNumber(456333796)
  BillingModeSummary get replicabillingmodesummary => $_getN(5);
  @$pb.TagNumber(456333796)
  set replicabillingmodesummary(BillingModeSummary value) =>
      $_setField(456333796, value);
  @$pb.TagNumber(456333796)
  $core.bool hasReplicabillingmodesummary() => $_has(5);
  @$pb.TagNumber(456333796)
  void clearReplicabillingmodesummary() => $_clearField(456333796);
  @$pb.TagNumber(456333796)
  BillingModeSummary ensureReplicabillingmodesummary() => $_ensure(5);

  @$pb.TagNumber(466739730)
  ReplicaStatus get replicastatus => $_getN(6);
  @$pb.TagNumber(466739730)
  set replicastatus(ReplicaStatus value) => $_setField(466739730, value);
  @$pb.TagNumber(466739730)
  $core.bool hasReplicastatus() => $_has(6);
  @$pb.TagNumber(466739730)
  void clearReplicastatus() => $_clearField(466739730);

  @$pb.TagNumber(519483614)
  TableClassSummary get replicatableclasssummary => $_getN(7);
  @$pb.TagNumber(519483614)
  set replicatableclasssummary(TableClassSummary value) =>
      $_setField(519483614, value);
  @$pb.TagNumber(519483614)
  $core.bool hasReplicatableclasssummary() => $_has(7);
  @$pb.TagNumber(519483614)
  void clearReplicatableclasssummary() => $_clearField(519483614);
  @$pb.TagNumber(519483614)
  TableClassSummary ensureReplicatableclasssummary() => $_ensure(7);

  @$pb.TagNumber(521234416)
  $pb.PbList<ReplicaGlobalSecondaryIndexSettingsDescription>
      get replicaglobalsecondaryindexsettings => $_getList(8);
}

class ReplicaSettingsUpdate extends $pb.GeneratedMessage {
  factory ReplicaSettingsUpdate({
    $fixnum.Int64? replicaprovisionedreadcapacityunits,
    $core.String? regionname,
    $core.Iterable<ReplicaGlobalSecondaryIndexSettingsUpdate>?
        replicaglobalsecondaryindexsettingsupdate,
    AutoScalingSettingsUpdate?
        replicaprovisionedreadcapacityautoscalingsettingsupdate,
    TableClass? replicatableclass,
  }) {
    final result = create();
    if (replicaprovisionedreadcapacityunits != null)
      result.replicaprovisionedreadcapacityunits =
          replicaprovisionedreadcapacityunits;
    if (regionname != null) result.regionname = regionname;
    if (replicaglobalsecondaryindexsettingsupdate != null)
      result.replicaglobalsecondaryindexsettingsupdate
          .addAll(replicaglobalsecondaryindexsettingsupdate);
    if (replicaprovisionedreadcapacityautoscalingsettingsupdate != null)
      result.replicaprovisionedreadcapacityautoscalingsettingsupdate =
          replicaprovisionedreadcapacityautoscalingsettingsupdate;
    if (replicatableclass != null) result.replicatableclass = replicatableclass;
    return result;
  }

  ReplicaSettingsUpdate._();

  factory ReplicaSettingsUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaSettingsUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaSettingsUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(
        82081083, _omitFieldNames ? '' : 'replicaprovisionedreadcapacityunits')
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..pPM<ReplicaGlobalSecondaryIndexSettingsUpdate>(116195935,
        _omitFieldNames ? '' : 'replicaglobalsecondaryindexsettingsupdate',
        subBuilder: ReplicaGlobalSecondaryIndexSettingsUpdate.create)
    ..aOM<AutoScalingSettingsUpdate>(
        245262702,
        _omitFieldNames
            ? ''
            : 'replicaprovisionedreadcapacityautoscalingsettingsupdate',
        subBuilder: AutoScalingSettingsUpdate.create)
    ..aE<TableClass>(248679204, _omitFieldNames ? '' : 'replicatableclass',
        enumValues: TableClass.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaSettingsUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaSettingsUpdate copyWith(
          void Function(ReplicaSettingsUpdate) updates) =>
      super.copyWith((message) => updates(message as ReplicaSettingsUpdate))
          as ReplicaSettingsUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaSettingsUpdate create() => ReplicaSettingsUpdate._();
  @$core.override
  ReplicaSettingsUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaSettingsUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaSettingsUpdate>(create);
  static ReplicaSettingsUpdate? _defaultInstance;

  @$pb.TagNumber(82081083)
  $fixnum.Int64 get replicaprovisionedreadcapacityunits => $_getI64(0);
  @$pb.TagNumber(82081083)
  set replicaprovisionedreadcapacityunits($fixnum.Int64 value) =>
      $_setInt64(0, value);
  @$pb.TagNumber(82081083)
  $core.bool hasReplicaprovisionedreadcapacityunits() => $_has(0);
  @$pb.TagNumber(82081083)
  void clearReplicaprovisionedreadcapacityunits() => $_clearField(82081083);

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(1);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(1);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(116195935)
  $pb.PbList<ReplicaGlobalSecondaryIndexSettingsUpdate>
      get replicaglobalsecondaryindexsettingsupdate => $_getList(2);

  @$pb.TagNumber(245262702)
  AutoScalingSettingsUpdate
      get replicaprovisionedreadcapacityautoscalingsettingsupdate => $_getN(3);
  @$pb.TagNumber(245262702)
  set replicaprovisionedreadcapacityautoscalingsettingsupdate(
          AutoScalingSettingsUpdate value) =>
      $_setField(245262702, value);
  @$pb.TagNumber(245262702)
  $core.bool hasReplicaprovisionedreadcapacityautoscalingsettingsupdate() =>
      $_has(3);
  @$pb.TagNumber(245262702)
  void clearReplicaprovisionedreadcapacityautoscalingsettingsupdate() =>
      $_clearField(245262702);
  @$pb.TagNumber(245262702)
  AutoScalingSettingsUpdate
      ensureReplicaprovisionedreadcapacityautoscalingsettingsupdate() =>
          $_ensure(3);

  @$pb.TagNumber(248679204)
  TableClass get replicatableclass => $_getN(4);
  @$pb.TagNumber(248679204)
  set replicatableclass(TableClass value) => $_setField(248679204, value);
  @$pb.TagNumber(248679204)
  $core.bool hasReplicatableclass() => $_has(4);
  @$pb.TagNumber(248679204)
  void clearReplicatableclass() => $_clearField(248679204);
}

class ReplicaUpdate extends $pb.GeneratedMessage {
  factory ReplicaUpdate({
    DeleteReplicaAction? delete,
    CreateReplicaAction? create_420340862,
  }) {
    final result = create();
    if (delete != null) result.delete = delete;
    if (create_420340862 != null) result.create_420340862 = create_420340862;
    return result;
  }

  ReplicaUpdate._();

  factory ReplicaUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicaUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicaUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<DeleteReplicaAction>(395831915, _omitFieldNames ? '' : 'delete',
        subBuilder: DeleteReplicaAction.create)
    ..aOM<CreateReplicaAction>(420340862, _omitFieldNames ? '' : 'create',
        subBuilder: CreateReplicaAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicaUpdate copyWith(void Function(ReplicaUpdate) updates) =>
      super.copyWith((message) => updates(message as ReplicaUpdate))
          as ReplicaUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicaUpdate create() => ReplicaUpdate._();
  @$core.override
  ReplicaUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicaUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicaUpdate>(create);
  static ReplicaUpdate? _defaultInstance;

  @$pb.TagNumber(395831915)
  DeleteReplicaAction get delete => $_getN(0);
  @$pb.TagNumber(395831915)
  set delete(DeleteReplicaAction value) => $_setField(395831915, value);
  @$pb.TagNumber(395831915)
  $core.bool hasDelete() => $_has(0);
  @$pb.TagNumber(395831915)
  void clearDelete() => $_clearField(395831915);
  @$pb.TagNumber(395831915)
  DeleteReplicaAction ensureDelete() => $_ensure(0);

  @$pb.TagNumber(420340862)
  CreateReplicaAction get create_420340862 => $_getN(1);
  @$pb.TagNumber(420340862)
  set create_420340862(CreateReplicaAction value) =>
      $_setField(420340862, value);
  @$pb.TagNumber(420340862)
  $core.bool hasCreate_420340862() => $_has(1);
  @$pb.TagNumber(420340862)
  void clearCreate_420340862() => $_clearField(420340862);
  @$pb.TagNumber(420340862)
  CreateReplicaAction ensureCreate_420340862() => $_ensure(1);
}

class ReplicatedWriteConflictException extends $pb.GeneratedMessage {
  factory ReplicatedWriteConflictException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ReplicatedWriteConflictException._();

  factory ReplicatedWriteConflictException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicatedWriteConflictException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicatedWriteConflictException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicatedWriteConflictException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicatedWriteConflictException copyWith(
          void Function(ReplicatedWriteConflictException) updates) =>
      super.copyWith(
              (message) => updates(message as ReplicatedWriteConflictException))
          as ReplicatedWriteConflictException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicatedWriteConflictException create() =>
      ReplicatedWriteConflictException._();
  @$core.override
  ReplicatedWriteConflictException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicatedWriteConflictException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicatedWriteConflictException>(
          create);
  static ReplicatedWriteConflictException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ReplicationGroupUpdate extends $pb.GeneratedMessage {
  factory ReplicationGroupUpdate({
    UpdateReplicationGroupMemberAction? update,
    DeleteReplicationGroupMemberAction? delete,
    CreateReplicationGroupMemberAction? create_420340862,
  }) {
    final result = create();
    if (update != null) result.update = update;
    if (delete != null) result.delete = delete;
    if (create_420340862 != null) result.create_420340862 = create_420340862;
    return result;
  }

  ReplicationGroupUpdate._();

  factory ReplicationGroupUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplicationGroupUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplicationGroupUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<UpdateReplicationGroupMemberAction>(
        237178517, _omitFieldNames ? '' : 'update',
        subBuilder: UpdateReplicationGroupMemberAction.create)
    ..aOM<DeleteReplicationGroupMemberAction>(
        395831915, _omitFieldNames ? '' : 'delete',
        subBuilder: DeleteReplicationGroupMemberAction.create)
    ..aOM<CreateReplicationGroupMemberAction>(
        420340862, _omitFieldNames ? '' : 'create',
        subBuilder: CreateReplicationGroupMemberAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicationGroupUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplicationGroupUpdate copyWith(
          void Function(ReplicationGroupUpdate) updates) =>
      super.copyWith((message) => updates(message as ReplicationGroupUpdate))
          as ReplicationGroupUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplicationGroupUpdate create() => ReplicationGroupUpdate._();
  @$core.override
  ReplicationGroupUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplicationGroupUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplicationGroupUpdate>(create);
  static ReplicationGroupUpdate? _defaultInstance;

  @$pb.TagNumber(237178517)
  UpdateReplicationGroupMemberAction get update => $_getN(0);
  @$pb.TagNumber(237178517)
  set update(UpdateReplicationGroupMemberAction value) =>
      $_setField(237178517, value);
  @$pb.TagNumber(237178517)
  $core.bool hasUpdate() => $_has(0);
  @$pb.TagNumber(237178517)
  void clearUpdate() => $_clearField(237178517);
  @$pb.TagNumber(237178517)
  UpdateReplicationGroupMemberAction ensureUpdate() => $_ensure(0);

  @$pb.TagNumber(395831915)
  DeleteReplicationGroupMemberAction get delete => $_getN(1);
  @$pb.TagNumber(395831915)
  set delete(DeleteReplicationGroupMemberAction value) =>
      $_setField(395831915, value);
  @$pb.TagNumber(395831915)
  $core.bool hasDelete() => $_has(1);
  @$pb.TagNumber(395831915)
  void clearDelete() => $_clearField(395831915);
  @$pb.TagNumber(395831915)
  DeleteReplicationGroupMemberAction ensureDelete() => $_ensure(1);

  @$pb.TagNumber(420340862)
  CreateReplicationGroupMemberAction get create_420340862 => $_getN(2);
  @$pb.TagNumber(420340862)
  set create_420340862(CreateReplicationGroupMemberAction value) =>
      $_setField(420340862, value);
  @$pb.TagNumber(420340862)
  $core.bool hasCreate_420340862() => $_has(2);
  @$pb.TagNumber(420340862)
  void clearCreate_420340862() => $_clearField(420340862);
  @$pb.TagNumber(420340862)
  CreateReplicationGroupMemberAction ensureCreate_420340862() => $_ensure(2);
}

class RequestLimitExceeded extends $pb.GeneratedMessage {
  factory RequestLimitExceeded({
    $core.String? message,
    $core.Iterable<ThrottlingReason>? throttlingreasons,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (throttlingreasons != null)
      result.throttlingreasons.addAll(throttlingreasons);
    return result;
  }

  RequestLimitExceeded._();

  factory RequestLimitExceeded.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RequestLimitExceeded.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RequestLimitExceeded',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..pPM<ThrottlingReason>(
        284436852, _omitFieldNames ? '' : 'throttlingreasons',
        subBuilder: ThrottlingReason.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestLimitExceeded clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RequestLimitExceeded copyWith(void Function(RequestLimitExceeded) updates) =>
      super.copyWith((message) => updates(message as RequestLimitExceeded))
          as RequestLimitExceeded;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RequestLimitExceeded create() => RequestLimitExceeded._();
  @$core.override
  RequestLimitExceeded createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RequestLimitExceeded getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RequestLimitExceeded>(create);
  static RequestLimitExceeded? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(284436852)
  $pb.PbList<ThrottlingReason> get throttlingreasons => $_getList(1);
}

class ResourceInUseException extends $pb.GeneratedMessage {
  factory ResourceInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourceInUseException._();

  factory ResourceInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceInUseException copyWith(
          void Function(ResourceInUseException) updates) =>
      super.copyWith((message) => updates(message as ResourceInUseException))
          as ResourceInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceInUseException create() => ResourceInUseException._();
  @$core.override
  ResourceInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceInUseException>(create);
  static ResourceInUseException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
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

class RestoreSummary extends $pb.GeneratedMessage {
  factory RestoreSummary({
    $core.String? sourcebackuparn,
    $core.String? sourcetablearn,
    $core.String? restoredatetime,
    $core.bool? restoreinprogress,
  }) {
    final result = create();
    if (sourcebackuparn != null) result.sourcebackuparn = sourcebackuparn;
    if (sourcetablearn != null) result.sourcetablearn = sourcetablearn;
    if (restoredatetime != null) result.restoredatetime = restoredatetime;
    if (restoreinprogress != null) result.restoreinprogress = restoreinprogress;
    return result;
  }

  RestoreSummary._();

  factory RestoreSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(16805842, _omitFieldNames ? '' : 'sourcebackuparn')
    ..aOS(124215796, _omitFieldNames ? '' : 'sourcetablearn')
    ..aOS(176772301, _omitFieldNames ? '' : 'restoredatetime')
    ..aOB(508499598, _omitFieldNames ? '' : 'restoreinprogress')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreSummary copyWith(void Function(RestoreSummary) updates) =>
      super.copyWith((message) => updates(message as RestoreSummary))
          as RestoreSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreSummary create() => RestoreSummary._();
  @$core.override
  RestoreSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreSummary>(create);
  static RestoreSummary? _defaultInstance;

  @$pb.TagNumber(16805842)
  $core.String get sourcebackuparn => $_getSZ(0);
  @$pb.TagNumber(16805842)
  set sourcebackuparn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(16805842)
  $core.bool hasSourcebackuparn() => $_has(0);
  @$pb.TagNumber(16805842)
  void clearSourcebackuparn() => $_clearField(16805842);

  @$pb.TagNumber(124215796)
  $core.String get sourcetablearn => $_getSZ(1);
  @$pb.TagNumber(124215796)
  set sourcetablearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(124215796)
  $core.bool hasSourcetablearn() => $_has(1);
  @$pb.TagNumber(124215796)
  void clearSourcetablearn() => $_clearField(124215796);

  @$pb.TagNumber(176772301)
  $core.String get restoredatetime => $_getSZ(2);
  @$pb.TagNumber(176772301)
  set restoredatetime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(176772301)
  $core.bool hasRestoredatetime() => $_has(2);
  @$pb.TagNumber(176772301)
  void clearRestoredatetime() => $_clearField(176772301);

  @$pb.TagNumber(508499598)
  $core.bool get restoreinprogress => $_getBF(3);
  @$pb.TagNumber(508499598)
  set restoreinprogress($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(508499598)
  $core.bool hasRestoreinprogress() => $_has(3);
  @$pb.TagNumber(508499598)
  void clearRestoreinprogress() => $_clearField(508499598);
}

class RestoreTableFromBackupInput extends $pb.GeneratedMessage {
  factory RestoreTableFromBackupInput({
    $core.String? targettablename,
    OnDemandThroughput? ondemandthroughputoverride,
    BillingMode? billingmodeoverride,
    $core.Iterable<GlobalSecondaryIndex>? globalsecondaryindexoverride,
    $core.String? backuparn,
    ProvisionedThroughput? provisionedthroughputoverride,
    SSESpecification? ssespecificationoverride,
    $core.Iterable<LocalSecondaryIndex>? localsecondaryindexoverride,
  }) {
    final result = create();
    if (targettablename != null) result.targettablename = targettablename;
    if (ondemandthroughputoverride != null)
      result.ondemandthroughputoverride = ondemandthroughputoverride;
    if (billingmodeoverride != null)
      result.billingmodeoverride = billingmodeoverride;
    if (globalsecondaryindexoverride != null)
      result.globalsecondaryindexoverride.addAll(globalsecondaryindexoverride);
    if (backuparn != null) result.backuparn = backuparn;
    if (provisionedthroughputoverride != null)
      result.provisionedthroughputoverride = provisionedthroughputoverride;
    if (ssespecificationoverride != null)
      result.ssespecificationoverride = ssespecificationoverride;
    if (localsecondaryindexoverride != null)
      result.localsecondaryindexoverride.addAll(localsecondaryindexoverride);
    return result;
  }

  RestoreTableFromBackupInput._();

  factory RestoreTableFromBackupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreTableFromBackupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreTableFromBackupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(298767720, _omitFieldNames ? '' : 'targettablename')
    ..aOM<OnDemandThroughput>(
        317165234, _omitFieldNames ? '' : 'ondemandthroughputoverride',
        subBuilder: OnDemandThroughput.create)
    ..aE<BillingMode>(357784560, _omitFieldNames ? '' : 'billingmodeoverride',
        enumValues: BillingMode.values)
    ..pPM<GlobalSecondaryIndex>(
        369844021, _omitFieldNames ? '' : 'globalsecondaryindexoverride',
        subBuilder: GlobalSecondaryIndex.create)
    ..aOS(370874339, _omitFieldNames ? '' : 'backuparn')
    ..aOM<ProvisionedThroughput>(
        413332116, _omitFieldNames ? '' : 'provisionedthroughputoverride',
        subBuilder: ProvisionedThroughput.create)
    ..aOM<SSESpecification>(
        421570884, _omitFieldNames ? '' : 'ssespecificationoverride',
        subBuilder: SSESpecification.create)
    ..pPM<LocalSecondaryIndex>(
        431784607, _omitFieldNames ? '' : 'localsecondaryindexoverride',
        subBuilder: LocalSecondaryIndex.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreTableFromBackupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreTableFromBackupInput copyWith(
          void Function(RestoreTableFromBackupInput) updates) =>
      super.copyWith(
              (message) => updates(message as RestoreTableFromBackupInput))
          as RestoreTableFromBackupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreTableFromBackupInput create() =>
      RestoreTableFromBackupInput._();
  @$core.override
  RestoreTableFromBackupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreTableFromBackupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreTableFromBackupInput>(create);
  static RestoreTableFromBackupInput? _defaultInstance;

  @$pb.TagNumber(298767720)
  $core.String get targettablename => $_getSZ(0);
  @$pb.TagNumber(298767720)
  set targettablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(298767720)
  $core.bool hasTargettablename() => $_has(0);
  @$pb.TagNumber(298767720)
  void clearTargettablename() => $_clearField(298767720);

  @$pb.TagNumber(317165234)
  OnDemandThroughput get ondemandthroughputoverride => $_getN(1);
  @$pb.TagNumber(317165234)
  set ondemandthroughputoverride(OnDemandThroughput value) =>
      $_setField(317165234, value);
  @$pb.TagNumber(317165234)
  $core.bool hasOndemandthroughputoverride() => $_has(1);
  @$pb.TagNumber(317165234)
  void clearOndemandthroughputoverride() => $_clearField(317165234);
  @$pb.TagNumber(317165234)
  OnDemandThroughput ensureOndemandthroughputoverride() => $_ensure(1);

  @$pb.TagNumber(357784560)
  BillingMode get billingmodeoverride => $_getN(2);
  @$pb.TagNumber(357784560)
  set billingmodeoverride(BillingMode value) => $_setField(357784560, value);
  @$pb.TagNumber(357784560)
  $core.bool hasBillingmodeoverride() => $_has(2);
  @$pb.TagNumber(357784560)
  void clearBillingmodeoverride() => $_clearField(357784560);

  @$pb.TagNumber(369844021)
  $pb.PbList<GlobalSecondaryIndex> get globalsecondaryindexoverride =>
      $_getList(3);

  @$pb.TagNumber(370874339)
  $core.String get backuparn => $_getSZ(4);
  @$pb.TagNumber(370874339)
  set backuparn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(370874339)
  $core.bool hasBackuparn() => $_has(4);
  @$pb.TagNumber(370874339)
  void clearBackuparn() => $_clearField(370874339);

  @$pb.TagNumber(413332116)
  ProvisionedThroughput get provisionedthroughputoverride => $_getN(5);
  @$pb.TagNumber(413332116)
  set provisionedthroughputoverride(ProvisionedThroughput value) =>
      $_setField(413332116, value);
  @$pb.TagNumber(413332116)
  $core.bool hasProvisionedthroughputoverride() => $_has(5);
  @$pb.TagNumber(413332116)
  void clearProvisionedthroughputoverride() => $_clearField(413332116);
  @$pb.TagNumber(413332116)
  ProvisionedThroughput ensureProvisionedthroughputoverride() => $_ensure(5);

  @$pb.TagNumber(421570884)
  SSESpecification get ssespecificationoverride => $_getN(6);
  @$pb.TagNumber(421570884)
  set ssespecificationoverride(SSESpecification value) =>
      $_setField(421570884, value);
  @$pb.TagNumber(421570884)
  $core.bool hasSsespecificationoverride() => $_has(6);
  @$pb.TagNumber(421570884)
  void clearSsespecificationoverride() => $_clearField(421570884);
  @$pb.TagNumber(421570884)
  SSESpecification ensureSsespecificationoverride() => $_ensure(6);

  @$pb.TagNumber(431784607)
  $pb.PbList<LocalSecondaryIndex> get localsecondaryindexoverride =>
      $_getList(7);
}

class RestoreTableFromBackupOutput extends $pb.GeneratedMessage {
  factory RestoreTableFromBackupOutput({
    TableDescription? tabledescription,
  }) {
    final result = create();
    if (tabledescription != null) result.tabledescription = tabledescription;
    return result;
  }

  RestoreTableFromBackupOutput._();

  factory RestoreTableFromBackupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreTableFromBackupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreTableFromBackupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TableDescription>(19962388, _omitFieldNames ? '' : 'tabledescription',
        subBuilder: TableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreTableFromBackupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreTableFromBackupOutput copyWith(
          void Function(RestoreTableFromBackupOutput) updates) =>
      super.copyWith(
              (message) => updates(message as RestoreTableFromBackupOutput))
          as RestoreTableFromBackupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreTableFromBackupOutput create() =>
      RestoreTableFromBackupOutput._();
  @$core.override
  RestoreTableFromBackupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreTableFromBackupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreTableFromBackupOutput>(create);
  static RestoreTableFromBackupOutput? _defaultInstance;

  @$pb.TagNumber(19962388)
  TableDescription get tabledescription => $_getN(0);
  @$pb.TagNumber(19962388)
  set tabledescription(TableDescription value) => $_setField(19962388, value);
  @$pb.TagNumber(19962388)
  $core.bool hasTabledescription() => $_has(0);
  @$pb.TagNumber(19962388)
  void clearTabledescription() => $_clearField(19962388);
  @$pb.TagNumber(19962388)
  TableDescription ensureTabledescription() => $_ensure(0);
}

class RestoreTableToPointInTimeInput extends $pb.GeneratedMessage {
  factory RestoreTableToPointInTimeInput({
    $core.String? sourcetablearn,
    $core.String? sourcetablename,
    $core.String? restoredatetime,
    $core.String? targettablename,
    OnDemandThroughput? ondemandthroughputoverride,
    BillingMode? billingmodeoverride,
    $core.Iterable<GlobalSecondaryIndex>? globalsecondaryindexoverride,
    ProvisionedThroughput? provisionedthroughputoverride,
    SSESpecification? ssespecificationoverride,
    $core.Iterable<LocalSecondaryIndex>? localsecondaryindexoverride,
    $core.bool? uselatestrestorabletime,
  }) {
    final result = create();
    if (sourcetablearn != null) result.sourcetablearn = sourcetablearn;
    if (sourcetablename != null) result.sourcetablename = sourcetablename;
    if (restoredatetime != null) result.restoredatetime = restoredatetime;
    if (targettablename != null) result.targettablename = targettablename;
    if (ondemandthroughputoverride != null)
      result.ondemandthroughputoverride = ondemandthroughputoverride;
    if (billingmodeoverride != null)
      result.billingmodeoverride = billingmodeoverride;
    if (globalsecondaryindexoverride != null)
      result.globalsecondaryindexoverride.addAll(globalsecondaryindexoverride);
    if (provisionedthroughputoverride != null)
      result.provisionedthroughputoverride = provisionedthroughputoverride;
    if (ssespecificationoverride != null)
      result.ssespecificationoverride = ssespecificationoverride;
    if (localsecondaryindexoverride != null)
      result.localsecondaryindexoverride.addAll(localsecondaryindexoverride);
    if (uselatestrestorabletime != null)
      result.uselatestrestorabletime = uselatestrestorabletime;
    return result;
  }

  RestoreTableToPointInTimeInput._();

  factory RestoreTableToPointInTimeInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreTableToPointInTimeInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreTableToPointInTimeInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(124215796, _omitFieldNames ? '' : 'sourcetablearn')
    ..aOS(165368952, _omitFieldNames ? '' : 'sourcetablename')
    ..aOS(176772301, _omitFieldNames ? '' : 'restoredatetime')
    ..aOS(298767720, _omitFieldNames ? '' : 'targettablename')
    ..aOM<OnDemandThroughput>(
        317165234, _omitFieldNames ? '' : 'ondemandthroughputoverride',
        subBuilder: OnDemandThroughput.create)
    ..aE<BillingMode>(357784560, _omitFieldNames ? '' : 'billingmodeoverride',
        enumValues: BillingMode.values)
    ..pPM<GlobalSecondaryIndex>(
        369844021, _omitFieldNames ? '' : 'globalsecondaryindexoverride',
        subBuilder: GlobalSecondaryIndex.create)
    ..aOM<ProvisionedThroughput>(
        413332116, _omitFieldNames ? '' : 'provisionedthroughputoverride',
        subBuilder: ProvisionedThroughput.create)
    ..aOM<SSESpecification>(
        421570884, _omitFieldNames ? '' : 'ssespecificationoverride',
        subBuilder: SSESpecification.create)
    ..pPM<LocalSecondaryIndex>(
        431784607, _omitFieldNames ? '' : 'localsecondaryindexoverride',
        subBuilder: LocalSecondaryIndex.create)
    ..aOB(434512618, _omitFieldNames ? '' : 'uselatestrestorabletime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreTableToPointInTimeInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreTableToPointInTimeInput copyWith(
          void Function(RestoreTableToPointInTimeInput) updates) =>
      super.copyWith(
              (message) => updates(message as RestoreTableToPointInTimeInput))
          as RestoreTableToPointInTimeInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreTableToPointInTimeInput create() =>
      RestoreTableToPointInTimeInput._();
  @$core.override
  RestoreTableToPointInTimeInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreTableToPointInTimeInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreTableToPointInTimeInput>(create);
  static RestoreTableToPointInTimeInput? _defaultInstance;

  @$pb.TagNumber(124215796)
  $core.String get sourcetablearn => $_getSZ(0);
  @$pb.TagNumber(124215796)
  set sourcetablearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(124215796)
  $core.bool hasSourcetablearn() => $_has(0);
  @$pb.TagNumber(124215796)
  void clearSourcetablearn() => $_clearField(124215796);

  @$pb.TagNumber(165368952)
  $core.String get sourcetablename => $_getSZ(1);
  @$pb.TagNumber(165368952)
  set sourcetablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(165368952)
  $core.bool hasSourcetablename() => $_has(1);
  @$pb.TagNumber(165368952)
  void clearSourcetablename() => $_clearField(165368952);

  @$pb.TagNumber(176772301)
  $core.String get restoredatetime => $_getSZ(2);
  @$pb.TagNumber(176772301)
  set restoredatetime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(176772301)
  $core.bool hasRestoredatetime() => $_has(2);
  @$pb.TagNumber(176772301)
  void clearRestoredatetime() => $_clearField(176772301);

  @$pb.TagNumber(298767720)
  $core.String get targettablename => $_getSZ(3);
  @$pb.TagNumber(298767720)
  set targettablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(298767720)
  $core.bool hasTargettablename() => $_has(3);
  @$pb.TagNumber(298767720)
  void clearTargettablename() => $_clearField(298767720);

  @$pb.TagNumber(317165234)
  OnDemandThroughput get ondemandthroughputoverride => $_getN(4);
  @$pb.TagNumber(317165234)
  set ondemandthroughputoverride(OnDemandThroughput value) =>
      $_setField(317165234, value);
  @$pb.TagNumber(317165234)
  $core.bool hasOndemandthroughputoverride() => $_has(4);
  @$pb.TagNumber(317165234)
  void clearOndemandthroughputoverride() => $_clearField(317165234);
  @$pb.TagNumber(317165234)
  OnDemandThroughput ensureOndemandthroughputoverride() => $_ensure(4);

  @$pb.TagNumber(357784560)
  BillingMode get billingmodeoverride => $_getN(5);
  @$pb.TagNumber(357784560)
  set billingmodeoverride(BillingMode value) => $_setField(357784560, value);
  @$pb.TagNumber(357784560)
  $core.bool hasBillingmodeoverride() => $_has(5);
  @$pb.TagNumber(357784560)
  void clearBillingmodeoverride() => $_clearField(357784560);

  @$pb.TagNumber(369844021)
  $pb.PbList<GlobalSecondaryIndex> get globalsecondaryindexoverride =>
      $_getList(6);

  @$pb.TagNumber(413332116)
  ProvisionedThroughput get provisionedthroughputoverride => $_getN(7);
  @$pb.TagNumber(413332116)
  set provisionedthroughputoverride(ProvisionedThroughput value) =>
      $_setField(413332116, value);
  @$pb.TagNumber(413332116)
  $core.bool hasProvisionedthroughputoverride() => $_has(7);
  @$pb.TagNumber(413332116)
  void clearProvisionedthroughputoverride() => $_clearField(413332116);
  @$pb.TagNumber(413332116)
  ProvisionedThroughput ensureProvisionedthroughputoverride() => $_ensure(7);

  @$pb.TagNumber(421570884)
  SSESpecification get ssespecificationoverride => $_getN(8);
  @$pb.TagNumber(421570884)
  set ssespecificationoverride(SSESpecification value) =>
      $_setField(421570884, value);
  @$pb.TagNumber(421570884)
  $core.bool hasSsespecificationoverride() => $_has(8);
  @$pb.TagNumber(421570884)
  void clearSsespecificationoverride() => $_clearField(421570884);
  @$pb.TagNumber(421570884)
  SSESpecification ensureSsespecificationoverride() => $_ensure(8);

  @$pb.TagNumber(431784607)
  $pb.PbList<LocalSecondaryIndex> get localsecondaryindexoverride =>
      $_getList(9);

  @$pb.TagNumber(434512618)
  $core.bool get uselatestrestorabletime => $_getBF(10);
  @$pb.TagNumber(434512618)
  set uselatestrestorabletime($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(434512618)
  $core.bool hasUselatestrestorabletime() => $_has(10);
  @$pb.TagNumber(434512618)
  void clearUselatestrestorabletime() => $_clearField(434512618);
}

class RestoreTableToPointInTimeOutput extends $pb.GeneratedMessage {
  factory RestoreTableToPointInTimeOutput({
    TableDescription? tabledescription,
  }) {
    final result = create();
    if (tabledescription != null) result.tabledescription = tabledescription;
    return result;
  }

  RestoreTableToPointInTimeOutput._();

  factory RestoreTableToPointInTimeOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreTableToPointInTimeOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreTableToPointInTimeOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TableDescription>(19962388, _omitFieldNames ? '' : 'tabledescription',
        subBuilder: TableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreTableToPointInTimeOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreTableToPointInTimeOutput copyWith(
          void Function(RestoreTableToPointInTimeOutput) updates) =>
      super.copyWith(
              (message) => updates(message as RestoreTableToPointInTimeOutput))
          as RestoreTableToPointInTimeOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreTableToPointInTimeOutput create() =>
      RestoreTableToPointInTimeOutput._();
  @$core.override
  RestoreTableToPointInTimeOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreTableToPointInTimeOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreTableToPointInTimeOutput>(
          create);
  static RestoreTableToPointInTimeOutput? _defaultInstance;

  @$pb.TagNumber(19962388)
  TableDescription get tabledescription => $_getN(0);
  @$pb.TagNumber(19962388)
  set tabledescription(TableDescription value) => $_setField(19962388, value);
  @$pb.TagNumber(19962388)
  $core.bool hasTabledescription() => $_has(0);
  @$pb.TagNumber(19962388)
  void clearTabledescription() => $_clearField(19962388);
  @$pb.TagNumber(19962388)
  TableDescription ensureTabledescription() => $_ensure(0);
}

class S3BucketSource extends $pb.GeneratedMessage {
  factory S3BucketSource({
    $core.String? s3bucket,
    $core.String? s3keyprefix,
    $core.String? s3bucketowner,
  }) {
    final result = create();
    if (s3bucket != null) result.s3bucket = s3bucket;
    if (s3keyprefix != null) result.s3keyprefix = s3keyprefix;
    if (s3bucketowner != null) result.s3bucketowner = s3bucketowner;
    return result;
  }

  S3BucketSource._();

  factory S3BucketSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory S3BucketSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'S3BucketSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(114031434, _omitFieldNames ? '' : 's3bucket')
    ..aOS(206015359, _omitFieldNames ? '' : 's3keyprefix')
    ..aOS(351576129, _omitFieldNames ? '' : 's3bucketowner')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3BucketSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  S3BucketSource copyWith(void Function(S3BucketSource) updates) =>
      super.copyWith((message) => updates(message as S3BucketSource))
          as S3BucketSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static S3BucketSource create() => S3BucketSource._();
  @$core.override
  S3BucketSource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static S3BucketSource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<S3BucketSource>(create);
  static S3BucketSource? _defaultInstance;

  @$pb.TagNumber(114031434)
  $core.String get s3bucket => $_getSZ(0);
  @$pb.TagNumber(114031434)
  set s3bucket($core.String value) => $_setString(0, value);
  @$pb.TagNumber(114031434)
  $core.bool hasS3bucket() => $_has(0);
  @$pb.TagNumber(114031434)
  void clearS3bucket() => $_clearField(114031434);

  @$pb.TagNumber(206015359)
  $core.String get s3keyprefix => $_getSZ(1);
  @$pb.TagNumber(206015359)
  set s3keyprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(206015359)
  $core.bool hasS3keyprefix() => $_has(1);
  @$pb.TagNumber(206015359)
  void clearS3keyprefix() => $_clearField(206015359);

  @$pb.TagNumber(351576129)
  $core.String get s3bucketowner => $_getSZ(2);
  @$pb.TagNumber(351576129)
  set s3bucketowner($core.String value) => $_setString(2, value);
  @$pb.TagNumber(351576129)
  $core.bool hasS3bucketowner() => $_has(2);
  @$pb.TagNumber(351576129)
  void clearS3bucketowner() => $_clearField(351576129);
}

class SSEDescription extends $pb.GeneratedMessage {
  factory SSEDescription({
    SSEStatus? status,
    $core.String? kmsmasterkeyarn,
    SSEType? ssetype,
    $core.String? inaccessibleencryptiondatetime,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (kmsmasterkeyarn != null) result.kmsmasterkeyarn = kmsmasterkeyarn;
    if (ssetype != null) result.ssetype = ssetype;
    if (inaccessibleencryptiondatetime != null)
      result.inaccessibleencryptiondatetime = inaccessibleencryptiondatetime;
    return result;
  }

  SSEDescription._();

  factory SSEDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SSEDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SSEDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<SSEStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: SSEStatus.values)
    ..aOS(281937987, _omitFieldNames ? '' : 'kmsmasterkeyarn')
    ..aE<SSEType>(431846435, _omitFieldNames ? '' : 'ssetype',
        enumValues: SSEType.values)
    ..aOS(478188797, _omitFieldNames ? '' : 'inaccessibleencryptiondatetime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SSEDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SSEDescription copyWith(void Function(SSEDescription) updates) =>
      super.copyWith((message) => updates(message as SSEDescription))
          as SSEDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SSEDescription create() => SSEDescription._();
  @$core.override
  SSEDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SSEDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SSEDescription>(create);
  static SSEDescription? _defaultInstance;

  @$pb.TagNumber(6222352)
  SSEStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(SSEStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(281937987)
  $core.String get kmsmasterkeyarn => $_getSZ(1);
  @$pb.TagNumber(281937987)
  set kmsmasterkeyarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(281937987)
  $core.bool hasKmsmasterkeyarn() => $_has(1);
  @$pb.TagNumber(281937987)
  void clearKmsmasterkeyarn() => $_clearField(281937987);

  @$pb.TagNumber(431846435)
  SSEType get ssetype => $_getN(2);
  @$pb.TagNumber(431846435)
  set ssetype(SSEType value) => $_setField(431846435, value);
  @$pb.TagNumber(431846435)
  $core.bool hasSsetype() => $_has(2);
  @$pb.TagNumber(431846435)
  void clearSsetype() => $_clearField(431846435);

  @$pb.TagNumber(478188797)
  $core.String get inaccessibleencryptiondatetime => $_getSZ(3);
  @$pb.TagNumber(478188797)
  set inaccessibleencryptiondatetime($core.String value) =>
      $_setString(3, value);
  @$pb.TagNumber(478188797)
  $core.bool hasInaccessibleencryptiondatetime() => $_has(3);
  @$pb.TagNumber(478188797)
  void clearInaccessibleencryptiondatetime() => $_clearField(478188797);
}

class SSESpecification extends $pb.GeneratedMessage {
  factory SSESpecification({
    SSEType? ssetype,
    $core.bool? enabled,
    $core.String? kmsmasterkeyid,
  }) {
    final result = create();
    if (ssetype != null) result.ssetype = ssetype;
    if (enabled != null) result.enabled = enabled;
    if (kmsmasterkeyid != null) result.kmsmasterkeyid = kmsmasterkeyid;
    return result;
  }

  SSESpecification._();

  factory SSESpecification.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SSESpecification.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SSESpecification',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<SSEType>(431846435, _omitFieldNames ? '' : 'ssetype',
        enumValues: SSEType.values)
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..aOS(521459443, _omitFieldNames ? '' : 'kmsmasterkeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SSESpecification clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SSESpecification copyWith(void Function(SSESpecification) updates) =>
      super.copyWith((message) => updates(message as SSESpecification))
          as SSESpecification;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SSESpecification create() => SSESpecification._();
  @$core.override
  SSESpecification createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SSESpecification getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SSESpecification>(create);
  static SSESpecification? _defaultInstance;

  @$pb.TagNumber(431846435)
  SSEType get ssetype => $_getN(0);
  @$pb.TagNumber(431846435)
  set ssetype(SSEType value) => $_setField(431846435, value);
  @$pb.TagNumber(431846435)
  $core.bool hasSsetype() => $_has(0);
  @$pb.TagNumber(431846435)
  void clearSsetype() => $_clearField(431846435);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);

  @$pb.TagNumber(521459443)
  $core.String get kmsmasterkeyid => $_getSZ(2);
  @$pb.TagNumber(521459443)
  set kmsmasterkeyid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(521459443)
  $core.bool hasKmsmasterkeyid() => $_has(2);
  @$pb.TagNumber(521459443)
  void clearKmsmasterkeyid() => $_clearField(521459443);
}

class ScanInput extends $pb.GeneratedMessage {
  factory ScanInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.String? filterexpression,
    $core.String? indexname,
    $core.int? totalsegments,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.String? projectionexpression,
    ConditionalOperator? conditionaloperator,
    $core.String? tablename,
    $core.Iterable<$core.MapEntry<$core.String, Condition>>? scanfilter,
    $core.int? segment,
    $core.Iterable<$core.String>? attributestoget,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        exclusivestartkey,
    $core.int? limit,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
    Select? select,
    $core.bool? consistentread,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (filterexpression != null) result.filterexpression = filterexpression;
    if (indexname != null) result.indexname = indexname;
    if (totalsegments != null) result.totalsegments = totalsegments;
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (projectionexpression != null)
      result.projectionexpression = projectionexpression;
    if (conditionaloperator != null)
      result.conditionaloperator = conditionaloperator;
    if (tablename != null) result.tablename = tablename;
    if (scanfilter != null) result.scanfilter.addEntries(scanfilter);
    if (segment != null) result.segment = segment;
    if (attributestoget != null) result.attributestoget.addAll(attributestoget);
    if (exclusivestartkey != null)
      result.exclusivestartkey.addEntries(exclusivestartkey);
    if (limit != null) result.limit = limit;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    if (select != null) result.select = select;
    if (consistentread != null) result.consistentread = consistentread;
    return result;
  }

  ScanInput._();

  factory ScanInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScanInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScanInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..aOS(68019610, _omitFieldNames ? '' : 'filterexpression')
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aI(149136904, _omitFieldNames ? '' : 'totalsegments')
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'ScanInput.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(150730243, _omitFieldNames ? '' : 'projectionexpression')
    ..aE<ConditionalOperator>(
        172066260, _omitFieldNames ? '' : 'conditionaloperator',
        enumValues: ConditionalOperator.values)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..m<$core.String, Condition>(272885755, _omitFieldNames ? '' : 'scanfilter',
        entryClassName: 'ScanInput.ScanfilterEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: Condition.create,
        valueDefaultOrMaker: Condition.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aI(279654279, _omitFieldNames ? '' : 'segment')
    ..pPS(311382592, _omitFieldNames ? '' : 'attributestoget')
    ..m<$core.String, AttributeValue>(
        348137297, _omitFieldNames ? '' : 'exclusivestartkey',
        entryClassName: 'ScanInput.ExclusivestartkeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'ScanInput.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<Select>(512305998, _omitFieldNames ? '' : 'select',
        enumValues: Select.values)
    ..aOB(531556994, _omitFieldNames ? '' : 'consistentread')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScanInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScanInput copyWith(void Function(ScanInput) updates) =>
      super.copyWith((message) => updates(message as ScanInput)) as ScanInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScanInput create() => ScanInput._();
  @$core.override
  ScanInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScanInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ScanInput>(create);
  static ScanInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(68019610)
  $core.String get filterexpression => $_getSZ(1);
  @$pb.TagNumber(68019610)
  set filterexpression($core.String value) => $_setString(1, value);
  @$pb.TagNumber(68019610)
  $core.bool hasFilterexpression() => $_has(1);
  @$pb.TagNumber(68019610)
  void clearFilterexpression() => $_clearField(68019610);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(2);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(2);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(149136904)
  $core.int get totalsegments => $_getIZ(3);
  @$pb.TagNumber(149136904)
  set totalsegments($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(149136904)
  $core.bool hasTotalsegments() => $_has(3);
  @$pb.TagNumber(149136904)
  void clearTotalsegments() => $_clearField(149136904);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(4);

  @$pb.TagNumber(150730243)
  $core.String get projectionexpression => $_getSZ(5);
  @$pb.TagNumber(150730243)
  set projectionexpression($core.String value) => $_setString(5, value);
  @$pb.TagNumber(150730243)
  $core.bool hasProjectionexpression() => $_has(5);
  @$pb.TagNumber(150730243)
  void clearProjectionexpression() => $_clearField(150730243);

  @$pb.TagNumber(172066260)
  ConditionalOperator get conditionaloperator => $_getN(6);
  @$pb.TagNumber(172066260)
  set conditionaloperator(ConditionalOperator value) =>
      $_setField(172066260, value);
  @$pb.TagNumber(172066260)
  $core.bool hasConditionaloperator() => $_has(6);
  @$pb.TagNumber(172066260)
  void clearConditionaloperator() => $_clearField(172066260);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(7);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(7, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(7);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(272885755)
  $pb.PbMap<$core.String, Condition> get scanfilter => $_getMap(8);

  @$pb.TagNumber(279654279)
  $core.int get segment => $_getIZ(9);
  @$pb.TagNumber(279654279)
  set segment($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(279654279)
  $core.bool hasSegment() => $_has(9);
  @$pb.TagNumber(279654279)
  void clearSegment() => $_clearField(279654279);

  @$pb.TagNumber(311382592)
  $pb.PbList<$core.String> get attributestoget => $_getList(10);

  @$pb.TagNumber(348137297)
  $pb.PbMap<$core.String, AttributeValue> get exclusivestartkey => $_getMap(11);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(12);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(12, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(12);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(13);

  @$pb.TagNumber(512305998)
  Select get select => $_getN(14);
  @$pb.TagNumber(512305998)
  set select(Select value) => $_setField(512305998, value);
  @$pb.TagNumber(512305998)
  $core.bool hasSelect() => $_has(14);
  @$pb.TagNumber(512305998)
  void clearSelect() => $_clearField(512305998);

  @$pb.TagNumber(531556994)
  $core.bool get consistentread => $_getBF(15);
  @$pb.TagNumber(531556994)
  set consistentread($core.bool value) => $_setBool(15, value);
  @$pb.TagNumber(531556994)
  $core.bool hasConsistentread() => $_has(15);
  @$pb.TagNumber(531556994)
  void clearConsistentread() => $_clearField(531556994);
}

class ScanOutput extends $pb.GeneratedMessage {
  factory ScanOutput({
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? items,
    $core.int? count,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        lastevaluatedkey,
    ConsumedCapacity? consumedcapacity,
    $core.int? scannedcount,
  }) {
    final result = create();
    if (items != null) result.items.addEntries(items);
    if (count != null) result.count = count;
    if (lastevaluatedkey != null)
      result.lastevaluatedkey.addEntries(lastevaluatedkey);
    if (consumedcapacity != null) result.consumedcapacity = consumedcapacity;
    if (scannedcount != null) result.scannedcount = scannedcount;
    return result;
  }

  ScanOutput._();

  factory ScanOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScanOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScanOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, AttributeValue>(3553328, _omitFieldNames ? '' : 'items',
        entryClassName: 'ScanOutput.ItemsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aI(31963285, _omitFieldNames ? '' : 'count')
    ..m<$core.String, AttributeValue>(
        54319830, _omitFieldNames ? '' : 'lastevaluatedkey',
        entryClassName: 'ScanOutput.LastevaluatedkeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..aI(531161315, _omitFieldNames ? '' : 'scannedcount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScanOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScanOutput copyWith(void Function(ScanOutput) updates) =>
      super.copyWith((message) => updates(message as ScanOutput)) as ScanOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScanOutput create() => ScanOutput._();
  @$core.override
  ScanOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScanOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScanOutput>(create);
  static ScanOutput? _defaultInstance;

  @$pb.TagNumber(3553328)
  $pb.PbMap<$core.String, AttributeValue> get items => $_getMap(0);

  @$pb.TagNumber(31963285)
  $core.int get count => $_getIZ(1);
  @$pb.TagNumber(31963285)
  set count($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(31963285)
  $core.bool hasCount() => $_has(1);
  @$pb.TagNumber(31963285)
  void clearCount() => $_clearField(31963285);

  @$pb.TagNumber(54319830)
  $pb.PbMap<$core.String, AttributeValue> get lastevaluatedkey => $_getMap(2);

  @$pb.TagNumber(449336620)
  ConsumedCapacity get consumedcapacity => $_getN(3);
  @$pb.TagNumber(449336620)
  set consumedcapacity(ConsumedCapacity value) => $_setField(449336620, value);
  @$pb.TagNumber(449336620)
  $core.bool hasConsumedcapacity() => $_has(3);
  @$pb.TagNumber(449336620)
  void clearConsumedcapacity() => $_clearField(449336620);
  @$pb.TagNumber(449336620)
  ConsumedCapacity ensureConsumedcapacity() => $_ensure(3);

  @$pb.TagNumber(531161315)
  $core.int get scannedcount => $_getIZ(4);
  @$pb.TagNumber(531161315)
  set scannedcount($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(531161315)
  $core.bool hasScannedcount() => $_has(4);
  @$pb.TagNumber(531161315)
  void clearScannedcount() => $_clearField(531161315);
}

class SourceTableDetails extends $pb.GeneratedMessage {
  factory SourceTableDetails({
    ProvisionedThroughput? provisionedthroughput,
    $fixnum.Int64? itemcount,
    $core.String? tablecreationdatetime,
    BillingMode? billingmode,
    $fixnum.Int64? tablesizebytes,
    $core.String? tablename,
    $core.Iterable<KeySchemaElement>? keyschema,
    $core.String? tablearn,
    $core.String? tableid,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (itemcount != null) result.itemcount = itemcount;
    if (tablecreationdatetime != null)
      result.tablecreationdatetime = tablecreationdatetime;
    if (billingmode != null) result.billingmode = billingmode;
    if (tablesizebytes != null) result.tablesizebytes = tablesizebytes;
    if (tablename != null) result.tablename = tablename;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (tablearn != null) result.tablearn = tablearn;
    if (tableid != null) result.tableid = tableid;
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  SourceTableDetails._();

  factory SourceTableDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SourceTableDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SourceTableDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughput>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughput.create)
    ..aInt64(26280022, _omitFieldNames ? '' : 'itemcount')
    ..aOS(152706380, _omitFieldNames ? '' : 'tablecreationdatetime')
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aInt64(220631294, _omitFieldNames ? '' : 'tablesizebytes')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..aOS(449893011, _omitFieldNames ? '' : 'tableid')
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SourceTableDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SourceTableDetails copyWith(void Function(SourceTableDetails) updates) =>
      super.copyWith((message) => updates(message as SourceTableDetails))
          as SourceTableDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SourceTableDetails create() => SourceTableDetails._();
  @$core.override
  SourceTableDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SourceTableDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SourceTableDetails>(create);
  static SourceTableDetails? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughput get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughput value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughput ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(26280022)
  $fixnum.Int64 get itemcount => $_getI64(1);
  @$pb.TagNumber(26280022)
  set itemcount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(26280022)
  $core.bool hasItemcount() => $_has(1);
  @$pb.TagNumber(26280022)
  void clearItemcount() => $_clearField(26280022);

  @$pb.TagNumber(152706380)
  $core.String get tablecreationdatetime => $_getSZ(2);
  @$pb.TagNumber(152706380)
  set tablecreationdatetime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(152706380)
  $core.bool hasTablecreationdatetime() => $_has(2);
  @$pb.TagNumber(152706380)
  void clearTablecreationdatetime() => $_clearField(152706380);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(3);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(3);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(220631294)
  $fixnum.Int64 get tablesizebytes => $_getI64(4);
  @$pb.TagNumber(220631294)
  set tablesizebytes($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(220631294)
  $core.bool hasTablesizebytes() => $_has(4);
  @$pb.TagNumber(220631294)
  void clearTablesizebytes() => $_clearField(220631294);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(5);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(5, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(5);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(6);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(7);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(7, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(7);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);

  @$pb.TagNumber(449893011)
  $core.String get tableid => $_getSZ(8);
  @$pb.TagNumber(449893011)
  set tableid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(449893011)
  $core.bool hasTableid() => $_has(8);
  @$pb.TagNumber(449893011)
  void clearTableid() => $_clearField(449893011);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(9);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(9);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(9);
}

class SourceTableFeatureDetails extends $pb.GeneratedMessage {
  factory SourceTableFeatureDetails({
    SSEDescription? ssedescription,
    $core.Iterable<LocalSecondaryIndexInfo>? localsecondaryindexes,
    StreamSpecification? streamdescription,
    TimeToLiveDescription? timetolivedescription,
    $core.Iterable<GlobalSecondaryIndexInfo>? globalsecondaryindexes,
  }) {
    final result = create();
    if (ssedescription != null) result.ssedescription = ssedescription;
    if (localsecondaryindexes != null)
      result.localsecondaryindexes.addAll(localsecondaryindexes);
    if (streamdescription != null) result.streamdescription = streamdescription;
    if (timetolivedescription != null)
      result.timetolivedescription = timetolivedescription;
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addAll(globalsecondaryindexes);
    return result;
  }

  SourceTableFeatureDetails._();

  factory SourceTableFeatureDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SourceTableFeatureDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SourceTableFeatureDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<SSEDescription>(350068773, _omitFieldNames ? '' : 'ssedescription',
        subBuilder: SSEDescription.create)
    ..pPM<LocalSecondaryIndexInfo>(
        362339959, _omitFieldNames ? '' : 'localsecondaryindexes',
        subBuilder: LocalSecondaryIndexInfo.create)
    ..aOM<StreamSpecification>(
        363737034, _omitFieldNames ? '' : 'streamdescription',
        subBuilder: StreamSpecification.create)
    ..aOM<TimeToLiveDescription>(
        367590876, _omitFieldNames ? '' : 'timetolivedescription',
        subBuilder: TimeToLiveDescription.create)
    ..pPM<GlobalSecondaryIndexInfo>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        subBuilder: GlobalSecondaryIndexInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SourceTableFeatureDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SourceTableFeatureDetails copyWith(
          void Function(SourceTableFeatureDetails) updates) =>
      super.copyWith((message) => updates(message as SourceTableFeatureDetails))
          as SourceTableFeatureDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SourceTableFeatureDetails create() => SourceTableFeatureDetails._();
  @$core.override
  SourceTableFeatureDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SourceTableFeatureDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SourceTableFeatureDetails>(create);
  static SourceTableFeatureDetails? _defaultInstance;

  @$pb.TagNumber(350068773)
  SSEDescription get ssedescription => $_getN(0);
  @$pb.TagNumber(350068773)
  set ssedescription(SSEDescription value) => $_setField(350068773, value);
  @$pb.TagNumber(350068773)
  $core.bool hasSsedescription() => $_has(0);
  @$pb.TagNumber(350068773)
  void clearSsedescription() => $_clearField(350068773);
  @$pb.TagNumber(350068773)
  SSEDescription ensureSsedescription() => $_ensure(0);

  @$pb.TagNumber(362339959)
  $pb.PbList<LocalSecondaryIndexInfo> get localsecondaryindexes => $_getList(1);

  @$pb.TagNumber(363737034)
  StreamSpecification get streamdescription => $_getN(2);
  @$pb.TagNumber(363737034)
  set streamdescription(StreamSpecification value) =>
      $_setField(363737034, value);
  @$pb.TagNumber(363737034)
  $core.bool hasStreamdescription() => $_has(2);
  @$pb.TagNumber(363737034)
  void clearStreamdescription() => $_clearField(363737034);
  @$pb.TagNumber(363737034)
  StreamSpecification ensureStreamdescription() => $_ensure(2);

  @$pb.TagNumber(367590876)
  TimeToLiveDescription get timetolivedescription => $_getN(3);
  @$pb.TagNumber(367590876)
  set timetolivedescription(TimeToLiveDescription value) =>
      $_setField(367590876, value);
  @$pb.TagNumber(367590876)
  $core.bool hasTimetolivedescription() => $_has(3);
  @$pb.TagNumber(367590876)
  void clearTimetolivedescription() => $_clearField(367590876);
  @$pb.TagNumber(367590876)
  TimeToLiveDescription ensureTimetolivedescription() => $_ensure(3);

  @$pb.TagNumber(409156905)
  $pb.PbList<GlobalSecondaryIndexInfo> get globalsecondaryindexes =>
      $_getList(4);
}

class StreamSpecification extends $pb.GeneratedMessage {
  factory StreamSpecification({
    $core.bool? streamenabled,
    StreamViewType? streamviewtype,
  }) {
    final result = create();
    if (streamenabled != null) result.streamenabled = streamenabled;
    if (streamviewtype != null) result.streamviewtype = streamviewtype;
    return result;
  }

  StreamSpecification._();

  factory StreamSpecification.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StreamSpecification.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StreamSpecification',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOB(266707711, _omitFieldNames ? '' : 'streamenabled')
    ..aE<StreamViewType>(380488241, _omitFieldNames ? '' : 'streamviewtype',
        enumValues: StreamViewType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamSpecification clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StreamSpecification copyWith(void Function(StreamSpecification) updates) =>
      super.copyWith((message) => updates(message as StreamSpecification))
          as StreamSpecification;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StreamSpecification create() => StreamSpecification._();
  @$core.override
  StreamSpecification createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StreamSpecification getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StreamSpecification>(create);
  static StreamSpecification? _defaultInstance;

  @$pb.TagNumber(266707711)
  $core.bool get streamenabled => $_getBF(0);
  @$pb.TagNumber(266707711)
  set streamenabled($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(266707711)
  $core.bool hasStreamenabled() => $_has(0);
  @$pb.TagNumber(266707711)
  void clearStreamenabled() => $_clearField(266707711);

  @$pb.TagNumber(380488241)
  StreamViewType get streamviewtype => $_getN(1);
  @$pb.TagNumber(380488241)
  set streamviewtype(StreamViewType value) => $_setField(380488241, value);
  @$pb.TagNumber(380488241)
  $core.bool hasStreamviewtype() => $_has(1);
  @$pb.TagNumber(380488241)
  void clearStreamviewtype() => $_clearField(380488241);
}

class TableAlreadyExistsException extends $pb.GeneratedMessage {
  factory TableAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TableAlreadyExistsException._();

  factory TableAlreadyExistsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableAlreadyExistsException copyWith(
          void Function(TableAlreadyExistsException) updates) =>
      super.copyWith(
              (message) => updates(message as TableAlreadyExistsException))
          as TableAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableAlreadyExistsException create() =>
      TableAlreadyExistsException._();
  @$core.override
  TableAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableAlreadyExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableAlreadyExistsException>(create);
  static TableAlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TableAutoScalingDescription extends $pb.GeneratedMessage {
  factory TableAutoScalingDescription({
    TableStatus? tablestatus,
    $core.String? tablename,
    $core.Iterable<ReplicaAutoScalingDescription>? replicas,
  }) {
    final result = create();
    if (tablestatus != null) result.tablestatus = tablestatus;
    if (tablename != null) result.tablename = tablename;
    if (replicas != null) result.replicas.addAll(replicas);
    return result;
  }

  TableAutoScalingDescription._();

  factory TableAutoScalingDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableAutoScalingDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableAutoScalingDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<TableStatus>(207908810, _omitFieldNames ? '' : 'tablestatus',
        enumValues: TableStatus.values)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..pPM<ReplicaAutoScalingDescription>(
        306066781, _omitFieldNames ? '' : 'replicas',
        subBuilder: ReplicaAutoScalingDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableAutoScalingDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableAutoScalingDescription copyWith(
          void Function(TableAutoScalingDescription) updates) =>
      super.copyWith(
              (message) => updates(message as TableAutoScalingDescription))
          as TableAutoScalingDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableAutoScalingDescription create() =>
      TableAutoScalingDescription._();
  @$core.override
  TableAutoScalingDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableAutoScalingDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableAutoScalingDescription>(create);
  static TableAutoScalingDescription? _defaultInstance;

  @$pb.TagNumber(207908810)
  TableStatus get tablestatus => $_getN(0);
  @$pb.TagNumber(207908810)
  set tablestatus(TableStatus value) => $_setField(207908810, value);
  @$pb.TagNumber(207908810)
  $core.bool hasTablestatus() => $_has(0);
  @$pb.TagNumber(207908810)
  void clearTablestatus() => $_clearField(207908810);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(306066781)
  $pb.PbList<ReplicaAutoScalingDescription> get replicas => $_getList(2);
}

class TableClassSummary extends $pb.GeneratedMessage {
  factory TableClassSummary({
    TableClass? tableclass,
    $core.String? lastupdatedatetime,
  }) {
    final result = create();
    if (tableclass != null) result.tableclass = tableclass;
    if (lastupdatedatetime != null)
      result.lastupdatedatetime = lastupdatedatetime;
    return result;
  }

  TableClassSummary._();

  factory TableClassSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableClassSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableClassSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<TableClass>(342890498, _omitFieldNames ? '' : 'tableclass',
        enumValues: TableClass.values)
    ..aOS(452274318, _omitFieldNames ? '' : 'lastupdatedatetime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableClassSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableClassSummary copyWith(void Function(TableClassSummary) updates) =>
      super.copyWith((message) => updates(message as TableClassSummary))
          as TableClassSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableClassSummary create() => TableClassSummary._();
  @$core.override
  TableClassSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableClassSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableClassSummary>(create);
  static TableClassSummary? _defaultInstance;

  @$pb.TagNumber(342890498)
  TableClass get tableclass => $_getN(0);
  @$pb.TagNumber(342890498)
  set tableclass(TableClass value) => $_setField(342890498, value);
  @$pb.TagNumber(342890498)
  $core.bool hasTableclass() => $_has(0);
  @$pb.TagNumber(342890498)
  void clearTableclass() => $_clearField(342890498);

  @$pb.TagNumber(452274318)
  $core.String get lastupdatedatetime => $_getSZ(1);
  @$pb.TagNumber(452274318)
  set lastupdatedatetime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(452274318)
  $core.bool hasLastupdatedatetime() => $_has(1);
  @$pb.TagNumber(452274318)
  void clearLastupdatedatetime() => $_clearField(452274318);
}

class TableCreationParameters extends $pb.GeneratedMessage {
  factory TableCreationParameters({
    ProvisionedThroughput? provisionedthroughput,
    SSESpecification? ssespecification,
    BillingMode? billingmode,
    $core.String? tablename,
    $core.Iterable<KeySchemaElement>? keyschema,
    $core.Iterable<GlobalSecondaryIndex>? globalsecondaryindexes,
    $core.Iterable<AttributeDefinition>? attributedefinitions,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (ssespecification != null) result.ssespecification = ssespecification;
    if (billingmode != null) result.billingmode = billingmode;
    if (tablename != null) result.tablename = tablename;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addAll(globalsecondaryindexes);
    if (attributedefinitions != null)
      result.attributedefinitions.addAll(attributedefinitions);
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  TableCreationParameters._();

  factory TableCreationParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableCreationParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableCreationParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughput>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughput.create)
    ..aOM<SSESpecification>(31692444, _omitFieldNames ? '' : 'ssespecification',
        subBuilder: SSESpecification.create)
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..pPM<GlobalSecondaryIndex>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        subBuilder: GlobalSecondaryIndex.create)
    ..pPM<AttributeDefinition>(
        414687108, _omitFieldNames ? '' : 'attributedefinitions',
        subBuilder: AttributeDefinition.create)
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableCreationParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableCreationParameters copyWith(
          void Function(TableCreationParameters) updates) =>
      super.copyWith((message) => updates(message as TableCreationParameters))
          as TableCreationParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableCreationParameters create() => TableCreationParameters._();
  @$core.override
  TableCreationParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableCreationParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableCreationParameters>(create);
  static TableCreationParameters? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughput get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughput value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughput ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(31692444)
  SSESpecification get ssespecification => $_getN(1);
  @$pb.TagNumber(31692444)
  set ssespecification(SSESpecification value) => $_setField(31692444, value);
  @$pb.TagNumber(31692444)
  $core.bool hasSsespecification() => $_has(1);
  @$pb.TagNumber(31692444)
  void clearSsespecification() => $_clearField(31692444);
  @$pb.TagNumber(31692444)
  SSESpecification ensureSsespecification() => $_ensure(1);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(2);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(2);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(4);

  @$pb.TagNumber(409156905)
  $pb.PbList<GlobalSecondaryIndex> get globalsecondaryindexes => $_getList(5);

  @$pb.TagNumber(414687108)
  $pb.PbList<AttributeDefinition> get attributedefinitions => $_getList(6);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(7);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(7);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(7);
}

class TableDescription extends $pb.GeneratedMessage {
  factory TableDescription({
    ProvisionedThroughputDescription? provisionedthroughput,
    TableClassSummary? tableclasssummary,
    $core.Iterable<GlobalTableWitnessDescription>? globaltablewitnesses,
    GlobalTableSettingsReplicationMode? globaltablesettingsreplicationmode,
    $fixnum.Int64? itemcount,
    $core.String? creationdatetime,
    ArchivalSummary? archivalsummary,
    $core.String? globaltableversion,
    BillingModeSummary? billingmodesummary,
    $core.String? lateststreamarn,
    TableStatus? tablestatus,
    $fixnum.Int64? tablesizebytes,
    $core.bool? deletionprotectionenabled,
    $core.String? tablename,
    TableWarmThroughputDescription? warmthroughput,
    $core.Iterable<KeySchemaElement>? keyschema,
    $core.Iterable<ReplicaDescription>? replicas,
    $core.String? lateststreamlabel,
    RestoreSummary? restoresummary,
    SSEDescription? ssedescription,
    $core.Iterable<LocalSecondaryIndexDescription>? localsecondaryindexes,
    StreamSpecification? streamspecification,
    $core.Iterable<GlobalSecondaryIndexDescription>? globalsecondaryindexes,
    $core.Iterable<AttributeDefinition>? attributedefinitions,
    $core.String? tablearn,
    MultiRegionConsistency? multiregionconsistency,
    $core.String? tableid,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (tableclasssummary != null) result.tableclasssummary = tableclasssummary;
    if (globaltablewitnesses != null)
      result.globaltablewitnesses.addAll(globaltablewitnesses);
    if (globaltablesettingsreplicationmode != null)
      result.globaltablesettingsreplicationmode =
          globaltablesettingsreplicationmode;
    if (itemcount != null) result.itemcount = itemcount;
    if (creationdatetime != null) result.creationdatetime = creationdatetime;
    if (archivalsummary != null) result.archivalsummary = archivalsummary;
    if (globaltableversion != null)
      result.globaltableversion = globaltableversion;
    if (billingmodesummary != null)
      result.billingmodesummary = billingmodesummary;
    if (lateststreamarn != null) result.lateststreamarn = lateststreamarn;
    if (tablestatus != null) result.tablestatus = tablestatus;
    if (tablesizebytes != null) result.tablesizebytes = tablesizebytes;
    if (deletionprotectionenabled != null)
      result.deletionprotectionenabled = deletionprotectionenabled;
    if (tablename != null) result.tablename = tablename;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (keyschema != null) result.keyschema.addAll(keyschema);
    if (replicas != null) result.replicas.addAll(replicas);
    if (lateststreamlabel != null) result.lateststreamlabel = lateststreamlabel;
    if (restoresummary != null) result.restoresummary = restoresummary;
    if (ssedescription != null) result.ssedescription = ssedescription;
    if (localsecondaryindexes != null)
      result.localsecondaryindexes.addAll(localsecondaryindexes);
    if (streamspecification != null)
      result.streamspecification = streamspecification;
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addAll(globalsecondaryindexes);
    if (attributedefinitions != null)
      result.attributedefinitions.addAll(attributedefinitions);
    if (tablearn != null) result.tablearn = tablearn;
    if (multiregionconsistency != null)
      result.multiregionconsistency = multiregionconsistency;
    if (tableid != null) result.tableid = tableid;
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  TableDescription._();

  factory TableDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughputDescription>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughputDescription.create)
    ..aOM<TableClassSummary>(
        4371552, _omitFieldNames ? '' : 'tableclasssummary',
        subBuilder: TableClassSummary.create)
    ..pPM<GlobalTableWitnessDescription>(
        4521286, _omitFieldNames ? '' : 'globaltablewitnesses',
        subBuilder: GlobalTableWitnessDescription.create)
    ..aE<GlobalTableSettingsReplicationMode>(
        10446577, _omitFieldNames ? '' : 'globaltablesettingsreplicationmode',
        enumValues: GlobalTableSettingsReplicationMode.values)
    ..aInt64(26280022, _omitFieldNames ? '' : 'itemcount')
    ..aOS(48904698, _omitFieldNames ? '' : 'creationdatetime')
    ..aOM<ArchivalSummary>(52039658, _omitFieldNames ? '' : 'archivalsummary',
        subBuilder: ArchivalSummary.create)
    ..aOS(68234287, _omitFieldNames ? '' : 'globaltableversion')
    ..aOM<BillingModeSummary>(
        163529882, _omitFieldNames ? '' : 'billingmodesummary',
        subBuilder: BillingModeSummary.create)
    ..aOS(207365682, _omitFieldNames ? '' : 'lateststreamarn')
    ..aE<TableStatus>(207908810, _omitFieldNames ? '' : 'tablestatus',
        enumValues: TableStatus.values)
    ..aInt64(220631294, _omitFieldNames ? '' : 'tablesizebytes')
    ..aOB(259418450, _omitFieldNames ? '' : 'deletionprotectionenabled')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<TableWarmThroughputDescription>(
        290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: TableWarmThroughputDescription.create)
    ..pPM<KeySchemaElement>(293038056, _omitFieldNames ? '' : 'keyschema',
        subBuilder: KeySchemaElement.create)
    ..pPM<ReplicaDescription>(306066781, _omitFieldNames ? '' : 'replicas',
        subBuilder: ReplicaDescription.create)
    ..aOS(328475293, _omitFieldNames ? '' : 'lateststreamlabel')
    ..aOM<RestoreSummary>(330529648, _omitFieldNames ? '' : 'restoresummary',
        subBuilder: RestoreSummary.create)
    ..aOM<SSEDescription>(350068773, _omitFieldNames ? '' : 'ssedescription',
        subBuilder: SSEDescription.create)
    ..pPM<LocalSecondaryIndexDescription>(
        362339959, _omitFieldNames ? '' : 'localsecondaryindexes',
        subBuilder: LocalSecondaryIndexDescription.create)
    ..aOM<StreamSpecification>(
        403922627, _omitFieldNames ? '' : 'streamspecification',
        subBuilder: StreamSpecification.create)
    ..pPM<GlobalSecondaryIndexDescription>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        subBuilder: GlobalSecondaryIndexDescription.create)
    ..pPM<AttributeDefinition>(
        414687108, _omitFieldNames ? '' : 'attributedefinitions',
        subBuilder: AttributeDefinition.create)
    ..aOS(431669347, _omitFieldNames ? '' : 'tablearn')
    ..aE<MultiRegionConsistency>(
        446019131, _omitFieldNames ? '' : 'multiregionconsistency',
        enumValues: MultiRegionConsistency.values)
    ..aOS(449893011, _omitFieldNames ? '' : 'tableid')
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableDescription copyWith(void Function(TableDescription) updates) =>
      super.copyWith((message) => updates(message as TableDescription))
          as TableDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableDescription create() => TableDescription._();
  @$core.override
  TableDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableDescription>(create);
  static TableDescription? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughputDescription get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughputDescription value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughputDescription ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(4371552)
  TableClassSummary get tableclasssummary => $_getN(1);
  @$pb.TagNumber(4371552)
  set tableclasssummary(TableClassSummary value) => $_setField(4371552, value);
  @$pb.TagNumber(4371552)
  $core.bool hasTableclasssummary() => $_has(1);
  @$pb.TagNumber(4371552)
  void clearTableclasssummary() => $_clearField(4371552);
  @$pb.TagNumber(4371552)
  TableClassSummary ensureTableclasssummary() => $_ensure(1);

  @$pb.TagNumber(4521286)
  $pb.PbList<GlobalTableWitnessDescription> get globaltablewitnesses =>
      $_getList(2);

  @$pb.TagNumber(10446577)
  GlobalTableSettingsReplicationMode get globaltablesettingsreplicationmode =>
      $_getN(3);
  @$pb.TagNumber(10446577)
  set globaltablesettingsreplicationmode(
          GlobalTableSettingsReplicationMode value) =>
      $_setField(10446577, value);
  @$pb.TagNumber(10446577)
  $core.bool hasGlobaltablesettingsreplicationmode() => $_has(3);
  @$pb.TagNumber(10446577)
  void clearGlobaltablesettingsreplicationmode() => $_clearField(10446577);

  @$pb.TagNumber(26280022)
  $fixnum.Int64 get itemcount => $_getI64(4);
  @$pb.TagNumber(26280022)
  set itemcount($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(26280022)
  $core.bool hasItemcount() => $_has(4);
  @$pb.TagNumber(26280022)
  void clearItemcount() => $_clearField(26280022);

  @$pb.TagNumber(48904698)
  $core.String get creationdatetime => $_getSZ(5);
  @$pb.TagNumber(48904698)
  set creationdatetime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(48904698)
  $core.bool hasCreationdatetime() => $_has(5);
  @$pb.TagNumber(48904698)
  void clearCreationdatetime() => $_clearField(48904698);

  @$pb.TagNumber(52039658)
  ArchivalSummary get archivalsummary => $_getN(6);
  @$pb.TagNumber(52039658)
  set archivalsummary(ArchivalSummary value) => $_setField(52039658, value);
  @$pb.TagNumber(52039658)
  $core.bool hasArchivalsummary() => $_has(6);
  @$pb.TagNumber(52039658)
  void clearArchivalsummary() => $_clearField(52039658);
  @$pb.TagNumber(52039658)
  ArchivalSummary ensureArchivalsummary() => $_ensure(6);

  @$pb.TagNumber(68234287)
  $core.String get globaltableversion => $_getSZ(7);
  @$pb.TagNumber(68234287)
  set globaltableversion($core.String value) => $_setString(7, value);
  @$pb.TagNumber(68234287)
  $core.bool hasGlobaltableversion() => $_has(7);
  @$pb.TagNumber(68234287)
  void clearGlobaltableversion() => $_clearField(68234287);

  @$pb.TagNumber(163529882)
  BillingModeSummary get billingmodesummary => $_getN(8);
  @$pb.TagNumber(163529882)
  set billingmodesummary(BillingModeSummary value) =>
      $_setField(163529882, value);
  @$pb.TagNumber(163529882)
  $core.bool hasBillingmodesummary() => $_has(8);
  @$pb.TagNumber(163529882)
  void clearBillingmodesummary() => $_clearField(163529882);
  @$pb.TagNumber(163529882)
  BillingModeSummary ensureBillingmodesummary() => $_ensure(8);

  @$pb.TagNumber(207365682)
  $core.String get lateststreamarn => $_getSZ(9);
  @$pb.TagNumber(207365682)
  set lateststreamarn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(207365682)
  $core.bool hasLateststreamarn() => $_has(9);
  @$pb.TagNumber(207365682)
  void clearLateststreamarn() => $_clearField(207365682);

  @$pb.TagNumber(207908810)
  TableStatus get tablestatus => $_getN(10);
  @$pb.TagNumber(207908810)
  set tablestatus(TableStatus value) => $_setField(207908810, value);
  @$pb.TagNumber(207908810)
  $core.bool hasTablestatus() => $_has(10);
  @$pb.TagNumber(207908810)
  void clearTablestatus() => $_clearField(207908810);

  @$pb.TagNumber(220631294)
  $fixnum.Int64 get tablesizebytes => $_getI64(11);
  @$pb.TagNumber(220631294)
  set tablesizebytes($fixnum.Int64 value) => $_setInt64(11, value);
  @$pb.TagNumber(220631294)
  $core.bool hasTablesizebytes() => $_has(11);
  @$pb.TagNumber(220631294)
  void clearTablesizebytes() => $_clearField(220631294);

  @$pb.TagNumber(259418450)
  $core.bool get deletionprotectionenabled => $_getBF(12);
  @$pb.TagNumber(259418450)
  set deletionprotectionenabled($core.bool value) => $_setBool(12, value);
  @$pb.TagNumber(259418450)
  $core.bool hasDeletionprotectionenabled() => $_has(12);
  @$pb.TagNumber(259418450)
  void clearDeletionprotectionenabled() => $_clearField(259418450);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(13);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(13, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(13);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(290598659)
  TableWarmThroughputDescription get warmthroughput => $_getN(14);
  @$pb.TagNumber(290598659)
  set warmthroughput(TableWarmThroughputDescription value) =>
      $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(14);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  TableWarmThroughputDescription ensureWarmthroughput() => $_ensure(14);

  @$pb.TagNumber(293038056)
  $pb.PbList<KeySchemaElement> get keyschema => $_getList(15);

  @$pb.TagNumber(306066781)
  $pb.PbList<ReplicaDescription> get replicas => $_getList(16);

  @$pb.TagNumber(328475293)
  $core.String get lateststreamlabel => $_getSZ(17);
  @$pb.TagNumber(328475293)
  set lateststreamlabel($core.String value) => $_setString(17, value);
  @$pb.TagNumber(328475293)
  $core.bool hasLateststreamlabel() => $_has(17);
  @$pb.TagNumber(328475293)
  void clearLateststreamlabel() => $_clearField(328475293);

  @$pb.TagNumber(330529648)
  RestoreSummary get restoresummary => $_getN(18);
  @$pb.TagNumber(330529648)
  set restoresummary(RestoreSummary value) => $_setField(330529648, value);
  @$pb.TagNumber(330529648)
  $core.bool hasRestoresummary() => $_has(18);
  @$pb.TagNumber(330529648)
  void clearRestoresummary() => $_clearField(330529648);
  @$pb.TagNumber(330529648)
  RestoreSummary ensureRestoresummary() => $_ensure(18);

  @$pb.TagNumber(350068773)
  SSEDescription get ssedescription => $_getN(19);
  @$pb.TagNumber(350068773)
  set ssedescription(SSEDescription value) => $_setField(350068773, value);
  @$pb.TagNumber(350068773)
  $core.bool hasSsedescription() => $_has(19);
  @$pb.TagNumber(350068773)
  void clearSsedescription() => $_clearField(350068773);
  @$pb.TagNumber(350068773)
  SSEDescription ensureSsedescription() => $_ensure(19);

  @$pb.TagNumber(362339959)
  $pb.PbList<LocalSecondaryIndexDescription> get localsecondaryindexes =>
      $_getList(20);

  @$pb.TagNumber(403922627)
  StreamSpecification get streamspecification => $_getN(21);
  @$pb.TagNumber(403922627)
  set streamspecification(StreamSpecification value) =>
      $_setField(403922627, value);
  @$pb.TagNumber(403922627)
  $core.bool hasStreamspecification() => $_has(21);
  @$pb.TagNumber(403922627)
  void clearStreamspecification() => $_clearField(403922627);
  @$pb.TagNumber(403922627)
  StreamSpecification ensureStreamspecification() => $_ensure(21);

  @$pb.TagNumber(409156905)
  $pb.PbList<GlobalSecondaryIndexDescription> get globalsecondaryindexes =>
      $_getList(22);

  @$pb.TagNumber(414687108)
  $pb.PbList<AttributeDefinition> get attributedefinitions => $_getList(23);

  @$pb.TagNumber(431669347)
  $core.String get tablearn => $_getSZ(24);
  @$pb.TagNumber(431669347)
  set tablearn($core.String value) => $_setString(24, value);
  @$pb.TagNumber(431669347)
  $core.bool hasTablearn() => $_has(24);
  @$pb.TagNumber(431669347)
  void clearTablearn() => $_clearField(431669347);

  @$pb.TagNumber(446019131)
  MultiRegionConsistency get multiregionconsistency => $_getN(25);
  @$pb.TagNumber(446019131)
  set multiregionconsistency(MultiRegionConsistency value) =>
      $_setField(446019131, value);
  @$pb.TagNumber(446019131)
  $core.bool hasMultiregionconsistency() => $_has(25);
  @$pb.TagNumber(446019131)
  void clearMultiregionconsistency() => $_clearField(446019131);

  @$pb.TagNumber(449893011)
  $core.String get tableid => $_getSZ(26);
  @$pb.TagNumber(449893011)
  set tableid($core.String value) => $_setString(26, value);
  @$pb.TagNumber(449893011)
  $core.bool hasTableid() => $_has(26);
  @$pb.TagNumber(449893011)
  void clearTableid() => $_clearField(449893011);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(27);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(27);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(27);
}

class TableInUseException extends $pb.GeneratedMessage {
  factory TableInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TableInUseException._();

  factory TableInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableInUseException copyWith(void Function(TableInUseException) updates) =>
      super.copyWith((message) => updates(message as TableInUseException))
          as TableInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableInUseException create() => TableInUseException._();
  @$core.override
  TableInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableInUseException>(create);
  static TableInUseException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TableNotFoundException extends $pb.GeneratedMessage {
  factory TableNotFoundException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TableNotFoundException._();

  factory TableNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableNotFoundException copyWith(
          void Function(TableNotFoundException) updates) =>
      super.copyWith((message) => updates(message as TableNotFoundException))
          as TableNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableNotFoundException create() => TableNotFoundException._();
  @$core.override
  TableNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableNotFoundException>(create);
  static TableNotFoundException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TableWarmThroughputDescription extends $pb.GeneratedMessage {
  factory TableWarmThroughputDescription({
    TableStatus? status,
    $fixnum.Int64? readunitspersecond,
    $fixnum.Int64? writeunitspersecond,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (readunitspersecond != null)
      result.readunitspersecond = readunitspersecond;
    if (writeunitspersecond != null)
      result.writeunitspersecond = writeunitspersecond;
    return result;
  }

  TableWarmThroughputDescription._();

  factory TableWarmThroughputDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TableWarmThroughputDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TableWarmThroughputDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<TableStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: TableStatus.values)
    ..aInt64(11400732, _omitFieldNames ? '' : 'readunitspersecond')
    ..aInt64(339770127, _omitFieldNames ? '' : 'writeunitspersecond')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableWarmThroughputDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TableWarmThroughputDescription copyWith(
          void Function(TableWarmThroughputDescription) updates) =>
      super.copyWith(
              (message) => updates(message as TableWarmThroughputDescription))
          as TableWarmThroughputDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TableWarmThroughputDescription create() =>
      TableWarmThroughputDescription._();
  @$core.override
  TableWarmThroughputDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TableWarmThroughputDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TableWarmThroughputDescription>(create);
  static TableWarmThroughputDescription? _defaultInstance;

  @$pb.TagNumber(6222352)
  TableStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(TableStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(11400732)
  $fixnum.Int64 get readunitspersecond => $_getI64(1);
  @$pb.TagNumber(11400732)
  set readunitspersecond($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(11400732)
  $core.bool hasReadunitspersecond() => $_has(1);
  @$pb.TagNumber(11400732)
  void clearReadunitspersecond() => $_clearField(11400732);

  @$pb.TagNumber(339770127)
  $fixnum.Int64 get writeunitspersecond => $_getI64(2);
  @$pb.TagNumber(339770127)
  set writeunitspersecond($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(339770127)
  $core.bool hasWriteunitspersecond() => $_has(2);
  @$pb.TagNumber(339770127)
  void clearWriteunitspersecond() => $_clearField(339770127);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
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
  $pb.PbList<Tag> get tags => $_getList(1);
}

class ThrottlingException extends $pb.GeneratedMessage {
  factory ThrottlingException({
    $core.String? message,
    $core.Iterable<ThrottlingReason>? throttlingreasons,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (throttlingreasons != null)
      result.throttlingreasons.addAll(throttlingreasons);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..pPM<ThrottlingReason>(
        170616084, _omitFieldNames ? '' : 'throttlingreasons',
        subBuilder: ThrottlingReason.create)
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(170616084)
  $pb.PbList<ThrottlingReason> get throttlingreasons => $_getList(1);
}

class ThrottlingReason extends $pb.GeneratedMessage {
  factory ThrottlingReason({
    $core.String? resource,
    $core.String? reason,
  }) {
    final result = create();
    if (resource != null) result.resource = resource;
    if (reason != null) result.reason = reason;
    return result;
  }

  ThrottlingReason._();

  factory ThrottlingReason.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ThrottlingReason.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ThrottlingReason',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(165642230, _omitFieldNames ? '' : 'resource')
    ..aOS(413359642, _omitFieldNames ? '' : 'reason')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ThrottlingReason clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ThrottlingReason copyWith(void Function(ThrottlingReason) updates) =>
      super.copyWith((message) => updates(message as ThrottlingReason))
          as ThrottlingReason;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ThrottlingReason create() => ThrottlingReason._();
  @$core.override
  ThrottlingReason createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ThrottlingReason getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ThrottlingReason>(create);
  static ThrottlingReason? _defaultInstance;

  @$pb.TagNumber(165642230)
  $core.String get resource => $_getSZ(0);
  @$pb.TagNumber(165642230)
  set resource($core.String value) => $_setString(0, value);
  @$pb.TagNumber(165642230)
  $core.bool hasResource() => $_has(0);
  @$pb.TagNumber(165642230)
  void clearResource() => $_clearField(165642230);

  @$pb.TagNumber(413359642)
  $core.String get reason => $_getSZ(1);
  @$pb.TagNumber(413359642)
  set reason($core.String value) => $_setString(1, value);
  @$pb.TagNumber(413359642)
  $core.bool hasReason() => $_has(1);
  @$pb.TagNumber(413359642)
  void clearReason() => $_clearField(413359642);
}

class TimeToLiveDescription extends $pb.GeneratedMessage {
  factory TimeToLiveDescription({
    TimeToLiveStatus? timetolivestatus,
    $core.String? attributename,
  }) {
    final result = create();
    if (timetolivestatus != null) result.timetolivestatus = timetolivestatus;
    if (attributename != null) result.attributename = attributename;
    return result;
  }

  TimeToLiveDescription._();

  factory TimeToLiveDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TimeToLiveDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TimeToLiveDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<TimeToLiveStatus>(279467554, _omitFieldNames ? '' : 'timetolivestatus',
        enumValues: TimeToLiveStatus.values)
    ..aOS(352717485, _omitFieldNames ? '' : 'attributename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeToLiveDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeToLiveDescription copyWith(
          void Function(TimeToLiveDescription) updates) =>
      super.copyWith((message) => updates(message as TimeToLiveDescription))
          as TimeToLiveDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TimeToLiveDescription create() => TimeToLiveDescription._();
  @$core.override
  TimeToLiveDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TimeToLiveDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TimeToLiveDescription>(create);
  static TimeToLiveDescription? _defaultInstance;

  @$pb.TagNumber(279467554)
  TimeToLiveStatus get timetolivestatus => $_getN(0);
  @$pb.TagNumber(279467554)
  set timetolivestatus(TimeToLiveStatus value) => $_setField(279467554, value);
  @$pb.TagNumber(279467554)
  $core.bool hasTimetolivestatus() => $_has(0);
  @$pb.TagNumber(279467554)
  void clearTimetolivestatus() => $_clearField(279467554);

  @$pb.TagNumber(352717485)
  $core.String get attributename => $_getSZ(1);
  @$pb.TagNumber(352717485)
  set attributename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(352717485)
  $core.bool hasAttributename() => $_has(1);
  @$pb.TagNumber(352717485)
  void clearAttributename() => $_clearField(352717485);
}

class TimeToLiveSpecification extends $pb.GeneratedMessage {
  factory TimeToLiveSpecification({
    $core.String? attributename,
    $core.bool? enabled,
  }) {
    final result = create();
    if (attributename != null) result.attributename = attributename;
    if (enabled != null) result.enabled = enabled;
    return result;
  }

  TimeToLiveSpecification._();

  factory TimeToLiveSpecification.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TimeToLiveSpecification.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TimeToLiveSpecification',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(352717485, _omitFieldNames ? '' : 'attributename')
    ..aOB(478602303, _omitFieldNames ? '' : 'enabled')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeToLiveSpecification clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TimeToLiveSpecification copyWith(
          void Function(TimeToLiveSpecification) updates) =>
      super.copyWith((message) => updates(message as TimeToLiveSpecification))
          as TimeToLiveSpecification;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TimeToLiveSpecification create() => TimeToLiveSpecification._();
  @$core.override
  TimeToLiveSpecification createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TimeToLiveSpecification getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TimeToLiveSpecification>(create);
  static TimeToLiveSpecification? _defaultInstance;

  @$pb.TagNumber(352717485)
  $core.String get attributename => $_getSZ(0);
  @$pb.TagNumber(352717485)
  set attributename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(352717485)
  $core.bool hasAttributename() => $_has(0);
  @$pb.TagNumber(352717485)
  void clearAttributename() => $_clearField(352717485);

  @$pb.TagNumber(478602303)
  $core.bool get enabled => $_getBF(1);
  @$pb.TagNumber(478602303)
  set enabled($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(478602303)
  $core.bool hasEnabled() => $_has(1);
  @$pb.TagNumber(478602303)
  void clearEnabled() => $_clearField(478602303);
}

class TransactGetItem extends $pb.GeneratedMessage {
  factory TransactGetItem({
    Get? get,
  }) {
    final result = create();
    if (get != null) result.get = get;
    return result;
  }

  TransactGetItem._();

  factory TransactGetItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactGetItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactGetItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<Get>(379010808, _omitFieldNames ? '' : 'get', subBuilder: Get.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactGetItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactGetItem copyWith(void Function(TransactGetItem) updates) =>
      super.copyWith((message) => updates(message as TransactGetItem))
          as TransactGetItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactGetItem create() => TransactGetItem._();
  @$core.override
  TransactGetItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactGetItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactGetItem>(create);
  static TransactGetItem? _defaultInstance;

  @$pb.TagNumber(379010808)
  Get get get => $_getN(0);
  @$pb.TagNumber(379010808)
  set get(Get value) => $_setField(379010808, value);
  @$pb.TagNumber(379010808)
  $core.bool hasGet() => $_has(0);
  @$pb.TagNumber(379010808)
  void clearGet() => $_clearField(379010808);
  @$pb.TagNumber(379010808)
  Get ensureGet() => $_ensure(0);
}

class TransactGetItemsInput extends $pb.GeneratedMessage {
  factory TransactGetItemsInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<TransactGetItem>? transactitems,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (transactitems != null) result.transactitems.addAll(transactitems);
    return result;
  }

  TransactGetItemsInput._();

  factory TransactGetItemsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactGetItemsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactGetItemsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..pPM<TransactGetItem>(506245290, _omitFieldNames ? '' : 'transactitems',
        subBuilder: TransactGetItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactGetItemsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactGetItemsInput copyWith(
          void Function(TransactGetItemsInput) updates) =>
      super.copyWith((message) => updates(message as TransactGetItemsInput))
          as TransactGetItemsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactGetItemsInput create() => TransactGetItemsInput._();
  @$core.override
  TransactGetItemsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactGetItemsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactGetItemsInput>(create);
  static TransactGetItemsInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(506245290)
  $pb.PbList<TransactGetItem> get transactitems => $_getList(1);
}

class TransactGetItemsOutput extends $pb.GeneratedMessage {
  factory TransactGetItemsOutput({
    $core.Iterable<ItemResponse>? responses,
    $core.Iterable<ConsumedCapacity>? consumedcapacity,
  }) {
    final result = create();
    if (responses != null) result.responses.addAll(responses);
    if (consumedcapacity != null)
      result.consumedcapacity.addAll(consumedcapacity);
    return result;
  }

  TransactGetItemsOutput._();

  factory TransactGetItemsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactGetItemsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactGetItemsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<ItemResponse>(114852856, _omitFieldNames ? '' : 'responses',
        subBuilder: ItemResponse.create)
    ..pPM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactGetItemsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactGetItemsOutput copyWith(
          void Function(TransactGetItemsOutput) updates) =>
      super.copyWith((message) => updates(message as TransactGetItemsOutput))
          as TransactGetItemsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactGetItemsOutput create() => TransactGetItemsOutput._();
  @$core.override
  TransactGetItemsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactGetItemsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactGetItemsOutput>(create);
  static TransactGetItemsOutput? _defaultInstance;

  @$pb.TagNumber(114852856)
  $pb.PbList<ItemResponse> get responses => $_getList(0);

  @$pb.TagNumber(449336620)
  $pb.PbList<ConsumedCapacity> get consumedcapacity => $_getList(1);
}

class TransactWriteItem extends $pb.GeneratedMessage {
  factory TransactWriteItem({
    ConditionCheck? conditioncheck,
    Update? update,
    Put? put,
    Delete? delete,
  }) {
    final result = create();
    if (conditioncheck != null) result.conditioncheck = conditioncheck;
    if (update != null) result.update = update;
    if (put != null) result.put = put;
    if (delete != null) result.delete = delete;
    return result;
  }

  TransactWriteItem._();

  factory TransactWriteItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactWriteItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactWriteItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ConditionCheck>(167629711, _omitFieldNames ? '' : 'conditioncheck',
        subBuilder: ConditionCheck.create)
    ..aOM<Update>(237178517, _omitFieldNames ? '' : 'update',
        subBuilder: Update.create)
    ..aOM<Put>(242822543, _omitFieldNames ? '' : 'put', subBuilder: Put.create)
    ..aOM<Delete>(395831915, _omitFieldNames ? '' : 'delete',
        subBuilder: Delete.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactWriteItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactWriteItem copyWith(void Function(TransactWriteItem) updates) =>
      super.copyWith((message) => updates(message as TransactWriteItem))
          as TransactWriteItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactWriteItem create() => TransactWriteItem._();
  @$core.override
  TransactWriteItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactWriteItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactWriteItem>(create);
  static TransactWriteItem? _defaultInstance;

  @$pb.TagNumber(167629711)
  ConditionCheck get conditioncheck => $_getN(0);
  @$pb.TagNumber(167629711)
  set conditioncheck(ConditionCheck value) => $_setField(167629711, value);
  @$pb.TagNumber(167629711)
  $core.bool hasConditioncheck() => $_has(0);
  @$pb.TagNumber(167629711)
  void clearConditioncheck() => $_clearField(167629711);
  @$pb.TagNumber(167629711)
  ConditionCheck ensureConditioncheck() => $_ensure(0);

  @$pb.TagNumber(237178517)
  Update get update => $_getN(1);
  @$pb.TagNumber(237178517)
  set update(Update value) => $_setField(237178517, value);
  @$pb.TagNumber(237178517)
  $core.bool hasUpdate() => $_has(1);
  @$pb.TagNumber(237178517)
  void clearUpdate() => $_clearField(237178517);
  @$pb.TagNumber(237178517)
  Update ensureUpdate() => $_ensure(1);

  @$pb.TagNumber(242822543)
  Put get put => $_getN(2);
  @$pb.TagNumber(242822543)
  set put(Put value) => $_setField(242822543, value);
  @$pb.TagNumber(242822543)
  $core.bool hasPut() => $_has(2);
  @$pb.TagNumber(242822543)
  void clearPut() => $_clearField(242822543);
  @$pb.TagNumber(242822543)
  Put ensurePut() => $_ensure(2);

  @$pb.TagNumber(395831915)
  Delete get delete => $_getN(3);
  @$pb.TagNumber(395831915)
  set delete(Delete value) => $_setField(395831915, value);
  @$pb.TagNumber(395831915)
  $core.bool hasDelete() => $_has(3);
  @$pb.TagNumber(395831915)
  void clearDelete() => $_clearField(395831915);
  @$pb.TagNumber(395831915)
  Delete ensureDelete() => $_ensure(3);
}

class TransactWriteItemsInput extends $pb.GeneratedMessage {
  factory TransactWriteItemsInput({
    ReturnConsumedCapacity? returnconsumedcapacity,
    ReturnItemCollectionMetrics? returnitemcollectionmetrics,
    $core.String? clientrequesttoken,
    $core.Iterable<TransactWriteItem>? transactitems,
  }) {
    final result = create();
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (returnitemcollectionmetrics != null)
      result.returnitemcollectionmetrics = returnitemcollectionmetrics;
    if (clientrequesttoken != null)
      result.clientrequesttoken = clientrequesttoken;
    if (transactitems != null) result.transactitems.addAll(transactitems);
    return result;
  }

  TransactWriteItemsInput._();

  factory TransactWriteItemsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactWriteItemsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactWriteItemsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..aE<ReturnItemCollectionMetrics>(
        255507354, _omitFieldNames ? '' : 'returnitemcollectionmetrics',
        enumValues: ReturnItemCollectionMetrics.values)
    ..aOS(455653361, _omitFieldNames ? '' : 'clientrequesttoken')
    ..pPM<TransactWriteItem>(506245290, _omitFieldNames ? '' : 'transactitems',
        subBuilder: TransactWriteItem.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactWriteItemsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactWriteItemsInput copyWith(
          void Function(TransactWriteItemsInput) updates) =>
      super.copyWith((message) => updates(message as TransactWriteItemsInput))
          as TransactWriteItemsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactWriteItemsInput create() => TransactWriteItemsInput._();
  @$core.override
  TransactWriteItemsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactWriteItemsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactWriteItemsInput>(create);
  static TransactWriteItemsInput? _defaultInstance;

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(0);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(0);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(255507354)
  ReturnItemCollectionMetrics get returnitemcollectionmetrics => $_getN(1);
  @$pb.TagNumber(255507354)
  set returnitemcollectionmetrics(ReturnItemCollectionMetrics value) =>
      $_setField(255507354, value);
  @$pb.TagNumber(255507354)
  $core.bool hasReturnitemcollectionmetrics() => $_has(1);
  @$pb.TagNumber(255507354)
  void clearReturnitemcollectionmetrics() => $_clearField(255507354);

  @$pb.TagNumber(455653361)
  $core.String get clientrequesttoken => $_getSZ(2);
  @$pb.TagNumber(455653361)
  set clientrequesttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(455653361)
  $core.bool hasClientrequesttoken() => $_has(2);
  @$pb.TagNumber(455653361)
  void clearClientrequesttoken() => $_clearField(455653361);

  @$pb.TagNumber(506245290)
  $pb.PbList<TransactWriteItem> get transactitems => $_getList(3);
}

class TransactWriteItemsOutput extends $pb.GeneratedMessage {
  factory TransactWriteItemsOutput({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        itemcollectionmetrics,
    $core.Iterable<ConsumedCapacity>? consumedcapacity,
  }) {
    final result = create();
    if (itemcollectionmetrics != null)
      result.itemcollectionmetrics.addEntries(itemcollectionmetrics);
    if (consumedcapacity != null)
      result.consumedcapacity.addAll(consumedcapacity);
    return result;
  }

  TransactWriteItemsOutput._();

  factory TransactWriteItemsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactWriteItemsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactWriteItemsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        185317452, _omitFieldNames ? '' : 'itemcollectionmetrics',
        entryClassName: 'TransactWriteItemsOutput.ItemcollectionmetricsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..pPM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactWriteItemsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactWriteItemsOutput copyWith(
          void Function(TransactWriteItemsOutput) updates) =>
      super.copyWith((message) => updates(message as TransactWriteItemsOutput))
          as TransactWriteItemsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactWriteItemsOutput create() => TransactWriteItemsOutput._();
  @$core.override
  TransactWriteItemsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactWriteItemsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactWriteItemsOutput>(create);
  static TransactWriteItemsOutput? _defaultInstance;

  @$pb.TagNumber(185317452)
  $pb.PbMap<$core.String, $core.String> get itemcollectionmetrics =>
      $_getMap(0);

  @$pb.TagNumber(449336620)
  $pb.PbList<ConsumedCapacity> get consumedcapacity => $_getList(1);
}

class TransactionCanceledException extends $pb.GeneratedMessage {
  factory TransactionCanceledException({
    $core.String? message,
    $core.Iterable<CancellationReason>? cancellationreasons,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (cancellationreasons != null)
      result.cancellationreasons.addAll(cancellationreasons);
    return result;
  }

  TransactionCanceledException._();

  factory TransactionCanceledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactionCanceledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactionCanceledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..pPM<CancellationReason>(
        346602284, _omitFieldNames ? '' : 'cancellationreasons',
        subBuilder: CancellationReason.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactionCanceledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactionCanceledException copyWith(
          void Function(TransactionCanceledException) updates) =>
      super.copyWith(
              (message) => updates(message as TransactionCanceledException))
          as TransactionCanceledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactionCanceledException create() =>
      TransactionCanceledException._();
  @$core.override
  TransactionCanceledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactionCanceledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactionCanceledException>(create);
  static TransactionCanceledException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(346602284)
  $pb.PbList<CancellationReason> get cancellationreasons => $_getList(1);
}

class TransactionConflictException extends $pb.GeneratedMessage {
  factory TransactionConflictException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TransactionConflictException._();

  factory TransactionConflictException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactionConflictException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactionConflictException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactionConflictException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactionConflictException copyWith(
          void Function(TransactionConflictException) updates) =>
      super.copyWith(
              (message) => updates(message as TransactionConflictException))
          as TransactionConflictException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactionConflictException create() =>
      TransactionConflictException._();
  @$core.override
  TransactionConflictException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactionConflictException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactionConflictException>(create);
  static TransactionConflictException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TransactionInProgressException extends $pb.GeneratedMessage {
  factory TransactionInProgressException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TransactionInProgressException._();

  factory TransactionInProgressException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactionInProgressException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactionInProgressException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactionInProgressException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactionInProgressException copyWith(
          void Function(TransactionInProgressException) updates) =>
      super.copyWith(
              (message) => updates(message as TransactionInProgressException))
          as TransactionInProgressException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactionInProgressException create() =>
      TransactionInProgressException._();
  @$core.override
  TransactionInProgressException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactionInProgressException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactionInProgressException>(create);
  static TransactionInProgressException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
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

class Update extends $pb.GeneratedMessage {
  factory Update({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? key,
    $core.String? tablename,
    $core.String? updateexpression,
    $core.String? conditionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (key != null) result.key.addEntries(key);
    if (tablename != null) result.tablename = tablename;
    if (updateexpression != null) result.updateexpression = updateexpression;
    if (conditionexpression != null)
      result.conditionexpression = conditionexpression;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    return result;
  }

  Update._();

  factory Update.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Update.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Update',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'Update.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, AttributeValue>(219859213, _omitFieldNames ? '' : 'key',
        entryClassName: 'Update.KeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(389263879, _omitFieldNames ? '' : 'updateexpression')
    ..aOS(409657405, _omitFieldNames ? '' : 'conditionexpression')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'Update.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Update clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Update copyWith(void Function(Update) updates) =>
      super.copyWith((message) => updates(message as Update)) as Update;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Update create() => Update._();
  @$core.override
  Update createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Update getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Update>(create);
  static Update? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(1);

  @$pb.TagNumber(219859213)
  $pb.PbMap<$core.String, AttributeValue> get key => $_getMap(2);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(389263879)
  $core.String get updateexpression => $_getSZ(4);
  @$pb.TagNumber(389263879)
  set updateexpression($core.String value) => $_setString(4, value);
  @$pb.TagNumber(389263879)
  $core.bool hasUpdateexpression() => $_has(4);
  @$pb.TagNumber(389263879)
  void clearUpdateexpression() => $_clearField(389263879);

  @$pb.TagNumber(409657405)
  $core.String get conditionexpression => $_getSZ(5);
  @$pb.TagNumber(409657405)
  set conditionexpression($core.String value) => $_setString(5, value);
  @$pb.TagNumber(409657405)
  $core.bool hasConditionexpression() => $_has(5);
  @$pb.TagNumber(409657405)
  void clearConditionexpression() => $_clearField(409657405);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(6);
}

class UpdateContinuousBackupsInput extends $pb.GeneratedMessage {
  factory UpdateContinuousBackupsInput({
    $core.String? tablename,
    PointInTimeRecoverySpecification? pointintimerecoveryspecification,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    if (pointintimerecoveryspecification != null)
      result.pointintimerecoveryspecification =
          pointintimerecoveryspecification;
    return result;
  }

  UpdateContinuousBackupsInput._();

  factory UpdateContinuousBackupsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateContinuousBackupsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateContinuousBackupsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<PointInTimeRecoverySpecification>(
        395539950, _omitFieldNames ? '' : 'pointintimerecoveryspecification',
        subBuilder: PointInTimeRecoverySpecification.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateContinuousBackupsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateContinuousBackupsInput copyWith(
          void Function(UpdateContinuousBackupsInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateContinuousBackupsInput))
          as UpdateContinuousBackupsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateContinuousBackupsInput create() =>
      UpdateContinuousBackupsInput._();
  @$core.override
  UpdateContinuousBackupsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateContinuousBackupsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateContinuousBackupsInput>(create);
  static UpdateContinuousBackupsInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(395539950)
  PointInTimeRecoverySpecification get pointintimerecoveryspecification =>
      $_getN(1);
  @$pb.TagNumber(395539950)
  set pointintimerecoveryspecification(
          PointInTimeRecoverySpecification value) =>
      $_setField(395539950, value);
  @$pb.TagNumber(395539950)
  $core.bool hasPointintimerecoveryspecification() => $_has(1);
  @$pb.TagNumber(395539950)
  void clearPointintimerecoveryspecification() => $_clearField(395539950);
  @$pb.TagNumber(395539950)
  PointInTimeRecoverySpecification ensurePointintimerecoveryspecification() =>
      $_ensure(1);
}

class UpdateContinuousBackupsOutput extends $pb.GeneratedMessage {
  factory UpdateContinuousBackupsOutput({
    ContinuousBackupsDescription? continuousbackupsdescription,
  }) {
    final result = create();
    if (continuousbackupsdescription != null)
      result.continuousbackupsdescription = continuousbackupsdescription;
    return result;
  }

  UpdateContinuousBackupsOutput._();

  factory UpdateContinuousBackupsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateContinuousBackupsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateContinuousBackupsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ContinuousBackupsDescription>(
        446030650, _omitFieldNames ? '' : 'continuousbackupsdescription',
        subBuilder: ContinuousBackupsDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateContinuousBackupsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateContinuousBackupsOutput copyWith(
          void Function(UpdateContinuousBackupsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateContinuousBackupsOutput))
          as UpdateContinuousBackupsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateContinuousBackupsOutput create() =>
      UpdateContinuousBackupsOutput._();
  @$core.override
  UpdateContinuousBackupsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateContinuousBackupsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateContinuousBackupsOutput>(create);
  static UpdateContinuousBackupsOutput? _defaultInstance;

  @$pb.TagNumber(446030650)
  ContinuousBackupsDescription get continuousbackupsdescription => $_getN(0);
  @$pb.TagNumber(446030650)
  set continuousbackupsdescription(ContinuousBackupsDescription value) =>
      $_setField(446030650, value);
  @$pb.TagNumber(446030650)
  $core.bool hasContinuousbackupsdescription() => $_has(0);
  @$pb.TagNumber(446030650)
  void clearContinuousbackupsdescription() => $_clearField(446030650);
  @$pb.TagNumber(446030650)
  ContinuousBackupsDescription ensureContinuousbackupsdescription() =>
      $_ensure(0);
}

class UpdateContributorInsightsInput extends $pb.GeneratedMessage {
  factory UpdateContributorInsightsInput({
    ContributorInsightsAction? contributorinsightsaction,
    ContributorInsightsMode? contributorinsightsmode,
    $core.String? indexname,
    $core.String? tablename,
  }) {
    final result = create();
    if (contributorinsightsaction != null)
      result.contributorinsightsaction = contributorinsightsaction;
    if (contributorinsightsmode != null)
      result.contributorinsightsmode = contributorinsightsmode;
    if (indexname != null) result.indexname = indexname;
    if (tablename != null) result.tablename = tablename;
    return result;
  }

  UpdateContributorInsightsInput._();

  factory UpdateContributorInsightsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateContributorInsightsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateContributorInsightsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ContributorInsightsAction>(
        81909182, _omitFieldNames ? '' : 'contributorinsightsaction',
        enumValues: ContributorInsightsAction.values)
    ..aE<ContributorInsightsMode>(
        86700161, _omitFieldNames ? '' : 'contributorinsightsmode',
        enumValues: ContributorInsightsMode.values)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateContributorInsightsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateContributorInsightsInput copyWith(
          void Function(UpdateContributorInsightsInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateContributorInsightsInput))
          as UpdateContributorInsightsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateContributorInsightsInput create() =>
      UpdateContributorInsightsInput._();
  @$core.override
  UpdateContributorInsightsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateContributorInsightsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateContributorInsightsInput>(create);
  static UpdateContributorInsightsInput? _defaultInstance;

  @$pb.TagNumber(81909182)
  ContributorInsightsAction get contributorinsightsaction => $_getN(0);
  @$pb.TagNumber(81909182)
  set contributorinsightsaction(ContributorInsightsAction value) =>
      $_setField(81909182, value);
  @$pb.TagNumber(81909182)
  $core.bool hasContributorinsightsaction() => $_has(0);
  @$pb.TagNumber(81909182)
  void clearContributorinsightsaction() => $_clearField(81909182);

  @$pb.TagNumber(86700161)
  ContributorInsightsMode get contributorinsightsmode => $_getN(1);
  @$pb.TagNumber(86700161)
  set contributorinsightsmode(ContributorInsightsMode value) =>
      $_setField(86700161, value);
  @$pb.TagNumber(86700161)
  $core.bool hasContributorinsightsmode() => $_has(1);
  @$pb.TagNumber(86700161)
  void clearContributorinsightsmode() => $_clearField(86700161);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(2);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(2);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(3);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(3, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(3);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);
}

class UpdateContributorInsightsOutput extends $pb.GeneratedMessage {
  factory UpdateContributorInsightsOutput({
    ContributorInsightsMode? contributorinsightsmode,
    $core.String? indexname,
    $core.String? tablename,
    ContributorInsightsStatus? contributorinsightsstatus,
  }) {
    final result = create();
    if (contributorinsightsmode != null)
      result.contributorinsightsmode = contributorinsightsmode;
    if (indexname != null) result.indexname = indexname;
    if (tablename != null) result.tablename = tablename;
    if (contributorinsightsstatus != null)
      result.contributorinsightsstatus = contributorinsightsstatus;
    return result;
  }

  UpdateContributorInsightsOutput._();

  factory UpdateContributorInsightsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateContributorInsightsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateContributorInsightsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ContributorInsightsMode>(
        86700161, _omitFieldNames ? '' : 'contributorinsightsmode',
        enumValues: ContributorInsightsMode.values)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aE<ContributorInsightsStatus>(
        363347282, _omitFieldNames ? '' : 'contributorinsightsstatus',
        enumValues: ContributorInsightsStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateContributorInsightsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateContributorInsightsOutput copyWith(
          void Function(UpdateContributorInsightsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateContributorInsightsOutput))
          as UpdateContributorInsightsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateContributorInsightsOutput create() =>
      UpdateContributorInsightsOutput._();
  @$core.override
  UpdateContributorInsightsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateContributorInsightsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateContributorInsightsOutput>(
          create);
  static UpdateContributorInsightsOutput? _defaultInstance;

  @$pb.TagNumber(86700161)
  ContributorInsightsMode get contributorinsightsmode => $_getN(0);
  @$pb.TagNumber(86700161)
  set contributorinsightsmode(ContributorInsightsMode value) =>
      $_setField(86700161, value);
  @$pb.TagNumber(86700161)
  $core.bool hasContributorinsightsmode() => $_has(0);
  @$pb.TagNumber(86700161)
  void clearContributorinsightsmode() => $_clearField(86700161);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(2);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(2);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(363347282)
  ContributorInsightsStatus get contributorinsightsstatus => $_getN(3);
  @$pb.TagNumber(363347282)
  set contributorinsightsstatus(ContributorInsightsStatus value) =>
      $_setField(363347282, value);
  @$pb.TagNumber(363347282)
  $core.bool hasContributorinsightsstatus() => $_has(3);
  @$pb.TagNumber(363347282)
  void clearContributorinsightsstatus() => $_clearField(363347282);
}

class UpdateGlobalSecondaryIndexAction extends $pb.GeneratedMessage {
  factory UpdateGlobalSecondaryIndexAction({
    ProvisionedThroughput? provisionedthroughput,
    $core.String? indexname,
    WarmThroughput? warmthroughput,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (indexname != null) result.indexname = indexname;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  UpdateGlobalSecondaryIndexAction._();

  factory UpdateGlobalSecondaryIndexAction.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateGlobalSecondaryIndexAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateGlobalSecondaryIndexAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughput>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughput.create)
    ..aOS(102427281, _omitFieldNames ? '' : 'indexname')
    ..aOM<WarmThroughput>(290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: WarmThroughput.create)
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalSecondaryIndexAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalSecondaryIndexAction copyWith(
          void Function(UpdateGlobalSecondaryIndexAction) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateGlobalSecondaryIndexAction))
          as UpdateGlobalSecondaryIndexAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateGlobalSecondaryIndexAction create() =>
      UpdateGlobalSecondaryIndexAction._();
  @$core.override
  UpdateGlobalSecondaryIndexAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateGlobalSecondaryIndexAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateGlobalSecondaryIndexAction>(
          create);
  static UpdateGlobalSecondaryIndexAction? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughput get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughput value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughput ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(102427281)
  $core.String get indexname => $_getSZ(1);
  @$pb.TagNumber(102427281)
  set indexname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(102427281)
  $core.bool hasIndexname() => $_has(1);
  @$pb.TagNumber(102427281)
  void clearIndexname() => $_clearField(102427281);

  @$pb.TagNumber(290598659)
  WarmThroughput get warmthroughput => $_getN(2);
  @$pb.TagNumber(290598659)
  set warmthroughput(WarmThroughput value) => $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(2);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  WarmThroughput ensureWarmthroughput() => $_ensure(2);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(3);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(3);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(3);
}

class UpdateGlobalTableInput extends $pb.GeneratedMessage {
  factory UpdateGlobalTableInput({
    $core.Iterable<ReplicaUpdate>? replicaupdates,
    $core.String? globaltablename,
  }) {
    final result = create();
    if (replicaupdates != null) result.replicaupdates.addAll(replicaupdates);
    if (globaltablename != null) result.globaltablename = globaltablename;
    return result;
  }

  UpdateGlobalTableInput._();

  factory UpdateGlobalTableInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateGlobalTableInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateGlobalTableInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<ReplicaUpdate>(260731936, _omitFieldNames ? '' : 'replicaupdates',
        subBuilder: ReplicaUpdate.create)
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalTableInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalTableInput copyWith(
          void Function(UpdateGlobalTableInput) updates) =>
      super.copyWith((message) => updates(message as UpdateGlobalTableInput))
          as UpdateGlobalTableInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateGlobalTableInput create() => UpdateGlobalTableInput._();
  @$core.override
  UpdateGlobalTableInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateGlobalTableInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateGlobalTableInput>(create);
  static UpdateGlobalTableInput? _defaultInstance;

  @$pb.TagNumber(260731936)
  $pb.PbList<ReplicaUpdate> get replicaupdates => $_getList(0);

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(1);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(1);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);
}

class UpdateGlobalTableOutput extends $pb.GeneratedMessage {
  factory UpdateGlobalTableOutput({
    GlobalTableDescription? globaltabledescription,
  }) {
    final result = create();
    if (globaltabledescription != null)
      result.globaltabledescription = globaltabledescription;
    return result;
  }

  UpdateGlobalTableOutput._();

  factory UpdateGlobalTableOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateGlobalTableOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateGlobalTableOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<GlobalTableDescription>(
        342116405, _omitFieldNames ? '' : 'globaltabledescription',
        subBuilder: GlobalTableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalTableOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalTableOutput copyWith(
          void Function(UpdateGlobalTableOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateGlobalTableOutput))
          as UpdateGlobalTableOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateGlobalTableOutput create() => UpdateGlobalTableOutput._();
  @$core.override
  UpdateGlobalTableOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateGlobalTableOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateGlobalTableOutput>(create);
  static UpdateGlobalTableOutput? _defaultInstance;

  @$pb.TagNumber(342116405)
  GlobalTableDescription get globaltabledescription => $_getN(0);
  @$pb.TagNumber(342116405)
  set globaltabledescription(GlobalTableDescription value) =>
      $_setField(342116405, value);
  @$pb.TagNumber(342116405)
  $core.bool hasGlobaltabledescription() => $_has(0);
  @$pb.TagNumber(342116405)
  void clearGlobaltabledescription() => $_clearField(342116405);
  @$pb.TagNumber(342116405)
  GlobalTableDescription ensureGlobaltabledescription() => $_ensure(0);
}

class UpdateGlobalTableSettingsInput extends $pb.GeneratedMessage {
  factory UpdateGlobalTableSettingsInput({
    $fixnum.Int64? globaltableprovisionedwritecapacityunits,
    $core.String? globaltablename,
    BillingMode? globaltablebillingmode,
    $core.Iterable<GlobalTableGlobalSecondaryIndexSettingsUpdate>?
        globaltableglobalsecondaryindexsettingsupdate,
    AutoScalingSettingsUpdate?
        globaltableprovisionedwritecapacityautoscalingsettingsupdate,
    $core.Iterable<ReplicaSettingsUpdate>? replicasettingsupdate,
  }) {
    final result = create();
    if (globaltableprovisionedwritecapacityunits != null)
      result.globaltableprovisionedwritecapacityunits =
          globaltableprovisionedwritecapacityunits;
    if (globaltablename != null) result.globaltablename = globaltablename;
    if (globaltablebillingmode != null)
      result.globaltablebillingmode = globaltablebillingmode;
    if (globaltableglobalsecondaryindexsettingsupdate != null)
      result.globaltableglobalsecondaryindexsettingsupdate
          .addAll(globaltableglobalsecondaryindexsettingsupdate);
    if (globaltableprovisionedwritecapacityautoscalingsettingsupdate != null)
      result.globaltableprovisionedwritecapacityautoscalingsettingsupdate =
          globaltableprovisionedwritecapacityautoscalingsettingsupdate;
    if (replicasettingsupdate != null)
      result.replicasettingsupdate.addAll(replicasettingsupdate);
    return result;
  }

  UpdateGlobalTableSettingsInput._();

  factory UpdateGlobalTableSettingsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateGlobalTableSettingsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateGlobalTableSettingsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(127172283,
        _omitFieldNames ? '' : 'globaltableprovisionedwritecapacityunits')
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..aE<BillingMode>(
        285868859, _omitFieldNames ? '' : 'globaltablebillingmode',
        enumValues: BillingMode.values)
    ..pPM<GlobalTableGlobalSecondaryIndexSettingsUpdate>(368411116,
        _omitFieldNames ? '' : 'globaltableglobalsecondaryindexsettingsupdate',
        subBuilder: GlobalTableGlobalSecondaryIndexSettingsUpdate.create)
    ..aOM<AutoScalingSettingsUpdate>(
        435219182,
        _omitFieldNames
            ? ''
            : 'globaltableprovisionedwritecapacityautoscalingsettingsupdate',
        subBuilder: AutoScalingSettingsUpdate.create)
    ..pPM<ReplicaSettingsUpdate>(
        463498612, _omitFieldNames ? '' : 'replicasettingsupdate',
        subBuilder: ReplicaSettingsUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalTableSettingsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalTableSettingsInput copyWith(
          void Function(UpdateGlobalTableSettingsInput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateGlobalTableSettingsInput))
          as UpdateGlobalTableSettingsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateGlobalTableSettingsInput create() =>
      UpdateGlobalTableSettingsInput._();
  @$core.override
  UpdateGlobalTableSettingsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateGlobalTableSettingsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateGlobalTableSettingsInput>(create);
  static UpdateGlobalTableSettingsInput? _defaultInstance;

  @$pb.TagNumber(127172283)
  $fixnum.Int64 get globaltableprovisionedwritecapacityunits => $_getI64(0);
  @$pb.TagNumber(127172283)
  set globaltableprovisionedwritecapacityunits($fixnum.Int64 value) =>
      $_setInt64(0, value);
  @$pb.TagNumber(127172283)
  $core.bool hasGlobaltableprovisionedwritecapacityunits() => $_has(0);
  @$pb.TagNumber(127172283)
  void clearGlobaltableprovisionedwritecapacityunits() =>
      $_clearField(127172283);

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(1);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(1);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);

  @$pb.TagNumber(285868859)
  BillingMode get globaltablebillingmode => $_getN(2);
  @$pb.TagNumber(285868859)
  set globaltablebillingmode(BillingMode value) => $_setField(285868859, value);
  @$pb.TagNumber(285868859)
  $core.bool hasGlobaltablebillingmode() => $_has(2);
  @$pb.TagNumber(285868859)
  void clearGlobaltablebillingmode() => $_clearField(285868859);

  @$pb.TagNumber(368411116)
  $pb.PbList<GlobalTableGlobalSecondaryIndexSettingsUpdate>
      get globaltableglobalsecondaryindexsettingsupdate => $_getList(3);

  @$pb.TagNumber(435219182)
  AutoScalingSettingsUpdate
      get globaltableprovisionedwritecapacityautoscalingsettingsupdate =>
          $_getN(4);
  @$pb.TagNumber(435219182)
  set globaltableprovisionedwritecapacityautoscalingsettingsupdate(
          AutoScalingSettingsUpdate value) =>
      $_setField(435219182, value);
  @$pb.TagNumber(435219182)
  $core.bool
      hasGlobaltableprovisionedwritecapacityautoscalingsettingsupdate() =>
          $_has(4);
  @$pb.TagNumber(435219182)
  void clearGlobaltableprovisionedwritecapacityautoscalingsettingsupdate() =>
      $_clearField(435219182);
  @$pb.TagNumber(435219182)
  AutoScalingSettingsUpdate
      ensureGlobaltableprovisionedwritecapacityautoscalingsettingsupdate() =>
          $_ensure(4);

  @$pb.TagNumber(463498612)
  $pb.PbList<ReplicaSettingsUpdate> get replicasettingsupdate => $_getList(5);
}

class UpdateGlobalTableSettingsOutput extends $pb.GeneratedMessage {
  factory UpdateGlobalTableSettingsOutput({
    $core.String? globaltablename,
    $core.Iterable<ReplicaSettingsDescription>? replicasettings,
  }) {
    final result = create();
    if (globaltablename != null) result.globaltablename = globaltablename;
    if (replicasettings != null) result.replicasettings.addAll(replicasettings);
    return result;
  }

  UpdateGlobalTableSettingsOutput._();

  factory UpdateGlobalTableSettingsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateGlobalTableSettingsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateGlobalTableSettingsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(283759402, _omitFieldNames ? '' : 'globaltablename')
    ..pPM<ReplicaSettingsDescription>(
        288882107, _omitFieldNames ? '' : 'replicasettings',
        subBuilder: ReplicaSettingsDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalTableSettingsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGlobalTableSettingsOutput copyWith(
          void Function(UpdateGlobalTableSettingsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateGlobalTableSettingsOutput))
          as UpdateGlobalTableSettingsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateGlobalTableSettingsOutput create() =>
      UpdateGlobalTableSettingsOutput._();
  @$core.override
  UpdateGlobalTableSettingsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateGlobalTableSettingsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateGlobalTableSettingsOutput>(
          create);
  static UpdateGlobalTableSettingsOutput? _defaultInstance;

  @$pb.TagNumber(283759402)
  $core.String get globaltablename => $_getSZ(0);
  @$pb.TagNumber(283759402)
  set globaltablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(283759402)
  $core.bool hasGlobaltablename() => $_has(0);
  @$pb.TagNumber(283759402)
  void clearGlobaltablename() => $_clearField(283759402);

  @$pb.TagNumber(288882107)
  $pb.PbList<ReplicaSettingsDescription> get replicasettings => $_getList(1);
}

class UpdateItemInput extends $pb.GeneratedMessage {
  factory UpdateItemInput({
    ReturnValuesOnConditionCheckFailure? returnvaluesonconditioncheckfailure,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValueUpdate>>?
        attributeupdates,
    ReturnConsumedCapacity? returnconsumedcapacity,
    $core.Iterable<$core.MapEntry<$core.String, ExpectedAttributeValue>>?
        expected,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        expressionattributenames,
    ConditionalOperator? conditionaloperator,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? key,
    ReturnItemCollectionMetrics? returnitemcollectionmetrics,
    $core.String? tablename,
    $core.String? updateexpression,
    ReturnValue? returnvalues,
    $core.String? conditionexpression,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>?
        expressionattributevalues,
  }) {
    final result = create();
    if (returnvaluesonconditioncheckfailure != null)
      result.returnvaluesonconditioncheckfailure =
          returnvaluesonconditioncheckfailure;
    if (attributeupdates != null)
      result.attributeupdates.addEntries(attributeupdates);
    if (returnconsumedcapacity != null)
      result.returnconsumedcapacity = returnconsumedcapacity;
    if (expected != null) result.expected.addEntries(expected);
    if (expressionattributenames != null)
      result.expressionattributenames.addEntries(expressionattributenames);
    if (conditionaloperator != null)
      result.conditionaloperator = conditionaloperator;
    if (key != null) result.key.addEntries(key);
    if (returnitemcollectionmetrics != null)
      result.returnitemcollectionmetrics = returnitemcollectionmetrics;
    if (tablename != null) result.tablename = tablename;
    if (updateexpression != null) result.updateexpression = updateexpression;
    if (returnvalues != null) result.returnvalues = returnvalues;
    if (conditionexpression != null)
      result.conditionexpression = conditionexpression;
    if (expressionattributevalues != null)
      result.expressionattributevalues.addEntries(expressionattributevalues);
    return result;
  }

  UpdateItemInput._();

  factory UpdateItemInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateItemInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateItemInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ReturnValuesOnConditionCheckFailure>(
        4213728, _omitFieldNames ? '' : 'returnvaluesonconditioncheckfailure',
        enumValues: ReturnValuesOnConditionCheckFailure.values)
    ..m<$core.String, AttributeValueUpdate>(
        32403512, _omitFieldNames ? '' : 'attributeupdates',
        entryClassName: 'UpdateItemInput.AttributeupdatesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValueUpdate.create,
        valueDefaultOrMaker: AttributeValueUpdate.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<ReturnConsumedCapacity>(
        43545598, _omitFieldNames ? '' : 'returnconsumedcapacity',
        enumValues: ReturnConsumedCapacity.values)
    ..m<$core.String, ExpectedAttributeValue>(
        106557946, _omitFieldNames ? '' : 'expected',
        entryClassName: 'UpdateItemInput.ExpectedEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: ExpectedAttributeValue.create,
        valueDefaultOrMaker: ExpectedAttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..m<$core.String, $core.String>(
        150228092, _omitFieldNames ? '' : 'expressionattributenames',
        entryClassName: 'UpdateItemInput.ExpressionattributenamesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<ConditionalOperator>(
        172066260, _omitFieldNames ? '' : 'conditionaloperator',
        enumValues: ConditionalOperator.values)
    ..m<$core.String, AttributeValue>(219859213, _omitFieldNames ? '' : 'key',
        entryClassName: 'UpdateItemInput.KeyEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aE<ReturnItemCollectionMetrics>(
        255507354, _omitFieldNames ? '' : 'returnitemcollectionmetrics',
        enumValues: ReturnItemCollectionMetrics.values)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(389263879, _omitFieldNames ? '' : 'updateexpression')
    ..aE<ReturnValue>(402960198, _omitFieldNames ? '' : 'returnvalues',
        enumValues: ReturnValue.values)
    ..aOS(409657405, _omitFieldNames ? '' : 'conditionexpression')
    ..m<$core.String, AttributeValue>(
        484970072, _omitFieldNames ? '' : 'expressionattributevalues',
        entryClassName: 'UpdateItemInput.ExpressionattributevaluesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateItemInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateItemInput copyWith(void Function(UpdateItemInput) updates) =>
      super.copyWith((message) => updates(message as UpdateItemInput))
          as UpdateItemInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateItemInput create() => UpdateItemInput._();
  @$core.override
  UpdateItemInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateItemInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateItemInput>(create);
  static UpdateItemInput? _defaultInstance;

  @$pb.TagNumber(4213728)
  ReturnValuesOnConditionCheckFailure get returnvaluesonconditioncheckfailure =>
      $_getN(0);
  @$pb.TagNumber(4213728)
  set returnvaluesonconditioncheckfailure(
          ReturnValuesOnConditionCheckFailure value) =>
      $_setField(4213728, value);
  @$pb.TagNumber(4213728)
  $core.bool hasReturnvaluesonconditioncheckfailure() => $_has(0);
  @$pb.TagNumber(4213728)
  void clearReturnvaluesonconditioncheckfailure() => $_clearField(4213728);

  @$pb.TagNumber(32403512)
  $pb.PbMap<$core.String, AttributeValueUpdate> get attributeupdates =>
      $_getMap(1);

  @$pb.TagNumber(43545598)
  ReturnConsumedCapacity get returnconsumedcapacity => $_getN(2);
  @$pb.TagNumber(43545598)
  set returnconsumedcapacity(ReturnConsumedCapacity value) =>
      $_setField(43545598, value);
  @$pb.TagNumber(43545598)
  $core.bool hasReturnconsumedcapacity() => $_has(2);
  @$pb.TagNumber(43545598)
  void clearReturnconsumedcapacity() => $_clearField(43545598);

  @$pb.TagNumber(106557946)
  $pb.PbMap<$core.String, ExpectedAttributeValue> get expected => $_getMap(3);

  @$pb.TagNumber(150228092)
  $pb.PbMap<$core.String, $core.String> get expressionattributenames =>
      $_getMap(4);

  @$pb.TagNumber(172066260)
  ConditionalOperator get conditionaloperator => $_getN(5);
  @$pb.TagNumber(172066260)
  set conditionaloperator(ConditionalOperator value) =>
      $_setField(172066260, value);
  @$pb.TagNumber(172066260)
  $core.bool hasConditionaloperator() => $_has(5);
  @$pb.TagNumber(172066260)
  void clearConditionaloperator() => $_clearField(172066260);

  @$pb.TagNumber(219859213)
  $pb.PbMap<$core.String, AttributeValue> get key => $_getMap(6);

  @$pb.TagNumber(255507354)
  ReturnItemCollectionMetrics get returnitemcollectionmetrics => $_getN(7);
  @$pb.TagNumber(255507354)
  set returnitemcollectionmetrics(ReturnItemCollectionMetrics value) =>
      $_setField(255507354, value);
  @$pb.TagNumber(255507354)
  $core.bool hasReturnitemcollectionmetrics() => $_has(7);
  @$pb.TagNumber(255507354)
  void clearReturnitemcollectionmetrics() => $_clearField(255507354);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(8);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(8, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(8);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(389263879)
  $core.String get updateexpression => $_getSZ(9);
  @$pb.TagNumber(389263879)
  set updateexpression($core.String value) => $_setString(9, value);
  @$pb.TagNumber(389263879)
  $core.bool hasUpdateexpression() => $_has(9);
  @$pb.TagNumber(389263879)
  void clearUpdateexpression() => $_clearField(389263879);

  @$pb.TagNumber(402960198)
  ReturnValue get returnvalues => $_getN(10);
  @$pb.TagNumber(402960198)
  set returnvalues(ReturnValue value) => $_setField(402960198, value);
  @$pb.TagNumber(402960198)
  $core.bool hasReturnvalues() => $_has(10);
  @$pb.TagNumber(402960198)
  void clearReturnvalues() => $_clearField(402960198);

  @$pb.TagNumber(409657405)
  $core.String get conditionexpression => $_getSZ(11);
  @$pb.TagNumber(409657405)
  set conditionexpression($core.String value) => $_setString(11, value);
  @$pb.TagNumber(409657405)
  $core.bool hasConditionexpression() => $_has(11);
  @$pb.TagNumber(409657405)
  void clearConditionexpression() => $_clearField(409657405);

  @$pb.TagNumber(484970072)
  $pb.PbMap<$core.String, AttributeValue> get expressionattributevalues =>
      $_getMap(12);
}

class UpdateItemOutput extends $pb.GeneratedMessage {
  factory UpdateItemOutput({
    ItemCollectionMetrics? itemcollectionmetrics,
    $core.Iterable<$core.MapEntry<$core.String, AttributeValue>>? attributes,
    ConsumedCapacity? consumedcapacity,
  }) {
    final result = create();
    if (itemcollectionmetrics != null)
      result.itemcollectionmetrics = itemcollectionmetrics;
    if (attributes != null) result.attributes.addEntries(attributes);
    if (consumedcapacity != null) result.consumedcapacity = consumedcapacity;
    return result;
  }

  UpdateItemOutput._();

  factory UpdateItemOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateItemOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateItemOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ItemCollectionMetrics>(
        185317452, _omitFieldNames ? '' : 'itemcollectionmetrics',
        subBuilder: ItemCollectionMetrics.create)
    ..m<$core.String, AttributeValue>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'UpdateItemOutput.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: AttributeValue.create,
        valueDefaultOrMaker: AttributeValue.getDefault,
        packageName: const $pb.PackageName('dynamodb'))
    ..aOM<ConsumedCapacity>(
        449336620, _omitFieldNames ? '' : 'consumedcapacity',
        subBuilder: ConsumedCapacity.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateItemOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateItemOutput copyWith(void Function(UpdateItemOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateItemOutput))
          as UpdateItemOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateItemOutput create() => UpdateItemOutput._();
  @$core.override
  UpdateItemOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateItemOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateItemOutput>(create);
  static UpdateItemOutput? _defaultInstance;

  @$pb.TagNumber(185317452)
  ItemCollectionMetrics get itemcollectionmetrics => $_getN(0);
  @$pb.TagNumber(185317452)
  set itemcollectionmetrics(ItemCollectionMetrics value) =>
      $_setField(185317452, value);
  @$pb.TagNumber(185317452)
  $core.bool hasItemcollectionmetrics() => $_has(0);
  @$pb.TagNumber(185317452)
  void clearItemcollectionmetrics() => $_clearField(185317452);
  @$pb.TagNumber(185317452)
  ItemCollectionMetrics ensureItemcollectionmetrics() => $_ensure(0);

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, AttributeValue> get attributes => $_getMap(1);

  @$pb.TagNumber(449336620)
  ConsumedCapacity get consumedcapacity => $_getN(2);
  @$pb.TagNumber(449336620)
  set consumedcapacity(ConsumedCapacity value) => $_setField(449336620, value);
  @$pb.TagNumber(449336620)
  $core.bool hasConsumedcapacity() => $_has(2);
  @$pb.TagNumber(449336620)
  void clearConsumedcapacity() => $_clearField(449336620);
  @$pb.TagNumber(449336620)
  ConsumedCapacity ensureConsumedcapacity() => $_ensure(2);
}

class UpdateKinesisStreamingConfiguration extends $pb.GeneratedMessage {
  factory UpdateKinesisStreamingConfiguration({
    ApproximateCreationDateTimePrecision? approximatecreationdatetimeprecision,
  }) {
    final result = create();
    if (approximatecreationdatetimeprecision != null)
      result.approximatecreationdatetimeprecision =
          approximatecreationdatetimeprecision;
    return result;
  }

  UpdateKinesisStreamingConfiguration._();

  factory UpdateKinesisStreamingConfiguration.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateKinesisStreamingConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateKinesisStreamingConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aE<ApproximateCreationDateTimePrecision>(392293352,
        _omitFieldNames ? '' : 'approximatecreationdatetimeprecision',
        enumValues: ApproximateCreationDateTimePrecision.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateKinesisStreamingConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateKinesisStreamingConfiguration copyWith(
          void Function(UpdateKinesisStreamingConfiguration) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateKinesisStreamingConfiguration))
          as UpdateKinesisStreamingConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateKinesisStreamingConfiguration create() =>
      UpdateKinesisStreamingConfiguration._();
  @$core.override
  UpdateKinesisStreamingConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateKinesisStreamingConfiguration getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateKinesisStreamingConfiguration>(create);
  static UpdateKinesisStreamingConfiguration? _defaultInstance;

  @$pb.TagNumber(392293352)
  ApproximateCreationDateTimePrecision
      get approximatecreationdatetimeprecision => $_getN(0);
  @$pb.TagNumber(392293352)
  set approximatecreationdatetimeprecision(
          ApproximateCreationDateTimePrecision value) =>
      $_setField(392293352, value);
  @$pb.TagNumber(392293352)
  $core.bool hasApproximatecreationdatetimeprecision() => $_has(0);
  @$pb.TagNumber(392293352)
  void clearApproximatecreationdatetimeprecision() => $_clearField(392293352);
}

class UpdateKinesisStreamingDestinationInput extends $pb.GeneratedMessage {
  factory UpdateKinesisStreamingDestinationInput({
    UpdateKinesisStreamingConfiguration? updatekinesisstreamingconfiguration,
    $core.String? tablename,
    $core.String? streamarn,
  }) {
    final result = create();
    if (updatekinesisstreamingconfiguration != null)
      result.updatekinesisstreamingconfiguration =
          updatekinesisstreamingconfiguration;
    if (tablename != null) result.tablename = tablename;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  UpdateKinesisStreamingDestinationInput._();

  factory UpdateKinesisStreamingDestinationInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateKinesisStreamingDestinationInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateKinesisStreamingDestinationInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<UpdateKinesisStreamingConfiguration>(
        239134845, _omitFieldNames ? '' : 'updatekinesisstreamingconfiguration',
        subBuilder: UpdateKinesisStreamingConfiguration.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOS(513423709, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateKinesisStreamingDestinationInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateKinesisStreamingDestinationInput copyWith(
          void Function(UpdateKinesisStreamingDestinationInput) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateKinesisStreamingDestinationInput))
          as UpdateKinesisStreamingDestinationInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateKinesisStreamingDestinationInput create() =>
      UpdateKinesisStreamingDestinationInput._();
  @$core.override
  UpdateKinesisStreamingDestinationInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateKinesisStreamingDestinationInput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateKinesisStreamingDestinationInput>(create);
  static UpdateKinesisStreamingDestinationInput? _defaultInstance;

  @$pb.TagNumber(239134845)
  UpdateKinesisStreamingConfiguration get updatekinesisstreamingconfiguration =>
      $_getN(0);
  @$pb.TagNumber(239134845)
  set updatekinesisstreamingconfiguration(
          UpdateKinesisStreamingConfiguration value) =>
      $_setField(239134845, value);
  @$pb.TagNumber(239134845)
  $core.bool hasUpdatekinesisstreamingconfiguration() => $_has(0);
  @$pb.TagNumber(239134845)
  void clearUpdatekinesisstreamingconfiguration() => $_clearField(239134845);
  @$pb.TagNumber(239134845)
  UpdateKinesisStreamingConfiguration
      ensureUpdatekinesisstreamingconfiguration() => $_ensure(0);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(513423709)
  $core.String get streamarn => $_getSZ(2);
  @$pb.TagNumber(513423709)
  set streamarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(513423709)
  $core.bool hasStreamarn() => $_has(2);
  @$pb.TagNumber(513423709)
  void clearStreamarn() => $_clearField(513423709);
}

class UpdateKinesisStreamingDestinationOutput extends $pb.GeneratedMessage {
  factory UpdateKinesisStreamingDestinationOutput({
    UpdateKinesisStreamingConfiguration? updatekinesisstreamingconfiguration,
    $core.String? tablename,
    DestinationStatus? destinationstatus,
    $core.String? streamarn,
  }) {
    final result = create();
    if (updatekinesisstreamingconfiguration != null)
      result.updatekinesisstreamingconfiguration =
          updatekinesisstreamingconfiguration;
    if (tablename != null) result.tablename = tablename;
    if (destinationstatus != null) result.destinationstatus = destinationstatus;
    if (streamarn != null) result.streamarn = streamarn;
    return result;
  }

  UpdateKinesisStreamingDestinationOutput._();

  factory UpdateKinesisStreamingDestinationOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateKinesisStreamingDestinationOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateKinesisStreamingDestinationOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<UpdateKinesisStreamingConfiguration>(
        239134845, _omitFieldNames ? '' : 'updatekinesisstreamingconfiguration',
        subBuilder: UpdateKinesisStreamingConfiguration.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aE<DestinationStatus>(
        381248234, _omitFieldNames ? '' : 'destinationstatus',
        enumValues: DestinationStatus.values)
    ..aOS(513423709, _omitFieldNames ? '' : 'streamarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateKinesisStreamingDestinationOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateKinesisStreamingDestinationOutput copyWith(
          void Function(UpdateKinesisStreamingDestinationOutput) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateKinesisStreamingDestinationOutput))
          as UpdateKinesisStreamingDestinationOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateKinesisStreamingDestinationOutput create() =>
      UpdateKinesisStreamingDestinationOutput._();
  @$core.override
  UpdateKinesisStreamingDestinationOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateKinesisStreamingDestinationOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateKinesisStreamingDestinationOutput>(create);
  static UpdateKinesisStreamingDestinationOutput? _defaultInstance;

  @$pb.TagNumber(239134845)
  UpdateKinesisStreamingConfiguration get updatekinesisstreamingconfiguration =>
      $_getN(0);
  @$pb.TagNumber(239134845)
  set updatekinesisstreamingconfiguration(
          UpdateKinesisStreamingConfiguration value) =>
      $_setField(239134845, value);
  @$pb.TagNumber(239134845)
  $core.bool hasUpdatekinesisstreamingconfiguration() => $_has(0);
  @$pb.TagNumber(239134845)
  void clearUpdatekinesisstreamingconfiguration() => $_clearField(239134845);
  @$pb.TagNumber(239134845)
  UpdateKinesisStreamingConfiguration
      ensureUpdatekinesisstreamingconfiguration() => $_ensure(0);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(1);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(1);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(381248234)
  DestinationStatus get destinationstatus => $_getN(2);
  @$pb.TagNumber(381248234)
  set destinationstatus(DestinationStatus value) =>
      $_setField(381248234, value);
  @$pb.TagNumber(381248234)
  $core.bool hasDestinationstatus() => $_has(2);
  @$pb.TagNumber(381248234)
  void clearDestinationstatus() => $_clearField(381248234);

  @$pb.TagNumber(513423709)
  $core.String get streamarn => $_getSZ(3);
  @$pb.TagNumber(513423709)
  set streamarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(513423709)
  $core.bool hasStreamarn() => $_has(3);
  @$pb.TagNumber(513423709)
  void clearStreamarn() => $_clearField(513423709);
}

class UpdateReplicationGroupMemberAction extends $pb.GeneratedMessage {
  factory UpdateReplicationGroupMemberAction({
    $core.String? regionname,
    OnDemandThroughputOverride? ondemandthroughputoverride,
    $core.Iterable<ReplicaGlobalSecondaryIndex>? globalsecondaryindexes,
    ProvisionedThroughputOverride? provisionedthroughputoverride,
    TableClass? tableclassoverride,
    $core.String? kmsmasterkeyid,
  }) {
    final result = create();
    if (regionname != null) result.regionname = regionname;
    if (ondemandthroughputoverride != null)
      result.ondemandthroughputoverride = ondemandthroughputoverride;
    if (globalsecondaryindexes != null)
      result.globalsecondaryindexes.addAll(globalsecondaryindexes);
    if (provisionedthroughputoverride != null)
      result.provisionedthroughputoverride = provisionedthroughputoverride;
    if (tableclassoverride != null)
      result.tableclassoverride = tableclassoverride;
    if (kmsmasterkeyid != null) result.kmsmasterkeyid = kmsmasterkeyid;
    return result;
  }

  UpdateReplicationGroupMemberAction._();

  factory UpdateReplicationGroupMemberAction.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateReplicationGroupMemberAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateReplicationGroupMemberAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(112086463, _omitFieldNames ? '' : 'regionname')
    ..aOM<OnDemandThroughputOverride>(
        317165234, _omitFieldNames ? '' : 'ondemandthroughputoverride',
        subBuilder: OnDemandThroughputOverride.create)
    ..pPM<ReplicaGlobalSecondaryIndex>(
        409156905, _omitFieldNames ? '' : 'globalsecondaryindexes',
        subBuilder: ReplicaGlobalSecondaryIndex.create)
    ..aOM<ProvisionedThroughputOverride>(
        413332116, _omitFieldNames ? '' : 'provisionedthroughputoverride',
        subBuilder: ProvisionedThroughputOverride.create)
    ..aE<TableClass>(415569842, _omitFieldNames ? '' : 'tableclassoverride',
        enumValues: TableClass.values)
    ..aOS(521459443, _omitFieldNames ? '' : 'kmsmasterkeyid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateReplicationGroupMemberAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateReplicationGroupMemberAction copyWith(
          void Function(UpdateReplicationGroupMemberAction) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateReplicationGroupMemberAction))
          as UpdateReplicationGroupMemberAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateReplicationGroupMemberAction create() =>
      UpdateReplicationGroupMemberAction._();
  @$core.override
  UpdateReplicationGroupMemberAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateReplicationGroupMemberAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateReplicationGroupMemberAction>(
          create);
  static UpdateReplicationGroupMemberAction? _defaultInstance;

  @$pb.TagNumber(112086463)
  $core.String get regionname => $_getSZ(0);
  @$pb.TagNumber(112086463)
  set regionname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(112086463)
  $core.bool hasRegionname() => $_has(0);
  @$pb.TagNumber(112086463)
  void clearRegionname() => $_clearField(112086463);

  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride get ondemandthroughputoverride => $_getN(1);
  @$pb.TagNumber(317165234)
  set ondemandthroughputoverride(OnDemandThroughputOverride value) =>
      $_setField(317165234, value);
  @$pb.TagNumber(317165234)
  $core.bool hasOndemandthroughputoverride() => $_has(1);
  @$pb.TagNumber(317165234)
  void clearOndemandthroughputoverride() => $_clearField(317165234);
  @$pb.TagNumber(317165234)
  OnDemandThroughputOverride ensureOndemandthroughputoverride() => $_ensure(1);

  @$pb.TagNumber(409156905)
  $pb.PbList<ReplicaGlobalSecondaryIndex> get globalsecondaryindexes =>
      $_getList(2);

  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride get provisionedthroughputoverride => $_getN(3);
  @$pb.TagNumber(413332116)
  set provisionedthroughputoverride(ProvisionedThroughputOverride value) =>
      $_setField(413332116, value);
  @$pb.TagNumber(413332116)
  $core.bool hasProvisionedthroughputoverride() => $_has(3);
  @$pb.TagNumber(413332116)
  void clearProvisionedthroughputoverride() => $_clearField(413332116);
  @$pb.TagNumber(413332116)
  ProvisionedThroughputOverride ensureProvisionedthroughputoverride() =>
      $_ensure(3);

  @$pb.TagNumber(415569842)
  TableClass get tableclassoverride => $_getN(4);
  @$pb.TagNumber(415569842)
  set tableclassoverride(TableClass value) => $_setField(415569842, value);
  @$pb.TagNumber(415569842)
  $core.bool hasTableclassoverride() => $_has(4);
  @$pb.TagNumber(415569842)
  void clearTableclassoverride() => $_clearField(415569842);

  @$pb.TagNumber(521459443)
  $core.String get kmsmasterkeyid => $_getSZ(5);
  @$pb.TagNumber(521459443)
  set kmsmasterkeyid($core.String value) => $_setString(5, value);
  @$pb.TagNumber(521459443)
  $core.bool hasKmsmasterkeyid() => $_has(5);
  @$pb.TagNumber(521459443)
  void clearKmsmasterkeyid() => $_clearField(521459443);
}

class UpdateTableInput extends $pb.GeneratedMessage {
  factory UpdateTableInput({
    ProvisionedThroughput? provisionedthroughput,
    GlobalTableSettingsReplicationMode? globaltablesettingsreplicationmode,
    SSESpecification? ssespecification,
    BillingMode? billingmode,
    $core.bool? deletionprotectionenabled,
    $core.Iterable<ReplicationGroupUpdate>? replicaupdates,
    $core.Iterable<GlobalSecondaryIndexUpdate>? globalsecondaryindexupdates,
    $core.Iterable<GlobalTableWitnessGroupUpdate>? globaltablewitnessupdates,
    $core.String? tablename,
    WarmThroughput? warmthroughput,
    TableClass? tableclass,
    StreamSpecification? streamspecification,
    $core.Iterable<AttributeDefinition>? attributedefinitions,
    MultiRegionConsistency? multiregionconsistency,
    OnDemandThroughput? ondemandthroughput,
  }) {
    final result = create();
    if (provisionedthroughput != null)
      result.provisionedthroughput = provisionedthroughput;
    if (globaltablesettingsreplicationmode != null)
      result.globaltablesettingsreplicationmode =
          globaltablesettingsreplicationmode;
    if (ssespecification != null) result.ssespecification = ssespecification;
    if (billingmode != null) result.billingmode = billingmode;
    if (deletionprotectionenabled != null)
      result.deletionprotectionenabled = deletionprotectionenabled;
    if (replicaupdates != null) result.replicaupdates.addAll(replicaupdates);
    if (globalsecondaryindexupdates != null)
      result.globalsecondaryindexupdates.addAll(globalsecondaryindexupdates);
    if (globaltablewitnessupdates != null)
      result.globaltablewitnessupdates.addAll(globaltablewitnessupdates);
    if (tablename != null) result.tablename = tablename;
    if (warmthroughput != null) result.warmthroughput = warmthroughput;
    if (tableclass != null) result.tableclass = tableclass;
    if (streamspecification != null)
      result.streamspecification = streamspecification;
    if (attributedefinitions != null)
      result.attributedefinitions.addAll(attributedefinitions);
    if (multiregionconsistency != null)
      result.multiregionconsistency = multiregionconsistency;
    if (ondemandthroughput != null)
      result.ondemandthroughput = ondemandthroughput;
    return result;
  }

  UpdateTableInput._();

  factory UpdateTableInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTableInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTableInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<ProvisionedThroughput>(
        1757580, _omitFieldNames ? '' : 'provisionedthroughput',
        subBuilder: ProvisionedThroughput.create)
    ..aE<GlobalTableSettingsReplicationMode>(
        10446577, _omitFieldNames ? '' : 'globaltablesettingsreplicationmode',
        enumValues: GlobalTableSettingsReplicationMode.values)
    ..aOM<SSESpecification>(31692444, _omitFieldNames ? '' : 'ssespecification',
        subBuilder: SSESpecification.create)
    ..aE<BillingMode>(184162880, _omitFieldNames ? '' : 'billingmode',
        enumValues: BillingMode.values)
    ..aOB(259418450, _omitFieldNames ? '' : 'deletionprotectionenabled')
    ..pPM<ReplicationGroupUpdate>(
        260731936, _omitFieldNames ? '' : 'replicaupdates',
        subBuilder: ReplicationGroupUpdate.create)
    ..pPM<GlobalSecondaryIndexUpdate>(
        265760923, _omitFieldNames ? '' : 'globalsecondaryindexupdates',
        subBuilder: GlobalSecondaryIndexUpdate.create)
    ..pPM<GlobalTableWitnessGroupUpdate>(
        269201458, _omitFieldNames ? '' : 'globaltablewitnessupdates',
        subBuilder: GlobalTableWitnessGroupUpdate.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<WarmThroughput>(290598659, _omitFieldNames ? '' : 'warmthroughput',
        subBuilder: WarmThroughput.create)
    ..aE<TableClass>(342890498, _omitFieldNames ? '' : 'tableclass',
        enumValues: TableClass.values)
    ..aOM<StreamSpecification>(
        403922627, _omitFieldNames ? '' : 'streamspecification',
        subBuilder: StreamSpecification.create)
    ..pPM<AttributeDefinition>(
        414687108, _omitFieldNames ? '' : 'attributedefinitions',
        subBuilder: AttributeDefinition.create)
    ..aE<MultiRegionConsistency>(
        446019131, _omitFieldNames ? '' : 'multiregionconsistency',
        enumValues: MultiRegionConsistency.values)
    ..aOM<OnDemandThroughput>(
        481734402, _omitFieldNames ? '' : 'ondemandthroughput',
        subBuilder: OnDemandThroughput.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableInput copyWith(void Function(UpdateTableInput) updates) =>
      super.copyWith((message) => updates(message as UpdateTableInput))
          as UpdateTableInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTableInput create() => UpdateTableInput._();
  @$core.override
  UpdateTableInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTableInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTableInput>(create);
  static UpdateTableInput? _defaultInstance;

  @$pb.TagNumber(1757580)
  ProvisionedThroughput get provisionedthroughput => $_getN(0);
  @$pb.TagNumber(1757580)
  set provisionedthroughput(ProvisionedThroughput value) =>
      $_setField(1757580, value);
  @$pb.TagNumber(1757580)
  $core.bool hasProvisionedthroughput() => $_has(0);
  @$pb.TagNumber(1757580)
  void clearProvisionedthroughput() => $_clearField(1757580);
  @$pb.TagNumber(1757580)
  ProvisionedThroughput ensureProvisionedthroughput() => $_ensure(0);

  @$pb.TagNumber(10446577)
  GlobalTableSettingsReplicationMode get globaltablesettingsreplicationmode =>
      $_getN(1);
  @$pb.TagNumber(10446577)
  set globaltablesettingsreplicationmode(
          GlobalTableSettingsReplicationMode value) =>
      $_setField(10446577, value);
  @$pb.TagNumber(10446577)
  $core.bool hasGlobaltablesettingsreplicationmode() => $_has(1);
  @$pb.TagNumber(10446577)
  void clearGlobaltablesettingsreplicationmode() => $_clearField(10446577);

  @$pb.TagNumber(31692444)
  SSESpecification get ssespecification => $_getN(2);
  @$pb.TagNumber(31692444)
  set ssespecification(SSESpecification value) => $_setField(31692444, value);
  @$pb.TagNumber(31692444)
  $core.bool hasSsespecification() => $_has(2);
  @$pb.TagNumber(31692444)
  void clearSsespecification() => $_clearField(31692444);
  @$pb.TagNumber(31692444)
  SSESpecification ensureSsespecification() => $_ensure(2);

  @$pb.TagNumber(184162880)
  BillingMode get billingmode => $_getN(3);
  @$pb.TagNumber(184162880)
  set billingmode(BillingMode value) => $_setField(184162880, value);
  @$pb.TagNumber(184162880)
  $core.bool hasBillingmode() => $_has(3);
  @$pb.TagNumber(184162880)
  void clearBillingmode() => $_clearField(184162880);

  @$pb.TagNumber(259418450)
  $core.bool get deletionprotectionenabled => $_getBF(4);
  @$pb.TagNumber(259418450)
  set deletionprotectionenabled($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(259418450)
  $core.bool hasDeletionprotectionenabled() => $_has(4);
  @$pb.TagNumber(259418450)
  void clearDeletionprotectionenabled() => $_clearField(259418450);

  @$pb.TagNumber(260731936)
  $pb.PbList<ReplicationGroupUpdate> get replicaupdates => $_getList(5);

  @$pb.TagNumber(265760923)
  $pb.PbList<GlobalSecondaryIndexUpdate> get globalsecondaryindexupdates =>
      $_getList(6);

  @$pb.TagNumber(269201458)
  $pb.PbList<GlobalTableWitnessGroupUpdate> get globaltablewitnessupdates =>
      $_getList(7);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(8);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(8, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(8);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(290598659)
  WarmThroughput get warmthroughput => $_getN(9);
  @$pb.TagNumber(290598659)
  set warmthroughput(WarmThroughput value) => $_setField(290598659, value);
  @$pb.TagNumber(290598659)
  $core.bool hasWarmthroughput() => $_has(9);
  @$pb.TagNumber(290598659)
  void clearWarmthroughput() => $_clearField(290598659);
  @$pb.TagNumber(290598659)
  WarmThroughput ensureWarmthroughput() => $_ensure(9);

  @$pb.TagNumber(342890498)
  TableClass get tableclass => $_getN(10);
  @$pb.TagNumber(342890498)
  set tableclass(TableClass value) => $_setField(342890498, value);
  @$pb.TagNumber(342890498)
  $core.bool hasTableclass() => $_has(10);
  @$pb.TagNumber(342890498)
  void clearTableclass() => $_clearField(342890498);

  @$pb.TagNumber(403922627)
  StreamSpecification get streamspecification => $_getN(11);
  @$pb.TagNumber(403922627)
  set streamspecification(StreamSpecification value) =>
      $_setField(403922627, value);
  @$pb.TagNumber(403922627)
  $core.bool hasStreamspecification() => $_has(11);
  @$pb.TagNumber(403922627)
  void clearStreamspecification() => $_clearField(403922627);
  @$pb.TagNumber(403922627)
  StreamSpecification ensureStreamspecification() => $_ensure(11);

  @$pb.TagNumber(414687108)
  $pb.PbList<AttributeDefinition> get attributedefinitions => $_getList(12);

  @$pb.TagNumber(446019131)
  MultiRegionConsistency get multiregionconsistency => $_getN(13);
  @$pb.TagNumber(446019131)
  set multiregionconsistency(MultiRegionConsistency value) =>
      $_setField(446019131, value);
  @$pb.TagNumber(446019131)
  $core.bool hasMultiregionconsistency() => $_has(13);
  @$pb.TagNumber(446019131)
  void clearMultiregionconsistency() => $_clearField(446019131);

  @$pb.TagNumber(481734402)
  OnDemandThroughput get ondemandthroughput => $_getN(14);
  @$pb.TagNumber(481734402)
  set ondemandthroughput(OnDemandThroughput value) =>
      $_setField(481734402, value);
  @$pb.TagNumber(481734402)
  $core.bool hasOndemandthroughput() => $_has(14);
  @$pb.TagNumber(481734402)
  void clearOndemandthroughput() => $_clearField(481734402);
  @$pb.TagNumber(481734402)
  OnDemandThroughput ensureOndemandthroughput() => $_ensure(14);
}

class UpdateTableOutput extends $pb.GeneratedMessage {
  factory UpdateTableOutput({
    TableDescription? tabledescription,
  }) {
    final result = create();
    if (tabledescription != null) result.tabledescription = tabledescription;
    return result;
  }

  UpdateTableOutput._();

  factory UpdateTableOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTableOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTableOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TableDescription>(19962388, _omitFieldNames ? '' : 'tabledescription',
        subBuilder: TableDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableOutput copyWith(void Function(UpdateTableOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateTableOutput))
          as UpdateTableOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTableOutput create() => UpdateTableOutput._();
  @$core.override
  UpdateTableOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTableOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTableOutput>(create);
  static UpdateTableOutput? _defaultInstance;

  @$pb.TagNumber(19962388)
  TableDescription get tabledescription => $_getN(0);
  @$pb.TagNumber(19962388)
  set tabledescription(TableDescription value) => $_setField(19962388, value);
  @$pb.TagNumber(19962388)
  $core.bool hasTabledescription() => $_has(0);
  @$pb.TagNumber(19962388)
  void clearTabledescription() => $_clearField(19962388);
  @$pb.TagNumber(19962388)
  TableDescription ensureTabledescription() => $_ensure(0);
}

class UpdateTableReplicaAutoScalingInput extends $pb.GeneratedMessage {
  factory UpdateTableReplicaAutoScalingInput({
    $core.Iterable<ReplicaAutoScalingUpdate>? replicaupdates,
    $core.Iterable<GlobalSecondaryIndexAutoScalingUpdate>?
        globalsecondaryindexupdates,
    $core.String? tablename,
    AutoScalingSettingsUpdate? provisionedwritecapacityautoscalingupdate,
  }) {
    final result = create();
    if (replicaupdates != null) result.replicaupdates.addAll(replicaupdates);
    if (globalsecondaryindexupdates != null)
      result.globalsecondaryindexupdates.addAll(globalsecondaryindexupdates);
    if (tablename != null) result.tablename = tablename;
    if (provisionedwritecapacityautoscalingupdate != null)
      result.provisionedwritecapacityautoscalingupdate =
          provisionedwritecapacityautoscalingupdate;
    return result;
  }

  UpdateTableReplicaAutoScalingInput._();

  factory UpdateTableReplicaAutoScalingInput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTableReplicaAutoScalingInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTableReplicaAutoScalingInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..pPM<ReplicaAutoScalingUpdate>(
        260731936, _omitFieldNames ? '' : 'replicaupdates',
        subBuilder: ReplicaAutoScalingUpdate.create)
    ..pPM<GlobalSecondaryIndexAutoScalingUpdate>(
        265760923, _omitFieldNames ? '' : 'globalsecondaryindexupdates',
        subBuilder: GlobalSecondaryIndexAutoScalingUpdate.create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<AutoScalingSettingsUpdate>(316618294,
        _omitFieldNames ? '' : 'provisionedwritecapacityautoscalingupdate',
        subBuilder: AutoScalingSettingsUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableReplicaAutoScalingInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableReplicaAutoScalingInput copyWith(
          void Function(UpdateTableReplicaAutoScalingInput) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateTableReplicaAutoScalingInput))
          as UpdateTableReplicaAutoScalingInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTableReplicaAutoScalingInput create() =>
      UpdateTableReplicaAutoScalingInput._();
  @$core.override
  UpdateTableReplicaAutoScalingInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTableReplicaAutoScalingInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTableReplicaAutoScalingInput>(
          create);
  static UpdateTableReplicaAutoScalingInput? _defaultInstance;

  @$pb.TagNumber(260731936)
  $pb.PbList<ReplicaAutoScalingUpdate> get replicaupdates => $_getList(0);

  @$pb.TagNumber(265760923)
  $pb.PbList<GlobalSecondaryIndexAutoScalingUpdate>
      get globalsecondaryindexupdates => $_getList(1);

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(2);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(2);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(316618294)
  AutoScalingSettingsUpdate get provisionedwritecapacityautoscalingupdate =>
      $_getN(3);
  @$pb.TagNumber(316618294)
  set provisionedwritecapacityautoscalingupdate(
          AutoScalingSettingsUpdate value) =>
      $_setField(316618294, value);
  @$pb.TagNumber(316618294)
  $core.bool hasProvisionedwritecapacityautoscalingupdate() => $_has(3);
  @$pb.TagNumber(316618294)
  void clearProvisionedwritecapacityautoscalingupdate() =>
      $_clearField(316618294);
  @$pb.TagNumber(316618294)
  AutoScalingSettingsUpdate ensureProvisionedwritecapacityautoscalingupdate() =>
      $_ensure(3);
}

class UpdateTableReplicaAutoScalingOutput extends $pb.GeneratedMessage {
  factory UpdateTableReplicaAutoScalingOutput({
    TableAutoScalingDescription? tableautoscalingdescription,
  }) {
    final result = create();
    if (tableautoscalingdescription != null)
      result.tableautoscalingdescription = tableautoscalingdescription;
    return result;
  }

  UpdateTableReplicaAutoScalingOutput._();

  factory UpdateTableReplicaAutoScalingOutput.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTableReplicaAutoScalingOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTableReplicaAutoScalingOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TableAutoScalingDescription>(
        490422806, _omitFieldNames ? '' : 'tableautoscalingdescription',
        subBuilder: TableAutoScalingDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableReplicaAutoScalingOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTableReplicaAutoScalingOutput copyWith(
          void Function(UpdateTableReplicaAutoScalingOutput) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateTableReplicaAutoScalingOutput))
          as UpdateTableReplicaAutoScalingOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTableReplicaAutoScalingOutput create() =>
      UpdateTableReplicaAutoScalingOutput._();
  @$core.override
  UpdateTableReplicaAutoScalingOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTableReplicaAutoScalingOutput getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateTableReplicaAutoScalingOutput>(create);
  static UpdateTableReplicaAutoScalingOutput? _defaultInstance;

  @$pb.TagNumber(490422806)
  TableAutoScalingDescription get tableautoscalingdescription => $_getN(0);
  @$pb.TagNumber(490422806)
  set tableautoscalingdescription(TableAutoScalingDescription value) =>
      $_setField(490422806, value);
  @$pb.TagNumber(490422806)
  $core.bool hasTableautoscalingdescription() => $_has(0);
  @$pb.TagNumber(490422806)
  void clearTableautoscalingdescription() => $_clearField(490422806);
  @$pb.TagNumber(490422806)
  TableAutoScalingDescription ensureTableautoscalingdescription() =>
      $_ensure(0);
}

class UpdateTimeToLiveInput extends $pb.GeneratedMessage {
  factory UpdateTimeToLiveInput({
    $core.String? tablename,
    TimeToLiveSpecification? timetolivespecification,
  }) {
    final result = create();
    if (tablename != null) result.tablename = tablename;
    if (timetolivespecification != null)
      result.timetolivespecification = timetolivespecification;
    return result;
  }

  UpdateTimeToLiveInput._();

  factory UpdateTimeToLiveInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTimeToLiveInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTimeToLiveInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOS(272020061, _omitFieldNames ? '' : 'tablename')
    ..aOM<TimeToLiveSpecification>(
        315882193, _omitFieldNames ? '' : 'timetolivespecification',
        subBuilder: TimeToLiveSpecification.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTimeToLiveInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTimeToLiveInput copyWith(
          void Function(UpdateTimeToLiveInput) updates) =>
      super.copyWith((message) => updates(message as UpdateTimeToLiveInput))
          as UpdateTimeToLiveInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTimeToLiveInput create() => UpdateTimeToLiveInput._();
  @$core.override
  UpdateTimeToLiveInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTimeToLiveInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTimeToLiveInput>(create);
  static UpdateTimeToLiveInput? _defaultInstance;

  @$pb.TagNumber(272020061)
  $core.String get tablename => $_getSZ(0);
  @$pb.TagNumber(272020061)
  set tablename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(272020061)
  $core.bool hasTablename() => $_has(0);
  @$pb.TagNumber(272020061)
  void clearTablename() => $_clearField(272020061);

  @$pb.TagNumber(315882193)
  TimeToLiveSpecification get timetolivespecification => $_getN(1);
  @$pb.TagNumber(315882193)
  set timetolivespecification(TimeToLiveSpecification value) =>
      $_setField(315882193, value);
  @$pb.TagNumber(315882193)
  $core.bool hasTimetolivespecification() => $_has(1);
  @$pb.TagNumber(315882193)
  void clearTimetolivespecification() => $_clearField(315882193);
  @$pb.TagNumber(315882193)
  TimeToLiveSpecification ensureTimetolivespecification() => $_ensure(1);
}

class UpdateTimeToLiveOutput extends $pb.GeneratedMessage {
  factory UpdateTimeToLiveOutput({
    TimeToLiveSpecification? timetolivespecification,
  }) {
    final result = create();
    if (timetolivespecification != null)
      result.timetolivespecification = timetolivespecification;
    return result;
  }

  UpdateTimeToLiveOutput._();

  factory UpdateTimeToLiveOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTimeToLiveOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTimeToLiveOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<TimeToLiveSpecification>(
        315882193, _omitFieldNames ? '' : 'timetolivespecification',
        subBuilder: TimeToLiveSpecification.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTimeToLiveOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTimeToLiveOutput copyWith(
          void Function(UpdateTimeToLiveOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateTimeToLiveOutput))
          as UpdateTimeToLiveOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTimeToLiveOutput create() => UpdateTimeToLiveOutput._();
  @$core.override
  UpdateTimeToLiveOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTimeToLiveOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTimeToLiveOutput>(create);
  static UpdateTimeToLiveOutput? _defaultInstance;

  @$pb.TagNumber(315882193)
  TimeToLiveSpecification get timetolivespecification => $_getN(0);
  @$pb.TagNumber(315882193)
  set timetolivespecification(TimeToLiveSpecification value) =>
      $_setField(315882193, value);
  @$pb.TagNumber(315882193)
  $core.bool hasTimetolivespecification() => $_has(0);
  @$pb.TagNumber(315882193)
  void clearTimetolivespecification() => $_clearField(315882193);
  @$pb.TagNumber(315882193)
  TimeToLiveSpecification ensureTimetolivespecification() => $_ensure(0);
}

class WarmThroughput extends $pb.GeneratedMessage {
  factory WarmThroughput({
    $fixnum.Int64? readunitspersecond,
    $fixnum.Int64? writeunitspersecond,
  }) {
    final result = create();
    if (readunitspersecond != null)
      result.readunitspersecond = readunitspersecond;
    if (writeunitspersecond != null)
      result.writeunitspersecond = writeunitspersecond;
    return result;
  }

  WarmThroughput._();

  factory WarmThroughput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WarmThroughput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WarmThroughput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aInt64(11400732, _omitFieldNames ? '' : 'readunitspersecond')
    ..aInt64(339770127, _omitFieldNames ? '' : 'writeunitspersecond')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WarmThroughput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WarmThroughput copyWith(void Function(WarmThroughput) updates) =>
      super.copyWith((message) => updates(message as WarmThroughput))
          as WarmThroughput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WarmThroughput create() => WarmThroughput._();
  @$core.override
  WarmThroughput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WarmThroughput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WarmThroughput>(create);
  static WarmThroughput? _defaultInstance;

  @$pb.TagNumber(11400732)
  $fixnum.Int64 get readunitspersecond => $_getI64(0);
  @$pb.TagNumber(11400732)
  set readunitspersecond($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(11400732)
  $core.bool hasReadunitspersecond() => $_has(0);
  @$pb.TagNumber(11400732)
  void clearReadunitspersecond() => $_clearField(11400732);

  @$pb.TagNumber(339770127)
  $fixnum.Int64 get writeunitspersecond => $_getI64(1);
  @$pb.TagNumber(339770127)
  set writeunitspersecond($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(339770127)
  $core.bool hasWriteunitspersecond() => $_has(1);
  @$pb.TagNumber(339770127)
  void clearWriteunitspersecond() => $_clearField(339770127);
}

class WriteRequest extends $pb.GeneratedMessage {
  factory WriteRequest({
    PutRequest? putrequest,
    DeleteRequest? deleterequest,
  }) {
    final result = create();
    if (putrequest != null) result.putrequest = putrequest;
    if (deleterequest != null) result.deleterequest = deleterequest;
    return result;
  }

  WriteRequest._();

  factory WriteRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WriteRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WriteRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'dynamodb'),
      createEmptyInstance: create)
    ..aOM<PutRequest>(443252836, _omitFieldNames ? '' : 'putrequest',
        subBuilder: PutRequest.create)
    ..aOM<DeleteRequest>(477809064, _omitFieldNames ? '' : 'deleterequest',
        subBuilder: DeleteRequest.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WriteRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WriteRequest copyWith(void Function(WriteRequest) updates) =>
      super.copyWith((message) => updates(message as WriteRequest))
          as WriteRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WriteRequest create() => WriteRequest._();
  @$core.override
  WriteRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WriteRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WriteRequest>(create);
  static WriteRequest? _defaultInstance;

  @$pb.TagNumber(443252836)
  PutRequest get putrequest => $_getN(0);
  @$pb.TagNumber(443252836)
  set putrequest(PutRequest value) => $_setField(443252836, value);
  @$pb.TagNumber(443252836)
  $core.bool hasPutrequest() => $_has(0);
  @$pb.TagNumber(443252836)
  void clearPutrequest() => $_clearField(443252836);
  @$pb.TagNumber(443252836)
  PutRequest ensurePutrequest() => $_ensure(0);

  @$pb.TagNumber(477809064)
  DeleteRequest get deleterequest => $_getN(1);
  @$pb.TagNumber(477809064)
  set deleterequest(DeleteRequest value) => $_setField(477809064, value);
  @$pb.TagNumber(477809064)
  $core.bool hasDeleterequest() => $_has(1);
  @$pb.TagNumber(477809064)
  void clearDeleterequest() => $_clearField(477809064);
  @$pb.TagNumber(477809064)
  DeleteRequest ensureDeleterequest() => $_ensure(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
