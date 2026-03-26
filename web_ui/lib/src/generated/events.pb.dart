// This is a generated file - do not edit.
//
// Generated from events.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'events.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'events.pbenum.dart';

class ActivateEventSourceRequest extends $pb.GeneratedMessage {
  factory ActivateEventSourceRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  ActivateEventSourceRequest._();

  factory ActivateEventSourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ActivateEventSourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ActivateEventSourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivateEventSourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ActivateEventSourceRequest copyWith(
          void Function(ActivateEventSourceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ActivateEventSourceRequest))
          as ActivateEventSourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ActivateEventSourceRequest create() => ActivateEventSourceRequest._();
  @$core.override
  ActivateEventSourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ActivateEventSourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ActivateEventSourceRequest>(create);
  static ActivateEventSourceRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class ApiDestination extends $pb.GeneratedMessage {
  factory ApiDestination({
    ApiDestinationState? apidestinationstate,
    $core.String? apidestinationarn,
    $core.String? creationtime,
    $core.String? connectionarn,
    $core.String? lastmodifiedtime,
    $core.String? name,
    $core.int? invocationratelimitpersecond,
    ApiDestinationHttpMethod? httpmethod,
    $core.String? invocationendpoint,
  }) {
    final result = create();
    if (apidestinationstate != null)
      result.apidestinationstate = apidestinationstate;
    if (apidestinationarn != null) result.apidestinationarn = apidestinationarn;
    if (creationtime != null) result.creationtime = creationtime;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (name != null) result.name = name;
    if (invocationratelimitpersecond != null)
      result.invocationratelimitpersecond = invocationratelimitpersecond;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (invocationendpoint != null)
      result.invocationendpoint = invocationendpoint;
    return result;
  }

  ApiDestination._();

  factory ApiDestination.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ApiDestination.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ApiDestination',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aE<ApiDestinationState>(
        13153343, _omitFieldNames ? '' : 'apidestinationstate',
        enumValues: ApiDestinationState.values)
    ..aOS(90996885, _omitFieldNames ? '' : 'apidestinationarn')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aI(295327816, _omitFieldNames ? '' : 'invocationratelimitpersecond')
    ..aE<ApiDestinationHttpMethod>(
        398394961, _omitFieldNames ? '' : 'httpmethod',
        enumValues: ApiDestinationHttpMethod.values)
    ..aOS(411800759, _omitFieldNames ? '' : 'invocationendpoint')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiDestination clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ApiDestination copyWith(void Function(ApiDestination) updates) =>
      super.copyWith((message) => updates(message as ApiDestination))
          as ApiDestination;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ApiDestination create() => ApiDestination._();
  @$core.override
  ApiDestination createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ApiDestination getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ApiDestination>(create);
  static ApiDestination? _defaultInstance;

  @$pb.TagNumber(13153343)
  ApiDestinationState get apidestinationstate => $_getN(0);
  @$pb.TagNumber(13153343)
  set apidestinationstate(ApiDestinationState value) =>
      $_setField(13153343, value);
  @$pb.TagNumber(13153343)
  $core.bool hasApidestinationstate() => $_has(0);
  @$pb.TagNumber(13153343)
  void clearApidestinationstate() => $_clearField(13153343);

