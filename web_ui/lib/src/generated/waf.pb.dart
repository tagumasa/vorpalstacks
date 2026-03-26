// This is a generated file - do not edit.
//
// Generated from waf.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'waf.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'waf.pbenum.dart';

class ActivatedRule extends $pb.GeneratedMessage {
  factory ActivatedRule({
    $core.int? priority,
    $core.Iterable<ExcludedRule>? excludedrules,
    WafAction? action,
    WafRuleType? type,
    $core.String? ruleid,
    WafOverrideAction? overrideaction,
  }) {
    final result = create();
    if (priority != null) result.priority = priority;
    if (excludedrules != null) result.excludedrules.addAll(excludedrules);
    if (action != null) result.action = action;
    if (type != null) result.type = type;
    if (ruleid != null) result.ruleid = ruleid;
    if (overrideaction != null) result.overrideaction = overrideaction;
    return result;
  }

  ActivatedRule._();

  factory ActivatedRule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivatedRule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivatedRule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(109944618, _omitFieldNames ? '' : 'priority')
    ..pPM<ExcludedRule>(129929071, _omitFieldNames ? '' : 'excludedrules',
        subBuilder: ExcludedRule.create)
    ..aOM<WafAction>(175614240, _omitFieldNames ? '' : 'action',
        subBuilder: WafAction.create)
    ..aE<WafRuleType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: WafRuleType.values)
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..aOM<WafOverrideAction>(515842888, _omitFieldNames ? '' : 'overrideaction',
        subBuilder: WafOverrideAction.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivatedRule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivatedRule copyWith(void Function(ActivatedRule) updates) =>
      super.copyWith((message) => updates(message as ActivatedRule))
          as ActivatedRule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivatedRule create() => ActivatedRule._();
  @$core.override
  ActivatedRule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivatedRule getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivatedRule>(create);
  static ActivatedRule? _defaultInstance;

  @$pb.TagNumber(109944618)
  $core.int get priority => $_getIZ(0);
  @$pb.TagNumber(109944618)
  set priority($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(109944618)
  $core.bool hasPriority() => $_has(0);
  @$pb.TagNumber(109944618)
  void clearPriority() => $_clearField(109944618);

  @$pb.TagNumber(129929071)
  $pb.PbList<ExcludedRule> get excludedrules => $_getList(1);

  @$pb.TagNumber(175614240)
  WafAction get action => $_getN(2);
  @$pb.TagNumber(175614240)
  set action(WafAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(2);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
  @$pb.TagNumber(175614240)
  WafAction ensureAction() => $_ensure(2);

  @$pb.TagNumber(290836590)
  WafRuleType get type => $_getN(3);
  @$pb.TagNumber(290836590)
  set type(WafRuleType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(3);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(4);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(4);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);

  @$pb.TagNumber(515842888)
  WafOverrideAction get overrideaction => $_getN(5);
  @$pb.TagNumber(515842888)
  set overrideaction(WafOverrideAction value) => $_setField(515842888, value);
  @$pb.TagNumber(515842888)
  $core.bool hasOverrideaction() => $_has(5);
  @$pb.TagNumber(515842888)
  void clearOverrideaction() => $_clearField(515842888);
  @$pb.TagNumber(515842888)
  WafOverrideAction ensureOverrideaction() => $_ensure(5);
}

class ByteMatchSet extends $pb.GeneratedMessage {
  factory ByteMatchSet({
    $core.Iterable<ByteMatchTuple>? bytematchtuples,
    $core.String? bytematchsetid,
    $core.String? name,
  }) {
    final result = create();
    if (bytematchtuples != null) result.bytematchtuples.addAll(bytematchtuples);
    if (bytematchsetid != null) result.bytematchsetid = bytematchsetid;
    if (name != null) result.name = name;
    return result;
  }

  ByteMatchSet._();

  factory ByteMatchSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ByteMatchSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ByteMatchSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<ByteMatchTuple>(3234744, _omitFieldNames ? '' : 'bytematchtuples',
        subBuilder: ByteMatchTuple.create)
    ..aOS(62033398, _omitFieldNames ? '' : 'bytematchsetid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchSet copyWith(void Function(ByteMatchSet) updates) =>
      super.copyWith((message) => updates(message as ByteMatchSet))
          as ByteMatchSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ByteMatchSet create() => ByteMatchSet._();
  @$core.override
  ByteMatchSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ByteMatchSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ByteMatchSet>(create);
  static ByteMatchSet? _defaultInstance;

  @$pb.TagNumber(3234744)
  $pb.PbList<ByteMatchTuple> get bytematchtuples => $_getList(0);

  @$pb.TagNumber(62033398)
  $core.String get bytematchsetid => $_getSZ(1);
  @$pb.TagNumber(62033398)
  set bytematchsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(62033398)
  $core.bool hasBytematchsetid() => $_has(1);
  @$pb.TagNumber(62033398)
  void clearBytematchsetid() => $_clearField(62033398);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ByteMatchSetSummary extends $pb.GeneratedMessage {
  factory ByteMatchSetSummary({
    $core.String? bytematchsetid,
    $core.String? name,
  }) {
    final result = create();
    if (bytematchsetid != null) result.bytematchsetid = bytematchsetid;
    if (name != null) result.name = name;
    return result;
  }

  ByteMatchSetSummary._();

  factory ByteMatchSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ByteMatchSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ByteMatchSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(62033398, _omitFieldNames ? '' : 'bytematchsetid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchSetSummary copyWith(void Function(ByteMatchSetSummary) updates) =>
      super.copyWith((message) => updates(message as ByteMatchSetSummary))
          as ByteMatchSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ByteMatchSetSummary create() => ByteMatchSetSummary._();
  @$core.override
  ByteMatchSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ByteMatchSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ByteMatchSetSummary>(create);
  static ByteMatchSetSummary? _defaultInstance;

  @$pb.TagNumber(62033398)
  $core.String get bytematchsetid => $_getSZ(0);
  @$pb.TagNumber(62033398)
  set bytematchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(62033398)
  $core.bool hasBytematchsetid() => $_has(0);
  @$pb.TagNumber(62033398)
  void clearBytematchsetid() => $_clearField(62033398);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ByteMatchSetUpdate extends $pb.GeneratedMessage {
  factory ByteMatchSetUpdate({
    ByteMatchTuple? bytematchtuple,
    ChangeAction? action,
  }) {
    final result = create();
    if (bytematchtuple != null) result.bytematchtuple = bytematchtuple;
    if (action != null) result.action = action;
    return result;
  }

  ByteMatchSetUpdate._();

  factory ByteMatchSetUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ByteMatchSetUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ByteMatchSetUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<ByteMatchTuple>(12788671, _omitFieldNames ? '' : 'bytematchtuple',
        subBuilder: ByteMatchTuple.create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchSetUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchSetUpdate copyWith(void Function(ByteMatchSetUpdate) updates) =>
      super.copyWith((message) => updates(message as ByteMatchSetUpdate))
          as ByteMatchSetUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ByteMatchSetUpdate create() => ByteMatchSetUpdate._();
  @$core.override
  ByteMatchSetUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ByteMatchSetUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ByteMatchSetUpdate>(create);
  static ByteMatchSetUpdate? _defaultInstance;

  @$pb.TagNumber(12788671)
  ByteMatchTuple get bytematchtuple => $_getN(0);
  @$pb.TagNumber(12788671)
  set bytematchtuple(ByteMatchTuple value) => $_setField(12788671, value);
  @$pb.TagNumber(12788671)
  $core.bool hasBytematchtuple() => $_has(0);
  @$pb.TagNumber(12788671)
  void clearBytematchtuple() => $_clearField(12788671);
  @$pb.TagNumber(12788671)
  ByteMatchTuple ensureBytematchtuple() => $_ensure(0);

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(1);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(1);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
}

class ByteMatchTuple extends $pb.GeneratedMessage {
  factory ByteMatchTuple({
    PositionalConstraint? positionalconstraint,
    FieldToMatch? fieldtomatch,
    TextTransformation? texttransformation,
    $core.List<$core.int>? targetstring,
  }) {
    final result = create();
    if (positionalconstraint != null)
      result.positionalconstraint = positionalconstraint;
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    if (texttransformation != null)
      result.texttransformation = texttransformation;
    if (targetstring != null) result.targetstring = targetstring;
    return result;
  }

  ByteMatchTuple._();

  factory ByteMatchTuple.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ByteMatchTuple.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ByteMatchTuple',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<PositionalConstraint>(
        260734525, _omitFieldNames ? '' : 'positionalconstraint',
        enumValues: PositionalConstraint.values)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..aE<TextTransformation>(
        375493608, _omitFieldNames ? '' : 'texttransformation',
        enumValues: TextTransformation.values)
    ..a<$core.List<$core.int>>(
        491118258, _omitFieldNames ? '' : 'targetstring', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchTuple clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ByteMatchTuple copyWith(void Function(ByteMatchTuple) updates) =>
      super.copyWith((message) => updates(message as ByteMatchTuple))
          as ByteMatchTuple;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ByteMatchTuple create() => ByteMatchTuple._();
  @$core.override
  ByteMatchTuple createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ByteMatchTuple getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ByteMatchTuple>(create);
  static ByteMatchTuple? _defaultInstance;

  @$pb.TagNumber(260734525)
  PositionalConstraint get positionalconstraint => $_getN(0);
  @$pb.TagNumber(260734525)
  set positionalconstraint(PositionalConstraint value) =>
      $_setField(260734525, value);
  @$pb.TagNumber(260734525)
  $core.bool hasPositionalconstraint() => $_has(0);
  @$pb.TagNumber(260734525)
  void clearPositionalconstraint() => $_clearField(260734525);

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

  @$pb.TagNumber(375493608)
  TextTransformation get texttransformation => $_getN(2);
  @$pb.TagNumber(375493608)
  set texttransformation(TextTransformation value) =>
      $_setField(375493608, value);
  @$pb.TagNumber(375493608)
  $core.bool hasTexttransformation() => $_has(2);
  @$pb.TagNumber(375493608)
  void clearTexttransformation() => $_clearField(375493608);

  @$pb.TagNumber(491118258)
  $core.List<$core.int> get targetstring => $_getN(3);
  @$pb.TagNumber(491118258)
  set targetstring($core.List<$core.int> value) => $_setBytes(3, value);
  @$pb.TagNumber(491118258)
  $core.bool hasTargetstring() => $_has(3);
  @$pb.TagNumber(491118258)
  void clearTargetstring() => $_clearField(491118258);
}

class CreateByteMatchSetRequest extends $pb.GeneratedMessage {
  factory CreateByteMatchSetRequest({
    $core.String? changetoken,
    $core.String? name,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (name != null) result.name = name;
    return result;
  }

  CreateByteMatchSetRequest._();

  factory CreateByteMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateByteMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateByteMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateByteMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateByteMatchSetRequest copyWith(
          void Function(CreateByteMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as CreateByteMatchSetRequest))
          as CreateByteMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateByteMatchSetRequest create() => CreateByteMatchSetRequest._();
  @$core.override
  CreateByteMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateByteMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateByteMatchSetRequest>(create);
  static CreateByteMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateByteMatchSetResponse extends $pb.GeneratedMessage {
  factory CreateByteMatchSetResponse({
    ByteMatchSet? bytematchset,
    $core.String? changetoken,
  }) {
    final result = create();
    if (bytematchset != null) result.bytematchset = bytematchset;
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  CreateByteMatchSetResponse._();

  factory CreateByteMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateByteMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateByteMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<ByteMatchSet>(9496489, _omitFieldNames ? '' : 'bytematchset',
        subBuilder: ByteMatchSet.create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateByteMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateByteMatchSetResponse copyWith(
          void Function(CreateByteMatchSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateByteMatchSetResponse))
          as CreateByteMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateByteMatchSetResponse create() => CreateByteMatchSetResponse._();
  @$core.override
  CreateByteMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateByteMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateByteMatchSetResponse>(create);
  static CreateByteMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(9496489)
  ByteMatchSet get bytematchset => $_getN(0);
  @$pb.TagNumber(9496489)
  set bytematchset(ByteMatchSet value) => $_setField(9496489, value);
  @$pb.TagNumber(9496489)
  $core.bool hasBytematchset() => $_has(0);
  @$pb.TagNumber(9496489)
  void clearBytematchset() => $_clearField(9496489);
  @$pb.TagNumber(9496489)
  ByteMatchSet ensureBytematchset() => $_ensure(0);

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(1);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(1);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class CreateGeoMatchSetRequest extends $pb.GeneratedMessage {
  factory CreateGeoMatchSetRequest({
    $core.String? changetoken,
    $core.String? name,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (name != null) result.name = name;
    return result;
  }

  CreateGeoMatchSetRequest._();

  factory CreateGeoMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateGeoMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateGeoMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGeoMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGeoMatchSetRequest copyWith(
          void Function(CreateGeoMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as CreateGeoMatchSetRequest))
          as CreateGeoMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateGeoMatchSetRequest create() => CreateGeoMatchSetRequest._();
  @$core.override
  CreateGeoMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateGeoMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateGeoMatchSetRequest>(create);
  static CreateGeoMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateGeoMatchSetResponse extends $pb.GeneratedMessage {
  factory CreateGeoMatchSetResponse({
    $core.String? changetoken,
    GeoMatchSet? geomatchset,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (geomatchset != null) result.geomatchset = geomatchset;
    return result;
  }

  CreateGeoMatchSetResponse._();

  factory CreateGeoMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateGeoMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateGeoMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<GeoMatchSet>(84933666, _omitFieldNames ? '' : 'geomatchset',
        subBuilder: GeoMatchSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGeoMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateGeoMatchSetResponse copyWith(
          void Function(CreateGeoMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as CreateGeoMatchSetResponse))
          as CreateGeoMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateGeoMatchSetResponse create() => CreateGeoMatchSetResponse._();
  @$core.override
  CreateGeoMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateGeoMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateGeoMatchSetResponse>(create);
  static CreateGeoMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(84933666)
  GeoMatchSet get geomatchset => $_getN(1);
  @$pb.TagNumber(84933666)
  set geomatchset(GeoMatchSet value) => $_setField(84933666, value);
  @$pb.TagNumber(84933666)
  $core.bool hasGeomatchset() => $_has(1);
  @$pb.TagNumber(84933666)
  void clearGeomatchset() => $_clearField(84933666);
  @$pb.TagNumber(84933666)
  GeoMatchSet ensureGeomatchset() => $_ensure(1);
}

class CreateIPSetRequest extends $pb.GeneratedMessage {
  factory CreateIPSetRequest({
    $core.String? changetoken,
    $core.String? name,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (name != null) result.name = name;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateIPSetResponse extends $pb.GeneratedMessage {
  factory CreateIPSetResponse({
    $core.String? changetoken,
    IPSet? ipset,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (ipset != null) result.ipset = ipset;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<IPSet>(436412565, _omitFieldNames ? '' : 'ipset',
        subBuilder: IPSet.create)
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

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

class CreateRateBasedRuleRequest extends $pb.GeneratedMessage {
  factory CreateRateBasedRuleRequest({
    RateKey? ratekey,
    $fixnum.Int64? ratelimit,
    $core.String? changetoken,
    $core.String? metricname,
    $core.String? name,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (ratekey != null) result.ratekey = ratekey;
    if (ratelimit != null) result.ratelimit = ratelimit;
    if (changetoken != null) result.changetoken = changetoken;
    if (metricname != null) result.metricname = metricname;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  CreateRateBasedRuleRequest._();

  factory CreateRateBasedRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRateBasedRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRateBasedRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<RateKey>(31421727, _omitFieldNames ? '' : 'ratekey',
        enumValues: RateKey.values)
    ..aInt64(70042947, _omitFieldNames ? '' : 'ratelimit')
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRateBasedRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRateBasedRuleRequest copyWith(
          void Function(CreateRateBasedRuleRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateRateBasedRuleRequest))
          as CreateRateBasedRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRateBasedRuleRequest create() => CreateRateBasedRuleRequest._();
  @$core.override
  CreateRateBasedRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRateBasedRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRateBasedRuleRequest>(create);
  static CreateRateBasedRuleRequest? _defaultInstance;

  @$pb.TagNumber(31421727)
  RateKey get ratekey => $_getN(0);
  @$pb.TagNumber(31421727)
  set ratekey(RateKey value) => $_setField(31421727, value);
  @$pb.TagNumber(31421727)
  $core.bool hasRatekey() => $_has(0);
  @$pb.TagNumber(31421727)
  void clearRatekey() => $_clearField(31421727);

  @$pb.TagNumber(70042947)
  $fixnum.Int64 get ratelimit => $_getI64(1);
  @$pb.TagNumber(70042947)
  set ratelimit($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(70042947)
  $core.bool hasRatelimit() => $_has(1);
  @$pb.TagNumber(70042947)
  void clearRatelimit() => $_clearField(70042947);

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(2);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(2);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(3);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(3);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(5);
}

class CreateRateBasedRuleResponse extends $pb.GeneratedMessage {
  factory CreateRateBasedRuleResponse({
    $core.String? changetoken,
    RateBasedRule? rule,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (rule != null) result.rule = rule;
    return result;
  }

  CreateRateBasedRuleResponse._();

  factory CreateRateBasedRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRateBasedRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRateBasedRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<RateBasedRule>(475696372, _omitFieldNames ? '' : 'rule',
        subBuilder: RateBasedRule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRateBasedRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRateBasedRuleResponse copyWith(
          void Function(CreateRateBasedRuleResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateRateBasedRuleResponse))
          as CreateRateBasedRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRateBasedRuleResponse create() =>
      CreateRateBasedRuleResponse._();
  @$core.override
  CreateRateBasedRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRateBasedRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRateBasedRuleResponse>(create);
  static CreateRateBasedRuleResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(475696372)
  RateBasedRule get rule => $_getN(1);
  @$pb.TagNumber(475696372)
  set rule(RateBasedRule value) => $_setField(475696372, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(1);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);
  @$pb.TagNumber(475696372)
  RateBasedRule ensureRule() => $_ensure(1);
}

class CreateRegexMatchSetRequest extends $pb.GeneratedMessage {
  factory CreateRegexMatchSetRequest({
    $core.String? changetoken,
    $core.String? name,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (name != null) result.name = name;
    return result;
  }

  CreateRegexMatchSetRequest._();

  factory CreateRegexMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRegexMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRegexMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRegexMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRegexMatchSetRequest copyWith(
          void Function(CreateRegexMatchSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateRegexMatchSetRequest))
          as CreateRegexMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRegexMatchSetRequest create() => CreateRegexMatchSetRequest._();
  @$core.override
  CreateRegexMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRegexMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRegexMatchSetRequest>(create);
  static CreateRegexMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateRegexMatchSetResponse extends $pb.GeneratedMessage {
  factory CreateRegexMatchSetResponse({
    $core.String? changetoken,
    RegexMatchSet? regexmatchset,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (regexmatchset != null) result.regexmatchset = regexmatchset;
    return result;
  }

  CreateRegexMatchSetResponse._();

  factory CreateRegexMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRegexMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRegexMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<RegexMatchSet>(534067456, _omitFieldNames ? '' : 'regexmatchset',
        subBuilder: RegexMatchSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRegexMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRegexMatchSetResponse copyWith(
          void Function(CreateRegexMatchSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateRegexMatchSetResponse))
          as CreateRegexMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRegexMatchSetResponse create() =>
      CreateRegexMatchSetResponse._();
  @$core.override
  CreateRegexMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRegexMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRegexMatchSetResponse>(create);
  static CreateRegexMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(534067456)
  RegexMatchSet get regexmatchset => $_getN(1);
  @$pb.TagNumber(534067456)
  set regexmatchset(RegexMatchSet value) => $_setField(534067456, value);
  @$pb.TagNumber(534067456)
  $core.bool hasRegexmatchset() => $_has(1);
  @$pb.TagNumber(534067456)
  void clearRegexmatchset() => $_clearField(534067456);
  @$pb.TagNumber(534067456)
  RegexMatchSet ensureRegexmatchset() => $_ensure(1);
}

class CreateRegexPatternSetRequest extends $pb.GeneratedMessage {
  factory CreateRegexPatternSetRequest({
    $core.String? changetoken,
    $core.String? name,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (name != null) result.name = name;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateRegexPatternSetResponse extends $pb.GeneratedMessage {
  factory CreateRegexPatternSetResponse({
    RegexPatternSet? regexpatternset,
    $core.String? changetoken,
  }) {
    final result = create();
    if (regexpatternset != null) result.regexpatternset = regexpatternset;
    if (changetoken != null) result.changetoken = changetoken;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<RegexPatternSet>(9374915, _omitFieldNames ? '' : 'regexpatternset',
        subBuilder: RegexPatternSet.create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(1);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(1);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class CreateRuleGroupRequest extends $pb.GeneratedMessage {
  factory CreateRuleGroupRequest({
    $core.String? changetoken,
    $core.String? metricname,
    $core.String? name,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (metricname != null) result.metricname = metricname;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(1);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(1);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(3);
}

class CreateRuleGroupResponse extends $pb.GeneratedMessage {
  factory CreateRuleGroupResponse({
    $core.String? changetoken,
    RuleGroup? rulegroup,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (rulegroup != null) result.rulegroup = rulegroup;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<RuleGroup>(267398571, _omitFieldNames ? '' : 'rulegroup',
        subBuilder: RuleGroup.create)
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

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

class CreateRuleRequest extends $pb.GeneratedMessage {
  factory CreateRuleRequest({
    $core.String? changetoken,
    $core.String? metricname,
    $core.String? name,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (metricname != null) result.metricname = metricname;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  CreateRuleRequest._();

  factory CreateRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRuleRequest copyWith(void Function(CreateRuleRequest) updates) =>
      super.copyWith((message) => updates(message as CreateRuleRequest))
          as CreateRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRuleRequest create() => CreateRuleRequest._();
  @$core.override
  CreateRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRuleRequest>(create);
  static CreateRuleRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(1);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(1);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(3);
}

class CreateRuleResponse extends $pb.GeneratedMessage {
  factory CreateRuleResponse({
    $core.String? changetoken,
    Rule? rule,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (rule != null) result.rule = rule;
    return result;
  }

  CreateRuleResponse._();

  factory CreateRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<Rule>(475696372, _omitFieldNames ? '' : 'rule',
        subBuilder: Rule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateRuleResponse copyWith(void Function(CreateRuleResponse) updates) =>
      super.copyWith((message) => updates(message as CreateRuleResponse))
          as CreateRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateRuleResponse create() => CreateRuleResponse._();
  @$core.override
  CreateRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateRuleResponse>(create);
  static CreateRuleResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(475696372)
  Rule get rule => $_getN(1);
  @$pb.TagNumber(475696372)
  set rule(Rule value) => $_setField(475696372, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(1);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);
  @$pb.TagNumber(475696372)
  Rule ensureRule() => $_ensure(1);
}

class CreateSizeConstraintSetRequest extends $pb.GeneratedMessage {
  factory CreateSizeConstraintSetRequest({
    $core.String? changetoken,
    $core.String? name,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (name != null) result.name = name;
    return result;
  }

  CreateSizeConstraintSetRequest._();

  factory CreateSizeConstraintSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateSizeConstraintSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateSizeConstraintSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSizeConstraintSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSizeConstraintSetRequest copyWith(
          void Function(CreateSizeConstraintSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateSizeConstraintSetRequest))
          as CreateSizeConstraintSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateSizeConstraintSetRequest create() =>
      CreateSizeConstraintSetRequest._();
  @$core.override
  CreateSizeConstraintSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateSizeConstraintSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateSizeConstraintSetRequest>(create);
  static CreateSizeConstraintSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateSizeConstraintSetResponse extends $pb.GeneratedMessage {
  factory CreateSizeConstraintSetResponse({
    $core.String? changetoken,
    SizeConstraintSet? sizeconstraintset,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (sizeconstraintset != null) result.sizeconstraintset = sizeconstraintset;
    return result;
  }

  CreateSizeConstraintSetResponse._();

  factory CreateSizeConstraintSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateSizeConstraintSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateSizeConstraintSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<SizeConstraintSet>(
        166385770, _omitFieldNames ? '' : 'sizeconstraintset',
        subBuilder: SizeConstraintSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSizeConstraintSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSizeConstraintSetResponse copyWith(
          void Function(CreateSizeConstraintSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateSizeConstraintSetResponse))
          as CreateSizeConstraintSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateSizeConstraintSetResponse create() =>
      CreateSizeConstraintSetResponse._();
  @$core.override
  CreateSizeConstraintSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateSizeConstraintSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateSizeConstraintSetResponse>(
          create);
  static CreateSizeConstraintSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(166385770)
  SizeConstraintSet get sizeconstraintset => $_getN(1);
  @$pb.TagNumber(166385770)
  set sizeconstraintset(SizeConstraintSet value) =>
      $_setField(166385770, value);
  @$pb.TagNumber(166385770)
  $core.bool hasSizeconstraintset() => $_has(1);
  @$pb.TagNumber(166385770)
  void clearSizeconstraintset() => $_clearField(166385770);
  @$pb.TagNumber(166385770)
  SizeConstraintSet ensureSizeconstraintset() => $_ensure(1);
}

class CreateSqlInjectionMatchSetRequest extends $pb.GeneratedMessage {
  factory CreateSqlInjectionMatchSetRequest({
    $core.String? changetoken,
    $core.String? name,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (name != null) result.name = name;
    return result;
  }

  CreateSqlInjectionMatchSetRequest._();

  factory CreateSqlInjectionMatchSetRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateSqlInjectionMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateSqlInjectionMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSqlInjectionMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSqlInjectionMatchSetRequest copyWith(
          void Function(CreateSqlInjectionMatchSetRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreateSqlInjectionMatchSetRequest))
          as CreateSqlInjectionMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateSqlInjectionMatchSetRequest create() =>
      CreateSqlInjectionMatchSetRequest._();
  @$core.override
  CreateSqlInjectionMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateSqlInjectionMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateSqlInjectionMatchSetRequest>(
          create);
  static CreateSqlInjectionMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateSqlInjectionMatchSetResponse extends $pb.GeneratedMessage {
  factory CreateSqlInjectionMatchSetResponse({
    SqlInjectionMatchSet? sqlinjectionmatchset,
    $core.String? changetoken,
  }) {
    final result = create();
    if (sqlinjectionmatchset != null)
      result.sqlinjectionmatchset = sqlinjectionmatchset;
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  CreateSqlInjectionMatchSetResponse._();

  factory CreateSqlInjectionMatchSetResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateSqlInjectionMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateSqlInjectionMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<SqlInjectionMatchSet>(
        60356588, _omitFieldNames ? '' : 'sqlinjectionmatchset',
        subBuilder: SqlInjectionMatchSet.create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSqlInjectionMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateSqlInjectionMatchSetResponse copyWith(
          void Function(CreateSqlInjectionMatchSetResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreateSqlInjectionMatchSetResponse))
          as CreateSqlInjectionMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateSqlInjectionMatchSetResponse create() =>
      CreateSqlInjectionMatchSetResponse._();
  @$core.override
  CreateSqlInjectionMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateSqlInjectionMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateSqlInjectionMatchSetResponse>(
          create);
  static CreateSqlInjectionMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(60356588)
  SqlInjectionMatchSet get sqlinjectionmatchset => $_getN(0);
  @$pb.TagNumber(60356588)
  set sqlinjectionmatchset(SqlInjectionMatchSet value) =>
      $_setField(60356588, value);
  @$pb.TagNumber(60356588)
  $core.bool hasSqlinjectionmatchset() => $_has(0);
  @$pb.TagNumber(60356588)
  void clearSqlinjectionmatchset() => $_clearField(60356588);
  @$pb.TagNumber(60356588)
  SqlInjectionMatchSet ensureSqlinjectionmatchset() => $_ensure(0);

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(1);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(1);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class CreateWebACLMigrationStackRequest extends $pb.GeneratedMessage {
  factory CreateWebACLMigrationStackRequest({
    $core.String? webaclid,
    $core.bool? ignoreunsupportedtype,
    $core.String? s3bucketname,
  }) {
    final result = create();
    if (webaclid != null) result.webaclid = webaclid;
    if (ignoreunsupportedtype != null)
      result.ignoreunsupportedtype = ignoreunsupportedtype;
    if (s3bucketname != null) result.s3bucketname = s3bucketname;
    return result;
  }

  CreateWebACLMigrationStackRequest._();

  factory CreateWebACLMigrationStackRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateWebACLMigrationStackRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateWebACLMigrationStackRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(161274579, _omitFieldNames ? '' : 'webaclid')
    ..aOB(171400459, _omitFieldNames ? '' : 'ignoreunsupportedtype')
    ..aOS(320495427, _omitFieldNames ? '' : 's3bucketname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWebACLMigrationStackRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWebACLMigrationStackRequest copyWith(
          void Function(CreateWebACLMigrationStackRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreateWebACLMigrationStackRequest))
          as CreateWebACLMigrationStackRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWebACLMigrationStackRequest create() =>
      CreateWebACLMigrationStackRequest._();
  @$core.override
  CreateWebACLMigrationStackRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateWebACLMigrationStackRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateWebACLMigrationStackRequest>(
          create);
  static CreateWebACLMigrationStackRequest? _defaultInstance;

  @$pb.TagNumber(161274579)
  $core.String get webaclid => $_getSZ(0);
  @$pb.TagNumber(161274579)
  set webaclid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(161274579)
  $core.bool hasWebaclid() => $_has(0);
  @$pb.TagNumber(161274579)
  void clearWebaclid() => $_clearField(161274579);

  @$pb.TagNumber(171400459)
  $core.bool get ignoreunsupportedtype => $_getBF(1);
  @$pb.TagNumber(171400459)
  set ignoreunsupportedtype($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(171400459)
  $core.bool hasIgnoreunsupportedtype() => $_has(1);
  @$pb.TagNumber(171400459)
  void clearIgnoreunsupportedtype() => $_clearField(171400459);

  @$pb.TagNumber(320495427)
  $core.String get s3bucketname => $_getSZ(2);
  @$pb.TagNumber(320495427)
  set s3bucketname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(320495427)
  $core.bool hasS3bucketname() => $_has(2);
  @$pb.TagNumber(320495427)
  void clearS3bucketname() => $_clearField(320495427);
}

class CreateWebACLMigrationStackResponse extends $pb.GeneratedMessage {
  factory CreateWebACLMigrationStackResponse({
    $core.String? s3objecturl,
  }) {
    final result = create();
    if (s3objecturl != null) result.s3objecturl = s3objecturl;
    return result;
  }

  CreateWebACLMigrationStackResponse._();

  factory CreateWebACLMigrationStackResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateWebACLMigrationStackResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateWebACLMigrationStackResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(489662762, _omitFieldNames ? '' : 's3objecturl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWebACLMigrationStackResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWebACLMigrationStackResponse copyWith(
          void Function(CreateWebACLMigrationStackResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreateWebACLMigrationStackResponse))
          as CreateWebACLMigrationStackResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWebACLMigrationStackResponse create() =>
      CreateWebACLMigrationStackResponse._();
  @$core.override
  CreateWebACLMigrationStackResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateWebACLMigrationStackResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateWebACLMigrationStackResponse>(
          create);
  static CreateWebACLMigrationStackResponse? _defaultInstance;

  @$pb.TagNumber(489662762)
  $core.String get s3objecturl => $_getSZ(0);
  @$pb.TagNumber(489662762)
  set s3objecturl($core.String value) => $_setString(0, value);
  @$pb.TagNumber(489662762)
  $core.bool hasS3objecturl() => $_has(0);
  @$pb.TagNumber(489662762)
  void clearS3objecturl() => $_clearField(489662762);
}

class CreateWebACLRequest extends $pb.GeneratedMessage {
  factory CreateWebACLRequest({
    $core.String? changetoken,
    $core.String? metricname,
    $core.String? name,
    WafAction? defaultaction,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (metricname != null) result.metricname = metricname;
    if (name != null) result.name = name;
    if (defaultaction != null) result.defaultaction = defaultaction;
    if (tags != null) result.tags.addAll(tags);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<WafAction>(322663861, _omitFieldNames ? '' : 'defaultaction',
        subBuilder: WafAction.create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(1);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(1);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(322663861)
  WafAction get defaultaction => $_getN(3);
  @$pb.TagNumber(322663861)
  set defaultaction(WafAction value) => $_setField(322663861, value);
  @$pb.TagNumber(322663861)
  $core.bool hasDefaultaction() => $_has(3);
  @$pb.TagNumber(322663861)
  void clearDefaultaction() => $_clearField(322663861);
  @$pb.TagNumber(322663861)
  WafAction ensureDefaultaction() => $_ensure(3);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(4);
}

class CreateWebACLResponse extends $pb.GeneratedMessage {
  factory CreateWebACLResponse({
    $core.String? changetoken,
    WebACL? webacl,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (webacl != null) result.webacl = webacl;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<WebACL>(343373504, _omitFieldNames ? '' : 'webacl',
        subBuilder: WebACL.create)
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(343373504)
  WebACL get webacl => $_getN(1);
  @$pb.TagNumber(343373504)
  set webacl(WebACL value) => $_setField(343373504, value);
  @$pb.TagNumber(343373504)
  $core.bool hasWebacl() => $_has(1);
  @$pb.TagNumber(343373504)
  void clearWebacl() => $_clearField(343373504);
  @$pb.TagNumber(343373504)
  WebACL ensureWebacl() => $_ensure(1);
}

class CreateXssMatchSetRequest extends $pb.GeneratedMessage {
  factory CreateXssMatchSetRequest({
    $core.String? changetoken,
    $core.String? name,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (name != null) result.name = name;
    return result;
  }

  CreateXssMatchSetRequest._();

  factory CreateXssMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateXssMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateXssMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateXssMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateXssMatchSetRequest copyWith(
          void Function(CreateXssMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as CreateXssMatchSetRequest))
          as CreateXssMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateXssMatchSetRequest create() => CreateXssMatchSetRequest._();
  @$core.override
  CreateXssMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateXssMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateXssMatchSetRequest>(create);
  static CreateXssMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateXssMatchSetResponse extends $pb.GeneratedMessage {
  factory CreateXssMatchSetResponse({
    $core.String? changetoken,
    XssMatchSet? xssmatchset,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (xssmatchset != null) result.xssmatchset = xssmatchset;
    return result;
  }

  CreateXssMatchSetResponse._();

  factory CreateXssMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateXssMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateXssMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOM<XssMatchSet>(169555847, _omitFieldNames ? '' : 'xssmatchset',
        subBuilder: XssMatchSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateXssMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateXssMatchSetResponse copyWith(
          void Function(CreateXssMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as CreateXssMatchSetResponse))
          as CreateXssMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateXssMatchSetResponse create() => CreateXssMatchSetResponse._();
  @$core.override
  CreateXssMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateXssMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateXssMatchSetResponse>(create);
  static CreateXssMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(169555847)
  XssMatchSet get xssmatchset => $_getN(1);
  @$pb.TagNumber(169555847)
  set xssmatchset(XssMatchSet value) => $_setField(169555847, value);
  @$pb.TagNumber(169555847)
  $core.bool hasXssmatchset() => $_has(1);
  @$pb.TagNumber(169555847)
  void clearXssmatchset() => $_clearField(169555847);
  @$pb.TagNumber(169555847)
  XssMatchSet ensureXssmatchset() => $_ensure(1);
}

class DeleteByteMatchSetRequest extends $pb.GeneratedMessage {
  factory DeleteByteMatchSetRequest({
    $core.String? bytematchsetid,
    $core.String? changetoken,
  }) {
    final result = create();
    if (bytematchsetid != null) result.bytematchsetid = bytematchsetid;
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteByteMatchSetRequest._();

  factory DeleteByteMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteByteMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteByteMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(62033398, _omitFieldNames ? '' : 'bytematchsetid')
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteByteMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteByteMatchSetRequest copyWith(
          void Function(DeleteByteMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteByteMatchSetRequest))
          as DeleteByteMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteByteMatchSetRequest create() => DeleteByteMatchSetRequest._();
  @$core.override
  DeleteByteMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteByteMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteByteMatchSetRequest>(create);
  static DeleteByteMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(62033398)
  $core.String get bytematchsetid => $_getSZ(0);
  @$pb.TagNumber(62033398)
  set bytematchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(62033398)
  $core.bool hasBytematchsetid() => $_has(0);
  @$pb.TagNumber(62033398)
  void clearBytematchsetid() => $_clearField(62033398);

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(1);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(1);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteByteMatchSetResponse extends $pb.GeneratedMessage {
  factory DeleteByteMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteByteMatchSetResponse._();

  factory DeleteByteMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteByteMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteByteMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteByteMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteByteMatchSetResponse copyWith(
          void Function(DeleteByteMatchSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteByteMatchSetResponse))
          as DeleteByteMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteByteMatchSetResponse create() => DeleteByteMatchSetResponse._();
  @$core.override
  DeleteByteMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteByteMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteByteMatchSetResponse>(create);
  static DeleteByteMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteGeoMatchSetRequest extends $pb.GeneratedMessage {
  factory DeleteGeoMatchSetRequest({
    $core.String? changetoken,
    $core.String? geomatchsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (geomatchsetid != null) result.geomatchsetid = geomatchsetid;
    return result;
  }

  DeleteGeoMatchSetRequest._();

  factory DeleteGeoMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteGeoMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteGeoMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(514346837, _omitFieldNames ? '' : 'geomatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGeoMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGeoMatchSetRequest copyWith(
          void Function(DeleteGeoMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteGeoMatchSetRequest))
          as DeleteGeoMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteGeoMatchSetRequest create() => DeleteGeoMatchSetRequest._();
  @$core.override
  DeleteGeoMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteGeoMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteGeoMatchSetRequest>(create);
  static DeleteGeoMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(514346837)
  $core.String get geomatchsetid => $_getSZ(1);
  @$pb.TagNumber(514346837)
  set geomatchsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(514346837)
  $core.bool hasGeomatchsetid() => $_has(1);
  @$pb.TagNumber(514346837)
  void clearGeomatchsetid() => $_clearField(514346837);
}

class DeleteGeoMatchSetResponse extends $pb.GeneratedMessage {
  factory DeleteGeoMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteGeoMatchSetResponse._();

  factory DeleteGeoMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteGeoMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteGeoMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGeoMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteGeoMatchSetResponse copyWith(
          void Function(DeleteGeoMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteGeoMatchSetResponse))
          as DeleteGeoMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteGeoMatchSetResponse create() => DeleteGeoMatchSetResponse._();
  @$core.override
  DeleteGeoMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteGeoMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteGeoMatchSetResponse>(create);
  static DeleteGeoMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteIPSetRequest extends $pb.GeneratedMessage {
  factory DeleteIPSetRequest({
    $core.String? changetoken,
    $core.String? ipsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (ipsetid != null) result.ipsetid = ipsetid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(346763066, _omitFieldNames ? '' : 'ipsetid')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(346763066)
  $core.String get ipsetid => $_getSZ(1);
  @$pb.TagNumber(346763066)
  set ipsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346763066)
  $core.bool hasIpsetid() => $_has(1);
  @$pb.TagNumber(346763066)
  void clearIpsetid() => $_clearField(346763066);
}

class DeleteIPSetResponse extends $pb.GeneratedMessage {
  factory DeleteIPSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteIPSetResponse._();

  factory DeleteIPSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteIPSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteIPSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteLoggingConfigurationRequest extends $pb.GeneratedMessage {
  factory DeleteLoggingConfigurationRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class DeleteRateBasedRuleRequest extends $pb.GeneratedMessage {
  factory DeleteRateBasedRuleRequest({
    $core.String? changetoken,
    $core.String? ruleid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (ruleid != null) result.ruleid = ruleid;
    return result;
  }

  DeleteRateBasedRuleRequest._();

  factory DeleteRateBasedRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRateBasedRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRateBasedRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRateBasedRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRateBasedRuleRequest copyWith(
          void Function(DeleteRateBasedRuleRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteRateBasedRuleRequest))
          as DeleteRateBasedRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRateBasedRuleRequest create() => DeleteRateBasedRuleRequest._();
  @$core.override
  DeleteRateBasedRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRateBasedRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRateBasedRuleRequest>(create);
  static DeleteRateBasedRuleRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(1);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(1);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class DeleteRateBasedRuleResponse extends $pb.GeneratedMessage {
  factory DeleteRateBasedRuleResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteRateBasedRuleResponse._();

  factory DeleteRateBasedRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRateBasedRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRateBasedRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRateBasedRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRateBasedRuleResponse copyWith(
          void Function(DeleteRateBasedRuleResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteRateBasedRuleResponse))
          as DeleteRateBasedRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRateBasedRuleResponse create() =>
      DeleteRateBasedRuleResponse._();
  @$core.override
  DeleteRateBasedRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRateBasedRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRateBasedRuleResponse>(create);
  static DeleteRateBasedRuleResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteRegexMatchSetRequest extends $pb.GeneratedMessage {
  factory DeleteRegexMatchSetRequest({
    $core.String? changetoken,
    $core.String? regexmatchsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (regexmatchsetid != null) result.regexmatchsetid = regexmatchsetid;
    return result;
  }

  DeleteRegexMatchSetRequest._();

  factory DeleteRegexMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRegexMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRegexMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(82287635, _omitFieldNames ? '' : 'regexmatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRegexMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRegexMatchSetRequest copyWith(
          void Function(DeleteRegexMatchSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteRegexMatchSetRequest))
          as DeleteRegexMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRegexMatchSetRequest create() => DeleteRegexMatchSetRequest._();
  @$core.override
  DeleteRegexMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRegexMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRegexMatchSetRequest>(create);
  static DeleteRegexMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(82287635)
  $core.String get regexmatchsetid => $_getSZ(1);
  @$pb.TagNumber(82287635)
  set regexmatchsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82287635)
  $core.bool hasRegexmatchsetid() => $_has(1);
  @$pb.TagNumber(82287635)
  void clearRegexmatchsetid() => $_clearField(82287635);
}

class DeleteRegexMatchSetResponse extends $pb.GeneratedMessage {
  factory DeleteRegexMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteRegexMatchSetResponse._();

  factory DeleteRegexMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRegexMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRegexMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRegexMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRegexMatchSetResponse copyWith(
          void Function(DeleteRegexMatchSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteRegexMatchSetResponse))
          as DeleteRegexMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRegexMatchSetResponse create() =>
      DeleteRegexMatchSetResponse._();
  @$core.override
  DeleteRegexMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRegexMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRegexMatchSetResponse>(create);
  static DeleteRegexMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteRegexPatternSetRequest extends $pb.GeneratedMessage {
  factory DeleteRegexPatternSetRequest({
    $core.String? changetoken,
    $core.String? regexpatternsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (regexpatternsetid != null) result.regexpatternsetid = regexpatternsetid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(96220912, _omitFieldNames ? '' : 'regexpatternsetid')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(96220912)
  $core.String get regexpatternsetid => $_getSZ(1);
  @$pb.TagNumber(96220912)
  set regexpatternsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(96220912)
  $core.bool hasRegexpatternsetid() => $_has(1);
  @$pb.TagNumber(96220912)
  void clearRegexpatternsetid() => $_clearField(96220912);
}

class DeleteRegexPatternSetResponse extends $pb.GeneratedMessage {
  factory DeleteRegexPatternSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteRegexPatternSetResponse._();

  factory DeleteRegexPatternSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRegexPatternSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRegexPatternSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteRuleGroupRequest extends $pb.GeneratedMessage {
  factory DeleteRuleGroupRequest({
    $core.String? changetoken,
    $core.String? rulegroupid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (rulegroupid != null) result.rulegroupid = rulegroupid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(287757320, _omitFieldNames ? '' : 'rulegroupid')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(287757320)
  $core.String get rulegroupid => $_getSZ(1);
  @$pb.TagNumber(287757320)
  set rulegroupid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(287757320)
  $core.bool hasRulegroupid() => $_has(1);
  @$pb.TagNumber(287757320)
  void clearRulegroupid() => $_clearField(287757320);
}

class DeleteRuleGroupResponse extends $pb.GeneratedMessage {
  factory DeleteRuleGroupResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteRuleGroupResponse._();

  factory DeleteRuleGroupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRuleGroupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRuleGroupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteRuleRequest extends $pb.GeneratedMessage {
  factory DeleteRuleRequest({
    $core.String? changetoken,
    $core.String? ruleid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (ruleid != null) result.ruleid = ruleid;
    return result;
  }

  DeleteRuleRequest._();

  factory DeleteRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRuleRequest copyWith(void Function(DeleteRuleRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteRuleRequest))
          as DeleteRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRuleRequest create() => DeleteRuleRequest._();
  @$core.override
  DeleteRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRuleRequest>(create);
  static DeleteRuleRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(1);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(1);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class DeleteRuleResponse extends $pb.GeneratedMessage {
  factory DeleteRuleResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteRuleResponse._();

  factory DeleteRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteRuleResponse copyWith(void Function(DeleteRuleResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteRuleResponse))
          as DeleteRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteRuleResponse create() => DeleteRuleResponse._();
  @$core.override
  DeleteRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteRuleResponse>(create);
  static DeleteRuleResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteSizeConstraintSetRequest extends $pb.GeneratedMessage {
  factory DeleteSizeConstraintSetRequest({
    $core.String? changetoken,
    $core.String? sizeconstraintsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (sizeconstraintsetid != null)
      result.sizeconstraintsetid = sizeconstraintsetid;
    return result;
  }

  DeleteSizeConstraintSetRequest._();

  factory DeleteSizeConstraintSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteSizeConstraintSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteSizeConstraintSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(277959757, _omitFieldNames ? '' : 'sizeconstraintsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSizeConstraintSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSizeConstraintSetRequest copyWith(
          void Function(DeleteSizeConstraintSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteSizeConstraintSetRequest))
          as DeleteSizeConstraintSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteSizeConstraintSetRequest create() =>
      DeleteSizeConstraintSetRequest._();
  @$core.override
  DeleteSizeConstraintSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteSizeConstraintSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteSizeConstraintSetRequest>(create);
  static DeleteSizeConstraintSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(277959757)
  $core.String get sizeconstraintsetid => $_getSZ(1);
  @$pb.TagNumber(277959757)
  set sizeconstraintsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(277959757)
  $core.bool hasSizeconstraintsetid() => $_has(1);
  @$pb.TagNumber(277959757)
  void clearSizeconstraintsetid() => $_clearField(277959757);
}

class DeleteSizeConstraintSetResponse extends $pb.GeneratedMessage {
  factory DeleteSizeConstraintSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteSizeConstraintSetResponse._();

  factory DeleteSizeConstraintSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteSizeConstraintSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteSizeConstraintSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSizeConstraintSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSizeConstraintSetResponse copyWith(
          void Function(DeleteSizeConstraintSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteSizeConstraintSetResponse))
          as DeleteSizeConstraintSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteSizeConstraintSetResponse create() =>
      DeleteSizeConstraintSetResponse._();
  @$core.override
  DeleteSizeConstraintSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteSizeConstraintSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteSizeConstraintSetResponse>(
          create);
  static DeleteSizeConstraintSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteSqlInjectionMatchSetRequest extends $pb.GeneratedMessage {
  factory DeleteSqlInjectionMatchSetRequest({
    $core.String? changetoken,
    $core.String? sqlinjectionmatchsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (sqlinjectionmatchsetid != null)
      result.sqlinjectionmatchsetid = sqlinjectionmatchsetid;
    return result;
  }

  DeleteSqlInjectionMatchSetRequest._();

  factory DeleteSqlInjectionMatchSetRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteSqlInjectionMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteSqlInjectionMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(386225735, _omitFieldNames ? '' : 'sqlinjectionmatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSqlInjectionMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSqlInjectionMatchSetRequest copyWith(
          void Function(DeleteSqlInjectionMatchSetRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteSqlInjectionMatchSetRequest))
          as DeleteSqlInjectionMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteSqlInjectionMatchSetRequest create() =>
      DeleteSqlInjectionMatchSetRequest._();
  @$core.override
  DeleteSqlInjectionMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteSqlInjectionMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteSqlInjectionMatchSetRequest>(
          create);
  static DeleteSqlInjectionMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(386225735)
  $core.String get sqlinjectionmatchsetid => $_getSZ(1);
  @$pb.TagNumber(386225735)
  set sqlinjectionmatchsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(386225735)
  $core.bool hasSqlinjectionmatchsetid() => $_has(1);
  @$pb.TagNumber(386225735)
  void clearSqlinjectionmatchsetid() => $_clearField(386225735);
}

class DeleteSqlInjectionMatchSetResponse extends $pb.GeneratedMessage {
  factory DeleteSqlInjectionMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteSqlInjectionMatchSetResponse._();

  factory DeleteSqlInjectionMatchSetResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteSqlInjectionMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteSqlInjectionMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSqlInjectionMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteSqlInjectionMatchSetResponse copyWith(
          void Function(DeleteSqlInjectionMatchSetResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteSqlInjectionMatchSetResponse))
          as DeleteSqlInjectionMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteSqlInjectionMatchSetResponse create() =>
      DeleteSqlInjectionMatchSetResponse._();
  @$core.override
  DeleteSqlInjectionMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteSqlInjectionMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteSqlInjectionMatchSetResponse>(
          create);
  static DeleteSqlInjectionMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteWebACLRequest extends $pb.GeneratedMessage {
  factory DeleteWebACLRequest({
    $core.String? changetoken,
    $core.String? webaclid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (webaclid != null) result.webaclid = webaclid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(161274579, _omitFieldNames ? '' : 'webaclid')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(161274579)
  $core.String get webaclid => $_getSZ(1);
  @$pb.TagNumber(161274579)
  set webaclid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(161274579)
  $core.bool hasWebaclid() => $_has(1);
  @$pb.TagNumber(161274579)
  void clearWebaclid() => $_clearField(161274579);
}

class DeleteWebACLResponse extends $pb.GeneratedMessage {
  factory DeleteWebACLResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteWebACLResponse._();

  factory DeleteWebACLResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteWebACLResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteWebACLResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteXssMatchSetRequest extends $pb.GeneratedMessage {
  factory DeleteXssMatchSetRequest({
    $core.String? xssmatchsetid,
    $core.String? changetoken,
  }) {
    final result = create();
    if (xssmatchsetid != null) result.xssmatchsetid = xssmatchsetid;
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteXssMatchSetRequest._();

  factory DeleteXssMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteXssMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteXssMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(56643900, _omitFieldNames ? '' : 'xssmatchsetid')
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteXssMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteXssMatchSetRequest copyWith(
          void Function(DeleteXssMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteXssMatchSetRequest))
          as DeleteXssMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteXssMatchSetRequest create() => DeleteXssMatchSetRequest._();
  @$core.override
  DeleteXssMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteXssMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteXssMatchSetRequest>(create);
  static DeleteXssMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(56643900)
  $core.String get xssmatchsetid => $_getSZ(0);
  @$pb.TagNumber(56643900)
  set xssmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56643900)
  $core.bool hasXssmatchsetid() => $_has(0);
  @$pb.TagNumber(56643900)
  void clearXssmatchsetid() => $_clearField(56643900);

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(1);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(1);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class DeleteXssMatchSetResponse extends $pb.GeneratedMessage {
  factory DeleteXssMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  DeleteXssMatchSetResponse._();

  factory DeleteXssMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteXssMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteXssMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteXssMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteXssMatchSetResponse copyWith(
          void Function(DeleteXssMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteXssMatchSetResponse))
          as DeleteXssMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteXssMatchSetResponse create() => DeleteXssMatchSetResponse._();
  @$core.override
  DeleteXssMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteXssMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteXssMatchSetResponse>(create);
  static DeleteXssMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class ExcludedRule extends $pb.GeneratedMessage {
  factory ExcludedRule({
    $core.String? ruleid,
  }) {
    final result = create();
    if (ruleid != null) result.ruleid = ruleid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
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

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(0);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(0);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class FieldToMatch extends $pb.GeneratedMessage {
  factory FieldToMatch({
    MatchFieldType? type,
    $core.String? data,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (data != null) result.data = data;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<MatchFieldType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: MatchFieldType.values)
    ..aOS(525498822, _omitFieldNames ? '' : 'data')
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

  @$pb.TagNumber(290836590)
  MatchFieldType get type => $_getN(0);
  @$pb.TagNumber(290836590)
  set type(MatchFieldType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(525498822)
  $core.String get data => $_getSZ(1);
  @$pb.TagNumber(525498822)
  set data($core.String value) => $_setString(1, value);
  @$pb.TagNumber(525498822)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(525498822)
  void clearData() => $_clearField(525498822);
}

class GeoMatchConstraint extends $pb.GeneratedMessage {
  factory GeoMatchConstraint({
    GeoMatchConstraintValue? value,
    GeoMatchConstraintType? type,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (type != null) result.type = type;
    return result;
  }

  GeoMatchConstraint._();

  factory GeoMatchConstraint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GeoMatchConstraint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoMatchConstraint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<GeoMatchConstraintValue>(289929579, _omitFieldNames ? '' : 'value',
        enumValues: GeoMatchConstraintValue.values)
    ..aE<GeoMatchConstraintType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: GeoMatchConstraintType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchConstraint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchConstraint copyWith(void Function(GeoMatchConstraint) updates) =>
      super.copyWith((message) => updates(message as GeoMatchConstraint))
          as GeoMatchConstraint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoMatchConstraint create() => GeoMatchConstraint._();
  @$core.override
  GeoMatchConstraint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GeoMatchConstraint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GeoMatchConstraint>(create);
  static GeoMatchConstraint? _defaultInstance;

  @$pb.TagNumber(289929579)
  GeoMatchConstraintValue get value => $_getN(0);
  @$pb.TagNumber(289929579)
  set value(GeoMatchConstraintValue value) => $_setField(289929579, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(290836590)
  GeoMatchConstraintType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(GeoMatchConstraintType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class GeoMatchSet extends $pb.GeneratedMessage {
  factory GeoMatchSet({
    $core.Iterable<GeoMatchConstraint>? geomatchconstraints,
    $core.String? name,
    $core.String? geomatchsetid,
  }) {
    final result = create();
    if (geomatchconstraints != null)
      result.geomatchconstraints.addAll(geomatchconstraints);
    if (name != null) result.name = name;
    if (geomatchsetid != null) result.geomatchsetid = geomatchsetid;
    return result;
  }

  GeoMatchSet._();

  factory GeoMatchSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GeoMatchSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoMatchSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<GeoMatchConstraint>(
        161026282, _omitFieldNames ? '' : 'geomatchconstraints',
        subBuilder: GeoMatchConstraint.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(514346837, _omitFieldNames ? '' : 'geomatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchSet copyWith(void Function(GeoMatchSet) updates) =>
      super.copyWith((message) => updates(message as GeoMatchSet))
          as GeoMatchSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoMatchSet create() => GeoMatchSet._();
  @$core.override
  GeoMatchSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GeoMatchSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GeoMatchSet>(create);
  static GeoMatchSet? _defaultInstance;

  @$pb.TagNumber(161026282)
  $pb.PbList<GeoMatchConstraint> get geomatchconstraints => $_getList(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(514346837)
  $core.String get geomatchsetid => $_getSZ(2);
  @$pb.TagNumber(514346837)
  set geomatchsetid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(514346837)
  $core.bool hasGeomatchsetid() => $_has(2);
  @$pb.TagNumber(514346837)
  void clearGeomatchsetid() => $_clearField(514346837);
}

class GeoMatchSetSummary extends $pb.GeneratedMessage {
  factory GeoMatchSetSummary({
    $core.String? name,
    $core.String? geomatchsetid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (geomatchsetid != null) result.geomatchsetid = geomatchsetid;
    return result;
  }

  GeoMatchSetSummary._();

  factory GeoMatchSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GeoMatchSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoMatchSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(514346837, _omitFieldNames ? '' : 'geomatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchSetSummary copyWith(void Function(GeoMatchSetSummary) updates) =>
      super.copyWith((message) => updates(message as GeoMatchSetSummary))
          as GeoMatchSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoMatchSetSummary create() => GeoMatchSetSummary._();
  @$core.override
  GeoMatchSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GeoMatchSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GeoMatchSetSummary>(create);
  static GeoMatchSetSummary? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(514346837)
  $core.String get geomatchsetid => $_getSZ(1);
  @$pb.TagNumber(514346837)
  set geomatchsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(514346837)
  $core.bool hasGeomatchsetid() => $_has(1);
  @$pb.TagNumber(514346837)
  void clearGeomatchsetid() => $_clearField(514346837);
}

class GeoMatchSetUpdate extends $pb.GeneratedMessage {
  factory GeoMatchSetUpdate({
    GeoMatchConstraint? geomatchconstraint,
    ChangeAction? action,
  }) {
    final result = create();
    if (geomatchconstraint != null)
      result.geomatchconstraint = geomatchconstraint;
    if (action != null) result.action = action;
    return result;
  }

  GeoMatchSetUpdate._();

  factory GeoMatchSetUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GeoMatchSetUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoMatchSetUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<GeoMatchConstraint>(
        52896097, _omitFieldNames ? '' : 'geomatchconstraint',
        subBuilder: GeoMatchConstraint.create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchSetUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoMatchSetUpdate copyWith(void Function(GeoMatchSetUpdate) updates) =>
      super.copyWith((message) => updates(message as GeoMatchSetUpdate))
          as GeoMatchSetUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoMatchSetUpdate create() => GeoMatchSetUpdate._();
  @$core.override
  GeoMatchSetUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GeoMatchSetUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GeoMatchSetUpdate>(create);
  static GeoMatchSetUpdate? _defaultInstance;

  @$pb.TagNumber(52896097)
  GeoMatchConstraint get geomatchconstraint => $_getN(0);
  @$pb.TagNumber(52896097)
  set geomatchconstraint(GeoMatchConstraint value) =>
      $_setField(52896097, value);
  @$pb.TagNumber(52896097)
  $core.bool hasGeomatchconstraint() => $_has(0);
  @$pb.TagNumber(52896097)
  void clearGeomatchconstraint() => $_clearField(52896097);
  @$pb.TagNumber(52896097)
  GeoMatchConstraint ensureGeomatchconstraint() => $_ensure(0);

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(1);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(1);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
}

class GetByteMatchSetRequest extends $pb.GeneratedMessage {
  factory GetByteMatchSetRequest({
    $core.String? bytematchsetid,
  }) {
    final result = create();
    if (bytematchsetid != null) result.bytematchsetid = bytematchsetid;
    return result;
  }

  GetByteMatchSetRequest._();

  factory GetByteMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetByteMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetByteMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(62033398, _omitFieldNames ? '' : 'bytematchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByteMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByteMatchSetRequest copyWith(
          void Function(GetByteMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as GetByteMatchSetRequest))
          as GetByteMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetByteMatchSetRequest create() => GetByteMatchSetRequest._();
  @$core.override
  GetByteMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetByteMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetByteMatchSetRequest>(create);
  static GetByteMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(62033398)
  $core.String get bytematchsetid => $_getSZ(0);
  @$pb.TagNumber(62033398)
  set bytematchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(62033398)
  $core.bool hasBytematchsetid() => $_has(0);
  @$pb.TagNumber(62033398)
  void clearBytematchsetid() => $_clearField(62033398);
}

class GetByteMatchSetResponse extends $pb.GeneratedMessage {
  factory GetByteMatchSetResponse({
    ByteMatchSet? bytematchset,
  }) {
    final result = create();
    if (bytematchset != null) result.bytematchset = bytematchset;
    return result;
  }

  GetByteMatchSetResponse._();

  factory GetByteMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetByteMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetByteMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<ByteMatchSet>(9496489, _omitFieldNames ? '' : 'bytematchset',
        subBuilder: ByteMatchSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByteMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetByteMatchSetResponse copyWith(
          void Function(GetByteMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as GetByteMatchSetResponse))
          as GetByteMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetByteMatchSetResponse create() => GetByteMatchSetResponse._();
  @$core.override
  GetByteMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetByteMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetByteMatchSetResponse>(create);
  static GetByteMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(9496489)
  ByteMatchSet get bytematchset => $_getN(0);
  @$pb.TagNumber(9496489)
  set bytematchset(ByteMatchSet value) => $_setField(9496489, value);
  @$pb.TagNumber(9496489)
  $core.bool hasBytematchset() => $_has(0);
  @$pb.TagNumber(9496489)
  void clearBytematchset() => $_clearField(9496489);
  @$pb.TagNumber(9496489)
  ByteMatchSet ensureBytematchset() => $_ensure(0);
}

class GetChangeTokenRequest extends $pb.GeneratedMessage {
  factory GetChangeTokenRequest() => create();

  GetChangeTokenRequest._();

  factory GetChangeTokenRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetChangeTokenRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetChangeTokenRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeTokenRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeTokenRequest copyWith(
          void Function(GetChangeTokenRequest) updates) =>
      super.copyWith((message) => updates(message as GetChangeTokenRequest))
          as GetChangeTokenRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChangeTokenRequest create() => GetChangeTokenRequest._();
  @$core.override
  GetChangeTokenRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetChangeTokenRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetChangeTokenRequest>(create);
  static GetChangeTokenRequest? _defaultInstance;
}

class GetChangeTokenResponse extends $pb.GeneratedMessage {
  factory GetChangeTokenResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  GetChangeTokenResponse._();

  factory GetChangeTokenResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetChangeTokenResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetChangeTokenResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeTokenResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeTokenResponse copyWith(
          void Function(GetChangeTokenResponse) updates) =>
      super.copyWith((message) => updates(message as GetChangeTokenResponse))
          as GetChangeTokenResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChangeTokenResponse create() => GetChangeTokenResponse._();
  @$core.override
  GetChangeTokenResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetChangeTokenResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetChangeTokenResponse>(create);
  static GetChangeTokenResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class GetChangeTokenStatusRequest extends $pb.GeneratedMessage {
  factory GetChangeTokenStatusRequest({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  GetChangeTokenStatusRequest._();

  factory GetChangeTokenStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetChangeTokenStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetChangeTokenStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeTokenStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeTokenStatusRequest copyWith(
          void Function(GetChangeTokenStatusRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetChangeTokenStatusRequest))
          as GetChangeTokenStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChangeTokenStatusRequest create() =>
      GetChangeTokenStatusRequest._();
  @$core.override
  GetChangeTokenStatusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetChangeTokenStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetChangeTokenStatusRequest>(create);
  static GetChangeTokenStatusRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class GetChangeTokenStatusResponse extends $pb.GeneratedMessage {
  factory GetChangeTokenStatusResponse({
    ChangeTokenStatus? changetokenstatus,
  }) {
    final result = create();
    if (changetokenstatus != null) result.changetokenstatus = changetokenstatus;
    return result;
  }

  GetChangeTokenStatusResponse._();

  factory GetChangeTokenStatusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetChangeTokenStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetChangeTokenStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ChangeTokenStatus>(
        37644373, _omitFieldNames ? '' : 'changetokenstatus',
        enumValues: ChangeTokenStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeTokenStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeTokenStatusResponse copyWith(
          void Function(GetChangeTokenStatusResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetChangeTokenStatusResponse))
          as GetChangeTokenStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChangeTokenStatusResponse create() =>
      GetChangeTokenStatusResponse._();
  @$core.override
  GetChangeTokenStatusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetChangeTokenStatusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetChangeTokenStatusResponse>(create);
  static GetChangeTokenStatusResponse? _defaultInstance;

  @$pb.TagNumber(37644373)
  ChangeTokenStatus get changetokenstatus => $_getN(0);
  @$pb.TagNumber(37644373)
  set changetokenstatus(ChangeTokenStatus value) => $_setField(37644373, value);
  @$pb.TagNumber(37644373)
  $core.bool hasChangetokenstatus() => $_has(0);
  @$pb.TagNumber(37644373)
  void clearChangetokenstatus() => $_clearField(37644373);
}

class GetGeoMatchSetRequest extends $pb.GeneratedMessage {
  factory GetGeoMatchSetRequest({
    $core.String? geomatchsetid,
  }) {
    final result = create();
    if (geomatchsetid != null) result.geomatchsetid = geomatchsetid;
    return result;
  }

  GetGeoMatchSetRequest._();

  factory GetGeoMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetGeoMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetGeoMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(514346837, _omitFieldNames ? '' : 'geomatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGeoMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGeoMatchSetRequest copyWith(
          void Function(GetGeoMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as GetGeoMatchSetRequest))
          as GetGeoMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetGeoMatchSetRequest create() => GetGeoMatchSetRequest._();
  @$core.override
  GetGeoMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetGeoMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetGeoMatchSetRequest>(create);
  static GetGeoMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(514346837)
  $core.String get geomatchsetid => $_getSZ(0);
  @$pb.TagNumber(514346837)
  set geomatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(514346837)
  $core.bool hasGeomatchsetid() => $_has(0);
  @$pb.TagNumber(514346837)
  void clearGeomatchsetid() => $_clearField(514346837);
}

class GetGeoMatchSetResponse extends $pb.GeneratedMessage {
  factory GetGeoMatchSetResponse({
    GeoMatchSet? geomatchset,
  }) {
    final result = create();
    if (geomatchset != null) result.geomatchset = geomatchset;
    return result;
  }

  GetGeoMatchSetResponse._();

  factory GetGeoMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetGeoMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetGeoMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<GeoMatchSet>(84933666, _omitFieldNames ? '' : 'geomatchset',
        subBuilder: GeoMatchSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGeoMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGeoMatchSetResponse copyWith(
          void Function(GetGeoMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as GetGeoMatchSetResponse))
          as GetGeoMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetGeoMatchSetResponse create() => GetGeoMatchSetResponse._();
  @$core.override
  GetGeoMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetGeoMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetGeoMatchSetResponse>(create);
  static GetGeoMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(84933666)
  GeoMatchSet get geomatchset => $_getN(0);
  @$pb.TagNumber(84933666)
  set geomatchset(GeoMatchSet value) => $_setField(84933666, value);
  @$pb.TagNumber(84933666)
  $core.bool hasGeomatchset() => $_has(0);
  @$pb.TagNumber(84933666)
  void clearGeomatchset() => $_clearField(84933666);
  @$pb.TagNumber(84933666)
  GeoMatchSet ensureGeomatchset() => $_ensure(0);
}

class GetIPSetRequest extends $pb.GeneratedMessage {
  factory GetIPSetRequest({
    $core.String? ipsetid,
  }) {
    final result = create();
    if (ipsetid != null) result.ipsetid = ipsetid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(346763066, _omitFieldNames ? '' : 'ipsetid')
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

  @$pb.TagNumber(346763066)
  $core.String get ipsetid => $_getSZ(0);
  @$pb.TagNumber(346763066)
  set ipsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346763066)
  $core.bool hasIpsetid() => $_has(0);
  @$pb.TagNumber(346763066)
  void clearIpsetid() => $_clearField(346763066);
}

class GetIPSetResponse extends $pb.GeneratedMessage {
  factory GetIPSetResponse({
    IPSet? ipset,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(436412565)
  IPSet get ipset => $_getN(0);
  @$pb.TagNumber(436412565)
  set ipset(IPSet value) => $_setField(436412565, value);
  @$pb.TagNumber(436412565)
  $core.bool hasIpset() => $_has(0);
  @$pb.TagNumber(436412565)
  void clearIpset() => $_clearField(436412565);
  @$pb.TagNumber(436412565)
  IPSet ensureIpset() => $_ensure(0);
}

class GetLoggingConfigurationRequest extends $pb.GeneratedMessage {
  factory GetLoggingConfigurationRequest({
    $core.String? resourcearn,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(0);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class GetRateBasedRuleManagedKeysRequest extends $pb.GeneratedMessage {
  factory GetRateBasedRuleManagedKeysRequest({
    $core.String? ruleid,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (ruleid != null) result.ruleid = ruleid;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  GetRateBasedRuleManagedKeysRequest._();

  factory GetRateBasedRuleManagedKeysRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRateBasedRuleManagedKeysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRateBasedRuleManagedKeysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedRuleManagedKeysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedRuleManagedKeysRequest copyWith(
          void Function(GetRateBasedRuleManagedKeysRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetRateBasedRuleManagedKeysRequest))
          as GetRateBasedRuleManagedKeysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRateBasedRuleManagedKeysRequest create() =>
      GetRateBasedRuleManagedKeysRequest._();
  @$core.override
  GetRateBasedRuleManagedKeysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRateBasedRuleManagedKeysRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRateBasedRuleManagedKeysRequest>(
          create);
  static GetRateBasedRuleManagedKeysRequest? _defaultInstance;

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(0);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(0);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class GetRateBasedRuleManagedKeysResponse extends $pb.GeneratedMessage {
  factory GetRateBasedRuleManagedKeysResponse({
    $core.Iterable<$core.String>? managedkeys,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (managedkeys != null) result.managedkeys.addAll(managedkeys);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  GetRateBasedRuleManagedKeysResponse._();

  factory GetRateBasedRuleManagedKeysResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRateBasedRuleManagedKeysResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRateBasedRuleManagedKeysResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPS(225455355, _omitFieldNames ? '' : 'managedkeys')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedRuleManagedKeysResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedRuleManagedKeysResponse copyWith(
          void Function(GetRateBasedRuleManagedKeysResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetRateBasedRuleManagedKeysResponse))
          as GetRateBasedRuleManagedKeysResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRateBasedRuleManagedKeysResponse create() =>
      GetRateBasedRuleManagedKeysResponse._();
  @$core.override
  GetRateBasedRuleManagedKeysResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRateBasedRuleManagedKeysResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetRateBasedRuleManagedKeysResponse>(create);
  static GetRateBasedRuleManagedKeysResponse? _defaultInstance;

  @$pb.TagNumber(225455355)
  $pb.PbList<$core.String> get managedkeys => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class GetRateBasedRuleRequest extends $pb.GeneratedMessage {
  factory GetRateBasedRuleRequest({
    $core.String? ruleid,
  }) {
    final result = create();
    if (ruleid != null) result.ruleid = ruleid;
    return result;
  }

  GetRateBasedRuleRequest._();

  factory GetRateBasedRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRateBasedRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRateBasedRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedRuleRequest copyWith(
          void Function(GetRateBasedRuleRequest) updates) =>
      super.copyWith((message) => updates(message as GetRateBasedRuleRequest))
          as GetRateBasedRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRateBasedRuleRequest create() => GetRateBasedRuleRequest._();
  @$core.override
  GetRateBasedRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRateBasedRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRateBasedRuleRequest>(create);
  static GetRateBasedRuleRequest? _defaultInstance;

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(0);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(0);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class GetRateBasedRuleResponse extends $pb.GeneratedMessage {
  factory GetRateBasedRuleResponse({
    RateBasedRule? rule,
  }) {
    final result = create();
    if (rule != null) result.rule = rule;
    return result;
  }

  GetRateBasedRuleResponse._();

  factory GetRateBasedRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRateBasedRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRateBasedRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<RateBasedRule>(475696372, _omitFieldNames ? '' : 'rule',
        subBuilder: RateBasedRule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRateBasedRuleResponse copyWith(
          void Function(GetRateBasedRuleResponse) updates) =>
      super.copyWith((message) => updates(message as GetRateBasedRuleResponse))
          as GetRateBasedRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRateBasedRuleResponse create() => GetRateBasedRuleResponse._();
  @$core.override
  GetRateBasedRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRateBasedRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRateBasedRuleResponse>(create);
  static GetRateBasedRuleResponse? _defaultInstance;

  @$pb.TagNumber(475696372)
  RateBasedRule get rule => $_getN(0);
  @$pb.TagNumber(475696372)
  set rule(RateBasedRule value) => $_setField(475696372, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(0);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);
  @$pb.TagNumber(475696372)
  RateBasedRule ensureRule() => $_ensure(0);
}

class GetRegexMatchSetRequest extends $pb.GeneratedMessage {
  factory GetRegexMatchSetRequest({
    $core.String? regexmatchsetid,
  }) {
    final result = create();
    if (regexmatchsetid != null) result.regexmatchsetid = regexmatchsetid;
    return result;
  }

  GetRegexMatchSetRequest._();

  factory GetRegexMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRegexMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRegexMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82287635, _omitFieldNames ? '' : 'regexmatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRegexMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRegexMatchSetRequest copyWith(
          void Function(GetRegexMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as GetRegexMatchSetRequest))
          as GetRegexMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRegexMatchSetRequest create() => GetRegexMatchSetRequest._();
  @$core.override
  GetRegexMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRegexMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRegexMatchSetRequest>(create);
  static GetRegexMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(82287635)
  $core.String get regexmatchsetid => $_getSZ(0);
  @$pb.TagNumber(82287635)
  set regexmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82287635)
  $core.bool hasRegexmatchsetid() => $_has(0);
  @$pb.TagNumber(82287635)
  void clearRegexmatchsetid() => $_clearField(82287635);
}

class GetRegexMatchSetResponse extends $pb.GeneratedMessage {
  factory GetRegexMatchSetResponse({
    RegexMatchSet? regexmatchset,
  }) {
    final result = create();
    if (regexmatchset != null) result.regexmatchset = regexmatchset;
    return result;
  }

  GetRegexMatchSetResponse._();

  factory GetRegexMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRegexMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRegexMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<RegexMatchSet>(534067456, _omitFieldNames ? '' : 'regexmatchset',
        subBuilder: RegexMatchSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRegexMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRegexMatchSetResponse copyWith(
          void Function(GetRegexMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as GetRegexMatchSetResponse))
          as GetRegexMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRegexMatchSetResponse create() => GetRegexMatchSetResponse._();
  @$core.override
  GetRegexMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRegexMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRegexMatchSetResponse>(create);
  static GetRegexMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(534067456)
  RegexMatchSet get regexmatchset => $_getN(0);
  @$pb.TagNumber(534067456)
  set regexmatchset(RegexMatchSet value) => $_setField(534067456, value);
  @$pb.TagNumber(534067456)
  $core.bool hasRegexmatchset() => $_has(0);
  @$pb.TagNumber(534067456)
  void clearRegexmatchset() => $_clearField(534067456);
  @$pb.TagNumber(534067456)
  RegexMatchSet ensureRegexmatchset() => $_ensure(0);
}

class GetRegexPatternSetRequest extends $pb.GeneratedMessage {
  factory GetRegexPatternSetRequest({
    $core.String? regexpatternsetid,
  }) {
    final result = create();
    if (regexpatternsetid != null) result.regexpatternsetid = regexpatternsetid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(96220912, _omitFieldNames ? '' : 'regexpatternsetid')
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

  @$pb.TagNumber(96220912)
  $core.String get regexpatternsetid => $_getSZ(0);
  @$pb.TagNumber(96220912)
  set regexpatternsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(96220912)
  $core.bool hasRegexpatternsetid() => $_has(0);
  @$pb.TagNumber(96220912)
  void clearRegexpatternsetid() => $_clearField(96220912);
}

class GetRegexPatternSetResponse extends $pb.GeneratedMessage {
  factory GetRegexPatternSetResponse({
    RegexPatternSet? regexpatternset,
  }) {
    final result = create();
    if (regexpatternset != null) result.regexpatternset = regexpatternset;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<RegexPatternSet>(9374915, _omitFieldNames ? '' : 'regexpatternset',
        subBuilder: RegexPatternSet.create)
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
}

class GetRuleGroupRequest extends $pb.GeneratedMessage {
  factory GetRuleGroupRequest({
    $core.String? rulegroupid,
  }) {
    final result = create();
    if (rulegroupid != null) result.rulegroupid = rulegroupid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(287757320, _omitFieldNames ? '' : 'rulegroupid')
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

  @$pb.TagNumber(287757320)
  $core.String get rulegroupid => $_getSZ(0);
  @$pb.TagNumber(287757320)
  set rulegroupid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(287757320)
  $core.bool hasRulegroupid() => $_has(0);
  @$pb.TagNumber(287757320)
  void clearRulegroupid() => $_clearField(287757320);
}

class GetRuleGroupResponse extends $pb.GeneratedMessage {
  factory GetRuleGroupResponse({
    RuleGroup? rulegroup,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(267398571)
  RuleGroup get rulegroup => $_getN(0);
  @$pb.TagNumber(267398571)
  set rulegroup(RuleGroup value) => $_setField(267398571, value);
  @$pb.TagNumber(267398571)
  $core.bool hasRulegroup() => $_has(0);
  @$pb.TagNumber(267398571)
  void clearRulegroup() => $_clearField(267398571);
  @$pb.TagNumber(267398571)
  RuleGroup ensureRulegroup() => $_ensure(0);
}

class GetRuleRequest extends $pb.GeneratedMessage {
  factory GetRuleRequest({
    $core.String? ruleid,
  }) {
    final result = create();
    if (ruleid != null) result.ruleid = ruleid;
    return result;
  }

  GetRuleRequest._();

  factory GetRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRuleRequest copyWith(void Function(GetRuleRequest) updates) =>
      super.copyWith((message) => updates(message as GetRuleRequest))
          as GetRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRuleRequest create() => GetRuleRequest._();
  @$core.override
  GetRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRuleRequest>(create);
  static GetRuleRequest? _defaultInstance;

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(0);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(0);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class GetRuleResponse extends $pb.GeneratedMessage {
  factory GetRuleResponse({
    Rule? rule,
  }) {
    final result = create();
    if (rule != null) result.rule = rule;
    return result;
  }

  GetRuleResponse._();

  factory GetRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<Rule>(475696372, _omitFieldNames ? '' : 'rule',
        subBuilder: Rule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRuleResponse copyWith(void Function(GetRuleResponse) updates) =>
      super.copyWith((message) => updates(message as GetRuleResponse))
          as GetRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRuleResponse create() => GetRuleResponse._();
  @$core.override
  GetRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRuleResponse>(create);
  static GetRuleResponse? _defaultInstance;

  @$pb.TagNumber(475696372)
  Rule get rule => $_getN(0);
  @$pb.TagNumber(475696372)
  set rule(Rule value) => $_setField(475696372, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(0);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);
  @$pb.TagNumber(475696372)
  Rule ensureRule() => $_ensure(0);
}

class GetSampledRequestsRequest extends $pb.GeneratedMessage {
  factory GetSampledRequestsRequest({
    TimeWindow? timewindow,
    $core.String? webaclid,
    $core.String? ruleid,
    $fixnum.Int64? maxitems,
  }) {
    final result = create();
    if (timewindow != null) result.timewindow = timewindow;
    if (webaclid != null) result.webaclid = webaclid;
    if (ruleid != null) result.ruleid = ruleid;
    if (maxitems != null) result.maxitems = maxitems;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<TimeWindow>(140543513, _omitFieldNames ? '' : 'timewindow',
        subBuilder: TimeWindow.create)
    ..aOS(201008723, _omitFieldNames ? '' : 'webaclid')
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..aInt64(506899220, _omitFieldNames ? '' : 'maxitems')
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

  @$pb.TagNumber(140543513)
  TimeWindow get timewindow => $_getN(0);
  @$pb.TagNumber(140543513)
  set timewindow(TimeWindow value) => $_setField(140543513, value);
  @$pb.TagNumber(140543513)
  $core.bool hasTimewindow() => $_has(0);
  @$pb.TagNumber(140543513)
  void clearTimewindow() => $_clearField(140543513);
  @$pb.TagNumber(140543513)
  TimeWindow ensureTimewindow() => $_ensure(0);

  @$pb.TagNumber(201008723)
  $core.String get webaclid => $_getSZ(1);
  @$pb.TagNumber(201008723)
  set webaclid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(201008723)
  $core.bool hasWebaclid() => $_has(1);
  @$pb.TagNumber(201008723)
  void clearWebaclid() => $_clearField(201008723);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(2);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(2);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);

  @$pb.TagNumber(506899220)
  $fixnum.Int64 get maxitems => $_getI64(3);
  @$pb.TagNumber(506899220)
  set maxitems($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class GetSizeConstraintSetRequest extends $pb.GeneratedMessage {
  factory GetSizeConstraintSetRequest({
    $core.String? sizeconstraintsetid,
  }) {
    final result = create();
    if (sizeconstraintsetid != null)
      result.sizeconstraintsetid = sizeconstraintsetid;
    return result;
  }

  GetSizeConstraintSetRequest._();

  factory GetSizeConstraintSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSizeConstraintSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSizeConstraintSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(277959757, _omitFieldNames ? '' : 'sizeconstraintsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSizeConstraintSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSizeConstraintSetRequest copyWith(
          void Function(GetSizeConstraintSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetSizeConstraintSetRequest))
          as GetSizeConstraintSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSizeConstraintSetRequest create() =>
      GetSizeConstraintSetRequest._();
  @$core.override
  GetSizeConstraintSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSizeConstraintSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSizeConstraintSetRequest>(create);
  static GetSizeConstraintSetRequest? _defaultInstance;

  @$pb.TagNumber(277959757)
  $core.String get sizeconstraintsetid => $_getSZ(0);
  @$pb.TagNumber(277959757)
  set sizeconstraintsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(277959757)
  $core.bool hasSizeconstraintsetid() => $_has(0);
  @$pb.TagNumber(277959757)
  void clearSizeconstraintsetid() => $_clearField(277959757);
}

class GetSizeConstraintSetResponse extends $pb.GeneratedMessage {
  factory GetSizeConstraintSetResponse({
    SizeConstraintSet? sizeconstraintset,
  }) {
    final result = create();
    if (sizeconstraintset != null) result.sizeconstraintset = sizeconstraintset;
    return result;
  }

  GetSizeConstraintSetResponse._();

  factory GetSizeConstraintSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSizeConstraintSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSizeConstraintSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<SizeConstraintSet>(
        166385770, _omitFieldNames ? '' : 'sizeconstraintset',
        subBuilder: SizeConstraintSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSizeConstraintSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSizeConstraintSetResponse copyWith(
          void Function(GetSizeConstraintSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetSizeConstraintSetResponse))
          as GetSizeConstraintSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSizeConstraintSetResponse create() =>
      GetSizeConstraintSetResponse._();
  @$core.override
  GetSizeConstraintSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSizeConstraintSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSizeConstraintSetResponse>(create);
  static GetSizeConstraintSetResponse? _defaultInstance;

  @$pb.TagNumber(166385770)
  SizeConstraintSet get sizeconstraintset => $_getN(0);
  @$pb.TagNumber(166385770)
  set sizeconstraintset(SizeConstraintSet value) =>
      $_setField(166385770, value);
  @$pb.TagNumber(166385770)
  $core.bool hasSizeconstraintset() => $_has(0);
  @$pb.TagNumber(166385770)
  void clearSizeconstraintset() => $_clearField(166385770);
  @$pb.TagNumber(166385770)
  SizeConstraintSet ensureSizeconstraintset() => $_ensure(0);
}

class GetSqlInjectionMatchSetRequest extends $pb.GeneratedMessage {
  factory GetSqlInjectionMatchSetRequest({
    $core.String? sqlinjectionmatchsetid,
  }) {
    final result = create();
    if (sqlinjectionmatchsetid != null)
      result.sqlinjectionmatchsetid = sqlinjectionmatchsetid;
    return result;
  }

  GetSqlInjectionMatchSetRequest._();

  factory GetSqlInjectionMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSqlInjectionMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSqlInjectionMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(386225735, _omitFieldNames ? '' : 'sqlinjectionmatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSqlInjectionMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSqlInjectionMatchSetRequest copyWith(
          void Function(GetSqlInjectionMatchSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetSqlInjectionMatchSetRequest))
          as GetSqlInjectionMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSqlInjectionMatchSetRequest create() =>
      GetSqlInjectionMatchSetRequest._();
  @$core.override
  GetSqlInjectionMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSqlInjectionMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSqlInjectionMatchSetRequest>(create);
  static GetSqlInjectionMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(386225735)
  $core.String get sqlinjectionmatchsetid => $_getSZ(0);
  @$pb.TagNumber(386225735)
  set sqlinjectionmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(386225735)
  $core.bool hasSqlinjectionmatchsetid() => $_has(0);
  @$pb.TagNumber(386225735)
  void clearSqlinjectionmatchsetid() => $_clearField(386225735);
}

class GetSqlInjectionMatchSetResponse extends $pb.GeneratedMessage {
  factory GetSqlInjectionMatchSetResponse({
    SqlInjectionMatchSet? sqlinjectionmatchset,
  }) {
    final result = create();
    if (sqlinjectionmatchset != null)
      result.sqlinjectionmatchset = sqlinjectionmatchset;
    return result;
  }

  GetSqlInjectionMatchSetResponse._();

  factory GetSqlInjectionMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetSqlInjectionMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetSqlInjectionMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<SqlInjectionMatchSet>(
        60356588, _omitFieldNames ? '' : 'sqlinjectionmatchset',
        subBuilder: SqlInjectionMatchSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSqlInjectionMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetSqlInjectionMatchSetResponse copyWith(
          void Function(GetSqlInjectionMatchSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetSqlInjectionMatchSetResponse))
          as GetSqlInjectionMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetSqlInjectionMatchSetResponse create() =>
      GetSqlInjectionMatchSetResponse._();
  @$core.override
  GetSqlInjectionMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetSqlInjectionMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetSqlInjectionMatchSetResponse>(
          create);
  static GetSqlInjectionMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(60356588)
  SqlInjectionMatchSet get sqlinjectionmatchset => $_getN(0);
  @$pb.TagNumber(60356588)
  set sqlinjectionmatchset(SqlInjectionMatchSet value) =>
      $_setField(60356588, value);
  @$pb.TagNumber(60356588)
  $core.bool hasSqlinjectionmatchset() => $_has(0);
  @$pb.TagNumber(60356588)
  void clearSqlinjectionmatchset() => $_clearField(60356588);
  @$pb.TagNumber(60356588)
  SqlInjectionMatchSet ensureSqlinjectionmatchset() => $_ensure(0);
}

class GetWebACLRequest extends $pb.GeneratedMessage {
  factory GetWebACLRequest({
    $core.String? webaclid,
  }) {
    final result = create();
    if (webaclid != null) result.webaclid = webaclid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(161274579, _omitFieldNames ? '' : 'webaclid')
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

  @$pb.TagNumber(161274579)
  $core.String get webaclid => $_getSZ(0);
  @$pb.TagNumber(161274579)
  set webaclid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(161274579)
  $core.bool hasWebaclid() => $_has(0);
  @$pb.TagNumber(161274579)
  void clearWebaclid() => $_clearField(161274579);
}

class GetWebACLResponse extends $pb.GeneratedMessage {
  factory GetWebACLResponse({
    WebACL? webacl,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

class GetXssMatchSetRequest extends $pb.GeneratedMessage {
  factory GetXssMatchSetRequest({
    $core.String? xssmatchsetid,
  }) {
    final result = create();
    if (xssmatchsetid != null) result.xssmatchsetid = xssmatchsetid;
    return result;
  }

  GetXssMatchSetRequest._();

  factory GetXssMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetXssMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetXssMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(56643900, _omitFieldNames ? '' : 'xssmatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetXssMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetXssMatchSetRequest copyWith(
          void Function(GetXssMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as GetXssMatchSetRequest))
          as GetXssMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetXssMatchSetRequest create() => GetXssMatchSetRequest._();
  @$core.override
  GetXssMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetXssMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetXssMatchSetRequest>(create);
  static GetXssMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(56643900)
  $core.String get xssmatchsetid => $_getSZ(0);
  @$pb.TagNumber(56643900)
  set xssmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56643900)
  $core.bool hasXssmatchsetid() => $_has(0);
  @$pb.TagNumber(56643900)
  void clearXssmatchsetid() => $_clearField(56643900);
}

class GetXssMatchSetResponse extends $pb.GeneratedMessage {
  factory GetXssMatchSetResponse({
    XssMatchSet? xssmatchset,
  }) {
    final result = create();
    if (xssmatchset != null) result.xssmatchset = xssmatchset;
    return result;
  }

  GetXssMatchSetResponse._();

  factory GetXssMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetXssMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetXssMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<XssMatchSet>(169555847, _omitFieldNames ? '' : 'xssmatchset',
        subBuilder: XssMatchSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetXssMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetXssMatchSetResponse copyWith(
          void Function(GetXssMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as GetXssMatchSetResponse))
          as GetXssMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetXssMatchSetResponse create() => GetXssMatchSetResponse._();
  @$core.override
  GetXssMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetXssMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetXssMatchSetResponse>(create);
  static GetXssMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(169555847)
  XssMatchSet get xssmatchset => $_getN(0);
  @$pb.TagNumber(169555847)
  set xssmatchset(XssMatchSet value) => $_setField(169555847, value);
  @$pb.TagNumber(169555847)
  $core.bool hasXssmatchset() => $_has(0);
  @$pb.TagNumber(169555847)
  void clearXssmatchset() => $_clearField(169555847);
  @$pb.TagNumber(169555847)
  XssMatchSet ensureXssmatchset() => $_ensure(0);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class IPSet extends $pb.GeneratedMessage {
  factory IPSet({
    $core.String? name,
    $core.Iterable<IPSetDescriptor>? ipsetdescriptors,
    $core.String? ipsetid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (ipsetdescriptors != null)
      result.ipsetdescriptors.addAll(ipsetdescriptors);
    if (ipsetid != null) result.ipsetid = ipsetid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<IPSetDescriptor>(320792515, _omitFieldNames ? '' : 'ipsetdescriptors',
        subBuilder: IPSetDescriptor.create)
    ..aOS(346763066, _omitFieldNames ? '' : 'ipsetid')
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

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(320792515)
  $pb.PbList<IPSetDescriptor> get ipsetdescriptors => $_getList(1);

  @$pb.TagNumber(346763066)
  $core.String get ipsetid => $_getSZ(2);
  @$pb.TagNumber(346763066)
  set ipsetid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(346763066)
  $core.bool hasIpsetid() => $_has(2);
  @$pb.TagNumber(346763066)
  void clearIpsetid() => $_clearField(346763066);
}

class IPSetDescriptor extends $pb.GeneratedMessage {
  factory IPSetDescriptor({
    $core.String? value,
    IPSetDescriptorType? type,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (type != null) result.type = type;
    return result;
  }

  IPSetDescriptor._();

  factory IPSetDescriptor.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IPSetDescriptor.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IPSetDescriptor',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..aE<IPSetDescriptorType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: IPSetDescriptorType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetDescriptor clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetDescriptor copyWith(void Function(IPSetDescriptor) updates) =>
      super.copyWith((message) => updates(message as IPSetDescriptor))
          as IPSetDescriptor;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IPSetDescriptor create() => IPSetDescriptor._();
  @$core.override
  IPSetDescriptor createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IPSetDescriptor getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IPSetDescriptor>(create);
  static IPSetDescriptor? _defaultInstance;

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(290836590)
  IPSetDescriptorType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(IPSetDescriptorType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class IPSetSummary extends $pb.GeneratedMessage {
  factory IPSetSummary({
    $core.String? name,
    $core.String? ipsetid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (ipsetid != null) result.ipsetid = ipsetid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(346763066, _omitFieldNames ? '' : 'ipsetid')
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

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(346763066)
  $core.String get ipsetid => $_getSZ(1);
  @$pb.TagNumber(346763066)
  set ipsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346763066)
  $core.bool hasIpsetid() => $_has(1);
  @$pb.TagNumber(346763066)
  void clearIpsetid() => $_clearField(346763066);
}

class IPSetUpdate extends $pb.GeneratedMessage {
  factory IPSetUpdate({
    ChangeAction? action,
    IPSetDescriptor? ipsetdescriptor,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (ipsetdescriptor != null) result.ipsetdescriptor = ipsetdescriptor;
    return result;
  }

  IPSetUpdate._();

  factory IPSetUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IPSetUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IPSetUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..aOM<IPSetDescriptor>(226351622, _omitFieldNames ? '' : 'ipsetdescriptor',
        subBuilder: IPSetDescriptor.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IPSetUpdate copyWith(void Function(IPSetUpdate) updates) =>
      super.copyWith((message) => updates(message as IPSetUpdate))
          as IPSetUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IPSetUpdate create() => IPSetUpdate._();
  @$core.override
  IPSetUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IPSetUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IPSetUpdate>(create);
  static IPSetUpdate? _defaultInstance;

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(226351622)
  IPSetDescriptor get ipsetdescriptor => $_getN(1);
  @$pb.TagNumber(226351622)
  set ipsetdescriptor(IPSetDescriptor value) => $_setField(226351622, value);
  @$pb.TagNumber(226351622)
  $core.bool hasIpsetdescriptor() => $_has(1);
  @$pb.TagNumber(226351622)
  void clearIpsetdescriptor() => $_clearField(226351622);
  @$pb.TagNumber(226351622)
  IPSetDescriptor ensureIpsetdescriptor() => $_ensure(1);
}

class ListActivatedRulesInRuleGroupRequest extends $pb.GeneratedMessage {
  factory ListActivatedRulesInRuleGroupRequest({
    $core.String? rulegroupid,
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (rulegroupid != null) result.rulegroupid = rulegroupid;
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListActivatedRulesInRuleGroupRequest._();

  factory ListActivatedRulesInRuleGroupRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListActivatedRulesInRuleGroupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListActivatedRulesInRuleGroupRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(287757320, _omitFieldNames ? '' : 'rulegroupid')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListActivatedRulesInRuleGroupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListActivatedRulesInRuleGroupRequest copyWith(
          void Function(ListActivatedRulesInRuleGroupRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ListActivatedRulesInRuleGroupRequest))
          as ListActivatedRulesInRuleGroupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListActivatedRulesInRuleGroupRequest create() =>
      ListActivatedRulesInRuleGroupRequest._();
  @$core.override
  ListActivatedRulesInRuleGroupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListActivatedRulesInRuleGroupRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListActivatedRulesInRuleGroupRequest>(create);
  static ListActivatedRulesInRuleGroupRequest? _defaultInstance;

  @$pb.TagNumber(287757320)
  $core.String get rulegroupid => $_getSZ(0);
  @$pb.TagNumber(287757320)
  set rulegroupid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(287757320)
  $core.bool hasRulegroupid() => $_has(0);
  @$pb.TagNumber(287757320)
  void clearRulegroupid() => $_clearField(287757320);

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

class ListActivatedRulesInRuleGroupResponse extends $pb.GeneratedMessage {
  factory ListActivatedRulesInRuleGroupResponse({
    $core.Iterable<ActivatedRule>? activatedrules,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (activatedrules != null) result.activatedrules.addAll(activatedrules);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListActivatedRulesInRuleGroupResponse._();

  factory ListActivatedRulesInRuleGroupResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListActivatedRulesInRuleGroupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListActivatedRulesInRuleGroupResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<ActivatedRule>(75723620, _omitFieldNames ? '' : 'activatedrules',
        subBuilder: ActivatedRule.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListActivatedRulesInRuleGroupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListActivatedRulesInRuleGroupResponse copyWith(
          void Function(ListActivatedRulesInRuleGroupResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListActivatedRulesInRuleGroupResponse))
          as ListActivatedRulesInRuleGroupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListActivatedRulesInRuleGroupResponse create() =>
      ListActivatedRulesInRuleGroupResponse._();
  @$core.override
  ListActivatedRulesInRuleGroupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListActivatedRulesInRuleGroupResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListActivatedRulesInRuleGroupResponse>(create);
  static ListActivatedRulesInRuleGroupResponse? _defaultInstance;

  @$pb.TagNumber(75723620)
  $pb.PbList<ActivatedRule> get activatedrules => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListByteMatchSetsRequest extends $pb.GeneratedMessage {
  factory ListByteMatchSetsRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListByteMatchSetsRequest._();

  factory ListByteMatchSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListByteMatchSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListByteMatchSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListByteMatchSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListByteMatchSetsRequest copyWith(
          void Function(ListByteMatchSetsRequest) updates) =>
      super.copyWith((message) => updates(message as ListByteMatchSetsRequest))
          as ListByteMatchSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListByteMatchSetsRequest create() => ListByteMatchSetsRequest._();
  @$core.override
  ListByteMatchSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListByteMatchSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListByteMatchSetsRequest>(create);
  static ListByteMatchSetsRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListByteMatchSetsResponse extends $pb.GeneratedMessage {
  factory ListByteMatchSetsResponse({
    $core.Iterable<ByteMatchSetSummary>? bytematchsets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (bytematchsets != null) result.bytematchsets.addAll(bytematchsets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListByteMatchSetsResponse._();

  factory ListByteMatchSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListByteMatchSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListByteMatchSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<ByteMatchSetSummary>(
        521993666, _omitFieldNames ? '' : 'bytematchsets',
        subBuilder: ByteMatchSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListByteMatchSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListByteMatchSetsResponse copyWith(
          void Function(ListByteMatchSetsResponse) updates) =>
      super.copyWith((message) => updates(message as ListByteMatchSetsResponse))
          as ListByteMatchSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListByteMatchSetsResponse create() => ListByteMatchSetsResponse._();
  @$core.override
  ListByteMatchSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListByteMatchSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListByteMatchSetsResponse>(create);
  static ListByteMatchSetsResponse? _defaultInstance;

  @$pb.TagNumber(521993666)
  $pb.PbList<ByteMatchSetSummary> get bytematchsets => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListGeoMatchSetsRequest extends $pb.GeneratedMessage {
  factory ListGeoMatchSetsRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListGeoMatchSetsRequest._();

  factory ListGeoMatchSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListGeoMatchSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListGeoMatchSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGeoMatchSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGeoMatchSetsRequest copyWith(
          void Function(ListGeoMatchSetsRequest) updates) =>
      super.copyWith((message) => updates(message as ListGeoMatchSetsRequest))
          as ListGeoMatchSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListGeoMatchSetsRequest create() => ListGeoMatchSetsRequest._();
  @$core.override
  ListGeoMatchSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListGeoMatchSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListGeoMatchSetsRequest>(create);
  static ListGeoMatchSetsRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListGeoMatchSetsResponse extends $pb.GeneratedMessage {
  factory ListGeoMatchSetsResponse({
    $core.Iterable<GeoMatchSetSummary>? geomatchsets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (geomatchsets != null) result.geomatchsets.addAll(geomatchsets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListGeoMatchSetsResponse._();

  factory ListGeoMatchSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListGeoMatchSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListGeoMatchSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<GeoMatchSetSummary>(170538263, _omitFieldNames ? '' : 'geomatchsets',
        subBuilder: GeoMatchSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGeoMatchSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGeoMatchSetsResponse copyWith(
          void Function(ListGeoMatchSetsResponse) updates) =>
      super.copyWith((message) => updates(message as ListGeoMatchSetsResponse))
          as ListGeoMatchSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListGeoMatchSetsResponse create() => ListGeoMatchSetsResponse._();
  @$core.override
  ListGeoMatchSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListGeoMatchSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListGeoMatchSetsResponse>(create);
  static ListGeoMatchSetsResponse? _defaultInstance;

  @$pb.TagNumber(170538263)
  $pb.PbList<GeoMatchSetSummary> get geomatchsets => $_getList(0);

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
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class ListRateBasedRulesRequest extends $pb.GeneratedMessage {
  factory ListRateBasedRulesRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRateBasedRulesRequest._();

  factory ListRateBasedRulesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRateBasedRulesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRateBasedRulesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRateBasedRulesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRateBasedRulesRequest copyWith(
          void Function(ListRateBasedRulesRequest) updates) =>
      super.copyWith((message) => updates(message as ListRateBasedRulesRequest))
          as ListRateBasedRulesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRateBasedRulesRequest create() => ListRateBasedRulesRequest._();
  @$core.override
  ListRateBasedRulesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRateBasedRulesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRateBasedRulesRequest>(create);
  static ListRateBasedRulesRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListRateBasedRulesResponse extends $pb.GeneratedMessage {
  factory ListRateBasedRulesResponse({
    $core.Iterable<RuleSummary>? rules,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (rules != null) result.rules.addAll(rules);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRateBasedRulesResponse._();

  factory ListRateBasedRulesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRateBasedRulesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRateBasedRulesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<RuleSummary>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: RuleSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRateBasedRulesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRateBasedRulesResponse copyWith(
          void Function(ListRateBasedRulesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListRateBasedRulesResponse))
          as ListRateBasedRulesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRateBasedRulesResponse create() => ListRateBasedRulesResponse._();
  @$core.override
  ListRateBasedRulesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRateBasedRulesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRateBasedRulesResponse>(create);
  static ListRateBasedRulesResponse? _defaultInstance;

  @$pb.TagNumber(42675585)
  $pb.PbList<RuleSummary> get rules => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListRegexMatchSetsRequest extends $pb.GeneratedMessage {
  factory ListRegexMatchSetsRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRegexMatchSetsRequest._();

  factory ListRegexMatchSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRegexMatchSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRegexMatchSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRegexMatchSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRegexMatchSetsRequest copyWith(
          void Function(ListRegexMatchSetsRequest) updates) =>
      super.copyWith((message) => updates(message as ListRegexMatchSetsRequest))
          as ListRegexMatchSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRegexMatchSetsRequest create() => ListRegexMatchSetsRequest._();
  @$core.override
  ListRegexMatchSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRegexMatchSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRegexMatchSetsRequest>(create);
  static ListRegexMatchSetsRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListRegexMatchSetsResponse extends $pb.GeneratedMessage {
  factory ListRegexMatchSetsResponse({
    $core.Iterable<RegexMatchSetSummary>? regexmatchsets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (regexmatchsets != null) result.regexmatchsets.addAll(regexmatchsets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRegexMatchSetsResponse._();

  factory ListRegexMatchSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRegexMatchSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRegexMatchSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<RegexMatchSetSummary>(
        145228901, _omitFieldNames ? '' : 'regexmatchsets',
        subBuilder: RegexMatchSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRegexMatchSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRegexMatchSetsResponse copyWith(
          void Function(ListRegexMatchSetsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListRegexMatchSetsResponse))
          as ListRegexMatchSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRegexMatchSetsResponse create() => ListRegexMatchSetsResponse._();
  @$core.override
  ListRegexMatchSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRegexMatchSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRegexMatchSetsResponse>(create);
  static ListRegexMatchSetsResponse? _defaultInstance;

  @$pb.TagNumber(145228901)
  $pb.PbList<RegexMatchSetSummary> get regexmatchsets => $_getList(0);

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
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class ListRuleGroupsRequest extends $pb.GeneratedMessage {
  factory ListRuleGroupsRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class ListRulesRequest extends $pb.GeneratedMessage {
  factory ListRulesRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRulesRequest._();

  factory ListRulesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRulesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRulesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRulesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRulesRequest copyWith(void Function(ListRulesRequest) updates) =>
      super.copyWith((message) => updates(message as ListRulesRequest))
          as ListRulesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRulesRequest create() => ListRulesRequest._();
  @$core.override
  ListRulesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRulesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRulesRequest>(create);
  static ListRulesRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListRulesResponse extends $pb.GeneratedMessage {
  factory ListRulesResponse({
    $core.Iterable<RuleSummary>? rules,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (rules != null) result.rules.addAll(rules);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListRulesResponse._();

  factory ListRulesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRulesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRulesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<RuleSummary>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: RuleSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRulesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRulesResponse copyWith(void Function(ListRulesResponse) updates) =>
      super.copyWith((message) => updates(message as ListRulesResponse))
          as ListRulesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRulesResponse create() => ListRulesResponse._();
  @$core.override
  ListRulesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRulesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRulesResponse>(create);
  static ListRulesResponse? _defaultInstance;

  @$pb.TagNumber(42675585)
  $pb.PbList<RuleSummary> get rules => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListSizeConstraintSetsRequest extends $pb.GeneratedMessage {
  factory ListSizeConstraintSetsRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListSizeConstraintSetsRequest._();

  factory ListSizeConstraintSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSizeConstraintSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSizeConstraintSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSizeConstraintSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSizeConstraintSetsRequest copyWith(
          void Function(ListSizeConstraintSetsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListSizeConstraintSetsRequest))
          as ListSizeConstraintSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSizeConstraintSetsRequest create() =>
      ListSizeConstraintSetsRequest._();
  @$core.override
  ListSizeConstraintSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSizeConstraintSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSizeConstraintSetsRequest>(create);
  static ListSizeConstraintSetsRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListSizeConstraintSetsResponse extends $pb.GeneratedMessage {
  factory ListSizeConstraintSetsResponse({
    $core.Iterable<SizeConstraintSetSummary>? sizeconstraintsets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (sizeconstraintsets != null)
      result.sizeconstraintsets.addAll(sizeconstraintsets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListSizeConstraintSetsResponse._();

  factory ListSizeConstraintSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSizeConstraintSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSizeConstraintSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<SizeConstraintSetSummary>(
        380776687, _omitFieldNames ? '' : 'sizeconstraintsets',
        subBuilder: SizeConstraintSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSizeConstraintSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSizeConstraintSetsResponse copyWith(
          void Function(ListSizeConstraintSetsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListSizeConstraintSetsResponse))
          as ListSizeConstraintSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSizeConstraintSetsResponse create() =>
      ListSizeConstraintSetsResponse._();
  @$core.override
  ListSizeConstraintSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSizeConstraintSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSizeConstraintSetsResponse>(create);
  static ListSizeConstraintSetsResponse? _defaultInstance;

  @$pb.TagNumber(380776687)
  $pb.PbList<SizeConstraintSetSummary> get sizeconstraintsets => $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListSqlInjectionMatchSetsRequest extends $pb.GeneratedMessage {
  factory ListSqlInjectionMatchSetsRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListSqlInjectionMatchSetsRequest._();

  factory ListSqlInjectionMatchSetsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSqlInjectionMatchSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSqlInjectionMatchSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSqlInjectionMatchSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSqlInjectionMatchSetsRequest copyWith(
          void Function(ListSqlInjectionMatchSetsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListSqlInjectionMatchSetsRequest))
          as ListSqlInjectionMatchSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSqlInjectionMatchSetsRequest create() =>
      ListSqlInjectionMatchSetsRequest._();
  @$core.override
  ListSqlInjectionMatchSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSqlInjectionMatchSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSqlInjectionMatchSetsRequest>(
          create);
  static ListSqlInjectionMatchSetsRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListSqlInjectionMatchSetsResponse extends $pb.GeneratedMessage {
  factory ListSqlInjectionMatchSetsResponse({
    $core.Iterable<SqlInjectionMatchSetSummary>? sqlinjectionmatchsets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (sqlinjectionmatchsets != null)
      result.sqlinjectionmatchsets.addAll(sqlinjectionmatchsets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListSqlInjectionMatchSetsResponse._();

  factory ListSqlInjectionMatchSetsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSqlInjectionMatchSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSqlInjectionMatchSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<SqlInjectionMatchSetSummary>(
        30262345, _omitFieldNames ? '' : 'sqlinjectionmatchsets',
        subBuilder: SqlInjectionMatchSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSqlInjectionMatchSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSqlInjectionMatchSetsResponse copyWith(
          void Function(ListSqlInjectionMatchSetsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListSqlInjectionMatchSetsResponse))
          as ListSqlInjectionMatchSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSqlInjectionMatchSetsResponse create() =>
      ListSqlInjectionMatchSetsResponse._();
  @$core.override
  ListSqlInjectionMatchSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSqlInjectionMatchSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSqlInjectionMatchSetsResponse>(
          create);
  static ListSqlInjectionMatchSetsResponse? _defaultInstance;

  @$pb.TagNumber(30262345)
  $pb.PbList<SqlInjectionMatchSetSummary> get sqlinjectionmatchsets =>
      $_getList(0);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListSubscribedRuleGroupsRequest extends $pb.GeneratedMessage {
  factory ListSubscribedRuleGroupsRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListSubscribedRuleGroupsRequest._();

  factory ListSubscribedRuleGroupsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSubscribedRuleGroupsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSubscribedRuleGroupsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscribedRuleGroupsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscribedRuleGroupsRequest copyWith(
          void Function(ListSubscribedRuleGroupsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListSubscribedRuleGroupsRequest))
          as ListSubscribedRuleGroupsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSubscribedRuleGroupsRequest create() =>
      ListSubscribedRuleGroupsRequest._();
  @$core.override
  ListSubscribedRuleGroupsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSubscribedRuleGroupsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSubscribedRuleGroupsRequest>(
          create);
  static ListSubscribedRuleGroupsRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListSubscribedRuleGroupsResponse extends $pb.GeneratedMessage {
  factory ListSubscribedRuleGroupsResponse({
    $core.Iterable<SubscribedRuleGroupSummary>? rulegroups,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (rulegroups != null) result.rulegroups.addAll(rulegroups);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListSubscribedRuleGroupsResponse._();

  factory ListSubscribedRuleGroupsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSubscribedRuleGroupsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSubscribedRuleGroupsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<SubscribedRuleGroupSummary>(
        270019740, _omitFieldNames ? '' : 'rulegroups',
        subBuilder: SubscribedRuleGroupSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscribedRuleGroupsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSubscribedRuleGroupsResponse copyWith(
          void Function(ListSubscribedRuleGroupsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListSubscribedRuleGroupsResponse))
          as ListSubscribedRuleGroupsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSubscribedRuleGroupsResponse create() =>
      ListSubscribedRuleGroupsResponse._();
  @$core.override
  ListSubscribedRuleGroupsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSubscribedRuleGroupsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSubscribedRuleGroupsResponse>(
          create);
  static ListSubscribedRuleGroupsResponse? _defaultInstance;

  @$pb.TagNumber(270019740)
  $pb.PbList<SubscribedRuleGroupSummary> get rulegroups => $_getList(0);

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
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

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class ListXssMatchSetsRequest extends $pb.GeneratedMessage {
  factory ListXssMatchSetsRequest({
    $core.int? limit,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (limit != null) result.limit = limit;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListXssMatchSetsRequest._();

  factory ListXssMatchSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListXssMatchSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListXssMatchSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListXssMatchSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListXssMatchSetsRequest copyWith(
          void Function(ListXssMatchSetsRequest) updates) =>
      super.copyWith((message) => updates(message as ListXssMatchSetsRequest))
          as ListXssMatchSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListXssMatchSetsRequest create() => ListXssMatchSetsRequest._();
  @$core.override
  ListXssMatchSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListXssMatchSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListXssMatchSetsRequest>(create);
  static ListXssMatchSetsRequest? _defaultInstance;

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(0);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(0);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(1);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(1);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListXssMatchSetsResponse extends $pb.GeneratedMessage {
  factory ListXssMatchSetsResponse({
    $core.Iterable<XssMatchSetSummary>? xssmatchsets,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (xssmatchsets != null) result.xssmatchsets.addAll(xssmatchsets);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListXssMatchSetsResponse._();

  factory ListXssMatchSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListXssMatchSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListXssMatchSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<XssMatchSetSummary>(500766384, _omitFieldNames ? '' : 'xssmatchsets',
        subBuilder: XssMatchSetSummary.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListXssMatchSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListXssMatchSetsResponse copyWith(
          void Function(ListXssMatchSetsResponse) updates) =>
      super.copyWith((message) => updates(message as ListXssMatchSetsResponse))
          as ListXssMatchSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListXssMatchSetsResponse create() => ListXssMatchSetsResponse._();
  @$core.override
  ListXssMatchSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListXssMatchSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListXssMatchSetsResponse>(create);
  static ListXssMatchSetsResponse? _defaultInstance;

  @$pb.TagNumber(500766384)
  $pb.PbList<XssMatchSetSummary> get xssmatchsets => $_getList(0);

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
    $core.Iterable<FieldToMatch>? redactedfields,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (logdestinationconfigs != null)
      result.logdestinationconfigs.addAll(logdestinationconfigs);
    if (redactedfields != null) result.redactedfields.addAll(redactedfields);
    if (resourcearn != null) result.resourcearn = resourcearn;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPS(22070207, _omitFieldNames ? '' : 'logdestinationconfigs')
    ..pPM<FieldToMatch>(48530737, _omitFieldNames ? '' : 'redactedfields',
        subBuilder: FieldToMatch.create)
    ..aOS(364280877, _omitFieldNames ? '' : 'resourcearn')
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

  @$pb.TagNumber(48530737)
  $pb.PbList<FieldToMatch> get redactedfields => $_getList(1);

  @$pb.TagNumber(364280877)
  $core.String get resourcearn => $_getSZ(2);
  @$pb.TagNumber(364280877)
  set resourcearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(364280877)
  $core.bool hasResourcearn() => $_has(2);
  @$pb.TagNumber(364280877)
  void clearResourcearn() => $_clearField(364280877);
}

class Predicate extends $pb.GeneratedMessage {
  factory Predicate({
    $core.bool? negated,
    PredicateType? type,
    $core.String? dataid,
  }) {
    final result = create();
    if (negated != null) result.negated = negated;
    if (type != null) result.type = type;
    if (dataid != null) result.dataid = dataid;
    return result;
  }

  Predicate._();

  factory Predicate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Predicate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Predicate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOB(204559310, _omitFieldNames ? '' : 'negated')
    ..aE<PredicateType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: PredicateType.values)
    ..aOS(500692225, _omitFieldNames ? '' : 'dataid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Predicate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Predicate copyWith(void Function(Predicate) updates) =>
      super.copyWith((message) => updates(message as Predicate)) as Predicate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Predicate create() => Predicate._();
  @$core.override
  Predicate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Predicate getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Predicate>(create);
  static Predicate? _defaultInstance;

  @$pb.TagNumber(204559310)
  $core.bool get negated => $_getBF(0);
  @$pb.TagNumber(204559310)
  set negated($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(204559310)
  $core.bool hasNegated() => $_has(0);
  @$pb.TagNumber(204559310)
  void clearNegated() => $_clearField(204559310);

  @$pb.TagNumber(290836590)
  PredicateType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(PredicateType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(500692225)
  $core.String get dataid => $_getSZ(2);
  @$pb.TagNumber(500692225)
  set dataid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(500692225)
  $core.bool hasDataid() => $_has(2);
  @$pb.TagNumber(500692225)
  void clearDataid() => $_clearField(500692225);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class RateBasedRule extends $pb.GeneratedMessage {
  factory RateBasedRule({
    RateKey? ratekey,
    $fixnum.Int64? ratelimit,
    $core.String? metricname,
    $core.String? name,
    $core.String? ruleid,
    $core.Iterable<Predicate>? matchpredicates,
  }) {
    final result = create();
    if (ratekey != null) result.ratekey = ratekey;
    if (ratelimit != null) result.ratelimit = ratelimit;
    if (metricname != null) result.metricname = metricname;
    if (name != null) result.name = name;
    if (ruleid != null) result.ruleid = ruleid;
    if (matchpredicates != null) result.matchpredicates.addAll(matchpredicates);
    return result;
  }

  RateBasedRule._();

  factory RateBasedRule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RateBasedRule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RateBasedRule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<RateKey>(31421727, _omitFieldNames ? '' : 'ratekey',
        enumValues: RateKey.values)
    ..aInt64(70042947, _omitFieldNames ? '' : 'ratelimit')
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..pPM<Predicate>(480711449, _omitFieldNames ? '' : 'matchpredicates',
        subBuilder: Predicate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateBasedRule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RateBasedRule copyWith(void Function(RateBasedRule) updates) =>
      super.copyWith((message) => updates(message as RateBasedRule))
          as RateBasedRule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RateBasedRule create() => RateBasedRule._();
  @$core.override
  RateBasedRule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RateBasedRule getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RateBasedRule>(create);
  static RateBasedRule? _defaultInstance;

  @$pb.TagNumber(31421727)
  RateKey get ratekey => $_getN(0);
  @$pb.TagNumber(31421727)
  set ratekey(RateKey value) => $_setField(31421727, value);
  @$pb.TagNumber(31421727)
  $core.bool hasRatekey() => $_has(0);
  @$pb.TagNumber(31421727)
  void clearRatekey() => $_clearField(31421727);

  @$pb.TagNumber(70042947)
  $fixnum.Int64 get ratelimit => $_getI64(1);
  @$pb.TagNumber(70042947)
  set ratelimit($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(70042947)
  $core.bool hasRatelimit() => $_has(1);
  @$pb.TagNumber(70042947)
  void clearRatelimit() => $_clearField(70042947);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(2);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(2);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(4);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(4);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);

  @$pb.TagNumber(480711449)
  $pb.PbList<Predicate> get matchpredicates => $_getList(5);
}

class RegexMatchSet extends $pb.GeneratedMessage {
  factory RegexMatchSet({
    $core.String? regexmatchsetid,
    $core.Iterable<RegexMatchTuple>? regexmatchtuples,
    $core.String? name,
  }) {
    final result = create();
    if (regexmatchsetid != null) result.regexmatchsetid = regexmatchsetid;
    if (regexmatchtuples != null)
      result.regexmatchtuples.addAll(regexmatchtuples);
    if (name != null) result.name = name;
    return result;
  }

  RegexMatchSet._();

  factory RegexMatchSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexMatchSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexMatchSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82287635, _omitFieldNames ? '' : 'regexmatchsetid')
    ..pPM<RegexMatchTuple>(252918587, _omitFieldNames ? '' : 'regexmatchtuples',
        subBuilder: RegexMatchTuple.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchSet copyWith(void Function(RegexMatchSet) updates) =>
      super.copyWith((message) => updates(message as RegexMatchSet))
          as RegexMatchSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexMatchSet create() => RegexMatchSet._();
  @$core.override
  RegexMatchSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexMatchSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexMatchSet>(create);
  static RegexMatchSet? _defaultInstance;

  @$pb.TagNumber(82287635)
  $core.String get regexmatchsetid => $_getSZ(0);
  @$pb.TagNumber(82287635)
  set regexmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82287635)
  $core.bool hasRegexmatchsetid() => $_has(0);
  @$pb.TagNumber(82287635)
  void clearRegexmatchsetid() => $_clearField(82287635);

  @$pb.TagNumber(252918587)
  $pb.PbList<RegexMatchTuple> get regexmatchtuples => $_getList(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class RegexMatchSetSummary extends $pb.GeneratedMessage {
  factory RegexMatchSetSummary({
    $core.String? regexmatchsetid,
    $core.String? name,
  }) {
    final result = create();
    if (regexmatchsetid != null) result.regexmatchsetid = regexmatchsetid;
    if (name != null) result.name = name;
    return result;
  }

  RegexMatchSetSummary._();

  factory RegexMatchSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexMatchSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexMatchSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82287635, _omitFieldNames ? '' : 'regexmatchsetid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchSetSummary copyWith(void Function(RegexMatchSetSummary) updates) =>
      super.copyWith((message) => updates(message as RegexMatchSetSummary))
          as RegexMatchSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexMatchSetSummary create() => RegexMatchSetSummary._();
  @$core.override
  RegexMatchSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexMatchSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexMatchSetSummary>(create);
  static RegexMatchSetSummary? _defaultInstance;

  @$pb.TagNumber(82287635)
  $core.String get regexmatchsetid => $_getSZ(0);
  @$pb.TagNumber(82287635)
  set regexmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82287635)
  $core.bool hasRegexmatchsetid() => $_has(0);
  @$pb.TagNumber(82287635)
  void clearRegexmatchsetid() => $_clearField(82287635);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class RegexMatchSetUpdate extends $pb.GeneratedMessage {
  factory RegexMatchSetUpdate({
    ChangeAction? action,
    RegexMatchTuple? regexmatchtuple,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (regexmatchtuple != null) result.regexmatchtuple = regexmatchtuple;
    return result;
  }

  RegexMatchSetUpdate._();

  factory RegexMatchSetUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexMatchSetUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexMatchSetUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..aOM<RegexMatchTuple>(364397678, _omitFieldNames ? '' : 'regexmatchtuple',
        subBuilder: RegexMatchTuple.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchSetUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchSetUpdate copyWith(void Function(RegexMatchSetUpdate) updates) =>
      super.copyWith((message) => updates(message as RegexMatchSetUpdate))
          as RegexMatchSetUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexMatchSetUpdate create() => RegexMatchSetUpdate._();
  @$core.override
  RegexMatchSetUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexMatchSetUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexMatchSetUpdate>(create);
  static RegexMatchSetUpdate? _defaultInstance;

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(364397678)
  RegexMatchTuple get regexmatchtuple => $_getN(1);
  @$pb.TagNumber(364397678)
  set regexmatchtuple(RegexMatchTuple value) => $_setField(364397678, value);
  @$pb.TagNumber(364397678)
  $core.bool hasRegexmatchtuple() => $_has(1);
  @$pb.TagNumber(364397678)
  void clearRegexmatchtuple() => $_clearField(364397678);
  @$pb.TagNumber(364397678)
  RegexMatchTuple ensureRegexmatchtuple() => $_ensure(1);
}

class RegexMatchTuple extends $pb.GeneratedMessage {
  factory RegexMatchTuple({
    $core.String? regexpatternsetid,
    FieldToMatch? fieldtomatch,
    TextTransformation? texttransformation,
  }) {
    final result = create();
    if (regexpatternsetid != null) result.regexpatternsetid = regexpatternsetid;
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    if (texttransformation != null)
      result.texttransformation = texttransformation;
    return result;
  }

  RegexMatchTuple._();

  factory RegexMatchTuple.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexMatchTuple.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexMatchTuple',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(96220912, _omitFieldNames ? '' : 'regexpatternsetid')
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..aE<TextTransformation>(
        375493608, _omitFieldNames ? '' : 'texttransformation',
        enumValues: TextTransformation.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchTuple clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexMatchTuple copyWith(void Function(RegexMatchTuple) updates) =>
      super.copyWith((message) => updates(message as RegexMatchTuple))
          as RegexMatchTuple;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexMatchTuple create() => RegexMatchTuple._();
  @$core.override
  RegexMatchTuple createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexMatchTuple getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexMatchTuple>(create);
  static RegexMatchTuple? _defaultInstance;

  @$pb.TagNumber(96220912)
  $core.String get regexpatternsetid => $_getSZ(0);
  @$pb.TagNumber(96220912)
  set regexpatternsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(96220912)
  $core.bool hasRegexpatternsetid() => $_has(0);
  @$pb.TagNumber(96220912)
  void clearRegexpatternsetid() => $_clearField(96220912);

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

  @$pb.TagNumber(375493608)
  TextTransformation get texttransformation => $_getN(2);
  @$pb.TagNumber(375493608)
  set texttransformation(TextTransformation value) =>
      $_setField(375493608, value);
  @$pb.TagNumber(375493608)
  $core.bool hasTexttransformation() => $_has(2);
  @$pb.TagNumber(375493608)
  void clearTexttransformation() => $_clearField(375493608);
}

class RegexPatternSet extends $pb.GeneratedMessage {
  factory RegexPatternSet({
    $core.String? regexpatternsetid,
    $core.String? name,
    $core.Iterable<$core.String>? regexpatternstrings,
  }) {
    final result = create();
    if (regexpatternsetid != null) result.regexpatternsetid = regexpatternsetid;
    if (name != null) result.name = name;
    if (regexpatternstrings != null)
      result.regexpatternstrings.addAll(regexpatternstrings);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(96220912, _omitFieldNames ? '' : 'regexpatternsetid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPS(488771049, _omitFieldNames ? '' : 'regexpatternstrings')
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

  @$pb.TagNumber(96220912)
  $core.String get regexpatternsetid => $_getSZ(0);
  @$pb.TagNumber(96220912)
  set regexpatternsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(96220912)
  $core.bool hasRegexpatternsetid() => $_has(0);
  @$pb.TagNumber(96220912)
  void clearRegexpatternsetid() => $_clearField(96220912);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(488771049)
  $pb.PbList<$core.String> get regexpatternstrings => $_getList(2);
}

class RegexPatternSetSummary extends $pb.GeneratedMessage {
  factory RegexPatternSetSummary({
    $core.String? regexpatternsetid,
    $core.String? name,
  }) {
    final result = create();
    if (regexpatternsetid != null) result.regexpatternsetid = regexpatternsetid;
    if (name != null) result.name = name;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(96220912, _omitFieldNames ? '' : 'regexpatternsetid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
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

  @$pb.TagNumber(96220912)
  $core.String get regexpatternsetid => $_getSZ(0);
  @$pb.TagNumber(96220912)
  set regexpatternsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(96220912)
  $core.bool hasRegexpatternsetid() => $_has(0);
  @$pb.TagNumber(96220912)
  void clearRegexpatternsetid() => $_clearField(96220912);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class RegexPatternSetUpdate extends $pb.GeneratedMessage {
  factory RegexPatternSetUpdate({
    ChangeAction? action,
    $core.String? regexpatternstring,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (regexpatternstring != null)
      result.regexpatternstring = regexpatternstring;
    return result;
  }

  RegexPatternSetUpdate._();

  factory RegexPatternSetUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RegexPatternSetUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RegexPatternSetUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..aOS(279972620, _omitFieldNames ? '' : 'regexpatternstring')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexPatternSetUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RegexPatternSetUpdate copyWith(
          void Function(RegexPatternSetUpdate) updates) =>
      super.copyWith((message) => updates(message as RegexPatternSetUpdate))
          as RegexPatternSetUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RegexPatternSetUpdate create() => RegexPatternSetUpdate._();
  @$core.override
  RegexPatternSetUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RegexPatternSetUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RegexPatternSetUpdate>(create);
  static RegexPatternSetUpdate? _defaultInstance;

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(279972620)
  $core.String get regexpatternstring => $_getSZ(1);
  @$pb.TagNumber(279972620)
  set regexpatternstring($core.String value) => $_setString(1, value);
  @$pb.TagNumber(279972620)
  $core.bool hasRegexpatternstring() => $_has(1);
  @$pb.TagNumber(279972620)
  void clearRegexpatternstring() => $_clearField(279972620);
}

class Rule extends $pb.GeneratedMessage {
  factory Rule({
    $core.String? metricname,
    $core.String? name,
    $core.Iterable<Predicate>? predicates,
    $core.String? ruleid,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (name != null) result.name = name;
    if (predicates != null) result.predicates.addAll(predicates);
    if (ruleid != null) result.ruleid = ruleid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Predicate>(315759288, _omitFieldNames ? '' : 'predicates',
        subBuilder: Predicate.create)
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
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

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(315759288)
  $pb.PbList<Predicate> get predicates => $_getList(2);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(3);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(3);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class RuleGroup extends $pb.GeneratedMessage {
  factory RuleGroup({
    $core.String? metricname,
    $core.String? name,
    $core.String? rulegroupid,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (name != null) result.name = name;
    if (rulegroupid != null) result.rulegroupid = rulegroupid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(287757320, _omitFieldNames ? '' : 'rulegroupid')
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

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(287757320)
  $core.String get rulegroupid => $_getSZ(2);
  @$pb.TagNumber(287757320)
  set rulegroupid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(287757320)
  $core.bool hasRulegroupid() => $_has(2);
  @$pb.TagNumber(287757320)
  void clearRulegroupid() => $_clearField(287757320);
}

class RuleGroupSummary extends $pb.GeneratedMessage {
  factory RuleGroupSummary({
    $core.String? name,
    $core.String? rulegroupid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (rulegroupid != null) result.rulegroupid = rulegroupid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(287757320, _omitFieldNames ? '' : 'rulegroupid')
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

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(287757320)
  $core.String get rulegroupid => $_getSZ(1);
  @$pb.TagNumber(287757320)
  set rulegroupid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(287757320)
  $core.bool hasRulegroupid() => $_has(1);
  @$pb.TagNumber(287757320)
  void clearRulegroupid() => $_clearField(287757320);
}

class RuleGroupUpdate extends $pb.GeneratedMessage {
  factory RuleGroupUpdate({
    ChangeAction? action,
    ActivatedRule? activatedrule,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (activatedrule != null) result.activatedrule = activatedrule;
    return result;
  }

  RuleGroupUpdate._();

  factory RuleGroupUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RuleGroupUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RuleGroupUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..aOM<ActivatedRule>(361834627, _omitFieldNames ? '' : 'activatedrule',
        subBuilder: ActivatedRule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleGroupUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleGroupUpdate copyWith(void Function(RuleGroupUpdate) updates) =>
      super.copyWith((message) => updates(message as RuleGroupUpdate))
          as RuleGroupUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RuleGroupUpdate create() => RuleGroupUpdate._();
  @$core.override
  RuleGroupUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RuleGroupUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RuleGroupUpdate>(create);
  static RuleGroupUpdate? _defaultInstance;

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(361834627)
  ActivatedRule get activatedrule => $_getN(1);
  @$pb.TagNumber(361834627)
  set activatedrule(ActivatedRule value) => $_setField(361834627, value);
  @$pb.TagNumber(361834627)
  $core.bool hasActivatedrule() => $_has(1);
  @$pb.TagNumber(361834627)
  void clearActivatedrule() => $_clearField(361834627);
  @$pb.TagNumber(361834627)
  ActivatedRule ensureActivatedrule() => $_ensure(1);
}

class RuleSummary extends $pb.GeneratedMessage {
  factory RuleSummary({
    $core.String? name,
    $core.String? ruleid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (ruleid != null) result.ruleid = ruleid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
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

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(1);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(1);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class RuleUpdate extends $pb.GeneratedMessage {
  factory RuleUpdate({
    ChangeAction? action,
    Predicate? predicate,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (predicate != null) result.predicate = predicate;
    return result;
  }

  RuleUpdate._();

  factory RuleUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RuleUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RuleUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..aOM<Predicate>(517130431, _omitFieldNames ? '' : 'predicate',
        subBuilder: Predicate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RuleUpdate copyWith(void Function(RuleUpdate) updates) =>
      super.copyWith((message) => updates(message as RuleUpdate)) as RuleUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RuleUpdate create() => RuleUpdate._();
  @$core.override
  RuleUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RuleUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RuleUpdate>(create);
  static RuleUpdate? _defaultInstance;

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(517130431)
  Predicate get predicate => $_getN(1);
  @$pb.TagNumber(517130431)
  set predicate(Predicate value) => $_setField(517130431, value);
  @$pb.TagNumber(517130431)
  $core.bool hasPredicate() => $_has(1);
  @$pb.TagNumber(517130431)
  void clearPredicate() => $_clearField(517130431);
  @$pb.TagNumber(517130431)
  Predicate ensurePredicate() => $_ensure(1);
}

class SampledHTTPRequest extends $pb.GeneratedMessage {
  factory SampledHTTPRequest({
    HTTPRequest? request,
    $core.String? rulewithinrulegroup,
    $core.String? timestamp,
    $core.String? action,
    $fixnum.Int64? weight,
  }) {
    final result = create();
    if (request != null) result.request = request;
    if (rulewithinrulegroup != null)
      result.rulewithinrulegroup = rulewithinrulegroup;
    if (timestamp != null) result.timestamp = timestamp;
    if (action != null) result.action = action;
    if (weight != null) result.weight = weight;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<HTTPRequest>(38093139, _omitFieldNames ? '' : 'request',
        subBuilder: HTTPRequest.create)
    ..aOS(58580650, _omitFieldNames ? '' : 'rulewithinrulegroup')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aOS(175614240, _omitFieldNames ? '' : 'action')
    ..aInt64(422581466, _omitFieldNames ? '' : 'weight')
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

  @$pb.TagNumber(38093139)
  HTTPRequest get request => $_getN(0);
  @$pb.TagNumber(38093139)
  set request(HTTPRequest value) => $_setField(38093139, value);
  @$pb.TagNumber(38093139)
  $core.bool hasRequest() => $_has(0);
  @$pb.TagNumber(38093139)
  void clearRequest() => $_clearField(38093139);
  @$pb.TagNumber(38093139)
  HTTPRequest ensureRequest() => $_ensure(0);

  @$pb.TagNumber(58580650)
  $core.String get rulewithinrulegroup => $_getSZ(1);
  @$pb.TagNumber(58580650)
  set rulewithinrulegroup($core.String value) => $_setString(1, value);
  @$pb.TagNumber(58580650)
  $core.bool hasRulewithinrulegroup() => $_has(1);
  @$pb.TagNumber(58580650)
  void clearRulewithinrulegroup() => $_clearField(58580650);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(2);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(2);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(175614240)
  $core.String get action => $_getSZ(3);
  @$pb.TagNumber(175614240)
  set action($core.String value) => $_setString(3, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(3);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(422581466)
  $fixnum.Int64 get weight => $_getI64(4);
  @$pb.TagNumber(422581466)
  set weight($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(422581466)
  $core.bool hasWeight() => $_has(4);
  @$pb.TagNumber(422581466)
  void clearWeight() => $_clearField(422581466);
}

class SizeConstraint extends $pb.GeneratedMessage {
  factory SizeConstraint({
    ComparisonOperator? comparisonoperator,
    $fixnum.Int64? size,
    FieldToMatch? fieldtomatch,
    TextTransformation? texttransformation,
  }) {
    final result = create();
    if (comparisonoperator != null)
      result.comparisonoperator = comparisonoperator;
    if (size != null) result.size = size;
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    if (texttransformation != null)
      result.texttransformation = texttransformation;
    return result;
  }

  SizeConstraint._();

  factory SizeConstraint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SizeConstraint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SizeConstraint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ComparisonOperator>(
        7059861, _omitFieldNames ? '' : 'comparisonoperator',
        enumValues: ComparisonOperator.values)
    ..aInt64(105352829, _omitFieldNames ? '' : 'size')
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..aE<TextTransformation>(
        375493608, _omitFieldNames ? '' : 'texttransformation',
        enumValues: TextTransformation.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraint copyWith(void Function(SizeConstraint) updates) =>
      super.copyWith((message) => updates(message as SizeConstraint))
          as SizeConstraint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SizeConstraint create() => SizeConstraint._();
  @$core.override
  SizeConstraint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SizeConstraint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SizeConstraint>(create);
  static SizeConstraint? _defaultInstance;

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

  @$pb.TagNumber(375493608)
  TextTransformation get texttransformation => $_getN(3);
  @$pb.TagNumber(375493608)
  set texttransformation(TextTransformation value) =>
      $_setField(375493608, value);
  @$pb.TagNumber(375493608)
  $core.bool hasTexttransformation() => $_has(3);
  @$pb.TagNumber(375493608)
  void clearTexttransformation() => $_clearField(375493608);
}

class SizeConstraintSet extends $pb.GeneratedMessage {
  factory SizeConstraintSet({
    $core.String? name,
    $core.String? sizeconstraintsetid,
    $core.Iterable<SizeConstraint>? sizeconstraints,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (sizeconstraintsetid != null)
      result.sizeconstraintsetid = sizeconstraintsetid;
    if (sizeconstraints != null) result.sizeconstraints.addAll(sizeconstraints);
    return result;
  }

  SizeConstraintSet._();

  factory SizeConstraintSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SizeConstraintSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SizeConstraintSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(277959757, _omitFieldNames ? '' : 'sizeconstraintsetid')
    ..pPM<SizeConstraint>(498335321, _omitFieldNames ? '' : 'sizeconstraints',
        subBuilder: SizeConstraint.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraintSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraintSet copyWith(void Function(SizeConstraintSet) updates) =>
      super.copyWith((message) => updates(message as SizeConstraintSet))
          as SizeConstraintSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SizeConstraintSet create() => SizeConstraintSet._();
  @$core.override
  SizeConstraintSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SizeConstraintSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SizeConstraintSet>(create);
  static SizeConstraintSet? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(277959757)
  $core.String get sizeconstraintsetid => $_getSZ(1);
  @$pb.TagNumber(277959757)
  set sizeconstraintsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(277959757)
  $core.bool hasSizeconstraintsetid() => $_has(1);
  @$pb.TagNumber(277959757)
  void clearSizeconstraintsetid() => $_clearField(277959757);

  @$pb.TagNumber(498335321)
  $pb.PbList<SizeConstraint> get sizeconstraints => $_getList(2);
}

class SizeConstraintSetSummary extends $pb.GeneratedMessage {
  factory SizeConstraintSetSummary({
    $core.String? name,
    $core.String? sizeconstraintsetid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (sizeconstraintsetid != null)
      result.sizeconstraintsetid = sizeconstraintsetid;
    return result;
  }

  SizeConstraintSetSummary._();

  factory SizeConstraintSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SizeConstraintSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SizeConstraintSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(277959757, _omitFieldNames ? '' : 'sizeconstraintsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraintSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraintSetSummary copyWith(
          void Function(SizeConstraintSetSummary) updates) =>
      super.copyWith((message) => updates(message as SizeConstraintSetSummary))
          as SizeConstraintSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SizeConstraintSetSummary create() => SizeConstraintSetSummary._();
  @$core.override
  SizeConstraintSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SizeConstraintSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SizeConstraintSetSummary>(create);
  static SizeConstraintSetSummary? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(277959757)
  $core.String get sizeconstraintsetid => $_getSZ(1);
  @$pb.TagNumber(277959757)
  set sizeconstraintsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(277959757)
  $core.bool hasSizeconstraintsetid() => $_has(1);
  @$pb.TagNumber(277959757)
  void clearSizeconstraintsetid() => $_clearField(277959757);
}

class SizeConstraintSetUpdate extends $pb.GeneratedMessage {
  factory SizeConstraintSetUpdate({
    SizeConstraint? sizeconstraint,
    ChangeAction? action,
  }) {
    final result = create();
    if (sizeconstraint != null) result.sizeconstraint = sizeconstraint;
    if (action != null) result.action = action;
    return result;
  }

  SizeConstraintSetUpdate._();

  factory SizeConstraintSetUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SizeConstraintSetUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SizeConstraintSetUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<SizeConstraint>(19554108, _omitFieldNames ? '' : 'sizeconstraint',
        subBuilder: SizeConstraint.create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraintSetUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SizeConstraintSetUpdate copyWith(
          void Function(SizeConstraintSetUpdate) updates) =>
      super.copyWith((message) => updates(message as SizeConstraintSetUpdate))
          as SizeConstraintSetUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SizeConstraintSetUpdate create() => SizeConstraintSetUpdate._();
  @$core.override
  SizeConstraintSetUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SizeConstraintSetUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SizeConstraintSetUpdate>(create);
  static SizeConstraintSetUpdate? _defaultInstance;

  @$pb.TagNumber(19554108)
  SizeConstraint get sizeconstraint => $_getN(0);
  @$pb.TagNumber(19554108)
  set sizeconstraint(SizeConstraint value) => $_setField(19554108, value);
  @$pb.TagNumber(19554108)
  $core.bool hasSizeconstraint() => $_has(0);
  @$pb.TagNumber(19554108)
  void clearSizeconstraint() => $_clearField(19554108);
  @$pb.TagNumber(19554108)
  SizeConstraint ensureSizeconstraint() => $_ensure(0);

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(1);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(1);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
}

class SqlInjectionMatchSet extends $pb.GeneratedMessage {
  factory SqlInjectionMatchSet({
    $core.Iterable<SqlInjectionMatchTuple>? sqlinjectionmatchtuples,
    $core.String? name,
    $core.String? sqlinjectionmatchsetid,
  }) {
    final result = create();
    if (sqlinjectionmatchtuples != null)
      result.sqlinjectionmatchtuples.addAll(sqlinjectionmatchtuples);
    if (name != null) result.name = name;
    if (sqlinjectionmatchsetid != null)
      result.sqlinjectionmatchsetid = sqlinjectionmatchsetid;
    return result;
  }

  SqlInjectionMatchSet._();

  factory SqlInjectionMatchSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SqlInjectionMatchSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SqlInjectionMatchSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<SqlInjectionMatchTuple>(
        74495631, _omitFieldNames ? '' : 'sqlinjectionmatchtuples',
        subBuilder: SqlInjectionMatchTuple.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(386225735, _omitFieldNames ? '' : 'sqlinjectionmatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqlInjectionMatchSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqlInjectionMatchSet copyWith(void Function(SqlInjectionMatchSet) updates) =>
      super.copyWith((message) => updates(message as SqlInjectionMatchSet))
          as SqlInjectionMatchSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SqlInjectionMatchSet create() => SqlInjectionMatchSet._();
  @$core.override
  SqlInjectionMatchSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SqlInjectionMatchSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SqlInjectionMatchSet>(create);
  static SqlInjectionMatchSet? _defaultInstance;

  @$pb.TagNumber(74495631)
  $pb.PbList<SqlInjectionMatchTuple> get sqlinjectionmatchtuples =>
      $_getList(0);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(386225735)
  $core.String get sqlinjectionmatchsetid => $_getSZ(2);
  @$pb.TagNumber(386225735)
  set sqlinjectionmatchsetid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(386225735)
  $core.bool hasSqlinjectionmatchsetid() => $_has(2);
  @$pb.TagNumber(386225735)
  void clearSqlinjectionmatchsetid() => $_clearField(386225735);
}

class SqlInjectionMatchSetSummary extends $pb.GeneratedMessage {
  factory SqlInjectionMatchSetSummary({
    $core.String? name,
    $core.String? sqlinjectionmatchsetid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (sqlinjectionmatchsetid != null)
      result.sqlinjectionmatchsetid = sqlinjectionmatchsetid;
    return result;
  }

  SqlInjectionMatchSetSummary._();

  factory SqlInjectionMatchSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SqlInjectionMatchSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SqlInjectionMatchSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(386225735, _omitFieldNames ? '' : 'sqlinjectionmatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqlInjectionMatchSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqlInjectionMatchSetSummary copyWith(
          void Function(SqlInjectionMatchSetSummary) updates) =>
      super.copyWith(
              (message) => updates(message as SqlInjectionMatchSetSummary))
          as SqlInjectionMatchSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SqlInjectionMatchSetSummary create() =>
      SqlInjectionMatchSetSummary._();
  @$core.override
  SqlInjectionMatchSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SqlInjectionMatchSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SqlInjectionMatchSetSummary>(create);
  static SqlInjectionMatchSetSummary? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(386225735)
  $core.String get sqlinjectionmatchsetid => $_getSZ(1);
  @$pb.TagNumber(386225735)
  set sqlinjectionmatchsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(386225735)
  $core.bool hasSqlinjectionmatchsetid() => $_has(1);
  @$pb.TagNumber(386225735)
  void clearSqlinjectionmatchsetid() => $_clearField(386225735);
}

class SqlInjectionMatchSetUpdate extends $pb.GeneratedMessage {
  factory SqlInjectionMatchSetUpdate({
    SqlInjectionMatchTuple? sqlinjectionmatchtuple,
    ChangeAction? action,
  }) {
    final result = create();
    if (sqlinjectionmatchtuple != null)
      result.sqlinjectionmatchtuple = sqlinjectionmatchtuple;
    if (action != null) result.action = action;
    return result;
  }

  SqlInjectionMatchSetUpdate._();

  factory SqlInjectionMatchSetUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SqlInjectionMatchSetUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SqlInjectionMatchSetUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<SqlInjectionMatchTuple>(
        88358794, _omitFieldNames ? '' : 'sqlinjectionmatchtuple',
        subBuilder: SqlInjectionMatchTuple.create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqlInjectionMatchSetUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqlInjectionMatchSetUpdate copyWith(
          void Function(SqlInjectionMatchSetUpdate) updates) =>
      super.copyWith(
              (message) => updates(message as SqlInjectionMatchSetUpdate))
          as SqlInjectionMatchSetUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SqlInjectionMatchSetUpdate create() => SqlInjectionMatchSetUpdate._();
  @$core.override
  SqlInjectionMatchSetUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SqlInjectionMatchSetUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SqlInjectionMatchSetUpdate>(create);
  static SqlInjectionMatchSetUpdate? _defaultInstance;

  @$pb.TagNumber(88358794)
  SqlInjectionMatchTuple get sqlinjectionmatchtuple => $_getN(0);
  @$pb.TagNumber(88358794)
  set sqlinjectionmatchtuple(SqlInjectionMatchTuple value) =>
      $_setField(88358794, value);
  @$pb.TagNumber(88358794)
  $core.bool hasSqlinjectionmatchtuple() => $_has(0);
  @$pb.TagNumber(88358794)
  void clearSqlinjectionmatchtuple() => $_clearField(88358794);
  @$pb.TagNumber(88358794)
  SqlInjectionMatchTuple ensureSqlinjectionmatchtuple() => $_ensure(0);

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(1);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(1);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
}

class SqlInjectionMatchTuple extends $pb.GeneratedMessage {
  factory SqlInjectionMatchTuple({
    FieldToMatch? fieldtomatch,
    TextTransformation? texttransformation,
  }) {
    final result = create();
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    if (texttransformation != null)
      result.texttransformation = texttransformation;
    return result;
  }

  SqlInjectionMatchTuple._();

  factory SqlInjectionMatchTuple.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SqlInjectionMatchTuple.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SqlInjectionMatchTuple',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..aE<TextTransformation>(
        375493608, _omitFieldNames ? '' : 'texttransformation',
        enumValues: TextTransformation.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqlInjectionMatchTuple clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqlInjectionMatchTuple copyWith(
          void Function(SqlInjectionMatchTuple) updates) =>
      super.copyWith((message) => updates(message as SqlInjectionMatchTuple))
          as SqlInjectionMatchTuple;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SqlInjectionMatchTuple create() => SqlInjectionMatchTuple._();
  @$core.override
  SqlInjectionMatchTuple createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SqlInjectionMatchTuple getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SqlInjectionMatchTuple>(create);
  static SqlInjectionMatchTuple? _defaultInstance;

  @$pb.TagNumber(338348372)
  FieldToMatch get fieldtomatch => $_getN(0);
  @$pb.TagNumber(338348372)
  set fieldtomatch(FieldToMatch value) => $_setField(338348372, value);
  @$pb.TagNumber(338348372)
  $core.bool hasFieldtomatch() => $_has(0);
  @$pb.TagNumber(338348372)
  void clearFieldtomatch() => $_clearField(338348372);
  @$pb.TagNumber(338348372)
  FieldToMatch ensureFieldtomatch() => $_ensure(0);

  @$pb.TagNumber(375493608)
  TextTransformation get texttransformation => $_getN(1);
  @$pb.TagNumber(375493608)
  set texttransformation(TextTransformation value) =>
      $_setField(375493608, value);
  @$pb.TagNumber(375493608)
  $core.bool hasTexttransformation() => $_has(1);
  @$pb.TagNumber(375493608)
  void clearTexttransformation() => $_clearField(375493608);
}

class SubscribedRuleGroupSummary extends $pb.GeneratedMessage {
  factory SubscribedRuleGroupSummary({
    $core.String? metricname,
    $core.String? name,
    $core.String? rulegroupid,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (name != null) result.name = name;
    if (rulegroupid != null) result.rulegroupid = rulegroupid;
    return result;
  }

  SubscribedRuleGroupSummary._();

  factory SubscribedRuleGroupSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubscribedRuleGroupSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubscribedRuleGroupSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(287757320, _omitFieldNames ? '' : 'rulegroupid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribedRuleGroupSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubscribedRuleGroupSummary copyWith(
          void Function(SubscribedRuleGroupSummary) updates) =>
      super.copyWith(
              (message) => updates(message as SubscribedRuleGroupSummary))
          as SubscribedRuleGroupSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubscribedRuleGroupSummary create() => SubscribedRuleGroupSummary._();
  @$core.override
  SubscribedRuleGroupSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubscribedRuleGroupSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubscribedRuleGroupSummary>(create);
  static SubscribedRuleGroupSummary? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(287757320)
  $core.String get rulegroupid => $_getSZ(2);
  @$pb.TagNumber(287757320)
  set rulegroupid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(287757320)
  $core.bool hasRulegroupid() => $_has(2);
  @$pb.TagNumber(287757320)
  void clearRulegroupid() => $_clearField(287757320);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class UpdateByteMatchSetRequest extends $pb.GeneratedMessage {
  factory UpdateByteMatchSetRequest({
    $core.String? bytematchsetid,
    $core.String? changetoken,
    $core.Iterable<ByteMatchSetUpdate>? updates,
  }) {
    final result = create();
    if (bytematchsetid != null) result.bytematchsetid = bytematchsetid;
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    return result;
  }

  UpdateByteMatchSetRequest._();

  factory UpdateByteMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateByteMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateByteMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(62033398, _omitFieldNames ? '' : 'bytematchsetid')
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<ByteMatchSetUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: ByteMatchSetUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateByteMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateByteMatchSetRequest copyWith(
          void Function(UpdateByteMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateByteMatchSetRequest))
          as UpdateByteMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateByteMatchSetRequest create() => UpdateByteMatchSetRequest._();
  @$core.override
  UpdateByteMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateByteMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateByteMatchSetRequest>(create);
  static UpdateByteMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(62033398)
  $core.String get bytematchsetid => $_getSZ(0);
  @$pb.TagNumber(62033398)
  set bytematchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(62033398)
  $core.bool hasBytematchsetid() => $_has(0);
  @$pb.TagNumber(62033398)
  void clearBytematchsetid() => $_clearField(62033398);

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(1);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(1);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<ByteMatchSetUpdate> get updates => $_getList(2);
}

class UpdateByteMatchSetResponse extends $pb.GeneratedMessage {
  factory UpdateByteMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  UpdateByteMatchSetResponse._();

  factory UpdateByteMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateByteMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateByteMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateByteMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateByteMatchSetResponse copyWith(
          void Function(UpdateByteMatchSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateByteMatchSetResponse))
          as UpdateByteMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateByteMatchSetResponse create() => UpdateByteMatchSetResponse._();
  @$core.override
  UpdateByteMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateByteMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateByteMatchSetResponse>(create);
  static UpdateByteMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateGeoMatchSetRequest extends $pb.GeneratedMessage {
  factory UpdateGeoMatchSetRequest({
    $core.String? changetoken,
    $core.Iterable<GeoMatchSetUpdate>? updates,
    $core.String? geomatchsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    if (geomatchsetid != null) result.geomatchsetid = geomatchsetid;
    return result;
  }

  UpdateGeoMatchSetRequest._();

  factory UpdateGeoMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateGeoMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateGeoMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<GeoMatchSetUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: GeoMatchSetUpdate.create)
    ..aOS(514346837, _omitFieldNames ? '' : 'geomatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGeoMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGeoMatchSetRequest copyWith(
          void Function(UpdateGeoMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateGeoMatchSetRequest))
          as UpdateGeoMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateGeoMatchSetRequest create() => UpdateGeoMatchSetRequest._();
  @$core.override
  UpdateGeoMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateGeoMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateGeoMatchSetRequest>(create);
  static UpdateGeoMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<GeoMatchSetUpdate> get updates => $_getList(1);

  @$pb.TagNumber(514346837)
  $core.String get geomatchsetid => $_getSZ(2);
  @$pb.TagNumber(514346837)
  set geomatchsetid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(514346837)
  $core.bool hasGeomatchsetid() => $_has(2);
  @$pb.TagNumber(514346837)
  void clearGeomatchsetid() => $_clearField(514346837);
}

class UpdateGeoMatchSetResponse extends $pb.GeneratedMessage {
  factory UpdateGeoMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  UpdateGeoMatchSetResponse._();

  factory UpdateGeoMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateGeoMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateGeoMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGeoMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateGeoMatchSetResponse copyWith(
          void Function(UpdateGeoMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateGeoMatchSetResponse))
          as UpdateGeoMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateGeoMatchSetResponse create() => UpdateGeoMatchSetResponse._();
  @$core.override
  UpdateGeoMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateGeoMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateGeoMatchSetResponse>(create);
  static UpdateGeoMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateIPSetRequest extends $pb.GeneratedMessage {
  factory UpdateIPSetRequest({
    $core.String? changetoken,
    $core.Iterable<IPSetUpdate>? updates,
    $core.String? ipsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    if (ipsetid != null) result.ipsetid = ipsetid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<IPSetUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: IPSetUpdate.create)
    ..aOS(346763066, _omitFieldNames ? '' : 'ipsetid')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<IPSetUpdate> get updates => $_getList(1);

  @$pb.TagNumber(346763066)
  $core.String get ipsetid => $_getSZ(2);
  @$pb.TagNumber(346763066)
  set ipsetid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(346763066)
  $core.bool hasIpsetid() => $_has(2);
  @$pb.TagNumber(346763066)
  void clearIpsetid() => $_clearField(346763066);
}

class UpdateIPSetResponse extends $pb.GeneratedMessage {
  factory UpdateIPSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateRateBasedRuleRequest extends $pb.GeneratedMessage {
  factory UpdateRateBasedRuleRequest({
    $fixnum.Int64? ratelimit,
    $core.String? changetoken,
    $core.Iterable<RuleUpdate>? updates,
    $core.String? ruleid,
  }) {
    final result = create();
    if (ratelimit != null) result.ratelimit = ratelimit;
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    if (ruleid != null) result.ruleid = ruleid;
    return result;
  }

  UpdateRateBasedRuleRequest._();

  factory UpdateRateBasedRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRateBasedRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRateBasedRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aInt64(70042947, _omitFieldNames ? '' : 'ratelimit')
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<RuleUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: RuleUpdate.create)
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRateBasedRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRateBasedRuleRequest copyWith(
          void Function(UpdateRateBasedRuleRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateRateBasedRuleRequest))
          as UpdateRateBasedRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRateBasedRuleRequest create() => UpdateRateBasedRuleRequest._();
  @$core.override
  UpdateRateBasedRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRateBasedRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRateBasedRuleRequest>(create);
  static UpdateRateBasedRuleRequest? _defaultInstance;

  @$pb.TagNumber(70042947)
  $fixnum.Int64 get ratelimit => $_getI64(0);
  @$pb.TagNumber(70042947)
  set ratelimit($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(70042947)
  $core.bool hasRatelimit() => $_has(0);
  @$pb.TagNumber(70042947)
  void clearRatelimit() => $_clearField(70042947);

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(1);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(1);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<RuleUpdate> get updates => $_getList(2);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(3);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(3);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class UpdateRateBasedRuleResponse extends $pb.GeneratedMessage {
  factory UpdateRateBasedRuleResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  UpdateRateBasedRuleResponse._();

  factory UpdateRateBasedRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRateBasedRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRateBasedRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRateBasedRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRateBasedRuleResponse copyWith(
          void Function(UpdateRateBasedRuleResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateRateBasedRuleResponse))
          as UpdateRateBasedRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRateBasedRuleResponse create() =>
      UpdateRateBasedRuleResponse._();
  @$core.override
  UpdateRateBasedRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRateBasedRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRateBasedRuleResponse>(create);
  static UpdateRateBasedRuleResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateRegexMatchSetRequest extends $pb.GeneratedMessage {
  factory UpdateRegexMatchSetRequest({
    $core.String? changetoken,
    $core.String? regexmatchsetid,
    $core.Iterable<RegexMatchSetUpdate>? updates,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (regexmatchsetid != null) result.regexmatchsetid = regexmatchsetid;
    if (updates != null) result.updates.addAll(updates);
    return result;
  }

  UpdateRegexMatchSetRequest._();

  factory UpdateRegexMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRegexMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRegexMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(82287635, _omitFieldNames ? '' : 'regexmatchsetid')
    ..pPM<RegexMatchSetUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: RegexMatchSetUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRegexMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRegexMatchSetRequest copyWith(
          void Function(UpdateRegexMatchSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateRegexMatchSetRequest))
          as UpdateRegexMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRegexMatchSetRequest create() => UpdateRegexMatchSetRequest._();
  @$core.override
  UpdateRegexMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRegexMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRegexMatchSetRequest>(create);
  static UpdateRegexMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(82287635)
  $core.String get regexmatchsetid => $_getSZ(1);
  @$pb.TagNumber(82287635)
  set regexmatchsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82287635)
  $core.bool hasRegexmatchsetid() => $_has(1);
  @$pb.TagNumber(82287635)
  void clearRegexmatchsetid() => $_clearField(82287635);

  @$pb.TagNumber(137393574)
  $pb.PbList<RegexMatchSetUpdate> get updates => $_getList(2);
}

class UpdateRegexMatchSetResponse extends $pb.GeneratedMessage {
  factory UpdateRegexMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  UpdateRegexMatchSetResponse._();

  factory UpdateRegexMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRegexMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRegexMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRegexMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRegexMatchSetResponse copyWith(
          void Function(UpdateRegexMatchSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateRegexMatchSetResponse))
          as UpdateRegexMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRegexMatchSetResponse create() =>
      UpdateRegexMatchSetResponse._();
  @$core.override
  UpdateRegexMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRegexMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRegexMatchSetResponse>(create);
  static UpdateRegexMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateRegexPatternSetRequest extends $pb.GeneratedMessage {
  factory UpdateRegexPatternSetRequest({
    $core.String? changetoken,
    $core.String? regexpatternsetid,
    $core.Iterable<RegexPatternSetUpdate>? updates,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (regexpatternsetid != null) result.regexpatternsetid = regexpatternsetid;
    if (updates != null) result.updates.addAll(updates);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..aOS(96220912, _omitFieldNames ? '' : 'regexpatternsetid')
    ..pPM<RegexPatternSetUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: RegexPatternSetUpdate.create)
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(96220912)
  $core.String get regexpatternsetid => $_getSZ(1);
  @$pb.TagNumber(96220912)
  set regexpatternsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(96220912)
  $core.bool hasRegexpatternsetid() => $_has(1);
  @$pb.TagNumber(96220912)
  void clearRegexpatternsetid() => $_clearField(96220912);

  @$pb.TagNumber(137393574)
  $pb.PbList<RegexPatternSetUpdate> get updates => $_getList(2);
}

class UpdateRegexPatternSetResponse extends $pb.GeneratedMessage {
  factory UpdateRegexPatternSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateRuleGroupRequest extends $pb.GeneratedMessage {
  factory UpdateRuleGroupRequest({
    $core.String? changetoken,
    $core.Iterable<RuleGroupUpdate>? updates,
    $core.String? rulegroupid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    if (rulegroupid != null) result.rulegroupid = rulegroupid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<RuleGroupUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: RuleGroupUpdate.create)
    ..aOS(287757320, _omitFieldNames ? '' : 'rulegroupid')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<RuleGroupUpdate> get updates => $_getList(1);

  @$pb.TagNumber(287757320)
  $core.String get rulegroupid => $_getSZ(2);
  @$pb.TagNumber(287757320)
  set rulegroupid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(287757320)
  $core.bool hasRulegroupid() => $_has(2);
  @$pb.TagNumber(287757320)
  void clearRulegroupid() => $_clearField(287757320);
}

class UpdateRuleGroupResponse extends $pb.GeneratedMessage {
  factory UpdateRuleGroupResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateRuleRequest extends $pb.GeneratedMessage {
  factory UpdateRuleRequest({
    $core.String? changetoken,
    $core.Iterable<RuleUpdate>? updates,
    $core.String? ruleid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    if (ruleid != null) result.ruleid = ruleid;
    return result;
  }

  UpdateRuleRequest._();

  factory UpdateRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<RuleUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: RuleUpdate.create)
    ..aOS(430449567, _omitFieldNames ? '' : 'ruleid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRuleRequest copyWith(void Function(UpdateRuleRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateRuleRequest))
          as UpdateRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRuleRequest create() => UpdateRuleRequest._();
  @$core.override
  UpdateRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRuleRequest>(create);
  static UpdateRuleRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<RuleUpdate> get updates => $_getList(1);

  @$pb.TagNumber(430449567)
  $core.String get ruleid => $_getSZ(2);
  @$pb.TagNumber(430449567)
  set ruleid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(430449567)
  $core.bool hasRuleid() => $_has(2);
  @$pb.TagNumber(430449567)
  void clearRuleid() => $_clearField(430449567);
}

class UpdateRuleResponse extends $pb.GeneratedMessage {
  factory UpdateRuleResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  UpdateRuleResponse._();

  factory UpdateRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateRuleResponse copyWith(void Function(UpdateRuleResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateRuleResponse))
          as UpdateRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateRuleResponse create() => UpdateRuleResponse._();
  @$core.override
  UpdateRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateRuleResponse>(create);
  static UpdateRuleResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateSizeConstraintSetRequest extends $pb.GeneratedMessage {
  factory UpdateSizeConstraintSetRequest({
    $core.String? changetoken,
    $core.Iterable<SizeConstraintSetUpdate>? updates,
    $core.String? sizeconstraintsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    if (sizeconstraintsetid != null)
      result.sizeconstraintsetid = sizeconstraintsetid;
    return result;
  }

  UpdateSizeConstraintSetRequest._();

  factory UpdateSizeConstraintSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateSizeConstraintSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateSizeConstraintSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<SizeConstraintSetUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: SizeConstraintSetUpdate.create)
    ..aOS(277959757, _omitFieldNames ? '' : 'sizeconstraintsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSizeConstraintSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSizeConstraintSetRequest copyWith(
          void Function(UpdateSizeConstraintSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateSizeConstraintSetRequest))
          as UpdateSizeConstraintSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateSizeConstraintSetRequest create() =>
      UpdateSizeConstraintSetRequest._();
  @$core.override
  UpdateSizeConstraintSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateSizeConstraintSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateSizeConstraintSetRequest>(create);
  static UpdateSizeConstraintSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<SizeConstraintSetUpdate> get updates => $_getList(1);

  @$pb.TagNumber(277959757)
  $core.String get sizeconstraintsetid => $_getSZ(2);
  @$pb.TagNumber(277959757)
  set sizeconstraintsetid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(277959757)
  $core.bool hasSizeconstraintsetid() => $_has(2);
  @$pb.TagNumber(277959757)
  void clearSizeconstraintsetid() => $_clearField(277959757);
}

class UpdateSizeConstraintSetResponse extends $pb.GeneratedMessage {
  factory UpdateSizeConstraintSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  UpdateSizeConstraintSetResponse._();

  factory UpdateSizeConstraintSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateSizeConstraintSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateSizeConstraintSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSizeConstraintSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSizeConstraintSetResponse copyWith(
          void Function(UpdateSizeConstraintSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateSizeConstraintSetResponse))
          as UpdateSizeConstraintSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateSizeConstraintSetResponse create() =>
      UpdateSizeConstraintSetResponse._();
  @$core.override
  UpdateSizeConstraintSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateSizeConstraintSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateSizeConstraintSetResponse>(
          create);
  static UpdateSizeConstraintSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateSqlInjectionMatchSetRequest extends $pb.GeneratedMessage {
  factory UpdateSqlInjectionMatchSetRequest({
    $core.String? changetoken,
    $core.Iterable<SqlInjectionMatchSetUpdate>? updates,
    $core.String? sqlinjectionmatchsetid,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    if (sqlinjectionmatchsetid != null)
      result.sqlinjectionmatchsetid = sqlinjectionmatchsetid;
    return result;
  }

  UpdateSqlInjectionMatchSetRequest._();

  factory UpdateSqlInjectionMatchSetRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateSqlInjectionMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateSqlInjectionMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<SqlInjectionMatchSetUpdate>(
        137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: SqlInjectionMatchSetUpdate.create)
    ..aOS(386225735, _omitFieldNames ? '' : 'sqlinjectionmatchsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSqlInjectionMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSqlInjectionMatchSetRequest copyWith(
          void Function(UpdateSqlInjectionMatchSetRequest) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateSqlInjectionMatchSetRequest))
          as UpdateSqlInjectionMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateSqlInjectionMatchSetRequest create() =>
      UpdateSqlInjectionMatchSetRequest._();
  @$core.override
  UpdateSqlInjectionMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateSqlInjectionMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateSqlInjectionMatchSetRequest>(
          create);
  static UpdateSqlInjectionMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<SqlInjectionMatchSetUpdate> get updates => $_getList(1);

  @$pb.TagNumber(386225735)
  $core.String get sqlinjectionmatchsetid => $_getSZ(2);
  @$pb.TagNumber(386225735)
  set sqlinjectionmatchsetid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(386225735)
  $core.bool hasSqlinjectionmatchsetid() => $_has(2);
  @$pb.TagNumber(386225735)
  void clearSqlinjectionmatchsetid() => $_clearField(386225735);
}

class UpdateSqlInjectionMatchSetResponse extends $pb.GeneratedMessage {
  factory UpdateSqlInjectionMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  UpdateSqlInjectionMatchSetResponse._();

  factory UpdateSqlInjectionMatchSetResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateSqlInjectionMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateSqlInjectionMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSqlInjectionMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateSqlInjectionMatchSetResponse copyWith(
          void Function(UpdateSqlInjectionMatchSetResponse) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateSqlInjectionMatchSetResponse))
          as UpdateSqlInjectionMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateSqlInjectionMatchSetResponse create() =>
      UpdateSqlInjectionMatchSetResponse._();
  @$core.override
  UpdateSqlInjectionMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateSqlInjectionMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateSqlInjectionMatchSetResponse>(
          create);
  static UpdateSqlInjectionMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateWebACLRequest extends $pb.GeneratedMessage {
  factory UpdateWebACLRequest({
    $core.String? changetoken,
    $core.Iterable<WebACLUpdate>? updates,
    $core.String? webaclid,
    WafAction? defaultaction,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    if (webaclid != null) result.webaclid = webaclid;
    if (defaultaction != null) result.defaultaction = defaultaction;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<WebACLUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: WebACLUpdate.create)
    ..aOS(161274579, _omitFieldNames ? '' : 'webaclid')
    ..aOM<WafAction>(322663861, _omitFieldNames ? '' : 'defaultaction',
        subBuilder: WafAction.create)
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<WebACLUpdate> get updates => $_getList(1);

  @$pb.TagNumber(161274579)
  $core.String get webaclid => $_getSZ(2);
  @$pb.TagNumber(161274579)
  set webaclid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(161274579)
  $core.bool hasWebaclid() => $_has(2);
  @$pb.TagNumber(161274579)
  void clearWebaclid() => $_clearField(161274579);

  @$pb.TagNumber(322663861)
  WafAction get defaultaction => $_getN(3);
  @$pb.TagNumber(322663861)
  set defaultaction(WafAction value) => $_setField(322663861, value);
  @$pb.TagNumber(322663861)
  $core.bool hasDefaultaction() => $_has(3);
  @$pb.TagNumber(322663861)
  void clearDefaultaction() => $_clearField(322663861);
  @$pb.TagNumber(322663861)
  WafAction ensureDefaultaction() => $_ensure(3);
}

class UpdateWebACLResponse extends $pb.GeneratedMessage {
  factory UpdateWebACLResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
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

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class UpdateXssMatchSetRequest extends $pb.GeneratedMessage {
  factory UpdateXssMatchSetRequest({
    $core.String? xssmatchsetid,
    $core.String? changetoken,
    $core.Iterable<XssMatchSetUpdate>? updates,
  }) {
    final result = create();
    if (xssmatchsetid != null) result.xssmatchsetid = xssmatchsetid;
    if (changetoken != null) result.changetoken = changetoken;
    if (updates != null) result.updates.addAll(updates);
    return result;
  }

  UpdateXssMatchSetRequest._();

  factory UpdateXssMatchSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateXssMatchSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateXssMatchSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(56643900, _omitFieldNames ? '' : 'xssmatchsetid')
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..pPM<XssMatchSetUpdate>(137393574, _omitFieldNames ? '' : 'updates',
        subBuilder: XssMatchSetUpdate.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateXssMatchSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateXssMatchSetRequest copyWith(
          void Function(UpdateXssMatchSetRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateXssMatchSetRequest))
          as UpdateXssMatchSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateXssMatchSetRequest create() => UpdateXssMatchSetRequest._();
  @$core.override
  UpdateXssMatchSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateXssMatchSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateXssMatchSetRequest>(create);
  static UpdateXssMatchSetRequest? _defaultInstance;

  @$pb.TagNumber(56643900)
  $core.String get xssmatchsetid => $_getSZ(0);
  @$pb.TagNumber(56643900)
  set xssmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56643900)
  $core.bool hasXssmatchsetid() => $_has(0);
  @$pb.TagNumber(56643900)
  void clearXssmatchsetid() => $_clearField(56643900);

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(1);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(1);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);

  @$pb.TagNumber(137393574)
  $pb.PbList<XssMatchSetUpdate> get updates => $_getList(2);
}

class UpdateXssMatchSetResponse extends $pb.GeneratedMessage {
  factory UpdateXssMatchSetResponse({
    $core.String? changetoken,
  }) {
    final result = create();
    if (changetoken != null) result.changetoken = changetoken;
    return result;
  }

  UpdateXssMatchSetResponse._();

  factory UpdateXssMatchSetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateXssMatchSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateXssMatchSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(73590123, _omitFieldNames ? '' : 'changetoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateXssMatchSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateXssMatchSetResponse copyWith(
          void Function(UpdateXssMatchSetResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateXssMatchSetResponse))
          as UpdateXssMatchSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateXssMatchSetResponse create() => UpdateXssMatchSetResponse._();
  @$core.override
  UpdateXssMatchSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateXssMatchSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateXssMatchSetResponse>(create);
  static UpdateXssMatchSetResponse? _defaultInstance;

  @$pb.TagNumber(73590123)
  $core.String get changetoken => $_getSZ(0);
  @$pb.TagNumber(73590123)
  set changetoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73590123)
  $core.bool hasChangetoken() => $_has(0);
  @$pb.TagNumber(73590123)
  void clearChangetoken() => $_clearField(73590123);
}

class WAFBadRequestException extends $pb.GeneratedMessage {
  factory WAFBadRequestException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFBadRequestException._();

  factory WAFBadRequestException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFBadRequestException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFBadRequestException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFBadRequestException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFBadRequestException copyWith(
          void Function(WAFBadRequestException) updates) =>
      super.copyWith((message) => updates(message as WAFBadRequestException))
          as WAFBadRequestException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFBadRequestException create() => WAFBadRequestException._();
  @$core.override
  WAFBadRequestException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFBadRequestException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFBadRequestException>(create);
  static WAFBadRequestException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFDisallowedNameException extends $pb.GeneratedMessage {
  factory WAFDisallowedNameException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFDisallowedNameException._();

  factory WAFDisallowedNameException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFDisallowedNameException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFDisallowedNameException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFDisallowedNameException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFDisallowedNameException copyWith(
          void Function(WAFDisallowedNameException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFDisallowedNameException))
          as WAFDisallowedNameException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFDisallowedNameException create() => WAFDisallowedNameException._();
  @$core.override
  WAFDisallowedNameException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFDisallowedNameException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFDisallowedNameException>(create);
  static WAFDisallowedNameException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFEntityMigrationException extends $pb.GeneratedMessage {
  factory WAFEntityMigrationException({
    MigrationErrorType? migrationerrortype,
    $core.String? message,
    $core.String? migrationerrorreason,
  }) {
    final result = create();
    if (migrationerrortype != null)
      result.migrationerrortype = migrationerrortype;
    if (message != null) result.message = message;
    if (migrationerrorreason != null)
      result.migrationerrorreason = migrationerrorreason;
    return result;
  }

  WAFEntityMigrationException._();

  factory WAFEntityMigrationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFEntityMigrationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFEntityMigrationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<MigrationErrorType>(
        53702508, _omitFieldNames ? '' : 'migrationerrortype',
        enumValues: MigrationErrorType.values)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..aOS(351537368, _omitFieldNames ? '' : 'migrationerrorreason')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFEntityMigrationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFEntityMigrationException copyWith(
          void Function(WAFEntityMigrationException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFEntityMigrationException))
          as WAFEntityMigrationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFEntityMigrationException create() =>
      WAFEntityMigrationException._();
  @$core.override
  WAFEntityMigrationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFEntityMigrationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFEntityMigrationException>(create);
  static WAFEntityMigrationException? _defaultInstance;

  @$pb.TagNumber(53702508)
  MigrationErrorType get migrationerrortype => $_getN(0);
  @$pb.TagNumber(53702508)
  set migrationerrortype(MigrationErrorType value) =>
      $_setField(53702508, value);
  @$pb.TagNumber(53702508)
  $core.bool hasMigrationerrortype() => $_has(0);
  @$pb.TagNumber(53702508)
  void clearMigrationerrortype() => $_clearField(53702508);

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);

  @$pb.TagNumber(351537368)
  $core.String get migrationerrorreason => $_getSZ(2);
  @$pb.TagNumber(351537368)
  set migrationerrorreason($core.String value) => $_setString(2, value);
  @$pb.TagNumber(351537368)
  $core.bool hasMigrationerrorreason() => $_has(2);
  @$pb.TagNumber(351537368)
  void clearMigrationerrorreason() => $_clearField(351537368);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFInvalidAccountException extends $pb.GeneratedMessage {
  factory WAFInvalidAccountException() => create();

  WAFInvalidAccountException._();

  factory WAFInvalidAccountException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFInvalidAccountException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFInvalidAccountException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidAccountException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidAccountException copyWith(
          void Function(WAFInvalidAccountException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFInvalidAccountException))
          as WAFInvalidAccountException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFInvalidAccountException create() => WAFInvalidAccountException._();
  @$core.override
  WAFInvalidAccountException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFInvalidAccountException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFInvalidAccountException>(create);
  static WAFInvalidAccountException? _defaultInstance;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFInvalidParameterException extends $pb.GeneratedMessage {
  factory WAFInvalidParameterException({
    ParameterExceptionField? field_125985384,
    $core.String? parameter,
    ParameterExceptionReason? reason,
  }) {
    final result = create();
    if (field_125985384 != null) result.field_125985384 = field_125985384;
    if (parameter != null) result.parameter = parameter;
    if (reason != null) result.reason = reason;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ParameterExceptionField>(125985384, _omitFieldNames ? '' : 'field',
        enumValues: ParameterExceptionField.values)
    ..aOS(363921681, _omitFieldNames ? '' : 'parameter')
    ..aE<ParameterExceptionReason>(413359642, _omitFieldNames ? '' : 'reason',
        enumValues: ParameterExceptionReason.values)
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

  @$pb.TagNumber(125985384)
  ParameterExceptionField get field_125985384 => $_getN(0);
  @$pb.TagNumber(125985384)
  set field_125985384(ParameterExceptionField value) =>
      $_setField(125985384, value);
  @$pb.TagNumber(125985384)
  $core.bool hasField_125985384() => $_has(0);
  @$pb.TagNumber(125985384)
  void clearField_125985384() => $_clearField(125985384);

  @$pb.TagNumber(363921681)
  $core.String get parameter => $_getSZ(1);
  @$pb.TagNumber(363921681)
  set parameter($core.String value) => $_setString(1, value);
  @$pb.TagNumber(363921681)
  $core.bool hasParameter() => $_has(1);
  @$pb.TagNumber(363921681)
  void clearParameter() => $_clearField(363921681);

  @$pb.TagNumber(413359642)
  ParameterExceptionReason get reason => $_getN(2);
  @$pb.TagNumber(413359642)
  set reason(ParameterExceptionReason value) => $_setField(413359642, value);
  @$pb.TagNumber(413359642)
  $core.bool hasReason() => $_has(2);
  @$pb.TagNumber(413359642)
  void clearReason() => $_clearField(413359642);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFInvalidRegexPatternException extends $pb.GeneratedMessage {
  factory WAFInvalidRegexPatternException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFInvalidRegexPatternException._();

  factory WAFInvalidRegexPatternException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFInvalidRegexPatternException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFInvalidRegexPatternException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidRegexPatternException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFInvalidRegexPatternException copyWith(
          void Function(WAFInvalidRegexPatternException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFInvalidRegexPatternException))
          as WAFInvalidRegexPatternException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFInvalidRegexPatternException create() =>
      WAFInvalidRegexPatternException._();
  @$core.override
  WAFInvalidRegexPatternException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFInvalidRegexPatternException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFInvalidRegexPatternException>(
          create);
  static WAFInvalidRegexPatternException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFLimitsExceededException extends $pb.GeneratedMessage {
  factory WAFLimitsExceededException({
    $core.String? message,
  }) {
    final result = create();
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFNonEmptyEntityException extends $pb.GeneratedMessage {
  factory WAFNonEmptyEntityException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFNonEmptyEntityException._();

  factory WAFNonEmptyEntityException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFNonEmptyEntityException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFNonEmptyEntityException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFNonEmptyEntityException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFNonEmptyEntityException copyWith(
          void Function(WAFNonEmptyEntityException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFNonEmptyEntityException))
          as WAFNonEmptyEntityException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFNonEmptyEntityException create() => WAFNonEmptyEntityException._();
  @$core.override
  WAFNonEmptyEntityException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFNonEmptyEntityException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFNonEmptyEntityException>(create);
  static WAFNonEmptyEntityException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFNonexistentContainerException extends $pb.GeneratedMessage {
  factory WAFNonexistentContainerException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFNonexistentContainerException._();

  factory WAFNonexistentContainerException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFNonexistentContainerException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFNonexistentContainerException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFNonexistentContainerException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFNonexistentContainerException copyWith(
          void Function(WAFNonexistentContainerException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFNonexistentContainerException))
          as WAFNonexistentContainerException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFNonexistentContainerException create() =>
      WAFNonexistentContainerException._();
  @$core.override
  WAFNonexistentContainerException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFNonexistentContainerException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFNonexistentContainerException>(
          create);
  static WAFNonexistentContainerException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WAFReferencedItemException extends $pb.GeneratedMessage {
  factory WAFReferencedItemException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFReferencedItemException._();

  factory WAFReferencedItemException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFReferencedItemException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFReferencedItemException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFReferencedItemException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFReferencedItemException copyWith(
          void Function(WAFReferencedItemException) updates) =>
      super.copyWith(
              (message) => updates(message as WAFReferencedItemException))
          as WAFReferencedItemException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFReferencedItemException create() => WAFReferencedItemException._();
  @$core.override
  WAFReferencedItemException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFReferencedItemException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFReferencedItemException>(create);
  static WAFReferencedItemException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
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

class WAFStaleDataException extends $pb.GeneratedMessage {
  factory WAFStaleDataException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  WAFStaleDataException._();

  factory WAFStaleDataException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WAFStaleDataException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WAFStaleDataException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFStaleDataException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WAFStaleDataException copyWith(
          void Function(WAFStaleDataException) updates) =>
      super.copyWith((message) => updates(message as WAFStaleDataException))
          as WAFStaleDataException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WAFStaleDataException create() => WAFStaleDataException._();
  @$core.override
  WAFStaleDataException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WAFStaleDataException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WAFStaleDataException>(create);
  static WAFStaleDataException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class WafAction extends $pb.GeneratedMessage {
  factory WafAction({
    WafActionType? type,
  }) {
    final result = create();
    if (type != null) result.type = type;
    return result;
  }

  WafAction._();

  factory WafAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WafAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WafAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<WafActionType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: WafActionType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WafAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WafAction copyWith(void Function(WafAction) updates) =>
      super.copyWith((message) => updates(message as WafAction)) as WafAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WafAction create() => WafAction._();
  @$core.override
  WafAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WafAction getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<WafAction>(create);
  static WafAction? _defaultInstance;

  @$pb.TagNumber(290836590)
  WafActionType get type => $_getN(0);
  @$pb.TagNumber(290836590)
  set type(WafActionType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class WafOverrideAction extends $pb.GeneratedMessage {
  factory WafOverrideAction({
    WafOverrideActionType? type,
  }) {
    final result = create();
    if (type != null) result.type = type;
    return result;
  }

  WafOverrideAction._();

  factory WafOverrideAction.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WafOverrideAction.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WafOverrideAction',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<WafOverrideActionType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: WafOverrideActionType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WafOverrideAction clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WafOverrideAction copyWith(void Function(WafOverrideAction) updates) =>
      super.copyWith((message) => updates(message as WafOverrideAction))
          as WafOverrideAction;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WafOverrideAction create() => WafOverrideAction._();
  @$core.override
  WafOverrideAction createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WafOverrideAction getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WafOverrideAction>(create);
  static WafOverrideAction? _defaultInstance;

  @$pb.TagNumber(290836590)
  WafOverrideActionType get type => $_getN(0);
  @$pb.TagNumber(290836590)
  set type(WafOverrideActionType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class WebACL extends $pb.GeneratedMessage {
  factory WebACL({
    $core.Iterable<ActivatedRule>? rules,
    $core.String? webaclarn,
    $core.String? metricname,
    $core.String? webaclid,
    $core.String? name,
    WafAction? defaultaction,
  }) {
    final result = create();
    if (rules != null) result.rules.addAll(rules);
    if (webaclarn != null) result.webaclarn = webaclarn;
    if (metricname != null) result.metricname = metricname;
    if (webaclid != null) result.webaclid = webaclid;
    if (name != null) result.name = name;
    if (defaultaction != null) result.defaultaction = defaultaction;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..pPM<ActivatedRule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: ActivatedRule.create)
    ..aOS(82506659, _omitFieldNames ? '' : 'webaclarn')
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(161274579, _omitFieldNames ? '' : 'webaclid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<WafAction>(322663861, _omitFieldNames ? '' : 'defaultaction',
        subBuilder: WafAction.create)
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

  @$pb.TagNumber(42675585)
  $pb.PbList<ActivatedRule> get rules => $_getList(0);

  @$pb.TagNumber(82506659)
  $core.String get webaclarn => $_getSZ(1);
  @$pb.TagNumber(82506659)
  set webaclarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82506659)
  $core.bool hasWebaclarn() => $_has(1);
  @$pb.TagNumber(82506659)
  void clearWebaclarn() => $_clearField(82506659);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(2);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(2);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(161274579)
  $core.String get webaclid => $_getSZ(3);
  @$pb.TagNumber(161274579)
  set webaclid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(161274579)
  $core.bool hasWebaclid() => $_has(3);
  @$pb.TagNumber(161274579)
  void clearWebaclid() => $_clearField(161274579);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(322663861)
  WafAction get defaultaction => $_getN(5);
  @$pb.TagNumber(322663861)
  set defaultaction(WafAction value) => $_setField(322663861, value);
  @$pb.TagNumber(322663861)
  $core.bool hasDefaultaction() => $_has(5);
  @$pb.TagNumber(322663861)
  void clearDefaultaction() => $_clearField(322663861);
  @$pb.TagNumber(322663861)
  WafAction ensureDefaultaction() => $_ensure(5);
}

class WebACLSummary extends $pb.GeneratedMessage {
  factory WebACLSummary({
    $core.String? webaclid,
    $core.String? name,
  }) {
    final result = create();
    if (webaclid != null) result.webaclid = webaclid;
    if (name != null) result.name = name;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(161274579, _omitFieldNames ? '' : 'webaclid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
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

  @$pb.TagNumber(161274579)
  $core.String get webaclid => $_getSZ(0);
  @$pb.TagNumber(161274579)
  set webaclid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(161274579)
  $core.bool hasWebaclid() => $_has(0);
  @$pb.TagNumber(161274579)
  void clearWebaclid() => $_clearField(161274579);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class WebACLUpdate extends $pb.GeneratedMessage {
  factory WebACLUpdate({
    ChangeAction? action,
    ActivatedRule? activatedrule,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (activatedrule != null) result.activatedrule = activatedrule;
    return result;
  }

  WebACLUpdate._();

  factory WebACLUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory WebACLUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'WebACLUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..aOM<ActivatedRule>(361834627, _omitFieldNames ? '' : 'activatedrule',
        subBuilder: ActivatedRule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WebACLUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  WebACLUpdate copyWith(void Function(WebACLUpdate) updates) =>
      super.copyWith((message) => updates(message as WebACLUpdate))
          as WebACLUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static WebACLUpdate create() => WebACLUpdate._();
  @$core.override
  WebACLUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static WebACLUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<WebACLUpdate>(create);
  static WebACLUpdate? _defaultInstance;

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(361834627)
  ActivatedRule get activatedrule => $_getN(1);
  @$pb.TagNumber(361834627)
  set activatedrule(ActivatedRule value) => $_setField(361834627, value);
  @$pb.TagNumber(361834627)
  $core.bool hasActivatedrule() => $_has(1);
  @$pb.TagNumber(361834627)
  void clearActivatedrule() => $_clearField(361834627);
  @$pb.TagNumber(361834627)
  ActivatedRule ensureActivatedrule() => $_ensure(1);
}

class XssMatchSet extends $pb.GeneratedMessage {
  factory XssMatchSet({
    $core.String? xssmatchsetid,
    $core.String? name,
    $core.Iterable<XssMatchTuple>? xssmatchtuples,
  }) {
    final result = create();
    if (xssmatchsetid != null) result.xssmatchsetid = xssmatchsetid;
    if (name != null) result.name = name;
    if (xssmatchtuples != null) result.xssmatchtuples.addAll(xssmatchtuples);
    return result;
  }

  XssMatchSet._();

  factory XssMatchSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XssMatchSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XssMatchSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(56643900, _omitFieldNames ? '' : 'xssmatchsetid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<XssMatchTuple>(302316050, _omitFieldNames ? '' : 'xssmatchtuples',
        subBuilder: XssMatchTuple.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchSet copyWith(void Function(XssMatchSet) updates) =>
      super.copyWith((message) => updates(message as XssMatchSet))
          as XssMatchSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XssMatchSet create() => XssMatchSet._();
  @$core.override
  XssMatchSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XssMatchSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XssMatchSet>(create);
  static XssMatchSet? _defaultInstance;

  @$pb.TagNumber(56643900)
  $core.String get xssmatchsetid => $_getSZ(0);
  @$pb.TagNumber(56643900)
  set xssmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56643900)
  $core.bool hasXssmatchsetid() => $_has(0);
  @$pb.TagNumber(56643900)
  void clearXssmatchsetid() => $_clearField(56643900);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(302316050)
  $pb.PbList<XssMatchTuple> get xssmatchtuples => $_getList(2);
}

class XssMatchSetSummary extends $pb.GeneratedMessage {
  factory XssMatchSetSummary({
    $core.String? xssmatchsetid,
    $core.String? name,
  }) {
    final result = create();
    if (xssmatchsetid != null) result.xssmatchsetid = xssmatchsetid;
    if (name != null) result.name = name;
    return result;
  }

  XssMatchSetSummary._();

  factory XssMatchSetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XssMatchSetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XssMatchSetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOS(56643900, _omitFieldNames ? '' : 'xssmatchsetid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchSetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchSetSummary copyWith(void Function(XssMatchSetSummary) updates) =>
      super.copyWith((message) => updates(message as XssMatchSetSummary))
          as XssMatchSetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XssMatchSetSummary create() => XssMatchSetSummary._();
  @$core.override
  XssMatchSetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XssMatchSetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XssMatchSetSummary>(create);
  static XssMatchSetSummary? _defaultInstance;

  @$pb.TagNumber(56643900)
  $core.String get xssmatchsetid => $_getSZ(0);
  @$pb.TagNumber(56643900)
  set xssmatchsetid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56643900)
  $core.bool hasXssmatchsetid() => $_has(0);
  @$pb.TagNumber(56643900)
  void clearXssmatchsetid() => $_clearField(56643900);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class XssMatchSetUpdate extends $pb.GeneratedMessage {
  factory XssMatchSetUpdate({
    XssMatchTuple? xssmatchtuple,
    ChangeAction? action,
  }) {
    final result = create();
    if (xssmatchtuple != null) result.xssmatchtuple = xssmatchtuple;
    if (action != null) result.action = action;
    return result;
  }

  XssMatchSetUpdate._();

  factory XssMatchSetUpdate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XssMatchSetUpdate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XssMatchSetUpdate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<XssMatchTuple>(46918713, _omitFieldNames ? '' : 'xssmatchtuple',
        subBuilder: XssMatchTuple.create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchSetUpdate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchSetUpdate copyWith(void Function(XssMatchSetUpdate) updates) =>
      super.copyWith((message) => updates(message as XssMatchSetUpdate))
          as XssMatchSetUpdate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XssMatchSetUpdate create() => XssMatchSetUpdate._();
  @$core.override
  XssMatchSetUpdate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XssMatchSetUpdate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XssMatchSetUpdate>(create);
  static XssMatchSetUpdate? _defaultInstance;

  @$pb.TagNumber(46918713)
  XssMatchTuple get xssmatchtuple => $_getN(0);
  @$pb.TagNumber(46918713)
  set xssmatchtuple(XssMatchTuple value) => $_setField(46918713, value);
  @$pb.TagNumber(46918713)
  $core.bool hasXssmatchtuple() => $_has(0);
  @$pb.TagNumber(46918713)
  void clearXssmatchtuple() => $_clearField(46918713);
  @$pb.TagNumber(46918713)
  XssMatchTuple ensureXssmatchtuple() => $_ensure(0);

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(1);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(1);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);
}

class XssMatchTuple extends $pb.GeneratedMessage {
  factory XssMatchTuple({
    FieldToMatch? fieldtomatch,
    TextTransformation? texttransformation,
  }) {
    final result = create();
    if (fieldtomatch != null) result.fieldtomatch = fieldtomatch;
    if (texttransformation != null)
      result.texttransformation = texttransformation;
    return result;
  }

  XssMatchTuple._();

  factory XssMatchTuple.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory XssMatchTuple.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'XssMatchTuple',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'waf'),
      createEmptyInstance: create)
    ..aOM<FieldToMatch>(338348372, _omitFieldNames ? '' : 'fieldtomatch',
        subBuilder: FieldToMatch.create)
    ..aE<TextTransformation>(
        375493608, _omitFieldNames ? '' : 'texttransformation',
        enumValues: TextTransformation.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchTuple clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  XssMatchTuple copyWith(void Function(XssMatchTuple) updates) =>
      super.copyWith((message) => updates(message as XssMatchTuple))
          as XssMatchTuple;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static XssMatchTuple create() => XssMatchTuple._();
  @$core.override
  XssMatchTuple createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static XssMatchTuple getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<XssMatchTuple>(create);
  static XssMatchTuple? _defaultInstance;

  @$pb.TagNumber(338348372)
  FieldToMatch get fieldtomatch => $_getN(0);
  @$pb.TagNumber(338348372)
  set fieldtomatch(FieldToMatch value) => $_setField(338348372, value);
  @$pb.TagNumber(338348372)
  $core.bool hasFieldtomatch() => $_has(0);
  @$pb.TagNumber(338348372)
  void clearFieldtomatch() => $_clearField(338348372);
  @$pb.TagNumber(338348372)
  FieldToMatch ensureFieldtomatch() => $_ensure(0);

  @$pb.TagNumber(375493608)
  TextTransformation get texttransformation => $_getN(1);
  @$pb.TagNumber(375493608)
  set texttransformation(TextTransformation value) =>
      $_setField(375493608, value);
  @$pb.TagNumber(375493608)
  $core.bool hasTexttransformation() => $_has(1);
  @$pb.TagNumber(375493608)
  void clearTexttransformation() => $_clearField(375493608);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
