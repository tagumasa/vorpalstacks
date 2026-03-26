// This is a generated file - do not edit.
//
// Generated from route53.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'route53.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'route53.pbenum.dart';

class AccountLimit extends $pb.GeneratedMessage {
  factory AccountLimit({
    $fixnum.Int64? value,
    AccountLimitType? type,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (type != null) result.type = type;
    return result;
  }

  AccountLimit._();

  factory AccountLimit.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccountLimit.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountLimit',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(289929579, _omitFieldNames ? '' : 'value')
    ..aE<AccountLimitType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: AccountLimitType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountLimit clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountLimit copyWith(void Function(AccountLimit) updates) =>
      super.copyWith((message) => updates(message as AccountLimit))
          as AccountLimit;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccountLimit create() => AccountLimit._();
  @$core.override
  AccountLimit createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccountLimit getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountLimit>(create);
  static AccountLimit? _defaultInstance;

  @$pb.TagNumber(289929579)
  $fixnum.Int64 get value => $_getI64(0);
  @$pb.TagNumber(289929579)
  set value($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(290836590)
  AccountLimitType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(AccountLimitType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class ActivateKeySigningKeyRequest extends $pb.GeneratedMessage {
  factory ActivateKeySigningKeyRequest({
    $core.String? name,
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  ActivateKeySigningKeyRequest._();

  factory ActivateKeySigningKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivateKeySigningKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivateKeySigningKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivateKeySigningKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivateKeySigningKeyRequest copyWith(
          void Function(ActivateKeySigningKeyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ActivateKeySigningKeyRequest))
          as ActivateKeySigningKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivateKeySigningKeyRequest create() =>
      ActivateKeySigningKeyRequest._();
  @$core.override
  ActivateKeySigningKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivateKeySigningKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivateKeySigningKeyRequest>(create);
  static ActivateKeySigningKeyRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class ActivateKeySigningKeyResponse extends $pb.GeneratedMessage {
  factory ActivateKeySigningKeyResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  ActivateKeySigningKeyResponse._();

  factory ActivateKeySigningKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivateKeySigningKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivateKeySigningKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivateKeySigningKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivateKeySigningKeyResponse copyWith(
          void Function(ActivateKeySigningKeyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ActivateKeySigningKeyResponse))
          as ActivateKeySigningKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivateKeySigningKeyResponse create() =>
      ActivateKeySigningKeyResponse._();
  @$core.override
  ActivateKeySigningKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivateKeySigningKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivateKeySigningKeyResponse>(create);
  static ActivateKeySigningKeyResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class AlarmIdentifier extends $pb.GeneratedMessage {
  factory AlarmIdentifier({
    CloudWatchRegion? region,
    $core.String? name,
  }) {
    final result = create();
    if (region != null) result.region = region;
    if (name != null) result.name = name;
    return result;
  }

  AlarmIdentifier._();

  factory AlarmIdentifier.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AlarmIdentifier.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AlarmIdentifier',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<CloudWatchRegion>(154040478, _omitFieldNames ? '' : 'region',
        enumValues: CloudWatchRegion.values)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlarmIdentifier clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlarmIdentifier copyWith(void Function(AlarmIdentifier) updates) =>
      super.copyWith((message) => updates(message as AlarmIdentifier))
          as AlarmIdentifier;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AlarmIdentifier create() => AlarmIdentifier._();
  @$core.override
  AlarmIdentifier createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AlarmIdentifier getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AlarmIdentifier>(create);
  static AlarmIdentifier? _defaultInstance;

  @$pb.TagNumber(154040478)
  CloudWatchRegion get region => $_getN(0);
  @$pb.TagNumber(154040478)
  set region(CloudWatchRegion value) => $_setField(154040478, value);
  @$pb.TagNumber(154040478)
  $core.bool hasRegion() => $_has(0);
  @$pb.TagNumber(154040478)
  void clearRegion() => $_clearField(154040478);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class AliasTarget extends $pb.GeneratedMessage {
  factory AliasTarget({
    $core.String? dnsname,
    $core.String? hostedzoneid,
    $core.bool? evaluatetargethealth,
  }) {
    final result = create();
    if (dnsname != null) result.dnsname = dnsname;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (evaluatetargethealth != null)
      result.evaluatetargethealth = evaluatetargethealth;
    return result;
  }

  AliasTarget._();

  factory AliasTarget.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AliasTarget.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AliasTarget',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(171901432, _omitFieldNames ? '' : 'dnsname')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOB(409666542, _omitFieldNames ? '' : 'evaluatetargethealth')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AliasTarget clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AliasTarget copyWith(void Function(AliasTarget) updates) =>
      super.copyWith((message) => updates(message as AliasTarget))
          as AliasTarget;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AliasTarget create() => AliasTarget._();
  @$core.override
  AliasTarget createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AliasTarget getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AliasTarget>(create);
  static AliasTarget? _defaultInstance;

  @$pb.TagNumber(171901432)
  $core.String get dnsname => $_getSZ(0);
  @$pb.TagNumber(171901432)
  set dnsname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(171901432)
  $core.bool hasDnsname() => $_has(0);
  @$pb.TagNumber(171901432)
  void clearDnsname() => $_clearField(171901432);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(409666542)
  $core.bool get evaluatetargethealth => $_getBF(2);
  @$pb.TagNumber(409666542)
  set evaluatetargethealth($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(409666542)
  $core.bool hasEvaluatetargethealth() => $_has(2);
  @$pb.TagNumber(409666542)
  void clearEvaluatetargethealth() => $_clearField(409666542);
}

class AssociateVPCWithHostedZoneRequest extends $pb.GeneratedMessage {
  factory AssociateVPCWithHostedZoneRequest({
    $core.String? hostedzoneid,
    $core.String? comment,
    VPC? vpc,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (comment != null) result.comment = comment;
    if (vpc != null) result.vpc = vpc;
    return result;
  }

  AssociateVPCWithHostedZoneRequest._();

  factory AssociateVPCWithHostedZoneRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssociateVPCWithHostedZoneRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssociateVPCWithHostedZoneRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..aOM<VPC>(506158953, _omitFieldNames ? '' : 'vpc', subBuilder: VPC.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociateVPCWithHostedZoneRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociateVPCWithHostedZoneRequest copyWith(
          void Function(AssociateVPCWithHostedZoneRequest) updates) =>
      super.copyWith((message) =>
              updates(message as AssociateVPCWithHostedZoneRequest))
          as AssociateVPCWithHostedZoneRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssociateVPCWithHostedZoneRequest create() =>
      AssociateVPCWithHostedZoneRequest._();
  @$core.override
  AssociateVPCWithHostedZoneRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssociateVPCWithHostedZoneRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssociateVPCWithHostedZoneRequest>(
          create);
  static AssociateVPCWithHostedZoneRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(1);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(1, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(1);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);

  @$pb.TagNumber(506158953)
  VPC get vpc => $_getN(2);
  @$pb.TagNumber(506158953)
  set vpc(VPC value) => $_setField(506158953, value);
  @$pb.TagNumber(506158953)
  $core.bool hasVpc() => $_has(2);
  @$pb.TagNumber(506158953)
  void clearVpc() => $_clearField(506158953);
  @$pb.TagNumber(506158953)
  VPC ensureVpc() => $_ensure(2);
}

class AssociateVPCWithHostedZoneResponse extends $pb.GeneratedMessage {
  factory AssociateVPCWithHostedZoneResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  AssociateVPCWithHostedZoneResponse._();

  factory AssociateVPCWithHostedZoneResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AssociateVPCWithHostedZoneResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AssociateVPCWithHostedZoneResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociateVPCWithHostedZoneResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AssociateVPCWithHostedZoneResponse copyWith(
          void Function(AssociateVPCWithHostedZoneResponse) updates) =>
      super.copyWith((message) =>
              updates(message as AssociateVPCWithHostedZoneResponse))
          as AssociateVPCWithHostedZoneResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AssociateVPCWithHostedZoneResponse create() =>
      AssociateVPCWithHostedZoneResponse._();
  @$core.override
  AssociateVPCWithHostedZoneResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AssociateVPCWithHostedZoneResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AssociateVPCWithHostedZoneResponse>(
          create);
  static AssociateVPCWithHostedZoneResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class Change extends $pb.GeneratedMessage {
  factory Change({
    ChangeAction? action,
    ResourceRecordSet? resourcerecordset,
  }) {
    final result = create();
    if (action != null) result.action = action;
    if (resourcerecordset != null) result.resourcerecordset = resourcerecordset;
    return result;
  }

  Change._();

  factory Change.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Change.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Change',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<ChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: ChangeAction.values)
    ..aOM<ResourceRecordSet>(
        344937781, _omitFieldNames ? '' : 'resourcerecordset',
        subBuilder: ResourceRecordSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Change clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Change copyWith(void Function(Change) updates) =>
      super.copyWith((message) => updates(message as Change)) as Change;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Change create() => Change._();
  @$core.override
  Change createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Change getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Change>(create);
  static Change? _defaultInstance;

  @$pb.TagNumber(175614240)
  ChangeAction get action => $_getN(0);
  @$pb.TagNumber(175614240)
  set action(ChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(0);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(344937781)
  ResourceRecordSet get resourcerecordset => $_getN(1);
  @$pb.TagNumber(344937781)
  set resourcerecordset(ResourceRecordSet value) =>
      $_setField(344937781, value);
  @$pb.TagNumber(344937781)
  $core.bool hasResourcerecordset() => $_has(1);
  @$pb.TagNumber(344937781)
  void clearResourcerecordset() => $_clearField(344937781);
  @$pb.TagNumber(344937781)
  ResourceRecordSet ensureResourcerecordset() => $_ensure(1);
}

class ChangeBatch extends $pb.GeneratedMessage {
  factory ChangeBatch({
    $core.String? comment,
    $core.Iterable<Change>? changes,
  }) {
    final result = create();
    if (comment != null) result.comment = comment;
    if (changes != null) result.changes.addAll(changes);
    return result;
  }

  ChangeBatch._();

  factory ChangeBatch.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeBatch.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeBatch',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..pPM<Change>(516230891, _omitFieldNames ? '' : 'changes',
        subBuilder: Change.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeBatch clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeBatch copyWith(void Function(ChangeBatch) updates) =>
      super.copyWith((message) => updates(message as ChangeBatch))
          as ChangeBatch;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeBatch create() => ChangeBatch._();
  @$core.override
  ChangeBatch createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeBatch getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeBatch>(create);
  static ChangeBatch? _defaultInstance;

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(0);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(0, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(0);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);

  @$pb.TagNumber(516230891)
  $pb.PbList<Change> get changes => $_getList(1);
}

class ChangeCidrCollectionRequest extends $pb.GeneratedMessage {
  factory ChangeCidrCollectionRequest({
    $fixnum.Int64? collectionversion,
    $core.String? id,
    $core.Iterable<CidrCollectionChange>? changes,
  }) {
    final result = create();
    if (collectionversion != null) result.collectionversion = collectionversion;
    if (id != null) result.id = id;
    if (changes != null) result.changes.addAll(changes);
    return result;
  }

  ChangeCidrCollectionRequest._();

  factory ChangeCidrCollectionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeCidrCollectionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeCidrCollectionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(95829948, _omitFieldNames ? '' : 'collectionversion')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..pPM<CidrCollectionChange>(516230891, _omitFieldNames ? '' : 'changes',
        subBuilder: CidrCollectionChange.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeCidrCollectionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeCidrCollectionRequest copyWith(
          void Function(ChangeCidrCollectionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ChangeCidrCollectionRequest))
          as ChangeCidrCollectionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeCidrCollectionRequest create() =>
      ChangeCidrCollectionRequest._();
  @$core.override
  ChangeCidrCollectionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeCidrCollectionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeCidrCollectionRequest>(create);
  static ChangeCidrCollectionRequest? _defaultInstance;

  @$pb.TagNumber(95829948)
  $fixnum.Int64 get collectionversion => $_getI64(0);
  @$pb.TagNumber(95829948)
  set collectionversion($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(95829948)
  $core.bool hasCollectionversion() => $_has(0);
  @$pb.TagNumber(95829948)
  void clearCollectionversion() => $_clearField(95829948);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(516230891)
  $pb.PbList<CidrCollectionChange> get changes => $_getList(2);
}

class ChangeCidrCollectionResponse extends $pb.GeneratedMessage {
  factory ChangeCidrCollectionResponse({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  ChangeCidrCollectionResponse._();

  factory ChangeCidrCollectionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeCidrCollectionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeCidrCollectionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeCidrCollectionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeCidrCollectionResponse copyWith(
          void Function(ChangeCidrCollectionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ChangeCidrCollectionResponse))
          as ChangeCidrCollectionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeCidrCollectionResponse create() =>
      ChangeCidrCollectionResponse._();
  @$core.override
  ChangeCidrCollectionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeCidrCollectionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeCidrCollectionResponse>(create);
  static ChangeCidrCollectionResponse? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class ChangeInfo extends $pb.GeneratedMessage {
  factory ChangeInfo({
    ChangeStatus? status,
    $core.String? submittedat,
    $core.String? id,
    $core.String? comment,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (submittedat != null) result.submittedat = submittedat;
    if (id != null) result.id = id;
    if (comment != null) result.comment = comment;
    return result;
  }

  ChangeInfo._();

  factory ChangeInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<ChangeStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: ChangeStatus.values)
    ..aOS(343958936, _omitFieldNames ? '' : 'submittedat')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeInfo copyWith(void Function(ChangeInfo) updates) =>
      super.copyWith((message) => updates(message as ChangeInfo)) as ChangeInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeInfo create() => ChangeInfo._();
  @$core.override
  ChangeInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeInfo>(create);
  static ChangeInfo? _defaultInstance;

  @$pb.TagNumber(6222352)
  ChangeStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(ChangeStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(343958936)
  $core.String get submittedat => $_getSZ(1);
  @$pb.TagNumber(343958936)
  set submittedat($core.String value) => $_setString(1, value);
  @$pb.TagNumber(343958936)
  $core.bool hasSubmittedat() => $_has(1);
  @$pb.TagNumber(343958936)
  void clearSubmittedat() => $_clearField(343958936);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(3);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(3, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(3);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);
}

class ChangeResourceRecordSetsRequest extends $pb.GeneratedMessage {
  factory ChangeResourceRecordSetsRequest({
    ChangeBatch? changebatch,
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (changebatch != null) result.changebatch = changebatch;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  ChangeResourceRecordSetsRequest._();

  factory ChangeResourceRecordSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeResourceRecordSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeResourceRecordSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeBatch>(23600436, _omitFieldNames ? '' : 'changebatch',
        subBuilder: ChangeBatch.create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeResourceRecordSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeResourceRecordSetsRequest copyWith(
          void Function(ChangeResourceRecordSetsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ChangeResourceRecordSetsRequest))
          as ChangeResourceRecordSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeResourceRecordSetsRequest create() =>
      ChangeResourceRecordSetsRequest._();
  @$core.override
  ChangeResourceRecordSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeResourceRecordSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeResourceRecordSetsRequest>(
          create);
  static ChangeResourceRecordSetsRequest? _defaultInstance;

  @$pb.TagNumber(23600436)
  ChangeBatch get changebatch => $_getN(0);
  @$pb.TagNumber(23600436)
  set changebatch(ChangeBatch value) => $_setField(23600436, value);
  @$pb.TagNumber(23600436)
  $core.bool hasChangebatch() => $_has(0);
  @$pb.TagNumber(23600436)
  void clearChangebatch() => $_clearField(23600436);
  @$pb.TagNumber(23600436)
  ChangeBatch ensureChangebatch() => $_ensure(0);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class ChangeResourceRecordSetsResponse extends $pb.GeneratedMessage {
  factory ChangeResourceRecordSetsResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  ChangeResourceRecordSetsResponse._();

  factory ChangeResourceRecordSetsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeResourceRecordSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeResourceRecordSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeResourceRecordSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeResourceRecordSetsResponse copyWith(
          void Function(ChangeResourceRecordSetsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ChangeResourceRecordSetsResponse))
          as ChangeResourceRecordSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeResourceRecordSetsResponse create() =>
      ChangeResourceRecordSetsResponse._();
  @$core.override
  ChangeResourceRecordSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeResourceRecordSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeResourceRecordSetsResponse>(
          create);
  static ChangeResourceRecordSetsResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class ChangeTagsForResourceRequest extends $pb.GeneratedMessage {
  factory ChangeTagsForResourceRequest({
    $core.Iterable<$core.String>? removetagkeys,
    TagResourceType? resourcetype,
    $core.Iterable<Tag>? addtags,
    $core.String? resourceid,
  }) {
    final result = create();
    if (removetagkeys != null) result.removetagkeys.addAll(removetagkeys);
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (addtags != null) result.addtags.addAll(addtags);
    if (resourceid != null) result.resourceid = resourceid;
    return result;
  }

  ChangeTagsForResourceRequest._();

  factory ChangeTagsForResourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeTagsForResourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeTagsForResourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPS(128174116, _omitFieldNames ? '' : 'removetagkeys')
    ..aE<TagResourceType>(301342558, _omitFieldNames ? '' : 'resourcetype',
        enumValues: TagResourceType.values)
    ..pPM<Tag>(503857054, _omitFieldNames ? '' : 'addtags',
        subBuilder: Tag.create)
    ..aOS(526146833, _omitFieldNames ? '' : 'resourceid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeTagsForResourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeTagsForResourceRequest copyWith(
          void Function(ChangeTagsForResourceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ChangeTagsForResourceRequest))
          as ChangeTagsForResourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeTagsForResourceRequest create() =>
      ChangeTagsForResourceRequest._();
  @$core.override
  ChangeTagsForResourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeTagsForResourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeTagsForResourceRequest>(create);
  static ChangeTagsForResourceRequest? _defaultInstance;

  @$pb.TagNumber(128174116)
  $pb.PbList<$core.String> get removetagkeys => $_getList(0);

  @$pb.TagNumber(301342558)
  TagResourceType get resourcetype => $_getN(1);
  @$pb.TagNumber(301342558)
  set resourcetype(TagResourceType value) => $_setField(301342558, value);
  @$pb.TagNumber(301342558)
  $core.bool hasResourcetype() => $_has(1);
  @$pb.TagNumber(301342558)
  void clearResourcetype() => $_clearField(301342558);

  @$pb.TagNumber(503857054)
  $pb.PbList<Tag> get addtags => $_getList(2);

  @$pb.TagNumber(526146833)
  $core.String get resourceid => $_getSZ(3);
  @$pb.TagNumber(526146833)
  set resourceid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(526146833)
  $core.bool hasResourceid() => $_has(3);
  @$pb.TagNumber(526146833)
  void clearResourceid() => $_clearField(526146833);
}

class ChangeTagsForResourceResponse extends $pb.GeneratedMessage {
  factory ChangeTagsForResourceResponse() => create();

  ChangeTagsForResourceResponse._();

  factory ChangeTagsForResourceResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ChangeTagsForResourceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ChangeTagsForResourceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeTagsForResourceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ChangeTagsForResourceResponse copyWith(
          void Function(ChangeTagsForResourceResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ChangeTagsForResourceResponse))
          as ChangeTagsForResourceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ChangeTagsForResourceResponse create() =>
      ChangeTagsForResourceResponse._();
  @$core.override
  ChangeTagsForResourceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ChangeTagsForResourceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ChangeTagsForResourceResponse>(create);
  static ChangeTagsForResourceResponse? _defaultInstance;
}

class CidrBlockInUseException extends $pb.GeneratedMessage {
  factory CidrBlockInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CidrBlockInUseException._();

  factory CidrBlockInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CidrBlockInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CidrBlockInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrBlockInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrBlockInUseException copyWith(
          void Function(CidrBlockInUseException) updates) =>
      super.copyWith((message) => updates(message as CidrBlockInUseException))
          as CidrBlockInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CidrBlockInUseException create() => CidrBlockInUseException._();
  @$core.override
  CidrBlockInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CidrBlockInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CidrBlockInUseException>(create);
  static CidrBlockInUseException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CidrBlockSummary extends $pb.GeneratedMessage {
  factory CidrBlockSummary({
    $core.String? cidrblock,
    $core.String? locationname,
  }) {
    final result = create();
    if (cidrblock != null) result.cidrblock = cidrblock;
    if (locationname != null) result.locationname = locationname;
    return result;
  }

  CidrBlockSummary._();

  factory CidrBlockSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CidrBlockSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CidrBlockSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(56494185, _omitFieldNames ? '' : 'cidrblock')
    ..aOS(158186566, _omitFieldNames ? '' : 'locationname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrBlockSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrBlockSummary copyWith(void Function(CidrBlockSummary) updates) =>
      super.copyWith((message) => updates(message as CidrBlockSummary))
          as CidrBlockSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CidrBlockSummary create() => CidrBlockSummary._();
  @$core.override
  CidrBlockSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CidrBlockSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CidrBlockSummary>(create);
  static CidrBlockSummary? _defaultInstance;

  @$pb.TagNumber(56494185)
  $core.String get cidrblock => $_getSZ(0);
  @$pb.TagNumber(56494185)
  set cidrblock($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56494185)
  $core.bool hasCidrblock() => $_has(0);
  @$pb.TagNumber(56494185)
  void clearCidrblock() => $_clearField(56494185);

  @$pb.TagNumber(158186566)
  $core.String get locationname => $_getSZ(1);
  @$pb.TagNumber(158186566)
  set locationname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(158186566)
  $core.bool hasLocationname() => $_has(1);
  @$pb.TagNumber(158186566)
  void clearLocationname() => $_clearField(158186566);
}

class CidrCollection extends $pb.GeneratedMessage {
  factory CidrCollection({
    $core.String? name,
    $core.String? id,
    $core.String? arn,
    $fixnum.Int64? version,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    if (version != null) result.version = version;
    return result;
  }

  CidrCollection._();

  factory CidrCollection.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CidrCollection.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CidrCollection',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aInt64(500028728, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollection clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollection copyWith(void Function(CidrCollection) updates) =>
      super.copyWith((message) => updates(message as CidrCollection))
          as CidrCollection;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CidrCollection create() => CidrCollection._();
  @$core.override
  CidrCollection createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CidrCollection getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CidrCollection>(create);
  static CidrCollection? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(500028728)
  $fixnum.Int64 get version => $_getI64(3);
  @$pb.TagNumber(500028728)
  set version($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(3);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);
}

class CidrCollectionAlreadyExistsException extends $pb.GeneratedMessage {
  factory CidrCollectionAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CidrCollectionAlreadyExistsException._();

  factory CidrCollectionAlreadyExistsException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CidrCollectionAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CidrCollectionAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollectionAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollectionAlreadyExistsException copyWith(
          void Function(CidrCollectionAlreadyExistsException) updates) =>
      super.copyWith((message) =>
              updates(message as CidrCollectionAlreadyExistsException))
          as CidrCollectionAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CidrCollectionAlreadyExistsException create() =>
      CidrCollectionAlreadyExistsException._();
  @$core.override
  CidrCollectionAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CidrCollectionAlreadyExistsException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CidrCollectionAlreadyExistsException>(create);
  static CidrCollectionAlreadyExistsException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CidrCollectionChange extends $pb.GeneratedMessage {
  factory CidrCollectionChange({
    $core.String? locationname,
    CidrCollectionChangeAction? action,
    $core.Iterable<$core.String>? cidrlist,
  }) {
    final result = create();
    if (locationname != null) result.locationname = locationname;
    if (action != null) result.action = action;
    if (cidrlist != null) result.cidrlist.addAll(cidrlist);
    return result;
  }

  CidrCollectionChange._();

  factory CidrCollectionChange.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CidrCollectionChange.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CidrCollectionChange',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(158186566, _omitFieldNames ? '' : 'locationname')
    ..aE<CidrCollectionChangeAction>(175614240, _omitFieldNames ? '' : 'action',
        enumValues: CidrCollectionChangeAction.values)
    ..pPS(288964704, _omitFieldNames ? '' : 'cidrlist')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollectionChange clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollectionChange copyWith(void Function(CidrCollectionChange) updates) =>
      super.copyWith((message) => updates(message as CidrCollectionChange))
          as CidrCollectionChange;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CidrCollectionChange create() => CidrCollectionChange._();
  @$core.override
  CidrCollectionChange createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CidrCollectionChange getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CidrCollectionChange>(create);
  static CidrCollectionChange? _defaultInstance;

  @$pb.TagNumber(158186566)
  $core.String get locationname => $_getSZ(0);
  @$pb.TagNumber(158186566)
  set locationname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(158186566)
  $core.bool hasLocationname() => $_has(0);
  @$pb.TagNumber(158186566)
  void clearLocationname() => $_clearField(158186566);

  @$pb.TagNumber(175614240)
  CidrCollectionChangeAction get action => $_getN(1);
  @$pb.TagNumber(175614240)
  set action(CidrCollectionChangeAction value) => $_setField(175614240, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(1);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(288964704)
  $pb.PbList<$core.String> get cidrlist => $_getList(2);
}

class CidrCollectionInUseException extends $pb.GeneratedMessage {
  factory CidrCollectionInUseException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CidrCollectionInUseException._();

  factory CidrCollectionInUseException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CidrCollectionInUseException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CidrCollectionInUseException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollectionInUseException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollectionInUseException copyWith(
          void Function(CidrCollectionInUseException) updates) =>
      super.copyWith(
              (message) => updates(message as CidrCollectionInUseException))
          as CidrCollectionInUseException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CidrCollectionInUseException create() =>
      CidrCollectionInUseException._();
  @$core.override
  CidrCollectionInUseException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CidrCollectionInUseException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CidrCollectionInUseException>(create);
  static CidrCollectionInUseException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CidrCollectionVersionMismatchException extends $pb.GeneratedMessage {
  factory CidrCollectionVersionMismatchException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  CidrCollectionVersionMismatchException._();

  factory CidrCollectionVersionMismatchException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CidrCollectionVersionMismatchException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CidrCollectionVersionMismatchException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollectionVersionMismatchException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrCollectionVersionMismatchException copyWith(
          void Function(CidrCollectionVersionMismatchException) updates) =>
      super.copyWith((message) =>
              updates(message as CidrCollectionVersionMismatchException))
          as CidrCollectionVersionMismatchException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CidrCollectionVersionMismatchException create() =>
      CidrCollectionVersionMismatchException._();
  @$core.override
  CidrCollectionVersionMismatchException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CidrCollectionVersionMismatchException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CidrCollectionVersionMismatchException>(create);
  static CidrCollectionVersionMismatchException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class CidrRoutingConfig extends $pb.GeneratedMessage {
  factory CidrRoutingConfig({
    $core.String? collectionid,
    $core.String? locationname,
  }) {
    final result = create();
    if (collectionid != null) result.collectionid = collectionid;
    if (locationname != null) result.locationname = locationname;
    return result;
  }

  CidrRoutingConfig._();

  factory CidrRoutingConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CidrRoutingConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CidrRoutingConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(128052453, _omitFieldNames ? '' : 'collectionid')
    ..aOS(158186566, _omitFieldNames ? '' : 'locationname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrRoutingConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CidrRoutingConfig copyWith(void Function(CidrRoutingConfig) updates) =>
      super.copyWith((message) => updates(message as CidrRoutingConfig))
          as CidrRoutingConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CidrRoutingConfig create() => CidrRoutingConfig._();
  @$core.override
  CidrRoutingConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CidrRoutingConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CidrRoutingConfig>(create);
  static CidrRoutingConfig? _defaultInstance;

  @$pb.TagNumber(128052453)
  $core.String get collectionid => $_getSZ(0);
  @$pb.TagNumber(128052453)
  set collectionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(128052453)
  $core.bool hasCollectionid() => $_has(0);
  @$pb.TagNumber(128052453)
  void clearCollectionid() => $_clearField(128052453);

  @$pb.TagNumber(158186566)
  $core.String get locationname => $_getSZ(1);
  @$pb.TagNumber(158186566)
  set locationname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(158186566)
  $core.bool hasLocationname() => $_has(1);
  @$pb.TagNumber(158186566)
  void clearLocationname() => $_clearField(158186566);
}

class CloudWatchAlarmConfiguration extends $pb.GeneratedMessage {
  factory CloudWatchAlarmConfiguration({
    ComparisonOperator? comparisonoperator,
    Statistic? statistic,
    $core.String? metricname,
    $core.int? period,
    $core.int? evaluationperiods,
    $core.String? namespace,
    $core.double? threshold,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (comparisonoperator != null)
      result.comparisonoperator = comparisonoperator;
    if (statistic != null) result.statistic = statistic;
    if (metricname != null) result.metricname = metricname;
    if (period != null) result.period = period;
    if (evaluationperiods != null) result.evaluationperiods = evaluationperiods;
    if (namespace != null) result.namespace = namespace;
    if (threshold != null) result.threshold = threshold;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  CloudWatchAlarmConfiguration._();

  factory CloudWatchAlarmConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CloudWatchAlarmConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CloudWatchAlarmConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<ComparisonOperator>(
        7059861, _omitFieldNames ? '' : 'comparisonoperator',
        enumValues: ComparisonOperator.values)
    ..aE<Statistic>(67293470, _omitFieldNames ? '' : 'statistic',
        enumValues: Statistic.values)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aI(214794856, _omitFieldNames ? '' : 'evaluationperiods')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..aD(394884505, _omitFieldNames ? '' : 'threshold')
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchAlarmConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CloudWatchAlarmConfiguration copyWith(
          void Function(CloudWatchAlarmConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as CloudWatchAlarmConfiguration))
          as CloudWatchAlarmConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CloudWatchAlarmConfiguration create() =>
      CloudWatchAlarmConfiguration._();
  @$core.override
  CloudWatchAlarmConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CloudWatchAlarmConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CloudWatchAlarmConfiguration>(create);
  static CloudWatchAlarmConfiguration? _defaultInstance;

  @$pb.TagNumber(7059861)
  ComparisonOperator get comparisonoperator => $_getN(0);
  @$pb.TagNumber(7059861)
  set comparisonoperator(ComparisonOperator value) =>
      $_setField(7059861, value);
  @$pb.TagNumber(7059861)
  $core.bool hasComparisonoperator() => $_has(0);
  @$pb.TagNumber(7059861)
  void clearComparisonoperator() => $_clearField(7059861);

  @$pb.TagNumber(67293470)
  Statistic get statistic => $_getN(1);
  @$pb.TagNumber(67293470)
  set statistic(Statistic value) => $_setField(67293470, value);
  @$pb.TagNumber(67293470)
  $core.bool hasStatistic() => $_has(1);
  @$pb.TagNumber(67293470)
  void clearStatistic() => $_clearField(67293470);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(2);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(2);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(3);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(3);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(214794856)
  $core.int get evaluationperiods => $_getIZ(4);
  @$pb.TagNumber(214794856)
  set evaluationperiods($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(214794856)
  $core.bool hasEvaluationperiods() => $_has(4);
  @$pb.TagNumber(214794856)
  void clearEvaluationperiods() => $_clearField(214794856);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(5);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(5, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(5);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(394884505)
  $core.double get threshold => $_getN(6);
  @$pb.TagNumber(394884505)
  set threshold($core.double value) => $_setDouble(6, value);
  @$pb.TagNumber(394884505)
  $core.bool hasThreshold() => $_has(6);
  @$pb.TagNumber(394884505)
  void clearThreshold() => $_clearField(394884505);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(7);
}

class CollectionSummary extends $pb.GeneratedMessage {
  factory CollectionSummary({
    $core.String? name,
    $core.String? id,
    $core.String? arn,
    $fixnum.Int64? version,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (id != null) result.id = id;
    if (arn != null) result.arn = arn;
    if (version != null) result.version = version;
    return result;
  }

  CollectionSummary._();

  factory CollectionSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CollectionSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CollectionSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aInt64(500028728, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CollectionSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CollectionSummary copyWith(void Function(CollectionSummary) updates) =>
      super.copyWith((message) => updates(message as CollectionSummary))
          as CollectionSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CollectionSummary create() => CollectionSummary._();
  @$core.override
  CollectionSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CollectionSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CollectionSummary>(create);
  static CollectionSummary? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(2);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(2);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(500028728)
  $fixnum.Int64 get version => $_getI64(3);
  @$pb.TagNumber(500028728)
  set version($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(3);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);
}

class ConcurrentModification extends $pb.GeneratedMessage {
  factory ConcurrentModification({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ConcurrentModification._();

  factory ConcurrentModification.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConcurrentModification.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConcurrentModification',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConcurrentModification clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConcurrentModification copyWith(
          void Function(ConcurrentModification) updates) =>
      super.copyWith((message) => updates(message as ConcurrentModification))
          as ConcurrentModification;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConcurrentModification create() => ConcurrentModification._();
  @$core.override
  ConcurrentModification createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConcurrentModification getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConcurrentModification>(create);
  static ConcurrentModification? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ConflictingDomainExists extends $pb.GeneratedMessage {
  factory ConflictingDomainExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ConflictingDomainExists._();

  factory ConflictingDomainExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConflictingDomainExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConflictingDomainExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConflictingDomainExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConflictingDomainExists copyWith(
          void Function(ConflictingDomainExists) updates) =>
      super.copyWith((message) => updates(message as ConflictingDomainExists))
          as ConflictingDomainExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConflictingDomainExists create() => ConflictingDomainExists._();
  @$core.override
  ConflictingDomainExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConflictingDomainExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConflictingDomainExists>(create);
  static ConflictingDomainExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ConflictingTypes extends $pb.GeneratedMessage {
  factory ConflictingTypes({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ConflictingTypes._();

  factory ConflictingTypes.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConflictingTypes.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConflictingTypes',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConflictingTypes clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConflictingTypes copyWith(void Function(ConflictingTypes) updates) =>
      super.copyWith((message) => updates(message as ConflictingTypes))
          as ConflictingTypes;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConflictingTypes create() => ConflictingTypes._();
  @$core.override
  ConflictingTypes createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConflictingTypes getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConflictingTypes>(create);
  static ConflictingTypes? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class Coordinates extends $pb.GeneratedMessage {
  factory Coordinates({
    $core.String? longitude,
    $core.String? latitude,
  }) {
    final result = create();
    if (longitude != null) result.longitude = longitude;
    if (latitude != null) result.latitude = latitude;
    return result;
  }

  Coordinates._();

  factory Coordinates.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Coordinates.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Coordinates',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(119104751, _omitFieldNames ? '' : 'longitude')
    ..aOS(226378086, _omitFieldNames ? '' : 'latitude')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Coordinates clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Coordinates copyWith(void Function(Coordinates) updates) =>
      super.copyWith((message) => updates(message as Coordinates))
          as Coordinates;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Coordinates create() => Coordinates._();
  @$core.override
  Coordinates createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Coordinates getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Coordinates>(create);
  static Coordinates? _defaultInstance;

  @$pb.TagNumber(119104751)
  $core.String get longitude => $_getSZ(0);
  @$pb.TagNumber(119104751)
  set longitude($core.String value) => $_setString(0, value);
  @$pb.TagNumber(119104751)
  $core.bool hasLongitude() => $_has(0);
  @$pb.TagNumber(119104751)
  void clearLongitude() => $_clearField(119104751);

  @$pb.TagNumber(226378086)
  $core.String get latitude => $_getSZ(1);
  @$pb.TagNumber(226378086)
  set latitude($core.String value) => $_setString(1, value);
  @$pb.TagNumber(226378086)
  $core.bool hasLatitude() => $_has(1);
  @$pb.TagNumber(226378086)
  void clearLatitude() => $_clearField(226378086);
}

class CreateCidrCollectionRequest extends $pb.GeneratedMessage {
  factory CreateCidrCollectionRequest({
    $core.String? callerreference,
    $core.String? name,
  }) {
    final result = create();
    if (callerreference != null) result.callerreference = callerreference;
    if (name != null) result.name = name;
    return result;
  }

  CreateCidrCollectionRequest._();

  factory CreateCidrCollectionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateCidrCollectionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateCidrCollectionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(151211160, _omitFieldNames ? '' : 'callerreference')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCidrCollectionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCidrCollectionRequest copyWith(
          void Function(CreateCidrCollectionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateCidrCollectionRequest))
          as CreateCidrCollectionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateCidrCollectionRequest create() =>
      CreateCidrCollectionRequest._();
  @$core.override
  CreateCidrCollectionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateCidrCollectionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateCidrCollectionRequest>(create);
  static CreateCidrCollectionRequest? _defaultInstance;

  @$pb.TagNumber(151211160)
  $core.String get callerreference => $_getSZ(0);
  @$pb.TagNumber(151211160)
  set callerreference($core.String value) => $_setString(0, value);
  @$pb.TagNumber(151211160)
  $core.bool hasCallerreference() => $_has(0);
  @$pb.TagNumber(151211160)
  void clearCallerreference() => $_clearField(151211160);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class CreateCidrCollectionResponse extends $pb.GeneratedMessage {
  factory CreateCidrCollectionResponse({
    CidrCollection? collection,
    $core.String? location,
  }) {
    final result = create();
    if (collection != null) result.collection = collection;
    if (location != null) result.location = location;
    return result;
  }

  CreateCidrCollectionResponse._();

  factory CreateCidrCollectionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateCidrCollectionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateCidrCollectionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<CidrCollection>(114272434, _omitFieldNames ? '' : 'collection',
        subBuilder: CidrCollection.create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCidrCollectionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateCidrCollectionResponse copyWith(
          void Function(CreateCidrCollectionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateCidrCollectionResponse))
          as CreateCidrCollectionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateCidrCollectionResponse create() =>
      CreateCidrCollectionResponse._();
  @$core.override
  CreateCidrCollectionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateCidrCollectionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateCidrCollectionResponse>(create);
  static CreateCidrCollectionResponse? _defaultInstance;

  @$pb.TagNumber(114272434)
  CidrCollection get collection => $_getN(0);
  @$pb.TagNumber(114272434)
  set collection(CidrCollection value) => $_setField(114272434, value);
  @$pb.TagNumber(114272434)
  $core.bool hasCollection() => $_has(0);
  @$pb.TagNumber(114272434)
  void clearCollection() => $_clearField(114272434);
  @$pb.TagNumber(114272434)
  CidrCollection ensureCollection() => $_ensure(0);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(1);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(1, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(1);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);
}

class CreateHealthCheckRequest extends $pb.GeneratedMessage {
  factory CreateHealthCheckRequest({
    HealthCheckConfig? healthcheckconfig,
    $core.String? callerreference,
  }) {
    final result = create();
    if (healthcheckconfig != null) result.healthcheckconfig = healthcheckconfig;
    if (callerreference != null) result.callerreference = callerreference;
    return result;
  }

  CreateHealthCheckRequest._();

  factory CreateHealthCheckRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateHealthCheckRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateHealthCheckRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<HealthCheckConfig>(
        83111204, _omitFieldNames ? '' : 'healthcheckconfig',
        subBuilder: HealthCheckConfig.create)
    ..aOS(151211160, _omitFieldNames ? '' : 'callerreference')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateHealthCheckRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateHealthCheckRequest copyWith(
          void Function(CreateHealthCheckRequest) updates) =>
      super.copyWith((message) => updates(message as CreateHealthCheckRequest))
          as CreateHealthCheckRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateHealthCheckRequest create() => CreateHealthCheckRequest._();
  @$core.override
  CreateHealthCheckRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateHealthCheckRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateHealthCheckRequest>(create);
  static CreateHealthCheckRequest? _defaultInstance;

  @$pb.TagNumber(83111204)
  HealthCheckConfig get healthcheckconfig => $_getN(0);
  @$pb.TagNumber(83111204)
  set healthcheckconfig(HealthCheckConfig value) => $_setField(83111204, value);
  @$pb.TagNumber(83111204)
  $core.bool hasHealthcheckconfig() => $_has(0);
  @$pb.TagNumber(83111204)
  void clearHealthcheckconfig() => $_clearField(83111204);
  @$pb.TagNumber(83111204)
  HealthCheckConfig ensureHealthcheckconfig() => $_ensure(0);

  @$pb.TagNumber(151211160)
  $core.String get callerreference => $_getSZ(1);
  @$pb.TagNumber(151211160)
  set callerreference($core.String value) => $_setString(1, value);
  @$pb.TagNumber(151211160)
  $core.bool hasCallerreference() => $_has(1);
  @$pb.TagNumber(151211160)
  void clearCallerreference() => $_clearField(151211160);
}

class CreateHealthCheckResponse extends $pb.GeneratedMessage {
  factory CreateHealthCheckResponse({
    HealthCheck? healthcheck,
    $core.String? location,
  }) {
    final result = create();
    if (healthcheck != null) result.healthcheck = healthcheck;
    if (location != null) result.location = location;
    return result;
  }

  CreateHealthCheckResponse._();

  factory CreateHealthCheckResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateHealthCheckResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateHealthCheckResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<HealthCheck>(377540610, _omitFieldNames ? '' : 'healthcheck',
        subBuilder: HealthCheck.create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateHealthCheckResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateHealthCheckResponse copyWith(
          void Function(CreateHealthCheckResponse) updates) =>
      super.copyWith((message) => updates(message as CreateHealthCheckResponse))
          as CreateHealthCheckResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateHealthCheckResponse create() => CreateHealthCheckResponse._();
  @$core.override
  CreateHealthCheckResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateHealthCheckResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateHealthCheckResponse>(create);
  static CreateHealthCheckResponse? _defaultInstance;

  @$pb.TagNumber(377540610)
  HealthCheck get healthcheck => $_getN(0);
  @$pb.TagNumber(377540610)
  set healthcheck(HealthCheck value) => $_setField(377540610, value);
  @$pb.TagNumber(377540610)
  $core.bool hasHealthcheck() => $_has(0);
  @$pb.TagNumber(377540610)
  void clearHealthcheck() => $_clearField(377540610);
  @$pb.TagNumber(377540610)
  HealthCheck ensureHealthcheck() => $_ensure(0);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(1);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(1, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(1);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);
}

class CreateHostedZoneRequest extends $pb.GeneratedMessage {
  factory CreateHostedZoneRequest({
    HostedZoneConfig? hostedzoneconfig,
    $core.String? callerreference,
    $core.String? name,
    $core.String? delegationsetid,
    VPC? vpc,
  }) {
    final result = create();
    if (hostedzoneconfig != null) result.hostedzoneconfig = hostedzoneconfig;
    if (callerreference != null) result.callerreference = callerreference;
    if (name != null) result.name = name;
    if (delegationsetid != null) result.delegationsetid = delegationsetid;
    if (vpc != null) result.vpc = vpc;
    return result;
  }

  CreateHostedZoneRequest._();

  factory CreateHostedZoneRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateHostedZoneRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateHostedZoneRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<HostedZoneConfig>(881519, _omitFieldNames ? '' : 'hostedzoneconfig',
        subBuilder: HostedZoneConfig.create)
    ..aOS(151211160, _omitFieldNames ? '' : 'callerreference')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(307328801, _omitFieldNames ? '' : 'delegationsetid')
    ..aOM<VPC>(506158953, _omitFieldNames ? '' : 'vpc', subBuilder: VPC.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateHostedZoneRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateHostedZoneRequest copyWith(
          void Function(CreateHostedZoneRequest) updates) =>
      super.copyWith((message) => updates(message as CreateHostedZoneRequest))
          as CreateHostedZoneRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateHostedZoneRequest create() => CreateHostedZoneRequest._();
  @$core.override
  CreateHostedZoneRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateHostedZoneRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateHostedZoneRequest>(create);
  static CreateHostedZoneRequest? _defaultInstance;

  @$pb.TagNumber(881519)
  HostedZoneConfig get hostedzoneconfig => $_getN(0);
  @$pb.TagNumber(881519)
  set hostedzoneconfig(HostedZoneConfig value) => $_setField(881519, value);
  @$pb.TagNumber(881519)
  $core.bool hasHostedzoneconfig() => $_has(0);
  @$pb.TagNumber(881519)
  void clearHostedzoneconfig() => $_clearField(881519);
  @$pb.TagNumber(881519)
  HostedZoneConfig ensureHostedzoneconfig() => $_ensure(0);

  @$pb.TagNumber(151211160)
  $core.String get callerreference => $_getSZ(1);
  @$pb.TagNumber(151211160)
  set callerreference($core.String value) => $_setString(1, value);
  @$pb.TagNumber(151211160)
  $core.bool hasCallerreference() => $_has(1);
  @$pb.TagNumber(151211160)
  void clearCallerreference() => $_clearField(151211160);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(307328801)
  $core.String get delegationsetid => $_getSZ(3);
  @$pb.TagNumber(307328801)
  set delegationsetid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(307328801)
  $core.bool hasDelegationsetid() => $_has(3);
  @$pb.TagNumber(307328801)
  void clearDelegationsetid() => $_clearField(307328801);

  @$pb.TagNumber(506158953)
  VPC get vpc => $_getN(4);
  @$pb.TagNumber(506158953)
  set vpc(VPC value) => $_setField(506158953, value);
  @$pb.TagNumber(506158953)
  $core.bool hasVpc() => $_has(4);
  @$pb.TagNumber(506158953)
  void clearVpc() => $_clearField(506158953);
  @$pb.TagNumber(506158953)
  VPC ensureVpc() => $_ensure(4);
}

class CreateHostedZoneResponse extends $pb.GeneratedMessage {
  factory CreateHostedZoneResponse({
    DelegationSet? delegationset,
    ChangeInfo? changeinfo,
    $core.String? location,
    HostedZone? hostedzone,
    VPC? vpc,
  }) {
    final result = create();
    if (delegationset != null) result.delegationset = delegationset;
    if (changeinfo != null) result.changeinfo = changeinfo;
    if (location != null) result.location = location;
    if (hostedzone != null) result.hostedzone = hostedzone;
    if (vpc != null) result.vpc = vpc;
    return result;
  }

  CreateHostedZoneResponse._();

  factory CreateHostedZoneResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateHostedZoneResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateHostedZoneResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<DelegationSet>(190771750, _omitFieldNames ? '' : 'delegationset',
        subBuilder: DelegationSet.create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..aOM<HostedZone>(465689249, _omitFieldNames ? '' : 'hostedzone',
        subBuilder: HostedZone.create)
    ..aOM<VPC>(506158953, _omitFieldNames ? '' : 'vpc', subBuilder: VPC.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateHostedZoneResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateHostedZoneResponse copyWith(
          void Function(CreateHostedZoneResponse) updates) =>
      super.copyWith((message) => updates(message as CreateHostedZoneResponse))
          as CreateHostedZoneResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateHostedZoneResponse create() => CreateHostedZoneResponse._();
  @$core.override
  CreateHostedZoneResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateHostedZoneResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateHostedZoneResponse>(create);
  static CreateHostedZoneResponse? _defaultInstance;

  @$pb.TagNumber(190771750)
  DelegationSet get delegationset => $_getN(0);
  @$pb.TagNumber(190771750)
  set delegationset(DelegationSet value) => $_setField(190771750, value);
  @$pb.TagNumber(190771750)
  $core.bool hasDelegationset() => $_has(0);
  @$pb.TagNumber(190771750)
  void clearDelegationset() => $_clearField(190771750);
  @$pb.TagNumber(190771750)
  DelegationSet ensureDelegationset() => $_ensure(0);

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(1);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(1);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(1);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(2);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(2, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(2);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);

  @$pb.TagNumber(465689249)
  HostedZone get hostedzone => $_getN(3);
  @$pb.TagNumber(465689249)
  set hostedzone(HostedZone value) => $_setField(465689249, value);
  @$pb.TagNumber(465689249)
  $core.bool hasHostedzone() => $_has(3);
  @$pb.TagNumber(465689249)
  void clearHostedzone() => $_clearField(465689249);
  @$pb.TagNumber(465689249)
  HostedZone ensureHostedzone() => $_ensure(3);

  @$pb.TagNumber(506158953)
  VPC get vpc => $_getN(4);
  @$pb.TagNumber(506158953)
  set vpc(VPC value) => $_setField(506158953, value);
  @$pb.TagNumber(506158953)
  $core.bool hasVpc() => $_has(4);
  @$pb.TagNumber(506158953)
  void clearVpc() => $_clearField(506158953);
  @$pb.TagNumber(506158953)
  VPC ensureVpc() => $_ensure(4);
}

class CreateKeySigningKeyRequest extends $pb.GeneratedMessage {
  factory CreateKeySigningKeyRequest({
    $core.String? status,
    $core.String? callerreference,
    $core.String? name,
    $core.String? hostedzoneid,
    $core.String? keymanagementservicearn,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (callerreference != null) result.callerreference = callerreference;
    if (name != null) result.name = name;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (keymanagementservicearn != null)
      result.keymanagementservicearn = keymanagementservicearn;
    return result;
  }

  CreateKeySigningKeyRequest._();

  factory CreateKeySigningKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateKeySigningKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateKeySigningKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(6222352, _omitFieldNames ? '' : 'status')
    ..aOS(151211160, _omitFieldNames ? '' : 'callerreference')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(360687398, _omitFieldNames ? '' : 'keymanagementservicearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateKeySigningKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateKeySigningKeyRequest copyWith(
          void Function(CreateKeySigningKeyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateKeySigningKeyRequest))
          as CreateKeySigningKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateKeySigningKeyRequest create() => CreateKeySigningKeyRequest._();
  @$core.override
  CreateKeySigningKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateKeySigningKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateKeySigningKeyRequest>(create);
  static CreateKeySigningKeyRequest? _defaultInstance;

  @$pb.TagNumber(6222352)
  $core.String get status => $_getSZ(0);
  @$pb.TagNumber(6222352)
  set status($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(151211160)
  $core.String get callerreference => $_getSZ(1);
  @$pb.TagNumber(151211160)
  set callerreference($core.String value) => $_setString(1, value);
  @$pb.TagNumber(151211160)
  $core.bool hasCallerreference() => $_has(1);
  @$pb.TagNumber(151211160)
  void clearCallerreference() => $_clearField(151211160);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(3);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(3, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(3);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(360687398)
  $core.String get keymanagementservicearn => $_getSZ(4);
  @$pb.TagNumber(360687398)
  set keymanagementservicearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(360687398)
  $core.bool hasKeymanagementservicearn() => $_has(4);
  @$pb.TagNumber(360687398)
  void clearKeymanagementservicearn() => $_clearField(360687398);
}

class CreateKeySigningKeyResponse extends $pb.GeneratedMessage {
  factory CreateKeySigningKeyResponse({
    KeySigningKey? keysigningkey,
    ChangeInfo? changeinfo,
    $core.String? location,
  }) {
    final result = create();
    if (keysigningkey != null) result.keysigningkey = keysigningkey;
    if (changeinfo != null) result.changeinfo = changeinfo;
    if (location != null) result.location = location;
    return result;
  }

  CreateKeySigningKeyResponse._();

  factory CreateKeySigningKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateKeySigningKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateKeySigningKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<KeySigningKey>(91098059, _omitFieldNames ? '' : 'keysigningkey',
        subBuilder: KeySigningKey.create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateKeySigningKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateKeySigningKeyResponse copyWith(
          void Function(CreateKeySigningKeyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateKeySigningKeyResponse))
          as CreateKeySigningKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateKeySigningKeyResponse create() =>
      CreateKeySigningKeyResponse._();
  @$core.override
  CreateKeySigningKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateKeySigningKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateKeySigningKeyResponse>(create);
  static CreateKeySigningKeyResponse? _defaultInstance;

  @$pb.TagNumber(91098059)
  KeySigningKey get keysigningkey => $_getN(0);
  @$pb.TagNumber(91098059)
  set keysigningkey(KeySigningKey value) => $_setField(91098059, value);
  @$pb.TagNumber(91098059)
  $core.bool hasKeysigningkey() => $_has(0);
  @$pb.TagNumber(91098059)
  void clearKeysigningkey() => $_clearField(91098059);
  @$pb.TagNumber(91098059)
  KeySigningKey ensureKeysigningkey() => $_ensure(0);

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(1);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(1);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(1);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(2);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(2, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(2);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);
}

class CreateQueryLoggingConfigRequest extends $pb.GeneratedMessage {
  factory CreateQueryLoggingConfigRequest({
    $core.String? hostedzoneid,
    $core.String? cloudwatchlogsloggrouparn,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (cloudwatchlogsloggrouparn != null)
      result.cloudwatchlogsloggrouparn = cloudwatchlogsloggrouparn;
    return result;
  }

  CreateQueryLoggingConfigRequest._();

  factory CreateQueryLoggingConfigRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateQueryLoggingConfigRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateQueryLoggingConfigRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(519613355, _omitFieldNames ? '' : 'cloudwatchlogsloggrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateQueryLoggingConfigRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateQueryLoggingConfigRequest copyWith(
          void Function(CreateQueryLoggingConfigRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateQueryLoggingConfigRequest))
          as CreateQueryLoggingConfigRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateQueryLoggingConfigRequest create() =>
      CreateQueryLoggingConfigRequest._();
  @$core.override
  CreateQueryLoggingConfigRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateQueryLoggingConfigRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateQueryLoggingConfigRequest>(
          create);
  static CreateQueryLoggingConfigRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(519613355)
  $core.String get cloudwatchlogsloggrouparn => $_getSZ(1);
  @$pb.TagNumber(519613355)
  set cloudwatchlogsloggrouparn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(519613355)
  $core.bool hasCloudwatchlogsloggrouparn() => $_has(1);
  @$pb.TagNumber(519613355)
  void clearCloudwatchlogsloggrouparn() => $_clearField(519613355);
}

class CreateQueryLoggingConfigResponse extends $pb.GeneratedMessage {
  factory CreateQueryLoggingConfigResponse({
    $core.String? location,
    QueryLoggingConfig? queryloggingconfig,
  }) {
    final result = create();
    if (location != null) result.location = location;
    if (queryloggingconfig != null)
      result.queryloggingconfig = queryloggingconfig;
    return result;
  }

  CreateQueryLoggingConfigResponse._();

  factory CreateQueryLoggingConfigResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateQueryLoggingConfigResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateQueryLoggingConfigResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..aOM<QueryLoggingConfig>(
        492751675, _omitFieldNames ? '' : 'queryloggingconfig',
        subBuilder: QueryLoggingConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateQueryLoggingConfigResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateQueryLoggingConfigResponse copyWith(
          void Function(CreateQueryLoggingConfigResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateQueryLoggingConfigResponse))
          as CreateQueryLoggingConfigResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateQueryLoggingConfigResponse create() =>
      CreateQueryLoggingConfigResponse._();
  @$core.override
  CreateQueryLoggingConfigResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateQueryLoggingConfigResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateQueryLoggingConfigResponse>(
          create);
  static CreateQueryLoggingConfigResponse? _defaultInstance;

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(0);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(0, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(0);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);

  @$pb.TagNumber(492751675)
  QueryLoggingConfig get queryloggingconfig => $_getN(1);
  @$pb.TagNumber(492751675)
  set queryloggingconfig(QueryLoggingConfig value) =>
      $_setField(492751675, value);
  @$pb.TagNumber(492751675)
  $core.bool hasQueryloggingconfig() => $_has(1);
  @$pb.TagNumber(492751675)
  void clearQueryloggingconfig() => $_clearField(492751675);
  @$pb.TagNumber(492751675)
  QueryLoggingConfig ensureQueryloggingconfig() => $_ensure(1);
}

class CreateReusableDelegationSetRequest extends $pb.GeneratedMessage {
  factory CreateReusableDelegationSetRequest({
    $core.String? callerreference,
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (callerreference != null) result.callerreference = callerreference;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  CreateReusableDelegationSetRequest._();

  factory CreateReusableDelegationSetRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateReusableDelegationSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateReusableDelegationSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(151211160, _omitFieldNames ? '' : 'callerreference')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateReusableDelegationSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateReusableDelegationSetRequest copyWith(
          void Function(CreateReusableDelegationSetRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreateReusableDelegationSetRequest))
          as CreateReusableDelegationSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateReusableDelegationSetRequest create() =>
      CreateReusableDelegationSetRequest._();
  @$core.override
  CreateReusableDelegationSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateReusableDelegationSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateReusableDelegationSetRequest>(
          create);
  static CreateReusableDelegationSetRequest? _defaultInstance;

  @$pb.TagNumber(151211160)
  $core.String get callerreference => $_getSZ(0);
  @$pb.TagNumber(151211160)
  set callerreference($core.String value) => $_setString(0, value);
  @$pb.TagNumber(151211160)
  $core.bool hasCallerreference() => $_has(0);
  @$pb.TagNumber(151211160)
  void clearCallerreference() => $_clearField(151211160);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class CreateReusableDelegationSetResponse extends $pb.GeneratedMessage {
  factory CreateReusableDelegationSetResponse({
    DelegationSet? delegationset,
    $core.String? location,
  }) {
    final result = create();
    if (delegationset != null) result.delegationset = delegationset;
    if (location != null) result.location = location;
    return result;
  }

  CreateReusableDelegationSetResponse._();

  factory CreateReusableDelegationSetResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateReusableDelegationSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateReusableDelegationSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<DelegationSet>(190771750, _omitFieldNames ? '' : 'delegationset',
        subBuilder: DelegationSet.create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateReusableDelegationSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateReusableDelegationSetResponse copyWith(
          void Function(CreateReusableDelegationSetResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreateReusableDelegationSetResponse))
          as CreateReusableDelegationSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateReusableDelegationSetResponse create() =>
      CreateReusableDelegationSetResponse._();
  @$core.override
  CreateReusableDelegationSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateReusableDelegationSetResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateReusableDelegationSetResponse>(create);
  static CreateReusableDelegationSetResponse? _defaultInstance;

  @$pb.TagNumber(190771750)
  DelegationSet get delegationset => $_getN(0);
  @$pb.TagNumber(190771750)
  set delegationset(DelegationSet value) => $_setField(190771750, value);
  @$pb.TagNumber(190771750)
  $core.bool hasDelegationset() => $_has(0);
  @$pb.TagNumber(190771750)
  void clearDelegationset() => $_clearField(190771750);
  @$pb.TagNumber(190771750)
  DelegationSet ensureDelegationset() => $_ensure(0);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(1);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(1, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(1);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);
}

class CreateTrafficPolicyInstanceRequest extends $pb.GeneratedMessage {
  factory CreateTrafficPolicyInstanceRequest({
    $core.String? trafficpolicyid,
    $core.String? name,
    $core.String? hostedzoneid,
    $core.int? trafficpolicyversion,
    $fixnum.Int64? ttl,
  }) {
    final result = create();
    if (trafficpolicyid != null) result.trafficpolicyid = trafficpolicyid;
    if (name != null) result.name = name;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (trafficpolicyversion != null)
      result.trafficpolicyversion = trafficpolicyversion;
    if (ttl != null) result.ttl = ttl;
    return result;
  }

  CreateTrafficPolicyInstanceRequest._();

  factory CreateTrafficPolicyInstanceRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTrafficPolicyInstanceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTrafficPolicyInstanceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(40235222, _omitFieldNames ? '' : 'trafficpolicyid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aI(479078485, _omitFieldNames ? '' : 'trafficpolicyversion')
    ..aInt64(526904700, _omitFieldNames ? '' : 'ttl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyInstanceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyInstanceRequest copyWith(
          void Function(CreateTrafficPolicyInstanceRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreateTrafficPolicyInstanceRequest))
          as CreateTrafficPolicyInstanceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyInstanceRequest create() =>
      CreateTrafficPolicyInstanceRequest._();
  @$core.override
  CreateTrafficPolicyInstanceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyInstanceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTrafficPolicyInstanceRequest>(
          create);
  static CreateTrafficPolicyInstanceRequest? _defaultInstance;

  @$pb.TagNumber(40235222)
  $core.String get trafficpolicyid => $_getSZ(0);
  @$pb.TagNumber(40235222)
  set trafficpolicyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(40235222)
  $core.bool hasTrafficpolicyid() => $_has(0);
  @$pb.TagNumber(40235222)
  void clearTrafficpolicyid() => $_clearField(40235222);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(2);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(2);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(479078485)
  $core.int get trafficpolicyversion => $_getIZ(3);
  @$pb.TagNumber(479078485)
  set trafficpolicyversion($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(479078485)
  $core.bool hasTrafficpolicyversion() => $_has(3);
  @$pb.TagNumber(479078485)
  void clearTrafficpolicyversion() => $_clearField(479078485);

  @$pb.TagNumber(526904700)
  $fixnum.Int64 get ttl => $_getI64(4);
  @$pb.TagNumber(526904700)
  set ttl($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(526904700)
  $core.bool hasTtl() => $_has(4);
  @$pb.TagNumber(526904700)
  void clearTtl() => $_clearField(526904700);
}

class CreateTrafficPolicyInstanceResponse extends $pb.GeneratedMessage {
  factory CreateTrafficPolicyInstanceResponse({
    TrafficPolicyInstance? trafficpolicyinstance,
    $core.String? location,
  }) {
    final result = create();
    if (trafficpolicyinstance != null)
      result.trafficpolicyinstance = trafficpolicyinstance;
    if (location != null) result.location = location;
    return result;
  }

  CreateTrafficPolicyInstanceResponse._();

  factory CreateTrafficPolicyInstanceResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTrafficPolicyInstanceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTrafficPolicyInstanceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<TrafficPolicyInstance>(
        205651476, _omitFieldNames ? '' : 'trafficpolicyinstance',
        subBuilder: TrafficPolicyInstance.create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyInstanceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyInstanceResponse copyWith(
          void Function(CreateTrafficPolicyInstanceResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreateTrafficPolicyInstanceResponse))
          as CreateTrafficPolicyInstanceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyInstanceResponse create() =>
      CreateTrafficPolicyInstanceResponse._();
  @$core.override
  CreateTrafficPolicyInstanceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyInstanceResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateTrafficPolicyInstanceResponse>(create);
  static CreateTrafficPolicyInstanceResponse? _defaultInstance;

  @$pb.TagNumber(205651476)
  TrafficPolicyInstance get trafficpolicyinstance => $_getN(0);
  @$pb.TagNumber(205651476)
  set trafficpolicyinstance(TrafficPolicyInstance value) =>
      $_setField(205651476, value);
  @$pb.TagNumber(205651476)
  $core.bool hasTrafficpolicyinstance() => $_has(0);
  @$pb.TagNumber(205651476)
  void clearTrafficpolicyinstance() => $_clearField(205651476);
  @$pb.TagNumber(205651476)
  TrafficPolicyInstance ensureTrafficpolicyinstance() => $_ensure(0);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(1);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(1, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(1);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);
}

class CreateTrafficPolicyRequest extends $pb.GeneratedMessage {
  factory CreateTrafficPolicyRequest({
    $core.String? name,
    $core.String? document,
    $core.String? comment,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (document != null) result.document = document;
    if (comment != null) result.comment = comment;
    return result;
  }

  CreateTrafficPolicyRequest._();

  factory CreateTrafficPolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTrafficPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTrafficPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(407108341, _omitFieldNames ? '' : 'document')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyRequest copyWith(
          void Function(CreateTrafficPolicyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateTrafficPolicyRequest))
          as CreateTrafficPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyRequest create() => CreateTrafficPolicyRequest._();
  @$core.override
  CreateTrafficPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTrafficPolicyRequest>(create);
  static CreateTrafficPolicyRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(407108341)
  $core.String get document => $_getSZ(1);
  @$pb.TagNumber(407108341)
  set document($core.String value) => $_setString(1, value);
  @$pb.TagNumber(407108341)
  $core.bool hasDocument() => $_has(1);
  @$pb.TagNumber(407108341)
  void clearDocument() => $_clearField(407108341);

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(2);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(2, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(2);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);
}

class CreateTrafficPolicyResponse extends $pb.GeneratedMessage {
  factory CreateTrafficPolicyResponse({
    TrafficPolicy? trafficpolicy,
    $core.String? location,
  }) {
    final result = create();
    if (trafficpolicy != null) result.trafficpolicy = trafficpolicy;
    if (location != null) result.location = location;
    return result;
  }

  CreateTrafficPolicyResponse._();

  factory CreateTrafficPolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTrafficPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTrafficPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<TrafficPolicy>(154595657, _omitFieldNames ? '' : 'trafficpolicy',
        subBuilder: TrafficPolicy.create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyResponse copyWith(
          void Function(CreateTrafficPolicyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateTrafficPolicyResponse))
          as CreateTrafficPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyResponse create() =>
      CreateTrafficPolicyResponse._();
  @$core.override
  CreateTrafficPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTrafficPolicyResponse>(create);
  static CreateTrafficPolicyResponse? _defaultInstance;

  @$pb.TagNumber(154595657)
  TrafficPolicy get trafficpolicy => $_getN(0);
  @$pb.TagNumber(154595657)
  set trafficpolicy(TrafficPolicy value) => $_setField(154595657, value);
  @$pb.TagNumber(154595657)
  $core.bool hasTrafficpolicy() => $_has(0);
  @$pb.TagNumber(154595657)
  void clearTrafficpolicy() => $_clearField(154595657);
  @$pb.TagNumber(154595657)
  TrafficPolicy ensureTrafficpolicy() => $_ensure(0);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(1);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(1, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(1);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);
}

class CreateTrafficPolicyVersionRequest extends $pb.GeneratedMessage {
  factory CreateTrafficPolicyVersionRequest({
    $core.String? id,
    $core.String? document,
    $core.String? comment,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (document != null) result.document = document;
    if (comment != null) result.comment = comment;
    return result;
  }

  CreateTrafficPolicyVersionRequest._();

  factory CreateTrafficPolicyVersionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTrafficPolicyVersionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTrafficPolicyVersionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(407108341, _omitFieldNames ? '' : 'document')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyVersionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyVersionRequest copyWith(
          void Function(CreateTrafficPolicyVersionRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreateTrafficPolicyVersionRequest))
          as CreateTrafficPolicyVersionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyVersionRequest create() =>
      CreateTrafficPolicyVersionRequest._();
  @$core.override
  CreateTrafficPolicyVersionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyVersionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTrafficPolicyVersionRequest>(
          create);
  static CreateTrafficPolicyVersionRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(407108341)
  $core.String get document => $_getSZ(1);
  @$pb.TagNumber(407108341)
  set document($core.String value) => $_setString(1, value);
  @$pb.TagNumber(407108341)
  $core.bool hasDocument() => $_has(1);
  @$pb.TagNumber(407108341)
  void clearDocument() => $_clearField(407108341);

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(2);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(2, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(2);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);
}

class CreateTrafficPolicyVersionResponse extends $pb.GeneratedMessage {
  factory CreateTrafficPolicyVersionResponse({
    TrafficPolicy? trafficpolicy,
    $core.String? location,
  }) {
    final result = create();
    if (trafficpolicy != null) result.trafficpolicy = trafficpolicy;
    if (location != null) result.location = location;
    return result;
  }

  CreateTrafficPolicyVersionResponse._();

  factory CreateTrafficPolicyVersionResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateTrafficPolicyVersionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateTrafficPolicyVersionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<TrafficPolicy>(154595657, _omitFieldNames ? '' : 'trafficpolicy',
        subBuilder: TrafficPolicy.create)
    ..aOS(465604039, _omitFieldNames ? '' : 'location')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyVersionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateTrafficPolicyVersionResponse copyWith(
          void Function(CreateTrafficPolicyVersionResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreateTrafficPolicyVersionResponse))
          as CreateTrafficPolicyVersionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyVersionResponse create() =>
      CreateTrafficPolicyVersionResponse._();
  @$core.override
  CreateTrafficPolicyVersionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateTrafficPolicyVersionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateTrafficPolicyVersionResponse>(
          create);
  static CreateTrafficPolicyVersionResponse? _defaultInstance;

  @$pb.TagNumber(154595657)
  TrafficPolicy get trafficpolicy => $_getN(0);
  @$pb.TagNumber(154595657)
  set trafficpolicy(TrafficPolicy value) => $_setField(154595657, value);
  @$pb.TagNumber(154595657)
  $core.bool hasTrafficpolicy() => $_has(0);
  @$pb.TagNumber(154595657)
  void clearTrafficpolicy() => $_clearField(154595657);
  @$pb.TagNumber(154595657)
  TrafficPolicy ensureTrafficpolicy() => $_ensure(0);

  @$pb.TagNumber(465604039)
  $core.String get location => $_getSZ(1);
  @$pb.TagNumber(465604039)
  set location($core.String value) => $_setString(1, value);
  @$pb.TagNumber(465604039)
  $core.bool hasLocation() => $_has(1);
  @$pb.TagNumber(465604039)
  void clearLocation() => $_clearField(465604039);
}

class CreateVPCAssociationAuthorizationRequest extends $pb.GeneratedMessage {
  factory CreateVPCAssociationAuthorizationRequest({
    $core.String? hostedzoneid,
    VPC? vpc,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (vpc != null) result.vpc = vpc;
    return result;
  }

  CreateVPCAssociationAuthorizationRequest._();

  factory CreateVPCAssociationAuthorizationRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateVPCAssociationAuthorizationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateVPCAssociationAuthorizationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOM<VPC>(506158953, _omitFieldNames ? '' : 'vpc', subBuilder: VPC.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateVPCAssociationAuthorizationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateVPCAssociationAuthorizationRequest copyWith(
          void Function(CreateVPCAssociationAuthorizationRequest) updates) =>
      super.copyWith((message) =>
              updates(message as CreateVPCAssociationAuthorizationRequest))
          as CreateVPCAssociationAuthorizationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateVPCAssociationAuthorizationRequest create() =>
      CreateVPCAssociationAuthorizationRequest._();
  @$core.override
  CreateVPCAssociationAuthorizationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateVPCAssociationAuthorizationRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateVPCAssociationAuthorizationRequest>(create);
  static CreateVPCAssociationAuthorizationRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(506158953)
  VPC get vpc => $_getN(1);
  @$pb.TagNumber(506158953)
  set vpc(VPC value) => $_setField(506158953, value);
  @$pb.TagNumber(506158953)
  $core.bool hasVpc() => $_has(1);
  @$pb.TagNumber(506158953)
  void clearVpc() => $_clearField(506158953);
  @$pb.TagNumber(506158953)
  VPC ensureVpc() => $_ensure(1);
}

class CreateVPCAssociationAuthorizationResponse extends $pb.GeneratedMessage {
  factory CreateVPCAssociationAuthorizationResponse({
    $core.String? hostedzoneid,
    VPC? vpc,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (vpc != null) result.vpc = vpc;
    return result;
  }

  CreateVPCAssociationAuthorizationResponse._();

  factory CreateVPCAssociationAuthorizationResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateVPCAssociationAuthorizationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateVPCAssociationAuthorizationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOM<VPC>(506158953, _omitFieldNames ? '' : 'vpc', subBuilder: VPC.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateVPCAssociationAuthorizationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateVPCAssociationAuthorizationResponse copyWith(
          void Function(CreateVPCAssociationAuthorizationResponse) updates) =>
      super.copyWith((message) =>
              updates(message as CreateVPCAssociationAuthorizationResponse))
          as CreateVPCAssociationAuthorizationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateVPCAssociationAuthorizationResponse create() =>
      CreateVPCAssociationAuthorizationResponse._();
  @$core.override
  CreateVPCAssociationAuthorizationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateVPCAssociationAuthorizationResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateVPCAssociationAuthorizationResponse>(create);
  static CreateVPCAssociationAuthorizationResponse? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(506158953)
  VPC get vpc => $_getN(1);
  @$pb.TagNumber(506158953)
  set vpc(VPC value) => $_setField(506158953, value);
  @$pb.TagNumber(506158953)
  $core.bool hasVpc() => $_has(1);
  @$pb.TagNumber(506158953)
  void clearVpc() => $_clearField(506158953);
  @$pb.TagNumber(506158953)
  VPC ensureVpc() => $_ensure(1);
}

class DNSSECNotFound extends $pb.GeneratedMessage {
  factory DNSSECNotFound({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DNSSECNotFound._();

  factory DNSSECNotFound.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DNSSECNotFound.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DNSSECNotFound',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DNSSECNotFound clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DNSSECNotFound copyWith(void Function(DNSSECNotFound) updates) =>
      super.copyWith((message) => updates(message as DNSSECNotFound))
          as DNSSECNotFound;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DNSSECNotFound create() => DNSSECNotFound._();
  @$core.override
  DNSSECNotFound createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DNSSECNotFound getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DNSSECNotFound>(create);
  static DNSSECNotFound? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DNSSECStatus extends $pb.GeneratedMessage {
  factory DNSSECStatus({
    $core.String? statusmessage,
    $core.String? servesignature,
  }) {
    final result = create();
    if (statusmessage != null) result.statusmessage = statusmessage;
    if (servesignature != null) result.servesignature = servesignature;
    return result;
  }

  DNSSECStatus._();

  factory DNSSECStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DNSSECStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DNSSECStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(72590095, _omitFieldNames ? '' : 'statusmessage')
    ..aOS(466239599, _omitFieldNames ? '' : 'servesignature')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DNSSECStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DNSSECStatus copyWith(void Function(DNSSECStatus) updates) =>
      super.copyWith((message) => updates(message as DNSSECStatus))
          as DNSSECStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DNSSECStatus create() => DNSSECStatus._();
  @$core.override
  DNSSECStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DNSSECStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DNSSECStatus>(create);
  static DNSSECStatus? _defaultInstance;

  @$pb.TagNumber(72590095)
  $core.String get statusmessage => $_getSZ(0);
  @$pb.TagNumber(72590095)
  set statusmessage($core.String value) => $_setString(0, value);
  @$pb.TagNumber(72590095)
  $core.bool hasStatusmessage() => $_has(0);
  @$pb.TagNumber(72590095)
  void clearStatusmessage() => $_clearField(72590095);

  @$pb.TagNumber(466239599)
  $core.String get servesignature => $_getSZ(1);
  @$pb.TagNumber(466239599)
  set servesignature($core.String value) => $_setString(1, value);
  @$pb.TagNumber(466239599)
  $core.bool hasServesignature() => $_has(1);
  @$pb.TagNumber(466239599)
  void clearServesignature() => $_clearField(466239599);
}

class DeactivateKeySigningKeyRequest extends $pb.GeneratedMessage {
  factory DeactivateKeySigningKeyRequest({
    $core.String? name,
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  DeactivateKeySigningKeyRequest._();

  factory DeactivateKeySigningKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeactivateKeySigningKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeactivateKeySigningKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeactivateKeySigningKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeactivateKeySigningKeyRequest copyWith(
          void Function(DeactivateKeySigningKeyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeactivateKeySigningKeyRequest))
          as DeactivateKeySigningKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeactivateKeySigningKeyRequest create() =>
      DeactivateKeySigningKeyRequest._();
  @$core.override
  DeactivateKeySigningKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeactivateKeySigningKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeactivateKeySigningKeyRequest>(create);
  static DeactivateKeySigningKeyRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class DeactivateKeySigningKeyResponse extends $pb.GeneratedMessage {
  factory DeactivateKeySigningKeyResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  DeactivateKeySigningKeyResponse._();

  factory DeactivateKeySigningKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeactivateKeySigningKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeactivateKeySigningKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeactivateKeySigningKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeactivateKeySigningKeyResponse copyWith(
          void Function(DeactivateKeySigningKeyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeactivateKeySigningKeyResponse))
          as DeactivateKeySigningKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeactivateKeySigningKeyResponse create() =>
      DeactivateKeySigningKeyResponse._();
  @$core.override
  DeactivateKeySigningKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeactivateKeySigningKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeactivateKeySigningKeyResponse>(
          create);
  static DeactivateKeySigningKeyResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class DelegationSet extends $pb.GeneratedMessage {
  factory DelegationSet({
    $core.String? callerreference,
    $core.Iterable<$core.String>? nameservers,
    $core.String? id,
  }) {
    final result = create();
    if (callerreference != null) result.callerreference = callerreference;
    if (nameservers != null) result.nameservers.addAll(nameservers);
    if (id != null) result.id = id;
    return result;
  }

  DelegationSet._();

  factory DelegationSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DelegationSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DelegationSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(151211160, _omitFieldNames ? '' : 'callerreference')
    ..pPS(340971511, _omitFieldNames ? '' : 'nameservers')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSet copyWith(void Function(DelegationSet) updates) =>
      super.copyWith((message) => updates(message as DelegationSet))
          as DelegationSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DelegationSet create() => DelegationSet._();
  @$core.override
  DelegationSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DelegationSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DelegationSet>(create);
  static DelegationSet? _defaultInstance;

  @$pb.TagNumber(151211160)
  $core.String get callerreference => $_getSZ(0);
  @$pb.TagNumber(151211160)
  set callerreference($core.String value) => $_setString(0, value);
  @$pb.TagNumber(151211160)
  $core.bool hasCallerreference() => $_has(0);
  @$pb.TagNumber(151211160)
  void clearCallerreference() => $_clearField(151211160);

  @$pb.TagNumber(340971511)
  $pb.PbList<$core.String> get nameservers => $_getList(1);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DelegationSetAlreadyCreated extends $pb.GeneratedMessage {
  factory DelegationSetAlreadyCreated({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DelegationSetAlreadyCreated._();

  factory DelegationSetAlreadyCreated.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DelegationSetAlreadyCreated.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DelegationSetAlreadyCreated',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetAlreadyCreated clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetAlreadyCreated copyWith(
          void Function(DelegationSetAlreadyCreated) updates) =>
      super.copyWith(
              (message) => updates(message as DelegationSetAlreadyCreated))
          as DelegationSetAlreadyCreated;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DelegationSetAlreadyCreated create() =>
      DelegationSetAlreadyCreated._();
  @$core.override
  DelegationSetAlreadyCreated createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DelegationSetAlreadyCreated getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DelegationSetAlreadyCreated>(create);
  static DelegationSetAlreadyCreated? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DelegationSetAlreadyReusable extends $pb.GeneratedMessage {
  factory DelegationSetAlreadyReusable({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DelegationSetAlreadyReusable._();

  factory DelegationSetAlreadyReusable.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DelegationSetAlreadyReusable.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DelegationSetAlreadyReusable',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetAlreadyReusable clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetAlreadyReusable copyWith(
          void Function(DelegationSetAlreadyReusable) updates) =>
      super.copyWith(
              (message) => updates(message as DelegationSetAlreadyReusable))
          as DelegationSetAlreadyReusable;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DelegationSetAlreadyReusable create() =>
      DelegationSetAlreadyReusable._();
  @$core.override
  DelegationSetAlreadyReusable createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DelegationSetAlreadyReusable getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DelegationSetAlreadyReusable>(create);
  static DelegationSetAlreadyReusable? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DelegationSetInUse extends $pb.GeneratedMessage {
  factory DelegationSetInUse({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DelegationSetInUse._();

  factory DelegationSetInUse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DelegationSetInUse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DelegationSetInUse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetInUse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetInUse copyWith(void Function(DelegationSetInUse) updates) =>
      super.copyWith((message) => updates(message as DelegationSetInUse))
          as DelegationSetInUse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DelegationSetInUse create() => DelegationSetInUse._();
  @$core.override
  DelegationSetInUse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DelegationSetInUse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DelegationSetInUse>(create);
  static DelegationSetInUse? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DelegationSetNotAvailable extends $pb.GeneratedMessage {
  factory DelegationSetNotAvailable({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DelegationSetNotAvailable._();

  factory DelegationSetNotAvailable.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DelegationSetNotAvailable.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DelegationSetNotAvailable',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetNotAvailable clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetNotAvailable copyWith(
          void Function(DelegationSetNotAvailable) updates) =>
      super.copyWith((message) => updates(message as DelegationSetNotAvailable))
          as DelegationSetNotAvailable;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DelegationSetNotAvailable create() => DelegationSetNotAvailable._();
  @$core.override
  DelegationSetNotAvailable createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DelegationSetNotAvailable getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DelegationSetNotAvailable>(create);
  static DelegationSetNotAvailable? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DelegationSetNotReusable extends $pb.GeneratedMessage {
  factory DelegationSetNotReusable({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DelegationSetNotReusable._();

  factory DelegationSetNotReusable.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DelegationSetNotReusable.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DelegationSetNotReusable',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetNotReusable clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DelegationSetNotReusable copyWith(
          void Function(DelegationSetNotReusable) updates) =>
      super.copyWith((message) => updates(message as DelegationSetNotReusable))
          as DelegationSetNotReusable;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DelegationSetNotReusable create() => DelegationSetNotReusable._();
  @$core.override
  DelegationSetNotReusable createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DelegationSetNotReusable getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DelegationSetNotReusable>(create);
  static DelegationSetNotReusable? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DeleteCidrCollectionRequest extends $pb.GeneratedMessage {
  factory DeleteCidrCollectionRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  DeleteCidrCollectionRequest._();

  factory DeleteCidrCollectionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteCidrCollectionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteCidrCollectionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCidrCollectionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCidrCollectionRequest copyWith(
          void Function(DeleteCidrCollectionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteCidrCollectionRequest))
          as DeleteCidrCollectionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteCidrCollectionRequest create() =>
      DeleteCidrCollectionRequest._();
  @$core.override
  DeleteCidrCollectionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteCidrCollectionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteCidrCollectionRequest>(create);
  static DeleteCidrCollectionRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteCidrCollectionResponse extends $pb.GeneratedMessage {
  factory DeleteCidrCollectionResponse() => create();

  DeleteCidrCollectionResponse._();

  factory DeleteCidrCollectionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteCidrCollectionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteCidrCollectionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCidrCollectionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteCidrCollectionResponse copyWith(
          void Function(DeleteCidrCollectionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteCidrCollectionResponse))
          as DeleteCidrCollectionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteCidrCollectionResponse create() =>
      DeleteCidrCollectionResponse._();
  @$core.override
  DeleteCidrCollectionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteCidrCollectionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteCidrCollectionResponse>(create);
  static DeleteCidrCollectionResponse? _defaultInstance;
}

class DeleteHealthCheckRequest extends $pb.GeneratedMessage {
  factory DeleteHealthCheckRequest({
    $core.String? healthcheckid,
  }) {
    final result = create();
    if (healthcheckid != null) result.healthcheckid = healthcheckid;
    return result;
  }

  DeleteHealthCheckRequest._();

  factory DeleteHealthCheckRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteHealthCheckRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteHealthCheckRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(312971637, _omitFieldNames ? '' : 'healthcheckid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHealthCheckRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHealthCheckRequest copyWith(
          void Function(DeleteHealthCheckRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteHealthCheckRequest))
          as DeleteHealthCheckRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteHealthCheckRequest create() => DeleteHealthCheckRequest._();
  @$core.override
  DeleteHealthCheckRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteHealthCheckRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteHealthCheckRequest>(create);
  static DeleteHealthCheckRequest? _defaultInstance;

  @$pb.TagNumber(312971637)
  $core.String get healthcheckid => $_getSZ(0);
  @$pb.TagNumber(312971637)
  set healthcheckid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(312971637)
  $core.bool hasHealthcheckid() => $_has(0);
  @$pb.TagNumber(312971637)
  void clearHealthcheckid() => $_clearField(312971637);
}

class DeleteHealthCheckResponse extends $pb.GeneratedMessage {
  factory DeleteHealthCheckResponse() => create();

  DeleteHealthCheckResponse._();

  factory DeleteHealthCheckResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteHealthCheckResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteHealthCheckResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHealthCheckResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHealthCheckResponse copyWith(
          void Function(DeleteHealthCheckResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteHealthCheckResponse))
          as DeleteHealthCheckResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteHealthCheckResponse create() => DeleteHealthCheckResponse._();
  @$core.override
  DeleteHealthCheckResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteHealthCheckResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteHealthCheckResponse>(create);
  static DeleteHealthCheckResponse? _defaultInstance;
}

class DeleteHostedZoneRequest extends $pb.GeneratedMessage {
  factory DeleteHostedZoneRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  DeleteHostedZoneRequest._();

  factory DeleteHostedZoneRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteHostedZoneRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteHostedZoneRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHostedZoneRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHostedZoneRequest copyWith(
          void Function(DeleteHostedZoneRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteHostedZoneRequest))
          as DeleteHostedZoneRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteHostedZoneRequest create() => DeleteHostedZoneRequest._();
  @$core.override
  DeleteHostedZoneRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteHostedZoneRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteHostedZoneRequest>(create);
  static DeleteHostedZoneRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteHostedZoneResponse extends $pb.GeneratedMessage {
  factory DeleteHostedZoneResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  DeleteHostedZoneResponse._();

  factory DeleteHostedZoneResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteHostedZoneResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteHostedZoneResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHostedZoneResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteHostedZoneResponse copyWith(
          void Function(DeleteHostedZoneResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteHostedZoneResponse))
          as DeleteHostedZoneResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteHostedZoneResponse create() => DeleteHostedZoneResponse._();
  @$core.override
  DeleteHostedZoneResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteHostedZoneResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteHostedZoneResponse>(create);
  static DeleteHostedZoneResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class DeleteKeySigningKeyRequest extends $pb.GeneratedMessage {
  factory DeleteKeySigningKeyRequest({
    $core.String? name,
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  DeleteKeySigningKeyRequest._();

  factory DeleteKeySigningKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteKeySigningKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteKeySigningKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteKeySigningKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteKeySigningKeyRequest copyWith(
          void Function(DeleteKeySigningKeyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteKeySigningKeyRequest))
          as DeleteKeySigningKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteKeySigningKeyRequest create() => DeleteKeySigningKeyRequest._();
  @$core.override
  DeleteKeySigningKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteKeySigningKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteKeySigningKeyRequest>(create);
  static DeleteKeySigningKeyRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class DeleteKeySigningKeyResponse extends $pb.GeneratedMessage {
  factory DeleteKeySigningKeyResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  DeleteKeySigningKeyResponse._();

  factory DeleteKeySigningKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteKeySigningKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteKeySigningKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteKeySigningKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteKeySigningKeyResponse copyWith(
          void Function(DeleteKeySigningKeyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteKeySigningKeyResponse))
          as DeleteKeySigningKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteKeySigningKeyResponse create() =>
      DeleteKeySigningKeyResponse._();
  @$core.override
  DeleteKeySigningKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteKeySigningKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteKeySigningKeyResponse>(create);
  static DeleteKeySigningKeyResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class DeleteQueryLoggingConfigRequest extends $pb.GeneratedMessage {
  factory DeleteQueryLoggingConfigRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  DeleteQueryLoggingConfigRequest._();

  factory DeleteQueryLoggingConfigRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteQueryLoggingConfigRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteQueryLoggingConfigRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteQueryLoggingConfigRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteQueryLoggingConfigRequest copyWith(
          void Function(DeleteQueryLoggingConfigRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteQueryLoggingConfigRequest))
          as DeleteQueryLoggingConfigRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteQueryLoggingConfigRequest create() =>
      DeleteQueryLoggingConfigRequest._();
  @$core.override
  DeleteQueryLoggingConfigRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteQueryLoggingConfigRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteQueryLoggingConfigRequest>(
          create);
  static DeleteQueryLoggingConfigRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteQueryLoggingConfigResponse extends $pb.GeneratedMessage {
  factory DeleteQueryLoggingConfigResponse() => create();

  DeleteQueryLoggingConfigResponse._();

  factory DeleteQueryLoggingConfigResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteQueryLoggingConfigResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteQueryLoggingConfigResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteQueryLoggingConfigResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteQueryLoggingConfigResponse copyWith(
          void Function(DeleteQueryLoggingConfigResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteQueryLoggingConfigResponse))
          as DeleteQueryLoggingConfigResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteQueryLoggingConfigResponse create() =>
      DeleteQueryLoggingConfigResponse._();
  @$core.override
  DeleteQueryLoggingConfigResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteQueryLoggingConfigResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteQueryLoggingConfigResponse>(
          create);
  static DeleteQueryLoggingConfigResponse? _defaultInstance;
}

class DeleteReusableDelegationSetRequest extends $pb.GeneratedMessage {
  factory DeleteReusableDelegationSetRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  DeleteReusableDelegationSetRequest._();

  factory DeleteReusableDelegationSetRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteReusableDelegationSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteReusableDelegationSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteReusableDelegationSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteReusableDelegationSetRequest copyWith(
          void Function(DeleteReusableDelegationSetRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteReusableDelegationSetRequest))
          as DeleteReusableDelegationSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteReusableDelegationSetRequest create() =>
      DeleteReusableDelegationSetRequest._();
  @$core.override
  DeleteReusableDelegationSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteReusableDelegationSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteReusableDelegationSetRequest>(
          create);
  static DeleteReusableDelegationSetRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteReusableDelegationSetResponse extends $pb.GeneratedMessage {
  factory DeleteReusableDelegationSetResponse() => create();

  DeleteReusableDelegationSetResponse._();

  factory DeleteReusableDelegationSetResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteReusableDelegationSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteReusableDelegationSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteReusableDelegationSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteReusableDelegationSetResponse copyWith(
          void Function(DeleteReusableDelegationSetResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteReusableDelegationSetResponse))
          as DeleteReusableDelegationSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteReusableDelegationSetResponse create() =>
      DeleteReusableDelegationSetResponse._();
  @$core.override
  DeleteReusableDelegationSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteReusableDelegationSetResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeleteReusableDelegationSetResponse>(create);
  static DeleteReusableDelegationSetResponse? _defaultInstance;
}

class DeleteTrafficPolicyInstanceRequest extends $pb.GeneratedMessage {
  factory DeleteTrafficPolicyInstanceRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  DeleteTrafficPolicyInstanceRequest._();

  factory DeleteTrafficPolicyInstanceRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTrafficPolicyInstanceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTrafficPolicyInstanceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrafficPolicyInstanceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrafficPolicyInstanceRequest copyWith(
          void Function(DeleteTrafficPolicyInstanceRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteTrafficPolicyInstanceRequest))
          as DeleteTrafficPolicyInstanceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTrafficPolicyInstanceRequest create() =>
      DeleteTrafficPolicyInstanceRequest._();
  @$core.override
  DeleteTrafficPolicyInstanceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTrafficPolicyInstanceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTrafficPolicyInstanceRequest>(
          create);
  static DeleteTrafficPolicyInstanceRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class DeleteTrafficPolicyInstanceResponse extends $pb.GeneratedMessage {
  factory DeleteTrafficPolicyInstanceResponse() => create();

  DeleteTrafficPolicyInstanceResponse._();

  factory DeleteTrafficPolicyInstanceResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTrafficPolicyInstanceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTrafficPolicyInstanceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrafficPolicyInstanceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrafficPolicyInstanceResponse copyWith(
          void Function(DeleteTrafficPolicyInstanceResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteTrafficPolicyInstanceResponse))
          as DeleteTrafficPolicyInstanceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTrafficPolicyInstanceResponse create() =>
      DeleteTrafficPolicyInstanceResponse._();
  @$core.override
  DeleteTrafficPolicyInstanceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTrafficPolicyInstanceResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeleteTrafficPolicyInstanceResponse>(create);
  static DeleteTrafficPolicyInstanceResponse? _defaultInstance;
}

class DeleteTrafficPolicyRequest extends $pb.GeneratedMessage {
  factory DeleteTrafficPolicyRequest({
    $core.String? id,
    $core.int? version,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (version != null) result.version = version;
    return result;
  }

  DeleteTrafficPolicyRequest._();

  factory DeleteTrafficPolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTrafficPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTrafficPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aI(500028728, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrafficPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrafficPolicyRequest copyWith(
          void Function(DeleteTrafficPolicyRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteTrafficPolicyRequest))
          as DeleteTrafficPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTrafficPolicyRequest create() => DeleteTrafficPolicyRequest._();
  @$core.override
  DeleteTrafficPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTrafficPolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTrafficPolicyRequest>(create);
  static DeleteTrafficPolicyRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(500028728)
  $core.int get version => $_getIZ(1);
  @$pb.TagNumber(500028728)
  set version($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(1);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);
}

class DeleteTrafficPolicyResponse extends $pb.GeneratedMessage {
  factory DeleteTrafficPolicyResponse() => create();

  DeleteTrafficPolicyResponse._();

  factory DeleteTrafficPolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteTrafficPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteTrafficPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrafficPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteTrafficPolicyResponse copyWith(
          void Function(DeleteTrafficPolicyResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteTrafficPolicyResponse))
          as DeleteTrafficPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteTrafficPolicyResponse create() =>
      DeleteTrafficPolicyResponse._();
  @$core.override
  DeleteTrafficPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteTrafficPolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteTrafficPolicyResponse>(create);
  static DeleteTrafficPolicyResponse? _defaultInstance;
}

class DeleteVPCAssociationAuthorizationRequest extends $pb.GeneratedMessage {
  factory DeleteVPCAssociationAuthorizationRequest({
    $core.String? hostedzoneid,
    VPC? vpc,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (vpc != null) result.vpc = vpc;
    return result;
  }

  DeleteVPCAssociationAuthorizationRequest._();

  factory DeleteVPCAssociationAuthorizationRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteVPCAssociationAuthorizationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteVPCAssociationAuthorizationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOM<VPC>(506158953, _omitFieldNames ? '' : 'vpc', subBuilder: VPC.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteVPCAssociationAuthorizationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteVPCAssociationAuthorizationRequest copyWith(
          void Function(DeleteVPCAssociationAuthorizationRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteVPCAssociationAuthorizationRequest))
          as DeleteVPCAssociationAuthorizationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteVPCAssociationAuthorizationRequest create() =>
      DeleteVPCAssociationAuthorizationRequest._();
  @$core.override
  DeleteVPCAssociationAuthorizationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteVPCAssociationAuthorizationRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeleteVPCAssociationAuthorizationRequest>(create);
  static DeleteVPCAssociationAuthorizationRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(506158953)
  VPC get vpc => $_getN(1);
  @$pb.TagNumber(506158953)
  set vpc(VPC value) => $_setField(506158953, value);
  @$pb.TagNumber(506158953)
  $core.bool hasVpc() => $_has(1);
  @$pb.TagNumber(506158953)
  void clearVpc() => $_clearField(506158953);
  @$pb.TagNumber(506158953)
  VPC ensureVpc() => $_ensure(1);
}

class DeleteVPCAssociationAuthorizationResponse extends $pb.GeneratedMessage {
  factory DeleteVPCAssociationAuthorizationResponse() => create();

  DeleteVPCAssociationAuthorizationResponse._();

  factory DeleteVPCAssociationAuthorizationResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteVPCAssociationAuthorizationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteVPCAssociationAuthorizationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteVPCAssociationAuthorizationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteVPCAssociationAuthorizationResponse copyWith(
          void Function(DeleteVPCAssociationAuthorizationResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DeleteVPCAssociationAuthorizationResponse))
          as DeleteVPCAssociationAuthorizationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteVPCAssociationAuthorizationResponse create() =>
      DeleteVPCAssociationAuthorizationResponse._();
  @$core.override
  DeleteVPCAssociationAuthorizationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteVPCAssociationAuthorizationResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DeleteVPCAssociationAuthorizationResponse>(create);
  static DeleteVPCAssociationAuthorizationResponse? _defaultInstance;
}

class Dimension extends $pb.GeneratedMessage {
  factory Dimension({
    $core.String? name,
    $core.String? value,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (value != null) result.value = value;
    return result;
  }

  Dimension._();

  factory Dimension.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Dimension.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Dimension',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Dimension clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Dimension copyWith(void Function(Dimension) updates) =>
      super.copyWith((message) => updates(message as Dimension)) as Dimension;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Dimension create() => Dimension._();
  @$core.override
  Dimension createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Dimension getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Dimension>(create);
  static Dimension? _defaultInstance;

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

class DisableHostedZoneDNSSECRequest extends $pb.GeneratedMessage {
  factory DisableHostedZoneDNSSECRequest({
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  DisableHostedZoneDNSSECRequest._();

  factory DisableHostedZoneDNSSECRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableHostedZoneDNSSECRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableHostedZoneDNSSECRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableHostedZoneDNSSECRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableHostedZoneDNSSECRequest copyWith(
          void Function(DisableHostedZoneDNSSECRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DisableHostedZoneDNSSECRequest))
          as DisableHostedZoneDNSSECRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableHostedZoneDNSSECRequest create() =>
      DisableHostedZoneDNSSECRequest._();
  @$core.override
  DisableHostedZoneDNSSECRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableHostedZoneDNSSECRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableHostedZoneDNSSECRequest>(create);
  static DisableHostedZoneDNSSECRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class DisableHostedZoneDNSSECResponse extends $pb.GeneratedMessage {
  factory DisableHostedZoneDNSSECResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  DisableHostedZoneDNSSECResponse._();

  factory DisableHostedZoneDNSSECResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableHostedZoneDNSSECResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableHostedZoneDNSSECResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableHostedZoneDNSSECResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableHostedZoneDNSSECResponse copyWith(
          void Function(DisableHostedZoneDNSSECResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DisableHostedZoneDNSSECResponse))
          as DisableHostedZoneDNSSECResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableHostedZoneDNSSECResponse create() =>
      DisableHostedZoneDNSSECResponse._();
  @$core.override
  DisableHostedZoneDNSSECResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableHostedZoneDNSSECResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableHostedZoneDNSSECResponse>(
          create);
  static DisableHostedZoneDNSSECResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class DisassociateVPCFromHostedZoneRequest extends $pb.GeneratedMessage {
  factory DisassociateVPCFromHostedZoneRequest({
    $core.String? hostedzoneid,
    $core.String? comment,
    VPC? vpc,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (comment != null) result.comment = comment;
    if (vpc != null) result.vpc = vpc;
    return result;
  }

  DisassociateVPCFromHostedZoneRequest._();

  factory DisassociateVPCFromHostedZoneRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisassociateVPCFromHostedZoneRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisassociateVPCFromHostedZoneRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..aOM<VPC>(506158953, _omitFieldNames ? '' : 'vpc', subBuilder: VPC.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisassociateVPCFromHostedZoneRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisassociateVPCFromHostedZoneRequest copyWith(
          void Function(DisassociateVPCFromHostedZoneRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DisassociateVPCFromHostedZoneRequest))
          as DisassociateVPCFromHostedZoneRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisassociateVPCFromHostedZoneRequest create() =>
      DisassociateVPCFromHostedZoneRequest._();
  @$core.override
  DisassociateVPCFromHostedZoneRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisassociateVPCFromHostedZoneRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DisassociateVPCFromHostedZoneRequest>(create);
  static DisassociateVPCFromHostedZoneRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(1);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(1, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(1);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);

  @$pb.TagNumber(506158953)
  VPC get vpc => $_getN(2);
  @$pb.TagNumber(506158953)
  set vpc(VPC value) => $_setField(506158953, value);
  @$pb.TagNumber(506158953)
  $core.bool hasVpc() => $_has(2);
  @$pb.TagNumber(506158953)
  void clearVpc() => $_clearField(506158953);
  @$pb.TagNumber(506158953)
  VPC ensureVpc() => $_ensure(2);
}

class DisassociateVPCFromHostedZoneResponse extends $pb.GeneratedMessage {
  factory DisassociateVPCFromHostedZoneResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  DisassociateVPCFromHostedZoneResponse._();

  factory DisassociateVPCFromHostedZoneResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisassociateVPCFromHostedZoneResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisassociateVPCFromHostedZoneResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisassociateVPCFromHostedZoneResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisassociateVPCFromHostedZoneResponse copyWith(
          void Function(DisassociateVPCFromHostedZoneResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DisassociateVPCFromHostedZoneResponse))
          as DisassociateVPCFromHostedZoneResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisassociateVPCFromHostedZoneResponse create() =>
      DisassociateVPCFromHostedZoneResponse._();
  @$core.override
  DisassociateVPCFromHostedZoneResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisassociateVPCFromHostedZoneResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          DisassociateVPCFromHostedZoneResponse>(create);
  static DisassociateVPCFromHostedZoneResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class EnableHostedZoneDNSSECRequest extends $pb.GeneratedMessage {
  factory EnableHostedZoneDNSSECRequest({
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  EnableHostedZoneDNSSECRequest._();

  factory EnableHostedZoneDNSSECRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableHostedZoneDNSSECRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableHostedZoneDNSSECRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableHostedZoneDNSSECRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableHostedZoneDNSSECRequest copyWith(
          void Function(EnableHostedZoneDNSSECRequest) updates) =>
      super.copyWith(
              (message) => updates(message as EnableHostedZoneDNSSECRequest))
          as EnableHostedZoneDNSSECRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableHostedZoneDNSSECRequest create() =>
      EnableHostedZoneDNSSECRequest._();
  @$core.override
  EnableHostedZoneDNSSECRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableHostedZoneDNSSECRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableHostedZoneDNSSECRequest>(create);
  static EnableHostedZoneDNSSECRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class EnableHostedZoneDNSSECResponse extends $pb.GeneratedMessage {
  factory EnableHostedZoneDNSSECResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  EnableHostedZoneDNSSECResponse._();

  factory EnableHostedZoneDNSSECResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableHostedZoneDNSSECResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableHostedZoneDNSSECResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableHostedZoneDNSSECResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableHostedZoneDNSSECResponse copyWith(
          void Function(EnableHostedZoneDNSSECResponse) updates) =>
      super.copyWith(
              (message) => updates(message as EnableHostedZoneDNSSECResponse))
          as EnableHostedZoneDNSSECResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableHostedZoneDNSSECResponse create() =>
      EnableHostedZoneDNSSECResponse._();
  @$core.override
  EnableHostedZoneDNSSECResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableHostedZoneDNSSECResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableHostedZoneDNSSECResponse>(create);
  static EnableHostedZoneDNSSECResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class GeoLocation extends $pb.GeneratedMessage {
  factory GeoLocation({
    $core.String? continentcode,
    $core.String? subdivisioncode,
    $core.String? countrycode,
  }) {
    final result = create();
    if (continentcode != null) result.continentcode = continentcode;
    if (subdivisioncode != null) result.subdivisioncode = subdivisioncode;
    if (countrycode != null) result.countrycode = countrycode;
    return result;
  }

  GeoLocation._();

  factory GeoLocation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GeoLocation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoLocation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(180194347, _omitFieldNames ? '' : 'continentcode')
    ..aOS(444328268, _omitFieldNames ? '' : 'subdivisioncode')
    ..aOS(485287065, _omitFieldNames ? '' : 'countrycode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoLocation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoLocation copyWith(void Function(GeoLocation) updates) =>
      super.copyWith((message) => updates(message as GeoLocation))
          as GeoLocation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoLocation create() => GeoLocation._();
  @$core.override
  GeoLocation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GeoLocation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GeoLocation>(create);
  static GeoLocation? _defaultInstance;

  @$pb.TagNumber(180194347)
  $core.String get continentcode => $_getSZ(0);
  @$pb.TagNumber(180194347)
  set continentcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(180194347)
  $core.bool hasContinentcode() => $_has(0);
  @$pb.TagNumber(180194347)
  void clearContinentcode() => $_clearField(180194347);

  @$pb.TagNumber(444328268)
  $core.String get subdivisioncode => $_getSZ(1);
  @$pb.TagNumber(444328268)
  set subdivisioncode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(444328268)
  $core.bool hasSubdivisioncode() => $_has(1);
  @$pb.TagNumber(444328268)
  void clearSubdivisioncode() => $_clearField(444328268);

  @$pb.TagNumber(485287065)
  $core.String get countrycode => $_getSZ(2);
  @$pb.TagNumber(485287065)
  set countrycode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(485287065)
  $core.bool hasCountrycode() => $_has(2);
  @$pb.TagNumber(485287065)
  void clearCountrycode() => $_clearField(485287065);
}

class GeoLocationDetails extends $pb.GeneratedMessage {
  factory GeoLocationDetails({
    $core.String? countryname,
    $core.String? continentcode,
    $core.String? continentname,
    $core.String? subdivisionname,
    $core.String? subdivisioncode,
    $core.String? countrycode,
  }) {
    final result = create();
    if (countryname != null) result.countryname = countryname;
    if (continentcode != null) result.continentcode = continentcode;
    if (continentname != null) result.continentname = continentname;
    if (subdivisionname != null) result.subdivisionname = subdivisionname;
    if (subdivisioncode != null) result.subdivisioncode = subdivisioncode;
    if (countrycode != null) result.countrycode = countrycode;
    return result;
  }

  GeoLocationDetails._();

  factory GeoLocationDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GeoLocationDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoLocationDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(16070667, _omitFieldNames ? '' : 'countryname')
    ..aOS(180194347, _omitFieldNames ? '' : 'continentcode')
    ..aOS(203176953, _omitFieldNames ? '' : 'continentname')
    ..aOS(412782294, _omitFieldNames ? '' : 'subdivisionname')
    ..aOS(444328268, _omitFieldNames ? '' : 'subdivisioncode')
    ..aOS(485287065, _omitFieldNames ? '' : 'countrycode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoLocationDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoLocationDetails copyWith(void Function(GeoLocationDetails) updates) =>
      super.copyWith((message) => updates(message as GeoLocationDetails))
          as GeoLocationDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoLocationDetails create() => GeoLocationDetails._();
  @$core.override
  GeoLocationDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GeoLocationDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GeoLocationDetails>(create);
  static GeoLocationDetails? _defaultInstance;

  @$pb.TagNumber(16070667)
  $core.String get countryname => $_getSZ(0);
  @$pb.TagNumber(16070667)
  set countryname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(16070667)
  $core.bool hasCountryname() => $_has(0);
  @$pb.TagNumber(16070667)
  void clearCountryname() => $_clearField(16070667);

  @$pb.TagNumber(180194347)
  $core.String get continentcode => $_getSZ(1);
  @$pb.TagNumber(180194347)
  set continentcode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(180194347)
  $core.bool hasContinentcode() => $_has(1);
  @$pb.TagNumber(180194347)
  void clearContinentcode() => $_clearField(180194347);

  @$pb.TagNumber(203176953)
  $core.String get continentname => $_getSZ(2);
  @$pb.TagNumber(203176953)
  set continentname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(203176953)
  $core.bool hasContinentname() => $_has(2);
  @$pb.TagNumber(203176953)
  void clearContinentname() => $_clearField(203176953);

  @$pb.TagNumber(412782294)
  $core.String get subdivisionname => $_getSZ(3);
  @$pb.TagNumber(412782294)
  set subdivisionname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(412782294)
  $core.bool hasSubdivisionname() => $_has(3);
  @$pb.TagNumber(412782294)
  void clearSubdivisionname() => $_clearField(412782294);

  @$pb.TagNumber(444328268)
  $core.String get subdivisioncode => $_getSZ(4);
  @$pb.TagNumber(444328268)
  set subdivisioncode($core.String value) => $_setString(4, value);
  @$pb.TagNumber(444328268)
  $core.bool hasSubdivisioncode() => $_has(4);
  @$pb.TagNumber(444328268)
  void clearSubdivisioncode() => $_clearField(444328268);

  @$pb.TagNumber(485287065)
  $core.String get countrycode => $_getSZ(5);
  @$pb.TagNumber(485287065)
  set countrycode($core.String value) => $_setString(5, value);
  @$pb.TagNumber(485287065)
  $core.bool hasCountrycode() => $_has(5);
  @$pb.TagNumber(485287065)
  void clearCountrycode() => $_clearField(485287065);
}

class GeoProximityLocation extends $pb.GeneratedMessage {
  factory GeoProximityLocation({
    $core.String? localzonegroup,
    $core.int? bias,
    Coordinates? coordinates,
    $core.String? awsregion,
  }) {
    final result = create();
    if (localzonegroup != null) result.localzonegroup = localzonegroup;
    if (bias != null) result.bias = bias;
    if (coordinates != null) result.coordinates = coordinates;
    if (awsregion != null) result.awsregion = awsregion;
    return result;
  }

  GeoProximityLocation._();

  factory GeoProximityLocation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GeoProximityLocation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GeoProximityLocation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(60150354, _omitFieldNames ? '' : 'localzonegroup')
    ..aI(60849893, _omitFieldNames ? '' : 'bias')
    ..aOM<Coordinates>(231283401, _omitFieldNames ? '' : 'coordinates',
        subBuilder: Coordinates.create)
    ..aOS(245430451, _omitFieldNames ? '' : 'awsregion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoProximityLocation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GeoProximityLocation copyWith(void Function(GeoProximityLocation) updates) =>
      super.copyWith((message) => updates(message as GeoProximityLocation))
          as GeoProximityLocation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GeoProximityLocation create() => GeoProximityLocation._();
  @$core.override
  GeoProximityLocation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GeoProximityLocation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GeoProximityLocation>(create);
  static GeoProximityLocation? _defaultInstance;

  @$pb.TagNumber(60150354)
  $core.String get localzonegroup => $_getSZ(0);
  @$pb.TagNumber(60150354)
  set localzonegroup($core.String value) => $_setString(0, value);
  @$pb.TagNumber(60150354)
  $core.bool hasLocalzonegroup() => $_has(0);
  @$pb.TagNumber(60150354)
  void clearLocalzonegroup() => $_clearField(60150354);

  @$pb.TagNumber(60849893)
  $core.int get bias => $_getIZ(1);
  @$pb.TagNumber(60849893)
  set bias($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(60849893)
  $core.bool hasBias() => $_has(1);
  @$pb.TagNumber(60849893)
  void clearBias() => $_clearField(60849893);

  @$pb.TagNumber(231283401)
  Coordinates get coordinates => $_getN(2);
  @$pb.TagNumber(231283401)
  set coordinates(Coordinates value) => $_setField(231283401, value);
  @$pb.TagNumber(231283401)
  $core.bool hasCoordinates() => $_has(2);
  @$pb.TagNumber(231283401)
  void clearCoordinates() => $_clearField(231283401);
  @$pb.TagNumber(231283401)
  Coordinates ensureCoordinates() => $_ensure(2);

  @$pb.TagNumber(245430451)
  $core.String get awsregion => $_getSZ(3);
  @$pb.TagNumber(245430451)
  set awsregion($core.String value) => $_setString(3, value);
  @$pb.TagNumber(245430451)
  $core.bool hasAwsregion() => $_has(3);
  @$pb.TagNumber(245430451)
  void clearAwsregion() => $_clearField(245430451);
}

class GetAccountLimitRequest extends $pb.GeneratedMessage {
  factory GetAccountLimitRequest({
    AccountLimitType? type,
  }) {
    final result = create();
    if (type != null) result.type = type;
    return result;
  }

  GetAccountLimitRequest._();

  factory GetAccountLimitRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAccountLimitRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAccountLimitRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<AccountLimitType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: AccountLimitType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountLimitRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountLimitRequest copyWith(
          void Function(GetAccountLimitRequest) updates) =>
      super.copyWith((message) => updates(message as GetAccountLimitRequest))
          as GetAccountLimitRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccountLimitRequest create() => GetAccountLimitRequest._();
  @$core.override
  GetAccountLimitRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAccountLimitRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAccountLimitRequest>(create);
  static GetAccountLimitRequest? _defaultInstance;

  @$pb.TagNumber(290836590)
  AccountLimitType get type => $_getN(0);
  @$pb.TagNumber(290836590)
  set type(AccountLimitType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class GetAccountLimitResponse extends $pb.GeneratedMessage {
  factory GetAccountLimitResponse({
    $fixnum.Int64? count,
    AccountLimit? limit,
  }) {
    final result = create();
    if (count != null) result.count = count;
    if (limit != null) result.limit = limit;
    return result;
  }

  GetAccountLimitResponse._();

  factory GetAccountLimitResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAccountLimitResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAccountLimitResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(31963285, _omitFieldNames ? '' : 'count')
    ..aOM<AccountLimit>(412502741, _omitFieldNames ? '' : 'limit',
        subBuilder: AccountLimit.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountLimitResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountLimitResponse copyWith(
          void Function(GetAccountLimitResponse) updates) =>
      super.copyWith((message) => updates(message as GetAccountLimitResponse))
          as GetAccountLimitResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccountLimitResponse create() => GetAccountLimitResponse._();
  @$core.override
  GetAccountLimitResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAccountLimitResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAccountLimitResponse>(create);
  static GetAccountLimitResponse? _defaultInstance;

  @$pb.TagNumber(31963285)
  $fixnum.Int64 get count => $_getI64(0);
  @$pb.TagNumber(31963285)
  set count($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(31963285)
  $core.bool hasCount() => $_has(0);
  @$pb.TagNumber(31963285)
  void clearCount() => $_clearField(31963285);

  @$pb.TagNumber(412502741)
  AccountLimit get limit => $_getN(1);
  @$pb.TagNumber(412502741)
  set limit(AccountLimit value) => $_setField(412502741, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
  @$pb.TagNumber(412502741)
  AccountLimit ensureLimit() => $_ensure(1);
}

class GetChangeRequest extends $pb.GeneratedMessage {
  factory GetChangeRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  GetChangeRequest._();

  factory GetChangeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetChangeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetChangeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeRequest copyWith(void Function(GetChangeRequest) updates) =>
      super.copyWith((message) => updates(message as GetChangeRequest))
          as GetChangeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChangeRequest create() => GetChangeRequest._();
  @$core.override
  GetChangeRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetChangeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetChangeRequest>(create);
  static GetChangeRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class GetChangeResponse extends $pb.GeneratedMessage {
  factory GetChangeResponse({
    ChangeInfo? changeinfo,
  }) {
    final result = create();
    if (changeinfo != null) result.changeinfo = changeinfo;
    return result;
  }

  GetChangeResponse._();

  factory GetChangeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetChangeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetChangeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ChangeInfo>(436394734, _omitFieldNames ? '' : 'changeinfo',
        subBuilder: ChangeInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetChangeResponse copyWith(void Function(GetChangeResponse) updates) =>
      super.copyWith((message) => updates(message as GetChangeResponse))
          as GetChangeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetChangeResponse create() => GetChangeResponse._();
  @$core.override
  GetChangeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetChangeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetChangeResponse>(create);
  static GetChangeResponse? _defaultInstance;

  @$pb.TagNumber(436394734)
  ChangeInfo get changeinfo => $_getN(0);
  @$pb.TagNumber(436394734)
  set changeinfo(ChangeInfo value) => $_setField(436394734, value);
  @$pb.TagNumber(436394734)
  $core.bool hasChangeinfo() => $_has(0);
  @$pb.TagNumber(436394734)
  void clearChangeinfo() => $_clearField(436394734);
  @$pb.TagNumber(436394734)
  ChangeInfo ensureChangeinfo() => $_ensure(0);
}

class GetCheckerIpRangesRequest extends $pb.GeneratedMessage {
  factory GetCheckerIpRangesRequest() => create();

  GetCheckerIpRangesRequest._();

  factory GetCheckerIpRangesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCheckerIpRangesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCheckerIpRangesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCheckerIpRangesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCheckerIpRangesRequest copyWith(
          void Function(GetCheckerIpRangesRequest) updates) =>
      super.copyWith((message) => updates(message as GetCheckerIpRangesRequest))
          as GetCheckerIpRangesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCheckerIpRangesRequest create() => GetCheckerIpRangesRequest._();
  @$core.override
  GetCheckerIpRangesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCheckerIpRangesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCheckerIpRangesRequest>(create);
  static GetCheckerIpRangesRequest? _defaultInstance;
}

class GetCheckerIpRangesResponse extends $pb.GeneratedMessage {
  factory GetCheckerIpRangesResponse({
    $core.Iterable<$core.String>? checkeripranges,
  }) {
    final result = create();
    if (checkeripranges != null) result.checkeripranges.addAll(checkeripranges);
    return result;
  }

  GetCheckerIpRangesResponse._();

  factory GetCheckerIpRangesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCheckerIpRangesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCheckerIpRangesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPS(376098860, _omitFieldNames ? '' : 'checkeripranges')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCheckerIpRangesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCheckerIpRangesResponse copyWith(
          void Function(GetCheckerIpRangesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetCheckerIpRangesResponse))
          as GetCheckerIpRangesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCheckerIpRangesResponse create() => GetCheckerIpRangesResponse._();
  @$core.override
  GetCheckerIpRangesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCheckerIpRangesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCheckerIpRangesResponse>(create);
  static GetCheckerIpRangesResponse? _defaultInstance;

  @$pb.TagNumber(376098860)
  $pb.PbList<$core.String> get checkeripranges => $_getList(0);
}

class GetDNSSECRequest extends $pb.GeneratedMessage {
  factory GetDNSSECRequest({
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  GetDNSSECRequest._();

  factory GetDNSSECRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDNSSECRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDNSSECRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDNSSECRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDNSSECRequest copyWith(void Function(GetDNSSECRequest) updates) =>
      super.copyWith((message) => updates(message as GetDNSSECRequest))
          as GetDNSSECRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDNSSECRequest create() => GetDNSSECRequest._();
  @$core.override
  GetDNSSECRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDNSSECRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDNSSECRequest>(create);
  static GetDNSSECRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class GetDNSSECResponse extends $pb.GeneratedMessage {
  factory GetDNSSECResponse({
    DNSSECStatus? status,
    $core.Iterable<KeySigningKey>? keysigningkeys,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (keysigningkeys != null) result.keysigningkeys.addAll(keysigningkeys);
    return result;
  }

  GetDNSSECResponse._();

  factory GetDNSSECResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDNSSECResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDNSSECResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<DNSSECStatus>(6222352, _omitFieldNames ? '' : 'status',
        subBuilder: DNSSECStatus.create)
    ..pPM<KeySigningKey>(87847996, _omitFieldNames ? '' : 'keysigningkeys',
        subBuilder: KeySigningKey.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDNSSECResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDNSSECResponse copyWith(void Function(GetDNSSECResponse) updates) =>
      super.copyWith((message) => updates(message as GetDNSSECResponse))
          as GetDNSSECResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDNSSECResponse create() => GetDNSSECResponse._();
  @$core.override
  GetDNSSECResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDNSSECResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDNSSECResponse>(create);
  static GetDNSSECResponse? _defaultInstance;

  @$pb.TagNumber(6222352)
  DNSSECStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(DNSSECStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);
  @$pb.TagNumber(6222352)
  DNSSECStatus ensureStatus() => $_ensure(0);

  @$pb.TagNumber(87847996)
  $pb.PbList<KeySigningKey> get keysigningkeys => $_getList(1);
}

class GetGeoLocationRequest extends $pb.GeneratedMessage {
  factory GetGeoLocationRequest({
    $core.String? continentcode,
    $core.String? subdivisioncode,
    $core.String? countrycode,
  }) {
    final result = create();
    if (continentcode != null) result.continentcode = continentcode;
    if (subdivisioncode != null) result.subdivisioncode = subdivisioncode;
    if (countrycode != null) result.countrycode = countrycode;
    return result;
  }

  GetGeoLocationRequest._();

  factory GetGeoLocationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetGeoLocationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetGeoLocationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(180194347, _omitFieldNames ? '' : 'continentcode')
    ..aOS(444328268, _omitFieldNames ? '' : 'subdivisioncode')
    ..aOS(485287065, _omitFieldNames ? '' : 'countrycode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGeoLocationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGeoLocationRequest copyWith(
          void Function(GetGeoLocationRequest) updates) =>
      super.copyWith((message) => updates(message as GetGeoLocationRequest))
          as GetGeoLocationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetGeoLocationRequest create() => GetGeoLocationRequest._();
  @$core.override
  GetGeoLocationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetGeoLocationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetGeoLocationRequest>(create);
  static GetGeoLocationRequest? _defaultInstance;

  @$pb.TagNumber(180194347)
  $core.String get continentcode => $_getSZ(0);
  @$pb.TagNumber(180194347)
  set continentcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(180194347)
  $core.bool hasContinentcode() => $_has(0);
  @$pb.TagNumber(180194347)
  void clearContinentcode() => $_clearField(180194347);

  @$pb.TagNumber(444328268)
  $core.String get subdivisioncode => $_getSZ(1);
  @$pb.TagNumber(444328268)
  set subdivisioncode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(444328268)
  $core.bool hasSubdivisioncode() => $_has(1);
  @$pb.TagNumber(444328268)
  void clearSubdivisioncode() => $_clearField(444328268);

  @$pb.TagNumber(485287065)
  $core.String get countrycode => $_getSZ(2);
  @$pb.TagNumber(485287065)
  set countrycode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(485287065)
  $core.bool hasCountrycode() => $_has(2);
  @$pb.TagNumber(485287065)
  void clearCountrycode() => $_clearField(485287065);
}

class GetGeoLocationResponse extends $pb.GeneratedMessage {
  factory GetGeoLocationResponse({
    GeoLocationDetails? geolocationdetails,
  }) {
    final result = create();
    if (geolocationdetails != null)
      result.geolocationdetails = geolocationdetails;
    return result;
  }

  GetGeoLocationResponse._();

  factory GetGeoLocationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetGeoLocationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetGeoLocationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<GeoLocationDetails>(
        70816602, _omitFieldNames ? '' : 'geolocationdetails',
        subBuilder: GeoLocationDetails.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGeoLocationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetGeoLocationResponse copyWith(
          void Function(GetGeoLocationResponse) updates) =>
      super.copyWith((message) => updates(message as GetGeoLocationResponse))
          as GetGeoLocationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetGeoLocationResponse create() => GetGeoLocationResponse._();
  @$core.override
  GetGeoLocationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetGeoLocationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetGeoLocationResponse>(create);
  static GetGeoLocationResponse? _defaultInstance;

  @$pb.TagNumber(70816602)
  GeoLocationDetails get geolocationdetails => $_getN(0);
  @$pb.TagNumber(70816602)
  set geolocationdetails(GeoLocationDetails value) =>
      $_setField(70816602, value);
  @$pb.TagNumber(70816602)
  $core.bool hasGeolocationdetails() => $_has(0);
  @$pb.TagNumber(70816602)
  void clearGeolocationdetails() => $_clearField(70816602);
  @$pb.TagNumber(70816602)
  GeoLocationDetails ensureGeolocationdetails() => $_ensure(0);
}

class GetHealthCheckCountRequest extends $pb.GeneratedMessage {
  factory GetHealthCheckCountRequest() => create();

  GetHealthCheckCountRequest._();

  factory GetHealthCheckCountRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHealthCheckCountRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHealthCheckCountRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckCountRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckCountRequest copyWith(
          void Function(GetHealthCheckCountRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetHealthCheckCountRequest))
          as GetHealthCheckCountRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHealthCheckCountRequest create() => GetHealthCheckCountRequest._();
  @$core.override
  GetHealthCheckCountRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHealthCheckCountRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHealthCheckCountRequest>(create);
  static GetHealthCheckCountRequest? _defaultInstance;
}

class GetHealthCheckCountResponse extends $pb.GeneratedMessage {
  factory GetHealthCheckCountResponse({
    $fixnum.Int64? healthcheckcount,
  }) {
    final result = create();
    if (healthcheckcount != null) result.healthcheckcount = healthcheckcount;
    return result;
  }

  GetHealthCheckCountResponse._();

  factory GetHealthCheckCountResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHealthCheckCountResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHealthCheckCountResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(332677785, _omitFieldNames ? '' : 'healthcheckcount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckCountResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckCountResponse copyWith(
          void Function(GetHealthCheckCountResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetHealthCheckCountResponse))
          as GetHealthCheckCountResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHealthCheckCountResponse create() =>
      GetHealthCheckCountResponse._();
  @$core.override
  GetHealthCheckCountResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHealthCheckCountResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHealthCheckCountResponse>(create);
  static GetHealthCheckCountResponse? _defaultInstance;

  @$pb.TagNumber(332677785)
  $fixnum.Int64 get healthcheckcount => $_getI64(0);
  @$pb.TagNumber(332677785)
  set healthcheckcount($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(332677785)
  $core.bool hasHealthcheckcount() => $_has(0);
  @$pb.TagNumber(332677785)
  void clearHealthcheckcount() => $_clearField(332677785);
}

class GetHealthCheckLastFailureReasonRequest extends $pb.GeneratedMessage {
  factory GetHealthCheckLastFailureReasonRequest({
    $core.String? healthcheckid,
  }) {
    final result = create();
    if (healthcheckid != null) result.healthcheckid = healthcheckid;
    return result;
  }

  GetHealthCheckLastFailureReasonRequest._();

  factory GetHealthCheckLastFailureReasonRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHealthCheckLastFailureReasonRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHealthCheckLastFailureReasonRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(312971637, _omitFieldNames ? '' : 'healthcheckid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckLastFailureReasonRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckLastFailureReasonRequest copyWith(
          void Function(GetHealthCheckLastFailureReasonRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetHealthCheckLastFailureReasonRequest))
          as GetHealthCheckLastFailureReasonRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHealthCheckLastFailureReasonRequest create() =>
      GetHealthCheckLastFailureReasonRequest._();
  @$core.override
  GetHealthCheckLastFailureReasonRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHealthCheckLastFailureReasonRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetHealthCheckLastFailureReasonRequest>(create);
  static GetHealthCheckLastFailureReasonRequest? _defaultInstance;

  @$pb.TagNumber(312971637)
  $core.String get healthcheckid => $_getSZ(0);
  @$pb.TagNumber(312971637)
  set healthcheckid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(312971637)
  $core.bool hasHealthcheckid() => $_has(0);
  @$pb.TagNumber(312971637)
  void clearHealthcheckid() => $_clearField(312971637);
}

class GetHealthCheckLastFailureReasonResponse extends $pb.GeneratedMessage {
  factory GetHealthCheckLastFailureReasonResponse({
    $core.Iterable<HealthCheckObservation>? healthcheckobservations,
  }) {
    final result = create();
    if (healthcheckobservations != null)
      result.healthcheckobservations.addAll(healthcheckobservations);
    return result;
  }

  GetHealthCheckLastFailureReasonResponse._();

  factory GetHealthCheckLastFailureReasonResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHealthCheckLastFailureReasonResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHealthCheckLastFailureReasonResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<HealthCheckObservation>(
        391890103, _omitFieldNames ? '' : 'healthcheckobservations',
        subBuilder: HealthCheckObservation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckLastFailureReasonResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckLastFailureReasonResponse copyWith(
          void Function(GetHealthCheckLastFailureReasonResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetHealthCheckLastFailureReasonResponse))
          as GetHealthCheckLastFailureReasonResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHealthCheckLastFailureReasonResponse create() =>
      GetHealthCheckLastFailureReasonResponse._();
  @$core.override
  GetHealthCheckLastFailureReasonResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHealthCheckLastFailureReasonResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetHealthCheckLastFailureReasonResponse>(create);
  static GetHealthCheckLastFailureReasonResponse? _defaultInstance;

  @$pb.TagNumber(391890103)
  $pb.PbList<HealthCheckObservation> get healthcheckobservations =>
      $_getList(0);
}

class GetHealthCheckRequest extends $pb.GeneratedMessage {
  factory GetHealthCheckRequest({
    $core.String? healthcheckid,
  }) {
    final result = create();
    if (healthcheckid != null) result.healthcheckid = healthcheckid;
    return result;
  }

  GetHealthCheckRequest._();

  factory GetHealthCheckRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHealthCheckRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHealthCheckRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(312971637, _omitFieldNames ? '' : 'healthcheckid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckRequest copyWith(
          void Function(GetHealthCheckRequest) updates) =>
      super.copyWith((message) => updates(message as GetHealthCheckRequest))
          as GetHealthCheckRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHealthCheckRequest create() => GetHealthCheckRequest._();
  @$core.override
  GetHealthCheckRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHealthCheckRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHealthCheckRequest>(create);
  static GetHealthCheckRequest? _defaultInstance;

  @$pb.TagNumber(312971637)
  $core.String get healthcheckid => $_getSZ(0);
  @$pb.TagNumber(312971637)
  set healthcheckid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(312971637)
  $core.bool hasHealthcheckid() => $_has(0);
  @$pb.TagNumber(312971637)
  void clearHealthcheckid() => $_clearField(312971637);
}

class GetHealthCheckResponse extends $pb.GeneratedMessage {
  factory GetHealthCheckResponse({
    HealthCheck? healthcheck,
  }) {
    final result = create();
    if (healthcheck != null) result.healthcheck = healthcheck;
    return result;
  }

  GetHealthCheckResponse._();

  factory GetHealthCheckResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHealthCheckResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHealthCheckResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<HealthCheck>(377540610, _omitFieldNames ? '' : 'healthcheck',
        subBuilder: HealthCheck.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckResponse copyWith(
          void Function(GetHealthCheckResponse) updates) =>
      super.copyWith((message) => updates(message as GetHealthCheckResponse))
          as GetHealthCheckResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHealthCheckResponse create() => GetHealthCheckResponse._();
  @$core.override
  GetHealthCheckResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHealthCheckResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHealthCheckResponse>(create);
  static GetHealthCheckResponse? _defaultInstance;

  @$pb.TagNumber(377540610)
  HealthCheck get healthcheck => $_getN(0);
  @$pb.TagNumber(377540610)
  set healthcheck(HealthCheck value) => $_setField(377540610, value);
  @$pb.TagNumber(377540610)
  $core.bool hasHealthcheck() => $_has(0);
  @$pb.TagNumber(377540610)
  void clearHealthcheck() => $_clearField(377540610);
  @$pb.TagNumber(377540610)
  HealthCheck ensureHealthcheck() => $_ensure(0);
}

class GetHealthCheckStatusRequest extends $pb.GeneratedMessage {
  factory GetHealthCheckStatusRequest({
    $core.String? healthcheckid,
  }) {
    final result = create();
    if (healthcheckid != null) result.healthcheckid = healthcheckid;
    return result;
  }

  GetHealthCheckStatusRequest._();

  factory GetHealthCheckStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHealthCheckStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHealthCheckStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(312971637, _omitFieldNames ? '' : 'healthcheckid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckStatusRequest copyWith(
          void Function(GetHealthCheckStatusRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetHealthCheckStatusRequest))
          as GetHealthCheckStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHealthCheckStatusRequest create() =>
      GetHealthCheckStatusRequest._();
  @$core.override
  GetHealthCheckStatusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHealthCheckStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHealthCheckStatusRequest>(create);
  static GetHealthCheckStatusRequest? _defaultInstance;

  @$pb.TagNumber(312971637)
  $core.String get healthcheckid => $_getSZ(0);
  @$pb.TagNumber(312971637)
  set healthcheckid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(312971637)
  $core.bool hasHealthcheckid() => $_has(0);
  @$pb.TagNumber(312971637)
  void clearHealthcheckid() => $_clearField(312971637);
}

class GetHealthCheckStatusResponse extends $pb.GeneratedMessage {
  factory GetHealthCheckStatusResponse({
    $core.Iterable<HealthCheckObservation>? healthcheckobservations,
  }) {
    final result = create();
    if (healthcheckobservations != null)
      result.healthcheckobservations.addAll(healthcheckobservations);
    return result;
  }

  GetHealthCheckStatusResponse._();

  factory GetHealthCheckStatusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHealthCheckStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHealthCheckStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<HealthCheckObservation>(
        391890103, _omitFieldNames ? '' : 'healthcheckobservations',
        subBuilder: HealthCheckObservation.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHealthCheckStatusResponse copyWith(
          void Function(GetHealthCheckStatusResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetHealthCheckStatusResponse))
          as GetHealthCheckStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHealthCheckStatusResponse create() =>
      GetHealthCheckStatusResponse._();
  @$core.override
  GetHealthCheckStatusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHealthCheckStatusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHealthCheckStatusResponse>(create);
  static GetHealthCheckStatusResponse? _defaultInstance;

  @$pb.TagNumber(391890103)
  $pb.PbList<HealthCheckObservation> get healthcheckobservations =>
      $_getList(0);
}

class GetHostedZoneCountRequest extends $pb.GeneratedMessage {
  factory GetHostedZoneCountRequest() => create();

  GetHostedZoneCountRequest._();

  factory GetHostedZoneCountRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHostedZoneCountRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHostedZoneCountRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneCountRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneCountRequest copyWith(
          void Function(GetHostedZoneCountRequest) updates) =>
      super.copyWith((message) => updates(message as GetHostedZoneCountRequest))
          as GetHostedZoneCountRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHostedZoneCountRequest create() => GetHostedZoneCountRequest._();
  @$core.override
  GetHostedZoneCountRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHostedZoneCountRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHostedZoneCountRequest>(create);
  static GetHostedZoneCountRequest? _defaultInstance;
}

class GetHostedZoneCountResponse extends $pb.GeneratedMessage {
  factory GetHostedZoneCountResponse({
    $fixnum.Int64? hostedzonecount,
  }) {
    final result = create();
    if (hostedzonecount != null) result.hostedzonecount = hostedzonecount;
    return result;
  }

  GetHostedZoneCountResponse._();

  factory GetHostedZoneCountResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHostedZoneCountResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHostedZoneCountResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(93723536, _omitFieldNames ? '' : 'hostedzonecount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneCountResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneCountResponse copyWith(
          void Function(GetHostedZoneCountResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetHostedZoneCountResponse))
          as GetHostedZoneCountResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHostedZoneCountResponse create() => GetHostedZoneCountResponse._();
  @$core.override
  GetHostedZoneCountResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHostedZoneCountResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHostedZoneCountResponse>(create);
  static GetHostedZoneCountResponse? _defaultInstance;

  @$pb.TagNumber(93723536)
  $fixnum.Int64 get hostedzonecount => $_getI64(0);
  @$pb.TagNumber(93723536)
  set hostedzonecount($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(93723536)
  $core.bool hasHostedzonecount() => $_has(0);
  @$pb.TagNumber(93723536)
  void clearHostedzonecount() => $_clearField(93723536);
}

class GetHostedZoneLimitRequest extends $pb.GeneratedMessage {
  factory GetHostedZoneLimitRequest({
    HostedZoneLimitType? type,
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  GetHostedZoneLimitRequest._();

  factory GetHostedZoneLimitRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHostedZoneLimitRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHostedZoneLimitRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<HostedZoneLimitType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: HostedZoneLimitType.values)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneLimitRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneLimitRequest copyWith(
          void Function(GetHostedZoneLimitRequest) updates) =>
      super.copyWith((message) => updates(message as GetHostedZoneLimitRequest))
          as GetHostedZoneLimitRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHostedZoneLimitRequest create() => GetHostedZoneLimitRequest._();
  @$core.override
  GetHostedZoneLimitRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHostedZoneLimitRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHostedZoneLimitRequest>(create);
  static GetHostedZoneLimitRequest? _defaultInstance;

  @$pb.TagNumber(290836590)
  HostedZoneLimitType get type => $_getN(0);
  @$pb.TagNumber(290836590)
  set type(HostedZoneLimitType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class GetHostedZoneLimitResponse extends $pb.GeneratedMessage {
  factory GetHostedZoneLimitResponse({
    $fixnum.Int64? count,
    HostedZoneLimit? limit,
  }) {
    final result = create();
    if (count != null) result.count = count;
    if (limit != null) result.limit = limit;
    return result;
  }

  GetHostedZoneLimitResponse._();

  factory GetHostedZoneLimitResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHostedZoneLimitResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHostedZoneLimitResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(31963285, _omitFieldNames ? '' : 'count')
    ..aOM<HostedZoneLimit>(412502741, _omitFieldNames ? '' : 'limit',
        subBuilder: HostedZoneLimit.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneLimitResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneLimitResponse copyWith(
          void Function(GetHostedZoneLimitResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetHostedZoneLimitResponse))
          as GetHostedZoneLimitResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHostedZoneLimitResponse create() => GetHostedZoneLimitResponse._();
  @$core.override
  GetHostedZoneLimitResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHostedZoneLimitResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHostedZoneLimitResponse>(create);
  static GetHostedZoneLimitResponse? _defaultInstance;

  @$pb.TagNumber(31963285)
  $fixnum.Int64 get count => $_getI64(0);
  @$pb.TagNumber(31963285)
  set count($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(31963285)
  $core.bool hasCount() => $_has(0);
  @$pb.TagNumber(31963285)
  void clearCount() => $_clearField(31963285);

  @$pb.TagNumber(412502741)
  HostedZoneLimit get limit => $_getN(1);
  @$pb.TagNumber(412502741)
  set limit(HostedZoneLimit value) => $_setField(412502741, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
  @$pb.TagNumber(412502741)
  HostedZoneLimit ensureLimit() => $_ensure(1);
}

class GetHostedZoneRequest extends $pb.GeneratedMessage {
  factory GetHostedZoneRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  GetHostedZoneRequest._();

  factory GetHostedZoneRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHostedZoneRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHostedZoneRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneRequest copyWith(void Function(GetHostedZoneRequest) updates) =>
      super.copyWith((message) => updates(message as GetHostedZoneRequest))
          as GetHostedZoneRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHostedZoneRequest create() => GetHostedZoneRequest._();
  @$core.override
  GetHostedZoneRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHostedZoneRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHostedZoneRequest>(create);
  static GetHostedZoneRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class GetHostedZoneResponse extends $pb.GeneratedMessage {
  factory GetHostedZoneResponse({
    DelegationSet? delegationset,
    $core.Iterable<VPC>? vpcs,
    HostedZone? hostedzone,
  }) {
    final result = create();
    if (delegationset != null) result.delegationset = delegationset;
    if (vpcs != null) result.vpcs.addAll(vpcs);
    if (hostedzone != null) result.hostedzone = hostedzone;
    return result;
  }

  GetHostedZoneResponse._();

  factory GetHostedZoneResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetHostedZoneResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetHostedZoneResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<DelegationSet>(190771750, _omitFieldNames ? '' : 'delegationset',
        subBuilder: DelegationSet.create)
    ..pPM<VPC>(424064898, _omitFieldNames ? '' : 'vpcs', subBuilder: VPC.create)
    ..aOM<HostedZone>(465689249, _omitFieldNames ? '' : 'hostedzone',
        subBuilder: HostedZone.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetHostedZoneResponse copyWith(
          void Function(GetHostedZoneResponse) updates) =>
      super.copyWith((message) => updates(message as GetHostedZoneResponse))
          as GetHostedZoneResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetHostedZoneResponse create() => GetHostedZoneResponse._();
  @$core.override
  GetHostedZoneResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetHostedZoneResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetHostedZoneResponse>(create);
  static GetHostedZoneResponse? _defaultInstance;

  @$pb.TagNumber(190771750)
  DelegationSet get delegationset => $_getN(0);
  @$pb.TagNumber(190771750)
  set delegationset(DelegationSet value) => $_setField(190771750, value);
  @$pb.TagNumber(190771750)
  $core.bool hasDelegationset() => $_has(0);
  @$pb.TagNumber(190771750)
  void clearDelegationset() => $_clearField(190771750);
  @$pb.TagNumber(190771750)
  DelegationSet ensureDelegationset() => $_ensure(0);

  @$pb.TagNumber(424064898)
  $pb.PbList<VPC> get vpcs => $_getList(1);

  @$pb.TagNumber(465689249)
  HostedZone get hostedzone => $_getN(2);
  @$pb.TagNumber(465689249)
  set hostedzone(HostedZone value) => $_setField(465689249, value);
  @$pb.TagNumber(465689249)
  $core.bool hasHostedzone() => $_has(2);
  @$pb.TagNumber(465689249)
  void clearHostedzone() => $_clearField(465689249);
  @$pb.TagNumber(465689249)
  HostedZone ensureHostedzone() => $_ensure(2);
}

class GetQueryLoggingConfigRequest extends $pb.GeneratedMessage {
  factory GetQueryLoggingConfigRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  GetQueryLoggingConfigRequest._();

  factory GetQueryLoggingConfigRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryLoggingConfigRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryLoggingConfigRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryLoggingConfigRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryLoggingConfigRequest copyWith(
          void Function(GetQueryLoggingConfigRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetQueryLoggingConfigRequest))
          as GetQueryLoggingConfigRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryLoggingConfigRequest create() =>
      GetQueryLoggingConfigRequest._();
  @$core.override
  GetQueryLoggingConfigRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryLoggingConfigRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryLoggingConfigRequest>(create);
  static GetQueryLoggingConfigRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class GetQueryLoggingConfigResponse extends $pb.GeneratedMessage {
  factory GetQueryLoggingConfigResponse({
    QueryLoggingConfig? queryloggingconfig,
  }) {
    final result = create();
    if (queryloggingconfig != null)
      result.queryloggingconfig = queryloggingconfig;
    return result;
  }

  GetQueryLoggingConfigResponse._();

  factory GetQueryLoggingConfigResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetQueryLoggingConfigResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetQueryLoggingConfigResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<QueryLoggingConfig>(
        492751675, _omitFieldNames ? '' : 'queryloggingconfig',
        subBuilder: QueryLoggingConfig.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryLoggingConfigResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetQueryLoggingConfigResponse copyWith(
          void Function(GetQueryLoggingConfigResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetQueryLoggingConfigResponse))
          as GetQueryLoggingConfigResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetQueryLoggingConfigResponse create() =>
      GetQueryLoggingConfigResponse._();
  @$core.override
  GetQueryLoggingConfigResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetQueryLoggingConfigResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetQueryLoggingConfigResponse>(create);
  static GetQueryLoggingConfigResponse? _defaultInstance;

  @$pb.TagNumber(492751675)
  QueryLoggingConfig get queryloggingconfig => $_getN(0);
  @$pb.TagNumber(492751675)
  set queryloggingconfig(QueryLoggingConfig value) =>
      $_setField(492751675, value);
  @$pb.TagNumber(492751675)
  $core.bool hasQueryloggingconfig() => $_has(0);
  @$pb.TagNumber(492751675)
  void clearQueryloggingconfig() => $_clearField(492751675);
  @$pb.TagNumber(492751675)
  QueryLoggingConfig ensureQueryloggingconfig() => $_ensure(0);
}

class GetReusableDelegationSetLimitRequest extends $pb.GeneratedMessage {
  factory GetReusableDelegationSetLimitRequest({
    ReusableDelegationSetLimitType? type,
    $core.String? delegationsetid,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (delegationsetid != null) result.delegationsetid = delegationsetid;
    return result;
  }

  GetReusableDelegationSetLimitRequest._();

  factory GetReusableDelegationSetLimitRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetReusableDelegationSetLimitRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetReusableDelegationSetLimitRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<ReusableDelegationSetLimitType>(
        290836590, _omitFieldNames ? '' : 'type',
        enumValues: ReusableDelegationSetLimitType.values)
    ..aOS(307328801, _omitFieldNames ? '' : 'delegationsetid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetReusableDelegationSetLimitRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetReusableDelegationSetLimitRequest copyWith(
          void Function(GetReusableDelegationSetLimitRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetReusableDelegationSetLimitRequest))
          as GetReusableDelegationSetLimitRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetReusableDelegationSetLimitRequest create() =>
      GetReusableDelegationSetLimitRequest._();
  @$core.override
  GetReusableDelegationSetLimitRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetReusableDelegationSetLimitRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetReusableDelegationSetLimitRequest>(create);
  static GetReusableDelegationSetLimitRequest? _defaultInstance;

  @$pb.TagNumber(290836590)
  ReusableDelegationSetLimitType get type => $_getN(0);
  @$pb.TagNumber(290836590)
  set type(ReusableDelegationSetLimitType value) =>
      $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(307328801)
  $core.String get delegationsetid => $_getSZ(1);
  @$pb.TagNumber(307328801)
  set delegationsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(307328801)
  $core.bool hasDelegationsetid() => $_has(1);
  @$pb.TagNumber(307328801)
  void clearDelegationsetid() => $_clearField(307328801);
}

class GetReusableDelegationSetLimitResponse extends $pb.GeneratedMessage {
  factory GetReusableDelegationSetLimitResponse({
    $fixnum.Int64? count,
    ReusableDelegationSetLimit? limit,
  }) {
    final result = create();
    if (count != null) result.count = count;
    if (limit != null) result.limit = limit;
    return result;
  }

  GetReusableDelegationSetLimitResponse._();

  factory GetReusableDelegationSetLimitResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetReusableDelegationSetLimitResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetReusableDelegationSetLimitResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(31963285, _omitFieldNames ? '' : 'count')
    ..aOM<ReusableDelegationSetLimit>(412502741, _omitFieldNames ? '' : 'limit',
        subBuilder: ReusableDelegationSetLimit.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetReusableDelegationSetLimitResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetReusableDelegationSetLimitResponse copyWith(
          void Function(GetReusableDelegationSetLimitResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetReusableDelegationSetLimitResponse))
          as GetReusableDelegationSetLimitResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetReusableDelegationSetLimitResponse create() =>
      GetReusableDelegationSetLimitResponse._();
  @$core.override
  GetReusableDelegationSetLimitResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetReusableDelegationSetLimitResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetReusableDelegationSetLimitResponse>(create);
  static GetReusableDelegationSetLimitResponse? _defaultInstance;

  @$pb.TagNumber(31963285)
  $fixnum.Int64 get count => $_getI64(0);
  @$pb.TagNumber(31963285)
  set count($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(31963285)
  $core.bool hasCount() => $_has(0);
  @$pb.TagNumber(31963285)
  void clearCount() => $_clearField(31963285);

  @$pb.TagNumber(412502741)
  ReusableDelegationSetLimit get limit => $_getN(1);
  @$pb.TagNumber(412502741)
  set limit(ReusableDelegationSetLimit value) => $_setField(412502741, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
  @$pb.TagNumber(412502741)
  ReusableDelegationSetLimit ensureLimit() => $_ensure(1);
}

class GetReusableDelegationSetRequest extends $pb.GeneratedMessage {
  factory GetReusableDelegationSetRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  GetReusableDelegationSetRequest._();

  factory GetReusableDelegationSetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetReusableDelegationSetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetReusableDelegationSetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetReusableDelegationSetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetReusableDelegationSetRequest copyWith(
          void Function(GetReusableDelegationSetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetReusableDelegationSetRequest))
          as GetReusableDelegationSetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetReusableDelegationSetRequest create() =>
      GetReusableDelegationSetRequest._();
  @$core.override
  GetReusableDelegationSetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetReusableDelegationSetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetReusableDelegationSetRequest>(
          create);
  static GetReusableDelegationSetRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class GetReusableDelegationSetResponse extends $pb.GeneratedMessage {
  factory GetReusableDelegationSetResponse({
    DelegationSet? delegationset,
  }) {
    final result = create();
    if (delegationset != null) result.delegationset = delegationset;
    return result;
  }

  GetReusableDelegationSetResponse._();

  factory GetReusableDelegationSetResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetReusableDelegationSetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetReusableDelegationSetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<DelegationSet>(190771750, _omitFieldNames ? '' : 'delegationset',
        subBuilder: DelegationSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetReusableDelegationSetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetReusableDelegationSetResponse copyWith(
          void Function(GetReusableDelegationSetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetReusableDelegationSetResponse))
          as GetReusableDelegationSetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetReusableDelegationSetResponse create() =>
      GetReusableDelegationSetResponse._();
  @$core.override
  GetReusableDelegationSetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetReusableDelegationSetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetReusableDelegationSetResponse>(
          create);
  static GetReusableDelegationSetResponse? _defaultInstance;

  @$pb.TagNumber(190771750)
  DelegationSet get delegationset => $_getN(0);
  @$pb.TagNumber(190771750)
  set delegationset(DelegationSet value) => $_setField(190771750, value);
  @$pb.TagNumber(190771750)
  $core.bool hasDelegationset() => $_has(0);
  @$pb.TagNumber(190771750)
  void clearDelegationset() => $_clearField(190771750);
  @$pb.TagNumber(190771750)
  DelegationSet ensureDelegationset() => $_ensure(0);
}

class GetTrafficPolicyInstanceCountRequest extends $pb.GeneratedMessage {
  factory GetTrafficPolicyInstanceCountRequest() => create();

  GetTrafficPolicyInstanceCountRequest._();

  factory GetTrafficPolicyInstanceCountRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrafficPolicyInstanceCountRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrafficPolicyInstanceCountRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyInstanceCountRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyInstanceCountRequest copyWith(
          void Function(GetTrafficPolicyInstanceCountRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetTrafficPolicyInstanceCountRequest))
          as GetTrafficPolicyInstanceCountRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyInstanceCountRequest create() =>
      GetTrafficPolicyInstanceCountRequest._();
  @$core.override
  GetTrafficPolicyInstanceCountRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyInstanceCountRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetTrafficPolicyInstanceCountRequest>(create);
  static GetTrafficPolicyInstanceCountRequest? _defaultInstance;
}

class GetTrafficPolicyInstanceCountResponse extends $pb.GeneratedMessage {
  factory GetTrafficPolicyInstanceCountResponse({
    $core.int? trafficpolicyinstancecount,
  }) {
    final result = create();
    if (trafficpolicyinstancecount != null)
      result.trafficpolicyinstancecount = trafficpolicyinstancecount;
    return result;
  }

  GetTrafficPolicyInstanceCountResponse._();

  factory GetTrafficPolicyInstanceCountResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrafficPolicyInstanceCountResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrafficPolicyInstanceCountResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aI(190494003, _omitFieldNames ? '' : 'trafficpolicyinstancecount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyInstanceCountResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyInstanceCountResponse copyWith(
          void Function(GetTrafficPolicyInstanceCountResponse) updates) =>
      super.copyWith((message) =>
              updates(message as GetTrafficPolicyInstanceCountResponse))
          as GetTrafficPolicyInstanceCountResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyInstanceCountResponse create() =>
      GetTrafficPolicyInstanceCountResponse._();
  @$core.override
  GetTrafficPolicyInstanceCountResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyInstanceCountResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetTrafficPolicyInstanceCountResponse>(create);
  static GetTrafficPolicyInstanceCountResponse? _defaultInstance;

  @$pb.TagNumber(190494003)
  $core.int get trafficpolicyinstancecount => $_getIZ(0);
  @$pb.TagNumber(190494003)
  set trafficpolicyinstancecount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(190494003)
  $core.bool hasTrafficpolicyinstancecount() => $_has(0);
  @$pb.TagNumber(190494003)
  void clearTrafficpolicyinstancecount() => $_clearField(190494003);
}

class GetTrafficPolicyInstanceRequest extends $pb.GeneratedMessage {
  factory GetTrafficPolicyInstanceRequest({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  GetTrafficPolicyInstanceRequest._();

  factory GetTrafficPolicyInstanceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrafficPolicyInstanceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrafficPolicyInstanceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyInstanceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyInstanceRequest copyWith(
          void Function(GetTrafficPolicyInstanceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetTrafficPolicyInstanceRequest))
          as GetTrafficPolicyInstanceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyInstanceRequest create() =>
      GetTrafficPolicyInstanceRequest._();
  @$core.override
  GetTrafficPolicyInstanceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyInstanceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTrafficPolicyInstanceRequest>(
          create);
  static GetTrafficPolicyInstanceRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);
}

class GetTrafficPolicyInstanceResponse extends $pb.GeneratedMessage {
  factory GetTrafficPolicyInstanceResponse({
    TrafficPolicyInstance? trafficpolicyinstance,
  }) {
    final result = create();
    if (trafficpolicyinstance != null)
      result.trafficpolicyinstance = trafficpolicyinstance;
    return result;
  }

  GetTrafficPolicyInstanceResponse._();

  factory GetTrafficPolicyInstanceResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrafficPolicyInstanceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrafficPolicyInstanceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<TrafficPolicyInstance>(
        205651476, _omitFieldNames ? '' : 'trafficpolicyinstance',
        subBuilder: TrafficPolicyInstance.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyInstanceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyInstanceResponse copyWith(
          void Function(GetTrafficPolicyInstanceResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetTrafficPolicyInstanceResponse))
          as GetTrafficPolicyInstanceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyInstanceResponse create() =>
      GetTrafficPolicyInstanceResponse._();
  @$core.override
  GetTrafficPolicyInstanceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyInstanceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTrafficPolicyInstanceResponse>(
          create);
  static GetTrafficPolicyInstanceResponse? _defaultInstance;

  @$pb.TagNumber(205651476)
  TrafficPolicyInstance get trafficpolicyinstance => $_getN(0);
  @$pb.TagNumber(205651476)
  set trafficpolicyinstance(TrafficPolicyInstance value) =>
      $_setField(205651476, value);
  @$pb.TagNumber(205651476)
  $core.bool hasTrafficpolicyinstance() => $_has(0);
  @$pb.TagNumber(205651476)
  void clearTrafficpolicyinstance() => $_clearField(205651476);
  @$pb.TagNumber(205651476)
  TrafficPolicyInstance ensureTrafficpolicyinstance() => $_ensure(0);
}

class GetTrafficPolicyRequest extends $pb.GeneratedMessage {
  factory GetTrafficPolicyRequest({
    $core.String? id,
    $core.int? version,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (version != null) result.version = version;
    return result;
  }

  GetTrafficPolicyRequest._();

  factory GetTrafficPolicyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrafficPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrafficPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aI(500028728, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyRequest copyWith(
          void Function(GetTrafficPolicyRequest) updates) =>
      super.copyWith((message) => updates(message as GetTrafficPolicyRequest))
          as GetTrafficPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyRequest create() => GetTrafficPolicyRequest._();
  @$core.override
  GetTrafficPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTrafficPolicyRequest>(create);
  static GetTrafficPolicyRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(500028728)
  $core.int get version => $_getIZ(1);
  @$pb.TagNumber(500028728)
  set version($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(1);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);
}

class GetTrafficPolicyResponse extends $pb.GeneratedMessage {
  factory GetTrafficPolicyResponse({
    TrafficPolicy? trafficpolicy,
  }) {
    final result = create();
    if (trafficpolicy != null) result.trafficpolicy = trafficpolicy;
    return result;
  }

  GetTrafficPolicyResponse._();

  factory GetTrafficPolicyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTrafficPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTrafficPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<TrafficPolicy>(154595657, _omitFieldNames ? '' : 'trafficpolicy',
        subBuilder: TrafficPolicy.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTrafficPolicyResponse copyWith(
          void Function(GetTrafficPolicyResponse) updates) =>
      super.copyWith((message) => updates(message as GetTrafficPolicyResponse))
          as GetTrafficPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyResponse create() => GetTrafficPolicyResponse._();
  @$core.override
  GetTrafficPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTrafficPolicyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTrafficPolicyResponse>(create);
  static GetTrafficPolicyResponse? _defaultInstance;

  @$pb.TagNumber(154595657)
  TrafficPolicy get trafficpolicy => $_getN(0);
  @$pb.TagNumber(154595657)
  set trafficpolicy(TrafficPolicy value) => $_setField(154595657, value);
  @$pb.TagNumber(154595657)
  $core.bool hasTrafficpolicy() => $_has(0);
  @$pb.TagNumber(154595657)
  void clearTrafficpolicy() => $_clearField(154595657);
  @$pb.TagNumber(154595657)
  TrafficPolicy ensureTrafficpolicy() => $_ensure(0);
}

class HealthCheck extends $pb.GeneratedMessage {
  factory HealthCheck({
    HealthCheckConfig? healthcheckconfig,
    $fixnum.Int64? healthcheckversion,
    $core.String? callerreference,
    CloudWatchAlarmConfiguration? cloudwatchalarmconfiguration,
    $core.String? id,
    LinkedService? linkedservice,
  }) {
    final result = create();
    if (healthcheckconfig != null) result.healthcheckconfig = healthcheckconfig;
    if (healthcheckversion != null)
      result.healthcheckversion = healthcheckversion;
    if (callerreference != null) result.callerreference = callerreference;
    if (cloudwatchalarmconfiguration != null)
      result.cloudwatchalarmconfiguration = cloudwatchalarmconfiguration;
    if (id != null) result.id = id;
    if (linkedservice != null) result.linkedservice = linkedservice;
    return result;
  }

  HealthCheck._();

  factory HealthCheck.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealthCheck.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealthCheck',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<HealthCheckConfig>(
        83111204, _omitFieldNames ? '' : 'healthcheckconfig',
        subBuilder: HealthCheckConfig.create)
    ..aInt64(89568396, _omitFieldNames ? '' : 'healthcheckversion')
    ..aOS(151211160, _omitFieldNames ? '' : 'callerreference')
    ..aOM<CloudWatchAlarmConfiguration>(
        358915841, _omitFieldNames ? '' : 'cloudwatchalarmconfiguration',
        subBuilder: CloudWatchAlarmConfiguration.create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOM<LinkedService>(438614164, _omitFieldNames ? '' : 'linkedservice',
        subBuilder: LinkedService.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheck clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheck copyWith(void Function(HealthCheck) updates) =>
      super.copyWith((message) => updates(message as HealthCheck))
          as HealthCheck;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealthCheck create() => HealthCheck._();
  @$core.override
  HealthCheck createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealthCheck getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealthCheck>(create);
  static HealthCheck? _defaultInstance;

  @$pb.TagNumber(83111204)
  HealthCheckConfig get healthcheckconfig => $_getN(0);
  @$pb.TagNumber(83111204)
  set healthcheckconfig(HealthCheckConfig value) => $_setField(83111204, value);
  @$pb.TagNumber(83111204)
  $core.bool hasHealthcheckconfig() => $_has(0);
  @$pb.TagNumber(83111204)
  void clearHealthcheckconfig() => $_clearField(83111204);
  @$pb.TagNumber(83111204)
  HealthCheckConfig ensureHealthcheckconfig() => $_ensure(0);

  @$pb.TagNumber(89568396)
  $fixnum.Int64 get healthcheckversion => $_getI64(1);
  @$pb.TagNumber(89568396)
  set healthcheckversion($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(89568396)
  $core.bool hasHealthcheckversion() => $_has(1);
  @$pb.TagNumber(89568396)
  void clearHealthcheckversion() => $_clearField(89568396);

  @$pb.TagNumber(151211160)
  $core.String get callerreference => $_getSZ(2);
  @$pb.TagNumber(151211160)
  set callerreference($core.String value) => $_setString(2, value);
  @$pb.TagNumber(151211160)
  $core.bool hasCallerreference() => $_has(2);
  @$pb.TagNumber(151211160)
  void clearCallerreference() => $_clearField(151211160);

  @$pb.TagNumber(358915841)
  CloudWatchAlarmConfiguration get cloudwatchalarmconfiguration => $_getN(3);
  @$pb.TagNumber(358915841)
  set cloudwatchalarmconfiguration(CloudWatchAlarmConfiguration value) =>
      $_setField(358915841, value);
  @$pb.TagNumber(358915841)
  $core.bool hasCloudwatchalarmconfiguration() => $_has(3);
  @$pb.TagNumber(358915841)
  void clearCloudwatchalarmconfiguration() => $_clearField(358915841);
  @$pb.TagNumber(358915841)
  CloudWatchAlarmConfiguration ensureCloudwatchalarmconfiguration() =>
      $_ensure(3);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(438614164)
  LinkedService get linkedservice => $_getN(5);
  @$pb.TagNumber(438614164)
  set linkedservice(LinkedService value) => $_setField(438614164, value);
  @$pb.TagNumber(438614164)
  $core.bool hasLinkedservice() => $_has(5);
  @$pb.TagNumber(438614164)
  void clearLinkedservice() => $_clearField(438614164);
  @$pb.TagNumber(438614164)
  LinkedService ensureLinkedservice() => $_ensure(5);
}

class HealthCheckAlreadyExists extends $pb.GeneratedMessage {
  factory HealthCheckAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  HealthCheckAlreadyExists._();

  factory HealthCheckAlreadyExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealthCheckAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealthCheckAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckAlreadyExists copyWith(
          void Function(HealthCheckAlreadyExists) updates) =>
      super.copyWith((message) => updates(message as HealthCheckAlreadyExists))
          as HealthCheckAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealthCheckAlreadyExists create() => HealthCheckAlreadyExists._();
  @$core.override
  HealthCheckAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealthCheckAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealthCheckAlreadyExists>(create);
  static HealthCheckAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class HealthCheckConfig extends $pb.GeneratedMessage {
  factory HealthCheckConfig({
    $core.Iterable<HealthCheckRegion>? regions,
    $core.int? port,
    $core.bool? inverted,
    $core.bool? enablesni,
    $core.bool? measurelatency,
    $core.String? resourcepath,
    $core.String? ipaddress,
    $core.int? failurethreshold,
    $core.String? routingcontrolarn,
    $core.int? healththreshold,
    HealthCheckType? type,
    $core.String? searchstring,
    $core.int? requestinterval,
    $core.String? fullyqualifieddomainname,
    $core.Iterable<$core.String>? childhealthchecks,
    InsufficientDataHealthStatus? insufficientdatahealthstatus,
    $core.bool? disabled,
    AlarmIdentifier? alarmidentifier,
  }) {
    final result = create();
    if (regions != null) result.regions.addAll(regions);
    if (port != null) result.port = port;
    if (inverted != null) result.inverted = inverted;
    if (enablesni != null) result.enablesni = enablesni;
    if (measurelatency != null) result.measurelatency = measurelatency;
    if (resourcepath != null) result.resourcepath = resourcepath;
    if (ipaddress != null) result.ipaddress = ipaddress;
    if (failurethreshold != null) result.failurethreshold = failurethreshold;
    if (routingcontrolarn != null) result.routingcontrolarn = routingcontrolarn;
    if (healththreshold != null) result.healththreshold = healththreshold;
    if (type != null) result.type = type;
    if (searchstring != null) result.searchstring = searchstring;
    if (requestinterval != null) result.requestinterval = requestinterval;
    if (fullyqualifieddomainname != null)
      result.fullyqualifieddomainname = fullyqualifieddomainname;
    if (childhealthchecks != null)
      result.childhealthchecks.addAll(childhealthchecks);
    if (insufficientdatahealthstatus != null)
      result.insufficientdatahealthstatus = insufficientdatahealthstatus;
    if (disabled != null) result.disabled = disabled;
    if (alarmidentifier != null) result.alarmidentifier = alarmidentifier;
    return result;
  }

  HealthCheckConfig._();

  factory HealthCheckConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealthCheckConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealthCheckConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pc<HealthCheckRegion>(
        36200107, _omitFieldNames ? '' : 'regions', $pb.PbFieldType.KE,
        valueOf: HealthCheckRegion.valueOf,
        enumValues: HealthCheckRegion.values,
        defaultEnumValue: HealthCheckRegion.HEALTH_CHECK_REGION_AP_NORTHEAST_1)
    ..aI(46480583, _omitFieldNames ? '' : 'port')
    ..aOB(55175513, _omitFieldNames ? '' : 'inverted')
    ..aOB(70122887, _omitFieldNames ? '' : 'enablesni')
    ..aOB(87136848, _omitFieldNames ? '' : 'measurelatency')
    ..aOS(117584551, _omitFieldNames ? '' : 'resourcepath')
    ..aOS(169333741, _omitFieldNames ? '' : 'ipaddress')
    ..aI(176846565, _omitFieldNames ? '' : 'failurethreshold')
    ..aOS(206883790, _omitFieldNames ? '' : 'routingcontrolarn')
    ..aI(215873163, _omitFieldNames ? '' : 'healththreshold')
    ..aE<HealthCheckType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: HealthCheckType.values)
    ..aOS(318687365, _omitFieldNames ? '' : 'searchstring')
    ..aI(350673112, _omitFieldNames ? '' : 'requestinterval')
    ..aOS(459321509, _omitFieldNames ? '' : 'fullyqualifieddomainname')
    ..pPS(485535935, _omitFieldNames ? '' : 'childhealthchecks')
    ..aE<InsufficientDataHealthStatus>(
        493115723, _omitFieldNames ? '' : 'insufficientdatahealthstatus',
        enumValues: InsufficientDataHealthStatus.values)
    ..aOB(533633318, _omitFieldNames ? '' : 'disabled')
    ..aOM<AlarmIdentifier>(536124346, _omitFieldNames ? '' : 'alarmidentifier',
        subBuilder: AlarmIdentifier.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckConfig copyWith(void Function(HealthCheckConfig) updates) =>
      super.copyWith((message) => updates(message as HealthCheckConfig))
          as HealthCheckConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealthCheckConfig create() => HealthCheckConfig._();
  @$core.override
  HealthCheckConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealthCheckConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealthCheckConfig>(create);
  static HealthCheckConfig? _defaultInstance;

  @$pb.TagNumber(36200107)
  $pb.PbList<HealthCheckRegion> get regions => $_getList(0);

  @$pb.TagNumber(46480583)
  $core.int get port => $_getIZ(1);
  @$pb.TagNumber(46480583)
  set port($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(46480583)
  $core.bool hasPort() => $_has(1);
  @$pb.TagNumber(46480583)
  void clearPort() => $_clearField(46480583);

  @$pb.TagNumber(55175513)
  $core.bool get inverted => $_getBF(2);
  @$pb.TagNumber(55175513)
  set inverted($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(55175513)
  $core.bool hasInverted() => $_has(2);
  @$pb.TagNumber(55175513)
  void clearInverted() => $_clearField(55175513);

  @$pb.TagNumber(70122887)
  $core.bool get enablesni => $_getBF(3);
  @$pb.TagNumber(70122887)
  set enablesni($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(70122887)
  $core.bool hasEnablesni() => $_has(3);
  @$pb.TagNumber(70122887)
  void clearEnablesni() => $_clearField(70122887);

  @$pb.TagNumber(87136848)
  $core.bool get measurelatency => $_getBF(4);
  @$pb.TagNumber(87136848)
  set measurelatency($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(87136848)
  $core.bool hasMeasurelatency() => $_has(4);
  @$pb.TagNumber(87136848)
  void clearMeasurelatency() => $_clearField(87136848);

  @$pb.TagNumber(117584551)
  $core.String get resourcepath => $_getSZ(5);
  @$pb.TagNumber(117584551)
  set resourcepath($core.String value) => $_setString(5, value);
  @$pb.TagNumber(117584551)
  $core.bool hasResourcepath() => $_has(5);
  @$pb.TagNumber(117584551)
  void clearResourcepath() => $_clearField(117584551);

  @$pb.TagNumber(169333741)
  $core.String get ipaddress => $_getSZ(6);
  @$pb.TagNumber(169333741)
  set ipaddress($core.String value) => $_setString(6, value);
  @$pb.TagNumber(169333741)
  $core.bool hasIpaddress() => $_has(6);
  @$pb.TagNumber(169333741)
  void clearIpaddress() => $_clearField(169333741);

  @$pb.TagNumber(176846565)
  $core.int get failurethreshold => $_getIZ(7);
  @$pb.TagNumber(176846565)
  set failurethreshold($core.int value) => $_setSignedInt32(7, value);
  @$pb.TagNumber(176846565)
  $core.bool hasFailurethreshold() => $_has(7);
  @$pb.TagNumber(176846565)
  void clearFailurethreshold() => $_clearField(176846565);

  @$pb.TagNumber(206883790)
  $core.String get routingcontrolarn => $_getSZ(8);
  @$pb.TagNumber(206883790)
  set routingcontrolarn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(206883790)
  $core.bool hasRoutingcontrolarn() => $_has(8);
  @$pb.TagNumber(206883790)
  void clearRoutingcontrolarn() => $_clearField(206883790);

  @$pb.TagNumber(215873163)
  $core.int get healththreshold => $_getIZ(9);
  @$pb.TagNumber(215873163)
  set healththreshold($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(215873163)
  $core.bool hasHealththreshold() => $_has(9);
  @$pb.TagNumber(215873163)
  void clearHealththreshold() => $_clearField(215873163);

  @$pb.TagNumber(290836590)
  HealthCheckType get type => $_getN(10);
  @$pb.TagNumber(290836590)
  set type(HealthCheckType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(10);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(318687365)
  $core.String get searchstring => $_getSZ(11);
  @$pb.TagNumber(318687365)
  set searchstring($core.String value) => $_setString(11, value);
  @$pb.TagNumber(318687365)
  $core.bool hasSearchstring() => $_has(11);
  @$pb.TagNumber(318687365)
  void clearSearchstring() => $_clearField(318687365);

  @$pb.TagNumber(350673112)
  $core.int get requestinterval => $_getIZ(12);
  @$pb.TagNumber(350673112)
  set requestinterval($core.int value) => $_setSignedInt32(12, value);
  @$pb.TagNumber(350673112)
  $core.bool hasRequestinterval() => $_has(12);
  @$pb.TagNumber(350673112)
  void clearRequestinterval() => $_clearField(350673112);

  @$pb.TagNumber(459321509)
  $core.String get fullyqualifieddomainname => $_getSZ(13);
  @$pb.TagNumber(459321509)
  set fullyqualifieddomainname($core.String value) => $_setString(13, value);
  @$pb.TagNumber(459321509)
  $core.bool hasFullyqualifieddomainname() => $_has(13);
  @$pb.TagNumber(459321509)
  void clearFullyqualifieddomainname() => $_clearField(459321509);

  @$pb.TagNumber(485535935)
  $pb.PbList<$core.String> get childhealthchecks => $_getList(14);

  @$pb.TagNumber(493115723)
  InsufficientDataHealthStatus get insufficientdatahealthstatus => $_getN(15);
  @$pb.TagNumber(493115723)
  set insufficientdatahealthstatus(InsufficientDataHealthStatus value) =>
      $_setField(493115723, value);
  @$pb.TagNumber(493115723)
  $core.bool hasInsufficientdatahealthstatus() => $_has(15);
  @$pb.TagNumber(493115723)
  void clearInsufficientdatahealthstatus() => $_clearField(493115723);

  @$pb.TagNumber(533633318)
  $core.bool get disabled => $_getBF(16);
  @$pb.TagNumber(533633318)
  set disabled($core.bool value) => $_setBool(16, value);
  @$pb.TagNumber(533633318)
  $core.bool hasDisabled() => $_has(16);
  @$pb.TagNumber(533633318)
  void clearDisabled() => $_clearField(533633318);

  @$pb.TagNumber(536124346)
  AlarmIdentifier get alarmidentifier => $_getN(17);
  @$pb.TagNumber(536124346)
  set alarmidentifier(AlarmIdentifier value) => $_setField(536124346, value);
  @$pb.TagNumber(536124346)
  $core.bool hasAlarmidentifier() => $_has(17);
  @$pb.TagNumber(536124346)
  void clearAlarmidentifier() => $_clearField(536124346);
  @$pb.TagNumber(536124346)
  AlarmIdentifier ensureAlarmidentifier() => $_ensure(17);
}

class HealthCheckInUse extends $pb.GeneratedMessage {
  factory HealthCheckInUse({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  HealthCheckInUse._();

  factory HealthCheckInUse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealthCheckInUse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealthCheckInUse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckInUse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckInUse copyWith(void Function(HealthCheckInUse) updates) =>
      super.copyWith((message) => updates(message as HealthCheckInUse))
          as HealthCheckInUse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealthCheckInUse create() => HealthCheckInUse._();
  @$core.override
  HealthCheckInUse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealthCheckInUse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealthCheckInUse>(create);
  static HealthCheckInUse? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class HealthCheckObservation extends $pb.GeneratedMessage {
  factory HealthCheckObservation({
    StatusReport? statusreport,
    HealthCheckRegion? region,
    $core.String? ipaddress,
  }) {
    final result = create();
    if (statusreport != null) result.statusreport = statusreport;
    if (region != null) result.region = region;
    if (ipaddress != null) result.ipaddress = ipaddress;
    return result;
  }

  HealthCheckObservation._();

  factory HealthCheckObservation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealthCheckObservation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealthCheckObservation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<StatusReport>(27958834, _omitFieldNames ? '' : 'statusreport',
        subBuilder: StatusReport.create)
    ..aE<HealthCheckRegion>(154040478, _omitFieldNames ? '' : 'region',
        enumValues: HealthCheckRegion.values)
    ..aOS(169333741, _omitFieldNames ? '' : 'ipaddress')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckObservation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckObservation copyWith(
          void Function(HealthCheckObservation) updates) =>
      super.copyWith((message) => updates(message as HealthCheckObservation))
          as HealthCheckObservation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealthCheckObservation create() => HealthCheckObservation._();
  @$core.override
  HealthCheckObservation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealthCheckObservation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealthCheckObservation>(create);
  static HealthCheckObservation? _defaultInstance;

  @$pb.TagNumber(27958834)
  StatusReport get statusreport => $_getN(0);
  @$pb.TagNumber(27958834)
  set statusreport(StatusReport value) => $_setField(27958834, value);
  @$pb.TagNumber(27958834)
  $core.bool hasStatusreport() => $_has(0);
  @$pb.TagNumber(27958834)
  void clearStatusreport() => $_clearField(27958834);
  @$pb.TagNumber(27958834)
  StatusReport ensureStatusreport() => $_ensure(0);

  @$pb.TagNumber(154040478)
  HealthCheckRegion get region => $_getN(1);
  @$pb.TagNumber(154040478)
  set region(HealthCheckRegion value) => $_setField(154040478, value);
  @$pb.TagNumber(154040478)
  $core.bool hasRegion() => $_has(1);
  @$pb.TagNumber(154040478)
  void clearRegion() => $_clearField(154040478);

  @$pb.TagNumber(169333741)
  $core.String get ipaddress => $_getSZ(2);
  @$pb.TagNumber(169333741)
  set ipaddress($core.String value) => $_setString(2, value);
  @$pb.TagNumber(169333741)
  $core.bool hasIpaddress() => $_has(2);
  @$pb.TagNumber(169333741)
  void clearIpaddress() => $_clearField(169333741);
}

class HealthCheckVersionMismatch extends $pb.GeneratedMessage {
  factory HealthCheckVersionMismatch({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  HealthCheckVersionMismatch._();

  factory HealthCheckVersionMismatch.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealthCheckVersionMismatch.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealthCheckVersionMismatch',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckVersionMismatch clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthCheckVersionMismatch copyWith(
          void Function(HealthCheckVersionMismatch) updates) =>
      super.copyWith(
              (message) => updates(message as HealthCheckVersionMismatch))
          as HealthCheckVersionMismatch;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealthCheckVersionMismatch create() => HealthCheckVersionMismatch._();
  @$core.override
  HealthCheckVersionMismatch createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealthCheckVersionMismatch getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealthCheckVersionMismatch>(create);
  static HealthCheckVersionMismatch? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class HostedZone extends $pb.GeneratedMessage {
  factory HostedZone({
    $core.String? callerreference,
    HostedZoneConfig? config,
    HostedZoneFeatures? features,
    $core.String? name,
    $fixnum.Int64? resourcerecordsetcount,
    $core.String? id,
    LinkedService? linkedservice,
  }) {
    final result = create();
    if (callerreference != null) result.callerreference = callerreference;
    if (config != null) result.config = config;
    if (features != null) result.features = features;
    if (name != null) result.name = name;
    if (resourcerecordsetcount != null)
      result.resourcerecordsetcount = resourcerecordsetcount;
    if (id != null) result.id = id;
    if (linkedservice != null) result.linkedservice = linkedservice;
    return result;
  }

  HostedZone._();

  factory HostedZone.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZone.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZone',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(151211160, _omitFieldNames ? '' : 'callerreference')
    ..aOM<HostedZoneConfig>(169009384, _omitFieldNames ? '' : 'config',
        subBuilder: HostedZoneConfig.create)
    ..aOM<HostedZoneFeatures>(205993963, _omitFieldNames ? '' : 'features',
        subBuilder: HostedZoneFeatures.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aInt64(334391116, _omitFieldNames ? '' : 'resourcerecordsetcount')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOM<LinkedService>(438614164, _omitFieldNames ? '' : 'linkedservice',
        subBuilder: LinkedService.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZone clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZone copyWith(void Function(HostedZone) updates) =>
      super.copyWith((message) => updates(message as HostedZone)) as HostedZone;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZone create() => HostedZone._();
  @$core.override
  HostedZone createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZone getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZone>(create);
  static HostedZone? _defaultInstance;

  @$pb.TagNumber(151211160)
  $core.String get callerreference => $_getSZ(0);
  @$pb.TagNumber(151211160)
  set callerreference($core.String value) => $_setString(0, value);
  @$pb.TagNumber(151211160)
  $core.bool hasCallerreference() => $_has(0);
  @$pb.TagNumber(151211160)
  void clearCallerreference() => $_clearField(151211160);

  @$pb.TagNumber(169009384)
  HostedZoneConfig get config => $_getN(1);
  @$pb.TagNumber(169009384)
  set config(HostedZoneConfig value) => $_setField(169009384, value);
  @$pb.TagNumber(169009384)
  $core.bool hasConfig() => $_has(1);
  @$pb.TagNumber(169009384)
  void clearConfig() => $_clearField(169009384);
  @$pb.TagNumber(169009384)
  HostedZoneConfig ensureConfig() => $_ensure(1);

  @$pb.TagNumber(205993963)
  HostedZoneFeatures get features => $_getN(2);
  @$pb.TagNumber(205993963)
  set features(HostedZoneFeatures value) => $_setField(205993963, value);
  @$pb.TagNumber(205993963)
  $core.bool hasFeatures() => $_has(2);
  @$pb.TagNumber(205993963)
  void clearFeatures() => $_clearField(205993963);
  @$pb.TagNumber(205993963)
  HostedZoneFeatures ensureFeatures() => $_ensure(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(334391116)
  $fixnum.Int64 get resourcerecordsetcount => $_getI64(4);
  @$pb.TagNumber(334391116)
  set resourcerecordsetcount($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(334391116)
  $core.bool hasResourcerecordsetcount() => $_has(4);
  @$pb.TagNumber(334391116)
  void clearResourcerecordsetcount() => $_clearField(334391116);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(5);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(5, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(5);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(438614164)
  LinkedService get linkedservice => $_getN(6);
  @$pb.TagNumber(438614164)
  set linkedservice(LinkedService value) => $_setField(438614164, value);
  @$pb.TagNumber(438614164)
  $core.bool hasLinkedservice() => $_has(6);
  @$pb.TagNumber(438614164)
  void clearLinkedservice() => $_clearField(438614164);
  @$pb.TagNumber(438614164)
  LinkedService ensureLinkedservice() => $_ensure(6);
}

class HostedZoneAlreadyExists extends $pb.GeneratedMessage {
  factory HostedZoneAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  HostedZoneAlreadyExists._();

  factory HostedZoneAlreadyExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneAlreadyExists copyWith(
          void Function(HostedZoneAlreadyExists) updates) =>
      super.copyWith((message) => updates(message as HostedZoneAlreadyExists))
          as HostedZoneAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneAlreadyExists create() => HostedZoneAlreadyExists._();
  @$core.override
  HostedZoneAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneAlreadyExists>(create);
  static HostedZoneAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class HostedZoneConfig extends $pb.GeneratedMessage {
  factory HostedZoneConfig({
    $core.String? comment,
    $core.bool? privatezone,
  }) {
    final result = create();
    if (comment != null) result.comment = comment;
    if (privatezone != null) result.privatezone = privatezone;
    return result;
  }

  HostedZoneConfig._();

  factory HostedZoneConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..aOB(506703867, _omitFieldNames ? '' : 'privatezone')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneConfig copyWith(void Function(HostedZoneConfig) updates) =>
      super.copyWith((message) => updates(message as HostedZoneConfig))
          as HostedZoneConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneConfig create() => HostedZoneConfig._();
  @$core.override
  HostedZoneConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneConfig>(create);
  static HostedZoneConfig? _defaultInstance;

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(0);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(0, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(0);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);

  @$pb.TagNumber(506703867)
  $core.bool get privatezone => $_getBF(1);
  @$pb.TagNumber(506703867)
  set privatezone($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(506703867)
  $core.bool hasPrivatezone() => $_has(1);
  @$pb.TagNumber(506703867)
  void clearPrivatezone() => $_clearField(506703867);
}

class HostedZoneFailureReasons extends $pb.GeneratedMessage {
  factory HostedZoneFailureReasons({
    $core.String? acceleratedrecovery,
  }) {
    final result = create();
    if (acceleratedrecovery != null)
      result.acceleratedrecovery = acceleratedrecovery;
    return result;
  }

  HostedZoneFailureReasons._();

  factory HostedZoneFailureReasons.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneFailureReasons.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneFailureReasons',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(356087092, _omitFieldNames ? '' : 'acceleratedrecovery')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneFailureReasons clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneFailureReasons copyWith(
          void Function(HostedZoneFailureReasons) updates) =>
      super.copyWith((message) => updates(message as HostedZoneFailureReasons))
          as HostedZoneFailureReasons;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneFailureReasons create() => HostedZoneFailureReasons._();
  @$core.override
  HostedZoneFailureReasons createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneFailureReasons getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneFailureReasons>(create);
  static HostedZoneFailureReasons? _defaultInstance;

  @$pb.TagNumber(356087092)
  $core.String get acceleratedrecovery => $_getSZ(0);
  @$pb.TagNumber(356087092)
  set acceleratedrecovery($core.String value) => $_setString(0, value);
  @$pb.TagNumber(356087092)
  $core.bool hasAcceleratedrecovery() => $_has(0);
  @$pb.TagNumber(356087092)
  void clearAcceleratedrecovery() => $_clearField(356087092);
}

class HostedZoneFeatures extends $pb.GeneratedMessage {
  factory HostedZoneFeatures({
    HostedZoneFailureReasons? failurereasons,
    AcceleratedRecoveryStatus? acceleratedrecoverystatus,
  }) {
    final result = create();
    if (failurereasons != null) result.failurereasons = failurereasons;
    if (acceleratedrecoverystatus != null)
      result.acceleratedrecoverystatus = acceleratedrecoverystatus;
    return result;
  }

  HostedZoneFeatures._();

  factory HostedZoneFeatures.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneFeatures.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneFeatures',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<HostedZoneFailureReasons>(
        445146219, _omitFieldNames ? '' : 'failurereasons',
        subBuilder: HostedZoneFailureReasons.create)
    ..aE<AcceleratedRecoveryStatus>(
        507534246, _omitFieldNames ? '' : 'acceleratedrecoverystatus',
        enumValues: AcceleratedRecoveryStatus.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneFeatures clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneFeatures copyWith(void Function(HostedZoneFeatures) updates) =>
      super.copyWith((message) => updates(message as HostedZoneFeatures))
          as HostedZoneFeatures;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneFeatures create() => HostedZoneFeatures._();
  @$core.override
  HostedZoneFeatures createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneFeatures getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneFeatures>(create);
  static HostedZoneFeatures? _defaultInstance;

  @$pb.TagNumber(445146219)
  HostedZoneFailureReasons get failurereasons => $_getN(0);
  @$pb.TagNumber(445146219)
  set failurereasons(HostedZoneFailureReasons value) =>
      $_setField(445146219, value);
  @$pb.TagNumber(445146219)
  $core.bool hasFailurereasons() => $_has(0);
  @$pb.TagNumber(445146219)
  void clearFailurereasons() => $_clearField(445146219);
  @$pb.TagNumber(445146219)
  HostedZoneFailureReasons ensureFailurereasons() => $_ensure(0);

  @$pb.TagNumber(507534246)
  AcceleratedRecoveryStatus get acceleratedrecoverystatus => $_getN(1);
  @$pb.TagNumber(507534246)
  set acceleratedrecoverystatus(AcceleratedRecoveryStatus value) =>
      $_setField(507534246, value);
  @$pb.TagNumber(507534246)
  $core.bool hasAcceleratedrecoverystatus() => $_has(1);
  @$pb.TagNumber(507534246)
  void clearAcceleratedrecoverystatus() => $_clearField(507534246);
}

class HostedZoneLimit extends $pb.GeneratedMessage {
  factory HostedZoneLimit({
    $fixnum.Int64? value,
    HostedZoneLimitType? type,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (type != null) result.type = type;
    return result;
  }

  HostedZoneLimit._();

  factory HostedZoneLimit.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneLimit.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneLimit',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(289929579, _omitFieldNames ? '' : 'value')
    ..aE<HostedZoneLimitType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: HostedZoneLimitType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneLimit clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneLimit copyWith(void Function(HostedZoneLimit) updates) =>
      super.copyWith((message) => updates(message as HostedZoneLimit))
          as HostedZoneLimit;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneLimit create() => HostedZoneLimit._();
  @$core.override
  HostedZoneLimit createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneLimit getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneLimit>(create);
  static HostedZoneLimit? _defaultInstance;

  @$pb.TagNumber(289929579)
  $fixnum.Int64 get value => $_getI64(0);
  @$pb.TagNumber(289929579)
  set value($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(290836590)
  HostedZoneLimitType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(HostedZoneLimitType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class HostedZoneNotEmpty extends $pb.GeneratedMessage {
  factory HostedZoneNotEmpty({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  HostedZoneNotEmpty._();

  factory HostedZoneNotEmpty.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneNotEmpty.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneNotEmpty',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneNotEmpty clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneNotEmpty copyWith(void Function(HostedZoneNotEmpty) updates) =>
      super.copyWith((message) => updates(message as HostedZoneNotEmpty))
          as HostedZoneNotEmpty;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneNotEmpty create() => HostedZoneNotEmpty._();
  @$core.override
  HostedZoneNotEmpty createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneNotEmpty getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneNotEmpty>(create);
  static HostedZoneNotEmpty? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class HostedZoneNotFound extends $pb.GeneratedMessage {
  factory HostedZoneNotFound({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  HostedZoneNotFound._();

  factory HostedZoneNotFound.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneNotFound.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneNotFound',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneNotFound clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneNotFound copyWith(void Function(HostedZoneNotFound) updates) =>
      super.copyWith((message) => updates(message as HostedZoneNotFound))
          as HostedZoneNotFound;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneNotFound create() => HostedZoneNotFound._();
  @$core.override
  HostedZoneNotFound createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneNotFound getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneNotFound>(create);
  static HostedZoneNotFound? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class HostedZoneNotPrivate extends $pb.GeneratedMessage {
  factory HostedZoneNotPrivate({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  HostedZoneNotPrivate._();

  factory HostedZoneNotPrivate.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneNotPrivate.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneNotPrivate',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneNotPrivate clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneNotPrivate copyWith(void Function(HostedZoneNotPrivate) updates) =>
      super.copyWith((message) => updates(message as HostedZoneNotPrivate))
          as HostedZoneNotPrivate;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneNotPrivate create() => HostedZoneNotPrivate._();
  @$core.override
  HostedZoneNotPrivate createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneNotPrivate getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneNotPrivate>(create);
  static HostedZoneNotPrivate? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class HostedZoneOwner extends $pb.GeneratedMessage {
  factory HostedZoneOwner({
    $core.String? owningaccount,
    $core.String? owningservice,
  }) {
    final result = create();
    if (owningaccount != null) result.owningaccount = owningaccount;
    if (owningservice != null) result.owningservice = owningservice;
    return result;
  }

  HostedZoneOwner._();

  factory HostedZoneOwner.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneOwner.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneOwner',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(339968073, _omitFieldNames ? '' : 'owningaccount')
    ..aOS(462487817, _omitFieldNames ? '' : 'owningservice')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneOwner clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneOwner copyWith(void Function(HostedZoneOwner) updates) =>
      super.copyWith((message) => updates(message as HostedZoneOwner))
          as HostedZoneOwner;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneOwner create() => HostedZoneOwner._();
  @$core.override
  HostedZoneOwner createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneOwner getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneOwner>(create);
  static HostedZoneOwner? _defaultInstance;

  @$pb.TagNumber(339968073)
  $core.String get owningaccount => $_getSZ(0);
  @$pb.TagNumber(339968073)
  set owningaccount($core.String value) => $_setString(0, value);
  @$pb.TagNumber(339968073)
  $core.bool hasOwningaccount() => $_has(0);
  @$pb.TagNumber(339968073)
  void clearOwningaccount() => $_clearField(339968073);

  @$pb.TagNumber(462487817)
  $core.String get owningservice => $_getSZ(1);
  @$pb.TagNumber(462487817)
  set owningservice($core.String value) => $_setString(1, value);
  @$pb.TagNumber(462487817)
  $core.bool hasOwningservice() => $_has(1);
  @$pb.TagNumber(462487817)
  void clearOwningservice() => $_clearField(462487817);
}

class HostedZonePartiallyDelegated extends $pb.GeneratedMessage {
  factory HostedZonePartiallyDelegated({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  HostedZonePartiallyDelegated._();

  factory HostedZonePartiallyDelegated.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZonePartiallyDelegated.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZonePartiallyDelegated',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZonePartiallyDelegated clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZonePartiallyDelegated copyWith(
          void Function(HostedZonePartiallyDelegated) updates) =>
      super.copyWith(
              (message) => updates(message as HostedZonePartiallyDelegated))
          as HostedZonePartiallyDelegated;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZonePartiallyDelegated create() =>
      HostedZonePartiallyDelegated._();
  @$core.override
  HostedZonePartiallyDelegated createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZonePartiallyDelegated getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZonePartiallyDelegated>(create);
  static HostedZonePartiallyDelegated? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class HostedZoneSummary extends $pb.GeneratedMessage {
  factory HostedZoneSummary({
    $core.String? name,
    $core.String? hostedzoneid,
    HostedZoneOwner? owner,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (owner != null) result.owner = owner;
    return result;
  }

  HostedZoneSummary._();

  factory HostedZoneSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HostedZoneSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HostedZoneSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOM<HostedZoneOwner>(455261813, _omitFieldNames ? '' : 'owner',
        subBuilder: HostedZoneOwner.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HostedZoneSummary copyWith(void Function(HostedZoneSummary) updates) =>
      super.copyWith((message) => updates(message as HostedZoneSummary))
          as HostedZoneSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HostedZoneSummary create() => HostedZoneSummary._();
  @$core.override
  HostedZoneSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HostedZoneSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HostedZoneSummary>(create);
  static HostedZoneSummary? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(455261813)
  HostedZoneOwner get owner => $_getN(2);
  @$pb.TagNumber(455261813)
  set owner(HostedZoneOwner value) => $_setField(455261813, value);
  @$pb.TagNumber(455261813)
  $core.bool hasOwner() => $_has(2);
  @$pb.TagNumber(455261813)
  void clearOwner() => $_clearField(455261813);
  @$pb.TagNumber(455261813)
  HostedZoneOwner ensureOwner() => $_ensure(2);
}

class IncompatibleVersion extends $pb.GeneratedMessage {
  factory IncompatibleVersion({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IncompatibleVersion._();

  factory IncompatibleVersion.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IncompatibleVersion.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IncompatibleVersion',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncompatibleVersion clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IncompatibleVersion copyWith(void Function(IncompatibleVersion) updates) =>
      super.copyWith((message) => updates(message as IncompatibleVersion))
          as IncompatibleVersion;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IncompatibleVersion create() => IncompatibleVersion._();
  @$core.override
  IncompatibleVersion createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IncompatibleVersion getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IncompatibleVersion>(create);
  static IncompatibleVersion? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InsufficientCloudWatchLogsResourcePolicy extends $pb.GeneratedMessage {
  factory InsufficientCloudWatchLogsResourcePolicy({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InsufficientCloudWatchLogsResourcePolicy._();

  factory InsufficientCloudWatchLogsResourcePolicy.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsufficientCloudWatchLogsResourcePolicy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsufficientCloudWatchLogsResourcePolicy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientCloudWatchLogsResourcePolicy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsufficientCloudWatchLogsResourcePolicy copyWith(
          void Function(InsufficientCloudWatchLogsResourcePolicy) updates) =>
      super.copyWith((message) =>
              updates(message as InsufficientCloudWatchLogsResourcePolicy))
          as InsufficientCloudWatchLogsResourcePolicy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsufficientCloudWatchLogsResourcePolicy create() =>
      InsufficientCloudWatchLogsResourcePolicy._();
  @$core.override
  InsufficientCloudWatchLogsResourcePolicy createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsufficientCloudWatchLogsResourcePolicy getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InsufficientCloudWatchLogsResourcePolicy>(create);
  static InsufficientCloudWatchLogsResourcePolicy? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidArgument extends $pb.GeneratedMessage {
  factory InvalidArgument({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidArgument._();

  factory InvalidArgument.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidArgument.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidArgument',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArgument clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidArgument copyWith(void Function(InvalidArgument) updates) =>
      super.copyWith((message) => updates(message as InvalidArgument))
          as InvalidArgument;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidArgument create() => InvalidArgument._();
  @$core.override
  InvalidArgument createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidArgument getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidArgument>(create);
  static InvalidArgument? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidChangeBatch extends $pb.GeneratedMessage {
  factory InvalidChangeBatch({
    $core.Iterable<$core.String>? messages,
    $core.String? message,
  }) {
    final result = create();
    if (messages != null) result.messages.addAll(messages);
    if (message != null) result.message = message;
    return result;
  }

  InvalidChangeBatch._();

  factory InvalidChangeBatch.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidChangeBatch.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidChangeBatch',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPS(230838, _omitFieldNames ? '' : 'messages')
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidChangeBatch clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidChangeBatch copyWith(void Function(InvalidChangeBatch) updates) =>
      super.copyWith((message) => updates(message as InvalidChangeBatch))
          as InvalidChangeBatch;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidChangeBatch create() => InvalidChangeBatch._();
  @$core.override
  InvalidChangeBatch createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidChangeBatch getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidChangeBatch>(create);
  static InvalidChangeBatch? _defaultInstance;

  @$pb.TagNumber(230838)
  $pb.PbList<$core.String> get messages => $_getList(0);

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidDomainName extends $pb.GeneratedMessage {
  factory InvalidDomainName({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidDomainName._();

  factory InvalidDomainName.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidDomainName.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidDomainName',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidDomainName clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidDomainName copyWith(void Function(InvalidDomainName) updates) =>
      super.copyWith((message) => updates(message as InvalidDomainName))
          as InvalidDomainName;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidDomainName create() => InvalidDomainName._();
  @$core.override
  InvalidDomainName createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidDomainName getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidDomainName>(create);
  static InvalidDomainName? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidInput extends $pb.GeneratedMessage {
  factory InvalidInput({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidInput._();

  factory InvalidInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidInput copyWith(void Function(InvalidInput) updates) =>
      super.copyWith((message) => updates(message as InvalidInput))
          as InvalidInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidInput create() => InvalidInput._();
  @$core.override
  InvalidInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidInput>(create);
  static InvalidInput? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidKMSArn extends $pb.GeneratedMessage {
  factory InvalidKMSArn({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidKMSArn._();

  factory InvalidKMSArn.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidKMSArn.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidKMSArn',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKMSArn clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKMSArn copyWith(void Function(InvalidKMSArn) updates) =>
      super.copyWith((message) => updates(message as InvalidKMSArn))
          as InvalidKMSArn;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidKMSArn create() => InvalidKMSArn._();
  @$core.override
  InvalidKMSArn createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidKMSArn getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidKMSArn>(create);
  static InvalidKMSArn? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidKeySigningKeyName extends $pb.GeneratedMessage {
  factory InvalidKeySigningKeyName({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidKeySigningKeyName._();

  factory InvalidKeySigningKeyName.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidKeySigningKeyName.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidKeySigningKeyName',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKeySigningKeyName clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKeySigningKeyName copyWith(
          void Function(InvalidKeySigningKeyName) updates) =>
      super.copyWith((message) => updates(message as InvalidKeySigningKeyName))
          as InvalidKeySigningKeyName;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidKeySigningKeyName create() => InvalidKeySigningKeyName._();
  @$core.override
  InvalidKeySigningKeyName createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidKeySigningKeyName getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidKeySigningKeyName>(create);
  static InvalidKeySigningKeyName? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidKeySigningKeyStatus extends $pb.GeneratedMessage {
  factory InvalidKeySigningKeyStatus({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidKeySigningKeyStatus._();

  factory InvalidKeySigningKeyStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidKeySigningKeyStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidKeySigningKeyStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKeySigningKeyStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidKeySigningKeyStatus copyWith(
          void Function(InvalidKeySigningKeyStatus) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidKeySigningKeyStatus))
          as InvalidKeySigningKeyStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidKeySigningKeyStatus create() => InvalidKeySigningKeyStatus._();
  @$core.override
  InvalidKeySigningKeyStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidKeySigningKeyStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidKeySigningKeyStatus>(create);
  static InvalidKeySigningKeyStatus? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidPaginationToken extends $pb.GeneratedMessage {
  factory InvalidPaginationToken({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidPaginationToken._();

  factory InvalidPaginationToken.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidPaginationToken.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidPaginationToken',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidPaginationToken clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidPaginationToken copyWith(
          void Function(InvalidPaginationToken) updates) =>
      super.copyWith((message) => updates(message as InvalidPaginationToken))
          as InvalidPaginationToken;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidPaginationToken create() => InvalidPaginationToken._();
  @$core.override
  InvalidPaginationToken createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidPaginationToken getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidPaginationToken>(create);
  static InvalidPaginationToken? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidSigningStatus extends $pb.GeneratedMessage {
  factory InvalidSigningStatus({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidSigningStatus._();

  factory InvalidSigningStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidSigningStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidSigningStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSigningStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidSigningStatus copyWith(void Function(InvalidSigningStatus) updates) =>
      super.copyWith((message) => updates(message as InvalidSigningStatus))
          as InvalidSigningStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidSigningStatus create() => InvalidSigningStatus._();
  @$core.override
  InvalidSigningStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidSigningStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidSigningStatus>(create);
  static InvalidSigningStatus? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidTrafficPolicyDocument extends $pb.GeneratedMessage {
  factory InvalidTrafficPolicyDocument({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidTrafficPolicyDocument._();

  factory InvalidTrafficPolicyDocument.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidTrafficPolicyDocument.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidTrafficPolicyDocument',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTrafficPolicyDocument clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidTrafficPolicyDocument copyWith(
          void Function(InvalidTrafficPolicyDocument) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidTrafficPolicyDocument))
          as InvalidTrafficPolicyDocument;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidTrafficPolicyDocument create() =>
      InvalidTrafficPolicyDocument._();
  @$core.override
  InvalidTrafficPolicyDocument createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidTrafficPolicyDocument getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidTrafficPolicyDocument>(create);
  static InvalidTrafficPolicyDocument? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidVPCId extends $pb.GeneratedMessage {
  factory InvalidVPCId({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidVPCId._();

  factory InvalidVPCId.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidVPCId.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidVPCId',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidVPCId clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidVPCId copyWith(void Function(InvalidVPCId) updates) =>
      super.copyWith((message) => updates(message as InvalidVPCId))
          as InvalidVPCId;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidVPCId create() => InvalidVPCId._();
  @$core.override
  InvalidVPCId createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidVPCId getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidVPCId>(create);
  static InvalidVPCId? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KeySigningKey extends $pb.GeneratedMessage {
  factory KeySigningKey({
    $core.String? status,
    $core.String? lastmodifieddate,
    $core.String? dnskeyrecord,
    $core.String? statusmessage,
    $core.String? dsrecord,
    $core.String? publickey,
    $core.String? kmsarn,
    $core.String? digestalgorithmmnemonic,
    $core.int? keytag,
    $core.int? digestalgorithmtype,
    $core.String? name,
    $core.String? digestvalue,
    $core.String? signingalgorithmmnemonic,
    $core.int? signingalgorithmtype,
    $core.String? createddate,
    $core.int? flag,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (lastmodifieddate != null) result.lastmodifieddate = lastmodifieddate;
    if (dnskeyrecord != null) result.dnskeyrecord = dnskeyrecord;
    if (statusmessage != null) result.statusmessage = statusmessage;
    if (dsrecord != null) result.dsrecord = dsrecord;
    if (publickey != null) result.publickey = publickey;
    if (kmsarn != null) result.kmsarn = kmsarn;
    if (digestalgorithmmnemonic != null)
      result.digestalgorithmmnemonic = digestalgorithmmnemonic;
    if (keytag != null) result.keytag = keytag;
    if (digestalgorithmtype != null)
      result.digestalgorithmtype = digestalgorithmtype;
    if (name != null) result.name = name;
    if (digestvalue != null) result.digestvalue = digestvalue;
    if (signingalgorithmmnemonic != null)
      result.signingalgorithmmnemonic = signingalgorithmmnemonic;
    if (signingalgorithmtype != null)
      result.signingalgorithmtype = signingalgorithmtype;
    if (createddate != null) result.createddate = createddate;
    if (flag != null) result.flag = flag;
    return result;
  }

  KeySigningKey._();

  factory KeySigningKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeySigningKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeySigningKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(6222352, _omitFieldNames ? '' : 'status')
    ..aOS(24249427, _omitFieldNames ? '' : 'lastmodifieddate')
    ..aOS(33128395, _omitFieldNames ? '' : 'dnskeyrecord')
    ..aOS(72590095, _omitFieldNames ? '' : 'statusmessage')
    ..aOS(123290342, _omitFieldNames ? '' : 'dsrecord')
    ..aOS(167335776, _omitFieldNames ? '' : 'publickey')
    ..aOS(205584974, _omitFieldNames ? '' : 'kmsarn')
    ..aOS(209907099, _omitFieldNames ? '' : 'digestalgorithmmnemonic')
    ..aI(245478853, _omitFieldNames ? '' : 'keytag')
    ..aI(255175667, _omitFieldNames ? '' : 'digestalgorithmtype')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(277868387, _omitFieldNames ? '' : 'digestvalue')
    ..aOS(379080734, _omitFieldNames ? '' : 'signingalgorithmmnemonic')
    ..aI(407768874, _omitFieldNames ? '' : 'signingalgorithmtype')
    ..aOS(416929840, _omitFieldNames ? '' : 'createddate')
    ..aI(504820984, _omitFieldNames ? '' : 'flag')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKey copyWith(void Function(KeySigningKey) updates) =>
      super.copyWith((message) => updates(message as KeySigningKey))
          as KeySigningKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeySigningKey create() => KeySigningKey._();
  @$core.override
  KeySigningKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeySigningKey getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeySigningKey>(create);
  static KeySigningKey? _defaultInstance;

  @$pb.TagNumber(6222352)
  $core.String get status => $_getSZ(0);
  @$pb.TagNumber(6222352)
  set status($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(24249427)
  $core.String get lastmodifieddate => $_getSZ(1);
  @$pb.TagNumber(24249427)
  set lastmodifieddate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(24249427)
  $core.bool hasLastmodifieddate() => $_has(1);
  @$pb.TagNumber(24249427)
  void clearLastmodifieddate() => $_clearField(24249427);

  @$pb.TagNumber(33128395)
  $core.String get dnskeyrecord => $_getSZ(2);
  @$pb.TagNumber(33128395)
  set dnskeyrecord($core.String value) => $_setString(2, value);
  @$pb.TagNumber(33128395)
  $core.bool hasDnskeyrecord() => $_has(2);
  @$pb.TagNumber(33128395)
  void clearDnskeyrecord() => $_clearField(33128395);

  @$pb.TagNumber(72590095)
  $core.String get statusmessage => $_getSZ(3);
  @$pb.TagNumber(72590095)
  set statusmessage($core.String value) => $_setString(3, value);
  @$pb.TagNumber(72590095)
  $core.bool hasStatusmessage() => $_has(3);
  @$pb.TagNumber(72590095)
  void clearStatusmessage() => $_clearField(72590095);

  @$pb.TagNumber(123290342)
  $core.String get dsrecord => $_getSZ(4);
  @$pb.TagNumber(123290342)
  set dsrecord($core.String value) => $_setString(4, value);
  @$pb.TagNumber(123290342)
  $core.bool hasDsrecord() => $_has(4);
  @$pb.TagNumber(123290342)
  void clearDsrecord() => $_clearField(123290342);

  @$pb.TagNumber(167335776)
  $core.String get publickey => $_getSZ(5);
  @$pb.TagNumber(167335776)
  set publickey($core.String value) => $_setString(5, value);
  @$pb.TagNumber(167335776)
  $core.bool hasPublickey() => $_has(5);
  @$pb.TagNumber(167335776)
  void clearPublickey() => $_clearField(167335776);

  @$pb.TagNumber(205584974)
  $core.String get kmsarn => $_getSZ(6);
  @$pb.TagNumber(205584974)
  set kmsarn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(205584974)
  $core.bool hasKmsarn() => $_has(6);
  @$pb.TagNumber(205584974)
  void clearKmsarn() => $_clearField(205584974);

  @$pb.TagNumber(209907099)
  $core.String get digestalgorithmmnemonic => $_getSZ(7);
  @$pb.TagNumber(209907099)
  set digestalgorithmmnemonic($core.String value) => $_setString(7, value);
  @$pb.TagNumber(209907099)
  $core.bool hasDigestalgorithmmnemonic() => $_has(7);
  @$pb.TagNumber(209907099)
  void clearDigestalgorithmmnemonic() => $_clearField(209907099);

  @$pb.TagNumber(245478853)
  $core.int get keytag => $_getIZ(8);
  @$pb.TagNumber(245478853)
  set keytag($core.int value) => $_setSignedInt32(8, value);
  @$pb.TagNumber(245478853)
  $core.bool hasKeytag() => $_has(8);
  @$pb.TagNumber(245478853)
  void clearKeytag() => $_clearField(245478853);

  @$pb.TagNumber(255175667)
  $core.int get digestalgorithmtype => $_getIZ(9);
  @$pb.TagNumber(255175667)
  set digestalgorithmtype($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(255175667)
  $core.bool hasDigestalgorithmtype() => $_has(9);
  @$pb.TagNumber(255175667)
  void clearDigestalgorithmtype() => $_clearField(255175667);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(10);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(10, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(10);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(277868387)
  $core.String get digestvalue => $_getSZ(11);
  @$pb.TagNumber(277868387)
  set digestvalue($core.String value) => $_setString(11, value);
  @$pb.TagNumber(277868387)
  $core.bool hasDigestvalue() => $_has(11);
  @$pb.TagNumber(277868387)
  void clearDigestvalue() => $_clearField(277868387);

  @$pb.TagNumber(379080734)
  $core.String get signingalgorithmmnemonic => $_getSZ(12);
  @$pb.TagNumber(379080734)
  set signingalgorithmmnemonic($core.String value) => $_setString(12, value);
  @$pb.TagNumber(379080734)
  $core.bool hasSigningalgorithmmnemonic() => $_has(12);
  @$pb.TagNumber(379080734)
  void clearSigningalgorithmmnemonic() => $_clearField(379080734);

  @$pb.TagNumber(407768874)
  $core.int get signingalgorithmtype => $_getIZ(13);
  @$pb.TagNumber(407768874)
  set signingalgorithmtype($core.int value) => $_setSignedInt32(13, value);
  @$pb.TagNumber(407768874)
  $core.bool hasSigningalgorithmtype() => $_has(13);
  @$pb.TagNumber(407768874)
  void clearSigningalgorithmtype() => $_clearField(407768874);

  @$pb.TagNumber(416929840)
  $core.String get createddate => $_getSZ(14);
  @$pb.TagNumber(416929840)
  set createddate($core.String value) => $_setString(14, value);
  @$pb.TagNumber(416929840)
  $core.bool hasCreateddate() => $_has(14);
  @$pb.TagNumber(416929840)
  void clearCreateddate() => $_clearField(416929840);

  @$pb.TagNumber(504820984)
  $core.int get flag => $_getIZ(15);
  @$pb.TagNumber(504820984)
  set flag($core.int value) => $_setSignedInt32(15, value);
  @$pb.TagNumber(504820984)
  $core.bool hasFlag() => $_has(15);
  @$pb.TagNumber(504820984)
  void clearFlag() => $_clearField(504820984);
}

class KeySigningKeyAlreadyExists extends $pb.GeneratedMessage {
  factory KeySigningKeyAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KeySigningKeyAlreadyExists._();

  factory KeySigningKeyAlreadyExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeySigningKeyAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeySigningKeyAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKeyAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKeyAlreadyExists copyWith(
          void Function(KeySigningKeyAlreadyExists) updates) =>
      super.copyWith(
              (message) => updates(message as KeySigningKeyAlreadyExists))
          as KeySigningKeyAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeySigningKeyAlreadyExists create() => KeySigningKeyAlreadyExists._();
  @$core.override
  KeySigningKeyAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeySigningKeyAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeySigningKeyAlreadyExists>(create);
  static KeySigningKeyAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KeySigningKeyInParentDSRecord extends $pb.GeneratedMessage {
  factory KeySigningKeyInParentDSRecord({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KeySigningKeyInParentDSRecord._();

  factory KeySigningKeyInParentDSRecord.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeySigningKeyInParentDSRecord.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeySigningKeyInParentDSRecord',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKeyInParentDSRecord clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKeyInParentDSRecord copyWith(
          void Function(KeySigningKeyInParentDSRecord) updates) =>
      super.copyWith(
              (message) => updates(message as KeySigningKeyInParentDSRecord))
          as KeySigningKeyInParentDSRecord;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeySigningKeyInParentDSRecord create() =>
      KeySigningKeyInParentDSRecord._();
  @$core.override
  KeySigningKeyInParentDSRecord createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeySigningKeyInParentDSRecord getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeySigningKeyInParentDSRecord>(create);
  static KeySigningKeyInParentDSRecord? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KeySigningKeyInUse extends $pb.GeneratedMessage {
  factory KeySigningKeyInUse({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KeySigningKeyInUse._();

  factory KeySigningKeyInUse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeySigningKeyInUse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeySigningKeyInUse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKeyInUse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKeyInUse copyWith(void Function(KeySigningKeyInUse) updates) =>
      super.copyWith((message) => updates(message as KeySigningKeyInUse))
          as KeySigningKeyInUse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeySigningKeyInUse create() => KeySigningKeyInUse._();
  @$core.override
  KeySigningKeyInUse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeySigningKeyInUse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<KeySigningKeyInUse>(create);
  static KeySigningKeyInUse? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class KeySigningKeyWithActiveStatusNotFound extends $pb.GeneratedMessage {
  factory KeySigningKeyWithActiveStatusNotFound({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  KeySigningKeyWithActiveStatusNotFound._();

  factory KeySigningKeyWithActiveStatusNotFound.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory KeySigningKeyWithActiveStatusNotFound.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'KeySigningKeyWithActiveStatusNotFound',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKeyWithActiveStatusNotFound clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  KeySigningKeyWithActiveStatusNotFound copyWith(
          void Function(KeySigningKeyWithActiveStatusNotFound) updates) =>
      super.copyWith((message) =>
              updates(message as KeySigningKeyWithActiveStatusNotFound))
          as KeySigningKeyWithActiveStatusNotFound;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static KeySigningKeyWithActiveStatusNotFound create() =>
      KeySigningKeyWithActiveStatusNotFound._();
  @$core.override
  KeySigningKeyWithActiveStatusNotFound createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static KeySigningKeyWithActiveStatusNotFound getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          KeySigningKeyWithActiveStatusNotFound>(create);
  static KeySigningKeyWithActiveStatusNotFound? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class LastVPCAssociation extends $pb.GeneratedMessage {
  factory LastVPCAssociation({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  LastVPCAssociation._();

  factory LastVPCAssociation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LastVPCAssociation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LastVPCAssociation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LastVPCAssociation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LastVPCAssociation copyWith(void Function(LastVPCAssociation) updates) =>
      super.copyWith((message) => updates(message as LastVPCAssociation))
          as LastVPCAssociation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LastVPCAssociation create() => LastVPCAssociation._();
  @$core.override
  LastVPCAssociation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LastVPCAssociation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LastVPCAssociation>(create);
  static LastVPCAssociation? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class LimitsExceeded extends $pb.GeneratedMessage {
  factory LimitsExceeded({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  LimitsExceeded._();

  factory LimitsExceeded.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LimitsExceeded.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LimitsExceeded',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LimitsExceeded clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LimitsExceeded copyWith(void Function(LimitsExceeded) updates) =>
      super.copyWith((message) => updates(message as LimitsExceeded))
          as LimitsExceeded;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LimitsExceeded create() => LimitsExceeded._();
  @$core.override
  LimitsExceeded createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LimitsExceeded getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LimitsExceeded>(create);
  static LimitsExceeded? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class LinkedService extends $pb.GeneratedMessage {
  factory LinkedService({
    $core.String? description,
    $core.String? serviceprincipal,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (serviceprincipal != null) result.serviceprincipal = serviceprincipal;
    return result;
  }

  LinkedService._();

  factory LinkedService.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LinkedService.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LinkedService',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(146694383, _omitFieldNames ? '' : 'serviceprincipal')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LinkedService clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LinkedService copyWith(void Function(LinkedService) updates) =>
      super.copyWith((message) => updates(message as LinkedService))
          as LinkedService;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LinkedService create() => LinkedService._();
  @$core.override
  LinkedService createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LinkedService getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LinkedService>(create);
  static LinkedService? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(146694383)
  $core.String get serviceprincipal => $_getSZ(1);
  @$pb.TagNumber(146694383)
  set serviceprincipal($core.String value) => $_setString(1, value);
  @$pb.TagNumber(146694383)
  $core.bool hasServiceprincipal() => $_has(1);
  @$pb.TagNumber(146694383)
  void clearServiceprincipal() => $_clearField(146694383);
}

class ListCidrBlocksRequest extends $pb.GeneratedMessage {
  factory ListCidrBlocksRequest({
    $core.String? collectionid,
    $core.String? locationname,
    $core.String? nexttoken,
    $core.String? maxresults,
  }) {
    final result = create();
    if (collectionid != null) result.collectionid = collectionid;
    if (locationname != null) result.locationname = locationname;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListCidrBlocksRequest._();

  factory ListCidrBlocksRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCidrBlocksRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCidrBlocksRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(128052453, _omitFieldNames ? '' : 'collectionid')
    ..aOS(158186566, _omitFieldNames ? '' : 'locationname')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrBlocksRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrBlocksRequest copyWith(
          void Function(ListCidrBlocksRequest) updates) =>
      super.copyWith((message) => updates(message as ListCidrBlocksRequest))
          as ListCidrBlocksRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCidrBlocksRequest create() => ListCidrBlocksRequest._();
  @$core.override
  ListCidrBlocksRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCidrBlocksRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCidrBlocksRequest>(create);
  static ListCidrBlocksRequest? _defaultInstance;

  @$pb.TagNumber(128052453)
  $core.String get collectionid => $_getSZ(0);
  @$pb.TagNumber(128052453)
  set collectionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(128052453)
  $core.bool hasCollectionid() => $_has(0);
  @$pb.TagNumber(128052453)
  void clearCollectionid() => $_clearField(128052453);

  @$pb.TagNumber(158186566)
  $core.String get locationname => $_getSZ(1);
  @$pb.TagNumber(158186566)
  set locationname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(158186566)
  $core.bool hasLocationname() => $_has(1);
  @$pb.TagNumber(158186566)
  void clearLocationname() => $_clearField(158186566);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.String get maxresults => $_getSZ(3);
  @$pb.TagNumber(275174450)
  set maxresults($core.String value) => $_setString(3, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(3);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListCidrBlocksResponse extends $pb.GeneratedMessage {
  factory ListCidrBlocksResponse({
    $core.Iterable<CidrBlockSummary>? cidrblocks,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (cidrblocks != null) result.cidrblocks.addAll(cidrblocks);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListCidrBlocksResponse._();

  factory ListCidrBlocksResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCidrBlocksResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCidrBlocksResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<CidrBlockSummary>(134660738, _omitFieldNames ? '' : 'cidrblocks',
        subBuilder: CidrBlockSummary.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrBlocksResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrBlocksResponse copyWith(
          void Function(ListCidrBlocksResponse) updates) =>
      super.copyWith((message) => updates(message as ListCidrBlocksResponse))
          as ListCidrBlocksResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCidrBlocksResponse create() => ListCidrBlocksResponse._();
  @$core.override
  ListCidrBlocksResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCidrBlocksResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCidrBlocksResponse>(create);
  static ListCidrBlocksResponse? _defaultInstance;

  @$pb.TagNumber(134660738)
  $pb.PbList<CidrBlockSummary> get cidrblocks => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListCidrCollectionsRequest extends $pb.GeneratedMessage {
  factory ListCidrCollectionsRequest({
    $core.String? nexttoken,
    $core.String? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListCidrCollectionsRequest._();

  factory ListCidrCollectionsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCidrCollectionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCidrCollectionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrCollectionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrCollectionsRequest copyWith(
          void Function(ListCidrCollectionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListCidrCollectionsRequest))
          as ListCidrCollectionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCidrCollectionsRequest create() => ListCidrCollectionsRequest._();
  @$core.override
  ListCidrCollectionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCidrCollectionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCidrCollectionsRequest>(create);
  static ListCidrCollectionsRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.String get maxresults => $_getSZ(1);
  @$pb.TagNumber(275174450)
  set maxresults($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListCidrCollectionsResponse extends $pb.GeneratedMessage {
  factory ListCidrCollectionsResponse({
    $core.String? nexttoken,
    $core.Iterable<CollectionSummary>? cidrcollections,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (cidrcollections != null) result.cidrcollections.addAll(cidrcollections);
    return result;
  }

  ListCidrCollectionsResponse._();

  factory ListCidrCollectionsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCidrCollectionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCidrCollectionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<CollectionSummary>(
        502378541, _omitFieldNames ? '' : 'cidrcollections',
        subBuilder: CollectionSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrCollectionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrCollectionsResponse copyWith(
          void Function(ListCidrCollectionsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListCidrCollectionsResponse))
          as ListCidrCollectionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCidrCollectionsResponse create() =>
      ListCidrCollectionsResponse._();
  @$core.override
  ListCidrCollectionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCidrCollectionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCidrCollectionsResponse>(create);
  static ListCidrCollectionsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(502378541)
  $pb.PbList<CollectionSummary> get cidrcollections => $_getList(1);
}

class ListCidrLocationsRequest extends $pb.GeneratedMessage {
  factory ListCidrLocationsRequest({
    $core.String? collectionid,
    $core.String? nexttoken,
    $core.String? maxresults,
  }) {
    final result = create();
    if (collectionid != null) result.collectionid = collectionid;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListCidrLocationsRequest._();

  factory ListCidrLocationsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCidrLocationsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCidrLocationsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(128052453, _omitFieldNames ? '' : 'collectionid')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrLocationsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrLocationsRequest copyWith(
          void Function(ListCidrLocationsRequest) updates) =>
      super.copyWith((message) => updates(message as ListCidrLocationsRequest))
          as ListCidrLocationsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCidrLocationsRequest create() => ListCidrLocationsRequest._();
  @$core.override
  ListCidrLocationsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCidrLocationsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCidrLocationsRequest>(create);
  static ListCidrLocationsRequest? _defaultInstance;

  @$pb.TagNumber(128052453)
  $core.String get collectionid => $_getSZ(0);
  @$pb.TagNumber(128052453)
  set collectionid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(128052453)
  $core.bool hasCollectionid() => $_has(0);
  @$pb.TagNumber(128052453)
  void clearCollectionid() => $_clearField(128052453);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.String get maxresults => $_getSZ(2);
  @$pb.TagNumber(275174450)
  set maxresults($core.String value) => $_setString(2, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(2);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListCidrLocationsResponse extends $pb.GeneratedMessage {
  factory ListCidrLocationsResponse({
    $core.String? nexttoken,
    $core.Iterable<LocationSummary>? cidrlocations,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (cidrlocations != null) result.cidrlocations.addAll(cidrlocations);
    return result;
  }

  ListCidrLocationsResponse._();

  factory ListCidrLocationsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListCidrLocationsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListCidrLocationsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<LocationSummary>(480481962, _omitFieldNames ? '' : 'cidrlocations',
        subBuilder: LocationSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrLocationsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListCidrLocationsResponse copyWith(
          void Function(ListCidrLocationsResponse) updates) =>
      super.copyWith((message) => updates(message as ListCidrLocationsResponse))
          as ListCidrLocationsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListCidrLocationsResponse create() => ListCidrLocationsResponse._();
  @$core.override
  ListCidrLocationsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListCidrLocationsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListCidrLocationsResponse>(create);
  static ListCidrLocationsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(480481962)
  $pb.PbList<LocationSummary> get cidrlocations => $_getList(1);
}

class ListGeoLocationsRequest extends $pb.GeneratedMessage {
  factory ListGeoLocationsRequest({
    $core.String? startcountrycode,
    $core.String? startsubdivisioncode,
    $core.String? startcontinentcode,
    $core.String? maxitems,
  }) {
    final result = create();
    if (startcountrycode != null) result.startcountrycode = startcountrycode;
    if (startsubdivisioncode != null)
      result.startsubdivisioncode = startsubdivisioncode;
    if (startcontinentcode != null)
      result.startcontinentcode = startcontinentcode;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListGeoLocationsRequest._();

  factory ListGeoLocationsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListGeoLocationsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListGeoLocationsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(182995791, _omitFieldNames ? '' : 'startcountrycode')
    ..aOS(206846290, _omitFieldNames ? '' : 'startsubdivisioncode')
    ..aOS(272774541, _omitFieldNames ? '' : 'startcontinentcode')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGeoLocationsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGeoLocationsRequest copyWith(
          void Function(ListGeoLocationsRequest) updates) =>
      super.copyWith((message) => updates(message as ListGeoLocationsRequest))
          as ListGeoLocationsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListGeoLocationsRequest create() => ListGeoLocationsRequest._();
  @$core.override
  ListGeoLocationsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListGeoLocationsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListGeoLocationsRequest>(create);
  static ListGeoLocationsRequest? _defaultInstance;

  @$pb.TagNumber(182995791)
  $core.String get startcountrycode => $_getSZ(0);
  @$pb.TagNumber(182995791)
  set startcountrycode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(182995791)
  $core.bool hasStartcountrycode() => $_has(0);
  @$pb.TagNumber(182995791)
  void clearStartcountrycode() => $_clearField(182995791);

  @$pb.TagNumber(206846290)
  $core.String get startsubdivisioncode => $_getSZ(1);
  @$pb.TagNumber(206846290)
  set startsubdivisioncode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(206846290)
  $core.bool hasStartsubdivisioncode() => $_has(1);
  @$pb.TagNumber(206846290)
  void clearStartsubdivisioncode() => $_clearField(206846290);

  @$pb.TagNumber(272774541)
  $core.String get startcontinentcode => $_getSZ(2);
  @$pb.TagNumber(272774541)
  set startcontinentcode($core.String value) => $_setString(2, value);
  @$pb.TagNumber(272774541)
  $core.bool hasStartcontinentcode() => $_has(2);
  @$pb.TagNumber(272774541)
  void clearStartcontinentcode() => $_clearField(272774541);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListGeoLocationsResponse extends $pb.GeneratedMessage {
  factory ListGeoLocationsResponse({
    $core.String? nextcontinentcode,
    $core.String? nextcountrycode,
    $core.bool? istruncated,
    $core.String? nextsubdivisioncode,
    $core.Iterable<GeoLocationDetails>? geolocationdetailslist,
    $core.String? maxitems,
  }) {
    final result = create();
    if (nextcontinentcode != null) result.nextcontinentcode = nextcontinentcode;
    if (nextcountrycode != null) result.nextcountrycode = nextcountrycode;
    if (istruncated != null) result.istruncated = istruncated;
    if (nextsubdivisioncode != null)
      result.nextsubdivisioncode = nextsubdivisioncode;
    if (geolocationdetailslist != null)
      result.geolocationdetailslist.addAll(geolocationdetailslist);
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListGeoLocationsResponse._();

  factory ListGeoLocationsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListGeoLocationsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListGeoLocationsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(14234018, _omitFieldNames ? '' : 'nextcontinentcode')
    ..aOS(147948496, _omitFieldNames ? '' : 'nextcountrycode')
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(385449253, _omitFieldNames ? '' : 'nextsubdivisioncode')
    ..pPM<GeoLocationDetails>(
        448447430, _omitFieldNames ? '' : 'geolocationdetailslist',
        subBuilder: GeoLocationDetails.create)
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGeoLocationsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListGeoLocationsResponse copyWith(
          void Function(ListGeoLocationsResponse) updates) =>
      super.copyWith((message) => updates(message as ListGeoLocationsResponse))
          as ListGeoLocationsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListGeoLocationsResponse create() => ListGeoLocationsResponse._();
  @$core.override
  ListGeoLocationsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListGeoLocationsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListGeoLocationsResponse>(create);
  static ListGeoLocationsResponse? _defaultInstance;

  @$pb.TagNumber(14234018)
  $core.String get nextcontinentcode => $_getSZ(0);
  @$pb.TagNumber(14234018)
  set nextcontinentcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(14234018)
  $core.bool hasNextcontinentcode() => $_has(0);
  @$pb.TagNumber(14234018)
  void clearNextcontinentcode() => $_clearField(14234018);

  @$pb.TagNumber(147948496)
  $core.String get nextcountrycode => $_getSZ(1);
  @$pb.TagNumber(147948496)
  set nextcountrycode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(147948496)
  $core.bool hasNextcountrycode() => $_has(1);
  @$pb.TagNumber(147948496)
  void clearNextcountrycode() => $_clearField(147948496);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(2);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(2);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(385449253)
  $core.String get nextsubdivisioncode => $_getSZ(3);
  @$pb.TagNumber(385449253)
  set nextsubdivisioncode($core.String value) => $_setString(3, value);
  @$pb.TagNumber(385449253)
  $core.bool hasNextsubdivisioncode() => $_has(3);
  @$pb.TagNumber(385449253)
  void clearNextsubdivisioncode() => $_clearField(385449253);

  @$pb.TagNumber(448447430)
  $pb.PbList<GeoLocationDetails> get geolocationdetailslist => $_getList(4);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(5);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(5, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(5);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListHealthChecksRequest extends $pb.GeneratedMessage {
  factory ListHealthChecksRequest({
    $core.String? marker,
    $core.String? maxitems,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListHealthChecksRequest._();

  factory ListHealthChecksRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListHealthChecksRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHealthChecksRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHealthChecksRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHealthChecksRequest copyWith(
          void Function(ListHealthChecksRequest) updates) =>
      super.copyWith((message) => updates(message as ListHealthChecksRequest))
          as ListHealthChecksRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListHealthChecksRequest create() => ListHealthChecksRequest._();
  @$core.override
  ListHealthChecksRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListHealthChecksRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHealthChecksRequest>(create);
  static ListHealthChecksRequest? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(1);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(1, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(1);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListHealthChecksResponse extends $pb.GeneratedMessage {
  factory ListHealthChecksResponse({
    $core.String? marker,
    $core.bool? istruncated,
    $core.String? maxitems,
    $core.Iterable<HealthCheck>? healthchecks,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (istruncated != null) result.istruncated = istruncated;
    if (maxitems != null) result.maxitems = maxitems;
    if (healthchecks != null) result.healthchecks.addAll(healthchecks);
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListHealthChecksResponse._();

  factory ListHealthChecksResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListHealthChecksResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHealthChecksResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..pPM<HealthCheck>(516432759, _omitFieldNames ? '' : 'healthchecks',
        subBuilder: HealthCheck.create)
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHealthChecksResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHealthChecksResponse copyWith(
          void Function(ListHealthChecksResponse) updates) =>
      super.copyWith((message) => updates(message as ListHealthChecksResponse))
          as ListHealthChecksResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListHealthChecksResponse create() => ListHealthChecksResponse._();
  @$core.override
  ListHealthChecksResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListHealthChecksResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHealthChecksResponse>(create);
  static ListHealthChecksResponse? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(1);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(1);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(2);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(2, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(2);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);

  @$pb.TagNumber(516432759)
  $pb.PbList<HealthCheck> get healthchecks => $_getList(3);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(4);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(4, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(4);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListHostedZonesByNameRequest extends $pb.GeneratedMessage {
  factory ListHostedZonesByNameRequest({
    $core.String? dnsname,
    $core.String? hostedzoneid,
    $core.String? maxitems,
  }) {
    final result = create();
    if (dnsname != null) result.dnsname = dnsname;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListHostedZonesByNameRequest._();

  factory ListHostedZonesByNameRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListHostedZonesByNameRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHostedZonesByNameRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(171901432, _omitFieldNames ? '' : 'dnsname')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesByNameRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesByNameRequest copyWith(
          void Function(ListHostedZonesByNameRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListHostedZonesByNameRequest))
          as ListHostedZonesByNameRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListHostedZonesByNameRequest create() =>
      ListHostedZonesByNameRequest._();
  @$core.override
  ListHostedZonesByNameRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListHostedZonesByNameRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHostedZonesByNameRequest>(create);
  static ListHostedZonesByNameRequest? _defaultInstance;

  @$pb.TagNumber(171901432)
  $core.String get dnsname => $_getSZ(0);
  @$pb.TagNumber(171901432)
  set dnsname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(171901432)
  $core.bool hasDnsname() => $_has(0);
  @$pb.TagNumber(171901432)
  void clearDnsname() => $_clearField(171901432);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(2);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(2, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(2);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListHostedZonesByNameResponse extends $pb.GeneratedMessage {
  factory ListHostedZonesByNameResponse({
    $core.Iterable<HostedZone>? hostedzones,
    $core.String? nexthostedzoneid,
    $core.String? dnsname,
    $core.bool? istruncated,
    $core.String? hostedzoneid,
    $core.String? nextdnsname,
    $core.String? maxitems,
  }) {
    final result = create();
    if (hostedzones != null) result.hostedzones.addAll(hostedzones);
    if (nexthostedzoneid != null) result.nexthostedzoneid = nexthostedzoneid;
    if (dnsname != null) result.dnsname = dnsname;
    if (istruncated != null) result.istruncated = istruncated;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (nextdnsname != null) result.nextdnsname = nextdnsname;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListHostedZonesByNameResponse._();

  factory ListHostedZonesByNameResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListHostedZonesByNameResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHostedZonesByNameResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<HostedZone>(86735402, _omitFieldNames ? '' : 'hostedzones',
        subBuilder: HostedZone.create)
    ..aOS(162450165, _omitFieldNames ? '' : 'nexthostedzoneid')
    ..aOS(171901432, _omitFieldNames ? '' : 'dnsname')
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(488331797, _omitFieldNames ? '' : 'nextdnsname')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesByNameResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesByNameResponse copyWith(
          void Function(ListHostedZonesByNameResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListHostedZonesByNameResponse))
          as ListHostedZonesByNameResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListHostedZonesByNameResponse create() =>
      ListHostedZonesByNameResponse._();
  @$core.override
  ListHostedZonesByNameResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListHostedZonesByNameResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHostedZonesByNameResponse>(create);
  static ListHostedZonesByNameResponse? _defaultInstance;

  @$pb.TagNumber(86735402)
  $pb.PbList<HostedZone> get hostedzones => $_getList(0);

  @$pb.TagNumber(162450165)
  $core.String get nexthostedzoneid => $_getSZ(1);
  @$pb.TagNumber(162450165)
  set nexthostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(162450165)
  $core.bool hasNexthostedzoneid() => $_has(1);
  @$pb.TagNumber(162450165)
  void clearNexthostedzoneid() => $_clearField(162450165);

  @$pb.TagNumber(171901432)
  $core.String get dnsname => $_getSZ(2);
  @$pb.TagNumber(171901432)
  set dnsname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(171901432)
  $core.bool hasDnsname() => $_has(2);
  @$pb.TagNumber(171901432)
  void clearDnsname() => $_clearField(171901432);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(3);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(3);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(4);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(4);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(488331797)
  $core.String get nextdnsname => $_getSZ(5);
  @$pb.TagNumber(488331797)
  set nextdnsname($core.String value) => $_setString(5, value);
  @$pb.TagNumber(488331797)
  $core.bool hasNextdnsname() => $_has(5);
  @$pb.TagNumber(488331797)
  void clearNextdnsname() => $_clearField(488331797);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(6);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(6, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(6);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListHostedZonesByVPCRequest extends $pb.GeneratedMessage {
  factory ListHostedZonesByVPCRequest({
    $core.String? nexttoken,
    $core.String? vpcid,
    VPCRegion? vpcregion,
    $core.String? maxitems,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (vpcid != null) result.vpcid = vpcid;
    if (vpcregion != null) result.vpcregion = vpcregion;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListHostedZonesByVPCRequest._();

  factory ListHostedZonesByVPCRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListHostedZonesByVPCRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHostedZonesByVPCRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(325135798, _omitFieldNames ? '' : 'vpcid')
    ..aE<VPCRegion>(474180765, _omitFieldNames ? '' : 'vpcregion',
        enumValues: VPCRegion.values)
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesByVPCRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesByVPCRequest copyWith(
          void Function(ListHostedZonesByVPCRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListHostedZonesByVPCRequest))
          as ListHostedZonesByVPCRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListHostedZonesByVPCRequest create() =>
      ListHostedZonesByVPCRequest._();
  @$core.override
  ListHostedZonesByVPCRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListHostedZonesByVPCRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHostedZonesByVPCRequest>(create);
  static ListHostedZonesByVPCRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(325135798)
  $core.String get vpcid => $_getSZ(1);
  @$pb.TagNumber(325135798)
  set vpcid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(325135798)
  $core.bool hasVpcid() => $_has(1);
  @$pb.TagNumber(325135798)
  void clearVpcid() => $_clearField(325135798);

  @$pb.TagNumber(474180765)
  VPCRegion get vpcregion => $_getN(2);
  @$pb.TagNumber(474180765)
  set vpcregion(VPCRegion value) => $_setField(474180765, value);
  @$pb.TagNumber(474180765)
  $core.bool hasVpcregion() => $_has(2);
  @$pb.TagNumber(474180765)
  void clearVpcregion() => $_clearField(474180765);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListHostedZonesByVPCResponse extends $pb.GeneratedMessage {
  factory ListHostedZonesByVPCResponse({
    $core.Iterable<HostedZoneSummary>? hostedzonesummaries,
    $core.String? nexttoken,
    $core.String? maxitems,
  }) {
    final result = create();
    if (hostedzonesummaries != null)
      result.hostedzonesummaries.addAll(hostedzonesummaries);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListHostedZonesByVPCResponse._();

  factory ListHostedZonesByVPCResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListHostedZonesByVPCResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHostedZonesByVPCResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<HostedZoneSummary>(
        111631021, _omitFieldNames ? '' : 'hostedzonesummaries',
        subBuilder: HostedZoneSummary.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesByVPCResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesByVPCResponse copyWith(
          void Function(ListHostedZonesByVPCResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListHostedZonesByVPCResponse))
          as ListHostedZonesByVPCResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListHostedZonesByVPCResponse create() =>
      ListHostedZonesByVPCResponse._();
  @$core.override
  ListHostedZonesByVPCResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListHostedZonesByVPCResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHostedZonesByVPCResponse>(create);
  static ListHostedZonesByVPCResponse? _defaultInstance;

  @$pb.TagNumber(111631021)
  $pb.PbList<HostedZoneSummary> get hostedzonesummaries => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(2);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(2, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(2);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListHostedZonesRequest extends $pb.GeneratedMessage {
  factory ListHostedZonesRequest({
    $core.String? marker,
    $core.String? delegationsetid,
    HostedZoneType? hostedzonetype,
    $core.String? maxitems,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (delegationsetid != null) result.delegationsetid = delegationsetid;
    if (hostedzonetype != null) result.hostedzonetype = hostedzonetype;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListHostedZonesRequest._();

  factory ListHostedZonesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListHostedZonesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHostedZonesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(307328801, _omitFieldNames ? '' : 'delegationsetid')
    ..aE<HostedZoneType>(409319401, _omitFieldNames ? '' : 'hostedzonetype',
        enumValues: HostedZoneType.values)
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesRequest copyWith(
          void Function(ListHostedZonesRequest) updates) =>
      super.copyWith((message) => updates(message as ListHostedZonesRequest))
          as ListHostedZonesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListHostedZonesRequest create() => ListHostedZonesRequest._();
  @$core.override
  ListHostedZonesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListHostedZonesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHostedZonesRequest>(create);
  static ListHostedZonesRequest? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(307328801)
  $core.String get delegationsetid => $_getSZ(1);
  @$pb.TagNumber(307328801)
  set delegationsetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(307328801)
  $core.bool hasDelegationsetid() => $_has(1);
  @$pb.TagNumber(307328801)
  void clearDelegationsetid() => $_clearField(307328801);

  @$pb.TagNumber(409319401)
  HostedZoneType get hostedzonetype => $_getN(2);
  @$pb.TagNumber(409319401)
  set hostedzonetype(HostedZoneType value) => $_setField(409319401, value);
  @$pb.TagNumber(409319401)
  $core.bool hasHostedzonetype() => $_has(2);
  @$pb.TagNumber(409319401)
  void clearHostedzonetype() => $_clearField(409319401);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListHostedZonesResponse extends $pb.GeneratedMessage {
  factory ListHostedZonesResponse({
    $core.Iterable<HostedZone>? hostedzones,
    $core.String? marker,
    $core.bool? istruncated,
    $core.String? maxitems,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (hostedzones != null) result.hostedzones.addAll(hostedzones);
    if (marker != null) result.marker = marker;
    if (istruncated != null) result.istruncated = istruncated;
    if (maxitems != null) result.maxitems = maxitems;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListHostedZonesResponse._();

  factory ListHostedZonesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListHostedZonesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListHostedZonesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<HostedZone>(86735402, _omitFieldNames ? '' : 'hostedzones',
        subBuilder: HostedZone.create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListHostedZonesResponse copyWith(
          void Function(ListHostedZonesResponse) updates) =>
      super.copyWith((message) => updates(message as ListHostedZonesResponse))
          as ListHostedZonesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListHostedZonesResponse create() => ListHostedZonesResponse._();
  @$core.override
  ListHostedZonesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListHostedZonesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListHostedZonesResponse>(create);
  static ListHostedZonesResponse? _defaultInstance;

  @$pb.TagNumber(86735402)
  $pb.PbList<HostedZone> get hostedzones => $_getList(0);

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(1);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(1);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(2);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(2);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(4);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(4, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(4);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListQueryLoggingConfigsRequest extends $pb.GeneratedMessage {
  factory ListQueryLoggingConfigsRequest({
    $core.String? nexttoken,
    $core.String? maxresults,
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  ListQueryLoggingConfigsRequest._();

  factory ListQueryLoggingConfigsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueryLoggingConfigsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueryLoggingConfigsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueryLoggingConfigsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueryLoggingConfigsRequest copyWith(
          void Function(ListQueryLoggingConfigsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListQueryLoggingConfigsRequest))
          as ListQueryLoggingConfigsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueryLoggingConfigsRequest create() =>
      ListQueryLoggingConfigsRequest._();
  @$core.override
  ListQueryLoggingConfigsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueryLoggingConfigsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueryLoggingConfigsRequest>(create);
  static ListQueryLoggingConfigsRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.String get maxresults => $_getSZ(1);
  @$pb.TagNumber(275174450)
  set maxresults($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(2);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(2);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class ListQueryLoggingConfigsResponse extends $pb.GeneratedMessage {
  factory ListQueryLoggingConfigsResponse({
    $core.Iterable<QueryLoggingConfig>? queryloggingconfigs,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (queryloggingconfigs != null)
      result.queryloggingconfigs.addAll(queryloggingconfigs);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListQueryLoggingConfigsResponse._();

  factory ListQueryLoggingConfigsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListQueryLoggingConfigsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListQueryLoggingConfigsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<QueryLoggingConfig>(
        87688172, _omitFieldNames ? '' : 'queryloggingconfigs',
        subBuilder: QueryLoggingConfig.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueryLoggingConfigsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListQueryLoggingConfigsResponse copyWith(
          void Function(ListQueryLoggingConfigsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListQueryLoggingConfigsResponse))
          as ListQueryLoggingConfigsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListQueryLoggingConfigsResponse create() =>
      ListQueryLoggingConfigsResponse._();
  @$core.override
  ListQueryLoggingConfigsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListQueryLoggingConfigsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListQueryLoggingConfigsResponse>(
          create);
  static ListQueryLoggingConfigsResponse? _defaultInstance;

  @$pb.TagNumber(87688172)
  $pb.PbList<QueryLoggingConfig> get queryloggingconfigs => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListResourceRecordSetsRequest extends $pb.GeneratedMessage {
  factory ListResourceRecordSetsRequest({
    $core.String? startrecordname,
    $core.String? hostedzoneid,
    RRType? startrecordtype,
    $core.String? maxitems,
    $core.String? startrecordidentifier,
  }) {
    final result = create();
    if (startrecordname != null) result.startrecordname = startrecordname;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (startrecordtype != null) result.startrecordtype = startrecordtype;
    if (maxitems != null) result.maxitems = maxitems;
    if (startrecordidentifier != null)
      result.startrecordidentifier = startrecordidentifier;
    return result;
  }

  ListResourceRecordSetsRequest._();

  factory ListResourceRecordSetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListResourceRecordSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListResourceRecordSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(145299062, _omitFieldNames ? '' : 'startrecordname')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aE<RRType>(408714791, _omitFieldNames ? '' : 'startrecordtype',
        enumValues: RRType.values)
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..aOS(518502950, _omitFieldNames ? '' : 'startrecordidentifier')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourceRecordSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourceRecordSetsRequest copyWith(
          void Function(ListResourceRecordSetsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListResourceRecordSetsRequest))
          as ListResourceRecordSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListResourceRecordSetsRequest create() =>
      ListResourceRecordSetsRequest._();
  @$core.override
  ListResourceRecordSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListResourceRecordSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListResourceRecordSetsRequest>(create);
  static ListResourceRecordSetsRequest? _defaultInstance;

  @$pb.TagNumber(145299062)
  $core.String get startrecordname => $_getSZ(0);
  @$pb.TagNumber(145299062)
  set startrecordname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(145299062)
  $core.bool hasStartrecordname() => $_has(0);
  @$pb.TagNumber(145299062)
  void clearStartrecordname() => $_clearField(145299062);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(408714791)
  RRType get startrecordtype => $_getN(2);
  @$pb.TagNumber(408714791)
  set startrecordtype(RRType value) => $_setField(408714791, value);
  @$pb.TagNumber(408714791)
  $core.bool hasStartrecordtype() => $_has(2);
  @$pb.TagNumber(408714791)
  void clearStartrecordtype() => $_clearField(408714791);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);

  @$pb.TagNumber(518502950)
  $core.String get startrecordidentifier => $_getSZ(4);
  @$pb.TagNumber(518502950)
  set startrecordidentifier($core.String value) => $_setString(4, value);
  @$pb.TagNumber(518502950)
  $core.bool hasStartrecordidentifier() => $_has(4);
  @$pb.TagNumber(518502950)
  void clearStartrecordidentifier() => $_clearField(518502950);
}

class ListResourceRecordSetsResponse extends $pb.GeneratedMessage {
  factory ListResourceRecordSetsResponse({
    $core.Iterable<ResourceRecordSet>? resourcerecordsets,
    RRType? nextrecordtype,
    $core.String? nextrecordname,
    $core.bool? istruncated,
    $core.String? nextrecordidentifier,
    $core.String? maxitems,
  }) {
    final result = create();
    if (resourcerecordsets != null)
      result.resourcerecordsets.addAll(resourcerecordsets);
    if (nextrecordtype != null) result.nextrecordtype = nextrecordtype;
    if (nextrecordname != null) result.nextrecordname = nextrecordname;
    if (istruncated != null) result.istruncated = istruncated;
    if (nextrecordidentifier != null)
      result.nextrecordidentifier = nextrecordidentifier;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListResourceRecordSetsResponse._();

  factory ListResourceRecordSetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListResourceRecordSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListResourceRecordSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<ResourceRecordSet>(
        77807302, _omitFieldNames ? '' : 'resourcerecordsets',
        subBuilder: ResourceRecordSet.create)
    ..aE<RRType>(97817846, _omitFieldNames ? '' : 'nextrecordtype',
        enumValues: RRType.values)
    ..aOS(131258783, _omitFieldNames ? '' : 'nextrecordname')
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(424069527, _omitFieldNames ? '' : 'nextrecordidentifier')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourceRecordSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListResourceRecordSetsResponse copyWith(
          void Function(ListResourceRecordSetsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListResourceRecordSetsResponse))
          as ListResourceRecordSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListResourceRecordSetsResponse create() =>
      ListResourceRecordSetsResponse._();
  @$core.override
  ListResourceRecordSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListResourceRecordSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListResourceRecordSetsResponse>(create);
  static ListResourceRecordSetsResponse? _defaultInstance;

  @$pb.TagNumber(77807302)
  $pb.PbList<ResourceRecordSet> get resourcerecordsets => $_getList(0);

  @$pb.TagNumber(97817846)
  RRType get nextrecordtype => $_getN(1);
  @$pb.TagNumber(97817846)
  set nextrecordtype(RRType value) => $_setField(97817846, value);
  @$pb.TagNumber(97817846)
  $core.bool hasNextrecordtype() => $_has(1);
  @$pb.TagNumber(97817846)
  void clearNextrecordtype() => $_clearField(97817846);

  @$pb.TagNumber(131258783)
  $core.String get nextrecordname => $_getSZ(2);
  @$pb.TagNumber(131258783)
  set nextrecordname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(131258783)
  $core.bool hasNextrecordname() => $_has(2);
  @$pb.TagNumber(131258783)
  void clearNextrecordname() => $_clearField(131258783);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(3);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(3);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(424069527)
  $core.String get nextrecordidentifier => $_getSZ(4);
  @$pb.TagNumber(424069527)
  set nextrecordidentifier($core.String value) => $_setString(4, value);
  @$pb.TagNumber(424069527)
  $core.bool hasNextrecordidentifier() => $_has(4);
  @$pb.TagNumber(424069527)
  void clearNextrecordidentifier() => $_clearField(424069527);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(5);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(5, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(5);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListReusableDelegationSetsRequest extends $pb.GeneratedMessage {
  factory ListReusableDelegationSetsRequest({
    $core.String? marker,
    $core.String? maxitems,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListReusableDelegationSetsRequest._();

  factory ListReusableDelegationSetsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListReusableDelegationSetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListReusableDelegationSetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListReusableDelegationSetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListReusableDelegationSetsRequest copyWith(
          void Function(ListReusableDelegationSetsRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ListReusableDelegationSetsRequest))
          as ListReusableDelegationSetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListReusableDelegationSetsRequest create() =>
      ListReusableDelegationSetsRequest._();
  @$core.override
  ListReusableDelegationSetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListReusableDelegationSetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListReusableDelegationSetsRequest>(
          create);
  static ListReusableDelegationSetsRequest? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(1);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(1, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(1);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListReusableDelegationSetsResponse extends $pb.GeneratedMessage {
  factory ListReusableDelegationSetsResponse({
    $core.String? marker,
    $core.bool? istruncated,
    $core.Iterable<DelegationSet>? delegationsets,
    $core.String? maxitems,
    $core.String? nextmarker,
  }) {
    final result = create();
    if (marker != null) result.marker = marker;
    if (istruncated != null) result.istruncated = istruncated;
    if (delegationsets != null) result.delegationsets.addAll(delegationsets);
    if (maxitems != null) result.maxitems = maxitems;
    if (nextmarker != null) result.nextmarker = nextmarker;
    return result;
  }

  ListReusableDelegationSetsResponse._();

  factory ListReusableDelegationSetsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListReusableDelegationSetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListReusableDelegationSetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(89353912, _omitFieldNames ? '' : 'marker')
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..pPM<DelegationSet>(477592931, _omitFieldNames ? '' : 'delegationsets',
        subBuilder: DelegationSet.create)
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..aOS(531333283, _omitFieldNames ? '' : 'nextmarker')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListReusableDelegationSetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListReusableDelegationSetsResponse copyWith(
          void Function(ListReusableDelegationSetsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListReusableDelegationSetsResponse))
          as ListReusableDelegationSetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListReusableDelegationSetsResponse create() =>
      ListReusableDelegationSetsResponse._();
  @$core.override
  ListReusableDelegationSetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListReusableDelegationSetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListReusableDelegationSetsResponse>(
          create);
  static ListReusableDelegationSetsResponse? _defaultInstance;

  @$pb.TagNumber(89353912)
  $core.String get marker => $_getSZ(0);
  @$pb.TagNumber(89353912)
  set marker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(89353912)
  $core.bool hasMarker() => $_has(0);
  @$pb.TagNumber(89353912)
  void clearMarker() => $_clearField(89353912);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(1);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(1);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(477592931)
  $pb.PbList<DelegationSet> get delegationsets => $_getList(2);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);

  @$pb.TagNumber(531333283)
  $core.String get nextmarker => $_getSZ(4);
  @$pb.TagNumber(531333283)
  set nextmarker($core.String value) => $_setString(4, value);
  @$pb.TagNumber(531333283)
  $core.bool hasNextmarker() => $_has(4);
  @$pb.TagNumber(531333283)
  void clearNextmarker() => $_clearField(531333283);
}

class ListTagsForResourceRequest extends $pb.GeneratedMessage {
  factory ListTagsForResourceRequest({
    TagResourceType? resourcetype,
    $core.String? resourceid,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (resourceid != null) result.resourceid = resourceid;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<TagResourceType>(301342558, _omitFieldNames ? '' : 'resourcetype',
        enumValues: TagResourceType.values)
    ..aOS(526146833, _omitFieldNames ? '' : 'resourceid')
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

  @$pb.TagNumber(301342558)
  TagResourceType get resourcetype => $_getN(0);
  @$pb.TagNumber(301342558)
  set resourcetype(TagResourceType value) => $_setField(301342558, value);
  @$pb.TagNumber(301342558)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(301342558)
  void clearResourcetype() => $_clearField(301342558);

  @$pb.TagNumber(526146833)
  $core.String get resourceid => $_getSZ(1);
  @$pb.TagNumber(526146833)
  set resourceid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(526146833)
  $core.bool hasResourceid() => $_has(1);
  @$pb.TagNumber(526146833)
  void clearResourceid() => $_clearField(526146833);
}

class ListTagsForResourceResponse extends $pb.GeneratedMessage {
  factory ListTagsForResourceResponse({
    ResourceTagSet? resourcetagset,
  }) {
    final result = create();
    if (resourcetagset != null) result.resourcetagset = resourcetagset;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<ResourceTagSet>(249268514, _omitFieldNames ? '' : 'resourcetagset',
        subBuilder: ResourceTagSet.create)
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

  @$pb.TagNumber(249268514)
  ResourceTagSet get resourcetagset => $_getN(0);
  @$pb.TagNumber(249268514)
  set resourcetagset(ResourceTagSet value) => $_setField(249268514, value);
  @$pb.TagNumber(249268514)
  $core.bool hasResourcetagset() => $_has(0);
  @$pb.TagNumber(249268514)
  void clearResourcetagset() => $_clearField(249268514);
  @$pb.TagNumber(249268514)
  ResourceTagSet ensureResourcetagset() => $_ensure(0);
}

class ListTagsForResourcesRequest extends $pb.GeneratedMessage {
  factory ListTagsForResourcesRequest({
    $core.Iterable<$core.String>? resourceids,
    TagResourceType? resourcetype,
  }) {
    final result = create();
    if (resourceids != null) result.resourceids.addAll(resourceids);
    if (resourcetype != null) result.resourcetype = resourcetype;
    return result;
  }

  ListTagsForResourcesRequest._();

  factory ListTagsForResourcesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForResourcesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForResourcesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPS(23528154, _omitFieldNames ? '' : 'resourceids')
    ..aE<TagResourceType>(301342558, _omitFieldNames ? '' : 'resourcetype',
        enumValues: TagResourceType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourcesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourcesRequest copyWith(
          void Function(ListTagsForResourcesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListTagsForResourcesRequest))
          as ListTagsForResourcesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForResourcesRequest create() =>
      ListTagsForResourcesRequest._();
  @$core.override
  ListTagsForResourcesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForResourcesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForResourcesRequest>(create);
  static ListTagsForResourcesRequest? _defaultInstance;

  @$pb.TagNumber(23528154)
  $pb.PbList<$core.String> get resourceids => $_getList(0);

  @$pb.TagNumber(301342558)
  TagResourceType get resourcetype => $_getN(1);
  @$pb.TagNumber(301342558)
  set resourcetype(TagResourceType value) => $_setField(301342558, value);
  @$pb.TagNumber(301342558)
  $core.bool hasResourcetype() => $_has(1);
  @$pb.TagNumber(301342558)
  void clearResourcetype() => $_clearField(301342558);
}

class ListTagsForResourcesResponse extends $pb.GeneratedMessage {
  factory ListTagsForResourcesResponse({
    $core.Iterable<ResourceTagSet>? resourcetagsets,
  }) {
    final result = create();
    if (resourcetagsets != null) result.resourcetagsets.addAll(resourcetagsets);
    return result;
  }

  ListTagsForResourcesResponse._();

  factory ListTagsForResourcesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForResourcesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForResourcesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<ResourceTagSet>(362359831, _omitFieldNames ? '' : 'resourcetagsets',
        subBuilder: ResourceTagSet.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourcesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourcesResponse copyWith(
          void Function(ListTagsForResourcesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListTagsForResourcesResponse))
          as ListTagsForResourcesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForResourcesResponse create() =>
      ListTagsForResourcesResponse._();
  @$core.override
  ListTagsForResourcesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForResourcesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForResourcesResponse>(create);
  static ListTagsForResourcesResponse? _defaultInstance;

  @$pb.TagNumber(362359831)
  $pb.PbList<ResourceTagSet> get resourcetagsets => $_getList(0);
}

class ListTrafficPoliciesRequest extends $pb.GeneratedMessage {
  factory ListTrafficPoliciesRequest({
    $core.String? trafficpolicyidmarker,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicyidmarker != null)
      result.trafficpolicyidmarker = trafficpolicyidmarker;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPoliciesRequest._();

  factory ListTrafficPoliciesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPoliciesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPoliciesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(426883336, _omitFieldNames ? '' : 'trafficpolicyidmarker')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPoliciesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPoliciesRequest copyWith(
          void Function(ListTrafficPoliciesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListTrafficPoliciesRequest))
          as ListTrafficPoliciesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPoliciesRequest create() => ListTrafficPoliciesRequest._();
  @$core.override
  ListTrafficPoliciesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPoliciesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTrafficPoliciesRequest>(create);
  static ListTrafficPoliciesRequest? _defaultInstance;

  @$pb.TagNumber(426883336)
  $core.String get trafficpolicyidmarker => $_getSZ(0);
  @$pb.TagNumber(426883336)
  set trafficpolicyidmarker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(426883336)
  $core.bool hasTrafficpolicyidmarker() => $_has(0);
  @$pb.TagNumber(426883336)
  void clearTrafficpolicyidmarker() => $_clearField(426883336);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(1);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(1, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(1);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPoliciesResponse extends $pb.GeneratedMessage {
  factory ListTrafficPoliciesResponse({
    $core.Iterable<TrafficPolicySummary>? trafficpolicysummaries,
    $core.bool? istruncated,
    $core.String? trafficpolicyidmarker,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicysummaries != null)
      result.trafficpolicysummaries.addAll(trafficpolicysummaries);
    if (istruncated != null) result.istruncated = istruncated;
    if (trafficpolicyidmarker != null)
      result.trafficpolicyidmarker = trafficpolicyidmarker;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPoliciesResponse._();

  factory ListTrafficPoliciesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPoliciesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPoliciesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<TrafficPolicySummary>(
        206945717, _omitFieldNames ? '' : 'trafficpolicysummaries',
        subBuilder: TrafficPolicySummary.create)
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(426883336, _omitFieldNames ? '' : 'trafficpolicyidmarker')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPoliciesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPoliciesResponse copyWith(
          void Function(ListTrafficPoliciesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListTrafficPoliciesResponse))
          as ListTrafficPoliciesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPoliciesResponse create() =>
      ListTrafficPoliciesResponse._();
  @$core.override
  ListTrafficPoliciesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPoliciesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTrafficPoliciesResponse>(create);
  static ListTrafficPoliciesResponse? _defaultInstance;

  @$pb.TagNumber(206945717)
  $pb.PbList<TrafficPolicySummary> get trafficpolicysummaries => $_getList(0);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(1);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(1);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(426883336)
  $core.String get trafficpolicyidmarker => $_getSZ(2);
  @$pb.TagNumber(426883336)
  set trafficpolicyidmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(426883336)
  $core.bool hasTrafficpolicyidmarker() => $_has(2);
  @$pb.TagNumber(426883336)
  void clearTrafficpolicyidmarker() => $_clearField(426883336);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPolicyInstancesByHostedZoneRequest
    extends $pb.GeneratedMessage {
  factory ListTrafficPolicyInstancesByHostedZoneRequest({
    RRType? trafficpolicyinstancetypemarker,
    $core.String? trafficpolicyinstancenamemarker,
    $core.String? hostedzoneid,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicyinstancetypemarker != null)
      result.trafficpolicyinstancetypemarker = trafficpolicyinstancetypemarker;
    if (trafficpolicyinstancenamemarker != null)
      result.trafficpolicyinstancenamemarker = trafficpolicyinstancenamemarker;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPolicyInstancesByHostedZoneRequest._();

  factory ListTrafficPolicyInstancesByHostedZoneRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPolicyInstancesByHostedZoneRequest.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPolicyInstancesByHostedZoneRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<RRType>(
        30935978, _omitFieldNames ? '' : 'trafficpolicyinstancetypemarker',
        enumValues: RRType.values)
    ..aOS(136215379, _omitFieldNames ? '' : 'trafficpolicyinstancenamemarker')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesByHostedZoneRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesByHostedZoneRequest copyWith(
          void Function(ListTrafficPolicyInstancesByHostedZoneRequest)
              updates) =>
      super.copyWith((message) =>
              updates(message as ListTrafficPolicyInstancesByHostedZoneRequest))
          as ListTrafficPolicyInstancesByHostedZoneRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesByHostedZoneRequest create() =>
      ListTrafficPolicyInstancesByHostedZoneRequest._();
  @$core.override
  ListTrafficPolicyInstancesByHostedZoneRequest createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesByHostedZoneRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListTrafficPolicyInstancesByHostedZoneRequest>(create);
  static ListTrafficPolicyInstancesByHostedZoneRequest? _defaultInstance;

  @$pb.TagNumber(30935978)
  RRType get trafficpolicyinstancetypemarker => $_getN(0);
  @$pb.TagNumber(30935978)
  set trafficpolicyinstancetypemarker(RRType value) =>
      $_setField(30935978, value);
  @$pb.TagNumber(30935978)
  $core.bool hasTrafficpolicyinstancetypemarker() => $_has(0);
  @$pb.TagNumber(30935978)
  void clearTrafficpolicyinstancetypemarker() => $_clearField(30935978);

  @$pb.TagNumber(136215379)
  $core.String get trafficpolicyinstancenamemarker => $_getSZ(1);
  @$pb.TagNumber(136215379)
  set trafficpolicyinstancenamemarker($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(136215379)
  $core.bool hasTrafficpolicyinstancenamemarker() => $_has(1);
  @$pb.TagNumber(136215379)
  void clearTrafficpolicyinstancenamemarker() => $_clearField(136215379);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(2);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(2);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPolicyInstancesByHostedZoneResponse
    extends $pb.GeneratedMessage {
  factory ListTrafficPolicyInstancesByHostedZoneResponse({
    RRType? trafficpolicyinstancetypemarker,
    $core.String? trafficpolicyinstancenamemarker,
    $core.Iterable<TrafficPolicyInstance>? trafficpolicyinstances,
    $core.bool? istruncated,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicyinstancetypemarker != null)
      result.trafficpolicyinstancetypemarker = trafficpolicyinstancetypemarker;
    if (trafficpolicyinstancenamemarker != null)
      result.trafficpolicyinstancenamemarker = trafficpolicyinstancenamemarker;
    if (trafficpolicyinstances != null)
      result.trafficpolicyinstances.addAll(trafficpolicyinstances);
    if (istruncated != null) result.istruncated = istruncated;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPolicyInstancesByHostedZoneResponse._();

  factory ListTrafficPolicyInstancesByHostedZoneResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPolicyInstancesByHostedZoneResponse.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPolicyInstancesByHostedZoneResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<RRType>(
        30935978, _omitFieldNames ? '' : 'trafficpolicyinstancetypemarker',
        enumValues: RRType.values)
    ..aOS(136215379, _omitFieldNames ? '' : 'trafficpolicyinstancenamemarker')
    ..pPM<TrafficPolicyInstance>(
        199455009, _omitFieldNames ? '' : 'trafficpolicyinstances',
        subBuilder: TrafficPolicyInstance.create)
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesByHostedZoneResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesByHostedZoneResponse copyWith(
          void Function(ListTrafficPolicyInstancesByHostedZoneResponse)
              updates) =>
      super.copyWith((message) => updates(
              message as ListTrafficPolicyInstancesByHostedZoneResponse))
          as ListTrafficPolicyInstancesByHostedZoneResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesByHostedZoneResponse create() =>
      ListTrafficPolicyInstancesByHostedZoneResponse._();
  @$core.override
  ListTrafficPolicyInstancesByHostedZoneResponse createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesByHostedZoneResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListTrafficPolicyInstancesByHostedZoneResponse>(create);
  static ListTrafficPolicyInstancesByHostedZoneResponse? _defaultInstance;

  @$pb.TagNumber(30935978)
  RRType get trafficpolicyinstancetypemarker => $_getN(0);
  @$pb.TagNumber(30935978)
  set trafficpolicyinstancetypemarker(RRType value) =>
      $_setField(30935978, value);
  @$pb.TagNumber(30935978)
  $core.bool hasTrafficpolicyinstancetypemarker() => $_has(0);
  @$pb.TagNumber(30935978)
  void clearTrafficpolicyinstancetypemarker() => $_clearField(30935978);

  @$pb.TagNumber(136215379)
  $core.String get trafficpolicyinstancenamemarker => $_getSZ(1);
  @$pb.TagNumber(136215379)
  set trafficpolicyinstancenamemarker($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(136215379)
  $core.bool hasTrafficpolicyinstancenamemarker() => $_has(1);
  @$pb.TagNumber(136215379)
  void clearTrafficpolicyinstancenamemarker() => $_clearField(136215379);

  @$pb.TagNumber(199455009)
  $pb.PbList<TrafficPolicyInstance> get trafficpolicyinstances => $_getList(2);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(3);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(3);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(4);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(4, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(4);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPolicyInstancesByPolicyRequest extends $pb.GeneratedMessage {
  factory ListTrafficPolicyInstancesByPolicyRequest({
    RRType? trafficpolicyinstancetypemarker,
    $core.String? trafficpolicyid,
    $core.String? trafficpolicyinstancenamemarker,
    $core.String? hostedzoneidmarker,
    $core.int? trafficpolicyversion,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicyinstancetypemarker != null)
      result.trafficpolicyinstancetypemarker = trafficpolicyinstancetypemarker;
    if (trafficpolicyid != null) result.trafficpolicyid = trafficpolicyid;
    if (trafficpolicyinstancenamemarker != null)
      result.trafficpolicyinstancenamemarker = trafficpolicyinstancenamemarker;
    if (hostedzoneidmarker != null)
      result.hostedzoneidmarker = hostedzoneidmarker;
    if (trafficpolicyversion != null)
      result.trafficpolicyversion = trafficpolicyversion;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPolicyInstancesByPolicyRequest._();

  factory ListTrafficPolicyInstancesByPolicyRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPolicyInstancesByPolicyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPolicyInstancesByPolicyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<RRType>(
        30935978, _omitFieldNames ? '' : 'trafficpolicyinstancetypemarker',
        enumValues: RRType.values)
    ..aOS(40235222, _omitFieldNames ? '' : 'trafficpolicyid')
    ..aOS(136215379, _omitFieldNames ? '' : 'trafficpolicyinstancenamemarker')
    ..aOS(475055952, _omitFieldNames ? '' : 'hostedzoneidmarker')
    ..aI(479078485, _omitFieldNames ? '' : 'trafficpolicyversion')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesByPolicyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesByPolicyRequest copyWith(
          void Function(ListTrafficPolicyInstancesByPolicyRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ListTrafficPolicyInstancesByPolicyRequest))
          as ListTrafficPolicyInstancesByPolicyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesByPolicyRequest create() =>
      ListTrafficPolicyInstancesByPolicyRequest._();
  @$core.override
  ListTrafficPolicyInstancesByPolicyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesByPolicyRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListTrafficPolicyInstancesByPolicyRequest>(create);
  static ListTrafficPolicyInstancesByPolicyRequest? _defaultInstance;

  @$pb.TagNumber(30935978)
  RRType get trafficpolicyinstancetypemarker => $_getN(0);
  @$pb.TagNumber(30935978)
  set trafficpolicyinstancetypemarker(RRType value) =>
      $_setField(30935978, value);
  @$pb.TagNumber(30935978)
  $core.bool hasTrafficpolicyinstancetypemarker() => $_has(0);
  @$pb.TagNumber(30935978)
  void clearTrafficpolicyinstancetypemarker() => $_clearField(30935978);

  @$pb.TagNumber(40235222)
  $core.String get trafficpolicyid => $_getSZ(1);
  @$pb.TagNumber(40235222)
  set trafficpolicyid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(40235222)
  $core.bool hasTrafficpolicyid() => $_has(1);
  @$pb.TagNumber(40235222)
  void clearTrafficpolicyid() => $_clearField(40235222);

  @$pb.TagNumber(136215379)
  $core.String get trafficpolicyinstancenamemarker => $_getSZ(2);
  @$pb.TagNumber(136215379)
  set trafficpolicyinstancenamemarker($core.String value) =>
      $_setString(2, value);
  @$pb.TagNumber(136215379)
  $core.bool hasTrafficpolicyinstancenamemarker() => $_has(2);
  @$pb.TagNumber(136215379)
  void clearTrafficpolicyinstancenamemarker() => $_clearField(136215379);

  @$pb.TagNumber(475055952)
  $core.String get hostedzoneidmarker => $_getSZ(3);
  @$pb.TagNumber(475055952)
  set hostedzoneidmarker($core.String value) => $_setString(3, value);
  @$pb.TagNumber(475055952)
  $core.bool hasHostedzoneidmarker() => $_has(3);
  @$pb.TagNumber(475055952)
  void clearHostedzoneidmarker() => $_clearField(475055952);

  @$pb.TagNumber(479078485)
  $core.int get trafficpolicyversion => $_getIZ(4);
  @$pb.TagNumber(479078485)
  set trafficpolicyversion($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(479078485)
  $core.bool hasTrafficpolicyversion() => $_has(4);
  @$pb.TagNumber(479078485)
  void clearTrafficpolicyversion() => $_clearField(479078485);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(5);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(5, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(5);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPolicyInstancesByPolicyResponse extends $pb.GeneratedMessage {
  factory ListTrafficPolicyInstancesByPolicyResponse({
    RRType? trafficpolicyinstancetypemarker,
    $core.String? trafficpolicyinstancenamemarker,
    $core.Iterable<TrafficPolicyInstance>? trafficpolicyinstances,
    $core.bool? istruncated,
    $core.String? hostedzoneidmarker,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicyinstancetypemarker != null)
      result.trafficpolicyinstancetypemarker = trafficpolicyinstancetypemarker;
    if (trafficpolicyinstancenamemarker != null)
      result.trafficpolicyinstancenamemarker = trafficpolicyinstancenamemarker;
    if (trafficpolicyinstances != null)
      result.trafficpolicyinstances.addAll(trafficpolicyinstances);
    if (istruncated != null) result.istruncated = istruncated;
    if (hostedzoneidmarker != null)
      result.hostedzoneidmarker = hostedzoneidmarker;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPolicyInstancesByPolicyResponse._();

  factory ListTrafficPolicyInstancesByPolicyResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPolicyInstancesByPolicyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPolicyInstancesByPolicyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<RRType>(
        30935978, _omitFieldNames ? '' : 'trafficpolicyinstancetypemarker',
        enumValues: RRType.values)
    ..aOS(136215379, _omitFieldNames ? '' : 'trafficpolicyinstancenamemarker')
    ..pPM<TrafficPolicyInstance>(
        199455009, _omitFieldNames ? '' : 'trafficpolicyinstances',
        subBuilder: TrafficPolicyInstance.create)
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(475055952, _omitFieldNames ? '' : 'hostedzoneidmarker')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesByPolicyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesByPolicyResponse copyWith(
          void Function(ListTrafficPolicyInstancesByPolicyResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListTrafficPolicyInstancesByPolicyResponse))
          as ListTrafficPolicyInstancesByPolicyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesByPolicyResponse create() =>
      ListTrafficPolicyInstancesByPolicyResponse._();
  @$core.override
  ListTrafficPolicyInstancesByPolicyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesByPolicyResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListTrafficPolicyInstancesByPolicyResponse>(create);
  static ListTrafficPolicyInstancesByPolicyResponse? _defaultInstance;

  @$pb.TagNumber(30935978)
  RRType get trafficpolicyinstancetypemarker => $_getN(0);
  @$pb.TagNumber(30935978)
  set trafficpolicyinstancetypemarker(RRType value) =>
      $_setField(30935978, value);
  @$pb.TagNumber(30935978)
  $core.bool hasTrafficpolicyinstancetypemarker() => $_has(0);
  @$pb.TagNumber(30935978)
  void clearTrafficpolicyinstancetypemarker() => $_clearField(30935978);

  @$pb.TagNumber(136215379)
  $core.String get trafficpolicyinstancenamemarker => $_getSZ(1);
  @$pb.TagNumber(136215379)
  set trafficpolicyinstancenamemarker($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(136215379)
  $core.bool hasTrafficpolicyinstancenamemarker() => $_has(1);
  @$pb.TagNumber(136215379)
  void clearTrafficpolicyinstancenamemarker() => $_clearField(136215379);

  @$pb.TagNumber(199455009)
  $pb.PbList<TrafficPolicyInstance> get trafficpolicyinstances => $_getList(2);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(3);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(3);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(475055952)
  $core.String get hostedzoneidmarker => $_getSZ(4);
  @$pb.TagNumber(475055952)
  set hostedzoneidmarker($core.String value) => $_setString(4, value);
  @$pb.TagNumber(475055952)
  $core.bool hasHostedzoneidmarker() => $_has(4);
  @$pb.TagNumber(475055952)
  void clearHostedzoneidmarker() => $_clearField(475055952);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(5);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(5, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(5);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPolicyInstancesRequest extends $pb.GeneratedMessage {
  factory ListTrafficPolicyInstancesRequest({
    RRType? trafficpolicyinstancetypemarker,
    $core.String? trafficpolicyinstancenamemarker,
    $core.String? hostedzoneidmarker,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicyinstancetypemarker != null)
      result.trafficpolicyinstancetypemarker = trafficpolicyinstancetypemarker;
    if (trafficpolicyinstancenamemarker != null)
      result.trafficpolicyinstancenamemarker = trafficpolicyinstancenamemarker;
    if (hostedzoneidmarker != null)
      result.hostedzoneidmarker = hostedzoneidmarker;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPolicyInstancesRequest._();

  factory ListTrafficPolicyInstancesRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPolicyInstancesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPolicyInstancesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<RRType>(
        30935978, _omitFieldNames ? '' : 'trafficpolicyinstancetypemarker',
        enumValues: RRType.values)
    ..aOS(136215379, _omitFieldNames ? '' : 'trafficpolicyinstancenamemarker')
    ..aOS(475055952, _omitFieldNames ? '' : 'hostedzoneidmarker')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesRequest copyWith(
          void Function(ListTrafficPolicyInstancesRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ListTrafficPolicyInstancesRequest))
          as ListTrafficPolicyInstancesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesRequest create() =>
      ListTrafficPolicyInstancesRequest._();
  @$core.override
  ListTrafficPolicyInstancesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTrafficPolicyInstancesRequest>(
          create);
  static ListTrafficPolicyInstancesRequest? _defaultInstance;

  @$pb.TagNumber(30935978)
  RRType get trafficpolicyinstancetypemarker => $_getN(0);
  @$pb.TagNumber(30935978)
  set trafficpolicyinstancetypemarker(RRType value) =>
      $_setField(30935978, value);
  @$pb.TagNumber(30935978)
  $core.bool hasTrafficpolicyinstancetypemarker() => $_has(0);
  @$pb.TagNumber(30935978)
  void clearTrafficpolicyinstancetypemarker() => $_clearField(30935978);

  @$pb.TagNumber(136215379)
  $core.String get trafficpolicyinstancenamemarker => $_getSZ(1);
  @$pb.TagNumber(136215379)
  set trafficpolicyinstancenamemarker($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(136215379)
  $core.bool hasTrafficpolicyinstancenamemarker() => $_has(1);
  @$pb.TagNumber(136215379)
  void clearTrafficpolicyinstancenamemarker() => $_clearField(136215379);

  @$pb.TagNumber(475055952)
  $core.String get hostedzoneidmarker => $_getSZ(2);
  @$pb.TagNumber(475055952)
  set hostedzoneidmarker($core.String value) => $_setString(2, value);
  @$pb.TagNumber(475055952)
  $core.bool hasHostedzoneidmarker() => $_has(2);
  @$pb.TagNumber(475055952)
  void clearHostedzoneidmarker() => $_clearField(475055952);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPolicyInstancesResponse extends $pb.GeneratedMessage {
  factory ListTrafficPolicyInstancesResponse({
    RRType? trafficpolicyinstancetypemarker,
    $core.String? trafficpolicyinstancenamemarker,
    $core.Iterable<TrafficPolicyInstance>? trafficpolicyinstances,
    $core.bool? istruncated,
    $core.String? hostedzoneidmarker,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicyinstancetypemarker != null)
      result.trafficpolicyinstancetypemarker = trafficpolicyinstancetypemarker;
    if (trafficpolicyinstancenamemarker != null)
      result.trafficpolicyinstancenamemarker = trafficpolicyinstancenamemarker;
    if (trafficpolicyinstances != null)
      result.trafficpolicyinstances.addAll(trafficpolicyinstances);
    if (istruncated != null) result.istruncated = istruncated;
    if (hostedzoneidmarker != null)
      result.hostedzoneidmarker = hostedzoneidmarker;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPolicyInstancesResponse._();

  factory ListTrafficPolicyInstancesResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPolicyInstancesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPolicyInstancesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<RRType>(
        30935978, _omitFieldNames ? '' : 'trafficpolicyinstancetypemarker',
        enumValues: RRType.values)
    ..aOS(136215379, _omitFieldNames ? '' : 'trafficpolicyinstancenamemarker')
    ..pPM<TrafficPolicyInstance>(
        199455009, _omitFieldNames ? '' : 'trafficpolicyinstances',
        subBuilder: TrafficPolicyInstance.create)
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(475055952, _omitFieldNames ? '' : 'hostedzoneidmarker')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyInstancesResponse copyWith(
          void Function(ListTrafficPolicyInstancesResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListTrafficPolicyInstancesResponse))
          as ListTrafficPolicyInstancesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesResponse create() =>
      ListTrafficPolicyInstancesResponse._();
  @$core.override
  ListTrafficPolicyInstancesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyInstancesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTrafficPolicyInstancesResponse>(
          create);
  static ListTrafficPolicyInstancesResponse? _defaultInstance;

  @$pb.TagNumber(30935978)
  RRType get trafficpolicyinstancetypemarker => $_getN(0);
  @$pb.TagNumber(30935978)
  set trafficpolicyinstancetypemarker(RRType value) =>
      $_setField(30935978, value);
  @$pb.TagNumber(30935978)
  $core.bool hasTrafficpolicyinstancetypemarker() => $_has(0);
  @$pb.TagNumber(30935978)
  void clearTrafficpolicyinstancetypemarker() => $_clearField(30935978);

  @$pb.TagNumber(136215379)
  $core.String get trafficpolicyinstancenamemarker => $_getSZ(1);
  @$pb.TagNumber(136215379)
  set trafficpolicyinstancenamemarker($core.String value) =>
      $_setString(1, value);
  @$pb.TagNumber(136215379)
  $core.bool hasTrafficpolicyinstancenamemarker() => $_has(1);
  @$pb.TagNumber(136215379)
  void clearTrafficpolicyinstancenamemarker() => $_clearField(136215379);

  @$pb.TagNumber(199455009)
  $pb.PbList<TrafficPolicyInstance> get trafficpolicyinstances => $_getList(2);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(3);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(3);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(475055952)
  $core.String get hostedzoneidmarker => $_getSZ(4);
  @$pb.TagNumber(475055952)
  set hostedzoneidmarker($core.String value) => $_setString(4, value);
  @$pb.TagNumber(475055952)
  $core.bool hasHostedzoneidmarker() => $_has(4);
  @$pb.TagNumber(475055952)
  void clearHostedzoneidmarker() => $_clearField(475055952);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(5);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(5, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(5);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPolicyVersionsRequest extends $pb.GeneratedMessage {
  factory ListTrafficPolicyVersionsRequest({
    $core.String? trafficpolicyversionmarker,
    $core.String? id,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicyversionmarker != null)
      result.trafficpolicyversionmarker = trafficpolicyversionmarker;
    if (id != null) result.id = id;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPolicyVersionsRequest._();

  factory ListTrafficPolicyVersionsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPolicyVersionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPolicyVersionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(228559295, _omitFieldNames ? '' : 'trafficpolicyversionmarker')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyVersionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyVersionsRequest copyWith(
          void Function(ListTrafficPolicyVersionsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListTrafficPolicyVersionsRequest))
          as ListTrafficPolicyVersionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyVersionsRequest create() =>
      ListTrafficPolicyVersionsRequest._();
  @$core.override
  ListTrafficPolicyVersionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyVersionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTrafficPolicyVersionsRequest>(
          create);
  static ListTrafficPolicyVersionsRequest? _defaultInstance;

  @$pb.TagNumber(228559295)
  $core.String get trafficpolicyversionmarker => $_getSZ(0);
  @$pb.TagNumber(228559295)
  set trafficpolicyversionmarker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(228559295)
  $core.bool hasTrafficpolicyversionmarker() => $_has(0);
  @$pb.TagNumber(228559295)
  void clearTrafficpolicyversionmarker() => $_clearField(228559295);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(2);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(2, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(2);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListTrafficPolicyVersionsResponse extends $pb.GeneratedMessage {
  factory ListTrafficPolicyVersionsResponse({
    $core.Iterable<TrafficPolicy>? trafficpolicies,
    $core.String? trafficpolicyversionmarker,
    $core.bool? istruncated,
    $core.String? maxitems,
  }) {
    final result = create();
    if (trafficpolicies != null) result.trafficpolicies.addAll(trafficpolicies);
    if (trafficpolicyversionmarker != null)
      result.trafficpolicyversionmarker = trafficpolicyversionmarker;
    if (istruncated != null) result.istruncated = istruncated;
    if (maxitems != null) result.maxitems = maxitems;
    return result;
  }

  ListTrafficPolicyVersionsResponse._();

  factory ListTrafficPolicyVersionsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTrafficPolicyVersionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTrafficPolicyVersionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pPM<TrafficPolicy>(111427421, _omitFieldNames ? '' : 'trafficpolicies',
        subBuilder: TrafficPolicy.create)
    ..aOS(228559295, _omitFieldNames ? '' : 'trafficpolicyversionmarker')
    ..aOB(242094042, _omitFieldNames ? '' : 'istruncated')
    ..aOS(506899220, _omitFieldNames ? '' : 'maxitems')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyVersionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTrafficPolicyVersionsResponse copyWith(
          void Function(ListTrafficPolicyVersionsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListTrafficPolicyVersionsResponse))
          as ListTrafficPolicyVersionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyVersionsResponse create() =>
      ListTrafficPolicyVersionsResponse._();
  @$core.override
  ListTrafficPolicyVersionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTrafficPolicyVersionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTrafficPolicyVersionsResponse>(
          create);
  static ListTrafficPolicyVersionsResponse? _defaultInstance;

  @$pb.TagNumber(111427421)
  $pb.PbList<TrafficPolicy> get trafficpolicies => $_getList(0);

  @$pb.TagNumber(228559295)
  $core.String get trafficpolicyversionmarker => $_getSZ(1);
  @$pb.TagNumber(228559295)
  set trafficpolicyversionmarker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(228559295)
  $core.bool hasTrafficpolicyversionmarker() => $_has(1);
  @$pb.TagNumber(228559295)
  void clearTrafficpolicyversionmarker() => $_clearField(228559295);

  @$pb.TagNumber(242094042)
  $core.bool get istruncated => $_getBF(2);
  @$pb.TagNumber(242094042)
  set istruncated($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(242094042)
  $core.bool hasIstruncated() => $_has(2);
  @$pb.TagNumber(242094042)
  void clearIstruncated() => $_clearField(242094042);

  @$pb.TagNumber(506899220)
  $core.String get maxitems => $_getSZ(3);
  @$pb.TagNumber(506899220)
  set maxitems($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506899220)
  $core.bool hasMaxitems() => $_has(3);
  @$pb.TagNumber(506899220)
  void clearMaxitems() => $_clearField(506899220);
}

class ListVPCAssociationAuthorizationsRequest extends $pb.GeneratedMessage {
  factory ListVPCAssociationAuthorizationsRequest({
    $core.String? nexttoken,
    $core.String? maxresults,
    $core.String? hostedzoneid,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    return result;
  }

  ListVPCAssociationAuthorizationsRequest._();

  factory ListVPCAssociationAuthorizationsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListVPCAssociationAuthorizationsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListVPCAssociationAuthorizationsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListVPCAssociationAuthorizationsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListVPCAssociationAuthorizationsRequest copyWith(
          void Function(ListVPCAssociationAuthorizationsRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ListVPCAssociationAuthorizationsRequest))
          as ListVPCAssociationAuthorizationsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListVPCAssociationAuthorizationsRequest create() =>
      ListVPCAssociationAuthorizationsRequest._();
  @$core.override
  ListVPCAssociationAuthorizationsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListVPCAssociationAuthorizationsRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListVPCAssociationAuthorizationsRequest>(create);
  static ListVPCAssociationAuthorizationsRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.String get maxresults => $_getSZ(1);
  @$pb.TagNumber(275174450)
  set maxresults($core.String value) => $_setString(1, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(2);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(2);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);
}

class ListVPCAssociationAuthorizationsResponse extends $pb.GeneratedMessage {
  factory ListVPCAssociationAuthorizationsResponse({
    $core.String? nexttoken,
    $core.String? hostedzoneid,
    $core.Iterable<VPC>? vpcs,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (vpcs != null) result.vpcs.addAll(vpcs);
    return result;
  }

  ListVPCAssociationAuthorizationsResponse._();

  factory ListVPCAssociationAuthorizationsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListVPCAssociationAuthorizationsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListVPCAssociationAuthorizationsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..pPM<VPC>(424064898, _omitFieldNames ? '' : 'vpcs', subBuilder: VPC.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListVPCAssociationAuthorizationsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListVPCAssociationAuthorizationsResponse copyWith(
          void Function(ListVPCAssociationAuthorizationsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListVPCAssociationAuthorizationsResponse))
          as ListVPCAssociationAuthorizationsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListVPCAssociationAuthorizationsResponse create() =>
      ListVPCAssociationAuthorizationsResponse._();
  @$core.override
  ListVPCAssociationAuthorizationsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListVPCAssociationAuthorizationsResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListVPCAssociationAuthorizationsResponse>(create);
  static ListVPCAssociationAuthorizationsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(1);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(1);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(424064898)
  $pb.PbList<VPC> get vpcs => $_getList(2);
}

class LocationSummary extends $pb.GeneratedMessage {
  factory LocationSummary({
    $core.String? locationname,
  }) {
    final result = create();
    if (locationname != null) result.locationname = locationname;
    return result;
  }

  LocationSummary._();

  factory LocationSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LocationSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LocationSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(158186566, _omitFieldNames ? '' : 'locationname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocationSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LocationSummary copyWith(void Function(LocationSummary) updates) =>
      super.copyWith((message) => updates(message as LocationSummary))
          as LocationSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LocationSummary create() => LocationSummary._();
  @$core.override
  LocationSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LocationSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LocationSummary>(create);
  static LocationSummary? _defaultInstance;

  @$pb.TagNumber(158186566)
  $core.String get locationname => $_getSZ(0);
  @$pb.TagNumber(158186566)
  set locationname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(158186566)
  $core.bool hasLocationname() => $_has(0);
  @$pb.TagNumber(158186566)
  void clearLocationname() => $_clearField(158186566);
}

class NoSuchChange extends $pb.GeneratedMessage {
  factory NoSuchChange({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchChange._();

  factory NoSuchChange.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchChange.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchChange',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchChange clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchChange copyWith(void Function(NoSuchChange) updates) =>
      super.copyWith((message) => updates(message as NoSuchChange))
          as NoSuchChange;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchChange create() => NoSuchChange._();
  @$core.override
  NoSuchChange createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchChange getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchChange>(create);
  static NoSuchChange? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchCidrCollectionException extends $pb.GeneratedMessage {
  factory NoSuchCidrCollectionException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchCidrCollectionException._();

  factory NoSuchCidrCollectionException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchCidrCollectionException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchCidrCollectionException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchCidrCollectionException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchCidrCollectionException copyWith(
          void Function(NoSuchCidrCollectionException) updates) =>
      super.copyWith(
              (message) => updates(message as NoSuchCidrCollectionException))
          as NoSuchCidrCollectionException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchCidrCollectionException create() =>
      NoSuchCidrCollectionException._();
  @$core.override
  NoSuchCidrCollectionException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchCidrCollectionException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchCidrCollectionException>(create);
  static NoSuchCidrCollectionException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class NoSuchCidrLocationException extends $pb.GeneratedMessage {
  factory NoSuchCidrLocationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchCidrLocationException._();

  factory NoSuchCidrLocationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchCidrLocationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchCidrLocationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchCidrLocationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchCidrLocationException copyWith(
          void Function(NoSuchCidrLocationException) updates) =>
      super.copyWith(
              (message) => updates(message as NoSuchCidrLocationException))
          as NoSuchCidrLocationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchCidrLocationException create() =>
      NoSuchCidrLocationException._();
  @$core.override
  NoSuchCidrLocationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchCidrLocationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchCidrLocationException>(create);
  static NoSuchCidrLocationException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class NoSuchCloudWatchLogsLogGroup extends $pb.GeneratedMessage {
  factory NoSuchCloudWatchLogsLogGroup({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchCloudWatchLogsLogGroup._();

  factory NoSuchCloudWatchLogsLogGroup.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchCloudWatchLogsLogGroup.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchCloudWatchLogsLogGroup',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchCloudWatchLogsLogGroup clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchCloudWatchLogsLogGroup copyWith(
          void Function(NoSuchCloudWatchLogsLogGroup) updates) =>
      super.copyWith(
              (message) => updates(message as NoSuchCloudWatchLogsLogGroup))
          as NoSuchCloudWatchLogsLogGroup;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchCloudWatchLogsLogGroup create() =>
      NoSuchCloudWatchLogsLogGroup._();
  @$core.override
  NoSuchCloudWatchLogsLogGroup createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchCloudWatchLogsLogGroup getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchCloudWatchLogsLogGroup>(create);
  static NoSuchCloudWatchLogsLogGroup? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchDelegationSet extends $pb.GeneratedMessage {
  factory NoSuchDelegationSet({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchDelegationSet._();

  factory NoSuchDelegationSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchDelegationSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchDelegationSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchDelegationSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchDelegationSet copyWith(void Function(NoSuchDelegationSet) updates) =>
      super.copyWith((message) => updates(message as NoSuchDelegationSet))
          as NoSuchDelegationSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchDelegationSet create() => NoSuchDelegationSet._();
  @$core.override
  NoSuchDelegationSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchDelegationSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchDelegationSet>(create);
  static NoSuchDelegationSet? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchGeoLocation extends $pb.GeneratedMessage {
  factory NoSuchGeoLocation({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchGeoLocation._();

  factory NoSuchGeoLocation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchGeoLocation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchGeoLocation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchGeoLocation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchGeoLocation copyWith(void Function(NoSuchGeoLocation) updates) =>
      super.copyWith((message) => updates(message as NoSuchGeoLocation))
          as NoSuchGeoLocation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchGeoLocation create() => NoSuchGeoLocation._();
  @$core.override
  NoSuchGeoLocation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchGeoLocation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchGeoLocation>(create);
  static NoSuchGeoLocation? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchHealthCheck extends $pb.GeneratedMessage {
  factory NoSuchHealthCheck({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchHealthCheck._();

  factory NoSuchHealthCheck.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchHealthCheck.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchHealthCheck',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchHealthCheck clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchHealthCheck copyWith(void Function(NoSuchHealthCheck) updates) =>
      super.copyWith((message) => updates(message as NoSuchHealthCheck))
          as NoSuchHealthCheck;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchHealthCheck create() => NoSuchHealthCheck._();
  @$core.override
  NoSuchHealthCheck createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchHealthCheck getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchHealthCheck>(create);
  static NoSuchHealthCheck? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchHostedZone extends $pb.GeneratedMessage {
  factory NoSuchHostedZone({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchHostedZone._();

  factory NoSuchHostedZone.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchHostedZone.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchHostedZone',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchHostedZone clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchHostedZone copyWith(void Function(NoSuchHostedZone) updates) =>
      super.copyWith((message) => updates(message as NoSuchHostedZone))
          as NoSuchHostedZone;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchHostedZone create() => NoSuchHostedZone._();
  @$core.override
  NoSuchHostedZone createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchHostedZone getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchHostedZone>(create);
  static NoSuchHostedZone? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchKeySigningKey extends $pb.GeneratedMessage {
  factory NoSuchKeySigningKey({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchKeySigningKey._();

  factory NoSuchKeySigningKey.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchKeySigningKey.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchKeySigningKey',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchKeySigningKey clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchKeySigningKey copyWith(void Function(NoSuchKeySigningKey) updates) =>
      super.copyWith((message) => updates(message as NoSuchKeySigningKey))
          as NoSuchKeySigningKey;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchKeySigningKey create() => NoSuchKeySigningKey._();
  @$core.override
  NoSuchKeySigningKey createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchKeySigningKey getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchKeySigningKey>(create);
  static NoSuchKeySigningKey? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchQueryLoggingConfig extends $pb.GeneratedMessage {
  factory NoSuchQueryLoggingConfig({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchQueryLoggingConfig._();

  factory NoSuchQueryLoggingConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchQueryLoggingConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchQueryLoggingConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchQueryLoggingConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchQueryLoggingConfig copyWith(
          void Function(NoSuchQueryLoggingConfig) updates) =>
      super.copyWith((message) => updates(message as NoSuchQueryLoggingConfig))
          as NoSuchQueryLoggingConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchQueryLoggingConfig create() => NoSuchQueryLoggingConfig._();
  @$core.override
  NoSuchQueryLoggingConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchQueryLoggingConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchQueryLoggingConfig>(create);
  static NoSuchQueryLoggingConfig? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchTrafficPolicy extends $pb.GeneratedMessage {
  factory NoSuchTrafficPolicy({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchTrafficPolicy._();

  factory NoSuchTrafficPolicy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchTrafficPolicy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchTrafficPolicy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchTrafficPolicy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchTrafficPolicy copyWith(void Function(NoSuchTrafficPolicy) updates) =>
      super.copyWith((message) => updates(message as NoSuchTrafficPolicy))
          as NoSuchTrafficPolicy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchTrafficPolicy create() => NoSuchTrafficPolicy._();
  @$core.override
  NoSuchTrafficPolicy createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchTrafficPolicy getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchTrafficPolicy>(create);
  static NoSuchTrafficPolicy? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class NoSuchTrafficPolicyInstance extends $pb.GeneratedMessage {
  factory NoSuchTrafficPolicyInstance({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  NoSuchTrafficPolicyInstance._();

  factory NoSuchTrafficPolicyInstance.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory NoSuchTrafficPolicyInstance.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'NoSuchTrafficPolicyInstance',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchTrafficPolicyInstance clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  NoSuchTrafficPolicyInstance copyWith(
          void Function(NoSuchTrafficPolicyInstance) updates) =>
      super.copyWith(
              (message) => updates(message as NoSuchTrafficPolicyInstance))
          as NoSuchTrafficPolicyInstance;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static NoSuchTrafficPolicyInstance create() =>
      NoSuchTrafficPolicyInstance._();
  @$core.override
  NoSuchTrafficPolicyInstance createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static NoSuchTrafficPolicyInstance getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<NoSuchTrafficPolicyInstance>(create);
  static NoSuchTrafficPolicyInstance? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
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

class PriorRequestNotComplete extends $pb.GeneratedMessage {
  factory PriorRequestNotComplete({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PriorRequestNotComplete._();

  factory PriorRequestNotComplete.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PriorRequestNotComplete.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PriorRequestNotComplete',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PriorRequestNotComplete clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PriorRequestNotComplete copyWith(
          void Function(PriorRequestNotComplete) updates) =>
      super.copyWith((message) => updates(message as PriorRequestNotComplete))
          as PriorRequestNotComplete;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PriorRequestNotComplete create() => PriorRequestNotComplete._();
  @$core.override
  PriorRequestNotComplete createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PriorRequestNotComplete getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PriorRequestNotComplete>(create);
  static PriorRequestNotComplete? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PublicZoneVPCAssociation extends $pb.GeneratedMessage {
  factory PublicZoneVPCAssociation({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PublicZoneVPCAssociation._();

  factory PublicZoneVPCAssociation.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PublicZoneVPCAssociation.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PublicZoneVPCAssociation',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublicZoneVPCAssociation clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PublicZoneVPCAssociation copyWith(
          void Function(PublicZoneVPCAssociation) updates) =>
      super.copyWith((message) => updates(message as PublicZoneVPCAssociation))
          as PublicZoneVPCAssociation;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublicZoneVPCAssociation create() => PublicZoneVPCAssociation._();
  @$core.override
  PublicZoneVPCAssociation createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PublicZoneVPCAssociation getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PublicZoneVPCAssociation>(create);
  static PublicZoneVPCAssociation? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class QueryLoggingConfig extends $pb.GeneratedMessage {
  factory QueryLoggingConfig({
    $core.String? hostedzoneid,
    $core.String? id,
    $core.String? cloudwatchlogsloggrouparn,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (id != null) result.id = id;
    if (cloudwatchlogsloggrouparn != null)
      result.cloudwatchlogsloggrouparn = cloudwatchlogsloggrouparn;
    return result;
  }

  QueryLoggingConfig._();

  factory QueryLoggingConfig.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryLoggingConfig.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryLoggingConfig',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(519613355, _omitFieldNames ? '' : 'cloudwatchlogsloggrouparn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLoggingConfig clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLoggingConfig copyWith(void Function(QueryLoggingConfig) updates) =>
      super.copyWith((message) => updates(message as QueryLoggingConfig))
          as QueryLoggingConfig;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryLoggingConfig create() => QueryLoggingConfig._();
  @$core.override
  QueryLoggingConfig createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryLoggingConfig getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryLoggingConfig>(create);
  static QueryLoggingConfig? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(519613355)
  $core.String get cloudwatchlogsloggrouparn => $_getSZ(2);
  @$pb.TagNumber(519613355)
  set cloudwatchlogsloggrouparn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(519613355)
  $core.bool hasCloudwatchlogsloggrouparn() => $_has(2);
  @$pb.TagNumber(519613355)
  void clearCloudwatchlogsloggrouparn() => $_clearField(519613355);
}

class QueryLoggingConfigAlreadyExists extends $pb.GeneratedMessage {
  factory QueryLoggingConfigAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  QueryLoggingConfigAlreadyExists._();

  factory QueryLoggingConfigAlreadyExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory QueryLoggingConfigAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'QueryLoggingConfigAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLoggingConfigAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLoggingConfigAlreadyExists copyWith(
          void Function(QueryLoggingConfigAlreadyExists) updates) =>
      super.copyWith(
              (message) => updates(message as QueryLoggingConfigAlreadyExists))
          as QueryLoggingConfigAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryLoggingConfigAlreadyExists create() =>
      QueryLoggingConfigAlreadyExists._();
  @$core.override
  QueryLoggingConfigAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static QueryLoggingConfigAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<QueryLoggingConfigAlreadyExists>(
          create);
  static QueryLoggingConfigAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ResourceRecord extends $pb.GeneratedMessage {
  factory ResourceRecord({
    $core.String? value,
  }) {
    final result = create();
    if (value != null) result.value = value;
    return result;
  }

  ResourceRecord._();

  factory ResourceRecord.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceRecord.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceRecord',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceRecord clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceRecord copyWith(void Function(ResourceRecord) updates) =>
      super.copyWith((message) => updates(message as ResourceRecord))
          as ResourceRecord;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceRecord create() => ResourceRecord._();
  @$core.override
  ResourceRecord createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceRecord getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceRecord>(create);
  static ResourceRecord? _defaultInstance;

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class ResourceRecordSet extends $pb.GeneratedMessage {
  factory ResourceRecordSet({
    ResourceRecordSetFailover? failover,
    GeoProximityLocation? geoproximitylocation,
    ResourceRecordSetRegion? region,
    $core.String? setidentifier,
    $core.String? trafficpolicyinstanceid,
    $core.String? name,
    GeoLocation? geolocation,
    RRType? type,
    $core.String? healthcheckid,
    AliasTarget? aliastarget,
    CidrRoutingConfig? cidrroutingconfig,
    $fixnum.Int64? weight,
    $core.bool? multivalueanswer,
    $core.Iterable<ResourceRecord>? resourcerecords,
    $fixnum.Int64? ttl,
  }) {
    final result = create();
    if (failover != null) result.failover = failover;
    if (geoproximitylocation != null)
      result.geoproximitylocation = geoproximitylocation;
    if (region != null) result.region = region;
    if (setidentifier != null) result.setidentifier = setidentifier;
    if (trafficpolicyinstanceid != null)
      result.trafficpolicyinstanceid = trafficpolicyinstanceid;
    if (name != null) result.name = name;
    if (geolocation != null) result.geolocation = geolocation;
    if (type != null) result.type = type;
    if (healthcheckid != null) result.healthcheckid = healthcheckid;
    if (aliastarget != null) result.aliastarget = aliastarget;
    if (cidrroutingconfig != null) result.cidrroutingconfig = cidrroutingconfig;
    if (weight != null) result.weight = weight;
    if (multivalueanswer != null) result.multivalueanswer = multivalueanswer;
    if (resourcerecords != null) result.resourcerecords.addAll(resourcerecords);
    if (ttl != null) result.ttl = ttl;
    return result;
  }

  ResourceRecordSet._();

  factory ResourceRecordSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceRecordSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceRecordSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<ResourceRecordSetFailover>(26793064, _omitFieldNames ? '' : 'failover',
        enumValues: ResourceRecordSetFailover.values)
    ..aOM<GeoProximityLocation>(
        94319785, _omitFieldNames ? '' : 'geoproximitylocation',
        subBuilder: GeoProximityLocation.create)
    ..aE<ResourceRecordSetRegion>(154040478, _omitFieldNames ? '' : 'region',
        enumValues: ResourceRecordSetRegion.values)
    ..aOS(201408985, _omitFieldNames ? '' : 'setidentifier')
    ..aOS(251421439, _omitFieldNames ? '' : 'trafficpolicyinstanceid')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<GeoLocation>(267973346, _omitFieldNames ? '' : 'geolocation',
        subBuilder: GeoLocation.create)
    ..aE<RRType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: RRType.values)
    ..aOS(312971637, _omitFieldNames ? '' : 'healthcheckid')
    ..aOM<AliasTarget>(317299867, _omitFieldNames ? '' : 'aliastarget',
        subBuilder: AliasTarget.create)
    ..aOM<CidrRoutingConfig>(
        357245246, _omitFieldNames ? '' : 'cidrroutingconfig',
        subBuilder: CidrRoutingConfig.create)
    ..aInt64(422581466, _omitFieldNames ? '' : 'weight')
    ..aOB(424105486, _omitFieldNames ? '' : 'multivalueanswer')
    ..pPM<ResourceRecord>(519418974, _omitFieldNames ? '' : 'resourcerecords',
        subBuilder: ResourceRecord.create)
    ..aInt64(526904700, _omitFieldNames ? '' : 'ttl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceRecordSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceRecordSet copyWith(void Function(ResourceRecordSet) updates) =>
      super.copyWith((message) => updates(message as ResourceRecordSet))
          as ResourceRecordSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceRecordSet create() => ResourceRecordSet._();
  @$core.override
  ResourceRecordSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceRecordSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceRecordSet>(create);
  static ResourceRecordSet? _defaultInstance;

  @$pb.TagNumber(26793064)
  ResourceRecordSetFailover get failover => $_getN(0);
  @$pb.TagNumber(26793064)
  set failover(ResourceRecordSetFailover value) => $_setField(26793064, value);
  @$pb.TagNumber(26793064)
  $core.bool hasFailover() => $_has(0);
  @$pb.TagNumber(26793064)
  void clearFailover() => $_clearField(26793064);

  @$pb.TagNumber(94319785)
  GeoProximityLocation get geoproximitylocation => $_getN(1);
  @$pb.TagNumber(94319785)
  set geoproximitylocation(GeoProximityLocation value) =>
      $_setField(94319785, value);
  @$pb.TagNumber(94319785)
  $core.bool hasGeoproximitylocation() => $_has(1);
  @$pb.TagNumber(94319785)
  void clearGeoproximitylocation() => $_clearField(94319785);
  @$pb.TagNumber(94319785)
  GeoProximityLocation ensureGeoproximitylocation() => $_ensure(1);

  @$pb.TagNumber(154040478)
  ResourceRecordSetRegion get region => $_getN(2);
  @$pb.TagNumber(154040478)
  set region(ResourceRecordSetRegion value) => $_setField(154040478, value);
  @$pb.TagNumber(154040478)
  $core.bool hasRegion() => $_has(2);
  @$pb.TagNumber(154040478)
  void clearRegion() => $_clearField(154040478);

  @$pb.TagNumber(201408985)
  $core.String get setidentifier => $_getSZ(3);
  @$pb.TagNumber(201408985)
  set setidentifier($core.String value) => $_setString(3, value);
  @$pb.TagNumber(201408985)
  $core.bool hasSetidentifier() => $_has(3);
  @$pb.TagNumber(201408985)
  void clearSetidentifier() => $_clearField(201408985);

  @$pb.TagNumber(251421439)
  $core.String get trafficpolicyinstanceid => $_getSZ(4);
  @$pb.TagNumber(251421439)
  set trafficpolicyinstanceid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(251421439)
  $core.bool hasTrafficpolicyinstanceid() => $_has(4);
  @$pb.TagNumber(251421439)
  void clearTrafficpolicyinstanceid() => $_clearField(251421439);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(267973346)
  GeoLocation get geolocation => $_getN(6);
  @$pb.TagNumber(267973346)
  set geolocation(GeoLocation value) => $_setField(267973346, value);
  @$pb.TagNumber(267973346)
  $core.bool hasGeolocation() => $_has(6);
  @$pb.TagNumber(267973346)
  void clearGeolocation() => $_clearField(267973346);
  @$pb.TagNumber(267973346)
  GeoLocation ensureGeolocation() => $_ensure(6);

  @$pb.TagNumber(290836590)
  RRType get type => $_getN(7);
  @$pb.TagNumber(290836590)
  set type(RRType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(7);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(312971637)
  $core.String get healthcheckid => $_getSZ(8);
  @$pb.TagNumber(312971637)
  set healthcheckid($core.String value) => $_setString(8, value);
  @$pb.TagNumber(312971637)
  $core.bool hasHealthcheckid() => $_has(8);
  @$pb.TagNumber(312971637)
  void clearHealthcheckid() => $_clearField(312971637);

  @$pb.TagNumber(317299867)
  AliasTarget get aliastarget => $_getN(9);
  @$pb.TagNumber(317299867)
  set aliastarget(AliasTarget value) => $_setField(317299867, value);
  @$pb.TagNumber(317299867)
  $core.bool hasAliastarget() => $_has(9);
  @$pb.TagNumber(317299867)
  void clearAliastarget() => $_clearField(317299867);
  @$pb.TagNumber(317299867)
  AliasTarget ensureAliastarget() => $_ensure(9);

  @$pb.TagNumber(357245246)
  CidrRoutingConfig get cidrroutingconfig => $_getN(10);
  @$pb.TagNumber(357245246)
  set cidrroutingconfig(CidrRoutingConfig value) =>
      $_setField(357245246, value);
  @$pb.TagNumber(357245246)
  $core.bool hasCidrroutingconfig() => $_has(10);
  @$pb.TagNumber(357245246)
  void clearCidrroutingconfig() => $_clearField(357245246);
  @$pb.TagNumber(357245246)
  CidrRoutingConfig ensureCidrroutingconfig() => $_ensure(10);

  @$pb.TagNumber(422581466)
  $fixnum.Int64 get weight => $_getI64(11);
  @$pb.TagNumber(422581466)
  set weight($fixnum.Int64 value) => $_setInt64(11, value);
  @$pb.TagNumber(422581466)
  $core.bool hasWeight() => $_has(11);
  @$pb.TagNumber(422581466)
  void clearWeight() => $_clearField(422581466);

  @$pb.TagNumber(424105486)
  $core.bool get multivalueanswer => $_getBF(12);
  @$pb.TagNumber(424105486)
  set multivalueanswer($core.bool value) => $_setBool(12, value);
  @$pb.TagNumber(424105486)
  $core.bool hasMultivalueanswer() => $_has(12);
  @$pb.TagNumber(424105486)
  void clearMultivalueanswer() => $_clearField(424105486);

  @$pb.TagNumber(519418974)
  $pb.PbList<ResourceRecord> get resourcerecords => $_getList(13);

  @$pb.TagNumber(526904700)
  $fixnum.Int64 get ttl => $_getI64(14);
  @$pb.TagNumber(526904700)
  set ttl($fixnum.Int64 value) => $_setInt64(14, value);
  @$pb.TagNumber(526904700)
  $core.bool hasTtl() => $_has(14);
  @$pb.TagNumber(526904700)
  void clearTtl() => $_clearField(526904700);
}

class ResourceTagSet extends $pb.GeneratedMessage {
  factory ResourceTagSet({
    TagResourceType? resourcetype,
    $core.Iterable<Tag>? tags,
    $core.String? resourceid,
  }) {
    final result = create();
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (tags != null) result.tags.addAll(tags);
    if (resourceid != null) result.resourceid = resourceid;
    return result;
  }

  ResourceTagSet._();

  factory ResourceTagSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceTagSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceTagSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aE<TagResourceType>(301342558, _omitFieldNames ? '' : 'resourcetype',
        enumValues: TagResourceType.values)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(526146833, _omitFieldNames ? '' : 'resourceid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceTagSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceTagSet copyWith(void Function(ResourceTagSet) updates) =>
      super.copyWith((message) => updates(message as ResourceTagSet))
          as ResourceTagSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceTagSet create() => ResourceTagSet._();
  @$core.override
  ResourceTagSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceTagSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceTagSet>(create);
  static ResourceTagSet? _defaultInstance;

  @$pb.TagNumber(301342558)
  TagResourceType get resourcetype => $_getN(0);
  @$pb.TagNumber(301342558)
  set resourcetype(TagResourceType value) => $_setField(301342558, value);
  @$pb.TagNumber(301342558)
  $core.bool hasResourcetype() => $_has(0);
  @$pb.TagNumber(301342558)
  void clearResourcetype() => $_clearField(301342558);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);

  @$pb.TagNumber(526146833)
  $core.String get resourceid => $_getSZ(2);
  @$pb.TagNumber(526146833)
  set resourceid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(526146833)
  $core.bool hasResourceid() => $_has(2);
  @$pb.TagNumber(526146833)
  void clearResourceid() => $_clearField(526146833);
}

class ReusableDelegationSetLimit extends $pb.GeneratedMessage {
  factory ReusableDelegationSetLimit({
    $fixnum.Int64? value,
    ReusableDelegationSetLimitType? type,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (type != null) result.type = type;
    return result;
  }

  ReusableDelegationSetLimit._();

  factory ReusableDelegationSetLimit.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReusableDelegationSetLimit.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReusableDelegationSetLimit',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aInt64(289929579, _omitFieldNames ? '' : 'value')
    ..aE<ReusableDelegationSetLimitType>(
        290836590, _omitFieldNames ? '' : 'type',
        enumValues: ReusableDelegationSetLimitType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReusableDelegationSetLimit clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReusableDelegationSetLimit copyWith(
          void Function(ReusableDelegationSetLimit) updates) =>
      super.copyWith(
              (message) => updates(message as ReusableDelegationSetLimit))
          as ReusableDelegationSetLimit;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReusableDelegationSetLimit create() => ReusableDelegationSetLimit._();
  @$core.override
  ReusableDelegationSetLimit createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReusableDelegationSetLimit getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReusableDelegationSetLimit>(create);
  static ReusableDelegationSetLimit? _defaultInstance;

  @$pb.TagNumber(289929579)
  $fixnum.Int64 get value => $_getI64(0);
  @$pb.TagNumber(289929579)
  set value($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(290836590)
  ReusableDelegationSetLimitType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(ReusableDelegationSetLimitType value) =>
      $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class StatusReport extends $pb.GeneratedMessage {
  factory StatusReport({
    $core.String? status,
    $core.String? checkedtime,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (checkedtime != null) result.checkedtime = checkedtime;
    return result;
  }

  StatusReport._();

  factory StatusReport.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StatusReport.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StatusReport',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(6222352, _omitFieldNames ? '' : 'status')
    ..aOS(152550560, _omitFieldNames ? '' : 'checkedtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StatusReport clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StatusReport copyWith(void Function(StatusReport) updates) =>
      super.copyWith((message) => updates(message as StatusReport))
          as StatusReport;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StatusReport create() => StatusReport._();
  @$core.override
  StatusReport createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StatusReport getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StatusReport>(create);
  static StatusReport? _defaultInstance;

  @$pb.TagNumber(6222352)
  $core.String get status => $_getSZ(0);
  @$pb.TagNumber(6222352)
  set status($core.String value) => $_setString(0, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(152550560)
  $core.String get checkedtime => $_getSZ(1);
  @$pb.TagNumber(152550560)
  set checkedtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(152550560)
  $core.bool hasCheckedtime() => $_has(1);
  @$pb.TagNumber(152550560)
  void clearCheckedtime() => $_clearField(152550560);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
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

class TestDNSAnswerRequest extends $pb.GeneratedMessage {
  factory TestDNSAnswerRequest({
    $core.String? resolverip,
    $core.String? recordname,
    $core.String? hostedzoneid,
    $core.String? edns0clientsubnetmask,
    RRType? recordtype,
    $core.String? edns0clientsubnetip,
  }) {
    final result = create();
    if (resolverip != null) result.resolverip = resolverip;
    if (recordname != null) result.recordname = recordname;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (edns0clientsubnetmask != null)
      result.edns0clientsubnetmask = edns0clientsubnetmask;
    if (recordtype != null) result.recordtype = recordtype;
    if (edns0clientsubnetip != null)
      result.edns0clientsubnetip = edns0clientsubnetip;
    return result;
  }

  TestDNSAnswerRequest._();

  factory TestDNSAnswerRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestDNSAnswerRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestDNSAnswerRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(95000907, _omitFieldNames ? '' : 'resolverip')
    ..aOS(204939016, _omitFieldNames ? '' : 'recordname')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(381322620, _omitFieldNames ? '' : 'edns0clientsubnetmask')
    ..aE<RRType>(441248261, _omitFieldNames ? '' : 'recordtype',
        enumValues: RRType.values)
    ..aOS(506730999, _omitFieldNames ? '' : 'edns0clientsubnetip')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestDNSAnswerRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestDNSAnswerRequest copyWith(void Function(TestDNSAnswerRequest) updates) =>
      super.copyWith((message) => updates(message as TestDNSAnswerRequest))
          as TestDNSAnswerRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestDNSAnswerRequest create() => TestDNSAnswerRequest._();
  @$core.override
  TestDNSAnswerRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestDNSAnswerRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestDNSAnswerRequest>(create);
  static TestDNSAnswerRequest? _defaultInstance;

  @$pb.TagNumber(95000907)
  $core.String get resolverip => $_getSZ(0);
  @$pb.TagNumber(95000907)
  set resolverip($core.String value) => $_setString(0, value);
  @$pb.TagNumber(95000907)
  $core.bool hasResolverip() => $_has(0);
  @$pb.TagNumber(95000907)
  void clearResolverip() => $_clearField(95000907);

  @$pb.TagNumber(204939016)
  $core.String get recordname => $_getSZ(1);
  @$pb.TagNumber(204939016)
  set recordname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(204939016)
  $core.bool hasRecordname() => $_has(1);
  @$pb.TagNumber(204939016)
  void clearRecordname() => $_clearField(204939016);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(2);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(2);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(381322620)
  $core.String get edns0clientsubnetmask => $_getSZ(3);
  @$pb.TagNumber(381322620)
  set edns0clientsubnetmask($core.String value) => $_setString(3, value);
  @$pb.TagNumber(381322620)
  $core.bool hasEdns0clientsubnetmask() => $_has(3);
  @$pb.TagNumber(381322620)
  void clearEdns0clientsubnetmask() => $_clearField(381322620);

  @$pb.TagNumber(441248261)
  RRType get recordtype => $_getN(4);
  @$pb.TagNumber(441248261)
  set recordtype(RRType value) => $_setField(441248261, value);
  @$pb.TagNumber(441248261)
  $core.bool hasRecordtype() => $_has(4);
  @$pb.TagNumber(441248261)
  void clearRecordtype() => $_clearField(441248261);

  @$pb.TagNumber(506730999)
  $core.String get edns0clientsubnetip => $_getSZ(5);
  @$pb.TagNumber(506730999)
  set edns0clientsubnetip($core.String value) => $_setString(5, value);
  @$pb.TagNumber(506730999)
  $core.bool hasEdns0clientsubnetip() => $_has(5);
  @$pb.TagNumber(506730999)
  void clearEdns0clientsubnetip() => $_clearField(506730999);
}

class TestDNSAnswerResponse extends $pb.GeneratedMessage {
  factory TestDNSAnswerResponse({
    $core.String? protocol,
    $core.String? recordname,
    RRType? recordtype,
    $core.String? responsecode,
    $core.String? nameserver,
    $core.Iterable<$core.String>? recorddata,
  }) {
    final result = create();
    if (protocol != null) result.protocol = protocol;
    if (recordname != null) result.recordname = recordname;
    if (recordtype != null) result.recordtype = recordtype;
    if (responsecode != null) result.responsecode = responsecode;
    if (nameserver != null) result.nameserver = nameserver;
    if (recorddata != null) result.recorddata.addAll(recorddata);
    return result;
  }

  TestDNSAnswerResponse._();

  factory TestDNSAnswerResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestDNSAnswerResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestDNSAnswerResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(173534166, _omitFieldNames ? '' : 'protocol')
    ..aOS(204939016, _omitFieldNames ? '' : 'recordname')
    ..aE<RRType>(441248261, _omitFieldNames ? '' : 'recordtype',
        enumValues: RRType.values)
    ..aOS(447553700, _omitFieldNames ? '' : 'responsecode')
    ..aOS(474180450, _omitFieldNames ? '' : 'nameserver')
    ..pPS(527821973, _omitFieldNames ? '' : 'recorddata')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestDNSAnswerResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestDNSAnswerResponse copyWith(
          void Function(TestDNSAnswerResponse) updates) =>
      super.copyWith((message) => updates(message as TestDNSAnswerResponse))
          as TestDNSAnswerResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestDNSAnswerResponse create() => TestDNSAnswerResponse._();
  @$core.override
  TestDNSAnswerResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestDNSAnswerResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestDNSAnswerResponse>(create);
  static TestDNSAnswerResponse? _defaultInstance;

  @$pb.TagNumber(173534166)
  $core.String get protocol => $_getSZ(0);
  @$pb.TagNumber(173534166)
  set protocol($core.String value) => $_setString(0, value);
  @$pb.TagNumber(173534166)
  $core.bool hasProtocol() => $_has(0);
  @$pb.TagNumber(173534166)
  void clearProtocol() => $_clearField(173534166);

  @$pb.TagNumber(204939016)
  $core.String get recordname => $_getSZ(1);
  @$pb.TagNumber(204939016)
  set recordname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(204939016)
  $core.bool hasRecordname() => $_has(1);
  @$pb.TagNumber(204939016)
  void clearRecordname() => $_clearField(204939016);

  @$pb.TagNumber(441248261)
  RRType get recordtype => $_getN(2);
  @$pb.TagNumber(441248261)
  set recordtype(RRType value) => $_setField(441248261, value);
  @$pb.TagNumber(441248261)
  $core.bool hasRecordtype() => $_has(2);
  @$pb.TagNumber(441248261)
  void clearRecordtype() => $_clearField(441248261);

  @$pb.TagNumber(447553700)
  $core.String get responsecode => $_getSZ(3);
  @$pb.TagNumber(447553700)
  set responsecode($core.String value) => $_setString(3, value);
  @$pb.TagNumber(447553700)
  $core.bool hasResponsecode() => $_has(3);
  @$pb.TagNumber(447553700)
  void clearResponsecode() => $_clearField(447553700);

  @$pb.TagNumber(474180450)
  $core.String get nameserver => $_getSZ(4);
  @$pb.TagNumber(474180450)
  set nameserver($core.String value) => $_setString(4, value);
  @$pb.TagNumber(474180450)
  $core.bool hasNameserver() => $_has(4);
  @$pb.TagNumber(474180450)
  void clearNameserver() => $_clearField(474180450);

  @$pb.TagNumber(527821973)
  $pb.PbList<$core.String> get recorddata => $_getList(5);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
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
}

class TooManyHealthChecks extends $pb.GeneratedMessage {
  factory TooManyHealthChecks({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyHealthChecks._();

  factory TooManyHealthChecks.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyHealthChecks.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyHealthChecks',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyHealthChecks clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyHealthChecks copyWith(void Function(TooManyHealthChecks) updates) =>
      super.copyWith((message) => updates(message as TooManyHealthChecks))
          as TooManyHealthChecks;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyHealthChecks create() => TooManyHealthChecks._();
  @$core.override
  TooManyHealthChecks createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyHealthChecks getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyHealthChecks>(create);
  static TooManyHealthChecks? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TooManyHostedZones extends $pb.GeneratedMessage {
  factory TooManyHostedZones({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyHostedZones._();

  factory TooManyHostedZones.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyHostedZones.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyHostedZones',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyHostedZones clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyHostedZones copyWith(void Function(TooManyHostedZones) updates) =>
      super.copyWith((message) => updates(message as TooManyHostedZones))
          as TooManyHostedZones;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyHostedZones create() => TooManyHostedZones._();
  @$core.override
  TooManyHostedZones createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyHostedZones getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyHostedZones>(create);
  static TooManyHostedZones? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TooManyKeySigningKeys extends $pb.GeneratedMessage {
  factory TooManyKeySigningKeys({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyKeySigningKeys._();

  factory TooManyKeySigningKeys.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyKeySigningKeys.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyKeySigningKeys',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyKeySigningKeys clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyKeySigningKeys copyWith(
          void Function(TooManyKeySigningKeys) updates) =>
      super.copyWith((message) => updates(message as TooManyKeySigningKeys))
          as TooManyKeySigningKeys;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyKeySigningKeys create() => TooManyKeySigningKeys._();
  @$core.override
  TooManyKeySigningKeys createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyKeySigningKeys getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyKeySigningKeys>(create);
  static TooManyKeySigningKeys? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TooManyTrafficPolicies extends $pb.GeneratedMessage {
  factory TooManyTrafficPolicies({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyTrafficPolicies._();

  factory TooManyTrafficPolicies.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyTrafficPolicies.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyTrafficPolicies',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTrafficPolicies clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTrafficPolicies copyWith(
          void Function(TooManyTrafficPolicies) updates) =>
      super.copyWith((message) => updates(message as TooManyTrafficPolicies))
          as TooManyTrafficPolicies;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyTrafficPolicies create() => TooManyTrafficPolicies._();
  @$core.override
  TooManyTrafficPolicies createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyTrafficPolicies getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyTrafficPolicies>(create);
  static TooManyTrafficPolicies? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TooManyTrafficPolicyInstances extends $pb.GeneratedMessage {
  factory TooManyTrafficPolicyInstances({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyTrafficPolicyInstances._();

  factory TooManyTrafficPolicyInstances.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyTrafficPolicyInstances.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyTrafficPolicyInstances',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTrafficPolicyInstances clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTrafficPolicyInstances copyWith(
          void Function(TooManyTrafficPolicyInstances) updates) =>
      super.copyWith(
              (message) => updates(message as TooManyTrafficPolicyInstances))
          as TooManyTrafficPolicyInstances;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyTrafficPolicyInstances create() =>
      TooManyTrafficPolicyInstances._();
  @$core.override
  TooManyTrafficPolicyInstances createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyTrafficPolicyInstances getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TooManyTrafficPolicyInstances>(create);
  static TooManyTrafficPolicyInstances? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TooManyTrafficPolicyVersionsForCurrentPolicy
    extends $pb.GeneratedMessage {
  factory TooManyTrafficPolicyVersionsForCurrentPolicy({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyTrafficPolicyVersionsForCurrentPolicy._();

  factory TooManyTrafficPolicyVersionsForCurrentPolicy.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyTrafficPolicyVersionsForCurrentPolicy.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyTrafficPolicyVersionsForCurrentPolicy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTrafficPolicyVersionsForCurrentPolicy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyTrafficPolicyVersionsForCurrentPolicy copyWith(
          void Function(TooManyTrafficPolicyVersionsForCurrentPolicy)
              updates) =>
      super.copyWith((message) =>
              updates(message as TooManyTrafficPolicyVersionsForCurrentPolicy))
          as TooManyTrafficPolicyVersionsForCurrentPolicy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyTrafficPolicyVersionsForCurrentPolicy create() =>
      TooManyTrafficPolicyVersionsForCurrentPolicy._();
  @$core.override
  TooManyTrafficPolicyVersionsForCurrentPolicy createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static TooManyTrafficPolicyVersionsForCurrentPolicy getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          TooManyTrafficPolicyVersionsForCurrentPolicy>(create);
  static TooManyTrafficPolicyVersionsForCurrentPolicy? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TooManyVPCAssociationAuthorizations extends $pb.GeneratedMessage {
  factory TooManyVPCAssociationAuthorizations({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TooManyVPCAssociationAuthorizations._();

  factory TooManyVPCAssociationAuthorizations.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TooManyVPCAssociationAuthorizations.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TooManyVPCAssociationAuthorizations',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyVPCAssociationAuthorizations clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TooManyVPCAssociationAuthorizations copyWith(
          void Function(TooManyVPCAssociationAuthorizations) updates) =>
      super.copyWith((message) =>
              updates(message as TooManyVPCAssociationAuthorizations))
          as TooManyVPCAssociationAuthorizations;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TooManyVPCAssociationAuthorizations create() =>
      TooManyVPCAssociationAuthorizations._();
  @$core.override
  TooManyVPCAssociationAuthorizations createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TooManyVPCAssociationAuthorizations getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          TooManyVPCAssociationAuthorizations>(create);
  static TooManyVPCAssociationAuthorizations? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TrafficPolicy extends $pb.GeneratedMessage {
  factory TrafficPolicy({
    $core.String? name,
    RRType? type,
    $core.String? id,
    $core.String? document,
    $core.String? comment,
    $core.int? version,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (id != null) result.id = id;
    if (document != null) result.document = document;
    if (comment != null) result.comment = comment;
    if (version != null) result.version = version;
    return result;
  }

  TrafficPolicy._();

  factory TrafficPolicy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrafficPolicy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrafficPolicy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<RRType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: RRType.values)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(407108341, _omitFieldNames ? '' : 'document')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..aI(500028728, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicy copyWith(void Function(TrafficPolicy) updates) =>
      super.copyWith((message) => updates(message as TrafficPolicy))
          as TrafficPolicy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrafficPolicy create() => TrafficPolicy._();
  @$core.override
  TrafficPolicy createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrafficPolicy getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrafficPolicy>(create);
  static TrafficPolicy? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  RRType get type => $_getN(1);
  @$pb.TagNumber(290836590)
  set type(RRType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(2);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(2, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(2);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(407108341)
  $core.String get document => $_getSZ(3);
  @$pb.TagNumber(407108341)
  set document($core.String value) => $_setString(3, value);
  @$pb.TagNumber(407108341)
  $core.bool hasDocument() => $_has(3);
  @$pb.TagNumber(407108341)
  void clearDocument() => $_clearField(407108341);

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(4);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(4, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(4);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);

  @$pb.TagNumber(500028728)
  $core.int get version => $_getIZ(5);
  @$pb.TagNumber(500028728)
  set version($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(5);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);
}

class TrafficPolicyAlreadyExists extends $pb.GeneratedMessage {
  factory TrafficPolicyAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TrafficPolicyAlreadyExists._();

  factory TrafficPolicyAlreadyExists.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrafficPolicyAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrafficPolicyAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicyAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicyAlreadyExists copyWith(
          void Function(TrafficPolicyAlreadyExists) updates) =>
      super.copyWith(
              (message) => updates(message as TrafficPolicyAlreadyExists))
          as TrafficPolicyAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrafficPolicyAlreadyExists create() => TrafficPolicyAlreadyExists._();
  @$core.override
  TrafficPolicyAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrafficPolicyAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrafficPolicyAlreadyExists>(create);
  static TrafficPolicyAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TrafficPolicyInUse extends $pb.GeneratedMessage {
  factory TrafficPolicyInUse({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TrafficPolicyInUse._();

  factory TrafficPolicyInUse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrafficPolicyInUse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrafficPolicyInUse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicyInUse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicyInUse copyWith(void Function(TrafficPolicyInUse) updates) =>
      super.copyWith((message) => updates(message as TrafficPolicyInUse))
          as TrafficPolicyInUse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrafficPolicyInUse create() => TrafficPolicyInUse._();
  @$core.override
  TrafficPolicyInUse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrafficPolicyInUse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrafficPolicyInUse>(create);
  static TrafficPolicyInUse? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TrafficPolicyInstance extends $pb.GeneratedMessage {
  factory TrafficPolicyInstance({
    $core.String? trafficpolicyid,
    RRType? trafficpolicytype,
    $core.String? message,
    $core.String? name,
    $core.String? hostedzoneid,
    $core.String? id,
    $core.int? trafficpolicyversion,
    $core.String? state,
    $fixnum.Int64? ttl,
  }) {
    final result = create();
    if (trafficpolicyid != null) result.trafficpolicyid = trafficpolicyid;
    if (trafficpolicytype != null) result.trafficpolicytype = trafficpolicytype;
    if (message != null) result.message = message;
    if (name != null) result.name = name;
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (id != null) result.id = id;
    if (trafficpolicyversion != null)
      result.trafficpolicyversion = trafficpolicyversion;
    if (state != null) result.state = state;
    if (ttl != null) result.ttl = ttl;
    return result;
  }

  TrafficPolicyInstance._();

  factory TrafficPolicyInstance.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrafficPolicyInstance.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrafficPolicyInstance',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(40235222, _omitFieldNames ? '' : 'trafficpolicyid')
    ..aE<RRType>(214345537, _omitFieldNames ? '' : 'trafficpolicytype',
        enumValues: RRType.values)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aI(479078485, _omitFieldNames ? '' : 'trafficpolicyversion')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..aInt64(526904700, _omitFieldNames ? '' : 'ttl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicyInstance clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicyInstance copyWith(
          void Function(TrafficPolicyInstance) updates) =>
      super.copyWith((message) => updates(message as TrafficPolicyInstance))
          as TrafficPolicyInstance;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrafficPolicyInstance create() => TrafficPolicyInstance._();
  @$core.override
  TrafficPolicyInstance createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrafficPolicyInstance getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrafficPolicyInstance>(create);
  static TrafficPolicyInstance? _defaultInstance;

  @$pb.TagNumber(40235222)
  $core.String get trafficpolicyid => $_getSZ(0);
  @$pb.TagNumber(40235222)
  set trafficpolicyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(40235222)
  $core.bool hasTrafficpolicyid() => $_has(0);
  @$pb.TagNumber(40235222)
  void clearTrafficpolicyid() => $_clearField(40235222);

  @$pb.TagNumber(214345537)
  RRType get trafficpolicytype => $_getN(1);
  @$pb.TagNumber(214345537)
  set trafficpolicytype(RRType value) => $_setField(214345537, value);
  @$pb.TagNumber(214345537)
  $core.bool hasTrafficpolicytype() => $_has(1);
  @$pb.TagNumber(214345537)
  void clearTrafficpolicytype() => $_clearField(214345537);

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(2);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(2, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(2);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(4);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(4);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(5);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(5, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(5);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(479078485)
  $core.int get trafficpolicyversion => $_getIZ(6);
  @$pb.TagNumber(479078485)
  set trafficpolicyversion($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(479078485)
  $core.bool hasTrafficpolicyversion() => $_has(6);
  @$pb.TagNumber(479078485)
  void clearTrafficpolicyversion() => $_clearField(479078485);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(7);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(7, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(7);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);

  @$pb.TagNumber(526904700)
  $fixnum.Int64 get ttl => $_getI64(8);
  @$pb.TagNumber(526904700)
  set ttl($fixnum.Int64 value) => $_setInt64(8, value);
  @$pb.TagNumber(526904700)
  $core.bool hasTtl() => $_has(8);
  @$pb.TagNumber(526904700)
  void clearTtl() => $_clearField(526904700);
}

class TrafficPolicyInstanceAlreadyExists extends $pb.GeneratedMessage {
  factory TrafficPolicyInstanceAlreadyExists({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  TrafficPolicyInstanceAlreadyExists._();

  factory TrafficPolicyInstanceAlreadyExists.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrafficPolicyInstanceAlreadyExists.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrafficPolicyInstanceAlreadyExists',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicyInstanceAlreadyExists clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicyInstanceAlreadyExists copyWith(
          void Function(TrafficPolicyInstanceAlreadyExists) updates) =>
      super.copyWith((message) =>
              updates(message as TrafficPolicyInstanceAlreadyExists))
          as TrafficPolicyInstanceAlreadyExists;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrafficPolicyInstanceAlreadyExists create() =>
      TrafficPolicyInstanceAlreadyExists._();
  @$core.override
  TrafficPolicyInstanceAlreadyExists createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrafficPolicyInstanceAlreadyExists getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrafficPolicyInstanceAlreadyExists>(
          create);
  static TrafficPolicyInstanceAlreadyExists? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class TrafficPolicySummary extends $pb.GeneratedMessage {
  factory TrafficPolicySummary({
    $core.int? trafficpolicycount,
    $core.String? name,
    RRType? type,
    $core.String? id,
    $core.int? latestversion,
  }) {
    final result = create();
    if (trafficpolicycount != null)
      result.trafficpolicycount = trafficpolicycount;
    if (name != null) result.name = name;
    if (type != null) result.type = type;
    if (id != null) result.id = id;
    if (latestversion != null) result.latestversion = latestversion;
    return result;
  }

  TrafficPolicySummary._();

  factory TrafficPolicySummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TrafficPolicySummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TrafficPolicySummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aI(157833448, _omitFieldNames ? '' : 'trafficpolicycount')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<RRType>(290836590, _omitFieldNames ? '' : 'type',
        enumValues: RRType.values)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aI(424864587, _omitFieldNames ? '' : 'latestversion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicySummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TrafficPolicySummary copyWith(void Function(TrafficPolicySummary) updates) =>
      super.copyWith((message) => updates(message as TrafficPolicySummary))
          as TrafficPolicySummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrafficPolicySummary create() => TrafficPolicySummary._();
  @$core.override
  TrafficPolicySummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TrafficPolicySummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TrafficPolicySummary>(create);
  static TrafficPolicySummary? _defaultInstance;

  @$pb.TagNumber(157833448)
  $core.int get trafficpolicycount => $_getIZ(0);
  @$pb.TagNumber(157833448)
  set trafficpolicycount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(157833448)
  $core.bool hasTrafficpolicycount() => $_has(0);
  @$pb.TagNumber(157833448)
  void clearTrafficpolicycount() => $_clearField(157833448);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(290836590)
  RRType get type => $_getN(2);
  @$pb.TagNumber(290836590)
  set type(RRType value) => $_setField(290836590, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(424864587)
  $core.int get latestversion => $_getIZ(4);
  @$pb.TagNumber(424864587)
  set latestversion($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(424864587)
  $core.bool hasLatestversion() => $_has(4);
  @$pb.TagNumber(424864587)
  void clearLatestversion() => $_clearField(424864587);
}

class UpdateHealthCheckRequest extends $pb.GeneratedMessage {
  factory UpdateHealthCheckRequest({
    $core.Iterable<ResettableElementName>? resetelements,
    $core.Iterable<HealthCheckRegion>? regions,
    $core.int? port,
    $core.bool? inverted,
    $core.bool? enablesni,
    $fixnum.Int64? healthcheckversion,
    $core.String? resourcepath,
    $core.String? ipaddress,
    $core.int? failurethreshold,
    $core.int? healththreshold,
    $core.String? healthcheckid,
    $core.String? searchstring,
    $core.String? fullyqualifieddomainname,
    $core.Iterable<$core.String>? childhealthchecks,
    InsufficientDataHealthStatus? insufficientdatahealthstatus,
    $core.bool? disabled,
    AlarmIdentifier? alarmidentifier,
  }) {
    final result = create();
    if (resetelements != null) result.resetelements.addAll(resetelements);
    if (regions != null) result.regions.addAll(regions);
    if (port != null) result.port = port;
    if (inverted != null) result.inverted = inverted;
    if (enablesni != null) result.enablesni = enablesni;
    if (healthcheckversion != null)
      result.healthcheckversion = healthcheckversion;
    if (resourcepath != null) result.resourcepath = resourcepath;
    if (ipaddress != null) result.ipaddress = ipaddress;
    if (failurethreshold != null) result.failurethreshold = failurethreshold;
    if (healththreshold != null) result.healththreshold = healththreshold;
    if (healthcheckid != null) result.healthcheckid = healthcheckid;
    if (searchstring != null) result.searchstring = searchstring;
    if (fullyqualifieddomainname != null)
      result.fullyqualifieddomainname = fullyqualifieddomainname;
    if (childhealthchecks != null)
      result.childhealthchecks.addAll(childhealthchecks);
    if (insufficientdatahealthstatus != null)
      result.insufficientdatahealthstatus = insufficientdatahealthstatus;
    if (disabled != null) result.disabled = disabled;
    if (alarmidentifier != null) result.alarmidentifier = alarmidentifier;
    return result;
  }

  UpdateHealthCheckRequest._();

  factory UpdateHealthCheckRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateHealthCheckRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateHealthCheckRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..pc<ResettableElementName>(
        16543458, _omitFieldNames ? '' : 'resetelements', $pb.PbFieldType.KE,
        valueOf: ResettableElementName.valueOf,
        enumValues: ResettableElementName.values,
        defaultEnumValue: ResettableElementName.RESETTABLE_ELEMENT_NAME_REGIONS)
    ..pc<HealthCheckRegion>(
        36200107, _omitFieldNames ? '' : 'regions', $pb.PbFieldType.KE,
        valueOf: HealthCheckRegion.valueOf,
        enumValues: HealthCheckRegion.values,
        defaultEnumValue: HealthCheckRegion.HEALTH_CHECK_REGION_AP_NORTHEAST_1)
    ..aI(46480583, _omitFieldNames ? '' : 'port')
    ..aOB(55175513, _omitFieldNames ? '' : 'inverted')
    ..aOB(70122887, _omitFieldNames ? '' : 'enablesni')
    ..aInt64(89568396, _omitFieldNames ? '' : 'healthcheckversion')
    ..aOS(117584551, _omitFieldNames ? '' : 'resourcepath')
    ..aOS(169333741, _omitFieldNames ? '' : 'ipaddress')
    ..aI(176846565, _omitFieldNames ? '' : 'failurethreshold')
    ..aI(215873163, _omitFieldNames ? '' : 'healththreshold')
    ..aOS(312971637, _omitFieldNames ? '' : 'healthcheckid')
    ..aOS(318687365, _omitFieldNames ? '' : 'searchstring')
    ..aOS(459321509, _omitFieldNames ? '' : 'fullyqualifieddomainname')
    ..pPS(485535935, _omitFieldNames ? '' : 'childhealthchecks')
    ..aE<InsufficientDataHealthStatus>(
        493115723, _omitFieldNames ? '' : 'insufficientdatahealthstatus',
        enumValues: InsufficientDataHealthStatus.values)
    ..aOB(533633318, _omitFieldNames ? '' : 'disabled')
    ..aOM<AlarmIdentifier>(536124346, _omitFieldNames ? '' : 'alarmidentifier',
        subBuilder: AlarmIdentifier.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHealthCheckRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHealthCheckRequest copyWith(
          void Function(UpdateHealthCheckRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateHealthCheckRequest))
          as UpdateHealthCheckRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateHealthCheckRequest create() => UpdateHealthCheckRequest._();
  @$core.override
  UpdateHealthCheckRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateHealthCheckRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateHealthCheckRequest>(create);
  static UpdateHealthCheckRequest? _defaultInstance;

  @$pb.TagNumber(16543458)
  $pb.PbList<ResettableElementName> get resetelements => $_getList(0);

  @$pb.TagNumber(36200107)
  $pb.PbList<HealthCheckRegion> get regions => $_getList(1);

  @$pb.TagNumber(46480583)
  $core.int get port => $_getIZ(2);
  @$pb.TagNumber(46480583)
  set port($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(46480583)
  $core.bool hasPort() => $_has(2);
  @$pb.TagNumber(46480583)
  void clearPort() => $_clearField(46480583);

  @$pb.TagNumber(55175513)
  $core.bool get inverted => $_getBF(3);
  @$pb.TagNumber(55175513)
  set inverted($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(55175513)
  $core.bool hasInverted() => $_has(3);
  @$pb.TagNumber(55175513)
  void clearInverted() => $_clearField(55175513);

  @$pb.TagNumber(70122887)
  $core.bool get enablesni => $_getBF(4);
  @$pb.TagNumber(70122887)
  set enablesni($core.bool value) => $_setBool(4, value);
  @$pb.TagNumber(70122887)
  $core.bool hasEnablesni() => $_has(4);
  @$pb.TagNumber(70122887)
  void clearEnablesni() => $_clearField(70122887);

  @$pb.TagNumber(89568396)
  $fixnum.Int64 get healthcheckversion => $_getI64(5);
  @$pb.TagNumber(89568396)
  set healthcheckversion($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(89568396)
  $core.bool hasHealthcheckversion() => $_has(5);
  @$pb.TagNumber(89568396)
  void clearHealthcheckversion() => $_clearField(89568396);

  @$pb.TagNumber(117584551)
  $core.String get resourcepath => $_getSZ(6);
  @$pb.TagNumber(117584551)
  set resourcepath($core.String value) => $_setString(6, value);
  @$pb.TagNumber(117584551)
  $core.bool hasResourcepath() => $_has(6);
  @$pb.TagNumber(117584551)
  void clearResourcepath() => $_clearField(117584551);

  @$pb.TagNumber(169333741)
  $core.String get ipaddress => $_getSZ(7);
  @$pb.TagNumber(169333741)
  set ipaddress($core.String value) => $_setString(7, value);
  @$pb.TagNumber(169333741)
  $core.bool hasIpaddress() => $_has(7);
  @$pb.TagNumber(169333741)
  void clearIpaddress() => $_clearField(169333741);

  @$pb.TagNumber(176846565)
  $core.int get failurethreshold => $_getIZ(8);
  @$pb.TagNumber(176846565)
  set failurethreshold($core.int value) => $_setSignedInt32(8, value);
  @$pb.TagNumber(176846565)
  $core.bool hasFailurethreshold() => $_has(8);
  @$pb.TagNumber(176846565)
  void clearFailurethreshold() => $_clearField(176846565);

  @$pb.TagNumber(215873163)
  $core.int get healththreshold => $_getIZ(9);
  @$pb.TagNumber(215873163)
  set healththreshold($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(215873163)
  $core.bool hasHealththreshold() => $_has(9);
  @$pb.TagNumber(215873163)
  void clearHealththreshold() => $_clearField(215873163);

  @$pb.TagNumber(312971637)
  $core.String get healthcheckid => $_getSZ(10);
  @$pb.TagNumber(312971637)
  set healthcheckid($core.String value) => $_setString(10, value);
  @$pb.TagNumber(312971637)
  $core.bool hasHealthcheckid() => $_has(10);
  @$pb.TagNumber(312971637)
  void clearHealthcheckid() => $_clearField(312971637);

  @$pb.TagNumber(318687365)
  $core.String get searchstring => $_getSZ(11);
  @$pb.TagNumber(318687365)
  set searchstring($core.String value) => $_setString(11, value);
  @$pb.TagNumber(318687365)
  $core.bool hasSearchstring() => $_has(11);
  @$pb.TagNumber(318687365)
  void clearSearchstring() => $_clearField(318687365);

  @$pb.TagNumber(459321509)
  $core.String get fullyqualifieddomainname => $_getSZ(12);
  @$pb.TagNumber(459321509)
  set fullyqualifieddomainname($core.String value) => $_setString(12, value);
  @$pb.TagNumber(459321509)
  $core.bool hasFullyqualifieddomainname() => $_has(12);
  @$pb.TagNumber(459321509)
  void clearFullyqualifieddomainname() => $_clearField(459321509);

  @$pb.TagNumber(485535935)
  $pb.PbList<$core.String> get childhealthchecks => $_getList(13);

  @$pb.TagNumber(493115723)
  InsufficientDataHealthStatus get insufficientdatahealthstatus => $_getN(14);
  @$pb.TagNumber(493115723)
  set insufficientdatahealthstatus(InsufficientDataHealthStatus value) =>
      $_setField(493115723, value);
  @$pb.TagNumber(493115723)
  $core.bool hasInsufficientdatahealthstatus() => $_has(14);
  @$pb.TagNumber(493115723)
  void clearInsufficientdatahealthstatus() => $_clearField(493115723);

  @$pb.TagNumber(533633318)
  $core.bool get disabled => $_getBF(15);
  @$pb.TagNumber(533633318)
  set disabled($core.bool value) => $_setBool(15, value);
  @$pb.TagNumber(533633318)
  $core.bool hasDisabled() => $_has(15);
  @$pb.TagNumber(533633318)
  void clearDisabled() => $_clearField(533633318);

  @$pb.TagNumber(536124346)
  AlarmIdentifier get alarmidentifier => $_getN(16);
  @$pb.TagNumber(536124346)
  set alarmidentifier(AlarmIdentifier value) => $_setField(536124346, value);
  @$pb.TagNumber(536124346)
  $core.bool hasAlarmidentifier() => $_has(16);
  @$pb.TagNumber(536124346)
  void clearAlarmidentifier() => $_clearField(536124346);
  @$pb.TagNumber(536124346)
  AlarmIdentifier ensureAlarmidentifier() => $_ensure(16);
}

class UpdateHealthCheckResponse extends $pb.GeneratedMessage {
  factory UpdateHealthCheckResponse({
    HealthCheck? healthcheck,
  }) {
    final result = create();
    if (healthcheck != null) result.healthcheck = healthcheck;
    return result;
  }

  UpdateHealthCheckResponse._();

  factory UpdateHealthCheckResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateHealthCheckResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateHealthCheckResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<HealthCheck>(377540610, _omitFieldNames ? '' : 'healthcheck',
        subBuilder: HealthCheck.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHealthCheckResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHealthCheckResponse copyWith(
          void Function(UpdateHealthCheckResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateHealthCheckResponse))
          as UpdateHealthCheckResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateHealthCheckResponse create() => UpdateHealthCheckResponse._();
  @$core.override
  UpdateHealthCheckResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateHealthCheckResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateHealthCheckResponse>(create);
  static UpdateHealthCheckResponse? _defaultInstance;

  @$pb.TagNumber(377540610)
  HealthCheck get healthcheck => $_getN(0);
  @$pb.TagNumber(377540610)
  set healthcheck(HealthCheck value) => $_setField(377540610, value);
  @$pb.TagNumber(377540610)
  $core.bool hasHealthcheck() => $_has(0);
  @$pb.TagNumber(377540610)
  void clearHealthcheck() => $_clearField(377540610);
  @$pb.TagNumber(377540610)
  HealthCheck ensureHealthcheck() => $_ensure(0);
}

class UpdateHostedZoneCommentRequest extends $pb.GeneratedMessage {
  factory UpdateHostedZoneCommentRequest({
    $core.String? id,
    $core.String? comment,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (comment != null) result.comment = comment;
    return result;
  }

  UpdateHostedZoneCommentRequest._();

  factory UpdateHostedZoneCommentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateHostedZoneCommentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateHostedZoneCommentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHostedZoneCommentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHostedZoneCommentRequest copyWith(
          void Function(UpdateHostedZoneCommentRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateHostedZoneCommentRequest))
          as UpdateHostedZoneCommentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateHostedZoneCommentRequest create() =>
      UpdateHostedZoneCommentRequest._();
  @$core.override
  UpdateHostedZoneCommentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateHostedZoneCommentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateHostedZoneCommentRequest>(create);
  static UpdateHostedZoneCommentRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(1);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(1, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(1);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);
}

class UpdateHostedZoneCommentResponse extends $pb.GeneratedMessage {
  factory UpdateHostedZoneCommentResponse({
    HostedZone? hostedzone,
  }) {
    final result = create();
    if (hostedzone != null) result.hostedzone = hostedzone;
    return result;
  }

  UpdateHostedZoneCommentResponse._();

  factory UpdateHostedZoneCommentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateHostedZoneCommentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateHostedZoneCommentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<HostedZone>(465689249, _omitFieldNames ? '' : 'hostedzone',
        subBuilder: HostedZone.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHostedZoneCommentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHostedZoneCommentResponse copyWith(
          void Function(UpdateHostedZoneCommentResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateHostedZoneCommentResponse))
          as UpdateHostedZoneCommentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateHostedZoneCommentResponse create() =>
      UpdateHostedZoneCommentResponse._();
  @$core.override
  UpdateHostedZoneCommentResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateHostedZoneCommentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateHostedZoneCommentResponse>(
          create);
  static UpdateHostedZoneCommentResponse? _defaultInstance;

  @$pb.TagNumber(465689249)
  HostedZone get hostedzone => $_getN(0);
  @$pb.TagNumber(465689249)
  set hostedzone(HostedZone value) => $_setField(465689249, value);
  @$pb.TagNumber(465689249)
  $core.bool hasHostedzone() => $_has(0);
  @$pb.TagNumber(465689249)
  void clearHostedzone() => $_clearField(465689249);
  @$pb.TagNumber(465689249)
  HostedZone ensureHostedzone() => $_ensure(0);
}

class UpdateHostedZoneFeaturesRequest extends $pb.GeneratedMessage {
  factory UpdateHostedZoneFeaturesRequest({
    $core.String? hostedzoneid,
    $core.bool? enableacceleratedrecovery,
  }) {
    final result = create();
    if (hostedzoneid != null) result.hostedzoneid = hostedzoneid;
    if (enableacceleratedrecovery != null)
      result.enableacceleratedrecovery = enableacceleratedrecovery;
    return result;
  }

  UpdateHostedZoneFeaturesRequest._();

  factory UpdateHostedZoneFeaturesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateHostedZoneFeaturesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateHostedZoneFeaturesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(346531710, _omitFieldNames ? '' : 'hostedzoneid')
    ..aOB(482177307, _omitFieldNames ? '' : 'enableacceleratedrecovery')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHostedZoneFeaturesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHostedZoneFeaturesRequest copyWith(
          void Function(UpdateHostedZoneFeaturesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateHostedZoneFeaturesRequest))
          as UpdateHostedZoneFeaturesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateHostedZoneFeaturesRequest create() =>
      UpdateHostedZoneFeaturesRequest._();
  @$core.override
  UpdateHostedZoneFeaturesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateHostedZoneFeaturesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateHostedZoneFeaturesRequest>(
          create);
  static UpdateHostedZoneFeaturesRequest? _defaultInstance;

  @$pb.TagNumber(346531710)
  $core.String get hostedzoneid => $_getSZ(0);
  @$pb.TagNumber(346531710)
  set hostedzoneid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(346531710)
  $core.bool hasHostedzoneid() => $_has(0);
  @$pb.TagNumber(346531710)
  void clearHostedzoneid() => $_clearField(346531710);

  @$pb.TagNumber(482177307)
  $core.bool get enableacceleratedrecovery => $_getBF(1);
  @$pb.TagNumber(482177307)
  set enableacceleratedrecovery($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(482177307)
  $core.bool hasEnableacceleratedrecovery() => $_has(1);
  @$pb.TagNumber(482177307)
  void clearEnableacceleratedrecovery() => $_clearField(482177307);
}

class UpdateHostedZoneFeaturesResponse extends $pb.GeneratedMessage {
  factory UpdateHostedZoneFeaturesResponse() => create();

  UpdateHostedZoneFeaturesResponse._();

  factory UpdateHostedZoneFeaturesResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateHostedZoneFeaturesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateHostedZoneFeaturesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHostedZoneFeaturesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateHostedZoneFeaturesResponse copyWith(
          void Function(UpdateHostedZoneFeaturesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateHostedZoneFeaturesResponse))
          as UpdateHostedZoneFeaturesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateHostedZoneFeaturesResponse create() =>
      UpdateHostedZoneFeaturesResponse._();
  @$core.override
  UpdateHostedZoneFeaturesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateHostedZoneFeaturesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateHostedZoneFeaturesResponse>(
          create);
  static UpdateHostedZoneFeaturesResponse? _defaultInstance;
}

class UpdateTrafficPolicyCommentRequest extends $pb.GeneratedMessage {
  factory UpdateTrafficPolicyCommentRequest({
    $core.String? id,
    $core.String? comment,
    $core.int? version,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (comment != null) result.comment = comment;
    if (version != null) result.version = version;
    return result;
  }

  UpdateTrafficPolicyCommentRequest._();

  factory UpdateTrafficPolicyCommentRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTrafficPolicyCommentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTrafficPolicyCommentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOS(407871487, _omitFieldNames ? '' : 'comment')
    ..aI(500028728, _omitFieldNames ? '' : 'version')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrafficPolicyCommentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrafficPolicyCommentRequest copyWith(
          void Function(UpdateTrafficPolicyCommentRequest) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateTrafficPolicyCommentRequest))
          as UpdateTrafficPolicyCommentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTrafficPolicyCommentRequest create() =>
      UpdateTrafficPolicyCommentRequest._();
  @$core.override
  UpdateTrafficPolicyCommentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTrafficPolicyCommentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTrafficPolicyCommentRequest>(
          create);
  static UpdateTrafficPolicyCommentRequest? _defaultInstance;

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(407871487)
  $core.String get comment => $_getSZ(1);
  @$pb.TagNumber(407871487)
  set comment($core.String value) => $_setString(1, value);
  @$pb.TagNumber(407871487)
  $core.bool hasComment() => $_has(1);
  @$pb.TagNumber(407871487)
  void clearComment() => $_clearField(407871487);

  @$pb.TagNumber(500028728)
  $core.int get version => $_getIZ(2);
  @$pb.TagNumber(500028728)
  set version($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(500028728)
  $core.bool hasVersion() => $_has(2);
  @$pb.TagNumber(500028728)
  void clearVersion() => $_clearField(500028728);
}

class UpdateTrafficPolicyCommentResponse extends $pb.GeneratedMessage {
  factory UpdateTrafficPolicyCommentResponse({
    TrafficPolicy? trafficpolicy,
  }) {
    final result = create();
    if (trafficpolicy != null) result.trafficpolicy = trafficpolicy;
    return result;
  }

  UpdateTrafficPolicyCommentResponse._();

  factory UpdateTrafficPolicyCommentResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTrafficPolicyCommentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTrafficPolicyCommentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<TrafficPolicy>(154595657, _omitFieldNames ? '' : 'trafficpolicy',
        subBuilder: TrafficPolicy.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrafficPolicyCommentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrafficPolicyCommentResponse copyWith(
          void Function(UpdateTrafficPolicyCommentResponse) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateTrafficPolicyCommentResponse))
          as UpdateTrafficPolicyCommentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTrafficPolicyCommentResponse create() =>
      UpdateTrafficPolicyCommentResponse._();
  @$core.override
  UpdateTrafficPolicyCommentResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTrafficPolicyCommentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTrafficPolicyCommentResponse>(
          create);
  static UpdateTrafficPolicyCommentResponse? _defaultInstance;

  @$pb.TagNumber(154595657)
  TrafficPolicy get trafficpolicy => $_getN(0);
  @$pb.TagNumber(154595657)
  set trafficpolicy(TrafficPolicy value) => $_setField(154595657, value);
  @$pb.TagNumber(154595657)
  $core.bool hasTrafficpolicy() => $_has(0);
  @$pb.TagNumber(154595657)
  void clearTrafficpolicy() => $_clearField(154595657);
  @$pb.TagNumber(154595657)
  TrafficPolicy ensureTrafficpolicy() => $_ensure(0);
}

class UpdateTrafficPolicyInstanceRequest extends $pb.GeneratedMessage {
  factory UpdateTrafficPolicyInstanceRequest({
    $core.String? trafficpolicyid,
    $core.String? id,
    $core.int? trafficpolicyversion,
    $fixnum.Int64? ttl,
  }) {
    final result = create();
    if (trafficpolicyid != null) result.trafficpolicyid = trafficpolicyid;
    if (id != null) result.id = id;
    if (trafficpolicyversion != null)
      result.trafficpolicyversion = trafficpolicyversion;
    if (ttl != null) result.ttl = ttl;
    return result;
  }

  UpdateTrafficPolicyInstanceRequest._();

  factory UpdateTrafficPolicyInstanceRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTrafficPolicyInstanceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTrafficPolicyInstanceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(40235222, _omitFieldNames ? '' : 'trafficpolicyid')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aI(479078485, _omitFieldNames ? '' : 'trafficpolicyversion')
    ..aInt64(526904700, _omitFieldNames ? '' : 'ttl')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrafficPolicyInstanceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrafficPolicyInstanceRequest copyWith(
          void Function(UpdateTrafficPolicyInstanceRequest) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateTrafficPolicyInstanceRequest))
          as UpdateTrafficPolicyInstanceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTrafficPolicyInstanceRequest create() =>
      UpdateTrafficPolicyInstanceRequest._();
  @$core.override
  UpdateTrafficPolicyInstanceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTrafficPolicyInstanceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateTrafficPolicyInstanceRequest>(
          create);
  static UpdateTrafficPolicyInstanceRequest? _defaultInstance;

  @$pb.TagNumber(40235222)
  $core.String get trafficpolicyid => $_getSZ(0);
  @$pb.TagNumber(40235222)
  set trafficpolicyid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(40235222)
  $core.bool hasTrafficpolicyid() => $_has(0);
  @$pb.TagNumber(40235222)
  void clearTrafficpolicyid() => $_clearField(40235222);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(479078485)
  $core.int get trafficpolicyversion => $_getIZ(2);
  @$pb.TagNumber(479078485)
  set trafficpolicyversion($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(479078485)
  $core.bool hasTrafficpolicyversion() => $_has(2);
  @$pb.TagNumber(479078485)
  void clearTrafficpolicyversion() => $_clearField(479078485);

  @$pb.TagNumber(526904700)
  $fixnum.Int64 get ttl => $_getI64(3);
  @$pb.TagNumber(526904700)
  set ttl($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(526904700)
  $core.bool hasTtl() => $_has(3);
  @$pb.TagNumber(526904700)
  void clearTtl() => $_clearField(526904700);
}

class UpdateTrafficPolicyInstanceResponse extends $pb.GeneratedMessage {
  factory UpdateTrafficPolicyInstanceResponse({
    TrafficPolicyInstance? trafficpolicyinstance,
  }) {
    final result = create();
    if (trafficpolicyinstance != null)
      result.trafficpolicyinstance = trafficpolicyinstance;
    return result;
  }

  UpdateTrafficPolicyInstanceResponse._();

  factory UpdateTrafficPolicyInstanceResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateTrafficPolicyInstanceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateTrafficPolicyInstanceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOM<TrafficPolicyInstance>(
        205651476, _omitFieldNames ? '' : 'trafficpolicyinstance',
        subBuilder: TrafficPolicyInstance.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrafficPolicyInstanceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateTrafficPolicyInstanceResponse copyWith(
          void Function(UpdateTrafficPolicyInstanceResponse) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateTrafficPolicyInstanceResponse))
          as UpdateTrafficPolicyInstanceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateTrafficPolicyInstanceResponse create() =>
      UpdateTrafficPolicyInstanceResponse._();
  @$core.override
  UpdateTrafficPolicyInstanceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateTrafficPolicyInstanceResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateTrafficPolicyInstanceResponse>(create);
  static UpdateTrafficPolicyInstanceResponse? _defaultInstance;

  @$pb.TagNumber(205651476)
  TrafficPolicyInstance get trafficpolicyinstance => $_getN(0);
  @$pb.TagNumber(205651476)
  set trafficpolicyinstance(TrafficPolicyInstance value) =>
      $_setField(205651476, value);
  @$pb.TagNumber(205651476)
  $core.bool hasTrafficpolicyinstance() => $_has(0);
  @$pb.TagNumber(205651476)
  void clearTrafficpolicyinstance() => $_clearField(205651476);
  @$pb.TagNumber(205651476)
  TrafficPolicyInstance ensureTrafficpolicyinstance() => $_ensure(0);
}

class VPC extends $pb.GeneratedMessage {
  factory VPC({
    $core.String? vpcid,
    VPCRegion? vpcregion,
  }) {
    final result = create();
    if (vpcid != null) result.vpcid = vpcid;
    if (vpcregion != null) result.vpcregion = vpcregion;
    return result;
  }

  VPC._();

  factory VPC.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VPC.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VPC',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(325135798, _omitFieldNames ? '' : 'vpcid')
    ..aE<VPCRegion>(474180765, _omitFieldNames ? '' : 'vpcregion',
        enumValues: VPCRegion.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VPC clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VPC copyWith(void Function(VPC) updates) =>
      super.copyWith((message) => updates(message as VPC)) as VPC;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VPC create() => VPC._();
  @$core.override
  VPC createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VPC getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VPC>(create);
  static VPC? _defaultInstance;

  @$pb.TagNumber(325135798)
  $core.String get vpcid => $_getSZ(0);
  @$pb.TagNumber(325135798)
  set vpcid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(325135798)
  $core.bool hasVpcid() => $_has(0);
  @$pb.TagNumber(325135798)
  void clearVpcid() => $_clearField(325135798);

  @$pb.TagNumber(474180765)
  VPCRegion get vpcregion => $_getN(1);
  @$pb.TagNumber(474180765)
  set vpcregion(VPCRegion value) => $_setField(474180765, value);
  @$pb.TagNumber(474180765)
  $core.bool hasVpcregion() => $_has(1);
  @$pb.TagNumber(474180765)
  void clearVpcregion() => $_clearField(474180765);
}

class VPCAssociationAuthorizationNotFound extends $pb.GeneratedMessage {
  factory VPCAssociationAuthorizationNotFound({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  VPCAssociationAuthorizationNotFound._();

  factory VPCAssociationAuthorizationNotFound.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VPCAssociationAuthorizationNotFound.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VPCAssociationAuthorizationNotFound',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VPCAssociationAuthorizationNotFound clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VPCAssociationAuthorizationNotFound copyWith(
          void Function(VPCAssociationAuthorizationNotFound) updates) =>
      super.copyWith((message) =>
              updates(message as VPCAssociationAuthorizationNotFound))
          as VPCAssociationAuthorizationNotFound;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VPCAssociationAuthorizationNotFound create() =>
      VPCAssociationAuthorizationNotFound._();
  @$core.override
  VPCAssociationAuthorizationNotFound createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VPCAssociationAuthorizationNotFound getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          VPCAssociationAuthorizationNotFound>(create);
  static VPCAssociationAuthorizationNotFound? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class VPCAssociationNotFound extends $pb.GeneratedMessage {
  factory VPCAssociationNotFound({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  VPCAssociationNotFound._();

  factory VPCAssociationNotFound.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VPCAssociationNotFound.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VPCAssociationNotFound',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'route53'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VPCAssociationNotFound clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VPCAssociationNotFound copyWith(
          void Function(VPCAssociationNotFound) updates) =>
      super.copyWith((message) => updates(message as VPCAssociationNotFound))
          as VPCAssociationNotFound;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VPCAssociationNotFound create() => VPCAssociationNotFound._();
  @$core.override
  VPCAssociationNotFound createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VPCAssociationNotFound getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<VPCAssociationNotFound>(create);
  static VPCAssociationNotFound? _defaultInstance;

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