  @$pb.TagNumber(90996885)
  $core.String get apidestinationarn => $_getSZ(1);
  @$pb.TagNumber(90996885)
  set apidestinationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(90996885)
  $core.bool hasApidestinationarn() => $_has(1);
  @$pb.TagNumber(90996885)
  void clearApidestinationarn() => $_clearField(90996885);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(3);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(3);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(4);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(4, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(4);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(5);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(5, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(5);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(295327816)
  $core.int get invocationratelimitpersecond => $_getIZ(6);
  @$pb.TagNumber(295327816)
  set invocationratelimitpersecond($core.int value) =>
      $_setSignedInt32(6, value);
  @$pb.TagNumber(295327816)
  $core.bool hasInvocationratelimitpersecond() => $_has(6);
  @$pb.TagNumber(295327816)
  void clearInvocationratelimitpersecond() => $_clearField(295327816);

  @$pb.TagNumber(398394961)
  ApiDestinationHttpMethod get httpmethod => $_getN(7);
  @$pb.TagNumber(398394961)
  set httpmethod(ApiDestinationHttpMethod value) =>
      $_setField(398394961, value);
  @$pb.TagNumber(398394961)
  $core.bool hasHttpmethod() => $_has(7);
  @$pb.TagNumber(398394961)
  void clearHttpmethod() => $_clearField(398394961);

  @$pb.TagNumber(411800759)
  $core.String get invocationendpoint => $_getSZ(8);
  @$pb.TagNumber(411800759)
  set invocationendpoint($core.String value) => $_setString(8, value);
  @$pb.TagNumber(411800759)
  $core.bool hasInvocationendpoint() => $_has(8);
  @$pb.TagNumber(411800759)
  void clearInvocationendpoint() => $_clearField(411800759);
}

class Archive extends $pb.GeneratedMessage {
  factory Archive({
    $core.String? archivename,
    $core.String? creationtime,
    $fixnum.Int64? eventcount,
    $core.int? retentiondays,
    $core.String? eventsourcearn,
    $core.String? statereason,
    $fixnum.Int64? sizebytes,
    ArchiveState? state,
  }) {
    final result = create();
    if (archivename != null) result.archivename = archivename;
    if (creationtime != null) result.creationtime = creationtime;
    if (eventcount != null) result.eventcount = eventcount;
    if (retentiondays != null) result.retentiondays = retentiondays;
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    if (statereason != null) result.statereason = statereason;
    if (sizebytes != null) result.sizebytes = sizebytes;
    if (state != null) result.state = state;
    return result;
  }

  Archive._();

  factory Archive.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Archive.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Archive',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(88048487, _omitFieldNames ? '' : 'archivename')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aInt64(128612839, _omitFieldNames ? '' : 'eventcount')
    ..aI(267894223, _omitFieldNames ? '' : 'retentiondays')
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aInt64(486677664, _omitFieldNames ? '' : 'sizebytes')
    ..aE<ArchiveState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ArchiveState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Archive clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Archive copyWith(void Function(Archive) updates) =>
      super.copyWith((message) => updates(message as Archive)) as Archive;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Archive create() => Archive._();
  @$core.override
  Archive createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Archive getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Archive>(create);
  static Archive? _defaultInstance;

  @$pb.TagNumber(88048487)
  $core.String get archivename => $_getSZ(0);
  @$pb.TagNumber(88048487)
  set archivename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88048487)
  $core.bool hasArchivename() => $_has(0);
  @$pb.TagNumber(88048487)
  void clearArchivename() => $_clearField(88048487);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(128612839)
  $fixnum.Int64 get eventcount => $_getI64(2);
  @$pb.TagNumber(128612839)
  set eventcount($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(128612839)
  $core.bool hasEventcount() => $_has(2);
  @$pb.TagNumber(128612839)
  void clearEventcount() => $_clearField(128612839);

  @$pb.TagNumber(267894223)
  $core.int get retentiondays => $_getIZ(3);
  @$pb.TagNumber(267894223)
  set retentiondays($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(267894223)
  $core.bool hasRetentiondays() => $_has(3);
  @$pb.TagNumber(267894223)
  void clearRetentiondays() => $_clearField(267894223);

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(4);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(4);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(5);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(5, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(5);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(486677664)
  $fixnum.Int64 get sizebytes => $_getI64(6);
  @$pb.TagNumber(486677664)
  set sizebytes($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(486677664)
  $core.bool hasSizebytes() => $_has(6);
  @$pb.TagNumber(486677664)
  void clearSizebytes() => $_clearField(486677664);

  @$pb.TagNumber(502047895)
  ArchiveState get state => $_getN(7);
  @$pb.TagNumber(502047895)
  set state(ArchiveState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(7);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class AwsVpcConfiguration extends $pb.GeneratedMessage {
  factory AwsVpcConfiguration({
    $core.Iterable<$core.String>? subnets,
    AssignPublicIp? assignpublicip,
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPS(414921506, _omitFieldNames ? '' : 'subnets')
    ..aE<AssignPublicIp>(461653589, _omitFieldNames ? '' : 'assignpublicip',
        enumValues: AssignPublicIp.values)
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
  AssignPublicIp get assignpublicip => $_getN(1);
  @$pb.TagNumber(461653589)
  set assignpublicip(AssignPublicIp value) => $_setField(461653589, value);
  @$pb.TagNumber(461653589)
  $core.bool hasAssignpublicip() => $_has(1);
  @$pb.TagNumber(461653589)
  void clearAssignpublicip() => $_clearField(461653589);

  @$pb.TagNumber(515282516)
  $pb.PbList<$core.String> get securitygroups => $_getList(2);
}

class BatchArrayProperties extends $pb.GeneratedMessage {
  factory BatchArrayProperties({
    $core.int? size,
  }) {
    final result = create();
    if (size != null) result.size = size;
    return result;
  }

  BatchArrayProperties._();

  factory BatchArrayProperties.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchArrayProperties.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchArrayProperties',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aI(105352829, _omitFieldNames ? '' : 'size')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchArrayProperties clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchArrayProperties copyWith(void Function(BatchArrayProperties) updates) =>
      super.copyWith((message) => updates(message as BatchArrayProperties))
          as BatchArrayProperties;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchArrayProperties create() => BatchArrayProperties._();
  @$core.override
  BatchArrayProperties createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchArrayProperties getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchArrayProperties>(create);
  static BatchArrayProperties? _defaultInstance;

  @$pb.TagNumber(105352829)
  $core.int get size => $_getIZ(0);
  @$pb.TagNumber(105352829)
  set size($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(105352829)
  $core.bool hasSize() => $_has(0);
  @$pb.TagNumber(105352829)
  void clearSize() => $_clearField(105352829);
}

class BatchParameters extends $pb.GeneratedMessage {
  factory BatchParameters({
    BatchRetryStrategy? retrystrategy,
    BatchArrayProperties? arrayproperties,
    $core.String? jobdefinition,
    $core.String? jobname,
  }) {
    final result = create();
    if (retrystrategy != null) result.retrystrategy = retrystrategy;
    if (arrayproperties != null) result.arrayproperties = arrayproperties;
    if (jobdefinition != null) result.jobdefinition = jobdefinition;
    if (jobname != null) result.jobname = jobname;
    return result;
  }

  BatchParameters._();

  factory BatchParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOM<BatchRetryStrategy>(105516073, _omitFieldNames ? '' : 'retrystrategy',
        subBuilder: BatchRetryStrategy.create)
    ..aOM<BatchArrayProperties>(
        230899444, _omitFieldNames ? '' : 'arrayproperties',
        subBuilder: BatchArrayProperties.create)
    ..aOS(420132006, _omitFieldNames ? '' : 'jobdefinition')
    ..aOS(498531160, _omitFieldNames ? '' : 'jobname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchParameters copyWith(void Function(BatchParameters) updates) =>
      super.copyWith((message) => updates(message as BatchParameters))
          as BatchParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchParameters create() => BatchParameters._();
  @$core.override
  BatchParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchParameters>(create);
  static BatchParameters? _defaultInstance;

  @$pb.TagNumber(105516073)
  BatchRetryStrategy get retrystrategy => $_getN(0);
  @$pb.TagNumber(105516073)
  set retrystrategy(BatchRetryStrategy value) => $_setField(105516073, value);
  @$pb.TagNumber(105516073)
  $core.bool hasRetrystrategy() => $_has(0);
  @$pb.TagNumber(105516073)
  void clearRetrystrategy() => $_clearField(105516073);
  @$pb.TagNumber(105516073)
  BatchRetryStrategy ensureRetrystrategy() => $_ensure(0);

  @$pb.TagNumber(230899444)
  BatchArrayProperties get arrayproperties => $_getN(1);
  @$pb.TagNumber(230899444)
  set arrayproperties(BatchArrayProperties value) =>
      $_setField(230899444, value);
  @$pb.TagNumber(230899444)
  $core.bool hasArrayproperties() => $_has(1);
  @$pb.TagNumber(230899444)
  void clearArrayproperties() => $_clearField(230899444);
  @$pb.TagNumber(230899444)
  BatchArrayProperties ensureArrayproperties() => $_ensure(1);

  @$pb.TagNumber(420132006)
  $core.String get jobdefinition => $_getSZ(2);
  @$pb.TagNumber(420132006)
  set jobdefinition($core.String value) => $_setString(2, value);
  @$pb.TagNumber(420132006)
  $core.bool hasJobdefinition() => $_has(2);
  @$pb.TagNumber(420132006)
  void clearJobdefinition() => $_clearField(420132006);

  @$pb.TagNumber(498531160)
  $core.String get jobname => $_getSZ(3);
  @$pb.TagNumber(498531160)
  set jobname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(498531160)
  $core.bool hasJobname() => $_has(3);
  @$pb.TagNumber(498531160)
  void clearJobname() => $_clearField(498531160);
}

class BatchRetryStrategy extends $pb.GeneratedMessage {
  factory BatchRetryStrategy({
    $core.int? attempts,
  }) {
    final result = create();
    if (attempts != null) result.attempts = attempts;
    return result;
  }

  BatchRetryStrategy._();

  factory BatchRetryStrategy.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BatchRetryStrategy.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BatchRetryStrategy',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aI(132533640, _omitFieldNames ? '' : 'attempts')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchRetryStrategy clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BatchRetryStrategy copyWith(void Function(BatchRetryStrategy) updates) =>
      super.copyWith((message) => updates(message as BatchRetryStrategy))
          as BatchRetryStrategy;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BatchRetryStrategy create() => BatchRetryStrategy._();
  @$core.override
  BatchRetryStrategy createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BatchRetryStrategy getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BatchRetryStrategy>(create);
  static BatchRetryStrategy? _defaultInstance;

  @$pb.TagNumber(132533640)
  $core.int get attempts => $_getIZ(0);
  @$pb.TagNumber(132533640)
  set attempts($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(132533640)
  $core.bool hasAttempts() => $_has(0);
  @$pb.TagNumber(132533640)
  void clearAttempts() => $_clearField(132533640);
}

class CancelReplayRequest extends $pb.GeneratedMessage {
  factory CancelReplayRequest({
    $core.String? replayname,
  }) {
    final result = create();
    if (replayname != null) result.replayname = replayname;
    return result;
  }

  CancelReplayRequest._();

  factory CancelReplayRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelReplayRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelReplayRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(442173850, _omitFieldNames ? '' : 'replayname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelReplayRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelReplayRequest copyWith(void Function(CancelReplayRequest) updates) =>
      super.copyWith((message) => updates(message as CancelReplayRequest))
          as CancelReplayRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelReplayRequest create() => CancelReplayRequest._();
  @$core.override
  CancelReplayRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelReplayRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelReplayRequest>(create);
  static CancelReplayRequest? _defaultInstance;

  @$pb.TagNumber(442173850)
  $core.String get replayname => $_getSZ(0);
  @$pb.TagNumber(442173850)
  set replayname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(442173850)
  $core.bool hasReplayname() => $_has(0);
  @$pb.TagNumber(442173850)
  void clearReplayname() => $_clearField(442173850);
}

class CancelReplayResponse extends $pb.GeneratedMessage {
  factory CancelReplayResponse({
    $core.String? replayarn,
    $core.String? statereason,
    ReplayState? state,
  }) {
    final result = create();
    if (replayarn != null) result.replayarn = replayarn;
    if (statereason != null) result.statereason = statereason;
    if (state != null) result.state = state;
    return result;
  }

  CancelReplayResponse._();

  factory CancelReplayResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CancelReplayResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CancelReplayResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(361946526, _omitFieldNames ? '' : 'replayarn')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aE<ReplayState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ReplayState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelReplayResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CancelReplayResponse copyWith(void Function(CancelReplayResponse) updates) =>
      super.copyWith((message) => updates(message as CancelReplayResponse))
          as CancelReplayResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CancelReplayResponse create() => CancelReplayResponse._();
  @$core.override
  CancelReplayResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CancelReplayResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CancelReplayResponse>(create);
  static CancelReplayResponse? _defaultInstance;

  @$pb.TagNumber(361946526)
  $core.String get replayarn => $_getSZ(0);
  @$pb.TagNumber(361946526)
  set replayarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(361946526)
  $core.bool hasReplayarn() => $_has(0);
  @$pb.TagNumber(361946526)
  void clearReplayarn() => $_clearField(361946526);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(1);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(1, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(1);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(502047895)
  ReplayState get state => $_getN(2);
  @$pb.TagNumber(502047895)
  set state(ReplayState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(2);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class Condition extends $pb.GeneratedMessage {
  factory Condition({
    $core.String? key,
    $core.String? value,
    $core.String? type,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    if (type != null) result.type = type;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..aOS(290836590, _omitFieldNames ? '' : 'type')
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

  @$pb.TagNumber(290836590)
  $core.String get type => $_getSZ(2);
  @$pb.TagNumber(290836590)
  set type($core.String value) => $_setString(2, value);
  @$pb.TagNumber(290836590)
  $core.bool hasType() => $_has(2);
  @$pb.TagNumber(290836590)
  void clearType() => $_clearField(290836590);
}

class Connection extends $pb.GeneratedMessage {
  factory Connection({
    $core.String? creationtime,
    $core.String? connectionarn,
    $core.String? lastmodifiedtime,
    $core.String? lastauthorizedtime,
    $core.String? name,
    $core.String? statereason,
    ConnectionState? connectionstate,
    ConnectionAuthorizationType? authorizationtype,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (lastauthorizedtime != null)
      result.lastauthorizedtime = lastauthorizedtime;
    if (name != null) result.name = name;
    if (statereason != null) result.statereason = statereason;
    if (connectionstate != null) result.connectionstate = connectionstate;
    if (authorizationtype != null) result.authorizationtype = authorizationtype;
    return result;
  }

  Connection._();

  factory Connection.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Connection.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Connection',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(250638066, _omitFieldNames ? '' : 'lastauthorizedtime')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aE<ConnectionState>(404323675, _omitFieldNames ? '' : 'connectionstate',
        enumValues: ConnectionState.values)
    ..aE<ConnectionAuthorizationType>(
        481976511, _omitFieldNames ? '' : 'authorizationtype',
        enumValues: ConnectionAuthorizationType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Connection clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Connection copyWith(void Function(Connection) updates) =>
      super.copyWith((message) => updates(message as Connection)) as Connection;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Connection create() => Connection._();
  @$core.override
  Connection createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Connection getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<Connection>(create);
  static Connection? _defaultInstance;

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(0);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(0);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(1);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(1);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(2);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(2);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(250638066)
  $core.String get lastauthorizedtime => $_getSZ(3);
  @$pb.TagNumber(250638066)
  set lastauthorizedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(250638066)
  $core.bool hasLastauthorizedtime() => $_has(3);
  @$pb.TagNumber(250638066)
  void clearLastauthorizedtime() => $_clearField(250638066);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(5);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(5, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(5);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(404323675)
  ConnectionState get connectionstate => $_getN(6);
  @$pb.TagNumber(404323675)
  set connectionstate(ConnectionState value) => $_setField(404323675, value);
  @$pb.TagNumber(404323675)
  $core.bool hasConnectionstate() => $_has(6);
  @$pb.TagNumber(404323675)
  void clearConnectionstate() => $_clearField(404323675);

  @$pb.TagNumber(481976511)
  ConnectionAuthorizationType get authorizationtype => $_getN(7);
  @$pb.TagNumber(481976511)
  set authorizationtype(ConnectionAuthorizationType value) =>
      $_setField(481976511, value);
  @$pb.TagNumber(481976511)
  $core.bool hasAuthorizationtype() => $_has(7);
  @$pb.TagNumber(481976511)
  void clearAuthorizationtype() => $_clearField(481976511);
}

class ConnectionApiKeyAuthResponseParameters extends $pb.GeneratedMessage {
  factory ConnectionApiKeyAuthResponseParameters({
    $core.String? apikeyname,
  }) {
    final result = create();
    if (apikeyname != null) result.apikeyname = apikeyname;
    return result;
  }

  ConnectionApiKeyAuthResponseParameters._();

  factory ConnectionApiKeyAuthResponseParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionApiKeyAuthResponseParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionApiKeyAuthResponseParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(43133502, _omitFieldNames ? '' : 'apikeyname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionApiKeyAuthResponseParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionApiKeyAuthResponseParameters copyWith(
          void Function(ConnectionApiKeyAuthResponseParameters) updates) =>
      super.copyWith((message) =>
              updates(message as ConnectionApiKeyAuthResponseParameters))
          as ConnectionApiKeyAuthResponseParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionApiKeyAuthResponseParameters create() =>
      ConnectionApiKeyAuthResponseParameters._();
  @$core.override
  ConnectionApiKeyAuthResponseParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionApiKeyAuthResponseParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ConnectionApiKeyAuthResponseParameters>(create);
  static ConnectionApiKeyAuthResponseParameters? _defaultInstance;

  @$pb.TagNumber(43133502)
  $core.String get apikeyname => $_getSZ(0);
  @$pb.TagNumber(43133502)
  set apikeyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(43133502)
  $core.bool hasApikeyname() => $_has(0);
  @$pb.TagNumber(43133502)
  void clearApikeyname() => $_clearField(43133502);
}

class ConnectionAuthResponseParameters extends $pb.GeneratedMessage {
  factory ConnectionAuthResponseParameters({
    ConnectionBasicAuthResponseParameters? basicauthparameters,
    ConnectionApiKeyAuthResponseParameters? apikeyauthparameters,
    ConnectionOAuthResponseParameters? oauthparameters,
    ConnectionHttpParameters? invocationhttpparameters,
  }) {
    final result = create();
    if (basicauthparameters != null)
      result.basicauthparameters = basicauthparameters;
    if (apikeyauthparameters != null)
      result.apikeyauthparameters = apikeyauthparameters;
    if (oauthparameters != null) result.oauthparameters = oauthparameters;
    if (invocationhttpparameters != null)
      result.invocationhttpparameters = invocationhttpparameters;
    return result;
  }

  ConnectionAuthResponseParameters._();

  factory ConnectionAuthResponseParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionAuthResponseParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionAuthResponseParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOM<ConnectionBasicAuthResponseParameters>(
        63965312, _omitFieldNames ? '' : 'basicauthparameters',
        subBuilder: ConnectionBasicAuthResponseParameters.create)
    ..aOM<ConnectionApiKeyAuthResponseParameters>(
        110622489, _omitFieldNames ? '' : 'apikeyauthparameters',
        subBuilder: ConnectionApiKeyAuthResponseParameters.create)
    ..aOM<ConnectionOAuthResponseParameters>(
        206836569, _omitFieldNames ? '' : 'oauthparameters',
        subBuilder: ConnectionOAuthResponseParameters.create)
    ..aOM<ConnectionHttpParameters>(
        499128572, _omitFieldNames ? '' : 'invocationhttpparameters',
        subBuilder: ConnectionHttpParameters.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionAuthResponseParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionAuthResponseParameters copyWith(
          void Function(ConnectionAuthResponseParameters) updates) =>
      super.copyWith(
              (message) => updates(message as ConnectionAuthResponseParameters))
          as ConnectionAuthResponseParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionAuthResponseParameters create() =>
      ConnectionAuthResponseParameters._();
  @$core.override
  ConnectionAuthResponseParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionAuthResponseParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectionAuthResponseParameters>(
          create);
  static ConnectionAuthResponseParameters? _defaultInstance;

  @$pb.TagNumber(63965312)
  ConnectionBasicAuthResponseParameters get basicauthparameters => $_getN(0);
  @$pb.TagNumber(63965312)
  set basicauthparameters(ConnectionBasicAuthResponseParameters value) =>
      $_setField(63965312, value);
  @$pb.TagNumber(63965312)
  $core.bool hasBasicauthparameters() => $_has(0);
  @$pb.TagNumber(63965312)
  void clearBasicauthparameters() => $_clearField(63965312);
  @$pb.TagNumber(63965312)
  ConnectionBasicAuthResponseParameters ensureBasicauthparameters() =>
      $_ensure(0);

  @$pb.TagNumber(110622489)
  ConnectionApiKeyAuthResponseParameters get apikeyauthparameters => $_getN(1);
  @$pb.TagNumber(110622489)
  set apikeyauthparameters(ConnectionApiKeyAuthResponseParameters value) =>
      $_setField(110622489, value);
  @$pb.TagNumber(110622489)
  $core.bool hasApikeyauthparameters() => $_has(1);
  @$pb.TagNumber(110622489)
  void clearApikeyauthparameters() => $_clearField(110622489);
  @$pb.TagNumber(110622489)
  ConnectionApiKeyAuthResponseParameters ensureApikeyauthparameters() =>
      $_ensure(1);

  @$pb.TagNumber(206836569)
  ConnectionOAuthResponseParameters get oauthparameters => $_getN(2);
  @$pb.TagNumber(206836569)
  set oauthparameters(ConnectionOAuthResponseParameters value) =>
      $_setField(206836569, value);
  @$pb.TagNumber(206836569)
  $core.bool hasOauthparameters() => $_has(2);
  @$pb.TagNumber(206836569)
  void clearOauthparameters() => $_clearField(206836569);
  @$pb.TagNumber(206836569)
  ConnectionOAuthResponseParameters ensureOauthparameters() => $_ensure(2);

  @$pb.TagNumber(499128572)
  ConnectionHttpParameters get invocationhttpparameters => $_getN(3);
  @$pb.TagNumber(499128572)
  set invocationhttpparameters(ConnectionHttpParameters value) =>
      $_setField(499128572, value);
  @$pb.TagNumber(499128572)
  $core.bool hasInvocationhttpparameters() => $_has(3);
  @$pb.TagNumber(499128572)
  void clearInvocationhttpparameters() => $_clearField(499128572);
  @$pb.TagNumber(499128572)
  ConnectionHttpParameters ensureInvocationhttpparameters() => $_ensure(3);
}

class ConnectionBasicAuthResponseParameters extends $pb.GeneratedMessage {
  factory ConnectionBasicAuthResponseParameters({
    $core.String? username,
  }) {
    final result = create();
    if (username != null) result.username = username;
    return result;
  }

  ConnectionBasicAuthResponseParameters._();

  factory ConnectionBasicAuthResponseParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionBasicAuthResponseParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionBasicAuthResponseParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(470340826, _omitFieldNames ? '' : 'username')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionBasicAuthResponseParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionBasicAuthResponseParameters copyWith(
          void Function(ConnectionBasicAuthResponseParameters) updates) =>
      super.copyWith((message) =>
              updates(message as ConnectionBasicAuthResponseParameters))
          as ConnectionBasicAuthResponseParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionBasicAuthResponseParameters create() =>
      ConnectionBasicAuthResponseParameters._();
  @$core.override
  ConnectionBasicAuthResponseParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionBasicAuthResponseParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ConnectionBasicAuthResponseParameters>(create);
  static ConnectionBasicAuthResponseParameters? _defaultInstance;

  @$pb.TagNumber(470340826)
  $core.String get username => $_getSZ(0);
  @$pb.TagNumber(470340826)
  set username($core.String value) => $_setString(0, value);
  @$pb.TagNumber(470340826)
  $core.bool hasUsername() => $_has(0);
  @$pb.TagNumber(470340826)
  void clearUsername() => $_clearField(470340826);
}

class ConnectionBodyParameter extends $pb.GeneratedMessage {
  factory ConnectionBodyParameter({
    $core.bool? isvaluesecret,
    $core.String? key,
    $core.String? value,
  }) {
    final result = create();
    if (isvaluesecret != null) result.isvaluesecret = isvaluesecret;
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    return result;
  }

  ConnectionBodyParameter._();

  factory ConnectionBodyParameter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionBodyParameter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionBodyParameter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOB(157884907, _omitFieldNames ? '' : 'isvaluesecret')
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionBodyParameter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionBodyParameter copyWith(
          void Function(ConnectionBodyParameter) updates) =>
      super.copyWith((message) => updates(message as ConnectionBodyParameter))
          as ConnectionBodyParameter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionBodyParameter create() => ConnectionBodyParameter._();
  @$core.override
  ConnectionBodyParameter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionBodyParameter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectionBodyParameter>(create);
  static ConnectionBodyParameter? _defaultInstance;

  @$pb.TagNumber(157884907)
  $core.bool get isvaluesecret => $_getBF(0);
  @$pb.TagNumber(157884907)
  set isvaluesecret($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(157884907)
  $core.bool hasIsvaluesecret() => $_has(0);
  @$pb.TagNumber(157884907)
  void clearIsvaluesecret() => $_clearField(157884907);

  @$pb.TagNumber(219859213)
  $core.String get key => $_getSZ(1);
  @$pb.TagNumber(219859213)
  set key($core.String value) => $_setString(1, value);
  @$pb.TagNumber(219859213)
  $core.bool hasKey() => $_has(1);
  @$pb.TagNumber(219859213)
  void clearKey() => $_clearField(219859213);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(2);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(2, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(2);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class ConnectionHeaderParameter extends $pb.GeneratedMessage {
  factory ConnectionHeaderParameter({
    $core.bool? isvaluesecret,
    $core.String? key,
    $core.String? value,
  }) {
    final result = create();
    if (isvaluesecret != null) result.isvaluesecret = isvaluesecret;
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    return result;
  }

  ConnectionHeaderParameter._();

  factory ConnectionHeaderParameter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionHeaderParameter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionHeaderParameter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOB(157884907, _omitFieldNames ? '' : 'isvaluesecret')
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionHeaderParameter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionHeaderParameter copyWith(
          void Function(ConnectionHeaderParameter) updates) =>
      super.copyWith((message) => updates(message as ConnectionHeaderParameter))
          as ConnectionHeaderParameter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionHeaderParameter create() => ConnectionHeaderParameter._();
  @$core.override
  ConnectionHeaderParameter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionHeaderParameter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectionHeaderParameter>(create);
  static ConnectionHeaderParameter? _defaultInstance;

  @$pb.TagNumber(157884907)
  $core.bool get isvaluesecret => $_getBF(0);
  @$pb.TagNumber(157884907)
  set isvaluesecret($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(157884907)
  $core.bool hasIsvaluesecret() => $_has(0);
  @$pb.TagNumber(157884907)
  void clearIsvaluesecret() => $_clearField(157884907);

  @$pb.TagNumber(219859213)
  $core.String get key => $_getSZ(1);
  @$pb.TagNumber(219859213)
  set key($core.String value) => $_setString(1, value);
  @$pb.TagNumber(219859213)
  $core.bool hasKey() => $_has(1);
  @$pb.TagNumber(219859213)
  void clearKey() => $_clearField(219859213);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(2);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(2, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(2);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class ConnectionHttpParameters extends $pb.GeneratedMessage {
  factory ConnectionHttpParameters({
    $core.Iterable<ConnectionHeaderParameter>? headerparameters,
    $core.Iterable<ConnectionQueryStringParameter>? querystringparameters,
    $core.Iterable<ConnectionBodyParameter>? bodyparameters,
  }) {
    final result = create();
    if (headerparameters != null)
      result.headerparameters.addAll(headerparameters);
    if (querystringparameters != null)
      result.querystringparameters.addAll(querystringparameters);
    if (bodyparameters != null) result.bodyparameters.addAll(bodyparameters);
    return result;
  }

  ConnectionHttpParameters._();

  factory ConnectionHttpParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionHttpParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionHttpParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<ConnectionHeaderParameter>(
        148944757, _omitFieldNames ? '' : 'headerparameters',
        subBuilder: ConnectionHeaderParameter.create)
    ..pPM<ConnectionQueryStringParameter>(
        258002263, _omitFieldNames ? '' : 'querystringparameters',
        subBuilder: ConnectionQueryStringParameter.create)
    ..pPM<ConnectionBodyParameter>(
        531363370, _omitFieldNames ? '' : 'bodyparameters',
        subBuilder: ConnectionBodyParameter.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionHttpParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionHttpParameters copyWith(
          void Function(ConnectionHttpParameters) updates) =>
      super.copyWith((message) => updates(message as ConnectionHttpParameters))
          as ConnectionHttpParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionHttpParameters create() => ConnectionHttpParameters._();
  @$core.override
  ConnectionHttpParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionHttpParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectionHttpParameters>(create);
  static ConnectionHttpParameters? _defaultInstance;

  @$pb.TagNumber(148944757)
  $pb.PbList<ConnectionHeaderParameter> get headerparameters => $_getList(0);

  @$pb.TagNumber(258002263)
  $pb.PbList<ConnectionQueryStringParameter> get querystringparameters =>
      $_getList(1);

  @$pb.TagNumber(531363370)
  $pb.PbList<ConnectionBodyParameter> get bodyparameters => $_getList(2);
}

class ConnectionOAuthClientResponseParameters extends $pb.GeneratedMessage {
  factory ConnectionOAuthClientResponseParameters({
    $core.String? clientid,
  }) {
    final result = create();
    if (clientid != null) result.clientid = clientid;
    return result;
  }

  ConnectionOAuthClientResponseParameters._();

  factory ConnectionOAuthClientResponseParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionOAuthClientResponseParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionOAuthClientResponseParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(448889284, _omitFieldNames ? '' : 'clientid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionOAuthClientResponseParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionOAuthClientResponseParameters copyWith(
          void Function(ConnectionOAuthClientResponseParameters) updates) =>
      super.copyWith((message) =>
              updates(message as ConnectionOAuthClientResponseParameters))
          as ConnectionOAuthClientResponseParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionOAuthClientResponseParameters create() =>
      ConnectionOAuthClientResponseParameters._();
  @$core.override
  ConnectionOAuthClientResponseParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionOAuthClientResponseParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ConnectionOAuthClientResponseParameters>(create);
  static ConnectionOAuthClientResponseParameters? _defaultInstance;

  @$pb.TagNumber(448889284)
  $core.String get clientid => $_getSZ(0);
  @$pb.TagNumber(448889284)
  set clientid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(448889284)
  $core.bool hasClientid() => $_has(0);
  @$pb.TagNumber(448889284)
  void clearClientid() => $_clearField(448889284);
}

class ConnectionOAuthResponseParameters extends $pb.GeneratedMessage {
  factory ConnectionOAuthResponseParameters({
    ConnectionHttpParameters? oauthhttpparameters,
    ConnectionOAuthClientResponseParameters? clientparameters,
    ConnectionOAuthHttpMethod? httpmethod,
    $core.String? authorizationendpoint,
  }) {
    final result = create();
    if (oauthhttpparameters != null)
      result.oauthhttpparameters = oauthhttpparameters;
    if (clientparameters != null) result.clientparameters = clientparameters;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (authorizationendpoint != null)
      result.authorizationendpoint = authorizationendpoint;
    return result;
  }

  ConnectionOAuthResponseParameters._();

  factory ConnectionOAuthResponseParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionOAuthResponseParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionOAuthResponseParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOM<ConnectionHttpParameters>(
        10294537, _omitFieldNames ? '' : 'oauthhttpparameters',
        subBuilder: ConnectionHttpParameters.create)
    ..aOM<ConnectionOAuthClientResponseParameters>(
        246864127, _omitFieldNames ? '' : 'clientparameters',
        subBuilder: ConnectionOAuthClientResponseParameters.create)
    ..aE<ConnectionOAuthHttpMethod>(
        398394961, _omitFieldNames ? '' : 'httpmethod',
        enumValues: ConnectionOAuthHttpMethod.values)
    ..aOS(427938596, _omitFieldNames ? '' : 'authorizationendpoint')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionOAuthResponseParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionOAuthResponseParameters copyWith(
          void Function(ConnectionOAuthResponseParameters) updates) =>
      super.copyWith((message) =>
              updates(message as ConnectionOAuthResponseParameters))
          as ConnectionOAuthResponseParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionOAuthResponseParameters create() =>
      ConnectionOAuthResponseParameters._();
  @$core.override
  ConnectionOAuthResponseParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionOAuthResponseParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectionOAuthResponseParameters>(
          create);
  static ConnectionOAuthResponseParameters? _defaultInstance;

  @$pb.TagNumber(10294537)
  ConnectionHttpParameters get oauthhttpparameters => $_getN(0);
  @$pb.TagNumber(10294537)
  set oauthhttpparameters(ConnectionHttpParameters value) =>
      $_setField(10294537, value);
  @$pb.TagNumber(10294537)
  $core.bool hasOauthhttpparameters() => $_has(0);
  @$pb.TagNumber(10294537)
  void clearOauthhttpparameters() => $_clearField(10294537);
  @$pb.TagNumber(10294537)
  ConnectionHttpParameters ensureOauthhttpparameters() => $_ensure(0);

  @$pb.TagNumber(246864127)
  ConnectionOAuthClientResponseParameters get clientparameters => $_getN(1);
  @$pb.TagNumber(246864127)
  set clientparameters(ConnectionOAuthClientResponseParameters value) =>
      $_setField(246864127, value);
  @$pb.TagNumber(246864127)
  $core.bool hasClientparameters() => $_has(1);
  @$pb.TagNumber(246864127)
  void clearClientparameters() => $_clearField(246864127);
  @$pb.TagNumber(246864127)
  ConnectionOAuthClientResponseParameters ensureClientparameters() =>
      $_ensure(1);

  @$pb.TagNumber(398394961)
  ConnectionOAuthHttpMethod get httpmethod => $_getN(2);
  @$pb.TagNumber(398394961)
  set httpmethod(ConnectionOAuthHttpMethod value) =>
      $_setField(398394961, value);
  @$pb.TagNumber(398394961)
  $core.bool hasHttpmethod() => $_has(2);
  @$pb.TagNumber(398394961)
  void clearHttpmethod() => $_clearField(398394961);

  @$pb.TagNumber(427938596)
  $core.String get authorizationendpoint => $_getSZ(3);
  @$pb.TagNumber(427938596)
  set authorizationendpoint($core.String value) => $_setString(3, value);
  @$pb.TagNumber(427938596)
  $core.bool hasAuthorizationendpoint() => $_has(3);
  @$pb.TagNumber(427938596)
  void clearAuthorizationendpoint() => $_clearField(427938596);
}

class ConnectionQueryStringParameter extends $pb.GeneratedMessage {
  factory ConnectionQueryStringParameter({
    $core.bool? isvaluesecret,
    $core.String? key,
    $core.String? value,
  }) {
    final result = create();
    if (isvaluesecret != null) result.isvaluesecret = isvaluesecret;
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    return result;
  }

  ConnectionQueryStringParameter._();

  factory ConnectionQueryStringParameter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionQueryStringParameter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionQueryStringParameter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOB(157884907, _omitFieldNames ? '' : 'isvaluesecret')
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionQueryStringParameter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionQueryStringParameter copyWith(
          void Function(ConnectionQueryStringParameter) updates) =>
      super.copyWith(
              (message) => updates(message as ConnectionQueryStringParameter))
          as ConnectionQueryStringParameter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionQueryStringParameter create() =>
      ConnectionQueryStringParameter._();
  @$core.override
  ConnectionQueryStringParameter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionQueryStringParameter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectionQueryStringParameter>(create);
  static ConnectionQueryStringParameter? _defaultInstance;

  @$pb.TagNumber(157884907)
  $core.bool get isvaluesecret => $_getBF(0);
  @$pb.TagNumber(157884907)
  set isvaluesecret($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(157884907)
  $core.bool hasIsvaluesecret() => $_has(0);
  @$pb.TagNumber(157884907)
  void clearIsvaluesecret() => $_clearField(157884907);

  @$pb.TagNumber(219859213)
  $core.String get key => $_getSZ(1);
  @$pb.TagNumber(219859213)
  set key($core.String value) => $_setString(1, value);
  @$pb.TagNumber(219859213)
  $core.bool hasKey() => $_has(1);
  @$pb.TagNumber(219859213)
  void clearKey() => $_clearField(219859213);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(2);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(2, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(2);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class CreateApiDestinationRequest extends $pb.GeneratedMessage {
  factory CreateApiDestinationRequest({
    $core.String? description,
    $core.String? connectionarn,
    $core.String? name,
    $core.int? invocationratelimitpersecond,
    ApiDestinationHttpMethod? httpmethod,
    $core.String? invocationendpoint,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (name != null) result.name = name;
    if (invocationratelimitpersecond != null)
      result.invocationratelimitpersecond = invocationratelimitpersecond;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (invocationendpoint != null)
      result.invocationendpoint = invocationendpoint;
    return result;
  }

  CreateApiDestinationRequest._();

  factory CreateApiDestinationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateApiDestinationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateApiDestinationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aI(295327816, _omitFieldNames ? '' : 'invocationratelimitpersecond')
    ..aE<ApiDestinationHttpMethod>(
        398394961, _omitFieldNames ? '' : 'httpmethod',
        enumValues: ApiDestinationHttpMethod.values)
    ..aOS(411800759, _omitFieldNames ? '' : 'invocationendpoint')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateApiDestinationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateApiDestinationRequest copyWith(
          void Function(CreateApiDestinationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreateApiDestinationRequest))
          as CreateApiDestinationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateApiDestinationRequest create() =>
      CreateApiDestinationRequest._();
  @$core.override
  CreateApiDestinationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateApiDestinationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateApiDestinationRequest>(create);
  static CreateApiDestinationRequest? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(1);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(1);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(295327816)
  $core.int get invocationratelimitpersecond => $_getIZ(3);
  @$pb.TagNumber(295327816)
  set invocationratelimitpersecond($core.int value) =>
      $_setSignedInt32(3, value);
  @$pb.TagNumber(295327816)
  $core.bool hasInvocationratelimitpersecond() => $_has(3);
  @$pb.TagNumber(295327816)
  void clearInvocationratelimitpersecond() => $_clearField(295327816);

  @$pb.TagNumber(398394961)
  ApiDestinationHttpMethod get httpmethod => $_getN(4);
  @$pb.TagNumber(398394961)
  set httpmethod(ApiDestinationHttpMethod value) =>
      $_setField(398394961, value);
  @$pb.TagNumber(398394961)
  $core.bool hasHttpmethod() => $_has(4);
  @$pb.TagNumber(398394961)
  void clearHttpmethod() => $_clearField(398394961);

  @$pb.TagNumber(411800759)
  $core.String get invocationendpoint => $_getSZ(5);
  @$pb.TagNumber(411800759)
  set invocationendpoint($core.String value) => $_setString(5, value);
  @$pb.TagNumber(411800759)
  $core.bool hasInvocationendpoint() => $_has(5);
  @$pb.TagNumber(411800759)
  void clearInvocationendpoint() => $_clearField(411800759);
}

class CreateApiDestinationResponse extends $pb.GeneratedMessage {
  factory CreateApiDestinationResponse({
    ApiDestinationState? apidestinationstate,
    $core.String? apidestinationarn,
    $core.String? creationtime,
    $core.String? lastmodifiedtime,
  }) {
    final result = create();
    if (apidestinationstate != null)
      result.apidestinationstate = apidestinationstate;
    if (apidestinationarn != null) result.apidestinationarn = apidestinationarn;
    if (creationtime != null) result.creationtime = creationtime;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    return result;
  }

  CreateApiDestinationResponse._();

  factory CreateApiDestinationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateApiDestinationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateApiDestinationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aE<ApiDestinationState>(
        13153343, _omitFieldNames ? '' : 'apidestinationstate',
        enumValues: ApiDestinationState.values)
    ..aOS(90996885, _omitFieldNames ? '' : 'apidestinationarn')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateApiDestinationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateApiDestinationResponse copyWith(
          void Function(CreateApiDestinationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreateApiDestinationResponse))
          as CreateApiDestinationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateApiDestinationResponse create() =>
      CreateApiDestinationResponse._();
  @$core.override
  CreateApiDestinationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateApiDestinationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateApiDestinationResponse>(create);
  static CreateApiDestinationResponse? _defaultInstance;

  @$pb.TagNumber(13153343)
  ApiDestinationState get apidestinationstate => $_getN(0);
  @$pb.TagNumber(13153343)
  set apidestinationstate(ApiDestinationState value) =>
      $_setField(13153343, value);
  @$pb.TagNumber(13153343)
  $core.bool hasApidestinationstate() => $_has(0);
  @$pb.TagNumber(13153343)
  void clearApidestinationstate() => $_clearField(13153343);

  @$pb.TagNumber(90996885)
  $core.String get apidestinationarn => $_getSZ(1);
  @$pb.TagNumber(90996885)
  set apidestinationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(90996885)
  $core.bool hasApidestinationarn() => $_has(1);
  @$pb.TagNumber(90996885)
  void clearApidestinationarn() => $_clearField(90996885);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(3);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(3);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);
}

class CreateArchiveRequest extends $pb.GeneratedMessage {
  factory CreateArchiveRequest({
    $core.String? archivename,
    $core.String? description,
    $core.String? eventpattern,
    $core.int? retentiondays,
    $core.String? eventsourcearn,
  }) {
    final result = create();
    if (archivename != null) result.archivename = archivename;
    if (description != null) result.description = description;
    if (eventpattern != null) result.eventpattern = eventpattern;
    if (retentiondays != null) result.retentiondays = retentiondays;
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    return result;
  }

  CreateArchiveRequest._();

  factory CreateArchiveRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateArchiveRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateArchiveRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(88048487, _omitFieldNames ? '' : 'archivename')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(233487232, _omitFieldNames ? '' : 'eventpattern')
    ..aI(267894223, _omitFieldNames ? '' : 'retentiondays')
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateArchiveRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateArchiveRequest copyWith(void Function(CreateArchiveRequest) updates) =>
      super.copyWith((message) => updates(message as CreateArchiveRequest))
          as CreateArchiveRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateArchiveRequest create() => CreateArchiveRequest._();
  @$core.override
  CreateArchiveRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateArchiveRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateArchiveRequest>(create);
  static CreateArchiveRequest? _defaultInstance;

  @$pb.TagNumber(88048487)
  $core.String get archivename => $_getSZ(0);
  @$pb.TagNumber(88048487)
  set archivename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88048487)
  $core.bool hasArchivename() => $_has(0);
  @$pb.TagNumber(88048487)
  void clearArchivename() => $_clearField(88048487);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(233487232)
  $core.String get eventpattern => $_getSZ(2);
  @$pb.TagNumber(233487232)
  set eventpattern($core.String value) => $_setString(2, value);
  @$pb.TagNumber(233487232)
  $core.bool hasEventpattern() => $_has(2);
  @$pb.TagNumber(233487232)
  void clearEventpattern() => $_clearField(233487232);

  @$pb.TagNumber(267894223)
  $core.int get retentiondays => $_getIZ(3);
  @$pb.TagNumber(267894223)
  set retentiondays($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(267894223)
  $core.bool hasRetentiondays() => $_has(3);
  @$pb.TagNumber(267894223)
  void clearRetentiondays() => $_clearField(267894223);

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(4);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(4);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);
}

class CreateArchiveResponse extends $pb.GeneratedMessage {
  factory CreateArchiveResponse({
    $core.String? archivearn,
    $core.String? creationtime,
    $core.String? statereason,
    ArchiveState? state,
  }) {
    final result = create();
    if (archivearn != null) result.archivearn = archivearn;
    if (creationtime != null) result.creationtime = creationtime;
    if (statereason != null) result.statereason = statereason;
    if (state != null) result.state = state;
    return result;
  }

  CreateArchiveResponse._();

  factory CreateArchiveResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateArchiveResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateArchiveResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(56866685, _omitFieldNames ? '' : 'archivearn')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aE<ArchiveState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ArchiveState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateArchiveResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateArchiveResponse copyWith(
          void Function(CreateArchiveResponse) updates) =>
      super.copyWith((message) => updates(message as CreateArchiveResponse))
          as CreateArchiveResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateArchiveResponse create() => CreateArchiveResponse._();
  @$core.override
  CreateArchiveResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateArchiveResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateArchiveResponse>(create);
  static CreateArchiveResponse? _defaultInstance;

  @$pb.TagNumber(56866685)
  $core.String get archivearn => $_getSZ(0);
  @$pb.TagNumber(56866685)
  set archivearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56866685)
  $core.bool hasArchivearn() => $_has(0);
  @$pb.TagNumber(56866685)
  void clearArchivearn() => $_clearField(56866685);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(2);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(2, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(2);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(502047895)
  ArchiveState get state => $_getN(3);
  @$pb.TagNumber(502047895)
  set state(ArchiveState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(3);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class CreateConnectionApiKeyAuthRequestParameters extends $pb.GeneratedMessage {
  factory CreateConnectionApiKeyAuthRequestParameters({
    $core.String? apikeyname,
    $core.String? apikeyvalue,
  }) {
    final result = create();
    if (apikeyname != null) result.apikeyname = apikeyname;
    if (apikeyvalue != null) result.apikeyvalue = apikeyvalue;
    return result;
  }

  CreateConnectionApiKeyAuthRequestParameters._();

  factory CreateConnectionApiKeyAuthRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateConnectionApiKeyAuthRequestParameters.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateConnectionApiKeyAuthRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(43133502, _omitFieldNames ? '' : 'apikeyname')
    ..aOS(93786144, _omitFieldNames ? '' : 'apikeyvalue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionApiKeyAuthRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionApiKeyAuthRequestParameters copyWith(
          void Function(CreateConnectionApiKeyAuthRequestParameters) updates) =>
      super.copyWith((message) =>
              updates(message as CreateConnectionApiKeyAuthRequestParameters))
          as CreateConnectionApiKeyAuthRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateConnectionApiKeyAuthRequestParameters create() =>
      CreateConnectionApiKeyAuthRequestParameters._();
  @$core.override
  CreateConnectionApiKeyAuthRequestParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateConnectionApiKeyAuthRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateConnectionApiKeyAuthRequestParameters>(create);
  static CreateConnectionApiKeyAuthRequestParameters? _defaultInstance;

  @$pb.TagNumber(43133502)
  $core.String get apikeyname => $_getSZ(0);
  @$pb.TagNumber(43133502)
  set apikeyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(43133502)
  $core.bool hasApikeyname() => $_has(0);
  @$pb.TagNumber(43133502)
  void clearApikeyname() => $_clearField(43133502);

  @$pb.TagNumber(93786144)
  $core.String get apikeyvalue => $_getSZ(1);
  @$pb.TagNumber(93786144)
  set apikeyvalue($core.String value) => $_setString(1, value);
  @$pb.TagNumber(93786144)
  $core.bool hasApikeyvalue() => $_has(1);
  @$pb.TagNumber(93786144)
  void clearApikeyvalue() => $_clearField(93786144);
}

class CreateConnectionAuthRequestParameters extends $pb.GeneratedMessage {
  factory CreateConnectionAuthRequestParameters({
    CreateConnectionBasicAuthRequestParameters? basicauthparameters,
    CreateConnectionApiKeyAuthRequestParameters? apikeyauthparameters,
    CreateConnectionOAuthRequestParameters? oauthparameters,
    ConnectionHttpParameters? invocationhttpparameters,
  }) {
    final result = create();
    if (basicauthparameters != null)
      result.basicauthparameters = basicauthparameters;
    if (apikeyauthparameters != null)
      result.apikeyauthparameters = apikeyauthparameters;
    if (oauthparameters != null) result.oauthparameters = oauthparameters;
    if (invocationhttpparameters != null)
      result.invocationhttpparameters = invocationhttpparameters;
    return result;
  }

  CreateConnectionAuthRequestParameters._();

  factory CreateConnectionAuthRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateConnectionAuthRequestParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateConnectionAuthRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOM<CreateConnectionBasicAuthRequestParameters>(
        63965312, _omitFieldNames ? '' : 'basicauthparameters',
        subBuilder: CreateConnectionBasicAuthRequestParameters.create)
    ..aOM<CreateConnectionApiKeyAuthRequestParameters>(
        110622489, _omitFieldNames ? '' : 'apikeyauthparameters',
        subBuilder: CreateConnectionApiKeyAuthRequestParameters.create)
    ..aOM<CreateConnectionOAuthRequestParameters>(
        206836569, _omitFieldNames ? '' : 'oauthparameters',
        subBuilder: CreateConnectionOAuthRequestParameters.create)
    ..aOM<ConnectionHttpParameters>(
        499128572, _omitFieldNames ? '' : 'invocationhttpparameters',
        subBuilder: ConnectionHttpParameters.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionAuthRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionAuthRequestParameters copyWith(
          void Function(CreateConnectionAuthRequestParameters) updates) =>
      super.copyWith((message) =>
              updates(message as CreateConnectionAuthRequestParameters))
          as CreateConnectionAuthRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateConnectionAuthRequestParameters create() =>
      CreateConnectionAuthRequestParameters._();
  @$core.override
  CreateConnectionAuthRequestParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateConnectionAuthRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateConnectionAuthRequestParameters>(create);
  static CreateConnectionAuthRequestParameters? _defaultInstance;

  @$pb.TagNumber(63965312)
  CreateConnectionBasicAuthRequestParameters get basicauthparameters =>
      $_getN(0);
  @$pb.TagNumber(63965312)
  set basicauthparameters(CreateConnectionBasicAuthRequestParameters value) =>
      $_setField(63965312, value);
  @$pb.TagNumber(63965312)
  $core.bool hasBasicauthparameters() => $_has(0);
  @$pb.TagNumber(63965312)
  void clearBasicauthparameters() => $_clearField(63965312);
  @$pb.TagNumber(63965312)
  CreateConnectionBasicAuthRequestParameters ensureBasicauthparameters() =>
      $_ensure(0);

  @$pb.TagNumber(110622489)
  CreateConnectionApiKeyAuthRequestParameters get apikeyauthparameters =>
      $_getN(1);
  @$pb.TagNumber(110622489)
  set apikeyauthparameters(CreateConnectionApiKeyAuthRequestParameters value) =>
      $_setField(110622489, value);
  @$pb.TagNumber(110622489)
  $core.bool hasApikeyauthparameters() => $_has(1);
  @$pb.TagNumber(110622489)
  void clearApikeyauthparameters() => $_clearField(110622489);
  @$pb.TagNumber(110622489)
  CreateConnectionApiKeyAuthRequestParameters ensureApikeyauthparameters() =>
      $_ensure(1);

  @$pb.TagNumber(206836569)
  CreateConnectionOAuthRequestParameters get oauthparameters => $_getN(2);
  @$pb.TagNumber(206836569)
  set oauthparameters(CreateConnectionOAuthRequestParameters value) =>
      $_setField(206836569, value);
  @$pb.TagNumber(206836569)
  $core.bool hasOauthparameters() => $_has(2);
  @$pb.TagNumber(206836569)
  void clearOauthparameters() => $_clearField(206836569);
  @$pb.TagNumber(206836569)
  CreateConnectionOAuthRequestParameters ensureOauthparameters() => $_ensure(2);

  @$pb.TagNumber(499128572)
  ConnectionHttpParameters get invocationhttpparameters => $_getN(3);
  @$pb.TagNumber(499128572)
  set invocationhttpparameters(ConnectionHttpParameters value) =>
      $_setField(499128572, value);
  @$pb.TagNumber(499128572)
  $core.bool hasInvocationhttpparameters() => $_has(3);
  @$pb.TagNumber(499128572)
  void clearInvocationhttpparameters() => $_clearField(499128572);
  @$pb.TagNumber(499128572)
  ConnectionHttpParameters ensureInvocationhttpparameters() => $_ensure(3);
}

class CreateConnectionBasicAuthRequestParameters extends $pb.GeneratedMessage {
  factory CreateConnectionBasicAuthRequestParameters({
    $core.String? password,
    $core.String? username,
  }) {
    final result = create();
    if (password != null) result.password = password;
    if (username != null) result.username = username;
    return result;
  }

  CreateConnectionBasicAuthRequestParameters._();

  factory CreateConnectionBasicAuthRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateConnectionBasicAuthRequestParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateConnectionBasicAuthRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(214108217, _omitFieldNames ? '' : 'password')
    ..aOS(470340826, _omitFieldNames ? '' : 'username')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionBasicAuthRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionBasicAuthRequestParameters copyWith(
          void Function(CreateConnectionBasicAuthRequestParameters) updates) =>
      super.copyWith((message) =>
              updates(message as CreateConnectionBasicAuthRequestParameters))
          as CreateConnectionBasicAuthRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateConnectionBasicAuthRequestParameters create() =>
      CreateConnectionBasicAuthRequestParameters._();
  @$core.override
  CreateConnectionBasicAuthRequestParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateConnectionBasicAuthRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateConnectionBasicAuthRequestParameters>(create);
  static CreateConnectionBasicAuthRequestParameters? _defaultInstance;

  @$pb.TagNumber(214108217)
  $core.String get password => $_getSZ(0);
  @$pb.TagNumber(214108217)
  set password($core.String value) => $_setString(0, value);
  @$pb.TagNumber(214108217)
  $core.bool hasPassword() => $_has(0);
  @$pb.TagNumber(214108217)
  void clearPassword() => $_clearField(214108217);

  @$pb.TagNumber(470340826)
  $core.String get username => $_getSZ(1);
  @$pb.TagNumber(470340826)
  set username($core.String value) => $_setString(1, value);
  @$pb.TagNumber(470340826)
  $core.bool hasUsername() => $_has(1);
  @$pb.TagNumber(470340826)
  void clearUsername() => $_clearField(470340826);
}

class CreateConnectionOAuthClientRequestParameters
    extends $pb.GeneratedMessage {
  factory CreateConnectionOAuthClientRequestParameters({
    $core.String? clientid,
    $core.String? clientsecret,
  }) {
    final result = create();
    if (clientid != null) result.clientid = clientid;
    if (clientsecret != null) result.clientsecret = clientsecret;
    return result;
  }

  CreateConnectionOAuthClientRequestParameters._();

  factory CreateConnectionOAuthClientRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateConnectionOAuthClientRequestParameters.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateConnectionOAuthClientRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(448889284, _omitFieldNames ? '' : 'clientid')
    ..aOS(500734711, _omitFieldNames ? '' : 'clientsecret')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionOAuthClientRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionOAuthClientRequestParameters copyWith(
          void Function(CreateConnectionOAuthClientRequestParameters)
              updates) =>
      super.copyWith((message) =>
              updates(message as CreateConnectionOAuthClientRequestParameters))
          as CreateConnectionOAuthClientRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateConnectionOAuthClientRequestParameters create() =>
      CreateConnectionOAuthClientRequestParameters._();
  @$core.override
  CreateConnectionOAuthClientRequestParameters createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static CreateConnectionOAuthClientRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateConnectionOAuthClientRequestParameters>(create);
  static CreateConnectionOAuthClientRequestParameters? _defaultInstance;

  @$pb.TagNumber(448889284)
  $core.String get clientid => $_getSZ(0);
  @$pb.TagNumber(448889284)
  set clientid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(448889284)
  $core.bool hasClientid() => $_has(0);
  @$pb.TagNumber(448889284)
  void clearClientid() => $_clearField(448889284);

  @$pb.TagNumber(500734711)
  $core.String get clientsecret => $_getSZ(1);
  @$pb.TagNumber(500734711)
  set clientsecret($core.String value) => $_setString(1, value);
  @$pb.TagNumber(500734711)
  $core.bool hasClientsecret() => $_has(1);
  @$pb.TagNumber(500734711)
  void clearClientsecret() => $_clearField(500734711);
}

class CreateConnectionOAuthRequestParameters extends $pb.GeneratedMessage {
  factory CreateConnectionOAuthRequestParameters({
    ConnectionHttpParameters? oauthhttpparameters,
    CreateConnectionOAuthClientRequestParameters? clientparameters,
    ConnectionOAuthHttpMethod? httpmethod,
    $core.String? authorizationendpoint,
  }) {
    final result = create();
    if (oauthhttpparameters != null)
      result.oauthhttpparameters = oauthhttpparameters;
    if (clientparameters != null) result.clientparameters = clientparameters;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (authorizationendpoint != null)
      result.authorizationendpoint = authorizationendpoint;
    return result;
  }

  CreateConnectionOAuthRequestParameters._();

  factory CreateConnectionOAuthRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateConnectionOAuthRequestParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateConnectionOAuthRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOM<ConnectionHttpParameters>(
        10294537, _omitFieldNames ? '' : 'oauthhttpparameters',
        subBuilder: ConnectionHttpParameters.create)
    ..aOM<CreateConnectionOAuthClientRequestParameters>(
        246864127, _omitFieldNames ? '' : 'clientparameters',
        subBuilder: CreateConnectionOAuthClientRequestParameters.create)
    ..aE<ConnectionOAuthHttpMethod>(
        398394961, _omitFieldNames ? '' : 'httpmethod',
        enumValues: ConnectionOAuthHttpMethod.values)
    ..aOS(427938596, _omitFieldNames ? '' : 'authorizationendpoint')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionOAuthRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionOAuthRequestParameters copyWith(
          void Function(CreateConnectionOAuthRequestParameters) updates) =>
      super.copyWith((message) =>
              updates(message as CreateConnectionOAuthRequestParameters))
          as CreateConnectionOAuthRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateConnectionOAuthRequestParameters create() =>
      CreateConnectionOAuthRequestParameters._();
  @$core.override
  CreateConnectionOAuthRequestParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateConnectionOAuthRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          CreateConnectionOAuthRequestParameters>(create);
  static CreateConnectionOAuthRequestParameters? _defaultInstance;

  @$pb.TagNumber(10294537)
  ConnectionHttpParameters get oauthhttpparameters => $_getN(0);
  @$pb.TagNumber(10294537)
  set oauthhttpparameters(ConnectionHttpParameters value) =>
      $_setField(10294537, value);
  @$pb.TagNumber(10294537)
  $core.bool hasOauthhttpparameters() => $_has(0);
  @$pb.TagNumber(10294537)
  void clearOauthhttpparameters() => $_clearField(10294537);
  @$pb.TagNumber(10294537)
  ConnectionHttpParameters ensureOauthhttpparameters() => $_ensure(0);

  @$pb.TagNumber(246864127)
  CreateConnectionOAuthClientRequestParameters get clientparameters =>
      $_getN(1);
  @$pb.TagNumber(246864127)
  set clientparameters(CreateConnectionOAuthClientRequestParameters value) =>
      $_setField(246864127, value);
  @$pb.TagNumber(246864127)
  $core.bool hasClientparameters() => $_has(1);
  @$pb.TagNumber(246864127)
  void clearClientparameters() => $_clearField(246864127);
  @$pb.TagNumber(246864127)
  CreateConnectionOAuthClientRequestParameters ensureClientparameters() =>
      $_ensure(1);

  @$pb.TagNumber(398394961)
  ConnectionOAuthHttpMethod get httpmethod => $_getN(2);
  @$pb.TagNumber(398394961)
  set httpmethod(ConnectionOAuthHttpMethod value) =>
      $_setField(398394961, value);
  @$pb.TagNumber(398394961)
  $core.bool hasHttpmethod() => $_has(2);
  @$pb.TagNumber(398394961)
  void clearHttpmethod() => $_clearField(398394961);

  @$pb.TagNumber(427938596)
  $core.String get authorizationendpoint => $_getSZ(3);
  @$pb.TagNumber(427938596)
  set authorizationendpoint($core.String value) => $_setString(3, value);
  @$pb.TagNumber(427938596)
  $core.bool hasAuthorizationendpoint() => $_has(3);
  @$pb.TagNumber(427938596)
  void clearAuthorizationendpoint() => $_clearField(427938596);
}

class CreateConnectionRequest extends $pb.GeneratedMessage {
  factory CreateConnectionRequest({
    $core.String? description,
    CreateConnectionAuthRequestParameters? authparameters,
    $core.String? name,
    ConnectionAuthorizationType? authorizationtype,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (authparameters != null) result.authparameters = authparameters;
    if (name != null) result.name = name;
    if (authorizationtype != null) result.authorizationtype = authorizationtype;
    return result;
  }

  CreateConnectionRequest._();

  factory CreateConnectionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateConnectionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateConnectionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOM<CreateConnectionAuthRequestParameters>(
        258276552, _omitFieldNames ? '' : 'authparameters',
        subBuilder: CreateConnectionAuthRequestParameters.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<ConnectionAuthorizationType>(
        481976511, _omitFieldNames ? '' : 'authorizationtype',
        enumValues: ConnectionAuthorizationType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionRequest copyWith(
          void Function(CreateConnectionRequest) updates) =>
      super.copyWith((message) => updates(message as CreateConnectionRequest))
          as CreateConnectionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateConnectionRequest create() => CreateConnectionRequest._();
  @$core.override
  CreateConnectionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateConnectionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateConnectionRequest>(create);
  static CreateConnectionRequest? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(258276552)
  CreateConnectionAuthRequestParameters get authparameters => $_getN(1);
  @$pb.TagNumber(258276552)
  set authparameters(CreateConnectionAuthRequestParameters value) =>
      $_setField(258276552, value);
  @$pb.TagNumber(258276552)
  $core.bool hasAuthparameters() => $_has(1);
  @$pb.TagNumber(258276552)
  void clearAuthparameters() => $_clearField(258276552);
  @$pb.TagNumber(258276552)
  CreateConnectionAuthRequestParameters ensureAuthparameters() => $_ensure(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(481976511)
  ConnectionAuthorizationType get authorizationtype => $_getN(3);
  @$pb.TagNumber(481976511)
  set authorizationtype(ConnectionAuthorizationType value) =>
      $_setField(481976511, value);
  @$pb.TagNumber(481976511)
  $core.bool hasAuthorizationtype() => $_has(3);
  @$pb.TagNumber(481976511)
  void clearAuthorizationtype() => $_clearField(481976511);
}

class CreateConnectionResponse extends $pb.GeneratedMessage {
  factory CreateConnectionResponse({
    $core.String? creationtime,
    $core.String? connectionarn,
    $core.String? lastmodifiedtime,
    ConnectionState? connectionstate,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (connectionstate != null) result.connectionstate = connectionstate;
    return result;
  }

  CreateConnectionResponse._();

  factory CreateConnectionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateConnectionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateConnectionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aE<ConnectionState>(404323675, _omitFieldNames ? '' : 'connectionstate',
        enumValues: ConnectionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateConnectionResponse copyWith(
          void Function(CreateConnectionResponse) updates) =>
      super.copyWith((message) => updates(message as CreateConnectionResponse))
          as CreateConnectionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateConnectionResponse create() => CreateConnectionResponse._();
  @$core.override
  CreateConnectionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateConnectionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateConnectionResponse>(create);
  static CreateConnectionResponse? _defaultInstance;

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(0);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(0);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(1);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(1);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(2);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(2);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(404323675)
  ConnectionState get connectionstate => $_getN(3);
  @$pb.TagNumber(404323675)
  set connectionstate(ConnectionState value) => $_setField(404323675, value);
  @$pb.TagNumber(404323675)
  $core.bool hasConnectionstate() => $_has(3);
  @$pb.TagNumber(404323675)
  void clearConnectionstate() => $_clearField(404323675);
}

class CreateEventBusRequest extends $pb.GeneratedMessage {
  factory CreateEventBusRequest({
    $core.String? name,
    $core.Iterable<Tag>? tags,
    $core.String? eventsourcename,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (tags != null) result.tags.addAll(tags);
    if (eventsourcename != null) result.eventsourcename = eventsourcename;
    return result;
  }

  CreateEventBusRequest._();

  factory CreateEventBusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateEventBusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateEventBusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(427669794, _omitFieldNames ? '' : 'eventsourcename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEventBusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEventBusRequest copyWith(
          void Function(CreateEventBusRequest) updates) =>
      super.copyWith((message) => updates(message as CreateEventBusRequest))
          as CreateEventBusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateEventBusRequest create() => CreateEventBusRequest._();
  @$core.override
  CreateEventBusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateEventBusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateEventBusRequest>(create);
  static CreateEventBusRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);

  @$pb.TagNumber(427669794)
  $core.String get eventsourcename => $_getSZ(2);
  @$pb.TagNumber(427669794)
  set eventsourcename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(427669794)
  $core.bool hasEventsourcename() => $_has(2);
  @$pb.TagNumber(427669794)
  void clearEventsourcename() => $_clearField(427669794);
}

class CreateEventBusResponse extends $pb.GeneratedMessage {
  factory CreateEventBusResponse({
    $core.String? eventbusarn,
  }) {
    final result = create();
    if (eventbusarn != null) result.eventbusarn = eventbusarn;
    return result;
  }

  CreateEventBusResponse._();

  factory CreateEventBusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateEventBusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateEventBusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(515755843, _omitFieldNames ? '' : 'eventbusarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEventBusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateEventBusResponse copyWith(
          void Function(CreateEventBusResponse) updates) =>
      super.copyWith((message) => updates(message as CreateEventBusResponse))
          as CreateEventBusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateEventBusResponse create() => CreateEventBusResponse._();
  @$core.override
  CreateEventBusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateEventBusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateEventBusResponse>(create);
  static CreateEventBusResponse? _defaultInstance;

  @$pb.TagNumber(515755843)
  $core.String get eventbusarn => $_getSZ(0);
  @$pb.TagNumber(515755843)
  set eventbusarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(515755843)
  $core.bool hasEventbusarn() => $_has(0);
  @$pb.TagNumber(515755843)
  void clearEventbusarn() => $_clearField(515755843);
}

class CreatePartnerEventSourceRequest extends $pb.GeneratedMessage {
  factory CreatePartnerEventSourceRequest({
    $core.String? name,
    $core.String? account,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (account != null) result.account = account;
    return result;
  }

  CreatePartnerEventSourceRequest._();

  factory CreatePartnerEventSourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePartnerEventSourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePartnerEventSourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(435725053, _omitFieldNames ? '' : 'account')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePartnerEventSourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePartnerEventSourceRequest copyWith(
          void Function(CreatePartnerEventSourceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as CreatePartnerEventSourceRequest))
          as CreatePartnerEventSourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePartnerEventSourceRequest create() =>
      CreatePartnerEventSourceRequest._();
  @$core.override
  CreatePartnerEventSourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePartnerEventSourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePartnerEventSourceRequest>(
          create);
  static CreatePartnerEventSourceRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(435725053)
  $core.String get account => $_getSZ(1);
  @$pb.TagNumber(435725053)
  set account($core.String value) => $_setString(1, value);
  @$pb.TagNumber(435725053)
  $core.bool hasAccount() => $_has(1);
  @$pb.TagNumber(435725053)
  void clearAccount() => $_clearField(435725053);
}

class CreatePartnerEventSourceResponse extends $pb.GeneratedMessage {
  factory CreatePartnerEventSourceResponse({
    $core.String? eventsourcearn,
  }) {
    final result = create();
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    return result;
  }

  CreatePartnerEventSourceResponse._();

  factory CreatePartnerEventSourceResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreatePartnerEventSourceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreatePartnerEventSourceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePartnerEventSourceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreatePartnerEventSourceResponse copyWith(
          void Function(CreatePartnerEventSourceResponse) updates) =>
      super.copyWith(
              (message) => updates(message as CreatePartnerEventSourceResponse))
          as CreatePartnerEventSourceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreatePartnerEventSourceResponse create() =>
      CreatePartnerEventSourceResponse._();
  @$core.override
  CreatePartnerEventSourceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreatePartnerEventSourceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreatePartnerEventSourceResponse>(
          create);
  static CreatePartnerEventSourceResponse? _defaultInstance;

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(0);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(0);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);
}

class DeactivateEventSourceRequest extends $pb.GeneratedMessage {
  factory DeactivateEventSourceRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DeactivateEventSourceRequest._();

  factory DeactivateEventSourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeactivateEventSourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeactivateEventSourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeactivateEventSourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeactivateEventSourceRequest copyWith(
          void Function(DeactivateEventSourceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeactivateEventSourceRequest))
          as DeactivateEventSourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeactivateEventSourceRequest create() =>
      DeactivateEventSourceRequest._();
  @$core.override
  DeactivateEventSourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeactivateEventSourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeactivateEventSourceRequest>(create);
  static DeactivateEventSourceRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class DeauthorizeConnectionRequest extends $pb.GeneratedMessage {
  factory DeauthorizeConnectionRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DeauthorizeConnectionRequest._();

  factory DeauthorizeConnectionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeauthorizeConnectionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeauthorizeConnectionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeauthorizeConnectionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeauthorizeConnectionRequest copyWith(
          void Function(DeauthorizeConnectionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeauthorizeConnectionRequest))
          as DeauthorizeConnectionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeauthorizeConnectionRequest create() =>
      DeauthorizeConnectionRequest._();
  @$core.override
  DeauthorizeConnectionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeauthorizeConnectionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeauthorizeConnectionRequest>(create);
  static DeauthorizeConnectionRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DeauthorizeConnectionResponse extends $pb.GeneratedMessage {
  factory DeauthorizeConnectionResponse({
    $core.String? creationtime,
    $core.String? connectionarn,
    $core.String? lastmodifiedtime,
    $core.String? lastauthorizedtime,
    ConnectionState? connectionstate,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (lastauthorizedtime != null)
      result.lastauthorizedtime = lastauthorizedtime;
    if (connectionstate != null) result.connectionstate = connectionstate;
    return result;
  }

  DeauthorizeConnectionResponse._();

  factory DeauthorizeConnectionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeauthorizeConnectionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeauthorizeConnectionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(250638066, _omitFieldNames ? '' : 'lastauthorizedtime')
    ..aE<ConnectionState>(404323675, _omitFieldNames ? '' : 'connectionstate',
        enumValues: ConnectionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeauthorizeConnectionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeauthorizeConnectionResponse copyWith(
          void Function(DeauthorizeConnectionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeauthorizeConnectionResponse))
          as DeauthorizeConnectionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeauthorizeConnectionResponse create() =>
      DeauthorizeConnectionResponse._();
  @$core.override
  DeauthorizeConnectionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeauthorizeConnectionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeauthorizeConnectionResponse>(create);
  static DeauthorizeConnectionResponse? _defaultInstance;

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(0);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(0);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(1);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(1);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(2);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(2);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(250638066)
  $core.String get lastauthorizedtime => $_getSZ(3);
  @$pb.TagNumber(250638066)
  set lastauthorizedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(250638066)
  $core.bool hasLastauthorizedtime() => $_has(3);
  @$pb.TagNumber(250638066)
  void clearLastauthorizedtime() => $_clearField(250638066);

  @$pb.TagNumber(404323675)
  ConnectionState get connectionstate => $_getN(4);
  @$pb.TagNumber(404323675)
  set connectionstate(ConnectionState value) => $_setField(404323675, value);
  @$pb.TagNumber(404323675)
  $core.bool hasConnectionstate() => $_has(4);
  @$pb.TagNumber(404323675)
  void clearConnectionstate() => $_clearField(404323675);
}

class DeleteApiDestinationRequest extends $pb.GeneratedMessage {
  factory DeleteApiDestinationRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DeleteApiDestinationRequest._();

  factory DeleteApiDestinationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteApiDestinationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteApiDestinationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteApiDestinationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteApiDestinationRequest copyWith(
          void Function(DeleteApiDestinationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteApiDestinationRequest))
          as DeleteApiDestinationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteApiDestinationRequest create() =>
      DeleteApiDestinationRequest._();
  @$core.override
  DeleteApiDestinationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteApiDestinationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteApiDestinationRequest>(create);
  static DeleteApiDestinationRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DeleteApiDestinationResponse extends $pb.GeneratedMessage {
  factory DeleteApiDestinationResponse() => create();

  DeleteApiDestinationResponse._();

  factory DeleteApiDestinationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteApiDestinationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteApiDestinationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteApiDestinationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteApiDestinationResponse copyWith(
          void Function(DeleteApiDestinationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteApiDestinationResponse))
          as DeleteApiDestinationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteApiDestinationResponse create() =>
      DeleteApiDestinationResponse._();
  @$core.override
  DeleteApiDestinationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteApiDestinationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteApiDestinationResponse>(create);
  static DeleteApiDestinationResponse? _defaultInstance;
}

class DeleteArchiveRequest extends $pb.GeneratedMessage {
  factory DeleteArchiveRequest({
    $core.String? archivename,
  }) {
    final result = create();
    if (archivename != null) result.archivename = archivename;
    return result;
  }

  DeleteArchiveRequest._();

  factory DeleteArchiveRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteArchiveRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteArchiveRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(88048487, _omitFieldNames ? '' : 'archivename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteArchiveRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteArchiveRequest copyWith(void Function(DeleteArchiveRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteArchiveRequest))
          as DeleteArchiveRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteArchiveRequest create() => DeleteArchiveRequest._();
  @$core.override
  DeleteArchiveRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteArchiveRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteArchiveRequest>(create);
  static DeleteArchiveRequest? _defaultInstance;

  @$pb.TagNumber(88048487)
  $core.String get archivename => $_getSZ(0);
  @$pb.TagNumber(88048487)
  set archivename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88048487)
  $core.bool hasArchivename() => $_has(0);
  @$pb.TagNumber(88048487)
  void clearArchivename() => $_clearField(88048487);
}

class DeleteArchiveResponse extends $pb.GeneratedMessage {
  factory DeleteArchiveResponse() => create();

  DeleteArchiveResponse._();

  factory DeleteArchiveResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteArchiveResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteArchiveResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteArchiveResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteArchiveResponse copyWith(
          void Function(DeleteArchiveResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteArchiveResponse))
          as DeleteArchiveResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteArchiveResponse create() => DeleteArchiveResponse._();
  @$core.override
  DeleteArchiveResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteArchiveResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteArchiveResponse>(create);
  static DeleteArchiveResponse? _defaultInstance;
}

class DeleteConnectionRequest extends $pb.GeneratedMessage {
  factory DeleteConnectionRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DeleteConnectionRequest._();

  factory DeleteConnectionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteConnectionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteConnectionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteConnectionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteConnectionRequest copyWith(
          void Function(DeleteConnectionRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteConnectionRequest))
          as DeleteConnectionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteConnectionRequest create() => DeleteConnectionRequest._();
  @$core.override
  DeleteConnectionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteConnectionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteConnectionRequest>(create);
  static DeleteConnectionRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DeleteConnectionResponse extends $pb.GeneratedMessage {
  factory DeleteConnectionResponse({
    $core.String? creationtime,
    $core.String? connectionarn,
    $core.String? lastmodifiedtime,
    $core.String? lastauthorizedtime,
    ConnectionState? connectionstate,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (lastauthorizedtime != null)
      result.lastauthorizedtime = lastauthorizedtime;
    if (connectionstate != null) result.connectionstate = connectionstate;
    return result;
  }

  DeleteConnectionResponse._();

  factory DeleteConnectionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteConnectionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteConnectionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(250638066, _omitFieldNames ? '' : 'lastauthorizedtime')
    ..aE<ConnectionState>(404323675, _omitFieldNames ? '' : 'connectionstate',
        enumValues: ConnectionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteConnectionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteConnectionResponse copyWith(
          void Function(DeleteConnectionResponse) updates) =>
      super.copyWith((message) => updates(message as DeleteConnectionResponse))
          as DeleteConnectionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteConnectionResponse create() => DeleteConnectionResponse._();
  @$core.override
  DeleteConnectionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteConnectionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteConnectionResponse>(create);
  static DeleteConnectionResponse? _defaultInstance;

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(0);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(0);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(1);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(1);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(2);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(2);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(250638066)
  $core.String get lastauthorizedtime => $_getSZ(3);
  @$pb.TagNumber(250638066)
  set lastauthorizedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(250638066)
  $core.bool hasLastauthorizedtime() => $_has(3);
  @$pb.TagNumber(250638066)
  void clearLastauthorizedtime() => $_clearField(250638066);

  @$pb.TagNumber(404323675)
  ConnectionState get connectionstate => $_getN(4);
  @$pb.TagNumber(404323675)
  set connectionstate(ConnectionState value) => $_setField(404323675, value);
  @$pb.TagNumber(404323675)
  $core.bool hasConnectionstate() => $_has(4);
  @$pb.TagNumber(404323675)
  void clearConnectionstate() => $_clearField(404323675);
}

class DeleteEventBusRequest extends $pb.GeneratedMessage {
  factory DeleteEventBusRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DeleteEventBusRequest._();

  factory DeleteEventBusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteEventBusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteEventBusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteEventBusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteEventBusRequest copyWith(
          void Function(DeleteEventBusRequest) updates) =>
      super.copyWith((message) => updates(message as DeleteEventBusRequest))
          as DeleteEventBusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteEventBusRequest create() => DeleteEventBusRequest._();
  @$core.override
  DeleteEventBusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteEventBusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteEventBusRequest>(create);
  static DeleteEventBusRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DeletePartnerEventSourceRequest extends $pb.GeneratedMessage {
  factory DeletePartnerEventSourceRequest({
    $core.String? name,
    $core.String? account,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (account != null) result.account = account;
    return result;
  }

  DeletePartnerEventSourceRequest._();

  factory DeletePartnerEventSourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeletePartnerEventSourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeletePartnerEventSourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(435725053, _omitFieldNames ? '' : 'account')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePartnerEventSourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeletePartnerEventSourceRequest copyWith(
          void Function(DeletePartnerEventSourceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DeletePartnerEventSourceRequest))
          as DeletePartnerEventSourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeletePartnerEventSourceRequest create() =>
      DeletePartnerEventSourceRequest._();
  @$core.override
  DeletePartnerEventSourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeletePartnerEventSourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeletePartnerEventSourceRequest>(
          create);
  static DeletePartnerEventSourceRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(435725053)
  $core.String get account => $_getSZ(1);
  @$pb.TagNumber(435725053)
  set account($core.String value) => $_setString(1, value);
  @$pb.TagNumber(435725053)
  $core.bool hasAccount() => $_has(1);
  @$pb.TagNumber(435725053)
  void clearAccount() => $_clearField(435725053);
}

class DeleteRuleRequest extends $pb.GeneratedMessage {
  factory DeleteRuleRequest({
    $core.String? name,
    $core.String? eventbusname,
    $core.bool? force,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (force != null) result.force = force;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aOB(526711333, _omitFieldNames ? '' : 'force')
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

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(1);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(1);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(526711333)
  $core.bool get force => $_getBF(2);
  @$pb.TagNumber(526711333)
  set force($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(526711333)
  $core.bool hasForce() => $_has(2);
  @$pb.TagNumber(526711333)
  void clearForce() => $_clearField(526711333);
}

class DescribeApiDestinationRequest extends $pb.GeneratedMessage {
  factory DescribeApiDestinationRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DescribeApiDestinationRequest._();

  factory DescribeApiDestinationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeApiDestinationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeApiDestinationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeApiDestinationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeApiDestinationRequest copyWith(
          void Function(DescribeApiDestinationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeApiDestinationRequest))
          as DescribeApiDestinationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeApiDestinationRequest create() =>
      DescribeApiDestinationRequest._();
  @$core.override
  DescribeApiDestinationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeApiDestinationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeApiDestinationRequest>(create);
  static DescribeApiDestinationRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DescribeApiDestinationResponse extends $pb.GeneratedMessage {
  factory DescribeApiDestinationResponse({
    ApiDestinationState? apidestinationstate,
    $core.String? apidestinationarn,
    $core.String? creationtime,
    $core.String? description,
    $core.String? connectionarn,
    $core.String? lastmodifiedtime,
    $core.String? name,
    $core.int? invocationratelimitpersecond,
    ApiDestinationHttpMethod? httpmethod,
    $core.String? invocationendpoint,
  }) {
    final result = create();
    if (apidestinationstate != null)
      result.apidestinationstate = apidestinationstate;
    if (apidestinationarn != null) result.apidestinationarn = apidestinationarn;
    if (creationtime != null) result.creationtime = creationtime;
    if (description != null) result.description = description;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (name != null) result.name = name;
    if (invocationratelimitpersecond != null)
      result.invocationratelimitpersecond = invocationratelimitpersecond;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (invocationendpoint != null)
      result.invocationendpoint = invocationendpoint;
    return result;
  }

  DescribeApiDestinationResponse._();

  factory DescribeApiDestinationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeApiDestinationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeApiDestinationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aE<ApiDestinationState>(
        13153343, _omitFieldNames ? '' : 'apidestinationstate',
        enumValues: ApiDestinationState.values)
    ..aOS(90996885, _omitFieldNames ? '' : 'apidestinationarn')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aI(295327816, _omitFieldNames ? '' : 'invocationratelimitpersecond')
    ..aE<ApiDestinationHttpMethod>(
        398394961, _omitFieldNames ? '' : 'httpmethod',
        enumValues: ApiDestinationHttpMethod.values)
    ..aOS(411800759, _omitFieldNames ? '' : 'invocationendpoint')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeApiDestinationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeApiDestinationResponse copyWith(
          void Function(DescribeApiDestinationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeApiDestinationResponse))
          as DescribeApiDestinationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeApiDestinationResponse create() =>
      DescribeApiDestinationResponse._();
  @$core.override
  DescribeApiDestinationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeApiDestinationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeApiDestinationResponse>(create);
  static DescribeApiDestinationResponse? _defaultInstance;

  @$pb.TagNumber(13153343)
  ApiDestinationState get apidestinationstate => $_getN(0);
  @$pb.TagNumber(13153343)
  set apidestinationstate(ApiDestinationState value) =>
      $_setField(13153343, value);
  @$pb.TagNumber(13153343)
  $core.bool hasApidestinationstate() => $_has(0);
  @$pb.TagNumber(13153343)
  void clearApidestinationstate() => $_clearField(13153343);

  @$pb.TagNumber(90996885)
  $core.String get apidestinationarn => $_getSZ(1);
  @$pb.TagNumber(90996885)
  set apidestinationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(90996885)
  $core.bool hasApidestinationarn() => $_has(1);
  @$pb.TagNumber(90996885)
  void clearApidestinationarn() => $_clearField(90996885);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(4);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(4);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(5);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(5);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(6);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(6, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(6);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(295327816)
  $core.int get invocationratelimitpersecond => $_getIZ(7);
  @$pb.TagNumber(295327816)
  set invocationratelimitpersecond($core.int value) =>
      $_setSignedInt32(7, value);
  @$pb.TagNumber(295327816)
  $core.bool hasInvocationratelimitpersecond() => $_has(7);
  @$pb.TagNumber(295327816)
  void clearInvocationratelimitpersecond() => $_clearField(295327816);

  @$pb.TagNumber(398394961)
  ApiDestinationHttpMethod get httpmethod => $_getN(8);
  @$pb.TagNumber(398394961)
  set httpmethod(ApiDestinationHttpMethod value) =>
      $_setField(398394961, value);
  @$pb.TagNumber(398394961)
  $core.bool hasHttpmethod() => $_has(8);
  @$pb.TagNumber(398394961)
  void clearHttpmethod() => $_clearField(398394961);

  @$pb.TagNumber(411800759)
  $core.String get invocationendpoint => $_getSZ(9);
  @$pb.TagNumber(411800759)
  set invocationendpoint($core.String value) => $_setString(9, value);
  @$pb.TagNumber(411800759)
  $core.bool hasInvocationendpoint() => $_has(9);
  @$pb.TagNumber(411800759)
  void clearInvocationendpoint() => $_clearField(411800759);
}

class DescribeArchiveRequest extends $pb.GeneratedMessage {
  factory DescribeArchiveRequest({
    $core.String? archivename,
  }) {
    final result = create();
    if (archivename != null) result.archivename = archivename;
    return result;
  }

  DescribeArchiveRequest._();

  factory DescribeArchiveRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeArchiveRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeArchiveRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(88048487, _omitFieldNames ? '' : 'archivename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeArchiveRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeArchiveRequest copyWith(
          void Function(DescribeArchiveRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeArchiveRequest))
          as DescribeArchiveRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeArchiveRequest create() => DescribeArchiveRequest._();
  @$core.override
  DescribeArchiveRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeArchiveRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeArchiveRequest>(create);
  static DescribeArchiveRequest? _defaultInstance;

  @$pb.TagNumber(88048487)
  $core.String get archivename => $_getSZ(0);
  @$pb.TagNumber(88048487)
  set archivename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88048487)
  $core.bool hasArchivename() => $_has(0);
  @$pb.TagNumber(88048487)
  void clearArchivename() => $_clearField(88048487);
}

class DescribeArchiveResponse extends $pb.GeneratedMessage {
  factory DescribeArchiveResponse({
    $core.String? archivearn,
    $core.String? archivename,
    $core.String? creationtime,
    $core.String? description,
    $fixnum.Int64? eventcount,
    $core.String? eventpattern,
    $core.int? retentiondays,
    $core.String? eventsourcearn,
    $core.String? statereason,
    $fixnum.Int64? sizebytes,
    ArchiveState? state,
  }) {
    final result = create();
    if (archivearn != null) result.archivearn = archivearn;
    if (archivename != null) result.archivename = archivename;
    if (creationtime != null) result.creationtime = creationtime;
    if (description != null) result.description = description;
    if (eventcount != null) result.eventcount = eventcount;
    if (eventpattern != null) result.eventpattern = eventpattern;
    if (retentiondays != null) result.retentiondays = retentiondays;
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    if (statereason != null) result.statereason = statereason;
    if (sizebytes != null) result.sizebytes = sizebytes;
    if (state != null) result.state = state;
    return result;
  }

  DescribeArchiveResponse._();

  factory DescribeArchiveResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeArchiveResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeArchiveResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(56866685, _omitFieldNames ? '' : 'archivearn')
    ..aOS(88048487, _omitFieldNames ? '' : 'archivename')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aInt64(128612839, _omitFieldNames ? '' : 'eventcount')
    ..aOS(233487232, _omitFieldNames ? '' : 'eventpattern')
    ..aI(267894223, _omitFieldNames ? '' : 'retentiondays')
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aInt64(486677664, _omitFieldNames ? '' : 'sizebytes')
    ..aE<ArchiveState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ArchiveState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeArchiveResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeArchiveResponse copyWith(
          void Function(DescribeArchiveResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeArchiveResponse))
          as DescribeArchiveResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeArchiveResponse create() => DescribeArchiveResponse._();
  @$core.override
  DescribeArchiveResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeArchiveResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeArchiveResponse>(create);
  static DescribeArchiveResponse? _defaultInstance;

  @$pb.TagNumber(56866685)
  $core.String get archivearn => $_getSZ(0);
  @$pb.TagNumber(56866685)
  set archivearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56866685)
  $core.bool hasArchivearn() => $_has(0);
  @$pb.TagNumber(56866685)
  void clearArchivearn() => $_clearField(56866685);

  @$pb.TagNumber(88048487)
  $core.String get archivename => $_getSZ(1);
  @$pb.TagNumber(88048487)
  set archivename($core.String value) => $_setString(1, value);
  @$pb.TagNumber(88048487)
  $core.bool hasArchivename() => $_has(1);
  @$pb.TagNumber(88048487)
  void clearArchivename() => $_clearField(88048487);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(128612839)
  $fixnum.Int64 get eventcount => $_getI64(4);
  @$pb.TagNumber(128612839)
  set eventcount($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(128612839)
  $core.bool hasEventcount() => $_has(4);
  @$pb.TagNumber(128612839)
  void clearEventcount() => $_clearField(128612839);

  @$pb.TagNumber(233487232)
  $core.String get eventpattern => $_getSZ(5);
  @$pb.TagNumber(233487232)
  set eventpattern($core.String value) => $_setString(5, value);
  @$pb.TagNumber(233487232)
  $core.bool hasEventpattern() => $_has(5);
  @$pb.TagNumber(233487232)
  void clearEventpattern() => $_clearField(233487232);

  @$pb.TagNumber(267894223)
  $core.int get retentiondays => $_getIZ(6);
  @$pb.TagNumber(267894223)
  set retentiondays($core.int value) => $_setSignedInt32(6, value);
  @$pb.TagNumber(267894223)
  $core.bool hasRetentiondays() => $_has(6);
  @$pb.TagNumber(267894223)
  void clearRetentiondays() => $_clearField(267894223);

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(7);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(7, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(7);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(8);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(8, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(8);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(486677664)
  $fixnum.Int64 get sizebytes => $_getI64(9);
  @$pb.TagNumber(486677664)
  set sizebytes($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(486677664)
  $core.bool hasSizebytes() => $_has(9);
  @$pb.TagNumber(486677664)
  void clearSizebytes() => $_clearField(486677664);

  @$pb.TagNumber(502047895)
  ArchiveState get state => $_getN(10);
  @$pb.TagNumber(502047895)
  set state(ArchiveState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(10);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class DescribeConnectionRequest extends $pb.GeneratedMessage {
  factory DescribeConnectionRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DescribeConnectionRequest._();

  factory DescribeConnectionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeConnectionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeConnectionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeConnectionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeConnectionRequest copyWith(
          void Function(DescribeConnectionRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeConnectionRequest))
          as DescribeConnectionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeConnectionRequest create() => DescribeConnectionRequest._();
  @$core.override
  DescribeConnectionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeConnectionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeConnectionRequest>(create);
  static DescribeConnectionRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DescribeConnectionResponse extends $pb.GeneratedMessage {
  factory DescribeConnectionResponse({
    $core.String? creationtime,
    $core.String? description,
    $core.String? connectionarn,
    $core.String? lastmodifiedtime,
    $core.String? secretarn,
    $core.String? lastauthorizedtime,
    ConnectionAuthResponseParameters? authparameters,
    $core.String? name,
    $core.String? statereason,
    ConnectionState? connectionstate,
    ConnectionAuthorizationType? authorizationtype,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (description != null) result.description = description;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (secretarn != null) result.secretarn = secretarn;
    if (lastauthorizedtime != null)
      result.lastauthorizedtime = lastauthorizedtime;
    if (authparameters != null) result.authparameters = authparameters;
    if (name != null) result.name = name;
    if (statereason != null) result.statereason = statereason;
    if (connectionstate != null) result.connectionstate = connectionstate;
    if (authorizationtype != null) result.authorizationtype = authorizationtype;
    return result;
  }

  DescribeConnectionResponse._();

  factory DescribeConnectionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeConnectionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeConnectionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(241012025, _omitFieldNames ? '' : 'secretarn')
    ..aOS(250638066, _omitFieldNames ? '' : 'lastauthorizedtime')
    ..aOM<ConnectionAuthResponseParameters>(
        258276552, _omitFieldNames ? '' : 'authparameters',
        subBuilder: ConnectionAuthResponseParameters.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aE<ConnectionState>(404323675, _omitFieldNames ? '' : 'connectionstate',
        enumValues: ConnectionState.values)
    ..aE<ConnectionAuthorizationType>(
        481976511, _omitFieldNames ? '' : 'authorizationtype',
        enumValues: ConnectionAuthorizationType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeConnectionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeConnectionResponse copyWith(
          void Function(DescribeConnectionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeConnectionResponse))
          as DescribeConnectionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeConnectionResponse create() => DescribeConnectionResponse._();
  @$core.override
  DescribeConnectionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeConnectionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeConnectionResponse>(create);
  static DescribeConnectionResponse? _defaultInstance;

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(0);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(0);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(2);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(2);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(3);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(3);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(241012025)
  $core.String get secretarn => $_getSZ(4);
  @$pb.TagNumber(241012025)
  set secretarn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(241012025)
  $core.bool hasSecretarn() => $_has(4);
  @$pb.TagNumber(241012025)
  void clearSecretarn() => $_clearField(241012025);

  @$pb.TagNumber(250638066)
  $core.String get lastauthorizedtime => $_getSZ(5);
  @$pb.TagNumber(250638066)
  set lastauthorizedtime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(250638066)
  $core.bool hasLastauthorizedtime() => $_has(5);
  @$pb.TagNumber(250638066)
  void clearLastauthorizedtime() => $_clearField(250638066);

  @$pb.TagNumber(258276552)
  ConnectionAuthResponseParameters get authparameters => $_getN(6);
  @$pb.TagNumber(258276552)
  set authparameters(ConnectionAuthResponseParameters value) =>
      $_setField(258276552, value);
  @$pb.TagNumber(258276552)
  $core.bool hasAuthparameters() => $_has(6);
  @$pb.TagNumber(258276552)
  void clearAuthparameters() => $_clearField(258276552);
  @$pb.TagNumber(258276552)
  ConnectionAuthResponseParameters ensureAuthparameters() => $_ensure(6);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(7);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(7, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(7);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(8);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(8, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(8);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(404323675)
  ConnectionState get connectionstate => $_getN(9);
  @$pb.TagNumber(404323675)
  set connectionstate(ConnectionState value) => $_setField(404323675, value);
  @$pb.TagNumber(404323675)
  $core.bool hasConnectionstate() => $_has(9);
  @$pb.TagNumber(404323675)
  void clearConnectionstate() => $_clearField(404323675);

  @$pb.TagNumber(481976511)
  ConnectionAuthorizationType get authorizationtype => $_getN(10);
  @$pb.TagNumber(481976511)
  set authorizationtype(ConnectionAuthorizationType value) =>
      $_setField(481976511, value);
  @$pb.TagNumber(481976511)
  $core.bool hasAuthorizationtype() => $_has(10);
  @$pb.TagNumber(481976511)
  void clearAuthorizationtype() => $_clearField(481976511);
}

class DescribeEventBusRequest extends $pb.GeneratedMessage {
  factory DescribeEventBusRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DescribeEventBusRequest._();

  factory DescribeEventBusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeEventBusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeEventBusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEventBusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEventBusRequest copyWith(
          void Function(DescribeEventBusRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeEventBusRequest))
          as DescribeEventBusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeEventBusRequest create() => DescribeEventBusRequest._();
  @$core.override
  DescribeEventBusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeEventBusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeEventBusRequest>(create);
  static DescribeEventBusRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DescribeEventBusResponse extends $pb.GeneratedMessage {
  factory DescribeEventBusResponse({
    $core.String? name,
    $core.String? arn,
    $core.String? policy,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    if (policy != null) result.policy = policy;
    return result;
  }

  DescribeEventBusResponse._();

  factory DescribeEventBusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeEventBusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeEventBusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEventBusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEventBusResponse copyWith(
          void Function(DescribeEventBusResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeEventBusResponse))
          as DescribeEventBusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeEventBusResponse create() => DescribeEventBusResponse._();
  @$core.override
  DescribeEventBusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeEventBusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeEventBusResponse>(create);
  static DescribeEventBusResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(2);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(2);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class DescribeEventSourceRequest extends $pb.GeneratedMessage {
  factory DescribeEventSourceRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DescribeEventSourceRequest._();

  factory DescribeEventSourceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeEventSourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeEventSourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEventSourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEventSourceRequest copyWith(
          void Function(DescribeEventSourceRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeEventSourceRequest))
          as DescribeEventSourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeEventSourceRequest create() => DescribeEventSourceRequest._();
  @$core.override
  DescribeEventSourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeEventSourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeEventSourceRequest>(create);
  static DescribeEventSourceRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DescribeEventSourceResponse extends $pb.GeneratedMessage {
  factory DescribeEventSourceResponse({
    $core.String? createdby,
    $core.String? expirationtime,
    $core.String? creationtime,
    $core.String? name,
    $core.String? arn,
    EventSourceState? state,
  }) {
    final result = create();
    if (createdby != null) result.createdby = createdby;
    if (expirationtime != null) result.expirationtime = expirationtime;
    if (creationtime != null) result.creationtime = creationtime;
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    if (state != null) result.state = state;
    return result;
  }

  DescribeEventSourceResponse._();

  factory DescribeEventSourceResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeEventSourceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeEventSourceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(73491847, _omitFieldNames ? '' : 'createdby')
    ..aOS(93473378, _omitFieldNames ? '' : 'expirationtime')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aE<EventSourceState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: EventSourceState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEventSourceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeEventSourceResponse copyWith(
          void Function(DescribeEventSourceResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeEventSourceResponse))
          as DescribeEventSourceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeEventSourceResponse create() =>
      DescribeEventSourceResponse._();
  @$core.override
  DescribeEventSourceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeEventSourceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeEventSourceResponse>(create);
  static DescribeEventSourceResponse? _defaultInstance;

  @$pb.TagNumber(73491847)
  $core.String get createdby => $_getSZ(0);
  @$pb.TagNumber(73491847)
  set createdby($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73491847)
  $core.bool hasCreatedby() => $_has(0);
  @$pb.TagNumber(73491847)
  void clearCreatedby() => $_clearField(73491847);

  @$pb.TagNumber(93473378)
  $core.String get expirationtime => $_getSZ(1);
  @$pb.TagNumber(93473378)
  set expirationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(93473378)
  $core.bool hasExpirationtime() => $_has(1);
  @$pb.TagNumber(93473378)
  void clearExpirationtime() => $_clearField(93473378);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(502047895)
  EventSourceState get state => $_getN(5);
  @$pb.TagNumber(502047895)
  set state(EventSourceState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class DescribePartnerEventSourceRequest extends $pb.GeneratedMessage {
  factory DescribePartnerEventSourceRequest({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DescribePartnerEventSourceRequest._();

  factory DescribePartnerEventSourceRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribePartnerEventSourceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribePartnerEventSourceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribePartnerEventSourceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribePartnerEventSourceRequest copyWith(
          void Function(DescribePartnerEventSourceRequest) updates) =>
      super.copyWith((message) =>
              updates(message as DescribePartnerEventSourceRequest))
          as DescribePartnerEventSourceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribePartnerEventSourceRequest create() =>
      DescribePartnerEventSourceRequest._();
  @$core.override
  DescribePartnerEventSourceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribePartnerEventSourceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribePartnerEventSourceRequest>(
          create);
  static DescribePartnerEventSourceRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DescribePartnerEventSourceResponse extends $pb.GeneratedMessage {
  factory DescribePartnerEventSourceResponse({
    $core.String? name,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    return result;
  }

  DescribePartnerEventSourceResponse._();

  factory DescribePartnerEventSourceResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribePartnerEventSourceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribePartnerEventSourceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribePartnerEventSourceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribePartnerEventSourceResponse copyWith(
          void Function(DescribePartnerEventSourceResponse) updates) =>
      super.copyWith((message) =>
              updates(message as DescribePartnerEventSourceResponse))
          as DescribePartnerEventSourceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribePartnerEventSourceResponse create() =>
      DescribePartnerEventSourceResponse._();
  @$core.override
  DescribePartnerEventSourceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribePartnerEventSourceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribePartnerEventSourceResponse>(
          create);
  static DescribePartnerEventSourceResponse? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class DescribeReplayRequest extends $pb.GeneratedMessage {
  factory DescribeReplayRequest({
    $core.String? replayname,
  }) {
    final result = create();
    if (replayname != null) result.replayname = replayname;
    return result;
  }

  DescribeReplayRequest._();

  factory DescribeReplayRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeReplayRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeReplayRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(442173850, _omitFieldNames ? '' : 'replayname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeReplayRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeReplayRequest copyWith(
          void Function(DescribeReplayRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeReplayRequest))
          as DescribeReplayRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeReplayRequest create() => DescribeReplayRequest._();
  @$core.override
  DescribeReplayRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeReplayRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeReplayRequest>(create);
  static DescribeReplayRequest? _defaultInstance;

  @$pb.TagNumber(442173850)
  $core.String get replayname => $_getSZ(0);
  @$pb.TagNumber(442173850)
  set replayname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(442173850)
  $core.bool hasReplayname() => $_has(0);
  @$pb.TagNumber(442173850)
  void clearReplayname() => $_clearField(442173850);
}

class DescribeReplayResponse extends $pb.GeneratedMessage {
  factory DescribeReplayResponse({
    $core.String? eventlastreplayedtime,
    $core.String? eventendtime,
    $core.String? eventstarttime,
    $core.String? description,
    $core.String? replayendtime,
    $core.String? eventsourcearn,
    $core.String? replayarn,
    $core.String? statereason,
    $core.String? replayname,
    ReplayDestination? destination,
    ReplayState? state,
    $core.String? replaystarttime,
  }) {
    final result = create();
    if (eventlastreplayedtime != null)
      result.eventlastreplayedtime = eventlastreplayedtime;
    if (eventendtime != null) result.eventendtime = eventendtime;
    if (eventstarttime != null) result.eventstarttime = eventstarttime;
    if (description != null) result.description = description;
    if (replayendtime != null) result.replayendtime = replayendtime;
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    if (replayarn != null) result.replayarn = replayarn;
    if (statereason != null) result.statereason = statereason;
    if (replayname != null) result.replayname = replayname;
    if (destination != null) result.destination = destination;
    if (state != null) result.state = state;
    if (replaystarttime != null) result.replaystarttime = replaystarttime;
    return result;
  }

  DescribeReplayResponse._();

  factory DescribeReplayResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeReplayResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeReplayResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(7759073, _omitFieldNames ? '' : 'eventlastreplayedtime')
    ..aOS(30229582, _omitFieldNames ? '' : 'eventendtime')
    ..aOS(103680757, _omitFieldNames ? '' : 'eventstarttime')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(304261199, _omitFieldNames ? '' : 'replayendtime')
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..aOS(361946526, _omitFieldNames ? '' : 'replayarn')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aOS(442173850, _omitFieldNames ? '' : 'replayname')
    ..aOM<ReplayDestination>(457443680, _omitFieldNames ? '' : 'destination',
        subBuilder: ReplayDestination.create)
    ..aE<ReplayState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ReplayState.values)
    ..aOS(503580492, _omitFieldNames ? '' : 'replaystarttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeReplayResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeReplayResponse copyWith(
          void Function(DescribeReplayResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeReplayResponse))
          as DescribeReplayResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeReplayResponse create() => DescribeReplayResponse._();
  @$core.override
  DescribeReplayResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeReplayResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeReplayResponse>(create);
  static DescribeReplayResponse? _defaultInstance;

  @$pb.TagNumber(7759073)
  $core.String get eventlastreplayedtime => $_getSZ(0);
  @$pb.TagNumber(7759073)
  set eventlastreplayedtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7759073)
  $core.bool hasEventlastreplayedtime() => $_has(0);
  @$pb.TagNumber(7759073)
  void clearEventlastreplayedtime() => $_clearField(7759073);

  @$pb.TagNumber(30229582)
  $core.String get eventendtime => $_getSZ(1);
  @$pb.TagNumber(30229582)
  set eventendtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(30229582)
  $core.bool hasEventendtime() => $_has(1);
  @$pb.TagNumber(30229582)
  void clearEventendtime() => $_clearField(30229582);

  @$pb.TagNumber(103680757)
  $core.String get eventstarttime => $_getSZ(2);
  @$pb.TagNumber(103680757)
  set eventstarttime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103680757)
  $core.bool hasEventstarttime() => $_has(2);
  @$pb.TagNumber(103680757)
  void clearEventstarttime() => $_clearField(103680757);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(304261199)
  $core.String get replayendtime => $_getSZ(4);
  @$pb.TagNumber(304261199)
  set replayendtime($core.String value) => $_setString(4, value);
  @$pb.TagNumber(304261199)
  $core.bool hasReplayendtime() => $_has(4);
  @$pb.TagNumber(304261199)
  void clearReplayendtime() => $_clearField(304261199);

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(5);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(5);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);

  @$pb.TagNumber(361946526)
  $core.String get replayarn => $_getSZ(6);
  @$pb.TagNumber(361946526)
  set replayarn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(361946526)
  $core.bool hasReplayarn() => $_has(6);
  @$pb.TagNumber(361946526)
  void clearReplayarn() => $_clearField(361946526);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(7);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(7, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(7);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(442173850)
  $core.String get replayname => $_getSZ(8);
  @$pb.TagNumber(442173850)
  set replayname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(442173850)
  $core.bool hasReplayname() => $_has(8);
  @$pb.TagNumber(442173850)
  void clearReplayname() => $_clearField(442173850);

  @$pb.TagNumber(457443680)
  ReplayDestination get destination => $_getN(9);
  @$pb.TagNumber(457443680)
  set destination(ReplayDestination value) => $_setField(457443680, value);
  @$pb.TagNumber(457443680)
  $core.bool hasDestination() => $_has(9);
  @$pb.TagNumber(457443680)
  void clearDestination() => $_clearField(457443680);
  @$pb.TagNumber(457443680)
  ReplayDestination ensureDestination() => $_ensure(9);

  @$pb.TagNumber(502047895)
  ReplayState get state => $_getN(10);
  @$pb.TagNumber(502047895)
  set state(ReplayState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(10);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);

  @$pb.TagNumber(503580492)
  $core.String get replaystarttime => $_getSZ(11);
  @$pb.TagNumber(503580492)
  set replaystarttime($core.String value) => $_setString(11, value);
  @$pb.TagNumber(503580492)
  $core.bool hasReplaystarttime() => $_has(11);
  @$pb.TagNumber(503580492)
  void clearReplaystarttime() => $_clearField(503580492);
}

class DescribeRuleRequest extends $pb.GeneratedMessage {
  factory DescribeRuleRequest({
    $core.String? name,
    $core.String? eventbusname,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (eventbusname != null) result.eventbusname = eventbusname;
    return result;
  }

  DescribeRuleRequest._();

  factory DescribeRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeRuleRequest copyWith(void Function(DescribeRuleRequest) updates) =>
      super.copyWith((message) => updates(message as DescribeRuleRequest))
          as DescribeRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeRuleRequest create() => DescribeRuleRequest._();
  @$core.override
  DescribeRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeRuleRequest>(create);
  static DescribeRuleRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(1);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(1);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);
}

class DescribeRuleResponse extends $pb.GeneratedMessage {
  factory DescribeRuleResponse({
    $core.String? createdby,
    $core.String? description,
    $core.String? eventpattern,
    $core.String? name,
    $core.String? rolearn,
    $core.String? arn,
    $core.String? scheduleexpression,
    $core.String? eventbusname,
    $core.String? managedby,
    RuleState? state,
  }) {
    final result = create();
    if (createdby != null) result.createdby = createdby;
    if (description != null) result.description = description;
    if (eventpattern != null) result.eventpattern = eventpattern;
    if (name != null) result.name = name;
    if (rolearn != null) result.rolearn = rolearn;
    if (arn != null) result.arn = arn;
    if (scheduleexpression != null)
      result.scheduleexpression = scheduleexpression;
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (managedby != null) result.managedby = managedby;
    if (state != null) result.state = state;
    return result;
  }

  DescribeRuleResponse._();

  factory DescribeRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(73491847, _omitFieldNames ? '' : 'createdby')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(233487232, _omitFieldNames ? '' : 'eventpattern')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(446089471, _omitFieldNames ? '' : 'scheduleexpression')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aOS(455511232, _omitFieldNames ? '' : 'managedby')
    ..aE<RuleState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: RuleState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeRuleResponse copyWith(void Function(DescribeRuleResponse) updates) =>
      super.copyWith((message) => updates(message as DescribeRuleResponse))
          as DescribeRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeRuleResponse create() => DescribeRuleResponse._();
  @$core.override
  DescribeRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeRuleResponse>(create);
  static DescribeRuleResponse? _defaultInstance;

  @$pb.TagNumber(73491847)
  $core.String get createdby => $_getSZ(0);
  @$pb.TagNumber(73491847)
  set createdby($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73491847)
  $core.bool hasCreatedby() => $_has(0);
  @$pb.TagNumber(73491847)
  void clearCreatedby() => $_clearField(73491847);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(233487232)
  $core.String get eventpattern => $_getSZ(2);
  @$pb.TagNumber(233487232)
  set eventpattern($core.String value) => $_setString(2, value);
  @$pb.TagNumber(233487232)
  $core.bool hasEventpattern() => $_has(2);
  @$pb.TagNumber(233487232)
  void clearEventpattern() => $_clearField(233487232);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(4);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(4);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(5);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(5);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(446089471)
  $core.String get scheduleexpression => $_getSZ(6);
  @$pb.TagNumber(446089471)
  set scheduleexpression($core.String value) => $_setString(6, value);
  @$pb.TagNumber(446089471)
  $core.bool hasScheduleexpression() => $_has(6);
  @$pb.TagNumber(446089471)
  void clearScheduleexpression() => $_clearField(446089471);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(7);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(7, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(7);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(455511232)
  $core.String get managedby => $_getSZ(8);
  @$pb.TagNumber(455511232)
  set managedby($core.String value) => $_setString(8, value);
  @$pb.TagNumber(455511232)
  $core.bool hasManagedby() => $_has(8);
  @$pb.TagNumber(455511232)
  void clearManagedby() => $_clearField(455511232);

  @$pb.TagNumber(502047895)
  RuleState get state => $_getN(9);
  @$pb.TagNumber(502047895)
  set state(RuleState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(9);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class DisableRuleRequest extends $pb.GeneratedMessage {
  factory DisableRuleRequest({
    $core.String? name,
    $core.String? eventbusname,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (eventbusname != null) result.eventbusname = eventbusname;
    return result;
  }

  DisableRuleRequest._();

  factory DisableRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableRuleRequest copyWith(void Function(DisableRuleRequest) updates) =>
      super.copyWith((message) => updates(message as DisableRuleRequest))
          as DisableRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableRuleRequest create() => DisableRuleRequest._();
  @$core.override
  DisableRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableRuleRequest>(create);
  static DisableRuleRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(1);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(1);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);
}

class EcsParameters extends $pb.GeneratedMessage {
  factory EcsParameters({
    $core.Iterable<PlacementStrategy>? placementstrategy,
    $core.String? taskdefinitionarn,
    $core.String? group,
    $core.String? platformversion,
    $core.bool? enableecsmanagedtags,
    LaunchType? launchtype,
    NetworkConfiguration? networkconfiguration,
    $core.Iterable<PlacementConstraint>? placementconstraints,
    $core.Iterable<CapacityProviderStrategyItem>? capacityproviderstrategy,
    $core.String? referenceid,
    $core.Iterable<Tag>? tags,
    $core.int? taskcount,
    PropagateTags? propagatetags,
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
    if (tags != null) result.tags.addAll(tags);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<PlacementStrategy>(
        25036678, _omitFieldNames ? '' : 'placementstrategy',
        subBuilder: PlacementStrategy.create)
    ..aOS(82234477, _omitFieldNames ? '' : 'taskdefinitionarn')
    ..aOS(91525165, _omitFieldNames ? '' : 'group')
    ..aOS(139924287, _omitFieldNames ? '' : 'platformversion')
    ..aOB(146161174, _omitFieldNames ? '' : 'enableecsmanagedtags')
    ..aE<LaunchType>(184333335, _omitFieldNames ? '' : 'launchtype',
        enumValues: LaunchType.values)
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
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aI(398407508, _omitFieldNames ? '' : 'taskcount')
    ..aE<PropagateTags>(405557622, _omitFieldNames ? '' : 'propagatetags',
        enumValues: PropagateTags.values)
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
  LaunchType get launchtype => $_getN(5);
  @$pb.TagNumber(184333335)
  set launchtype(LaunchType value) => $_setField(184333335, value);
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
  $pb.PbList<Tag> get tags => $_getList(10);

  @$pb.TagNumber(398407508)
  $core.int get taskcount => $_getIZ(11);
  @$pb.TagNumber(398407508)
  set taskcount($core.int value) => $_setSignedInt32(11, value);
  @$pb.TagNumber(398407508)
  $core.bool hasTaskcount() => $_has(11);
  @$pb.TagNumber(398407508)
  void clearTaskcount() => $_clearField(398407508);

  @$pb.TagNumber(405557622)
  PropagateTags get propagatetags => $_getN(12);
  @$pb.TagNumber(405557622)
  set propagatetags(PropagateTags value) => $_setField(405557622, value);
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

class EnableRuleRequest extends $pb.GeneratedMessage {
  factory EnableRuleRequest({
    $core.String? name,
    $core.String? eventbusname,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (eventbusname != null) result.eventbusname = eventbusname;
    return result;
  }

  EnableRuleRequest._();

  factory EnableRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableRuleRequest copyWith(void Function(EnableRuleRequest) updates) =>
      super.copyWith((message) => updates(message as EnableRuleRequest))
          as EnableRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableRuleRequest create() => EnableRuleRequest._();
  @$core.override
  EnableRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableRuleRequest>(create);
  static EnableRuleRequest? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(1);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(1);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);
}

class EventBus extends $pb.GeneratedMessage {
  factory EventBus({
    $core.String? name,
    $core.String? arn,
    $core.String? policy,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    if (policy != null) result.policy = policy;
    return result;
  }

  EventBus._();

  factory EventBus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventBus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventBus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventBus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventBus copyWith(void Function(EventBus) updates) =>
      super.copyWith((message) => updates(message as EventBus)) as EventBus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventBus create() => EventBus._();
  @$core.override
  EventBus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventBus getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<EventBus>(create);
  static EventBus? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(2);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(2);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class EventSource extends $pb.GeneratedMessage {
  factory EventSource({
    $core.String? createdby,
    $core.String? expirationtime,
    $core.String? creationtime,
    $core.String? name,
    $core.String? arn,
    EventSourceState? state,
  }) {
    final result = create();
    if (createdby != null) result.createdby = createdby;
    if (expirationtime != null) result.expirationtime = expirationtime;
    if (creationtime != null) result.creationtime = creationtime;
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    if (state != null) result.state = state;
    return result;
  }

  EventSource._();

  factory EventSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EventSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EventSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(73491847, _omitFieldNames ? '' : 'createdby')
    ..aOS(93473378, _omitFieldNames ? '' : 'expirationtime')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aE<EventSourceState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: EventSourceState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EventSource copyWith(void Function(EventSource) updates) =>
      super.copyWith((message) => updates(message as EventSource))
          as EventSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EventSource create() => EventSource._();
  @$core.override
  EventSource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EventSource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EventSource>(create);
  static EventSource? _defaultInstance;

  @$pb.TagNumber(73491847)
  $core.String get createdby => $_getSZ(0);
  @$pb.TagNumber(73491847)
  set createdby($core.String value) => $_setString(0, value);
  @$pb.TagNumber(73491847)
  $core.bool hasCreatedby() => $_has(0);
  @$pb.TagNumber(73491847)
  void clearCreatedby() => $_clearField(73491847);

  @$pb.TagNumber(93473378)
  $core.String get expirationtime => $_getSZ(1);
  @$pb.TagNumber(93473378)
  set expirationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(93473378)
  $core.bool hasExpirationtime() => $_has(1);
  @$pb.TagNumber(93473378)
  void clearExpirationtime() => $_clearField(93473378);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(502047895)
  EventSourceState get state => $_getN(5);
  @$pb.TagNumber(502047895)
  set state(EventSourceState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class HttpParameters extends $pb.GeneratedMessage {
  factory HttpParameters({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        headerparameters,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        querystringparameters,
    $core.Iterable<$core.String>? pathparametervalues,
  }) {
    final result = create();
    if (headerparameters != null)
      result.headerparameters.addEntries(headerparameters);
    if (querystringparameters != null)
      result.querystringparameters.addEntries(querystringparameters);
    if (pathparametervalues != null)
      result.pathparametervalues.addAll(pathparametervalues);
    return result;
  }

  HttpParameters._();

  factory HttpParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HttpParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HttpParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        148944757, _omitFieldNames ? '' : 'headerparameters',
        entryClassName: 'HttpParameters.HeaderparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('events'))
    ..m<$core.String, $core.String>(
        258002263, _omitFieldNames ? '' : 'querystringparameters',
        entryClassName: 'HttpParameters.QuerystringparametersEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('events'))
    ..pPS(489440856, _omitFieldNames ? '' : 'pathparametervalues')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HttpParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HttpParameters copyWith(void Function(HttpParameters) updates) =>
      super.copyWith((message) => updates(message as HttpParameters))
          as HttpParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HttpParameters create() => HttpParameters._();
  @$core.override
  HttpParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HttpParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HttpParameters>(create);
  static HttpParameters? _defaultInstance;

  @$pb.TagNumber(148944757)
  $pb.PbMap<$core.String, $core.String> get headerparameters => $_getMap(0);

  @$pb.TagNumber(258002263)
  $pb.PbMap<$core.String, $core.String> get querystringparameters =>
      $_getMap(1);

  @$pb.TagNumber(489440856)
  $pb.PbList<$core.String> get pathparametervalues => $_getList(2);
}

class IllegalStatusException extends $pb.GeneratedMessage {
  factory IllegalStatusException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  IllegalStatusException._();

  factory IllegalStatusException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IllegalStatusException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IllegalStatusException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IllegalStatusException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IllegalStatusException copyWith(
          void Function(IllegalStatusException) updates) =>
      super.copyWith((message) => updates(message as IllegalStatusException))
          as IllegalStatusException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IllegalStatusException create() => IllegalStatusException._();
  @$core.override
  IllegalStatusException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IllegalStatusException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IllegalStatusException>(create);
  static IllegalStatusException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InputTransformer extends $pb.GeneratedMessage {
  factory InputTransformer({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? inputpathsmap,
    $core.String? inputtemplate,
  }) {
    final result = create();
    if (inputpathsmap != null) result.inputpathsmap.addEntries(inputpathsmap);
    if (inputtemplate != null) result.inputtemplate = inputtemplate;
    return result;
  }

  InputTransformer._();

  factory InputTransformer.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InputTransformer.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InputTransformer',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        154989658, _omitFieldNames ? '' : 'inputpathsmap',
        entryClassName: 'InputTransformer.InputpathsmapEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('events'))
    ..aOS(505971626, _omitFieldNames ? '' : 'inputtemplate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InputTransformer clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InputTransformer copyWith(void Function(InputTransformer) updates) =>
      super.copyWith((message) => updates(message as InputTransformer))
          as InputTransformer;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InputTransformer create() => InputTransformer._();
  @$core.override
  InputTransformer createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InputTransformer getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InputTransformer>(create);
  static InputTransformer? _defaultInstance;

  @$pb.TagNumber(154989658)
  $pb.PbMap<$core.String, $core.String> get inputpathsmap => $_getMap(0);

  @$pb.TagNumber(505971626)
  $core.String get inputtemplate => $_getSZ(1);
  @$pb.TagNumber(505971626)
  set inputtemplate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(505971626)
  $core.bool hasInputtemplate() => $_has(1);
  @$pb.TagNumber(505971626)
  void clearInputtemplate() => $_clearField(505971626);
}

class InternalException extends $pb.GeneratedMessage {
  factory InternalException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InternalException._();

  factory InternalException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InternalException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InternalException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalException copyWith(void Function(InternalException) updates) =>
      super.copyWith((message) => updates(message as InternalException))
          as InternalException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InternalException create() => InternalException._();
  @$core.override
  InternalException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InternalException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InternalException>(create);
  static InternalException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidEventPatternException extends $pb.GeneratedMessage {
  factory InvalidEventPatternException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidEventPatternException._();

  factory InvalidEventPatternException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidEventPatternException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidEventPatternException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventPatternException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidEventPatternException copyWith(
          void Function(InvalidEventPatternException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidEventPatternException))
          as InvalidEventPatternException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidEventPatternException create() =>
      InvalidEventPatternException._();
  @$core.override
  InvalidEventPatternException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidEventPatternException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidEventPatternException>(create);
  static InvalidEventPatternException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class KinesisParameters extends $pb.GeneratedMessage {
  factory KinesisParameters({
    $core.String? partitionkeypath,
  }) {
    final result = create();
    if (partitionkeypath != null) result.partitionkeypath = partitionkeypath;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(458722216, _omitFieldNames ? '' : 'partitionkeypath')
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

  @$pb.TagNumber(458722216)
  $core.String get partitionkeypath => $_getSZ(0);
  @$pb.TagNumber(458722216)
  set partitionkeypath($core.String value) => $_setString(0, value);
  @$pb.TagNumber(458722216)
  $core.bool hasPartitionkeypath() => $_has(0);
  @$pb.TagNumber(458722216)
  void clearPartitionkeypath() => $_clearField(458722216);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class ListApiDestinationsRequest extends $pb.GeneratedMessage {
  factory ListApiDestinationsRequest({
    $core.String? connectionarn,
    $core.String? nexttoken,
    $core.String? nameprefix,
    $core.int? limit,
  }) {
    final result = create();
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListApiDestinationsRequest._();

  factory ListApiDestinationsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListApiDestinationsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListApiDestinationsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListApiDestinationsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListApiDestinationsRequest copyWith(
          void Function(ListApiDestinationsRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListApiDestinationsRequest))
          as ListApiDestinationsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListApiDestinationsRequest create() => ListApiDestinationsRequest._();
  @$core.override
  ListApiDestinationsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListApiDestinationsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListApiDestinationsRequest>(create);
  static ListApiDestinationsRequest? _defaultInstance;

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(0);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(0);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(2);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(2, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(2);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListApiDestinationsResponse extends $pb.GeneratedMessage {
  factory ListApiDestinationsResponse({
    $core.String? nexttoken,
    $core.Iterable<ApiDestination>? apidestinations,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (apidestinations != null) result.apidestinations.addAll(apidestinations);
    return result;
  }

  ListApiDestinationsResponse._();

  factory ListApiDestinationsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListApiDestinationsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListApiDestinationsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ApiDestination>(321414539, _omitFieldNames ? '' : 'apidestinations',
        subBuilder: ApiDestination.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListApiDestinationsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListApiDestinationsResponse copyWith(
          void Function(ListApiDestinationsResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListApiDestinationsResponse))
          as ListApiDestinationsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListApiDestinationsResponse create() =>
      ListApiDestinationsResponse._();
  @$core.override
  ListApiDestinationsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListApiDestinationsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListApiDestinationsResponse>(create);
  static ListApiDestinationsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(321414539)
  $pb.PbList<ApiDestination> get apidestinations => $_getList(1);
}

class ListArchivesRequest extends $pb.GeneratedMessage {
  factory ListArchivesRequest({
    $core.String? nexttoken,
    $core.String? eventsourcearn,
    $core.String? nameprefix,
    $core.int? limit,
    ArchiveState? state,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (limit != null) result.limit = limit;
    if (state != null) result.state = state;
    return result;
  }

  ListArchivesRequest._();

  factory ListArchivesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListArchivesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListArchivesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aE<ArchiveState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ArchiveState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListArchivesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListArchivesRequest copyWith(void Function(ListArchivesRequest) updates) =>
      super.copyWith((message) => updates(message as ListArchivesRequest))
          as ListArchivesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListArchivesRequest create() => ListArchivesRequest._();
  @$core.override
  ListArchivesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListArchivesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListArchivesRequest>(create);
  static ListArchivesRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(1);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(1);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(2);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(2, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(2);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(502047895)
  ArchiveState get state => $_getN(4);
  @$pb.TagNumber(502047895)
  set state(ArchiveState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(4);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class ListArchivesResponse extends $pb.GeneratedMessage {
  factory ListArchivesResponse({
    $core.Iterable<Archive>? archives,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (archives != null) result.archives.addAll(archives);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListArchivesResponse._();

  factory ListArchivesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListArchivesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListArchivesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<Archive>(73167779, _omitFieldNames ? '' : 'archives',
        subBuilder: Archive.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListArchivesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListArchivesResponse copyWith(void Function(ListArchivesResponse) updates) =>
      super.copyWith((message) => updates(message as ListArchivesResponse))
          as ListArchivesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListArchivesResponse create() => ListArchivesResponse._();
  @$core.override
  ListArchivesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListArchivesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListArchivesResponse>(create);
  static ListArchivesResponse? _defaultInstance;

  @$pb.TagNumber(73167779)
  $pb.PbList<Archive> get archives => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListConnectionsRequest extends $pb.GeneratedMessage {
  factory ListConnectionsRequest({
    $core.String? nexttoken,
    $core.String? nameprefix,
    ConnectionState? connectionstate,
    $core.int? limit,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (connectionstate != null) result.connectionstate = connectionstate;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListConnectionsRequest._();

  factory ListConnectionsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListConnectionsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListConnectionsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aE<ConnectionState>(404323675, _omitFieldNames ? '' : 'connectionstate',
        enumValues: ConnectionState.values)
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListConnectionsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListConnectionsRequest copyWith(
          void Function(ListConnectionsRequest) updates) =>
      super.copyWith((message) => updates(message as ListConnectionsRequest))
          as ListConnectionsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListConnectionsRequest create() => ListConnectionsRequest._();
  @$core.override
  ListConnectionsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListConnectionsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListConnectionsRequest>(create);
  static ListConnectionsRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(1);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(1);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(404323675)
  ConnectionState get connectionstate => $_getN(2);
  @$pb.TagNumber(404323675)
  set connectionstate(ConnectionState value) => $_setField(404323675, value);
  @$pb.TagNumber(404323675)
  $core.bool hasConnectionstate() => $_has(2);
  @$pb.TagNumber(404323675)
  void clearConnectionstate() => $_clearField(404323675);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListConnectionsResponse extends $pb.GeneratedMessage {
  factory ListConnectionsResponse({
    $core.Iterable<Connection>? connections,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (connections != null) result.connections.addAll(connections);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListConnectionsResponse._();

  factory ListConnectionsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListConnectionsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListConnectionsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<Connection>(98252351, _omitFieldNames ? '' : 'connections',
        subBuilder: Connection.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListConnectionsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListConnectionsResponse copyWith(
          void Function(ListConnectionsResponse) updates) =>
      super.copyWith((message) => updates(message as ListConnectionsResponse))
          as ListConnectionsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListConnectionsResponse create() => ListConnectionsResponse._();
  @$core.override
  ListConnectionsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListConnectionsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListConnectionsResponse>(create);
  static ListConnectionsResponse? _defaultInstance;

  @$pb.TagNumber(98252351)
  $pb.PbList<Connection> get connections => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListEventBusesRequest extends $pb.GeneratedMessage {
  factory ListEventBusesRequest({
    $core.String? nexttoken,
    $core.String? nameprefix,
    $core.int? limit,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListEventBusesRequest._();

  factory ListEventBusesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEventBusesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEventBusesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventBusesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventBusesRequest copyWith(
          void Function(ListEventBusesRequest) updates) =>
      super.copyWith((message) => updates(message as ListEventBusesRequest))
          as ListEventBusesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEventBusesRequest create() => ListEventBusesRequest._();
  @$core.override
  ListEventBusesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEventBusesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListEventBusesRequest>(create);
  static ListEventBusesRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(1);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(1);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListEventBusesResponse extends $pb.GeneratedMessage {
  factory ListEventBusesResponse({
    $core.Iterable<EventBus>? eventbuses,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (eventbuses != null) result.eventbuses.addAll(eventbuses);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  ListEventBusesResponse._();

  factory ListEventBusesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEventBusesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEventBusesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<EventBus>(125356616, _omitFieldNames ? '' : 'eventbuses',
        subBuilder: EventBus.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventBusesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventBusesResponse copyWith(
          void Function(ListEventBusesResponse) updates) =>
      super.copyWith((message) => updates(message as ListEventBusesResponse))
          as ListEventBusesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEventBusesResponse create() => ListEventBusesResponse._();
  @$core.override
  ListEventBusesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEventBusesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListEventBusesResponse>(create);
  static ListEventBusesResponse? _defaultInstance;

  @$pb.TagNumber(125356616)
  $pb.PbList<EventBus> get eventbuses => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class ListEventSourcesRequest extends $pb.GeneratedMessage {
  factory ListEventSourcesRequest({
    $core.String? nexttoken,
    $core.String? nameprefix,
    $core.int? limit,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListEventSourcesRequest._();

  factory ListEventSourcesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEventSourcesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEventSourcesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventSourcesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventSourcesRequest copyWith(
          void Function(ListEventSourcesRequest) updates) =>
      super.copyWith((message) => updates(message as ListEventSourcesRequest))
          as ListEventSourcesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEventSourcesRequest create() => ListEventSourcesRequest._();
  @$core.override
  ListEventSourcesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEventSourcesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListEventSourcesRequest>(create);
  static ListEventSourcesRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(1);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(1);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListEventSourcesResponse extends $pb.GeneratedMessage {
  factory ListEventSourcesResponse({
    $core.String? nexttoken,
    $core.Iterable<EventSource>? eventsources,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (eventsources != null) result.eventsources.addAll(eventsources);
    return result;
  }

  ListEventSourcesResponse._();

  factory ListEventSourcesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListEventSourcesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListEventSourcesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<EventSource>(368674668, _omitFieldNames ? '' : 'eventsources',
        subBuilder: EventSource.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventSourcesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListEventSourcesResponse copyWith(
          void Function(ListEventSourcesResponse) updates) =>
      super.copyWith((message) => updates(message as ListEventSourcesResponse))
          as ListEventSourcesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListEventSourcesResponse create() => ListEventSourcesResponse._();
  @$core.override
  ListEventSourcesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListEventSourcesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListEventSourcesResponse>(create);
  static ListEventSourcesResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(368674668)
  $pb.PbList<EventSource> get eventsources => $_getList(1);
}

class ListPartnerEventSourceAccountsRequest extends $pb.GeneratedMessage {
  factory ListPartnerEventSourceAccountsRequest({
    $core.String? nexttoken,
    $core.int? limit,
    $core.String? eventsourcename,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (limit != null) result.limit = limit;
    if (eventsourcename != null) result.eventsourcename = eventsourcename;
    return result;
  }

  ListPartnerEventSourceAccountsRequest._();

  factory ListPartnerEventSourceAccountsRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPartnerEventSourceAccountsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPartnerEventSourceAccountsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(427669794, _omitFieldNames ? '' : 'eventsourcename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPartnerEventSourceAccountsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPartnerEventSourceAccountsRequest copyWith(
          void Function(ListPartnerEventSourceAccountsRequest) updates) =>
      super.copyWith((message) =>
              updates(message as ListPartnerEventSourceAccountsRequest))
          as ListPartnerEventSourceAccountsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPartnerEventSourceAccountsRequest create() =>
      ListPartnerEventSourceAccountsRequest._();
  @$core.override
  ListPartnerEventSourceAccountsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPartnerEventSourceAccountsRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListPartnerEventSourceAccountsRequest>(create);
  static ListPartnerEventSourceAccountsRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(427669794)
  $core.String get eventsourcename => $_getSZ(2);
  @$pb.TagNumber(427669794)
  set eventsourcename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(427669794)
  $core.bool hasEventsourcename() => $_has(2);
  @$pb.TagNumber(427669794)
  void clearEventsourcename() => $_clearField(427669794);
}

class ListPartnerEventSourceAccountsResponse extends $pb.GeneratedMessage {
  factory ListPartnerEventSourceAccountsResponse({
    $core.String? nexttoken,
    $core.Iterable<PartnerEventSourceAccount>? partnereventsourceaccounts,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (partnereventsourceaccounts != null)
      result.partnereventsourceaccounts.addAll(partnereventsourceaccounts);
    return result;
  }

  ListPartnerEventSourceAccountsResponse._();

  factory ListPartnerEventSourceAccountsResponse.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPartnerEventSourceAccountsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPartnerEventSourceAccountsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<PartnerEventSourceAccount>(
        334415295, _omitFieldNames ? '' : 'partnereventsourceaccounts',
        subBuilder: PartnerEventSourceAccount.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPartnerEventSourceAccountsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPartnerEventSourceAccountsResponse copyWith(
          void Function(ListPartnerEventSourceAccountsResponse) updates) =>
      super.copyWith((message) =>
              updates(message as ListPartnerEventSourceAccountsResponse))
          as ListPartnerEventSourceAccountsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPartnerEventSourceAccountsResponse create() =>
      ListPartnerEventSourceAccountsResponse._();
  @$core.override
  ListPartnerEventSourceAccountsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPartnerEventSourceAccountsResponse getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          ListPartnerEventSourceAccountsResponse>(create);
  static ListPartnerEventSourceAccountsResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(334415295)
  $pb.PbList<PartnerEventSourceAccount> get partnereventsourceaccounts =>
      $_getList(1);
}

class ListPartnerEventSourcesRequest extends $pb.GeneratedMessage {
  factory ListPartnerEventSourcesRequest({
    $core.String? nexttoken,
    $core.String? nameprefix,
    $core.int? limit,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (limit != null) result.limit = limit;
    return result;
  }

  ListPartnerEventSourcesRequest._();

  factory ListPartnerEventSourcesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPartnerEventSourcesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPartnerEventSourcesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPartnerEventSourcesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPartnerEventSourcesRequest copyWith(
          void Function(ListPartnerEventSourcesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListPartnerEventSourcesRequest))
          as ListPartnerEventSourcesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPartnerEventSourcesRequest create() =>
      ListPartnerEventSourcesRequest._();
  @$core.override
  ListPartnerEventSourcesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPartnerEventSourcesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPartnerEventSourcesRequest>(create);
  static ListPartnerEventSourcesRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(1);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(1);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);
}

class ListPartnerEventSourcesResponse extends $pb.GeneratedMessage {
  factory ListPartnerEventSourcesResponse({
    $core.String? nexttoken,
    $core.Iterable<PartnerEventSource>? partnereventsources,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (partnereventsources != null)
      result.partnereventsources.addAll(partnereventsources);
    return result;
  }

  ListPartnerEventSourcesResponse._();

  factory ListPartnerEventSourcesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPartnerEventSourcesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPartnerEventSourcesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<PartnerEventSource>(
        502828176, _omitFieldNames ? '' : 'partnereventsources',
        subBuilder: PartnerEventSource.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPartnerEventSourcesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPartnerEventSourcesResponse copyWith(
          void Function(ListPartnerEventSourcesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListPartnerEventSourcesResponse))
          as ListPartnerEventSourcesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPartnerEventSourcesResponse create() =>
      ListPartnerEventSourcesResponse._();
  @$core.override
  ListPartnerEventSourcesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPartnerEventSourcesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPartnerEventSourcesResponse>(
          create);
  static ListPartnerEventSourcesResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(502828176)
  $pb.PbList<PartnerEventSource> get partnereventsources => $_getList(1);
}

class ListReplaysRequest extends $pb.GeneratedMessage {
  factory ListReplaysRequest({
    $core.String? nexttoken,
    $core.String? eventsourcearn,
    $core.String? nameprefix,
    $core.int? limit,
    ReplayState? state,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (limit != null) result.limit = limit;
    if (state != null) result.state = state;
    return result;
  }

  ListReplaysRequest._();

  factory ListReplaysRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListReplaysRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListReplaysRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aE<ReplayState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ReplayState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListReplaysRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListReplaysRequest copyWith(void Function(ListReplaysRequest) updates) =>
      super.copyWith((message) => updates(message as ListReplaysRequest))
          as ListReplaysRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListReplaysRequest create() => ListReplaysRequest._();
  @$core.override
  ListReplaysRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListReplaysRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListReplaysRequest>(create);
  static ListReplaysRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(1);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(1);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(2);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(2, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(2);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(3);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(3);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(502047895)
  ReplayState get state => $_getN(4);
  @$pb.TagNumber(502047895)
  set state(ReplayState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(4);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class ListReplaysResponse extends $pb.GeneratedMessage {
  factory ListReplaysResponse({
    $core.String? nexttoken,
    $core.Iterable<Replay>? replays,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (replays != null) result.replays.addAll(replays);
    return result;
  }

  ListReplaysResponse._();

  factory ListReplaysResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListReplaysResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListReplaysResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Replay>(312313364, _omitFieldNames ? '' : 'replays',
        subBuilder: Replay.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListReplaysResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListReplaysResponse copyWith(void Function(ListReplaysResponse) updates) =>
      super.copyWith((message) => updates(message as ListReplaysResponse))
          as ListReplaysResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListReplaysResponse create() => ListReplaysResponse._();
  @$core.override
  ListReplaysResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListReplaysResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListReplaysResponse>(create);
  static ListReplaysResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(312313364)
  $pb.PbList<Replay> get replays => $_getList(1);
}

class ListRuleNamesByTargetRequest extends $pb.GeneratedMessage {
  factory ListRuleNamesByTargetRequest({
    $core.String? nexttoken,
    $core.String? targetarn,
    $core.int? limit,
    $core.String? eventbusname,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (targetarn != null) result.targetarn = targetarn;
    if (limit != null) result.limit = limit;
    if (eventbusname != null) result.eventbusname = eventbusname;
    return result;
  }

  ListRuleNamesByTargetRequest._();

  factory ListRuleNamesByTargetRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRuleNamesByTargetRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRuleNamesByTargetRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(217664144, _omitFieldNames ? '' : 'targetarn')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRuleNamesByTargetRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRuleNamesByTargetRequest copyWith(
          void Function(ListRuleNamesByTargetRequest) updates) =>
      super.copyWith(
              (message) => updates(message as ListRuleNamesByTargetRequest))
          as ListRuleNamesByTargetRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRuleNamesByTargetRequest create() =>
      ListRuleNamesByTargetRequest._();
  @$core.override
  ListRuleNamesByTargetRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRuleNamesByTargetRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRuleNamesByTargetRequest>(create);
  static ListRuleNamesByTargetRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(217664144)
  $core.String get targetarn => $_getSZ(1);
  @$pb.TagNumber(217664144)
  set targetarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(217664144)
  $core.bool hasTargetarn() => $_has(1);
  @$pb.TagNumber(217664144)
  void clearTargetarn() => $_clearField(217664144);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(3);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(3);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);
}

class ListRuleNamesByTargetResponse extends $pb.GeneratedMessage {
  factory ListRuleNamesByTargetResponse({
    $core.String? nexttoken,
    $core.Iterable<$core.String>? rulenames,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (rulenames != null) result.rulenames.addAll(rulenames);
    return result;
  }

  ListRuleNamesByTargetResponse._();

  factory ListRuleNamesByTargetResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListRuleNamesByTargetResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListRuleNamesByTargetResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPS(267949170, _omitFieldNames ? '' : 'rulenames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRuleNamesByTargetResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListRuleNamesByTargetResponse copyWith(
          void Function(ListRuleNamesByTargetResponse) updates) =>
      super.copyWith(
              (message) => updates(message as ListRuleNamesByTargetResponse))
          as ListRuleNamesByTargetResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListRuleNamesByTargetResponse create() =>
      ListRuleNamesByTargetResponse._();
  @$core.override
  ListRuleNamesByTargetResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListRuleNamesByTargetResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListRuleNamesByTargetResponse>(create);
  static ListRuleNamesByTargetResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(267949170)
  $pb.PbList<$core.String> get rulenames => $_getList(1);
}

class ListRulesRequest extends $pb.GeneratedMessage {
  factory ListRulesRequest({
    $core.String? nexttoken,
    $core.String? nameprefix,
    $core.int? limit,
    $core.String? eventbusname,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (nameprefix != null) result.nameprefix = nameprefix;
    if (limit != null) result.limit = limit;
    if (eventbusname != null) result.eventbusname = eventbusname;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(361707931, _omitFieldNames ? '' : 'nameprefix')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
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

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(361707931)
  $core.String get nameprefix => $_getSZ(1);
  @$pb.TagNumber(361707931)
  set nameprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(361707931)
  $core.bool hasNameprefix() => $_has(1);
  @$pb.TagNumber(361707931)
  void clearNameprefix() => $_clearField(361707931);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(2);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(2);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(3);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(3);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);
}

class ListRulesResponse extends $pb.GeneratedMessage {
  factory ListRulesResponse({
    $core.Iterable<Rule>? rules,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (rules != null) result.rules.addAll(rules);
    if (nexttoken != null) result.nexttoken = nexttoken;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<Rule>(42675585, _omitFieldNames ? '' : 'rules',
        subBuilder: Rule.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
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
  $pb.PbList<Rule> get rules => $_getList(0);

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class ListTargetsByRuleRequest extends $pb.GeneratedMessage {
  factory ListTargetsByRuleRequest({
    $core.String? nexttoken,
    $core.int? limit,
    $core.String? eventbusname,
    $core.String? rule,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (limit != null) result.limit = limit;
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (rule != null) result.rule = rule;
    return result;
  }

  ListTargetsByRuleRequest._();

  factory ListTargetsByRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTargetsByRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTargetsByRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(412502741, _omitFieldNames ? '' : 'limit')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aOS(475696372, _omitFieldNames ? '' : 'rule')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTargetsByRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTargetsByRuleRequest copyWith(
          void Function(ListTargetsByRuleRequest) updates) =>
      super.copyWith((message) => updates(message as ListTargetsByRuleRequest))
          as ListTargetsByRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTargetsByRuleRequest create() => ListTargetsByRuleRequest._();
  @$core.override
  ListTargetsByRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTargetsByRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTargetsByRuleRequest>(create);
  static ListTargetsByRuleRequest? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(412502741)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(412502741)
  set limit($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(412502741)
  $core.bool hasLimit() => $_has(1);
  @$pb.TagNumber(412502741)
  void clearLimit() => $_clearField(412502741);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(2);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(2);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(475696372)
  $core.String get rule => $_getSZ(3);
  @$pb.TagNumber(475696372)
  set rule($core.String value) => $_setString(3, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(3);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);
}

class ListTargetsByRuleResponse extends $pb.GeneratedMessage {
  factory ListTargetsByRuleResponse({
    $core.String? nexttoken,
    $core.Iterable<Target>? targets,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (targets != null) result.targets.addAll(targets);
    return result;
  }

  ListTargetsByRuleResponse._();

  factory ListTargetsByRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTargetsByRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTargetsByRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Target>(262180226, _omitFieldNames ? '' : 'targets',
        subBuilder: Target.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTargetsByRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTargetsByRuleResponse copyWith(
          void Function(ListTargetsByRuleResponse) updates) =>
      super.copyWith((message) => updates(message as ListTargetsByRuleResponse))
          as ListTargetsByRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTargetsByRuleResponse create() => ListTargetsByRuleResponse._();
  @$core.override
  ListTargetsByRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTargetsByRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTargetsByRuleResponse>(create);
  static ListTargetsByRuleResponse? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(262180226)
  $pb.PbList<Target> get targets => $_getList(1);
}

class ManagedRuleException extends $pb.GeneratedMessage {
  factory ManagedRuleException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ManagedRuleException._();

  factory ManagedRuleException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleException copyWith(void Function(ManagedRuleException) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleException))
          as ManagedRuleException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleException create() => ManagedRuleException._();
  @$core.override
  ManagedRuleException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleException>(create);
  static ManagedRuleException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class OperationDisabledException extends $pb.GeneratedMessage {
  factory OperationDisabledException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  OperationDisabledException._();

  factory OperationDisabledException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory OperationDisabledException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'OperationDisabledException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OperationDisabledException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  OperationDisabledException copyWith(
          void Function(OperationDisabledException) updates) =>
      super.copyWith(
              (message) => updates(message as OperationDisabledException))
          as OperationDisabledException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static OperationDisabledException create() => OperationDisabledException._();
  @$core.override
  OperationDisabledException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static OperationDisabledException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<OperationDisabledException>(create);
  static OperationDisabledException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PartnerEventSource extends $pb.GeneratedMessage {
  factory PartnerEventSource({
    $core.String? name,
    $core.String? arn,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (arn != null) result.arn = arn;
    return result;
  }

  PartnerEventSource._();

  factory PartnerEventSource.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PartnerEventSource.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PartnerEventSource',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartnerEventSource clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartnerEventSource copyWith(void Function(PartnerEventSource) updates) =>
      super.copyWith((message) => updates(message as PartnerEventSource))
          as PartnerEventSource;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PartnerEventSource create() => PartnerEventSource._();
  @$core.override
  PartnerEventSource createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PartnerEventSource getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PartnerEventSource>(create);
  static PartnerEventSource? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class PartnerEventSourceAccount extends $pb.GeneratedMessage {
  factory PartnerEventSourceAccount({
    $core.String? expirationtime,
    $core.String? creationtime,
    $core.String? account,
    EventSourceState? state,
  }) {
    final result = create();
    if (expirationtime != null) result.expirationtime = expirationtime;
    if (creationtime != null) result.creationtime = creationtime;
    if (account != null) result.account = account;
    if (state != null) result.state = state;
    return result;
  }

  PartnerEventSourceAccount._();

  factory PartnerEventSourceAccount.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PartnerEventSourceAccount.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PartnerEventSourceAccount',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(93473378, _omitFieldNames ? '' : 'expirationtime')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(435725053, _omitFieldNames ? '' : 'account')
    ..aE<EventSourceState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: EventSourceState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartnerEventSourceAccount clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartnerEventSourceAccount copyWith(
          void Function(PartnerEventSourceAccount) updates) =>
      super.copyWith((message) => updates(message as PartnerEventSourceAccount))
          as PartnerEventSourceAccount;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PartnerEventSourceAccount create() => PartnerEventSourceAccount._();
  @$core.override
  PartnerEventSourceAccount createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PartnerEventSourceAccount getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PartnerEventSourceAccount>(create);
  static PartnerEventSourceAccount? _defaultInstance;

  @$pb.TagNumber(93473378)
  $core.String get expirationtime => $_getSZ(0);
  @$pb.TagNumber(93473378)
  set expirationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(93473378)
  $core.bool hasExpirationtime() => $_has(0);
  @$pb.TagNumber(93473378)
  void clearExpirationtime() => $_clearField(93473378);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(435725053)
  $core.String get account => $_getSZ(2);
  @$pb.TagNumber(435725053)
  set account($core.String value) => $_setString(2, value);
  @$pb.TagNumber(435725053)
  $core.bool hasAccount() => $_has(2);
  @$pb.TagNumber(435725053)
  void clearAccount() => $_clearField(435725053);

  @$pb.TagNumber(502047895)
  EventSourceState get state => $_getN(3);
  @$pb.TagNumber(502047895)
  set state(EventSourceState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(3);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class PlacementConstraint extends $pb.GeneratedMessage {
  factory PlacementConstraint({
    $core.String? expression,
    PlacementConstraintType? type,
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(253079532, _omitFieldNames ? '' : 'expression')
    ..aE<PlacementConstraintType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: PlacementConstraintType.values)
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
  PlacementConstraintType get type => $_getN(1);
  @$pb.TagNumber(287830350)
  set type(PlacementConstraintType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);
}

class PlacementStrategy extends $pb.GeneratedMessage {
  factory PlacementStrategy({
    $core.String? field_125985384,
    PlacementStrategyType? type,
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(125985384, _omitFieldNames ? '' : 'field')
    ..aE<PlacementStrategyType>(287830350, _omitFieldNames ? '' : 'type',
        enumValues: PlacementStrategyType.values)
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
  PlacementStrategyType get type => $_getN(1);
  @$pb.TagNumber(287830350)
  set type(PlacementStrategyType value) => $_setField(287830350, value);
  @$pb.TagNumber(287830350)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(287830350)
  void clearType() => $_clearField(287830350);
}

class PolicyLengthExceededException extends $pb.GeneratedMessage {
  factory PolicyLengthExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  PolicyLengthExceededException._();

  factory PolicyLengthExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PolicyLengthExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PolicyLengthExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PolicyLengthExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PolicyLengthExceededException copyWith(
          void Function(PolicyLengthExceededException) updates) =>
      super.copyWith(
              (message) => updates(message as PolicyLengthExceededException))
          as PolicyLengthExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PolicyLengthExceededException create() =>
      PolicyLengthExceededException._();
  @$core.override
  PolicyLengthExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PolicyLengthExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PolicyLengthExceededException>(create);
  static PolicyLengthExceededException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class PutEventsRequest extends $pb.GeneratedMessage {
  factory PutEventsRequest({
    $core.Iterable<PutEventsRequestEntry>? entries,
  }) {
    final result = create();
    if (entries != null) result.entries.addAll(entries);
    return result;
  }

  PutEventsRequest._();

  factory PutEventsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutEventsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutEventsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<PutEventsRequestEntry>(481075860, _omitFieldNames ? '' : 'entries',
        subBuilder: PutEventsRequestEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventsRequest copyWith(void Function(PutEventsRequest) updates) =>
      super.copyWith((message) => updates(message as PutEventsRequest))
          as PutEventsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutEventsRequest create() => PutEventsRequest._();
  @$core.override
  PutEventsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutEventsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutEventsRequest>(create);
  static PutEventsRequest? _defaultInstance;

  @$pb.TagNumber(481075860)
  $pb.PbList<PutEventsRequestEntry> get entries => $_getList(0);
}

class PutEventsRequestEntry extends $pb.GeneratedMessage {
  factory PutEventsRequestEntry({
    $core.String? detailtype,
    $core.String? source,
    $core.String? traceheader,
    $core.String? detail,
    $core.Iterable<$core.String>? resources,
    $core.String? eventbusname,
    $core.String? time,
  }) {
    final result = create();
    if (detailtype != null) result.detailtype = detailtype;
    if (source != null) result.source = source;
    if (traceheader != null) result.traceheader = traceheader;
    if (detail != null) result.detail = detail;
    if (resources != null) result.resources.addAll(resources);
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (time != null) result.time = time;
    return result;
  }

  PutEventsRequestEntry._();

  factory PutEventsRequestEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutEventsRequestEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutEventsRequestEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(11758589, _omitFieldNames ? '' : 'detailtype')
    ..aOS(31630329, _omitFieldNames ? '' : 'source')
    ..aOS(303628800, _omitFieldNames ? '' : 'traceheader')
    ..aOS(340030389, _omitFieldNames ? '' : 'detail')
    ..pPS(358918291, _omitFieldNames ? '' : 'resources')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aOS(535094277, _omitFieldNames ? '' : 'time')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventsRequestEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventsRequestEntry copyWith(
          void Function(PutEventsRequestEntry) updates) =>
      super.copyWith((message) => updates(message as PutEventsRequestEntry))
          as PutEventsRequestEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutEventsRequestEntry create() => PutEventsRequestEntry._();
  @$core.override
  PutEventsRequestEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutEventsRequestEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutEventsRequestEntry>(create);
  static PutEventsRequestEntry? _defaultInstance;

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

  @$pb.TagNumber(303628800)
  $core.String get traceheader => $_getSZ(2);
  @$pb.TagNumber(303628800)
  set traceheader($core.String value) => $_setString(2, value);
  @$pb.TagNumber(303628800)
  $core.bool hasTraceheader() => $_has(2);
  @$pb.TagNumber(303628800)
  void clearTraceheader() => $_clearField(303628800);

  @$pb.TagNumber(340030389)
  $core.String get detail => $_getSZ(3);
  @$pb.TagNumber(340030389)
  set detail($core.String value) => $_setString(3, value);
  @$pb.TagNumber(340030389)
  $core.bool hasDetail() => $_has(3);
  @$pb.TagNumber(340030389)
  void clearDetail() => $_clearField(340030389);

  @$pb.TagNumber(358918291)
  $pb.PbList<$core.String> get resources => $_getList(4);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(5);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(5, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(5);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(535094277)
  $core.String get time => $_getSZ(6);
  @$pb.TagNumber(535094277)
  set time($core.String value) => $_setString(6, value);
  @$pb.TagNumber(535094277)
  $core.bool hasTime() => $_has(6);
  @$pb.TagNumber(535094277)
  void clearTime() => $_clearField(535094277);
}

class PutEventsResponse extends $pb.GeneratedMessage {
  factory PutEventsResponse({
    $core.int? failedentrycount,
    $core.Iterable<PutEventsResultEntry>? entries,
  }) {
    final result = create();
    if (failedentrycount != null) result.failedentrycount = failedentrycount;
    if (entries != null) result.entries.addAll(entries);
    return result;
  }

  PutEventsResponse._();

  factory PutEventsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutEventsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutEventsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aI(458519576, _omitFieldNames ? '' : 'failedentrycount')
    ..pPM<PutEventsResultEntry>(481075860, _omitFieldNames ? '' : 'entries',
        subBuilder: PutEventsResultEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventsResponse copyWith(void Function(PutEventsResponse) updates) =>
      super.copyWith((message) => updates(message as PutEventsResponse))
          as PutEventsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutEventsResponse create() => PutEventsResponse._();
  @$core.override
  PutEventsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutEventsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutEventsResponse>(create);
  static PutEventsResponse? _defaultInstance;

  @$pb.TagNumber(458519576)
  $core.int get failedentrycount => $_getIZ(0);
  @$pb.TagNumber(458519576)
  set failedentrycount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(458519576)
  $core.bool hasFailedentrycount() => $_has(0);
  @$pb.TagNumber(458519576)
  void clearFailedentrycount() => $_clearField(458519576);

  @$pb.TagNumber(481075860)
  $pb.PbList<PutEventsResultEntry> get entries => $_getList(1);
}

class PutEventsResultEntry extends $pb.GeneratedMessage {
  factory PutEventsResultEntry({
    $core.String? errorcode,
    $core.String? eventid,
    $core.String? errormessage,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (eventid != null) result.eventid = eventid;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  PutEventsResultEntry._();

  factory PutEventsResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutEventsResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutEventsResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(376916819, _omitFieldNames ? '' : 'eventid')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventsResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutEventsResultEntry copyWith(void Function(PutEventsResultEntry) updates) =>
      super.copyWith((message) => updates(message as PutEventsResultEntry))
          as PutEventsResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutEventsResultEntry create() => PutEventsResultEntry._();
  @$core.override
  PutEventsResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutEventsResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutEventsResultEntry>(create);
  static PutEventsResultEntry? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(376916819)
  $core.String get eventid => $_getSZ(1);
  @$pb.TagNumber(376916819)
  set eventid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(376916819)
  $core.bool hasEventid() => $_has(1);
  @$pb.TagNumber(376916819)
  void clearEventid() => $_clearField(376916819);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(2);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(2, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(2);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class PutPartnerEventsRequest extends $pb.GeneratedMessage {
  factory PutPartnerEventsRequest({
    $core.Iterable<PutPartnerEventsRequestEntry>? entries,
  }) {
    final result = create();
    if (entries != null) result.entries.addAll(entries);
    return result;
  }

  PutPartnerEventsRequest._();

  factory PutPartnerEventsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutPartnerEventsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutPartnerEventsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<PutPartnerEventsRequestEntry>(
        481075860, _omitFieldNames ? '' : 'entries',
        subBuilder: PutPartnerEventsRequestEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPartnerEventsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPartnerEventsRequest copyWith(
          void Function(PutPartnerEventsRequest) updates) =>
      super.copyWith((message) => updates(message as PutPartnerEventsRequest))
          as PutPartnerEventsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutPartnerEventsRequest create() => PutPartnerEventsRequest._();
  @$core.override
  PutPartnerEventsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutPartnerEventsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutPartnerEventsRequest>(create);
  static PutPartnerEventsRequest? _defaultInstance;

  @$pb.TagNumber(481075860)
  $pb.PbList<PutPartnerEventsRequestEntry> get entries => $_getList(0);
}

class PutPartnerEventsRequestEntry extends $pb.GeneratedMessage {
  factory PutPartnerEventsRequestEntry({
    $core.String? detailtype,
    $core.String? source,
    $core.String? detail,
    $core.Iterable<$core.String>? resources,
    $core.String? time,
  }) {
    final result = create();
    if (detailtype != null) result.detailtype = detailtype;
    if (source != null) result.source = source;
    if (detail != null) result.detail = detail;
    if (resources != null) result.resources.addAll(resources);
    if (time != null) result.time = time;
    return result;
  }

  PutPartnerEventsRequestEntry._();

  factory PutPartnerEventsRequestEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutPartnerEventsRequestEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutPartnerEventsRequestEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(11758589, _omitFieldNames ? '' : 'detailtype')
    ..aOS(31630329, _omitFieldNames ? '' : 'source')
    ..aOS(340030389, _omitFieldNames ? '' : 'detail')
    ..pPS(358918291, _omitFieldNames ? '' : 'resources')
    ..aOS(535094277, _omitFieldNames ? '' : 'time')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPartnerEventsRequestEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPartnerEventsRequestEntry copyWith(
          void Function(PutPartnerEventsRequestEntry) updates) =>
      super.copyWith(
              (message) => updates(message as PutPartnerEventsRequestEntry))
          as PutPartnerEventsRequestEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutPartnerEventsRequestEntry create() =>
      PutPartnerEventsRequestEntry._();
  @$core.override
  PutPartnerEventsRequestEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutPartnerEventsRequestEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutPartnerEventsRequestEntry>(create);
  static PutPartnerEventsRequestEntry? _defaultInstance;

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

  @$pb.TagNumber(340030389)
  $core.String get detail => $_getSZ(2);
  @$pb.TagNumber(340030389)
  set detail($core.String value) => $_setString(2, value);
  @$pb.TagNumber(340030389)
  $core.bool hasDetail() => $_has(2);
  @$pb.TagNumber(340030389)
  void clearDetail() => $_clearField(340030389);

  @$pb.TagNumber(358918291)
  $pb.PbList<$core.String> get resources => $_getList(3);

  @$pb.TagNumber(535094277)
  $core.String get time => $_getSZ(4);
  @$pb.TagNumber(535094277)
  set time($core.String value) => $_setString(4, value);
  @$pb.TagNumber(535094277)
  $core.bool hasTime() => $_has(4);
  @$pb.TagNumber(535094277)
  void clearTime() => $_clearField(535094277);
}

class PutPartnerEventsResponse extends $pb.GeneratedMessage {
  factory PutPartnerEventsResponse({
    $core.int? failedentrycount,
    $core.Iterable<PutPartnerEventsResultEntry>? entries,
  }) {
    final result = create();
    if (failedentrycount != null) result.failedentrycount = failedentrycount;
    if (entries != null) result.entries.addAll(entries);
    return result;
  }

  PutPartnerEventsResponse._();

  factory PutPartnerEventsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutPartnerEventsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutPartnerEventsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aI(458519576, _omitFieldNames ? '' : 'failedentrycount')
    ..pPM<PutPartnerEventsResultEntry>(
        481075860, _omitFieldNames ? '' : 'entries',
        subBuilder: PutPartnerEventsResultEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPartnerEventsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPartnerEventsResponse copyWith(
          void Function(PutPartnerEventsResponse) updates) =>
      super.copyWith((message) => updates(message as PutPartnerEventsResponse))
          as PutPartnerEventsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutPartnerEventsResponse create() => PutPartnerEventsResponse._();
  @$core.override
  PutPartnerEventsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutPartnerEventsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutPartnerEventsResponse>(create);
  static PutPartnerEventsResponse? _defaultInstance;

  @$pb.TagNumber(458519576)
  $core.int get failedentrycount => $_getIZ(0);
  @$pb.TagNumber(458519576)
  set failedentrycount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(458519576)
  $core.bool hasFailedentrycount() => $_has(0);
  @$pb.TagNumber(458519576)
  void clearFailedentrycount() => $_clearField(458519576);

  @$pb.TagNumber(481075860)
  $pb.PbList<PutPartnerEventsResultEntry> get entries => $_getList(1);
}

class PutPartnerEventsResultEntry extends $pb.GeneratedMessage {
  factory PutPartnerEventsResultEntry({
    $core.String? errorcode,
    $core.String? eventid,
    $core.String? errormessage,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (eventid != null) result.eventid = eventid;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  PutPartnerEventsResultEntry._();

  factory PutPartnerEventsResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutPartnerEventsResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutPartnerEventsResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(376916819, _omitFieldNames ? '' : 'eventid')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPartnerEventsResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPartnerEventsResultEntry copyWith(
          void Function(PutPartnerEventsResultEntry) updates) =>
      super.copyWith(
              (message) => updates(message as PutPartnerEventsResultEntry))
          as PutPartnerEventsResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutPartnerEventsResultEntry create() =>
      PutPartnerEventsResultEntry._();
  @$core.override
  PutPartnerEventsResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutPartnerEventsResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutPartnerEventsResultEntry>(create);
  static PutPartnerEventsResultEntry? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(376916819)
  $core.String get eventid => $_getSZ(1);
  @$pb.TagNumber(376916819)
  set eventid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(376916819)
  $core.bool hasEventid() => $_has(1);
  @$pb.TagNumber(376916819)
  void clearEventid() => $_clearField(376916819);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(2);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(2, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(2);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class PutPermissionRequest extends $pb.GeneratedMessage {
  factory PutPermissionRequest({
    $core.String? statementid,
    $core.String? action,
    Condition? condition,
    $core.String? principal,
    $core.String? eventbusname,
    $core.String? policy,
  }) {
    final result = create();
    if (statementid != null) result.statementid = statementid;
    if (action != null) result.action = action;
    if (condition != null) result.condition = condition;
    if (principal != null) result.principal = principal;
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (policy != null) result.policy = policy;
    return result;
  }

  PutPermissionRequest._();

  factory PutPermissionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutPermissionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutPermissionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(169602348, _omitFieldNames ? '' : 'statementid')
    ..aOS(175614240, _omitFieldNames ? '' : 'action')
    ..aOM<Condition>(212017247, _omitFieldNames ? '' : 'condition',
        subBuilder: Condition.create)
    ..aOS(361640138, _omitFieldNames ? '' : 'principal')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aOS(471611296, _omitFieldNames ? '' : 'policy')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPermissionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutPermissionRequest copyWith(void Function(PutPermissionRequest) updates) =>
      super.copyWith((message) => updates(message as PutPermissionRequest))
          as PutPermissionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutPermissionRequest create() => PutPermissionRequest._();
  @$core.override
  PutPermissionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutPermissionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutPermissionRequest>(create);
  static PutPermissionRequest? _defaultInstance;

  @$pb.TagNumber(169602348)
  $core.String get statementid => $_getSZ(0);
  @$pb.TagNumber(169602348)
  set statementid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(169602348)
  $core.bool hasStatementid() => $_has(0);
  @$pb.TagNumber(169602348)
  void clearStatementid() => $_clearField(169602348);

  @$pb.TagNumber(175614240)
  $core.String get action => $_getSZ(1);
  @$pb.TagNumber(175614240)
  set action($core.String value) => $_setString(1, value);
  @$pb.TagNumber(175614240)
  $core.bool hasAction() => $_has(1);
  @$pb.TagNumber(175614240)
  void clearAction() => $_clearField(175614240);

  @$pb.TagNumber(212017247)
  Condition get condition => $_getN(2);
  @$pb.TagNumber(212017247)
  set condition(Condition value) => $_setField(212017247, value);
  @$pb.TagNumber(212017247)
  $core.bool hasCondition() => $_has(2);
  @$pb.TagNumber(212017247)
  void clearCondition() => $_clearField(212017247);
  @$pb.TagNumber(212017247)
  Condition ensureCondition() => $_ensure(2);

  @$pb.TagNumber(361640138)
  $core.String get principal => $_getSZ(3);
  @$pb.TagNumber(361640138)
  set principal($core.String value) => $_setString(3, value);
  @$pb.TagNumber(361640138)
  $core.bool hasPrincipal() => $_has(3);
  @$pb.TagNumber(361640138)
  void clearPrincipal() => $_clearField(361640138);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(4);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(4);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(471611296)
  $core.String get policy => $_getSZ(5);
  @$pb.TagNumber(471611296)
  set policy($core.String value) => $_setString(5, value);
  @$pb.TagNumber(471611296)
  $core.bool hasPolicy() => $_has(5);
  @$pb.TagNumber(471611296)
  void clearPolicy() => $_clearField(471611296);
}

class PutRuleRequest extends $pb.GeneratedMessage {
  factory PutRuleRequest({
    $core.String? description,
    $core.String? eventpattern,
    $core.String? name,
    $core.String? rolearn,
    $core.Iterable<Tag>? tags,
    $core.String? scheduleexpression,
    $core.String? eventbusname,
    RuleState? state,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (eventpattern != null) result.eventpattern = eventpattern;
    if (name != null) result.name = name;
    if (rolearn != null) result.rolearn = rolearn;
    if (tags != null) result.tags.addAll(tags);
    if (scheduleexpression != null)
      result.scheduleexpression = scheduleexpression;
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (state != null) result.state = state;
    return result;
  }

  PutRuleRequest._();

  factory PutRuleRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRuleRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRuleRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(233487232, _omitFieldNames ? '' : 'eventpattern')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(446089471, _omitFieldNames ? '' : 'scheduleexpression')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aE<RuleState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: RuleState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRuleRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRuleRequest copyWith(void Function(PutRuleRequest) updates) =>
      super.copyWith((message) => updates(message as PutRuleRequest))
          as PutRuleRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRuleRequest create() => PutRuleRequest._();
  @$core.override
  PutRuleRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRuleRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRuleRequest>(create);
  static PutRuleRequest? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(233487232)
  $core.String get eventpattern => $_getSZ(1);
  @$pb.TagNumber(233487232)
  set eventpattern($core.String value) => $_setString(1, value);
  @$pb.TagNumber(233487232)
  $core.bool hasEventpattern() => $_has(1);
  @$pb.TagNumber(233487232)
  void clearEventpattern() => $_clearField(233487232);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(3);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(3);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(4);

  @$pb.TagNumber(446089471)
  $core.String get scheduleexpression => $_getSZ(5);
  @$pb.TagNumber(446089471)
  set scheduleexpression($core.String value) => $_setString(5, value);
  @$pb.TagNumber(446089471)
  $core.bool hasScheduleexpression() => $_has(5);
  @$pb.TagNumber(446089471)
  void clearScheduleexpression() => $_clearField(446089471);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(6);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(6);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(502047895)
  RuleState get state => $_getN(7);
  @$pb.TagNumber(502047895)
  set state(RuleState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(7);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class PutRuleResponse extends $pb.GeneratedMessage {
  factory PutRuleResponse({
    $core.String? rulearn,
  }) {
    final result = create();
    if (rulearn != null) result.rulearn = rulearn;
    return result;
  }

  PutRuleResponse._();

  factory PutRuleResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutRuleResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutRuleResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(217507655, _omitFieldNames ? '' : 'rulearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRuleResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutRuleResponse copyWith(void Function(PutRuleResponse) updates) =>
      super.copyWith((message) => updates(message as PutRuleResponse))
          as PutRuleResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutRuleResponse create() => PutRuleResponse._();
  @$core.override
  PutRuleResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutRuleResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutRuleResponse>(create);
  static PutRuleResponse? _defaultInstance;

  @$pb.TagNumber(217507655)
  $core.String get rulearn => $_getSZ(0);
  @$pb.TagNumber(217507655)
  set rulearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(217507655)
  $core.bool hasRulearn() => $_has(0);
  @$pb.TagNumber(217507655)
  void clearRulearn() => $_clearField(217507655);
}

class PutTargetsRequest extends $pb.GeneratedMessage {
  factory PutTargetsRequest({
    $core.Iterable<Target>? targets,
    $core.String? eventbusname,
    $core.String? rule,
  }) {
    final result = create();
    if (targets != null) result.targets.addAll(targets);
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (rule != null) result.rule = rule;
    return result;
  }

  PutTargetsRequest._();

  factory PutTargetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutTargetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutTargetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<Target>(262180226, _omitFieldNames ? '' : 'targets',
        subBuilder: Target.create)
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aOS(475696372, _omitFieldNames ? '' : 'rule')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutTargetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutTargetsRequest copyWith(void Function(PutTargetsRequest) updates) =>
      super.copyWith((message) => updates(message as PutTargetsRequest))
          as PutTargetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutTargetsRequest create() => PutTargetsRequest._();
  @$core.override
  PutTargetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutTargetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutTargetsRequest>(create);
  static PutTargetsRequest? _defaultInstance;

  @$pb.TagNumber(262180226)
  $pb.PbList<Target> get targets => $_getList(0);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(1);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(1);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(475696372)
  $core.String get rule => $_getSZ(2);
  @$pb.TagNumber(475696372)
  set rule($core.String value) => $_setString(2, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(2);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);
}

class PutTargetsResponse extends $pb.GeneratedMessage {
  factory PutTargetsResponse({
    $core.Iterable<PutTargetsResultEntry>? failedentries,
    $core.int? failedentrycount,
  }) {
    final result = create();
    if (failedentries != null) result.failedentries.addAll(failedentries);
    if (failedentrycount != null) result.failedentrycount = failedentrycount;
    return result;
  }

  PutTargetsResponse._();

  factory PutTargetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutTargetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutTargetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<PutTargetsResultEntry>(
        86416685, _omitFieldNames ? '' : 'failedentries',
        subBuilder: PutTargetsResultEntry.create)
    ..aI(458519576, _omitFieldNames ? '' : 'failedentrycount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutTargetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutTargetsResponse copyWith(void Function(PutTargetsResponse) updates) =>
      super.copyWith((message) => updates(message as PutTargetsResponse))
          as PutTargetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutTargetsResponse create() => PutTargetsResponse._();
  @$core.override
  PutTargetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutTargetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutTargetsResponse>(create);
  static PutTargetsResponse? _defaultInstance;

  @$pb.TagNumber(86416685)
  $pb.PbList<PutTargetsResultEntry> get failedentries => $_getList(0);

  @$pb.TagNumber(458519576)
  $core.int get failedentrycount => $_getIZ(1);
  @$pb.TagNumber(458519576)
  set failedentrycount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(458519576)
  $core.bool hasFailedentrycount() => $_has(1);
  @$pb.TagNumber(458519576)
  void clearFailedentrycount() => $_clearField(458519576);
}

class PutTargetsResultEntry extends $pb.GeneratedMessage {
  factory PutTargetsResultEntry({
    $core.String? errorcode,
    $core.String? targetid,
    $core.String? errormessage,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (targetid != null) result.targetid = targetid;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  PutTargetsResultEntry._();

  factory PutTargetsResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutTargetsResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutTargetsResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(46993334, _omitFieldNames ? '' : 'targetid')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutTargetsResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutTargetsResultEntry copyWith(
          void Function(PutTargetsResultEntry) updates) =>
      super.copyWith((message) => updates(message as PutTargetsResultEntry))
          as PutTargetsResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutTargetsResultEntry create() => PutTargetsResultEntry._();
  @$core.override
  PutTargetsResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutTargetsResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutTargetsResultEntry>(create);
  static PutTargetsResultEntry? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(46993334)
  $core.String get targetid => $_getSZ(1);
  @$pb.TagNumber(46993334)
  set targetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(46993334)
  $core.bool hasTargetid() => $_has(1);
  @$pb.TagNumber(46993334)
  void clearTargetid() => $_clearField(46993334);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(2);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(2, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(2);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class RedshiftDataParameters extends $pb.GeneratedMessage {
  factory RedshiftDataParameters({
    $core.bool? withevent,
    $core.String? sql,
    $core.String? statementname,
    $core.String? database,
    $core.String? dbuser,
    $core.String? secretmanagerarn,
  }) {
    final result = create();
    if (withevent != null) result.withevent = withevent;
    if (sql != null) result.sql = sql;
    if (statementname != null) result.statementname = statementname;
    if (database != null) result.database = database;
    if (dbuser != null) result.dbuser = dbuser;
    if (secretmanagerarn != null) result.secretmanagerarn = secretmanagerarn;
    return result;
  }

  RedshiftDataParameters._();

  factory RedshiftDataParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RedshiftDataParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RedshiftDataParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOB(5518724, _omitFieldNames ? '' : 'withevent')
    ..aOS(13608736, _omitFieldNames ? '' : 'sql')
    ..aOS(23047926, _omitFieldNames ? '' : 'statementname')
    ..aOS(278147289, _omitFieldNames ? '' : 'database')
    ..aOS(392658689, _omitFieldNames ? '' : 'dbuser')
    ..aOS(530394088, _omitFieldNames ? '' : 'secretmanagerarn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedshiftDataParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RedshiftDataParameters copyWith(
          void Function(RedshiftDataParameters) updates) =>
      super.copyWith((message) => updates(message as RedshiftDataParameters))
          as RedshiftDataParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RedshiftDataParameters create() => RedshiftDataParameters._();
  @$core.override
  RedshiftDataParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RedshiftDataParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RedshiftDataParameters>(create);
  static RedshiftDataParameters? _defaultInstance;

  @$pb.TagNumber(5518724)
  $core.bool get withevent => $_getBF(0);
  @$pb.TagNumber(5518724)
  set withevent($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(5518724)
  $core.bool hasWithevent() => $_has(0);
  @$pb.TagNumber(5518724)
  void clearWithevent() => $_clearField(5518724);

  @$pb.TagNumber(13608736)
  $core.String get sql => $_getSZ(1);
  @$pb.TagNumber(13608736)
  set sql($core.String value) => $_setString(1, value);
  @$pb.TagNumber(13608736)
  $core.bool hasSql() => $_has(1);
  @$pb.TagNumber(13608736)
  void clearSql() => $_clearField(13608736);

  @$pb.TagNumber(23047926)
  $core.String get statementname => $_getSZ(2);
  @$pb.TagNumber(23047926)
  set statementname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(23047926)
  $core.bool hasStatementname() => $_has(2);
  @$pb.TagNumber(23047926)
  void clearStatementname() => $_clearField(23047926);

  @$pb.TagNumber(278147289)
  $core.String get database => $_getSZ(3);
  @$pb.TagNumber(278147289)
  set database($core.String value) => $_setString(3, value);
  @$pb.TagNumber(278147289)
  $core.bool hasDatabase() => $_has(3);
  @$pb.TagNumber(278147289)
  void clearDatabase() => $_clearField(278147289);

  @$pb.TagNumber(392658689)
  $core.String get dbuser => $_getSZ(4);
  @$pb.TagNumber(392658689)
  set dbuser($core.String value) => $_setString(4, value);
  @$pb.TagNumber(392658689)
  $core.bool hasDbuser() => $_has(4);
  @$pb.TagNumber(392658689)
  void clearDbuser() => $_clearField(392658689);

  @$pb.TagNumber(530394088)
  $core.String get secretmanagerarn => $_getSZ(5);
  @$pb.TagNumber(530394088)
  set secretmanagerarn($core.String value) => $_setString(5, value);
  @$pb.TagNumber(530394088)
  $core.bool hasSecretmanagerarn() => $_has(5);
  @$pb.TagNumber(530394088)
  void clearSecretmanagerarn() => $_clearField(530394088);
}

class RemovePermissionRequest extends $pb.GeneratedMessage {
  factory RemovePermissionRequest({
    $core.bool? removeallpermissions,
    $core.String? statementid,
    $core.String? eventbusname,
  }) {
    final result = create();
    if (removeallpermissions != null)
      result.removeallpermissions = removeallpermissions;
    if (statementid != null) result.statementid = statementid;
    if (eventbusname != null) result.eventbusname = eventbusname;
    return result;
  }

  RemovePermissionRequest._();

  factory RemovePermissionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemovePermissionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemovePermissionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOB(8309441, _omitFieldNames ? '' : 'removeallpermissions')
    ..aOS(169602348, _omitFieldNames ? '' : 'statementid')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemovePermissionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemovePermissionRequest copyWith(
          void Function(RemovePermissionRequest) updates) =>
      super.copyWith((message) => updates(message as RemovePermissionRequest))
          as RemovePermissionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemovePermissionRequest create() => RemovePermissionRequest._();
  @$core.override
  RemovePermissionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemovePermissionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemovePermissionRequest>(create);
  static RemovePermissionRequest? _defaultInstance;

  @$pb.TagNumber(8309441)
  $core.bool get removeallpermissions => $_getBF(0);
  @$pb.TagNumber(8309441)
  set removeallpermissions($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(8309441)
  $core.bool hasRemoveallpermissions() => $_has(0);
  @$pb.TagNumber(8309441)
  void clearRemoveallpermissions() => $_clearField(8309441);

  @$pb.TagNumber(169602348)
  $core.String get statementid => $_getSZ(1);
  @$pb.TagNumber(169602348)
  set statementid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(169602348)
  $core.bool hasStatementid() => $_has(1);
  @$pb.TagNumber(169602348)
  void clearStatementid() => $_clearField(169602348);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(2);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(2);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);
}

class RemoveTargetsRequest extends $pb.GeneratedMessage {
  factory RemoveTargetsRequest({
    $core.Iterable<$core.String>? ids,
    $core.String? eventbusname,
    $core.String? rule,
    $core.bool? force,
  }) {
    final result = create();
    if (ids != null) result.ids.addAll(ids);
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (rule != null) result.rule = rule;
    if (force != null) result.force = force;
    return result;
  }

  RemoveTargetsRequest._();

  factory RemoveTargetsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveTargetsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveTargetsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPS(56356874, _omitFieldNames ? '' : 'ids')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aOS(475696372, _omitFieldNames ? '' : 'rule')
    ..aOB(526711333, _omitFieldNames ? '' : 'force')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTargetsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTargetsRequest copyWith(void Function(RemoveTargetsRequest) updates) =>
      super.copyWith((message) => updates(message as RemoveTargetsRequest))
          as RemoveTargetsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveTargetsRequest create() => RemoveTargetsRequest._();
  @$core.override
  RemoveTargetsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveTargetsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveTargetsRequest>(create);
  static RemoveTargetsRequest? _defaultInstance;

  @$pb.TagNumber(56356874)
  $pb.PbList<$core.String> get ids => $_getList(0);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(1);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(1);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(475696372)
  $core.String get rule => $_getSZ(2);
  @$pb.TagNumber(475696372)
  set rule($core.String value) => $_setString(2, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(2);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);

  @$pb.TagNumber(526711333)
  $core.bool get force => $_getBF(3);
  @$pb.TagNumber(526711333)
  set force($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(526711333)
  $core.bool hasForce() => $_has(3);
  @$pb.TagNumber(526711333)
  void clearForce() => $_clearField(526711333);
}

class RemoveTargetsResponse extends $pb.GeneratedMessage {
  factory RemoveTargetsResponse({
    $core.Iterable<RemoveTargetsResultEntry>? failedentries,
    $core.int? failedentrycount,
  }) {
    final result = create();
    if (failedentries != null) result.failedentries.addAll(failedentries);
    if (failedentrycount != null) result.failedentrycount = failedentrycount;
    return result;
  }

  RemoveTargetsResponse._();

  factory RemoveTargetsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveTargetsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveTargetsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<RemoveTargetsResultEntry>(
        86416685, _omitFieldNames ? '' : 'failedentries',
        subBuilder: RemoveTargetsResultEntry.create)
    ..aI(458519576, _omitFieldNames ? '' : 'failedentrycount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTargetsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTargetsResponse copyWith(
          void Function(RemoveTargetsResponse) updates) =>
      super.copyWith((message) => updates(message as RemoveTargetsResponse))
          as RemoveTargetsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveTargetsResponse create() => RemoveTargetsResponse._();
  @$core.override
  RemoveTargetsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveTargetsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveTargetsResponse>(create);
  static RemoveTargetsResponse? _defaultInstance;

  @$pb.TagNumber(86416685)
  $pb.PbList<RemoveTargetsResultEntry> get failedentries => $_getList(0);

  @$pb.TagNumber(458519576)
  $core.int get failedentrycount => $_getIZ(1);
  @$pb.TagNumber(458519576)
  set failedentrycount($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(458519576)
  $core.bool hasFailedentrycount() => $_has(1);
  @$pb.TagNumber(458519576)
  void clearFailedentrycount() => $_clearField(458519576);
}

class RemoveTargetsResultEntry extends $pb.GeneratedMessage {
  factory RemoveTargetsResultEntry({
    $core.String? errorcode,
    $core.String? targetid,
    $core.String? errormessage,
  }) {
    final result = create();
    if (errorcode != null) result.errorcode = errorcode;
    if (targetid != null) result.targetid = targetid;
    if (errormessage != null) result.errormessage = errormessage;
    return result;
  }

  RemoveTargetsResultEntry._();

  factory RemoveTargetsResultEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RemoveTargetsResultEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RemoveTargetsResultEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(34663193, _omitFieldNames ? '' : 'errorcode')
    ..aOS(46993334, _omitFieldNames ? '' : 'targetid')
    ..aOS(518702377, _omitFieldNames ? '' : 'errormessage')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTargetsResultEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RemoveTargetsResultEntry copyWith(
          void Function(RemoveTargetsResultEntry) updates) =>
      super.copyWith((message) => updates(message as RemoveTargetsResultEntry))
          as RemoveTargetsResultEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RemoveTargetsResultEntry create() => RemoveTargetsResultEntry._();
  @$core.override
  RemoveTargetsResultEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RemoveTargetsResultEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RemoveTargetsResultEntry>(create);
  static RemoveTargetsResultEntry? _defaultInstance;

  @$pb.TagNumber(34663193)
  $core.String get errorcode => $_getSZ(0);
  @$pb.TagNumber(34663193)
  set errorcode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(34663193)
  $core.bool hasErrorcode() => $_has(0);
  @$pb.TagNumber(34663193)
  void clearErrorcode() => $_clearField(34663193);

  @$pb.TagNumber(46993334)
  $core.String get targetid => $_getSZ(1);
  @$pb.TagNumber(46993334)
  set targetid($core.String value) => $_setString(1, value);
  @$pb.TagNumber(46993334)
  $core.bool hasTargetid() => $_has(1);
  @$pb.TagNumber(46993334)
  void clearTargetid() => $_clearField(46993334);

  @$pb.TagNumber(518702377)
  $core.String get errormessage => $_getSZ(2);
  @$pb.TagNumber(518702377)
  set errormessage($core.String value) => $_setString(2, value);
  @$pb.TagNumber(518702377)
  $core.bool hasErrormessage() => $_has(2);
  @$pb.TagNumber(518702377)
  void clearErrormessage() => $_clearField(518702377);
}

class Replay extends $pb.GeneratedMessage {
  factory Replay({
    $core.String? eventlastreplayedtime,
    $core.String? eventendtime,
    $core.String? eventstarttime,
    $core.String? replayendtime,
    $core.String? eventsourcearn,
    $core.String? statereason,
    $core.String? replayname,
    ReplayState? state,
    $core.String? replaystarttime,
  }) {
    final result = create();
    if (eventlastreplayedtime != null)
      result.eventlastreplayedtime = eventlastreplayedtime;
    if (eventendtime != null) result.eventendtime = eventendtime;
    if (eventstarttime != null) result.eventstarttime = eventstarttime;
    if (replayendtime != null) result.replayendtime = replayendtime;
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    if (statereason != null) result.statereason = statereason;
    if (replayname != null) result.replayname = replayname;
    if (state != null) result.state = state;
    if (replaystarttime != null) result.replaystarttime = replaystarttime;
    return result;
  }

  Replay._();

  factory Replay.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Replay.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Replay',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(7759073, _omitFieldNames ? '' : 'eventlastreplayedtime')
    ..aOS(30229582, _omitFieldNames ? '' : 'eventendtime')
    ..aOS(103680757, _omitFieldNames ? '' : 'eventstarttime')
    ..aOS(304261199, _omitFieldNames ? '' : 'replayendtime')
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aOS(442173850, _omitFieldNames ? '' : 'replayname')
    ..aE<ReplayState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ReplayState.values)
    ..aOS(503580492, _omitFieldNames ? '' : 'replaystarttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Replay clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Replay copyWith(void Function(Replay) updates) =>
      super.copyWith((message) => updates(message as Replay)) as Replay;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Replay create() => Replay._();
  @$core.override
  Replay createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Replay getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Replay>(create);
  static Replay? _defaultInstance;

  @$pb.TagNumber(7759073)
  $core.String get eventlastreplayedtime => $_getSZ(0);
  @$pb.TagNumber(7759073)
  set eventlastreplayedtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(7759073)
  $core.bool hasEventlastreplayedtime() => $_has(0);
  @$pb.TagNumber(7759073)
  void clearEventlastreplayedtime() => $_clearField(7759073);

  @$pb.TagNumber(30229582)
  $core.String get eventendtime => $_getSZ(1);
  @$pb.TagNumber(30229582)
  set eventendtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(30229582)
  $core.bool hasEventendtime() => $_has(1);
  @$pb.TagNumber(30229582)
  void clearEventendtime() => $_clearField(30229582);

  @$pb.TagNumber(103680757)
  $core.String get eventstarttime => $_getSZ(2);
  @$pb.TagNumber(103680757)
  set eventstarttime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103680757)
  $core.bool hasEventstarttime() => $_has(2);
  @$pb.TagNumber(103680757)
  void clearEventstarttime() => $_clearField(103680757);

  @$pb.TagNumber(304261199)
  $core.String get replayendtime => $_getSZ(3);
  @$pb.TagNumber(304261199)
  set replayendtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(304261199)
  $core.bool hasReplayendtime() => $_has(3);
  @$pb.TagNumber(304261199)
  void clearReplayendtime() => $_clearField(304261199);

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(4);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(4);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(5);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(5, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(5);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(442173850)
  $core.String get replayname => $_getSZ(6);
  @$pb.TagNumber(442173850)
  set replayname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(442173850)
  $core.bool hasReplayname() => $_has(6);
  @$pb.TagNumber(442173850)
  void clearReplayname() => $_clearField(442173850);

  @$pb.TagNumber(502047895)
  ReplayState get state => $_getN(7);
  @$pb.TagNumber(502047895)
  set state(ReplayState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(7);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);

  @$pb.TagNumber(503580492)
  $core.String get replaystarttime => $_getSZ(8);
  @$pb.TagNumber(503580492)
  set replaystarttime($core.String value) => $_setString(8, value);
  @$pb.TagNumber(503580492)
  $core.bool hasReplaystarttime() => $_has(8);
  @$pb.TagNumber(503580492)
  void clearReplaystarttime() => $_clearField(503580492);
}

class ReplayDestination extends $pb.GeneratedMessage {
  factory ReplayDestination({
    $core.Iterable<$core.String>? filterarns,
    $core.String? arn,
  }) {
    final result = create();
    if (filterarns != null) result.filterarns.addAll(filterarns);
    if (arn != null) result.arn = arn;
    return result;
  }

  ReplayDestination._();

  factory ReplayDestination.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ReplayDestination.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ReplayDestination',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPS(135157548, _omitFieldNames ? '' : 'filterarns')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplayDestination clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ReplayDestination copyWith(void Function(ReplayDestination) updates) =>
      super.copyWith((message) => updates(message as ReplayDestination))
          as ReplayDestination;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ReplayDestination create() => ReplayDestination._();
  @$core.override
  ReplayDestination createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ReplayDestination getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ReplayDestination>(create);
  static ReplayDestination? _defaultInstance;

  @$pb.TagNumber(135157548)
  $pb.PbList<$core.String> get filterarns => $_getList(0);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(1);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(1);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class ResourceAlreadyExistsException extends $pb.GeneratedMessage {
  factory ResourceAlreadyExistsException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourceAlreadyExistsException._();

  factory ResourceAlreadyExistsException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceAlreadyExistsException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceAlreadyExistsException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceAlreadyExistsException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceAlreadyExistsException copyWith(
          void Function(ResourceAlreadyExistsException) updates) =>
      super.copyWith(
              (message) => updates(message as ResourceAlreadyExistsException))
          as ResourceAlreadyExistsException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceAlreadyExistsException create() =>
      ResourceAlreadyExistsException._();
  @$core.override
  ResourceAlreadyExistsException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceAlreadyExistsException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceAlreadyExistsException>(create);
  static ResourceAlreadyExistsException? _defaultInstance;

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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class Rule extends $pb.GeneratedMessage {
  factory Rule({
    $core.String? description,
    $core.String? eventpattern,
    $core.String? name,
    $core.String? rolearn,
    $core.String? arn,
    $core.String? scheduleexpression,
    $core.String? eventbusname,
    $core.String? managedby,
    RuleState? state,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (eventpattern != null) result.eventpattern = eventpattern;
    if (name != null) result.name = name;
    if (rolearn != null) result.rolearn = rolearn;
    if (arn != null) result.arn = arn;
    if (scheduleexpression != null)
      result.scheduleexpression = scheduleexpression;
    if (eventbusname != null) result.eventbusname = eventbusname;
    if (managedby != null) result.managedby = managedby;
    if (state != null) result.state = state;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(233487232, _omitFieldNames ? '' : 'eventpattern')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(446089471, _omitFieldNames ? '' : 'scheduleexpression')
    ..aOS(448449725, _omitFieldNames ? '' : 'eventbusname')
    ..aOS(455511232, _omitFieldNames ? '' : 'managedby')
    ..aE<RuleState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: RuleState.values)
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

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(233487232)
  $core.String get eventpattern => $_getSZ(1);
  @$pb.TagNumber(233487232)
  set eventpattern($core.String value) => $_setString(1, value);
  @$pb.TagNumber(233487232)
  $core.bool hasEventpattern() => $_has(1);
  @$pb.TagNumber(233487232)
  void clearEventpattern() => $_clearField(233487232);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(3);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(3);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(446089471)
  $core.String get scheduleexpression => $_getSZ(5);
  @$pb.TagNumber(446089471)
  set scheduleexpression($core.String value) => $_setString(5, value);
  @$pb.TagNumber(446089471)
  $core.bool hasScheduleexpression() => $_has(5);
  @$pb.TagNumber(446089471)
  void clearScheduleexpression() => $_clearField(446089471);

  @$pb.TagNumber(448449725)
  $core.String get eventbusname => $_getSZ(6);
  @$pb.TagNumber(448449725)
  set eventbusname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(448449725)
  $core.bool hasEventbusname() => $_has(6);
  @$pb.TagNumber(448449725)
  void clearEventbusname() => $_clearField(448449725);

  @$pb.TagNumber(455511232)
  $core.String get managedby => $_getSZ(7);
  @$pb.TagNumber(455511232)
  set managedby($core.String value) => $_setString(7, value);
  @$pb.TagNumber(455511232)
  $core.bool hasManagedby() => $_has(7);
  @$pb.TagNumber(455511232)
  void clearManagedby() => $_clearField(455511232);

  @$pb.TagNumber(502047895)
  RuleState get state => $_getN(8);
  @$pb.TagNumber(502047895)
  set state(RuleState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(8);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class RunCommandParameters extends $pb.GeneratedMessage {
  factory RunCommandParameters({
    $core.Iterable<RunCommandTarget>? runcommandtargets,
  }) {
    final result = create();
    if (runcommandtargets != null)
      result.runcommandtargets.addAll(runcommandtargets);
    return result;
  }

  RunCommandParameters._();

  factory RunCommandParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RunCommandParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RunCommandParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..pPM<RunCommandTarget>(
        422375352, _omitFieldNames ? '' : 'runcommandtargets',
        subBuilder: RunCommandTarget.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RunCommandParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RunCommandParameters copyWith(void Function(RunCommandParameters) updates) =>
      super.copyWith((message) => updates(message as RunCommandParameters))
          as RunCommandParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RunCommandParameters create() => RunCommandParameters._();
  @$core.override
  RunCommandParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RunCommandParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RunCommandParameters>(create);
  static RunCommandParameters? _defaultInstance;

  @$pb.TagNumber(422375352)
  $pb.PbList<RunCommandTarget> get runcommandtargets => $_getList(0);
}

class RunCommandTarget extends $pb.GeneratedMessage {
  factory RunCommandTarget({
    $core.String? key,
    $core.Iterable<$core.String>? values,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (values != null) result.values.addAll(values);
    return result;
  }

  RunCommandTarget._();

  factory RunCommandTarget.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RunCommandTarget.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RunCommandTarget',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..pPS(223158876, _omitFieldNames ? '' : 'values')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RunCommandTarget clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RunCommandTarget copyWith(void Function(RunCommandTarget) updates) =>
      super.copyWith((message) => updates(message as RunCommandTarget))
          as RunCommandTarget;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RunCommandTarget create() => RunCommandTarget._();
  @$core.override
  RunCommandTarget createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RunCommandTarget getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RunCommandTarget>(create);
  static RunCommandTarget? _defaultInstance;

  @$pb.TagNumber(219859213)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(219859213)
  set key($core.String value) => $_setString(0, value);
  @$pb.TagNumber(219859213)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(219859213)
  void clearKey() => $_clearField(219859213);

  @$pb.TagNumber(223158876)
  $pb.PbList<$core.String> get values => $_getList(1);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class StartReplayRequest extends $pb.GeneratedMessage {
  factory StartReplayRequest({
    $core.String? eventendtime,
    $core.String? eventstarttime,
    $core.String? description,
    $core.String? eventsourcearn,
    $core.String? replayname,
    ReplayDestination? destination,
  }) {
    final result = create();
    if (eventendtime != null) result.eventendtime = eventendtime;
    if (eventstarttime != null) result.eventstarttime = eventstarttime;
    if (description != null) result.description = description;
    if (eventsourcearn != null) result.eventsourcearn = eventsourcearn;
    if (replayname != null) result.replayname = replayname;
    if (destination != null) result.destination = destination;
    return result;
  }

  StartReplayRequest._();

  factory StartReplayRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartReplayRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartReplayRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(30229582, _omitFieldNames ? '' : 'eventendtime')
    ..aOS(103680757, _omitFieldNames ? '' : 'eventstarttime')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(306357574, _omitFieldNames ? '' : 'eventsourcearn')
    ..aOS(442173850, _omitFieldNames ? '' : 'replayname')
    ..aOM<ReplayDestination>(457443680, _omitFieldNames ? '' : 'destination',
        subBuilder: ReplayDestination.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartReplayRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartReplayRequest copyWith(void Function(StartReplayRequest) updates) =>
      super.copyWith((message) => updates(message as StartReplayRequest))
          as StartReplayRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartReplayRequest create() => StartReplayRequest._();
  @$core.override
  StartReplayRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartReplayRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartReplayRequest>(create);
  static StartReplayRequest? _defaultInstance;

  @$pb.TagNumber(30229582)
  $core.String get eventendtime => $_getSZ(0);
  @$pb.TagNumber(30229582)
  set eventendtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30229582)
  $core.bool hasEventendtime() => $_has(0);
  @$pb.TagNumber(30229582)
  void clearEventendtime() => $_clearField(30229582);

  @$pb.TagNumber(103680757)
  $core.String get eventstarttime => $_getSZ(1);
  @$pb.TagNumber(103680757)
  set eventstarttime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103680757)
  $core.bool hasEventstarttime() => $_has(1);
  @$pb.TagNumber(103680757)
  void clearEventstarttime() => $_clearField(103680757);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(306357574)
  $core.String get eventsourcearn => $_getSZ(3);
  @$pb.TagNumber(306357574)
  set eventsourcearn($core.String value) => $_setString(3, value);
  @$pb.TagNumber(306357574)
  $core.bool hasEventsourcearn() => $_has(3);
  @$pb.TagNumber(306357574)
  void clearEventsourcearn() => $_clearField(306357574);

  @$pb.TagNumber(442173850)
  $core.String get replayname => $_getSZ(4);
  @$pb.TagNumber(442173850)
  set replayname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(442173850)
  $core.bool hasReplayname() => $_has(4);
  @$pb.TagNumber(442173850)
  void clearReplayname() => $_clearField(442173850);

  @$pb.TagNumber(457443680)
  ReplayDestination get destination => $_getN(5);
  @$pb.TagNumber(457443680)
  set destination(ReplayDestination value) => $_setField(457443680, value);
  @$pb.TagNumber(457443680)
  $core.bool hasDestination() => $_has(5);
  @$pb.TagNumber(457443680)
  void clearDestination() => $_clearField(457443680);
  @$pb.TagNumber(457443680)
  ReplayDestination ensureDestination() => $_ensure(5);
}

class StartReplayResponse extends $pb.GeneratedMessage {
  factory StartReplayResponse({
    $core.String? replayarn,
    $core.String? statereason,
    ReplayState? state,
    $core.String? replaystarttime,
  }) {
    final result = create();
    if (replayarn != null) result.replayarn = replayarn;
    if (statereason != null) result.statereason = statereason;
    if (state != null) result.state = state;
    if (replaystarttime != null) result.replaystarttime = replaystarttime;
    return result;
  }

  StartReplayResponse._();

  factory StartReplayResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartReplayResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartReplayResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(361946526, _omitFieldNames ? '' : 'replayarn')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aE<ReplayState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ReplayState.values)
    ..aOS(503580492, _omitFieldNames ? '' : 'replaystarttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartReplayResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartReplayResponse copyWith(void Function(StartReplayResponse) updates) =>
      super.copyWith((message) => updates(message as StartReplayResponse))
          as StartReplayResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartReplayResponse create() => StartReplayResponse._();
  @$core.override
  StartReplayResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartReplayResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartReplayResponse>(create);
  static StartReplayResponse? _defaultInstance;

  @$pb.TagNumber(361946526)
  $core.String get replayarn => $_getSZ(0);
  @$pb.TagNumber(361946526)
  set replayarn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(361946526)
  $core.bool hasReplayarn() => $_has(0);
  @$pb.TagNumber(361946526)
  void clearReplayarn() => $_clearField(361946526);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(1);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(1, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(1);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(502047895)
  ReplayState get state => $_getN(2);
  @$pb.TagNumber(502047895)
  set state(ReplayState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(2);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);

  @$pb.TagNumber(503580492)
  $core.String get replaystarttime => $_getSZ(3);
  @$pb.TagNumber(503580492)
  set replaystarttime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(503580492)
  $core.bool hasReplaystarttime() => $_has(3);
  @$pb.TagNumber(503580492)
  void clearReplaystarttime() => $_clearField(503580492);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class Target extends $pb.GeneratedMessage {
  factory Target({
    InputTransformer? inputtransformer,
    KinesisParameters? kinesisparameters,
    DeadLetterConfig? deadletterconfig,
    SqsParameters? sqsparameters,
    $core.String? inputpath,
    BatchParameters? batchparameters,
    RedshiftDataParameters? redshiftdataparameters,
    RetryPolicy? retrypolicy,
    $core.String? rolearn,
    RunCommandParameters? runcommandparameters,
    HttpParameters? httpparameters,
    $core.String? id,
    $core.String? arn,
    EcsParameters? ecsparameters,
    SageMakerPipelineParameters? sagemakerpipelineparameters,
    $core.String? input,
  }) {
    final result = create();
    if (inputtransformer != null) result.inputtransformer = inputtransformer;
    if (kinesisparameters != null) result.kinesisparameters = kinesisparameters;
    if (deadletterconfig != null) result.deadletterconfig = deadletterconfig;
    if (sqsparameters != null) result.sqsparameters = sqsparameters;
    if (inputpath != null) result.inputpath = inputpath;
    if (batchparameters != null) result.batchparameters = batchparameters;
    if (redshiftdataparameters != null)
      result.redshiftdataparameters = redshiftdataparameters;
    if (retrypolicy != null) result.retrypolicy = retrypolicy;
    if (rolearn != null) result.rolearn = rolearn;
    if (runcommandparameters != null)
      result.runcommandparameters = runcommandparameters;
    if (httpparameters != null) result.httpparameters = httpparameters;
    if (id != null) result.id = id;
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOM<InputTransformer>(31280095, _omitFieldNames ? '' : 'inputtransformer',
        subBuilder: InputTransformer.create)
    ..aOM<KinesisParameters>(
        70111902, _omitFieldNames ? '' : 'kinesisparameters',
        subBuilder: KinesisParameters.create)
    ..aOM<DeadLetterConfig>(79786642, _omitFieldNames ? '' : 'deadletterconfig',
        subBuilder: DeadLetterConfig.create)
    ..aOM<SqsParameters>(91345999, _omitFieldNames ? '' : 'sqsparameters',
        subBuilder: SqsParameters.create)
    ..aOS(223809421, _omitFieldNames ? '' : 'inputpath')
    ..aOM<BatchParameters>(242480292, _omitFieldNames ? '' : 'batchparameters',
        subBuilder: BatchParameters.create)
    ..aOM<RedshiftDataParameters>(
        254855317, _omitFieldNames ? '' : 'redshiftdataparameters',
        subBuilder: RedshiftDataParameters.create)
    ..aOM<RetryPolicy>(266827188, _omitFieldNames ? '' : 'retrypolicy',
        subBuilder: RetryPolicy.create)
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..aOM<RunCommandParameters>(
        351941188, _omitFieldNames ? '' : 'runcommandparameters',
        subBuilder: RunCommandParameters.create)
    ..aOM<HttpParameters>(367051514, _omitFieldNames ? '' : 'httpparameters',
        subBuilder: HttpParameters.create)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
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

  @$pb.TagNumber(31280095)
  InputTransformer get inputtransformer => $_getN(0);
  @$pb.TagNumber(31280095)
  set inputtransformer(InputTransformer value) => $_setField(31280095, value);
  @$pb.TagNumber(31280095)
  $core.bool hasInputtransformer() => $_has(0);
  @$pb.TagNumber(31280095)
  void clearInputtransformer() => $_clearField(31280095);
  @$pb.TagNumber(31280095)
  InputTransformer ensureInputtransformer() => $_ensure(0);

  @$pb.TagNumber(70111902)
  KinesisParameters get kinesisparameters => $_getN(1);
  @$pb.TagNumber(70111902)
  set kinesisparameters(KinesisParameters value) => $_setField(70111902, value);
  @$pb.TagNumber(70111902)
  $core.bool hasKinesisparameters() => $_has(1);
  @$pb.TagNumber(70111902)
  void clearKinesisparameters() => $_clearField(70111902);
  @$pb.TagNumber(70111902)
  KinesisParameters ensureKinesisparameters() => $_ensure(1);

  @$pb.TagNumber(79786642)
  DeadLetterConfig get deadletterconfig => $_getN(2);
  @$pb.TagNumber(79786642)
  set deadletterconfig(DeadLetterConfig value) => $_setField(79786642, value);
  @$pb.TagNumber(79786642)
  $core.bool hasDeadletterconfig() => $_has(2);
  @$pb.TagNumber(79786642)
  void clearDeadletterconfig() => $_clearField(79786642);
  @$pb.TagNumber(79786642)
  DeadLetterConfig ensureDeadletterconfig() => $_ensure(2);

  @$pb.TagNumber(91345999)
  SqsParameters get sqsparameters => $_getN(3);
  @$pb.TagNumber(91345999)
  set sqsparameters(SqsParameters value) => $_setField(91345999, value);
  @$pb.TagNumber(91345999)
  $core.bool hasSqsparameters() => $_has(3);
  @$pb.TagNumber(91345999)
  void clearSqsparameters() => $_clearField(91345999);
  @$pb.TagNumber(91345999)
  SqsParameters ensureSqsparameters() => $_ensure(3);

  @$pb.TagNumber(223809421)
  $core.String get inputpath => $_getSZ(4);
  @$pb.TagNumber(223809421)
  set inputpath($core.String value) => $_setString(4, value);
  @$pb.TagNumber(223809421)
  $core.bool hasInputpath() => $_has(4);
  @$pb.TagNumber(223809421)
  void clearInputpath() => $_clearField(223809421);

  @$pb.TagNumber(242480292)
  BatchParameters get batchparameters => $_getN(5);
  @$pb.TagNumber(242480292)
  set batchparameters(BatchParameters value) => $_setField(242480292, value);
  @$pb.TagNumber(242480292)
  $core.bool hasBatchparameters() => $_has(5);
  @$pb.TagNumber(242480292)
  void clearBatchparameters() => $_clearField(242480292);
  @$pb.TagNumber(242480292)
  BatchParameters ensureBatchparameters() => $_ensure(5);

  @$pb.TagNumber(254855317)
  RedshiftDataParameters get redshiftdataparameters => $_getN(6);
  @$pb.TagNumber(254855317)
  set redshiftdataparameters(RedshiftDataParameters value) =>
      $_setField(254855317, value);
  @$pb.TagNumber(254855317)
  $core.bool hasRedshiftdataparameters() => $_has(6);
  @$pb.TagNumber(254855317)
  void clearRedshiftdataparameters() => $_clearField(254855317);
  @$pb.TagNumber(254855317)
  RedshiftDataParameters ensureRedshiftdataparameters() => $_ensure(6);

  @$pb.TagNumber(266827188)
  RetryPolicy get retrypolicy => $_getN(7);
  @$pb.TagNumber(266827188)
  set retrypolicy(RetryPolicy value) => $_setField(266827188, value);
  @$pb.TagNumber(266827188)
  $core.bool hasRetrypolicy() => $_has(7);
  @$pb.TagNumber(266827188)
  void clearRetrypolicy() => $_clearField(266827188);
  @$pb.TagNumber(266827188)
  RetryPolicy ensureRetrypolicy() => $_ensure(7);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(8);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(8);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(351941188)
  RunCommandParameters get runcommandparameters => $_getN(9);
  @$pb.TagNumber(351941188)
  set runcommandparameters(RunCommandParameters value) =>
      $_setField(351941188, value);
  @$pb.TagNumber(351941188)
  $core.bool hasRuncommandparameters() => $_has(9);
  @$pb.TagNumber(351941188)
  void clearRuncommandparameters() => $_clearField(351941188);
  @$pb.TagNumber(351941188)
  RunCommandParameters ensureRuncommandparameters() => $_ensure(9);

  @$pb.TagNumber(367051514)
  HttpParameters get httpparameters => $_getN(10);
  @$pb.TagNumber(367051514)
  set httpparameters(HttpParameters value) => $_setField(367051514, value);
  @$pb.TagNumber(367051514)
  $core.bool hasHttpparameters() => $_has(10);
  @$pb.TagNumber(367051514)
  void clearHttpparameters() => $_clearField(367051514);
  @$pb.TagNumber(367051514)
  HttpParameters ensureHttpparameters() => $_ensure(10);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(11);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(11, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(11);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(12);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(12, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(12);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(501521183)
  EcsParameters get ecsparameters => $_getN(13);
  @$pb.TagNumber(501521183)
  set ecsparameters(EcsParameters value) => $_setField(501521183, value);
  @$pb.TagNumber(501521183)
  $core.bool hasEcsparameters() => $_has(13);
  @$pb.TagNumber(501521183)
  void clearEcsparameters() => $_clearField(501521183);
  @$pb.TagNumber(501521183)
  EcsParameters ensureEcsparameters() => $_ensure(13);

  @$pb.TagNumber(513800606)
  SageMakerPipelineParameters get sagemakerpipelineparameters => $_getN(14);
  @$pb.TagNumber(513800606)
  set sagemakerpipelineparameters(SageMakerPipelineParameters value) =>
      $_setField(513800606, value);
  @$pb.TagNumber(513800606)
  $core.bool hasSagemakerpipelineparameters() => $_has(14);
  @$pb.TagNumber(513800606)
  void clearSagemakerpipelineparameters() => $_clearField(513800606);
  @$pb.TagNumber(513800606)
  SageMakerPipelineParameters ensureSagemakerpipelineparameters() =>
      $_ensure(14);

  @$pb.TagNumber(529785116)
  $core.String get input => $_getSZ(15);
  @$pb.TagNumber(529785116)
  set input($core.String value) => $_setString(15, value);
  @$pb.TagNumber(529785116)
  $core.bool hasInput() => $_has(15);
  @$pb.TagNumber(529785116)
  void clearInput() => $_clearField(529785116);
}

class TestEventPatternRequest extends $pb.GeneratedMessage {
  factory TestEventPatternRequest({
    $core.String? eventpattern,
    $core.String? event,
  }) {
    final result = create();
    if (eventpattern != null) result.eventpattern = eventpattern;
    if (event != null) result.event = event;
    return result;
  }

  TestEventPatternRequest._();

  factory TestEventPatternRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestEventPatternRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestEventPatternRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(233487232, _omitFieldNames ? '' : 'eventpattern')
    ..aOS(271274816, _omitFieldNames ? '' : 'event')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestEventPatternRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestEventPatternRequest copyWith(
          void Function(TestEventPatternRequest) updates) =>
      super.copyWith((message) => updates(message as TestEventPatternRequest))
          as TestEventPatternRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestEventPatternRequest create() => TestEventPatternRequest._();
  @$core.override
  TestEventPatternRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestEventPatternRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestEventPatternRequest>(create);
  static TestEventPatternRequest? _defaultInstance;

  @$pb.TagNumber(233487232)
  $core.String get eventpattern => $_getSZ(0);
  @$pb.TagNumber(233487232)
  set eventpattern($core.String value) => $_setString(0, value);
  @$pb.TagNumber(233487232)
  $core.bool hasEventpattern() => $_has(0);
  @$pb.TagNumber(233487232)
  void clearEventpattern() => $_clearField(233487232);

  @$pb.TagNumber(271274816)
  $core.String get event => $_getSZ(1);
  @$pb.TagNumber(271274816)
  set event($core.String value) => $_setString(1, value);
  @$pb.TagNumber(271274816)
  $core.bool hasEvent() => $_has(1);
  @$pb.TagNumber(271274816)
  void clearEvent() => $_clearField(271274816);
}

class TestEventPatternResponse extends $pb.GeneratedMessage {
  factory TestEventPatternResponse({
    $core.bool? result,
  }) {
    final result$ = create();
    if (result != null) result$.result = result;
    return result$;
  }

  TestEventPatternResponse._();

  factory TestEventPatternResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestEventPatternResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestEventPatternResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOB(273346629, _omitFieldNames ? '' : 'result')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestEventPatternResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestEventPatternResponse copyWith(
          void Function(TestEventPatternResponse) updates) =>
      super.copyWith((message) => updates(message as TestEventPatternResponse))
          as TestEventPatternResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestEventPatternResponse create() => TestEventPatternResponse._();
  @$core.override
  TestEventPatternResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestEventPatternResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestEventPatternResponse>(create);
  static TestEventPatternResponse? _defaultInstance;

  @$pb.TagNumber(273346629)
  $core.bool get result => $_getBF(0);
  @$pb.TagNumber(273346629)
  set result($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(273346629)
  $core.bool hasResult() => $_has(0);
  @$pb.TagNumber(273346629)
  void clearResult() => $_clearField(273346629);
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
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

class UpdateApiDestinationRequest extends $pb.GeneratedMessage {
  factory UpdateApiDestinationRequest({
    $core.String? description,
    $core.String? connectionarn,
    $core.String? name,
    $core.int? invocationratelimitpersecond,
    ApiDestinationHttpMethod? httpmethod,
    $core.String? invocationendpoint,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (name != null) result.name = name;
    if (invocationratelimitpersecond != null)
      result.invocationratelimitpersecond = invocationratelimitpersecond;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (invocationendpoint != null)
      result.invocationendpoint = invocationendpoint;
    return result;
  }

  UpdateApiDestinationRequest._();

  factory UpdateApiDestinationRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateApiDestinationRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateApiDestinationRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aI(295327816, _omitFieldNames ? '' : 'invocationratelimitpersecond')
    ..aE<ApiDestinationHttpMethod>(
        398394961, _omitFieldNames ? '' : 'httpmethod',
        enumValues: ApiDestinationHttpMethod.values)
    ..aOS(411800759, _omitFieldNames ? '' : 'invocationendpoint')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateApiDestinationRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateApiDestinationRequest copyWith(
          void Function(UpdateApiDestinationRequest) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateApiDestinationRequest))
          as UpdateApiDestinationRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateApiDestinationRequest create() =>
      UpdateApiDestinationRequest._();
  @$core.override
  UpdateApiDestinationRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateApiDestinationRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateApiDestinationRequest>(create);
  static UpdateApiDestinationRequest? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(1);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(1);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(295327816)
  $core.int get invocationratelimitpersecond => $_getIZ(3);
  @$pb.TagNumber(295327816)
  set invocationratelimitpersecond($core.int value) =>
      $_setSignedInt32(3, value);
  @$pb.TagNumber(295327816)
  $core.bool hasInvocationratelimitpersecond() => $_has(3);
  @$pb.TagNumber(295327816)
  void clearInvocationratelimitpersecond() => $_clearField(295327816);

  @$pb.TagNumber(398394961)
  ApiDestinationHttpMethod get httpmethod => $_getN(4);
  @$pb.TagNumber(398394961)
  set httpmethod(ApiDestinationHttpMethod value) =>
      $_setField(398394961, value);
  @$pb.TagNumber(398394961)
  $core.bool hasHttpmethod() => $_has(4);
  @$pb.TagNumber(398394961)
  void clearHttpmethod() => $_clearField(398394961);

  @$pb.TagNumber(411800759)
  $core.String get invocationendpoint => $_getSZ(5);
  @$pb.TagNumber(411800759)
  set invocationendpoint($core.String value) => $_setString(5, value);
  @$pb.TagNumber(411800759)
  $core.bool hasInvocationendpoint() => $_has(5);
  @$pb.TagNumber(411800759)
  void clearInvocationendpoint() => $_clearField(411800759);
}

class UpdateApiDestinationResponse extends $pb.GeneratedMessage {
  factory UpdateApiDestinationResponse({
    ApiDestinationState? apidestinationstate,
    $core.String? apidestinationarn,
    $core.String? creationtime,
    $core.String? lastmodifiedtime,
  }) {
    final result = create();
    if (apidestinationstate != null)
      result.apidestinationstate = apidestinationstate;
    if (apidestinationarn != null) result.apidestinationarn = apidestinationarn;
    if (creationtime != null) result.creationtime = creationtime;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    return result;
  }

  UpdateApiDestinationResponse._();

  factory UpdateApiDestinationResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateApiDestinationResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateApiDestinationResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aE<ApiDestinationState>(
        13153343, _omitFieldNames ? '' : 'apidestinationstate',
        enumValues: ApiDestinationState.values)
    ..aOS(90996885, _omitFieldNames ? '' : 'apidestinationarn')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateApiDestinationResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateApiDestinationResponse copyWith(
          void Function(UpdateApiDestinationResponse) updates) =>
      super.copyWith(
              (message) => updates(message as UpdateApiDestinationResponse))
          as UpdateApiDestinationResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateApiDestinationResponse create() =>
      UpdateApiDestinationResponse._();
  @$core.override
  UpdateApiDestinationResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateApiDestinationResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateApiDestinationResponse>(create);
  static UpdateApiDestinationResponse? _defaultInstance;

  @$pb.TagNumber(13153343)
  ApiDestinationState get apidestinationstate => $_getN(0);
  @$pb.TagNumber(13153343)
  set apidestinationstate(ApiDestinationState value) =>
      $_setField(13153343, value);
  @$pb.TagNumber(13153343)
  $core.bool hasApidestinationstate() => $_has(0);
  @$pb.TagNumber(13153343)
  void clearApidestinationstate() => $_clearField(13153343);

  @$pb.TagNumber(90996885)
  $core.String get apidestinationarn => $_getSZ(1);
  @$pb.TagNumber(90996885)
  set apidestinationarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(90996885)
  $core.bool hasApidestinationarn() => $_has(1);
  @$pb.TagNumber(90996885)
  void clearApidestinationarn() => $_clearField(90996885);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(2);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(2);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(3);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(3);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);
}

class UpdateArchiveRequest extends $pb.GeneratedMessage {
  factory UpdateArchiveRequest({
    $core.String? archivename,
    $core.String? description,
    $core.String? eventpattern,
    $core.int? retentiondays,
  }) {
    final result = create();
    if (archivename != null) result.archivename = archivename;
    if (description != null) result.description = description;
    if (eventpattern != null) result.eventpattern = eventpattern;
    if (retentiondays != null) result.retentiondays = retentiondays;
    return result;
  }

  UpdateArchiveRequest._();

  factory UpdateArchiveRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateArchiveRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateArchiveRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(88048487, _omitFieldNames ? '' : 'archivename')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(233487232, _omitFieldNames ? '' : 'eventpattern')
    ..aI(267894223, _omitFieldNames ? '' : 'retentiondays')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateArchiveRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateArchiveRequest copyWith(void Function(UpdateArchiveRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateArchiveRequest))
          as UpdateArchiveRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateArchiveRequest create() => UpdateArchiveRequest._();
  @$core.override
  UpdateArchiveRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateArchiveRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateArchiveRequest>(create);
  static UpdateArchiveRequest? _defaultInstance;

  @$pb.TagNumber(88048487)
  $core.String get archivename => $_getSZ(0);
  @$pb.TagNumber(88048487)
  set archivename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(88048487)
  $core.bool hasArchivename() => $_has(0);
  @$pb.TagNumber(88048487)
  void clearArchivename() => $_clearField(88048487);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(1);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(1, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(1);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(233487232)
  $core.String get eventpattern => $_getSZ(2);
  @$pb.TagNumber(233487232)
  set eventpattern($core.String value) => $_setString(2, value);
  @$pb.TagNumber(233487232)
  $core.bool hasEventpattern() => $_has(2);
  @$pb.TagNumber(233487232)
  void clearEventpattern() => $_clearField(233487232);

  @$pb.TagNumber(267894223)
  $core.int get retentiondays => $_getIZ(3);
  @$pb.TagNumber(267894223)
  set retentiondays($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(267894223)
  $core.bool hasRetentiondays() => $_has(3);
  @$pb.TagNumber(267894223)
  void clearRetentiondays() => $_clearField(267894223);
}

class UpdateArchiveResponse extends $pb.GeneratedMessage {
  factory UpdateArchiveResponse({
    $core.String? archivearn,
    $core.String? creationtime,
    $core.String? statereason,
    ArchiveState? state,
  }) {
    final result = create();
    if (archivearn != null) result.archivearn = archivearn;
    if (creationtime != null) result.creationtime = creationtime;
    if (statereason != null) result.statereason = statereason;
    if (state != null) result.state = state;
    return result;
  }

  UpdateArchiveResponse._();

  factory UpdateArchiveResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateArchiveResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateArchiveResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(56866685, _omitFieldNames ? '' : 'archivearn')
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..aE<ArchiveState>(502047895, _omitFieldNames ? '' : 'state',
        enumValues: ArchiveState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateArchiveResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateArchiveResponse copyWith(
          void Function(UpdateArchiveResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateArchiveResponse))
          as UpdateArchiveResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateArchiveResponse create() => UpdateArchiveResponse._();
  @$core.override
  UpdateArchiveResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateArchiveResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateArchiveResponse>(create);
  static UpdateArchiveResponse? _defaultInstance;

  @$pb.TagNumber(56866685)
  $core.String get archivearn => $_getSZ(0);
  @$pb.TagNumber(56866685)
  set archivearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(56866685)
  $core.bool hasArchivearn() => $_has(0);
  @$pb.TagNumber(56866685)
  void clearArchivearn() => $_clearField(56866685);

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(1);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(1);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(2);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(2, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(2);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(502047895)
  ArchiveState get state => $_getN(3);
  @$pb.TagNumber(502047895)
  set state(ArchiveState value) => $_setField(502047895, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(3);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class UpdateConnectionApiKeyAuthRequestParameters extends $pb.GeneratedMessage {
  factory UpdateConnectionApiKeyAuthRequestParameters({
    $core.String? apikeyname,
    $core.String? apikeyvalue,
  }) {
    final result = create();
    if (apikeyname != null) result.apikeyname = apikeyname;
    if (apikeyvalue != null) result.apikeyvalue = apikeyvalue;
    return result;
  }

  UpdateConnectionApiKeyAuthRequestParameters._();

  factory UpdateConnectionApiKeyAuthRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateConnectionApiKeyAuthRequestParameters.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateConnectionApiKeyAuthRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(43133502, _omitFieldNames ? '' : 'apikeyname')
    ..aOS(93786144, _omitFieldNames ? '' : 'apikeyvalue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionApiKeyAuthRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionApiKeyAuthRequestParameters copyWith(
          void Function(UpdateConnectionApiKeyAuthRequestParameters) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateConnectionApiKeyAuthRequestParameters))
          as UpdateConnectionApiKeyAuthRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateConnectionApiKeyAuthRequestParameters create() =>
      UpdateConnectionApiKeyAuthRequestParameters._();
  @$core.override
  UpdateConnectionApiKeyAuthRequestParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateConnectionApiKeyAuthRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateConnectionApiKeyAuthRequestParameters>(create);
  static UpdateConnectionApiKeyAuthRequestParameters? _defaultInstance;

  @$pb.TagNumber(43133502)
  $core.String get apikeyname => $_getSZ(0);
  @$pb.TagNumber(43133502)
  set apikeyname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(43133502)
  $core.bool hasApikeyname() => $_has(0);
  @$pb.TagNumber(43133502)
  void clearApikeyname() => $_clearField(43133502);

  @$pb.TagNumber(93786144)
  $core.String get apikeyvalue => $_getSZ(1);
  @$pb.TagNumber(93786144)
  set apikeyvalue($core.String value) => $_setString(1, value);
  @$pb.TagNumber(93786144)
  $core.bool hasApikeyvalue() => $_has(1);
  @$pb.TagNumber(93786144)
  void clearApikeyvalue() => $_clearField(93786144);
}

class UpdateConnectionAuthRequestParameters extends $pb.GeneratedMessage {
  factory UpdateConnectionAuthRequestParameters({
    UpdateConnectionBasicAuthRequestParameters? basicauthparameters,
    UpdateConnectionApiKeyAuthRequestParameters? apikeyauthparameters,
    UpdateConnectionOAuthRequestParameters? oauthparameters,
    ConnectionHttpParameters? invocationhttpparameters,
  }) {
    final result = create();
    if (basicauthparameters != null)
      result.basicauthparameters = basicauthparameters;
    if (apikeyauthparameters != null)
      result.apikeyauthparameters = apikeyauthparameters;
    if (oauthparameters != null) result.oauthparameters = oauthparameters;
    if (invocationhttpparameters != null)
      result.invocationhttpparameters = invocationhttpparameters;
    return result;
  }

  UpdateConnectionAuthRequestParameters._();

  factory UpdateConnectionAuthRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateConnectionAuthRequestParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateConnectionAuthRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOM<UpdateConnectionBasicAuthRequestParameters>(
        63965312, _omitFieldNames ? '' : 'basicauthparameters',
        subBuilder: UpdateConnectionBasicAuthRequestParameters.create)
    ..aOM<UpdateConnectionApiKeyAuthRequestParameters>(
        110622489, _omitFieldNames ? '' : 'apikeyauthparameters',
        subBuilder: UpdateConnectionApiKeyAuthRequestParameters.create)
    ..aOM<UpdateConnectionOAuthRequestParameters>(
        206836569, _omitFieldNames ? '' : 'oauthparameters',
        subBuilder: UpdateConnectionOAuthRequestParameters.create)
    ..aOM<ConnectionHttpParameters>(
        499128572, _omitFieldNames ? '' : 'invocationhttpparameters',
        subBuilder: ConnectionHttpParameters.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionAuthRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionAuthRequestParameters copyWith(
          void Function(UpdateConnectionAuthRequestParameters) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateConnectionAuthRequestParameters))
          as UpdateConnectionAuthRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateConnectionAuthRequestParameters create() =>
      UpdateConnectionAuthRequestParameters._();
  @$core.override
  UpdateConnectionAuthRequestParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateConnectionAuthRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateConnectionAuthRequestParameters>(create);
  static UpdateConnectionAuthRequestParameters? _defaultInstance;

  @$pb.TagNumber(63965312)
  UpdateConnectionBasicAuthRequestParameters get basicauthparameters =>
      $_getN(0);
  @$pb.TagNumber(63965312)
  set basicauthparameters(UpdateConnectionBasicAuthRequestParameters value) =>
      $_setField(63965312, value);
  @$pb.TagNumber(63965312)
  $core.bool hasBasicauthparameters() => $_has(0);
  @$pb.TagNumber(63965312)
  void clearBasicauthparameters() => $_clearField(63965312);
  @$pb.TagNumber(63965312)
  UpdateConnectionBasicAuthRequestParameters ensureBasicauthparameters() =>
      $_ensure(0);

  @$pb.TagNumber(110622489)
  UpdateConnectionApiKeyAuthRequestParameters get apikeyauthparameters =>
      $_getN(1);
  @$pb.TagNumber(110622489)
  set apikeyauthparameters(UpdateConnectionApiKeyAuthRequestParameters value) =>
      $_setField(110622489, value);
  @$pb.TagNumber(110622489)
  $core.bool hasApikeyauthparameters() => $_has(1);
  @$pb.TagNumber(110622489)
  void clearApikeyauthparameters() => $_clearField(110622489);
  @$pb.TagNumber(110622489)
  UpdateConnectionApiKeyAuthRequestParameters ensureApikeyauthparameters() =>
      $_ensure(1);

  @$pb.TagNumber(206836569)
  UpdateConnectionOAuthRequestParameters get oauthparameters => $_getN(2);
  @$pb.TagNumber(206836569)
  set oauthparameters(UpdateConnectionOAuthRequestParameters value) =>
      $_setField(206836569, value);
  @$pb.TagNumber(206836569)
  $core.bool hasOauthparameters() => $_has(2);
  @$pb.TagNumber(206836569)
  void clearOauthparameters() => $_clearField(206836569);
  @$pb.TagNumber(206836569)
  UpdateConnectionOAuthRequestParameters ensureOauthparameters() => $_ensure(2);

  @$pb.TagNumber(499128572)
  ConnectionHttpParameters get invocationhttpparameters => $_getN(3);
  @$pb.TagNumber(499128572)
  set invocationhttpparameters(ConnectionHttpParameters value) =>
      $_setField(499128572, value);
  @$pb.TagNumber(499128572)
  $core.bool hasInvocationhttpparameters() => $_has(3);
  @$pb.TagNumber(499128572)
  void clearInvocationhttpparameters() => $_clearField(499128572);
  @$pb.TagNumber(499128572)
  ConnectionHttpParameters ensureInvocationhttpparameters() => $_ensure(3);
}

class UpdateConnectionBasicAuthRequestParameters extends $pb.GeneratedMessage {
  factory UpdateConnectionBasicAuthRequestParameters({
    $core.String? password,
    $core.String? username,
  }) {
    final result = create();
    if (password != null) result.password = password;
    if (username != null) result.username = username;
    return result;
  }

  UpdateConnectionBasicAuthRequestParameters._();

  factory UpdateConnectionBasicAuthRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateConnectionBasicAuthRequestParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateConnectionBasicAuthRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(214108217, _omitFieldNames ? '' : 'password')
    ..aOS(470340826, _omitFieldNames ? '' : 'username')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionBasicAuthRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionBasicAuthRequestParameters copyWith(
          void Function(UpdateConnectionBasicAuthRequestParameters) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateConnectionBasicAuthRequestParameters))
          as UpdateConnectionBasicAuthRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateConnectionBasicAuthRequestParameters create() =>
      UpdateConnectionBasicAuthRequestParameters._();
  @$core.override
  UpdateConnectionBasicAuthRequestParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateConnectionBasicAuthRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateConnectionBasicAuthRequestParameters>(create);
  static UpdateConnectionBasicAuthRequestParameters? _defaultInstance;

  @$pb.TagNumber(214108217)
  $core.String get password => $_getSZ(0);
  @$pb.TagNumber(214108217)
  set password($core.String value) => $_setString(0, value);
  @$pb.TagNumber(214108217)
  $core.bool hasPassword() => $_has(0);
  @$pb.TagNumber(214108217)
  void clearPassword() => $_clearField(214108217);

  @$pb.TagNumber(470340826)
  $core.String get username => $_getSZ(1);
  @$pb.TagNumber(470340826)
  set username($core.String value) => $_setString(1, value);
  @$pb.TagNumber(470340826)
  $core.bool hasUsername() => $_has(1);
  @$pb.TagNumber(470340826)
  void clearUsername() => $_clearField(470340826);
}

class UpdateConnectionOAuthClientRequestParameters
    extends $pb.GeneratedMessage {
  factory UpdateConnectionOAuthClientRequestParameters({
    $core.String? clientid,
    $core.String? clientsecret,
  }) {
    final result = create();
    if (clientid != null) result.clientid = clientid;
    if (clientsecret != null) result.clientsecret = clientsecret;
    return result;
  }

  UpdateConnectionOAuthClientRequestParameters._();

  factory UpdateConnectionOAuthClientRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateConnectionOAuthClientRequestParameters.fromJson(
          $core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateConnectionOAuthClientRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(448889284, _omitFieldNames ? '' : 'clientid')
    ..aOS(500734711, _omitFieldNames ? '' : 'clientsecret')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionOAuthClientRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionOAuthClientRequestParameters copyWith(
          void Function(UpdateConnectionOAuthClientRequestParameters)
              updates) =>
      super.copyWith((message) =>
              updates(message as UpdateConnectionOAuthClientRequestParameters))
          as UpdateConnectionOAuthClientRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateConnectionOAuthClientRequestParameters create() =>
      UpdateConnectionOAuthClientRequestParameters._();
  @$core.override
  UpdateConnectionOAuthClientRequestParameters createEmptyInstance() =>
      create();
  @$core.pragma('dart2js:noInline')
  static UpdateConnectionOAuthClientRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateConnectionOAuthClientRequestParameters>(create);
  static UpdateConnectionOAuthClientRequestParameters? _defaultInstance;

  @$pb.TagNumber(448889284)
  $core.String get clientid => $_getSZ(0);
  @$pb.TagNumber(448889284)
  set clientid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(448889284)
  $core.bool hasClientid() => $_has(0);
  @$pb.TagNumber(448889284)
  void clearClientid() => $_clearField(448889284);

  @$pb.TagNumber(500734711)
  $core.String get clientsecret => $_getSZ(1);
  @$pb.TagNumber(500734711)
  set clientsecret($core.String value) => $_setString(1, value);
  @$pb.TagNumber(500734711)
  $core.bool hasClientsecret() => $_has(1);
  @$pb.TagNumber(500734711)
  void clearClientsecret() => $_clearField(500734711);
}

class UpdateConnectionOAuthRequestParameters extends $pb.GeneratedMessage {
  factory UpdateConnectionOAuthRequestParameters({
    ConnectionHttpParameters? oauthhttpparameters,
    UpdateConnectionOAuthClientRequestParameters? clientparameters,
    ConnectionOAuthHttpMethod? httpmethod,
    $core.String? authorizationendpoint,
  }) {
    final result = create();
    if (oauthhttpparameters != null)
      result.oauthhttpparameters = oauthhttpparameters;
    if (clientparameters != null) result.clientparameters = clientparameters;
    if (httpmethod != null) result.httpmethod = httpmethod;
    if (authorizationendpoint != null)
      result.authorizationendpoint = authorizationendpoint;
    return result;
  }

  UpdateConnectionOAuthRequestParameters._();

  factory UpdateConnectionOAuthRequestParameters.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateConnectionOAuthRequestParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateConnectionOAuthRequestParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOM<ConnectionHttpParameters>(
        10294537, _omitFieldNames ? '' : 'oauthhttpparameters',
        subBuilder: ConnectionHttpParameters.create)
    ..aOM<UpdateConnectionOAuthClientRequestParameters>(
        246864127, _omitFieldNames ? '' : 'clientparameters',
        subBuilder: UpdateConnectionOAuthClientRequestParameters.create)
    ..aE<ConnectionOAuthHttpMethod>(
        398394961, _omitFieldNames ? '' : 'httpmethod',
        enumValues: ConnectionOAuthHttpMethod.values)
    ..aOS(427938596, _omitFieldNames ? '' : 'authorizationendpoint')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionOAuthRequestParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionOAuthRequestParameters copyWith(
          void Function(UpdateConnectionOAuthRequestParameters) updates) =>
      super.copyWith((message) =>
              updates(message as UpdateConnectionOAuthRequestParameters))
          as UpdateConnectionOAuthRequestParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateConnectionOAuthRequestParameters create() =>
      UpdateConnectionOAuthRequestParameters._();
  @$core.override
  UpdateConnectionOAuthRequestParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateConnectionOAuthRequestParameters getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          UpdateConnectionOAuthRequestParameters>(create);
  static UpdateConnectionOAuthRequestParameters? _defaultInstance;

  @$pb.TagNumber(10294537)
  ConnectionHttpParameters get oauthhttpparameters => $_getN(0);
  @$pb.TagNumber(10294537)
  set oauthhttpparameters(ConnectionHttpParameters value) =>
      $_setField(10294537, value);
  @$pb.TagNumber(10294537)
  $core.bool hasOauthhttpparameters() => $_has(0);
  @$pb.TagNumber(10294537)
  void clearOauthhttpparameters() => $_clearField(10294537);
  @$pb.TagNumber(10294537)
  ConnectionHttpParameters ensureOauthhttpparameters() => $_ensure(0);

  @$pb.TagNumber(246864127)
  UpdateConnectionOAuthClientRequestParameters get clientparameters =>
      $_getN(1);
  @$pb.TagNumber(246864127)
  set clientparameters(UpdateConnectionOAuthClientRequestParameters value) =>
      $_setField(246864127, value);
  @$pb.TagNumber(246864127)
  $core.bool hasClientparameters() => $_has(1);
  @$pb.TagNumber(246864127)
  void clearClientparameters() => $_clearField(246864127);
  @$pb.TagNumber(246864127)
  UpdateConnectionOAuthClientRequestParameters ensureClientparameters() =>
      $_ensure(1);

  @$pb.TagNumber(398394961)
  ConnectionOAuthHttpMethod get httpmethod => $_getN(2);
  @$pb.TagNumber(398394961)
  set httpmethod(ConnectionOAuthHttpMethod value) =>
      $_setField(398394961, value);
  @$pb.TagNumber(398394961)
  $core.bool hasHttpmethod() => $_has(2);
  @$pb.TagNumber(398394961)
  void clearHttpmethod() => $_clearField(398394961);

  @$pb.TagNumber(427938596)
  $core.String get authorizationendpoint => $_getSZ(3);
  @$pb.TagNumber(427938596)
  set authorizationendpoint($core.String value) => $_setString(3, value);
  @$pb.TagNumber(427938596)
  $core.bool hasAuthorizationendpoint() => $_has(3);
  @$pb.TagNumber(427938596)
  void clearAuthorizationendpoint() => $_clearField(427938596);
}

class UpdateConnectionRequest extends $pb.GeneratedMessage {
  factory UpdateConnectionRequest({
    $core.String? description,
    UpdateConnectionAuthRequestParameters? authparameters,
    $core.String? name,
    ConnectionAuthorizationType? authorizationtype,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (authparameters != null) result.authparameters = authparameters;
    if (name != null) result.name = name;
    if (authorizationtype != null) result.authorizationtype = authorizationtype;
    return result;
  }

  UpdateConnectionRequest._();

  factory UpdateConnectionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateConnectionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateConnectionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOM<UpdateConnectionAuthRequestParameters>(
        258276552, _omitFieldNames ? '' : 'authparameters',
        subBuilder: UpdateConnectionAuthRequestParameters.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aE<ConnectionAuthorizationType>(
        481976511, _omitFieldNames ? '' : 'authorizationtype',
        enumValues: ConnectionAuthorizationType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionRequest copyWith(
          void Function(UpdateConnectionRequest) updates) =>
      super.copyWith((message) => updates(message as UpdateConnectionRequest))
          as UpdateConnectionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateConnectionRequest create() => UpdateConnectionRequest._();
  @$core.override
  UpdateConnectionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateConnectionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateConnectionRequest>(create);
  static UpdateConnectionRequest? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(258276552)
  UpdateConnectionAuthRequestParameters get authparameters => $_getN(1);
  @$pb.TagNumber(258276552)
  set authparameters(UpdateConnectionAuthRequestParameters value) =>
      $_setField(258276552, value);
  @$pb.TagNumber(258276552)
  $core.bool hasAuthparameters() => $_has(1);
  @$pb.TagNumber(258276552)
  void clearAuthparameters() => $_clearField(258276552);
  @$pb.TagNumber(258276552)
  UpdateConnectionAuthRequestParameters ensureAuthparameters() => $_ensure(1);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(481976511)
  ConnectionAuthorizationType get authorizationtype => $_getN(3);
  @$pb.TagNumber(481976511)
  set authorizationtype(ConnectionAuthorizationType value) =>
      $_setField(481976511, value);
  @$pb.TagNumber(481976511)
  $core.bool hasAuthorizationtype() => $_has(3);
  @$pb.TagNumber(481976511)
  void clearAuthorizationtype() => $_clearField(481976511);
}

class UpdateConnectionResponse extends $pb.GeneratedMessage {
  factory UpdateConnectionResponse({
    $core.String? creationtime,
    $core.String? connectionarn,
    $core.String? lastmodifiedtime,
    $core.String? lastauthorizedtime,
    ConnectionState? connectionstate,
  }) {
    final result = create();
    if (creationtime != null) result.creationtime = creationtime;
    if (connectionarn != null) result.connectionarn = connectionarn;
    if (lastmodifiedtime != null) result.lastmodifiedtime = lastmodifiedtime;
    if (lastauthorizedtime != null)
      result.lastauthorizedtime = lastauthorizedtime;
    if (connectionstate != null) result.connectionstate = connectionstate;
    return result;
  }

  UpdateConnectionResponse._();

  factory UpdateConnectionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UpdateConnectionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UpdateConnectionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'events'),
      createEmptyInstance: create)
    ..aOS(103458790, _omitFieldNames ? '' : 'creationtime')
    ..aOS(187631553, _omitFieldNames ? '' : 'connectionarn')
    ..aOS(236912992, _omitFieldNames ? '' : 'lastmodifiedtime')
    ..aOS(250638066, _omitFieldNames ? '' : 'lastauthorizedtime')
    ..aE<ConnectionState>(404323675, _omitFieldNames ? '' : 'connectionstate',
        enumValues: ConnectionState.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UpdateConnectionResponse copyWith(
          void Function(UpdateConnectionResponse) updates) =>
      super.copyWith((message) => updates(message as UpdateConnectionResponse))
          as UpdateConnectionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UpdateConnectionResponse create() => UpdateConnectionResponse._();
  @$core.override
  UpdateConnectionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UpdateConnectionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UpdateConnectionResponse>(create);
  static UpdateConnectionResponse? _defaultInstance;

  @$pb.TagNumber(103458790)
  $core.String get creationtime => $_getSZ(0);
  @$pb.TagNumber(103458790)
  set creationtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(103458790)
  $core.bool hasCreationtime() => $_has(0);
  @$pb.TagNumber(103458790)
  void clearCreationtime() => $_clearField(103458790);

  @$pb.TagNumber(187631553)
  $core.String get connectionarn => $_getSZ(1);
  @$pb.TagNumber(187631553)
  set connectionarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(187631553)
  $core.bool hasConnectionarn() => $_has(1);
  @$pb.TagNumber(187631553)
  void clearConnectionarn() => $_clearField(187631553);

  @$pb.TagNumber(236912992)
  $core.String get lastmodifiedtime => $_getSZ(2);
  @$pb.TagNumber(236912992)
  set lastmodifiedtime($core.String value) => $_setString(2, value);
  @$pb.TagNumber(236912992)
  $core.bool hasLastmodifiedtime() => $_has(2);
  @$pb.TagNumber(236912992)
  void clearLastmodifiedtime() => $_clearField(236912992);

  @$pb.TagNumber(250638066)
  $core.String get lastauthorizedtime => $_getSZ(3);
  @$pb.TagNumber(250638066)
  set lastauthorizedtime($core.String value) => $_setString(3, value);
  @$pb.TagNumber(250638066)
  $core.bool hasLastauthorizedtime() => $_has(3);
  @$pb.TagNumber(250638066)
  void clearLastauthorizedtime() => $_clearField(250638066);

  @$pb.TagNumber(404323675)
  ConnectionState get connectionstate => $_getN(4);
  @$pb.TagNumber(404323675)
  set connectionstate(ConnectionState value) => $_setField(404323675, value);
  @$pb.TagNumber(404323675)
  $core.bool hasConnectionstate() => $_has(4);
  @$pb.TagNumber(404323675)
  void clearConnectionstate() => $_clearField(404323675);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
