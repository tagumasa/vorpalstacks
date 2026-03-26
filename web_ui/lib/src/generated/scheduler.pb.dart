// This is a generated file - do not edit.
//
// Generated from scheduler.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class AwsVpcConfiguration extends $pb.GeneratedMessage {
  factory AwsVpcConfiguration({
    $core.Iterable<$core.String>? subnets,
    $core.String? assignpublicip,
    $core.Iterable<$core.String>? securitygroups,
  }) {
    final result = create();
    if (subnets != null) result.subnets.addAll(subnets);
    if (assignpublicip != null) result.assignpublicip = assignpublicip;
    if (securitygroups != null) result.securitygroups.addAll(securitygroups);
    return result;
  }

  AwsVpcConfiguration._();

  factory AwsVpcConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AwsVpcConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AwsVpcConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..pPS(414921506, _omitFieldNames ? '' : 'subnets')
    ..aOS(461653589, _omitFieldNames ? '' : 'assignpublicip')
    ..pPS(515282516, _omitFieldNames ? '' : 'securitygroups')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AwsVpcConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AwsVpcConfiguration copyWith(void Function(AwsVpcConfiguration) updates) =>
      super.copyWith((message) => updates(message as AwsVpcConfiguration))
          as AwsVpcConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AwsVpcConfiguration create() => AwsVpcConfiguration._();
  @$core.override
  AwsVpcConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AwsVpcConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AwsVpcConfiguration>(create);
  static AwsVpcConfiguration? _defaultInstance;

  @$pb.TagNumber(414921506)
  $pb.PbList<$core.String> get subnets => $_getList(0);

  @$pb.TagNumber(461653589)
  $core.String get assignpublicip => $_getSZ(1);
  @$pb.TagNumber(461653589)
  set assignpublicip($core.String value) => $_setString(1, value);
  @$pb.TagNumber(461653589)
  $core.bool hasAssignpublicip() => $_has(1);
  @$pb.TagNumber(461653589)
  void clearAssignpublicip() => $_clearField(461653589);

  @$pb.TagNumber(515282516)
  $pb.PbList<$core.String> get securitygroups => $_getList(2);
}

class CapacityProviderStrategyItem extends $pb.GeneratedMessage {
  factory CapacityProviderStrategyItem({
    $core.String? capacityprovider,
    $core.int? weight,
    $core.int? base,
  }) {
    final result = create();
    if (capacityprovider != null) result.capacityprovider = capacityprovider;
    if (weight != null) result.weight = weight;
    if (base != null) result.base = base;
    return result;
  }

  CapacityProviderStrategyItem._();

  factory CapacityProviderStrategyItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CapacityProviderStrategyItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CapacityProviderStrategyItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(109086449, _omitFieldNames ? '' : 'capacityprovider')
    ..aI(278961850, _omitFieldNames ? '' : 'weight')
    ..aI(500995289, _omitFieldNames ? '' : 'base')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityProviderStrategyItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CapacityProviderStrategyItem copyWith(
          void Function(CapacityProviderStrategyItem) updates) =>
      super.copyWith(
              (message) => updates(message as CapacityProviderStrategyItem))
          as CapacityProviderStrategyItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CapacityProviderStrategyItem create() =>
      CapacityProviderStrategyItem._();
  @$core.override
  CapacityProviderStrategyItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CapacityProviderStrategyItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CapacityProviderStrategyItem>(create);
  static CapacityProviderStrategyItem? _defaultInstance;

  @$pb.TagNumber(109086449)
  $core.String get capacityprovider => $_getSZ(0);
  @$pb.TagNumber(109086449)
  set capacityprovider($core.String value) => $_setString(0, value);
  @$pb.TagNumber(109086449)
  $core.bool hasCapacityprovider() => $_has(0);
  @$pb.TagNumber(109086449)
  void clearCapacityprovider() => $_clearField(109086449);

  @$pb.TagNumber(278961850)
  $core.int get weight => $_getIZ(1);
  @$pb.TagNumber(278961850)
  set weight($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(278961850)
  $core.bool hasWeight() => $_has(1);
  @$pb.TagNumber(278961850)
  void clearWeight() => $_clearField(278961850);

  @$pb.TagNumber(500995289)
  $core.int get base => $_getIZ(2);
  @$pb.TagNumber(500995289)
  set base($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(500995289)
  $core.bool hasBase() => $_has(2);
  @$pb.TagNumber(500995289)
  void clearBase() => $_clearField(500995289);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CreateScheduleGroupInput extends $pb.GeneratedMessage {
  factory CreateScheduleGroupInput({
    $core.String? clienttoken,
    $core.String? name,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  CreateScheduleGroupInput._();

  factory CreateScheduleGroupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateScheduleGroupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateScheduleGroupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduleGroupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduleGroupInput copyWith(
          void Function(CreateScheduleGroupInput) updates) =>
      super.copyWith((message) => updates(message as CreateScheduleGroupInput))
          as CreateScheduleGroupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateScheduleGroupInput create() => CreateScheduleGroupInput._();
  @$core.override
  CreateScheduleGroupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateScheduleGroupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateScheduleGroupInput>(create);
  static CreateScheduleGroupInput? _defaultInstance;

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(0);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(0);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

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
}

class CreateScheduleGroupOutput extends $pb.GeneratedMessage {
  factory CreateScheduleGroupOutput({
    $core.String? schedulegrouparn,
  }) {
    final result = create();
    if (schedulegrouparn != null) result.schedulegrouparn = schedulegrouparn;
    return result;
  }

  CreateScheduleGroupOutput._();

  factory CreateScheduleGroupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateScheduleGroupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateScheduleGroupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(148398979, _omitFieldNames ? '' : 'schedulegrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduleGroupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduleGroupOutput copyWith(
          void Function(CreateScheduleGroupOutput) updates) =>
      super.copyWith((message) => updates(message as CreateScheduleGroupOutput))
          as CreateScheduleGroupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateScheduleGroupOutput create() => CreateScheduleGroupOutput._();
  @$core.override
  CreateScheduleGroupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateScheduleGroupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateScheduleGroupOutput>(create);
  static CreateScheduleGroupOutput? _defaultInstance;

  @$pb.TagNumber(148398979)
  $core.String get schedulegrouparn => $_getSZ(0);
  @$pb.TagNumber(148398979)
  set schedulegrouparn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(148398979)
  $core.bool hasSchedulegrouparn() => $_has(0);
  @$pb.TagNumber(148398979)
  void clearSchedulegrouparn() => $_clearField(148398979);
}

class CreateScheduleInput extends $pb.GeneratedMessage {
  factory CreateScheduleInput({
    FlexibleTimeWindow? flexibletimewindow,
    $core.String? enddate,
    $core.String? kmskeyarn,
    $core.String? description,
    $core.String? clienttoken,
    Target? target,
    $core.String? name,
    $core.String? actionaftercompletion,
    $core.String? groupname,
    $core.String? scheduleexpressiontimezone,
    $core.String? startdate,
    $core.String? scheduleexpression,
    $core.String? state,
  }) {
    final result = create();
    if (flexibletimewindow != null)
      result.flexibletimewindow = flexibletimewindow;
    if (enddate != null) result.enddate = enddate;
    if (kmskeyarn != null) result.kmskeyarn = kmskeyarn;
    if (description != null) result.description = description;
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (target != null) result.target = target;
    if (name != null) result.name = name;
    if (actionaftercompletion != null)
      result.actionaftercompletion = actionaftercompletion;
    if (groupname != null) result.groupname = groupname;
    if (scheduleexpressiontimezone != null)
      result.scheduleexpressiontimezone = scheduleexpressiontimezone;
    if (startdate != null) result.startdate = startdate;
    if (scheduleexpression != null)
      result.scheduleexpression = scheduleexpression;
    if (state != null) result.state = state;
    return result;
  }

  CreateScheduleInput._();

  factory CreateScheduleInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateScheduleInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateScheduleInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOM<FlexibleTimeWindow>(
        2518952, _omitFieldNames ? '' : 'flexibletimewindow',
        subBuilder: FlexibleTimeWindow.create)
    ..aOS(77486543, _omitFieldNames ? '' : 'enddate')
    ..aOS(110041649, _omitFieldNames ? '' : 'kmskeyarn')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOM<Target>(191361385, _omitFieldNames ? '' : 'target',
        subBuilder: Target.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(282350906, _omitFieldNames ? '' : 'actionaftercompletion')
    ..aOS(357049672, _omitFieldNames ? '' : 'groupname')
    ..aOS(400730326, _omitFieldNames ? '' : 'scheduleexpressiontimezone')
    ..aOS(445135996, _omitFieldNames ? '' : 'startdate')
    ..aOS(446089471, _omitFieldNames ? '' : 'scheduleexpression')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduleInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduleInput copyWith(void Function(CreateScheduleInput) updates) =>
      super.copyWith((message) => updates(message as CreateScheduleInput))
          as CreateScheduleInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateScheduleInput create() => CreateScheduleInput._();
  @$core.override
  CreateScheduleInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateScheduleInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateScheduleInput>(create);
  static CreateScheduleInput? _defaultInstance;

  @$pb.TagNumber(2518952)
  FlexibleTimeWindow get flexibletimewindow => $_getN(0);
  @$pb.TagNumber(2518952)
  set flexibletimewindow(FlexibleTimeWindow value) =>
      $_setField(2518952, value);
  @$pb.TagNumber(2518952)
  $core.bool hasFlexibletimewindow() => $_has(0);
  @$pb.TagNumber(2518952)
  void clearFlexibletimewindow() => $_clearField(2518952);
  @$pb.TagNumber(2518952)
  FlexibleTimeWindow ensureFlexibletimewindow() => $_ensure(0);

  @$pb.TagNumber(77486543)
  $core.String get enddate => $_getSZ(1);
  @$pb.TagNumber(77486543)
  set enddate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(77486543)
  $core.bool hasEnddate() => $_has(1);
  @$pb.TagNumber(77486543)
  void clearEnddate() => $_clearField(77486543);

  @$pb.TagNumber(110041649)
  $core.String get kmskeyarn => $_getSZ(2);
  @$pb.TagNumber(110041649)
  set kmskeyarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(110041649)
  $core.bool hasKmskeyarn() => $_has(2);
  @$pb.TagNumber(110041649)
  void clearKmskeyarn() => $_clearField(110041649);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(4);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(4);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(191361385)
  Target get target => $_getN(5);
  @$pb.TagNumber(191361385)
  set target(Target value) => $_setField(191361385, value);
  @$pb.TagNumber(191361385)
  $core.bool hasTarget() => $_has(5);
  @$pb.TagNumber(191361385)
  void clearTarget() => $_clearField(191361385);
  @$pb.TagNumber(191361385)
  Target ensureTarget() => $_ensure(5);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(282350906)
  $core.String get actionaftercompletion => $_getSZ(7);
  @$pb.TagNumber(282350906)
  set actionaftercompletion($core.String value) => $_setString(7, value);
  @$pb.TagNumber(282350906)
  $core.bool hasActionaftercompletion() => $_has(7);
  @$pb.TagNumber(282350906)
  void clearActionaftercompletion() => $_clearField(282350906);

  @$pb.TagNumber(357049672)
  $core.String get groupname => $_getSZ(8);
  @$pb.TagNumber(357049672)
  set groupname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(357049672)
  $core.bool hasGroupname() => $_has(8);
  @$pb.TagNumber(357049672)
  void clearGroupname() => $_clearField(357049672);

  @$pb.TagNumber(400730326)
  $core.String get scheduleexpressiontimezone => $_getSZ(9);
  @$pb.TagNumber(400730326)
  set scheduleexpressiontimezone($core.String value) => $_setString(9, value);
  @$pb.TagNumber(400730326)
  $core.bool hasScheduleexpressiontimezone() => $_has(9);
  @$pb.TagNumber(400730326)
  void clearScheduleexpressiontimezone() => $_clearField(400730326);

  @$pb.TagNumber(445135996)
  $core.String get startdate => $_getSZ(10);
  @$pb.TagNumber(445135996)
  set startdate($core.String value) => $_setString(10, value);
  @$pb.TagNumber(445135996)
  $core.bool hasStartdate() => $_has(10);
  @$pb.TagNumber(445135996)
  void clearStartdate() => $_clearField(445135996);

  @$pb.TagNumber(446089471)
  $core.String get scheduleexpression => $_getSZ(11);
  @$pb.TagNumber(446089471)
  set scheduleexpression($core.String value) => $_setString(11, value);
  @$pb.TagNumber(446089471)
  $core.bool hasScheduleexpression() => $_has(11);
  @$pb.TagNumber(446089471)
  void clearScheduleexpression() => $_clearField(446089471);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(12);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(12, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(12);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class CreateScheduleOutput extends $pb.GeneratedMessage {
  factory CreateScheduleOutput({
    $core.String? schedulearn,
  }) {
    final result = create();
    if (schedulearn != null) result.schedulearn = schedulearn;
    return result;
  }

  CreateScheduleOutput._();

  factory CreateScheduleOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateScheduleOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateScheduleOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(178445188, _omitFieldNames ? '' : 'schedulearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduleOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateScheduleOutput copyWith(void Function(CreateScheduleOutput) updates) =>
      super.copyWith((message) => updates(message as CreateScheduleOutput))
          as CreateScheduleOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateScheduleOutput create() => CreateScheduleOutput._();
  @$core.override
  CreateScheduleOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateScheduleOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateScheduleOutput>(create);
  static CreateScheduleOutput? _defaultInstance;

  @$pb.TagNumber(178445188)
  $core.String get schedulearn => $_getSZ(0);
  @$pb.TagNumber(178445188)
  set schedulearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(178445188)
  $core.bool hasSchedulearn() => $_has(0);
  @$pb.TagNumber(178445188)
  void clearSchedulearn() => $_clearField(178445188);
}

class DeadLetterConfig extends $pb.GeneratedMessage {
  factory DeadLetterConfig({
    $core.String? arn,
  }) {
    final result = create();
    if (arn != null) result.arn = arn;
    return result;
  }

  DeadLetterConfig._();

  factory DeadLetterConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeadLetterConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeadLetterConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeadLetterConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeadLetterConfig copyWith(void Function(DeadLetterConfig) updates) =>
      super.copyWith((message) => updates(message as DeadLetterConfig))
          as DeadLetterConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeadLetterConfig create() => DeadLetterConfig._();
  @$core.override
  DeadLetterConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeadLetterConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeadLetterConfig>(create);
  static DeadLetterConfig? _defaultInstance;

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(0);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(0);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class DeleteScheduleGroupInput extends $pb.GeneratedMessage {
  factory DeleteScheduleGroupInput({
    $core.String? clienttoken,
    $core.String? name,
  }) {
    final result = create();
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (name != null) result.name = name;
    return result;
  }

  DeleteScheduleGroupInput._();

  factory DeleteScheduleGroupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteScheduleGroupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteScheduleGroupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduleGroupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduleGroupInput copyWith(
          void Function(DeleteScheduleGroupInput) updates) =>
      super.copyWith((message) => updates(message as DeleteScheduleGroupInput))
          as DeleteScheduleGroupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteScheduleGroupInput create() => DeleteScheduleGroupInput._();
  @$core.override
  DeleteScheduleGroupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteScheduleGroupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteScheduleGroupInput>(create);
  static DeleteScheduleGroupInput? _defaultInstance;

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(0);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(0);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DeleteScheduleGroupOutput extends $pb.GeneratedMessage {
  factory DeleteScheduleGroupOutput() => create();

  DeleteScheduleGroupOutput._();

  factory DeleteScheduleGroupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteScheduleGroupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteScheduleGroupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduleGroupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduleGroupOutput copyWith(
          void Function(DeleteScheduleGroupOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteScheduleGroupOutput))
          as DeleteScheduleGroupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteScheduleGroupOutput create() => DeleteScheduleGroupOutput._();
  @$core.override
  DeleteScheduleGroupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteScheduleGroupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteScheduleGroupOutput>(create);
  static DeleteScheduleGroupOutput? _defaultInstance;
}

class DeleteScheduleInput extends $pb.GeneratedMessage {
  factory DeleteScheduleInput({
    $core.String? clienttoken,
    $core.String? name,
    $core.String? groupname,
  }) {
    final result = create();
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (name != null) result.name = name;
    if (groupname != null) result.groupname = groupname;
    return result;
  }

  DeleteScheduleInput._();

  factory DeleteScheduleInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteScheduleInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteScheduleInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(357049672, _omitFieldNames ? '' : 'groupname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduleInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduleInput copyWith(void Function(DeleteScheduleInput) updates) =>
      super.copyWith((message) => updates(message as DeleteScheduleInput))
          as DeleteScheduleInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteScheduleInput create() => DeleteScheduleInput._();
  @$core.override
  DeleteScheduleInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteScheduleInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteScheduleInput>(create);
  static DeleteScheduleInput? _defaultInstance;

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(0);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(0);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(357049672)
  $core.String get groupname => $_getSZ(2);
  @$pb.TagNumber(357049672)
  set groupname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(357049672)
  $core.bool hasGroupname() => $_has(2);
  @$pb.TagNumber(357049672)
  void clearGroupname() => $_clearField(357049672);
}

class DeleteScheduleOutput extends $pb.GeneratedMessage {
  factory DeleteScheduleOutput() => create();

  DeleteScheduleOutput._();

  factory DeleteScheduleOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteScheduleOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteScheduleOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduleOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteScheduleOutput copyWith(void Function(DeleteScheduleOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteScheduleOutput))
          as DeleteScheduleOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteScheduleOutput create() => DeleteScheduleOutput._();
  @$core.override
  DeleteScheduleOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteScheduleOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteScheduleOutput>(create);
  static DeleteScheduleOutput? _defaultInstance;
}

class EcsParameters extends $pb.GeneratedMessage {
  factory EcsParameters({
    $core.Iterable<PlacementStrategy>? placementstrategy,
    $core.String? taskdefinitionarn,
    $core.String? group,
    $core.String? platformversion,
    $core.bool? enableecsmanagedtags,
    $core.String? launchtype,
    NetworkConfiguration? networkconfiguration,
    $core.Iterable<PlacementConstraint>? placementconstraints,
    $core.Iterable<CapacityProviderStrategyItem>? capacityproviderstrategy,
    $core.String? referenceid,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? tags,
    $core.int? taskcount,
    $core.String? propagatetags,
    $core.bool? enableexecutecommand,
  }) {
    final result = create();
    if (placementstrategy != null)
      result.placementstrategy.addAll(placementstrategy);
    if (taskdefinitionarn != null) result.taskdefinitionarn = taskdefinitionarn;
    if (group != null) result.group = group;
    if (platformversion != null) result.platformversion = platformversion;
    if (enableecsmanagedtags != null)
      result.enableecsmanagedtags = enableecsmanagedtags;
    if (launchtype != null) result.launchtype = launchtype;
    if (networkconfiguration != null)
      result.networkconfiguration = networkconfiguration;
    if (placementconstraints != null)
      result.placementconstraints.addAll(placementconstraints);
    if (capacityproviderstrategy != null)
      result.capacityproviderstrategy.addAll(capacityproviderstrategy);
    if (referenceid != null) result.referenceid = referenceid;
    if (tags != null) result.tags.addEntries(tags);
    if (taskcount != null) result.taskcount = taskcount;
    if (propagatetags != null) result.propagatetags = propagatetags;
    if (enableexecutecommand != null)
      result.enableexecutecommand = enableexecutecommand;
    return result;
  }

  EcsParameters._();

  factory EcsParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EcsParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EcsParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..pPM<PlacementStrategy>(
        25036678, _omitFieldNames ? '' : 'placementstrategy',
        subBuilder: PlacementStrategy.create)
    ..aOS(82234477, _omitFieldNames ? '' : 'taskdefinitionarn')
    ..aOS(91525165, _omitFieldNames ? '' : 'group')
    ..aOS(139924287, _omitFieldNames ? '' : 'platformversion')
    ..aOB(146161174, _omitFieldNames ? '' : 'enableecsmanagedtags')
    ..aOS(184333335, _omitFieldNames ? '' : 'launchtype')
    ..aOM<NetworkConfiguration>(
        240088634, _omitFieldNames ? '' : 'networkconfiguration',
        subBuilder: NetworkConfiguration.create)
    ..pPM<PlacementConstraint>(
        248464365, _omitFieldNames ? '' : 'placementconstraints',
        subBuilder: PlacementConstraint.create)
    ..pPM<CapacityProviderStrategyItem>(
        273957206, _omitFieldNames ? '' : 'capacityproviderstrategy',
        subBuilder: CapacityProviderStrategyItem.create)
    ..aOS(291739032, _omitFieldNames ? '' : 'referenceid')
    ..m<$core.String, $core.String>(381526209, _omitFieldNames ? '' : 'tags',
        entryClassName: 'EcsParameters.TagsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('scheduler'))
    ..aI(398407508, _omitFieldNames ? '' : 'taskcount')
    ..aOS(405557622, _omitFieldNames ? '' : 'propagatetags')
    ..aOB(451374779, _omitFieldNames ? '' : 'enableexecutecommand')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EcsParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EcsParameters copyWith(void Function(EcsParameters) updates) =>
      super.copyWith((message) => updates(message as EcsParameters))
          as EcsParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EcsParameters create() => EcsParameters._();
  @$core.override
  EcsParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EcsParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EcsParameters>(create);
  static EcsParameters? _defaultInstance;

  @$pb.TagNumber(25036678)
  $pb.PbList<PlacementStrategy> get placementstrategy => $_getList(0);

  @$pb.TagNumber(82234477)
  $core.String get taskdefinitionarn => $_getSZ(1);
  @$pb.TagNumber(82234477)
  set taskdefinitionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82234477)
  $core.bool hasTaskdefinitionarn() => $_has(1);
  @$pb.TagNumber(82234477)
  void clearTaskdefinitionarn() => $_clearField(82234477);

  @$pb.TagNumber(91525165)
  $core.String get group => $_getSZ(2);
  @$pb.TagNumber(91525165)
  set group($core.String value) => $_setString(2, value);
  @$pb.TagNumber(91525165)
  $core.bool hasGroup() => $_has(2);
  @$pb.TagNumber(91525165)
  void clearGroup() => $_clearField(91525165);

  @$pb.TagNumber(139924287)
  $core.String get platformversion => $_getSZ(3);
  @$pb.TagNumber(139924287)
  set platformversion($core.String value) => $_setString(3, value);
  @$pb.TagNumber(139924287)
  $core.bool hasPlatformversion() => $_has(3);
  @$pb.TagNumber(139924287)
  void clearPlatformversion() => $_clearField(139924287);

  @$pb.TagNumber(146161174)
  $core.bool get enableecsmanagedtags => $_getBF(4);
  @$pb.TagNumber(146161174)
  set enableecsmanagedtags($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(146161174)
  $core.bool hasEnableecsmanagedtags() => $_has(4);
  @$pb.TagNumber(146161174)
  void clearEnableecsmanagedtags() => $_clearField(146161174);

  @$pb.TagNumber(184333335)
  $core.String get launchtype => $_getSZ(5);
  @$pb.TagNumber(184333335)
  set launchtype($core.String value) => $_setString(5, value);
  @$pb.TagNumber(184333335)
  $core.bool hasLaunchtype() => $_has(5);
  @$pb.TagNumber(184333335)
  void clearLaunchtype() => $_clearField(184333335);

  @$pb.TagNumber(240088634)
  NetworkConfiguration get networkconfiguration => $_getN(6);
  @$pb.TagNumber(240088634)
  set networkconfiguration(NetworkConfiguration value) =>
      $_setField(240088634, value);
  @$pb.TagNumber(240088634)
  $core.bool hasNetworkconfiguration() => $_has(6);
  @$pb.TagNumber(240088634)
  void clearNetworkconfiguration() => $_clearField(240088634);
  @$pb.TagNumber(240088634)
  NetworkConfiguration ensureNetworkconfiguration() => $_ensure(6);

  @$pb.TagNumber(248464365)
  $pb.PbList<PlacementConstraint> get placementconstraints => $_getList(7);

  @$pb.TagNumber(273957206)
  $pb.PbList<CapacityProviderStrategyItem> get capacityproviderstrategy =>
      $_getList(8);

  @$pb.TagNumber(291739032)
  $core.String get referenceid => $_getSZ(9);
  @$pb.TagNumber(291739032)
  set referenceid($core.String value) => $_setString(9, value);
  @$pb.TagNumber(291739032)
  $core.bool hasReferenceid() => $_has(9);
  @$pb.TagNumber(291739032)
  void clearReferenceid() => $_clearField(291739032);

  @$pb.TagNumber(381526209)
  $pb.PbMap<$core.String, $core.String> get tags => $_getMap(10);

  @$pb.TagNumber(398407508)
  $core.int get taskcount => $_getIZ(11);
  @$pb.TagNumber(398407508)
  set taskcount($core.int value) => $_setSignedInt32(11, value);
  @$pb.TagNumber(398407508)
  $core.bool hasTaskcount() => $_has(11);
  @$pb.TagNumber(398407508)
  void clearTaskcount() => $_clearField(398407508);

  @$pb.TagNumber(405557622)
  $core.String get propagatetags => $_getSZ(12);
  @$pb.TagNumber(405557622)
  set propagatetags($core.String value) => $_setString(12, value);
  @$pb.TagNumber(405557622)
  $core.bool hasPropagatetags() => $_has(12);
  @$pb.TagNumber(405557622)
  void clearPropagatetags() => $_clearField(405557622);

  @$pb.TagNumber(451374779)
  $core.bool get enableexecutecommand => $_getBF(13);
  @$pb.TagNumber(451374779)
  set enableexecutecommand($core.bool value) => $_setBool(13, value);
  @$pb.TagNumber(451374779)
  $core.bool hasEnableexecutecommand() => $_has(13);
  @$pb.TagNumber(451374779)
  void clearEnableexecutecommand() => $_clearField(451374779);
}

class EventBridgeParameters extends $pb.GeneratedMessage {
  factory EventBridgeParameters({
    $core.String? detailtype,
    $core.String? source,
  }) {
    final result = create();
    if (detailtype != null) result.detailtype = detailtype;
    if (source != null) result.source = source;
    return result;
  }

  EventBridgeParameters._();

  factory EventBridgeParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventBridgeParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventBridgeParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(11758589, _omitFieldNames ? '' : 'detailtype')
    ..aOS(31630329, _omitFieldNames ? '' : 'source')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventBridgeParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventBridgeParameters copyWith(
          void Function(EventBridgeParameters) updates) =>
      super.copyWith((message) => updates(message as EventBridgeParameters))
          as EventBridgeParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventBridgeParameters create() => EventBridgeParameters._();
  @$core.override
  EventBridgeParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventBridgeParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EventBridgeParameters>(create);
  static EventBridgeParameters? _defaultInstance;

  @$pb.TagNumber(11758589)
  $core.String get detailtype => $_getSZ(0);
  @$pb.TagNumber(11758589)
  set detailtype($core.String value) => $_setString(0, value);
  @$pb.TagNumber(11758589)
  $core.bool hasDetailtype() => $_has(0);
  @$pb.TagNumber(11758589)
  void clearDetailtype() => $_clearField(11758589);

  @$pb.TagNumber(31630329)
  $core.String get source => $_getSZ(1);
  @$pb.TagNumber(31630329)
  set source($core.String value) => $_setString(1, value);
  @$pb.TagNumber(31630329)
  $core.bool hasSource() => $_has(1);
  @$pb.TagNumber(31630329)
  void clearSource() => $_clearField(31630329);
}

class FlexibleTimeWindow extends $pb.GeneratedMessage {
  factory FlexibleTimeWindow({
    $core.String? mode,
    $core.int? maximumwindowinminutes,
  }) {
    final result = create();
    if (mode != null) result.mode = mode;
    if (maximumwindowinminutes != null)
      result.maximumwindowinminutes = maximumwindowinminutes;
    return result;
  }

  FlexibleTimeWindow._();

  factory FlexibleTimeWindow.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FlexibleTimeWindow.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FlexibleTimeWindow',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(323909427, _omitFieldNames ? '' : 'mode')
    ..aI(482931584, _omitFieldNames ? '' : 'maximumwindowinminutes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FlexibleTimeWindow clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FlexibleTimeWindow copyWith(void Function(FlexibleTimeWindow) updates) =>
      super.copyWith((message) => updates(message as FlexibleTimeWindow))
          as FlexibleTimeWindow;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FlexibleTimeWindow create() => FlexibleTimeWindow._();
  @$core.override
  FlexibleTimeWindow createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FlexibleTimeWindow getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FlexibleTimeWindow>(create);
  static FlexibleTimeWindow? _defaultInstance;

  @$pb.TagNumber(323909427)
  $core.String get mode => $_getSZ(0);
  @$pb.TagNumber(323909427)
  set mode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(323909427)
  $core.bool hasMode() => $_has(0);
  @$pb.TagNumber(323909427)
  void clearMode() => $_clearField(323909427);

  @$pb.TagNumber(482931584)
  $core.int get maximumwindowinminutes => $_getIZ(1);
  @$pb.TagNumber(482931584)
  set maximumwindowinminutes($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(482931584)
  $core.bool hasMaximumwindowinminutes() => $_has(1);
  @$pb.TagNumber(482931584)
  void clearMaximumwindowinminutes() => $_clearField(482931584);
}

class GetScheduleGroupInput extends $pb.GeneratedMessage {
  factory GetScheduleGroupInput({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  GetScheduleGroupInput._();

  factory GetScheduleGroupInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetScheduleGroupInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetScheduleGroupInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetScheduleGroupInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetScheduleGroupInput copyWith(
          void Function(GetScheduleGroupInput) updates) =>
      super.copyWith((message) => updates(message as GetScheduleGroupInput))
          as GetScheduleGroupInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetScheduleGroupInput create() => GetScheduleGroupInput._();
  @$core.override
  GetScheduleGroupInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetScheduleGroupInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetScheduleGroupInput>(create);
  static GetScheduleGroupInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class GetScheduleGroupOutput extends $pb.GeneratedMessage {
  factory GetScheduleGroupOutput({
    $core.String? lastmodificationdate,
    $core.String? name,
    $core.String? creationdate,
    $core.String? arn,
    $core.String? state,
  }) {
    final result = create();
    if (lastmodificationdate != null)
      result.lastmodificationdate = lastmodificationdate;
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (arn != null) result.arn = arn;
    if (state != null) result.state = state;
    return result;
  }

  GetScheduleGroupOutput._();

  factory GetScheduleGroupOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetScheduleGroupOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetScheduleGroupOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(4750358, _omitFieldNames ? '' : 'lastmodificationdate')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetScheduleGroupOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetScheduleGroupOutput copyWith(
          void Function(GetScheduleGroupOutput) updates) =>
      super.copyWith((message) => updates(message as GetScheduleGroupOutput))
          as GetScheduleGroupOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetScheduleGroupOutput create() => GetScheduleGroupOutput._();
  @$core.override
  GetScheduleGroupOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetScheduleGroupOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetScheduleGroupOutput>(create);
  static GetScheduleGroupOutput? _defaultInstance;

  @$pb.TagNumber(4750358)
  $core.String get lastmodificationdate => $_getSZ(0);
  @$pb.TagNumber(4750358)
  set lastmodificationdate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(4750358)
  $core.bool hasLastmodificationdate() => $_has(0);
  @$pb.TagNumber(4750358)
  void clearLastmodificationdate() => $_clearField(4750358);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(2);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(2);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(3);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(3);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(4);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(4, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(4);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class GetScheduleInput extends $pb.GeneratedMessage {
  factory GetScheduleInput({
    $core.String? name,
    $core.String? groupname,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (groupname != null) result.groupname = groupname;
    return result;
  }

  GetScheduleInput._();

  factory GetScheduleInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetScheduleInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetScheduleInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(357049672, _omitFieldNames ? '' : 'groupname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetScheduleInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetScheduleInput copyWith(void Function(GetScheduleInput) updates) =>
      super.copyWith((message) => updates(message as GetScheduleInput))
          as GetScheduleInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetScheduleInput create() => GetScheduleInput._();
  @$core.override
  GetScheduleInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetScheduleInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetScheduleInput>(create);
  static GetScheduleInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(357049672)
  $core.String get groupname => $_getSZ(1);
  @$pb.TagNumber(357049672)
  set groupname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(357049672)
  $core.bool hasGroupname() => $_has(1);
  @$pb.TagNumber(357049672)
  void clearGroupname() => $_clearField(357049672);
}

class GetScheduleOutput extends $pb.GeneratedMessage {
  factory GetScheduleOutput({
    FlexibleTimeWindow? flexibletimewindow,
    $core.String? lastmodificationdate,
    $core.String? enddate,
    $core.String? kmskeyarn,
    $core.String? description,
    Target? target,
    $core.String? name,
    $core.String? actionaftercompletion,
    $core.String? creationdate,
    $core.String? groupname,
    $core.String? scheduleexpressiontimezone,
    $core.String? arn,
    $core.String? startdate,
    $core.String? scheduleexpression,
    $core.String? state,
  }) {
    final result = create();
    if (flexibletimewindow != null)
      result.flexibletimewindow = flexibletimewindow;
    if (lastmodificationdate != null)
      result.lastmodificationdate = lastmodificationdate;
    if (enddate != null) result.enddate = enddate;
    if (kmskeyarn != null) result.kmskeyarn = kmskeyarn;
    if (description != null) result.description = description;
    if (target != null) result.target = target;
    if (name != null) result.name = name;
    if (actionaftercompletion != null)
      result.actionaftercompletion = actionaftercompletion;
    if (creationdate != null) result.creationdate = creationdate;
    if (groupname != null) result.groupname = groupname;
    if (scheduleexpressiontimezone != null)
      result.scheduleexpressiontimezone = scheduleexpressiontimezone;
    if (arn != null) result.arn = arn;
    if (startdate != null) result.startdate = startdate;
    if (scheduleexpression != null)
      result.scheduleexpression = scheduleexpression;
    if (state != null) result.state = state;
    return result;
  }

  GetScheduleOutput._();

  factory GetScheduleOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetScheduleOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetScheduleOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOM<FlexibleTimeWindow>(
        2518952, _omitFieldNames ? '' : 'flexibletimewindow',
        subBuilder: FlexibleTimeWindow.create)
    ..aOS(4750358, _omitFieldNames ? '' : 'lastmodificationdate')
    ..aOS(77486543, _omitFieldNames ? '' : 'enddate')
    ..aOS(110041649, _omitFieldNames ? '' : 'kmskeyarn')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOM<Target>(191361385, _omitFieldNames ? '' : 'target',
        subBuilder: Target.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(282350906, _omitFieldNames ? '' : 'actionaftercompletion')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aOS(357049672, _omitFieldNames ? '' : 'groupname')
    ..aOS(400730326, _omitFieldNames ? '' : 'scheduleexpressiontimezone')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(445135996, _omitFieldNames ? '' : 'startdate')
    ..aOS(446089471, _omitFieldNames ? '' : 'scheduleexpression')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetScheduleOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetScheduleOutput copyWith(void Function(GetScheduleOutput) updates) =>
      super.copyWith((message) => updates(message as GetScheduleOutput))
          as GetScheduleOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetScheduleOutput create() => GetScheduleOutput._();
  @$core.override
  GetScheduleOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetScheduleOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetScheduleOutput>(create);
  static GetScheduleOutput? _defaultInstance;

  @$pb.TagNumber(2518952)
  FlexibleTimeWindow get flexibletimewindow => $_getN(0);
  @$pb.TagNumber(2518952)
  set flexibletimewindow(FlexibleTimeWindow value) =>
      $_setField(2518952, value);
  @$pb.TagNumber(2518952)
  $core.bool hasFlexibletimewindow() => $_has(0);
  @$pb.TagNumber(2518952)
  void clearFlexibletimewindow() => $_clearField(2518952);
  @$pb.TagNumber(2518952)
  FlexibleTimeWindow ensureFlexibletimewindow() => $_ensure(0);

  @$pb.TagNumber(4750358)
  $core.String get lastmodificationdate => $_getSZ(1);
  @$pb.TagNumber(4750358)
  set lastmodificationdate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(4750358)
  $core.bool hasLastmodificationdate() => $_has(1);
  @$pb.TagNumber(4750358)
  void clearLastmodificationdate() => $_clearField(4750358);

  @$pb.TagNumber(77486543)
  $core.String get enddate => $_getSZ(2);
  @$pb.TagNumber(77486543)
  set enddate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(77486543)
  $core.bool hasEnddate() => $_has(2);
  @$pb.TagNumber(77486543)
  void clearEnddate() => $_clearField(77486543);

  @$pb.TagNumber(110041649)
  $core.String get kmskeyarn => $_getSZ(3);
  @$pb.TagNumber(110041649)
  set kmskeyarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(110041649)
  $core.bool hasKmskeyarn() => $_has(3);
  @$pb.TagNumber(110041649)
  void clearKmskeyarn() => $_clearField(110041649);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(4);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(4, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(4);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(191361385)
  Target get target => $_getN(5);
  @$pb.TagNumber(191361385)
  set target(Target value) => $_setField(191361385, value);
  @$pb.TagNumber(191361385)
  $core.bool hasTarget() => $_has(5);
  @$pb.TagNumber(191361385)
  void clearTarget() => $_clearField(191361385);
  @$pb.TagNumber(191361385)
  Target ensureTarget() => $_ensure(5);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(282350906)
  $core.String get actionaftercompletion => $_getSZ(7);
  @$pb.TagNumber(282350906)
  set actionaftercompletion($core.String value) => $_setString(7, value);
  @$pb.TagNumber(282350906)
  $core.bool hasActionaftercompletion() => $_has(7);
  @$pb.TagNumber(282350906)
  void clearActionaftercompletion() => $_clearField(282350906);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(8);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(8, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(8);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(357049672)
  $core.String get groupname => $_getSZ(9);
  @$pb.TagNumber(357049672)
  set groupname($core.String value) => $_setString(9, value);
  @$pb.TagNumber(357049672)
  $core.bool hasGroupname() => $_has(9);
  @$pb.TagNumber(357049672)
  void clearGroupname() => $_clearField(357049672);

  @$pb.TagNumber(400730326)
  $core.String get scheduleexpressiontimezone => $_getSZ(10);
  @$pb.TagNumber(400730326)
  set scheduleexpressiontimezone($core.String value) => $_setString(10, value);
  @$pb.TagNumber(400730326)
  $core.bool hasScheduleexpressiontimezone() => $_has(10);
  @$pb.TagNumber(400730326)
  void clearScheduleexpressiontimezone() => $_clearField(400730326);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(11);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(11, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(11);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(445135996)
  $core.String get startdate => $_getSZ(12);
  @$pb.TagNumber(445135996)
  set startdate($core.String value) => $_setString(12, value);
  @$pb.TagNumber(445135996)
  $core.bool hasStartdate() => $_has(12);
  @$pb.TagNumber(445135996)
  void clearStartdate() => $_clearField(445135996);

  @$pb.TagNumber(446089471)
  $core.String get scheduleexpression => $_getSZ(13);
  @$pb.TagNumber(446089471)
  set scheduleexpression($core.String value) => $_setString(13, value);
  @$pb.TagNumber(446089471)
  $core.bool hasScheduleexpression() => $_has(13);
  @$pb.TagNumber(446089471)
  void clearScheduleexpression() => $_clearField(446089471);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(14);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(14, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(14);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class InternalServerException extends $pb.GeneratedMessage {
  factory InternalServerException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InternalServerException._();

  factory InternalServerException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InternalServerException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InternalServerException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServerException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServerException copyWith(
          void Function(InternalServerException) updates) =>
      super.copyWith((message) => updates(message as InternalServerException))
          as InternalServerException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InternalServerException create() => InternalServerException._();
  @$core.override
  InternalServerException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InternalServerException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InternalServerException>(create);
  static InternalServerException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class KinesisParameters extends $pb.GeneratedMessage {
  factory KinesisParameters({
    $core.String? partitionkey,
  }) {
    final result = create();
    if (partitionkey != null) result.partitionkey = partitionkey;
    return result;
  }

  KinesisParameters._();

  factory KinesisParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KinesisParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KinesisParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(379379617, _omitFieldNames ? '' : 'partitionkey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KinesisParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KinesisParameters copyWith(void Function(KinesisParameters) updates) =>
      super.copyWith((message) => updates(message as KinesisParameters))
          as KinesisParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KinesisParameters create() => KinesisParameters._();
  @$core.override
  KinesisParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KinesisParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KinesisParameters>(create);
  static KinesisParameters? _defaultInstance;

  @$pb.TagNumber(379379617)
  $core.String get partitionkey => $_getSZ(0);
  @$pb.TagNumber(379379617)
  set partitionkey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(379379617)
  $core.bool hasPartitionkey() => $_has(0);
  @$pb.TagNumber(379379617)
  void clearPartitionkey() => $_clearField(379379617);
}

class ListScheduleGroupsInput extends $pb.GeneratedMessage {
  factory ListScheduleGroupsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? nameprefix,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (nameprefix != null) result.nameprefix = nameprefix;
    return result;
  }

  ListScheduleGroupsInput._();

  factory ListScheduleGroupsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListScheduleGroupsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListScheduleGroupsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListScheduleGroupsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListScheduleGroupsInput copyWith(
          void Function(ListScheduleGroupsInput) updates) =>
      super.copyWith((message) => updates(message as ListScheduleGroupsInput))
          as ListScheduleGroupsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListScheduleGroupsInput create() => ListScheduleGroupsInput._();
  @$core.override
  ListScheduleGroupsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListScheduleGroupsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListScheduleGroupsInput>(create);
  static ListScheduleGroupsInput? _defaultInstance;

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

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(2);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(2, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(2);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);
}

class ListScheduleGroupsOutput extends $pb.GeneratedMessage {
  factory ListScheduleGroupsOutput({
    $core.Iterable<ScheduleGroupSummary>? schedulegroups,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (schedulegroups != null) result.schedulegroups.addAll(schedulegroups);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListScheduleGroupsOutput._();

  factory ListScheduleGroupsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListScheduleGroupsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListScheduleGroupsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..pPM<ScheduleGroupSummary>(
        136082885, _omitFieldNames ? '' : 'schedulegroups',
        subBuilder: ScheduleGroupSummary.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListScheduleGroupsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListScheduleGroupsOutput copyWith(
          void Function(ListScheduleGroupsOutput) updates) =>
      super.copyWith((message) => updates(message as ListScheduleGroupsOutput))
          as ListScheduleGroupsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListScheduleGroupsOutput create() => ListScheduleGroupsOutput._();
  @$core.override
  ListScheduleGroupsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListScheduleGroupsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListScheduleGroupsOutput>(create);
  static ListScheduleGroupsOutput? _defaultInstance;

  @$pb.TagNumber(136082885)
  $pb.PbList<ScheduleGroupSummary> get schedulegroups => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListSchedulesInput extends $pb.GeneratedMessage {
  factory ListSchedulesInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? groupname,
    $core.String? nameprefix,
    $core.String? state,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (groupname != null) result.groupname = groupname;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (state != null) result.state = state;
    return result;
  }

  ListSchedulesInput._();

  factory ListSchedulesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSchedulesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSchedulesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(357049672, _omitFieldNames ? '' : 'groupname')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSchedulesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSchedulesInput copyWith(void Function(ListSchedulesInput) updates) =>
      super.copyWith((message) => updates(message as ListSchedulesInput))
          as ListSchedulesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSchedulesInput create() => ListSchedulesInput._();
  @$core.override
  ListSchedulesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSchedulesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSchedulesInput>(create);
  static ListSchedulesInput? _defaultInstance;

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

  @$pb.TagNumber(357049672)
  $core.String get groupname => $_getSZ(2);
  @$pb.TagNumber(357049672)
  set groupname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(357049672)
  $core.bool hasGroupname() => $_has(2);
  @$pb.TagNumber(357049672)
  void clearGroupname() => $_clearField(357049672);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(3);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(3, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(3);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(4);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(4, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(4);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class ListSchedulesOutput extends $pb.GeneratedMessage {
  factory ListSchedulesOutput({
    $core.Iterable<ScheduleSummary>? schedules,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (schedules != null) result.schedules.addAll(schedules);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListSchedulesOutput._();

  factory ListSchedulesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListSchedulesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListSchedulesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..pPM<ScheduleSummary>(18925646, _omitFieldNames ? '' : 'schedules',
        subBuilder: ScheduleSummary.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSchedulesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListSchedulesOutput copyWith(void Function(ListSchedulesOutput) updates) =>
      super.copyWith((message) => updates(message as ListSchedulesOutput))
          as ListSchedulesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListSchedulesOutput create() => ListSchedulesOutput._();
  @$core.override
  ListSchedulesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListSchedulesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListSchedulesOutput>(create);
  static ListSchedulesOutput? _defaultInstance;

  @$pb.TagNumber(18925646)
  $pb.PbList<ScheduleSummary> get schedules => $_getList(0);

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
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

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(0);
}

class NetworkConfiguration extends $pb.GeneratedMessage {
  factory NetworkConfiguration({
    AwsVpcConfiguration? awsvpcconfiguration,
  }) {
    final result = create();
    if (awsvpcconfiguration != null)
      result.awsvpcconfiguration = awsvpcconfiguration;
    return result;
  }

  NetworkConfiguration._();

  factory NetworkConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NetworkConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NetworkConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOM<AwsVpcConfiguration>(
        223852630, _omitFieldNames ? '' : 'awsvpcconfiguration',
        subBuilder: AwsVpcConfiguration.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NetworkConfiguration copyWith(void Function(NetworkConfiguration) updates) =>
      super.copyWith((message) => updates(message as NetworkConfiguration))
          as NetworkConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NetworkConfiguration create() => NetworkConfiguration._();
  @$core.override
  NetworkConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NetworkConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NetworkConfiguration>(create);
  static NetworkConfiguration? _defaultInstance;

  @$pb.TagNumber(223852630)
  AwsVpcConfiguration get awsvpcconfiguration => $_getN(0);
  @$pb.TagNumber(223852630)
  set awsvpcconfiguration(AwsVpcConfiguration value) =>
      $_setField(223852630, value);
  @$pb.TagNumber(223852630)
  $core.bool hasAwsvpcconfiguration() => $_has(0);
  @$pb.TagNumber(223852630)
  void clearAwsvpcconfiguration() => $_clearField(223852630);
  @$pb.TagNumber(223852630)
  AwsVpcConfiguration ensureAwsvpcconfiguration() => $_ensure(0);
}

class PlacementConstraint extends $pb.GeneratedMessage {
  factory PlacementConstraint({
    $core.String? expression,
    $core.String? type,
  }) {
    final result = create();
    if (expression != null) result.expression = expression;
    if (type != null) result.type = type;
    return result;
  }

  PlacementConstraint._();

  factory PlacementConstraint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PlacementConstraint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PlacementConstraint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(253079532, _omitFieldNames ? '' : 'expression')
    ..aOS(287830350, _omitFieldNames ? '' : 'type')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlacementConstraint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlacementConstraint copyWith(void Function(PlacementConstraint) updates) =>
      super.copyWith((message) => updates(message as PlacementConstraint))
          as PlacementConstraint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PlacementConstraint create() => PlacementConstraint._();
  @$core.override
  PlacementConstraint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PlacementConstraint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PlacementConstraint>(create);
  static PlacementConstraint? _defaultInstance;

  @$pb.TagNumber(253079532)
  $core.String get expression => $_getSZ(0);
  @$pb.TagNumber(253079532)
  set expression($core.String value) => $_setString(0, value);
  @$pb.TagNumber(253079532)
  $core.bool hasExpression() => $_has(0);
  @$pb.TagNumber(253079532)
  void clearExpression() => $_clearField(253079532);

  @$pb.TagNumber(287830350)
  $core.String get type => $_getSZ(1);
  @$pb.TagNumber(287830350)
  set type($core.String value) => $_setString(1, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);
}

class PlacementStrategy extends $pb.GeneratedMessage {
  factory PlacementStrategy({
    $core.String? field_125985384,
    $core.String? type,
  }) {
    final result = create();
    if (field_125985384 != null) result.field_125985384 = field_125985384;
    if (type != null) result.type = type;
    return result;
  }

  PlacementStrategy._();

  factory PlacementStrategy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PlacementStrategy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PlacementStrategy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(125985384, _omitFieldNames ? '' : 'field')
    ..aOS(287830350, _omitFieldNames ? '' : 'type')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlacementStrategy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PlacementStrategy copyWith(void Function(PlacementStrategy) updates) =>
      super.copyWith((message) => updates(message as PlacementStrategy))
          as PlacementStrategy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PlacementStrategy create() => PlacementStrategy._();
  @$core.override
  PlacementStrategy createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PlacementStrategy getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PlacementStrategy>(create);
  static PlacementStrategy? _defaultInstance;

  @$pb.TagNumber(125985384)
  $core.String get field_125985384 => $_getSZ(0);
  @$pb.TagNumber(125985384)
  set field_125985384($core.String value) => $_setString(0, value);
  @$pb.TagNumber(125985384)
  $core.bool hasField_125985384() => $_has(0);
  @$pb.TagNumber(125985384)
  void clearField_125985384() => $_clearField(125985384);

  @$pb.TagNumber(287830350)
  $core.String get type => $_getSZ(1);
  @$pb.TagNumber(287830350)
  set type($core.String value) => $_setString(1, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
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

class RetryPolicy extends $pb.GeneratedMessage {
  factory RetryPolicy({
    $core.int? maximumretryattempts,
    $core.int? maximumeventageinseconds,
  }) {
    final result = create();
    if (maximumretryattempts != null)
      result.maximumretryattempts = maximumretryattempts;
    if (maximumeventageinseconds != null)
      result.maximumeventageinseconds = maximumeventageinseconds;
    return result;
  }

  RetryPolicy._();

  factory RetryPolicy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RetryPolicy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RetryPolicy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aI(112088128, _omitFieldNames ? '' : 'maximumretryattempts')
    ..aI(393041563, _omitFieldNames ? '' : 'maximumeventageinseconds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RetryPolicy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RetryPolicy copyWith(void Function(RetryPolicy) updates) =>
      super.copyWith((message) => updates(message as RetryPolicy))
          as RetryPolicy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RetryPolicy create() => RetryPolicy._();
  @$core.override
  RetryPolicy createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RetryPolicy getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RetryPolicy>(create);
  static RetryPolicy? _defaultInstance;

  @$pb.TagNumber(112088128)
  $core.int get maximumretryattempts => $_getIZ(0);
  @$pb.TagNumber(112088128)
  set maximumretryattempts($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(112088128)
  $core.bool hasMaximumretryattempts() => $_has(0);
  @$pb.TagNumber(112088128)
  void clearMaximumretryattempts() => $_clearField(112088128);

  @$pb.TagNumber(393041563)
  $core.int get maximumeventageinseconds => $_getIZ(1);
  @$pb.TagNumber(393041563)
  set maximumeventageinseconds($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(393041563)
  $core.bool hasMaximumeventageinseconds() => $_has(1);
  @$pb.TagNumber(393041563)
  void clearMaximumeventageinseconds() => $_clearField(393041563);
}

class SageMakerPipelineParameter extends $pb.GeneratedMessage {
  factory SageMakerPipelineParameter({
    $core.String? name,
    $core.String? value,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (value != null) result.value = value;
    return result;
  }

  SageMakerPipelineParameter._();

  factory SageMakerPipelineParameter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SageMakerPipelineParameter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SageMakerPipelineParameter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SageMakerPipelineParameter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SageMakerPipelineParameter copyWith(
          void Function(SageMakerPipelineParameter) updates) =>
      super.copyWith(
              (message) => updates(message as SageMakerPipelineParameter))
          as SageMakerPipelineParameter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SageMakerPipelineParameter create() => SageMakerPipelineParameter._();
  @$core.override
  SageMakerPipelineParameter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SageMakerPipelineParameter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SageMakerPipelineParameter>(create);
  static SageMakerPipelineParameter? _defaultInstance;

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

class SageMakerPipelineParameters extends $pb.GeneratedMessage {
  factory SageMakerPipelineParameters({
    $core.Iterable<SageMakerPipelineParameter>? pipelineparameterlist,
  }) {
    final result = create();
    if (pipelineparameterlist != null)
      result.pipelineparameterlist.addAll(pipelineparameterlist);
    return result;
  }

  SageMakerPipelineParameters._();

  factory SageMakerPipelineParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SageMakerPipelineParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SageMakerPipelineParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..pPM<SageMakerPipelineParameter>(
        198270119, _omitFieldNames ? '' : 'pipelineparameterlist',
        subBuilder: SageMakerPipelineParameter.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SageMakerPipelineParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SageMakerPipelineParameters copyWith(
          void Function(SageMakerPipelineParameters) updates) =>
      super.copyWith(
              (message) => updates(message as SageMakerPipelineParameters))
          as SageMakerPipelineParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SageMakerPipelineParameters create() =>
      SageMakerPipelineParameters._();
  @$core.override
  SageMakerPipelineParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SageMakerPipelineParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SageMakerPipelineParameters>(create);
  static SageMakerPipelineParameters? _defaultInstance;

  @$pb.TagNumber(198270119)
  $pb.PbList<SageMakerPipelineParameter> get pipelineparameterlist =>
      $_getList(0);
}

class ScheduleGroupSummary extends $pb.GeneratedMessage {
  factory ScheduleGroupSummary({
    $core.String? lastmodificationdate,
    $core.String? name,
    $core.String? creationdate,
    $core.String? arn,
    $core.String? state,
  }) {
    final result = create();
    if (lastmodificationdate != null)
      result.lastmodificationdate = lastmodificationdate;
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (arn != null) result.arn = arn;
    if (state != null) result.state = state;
    return result;
  }

  ScheduleGroupSummary._();

  factory ScheduleGroupSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduleGroupSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduleGroupSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(4750358, _omitFieldNames ? '' : 'lastmodificationdate')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleGroupSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleGroupSummary copyWith(void Function(ScheduleGroupSummary) updates) =>
      super.copyWith((message) => updates(message as ScheduleGroupSummary))
          as ScheduleGroupSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduleGroupSummary create() => ScheduleGroupSummary._();
  @$core.override
  ScheduleGroupSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduleGroupSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduleGroupSummary>(create);
  static ScheduleGroupSummary? _defaultInstance;

  @$pb.TagNumber(4750358)
  $core.String get lastmodificationdate => $_getSZ(0);
  @$pb.TagNumber(4750358)
  set lastmodificationdate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(4750358)
  $core.bool hasLastmodificationdate() => $_has(0);
  @$pb.TagNumber(4750358)
  void clearLastmodificationdate() => $_clearField(4750358);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(2);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(2, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(2);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(3);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(3);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(4);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(4, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(4);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class ScheduleSummary extends $pb.GeneratedMessage {
  factory ScheduleSummary({
    $core.String? lastmodificationdate,
    TargetSummary? target,
    $core.String? name,
    $core.String? creationdate,
    $core.String? groupname,
    $core.String? arn,
    $core.String? state,
  }) {
    final result = create();
    if (lastmodificationdate != null)
      result.lastmodificationdate = lastmodificationdate;
    if (target != null) result.target = target;
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (groupname != null) result.groupname = groupname;
    if (arn != null) result.arn = arn;
    if (state != null) result.state = state;
    return result;
  }

  ScheduleSummary._();

  factory ScheduleSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ScheduleSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ScheduleSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(4750358, _omitFieldNames ? '' : 'lastmodificationdate')
    ..aOM<TargetSummary>(191361385, _omitFieldNames ? '' : 'target',
        subBuilder: TargetSummary.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aOS(357049672, _omitFieldNames ? '' : 'groupname')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ScheduleSummary copyWith(void Function(ScheduleSummary) updates) =>
      super.copyWith((message) => updates(message as ScheduleSummary))
          as ScheduleSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ScheduleSummary create() => ScheduleSummary._();
  @$core.override
  ScheduleSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ScheduleSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ScheduleSummary>(create);
  static ScheduleSummary? _defaultInstance;

  @$pb.TagNumber(4750358)
  $core.String get lastmodificationdate => $_getSZ(0);
  @$pb.TagNumber(4750358)
  set lastmodificationdate($core.String value) => $_setString(0, value);
  @$pb.TagNumber(4750358)
  $core.bool hasLastmodificationdate() => $_has(0);
  @$pb.TagNumber(4750358)
  void clearLastmodificationdate() => $_clearField(4750358);

  @$pb.TagNumber(191361385)
  TargetSummary get target => $_getN(1);
  @$pb.TagNumber(191361385)
  set target(TargetSummary value) => $_setField(191361385, value);
  @$pb.TagNumber(191361385)
  $core.bool hasTarget() => $_has(1);
  @$pb.TagNumber(191361385)
  void clearTarget() => $_clearField(191361385);
  @$pb.TagNumber(191361385)
  TargetSummary ensureTarget() => $_ensure(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(3);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(3);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(357049672)
  $core.String get groupname => $_getSZ(4);
  @$pb.TagNumber(357049672)
  set groupname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(357049672)
  $core.bool hasGroupname() => $_has(4);
  @$pb.TagNumber(357049672)
  void clearGroupname() => $_clearField(357049672);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(5);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(5);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(6);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(6, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(6);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class SqsParameters extends $pb.GeneratedMessage {
  factory SqsParameters({
    $core.String? messagegroupid,
  }) {
    final result = create();
    if (messagegroupid != null) result.messagegroupid = messagegroupid;
    return result;
  }

  SqsParameters._();

  factory SqsParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SqsParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SqsParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(419537435, _omitFieldNames ? '' : 'messagegroupid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqsParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SqsParameters copyWith(void Function(SqsParameters) updates) =>
      super.copyWith((message) => updates(message as SqsParameters))
          as SqsParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SqsParameters create() => SqsParameters._();
  @$core.override
  SqsParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SqsParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SqsParameters>(create);
  static SqsParameters? _defaultInstance;

  @$pb.TagNumber(419537435)
  $core.String get messagegroupid => $_getSZ(0);
  @$pb.TagNumber(419537435)
  set messagegroupid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(419537435)
  $core.bool hasMessagegroupid() => $_has(0);
  @$pb.TagNumber(419537435)
  void clearMessagegroupid() => $_clearField(419537435);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
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

class Target extends $pb.GeneratedMessage {
  factory Target({
    KinesisParameters? kinesisparameters,
    DeadLetterConfig? deadletterconfig,
    SqsParameters? sqsparameters,
    RetryPolicy? retrypolicy,
    EventBridgeParameters? eventbridgeparameters,
    $core.String? rolearn,
    $core.String? arn,
    EcsParameters? ecsparameters,
    SageMakerPipelineParameters? sagemakerpipelineparameters,
    $core.String? input,
  }) {
    final result = create();
    if (kinesisparameters != null) result.kinesisparameters = kinesisparameters;
    if (deadletterconfig != null) result.deadletterconfig = deadletterconfig;
    if (sqsparameters != null) result.sqsparameters = sqsparameters;
    if (retrypolicy != null) result.retrypolicy = retrypolicy;
    if (eventbridgeparameters != null)
      result.eventbridgeparameters = eventbridgeparameters;
    if (rolearn != null) result.rolearn = rolearn;
    if (arn != null) result.arn = arn;
    if (ecsparameters != null) result.ecsparameters = ecsparameters;
    if (sagemakerpipelineparameters != null)
      result.sagemakerpipelineparameters = sagemakerpipelineparameters;
    if (input != null) result.input = input;
    return result;
  }

  Target._();

  factory Target.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Target.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Target',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOM<KinesisParameters>(
        70111902, _omitFieldNames ? '' : 'kinesisparameters',
        subBuilder: KinesisParameters.create)
    ..aOM<DeadLetterConfig>(79786642, _omitFieldNames ? '' : 'deadletterconfig',
        subBuilder: DeadLetterConfig.create)
    ..aOM<SqsParameters>(91345999, _omitFieldNames ? '' : 'sqsparameters',
        subBuilder: SqsParameters.create)
    ..aOM<RetryPolicy>(266827188, _omitFieldNames ? '' : 'retrypolicy',
        subBuilder: RetryPolicy.create)
    ..aOM<EventBridgeParameters>(
        285439471, _omitFieldNames ? '' : 'eventbridgeparameters',
        subBuilder: EventBridgeParameters.create)
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOM<EcsParameters>(501521183, _omitFieldNames ? '' : 'ecsparameters',
        subBuilder: EcsParameters.create)
    ..aOM<SageMakerPipelineParameters>(
        513800606, _omitFieldNames ? '' : 'sagemakerpipelineparameters',
        subBuilder: SageMakerPipelineParameters.create)
    ..aOS(529785116, _omitFieldNames ? '' : 'input')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Target clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Target copyWith(void Function(Target) updates) =>
      super.copyWith((message) => updates(message as Target)) as Target;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Target create() => Target._();
  @$core.override
  Target createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Target getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Target>(create);
  static Target? _defaultInstance;

  @$pb.TagNumber(70111902)
  KinesisParameters get kinesisparameters => $_getN(0);
  @$pb.TagNumber(70111902)
  set kinesisparameters(KinesisParameters value) => $_setField(70111902, value);
  @$pb.TagNumber(70111902)
  $core.bool hasKinesisparameters() => $_has(0);
  @$pb.TagNumber(70111902)
  void clearKinesisparameters() => $_clearField(70111902);
  @$pb.TagNumber(70111902)
  KinesisParameters ensureKinesisparameters() => $_ensure(0);

  @$pb.TagNumber(79786642)
  DeadLetterConfig get deadletterconfig => $_getN(1);
  @$pb.TagNumber(79786642)
  set deadletterconfig(DeadLetterConfig value) => $_setField(79786642, value);
  @$pb.TagNumber(79786642)
  $core.bool hasDeadletterconfig() => $_has(1);
  @$pb.TagNumber(79786642)
  void clearDeadletterconfig() => $_clearField(79786642);
  @$pb.TagNumber(79786642)
  DeadLetterConfig ensureDeadletterconfig() => $_ensure(1);

  @$pb.TagNumber(91345999)
  SqsParameters get sqsparameters => $_getN(2);
  @$pb.TagNumber(91345999)
  set sqsparameters(SqsParameters value) => $_setField(91345999, value);
  @$pb.TagNumber(91345999)
  $core.bool hasSqsparameters() => $_has(2);
  @$pb.TagNumber(91345999)
  void clearSqsparameters() => $_clearField(91345999);
  @$pb.TagNumber(91345999)
  SqsParameters ensureSqsparameters() => $_ensure(2);

  @$pb.TagNumber(266827188)
  RetryPolicy get retrypolicy => $_getN(3);
  @$pb.TagNumber(266827188)
  set retrypolicy(RetryPolicy value) => $_setField(266827188, value);
  @$pb.TagNumber(266827188)
  $core.bool hasRetrypolicy() => $_has(3);
  @$pb.TagNumber(266827188)
  void clearRetrypolicy() => $_clearField(266827188);
  @$pb.TagNumber(266827188)
  RetryPolicy ensureRetrypolicy() => $_ensure(3);

  @$pb.TagNumber(285439471)
  EventBridgeParameters get eventbridgeparameters => $_getN(4);
  @$pb.TagNumber(285439471)
  set eventbridgeparameters(EventBridgeParameters value) =>
      $_setField(285439471, value);
  @$pb.TagNumber(285439471)
  $core.bool hasEventbridgeparameters() => $_has(4);
  @$pb.TagNumber(285439471)
  void clearEventbridgeparameters() => $_clearField(285439471);
  @$pb.TagNumber(285439471)
  EventBridgeParameters ensureEventbridgeparameters() => $_ensure(4);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(5);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(5);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(6);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(6);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(501521183)
  EcsParameters get ecsparameters => $_getN(7);
  @$pb.TagNumber(501521183)
  set ecsparameters(EcsParameters value) => $_setField(501521183, value);
  @$pb.TagNumber(501521183)
  $core.bool hasEcsparameters() => $_has(7);
  @$pb.TagNumber(501521183)
  void clearEcsparameters() => $_clearField(501521183);
  @$pb.TagNumber(501521183)
  EcsParameters ensureEcsparameters() => $_ensure(7);

  @$pb.TagNumber(513800606)
  SageMakerPipelineParameters get sagemakerpipelineparameters => $_getN(8);
  @$pb.TagNumber(513800606)
  set sagemakerpipelineparameters(SageMakerPipelineParameters value) =>
      $_setField(513800606, value);
  @$pb.TagNumber(513800606)
  $core.bool hasSagemakerpipelineparameters() => $_has(8);
  @$pb.TagNumber(513800606)
  void clearSagemakerpipelineparameters() => $_clearField(513800606);
  @$pb.TagNumber(513800606)
  SageMakerPipelineParameters ensureSagemakerpipelineparameters() =>
      $_ensure(8);

  @$pb.TagNumber(529785116)
  $core.String get input => $_getSZ(9);
  @$pb.TagNumber(529785116)
  set input($core.String value) => $_setString(9, value);
  @$pb.TagNumber(529785116)
  $core.bool hasInput() => $_has(9);
  @$pb.TagNumber(529785116)
  void clearInput() => $_clearField(529785116);
}

class TargetSummary extends $pb.GeneratedMessage {
  factory TargetSummary({
    $core.String? arn,
  }) {
    final result = create();
    if (arn != null) result.arn = arn;
    return result;
  }

  TargetSummary._();

  factory TargetSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TargetSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TargetSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TargetSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TargetSummary copyWith(void Function(TargetSummary) updates) =>
      super.copyWith((message) => updates(message as TargetSummary))
          as TargetSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TargetSummary create() => TargetSummary._();
  @$core.override
  TargetSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TargetSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TargetSummary>(create);
  static TargetSummary? _defaultInstance;

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(0);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(0);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class ThrottlingException extends $pb.GeneratedMessage {
  factory ThrottlingException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
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

class UpdateScheduleInput extends $pb.GeneratedMessage {
  factory UpdateScheduleInput({
    FlexibleTimeWindow? flexibletimewindow,
    $core.String? enddate,
    $core.String? kmskeyarn,
    $core.String? description,
    $core.String? clienttoken,
    Target? target,
    $core.String? name,
    $core.String? actionaftercompletion,
    $core.String? groupname,
    $core.String? scheduleexpressiontimezone,
    $core.String? startdate,
    $core.String? scheduleexpression,
    $core.String? state,
  }) {
    final result = create();
    if (flexibletimewindow != null)
      result.flexibletimewindow = flexibletimewindow;
    if (enddate != null) result.enddate = enddate;
    if (kmskeyarn != null) result.kmskeyarn = kmskeyarn;
    if (description != null) result.description = description;
    if (clienttoken != null) result.clienttoken = clienttoken;
    if (target != null) result.target = target;
    if (name != null) result.name = name;
    if (actionaftercompletion != null)
      result.actionaftercompletion = actionaftercompletion;
    if (groupname != null) result.groupname = groupname;
    if (scheduleexpressiontimezone != null)
      result.scheduleexpressiontimezone = scheduleexpressiontimezone;
    if (startdate != null) result.startdate = startdate;
    if (scheduleexpression != null)
      result.scheduleexpression = scheduleexpression;
    if (state != null) result.state = state;
    return result;
  }

  UpdateScheduleInput._();

  factory UpdateScheduleInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateScheduleInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateScheduleInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOM<FlexibleTimeWindow>(
        2518952, _omitFieldNames ? '' : 'flexibletimewindow',
        subBuilder: FlexibleTimeWindow.create)
    ..aOS(77486543, _omitFieldNames ? '' : 'enddate')
    ..aOS(110041649, _omitFieldNames ? '' : 'kmskeyarn')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(137297356, _omitFieldNames ? '' : 'clienttoken')
    ..aOM<Target>(191361385, _omitFieldNames ? '' : 'target',
        subBuilder: Target.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(282350906, _omitFieldNames ? '' : 'actionaftercompletion')
    ..aOS(357049672, _omitFieldNames ? '' : 'groupname')
    ..aOS(400730326, _omitFieldNames ? '' : 'scheduleexpressiontimezone')
    ..aOS(445135996, _omitFieldNames ? '' : 'startdate')
    ..aOS(446089471, _omitFieldNames ? '' : 'scheduleexpression')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateScheduleInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateScheduleInput copyWith(void Function(UpdateScheduleInput) updates) =>
      super.copyWith((message) => updates(message as UpdateScheduleInput))
          as UpdateScheduleInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateScheduleInput create() => UpdateScheduleInput._();
  @$core.override
  UpdateScheduleInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateScheduleInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateScheduleInput>(create);
  static UpdateScheduleInput? _defaultInstance;

  @$pb.TagNumber(2518952)
  FlexibleTimeWindow get flexibletimewindow => $_getN(0);
  @$pb.TagNumber(2518952)
  set flexibletimewindow(FlexibleTimeWindow value) =>
      $_setField(2518952, value);
  @$pb.TagNumber(2518952)
  $core.bool hasFlexibletimewindow() => $_has(0);
  @$pb.TagNumber(2518952)
  void clearFlexibletimewindow() => $_clearField(2518952);
  @$pb.TagNumber(2518952)
  FlexibleTimeWindow ensureFlexibletimewindow() => $_ensure(0);

  @$pb.TagNumber(77486543)
  $core.String get enddate => $_getSZ(1);
  @$pb.TagNumber(77486543)
  set enddate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(77486543)
  $core.bool hasEnddate() => $_has(1);
  @$pb.TagNumber(77486543)
  void clearEnddate() => $_clearField(77486543);

  @$pb.TagNumber(110041649)
  $core.String get kmskeyarn => $_getSZ(2);
  @$pb.TagNumber(110041649)
  set kmskeyarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(110041649)
  $core.bool hasKmskeyarn() => $_has(2);
  @$pb.TagNumber(110041649)
  void clearKmskeyarn() => $_clearField(110041649);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(137297356)
  $core.String get clienttoken => $_getSZ(4);
  @$pb.TagNumber(137297356)
  set clienttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(137297356)
  $core.bool hasClienttoken() => $_has(4);
  @$pb.TagNumber(137297356)
  void clearClienttoken() => $_clearField(137297356);

  @$pb.TagNumber(191361385)
  Target get target => $_getN(5);
  @$pb.TagNumber(191361385)
  set target(Target value) => $_setField(191361385, value);
  @$pb.TagNumber(191361385)
  $core.bool hasTarget() => $_has(5);
  @$pb.TagNumber(191361385)
  void clearTarget() => $_clearField(191361385);
  @$pb.TagNumber(191361385)
  Target ensureTarget() => $_ensure(5);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(282350906)
  $core.String get actionaftercompletion => $_getSZ(7);
  @$pb.TagNumber(282350906)
  set actionaftercompletion($core.String value) => $_setString(7, value);
  @$pb.TagNumber(282350906)
  $core.bool hasActionaftercompletion() => $_has(7);
  @$pb.TagNumber(282350906)
  void clearActionaftercompletion() => $_clearField(282350906);

  @$pb.TagNumber(357049672)
  $core.String get groupname => $_getSZ(8);
  @$pb.TagNumber(357049672)
  set groupname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(357049672)
  $core.bool hasGroupname() => $_has(8);
  @$pb.TagNumber(357049672)
  void clearGroupname() => $_clearField(357049672);

  @$pb.TagNumber(400730326)
  $core.String get scheduleexpressiontimezone => $_getSZ(9);
  @$pb.TagNumber(400730326)
  set scheduleexpressiontimezone($core.String value) => $_setString(9, value);
  @$pb.TagNumber(400730326)
  $core.bool hasScheduleexpressiontimezone() => $_has(9);
  @$pb.TagNumber(400730326)
  void clearScheduleexpressiontimezone() => $_clearField(400730326);

  @$pb.TagNumber(445135996)
  $core.String get startdate => $_getSZ(10);
  @$pb.TagNumber(445135996)
  set startdate($core.String value) => $_setString(10, value);
  @$pb.TagNumber(445135996)
  $core.bool hasStartdate() => $_has(10);
  @$pb.TagNumber(445135996)
  void clearStartdate() => $_clearField(445135996);

  @$pb.TagNumber(446089471)
  $core.String get scheduleexpression => $_getSZ(11);
  @$pb.TagNumber(446089471)
  set scheduleexpression($core.String value) => $_setString(11, value);
  @$pb.TagNumber(446089471)
  $core.bool hasScheduleexpression() => $_has(11);
  @$pb.TagNumber(446089471)
  void clearScheduleexpression() => $_clearField(446089471);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(12);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(12, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(12);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class UpdateScheduleOutput extends $pb.GeneratedMessage {
  factory UpdateScheduleOutput({
    $core.String? schedulearn,
  }) {
    final result = create();
    if (schedulearn != null) result.schedulearn = schedulearn;
    return result;
  }

  UpdateScheduleOutput._();

  factory UpdateScheduleOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateScheduleOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateScheduleOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
      createEmptyInstance: create)
    ..aOS(178445188, _omitFieldNames ? '' : 'schedulearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateScheduleOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateScheduleOutput copyWith(void Function(UpdateScheduleOutput) updates) =>
      super.copyWith((message) => updates(message as UpdateScheduleOutput))
          as UpdateScheduleOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateScheduleOutput create() => UpdateScheduleOutput._();
  @$core.override
  UpdateScheduleOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateScheduleOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateScheduleOutput>(create);
  static UpdateScheduleOutput? _defaultInstance;

  @$pb.TagNumber(178445188)
  $core.String get schedulearn => $_getSZ(0);
  @$pb.TagNumber(178445188)
  set schedulearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(178445188)
  $core.bool hasSchedulearn() => $_has(0);
  @$pb.TagNumber(178445188)
  void clearSchedulearn() => $_clearField(178445188);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'scheduler'),
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

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
